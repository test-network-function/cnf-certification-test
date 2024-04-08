// Copyright (C) 2022-2024 Red Hat, Inc.
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
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/test-network-function/cnf-certification-test/internal/clientsholder"
	"github.com/test-network-function/cnf-certification-test/internal/log"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	v1 "k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	k8sfake "k8s.io/client-go/kubernetes/fake"
)

func TestPod_CheckResourceOnly2MiHugePages(t *testing.T) {
	tests := []struct {
		name string
		aPod Pod
		want bool
	}{
		{
			name: "pass",
			aPod: *generatePod(10, 10, 0, 0),
			want: true,
		},
		{
			name: "fail",
			aPod: *generatePod(10, 10, 1, 1),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := tt.aPod
			got := p.CheckResourceHugePagesSize(HugePages2Mi)
			if got != tt.want {
				t.Errorf("Pod.CheckResourceHugePagesSize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPod_CheckResourceOnly1GiHugePages(t *testing.T) {
	tests := []struct {
		name string
		aPod Pod
		want bool
	}{
		{
			name: "pass",
			aPod: *generatePod(0, 0, 1, 1),
			want: true,
		},
		{
			name: "fail",
			aPod: *generatePod(10, 10, 1, 1),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := tt.aPod
			got := p.CheckResourceHugePagesSize(HugePages1Gi)
			if got != tt.want {
				t.Errorf("Pod.CheckResourceHugePagesSize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func generatePod(requestsValue2M, limitsValue2M, requestsValue1G, limitsValue1G int64) *Pod {
	aPod := Pod{
		Containers: []*Container{
			{
				Container: &corev1.Container{
					Name: "test1",
					Resources: corev1.ResourceRequirements{
						Requests: corev1.ResourceList{},
						Limits:   corev1.ResourceList{}}},
			},
		},
	}
	var aQuantity v1.Quantity
	if requestsValue2M != 0 {
		aQuantity.Set(requestsValue2M)
		aPod.Containers[0].Resources.Requests[HugePages2Mi] = aQuantity
	}
	if limitsValue2M != 0 {
		aQuantity.Set(limitsValue2M)
		aPod.Containers[0].Resources.Limits[HugePages2Mi] = aQuantity
	}

	if requestsValue1G != 0 {
		aQuantity.Set(requestsValue1G)
		aPod.Containers[0].Resources.Requests[HugePages1Gi] = aQuantity
	}
	if limitsValue1G != 0 {
		aQuantity.Set(limitsValue1G)
		aPod.Containers[0].Resources.Limits[HugePages1Gi] = aQuantity
	}
	return &aPod
}

func TestIsAffinityCompliantPods(t *testing.T) {
	testCases := []struct {
		testPod      Pod
		resultErrStr error
		isCompliant  bool
	}{
		{ // Test Case #1 - Affinity is nil, AffinityRequired label is set, fail
			testPod: Pod{
				Pod: &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{
							AffinityRequiredKey: "true",
						},
					},
					Spec: corev1.PodSpec{
						Affinity: nil,
					},
				},
			},
			resultErrStr: errors.New("has been found with an AffinityRequired flag but is missing corresponding affinity rules"),
			isCompliant:  false,
		},
		{ // Test Case #2 - Affinity is not nil, but PodAffinity/NodeAffinity are also not set, fail
			testPod: Pod{
				Pod: &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{
							AffinityRequiredKey: "true",
						},
					},
					Spec: corev1.PodSpec{
						Affinity: &corev1.Affinity{}, // not nil
					},
				},
			},
			resultErrStr: errors.New("has been found with an AffinityRequired flag but is missing corresponding pod/node affinity rules"),
			isCompliant:  false,
		},
		{ // Test Case #3 - Affinity is not nil, but anti-affinity rule is set which defeats the purpose of the Required flag
			testPod: Pod{
				Pod: &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{
							AffinityRequiredKey: "true",
						},
					},
					Spec: corev1.PodSpec{
						Affinity: &corev1.Affinity{
							PodAntiAffinity: &corev1.PodAntiAffinity{},
						},
					},
				},
			},
			resultErrStr: errors.New("has been found with an AffinityRequired flag but has anti-affinity rules"),
			isCompliant:  false,
		},
	}

	for _, tc := range testCases {
		result, testErr := tc.testPod.IsAffinityCompliant()
		assert.Contains(t, testErr.Error(), tc.resultErrStr.Error())
		assert.Equal(t, tc.isCompliant, result)
	}
}

func TestPodString(t *testing.T) {
	p := Pod{
		Pod: &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test1",
				Namespace: "testNS",
			},
		},
	}
	assert.Equal(t, "pod: test1 ns: testNS", p.String())
}

