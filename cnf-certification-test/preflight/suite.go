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

package preflight

import (
	"fmt"
	"strings"

	plibRuntime "github.com/redhat-openshift-ecosystem/openshift-preflight/certification"
	"github.com/sirupsen/logrus"
	"github.com/test-network-function/cnf-certification-test/cnf-certification-test/common"
	"github.com/test-network-function/cnf-certification-test/cnf-certification-test/identifiers"
	"github.com/test-network-function/cnf-certification-test/pkg/checksdb"
	"github.com/test-network-function/cnf-certification-test/pkg/configuration"
	"github.com/test-network-function/cnf-certification-test/pkg/provider"
	"github.com/test-network-function/cnf-certification-test/pkg/testhelper"
	"github.com/test-network-function/cnf-certification-test/pkg/tnf"
)

var (
	env provider.TestEnvironment

	beforeEachFn = func(check *checksdb.Check) error {
		logrus.Infof("Check %s: getting test environment.", check.ID)
		env = provider.GetTestEnvironment()
		return nil
	}
)

func labelsAllowTestRun(labelFilter string, allowedLabels []string) bool {
	for _, label := range allowedLabels {
		if strings.Contains(labelFilter, label) {
			return true
		}
	}
	return false
}

// Returns true if the preflight checks should run.
// Conditions: (1) the labels expr should contain any of the preflight tags/labels & (2) the
// preflight dockerconfig file must exist.
// This is just a hack to avoid running the preflight.LoadChecks() if it's not necessary
// since that function is actually running all the preflight lib's checks, which can take some
// time to finish. When they're finished, a checksdb.Check is created for each preflight lib's
// check that has run. The CheckFn will simply store the result.
func ShouldRun(labelsExpr string) bool {
	env = provider.GetTestEnvironment()
	preflightAllowedLabels := []string{common.PreflightTestKey, identifiers.TagPreflight}

	if !labelsAllowTestRun(labelsExpr, preflightAllowedLabels) {
		return false
	}

	// Add safeguard against running the preflight tests if the docker config does not exist.
	preflightDockerConfigFile := configuration.GetTestParameters().PfltDockerconfig
	if preflightDockerConfigFile == "" || preflightDockerConfigFile == "NA" {
		logrus.Warn("Skipping the preflight suite because the Docker Config file is not provided.")
		env.SkipPreflight = true
	}

	return true
}

func LoadChecks() {
	logrus.Debugf("Entering %s suite", common.PreflightTestKey)

	// As the preflight lib's checks need to run here, we need to get the test environment now.
	env = provider.GetTestEnvironment()

	checksGroup := checksdb.NewChecksGroup(common.PreflightTestKey).
		WithBeforeEachFn(beforeEachFn)

	testPreflightContainers(checksGroup, &env)
	if provider.IsOCPCluster() {
		logrus.Debugf("OCP cluster detected, allowing operator tests to run")
		testPreflightOperators(checksGroup, &env)
	} else {
		logrus.Debugf("Skipping the preflight operators test because it requires an OCP cluster to run against")
	}
}

func testPreflightOperators(checksGroup *checksdb.ChecksGroup, env *provider.TestEnvironment) {
	// Loop through all of the operators, run preflight, and set their results into their respective object
	for _, op := range env.Operators {
		// Note: We are not using a cache here for the operator bundle images because
		// in-general you are only going to have an operator installed once in a cluster.
		err := op.SetPreflightResults(env)
		if err != nil {
			logrus.Fatalf("failed running preflight on operator: %s error: %v", op.Name, err)
		}
	}

	logrus.Infof("Completed running preflight operator tests for %d operators", len(env.Operators))

	// Handle Operator-based preflight tests
	// Note: We only care about the `testEntry` variable below because we need its 'Description' and 'Suggestion' variables.
	for testName, testEntry := range getUniqueTestEntriesFromOperatorResults(env.Operators) {
		logrus.Infof("Testing operator ginkgo test: %s", testName)
		generatePreflightOperatorGinkgoTest(checksGroup, testName, testEntry.Metadata().Description, testEntry.Help().Suggestion, env.Operators)
	}
}

func testPreflightContainers(checksGroup *checksdb.ChecksGroup, env *provider.TestEnvironment) {
	// Using a cache to prevent unnecessary processing of images if we already have the results available
	preflightImageCache := make(map[string]plibRuntime.Results)

	// Loop through all of the containers, run preflight, and set their results into their respective objects
	for _, cut := range env.Containers {
		logrus.Debugf("Running preflight container tests for: %s", cut.Name)
		err := cut.SetPreflightResults(preflightImageCache, env)
		if err != nil {
			logrus.Fatalf("failed running preflight on image: %s error: %v", cut.Image, err)
		}
	}

	logrus.Infof("Completed running preflight container tests for %d containers", len(env.Containers))

	// Handle Container-based preflight tests
	// Note: We only care about the `testEntry` variable below because we need its 'Description' and 'Suggestion' variables.
	for testName, testEntry := range getUniqueTestEntriesFromContainerResults(env.Containers) {
		logrus.Infof("Testing container ginkgo test: %s", testName)
		generatePreflightContainerGinkgoTest(checksGroup, testName, testEntry.Metadata().Description, testEntry.Help().Suggestion, env.Containers)
	}
}

