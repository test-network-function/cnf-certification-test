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

package accesscontrol

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/operator-framework/api/pkg/operators/v1alpha1"
	"github.com/test-network-function/cnf-certification-test/cnf-certification-test/accesscontrol/namespace"
	"github.com/test-network-function/cnf-certification-test/cnf-certification-test/accesscontrol/resources"
	"github.com/test-network-function/cnf-certification-test/cnf-certification-test/accesscontrol/securitycontextcontainer"
	"github.com/test-network-function/cnf-certification-test/cnf-certification-test/common"
	"github.com/test-network-function/cnf-certification-test/cnf-certification-test/common/rbac"
	"github.com/test-network-function/cnf-certification-test/cnf-certification-test/identifiers"
	"github.com/test-network-function/cnf-certification-test/cnf-certification-test/networking/netutil"
	"github.com/test-network-function/cnf-certification-test/cnf-certification-test/networking/services"
	"github.com/test-network-function/cnf-certification-test/internal/clientsholder"
	"github.com/test-network-function/cnf-certification-test/internal/crclient"
	"github.com/test-network-function/cnf-certification-test/internal/log"
	"github.com/test-network-function/cnf-certification-test/pkg/checksdb"
	"github.com/test-network-function/cnf-certification-test/pkg/podhelper"
	"github.com/test-network-function/cnf-certification-test/pkg/provider"
	"github.com/test-network-function/cnf-certification-test/pkg/stringhelper"
	"github.com/test-network-function/cnf-certification-test/pkg/testhelper"
	rbacv1 "k8s.io/api/rbac/v1"
)

const (
	nodePort              = "NodePort"
	defaultServiceAccount = "default"
)

var (
	invalidNamespacePrefixes = []string{
		"default",
		"openshift-",
		"istio-",
		"aspenmesh-",
	}
)

var (
	env provider.TestEnvironment

	beforeEachFn = func(check *checksdb.Check) error {
		env = provider.GetTestEnvironment()
		return nil
	}
)

