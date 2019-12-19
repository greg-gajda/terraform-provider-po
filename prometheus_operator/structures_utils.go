package prometheus_operator

import (
	po_types "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"k8s.io/apimachinery/pkg/util/intstr"
	"strconv"
	"strings"
)

func expandRuleGroup(groups []interface{}) ([]po_types.RuleGroup, error) {
	if len(groups) == 0 {
		return []po_types.RuleGroup{}, nil
	}
	obj := make([]po_types.RuleGroup, len(groups))
	for i, e := range groups {
		in := e.(map[string]interface{})
		if name, ok := in["name"]; ok {
			obj[i].Name = name.(string)
		}
		if interval, ok := in["interval"]; ok {
			obj[i].Interval = interval.(string)
		}
		if v, ok := in["rules"].([]interface{}); ok && len(v) > 0 {
			rules, err := expandRules(v)
			if err != nil {
				return obj, err
			}
			obj[i].Rules = rules
		}
	}
	return obj, nil
}

func flattenRuleGroup(in []po_types.RuleGroup) ([]interface{}, error) {
	att := make([]interface{}, len(in))
	for i, v := range in {
		out := make(map[string]interface{})
		out["name"] = v.Name
		out["interval"] = v.Interval
		out["rules"] = flattenRules(v.Rules)
		att[i] = out
	}
	return att, nil
}


func expandRules(rules []interface{}) ([]po_types.Rule, error) {
	if len(rules) == 0 {
		return []po_types.Rule{}, nil
	}
	obj := make([]po_types.Rule, len(rules))
	for i, e := range rules {
		in := e.(map[string]interface{})
		if alert, ok := in["alert"]; ok {
			obj[i].Alert = alert.(string)
		}
		if expr, ok := in["expr"]; ok {
			if num, err := strconv.Atoi(expr.(string)); err == nil {
				obj[i].Expr = intstr.FromInt(num)
			} else {
				obj[i].Expr = intstr.FromString(expr.(string))
			}
		}
		if f, ok := in["for"]; ok {
			obj[i].For = f.(string)
		}
		if record, ok := in["record"]; ok {
			obj[i].Record = record.(string)
		}
		if labels, ok := in["labels"].(map[string]interface{}); ok && len(labels) > 0 {
			obj[i].Labels = expandStringMap(in["labels"].(map[string]interface{}))
		}
		if a, ok := in["annotations"].(map[string]interface{}); ok && len(a) > 0 {
			obj[i].Annotations = expandStringMap(in["annotations"].(map[string]interface{}))
		}
	}
	return obj, nil
}

func expandStringMapTrim(m map[string]interface{}) map[string]string {
	result := make(map[string]string)
	for k, v := range m {
		result[k] = strings.TrimSpace(v.(string))
	}
	return result
}

func flattenRules(in []po_types.Rule) []interface{} {
	att := make([]interface{}, len(in))
	for i, v := range in {
		out := make(map[string]interface{})
		out["alert"] = v.Alert
		out["record"] = v.Record
		out["for"] = v.For
		out["labels"] = v.Labels
		out["annotations"] = v.Annotations
		out["expr"] = v.Expr.StrVal
		att[i] = out
	}
	return att
}

func expandAlertingSpec(l []interface{}) (*po_types.AlertingSpec, error) {
	obj := &po_types.AlertingSpec{}
	if len(l) == 0 || l[0] == nil {
		return obj, nil
	}
	in := l[0].(map[string]interface{})
	if v, ok := in["alertmanagers"].([]interface{}); ok && len(v) > 0 {
		endpoints, err := expandAlertmanagersEndpoints(v)
		if err != nil {
			return obj, err
		}
		obj.Alertmanagers = endpoints
	}
	return obj, nil
}

func flattenAlertingSpec(in *po_types.AlertingSpec) ([]interface{}, error) {
	att := make(map[string]interface{})
	if in != nil && len(in.Alertmanagers) > 0 {
		endpoints, err := flattenAlertmanagersEndpoints(in.Alertmanagers)
		if err != nil {
			return nil, err
		}
		att["alertmanagers"] = endpoints
	}
	return []interface{}{att}, nil
}

