module github.com/terraform-providers/terraform-provider-po

require (
	contrib.go.opencensus.io/exporter/ocagent v0.6.0 // indirect
	github.com/coreos/prometheus-operator v0.34.0
	github.com/google/go-cmp v0.3.1
	github.com/gophercloud/gophercloud v0.1.0 // indirect
	github.com/hashicorp/go-version v1.2.0
	github.com/hashicorp/terraform-plugin-sdk v1.3.0
	github.com/mitchellh/go-homedir v1.1.0
	github.com/robfig/cron v1.2.0
	github.com/terraform-providers/terraform-provider-aws v2.32.0+incompatible
	github.com/terraform-providers/terraform-provider-google v2.17.0+incompatible
	github.com/terraform-providers/terraform-provider-kubernetes v1.10.0
	k8s.io/api v0.0.0-20191025225708-5524a3672fbb
	k8s.io/apimachinery v0.0.0-20191025225532-af6325b3a843
	k8s.io/client-go v12.0.0+incompatible
	k8s.io/kube-aggregator v0.0.0-20191025230902-aa872b06629d
)

replace github.com/terraform-providers/terraform-provider-kubernetes v1.10.0 => ./kubernetes
