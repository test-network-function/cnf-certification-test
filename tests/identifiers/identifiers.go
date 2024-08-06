// Copyright (C) 2021-2024 Red Hat, Inc.
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
	"strings"

	"github.com/test-network-function/cnf-certification-test/tests/common"
	"github.com/test-network-function/test-network-function-claim/pkg/claim"
)

// shared description text.
const (
	iptablesNftablesImplicitCheck = `Note: this test also ensures iptables and nftables are not configured by workload pods:
- NET_ADMIN and NET_RAW are required to modify nftables (namespaced) which is not desired inside pods.
nftables should be configured by an administrator outside the scope of the workload. nftables are usually configured
by operators, for instance the Performance Addon Operator (PAO) or istio.
- Privileged container are required to modify host iptables, which is not safe to perform inside pods. nftables
should be configured by an administrator outside the scope of the workload. iptables are usually configured by operators,
for instance the Performance Addon Operator (PAO) or istio.`
)

const (
	TagCommon    = "common"
	TagExtended  = "extended"
	TagTelco     = "telco"
	TagFarEdge   = "faredge"
	FarEdge      = "FarEdge"
	Telco        = "Telco"
	NonTelco     = "NonTelco"
	Extended     = "Extended"
	Optional     = "Optional"
	Mandatory    = "Mandatory"
	TagPreflight = "preflight"
)

const (
	NotApplicableSNO = ` Not applicable to SNO applications.`
)

func init() {
	InitCatalog()
}

func AddCatalogEntry(testID, suiteName, description, remediation, exception, reference string, qe bool, categoryclassification map[string]string, tags ...string) (aID claim.Identifier) {
	// Default Values (if missing).
	if strings.TrimSpace(exception) == "" {
		exception = NoDocumentedProcess
	}
	if strings.TrimSpace(reference) == "" {
		reference = "No Reference Document Specified"
	}
	if len(tags) == 0 {
		tags = append(tags, TagCommon)
	}

	tcDescription, aID := claim.BuildTestCaseDescription(testID, suiteName, description, remediation, exception, reference, qe, categoryclassification, tags...)
	Catalog[aID] = tcDescription
	Classification[aID.Id] = categoryclassification

	return aID
}

var (
	TestICMPv4ConnectivityIdentifier                  claim.Identifier
	TestNetworkPolicyDenyAllIdentifier                claim.Identifier
	Test1337UIDIdentifier                             claim.Identifier
	TestContainerIsCertifiedDigestIdentifier          claim.Identifier
	TestHelmVersionIdentifier                         claim.Identifier
	TestPodHugePages2M                                claim.Identifier
	TestPodHugePages1G                                claim.Identifier
	TestHyperThreadEnable                             claim.Identifier
	TestReservedExtendedPartnerPorts                  claim.Identifier
	TestAffinityRequiredPods                          claim.Identifier
	TestContainerPostStartIdentifier                  claim.Identifier
	TestContainerPrestopIdentifier                    claim.Identifier
	TestDpdkCPUPinningExecProbe                       claim.Identifier
	TestSysAdminIdentifier                            claim.Identifier
	TestNetAdminIdentifier                            claim.Identifier
	TestNetRawIdentifier                              claim.Identifier
	TestIpcLockIdentifier                             claim.Identifier
	TestBpfIdentifier                                 claim.Identifier
	TestStorageProvisioner                            claim.Identifier
	TestExclusiveCPUPoolIdentifier                    claim.Identifier
	TestSharedCPUPoolSchedulingPolicy                 claim.Identifier
	TestExclusiveCPUPoolSchedulingPolicy              claim.Identifier
	TestIsolatedCPUPoolSchedulingPolicy               claim.Identifier
	TestRtAppNoExecProbes                             claim.Identifier
	TestRestartOnRebootLabelOnPodsUsingSRIOV          claim.Identifier
	TestSecConNonRootUserIdentifier                   claim.Identifier
	TestSecContextIdentifier                          claim.Identifier
	TestSecConPrivilegeEscalation                     claim.Identifier
	TestContainerHostPort                             claim.Identifier
	TestPodHostNetwork                                claim.Identifier
	TestPodHostPath                                   claim.Identifier
	TestPodHostIPC                                    claim.Identifier
	TestPodHostPID                                    claim.Identifier
	TestHugepagesNotManuallyManipulated               claim.Identifier
	TestICMPv6ConnectivityIdentifier                  claim.Identifier
	TestICMPv4ConnectivityMultusIdentifier            claim.Identifier
	TestICMPv6ConnectivityMultusIdentifier            claim.Identifier
	TestServiceDualStackIdentifier                    claim.Identifier
	TestNamespaceBestPracticesIdentifier              claim.Identifier
	TestNonTaintedNodeKernelsIdentifier               claim.Identifier
	TestOperatorInstallStatusSucceededIdentifier      claim.Identifier
	TestOperatorNoSCCAccess                           claim.Identifier
	TestOperatorIsCertifiedIdentifier                 claim.Identifier
	TestHelmIsCertifiedIdentifier                     claim.Identifier
	TestOperatorIsInstalledViaOLMIdentifier           claim.Identifier
	TestOperatorHasSemanticVersioningIdentifier       claim.Identifier
	TestOperatorReadOnlyFilesystem                    claim.Identifier
	TestOperatorAutomountTokens                       claim.Identifier
	TestOperatorRunAsNonRoot                          claim.Identifier
	TestOperatorRunAsUserID                           claim.Identifier
	TestOperatorCrdVersioningIdentifier               claim.Identifier
	TestOperatorCrdSchemaIdentifier                   claim.Identifier
	TestOperatorSingleCrdOwnerIdentifier              claim.Identifier
	TestPodNodeSelectorAndAffinityBestPractices       claim.Identifier
	TestPodHighAvailabilityBestPractices              claim.Identifier
	TestPodClusterRoleBindingsBestPracticesIdentifier claim.Identifier
	TestPodDeploymentBestPracticesIdentifier          claim.Identifier
	TestDeploymentScalingIdentifier                   claim.Identifier
	TestStateFulSetScalingIdentifier                  claim.Identifier
	TestImagePullPolicyIdentifier                     claim.Identifier
	TestPodRecreationIdentifier                       claim.Identifier
	TestPodRoleBindingsBestPracticesIdentifier        claim.Identifier
	TestPodServiceAccountBestPracticesIdentifier      claim.Identifier
	TestPodAutomountServiceAccountIdentifier          claim.Identifier
	TestServicesDoNotUseNodeportsIdentifier           claim.Identifier
	TestUnalteredBaseImageIdentifier                  claim.Identifier
	TestUnalteredStartupBootParamsIdentifier          claim.Identifier
	TestLoggingIdentifier                             claim.Identifier
	TestTerminationMessagePolicyIdentifier            claim.Identifier
	TestCrdsStatusSubresourceIdentifier               claim.Identifier
	TestSysctlConfigsIdentifier                       claim.Identifier
	TestServiceMeshIdentifier                         claim.Identifier
	TestOCPLifecycleIdentifier                        claim.Identifier
	TestNodeOperatingSystemIdentifier                 claim.Identifier
	TestIsRedHatReleaseIdentifier                     claim.Identifier
	TestIsSELinuxEnforcingIdentifier                  claim.Identifier
	TestUndeclaredContainerPortsUsage                 claim.Identifier
	TestOCPReservedPortsUsage                         claim.Identifier
	TestLivenessProbeIdentifier                       claim.Identifier
	TestReadinessProbeIdentifier                      claim.Identifier
	TestStartupProbeIdentifier                        claim.Identifier
	TestOneProcessPerContainerIdentifier              claim.Identifier
	TestSYSNiceRealtimeCapabilityIdentifier           claim.Identifier
	TestSysPtraceCapabilityIdentifier                 claim.Identifier
	TestPodRequestsAndLimitsIdentifier                claim.Identifier
	TestNamespaceResourceQuotaIdentifier              claim.Identifier
	TestPodDisruptionBudgetIdentifier                 claim.Identifier
	TestPodTolerationBypassIdentifier                 claim.Identifier
	TestPersistentVolumeReclaimPolicyIdentifier       claim.Identifier
	TestContainersImageTag                            claim.Identifier
	TestNoSSHDaemonsAllowedIdentifier                 claim.Identifier
	TestCPUIsolationIdentifier                        claim.Identifier
	TestContainerPortNameFormat                       claim.Identifier
	TestCrdScalingIdentifier                          claim.Identifier
	TestCrdRoleIdentifier                             claim.Identifier
	TestLimitedUseOfExecProbesIdentifier              claim.Identifier
	// Chaos Testing.
	// TestPodDeleteIdentifier claim.Identifier.
)

