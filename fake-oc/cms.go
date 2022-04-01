package main

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/homedir"
)

func listConfigMapsCmd() *cobra.Command {
	var kubeCfg, project string
	podsCmd := &cobra.Command{
		Use:   "cms",
		Short: "List ConfigMap",
		Run: func(cmd *cobra.Command, args []string) {
			doListConfigMaps(kubeCfg, project)
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

func doListConfigMaps(kubeconfig, project string) {
	clientset := GetKubeClient(kubeconfig)
	namespace := project
	cfgMaps, err := clientset.CoreV1().ConfigMaps(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("There are %d ConfigMaps in the project %s\n\n", len(cfgMaps.Items), namespace)
	for _, c := range cfgMaps.Items {
		fmt.Print(c.ObjectMeta.Name + "\t")
		fmt.Println()
	}

}
