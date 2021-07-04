package ocean_gke_launch_spec

import "github.com/spotinst/terraform-provider-spotinst/spotinst/commons"

type LabelField string
type MetadataField string

const (
	LabelKey   LabelField = "key"
	LabelValue LabelField = "value"

	MetadataKey   MetadataField = "key"
	MetadataValue MetadataField = "value"

	TaintKey    MetadataField = "key"
	TaintValue  MetadataField = "value"
	TaintEffect MetadataField = "effect"
)

const (
	CPUPerUnit    commons.FieldName = "cpu_per_unit"
	GPUPerUnit    commons.FieldName = "gpu_per_unit"
	MemoryPerUnit commons.FieldName = "memory_per_unit"
	NumOfUnits    commons.FieldName = "num_of_units"
)

const (
	OceanId            commons.FieldName = "ocean_id"
	SourceImage        commons.FieldName = "source_image"
	Metadata           commons.FieldName = "metadata"
	Labels             commons.FieldName = "labels"
	Taints             commons.FieldName = "taints"
	AutoscaleHeadrooms commons.FieldName = "autoscale_headrooms"
	RestrictScaleDown  commons.FieldName = "restrict_scale_down"
	RootVolumeType     commons.FieldName = "root_volume_type"
	RootVolumeSizeInGB commons.FieldName = "root_volume_size"
	InstanceTypes      commons.FieldName = "instance_types"
)

const (
	NodePoolName commons.FieldName = "node_pool_name"
)
