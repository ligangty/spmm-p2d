package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
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
	// rootCmd.AddCommand(newLoginCmd())
	rootCmd.AddCommand(PodCmd())
	rootCmd.AddCommand(listPVCsCmd())
	rootCmd.AddCommand(listConfigMapsCmd())
	rootCmd.AddCommand(listDeploymentsCmd())
	rootCmd.AddCommand(listSTSCmd())
	return rootCmd
}
