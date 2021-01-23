package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

var (
	listUsage   = "list [flags]"
	listShort   = "List All Pods That Running Java Application"
	listLong    = "List All Pods That Running Java Application"
	listExample = ""
)

type JavaPodFinder struct {
}

func NewListCmd(streams genericclioptions.IOStreams) *cobra.Command {
	cmd := &cobra.Command{
		Use:     listUsage,
		Short:   listShort,
		Long:    listLong,
		Example: listExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			_, _ = fmt.Fprintln(streams.Out, "todo")
			return nil
		},
	}
	return cmd
}
