package prometheus_operator

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	po_types "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccPrometheusOperatorServiceMonitor_basic(t *testing.T) {
	var sm po_types.ServiceMonitor
	name := fmt.Sprintf("tf-acc-test:%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))
	namespace := "monitoring"

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "po_service_monitor.test",
		Providers:     testAccProviders,
		CheckDestroy:  testAccPrometheusOperatorServiceMonitorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPrometheusOperatorServiceMonitorConfig_basic(name, namespace),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccPrometheusOperatorServiceMonitorExists("po_service_monitor.test", &sm),
					resource.TestCheckResourceAttr("po_service_monitor.test", "metadata.0.name", name),
					resource.TestCheckResourceAttr("po_service_monitor.test", "metadata.0.name", "monitoring"),
					resource.TestCheckResourceAttr("po_service_monitor.test", "metadata.0.labels.%", "1"),
					resource.TestCheckResourceAttr("po_service_monitor.test", "metadata.0.labels.k8s-app", name),
					resource.TestCheckResourceAttrSet("po_service_monitor.test", "metadata.0.generation"),
					resource.TestCheckResourceAttrSet("po_service_monitor.test", "metadata.0.resource_version"),
					resource.TestCheckResourceAttrSet("po_service_monitor.test", "metadata.0.self_link"),
					resource.TestCheckResourceAttrSet("po_service_monitor.test", "metadata.0.uid"),
					resource.TestCheckResourceAttr("po_service_monitor.test", "spec.0.endpoints.0.port", "http-metrics"),
					resource.TestCheckResourceAttr("po_service_monitor.test", "spec.0.endpoints.0.interval", "30s"),
					resource.TestCheckResourceAttr("po_service_monitor.test", "spec.0.job_label", "k8s-app"),
					resource.TestCheckResourceAttr("po_service_monitor.test", "spec.0.namespace_selector.0.match_names.#", "1"),
					resource.TestCheckResourceAttr("po_service_monitor.test", "spec.0.namespace_selector.0.match_names.0", "kube-system"),
					resource.TestCheckResourceAttr("po_service_monitor.test", "spec.0.selector.0.match_labels.0.k8s-app", name),
				),
			},
		},
	})
}

func TestAccPrometheusOperatorServiceMonitor_importBasic(t *testing.T) {
	resourceName := "po_service_monitor.test"
	name := fmt.Sprintf("tf-acc-test:%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))
	namespace := "monitoring"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPrometheusOperatorServiceMonitorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPrometheusOperatorServiceMonitorConfig_basic(name, namespace),
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


func testAccPrometheusOperatorServiceMonitorConfig_basic(name, namespace string) string {
	return fmt.Sprintf(`
resource "po_service_monitor" "test" {
  metadata {
    name = "%[1]s"
    namespace = "%[2]s"
    labels = {
      "k8s-app" = "%[1]s"
    }
  }
  spec {
    endpoints {
      port = "http-metrics"
      interval = "30s"
    }
    job_label = "k8s-app"
    namespace_selector {
      match_names = ["kube-system"]
    }
    selector {
      match_labels = {
        "k8s-app" = "%[1]s"
      }
    }
  }
}`, name, namespace)
}

func testAccPrometheusOperatorServiceMonitorExists(n string, obj *po_types.ServiceMonitor) resource.TestCheckFunc {
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

		out, err := conn.ServiceMonitors(namespace).Get(context.TODO(), name, meta_v1.GetOptions{})
		if err != nil {
			return err
		}

		*obj = *out
		return nil
	}
}


func testAccPrometheusOperatorServiceMonitorDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*KubeClientsets).MonitoringClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "po_service_monitor" {
			continue
		}

		namespace, name, err := idParts(rs.Primary.ID)
		if err != nil {
			return err
		}

		resp, err := conn.ServiceMonitors(namespace).Get(context.TODO(), name, meta_v1.GetOptions{})
		if err == nil {
			if resp.Name == rs.Primary.ID {
				return fmt.Errorf("Service still exists: %s", rs.Primary.ID)
			}
		}
	}
	return nil
}
