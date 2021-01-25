/*
 * Copyright (c) 2020-2021 Ubisoft-Potato
 *
 * Anti 996 License Version 1.0 (Draft)
 *
 * Permission is hereby granted to any individual or legal entity obtaining a copy
 * of this licensed work (including the source code, documentation and/or related
 * items, hereinafter collectively referred to as the "licensed work"), free of
 * charge, to deal with the licensed work for any purpose, including without
 * limitation, the rights to use, reproduce, modify, prepare derivative works of,
 * publish, distribute and sublicense the licensed work, subject to the following
 * conditions:
 *
 * 1.  The individual or the legal entity must conspicuously display, without
 *     modification, this License on each redistributed or derivative copy of the
 *     Licensed Work.
 *
 * 2.  The individual or the legal entity must strictly comply with all applicable
 *     laws, regulations, rules and standards of the jurisdiction relating to
 *     labor and employment where the individual is physically located or where
 *     the individual was born or naturalized; or where the legal entity is
 *     registered or is operating (whichever is stricter). In case that the
 *     jurisdiction has no such laws, regulations, rules and standards or its
 *     laws, regulations, rules and standards are unenforceable, the individual
 *     or the legal entity are required to comply with Core International Labor
 *     Standards.
 *
 * 3.  The individual or the legal entity shall not induce or force its
 *     employee(s), whether full-time or part-time, or its independent
 *     contractor(s), in any methods, to agree in oral or written form,
 *     to directly or indirectly restrict, weaken or relinquish his or
 *     her rights or remedies under such laws, regulations, rules and
 *     standards relating to labor and employment as mentioned above,
 *     no matter whether such written or oral agreement are enforceable
 *     under the laws of the said jurisdiction, nor shall such individual
 *     or the legal entity limit, in any methods, the rights of its employee(s)
 *     or independent contractor(s) from reporting or complaining to the copyright
 *     holder or relevant authorities monitoring the compliance of the license
 *     about its violation(s) of the said license.
 *
 * THE LICENSED WORK IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
 * FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE COPYRIGHT
 * HOLDER BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
 * OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN ANY WAY CONNECTION
 * WITH THE LICENSED WORK OR THE USE OR OTHER DEALINGS IN THE LICENSED WORK.
 *
 */

package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/client-go/util/homedir"
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

// New kubectl-java main cmd
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
			err = initBeforeExecute(options)
			return
		},
	}
	// add flags
	options.configFlags.AddFlags(rootCmd.PersistentFlags())
	// find java pod cmd
	rootCmd.AddCommand(NewListCmd(podFinder))

	return rootCmd
}

//cmd options
type KubeJavaAppOptions struct {
	configFlags   *genericclioptions.ConfigFlags
	userKubConfig clientcmdapi.Config
	clientSet     *kubernetes.Clientset
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

func initBeforeExecute(options *KubeJavaAppOptions) error {
	kubeConfigLoader := options.configFlags.ToRawKubeConfigLoader()
	userKubConfig, rawErr := kubeConfigLoader.RawConfig()
	if rawErr != nil {
		return rawErr
	}
	restConfig, _ := kubeConfigLoader.ClientConfig()
	clientSet, clientErr := kubernetes.NewForConfig(restConfig)
	if clientErr != nil {
		return clientErr
	}
	options.userKubConfig = userKubConfig
	options.clientSet = clientSet
	return nil
}
