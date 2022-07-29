// Copyright (C) 2022 Red Hat, Inc.
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

package autodiscover

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/test-network-function/cnf-certification-test/internal/clientsholder"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestGetLimitRanges(t *testing.T) {
	generateLimitRange := func(name string) *corev1.LimitRange {
		return &corev1.LimitRange{
			ObjectMeta: metav1.ObjectMeta{
				Name: name,
				// Namespace: "testNamespace",
			},
			Spec: corev1.LimitRangeSpec{
				Limits: []corev1.LimitRangeItem{
					{
						Type: corev1.LimitTypeContainer,
					},
				},
			},
		}
	}

	testCases := []struct {
		lrName      string
		expectedLRs []corev1.LimitRange
	}{
		{
			lrName: "test1",
			expectedLRs: []corev1.LimitRange{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "test1",
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		var testRuntimeObjects []runtime.Object
		testRuntimeObjects = append(testRuntimeObjects, generateLimitRange(tc.lrName))
		oc := clientsholder.GetTestClientsHolder(testRuntimeObjects)
		limitRanges, err := getLimitRanges(oc.K8sClient.CoreV1())
		assert.Nil(t, err)
		assert.Equal(t, tc.expectedLRs[0].Name, limitRanges[0].Name)
	}
}

func TestGetResourceQuotas(t *testing.T) {
	generateResourceQuota := func(name string) *corev1.ResourceQuota {
		return &corev1.ResourceQuota{
			ObjectMeta: metav1.ObjectMeta{
				Name: name,
				// Namespace: "testNamespace",
			},
			Spec: corev1.ResourceQuotaSpec{},
		}
	}

	testCases := []struct {
		rqName      string
		expectedRQs []corev1.ResourceQuota
	}{
		{
			rqName: "test1",
			expectedRQs: []corev1.ResourceQuota{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "test1",
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		var testRuntimeObjects []runtime.Object
		testRuntimeObjects = append(testRuntimeObjects, generateResourceQuota(tc.rqName))
		oc := clientsholder.GetTestClientsHolder(testRuntimeObjects)
		resourceQuotas, err := getResourceQuotas(oc.K8sClient.CoreV1())
		assert.Nil(t, err)
		assert.Equal(t, tc.expectedRQs[0].Name, resourceQuotas[0].Name)
	}
}
