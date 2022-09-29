package prometheus_operator

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	k8s "github.com/hashicorp/terraform-provider-kubernetes"
	api "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func namespacedMetadataSchema(objectName string, generatableName bool) *schema.Schema {
	return k8s.NamespacedMetadataSchema(objectName, generatableName)
}

func securityContextSchema() *schema.Resource {
	return k8s.SecurityContextSchema()
}

func resourcesField() map[string]*schema.Schema {
	return k8s.ResourcesField()
}


func flattenContainerSecurityContext(in *v1.SecurityContext) []interface{} {
	return k8s.FlattenContainerSecurityContext(in)
}

func expandContainerSecurityContext(l []interface{}) *v1.SecurityContext {
	return k8s.ExpandContainerSecurityContext(l)
}

func flattenPodSecurityContext(in *v1.PodSecurityContext) []interface{} {
	return k8s.FlattenPodSecurityContext(in)
}

func expandPodSecurityContext(l []interface{}) *v1.PodSecurityContext {
	return k8s.ExpandPodSecurityContext(l)
}

func expandVolumes(volumes []interface{}) ([]v1.Volume, error) {
	return k8s.ExpandVolumes(volumes)
}

func expandTolerations(tolerations []interface{}) ([]*v1.Toleration, error) {
	return k8s.ExpandTolerations(tolerations)
}

func expandContainers(ctrs []interface{}) ([]v1.Container, error) {
	return k8s.ExpandContainers(ctrs)
}

func flattenContainers(in []v1.Container) ([]interface{}, error) {
	return k8s.FlattenContainers(in)
}

func expandContainerResourceRequirements(l []interface{}) (*v1.ResourceRequirements, error) {
	return k8s.ExpandContainerResourceRequirements(l)
}

func expandContainerVolumeMounts(in []interface{}) ([]v1.VolumeMount, error) {
	return k8s.ExpandContainerVolumeMounts(in)
}



func diffStringMap(pathPrefix string, oldV, newV map[string]interface{}) k8s.PatchOperations {
	return k8s.DiffStringMap(pathPrefix, oldV, newV)
}

func containerFields(isUpdatable, isInitContainer bool) map[string]*schema.Schema {
	return k8s.ContainerFields(isUpdatable, isInitContainer)
}

func volumeSchema(isUpdatable bool) *schema.Resource {
	return k8s.VolumeSchema(isUpdatable)
}

func volumeMountFields() map[string]*schema.Schema {
	return k8s.VolumeMountFields()
}

func idParts(id string) (string, string, error) {
	return k8s.IdParts(id)
}

func buildId(meta metav1.ObjectMeta) string {
	return k8s.BuildId(meta)
}

func expandMetadata(in []interface{}) metav1.ObjectMeta {
	return k8s.ExpandMetadata(in)
}

func patchMetadata(keyPrefix, pathPrefix string, d *schema.ResourceData) k8s.PatchOperations {
	return k8s.PatchMetadata(keyPrefix, pathPrefix, d)
}

func expandStringMap(m map[string]interface{}) map[string]string {
	return k8s.ExpandStringMap(m)
}

func expandBase64MapToByteMap(m map[string]interface{}) map[string][]byte {
	return k8s.ExpandBase64MapToByteMap(m)
}

func expandStringMapToByteMap(m map[string]interface{}) map[string][]byte {
	return k8s.ExpandStringMapToByteMap(m)
}

func expandStringSlice(s []interface{}) []string {
	return k8s.ExpandStringSlice(s)
}

func flattenMetadata(meta metav1.ObjectMeta, d *schema.ResourceData, metaPrefix ...string) []interface{} {
	return k8s.FlattenMetadata(meta, d, metaPrefix...)
}

func removeInternalKeys(m map[string]string, d map[string]interface{}) map[string]string {
	return k8s.RemoveInternalKeys(m, d)
}

func isKeyInMap(key string, d map[string]interface{}) bool {
	return k8s.IsKeyInMap(key, d)
}

func isInternalKey(annotationKey string) bool {
	return k8s.IsInternalKey(annotationKey)
}

func flattenByteMapToBase64Map(m map[string][]byte) map[string]string {
	return k8s.FlattenByteMapToBase64Map(m)
}

func flattenByteMapToStringMap(m map[string][]byte) map[string]string {
	return k8s.FlattenByteMapToStringMap(m)
}

func ptrToString(s string) *string {
	return k8s.PtrToString(s)
}

func ptrToInt(i int) *int {
	return k8s.PtrToInt(i)
}

func ptrToBool(b bool) *bool {
	return k8s.PtrToBool(b)
}

func ptrToInt32(i int32) *int32 {
	return k8s.PtrToInt32(i)
}

func ptrToInt64(i int64) *int64 {
	return k8s.PtrToInt64(i)
}

func sliceOfString(slice []interface{}) []string {
	return k8s.SliceOfString(slice)
}

func base64EncodeStringMap(m map[string]interface{}) map[string]interface{} {
	return k8s.Base64EncodeStringMap(m)
}

func flattenResourceList(l api.ResourceList) map[string]string {
	return k8s.FlattenResourceList(l)
}

func expandMapToResourceList(m map[string]interface{}) (*api.ResourceList, error) {
	return k8s.ExpandMapToResourceList(m)
}

