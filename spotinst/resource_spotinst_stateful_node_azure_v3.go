package spotinst

import (
	"context"
	"fmt"
	"github.com/spotinst/terraform-provider-spotinst/spotinst/azure_v3/stateful_node_azure"
	"github.com/spotinst/terraform-provider-spotinst/spotinst/azure_v3/stateful_node_azure_extension"
	"github.com/spotinst/terraform-provider-spotinst/spotinst/azure_v3/stateful_node_azure_health"
	"github.com/spotinst/terraform-provider-spotinst/spotinst/azure_v3/stateful_node_azure_image"
	"github.com/spotinst/terraform-provider-spotinst/spotinst/azure_v3/stateful_node_azure_launch_spec"
	"github.com/spotinst/terraform-provider-spotinst/spotinst/azure_v3/stateful_node_azure_load_balancer"
	"github.com/spotinst/terraform-provider-spotinst/spotinst/azure_v3/stateful_node_azure_login"
	"github.com/spotinst/terraform-provider-spotinst/spotinst/azure_v3/stateful_node_azure_network"
	"github.com/spotinst/terraform-provider-spotinst/spotinst/azure_v3/stateful_node_azure_persistence"
	"github.com/spotinst/terraform-provider-spotinst/spotinst/azure_v3/stateful_node_azure_scheduling"
	"github.com/spotinst/terraform-provider-spotinst/spotinst/azure_v3/stateful_node_azure_secret"
	"github.com/spotinst/terraform-provider-spotinst/spotinst/azure_v3/stateful_node_azure_strategy"
	"github.com/spotinst/terraform-provider-spotinst/spotinst/azure_v3/stateful_node_azure_vm_sizes"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	v3 "github.com/spotinst/spotinst-sdk-go/service/stateful/providers/azure"
	"github.com/spotinst/spotinst-sdk-go/spotinst"
	"github.com/spotinst/spotinst-sdk-go/spotinst/client"
	"github.com/spotinst/terraform-provider-spotinst/spotinst/commons"
)

func resourceSpotinstStatefulNodeAzureV3() *schema.Resource {
	setupStatefulNodeAzureV3Resource()

	return &schema.Resource{
		Create: resourceSpotinstStatefulNodeAzureV3Create,
		Read:   resourceSpotinstStatefulNodeAzureV3Read,
		Update: resourceSpotinstStatefulNodeAzureV3Update,
		Delete: resourceSpotinstStatefulNodeAzureV3Delete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: commons.StatefulNodeAzureV3Resource.GetSchemaMap(),
	}
}

func setupStatefulNodeAzureV3Resource() {
	fieldsMap := make(map[commons.FieldName]*commons.GenericField)

	stateful_node_azure.Setup(fieldsMap)
	stateful_node_azure_strategy.Setup(fieldsMap)
	stateful_node_azure_launch_spec.Setup(fieldsMap)
	stateful_node_azure_image.Setup(fieldsMap)
	stateful_node_azure_network.Setup(fieldsMap)
	stateful_node_azure_login.Setup(fieldsMap)
	stateful_node_azure_load_balancer.Setup(fieldsMap)
	stateful_node_azure_extension.Setup(fieldsMap)
	stateful_node_azure_secret.Setup(fieldsMap)
	stateful_node_azure_vm_sizes.Setup(fieldsMap)
	stateful_node_azure_persistence.Setup(fieldsMap)
	stateful_node_azure_scheduling.Setup(fieldsMap)
	stateful_node_azure_health.Setup(fieldsMap)

	commons.StatefulNodeAzureV3Resource = commons.NewStatefulNodeAzureV3Resource(fieldsMap)
}

