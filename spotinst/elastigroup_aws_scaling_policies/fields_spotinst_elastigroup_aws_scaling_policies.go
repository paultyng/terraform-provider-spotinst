package elastigroup_aws_scaling_policies

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/spotinst/spotinst-sdk-go/service/elastigroup/providers/aws"
	"github.com/spotinst/spotinst-sdk-go/spotinst"
	"github.com/spotinst/terraform-provider-spotinst/spotinst/commons"
)

func Setup(fieldsMap map[commons.FieldName]*commons.GenericField) {

	fieldsMap[ScalingUpPolicy] = commons.NewGenericField(
		commons.ElastigroupAWSScalingPolicies,
		ScalingUpPolicy,
		upDownScalingPolicySchema(),
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			egWrapper := resourceObject.(*commons.ElastigroupWrapper)
			elastigroup := egWrapper.GetElastigroup()
			var policiesResult []interface{} = nil
			if elastigroup.Scaling != nil && elastigroup.Scaling.Up != nil {
				scaleUpPolicies := elastigroup.Scaling.Up
				policiesResult = flattenAWSGroupScalingPolicy(scaleUpPolicies)
			}
			if err := resourceData.Set(string(ScalingUpPolicy), policiesResult); err != nil {
				return fmt.Errorf(string(commons.FailureFieldReadPattern), string(ScalingUpPolicy), err)
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			egWrapper := resourceObject.(*commons.ElastigroupWrapper)
			elastigroup := egWrapper.GetElastigroup()
			if v, ok := resourceData.GetOk(string(ScalingUpPolicy)); ok {
				if policies, err := expandAWSGroupScalingPolicies(v); err != nil {
					return err
				} else {
					elastigroup.Scaling.SetUp(policies)
				}
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			egWrapper := resourceObject.(*commons.ElastigroupWrapper)
			elastigroup := egWrapper.GetElastigroup()
			var value []*aws.ScalingPolicy = nil
			if v, ok := resourceData.GetOk(string(ScalingUpPolicy)); ok && v != nil {
				if policies, err := expandAWSGroupScalingPolicies(v); err != nil {
					return err
				} else {
					value = policies
				}
			}
			elastigroup.Scaling.SetUp(value)
			return nil
		},
		nil,
	)

	fieldsMap[ScalingDownPolicy] = commons.NewGenericField(
		commons.ElastigroupAWSScalingPolicies,
		ScalingDownPolicy,
		upDownScalingPolicySchema(),
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			egWrapper := resourceObject.(*commons.ElastigroupWrapper)
			elastigroup := egWrapper.GetElastigroup()
			var policiesResult []interface{} = nil
			if elastigroup.Scaling != nil && elastigroup.Scaling.Down != nil {
				scaleDownPolicies := elastigroup.Scaling.Down
				policiesResult = flattenAWSGroupScalingPolicy(scaleDownPolicies)
			}
			if err := resourceData.Set(string(ScalingDownPolicy), policiesResult); err != nil {
				return fmt.Errorf(string(commons.FailureFieldReadPattern), string(ScalingDownPolicy), err)
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			egWrapper := resourceObject.(*commons.ElastigroupWrapper)
			elastigroup := egWrapper.GetElastigroup()
			if v, ok := resourceData.GetOk(string(ScalingDownPolicy)); ok {
				if policies, err := expandAWSGroupScalingPolicies(v); err != nil {
					return err
				} else {
					elastigroup.Scaling.SetDown(policies)
				}
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			egWrapper := resourceObject.(*commons.ElastigroupWrapper)
			elastigroup := egWrapper.GetElastigroup()
			var value []*aws.ScalingPolicy = nil
			if v, ok := resourceData.GetOk(string(ScalingDownPolicy)); ok && v != nil {
				if policies, err := expandAWSGroupScalingPolicies(v); err != nil {
					return err
				} else {
					value = policies
				}
			}
			elastigroup.Scaling.SetDown(value)
			return nil
		},
		nil,
	)

	fieldsMap[ScalingTargetPolicy] = commons.NewGenericField(
		commons.ElastigroupAWSScalingPolicies,
		ScalingTargetPolicy,
		targetScalingPolicySchema(),
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			egWrapper := resourceObject.(*commons.ElastigroupWrapper)
			elastigroup := egWrapper.GetElastigroup()
			var policiesResult []interface{} = nil
			if elastigroup.Scaling != nil && elastigroup.Scaling.Target != nil {
				scaleTargetPolicies := elastigroup.Scaling.Target
				policiesResult = flattenAWSGroupScalingPolicy(scaleTargetPolicies)
			}
			if err := resourceData.Set(string(ScalingTargetPolicy), policiesResult); err != nil {
				return fmt.Errorf(string(commons.FailureFieldReadPattern), string(ScalingTargetPolicy), err)
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			egWrapper := resourceObject.(*commons.ElastigroupWrapper)
			elastigroup := egWrapper.GetElastigroup()
			if v, ok := resourceData.GetOk(string(ScalingTargetPolicy)); ok {
				if policies, err := expandAWSGroupScalingPolicies(v); err != nil {
					return err
				} else {
					elastigroup.Scaling.SetTarget(policies)
				}
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			egWrapper := resourceObject.(*commons.ElastigroupWrapper)
			elastigroup := egWrapper.GetElastigroup()
			var value []*aws.ScalingPolicy = nil
			if v, ok := resourceData.GetOk(string(ScalingTargetPolicy)); ok && v != nil {
				if policies, err := expandAWSGroupScalingPolicies(v); err != nil {
					return err
				} else {
					value = policies
				}
			}
			elastigroup.Scaling.SetTarget(value)
			return nil
		},
		nil,
	)

	fieldsMap[MultipleMetrics] = commons.NewGenericField(
		commons.ElastigroupAWSScalingPolicies,
		MultipleMetrics,
		&schema.Schema{
			Type:     schema.TypeSet,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					string(Expressions): {
						Type:     schema.TypeSet,
						Optional: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								string(Expression): {
									Type:     schema.TypeString,
									Required: true,
								},

								string(Name): {
									Type:     schema.TypeString,
									Required: true,
								},
							},
						},
					},

					string(Metrics): {
						Type:     schema.TypeSet,
						Optional: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								string(Dimensions): {
									Type: schema.TypeList,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											string(DimensionName): {
												Type:     schema.TypeString,
												Required: true,
											},

											string(DimensionValue): {
												Type:     schema.TypeString,
												Optional: true,
											},
										},
									},
									Optional: true,
								},

								string(ExtendedStatistic): {
									Type:     schema.TypeString,
									Optional: true,
								},
								string(MetricName): {
									Type:     schema.TypeString,
									Required: true,
								},
								string(Name): {
									Type:     schema.TypeString,
									Required: true,
								},
								string(Namespace): {
									Type:     schema.TypeString,
									Required: true,
								},
								string(Statistic): {
									Type:     schema.TypeString,
									Optional: true,
								},
								string(Unit): {
									Type:     schema.TypeString,
									Optional: true,
								},
							},
						},
					},
				},
			},
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			egWrapper := resourceObject.(*commons.ElastigroupWrapper)
			elastigroup := egWrapper.GetElastigroup()
			var policiesResult []interface{} = nil
			if elastigroup.Scaling != nil && elastigroup.Scaling.MultipleMetrics != nil {
				multipleMetrics := elastigroup.Scaling.MultipleMetrics
				policiesResult = flattenMultipleMetrics(multipleMetrics)
			}
			if err := resourceData.Set(string(MultipleMetrics), policiesResult); err != nil {
				return fmt.Errorf(string(commons.FailureFieldReadPattern), string(MultipleMetrics), err)
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			egWrapper := resourceObject.(*commons.ElastigroupWrapper)
			elastigroup := egWrapper.GetElastigroup()
			if v, ok := resourceData.GetOk(string(MultipleMetrics)); ok {
				if multipleMetrics, err := expandMultipleMetrics(v); err != nil {
					return err
				} else {
					elastigroup.Scaling.SetMultipleMetrics(multipleMetrics)
				}
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			egWrapper := resourceObject.(*commons.ElastigroupWrapper)
			elastigroup := egWrapper.GetElastigroup()
			var value *aws.MultipleMetrics = nil
			if v, ok := resourceData.GetOk(string(MultipleMetrics)); ok && v != nil {
				if multipleMetrics, err := expandMultipleMetrics(v); err != nil {
					return err
				} else {
					value = multipleMetrics
				}
			}
			elastigroup.Scaling.SetMultipleMetrics(value)
			return nil
		},
		nil,
	)

}

func baseScalingPolicySchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				string(PolicyName): {
					Type:     schema.TypeString,
					Required: true,
				},

				string(MetricName): {
					Type:     schema.TypeString,
					Required: true,
				},

				string(Namespace): {
					Type:     schema.TypeString,
					Required: true,
				},

				string(Source): {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},

				string(Statistic): {
					Type:     schema.TypeString,
					Optional: true,
					Default:  "average",
				},

				string(Unit): {
					Type:     schema.TypeString,
					Optional: true,
					Default:  "percent",
				},

				string(Cooldown): {
					Type:     schema.TypeInt,
					Optional: true,
					Default:  300,
				},

				string(Dimensions): {
					Type: schema.TypeList,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							string(DimensionName): {
								Type:     schema.TypeString,
								Required: true,
							},

							string(DimensionValue): {
								Type:     schema.TypeString,
								Optional: true,
							},
						},
					},
					Optional: true,
				},
			},
		},
	}
}

func upDownScalingPolicySchema() *schema.Schema {
	o := baseScalingPolicySchema()
	s := o.Elem.(*schema.Resource).Schema

	s[string(Threshold)] = &schema.Schema{
		Type:     schema.TypeFloat,
		Optional: true,
	}

	s[string(Adjustment)] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	}

	s[string(MinTargetCapacity)] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	}

	s[string(MaxTargetCapacity)] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	}

	s[string(Operator)] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
		Computed: true,
	}

	s[string(EvaluationPeriods)] = &schema.Schema{
		Type:     schema.TypeInt,
		Optional: true,
		Computed: true,
	}

	s[string(Period)] = &schema.Schema{
		Type:     schema.TypeInt,
		Optional: true,
		Computed: true,
	}

	s[string(Minimum)] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	}

	s[string(Maximum)] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	}

	s[string(Target)] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	}

	s[string(ActionType)] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	}

	s[string(IsEnabled)] = &schema.Schema{
		Type:     schema.TypeBool,
		Optional: true,
		Default:  true,
	}

	s[string(StepAdjustments)] = &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				string(Action): {
					Type:     schema.TypeSet,
					Required: true,
					MaxItems: 1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							string(Adjustment): {
								Type:     schema.TypeString,
								Optional: true,
							},

							string(Maximum): {
								Type:     schema.TypeString,
								Optional: true,
							},

							string(Minimum): {
								Type:     schema.TypeString,
								Optional: true,
							},

							string(MaxTargetCapacity): {
								Type:     schema.TypeString,
								Optional: true,
							},

							string(MinTargetCapacity): {
								Type:     schema.TypeString,
								Optional: true,
							},

							string(Target): {
								Type:     schema.TypeString,
								Optional: true,
							},

							string(Type): {
								Type:     schema.TypeString,
								Required: true,
							},
						},
					},
				},

				string(Threshold): {
					Type:     schema.TypeInt,
					Required: true,
				},
			},
		},
	}

	return o
}