// LoadChecks loads all the checks.
//
//nolint:funlen
func LoadChecks() {
	log.Debug("Loading %s suite checks", common.AccessControlTestKey)

	checksGroup := checksdb.NewChecksGroup(common.AccessControlTestKey).
		WithBeforeEachFn(beforeEachFn)

	checksGroup.Add(checksdb.NewCheck(identifiers.GetTestIDAndLabels(identifiers.TestSecContextIdentifier)).
		WithSkipCheckFn(testhelper.GetNoContainersUnderTestSkipFn(&env)).
		WithCheckFn(func(c *checksdb.Check) error {
			testContainerSCC(c, &env)
			return nil
		}))

	checksGroup.Add(checksdb.NewCheck(identifiers.GetTestIDAndLabels(identifiers.TestSysAdminIdentifier)).
		WithSkipCheckFn(testhelper.GetNoContainersUnderTestSkipFn(&env)).
		WithCheckFn(func(c *checksdb.Check) error {
			testSysAdminCapability(c, &env)
			return nil
		}))

	checksGroup.Add(checksdb.NewCheck(identifiers.GetTestIDAndLabels(identifiers.TestNetAdminIdentifier)).
		WithSkipCheckFn(testhelper.GetNoContainersUnderTestSkipFn(&env)).
		WithCheckFn(func(c *checksdb.Check) error {
			testNetAdminCapability(c, &env)
			return nil
		}))

	checksGroup.Add(checksdb.NewCheck(identifiers.GetTestIDAndLabels(identifiers.TestNetRawIdentifier)).
		WithSkipCheckFn(testhelper.GetNoContainersUnderTestSkipFn(&env)).
		WithCheckFn(func(c *checksdb.Check) error {
			testNetRawCapability(c, &env)
			return nil
		}))

	checksGroup.Add(checksdb.NewCheck(identifiers.GetTestIDAndLabels(identifiers.TestIpcLockIdentifier)).
		WithSkipCheckFn(testhelper.GetNoContainersUnderTestSkipFn(&env)).
		WithCheckFn(func(c *checksdb.Check) error {
			testIpcLockCapability(c, &env)
			return nil
		}))

	checksGroup.Add(checksdb.NewCheck(identifiers.GetTestIDAndLabels(identifiers.TestBpfIdentifier)).
		WithSkipCheckFn(testhelper.GetNoContainersUnderTestSkipFn(&env)).
		WithCheckFn(func(c *checksdb.Check) error {
			testBpfCapability(c, &env)
			return nil
		}))

	checksGroup.Add(checksdb.NewCheck(identifiers.GetTestIDAndLabels(identifiers.TestSecConNonRootUserIdentifier)).
		WithSkipCheckFn(testhelper.GetNoContainersUnderTestSkipFn(&env)).
		WithCheckFn(func(c *checksdb.Check) error {
			testSecConRootUser(c, &env)
			return nil
		}))

	checksGroup.Add(checksdb.NewCheck(identifiers.GetTestIDAndLabels(identifiers.TestSecConPrivilegeEscalation)).
		WithSkipCheckFn(testhelper.GetNoContainersUnderTestSkipFn(&env)).
		WithCheckFn(func(c *checksdb.Check) error {
			testSecConPrivilegeEscalation(c, &env)
			return nil
		}))

	checksGroup.Add(checksdb.NewCheck(identifiers.GetTestIDAndLabels(identifiers.TestContainerHostPort)).
		WithSkipCheckFn(testhelper.GetNoContainersUnderTestSkipFn(&env)).
		WithCheckFn(func(c *checksdb.Check) error {
			testContainerHostPort(c, &env)
			return nil
		}))

	checksGroup.Add(checksdb.NewCheck(identifiers.GetTestIDAndLabels(identifiers.TestPodHostNetwork)).
		WithSkipCheckFn(testhelper.GetNoPodsUnderTestSkipFn(&env)).
		WithCheckFn(func(c *checksdb.Check) error {
			testPodHostNetwork(c, &env)
			return nil
		}))

	checksGroup.Add(checksdb.NewCheck(identifiers.GetTestIDAndLabels(identifiers.TestPodHostPath)).
		WithSkipCheckFn(testhelper.GetNoPodsUnderTestSkipFn(&env)).
		WithCheckFn(func(c *checksdb.Check) error {
			testPodHostPath(c, &env)
			return nil
		}))

	checksGroup.Add(checksdb.NewCheck(identifiers.GetTestIDAndLabels(identifiers.TestPodHostIPC)).
		WithSkipCheckFn(testhelper.GetNoPodsUnderTestSkipFn(&env)).
		WithCheckFn(func(c *checksdb.Check) error {
			testPodHostIPC(c, &env)
			return nil
		}))

	checksGroup.Add(checksdb.NewCheck(identifiers.GetTestIDAndLabels(identifiers.TestPodHostPID)).
		WithSkipCheckFn(testhelper.GetNoPodsUnderTestSkipFn(&env)).
		WithCheckFn(func(c *checksdb.Check) error {
			testPodHostPID(c, &env)
			return nil
		}))

	checksGroup.Add(checksdb.NewCheck(identifiers.GetTestIDAndLabels(identifiers.TestNamespaceBestPracticesIdentifier)).
		WithSkipCheckFn(testhelper.GetNoNamespacesSkipFn(&env)).
		WithCheckFn(func(c *checksdb.Check) error {
			testNamespace(c, &env)
			return nil
		}))

	checksGroup.Add(checksdb.NewCheck(identifiers.GetTestIDAndLabels(identifiers.TestPodServiceAccountBestPracticesIdentifier)).
		WithSkipCheckFn(testhelper.GetNoPodsUnderTestSkipFn(&env)).
		WithCheckFn(func(c *checksdb.Check) error {
			testPodServiceAccount(c, &env)
			return nil
		}))

	checksGroup.Add(checksdb.NewCheck(identifiers.GetTestIDAndLabels(identifiers.TestPodRoleBindingsBestPracticesIdentifier)).
		WithSkipCheckFn(testhelper.GetNoPodsUnderTestSkipFn(&env)).
		WithCheckFn(func(c *checksdb.Check) error {
			testPodRoleBindings(c, &env)
			return nil
		}))

	checksGroup.Add(checksdb.NewCheck(identifiers.GetTestIDAndLabels(identifiers.TestPodClusterRoleBindingsBestPracticesIdentifier)).
		WithSkipCheckFn(testhelper.GetNoPodsUnderTestSkipFn(&env)).
		WithCheckFn(func(c *checksdb.Check) error {
			testPodClusterRoleBindings(c, &env)
			return nil
		}))

	checksGroup.Add(checksdb.NewCheck(identifiers.GetTestIDAndLabels(identifiers.TestPodAutomountServiceAccountIdentifier)).
		WithSkipCheckFn(testhelper.GetNoPodsUnderTestSkipFn(&env)).
		WithCheckFn(func(c *checksdb.Check) error {
			testAutomountServiceToken(c, &env)
			return nil
		}))

	checksGroup.Add(checksdb.NewCheck(identifiers.GetTestIDAndLabels(identifiers.TestOneProcessPerContainerIdentifier)).
		WithSkipCheckFn(testhelper.GetNoContainersUnderTestSkipFn(&env), testhelper.GetDaemonSetFailedToSpawnSkipFn(&env)).
		WithCheckFn(func(c *checksdb.Check) error {
			testOneProcessPerContainer(c, &env)
			return nil
		}))

	checksGroup.Add(checksdb.NewCheck(identifiers.GetTestIDAndLabels(identifiers.TestSYSNiceRealtimeCapabilityIdentifier)).
		WithSkipCheckFn(testhelper.GetNoContainersUnderTestSkipFn(&env)).
		WithCheckFn(func(c *checksdb.Check) error {
			testSYSNiceRealtimeCapability(c, &env)
			return nil
		}))

	checksGroup.Add(checksdb.NewCheck(identifiers.GetTestIDAndLabels(identifiers.TestSysPtraceCapabilityIdentifier)).
		WithSkipCheckFn(testhelper.GetSharedProcessNamespacePodsSkipFn(&env)).
		WithCheckFn(func(c *checksdb.Check) error {
			testSysPtraceCapability(c, &env)
			return nil
		}))

	checksGroup.Add(checksdb.NewCheck(identifiers.GetTestIDAndLabels(identifiers.TestNamespaceResourceQuotaIdentifier)).
		WithSkipCheckFn(testhelper.GetNoPodsUnderTestSkipFn(&env)).
		WithCheckFn(func(c *checksdb.Check) error {
			testNamespaceResourceQuota(c, &env)
			return nil
		}))

	checksGroup.Add(checksdb.NewCheck(identifiers.GetTestIDAndLabels(identifiers.TestNoSSHDaemonsAllowedIdentifier)).
		WithSkipCheckFn(testhelper.GetDaemonSetFailedToSpawnSkipFn(&env), testhelper.GetNoContainersUnderTestSkipFn(&env)).
		WithCheckFn(func(c *checksdb.Check) error {
			testNoSSHDaemonsAllowed(c, &env)
			return nil
		}))

	checksGroup.Add(checksdb.NewCheck(identifiers.GetTestIDAndLabels(identifiers.TestPodRequestsAndLimitsIdentifier)).
		WithSkipCheckFn(testhelper.GetNoPodsUnderTestSkipFn(&env)).
		WithCheckFn(func(c *checksdb.Check) error {
			testPodRequestsAndLimits(c, &env)
			return nil
		}))

	checksGroup.Add(checksdb.NewCheck(identifiers.GetTestIDAndLabels(identifiers.Test1337UIDIdentifier)).
		WithSkipCheckFn(testhelper.GetNoPodsUnderTestSkipFn(&env)).
		WithCheckFn(func(c *checksdb.Check) error {
			test1337UIDs(c, &env)
			return nil
		}))

	checksGroup.Add(checksdb.NewCheck(identifiers.GetTestIDAndLabels(identifiers.TestServicesDoNotUseNodeportsIdentifier)).
		WithSkipCheckFn(testhelper.GetNoServicesUnderTestSkipFn(&env)).
		WithCheckFn(func(c *checksdb.Check) error {
			testNodePort(c, &env)
			return nil
		}))

	checksGroup.Add(checksdb.NewCheck(identifiers.GetTestIDAndLabels(identifiers.TestCrdRoleIdentifier)).
		WithSkipCheckFn(testhelper.GetNoCrdsUnderTestSkipFn(&env), testhelper.GetNoNamespacesSkipFn(&env), testhelper.GetNoRolesSkipFn(&env)).
		WithCheckFn(func(c *checksdb.Check) error {
			testCrdRoles(c, &env)
			return nil
		}))
}

// checkForbiddenCapability checks if containers use a forbidden capability.
// Returns:
//   - compliantObjects []*testhelper.ReportObject : Slice containing report objects for containers compliant with the capability restrictions.
//   - nonCompliantObjects []*testhelper.ReportObject : Slice containing report objects for containers not compliant with the capability restrictions.
func checkForbiddenCapability(containers []*provider.Container, capability string, logger *log.Logger) (compliantObjects, nonCompliantObjects []*testhelper.ReportObject) {
	for _, cut := range containers {
		logger.Info("Testing Container %q", cut)
		compliant := true

		switch {
		case cut.SecurityContext == nil:
		case cut.SecurityContext.Capabilities == nil:
		case strings.Contains(cut.SecurityContext.Capabilities.String(), capability):
			compliant = false
		}

		if compliant {
			logger.Info("Container %q does not use non-compliant capability %q", cut, capability)
			compliantObjects = append(compliantObjects, testhelper.NewContainerReportObject(cut.Namespace, cut.Podname, cut.Name, "No forbidden capability "+capability+" detected in container", true))
		} else {
			logger.Error("Non compliant %q capability detected in container %q. All container caps: %q", capability, cut, cut.SecurityContext.Capabilities)
			nonCompliantObjects = append(nonCompliantObjects, testhelper.NewContainerReportObject(cut.Namespace, cut.Podname, cut.Name, "Non compliant capability "+capability+" in container", false).AddField(testhelper.SCCCapability, capability))
		}
	}
	return compliantObjects, nonCompliantObjects
}

