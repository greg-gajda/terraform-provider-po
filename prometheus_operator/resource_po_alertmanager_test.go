package prometheus_operator

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	po_types "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccPrometheusOperatorAlertmanager_basic(t *testing.T) {
	var am po_types.Alertmanager
	name := fmt.Sprintf("tf-acc-test:%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "po_alertmanager.test",
		Providers:     testAccProviders,
		CheckDestroy:  testAccPrometheusOperatorAlertmanagerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPrometheusOperatorAlertmanagerConfig_basic(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccPrometheusOperatorAlertmanagerExists("po_alertmanager.test", &am),
					resource.TestCheckResourceAttr("po_alertmanager.test", "metadata.0.name", name),
					resource.TestCheckResourceAttr("po_alertmanager.test", "metadata.0.namespace", "monitoring"),
					resource.TestCheckResourceAttr("po_alertmanager.test", "metadata.0.labels.%", "1"),
					resource.TestCheckResourceAttr("po_alertmanager.test", "metadata.0.labels.alertmanager", "main"),
					resource.TestCheckResourceAttrSet("po_alertmanager.test", "metadata.0.generation"),
					resource.TestCheckResourceAttrSet("po_alertmanager.test", "metadata.0.resource_version"),
					resource.TestCheckResourceAttrSet("po_alertmanager.test", "metadata.0.self_link"),
					resource.TestCheckResourceAttrSet("po_alertmanager.test", "metadata.0.uid"),
					resource.TestCheckResourceAttr("po_alertmanager.test", "spec.0.replicas", "1"),
					resource.TestCheckResourceAttr("po_alertmanager.test", "spec.0.base_image", "quay.io/prometheus/alertmanager"),
					resource.TestCheckResourceAttr("po_alertmanager.test", "spec.0.service_account_name", "alertmanager-main"),
					resource.TestCheckResourceAttr("po_alertmanager.test", "spec.0.version", "v0.19.0"),
					resource.TestCheckResourceAttr("po_alertmanager.test", "spec.0.security_context.0.fs_group", "100"),
					resource.TestCheckResourceAttr("po_alertmanager.test", "spec.0.security_context.0.run_as_non_root", "true"),
					resource.TestCheckResourceAttr("po_alertmanager.test", "spec.0.security_context.0.run_as_user", "1000"),
					resource.TestCheckResourceAttr("po_alertmanager.test", "spec.0.security_context.0.fs_group", "2000"),
					resource.TestCheckResourceAttr("po_alertmanager.test", "spec.0.security_context.0.run_as_group", "3000"),
				),
			},
		},
	})
}

func TestAccPrometheusOperatorAlertmanager_importBasic(t *testing.T) {
	resourceName := "po_alertmanager.test"
	name := fmt.Sprintf("tf-acc-test:%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPrometheusOperatorAlertmanagerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPrometheusOperatorAlertmanagerConfig_basic(name),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"metadata.0.resource_version"},
			},
		},
	})
}


func testAccPrometheusOperatorAlertmanagerConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "po_alertmanager" "test" {
  metadata {
    name = "%s"
    namespace = "monitoring"
    labels = {
      alertmanager = "main"
    }
  }
  spec {
    replicas = 1
    base_image = "quay.io/prometheus/alertmanager"
    service_account_name = "alertmanager-main"
    version = "v0.19.0"
    security_context {
      run_as_non_root = true
      run_as_user = 1000
      fs_group = 2000
      run_as_group = 3000
    }
  }
}
`, name)
}

func testAccPrometheusOperatorAlertmanagerExists(n string, obj *po_types.Alertmanager) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		conn := testAccProvider.Meta().(*KubeClientsets).MonitoringClient

		namespace, name, err := idParts(rs.Primary.ID)
		if err != nil {
			return err
		}

		out, err := conn.Alertmanagers(namespace).Get(name, meta_v1.GetOptions{})
		if err != nil {
			return err
		}

		*obj = *out
		return nil
	}
}


func testAccPrometheusOperatorAlertmanagerDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*KubeClientsets).MonitoringClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "po_alertmanager" {
			continue
		}

		namespace, name, err := idParts(rs.Primary.ID)
		if err != nil {
			return err
		}

		resp, err := conn.Alertmanagers(namespace).Get(name, meta_v1.GetOptions{})
		if err == nil {
			if resp.Name == rs.Primary.ID {
				return fmt.Errorf("Service still exists: %s", rs.Primary.ID)
			}
		}
	}
	return nil
}
