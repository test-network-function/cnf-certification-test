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

package operator

import (
	"strings"

	"github.com/test-network-function/cnf-certification-test/cnf-certification-test/common"
	"github.com/test-network-function/cnf-certification-test/cnf-certification-test/common/rbac"
	"github.com/test-network-function/cnf-certification-test/cnf-certification-test/identifiers"
	"github.com/test-network-function/cnf-certification-test/cnf-certification-test/operator/phasecheck"
	"github.com/test-network-function/cnf-certification-test/internal/clientsholder"

	"github.com/test-network-function/cnf-certification-test/internal/log"
	"github.com/test-network-function/cnf-certification-test/pkg/checksdb"
	"github.com/test-network-function/cnf-certification-test/pkg/provider"
	"github.com/test-network-function/cnf-certification-test/pkg/testhelper"
	"github.com/test-network-function/cnf-certification-test/pkg/versions"
)

var (
	env provider.TestEnvironment

	beforeEachFn = func(check *checksdb.Check) error {
		env = provider.GetTestEnvironment()
		return nil
	}
)

//nolint:funlen
func LoadChecks() {
	log.Debug("Loading %s suite checks", common.OperatorTestKey)

	checksGroup := checksdb.NewChecksGroup(common.OperatorTestKey).
		WithBeforeEachFn(beforeEachFn)

	checksGroup.Add(checksdb.NewCheck(identifiers.GetTestIDAndLabels(identifiers.TestOperatorInstallStatusSucceededIdentifier)).
		WithSkipCheckFn(testhelper.GetNoOperatorsSkipFn(&env)).
		WithCheckFn(func(c *checksdb.Check) error {
			testOperatorInstallationPhaseSucceeded(c, &env)
			return nil
		}))

	checksGroup.Add(checksdb.NewCheck(identifiers.GetTestIDAndLabels(identifiers.TestOperatorNoSCCAccess)).
		WithSkipCheckFn(testhelper.GetNoOperatorsSkipFn(&env)).
		WithCheckFn(func(c *checksdb.Check) error {
			testOperatorInstallationAccessToSCC(c, &env)
			return nil
		}))

	checksGroup.Add(checksdb.NewCheck(identifiers.GetTestIDAndLabels(identifiers.TestOperatorIsInstalledViaOLMIdentifier)).
		WithSkipCheckFn(testhelper.GetNoOperatorsSkipFn(&env)).
		WithCheckFn(func(c *checksdb.Check) error {
			testOperatorOlmSubscription(c, &env)
			return nil
		}))

	checksGroup.Add(checksdb.NewCheck(identifiers.GetTestIDAndLabels(identifiers.TestOperatorHasSemanticVersioningIdentifier)).
		WithSkipCheckFn(testhelper.GetNoOperatorsSkipFn(&env)).
		WithCheckFn(func(c *checksdb.Check) error {
			testOperatorSemanticVersioning(c, &env)
			return nil
		}))

	checksGroup.Add(checksdb.NewCheck(identifiers.GetTestIDAndLabels(identifiers.TestOperatorCrdVersioningIdentifier)).
		WithSkipCheckFn(testhelper.GetNoOperatorCrdsSkipFn(&env)).
		WithCheckFn(func(c *checksdb.Check) error {
			testOperatorCrdVersioning(c, &env)
			return nil
		}))

	checksGroup.Add(checksdb.NewCheck(identifiers.GetTestIDAndLabels(identifiers.TestOperatorCrdSchemaIdentifier)).
		WithSkipCheckFn(testhelper.GetNoOperatorCrdsSkipFn(&env)).
		WithCheckFn(func(c *checksdb.Check) error {
			testOperatorCrdOpenAPISpec(c, &env)
			return nil
		}))

	checksGroup.Add(checksdb.NewCheck(identifiers.GetTestIDAndLabels(identifiers.TestOperatorSingleCrdOwnerIdentifier)).
		WithSkipCheckFn(testhelper.GetNoOperatorsSkipFn(&env)).
		WithCheckFn(func(c *checksdb.Check) error {
			testOperatorSingleCrdOwner(c, &env)
			return nil
		}))

	checksGroup.Add(checksdb.NewCheck(identifiers.GetTestIDAndLabels(identifiers.TestOperatorRunAsUserID)).
		WithSkipCheckFn(testhelper.GetNoOperatorsSkipFn(&env)).
		WithCheckFn(func(c *checksdb.Check) error {
			testOperatorPodsRunAsUserID(c, &env)
			return nil
		}))
	checksGroup.Add(checksdb.NewCheck(identifiers.GetTestIDAndLabels(identifiers.TestOperatorRunAsNonRoot)).
		WithSkipCheckFn(testhelper.GetNoOperatorsSkipFn(&env)).
		WithCheckFn(func(c *checksdb.Check) error {
			testOperatorPodsRunAsNonRoot(c, &env)
			return nil
		}))

	checksGroup.Add(checksdb.NewCheck(identifiers.GetTestIDAndLabels(identifiers.TestOperatorAutomountTokens)).
		WithSkipCheckFn(testhelper.GetNoOperatorsSkipFn(&env)).
		WithCheckFn(func(c *checksdb.Check) error {
			testOperatorPodsAutomountTokens(c, &env)
			return nil
		}))

	checksGroup.Add(checksdb.NewCheck(identifiers.GetTestIDAndLabels(identifiers.TestOperatorReadOnlyFilesystem)).
		WithSkipCheckFn(testhelper.GetNoOperatorsSkipFn(&env)).
		WithCheckFn(func(c *checksdb.Check) error {
			testOperatorContainersReadOnlyFilesystem(c, &env)
			return nil
		}))
}

