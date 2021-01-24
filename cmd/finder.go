package cmd

import (
	"fmt"
	"github.com/cyka/kubectl-java/util"
	"github.com/jedib0t/go-pretty/table"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

var (
	// command info
	listUsage   = "list [flags]"
	listShort   = "List All Pods That Running Java Application"
	listLong    = "List All Pods That Running Java Application"
	listExample = ""

	// table printer
	writer = table.NewWriter()
)

type JavaPodFinder struct {
	genericclioptions.IOStreams
	options *KubeJavaAppOptions
}

func NewJavaPodFinder(IOStreams genericclioptions.IOStreams, options *KubeJavaAppOptions) *JavaPodFinder {
	return &JavaPodFinder{IOStreams: IOStreams, options: options}
}

func NewListCmd(finder *JavaPodFinder) *cobra.Command {

	cmd := &cobra.Command{
		Use:     listUsage,
		Short:   listShort,
		Long:    listLong,
		Example: listExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			finder.configTablePrint()
			fmt.Fprintln(finder.Out, writer.Render())
			return nil
		},
	}
	return cmd
}

func (f JavaPodFinder) configTablePrint() {
	kubConfig := f.options.userKubConfig
	writer.AppendHeader(table.Row{"context", "namespace"})
	_, _ = fmt.Fprintf(f.Out, "context: %s\t\tnamespace: %s\n",
		util.Cyan(util.GetCurrentContext(kubConfig)),
		util.Yellow(util.GetCurrentNameSpace(kubConfig)))
}
