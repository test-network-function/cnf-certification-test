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
	"bytes"
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/redhat-best-practices-for-k8s/certsuite/tests/common"
	"github.com/redhat-best-practices-for-k8s/certsuite/tests/identifiers"
	pdbv1 "github.com/redhat-best-practices-for-k8s/certsuite/tests/observability/pdb"

	"github.com/redhat-best-practices-for-k8s/certsuite/internal/clientsholder"
	"github.com/redhat-best-practices-for-k8s/certsuite/internal/log"
	"github.com/redhat-best-practices-for-k8s/certsuite/pkg/checksdb"
	"github.com/redhat-best-practices-for-k8s/certsuite/pkg/provider"
	"github.com/redhat-best-practices-for-k8s/certsuite/pkg/testhelper"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

var (
	env provider.TestEnvironment

	beforeEachFn = func(check *checksdb.Check) error {
		env = provider.GetTestEnvironment()
		return nil
	}
)

func LoadChecks() {
	log.Debug("Loading %s suite checks", common.ObservabilityTestKey)

	checksGroup := checksdb.NewChecksGroup(common.ObservabilityTestKey).
		WithBeforeEachFn(beforeEachFn)

	checksGroup.Add(checksdb.NewCheck(identifiers.GetTestIDAndLabels(identifiers.TestLoggingIdentifier)).
		WithSkipCheckFn(testhelper.GetNoContainersUnderTestSkipFn(&env)).
		WithCheckFn(func(c *checksdb.Check) error {
			testContainersLogging(c, &env)
			return nil
		}))

	checksGroup.Add(checksdb.NewCheck(identifiers.GetTestIDAndLabels(identifiers.TestCrdsStatusSubresourceIdentifier)).
		WithSkipCheckFn(testhelper.GetNoCrdsUnderTestSkipFn(&env)).
		WithCheckFn(func(c *checksdb.Check) error {
			testCrds(c, &env)
			return nil
		}))

	checksGroup.Add(checksdb.NewCheck(identifiers.GetTestIDAndLabels(identifiers.TestTerminationMessagePolicyIdentifier)).
		WithSkipCheckFn(testhelper.GetNoContainersUnderTestSkipFn(&env)).
		WithCheckFn(func(c *checksdb.Check) error {
			testTerminationMessagePolicy(c, &env)
			return nil
		}))

	checksGroup.Add(checksdb.NewCheck(identifiers.GetTestIDAndLabels(identifiers.TestPodDisruptionBudgetIdentifier)).
		WithSkipCheckFn(testhelper.GetNoDeploymentsUnderTestSkipFn(&env), testhelper.GetNoStatefulSetsUnderTestSkipFn(&env)).
		WithSkipModeAll().
		WithCheckFn(func(c *checksdb.Check) error {
			testPodDisruptionBudgets(c, &env)
			return nil
		}))

	checksGroup.Add(checksdb.NewCheck(identifiers.GetTestIDAndLabels(identifiers.TestAPICompatibilityWithNextOCPReleaseIdentifier)).
		WithCheckFn(func(c *checksdb.Check) error {
			testAPICompatibilityWithNextOCPRelease(c, &env)
			return nil
		}))
}

// containerHasLoggingOutput helper function to get the last line of logging output from
// a container. Returns true in case some output was found, false otherwise.
func containerHasLoggingOutput(cut *provider.Container) (bool, error) {
	ocpClient := clientsholder.GetClientsHolder()

	// K8s' API won't return lines that do not have the newline termination char, so
	// We need to ask for the last two lines.
	const tailLogLines = 2
	numLogLines := int64(tailLogLines)
	podLogOptions := corev1.PodLogOptions{TailLines: &numLogLines, Container: cut.Name}
	req := ocpClient.K8sClient.CoreV1().Pods(cut.Namespace).GetLogs(cut.Podname, &podLogOptions)

	podLogsReaderCloser, err := req.Stream(context.TODO())
	if err != nil {
		return false, fmt.Errorf("unable to get log streamer, err: %v", err)
	}

	defer podLogsReaderCloser.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, podLogsReaderCloser)
	if err != nil {
		return false, fmt.Errorf("unable to get log data, err: %v", err)
	}

	return buf.String() != "", nil
}

