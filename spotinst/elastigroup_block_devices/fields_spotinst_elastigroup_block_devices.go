package elastigroup_block_devices

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
func Setup(fieldsMap map[commons.FieldName]*commons.GenericField) {

	fieldsMap[EbsBlockDevice] = commons.NewGenericField(
		commons.ElastigroupBlockDevices,
		EbsBlockDevice,
		&schema.Schema{
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					string(DeleteOnTermination): &schema.Schema{
						Type:     schema.TypeBool,
						Optional: true,
						Computed: true,
					},

					string(DeviceName): &schema.Schema{
						Type:     schema.TypeString,
						Required: true,
					},

					string(Encrypted): &schema.Schema{
						Type:     schema.TypeBool,
						Optional: true,
						Computed: true,
					},

					string(Iops): &schema.Schema{
						Type:     schema.TypeInt,
						Optional: true,
					},

					string(SnapshotId): &schema.Schema{
						Type:     schema.TypeString,
						Optional: true,
					},

					string(VolumeSize): &schema.Schema{
						Type:     schema.TypeInt,
						Optional: true,
					},

					string(VolumeType): &schema.Schema{
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
				},
			},
			Set: hashAWSGroupEBSBlockDevice,
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			elastigroup := resourceObject.(*aws.Group)
			var ebsBlock []interface{} = nil
			if elastigroup.Compute != nil && elastigroup.Compute.LaunchSpecification != nil &&
				elastigroup.Compute.LaunchSpecification.BlockDeviceMappings != nil {

				blockDeviceMappings := elastigroup.Compute.LaunchSpecification.BlockDeviceMappings
				ebsBlock = flattenAWSGroupEBSBlockDevices(blockDeviceMappings)
			}
			if ebsBlock != nil {
				if err := resourceData.Set(string(EbsBlockDevice), ebsBlock); err != nil {
					return fmt.Errorf(string(commons.FailureFieldReadPattern), string(EbsBlockDevice), err)
				}
			} else {
				if err := resourceData.Set(string(EbsBlockDevice), []*aws.BlockDeviceMapping{}); err != nil {
					return fmt.Errorf(string(commons.FailureFieldReadPattern), string(EbsBlockDevice), err)
				}
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			elastigroup := resourceObject.(*aws.Group)
			if v, ok := resourceData.GetOk(string(EbsBlockDevice)); ok {
				if ebsDevices, err := expandAWSGroupEBSBlockDevices(v); err != nil {
					return err
				} else {
					if existingEphemeralDevices, err := extractBlockDevices(EphemeralBlockDevice, elastigroup, resourceData); err != nil {
						return err
					} else if existingEphemeralDevices != nil && len(existingEphemeralDevices) > 0 {
						ebsDevices = append(existingEphemeralDevices, ebsDevices...)
					}
					elastigroup.Compute.LaunchSpecification.SetBlockDeviceMappings(ebsDevices)
				}
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			elastigroup := resourceObject.(*aws.Group)
			if err := onUpdateBlockDevice(elastigroup, resourceData); err != nil {
				return err
			}
			return nil
		},
		nil,
	)

	fieldsMap[EphemeralBlockDevice] = commons.NewGenericField(
		commons.ElastigroupBlockDevices,
		EphemeralBlockDevice,
		&schema.Schema{
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					string(DeviceName): &schema.Schema{
						Type:     schema.TypeString,
						Required: true,
					},

					string(VirtualName): &schema.Schema{
						Type:     schema.TypeString,
						Required: true,
					},
				},
			},
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			elastigroup := resourceObject.(*aws.Group)
			var ephemeralBlock []interface{} = nil
			if elastigroup.Compute != nil && elastigroup.Compute.LaunchSpecification != nil &&
				elastigroup.Compute.LaunchSpecification.BlockDeviceMappings != nil {

				blockDeviceMappings := elastigroup.Compute.LaunchSpecification.BlockDeviceMappings
				ephemeralBlock = flattenAWSGroupEphemeralBlockDevices(blockDeviceMappings)
			}
			if ephemeralBlock != nil {
				if err := resourceData.Set(string(EphemeralBlockDevice), ephemeralBlock); err != nil {
					return fmt.Errorf(string(commons.FailureFieldReadPattern), string(EphemeralBlockDevice), err)
				}
			} else {
				if err := resourceData.Set(string(EphemeralBlockDevice), []*aws.BlockDeviceMapping{}); err != nil {
					return fmt.Errorf(string(commons.FailureFieldReadPattern), string(EphemeralBlockDevice), err)
				}
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			elastigroup := resourceObject.(*aws.Group)
			if v, ok := resourceData.GetOk(string(EphemeralBlockDevice)); ok {
				if ephemeralDevices, err := expandAWSGroupEphemeralBlockDevices(v); err != nil {
					return err
				} else {
					if existingEBSDevices, err := extractBlockDevices(EbsBlockDevice, elastigroup, resourceData); err == nil {
						return err
					} else if existingEBSDevices != nil && len(existingEBSDevices) > 0 {
						ephemeralDevices = append(existingEBSDevices, ephemeralDevices...)
					}
					elastigroup.Compute.LaunchSpecification.SetBlockDeviceMappings(ephemeralDevices)
				}
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			elastigroup := resourceObject.(*aws.Group)
			if err := onUpdateBlockDevice(elastigroup, resourceData); err != nil {
				return err
			}
			return nil
		},
		nil,
	)
}

//-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
//            Utils
//-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
func hashAWSGroupEBSBlockDevice(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	buf.WriteString(fmt.Sprintf("%s-", m[string(DeviceName)].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m[string(SnapshotId)].(string)))
	buf.WriteString(fmt.Sprintf("%d-", m[string(VolumeSize)].(int)))
	buf.WriteString(fmt.Sprintf("%t-", m[string(DeleteOnTermination)].(bool)))
	buf.WriteString(fmt.Sprintf("%t-", m[string(Encrypted)].(bool)))
	buf.WriteString(fmt.Sprintf("%d-", m[string(Iops)].(int)))
	return hashcode.String(buf.String())
}

func flattenAWSGroupEBSBlockDevices(devices []*aws.BlockDeviceMapping) []interface{} {
	result := make([]interface{}, 0, len(devices))
	for _, dev := range devices {
		if dev.EBS != nil {
			m := make(map[string]interface{})
			m[string(DeviceName)] = spotinst.StringValue(dev.DeviceName)
			m[string(DeleteOnTermination)] = spotinst.BoolValue(dev.EBS.DeleteOnTermination)
			m[string(Encrypted)] = spotinst.BoolValue(dev.EBS.Encrypted)
			m[string(Iops)] = spotinst.IntValue(dev.EBS.IOPS)
			m[string(SnapshotId)] = spotinst.StringValue(dev.EBS.SnapshotID)
			m[string(VolumeType)] = spotinst.StringValue(dev.EBS.VolumeType)
			m[string(VolumeSize)] = spotinst.IntValue(dev.EBS.VolumeSize)
			result = append(result, m)
		}
	}
	return result
}

func flattenAWSGroupEphemeralBlockDevices(devices []*aws.BlockDeviceMapping) []interface{} {
	result := make([]interface{}, 0, len(devices))
	for _, dev := range devices {
		if dev.EBS == nil {
			m := make(map[string]interface{})
			m[string(DeviceName)] = spotinst.StringValue(dev.DeviceName)
			m[string(VirtualName)] = spotinst.StringValue(dev.VirtualName)
			result = append(result, m)
		}
	}
	return result
}

func expandAWSGroupEBSBlockDevices(data interface{}) ([]*aws.BlockDeviceMapping, error) {
	list := data.(*schema.Set).List()
	devices := make([]*aws.BlockDeviceMapping, 0, len(list))
	for _, item := range list {
		m := item.(map[string]interface{})
		device := &aws.BlockDeviceMapping{EBS: &aws.EBS{}}

		if v, ok := m[string(DeviceName)].(string); ok && v != "" {
			device.SetDeviceName(spotinst.String(v))
		}

		if v, ok := m[string(DeleteOnTermination)].(bool); ok {
			device.EBS.SetDeleteOnTermination(spotinst.Bool(v))
		}

		if v, ok := m[string(Encrypted)].(bool); ok && v != false {
			device.EBS.SetEncrypted(spotinst.Bool(v))
		}

		if v, ok := m[string(SnapshotId)].(string); ok && v != "" {
			device.EBS.SetSnapshotId(spotinst.String(v))
		}

		if v, ok := m[string(VolumeType)].(string); ok && v != "" {
			device.EBS.SetVolumeType(spotinst.String(v))
		}

		if v, ok := m[string(VolumeSize)].(int); ok && v > 0 {
			device.EBS.SetVolumeSize(spotinst.Int(v))
		}

		if v, ok := m[string(Iops)].(int); ok && v > 0 {
			device.EBS.SetIOPS(spotinst.Int(v))
		}
		devices = append(devices, device)
	}

	return devices, nil
}

func expandAWSGroupEphemeralBlockDevices(data interface{}) ([]*aws.BlockDeviceMapping, error) {
	list := data.(*schema.Set).List()
	devices := make([]*aws.BlockDeviceMapping, 0, len(list))
	for _, item := range list {
		m := item.(map[string]interface{})
		device := &aws.BlockDeviceMapping{}

		if v, ok := m[string(DeviceName)].(string); ok && v != "" {
			device.SetDeviceName(spotinst.String(v))
		}

		if v, ok := m[string(VirtualName)].(string); ok && v != "" {
			device.SetVirtualName(spotinst.String(v))
		}
		devices = append(devices, device)
	}

	return devices, nil
}

func extractBlockDevices(
	fieldName commons.FieldName,
	elastigroup *aws.Group,
	resourceData *schema.ResourceData) ([]*aws.BlockDeviceMapping, error) {

	result := elastigroup.Compute.LaunchSpecification.BlockDeviceMappings

	var ebsDevices []*aws.BlockDeviceMapping = nil
	var ephemeralDevices []*aws.BlockDeviceMapping = nil

	if result != nil && len(result) > 0 {
		for _, device := range result {
			if device.EBS != nil {
				ebsDevices = append(ebsDevices, device)
			} else {
				ephemeralDevices = append(ephemeralDevices, device)
			}
		}
	}

	if v, ok := resourceData.GetOk(string(EbsBlockDevice)); ok && fieldName == EbsBlockDevice {
		if tfEbsDevices, err := expandAWSGroupEBSBlockDevices(v); err != nil {
			return nil, err
		} else {
			ebsDevices = append(tfEbsDevices, ebsDevices...)
		}
	}

	if v, ok := resourceData.GetOk(string(EphemeralBlockDevice)); ok && fieldName == EphemeralBlockDevice {
		if tfEphemeralDevices, err := expandAWSGroupEphemeralBlockDevices(v); err != nil {
			return nil, err
		} else {
			ephemeralDevices = append(tfEphemeralDevices, ephemeralDevices...)
		}
	}

	if fieldName == EbsBlockDevice {
		return ebsDevices, nil
	} else {
		return ephemeralDevices, nil
	}
}

var ephemeralBlockDeviceUpdated = false
var ebsBlockDeviceUpdated = false

func onUpdateBlockDevice(elastigroup *aws.Group, resourceData *schema.ResourceData) error {
	var ebsNullify = false
	var ephemeralNullify = false
	if !ebsBlockDeviceUpdated {
		if tfEBSDevices, err := extractBlockDevices(EbsBlockDevice, elastigroup, resourceData); err != nil {
			return err
		} else if tfEBSDevices != nil && len(tfEBSDevices) > 0 {
			existingMappings := elastigroup.Compute.LaunchSpecification.BlockDeviceMappings
			if existingMappings != nil && len(existingMappings) > 0 {
				tfEBSDevices = append(existingMappings, tfEBSDevices...)
			}
			elastigroup.Compute.LaunchSpecification.SetBlockDeviceMappings(tfEBSDevices)
		} else {
			ebsNullify = true
		}
		ebsBlockDeviceUpdated = true
	}
	if !ephemeralBlockDeviceUpdated {
		if tfEphemeralDevices, err := extractBlockDevices(EphemeralBlockDevice, elastigroup, resourceData); err != nil {
			return err
		} else if tfEphemeralDevices != nil && len(tfEphemeralDevices) > 0 {
			existingMappings := elastigroup.Compute.LaunchSpecification.BlockDeviceMappings
			if existingMappings != nil && len(existingMappings) > 0 {
				tfEphemeralDevices = append(existingMappings, tfEphemeralDevices...)
			}
			elastigroup.Compute.LaunchSpecification.SetBlockDeviceMappings(tfEphemeralDevices)
		} else {
			ephemeralNullify = true
		}
		ephemeralBlockDeviceUpdated = true
	}
	// Both fields share the same object structure, we need to nullify if and only if there are no items
	// from both types
	if ebsNullify && ephemeralNullify {
		elastigroup.Compute.LaunchSpecification.SetBlockDeviceMappings(nil)
	}
	return nil
}
