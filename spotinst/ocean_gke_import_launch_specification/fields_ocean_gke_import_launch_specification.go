package ocean_gke_import_launch_specification

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/spotinst/spotinst-sdk-go/spotinst"
	"github.com/spotinst/terraform-provider-spotinst/spotinst/commons"
)

func Setup(fieldsMap map[commons.FieldName]*commons.GenericField) {

	fieldsMap[RootVolumeType] = commons.NewGenericField(
		commons.OceanGKEImportLaunchSpecification,
		RootVolumeType,
		&schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			clusterWrapper := resourceObject.(*commons.GKEImportClusterWrapper)
			cluster := clusterWrapper.GetCluster()
			var result *string = nil
			if cluster != nil && cluster.Compute != nil && cluster.Compute.LaunchSpecification != nil && cluster.Compute.LaunchSpecification.RootVolumeType != nil {
				result = cluster.Compute.LaunchSpecification.RootVolumeType
			}
			if result != nil {
				if err := resourceData.Set(string(RootVolumeType), result); err != nil {
					return fmt.Errorf(string(commons.FailureFieldReadPattern), string(RootVolumeType), err)
				}
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			clusterWrapper := resourceObject.(*commons.GKEImportClusterWrapper)
			cluster := clusterWrapper.GetCluster()
			var rootVolumeType *string = nil

			if cluster != nil && cluster.Compute != nil && cluster.Compute.LaunchSpecification != nil {
				// get rootVolumeType from previous import step.
				rootVolumeType = cluster.Compute.LaunchSpecification.RootVolumeType

				// get rootVolumeType from user configuration.
				if v, ok := resourceData.GetOk(string(RootVolumeType)); ok {
					rootVolumeType = spotinst.String(v.(string))

					if rootVolumeType != nil {
						cluster.Compute.LaunchSpecification.SetRootVolumeType(rootVolumeType)
					}
				}

			}

			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			clusterWrapper := resourceObject.(*commons.GKEImportClusterWrapper)
			cluster := clusterWrapper.GetCluster()
			var rootVolumeType *string = nil
			if v, ok := resourceData.GetOk(string(RootVolumeType)); ok && v != "" {
				rootVolumeType = spotinst.String(v.(string))
			}
			cluster.Compute.LaunchSpecification.SetRootVolumeType(rootVolumeType)
			return nil
		},
		nil,
	)
}
