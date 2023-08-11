package stateful_node_azure

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/spotinst/spotinst-sdk-go/spotinst"
	"github.com/spotinst/terraform-provider-spotinst/spotinst/commons"
)

func Setup(fieldsMap map[commons.FieldName]*commons.GenericField) {

	fieldsMap[Name] = commons.NewGenericField(
		commons.StatefulNodeAzure,
		Name,
		&schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			snWrapper := resourceObject.(*commons.StatefulNodeAzureV3Wrapper)
			statefulNode := snWrapper.GetStatefulNode()
			var value *string = nil
			if statefulNode.Name != nil {
				value = statefulNode.Name
			}
			if err := resourceData.Set(string(Name), value); err != nil {
				return fmt.Errorf(string(commons.FailureFieldReadPattern), string(Name), err)
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			snWrapper := resourceObject.(*commons.StatefulNodeAzureV3Wrapper)
			statefulNode := snWrapper.GetStatefulNode()
			statefulNode.SetName(spotinst.String(resourceData.Get(string(Name)).(string)))
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			snWrapper := resourceObject.(*commons.StatefulNodeAzureV3Wrapper)
			statefulNode := snWrapper.GetStatefulNode()
			statefulNode.SetName(spotinst.String(resourceData.Get(string(Name)).(string)))
			return nil
		},
		nil,
	)

	fieldsMap[Region] = commons.NewGenericField(
		commons.StatefulNodeAzure,
		Region,
		&schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			snWrapper := resourceObject.(*commons.StatefulNodeAzureV3Wrapper)
			statefulNode := snWrapper.GetStatefulNode()
			var value *string = nil
			if statefulNode.Region != nil {
				value = statefulNode.Region
			}
			if err := resourceData.Set(string(Region), value); err != nil {
				return fmt.Errorf(string(commons.FailureFieldReadPattern), string(Region), err)
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			snWrapper := resourceObject.(*commons.StatefulNodeAzureV3Wrapper)
			statefulNode := snWrapper.GetStatefulNode()
			if v, ok := resourceData.GetOk(string(Region)); ok {
				statefulNode.SetRegion(spotinst.String(v.(string)))
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			snWrapper := resourceObject.(*commons.StatefulNodeAzureV3Wrapper)
			statefulNode := snWrapper.GetStatefulNode()
			if v, ok := resourceData.GetOk(string(Region)); ok {
				statefulNode.SetRegion(spotinst.String(v.(string)))
			}
			return nil
		},
		nil,
	)

	fieldsMap[ResourceGroupName] = commons.NewGenericField(
		commons.StatefulNodeAzure,
		ResourceGroupName,
		&schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			snWrapper := resourceObject.(*commons.StatefulNodeAzureV3Wrapper)
			statefulNode := snWrapper.GetStatefulNode()
			var value *string = nil
			if statefulNode.ResourceGroupName != nil {
				value = statefulNode.ResourceGroupName
			}
			if err := resourceData.Set(string(ResourceGroupName), value); err != nil {
				return fmt.Errorf(string(commons.FailureFieldReadPattern), string(ResourceGroupName), err)
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			snWrapper := resourceObject.(*commons.StatefulNodeAzureV3Wrapper)
			statefulNode := snWrapper.GetStatefulNode()
			statefulNode.SetResourceGroupName(spotinst.String(resourceData.Get(string(ResourceGroupName)).(string)))
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			err := fmt.Errorf(string(commons.FieldUpdateNotAllowedPattern),
				string(ResourceGroupName))
			return err
		},
		nil,
	)

	fieldsMap[Description] = commons.NewGenericField(
		commons.StatefulNodeAzure,
		Description,
		&schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			snWrapper := resourceObject.(*commons.StatefulNodeAzureV3Wrapper)
			statefulNode := snWrapper.GetStatefulNode()
			var value *string = nil
			if statefulNode.Description != nil {
				value = statefulNode.Description
			}
			if err := resourceData.Set(string(Description), value); err != nil {
				return fmt.Errorf(string(commons.FailureFieldReadPattern), string(Description), err)
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			snWrapper := resourceObject.(*commons.StatefulNodeAzureV3Wrapper)
			statefulNode := snWrapper.GetStatefulNode()
			if v, ok := resourceData.GetOk(string(Description)); ok && v != "" {
				statefulNode.SetDescription(spotinst.String(v.(string)))
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			snWrapper := resourceObject.(*commons.StatefulNodeAzureV3Wrapper)
			statefulNode := snWrapper.GetStatefulNode()
			if v, ok := resourceData.GetOk(string(Description)); ok && v != "" {
				statefulNode.SetDescription(spotinst.String(v.(string)))
			}
			return nil
		},
		nil,
	)

	fieldsMap[OS] = commons.NewGenericField(
		commons.StatefulNodeAzure,
		OS,
		&schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			snWrapper := resourceObject.(*commons.StatefulNodeAzureV3Wrapper)
			statefulNode := snWrapper.GetStatefulNode()
			var value *string = nil
			if statefulNode.Compute != nil && statefulNode.Compute.OS != nil {
				value = statefulNode.Compute.OS
			}
			if err := resourceData.Set(string(OS), value); err != nil {
				return fmt.Errorf(string(commons.FailureFieldReadPattern), string(OS), err)
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			snWrapper := resourceObject.(*commons.StatefulNodeAzureV3Wrapper)
			statefulNode := snWrapper.GetStatefulNode()
			if v, ok := resourceData.GetOk(string(OS)); ok && v != "" {
				statefulNode.Compute.SetOS(spotinst.String(v.(string)))
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			snWrapper := resourceObject.(*commons.StatefulNodeAzureV3Wrapper)
			statefulNode := snWrapper.GetStatefulNode()
			if v, ok := resourceData.GetOk(string(OS)); ok && v != "" {
				statefulNode.Compute.SetOS(spotinst.String(v.(string)))
			}
			return nil
		},
		nil,
	)

	fieldsMap[Zones] = commons.NewGenericField(
		commons.StatefulNodeAzure,
		Zones,
		&schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			snWrapper := resourceObject.(*commons.StatefulNodeAzureV3Wrapper)
			statefulNode := snWrapper.GetStatefulNode()

			var result []string
			if statefulNode.Compute != nil && statefulNode.Compute.Zones != nil {
				result = append(result, statefulNode.Compute.Zones...)
				if err := resourceData.Set(string(Zones), result); err != nil {
					return fmt.Errorf(string(commons.FailureFieldReadPattern), string(Zones), err)
				}
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			snWrapper := resourceObject.(*commons.StatefulNodeAzureV3Wrapper)
			statefulNode := snWrapper.GetStatefulNode()
			if v, ok := resourceData.Get(string(Zones)).([]interface{}); ok && v != nil {
				if zones, err := expandZones(v); err != nil {
					return err
				} else {
					statefulNode.Compute.SetZones(zones)
				}
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			snWrapper := resourceObject.(*commons.StatefulNodeAzureV3Wrapper)
			statefulNode := snWrapper.GetStatefulNode()
			if v, ok := resourceData.Get(string(Zones)).([]interface{}); ok && v != nil {
				if zones, err := expandZones(v); err != nil {
					return err
				} else {
					statefulNode.Compute.SetZones(zones)
				}
			}
			return nil
		},
		nil,
	)

	fieldsMap[PreferredZone] = commons.NewGenericField(
		commons.StatefulNodeAzure,
		PreferredZone,
		&schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			snWrapper := resourceObject.(*commons.StatefulNodeAzureV3Wrapper)
			statefulNode := snWrapper.GetStatefulNode()
			var value *string = nil
			if statefulNode.Compute != nil && statefulNode.Compute.OS != nil {
				value = statefulNode.Compute.PreferredZone
			}
			if err := resourceData.Set(string(PreferredZone), value); err != nil {
				return fmt.Errorf(string(commons.FailureFieldReadPattern), string(PreferredZone), err)
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			snWrapper := resourceObject.(*commons.StatefulNodeAzureV3Wrapper)
			statefulNode := snWrapper.GetStatefulNode()
			if v, ok := resourceData.GetOk(string(PreferredZone)); ok && v != "" {
				statefulNode.Compute.SetPreferredZone(spotinst.String(v.(string)))
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			snWrapper := resourceObject.(*commons.StatefulNodeAzureV3Wrapper)
			statefulNode := snWrapper.GetStatefulNode()
			if v, ok := resourceData.GetOk(string(PreferredZone)); ok && v != "" {
				statefulNode.Compute.SetPreferredZone(spotinst.String(v.(string)))
			}
			return nil
		},
		nil,
	)

	fieldsMap[Delete] = commons.NewGenericField(
		commons.StatefulNodeAzure,
		Delete,
		&schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					string(ShouldTerminateVm): {
						Type:     schema.TypeBool,
						Required: true,
					},
					string(ShouldDeallocateNetwork): {
						Type:     schema.TypeBool,
						Required: true,
					},
					string(NetworkTTLInHours): {
						Type:     schema.TypeInt,
						Optional: true,
					},
					string(ShouldDeallocateDisk): {
						Type:     schema.TypeBool,
						Required: true,
					},
					string(DiskTTLInHours): {
						Type:     schema.TypeInt,
						Optional: true,
					},
					string(ShouldDeallocateSnapshot): {
						Type:     schema.TypeBool,
						Required: true,
					},
					string(SnapshotTTLInHours): {
						Type:     schema.TypeInt,
						Optional: true,
					},
					string(ShouldDeallocatePublicIP): {
						Type:     schema.TypeBool,
						Required: true,
					},
					string(PublicIPTTLInHours): {
						Type:     schema.TypeInt,
						Optional: true,
					},
				},
			},
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			return nil
		},
		nil,
	)

	fieldsMap[AttachDataDisk] = commons.NewGenericField(
		commons.StatefulNodeAzure,
		AttachDataDisk,
		&schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					string(AttachDataDiskName): {
						Type:     schema.TypeString,
						Required: true,
					},
					string(AttachDataDiskResourceGroupName): {
						Type:     schema.TypeString,
						Required: true,
					},
					string(AttachStorageAccountType): {
						Type:     schema.TypeString,
						Required: true,
					},
					string(AttachSizeGB): {
						Type:     schema.TypeInt,
						Required: true,
					},
					string(AttachLUN): {
						Type:     schema.TypeInt,
						Optional: true,
					},
					string(AttachZone): {
						Type:     schema.TypeString,
						Optional: true,
					},
				},
			},
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			return nil
		},
		nil,
	)

	fieldsMap[DetachDataDisk] = commons.NewGenericField(
		commons.StatefulNodeAzure,
		DetachDataDisk,
		&schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					string(DetachDataDiskName): {
						Type:     schema.TypeString,
						Required: true,
					},
					string(DetachDataDiskResourceGroupName): {
						Type:     schema.TypeString,
						Required: true,
					},
					string(DetachShouldDeallocate): {
						Type:     schema.TypeBool,
						Required: true,
					},
					string(DetachTTLInHours): {
						Type:     schema.TypeInt,
						Optional: true,
					},
				},
			},
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			return nil
		},
		nil,
	)

	fieldsMap[UpdateState] = commons.NewGenericField(
		commons.StatefulNodeAzure,
		UpdateState,
		&schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					string(State): {
						Type:     schema.TypeString,
						Required: true,
					},
				},
			},
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			return nil
		},
		nil,
	)

	fieldsMap[ImportVM] = commons.NewGenericField(
		commons.StatefulNodeAzure,
		ImportVM,
		&schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					string(ImportVMResourceGroupName): {
						Type:     schema.TypeString,
						Required: true,
					},
					string(ImportVMOriginalVMName): {
						Type:     schema.TypeString,
						Required: true,
					},
					string(ImportVMDrainingTimeout): {
						Type:     schema.TypeInt,
						Optional: true,
					},
					string(ImportVMResourcesRetentionTime): {
						Type:     schema.TypeInt,
						Optional: true,
					},
				},
			},
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			return nil
		},
		nil,
	)

}

func expandZones(data interface{}) ([]string, error) {
	list := data.([]interface{})
	result := make([]string, 0, len(list))

	for _, v := range list {
		if zone, ok := v.(string); ok && zone != "" {
			result = append(result, zone)
		}
	}

	return result, nil
}
