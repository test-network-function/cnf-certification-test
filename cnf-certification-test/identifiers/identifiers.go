// Copyright (C) 2021 Red Hat, Inc.
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

package identifiers

import (
	"fmt"
	"strings"

	"github.com/test-network-function/cnf-certification-test/cnf-certification-test/common"
	"github.com/test-network-function/test-network-function-claim/pkg/claim"
)

const (
	informativeResult        = "informative"
	normativeResult          = "normative"
	url                      = "http://test-network-function.com/testcases"
	versionOne               = "v1.0.0"
	bestPracticeDocV1dot3URL = "https://connect.redhat.com/sites/default/files/2022-05/Cloud%20Native%20Network%20Function%20Requirements%201-3.pdf"
	bestPracticeDocV1dot4URL = "https://TODO" // TODO: Fill in this variable with the new v1.4 document when available.
)

// TestCaseDescription describes a JUnit test case.
type TestCaseDescription struct {
	// Identifier is the unique test identifier.
	Identifier claim.Identifier `json:"identifier" yaml:"identifier"`

	// Description is a helpful description of the purpose of the test case.
	Description string `json:"description" yaml:"description"`

	// Remediation is an optional suggested remediation for passing the test.
	Remediation string `json:"remediation,omitempty" yaml:"remediation,omitempty"`

	// Type is the type of the test (i.e., normative).
	Type string `json:"type" yaml:"type"`

	// BestPracticeReference is a helpful best practice references of the test case.
	BestPracticeReference string `json:"BestPracticeReference" yaml:"BestPracticeReference"`

	// ExceptionProcess will show any possible exception processes documented for partners to follow.
	ExceptionProcess string `json:"exceptionProcess,omitempty" yaml:"exceptionProcess,omitempty"`
}

func formTestURL(suite, name string) string {
	return fmt.Sprintf("%s/%s/%s", url, suite, name)
}

