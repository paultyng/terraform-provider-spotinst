package spotinst

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/spotinst/spotinst-sdk-go/service/managedinstance/providers/aws"
	"github.com/spotinst/spotinst-sdk-go/spotinst"
	"github.com/spotinst/spotinst-sdk-go/spotinst/client"
	"github.com/spotinst/terraform-provider-spotinst/spotinst/commons"
	"github.com/spotinst/terraform-provider-spotinst/spotinst/commons/managed_instance_aws"
	"github.com/spotinst/terraform-provider-spotinst/spotinst/commons/managed_instance_aws_compute"
	"github.com/spotinst/terraform-provider-spotinst/spotinst/commons/managed_instance_aws_compute_instance_type"
	"github.com/spotinst/terraform-provider-spotinst/spotinst/commons/managed_instance_aws_compute_launchspecification"
	"github.com/spotinst/terraform-provider-spotinst/spotinst/commons/managed_instance_aws_integrations"
	"github.com/spotinst/terraform-provider-spotinst/spotinst/commons/managed_instance_healthcheck"
	"github.com/spotinst/terraform-provider-spotinst/spotinst/commons/managed_instance_persistence"
	"github.com/spotinst/terraform-provider-spotinst/spotinst/commons/managed_instance_scheduling"
	"github.com/spotinst/terraform-provider-spotinst/spotinst/commons/managed_instance_strategy"
	"github.com/spotinst/terraform-provider-spotinst/spotinst/commons/managed_instances_aws_compute_launchspecification_networkinterfaces"
)

func resourceSpotinstMangedInstanceAWS() *schema.Resource {
	setupMangedInstanceResource()

	return &schema.Resource{
		CreateContext: resourceSpotinstManagedInstanceAWSCreate,
		ReadContext:   resourceSpotinstManagedInstanceAWSRead,
		UpdateContext: resourceSpotinstManagedInstanceAWSUpdate,
		DeleteContext: resourceSpotinstManagedInstanceAWSDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: commons.ManagedInstanceResource.GetSchemaMap(),
	}
}

func setupMangedInstanceResource() {
	fieldsMap := make(map[commons.FieldName]*commons.GenericField)

	managed_instance_aws.Setup(fieldsMap)
	managed_instance_strategy.Setup(fieldsMap)
	managed_instance_persistence.Setup(fieldsMap)
	managed_instance_healthcheck.Setup(fieldsMap)
	managed_instance_aws_compute.Setup(fieldsMap)
	managed_instance_aws_integrations.Setup(fieldsMap)
	managed_instance_scheduling.Setup(fieldsMap)
	managed_instances_aws_compute_launchspecification_networkinterfaces.Setup(fieldsMap)
	managed_instance_aws_compute_launchspecification.Setup(fieldsMap)
	managed_instance_aws_compute_instance_type.Setup(fieldsMap)

	commons.ManagedInstanceResource = commons.NewManagedInstanceResource(fieldsMap)
}

const ErrCodeManagedInstanceDoesntExist = "MANAGED_INSTANCE_DOESNT_EXIST"

func resourceSpotinstManagedInstanceAWSRead(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	id := resourceData.Id()
	log.Printf(string(commons.ResourceOnRead),
		commons.ManagedInstanceResource.GetName(), id)

	input := &aws.ReadManagedInstanceInput{ManagedInstanceID: spotinst.String(id)}
	resp, err := meta.(*Client).managedInstance.CloudProviderAWS().Read(context.Background(), input)
	if err != nil {
		// If the group was not found, return nil so that we can show
		// that the group does not exist
		if errs, ok := err.(client.Errors); ok && len(errs) > 0 {
			for _, err := range errs {
				if err.Code == ErrCodeManagedInstanceDoesntExist {
					resourceData.SetId("")
					return nil
				}
			}
		}

		// Some other error, report it.
		return diag.Errorf("failed to read ManagedInstance: %s", err)
	}

	// If nothing was found, then return no state.
	managedInstanceResponse := resp.ManagedInstance
	if managedInstanceResponse == nil {
		resourceData.SetId("")
		return nil
	}

	if err := commons.ManagedInstanceResource.OnRead(managedInstanceResponse, resourceData, meta); err != nil {
		return diag.FromErr(err)
	}
	log.Printf("===> ManagedInstance read successfully: %s <===", id)
	return nil
}