func resourceSpotinstStatefulNodeAzureV3Create(resourceData *schema.ResourceData, meta interface{}) error {
	log.Printf(string(commons.ResourceOnCreate),
		commons.StatefulNodeAzureV3Resource.GetName())

	statefulNode, err := commons.StatefulNodeAzureV3Resource.OnCreate(resourceData, meta)
	if err != nil {
		return err
	}

	if importVMConfig, ok := resourceData.GetOk(string(stateful_node_azure.ImportVM)); ok {

		importVMStatefulNodeInput, err := expandStatefulNodeAzureImportVMConfig(importVMConfig, statefulNode)
		if err != nil {
			return fmt.Errorf("stateful node/azure: failed expanding import vm configuration: %v", err)
		}

		statefulNodeId, err := createAzureV3StatefulNodeImportVM(importVMStatefulNodeInput, meta.(*Client))
		if err != nil {
			return err
		}

		resourceData.SetId(spotinst.StringValue(statefulNodeId))
		log.Printf("===> Stateful node using import vm created successfully: %s <===", resourceData.Id())

	} else {
		statefulNodeId, err := createAzureV3StatefulNode(statefulNode, meta.(*Client))
		if err != nil {
			return err
		}

		resourceData.SetId(spotinst.StringValue(statefulNodeId))
		log.Printf("===> Stateful node created successfully: %s <===", resourceData.Id())
	}

	return resourceSpotinstStatefulNodeAzureV3Read(resourceData, meta)
}

func expandStatefulNodeAzureImportVMConfig(data interface{}, statefulNode *v3.StatefulNode) (*v3.ImportVMStatefulNodeInput, error) {
	spec := &v3.ImportVMStatefulNodeInput{
		StatefulNodeImport: &v3.StatefulNodeImport{
			StatefulNode: statefulNode,
		},
	}

	list := data.([]interface{})
	if list != nil && list[0] != nil {
		m := list[0].(map[string]interface{})

		if v, ok := m[string(stateful_node_azure.ImportVMOriginalVMName)].(string); ok && v != "" {
			spec.StatefulNodeImport.OriginalVMName = spotinst.String(v)
		}

		if v, ok := m[string(stateful_node_azure.ImportVMResourceGroupName)].(string); ok && v != "" {
			spec.StatefulNodeImport.ResourceGroupName = spotinst.String(v)
		}

		if v, ok := m[string(stateful_node_azure.ImportVMDrainingTimeout)].(int); ok && v >= 0 {
			spec.StatefulNodeImport.DrainingTimeout = spotinst.Int(v)
		}

		if v, ok := m[string(stateful_node_azure.ImportVMResourcesRetentionTime)].(int); ok && v >= 0 {
			spec.StatefulNodeImport.ResourcesRetentionTime = spotinst.Int(v)
		}
	}

	return spec, nil
}

