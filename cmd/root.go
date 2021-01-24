package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/util/homedir"

	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

var (
	usage            = "kubectl-java [command]"
	version          = "v1.0.0"
	shortDescription = "The Kubectl Plugin For Java Application"
	longDescription  = "The Kubectl Plugin For Java Application"
	exampleUsage     = `
	# list pods that running java application
	%[1]s java list`
)

type KubeJavaAppOptions struct {
	configFlags   *genericclioptions.ConfigFlags
	userKubConfig clientcmdapi.Config
	genericclioptions.IOStreams
}

func NewKubeJavaAppOptions(IOStreams genericclioptions.IOStreams) *KubeJavaAppOptions {
	flags := genericclioptions.NewConfigFlags(true)
	flags.KubeConfig = getLocalKubeConfigPath()
	return &KubeJavaAppOptions{
		configFlags: flags,
		IOStreams:   IOStreams,
	}
}

func NewKubeJavaCmd(streams genericclioptions.IOStreams) *cobra.Command {
	options := NewKubeJavaAppOptions(streams)
	podFinder := NewJavaPodFinder(streams, options)

	rootCmd := &cobra.Command{
		Use:     usage,
		Short:   shortDescription,
		Long:    longDescription,
		Version: version,
		Example: fmt.Sprintf(exampleUsage, "kubectl"),
		PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
			options.userKubConfig, err = options.configFlags.ToRawKubeConfigLoader().RawConfig()
			return
		},
	}
	// add flags
	options.configFlags.AddFlags(rootCmd.PersistentFlags())
	// find java pod cmd
	rootCmd.AddCommand(NewListCmd(podFinder))

	return rootCmd
}

func build() {

}

func getLocalKubeConfigPath() *string {
	var kubeConfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeConfig = stringPtr(filepath.Join(home, ".kube", "config"))
	}
	return kubeConfig
}

func stringPtr(val string) *string {
	return &val
}