func testContainersLogging(check *checksdb.Check, env *provider.TestEnvironment) {
	// Iterate through all the CUTs to get their log output. The TC checks that at least
	// one log line is found.
	var compliantObjects []*testhelper.ReportObject
	var nonCompliantObjects []*testhelper.ReportObject
	for _, cut := range env.Containers {
		check.LogInfo("Testing Container %q", cut)
		hasLoggingOutput, err := containerHasLoggingOutput(cut)
		if err != nil {
			check.LogError("Failed to get %q log output, err: %v", cut, err)
			nonCompliantObjects = append(nonCompliantObjects,
				testhelper.NewContainerReportObject(cut.Namespace, cut.Podname, cut.Name, "Could not get log output", false))
			continue
		}

		if !hasLoggingOutput {
			check.LogError("Container %q does not have any line of log to stderr/stdout", cut)
			nonCompliantObjects = append(nonCompliantObjects,
				testhelper.NewContainerReportObject(cut.Namespace, cut.Podname, cut.Name, "No log line to stderr/stdout found", false))
		} else {
			check.LogInfo("Container %q has some logging output", cut)
			compliantObjects = append(compliantObjects,
				testhelper.NewContainerReportObject(cut.Namespace, cut.Podname, cut.Name, "Found log line to stderr/stdout", true))
		}
	}

	check.SetResult(compliantObjects, nonCompliantObjects)
}

// testCrds testing if crds have a status sub resource set
func testCrds(check *checksdb.Check, env *provider.TestEnvironment) {
	var compliantObjects []*testhelper.ReportObject
	var nonCompliantObjects []*testhelper.ReportObject
	for _, crd := range env.Crds {
		check.LogInfo("Testing CRD: %s", crd.Name)
		for _, ver := range crd.Spec.Versions {
			if _, ok := ver.Schema.OpenAPIV3Schema.Properties["status"]; !ok {
				check.LogError("CRD: %s, version: %s does not have a status subresource", crd.Name, ver.Name)
				nonCompliantObjects = append(nonCompliantObjects,
					testhelper.NewReportObject("Crd does not have a status sub resource set", testhelper.CustomResourceDefinitionType, false).
						AddField(testhelper.CustomResourceDefinitionName, crd.Name).
						AddField(testhelper.CustomResourceDefinitionVersion, ver.Name))
			} else {
				check.LogInfo("CRD: %s, version: %s has a status subresource", crd.Name, ver.Name)
				compliantObjects = append(compliantObjects,
					testhelper.NewReportObject("Crd has a status sub resource set", testhelper.CustomResourceDefinitionType, true).
						AddField(testhelper.CustomResourceDefinitionName, crd.Name).
						AddField(testhelper.CustomResourceDefinitionVersion, ver.Name))
			}
		}
	}

	check.SetResult(compliantObjects, nonCompliantObjects)
}

// testTerminationMessagePolicy tests to make sure that pods
func testTerminationMessagePolicy(check *checksdb.Check, env *provider.TestEnvironment) {
	var compliantObjects []*testhelper.ReportObject
	var nonCompliantObjects []*testhelper.ReportObject

	for _, cut := range env.Containers {
		check.LogInfo("Testing Container %q", cut)
		if cut.TerminationMessagePolicy != corev1.TerminationMessageFallbackToLogsOnError {
			check.LogError("Container %q does not have a TerminationMessagePolicy: FallbackToLogsOnError (has %s)", cut, cut.TerminationMessagePolicy)
			nonCompliantObjects = append(nonCompliantObjects,
				testhelper.NewContainerReportObject(cut.Namespace, cut.Podname, cut.Name, "TerminationMessagePolicy is not FallbackToLogsOnError", false))
		} else {
			check.LogInfo("Container %q has a TerminationMessagePolicy: FallbackToLogsOnError", cut)
			compliantObjects = append(compliantObjects,
				testhelper.NewContainerReportObject(cut.Namespace, cut.Podname, cut.Name, "TerminationMessagePolicy is FallbackToLogsOnError", true))
		}
	}

	check.SetResult(compliantObjects, nonCompliantObjects)
}

