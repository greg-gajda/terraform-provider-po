package prometheus_operator

import (
	"fmt"
	po_types "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	po_v1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	pkgApi "k8s.io/apimachinery/pkg/types"
	"log"
	"context"
)

func resourcePOPrometheusRule() *schema.Resource {
	return &schema.Resource{
		Create: resourcePOPrometheusRuleCreate,
		Read:   resourcePOPrometheusRuleRead,
		Exists: resourcePOPrometheusRuleExists,
		Update: resourcePOPrometheusRuleUpdate,
		Delete: resourcePOPrometheusRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"metadata": namespacedMetadataSchema("prometheus rule", true),
			"spec": {
				Type:        schema.TypeList,
				Description: "Spec defines the specification of the desired behavior of the deployment. More info https://github.com/prometheus-operator/prometheus-operator/blob/master/Documentation/api.md#prometheusrulespec",
				Required:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"groups": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Content of Prometheus rule file. More info https://github.com/prometheus-operator/prometheus-operator/blob/master/Documentation/api.md#rulegroup",
							Elem: &schema.Resource{
								Schema: RuleGroupSchema(),
							},
						},
					},
				},
			},
		},
	}
}

func resourcePOPrometheusRuleCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*KubeClientsets).MonitoringClient
	metadata := expandMetadata(d.Get("metadata").([]interface{}))

	spec, err := expandPrometheusRuleSpec(d.Get("spec").([]interface{}))
	if err != nil {
		return err
	}

	monitor := po_types.PrometheusRule{
		ObjectMeta: metadata,
		Spec:       *spec,
	}

	log.Printf("[INFO] Creating PrometheusRule custom resource: %#v", monitor)
	out, err := conn.PrometheusRules(metadata.Namespace).Create(context.Background(), &monitor, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("Failed to create PrometheusRule: %s", err)
	}

	log.Printf("[INFO] Submitted new PrometheusRule custom resource: %#v", out)

	d.SetId(buildId(out.ObjectMeta))

	return resourcePOPrometheusRuleRead(d, meta)
}

func resourcePOPrometheusRuleExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	conn := meta.(*KubeClientsets).MonitoringClient
	namespace, name, err := idParts(d.Id())
	if err != nil {
		return false, err
	}

	log.Printf("[INFO] Checking PrometheusRule custom resource %s", name)
	_, err = conn.PrometheusRules(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return false, nil
		}
		log.Printf("[DEBUG] Received error: %#v", err)
	}
	return true, err
}

func resourcePOPrometheusRuleRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*KubeClientsets).MonitoringClient
	namespace, name, err := idParts(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Reading PrometheusRule custom resource %s", name)
	am, err := conn.PrometheusRules(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		switch {
		case errors.IsNotFound(err):
			log.Printf("[DEBUG] PrometheusRule %q was not found in Namespace %q - removing from state!", namespace, name)
			d.SetId("")
			return nil
		default:
			log.Printf("[DEBUG] Error reading PrometheusRule: %#v", err)
			return err
		}
	}
	log.Printf("[INFO] Received PrometheusRule: %#v", am)

	if d.Set("metadata", flattenMetadata(am.ObjectMeta, d)) != nil {
		return fmt.Errorf("Error setting `metadata`: %+v", err)
	}
	spec, err := flattenPrometheusRuleSpec(am.Spec, d)

	d.Set("spec", spec)
	if err != nil {
		return fmt.Errorf("Failed to set PrometheusRule spec: %s", err)
	}
	return nil
}

func resourcePOPrometheusRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*KubeClientsets).MonitoringClient
	namespace, name, err := idParts(d.Id())
	if err != nil {
		return err
	}
	ops := patchMetadata("metadata.0.", "/metadata/", d)

	if d.HasChange("spec") {
		log.Println("[TRACE] PrometheusRule.Spec has changes")
		spec, err := expandPrometheusRuleSpec(d.Get("spec").([]interface{}))
		if err != nil {
			return err
		}
		ops = append(ops, replace(spec))
	}

	data, err := ops.MarshalJSON()
	if err != nil {
		return fmt.Errorf("Failed to marshal update operations for PrometheusRule: %s", err)
	}
	log.Printf("[INFO] Updating PrometheusRule %q: %v", name, string(data))
	out, err := conn.PrometheusRules(namespace).Patch(context.Background(), name, pkgApi.JSONPatchType, data, metav1.PatchOptions{})
	if err != nil {
		return fmt.Errorf("Failed to update PrometheusRule: %s", err)
	}
	log.Printf("[INFO] Submitted updated PrometheusRule: %#v", out)

	return resourcePOPrometheusRuleRead(d, meta)
}

func resourcePOPrometheusRuleDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*KubeClientsets).MonitoringClient
	namespace, name, err := idParts(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleting PrometheusRule: %q", name)
	err = conn.PrometheusRules(namespace).Delete(context.Background(), name, metav1.DeleteOptions{})
	if err != nil {
		return err
	}

	log.Printf("[INFO] PrometheusRule %s deleted", name)

	d.SetId("")

	return nil
}

func expandPrometheusRuleSpec(groups []interface{}) (*po_types.PrometheusRuleSpec, error) {
	obj := &po_types.PrometheusRuleSpec{}
	if len(groups) == 0 || groups[0] == nil {
		return obj, nil
	}
	in := groups[0].(map[string]interface{})

	if v, ok := in["groups"].([]interface{}); ok && len(v) > 0 {
		g, err := expandRuleGroup(v)
		if err != nil {
			return obj, err
		}
		obj.Groups = g
	}

	return obj, nil
}

func flattenPrometheusRuleSpec(spec po_v1.PrometheusRuleSpec, d *schema.ResourceData) ([]interface{}, error) {
	att := make(map[string]interface{})

	groups, err := flattenRuleGroup(spec.Groups)
	if err != nil {
		return nil, err
	}
	att["groups"] = groups

	return []interface{}{att}, nil
}
