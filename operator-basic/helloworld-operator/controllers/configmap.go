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

func getOrNewConfigMap(ctx context.Context, h *testv1alpha1.Hello, client client.Client, labels map[string]string, logger logr.Logger) ([]*corev1.ConfigMap, error) {
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
	logger.Info("Start creating Hello configmaps")
	spacedName := types.NamespacedName{
		Namespace: cfg.ObjectMeta.Namespace,
		Name:      cfg.ObjectMeta.Name,
	}
	if err := client.Get(ctx, spacedName, cfg); err != nil {
		if errors.IsNotFound(err) {
			logger.Info("Configmap not found")
			if err = client.Create(context.TODO(), cfg); err != nil {
				return nil, err
			}
			logger.Info("Configmap created")
		} else {
			return nil, err
		}
	}

	return []*corev1.ConfigMap{cfg}, nil
}
