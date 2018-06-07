package elastigroup_aws

import "github.com/terraform-providers/terraform-provider-spotinst/spotinst/commons"

type BalancerType string

const (
	BalancerTypeClassic         BalancerType = "CLASSIC"
	BalancerTypeTargetGroup     BalancerType = "TARGET_GROUP"
	BalancerTypeMultaiTargetSet BalancerType = "MULTAI_TARGET_SET"
)

type TagField string

const (
	TagKey   TagField = "key"
	TagValue TagField = "value"
)

const (
	ResourceOnRoll commons.LogFormat = "onRoll() -> started for group %v..."
)

const (
	Name        commons.FieldName = "name"
	Description commons.FieldName = "description"
	Product     commons.FieldName = "product"

	MaxSize         commons.FieldName = "max_size"
	MinSize         commons.FieldName = "min_size"
	DesiredCapacity commons.FieldName = "desired_capacity"
	CapacityUnit    commons.FieldName = "capacity_unit"

	HealthCheckGracePeriod                        commons.FieldName = "health_check_grace_period"
	HealthCheckType                               commons.FieldName = "health_check_type"
	HealthCheckUnhealthyDurationBeforeReplacement commons.FieldName = "health_check_unhealthy_duration_before_replacement"

	Region               commons.FieldName = "region"
	SubnetIds            commons.FieldName = "subnet_ids"
	AvailabilityZones    commons.FieldName = "availability_zones"
	ElasticLoadBalancers commons.FieldName = "elastic_load_balancers"
	TargetGroupArns      commons.FieldName = "target_group_arns"
	MultaiTargetSets     commons.FieldName = "multai_target_sets"
	MultaiTargetSetId    commons.FieldName = "target_set_id"
	MultaiBalancerId     commons.FieldName = "balancer_id"
	Tags                 commons.FieldName = "tags"

	RevertToSpot commons.FieldName = "revert_to_spot"
	PerformAt    commons.FieldName = "perform_at"
	TimeWindow   commons.FieldName = "time_windows"

	// ***********************************************************************
	// ********************* Spotinst Unique Properties **********************
	// ***********************************************************************

	ElasticIps commons.FieldName = "elastic_ips"

	Signal        commons.FieldName = "signal"
	SignalName    commons.FieldName = "name"
	SignalTimeout commons.FieldName = "timeout"

	UpdatePolicy         commons.FieldName = "update_policy"
	ShouldResumeStateful commons.FieldName = "should_resume_stateful"
	ShouldRoll           commons.FieldName = "should_roll"

	RollConfig          commons.FieldName = "roll_config"
	BatchSizePercentage commons.FieldName = "batch_size_percentage"
	GracePeriod         commons.FieldName = "grace_period"
)
