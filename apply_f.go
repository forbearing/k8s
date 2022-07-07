package k8s

import (
	"bytes"
	"context"
	"io/ioutil"
	"regexp"

	"github.com/forbearing/k8s/clusterrole"
	"github.com/forbearing/k8s/clusterrolebinding"
	"github.com/forbearing/k8s/configmap"
	"github.com/forbearing/k8s/cronjob"
	"github.com/forbearing/k8s/daemonset"
	"github.com/forbearing/k8s/deployment"
	"github.com/forbearing/k8s/ingress"
	"github.com/forbearing/k8s/ingressclass"
	"github.com/forbearing/k8s/job"
	"github.com/forbearing/k8s/namespace"
	"github.com/forbearing/k8s/networkpolicy"
	"github.com/forbearing/k8s/persistentvolume"
	"github.com/forbearing/k8s/persistentvolumeclaim"
	"github.com/forbearing/k8s/pod"
	"github.com/forbearing/k8s/replicaset"
	"github.com/forbearing/k8s/replicationcontroller"
	"github.com/forbearing/k8s/role"
	"github.com/forbearing/k8s/rolebinding"
	"github.com/forbearing/k8s/secret"
	"github.com/forbearing/k8s/service"
	"github.com/forbearing/k8s/serviceaccount"
	"github.com/forbearing/k8s/statefulset"
	"github.com/forbearing/k8s/storageclass"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	networking "k8s.io/api/networking/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	storagev1 "k8s.io/api/storage/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/klog"
)