func testSysAdminCapability(check *checksdb.Check, env *provider.TestEnvironment) {
	compliantObjects, nonCompliantObjects := checkForbiddenCapability(env.Containers, "SYS_ADMIN", check.GetLogger())
	check.SetResult(compliantObjects, nonCompliantObjects)
}

func testNetAdminCapability(check *checksdb.Check, env *provider.TestEnvironment) {
	compliantObjects, nonCompliantObjects := checkForbiddenCapability(env.Containers, "NET_ADMIN", check.GetLogger())
	check.SetResult(compliantObjects, nonCompliantObjects)
}

func testNetRawCapability(check *checksdb.Check, env *provider.TestEnvironment) {
	compliantObjects, nonCompliantObjects := checkForbiddenCapability(env.Containers, "NET_RAW", check.GetLogger())
	check.SetResult(compliantObjects, nonCompliantObjects)
}

func testIpcLockCapability(check *checksdb.Check, env *provider.TestEnvironment) {
	compliantObjects, nonCompliantObjects := checkForbiddenCapability(env.Containers, "IPC_LOCK", check.GetLogger())
	check.SetResult(compliantObjects, nonCompliantObjects)
}

func testBpfCapability(check *checksdb.Check, env *provider.TestEnvironment) {
	compliantObjects, nonCompliantObjects := checkForbiddenCapability(env.Containers, "BPF", check.GetLogger())
	check.SetResult(compliantObjects, nonCompliantObjects)
}

// testSecConRootUser verifies that the container is not running as root
func testSecConRootUser(check *checksdb.Check, env *provider.TestEnvironment) {
	var compliantObjects []*testhelper.ReportObject
	var nonCompliantObjects []*testhelper.ReportObject
	for _, put := range env.Pods {
		check.LogInfo("Testing Pod %q", put)
		if put.IsRunAsUserID(0) {
			check.LogError("Root user detected (RunAsUser uid=0) in Pod %q", put)
			nonCompliantObjects = append(nonCompliantObjects, testhelper.NewPodReportObject(put.Namespace, put.Name, "Root User detected (RunAsUser uid=0)", false))
		} else {
			check.LogInfo("Non-root user detected in Pod %q", put)
			compliantObjects = append(compliantObjects, testhelper.NewPodReportObject(put.Namespace, put.Name, "Root User not detected (RunAsUser uid=0)", true))
		}

		for idx := range put.Spec.Containers {
			cut := &(put.Spec.Containers[idx])
			check.LogInfo("Testing Container %q", cut)
			// Check the container level RunAsUser parameter
			if cut.SecurityContext != nil && cut.SecurityContext.RunAsUser != nil {
				if *(cut.SecurityContext.RunAsUser) == 0 {
					check.LogError("Root user detected (RunAsUser uid=0) in Container %q (%q)", cut.Name, put)
					nonCompliantObjects = append(nonCompliantObjects, testhelper.NewContainerReportObject(put.Namespace, put.Name, cut.Name, "Root User detected (RunAsUser uid=0)", false))
				} else {
					check.LogInfo("Non-root user detected (RunAsUser uid=0) in Container %q (%q)", cut, put)
					compliantObjects = append(compliantObjects, testhelper.NewContainerReportObject(put.Namespace, put.Name, cut.Name, "Root User not detected (RunAsUser uid=0)", true))
				}
			}
		}
	}

	check.SetResult(compliantObjects, nonCompliantObjects)
}

// testSecConPrivilegeEscalation verifies that the container is not allowed privilege escalation
func testSecConPrivilegeEscalation(check *checksdb.Check, env *provider.TestEnvironment) {
	var compliantObjects []*testhelper.ReportObject
	var nonCompliantObjects []*testhelper.ReportObject
	for _, cut := range env.Containers {
		check.LogInfo("Testing Container %q", cut)
		privEscFound := false
		if cut.SecurityContext != nil && cut.SecurityContext.AllowPrivilegeEscalation != nil {
			if *(cut.SecurityContext.AllowPrivilegeEscalation) {
				check.LogError("AllowPrivilegeEscalation is set to true in Container %q.", cut)
				nonCompliantObjects = append(nonCompliantObjects, testhelper.NewContainerReportObject(cut.Namespace, cut.Podname, cut.Name, "AllowPrivilegeEscalation is set to true", false))
				privEscFound = true
			}
		}

		if !privEscFound {
			check.LogInfo("AllowPrivilegeEscalation is set to false in Container %q.", cut)
			compliantObjects = append(compliantObjects, testhelper.NewContainerReportObject(cut.Namespace, cut.Podname, cut.Name, "AllowPrivilegeEscalation is not set to true", true))
		}
	}

	check.SetResult(compliantObjects, nonCompliantObjects)
}

// testContainerHostPort tests that containers are not configured with host port privileges
func testContainerHostPort(check *checksdb.Check, env *provider.TestEnvironment) {
	var compliantObjects []*testhelper.ReportObject
	var nonCompliantObjects []*testhelper.ReportObject
	for _, cut := range env.Containers {
		check.LogInfo("Testing Container %q", cut)
		hostPortFound := false
		for _, aPort := range cut.Ports {
			if aPort.HostPort != 0 {
				check.LogError("Host port %d is configured in Container %q.", aPort.HostPort, cut)
				nonCompliantObjects = append(nonCompliantObjects, testhelper.NewContainerReportObject(cut.Namespace, cut.Podname, cut.Name, "Host port is configured", false).
					SetType(testhelper.HostPortType).
					AddField(testhelper.PortNumber, strconv.Itoa(int(aPort.HostPort))))
				hostPortFound = true
			}
		}

		if !hostPortFound {
			check.LogInfo("Host port not configured in Container %q.", cut)
			compliantObjects = append(compliantObjects, testhelper.NewContainerReportObject(cut.Namespace, cut.Podname, cut.Name, "Host port is not configured", true))
		}
	}

	check.SetResult(compliantObjects, nonCompliantObjects)
}

// testPodHostNetwork verifies that the pod hostNetwork parameter is not set to true
func testPodHostNetwork(check *checksdb.Check, env *provider.TestEnvironment) {
	var compliantObjects []*testhelper.ReportObject
	var nonCompliantObjects []*testhelper.ReportObject
	for _, put := range env.Pods {
		check.LogInfo("Testing Pod %q", put)
		if put.Spec.HostNetwork {
			check.LogError("Host network is set to true in Pod %q.", put)
			nonCompliantObjects = append(nonCompliantObjects, testhelper.NewPodReportObject(put.Namespace, put.Name, "Host network is set to true", false))
		} else {
			check.LogInfo("Host network is set to false in Pod %q.", put)
			compliantObjects = append(compliantObjects, testhelper.NewPodReportObject(put.Namespace, put.Name, "Host network is not set to true", true))
		}
	}

	check.SetResult(compliantObjects, nonCompliantObjects)
}

