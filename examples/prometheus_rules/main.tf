variable "namespace" { default = "monitoring" }
variable "k8s_cluster" {}

provider "po" {
  config_context_cluster = var.k8s_cluster
}

resource "po_prometheus_rule" "prometheus_rules" {
  metadata {
    name ="prometheus-k8s-rules"
    namespace = var.namespace
    labels = {
      prometheus = "k8s"
      role = "alert-rules"
    }
  }
  spec {
    groups {
      name = "general.rules"
      rules {
        alert = "TargetDown"
        annotations = {
          message = <<EOF
            '{{ printf "%.4g" $value }}% of the {{ $labels.job }} targets in {{ $labels.namespace }} namespace are down.'
          EOF
        }
        expr = <<EOF
          100 * (count(up == 0) BY (job, namespace, service) / count(up) BY (job, namespace, service)) > 10
        EOF
        for = "10m"
        labels = {
          severity = "warning"
        }
      }
      rules {
        alert = "Watchdog"
        annotations = {
          message = <<EOF
            This is an alert meant to ensure that the entire alerting pipeline is functional.
            This alert is always firing, therefore it should always be firing in Alertmanager
            and always fire against a receiver. There are integrations with various notification
            mechanisms that send a notification when this alert is not firing. For example the
            "DeadMansSnitch" integration in PagerDuty.
          EOF
        }
        expr = "vector(1)"
        labels = {
          severity = "none"
        }
      }
    }
  }
}