func targetScalingPolicySchema() *schema.Schema {
	o := baseScalingPolicySchema()
	s := o.Elem.(*schema.Resource).Schema

	s[string(Target)] = &schema.Schema{
		Type:     schema.TypeFloat,
		Required: true,
	}

	s[string(PredictiveMode)] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	}

	s[string(MaxCapacityPerScale)] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	}

	return o
}

func expandMultipleMetrics(data interface{}) (*aws.MultipleMetrics, error) {
	list := data.(*schema.Set).List()
	multipleMetrics := &aws.MultipleMetrics{}

	if list != nil && list[0] != nil {
		m := list[0].(map[string]interface{})

		if v, ok := m[string(Expressions)]; ok {
			expressions := expandExpressions(v.(interface{}))
			if len(expressions) > 0 {
				multipleMetrics.SetExpressions(expressions)
			}
		}

		if v, ok := m[string(Metrics)]; ok {
			if metrics, err := expandMetrics(v.(interface{})); err != nil {
				return nil, err
			} else {
				multipleMetrics.SetMetrics(metrics)
			}
		}

		if multipleMetrics.Expressions != nil || multipleMetrics.Metrics != nil {
			return multipleMetrics, nil
		}

	}

	return nil, nil
}

func expandAWSGroupScalingPolicies(data interface{}) ([]*aws.ScalingPolicy, error) {
	list := data.(*schema.Set).List()
	policies := make([]*aws.ScalingPolicy, 0, len(list))
	for _, item := range list {
		m := item.(map[string]interface{})
		policy := &aws.ScalingPolicy{}

		if v, ok := m[string(PolicyName)].(string); ok && v != "" {
			policy.SetPolicyName(spotinst.String(v))
		}

		if v, ok := m[string(MetricName)].(string); ok && v != "" {
			policy.SetMetricName(spotinst.String(v))
		}

		if v, ok := m[string(Namespace)].(string); ok && v != "" {
			policy.SetNamespace(spotinst.String(v))
		}

		if v, ok := m[string(Source)].(string); ok && v != "" {
			policy.SetSource(spotinst.String(v))
		}

		if v, ok := m[string(Statistic)].(string); ok && v != "" {
			policy.SetStatistic(spotinst.String(v))
		}

		if v, ok := m[string(Unit)].(string); ok && v != "" {
			policy.SetUnit(spotinst.String(v))
		}

		if v, ok := m[string(Threshold)].(float64); ok && v > 0 {
			policy.SetThreshold(spotinst.Float64(v))
		}

		if v, ok := m[string(Operator)].(string); ok && v != "" {
			policy.SetOperator(spotinst.String(v))
		}

		if v, ok := m[string(Period)].(int); ok && v > 0 {
			policy.SetPeriod(spotinst.Int(v))
		}

		if v, ok := m[string(EvaluationPeriods)].(int); ok && v > 0 {
			policy.SetEvaluationPeriods(spotinst.Int(v))
		}

		if v, ok := m[string(Cooldown)].(int); ok && v > 0 {
			policy.SetCooldown(spotinst.Int(v))
		}

		if v, ok := m[string(IsEnabled)].(bool); ok {
			policy.SetIsEnabled(spotinst.Bool(v))
		}

		if v, ok := m[string(Dimensions)]; ok {
			dimensions := expandAWSGroupScalingPolicyDimensions(v.(interface{}))
			if len(dimensions) > 0 {
				policy.SetDimensions(dimensions)
			}
		}

		if v, ok := m[string(ActionType)].(string); ok && v != "" {
			action := &aws.Action{}
			action.SetType(spotinst.String(v))

			if v, ok := m[string(Adjustment)].(string); ok && v != "" {
				action.SetAdjustment(spotinst.String(v))
			}

			if v, ok := m[string(MinTargetCapacity)].(string); ok && v != "" {
				action.SetMinTargetCapacity(spotinst.String(v))
			}

			if v, ok := m[string(MaxTargetCapacity)].(string); ok && v != "" {
				action.SetMaxTargetCapacity(spotinst.String(v))
			}

			if v, ok := m[string(Minimum)].(string); ok && v != "" {
				action.SetMinimum(spotinst.String(v))
			}

			if v, ok := m[string(Maximum)].(string); ok && v != "" {
				action.SetMaximum(spotinst.String(v))
			}

			if v, ok := m[string(Target)].(string); ok && v != "" {
				action.SetTarget(spotinst.String(v))
			}

			policy.SetAction(action)
		} else {
			if v, ok := m[string(Adjustment)].(int); ok && v > 0 {
				policy.SetAdjustment(spotinst.Int(v))
			}

			if v, ok := m[string(MinTargetCapacity)].(int); ok && v > 0 {
				policy.SetMinTargetCapacity(spotinst.Int(v))
			}

			if v, ok := m[string(MaxTargetCapacity)].(int); ok && v > 0 {
				policy.SetMaxTargetCapacity(spotinst.Int(v))
			}
		}

		if v, ok := m[string(StepAdjustments)]; ok {
			stepAdjustments := expandAWSGroupScalingPolicyStepAdjustments(v.(interface{}))
			if len(stepAdjustments) > 0 {
				policy.SetStepAdjustments(stepAdjustments)
			}
		}

		// Target scaling policy?
		if policy.Threshold == nil {
			if v, ok := m[string(Target)].(float64); ok && v >= 0 {
				policy.SetTarget(spotinst.Float64(v))
			}

			if v, ok := m[string(PredictiveMode)].(string); ok && v != "" {
				policy.SetPredictive(&aws.Predictive{Mode: spotinst.String(v)})
			}

			if v, ok := m[string(MaxCapacityPerScale)].(string); ok && v != "" {
				policy.SetMaxCapacityPerScale(spotinst.String(v))
			}
		}

		if policy.Namespace != nil {
			policies = append(policies, policy)
		}
	}

	return policies, nil
}