// testPodHostPath verifies that the pod hostpath parameter is not set to true
func testPodHostPath(check *checksdb.Check, env *provider.TestEnvironment) {
	var compliantObjects []*testhelper.ReportObject
	var nonCompliantObjects []*testhelper.ReportObject
	for _, put := range env.Pods {
		check.LogInfo("Testing Pod %q", put)
		podIsCompliant := true
		for idx := range put.Spec.Volumes {
			vol := &put.Spec.Volumes[idx]
			if vol.HostPath != nil && vol.HostPath.Path != "" {
				check.LogError("Hostpath path: %q is set in Pod %q.", vol.HostPath.Path, put)
				nonCompliantObjects = append(nonCompliantObjects, testhelper.NewPodReportObject(put.Namespace, put.Name, "Hostpath path is set", false).
					SetType(testhelper.HostPathType).
					AddField(testhelper.Path, vol.HostPath.Path))
				podIsCompliant = false
			}
		}
		if podIsCompliant {
			check.LogError("Hostpath path not set in Pod %q.", put)
			compliantObjects = append(compliantObjects, testhelper.NewPodReportObject(put.Namespace, put.Name, "Hostpath path is not set", true))
		}
	}

	check.SetResult(compliantObjects, nonCompliantObjects)
}

// testPodHostIPC verifies that the pod hostIpc parameter is not set to true
func testPodHostIPC(check *checksdb.Check, env *provider.TestEnvironment) {
	var compliantObjects []*testhelper.ReportObject
	var nonCompliantObjects []*testhelper.ReportObject
	for _, put := range env.Pods {
		check.LogInfo("Testing Pod %q", put)
		if put.Spec.HostIPC {
			check.LogError("HostIpc is set in Pod %q.", put)
			nonCompliantObjects = append(nonCompliantObjects, testhelper.NewPodReportObject(put.Namespace, put.Name, "HostIpc is set to true", false))
		} else {
			check.LogInfo("HostIpc not set in Pod %q.", put)
			compliantObjects = append(compliantObjects, testhelper.NewPodReportObject(put.Namespace, put.Name, "HostIpc is not set to true", true))
		}
	}

	check.SetResult(compliantObjects, nonCompliantObjects)
}

// testPodHostPID verifies that the pod hostPid parameter is not set to true
func testPodHostPID(check *checksdb.Check, env *provider.TestEnvironment) {
	var compliantObjects []*testhelper.ReportObject
	var nonCompliantObjects []*testhelper.ReportObject
	for _, put := range env.Pods {
		check.LogInfo("Testing Pod %q", put)
		if put.Spec.HostPID {
			check.LogError("HostPid is set in Pod %q.", put)
			nonCompliantObjects = append(nonCompliantObjects, testhelper.NewPodReportObject(put.Namespace, put.Name, "HostPid is set to true", false))
		} else {
			check.LogInfo("HostPid not set in Pod %q.", put)
			compliantObjects = append(compliantObjects, testhelper.NewPodReportObject(put.Namespace, put.Name, "HostPid is not set to true", true))
		}
	}

	check.SetResult(compliantObjects, nonCompliantObjects)
}

// testNamespace Tests namespaces for invalid prefixes and CRs that are not defined in namespaces under test.
func testNamespace(check *checksdb.Check, env *provider.TestEnvironment) {
	var compliantObjects []*testhelper.ReportObject
	var nonCompliantObjects []*testhelper.ReportObject
	for _, namespace := range env.Namespaces {
		check.LogInfo("Testing namespace %q", namespace)
		namespaceCompliant := true
		for _, invalidPrefix := range invalidNamespacePrefixes {
			if strings.HasPrefix(namespace, invalidPrefix) {
				check.LogError("Namespace %q has invalid prefix %q", namespace, invalidPrefix)
				nonCompliantObjects = append(nonCompliantObjects, testhelper.NewNamespacedReportObject("Namespace has invalid prefix", testhelper.Namespace, false, namespace))
				namespaceCompliant = false
				break // Break out of the loop if we find an invalid prefix
			}
		}
		if namespaceCompliant {
			check.LogInfo("Namespace %q has valid prefix", namespace)
			compliantObjects = append(compliantObjects, testhelper.NewNamespacedReportObject("Namespace has valid prefix", testhelper.Namespace, true, namespace))
		}
	}
	if failedNamespacesNum := len(nonCompliantObjects); failedNamespacesNum > 0 {
		check.SetResult(compliantObjects, nonCompliantObjects)
	}

	invalidCrs, err := namespace.TestCrsNamespaces(env.Crds, env.Namespaces, check.GetLogger())
	if err != nil {
		check.LogError("Error while testing CRs namespaces, err=%v", err)
		return
	}

	invalidCrsNum := namespace.GetInvalidCRsNum(invalidCrs, check.GetLogger())
	if invalidCrsNum > 0 {
		nonCompliantObjects = append(nonCompliantObjects, testhelper.NewReportObject("CRs are not in the configured namespaces", testhelper.Namespace, false))
	} else {
		compliantObjects = append(compliantObjects, testhelper.NewReportObject("CRs are in the configured namespaces", testhelper.Namespace, true))
	}
	check.SetResult(compliantObjects, nonCompliantObjects)
}

// testPodServiceAccount verifies that the pod utilizes a valid service account
func testPodServiceAccount(check *checksdb.Check, env *provider.TestEnvironment) {
	var compliantObjects []*testhelper.ReportObject
	var nonCompliantObjects []*testhelper.ReportObject
	for _, put := range env.Pods {
		check.LogInfo("Testing Pod %q", put)
		if put.Spec.ServiceAccountName == defaultServiceAccount {
			check.LogError("Pod %q does not have a valid service account name (uses the default service account instead).", put)
			nonCompliantObjects = append(nonCompliantObjects, testhelper.NewPodReportObject(put.Namespace, put.Name, "Pod does not have a valid service account name", false))
		} else {
			check.LogInfo("Pod %q has a valid service account name", put)
			compliantObjects = append(compliantObjects, testhelper.NewPodReportObject(put.Namespace, put.Name, "Pod has a service account name", true))
		}
	}

	check.SetResult(compliantObjects, nonCompliantObjects)
}