//nolint:funlen
func testPodDisruptionBudgets(check *checksdb.Check, env *provider.TestEnvironment) {
	var compliantObjects []*testhelper.ReportObject
	var nonCompliantObjects []*testhelper.ReportObject

	// Loop through all of the of Deployments and StatefulSets and check if the PDBs are valid
	for _, d := range env.Deployments {
		check.LogInfo("Testing Deployment %q", d.ToString())
		deploymentSelector := labels.Set(d.Spec.Template.Labels)
		pdbFound := false
		for pdbIndex := range env.PodDisruptionBudgets {
			pdb := &env.PodDisruptionBudgets[pdbIndex]
			if pdb.Namespace != d.Namespace {
				continue
			}
			pdbSelector, err := metav1.LabelSelectorAsSelector(pdb.Spec.Selector)
			if err != nil {
				check.LogError("Could not convert the PDB %q label selector to selector, err: %v", pdbSelector, err)
				continue
			}
			if pdbSelector.Matches(deploymentSelector) {
				pdbFound = true
				if ok, err := pdbv1.CheckPDBIsValid(pdb, d.Spec.Replicas); !ok {
					check.LogError("PDB %q is not valid for Deployment %q, err: %v", pdb.Name, d.Name, err)
					nonCompliantObjects = append(nonCompliantObjects, testhelper.NewReportObject(fmt.Sprintf("Invalid PodDisruptionBudget config: %v", err), testhelper.DeploymentType, false).
						AddField(testhelper.DeploymentName, d.Name).
						AddField(testhelper.Namespace, d.Namespace).
						AddField(testhelper.PodDisruptionBudgetReference, pdb.Name))
				} else {
					check.LogInfo("PDB %q is valid for Deployment: %q", pdb.Name, d.Name)
					compliantObjects = append(compliantObjects, testhelper.NewReportObject("Deployment: references PodDisruptionBudget", testhelper.DeploymentType, true).
						AddField(testhelper.DeploymentName, d.Name).
						AddField(testhelper.Namespace, d.Namespace).
						AddField(testhelper.PodDisruptionBudgetReference, pdb.Name))
				}
			}
		}
		if !pdbFound {
			check.LogError("Deployment %q is missing a corresponding PodDisruptionBudget", d.ToString())
			nonCompliantObjects = append(nonCompliantObjects, testhelper.NewReportObject("Deployment is missing a corresponding PodDisruptionBudget", testhelper.DeploymentType, false).
				AddField(testhelper.DeploymentName, d.Name).
				AddField(testhelper.Namespace, d.Namespace))
		}
	}

	for _, s := range env.StatefulSets {
		check.LogInfo("Testing StatefulSet %q", s.ToString())
		statefulSetSelector := labels.Set(s.Spec.Template.Labels)
		pdbFound := false
		for pdbIndex := range env.PodDisruptionBudgets {
			pdb := &env.PodDisruptionBudgets[pdbIndex]
			if pdb.Namespace != s.Namespace {
				continue
			}
			pdbSelector, err := metav1.LabelSelectorAsSelector(pdb.Spec.Selector)
			if err != nil {
				check.LogError("Could not convert the PDB %q label selector to selector, err: %v", pdbSelector, err)
				continue
			}
			if pdbSelector.Matches(statefulSetSelector) {
				pdbFound = true
				if ok, err := pdbv1.CheckPDBIsValid(pdb, s.Spec.Replicas); !ok {
					check.LogError("PDB %q is not valid for StatefulSet %q, err: %v", pdb.Name, s.Name, err)
					nonCompliantObjects = append(nonCompliantObjects, testhelper.NewReportObject(fmt.Sprintf("Invalid PodDisruptionBudget config: %v", err), testhelper.StatefulSetType, false).
						AddField(testhelper.StatefulSetName, s.Name).
						AddField(testhelper.Namespace, s.Namespace).
						AddField(testhelper.PodDisruptionBudgetReference, pdb.Name))
				} else {
					check.LogInfo("PDB %q is valid for StatefulSet: %q", pdb.Name, s.Name)
					compliantObjects = append(compliantObjects, testhelper.NewReportObject("StatefulSet: references PodDisruptionBudget", testhelper.StatefulSetType, true).
						AddField(testhelper.StatefulSetName, s.Name).
						AddField(testhelper.Namespace, s.Namespace).
						AddField(testhelper.PodDisruptionBudgetReference, pdb.Name))
				}
			}
		}
		if !pdbFound {
			check.LogError("StatefulSet %q is missing a corresponding PodDisruptionBudget", s.ToString())
			nonCompliantObjects = append(nonCompliantObjects, testhelper.NewReportObject("StatefulSet is missing a corresponding PodDisruptionBudget", testhelper.StatefulSetType, false).
				AddField(testhelper.StatefulSetName, s.Name).
				AddField(testhelper.Namespace, s.Namespace))
		}
	}

	check.SetResult(compliantObjects, nonCompliantObjects)
}

