package ocean_aks_np_auto_scaler

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/spotinst/spotinst-sdk-go/service/ocean/providers/azure_np"
	"github.com/spotinst/spotinst-sdk-go/spotinst"
	"github.com/spotinst/terraform-provider-spotinst/spotinst/commons"
)

func Setup(fieldsMap map[commons.FieldName]*commons.GenericField) {
	fieldsMap[AutoScaler] = commons.NewGenericField(
		commons.OceanAKSNPAutoScaler,
		AutoScaler,
		&schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					string(AutoscaleIsEnabled): {
						Type:     schema.TypeBool,
						Optional: true,
					},
					string(ResourceLimits): {
						Type:     schema.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								string(MaxVCPU): {
									Type:     schema.TypeInt,
									Optional: true,
									// here -1 is used to set MaxVCPU field to null when the customer doesn't want to set this param,
									//as terraform set it 0 for integer type param by default
									Default: -1,
								},
								string(MaxMemoryGib): {
									Type:     schema.TypeInt,
									Optional: true,
									// here -1 is used to set MaxMemoryGib field to null when the customer doesn't want to set this param,
									//as terraform set it 0 for integer type param by default
									Default: -1,
								},
							},
						},
					},
					string(Down): {
						Type:     schema.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								string(MaxScaleDownPercentage): {
									Type:     schema.TypeInt,
									Optional: true,
									// here -1 is used to set MaxScaleDownPercentage field to null when the customer
									//doesn't want to set this param, as terraform set it 0 for integer type param by default
									Default: -1,
								},
							},
						},
					},
					string(Headroom): {
						Type:     schema.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								string(Automatic): {
									Type:     schema.TypeList,
									Optional: true,
									MaxItems: 1,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											string(Percentage): {
												Type:     schema.TypeInt,
												Optional: true,
												Default:  -1,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},

		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			clusterWrapper := resourceObject.(*commons.AKSNPClusterWrapper)
			cluster := clusterWrapper.GetNPCluster()
			var result []interface{} = nil

			if cluster != nil && cluster.AutoScaler != nil {
				result = flattenAutoScaler(cluster.AutoScaler)
			}
			if len(result) > 0 {
				if err := resourceData.Set(string(AutoScaler), result); err != nil {
					return fmt.Errorf(string(commons.FailureFieldReadPattern), string(AutoScaler), err)
				}
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			clusterWrapper := resourceObject.(*commons.AKSNPClusterWrapper)
			cluster := clusterWrapper.GetNPCluster()
			var value *azure_np.AutoScaler = nil

			if v, ok := resourceData.GetOkExists(string(AutoScaler)); ok {
				if autoScaler, err := expandAutoScaler(v); err != nil {
					return err
				} else {
					value = autoScaler
				}
			}
			cluster.SetAutoScaler(value)
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			clusterWrapper := resourceObject.(*commons.AKSNPClusterWrapper)
			cluster := clusterWrapper.GetNPCluster()
			var value *azure_np.AutoScaler = nil
			if v, ok := resourceData.GetOkExists(string(AutoScaler)); ok {
				if autoScaler, err := expandAutoScaler(v); err != nil {
					return err
				} else {
					value = autoScaler
				}
			}
			cluster.SetAutoScaler(value)
			return nil
		},
		nil,
	)
}

func expandAutoScaler(data interface{}) (*azure_np.AutoScaler, error) {
	if list := data.([]interface{}); len(list) > 0 {
		autoScaler := &azure_np.AutoScaler{}

		if list[0] != nil {
			m := list[0].(map[string]interface{})

			if v, ok := m[string(AutoscaleIsEnabled)].(bool); ok {
				autoScaler.SetIsEnabled(spotinst.Bool(v))
			}

			if v, ok := m[string(ResourceLimits)]; ok && v != nil {

				resLimits, err := expandResourceLimits(v)
				if err != nil {
					return nil, err
				}
				if resLimits != nil {
					autoScaler.SetResourceLimits(resLimits)
				} else {
					log.Printf("resLimits == nil")
					autoScaler.ResourceLimits = nil
				}
			}

			if v, ok := m[string(Down)]; ok {
				down, err := expandDown(v)
				if err != nil {
					return nil, err
				}
				if down != nil {
					autoScaler.SetDown(down)
				} else {
					autoScaler.Down = nil
				}
			}

			if v, ok := m[string(Headroom)]; ok {
				headroom, err := expandHeadroom(v)
				if err != nil {
					return nil, err
				}
				if headroom != nil {
					autoScaler.SetHeadroom(headroom)
				} else {
					autoScaler.Headroom = nil
				}
			}
		}
		return autoScaler, nil
	}
	return nil, nil
}

