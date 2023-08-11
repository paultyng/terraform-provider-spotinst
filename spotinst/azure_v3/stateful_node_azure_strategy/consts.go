package stateful_node_azure_strategy

import "github.com/spotinst/terraform-provider-spotinst/spotinst/commons"

const (
	Strategy            commons.FieldName = "strategy"
	PreferredLifecycle  commons.FieldName = "preferred_life_cycle"
	FallbackToOnDemand  commons.FieldName = "fallback_to_on_demand"
	DrainingTimeout     commons.FieldName = "draining_timeout"
	RevertToSpot        commons.FieldName = "revert_to_spot"
	PerformAt           commons.FieldName = "perform_at"
	OptimizationWindows commons.FieldName = "optimization_windows"
)

const (
	Signal  commons.FieldName = "signal"
	Type    commons.FieldName = "type"
	Timeout commons.FieldName = "timeout"
)

const (
	CapacityReservation       commons.FieldName = "capacity_reservation"
	ShouldUtilize             commons.FieldName = "should_utilize"
	UtilizationStrategy       commons.FieldName = "utilization_strategy"
	CapacityReservationGroups commons.FieldName = "capacity_reservation_groups"
	CRGName                   commons.FieldName = "crg_name"
	CRGResourceGroupName      commons.FieldName = "crg_resource_group_name"
	CRGShouldPrioritize       commons.FieldName = "crg_should_prioritize"
)
