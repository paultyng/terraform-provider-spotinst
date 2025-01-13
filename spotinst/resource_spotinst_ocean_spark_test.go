package spotinst

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/spotinst/spotinst-sdk-go/service/ocean/spark"
	"github.com/spotinst/spotinst-sdk-go/spotinst"

	"github.com/spotinst/terraform-provider-spotinst/spotinst/commons"
)

const (
	testOceanClusterID = "o-7cb5c0b8"
)

var oceanClusterID = getOceanClusterID() // NOTE: This needs to be an existing ocean cluster

func init() {
	resource.AddTestSweepers("spotinst_ocean_spark", &resource.Sweeper{
		Name: "spotinst_ocean_spark",
		F:    testSweepOceanSpark,
	})
}

func getOceanClusterID() string {
	// Prefer environment variable
	oceanClusterID := os.Getenv("TEST_ACC_OCEAN_SPARK_OCEAN_ID")
	if oceanClusterID == "" {
		// Default to hardcoded ID
		oceanClusterID = testOceanClusterID
	}

	return oceanClusterID
}

func testSweepOceanSpark(_ string) error {
	client, err := getProviderClient("aws")
	if err != nil {
		return fmt.Errorf("error getting client: %v", err)
	}

	conn := client.(*Client).ocean.Spark()
	input := &spark.ListClustersInput{}
	if resp, err := conn.ListClusters(context.Background(), input); err != nil {
		return fmt.Errorf("error getting list of clusters to sweep")
	} else {
		if len(resp.Clusters) == 0 {
			log.Printf("[INFO] No clusters to sweep")
		}
		for _, cluster := range resp.Clusters {
			// Our test clusters should have a controller cluster ID starting with "tf-test-acc-"
			if strings.Contains(spotinst.StringValue(cluster.ControllerClusterID), "tf-test-acc-") {
				if _, err := conn.DeleteCluster(context.Background(), &spark.DeleteClusterInput{ClusterID: cluster.ID}); err != nil {
					return fmt.Errorf("unable to delete cluster %v in sweep", spotinst.StringValue(cluster.ID))
				} else {
					log.Printf("Sweeper deleted %v\n", spotinst.StringValue(cluster.ID))
				}
			}
		}
	}
	return nil
}

func createOceanSparkResourceName(name string) string {
	return fmt.Sprintf("%v.%v", string(commons.OceanSparkResourceName), name)
}

func testOceanSparkAWSDestroy(s *terraform.State) error {
	client := testAccProviderAWS.Meta().(*Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != string(commons.OceanSparkResourceName) {
			continue
		}
		input := &spark.ReadClusterInput{ClusterID: spotinst.String(rs.Primary.ID)}
		resp, err := client.ocean.Spark().ReadCluster(context.Background(), input)
		if err == nil && resp != nil && resp.Cluster != nil {
			return fmt.Errorf("cluster still exists")
		}
	}
	return nil
}

func testCheckOceanSparkAttributes(cluster *spark.Cluster, oceanClusterID string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if spotinst.StringValue(cluster.OceanClusterID) != oceanClusterID {
			return fmt.Errorf("bad content: %v", cluster.OceanClusterID)
		}
		return nil
	}
}

func testCheckOceanSparkExists(cluster *spark.Cluster, resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no resource ID is set")
		}
		client := testAccProviderAWS.Meta().(*Client)
		input := &spark.ReadClusterInput{ClusterID: spotinst.String(rs.Primary.ID)}
		resp, err := client.ocean.Spark().ReadCluster(context.Background(), input)
		if err != nil {
			return err
		}
		if spotinst.StringValue(resp.Cluster.OceanClusterID) != rs.Primary.Attributes["ocean_cluster_id"] {
			return fmt.Errorf("Cluster not found: %+v,\n %+v\n", resp.Cluster, rs.Primary.Attributes)
		}
		*cluster = *resp.Cluster
		return nil
	}
}

type SparkClusterConfigMetadata struct {
	oceanClusterID string
	fieldsToAppend string
}

func createOceanSparkTerraform(sccm *SparkClusterConfigMetadata) string {
	if sccm == nil {
		return ""
	}

	format := testBaseSparkConfig
	template := fmt.Sprintf(format,
		sccm.oceanClusterID,
		sccm.oceanClusterID,
		sccm.fieldsToAppend,
	)

	log.Printf("Terraform [%v] template:\n%v", sccm.oceanClusterID, template)
	return template
}

