package spotinst

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	azurev3 "github.com/spotinst/spotinst-sdk-go/service/stateful/providers/azure"
	"github.com/spotinst/spotinst-sdk-go/spotinst"
	"github.com/spotinst/terraform-provider-spotinst/spotinst/commons"
	"log"
	"testing"
)

func createStatefulNodeAzureV3ResourceName(name string) string {
	return fmt.Sprintf("%v.%v", string(commons.StatefulNodeAzureV3ResourceName), name)
}

func testStatefulNodeAzureV3Destroy(s *terraform.State) error {
	client := testAccProviderAzure.Meta().(*Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != string(commons.StatefulNodeAzureV3ResourceName) {
			continue
		}
		input := &azurev3.GetStatefulNodeStateInput{ID: spotinst.String(rs.Primary.ID)}
		resp, err := client.statefulNode.CloudProviderAzure().GetState(context.Background(), input)
		if err == nil && resp != nil && resp.StatefulNodeState != nil {
			statefulNodeState := spotinst.StringValue(resp.StatefulNodeState.Status)
			if statefulNodeState != "DEALLOCATE" && statefulNodeState != "DEALLOCATING" {
				return fmt.Errorf("stateful node still exists! stateful node state = %s", statefulNodeState)
			}
		}
	}
	return nil
}

func testCheckStatefulNodeAzureV3Attributes(statefulNode *azurev3.StatefulNode, expectedName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if spotinst.StringValue(statefulNode.Name) != expectedName {
			return fmt.Errorf("bad content: %v", statefulNode.Name)
		}
		return nil
	}
}

func testCheckStatefulNodeAzureV3Exists(statefulNode *azurev3.StatefulNode, resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no resource ID is set")
		}
		client := testAccProviderAzure.Meta().(*Client)
		input := &azurev3.ReadStatefulNodeInput{ID: spotinst.String(rs.Primary.ID)}
		resp, err := client.statefulNode.CloudProviderAzure().Read(context.Background(), input)
		if err != nil {
			return err
		}
		*statefulNode = *resp.StatefulNode
		return nil
	}
}

type AzureV3StatefulNodeConfigMetadata struct {
	statefulNodeName     string
	acdIdentifier        string
	controllerClusterID  string
	provider             string
	strategy             string
	autoScaler           string
	health               string
	vmSizes              string
	osDisk               string
	dataDisk             string
	image                string
	network              string
	extensions           string
	login                string
	persistence          string
	secret               string
	scheduling           string
	tag                  string
	bootDiagnostics      string
	variables            string
	updateBaselineFields bool
}

func createStatefulNodeAzureV3Terraform(StatefulNodeMeta *AzureV3StatefulNodeConfigMetadata) string {
	if StatefulNodeMeta == nil {
		return ""
	}

	if StatefulNodeMeta.provider == "" {
		StatefulNodeMeta.provider = "azure"
	}

	//TODO - need a launchSpecification test?

	if StatefulNodeMeta.strategy == "" {
		StatefulNodeMeta.strategy = testStrategyStatefulNodeAzureV3Config_Create
	}

	if StatefulNodeMeta.health == "" {
		StatefulNodeMeta.health = testHealthStatefulNodeAzureV3Config_Create
	}

	if StatefulNodeMeta.vmSizes == "" {
		StatefulNodeMeta.vmSizes = testVMSizesStatefulNodeAzureV3Config_Create
	}
	//
	if StatefulNodeMeta.image == "" {
		StatefulNodeMeta.image = testImageStatefulNodeAzureV3Config_Create
	}

	if StatefulNodeMeta.network == "" {
		StatefulNodeMeta.network = testNetworkStatefulNodeAzureV3Config_Create
	}

	if StatefulNodeMeta.osDisk == "" {
		StatefulNodeMeta.osDisk = testOSDiskStatefulNodeAzureV3Config_Create
	}

	if StatefulNodeMeta.dataDisk == "" {
		StatefulNodeMeta.dataDisk = testDataDiskStatefulNodeAzureV3Config_Create
	}

	if StatefulNodeMeta.extensions == "" {
		StatefulNodeMeta.extensions = testExtensionsStatefulNodeAzureV3Config_Create
	}

	if StatefulNodeMeta.login == "" {
		StatefulNodeMeta.login = testAzureV3LoginStatefulNodeConfig_Create
	}

	if StatefulNodeMeta.persistence == "" {
		StatefulNodeMeta.persistence = testPersistenceStatefulNodeAzureV3Config_Create
	}
	//
	if StatefulNodeMeta.scheduling == "" {
		StatefulNodeMeta.scheduling = testSchedulingStatefulNodeAzureV3Config_Create
	}

	//if StatefulNodeMeta.secret == "" {
	//	StatefulNodeMeta.secret = testSecretsStatefulNodeAzureV3Config_Create
	//}

	if StatefulNodeMeta.tag == "" {
		StatefulNodeMeta.tag = testTagStatefulNodeAzureV3Config_Create
	}

	if StatefulNodeMeta.bootDiagnostics == "" {
		StatefulNodeMeta.bootDiagnostics = testBootDiagnosticsStatefulNodeAzureV3Config_Create
	}

	template := `
	provider "azure" {
	token = "fake"
	account = "fake"
	}
	`
	if StatefulNodeMeta.updateBaselineFields {
		format := testBaselineStatefulNodeAzureV3Config_Update
		template += fmt.Sprintf(format,
			StatefulNodeMeta.statefulNodeName,
			StatefulNodeMeta.provider,
			StatefulNodeMeta.statefulNodeName,
			StatefulNodeMeta.login,
			StatefulNodeMeta.strategy,
			StatefulNodeMeta.image,
			StatefulNodeMeta.extensions,
			StatefulNodeMeta.network,
			StatefulNodeMeta.osDisk,
			StatefulNodeMeta.dataDisk,
			StatefulNodeMeta.health,
			StatefulNodeMeta.vmSizes,
			StatefulNodeMeta.persistence,
			StatefulNodeMeta.scheduling,
			//StatefulNodeMeta.secret,
			StatefulNodeMeta.tag,
			StatefulNodeMeta.bootDiagnostics,
		)
	} else {
		format := testBaselineStatefulNodeAzureV3Config_Create
		template += fmt.Sprintf(format,
			StatefulNodeMeta.statefulNodeName,
			StatefulNodeMeta.provider,
			StatefulNodeMeta.statefulNodeName,
			StatefulNodeMeta.login,
			StatefulNodeMeta.strategy,
			StatefulNodeMeta.image,
			StatefulNodeMeta.extensions,
			StatefulNodeMeta.network,
			StatefulNodeMeta.osDisk,
			StatefulNodeMeta.dataDisk,
			StatefulNodeMeta.health,
			StatefulNodeMeta.vmSizes,
			StatefulNodeMeta.persistence,
			StatefulNodeMeta.scheduling,
			//StatefulNodeMeta.secret,
			StatefulNodeMeta.tag,
			StatefulNodeMeta.bootDiagnostics,
		)

	}

	if StatefulNodeMeta.variables != "" {
		template = StatefulNodeMeta.variables + "\n" + template
	}

	log.Printf("Terraform [%v] template:\n%v", StatefulNodeMeta.statefulNodeName, template)
	return template
}

