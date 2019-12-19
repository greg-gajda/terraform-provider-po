package prometheus_operator

import (
	"fmt"
	po_types "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	v1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	pkgApi "k8s.io/apimachinery/pkg/types"
	"log"
)

func resourcePOPrometheus() *schema.Resource {
	return &schema.Resource{
		Create: resourcePOPrometheusCreate,
		Read:   resourcePOPrometheusRead,
		Exists: resourcePOPrometheusExists,
		Update: resourcePOPrometheusUpdate,
		Delete: resourcePOPrometheusDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"metadata": namespacedMetadataSchema("prometheus", true),
			"spec": {
				Type:        schema.TypeList,
				Description: "Spec defines the specification of the desired behavior of the deployment. More info: https://github.com/coreos/prometheus-operator/blob/master/Documentation/api.md#alertmanager",
				Required:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"base_image": {
							Type:        schema.TypeString,
							Description: "Base image that is used to deploy pods, without tag. More info: https://github.com/coreos/prometheus-operator/blob/master/Documentation/api.md#alertmanager",
							Optional:    true,
							ForceNew:    true,
							Default:     "quay.io/prometheus/prometheus",
						},
						"image": {
							Type:        schema.TypeString,
							Description: "Image if specified has precedence over baseImage, tag and sha combinations. Specifying the version is still necessary to ensure the Prometheus Operator knows what version of Alertmanager is being configured.",
							Optional:    true,
							ForceNew:    true,
						},
						"service_monitor_selector": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "ServiceMonitors to be selected for target discovery.",
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: labelSelectorFields(true),
							},
						},
						"service_monitor_namespace_selector": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Namespaces to be selected for ServiceMonitor discovery. If nil, only check own namespace.",
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: labelSelectorFields(true),
							},
						},
						"pod_monitor_selector": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Experimental PodMonitors to be selected for target discovery.",
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: labelSelectorFields(true),
							},
						},
						"pod_monitor_namespace_selector": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "podMonitorNamespaceSelector",
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: labelSelectorFields(true),
							},
						},
						"secrets": {
							Type:        schema.TypeList,
							Description: "Secrets is a list of Secrets in the same namespace as the Prometheus object, which shall be mounted into the Prometheus Pods. The Secrets are mounted into /etc/prometheus/secrets/.",
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"config_maps": {
							Type:        schema.TypeList,
							Description: "ConfigMaps is a list of ConfigMaps in the same namespace as the Prometheus object, which shall be mounted into the Prometheus Pods. The ConfigMaps are mounted into /etc/prometheus/configmaps/.",
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"external_url": {
							Type:        schema.TypeString,
							Description: "The external URL the Prometheus instances will be available under. This is necessary to generate correct URLs. This is necessary if Prometheus is not served from root of a DNS name. More info: https://github.com/coreos/prometheus-operator/blob/master/Documentation/api.md#alertmanager",
							Optional:    true,
						},
						"service_account_name": {
							Type:        schema.TypeString,
							Description: "ServiceAccountName is the name of the ServiceAccount to use to run the Alertmanager Pods. More info: https://github.com/coreos/prometheus-operator/blob/master/Documentation/api.md#alertmanager",
							Optional:    true,
						},
						"paused": {
							Type:        schema.TypeBool,
							Description: "If set to true all actions on the underlaying managed objects are not goint to be performed, except for delete actions. More info: https://github.com/coreos/prometheus-operator/blob/master/Documentation/api.md#alertmanager",
							Optional:    true,
							Default:     false,
						},
						"replicas": {
							Type:        schema.TypeInt,
							Description: "The number of desired replicas. Defaults to 2.",
							Optional:    true,
							Default:     2,
						},
						"version": {
							Type:        schema.TypeString,
							Description: "Prometheus version to be used.",
							Optional:    true,
							ForceNew:    true,
						},
						"tag": {
							Type:        schema.TypeString,
							Description: "Tag of Alertmanager container image to be deployed. Defaults to the value of version. Version is ignored if Tag is set",
							Optional:    true,
							ForceNew:    true,
						},
						"sha": {
							Type:        schema.TypeString,
							Description: "SHA of Alertmanager container image to be deployed. Defaults to the value of version. Similar to a tag, but the SHA explicitly deploys an immutable container image. Version and Tag are ignored if SHA is set.",
							Optional:    true,
							ForceNew:    true,
						},
						"retention": {
							Type:        schema.TypeString,
							Description: "Time duration Prometheus shall retain data for. Default is '24h', and must match the regular expression [0-9]+(ms|s|m|h|d|w|y) (milliseconds seconds minutes hours days weeks years)",
							Optional:    true,
							Default:     "24h",
						},
						"retention_size": {
							Type:        schema.TypeString,
							Description: "Maximum amount of disk space used by blocks.",
							Optional:    true,
						},
						"port_name": {
							Type:        schema.TypeString,
							Description: "Port name used for the pods and governing service. More info: https://github.com/coreos/prometheus-operator/blob/master/Documentation/api.md#alertmanager",
							Optional:    true,
						},
						"priority_class_name": {
							Type:        schema.TypeString,
							Description: "Priority class assigned to the alertmanager Pods. More info: https://github.com/coreos/prometheus-operator/blob/master/Documentation/api.md#alertmanager",
							Optional:    true,
						},
						"listen_local": {
							Type:        schema.TypeBool,
							Description: "ListenLocal makes the Alertmanager server listen on loopback. More info: https://github.com/coreos/prometheus-operator/blob/master/Documentation/api.md#alertmanager",
							Optional:    true,
							Default:     false,
						},
						"container": {
							Type:        schema.TypeList,
							Optional:    true,
							ForceNew:    true,
							Description: "Containers allows injecting additional containers. More info: https://github.com/coreos/prometheus-operator/blob/master/Documentation/api.md#alertmanager",
							Elem: &schema.Resource{
								Schema: containerFields(true, false),
							},
						},
						"init_container": {
							Type:        schema.TypeList,
							Optional:    true,
							ForceNew:    true,
							Description: "InitContainers allows adding initContainers to the pod definition. More info: https://github.com/coreos/prometheus-operator/blob/master/Documentation/api.md#alertmanager",
							Elem: &schema.Resource{
								Schema: containerFields(true, true),
							},
						},
						"node_selector": {
							Type:        schema.TypeMap,
							Optional:    true,
							Description: "Define which Nodes the Alertmanager Pods are scheduled on. More info: https://github.com/coreos/prometheus-operator/blob/master/Documentation/api.md#alertmanager",
						},
						"security_context": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "SecurityContext holds pod-level security attributes and common container settings. Optional: Defaults to empty. . More info: http://releases.k8s.io/HEAD/docs/design/security_context.md",
							Elem: &schema.Resource{
								Schema: SecurityContextSchema(),
							},
						},
						"resources": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Computed:    true,
							Description: "Define resources requests and limits for single Alertmanager Pods. More info: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.11/#resourcerequirements-v1-core",
							Elem: &schema.Resource{
								Schema: resourcesField(),
							},
						},
						"volume": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Description: "List of volumes that can be mounted by containers belonging to the pod. More info: http://kubernetes.io/docs/user-guide/volumes",
							Elem:        volumeSchema(),
						},
						"toleration": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "If specified, the pod's toleration. Optional: Defaults to empty",
							Elem: &schema.Resource{
								Schema: TolerationSchema(),
							},
						},
						"alerting": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Define details regarding alerting.",
							Elem: &schema.Resource{
								Schema: AlertingSchema(),
							},
						},
						"rule_selector": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "A selector to select which PrometheusRules to mount for loading alerting rules from. More info: https://github.com/coreos/prometheus-operator/blob/master/Documentation/api.md#prometheusspec",
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: labelSelectorFields(true),
							},
						},
						"rule_namespace_selector": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Namespaces to be selected for PrometheusRules discovery. If unspecified, only the same namespace as the Prometheus object is in is used.",
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: labelSelectorFields(true),
							},
						},
					},
				},
			},
		},
	}
}

func resourcePOPrometheusCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*KubeClientsets).MonitoringClient
	metadata := expandMetadata(d.Get("metadata").([]interface{}))

	spec, err := expandPrometheusSpec(d.Get("spec").([]interface{}))
	if err != nil {
		return err
	}

	prometheus := po_types.Prometheus{
		ObjectMeta: metadata,
		Spec:       *spec,
	}

	log.Printf("[INFO] Creating Prometheus custom resource: %#v", prometheus)
	out, err := conn.Prometheuses(metadata.Namespace).Create(&prometheus)
	if err != nil {
		return fmt.Errorf("Failed to create Prometheus: %s", err)
	}

	log.Printf("[INFO] Submitted new Prometheus custom resource: %#v", out)

	d.SetId(buildId(out.ObjectMeta))

	return resourcePOPrometheusRead(d, meta)
}

func resourcePOPrometheusExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	conn := meta.(*KubeClientsets).MonitoringClient
	namespace, name, err := idParts(d.Id())
	if err != nil {
		return false, err
	}

	log.Printf("[INFO] Checking Prometheus custom resource %s", name)
	_, err = conn.Prometheuses(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return false, nil
		}
		log.Printf("[DEBUG] Received error: %#v", err)
	}
	return true, err
}

func resourcePOPrometheusRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*KubeClientsets).MonitoringClient
	namespace, name, err := idParts(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Reading Prometheus %s", name)
	am, err := conn.Prometheuses(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		switch {
		case errors.IsNotFound(err):
			log.Printf("[DEBUG] Prometheus %q was not found in Namespace %q - removing from state!", namespace, name)
			d.SetId("")
			return nil
		default:
			log.Printf("[DEBUG] Error reading Prometheus: %#v", err)
			return err
		}
	}
	log.Printf("[INFO] Received Prometheus: %#v", am)

	if d.Set("metadata", flattenMetadata(am.ObjectMeta, d)) != nil {
		return fmt.Errorf("Error setting `metadata`: %+v", err)
	}
	spec, err := flattenPrometheusSpec(am.Spec)

	d.Set("spec", spec)
	if err != nil {
		return fmt.Errorf("Failed to set Prometheus spec: %s", err)
	}
	return nil
}

func resourcePOPrometheusUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*KubeClientsets).MonitoringClient
	namespace, name, err := idParts(d.Id())
	if err != nil {
		return err
	}
	ops := patchMetadata("metadata.0.", "/metadata/", d)

	if d.HasChange("spec") {
		log.Println("[TRACE] Prometheus.Spec has changes")
		spec, err := expandPrometheusSpec(d.Get("spec").([]interface{}))
		if err != nil {
			return err
		}
		ops = append(ops, replace(spec))
	}

	data, err := ops.MarshalJSON()
	if err != nil {
		return fmt.Errorf("Failed to marshal update operations for Prometheus: %s", err)
	}
	log.Printf("[INFO] Updating Prometheus %q: %v", name, string(data))
	out, err := conn.Prometheuses(namespace).Patch(name, pkgApi.JSONPatchType, data)
	if err != nil {
		return fmt.Errorf("Failed to update Prometheus: %s", err)
	}
	log.Printf("[INFO] Submitted updated Prometheus: %#v", out)

	return resourcePOPrometheusRead(d, meta)
}

func resourcePOPrometheusDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*KubeClientsets).MonitoringClient
	namespace, name, err := idParts(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleting Prometheus: %q", name)
	err = conn.Prometheuses(namespace).Delete(name, &metav1.DeleteOptions{})
	if err != nil {
		return err
	}

	log.Printf("[INFO] Prometheus %s deleted", name)

	d.SetId("")

	return nil
}