func flattenPersistentVolumeAccessModes(in []api.PersistentVolumeAccessMode) *schema.Set {
	return k8s.FlattenPersistentVolumeAccessModes(in)
}

func expandPersistentVolumeAccessModes(s []interface{}) []api.PersistentVolumeAccessMode {
	return k8s.ExpandPersistentVolumeAccessModes(s)
}

func flattenResourceQuotaSpec(in api.ResourceQuotaSpec) []interface{} {
	return k8s.FlattenResourceQuotaSpec(in)
}

func expandResourceQuotaSpec(s []interface{}) (*api.ResourceQuotaSpec, error) {
	return k8s.ExpandResourceQuotaSpec(s)
}

func flattenResourceQuotaScopes(in []api.ResourceQuotaScope) *schema.Set {
	return k8s.FlattenResourceQuotaScopes(in)
}

func expandResourceQuotaScopes(s []interface{}) []api.ResourceQuotaScope {
	return k8s.ExpandResourceQuotaScopes(s)
}

func newStringSet(f schema.SchemaSetFunc, in []string) *schema.Set {
	return k8s.NewStringSet(f, in)
}
func newInt64Set(f schema.SchemaSetFunc, in []int64) *schema.Set {
	return k8s.NewInt64Set(f, in)
}

func resourceListEquals(x, y api.ResourceList) bool {
	return k8s.ResourceListEquals(x, y)
}

func expandLimitRangeSpec(s []interface{}, isNew bool) (*api.LimitRangeSpec, error) {
	return k8s.ExpandLimitRangeSpec(s, isNew)
}

func flattenLimitRangeSpec(in api.LimitRangeSpec) []interface{} {
	return k8s.FlattenLimitRangeSpec(in)
}

func schemaSetToStringArray(set *schema.Set) []string {
	return k8s.SchemaSetToStringArray(set)
}

func schemaSetToInt64Array(set *schema.Set) []int64 {
	return k8s.SchemaSetToInt64Array(set)
}
func flattenLabelSelectorRequirementList(l []metav1.LabelSelectorRequirement) []interface{} {
	return k8s.FlattenLabelSelectorRequirementList(l)
}

func flattenLocalObjectReferenceArray(in []api.LocalObjectReference) []interface{} {
	return k8s.FlattenLocalObjectReferenceArray(in)
}

func expandLocalObjectReferenceArray(in []interface{}) []api.LocalObjectReference {
	return k8s.ExpandLocalObjectReferenceArray(in)
}

func flattenServiceAccountSecrets(in []api.ObjectReference, defaultSecretName string) []interface{} {
	return k8s.FlattenServiceAccountSecrets(in, defaultSecretName)
}

func expandServiceAccountSecrets(in []interface{}, defaultSecretName string) []api.ObjectReference {
	return k8s.ExpandServiceAccountSecrets(in, defaultSecretName)
}

func flattenNodeSelectorRequirementList(in []api.NodeSelectorRequirement) []map[string]interface{} {
	return k8s.FlattenNodeSelectorRequirementList(in)
}

func expandNodeSelectorRequirementList(in []interface{}) []api.NodeSelectorRequirement {
	return k8s.ExpandNodeSelectorRequirementList(in)
}

func flattenNodeSelectorTerm(in api.NodeSelectorTerm) []interface{} {
	return k8s.FlattenNodeSelectorTerm(in)
}

func expandNodeSelectorTerm(l []interface{}) *api.NodeSelectorTerm {
	return k8s.ExpandNodeSelectorTerm(l)
}

func flattenNodeSelectorTerms(in []api.NodeSelectorTerm) []interface{} {
	return k8s.FlattenNodeSelectorTerms(in)
}

func expandNodeSelectorTerms(l []interface{}) []api.NodeSelectorTerm {
	return k8s.ExpandNodeSelectorTerms(l)
}


func validateTypeStringNullableInt(v interface{}, k string) (ws []string, es []error) {
	return k8s.ValidateTypeStringNullableInt(v, k)
}


func replace(spec interface{}) *k8s.ReplaceOperation {
	return &k8s.ReplaceOperation{
		Path:  "/spec",
		Value: spec,
	}
}

func seLinuxOptionsField() map[string]*schema.Schema {
	return k8s.SeLinuxOptionsField()
}

func flattenSecretKeyRef(in *v1.SecretKeySelector) []interface{} {
	return k8s.FlattenSecretKeyRef(in)
}

func flattenConfigMapKeyRef(in *v1.ConfigMapKeySelector) []interface{} {
	return k8s.FlattenConfigMapKeyRef(in)
}

func expandSecretKeyRef(r []interface{}) (*v1.SecretKeySelector, error) {
	return k8s.ExpandSecretKeyRef(r)
}

func expandConfigMapKeyRef(r []interface{}) (*v1.ConfigMapKeySelector, error) {
	return k8s.ExpandConfigMapKeyRef(r)
}

func flattenLabelSelector(in *metav1.LabelSelector) []interface{} {
	return k8s.FlattenLabelSelector(in)
}

func expandLabelSelector(l []interface{}) *metav1.LabelSelector {
	return k8s.ExpandLabelSelector(l)
}

func labelSelectorFields(updatable bool) map[string]*schema.Schema {
	return k8s.LabelSelectorFields(updatable)
}

func validateLabels(value interface{}, key string) (ws []string, es []error) {
	return k8s.ValidateLabels(value, key)
}