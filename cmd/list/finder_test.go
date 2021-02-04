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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"reflect"
	"testing"

	"github.com/cyka/kubectl-java/util"
	"github.com/gosuri/uitable"
)

func Test_buildTableToPrint(t *testing.T) {
	type args struct {
		javaPods []JavaPod
	}

	buildTable := uitable.New()
	buildTable.AddRow(headers...)
	buildTable.AddRow([]interface{}{"podName", "node-01", util.Cyan("Running"), []string{util.HiCyan("nginx-web")}, "1.8.0"}...)

	tests := []struct {
		name string
		args args
		want *uitable.Table
	}{
		{
			name: "buildTable",
			args: args{
				javaPods: []JavaPod{
					{
						pod: corev1.Pod{
							ObjectMeta: metav1.ObjectMeta{
								Name: "podName",
							},
							Spec: corev1.PodSpec{
								NodeName: "node-01",
							},
							Status: corev1.PodStatus{
								Phase: corev1.PodRunning,
								ContainerStatuses: []corev1.ContainerStatus{
									{
										Name: "nginx-web",
									},
								},
							},
						},
						jdkVersion: "1.8.0",
					},
				},
			},
			want: buildTable,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildTableToPrint(tt.args.javaPods); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("buildTableToPrint() = %v, want %v", got, tt.want)
			}
		})
	}
}
