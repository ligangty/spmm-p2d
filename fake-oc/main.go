package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	rootCmd := newRootCmd()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func newRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "fake-oc",
		Short: "fake-oc is cli tools to fake openshift oc tools",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	rootCmd.AddCommand(newLoginCmd())
	return rootCmd
}

func newLoginCmd() *cobra.Command {
	loginCmd := &cobra.Command{
		Use:   "login ${kb_api_server}",
		Short: "Login to kube server",
		Run: func(cmd *cobra.Command, args []string) {
			user := getUsername(func() {
				cmd.Usage()
			})
			pass := getPassword(func() {
				cmd.Usage()
			})
			fmt.Printf("%s%s", user, pass)
		},
	}
	return loginCmd
}

func getUsername(usage func()) string {
	var username string
	prompt := promptui.Prompt{
		Label: "Please input your username",
	}

	uname, err := prompt.Run()
	if err != nil {
		fmt.Printf("Username input failed: %v\n", err)
		panic(err)
	}

	username = uname

	if IsEmptyString(username) {
		fmt.Println("Can not login to get access_token: username is missing")
		usage()
		os.Exit(1)
	}

	return username
}

func getPassword(usage func()) string {
	var password string
	prompt := promptui.Prompt{
		Label: "Please input your password",
		Mask:  '*',
	}

	pass, err := prompt.Run()

	if err != nil {
		fmt.Printf("Password input failed: %v\n", err)
		panic(err)
	}

	password = pass

	if IsEmptyString(password) {
		fmt.Println("Can not login to get access_token: password is missing")
		usage()
	}

	return password
}

func IsEmptyString(str string) bool {
	return len(strings.TrimSpace(str)) == 0
}

func example() {
	var kubeconfig *string
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	deploymentName := flag.String("deployment", "", "deployment name")
	imageName := flag.String("image", "", "new image name")
	appName := flag.String("app", "app", "application name")

	flag.Parse()
	if *deploymentName == "" {
		fmt.Println("You must specify the deployment name.")
		os.Exit(0)
	}
	if *imageName == "" {
		fmt.Println("You must specify the new image name.")
		os.Exit(0)
	}
	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	deployment, err := clientset.AppsV1beta1().Deployments("default").Get(*deploymentName, metav1.GetOptions{})
	if err != nil {
		panic(err.Error())
	}
	if errors.IsNotFound(err) {
		fmt.Printf("Deployment not found\n")
	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		fmt.Printf("Error getting deployment%v\n", statusError.ErrStatus.Message)
	} else if err != nil {
		panic(err.Error())
	} else {
		fmt.Printf("Found deployment\n")
		name := deployment.GetName()
		fmt.Println("name ->", name)
		containers := &deployment.Spec.Template.Spec.Containers
		found := false
		for i := range *containers {
			c := *containers
			if c[i].Name == *appName {
				found = true
				fmt.Println("Old image ->", c[i].Image)
				fmt.Println("New image ->", *imageName)
				c[i].Image = *imageName
			}
		}
		if found == false {
			fmt.Println("The application container not exist in the deployment pods.")
			os.Exit(0)
		}
		_, err := clientset.AppsV1beta1().Deployments("default").Update(deployment)
		if err != nil {
			panic(err.Error())
		}
	}
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
