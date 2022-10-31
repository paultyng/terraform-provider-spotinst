package spotinst

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/spotinst/spotinst-sdk-go/service/elastigroup/providers/aws"
	"github.com/spotinst/spotinst-sdk-go/spotinst"
	"github.com/spotinst/spotinst-sdk-go/spotinst/client"
	"github.com/spotinst/terraform-provider-spotinst/spotinst/commons"
	"github.com/spotinst/terraform-provider-spotinst/spotinst/elastigroup_aws_suspend_processes"
)

func resourceSpotinstElastigroupSuspendProcesses() *schema.Resource {
	setupSuspendProcesses()

	return &schema.Resource{
		CreateContext: resourceSpotinstAWSSuspendProcessesCreate,
		ReadContext:   resourceSpotinstAWSSuspendProcessesRead,
		UpdateContext: resourceSpotinstAWSSuspendProcessesUpdate,
		DeleteContext: resourceSpotinstAWSSuspendProcessesDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: commons.SuspendProcessesResource.GetSchemaMap(),
	}
}

func setupSuspendProcesses() {
	fieldsMap := make(map[commons.FieldName]*commons.GenericField)

	elastigroup_aws_suspend_processes.Setup(fieldsMap)

	commons.SuspendProcessesResource = commons.NewSuspendProcessesResource(fieldsMap)
}

const ErrCodeSuspendProcessesNotFound = "SUSPEND_PROCESSES_DOESNT_EXIST"

func resourceSpotinstAWSSuspendProcessesRead(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if resourceData.Id() == "" {
		resourceData.SetId(resourceData.Get(string(elastigroup_aws_suspend_processes.GroupID)).(string))
		resourceData.Set("group_id", resourceData.Get(string(elastigroup_aws_suspend_processes.GroupID)))
	}

	log.Printf(string(commons.ResourceOnRead), commons.SuspendProcessesResource.GetName(), resourceData.Id())

	input := &aws.ListSuspensionsInput{}
	gID := resourceData.Id()
	input.GroupID = &gID
	resp, err := meta.(*Client).elastigroup.CloudProviderAWS().ListSuspensions(context.Background(), input)
	if err != nil {
		// If the group was not found, return nil so that we can show
		// that the group is gone.
		if errs, ok := err.(client.Errors); ok && len(errs) > 0 {
			for _, err := range errs {
				if err.Code == ErrCodeSuspendProcessesNotFound {
					resourceData.SetId("")
					return nil
				}
			}
		}

		// Some other error, report it.
		return diag.Errorf("SUSPEND_PROCESSES:READ failed to read group: %s", err)
	}

	// If nothing was found, then return no state.
	spResponse := resp
	if spResponse == nil {
		resourceData.SetId("")
		return nil
	}

	if err := commons.SuspendProcessesResource.OnRead(spResponse.SuspendProcesses, resourceData, meta); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("===> Elastigroup read successfully: %s <===", resourceData.Id())
	return nil
}