func resourceSpotinstManagedInstanceAWSCreate(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf(string(commons.ResourceOnCreate),
		commons.ManagedInstanceResource.GetName())

	mangedInstance, err := commons.ManagedInstanceResource.OnCreate(resourceData, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	ManagedInstanceId, err := createManagedInstance(resourceData, mangedInstance, meta.(*Client))
	if err != nil {
		return diag.FromErr(err)
	}

	resourceData.SetId(spotinst.StringValue(ManagedInstanceId))

	log.Printf("===> ManagedInstance created successfully: %s <===", resourceData.Id())

	return resourceSpotinstManagedInstanceAWSRead(ctx, resourceData, meta)
}

func createManagedInstance(resourceData *schema.ResourceData, mangedInstance *aws.ManagedInstance, spotinstClient *Client) (*string, error) {
	if json, err := commons.ToJson(mangedInstance); err != nil {
		return nil, err
	} else {
		log.Printf("===> ManagedInstance create configuration: %s", json)
	}
	if v, ok := resourceData.Get(string(managed_instance_aws_compute_launchspecification.IAMInstanceProfile)).(string); ok && v != "" {
		time.Sleep(5 * time.Second)
	}

	var resp *aws.CreateManagedInstanceOutput = nil
	err := resource.RetryContext(context.Background(), time.Minute, func() *resource.RetryError {
		input := &aws.CreateManagedInstanceInput{ManagedInstance: mangedInstance}
		r, err := spotinstClient.managedInstance.CloudProviderAWS().Create(context.Background(), input)
		if err != nil {
			// Checks whether we should retry the group creation.
			if errs, ok := err.(client.Errors); ok && len(errs) > 0 {
				for _, err := range errs {
					if err.Code == "InvalidParameterValue" &&
						strings.Contains(err.Message, "Invalid IAM Instance Profile") {
						return resource.RetryableError(err)
					}
				}
			}

			// Some other error, report it.
			return resource.NonRetryableError(err)
		}
		resp = r
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("[ERROR] failed to create ManagedInstance: %s", err)
	}
	return resp.ManagedInstance.ID, nil
}

func resourceSpotinstManagedInstanceAWSUpdate(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	id := resourceData.Id()
	log.Printf(string(commons.ResourceOnUpdate),
		commons.ManagedInstanceResource.GetName(), id)

	shouldUpdate, managedInstance, err := commons.ManagedInstanceResource.OnUpdate(resourceData, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	if shouldUpdate {
		managedInstance.SetId(spotinst.String(id))
		if err := updateAWSManagedInstance(managedInstance, resourceData, meta); err != nil {
			return diag.FromErr(err)
		}
	}

	log.Printf("===> ManagedInstance updated successfully: %s <===", id)
	return resourceSpotinstManagedInstanceAWSRead(ctx, resourceData, meta)
}

func updateAWSManagedInstance(managedInstance *aws.ManagedInstance, resourceData *schema.ResourceData, meta interface{}) error {
	var input = &aws.UpdateManagedInstanceInput{
		ManagedInstance: managedInstance,
	}

	if instanceActions, exists := resourceData.GetOk(string(managed_instance_aws.ManagedInstanceAction)); exists {
		actionList := instanceActions.([]interface{})

		ctx := context.TODO()
		svc := meta.(*Client).managedInstance.CloudProviderAWS()

		for _, action := range actionList {
			var (
				actionMap  = action.(map[string]interface{})
				actionType = actionMap[string(managed_instance_aws.ActionType)].(string)
				err        error
			)
			switch strings.ToLower(actionType) {
			case "pause":
				err = pauseManagedInstance(ctx, svc, resourceData.Id())
			case "resume":
				err = resumeManagedInstance(ctx, svc, resourceData.Id())
			case "recycle":
				err = recycleManagedInstance(ctx, svc, resourceData.Id())
			default:
				err = fmt.Errorf("unsupported action %q on managed instance %q", actionType, resourceData.Id())
			}
			if err != nil {
				log.Printf("[ERROR] managed instance (%s) action failed with error: %v", resourceData.Id(), err)
				return err
			}
		}
	}

	if json, err := commons.ToJson(managedInstance); err != nil {
		return err
	} else {
		log.Printf("===> ManagedInstance update configuration: %s", json)
	}

	if _, err := meta.(*Client).managedInstance.CloudProviderAWS().Update(context.Background(), input); err != nil {
		return fmt.Errorf("[ERROR] Failed to update managed instance [%v]: %v", resourceData.Id(), err)
	}

	return nil
}

func pauseManagedInstance(ctx context.Context, svc aws.Service, instanceID string) error {
	log.Printf("Pausing managed instance (%s)", instanceID)

	input := &aws.PauseManagedInstanceInput{
		ManagedInstanceID: spotinst.String(instanceID),
	}
	_, err := svc.Pause(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to pause managed instance (%s): %v", instanceID, err)
	}

	log.Printf("Successfully paused managed instance (%s)", instanceID)
	return nil
}

func resumeManagedInstance(ctx context.Context, svc aws.Service, instanceID string) error {
	log.Printf("Resuming managed instance (%s)", instanceID)

	input := &aws.ResumeManagedInstanceInput{
		ManagedInstanceID: spotinst.String(instanceID),
	}
	_, err := svc.Resume(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to resume managed instance (%s): %v", instanceID, err)
	}

	log.Printf("Successfully resumed managed instance (%s)", instanceID)
	return nil
}

func recycleManagedInstance(ctx context.Context, svc aws.Service, instanceID string) error {
	log.Printf("Recycling managed instance (%s)", instanceID)

	input := &aws.RecycleManagedInstanceInput{
		ManagedInstanceID: spotinst.String(instanceID),
	}
	_, err := svc.Recycle(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to recycle managed instance (%s): %v", instanceID, err)
	}

	log.Printf("Successfully recycled managed instance (%s)", instanceID)
	return nil
}

func resourceSpotinstManagedInstanceAWSDelete(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	id := resourceData.Id()
	log.Printf(string(commons.ResourceOnDelete),
		commons.ManagedInstanceResource.GetName(), id)

	if err := deleteManagedInstance(resourceData, meta); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("===> ManagedInstance deleted successfully: %s <===", resourceData.Id())
	resourceData.SetId("")
	return nil
}

func deleteManagedInstance(resourceData *schema.ResourceData, meta interface{}) error {
	managedInstanceId := resourceData.Id()
	input := &aws.DeleteManagedInstanceInput{
		ManagedInstanceID: spotinst.String(managedInstanceId),
		AMIBackup: &aws.AMIBackup{
			ShouldDeleteImages: spotinst.Bool(true),
		},
		DeallocationConfig: &aws.DeallocationConfig{
			ShouldDeleteImages:            spotinst.Bool(true),
			ShouldTerminateInstance:       spotinst.Bool(true),
			ShouldDeleteVolumes:           spotinst.Bool(true),
			ShouldDeleteNetworkInterfaces: spotinst.Bool(true),
		},
	}
	if json, err := commons.ToJson(input); err != nil {
		return err
	} else {
		log.Printf("===> ManagedInstance delete configuration: %s", json)
	}

	if _, err := meta.(*Client).managedInstance.CloudProviderAWS().Delete(context.Background(), input); err != nil {
		return fmt.Errorf("[ERROR] onDelete() -> Failed to delete ManagedInstance: %s", err)
	}
	return nil
}