// region Stateful Node Azure: Baseline
func TestAccSpotinstStatefulNodeAzureV3_Baseline(t *testing.T) {
	statefulNodeName := "terraform-tests-do-not-delete"
	resourceName := createStatefulNodeAzureV3ResourceName(statefulNodeName)

	var node azurev3.StatefulNode
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t, "azure") },
		Providers:    TestAccProviders,
		CheckDestroy: testStatefulNodeAzureV3Destroy,

		Steps: []resource.TestStep{
			{
				Config: createStatefulNodeAzureV3Terraform(&AzureV3StatefulNodeConfigMetadata{statefulNodeName: statefulNodeName}),
				Check: resource.ComposeTestCheckFunc(
					testCheckStatefulNodeAzureV3Exists(&node, resourceName),
					testCheckStatefulNodeAzureV3Attributes(&node, statefulNodeName),
					resource.TestCheckResourceAttr(resourceName, "region", "eastus"),
					resource.TestCheckResourceAttr(resourceName, "description", "tamir-test-file-1"),
					resource.TestCheckResourceAttr(resourceName, "os", "Linux"),
				),
			},
			{
				Config: createStatefulNodeAzureV3Terraform(&AzureV3StatefulNodeConfigMetadata{
					statefulNodeName:     statefulNodeName,
					updateBaselineFields: true,
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "region", "eastus"),
					resource.TestCheckResourceAttr(resourceName, "description", "tamir-test-file-1-updated"),
					resource.TestCheckResourceAttr(resourceName, "os", "Linux"),
				),
			},
		},
	})
}

//TODO - resource group name confidential?
const testBaselineStatefulNodeAzureV3Config_Create = `
resource "` + string(commons.StatefulNodeAzureV3ResourceName) + `" "%v" {
provider = "%v"
name = "%v"
os = "Linux"
region = "eastus"
description = "tamir-test-file-1"
resource_group_name = "CoreReliabilityResourceGroup" 
%v
%v
%v
%v
%v
%v
%v
%v
%v
%v
%v
%v
%v

delete {
	should_terminate_vm = true
	network_should_deallocate = true
	network_ttl_in_hours      = 0
	disk_should_deallocate = true
	disk_ttl_in_hours      = 0
	snapshot_should_deallocate = true
	snapshot_ttl_in_hours      = 0
	public_ip_should_deallocate = true
	public_ip_ttl_in_hours      = 0
	}
}
`

//TODO - resource group name confidential?
const testBaselineStatefulNodeAzureV3Config_Update = `
resource "` + string(commons.StatefulNodeAzureV3ResourceName) + `" "%v" {
provider = "%v"
name = "%v"
os = "Linux"
region = "eastus"
description = "tamir-test-file-1-updated"
resource_group_name = "CoreReliabilityResourceGroup" 
%v
%v
%v
%v
%v
%v
%v
%v
%v
%v
%v
%v
%v

delete {
	should_terminate_vm = true
	network_should_deallocate = true
	network_ttl_in_hours      = 0
	disk_should_deallocate = true
	disk_ttl_in_hours      = 0
	snapshot_should_deallocate = true
	snapshot_ttl_in_hours      = 0
	public_ip_should_deallocate = true
	public_ip_ttl_in_hours      = 0
	}
}
`

//endregion

// region Stateful Node Azure : Login
func TestAccSpotinstStatefulNodeAzureV3_Login(t *testing.T) {
	statefulNodeName := "terraform-tests-do-not-delete"
	resourceName := createStatefulNodeAzureV3ResourceName(statefulNodeName)

	var node azurev3.StatefulNode
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t, "azure") },
		Providers:    TestAccProviders,
		CheckDestroy: testStatefulNodeAzureV3Destroy,

		Steps: []resource.TestStep{
			{
				Config: createStatefulNodeAzureV3Terraform(&AzureV3StatefulNodeConfigMetadata{statefulNodeName: statefulNodeName}),
				Check: resource.ComposeTestCheckFunc(
					testCheckStatefulNodeAzureV3Exists(&node, resourceName),
					testCheckStatefulNodeAzureV3Attributes(&node, statefulNodeName),
					resource.TestCheckResourceAttr(resourceName, "login.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "login.0.user_name", "tamir"),
					resource.TestCheckResourceAttr(resourceName, "login.0.ssh_public_key", "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDfWrinLRVHx+KB57pb1mEYBueGfPzyVa2qPpCPZYbpcuL45nDKU2B14twX91+/cJ2m7DmUa8LLk2EVwBW8FBTfg5Fuwj8+kTnk4PMo4G+T0UgFt7NuD47I5fxg3sD9WQFUbXlO44Flp+k5MHlv+hF8iHz/QRz2QDDKxPGLWM1mh10LtLz4T+im/73RviTgbJhCZQr0+Yx7Uz1ZlWkrPThLUa9/4Br5mKLk3zEYa8mbg4LblJXIgknFsZ3cXlqtN5WofxJEDLy9QiKMxDJ2PZfR73IscpWtPnAMZjcTf6aI02FKAg+iEs0mdh3bGVGLxNi5w32lWOiiqKKJGKa1ctWb"),
				),
			},
			{
				Config: createStatefulNodeAzureV3Terraform(&AzureV3StatefulNodeConfigMetadata{
					statefulNodeName:     statefulNodeName,
					login:                testAzureV3LoginStatefulNodeConfig_Update,
					updateBaselineFields: true,
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "login.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "login.0.user_name", "tamir"),
					resource.TestCheckResourceAttr(resourceName, "login.0.ssh_public_key", "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDfWrinLRVHx+KB57pb1mEYBueGfPzyVa2qPpCPZYbpcuL45nDKU2B14twX91+/cJ2m7DmUa8LLk2EVwBW8FBTfg5Fuwj8+kTnk4PMo4G+T0UgFt7NuD47I5fxg3sD9WQFUbXlO44Flp+k5MHlv+hF8iHz/QRz2QDDKxPGLWM1mh10LtLz4T+im/73RviTgbJhCZQr0+Yx7Uz1ZlWkrPThLUa9/4Br5mKLk3zEYa8mbg4LblJXIgknFsZ3cXlqtN5WofxJEDLy9QiKMxDJ2PZfR73IscpWtPnAMZjcTf6aI02FKAg+iEs0mdh3bGVGLxNi5w32lWOiiqKKJGKa1ctWb automation"),
				),
			},
		},
	})
}

