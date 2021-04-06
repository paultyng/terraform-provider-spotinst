package spotinst

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	v3 "github.com/spotinst/spotinst-sdk-go/service/elastigroup/providers/azure/v3"
	"github.com/spotinst/spotinst-sdk-go/spotinst"
	"github.com/spotinst/spotinst-sdk-go/spotinst/client"
	"github.com/spotinst/terraform-provider-spotinst/spotinst/azure_v3/elastigroup_azure"
	"github.com/spotinst/terraform-provider-spotinst/spotinst/azure_v3/elastigroup_azure_image"
	"github.com/spotinst/terraform-provider-spotinst/spotinst/azure_v3/elastigroup_azure_login"
	"github.com/spotinst/terraform-provider-spotinst/spotinst/azure_v3/elastigroup_azure_network"
	"github.com/spotinst/terraform-provider-spotinst/spotinst/azure_v3/elastigroup_azure_strategy"
	"github.com/spotinst/terraform-provider-spotinst/spotinst/azure_v3/elastigroup_azure_vm_sizes"
	"github.com/spotinst/terraform-provider-spotinst/spotinst/commons"
)

func resourceSpotinstElastigroupAzureV3() *schema.Resource {
	setupElastigroupAzureV3Resource()

	return &schema.Resource{
		Create: resourceSpotinstElastigroupAzureV3Create,
		Read:   resourceSpotinstElastigroupAzureV3Read,
		Update: resourceSpotinstElastigroupAzureV3Update,
		Delete: resourceSpotinstElastigroupAzureV3Delete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: commons.ElastigroupAzureV3Resource.GetSchemaMap(),
	}
}

func setupElastigroupAzureV3Resource() {
	fieldsMap := make(map[commons.FieldName]*commons.GenericField)

	elastigroup_azure.Setup(fieldsMap)
	elastigroup_azure_image.Setup(fieldsMap)
	elastigroup_azure_login.Setup(fieldsMap)
	elastigroup_azure_network.Setup(fieldsMap)
	elastigroup_azure_strategy.Setup(fieldsMap)
	elastigroup_azure_vm_sizes.Setup(fieldsMap)

	commons.ElastigroupAzureV3Resource = commons.NewElastigroupAzureV3Resource(fieldsMap)
}

//-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
//            Create
//-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
func resourceSpotinstElastigroupAzureV3Create(resourceData *schema.ResourceData, meta interface{}) error {
	log.Printf(string(commons.ResourceOnCreate),
		commons.ElastigroupAzureV3Resource.GetName())

	elastigroup, err := commons.ElastigroupAzureV3Resource.OnCreate(resourceData, meta)
	if err != nil {
		return err
	}

	groupId, err := createAzureV3Group(elastigroup, meta.(*Client))
	if err != nil {
		return err
	}

	resourceData.SetId(spotinst.StringValue(groupId))

	log.Printf("===> Elastigroup created successfully: %s <===", resourceData.Id())

	return resourceSpotinstElastigroupAzureV3Read(resourceData, meta)
}

func createAzureV3Group(group *v3.Group, spotinstClient *Client) (*string, error) {
	if json, err := commons.ToJson(group); err != nil {
		return nil, err
	} else {
		log.Printf("===> Group create configuration: %s", json)
	}

	input := &v3.CreateGroupInput{Group: group}

	var resp *v3.CreateGroupOutput = nil
	err := resource.Retry(time.Minute, func() *resource.RetryError {
		r, err := spotinstClient.elastigroup.CloudProviderAzureV3().Create(context.Background(), input)
		if err != nil {
			log.Printf("error: %v", err)
			// Some other error, report it.
			return resource.NonRetryableError(err)
		}
		resp = r
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("[ERROR] failed to create group: %s", err)
	}
	return resp.Group.ID, nil
}

//-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
//            Read
//-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
func resourceSpotinstElastigroupAzureV3Read(resourceData *schema.ResourceData, meta interface{}) error {
	id := resourceData.Id()
	log.Printf(string(commons.ResourceFieldOnRead),
		commons.ElastigroupAzureV3Resource.GetName(), id)

	input := &v3.ReadGroupInput{GroupID: spotinst.String(id)}
	resp, err := meta.(*Client).elastigroup.CloudProviderAzureV3().Read(context.Background(), input)
	if err != nil {
		// If the group was not found, return nil so that we can show
		// that the group does not exist
		if errs, ok := err.(client.Errors); ok && len(errs) > 0 {
			for _, err := range errs {
				if err.Code == ErrCodeGroupNotFound {
					resourceData.SetId("")
					return nil
				}
			}
		}

		// Some other error, report it.
		return fmt.Errorf("failed to read group: %s", err)
	}

	// If nothing was found, then return no state.
	groupResponse := resp.Group
	if groupResponse == nil {
		resourceData.SetId("")
		return nil
	}

	if err := commons.ElastigroupAzureV3Resource.OnRead(groupResponse, resourceData, meta); err != nil {
		return err
	}
	log.Printf("===> Elastigroup read successfully: %s <===", id)
	return nil
}

//-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
//            Update
//-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
func resourceSpotinstElastigroupAzureV3Update(resourceData *schema.ResourceData, meta interface{}) error {
	id := resourceData.Id()
	log.Printf(string(commons.ResourceOnUpdate),
		commons.ElastigroupAzureV3Resource.GetName(), id)

	shouldUpdate, elastigroup, err := commons.ElastigroupAzureV3Resource.OnUpdate(resourceData, meta)
	if err != nil {
		return err
	}

	if shouldUpdate {
		elastigroup.SetId(spotinst.String(id))
		if err := updateAzureV3Group(elastigroup, resourceData, meta); err != nil {
			return err
		}
	}

	log.Printf("===> Elastigroup updated successfully: %s <===", id)
	return resourceSpotinstElastigroupAzureV3Read(resourceData, meta)
}

func updateAzureV3Group(elastigroup *v3.Group, resourceData *schema.ResourceData, meta interface{}) error {
	var input = &v3.UpdateGroupInput{
		Group: elastigroup,
	}

	groupId := resourceData.Id()

	if json, err := commons.ToJson(elastigroup); err != nil {
		return err
	} else {
		log.Printf("===> Group update configuration: %s", json)
	}

	if _, err := meta.(*Client).elastigroup.CloudProviderAzureV3().Update(context.Background(), input); err != nil {
		return fmt.Errorf("[ERROR] Failed to update group [%v]: %v", groupId, err)
	}
	return nil
}

//-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
//            Delete
//-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
func resourceSpotinstElastigroupAzureV3Delete(resourceData *schema.ResourceData, meta interface{}) error {
	id := resourceData.Id()
	log.Printf(string(commons.ResourceOnDelete),
		commons.ElastigroupAzureV3Resource.GetName(), id)

	if err := deleteAzureV3Group(resourceData, meta); err != nil {
		return err
	}

	log.Printf("===> Elastigroup deleted successfully: %s <===", resourceData.Id())
	resourceData.SetId("")
	return nil
}

func deleteAzureV3Group(resourceData *schema.ResourceData, meta interface{}) error {
	groupId := resourceData.Id()
	input := &v3.DeleteGroupInput{
		GroupID: spotinst.String(groupId),
	}

	if json, err := commons.ToJson(input); err != nil {
		return err
	} else {
		log.Printf("===> Group delete configuration: %s", json)
	}

	if _, err := meta.(*Client).elastigroup.CloudProviderAzureV3().Delete(context.Background(), input); err != nil {
		return fmt.Errorf("[ERROR] onDelete() -> Failed to delete group: %s", err)
	}
	return nil
}