var (
	// TestIdToClaimId converts the testcase short ID to the claim identifier
	TestIDToClaimID = map[string]claim.Identifier{}

	// TestPodDeleteIdentifier tests for delete pod test
	TestPodDeleteIdentifier = claim.Identifier{
		Url:     formTestURL(common.ChaosTesting, "pod-delete"),
		Version: versionOne,
	}

	// BaseDomain for the test cases
	TestIDBaseDomain = url

	// TestSecConCapabilitiesIdentifier tests for non compliant security context capabilities
	TestSecConCapabilitiesIdentifier = claim.Identifier{
		Url:     formTestURL(common.AccessControlTestKey, "security-context-capabilities-check"),
		Version: versionOne,
	}
	// TestSecConNonRootUserIdentifier tests that pods or containers are not running with root permissions
	TestSecConNonRootUserIdentifier = claim.Identifier{
		Url:     formTestURL(common.AccessControlTestKey, "security-context-non-root-user-check"),
		Version: versionOne,
	}
	// TestSecPrivilegedEscalation tests that containers are not allowed privilege escalation
	TestSecConPrivilegeEscalation = claim.Identifier{
		Url:     formTestURL(common.AccessControlTestKey, "security-context-privilege-escalation"),
		Version: versionOne,
	}
	// TestSecPrivilegedEscalation tests that containers are not configured with host port privileges
	TestContainerHostPort = claim.Identifier{
		Url:     formTestURL(common.AccessControlTestKey, "container-host-port"),
		Version: versionOne,
	}
	// TestPodHostNetwork tests that pods do not configure hostnetwork to true
	TestPodHostNetwork = claim.Identifier{
		Url:     formTestURL(common.AccessControlTestKey, "pod-host-network"),
		Version: versionOne,
	}
	// TestPodHostPath tests that pods do not configure a hostpath volume
	TestPodHostPath = claim.Identifier{
		Url:     formTestURL(common.AccessControlTestKey, "pod-host-path"),
		Version: versionOne,
	}
	// TestPodHostPath tests that pods do not configure a hostpath volume
	TestPodHostIPC = claim.Identifier{
		Url:     formTestURL(common.AccessControlTestKey, "pod-host-ipc"),
		Version: versionOne,
	}
	// TestPodHostPath tests that pods do not configure a hostpath volume
	TestPodHostPID = claim.Identifier{
		Url:     formTestURL(common.AccessControlTestKey, "pod-host-pid"),
		Version: versionOne,
	}
	// TestContainerIsCertifiedIdentifier tests whether the container has passed Container Certification.
	TestContainerIsCertifiedIdentifier = claim.Identifier{
		Url:     formTestURL(common.AffiliatedCertTestKey, "container-is-certified"),
		Version: versionOne,
	}
	// TestHugepagesNotManuallyManipulated represents the test identifier testing hugepages have not been manipulated.
	TestHugepagesNotManuallyManipulated = claim.Identifier{
		Url:     formTestURL(common.PlatformAlterationTestKey, "hugepages-config"),
		Version: versionOne,
	}
	// TestICMPv4ConnectivityIdentifier tests icmpv4 connectivity.
	TestICMPv4ConnectivityIdentifier = claim.Identifier{
		Url:     formTestURL(common.NetworkingTestKey, "icmpv4-connectivity"),
		Version: versionOne,
	}
	// TestICMPv6ConnectivityIdentifier tests icmpv6 connectivity.
	TestICMPv6ConnectivityIdentifier = claim.Identifier{
		Url:     formTestURL(common.NetworkingTestKey, "icmpv6-connectivity"),
		Version: versionOne,
	}
	// TestICMPv4ConnectivityIdentifier tests icmpv4 Multus connectivity.
	TestICMPv4ConnectivityMultusIdentifier = claim.Identifier{
		Url:     formTestURL(common.NetworkingTestKey, "icmpv4-connectivity-multus"),
		Version: versionOne,
	}
	// TestICMPv6ConnectivityIdentifier tests icmpv6 Multus connectivity.
	TestICMPv6ConnectivityMultusIdentifier = claim.Identifier{
		Url:     formTestURL(common.NetworkingTestKey, "icmpv6-connectivity-multus"),
		Version: versionOne,
	}
	// TestNamespaceBestPracticesIdentifier ensures the namespace has followed best namespace practices.
	TestNamespaceBestPracticesIdentifier = claim.Identifier{
		Url:     formTestURL(common.AccessControlTestKey, "namespace"),
		Version: versionOne,
	}
	// TestNonTaintedNodeKernelsIdentifier is the identifier for the test checking tainted nodes.
	TestNonTaintedNodeKernelsIdentifier = claim.Identifier{
		Url:     formTestURL(common.PlatformAlterationTestKey, "tainted-node-kernel"),
		Version: versionOne,
	}
	// TestOperatorInstallStatusSucceededIdentifier tests Operator best practices.
	TestOperatorInstallStatusSucceededIdentifier = claim.Identifier{
		Url:     formTestURL(common.OperatorTestKey, "install-status-succeeded"),
		Version: versionOne,
	}
	// TestOperatorNoPrivileges tests Operator has no privileges on resources.
	TestOperatorNoPrivileges = claim.Identifier{
		Url:     formTestURL(common.OperatorTestKey, "install-status-no-privileges"),
		Version: versionOne,
	}
	// TestOperatorIsCertifiedIdentifier tests that an Operator has passed Operator certification.
	TestOperatorIsCertifiedIdentifier = claim.Identifier{
		Url:     formTestURL(common.AffiliatedCertTestKey, "operator-is-certified"),
		Version: versionOne,
	}
	TestHelmIsCertifiedIdentifier = claim.Identifier{
		Url:     formTestURL(common.AffiliatedCertTestKey, "helmchart-is-certified"),
		Version: versionOne,
	}
	// TestOperatorIsInstalledViaOLMIdentifier tests that an Operator is installed via OLM.
	TestOperatorIsInstalledViaOLMIdentifier = claim.Identifier{
		Url:     formTestURL(common.OperatorTestKey, "install-source"),
		Version: versionOne,
	}
	// TestPodNodeSelectorAndAffinityBestPractices is the test ensuring nodeSelector and nodeAffinity are not used by a
	// Pod.
	TestPodNodeSelectorAndAffinityBestPractices = claim.Identifier{
		Url:     formTestURL(common.LifecycleTestKey, "pod-scheduling"),
		Version: versionOne,
	}
	// TestPodHighAvailabilityBestPractices is the test ensuring podAntiAffinity are used by a
	// Pod when pod replica # are great than 1
	TestPodHighAvailabilityBestPractices = claim.Identifier{
		Url:     formTestURL(common.LifecycleTestKey, "pod-high-availability"),
		Version: versionOne,
	}

	// TestPodClusterRoleBindingsBestPracticesIdentifier ensures Pod crb best practices.
	TestPodClusterRoleBindingsBestPracticesIdentifier = claim.Identifier{
		Url:     formTestURL(common.AccessControlTestKey, "cluster-role-bindings"),
		Version: versionOne,
	}
	// TestPodDeploymentBestPracticesIdentifier ensures a CNF follows best Deployment practices.
	TestPodDeploymentBestPracticesIdentifier = claim.Identifier{
		Url:     formTestURL(common.LifecycleTestKey, "pod-owner-type"),
		Version: versionOne,
	}
	// TestDeploymentScalingIdentifier ensures deployment scale in/out operations work correctly.
	TestDeploymentScalingIdentifier = claim.Identifier{
		Url:     formTestURL(common.LifecycleTestKey, "deployment-scaling"),
		Version: versionOne,
	}
	// TestStateFulSetScalingIdentifier ensures statefulset scale in/out operations work correctly.
	TestStateFulSetScalingIdentifier = claim.Identifier{
		Url:     formTestURL(common.LifecycleTestKey, "statefulset-scaling"),
		Version: versionOne,
	}
	// TestImagePullPolicyIdentifier ensures represent image pull policy practices.
	TestImagePullPolicyIdentifier = claim.Identifier{
		Url:     formTestURL(common.LifecycleTestKey, "image-pull-policy"),
		Version: versionOne,
	}
	// TestPodRecreationIdentifier ensures recreation best practices.
	TestPodRecreationIdentifier = claim.Identifier{
		Url:     formTestURL(common.LifecycleTestKey, "pod-recreation"),
		Version: versionOne,
	}
	// TestPodRoleBindingsBestPracticesIdentifier represents rb best practices.
	TestPodRoleBindingsBestPracticesIdentifier = claim.Identifier{
		Url:     formTestURL(common.AccessControlTestKey, "pod-role-bindings"),
		Version: versionOne,
	}
	// TestPodServiceAccountBestPracticesIdentifier tests Pod SA best practices.
	TestPodServiceAccountBestPracticesIdentifier = claim.Identifier{
		Url:     formTestURL(common.AccessControlTestKey, "pod-service-account"),
		Version: versionOne,
	}
	//
	TestPodAutomountServiceAccountIdentifier = claim.Identifier{
		Url:     formTestURL(common.AccessControlTestKey, "pod-automount-service-account-token"),
		Version: versionOne,
	}
	// TestServicesDoNotUseNodeportsIdentifier ensures Services do not utilize NodePorts.
	TestServicesDoNotUseNodeportsIdentifier = claim.Identifier{
		Url:     formTestURL(common.NetworkingTestKey, "service-type"),
		Version: versionOne,
	}
	// TestUnalteredBaseImageIdentifier ensures the base image is not altered.
	TestUnalteredBaseImageIdentifier = claim.Identifier{
		Url:     formTestURL(common.PlatformAlterationTestKey, "base-image"),
		Version: versionOne,
	}
	// TestUnalteredStartupBootParamsIdentifier ensures startup boot params are not altered.
	TestUnalteredStartupBootParamsIdentifier = claim.Identifier{
		Url:     formTestURL(common.PlatformAlterationTestKey, "boot-params"),
		Version: versionOne,
	}
	// TestLoggingIdentifier ensures stderr/stdout are used
	TestLoggingIdentifier = claim.Identifier{
		Url:     formTestURL(common.ObservabilityTestKey, "container-logging"),
		Version: versionOne,
	}
	// TestTerminationMessagePolicy ensures pods have FallbackToLogsOnError set
	TestTerminationMessagePolicyIdentifier = claim.Identifier{
		Url:     formTestURL(common.ObservabilityTestKey, "termination-policy"),
		Version: versionOne,
	}
	// TestCrdsStatusSubresourceIdentifier ensures all CRDs have a valid status subresource
	TestCrdsStatusSubresourceIdentifier = claim.Identifier{
		Url:     formTestURL(common.ObservabilityTestKey, "crd-status"),
		Version: versionOne,
	}
	// TestShutdownIdentifier ensures pre-stop lifecycle is defined
	TestShutdownIdentifier = claim.Identifier{
		Url:     formTestURL(common.LifecycleTestKey, "container-shutdown"),
		Version: versionOne,
	}
	// TestSysctlConfigsIdentifier ensures that the node's sysctl configs are consistent with the MachineConfig CR
	TestSysctlConfigsIdentifier = claim.Identifier{
		Url:     formTestURL(common.PlatformAlterationTestKey, "sysctl-config"),
		Version: versionOne,
	}
	// TestServiceMesh checks if service mesh is exist.
	TestServiceMeshIdentifier = claim.Identifier{
		Url:     formTestURL(common.PlatformAlterationTestKey, "service-mesh-usage"),
		Version: versionOne,
	}
	// TestOCPLifecycleIdentifier ensures the OCP version of the cluster is within the valid lifecycle status
	TestOCPLifecycleIdentifier = claim.Identifier{
		Url:     formTestURL(common.PlatformAlterationTestKey, "ocp-lifecycle"),
		Version: versionOne,
	}
	// TestNodeOperatingSystemIdentifier ensures workers in the cluster have an operating system that is compatible with the installed version of OpenShift
	TestNodeOperatingSystemIdentifier = claim.Identifier{
		Url:     formTestURL(common.PlatformAlterationTestKey, "ocp-node-os-lifecycle"),
		Version: versionOne,
	}
	// TestScalingIdentifier ensures deployment scale in/out operations work correctly.
	TestScalingIdentifier = claim.Identifier{
		Url:     formTestURL(common.LifecycleTestKey, "scaling"),
		Version: versionOne,
	}
	// TestIsRedHatReleaseIdentifier ensures platform is defined
	TestIsRedHatReleaseIdentifier = claim.Identifier{
		Url:     formTestURL(common.PlatformAlterationTestKey, "isredhat-release"),
		Version: versionOne,
	}
	// TestIsSELinuxEnforcingIdentifier ensures selinux is in enforcing mode
	TestIsSELinuxEnforcingIdentifier = claim.Identifier{
		Url:     formTestURL(common.PlatformAlterationTestKey, "is-selinux-enforcing"),
		Version: versionOne,
	}
	TestUndeclaredContainerPortsUsage = claim.Identifier{
		Url:     formTestURL(common.NetworkingTestKey, "undeclared-container-ports-usage"),
		Version: versionOne,
	}
	TestOCPReservedPortsUsage = claim.Identifier{
		Url:     formTestURL(common.NetworkingTestKey, "ocp-reserved-ports-usage"),
		Version: versionOne,
	}
	TestLivenessProbeIdentifier = claim.Identifier{
		Url:     formTestURL(common.LifecycleTestKey, "liveness"),
		Version: versionOne,
	}
	TestReadinessProbeIdentifier = claim.Identifier{
		Url:     formTestURL(common.LifecycleTestKey, "readiness"),
		Version: versionOne,
	}
	// TestOneProcessPerContainerIdentifier ensures that only one process per container is running
	TestOneProcessPerContainerIdentifier = claim.Identifier{
		Url:     formTestURL(common.AccessControlTestKey, "one-process-per-container"),
		Version: versionOne,
	}
	TestSYSNiceRealtimeCapabilityIdentifier = claim.Identifier{
		Url:     formTestURL(common.AccessControlTestKey, "sys-nice-realtime-capability"),
		Version: versionOne,
	}
	// TestSysPtraceCapabilityIdentifier ensures that if process namespace sharing is enabled then the SYS_PTRACE capability is allowed
	TestSysPtraceCapabilityIdentifier = claim.Identifier{
		Url:     formTestURL(common.AccessControlTestKey, "sys-ptrace-capability"),
		Version: versionOne,
	}
	TestPodRequestsAndLimitsIdentifier = claim.Identifier{
		Url:     formTestURL(common.AccessControlTestKey, "requests-and-limits"),
		Version: versionOne,
	}
	TestNamespaceResourceQuotaIdentifier = claim.Identifier{
		Url:     formTestURL(common.AccessControlTestKey, "namespace-resource-quota"),
		Version: versionOne,
	}
	TestNamespaceLimitRangeIdentifier = claim.Identifier{
		Url:     formTestURL(common.AccessControlTestKey, "namespace-limit-range"),
		Version: versionOne,
	}
)

