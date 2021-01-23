package cmd

import (
	"fmt"
	"github.com/spf13/cobra"

	"k8s.io/cli-runtime/pkg/genericclioptions"

	"k8s.io/client-go/util/homedir"
)

var (
	exampleUsage = `
	# list pods that running java application
	%[1]s java list
`
)

type KubeJavaAppOptions struct {
	configFlags *genericclioptions.ConfigFlags
	//todo

	genericclioptions.IOStreams
}

func NewKubeJavaAppOptions(IOStreams genericclioptions.IOStreams) *KubeJavaAppOptions {
	return &KubeJavaAppOptions{
		configFlags: genericclioptions.NewConfigFlags(true),
		IOStreams:   IOStreams,
	}
}

func NewKubeJavaCmd(streams genericclioptions.IOStreams) *cobra.Command {
	options := NewKubeJavaAppOptions(streams)

	rootCmd := &cobra.Command{
		Use:     "kubectl-java",
		Long:    "",
		Example: fmt.Sprintf(exampleUsage, "kubectl"),
		Short:   "The Kubectl Plugin For Java Application",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, _ = fmt.Fprintln(options.IOStreams.Out, "hello kubectl!")
			return nil
		},
	}

	options.configFlags.AddFlags(rootCmd.Flags())

	return rootCmd
}

func init() {
	fmt.Printf("current home: %s\n", homedir.HomeDir())
}
