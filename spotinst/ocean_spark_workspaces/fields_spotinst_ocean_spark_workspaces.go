package ocean_spark_workspaces

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/spotinst/spotinst-sdk-go/service/ocean/spark"
	"github.com/spotinst/spotinst-sdk-go/spotinst"

	"github.com/spotinst/terraform-provider-spotinst/spotinst/commons"
)

func Setup(fieldsMap map[commons.FieldName]*commons.GenericField) {
	fieldsMap[StorageClassOverride] = commons.NewGenericField(
		commons.OceanSparkWorkspaces,
		Workspaces,
		&schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					string(StorageClassOverride): {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
				},
			},
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			clusterWrapper := resourceObject.(*commons.SparkClusterWrapper)
			cluster := clusterWrapper.GetCluster()
			var result []interface{} = nil
			if cluster.Config != nil && cluster.Config.Workspaces != nil {
				result = flattenWorkspaces(cluster.Config.Workspaces)
			}
			if len(result) > 0 {
				if err := resourceData.Set(string(Workspaces), result); err != nil {
					return fmt.Errorf(commons.FailureFieldReadPattern, string(Workspaces), err)
				}
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			clusterWrapper := resourceObject.(*commons.SparkClusterWrapper)
			cluster := clusterWrapper.GetCluster()
			if value, ok := resourceData.GetOk(string(Workspaces)); ok {
				if workspaces, err := expandWorkspaces(value); err != nil {
					return err
				} else {
					if cluster.Config == nil {
						cluster.Config = &spark.Config{}
					}
					cluster.Config.SetWorkspaces(workspaces)
				}
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			clusterWrapper := resourceObject.(*commons.SparkClusterWrapper)
			cluster := clusterWrapper.GetCluster()
			var value *spark.WorkspacesConfig = nil
			if v, ok := resourceData.GetOk(string(Workspaces)); ok {
				if workspaces, err := expandWorkspaces(v); err != nil {
					return err
				} else {
					value = workspaces
				}
			}
			if cluster.Config == nil {
				cluster.Config = &spark.Config{}
			}
			cluster.Config.SetWorkspaces(value)
			return nil
		},
		nil,
	)
}

func flattenWorkspaces(workspaces *spark.WorkspacesConfig) []interface{} {
	if workspaces == nil {
		return nil
	}
	result := make(map[string]interface{})
	result[string(StorageClassOverride)] = spotinst.StringValue(workspaces.StorageClassOverride)
	return []interface{}{result}
}

func expandWorkspaces(data interface{}) (*spark.WorkspacesConfig, error) {
	workspaces := &spark.WorkspacesConfig{}
	list := data.([]interface{})
	if list == nil || list[0] == nil {
		return workspaces, nil
	}
	m := list[0].(map[string]interface{})

	if v, ok := m[string(StorageClassOverride)].(string); ok {
		workspaces.SetStorageClassOverride(spotinst.String(v))
	}

	return workspaces, nil
}