func expandAlertmanagersEndpoints(endpoints []interface{}) ([]po_types.AlertmanagerEndpoints, error) {
	if len(endpoints) == 0 {
		return []po_types.AlertmanagerEndpoints{}, nil
	}
	obj := make([]po_types.AlertmanagerEndpoints, len(endpoints))
	for i, e := range endpoints {
		in := e.(map[string]interface{})
		if namespace, ok := in["namespace"]; ok {
			obj[i].Namespace = namespace.(string)
		}
		if name, ok := in["name"]; ok {
			obj[i].Name = name.(string)
		}
		if port, ok := in["port"]; ok {
			if num, err := strconv.Atoi(port.(string)); err == nil {
				obj[i].Port = intstr.FromInt(num)
			} else {
				obj[i].Port = intstr.FromString(port.(string))
			}
		}
		if path, ok := in["path_prefix"]; ok {
			obj[i].PathPrefix = path.(string)
		}
		if scheme, ok := in["scheme"]; ok {
			obj[i].Scheme = scheme.(string)
		}
		if tls, ok := in["tls_config"].([]interface{}); ok && len(tls) > 0 {
			tls, err := expandTLSConfig(tls)
			if err != nil {
				return obj, err
			}
			obj[i].TLSConfig = tls
		}
		if btf, ok := in["bearer_token_file"]; ok {
			obj[i].BearerTokenFile = btf.(string)
		}
	}
	return obj, nil
}

func flattenAlertmanagersEndpoints(in []po_types.AlertmanagerEndpoints) ([]interface{}, error) {
	att := make([]interface{}, len(in))
	for i, v := range in {
		e := make(map[string]interface{})
		e["name"] = v.Name
		e["namespace"] = v.Namespace
		e["port"] = v.Port
		e["path_prefix"] = v.PathPrefix
		e["scheme"] = v.Scheme
		if v.TLSConfig != nil {
			e["tls_config"] = flattenTLSConfig(v.TLSConfig)
		}
		e["bearer_token_file"] = v.BearerTokenFile
		att[i] = e
	}
	return att, nil
}

func expandNamespaceSelector(l []interface{}) (*po_types.NamespaceSelector, error) {
	obj := &po_types.NamespaceSelector{}
	if len(l) == 0 || l[0] == nil {
		return obj, nil
	}
	in := l[0].(map[string]interface{})
	if a, ok := in["any"]; ok {
		obj.Any = a.(bool)
	}
	if v, ok := in["match_names"].(*schema.Set); ok && v.Len() > 0 {
		obj.MatchNames = sliceOfString(v.List())
	}
	return obj, nil
}

func flattenNamespaceSelector(in *po_types.NamespaceSelector) []interface{} {
	att := make(map[string]interface{})
	if in != nil {
		att["any"] = in.Any
		if len(in.MatchNames) > 0 {
			att["match_names"] = newStringSet(schema.HashString, in.MatchNames)
		}
	}
	return []interface{}{att}
}

func expandRelabelConfig(conf []interface{}) ([]*po_types.RelabelConfig, error) {
	obj := make([]*po_types.RelabelConfig, len(conf))
	if len(conf) == 0 || conf[0] == nil {
		return obj, nil
	}
	for i, e := range conf {
		in := e.(map[string]interface{})
		obj[i] = &po_types.RelabelConfig{}
		if s, ok := in["separator"]; ok {
			obj[i].Separator = s.(string)
		}
		if tl, ok := in["target_label"]; ok {
			obj[i].TargetLabel = tl.(string)
		}
		if r, ok := in["regex"]; ok {
			obj[i].Regex = r.(string)
		}
		if r, ok := in["replacement"]; ok {
			obj[i].Replacement = r.(string)
		}
		if m, ok := in["modulus"]; ok {
			obj[i].Modulus = uint64(m.(int))
		}
		if a, ok := in["action"]; ok {
			obj[i].Action = a.(string)
		}
		if v, ok := in["source_labels"].(*schema.Set); ok && v.Len() > 0 {
			obj[i].SourceLabels = sliceOfString(v.List())
		}
	}
	return obj, nil
}