const testAzureV3LoginStatefulNodeConfig_Create = `
login {
	user_name = "tamir"
	ssh_public_key  = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDfWrinLRVHx+KB57pb1mEYBueGfPzyVa2qPpCPZYbpcuL45nDKU2B14twX91+/cJ2m7DmUa8LLk2EVwBW8FBTfg5Fuwj8+kTnk4PMo4G+T0UgFt7NuD47I5fxg3sD9WQFUbXlO44Flp+k5MHlv+hF8iHz/QRz2QDDKxPGLWM1mh10LtLz4T+im/73RviTgbJhCZQr0+Yx7Uz1ZlWkrPThLUa9/4Br5mKLk3zEYa8mbg4LblJXIgknFsZ3cXlqtN5WofxJEDLy9QiKMxDJ2PZfR73IscpWtPnAMZjcTf6aI02FKAg+iEs0mdh3bGVGLxNi5w32lWOiiqKKJGKa1ctWb"
}
`

const testAzureV3LoginStatefulNodeConfig_Update = `
login {
	user_name = "tamir"
	ssh_public_key  = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDfWrinLRVHx+KB57pb1mEYBueGfPzyVa2qPpCPZYbpcuL45nDKU2B14twX91+/cJ2m7DmUa8LLk2EVwBW8FBTfg5Fuwj8+kTnk4PMo4G+T0UgFt7NuD47I5fxg3sD9WQFUbXlO44Flp+k5MHlv+hF8iHz/QRz2QDDKxPGLWM1mh10LtLz4T+im/73RviTgbJhCZQr0+Yx7Uz1ZlWkrPThLUa9/4Br5mKLk3zEYa8mbg4LblJXIgknFsZ3cXlqtN5WofxJEDLy9QiKMxDJ2PZfR73IscpWtPnAMZjcTf6aI02FKAg+iEs0mdh3bGVGLxNi5w32lWOiiqKKJGKa1ctWb automation"
}
`

//endregion

// region Stateful Node Azure : Persistence
func TestAccSpotinstStatefulNodeAzureV3_Persistence(t *testing.T) {
	statefulNodeName := "terraform-tests-do-not-delete"
	resourceName := createStatefulNodeAzureV3ResourceName(statefulNodeName)

	var node azurev3.StatefulNode
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t, "azure") },
		Providers:    TestAccProviders,
		CheckDestroy: testStatefulNodeAzureV3Destroy,

		Steps: []resource.TestStep{
			{
				Config: createStatefulNodeAzureV3Terraform(&AzureV3StatefulNodeConfigMetadata{statefulNodeName: statefulNodeName}),
				Check: resource.ComposeTestCheckFunc(
					testCheckStatefulNodeAzureV3Exists(&node, resourceName),
					testCheckStatefulNodeAzureV3Attributes(&node, statefulNodeName),
					resource.TestCheckResourceAttr(resourceName, "should_persist_os_disk", "false"),
					resource.TestCheckResourceAttr(resourceName, "os_disk_persistence_mode", "reattach"),
					resource.TestCheckResourceAttr(resourceName, "should_persist_data_disks", "false"),
					resource.TestCheckResourceAttr(resourceName, "data_disks_persistence_mode", "reattach"),
					resource.TestCheckResourceAttr(resourceName, "should_persist_network", "true"),
					resource.TestCheckResourceAttr(resourceName, "should_persist_vm", "false"),
				),
			},
			{
				Config: createStatefulNodeAzureV3Terraform(&AzureV3StatefulNodeConfigMetadata{
					statefulNodeName:     statefulNodeName,
					persistence:          testPersistenceStatefulNodeAzureV3Config_Update,
					updateBaselineFields: true,
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "should_persist_os_disk", "false"),
					resource.TestCheckResourceAttr(resourceName, "os_disk_persistence_mode", "reattach"),
					resource.TestCheckResourceAttr(resourceName, "should_persist_data_disks", "false"),
					resource.TestCheckResourceAttr(resourceName, "data_disks_persistence_mode", "onLaunch"),
					resource.TestCheckResourceAttr(resourceName, "should_persist_network", "true"),
					resource.TestCheckResourceAttr(resourceName, "should_persist_vm", "false"),
				),
			},
		},
	})
}

const testPersistenceStatefulNodeAzureV3Config_Create = `
should_persist_os_disk = false
os_disk_persistence_mode = "reattach"
should_persist_data_disks = false
data_disks_persistence_mode = "reattach"
should_persist_network = true
should_persist_vm = false
`

const testPersistenceStatefulNodeAzureV3Config_Update = `
should_persist_os_disk = false
os_disk_persistence_mode = "reattach"
should_persist_data_disks = false
data_disks_persistence_mode = "onLaunch"
should_persist_network = true
should_persist_vm = false
`

//endregion

