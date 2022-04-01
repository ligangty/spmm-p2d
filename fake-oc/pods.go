package main

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/homedir"
)

func PodCmd() *cobra.Command {
	podCmd := &cobra.Command{
		Use:   "pod",
		Short: "Pod command",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	podCmd.AddCommand(listPodsCmd())
	podCmd.AddCommand(getPodCmd())
	return podCmd
}

func listPodsCmd() *cobra.Command {
	var kubeCfg, project string
	podsCmd := &cobra.Command{
		Use:   "list",
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

func getPodCmd() *cobra.Command {
	var kubeCfg, project string
	podsCmd := &cobra.Command{
		Use:   "get",
		Short: "Get Pod",
		Run: func(cmd *cobra.Command, args []string) {
			doGetPod(kubeCfg, args[0], project)
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

func doGetPod(kubeconfig, podname, project string) {
	clientset := GetKubeClient(kubeconfig)
	namespace := project
	pod, err := clientset.CoreV1().Pods(namespace).Get(context.TODO(), podname, metav1.GetOptions{})
	if err != nil {
		panic(err.Error())
	}
	data, err := pod.Marshal()
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(string(data))
}