func flattenRelabelConfig(in []*po_types.RelabelConfig) []interface{} {
	att := make([]interface{}, len(in))
	for i, v := range in {
		c := make(map[string]interface{})
		c["separator"] = v.Separator
		c["target_label"] = v.TargetLabel
		c["regex"] = v.Regex
		c["modulus"] = v.Modulus
		c["replacement"] = v.Replacement
		c["action"] = v.Action
		if len(v.SourceLabels) > 0 {
			c["source_labels"] = newStringSet(schema.HashString, v.SourceLabels)
		}
		att[i] = c
	}
	return att
}

func expandBasicAuth(l []interface{}) (*po_types.BasicAuth, error) {
	obj := &po_types.BasicAuth{}
	if len(l) == 0 || l[0] == nil {
		return obj, nil
	}
	in := l[0].(map[string]interface{})

	if v, ok := in["username"].([]interface{}); ok && len(v) > 0 {
		u, err := expandSecretKeyRef(v)
		if err != nil {
			return obj, err
		}
		obj.Username = *u
	}
	if v, ok := in["password"].([]interface{}); ok && len(v) > 0 {
		p, err := expandSecretKeyRef(v)
		if err != nil {
			return obj, err
		}
		obj.Password = *p
	}
	return obj, nil
}

func flattenBasicAuth(in *po_types.BasicAuth) []interface{} {
	att := make(map[string]interface{})
	att["username"] = flattenSecretKeyRef(&in.Username)
	att["password"] = flattenSecretKeyRef(&in.Password)
	return []interface{}{att}
}

func expandSecretOrConfigMap(l []interface{}) (*po_types.SecretOrConfigMap, error) {
	obj := &po_types.SecretOrConfigMap{}
	if len(l) == 0 || l[0] == nil {
		return obj, nil
	}
	in := l[0].(map[string]interface{})

	if v, ok := in["secret"].([]interface{}); ok && len(v) > 0 {
		s, err := expandSecretKeyRef(v)
		if err != nil {
			return obj, err
		}
		obj.Secret = s
	}
	if v, ok := in["config_map"].([]interface{}); ok && len(v) > 0 {
		cm, err := expandConfigMapKeyRef(v)
		if err != nil {
			return obj, err
		}
		obj.ConfigMap = cm
	}
	return obj, nil
}

func flattenSecretOrConfigMap(in *po_types.SecretOrConfigMap) []interface{} {
	att := make(map[string]interface{})
	if in != nil {
		if in.Secret != nil {
			att["secret"] = flattenSecretKeyRef(in.Secret)
		}
		if in.ConfigMap != nil {
			att["config_map"] = flattenConfigMapKeyRef(in.ConfigMap)
		}
	}
	return []interface{}{att}
}

func expandTLSConfig(l []interface{}) (*po_types.TLSConfig, error) {
	obj := &po_types.TLSConfig{}
	if len(l) == 0 || l[0] == nil {
		return obj, nil
	}
	in := l[0].(map[string]interface{})
	obj.CAFile = in["ca_file"].(string)
	if v, ok := in["ca"].([]interface{}); ok && len(v) > 0 {
		ca, err := expandSecretOrConfigMap(v)
		if err != nil {
			return obj, err
		}
		obj.CA = *ca
	}
	obj.CertFile = in["cert_file"].(string)
	if v, ok := in["cert"].([]interface{}); ok && len(v) > 0 {
		cert, err := expandSecretOrConfigMap(v)
		if err != nil {
			return obj, err
		}
		obj.Cert = *cert
	}
	obj.KeyFile = in["key_file"].(string)
	if v, ok := in["key_secret"].([]interface{}); ok && len(v) > 0 {
		ks, err := expandSecretKeyRef(v)
		if err != nil {
			return obj, err
		}
		obj.KeySecret = ks
	}
	obj.ServerName = in["server_name"].(string)
	obj.InsecureSkipVerify = in["insecure_skip_verify"].(bool)

	return obj, nil
}

func flattenTLSConfig(in *po_types.TLSConfig) []interface{} {
	att := make(map[string]interface{})
	att["ca_file"] = in.CAFile
	if in.CA.Secret != nil || in.CA.ConfigMap != nil {
		att["ca"] = flattenSecretOrConfigMap(&in.CA)
	}
	att["cert_file"] = in.CertFile
	if in.Cert.Secret != nil || in.Cert.ConfigMap != nil {
		att["cert"] = flattenSecretOrConfigMap(&in.Cert)
	}
	att["key_file"] = in.KeyFile
	if in.KeySecret != nil {
		att["key_secret"] = flattenSecretKeyRef(in.KeySecret)
	}
	att["server_name"] = in.ServerName
	att["insecure_skip_verify"] = in.InsecureSkipVerify
	return []interface{}{att}
}

