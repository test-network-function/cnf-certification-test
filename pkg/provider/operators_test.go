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
	"reflect"
	"testing"

	"github.com/operator-framework/api/pkg/lib/version"
	olmv1 "github.com/operator-framework/api/pkg/operators/v1"
	olmv1Alpha "github.com/operator-framework/api/pkg/operators/v1alpha1"
	"github.com/stretchr/testify/assert"
	"github.com/test-network-function/cnf-certification-test/internal/clientsholder"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
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
		TargetNamespaces: []string{},
	}
	assert.Equal(t, "csv: test1 ns:testNS subscription:sub1 targetNamespaces=[]", o.String())
}

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

	operatorGroupNs1 := olmv1.OperatorGroup{
		ObjectMeta: metav1.ObjectMeta{Name: "op1", Namespace: "ns1"},
		Spec:       olmv1.OperatorGroupSpec{TargetNamespaces: []string{"ns1"}},
	}
	operatorGroupNs2 := olmv1.OperatorGroup{
		ObjectMeta: metav1.ObjectMeta{Name: "op1", Namespace: "ns2"},
		Spec:       olmv1.OperatorGroupSpec{TargetNamespaces: []string{"ns2"}},
	}

	testCases := []struct {
		csvs              []*olmv1Alpha.ClusterServiceVersion
		subscriptions     []olmv1Alpha.Subscription
		installPlan       []*olmv1Alpha.InstallPlan
		catalogSource     []*olmv1Alpha.CatalogSource
		operatorGroups    []*olmv1.OperatorGroup
		expectedOperators []*Operator
		expectedErrorStr  string
	}{
		// ns1: csv1/subs1
		{
			csvs:              []*olmv1Alpha.ClusterServiceVersion{},
			subscriptions:     []olmv1Alpha.Subscription{subscription1},
			installPlan:       []*olmv1Alpha.InstallPlan{&ns1InstallPlan1},
			catalogSource:     []*olmv1Alpha.CatalogSource{&catalogSource1},
			operatorGroups:    []*olmv1.OperatorGroup{&operatorGroupNs1},
			expectedOperators: []*Operator{},
		},
		// ns1: csv1/subs1
		{
			csvs:           []*olmv1Alpha.ClusterServiceVersion{&op1Ns1},
			subscriptions:  []olmv1Alpha.Subscription{subscription1},
			installPlan:    []*olmv1Alpha.InstallPlan{&ns1InstallPlan1},
			catalogSource:  []*olmv1Alpha.CatalogSource{&catalogSource1},
			operatorGroups: []*olmv1.OperatorGroup{&operatorGroupNs1},
			expectedOperators: []*Operator{
				{
					Name:                  "op1.v1.0.1",
					Namespace:             "ns1",
					Csv:                   &op1Ns1,
					SubscriptionName:      "subs1",
					SubscriptionNamespace: "ns1",
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
					TargetNamespaces:   []string{"ns1"},
				},
			},
		},
		// ns1: csv1/subs1 - installPlan not found.
		{
			csvs:           []*olmv1Alpha.ClusterServiceVersion{&op1Ns1},
			subscriptions:  []olmv1Alpha.Subscription{subscription1},
			installPlan:    []*olmv1Alpha.InstallPlan{},
			catalogSource:  []*olmv1Alpha.CatalogSource{},
			operatorGroups: []*olmv1.OperatorGroup{&operatorGroupNs1},
			expectedOperators: []*Operator{
				{
					Name:                  "op1.v1.0.1",
					Namespace:             "ns1",
					Csv:                   &op1Ns1,
					SubscriptionName:      "subs1",
					SubscriptionNamespace: "ns1",
					InstallPlans:          nil,
					Package:               "op1",
					Org:                   "catalogSource1",
					Version:               "1.0.1",
					PackageFromCsvName:    "op1",
					TargetNamespaces:      []string{"ns1"},
				},
			},
		},
		// ns1: csv1/subs1 - bundleImage not found.
		{
			csvs:           []*olmv1Alpha.ClusterServiceVersion{&op1Ns1},
			subscriptions:  []olmv1Alpha.Subscription{subscription1},
			installPlan:    []*olmv1Alpha.InstallPlan{&ns1InstallPlan1},
			catalogSource:  []*olmv1Alpha.CatalogSource{},
			operatorGroups: []*olmv1.OperatorGroup{&operatorGroupNs1},
			expectedOperators: []*Operator{
				{
					Name:                  "op1.v1.0.1",
					Namespace:             "ns1",
					Csv:                   &op1Ns1,
					SubscriptionName:      "subs1",
					SubscriptionNamespace: "ns1",
					InstallPlans:          nil,
					Package:               "op1",
					Org:                   "catalogSource1",
					Version:               "1.0.1",
					PackageFromCsvName:    "op1",
					TargetNamespaces:      []string{"ns1"},
				},
			},
		},
		// ns1: csv1/subs1, ns2: csv2 (without subscription). The second csv
		// is he last which is added to the unique CSV names list, so that's
		// the one which will be bound to the operator struct.
		{
			csvs:           []*olmv1Alpha.ClusterServiceVersion{&op1Ns1, &op1Ns2},
			subscriptions:  []olmv1Alpha.Subscription{subscription1},
			installPlan:    []*olmv1Alpha.InstallPlan{&ns1InstallPlan1, &ns2InstallPlan1},
			catalogSource:  []*olmv1Alpha.CatalogSource{&catalogSource1, &catalogSource2},
			operatorGroups: []*olmv1.OperatorGroup{&operatorGroupNs1, &operatorGroupNs2},
			expectedOperators: []*Operator{
				{
					Name:                  "op1.v1.0.1",
					Namespace:             "ns2",
					Csv:                   &op1Ns2,
					SubscriptionName:      "subs1",
					SubscriptionNamespace: "ns1",
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
					TargetNamespaces:   []string{"ns1"},
				},
			},
		},
		// ns1: csv1, ns2: csv2/subs2. OperatorGroup in ns2 has ns1 as targetNamespace.
		{
			csvs:           []*olmv1Alpha.ClusterServiceVersion{&op1Ns1, &op1Ns2},
			subscriptions:  []olmv1Alpha.Subscription{subscription2},
			installPlan:    []*olmv1Alpha.InstallPlan{&ns1InstallPlan1, &ns2InstallPlan1},
			catalogSource:  []*olmv1Alpha.CatalogSource{&catalogSource1, &catalogSource2},
			operatorGroups: []*olmv1.OperatorGroup{&operatorGroupNs2},
			expectedOperators: []*Operator{
				{
					Name:                  "op1.v1.0.1",
					Namespace:             "ns2",
					Csv:                   &op1Ns2,
					SubscriptionName:      "subs2",
					SubscriptionNamespace: "ns2",
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
					TargetNamespaces:   []string{"ns2"},
				},
			},
		},
		// ns1: csv1, ns2: csv2/subs2 + csv3/subs3
		{
			csvs:           []*olmv1Alpha.ClusterServiceVersion{&op1Ns1, &op1Ns2, &op2Ns2},
			subscriptions:  []olmv1Alpha.Subscription{subscription2, subscription3},
			installPlan:    []*olmv1Alpha.InstallPlan{&ns1InstallPlan1, &ns2InstallPlan1, &ns2InstallPlan2},
			catalogSource:  []*olmv1Alpha.CatalogSource{&catalogSource1, &catalogSource2, &catalogSource3},
			operatorGroups: []*olmv1.OperatorGroup{&operatorGroupNs1, &operatorGroupNs2},
			expectedOperators: []*Operator{
				{
					Name:                  "op1.v1.0.1",
					Namespace:             "ns2",
					Csv:                   &op1Ns2,
					SubscriptionName:      "subs2",
					SubscriptionNamespace: "ns2",
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
					TargetNamespaces:   []string{"ns2"},
				},
				{
					Name:                  "op2.v2.0.2",
					Namespace:             "ns2",
					Csv:                   &op2Ns2,
					SubscriptionName:      "subs3",
					SubscriptionNamespace: "ns2",
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
					TargetNamespaces:   []string{"ns2"},
				},
			},
		},
	}

	for _, tc := range testCases {
		// Create runtime objects from operatorGroups for the fake OLM client.
		runtimeObjects := []runtime.Object{}
		for _, opGroup := range tc.operatorGroups {
			runtimeObjects = append(runtimeObjects, opGroup)
		}

		// Reinstantiate fake client and load runtimeObjects to OLM.
		_ = clientsholder.GetTestClientsHolder(nil)
		clientsholder.SetupFakeOlmClient(runtimeObjects)

		ops := createOperators(tc.csvs, tc.subscriptions, tc.installPlan, tc.catalogSource, false, true)
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

func createOperator(phase, packageName, aVersion, namespace string, targetNamespaces []string, isClusterWide bool) *Operator {
	aOperator := Operator{}
	aOperator.PackageFromCsvName = packageName
	aCsv := olmv1Alpha.ClusterServiceVersion{}
	aCsv.Status.Phase = olmv1Alpha.ClusterServiceVersionPhase(phase)
	aOperator.Csv = &aCsv
	aOperator.Namespace = namespace
	aOperator.Version = aVersion
	aOperator.TargetNamespaces = targetNamespaces
	aOperator.Phase = olmv1Alpha.ClusterServiceVersionPhase(phase)
	aOperator.IsClusterWide = isClusterWide
	return &aOperator
}

func createOperatorList() []*Operator {
	aOperatorList := []*Operator{}
	aOperatorList = append(aOperatorList,
		createOperator("Succeeded", "operator1", "0.0.1", "ns1", []string{"ns1"}, false),
		createOperator("Succeeded", "operator2", "0.0.1", "ns2", []string{"ns2"}, false),
		createOperator("Succeeded", "operator3", "1.0.1", "ns3", nil, true),
		createOperator("Succeeded", "operator3", "0.0.1", "ns4", nil, true),
		createOperator("InstallReady", "operator3", "0.0.1", "ns5", nil, true),
		createOperator("Succeeded", "operator3", "0.0.1", "ns6", nil, true),
		createOperator("Succeeded", "operator3", "0.0.1", "ns7", nil, true),
		createOperator("Succeeded", "operator3", "0.0.1", "ns8", nil, true),
		createOperator("Failed", "operator3", "0.0.1", "ns6", nil, true),
		createOperator("Failed", "operator4", "0.0.1", "ns1", []string{"ns1"}, false))
	return aOperatorList
}

func Test_getSummaryAllOperators(t *testing.T) {
	type args struct {
		operators []*Operator
	}
	tests := []struct {
		name        string
		args        args
		wantSummary []string
	}{
		{
			name: "ok",
			args: args{operators: createOperatorList()},
			wantSummary: []string{"Failed operator: operator3 ver: 0.0.1 (all namespaces)",
				"Failed operator: operator4 ver: 0.0.1 in ns: [ns1]",
				"InstallReady operator: operator3 ver: 0.0.1 (all namespaces)",
				"Succeeded operator: operator1 ver: 0.0.1 in ns: [ns1]",
				"Succeeded operator: operator2 ver: 0.0.1 in ns: [ns2]",
				"Succeeded operator: operator3 ver: 0.0.1 (all namespaces)",
				"Succeeded operator: operator3 ver: 1.0.1 (all namespaces)"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotSummary := getSummaryAllOperators(tt.args.operators); !reflect.DeepEqual(gotSummary, tt.wantSummary) {
				t.Errorf("getSummaryAllOperators() = %v, want %v", gotSummary, tt.wantSummary)
			}
		})
	}
}

func TestGetUniqueCsvListByName(t *testing.T) {
	// CSVs from 3 different fake operators:
	//  Namespace ns1: csv1 (op1) & csv2 (op2)
	//  Namespace ns2: csv3 (op1) & csv4 (op2)
	//  Namespace ns3: csv5 (op3)
	csv1 := createCsv("op1.v1.0.1", "ns1", 1, 0, 1)
	csv2 := createCsv("op1.v1.0.1", "ns2", 1, 0, 1)
	csv3 := createCsv("op2.v2.2.2", "ns1", 2, 2, 2)
	csv4 := createCsv("op2.v2.2.2", "ns2", 2, 2, 2)
	csv5 := createCsv("op3.v1.0.0", "ns3", 2, 0, 2)

	tcs := []struct {
		csvs               []*olmv1Alpha.ClusterServiceVersion
		expectedUniqueCsvs []*olmv1Alpha.ClusterServiceVersion
	}{
		{
			csvs:               nil,
			expectedUniqueCsvs: []*olmv1Alpha.ClusterServiceVersion{},
		},
		{
			csvs:               []*olmv1Alpha.ClusterServiceVersion{},
			expectedUniqueCsvs: []*olmv1Alpha.ClusterServiceVersion{},
		},
		{
			csvs:               []*olmv1Alpha.ClusterServiceVersion{&csv1},
			expectedUniqueCsvs: []*olmv1Alpha.ClusterServiceVersion{&csv1},
		},
		{
			csvs:               []*olmv1Alpha.ClusterServiceVersion{&csv1, &csv2, &csv3},
			expectedUniqueCsvs: []*olmv1Alpha.ClusterServiceVersion{&csv2, &csv3},
		},
		{
			csvs:               []*olmv1Alpha.ClusterServiceVersion{&csv1, &csv2, &csv3, &csv4},
			expectedUniqueCsvs: []*olmv1Alpha.ClusterServiceVersion{&csv2, &csv4},
		},
		{
			csvs:               []*olmv1Alpha.ClusterServiceVersion{&csv1, &csv2, &csv3, &csv4, &csv5},
			expectedUniqueCsvs: []*olmv1Alpha.ClusterServiceVersion{&csv2, &csv4, &csv5},
		},
	}

	for _, tc := range tcs {
		runtimeObjects := []runtime.Object{}
		for _, opGroup := range tc.csvs {
			runtimeObjects = append(runtimeObjects, opGroup)
		}

		// Reinstantiate fake client and load runtimeObjects to OLM.
		_ = clientsholder.GetTestClientsHolder(nil)
		clientsholder.SetupFakeOlmClient(runtimeObjects)

		uniqueCsvs := getUniqueCsvListByName(tc.csvs)
		assert.Equal(t, tc.expectedUniqueCsvs, uniqueCsvs)
	}
}

func TestIsValidSemanticVersioning(t *testing.T) {
	testCases := []struct {
		testVersion    string
		expectedResult bool
	}{
		{
			testVersion:    "1.1.1",
			expectedResult: true,
		},
		{
			testVersion:    "1.1.2-prerelease+meta",
			expectedResult: true,
		},
		{
			testVersion:    "1.1.2+meta",
			expectedResult: true,
		},
		{
			testVersion:    "1.0.0-alpha",
			expectedResult: true,
		},
		{
			testVersion:    "1.0.0-alpha.beta",
			expectedResult: true,
		},
		{
			testVersion:    "1.0.0-alpha.1",
			expectedResult: true,
		},
		{
			testVersion:    "1.0.0-alpha.0valid",
			expectedResult: true,
		},
		{
			testVersion:    "1.0.0-rc.1+build.1",
			expectedResult: true,
		},
		{
			testVersion:    "10.2.3-DEV-SNAPSHOT",
			expectedResult: true,
		},
		{
			testVersion:    "2.0.0+build.1848",
			expectedResult: true,
		},
		{
			testVersion:    "2.0.1-alpha.1227",
			expectedResult: true,
		},
		{
			testVersion:    "1.0.0-alpha+beta",
			expectedResult: true,
		},
		{
			testVersion:    "1.2.3----RC-SNAPSHOT.12.9.1--.12+788",
			expectedResult: true,
		},

		{
			testVersion:    "1.2.3----R-S.12.9.1--.12+meta",
			expectedResult: true,
		},
		{
			testVersion:    "1.2.3.4",
			expectedResult: false,
		},
		{
			testVersion:    "1.2.3.4.5",
			expectedResult: false,
		},
		{
			testVersion:    "1.2.3+",
			expectedResult: false,
		},
		{
			testVersion:    "hello.v1.1.1",
			expectedResult: false,
		},
	}

	for _, tc := range testCases {
		assert.Equal(t, tc.expectedResult, IsValidSemanticVersion(tc.testVersion))
	}
}
