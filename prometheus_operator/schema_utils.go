package prometheus_operator

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func RuleGroupSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "Name of Prometheus Rule.",
			Required:    true,
		},
		"interval": {
			Type:        schema.TypeString,
			Description: "",
			Optional:    true,
		},
		"rules": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Rules describes an alerting or recording rules",
			Elem: &schema.Resource{
				Schema: RuleSchema(),
			},
		},
	}
}

func RuleSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"record": {
			Type:        schema.TypeString,
			Description: "(empty)",
			Optional:    true,
		},
		"alert": {
			Type:        schema.TypeString,
			Description: "(empty)",
			Optional:    true,
		},
		"expr": {
			Type:        schema.TypeString,
			Description: "(empty)",
			Required:    true,
		},
		"for": {
			Type:        schema.TypeString,
			Description: "(empty)",
			Optional:    true,
		},
		"labels": {
			Type:         schema.TypeMap,
			Description:  "(empty)",
			Optional:     true,
			Elem:         &schema.Schema{Type: schema.TypeString},
			ValidateFunc: validateLabels,
		},
		"annotations": {
			Type:         schema.TypeMap,
			Description:  "(empty)",
			Optional:     true,
			Elem:         &schema.Schema{Type: schema.TypeString},
		},
	}
}

func AlertingSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"alertmanagers": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "AlertmanagerEndpoints Prometheus should fire alerts against.",
			Elem: &schema.Resource{
				Schema: AlertmanagerEndpointsSchema(),
			},
		},
	}
}

func AlertmanagerEndpointsSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"namespace": {
			Type:        schema.TypeString,
			Description: "Namespace of Endpoints object.",
			Required:    true,
		},
		"name": {
			Type:        schema.TypeString,
			Description: "Name of Endpoints object in Namespace.",
			Required:    true,
		},
		"port": {
			Type:        schema.TypeString,
			Description: "Port the Alertmanager API is exposed on.",
			Required:    true,
		},
		"path_prefix": {
			Type:        schema.TypeString,
			Description: "Prefix for the HTTP path alerts are pushed to.",
			Optional:    true,
		},
		"scheme": {
			Type:        schema.TypeString,
			Description: "Scheme to use when firing alerts.",
			Optional:    true,
		},
		"bearer_token_file": {
			Type:        schema.TypeString,
			Description: "BearerTokenFile to read from filesystem to use when authenticating to Alertmanager.",
			Optional:    true,
		},
		"api_version": {
			Type:        schema.TypeString,
			Description: "Version of the Alertmanager API that Prometheus uses to send alerts. It can be 'v1' or 'v2'.",
			Optional:    true,
		},
		"tls_config": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "TLS Config to use for alertmanager connection.",
			Elem: &schema.Resource{
				Schema: TLSConfigSchema(),
			},
		},
	}
}

func NamespaceSelectorSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"any": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Boolean describing whether all namespaces are selected in contrast to a list restricting them.",
		},
		"match_names": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Description: "List of namespace names",
		},
	}
}

func TLSConfigSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"ca_file": {
			Type:        schema.TypeString,
			Description: "Path to the CA cert in the Prometheus container to use for the targets.",
			Optional:    true,
		},
		"ca": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Struct containing the CA cert to use for the targets.",
			Elem: &schema.Resource{
				Schema: SecretOrConfigMapSchema(),
			},
		},
		"cert_file": {
			Type:        schema.TypeString,
			Description: "Path to the client cert file in the Prometheus container for the targets.",
			Optional:    true,
		},
		"cert": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Struct containing the client cert file for the targets.",
			Elem: &schema.Resource{
				Schema: SecretOrConfigMapSchema(),
			},
		},
		"key_file": {
			Type:        schema.TypeString,
			Description: "Path to the client key file in the Prometheus container for the targets.",
			Optional:    true,
		},
		"key_secret": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Secret containing data to use for the targets",
			Elem: &schema.Resource{
				Schema: SecretKeySelectorSchema(),
			},
		},
		"server_name": {
			Type:        schema.TypeString,
			Description: "Used to verify the hostname for the targets.",
			Optional:    true,
		},
		"insecure_skip_verify": {
			Type:        schema.TypeBool,
			Description: "Disable target certificate validation.",
			Optional:    true,
		},
	}
}

func SecretOrConfigMapSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"secret": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Secret containing data to use for the targets",
			Elem: &schema.Resource{
				Schema: SecretKeySelectorSchema(),
			},
		},
		"config_map": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Selects a key of a ConfigMap.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"key": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "The key to select.",
					},
					"name": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Name of the referent. More info: http://kubernetes.io/docs/user-guide/identifiers#names",
					},
				},
			},
		},
	}
}

func SecretKeySelectorSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"key": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The key of the secret to select from. Must be a valid secret key.",
		},
		"name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Name of the referent. More info: http://kubernetes.io/docs/user-guide/identifiers#names",
		},
	}
}

func BasicAuthSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"username": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "The secret in the service monitor namespace that contains the username for authentication.",
			Elem: &schema.Resource{
				Schema: SecretKeySelectorSchema(),
			},
		},
		"password": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "The secret in the service monitor namespace that contains the password for authentication.",
			Elem: &schema.Resource{
				Schema: SecretKeySelectorSchema(),
			},
		},
	}
}

func RelabelConfigSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"separator": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Separator placed between concatenated source label values. default is ';'.",
		},
		"target_label": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Label to which the resulting value is written in a replace action. It is mandatory for replace actions. Regex capture groups are available.",
		},
		"regex": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Regular expression against which the extracted value is matched. Default is '(.*)'",
		},
		"modulus": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Modulus to take of the hash of the source label values.",
		},
		"replacement": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Replacement value against which a regex replace is performed if the regular expression matches. Regex capture groups are available. Default is '$1'",
		},
		"action": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Action to perform based on regex matching. Default is 'replace'",
		},
		"source_labels": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Description: "The source labels select values from existing labels. More info: https://github.com/coreos/prometheus-operator/blob/master/Documentation/api.md#relabelconfig",
		},
	}
}

func EndpointSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"port": {
			Type:        schema.TypeString,
			Description: "Name of the port this endpoint refers to. Mutually exclusive with targetPort. More info: https://github.com/coreos/prometheus-operator/blob/master/Documentation/api.md#servicemonitorspec",
			Optional:    true,
		},
		"path": {
			Type:        schema.TypeString,
			Description: "HTTP path to scrape for metrics. 	More info https://github.com/coreos/prometheus-operator/blob/master/Documentation/api.md#servicemonitorspec",
			Optional:    true,
		},
		"scheme": {
			Type:        schema.TypeString,
			Description: "HTTP scheme to use for scraping.",
			Optional:    true,
		},
		"interval": {
			Type:        schema.TypeString,
			Description: "Interval at which metrics should be scraped.",
			Optional:    true,
		},
		"scrape_timeout": {
			Type:        schema.TypeString,
			Description: "Timeout after which the scrape is ended.",
			Optional:    true,
		},
		"tls_config": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "TLS configuration to use when scraping the endpoint. More info: https://github.com/coreos/prometheus-operator/blob/master/Documentation/api.md#tlsconfig",
			Elem: &schema.Resource{
				Schema: TLSConfigSchema(),
			},
		},
		"bearer_token_file": {
			Type:        schema.TypeString,
			Description: "File to read bearer token for scraping targets.",
			Optional:    true,
		},
		"bearer_token_secret": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Secret to mount to read bearer token for scraping targets. The secret needs to be in the same namespace as the service monitor and accessible by the Prometheus Operator.",
			Elem: &schema.Resource{
				Schema: SecretKeySelectorSchema(),
			},
		},
		"honor_labels": {
			Type:        schema.TypeBool,
			Description: "HonorLabels chooses the metric's labels on collisions with target labels.",
			Optional:    true,
		},
		"honor_timestamps": {
			Type:        schema.TypeBool,
			Description: "HonorTimestamps controls whether Prometheus respects the timestamps present in scraped data.",
			Optional:    true,
			Default:     true,
		},
		"basic_auth": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "BasicAuth allow an endpoint to authenticate over basic authentication More info: https://prometheus.io/docs/operating/configuration/#endpoints",
			Elem: &schema.Resource{
				Schema: BasicAuthSchema(),
			},
		},
		"metric_relabelings": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "MetricRelabelConfigs to apply to samples before ingestion.",
			Elem: &schema.Resource{
				Schema: RelabelConfigSchema(),
			},
		},
		"relabelings": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "RelabelConfigs to apply to samples before scraping. More info: https://prometheus.io/docs/prometheus/latest/configuration/configuration/#relabel_config",
			Elem: &schema.Resource{
				Schema: RelabelConfigSchema(),
			},
		},
		"proxy_url": {
			Type:        schema.TypeString,
			Description: "ProxyURL eg http://proxyserver:2195 Directs scrapes to proxy through this endpoint.",
			Optional:    true,
		},
	}
}

func TolerationSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"effect": {
			Type:         schema.TypeString,
			Description:  "Effect indicates the taint effect to match. Empty means match all taint effects. When specified, allowed values are NoSchedule, PreferNoSchedule and NoExecute.",
			Optional:     true,
			ValidateFunc: validation.StringInSlice([]string{"NoSchedule", "PreferNoSchedule", "NoExecute"}, false),
		},
		"key": {
			Type:        schema.TypeString,
			Description: "Key is the taint key that the toleration applies to. Empty means match all taint keys. If the key is empty, operator must be Exists; this combination means to match all values and all keys.",
			Optional:    true,
		},
		"operator": {
			Type:         schema.TypeString,
			Description:  "Operator represents a key's relationship to the value. Valid operators are Exists and Equal. Defaults to Equal. Exists is equivalent to wildcard for value, so that a pod can tolerate all taints of a particular category.",
			Default:      "Equal",
			Optional:     true,
			ValidateFunc: validation.StringInSlice([]string{"Exists", "Equal"}, false),
		},
		"toleration_seconds": {
			// Use TypeString to allow an "unspecified" value,
			Type:         schema.TypeString,
			Description:  "TolerationSeconds represents the period of time the toleration (which must be of effect NoExecute, otherwise this field is ignored) tolerates the taint. By default, it is not set, which means tolerate the taint forever (do not evict). Zero and negative values will be treated as 0 (evict immediately) by the system.",
			Optional:     true,
			ValidateFunc: validateTypeStringNullableInt,
		},
		"value": {
			Type:        schema.TypeString,
			Description: "Value is the taint value the toleration matches to. If the operator is Exists, the value should be empty, otherwise just a regular string.",
			Optional:    true,
		},
	}
}

func SecurityContextSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"fs_group": {
			Type:        schema.TypeInt,
			Description: "A special supplemental group that applies to all containers in a pod. Some volume types allow the Kubelet to change the ownership of that volume to be owned by the pod: 1. The owning GID will be the FSGroup 2. The setgid bit is set (new files created in the volume will be owned by FSGroup) 3. The permission bits are OR'd with rw-rw---- If unset, the Kubelet will not modify the ownership and permissions of any volume.",
			Optional:    true,
		},
		"run_as_group": {
			Type:        schema.TypeInt,
			Description: "The GID to run the entrypoint of the container process. Uses runtime default if unset. May also be set in SecurityContext. If set in both SecurityContext and PodSecurityContext, the value specified in SecurityContext takes precedence for that container.",
			Optional:    true,
		},
		"run_as_non_root": {
			Type:        schema.TypeBool,
			Description: "Indicates that the container must run as a non-root user. If true, the Kubelet will validate the image at runtime to ensure that it does not run as UID 0 (root) and fail to start the container if it does. If unset or false, no such validation will be performed. May also be set in SecurityContext. If set in both SecurityContext and PodSecurityContext, the value specified in SecurityContext takes precedence.",
			Optional:    true,
		},
		"run_as_user": {
			Type:        schema.TypeInt,
			Description: "The UID to run the entrypoint of the container process. Defaults to user specified in image metadata if unspecified. May also be set in SecurityContext. If set in both SecurityContext and PodSecurityContext, the value specified in SecurityContext takes precedence for that container.",
			Optional:    true,
		},
		"se_linux_options": {
			Type:        schema.TypeList,
			Description: "The SELinux context to be applied to all containers. If unspecified, the container runtime will allocate a random SELinux context for each container. May also be set in SecurityContext. If set in both SecurityContext and PodSecurityContext, the value specified in SecurityContext takes precedence for that container.",
			Optional:    true,
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: seLinuxOptionsField(),
			},
		},
		"supplemental_groups": {
			Type:        schema.TypeSet,
			Description: "A list of groups applied to the first process run in each container, in addition to the container's primary GID. If unspecified, no groups will be added to any container.",
			Optional:    true,
			Elem: &schema.Schema{
				Type: schema.TypeInt,
			},
		},
	}
}