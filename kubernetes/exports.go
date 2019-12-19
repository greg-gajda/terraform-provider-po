package kubernetes

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	api "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/api/core/v1"
)

func NamespacedMetadataSchema(objectName string, generatableName bool) *schema.Schema {
	return namespacedMetadataSchema(objectName, generatableName)
}

func SecurityContextSchema() *schema.Resource {
	return securityContextSchema()
}

func ResourcesField() map[string]*schema.Schema {
	return resourcesField()
}

func FlattenPodSecurityContext(in *v1.PodSecurityContext) []interface{} {
	return flattenPodSecurityContext(in)
}

func FlattenContainerSecurityContext(in *v1.SecurityContext) []interface{} {
	return flattenContainerSecurityContext(in)
}

func ExpandContainerSecurityContext(l []interface{}) *v1.SecurityContext {
	return expandContainerSecurityContext(l)
}

func ExpandPodSecurityContext(l []interface{}) *v1.PodSecurityContext {
	return expandPodSecurityContext(l)
}

func ExpandVolumes(volumes []interface{}) ([]v1.Volume, error) {
	return expandVolumes(volumes)
}

func ExpandTolerations(tolerations []interface{}) ([]*v1.Toleration, error) {
	return expandTolerations(tolerations)
}

func ExpandContainers(ctrs []interface{}) ([]v1.Container, error) {
	return expandContainers(ctrs)
}

func FlattenContainers(in []v1.Container) ([]interface{}, error) {
	return flattenContainers(in)
}

func ExpandContainerResourceRequirements(l []interface{}) (*v1.ResourceRequirements, error) {
	return expandContainerResourceRequirements(l)
}

func ExpandContainerVolumeMounts(in []interface{}) ([]v1.VolumeMount, error) {
	return expandContainerVolumeMounts(in)
}

func DiffStringMap(pathPrefix string, oldV, newV map[string]interface{}) PatchOperations {
	return diffStringMap(pathPrefix, oldV, newV)
}

func ContainerFields(isUpdatable, isInitContainer bool) map[string]*schema.Schema {
	return containerFields(isUpdatable, isInitContainer)
}

func VolumeSchema() *schema.Resource {
	return volumeSchema()
}

func VolumeMountFields() map[string]*schema.Schema {
	return volumeMountFields()
}

func IdParts(id string) (string, string, error) {
	return idParts(id)
}

func BuildId(meta metav1.ObjectMeta) string {
	return buildId(meta)
}

func ExpandMetadata(in []interface{}) metav1.ObjectMeta {
	return expandMetadata(in)
}

func PatchMetadata(keyPrefix, pathPrefix string, d *schema.ResourceData) PatchOperations {
	return patchMetadata(keyPrefix, pathPrefix, d)
}

func ExpandStringMap(m map[string]interface{}) map[string]string {
	return expandStringMap(m)
}

func ExpandBase64MapToByteMap(m map[string]interface{}) map[string][]byte {
	return expandBase64MapToByteMap(m)
}

func ExpandStringMapToByteMap(m map[string]interface{}) map[string][]byte {
	return expandStringMapToByteMap(m)
}

func ExpandStringSlice(s []interface{}) []string {
	return expandStringSlice(s)
}

func FlattenMetadata(meta metav1.ObjectMeta, d *schema.ResourceData, metaPrefix ...string) []interface{} {
	return flattenMetadata(meta, d, metaPrefix...)
}

func RemoveInternalKeys(m map[string]string, d map[string]interface{}) map[string]string {
	return removeInternalKeys(m, d)
}

func IsKeyInMap(key string, d map[string]interface{}) bool {
	return isKeyInMap(key, d)
}

func IsInternalKey(annotationKey string) bool {
	return isInternalKey(annotationKey)
}

func FlattenByteMapToBase64Map(m map[string][]byte) map[string]string {
	return flattenByteMapToBase64Map(m)
}

func FlattenByteMapToStringMap(m map[string][]byte) map[string]string {
	return flattenByteMapToStringMap(m)
}

func PtrToString(s string) *string {
	return ptrToString(s)
}

func PtrToInt(i int) *int {
	return ptrToInt(i)
}

func PtrToBool(b bool) *bool {
	return ptrToBool(b)
}

func PtrToInt32(i int32) *int32 {
	return ptrToInt32(i)
}

func PtrToInt64(i int64) *int64 {
	return ptrToInt64(i)
}

func SliceOfString(slice []interface{}) []string {
	return sliceOfString(slice)
}

func Base64EncodeStringMap(m map[string]interface{}) map[string]interface{} {
	return base64EncodeStringMap(m)
}

func FlattenResourceList(l api.ResourceList) map[string]string {
	return flattenResourceList(l)
}

func ExpandMapToResourceList(m map[string]interface{}) (*api.ResourceList, error) {
	return expandMapToResourceList(m)
}

