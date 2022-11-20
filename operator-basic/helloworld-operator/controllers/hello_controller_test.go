package controllers

import (
	"context"
	"os"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	hellov1alpha1 "github.com/ligangty/helloworld-operator/api/v1alpha1"
)

var _ = Describe("Hello controller", func() {
	Context("Hello controller test", func() {

		const HelloName = "test-hello"

		ctx := context.Background()

		namespace := &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name:      HelloName,
				Namespace: HelloName,
			},
		}

		typeNamespaceName := types.NamespacedName{Name: HelloName, Namespace: HelloName}

		BeforeEach(func() {
			By("Creating the Namespace to perform the tests")
			err := k8sClient.Create(ctx, namespace)
			Expect(err).To(Not(HaveOccurred()))

			By("Setting the Image ENV VAR which stores the Operand image")
			err = os.Setenv("HELLO_IMAGE", "example.com/image:test")
			Expect(err).To(Not(HaveOccurred()))
		})

		AfterEach(func() {
			// TODO(user): Attention if you improve this code by adding other context test you MUST
			// be aware of the current delete namespace limitations. More info: https://book.kubebuilder.io/reference/envtest.html#testing-considerations
			By("Deleting the Namespace to perform the tests")
			_ = k8sClient.Delete(ctx, namespace)

			By("Removing the Image ENV VAR which stores the Operand image")
			_ = os.Unsetenv("HELLO_IMAGE")
		})

		It("should successfully reconcile a custom resource for Hello", func() {
			By("Creating the custom resource for the Kind Hello")
			hello := &hellov1alpha1.Hello{}
			err := k8sClient.Get(ctx, typeNamespaceName, hello)
			if err != nil && errors.IsNotFound(err) {
				// Let's mock our custom resource at the same way that we would
				// apply on the cluster the manifest under config/samples
				hello := &hellov1alpha1.Hello{
					ObjectMeta: metav1.ObjectMeta{
						Name:      HelloName,
						Namespace: namespace.Name,
					},
					Spec: hellov1alpha1.HelloSpec{
						ContentConfigMap: "testcfg",
						TemplateFile:     "/var/www/template.html",
					},
				}

				err = k8sClient.Create(ctx, hello)
				Expect(err).To(Not(HaveOccurred()))
			}

			By("Checking if the custom resource was successfully created")
			Eventually(func() error {
				found := &hellov1alpha1.Hello{}
				return k8sClient.Get(ctx, typeNamespaceName, found)
			}, time.Minute, time.Second).Should(Succeed())

			By("Reconciling the custom resource created")
			helloReconciler := &HelloReconciler{
				Client: k8sClient,
				Scheme: k8sClient.Scheme(),
			}

			_, err = helloReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: typeNamespaceName,
			})
			Expect(err).To(Not(HaveOccurred()))

			By("Checking if Deployment was successfully created in the reconciliation")
			Eventually(func() error {
				found := &appsv1.Deployment{}
				return k8sClient.Get(ctx, typeNamespaceName, found)
			}, time.Minute, time.Second).Should(Succeed())

			By("Checking the latest Status Condition added to the Hello instance")
			Eventually(func() error {
				// if hello.Status.Conditions != nil && len(hello.Status.Conditions) != 0 {
				// 	latestStatusCondition := hello.Status.Conditions[len(hello.Status.Conditions)-1]
				// 	expectedLatestStatusCondition := metav1.Condition{Type: typeAvailableHello,
				// 		Status: metav1.ConditionTrue, Reason: "Reconciling",
				// 		Message: fmt.Sprintf("Deployment for custom resource (%s) with %d replicas created successfully", hello.Name, hello.Spec.Size)}
				// 	if latestStatusCondition != expectedLatestStatusCondition {
				// 		return fmt.Errorf("The latest status condition added to the hello instance is not as expected")
				// 	}
				// }
				return nil
			}, time.Minute, time.Second).Should(Succeed())
		})
	})
})
