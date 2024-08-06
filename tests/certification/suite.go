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

package certification

import (
	"context"
	"fmt"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/test-network-function/cnf-certification-test/tests/common"
	"github.com/test-network-function/cnf-certification-test/tests/identifiers"

	"github.com/test-network-function/cnf-certification-test/internal/clientsholder"
	"github.com/test-network-function/cnf-certification-test/internal/log"
	"github.com/test-network-function/cnf-certification-test/pkg/checksdb"
	"github.com/test-network-function/cnf-certification-test/pkg/provider"
	"github.com/test-network-function/cnf-certification-test/pkg/testhelper"
	"github.com/test-network-function/oct/pkg/certdb"
)

const (
	// timeout for eventually call.
	CertifiedOperator = "certified-operators"
	Online            = "online"
)

var (
	env       provider.TestEnvironment
	validator certdb.CertificationStatusValidator

	beforeEachFn = func(check *checksdb.Check) error {
		env = provider.GetTestEnvironment()

		var err error
		validator, err = certdb.GetValidator(env.GetOfflineDBPath())
		if err != nil {
			return fmt.Errorf("cannot access the certification DB, err: %v", err)
		}

		return nil
	}

	skipIfNoOperatorsFn = func() (bool, string) {
		if len(env.Operators) == 0 {
			return true, "There are no operators to check. Please check under test labels."
		}

		return false, ""
	}

	skipIfNoHelmChartReleasesFn = func() (bool, string) {
		if len(env.HelmChartReleases) == 0 {
			return true, "There are no helm chart releases to check."
		}

		return false, ""
	}
)

func LoadChecks() {
	log.Debug("Loading %s suite checks", common.AffiliatedCertTestKey)

	checksGroup := checksdb.NewChecksGroup(common.AffiliatedCertTestKey).
		WithBeforeEachFn(beforeEachFn)

	checksGroup.Add(checksdb.NewCheck(identifiers.GetTestIDAndLabels(identifiers.TestHelmVersionIdentifier)).
		WithSkipCheckFn(skipIfNoHelmChartReleasesFn).
		WithCheckFn(func(check *checksdb.Check) error {
			testHelmVersion(check)
			return nil
		}))

	checksGroup.Add(checksdb.NewCheck(identifiers.GetTestIDAndLabels(identifiers.TestOperatorIsCertifiedIdentifier)).
		WithSkipCheckFn(skipIfNoOperatorsFn).
		WithCheckFn(func(c *checksdb.Check) error {
			testAllOperatorCertified(c, &env, validator)
			return nil
		}))

	checksGroup.Add(checksdb.NewCheck(identifiers.GetTestIDAndLabels(identifiers.TestHelmIsCertifiedIdentifier)).
		WithSkipCheckFn(skipIfNoHelmChartReleasesFn).
		WithCheckFn(func(c *checksdb.Check) error {
			testHelmCertified(c, &env, validator)
			return nil
		}))

	checksGroup.Add(checksdb.NewCheck(identifiers.GetTestIDAndLabels(identifiers.TestContainerIsCertifiedDigestIdentifier)).
		WithSkipCheckFn(testhelper.GetNoContainersUnderTestSkipFn(&env)).
		WithCheckFn(func(c *checksdb.Check) error {
			testContainerCertificationStatusByDigest(c, &env, validator)
			return nil
		}))
}

func getContainersToQuery(env *provider.TestEnvironment) map[provider.ContainerImageIdentifier]bool {
	containersToQuery := make(map[provider.ContainerImageIdentifier]bool)
	for _, cut := range env.Containers {
		containersToQuery[cut.ContainerImageIdentifier] = true
	}
	return containersToQuery
}

func testContainerCertification(c provider.ContainerImageIdentifier, validator certdb.CertificationStatusValidator) bool {
	return validator.IsContainerCertified(c.Registry, c.Repository, c.Tag, c.Digest)
}

func testAllOperatorCertified(check *checksdb.Check, env *provider.TestEnvironment, validator certdb.CertificationStatusValidator) {
	operatorsUnderTest := env.Operators
	var compliantObjects []*testhelper.ReportObject
	var nonCompliantObjects []*testhelper.ReportObject

	ocpMinorVersion := ""
	if provider.IsOCPCluster() {
		// Converts	major.minor.patch version format to major.minor
		const majorMinorPatchCount = 3
		splitVersion := strings.SplitN(env.OpenshiftVersion, ".", majorMinorPatchCount)
		ocpMinorVersion = splitVersion[0] + "." + splitVersion[1]
	}
	for _, operator := range operatorsUnderTest {
		check.LogInfo("Testing Operator %q", operator)
		isCertified := validator.IsOperatorCertified(operator.Name, ocpMinorVersion, operator.Channel)
		if !isCertified {
			check.LogError("Operator %q (channel %q) failed to be certified for OpenShift %s", operator.Name, operator.Channel, ocpMinorVersion)
			nonCompliantObjects = append(nonCompliantObjects, testhelper.NewOperatorReportObject(operator.Namespace, operator.Name, "Operator failed to be certified for OpenShift", false).
				AddField(testhelper.OCPVersion, ocpMinorVersion).
				AddField(testhelper.OCPChannel, operator.Channel))
		} else {
			check.LogInfo("Operator %q (channel %q) is certified for OpenShift %s", operator.Name, operator.Channel, ocpMinorVersion)
			compliantObjects = append(compliantObjects, testhelper.NewOperatorReportObject(operator.Namespace, operator.Name, "Operator certified OK", true).
				AddField(testhelper.OCPVersion, ocpMinorVersion).
				AddField(testhelper.OCPChannel, operator.Channel))
		}
	}

	check.SetResult(compliantObjects, nonCompliantObjects)
}