// region Stateful Node Azure : Strategy
func TestAccSpotinstStatefulNodeAzureV3_Strategy(t *testing.T) {
	statefulNodeName := "test-acc-sn-azure-v3-strategy"
	resourceName := createStatefulNodeAzureV3ResourceName(statefulNodeName)

	var node azurev3.StatefulNode
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t, "azure") },
		Providers:    TestAccProviders,
		CheckDestroy: testStatefulNodeAzureV3Destroy,

		Steps: []resource.TestStep{
			{
				Config: createStatefulNodeAzureV3Terraform(&AzureV3StatefulNodeConfigMetadata{statefulNodeName: statefulNodeName,
					strategy: testStrategyStatefulNodeAzureV3Config_Create,
				}),
				Check: resource.ComposeTestCheckFunc(
					testCheckStatefulNodeAzureV3Exists(&node, resourceName),
					testCheckStatefulNodeAzureV3Attributes(&node, statefulNodeName),
					resource.TestCheckResourceAttr(resourceName, "strategy.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "strategy.0.fallback_to_on_demand", "true"),
					resource.TestCheckResourceAttr(resourceName, "strategy.0.draining_timeout", "40"),
					resource.TestCheckResourceAttr(resourceName, "strategy.0.revert_to_spot.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "strategy.0.revert_to_spot.0.perform_at", "always"),
				),
			},
			{
				Config: createStatefulNodeAzureV3Terraform(&AzureV3StatefulNodeConfigMetadata{
					statefulNodeName:     statefulNodeName,
					strategy:             testStrategyStatefulNodeAzureV3Config_Update,
					updateBaselineFields: true,
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "strategy.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "strategy.0.fallback_to_on_demand", "true"),
					resource.TestCheckResourceAttr(resourceName, "strategy.0.draining_timeout", "20"),
					resource.TestCheckResourceAttr(resourceName, "strategy.0.revert_to_spot.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "strategy.0.revert_to_spot.0.perform_at", "always"),
				),
			},
		},
	})
}

// TODO - another test that check signal?

const testStrategyStatefulNodeAzureV3Config_Create = `
strategy {
	draining_timeout =  40
	fallback_to_on_demand = true
	revert_to_spot {
		perform_at =  "always"
	}
}
`

const testStrategyStatefulNodeAzureV3Config_Update = `
strategy {
	draining_timeout =  20
	fallback_to_on_demand = true
	revert_to_spot {
		perform_at =  "always"
	}
}
`

//endregion

// region Stateful Node Azure : Health
func TestAccSpotinstStatefulNodeAzureV3_Health(t *testing.T) {
	statefulNodeName := "terraform-tests-do-not-delete"
	resourceName := createStatefulNodeAzureV3ResourceName(statefulNodeName)

	var node azurev3.StatefulNode
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t, "azure") },
		Providers:    TestAccProviders,
		CheckDestroy: testStatefulNodeAzureV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: createStatefulNodeAzureV3Terraform(&AzureV3StatefulNodeConfigMetadata{statefulNodeName: statefulNodeName,
					health: testHealthStatefulNodeAzureV3Config_Create,
				}),
				Check: resource.ComposeTestCheckFunc(
					testCheckStatefulNodeAzureV3Exists(&node, resourceName),
					testCheckStatefulNodeAzureV3Attributes(&node, statefulNodeName),
					resource.TestCheckResourceAttr(resourceName, "health.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "health.0.grace_period", "120"),
					resource.TestCheckResourceAttr(resourceName, "health.0.unhealthy_duration", "300"),
					resource.TestCheckResourceAttr(resourceName, "health.0.auto_healing", "true"),
					resource.TestCheckResourceAttr(resourceName, "health.0.health_check_types.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "health.0.health_check_types.0", "vmState"),
				),
			},
			{
				Config: createStatefulNodeAzureV3Terraform(&AzureV3StatefulNodeConfigMetadata{
					statefulNodeName:     statefulNodeName,
					health:               testHealthStatefulNodeAzureV3Config_Update,
					updateBaselineFields: true,
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "health.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "health.0.grace_period", "180"),
					resource.TestCheckResourceAttr(resourceName, "health.0.unhealthy_duration", "360"),
					resource.TestCheckResourceAttr(resourceName, "health.0.auto_healing", "true"),
					resource.TestCheckResourceAttr(resourceName, "health.0.health_check_types.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "health.0.health_check_types.0", "vmState"),
				),
			},
		},
	})
}

const testHealthStatefulNodeAzureV3Config_Create = `
health {
	health_check_types = ["vmState"]
	unhealthy_duration = "300"
	grace_period = "120"
	auto_healing = true
}
`

const testHealthStatefulNodeAzureV3Config_Update = `
health {
	health_check_types = ["vmState"]
	unhealthy_duration = "360"
	grace_period = "180"
	auto_healing = true
}
`

//endregion

