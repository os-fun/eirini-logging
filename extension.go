package main

import (
	"context"

	eirinix "github.com/SUSE/eirinix"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission/types"
)

type Extension struct{}

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

	// if pod == nil {
	// 	return admission.ErrorResponse(http.StatusBadRequest, errors.New("No pod could be decoded from the request"))
	// }

	// config, err := eiriniManager.GetKubeConnection()
	// if err != nil {
	// 	return admission.ErrorResponse(http.StatusBadRequest, errors.New("No pod could be decoded from the request"))
	// }

	// client, err := corev1client.NewForConfig(config)
	// if err != nil {
	// 	return admission.ErrorResponse(http.StatusBadRequest, errors.New("No pod could be decoded from the request"))
	// }

	// client.RESTClient().Get().URL().Host=

	podCopy := pod.DeepCopy()
	sidecar := corev1.Container{}
	sidecar.Name = "scf"
	sidecar.Image = "busybox"
	sidecar.Command = []string{"/bin/sh", "-c", "sleep 3450404030"}
	// if len(podCopy.Spec.ServiceAccountName) == 0 {
	// 	podCopy.Spec.ServiceAccountName = "default"

	// }

	//      /proc/1/fd/1
	for i := range podCopy.Spec.Containers {
		c := &podCopy.Spec.Containers[i]

		c.VolumeMounts = append(c.VolumeMounts, corev1.VolumeMount{Name: "sharedstdout", MountPath: "/var/log/eiriniapp"})

		c.Args = append(c.Args, "|", "tee", "/var/log/eiriniapp/logs")

		//c.Command = append(c.Command)
	}
	podCopy.Spec.Volumes = append(podCopy.Spec.Volumes, corev1.Volume{Name: "sharedstdout"})

	sidecar.VolumeMounts = append(sidecar.VolumeMounts, corev1.VolumeMount{Name: "sharedstdout", MountPath: "/var/log/eiriniapp"})

	podCopy.Spec.Containers = append(podCopy.Spec.Containers, sidecar)

	return admission.PatchResponse(pod, podCopy)
}