func TestContainsIstioProxy(t *testing.T) {
	testCases := []struct {
		testPod        Pod
		expectedOutput bool
	}{
		{
			testPod: Pod{
				Containers: []*Container{
					{
						Container: &corev1.Container{
							Name: "istio-proxy",
						},
					},
				},
			},
			expectedOutput: true,
		},
		{
			testPod: Pod{
				Containers: []*Container{
					{
						Container: &corev1.Container{
							Name: "not-istio-proxy",
						},
					},
				},
			},
			expectedOutput: false,
		},
	}

	for _, tc := range testCases {
		assert.Equal(t, tc.expectedOutput, tc.testPod.ContainsIstioProxy())
	}
}

func TestHasNodeSelector(t *testing.T) {
	testCases := []struct {
		testPod        Pod
		expectedOutput bool
	}{
		{ // Test #1 - Has a nil node selector and no NodeName, pass
			testPod: Pod{
				Pod: &corev1.Pod{
					Spec: corev1.PodSpec{
						NodeSelector: nil,
						NodeName:     "",
					},
				},
			},
			expectedOutput: false,
		},
		{ // Test #2 - Has a nodeSelector, fail
			testPod: Pod{
				Pod: &corev1.Pod{
					Spec: corev1.PodSpec{
						NodeSelector: map[string]string{
							"test1": "value1",
						},
					},
				},
			},
			expectedOutput: true,
		},
		{ // Test #3 - Has a nodeSelector initialized but it is empty, pass
			testPod: Pod{
				Pod: &corev1.Pod{
					Spec: corev1.PodSpec{
						NodeSelector: map[string]string{},
					},
				},
			},
			expectedOutput: false,
		},
	}

	for _, tc := range testCases {
		assert.Equal(t, tc.expectedOutput, tc.testPod.HasNodeSelector())
	}
}

func TestGetSRIOVNetworksNamesFromCNCFNetworks(t *testing.T) {
	testCases := []struct {
		networksAnnotation   string
		expectedNetworkNames []string
	}{
		{
			networksAnnotation:   "",
			expectedNetworkNames: []string{},
		},
		{
			networksAnnotation:   "net1,net2",
			expectedNetworkNames: []string{"net1", "net2"},
		},
		{
			networksAnnotation:   "   net1,       net2  ,net3",
			expectedNetworkNames: []string{"net1", "net2", "net3"},
		},
		{
			networksAnnotation:   `[{"name": "net1", "otherField" : "otherFieldValue1"}, {"name": "net2"}]`,
			expectedNetworkNames: []string{"net1", "net2"},
		},
	}

	for _, tc := range testCases {
		netNames := getCNCFNetworksNamesFromPodAnnotation(tc.networksAnnotation)
		assert.Equal(t, tc.expectedNetworkNames, netNames)
	}
}

func TestIssNetworkAttachmentDefinitionConfigTypeSRIOV(t *testing.T) {
	testCases := []struct {
		networkAttachmentDefinition string
		expectedNadTypeSriov        bool
		expectedErrorMsg            string
	}{
		// Single plugin mode:
		{
			networkAttachmentDefinition: "",
			expectedErrorMsg:            "failed to unmarshal cni config : unexpected end of JSON input",
		},
		{
			networkAttachmentDefinition: `{"cniVersion" : "0.4.0", "type" : "macvlan", "otherField": "true"}`,
			expectedNadTypeSriov:        false,
		},
		{
			networkAttachmentDefinition: `{"cniVersion" : "0.4.0", "type" : "sriov", "otherField": "true"}`,
			expectedNadTypeSriov:        true,
		},
		// Multi-plugin mode:
		{
			networkAttachmentDefinition: `{"cniVersion" : "0.4.0", "plugins" : [{"type": "mcvlan", "otherField": "true"}, {"type": "firewall"}]}`,
			expectedNadTypeSriov:        false,
		},
		{
			networkAttachmentDefinition: `{"cniVersion" : "0.4.0", "plugins" : [{"type": "mcvlan", "otherField": "true"}, {"type": "sriov", "otherfield": "false"}]}`,
			expectedNadTypeSriov:        true,
		},
	}

	for _, tc := range testCases {
		isTypeSriov, err := isNetworkAttachmentDefinitionConfigTypeSRIOV(tc.networkAttachmentDefinition)
		if err != nil {
			assert.Equal(t, tc.expectedErrorMsg, err.Error())
		} else {
			assert.Equal(t, tc.expectedNadTypeSriov, isTypeSriov)
		}
	}
}

