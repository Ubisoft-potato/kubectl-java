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
	"github.com/cyka/kubectl-java/util"
	"github.com/jedib0t/go-pretty/table"
	"github.com/jedib0t/go-pretty/text"
	"github.com/spf13/cobra"

	"k8s.io/cli-runtime/pkg/genericclioptions"
)

var (
	// command info
	listUsage   = "list [flags]"
	listShort   = "List All Pods That Running Java Application"
	listLong    = "List All Pods That Running Java Application"
	listExample = ""

	// table header
	header = table.Row{"node", "pod"}
)

//New kubectl-java list sub cmd
func NewListCmd(finder *JavaPodFinder) *cobra.Command {
	cmd := &cobra.Command{
		Use:     listUsage,
		Short:   listShort,
		Long:    listLong,
		Example: listExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			finder.addKubeConfigInfo()
			finder.render()
			return nil
		},
	}
	return cmd
}

type JavaPodFinder struct {
	genericclioptions.IOStreams
	options *KubeJavaAppOptions
	writer  table.Writer
}

func NewJavaPodFinder(IOStreams genericclioptions.IOStreams, options *KubeJavaAppOptions) *JavaPodFinder {
	writer := table.NewWriter()
	writer.SetOutputMirror(IOStreams.Out)
	customTableWriter(writer)
	return &JavaPodFinder{
		IOStreams: IOStreams,
		options:   options,
		writer:    writer,
	}
}

func (f *JavaPodFinder) addKubeConfigInfo() {
	kubConfig := f.options.userKubConfig
	currentContext, currentNameSpace, masterURL := util.GetCurrentConfigInfo(kubConfig)
	f.writer.SetTitle("Context:%s NameSpace:%s MasterURL:%s", currentContext, currentNameSpace, masterURL)
	f.writer.AppendRow(table.Row{"1.1.1.1", "order-service"})
}

func (f *JavaPodFinder) render() {
	f.writer.Render()
}

func customTableWriter(writer table.Writer) {
	writer.AppendHeader(header)
	writer.SetStyle(table.StyleColoredBright)
	writer.Style().Title.Align = text.AlignCenter
	writer.Style().Title.Colors = text.Colors{text.BgHiBlack, text.FgHiYellow}
	writer.SetColumnConfigs([]table.ColumnConfig{
		{
			Number:      1,
			Align:       text.AlignCenter,
			AlignHeader: text.AlignCenter,
			WidthMin:    30,
			WidthMax:    40,
		},
		{
			Number:      2,
			Align:       text.AlignCenter,
			AlignFooter: 0,
			AlignHeader: text.AlignCenter,
			WidthMin:    30,
			WidthMax:    40,
		},
	})
}
