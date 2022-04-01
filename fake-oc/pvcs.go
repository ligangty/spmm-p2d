package main

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/homedir"
)

func listPVCsCmd() *cobra.Command {
	var kubeCfg, project string
	podsCmd := &cobra.Command{
		Use:   "pvcs",
		Short: "List PersistentVolumesClaim",
		Run: func(cmd *cobra.Command, args []string) {
			doListPVCs(kubeCfg, project)
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

func doListPVCs(kubeconfig, project string) {
	clientset := GetKubeClient(kubeconfig)
	namespace := project
	pvcs, err := clientset.CoreV1().PersistentVolumeClaims(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("There are %d Persistent Volume claims in the project %s\n\n", len(pvcs.Items), namespace)
	for _, p := range pvcs.Items {
		fmt.Print(p.ObjectMeta.Name + "\t")
		fmt.Print(p.Status.Phase)
		fmt.Println()
	}

}