func TestIsUsingClusterRoleBinding(t *testing.T) {
	testCases := []struct {
		testPod                 Pod
		testClusterRoleBindings []rbacv1.ClusterRoleBinding
		testServiceAccounts     []corev1.ServiceAccount
		testResult              bool
		testroleRefName         string
		testErr                 error
	}{
		{ // Test Case #1 - Empty ServiceAccountName, return false
			testPod: Pod{
				Pod: &corev1.Pod{
					Spec: corev1.PodSpec{
						ServiceAccountName: "",
					},
				},
			},
			testClusterRoleBindings: []rbacv1.ClusterRoleBinding{},
			testServiceAccounts:     []corev1.ServiceAccount{},
			testResult:              false,
			testErr:                 nil,
		},
		{ // Test Case #2 - ServiceAccountName set, but no ClusterRoleBinding found, return false
			testPod: Pod{
				Pod: &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-pod",
						Namespace: "test-namespace",
					},
					Spec: corev1.PodSpec{
						ServiceAccountName: "test-service-account",
					},
				},
			},
			testClusterRoleBindings: []rbacv1.ClusterRoleBinding{},
			testServiceAccounts: []corev1.ServiceAccount{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-service-account",
						Namespace: "test-namespace",
					},
				},
			},
			testResult: false,
			testErr:    nil,
		},
		{ // Test Case #3 - ServiceAccountName set, ClusterRoleBinding found, return true
			testPod: Pod{
				Pod: &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-pod",
						Namespace: "test-namespace",
					},
					Spec: corev1.PodSpec{
						ServiceAccountName: "test-service-account",
					},
				},
			},
			testClusterRoleBindings: []rbacv1.ClusterRoleBinding{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "test-cluster-role-binding",
					},
					Subjects: []rbacv1.Subject{
						{
							Kind:      rbacv1.ServiceAccountKind,
							Name:      "test-service-account",
							Namespace: "test-namespace",
						},
					},
					RoleRef: rbacv1.RoleRef{
						Name: "test",
					},
				},
			},
			testServiceAccounts: []corev1.ServiceAccount{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-service-account",
						Namespace: "test-namespace",
					},
				},
			},
			testResult:      true,
			testroleRefName: "test",
			testErr:         nil,
		},
	}

	for _, tc := range testCases {
		// Create runtimeObjects based on the ServiceAccounts and ClusterRoleBindings
		// to pass to the clientsholder for testing
		var testRuntimeObjects []runtime.Object
		for _, sa := range tc.testServiceAccounts {
			saTemp := sa
			testRuntimeObjects = append(testRuntimeObjects, &saTemp)
		}

		for _, crb := range tc.testClusterRoleBindings {
			crbTemp := crb
			testRuntimeObjects = append(testRuntimeObjects, &crbTemp)
		}

		c := k8sfake.NewSimpleClientset(testRuntimeObjects...)
		clientsholder.SetTestK8sClientsHolder(c)
		var logArchive strings.Builder
		log.SetupLogger(&logArchive, "INFO")
		result, roleRefName, err := tc.testPod.IsUsingClusterRoleBinding(tc.testClusterRoleBindings, log.GetLogger())
		assert.Equal(t, tc.testResult, result)
		assert.Equal(t, tc.testroleRefName, roleRefName)
		assert.Equal(t, tc.testErr, err)
	}
}

func TestIsRunAsUserID(t *testing.T) {
	testCases := []struct {
		testPod        Pod
		testUID        int64
		expectedOutput bool
	}{
		{ // Test Case #1 - Empty SecurityContext, return false
			testPod: Pod{
				Pod: &corev1.Pod{
					Spec: corev1.PodSpec{
						SecurityContext: &corev1.PodSecurityContext{},
					},
				},
			},
			testUID:        1337,
			expectedOutput: false,
		},
		{ // Test Case #2 - SecurityContext.RunAsUser set to 1337, return true
			testPod: Pod{
				Pod: &corev1.Pod{
					Spec: corev1.PodSpec{
						SecurityContext: &corev1.PodSecurityContext{
							RunAsUser: func() *int64 {
								var uid int64 = 1337
								return &uid
							}(),
						},
					},
				},
			},
			testUID:        1337,
			expectedOutput: true,
		},
		{ // Test Case #3 - SecurityContext.RunAsUser set to 1336, return false
			testPod: Pod{
				Pod: &corev1.Pod{
					Spec: corev1.PodSpec{
						SecurityContext: &corev1.PodSecurityContext{
							RunAsUser: func() *int64 {
								var uid int64 = 1336
								return &uid
							}(),
						},
					},
				},
			},
			testUID:        1337,
			expectedOutput: false,
		},
	}

	for _, tc := range testCases {
		assert.Equal(t, tc.expectedOutput, tc.testPod.IsRunAsUserID(tc.testUID))
	}
}

func TestIsRunAsNonRoot(t *testing.T) {
	tests := []struct {
		name     string
		pod      *Pod
		expected bool
	}{
		{
			name: "All containers and pod set to run as non-root",
			pod: &Pod{
				Pod: &corev1.Pod{
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								SecurityContext: &corev1.SecurityContext{
									RunAsNonRoot: boolPtr(true),
								},
							},
						},
						SecurityContext: &corev1.PodSecurityContext{
							RunAsNonRoot: boolPtr(true),
						},
					},
				},
			},
			expected: true,
		},
		{
			name: "No containers, pod set to run as non-root",
			pod: &Pod{
				Pod: &corev1.Pod{
					Spec: corev1.PodSpec{
						SecurityContext: &corev1.PodSecurityContext{
							RunAsNonRoot: boolPtr(true),
						},
					},
				},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.pod.IsRunAsNonRoot()
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// Helper function to get a pointer to a bool value
func boolPtr(b bool) *bool {
	return &b
}
