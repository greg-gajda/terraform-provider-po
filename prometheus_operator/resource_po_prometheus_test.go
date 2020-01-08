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

func TestAccPrometheusOperatorPrometheus_basic(t *testing.T) {
	var p po_types.Prometheus
	name := fmt.Sprintf("tf-acc-test:%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))
	namespace := "monitoring"

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "po_prometheus.test",
		Providers:     testAccProviders,
		CheckDestroy:  testAccPrometheusOperatorPrometheusDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPrometheusOperatorPrometheusConfig_basic(name, namespace),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccPrometheusOperatorPrometheusExists("po_prometheus.test", &p),
					resource.TestCheckResourceAttr("po_prometheus.test", "metadata.0.name", name),
					resource.TestCheckResourceAttr("po_prometheus.test", "metadata.0.name", namespace),
					resource.TestCheckResourceAttr("po_prometheus.test", "metadata.0.labels.%", "1"),
					resource.TestCheckResourceAttr("po_prometheus.test", "metadata.0.labels.prometheus", "k8s"),
					resource.TestCheckResourceAttrSet("po_prometheus.test", "metadata.0.generation"),
					resource.TestCheckResourceAttrSet("po_prometheus.test", "metadata.0.resource_version"),
					resource.TestCheckResourceAttrSet("po_prometheus.test", "metadata.0.self_link"),
					resource.TestCheckResourceAttrSet("po_prometheus.test", "metadata.0.uid"),
					resource.TestCheckResourceAttr("po_prometheus.test", "spec.0.alerting.0.alertmanagers.0.name", "alertmanager-main"),
					resource.TestCheckResourceAttr("po_prometheus.test", "spec.0.alerting.0.alertmanagers.0.namespace", namespace),
					resource.TestCheckResourceAttr("po_prometheus.test", "spec.0.alerting.0.alertmanagers.0.port", "web"),
					resource.TestCheckResourceAttr("po_prometheus.test", "spec.0.retention", "30d"),
					resource.TestCheckResourceAttr("po_prometheus.test", "spec.0.resources.0.requests.0.memory", "400Mi"),
					resource.TestCheckResourceAttr("po_prometheus.test", "spec.0.rule_selector.0.match_labels.0.prometheus", "k8s"),
					resource.TestCheckResourceAttr("po_prometheus.test", "spec.0.rule_selector.0.match_labels.0.role", "alert-rules"),
					resource.TestCheckResourceAttr("po_prometheus.test", "spec.0.replicas", "1"),
					resource.TestCheckResourceAttr("po_prometheus.test", "spec.0.base_image", "quay.io/prometheus/prometheus"),
					resource.TestCheckResourceAttr("po_prometheus.test", "spec.0.service_account_name", "prometheus-k8s"),
					resource.TestCheckResourceAttr("po_prometheus.test", "spec.0.version", "v2.14.0"),
					resource.TestCheckResourceAttr("po_prometheus.test", "spec.0.security_context.0.fs_group", "100"),
					resource.TestCheckResourceAttr("po_prometheus.test", "spec.0.security_context.0.run_as_non_root", "true"),
					resource.TestCheckResourceAttr("po_prometheus.test", "spec.0.security_context.0.run_as_user", "1000"),
					resource.TestCheckResourceAttr("po_prometheus.test", "spec.0.security_context.0.fs_group", "2000"),
					resource.TestCheckResourceAttr("po_prometheus.test", "spec.0.security_context.0.run_as_group", "3000"),
				),
			},
		},
	})
}

func TestAccPrometheusOperatorPrometheus_importBasic(t *testing.T) {
	resourceName := "po_prometheus.test"
	name := fmt.Sprintf("tf-acc-test:%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))
	namespace := "monitoring"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPrometheusOperatorPrometheusDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPrometheusOperatorPrometheusConfig_basic(name, namespace),
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


func testAccPrometheusOperatorPrometheusConfig_basic(name, namespace string) string {
	return fmt.Sprintf(`
resource "po_prometheus" "test" {
  metadata {
    name = "%[1]s"
    namespace = "%[2]s"
    labels = {
      prometheus = "k8s"
    }
  }
  spec {
    alerting {
      alertmanagers {
        name = "alertmanager-main"
        namespace = "%[2]s"
        port = "web"
      }
    }
    retention = "30d"
    resources {
      requests {
        memory = "400Mi"
      }
    }
    rule_selector {
      match_labels = {
        prometheus = "k8s"
        role = "alert-rules"
      }
    }
    service_monitor_selector {}
    replicas = 1
    base_image = "quay.io/prometheus/prometheus"
    service_account_name = "prometheus-k8s"
    version = "v2.14.0"
    security_context {
      run_as_non_root = true
      run_as_user = 1000
      fs_group = 2000
      run_as_group = 3000
    }
  }
}
`, name, namespace)
}

func testAccPrometheusOperatorPrometheusExists(n string, obj *po_types.Prometheus) resource.TestCheckFunc {
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

		out, err := conn.Prometheuses(namespace).Get(name, meta_v1.GetOptions{})
		if err != nil {
			return err
		}

		*obj = *out
		return nil
	}
}


func testAccPrometheusOperatorPrometheusDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*KubeClientsets).MonitoringClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "po_prometheus" {
			continue
		}

		namespace, name, err := idParts(rs.Primary.ID)
		if err != nil {
			return err
		}

		resp, err := conn.Prometheuses(namespace).Get(name, meta_v1.GetOptions{})
		if err == nil {
			if resp.Name == rs.Primary.ID {
				return fmt.Errorf("Service still exists: %s", rs.Primary.ID)
			}
		}
	}
	return nil
}