func expandAWSGroupScalingPolicyDimensions(data interface{}) []*aws.Dimension {
	list := data.([]interface{})
	dimensions := make([]*aws.Dimension, 0, len(list))
	for _, v := range list {
		attr, ok := v.(map[string]interface{})
		if !ok {
			continue
		}
		if _, ok := attr[string(DimensionName)]; !ok {
			continue
		}

		if _, ok := attr[string(DimensionValue)]; !ok {
			continue
		}
		dimension := &aws.Dimension{
			Name:  spotinst.String(attr[string(DimensionName)].(string)),
			Value: spotinst.String(attr[string(DimensionValue)].(string)),
		}
		if (dimension.Name != nil) && (dimension.Value != nil) {
			dimensions = append(dimensions, dimension)
		}
	}
	return dimensions
}

func expandAWSGroupScalingPolicyStepAdjustments(data interface{}) []*aws.StepAdjustment {
	list := data.(*schema.Set).List()
	stepAdjustments := make([]*aws.StepAdjustment, 0, len(list))

	for _, item := range list {
		m := item.(map[string]interface{})
		stepAdjustment := &aws.StepAdjustment{}

		if v, ok := m[string(Threshold)].(int); ok && v > 0 {
			stepAdjustment.SetThreshold(spotinst.Int(v))
		}

		if v, ok := m[string(Action)]; ok {
			action := expandAWSGroupScalingPolicyStepAdjustmentsActions(v.(interface{}))
			if action != nil {
				stepAdjustment.SetAction(action)
			}
		}

		if (stepAdjustment.Threshold != nil) && (stepAdjustment.Action != nil) {
			stepAdjustments = append(stepAdjustments, stepAdjustment)
		}
	}
	return stepAdjustments
}