func expandResourceLimits(data interface{}) (*azure_np.ResourceLimits, error) {
	resLimits := &azure_np.ResourceLimits{}
	list := data.([]interface{})

	if list == nil || len(list) == 0 {
		return nil, nil
	}

	m := list[0].(map[string]interface{})

	if v, ok := m[string(MaxMemoryGib)].(int); ok {
		// here -1 is used to set MaxMemoryGib field to null when the customer doesn't want to set this param,
		//as terraform set it 0 for integer type param by default.
		if v == -1 {
			resLimits.SetMaxMemoryGib(nil)
		} else {
			resLimits.SetMaxMemoryGib(spotinst.Int(v))
		}
	}

	if v, ok := m[string(MaxVCPU)].(int); ok {
		//Here -1 is used to set MaxVCPU field to null when the customer doesn't want to set this param,
		//as terraform set it 0 for integer type param by default.
		if v == -1 {
			resLimits.SetMaxVcpu(nil)
		} else {
			resLimits.SetMaxVcpu(spotinst.Int(v))
		}
	}
	return resLimits, nil
}

func expandDown(data interface{}) (*azure_np.Down, error) {
	down := &azure_np.Down{}
	list := data.([]interface{})

	if list == nil || len(list) == 0 {
		return nil, nil
	}

	m := list[0].(map[string]interface{})

	if v, ok := m[string(MaxScaleDownPercentage)].(int); ok {
		// here -1 is used to set MaxScaleDownPercentage field to null when the customer doesn't want to set this param,
		//as terraform set it 0 for integer type param by default.
		if v == -1 {
			down.SetMaxScaleDownPercentage(nil)
		} else {
			down.SetMaxScaleDownPercentage(spotinst.Int(v))
		}
	}
	return down, nil
}

func expandHeadroom(data interface{}) (*azure_np.Headroom, error) {
	list := data.([]interface{})
	headroom := &azure_np.Headroom{}
	if list == nil || len(list) == 0 {
		return nil, nil
	}
	if list[0] != nil {
		m := list[0].(map[string]interface{})

		if v, ok := m[string(Automatic)]; ok {
			automatic, err := expandAutomatic(v)
			if err != nil {
				return nil, err
			}
			if automatic != nil {
				headroom.SetAutomatic(automatic)
			} else {
				headroom.Automatic = nil
			}
		}
	}
	return headroom, nil
}

func expandAutomatic(data interface{}) (*azure_np.Automatic, error) {
	automatic := &azure_np.Automatic{}
	list := data.([]interface{})

	if list == nil || len(list) == 0 {
		return nil, nil
	}
	m := list[0].(map[string]interface{})

	if v, ok := m[string(Percentage)].(int); ok {
		// here -1 is used to set Percentage field to null when the customer doesn't want to set this param,
		//as terraform set it 0 for integer type param by default.
		if v == -1 {
			automatic.SetPercentage(nil)
		} else {
			automatic.SetPercentage((spotinst.Int(v)))
		}
	}

	return automatic, nil
}

func flattenAutoScaler(autoScaler *azure_np.AutoScaler) []interface{} {
	result := make(map[string]interface{})
	result[string(AutoscaleIsEnabled)] = spotinst.BoolValue(autoScaler.IsEnabled)

	if autoScaler.Headroom != nil {
		result[string(Headroom)] = flattenHeadroom(autoScaler.Headroom)
	}

	if autoScaler.Down != nil {
		result[string(Down)] = flattenDown(autoScaler.Down)
	}

	if autoScaler.ResourceLimits != nil {
		result[string(ResourceLimits)] = flattenResourceLimits(autoScaler.ResourceLimits)
	}
	return []interface{}{result}
}

func flattenHeadroom(headroom *azure_np.Headroom) []interface{} {
	result := make(map[string]interface{})

	if headroom.Automatic != nil {
		result[string(Automatic)] = flattenAutomatic(headroom.Automatic)
	}
	return []interface{}{result}
}

func flattenDown(autoScaleDown *azure_np.Down) []interface{} {
	down := make(map[string]interface{})

	if autoScaleDown != nil {
		value := spotinst.Int(-1)
		down[string(MaxScaleDownPercentage)] = value

		if autoScaleDown.MaxScaleDownPercentage != nil {
			down[string(MaxScaleDownPercentage)] = spotinst.IntValue(autoScaleDown.MaxScaleDownPercentage)
		}
	}
	return []interface{}{down}
}

func flattenResourceLimits(autoScaleResourceLimits *azure_np.ResourceLimits) []interface{} {
	resourceLimits := make(map[string]interface{})
	if autoScaleResourceLimits != nil {
		value := spotinst.Int(-1)
		resourceLimits[string(MaxVCPU)] = value
		resourceLimits[string(MaxMemoryGib)] = value

		if autoScaleResourceLimits.MaxVCPU != nil {
			resourceLimits[string(MaxVCPU)] = spotinst.IntValue(autoScaleResourceLimits.MaxVCPU)
		}
		if autoScaleResourceLimits.MaxMemoryGib != nil {
			resourceLimits[string(MaxMemoryGib)] = spotinst.IntValue(autoScaleResourceLimits.MaxMemoryGib)
		}
	}

	return []interface{}{resourceLimits}
}

func flattenAutomatic(autoScaleAutomatic *azure_np.Automatic) []interface{} {
	automatic := make(map[string]interface{})
	if autoScaleAutomatic != nil {
		value := spotinst.Int(-1)
		automatic[string(Percentage)] = value
		if autoScaleAutomatic.Percentage != nil {
			automatic[string(Percentage)] = spotinst.IntValue(autoScaleAutomatic.Percentage)
		}
	}
	return []interface{}{automatic}
}