func formDescription(identifier claim.Identifier, description string) string {
	return fmt.Sprintf("%s %s", identifier.Url, description)
}

// XformToGinkgoItIdentifier transform the claim.Identifier into a test Id that can be used to skip
// specific tests
func XformToGinkgoItIdentifier(identifier claim.Identifier) string {
	return XformToGinkgoItIdentifierExtended(identifier, "")
}

// XformToGinkgoItIdentifierExtended transform the claim.Identifier into a test Id that can be used to skip
// specific tests
func XformToGinkgoItIdentifierExtended(identifier claim.Identifier, extra string) string {
	itID := strings.ReplaceAll(strings.TrimPrefix(identifier.Url, url+"/"), "/", "-")
	var key string
	if extra != "" {
		key = itID + "-" + extra
	} else {
		key = itID
	}
	TestIDToClaimID[key] = identifier
	return key
}

// It extracts the suite name and test name from a claim.Identifier based
// on the const url which contains a base domain
// From a claim.Identifier.url:
//   http://test-network-function.com/tests-case/SuitName/TestName
// It extracts SuitNAme and TestName

func GetSuiteAndTestFromIdentifier(identifier claim.Identifier) []string {
	result := strings.Split(identifier.Url, url+"/")
	const SPLITN = 2
	// len 2, the baseDomain can appear only once in the url
	// so it returns what you have previous and before basedomain
	if len(result) != SPLITN {
		return nil
	}
	return strings.Split(result[1], "/")
}

