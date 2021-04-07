package spotinst

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/spotinst/spotinst-sdk-go/service/ocean/providers/azure"
	"github.com/spotinst/spotinst-sdk-go/spotinst"
	"github.com/spotinst/spotinst-sdk-go/spotinst/client"
	"github.com/spotinst/terraform-provider-spotinst/spotinst/commons"
	"github.com/spotinst/terraform-provider-spotinst/spotinst/ocean_aks"
	"github.com/spotinst/terraform-provider-spotinst/spotinst/ocean_aks_login"
)

func resourceSpotinstOceanAKS() *schema.Resource {
	setupClusterAKSResource()

	return &schema.Resource{
		Create: resourceSpotinstClusterAKSCreate,
		Read:   resourceSpotinstClusterAKSRead,
		Update: resourceSpotinstClusterAKSUpdate,
		Delete: resourceSpotinstClusterAKSDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: commons.OceanAKSResource.GetSchemaMap(),
	}
}

func setupClusterAKSResource() {
	fieldsMap := make(map[commons.FieldName]*commons.GenericField)

	ocean_aks.Setup(fieldsMap)
	ocean_aks_login.Setup(fieldsMap)

	commons.OceanAKSResource = commons.NewOceanAKSResource(fieldsMap)
}

// region Create

func resourceSpotinstClusterAKSCreate(resourceData *schema.ResourceData, meta interface{}) error {
	log.Printf(string(commons.ResourceOnCreate), commons.OceanAKSResource.GetName())

	importedCluster, err := importAKSCluster(resourceData, meta.(*Client))
	if err != nil {
		return err
	}

	cluster, err := commons.OceanAKSResource.OnCreate(importedCluster, resourceData, meta)
	if err != nil {
		return err
	}

	clusterID, err := createAKSCluster(cluster, meta.(*Client))
	if err != nil {
		return err
	}

	resourceData.SetId(spotinst.StringValue(clusterID))
	log.Printf("ocean/aks: AKS cluster created successfully: %s", resourceData.Id())

	return resourceSpotinstClusterAKSRead(resourceData, meta)
}

