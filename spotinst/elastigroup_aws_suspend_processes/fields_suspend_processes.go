package elastigroup_aws_suspend_processes

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/spotinst/spotinst-sdk-go/service/elastigroup/providers/aws"
	"github.com/spotinst/spotinst-sdk-go/spotinst"
	"github.com/spotinst/terraform-provider-spotinst/spotinst/commons"
)

func Setup(fieldsMap map[commons.FieldName]*commons.GenericField) {

	fieldsMap[GroupID] = commons.NewGenericField(
		commons.SuspendProcesses,
		GroupID,
		&schema.Schema{
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			spWrapper := resourceObject.(*commons.SuspendProcessesWrapper)

			value := spotinst.String(resourceData.Get(string(GroupID)).(string))
			spWrapper.GroupID = value

			if err := resourceData.Set(string(GroupID), value); err != nil {
				return fmt.Errorf(commons.FailureFieldReadPattern, string(GroupID), err)
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			spWrapper := resourceObject.(*commons.SuspendProcessesWrapper)
			spWrapper.GroupID = spotinst.String(resourceData.Get(string(GroupID)).(string))
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			spWrapper := resourceObject.(*commons.SuspendProcessesWrapper)
			spWrapper.GroupID = spotinst.String(resourceData.Get(string(GroupID)).(string))
			return nil
		},
		nil,
	)

	fieldsMap[Suspension] = commons.NewGenericField(
		commons.SuspendProcesses,
		Suspension,
		&schema.Schema{
			Type:     schema.TypeList,
			Required: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{

					string(Name): {
						Type:     schema.TypeString,
						Required: true,
					},
				},
			},
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			spWrapper := resourceObject.(*commons.SuspendProcessesWrapper)
			var result []interface{} = nil
			if spWrapper.SuspendProcesses != nil {
				result = flattenSuspensions(spWrapper.SuspendProcesses.Suspensions)
			}
			if err := resourceData.Set(string(Suspension), result); err != nil {
				return fmt.Errorf(string(commons.FailureFieldReadPattern), string(Suspension), err)
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			spWrapper := resourceObject.(*commons.SuspendProcessesWrapper)
			if v, ok := resourceData.GetOk(string(Suspension)); ok {
				if v, err := expandSuspensions(v); err != nil {
					return err
				} else {
					spWrapper.SuspendProcesses.Suspensions = v
				}
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			spWrapper := resourceObject.(*commons.SuspendProcessesWrapper)
			var value []*aws.Suspension = nil
			if v, ok := resourceData.GetOk(string(Suspension)); ok {
				if v, err := expandSuspensions(v); err != nil {
					return err
				} else {
					value = v
				}
			}
			spWrapper.SuspendProcesses.Suspensions = value
			return nil
		},
		nil,
	)
}

func expandSuspensions(data interface{}) ([]*aws.Suspension, error) {
	list := data.([]interface{})
	suspensions := make([]*aws.Suspension, 0, len(list))

	for _, v := range list {
		attr, ok := v.(map[string]interface{})
		if !ok {
			continue
		}

		suspension := &aws.Suspension{}

		if v, ok := attr[string(Name)].(string); ok && v != "" {
			suspension.SetName(spotinst.String(v))
		}

		suspensions = append(suspensions, suspension)
	}

	return suspensions, nil
}

func flattenSuspensions(suspensions []*aws.Suspension) []interface{} {
	result := make([]interface{}, 0, len(suspensions))

	for _, suspension := range suspensions {
		m := make(map[string]interface{})
		m[string(Name)] = spotinst.StringValue(suspension.Name)
		result = append(result, m)
	}
	return result

}