// Catalog is the JUnit testcase catalog of tests.
var Catalog = map[claim.Identifier]TestCaseDescription{

	TestNodeOperatingSystemIdentifier: {
		Identifier: TestNodeOperatingSystemIdentifier,
		Type:       normativeResult,
		Description: formDescription(TestNodeOperatingSystemIdentifier, `Tests that the nodes running in the cluster have operating systems
			that are compatible with the deployed version of OpenShift.`),
		Remediation:           NodeOperatingSystemRemediation,
		ExceptionProcess:      NoDocumentedProcess,
		BestPracticeReference: bestPracticeDocV1dot3URL + " Section 7.9",
	},

	TestOCPLifecycleIdentifier: {
		Identifier:            TestOCPLifecycleIdentifier,
		Type:                  normativeResult,
		Description:           formDescription(TestOCPLifecycleIdentifier, `Tests that the running OCP version is not end of life.`),
		Remediation:           OCPLifecycleRemediation,
		ExceptionProcess:      NoDocumentedProcess,
		BestPracticeReference: bestPracticeDocV1dot3URL + " Section 7.9",
	},

	TestDeploymentScalingIdentifier: {
		Identifier: TestDeploymentScalingIdentifier,
		Type:       normativeResult,
		Description: formDescription(TestDeploymentScalingIdentifier,
			`tests that CNF deployments support scale in/out operations. 
			First, The test starts getting the current replicaCount (N) of the deployment/s with the Pod Under Test. Then, it executes the 
			scale-in oc command for (N-1) replicas. Lastly, it executes the scale-out oc command, restoring the original replicaCount of the deployment/s.
		    In case of deployments that are managed by HPA the test is changing the min and max value to deployment Replica - 1 during scale-in and the 
			original replicaCount again for both min/max during the scale-out stage. lastly its restoring the original min/max replica of the deployment/s`),
		Remediation:           DeploymentScalingRemediation,
		ExceptionProcess:      NoDocumentedProcess,
		BestPracticeReference: bestPracticeDocV1dot3URL + " Section 5.2",
	},
	TestStateFulSetScalingIdentifier: {
		Identifier: TestStateFulSetScalingIdentifier,
		Type:       normativeResult,
		Description: formDescription(TestStateFulSetScalingIdentifier,
			`tests that CNF statefulsets support scale in/out operations. 
			First, The test starts getting the current replicaCount (N) of the statefulset/s with the Pod Under Test. Then, it executes the 
			scale-in oc command for (N-1) replicas. Lastly, it executes the scale-out oc command, restoring the original replicaCount of the statefulset/s.
			In case of statefulsets that are managed by HPA the test is changing the min and max value to statefulset Replica - 1 during scale-in and the 
			original replicaCount again for both min/max during the scale-out stage. lastly its restoring the original min/max replica of the statefulset/s`),
		Remediation:           StatefulSetScalingRemediation,
		ExceptionProcess:      NoDocumentedProcess,
		BestPracticeReference: bestPracticeDocV1dot3URL + " Section 5.2",
	},
	TestSecConCapabilitiesIdentifier: {
		Identifier:       TestSecConCapabilitiesIdentifier,
		Type:             normativeResult,
		Remediation:      SecConCapabilitiesRemediation,
		ExceptionProcess: SecConCapabilitiesExceptionProcess,
		Description: formDescription(TestSecConCapabilitiesIdentifier,
			`Tests that the following capabilities are not granted:
			- NET_ADMIN
			- SYS_ADMIN 
			- NET_RAW
			- IPC_LOCK
`),
		BestPracticeReference: bestPracticeDocV1dot3URL + " Section 5.2",
	},
	// TestPodDeleteIdentifier: {
	// 	Identifier:  TestPodDeleteIdentifier,
	// 	Type:        normativeResult,
	// 	Remediation: `Make sure that the pods can be recreated successfully after deleting them`,
	// 	Description: formDescription(TestPodDeleteIdentifier,
	// 		`Using the litmus chaos operator, this test checks that pods are recreated successfully after deleting them.`),
	// 	BestPracticeReference: bestPracticeDocV1dot3URL + " Section 5.2",
	// },
	TestSecConNonRootUserIdentifier: {
		Identifier:       TestSecConNonRootUserIdentifier,
		Type:             normativeResult,
		Remediation:      SecConNonRootUserRemediation,
		ExceptionProcess: SecConNonRootUserExceptionProcess,
		Description: formDescription(TestSecConNonRootUserIdentifier,
			`Checks the security context runAsUser parameter in pods and containers to make sure it is not set to uid root(0)`),
		BestPracticeReference: bestPracticeDocV1dot3URL + " Section 5.2",
	},
	TestSecConPrivilegeEscalation: {
		Identifier:       TestSecConPrivilegeEscalation,
		Type:             normativeResult,
		Remediation:      SecConPrivilegeRemediation,
		ExceptionProcess: NoDocumentedProcess,
		Description: formDescription(TestSecConPrivilegeEscalation,
			`Checks if privileged escalation is enabled (AllowPrivilegeEscalation=true)`),
		BestPracticeReference: bestPracticeDocV1dot3URL + " Section 5.2",
	},
	TestContainerIsCertifiedIdentifier: {
		Identifier:  TestContainerIsCertifiedIdentifier,
		Type:        normativeResult,
		Remediation: ContainerIsCertifiedRemediation,
		Description: formDescription(TestContainerIsCertifiedIdentifier,
			`tests whether container images listed in the configuration file have passed the Red Hat Container Certification Program (CCP).`),
		BestPracticeReference: bestPracticeDocV1dot3URL + " Section 5.3.7",
		ExceptionProcess:      NoDocumentedProcess,
	},
	TestContainerHostPort: {
		Identifier:  TestContainerHostPort,
		Type:        informativeResult,
		Remediation: ContainerHostPortRemediation,
		Description: formDescription(TestContainerHostPort,
			`Verifies if containers define a hostPort.`),
		BestPracticeReference: bestPracticeDocV1dot3URL + " Section 5.3.6",
		ExceptionProcess:      NoDocumentedProcess,
	},
	TestPodHostNetwork: {
		Identifier:  TestPodHostNetwork,
		Type:        informativeResult,
		Remediation: PodHostNetworkRemediation,
		Description: formDescription(TestPodHostNetwork,
			`Verifies that the spec.HostNetwork parameter is set to false`),
		BestPracticeReference: bestPracticeDocV1dot3URL + " Section 5.3.6",
		ExceptionProcess:      NoDocumentedProcess,
	},
	TestPodHostPath: {
		Identifier:  TestPodHostPath,
		Type:        informativeResult,
		Remediation: PodHostNetworkRemediation,
		Description: formDescription(TestPodHostPath,
			`Verifies that the spec.HostPath parameter is not set (not present)`),
		BestPracticeReference: bestPracticeDocV1dot3URL + " Section 5.3.6",
		ExceptionProcess:      NoDocumentedProcess,
	},
	TestPodHostIPC: {
		Identifier:  TestPodHostIPC,
		Type:        informativeResult,
		Remediation: PodHostIPCRemediation,
		Description: formDescription(TestPodHostIPC,
			`Verifies that the spec.HostIpc parameter is set to false`),
		BestPracticeReference: bestPracticeDocV1dot3URL + " Section 5.3.6",
		ExceptionProcess:      NoDocumentedProcess,
	},
	TestPodHostPID: {
		Identifier:  TestPodHostPID,
		Type:        informativeResult,
		Remediation: PodHostPIDRemediation,
		Description: formDescription(TestPodHostPID,
			`Verifies that the spec.HostPid parameter is set to false`),
		BestPracticeReference: bestPracticeDocV1dot3URL + " Section 5.3.6",
		ExceptionProcess:      NoDocumentedProcess,
	},
	TestHugepagesNotManuallyManipulated: {
		Identifier:  TestHugepagesNotManuallyManipulated,
		Type:        normativeResult,
		Remediation: HugepagesNotManuallyManipulatedRemediation,
		Description: formDescription(TestHugepagesNotManuallyManipulated,
			`checks to see that HugePage settings have been configured through MachineConfig, and not manually on the
underlying Node.  This test case applies only to Nodes that are configured with the "worker" MachineConfigSet.  First,
the "worker" MachineConfig is polled, and the Hugepage settings are extracted.  Next, the underlying Nodes are polled
for configured HugePages through inspection of /proc/meminfo.  The results are compared, and the test passes only if
they are the same.`),
		BestPracticeReference: bestPracticeDocV1dot3URL + " Section 5.2",
		ExceptionProcess:      NoDocumentedProcess,
	},

	TestICMPv4ConnectivityIdentifier: {
		Identifier:  TestICMPv4ConnectivityIdentifier,
		Type:        normativeResult,
		Remediation: ICMPv4ConnectivityRemediation,
		Description: formDescription(TestICMPv4ConnectivityIdentifier,
			`checks that each CNF Container is able to communicate via ICMPv4 on the Default OpenShift network.  This
test case requires the Deployment of the debug daemonset.`),
		BestPracticeReference: bestPracticeDocV1dot3URL + " Section 5.2",
		ExceptionProcess:      NoDocumentedProcess,
	},

	TestICMPv6ConnectivityIdentifier: {
		Identifier:  TestICMPv6ConnectivityIdentifier,
		Type:        normativeResult,
		Remediation: ICMPv6ConnectivityRemediation,
		Description: formDescription(TestICMPv6ConnectivityIdentifier,
			`checks that each CNF Container is able to communicate via ICMPv6 on the Default OpenShift network.  This
test case requires the Deployment of the debug daemonset.`),
		BestPracticeReference: bestPracticeDocV1dot3URL + " Section 5.2",
		ExceptionProcess:      NoDocumentedProcess,
	},

	TestICMPv4ConnectivityMultusIdentifier: {
		Identifier:  TestICMPv4ConnectivityMultusIdentifier,
		Type:        normativeResult,
		Remediation: ICMPv4ConnectivityMultusRemediation,
		Description: formDescription(TestICMPv4ConnectivityMultusIdentifier,
			`checks that each CNF Container is able to communicate via ICMPv4 on the Multus network(s).  This
test case requires the Deployment of the debug daemonset.`),
		BestPracticeReference: bestPracticeDocV1dot3URL + " Section 5.2",
		ExceptionProcess:      NoDocumentedProcess,
	},

	TestICMPv6ConnectivityMultusIdentifier: {
		Identifier:  TestICMPv6ConnectivityMultusIdentifier,
		Type:        normativeResult,
		Remediation: ICMPv6ConnectivityMultusRemediation,
		Description: formDescription(TestICMPv6ConnectivityMultusIdentifier,
			`checks that each CNF Container is able to communicate via ICMPv6 on the Multus network(s).  This
test case requires the Deployment of the debug daemonset.`),
		BestPracticeReference: bestPracticeDocV1dot3URL + " Section 5.2",
		ExceptionProcess:      NoDocumentedProcess,
	},

	TestNamespaceBestPracticesIdentifier: {
		Identifier:  TestNamespaceBestPracticesIdentifier,
		Type:        normativeResult,
		Remediation: NamespaceBestPracticesRemediation,
		Description: formDescription(TestNamespaceBestPracticesIdentifier,
			`tests that all CNF's resources (PUTs and CRs) belong to valid namespaces. A valid namespace meets
the following conditions: (1) It was declared in the yaml config file under the targetNameSpaces
tag. (2) It doesn't have any of the following prefixes: default, openshift-, istio- and aspenmesh-`),
		BestPracticeReference: bestPracticeDocV1dot3URL + " Section 5.2, 16.3.8 and 16.3.9",
		ExceptionProcess:      NoDocumentedProcess,
	},

	TestNonTaintedNodeKernelsIdentifier: {
		Identifier:  TestNonTaintedNodeKernelsIdentifier,
		Type:        normativeResult,
		Remediation: NonTaintedNodeKernelsRemediation,
		Description: formDescription(TestNonTaintedNodeKernelsIdentifier,
			`ensures that the Node(s) hosting CNFs do not utilize tainted kernels. This test case is especially important
to support Highly Available CNFs, since when a CNF is re-instantiated on a backup Node, that Node's kernel may not have
the same hacks.'`),
		BestPracticeReference: bestPracticeDocV1dot3URL + " Section 5.2.14",
		ExceptionProcess:      NoDocumentedProcess,
	},

	TestOperatorInstallStatusSucceededIdentifier: {
		Identifier:  TestOperatorInstallStatusSucceededIdentifier,
		Type:        normativeResult,
		Remediation: OperatorInstallStatusSucceededRemediation,
		Description: formDescription(TestOperatorInstallStatusSucceededIdentifier,
			`Ensures that the target CNF operators report "Succeeded" as their installation status.`),
		BestPracticeReference: bestPracticeDocV1dot3URL + " Section 5.2.12 and 5.3.3",
		ExceptionProcess:      NoDocumentedProcess,
	},

	TestOperatorNoPrivileges: {
		Identifier:  TestOperatorNoPrivileges,
		Type:        normativeResult,
		Remediation: OperatorNoPrivilegesRemediation,
		Description: formDescription(TestOperatorNoPrivileges,
			`The operator is not installed with privileged rights. Test passes if clusterPermissions is not present in the CSV manifest or is present 
with no resourceNames under its rules.`),
		BestPracticeReference: bestPracticeDocV1dot3URL + " Section 5.2.12 and 5.3.3",
		ExceptionProcess:      NoDocumentedProcess,
	},

	TestOperatorIsCertifiedIdentifier: {
		Identifier:  TestOperatorIsCertifiedIdentifier,
		Type:        normativeResult,
		Remediation: OperatorIsCertifiedRemediation,
		Description: formDescription(TestOperatorIsCertifiedIdentifier,
			`tests whether CNF Operators listed in the configuration file have passed the Red Hat Operator Certification Program (OCP).`),
		BestPracticeReference: bestPracticeDocV1dot3URL + " Section 5.2.12 and 5.3.3",
		ExceptionProcess:      NoDocumentedProcess,
	},

	TestHelmIsCertifiedIdentifier: {
		Identifier:  TestHelmIsCertifiedIdentifier,
		Type:        normativeResult,
		Remediation: HelmIsCertifiedRemediation,
		Description: formDescription(TestHelmIsCertifiedIdentifier,
			`tests whether helm charts listed in the cluster passed the Red Hat Helm Certification Program.`),
		BestPracticeReference: bestPracticeDocV1dot3URL + " Section 5.2.12 and 5.3.3",
		ExceptionProcess:      NoDocumentedProcess,
	},

	TestOperatorIsInstalledViaOLMIdentifier: {
		Identifier:  TestOperatorIsInstalledViaOLMIdentifier,
		Type:        normativeResult,
		Remediation: OperatorIsInstalledViaOLMRemediation,
		Description: formDescription(TestOperatorIsInstalledViaOLMIdentifier,
			`tests whether a CNF Operator is installed via OLM.`),
		BestPracticeReference: bestPracticeDocV1dot3URL + " Section 5.2.12 and 5.3.3",
		ExceptionProcess:      NoDocumentedProcess,
	},

	TestPodNodeSelectorAndAffinityBestPractices: {
		Identifier:  TestPodNodeSelectorAndAffinityBestPractices,
		Type:        informativeResult,
		Remediation: PodNodeSelectorAndAffinityBestPracticesRemediation,
		Description: formDescription(TestPodNodeSelectorAndAffinityBestPractices,
			`ensures that CNF Pods do not specify nodeSelector or nodeAffinity.  In most cases, Pods should allow for
instantiation on any underlying Node.`),
		BestPracticeReference: bestPracticeDocV1dot3URL + " Section 5.2",
		ExceptionProcess:      NoDocumentedProcess,
	},

	TestPodHighAvailabilityBestPractices: {
		Identifier:  TestPodHighAvailabilityBestPractices,
		Type:        informativeResult,
		Remediation: PodHighAvailabilityBestPracticesRemediation,
		Description: formDescription(TestPodHighAvailabilityBestPractices,
			`ensures that CNF Pods specify podAntiAffinity rules and replica value is set to more than 1.`),
		BestPracticeReference: bestPracticeDocV1dot3URL + " Section 5.2",
		ExceptionProcess:      NoDocumentedProcess,
	},

	TestPodClusterRoleBindingsBestPracticesIdentifier: {
		Identifier:  TestPodClusterRoleBindingsBestPracticesIdentifier,
		Type:        normativeResult,
		Remediation: PodClusterRoleBindingsBestPracticesRemediation,
		Description: formDescription(TestPodClusterRoleBindingsBestPracticesIdentifier,
			`tests that a Pod does not specify ClusterRoleBindings.`),
		BestPracticeReference: bestPracticeDocV1dot3URL + " Section 5.2.10 and 5.3.6",
		ExceptionProcess:      NoDocumentedProcess,
	},

	TestPodDeploymentBestPracticesIdentifier: {
		Identifier:  TestPodDeploymentBestPracticesIdentifier,
		Type:        normativeResult,
		Remediation: PodDeploymentBestPracticesRemediation,
		Description: formDescription(TestPodDeploymentBestPracticesIdentifier,
			`tests that CNF Pod(s) are deployed as part of a ReplicaSet(s)/StatefulSet(s).`),
		BestPracticeReference: bestPracticeDocV1dot3URL + " Section 5.3.3 and 5.3.8",
		ExceptionProcess:      NoDocumentedProcess,
	},
	TestImagePullPolicyIdentifier: {
		Identifier:  TestImagePullPolicyIdentifier,
		Type:        normativeResult,
		Remediation: ImagePullPolicyRemediation,
		Description: formDescription(TestImagePullPolicyIdentifier,
			`Ensure that the containers under test are using IfNotPresent as Image Pull Policy..`),
		BestPracticeReference: bestPracticeDocV1dot3URL + "  Section 12.6",
		ExceptionProcess:      NoDocumentedProcess,
	},

	TestPodRoleBindingsBestPracticesIdentifier: {
		Identifier:  TestPodRoleBindingsBestPracticesIdentifier,
		Type:        normativeResult,
		Remediation: PodRoleBindingsBestPracticesRemediation,
		Description: formDescription(TestPodRoleBindingsBestPracticesIdentifier,
			`ensures that a CNF does not utilize RoleBinding(s) in a non-CNF Namespace.`),
		BestPracticeReference: bestPracticeDocV1dot3URL + " Section 5.3.3 and 5.3.5",
		ExceptionProcess:      NoDocumentedProcess,
	},

	TestPodServiceAccountBestPracticesIdentifier: {
		Identifier:  TestPodServiceAccountBestPracticesIdentifier,
		Type:        normativeResult,
		Remediation: PodServiceAccountBestPracticesRemediation,
		Description: formDescription(TestPodServiceAccountBestPracticesIdentifier,
			`tests that each CNF Pod utilizes a valid Service Account.`),
		BestPracticeReference: bestPracticeDocV1dot3URL + " Section 5.2.3 and 5.2.7",
		ExceptionProcess:      NoDocumentedProcess,
	},

	TestServicesDoNotUseNodeportsIdentifier: {
		Identifier:  TestServicesDoNotUseNodeportsIdentifier,
		Type:        normativeResult,
		Remediation: ServicesDoNotUseNodeportsRemediation,
		Description: formDescription(TestServicesDoNotUseNodeportsIdentifier,
			`tests that each CNF Service does not utilize NodePort(s).`),
		BestPracticeReference: bestPracticeDocV1dot3URL + " Section 5.3.1",
		ExceptionProcess:      NoDocumentedProcess,
	},

	TestUnalteredBaseImageIdentifier: {
		Identifier:  TestUnalteredBaseImageIdentifier,
		Type:        normativeResult,
		Remediation: UnalteredBaseImageRemediation,
		Description: formDescription(TestUnalteredBaseImageIdentifier,
			`ensures that the Container Base Image is not altered post-startup.  This test is a heuristic, and ensures
that there are no changes to the following directories:
1) /var/lib/rpm
2) /var/lib/dpkg
3) /bin
4) /sbin
5) /lib
6) /lib64
7) /usr/bin
8) /usr/sbin
9) /usr/lib
10) /usr/lib64`),
		BestPracticeReference: bestPracticeDocV1dot3URL + " Section 5.1.4",
		ExceptionProcess:      UnalteredBaseImageExceptionProcess,
	},

	TestUnalteredStartupBootParamsIdentifier: {
		Identifier:  TestUnalteredStartupBootParamsIdentifier,
		Type:        normativeResult,
		Remediation: UnalteredStartupBootParamsRemediation,
		Description: formDescription(TestUnalteredStartupBootParamsIdentifier,
			`tests that boot parameters are set through the MachineConfigOperator, and not set manually on the Node.`),
		BestPracticeReference: bestPracticeDocV1dot3URL + " Section 5.2.13 and 5.2.14",
		ExceptionProcess:      NoDocumentedProcess,
	},
	TestShutdownIdentifier: {
		Identifier: TestShutdownIdentifier,
		Type:       normativeResult,
		Description: formDescription(TestShutdownIdentifier,
			`Ensure that the containers lifecycle pre-stop management feature is configured.`),
		Remediation:           ShutdownRemediation,
		ExceptionProcess:      ShutdownExceptionProcess,
		BestPracticeReference: bestPracticeDocV1dot3URL + " Section 5.1.3, 12.2 and 12.5",
	},
	TestPodRecreationIdentifier: {
		Identifier: TestPodRecreationIdentifier,
		Type:       normativeResult,
		Description: formDescription(TestPodRecreationIdentifier,
			`tests that a CNF is configured to support High Availability.  
			First, this test cordons and drains a Node that hosts the CNF Pod.  
			Next, the test ensures that OpenShift can re-instantiate the Pod on another Node, 
			and that the actual replica count matches the desired replica count.`),
		Remediation:           PodRecreationRemediation,
		BestPracticeReference: bestPracticeDocV1dot3URL + " Section 5.2",
		ExceptionProcess:      NoDocumentedProcess,
	},
	TestSysctlConfigsIdentifier: {
		Identifier: TestSysctlConfigsIdentifier,
		Type:       normativeResult,
		Description: formDescription(TestPodRecreationIdentifier,
			`tests that no one has changed the node's sysctl configs after the node
			was created, the tests works by checking if the sysctl configs are consistent with the
			MachineConfig CR which defines how the node should be configured`),
		Remediation:           SysctlConfigsRemediation,
		BestPracticeReference: bestPracticeDocV1dot3URL + " Section 5.2",
		ExceptionProcess:      NoDocumentedProcess,
	},
	TestServiceMeshIdentifier: {
		Identifier: TestServiceMeshIdentifier,
		Type:       normativeResult,
		Description: formDescription(TestServiceMeshIdentifier,
			`verifies whether, if available, service mesh is actually being used by the CNF pods`),
		Remediation:           ServiceMeshRemediation,
		BestPracticeReference: bestPracticeDocV1dot3URL + " Section 5.2",
		ExceptionProcess:      NoDocumentedProcess,
	},
	TestScalingIdentifier: {
		Identifier: TestScalingIdentifier,
		Type:       normativeResult,
		Description: formDescription(TestScalingIdentifier,
			`tests that CNF deployments support scale in/out operations. 
			First, The test starts getting the current replicaCount (N) of the deployment/s with the Pod Under Test. Then, it executes the 
			scale-in oc command for (N-1) replicas. Lastly, it executes the scale-out oc command, restoring the original replicaCount of the deployment/s.`),
		Remediation:           ScalingRemediation,
		BestPracticeReference: bestPracticeDocV1dot3URL + " Section 5.2",
		ExceptionProcess:      NoDocumentedProcess,
	},
	TestIsRedHatReleaseIdentifier: {
		Identifier: TestIsRedHatReleaseIdentifier,
		Type:       normativeResult,
		Description: formDescription(TestIsRedHatReleaseIdentifier,
			`verifies if the container base image is redhat.`),
		Remediation:           IsRedHatReleaseRemediation,
		ExceptionProcess:      IsRedHatReleaseExceptionProcess,
		BestPracticeReference: bestPracticeDocV1dot3URL + " Section 5.2",
	},
	TestIsSELinuxEnforcingIdentifier: {
		Identifier: TestIsSELinuxEnforcingIdentifier,
		Type:       normativeResult,
		Description: formDescription(TestIsSELinuxEnforcingIdentifier,
			`verifies that all openshift platform/cluster nodes have selinux in "Enforcing" mode.`),
		Remediation:           IsSELinuxEnforcingRemediation,
		BestPracticeReference: bestPracticeDocV1dot3URL + " Section 10.3 Pod Security",
		ExceptionProcess:      NoDocumentedProcess,
	},
	TestUndeclaredContainerPortsUsage: {
		Identifier: TestUndeclaredContainerPortsUsage,
		Type:       normativeResult,
		Description: formDescription(TestUndeclaredContainerPortsUsage,
			`check that containers do not listen on ports that weren't declared in their specification`),
		Remediation:           UndeclaredContainerPortsRemediation,
		BestPracticeReference: bestPracticeDocV1dot3URL + " Section 16.3.1.1",
		ExceptionProcess:      NoDocumentedProcess,
	},
	TestOCPReservedPortsUsage: {
		Identifier: TestOCPReservedPortsUsage,
		Type:       normativeResult,
		Description: formDescription(TestOCPReservedPortsUsage,
			`check that containers do not listen on ports that are reserved by Openshift`),
		Remediation:           OCPReservedPortsUsageRemediation,
		BestPracticeReference: bestPracticeDocV1dot4URL + " Section 3.5.9",
		ExceptionProcess:      NoDocumentedProcess,
	},
	TestCrdsStatusSubresourceIdentifier: {
		Identifier: TestCrdsStatusSubresourceIdentifier,
		Type:       informativeResult,
		Description: formDescription(TestCrdsStatusSubresourceIdentifier,
			`checks that all CRDs have a status subresource specification (Spec.versions[].Schema.OpenAPIV3Schema.Properties["status"]).`),
		Remediation:           CrdsStatusSubresourceRemediation,
		BestPracticeReference: bestPracticeDocV1dot3URL + " Section 5.2",
		ExceptionProcess:      NoDocumentedProcess,
	},
	TestLoggingIdentifier: {
		Identifier: TestLoggingIdentifier,
		Type:       informativeResult,
		Description: formDescription(TestLoggingIdentifier,
			`check that all containers under test use standard input output and standard error when logging`),
		Remediation:           LoggingRemediation,
		BestPracticeReference: bestPracticeDocV1dot3URL + " Section 10.1",
		ExceptionProcess:      NoDocumentedProcess,
	},
	TestTerminationMessagePolicyIdentifier: {
		Identifier: TestTerminationMessagePolicyIdentifier,
		Type:       informativeResult,
		Description: formDescription(TestTerminationMessagePolicyIdentifier,
			`check that all containers are using terminationMessagePolicy: FallbackToLogsOnError`),
		Remediation:           TerminationMessagePolicyRemediation,
		BestPracticeReference: bestPracticeDocV1dot3URL + " Section 12.1",
		ExceptionProcess:      NoDocumentedProcess,
	},
	TestPodAutomountServiceAccountIdentifier: {
		Identifier: TestPodAutomountServiceAccountIdentifier,
		Type:       normativeResult,
		Description: formDescription(TestPodAutomountServiceAccountIdentifier,
			`check that all pods under test have automountServiceAccountToken set to false`),
		Remediation:           AutomountServiceTokenRemediation,
		ExceptionProcess:      AutomountServiceTokenExceptionProcess,
		BestPracticeReference: bestPracticeDocV1dot3URL + " Section 12.7",
	},
	TestLivenessProbeIdentifier: {
		Identifier: TestLivenessProbeIdentifier,
		Type:       normativeResult,
		Description: formDescription(TestLivenessProbeIdentifier, `check that all containers under test
		have liveness probe defined`),
		Remediation:           LivenessProbeRemediation,
		BestPracticeReference: bestPracticeDocV1dot3URL + " Section 5.2.16, 12.1 and 12.5",
		ExceptionProcess:      NoDocumentedProcess,
	},
	TestReadinessProbeIdentifier: {
		Identifier: TestReadinessProbeIdentifier,
		Type:       normativeResult,
		Description: formDescription(TestReadinessProbeIdentifier, `check that all containers under test
		have readiness probe defined`),
		Remediation:           ReadinessProbeRemediation,
		BestPracticeReference: bestPracticeDocV1dot3URL + " Section 5.2.16, 12.1 and 12.5",
		ExceptionProcess:      NoDocumentedProcess,
	},
	TestOneProcessPerContainerIdentifier: {
		Identifier: TestOneProcessPerContainerIdentifier,
		Type:       informativeResult,
		Description: formDescription(TestOneProcessPerContainerIdentifier, `check that all containers under test
		have only one process running`),
		Remediation:           OneProcessPerContainerRemediation,
		BestPracticeReference: bestPracticeDocV1dot3URL + " Section 10.8.3",
		ExceptionProcess:      NoDocumentedProcess,
	},
	TestSYSNiceRealtimeCapabilityIdentifier: {
		Identifier:            TestSYSNiceRealtimeCapabilityIdentifier,
		Type:                  informativeResult,
		Description:           formDescription(TestSYSNiceRealtimeCapabilityIdentifier, `Check that pods running on nodes with realtime kernel enabled have the SYS_NICE capability enabled in their spec.`),
		Remediation:           SYSNiceRealtimeCapabilityRemediation,
		BestPracticeReference: bestPracticeDocV1dot3URL + " Section 2.7.4",
		ExceptionProcess:      NoDocumentedProcess,
	},
	TestSysPtraceCapabilityIdentifier: {
		Identifier:            TestSysPtraceCapabilityIdentifier,
		Type:                  informativeResult,
		Description:           formDescription(TestSysPtraceCapabilityIdentifier, `check that if process namespace sharing is enabled for a Pod then the SYS_PTRACE capability is allowed`),
		Remediation:           SysPtraceCapabilityRemediation,
		BestPracticeReference: bestPracticeDocV1dot3URL + " Section 2.7.5",
		ExceptionProcess:      NoDocumentedProcess,
	},
	TestPodRequestsAndLimitsIdentifier: {
		Identifier:            TestPodRequestsAndLimitsIdentifier,
		Type:                  informativeResult,
		Description:           formDescription(TestPodRequestsAndLimitsIdentifier, `Check that containers have resource requests and limits specified in their spec.`),
		Remediation:           RequestsAndLimitsRemediation,
		BestPracticeReference: bestPracticeDocV1dot4URL + " Section 4.6.11",
	},
	TestNamespaceResourceQuotaIdentifier: {
		Identifier:            TestNamespaceResourceQuotaIdentifier,
		Type:                  informativeResult,
		Description:           formDescription(TestNamespaceResourceQuotaIdentifier, `Checks to see if CNF workload pods are running in namespaces that have resource quotas applied.`),
		Remediation:           NamespaceResourceQuotaRemediation,
		BestPracticeReference: bestPracticeDocV1dot3URL + " Section 4.6.8", // TODO Change this to v1.4 when available
	},
	TestNamespaceLimitRangeIdentifier: {
		Identifier:            TestNamespaceLimitRangeIdentifier,
		Type:                  informativeResult,
		Description:           formDescription(TestNamespaceLimitRangeIdentifier, `Checks to see if CNF workload pods are running in namespaces that have limit ranges applied.`),
		Remediation:           NamespaceLimitRangeRemediation,
		BestPracticeReference: bestPracticeDocV1dot3URL + " Section 4.6.8", // TODO Change this to v1.4 when available
	},
}
