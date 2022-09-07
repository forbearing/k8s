package restmapper

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	serializeryaml "k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	utilyaml "k8s.io/apimachinery/pkg/util/yaml"
)

/*
references:
    https://ymmt2005.hatenablog.com/entry/2020/04/14/An_example_of_using_dynamic_client_of_k8s.io/client-go#Mapping-between-GVK-and-GVR
*/

// FindGVR find the GroupVersionResource from signal yaml document or json document.
//
// Supported type are: map[string]interface{}, *unstructured.Unstructured,
// unstructured.Unstructured, string and []byte.
func FindGVR(restMapper meta.RESTMapper, obj interface{}) (schema.GroupVersionResource, error) {
	var gvr = schema.GroupVersionResource{}

	gvk, err := FindGVK(restMapper, obj)
	if err != nil {
		return gvr, err
	}

	restMapping, err := restMapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return gvr, err
	}
	return restMapping.Resource, nil
}

// FindGVK find the GroupVersionKind from signal yaml document or json document.
//
// Supported type are: map[string]interface{}, *unstructured.Unstructured,
// unstructured.Unstructured, string and []byte.
func FindGVK(restMapper meta.RESTMapper, obj interface{}) (schema.GroupVersionKind, error) {
	gvk := schema.GroupVersionKind{}

	switch val := obj.(type) {
	case map[string]interface{}:
		data, err := json.Marshal(val)
		if err != nil {
			return gvk, err
		}
		return findGVK(restMapper, data)
	case *unstructured.Unstructured:
		data, err := json.Marshal(val)
		if err != nil {
			return gvk, err
		}
		return findGVK(restMapper, data)
	case unstructured.Unstructured:
		data, err := json.Marshal(&val)
		if err != nil {
			return gvk, err
		}
		return findGVK(restMapper, data)
	case string:
		data, err := ioutil.ReadFile(val)
		if err != nil {
			return gvk, err
		}
		return findGVK(restMapper, data)
	case []byte:
		return findGVK(restMapper, val)
	default:
		return gvk, errors.New("type must be string, []byte, map[string]interface{}, *unstructured.Unstructured, unstructured.Unstructured or runtime.Object")
	}
}

// IsNamespaced check if the gvk is namespace scope.
func IsNamespaced(restMapper meta.RESTMapper, gvk schema.GroupVersionKind) (bool, error) {
	restMapping, err := restMapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return false, err
	}
	if restMapping.Scope.Name() == meta.RESTScopeNameNamespace {
		return true, nil
	}
	return false, nil
}

// GVKToGVR convert GVK to GVR.
func GVKToGVR(restMapper meta.RESTMapper, gvk schema.GroupVersionKind) (schema.GroupVersionResource, error) {
	restMapping, err := restMapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return schema.GroupVersionResource{}, err
	}
	return restMapping.Resource, nil
}

// findGVK find the GroupVersionKind from signal yaml document or json document.
func findGVK(restMapper meta.RESTMapper, data []byte) (schema.GroupVersionKind, error) {
	var (
		err error
		gvk = &schema.GroupVersionKind{}
		// RawExtension is used to hold extensions in external versions.
		rawExtension = runtime.RawExtension{}
	)

	// NewYAMLOrJSONDecoder returns a decoder that will process YAML documents
	// or JSON documents from the given reader as a stream. bufferSize determines
	// how far into the stream the decoder will look to figure out whether this
	// is a JSON stream (has whitespace followed by an open brace).
	decoder := utilyaml.NewYAMLOrJSONDecoder(bytes.NewReader(data), 4096)
	// Decode unmarshals the next object from the underlying stream into the
	// provide object, or returns an error.
	if err := decoder.Decode(&rawExtension); err != nil && err != io.EOF {
		return *gvk, err
	}

	// if the yaml object is empty just continue to the next one
	if len(rawExtension.Raw) == 0 || bytes.Equal(rawExtension.Raw, []byte("null")) {
		return *gvk, errors.New("the underlying serialization of this object is empty")
	}

	// NewDecodingSerializer adds YAML decoding support to a serializer that supports JSON.
	// json serializer runtime.Serializer --> runtime.Object, *schema.GroupVersionKind
	//fmt.Println(string(rawExtension.Raw))
	_, gvk, err = serializeryaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme).Decode(rawExtension.Raw, nil, nil)
	if err != nil {
		return *gvk, err
	}
	return *gvk, nil
}
