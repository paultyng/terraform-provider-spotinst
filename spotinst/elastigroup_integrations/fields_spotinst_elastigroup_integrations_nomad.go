package elastigroup_integrations

import (
	"bytes"
	"fmt"

	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/spotinst/spotinst-sdk-go/service/elastigroup/providers/aws"
	"github.com/spotinst/spotinst-sdk-go/spotinst"
	"github.com/terraform-providers/terraform-provider-spotinst/spotinst/commons"
)

//-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
//            Setup
//-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
func SetupNomad(fieldsMap map[commons.FieldName]*commons.GenericField) {

	fieldsMap[IntegrationNomad] = commons.NewGenericField(
		commons.ElastigroupIntegrations,
		IntegrationNomad,
		&schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					string(MasterHost): &schema.Schema{
						Type:     schema.TypeString,
						Required: true,
					},

					string(MasterPort): &schema.Schema{
						Type:     schema.TypeInt,
						Required: true,
					},

					string(AutoscaleIsEnabled): &schema.Schema{
						Type:     schema.TypeBool,
						Optional: true,
					},

					string(AutoscaleCooldown): &schema.Schema{
						Type:     schema.TypeInt,
						Optional: true,
					},

					string(AclToken): &schema.Schema{
						Type:     schema.TypeString,
						Optional: true,
					},

					string(AutoscaleHeadroom): &schema.Schema{
						Type:     schema.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								string(CpuPerUnit): &schema.Schema{
									Type:     schema.TypeInt,
									Optional: true,
								},

								string(MemoryPerUnit): &schema.Schema{
									Type:     schema.TypeInt,
									Optional: true,
								},

								string(NumOfUnits): &schema.Schema{
									Type:     schema.TypeInt,
									Optional: true,
								},
							},
						},
					},

					string(AutoscaleDown): &schema.Schema{
						Type:     schema.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								string(EvaluationPeriods): &schema.Schema{
									Type:     schema.TypeInt,
									Optional: true,
								},
							},
						},
					},

					string(AutoscaleConstraints): &schema.Schema{
						Type:     schema.TypeSet,
						Optional: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								string(Key): &schema.Schema{
									Type:      schema.TypeString,
									Required:  true,
									StateFunc: attrStateFunc,
								},

								string(Value): &schema.Schema{
									Type:     schema.TypeString,
									Required: true,
								},
							},
						},
						Set: hashKV,
					},
				},
			},
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			elastigroup := resourceObject.(*aws.Group)
			if v, ok := resourceData.GetOk(string(IntegrationNomad)); ok {
				if integration, err := expandAWSGroupNomadIntegration(v, false); err != nil {
					return err
				} else {
					elastigroup.Integration.SetNomad(integration)
				}
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			elastigroup := resourceObject.(*aws.Group)
			var value *aws.NomadIntegration = nil
			if v, ok := resourceData.GetOk(string(IntegrationNomad)); ok {
				if integration, err := expandAWSGroupNomadIntegration(v, true); err != nil {
					return err
				} else {
					value = integration
				}
			}
			elastigroup.Integration.SetNomad(value)
			return nil
		},
		nil,
	)
}

//-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
//            Utils
//-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
func expandAWSGroupNomadIntegration(data interface{}, nullify bool) (*aws.NomadIntegration, error) {
	integration := &aws.NomadIntegration{}
	list := data.([]interface{})
	if list == nil || list[0] == nil {
		return integration, nil
	}
	m := list[0].(map[string]interface{})

	if v, ok := m[string(MasterHost)].(string); ok && v != "" {
		integration.SetMasterHost(spotinst.String(v))
	}

	if v, ok := m[string(MasterPort)].(int); ok && v > 0 {
		integration.SetMasterPort(spotinst.Int(v))
	}

	if v, ok := m[string(AclToken)].(string); ok && v != "" {
		integration.SetAclToken(spotinst.String(v))
	} else if nullify {
		integration.SetAclToken(nil)
	}

	if v, ok := m[string(AutoscaleIsEnabled)].(bool); ok {
		if integration.AutoScale == nil {
			integration.SetAutoScale(&aws.AutoScale{})
		}
		integration.AutoScale.SetIsEnabled(spotinst.Bool(v))
	}

	if v, ok := m[string(AutoscaleCooldown)].(int); ok && v > 0 {
		if integration.AutoScale == nil {
			integration.SetAutoScale(&aws.AutoScale{})
		}
		integration.AutoScale.SetCooldown(spotinst.Int(v))
	}

	if v, ok := m[string(AutoscaleHeadroom)]; ok {
		headroom, err := expandAWSGroupAutoScaleHeadroom(v)
		if err != nil {
			return nil, err
		}
		if headroom != nil {
			if integration.AutoScale == nil {
				integration.SetAutoScale(&aws.AutoScale{})
			}
			integration.AutoScale.SetHeadroom(headroom)
		}
	}

	if v, ok := m[string(AutoscaleDown)]; ok {
		down, err := expandAWSGroupAutoScaleDown(v)
		if err != nil {
			return nil, err
		}
		if down != nil {
			if integration.AutoScale == nil {
				integration.SetAutoScale(&aws.AutoScale{})
			}
			integration.AutoScale.SetDown(down)
		}
	}

	if v, ok := m[string(AutoscaleConstraints)]; ok {
		consts, err := expandAWSGroupAutoScaleConstraints(v)
		if err != nil {
			return nil, err
		}
		if consts != nil {
			if integration.AutoScale == nil {
				integration.SetAutoScale(&aws.AutoScale{})
			}
			integration.AutoScale.SetConstraints(consts)
		}
	}
	return integration, nil
}

func attrStateFunc(v interface{}) string {
	switch s := v.(type) {
	case string:
		return fmt.Sprintf("${%s}", s)
	default:
		return ""
	}
}

func hashKV(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	buf.WriteString(fmt.Sprintf("%s-", m[string(Key)].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m[string(Value)].(string)))
	return hashcode.String(buf.String())
}