func createAKSCluster(cluster *azure.Cluster, spotinstClient *Client) (*string, error) {
	if json, err := commons.ToJson(cluster); err != nil {
		return nil, err
	} else {
		log.Printf("ocean/aks: cluster configuration: %s", json)
	}

	input := &azure.CreateClusterInput{
		Cluster: cluster,
	}

	var output *azure.CreateClusterOutput
	err := resource.Retry(time.Minute, func() *resource.RetryError {
		o, err := spotinstClient.ocean.CloudProviderAzure().CreateCluster(context.TODO(), input)
		if err != nil {
			// Some other error, report it.
			return resource.NonRetryableError(err)
		}
		output = o
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("ocean/aks: failed to create cluster: %v", err)
	}

	return output.Cluster.ID, nil
}

// endregion

// region Read

func resourceSpotinstClusterAKSRead(resourceData *schema.ResourceData, meta interface{}) error {
	clusterID := resourceData.Id()
	log.Printf(string(commons.ResourceOnRead), commons.OceanAKSResource.GetName(), clusterID)

	cluster, err := readAKSCluster(context.TODO(), clusterID, meta.(*Client))
	if err != nil {
		return err
	}

	// If nothing was found, return no state.
	if cluster == nil {
		resourceData.SetId("")
		return nil
	}

	// Expose the controller cluster identifier.
	if cluster.ControllerClusterID != nil {
		_ = resourceData.Set(string(ocean_aks.ControllerClusterID),
			spotinst.StringValue(cluster.ControllerClusterID))
	}

	if err := commons.OceanAKSResource.OnRead(cluster, resourceData, meta); err != nil {
		return err
	}

	log.Printf("ocean/aks: cluster read successfully: %s", clusterID)
	return nil
}

func readAKSCluster(ctx context.Context, clusterID string, spotinstClient *Client) (*azure.Cluster, error) {
	input := &azure.ReadClusterInput{
		ClusterID: spotinst.String(clusterID),
	}

	output, err := spotinstClient.ocean.CloudProviderAzure().ReadCluster(ctx, input)
	if err != nil {
		// If the virtualNodeGroup was not found, return nil so that we can show
		// that it does not exist
		if errs, ok := err.(client.Errors); ok && len(errs) > 0 {
			for _, err := range errs {
				if err.Code == ErrCodeClusterNotFound {
					return nil, nil
				}
			}
		}

		// Some other error, report it.
		return nil, fmt.Errorf("ocean/aks: failed to read cluster: %v", err)
	}

	return output.Cluster, nil
}

// endregion

// region Update

func resourceSpotinstClusterAKSUpdate(resourceData *schema.ResourceData, meta interface{}) error {
	clusterID := resourceData.Id()
	log.Printf(string(commons.ResourceOnUpdate), commons.OceanAKSResource.GetName(), clusterID)

	shouldUpdate, cluster, err := commons.OceanAKSResource.OnUpdate(resourceData, meta)
	if err != nil {
		return err
	}

	if shouldUpdate {
		cluster.SetId(spotinst.String(clusterID))
		if err := updateAKSCluster(cluster, meta.(*Client)); err != nil {
			return err
		}
	}

	log.Printf("ocean/aks: cluster updated successfully: %s", clusterID)
	return resourceSpotinstClusterAKSRead(resourceData, meta)
}

func updateAKSCluster(cluster *azure.Cluster, spotinstClient *Client) error {
	input := &azure.UpdateClusterInput{
		Cluster: cluster,
	}

	if json, err := commons.ToJson(cluster); err != nil {
		return err
	} else {
		log.Printf("ocean/aks: cluster update configuration: %s", json)
	}

	if _, err := spotinstClient.ocean.CloudProviderAzure().UpdateCluster(context.TODO(), input); err != nil {
		return fmt.Errorf("ocean/aks: failed to update cluster: %v", err)
	}

	return nil
}

// endregion

// region Delete

func resourceSpotinstClusterAKSDelete(resourceData *schema.ResourceData, meta interface{}) error {
	clusterID := resourceData.Id()
	log.Printf(string(commons.ResourceOnDelete), commons.OceanAKSResource.GetName(), clusterID)

	if err := deleteAKSCluster(clusterID, meta.(*Client)); err != nil {
		return err
	}

	log.Printf("ocean/aks: cluster deleted successfully: %s", clusterID)
	resourceData.SetId("")

	return nil
}

func deleteAKSCluster(clusterID string, spotinstClient *Client) error {
	input := &azure.DeleteClusterInput{
		ClusterID: spotinst.String(clusterID),
	}

	if json, err := commons.ToJson(input); err != nil {
		return err
	} else {
		log.Printf("ocean/aks: cluster delete configuration: %s", json)
	}

	if _, err := spotinstClient.ocean.CloudProviderAzure().DeleteCluster(context.TODO(), input); err != nil {
		return fmt.Errorf("ocean/aks: failed to delete cluster: %v", err)
	}

	return nil
}

// endregion

// region Import

func importAKSCluster(resourceData *schema.ResourceData, spotinstClient *Client) (*azure.Cluster, error) {
	input := &azure.ImportClusterInput{
		ACDIdentifier: spotinst.String(resourceData.Get("acd_identifier").(string)),
		Cluster: &azure.ImportCluster{
			Name: spotinst.String(resourceData.Get("name").(string)),
			AKS: &azure.AKS{
				Name:              spotinst.String(resourceData.Get("aks_name").(string)),
				ResourceGroupName: spotinst.String(resourceData.Get("aks_resource_group_name").(string)),
			}},
	}

	output, err := spotinstClient.ocean.CloudProviderAzure().ImportCluster(context.TODO(), input)
	if err != nil {
		// If the cluster was not found, return nil so that we can show
		// that the cluster is gone.
		if errs, ok := err.(client.Errors); ok && len(errs) > 0 {
			for _, err := range errs {
				if err.Code == ErrCodeClusterNotFound {
					resourceData.SetId("")
					return nil, err
				}
			}
		}
		// Some other error, report it.
		return nil, fmt.Errorf("ocean/aks: failed to import cluster: %v", err)
	}

	return output.Cluster, err
}

// endregion
