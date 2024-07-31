// Copyright (C) 2020-2024 Red Hat, Inc.
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

package observability

import (
	"testing"

	"github.com/redhat-best-practices-for-k8s/certsuite/pkg/provider"
	"github.com/redhat-best-practices-for-k8s/certsuite/pkg/testhelper"
	"github.com/stretchr/testify/assert"
)

func TestBuildToBeRemovedWorkloadAPIList(t *testing.T) {
	testCases := []struct {
		name             string
		apiRequestCounts []provider.APIRequestCount
		expected         []provider.ToBeRemovedWorkloadAPI
	}{
		{
			name: "Test to ensure proper conversion of APIRequestCount to ToBeRemovedWorkloadAPI",
			apiRequestCounts: []provider.APIRequestCount{
				{
					Metadata: struct {
						Name string `json:"name"`
					}{Name: "api1"},
					Status: struct {
						RemovedInRelease string `json:"removedInRelease"`
						Last24h          []struct {
							ByNode []struct {
								ByUser []struct {
									UserName  string `json:"username"`
									UserAgent string `json:"userAgent"`
								} `json:"byUser"`
							} `json:"byNode"`
						} `json:"last24h"`
					}{
						RemovedInRelease: "v1.20",
						Last24h: []struct {
							ByNode []struct {
								ByUser []struct {
									UserName  string `json:"username"`
									UserAgent string `json:"userAgent"`
								} `json:"byUser"`
							} `json:"byNode"`
						}{
							{
								ByNode: []struct {
									ByUser []struct {
										UserName  string `json:"username"`
										UserAgent string `json:"userAgent"`
									} `json:"byUser"`
								}{
									{
										ByUser: []struct {
											UserName  string `json:"username"`
											UserAgent string `json:"userAgent"`
										}{
											{UserName: "user1", UserAgent: "agent1"},
										},
									},
								},
							},
						},
					},
				},
			},
			expected: []provider.ToBeRemovedWorkloadAPI{
				{
					APIName:          "api1",
					RemovedInRelease: "v1.20",
					UserInfo:         map[string]struct{}{"UserName: user1, UserAgent: agent1": {}},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := buildToBeRemovedWorkloadAPIList(tc.apiRequestCounts)
			assert.ElementsMatch(t, tc.expected, result)
		})
	}
}

func TestEvaluateAPICompliance(t *testing.T) {
	// Mock toBeRemovedWorkloadAPIs
	toBeRemovedWorkloadAPIs := []provider.ToBeRemovedWorkloadAPI{
		{
			APIName:          "api1",
			RemovedInRelease: "1.21",
			UserInfo: map[string]struct{}{
				"UserName: user1, UserAgent: agent1": {},
			},
		},
		{
			APIName:          "api2",
			RemovedInRelease: "1.20",
			UserInfo: map[string]struct{}{
				"UserName: user2, UserAgent: agent2": {},
			},
		},
	}

	// Mock current Kubernetes version
	kubernetesVersion := "1.19"

	// Call the function
	compliantObjects, nonCompliantObjects := evaluateAPICompliance(toBeRemovedWorkloadAPIs, kubernetesVersion)

	// Helper to create a ReportObject with fields
	createReportObject := func(apiName, removedInRelease string, isCompliant bool, userInfo []string) *testhelper.ReportObject {
		obj := testhelper.NewReportObject("API", "StubType", isCompliant)
		obj.AddField("APIName", apiName)
		obj.AddField("RemovedInRelease", removedInRelease)
		for _, info := range userInfo {
			obj.AddField("UserInfo", info)
		}
		return obj
	}

	// Expected results
	expectedCompliantObjects := []*testhelper.ReportObject{
		// API removed in 1.21 is compliant with Kubernetes 1.20 = (current version 1.19  + 1)
		createReportObject("api1", "1.21", true, []string{"UserName: user1, UserAgent: agent1"}),
	}
	expectedNonCompliantObjects := []*testhelper.ReportObject{
		// API removed in 1.20 is non-compliant with Kubernetes 1.20 = (current version 1.19  + 1)
		createReportObject("api2", "1.20", false, []string{"UserName: user2, UserAgent: agent2"}),
	}

	// Verify the results
	assert.Len(t, compliantObjects, len(expectedCompliantObjects))
	for i, obj := range compliantObjects {
		assert.Equal(t, expectedCompliantObjects[i].ObjectType, obj.ObjectType)
		assert.ElementsMatch(t, expectedCompliantObjects[i].ObjectFieldsKeys, obj.ObjectFieldsKeys)
		assert.ElementsMatch(t, expectedCompliantObjects[i].ObjectFieldsValues, obj.ObjectFieldsValues)
	}

	assert.Len(t, nonCompliantObjects, len(expectedNonCompliantObjects))
	for i, obj := range nonCompliantObjects {
		assert.Equal(t, expectedNonCompliantObjects[i].ObjectType, obj.ObjectType)
		assert.ElementsMatch(t, expectedNonCompliantObjects[i].ObjectFieldsKeys, obj.ObjectFieldsKeys)
		assert.ElementsMatch(t, expectedNonCompliantObjects[i].ObjectFieldsValues, obj.ObjectFieldsValues)
	}
}
