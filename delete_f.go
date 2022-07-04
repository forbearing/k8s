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
	"github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	networking "k8s.io/api/networking/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	storagev1 "k8s.io/api/storage/v1"
)

// DeleteF like "kubectl delete -f filename"
func DeleteF(ctx context.Context, kubeconfig, filename string) (err error) {
	k8sResourceFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return
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
			logrus.Errorf("Decode error: %v\n", err)
			continue
		}
		switch object.(type) {
		case *corev1.Namespace:
			handler, err := namespace.New(ctx, kubeconfig)
			if err != nil {
				return err
			}
			ns, _ := handler.GetFromBytes(k8sResource)
			if err := handler.DeleteFromBytes(k8sResource); err != nil {
				logrus.Errorf("delete namespace %q failed: %v\n", ns.Name, err)
			}
		case *corev1.Service:
			handler, err := service.New(ctx, "", kubeconfig)
			if err != nil {
				return err
			}
			svc, _ := handler.GetFromBytes(k8sResource)
			if err := handler.DeleteFromBytes(k8sResource); err != nil {
				logrus.Errorf("delete service %q failed: %v\n", svc.Name, err)
			}
		case *corev1.ConfigMap:
			handler, err := configmap.New(ctx, "", kubeconfig)
			if err != nil {
				return err
			}
			cm, _ := handler.GetFromBytes(k8sResource)
			if err := handler.DeleteFromBytes(k8sResource); err != nil {
				logrus.Errorf("delete configmap %q failed: %v\n", cm.Name, err)
			}
		case *corev1.Secret:
			handler, err := secret.New(ctx, "", kubeconfig)
			if err != nil {
				return err
			}
			sec, _ := handler.GetFromBytes(k8sResource)
			if err := handler.DeleteFromBytes(k8sResource); err != nil {
				logrus.Errorf("delete secret %q failed: %v\n", sec.Name, err)
			}
		case *corev1.ServiceAccount:
			handler, err := serviceaccount.New(ctx, "", kubeconfig)
			if err != nil {
				return err
			}
			sa, _ := handler.GetFromBytes(k8sResource)
			if err := handler.DeleteFromBytes(k8sResource); err != nil {
				logrus.Errorf("delete serviceaccount %q failed: %v\n", sa.Name, err)
			}
		case *corev1.Pod:
			handler, err := pod.New(ctx, "", kubeconfig)
			if err != nil {
				return err
			}
			p, _ := handler.GetFromBytes(k8sResource)
			if err := handler.DeleteFromBytes(k8sResource); err != nil {
				logrus.Errorf("delete pod %q failed: %v\n", p.Name, err)
			}
		case *corev1.PersistentVolume:
			handler, err := persistentvolume.New(ctx, kubeconfig)
			if err != nil {
				return err
			}
			pv, _ := handler.GetFromBytes(k8sResource)
			if err := handler.DeleteFromBytes(k8sResource); err != nil {
				logrus.Errorf("delete persistentvolume %q failed: %v\n", pv.Name, err)
			}
		case *corev1.PersistentVolumeClaim:
			handler, err := persistentvolumeclaim.New(ctx, "", kubeconfig)
			if err != nil {
				return err
			}
			pvc, _ := handler.GetFromBytes(k8sResource)
			if err := handler.DeleteFromBytes(k8sResource); err != nil {
				logrus.Errorf("delete persistentvolumeclaim %q failed: %v\n", pvc.Name, err)
			}
		case *appsv1.Deployment:
			handler, err := deployment.New(ctx, "", kubeconfig)
			if err != nil {
				return err
			}
			deploy, _ := handler.GetFromBytes(k8sResource)
			if err := handler.DeleteFromBytes(k8sResource); err != nil {
				logrus.Errorf("delete deployment %q failed: %v\n", deploy.Name, err)
			}
		case *appsv1.StatefulSet:
			handler, err := statefulset.New(ctx, "", kubeconfig)
			if err != nil {
				return err
			}
			sts, _ := handler.GetFromBytes(k8sResource)
			if err := handler.DeleteFromBytes(k8sResource); err != nil {
				logrus.Errorf("delete statefulset %q failed: %v\n", sts.Name, err)
			}
		case *appsv1.DaemonSet:
			handler, err := daemonset.New(ctx, "", kubeconfig)
			if err != nil {
				return err
			}
			ds, _ := handler.GetFromBytes(k8sResource)
			if err := handler.DeleteFromBytes(k8sResource); err != nil {
				logrus.Errorf("delete daemonset %q failed: %v\n", ds.Name, err)
			}
		case *networking.Ingress:
			handler, err := ingress.New(ctx, "", kubeconfig)
			if err != nil {
				return err
			}
			ing, _ := handler.GetFromBytes(k8sResource)
			if err := handler.DeleteFromBytes(k8sResource); err != nil {
				logrus.Errorf("delete ingress %q failed: %v\n", ing.Name, err)
			}
		case *networking.IngressClass:
			handler, err := ingressclass.New(ctx, kubeconfig)
			if err != nil {
				return err
			}
			ingc, _ := handler.GetFromBytes(k8sResource)
			if err := handler.DeleteFromBytes(k8sResource); err != nil {
				logrus.Errorf("delete ingressclass %q failed: %v\n", ingc.Name, err)
			}
		case *networking.NetworkPolicy:
			handler, err := networkpolicy.New(ctx, "", kubeconfig)
			if err != nil {
				return err
			}
			netpol, _ := handler.GetFromBytes(k8sResource)
			if err := handler.DeleteFromBytes(k8sResource); err != nil {
				logrus.Errorf("delete networkpolicy %q failed: %v\n", netpol.Name, err)
			}
		case *batchv1.Job:
			handler, err := job.New(ctx, "", kubeconfig)
			if err != nil {
				return err
			}
			job, _ := handler.GetFromBytes(k8sResource)
			if err := handler.DeleteFromBytes(k8sResource); err != nil {
				logrus.Errorf("delete job %q failed: %v\n", job.Name, err)
			}
		case *batchv1.CronJob:
			handler, err := cronjob.New(ctx, "", kubeconfig)
			if err != nil {
				return err
			}
			cj, _ := handler.GetFromBytes(k8sResource)
			if err := handler.DeleteFromBytes(k8sResource); err != nil {
				logrus.Errorf("delete cronjob %q failed: %v\n", cj.Name, err)
			}
		case *rbacv1.Role:
			handler, err := role.New(ctx, "", kubeconfig)
			if err != nil {
				return err
			}
			r, _ := handler.GetFromBytes(k8sResource)
			if err := handler.DeleteFromBytes(k8sResource); err != nil {
				logrus.Errorf("delete role %q failed: %v\n", r.Name, err)
			}
		case *rbacv1.RoleBinding:
			handler, err := rolebinding.New(ctx, "", kubeconfig)
			if err != nil {
				return err
			}
			rb, _ := handler.CreateFromBytes(k8sResource)
			if err := handler.DeleteFromBytes(k8sResource); err != nil {
				logrus.Errorf("delete rolebinding %q failed", rb.Name)
			}
		case *rbacv1.ClusterRole:
			handler, err := clusterrole.New(ctx, kubeconfig)
			if err != nil {
				return err
			}
			cr, _ := handler.GetFromBytes(k8sResource)
			if err := handler.DeleteFromBytes(k8sResource); err != nil {
				logrus.Errorf("delete clusterrole %q failed: %v\n", cr.Name, err)
			}
		case *rbacv1.ClusterRoleBinding:
			handler, err := clusterrolebinding.New(ctx, kubeconfig)
			if err != nil {
				return err
			}
			crb, _ := handler.GetFromBytes(k8sResource)
			if err := handler.DeleteFromBytes(k8sResource); err != nil {
				logrus.Errorf("delete clusterrolebinding %q failed: %v\n", crb.Name, err)
			}
		case *appsv1.ReplicaSet:
			handler, err := replicaset.New(ctx, "", kubeconfig)
			if err != nil {
				return err
			}
			rs, _ := handler.GetFromBytes(k8sResource)
			if err := handler.DeleteFromBytes(k8sResource); err != nil {
				logrus.Errorf("delete replicaset %q failed: %v\n", rs.Name, err)
			}
		case *corev1.ReplicationController:
			handler, err := replicationcontroller.New(ctx, "", kubeconfig)
			if err != nil {
				return err
			}
			rc, _ := handler.GetFromBytes(k8sResource)
			if err := handler.DeleteFromBytes(k8sResource); err != nil {
				logrus.Errorf("delete replicationcontroller %q failed: %v\n", rc.Name, err)
			}
		case *storagev1.StorageClass:
			handler, err := storageclass.New(ctx, kubeconfig)
			if err != nil {
				return err
			}
			sc, _ := handler.GetFromBytes(k8sResource)
			if err := handler.DeleteFromBytes(k8sResource); err != nil {
				logrus.Errorf("delete storageclass %q failed: %v\n", sc.Name, err)
			}
		default:
		}
	}

	return
}
