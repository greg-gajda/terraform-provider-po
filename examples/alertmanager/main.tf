variable "k8s_cluster" {}
variable "namespace" { default = "monitoring" }
variable "no_of_replicas" { default = 3 }
variable "alertmanager_version" { default = "v0.19.0" }

provider "po" {
  config_context_cluster = var.k8s_cluster
}

resource "po_alertmanager" "alertmanager" {
  metadata {
    name = "main"
    namespace = var.namespace
    labels = {
      alertmanager = "main"
    }
  }
  spec {
    replicas = var.no_of_replicas
    base_image = "quay.io/prometheus/alertmanager"
    service_account_name = "alertmanager-main"
    version = var.alertmanager_version
    security_context {
      run_as_non_root = true
      run_as_user = 1000
      fs_group = 2000
      run_as_group = 3000
    }
  }
}
