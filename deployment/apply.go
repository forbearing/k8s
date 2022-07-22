package deployment

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"

	log "github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	serializeryaml "k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/apimachinery/pkg/util/yaml"
	utilyaml "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/restmapper"
)

// Apply applies deployment from type string, []byte, *appsv1.Deployment,
// appsv1.Deployment, runtime.Object or map[string]interface{}.
func (h *Handler) Apply(obj interface{}) (*appsv1.Deployment, error) {
	switch val := obj.(type) {
	case string:
		return h.ApplyFromFile(val)
	case []byte:
		return h.ApplyFromBytes(val)
	case *appsv1.Deployment:
		return h.ApplyFromObject(val)
	case appsv1.Deployment:
		return h.ApplyFromObject(&val)
	case runtime.Object:
		return h.ApplyFromObject(val)
	case map[string]interface{}:
		return h.ApplyFromUnstructured(val)
	default:
		return nil, ERR_TYPE_APPLY
	}
}

// ApplyFromFile applies deployment from yaml file.
func (h *Handler) ApplyFromFile(filename string) (deploy *appsv1.Deployment, err error) {
	deploy, err = h.CreateFromFile(filename)
	if k8serrors.IsAlreadyExists(err) { // if deployment already exist, update it.
		deploy, err = h.UpdateFromFile(filename)
	}
	return
}

// ApplyFromBytes pply deployment from bytes.
func (h *Handler) ApplyFromBytes(data []byte) (deploy *appsv1.Deployment, err error) {
	deploy, err = h.CreateFromBytes(data)
	if k8serrors.IsAlreadyExists(err) {
		deploy, err = h.UpdateFromBytes(data)
	}
	return
}

// ApplyFromObject applies deployment from runtime.Object.
func (h *Handler) ApplyFromObject(obj runtime.Object) (*appsv1.Deployment, error) {
	deploy, ok := obj.(*appsv1.Deployment)
	if !ok {
		return nil, fmt.Errorf("object is not *appsv1.Deployment")
	}
	return h.applyDeployment(deploy)
}

// ApplyFromUnstructured applies deployment from map[string]interface{}.
func (h *Handler) ApplyFromUnstructured(u map[string]interface{}) (*appsv1.Deployment, error) {
	deploy := &appsv1.Deployment{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, deploy)
	if err != nil {
		return nil, err
	}
	return h.applyDeployment(deploy)
}

// applyDeployment
func (h *Handler) applyDeployment(deploy *appsv1.Deployment) (*appsv1.Deployment, error) {
	//var namespace string
	//if len(deploy.Namespace) != 0 {
	//    namespace = deploy.Namespace
	//} else {
	//    namespace = h.namespace
	//}

	//// 这里不能再用 deploy 来接收 Create 的结果, 以为如果 Create 失败了, 返回的
	//// 是一个空的 deploy, 后续的 Update deployment 就会失败.
	//_, err := h.clientset.AppsV1().Deployments(namespace).Create(h.ctx, deploy, h.Options.CreateOptions)
	//if k8serrors.IsAlreadyExists(err) {
	//    deploy, err = h.clientset.AppsV1().Deployments(namespace).Update(h.ctx, deploy, h.Options.UpdateOptions)
	//}
	//return deploy, err
	_, err := h.createDeployment(deploy)
	if k8serrors.IsAlreadyExists(err) {
		//log.Println("create failed, update it.")
		return h.updateDeployment(deploy)
	}
	return nil, err
}

