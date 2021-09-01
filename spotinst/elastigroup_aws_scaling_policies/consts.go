package elastigroup_aws_scaling_policies

import "github.com/spotinst/terraform-provider-spotinst/spotinst/commons"

type DimensionField string

const (
	ScalingUpPolicy     commons.FieldName = "scaling_up_policy"
	ScalingDownPolicy   commons.FieldName = "scaling_down_policy"
	ScalingTargetPolicy commons.FieldName = "scaling_target_policy"
	MultipleMetrics     commons.FieldName = "multiple_metrics"

	Expressions commons.FieldName = "expressions"
	Metrics     commons.FieldName = "metrics"

	Expression commons.FieldName = "expression"
	Name       commons.FieldName = "name"

	PolicyName commons.FieldName = "policy_name"
	MetricName commons.FieldName = "metric_name"
	Namespace  commons.FieldName = "namespace"
	Source     commons.FieldName = "source"
	Statistic  commons.FieldName = "statistic"
	Unit       commons.FieldName = "unit"
	Cooldown   commons.FieldName = "cooldown"
	Dimensions commons.FieldName = "dimensions"

	Threshold           commons.FieldName = "threshold"
	Adjustment          commons.FieldName = "adjustment"
	MinTargetCapacity   commons.FieldName = "min_target_capacity"
	MaxTargetCapacity   commons.FieldName = "max_target_capacity"
	Operator            commons.FieldName = "operator"
	EvaluationPeriods   commons.FieldName = "evaluation_periods"
	Period              commons.FieldName = "period"
	Minimum             commons.FieldName = "minimum"
	Maximum             commons.FieldName = "maximum"
	Target              commons.FieldName = "target"
	ActionType          commons.FieldName = "action_type"
	IsEnabled           commons.FieldName = "is_enabled"
	PredictiveMode      commons.FieldName = "predictive_mode"
	MaxCapacityPerScale commons.FieldName = "max_capacity_per_scale"
	StepAdjustments     commons.FieldName = "step_adjustments"
	Action              commons.FieldName = "action"
	Type                commons.FieldName = "type"
	ExtendedStatistic   commons.FieldName = "extended_statistic"

	DimensionName  DimensionField = "name"
	DimensionValue DimensionField = "value"
)
