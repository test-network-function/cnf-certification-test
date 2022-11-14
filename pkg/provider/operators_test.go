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

package provider

import (
	"testing"

	"github.com/operator-framework/api/pkg/lib/version"
	olmv1Alpha "github.com/operator-framework/api/pkg/operators/v1alpha1"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestCsvToString(t *testing.T) {
	assert.Equal(t, "operator csv: test1 ns: testNS", CsvToString(&olmv1Alpha.ClusterServiceVersion{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test1",
			Namespace: "testNS",
		},
	}))
}

func TestOperatorString(t *testing.T) {
	o := Operator{
		Name:             "test1",
		Namespace:        "testNS",
		SubscriptionName: "sub1",
	}
	assert.Equal(t, "csv: test1 ns:testNS subscription:sub1", o.String())
}

//nolint:funlen
func TestCreateOperators(t *testing.T) {
	// op1 in namespace ns1
	op1Ns1 := createCsv("op1.v1.0.1", "ns1", 1, 0, 1)
	op1Ns2 := createCsv("op1.v1.0.1", "ns2", 1, 0, 1)
	op2Ns2 := createCsv("op2.v2.0.2", "ns2", 2, 0, 2)

	subscription1 := olmv1Alpha.Subscription{
		TypeMeta:   metav1.TypeMeta{Kind: "Subscription"},
		ObjectMeta: metav1.ObjectMeta{Name: "subs1", Namespace: "ns1"},
		Spec:       &olmv1Alpha.SubscriptionSpec{Package: "op1", CatalogSource: "catalogSource1"},
		Status:     olmv1Alpha.SubscriptionStatus{InstalledCSV: "op1.v1.0.1"},
	}
	subscription2 := olmv1Alpha.Subscription{
		TypeMeta:   metav1.TypeMeta{Kind: "Subscription"},
		ObjectMeta: metav1.ObjectMeta{Name: "subs2", Namespace: "ns2"},
		Spec:       &olmv1Alpha.SubscriptionSpec{Package: "op1", CatalogSource: "catalogSource2"},
		Status:     olmv1Alpha.SubscriptionStatus{InstalledCSV: "op1.v1.0.1"},
	}
	subscription3 := olmv1Alpha.Subscription{
		TypeMeta:   metav1.TypeMeta{Kind: "Subscription"},
		ObjectMeta: metav1.ObjectMeta{Name: "subs3", Namespace: "ns2"},
		Spec:       &olmv1Alpha.SubscriptionSpec{Package: "op2", CatalogSource: "catalogSource3"},
		Status:     olmv1Alpha.SubscriptionStatus{InstalledCSV: "op2.v2.0.2"},
	}

	testCases := []struct {
		csvs              []olmv1Alpha.ClusterServiceVersion
		subscriptions     []olmv1Alpha.Subscription
		installPlan       []*olmv1Alpha.InstallPlan
		catalogSource     []*olmv1Alpha.CatalogSource
		expectedOperators []*Operator
		expectedErrorStr  string
	}{
		// ns1: csv1/subs1
		{
			csvs:              []olmv1Alpha.ClusterServiceVersion{},
			subscriptions:     []olmv1Alpha.Subscription{subscription1},
			installPlan:       []*olmv1Alpha.InstallPlan{&ns1InstallPlan1},
			catalogSource:     []*olmv1Alpha.CatalogSource{&catalogSource1},
			expectedOperators: []*Operator{},
		},
		// ns1: csv1/subs1
		{
			csvs:          []olmv1Alpha.ClusterServiceVersion{op1Ns1},
			subscriptions: []olmv1Alpha.Subscription{subscription1},
			installPlan:   []*olmv1Alpha.InstallPlan{&ns1InstallPlan1},
			catalogSource: []*olmv1Alpha.CatalogSource{&catalogSource1},
			expectedOperators: []*Operator{
				{
					Name:             "op1.v1.0.1",
					Namespace:        "ns1",
					Csv:              &op1Ns1,
					SubscriptionName: "subs1",
					InstallPlans: []CsvInstallPlan{
						{
							Name:        "ns1Plan1",
							BundleImage: "lookuppath1",
							IndexImage:  "catalogSource1Image",
						},
					},
					Package:            "op1",
					Org:                "catalogSource1",
					Version:            "1.0.1",
					PackageFromCsvName: "op1",
				},
			},
		},
		// ns1: csv1/subs1 - installPlan not found.
		{
			csvs:          []olmv1Alpha.ClusterServiceVersion{op1Ns1},
			subscriptions: []olmv1Alpha.Subscription{subscription1},
			installPlan:   []*olmv1Alpha.InstallPlan{},
			catalogSource: []*olmv1Alpha.CatalogSource{},
			expectedOperators: []*Operator{
				{
					Name:               "op1.v1.0.1",
					Namespace:          "ns1",
					Csv:                &op1Ns1,
					SubscriptionName:   "subs1",
					InstallPlans:       nil,
					Package:            "op1",
					Org:                "catalogSource1",
					Version:            "1.0.1",
					PackageFromCsvName: "op1",
				},
			},
		},
		// ns1: csv1/subs1 - bundleImage not found.
		{
			csvs:          []olmv1Alpha.ClusterServiceVersion{op1Ns1},
			subscriptions: []olmv1Alpha.Subscription{subscription1},
			installPlan:   []*olmv1Alpha.InstallPlan{&ns1InstallPlan1},
			catalogSource: []*olmv1Alpha.CatalogSource{},
			expectedOperators: []*Operator{
				{
					Name:               "op1.v1.0.1",
					Namespace:          "ns1",
					Csv:                &op1Ns1,
					SubscriptionName:   "subs1",
					InstallPlans:       nil,
					Package:            "op1",
					Org:                "catalogSource1",
					Version:            "1.0.1",
					PackageFromCsvName: "op1",
				},
			},
		},
		// ns1: csv1/subs1, ns2: csv2 (without subscription)
		{
			csvs:          []olmv1Alpha.ClusterServiceVersion{op1Ns1, op1Ns2},
			subscriptions: []olmv1Alpha.Subscription{subscription1},
			installPlan:   []*olmv1Alpha.InstallPlan{&ns1InstallPlan1, &ns2InstallPlan1},
			catalogSource: []*olmv1Alpha.CatalogSource{&catalogSource1, &catalogSource2},
			expectedOperators: []*Operator{
				{
					Name:             "op1.v1.0.1",
					Namespace:        "ns1",
					Csv:              &op1Ns1,
					SubscriptionName: "subs1",
					InstallPlans: []CsvInstallPlan{
						{
							Name:        "ns1Plan1",
							BundleImage: "lookuppath1",
							IndexImage:  "catalogSource1Image",
						},
					},
					Package:            "op1",
					Org:                "catalogSource1",
					Version:            "1.0.1",
					PackageFromCsvName: "op1",
				},
				{
					Name:             "op1.v1.0.1",
					Namespace:        "ns2",
					Csv:              &op1Ns2,
					SubscriptionName: "",
					InstallPlans: []CsvInstallPlan{
						{
							Name:        "ns2Plan1",
							BundleImage: "lookuppath2",
							IndexImage:  "catalogSource2Image",
						},
					},
					Package:            "",
					Version:            "1.0.1",
					PackageFromCsvName: "op1",
				},
			},
		},
		// ns1: csv1/subs1, ns2: csv2/subs2
		{
			csvs:          []olmv1Alpha.ClusterServiceVersion{op1Ns1, op1Ns2},
			subscriptions: []olmv1Alpha.Subscription{subscription1, subscription2},
			installPlan:   []*olmv1Alpha.InstallPlan{&ns1InstallPlan1, &ns2InstallPlan1},
			catalogSource: []*olmv1Alpha.CatalogSource{&catalogSource1, &catalogSource2},
			expectedOperators: []*Operator{
				{
					Name:             "op1.v1.0.1",
					Namespace:        "ns1",
					Csv:              &op1Ns1,
					SubscriptionName: "subs1",
					InstallPlans: []CsvInstallPlan{
						{
							Name:        "ns1Plan1",
							BundleImage: "lookuppath1",
							IndexImage:  "catalogSource1Image",
						},
					},
					Package:            "op1",
					Org:                "catalogSource1",
					Version:            "1.0.1",
					PackageFromCsvName: "op1",
				},
				{
					Name:             "op1.v1.0.1",
					Namespace:        "ns2",
					Csv:              &op1Ns2,
					SubscriptionName: "subs2",
					InstallPlans: []CsvInstallPlan{
						{
							Name:        "ns2Plan1",
							BundleImage: "lookuppath2",
							IndexImage:  "catalogSource2Image",
						},
					},
					Package:            "op1",
					Org:                "catalogSource2",
					Version:            "1.0.1",
					PackageFromCsvName: "op1",
				},
			},
		},
		// ns1: csv1/subs1, ns2: csv2/subs2 + csv3/subs3
		{
			csvs:          []olmv1Alpha.ClusterServiceVersion{op1Ns1, op1Ns2, op2Ns2},
			subscriptions: []olmv1Alpha.Subscription{subscription1, subscription2, subscription3},
			installPlan:   []*olmv1Alpha.InstallPlan{&ns1InstallPlan1, &ns2InstallPlan1, &ns2InstallPlan2},
			catalogSource: []*olmv1Alpha.CatalogSource{&catalogSource1, &catalogSource2, &catalogSource3},
			expectedOperators: []*Operator{
				{
					Name:             "op1.v1.0.1",
					Namespace:        "ns1",
					Csv:              &op1Ns1,
					SubscriptionName: "subs1",
					InstallPlans: []CsvInstallPlan{
						{
							Name:        "ns1Plan1",
							BundleImage: "lookuppath1",
							IndexImage:  "catalogSource1Image",
						},
					},
					Package:            "op1",
					Org:                "catalogSource1",
					Version:            "1.0.1",
					PackageFromCsvName: "op1",
				},
				{
					Name:             "op1.v1.0.1",
					Namespace:        "ns2",
					Csv:              &op1Ns2,
					SubscriptionName: "subs2",
					InstallPlans: []CsvInstallPlan{
						{
							Name:        "ns2Plan1",
							BundleImage: "lookuppath2",
							IndexImage:  "catalogSource2Image",
						},
					},
					Package:            "op1",
					Org:                "catalogSource2",
					Version:            "1.0.1",
					PackageFromCsvName: "op1",
				},

				{
					Name:             "op2.v2.0.2",
					Namespace:        "ns2",
					Csv:              &op2Ns2,
					SubscriptionName: "subs3",
					InstallPlans: []CsvInstallPlan{
						{
							Name:        "ns2Plan2",
							BundleImage: "lookuppath3",
							IndexImage:  "catalogSource3Image",
						},
					},
					Package:            "op2",
					Org:                "catalogSource3",
					Version:            "2.0.2",
					PackageFromCsvName: "op2",
				},
			},
		},
	}

	for _, tc := range testCases {
		ops := createOperators(tc.csvs, tc.subscriptions, tc.installPlan, tc.catalogSource, true, false)
		assert.Equal(t, tc.expectedOperators, ops)
	}
}

func createCsv(name, namespace string, verMajor, verMinor, verPatch uint64) (aCsv olmv1Alpha.ClusterServiceVersion) {
	aCsv.Name = name
	aCsv.Namespace = namespace

	aVersion := version.OperatorVersion{}
	aVersion.Major = verMajor
	aVersion.Minor = verMinor
	aVersion.Patch = verPatch

	aCsv.Spec.Version = aVersion
	return aCsv
}

func Test_getCatalogSourceImageIndexFromInstallPlan(t *testing.T) {
	type args struct {
		installPlan       *olmv1Alpha.InstallPlan
		allCatalogSources []*olmv1Alpha.CatalogSource
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "ok",
			args:    args{installPlan: &ns1InstallPlan1, allCatalogSources: []*olmv1Alpha.CatalogSource{&catalogSource1}},
			want:    "catalogSource1Image",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getCatalogSourceImageIndexFromInstallPlan(tt.args.installPlan, tt.args.allCatalogSources)
			if (err != nil) != tt.wantErr {
				t.Errorf("getCatalogSourceImageIndexFromInstallPlan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getCatalogSourceImageIndexFromInstallPlan() = %v, want %v", got, tt.want)
			}
		})
	}
}
