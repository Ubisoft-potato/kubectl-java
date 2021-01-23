package main

import (
	"os"

	"github.com/cyka/kubectl-java/cmd"
	"github.com/spf13/pflag"

	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func main() {
	flags := pflag.NewFlagSet("kubectl-java", pflag.ExitOnError)
	pflag.CommandLine = flags

	javaCmd := cmd.NewKubeJavaCmd(genericclioptions.IOStreams{
		In:     os.Stdin,
		Out:    os.Stdout,
		ErrOut: os.Stderr,
	})

	if err := javaCmd.Execute(); err != nil {
		//todo handle err here
		os.Exit(1)
	}
}