// testPodRoleBindings verifies that the pod utilizes a valid role binding that does not cross non-CNF namespaces
//
//nolint:funlen
func testPodRoleBindings(check *checksdb.Check, env *provider.TestEnvironment) {
	var compliantObjects []*testhelper.ReportObject
	var nonCompliantObjects []*testhelper.ReportObject

	for _, put := range env.Pods {
		check.LogInfo("Testing Pod %q", put)
		podIsCompliant := true
		if put.Pod.Spec.ServiceAccountName == defaultServiceAccount {
			check.LogError("Pod %q has an empty or default serviceAccountName", put)
			// Add the pod to the non-compliant list
			nonCompliantObjects = append(nonCompliantObjects,
				testhelper.NewPodReportObject(put.Namespace, put.Name,
					"The serviceAccountName is either empty or default", false))
			podIsCompliant = false
		} else {
			check.LogInfo("Pod %q has a serviceAccountName %q, checking role bindings.", put, put.Spec.ServiceAccountName)
			// Loop through the rolebindings and check if they are from another namespace
			for rbIndex := range env.RoleBindings {
				// Short circuit if the role binding and the pod are in the same namespace.
				if env.RoleBindings[rbIndex].Namespace == put.Namespace {
					check.LogInfo("Pod %q and the role binding are in the same namespace", put)
					continue
				}
				// If we make it to this point, the role binding and the pod are in different namespaces.
				// We must check if the pod's service account is in the role binding's subjects.
				found := false
				for _, subject := range env.RoleBindings[rbIndex].Subjects {
					// If the subject is a service account and the service account is in the same namespace as one of the CNF's namespaces, then continue, this is allowed
					if subject.Kind == rbacv1.ServiceAccountKind &&
						subject.Namespace == put.Namespace &&
						subject.Name == put.Spec.ServiceAccountName &&
						stringhelper.StringInSlice[string](env.Namespaces, env.RoleBindings[rbIndex].Namespace, false) {
						continue
					}

					// Finally, if the subject is a service account and the service account is in the same namespace as the pod, then we have a failure
					if subject.Kind == rbacv1.ServiceAccountKind &&
						subject.Namespace == put.Namespace &&
						subject.Name == put.Spec.ServiceAccountName {
						check.LogError("Pod %q has the following role bindings that do not live in one of the CNF namespaces: %q", put, env.RoleBindings[rbIndex].Name)

						// Add the pod to the non-compliant list
						nonCompliantObjects = append(nonCompliantObjects,
							testhelper.NewPodReportObject(put.Namespace, put.Name,
								"The role bindings used by this pod do not live in one of the CNF namespaces", false).
								AddField(testhelper.RoleBindingName, env.RoleBindings[rbIndex].Name).
								AddField(testhelper.RoleBindingNamespace, env.RoleBindings[rbIndex].Namespace).
								AddField(testhelper.ServiceAccountName, put.Spec.ServiceAccountName).
								SetType(testhelper.PodRoleBinding))
						found = true
						podIsCompliant = false
						break
					}
				}
				// Break of out the loop if we found a role binding that is out of namespace
				if found {
					break
				}
			}
		}
		// Add pod to the compliant object list
		if podIsCompliant {
			check.LogInfo("All the role bindings used by Pod %q (applied by the service accounts) live in one of the CNF namespaces", put)
			compliantObjects = append(compliantObjects,
				testhelper.NewPodReportObject(put.Namespace, put.Name, "All the role bindings used by this pod (applied by the service accounts) live in one of the CNF namespaces", true))
		}
	}
	check.SetResult(compliantObjects, nonCompliantObjects)
}

// testPodClusterRoleBindings verifies that the pod does not use a cluster role binding
//
//nolint:dupl
func testPodClusterRoleBindings(check *checksdb.Check, env *provider.TestEnvironment) {
	var compliantObjects []*testhelper.ReportObject
	var nonCompliantObjects []*testhelper.ReportObject

	for _, put := range env.Pods {
		check.LogInfo("Testing Pod %q", put)
		result, roleRefName, err := put.IsUsingClusterRoleBinding(env.ClusterRoleBindings, check.GetLogger())
		if err != nil {
			check.LogError("Failed to determine if Pod %q is using a cluster role binding, err=%v", put, err)
			nonCompliantObjects = append(nonCompliantObjects, testhelper.NewPodReportObject(put.Namespace, put.Name, fmt.Sprintf("failed to determine if pod is using a cluster role binding: %v", err), false).
				AddField(testhelper.ClusterRoleName, roleRefName))
			continue
		}

		topOwners, err := put.GetTopOwner()
		if err != nil {
			check.LogError("Could not get top owners of Pod %q, err=%v", put, err)
			nonCompliantObjects = append(nonCompliantObjects, testhelper.NewPodReportObject(put.Namespace, put.Name, fmt.Sprintf("Error getting top owners of this pod, err=%s", err), false).
				AddField(testhelper.ClusterRoleName, roleRefName))
			continue
		}

		csvNamespace, csvName, isOwnedByClusterWideOperator := ownedByClusterWideOperator(topOwners, env)
		// Pod is using a cluster role binding but is owned by a cluster wide operator, so it is ok
		if isOwnedByClusterWideOperator && result {
			check.LogInfo("Pod %q is using a cluster role binding but is owned by a cluster-wide operator (Csv %q, namespace %q)", put, csvName, csvNamespace)
			compliantObjects = append(compliantObjects, testhelper.NewPodReportObject(put.Namespace, put.Name, "Pod is using a cluster role binding but owned by a cluster-wide operator", true))
			continue
		}
		if result {
			// Pod was found to be using a cluster role binding.  This is not allowed.
			// Flagging this pod as a failed pod.
			check.LogError("Pod %q is using a cluster role binding (roleRefName=%q)", put, roleRefName)
			nonCompliantObjects = append(nonCompliantObjects, testhelper.NewPodReportObject(put.Namespace, put.Name, "Pod is using a cluster role binding", false).
				AddField(testhelper.ClusterRoleName, roleRefName))
			continue
		}
		check.LogInfo("Pod %q is not using a cluster role binding", put)
		compliantObjects = append(compliantObjects, testhelper.NewPodReportObject(put.Namespace, put.Name, "Pod is not using a cluster role binding", true))
	}
	check.SetResult(compliantObjects, nonCompliantObjects)
}

// isCSVAndClusterWide checks if object identified by namespace and name is a CSV created by a cluster-wide operator
// Return:
//   - bool : true if object identified by namespace and name is a CSV created by a cluster-wide operator, otherwise return false
func isCSVAndClusterWide(aNamespace, name string, env *provider.TestEnvironment) bool {
	for _, op := range env.Operators {
		if op.Csv != nil &&
			op.Csv.Namespace == aNamespace &&
			op.Csv.Name == name &&
			(op.IsClusterWide || isInstallModeMultiNamespace(op.Csv.Spec.InstallModes)) {
			return true
		}
	}
	return false
}

// isInstallModeMultiNamespace checks if CSV install mode contains multi namespaces or all namespaces
// Return:
//   - bool : true if CSV install mode contains multi namespaces or all namespaces, otherwise return false
func isInstallModeMultiNamespace(installModes []v1alpha1.InstallMode) bool {
	for i := 0; i < len(installModes); i++ {
		if installModes[i].Type == v1alpha1.InstallModeTypeAllNamespaces {
			return true
		}
	}
	return false
}

// ownedByClusterWideOperator checks if one of the passed topOwners is a CSV that is installed by a cluster-wide operator.
// Return:
//   - bool: true if one of the passed topOwners is a CSV that is installed by a cluster-wide operator, otherwise return false
//   - name string : the name of the matching object, if found.
//   - aNamespace string : the namespace of the matching object, if found.
func ownedByClusterWideOperator(topOwners map[string]podhelper.TopOwner, env *provider.TestEnvironment) (aNamespace, name string, found bool) {
	for _, owner := range topOwners {
		if isCSVAndClusterWide(owner.Namespace, owner.Name, env) {
			return owner.Namespace, owner.Name, true
		}
	}
	return "", "", false
}