func expandAWSGroupScalingPolicyStepAdjustmentsActions(data interface{}) *aws.Action {
	list := data.(*schema.Set).List()
	action := &aws.Action{}

	if list != nil && list[0] != nil {
		m := list[0].(map[string]interface{})

		if v, ok := m[string(Adjustment)].(string); ok && v != "" {
			action.SetAdjustment(spotinst.String(v))
		}

		if v, ok := m[string(Maximum)].(string); ok && v != "" {
			action.SetMaximum(spotinst.String(v))
		}

		if v, ok := m[string(MaxTargetCapacity)].(string); ok && v != "" {
			action.SetMaxTargetCapacity(spotinst.String(v))
		}

		if v, ok := m[string(Minimum)].(string); ok && v != "" {
			action.SetMinimum(spotinst.String(v))
		}

		if v, ok := m[string(MinTargetCapacity)].(string); ok && v != "" {
			action.SetMinTargetCapacity(spotinst.String(v))
		}

		if v, ok := m[string(Target)].(string); ok && v != "" {
			action.SetTarget(spotinst.String(v))
		}

		if v, ok := m[string(Type)].(string); ok && v != "" {
			action.SetType(spotinst.String(v))
		}

		if action.Type != nil {
			return action
		}

	}
	return nil
}

func expandMetrics(data interface{}) ([]*aws.Metrics, error) {
	list := data.(*schema.Set).List()
	metrics := make([]*aws.Metrics, 0, len(list))
	for _, item := range list {
		m := item.(map[string]interface{})
		metric := &aws.Metrics{}

		if v, ok := m[string(Name)].(string); ok && v != "" {
			metric.SetName(spotinst.String(v))
		}

		if v, ok := m[string(Namespace)].(string); ok && v != "" {
			metric.SetNamespace(spotinst.String(v))
		}

		if v, ok := m[string(MetricName)].(string); ok && v != "" {
			metric.SetMetricName(spotinst.String(v))
		}

		if v, ok := m[string(ExtendedStatistic)].(string); ok && v != "" {
			metric.SetExtendedStatistic(spotinst.String(v))
		}

		if v, ok := m[string(Statistic)].(string); ok && v != "" {
			metric.SetStatistic(spotinst.String(v))
		}

		if v, ok := m[string(Unit)].(string); ok && v != "" {
			metric.SetUnit(spotinst.String(v))
		}

		if v, ok := m[string(Dimensions)]; ok {
			dimensions := expandAWSGroupScalingPolicyDimensions(v.(interface{}))
			if len(dimensions) > 0 {
				metric.SetDimensions(dimensions)
			}
		}

		if (metric.Name != nil) && (metric.Namespace != nil) && (metric.MetricName != nil) {
			metrics = append(metrics, metric)
		}
	}
	return metrics, nil
}