// This function check if the Operator CRD version follows K8s versioning
func testOperatorCrdVersioning(check *checksdb.Check, env *provider.TestEnvironment) {
	check.LogInfo("Starting testOperatorCrdVersioning")
	var compliantObjects []*testhelper.ReportObject
	var nonCompliantObjects []*testhelper.ReportObject

	for _, crd := range env.Crds {
		doesUseK8sVersioning := true
		nonCompliantVersion := ""

		for _, crdVersion := range crd.Spec.Versions {
			versionName := crdVersion.Name
			check.LogDebug("Checking for Operator CRD %s with version %s", crd.Name, versionName)

			if !versions.IsValidK8sVersion(versionName) {
				doesUseK8sVersioning = false
				nonCompliantVersion = versionName
				break
			}
		}

		if doesUseK8sVersioning {
			check.LogInfo("Operator CRD %s has valid K8s versioning ", crd.Name)
			compliantObjects = append(compliantObjects, testhelper.NewOperatorReportObject(crd.Namespace, crd.Name,
				"Operator CRD has valid K8s versioning ", true).AddField(testhelper.CrdVersion, crd.Name))
		} else {
			check.LogError("Operator CRD %s has invalid K8s versioning %s ", crd.Name, nonCompliantVersion)
			nonCompliantObjects = append(nonCompliantObjects, testhelper.NewOperatorReportObject(crd.Namespace, crd.Name,
				"Operator CRD has invalid K8s versioning ", false).AddField(testhelper.CrdVersion, crd.Name))
		}
	}
	check.SetResult(compliantObjects, nonCompliantObjects)
}

// This function checks if the operator CRD is defined with OpenAPI 3 specification
func testOperatorCrdOpenAPISpec(check *checksdb.Check, env *provider.TestEnvironment) {
	check.LogInfo("Starting testOperatorCrdOpenAPISpec")
	var compliantObjects []*testhelper.ReportObject
	var nonCompliantObjects []*testhelper.ReportObject

	for _, crd := range env.Crds {
		isCrdDefinedWithOpenAPI3Schema := false

		for _, version := range crd.Spec.Versions {
			crdSchema := version.Schema.String()

			containsOpenAPIV3SchemaSubstr := strings.Contains(strings.ToLower(crdSchema),
				strings.ToLower(testhelper.OpenAPIV3Schema))

			if containsOpenAPIV3SchemaSubstr {
				isCrdDefinedWithOpenAPI3Schema = true
				break
			}
		}

		if isCrdDefinedWithOpenAPI3Schema {
			check.LogInfo("Operator CRD %s is defined with OpenAPIV3 schema ", crd.Name)
			compliantObjects = append(compliantObjects, testhelper.NewOperatorReportObject(crd.Namespace, crd.Name,
				"Operator CRD is defined with OpenAPIV3 schema ", true).AddField(testhelper.OpenAPIV3Schema, crd.Name))
		} else {
			check.LogInfo("Operator CRD %s is not defined with OpenAPIV3 schema ", crd.Name)
			nonCompliantObjects = append(nonCompliantObjects, testhelper.NewOperatorReportObject(crd.Namespace, crd.Name,
				"Operator CRD is not defined with OpenAPIV3 schema ", false).AddField(testhelper.OpenAPIV3Schema, crd.Name))
		}
	}

	check.SetResult(compliantObjects, nonCompliantObjects)
}