// Don't Use This Method, Just for Testzng, May Be Removed.
// reserved it here as my study notes (hahaha).
func (h *Handler) __Apply(filename string) (deploy *appsv1.Deployment, err error) {
	var (
		data            []byte
		deployJson      []byte
		namespace       string
		bufferSize      = 500
		unstructuredMap map[string]interface{}
		unstructuredObj = &unstructured.Unstructured{}
	)
	deploy = &appsv1.Deployment{}
	if data, err = ioutil.ReadFile(filename); err != nil {
		return
	}
	if deployJson, err = yaml.ToJSON(data); err != nil {
		return
	}
	if err = json.Unmarshal(deployJson, deploy); err != nil {
		return
	}
	if len(deploy.Namespace) != 0 {
		namespace = deploy.Namespace
	} else {
		namespace = h.namespace
	}
	_ = namespace
	// NewYAMLOrJSONDecoder returns a decoder that will process YAML documents
	// or JSON documents from the given reader as a stream. bufferSize determines
	// how far into the stream the decoder will look to figure out whether this
	// is a JSON stream (has whitespace followed by an open brace).
	// yaml documents io.Reader  --> yaml decoder util/yaml.YAMLOrJSONDecoder
	decoder := utilyaml.NewYAMLOrJSONDecoder(bytes.NewReader(data), bufferSize)
	for {
		// RawExtension is used to hold extensions in external versions.
		var rawObject runtime.RawExtension
		// Decode reads a YAML document as JSON from the stream or returns an error.
		// The decoding rules match json.Unmarshal, not yaml.Unmarshal.
		// 用来判断文件内容是不是 yaml 格式. 如果 decoder.Decode 返回了错误, 说明文件内容
		// 不是 yaml 格式的. 如果返回 nil, 说明文件内容是 yaml 格式
		// yaml decoder util/yaml.YAMLOrJSONDecoder --> json serializer runtime.Serializer
		if err := decoder.Decode(&rawObject); err != nil {
			break
		}
		if len(rawObject.Raw) == 0 {
			// if the yaml object is empty just continue to the next one
			continue
		}
		// NewDecodingSerializer adds YAML decoding support to a serializer that supports JSON.
		// json serializer runtime.Serializer --> runtime.Object, *schema.GroupVersionKind
		object, gvk, err := serializeryaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme).Decode(rawObject.Raw, nil, nil)
		if err != nil {
			log.Error("NewDecodingSerializer error")
			log.Error(err)
			return nil, err
		}
		// runtime.Object --> map[string]interface{}
		unstructuredMap, err = runtime.DefaultUnstructuredConverter.ToUnstructured(object)
		if err != nil {
			return nil, err
		}
		// map[string]interface{} --> unstructured.Unstructured
		unstructuredObj = &unstructured.Unstructured{Object: unstructuredMap}

		// GetAPIGroupResources uses the provided discovery client to gather
		// discovery information and populate a slice of APIGroupResources.
		// DiscoveryInterface / DiscoveryClient --> []*APIGroupResources
		apiGroupResources, err := restmapper.GetAPIGroupResources(h.clientset.Discovery())
		if err != nil {
			log.Error("GetAPIGroupResources error")
			log.Error(err)
			return nil, err
		}

		// NewDiscoveryRESTMapper returns a PriorityRESTMapper based on the discovered
		// groups and resources passed in.
		// []*APIGroupResources --> meta.RESTMapper
		restMapper := restmapper.NewDiscoveryRESTMapper(apiGroupResources)

		// meta.RESTMapper -> meta.RESTMapping
		// RESTMapping identifies a preferred resource mapping for the provided group kind.
		restMapping, err := restMapper.RESTMapping(gvk.GroupKind(), gvk.Version)
		if err != nil {
			log.Error("RESTMapping error")
			log.Error(err)
			return nil, err
		}

		var dri dynamic.ResourceInterface
		// Scope contains the information needed to deal with REST Resources that are in a resource hierarchy
		// meta.RESTMapping.Resource --> shcema.GropuVersionResource
		if restMapping.Scope.Name() == meta.RESTScopeNameNamespace { // meta.RESTScopeNameNamespace is a const, and value is default
			if unstructuredObj.GetNamespace() == "" {
				unstructuredObj.SetNamespace("default")
			}
			dri = h.dynamicClient.Resource(restMapping.Resource).Namespace(unstructuredObj.GetNamespace())
		} else {
			dri = h.dynamicClient.Resource(restMapping.Resource)
		}
		_, err = dri.Create(context.Background(), unstructuredObj, metav1.CreateOptions{})
		if k8serrors.IsAlreadyExists(err) {
			_, err = dri.Update(context.Background(), unstructuredObj, metav1.UpdateOptions{})
		}
		if err != nil {
			log.Error("DynamicResourceInterface Apply error")
			log.Error(err)
			return nil, err
		}
	}
	//if err != io.EOF {
	//    log.Error("not io.EOF")
	//    log.Error(err)
	//    return nil, err
	//}

	if err = runtime.DefaultUnstructuredConverter.FromUnstructured(unstructuredObj.UnstructuredContent(), deploy); err != nil {
		log.Error("FromUnstructured error")
		log.Error(err)
		return nil, err
	}
	return deploy, nil
}
