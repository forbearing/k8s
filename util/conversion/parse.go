package conversion

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes/scheme"
)

// YamlOrJsonToObject will decode yaml/json documents to runtime.Object and GroupVersionKind.
func YamlOrJsonToObject(data []byte) (runtime.Object, *schema.GroupVersionKind, error) {
	// UniversalDeserializer can convert any stored data recognized by this factory
	// into a Go object that satisfies runtime.Object.
	// It does not perform conversion. It does not perform defaulting.
	return scheme.Codecs.UniversalDeserializer().Decode(data, nil, nil)

	//decoder := serializer.NewCodecFactory(scheme.Scheme).UniversalDecoder()
	//object := &v1.Namespace{}
	//err := runtime.DecodeInto(decoder, data, object)
	//if err != nil {
	//    return nil, err
	//}
	//return object, nil
}

// YamlToJson convert yaml documents to json documents, and those are k8s resources.
func YamlToJson(yamlData []byte) ([]byte, error) {
	return yaml.ToJSON(yamlData)

	//if yamlData == nil {
	//    return nil, nil
	//}

	//rawExtension := runtime.RawExtension{}
	//if err := yaml.NewYAMLToJSONDecoder(bytes.NewReader(yamlData)).Decode(&rawExtension); err != nil {
	//    return nil, err
	//}
	//return rawExtension.Raw, nil
}

//// JsonToYaml converts a signal json document to yaml document.
//func JsonToYaml(jsonData []byte) ([]byte, error) {
//}

//// JsonToYaml
//func JsonToYaml(jsonData []byte) ([]byte, error) {
//}

//func YamlOrJsonToUnstructured(data []byte) (*unstructured.Unstructured, *schema.GroupVersionKind, error) {

//}
