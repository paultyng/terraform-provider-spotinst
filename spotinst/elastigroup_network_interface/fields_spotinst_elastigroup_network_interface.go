package elastigroup_network_interface

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/spotinst/spotinst-sdk-go/service/elastigroup/providers/aws"
	"github.com/spotinst/spotinst-sdk-go/spotinst"
	"github.com/terraform-providers/terraform-provider-spotinst/spotinst/commons"
)

//-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
//            Setup
//-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
func Setup(fieldsMap map[commons.FieldName]*commons.GenericField) {

	fieldsMap[NetworkInterface] = commons.NewGenericField(
		commons.ElastigroupNetworkInterface,
		NetworkInterface,
		&schema.Schema{
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					string(Description): &schema.Schema{
						Type:     schema.TypeString,
						Required: true,
					},

					string(DeviceIndex): &schema.Schema{
						Type:     schema.TypeInt,
						Required: true,
					},

					string(SecondaryPrivateIpAddressCount): &schema.Schema{
						Type:     schema.TypeInt,
						Optional: true,
					},

					string(AssociatePublicIpAddress): &schema.Schema{
						Type:     schema.TypeBool,
						Optional: true,
					},

					string(DeleteOnTermination): &schema.Schema{
						Type:     schema.TypeBool,
						Optional: true,
						Computed: true,
					},

					string(NetworkInterfaceId): &schema.Schema{
						Type:     schema.TypeString,
						Optional: true,
					},

					string(PrivateIpAddress): &schema.Schema{
						Type:     schema.TypeString,
						Optional: true,
					},
				},
			},
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			elastigroup := resourceObject.(*aws.Group)
			var value []interface{} = nil
			if elastigroup.Compute != nil && elastigroup.Compute.LaunchSpecification != nil &&
				elastigroup.Compute.LaunchSpecification.NetworkInterfaces != nil {
				networkInterfaces := elastigroup.Compute.LaunchSpecification.NetworkInterfaces
				value = flattenAWSGroupNetworkInterfaces(networkInterfaces)
			}
			if value != nil {
				if err := resourceData.Set(string(NetworkInterface), value); err != nil {
					return fmt.Errorf(string(commons.FailureFieldReadPattern), string(NetworkInterface), err)
				}
			} else {
				if err := resourceData.Set(string(NetworkInterface), []*aws.NetworkInterface{}); err != nil {
					return fmt.Errorf(string(commons.FailureFieldReadPattern), string(NetworkInterface), err)
				}
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			elastigroup := resourceObject.(*aws.Group)
			if v, ok := resourceData.GetOk(string(NetworkInterface)); ok {
				if interfaces, err := expandAWSGroupNetworkInterfaces(v); err != nil {
					return err
				} else {
					elastigroup.Compute.LaunchSpecification.SetNetworkInterfaces(interfaces)
				}
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			elastigroup := resourceObject.(*aws.Group)
			var value []*aws.NetworkInterface = nil
			if v, ok := resourceData.GetOk(string(NetworkInterface)); ok {
				if interfaces, err := expandAWSGroupNetworkInterfaces(v); err != nil {
					return err
				} else {
					value = interfaces
				}
			}
			elastigroup.Compute.LaunchSpecification.SetNetworkInterfaces(value)
			return nil
		},
		nil,
	)
}

//-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
//            Utils
//-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
func flattenAWSGroupNetworkInterfaces(networkInterfaces []*aws.NetworkInterface) []interface{} {
	result := make([]interface{}, 0, len(networkInterfaces))
	for _, iface := range networkInterfaces {
		m := make(map[string]interface{})
		m[string(AssociatePublicIpAddress)] = spotinst.BoolValue(iface.AssociatePublicIPAddress)
		m[string(DeleteOnTermination)] = spotinst.BoolValue(iface.DeleteOnTermination)
		m[string(Description)] = spotinst.StringValue(iface.Description)
		m[string(DeviceIndex)] = spotinst.IntValue(iface.DeviceIndex)
		m[string(NetworkInterfaceId)] = spotinst.StringValue(iface.ID)
		m[string(PrivateIpAddress)] = spotinst.StringValue(iface.PrivateIPAddress)
		m[string(SecondaryPrivateIpAddressCount)] = spotinst.IntValue(iface.SecondaryPrivateIPAddressCount)
		result = append(result, m)
	}
	return result
}

func expandAWSGroupNetworkInterfaces(data interface{}) ([]*aws.NetworkInterface, error) {
	list := data.(*schema.Set).List()
	interfaces := make([]*aws.NetworkInterface, 0, len(list))
	for _, item := range list {
		m := item.(map[string]interface{})
		networkInterface := &aws.NetworkInterface{}

		if v, ok := m[string(NetworkInterfaceId)].(string); ok && v != "" {
			networkInterface.SetId(spotinst.String(v))
		}

		if v, ok := m[string(Description)].(string); ok && v != "" {
			networkInterface.SetDescription(spotinst.String(v))
		}

		if v, ok := m[string(DeviceIndex)].(int); ok && v >= 0 {
			networkInterface.SetDeviceIndex(spotinst.Int(v))
		}

		if v, ok := m[string(SecondaryPrivateIpAddressCount)].(int); ok && v >= 0 {
			networkInterface.SetSecondaryPrivateIPAddressCount(spotinst.Int(v))
		}

		if v, ok := m[string(AssociatePublicIpAddress)].(bool); ok {
			networkInterface.SetAssociatePublicIPAddress(spotinst.Bool(v))
		}

		if v, ok := m[string(DeleteOnTermination)].(bool); ok {
			networkInterface.SetDeleteOnTermination(spotinst.Bool(v))
		}

		if v, ok := m[string(PrivateIpAddress)].(string); ok && v != "" {
			networkInterface.SetPrivateIPAddress(spotinst.String(v))
		}

		interfaces = append(interfaces, networkInterface)
	}

	return interfaces, nil
}
