// Copyright (C) 2022-2023 Red Hat, Inc.
//
// This program is free software; you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation; either version 2 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along
// with this program; if not, write to the Free Software Foundation, Inc.,
// 51 Franklin Street, Fifth Floor, Boston, MA 02110-1301 USA.

package provider

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/test-network-function/cnf-certification-test/internal/log"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func TestIsIstioProxy(t *testing.T) {
	testCases := []struct {
		testContainer  Container
		expectedOutput bool
	}{
		{
			testContainer: Container{
				Container: &corev1.Container{
					Name: "istio-proxy",
				},
			},
			expectedOutput: true,
		},
		{
			testContainer: Container{
				Container: &corev1.Container{
					Name: "not-istio-proxy",
				},
			},
			expectedOutput: false,
		},
	}

	for _, tc := range testCases {
		assert.Equal(t, tc.expectedOutput, tc.testContainer.IsIstioProxy())
	}
}

func TestHasExecProbes(t *testing.T) {
	testCases := []struct {
		testContainer  Container
		expectedOutput bool
	}{
		{ // Test Case #1 - Container defines a LivenessProbe with Exec mechanism
			testContainer: Container{
				Container: &corev1.Container{
					LivenessProbe: &corev1.Probe{
						ProbeHandler: corev1.ProbeHandler{
							Exec: &corev1.ExecAction{
								Command: []string{"/bin/sh -c sleep 300"},
							},
						},
					},
				},
			},
			expectedOutput: true,
		},
		{ // Test Case #2 - Container defines a LivenessProbe with HTTP mechanism
			testContainer: Container{
				Container: &corev1.Container{
					LivenessProbe: &corev1.Probe{
						ProbeHandler: corev1.ProbeHandler{
							HTTPGet: &corev1.HTTPGetAction{
								Port: intstr.FromInt(10002),
							},
						},
					},
				},
			},
			expectedOutput: false,
		},
		{ // Test Case #3 - Container defines a LivenessProbe with HTTP and a ReadinessProbe with Exec
			testContainer: Container{
				Container: &corev1.Container{
					LivenessProbe: &corev1.Probe{
						ProbeHandler: corev1.ProbeHandler{
							HTTPGet: &corev1.HTTPGetAction{
								Port: intstr.FromInt(10002),
							},
						},
					},
					ReadinessProbe: &corev1.Probe{
						ProbeHandler: corev1.ProbeHandler{
							Exec: &corev1.ExecAction{
								Command: []string{"/bin/sh -c sleep 300"},
							},
						},
					},
				},
			},
			expectedOutput: true,
		},
		{ // Test Case #4 - Container defines a LivenessProbe and a ReadinessProbe with HTTP and a StartupProbe with Exec
			testContainer: Container{
				Container: &corev1.Container{
					LivenessProbe: &corev1.Probe{
						ProbeHandler: corev1.ProbeHandler{
							HTTPGet: &corev1.HTTPGetAction{
								Port: intstr.FromInt(10002),
							},
						},
					},
					ReadinessProbe: &corev1.Probe{
						ProbeHandler: corev1.ProbeHandler{
							HTTPGet: &corev1.HTTPGetAction{
								Port: intstr.FromInt(10005),
							},
						},
					},
					StartupProbe: &corev1.Probe{
						ProbeHandler: corev1.ProbeHandler{
							Exec: &corev1.ExecAction{
								Command: []string{"/bin/sh -c sleep 300"},
							},
						},
					},
				},
			},
			expectedOutput: true,
		},
	}

	for _, tc := range testCases {
		assert.Equal(t, tc.expectedOutput, tc.testContainer.HasExecProbes())
	}
}

func TestHasIgnoredContainerName(t *testing.T) {
	testCases := []struct {
		testContainer  Container
		expectedOutput bool
	}{
		{ // Test Case #1 - Container name is not ignored
			testContainer: Container{
				Container: &corev1.Container{
					Name: "not-ignored",
				},
			},
			expectedOutput: false,
		},
		{ // Test Case #2 - Container name is ignored
			testContainer: Container{
				Container: &corev1.Container{
					Name: "istio-proxy",
				},
			},
			expectedOutput: true,
		},
		{ // Test Case #3 - Container name contains ignored substring
			testContainer: Container{
				Container: &corev1.Container{
					Name: "istio-proxy-test",
				},
			},
			expectedOutput: true,
		},
	}

	for _, tc := range testCases {
		assert.Equal(t, tc.expectedOutput, tc.testContainer.HasIgnoredContainerName())
	}
}

func TestIsTagEmpty(t *testing.T) {
	testCases := []struct {
		testContainer  Container
		expectedOutput bool
	}{
		{ // Test Case #1 - Container image tag is empty
			testContainer: Container{
				ContainerImageIdentifier: ContainerImageIdentifier{
					Tag: "",
				},
			},
			expectedOutput: true,
		},
		{ // Test Case #2 - Container image tag is not empty
			testContainer: Container{
				ContainerImageIdentifier: ContainerImageIdentifier{
					Tag: "test",
				},
			},
			expectedOutput: false,
		},
	}

	for _, tc := range testCases {
		assert.Equal(t, tc.expectedOutput, tc.testContainer.IsTagEmpty())
	}
}

func TestIsreadOnlyRootFilessystem(t *testing.T) {
	trueVal := true
	falseVal := false
	testCases := []struct {
		testContainer  Container
		expectedOutput bool
	}{
		{
			testContainer: Container{
				Container: &corev1.Container{
					Name: "TestContainer1",
					SecurityContext: &corev1.SecurityContext{
						ReadOnlyRootFilesystem: &trueVal,
					},
				},
			},
			expectedOutput: true,
		},
		{
			testContainer: Container{
				Container: &corev1.Container{
					Name: "TestContainer2",
					SecurityContext: &corev1.SecurityContext{
						ReadOnlyRootFilesystem: &falseVal,
					},
				},
			},
			expectedOutput: false,
		},
	}

	for _, tc := range testCases {
		log.SetupLogger(os.Stdout, "INFO")
		actualOutput := tc.testContainer.IsReadOnlyRootFilesystem(log.GetLogger())
		assert.Equal(t, tc.expectedOutput, actualOutput)
	}
}

func TestIsContainerRunAsNonRoot(t *testing.T) {
	trueVal := true
	falseVal := false
	tests := []struct {
		name      string
		container Container
		expected  bool
	}{
		{
			name: "Container set to run as non-root",
			container: Container{
				Container: &corev1.Container{
					SecurityContext: &corev1.SecurityContext{
						RunAsNonRoot: &trueVal,
					},
				},
			},
			expected: true,
		},
		{
			name: "Container set to not run as non-root",
			container: Container{
				Container: &corev1.Container{
					SecurityContext: &corev1.SecurityContext{
						RunAsNonRoot: &falseVal,
					},
				},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.container.IsContainerRunAsNonRoot()
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}
