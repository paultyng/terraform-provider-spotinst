package elastigroup_aws

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/spotinst/spotinst-sdk-go/service/elastigroup/providers/aws"
	"github.com/spotinst/spotinst-sdk-go/spotinst"
	"github.com/spotinst/terraform-provider-spotinst/spotinst/commons"
)

func Setup(fieldsMap map[commons.FieldName]*commons.GenericField) {

	fieldsMap[Name] = commons.NewGenericField(
		commons.ElastigroupAWS,
		Name,
		&schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			egWrapper := resourceObject.(*commons.ElastigroupWrapper)
			elastigroup := egWrapper.GetElastigroup()
			var value *string = nil
			if elastigroup.Name != nil {
				value = elastigroup.Name
			}
			if err := resourceData.Set(string(Name), value); err != nil {
				return fmt.Errorf(string(commons.FailureFieldReadPattern), string(Name), err)
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			egWrapper := resourceObject.(*commons.ElastigroupWrapper)
			elastigroup := egWrapper.GetElastigroup()
			elastigroup.SetName(spotinst.String(resourceData.Get(string(Name)).(string)))
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			egWrapper := resourceObject.(*commons.ElastigroupWrapper)
			elastigroup := egWrapper.GetElastigroup()
			elastigroup.SetName(spotinst.String(resourceData.Get(string(Name)).(string)))
			return nil
		},
		nil,
	)

	fieldsMap[Description] = commons.NewGenericField(
		commons.ElastigroupAWS,
		Description,
		&schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			egWrapper := resourceObject.(*commons.ElastigroupWrapper)
			elastigroup := egWrapper.GetElastigroup()
			var value *string = nil
			if elastigroup.Description != nil {
				value = elastigroup.Description
			}
			if err := resourceData.Set(string(Description), value); err != nil {
				return fmt.Errorf(string(commons.FailureFieldReadPattern), string(Description), err)
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			egWrapper := resourceObject.(*commons.ElastigroupWrapper)
			elastigroup := egWrapper.GetElastigroup()
			if v, ok := resourceData.GetOk(string(Description)); ok && v != "" {
				elastigroup.SetDescription(spotinst.String(v.(string)))
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			egWrapper := resourceObject.(*commons.ElastigroupWrapper)
			elastigroup := egWrapper.GetElastigroup()
			if v, ok := resourceData.GetOk(string(Description)); ok && v != "" {
				elastigroup.SetDescription(spotinst.String(v.(string)))
			}
			return nil
		},
		nil,
	)

	fieldsMap[Product] = commons.NewGenericField(
		commons.ElastigroupAWS,
		Product,
		&schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			egWrapper := resourceObject.(*commons.ElastigroupWrapper)
			elastigroup := egWrapper.GetElastigroup()
			var value *string = nil
			if elastigroup.Compute != nil && elastigroup.Compute.Product != nil {
				value = elastigroup.Compute.Product
			}
			if err := resourceData.Set(string(Product), value); err != nil {
				return fmt.Errorf(string(commons.FailureFieldReadPattern), string(Product), err)
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			egWrapper := resourceObject.(*commons.ElastigroupWrapper)
			elastigroup := egWrapper.GetElastigroup()
			elastigroup.Compute.SetProduct(spotinst.String(resourceData.Get(string(Product)).(string)))
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			err := fmt.Errorf(string(commons.FieldUpdateNotAllowedPattern),
				string(Product))
			return err
		},
		nil,
	)

	fieldsMap[MaxSize] = commons.NewGenericField(
		commons.ElastigroupAWS,
		MaxSize,
		&schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
			Computed: true,
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			egWrapper := resourceObject.(*commons.ElastigroupWrapper)
			elastigroup := egWrapper.GetElastigroup()
			var value *int = nil
			if elastigroup.Capacity != nil && elastigroup.Capacity.Maximum != nil {
				value = elastigroup.Capacity.Maximum
			}
			if err := resourceData.Set(string(MaxSize), spotinst.IntValue(value)); err != nil {
				return fmt.Errorf(string(commons.FailureFieldReadPattern), string(MaxSize), err)
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			egWrapper := resourceObject.(*commons.ElastigroupWrapper)
			elastigroup := egWrapper.GetElastigroup()
			if v, ok := resourceData.Get(string(MaxSize)).(int); ok && v >= 0 {
				elastigroup.Capacity.SetMaximum(spotinst.Int(v))
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			egWrapper := resourceObject.(*commons.ElastigroupWrapper)
			elastigroup := egWrapper.GetElastigroup()
			if v, ok := resourceData.Get(string(MaxSize)).(int); ok && v >= 0 {
				elastigroup.Capacity.SetMaximum(spotinst.Int(v))
			}
			return nil
		},
		nil,
	)

	fieldsMap[MinSize] = commons.NewGenericField(
		commons.ElastigroupAWS,
		MinSize,
		&schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
			Computed: true,
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			egWrapper := resourceObject.(*commons.ElastigroupWrapper)
			elastigroup := egWrapper.GetElastigroup()
			var value *int = nil
			if elastigroup.Capacity != nil && elastigroup.Capacity.Minimum != nil {
				value = elastigroup.Capacity.Minimum
			}
			if err := resourceData.Set(string(MinSize), spotinst.IntValue(value)); err != nil {
				return fmt.Errorf(string(commons.FailureFieldReadPattern), string(MinSize), err)
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			egWrapper := resourceObject.(*commons.ElastigroupWrapper)
			elastigroup := egWrapper.GetElastigroup()
			if v, ok := resourceData.Get(string(MinSize)).(int); ok && v >= 0 {
				elastigroup.Capacity.SetMinimum(spotinst.Int(v))
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			egWrapper := resourceObject.(*commons.ElastigroupWrapper)
			elastigroup := egWrapper.GetElastigroup()
			if v, ok := resourceData.Get(string(MinSize)).(int); ok && v >= 0 {
				elastigroup.Capacity.SetMinimum(spotinst.Int(v))
			}
			return nil
		},
		nil,
	)

	fieldsMap[DesiredCapacity] = commons.NewGenericField(
		commons.ElastigroupAWS,
		DesiredCapacity,
		&schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			egWrapper := resourceObject.(*commons.ElastigroupWrapper)
			elastigroup := egWrapper.GetElastigroup()
			var value *int = nil
			if elastigroup.Capacity != nil && elastigroup.Capacity.Target != nil {
				value = elastigroup.Capacity.Target
			}
			if err := resourceData.Set(string(DesiredCapacity), spotinst.IntValue(value)); err != nil {
				return fmt.Errorf(string(commons.FailureFieldReadPattern), string(DesiredCapacity), err)
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			egWrapper := resourceObject.(*commons.ElastigroupWrapper)
			elastigroup := egWrapper.GetElastigroup()
			if v, ok := resourceData.Get(string(DesiredCapacity)).(int); ok && v >= 0 {
				elastigroup.Capacity.SetTarget(spotinst.Int(v))
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			egWrapper := resourceObject.(*commons.ElastigroupWrapper)
			elastigroup := egWrapper.GetElastigroup()
			if v, ok := resourceData.Get(string(DesiredCapacity)).(int); ok && v >= 0 {
				elastigroup.Capacity.SetTarget(spotinst.Int(v))
			}
			return nil
		},
		nil,
	)

	fieldsMap[CapacityUnit] = commons.NewGenericField(
		commons.ElastigroupAWS,
		CapacityUnit,
		&schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			egWrapper := resourceObject.(*commons.ElastigroupWrapper)
			elastigroup := egWrapper.GetElastigroup()
			var value *string = nil
			if elastigroup.Capacity != nil && elastigroup.Capacity.Unit != nil {
				value = elastigroup.Capacity.Unit
			}
			if err := resourceData.Set(string(CapacityUnit), spotinst.StringValue(value)); err != nil {
				return fmt.Errorf(string(commons.FailureFieldReadPattern), string(CapacityUnit), err)
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			egWrapper := resourceObject.(*commons.ElastigroupWrapper)
			elastigroup := egWrapper.GetElastigroup()
			if v, ok := resourceData.Get(string(CapacityUnit)).(string); ok && v != "" {
				elastigroup.Capacity.SetUnit(spotinst.String(v))
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			err := fmt.Errorf(string(commons.FieldUpdateNotAllowedPattern),
				string(CapacityUnit))
			return err
		},
		nil,
	)

	fieldsMap[HealthCheckGracePeriod] = commons.NewGenericField(
		commons.ElastigroupAWS,
		HealthCheckGracePeriod,
		&schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			egWrapper := resourceObject.(*commons.ElastigroupWrapper)
			elastigroup := egWrapper.GetElastigroup()
			var value *int = nil
			if elastigroup.Compute != nil && elastigroup.Compute.LaunchSpecification != nil &&
				elastigroup.Compute.LaunchSpecification.HealthCheckGracePeriod != nil {
				value = elastigroup.Compute.LaunchSpecification.HealthCheckGracePeriod
			}
			if err := resourceData.Set(string(HealthCheckGracePeriod), spotinst.IntValue(value)); err != nil {
				return fmt.Errorf(string(commons.FailureFieldReadPattern), string(HealthCheckGracePeriod), err)
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			egWrapper := resourceObject.(*commons.ElastigroupWrapper)
			elastigroup := egWrapper.GetElastigroup()
			if v, ok := resourceData.Get(string(HealthCheckGracePeriod)).(int); ok && v > 0 {
				elastigroup.Compute.LaunchSpecification.SetHealthCheckGracePeriod(spotinst.Int(v))
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			egWrapper := resourceObject.(*commons.ElastigroupWrapper)
			elastigroup := egWrapper.GetElastigroup()
			var value *int = nil
			if v, ok := resourceData.Get(string(HealthCheckGracePeriod)).(int); ok && v > 0 {
				value = spotinst.Int(v)
			}
			elastigroup.Compute.LaunchSpecification.SetHealthCheckGracePeriod(value)
			return nil
		},
		nil,
	)

	fieldsMap[HealthCheckType] = commons.NewGenericField(
		commons.ElastigroupAWS,
		HealthCheckType,
		&schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			egWrapper := resourceObject.(*commons.ElastigroupWrapper)
			elastigroup := egWrapper.GetElastigroup()
			var value *string = nil
			if elastigroup.Compute != nil && elastigroup.Compute.LaunchSpecification != nil &&
				elastigroup.Compute.LaunchSpecification.HealthCheckType != nil {
				value = elastigroup.Compute.LaunchSpecification.HealthCheckType
			}
			if err := resourceData.Set(string(HealthCheckType), spotinst.StringValue(value)); err != nil {
				return fmt.Errorf(string(commons.FailureFieldReadPattern), string(HealthCheckType), err)
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			egWrapper := resourceObject.(*commons.ElastigroupWrapper)
			elastigroup := egWrapper.GetElastigroup()
			if v, ok := resourceData.Get(string(HealthCheckType)).(string); ok && v != "" {
				elastigroup.Compute.LaunchSpecification.SetHealthCheckType(spotinst.String(v))
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			egWrapper := resourceObject.(*commons.ElastigroupWrapper)
			elastigroup := egWrapper.GetElastigroup()
			var value *string = nil
			if v, ok := resourceData.Get(string(HealthCheckType)).(string); ok && v != "" {
				value = spotinst.String(v)
			}
			elastigroup.Compute.LaunchSpecification.SetHealthCheckType(value)
			return nil
		},
		nil,
	)

	fieldsMap[HealthCheckUnhealthyDurationBeforeReplacement] = commons.NewGenericField(
		commons.ElastigroupAWS,
		HealthCheckUnhealthyDurationBeforeReplacement,
		&schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			egWrapper := resourceObject.(*commons.ElastigroupWrapper)
			elastigroup := egWrapper.GetElastigroup()
			var value *int = nil
			if elastigroup.Compute != nil && elastigroup.Compute.LaunchSpecification != nil &&
				elastigroup.Compute.LaunchSpecification.HealthCheckUnhealthyDurationBeforeReplacement != nil {
				value = elastigroup.Compute.LaunchSpecification.HealthCheckUnhealthyDurationBeforeReplacement
			}
			if err := resourceData.Set(string(HealthCheckUnhealthyDurationBeforeReplacement), spotinst.IntValue(value)); err != nil {
				return fmt.Errorf(string(commons.FailureFieldReadPattern), string(HealthCheckUnhealthyDurationBeforeReplacement), err)
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			egWrapper := resourceObject.(*commons.ElastigroupWrapper)
			elastigroup := egWrapper.GetElastigroup()
			if v, ok := resourceData.Get(string(HealthCheckUnhealthyDurationBeforeReplacement)).(int); ok && v > 0 {
				elastigroup.Compute.LaunchSpecification.SetHealthCheckUnhealthyDurationBeforeReplacement(spotinst.Int(v))
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			egWrapper := resourceObject.(*commons.ElastigroupWrapper)
			elastigroup := egWrapper.GetElastigroup()
			var value *int = nil
			if v, ok := resourceData.Get(string(HealthCheckUnhealthyDurationBeforeReplacement)).(int); ok && v > 0 {
				value = spotinst.Int(v)
			}
			elastigroup.Compute.LaunchSpecification.SetHealthCheckUnhealthyDurationBeforeReplacement(value)
			return nil
		},
		nil,
	)

	fieldsMap[Region] = commons.NewGenericField(
		commons.ElastigroupAWS,
		Region,
		&schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			egWrapper := resourceObject.(*commons.ElastigroupWrapper)
			elastigroup := egWrapper.GetElastigroup()
			var value *string = nil
			if elastigroup.Region != nil {
				value = elastigroup.Region
			}
			if err := resourceData.Set(string(Region), value); err != nil {
				return fmt.Errorf(string(commons.FailureFieldReadPattern), string(Region), err)
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			egWrapper := resourceObject.(*commons.ElastigroupWrapper)
			elastigroup := egWrapper.GetElastigroup()
			if v, ok := resourceData.GetOk(string(Region)); ok {
				elastigroup.SetRegion(spotinst.String(v.(string)))
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			egWrapper := resourceObject.(*commons.ElastigroupWrapper)
			elastigroup := egWrapper.GetElastigroup()
			if v, ok := resourceData.GetOk(string(Region)); ok {
				elastigroup.SetRegion(spotinst.String(v.(string)))
			}
			return nil
		},
		nil,
	)

	fieldsMap[SubnetIDs] = commons.NewGenericField(
		commons.ElastigroupAWS,
		SubnetIDs,
		&schema.Schema{
			Type:          schema.TypeList,
			Elem:          &schema.Schema{Type: schema.TypeString},
			ConflictsWith: []string{string(AvailabilityZones)},
			Optional:      true,
			Deprecated:    "This field will soon be deprecated and handled by availability_zones",
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			egWrapper := resourceObject.(*commons.ElastigroupWrapper)
			elastigroup := egWrapper.GetElastigroup()
			var value []string = nil
			if elastigroup.Compute != nil && elastigroup.Compute.SubnetIDs != nil {
				value = elastigroup.Compute.SubnetIDs
			}
			if err := resourceData.Set(string(SubnetIDs), value); err != nil {
				return fmt.Errorf(string(commons.FailureFieldReadPattern), string(SubnetIDs), err)
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			egWrapper := resourceObject.(*commons.ElastigroupWrapper)
			elastigroup := egWrapper.GetElastigroup()
			if value, ok := resourceData.GetOk(string(SubnetIDs)); ok && value != nil {
				if subnetIds, err := expandSubnetIDs(value); err != nil {
					return err
				} else {
					elastigroup.Compute.SetSubnetIDs(subnetIds)
				}
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			egWrapper := resourceObject.(*commons.ElastigroupWrapper)
			elastigroup := egWrapper.GetElastigroup()
			var value []string = nil
			if v, ok := resourceData.GetOk(string(SubnetIDs)); ok && v != nil {
				if subnetIds, err := expandSubnetIDs(v); err != nil {
					return err
				} else {
					value = subnetIds
				}
			}
			elastigroup.Compute.SetSubnetIDs(value)
			return nil
		},
		nil,
	)

	fieldsMap[PreferredAvailabilityZones] = commons.NewGenericField(
		commons.ElastigroupAWS,
		PreferredAvailabilityZones,
		&schema.Schema{
			Type:     schema.TypeList,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Optional: true,
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			egWrapper := resourceObject.(*commons.ElastigroupWrapper)
			elastigroup := egWrapper.GetElastigroup()
			var value []string = nil
			if elastigroup.Compute != nil && elastigroup.Compute.PreferredAvailabilityZones != nil {
				value = elastigroup.Compute.PreferredAvailabilityZones
			}
			if err := resourceData.Set(string(PreferredAvailabilityZones), value); err != nil {
				return fmt.Errorf(string(commons.FailureFieldReadPattern), string(PreferredAvailabilityZones), err)
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			egWrapper := resourceObject.(*commons.ElastigroupWrapper)
			elastigroup := egWrapper.GetElastigroup()
			if value, ok := resourceData.GetOk(string(PreferredAvailabilityZones)); ok {
				if preferredAZs, err := expandAWSGroupPreferredAvailabilityZones(value); err != nil {
					return err
				} else {
					elastigroup.Compute.SetPreferredAvailabilityZones(preferredAZs)
				}
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			egWrapper := resourceObject.(*commons.ElastigroupWrapper)
			elastigroup := egWrapper.GetElastigroup()
			var result []string = nil
			if value, ok := resourceData.GetOk(string(PreferredAvailabilityZones)); ok {
				if preferredAZs, err := expandAWSGroupPreferredAvailabilityZones(value); err != nil {
					return err
				} else {
					result = preferredAZs
				}
			}
			elastigroup.Compute.SetPreferredAvailabilityZones(result)
			return nil
		},
		nil,
	)

	fieldsMap[AvailabilityZones] = commons.NewGenericField(
		commons.ElastigroupAWS,
		AvailabilityZones,
		&schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					string(SubnetIDs): {
						Type:     schema.TypeList,
						Optional: true,
						Elem: &schema.Schema{
							Type: schema.TypeString},
					},
					string(AvailabilityZoneName): {
						Type:     schema.TypeString,
						Required: true,
					},
					string(PlacementGroupName): {
						Type:     schema.TypeString,
						Optional: true,
					},
				},
			},
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			egWrapper := resourceObject.(*commons.ElastigroupWrapper)
			elastigroup := egWrapper.GetElastigroup()
			var result []interface{} = nil
			if elastigroup.Compute != nil && elastigroup.Compute.AvailabilityZones != nil {
				az := elastigroup.Compute.AvailabilityZones
				result = flattenAvailabilityZones(az)
			}
			if result != nil {
				if err := resourceData.Set(string(AvailabilityZones), result); err != nil {
					return fmt.Errorf("failed to set availabilityZone configuration: %#v", err)
				}
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			egWrapper := resourceObject.(*commons.ElastigroupWrapper)
			elastigroup := egWrapper.GetElastigroup()
			if v, ok := resourceData.GetOk(string(AvailabilityZones)); ok {
				if availabilityZones, err := expandAvailabilityZones(v); err != nil {
					return err
				} else {
					elastigroup.Compute.SetAvailabilityZones(availabilityZones)
				}
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			egWrapper := resourceObject.(*commons.ElastigroupWrapper)
			elastigroup := egWrapper.GetElastigroup()
			var result []*aws.AvailabilityZone = nil
			if v, ok := resourceData.GetOk(string(AvailabilityZones)); ok {
				if availabilityZones, err := expandAvailabilityZones(v); err != nil {
					return err
				} else {
					result = availabilityZones
				}
			}
			elastigroup.Compute.SetAvailabilityZones(result)
			return nil
		},
		nil,
	)

	fieldsMap[ElasticLoadBalancers] = commons.NewGenericField(
		commons.ElastigroupAWS,
		ElasticLoadBalancers,
		&schema.Schema{
			Type:     schema.TypeList,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Optional: true,
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			egWrapper := resourceObject.(*commons.ElastigroupWrapper)
			elastigroup := egWrapper.GetElastigroup()
			var balNames []string = nil
			if elastigroup.Compute != nil && elastigroup.Compute.LaunchSpecification != nil &&
				elastigroup.Compute.LaunchSpecification.LoadBalancersConfig != nil &&
				elastigroup.Compute.LaunchSpecification.LoadBalancersConfig.LoadBalancers != nil {

				balancers := elastigroup.Compute.LaunchSpecification.LoadBalancersConfig.LoadBalancers
				for _, balancer := range balancers {
					balType := spotinst.StringValue(balancer.Type)
					if strings.ToUpper(balType) == string(BalancerTypeClassic) {
						balName := spotinst.StringValue(balancer.Name)
						balNames = append(balNames, balName)
					}
				}
			}
			return resourceData.Set(string(ElasticLoadBalancers), balNames)
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			egWrapper := resourceObject.(*commons.ElastigroupWrapper)
			elastigroup := egWrapper.GetElastigroup()
			if balNames, ok := resourceData.GetOk(string(ElasticLoadBalancers)); ok {
				var fn = func(name string) (*aws.LoadBalancer, error) {
					return &aws.LoadBalancer{
						Type: spotinst.String(strings.ToUpper(string(BalancerTypeClassic))),
						Name: spotinst.String(name),
					}, nil
				}
				if elbBalancers, err := expandBalancersContent(balNames, fn); err != nil {
					return err
				} else if len(elbBalancers) > 0 {
					existingBalancers := elastigroup.Compute.LaunchSpecification.LoadBalancersConfig.LoadBalancers
					if len(existingBalancers) > 0 {
						elbBalancers = append(elbBalancers, existingBalancers...)
					}
					elastigroup.Compute.LaunchSpecification.LoadBalancersConfig.SetLoadBalancers(elbBalancers)
				}
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			return onBalancersUpdate(resourceObject.(*commons.ElastigroupWrapper), resourceData)
		},
		nil,
	)

	fieldsMap[TargetGroupArns] = commons.NewGenericField(
		commons.ElastigroupAWS,
		TargetGroupArns,
		&schema.Schema{
			Type:     schema.TypeList,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Optional: true,
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			egWrapper := resourceObject.(*commons.ElastigroupWrapper)
			elastigroup := egWrapper.GetElastigroup()
			var tgArns []string = nil
			if elastigroup.Compute != nil && elastigroup.Compute.LaunchSpecification != nil &&
				elastigroup.Compute.LaunchSpecification.LoadBalancersConfig != nil &&
				elastigroup.Compute.LaunchSpecification.LoadBalancersConfig.LoadBalancers != nil {

				balancers := elastigroup.Compute.LaunchSpecification.LoadBalancersConfig.LoadBalancers
				for _, balancer := range balancers {
					balType := spotinst.StringValue(balancer.Type)
					if balType == string(BalancerTypeTargetGroup) {
						arn := spotinst.StringValue(balancer.Arn)
						tgArns = append(tgArns, arn)
					}
				}
			}
			return resourceData.Set(string(TargetGroupArns), tgArns)
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			egWrapper := resourceObject.(*commons.ElastigroupWrapper)
			elastigroup := egWrapper.GetElastigroup()
			if tgArns, ok := resourceData.GetOk(string(TargetGroupArns)); ok {
				var fn = func(arn string) (*aws.LoadBalancer, error) {
					// Name should be removed as a mandatory field in the future
					if name, err := extractTargetGroupFromArn(arn); err != nil {
						return nil, err
					} else {
						return &aws.LoadBalancer{
							Type: spotinst.String(strings.ToUpper(string(BalancerTypeTargetGroup))),
							Arn:  spotinst.String(arn),
							Name: spotinst.String(name),
						}, nil
					}
				}
				// Existing balancers appended if needed inside expand method
				if tgBalancers, err := expandBalancersContent(tgArns, fn); err != nil {
					return err
				} else {
					if len(tgBalancers) > 0 {
						existingBalancers := elastigroup.Compute.LaunchSpecification.LoadBalancersConfig.LoadBalancers
						if len(existingBalancers) > 0 {
							tgBalancers = append(tgBalancers, existingBalancers...)
						}
						elastigroup.Compute.LaunchSpecification.LoadBalancersConfig.SetLoadBalancers(tgBalancers)
					}
				}
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			egWrapper := resourceObject.(*commons.ElastigroupWrapper)
			if err := onBalancersUpdate(egWrapper, resourceData); err != nil {
				return err
			}
			return nil
		},
		nil,
	)

	fieldsMap[Tags] = commons.NewGenericField(
		commons.ElastigroupAWS,
		Tags,
		&schema.Schema{
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					string(TagKey): {
						Type:     schema.TypeString,
						Optional: true,
					},

					string(TagValue): {
						Type:     schema.TypeString,
						Optional: true,
					},
				},
			},
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			egWrapper := resourceObject.(*commons.ElastigroupWrapper)
			elastigroup := egWrapper.GetElastigroup()
			var result []interface{} = nil
			if elastigroup.Compute != nil && elastigroup.Compute.LaunchSpecification != nil &&
				elastigroup.Compute.LaunchSpecification.Tags != nil {
				tags := elastigroup.Compute.LaunchSpecification.Tags
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
			egWrapper := resourceObject.(*commons.ElastigroupWrapper)
			elastigroup := egWrapper.GetElastigroup()
			if value, ok := resourceData.GetOk(string(Tags)); ok {
				if tags, err := expandTags(value); err != nil {
					return err
				} else {
					elastigroup.Compute.LaunchSpecification.SetTags(tags)
				}
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			egWrapper := resourceObject.(*commons.ElastigroupWrapper)
			elastigroup := egWrapper.GetElastigroup()
			var tagsToAdd []*aws.Tag = nil
			if value, ok := resourceData.GetOk(string(Tags)); ok {
				if tags, err := expandTags(value); err != nil {
					return err
				} else {
					tagsToAdd = tags
				}
			}
			elastigroup.Compute.LaunchSpecification.SetTags(tagsToAdd)
			return nil
		},
		nil,
	)

	fieldsMap[ElasticIps] = commons.NewGenericField(
		commons.ElastigroupAWS,
		ElasticIps,
		&schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			egWrapper := resourceObject.(*commons.ElastigroupWrapper)
			elastigroup := egWrapper.GetElastigroup()
			var result []string
			if elastigroup.Compute != nil && elastigroup.Compute.ElasticIPs != nil {
				result = append(result, elastigroup.Compute.ElasticIPs...)
			}
			if err := resourceData.Set(string(ElasticIps), result); err != nil {
				return fmt.Errorf(string(commons.FailureFieldReadPattern), string(ElasticIps), err)
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			egWrapper := resourceObject.(*commons.ElastigroupWrapper)
			elastigroup := egWrapper.GetElastigroup()
			if value, ok := resourceData.GetOk(string(ElasticIps)); ok {
				if eips, err := expandAWSGroupElasticIPs(value); err != nil {
					return err
				} else {
					elastigroup.Compute.SetElasticIPs(eips)
				}
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			egWrapper := resourceObject.(*commons.ElastigroupWrapper)
			elastigroup := egWrapper.GetElastigroup()
			var result []string = nil
			if value, ok := resourceData.GetOk(string(ElasticIps)); ok {
				if eips, err := expandAWSGroupElasticIPs(value); err != nil {
					return err
				} else {
					result = eips
				}
			}
			elastigroup.Compute.SetElasticIPs(result)
			return nil
		},
		nil,
	)

	fieldsMap[RevertToSpot] = commons.NewGenericField(
		commons.ElastigroupAWS,
		RevertToSpot,
		&schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					string(PerformAt): {
						Type:     schema.TypeString,
						Required: true,
					},

					string(TimeWindow): {
						Type:     schema.TypeList,
						Optional: true,
						Elem:     &schema.Schema{Type: schema.TypeString},
					},
				},
			},
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			egWrapper := resourceObject.(*commons.ElastigroupWrapper)
			elastigroup := egWrapper.GetElastigroup()
			if elastigroup.Strategy != nil && elastigroup.Strategy.RevertToSpot != nil {
				rts := elastigroup.Strategy.RevertToSpot
				result := make(map[string]interface{})
				result[string(PerformAt)] = spotinst.StringValue(rts.PerformAt)
				result[string(TimeWindow)] = rts.TimeWindows
				revertToSpot := []interface{}{result}
				if err := resourceData.Set(string(RevertToSpot), revertToSpot); err != nil {
					return fmt.Errorf("failed to set revertToSpot configuration: %#v", err)
				}
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			egWrapper := resourceObject.(*commons.ElastigroupWrapper)
			elastigroup := egWrapper.GetElastigroup()
			if v, ok := resourceData.GetOk(string(RevertToSpot)); ok {
				if revertToSpot, err := expandAWSGroupRevertToSpot(v); err != nil {
					return err
				} else {
					elastigroup.Strategy.SetRevertToSpot(revertToSpot)
				}
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			egWrapper := resourceObject.(*commons.ElastigroupWrapper)
			elastigroup := egWrapper.GetElastigroup()
			var revertToSpot *aws.RevertToSpot = nil
			if v, ok := resourceData.GetOk(string(RevertToSpot)); ok {
				if rts, err := expandAWSGroupRevertToSpot(v); err != nil {
					return err
				} else {
					revertToSpot = rts
				}
			}
			elastigroup.Strategy.SetRevertToSpot(revertToSpot)
			return nil
		},
		nil,
	)

	fieldsMap[Signal] = commons.NewGenericField(
		commons.ElastigroupAWS,
		Signal,
		&schema.Schema{
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					string(SignalName): {
						Type:     schema.TypeString,
						Required: true,
					},

					string(SignalTimeout): {
						Type:     schema.TypeInt,
						Optional: true,
					},
				},
			},
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			egWrapper := resourceObject.(*commons.ElastigroupWrapper)
			elastigroup := egWrapper.GetElastigroup()
			var signalsToAdd []interface{} = nil
			if elastigroup.Strategy != nil && elastigroup.Strategy.Signals != nil {
				signals := elastigroup.Strategy.Signals
				signalsToAdd = make([]interface{}, 0, len(signals))
				for _, s := range signals {
					m := make(map[string]interface{})
					m[string(SignalName)] = spotinst.StringValue(s.Name)
					m[string(SignalTimeout)] = spotinst.IntValue(s.Timeout)
					signalsToAdd = append(signalsToAdd, m)
				}
			}
			if err := resourceData.Set(string(Signal), signalsToAdd); err != nil {
				return fmt.Errorf("failed to set signals configuration: %#v", err)
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			egWrapper := resourceObject.(*commons.ElastigroupWrapper)
			elastigroup := egWrapper.GetElastigroup()
			if v, ok := resourceData.GetOk(string(Signal)); ok {
				if signals, err := expandSignals(v); err != nil {
					return err
				} else {
					elastigroup.Strategy.SetSignals(signals)
				}
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			egWrapper := resourceObject.(*commons.ElastigroupWrapper)
			elastigroup := egWrapper.GetElastigroup()
			var signalsToAdd []*aws.Signal = nil
			if v, ok := resourceData.GetOk(string(Signal)); ok {
				if signals, err := expandSignals(v); err != nil {
					return err
				} else {
					signalsToAdd = signals
				}
			}
			if elastigroup.Strategy == nil {
				elastigroup.SetStrategy(&aws.Strategy{})
			}
			elastigroup.Strategy.SetSignals(signalsToAdd)
			return nil
		},
		nil,
	)

	fieldsMap[UpdatePolicy] = commons.NewGenericField(
		commons.ElastigroupAWS,
		UpdatePolicy,
		&schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					string(ShouldResumeStateful): {
						Type:     schema.TypeBool,
						Required: true,
					},

					string(AutoApplyTags): {
						Type:     schema.TypeBool,
						Optional: true,
					},

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

								string(GracePeriod): {
									Type:     schema.TypeInt,
									Optional: true,
									Default:  -1,
								},

								string(HealthCheckType): {
									Type:     schema.TypeString,
									Optional: true,
								},

								string(WaitForRollPct): {
									Type:     schema.TypeFloat,
									Optional: true,
								},

								string(WaitForRollTimeout): {
									Type:     schema.TypeInt,
									Optional: true,
								},

								string(Strategy): {
									Type:     schema.TypeList,
									Optional: true,
									MaxItems: 1,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											string(Action): {
												Type:     schema.TypeString,
												Required: true,
											},

											string(ShouldDrainInstances): {
												Type:     schema.TypeBool,
												Optional: true,
											},

											string(BatchMinHealthyPercentage): {
												Type:     schema.TypeInt,
												Optional: true,
												Default:  50,
											},

											string(OnFailure): {
												Type:     schema.TypeList,
												Optional: true,
												MaxItems: 1,
												Elem: &schema.Resource{
													Schema: map[string]*schema.Schema{
														string(ActionType): {
															Type:     schema.TypeString,
															Required: true,
														},

														string(ShouldHandleAllBatches): {
															Type:     schema.TypeBool,
															Optional: true,
														},

														string(BatchNum): {
															Type:     schema.TypeInt,
															Optional: true,
														},

														string(DrainingTimeout): {
															Type:     schema.TypeInt,
															Optional: true,
														},

														string(ShouldDecrementTargetCapacity): {
															Type:     schema.TypeBool,
															Optional: true,
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		nil, nil, nil, nil,
	)

	fieldsMap[WaitForCapacity] = commons.NewGenericField(
		commons.ElastigroupAWS,
		WaitForCapacity,
		&schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
		},
		nil, nil, nil, nil,
	)

	fieldsMap[WaitForCapacityTimeout] = commons.NewGenericField(
		commons.ElastigroupAWS,
		WaitForCapacityTimeout,
		&schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
		},
		nil, nil, nil, nil,
	)
}

var TargetGroupArnRegex = regexp.MustCompile(`arn:aws:elasticloadbalancing:.*:\d{12}:targetgroup/(.*)/.*`)

func expandAvailabilityZonesSlice(data interface{}) ([]*aws.AvailabilityZone, error) {
	list := data.([]interface{})
	zones := make([]*aws.AvailabilityZone, 0, len(list))
	for _, str := range list {
		if s, ok := str.(string); ok {
			parts := strings.Split(s, ":")
			zone := &aws.AvailabilityZone{}
			if len(parts) >= 1 && parts[0] != "" {
				zone.SetName(spotinst.String(parts[0]))
			}
			if len(parts) >= 2 && parts[1] != "" {
				zone.SetSubnetId(spotinst.String(parts[1]))
			}
			if len(parts) >= 3 && parts[2] != "" {
				zone.SetPlacementGroupName(spotinst.String(parts[2]))
			}
			zones = append(zones, zone)
		}
	}

	return zones, nil
}

func expandAWSGroupElasticIPs(data interface{}) ([]string, error) {
	list := data.([]interface{})
	eips := make([]string, 0, len(list))
	for _, str := range list {
		if eip, ok := str.(string); ok {
			eips = append(eips, eip)
		}
	}
	return eips, nil
}

func expandAWSGroupPreferredAvailabilityZones(data interface{}) ([]string, error) {
	list := data.([]interface{})
	result := make([]string, 0, len(list))

	for _, v := range list {
		if preferredAZ, ok := v.(string); ok && preferredAZ != "" {
			result = append(result, preferredAZ)
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

type CreateBalancerObjFunc func(id string) (*aws.LoadBalancer, error)

func expandBalancersContent(balancersIdentifiers interface{}, fn CreateBalancerObjFunc) ([]*aws.LoadBalancer, error) {
	if balancersIdentifiers == nil {
		return nil, nil
	}
	list := balancersIdentifiers.([]interface{})
	balancers := make([]*aws.LoadBalancer, 0, len(list))
	for _, str := range list {
		if id, ok := str.(string); ok && id != "" {
			if lb, err := fn(id); err != nil {
				return nil, err
			} else {
				balancers = append(balancers, lb)
			}
		}
	}
	return balancers, nil
}

func extractBalancers(
	balancerType BalancerType,
	elastigroup *aws.Group,
	resourceData *schema.ResourceData) ([]*aws.LoadBalancer, error) {

	existingBalancers := elastigroup.Compute.LaunchSpecification.LoadBalancersConfig.LoadBalancers

	var elbBalancers []*aws.LoadBalancer = nil
	var tgBalancers []*aws.LoadBalancer = nil

	if len(existingBalancers) > 0 {
		for _, balancer := range existingBalancers {
			balTypeStr := spotinst.StringValue(balancer.Type)

			switch balTypeStr {
			case string(BalancerTypeClassic):
				{
					elbBalancers = append(elbBalancers, balancer)
					break
				}
			case string(BalancerTypeTargetGroup):
				{
					tgBalancers = append(tgBalancers, balancer)
					break
				}
			}
		}
	}

	if elbNames, ok := resourceData.GetOk(string(ElasticLoadBalancers)); ok && balancerType == BalancerTypeClassic {
		var fn = func(name string) (*aws.LoadBalancer, error) {
			return &aws.LoadBalancer{
				Type: spotinst.String(strings.ToUpper(string(BalancerTypeClassic))),
				Name: spotinst.String(name),
			}, nil
		}
		if tfElbs, err := expandBalancersContent(elbNames, fn); err != nil {
			return nil, err
		} else {
			elbBalancers = append(tfElbs, elbBalancers...)
		}
	}

	if tgArns, ok := resourceData.GetOk(string(TargetGroupArns)); ok && balancerType == BalancerTypeTargetGroup {
		var fn = func(arn string) (*aws.LoadBalancer, error) {
			// Name should be removed as a mandatory field in the future
			if name, err := extractTargetGroupFromArn(arn); err != nil {
				return nil, err
			} else {
				return &aws.LoadBalancer{
					Type: spotinst.String(strings.ToUpper(string(BalancerTypeTargetGroup))),
					Arn:  spotinst.String(arn),
					Name: spotinst.String(name),
				}, nil
			}
		}
		if tfTargetGroups, err := expandBalancersContent(tgArns, fn); err != nil {
			return nil, err
		} else {
			tgBalancers = append(tfTargetGroups, tgBalancers...)
		}
	}

	var result []*aws.LoadBalancer = nil
	if balancerType == BalancerTypeClassic {
		result = elbBalancers
	} else if balancerType == BalancerTypeTargetGroup {
		result = tgBalancers
	}
	return result, nil
}

func onBalancersUpdate(egWrapper *commons.ElastigroupWrapper, resourceData *schema.ResourceData) error {
	var elbNullify = false
	var tgNullify = false

	elastigroup := egWrapper.GetElastigroup()

	if !egWrapper.StatusElbUpdated {
		if elbBalancers, err := extractBalancers(BalancerTypeClassic, elastigroup, resourceData); err != nil {
			return err
		} else if len(elbBalancers) > 0 {
			existingBalancers := elastigroup.Compute.LaunchSpecification.LoadBalancersConfig.LoadBalancers
			if len(existingBalancers) > 0 {
				elbBalancers = append(elbBalancers, existingBalancers...)
			}
			elastigroup.Compute.LaunchSpecification.LoadBalancersConfig.SetLoadBalancers(elbBalancers)
		} else {
			elbNullify = true
		}
		egWrapper.StatusElbUpdated = true
	}
	if !egWrapper.StatusTgUpdated {
		if tgBalancers, err := extractBalancers(BalancerTypeTargetGroup, elastigroup, resourceData); err != nil {
			return err
		} else if len(tgBalancers) > 0 {
			existingBalancers := elastigroup.Compute.LaunchSpecification.LoadBalancersConfig.LoadBalancers
			if len(existingBalancers) > 0 {
				tgBalancers = append(tgBalancers, existingBalancers...)
			}
			elastigroup.Compute.LaunchSpecification.LoadBalancersConfig.SetLoadBalancers(tgBalancers)
		} else {
			tgNullify = true
		}
		egWrapper.StatusTgUpdated = true
	}

	// All fields share the same object structure, we need to nullify if and only if there are no items
	// from all types
	if elbNullify && tgNullify {
		elastigroup.Compute.LaunchSpecification.LoadBalancersConfig.SetLoadBalancers(nil)
	}
	return nil
}

func expandSignals(data interface{}) ([]*aws.Signal, error) {
	list := data.(*schema.Set).List()
	signals := make([]*aws.Signal, 0, len(list))
	for _, item := range list {
		m := item.(map[string]interface{})
		signal := &aws.Signal{}

		if v, ok := m[string(SignalName)].(string); ok && v != "" {
			signal.SetName(spotinst.String(strings.ToUpper(v)))
		}

		if v, ok := m[string(SignalTimeout)].(int); ok && v > 0 {
			signal.SetTimeout(spotinst.Int(v))
		}
		signals = append(signals, signal)
	}

	return signals, nil
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

func expandAWSGroupRevertToSpot(data interface{}) (*aws.RevertToSpot, error) {
	revertToSpot := &aws.RevertToSpot{}
	list := data.([]interface{})
	if list != nil && list[0] != nil {
		m := list[0].(map[string]interface{})

		var performAt *string = nil
		if v, ok := m[string(PerformAt)].(string); ok {
			performAt = spotinst.String(v)
		}
		revertToSpot.SetPerformAt(performAt)

		var timeWindows []string = nil
		if v, ok := m[string(TimeWindow)].([]interface{}); ok && len(v) > 0 {
			ids := make([]string, 0, len(v))
			for _, id := range v {
				if v, ok := id.(string); ok && len(v) > 0 {
					ids = append(ids, v)
				}
			}
			timeWindows = ids
		}
		revertToSpot.SetTimeWindows(timeWindows)
	}
	//log.Printf("[DEBUG] Group revert to spot configuration: %s", stringutil.Stringify(revertToSpot))
	return revertToSpot, nil
}

func extractTargetGroupFromArn(arn string) (string, error) {
	name := ""
	success := false
	var match = TargetGroupArnRegex.FindStringSubmatch(arn)
	if len(match) == 2 {
		name = match[1]
		success = name != ""
	}
	if !success {
		return "", fmt.Errorf("cannot determine targret group name from target group arn")
	}
	return name, nil
}

func expandAvailabilityZones(data interface{}) ([]*aws.AvailabilityZone, error) {
	if list := data.([]interface{}); len(list) > 0 {
		availabilityZones := make([]*aws.AvailabilityZone, 0, len(list))
		for _, item := range list {
			m := item.(map[string]interface{})
			availabilityZone := &aws.AvailabilityZone{}

			if v, ok := m[string(AvailabilityZoneName)].(string); ok && v != "" {
				availabilityZone.SetName(spotinst.String(v))
			}

			if v, ok := m[string(PlacementGroupName)].(string); ok && v != "" {
				availabilityZone.SetPlacementGroupName(spotinst.String(v))
			}

			if v, ok := m[string(SubnetIDs)]; ok && len(v.([]interface{})) > 0 {
				if subnetIDs, err := expandSubnetIDs(v); err != nil {
					return nil, err
				} else {
					availabilityZone.SetSubnetIDs(subnetIDs)
				}
			}
			availabilityZones = append(availabilityZones, availabilityZone)
		}
		return availabilityZones, nil
	}
	return nil, nil

}

func flattenAvailabilityZones(availabilityZones []*aws.AvailabilityZone) []interface{} {
	result := make([]interface{}, 0, len(availabilityZones))

	for _, availabilityZone := range availabilityZones {
		m := make(map[string]interface{})
		if availabilityZone.SubnetIDs != nil {
			m[string(SubnetIDs)] = spotinst.StringSlice(availabilityZone.SubnetIDs)
		}
		m[string(AvailabilityZoneName)] = spotinst.StringValue(availabilityZone.Name)
		m[string(PlacementGroupName)] = spotinst.StringValue(availabilityZone.PlacementGroupName)

		result = append(result, m)
	}
	return result
}