func expandPrometheusSpec(prometheus []interface{}) (*po_types.PrometheusSpec, error) {
	obj := &po_types.PrometheusSpec{}
	if len(prometheus) == 0 || prometheus[0] == nil {
		return obj, nil
	}
	in := prometheus[0].(map[string]interface{})

	obj.BaseImage = in["base_image"].(string)
	if im, ok := in["image"]; ok {
		obj.Image = ptrToString(im.(string))
	}
	if sms, ok := in["service_monitor_selector"].([]interface{}); ok && len(sms) > 0 {
		selector := expandLabelSelector(sms)
		obj.ServiceMonitorSelector = selector
	}
	if smns, ok := in["service_monitor_namespace_selector"].([]interface{}); ok && len(smns) > 0 {
		selector := expandLabelSelector(smns)
		obj.ServiceMonitorNamespaceSelector = selector
	}
	if pms, ok := in["pod_monitor_selector"].([]interface{}); ok && len(pms) > 0 {
		selector := expandLabelSelector(pms)
		obj.PodMonitorSelector = selector
	}
	if pmns, ok := in["pod_monitor_namespace_selector"].([]interface{}); ok && len(pmns) > 0 {
		selector := expandLabelSelector(pmns)
		obj.PodMonitorNamespaceSelector = selector
	}
	if rs, ok := in["rule_selector"].([]interface{}); ok && len(rs) > 0 {
		selector := expandLabelSelector(rs)
		obj.RuleSelector = selector
	}
	if rns, ok := in["rule_namespace_selector"].([]interface{}); ok && len(rns) > 0 {
		selector := expandLabelSelector(rns)
		obj.RuleNamespaceSelector = selector
	}
	if sec, ok := in["secrets"]; ok {
		obj.Secrets = expandStringSlice(sec.([]interface{}))
	}
	if cm, ok := in["config_maps"]; ok {
		obj.ConfigMaps = expandStringSlice(cm.([]interface{}))
	}
	if ret, ok := in["retention"]; ok {
		obj.Retention = ret.(string)
	}
	if rets, ok := in["retention_size"]; ok {
		obj.RetentionSize = rets.(string)
	}
	obj.ExternalURL = in["external_url"].(string)
	obj.ServiceAccountName = in["service_account_name"].(string)
	obj.Paused = in["paused"].(bool)
	obj.Replicas = ptrToInt32(int32(in["replicas"].(int)))
	obj.Version  = in["version"].(string)
	obj.Tag  = in["tag"].(string)
	obj.SHA  = in["sha"].(string)
	obj.PriorityClassName = in["priority_class_name"].(string)
	obj.PortName = in["port_name"].(string)
	obj.ListenLocal = in["listen_local"].(bool)

	if v, ok := in["container"].([]interface{}); ok && len(v) > 0 {
		cs, err := expandContainers(v)
		if err != nil {
			return obj, err
		}
		obj.Containers = cs
	}
	if v, ok := in["init_container"].([]interface{}); ok && len(v) > 0 {
		cs, err := expandContainers(v)
		if err != nil {
			return obj, err
		}
		obj.InitContainers = cs
	}
	if v, ok := in["node_selector"].(map[string]interface{}); ok {
		nodeSelectors := make(map[string]string)
		for k, v := range v {
			if val, ok := v.(string); ok {
				nodeSelectors[k] = val
			}
		}
		obj.NodeSelector = nodeSelectors
	}
	if v, ok := in["resources"].([]interface{}); ok && len(v) > 0 {
		crr, err := expandContainerResourceRequirements(v)
		if err != nil {
			return obj, err
		}
		obj.Resources = *crr
	}
	if v, ok := in["security_context"].([]interface{}); ok && len(v) > 0 {
		obj.SecurityContext = expandPodSecurityContext(v)
	}
	if v, ok := in["toleration"].([]interface{}); ok && len(v) > 0 {
		ts, err := expandTolerations(v)
		if err != nil {
			return obj, err
		}
		for _, t := range ts {
			obj.Tolerations = append(obj.Tolerations, *t)
		}
	}
	if v, ok := in["volume"].([]interface{}); ok && len(v) > 0 {
		cs, err := expandVolumes(v)
		if err != nil {
			return obj, err
		}
		obj.Volumes = cs
	}
	if v, ok := in["alerting"].([]interface{}); ok && len(v) > 0 {
		a, err := expandAlertingSpec(v)
		if err != nil {
			return obj, err
		}
		obj.Alerting = a
	}

	return obj, nil
}