// region Stateful Node Azure : VMSizes
func TestAccSpotinstStatefulNodeAzureV3_VMSizes(t *testing.T) {
	statefulNodeName := "terraform-tests-do-not-delete"
	resourceName := createStatefulNodeAzureV3ResourceName(statefulNodeName)

	var node azurev3.StatefulNode
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t, "azure") },
		Providers:    TestAccProviders,
		CheckDestroy: testStatefulNodeAzureV3Destroy,

		Steps: []resource.TestStep{
			{
				Config: createStatefulNodeAzureV3Terraform(&AzureV3StatefulNodeConfigMetadata{statefulNodeName: statefulNodeName,
					vmSizes: testVMSizesStatefulNodeAzureV3Config_Create}),
				Check: resource.ComposeTestCheckFunc(
					testCheckStatefulNodeAzureV3Exists(&node, resourceName),
					testCheckStatefulNodeAzureV3Attributes(&node, statefulNodeName),
					resource.TestCheckResourceAttr(resourceName, "spot_sizes.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "spot_sizes.0", "standard_ds2_v2"),
					resource.TestCheckResourceAttr(resourceName, "od_sizes.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "od_sizes.0", "standard_ds1_v2"),
					resource.TestCheckResourceAttr(resourceName, "preferred_spot_sizes.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "preferred_spot_sizes.0", "standard_ds2_v2"),
				),
			},
			{
				Config: createStatefulNodeAzureV3Terraform(&AzureV3StatefulNodeConfigMetadata{
					statefulNodeName:     statefulNodeName,
					vmSizes:              testVMSizesStatefulNodeAzureV3Config_Update,
					updateBaselineFields: true,
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "spot_sizes.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "spot_sizes.0", "standard_ds3_v2"),
					resource.TestCheckResourceAttr(resourceName, "od_sizes.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "od_sizes.0", "standard_ds4_v2"),
					resource.TestCheckResourceAttr(resourceName, "preferred_spot_sizes.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "preferred_spot_sizes.0", "standard_ds3_v2"),
				),
			},
		},
	})
}

const testVMSizesStatefulNodeAzureV3Config_Create = `
spot_sizes = ["standard_ds2_v2"]
od_sizes = ["standard_ds1_v2"]
preferred_spot_sizes =  ["standard_ds2_v2"]
`

const testVMSizesStatefulNodeAzureV3Config_Update = `
spot_sizes = ["standard_ds3_v2"]
od_sizes = ["standard_ds4_v2"]
preferred_spot_sizes =  ["standard_ds3_v2"]
`

//endregion

// region Stateful Node Azure : Scheduling
func TestAccSpotinstStatefulNodeAzureV3_Scheduling(t *testing.T) {
	statefulNodeName := "terraform-tests-do-not-delete"
	resourceName := createStatefulNodeAzureV3ResourceName(statefulNodeName)

	var node azurev3.StatefulNode
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t, "azure") },
		Providers:    TestAccProviders,
		CheckDestroy: testStatefulNodeAzureV3Destroy,

		Steps: []resource.TestStep{
			{
				Config: createStatefulNodeAzureV3Terraform(&AzureV3StatefulNodeConfigMetadata{statefulNodeName: statefulNodeName}),
				Check: resource.ComposeTestCheckFunc(
					testCheckStatefulNodeAzureV3Exists(&node, resourceName),
					testCheckStatefulNodeAzureV3Attributes(&node, statefulNodeName),
					resource.TestCheckResourceAttr(resourceName, "scheduling_task.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "scheduling_task.1113301976.is_enabled", "true"), //TODO - get hashcode somehow
					resource.TestCheckResourceAttr(resourceName, "scheduling_task.1113301976.cron_expression", "44 10 * * *"),
					resource.TestCheckResourceAttr(resourceName, "scheduling_task.1113301976.type", "pause"),
				),
			},
			{
				Config: createStatefulNodeAzureV3Terraform(&AzureV3StatefulNodeConfigMetadata{
					statefulNodeName:     statefulNodeName,
					scheduling:           testSchedulingStatefulNodeAzureV3Config_Update,
					updateBaselineFields: true,
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "scheduling_task.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "scheduling_task.1113301976.is_enabled", "true"), //TODO - get hashcode somehow
					resource.TestCheckResourceAttr(resourceName, "scheduling_task.1113301976.cron_expression", "48 10 * * *"),
					resource.TestCheckResourceAttr(resourceName, "scheduling_task.1113301976.type", "resume"),
				),
			},
		},
	})
}

const testSchedulingStatefulNodeAzureV3Config_Create = `
  scheduling_task {
    is_enabled = true
	cron_expression = "44 10 * * *"
    type = "pause"
  }
`

const testSchedulingStatefulNodeAzureV3Config_Update = `
  scheduling_task {
    is_enabled = true
	cron_expression = "48 10 * * *"
    type = "resume"
  }
`

//endregion

// region Stateful Node Azure : Tag
func TestAccSpotinstStatefulNodeAzureV3_Tag(t *testing.T) {
	statefulNodeName := "terraform-tests-do-not-delete"
	resourceName := createStatefulNodeAzureV3ResourceName(statefulNodeName)

	var node azurev3.StatefulNode
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t, "azure") },
		Providers:    TestAccProviders,
		CheckDestroy: testStatefulNodeAzureV3Destroy,

		Steps: []resource.TestStep{
			{
				Config: createStatefulNodeAzureV3Terraform(&AzureV3StatefulNodeConfigMetadata{statefulNodeName: statefulNodeName}),
				Check: resource.ComposeTestCheckFunc(
					testCheckStatefulNodeAzureV3Exists(&node, resourceName),
					testCheckStatefulNodeAzureV3Attributes(&node, statefulNodeName),
					resource.TestCheckResourceAttr(resourceName, "tag.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "tag.1113301976.tag_key", "Creator"), //TODO - get hashcode somehow
					resource.TestCheckResourceAttr(resourceName, "tag.1113301976.tag_value", "tamiry@netapp.com"),
				),
			},
			{
				Config: createStatefulNodeAzureV3Terraform(&AzureV3StatefulNodeConfigMetadata{
					statefulNodeName:     statefulNodeName,
					tag:                  testTagStatefulNodeAzureV3Config_Update,
					updateBaselineFields: true,
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "tag.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "tag.1113301976.tag_key", "Creator"), //TODO - get hashcode somehow
					resource.TestCheckResourceAttr(resourceName, "tag.1113301976.tag_value", "talg@netapp.com"),
				),
			},
		},
	})
}

const testTagStatefulNodeAzureV3Config_Create = `
tag {
	tag_key = "Creator"
	tag_value = "tamiry@netapp.com"
}
`

const testTagStatefulNodeAzureV3Config_Update = `
tag {
	tag_key = "Creator"
	tag_value = "talg@netapp.com"
}
`

//endregion

// region Stateful Node Azure : Image
func TestAccSpotinstStatefulNodeAzureV3_Image(t *testing.T) {
	statefulNodeName := "terraform-tests-do-not-delete"
	resourceName := createStatefulNodeAzureV3ResourceName(statefulNodeName)

	var node azurev3.StatefulNode
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t, "azure") },
		Providers:    TestAccProviders,
		CheckDestroy: testStatefulNodeAzureV3Destroy,

		Steps: []resource.TestStep{
			{
				Config: createStatefulNodeAzureV3Terraform(&AzureV3StatefulNodeConfigMetadata{statefulNodeName: statefulNodeName,
					image: testImageStatefulNodeAzureV3Config_Create}),
				Check: resource.ComposeTestCheckFunc(
					testCheckStatefulNodeAzureV3Exists(&node, resourceName),
					testCheckStatefulNodeAzureV3Attributes(&node, statefulNodeName),
					resource.TestCheckResourceAttr(resourceName, "image.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "image.0.marketplace_image.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "image.0.marketplace_image.0.offer", "UbuntuServer"),
					resource.TestCheckResourceAttr(resourceName, "image.0.marketplace_image.0.publisher", "Canonical"),
					resource.TestCheckResourceAttr(resourceName, "image.0.marketplace_image.0.sku", "18.04-LTS"),
					resource.TestCheckResourceAttr(resourceName, "image.0.marketplace_image.0.version", "latest"),
					resource.TestCheckResourceAttr(resourceName, "image.0.custom_image.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "image.0.gallery.#", "0"),
				),
			},
			{
				Config: createStatefulNodeAzureV3Terraform(&AzureV3StatefulNodeConfigMetadata{
					statefulNodeName:     statefulNodeName,
					image:                testImageStatefulNodeAzureV3Config_Update,
					updateBaselineFields: true,
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "image.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "image.0.marketplace_image.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "image.0.marketplace_image.0.offer", "UbuntuServer"),
					resource.TestCheckResourceAttr(resourceName, "image.0.marketplace_image.0.publisher", "Canonical"),
					resource.TestCheckResourceAttr(resourceName, "image.0.marketplace_image.0.sku", "18.04-LTS"),
					resource.TestCheckResourceAttr(resourceName, "image.0.marketplace_image.0.version", "latest"),
					resource.TestCheckResourceAttr(resourceName, "image.0.custom_image.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "image.0.gallery.#", "0"),
				),
			},
		},
	})
}