// This function checks for semantic versioning of the installed operators
func testOperatorSemanticVersioning(check *checksdb.Check, env *provider.TestEnvironment) {
	check.LogInfo("Starting testOperatorSemanticVersioning")
	var compliantObjects []*testhelper.ReportObject
	var nonCompliantObjects []*testhelper.ReportObject

	for _, operator := range env.Operators {
		operatorVersion := operator.Version
		check.LogInfo("Testing Operator %q for version %s", operator, operatorVersion)

		if versions.IsValidSemanticVersion(operatorVersion) {
			check.LogInfo("Operator %q has a valid semantic version %s", operator, operatorVersion)
			compliantObjects = append(compliantObjects, testhelper.NewOperatorReportObject(operator.Namespace, operator.Name,
				"Operator has a valid semantic version ", true).AddField(testhelper.Version, operatorVersion))
		} else {
			check.LogError("Operator %q has an invalid semantic version %s", operator, operatorVersion)
			nonCompliantObjects = append(nonCompliantObjects, testhelper.NewOperatorReportObject(operator.Namespace, operator.Name,
				"Operator has an invalid semantic version ", false).AddField(testhelper.Version, operatorVersion))
		}
	}

	check.SetResult(compliantObjects, nonCompliantObjects)
}

func testOperatorInstallationPhaseSucceeded(check *checksdb.Check, env *provider.TestEnvironment) {
	var compliantObjects []*testhelper.ReportObject
	var nonCompliantObjects []*testhelper.ReportObject
	for _, op := range env.Operators {
		check.LogInfo("Testing Operator %q", op)
		if phasecheck.WaitOperatorReady(op.Csv) {
			check.LogInfo("Operator %q is in Succeeded phase", op)
			compliantObjects = append(compliantObjects, testhelper.NewOperatorReportObject(op.Namespace, op.Name,
				"Operator on Succeeded state ", true).AddField(testhelper.OperatorPhase, string(op.Csv.Status.Phase)))
		} else {
			check.LogError("Operator %q is not in Succeeded phase (phase=%q)", op, op.Csv.Status.Phase)
			nonCompliantObjects = append(nonCompliantObjects, testhelper.NewOperatorReportObject(op.Namespace, op.Name,
				"Operator not in Succeeded state ", false).AddField(testhelper.OperatorPhase, string(op.Csv.Status.Phase)))
		}
	}

	check.SetResult(compliantObjects, nonCompliantObjects)
}

func testOperatorInstallationAccessToSCC(check *checksdb.Check, env *provider.TestEnvironment) {
	var compliantObjects []*testhelper.ReportObject
	var nonCompliantObjects []*testhelper.ReportObject
	for i := range env.Operators {
		operator := env.Operators[i]
		csv := operator.Csv
		check.LogDebug("Checking operator %s", operator)
		clusterPermissions := csv.Spec.InstallStrategy.StrategySpec.ClusterPermissions
		if len(clusterPermissions) == 0 {
			check.LogInfo("No clusterPermissions found in %s's CSV", operator)
			compliantObjects = append(compliantObjects, testhelper.NewOperatorReportObject(operator.Namespace, operator.Name,
				"No RBAC rules for Security Context Constraints found in CSV (no clusterPermissions found)", true))
			continue
		}

		// Fails in case any cluster permission has a rule that refers to securitycontextconstraints.
		badRuleFound := false
		for permissionIndex := range clusterPermissions {
			permission := &clusterPermissions[permissionIndex]
			for ruleIndex := range permission.Rules {
				rule := &permission.Rules[ruleIndex]

				// Check whether the rule is for the security api group.
				securityGroupFound := false
				for _, group := range rule.APIGroups {
					if group == "*" || group == "security.openshift.io" {
						securityGroupFound = true
						break
					}
				}

				if !securityGroupFound {
					continue
				}

				// Now check whether it grants some access to securitycontextconstraint resources.
				for _, resource := range rule.Resources {
					if resource == "*" || resource == "securitycontextconstraints" {
						check.LogInfo("Operator %s has a rule (index %d) for service account %s to access cluster SCCs",
							operator, ruleIndex, permission.ServiceAccountName)
						// Keep reviewing other permissions' rules so we can log all the failing ones in the claim file.
						badRuleFound = true
						break
					}
				}
			}
		}

		if badRuleFound {
			nonCompliantObjects = append(nonCompliantObjects, testhelper.NewOperatorReportObject(operator.Namespace, operator.Name, "One or more RBAC rules for Security Context Constraints found in CSV", false))
		} else {
			compliantObjects = append(compliantObjects, testhelper.NewOperatorReportObject(operator.Namespace, operator.Name, "No RBAC rules for Security Context Constraints found in CSV", true))
		}
	}

	check.SetResult(compliantObjects, nonCompliantObjects)
}

