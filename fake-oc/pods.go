package main

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func listPodsCmd() *cobra.Command {
	var kubeCfg, project string
	podsCmd := &cobra.Command{
		Use:   "pods",
		Short: "List pods",
		Run: func(cmd *cobra.Command, args []string) {
			doListPods(kubeCfg, project)
			// fmt.Printf("%s%s", user, pass)
		},
	}
	if home := homedir.HomeDir(); home != "" {
		podsCmd.Flags().StringVarP(&kubeCfg, "kubeconfig", "k", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		podsCmd.Flags().StringVarP(&kubeCfg, "kubeconfig", "k", "", "absolute path to the kubeconfig file")
	}
	podsCmd.Flags().StringVarP(&project, "project", "p", "", "project which the pods belong to")
	podsCmd.MarkFlagRequired("project")
	return podsCmd
}

func doListPods(kubeconfig, project string) {

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(config)

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	namespace := project
	pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))
	for _, p := range pods.Items {
		fmt.Print(p.ObjectMeta.Name + "\t")
		fmt.Print(p.Status.Phase)
		fmt.Println()
	}

}