// Convert the list of APIRequestCount into a list of ToBeRemovedWorkloadAPI
// using the following filters:
// - status.removedInRelease is not nil
// - service accounts in status.last24.byNode.ByUser.username do not contain the 'openshift' or 'kube' regex
func buildToBeRemovedWorkloadAPIList(apiRequestCounts []provider.APIRequestCount) []provider.ToBeRemovedWorkloadAPI {
	var toBeRemovedWorkloadAPIs []provider.ToBeRemovedWorkloadAPI

	for _, obj := range apiRequestCounts {
		if obj.Status.RemovedInRelease != "" {
			workloadAPI := provider.ToBeRemovedWorkloadAPI{
				APIName:          obj.Metadata.Name,
				RemovedInRelease: obj.Status.RemovedInRelease,
				UserInfo:         make(map[string]struct{}),
			}

			for _, last24h := range obj.Status.Last24h {
				for _, byNode := range last24h.ByNode {
					if byNode.ByUser != nil {
						for _, byUser := range byNode.ByUser {
							if !strings.Contains(byUser.UserName, "openshift") && !strings.Contains(byUser.UserName, "kube") {
								userInfo := fmt.Sprintf("UserName: %s, UserAgent: %s", byUser.UserName, byUser.UserAgent)
								workloadAPI.UserInfo[userInfo] = struct{}{}
							}
						}
					}
				}
			}

			if len(workloadAPI.UserInfo) > 0 {
				toBeRemovedWorkloadAPIs = append(toBeRemovedWorkloadAPIs, workloadAPI)
			}
		}
	}

	return toBeRemovedWorkloadAPIs
}

// Evaluate workload API compliance with the next Kubernetes version
func evaluateAPICompliance(toBeRemovedWorkloadAPIs []provider.ToBeRemovedWorkloadAPI, kubernetesVersion string) (compliantObjects []*testhelper.ReportObject, nonCompliantObjects []*testhelper.ReportObject) {
	// Extract the major version from the current Kubernetes version
	minorVersion, err := strconv.Atoi(strings.Split(kubernetesVersion, ".")[1])
	if err != nil {
		fmt.Printf("Error parsing minor version from Kubernetes version: %s\n", err)
		return nil, nil
	}

	for _, api := range toBeRemovedWorkloadAPIs {
		removedMinorVersion, err := strconv.Atoi(strings.Split(api.RemovedInRelease, ".")[1])
		if err != nil {
			fmt.Printf("Error parsing minor version from APIRequestCount.status.removedInRelease: %s\n", err)
			// Skip this API if the version parsing fails
			continue
		}

		isCompliantWithNextK8sVersion := removedMinorVersion > minorVersion+1
		reportObject := testhelper.NewReportObject("API", "StubType", isCompliantWithNextK8sVersion)
		reportObject.AddField("APIName", api.APIName)
		reportObject.AddField("RemovedInRelease", api.RemovedInRelease)
		for userInfo := range api.UserInfo {
			reportObject.AddField("UserInfo", userInfo)
		}

		if isCompliantWithNextK8sVersion {
			compliantObjects = append(compliantObjects, reportObject)
		} else {
			nonCompliantObjects = append(nonCompliantObjects, reportObject)
		}
	}

	return compliantObjects, nonCompliantObjects
}

// Function to test API compatibility with the next OCP release
func testAPICompatibilityWithNextOCPRelease(check *checksdb.Check, env *provider.TestEnvironment) {
	check.LogInfo("Running testAPICompatibilityWithNextOCPRelease")

	isOCP := provider.IsOCPCluster()
	check.LogInfo("Is OCP: %v", isOCP)

	if isOCP {
		// Retrieve APIRequestCount objects with ocpClient
		ocpClient := clientsholder.GetClientsHolder().OcpClient.RESTClient()
		apiRequestCounts, err := provider.RetrieveAPIRequestCounts(ocpClient)
		if err != nil {
			check.LogError("Error retrieving APIRequestCount objects: %s", err)
			return
		}

		// Build the list of ToBeRemovedWorkloadAPI
		toBeRemovedWorkloadAPIs := buildToBeRemovedWorkloadAPIList(apiRequestCounts)

		// Evaluate API compliance with the next Kubernetes version
		compliantObjects, nonCompliantObjects := evaluateAPICompliance(toBeRemovedWorkloadAPIs, env.K8sVersion)

		// Add test results
		check.SetResult(compliantObjects, nonCompliantObjects)
	} else {
		check.LogInfo("The Kubernetes distribution is not OpenShift. Skipping API compatibility test.")
	}
}
