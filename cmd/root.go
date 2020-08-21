package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/klog"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "alert",
	Short: "error alert tool.",
	Long:  "error alert tool.",
}

func init() {
	klog.InitFlags(nil)
	flag.Parse()

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		klog.Fatalf(fmt.Sprintf("run root cmd err: %v", err))
		os.Exit(1)
	}
}