// testAutomountServiceToken checks if each pod uses the default service account name and if the token is explicitly set in the Pod's spec or if it is inherited from the associated ServiceAccount.
// It sets the result of a compliance check based on the analysis of lists of compliant and non-compliant objects.
func testAutomountServiceToken(check *checksdb.Check, env *provider.TestEnvironment) {
	var compliantObjects []*testhelper.ReportObject
	var nonCompliantObjects []*testhelper.ReportObject
	for _, put := range env.Pods {
		check.LogInfo("Testing Pod %q", put)
		if put.Spec.ServiceAccountName == defaultServiceAccount {
			check.LogError("Pod %q uses the default service account name.", put.Name)
			nonCompliantObjects = append(nonCompliantObjects, testhelper.NewPodReportObject(put.Namespace, put.Name, "Pod has been found with default service account name", false))
			break
		}

		// Evaluate the pod's automount service tokens and any attached service accounts
		client := clientsholder.GetClientsHolder()
		podPassed, newMsg := rbac.EvaluateAutomountTokens(client.K8sClient.CoreV1(), put.Pod)
		if !podPassed {
			check.LogError(newMsg)
			nonCompliantObjects = append(nonCompliantObjects, testhelper.NewPodReportObject(put.Namespace, put.Name, newMsg, false))
		} else {
			check.LogInfo("Pod %q does not have automount service tokens set to true", put)
			compliantObjects = append(compliantObjects, testhelper.NewPodReportObject(put.Namespace, put.Name, "Pod does not have automount service tokens set to true", true))
		}
	}

	check.SetResult(compliantObjects, nonCompliantObjects)
}

// testOneProcessPerContainer is a function that checks if each container(except Istio proxy containers) has only one process running.
// It sets the result of a compliance check based on the analysis of lists of compliant and non-compliant objects.
func testOneProcessPerContainer(check *checksdb.Check, env *provider.TestEnvironment) {
	var compliantObjects []*testhelper.ReportObject
	var nonCompliantObjects []*testhelper.ReportObject

	for _, cut := range env.Containers {
		check.LogInfo("Testing Container %q", cut)
		// the Istio sidecar container "istio-proxy" launches two processes: "pilot-agent" and "envoy"
		if cut.IsIstioProxy() {
			check.LogInfo("Skipping \"istio-proxy\" container")
			continue
		}
		debugPod := env.DebugPods[cut.NodeName]
		if debugPod == nil {
			check.LogError("Debug pod not found for node %q", cut.NodeName)
			return
		}
		ocpContext := clientsholder.NewContext(debugPod.Namespace, debugPod.Name, debugPod.Spec.Containers[0].Name)
		pid, err := crclient.GetPidFromContainer(cut, ocpContext)
		if err != nil {
			check.LogError("Could not get PID for Container %q, error: %v", cut, err)
			nonCompliantObjects = append(nonCompliantObjects, testhelper.NewContainerReportObject(cut.Namespace, cut.Podname, cut.Name, err.Error(), false))
			continue
		}

		nbProcesses, err := getNbOfProcessesInPidNamespace(ocpContext, pid, clientsholder.GetClientsHolder())
		if err != nil {
			check.LogError("Could not get number of processes for Container %q, error: %v", cut, err)
			nonCompliantObjects = append(nonCompliantObjects, testhelper.NewContainerReportObject(cut.Namespace, cut.Podname, cut.Name, err.Error(), false))
			continue
		}
		if nbProcesses > 1 {
			check.LogError("Container %q has more than one process running", cut)
			nonCompliantObjects = append(nonCompliantObjects, testhelper.NewContainerReportObject(cut.Namespace, cut.Podname, cut.Name, "Container has more than one process running", false))
		} else {
			check.LogInfo("Container %q has only one process running", cut)
			compliantObjects = append(compliantObjects, testhelper.NewContainerReportObject(cut.Namespace, cut.Podname, cut.Name, "Container has only one process running", true))
		}
	}

	check.SetResult(compliantObjects, nonCompliantObjects)
}

// testSYSNiceRealtimeCapability is a function that checks if each container running on a realtime kernel enabled node has the SYS_NICE capability.
// It sets the result of a compliance check based on the analysis of lists of compliant and non-compliant objects.
func testSYSNiceRealtimeCapability(check *checksdb.Check, env *provider.TestEnvironment) {
	var compliantObjects []*testhelper.ReportObject
	var nonCompliantObjects []*testhelper.ReportObject

	// Loop through all of the labeled containers and compare their security context capabilities and whether
	// or not the node's kernel is realtime enabled.
	for _, cut := range env.Containers {
		check.LogInfo("Testing Container %q", cut)
		n := env.Nodes[cut.NodeName]
		if !n.IsRTKernel() {
			check.LogInfo("Container is not running on a realtime kernel enabled node")
			compliantObjects = append(compliantObjects, testhelper.NewContainerReportObject(cut.Namespace, cut.Podname, cut.Name, "Container is not running on a realtime kernel enabled node", true))
			continue
		}
		if !strings.Contains(cut.SecurityContext.Capabilities.String(), "SYS_NICE") {
			check.LogError("Container %q has been found running on a realtime kernel enabled node without SYS_NICE capability.", cut)
			nonCompliantObjects = append(nonCompliantObjects, testhelper.NewContainerReportObject(cut.Namespace, cut.Podname, cut.Name, "Container is running on a realtime kernel enabled node without SYS_NICE capability", false))
		} else {
			check.LogInfo("Container is running on a realtime kernel enabled node with the SYS_NICE capability")
			compliantObjects = append(compliantObjects, testhelper.NewContainerReportObject(cut.Namespace, cut.Podname, cut.Name, "Container is running on a realtime kernel enabled node with the SYS_NICE capability", true))
		}
	}

	check.SetResult(compliantObjects, nonCompliantObjects)
}

// testSysPtraceCapability is a function that checks if each pod has process namespace sharing enabled and at least one container allowing the SYS_PTRACE capability.
// It sets the result of a compliance check based on the analysis of lists of compliant and non-compliant objects.
func testSysPtraceCapability(check *checksdb.Check, env *provider.TestEnvironment) {
	var compliantObjects []*testhelper.ReportObject
	var nonCompliantObjects []*testhelper.ReportObject
	for _, put := range env.GetShareProcessNamespacePods() {
		check.LogInfo("Testing Pod %q", put)
		sysPtraceEnabled := false
		for _, cut := range put.Containers {
			if cut.SecurityContext == nil ||
				cut.SecurityContext.Capabilities == nil ||
				len(cut.SecurityContext.Capabilities.Add) == 0 {
				continue
			}
			if stringhelper.StringInSlice(cut.SecurityContext.Capabilities.Add, "SYS_PTRACE", false) {
				check.LogInfo("Container %q defines the SYS_PTRACE capability", cut)
				sysPtraceEnabled = true
				break
			}
		}
		if !sysPtraceEnabled {
			check.LogError("Pod %q has process namespace sharing enabled but no container allowing the SYS_PTRACE capability.", put)
			nonCompliantObjects = append(nonCompliantObjects, testhelper.NewPodReportObject(put.Namespace, put.Name, "Pod has process namespace sharing enabled but no container allowing the SYS_PTRACE capability", false))
		} else {
			check.LogInfo("Pod %q has process namespace sharing enabled and at least one container allowing the SYS_PTRACE capability", put)
			compliantObjects = append(compliantObjects, testhelper.NewPodReportObject(put.Namespace, put.Name, "Pod has process namespace sharing enabled and at least one container allowing the SYS_PTRACE capability", true))
		}
	}
	check.SetResult(compliantObjects, nonCompliantObjects)
}