func testOperatorOlmSubscription(check *checksdb.Check, env *provider.TestEnvironment) {
	var compliantObjects []*testhelper.ReportObject
	var nonCompliantObjects []*testhelper.ReportObject

	for i := range env.Operators {
		operator := env.Operators[i]
		check.LogInfo("Testing Operator %q", operator)
		if operator.SubscriptionName == "" {
			check.LogError("OLM subscription not found for Operator %q", operator)
			nonCompliantObjects = append(nonCompliantObjects, testhelper.NewOperatorReportObject(env.Operators[i].Namespace, env.Operators[i].Name, "OLM subscription not found for operator, so it is not installed via OLM", false).
				AddField(testhelper.SubscriptionName, operator.SubscriptionName))
		} else {
			check.LogInfo("OLM subscription %q found for Operator %q", operator.SubscriptionName, operator)
			compliantObjects = append(compliantObjects, testhelper.NewOperatorReportObject(env.Operators[i].Namespace, env.Operators[i].Name, "install-status-no-privilege (subscription found)", true).
				AddField(testhelper.SubscriptionName, operator.SubscriptionName))
		}
	}

	check.SetResult(compliantObjects, nonCompliantObjects)
}

func testOperatorSingleCrdOwner(check *checksdb.Check, env *provider.TestEnvironment) {
	var compliantObjects []*testhelper.ReportObject
	var nonCompliantObjects []*testhelper.ReportObject

	// Map each CRD to a list of operators that own it
	crdOwners := map[string][]string{}
	for i := range env.Operators {
		operator := env.Operators[i]
		ownedCrds := operator.Csv.Spec.CustomResourceDefinitions.Owned
		for j := range ownedCrds {
			crdOwners[ownedCrds[j].Name] = append(crdOwners[ownedCrds[j].Name], operator.Name)
		}
	}

	// Flag those that are owned by more than one operator
	for crd, opList := range crdOwners {
		if len(opList) > 1 {
			check.LogError("CRD %q is owned by more than one operator (owners: %v)", crd, opList)
			nonCompliantObjects = append(nonCompliantObjects,
				testhelper.NewCrdReportObject(crd, "", "CRD is owned by more than one operator", false).
					AddField(testhelper.OperatorList, strings.Join(opList, ", ")))
		} else {
			check.LogDebug("CRD %q is owned by a single operator (%v)", crd, opList[0])
			compliantObjects = append(compliantObjects,
				testhelper.NewCrdReportObject(crd, "", "CRD is owned by a single operator", true).
					AddField(testhelper.OperatorName, opList[0]))
		}
	}

	check.SetResult(compliantObjects, nonCompliantObjects)
}

func testOperatorPodsRunAsUserID(check *checksdb.Check, env *provider.TestEnvironment) {
	var compliantObjects []*testhelper.ReportObject
	var nonCompliantObjects []*testhelper.ReportObject

	for csv, pods := range env.CSVToPodListMap {
		check.LogInfo("Csv: %q", csv)
		for _, pod := range pods {
			check.LogInfo("Testing Pod %q", pod)
			if pod.IsRunAsUserID(0) {
				check.LogError("Pod %q UserID of the pods created by the operator is 0", pod.Name)
				nonCompliantObjects = append(nonCompliantObjects, testhelper.NewPodReportObject(pod.Namespace, pod.Name, "Pod has been found with UserID is 0", false))
			} else {
				check.LogInfo("Pod %q UserID of the pods created by the operator is not 0", pod.Name)
				compliantObjects = append(compliantObjects, testhelper.NewPodReportObject(pod.Namespace, pod.Name, "Pod has been found with UserID is not 0", true))
			}
		}
	}
	check.SetResult(compliantObjects, nonCompliantObjects)
}