const testImageStatefulNodeAzureV3Config_Create = `
image {
	marketplace_image {
		publisher = "Canonical"
		offer = "UbuntuServer"
		sku = "18.04-LTS"
		version = "latest"
	}
}
`

const testImageStatefulNodeAzureV3Config_Update = `
image {
	marketplace_image {
		publisher = "Canonical"
		offer = "UbuntuServer"
		sku = "18.04-LTS"
		version = "latest"
	}
}
`

//endregion

// region Stateful Node Azure : Network
func TestAccSpotinstStatefulNodeAzureV3_Network(t *testing.T) {
	statefulNodeName := "terraform-tests-do-not-delete"
	resourceName := createStatefulNodeAzureV3ResourceName(statefulNodeName)

	var node azurev3.StatefulNode
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t, "azure") },
		Providers:    TestAccProviders,
		CheckDestroy: testStatefulNodeAzureV3Destroy,

		Steps: []resource.TestStep{
			{
				Config: createStatefulNodeAzureV3Terraform(&AzureV3StatefulNodeConfigMetadata{statefulNodeName: statefulNodeName,
					network: testNetworkStatefulNodeAzureV3Config_Create}),
				Check: resource.ComposeTestCheckFunc(
					testCheckStatefulNodeAzureV3Exists(&node, resourceName),
					testCheckStatefulNodeAzureV3Attributes(&node, statefulNodeName),
					resource.TestCheckResourceAttr(resourceName, "network.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "network.0.network_interface.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "network.0.network_interface.0.subnet_name", "default"),
					resource.TestCheckResourceAttr(resourceName, "network.0.network_interface.0.assign_public_ip", "true"),
					resource.TestCheckResourceAttr(resourceName, "network.0.network_interface.0.is_primary", "true"),
					resource.TestCheckResourceAttr(resourceName, "network.0.network_interface.0.public_ip_sku", "Standard"),
					resource.TestCheckResourceAttr(resourceName, "network.0.network_interface.0.network_security_group.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "network.0.network_interface.0.network_security_group.0.name", "core-reliability-network-security-group"),
					resource.TestCheckResourceAttr(resourceName, "network.0.network_interface.0.network_security_group.0.network_resource_group_name", "CoreReliabilityResourceGroup"),
					resource.TestCheckResourceAttr(resourceName, "network.0.network_interface.0.enable_ip_forwarding", "true"),
					resource.TestCheckResourceAttr(resourceName, "network.0.network_resource_group_name", "CoreReliabilityResourceGroup"), //TODO - confidential?
					resource.TestCheckResourceAttr(resourceName, "network.0.virtual_network_name", "CoreReliabilityVN"),
				),
			},
			{
				Config: createStatefulNodeAzureV3Terraform(&AzureV3StatefulNodeConfigMetadata{
					statefulNodeName:     statefulNodeName,
					network:              testNetworkStatefulNodeAzureV3Config_Update,
					updateBaselineFields: true,
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "network.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "network.0.network_interface.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "network.0.network_interface.0.subnet_name", "default"),
					resource.TestCheckResourceAttr(resourceName, "network.0.network_interface.0.assign_public_ip", "true"),
					resource.TestCheckResourceAttr(resourceName, "network.0.network_interface.0.is_primary", "true"),
					resource.TestCheckResourceAttr(resourceName, "network.0.network_interface.0.public_ip_sku", "Standard"),
					resource.TestCheckResourceAttr(resourceName, "network.0.network_interface.0.network_security_group.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "network.0.network_interface.0.network_security_group.0.name", "core-reliability-network-security-group"),
					resource.TestCheckResourceAttr(resourceName, "network.0.network_interface.0.network_security_group.0.network_resource_group_name", "CoreReliabilityResourceGroup"),
					resource.TestCheckResourceAttr(resourceName, "network.0.network_interface.0.enable_ip_forwarding", "true"),
					resource.TestCheckResourceAttr(resourceName, "network.0.network_resource_group_name", "CoreReliabilityResourceGroup"),
					resource.TestCheckResourceAttr(resourceName, "network.0.virtual_network_name", "CoreReliabilityVN"),
				),
			},
		},
	})
}

const testNetworkStatefulNodeAzureV3Config_Create = `
network {
	network_resource_group_name = "CoreReliabilityResourceGroup"
	virtual_network_name = "CoreReliabilityVN"
	network_interface {
		subnet_name = "default"
		assign_public_ip = true
		is_primary = true
		public_ip_sku = "Standard"
		network_security_group {
			name = "core-reliability-network-security-group"
			network_resource_group_name = "CoreReliabilityResourceGroup"
		}
		enable_ip_forwarding = true
	}
}
`

const testNetworkStatefulNodeAzureV3Config_Update = `
network {
	network_resource_group_name = "CoreReliabilityResourceGroup"
	virtual_network_name = "CoreReliabilityVN"
	network_interface {
		subnet_name = "default"
		assign_public_ip = true
		is_primary = true
		public_ip_sku = "Standard"
		network_security_group {
			name = "core-reliability-network-security-group"
			network_resource_group_name = "CoreReliabilityResourceGroup"
		}
		enable_ip_forwarding = true
	}
}
`

//endregion