func expandExpressions(data interface{}) []*aws.Expressions {
	list := data.(*schema.Set).List()
	expressions := make([]*aws.Expressions, 0, len(list))
	for _, item := range list {
		m := item.(map[string]interface{})
		expression := &aws.Expressions{}

		if v, ok := m[string(Name)].(string); ok && v != "" {
			expression.SetName(spotinst.String(v))
		}

		if v, ok := m[string(Expression)].(string); ok && v != "" {
			expression.SetExpression(spotinst.String(v))
		}

		if (expression.Name != nil) && (expression.Expression != nil) {
			expressions = append(expressions, expression)
		}
	}
	return expressions
}

func flattenAWSGroupScalingPolicy(policies []*aws.ScalingPolicy) []interface{} {
	result := make([]interface{}, 0, len(policies))
	for _, policy := range policies {
		m := make(map[string]interface{})
		m[string(PolicyName)] = spotinst.StringValue(policy.PolicyName)
		m[string(MetricName)] = spotinst.StringValue(policy.MetricName)
		m[string(Namespace)] = spotinst.StringValue(policy.Namespace)
		m[string(Source)] = spotinst.StringValue(policy.Source)
		m[string(Statistic)] = spotinst.StringValue(policy.Statistic)
		m[string(Unit)] = spotinst.StringValue(policy.Unit)
		m[string(Cooldown)] = spotinst.IntValue(policy.Cooldown)

		if policy.Dimensions != nil && len(policy.Dimensions) > 0 {
			dimMap := make([]interface{}, 0, len(policy.Dimensions))
			for _, dimension := range policy.Dimensions {
				d := make(map[string]interface{})
				d[string(DimensionName)] = spotinst.StringValue(dimension.Name)
				d[string(DimensionValue)] = spotinst.StringValue(dimension.Value)

				if (d[string(DimensionName)] != nil) && (d[string(DimensionValue)] != nil) {
					dimMap = append(dimMap, d)
				}
			}
			m[string(Dimensions)] = dimMap
		}

		if policy.Action != nil && policy.Action.Type != nil {
			m[string(ActionType)] = spotinst.StringValue(policy.Action.Type)
			m[string(Adjustment)] = spotinst.StringValue(policy.Action.Adjustment)
			m[string(MinTargetCapacity)] = spotinst.StringValue(policy.Action.MinTargetCapacity)
			m[string(MaxTargetCapacity)] = spotinst.StringValue(policy.Action.MaxTargetCapacity)
			m[string(Minimum)] = spotinst.StringValue(policy.Action.Minimum)
			m[string(Maximum)] = spotinst.StringValue(policy.Action.Maximum)
			m[string(Target)] = spotinst.StringValue(policy.Action.Target)
			m[string(EvaluationPeriods)] = spotinst.IntValue(policy.EvaluationPeriods)
			m[string(Period)] = spotinst.IntValue(policy.Period)
			m[string(Threshold)] = spotinst.Float64Value(policy.Threshold)
			m[string(Operator)] = spotinst.StringValue(policy.Operator)
		}

		// Target scaling policy?
		if policy.Threshold == nil && policy.StepAdjustments == nil {
			m[string(Target)] = spotinst.Float64Value(policy.Target)

			if policy.Predictive != nil && policy.Predictive.Mode != nil {
				m[string(PredictiveMode)] = spotinst.StringValue(policy.Predictive.Mode)
			}

			if policy.MaxCapacityPerScale != nil {
				m[string(MaxCapacityPerScale)] = spotinst.StringValue(policy.MaxCapacityPerScale)
			}
		} else {
			m[string(IsEnabled)] = spotinst.BoolValue(policy.IsEnabled)
		}

		if policy.StepAdjustments != nil && len(policy.StepAdjustments) > 0 {
			stepAdjMap := make([]interface{}, 0, len(policy.StepAdjustments))
			for _, step := range policy.StepAdjustments {
				s := make(map[string]interface{})
				s[string(Threshold)] = spotinst.IntValue(step.Threshold)

				s[string(Action)] = flattenAction(step.Action)

				if (s[string(Threshold)] != nil) && (s[string(Action)] != nil) {
					stepAdjMap = append(stepAdjMap, s)
				}
			}
			m[string(StepAdjustments)] = stepAdjMap
		}

		result = append(result, m)
	}
	return result
}

