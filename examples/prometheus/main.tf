variable "k8s_cluster" {}
variable "namespace" { default = "monitoring" }
variable "no_of_replicas" { default = 2 }
variable "prometheus_version" { default = "v2.14.0" }
variable "storage_retention" { default = "30d" }

provider "po" {
  config_context_cluster = var.k8s_cluster
}

resource "po_prometheus" "prometheus" {
  metadata {
    name = "k8s"
    namespace = var.namespace
    labels = {
      prometheus = "k8s"
    }
  }
  spec {
    alerting {
      alertmanagers {
        name = "alertmanager-main"
        namespace = var.namespace
        port = "web"
      }
    }
    retention = var.storage_retention
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
    replicas = var.no_of_replicas
    base_image = "quay.io/prometheus/prometheus"
    service_account_name = "prometheus-k8s"
    version = var.prometheus_version
    security_context {
      run_as_non_root = true
      run_as_user = 1000
      fs_group = 2000
      run_as_group = 3000
    }
  }
}