// ApplyF like "kubectl apply -f filename.yaml".
// the file could include multiple k8s resource.
// not support create "Node" from yaml file
func ApplyF(ctx context.Context, kubeconfig, filename string) error {
	k8sResourceFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	// remove all comments in the yaml file.
	removeComments := regexp.MustCompile(`#.*`)
	k8sResourceFile = removeComments.ReplaceAll(k8sResourceFile, []byte(""))
	// split yaml file by "---"
	k8sResourceItems := bytes.Split(k8sResourceFile, []byte("---"))

	for _, k8sResource := range k8sResourceItems {
		// ignore empty line
		if len(bytes.TrimSpace(k8sResource)) == 0 {
			continue
		}
		object, err := Decode(k8sResource)
		if err != nil {
			klog.V(4).Info("Decode error: %v", err)
			continue
		}
		switch object.(type) {
		case *corev1.Namespace:
			handler, err := namespace.New(ctx, kubeconfig)
			if err != nil {
				return err
			}
			if ns, err := handler.ApplyFromBytes(k8sResource); err != nil {
				klog.V(4).Info("apply namespace %q failed: %v\n", ns.Name, err)
			}
		case *corev1.Service:
			handler, err := service.New(ctx, "", kubeconfig)
			if err != nil {
				return err
			}
			if svc, err := handler.ApplyFromBytes(k8sResource); err != nil {
				klog.V(4).Info("apply service %q failed: %v\n", svc.Name, err)
			}
		case *corev1.ConfigMap:
			handler, err := configmap.New(ctx, "", kubeconfig)
			if err != nil {
				return err
			}
			if cm, err := handler.ApplyFromBytes(k8sResource); err != nil {
				klog.V(4).Info("apply configmap %q failed: %v\n", cm.Name, err)
			}
		case *corev1.Secret:
			handler, err := secret.New(ctx, "", kubeconfig)
			if err != nil {
				return err
			}
			if sec, err := handler.ApplyFromBytes(k8sResource); err != nil {
				klog.V(4).Info("apply secret %q failed: %v\n", sec.Name, err)
			}
		case *corev1.ServiceAccount:
			handler, err := serviceaccount.New(ctx, "", kubeconfig)
			if err != nil {
				return err
			}
			if sa, err := handler.ApplyFromBytes(k8sResource); err != nil {
				klog.V(4).Info("apply serviceaccount %q failed: %v\n", sa.Name, err)
			}
		case *corev1.Pod:
			handler, err := pod.New(ctx, "", kubeconfig)
			if err != nil {
				return err
			}
			if p, err := handler.ApplyFromBytes(k8sResource); err != nil {
				klog.V(4).Info("apply pod %q failed: %v\n", p.Name, err)
			}
		case *corev1.PersistentVolume:
			handler, err := persistentvolume.New(ctx, kubeconfig)
			if err != nil {
				return err
			}
			if pv, err := handler.ApplyFromBytes(k8sResource); err != nil {
				klog.V(4).Info("apply persistentvolume %q failed: %v\n", pv.Name, err)
			}
		case *corev1.PersistentVolumeClaim:
			handler, err := persistentvolumeclaim.New(ctx, "", kubeconfig)
			if err != nil {
				return err
			}
			if pvc, err := handler.ApplyFromBytes(k8sResource); err != nil {
				klog.V(4).Info("apply persistentvolumeclaim %q failed: %v\n", pvc.Name, err)
			}
		case *appsv1.Deployment:
			handler, err := deployment.New(ctx, "", kubeconfig)
			if err != nil {
				return err
			}
			if deploy, err := handler.ApplyFromBytes(k8sResource); err != nil {
				klog.V(4).Info("apply deployment %q failed: %v\n", deploy.Name, err)
			}
		case *appsv1.StatefulSet:
			handler, err := statefulset.New(ctx, "", kubeconfig)
			if err != nil {
				return err
			}
			if sts, err := handler.ApplyFromBytes(k8sResource); err != nil {
				klog.V(4).Info("apply statefulset %q failed: %v\n", sts.Name, err)
			}
		case *appsv1.DaemonSet:
			handler, err := daemonset.New(ctx, "", kubeconfig)
			if err != nil {
				return err
			}
			if ds, err := handler.ApplyFromBytes(k8sResource); err != nil {
				klog.V(4).Info("apply daemonset %q failed: %v\n", ds.Name, err)
			}
		case *networking.Ingress:
			handler, err := ingress.New(ctx, "", kubeconfig)
			if err != nil {
				return err
			}
			if ing, err := handler.ApplyFromBytes(k8sResource); err != nil {
				klog.V(4).Info("apply ingress %q failed: %v\n", ing.Name, err)
			}
		case *networking.IngressClass:
			handler, err := ingressclass.New(ctx, kubeconfig)
			if err != nil {
				return err
			}
			if ingc, err := handler.ApplyFromBytes(k8sResource); err != nil {
				klog.V(4).Info("apply ingressclass %q failed: %v\n", ingc.Name, err)
			}
		case *networking.NetworkPolicy:
			handler, err := networkpolicy.New(ctx, "", kubeconfig)
			if err != nil {
				return err
			}
			if netpol, err := handler.ApplyFromBytes(k8sResource); err != nil {
				klog.V(4).Info("apply networkpolicy %q failed: %v\n", netpol.Name, err)
			}
		case *batchv1.Job:
			handler, err := job.New(ctx, "", kubeconfig)
			if err != nil {
				return err
			}
			if j, err := handler.ApplyFromBytes(k8sResource); err != nil {
				klog.V(4).Info("apply job %q failed: %v\n", j.Name, err)
			}
		case *batchv1.CronJob:
			handler, err := cronjob.New(ctx, "", kubeconfig)
			if err != nil {
				return err
			}
			if cj, err := handler.ApplyFromBytes(k8sResource); err != nil {
				klog.V(4).Info("apply cronjob %q failed: %v\n", cj.Name, err)
			}
		case *rbacv1.Role:
			handler, err := role.New(ctx, "", kubeconfig)
			if err != nil {
				return err
			}
			if r, err := handler.ApplyFromBytes(k8sResource); err != nil {
				klog.V(4).Info("apply role %q failed: %v\n", r.Name, err)
			}
		case *rbacv1.RoleBinding:
			handler, err := rolebinding.New(ctx, "", kubeconfig)
			if err != nil {
				return err
			}
			if rb, err := handler.ApplyFromBytes(k8sResource); err != nil {
				klog.V(4).Info("apply rolebinding %q failed: %v\n", rb.Name, err)
			}
		case *rbacv1.ClusterRole:
			handler, err := clusterrole.New(ctx, kubeconfig)
			if err != nil {
				return err
			}
			if cr, err := handler.ApplyFromBytes(k8sResource); err != nil {
				klog.V(4).Info("apply clusterrole %q failed: %v\n", cr.Name, err)
			}
		case *rbacv1.ClusterRoleBinding:
			handler, err := clusterrolebinding.New(ctx, kubeconfig)
			if err != nil {
				return err
			}
			if crb, err := handler.ApplyFromBytes(k8sResource); err != nil {
				klog.V(4).Info("apply clusterrolebinding %q failed: %v\n", crb.Name, err)
			}
		case *corev1.ReplicationController:
			handler, err := replicationcontroller.New(ctx, "", kubeconfig)
			if err != nil {
				return err
			}
			if rc, err := handler.ApplyFromBytes(k8sResource); err != nil {
				klog.V(4).Info("apply replicationcontroller %q failed: %v\n", rc.Name, err)
			}
		case *appsv1.ReplicaSet:
			handler, err := replicaset.New(ctx, "", kubeconfig)
			if err != nil {
				return err
			}
			if rs, err := handler.ApplyFromBytes(k8sResource); err != nil {
				klog.V(4).Info("apply replicaset %q failed: %v\n", rs.Name, err)
			}
		case *storagev1.StorageClass:
			handler, err := storageclass.New(ctx, kubeconfig)
			if err != nil {
				return err
			}
			if sc, err := handler.ApplyFromBytes(k8sResource); err != nil {
				klog.V(4).Info("apply storageclass %q failed: %v\n", sc.Name, err)
			}
		default:
		}
	}

	return nil
}

func Decode(data []byte) (runtime.Object, error) {
	decode := scheme.Codecs.UniversalDeserializer().Decode
	object, _, err := decode(data, nil, nil)
	return object, err
}
