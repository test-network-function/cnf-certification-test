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

package configuration

import "time"

const (
	defaultDebugDaemonSetNamespace = "cnf-suite"
)

type SkipHelmChartList struct {
	// Name is the name of the `operator bundle package name` or `image-version` that you want to check if exists in the RedHat catalog.
	Name string `yaml:"name" json:"name"`
}

// AcceptedKernelTaintsInfo contains all certified operator request info.
type AcceptedKernelTaintsInfo struct {

	// Accepted modules that cause taints that we want to supply to the test suite.
	Module string `yaml:"module" json:"module"`
}

// SkipScalingTestDeploymentsInfo contains a list of names of deployments that should be skipped by the scaling tests to prevent issues.
type SkipScalingTestDeploymentsInfo struct {

	// Deployment name and namespace that can be skipped by the scaling tests.
	Name      string `yaml:"name" json:"name"`
	Namespace string `yaml:"namespace" json:"namespace"`
}

// SkipScalingTestStatefulSetsInfo contains a list of names of statefulsets that should be skipped by the scaling tests to prevent issues.
type SkipScalingTestStatefulSetsInfo struct {

	// StatefulSet name and namespace that can be skipped by the scaling tests.
	Name      string `yaml:"name" json:"name"`
	Namespace string `yaml:"namespace" json:"namespace"`
}

// Namespace struct defines namespace properties.
type Namespace struct {
	Name string `yaml:"name" json:"name"`
}

// CrdFilter defines a CustomResourceDefinition config filter.
type CrdFilter struct {
	NameSuffix string `yaml:"nameSuffix" json:"nameSuffix"`
	Scalable   bool   `yaml:"scalable" json:"scalable"`
}
type ManagedDeploymentsStatefulsets struct {
	Name string `yaml:"name" json:"name"`
}

// TestConfiguration provides test related configuration.
type TestConfiguration struct {
	// targetNameSpaces to be used in
	TargetNameSpaces []Namespace `yaml:"targetNameSpaces,omitempty" json:"targetNameSpaces,omitempty"`
	// labels identifying pods under test
	PodsUnderTestLabels []string `yaml:"podsUnderTestLabels,omitempty" json:"podsUnderTestLabels,omitempty"`
	// labels identifying operators unde test
	OperatorsUnderTestLabels []string `yaml:"operatorsUnderTestLabels,omitempty" json:"operatorsUnderTestLabels,omitempty"`
	// CRDs section.
	CrdFilters          []CrdFilter                      `yaml:"targetCrdFilters,omitempty" json:"targetCrdFilters,omitempty"`
	ManagedDeployments  []ManagedDeploymentsStatefulsets `yaml:"managedDeployments,omitempty" json:"managedDeployments,omitempty"`
	ManagedStatefulsets []ManagedDeploymentsStatefulsets `yaml:"managedStatefulsets,omitempty" json:"managedStatefulsets,omitempty"`

	// AcceptedKernelTaints
	AcceptedKernelTaints []AcceptedKernelTaintsInfo `yaml:"acceptedKernelTaints,omitempty" json:"acceptedKernelTaints,omitempty"`
	SkipHelmChartList    []SkipHelmChartList        `yaml:"skipHelmChartList,omitempty" json:"skipHelmChartList,omitempty"`
	// SkipScalingTestDeploymentNames
	SkipScalingTestDeployments []SkipScalingTestDeploymentsInfo `yaml:"skipScalingTestDeployments,omitempty" json:"skipScalingTestDeployments,omitempty"`
	// SkipScalingTestStatefulSetNames
	SkipScalingTestStatefulSets []SkipScalingTestStatefulSetsInfo `yaml:"skipScalingTestStatefulSets,omitempty" json:"skipScalingTestStatefulSets,omitempty"`
	ValidProtocolNames          []string                          `yaml:"validProtocolNames,omitempty" json:"validProtocolNames,omitempty"`
	ServicesIgnoreList          []string                          `yaml:"servicesignorelist,omitempty" json:"servicesignorelist,omitempty"`
	DebugDaemonSetNamespace     string                            `yaml:"debugDaemonSetNamespace,omitempty" json:"debugDaemonSetNamespace,omitempty"`
	// Collector's parameters
	ExecutedBy           string `yaml:"executedBy,omitempty" json:"executedBy,omitempty"`
	PartnerName          string `yaml:"partnerName,omitempty" json:"partnerName,omitempty"`
	CollectorAppPassword string `yaml:"collectorAppPassword,omitempty" json:"collectorAppPassword,omitempty"`
	CollectorAppEndpoint string `yaml:"collectorAppEndpoint,omitempty" json:"collectorAppEndpoint,omitempty"`
}

type TestParameters struct {
	Kubeconfig                    string
	ConfigFile                    string
	PfltDockerconfig              string
	OutputDir                     string
	LabelsFilter                  string
	LogLevel                      string
	OfflineDB                     string
	DaemonsetCPUReq               string
	DaemonsetCPULim               string
	DaemonsetMemReq               string
	DaemonsetMemLim               string
	SanitizeClaim                 bool
	TnfImageRepo                  string
	TnfDebugImage                 string
	NonIntrusiveOnly              bool
	AllowPreflightInsecure        bool
	IncludeWebFilesInOutputFolder bool
	OmitArtifactsZipFile          bool
	EnableDataCollection          bool
	EnableXMLCreation             bool
	ServerMode                    bool
	Timeout                       time.Duration
}