func FlattenPersistentVolumeAccessModes(in []api.PersistentVolumeAccessMode) *schema.Set {
	return flattenPersistentVolumeAccessModes(in)
}

func ExpandPersistentVolumeAccessModes(s []interface{}) []api.PersistentVolumeAccessMode {
	return expandPersistentVolumeAccessModes(s)
}

func FlattenResourceQuotaSpec(in api.ResourceQuotaSpec) []interface{} {
	return flattenResourceQuotaSpec(in)
}

func ExpandResourceQuotaSpec(s []interface{}) (*api.ResourceQuotaSpec, error) {
	return expandResourceQuotaSpec(s)
}

func FlattenResourceQuotaScopes(in []api.ResourceQuotaScope) *schema.Set {
	return flattenResourceQuotaScopes(in)
}

func ExpandResourceQuotaScopes(s []interface{}) []api.ResourceQuotaScope {
	return expandResourceQuotaScopes(s)
}

func NewStringSet(f schema.SchemaSetFunc, in []string) *schema.Set {
	return newStringSet(f, in)
}
func NewInt64Set(f schema.SchemaSetFunc, in []int64) *schema.Set {
	return newInt64Set(f, in)
}

func ResourceListEquals(x, y api.ResourceList) bool {
	return resourceListEquals(x, y)
}

func ExpandLimitRangeSpec(s []interface{}, isNew bool) (*api.LimitRangeSpec, error) {
	return expandLimitRangeSpec(s, isNew)
}

func FlattenLimitRangeSpec(in api.LimitRangeSpec) []interface{} {
	return flattenLimitRangeSpec(in)
}

func SchemaSetToStringArray(set *schema.Set) []string {
	return schemaSetToStringArray(set)
}

func SchemaSetToInt64Array(set *schema.Set) []int64 {
	return schemaSetToInt64Array(set)
}
func FlattenLabelSelectorRequirementList(l []metav1.LabelSelectorRequirement) []interface{} {
	return flattenLabelSelectorRequirementList(l)
}

func FlattenLocalObjectReferenceArray(in []api.LocalObjectReference) []interface{} {
	return flattenLocalObjectReferenceArray(in)
}

func ExpandLocalObjectReferenceArray(in []interface{}) []api.LocalObjectReference {
	return expandLocalObjectReferenceArray(in)
}

func FlattenServiceAccountSecrets(in []api.ObjectReference, defaultSecretName string) []interface{} {
	return flattenServiceAccountSecrets(in, defaultSecretName)
}

func ExpandServiceAccountSecrets(in []interface{}, defaultSecretName string) []api.ObjectReference {
	return expandServiceAccountSecrets(in, defaultSecretName)
}

func FlattenNodeSelectorRequirementList(in []api.NodeSelectorRequirement) []map[string]interface{} {
	return flattenNodeSelectorRequirementList(in)
}

func ExpandNodeSelectorRequirementList(in []interface{}) []api.NodeSelectorRequirement {
	return expandNodeSelectorRequirementList(in)
}

func FlattenNodeSelectorTerm(in api.NodeSelectorTerm) []interface{} {
	return flattenNodeSelectorTerm(in)
}

func ExpandNodeSelectorTerm(l []interface{}) *api.NodeSelectorTerm {
	return expandNodeSelectorTerm(l)
}

func FlattenNodeSelectorTerms(in []api.NodeSelectorTerm) []interface{} {
	return flattenNodeSelectorTerms(in)
}

func ExpandNodeSelectorTerms(l []interface{}) []api.NodeSelectorTerm {
	return expandNodeSelectorTerms(l)
}


func ValidateTypeStringNullableInt(v interface{}, k string) (ws []string, es []error) {
	return validateTypeStringNullableInt(v, k)
}

func SeLinuxOptionsField() map[string]*schema.Schema {
	return seLinuxOptionsField()
}

func ExpandSecretKeyRef(r []interface{}) (*v1.SecretKeySelector, error) {
	return expandSecretKeyRef(r)
}

func FlattenSecretKeyRef(in *v1.SecretKeySelector) []interface{} {
	return flattenSecretKeyRef(in)
}

func ExpandConfigMapKeyRef(r []interface{}) (*v1.ConfigMapKeySelector, error) {
	return expandConfigMapKeyRef(r)
}

func FlattenConfigMapKeyRef(in *v1.ConfigMapKeySelector) []interface{} {
	return flattenConfigMapKeyRef(in)
}

func FlattenLabelSelector(in *metav1.LabelSelector) []interface{} {
	return flattenLabelSelector(in)
}

func ExpandLabelSelector(l []interface{}) *metav1.LabelSelector {
	return expandLabelSelector(l)
}

func LabelSelectorFields(updatable bool) map[string]*schema.Schema {
	return labelSelectorFields(updatable)
}

func ValidateLabels(value interface{}, key string) (ws []string, es []error) {
	return validateLabels(value, key)
}