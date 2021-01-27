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
	"context"
	"fmt"
	"github.com/gosuri/uitable"

	"github.com/cyka/kubectl-java/util"
	"github.com/spf13/cobra"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

var (
	// command info
	listUsage   = "list [flags]"
	listShort   = "List Pods That Running Java Application"
	listLong    = "List Pods That Running Java Application"
	listExample = `
	just get pods that running java application
	`
)

var (
	headers = []interface{}{"NAME", "NODE", "STATUS", "CONTAINERS"}
)

type JavaPodFinder struct {
	options *KubeJavaAppOptions

	currentNameSpace string

	genericclioptions.IOStreams
}

//New kubectl-java list sub cmd
func NewListCmd(finder *JavaPodFinder) *cobra.Command {
	cmd := &cobra.Command{
		Use:     listUsage,
		Short:   listShort,
		Long:    listLong,
		Example: listExample,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			//TODO cmd should be processed by these step
			_ = finder.Complete(cmd)
			_ = finder.Validate()
			err = finder.Run()
			return
		},
	}
	return cmd
}

func NewJavaPodFinder(IOStreams genericclioptions.IOStreams, options *KubeJavaAppOptions) *JavaPodFinder {
	return &JavaPodFinder{
		options:   options,
		IOStreams: IOStreams,
	}
}

// handle user flags
func (f *JavaPodFinder) Complete(cmd *cobra.Command) error {
	namespace, err := cmd.Flags().GetString("namespace")
	if err != nil {
		return err
	}
	f.currentNameSpace = namespace
	f.printKubeConfigInfo()
	return nil
}

// validate the provided flags
func (f *JavaPodFinder) Validate() error {
	return nil
}

// execute the cmd
func (f *JavaPodFinder) Run() error {
	pods, err := f.findJavaPods()
	if err != nil {
		return err
	}
	tableToPrint := buildTableToPrint(pods)
	err = f.print(tableToPrint)
	return err
}

// print user kubeconfig
func (f *JavaPodFinder) printKubeConfigInfo() {
	kubConfig := f.options.userKubConfig
	currentContext, currentNameSpace, masterURL := util.GetCurrentConfigInfo(kubConfig)
	// namespace from user kubeconfig
	if len(f.currentNameSpace) == 0 {
		f.currentNameSpace = currentNameSpace
	}
	fmt.Printf("context:%s\tnameSpace:%s\tmaserURL:%s\n", util.Yellow(currentContext), util.Yellow(f.currentNameSpace), util.Yellow(masterURL))
}

// find pods that running java  application
func (f *JavaPodFinder) findJavaPods() ([]corev1.Pod, error) {
	podInfo, err := f.options.clientSet.
		CoreV1().
		Pods(f.currentNameSpace).
		List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	pods := podInfo.Items
	//TODO check the java pod and do filter
	return pods, nil
}

// print table to the console
func (f *JavaPodFinder) print(table *uitable.Table) error {
	// TODO add flags to control width
	table.MaxColWidth = 80
	_, err := fmt.Fprintln(f.Out, table)
	return err
}

// build table for printer
func buildTableToPrint(pods []corev1.Pod) *uitable.Table {
	length := len(headers)
	table := uitable.New()
	table.AddRow(headers...)
	rows := make([]metav1.TableRow, len(pods))
	for i, pod := range pods {
		podStatus := pod.Status
		containerStatuses := podStatus.ContainerStatuses
		row, containers := make([]interface{}, length, length), make([]string, len(containerStatuses))
		for index, status := range containerStatuses {
			containers[index] = util.HiCyan(status.Name)
		}
		// column: name
		row[0] = pod.Name
		// column: node
		row[1] = pod.Spec.NodeName
		// column: status
		row[2] = util.ColorizePodStatus(podStatus.Phase)
		// column: containers
		row[3] = containers
		// fill table cells with row
		rows[i].Cells = row
		table.AddRow(row...)
	}
	return table
}
