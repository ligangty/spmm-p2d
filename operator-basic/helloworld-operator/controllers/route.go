package controllers

import (
	"context"

	"github.com/go-logr/logr"
	testv1alpha1 "github.com/ligangty/helloworld-operator/api/v1alpha1"
	routev1 "github.com/openshift/api/route/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// // newRoute either returns the headless service or create it
// func newRoute(h *testv1alpha1.Hello, client client.Client, scheme *runtime.Scheme, labels map[string]string) (*routev1.Route, error) {
// 	route := &routev1.Route{}
// 	if err := client.Get(w, types.NamespacedName{Name: routeServiceName(w), Namespace: w.Namespace}, client, route); err != nil {
// 		if errors.IsNotFound(err) {
// 			if err := resources.Create(w, client, scheme, newRoute(w, labels)); err != nil {
// 				if errors.IsAlreadyExists(err) {
// 					return nil, nil
// 				}
// 				return nil, err
// 			}
// 			return nil, nil
// 		}
// 	}
// 	return route, nil
// }

// // DeleteExistingRoute delete the route if it exists. It returns true if the route is deleted.
// func DeleteExistingRoute(h *testv1alpha1.Hello, client client.Client) (bool, error) {
// 	route := &routev1.Route{}
// 	if err := resources.Get(w, types.NamespacedName{Name: routeServiceName(w), Namespace: w.Namespace}, client, route); err != nil {
// 		if errors.IsNotFound(err) {
// 			// route has been not found, nothing to do
// 			return false, nil
// 		}
// 		return false, err
// 	}
// 	// route has been found, let's delete it
// 	if err := resources.Delete(w, client, route); err != nil {
// 		return true, err
// 	}
// 	return false, nil
// }

func getOrNewRoute(ctx context.Context, h *testv1alpha1.Hello, client client.Client, serviceName string, scheme *runtime.Scheme, labels map[string]string, logger logr.Logger) (*routev1.Route, error) {
	weight := int32(100)

	route := &routev1.Route{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "route.openshift.io/v1",
			Kind:       "Route",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      routeServiceName(h),
			Namespace: h.Namespace,
			Labels:    labels,
		},
		Spec: routev1.RouteSpec{
			To: routev1.RouteTargetReference{
				Kind:   "Service",
				Name:   serviceName,
				Weight: &weight,
			},
			Port: &routev1.RoutePort{
				TargetPort: intstr.FromString("http"),
			},
		},
	}

	// Need to register scheme for route
	if err := controllerutil.SetControllerReference(h, &route.ObjectMeta, scheme); err != nil {
		logger.Error(err, "Failed to set controller reference for new route")
		return nil, err
	}
	logger.Info("Start creating Hello Route")
	spacedName := types.NamespacedName{
		Namespace: route.ObjectMeta.Namespace,
		Name:      route.ObjectMeta.Name,
	}
	if err := client.Get(ctx, spacedName, route); err != nil {
		if errors.IsNotFound(err) {
			if err := client.Create(ctx, route); err != nil {
				logger.Error(err, "Cannot create Hello Route")
				return nil, err
			}
			logger.Info("Hello Route created")
		} else {
			return nil, err
		}
	}
	return route, nil
}

// routeServiceName returns the name of the HTTP route
func routeServiceName(h *testv1alpha1.Hello) string {
	return h.Name + "-route"
}