// region Stateful Node Azure : OSDisk
func TestAccSpotinstStatefulNodeAzureV3_OSDisk(t *testing.T) {
	statefulNodeName := "terraform-tests-do-not-delete"
	resourceName := createStatefulNodeAzureV3ResourceName(statefulNodeName)

	var node azurev3.StatefulNode
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t, "azure") },
		Providers:    TestAccProviders,
		CheckDestroy: testStatefulNodeAzureV3Destroy,

		Steps: []resource.TestStep{
			{
				Config: createStatefulNodeAzureV3Terraform(&AzureV3StatefulNodeConfigMetadata{statefulNodeName: statefulNodeName}),
				Check: resource.ComposeTestCheckFunc(
					testCheckStatefulNodeAzureV3Exists(&node, resourceName),
					testCheckStatefulNodeAzureV3Attributes(&node, statefulNodeName),
					resource.TestCheckResourceAttr(resourceName, "os_disk.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "os_disk.0.size_gb", "30"),
					resource.TestCheckResourceAttr(resourceName, "os_disk.0.type", "Standard_LRS"),
				),
			},
			{
				Config: createStatefulNodeAzureV3Terraform(&AzureV3StatefulNodeConfigMetadata{
					statefulNodeName:     statefulNodeName,
					osDisk:               testOSDiskStatefulNodeAzureV3Config_Update,
					updateBaselineFields: true,
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "os_disk.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "os_disk.0.size_gb", "40"),
					resource.TestCheckResourceAttr(resourceName, "os_disk.0.type", "Standard_LRS"),
				),
			},
		},
	})
}

const testOSDiskStatefulNodeAzureV3Config_Create = `
os_disk {
	size_gb = 30
	type = "Standard_LRS"
}
`

const testOSDiskStatefulNodeAzureV3Config_Update = `
os_disk {
	size_gb = 40
	type = "Standard_LRS"
}
`

// region Stateful Node Azure : DataDisk
func TestAccSpotinstStatefulNodeAzureV3_DataDisk(t *testing.T) {
	statefulNodeName := "terraform-tests-do-not-delete"
	resourceName := createStatefulNodeAzureV3ResourceName(statefulNodeName)

	var node azurev3.StatefulNode
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t, "azure") },
		Providers:    TestAccProviders,
		CheckDestroy: testStatefulNodeAzureV3Destroy,

		Steps: []resource.TestStep{
			{
				Config: createStatefulNodeAzureV3Terraform(&AzureV3StatefulNodeConfigMetadata{statefulNodeName: statefulNodeName}),
				Check: resource.ComposeTestCheckFunc(
					testCheckStatefulNodeAzureV3Exists(&node, resourceName),
					testCheckStatefulNodeAzureV3Attributes(&node, statefulNodeName),
					resource.TestCheckResourceAttr(resourceName, "data_disk.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "data_disk.0.size_gb", "1"),
					resource.TestCheckResourceAttr(resourceName, "data_disk.0.lun", "1"),
					resource.TestCheckResourceAttr(resourceName, "data_disk.0.type", "Standard_LRS"),
				),
			},
			{
				Config: createStatefulNodeAzureV3Terraform(&AzureV3StatefulNodeConfigMetadata{
					statefulNodeName:     statefulNodeName,
					dataDisk:             testDataDiskStatefulNodeAzureV3Config_Update,
					updateBaselineFields: true,
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "data_disk.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "data_disk.0.size_gb", "2"),
					resource.TestCheckResourceAttr(resourceName, "data_disk.0.lun", "1"),
					resource.TestCheckResourceAttr(resourceName, "data_disk.0.type", "Standard_LRS"),
				),
			},
		},
	})
}

const testDataDiskStatefulNodeAzureV3Config_Create = `
data_disk {
	size_gb = 1
	lun = 1
	type = "Standard_LRS"
}
`

const testDataDiskStatefulNodeAzureV3Config_Update = `
data_disk {
	size_gb = 2
	lun = 1
	type = "Standard_LRS"
}
`

// region Stateful Node Azure : BootDiagnostics
func TestAccSpotinstStatefulNodeAzureV3_BootDiagnostics(t *testing.T) {
	statefulNodeName := "terraform-tests-do-not-delete"
	resourceName := createStatefulNodeAzureV3ResourceName(statefulNodeName)

	var node azurev3.StatefulNode
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t, "azure") },
		Providers:    TestAccProviders,
		CheckDestroy: testStatefulNodeAzureV3Destroy,

		Steps: []resource.TestStep{
			{
				Config: createStatefulNodeAzureV3Terraform(&AzureV3StatefulNodeConfigMetadata{statefulNodeName: statefulNodeName}),
				Check: resource.ComposeTestCheckFunc(
					testCheckStatefulNodeAzureV3Exists(&node, resourceName),
					testCheckStatefulNodeAzureV3Attributes(&node, statefulNodeName),
					resource.TestCheckResourceAttr(resourceName, "boot_diagnostics.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "boot_diagnostics.0.is_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "boot_diagnostics.0.type", "managed"),
					resource.TestCheckResourceAttr(resourceName, "boot_diagnostics.0.storage_url", "blob.core.windows.net/test123"),
				),
			},
			{
				Config: createStatefulNodeAzureV3Terraform(&AzureV3StatefulNodeConfigMetadata{
					statefulNodeName:     statefulNodeName,
					bootDiagnostics:      testBootDiagnosticsStatefulNodeAzureV3Config_Update,
					updateBaselineFields: true,
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "boot_diagnostics.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "boot_diagnostics.0.is_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "boot_diagnostics.0.type", "unmanaged"),
					resource.TestCheckResourceAttr(resourceName, "boot_diagnostics.0.storage_url", "blob.core.windows.net/test123"),
				),
			},
		},
	})
}

const testBootDiagnosticsStatefulNodeAzureV3Config_Create = `
boot_diagnostics {
	is_enabled = true
	type = "managed"
	storage_url = "blob.core.windows.net/test123"
}
`

const testBootDiagnosticsStatefulNodeAzureV3Config_Update = `
boot_diagnostics {
	is_enabled = true
	type = "unmanaged"
	storage_url = "blob.core.windows.net/test123"
}
`