const testBaseSparkConfig = `
resource "` + string(commons.OceanSparkResourceName) + `" "%v" {
  provider = "aws"

  ocean_cluster_id = "%v"

  %v
}
`

func TestAccSpotinstOceanSpark_noConfig(t *testing.T) {
	resourceName := createOceanSparkResourceName(oceanClusterID)

	var cluster spark.Cluster
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t, "aws") },
		Providers:    TestAccProviders,
		CheckDestroy: testOceanSparkAWSDestroy,

		Steps: []resource.TestStep{
			{
				Config: createOceanSparkTerraform(&SparkClusterConfigMetadata{
					oceanClusterID: oceanClusterID,
				}),
				Check: resource.ComposeTestCheckFunc(
					testCheckOceanSparkExists(&cluster, resourceName),
					testCheckOceanSparkAttributes(&cluster, oceanClusterID),
				),
			},
		},
	})
}

func TestAccSpotinstOceanSpark_withIngressConfig(t *testing.T) {
	resourceName := createOceanSparkResourceName(oceanClusterID)

	var cluster spark.Cluster
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t, "aws") },
		Providers:    TestAccProviders,
		CheckDestroy: testOceanSparkAWSDestroy,

		Steps: []resource.TestStep{
			{
				Config: createOceanSparkTerraform(&SparkClusterConfigMetadata{
					oceanClusterID: oceanClusterID,
					fieldsToAppend: testConfigWithIngressCreate,
				}),
				Check: resource.ComposeTestCheckFunc(
					testCheckOceanSparkExists(&cluster, resourceName),
					testCheckOceanSparkAttributes(&cluster, oceanClusterID),
					resource.TestCheckResourceAttr(resourceName, "ingress.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ingress.0.service_annotations.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "ingress.0.service_annotations.my-annotation-1", "my-annotation-value-1"),
					resource.TestCheckResourceAttr(resourceName, "ingress.0.service_annotations.my-annotation-2", "my-annotation-value-2"),
					resource.TestCheckResourceAttr(resourceName, "ingress.0.controller.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ingress.0.controller.0.managed", "true"),
					resource.TestCheckResourceAttr(resourceName, "ingress.0.load_balancer.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ingress.0.load_balancer.0.managed", "true"),
					resource.TestCheckResourceAttr(resourceName, "ingress.0.load_balancer.0.target_group_arn", "some-test-target-group-arn"),
					resource.TestCheckResourceAttr(resourceName, "ingress.0.load_balancer.0.service_annotations.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "ingress.0.load_balancer.0.service_annotations.my-lb-service-annotation-1", "my-lb-service-annotation-value-1"),
					resource.TestCheckResourceAttr(resourceName, "ingress.0.load_balancer.0.service_annotations.my-lb-service-annotation-2", "my-lb-service-annotation-value-2"),
					resource.TestCheckResourceAttr(resourceName, "ingress.0.custom_endpoint.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ingress.0.custom_endpoint.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "ingress.0.custom_endpoint.0.address", "valid-load.balancer"),
					resource.TestCheckResourceAttr(resourceName, "ingress.0.private_link.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ingress.0.private_link.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "ingress.0.private_link.0.vpc_endpoint_service", "test-vpc-endpoint-service"),
				),
			},
			{
				Config: createOceanSparkTerraform(&SparkClusterConfigMetadata{
					oceanClusterID: oceanClusterID,
					fieldsToAppend: testConfigWithIngressUpdate,
				}),
				Check: resource.ComposeTestCheckFunc(
					testCheckOceanSparkExists(&cluster, resourceName),
					testCheckOceanSparkAttributes(&cluster, oceanClusterID),
					resource.TestCheckResourceAttr(resourceName, "ingress.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ingress.0.service_annotations.%", "3"),
					resource.TestCheckResourceAttr(resourceName, "ingress.0.service_annotations.my-annotation-2", "my-annotation-value-2-updated"),
					resource.TestCheckResourceAttr(resourceName, "ingress.0.service_annotations.my-annotation-3", "my-annotation-value-3"),
					resource.TestCheckResourceAttr(resourceName, "ingress.0.service_annotations.my-new-annotation", "my-new-annotation-value"),
					resource.TestCheckResourceAttr(resourceName, "ingress.0.controller.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ingress.0.controller.0.managed", "true"),
					resource.TestCheckResourceAttr(resourceName, "ingress.0.load_balancer.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ingress.0.load_balancer.0.managed", "false"),
					resource.TestCheckResourceAttr(resourceName, "ingress.0.load_balancer.0.target_group_arn", "some-test-target-group-arn-active"),
					resource.TestCheckResourceAttr(resourceName, "ingress.0.load_balancer.0.service_annotations.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "ingress.0.load_balancer.0.service_annotations.my-lb-service-annotation-1", "my-lb-service-annotation-value-1"),
					resource.TestCheckResourceAttr(resourceName, "ingress.0.load_balancer.0.service_annotations.my-lb-service-annotation-3", "my-lb-service-annotation-value-3"),
					resource.TestCheckResourceAttr(resourceName, "ingress.0.custom_endpoint.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ingress.0.custom_endpoint.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "ingress.0.custom_endpoint.0.address", "active-load.balancer"),
					resource.TestCheckResourceAttr(resourceName, "ingress.0.private_link.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ingress.0.private_link.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "ingress.0.private_link.0.vpc_endpoint_service", "test-vpc-endpoint-service"),
				),
			},
			{
				Config: createOceanSparkTerraform(&SparkClusterConfigMetadata{
					oceanClusterID: oceanClusterID,
					fieldsToAppend: testConfigWithIngressUpdate2,
				}),
				Check: resource.ComposeTestCheckFunc(
					testCheckOceanSparkExists(&cluster, resourceName),
					testCheckOceanSparkAttributes(&cluster, oceanClusterID),
					resource.TestCheckResourceAttr(resourceName, "ingress.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ingress.0.service_annotations.%", "3"),
					resource.TestCheckResourceAttr(resourceName, "ingress.0.service_annotations.my-annotation-2", "my-annotation-value-2-updated"),
					resource.TestCheckResourceAttr(resourceName, "ingress.0.service_annotations.my-annotation-3", "my-annotation-value-3"),
					resource.TestCheckResourceAttr(resourceName, "ingress.0.service_annotations.my-new-annotation", "my-new-annotation-value"),
					resource.TestCheckResourceAttr(resourceName, "ingress.0.controller.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ingress.0.controller.0.managed", "false"),
					resource.TestCheckResourceAttr(resourceName, "ingress.0.load_balancer.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ingress.0.load_balancer.0.managed", "false"),
					resource.TestCheckResourceAttr(resourceName, "ingress.0.load_balancer.0.target_group_arn", "some-test-target-group-arn-inactive"),
					resource.TestCheckResourceAttr(resourceName, "ingress.0.load_balancer.0.service_annotations.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "ingress.0.custom_endpoint.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ingress.0.custom_endpoint.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "ingress.0.custom_endpoint.0.address", "inactive-load.balancer"),
					resource.TestCheckResourceAttr(resourceName, "ingress.0.private_link.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ingress.0.private_link.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "ingress.0.private_link.0.vpc_endpoint_service", "test-vpc-endpoint-service-active"),
				),
			},
			{
				// Reverts to default values if resources in terraform are empty
				Config: createOceanSparkTerraform(&SparkClusterConfigMetadata{
					oceanClusterID: oceanClusterID,
					fieldsToAppend: testConfigWithIngressEmptyFields,
				}),
				Check: resource.ComposeTestCheckFunc(
					testCheckOceanSparkExists(&cluster, resourceName),
					testCheckOceanSparkAttributes(&cluster, oceanClusterID),
					resource.TestCheckResourceAttr(resourceName, "ingress.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ingress.0.service_annotations.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "ingress.0.controller.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ingress.0.controller.0.managed", "true"),
					resource.TestCheckResourceAttr(resourceName, "ingress.0.load_balancer.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ingress.0.load_balancer.0.managed", "true"),
					resource.TestCheckResourceAttr(resourceName, "ingress.0.load_balancer.0.target_group_arn", ""),
					resource.TestCheckResourceAttr(resourceName, "ingress.0.load_balancer.0.service_annotations.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "ingress.0.custom_endpoint.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ingress.0.custom_endpoint.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "ingress.0.custom_endpoint.0.address", ""),
					resource.TestCheckResourceAttr(resourceName, "ingress.0.private_link.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ingress.0.private_link.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "ingress.0.private_link.0.vpc_endpoint_service", ""),
				),
			},
			{
				// Deletes config objects if resources not defined in terraform
				Config: createOceanSparkTerraform(&SparkClusterConfigMetadata{
					oceanClusterID: oceanClusterID,
					fieldsToAppend: testConfigWithIngressEmptyFields2,
				}),
				Check: resource.ComposeTestCheckFunc(
					testCheckOceanSparkExists(&cluster, resourceName),
					testCheckOceanSparkAttributes(&cluster, oceanClusterID),
					resource.TestCheckResourceAttr(resourceName, "ingress.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ingress.0.service_annotations.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "ingress.0.controller.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "ingress.0.load_balancer.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "ingress.0.custom_endpoint.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "ingress.0.private_link.#", "0"),
				),
			},
		},
	})
}

func TestAccSpotinstOceanSpark_withWebhookConfig(t *testing.T) {
	resourceName := createOceanSparkResourceName(oceanClusterID)

	var cluster spark.Cluster
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t, "aws") },
		Providers:    TestAccProviders,
		CheckDestroy: testOceanSparkAWSDestroy,

		Steps: []resource.TestStep{
			{
				Config: createOceanSparkTerraform(&SparkClusterConfigMetadata{
					oceanClusterID: oceanClusterID,
					fieldsToAppend: testConfigWithWebhookCreate,
				}),
				Check: resource.ComposeTestCheckFunc(
					testCheckOceanSparkExists(&cluster, resourceName),
					testCheckOceanSparkAttributes(&cluster, oceanClusterID),
					resource.TestCheckResourceAttr(resourceName, "webhook.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "webhook.0.use_host_network", "false"),
					resource.TestCheckResourceAttr(resourceName, "webhook.0.host_network_ports.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "webhook.0.host_network_ports.0", "12345"),
				),
			},
			{
				Config: createOceanSparkTerraform(&SparkClusterConfigMetadata{
					oceanClusterID: oceanClusterID,
					fieldsToAppend: testConfigWithWebhookUpdate,
				}),
				Check: resource.ComposeTestCheckFunc(
					testCheckOceanSparkExists(&cluster, resourceName),
					testCheckOceanSparkAttributes(&cluster, oceanClusterID),
					resource.TestCheckResourceAttr(resourceName, "webhook.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "webhook.0.use_host_network", "true"),
					resource.TestCheckResourceAttr(resourceName, "webhook.0.host_network_ports.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "webhook.0.host_network_ports.0", "12345"),
					resource.TestCheckResourceAttr(resourceName, "webhook.0.host_network_ports.1", "54321"),
				),
			},
			{
				Config: createOceanSparkTerraform(&SparkClusterConfigMetadata{
					oceanClusterID: oceanClusterID,
					fieldsToAppend: testConfigWithWebhookEmptyFields,
				}),
				Check: resource.ComposeTestCheckFunc(
					testCheckOceanSparkExists(&cluster, resourceName),
					testCheckOceanSparkAttributes(&cluster, oceanClusterID),
					resource.TestCheckResourceAttr(resourceName, "webhook.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "webhook.0.use_host_network", "false"),
					resource.TestCheckResourceAttr(resourceName, "webhook.0.host_network_ports.#", "0"),
				),
			},
		},
	})
}

func TestAccSpotinstOceanSpark_withComputeConfig(t *testing.T) {
	resourceName := createOceanSparkResourceName(oceanClusterID)

	var cluster spark.Cluster
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t, "aws") },
		Providers:    TestAccProviders,
		CheckDestroy: testOceanSparkAWSDestroy,

		Steps: []resource.TestStep{
			{
				Config: createOceanSparkTerraform(&SparkClusterConfigMetadata{
					oceanClusterID: oceanClusterID,
					fieldsToAppend: testConfigWithComputeCreate,
				}),
				Check: resource.ComposeTestCheckFunc(
					testCheckOceanSparkExists(&cluster, resourceName),
					testCheckOceanSparkAttributes(&cluster, oceanClusterID),
					resource.TestCheckResourceAttr(resourceName, "compute.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "compute.0.use_taints", "true"),
					resource.TestCheckResourceAttr(resourceName, "compute.0.create_vngs", "true"),
				),
			},
			{
				Config: createOceanSparkTerraform(&SparkClusterConfigMetadata{
					oceanClusterID: oceanClusterID,
					fieldsToAppend: testConfigWithComputeUpdate,
				}),
				Check: resource.ComposeTestCheckFunc(
					testCheckOceanSparkExists(&cluster, resourceName),
					testCheckOceanSparkAttributes(&cluster, oceanClusterID),
					resource.TestCheckResourceAttr(resourceName, "compute.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "compute.0.use_taints", "false"),
					resource.TestCheckResourceAttr(resourceName, "compute.0.create_vngs", "false"),
				),
			},
		},
	})
}

func TestAccSpotinstOceanSpark_withLogCollectionConfig(t *testing.T) {
	resourceName := createOceanSparkResourceName(oceanClusterID)

	var cluster spark.Cluster
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t, "aws") },
		Providers:    TestAccProviders,
		CheckDestroy: testOceanSparkAWSDestroy,

		Steps: []resource.TestStep{
			{
				Config: createOceanSparkTerraform(&SparkClusterConfigMetadata{
					oceanClusterID: oceanClusterID,
					fieldsToAppend: testConfigWithLogCollectionCreate,
				}),
				Check: resource.ComposeTestCheckFunc(
					testCheckOceanSparkExists(&cluster, resourceName),
					testCheckOceanSparkAttributes(&cluster, oceanClusterID),
					resource.TestCheckResourceAttr(resourceName, "log_collection.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "log_collection.0.collect_app_logs", "true"),
				),
			},
			{
				Config: createOceanSparkTerraform(&SparkClusterConfigMetadata{
					oceanClusterID: oceanClusterID,
					fieldsToAppend: testConfigWithLogCollectionUpdate,
				}),
				Check: resource.ComposeTestCheckFunc(
					testCheckOceanSparkExists(&cluster, resourceName),
					testCheckOceanSparkAttributes(&cluster, oceanClusterID),
					resource.TestCheckResourceAttr(resourceName, "log_collection.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "log_collection.0.collect_app_logs", "false"),
				),
			},
		},
	})
}

func TestAccSpotinstOceanSpark_withWorkspacesConfig(t *testing.T) {
	resourceName := createOceanSparkResourceName(oceanClusterID)

	var cluster spark.Cluster
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t, "aws") },
		Providers:    TestAccProviders,
		CheckDestroy: testOceanSparkAWSDestroy,

		Steps: []resource.TestStep{
			{
				Config: createOceanSparkTerraform(&SparkClusterConfigMetadata{
					oceanClusterID: oceanClusterID,
					fieldsToAppend: testConfigWithWorkspacesCreateWithOverrideSet,
				}),
				Check: resource.ComposeTestCheckFunc(
					testCheckOceanSparkExists(&cluster, resourceName),
					testCheckOceanSparkAttributes(&cluster, oceanClusterID),
					resource.TestCheckResourceAttr(resourceName, "workspaces.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "workspaces.0.storage_class_override", "gp2"),
				),
			},
			{
				Config: createOceanSparkTerraform(&SparkClusterConfigMetadata{
					oceanClusterID: oceanClusterID,
					fieldsToAppend: testConfigWithWorkspacesUpdateWithOverrideEmptyString,
				}),
				Check: resource.ComposeTestCheckFunc(
					testCheckOceanSparkExists(&cluster, resourceName),
					testCheckOceanSparkAttributes(&cluster, oceanClusterID),
					resource.TestCheckResourceAttr(resourceName, "workspaces.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "workspaces.0.storage_class_override", ""),
				),
			},
			{
				Config: createOceanSparkTerraform(&SparkClusterConfigMetadata{
					oceanClusterID: oceanClusterID,
					fieldsToAppend: testConfigWithWorkspacesUpdateWithOverrideSet,
				}),
				Check: resource.ComposeTestCheckFunc(
					testCheckOceanSparkExists(&cluster, resourceName),
					testCheckOceanSparkAttributes(&cluster, oceanClusterID),
					resource.TestCheckResourceAttr(resourceName, "workspaces.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "workspaces.0.storage_class_override", "pvc-ebs"),
				),
			},
			{
				Config: createOceanSparkTerraform(&SparkClusterConfigMetadata{
					oceanClusterID: oceanClusterID,
					fieldsToAppend: testConfigWithWorkspacesUpdateEmpty,
				}),
				Check: resource.ComposeTestCheckFunc(
					testCheckOceanSparkExists(&cluster, resourceName),
					testCheckOceanSparkAttributes(&cluster, oceanClusterID),
					resource.TestCheckResourceAttr(resourceName, "workspaces.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "workspaces.0.storage_class_override", ""),
				),
			},
		},
	})
}

func TestAccSpotinstOceanSpark_withSparkConfig(t *testing.T) {
	resourceName := createOceanSparkResourceName(oceanClusterID)

	var cluster spark.Cluster
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t, "aws") },
		Providers:    TestAccProviders,
		CheckDestroy: testOceanSparkAWSDestroy,

		Steps: []resource.TestStep{
			{
				Config: createOceanSparkTerraform(&SparkClusterConfigMetadata{
					oceanClusterID: oceanClusterID,
					fieldsToAppend: testConfigWithSparkConfigCreate,
				}),
				Check: resource.ComposeTestCheckFunc(
					testCheckOceanSparkExists(&cluster, resourceName),
					testCheckOceanSparkAttributes(&cluster, oceanClusterID),
					resource.TestCheckResourceAttr(resourceName, "spark.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "spark.0.additional_app_namespaces.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceName, "spark.0.additional_app_namespaces.*", "spark-apps-ns-1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "spark.0.additional_app_namespaces.*", "spark-apps-ns-2"),
				),
			},
			{
				Config: createOceanSparkTerraform(&SparkClusterConfigMetadata{
					oceanClusterID: oceanClusterID,
					fieldsToAppend: testConfigWithSparkConfigUpdate,
				}),
				Check: resource.ComposeTestCheckFunc(
					testCheckOceanSparkExists(&cluster, resourceName),
					testCheckOceanSparkAttributes(&cluster, oceanClusterID),
					resource.TestCheckResourceAttr(resourceName, "spark.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "spark.0.additional_app_namespaces.#", "3"),
					resource.TestCheckTypeSetElemAttr(resourceName, "spark.0.additional_app_namespaces.*", "spark-apps-ns-1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "spark.0.additional_app_namespaces.*", "spark-apps-ns-3"),
					resource.TestCheckTypeSetElemAttr(resourceName, "spark.0.additional_app_namespaces.*", "spark-apps-ns-4"),
				),
			},
			{
				Config: createOceanSparkTerraform(&SparkClusterConfigMetadata{
					oceanClusterID: oceanClusterID,
					fieldsToAppend: testConfigWithSparkConfigEmptyFields,
				}),
				Check: resource.ComposeTestCheckFunc(
					testCheckOceanSparkExists(&cluster, resourceName),
					testCheckOceanSparkAttributes(&cluster, oceanClusterID),
					resource.TestCheckResourceAttr(resourceName, "spark.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "spark.0.additional_app_namespaces.#", "0"),
				),
			},
		},
	})
}

func TestAccSpotinstOceanSpark_withSparkConfig_withDefaultNamespaceIncludedInAppNamespaceList_shouldIgnoreDefaultNamespace(t *testing.T) {
	resourceName := createOceanSparkResourceName(oceanClusterID)

	var cluster spark.Cluster
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t, "aws") },
		Providers:    TestAccProviders,
		CheckDestroy: testOceanSparkAWSDestroy,

		Steps: []resource.TestStep{
			{
				Config: createOceanSparkTerraform(&SparkClusterConfigMetadata{
					oceanClusterID: oceanClusterID,
					fieldsToAppend: testConfigWithSparkConfigWithDefaultNamespaceIncludedInNamespaceListCreate,
				}),
				Check: resource.ComposeTestCheckFunc(
					testCheckOceanSparkExists(&cluster, resourceName),
					testCheckOceanSparkAttributes(&cluster, oceanClusterID),
					resource.TestCheckResourceAttr(resourceName, "spark.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "spark.0.additional_app_namespaces.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceName, "spark.0.additional_app_namespaces.*", "spark-apps-ns-1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "spark.0.additional_app_namespaces.*", "spark-apps-ns-2"),
				),
			},
			{
				Config: createOceanSparkTerraform(&SparkClusterConfigMetadata{
					oceanClusterID: oceanClusterID,
					fieldsToAppend: testConfigWithSparkConfigWithDefaultNamespaceIncludedInNamespaceListUpdate,
				}),
				Check: resource.ComposeTestCheckFunc(
					testCheckOceanSparkExists(&cluster, resourceName),
					testCheckOceanSparkAttributes(&cluster, oceanClusterID),
					resource.TestCheckResourceAttr(resourceName, "spark.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "spark.0.additional_app_namespaces.#", "3"),
					resource.TestCheckTypeSetElemAttr(resourceName, "spark.0.additional_app_namespaces.*", "spark-apps-ns-1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "spark.0.additional_app_namespaces.*", "spark-apps-ns-3"),
					resource.TestCheckTypeSetElemAttr(resourceName, "spark.0.additional_app_namespaces.*", "spark-apps-ns-4"),
				),
			},
			{
				Config: createOceanSparkTerraform(&SparkClusterConfigMetadata{
					oceanClusterID: oceanClusterID,
					fieldsToAppend: testConfigWithSparkConfigWithDefaultNamespaceIncludedInNamespaceListUpdate2,
				}),
				Check: resource.ComposeTestCheckFunc(
					testCheckOceanSparkExists(&cluster, resourceName),
					testCheckOceanSparkAttributes(&cluster, oceanClusterID),
					resource.TestCheckResourceAttr(resourceName, "spark.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "spark.0.additional_app_namespaces.#", "3"),
					resource.TestCheckTypeSetElemAttr(resourceName, "spark.0.additional_app_namespaces.*", "spark-apps-ns-1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "spark.0.additional_app_namespaces.*", "spark-apps-ns-3"),
					resource.TestCheckTypeSetElemAttr(resourceName, "spark.0.additional_app_namespaces.*", "spark-apps-ns-4"),
				),
			},
			{
				Config: createOceanSparkTerraform(&SparkClusterConfigMetadata{
					oceanClusterID: oceanClusterID,
					fieldsToAppend: testConfigWithSparkConfigWithDefaultNamespaceIncludedInNamespaceListUpdate3,
				}),
				Check: resource.ComposeTestCheckFunc(
					testCheckOceanSparkExists(&cluster, resourceName),
					testCheckOceanSparkAttributes(&cluster, oceanClusterID),
					resource.TestCheckResourceAttr(resourceName, "spark.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "spark.0.additional_app_namespaces.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceName, "spark.0.additional_app_namespaces.*", "spark-apps-ns-5"),
					resource.TestCheckTypeSetElemAttr(resourceName, "spark.0.additional_app_namespaces.*", "spark-apps-ns-3"),
				),
			},
			{
				Config: createOceanSparkTerraform(&SparkClusterConfigMetadata{
					oceanClusterID: oceanClusterID,
					fieldsToAppend: testConfigWithSparkConfigWithDefaultNamespaceIncludedInNamespaceListEmptyFields,
				}),
				Check: resource.ComposeTestCheckFunc(
					testCheckOceanSparkExists(&cluster, resourceName),
					testCheckOceanSparkAttributes(&cluster, oceanClusterID),
					resource.TestCheckResourceAttr(resourceName, "spark.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "spark.0.additional_app_namespaces.#", "0"),
				),
			},
			/*
				This case does not work:
				spark {
					additional_app_namespaces = ["spark-apps"]
				 }

				i.e. where the list only contains one element, and that element is diff-suppressed.
				We get a list with length 1 (spark.0.additional_app_namespaces.# == 1), but it does not actually contain any elements.
				The test case also fails with the following:

				After applying this test step and performing a `terraform refresh`, the plan was not empty.
				~ spark {
				  ~ additional_app_namespaces = [
					  + null,
					]
				}

				We ignore this edge case here, and rely on the terraform module to prevent this.
				There we have validation that the default 'spark-apps' namespace should not be provided in the additional namespace list.
			*/
		},
	})
}

const testConfigWithIngressCreate = `
 ingress {

    service_annotations = {
     my-annotation-1 = "my-annotation-value-1"
     my-annotation-2 = "my-annotation-value-2"
	}

	controller {
	 managed = true
    }

	load_balancer {
	 managed = true
     target_group_arn = "some-test-target-group-arn"
	 service_annotations = {
      my-lb-service-annotation-1 = "my-lb-service-annotation-value-1"
      my-lb-service-annotation-2 = "my-lb-service-annotation-value-2"
	 }
	}

	custom_endpoint {
	 enabled = false
	 address = "valid-load.balancer"
	}

    private_link {
	 enabled = false
	 vpc_endpoint_service = "test-vpc-endpoint-service"
	}

 }
`

const testConfigWithIngressUpdate = `
 ingress {

    service_annotations = {
     my-new-annotation = "my-new-annotation-value"
     my-annotation-2 = "my-annotation-value-2-updated"
     my-annotation-3 = "my-annotation-value-3"
	}

	controller {
	 managed = true
    }

	load_balancer {
	 managed = false
     target_group_arn = "some-test-target-group-arn-active"
	 service_annotations = {
      my-lb-service-annotation-1 = "my-lb-service-annotation-value-1"
      my-lb-service-annotation-3 = "my-lb-service-annotation-value-3"
	 }
	}

	custom_endpoint {
	 enabled = true
	 address = "active-load.balancer"
	}

    private_link {
	 enabled = false
	 vpc_endpoint_service = "test-vpc-endpoint-service"
	}

 }
`

const testConfigWithIngressUpdate2 = `
 ingress {

    service_annotations = {
     my-new-annotation = "my-new-annotation-value"
     my-annotation-2 = "my-annotation-value-2-updated"
     my-annotation-3 = "my-annotation-value-3"
	}

	controller {
	 managed = false
    }

	load_balancer {
	 managed = false
     target_group_arn = "some-test-target-group-arn-inactive"
	}

	custom_endpoint {
	 enabled = false
	 address = "inactive-load.balancer"
	}

    private_link {
	 enabled = true
	 vpc_endpoint_service = "test-vpc-endpoint-service-active"
	}

 }
`

const testConfigWithIngressEmptyFields = `
 ingress {

	service_annotations = {}

	controller {

    }

	load_balancer {
	
	}

	custom_endpoint {
	
	}

    private_link {
	
	}

 }
`

const testConfigWithIngressEmptyFields2 = `
 ingress {

	service_annotations = {}

 }
`

const testConfigWithWebhookCreate = `
 webhook {

    use_host_network = false

	host_network_ports = [12345]

 }
`

const testConfigWithWebhookUpdate = `
 webhook {

    use_host_network = true

	host_network_ports = [12345, 54321]

 }
`

const testConfigWithWebhookEmptyFields = `
 webhook {

	use_host_network = false

	host_network_ports = []

 }
`

const testConfigWithSparkConfigCreate = `
 spark {

	additional_app_namespaces = ["spark-apps-ns-1","spark-apps-ns-2"]

 }
`

const testConfigWithSparkConfigUpdate = `
 spark {

	additional_app_namespaces = ["spark-apps-ns-1","spark-apps-ns-3","spark-apps-ns-4"]

 }
`

const testConfigWithSparkConfigEmptyFields = `
 spark {

	additional_app_namespaces = []

 }
`

const testConfigWithSparkConfigWithDefaultNamespaceIncludedInNamespaceListCreate = `
 spark {

	additional_app_namespaces = ["spark-apps","spark-apps-ns-1","spark-apps-ns-2"]

 }
`

const testConfigWithSparkConfigWithDefaultNamespaceIncludedInNamespaceListUpdate = `
 spark {

	additional_app_namespaces = ["spark-apps","spark-apps-ns-1","spark-apps-ns-3","spark-apps-ns-4"]

 }
`

const testConfigWithSparkConfigWithDefaultNamespaceIncludedInNamespaceListUpdate2 = `
 spark {

	additional_app_namespaces = ["spark-apps-ns-1","spark-apps-ns-3","spark-apps-ns-4"]

 }
`

const testConfigWithSparkConfigWithDefaultNamespaceIncludedInNamespaceListUpdate3 = `
 spark {

	additional_app_namespaces = ["spark-apps-ns-5","spark-apps-ns-3","spark-apps"]

 }
`

const testConfigWithSparkConfigWithDefaultNamespaceIncludedInNamespaceListEmptyFields = `
 spark {

	additional_app_namespaces = []

 }
`

const testConfigWithComputeCreate = `
 compute {

    use_taints = true

	create_vngs = true

 }
`

const testConfigWithComputeUpdate = `
 compute {

    use_taints = false

	create_vngs = false

 }
`

const testConfigWithLogCollectionCreate = `
 log_collection {

    collect_app_logs = true

 }
`

const testConfigWithLogCollectionUpdate = `
 log_collection {

    collect_app_logs = false

 }
`

const testConfigWithWorkspacesCreateWithOverrideSet = `
 workspaces {
  
    storage_class_override = "gp2"

 }
`

const testConfigWithWorkspacesUpdateWithOverrideSet = `
 workspaces {
  
    storage_class_override = "pvc-ebs"

 }
`

const testConfigWithWorkspacesUpdateWithOverrideEmptyString = `
 workspaces {
  
    storage_class_override = ""

 }
`

const testConfigWithWorkspacesUpdateEmpty = `
 workspaces {

 }
`