// testNamespaceResourceQuota is a function that checks if each pod is running in a namespace that has a ResourceQuota applied.
// It sets the result of a compliance check based on the analysis of lists of compliant and non-compliant objects.
func testNamespaceResourceQuota(check *checksdb.Check, env *provider.TestEnvironment) {
	var compliantObjects []*testhelper.ReportObject
	var nonCompliantObjects []*testhelper.ReportObject

	for _, put := range env.Pods {
		check.LogInfo("Testing Pod %q", put)
		// Look through all of the pods and compare their namespace to any potential
		// resource quotas
		foundPodNamespaceRQ := false
		for index := range env.ResourceQuotas {
			// We are just checking for the existence of the resource quota as of right now.
			// Read more about the resource quota object here:
			// https://kubernetes.io/docs/concepts/policy/resource-quotas/
			if put.Namespace == env.ResourceQuotas[index].Namespace {
				foundPodNamespaceRQ = true
				break
			}
		}

		if !foundPodNamespaceRQ {
			check.LogError("Pod %q is running in a namespace that does not have a ResourceQuota applied.", put)
			nonCompliantObjects = append(nonCompliantObjects, testhelper.NewPodReportObject(put.Namespace, put.Name, "Pod is running in a namespace that does not have a ResourceQuota applied", false))
		} else {
			check.LogInfo("Pod %q is running in a namespace that has a ResourceQuota applied.", put)
			compliantObjects = append(compliantObjects, testhelper.NewPodReportObject(put.Namespace, put.Name, "Pod is running in a namespace that has a ResourceQuota applied", true))
		}
	}

	check.SetResult(compliantObjects, nonCompliantObjects)
}

const (
	sshServicePortProtocol = "TCP"
)

// testNoSSHDaemonsAllowed is a function that checks if each pod is running an SSH daemon.
// It sets the result of a compliance check based on the analysis of lists of compliant and non-compliant objects.
func testNoSSHDaemonsAllowed(check *checksdb.Check, env *provider.TestEnvironment) {
	var compliantObjects []*testhelper.ReportObject
	var nonCompliantObjects []*testhelper.ReportObject

	for _, put := range env.Pods {
		check.LogInfo("Testing Pod %q", put)
		cut := put.Containers[0]

		// 1. Find SSH port
		port, err := netutil.GetSSHDaemonPort(cut)
		if err != nil {
			check.LogError("Could not get ssh daemon port on %q, err: %v", cut, err)
			nonCompliantObjects = append(nonCompliantObjects, testhelper.NewPodReportObject(put.Namespace, put.Name, "Failed to get the ssh port for pod", false))
			continue
		}

		if port == "" {
			check.LogInfo("Pod %q is not running an SSH daemon", put)
			compliantObjects = append(compliantObjects, testhelper.NewPodReportObject(put.Namespace, put.Name, "Pod is not running an SSH daemon", true))
			continue
		}

		sshServicePortNumber, err := strconv.ParseInt(port, 10, 32)
		if err != nil {
			check.LogError("Could not convert port %q from string to integer on Container %q", port, cut)
			nonCompliantObjects = append(nonCompliantObjects, testhelper.NewPodReportObject(put.Namespace, put.Name, "Failed to get the listening ports for pod", false))
			continue
		}

		// 2. Check if SSH port is listening
		sshPortInfo := netutil.PortInfo{PortNumber: int32(sshServicePortNumber), Protocol: sshServicePortProtocol}
		listeningPorts, err := netutil.GetListeningPorts(cut)
		if err != nil {
			check.LogError("Failed to get the listening ports for Pod %q, err: %v", put, err)
			nonCompliantObjects = append(nonCompliantObjects, testhelper.NewPodReportObject(put.Namespace, put.Name, "Failed to get the listening ports for pod", false))
			continue
		}

		if _, ok := listeningPorts[sshPortInfo]; ok {
			check.LogError("Pod %q is running an SSH daemon", put)
			nonCompliantObjects = append(nonCompliantObjects, testhelper.NewPodReportObject(put.Namespace, put.Name, "Pod is running an SSH daemon", false))
		} else {
			check.LogInfo("Pod %q is not running an SSH daemon", put)
			compliantObjects = append(compliantObjects, testhelper.NewPodReportObject(put.Namespace, put.Name, "Pod is not running an SSH daemon", true))
		}
	}

	check.SetResult(compliantObjects, nonCompliantObjects)
}

// testPodRequestsAndLimits is a function that checks if each container has resource requests and limits.
// It sets the result of a compliance check based on the analysis of lists of compliant and non-compliant objects.
func testPodRequestsAndLimits(check *checksdb.Check, env *provider.TestEnvironment) {
	var compliantObjects []*testhelper.ReportObject
	var nonCompliantObjects []*testhelper.ReportObject
	// Loop through the containers, looking for containers that are missing requests or limits.
	// These need to be defined in order to pass.
	for _, cut := range env.Containers {
		check.LogInfo("Testing Container %q", cut)
		if !resources.HasRequestsAndLimitsSet(cut, check.GetLogger()) {
			check.LogError("Container %q is missing resource requests or limits", cut)
			nonCompliantObjects = append(nonCompliantObjects, testhelper.NewContainerReportObject(cut.Namespace, cut.Podname, cut.Name, "Container is missing resource requests or limits", false))
		} else {
			check.LogInfo("Container %q has resource requests and limits", cut)
			compliantObjects = append(compliantObjects, testhelper.NewContainerReportObject(cut.Namespace, cut.Podname, cut.Name, "Container has resource requests and limits", true))
		}
	}

	check.SetResult(compliantObjects, nonCompliantObjects)
}

// test1337UIDs is a function that checks if each pod is using securityContext RunAsUser 1337.
// It sets the result of a compliance check based on the analysis of lists of compliant and non-compliant objects.
func test1337UIDs(check *checksdb.Check, env *provider.TestEnvironment) {
	// Note this test is only ran as part of the 'extended' test suite.
	var compliantObjects []*testhelper.ReportObject
	var nonCompliantObjects []*testhelper.ReportObject
	const leetNum = 1337
	for _, put := range env.Pods {
		check.LogInfo("Testing Pod %q", put)
		if put.IsRunAsUserID(leetNum) {
			check.LogError("Pod %q is using securityContext RunAsUser 1337", put)
			nonCompliantObjects = append(nonCompliantObjects, testhelper.NewPodReportObject(put.Namespace, put.Name, "Pod is using securityContext RunAsUser 1337", false))
		} else {
			check.LogInfo("Pod %q is not using securityContext RunAsUser 1337", put)
			compliantObjects = append(compliantObjects, testhelper.NewPodReportObject(put.Namespace, put.Name, "Pod is not using securityContext RunAsUser 1337", true))
		}
	}

	check.SetResult(compliantObjects, nonCompliantObjects)
}

