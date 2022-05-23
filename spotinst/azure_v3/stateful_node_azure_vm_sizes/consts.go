package stateful_node_azure_vm_sizes

import "github.com/spotinst/terraform-provider-spotinst/spotinst/commons"

const (
	OnDemandSizes      commons.FieldName = "od_sizes"
	SpotSizes          commons.FieldName = "spot_sizes"
	PreferredSpotSizes commons.FieldName = "preferred_spot_sizes"
)
