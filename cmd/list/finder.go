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

package list

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/cyka/kubectl-java/util"
	"github.com/gosuri/uitable"
	"github.com/spf13/cobra"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

var (
	// command info
	listUsage   = "list [flags]"
	listShort   = "List Pods That Running Java Application"
	listLong    = "List Pods That Running Java Application"
	listExample = `
	# just get pods that running java application
	%[1]s java list
	# specify the namespace
	%[1]s java list -n namespace
	# custom table column width
	%[1]s java list -w 80`
)

var (
	headers = []interface{}{"NAME", "NODE", "STATUS", "CONTAINERS", "JDK"}
	wg      = sync.WaitGroup{}
)

type JavaPodFinder struct {
	cmdFactory *util.CmdFactory
	executor   util.RemoteExecutor

	nameSpace string
	colWidth  uint

	genericclioptions.IOStreams
}

type JavaPod struct {
	pod        corev1.Pod
	jdkVersion string
}

//New kubectl-java list sub cmd
func NewListCmd(factory *util.CmdFactory, streams genericclioptions.IOStreams) *cobra.Command {
	f := &JavaPodFinder{
		cmdFactory: factory,
		IOStreams:  streams,
		executor:   util.DefaultRemoteExecutor{},
	}

	cmd := &cobra.Command{
		Use:     listUsage,
		Short:   listShort,
		Long:    listLong,
		Example: fmt.Sprintf(listExample, "kubectl"),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			//TODO cmd should be processed by these step
			_ = f.Complete(cmd)
			_ = f.Validate()
			err = f.Run()
			return
		},
	}

	cmd.Flags().UintVarP(&f.colWidth, "colWidth", "w", 80, "colWidth used to set the table column width")

	return cmd
}

// handle user flags
func (f *JavaPodFinder) Complete(cmd *cobra.Command) error {
	namespace, err := cmd.Flags().GetString("namespace")
	if err != nil {
		return err
	}
	f.nameSpace = namespace
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
	kubConfig := f.cmdFactory.UserKubConfig
	currentContext, currentNameSpace, masterURL := util.GetCurrentConfigInfo(kubConfig)
	// namespace from user kubeconfig
	if len(f.nameSpace) == 0 {
		f.nameSpace = currentNameSpace
	}
	_, _ = fmt.Fprintf(f.Out, "context:%s\tnamespace:%s\tmaserURL:%s\n", util.Yellow(currentContext), util.Yellow(f.nameSpace), util.Yellow(masterURL))
}

// find pods that running java  application
func (f *JavaPodFinder) findJavaPods() ([]JavaPod, error) {
	coreV1Client := f.cmdFactory.ClientSet.CoreV1()
	podInfo, err := coreV1Client.Pods(corev1.NamespaceAll).List(context.TODO(), metav1.ListOptions{})

	if err != nil {
		return nil, err
	}

	pods := podInfo.Items
	restClient := coreV1Client.RESTClient()
	javaPods, podChan := make([]JavaPod, 0, len(pods)), make(chan JavaPod)

	wg.Add(len(pods))
	for _, pod := range pods {
		go f.filterJavaPod(restClient, pod, podChan)
	}

	go func() {
		wg.Wait()
		close(podChan)
	}()

	for pod := range podChan {
		javaPods = append(javaPods, pod)
	}

	return javaPods, nil
}

// filter the pod that not running java application
func (f *JavaPodFinder) filterJavaPod(client rest.Interface, pod corev1.Pod, podChan chan JavaPod) {
	defer wg.Done()
	containers := pod.Spec.Containers
	containersRunningJavaApp := make([]corev1.Container, 0, len(containers))
	var jdkVersion string
	for _, container := range containers {
		var javaCmdOutput, javaCmdErrOutput bytes.Buffer
		req := client.Post().
			Resource("pods").
			SubResource("exec").
			Namespace(pod.Namespace).
			Name(pod.Name).
			VersionedParams(&corev1.PodExecOptions{
				TypeMeta:  metav1.TypeMeta{},
				Stdin:     false,
				Stdout:    true,
				Stderr:    true,
				TTY:       false,
				Container: container.Name,
				Command:   []string{"java", "-version"},
			}, scheme.ParameterCodec)
		err := f.executor.Execute("POST", req.URL(), f.cmdFactory.ClientConfig, nil, &javaCmdOutput, &javaCmdErrOutput, false, nil)
		if err != nil {
			_ = fmt.Errorf("command exec error: %s", err)
		}
		errOutput := javaCmdErrOutput.String()
		if strings.Contains(errOutput, "jdk") {
			containersRunningJavaApp = append(containersRunningJavaApp, container)
			jdkVersion = strings.Split(errOutput, "\n")[0]
		}
	}
	if len(containersRunningJavaApp) > 0 {
		pod.Spec.Containers = containersRunningJavaApp
		podChan <- JavaPod{
			pod:        pod,
			jdkVersion: jdkVersion,
		}
	}
}

// build table for printer
func buildTableToPrint(javaPods []JavaPod) *uitable.Table {
	cols := len(headers)
	table := uitable.New()
	table.AddRow(headers...)
	rows := make([]metav1.TableRow, len(javaPods))

	for i, javaPod := range javaPods {
		pod := javaPod.pod
		podStatus, podSpec := pod.Status, pod.Spec
		containerStatuses := podStatus.ContainerStatuses
		row, containers := make([]interface{}, cols, cols), make([]string, len(containerStatuses))

		for index, status := range containerStatuses {
			containers[index] = util.HiCyan(status.Name)
		}

		// column: name
		row[0] = pod.Name
		// column: node
		row[1] = podSpec.NodeName
		// column: status
		row[2] = util.ColorizePodStatus(podStatus.Phase)
		// column: containers
		row[3] = containers
		// column: jdkVersion
		row[4] = javaPod.jdkVersion
		// fill table cells with row
		rows[i].Cells = row
		table.AddRow(row...)
	}

	return table
}

// print table to the console
func (f *JavaPodFinder) print(table *uitable.Table) error {
	table.MaxColWidth = f.colWidth
	_, err := fmt.Fprintln(f.Out, table)
	return err
}
