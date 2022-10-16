package controllers

import (
	"context"
	"strings"

	"github.com/go-logr/logr"
	testv1alpha1 "github.com/ligangty/helloworld-operator/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func getOrNewDeploymet(ctx context.Context, h *testv1alpha1.Hello, client client.Client, configmaps []*corev1.ConfigMap, labels map[string]string, desiredReplicaSize int32, logger logr.Logger) (*appsv1.Deployment, error) {
	replicas := desiredReplicaSize
	image := "quay.io/ligangty/helloservice:latest"
	templateFile := h.Spec.TemplateFile
	if strings.TrimSpace(templateFile) == "" {
		templateFile = "/var/www/template.html"
	}

	deployment := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      h.Name,
			Namespace: h.Namespace,
			Labels:    labels,
			Annotations: map[string]string{
				"image.openshift.io/triggers": "[{ \"from\": { \"kind\":\"ImageStreamTag\", \"name\":\"" + image + "\"}, \"fieldPath\": \"spec.template.spec.containers[?(@.name==\\\"" + h.Name + "\\\")].image\"}]",
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      labels,
					Annotations: map[string]string{},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Name:  h.Name,
						Image: image,
						Ports: []corev1.ContainerPort{
							{
								ContainerPort: 8080,
								Name:          "http",
							},
						},
						Env: []corev1.EnvVar{
							{
								Name:  "TEMPLATE_PATH",
								Value: templateFile,
							},
						},
						VolumeMounts: newConfigMounts(configmaps),
						// // Resources
						// Resources: createResources(w.Spec.Resources),
					}},
					Volumes: newConfigVolumes(configmaps),
					// ServiceAccountName: h.Spec.ServiceAccountName,
				},
			},
		},
	}
	logger.Info("Start creating Hello deployment")
	spacedName := types.NamespacedName{
		Namespace: deployment.ObjectMeta.Namespace,
		Name:      deployment.ObjectMeta.Name,
	}
	if err := client.Get(ctx, spacedName, deployment); err != nil {
		if errors.IsNotFound(err) {
			if err := client.Create(ctx, deployment); err != nil {
				logger.Error(err, "Cannot create Hello deployment")
				return nil, err
			}
			logger.Info("Hello deployment created")
		} else {
			return nil, err
		}
	}
	return deployment, nil
}

func newConfigMounts(cfs []*corev1.ConfigMap) []corev1.VolumeMount {
	vols := []corev1.VolumeMount{}
	for _, cf := range cfs {
		vol := corev1.VolumeMount{
			MountPath: "/var/www/",
			Name:      "vol-" + cf.Name,
			ReadOnly:  true,
		}
		vols = append(vols, vol)
	}
	return vols
}

func newConfigVolumes(cfs []*corev1.ConfigMap) []corev1.Volume {
	vols := []corev1.Volume{}
	var mode *int32 = new(int32)
	*mode = 0420
	for _, cf := range cfs {
		vol := corev1.Volume{
			Name: "vol-" + cf.Name,
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: cf.Name,
					},
					DefaultMode: mode,
				},
			},
		}
		vols = append(vols, vol)
	}

	return vols
}
