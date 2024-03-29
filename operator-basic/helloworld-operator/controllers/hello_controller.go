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
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	testv1alpha1 "github.com/ligangty/helloworld-operator/api/v1alpha1"
	routev1 "github.com/openshift/api/route/v1"
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

	// Configmaps
	configmaps, err := getOrNewConfigMap(ctx, hello, r.Client, labels, logger)
	if err != nil {
		logger.Error(err, "Cannot create Configmaps!")
		return ctrl.Result{}, err
	}

	// Deployment
	_, err = getOrNewDeploymet(ctx, hello, r.Client, configmaps, labels, 1, logger)
	if err != nil {
		logger.Error(err, "Cannot create Deployment!")
		return ctrl.Result{}, err
	}

	// Service
	svc, err := getOrNewService(ctx, hello, r.Client, labels, logger)
	_ = svc //WARNING: ignore declared but not used error
	if err != nil {
		logger.Error(err, "Cannot create Service!")
		return ctrl.Result{}, err
	}

	// Route
	_, err = getOrNewRoute(ctx, hello, r.Client, svc.ObjectMeta.Name, r.Scheme, labels, logger)
	if err != nil {
		logger.Error(err, "Cannot create Route!")
		return ctrl.Result{}, err
	}

	return ctrl.Result{Requeue: true}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *HelloReconciler) SetupWithManager(mgr ctrl.Manager) error {
	// /Add route scheme
	if err := routev1.AddToScheme(mgr.GetScheme()); err != nil {
		return err
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&testv1alpha1.Hello{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.ConfigMap{}).
		Owns(&corev1.Service{}).
		Owns(&routev1.Route{}).
		Complete(r)
}
