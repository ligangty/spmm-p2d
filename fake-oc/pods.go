package main

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/homedir"
)

func listPodsCmd() *cobra.Command {
	var kubeCfg, project string
	podsCmd := &cobra.Command{
		Use:   "pods",
		Short: "List pods",
		Run: func(cmd *cobra.Command, args []string) {
			doListPods(kubeCfg, project)
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
	clientset := GetKubeClient(kubeconfig)
	namespace := project
	pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("There are %d pods in the project %s\n\n", len(pods.Items), namespace)
	for _, p := range pods.Items {
		fmt.Print(p.ObjectMeta.Name + "\t")
		fmt.Print(p.Status.Phase)
		fmt.Println()
	}

}