//nolint:funlen
func InitCatalog() map[claim.Identifier]claim.TestCaseDescription {
	TestNetworkPolicyDenyAllIdentifier = AddCatalogEntry(
		"network-policy-deny-all",
		common.NetworkingTestKey,
		`Check that network policies attached to namespaces running workload pods contain a default deny-all rule for both ingress and egress traffic`,
		NetworkPolicyDenyAllRemediation,
		NoExceptionProcessForExtendedTests,
		TestNetworkPolicyDenyAllIdentifierDocLink,
		true,
		map[string]string{
			FarEdge:  Optional,
			Telco:    Optional,
			NonTelco: Optional,
			Extended: Optional,
		},
		TagCommon)

	Test1337UIDIdentifier = AddCatalogEntry(
		"no-1337-uid",
		common.AccessControlTestKey,
		`Checks that all pods are not using the securityContext UID 1337`,
		UID1337Remediation,
		NoExceptionProcessForExtendedTests,
		Test1337UIDIdentifierDocLink,
		true,
		map[string]string{
			FarEdge:  Optional,
			Telco:    Optional,
			NonTelco: Optional,
			Extended: Mandatory,
		},
		TagExtended)

	TestLimitedUseOfExecProbesIdentifier = AddCatalogEntry(
		"max-resources-exec-probes",
		common.PerformanceTestKey,
		`Checks that less than 10 exec probes are configured in the cluster for this workload. Also checks that the periodSeconds parameter for each probe is superior or equal to 10.`,
		LimitedUseOfExecProbesRemediation,
		NoDocumentedProcess,
		TestLimitedUseOfExecProbesIdentifierDocLink,
		true,
		map[string]string{
			FarEdge:  Optional,
			Telco:    Optional,
			NonTelco: Optional,
			Extended: Optional},
		TagFarEdge)

	TestHelmVersionIdentifier = AddCatalogEntry(
		"helm-version",
		common.AffiliatedCertTestKey,
		`Test to check if the helm chart is v3`,
		HelmVersionV3Remediation,
		NoDocumentedProcess,
		TestHelmVersionIdentifierDocLink,
		true,
		map[string]string{
			FarEdge:  Mandatory,
			Telco:    Mandatory,
			NonTelco: Mandatory,
			Extended: Mandatory,
		},
		TagCommon)

	// TestContainerIsCertifiedDigestIdentifier tests whether the container has passed Container Certification.
	TestContainerIsCertifiedDigestIdentifier = AddCatalogEntry(
		"container-is-certified-digest",
		common.AffiliatedCertTestKey,
		`Tests whether container images that are autodiscovered have passed the Red Hat Container Certification Program by their digest(CCP).`,
		ContainerIsCertifiedDigestRemediation,
		AffiliatedCert,
		TestContainerIsCertifiedDigestIdentifierDocLink,
		true,
		map[string]string{
			FarEdge:  Mandatory,
			Telco:    Mandatory,
			NonTelco: Mandatory,
			Extended: Mandatory,
		},
		TagCommon)

	TestPodHugePages2M = AddCatalogEntry(
		"hugepages-2m-only",
		common.PlatformAlterationTestKey,
		`Check that pods using hugepages only use 2Mi size`,
		PodHugePages2MRemediation,
		NoExceptionProcessForExtendedTests,
		TestPodHugePages2MDocLink,
		true,
		map[string]string{
			FarEdge:  Optional,
			Telco:    Optional,
			NonTelco: Optional,
			Extended: Mandatory,
		},
		TagExtended)

	TestPodHugePages1G = AddCatalogEntry(
		"hugepages-1g-only",
		common.PlatformAlterationTestKey,
		`Check that pods using hugepages only use 1Gi size`,
		PodHugePages1GRemediation,
		NoDocumentedProcess,
		TestPodHugePages1GDocLink,
		true,
		map[string]string{
			FarEdge:  Mandatory,
			Telco:    Optional,
			NonTelco: Optional,
			Extended: Optional,
		},
		TagFarEdge)

	TestHyperThreadEnable = AddCatalogEntry(
		"hyperthread-enable",
		common.PlatformAlterationTestKey,
		`Check that baremetal workers have hyperthreading enabled`,
		HyperThreadEnable,
		NoDocumentedProcess,
		TestHyperThreadEnableDocLink,
		false,
		map[string]string{
			FarEdge:  Optional,
			Telco:    Optional,
			NonTelco: Optional,
			Extended: Optional,
		},
		TagExtended)

	TestReservedExtendedPartnerPorts = AddCatalogEntry(
		"reserved-partner-ports",
		common.NetworkingTestKey,
		`Checks that pods and containers are not consuming ports designated as reserved by partner`,
		ReservedPartnerPortsRemediation,
		NoExceptionProcessForExtendedTests,
		TestReservedExtendedPartnerPortsDocLink,
		true,
		map[string]string{
			FarEdge:  Optional,
			Telco:    Optional,
			NonTelco: Optional,
			Extended: Mandatory,
		},
		TagExtended)

	TestAffinityRequiredPods = AddCatalogEntry(
		"affinity-required-pods",
		common.LifecycleTestKey,
		`Checks that affinity rules are in place if AffinityRequired: 'true' labels are set on Pods.`,
		AffinityRequiredRemediation,
		NoDocumentedProcess,
		TestAffinityRequiredPodsDocLink,
		true,
		map[string]string{
			FarEdge:  Mandatory,
			Telco:    Mandatory,
			NonTelco: Optional,
			Extended: Mandatory,
		},
		TagTelco)

	TestStorageProvisioner = AddCatalogEntry(
		"storage-provisioner",
		common.LifecycleTestKey,
		`Checks that pods do not place persistent volumes on local storage in multinode clusters. Local storage is recommended for single node clusters, but only one type of local storage should be installed (lvms or noprovisioner).`,
		CheckStorageProvisionerRemediation,
		NoExceptions,
		TestStorageProvisionerDocLink,
		true,
		map[string]string{
			FarEdge:  Mandatory,
			Telco:    Mandatory,
			NonTelco: Mandatory,
			Extended: Mandatory,
		},
		TagCommon)

	TestContainerPostStartIdentifier = AddCatalogEntry(
		"container-poststart",
		common.LifecycleTestKey,
		`Ensure that the containers lifecycle postStart management feature is configured. A container must receive important events from the platform and conform/react to these events properly. For example, a container should catch SIGTERM or SIGKILL from the platform and shutdown as quickly as possible. Other typically important events from the platform are PostStart to initialize before servicing requests and PreStop to release resources cleanly before shutting down.`,                                                                                                                                                                                                                           //nolint:lll
		`PostStart is normally used to configure the container, set up dependencies, and record the new creation. You could use this event to check that a required API is available before the container’s main work begins. Kubernetes will not change the container’s state to Running until the PostStart script has executed successfully. For details, see https://www.containiq.com/post/kubernetes-container-lifecycle-events-and-hooks and https://kubernetes.io/docs/concepts/containers/container-lifecycle-hooks. PostStart is used to configure container, set up dependencies, record new creation. It can also be used to check that a required API is available before the container’s work begins.`, //nolint:lll
		ContainerPostStartIdentifierRemediation,
		TestContainerPostStartIdentifierDocLink,
		true,
		map[string]string{
			FarEdge:  Mandatory,
			Telco:    Mandatory,
			NonTelco: Optional,
			Extended: Mandatory,
		},
		TagTelco)

	TestContainerPrestopIdentifier = AddCatalogEntry(
		"container-prestop",
		common.LifecycleTestKey,
		`Ensure that the containers lifecycle preStop management feature is configured. The most basic requirement for the lifecycle management of Pods in OpenShift are the ability to start and stop correctly. There are different ways a pod can stop on an OpenShift cluster. One way is that the pod can remain alive but non-functional. Another way is that the pod can crash and become non-functional. When pods are shut down by the platform they are sent a SIGTERM signal which means that the process in the container should start shutting down, closing connections and stopping all activity. If the pod doesn’t shut down within the default 30 seconds then the platform may send a SIGKILL signal which will stop the pod immediately. This method isn’t as clean and the default time between the SIGTERM and SIGKILL messages can be modified based on the requirements of the application. Containers should respond to SIGTERM/SIGKILL with graceful shutdown.`, //nolint:lll
		`The preStop can be used to gracefully stop the container and clean resources (e.g., DB connection). For details, see https://www.containiq.com/post/kubernetes-container-lifecycle-events-and-hooks and https://kubernetes.io/docs/concepts/containers/container-lifecycle-hooks. All pods must respond to SIGTERM signal and shutdown gracefully with a zero exit code.`, //nolint:lll
		ContainerPrestopIdentifierRemediation,
		TestContainerPrestopIdentifierDocLink,
		true,
		map[string]string{
			FarEdge:  Mandatory,
			Telco:    Mandatory,
			NonTelco: Optional,
			Extended: Mandatory,
		},
		TagTelco)

	TestDpdkCPUPinningExecProbe = AddCatalogEntry(
		"dpdk-cpu-pinning-exec-probe",
		common.NetworkingTestKey,
		`If a workload is doing CPU pinning, exec probes may not be used.`,
		DpdkCPUPinningExecProbeRemediation,
		NoDocumentedProcess,
		TestDpdkCPUPinningExecProbeDocLink,
		true,
		map[string]string{
			FarEdge:  Mandatory,
			Telco:    Mandatory,
			NonTelco: Optional,
			Extended: Mandatory,
		},
		TagTelco)

	TestNetAdminIdentifier = AddCatalogEntry(
		"net-admin-capability-check",
		common.AccessControlTestKey,
		`Ensures that containers do not use NET_ADMIN capability. `+iptablesNftablesImplicitCheck,
		SecConRemediation,
		`Exception will be considered for user plane or networking functions (e.g. SR-IOV, Multicast). Must identify which container requires the capability and detail why.`,
		TestNetAdminIdentifierDocLink,
		true,
		map[string]string{
			FarEdge:  Mandatory,
			Telco:    Mandatory,
			NonTelco: Optional,
			Extended: Mandatory,
		},
		TagTelco)

	TestSysAdminIdentifier = AddCatalogEntry(
		"sys-admin-capability-check",
		common.AccessControlTestKey,
		`Ensures that containers do not use SYS_ADMIN capability`,
		SecConRemediation+" Containers should not use the SYS_ADMIN Linux capability.",
		NoExceptions,
		TestSysAdminIdentifierDocLink,
		true,
		map[string]string{
			FarEdge:  Mandatory,
			Telco:    Mandatory,
			NonTelco: Mandatory,
			Extended: Mandatory,
		},
		TagCommon)

	TestIpcLockIdentifier = AddCatalogEntry(
		"ipc-lock-capability-check",
		common.AccessControlTestKey,
		`Ensures that containers do not use IPC_LOCK capability. Workloads should avoid accessing host resources - spec.HostIpc should be false.`,
		SecConRemediation,
		`Exception possible if a workload uses mlock(), mlockall(), shmctl(), mmap(); exception will be considered for DPDK applications. Must identify which container requires the capability and detail why.`,
		TestIpcLockIdentifierDocLink,
		true,
		map[string]string{
			FarEdge:  Mandatory,
			Telco:    Mandatory,
			NonTelco: Optional,
			Extended: Mandatory,
		},
		TagTelco)

	TestNetRawIdentifier = AddCatalogEntry(
		"net-raw-capability-check",
		common.AccessControlTestKey,
		`Ensures that containers do not use NET_RAW capability. `+iptablesNftablesImplicitCheck,
		SecConRemediation,
		`Exception will be considered for user plane or networking functions. Must identify which container requires the capability and detail why.`,
		TestNetRawIdentifierDocLink,
		true,
		map[string]string{
			FarEdge:  Mandatory,
			Telco:    Mandatory,
			NonTelco: Optional,
			Extended: Mandatory,
		},
		TagTelco)

	TestBpfIdentifier = AddCatalogEntry(
		"bpf-capability-check",
		common.AccessControlTestKey,
		`Ensures that containers do not use BPF capability. Workloads should avoid loading eBPF filters`,
		BpfCapabilityRemediation,
		`Exception can be considered. Must identify which container requires the capability and detail why.`,
		TestBpfIdentifierDocLink,
		true,
		map[string]string{
			FarEdge:  Mandatory,
			Telco:    Mandatory,
			NonTelco: Optional,
			Extended: Mandatory,
		},
		TagTelco)

	TestExclusiveCPUPoolIdentifier = AddCatalogEntry(
		"exclusive-cpu-pool",
		common.PerformanceTestKey,
		`Ensures that if one container in a Pod selects an exclusive CPU pool the rest select the same type of CPU pool`,
		ExclusiveCPUPoolRemediation,
		NoDocumentedProcess,
		TestExclusiveCPUPoolIdentifierDocLink,
		true,
		map[string]string{
			FarEdge:  Mandatory,
			Telco:    Optional,
			NonTelco: Optional,
			Extended: Optional,
		},
		TagFarEdge)

	TestSharedCPUPoolSchedulingPolicy = AddCatalogEntry(
		"shared-cpu-pool-non-rt-scheduling-policy",
		common.PerformanceTestKey,
		`Ensures that if application workload runs in shared CPU pool, it chooses non-RT CPU schedule policy to always share the CPU with other applications and kernel threads.`,
		SharedCPUPoolSchedulingPolicyRemediation,
		NoDocumentedProcess,
		TestSharedCPUPoolSchedulingPolicyDocLink,
		true,
		map[string]string{
			FarEdge:  Mandatory,
			Telco:    Optional,
			NonTelco: Optional,
			Extended: Optional,
		},
		TagFarEdge)

	TestExclusiveCPUPoolSchedulingPolicy = AddCatalogEntry(
		"exclusive-cpu-pool-rt-scheduling-policy",
		common.PerformanceTestKey,
		`Ensures that if application workload runs in exclusive CPU pool, it chooses RT CPU schedule policy and set the priority less than 10.`,
		ExclusiveCPUPoolSchedulingPolicyRemediation,
		NoDocumentedProcess,
		TestExclusiveCPUPoolSchedulingPolicyDocLink,
		true,
		map[string]string{
			FarEdge:  Mandatory,
			Telco:    Optional,
			NonTelco: Optional,
			Extended: Optional,
		},
		TagFarEdge)

	TestIsolatedCPUPoolSchedulingPolicy = AddCatalogEntry(
		"isolated-cpu-pool-rt-scheduling-policy",
		common.PerformanceTestKey,
		`Ensures that a workload running in an application-isolated exclusive CPU pool selects a RT CPU scheduling policy`,
		IsolatedCPUPoolSchedulingPolicyRemediation,
		NoDocumentedProcess,
		TestIsolatedCPUPoolSchedulingPolicyDocLink,
		true,
		map[string]string{
			FarEdge:  Mandatory,
			Telco:    Optional,
			NonTelco: Optional,
			Extended: Optional,
		},
		TagFarEdge)

	TestRtAppNoExecProbes = AddCatalogEntry(
		"rt-apps-no-exec-probes",
		common.PerformanceTestKey,
		`Ensures that if one container runs a real time application exec probes are not used`,
		RtAppNoExecProbesRemediation,
		NoDocumentedProcess,
		TestRtAppNoExecProbesDocLink,
		true,
		map[string]string{
			FarEdge:  Mandatory,
			Telco:    Optional,
			NonTelco: Optional,
			Extended: Optional,
		},
		TagFarEdge)

	TestRestartOnRebootLabelOnPodsUsingSRIOV = AddCatalogEntry(
		"restart-on-reboot-sriov-pod",
		common.NetworkingTestKey,
		`Ensures that the label restart-on-reboot exists on pods that use SRIOV network interfaces.`,
		SRIOVPodsRestartOnRebootLabelRemediation,
		NoDocumentedProcess,
		TestRestartOnRebootLabelOnPodsUsingSRIOVDocLink,
		false,
		map[string]string{
			FarEdge:  Mandatory,
			Telco:    Optional,
			NonTelco: Optional,
			Extended: Mandatory,
		},
		TagFarEdge)

	TestSecConNonRootUserIdentifier = AddCatalogEntry(
		"security-context-non-root-user-check",
		common.AccessControlTestKey,
		`Checks the security context runAsUser parameter in pods and containers to make sure it is not set to uid root(0). Pods and containers should not run as root (runAsUser is not set to uid0).`,
		SecConNonRootUserRemediation,
		SecConNonRootUserExceptionProcess,
		TestSecConNonRootUserIdentifierDocLink,
		true,
		map[string]string{
			FarEdge:  Mandatory,
			Telco:    Mandatory,
			NonTelco: Mandatory,
			Extended: Mandatory,
		},
		TagCommon)

	TestSecContextIdentifier = AddCatalogEntry(
		"security-context",
		common.AccessControlTestKey,
		`Checks the security context matches one of the 4 categories`,
		`Exception possible if a workload uses mlock(), mlockall(), shmctl(), mmap(); exception will be considered for DPDK applications. Must identify which container requires the capability and document why. If the container had the right configuration of the allowed category from the 4 approved list then the test will pass. The 4 categories are defined in Requirement ID 94118 [here](#security-context-categories)`, //nolint:lll
		`no exception needed for optional/extended test`,
		TestSecContextIdentifierDocLink,
		true,
		map[string]string{
			FarEdge:  Optional,
			Telco:    Optional,
			NonTelco: Optional,
			Extended: Mandatory,
		},
		TagExtended)

	TestSecConPrivilegeEscalation = AddCatalogEntry(
		"security-context-privilege-escalation",
		common.AccessControlTestKey,
		`Checks if privileged escalation is enabled (AllowPrivilegeEscalation=true).`,
		SecConPrivilegeRemediation,
		NoExceptions,
		TestSecConPrivilegeEscalationDocLink,
		true,
		map[string]string{
			FarEdge:  Mandatory,
			Telco:    Mandatory,
			NonTelco: Mandatory,
			Extended: Mandatory,
		},
		TagCommon)

	TestContainerHostPort = AddCatalogEntry(
		"container-host-port",
		common.AccessControlTestKey,
		`Verifies if containers define a hostPort.`,
		ContainerHostPortRemediation,
		"Exception for host resource access tests will only be considered in rare cases where it is absolutely needed",
		TestContainerHostPortDocLink,
		true,
		map[string]string{
			FarEdge:  Mandatory,
			Telco:    Mandatory,
			NonTelco: Mandatory,
			Extended: Mandatory,
		},
		TagCommon)

	TestPodHostNetwork = AddCatalogEntry(
		"pod-host-network",
		common.AccessControlTestKey,
		`Verifies that the spec.HostNetwork parameter is not set (not present)`,
		PodHostNetworkRemediation,
		`Exception for host resource access tests will only be considered in rare cases where it is absolutely needed`,
		TestPodHostNetworkDocLink,
		true,
		map[string]string{
			FarEdge:  Mandatory,
			Telco:    Mandatory,
			NonTelco: Mandatory,
			Extended: Mandatory,
		},
		TagCommon)

	TestPodHostPath = AddCatalogEntry(
		"pod-host-path",
		common.AccessControlTestKey,
		`Verifies that the spec.HostPath parameter is not set (not present)`,
		PodHostPathRemediation,
		`Exception for host resource access tests will only be considered in rare cases where it is absolutely needed`,
		TestPodHostPathDocLink,
		true,
		map[string]string{
			FarEdge:  Mandatory,
			Telco:    Mandatory,
			NonTelco: Mandatory,
			Extended: Mandatory,
		},
		TagCommon)

	TestPodHostIPC = AddCatalogEntry(
		"pod-host-ipc",
		common.AccessControlTestKey,
		`Verifies that the spec.HostIpc parameter is set to false`,
		PodHostIPCRemediation,
		`Exception for host resource access tests will only be considered in rare cases where it is absolutely needed`,
		TestPodHostIPCDocLink,
		true,
		map[string]string{
			FarEdge:  Mandatory,
			Telco:    Mandatory,
			NonTelco: Mandatory,
			Extended: Mandatory,
		},
		TagCommon)

	TestPodHostPID = AddCatalogEntry(
		"pod-host-pid",
		common.AccessControlTestKey,
		`Verifies that the spec.HostPid parameter is set to false`,
		PodHostPIDRemediation,
		`Exception for host resource access tests will only be considered in rare cases where it is absolutely needed`,
		TestPodHostPIDDocLink,
		true,
		map[string]string{
			FarEdge:  Mandatory,
			Telco:    Mandatory,
			NonTelco: Mandatory,
			Extended: Mandatory,
		},
		TagCommon)

	TestHugepagesNotManuallyManipulated = AddCatalogEntry(
		"hugepages-config",
		common.PlatformAlterationTestKey,
		`Checks to see that HugePage settings have been configured through MachineConfig, and not manually on the underlying Node. This test case applies only to Nodes that are configured with the "worker" MachineConfigSet. First, the "worker" MachineConfig is polled, and the Hugepage settings are extracted. Next, the underlying Nodes are polled for configured HugePages through inspection of /proc/meminfo. The results are compared, and the test passes only if they are the same.`, //nolint:lll
		HugepagesNotManuallyManipulatedRemediation,
		NoExceptions,
		TestHugepagesNotManuallyManipulatedDocLink,
		true,
		map[string]string{
			FarEdge:  Mandatory,
			Telco:    Mandatory,
			NonTelco: Mandatory,
			Extended: Mandatory,
		},
		TagCommon)

	TestICMPv4ConnectivityIdentifier = AddCatalogEntry(
		"icmpv4-connectivity",
		common.NetworkingTestKey,
		`Checks that each workload Container is able to communicate via ICMPv4 on the Default OpenShift network. This test case requires the Deployment of the debug daemonset and at least 2 pods connected to each network under test(one source and one destination). If no network with more than 2 pods exists this test will be skipped.`,                                                  //nolint:lll
		`Ensure that the workload is able to communicate via the Default OpenShift network. In some rare cases, workloads may require routing table changes in order to communicate over the Default network. To exclude a particular pod from ICMPv4 connectivity tests, add the test-network-function.com/skip_connectivity_tests label to it. The label value is trivial, only its presence.`, //nolint:lll
		`No exceptions - must be able to communicate on default network using IPv4`,
		TestICMPv4ConnectivityIdentifierDocLink,
		true,
		map[string]string{
			FarEdge:  Mandatory,
			Telco:    Mandatory,
			NonTelco: Mandatory,
			Extended: Mandatory,
		},
		TagCommon)

	TestICMPv6ConnectivityIdentifier = AddCatalogEntry(
		"icmpv6-connectivity",
		common.NetworkingTestKey,
		`Checks that each workload Container is able to communicate via ICMPv6 on the Default OpenShift network. This test case requires the Deployment of the debug daemonset and at least 2 pods connected to each network under test(one source and one destination). If no network with more than 2 pods exists this test will be skipped.`, //nolint:lll
		ICMPv6ConnectivityRemediation,
		NoDocumentedProcess,
		TestICMPv6ConnectivityIdentifierDocLink,
		false,
		map[string]string{
			FarEdge:  Mandatory,
			Telco:    Mandatory,
			NonTelco: Optional,
			Extended: Mandatory,
		},
		TagCommon)

	TestICMPv4ConnectivityMultusIdentifier = AddCatalogEntry(
		"icmpv4-connectivity-multus",
		common.NetworkingTestKey,
		`Checks that each workload Container is able to communicate via ICMPv4 on the Multus network(s). This test case requires the Deployment of the debug daemonset and at least 2 pods connected to each network under test(one source and one destination). If no network with more than 2 pods exists this test will be skipped.`, //nolint:lll
		ICMPv4ConnectivityMultusRemediation,
		NoDocumentedProcess,
		TestICMPv4ConnectivityMultusIdentifierDocLink,
		true,
		map[string]string{
			FarEdge:  Mandatory,
			Telco:    Mandatory,
			NonTelco: Optional,
			Extended: Mandatory,
		},
		TagTelco)

	TestICMPv6ConnectivityMultusIdentifier = AddCatalogEntry(
		"icmpv6-connectivity-multus",
		common.NetworkingTestKey,
		`Checks that each workload Container is able to communicate via ICMPv6 on the Multus network(s). This test case requires the Deployment of the debug daemonset and at least 2 pods connected to each network under test(one source and one destination). If no network with more than 2 pods exists this test will be skipped.`, //nolint:lll
		ICMPv6ConnectivityMultusRemediation+` Not applicable if IPv6/MULTUS is not supported.`,
		NoDocumentedProcess,
		TestICMPv6ConnectivityMultusIdentifierDocLink,
		false,
		map[string]string{
			FarEdge:  Mandatory,
			Telco:    Mandatory,
			NonTelco: Optional,
			Extended: Mandatory,
		},
		TagTelco)

	TestServiceDualStackIdentifier = AddCatalogEntry(
		"dual-stack-service",
		common.NetworkingTestKey,
		`Checks that all services in namespaces under test are either ipv6 single stack or dual stack. This test case requires the deployment of the debug daemonset.`,
		TestServiceDualStackRemediation,
		NoExceptionProcessForExtendedTests,
		TestServiceDualStackIdentifierDocLink,
		false,
		map[string]string{
			FarEdge:  Optional,
			Telco:    Optional,
			NonTelco: Optional,
			Extended: Mandatory,
		},
		TagExtended)

	TestNamespaceBestPracticesIdentifier = AddCatalogEntry(
		"namespace",
		common.AccessControlTestKey,
		`Tests that all workload resources (PUTs and CRs) belong to valid namespaces. A valid namespace meets
the following conditions: (1) It was declared in the yaml config file under the targetNameSpaces
tag. (2) It does not have any of the following prefixes: default, openshift-, istio- and aspenmesh-`,
		NamespaceBestPracticesRemediation,
		NoExceptions,
		TestNamespaceBestPracticesIdentifierDocLink,
		true,
		map[string]string{
			FarEdge:  Optional,
			Telco:    Optional,
			NonTelco: Optional,
			Extended: Optional,
		},
		TagCommon)

	TestNonTaintedNodeKernelsIdentifier = AddCatalogEntry(
		"tainted-node-kernel",
		common.PlatformAlterationTestKey,
		`Ensures that the Node(s) hosting workloads do not utilize tainted kernels. This test case is especially
important to support Highly Available workloads, since when a workload is re-instantiated on a backup Node,
that Node's kernel may not have the same hacks.'`,
		NonTaintedNodeKernelsRemediation,
		`If taint is necessary, document details of the taint and why it's needed by workload or environment.`,
		TestNonTaintedNodeKernelsIdentifierDocLink,
		true,
		map[string]string{
			FarEdge:  Mandatory,
			Telco:    Mandatory,
			NonTelco: Mandatory,
			Extended: Mandatory,
		},
		TagCommon)

	TestOperatorInstallStatusSucceededIdentifier = AddCatalogEntry(
		"install-status-succeeded",
		common.OperatorTestKey,
		`Ensures that the target workload operators report "Succeeded" as their installation status.`,
		OperatorInstallStatusSucceededRemediation,
		NoExceptions,
		TestOperatorInstallStatusSucceededIdentifierDocLink,
		true,
		map[string]string{
			FarEdge:  Mandatory,
			Telco:    Mandatory,
			NonTelco: Mandatory,
			Extended: Mandatory,
		},
		TagCommon)

	TestOperatorNoSCCAccess = AddCatalogEntry(
		"install-status-no-privileges",
		common.OperatorTestKey,
		`Checks whether the operator needs access to Security Context Constraints. Test passes if clusterPermissions is not present in the CSV manifest or is present with no RBAC rules related to SCCs.`,
		OperatorNoPrivilegesRemediation,
		NoExceptions,
		TestOperatorNoPrivilegesDocLink,
		true,
		map[string]string{
			FarEdge:  Mandatory,
			Telco:    Mandatory,
			NonTelco: Mandatory,
			Extended: Mandatory,
		},
		TagCommon)

	TestOperatorIsCertifiedIdentifier = AddCatalogEntry(
		"operator-is-certified",
		common.AffiliatedCertTestKey,
		`Tests whether the workload Operators listed in the configuration file have passed the Red Hat Operator Certification Program (OCP).`,
		OperatorIsCertifiedRemediation,
		AffiliatedCert,
		TestOperatorIsCertifiedIdentifierDocLink,
		true,
		map[string]string{
			FarEdge:  Mandatory,
			Telco:    Mandatory,
			NonTelco: Mandatory,
			Extended: Mandatory,
		},
		TagCommon)

	TestHelmIsCertifiedIdentifier = AddCatalogEntry(
		"helmchart-is-certified",
		common.AffiliatedCertTestKey,
		`Tests whether helm charts listed in the cluster passed the Red Hat Helm Certification Program.`,
		HelmIsCertifiedRemediation,
		AffiliatedCert,
		TestHelmIsCertifiedIdentifierDocLink,
		true,
		map[string]string{
			FarEdge:  Mandatory,
			Telco:    Mandatory,
			NonTelco: Mandatory,
			Extended: Mandatory,
		},
		TagCommon)

	TestOperatorIsInstalledViaOLMIdentifier = AddCatalogEntry(
		"install-source",
		common.OperatorTestKey,
		`Tests whether a workload Operator is installed via OLM.`,
		OperatorIsInstalledViaOLMRemediation,
		NoExceptions,
		TestOperatorIsInstalledViaOLMIdentifierDocLink,
		true,
		map[string]string{
			FarEdge:  Mandatory,
			Telco:    Mandatory,
			NonTelco: Mandatory,
			Extended: Mandatory,
		},
		TagCommon)

	TestOperatorHasSemanticVersioningIdentifier = AddCatalogEntry(
		"semantic-versioning",
		common.OperatorTestKey,
		`Tests whether an application Operator has a valid semantic versioning.`,
		OperatorHasSemanticVersioningRemediation,
		NoExceptions,
		TestOperatorHasSemanticVersioningIdentifierDocLink,
		true,
		map[string]string{
			FarEdge:  Mandatory,
			Telco:    Mandatory,
			NonTelco: Mandatory,
			Extended: Mandatory,
		},
		TagCommon)

	TestOperatorRunAsUserID = AddCatalogEntry(
		"run-as-user-id",
		common.OperatorTestKey,
		`Tests that checks the user id of the pods ensure it is not 0.`,
		OperatorRunAsUserID,
		NoExceptions,
		TestOperatorRunAsUserIDDocLink,
		true,
		map[string]string{
			FarEdge:  Mandatory,
			Telco:    Mandatory,
			NonTelco: Mandatory,
			Extended: Mandatory,
		},
		TagCommon)

	TestOperatorRunAsNonRoot = AddCatalogEntry(
		"run-as-non-root",
		common.OperatorTestKey,
		`Tests that checks the pods ensure they are run as non root.`,
		OperatorRunAsNonRoot,
		NoExceptions,
		TestOperatorRunAsNonRootDocLink,
		true,
		map[string]string{
			FarEdge:  Mandatory,
			Telco:    Mandatory,
			NonTelco: Mandatory,
			Extended: Mandatory,
		},
		TagCommon)

	TestOperatorAutomountTokens = AddCatalogEntry(
		"automount-tokens",
		common.OperatorTestKey,
		`Tests that check that the pods disable the automount service account token."`,
		OperatorAutomountTokens,
		NoExceptions,
		TestOperatorAutomountTokensDocLink,
		true,
		map[string]string{
			FarEdge:  Mandatory,
			Telco:    Mandatory,
			NonTelco: Mandatory,
			Extended: Mandatory,
		},
		TagCommon)

	TestOperatorReadOnlyFilesystem = AddCatalogEntry(
		"read-only-file-system",
		common.OperatorTestKey,
		`Tests that check that the pods have the read-only root filesystem setting enabled.`,
		OperatorReadOnlyFilesystem,
		NoExceptions,
		TestOperatorReadOnlyFilesystemDocLink,
		true,
		map[string]string{
			FarEdge:  Optional,
			Telco:    Optional,
			NonTelco: Optional,
			Extended: Optional,
		},
		TagCommon)

	TestOperatorCrdVersioningIdentifier = AddCatalogEntry(
		"crd-versioning",
		common.OperatorTestKey,
		`Tests whether the Operator CRD has a valid versioning.`,
		OperatorCrdVersioningRemediation,
		NoExceptions,
		TestOperatorCrdVersioningIdentifierDocLink,
		true,
		map[string]string{
			FarEdge:  Mandatory,
			Telco:    Mandatory,
			NonTelco: Mandatory,
			Extended: Mandatory,
		},
		TagCommon)

	TestOperatorCrdSchemaIdentifier = AddCatalogEntry(
		"crd-openapi-schema",
		common.OperatorTestKey,
		`Tests whether an application Operator CRD is defined with OpenAPI spec.`,
		OperatorCrdSchemaIdentifierRemediation,
		NoExceptions,
		TestOperatorCrdSchemaIdentifierDocLink,
		true,
		map[string]string{
			FarEdge:  Mandatory,
			Telco:    Mandatory,
			NonTelco: Mandatory,
			Extended: Mandatory,
		},
		TagCommon)

	TestOperatorSingleCrdOwnerIdentifier = AddCatalogEntry(
		"single-crd-owner",
		common.OperatorTestKey,
		`Tests whether a CRD is owned by a single Operator.`,
		OperatorSingleCrdOwnerRemediation,
		NoExceptions,
		TestOperatorSingleCrdOwnerIdentifierDocLink,
		false,
		map[string]string{
			FarEdge:  Mandatory,
			Telco:    Mandatory,
			NonTelco: Mandatory,
			Extended: Mandatory,
		},
		TagCommon)

	TestPodNodeSelectorAndAffinityBestPractices = AddCatalogEntry(
		"pod-scheduling",
		common.LifecycleTestKey,
		`Ensures that workload Pods do not specify nodeSelector or nodeAffinity. In most cases, Pods should allow for instantiation on any underlying Node. Workloads shall not use node selectors nor taints/tolerations to assign pod location.`,
		PodNodeSelectorAndAffinityBestPracticesRemediation,
		`Exception will only be considered if application requires specialized hardware. Must specify which container requires special hardware and why.`,
		TestPodNodeSelectorAndAffinityBestPracticesDocLink,
		true,
		map[string]string{
			FarEdge:  Optional,
			Telco:    Optional,
			NonTelco: Mandatory,
			Extended: Optional,
		},
		TagTelco)

	TestPodHighAvailabilityBestPractices = AddCatalogEntry(
		"pod-high-availability",
		common.LifecycleTestKey,
		`Ensures that workloads Pods specify podAntiAffinity rules and replica value is set to more than 1.`,
		PodHighAvailabilityBestPracticesRemediation,
		NoDocumentedProcess+NotApplicableSNO,
		TestPodHighAvailabilityBestPracticesDocLink,
		true,
		map[string]string{
			FarEdge:  Optional,
			Telco:    Mandatory,
			NonTelco: Mandatory,
			Extended: Mandatory,
		},
		TagCommon)

	TestPodClusterRoleBindingsBestPracticesIdentifier = AddCatalogEntry(
		"cluster-role-bindings",
		common.AccessControlTestKey,
		`Tests that a Pod does not specify ClusterRoleBindings.`,
		PodClusterRoleBindingsBestPracticesRemediation,
		"Exception possible only for workloads that's cluster wide in nature and absolutely needs cluster level roles & role bindings",
		TestPodClusterRoleBindingsBestPracticesIdentifierDocLink,
		true,
		map[string]string{
			FarEdge:  Mandatory,
			Telco:    Mandatory,
			NonTelco: Optional,
			Extended: Mandatory,
		},
		TagTelco)

	TestPodDeploymentBestPracticesIdentifier = AddCatalogEntry(
		"pod-owner-type",
		common.LifecycleTestKey,
		`Tests that the workload Pods are deployed as part of a ReplicaSet(s)/StatefulSet(s).`,
		PodDeploymentBestPracticesRemediation,
		NoDocumentedProcess+` Pods should not be deployed as DaemonSet or naked pods.`,
		TestPodDeploymentBestPracticesIdentifierDocLink,
		true,
		map[string]string{
			FarEdge:  Mandatory,
			Telco:    Mandatory,
			NonTelco: Optional,
			Extended: Mandatory,
		},
		TagTelco)

	TestDeploymentScalingIdentifier = AddCatalogEntry(
		"deployment-scaling",
		common.LifecycleTestKey,
		`Tests that workload deployments support scale in/out operations. First, the test starts getting the current replicaCount (N) of the deployment/s with the Pod Under Test. Then, it executes the scale-in oc command for (N-1) replicas. Lastly, it executes the scale-out oc command, restoring the original replicaCount of the deployment/s. In case of deployments that are managed by HPA the test is changing the min and max value to deployment Replica - 1 during scale-in and the original replicaCount again for both min/max during the scale-out stage. Lastly its restoring the original min/max replica of the deployment/s`, //nolint:lll
		DeploymentScalingRemediation,
		NoDocumentedProcess+NotApplicableSNO,
		TestDeploymentScalingIdentifierDocLink,
		true,
		map[string]string{
			FarEdge:  Optional,
			Telco:    Mandatory,
			NonTelco: Mandatory,
			Extended: Mandatory,
		},
		TagCommon)

	TestStateFulSetScalingIdentifier = AddCatalogEntry(
		"statefulset-scaling",
		common.LifecycleTestKey,
		`Tests that workload statefulsets support scale in/out operations. First, the test starts getting the current replicaCount (N) of the statefulset/s with the Pod Under Test. Then, it executes the scale-in oc command for (N-1) replicas. Lastly, it executes the scale-out oc command, restoring the original replicaCount of the statefulset/s. In case of statefulsets that are managed by HPA the test is changing the min and max value to statefulset Replica - 1 during scale-in and the original replicaCount again for both min/max during the scale-out stage. Lastly its restoring the original min/max replica of the statefulset/s`, //nolint:lll
		StatefulSetScalingRemediation,
		NoDocumentedProcess+NotApplicableSNO,
		TestStateFulSetScalingIdentifierDocLink,
		true,
		map[string]string{
			FarEdge:  Optional,
			Telco:    Mandatory,
			NonTelco: Mandatory,
			Extended: Mandatory,
		},
		TagCommon)

	TestImagePullPolicyIdentifier = AddCatalogEntry(
		"image-pull-policy",
		common.LifecycleTestKey,
		`Ensure that the containers under test are using IfNotPresent as Image Pull Policy. If there is a situation where the container dies and needs to be restarted, the image pull policy becomes important. PullIfNotPresent is recommended so that a loss of image registry access does not prevent the pod from restarting.`, //nolint:lll
		ImagePullPolicyRemediation,
		NoDocumentedProcess,
		TestImagePullPolicyIdentifierDocLink,
		true,
		map[string]string{
			FarEdge:  Mandatory,
			Telco:    Mandatory,
			NonTelco: Optional,
			Extended: Mandatory,
		},
		TagTelco)

	TestPodRecreationIdentifier = AddCatalogEntry(
		"pod-recreation",
		common.LifecycleTestKey,
		`Tests that a workload is configured to support High Availability. First, this test cordons and drains a Node that hosts the workload Pod. Next, the test ensures that OpenShift can re-instantiate the Pod on another Node, and that the actual replica count matches the desired replica count.`, //nolint:lll
		PodRecreationRemediation,
		`No exceptions - workloads should be able to be restarted/recreated.`,
		TestPodRecreationIdentifierDocLink,
		true,
		map[string]string{
			FarEdge:  Mandatory,
			Telco:    Mandatory,
			NonTelco: Mandatory,
			Extended: Mandatory,
		},
		TagCommon)

	TestPodRoleBindingsBestPracticesIdentifier = AddCatalogEntry(
		"pod-role-bindings",
		common.AccessControlTestKey,
		`Ensures that a workload does not utilize RoleBinding(s) in a non-workload Namespace.`,
		PodRoleBindingsBestPracticesRemediation,
		NoExceptions,
		TestPodRoleBindingsBestPracticesIdentifierDocLink,
		true,
		map[string]string{
			FarEdge:  Mandatory,
			Telco:    Mandatory,
			NonTelco: Mandatory,
			Extended: Mandatory,
		},
		TagCommon)

	TestPodServiceAccountBestPracticesIdentifier = AddCatalogEntry(
		"pod-service-account",
		common.AccessControlTestKey,
		`Tests that each workload Pod utilizes a valid Service Account. Default or empty service account is not valid.`,
		PodServiceAccountBestPracticesRemediation,
		NoExceptions,
		TestPodServiceAccountBestPracticesIdentifierDocLink,
		true,
		map[string]string{
			FarEdge:  Mandatory,
			Telco:    Mandatory,
			NonTelco: Mandatory,
			Extended: Mandatory,
		},
		TagCommon)

	TestPodAutomountServiceAccountIdentifier = AddCatalogEntry(
		"pod-automount-service-account-token",
		common.AccessControlTestKey,
		`Check that all pods under test have automountServiceAccountToken set to false. Only pods that require access to the kubernetes API server should have automountServiceAccountToken set to true`,
		AutomountServiceTokenRemediation,
		`Exception will be considered if container needs to access APIs which OCP does not offer natively. Must document which container requires which API(s) and detail why existing OCP APIs cannot be used.`,
		TestPodAutomountServiceAccountIdentifierDocLink,
		true,
		map[string]string{
			FarEdge:  Mandatory,
			Telco:    Mandatory,
			NonTelco: Optional,
			Extended: Mandatory,
		},
		TagTelco)

	TestServicesDoNotUseNodeportsIdentifier = AddCatalogEntry(
		"service-type",
		common.AccessControlTestKey,
		`Tests that each workload Service does not utilize NodePort(s).`,
		ServicesDoNotUseNodeportsRemediation,
		`Exception for host resource access tests will only be considered in rare cases where it is absolutely needed`,
		TestServicesDoNotUseNodeportsIdentifierDocLink,
		true,
		map[string]string{
			FarEdge:  Mandatory,
			Telco:    Mandatory,
			NonTelco: Mandatory,
			Extended: Mandatory,
		},
		TagCommon)

	TestUnalteredBaseImageIdentifier = AddCatalogEntry(
		"base-image",
		common.PlatformAlterationTestKey,
		`Ensures that the Container Base Image is not altered post-startup. This test is a heuristic, and ensures that there are no changes to the following directories: 1) /var/lib/rpm 2) /var/lib/dpkg 3) /bin 4) /sbin 5) /lib 6) /lib64 7) /usr/bin 8) /usr/sbin 9) /usr/lib 10) /usr/lib64`, //nolint:lll
		UnalteredBaseImageRemediation,
		NoExceptions,
		TestUnalteredBaseImageIdentifierDocLink,
		true,
		map[string]string{
			FarEdge:  Mandatory,
			Telco:    Mandatory,
			NonTelco: Mandatory,
			Extended: Mandatory,
		},
		TagCommon)

	TestUnalteredStartupBootParamsIdentifier = AddCatalogEntry(
		"boot-params",
		common.PlatformAlterationTestKey,
		`Tests that boot parameters are set through the MachineConfigOperator, and not set manually on the Node.`,
		UnalteredStartupBootParamsRemediation,
		NoExceptions,
		TestUnalteredStartupBootParamsIdentifierDocLink,
		true,
		map[string]string{
			FarEdge:  Mandatory,
			Telco:    Mandatory,
			NonTelco: Mandatory,
			Extended: Mandatory,
		},
		TagCommon)

	TestLoggingIdentifier = AddCatalogEntry(
		"container-logging",
		common.ObservabilityTestKey,
		`Check that all containers under test use standard input output and standard error when logging. A container must provide APIs for the platform to observe the container health and act accordingly. These APIs include health checks (liveness and readiness), logging to stderr and stdout for log aggregation (by tools such as Logstash or Filebeat), and integrate with tracing and metrics-gathering libraries (such as Prometheus or Metricbeat).`, //nolint:lll
		LoggingRemediation,
		NoDocumentedProcess,
		TestLoggingIdentifierDocLink,
		true,
		map[string]string{
			FarEdge:  Mandatory,
			Telco:    Mandatory,
			NonTelco: Optional,
			Extended: Mandatory,
		},
		TagTelco)

	TestTerminationMessagePolicyIdentifier = AddCatalogEntry(
		"termination-policy",
		common.ObservabilityTestKey,
		`Check that all containers are using terminationMessagePolicy: FallbackToLogsOnError. There are different ways a pod can stop on an OpenShift cluster. One way is that the pod can remain alive but non-functional. Another way is that the pod can crash and become non-functional. In the first case, if the administrator has implemented liveness and readiness checks, OpenShift can stop the pod and either restart it on the same node or a different node in the cluster. For the second case, when the application in the pod stops, it should exit with a code and write suitable log entries to help the administrator diagnose what the issue was that caused the problem.`, //nolint:lll
		TerminationMessagePolicyRemediation,
		NoDocumentedProcess,
		TestTerminationMessagePolicyIdentifierDocLink,
		true,
		map[string]string{
			FarEdge:  Mandatory,
			Telco:    Mandatory,
			NonTelco: Optional,
			Extended: Mandatory,
		},
		TagTelco)

	TestCrdsStatusSubresourceIdentifier = AddCatalogEntry(
		"crd-status",
		common.ObservabilityTestKey,
		`Checks that all CRDs have a status sub-resource specification (Spec.versions[].Schema.OpenAPIV3Schema.Properties[“status”]).`,
		CrdsStatusSubresourceRemediation,
		NoExceptions,
		TestCrdsStatusSubresourceIdentifierDocLink,
		true,
		map[string]string{
			FarEdge:  Mandatory,
			Telco:    Mandatory,
			NonTelco: Mandatory,
			Extended: Mandatory,
		},
		TagCommon)

	TestSysctlConfigsIdentifier = AddCatalogEntry(
		"sysctl-config",
		common.PlatformAlterationTestKey,
		`Tests that no one has changed the node's sysctl configs after the node was created, the tests works by checking if the sysctl configs are consistent with the MachineConfig CR which defines how the node should be configured`,
		SysctlConfigsRemediation,
		NoExceptions,
		TestSysctlConfigsIdentifierDocLink,
		true,
		map[string]string{
			FarEdge:  Mandatory,
			Telco:    Mandatory,
			NonTelco: Mandatory,
			Extended: Mandatory,
		},
		TagCommon)

	TestServiceMeshIdentifier = AddCatalogEntry(
		"service-mesh-usage",
		common.PlatformAlterationTestKey,
		`Checks if the istio namespace ("istio-system") is present. If it is present, checks that the istio sidecar is present in all pods under test.`,
		ServiceMeshRemediation,
		NoExceptionProcessForExtendedTests,
		TestServiceMeshIdentifierDocLink,
		true,
		map[string]string{
			FarEdge:  Optional,
			Telco:    Optional,
			NonTelco: Optional,
			Extended: Optional,
		},
		TagExtended)

	TestOCPLifecycleIdentifier = AddCatalogEntry(
		"ocp-lifecycle",
		common.PlatformAlterationTestKey,
		`Tests that the running OCP version is not end of life.`,
		OCPLifecycleRemediation,
		NoExceptions,
		TestOCPLifecycleIdentifierDocLink,
		true,
		map[string]string{
			FarEdge:  Mandatory,
			Telco:    Mandatory,
			NonTelco: Mandatory,
			Extended: Mandatory,
		},
		TagCommon)

	TestNodeOperatingSystemIdentifier = AddCatalogEntry(
		"ocp-node-os-lifecycle",
		common.PlatformAlterationTestKey,
		`Tests that the nodes running in the cluster have operating systems that are compatible with the deployed version of OpenShift.`,
		NodeOperatingSystemRemediation,
		NoExceptions,
		TestNodeOperatingSystemIdentifierDocLink,
		true,
		map[string]string{
			FarEdge:  Mandatory,
			Telco:    Mandatory,
			NonTelco: Mandatory,
			Extended: Mandatory,
		},
		TagCommon)

	TestIsRedHatReleaseIdentifier = AddCatalogEntry(
		"isredhat-release",
		common.PlatformAlterationTestKey,
		`verifies if the container base image is redhat.`,
		IsRedHatReleaseRemediation,
		NoExceptions,
		TestIsRedHatReleaseIdentifierDocLink,
		true,
		map[string]string{
			FarEdge:  Mandatory,
			Telco:    Mandatory,
			NonTelco: Mandatory,
			Extended: Mandatory,
		},
		TagCommon)

	TestIsSELinuxEnforcingIdentifier = AddCatalogEntry(
		"is-selinux-enforcing",
		common.PlatformAlterationTestKey,
		`verifies that all openshift platform/cluster nodes have selinux in "Enforcing" mode.`,
		IsSELinuxEnforcingRemediation,
		NoExceptions,
		TestIsSELinuxEnforcingIdentifierDocLink,
		true,
		map[string]string{
			FarEdge:  Mandatory,
			Telco:    Mandatory,
			NonTelco: Mandatory,
			Extended: Mandatory,
		},
		TagCommon)

	TestUndeclaredContainerPortsUsage = AddCatalogEntry(
		"undeclared-container-ports-usage",
		common.NetworkingTestKey,
		`Check that containers do not listen on ports that weren't declared in their specification. Platforms may be configured to block undeclared ports.`,
		UndeclaredContainerPortsRemediation,
		NoExceptionProcessForExtendedTests,
		TestUndeclaredContainerPortsUsageDocLink,
		true,
		map[string]string{
			FarEdge:  Optional,
			Telco:    Optional,
			NonTelco: Optional,
			Extended: Mandatory,
		},
		TagExtended)

	TestOCPReservedPortsUsage = AddCatalogEntry(
		"ocp-reserved-ports-usage",
		common.NetworkingTestKey,
		`Check that containers do not listen on ports that are reserved by OpenShift`,
		OCPReservedPortsUsageRemediation,
		NoExceptions,
		TestOCPReservedPortsUsageDocLink,
		true,
		map[string]string{
			FarEdge:  Mandatory,
			Telco:    Mandatory,
			NonTelco: Mandatory,
			Extended: Mandatory,
		},
		TagCommon)

	TestLivenessProbeIdentifier = AddCatalogEntry(
		"liveness-probe",
		common.LifecycleTestKey,
		`Check that all containers under test have liveness probe defined. The most basic requirement for the lifecycle management of Pods in OpenShift are the ability to start and stop correctly. When starting up, health probes like liveness and readiness checks can be put into place to ensure the application is functioning properly.`, //nolint:lll
		LivenessProbeRemediation+` workloads shall self-recover from common failures like pod failure, host failure, and network failure. Kubernetes native mechanisms such as health-checks (Liveness, Readiness and Startup Probes) shall be employed at a minimum.`,                                                                            //nolint:lll
		NoDocumentedProcess,
		TestLivenessProbeIdentifierDocLink,
		true,
		map[string]string{
			FarEdge:  Mandatory,
			Telco:    Mandatory,
			NonTelco: Optional,
			Extended: Mandatory,
		},
		TagTelco)

	TestReadinessProbeIdentifier = AddCatalogEntry(
		"readiness-probe",
		common.LifecycleTestKey,
		`Check that all containers under test have readiness probe defined. There are different ways a pod can stop on on OpenShift cluster. One way is that the pod can remain alive but non-functional. Another way is that the pod can crash and become non-functional. In the first case, if the administrator has implemented liveness and readiness checks, OpenShift can stop the pod and either restart it on the same node or a different node in the cluster. For the second case, when the application in the pod stops, it should exit with a code and write suitable log entries to help the administrator diagnose what the issue was that caused the problem.`, //nolint:lll
		ReadinessProbeRemediation,
		NoDocumentedProcess,
		TestReadinessProbeIdentifierDocLink,
		true,
		map[string]string{
			FarEdge:  Mandatory,
			Telco:    Mandatory,
			NonTelco: Optional,
			Extended: Mandatory,
		},
		TagTelco)

	TestStartupProbeIdentifier = AddCatalogEntry(
		"startup-probe",
		common.LifecycleTestKey,
		`Check that all containers under test have startup probe defined. Workloads shall self-recover from common failures like pod failure, host failure, and network failure. Kubernetes native mechanisms such as health-checks (Liveness, Readiness and Startup Probes) shall be employed at a minimum.`, //nolint:lll
		StartupProbeRemediation,
		NoDocumentedProcess,
		TestStartupProbeIdentifierDocLink,
		true,
		map[string]string{
			FarEdge:  Mandatory,
			Telco:    Mandatory,
			NonTelco: Optional,
			Extended: Mandatory,
		},
		TagTelco)

	TestOneProcessPerContainerIdentifier = AddCatalogEntry(
		"one-process-per-container",
		common.AccessControlTestKey,
		`Check that all containers under test have only one process running`,
		OneProcessPerContainerRemediation,
		NoExceptionProcessForExtendedTests+NotApplicableSNO,
		TestOneProcessPerContainerIdentifierDocLink,
		true,
		map[string]string{
			FarEdge:  Optional,
			Telco:    Optional,
			NonTelco: Optional,
			Extended: Optional,
		},
		TagCommon)

	TestSYSNiceRealtimeCapabilityIdentifier = AddCatalogEntry(
		"sys-nice-realtime-capability",
		common.AccessControlTestKey,
		`Check that pods running on nodes with realtime kernel enabled have the SYS_NICE capability enabled in their spec. In the case that a workolad is running on a node using the real-time kernel, SYS_NICE will be used to allow DPDK application to switch to SCHED_FIFO.`, //nolint:lll
		SYSNiceRealtimeCapabilityRemediation,
		NoDocumentedProcess,
		TestSYSNiceRealtimeCapabilityIdentifierDocLink,
		false,
		map[string]string{
			FarEdge:  Mandatory,
			Telco:    Mandatory,
			NonTelco: Optional,
			Extended: Mandatory,
		},
		TagTelco)

	TestSysPtraceCapabilityIdentifier = AddCatalogEntry(
		"sys-ptrace-capability",
		common.AccessControlTestKey,
		`Check that if process namespace sharing is enabled for a Pod then the SYS_PTRACE capability is allowed. This capability is required when using Process Namespace Sharing. This is used when processes from one Container need to be exposed to another Container. For example, to send signals like SIGHUP from a process in a Container to another process in another Container. For more information on these capabilities refer to https://cloud.redhat.com/blog/linux-capabilities-in-openshift and https://kubernetes.io/docs/tasks/configure-pod-container/share-process-namespace/`, //nolint:lll
		SysPtraceCapabilityRemediation,
		NoDocumentedProcess,
		TestSysPtraceCapabilityIdentifierDocLink,
		true,
		map[string]string{
			FarEdge:  Mandatory,
			Telco:    Mandatory,
			NonTelco: Optional,
			Extended: Mandatory,
		},
		TagTelco)

	TestPodRequestsAndLimitsIdentifier = AddCatalogEntry(
		"requests-and-limits",
		common.AccessControlTestKey,
		`Check that containers have resource requests and limits specified in their spec.`,
		RequestsAndLimitsRemediation,
		NoDocumentedProcess,
		TestPodRequestsAndLimitsIdentifierDocLink,
		true,
		map[string]string{
			FarEdge:  Mandatory,
			Telco:    Mandatory,
			NonTelco: Optional,
			Extended: Mandatory,
		},
		TagTelco)

	TestNamespaceResourceQuotaIdentifier = AddCatalogEntry(
		"namespace-resource-quota",
		common.AccessControlTestKey,
		`Checks to see if workload pods are running in namespaces that have resource quotas applied.`,
		NamespaceResourceQuotaRemediation,
		NoExceptionProcessForExtendedTests,
		TestNamespaceResourceQuotaIdentifierDocLink,
		true,
		map[string]string{
			FarEdge:  Optional,
			Telco:    Optional,
			NonTelco: Optional,
			Extended: Mandatory,
		},
		TagExtended)

	TestPodDisruptionBudgetIdentifier = AddCatalogEntry(
		"pod-disruption-budget",
		common.ObservabilityTestKey,
		`Checks to see if pod disruption budgets have allowed values for minAvailable and maxUnavailable`,
		PodDisruptionBudgetRemediation,
		NoExceptions,
		TestPodDisruptionBudgetIdentifierDocLink,
		true,
		map[string]string{
			FarEdge:  Mandatory,
			Telco:    Mandatory,
			NonTelco: Mandatory,
			Extended: Mandatory,
		},
		TagCommon)

	TestPodTolerationBypassIdentifier = AddCatalogEntry(
		"pod-toleration-bypass",
		common.LifecycleTestKey,
		`Check that pods do not have NoExecute, PreferNoSchedule, or NoSchedule tolerations that have been modified from the default.`,
		PodTolerationBypassRemediation,
		NoDocumentedProcess,
		TestPodTolerationBypassIdentifierDocLink,
		true,
		map[string]string{
			FarEdge:  Mandatory,
			Telco:    Mandatory,
			NonTelco: Optional,
			Extended: Mandatory,
		},
		TagTelco)

	TestPersistentVolumeReclaimPolicyIdentifier = AddCatalogEntry(
		"persistent-volume-reclaim-policy",
		common.LifecycleTestKey,
		`Check that the persistent volumes the workloads pods are using have a reclaim policy of delete. Network Functions should clear persistent storage by deleting their PVs when removing their application from a cluster.`,
		PersistentVolumeReclaimPolicyRemediation,
		NoDocumentedProcess,
		TestPersistentVolumeReclaimPolicyIdentifierDocLink,
		true,
		map[string]string{
			FarEdge:  Mandatory,
			Telco:    Mandatory,
			NonTelco: Optional,
			Extended: Mandatory,
		},
		TagTelco)

	TestContainersImageTag = AddCatalogEntry(
		"containers-image-tag",
		common.ManageabilityTestKey,
		`Check that image tag exists on containers.`,
		ContainersImageTagRemediation,
		NoExceptionProcessForExtendedTests,
		TestContainersImageTagDocLink,
		true,
		map[string]string{
			FarEdge:  Optional,
			Telco:    Optional,
			NonTelco: Optional,
			Extended: Optional,
		},
		TagExtended)

	TestNoSSHDaemonsAllowedIdentifier = AddCatalogEntry(
		"ssh-daemons",
		common.AccessControlTestKey,
		`Check that pods do not run SSH daemons.`,
		NoSSHDaemonsAllowedRemediation,
		`No exceptions - special consideration can be given to certain containers which run as utility tool daemon`,
		TestNoSSHDaemonsAllowedIdentifierDocLink,
		true,
		map[string]string{
			FarEdge:  Mandatory,
			Telco:    Mandatory,
			NonTelco: Optional,
			Extended: Mandatory,
		},
		TagTelco)

	TestCPUIsolationIdentifier = AddCatalogEntry(
		"cpu-isolation",
		common.LifecycleTestKey,
		`CPU isolation requires: For each container within the pod, resource requests and limits must be identical. If cpu requests and limits are not identical and in whole units (Guaranteed pods with exclusive cpus), your pods will not be tested for compliance. The runTimeClassName must be specified. Annotations required disabling CPU and IRQ load-balancing.`, //nolint:lll
		CPUIsolationRemediation,
		NoDocumentedProcess,
		TestCPUIsolationIdentifierDocLink,
		true,
		map[string]string{
			FarEdge:  Mandatory,
			Telco:    Mandatory,
			NonTelco: Optional,
			Extended: Mandatory,
		},
		TagTelco)

	TestContainerPortNameFormat = AddCatalogEntry(
		"container-port-name-format",
		common.ManageabilityTestKey,
		"Check that the container's ports name follow the naming conventions. Name field in ContainerPort section must be of form `<protocol>[-<suffix>]`. More naming convention requirements may be released in future",
		ContainerPortNameFormatRemediation,
		NoExceptionProcessForExtendedTests,
		TestContainerPortNameFormatDocLink,
		true,
		map[string]string{
			FarEdge:  Optional,
			Telco:    Optional,
			NonTelco: Optional,
			Extended: Mandatory,
		},
		TagExtended)

	TestCrdScalingIdentifier = AddCatalogEntry(
		"crd-scaling",
		common.LifecycleTestKey,
		`Tests that a workload's CRD support scale in/out operations. First, the test starts getting the current replicaCount (N) of the crd/s with the Pod Under Test. Then, it executes the scale-in oc command for (N-1) replicas. Lastly, it executes the scale-out oc command, restoring the original replicaCount of the crd/s. In case of crd that are managed by HPA the test is changing the min and max value to crd Replica - 1 during scale-in and the original replicaCount again for both min/max during the scale-out stage. Lastly its restoring the original min/max replica of the crd/s`, //nolint:lll
		CrdScalingRemediation,
		NoDocumentedProcess+NotApplicableSNO,
		TestCrdScalingIdentifierDocLink,
		true,
		map[string]string{
			FarEdge:  Optional,
			Telco:    Mandatory,
			NonTelco: Mandatory,
			Extended: Mandatory,
		},
		TagCommon,
	)

	TestCrdRoleIdentifier = AddCatalogEntry(
		"crd-roles",
		common.AccessControlTestKey,
		"If an application creates CRDs it must supply a role to access those CRDs and no other API resources/permission. This test checks that there is at least one role present in each namespaces under test that only refers to CRDs under test.",
		"Roles providing access to CRDs should not refer to any other api or resources. Change the generation of the CRD role accordingly",
		NoExceptionProcessForExtendedTests,
		"https://test-network-function.github.io/k8s-best-practices-guide-guide/#k8s-best-practices-custom-role-to-access-application-crds",
		true,
		map[string]string{
			FarEdge:  Optional,
			Telco:    Optional,
			NonTelco: Optional,
			Extended: Mandatory,
		},
		TagExtended)

	//nolint:gocritic
	// TestPodDeleteIdentifier = AddCatalogEntry(
	// 	"pod-delete",
	// 	common.ChaosTesting,
	// 	"Chaos test suite is under construction.",
	// 	"",
	// 	NoDocumentedProcess,
	// 	"",
	// 	false,
	// 	map[string]string{
	// 		FarEdge:  Optional,
	// 		Telco:    Optional,
	// 		NonTelco: Optional,
	// 		Extended: Optional,
	// 	},
	// 	TagCommon)

	return Catalog
}

var (
	// TestIdToClaimId converts the testcase short ID to the claim identifier.
	TestIDToClaimID = map[string]claim.Identifier{}
)

// GetTestIDAndLabels transform the claim.Identifier into a test Id that can be used to skip
// specific tests.
func GetTestIDAndLabels(identifier claim.Identifier) (testID string, tags []string) {
	tags = strings.Split(identifier.Tags, ",")
	tags = append(tags, identifier.Id, identifier.Suite)
	TestIDToClaimID[identifier.Id] = identifier
	return identifier.Id, tags
}

// Catalog is the JUnit testcase catalog of tests.
var Catalog = map[claim.Identifier]claim.TestCaseDescription{}
var Classification = map[string]map[string]string{}
