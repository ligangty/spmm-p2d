/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/go-logr/logr"
	testv1alpha1 "github.com/ligangty/helloworld-operator/api/v1alpha1"
)

// HelloReconciler reconciles a Hello object
type HelloReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=test.ligangty.github.com,resources=hellos,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=test.ligangty.github.com,resources=hellos/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=test.ligangty.github.com,resources=hellos/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Hello object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *HelloReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// TODO(user): your logic here

	hello := &testv1alpha1.Hello{}
	err := r.Get(ctx, req.NamespacedName, hello)

	if err != nil {
		logger.Error(err, "Failed to get Hello resource")
		return ctrl.Result{}, err
	}

	labels := make(map[string]string)
	labels["app.kubernetes.io/name"] = hello.Name
	if hello.Labels != nil {
		for labelKey, labelValue := range hello.Labels {
			labels[labelKey] = labelValue
		}
	}

	configmaps, err := newConfigMap(hello, r.Client, labels, logger)
	if err != nil {
		logger.Error(err, "Cannot create configmaps")
		return ctrl.Result{}, err
	}
	deployment := newDeploymet(hello, configmaps, labels, 1)

	logger.Info("Start creating Hello deployment")
	spacedName := types.NamespacedName{
		Namespace: deployment.ObjectMeta.Namespace,
		Name:      deployment.ObjectMeta.Name,
	}
	if err := r.Client.Get(context.TODO(), spacedName, deployment); err != nil {
		if errors.IsNotFound(err) {
			if err := r.Client.Create(context.TODO(), deployment); err != nil {
				logger.Error(err, "Cannot create Hello deployment")
				return ctrl.Result{}, err
			}
		} else {
			return ctrl.Result{}, err
		}
	}

	logger.Info("Hello deployment created")
	return ctrl.Result{Requeue: true}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *HelloReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&testv1alpha1.Hello{}).
		Owns(&appsv1.Deployment{}).
		Complete(r)
}

func newDeploymet(h *testv1alpha1.Hello, configmaps []*corev1.ConfigMap, labels map[string]string, desiredReplicaSize int32) *appsv1.Deployment {
	replicas := desiredReplicaSize
	image := "quay.io/ligangty/helloservice:latest"

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
								Value: "/var/www/template.html",
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
	return deployment
}

func newConfigMap(h *testv1alpha1.Hello, client client.Client, labels map[string]string, logger logr.Logger) ([]*corev1.ConfigMap, error) {
	cfg := &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "ConfigMap",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      h.Name,
			Namespace: h.Namespace,
			Labels:    labels,
		},
		Data: map[string]string{
			"template.html": `
			<html>
				<head></head>
				<body>
					<h1>Hello, ${message}!</h1>
				</body>
			</html>

			`,
			"template.json": `
			{
				"message": "Hello, ${message}!"
			}
			
			`,
		},
	}
	spacedName := types.NamespacedName{
		Namespace: cfg.ObjectMeta.Namespace,
		Name:      cfg.ObjectMeta.Name,
	}
	if err := client.Get(context.TODO(), spacedName, cfg); err != nil {
		if errors.IsNotFound(err) {
			logger.Info("Resource not found")
			if err = client.Create(context.TODO(), cfg); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	return []*corev1.ConfigMap{cfg}, nil
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