// func generatePreflightContainerGinkgoTest(testName, testID string, tags []string, containers []*provider.Container) {
func generatePreflightContainerGinkgoTest(checksGroup *checksdb.ChecksGroup, testName, description, suggestion string, containers []*provider.Container) {
	// Based on a single test "name", we will be passing/failing in Ginkgo.
	// Brute force-ish type of method.

	// Store the test names into the Catalog map for results to be dynamically printed
	aID := identifiers.AddCatalogEntry(testName, common.PreflightTestKey, description, suggestion, "", "", false, map[string]string{
		identifiers.FarEdge:  identifiers.Optional,
		identifiers.Telco:    identifiers.Optional,
		identifiers.NonTelco: identifiers.Optional,
		identifiers.Extended: identifiers.Optional,
	}, identifiers.TagPreflight)

	testID, tags := identifiers.GetGinkgoTestIDAndLabels(aID)
	check := checksdb.NewCheck(testID, tags).
		WithSkipCheckFn(testhelper.GetNoContainersUnderTestSkipFn(&env)).
		WithCheckFn(func(c *checksdb.Check) error {
			var compliantObjects []*testhelper.ReportObject
			var nonCompliantObjects []*testhelper.ReportObject
			for _, cut := range containers {
				for _, r := range cut.PreflightResults.Passed {
					if r.Name() == testName {
						compliantObjects = append(compliantObjects, testhelper.NewContainerReportObject(cut.Namespace, cut.Podname, cut.Name, "Container has passed preflight test "+testName, true))
						logrus.Infof("%s has passed preflight test: %s", cut.String(), testName)
					}
				}
				for _, r := range cut.PreflightResults.Failed {
					if r.Name() == testName {
						nonCompliantObjects = append(nonCompliantObjects, testhelper.NewContainerReportObject(cut.Namespace, cut.Podname, cut.Name, "Container has failed preflight test "+testName, false))
						tnf.Logf(logrus.WarnLevel, "%s has failed preflight test: %s", cut, testName)
					}
				}
				for _, r := range cut.PreflightResults.Errors {
					if r.Name() == testName {
						nonCompliantObjects = append(nonCompliantObjects, testhelper.NewContainerReportObject(cut.Namespace, cut.Podname, cut.Name, fmt.Sprintf("Container has errored preflight test %s, err=%v", testName, r.Err), false))
						tnf.Logf(logrus.ErrorLevel, "%s has errored preflight test: %s", cut, testName)
					}
				}
			}

			c.SetResult(compliantObjects, nonCompliantObjects)
			return nil
		})

	checksGroup.Add(check)
}

func generatePreflightOperatorGinkgoTest(checksGroup *checksdb.ChecksGroup, testName, description, suggestion string, operators []*provider.Operator) {
	// Based on a single test "name", we will be passing/failing in Ginkgo.
	// Brute force-ish type of method.

	// Store the test names into the Catalog map for results to be dynamically printed
	aID := identifiers.AddCatalogEntry(testName, common.PreflightTestKey, description, suggestion, "", "", false, map[string]string{
		identifiers.FarEdge:  identifiers.Optional,
		identifiers.Telco:    identifiers.Optional,
		identifiers.NonTelco: identifiers.Optional,
		identifiers.Extended: identifiers.Optional,
	}, identifiers.TagPreflight)

	testID, tags := identifiers.GetGinkgoTestIDAndLabels(aID)
	check := checksdb.NewCheck(testID, tags).
		WithSkipCheckFn(testhelper.GetNoContainersUnderTestSkipFn(&env)).
		WithCheckFn(func(c *checksdb.Check) error {
			var compliantObjects []*testhelper.ReportObject
			var nonCompliantObjects []*testhelper.ReportObject

			for _, op := range operators {
				for _, r := range op.PreflightResults.Passed {
					if r.Name() == testName {
						logrus.Infof("%s has passed preflight test: %s", op.String(), testName)
						compliantObjects = append(compliantObjects, testhelper.NewOperatorReportObject(op.Namespace, op.Name, "Operator passed preflight test "+testName, true))
					}
				}
				for _, r := range op.PreflightResults.Failed {
					if r.Name() == testName {
						tnf.Logf(logrus.WarnLevel, "%s has failed preflight test: %s", op, testName)
						nonCompliantObjects = append(nonCompliantObjects, testhelper.NewOperatorReportObject(op.Namespace, op.Name, "Operator failed preflight test "+testName, false))
					}
				}
				for _, r := range op.PreflightResults.Errors {
					if r.Name() == testName {
						tnf.Logf(logrus.ErrorLevel, "%s has errored preflight test: %s", op, testName)
						nonCompliantObjects = append(nonCompliantObjects, testhelper.NewOperatorReportObject(op.Namespace, op.Name, "Operator has errored preflight test "+testName, false))
					}
				}
			}

			c.SetResult(compliantObjects, nonCompliantObjects)
			return nil
		})

	checksGroup.Add(check)
}

func getUniqueTestEntriesFromContainerResults(containers []*provider.Container) map[string]plibRuntime.Result {
	// If containers are sharing the same image, they should "presumably" have the same results returned from preflight.
	testEntries := make(map[string]plibRuntime.Result)
	for _, cut := range containers {
		for _, r := range cut.PreflightResults.Passed {
			testEntries[r.Name()] = r
		}
		// Failed Results have more information than the rest
		for _, r := range cut.PreflightResults.Failed {
			testEntries[r.Name()] = r
		}
		for _, r := range cut.PreflightResults.Errors {
			testEntries[r.Name()] = r
		}
	}

	return testEntries
}

func getUniqueTestEntriesFromOperatorResults(operators []*provider.Operator) map[string]plibRuntime.Result {
	testEntries := make(map[string]plibRuntime.Result)
	for _, op := range operators {
		for _, r := range op.PreflightResults.Passed {
			testEntries[r.Name()] = r
		}
		// Failed Results have more information than the rest
		for _, r := range op.PreflightResults.Failed {
			testEntries[r.Name()] = r
		}
		for _, r := range op.PreflightResults.Errors {
			testEntries[r.Name()] = r
		}
	}
	return testEntries
}