func testOperatorPodsRunAsNonRoot(check *checksdb.Check, env *provider.TestEnvironment) {
	var compliantObjects []*testhelper.ReportObject
	var nonCompliantObjects []*testhelper.ReportObject

	for csv, pods := range env.CSVToPodListMap {
		check.LogInfo("Csv: %q", csv)
		for _, pod := range pods {
			// We are looking through both the containers and the pods separately to make compliant and non-compliant objects.

			for _, c := range pod.Containers {
				if c.IsContainerRunAsNonRoot() {
					check.LogInfo("Container %q created by the operator is run as not root", c.Name)
					compliantObjects = append(compliantObjects, testhelper.NewPodReportObject(c.Namespace, c.Name, "Container has been found is run as not root", true))
				} else {
					check.LogError("Pod %q created by the operator is run as root", pod.Name)
					nonCompliantObjects = append(nonCompliantObjects, testhelper.NewPodReportObject(pod.Namespace, pod.Name, "Container has been found is run as root", false))
				}
			}

			if pod.IsRunAsNonRoot() {
				check.LogInfo("Pod %q created by the operator is run as not root", pod.Name)
				compliantObjects = append(compliantObjects, testhelper.NewPodReportObject(pod.Namespace, pod.Name, "Pod has been found is run as not root", true))
			} else {
				check.LogError("Pod %q created by the operator is run as root", pod.Name)
				nonCompliantObjects = append(nonCompliantObjects, testhelper.NewPodReportObject(pod.Namespace, pod.Name, "Pod has been found is run as root", false))
			}
		}
	}
	check.SetResult(compliantObjects, nonCompliantObjects)
}

func testOperatorPodsAutomountTokens(check *checksdb.Check, env *provider.TestEnvironment) {
	var compliantObjects []*testhelper.ReportObject
	var nonCompliantObjects []*testhelper.ReportObject

	for csv, pods := range env.CSVToPodListMap {
		check.LogInfo("Csv: %q", csv)
		for _, pod := range pods {
			check.LogInfo("Testing Pod %q", pod)
			// Evaluate the pod's automount service tokens and any attached service accounts
			client := clientsholder.GetClientsHolder()
			podPassed, newMsg := rbac.EvaluateAutomountTokens(client.K8sClient.CoreV1(), pod.Pod)
			if !podPassed {
				check.LogInfo("Pod %q have automount service tokens set to false", pod)
				compliantObjects = append(compliantObjects, testhelper.NewPodReportObject(pod.Namespace, pod.Name, "Pod %q created by the operator that the automount service account token set to false", true))
			} else {
				check.LogError(newMsg)
				nonCompliantObjects = append(nonCompliantObjects, testhelper.NewPodReportObject(pod.Namespace, pod.Name, newMsg, false))
			}
		}
	}
	check.SetResult(compliantObjects, nonCompliantObjects)
}

func testOperatorContainersReadOnlyFilesystem(check *checksdb.Check, env *provider.TestEnvironment) {
	var compliantObjects []*testhelper.ReportObject
	var nonCompliantObjects []*testhelper.ReportObject

	for csv, pods := range env.CSVToPodListMap {
		check.LogInfo("Csv: %q", csv)
		for _, pod := range pods {
			check.LogInfo("Testing Pod %q", pod)
			for _, cut := range pod.Containers {
				check.LogInfo("Testing Container %q", cut.Name)
				if cut.IsReadOnlyRootFilesystem(check.GetLogger()) {
					check.LogInfo("Pod %q container %q created by the operator is read only root file system.", pod.Name, cut.Name)
					compliantObjects = append(compliantObjects, testhelper.NewPodReportObject(pod.Namespace, pod.Name, "Pod has been found read only root file system", true))
				} else {
					check.LogError("Pod %q container %q created by the operator is read not only root file system.", pod.Name, cut.Name)
					nonCompliantObjects = append(nonCompliantObjects, testhelper.NewPodReportObject(pod.Namespace, pod.Name, "Pod has been found read not only root file system", false))
				}
			}
		}
		check.SetResult(compliantObjects, nonCompliantObjects)
	}
}
