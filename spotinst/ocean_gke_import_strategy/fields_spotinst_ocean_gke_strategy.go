package ocean_gke_import_strategy

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/spotinst/spotinst-sdk-go/service/ocean/providers/gcp"
	"github.com/spotinst/spotinst-sdk-go/spotinst"
	"github.com/spotinst/terraform-provider-spotinst/spotinst/commons"
)

func Setup(fieldsMap map[commons.FieldName]*commons.GenericField) {

	fieldsMap[Strategy] = commons.NewGenericField(
		commons.OceanGKEImportStrategy,
		Strategy,
		&schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					string(DrainingTimeout): {
						Type:     schema.TypeInt,
						Optional: true,
					},
					string(ProvisioningModel): {
						Type:     schema.TypeString,
						Optional: true,
					},
				},
			},
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			clusterWrapper := resourceObject.(*commons.GKEImportClusterWrapper)
			cluster := clusterWrapper.GetCluster()
			var result []interface{} = nil
			if cluster.Strategy != nil {
				strategy := cluster.Strategy
				result = flattenStrategy(strategy)
			}
			if result != nil {
				if err := resourceData.Set(string(Strategy), result); err != nil {
					return fmt.Errorf(string(commons.FailureFieldReadPattern), string(Strategy), err)
				}
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			clusterWrapper := resourceObject.(*commons.GKEImportClusterWrapper)
			cluster := clusterWrapper.GetCluster()

			if v, ok := resourceData.GetOk(string(Strategy)); ok {
				if strategy, err := expandStrategy(v); err != nil {
					return err
				} else {
					cluster.SetStrategy(strategy)
				}
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			clusterWrapper := resourceObject.(*commons.GKEImportClusterWrapper)
			cluster := clusterWrapper.GetCluster()
			var value *gcp.Strategy = nil

			if v, ok := resourceData.GetOk(string(Strategy)); ok {
				if strategy, err := expandStrategy(v); err != nil {
					return err
				} else {
					value = strategy
				}
			}
			cluster.SetStrategy(value)
			return nil
		},
		nil,
	)
}

func flattenStrategy(strategy *gcp.Strategy) []interface{} {
	var out []interface{}

	if strategy != nil {
		result := make(map[string]interface{})

		if strategy.DrainingTimeout != nil {
			result[string(DrainingTimeout)] = spotinst.IntValue(strategy.DrainingTimeout)
		}

		if strategy.ProvisioningModel != nil {
			result[string(ProvisioningModel)] = spotinst.StringValue(strategy.ProvisioningModel)
		}

		if len(result) > 0 {
			out = append(out, result)
		}
	}
	return out
}

func expandStrategy(data interface{}) (*gcp.Strategy, error) {
	if list := data.([]interface{}); len(list) > 0 {
		strategy := &gcp.Strategy{}

		if list[0] != nil {
			m := list[0].(map[string]interface{})

			if v, ok := m[string(DrainingTimeout)].(int); ok && v > 0 {
				strategy.SetDrainingTimeout(spotinst.Int(v))
			} else {
				strategy.SetDrainingTimeout(nil)
			}

			if v, ok := m[string(ProvisioningModel)].(string); ok && v != "" {
				strategy.SetProvisioningModel(spotinst.String(v))
			} else {
				strategy.SetProvisioningModel(nil)
			}
		}

		return strategy, nil
	}

	return nil, nil
}
