package prometheus_operator

import (
	"fmt"
	po_types "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	po_v1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	pkgApi "k8s.io/apimachinery/pkg/types"
	"log"
)

func resourcePOServiceMonitor() *schema.Resource {
	return &schema.Resource{
		Create: resourcePOServiceMonitorCreate,
		Read:   resourcePOServiceMonitorRead,
		Exists: resourcePOServiceMonitorExists,
		Update: resourcePOServiceMonitorUpdate,
		Delete: resourcePOServiceMonitorDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"metadata": namespacedMetadataSchema("service monitor", true),
			"spec": {
				Type:        schema.TypeList,
				Description: "Spec defines the specification of the desired behavior of the deployment. More info: https://github.com/coreos/prometheus-operator/blob/master/Documentation/api.md#servicemonitorspec",
				Required:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"job_label": {
							Type:        schema.TypeString,
							Description: "The label to use to retrieve the job name from. More info: https://github.com/coreos/prometheus-operator/blob/master/Documentation/api.md#servicemonitorspec",
							Optional:    true,
						},
						"target_labels": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "TargetLabels transfers labels on the Kubernetes Service onto the target. More info: https://github.com/coreos/prometheus-operator/blob/master/Documentation/api.md#servicemonitorspec",
						},
						"pod_target_labels": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "PodTargetLabels transfers labels on the Kubernetes Pod onto the target. More info: https://github.com/coreos/prometheus-operator/blob/master/Documentation/api.md#servicemonitorspec",
						},
						"endpoints": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "A list of endpoints allowed as part of this ServiceMonitor. More info: https://github.com/coreos/prometheus-operator/blob/master/Documentation/api.md#servicemonitorspec",
							Elem: &schema.Resource{
								Schema: EndpointSchema(),
							},
						},
						"selector": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Selector to select Endpoints objects.",
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: labelSelectorFields(true),
							},
						},
						"namespace_selector": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Selector to select which namespaces the Endpoints objects are discovered from.",
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: NamespaceSelectorSchema(),
							},
						},
						"sample_limit": {
							Type:        schema.TypeInt,
							Description: "SampleLimit defines per-scrape limit on number of scraped samples that will be accepted. More info: https://github.com/coreos/prometheus-operator/blob/master/Documentation/api.md#servicemonitorspec",
							Optional:    true,
						},
					},
				},
			},
		},
	}
}

func resourcePOServiceMonitorCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*KubeClientsets).MonitoringClient
	metadata := expandMetadata(d.Get("metadata").([]interface{}))

	spec, err := expandServiceMonitorSpec(d.Get("spec").([]interface{}))
	if err != nil {
		return err
	}

	monitor := po_types.ServiceMonitor{
		ObjectMeta: metadata,
		Spec:       *spec,
	}

	log.Printf("[INFO] Creating ServiceMonitor custom resource: %#v", monitor)
	out, err := conn.ServiceMonitors(metadata.Namespace).Create(&monitor)
	if err != nil {
		return fmt.Errorf("Failed to create ServiceMonitor: %s", err)
	}

	log.Printf("[INFO] Submitted new ServiceMonitor custom resource: %#v", out)

	d.SetId(buildId(out.ObjectMeta))

	return resourcePOServiceMonitorRead(d, meta)
}

func resourcePOServiceMonitorExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	conn := meta.(*KubeClientsets).MonitoringClient
	namespace, name, err := idParts(d.Id())
	if err != nil {
		return false, err
	}

	log.Printf("[INFO] Checking ServiceMonitor custom resource %s", name)
	_, err = conn.ServiceMonitors(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return false, nil
		}
		log.Printf("[DEBUG] Received error: %#v", err)
	}
	return true, err
}

func resourcePOServiceMonitorRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*KubeClientsets).MonitoringClient
	namespace, name, err := idParts(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Reading ServiceMonitor custom resource %s", name)
	am, err := conn.ServiceMonitors(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		switch {
		case errors.IsNotFound(err):
			log.Printf("[DEBUG] ServiceMonitor %q was not found in Namespace %q - removing from state!", namespace, name)
			d.SetId("")
			return nil
		default:
			log.Printf("[DEBUG] Error reading ServiceMonitor: %#v", err)
			return err
		}
	}
	log.Printf("[INFO] Received ServiceMonitor: %#v", am)

	if d.Set("metadata", flattenMetadata(am.ObjectMeta, d)) != nil {
		return fmt.Errorf("Error setting `metadata`: %+v", err)
	}
	spec, err := flattenServiceMonitorSpec(am.Spec, d)

	d.Set("spec", spec)
	if err != nil {
		return fmt.Errorf("Failed to set ServiceMonitor spec: %s", err)
	}
	return nil
}

func resourcePOServiceMonitorUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*KubeClientsets).MonitoringClient
	namespace, name, err := idParts(d.Id())
	if err != nil {
		return err
	}
	ops := patchMetadata("metadata.0.", "/metadata/", d)

	if d.HasChange("spec") {
		log.Println("[TRACE] ServiceMonitor.Spec has changes")
		spec, err := expandServiceMonitorSpec(d.Get("spec").([]interface{}))
		if err != nil {
			return err
		}
		ops = append(ops, replace(spec))
	}

	data, err := ops.MarshalJSON()
	if err != nil {
		return fmt.Errorf("Failed to marshal update operations for ServiceMonitor: %s", err)
	}
	log.Printf("[INFO] Updating ServiceMonitor %q: %v", name, string(data))
	out, err := conn.ServiceMonitors(namespace).Patch(name, pkgApi.JSONPatchType, data)
	if err != nil {
		return fmt.Errorf("Failed to update ServiceMonitor: %s", err)
	}
	log.Printf("[INFO] Submitted updated ServiceMonitor: %#v", out)

	return resourcePOServiceMonitorRead(d, meta)
}

func resourcePOServiceMonitorDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*KubeClientsets).MonitoringClient
	namespace, name, err := idParts(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleting ServiceMonitor: %q", name)
	err = conn.ServiceMonitors(namespace).Delete(name, &metav1.DeleteOptions{})
	if err != nil {
		return err
	}

	log.Printf("[INFO] ServiceMonitor %s deleted", name)

	d.SetId("")

	return nil
}

func expandServiceMonitorSpec(sm []interface{}) (*po_types.ServiceMonitorSpec, error) {
	obj := &po_types.ServiceMonitorSpec{}
	if len(sm) == 0 || sm[0] == nil {
		return obj, nil
	}
	in := sm[0].(map[string]interface{})

	obj.JobLabel = in["job_label"].(string)
	obj.SampleLimit = uint64(in["sample_limit"].(int))
	if tl, ok := in["target_labels"].([]interface{}); ok {
		obj.TargetLabels = expandStringSlice(tl)
	}
	if ptl, ok := in["pod_target_labels"].([]interface{}); ok {
		obj.PodTargetLabels = expandStringSlice(ptl)
	}
	if v, ok := in["endpoints"].([]interface{}); ok && len(v) > 0 {
		endpoints, err := expandEndpoints(v)
		if err != nil {
			return obj, err
		}
		obj.Endpoints = endpoints
	}
	if s, ok := in["selector"].([]interface{}); ok && len(s) > 0 {
		selector := expandLabelSelector(s)
		obj.Selector = *selector
	}
	if ns, ok := in["namespace_selector"].([]interface{}); ok && len(ns) > 0 {
		selector, err := expandNamespaceSelector(ns)
		if err != nil {
			return obj, err
		}
		obj.NamespaceSelector = *selector
	}
	return obj, nil
}

func flattenServiceMonitorSpec(spec po_v1.ServiceMonitorSpec, d *schema.ResourceData) ([]interface{}, error) {
	att := make(map[string]interface{})

	if spec.JobLabel != "" {
		att["job_label"] = spec.JobLabel
	}
	if len(spec.TargetLabels) > 0 {
		att["target_labels"] = spec.TargetLabels
	}
	if len(spec.PodTargetLabels) > 0 {
		att["pod_target_labels"] = spec.PodTargetLabels
	}
	att["sample_limit"] = int(spec.SampleLimit)

	endpoints, err := flattenEndpoints(spec.Endpoints)
	if err != nil {
		return nil, err
	}
	att["endpoints"] = endpoints
	att["selector"] = flattenLabelSelector(&spec.Selector)
	att["namespace_selector"] = flattenNamespaceSelector(&spec.NamespaceSelector)

	return []interface{}{att}, nil
}
