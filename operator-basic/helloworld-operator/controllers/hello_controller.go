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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

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

	deployment := newDeploymet(hello, labels, 1)

	logger.Info("Start creating Hello deployment")
	if err := r.Client.Create(context.TODO(), deployment); err != nil {
		logger.Error(err, "Cannot create Hello deployment")
		return ctrl.Result{}, err
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

func newDeploymet(h *testv1alpha1.Hello, labels map[string]string, desiredReplicaSize int32) *appsv1.Deployment {
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
						// // Resources
						// Resources: createResources(w.Spec.Resources),
					}},

					// ServiceAccountName: h.Spec.ServiceAccountName,
				},
			},
		},
	}
	return deployment
}
