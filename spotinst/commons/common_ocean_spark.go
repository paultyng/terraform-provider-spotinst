package commons

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/spotinst/spotinst-sdk-go/service/ocean/spark"
)

const (
	OceanSparkResourceName ResourceName = "spotinst_ocean_spark"
)

var OceanSparkResource *OceanSparkTerraformResource

type OceanSparkTerraformResource struct {
	GenericResource
}

type SparkClusterWrapper struct {
	cluster *spark.Cluster
}

func NewOceanSparkResource(fieldsMap map[FieldName]*GenericField) *OceanSparkTerraformResource {
	return &OceanSparkTerraformResource{
		GenericResource: GenericResource{
			resourceName: OceanSparkResourceName,
			fields:       NewGenericFields(fieldsMap),
		},
	}
}

func (res *OceanSparkTerraformResource) OnCreate(
	resourceData *schema.ResourceData,
	meta interface{}) (*spark.Cluster, error) {

	if res.fields == nil || res.fields.fieldsMap == nil || len(res.fields.fieldsMap) == 0 {
		return nil, fmt.Errorf("resource fields are nil or empty, cannot create")
	}

	clusterWrapper := NewSparkClusterWrapper()

	for _, field := range res.fields.fieldsMap {
		if field.onCreate == nil {
			continue
		}
		log.Printf(string(ResourceFieldOnCreate), field.resourceAffinity, field.fieldNameStr)
		if err := field.onCreate(clusterWrapper, resourceData, meta); err != nil {
			return nil, err
		}
	}
	return clusterWrapper.GetCluster(), nil
}

func (res *OceanSparkTerraformResource) OnRead(
	cluster *spark.Cluster,
	resourceData *schema.ResourceData,
	meta interface{}) error {

	if res.fields == nil || res.fields.fieldsMap == nil || len(res.fields.fieldsMap) == 0 {
		return fmt.Errorf("resource fields are nil or empty, cannot read")
	}

	clusterWrapper := NewSparkClusterWrapper()
	clusterWrapper.SetCluster(cluster)

	for _, field := range res.fields.fieldsMap {
		if field.onRead == nil {
			continue
		}
		log.Printf(string(ResourceFieldOnRead), field.resourceAffinity, field.fieldNameStr)
		if err := field.onRead(clusterWrapper, resourceData, meta); err != nil {
			return err
		}
	}

	return nil
}

func (res *OceanSparkTerraformResource) OnUpdate(
	resourceData *schema.ResourceData,
	meta interface{}) (bool, *spark.Cluster, error) {

	if res.fields == nil || res.fields.fieldsMap == nil || len(res.fields.fieldsMap) == 0 {
		return false, nil, fmt.Errorf("resource fields are nil or empty, cannot update")
	}

	clusterWrapper := NewSparkClusterWrapper()
	hasChanged := false

	for _, field := range res.fields.fieldsMap {
		if field.onUpdate == nil {
			continue
		}
		if field.hasFieldChange(resourceData, meta) {
			log.Printf(string(ResourceFieldOnUpdate), field.resourceAffinity, field.fieldNameStr)
			if err := field.onUpdate(clusterWrapper, resourceData, meta); err != nil {
				return false, nil, err
			}
			hasChanged = true
		}
	}

	return hasChanged, clusterWrapper.GetCluster(), nil
}

func NewSparkClusterWrapper() *SparkClusterWrapper {
	return &SparkClusterWrapper{
		cluster: &spark.Cluster{
			Config: &spark.Config{
				Ingress:       &spark.IngressConfig{},
				Webhook:       &spark.WebhookConfig{},
				Compute:       &spark.ComputeConfig{},
				Workspaces:    &spark.WorkspacesConfig{},
				LogCollection: &spark.LogCollectionConfig{},
			},
		},
	}
}

func (clusterWrapper *SparkClusterWrapper) GetCluster() *spark.Cluster {
	return clusterWrapper.cluster
}

func (clusterWrapper *SparkClusterWrapper) SetCluster(cluster *spark.Cluster) {
	clusterWrapper.cluster = cluster
}