// region Stateful Node Azure : Extensions
func TestAccSpotinstStatefulNodeAzureV3_Extensions(t *testing.T) {
	statefulNodeName := "terraform-tests-do-not-delete"
	resourceName := createStatefulNodeAzureV3ResourceName(statefulNodeName)

	var node azurev3.StatefulNode
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t, "azure") },
		Providers:    TestAccProviders,
		CheckDestroy: testStatefulNodeAzureV3Destroy,

		Steps: []resource.TestStep{
			{
				Config: createStatefulNodeAzureV3Terraform(&AzureV3StatefulNodeConfigMetadata{statefulNodeName: statefulNodeName,
					extensions: testExtensionsStatefulNodeAzureV3Config_Create}),
				Check: resource.ComposeTestCheckFunc(
					testCheckStatefulNodeAzureV3Exists(&node, resourceName),
					testCheckStatefulNodeAzureV3Attributes(&node, statefulNodeName),
					resource.TestCheckResourceAttr(resourceName, "extension.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "extension.1031128857.api_version", "2.0"), //TODO - get hashcode somehow
					resource.TestCheckResourceAttr(resourceName, "extension.1031128857.minor_version_auto_upgrade", "true"),
					resource.TestCheckResourceAttr(resourceName, "extension.1031128857.name", "terraform-extension"),
					resource.TestCheckResourceAttr(resourceName, "extension.1031128857.publisher", "Microsoft.Azure.Extensions"),
					resource.TestCheckResourceAttr(resourceName, "extension.1031128857.type", "Linux"),
					resource.TestCheckResourceAttr(resourceName, "extension.1031128857.protected_settings.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "extension.1031128857.protected_settings.0.script", "IyEvYmluL2Jhc2gKZWNobyAibmlyIiA+IC9ob21lL25pci9uaXIudHh0Cg=="), //ToDo check about field script under protected settings
				),
			},
			{
				Config: createStatefulNodeAzureV3Terraform(&AzureV3StatefulNodeConfigMetadata{
					statefulNodeName:     statefulNodeName,
					extensions:           testExtensionsStatefulNodeAzureV3Config_Update,
					updateBaselineFields: true,
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "extension.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "extension.1031128857.api_version", "2.0"), //TODO - get hashcode somehow
					resource.TestCheckResourceAttr(resourceName, "extension.1031128857.minor_version_auto_upgrade", "false"),
					resource.TestCheckResourceAttr(resourceName, "extension.1031128857.name", "terraform-extension"),
					resource.TestCheckResourceAttr(resourceName, "extension.1031128857.publisher", "Microsoft.Azure.Extensions"),
					resource.TestCheckResourceAttr(resourceName, "extension.1031128857.type", "Linux"),
					resource.TestCheckResourceAttr(resourceName, "extension.1031128857.protected_settings.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "extension.1031128857.protected_settings.0.script", "IyEvYmluL2Jhc2gKZWNobyAibmlyIiA+IC9ob21lL25pci9uaXIudHh0Cg=="), //ToDo check about field script under protected settings
				),
			},
		},
	})
}

const testExtensionsStatefulNodeAzureV3Config_Create = `
extension {
	api_version = "2.0"
	minor_version_auto_upgrade = true
	name = "terraform-extension"
	publisher = "Microsoft.Azure.Extensions"
	type = "Linux"
	protected_settings = {
		script = "IyEvYmluL2Jhc2gKZWNobyAibmlyIiA+IC9ob21lL25pci9uaXIudHh0Cg=="
	}
}
`
const testExtensionsStatefulNodeAzureV3Config_Update = `
extension {
	api_version = "2.0"
	minor_version_auto_upgrade = false
	name = "terraform-extension"
	publisher = "Microsoft.Azure.Extensions"
	type = "Linux"
	protected_settings = {
		script = "IyEvYmluL2Jhc2gKZWNobyAibmlyIiA+IC9ob21lL25pci9uaXIudHh0Cg=="
	}
}
`

//endregion

// region Stateful Node Azure : Secret
func TestAccSpotinstStatefulNodeAzureV3_Secret(t *testing.T) {
	statefulNodeName := "terraform-tests-do-not-delete"
	resourceName := createStatefulNodeAzureV3ResourceName(statefulNodeName)

	var node azurev3.StatefulNode
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t, "azure") },
		Providers:    TestAccProviders,
		CheckDestroy: testStatefulNodeAzureV3Destroy,

		Steps: []resource.TestStep{
			{
				Config: createStatefulNodeAzureV3Terraform(&AzureV3StatefulNodeConfigMetadata{statefulNodeName: statefulNodeName,
					extensions: testExtensionsStatefulNodeAzureV3Config_Create}),
				Check: resource.ComposeTestCheckFunc(
					testCheckStatefulNodeAzureV3Exists(&node, resourceName),
					testCheckStatefulNodeAzureV3Attributes(&node, statefulNodeName),
					resource.TestCheckResourceAttr(resourceName, "secrets.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "secrets.1031128857.source_vault.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "secrets.1031128857.source_vault.name", "core-reliability-source-vault-name"),
					resource.TestCheckResourceAttr(resourceName, "secrets.1031128857.source_vault.resource_group_name", "CoreReliabilityResourceGroup"),
					resource.TestCheckResourceAttr(resourceName, "secrets.1031128857.vault_certificates.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "secrets.1031128857.vault_certificates.certificate_url", "core-reliability-certificate-url"),
					resource.TestCheckResourceAttr(resourceName, "secrets.1031128857.vault_certificates.certificate_store", "core-reliability-certificate-store"),
				),
			},
			{
				Config: createStatefulNodeAzureV3Terraform(&AzureV3StatefulNodeConfigMetadata{
					statefulNodeName:     statefulNodeName,
					secret:               testSecretsStatefulNodeAzureV3Config_Update,
					updateBaselineFields: true,
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "secrets.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "secrets.1031128857.source_vault.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "secrets.1031128857.source_vault.name", "core-reliability-source-vault-name"),
					resource.TestCheckResourceAttr(resourceName, "secrets.1031128857.source_vault.resource_group_name", "CoreReliabilityResourceGroup3"),
					resource.TestCheckResourceAttr(resourceName, "secrets.1031128857.vault_certificates.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "secrets.1031128857.vault_certificates.certificate_url", "core-reliability-certificate-url"),
					resource.TestCheckResourceAttr(resourceName, "secrets.1031128857.vault_certificates.certificate_store", "core-reliability-certificate-store-is"),
				),
			},
		},
	})
}

const testSecretsStatefulNodeAzureV3Config_Create = `
    secrets {
      source_vault {
      resource_group_name = "CoreReliabilityResourceGroup"
      name = "core-reliability-source-vault-name"
	}
      vault_certificates {
      certificate_url = "core-reliability-certificate-url"
      certificate_store = "core-reliability-certificate-store"
	}
}
`
const testSecretsStatefulNodeAzureV3Config_Update = `
    secrets {
      source_vault {
      resource_group_name = "CoreReliabilityResourceGroup3"
      name = "core-reliability-source-vault-name"
	}
      vault_certificates {
      certificate_url = "core-reliability-certificate-url"
      certificate_store = "core-reliability-certificate-store-is"
	}
}
`

//endregion
