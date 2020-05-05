# This is a terraform provider for Prometheus operator deployments on Kubernetes

Why it may be useful to anyone?

Prometheus [Operator](https://coreos.com/operators/) makes the management of Prometheus based monitoring stack easier.
But as a consequence of that simplicity, new [CustomResourceDefinitions](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/#customresourcedefinitions) appear in your k8s cluster.

Quoting official docs:
The Prometheus Operator introduces additional resources in Kubernetes to declare the desired state of a Prometheus and Alertmanager cluster as well as the Prometheus configuration. The resources it introduces are:
* **Prometheus**
* **Alertmanager**
* **ServiceMonitor**

but there's also **PrometheusRule** and **PodMonitor**.

If you simply deliver your monitoring stack by running 
```
kubectl apply -f bundle.yaml
```
from the official [Prometheus Operator](https://github.com/coreos/prometheus-operator) Github repo, 
or by using [Helm chart](https://github.com/helm/charts/tree/master/stable/prometheus-operator),
you wouldn't find anything interesting in this repository.

But, if you're using [Terraform](https://www.terraform.io/) to manage your infrastructure, 
then it can be quite useful, as it delivers custom [Provider Plugin](https://www.terraform.io/docs/plugins/provider.html) for Prometheus Operator custom resources.

Content of /kubernetes folder is taken from [official Terraform Kubernetes provider](https://github.com/terraform-providers/terraform-provider-kubernetes).

To acquire a binary for your OS, simply clone the project and run `go build` command.
