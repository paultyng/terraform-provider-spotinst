package ocean_aws_launch_spec

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/spotinst/spotinst-sdk-go/service/ocean/providers/aws"
	"github.com/spotinst/spotinst-sdk-go/spotinst"
	"github.com/spotinst/terraform-provider-spotinst/spotinst/commons"
)

func Setup(fieldsMap map[commons.FieldName]*commons.GenericField) {

	fieldsMap[OceanID] = commons.NewGenericField(
		commons.OceanAWSLaunchSpec,
		OceanID,
		&schema.Schema{
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			launchSpecWrapper := resourceObject.(*commons.LaunchSpecWrapper)
			launchSpec := launchSpecWrapper.GetLaunchSpec()
			var value *string = nil
			if launchSpec.OceanID != nil {
				value = launchSpec.OceanID
			}
			if err := resourceData.Set(string(OceanID), value); err != nil {
				return fmt.Errorf(string(commons.FailureFieldReadPattern), string(OceanID), err)
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			launchSpecWrapper := resourceObject.(*commons.LaunchSpecWrapper)
			launchSpec := launchSpecWrapper.GetLaunchSpec()
			launchSpec.SetOceanId(spotinst.String(resourceData.Get(string(OceanID)).(string)))
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			launchSpecWrapper := resourceObject.(*commons.LaunchSpecWrapper)
			launchSpec := launchSpecWrapper.GetLaunchSpec()
			launchSpec.SetOceanId(spotinst.String(resourceData.Get(string(OceanID)).(string)))
			return nil
		},
		nil,
	)

	fieldsMap[ImageID] = commons.NewGenericField(
		commons.OceanAWSLaunchSpec,
		ImageID,
		&schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			launchSpecWrapper := resourceObject.(*commons.LaunchSpecWrapper)
			launchSpec := launchSpecWrapper.GetLaunchSpec()
			var value *string = nil
			if launchSpec.ImageID != nil {
				value = launchSpec.ImageID
			}
			if err := resourceData.Set(string(ImageID), value); err != nil {
				return fmt.Errorf(string(commons.FailureFieldReadPattern), string(ImageID), err)
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			launchSpecWrapper := resourceObject.(*commons.LaunchSpecWrapper)
			launchSpec := launchSpecWrapper.GetLaunchSpec()
			if value, ok := resourceData.GetOk(string(ImageID)); ok && value != nil {
				launchSpec.SetImageId(spotinst.String(resourceData.Get(string(ImageID)).(string)))
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			launchSpecWrapper := resourceObject.(*commons.LaunchSpecWrapper)
			launchSpec := launchSpecWrapper.GetLaunchSpec()
			var imageId *string = nil
			if value, ok := resourceData.GetOk(string(ImageID)); ok && value != nil {
				imageId = spotinst.String(resourceData.Get(string(ImageID)).(string))
			}
			launchSpec.SetImageId(imageId)
			return nil
		},
		nil,
	)

	fieldsMap[Name] = commons.NewGenericField(
		commons.OceanAWSLaunchSpec,
		Name,
		&schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			launchSpecWrapper := resourceObject.(*commons.LaunchSpecWrapper)
			launchSpec := launchSpecWrapper.GetLaunchSpec()
			var value *string = nil
			if launchSpec.Name != nil {
				value = launchSpec.Name
			}
			if err := resourceData.Set(string(Name), value); err != nil {
				return fmt.Errorf(string(commons.FailureFieldReadPattern), string(Name), err)
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			launchSpecWrapper := resourceObject.(*commons.LaunchSpecWrapper)
			launchSpec := launchSpecWrapper.GetLaunchSpec()
			if value, ok := resourceData.GetOk(string(Name)); ok && value != nil {
				launchSpec.SetName(spotinst.String(resourceData.Get(string(Name)).(string)))
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			launchSpecWrapper := resourceObject.(*commons.LaunchSpecWrapper)
			launchSpec := launchSpecWrapper.GetLaunchSpec()
			if value, ok := resourceData.GetOk(string(Name)); ok && value != nil {
				launchSpec.SetName(spotinst.String(resourceData.Get(string(Name)).(string)))
			}
			return nil
		},
		nil,
	)

	fieldsMap[UserData] = commons.NewGenericField(
		commons.OceanAWSLaunchSpec,
		UserData,
		&schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
				// Sometimes the EC2 API responds with the equivalent, empty SHA1 sum
				if (old == "da39a3ee5e6b4b0d3255bfef95601890afd80709" && new == "") ||
					(old == "" && new == "da39a3ee5e6b4b0d3255bfef95601890afd80709") {
					return true
				}
				return false
			},
			StateFunc: Base64StateFunc,
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			launchSpecWrapper := resourceObject.(*commons.LaunchSpecWrapper)
			launchSpec := launchSpecWrapper.GetLaunchSpec()
			var value = ""
			if launchSpec.UserData != nil {
				userData := launchSpec.UserData
				userDataValue := spotinst.StringValue(userData)
				if userDataValue != "" {
					if isBase64Encoded(resourceData.Get(string(UserData)).(string)) {
						value = userDataValue
					} else {
						decodedUserData, _ := base64.StdEncoding.DecodeString(userDataValue)
						value = string(decodedUserData)
					}
				}
			}
			if err := resourceData.Set(string(UserData), Base64StateFunc(value)); err != nil {
				return fmt.Errorf(string(commons.FailureFieldReadPattern), string(UserData), err)
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			launchSpecWrapper := resourceObject.(*commons.LaunchSpecWrapper)
			launchSpec := launchSpecWrapper.GetLaunchSpec()
			if v, ok := resourceData.Get(string(UserData)).(string); ok && v != "" {
				userData := spotinst.String(base64Encode(v))
				launchSpec.SetUserData(userData)
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			launchSpecWrapper := resourceObject.(*commons.LaunchSpecWrapper)
			launchSpec := launchSpecWrapper.GetLaunchSpec()
			var userData *string = nil
			if v, ok := resourceData.Get(string(UserData)).(string); ok && v != "" {
				userData = spotinst.String(base64Encode(v))
			}
			launchSpec.SetUserData(userData)
			return nil
		},
		nil,
	)

	fieldsMap[Labels] = commons.NewGenericField(
		commons.OceanAWSLaunchSpec,
		Labels,
		&schema.Schema{
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					string(LabelKey): {
						Type:     schema.TypeString,
						Required: true,
					},
					string(LabelValue): {
						Type:     schema.TypeString,
						Required: true,
					},
				},
			},
			Set: hashKV,
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			launchSpecWrapper := resourceObject.(*commons.LaunchSpecWrapper)
			launchSpec := launchSpecWrapper.GetLaunchSpec()
			var result []interface{} = nil
			if launchSpec.Labels != nil {
				labels := launchSpec.Labels
				result = flattenLabels(labels)
			}
			if result != nil {
				if err := resourceData.Set(string(Labels), result); err != nil {
					return fmt.Errorf(string(commons.FailureFieldReadPattern), string(Labels), err)
				}
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			launchSpecWrapper := resourceObject.(*commons.LaunchSpecWrapper)
			launchSpec := launchSpecWrapper.GetLaunchSpec()
			if value, ok := resourceData.GetOk(string(Labels)); ok {
				if labels, err := expandLabels(value); err != nil {
					return err
				} else {
					launchSpec.SetLabels(labels)
				}
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			launchSpecWrapper := resourceObject.(*commons.LaunchSpecWrapper)
			launchSpec := launchSpecWrapper.GetLaunchSpec()
			var labelList []*aws.Label = nil
			if value, ok := resourceData.GetOk(string(Labels)); ok {
				if labels, err := expandLabels(value); err != nil {
					return err
				} else {
					labelList = labels
				}
			}
			launchSpec.SetLabels(labelList)
			return nil
		},
		nil,
	)

	fieldsMap[Tags] = commons.NewGenericField(
		commons.OceanAWSLaunchConfiguration,
		Tags,
		&schema.Schema{
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					string(TagKey): {
						Type:     schema.TypeString,
						Required: true,
					},

					string(TagValue): {
						Type:     schema.TypeString,
						Required: true,
					},
				},
			},
			Set: hashKV,
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			launchSpecWrapper := resourceObject.(*commons.LaunchSpecWrapper)
			launchSpec := launchSpecWrapper.GetLaunchSpec()
			var result []interface{} = nil
			if launchSpec.Tags != nil {
				tags := launchSpec.Tags
				result = flattenTags(tags)
			}
			if result != nil {
				if err := resourceData.Set(string(Tags), result); err != nil {
					return fmt.Errorf(string(commons.FailureFieldReadPattern), string(Tags), err)
				}
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			launchSpecWrapper := resourceObject.(*commons.LaunchSpecWrapper)
			launchSpec := launchSpecWrapper.GetLaunchSpec()
			if value, ok := resourceData.GetOk(string(Tags)); ok {
				if tags, err := expandTags(value); err != nil {
					return err
				} else {
					launchSpec.SetTags(tags)
				}
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			launchSpecWrapper := resourceObject.(*commons.LaunchSpecWrapper)
			launchSpec := launchSpecWrapper.GetLaunchSpec()
			var tagsToAdd []*aws.Tag = nil
			if value, ok := resourceData.GetOk(string(Tags)); ok {
				if tags, err := expandTags(value); err != nil {
					return err
				} else {
					tagsToAdd = tags
				}
			}
			launchSpec.SetTags(tagsToAdd)
			return nil
		},
		nil,
	)

	fieldsMap[ElasticIpPool] = commons.NewGenericField(
		commons.OceanAWSLaunchSpec,
		ElasticIpPool,
		&schema.Schema{
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{

					string(TagSelector): {
						Type:     schema.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{

								string(TagSelectorKey): {
									Type:     schema.TypeString,
									Required: true,
								},

								string(TagSelectorValue): {
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
			launchSpecWrapper := resourceObject.(*commons.LaunchSpecWrapper)
			launchSpec := launchSpecWrapper.GetLaunchSpec()
			var result []interface{} = nil
			if launchSpec.ElasticIPPool != nil {
				elasticIpPool := launchSpec.ElasticIPPool
				result = flattenElasticIpPool(elasticIpPool)
			}
			if result != nil {
				if err := resourceData.Set(string(ElasticIpPool), result); err != nil {
					return fmt.Errorf(string(commons.FailureFieldReadPattern), string(ElasticIpPool), err)
				}
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			launchSpecWrapper := resourceObject.(*commons.LaunchSpecWrapper)
			launchSpec := launchSpecWrapper.GetLaunchSpec()
			if value, ok := resourceData.GetOk(string(ElasticIpPool)); ok {
				if elasticIpPool, err := expandElasticIpPool(value); err != nil {
					return err
				} else {
					launchSpec.SetElasticIPPool(elasticIpPool)
				}
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			launchSpecWrapper := resourceObject.(*commons.LaunchSpecWrapper)
			launchSpec := launchSpecWrapper.GetLaunchSpec()
			var value *aws.ElasticIPPool = nil

			if v, ok := resourceData.GetOk(string(ElasticIpPool)); ok {
				if elasticIpPool, err := expandElasticIpPool(v); err != nil {
					return err
				} else {
					value = elasticIpPool
				}
			}
			launchSpec.SetElasticIPPool(value)
			return nil
		},
		nil,
	)

	fieldsMap[BlockDeviceMappings] = commons.NewGenericField(
		commons.OceanAWSLaunchSpec,
		BlockDeviceMappings,
		&schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{

					string(DeviceName): {
						Type:     schema.TypeString,
						Optional: true,
					},

					string(Ebs): {
						Type:     schema.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &schema.Resource{

							Schema: map[string]*schema.Schema{

								string(DeleteOnTermination): {
									Type:     schema.TypeBool,
									Optional: true,
									Computed: true,
								},

								string(Encrypted): {
									Type:     schema.TypeBool,
									Optional: true,
									Computed: true,
								},

								string(IOPS): {
									Type:     schema.TypeInt,
									Optional: true,
								},

								string(KMSKeyID): {
									Type:     schema.TypeString,
									Optional: true,
								},

								string(SnapshotID): {
									Type:     schema.TypeString,
									Optional: true,
								},

								string(VolumeSize): {
									Type:     schema.TypeInt,
									Optional: true,
								},

								string(VolumeType): {
									Type:     schema.TypeString,
									Optional: true,
									Computed: true,
								},

								string(DynamicVolumeSize): {
									Type:     schema.TypeList,
									Optional: true,
									MaxItems: 1,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{

											string(BaseSize): {
												Type:     schema.TypeInt,
												Required: true,
											},

											string(Resource): {
												Type:     schema.TypeString,
												Required: true,
											},

											string(SizePerResourceUnit): {
												Type:     schema.TypeInt,
												Required: true,
											},
										},
									},
								},

								string(Throughput): {
									Type:     schema.TypeInt,
									Optional: true,
								},
							},
						},
					},

					string(NoDevice): {
						Type:     schema.TypeString,
						Optional: true,
					},

					string(VirtualName): {
						Type:     schema.TypeString,
						Optional: true,
					},
				},
			},
		},

		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			launchSpecWrapper := resourceObject.(*commons.LaunchSpecWrapper)
			launchSpec := launchSpecWrapper.GetLaunchSpec()
			var result []interface{} = nil

			if launchSpec != nil && launchSpec.BlockDeviceMappings != nil {
				result = flattenBlockDeviceMappings(launchSpec.BlockDeviceMappings)
			}

			if len(result) > 0 {
				if err := resourceData.Set(string(BlockDeviceMappings), result); err != nil {
					return fmt.Errorf(string(commons.FailureFieldReadPattern), string(BlockDeviceMappings), err)
				}
			}
			return nil
		},

		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			launchSpecWrapper := resourceObject.(*commons.LaunchSpecWrapper)
			launchSpec := launchSpecWrapper.GetLaunchSpec()
			if v, ok := resourceData.GetOk(string(BlockDeviceMappings)); ok {
				if v, err := expandBlockDeviceMappings(v); err != nil {
					return err
				} else {
					launchSpec.SetBlockDeviceMappings(v)
				}
			}
			return nil
		},

		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			launchSpecWrapper := resourceObject.(*commons.LaunchSpecWrapper)
			launchSpec := launchSpecWrapper.GetLaunchSpec()
			var value []*aws.BlockDeviceMapping = nil

			if v, ok := resourceData.GetOk(string(BlockDeviceMappings)); ok {
				if blockdevicemappings, err := expandBlockDeviceMappings(v); err != nil {
					return err
				} else {
					value = blockdevicemappings
				}
			}
			launchSpec.SetBlockDeviceMappings(value)
			return nil
		},
		nil,
	)

	fieldsMap[ResourceLimits] = commons.NewGenericField(
		commons.OceanAWSLaunchConfiguration,
		ResourceLimits,
		&schema.Schema{
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{

					string(MaxInstanceCount): {
						Type:     schema.TypeInt,
						Optional: true,
					},
					string(MinInstanceCount): {
						Type:     schema.TypeInt,
						Optional: true,
					},
				},
			},
		},

		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			launchSpecWrapper := resourceObject.(*commons.LaunchSpecWrapper)
			launchSpec := launchSpecWrapper.GetLaunchSpec()
			var result []interface{} = nil
			if launchSpec.ResourceLimits != nil {
				resourceLimits := launchSpec.ResourceLimits
				result = flattenResourceLimits(resourceLimits)
			}
			if result != nil {
				if err := resourceData.Set(string(ResourceLimits), result); err != nil {
					return fmt.Errorf(string(commons.FailureFieldReadPattern), string(ResourceLimits), err)
				}
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			launchSpecWrapper := resourceObject.(*commons.LaunchSpecWrapper)
			launchSpec := launchSpecWrapper.GetLaunchSpec()
			if value, ok := resourceData.GetOk(string(ResourceLimits)); ok {
				if resourceLimits, err := expandResourceLimits(value); err != nil {
					return err
				} else {
					launchSpec.SetResourceLimits(resourceLimits)
				}
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			launchSpecWrapper := resourceObject.(*commons.LaunchSpecWrapper)
			launchSpec := launchSpecWrapper.GetLaunchSpec()
			var value *aws.ResourceLimits = nil

			if v, ok := resourceData.GetOk(string(ResourceLimits)); ok {
				if resourceLimits, err := expandResourceLimits(v); err != nil {
					return err
				} else {
					value = resourceLimits
				}
			}
			launchSpec.SetResourceLimits(value)
			return nil
		},
		nil,
	)

	fieldsMap[SecurityGroups] = commons.NewGenericField(
		commons.OceanAWSLaunchSpec,
		SecurityGroups,
		&schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			launchSpecWrapper := resourceObject.(*commons.LaunchSpecWrapper)
			launchSpec := launchSpecWrapper.GetLaunchSpec()
			var value []string = nil
			if launchSpec.SecurityGroupIDs != nil {
				value = launchSpec.SecurityGroupIDs
			}
			if err := resourceData.Set(string(SecurityGroups), value); err != nil {
				return fmt.Errorf(string(commons.FailureFieldReadPattern), string(SecurityGroups), err)
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			launchSpecWrapper := resourceObject.(*commons.LaunchSpecWrapper)
			launchSpec := launchSpecWrapper.GetLaunchSpec()
			if v, ok := resourceData.Get(string(SecurityGroups)).([]interface{}); ok {
				ids := make([]string, len(v))
				for i, j := range v {
					ids[i] = j.(string)
				}
				launchSpec.SetSecurityGroupIDs(ids)
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			launchSpecWrapper := resourceObject.(*commons.LaunchSpecWrapper)
			launchSpec := launchSpecWrapper.GetLaunchSpec()
			if v, ok := resourceData.Get(string(SecurityGroups)).([]interface{}); ok {
				ids := make([]string, len(v))
				for i, j := range v {
					ids[i] = j.(string)
				}
				launchSpec.SetSecurityGroupIDs(ids)
			}
			return nil
		},
		nil,
	)

	fieldsMap[Taints] = commons.NewGenericField(
		commons.OceanAWSLaunchSpec,
		Taints,
		&schema.Schema{
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					string(TaintKey): {
						Type:     schema.TypeString,
						Required: true,
					},
					string(TaintValue): {
						Type:     schema.TypeString,
						Required: true,
					},
					string(Effect): {
						Type:     schema.TypeString,
						Required: true,
					},
				},
			},
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			launchSpecWrapper := resourceObject.(*commons.LaunchSpecWrapper)
			launchSpec := launchSpecWrapper.GetLaunchSpec()
			var result []interface{} = nil
			if launchSpec.Labels != nil {
				taints := launchSpec.Taints
				result = flattenTaints(taints)
			}
			if result != nil {
				if err := resourceData.Set(string(Taints), result); err != nil {
					return fmt.Errorf(string(commons.FailureFieldReadPattern), string(Taints), err)
				}
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			launchSpecWrapper := resourceObject.(*commons.LaunchSpecWrapper)
			launchSpec := launchSpecWrapper.GetLaunchSpec()
			if value, ok := resourceData.GetOk(string(Taints)); ok {
				if labels, err := expandTaints(value); err != nil {
					return err
				} else {
					launchSpec.SetTaints(labels)
				}
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			launchSpecWrapper := resourceObject.(*commons.LaunchSpecWrapper)
			launchSpec := launchSpecWrapper.GetLaunchSpec()
			var taintList []*aws.Taint = nil
			if value, ok := resourceData.GetOk(string(Taints)); ok {
				if taints, err := expandTaints(value); err != nil {
					return err
				} else {
					taintList = taints
				}
			}
			launchSpec.SetTaints(taintList)
			return nil
		},
		nil,
	)

	fieldsMap[IamInstanceProfile] = commons.NewGenericField(
		commons.OceanAWSLaunchSpec,
		IamInstanceProfile,
		&schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			launchSpecWrapper := resourceObject.(*commons.LaunchSpecWrapper)
			launchSpec := launchSpecWrapper.GetLaunchSpec()
			var value = ""
			if launchSpec.IAMInstanceProfile != nil {

				iam := launchSpec.IAMInstanceProfile
				if iam.ARN != nil {
					value = spotinst.StringValue(iam.ARN)
				} else if iam.Name != nil {
					value = spotinst.StringValue(iam.Name)
				}
			}
			if err := resourceData.Set(string(IamInstanceProfile), value); err != nil {
				return fmt.Errorf(string(commons.FailureFieldReadPattern), string(IamInstanceProfile), err)
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			launchSpecWrapper := resourceObject.(*commons.LaunchSpecWrapper)
			launchSpec := launchSpecWrapper.GetLaunchSpec()
			if v, ok := resourceData.Get(string(IamInstanceProfile)).(string); ok && v != "" {
				iamInstanceProf := &aws.IAMInstanceProfile{}
				if InstanceProfileArnRegex.MatchString(v) {
					iamInstanceProf.SetArn(spotinst.String(v))
				} else {
					iamInstanceProf.SetName(spotinst.String(v))
				}
				launchSpec.SetIAMInstanceProfile(iamInstanceProf)
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			launchSpecWrapper := resourceObject.(*commons.LaunchSpecWrapper)
			launchSpec := launchSpecWrapper.GetLaunchSpec()
			if v, ok := resourceData.Get(string(IamInstanceProfile)).(string); ok && v != "" {
				iamInstanceProf := &aws.IAMInstanceProfile{}
				if InstanceProfileArnRegex.MatchString(v) {
					iamInstanceProf.SetArn(spotinst.String(v))
				} else {
					iamInstanceProf.SetName(spotinst.String(v))
				}
				launchSpec.SetIAMInstanceProfile(iamInstanceProf)
			} else {
				launchSpec.SetIAMInstanceProfile(nil)
			}
			return nil
		},
		nil,
	)

	fieldsMap[AutoscaleHeadrooms] = commons.NewGenericField(
		commons.OceanAWSLaunchSpec,
		AutoscaleHeadrooms,
		&schema.Schema{
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					string(CPUPerUnit): {
						Type:     schema.TypeInt,
						Optional: true,
					},

					string(GPUPerUnit): {
						Type:     schema.TypeInt,
						Optional: true,
					},

					string(MemoryPerUnit): {
						Type:     schema.TypeInt,
						Optional: true,
					},

					string(NumOfUnits): {
						Type:     schema.TypeInt,
						Required: true,
					},
				},
			},
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			launchSpecWrapper := resourceObject.(*commons.LaunchSpecWrapper)
			launchSpec := launchSpecWrapper.GetLaunchSpec()
			var result []interface{} = nil
			if launchSpec.AutoScale != nil && launchSpec.AutoScale.Headrooms != nil {
				headrooms := launchSpec.AutoScale.Headrooms
				result = flattenHeadrooms(headrooms)
			}
			if result != nil {
				if err := resourceData.Set(string(AutoscaleHeadrooms), result); err != nil {
					return fmt.Errorf(string(commons.FailureFieldReadPattern), string(AutoscaleHeadrooms), err)
				}
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			launchSpecWrapper := resourceObject.(*commons.LaunchSpecWrapper)
			launchSpec := launchSpecWrapper.GetLaunchSpec()
			if value, ok := resourceData.GetOk(string(AutoscaleHeadrooms)); ok {
				if headrooms, err := expandHeadrooms(value); err != nil {
					return err
				} else {
					launchSpec.AutoScale.SetHeadrooms(headrooms)
				}
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			launchSpecWrapper := resourceObject.(*commons.LaunchSpecWrapper)
			launchSpec := launchSpecWrapper.GetLaunchSpec()
			var headroomList []*aws.AutoScaleHeadroom = nil
			if value, ok := resourceData.GetOk(string(AutoscaleHeadrooms)); ok {
				if expandedList, err := expandHeadrooms(value); err != nil {
					return err
				} else {
					headroomList = expandedList
				}
			}
			launchSpec.AutoScale.SetHeadrooms(headroomList)
			return nil
		},
		nil,
	)

	fieldsMap[SubnetIDs] = commons.NewGenericField(
		commons.OceanAWSLaunchSpec,
		SubnetIDs,
		&schema.Schema{
			Type:     schema.TypeList,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Optional: true,
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			launchSpecWrapper := resourceObject.(*commons.LaunchSpecWrapper)
			launchSpec := launchSpecWrapper.GetLaunchSpec()
			var value []string = nil
			if launchSpec.SubnetIDs != nil {
				value = launchSpec.SubnetIDs
			}
			if err := resourceData.Set(string(SubnetIDs), value); err != nil {
				return fmt.Errorf(string(commons.FailureFieldReadPattern), string(SubnetIDs), err)
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			launchSpecWrapper := resourceObject.(*commons.LaunchSpecWrapper)
			launchSpec := launchSpecWrapper.GetLaunchSpec()
			if v, ok := resourceData.GetOk(string(SubnetIDs)); ok {
				if subnetIDs, err := expandSubnetIDs(v); err != nil {
					return err
				} else {
					launchSpec.SetSubnetIDs(subnetIDs)
				}
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			launchSpecWrapper := resourceObject.(*commons.LaunchSpecWrapper)
			launchSpec := launchSpecWrapper.GetLaunchSpec()
			if v, ok := resourceData.GetOk(string(SubnetIDs)); ok {
				if subnetIDs, err := expandSubnetIDs(v); err != nil {
					return err
				} else {
					launchSpec.SetSubnetIDs(subnetIDs)
				}
			} else {
				launchSpec.SetSubnetIDs(nil)
			}

			return nil
		},
		nil,
	)

	fieldsMap[InstanceTypes] = commons.NewGenericField(
		commons.OceanAWSLaunchSpec,
		InstanceTypes,
		&schema.Schema{
			Type:     schema.TypeList,
			Elem:     &schema.Schema{Type: schema.TypeString},
			MinItems: 1,
			Optional: true,
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			launchSpecWrapper := resourceObject.(*commons.LaunchSpecWrapper)
			launchSpec := launchSpecWrapper.GetLaunchSpec()
			var value []string = nil
			if launchSpec.InstanceTypes != nil {
				value = launchSpec.InstanceTypes
			}
			if err := resourceData.Set(string(InstanceTypes), value); err != nil {
				return fmt.Errorf(string(commons.FailureFieldReadPattern), string(InstanceTypes), err)
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			launchSpecWrapper := resourceObject.(*commons.LaunchSpecWrapper)
			launchSpec := launchSpecWrapper.GetLaunchSpec()
			if v, ok := resourceData.GetOk(string(InstanceTypes)); ok {
				if instanceTypes, err := expandInstanceTypes(v); err != nil {
					return err
				} else {
					launchSpec.SetInstanceTypes(instanceTypes)
				}
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			launchSpecWrapper := resourceObject.(*commons.LaunchSpecWrapper)
			launchSpec := launchSpecWrapper.GetLaunchSpec()
			if v, ok := resourceData.GetOk(string(InstanceTypes)); ok {
				if instanceTypes, err := expandInstanceTypes(v); err != nil {
					return err
				} else {
					launchSpec.SetInstanceTypes(instanceTypes)
				}
			} else {
				launchSpec.SetInstanceTypes(nil)
			}

			return nil
		},
		nil,
	)

	fieldsMap[PreferredSpotTypes] = commons.NewGenericField(
		commons.OceanAWSLaunchSpec,
		PreferredSpotTypes,
		&schema.Schema{
			Type:     schema.TypeList,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Optional: true,
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			launchSpecWrapper := resourceObject.(*commons.LaunchSpecWrapper)
			launchSpec := launchSpecWrapper.GetLaunchSpec()
			var value []string = nil
			if launchSpec.PreferredSpotTypes != nil {
				value = launchSpec.PreferredSpotTypes
			}
			if err := resourceData.Set(string(PreferredSpotTypes), value); err != nil {
				return fmt.Errorf(string(commons.FailureFieldReadPattern), string(PreferredSpotTypes), err)
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			launchSpecWrapper := resourceObject.(*commons.LaunchSpecWrapper)
			launchSpec := launchSpecWrapper.GetLaunchSpec()
			if v, ok := resourceData.GetOk(string(PreferredSpotTypes)); ok {
				if preferredSpotTypes, err := expandPreferredSpotTypes(v); err != nil {
					return err
				} else {
					launchSpec.SetPreferredSpotTypes(preferredSpotTypes)
				}
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			launchSpecWrapper := resourceObject.(*commons.LaunchSpecWrapper)
			launchSpec := launchSpecWrapper.GetLaunchSpec()
			if v, ok := resourceData.GetOk(string(PreferredSpotTypes)); ok {
				if preferredSpotTypes, err := expandInstanceTypes(v); err != nil {
					return err
				} else {
					launchSpec.SetPreferredSpotTypes(preferredSpotTypes)
				}
			} else {
				launchSpec.SetPreferredSpotTypes(nil)
			}

			return nil
		},
		nil,
	)

	fieldsMap[RootVolumeSize] = commons.NewGenericField(
		commons.OceanAWSLaunchSpec,
		RootVolumeSize,
		&schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			launchSpecWrapper := resourceObject.(*commons.LaunchSpecWrapper)
			launchSpec := launchSpecWrapper.GetLaunchSpec()
			var value *int = nil
			if launchSpec.RootVolumeSize != nil {
				value = launchSpec.RootVolumeSize
			}
			if err := resourceData.Set(string(RootVolumeSize), spotinst.IntValue(value)); err != nil {
				return fmt.Errorf(string(commons.FailureFieldReadPattern), string(RootVolumeSize), err)
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			launchSpecWrapper := resourceObject.(*commons.LaunchSpecWrapper)
			launchSpec := launchSpecWrapper.GetLaunchSpec()
			if v, ok := resourceData.Get(string(RootVolumeSize)).(int); ok && v > 0 {
				launchSpec.SetRootVolumeSize(spotinst.Int(v))
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			launchSpecWrapper := resourceObject.(*commons.LaunchSpecWrapper)
			launchSpec := launchSpecWrapper.GetLaunchSpec()
			var value *int = nil
			if v, ok := resourceData.Get(string(RootVolumeSize)).(int); ok && v > 0 {
				value = spotinst.Int(v)
			}
			launchSpec.SetRootVolumeSize(value)
			return nil
		},
		nil,
	)

	fieldsMap[Strategy] = commons.NewGenericField(
		commons.OceanAWSLaunchSpec,
		Strategy,
		&schema.Schema{
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					string(SpotPercentage): {
						Type:         schema.TypeInt,
						Optional:     true,
						Default:      -1,
						ValidateFunc: validation.IntAtLeast(-1),
					},
				},
			},
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			launchSpecWrapper := resourceObject.(*commons.LaunchSpecWrapper)
			launchSpec := launchSpecWrapper.GetLaunchSpec()
			var result []interface{} = nil
			if launchSpec.Strategy != nil {
				strategy := launchSpec.Strategy
				result = flattenStrategy(strategy)
			}
			if result != nil {
				if err := resourceData.Set(string(Strategy), result); err != nil {
					return fmt.Errorf(string(commons.FailureFieldReadPattern), string(Strategy), err)
				}
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			launchSpecWrapper := resourceObject.(*commons.LaunchSpecWrapper)
			launchSpec := launchSpecWrapper.GetLaunchSpec()
			if value, ok := resourceData.GetOk(string(Strategy)); ok {
				if strategy, err := expandStrategy(value); err != nil {
					return err
				} else {
					launchSpec.SetStrategy(strategy)
				}
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			launchSpecWrapper := resourceObject.(*commons.LaunchSpecWrapper)
			launchSpec := launchSpecWrapper.GetLaunchSpec()
			var value *aws.LaunchSpecStrategy = nil

			if v, ok := resourceData.GetOk(string(Strategy)); ok {
				if strategy, err := expandStrategy(v); err != nil {
					return err
				} else {
					value = strategy
				}
			}
			launchSpec.SetStrategy(value)
			return nil
		},
		nil,
	)

	fieldsMap[AssociatePublicIPAddress] = commons.NewGenericField(
		commons.OceanAWSLaunchSpec,
		AssociatePublicIPAddress,
		&schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			launchSpecWrapper := resourceObject.(*commons.LaunchSpecWrapper)
			launchSpec := launchSpecWrapper.GetLaunchSpec()
			var value *bool = nil
			if launchSpec.AssociatePublicIPAddress != nil {
				value = launchSpec.AssociatePublicIPAddress
			}
			if value != nil {
				if err := resourceData.Set(string(AssociatePublicIPAddress), spotinst.BoolValue(value)); err != nil {
					return fmt.Errorf(string(commons.FailureFieldReadPattern), string(AssociatePublicIPAddress), err)
				}
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			launchSpecWrapper := resourceObject.(*commons.LaunchSpecWrapper)
			launchSpec := launchSpecWrapper.GetLaunchSpec()
			if v, ok := resourceData.GetOkExists(string(AssociatePublicIPAddress)); ok && v != nil {
				associatePublicIPAddress := spotinst.Bool(v.(bool))
				launchSpec.SetAssociatePublicIPAddress(associatePublicIPAddress)
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			launchSpecWrapper := resourceObject.(*commons.LaunchSpecWrapper)
			launchSpec := launchSpecWrapper.GetLaunchSpec()
			var associatePublicIPAddress *bool = nil
			if v, ok := resourceData.GetOkExists(string(AssociatePublicIPAddress)); ok && v != nil {
				associatePublicIPAddress = spotinst.Bool(v.(bool))
			}
			launchSpec.SetAssociatePublicIPAddress(associatePublicIPAddress)
			return nil
		},
		nil,
	)

	fieldsMap[RestrictScaleDown] = commons.NewGenericField(
		commons.OceanAWSLaunchSpec,
		RestrictScaleDown,
		&schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			launchSpecWrapper := resourceObject.(*commons.LaunchSpecWrapper)
			launchSpec := launchSpecWrapper.GetLaunchSpec()
			var value *bool = nil
			if launchSpec.RestrictScaleDown != nil {
				value = launchSpec.RestrictScaleDown
			}
			if value != nil {
				if err := resourceData.Set(string(RestrictScaleDown), spotinst.BoolValue(value)); err != nil {
					return fmt.Errorf(string(commons.FailureFieldReadPattern), string(RestrictScaleDown), err)
				}
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			launchSpecWrapper := resourceObject.(*commons.LaunchSpecWrapper)
			launchSpec := launchSpecWrapper.GetLaunchSpec()
			if v, ok := resourceData.GetOkExists(string(RestrictScaleDown)); ok && v != nil {
				restrictScaleDown := spotinst.Bool(v.(bool))
				launchSpec.SetRestrictScaleDown(restrictScaleDown)
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			launchSpecWrapper := resourceObject.(*commons.LaunchSpecWrapper)
			launchSpec := launchSpecWrapper.GetLaunchSpec()
			var restrictScaleDown *bool = nil
			if v, ok := resourceData.GetOkExists(string(RestrictScaleDown)); ok && v != nil {
				restrictScaleDown = spotinst.Bool(v.(bool))
			}
			launchSpec.SetRestrictScaleDown(restrictScaleDown)
			return nil
		},
		nil,
	)

	fieldsMap[CreateOptions] = commons.NewGenericField(
		commons.OceanAWSLaunchSpec,
		CreateOptions,
		&schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					string(InitialNodes): {
						Type:     schema.TypeInt,
						Required: true,
					},
				},
			},
		},
		nil, nil, nil, nil,
	)

	fieldsMap[UpdatePolicy] = commons.NewGenericField(
		commons.OceanAWSLaunchSpec,
		UpdatePolicy,
		&schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					string(ShouldRoll): {
						Type:     schema.TypeBool,
						Required: true,
					},

					string(RollConfig): {
						Type:     schema.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								string(BatchSizePercentage): {
									Type:     schema.TypeInt,
									Required: true,
								},
							},
						},
					},
				},
			},
		},
		nil, nil, nil, nil,
	)

	fieldsMap[DeleteOptions] = commons.NewGenericField(
		commons.OceanAWSLaunchSpec,
		DeleteOptions,
		&schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					string(ForceDelete): {
						Type:     schema.TypeBool,
						Required: true,
					},
				},
			},
		},
		nil, nil, nil, nil,
	)

	fieldsMap[SchedulingTask] = commons.NewGenericField(
		commons.OceanAWSLaunchSpec,
		SchedulingTask,
		&schema.Schema{
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{

					string(IsEnabled): {
						Type:     schema.TypeBool,
						Required: true,
					},
					string(CronExpression): {
						Type:     schema.TypeString,
						Required: true,
					},
					string(TaskType): {
						Type:     schema.TypeString,
						Required: true,
					},
					string(TaskHeadroom): {
						Type:     schema.TypeSet,
						Optional: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{

								string(CPUPerUnit): {
									Type:     schema.TypeInt,
									Optional: true,
								},

								string(GPUPerUnit): {
									Type:     schema.TypeInt,
									Optional: true,
								},

								string(MemoryPerUnit): {
									Type:     schema.TypeInt,
									Optional: true,
								},

								string(NumOfUnits): {
									Type:     schema.TypeInt,
									Required: true,
								},
							},
						},
					},
				},
			},
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			launchSpecWrapper := resourceObject.(*commons.LaunchSpecWrapper)
			launchSpec := launchSpecWrapper.GetLaunchSpec()
			var result []interface{} = nil
			if launchSpec.LaunchSpecScheduling != nil && launchSpec.LaunchSpecScheduling.Tasks != nil {
				tasks := launchSpec.LaunchSpecScheduling.Tasks
				result = flattenTasks(tasks)
			}
			if result != nil {
				if err := resourceData.Set(string(SchedulingTask), result); err != nil {
					return fmt.Errorf(string(commons.FailureFieldReadPattern), string(SchedulingTask), err)
				}
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			launchSpecWrapper := resourceObject.(*commons.LaunchSpecWrapper)
			launchSpec := launchSpecWrapper.GetLaunchSpec()
			if value, ok := resourceData.GetOk(string(SchedulingTask)); ok {
				if tasks, err := expandTasks(value); err != nil {
					return err
				} else {
					launchSpec.LaunchSpecScheduling.SetTasks(tasks)
				}
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			launchSpecWrapper := resourceObject.(*commons.LaunchSpecWrapper)
			launchSpec := launchSpecWrapper.GetLaunchSpec()
			var value []*aws.LaunchSpecTask = nil

			if v, ok := resourceData.GetOk(string(SchedulingTask)); ok {
				if tasks, err := expandTasks(v); err != nil {
					return err
				} else {
					value = tasks
				}
			}
			launchSpec.LaunchSpecScheduling.SetTasks(value)
			return nil
		},
		nil,
	)
}

func hashKV(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	buf.WriteString(fmt.Sprintf("%s-", m[string(LabelKey)].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m[string(LabelValue)].(string)))
	return hashcode.String(buf.String())
}

var InstanceProfileArnRegex = regexp.MustCompile(`arn:aws:iam::\d{12}:instance-profile/?[a-zA-Z_0-9+=,.@\-_/]+`)

func Base64StateFunc(v interface{}) string {
	if isBase64Encoded(v.(string)) {
		return v.(string)
	} else {
		return base64Encode(v.(string))
	}
}

// base64Encode encodes data if the input isn't already encoded using
// base64.StdEncoding.EncodeToString. If the input is already base64 encoded,
// return the original input unchanged.
func base64Encode(data string) string {
	// Check whether the data is already Base64 encoded; don't double-encode
	if isBase64Encoded(data) {
		return data
	}
	// data has not been encoded encode and return
	return base64.StdEncoding.EncodeToString([]byte(data))
}

func isBase64Encoded(data string) bool {
	_, err := base64.StdEncoding.DecodeString(data)
	return err == nil
}

func expandLabels(data interface{}) ([]*aws.Label, error) {
	list := data.(*schema.Set).List()
	labels := make([]*aws.Label, 0, len(list))
	for _, v := range list {
		attr, ok := v.(map[string]interface{})
		if !ok {
			continue
		}
		if _, ok := attr[string(LabelKey)]; !ok {
			return nil, errors.New("invalid label attributes: key missing")
		}

		if _, ok := attr[string(LabelValue)]; !ok {
			return nil, errors.New("invalid label attributes: value missing")
		}
		label := &aws.Label{
			Key:   spotinst.String(attr[string(LabelKey)].(string)),
			Value: spotinst.String(attr[string(LabelValue)].(string)),
		}
		labels = append(labels, label)
	}
	return labels, nil
}

func flattenLabels(labels []*aws.Label) []interface{} {
	result := make([]interface{}, 0, len(labels))
	for _, label := range labels {
		m := make(map[string]interface{})
		m[string(LabelKey)] = spotinst.StringValue(label.Key)
		m[string(LabelValue)] = spotinst.StringValue(label.Value)

		result = append(result, m)
	}
	return result
}

func expandTaints(data interface{}) ([]*aws.Taint, error) {
	list := data.(*schema.Set).List()
	taints := make([]*aws.Taint, 0, len(list))
	for _, v := range list {
		attr, ok := v.(map[string]interface{})
		if !ok {
			continue
		}
		if _, ok := attr[string(TaintKey)]; !ok {
			return nil, errors.New("invalid taint attributes: key missing")
		}

		if _, ok := attr[string(TaintValue)]; !ok {
			return nil, errors.New("invalid taint attributes: value missing")
		}

		if _, ok := attr[string(Effect)]; !ok {
			return nil, errors.New("invalid taint attributes: effect missing")
		}

		taint := &aws.Taint{
			Key:    spotinst.String(attr[string(TaintKey)].(string)),
			Value:  spotinst.String(attr[string(TaintValue)].(string)),
			Effect: spotinst.String(attr[string(Effect)].(string)),
		}
		taints = append(taints, taint)
	}
	return taints, nil
}

func flattenTaints(taints []*aws.Taint) []interface{} {
	result := make([]interface{}, 0, len(taints))
	for _, taint := range taints {
		m := make(map[string]interface{})
		m[string(TaintKey)] = spotinst.StringValue(taint.Key)
		m[string(TaintValue)] = spotinst.StringValue(taint.Value)
		m[string(Effect)] = spotinst.StringValue(taint.Effect)

		result = append(result, m)
	}
	return result
}

func expandHeadrooms(data interface{}) ([]*aws.AutoScaleHeadroom, error) {
	list := data.(*schema.Set).List()
	headrooms := make([]*aws.AutoScaleHeadroom, 0, len(list))

	for _, v := range list {
		attr, ok := v.(map[string]interface{})
		if !ok {
			continue
		}

		headroom := &aws.AutoScaleHeadroom{
			CPUPerUnit:    spotinst.Int(attr[string(CPUPerUnit)].(int)),
			GPUPerUnit:    spotinst.Int(attr[string(GPUPerUnit)].(int)),
			NumOfUnits:    spotinst.Int(attr[string(NumOfUnits)].(int)),
			MemoryPerUnit: spotinst.Int(attr[string(MemoryPerUnit)].(int)),
		}

		headrooms = append(headrooms, headroom)
	}
	return headrooms, nil
}

func flattenHeadrooms(headrooms []*aws.AutoScaleHeadroom) []interface{} {
	result := make([]interface{}, 0, len(headrooms))

	for _, headroom := range headrooms {
		m := make(map[string]interface{})
		m[string(CPUPerUnit)] = spotinst.IntValue(headroom.CPUPerUnit)
		m[string(GPUPerUnit)] = spotinst.IntValue(headroom.GPUPerUnit)
		m[string(NumOfUnits)] = spotinst.IntValue(headroom.NumOfUnits)
		m[string(MemoryPerUnit)] = spotinst.IntValue(headroom.MemoryPerUnit)

		result = append(result, m)
	}

	return result
}

func expandSubnetIDs(data interface{}) ([]string, error) {
	list := data.([]interface{})
	result := make([]string, 0, len(list))

	for _, v := range list {
		if subnetID, ok := v.(string); ok && subnetID != "" {
			result = append(result, subnetID)
		}
	}
	return result, nil
}

func expandInstanceTypes(data interface{}) ([]string, error) {
	list := data.([]interface{})
	result := make([]string, 0, len(list))

	for _, v := range list {
		if instanceTypes, ok := v.(string); ok && instanceTypes != "" {
			result = append(result, instanceTypes)
		}
	}
	return result, nil
}

func expandPreferredSpotTypes(data interface{}) ([]string, error) {
	list := data.([]interface{})
	result := make([]string, 0, len(list))

	for _, v := range list {
		if preferredSpotTypes, ok := v.(string); ok && preferredSpotTypes != "" {
			result = append(result, preferredSpotTypes)
		}
	}
	return result, nil
}

func flattenTags(tags []*aws.Tag) []interface{} {
	result := make([]interface{}, 0, len(tags))
	for _, tag := range tags {
		m := make(map[string]interface{})
		m[string(TagKey)] = spotinst.StringValue(tag.Key)
		m[string(TagValue)] = spotinst.StringValue(tag.Value)

		result = append(result, m)
	}
	return result
}

func expandTags(data interface{}) ([]*aws.Tag, error) {
	list := data.(*schema.Set).List()
	tags := make([]*aws.Tag, 0, len(list))
	for _, v := range list {
		attr, ok := v.(map[string]interface{})
		if !ok {
			continue
		}
		if _, ok := attr[string(TagKey)]; !ok {
			return nil, errors.New("invalid tag attributes: key missing")
		}

		if _, ok := attr[string(TagValue)]; !ok {
			return nil, errors.New("invalid tag attributes: value missing")
		}
		tag := &aws.Tag{
			Key:   spotinst.String(attr[string(TagKey)].(string)),
			Value: spotinst.String(attr[string(TagValue)].(string)),
		}
		tags = append(tags, tag)
	}
	return tags, nil
}

func expandElasticIpPool(data interface{}) (*aws.ElasticIPPool, error) {
	elasticIpPool := &aws.ElasticIPPool{}
	list := data.(*schema.Set).List()

	if list == nil || list[0] == nil {
		return elasticIpPool, nil
	}
	m := list[0].(map[string]interface{})

	if v, ok := m[string(TagSelector)]; ok {
		tagSelector, err := expandTagSelector(v)
		if err != nil {
			return nil, err
		}
		if tagSelector != nil {
			elasticIpPool.SetTagSelector(tagSelector)
		} else {
			elasticIpPool.SetTagSelector(nil)
		}
	}
	return elasticIpPool, nil

}

func expandTagSelector(data interface{}) (*aws.TagSelector, error) {
	if list := data.([]interface{}); len(list) > 0 {
		tagSelector := &aws.TagSelector{}
		if list != nil && list[0] != nil {
			m := list[0].(map[string]interface{})

			if v, ok := m[string(TagSelectorKey)].(string); ok && v != "" {
				tagSelector.SetTagKey(spotinst.String(v))
			}

			if v, ok := m[string(TagSelectorValue)].(string); ok && v != "" {
				tagSelector.SetTagValue(spotinst.String(v))
			}
		}
		return tagSelector, nil
	}

	return nil, nil
}

func flattenElasticIpPool(elasticIpPool *aws.ElasticIPPool) []interface{} {
	var out []interface{}

	if elasticIpPool != nil {
		result := make(map[string]interface{})

		if elasticIpPool.TagSelector != nil {
			result[string(TagSelector)] = flattenTagSelector(elasticIpPool.TagSelector)
		}
		if len(result) > 0 {
			out = append(out, result)
		}
	}

	return out
}

func expandBlockDeviceMappings(data interface{}) ([]*aws.BlockDeviceMapping, error) {

	list := data.([]interface{})
	bdms := make([]*aws.BlockDeviceMapping, 0, len(list))

	for _, v := range list {
		attr, ok := v.(map[string]interface{})
		bdm := &aws.BlockDeviceMapping{}

		if !ok {
			continue
		}

		if v, ok := attr[string(DeviceName)].(string); ok && v != "" {
			bdm.SetDeviceName(spotinst.String(v))
		}

		if r, ok := attr[string(Ebs)]; ok {
			if ebs, err := expandEbs(r); err != nil {
				return nil, err
			} else {
				bdm.SetEBS(ebs)
			}
		}

		if v, ok := attr[string(NoDevice)].(string); ok && v != "" {
			bdm.SetNoDevice(spotinst.String(v))
		}

		if v, ok := attr[string(VirtualName)].(string); ok && v != "" {
			bdm.SetVirtualName(spotinst.String(v))
		}
		bdms = append(bdms, bdm)
	}
	return bdms, nil
}

func expandEbs(data interface{}) (*aws.EBS, error) {

	ebs := &aws.EBS{}
	list := data.([]interface{})

	if list == nil || list[0] == nil {
		return ebs, nil
	}
	m := list[0].(map[string]interface{})

	if v, ok := m[string(DeleteOnTermination)].(bool); ok {
		ebs.SetDeleteOnTermination(spotinst.Bool(v))
	}

	if v, ok := m[string(Encrypted)].(bool); ok {
		ebs.SetEncrypted(spotinst.Bool(v))
	}

	if v, ok := m[string(IOPS)].(int); ok && v > 0 {
		ebs.SetIOPS(spotinst.Int(v))
	}

	if v, ok := m[string(KMSKeyID)].(string); ok && v != "" {
		ebs.SetKMSKeyId(spotinst.String(v))
	}

	if v, ok := m[string(SnapshotID)].(string); ok && v != "" {
		ebs.SetSnapshotId(spotinst.String(v))
	}

	if v, ok := m[string(VolumeSize)].(int); ok && v > 0 {
		ebs.SetVolumeSize(spotinst.Int(v))
	}

	if v, ok := m[string(VolumeType)].(string); ok && v != "" {
		ebs.SetVolumeType(spotinst.String(v))
	}

	if v, ok := m[string(DynamicVolumeSize)]; ok && v != nil {
		if dynamicVolumeSize, err := expandDynamicVolumeSize(v); err != nil {
			return nil, err
		} else {
			if dynamicVolumeSize != nil {
				ebs.SetDynamicVolumeSize(dynamicVolumeSize)
			}
		}
	}

	if v, ok := m[string(Throughput)].(int); ok && v > 0 {
		ebs.SetThroughput(spotinst.Int(v))
	}
	return ebs, nil
}

func expandDynamicVolumeSize(data interface{}) (*aws.DynamicVolumeSize, error) {
	if list := data.([]interface{}); len(list) > 0 {
		dvs := &aws.DynamicVolumeSize{}
		if list[0] != nil {
			m := list[0].(map[string]interface{})

			if v, ok := m[string(BaseSize)].(int); ok && v >= 0 {
				dvs.SetBaseSize(spotinst.Int(v))
			}

			if v, ok := m[string(Resource)].(string); ok && v != "" {
				dvs.SetResource(spotinst.String(v))
			}

			if v, ok := m[string(SizePerResourceUnit)].(int); ok && v >= 0 {
				dvs.SetSizePerResourceUnit(spotinst.Int(v))
			}
		}
		return dvs, nil
	}
	return nil, nil
}

func flattenBlockDeviceMappings(bdms []*aws.BlockDeviceMapping) []interface{} {
	result := make([]interface{}, 0, len(bdms))

	for _, bdm := range bdms {
		m := make(map[string]interface{})
		m[string(DeviceName)] = spotinst.StringValue(bdm.DeviceName)
		if bdm.EBS != nil {
			m[string(Ebs)] = flattenEbs(bdm.EBS)
		}
		m[string(NoDevice)] = spotinst.StringValue(bdm.NoDevice)
		m[string(VirtualName)] = spotinst.StringValue(bdm.VirtualName)
		result = append(result, m)
	}
	return result

}

func flattenEbs(ebs *aws.EBS) []interface{} {

	elasticBS := make(map[string]interface{})
	elasticBS[string(DeleteOnTermination)] = spotinst.BoolValue(ebs.DeleteOnTermination)
	elasticBS[string(Encrypted)] = spotinst.BoolValue(ebs.Encrypted)
	elasticBS[string(IOPS)] = spotinst.IntValue(ebs.IOPS)
	elasticBS[string(KMSKeyID)] = spotinst.StringValue(ebs.KMSKeyID)
	elasticBS[string(SnapshotID)] = spotinst.StringValue(ebs.SnapshotID)
	elasticBS[string(VolumeType)] = spotinst.StringValue(ebs.VolumeType)
	elasticBS[string(VolumeSize)] = spotinst.IntValue(ebs.VolumeSize)
	elasticBS[string(Throughput)] = spotinst.IntValue(ebs.Throughput)
	if ebs.DynamicVolumeSize != nil {
		elasticBS[string(DynamicVolumeSize)] = flattenDynamicVolumeSize(ebs.DynamicVolumeSize)
	}

	return []interface{}{elasticBS}
}

func flattenDynamicVolumeSize(dvs *aws.DynamicVolumeSize) interface{} {

	DynamicVS := make(map[string]interface{})
	DynamicVS[string(BaseSize)] = spotinst.IntValue(dvs.BaseSize)
	DynamicVS[string(Resource)] = spotinst.StringValue(dvs.Resource)
	DynamicVS[string(SizePerResourceUnit)] = spotinst.IntValue(dvs.SizePerResourceUnit)

	return []interface{}{DynamicVS}
}

func flattenResourceLimits(resourceLimits *aws.ResourceLimits) []interface{} {
	var out []interface{}

	if resourceLimits != nil {
		result := make(map[string]interface{})

		if resourceLimits.MaxInstanceCount != nil {
			result[string(MaxInstanceCount)] = spotinst.IntValue(resourceLimits.MaxInstanceCount)
		}
		if resourceLimits.MinInstanceCount != nil {
			result[string(MinInstanceCount)] = spotinst.IntValue(resourceLimits.MinInstanceCount)
		}
		if len(result) > 0 {
			out = append(out, result)
		}
	}

	return out
}

func flattenTagSelector(tagSelector *aws.TagSelector) []interface{} {
	m := make(map[string]interface{})
	m[string(TagSelectorKey)] = spotinst.StringValue(tagSelector.Key)
	m[string(TagSelectorValue)] = spotinst.StringValue(tagSelector.Value)

	return []interface{}{m}
}

func expandResourceLimits(data interface{}) (*aws.ResourceLimits, error) {
	if list := data.(*schema.Set).List(); len(list) > 0 {
		resLimits := &aws.ResourceLimits{}
		if list != nil && list[0] != nil {
			m := list[0].(map[string]interface{})

			if v, ok := m[string(MaxInstanceCount)].(int); ok && v > 0 {
				resLimits.SetMaxInstanceCount(spotinst.Int(v))
			} else {
				resLimits.SetMaxInstanceCount(nil)
			}

			if v, ok := m[string(MinInstanceCount)].(int); ok && v >= 0 {
				resLimits.SetMinInstanceCount(spotinst.Int(v))
			} else {
				resLimits.SetMinInstanceCount(nil)
			}
		}
		return resLimits, nil
	}

	return nil, nil
}

func expandStrategy(data interface{}) (*aws.LaunchSpecStrategy, error) {
	if list := data.(*schema.Set).List(); len(list) > 0 {
		strategy := &aws.LaunchSpecStrategy{}
		if list != nil && list[0] != nil {
			m := list[0].(map[string]interface{})

			if v, ok := m[string(SpotPercentage)].(int); ok && v > -1 {
				strategy.SetSpotPercentage(spotinst.Int(v))
			} else {
				strategy.SetSpotPercentage(nil)
			}
		}
		return strategy, nil
	}
	return nil, nil
}

func flattenStrategy(strategy *aws.LaunchSpecStrategy) []interface{} {
	var out []interface{}

	if strategy != nil {
		result := make(map[string]interface{})

		if strategy.SpotPercentage != nil {
			result[string(SpotPercentage)] = spotinst.IntValue(strategy.SpotPercentage)
		}
		if len(result) > 0 {
			out = append(out, result)
		}
	}
	return out
}

func flattenTasks(tasks []*aws.LaunchSpecTask) []interface{} {
	result := make([]interface{}, 0, len(tasks))

	for _, task := range tasks {
		m := make(map[string]interface{})
		m[string(IsEnabled)] = spotinst.BoolValue(task.IsEnabled)
		m[string(CronExpression)] = spotinst.StringValue(task.CronExpression)
		m[string(TaskType)] = spotinst.StringValue(task.TaskType)

		if task.Config != nil && task.Config.TaskHeadrooms != nil {
			m[string(TaskHeadroom)] = flattenTaskHeadroom(task.Config.TaskHeadrooms)
		}

		result = append(result, m)
	}

	return result
}

func flattenTaskHeadroom(headrooms []*aws.LaunchSpecTaskHeadroom) []interface{} {
	result := make([]interface{}, 0, len(headrooms))

	for _, headroom := range headrooms {
		m := make(map[string]interface{})
		m[string(CPUPerUnit)] = spotinst.IntValue(headroom.CPUPerUnit)
		m[string(GPUPerUnit)] = spotinst.IntValue(headroom.GPUPerUnit)
		m[string(NumOfUnits)] = spotinst.IntValue(headroom.NumOfUnits)
		m[string(MemoryPerUnit)] = spotinst.IntValue(headroom.MemoryPerUnit)

		result = append(result, m)
	}

	return result
}

func expandTasks(data interface{}) ([]*aws.LaunchSpecTask, error) {
	list := data.(*schema.Set).List()
	tasks := make([]*aws.LaunchSpecTask, 0, len(list))

	for _, v := range list {
		attr, ok := v.(map[string]interface{})
		task := &aws.LaunchSpecTask{}

		if !ok {
			continue
		}

		if v, ok := attr[string(IsEnabled)].(bool); ok {
			task.SetIsEnabled(spotinst.Bool(v))
		}

		if v, ok := attr[string(CronExpression)].(string); ok && v != "" {
			task.SetCronExpression(spotinst.String(v))
		}

		if v, ok := attr[string(TaskType)].(string); ok && v != "" {
			task.SetTaskType(spotinst.String(v))
		}

		if v, ok := attr[string(TaskHeadroom)]; ok {
			if config, err := expandTaskHeadroom(v); err != nil {
				return nil, err
			} else {
				task.SetTaskConfig(config)
			}
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

func expandTaskHeadroom(data interface{}) (*aws.TaskConfig, error) {
	list := data.(*schema.Set).List()
	headrooms := make([]*aws.LaunchSpecTaskHeadroom, 0, len(list))

	for _, v := range list {
		attr, ok := v.(map[string]interface{})
		headroom := &aws.LaunchSpecTaskHeadroom{}

		if !ok {
			continue
		}

		if v, ok := attr[string(CPUPerUnit)].(int); ok {
			headroom.SetCPUPerUnit(spotinst.Int(v))
		}

		if v, ok := attr[string(GPUPerUnit)].(int); ok {
			headroom.SetGPUPerUnit(spotinst.Int(v))
		}

		if v, ok := attr[string(NumOfUnits)].(int); ok {
			headroom.SetNumOfUnits(spotinst.Int(v))
		}

		if v, ok := attr[string(MemoryPerUnit)].(int); ok {
			headroom.SetMemoryPerUnit(spotinst.Int(v))
		}

		headrooms = append(headrooms, headroom)
	}

	taskConfig := &aws.TaskConfig{
		TaskHeadrooms: headrooms,
	}

	return taskConfig, nil
}