func testHelmCertified(check *checksdb.Check, env *provider.TestEnvironment, validator certdb.CertificationStatusValidator) {
	helmchartsReleases := env.HelmChartReleases

	// Collect all of the failed helm charts
	var compliantObjects []*testhelper.ReportObject
	var nonCompliantObjects []*testhelper.ReportObject
	for _, helm := range helmchartsReleases {
		check.LogInfo("Testing Helm Chart Release %q", helm.Name)
		if !validator.IsHelmChartCertified(helm, env.K8sVersion) {
			check.LogError("Helm Chart %q version %q is not certified.", helm.Name, helm.Chart.Metadata.Version)
			nonCompliantObjects = append(nonCompliantObjects, testhelper.NewHelmChartReportObject(helm.Namespace, helm.Name, "helm chart is not certified", false).
				SetType(testhelper.HelmVersionType).
				AddField(testhelper.Version, helm.Chart.Metadata.Version))
		} else {
			check.LogInfo("Helm Chart %q version %q is certified.", helm.Name, helm.Chart.Metadata.Version)
			compliantObjects = append(compliantObjects, testhelper.NewHelmChartReportObject(helm.Namespace, helm.Name, "helm chart is certified", true).
				SetType(testhelper.HelmVersionType).
				AddField(testhelper.Version, helm.Chart.Metadata.Version))
		}
	}

	check.SetResult(compliantObjects, nonCompliantObjects)
}

func testContainerCertificationStatusByDigest(check *checksdb.Check, env *provider.TestEnvironment, validator certdb.CertificationStatusValidator) {
	var compliantObjects []*testhelper.ReportObject
	var nonCompliantObjects []*testhelper.ReportObject
	for _, c := range env.Containers {
		check.LogInfo("Testing Container %q", c)
		switch {
		case c.ContainerImageIdentifier.Digest == "":
			check.LogError("Container %q is missing digest field, failing validation (repo=%q image=%q)", c, c.ContainerImageIdentifier.Registry, c.ContainerImageIdentifier.Repository)
			nonCompliantObjects = append(nonCompliantObjects, testhelper.NewContainerReportObject(c.Namespace, c.Podname, c.Name, "Missing digest field", false).
				AddField(testhelper.Repository, c.ContainerImageIdentifier.Registry).
				AddField(testhelper.ImageName, c.ContainerImageIdentifier.Repository).
				AddField(testhelper.ImageDigest, c.ContainerImageIdentifier.Digest))
		case !testContainerCertification(c.ContainerImageIdentifier, validator):
			check.LogError("Container %q digest not found in database, failing validation (repo=%q image=%q tag=%q digest=%q)", c,
				c.ContainerImageIdentifier.Registry, c.ContainerImageIdentifier.Repository,
				c.ContainerImageIdentifier.Tag, c.ContainerImageIdentifier.Digest)
			nonCompliantObjects = append(nonCompliantObjects, testhelper.NewContainerReportObject(c.Namespace, c.Podname, c.Name, "Digest not found in database", false).
				AddField(testhelper.Repository, c.ContainerImageIdentifier.Registry).
				AddField(testhelper.ImageName, c.ContainerImageIdentifier.Repository).
				AddField(testhelper.ImageDigest, c.ContainerImageIdentifier.Digest))
		default:
			check.LogInfo("Container %q digest found in database, image certified (repo=%q image=%q tag=%q digest=%q)", c,
				c.ContainerImageIdentifier.Registry, c.ContainerImageIdentifier.Repository,
				c.ContainerImageIdentifier.Tag, c.ContainerImageIdentifier.Digest)
			compliantObjects = append(compliantObjects, testhelper.NewContainerReportObject(c.Namespace, c.Podname, c.Name, "Container is certified", true))
		}
	}

	check.SetResult(compliantObjects, nonCompliantObjects)
}

func testHelmVersion(check *checksdb.Check) {
	var compliantObjects []*testhelper.ReportObject
	var nonCompliantObjects []*testhelper.ReportObject

	clients := clientsholder.GetClientsHolder()
	// Get the Tiller pod in the specified namespace
	podList, err := clients.K8sClient.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{
		LabelSelector: "app=helm,name=tiller",
	})
	if err != nil {
		check.LogError("Could not get Tiller pod, err=%v", err)
	}

	if len(podList.Items) == 0 {
		check.LogInfo("Tiller pod not found in any namespaces. Helm version is v3.")
		for _, helm := range env.HelmChartReleases {
			compliantObjects = append(compliantObjects, testhelper.NewHelmChartReportObject(helm.Namespace, helm.Name, "helm chart was installed with helm v3", true))
		}
	}

	check.LogError("Tiller pod found, Helm version is v2 but v3 required")
	for i := range podList.Items {
		nonCompliantObjects = append(nonCompliantObjects, testhelper.NewPodReportObject(podList.Items[i].Namespace, podList.Items[i].Name,
			"This pod is a Tiller pod. Helm Chart version is v2 but needs to be v3 due to the security risks associated with Tiller", false))
	}

	check.SetResult(compliantObjects, nonCompliantObjects)
}