func createAzureV3StatefulNodeImportVM(importVMStatefulNodeInput *v3.ImportVMStatefulNodeInput, spotinstClient *Client) (*string, error) {
	if json, err := commons.ToJson(importVMStatefulNodeInput); err != nil {
		return nil, err
	} else {
		log.Printf("===> Stateful node import vm create configuration: %s", json)
	}

	var resp *v3.ImportVMStatefulNodeOutput = nil
	err := resource.Retry(time.Minute, func() *resource.RetryError {
		r, err := spotinstClient.statefulNode.CloudProviderAzure().ImportVM(context.Background(), importVMStatefulNodeInput)
		if err != nil {
			log.Printf("error: %v", err)
			// Some other error, report it.
			return resource.NonRetryableError(err)
		}
		resp = r
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("[ERROR] failed to create stateful node using import vm: %s", err)
	}
	return resp.StatefulNodeImport.StatefulNode.ID, nil
}

func createAzureV3StatefulNode(statefulNode *v3.StatefulNode, spotinstClient *Client) (*string, error) {
	if json, err := commons.ToJson(statefulNode); err != nil {
		return nil, err
	} else {
		log.Printf("===> Stateful node create configuration: %s", json)
	}

	var resp *v3.CreateStatefulNodeOutput = nil
	err := resource.Retry(time.Minute, func() *resource.RetryError {
		input := &v3.CreateStatefulNodeInput{StatefulNode: statefulNode}
		r, err := spotinstClient.statefulNode.CloudProviderAzure().Create(context.Background(), input)
		if err != nil {
			log.Printf("error: %v", err)
			// Some other error, report it.
			return resource.NonRetryableError(err)
		}
		resp = r
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("[ERROR] failed to create stateful node: %s", err)
	}
	return resp.StatefulNode.ID, nil
}

func resourceSpotinstStatefulNodeAzureV3Read(resourceData *schema.ResourceData, meta interface{}) error {
	id := resourceData.Id()
	log.Printf(string(commons.ResourceFieldOnRead),
		commons.StatefulNodeAzureV3Resource.GetName(), id)

	input := &v3.ReadStatefulNodeInput{ID: spotinst.String(id)}
	resp, err := meta.(*Client).statefulNode.CloudProviderAzure().Read(context.Background(), input)
	if err != nil {
		// If the stateful node was not found, return nil so that we can show
		// that the stateful node does not exist
		if errs, ok := err.(client.Errors); ok && len(errs) > 0 {
			for _, err := range errs {
				if err.Code == ErrCodeGroupNotFound {
					resourceData.SetId("")
					return nil
				}
			}
		}

		// Some other error, report it.
		return fmt.Errorf("failed to read stateful node: %s", err)
	}

	// If nothing was found, then return no state.
	statefulNodeResponse := resp.StatefulNode
	if statefulNodeResponse == nil {
		resourceData.SetId("")
		return nil
	}

	if err := commons.StatefulNodeAzureV3Resource.OnRead(statefulNodeResponse, resourceData, meta); err != nil {
		return err
	}
	log.Printf("===> Stateful node read successfully: %s <===", id)
	return nil
}

func resourceSpotinstStatefulNodeAzureV3Update(resourceData *schema.ResourceData, meta interface{}) error {
	id := resourceData.Id()
	log.Printf(string(commons.ResourceOnUpdate),
		commons.StatefulNodeAzureV3Resource.GetName(), id)

	shouldUpdate, statefulNode, err := commons.StatefulNodeAzureV3Resource.OnUpdate(resourceData, meta)
	if err != nil {
		return err
	}

	if shouldUpdate {
		statefulNode.SetID(spotinst.String(id))
		if err := updateAzureV3StatefulNode(statefulNode, resourceData, meta); err != nil {
			return err
		}
	}

	log.Printf("===> Stateful node updated successfully: %s <===", id)
	return resourceSpotinstStatefulNodeAzureV3Read(resourceData, meta)
}

func updateAzureV3StatefulNode(statefulNode *v3.StatefulNode, resourceData *schema.ResourceData, meta interface{}) error {
	var input = &v3.UpdateStatefulNodeInput{
		StatefulNode: statefulNode,
	}

	statefulNodeId := resourceData.Id()
	var shouldUpdateState = false
	var shouldDetachDataDisk = false
	var shouldAttachDataDisk = false
	if updateState, ok := resourceData.GetOk(string(stateful_node_azure.UpdateState)); ok {
		list := updateState.([]interface{})
		if len(list) > 0 && list[0] != nil {
			shouldUpdateState = true
		}
	}

	if attachDataDisk, ok := resourceData.GetOk(string(stateful_node_azure.AttachDataDisk)); ok {
		list := attachDataDisk.([]interface{})
		if len(list) > 0 && list[0] != nil {
			shouldAttachDataDisk = true
		}
	}

	if detachDataDisk, ok := resourceData.GetOk(string(stateful_node_azure.DetachDataDisk)); ok {
		list := detachDataDisk.([]interface{})
		if len(list) > 0 && list[0] != nil {
			shouldDetachDataDisk = true
		}
	}

	if json, err := commons.ToJson(statefulNode); err != nil {
		return err
	} else {
		log.Printf("===> Stateful node update configuration: %s", json)
	}

	if _, err := meta.(*Client).statefulNode.CloudProviderAzure().Update(context.Background(), input); err != nil {
		return fmt.Errorf("[ERROR] Failed to update stateful node [%v]: %v", statefulNodeId, err)
	} else if shouldUpdateState {
		if err := updateStateAzureV3StatefulNode(resourceData, meta); err != nil {
			log.Printf("[ERROR] Stateful node [%v] state update failed, error: %v", statefulNodeId, err)
			return err
		}
	} else {
		log.Printf("onUpdate() -> Field [%v] is missing, skipping state update for stateful node",
			string(stateful_node_azure.UpdateState))
	}

	if shouldAttachDataDisk {
		if err := attachDataDiskAzureV3StatefulNode(resourceData, meta); err != nil {
			log.Printf("[ERROR] Stateful node [%v] attach data disk failed, error: %v", statefulNodeId, err)
			return err
		}
	} else {
		log.Printf("onUpdate() -> Field [%v] is missing, skipping attach data disk for stateful node",
			string(stateful_node_azure.AttachDataDisk))
	}

	if shouldDetachDataDisk {
		if err := detachDataDiskAzureV3StatefulNode(resourceData, meta); err != nil {
			log.Printf("[ERROR] Stateful node [%v] detach data disk failed, error: %v", statefulNodeId, err)
			return err
		}
	} else {
		log.Printf("onUpdate() -> Field [%v] is missing, skipping detach data disk for stateful node",
			string(stateful_node_azure.DetachDataDisk))
	}

	return nil
}

func updateStateAzureV3StatefulNode(resourceData *schema.ResourceData, meta interface{}) error {
	statefulNodeID := resourceData.Id()

	updateState, ok := resourceData.GetOk(string(stateful_node_azure.UpdateState))
	if !ok {
		return fmt.Errorf("stateful node/azure: missing update_state for stateful node %q", statefulNodeID)
	}

	list := updateState.([]interface{})
	if len(list) > 0 && list[0] != nil {
		updateStatefulNodeStateSchema := list[0].(map[string]interface{})
		if updateStatefulNodeStateSchema == nil {
			return fmt.Errorf("stateful node/azure: missing update state configuration, "+
				"skipping update state for stateful node %q", statefulNodeID)
		}

		updateStateSpec, err := expandStatefulNodeAzureUpdateStateConfig(updateStatefulNodeStateSchema, statefulNodeID)
		if err != nil {
			return fmt.Errorf("stateful node/azure: failed expanding state update "+
				"configuration for stateful node %q, error: %v", statefulNodeID, err)
		}

		updateStateJSON, err := commons.ToJson(updateStatefulNodeStateSchema)
		if err != nil {
			return fmt.Errorf("stateful node/azure: failed marshaling state update "+
				"configuration for stateful node %q, error: %v", statefulNodeID, err)
		}

		log.Printf("onUpdate() -> Updating stateful node [%v] with configuration %s", statefulNodeID, updateStateJSON)
		updateStateInput := &v3.UpdateStatefulNodeStateInput{ID: updateStateSpec.ID,
			StatefulNodeState: updateStateSpec.StatefulNodeState}
		if _, err = meta.(*Client).statefulNode.CloudProviderAzure().UpdateState(context.TODO(),
			updateStateInput); err != nil {
			return fmt.Errorf("onUpdate() -> State update failed for stateful node [%v], error: %v",
				statefulNodeID, err)
		}
		log.Printf("onUpdate() -> Successfully updated state for stateful node [%v]", statefulNodeID)
	}

	return nil
}

func attachDataDiskAzureV3StatefulNode(resourceData *schema.ResourceData, meta interface{}) error {
	statefulNodeID := resourceData.Id()

	attachDataDisk, ok := resourceData.GetOk(string(stateful_node_azure.AttachDataDisk))
	if !ok {
		return fmt.Errorf("stateful node/azure: missing attach_data_disk for stateful node %q", statefulNodeID)
	}

	list := attachDataDisk.([]interface{})
	if len(list) > 0 && list[0] != nil {
		attachDataDiskStatefulNodeSchema := list[0].(map[string]interface{})
		if attachDataDiskStatefulNodeSchema == nil {
			return fmt.Errorf("stateful node/azure: missing attach data disk configuration, "+
				"skipping attach data disk for stateful node %q", statefulNodeID)
		}

		attachDataDiskSpec, err := expandStatefulNodeAzureAttachDataDiskConfig(attachDataDiskStatefulNodeSchema, statefulNodeID)
		if err != nil {
			return fmt.Errorf("stateful node/azure: failed expanding attach data disk "+
				"configuration for stateful node %q, error: %v", statefulNodeID, err)
		}

		updateStateJSON, err := commons.ToJson(attachDataDiskStatefulNodeSchema)
		if err != nil {
			return fmt.Errorf("stateful node/azure: failed marshaling attach data disk "+
				"configuration for stateful node %q, error: %v", statefulNodeID, err)
		}

		log.Printf("onUpdate() -> Updating stateful node [%v] with configuration %s", statefulNodeID, updateStateJSON)
		attachDataDiskInput := &v3.AttachStatefulNodeDataDiskInput{ID: attachDataDiskSpec.ID,
			DataDiskName:              attachDataDiskSpec.DataDiskName,
			DataDiskResourceGroupName: attachDataDiskSpec.DataDiskResourceGroupName,
			StorageAccountType:        attachDataDiskSpec.StorageAccountType, SizeGB: attachDataDiskSpec.SizeGB,
			LUN: attachDataDiskSpec.LUN, Zone: attachDataDiskSpec.Zone}
		if _, err = meta.(*Client).statefulNode.CloudProviderAzure().AttachDataDisk(context.TODO(),
			attachDataDiskInput); err != nil {
			return fmt.Errorf("onUpdate() -> Attach data disk failed for stateful node [%v], error: %v",
				statefulNodeID, err)
		}
		log.Printf("onUpdate() -> Successfully attached data disk for stateful node [%v]", statefulNodeID)
	}

	return nil
}

func detachDataDiskAzureV3StatefulNode(resourceData *schema.ResourceData, meta interface{}) error {
	statefulNodeID := resourceData.Id()

	detachDataDisk, ok := resourceData.GetOk(string(stateful_node_azure.DetachDataDisk))
	if !ok {
		return fmt.Errorf("stateful node/azure: missing detach_data_disk for stateful node %q", statefulNodeID)
	}

	list := detachDataDisk.([]interface{})
	if len(list) > 0 && list[0] != nil {
		detachDataDiskStatefulNodeSchema := list[0].(map[string]interface{})
		if detachDataDiskStatefulNodeSchema == nil {
			return fmt.Errorf("stateful node/azure: missing detach data disk configuration, "+
				"skipping detach data disk for stateful node %q", statefulNodeID)
		}

		detachDataDiskSpec, err := expandStatefulNodeAzureDetachDataDiskConfig(detachDataDiskStatefulNodeSchema, statefulNodeID)
		if err != nil {
			return fmt.Errorf("stateful node/azure: failed expanding detach data disk "+
				"configuration for stateful node %q, error: %v", statefulNodeID, err)
		}

		updateStateJSON, err := commons.ToJson(detachDataDiskStatefulNodeSchema)
		if err != nil {
			return fmt.Errorf("stateful node/azure: failed marshaling detach data disk "+
				"configuration for stateful node %q, error: %v", statefulNodeID, err)
		}

		log.Printf("onUpdate() -> Updating stateful node [%v] with configuration %s", statefulNodeID, updateStateJSON)
		detachDataDiskInput := &v3.DetachStatefulNodeDataDiskInput{ID: detachDataDiskSpec.ID,
			DataDiskName:              detachDataDiskSpec.DataDiskName,
			DataDiskResourceGroupName: detachDataDiskSpec.DataDiskResourceGroupName,
			ShouldDeallocate:          detachDataDiskSpec.ShouldDeallocate}
		if _, err = meta.(*Client).statefulNode.CloudProviderAzure().DetachDataDisk(context.TODO(),
			detachDataDiskInput); err != nil {
			return fmt.Errorf("onUpdate() -> detach data disk failed for stateful node [%v], error: %v",
				statefulNodeID, err)
		}
		log.Printf("onUpdate() -> Successfully detached data disk for stateful node [%v]", statefulNodeID)
	}

	return nil
}

func expandStatefulNodeAzureUpdateStateConfig(data interface{}, statefulNodeID string) (*v3.UpdateStatefulNodeStateInput, error) {
	spec := &v3.UpdateStatefulNodeStateInput{
		ID: spotinst.String(statefulNodeID),
	}

	if data != nil {
		m := data.(map[string]interface{})

		if v, ok := m[string(stateful_node_azure.State)].(string); ok && v != "" {
			spec.StatefulNodeState = spotinst.String(v)
		}
	}

	return spec, nil
}

func expandStatefulNodeAzureAttachDataDiskConfig(data interface{},
	statefulNodeID string) (*v3.AttachStatefulNodeDataDiskInput, error) {
	spec := &v3.AttachStatefulNodeDataDiskInput{
		ID: spotinst.String(statefulNodeID),
	}

	if data != nil {
		m := data.(map[string]interface{})

		if v, ok := m[string(stateful_node_azure.AttachDataDiskName)].(string); ok && v != "" {
			spec.DataDiskName = spotinst.String(v)
		}

		if v, ok := m[string(stateful_node_azure.AttachDataDiskResourceGroupName)].(string); ok && v != "" {
			spec.DataDiskResourceGroupName = spotinst.String(v)
		}

		if v, ok := m[string(stateful_node_azure.AttachStorageAccountType)].(string); ok && v != "" {
			spec.StorageAccountType = spotinst.String(v)
		}

		if v, ok := m[string(stateful_node_azure.AttachSizeGB)].(int); ok && v > 0 {
			spec.SizeGB = spotinst.Int(v)
		}

		if v, ok := m[string(stateful_node_azure.AttachLUN)].(int); ok && v >= 0 {
			spec.LUN = spotinst.Int(v)
		}

		if v, ok := m[string(stateful_node_azure.AttachZone)].(string); ok && v != "" {
			spec.Zone = spotinst.String(v)
		}
	}

	return spec, nil
}

func expandStatefulNodeAzureDetachDataDiskConfig(data interface{}, statefulNodeID string) (*v3.DetachStatefulNodeDataDiskInput, error) {
	spec := &v3.DetachStatefulNodeDataDiskInput{
		ID: spotinst.String(statefulNodeID),
	}

	if data != nil {
		m := data.(map[string]interface{})

		if v, ok := m[string(stateful_node_azure.DetachDataDiskName)].(string); ok && v != "" {
			spec.DataDiskName = spotinst.String(v)
		}

		if v, ok := m[string(stateful_node_azure.DetachDataDiskResourceGroupName)].(string); ok && v != "" {
			spec.DataDiskResourceGroupName = spotinst.String(v)
		}

		if v, ok := m[string(stateful_node_azure.DetachShouldDeallocate)].(bool); ok {
			spec.ShouldDeallocate = spotinst.Bool(v)
		}

		if v, ok := m[string(stateful_node_azure.DetachTTLInHours)].(int); ok && v > 0 {
			spec.TTLInHours = spotinst.Int(v)
		}

	}

	return spec, nil
}

func resourceSpotinstStatefulNodeAzureV3Delete(resourceData *schema.ResourceData, meta interface{}) error {
	id := resourceData.Id()
	log.Printf(string(commons.ResourceOnDelete),
		commons.StatefulNodeAzureV3Resource.GetName(), id)

	if err := deleteAzureV3StatefulNode(resourceData, meta); err != nil {
		return err
	}

	log.Printf("===> Stateful node deleted successfully: %s <===", resourceData.Id())
	resourceData.SetId("")
	return nil
}

func deleteAzureV3StatefulNode(resourceData *schema.ResourceData, meta interface{}) error {
	statefulNodeId := resourceData.Id()
	input := &v3.DeleteStatefulNodeInput{
		ID: spotinst.String(statefulNodeId),
		DeallocationConfig: &v3.DeallocationConfig{
			ShouldTerminateVM: spotinst.Bool(resourceData.Get(string(stateful_node_azure.ShouldTerminateVm)).(bool)),
			NetworkDeallocationConfig: &v3.ResourceDeallocationConfig{
				ShouldDeallocate: spotinst.Bool(resourceData.Get(string(stateful_node_azure.NetworkShouldDeallocate)).(bool)),
				TTLInHours:       spotinst.Int(resourceData.Get(string(stateful_node_azure.NetworkTTLInHours)).(int)),
			},
			DiskDeallocationConfig: &v3.ResourceDeallocationConfig{
				ShouldDeallocate: spotinst.Bool(resourceData.Get(string(stateful_node_azure.DiskShouldDeallocate)).(bool)),
				TTLInHours:       spotinst.Int(resourceData.Get(string(stateful_node_azure.DiskTTLInHours)).(int)),
			},
			SnapshotDeallocationConfig: &v3.ResourceDeallocationConfig{
				ShouldDeallocate: spotinst.Bool(resourceData.Get(string(stateful_node_azure.SnapshotShouldDeallocate)).(bool)),
				TTLInHours:       spotinst.Int(resourceData.Get(string(stateful_node_azure.SnapshotTTLInHours)).(int)),
			},
			PublicIPDeallocationConfig: &v3.ResourceDeallocationConfig{
				ShouldDeallocate: spotinst.Bool(resourceData.Get(string(stateful_node_azure.PublicIPShouldDeallocate)).(bool)),
				TTLInHours:       spotinst.Int(resourceData.Get(string(stateful_node_azure.PublicIPTTLInHours)).(int)),
			},
		},
	}

	if json, err := commons.ToJson(input); err != nil {
		return err
	} else {
		log.Printf("===> Stateful node delete configuration: %s", json)
	}

	if _, err := meta.(*Client).statefulNode.CloudProviderAzure().Delete(context.Background(), input); err != nil {
		return fmt.Errorf("[ERROR] onDelete() -> Failed to delete stateful node: %s", err)
	}
	return nil
}
