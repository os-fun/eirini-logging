package main

import (
	"context"
	"net/http"

	"github.com/pkg/errors"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	eirinix "github.com/SUSE/eirinix"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	v1alpha1 "k8s.io/api/rbac/v1alpha1"
	rbac "k8s.io/client-go/kubernetes/typed/rbac/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission/types"
)

type Extension struct{ Namespace string }

func getVolume(name, path string) (v1.Volume, v1.VolumeMount) {
	mount := v1.VolumeMount{
		Name:      name,
		MountPath: path,
	}

	vol := v1.Volume{
		Name: name,
	}

	return vol, mount
}

func (ext *Extension) Handle(ctx context.Context, eiriniManager eirinix.Manager, pod *corev1.Pod, req types.Request) types.Response {

	if pod == nil {
		return admission.ErrorResponse(http.StatusBadRequest, errors.New("No pod could be decoded from the request"))
	}

	config, err := eiriniManager.GetKubeConnection()
	if err != nil {
		return admission.ErrorResponse(http.StatusBadRequest, errors.Wrap(err, "Failed getting the Kube connection"))
	}

	rbacClient, err := rbac.NewForConfig(config)
	if err != nil {
		return admission.ErrorResponse(http.StatusBadRequest, errors.Wrap(err, "Failed Creating RBAC Client"))
	}

	_, err = rbacClient.Roles(ext.Namespace).Create(&v1alpha1.Role{
		TypeMeta:   metav1.TypeMeta{Kind: "Role", APIVersion: "rbac.authorization.k8s.io/v1"},
		ObjectMeta: metav1.ObjectMeta{Namespace: ext.Namespace, Name: "role-" + pod.Name},
		Rules: []v1alpha1.PolicyRule{
			{
				ResourceNames: []string{pod.Name},
				Verbs:         []string{"get"},
				Resources:     []string{"pods", "pods/log"},
				APIGroups:     []string{""},
			},
		}})
	if err != nil {
		return admission.ErrorResponse(http.StatusBadRequest, errors.Wrap(err, "Failed Creating RBAC Role"))
	}

	_, err = rbacClient.RoleBindings(ext.Namespace).Create(&v1alpha1.RoleBinding{
		TypeMeta:   metav1.TypeMeta{Kind: "RoleBinding", APIVersion: "rbac.authorization.k8s.io/v1"},
		ObjectMeta: metav1.ObjectMeta{Namespace: ext.Namespace, Name: "role-binding-" + pod.Name},
		Subjects:   []v1alpha1.Subject{{Kind: "ServiceAccount", Name: "default", Namespace: ext.Namespace}},
		RoleRef: v1alpha1.RoleRef{
			Kind:     "Role",
			Name:     "role-" + pod.Name,
			APIGroup: "rbac.authorization.k8s.io",
		}})
	if err != nil {
		return admission.ErrorResponse(http.StatusBadRequest, errors.Wrap(err, "Failed Creating RBAC RoleBinding"))
	}

	podCopy := pod.DeepCopy()
	sidecar := corev1.Container{
		Name:    "eirini-logging",
		Image:   "cfcontainerization/drone-ci-image",
		Command: []string{"/bin/bash", "-c", "kubectl logs -f " + pod.Name + " -n " + ext.Namespace},
	}

	podCopy.Spec.Containers = append(podCopy.Spec.Containers, sidecar)

	return admission.PatchResponse(pod, podCopy)
}