func resourceSpotinstAWSSuspendProcessesCreate(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf(string(commons.ResourceOnCreate), commons.SuspendProcessesResource.GetName())

	suspendProcesses, err := commons.SuspendProcessesResource.OnCreate(resourceData, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	suspendProcessesId, err := createSuspendProcesses(resourceData, suspendProcesses, meta.(*Client))
	if err != nil {
		return diag.FromErr(err)
	}

	resourceData.SetId(spotinst.StringValue(suspendProcessesId))

	log.Printf("===> SuspendProcesses created successfully for Elastigroup: %s <===", resourceData.Id())

	return resourceSpotinstAWSSuspendProcessesRead(ctx, resourceData, meta)

}

func createSuspendProcesses(resourceData *schema.ResourceData, suspendProcesses *aws.SuspendProcesses, spotinstClient *Client) (*string, error) {
	if json, err := commons.ToJson(suspendProcesses); err != nil {
		return nil, err
	} else {
		log.Printf("===> SuspendProcesses create configuration: %s", json)
	}
	groupID := spotinst.String(resourceData.Get(string(elastigroup_aws_suspend_processes.GroupID)).(string))
	err := resource.RetryContext(context.Background(), time.Minute, func() *resource.RetryError {
		input := &aws.CreateSuspensionsInput{
			GroupID:     groupID,
			Suspensions: suspendProcesses.Suspensions,
		}
		_, err := spotinstClient.elastigroup.CloudProviderAWS().CreateSuspensions(context.Background(), input)
		if err != nil {
			// an error occurred, no retryable errors for this resource.
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("[ERROR] failed to create SuspendProcesses for Elastigroup: %s", err)
	}
	return groupID, nil

}

func resourceSpotinstAWSSuspendProcessesDelete(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	resourceId := resourceData.Id()
	log.Printf(string(commons.ResourceOnDelete), commons.SuspendProcessesResource.GetName(), resourceId)

	if err := deleteSuspendProcesses(resourceData, meta); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("===> SuspendProcesses deleted successfully for Elastigroup: %s <===", resourceData.Id())
	resourceData.SetId("")
	return nil
}

func deleteSuspendProcesses(resourceData *schema.ResourceData, meta interface{}) error {

	listInput := &aws.ListSuspensionsInput{}
	gID := resourceData.Id()
	listInput.GroupID = &gID

	curr, err := meta.(*Client).elastigroup.CloudProviderAWS().ListSuspensions(context.Background(), listInput)
	if err != nil {
		return fmt.Errorf("[ERROR] failed to update suspend processes: %v", err)
	}

	delInput := &aws.DeleteSuspensionsInput{}
	delInput.GroupID = &gID
	delInput.Processes = curr.SuspendProcesses.Processes

	if json, err := commons.ToJson(delInput); err != nil {
		return err
	} else {
		log.Printf("===> suspendProcesses delete configuration: %s", json)
	}

	if _, err := meta.(*Client).elastigroup.CloudProviderAWS().DeleteSuspensions(context.Background(), delInput); err != nil {
		return fmt.Errorf("[ERROR] onDelete() -> Failed to delete suspendProcesses for Elastigroup: %s", err)
	}
	return nil
}

func resourceSpotinstAWSSuspendProcessesUpdate(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	resourceId := resourceData.Id()
	log.Printf(string(commons.ResourceOnUpdate), commons.SuspendProcessesResource.GetName(), resourceId)

	shouldUpdate, suspendProcesses, err := commons.SuspendProcessesResource.OnUpdate(resourceData, meta)
	if err != nil {
		return diag.FromErr(err)
	}
	if shouldUpdate {
		if err := updateSuspendProcesses(suspendProcesses, resourceData, meta); err != nil {
			return diag.FromErr(err)
		}
	}
	log.Printf("===> SuspendProcesses updated successfully: %s <===", resourceId)
	return resourceSpotinstAWSSuspendProcessesRead(ctx, resourceData, meta)
}

func updateSuspendProcesses(suspendProcesses *aws.SuspendProcesses, resourceData *schema.ResourceData, meta interface{}) error {

	var input = &aws.SuspendProcesses{
		Suspensions: suspendProcesses.Suspensions,
		Processes:   nil,
	}

	var processesInput []string = nil
	for _, suspension := range input.Suspensions {
		processesInput = append(processesInput, *suspension.Name)
	}
	input.Processes = processesInput

	var processesToDelete []string = nil

	req := &aws.ListSuspensionsInput{}
	gID := resourceData.Id()
	req.GroupID = &gID

	curr, err := meta.(*Client).elastigroup.CloudProviderAWS().ListSuspensions(context.Background(), req)
	if err != nil {
		return fmt.Errorf("[ERROR] failed to update suspend processes: %v", err)
	}

	currProcesses := curr.SuspendProcesses.Processes

	for _, process := range currProcesses {
		if _, exists := findString(processesInput, process); !exists {
			processesToDelete = append(processesToDelete, process)
		}
	}

	var processesToCreate []string = nil

	for _, process := range processesInput {
		if _, exists := findString(currProcesses, process); !exists {
			processesToCreate = append(processesToCreate, process)
		}
	}

	if json, err := commons.ToJson(input); err != nil {
		return err
	} else {
		log.Printf("===> Suspend Processes update configuration: %s", json)
	}

	if processesToCreate != nil {
		createReqBody := &aws.CreateSuspensionsInput{Suspensions: suspendProcesses.Suspensions}
		groupIDInput := resourceData.Id()
		createReqBody.GroupID = &groupIDInput

		_, err = meta.(*Client).elastigroup.CloudProviderAWS().CreateSuspensions(context.Background(), createReqBody)
		if err != nil {
			return fmt.Errorf("[ERROR] failed to update suspend processes: %v", err)
		}
	}

	if processesToDelete != nil {
		deleteReqBody := &aws.DeleteSuspensionsInput{
			GroupID:   &gID,
			Processes: processesToDelete,
		}

		if _, err := meta.(*Client).elastigroup.CloudProviderAWS().DeleteSuspensions(context.Background(), deleteReqBody); err != nil {
			return fmt.Errorf("[ERROR] onDelete() -> Failed to delete suspendProcesses for Elastigroup: %s", err)
		}
	}

	return nil
}

func findString(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}