func flattenMultipleMetrics(multipleMetrics *aws.MultipleMetrics) []interface{} {
	result := make(map[string]interface{})

	result[string(Expressions)] = flattenExpressions(multipleMetrics.Expressions)
	result[string(Metrics)] = flattenMetrics(multipleMetrics.Metrics)

	return []interface{}{result}
}

func flattenExpressions(expressions []*aws.Expressions) []interface{} {
	result := make([]interface{}, 0, len(expressions))
	for _, expression := range expressions {
		m := make(map[string]interface{})

		m[string(Name)] = spotinst.StringValue(expression.Name)
		m[string(Expression)] = spotinst.StringValue(expression.Expression)

		result = append(result, m)
	}

	return result
}

func flattenMetrics(metrics []*aws.Metrics) []interface{} {
	result := make([]interface{}, 0, len(metrics))
	for _, metric := range metrics {
		m := make(map[string]interface{})

		m[string(ExtendedStatistic)] = spotinst.StringValue(metric.ExtendedStatistic)
		m[string(MetricName)] = spotinst.StringValue(metric.MetricName)
		m[string(Name)] = spotinst.StringValue(metric.Name)
		m[string(Namespace)] = spotinst.StringValue(metric.Namespace)
		m[string(Statistic)] = spotinst.StringValue(metric.Statistic)
		m[string(Unit)] = spotinst.StringValue(metric.Unit)

		if metric.Dimensions != nil && len(metric.Dimensions) > 0 {
			dimMap := make([]interface{}, 0, len(metric.Dimensions))
			for _, dimension := range metric.Dimensions {
				d := make(map[string]interface{})
				d[string(DimensionName)] = spotinst.StringValue(dimension.Name)
				d[string(DimensionValue)] = spotinst.StringValue(dimension.Value)

				if (d[string(DimensionName)] != nil) && (d[string(DimensionValue)] != nil) {
					dimMap = append(dimMap, d)
				}
			}
			m[string(Dimensions)] = dimMap
		}

		result = append(result, m)
	}

	return result
}

func flattenAction(action *aws.Action) []interface{} {
	result := make(map[string]interface{})

	result[string(Adjustment)] = spotinst.StringValue(action.Adjustment)
	result[string(MaxTargetCapacity)] = spotinst.StringValue(action.MaxTargetCapacity)
	result[string(Maximum)] = spotinst.StringValue(action.Maximum)
	result[string(MinTargetCapacity)] = spotinst.StringValue(action.MinTargetCapacity)
	result[string(Minimum)] = spotinst.StringValue(action.Minimum)
	result[string(Target)] = spotinst.StringValue(action.Target)
	result[string(Type)] = spotinst.StringValue(action.Type)

	return []interface{}{result}
}