func expandEndpoints(endpoints []interface{}) ([]po_types.Endpoint, error) {
	if len(endpoints) == 0 {
		return []po_types.Endpoint{}, nil
	}
	obj := make([]po_types.Endpoint, len(endpoints))
	for i, e := range endpoints {
		in := e.(map[string]interface{})
		if port, ok := in["port"]; ok {
			obj[i].Port = port.(string)
		}
		if path, ok := in["path"]; ok {
			obj[i].Path = path.(string)
		}
		if scheme, ok := in["scheme"]; ok {
			obj[i].Scheme = scheme.(string)
		}
		if interval, ok := in["interval"]; ok {
			obj[i].Interval = interval.(string)
		}
		if st, ok := in["scrape_timeout"]; ok {
			obj[i].ScrapeTimeout = st.(string)
		}
		if tls, ok := in["tls_config"].([]interface{}); ok && len(tls) > 0 {
			tls, err := expandTLSConfig(tls)
			if err != nil {
				return obj, err
			}
			obj[i].TLSConfig = tls
		}
		if btf, ok := in["bearer_token_file"]; ok {
			obj[i].BearerTokenFile = btf.(string)
		}
		if bts, ok := in["bearer_token_secret"].([]interface{}); ok && len(bts) > 0 {
			s, err := expandSecretKeyRef(bts)
			if err != nil {
				return obj, err
			}
			obj[i].BearerTokenSecret = *s
		}
		if hl, ok := in["honor_labels"]; ok {
			obj[i].HonorLabels = hl.(bool)
		}
		if ht, ok := in["honor_timestamps"]; ok {
			obj[i].HonorTimestamps = ptrToBool(ht.(bool))
		}
		if ba, ok := in["basic_auth"].([]interface{}); ok && len(ba) > 0 {
			ba, err := expandBasicAuth(ba)
			if err != nil {
				return obj, err
			}
			obj[i].BasicAuth = ba
		}
		if mrl, ok := in["metric_relabelings"].([]interface{}); ok && len(mrl) > 0 {
			c, err := expandRelabelConfig(mrl)
			if err != nil {
				return obj, err
			}
			obj[i].MetricRelabelConfigs = c
		}
		if re, ok := in["relabelings"].([]interface{}); ok && len(re) > 0 {
			c, err := expandRelabelConfig(re)
			if err != nil {
				return obj, err
			}
			obj[i].RelabelConfigs = c
		}
		if pu, ok := in["proxy_url"].(string); ok && pu != "" {
			obj[i].ProxyURL = ptrToString(pu)
		}
	}
	return obj, nil
}

func flattenEndpoints(in []po_types.Endpoint) ([]interface{}, error) {
	att := make([]interface{}, len(in))
	for i, v := range in {
		e := make(map[string]interface{})
		e["port"] = v.Port
		e["path"] = v.Path
		e["scheme"] = v.Scheme
		e["interval"] = v.Interval
		e["scrape_timeout"] = v.ScrapeTimeout
		if v.TLSConfig != nil {
			e["tls_config"] = flattenTLSConfig(v.TLSConfig)
		}
		e["bearer_token_file"] = v.BearerTokenFile
		if v.BearerTokenSecret.Key != "" || v.BearerTokenSecret.Name != "" {
			e["bearer_token_secret"] = flattenSecretKeyRef(&v.BearerTokenSecret)
		}
		e["honor_labels"] = v.HonorLabels
		if v.HonorTimestamps != nil {
			e["honor_timestamps"] = *v.HonorTimestamps
		}
		if v.BasicAuth != nil {
			e["basic_auth"] = flattenBasicAuth(v.BasicAuth)
		}
		e["metric_relabelings"] = flattenRelabelConfig(v.MetricRelabelConfigs)
		e["relabelings"] = flattenRelabelConfig(v.RelabelConfigs)
		if v.ProxyURL != nil {
			e["proxy_url"] = *v.ProxyURL
		}
		att[i] = e
	}
	return att, nil
}