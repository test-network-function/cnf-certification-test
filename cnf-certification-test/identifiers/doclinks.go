package identifiers

const (
	// Default Strings
	NoDocLinkExtended = "No Doc Link - Extended"
	NoDocLinkFarEdge  = "No Doc Link - Far Edge"
	NoDocLinkTelco    = "No Doc Link - Telco"
	NoDocLink         = "No Doc Link"

	// Networking Suite
	TestICMPv4ConnectivityIdentifierDocLink         = "https://test-network-function.github.io/cnf-best-practices-guide/#cnf-best-practices-ipv4-&-ipv6"
	TestNetworkPolicyDenyAllIdentifierDocLink       = "https://test-network-function.github.io/cnf-best-practices-guide/#cnf-best-practices-vrfs-aka-routing-instances"
	TestReservedExtendedPartnerPortsDocLink         = NoDocLinkExtended
	TestDpdkCPUPinningExecProbeDocLink              = "https://test-network-function.github.io/cnf-best-practices-guide/#cnf-best-practices-cpu-manager-pinning"
	TestRestartOnRebootLabelOnPodsUsingSRIOVDocLink = NoDocLinkFarEdge
	TestLimitedUseOfExecProbesIdentifierDocLink     = NoDocLinkFarEdge
	TestICMPv6ConnectivityIdentifierDocLink         = "https://test-network-function.github.io/cnf-best-practices-guide/#cnf-best-practices-ipv4-&-ipv6"
	TestICMPv4ConnectivityMultusIdentifierDocLink   = "https://test-network-function.github.io/cnf-best-practices-guide/#cnf-best-practices-high-level-cnf-expectations"
	TestICMPv6ConnectivityMultusIdentifierDocLink   = "https://test-network-function.github.io/cnf-best-practices-guide/#cnf-best-practices-high-level-cnf-expectations"
	TestServiceDualStackIdentifierDocLink           = "https://test-network-function.github.io/cnf-best-practices-guide/#cnf-best-practices-ipv4-&-ipv6"
	TestUndeclaredContainerPortsUsageDocLink        = "https://test-network-function.github.io/cnf-best-practices-guide/#cnf-best-practices-requirements-cnf-reqs"
	TestOCPReservedPortsUsageDocLink                = "https://test-network-function.github.io/cnf-best-practices-guide/#cnf-best-practices-ports-reserved-by-openshift"

	// Access Control Suite
	Test1337UIDIdentifierDocLink                             = NoDocLinkExtended
	TestNetAdminIdentifierDocLink                            = "https://test-network-function.github.io/cnf-best-practices-guide/#cnf-best-practices-net_admin"
	TestSysAdminIdentifierDocLink                            = "https://test-network-function.github.io/cnf-best-practices-guide/#cnf-best-practices-avoid-sys_admin"
	TestIpcLockIdentifierDocLink                             = "https://test-network-function.github.io/cnf-best-practices-guide/#cnf-best-practices-ipc_lock"
	TestNetRawIdentifierDocLink                              = "https://test-network-function.github.io/cnf-best-practices-guide/#cnf-best-practices-user-plane-cnfs"
	TestBpfIdentifierDocLink                                 = NoDocLinkTelco
	TestSecConNonRootUserIdentifierDocLink                   = "https://test-network-function.github.io/cnf-best-practices-guide/#cnf-best-practices-cnf-security"
	TestSecContextIdentifierDocLink                          = "https://test-network-function.github.io/cnf-best-practices-guide/#cnf-best-practices-cnf-security"
	TestSecConPrivilegeEscalationDocLink                     = "https://test-network-function.github.io/cnf-best-practices-guide/#cnf-best-practices-cnf-security"
	TestContainerHostPortDocLink                             = "https://test-network-function.github.io/cnf-best-practices-guide/#cnf-best-practices-avoid-accessing-resource-on-host"
	TestContainerHostNetworkDocLink                          = "https://test-network-function.github.io/cnf-best-practices-guide/#cnf-best-practices-avoid-the-host-network-namespace"
	TestPodHostNetworkDocLink                                = "https://test-network-function.github.io/cnf-best-practices-guide/#cnf-best-practices-avoid-the-host-network-namespace"
	TestPodHostPathDocLink                                   = "https://test-network-function.github.io/cnf-best-practices-guide/#cnf-best-practices-cnf-security"
	TestPodHostIPCDocLink                                    = "https://test-network-function.github.io/cnf-best-practices-guide/#cnf-best-practices-cnf-security"
	TestPodHostPIDDocLink                                    = "https://test-network-function.github.io/cnf-best-practices-guide/#cnf-best-practices-cnf-security"
	TestNamespaceBestPracticesIdentifierDocLink              = "https://test-network-function.github.io/cnf-best-practices-guide/#cnf-best-practices-requirements-cnf-reqs"
	TestPodClusterRoleBindingsBestPracticesIdentifierDocLink = "https://test-network-function.github.io/cnf-best-practices-guide/#cnf-best-practices-security-rbac"
	TestPodRoleBindingsBestPracticesIdentifierDocLink        = "https://test-network-function.github.io/cnf-best-practices-guide/#cnf-best-practices-security-rbac"
	TestPodServiceAccountBestPracticesIdentifierDocLink      = "https://test-network-function.github.io/cnf-best-practices-guide/#cnf-best-practices-scc-permissions-for-an-application"
	TestPodAutomountServiceAccountIdentifierDocLink          = "https://test-network-function.github.io/cnf-best-practices-guide/#cnf-best-practices-automount-services-for-pods"
	TestServicesDoNotUseNodeportsIdentifierDocLink           = "https://test-network-function.github.io/cnf-best-practices-guide/#cnf-best-practices-avoid-the-host-network-namespace"
	TestUnalteredBaseImageIdentifierDocLink                  = "https://test-network-function.github.io/cnf-best-practices-guide/#cnf-best-practices-image-standards"
	TestOneProcessPerContainerIdentifierDocLink              = "https://test-network-function.github.io/cnf-best-practices-guide/#cnf-best-practices-one-process-per-container"
	TestSYSNiceRealtimeCapabilityIdentifierDocLink           = "https://test-network-function.github.io/cnf-best-practices-guide/#cnf-best-practices-sys_nice"
	TestSysPtraceCapabilityIdentifierDocLink                 = "https://test-network-function.github.io/cnf-best-practices-guide/#cnf-best-practices-sys_ptrace"
	TestPodRequestsAndLimitsIdentifierDocLink                = "https://test-network-function.github.io/cnf-best-practices-guide/#cnf-best-practices-requests/limits"
	TestNamespaceResourceQuotaIdentifierDocLink              = "https://test-network-function.github.io/cnf-best-practices-guide/#cnf-best-practices-memory-allocation"
	TestNoSSHDaemonsAllowedIdentifierDocLink                 = "https://test-network-function.github.io/cnf-best-practices-guide/#cnf-best-practices-pod-interaction/configuration"

	// Affiliated Certification Suite
	TestHelmVersionIdentifierDocLink                = "https://test-network-function.github.io/cnf-best-practices-guide/#cnf-best-practices-helm"
	TestContainerIsCertifiedDigestIdentifierDocLink = "https://redhat-connect.gitbook.io/partner-guide-for-red-hat-openshift-and-container/certify-your-application/overview"
	TestContainerIsCertifiedIdentifierDocLink       = "https://redhat-connect.gitbook.io/partner-guide-for-red-hat-openshift-and-container/certify-your-application/overview"
	TestHelmIsCertifiedIdentifierDocLink            = "https://redhat-connect.gitbook.io/partner-guide-for-red-hat-openshift-and-container/certify-your-application/overview"

	// Platform Alteration Suite
	TestPodHugePages2MDocLink                       = "https://test-network-function.github.io/cnf-best-practices-guide/#cnf-best-practices-huge-pages"
	TestPodHugePages1GDocLink                       = NoDocLinkFarEdge
	TestHugepagesNotManuallyManipulatedDocLink      = "https://test-network-function.github.io/cnf-best-practices-guide/#cnf-best-practices-huge-pages"
	TestNonTaintedNodeKernelsIdentifierDocLink      = "https://test-network-function.github.io/cnf-best-practices-guide/#cnf-best-practices-high-level-cnf-expectations"
	TestUnalteredStartupBootParamsIdentifierDocLink = "https://test-network-function.github.io/cnf-best-practices-guide/#cnf-best-practices-host-os"
	TestSysctlConfigsIdentifierDocLink              = "https://test-network-function.github.io/cnf-best-practices-guide/#cnf-best-practices-cnf-security"
	TestServiceMeshIdentifierDocLink                = NoDocLinkExtended
	TestHyperThreadEnableDocLink                    = NoDocLinkExtended

	TestOCPLifecycleIdentifierDocLink        = "https://test-network-function.github.io/cnf-best-practices-guide/#cnf-best-practices-k8s"
	TestNodeOperatingSystemIdentifierDocLink = "https://test-network-function.github.io/cnf-best-practices-guide/#cnf-best-practices-host-os"
	TestIsRedHatReleaseIdentifierDocLink     = "https://test-network-function.github.io/cnf-best-practices-guide/#cnf-best-practices-base-images"
	TestIsSELinuxEnforcingIdentifierDocLink  = "https://test-network-function.github.io/cnf-best-practices-guide/#cnf-best-practices-pod-security"

	// Lifecycle Suite
	TestAffinityRequiredPodsDocLink                    = "https://test-network-function.github.io/cnf-best-practices-guide/#cnf-best-practices-high-level-cnf-expectations"
	TestStorageProvisionerDocLink                      = "https://test-network-function.github.io/cnf-best-practices-guide/#cnf-best-practices-local-storage"
	TestContainerPostStartIdentifierDocLink            = "https://test-network-function.github.io/cnf-best-practices-guide/#cnf-best-practices-cloud-native-design-best-practices"
	TestContainerPrestopIdentifierDocLink              = "https://test-network-function.github.io/cnf-best-practices-guide/#cnf-best-practices-cloud-native-design-best-practices"
	TestPodNodeSelectorAndAffinityBestPracticesDocLink = "https://test-network-function.github.io/cnf-best-practices-guide/#cnf-best-practices-high-level-cnf-expectations"
	TestPodHighAvailabilityBestPracticesDocLink        = "https://test-network-function.github.io/cnf-best-practices-guide/#cnf-best-practices-high-level-cnf-expectations"
	TestPodDeploymentBestPracticesIdentifierDocLink    = "https://test-network-function.github.io/cnf-best-practices-guide/#cnf-best-practices-no-naked-pods"
	TestDeploymentScalingIdentifierDocLink             = "https://test-network-function.github.io/cnf-best-practices-guide/#cnf-best-practices-high-level-cnf-expectations"
	TestStateFulSetScalingIdentifierDocLink            = "https://test-network-function.github.io/cnf-best-practices-guide/#cnf-best-practices-high-level-cnf-expectations"
	TestImagePullPolicyIdentifierDocLink               = "https://test-network-function.github.io/cnf-best-practices-guide/#cnf-best-practices-use-imagepullpolicy-if-not-present"
	TestPodRecreationIdentifierDocLink                 = "https://test-network-function.github.io/cnf-best-practices-guide/#cnf-best-practices-upgrade-expectations"
	TestLivenessProbeIdentifierDocLink                 = "https://test-network-function.github.io/cnf-best-practices-guide/#cnf-best-practices-high-level-cnf-expectations"
	TestReadinessProbeIdentifierDocLink                = "https://test-network-function.github.io/cnf-best-practices-guide/#cnf-best-practices-high-level-cnf-expectations"
	TestStartupProbeIdentifierDocLink                  = "https://test-network-function.github.io/cnf-best-practices-guide/#cnf-best-practices-pod-exit-status"
	//nolint:gosec
	TestPodTolerationBypassIdentifierDocLink           = "https://test-network-function.github.io/cnf-best-practices-guide/#cnf-best-practices-taints-and-tolerations"
	TestPersistentVolumeReclaimPolicyIdentifierDocLink = "https://test-network-function.github.io/cnf-best-practices-guide/#cnf-best-practices-csi"
	TestCPUIsolationIdentifierDocLink                  = "https://test-network-function.github.io/cnf-best-practices-guide/#cnf-best-practices-cpu-isolation"
	TestCrdScalingIdentifierDocLink                    = "https://test-network-function.github.io/cnf-best-practices-guide/#cnf-best-practices-high-level-cnf-expectations"

	// Performance Test Suite
	TestExclusiveCPUPoolIdentifierDocLink       = NoDocLinkFarEdge
	TestSharedCPUPoolSchedulingPolicyDocLink    = NoDocLinkFarEdge
	TestExclusiveCPUPoolSchedulingPolicyDocLink = NoDocLinkFarEdge
	TestIsolatedCPUPoolSchedulingPolicyDocLink  = NoDocLinkFarEdge
	TestRtAppNoExecProbesDocLink                = NoDocLinkFarEdge

	// Operator Test Suite
	DocOperatorRequirement                              = "https://test-network-function.github.io/cnf-best-practices-guide/#cnf-best-practices-cnf-operator-requirements"
	TestOperatorInstallStatusSucceededIdentifierDocLink = DocOperatorRequirement
	TestOperatorNoPrivilegesDocLink                     = DocOperatorRequirement
	TestOperatorIsCertifiedIdentifierDocLink            = DocOperatorRequirement
	TestOperatorIsInstalledViaOLMIdentifierDocLink      = DocOperatorRequirement
	TestOperatorHasSemanticVersioningIdentifierDocLink  = DocOperatorRequirement
	TestOperatorCrdSchemaIdentifierDocLink              = DocOperatorRequirement
	TestOperatorCrdVersioningIdentifierDocLink          = DocOperatorRequirement
	TestOperatorSingleCrdOwnerIdentifierDocLink         = DocOperatorRequirement
	TestOperatorRunAsUserIDDocLink                      = DocOperatorRequirement
	TestOperatorRunAsNonRootDocLink                     = DocOperatorRequirement
	TestOperatorAutomountTokensDocLink                  = DocOperatorRequirement
	TestOperatorReadOnlyFilesystemDocLink               = DocOperatorRequirement

	// Observability Test Suite
	TestLoggingIdentifierDocLink                  = "https://test-network-function.github.io/cnf-best-practices-guide/#cnf-best-practices-logging"
	TestTerminationMessagePolicyIdentifierDocLink = "https://test-network-function.github.io/cnf-best-practices-guide/#cnf-best-practices-pod-exit-status"
	TestCrdsStatusSubresourceIdentifierDocLink    = "https://test-network-function.github.io/cnf-best-practices-guide/#cnf-best-practices-cnf-operator-requirements"
	TestPodDisruptionBudgetIdentifierDocLink      = "https://test-network-function.github.io/cnf-best-practices-guide/#cnf-best-practices-upgrade-expectations"

	// Manageability Test Suite
	TestContainersImageTagDocLink      = "https://test-network-function.github.io/cnf-best-practices-guide/#cnf-best-practices-image-tagging"
	TestContainerPortNameFormatDocLink = "https://test-network-function.github.io/cnf-best-practices-guide/#cnf-best-practices-requirements-cnf-reqs"
)