func flattenPrometheusSpec(spec v1.PrometheusSpec/*, d *schema.ResourceData*/) ([]interface{}, error) {
	att := make(map[string]interface{})

	if spec.BaseImage != "" {
		att["base_image"] = spec.BaseImage
	}
	if spec.Image != nil {
		att["image"] = *spec.Image
	}
	if spec.ServiceMonitorSelector != nil {
		att["service_monitor_selector"] = flattenLabelSelector(spec.ServiceMonitorSelector)
	}
	if spec.ServiceMonitorNamespaceSelector != nil {
		att["service_monitor_namespace_selector"] = flattenLabelSelector(spec.ServiceMonitorNamespaceSelector)
	}
	if spec.PodMonitorSelector != nil {
		att["pod_monitor_selector"] = flattenLabelSelector(spec.PodMonitorSelector)
	}
	if spec.PodMonitorNamespaceSelector != nil {
		att["pod_monitor_namespace_selector"] = flattenLabelSelector(spec.PodMonitorNamespaceSelector)
	}
	if spec.RuleSelector != nil {
		att["rule_selector"] = flattenLabelSelector(spec.RuleSelector)
	}
	if spec.RuleNamespaceSelector != nil {
		att["rule_namespace_selector"] = flattenLabelSelector(spec.RuleNamespaceSelector)
	}
	if len(spec.Secrets) > 0 {
		att["secrets"] = spec.Secrets
	}
	if len(spec.ConfigMaps) > 0 {
		att["config_maps"] = spec.ConfigMaps
	}
	if spec.Retention != "" {
		att["retention"] = spec.Retention
	}
	if spec.RetentionSize != "" {
		att["retention_size"] = spec.RetentionSize
	}
	if spec.ExternalURL != "" {
		att["external_url"] = spec.ExternalURL
	}
	if spec.ServiceAccountName != "" {
		att["service_account_name"] = spec.ServiceAccountName
	}
	att["paused"] = spec.Paused
	if spec.Replicas != nil {
		att["replicas"] = *spec.Replicas
	}
	if spec.Version != "" {
		att["version"] = spec.Version
	}
	if spec.Tag != "" {
		att["tag"] = spec.Tag
	}
	if spec.SHA != "" {
		att["sha"] = spec.SHA
	}
	att["listen_local"] = spec.ListenLocal
	if spec.PriorityClassName != "" {
		att["priority_class_name"] = spec.PriorityClassName
	}
	if spec.PortName != "" {
		att["port_name"] = spec.PortName
	}
	if spec.SecurityContext != nil {
		att["security_context"] = flattenPodSecurityContext(spec.SecurityContext)
	}
	containers, err := flattenContainers(spec.Containers)
	if err != nil {
		return nil, err
	}
	att["container"] = containers
	initContainers, err := flattenContainers(spec.InitContainers)
	if err != nil {
		return nil, err
	}
	att["init_container"] = initContainers

	endpoints, err := flattenAlertingSpec(spec.Alerting)
	if err != nil {
		return nil, err
	}
	att["alerting"] = endpoints

	return []interface{}{att}, nil
}


