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

func TestAccPrometheusOperatorPrometheusRule_basic(t *testing.T) {
	var pr po_types.PrometheusRule
	name := fmt.Sprintf("tf-acc-test:%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))
	namespace := "monitoring"

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "po_prometheus_rule.test",
		Providers:     testAccProviders,
		CheckDestroy:  testAccPrometheusOperatorPrometheusRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPrometheusOperatorPrometheusRuleConfig_basic(name, namespace),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccPrometheusOperatorPrometheusRuleExists("po_prometheus_rule.test", &pr),
					resource.TestCheckResourceAttr("po_prometheus_rule.test", "metadata.0.name", name),
					resource.TestCheckResourceAttr("po_prometheus_rule.test", "metadata.0.namespace", namespace),
					resource.TestCheckResourceAttr("po_prometheus_rule.test", "metadata.0.labels.%", "2"),
					resource.TestCheckResourceAttr("po_prometheus_rule.test", "metadata.0.labels.prometheus", "k8s"),
					resource.TestCheckResourceAttr("po_prometheus_rule.test", "metadata.0.labels.role", "alert-rules"),
					resource.TestCheckResourceAttrSet("po_prometheus_rule.test", "metadata.0.generation"),
					resource.TestCheckResourceAttrSet("po_prometheus_rule.test", "metadata.0.resource_version"),
					resource.TestCheckResourceAttrSet("po_prometheus_rule.test", "metadata.0.self_link"),
					resource.TestCheckResourceAttrSet("po_prometheus_rule.test", "metadata.0.uid"),
					resource.TestCheckResourceAttr("po_prometheus_rule.test", "spec.0.groups.0.rules.#", "1"),
					resource.TestCheckResourceAttr("po_prometheus_rule.test", "spec.0.groups.0.rules.0.alert", "Watchdog"),
					resource.TestCheckResourceAttr("po_prometheus_rule.test", "spec.0.groups.0.rules.0.annotations.%", "1"),
					resource.TestCheckResourceAttr("po_prometheus_rule.test", "spec.0.groups.0.rules.0.annotations.message", "This is an alert meant ..."),
					resource.TestCheckResourceAttr("po_prometheus_rule.test", "spec.0.groups.0.rules.0.expr", "vector(1)"),
					resource.TestCheckResourceAttr("po_prometheus_rule.test", "spec.0.groups.0.rules.0.labels.%", "1"),
					resource.TestCheckResourceAttr("po_prometheus_rule.test", "spec.0.groups.0.rules.0.labels.severity", "none"),
				),
			},
		},
	})
}

func TestAccPrometheusOperatorPrometheusRule_importBasic(t *testing.T) {
	resourceName := "po_prometheus_rule.test"
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


func testAccPrometheusOperatorPrometheusRuleConfig_basic(name, namespace string) string {
	return fmt.Sprintf(`
resource "po_prometheus_rule" "test" {
  metadata {
    name ="%s"
    namespace = "%s"
    labels = {
      prometheus = "k8s"
      role = "alert-rules"
    }
  }
  spec {
    groups { 
      rules {
        alert = "Watchdog"
        annotations = {
          message = "This is an alert meant ..."
        }
        expr = "vector(1)"
        labels = {
          severity = "none"
        }
      }
    }
  }`, name, namespace)
}

func testAccPrometheusOperatorPrometheusRuleExists(n string, obj *po_types.PrometheusRule) resource.TestCheckFunc {
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

		out, err := conn.PrometheusRules(namespace).Get(name, meta_v1.GetOptions{})
		if err != nil {
			return err
		}

		*obj = *out
		return nil
	}
}

func testAccPrometheusOperatorPrometheusRuleDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*KubeClientsets).MonitoringClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "po_prometheus_rule" {
			continue
		}

		namespace, name, err := idParts(rs.Primary.ID)
		if err != nil {
			return err
		}

		resp, err := conn.PrometheusRules(namespace).Get(name, meta_v1.GetOptions{})
		if err == nil {
			if resp.Name == rs.Primary.ID {
				return fmt.Errorf("Service still exists: %s", rs.Primary.ID)
			}
		}
	}
	return nil
}