// testContainerSCC categorizes the containers under test into several categories of increasing privileges based on their SCC.
// Containers not compliant with the least privileged category fail this test.
// It sets the result of a compliance check based on the analysis of lists of compliant and non-compliant objects.
func testContainerSCC(check *checksdb.Check, env *provider.TestEnvironment) {
	var compliantObjects []*testhelper.ReportObject
	var nonCompliantObjects []*testhelper.ReportObject
	highLevelCat := securitycontextcontainer.CategoryID1
	for _, pod := range env.Pods {
		check.LogInfo("Testing Pod %q", pod)
		listCategory := securitycontextcontainer.CheckPod(pod)
		for _, cat := range listCategory {
			if cat.Category > securitycontextcontainer.CategoryID1NoUID0 {
				check.LogError("Category %q is NOT category 1 or category NoUID0", cat)
				aContainerOut := testhelper.NewContainerReportObject(cat.NameSpace, cat.Podname, cat.Containername, "container category is NOT category 1 or category NoUID0", false).
					SetType(testhelper.ContainerCategory).
					AddField(testhelper.Category, cat.Category.String())
				nonCompliantObjects = append(nonCompliantObjects, aContainerOut)
			} else {
				check.LogInfo("Category %q is category 1 or category NoUID0", cat)
				aContainerOut := testhelper.NewContainerReportObject(cat.NameSpace, cat.Podname, cat.Containername, "container category is category 1 or category NoUID0", true).
					SetType(testhelper.ContainerCategory).
					AddField(testhelper.Category, cat.Category.String())
				compliantObjects = append(compliantObjects, aContainerOut)
			}
			if cat.Category > highLevelCat {
				highLevelCat = cat.Category
			}
		}
	}
	aCNFOut := testhelper.NewReportObject("Overall CNF category", testhelper.CnfType, false).AddField(testhelper.Category, highLevelCat.String())
	compliantObjects = append(compliantObjects, aCNFOut)
	check.SetResult(compliantObjects, nonCompliantObjects)
}

// testNodePort is a function that checks for each service type if it is nodePort.
// It sets the result of a compliance check based on the analysis of lists of compliant and non-compliant objects.
func testNodePort(check *checksdb.Check, env *provider.TestEnvironment) {
	var compliantObjects []*testhelper.ReportObject
	var nonCompliantObjects []*testhelper.ReportObject
	for _, s := range env.Services {
		check.LogInfo("Testing %q", services.ToString(s))

		if s.Spec.Type == nodePort {
			check.LogError("Service %q (ns %q) type is nodePort", s.Name, s.Namespace)
			nonCompliantObjects = append(nonCompliantObjects, testhelper.NewReportObject("Service is type NodePort", testhelper.ServiceType, false).
				AddField(testhelper.Namespace, s.Namespace).
				AddField(testhelper.ServiceName, s.Name).
				AddField(testhelper.ServiceMode, string(s.Spec.Type)))
		} else {
			check.LogInfo("Service %q (ns %q) type is not nodePort (type=%q)", s.Name, s.Namespace, s.Spec.Type)
			compliantObjects = append(compliantObjects, testhelper.NewReportObject("Service is not type NodePort", testhelper.ServiceType, true).
				AddField(testhelper.Namespace, s.Namespace).
				AddField(testhelper.ServiceName, s.Name).
				AddField(testhelper.ServiceMode, string(s.Spec.Type)))
		}
	}
	check.SetResult(compliantObjects, nonCompliantObjects)
}

// testCrdRoles is a function that checks for each role applies only to CRDs under test.
// It sets the result of a compliance check based on the analysis of lists of compliant and non-compliant objects.
func testCrdRoles(check *checksdb.Check, env *provider.TestEnvironment) {
	var compliantObjects []*testhelper.ReportObject
	var nonCompliantObjects []*testhelper.ReportObject
	crdResources := rbac.GetCrdResources(env.Crds)
	for roleIndex := range env.Roles {
		if !stringhelper.StringInSlice[string](env.Namespaces, env.Roles[roleIndex].Namespace, false) {
			continue
		}

		allRules := rbac.GetAllRules(&env.Roles[roleIndex])

		matchingRules, nonMatchingRules := rbac.FilterRulesNonMatchingResources(allRules, crdResources)
		if len(matchingRules) == 0 {
			continue
		}
		for _, aRule := range matchingRules {
			check.LogInfo("Rule (resource-name=%q, resource-group=%q, verb=%q, role-name=%q) applies to CRDs under test",
				aRule.Resource.Name, aRule.Resource.Group, aRule.Verb, env.Roles[roleIndex].Name)
			compliantObjects = append(compliantObjects, testhelper.NewNamespacedReportObject("This applies to CRDs under test", testhelper.RoleRuleType, true, env.Roles[roleIndex].Namespace).
				AddField(testhelper.RoleName, env.Roles[roleIndex].Name).
				AddField(testhelper.Group, aRule.Resource.Group).
				AddField(testhelper.ResourceName, aRule.Resource.Name).
				AddField(testhelper.Verb, aRule.Verb))
		}
		for _, aRule := range nonMatchingRules {
			check.LogInfo("Rule (resource-name=%q, resource-group=%q, verb=%q, role-name=%q) does not apply to CRDs under test",
				aRule.Resource.Name, aRule.Resource.Group, aRule.Verb, env.Roles[roleIndex].Name)
			nonCompliantObjects = append(nonCompliantObjects, testhelper.NewNamespacedReportObject("This rule does not apply to CRDs under test", testhelper.RoleRuleType, false, env.Roles[roleIndex].Namespace).
				AddField(testhelper.RoleName, env.Roles[roleIndex].Name).
				AddField(testhelper.Group, aRule.Resource.Group).
				AddField(testhelper.ResourceName, aRule.Resource.Name).
				AddField(testhelper.Verb, aRule.Verb))
		}

		if len(nonMatchingRules) == 0 {
			check.LogInfo("Role %q rules only apply to CRDs under test", env.Roles[roleIndex].Name)
			compliantObjects = append(compliantObjects, testhelper.NewNamespacedNamedReportObject("This role's rules only apply to CRDs under test",
				testhelper.RoleType, true, env.Roles[roleIndex].Namespace, env.Roles[roleIndex].Name))
		} else {
			check.LogError("Role %q rules apply to a mix of CRDs under test and others.", env.Roles[roleIndex].Name)
			nonCompliantObjects = append(nonCompliantObjects, testhelper.NewNamespacedNamedReportObject("This role's rules apply to a mix of CRDs under test and others. See non compliant role rule objects.",
				testhelper.RoleType, false, env.Roles[roleIndex].Namespace, env.Roles[roleIndex].Name))
		}
	}
	check.SetResult(compliantObjects, nonCompliantObjects)
}
