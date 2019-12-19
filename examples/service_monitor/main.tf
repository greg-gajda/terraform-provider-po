variable "namespace" { default = "monitoring" }

resource "po_service_monitor" "kube_apiserver" {
  metadata {
    name = "kube-apiserver"
    namespace = var.namespace
    labels = {
      "k8s-app" = "apiserver"
    }
  }
  spec {
    endpoints {
      bearer_token_file = "/var/run/secrets/kubernetes.io/serviceaccount/token"
      interval = "30s"
      metric_relabelings {
        action = "drop"
        regex = "etcd_(debugging|disk|request|server).*"
        source_labels = ["__name__"]
      }
      metric_relabelings {
        action = "drop"
        regex = "apiserver_admission_controller_admission_latencies_seconds_.*"
        source_labels = ["__name__"]
      }
      metric_relabelings {
        action = "drop"
        regex = "apiserver_admission_step_admission_latencies_seconds_.*"
        source_labels = ["__name__"]
      }
      port = "https"
      scheme = "https"
      tls_config {
        ca_file = "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt"
        server_name = "kubernetes"
      }
    }
    job_label = "component"
    namespace_selector {
      match_names = ["default"]
    }
    selector {
      match_labels = {
        component = "apiserver"
        provider = "kubernetes"
      }
    }
  }
}

resource "po_service_monitor" "kube_controller_manager" {
  metadata {
    name = "kube-controller-manager"
    namespace = var.namespace
    labels = {
      "k8s-app" = "kube-controller-manager"
    }
  }
  spec {
    endpoints {
      port = "http-metrics"
      interval = "15s"
      metric_relabelings {
        action = "drop"
        regex = "etcd_(debugging|disk|request|server).*"
        source_labels = ["__name__"]
      }
    }
    job_label = "k8s-app"
    namespace_selector {
      match_names = ["kube-system"]
    }
    selector {
      match_labels = {
        "k8s-app" = "kube-controller-manager"
      }
    }
  }
}

resource "po_service_monitor" "kube_scheduler" {
  metadata {
    name = "kube-scheduler"
    namespace = var.namespace
    labels = {
      "k8s-app" = "kube-scheduler"
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
        "k8s-app" = "kube-scheduler"
      }
    }
  }
}
