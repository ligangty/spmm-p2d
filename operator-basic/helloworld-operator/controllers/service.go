package controllers

import (
	"context"

	"github.com/go-logr/logr"
	testv1alpha1 "github.com/ligangty/helloworld-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func getOrNewService(ctx context.Context, h *testv1alpha1.Hello, client client.Client, labels map[string]string, logger logr.Logger) (*corev1.Service, error) {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceName(h),
			Namespace: h.Namespace,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			Type:      corev1.ServiceTypeClusterIP,
			Selector:  labels,
			ClusterIP: corev1.ClusterIPNone,
			Ports: []corev1.ServicePort{
				{
					Name: "http",
					Port: 8080,
				},
			},
		},
	}
	logger.Info("Start creating Hello Service")
	spacedName := types.NamespacedName{
		Namespace: service.ObjectMeta.Namespace,
		Name:      service.ObjectMeta.Name,
	}
	if err := client.Get(ctx, spacedName, service); err != nil {
		if errors.IsNotFound(err) {
			if err := client.Create(ctx, service); err != nil {
				logger.Error(err, "Cannot create Hello Service!")
				return nil, err
			}
			logger.Info("Hello Service created")
		} else {
			return nil, err
		}
	}
	return service, nil
}

func serviceName(h *testv1alpha1.Hello) string {
	return h.Name + "-service"
}
