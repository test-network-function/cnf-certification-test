package identifiers

const (
	// Default Strings
	NoDocLinkExtended = "No Doc Link - Extended"
	NoDocLinkFarEdge  = "No Doc Link - Far Edge"
	NoDocLinkTelco    = "No Doc Link - Telco"
	NoDocLink         = "No Doc Link"

	// Networking Suite
	TestICMPv4ConnectivityIdentifierDocLink         = "https://redhat-best-practices-for-k8s.github.io/guide/#redhat-best-practices-for-k8s-ipv4-&-ipv6"
	TestNetworkPolicyDenyAllIdentifierDocLink       = "https://redhat-best-practices-for-k8s.github.io/guide/#redhat-best-practices-for-k8s-vrfs-aka-routing-instances"
	TestReservedExtendedPartnerPortsDocLink         = NoDocLinkExtended
	TestDpdkCPUPinningExecProbeDocLink              = "https://redhat-best-practices-for-k8s.github.io/guide/#redhat-best-practices-for-k8s-cpu-manager-pinning"
	TestRestartOnRebootLabelOnPodsUsingSRIOVDocLink = NoDocLinkFarEdge
	TestLimitedUseOfExecProbesIdentifierDocLink     = NoDocLinkFarEdge
	TestICMPv6ConnectivityIdentifierDocLink         = "https://redhat-best-practices-for-k8s.github.io/guide/#redhat-best-practices-for-k8s-ipv4-&-ipv6"
	TestICMPv4ConnectivityMultusIdentifierDocLink   = "https://redhat-best-practices-for-k8s.github.io/guide/#redhat-best-practices-for-k8s-high-level-cnf-expectations"
	TestICMPv6ConnectivityMultusIdentifierDocLink   = "https://redhat-best-practices-for-k8s.github.io/guide/#redhat-best-practices-for-k8s-high-level-cnf-expectations"
	TestServiceDualStackIdentifierDocLink           = "https://redhat-best-practices-for-k8s.github.io/guide/#redhat-best-practices-for-k8s-ipv4-&-ipv6"
	TestUndeclaredContainerPortsUsageDocLink        = "https://redhat-best-practices-for-k8s.github.io/guide/#redhat-best-practices-for-k8s-requirements-cnf-reqs"
	TestOCPReservedPortsUsageDocLink                = "https://redhat-best-practices-for-k8s.github.io/guide/#redhat-best-practices-for-k8s-ports-reserved-by-openshift"

	// Access Control Suite
	Test1337UIDIdentifierDocLink                             = NoDocLinkExtended
	TestNetAdminIdentifierDocLink                            = "https://redhat-best-practices-for-k8s.github.io/guide/#redhat-best-practices-for-k8s-net_admin"
	TestSysAdminIdentifierDocLink                            = "https://redhat-best-practices-for-k8s.github.io/guide/#redhat-best-practices-for-k8s-avoid-sys_admin"
	TestIpcLockIdentifierDocLink                             = "https://redhat-best-practices-for-k8s.github.io/guide/#redhat-best-practices-for-k8s-ipc_lock"
	TestNetRawIdentifierDocLink                              = "https://redhat-best-practices-for-k8s.github.io/guide/#redhat-best-practices-for-k8s-user-plane-cnfs"
	TestBpfIdentifierDocLink                                 = NoDocLinkTelco
	TestSecConNonRootUserIdentifierDocLink                   = "https://redhat-best-practices-for-k8s.github.io/guide/#redhat-best-practices-for-k8s-cnf-security"
	TestSecContextIdentifierDocLink                          = "https://redhat-best-practices-for-k8s.github.io/guide/#redhat-best-practices-for-k8s-cnf-security"
	TestSecConPrivilegeEscalationDocLink                     = "https://redhat-best-practices-for-k8s.github.io/guide/#redhat-best-practices-for-k8s-cnf-security"
	TestContainerHostPortDocLink                             = "https://redhat-best-practices-for-k8s.github.io/guide/#redhat-best-practices-for-k8s-avoid-accessing-resource-on-host"
	TestContainerHostNetworkDocLink                          = "https://redhat-best-practices-for-k8s.github.io/guide/#redhat-best-practices-for-k8s-avoid-the-host-network-namespace"
	TestPodHostNetworkDocLink                                = "https://redhat-best-practices-for-k8s.github.io/guide/#redhat-best-practices-for-k8s-avoid-the-host-network-namespace"
	TestPodHostPathDocLink                                   = "https://redhat-best-practices-for-k8s.github.io/guide/#redhat-best-practices-for-k8s-cnf-security"
	TestPodHostIPCDocLink                                    = "https://redhat-best-practices-for-k8s.github.io/guide/#redhat-best-practices-for-k8s-cnf-security"
	TestPodHostPIDDocLink                                    = "https://redhat-best-practices-for-k8s.github.io/guide/#redhat-best-practices-for-k8s-cnf-security"
	TestNamespaceBestPracticesIdentifierDocLink              = "https://redhat-best-practices-for-k8s.github.io/guide/#redhat-best-practices-for-k8s-requirements-cnf-reqs"
	TestPodClusterRoleBindingsBestPracticesIdentifierDocLink = "https://redhat-best-practices-for-k8s.github.io/guide/#redhat-best-practices-for-k8s-security-rbac"
	TestPodRoleBindingsBestPracticesIdentifierDocLink        = "https://redhat-best-practices-for-k8s.github.io/guide/#redhat-best-practices-for-k8s-security-rbac"
	TestPodServiceAccountBestPracticesIdentifierDocLink      = "https://redhat-best-practices-for-k8s.github.io/guide/#redhat-best-practices-for-k8s-scc-permissions-for-an-application"
	TestPodAutomountServiceAccountIdentifierDocLink          = "https://redhat-best-practices-for-k8s.github.io/guide/#redhat-best-practices-for-k8s-automount-services-for-pods"
	TestServicesDoNotUseNodeportsIdentifierDocLink           = "https://redhat-best-practices-for-k8s.github.io/guide/#redhat-best-practices-for-k8s-avoid-the-host-network-namespace"
	TestUnalteredBaseImageIdentifierDocLink                  = "https://redhat-best-practices-for-k8s.github.io/guide/#redhat-best-practices-for-k8s-image-standards"
	TestOneProcessPerContainerIdentifierDocLink              = "https://redhat-best-practices-for-k8s.github.io/guide/#redhat-best-practices-for-k8s-one-process-per-container"
	TestSYSNiceRealtimeCapabilityIdentifierDocLink           = "https://redhat-best-practices-for-k8s.github.io/guide/#redhat-best-practices-for-k8s-sys_nice"
	TestSysPtraceCapabilityIdentifierDocLink                 = "https://redhat-best-practices-for-k8s.github.io/guide/#redhat-best-practices-for-k8s-sys_ptrace"
	TestPodRequestsAndLimitsIdentifierDocLink                = "https://redhat-best-practices-for-k8s.github.io/guide/#redhat-best-practices-for-k8s-requests/limits"
	TestNamespaceResourceQuotaIdentifierDocLink              = "https://redhat-best-practices-for-k8s.github.io/guide/#redhat-best-practices-for-k8s-memory-allocation"
	TestNoSSHDaemonsAllowedIdentifierDocLink                 = "https://redhat-best-practices-for-k8s.github.io/guide/#redhat-best-practices-for-k8s-pod-interaction/configuration"

	// Affiliated Certification Suite
	TestHelmVersionIdentifierDocLink                = "https://redhat-best-practices-for-k8s.github.io/guide/#redhat-best-practices-for-k8s-helm"
	TestContainerIsCertifiedDigestIdentifierDocLink = "https://redhat-connect.gitbook.io/partner-guide-for-red-hat-openshift-and-container/certify-your-application/overview"
	TestContainerIsCertifiedIdentifierDocLink       = "https://redhat-connect.gitbook.io/partner-guide-for-red-hat-openshift-and-container/certify-your-application/overview"
	TestHelmIsCertifiedIdentifierDocLink            = "https://redhat-connect.gitbook.io/partner-guide-for-red-hat-openshift-and-container/certify-your-application/overview"

	// Platform Alteration Suite
	TestPodHugePages2MDocLink                       = "https://redhat-best-practices-for-k8s.github.io/guide/#redhat-best-practices-for-k8s-huge-pages"
	TestPodHugePages1GDocLink                       = NoDocLinkFarEdge
	TestHugepagesNotManuallyManipulatedDocLink      = "https://redhat-best-practices-for-k8s.github.io/guide/#redhat-best-practices-for-k8s-huge-pages"
	TestNonTaintedNodeKernelsIdentifierDocLink      = "https://redhat-best-practices-for-k8s.github.io/guide/#redhat-best-practices-for-k8s-high-level-cnf-expectations"
	TestUnalteredStartupBootParamsIdentifierDocLink = "https://redhat-best-practices-for-k8s.github.io/guide/#redhat-best-practices-for-k8s-host-os"
	TestSysctlConfigsIdentifierDocLink              = "https://redhat-best-practices-for-k8s.github.io/guide/#redhat-best-practices-for-k8s-cnf-security"
	TestServiceMeshIdentifierDocLink                = NoDocLinkExtended
	TestHyperThreadEnableDocLink                    = NoDocLinkExtended

	TestOCPLifecycleIdentifierDocLink        = "https://redhat-best-practices-for-k8s.github.io/guide/#redhat-best-practices-for-k8s-k8s"
	TestNodeOperatingSystemIdentifierDocLink = "https://redhat-best-practices-for-k8s.github.io/guide/#redhat-best-practices-for-k8s-host-os"
	TestIsRedHatReleaseIdentifierDocLink     = "https://redhat-best-practices-for-k8s.github.io/guide/#redhat-best-practices-for-k8s-base-images"
	TestIsSELinuxEnforcingIdentifierDocLink  = "https://redhat-best-practices-for-k8s.github.io/guide/#redhat-best-practices-for-k8s-pod-security"

	// Lifecycle Suite
	TestAffinityRequiredPodsDocLink                    = "https://redhat-best-practices-for-k8s.github.io/guide/#redhat-best-practices-for-k8s-high-level-cnf-expectations"
	TestStorageProvisionerDocLink                      = "https://redhat-best-practices-for-k8s.github.io/guide/#redhat-best-practices-for-k8s-local-storage"
	TestContainerPostStartIdentifierDocLink            = "https://redhat-best-practices-for-k8s.github.io/guide/#redhat-best-practices-for-k8s-cloud-native-design-best-practices"
	TestContainerPrestopIdentifierDocLink              = "https://redhat-best-practices-for-k8s.github.io/guide/#redhat-best-practices-for-k8s-cloud-native-design-best-practices"
	TestPodNodeSelectorAndAffinityBestPracticesDocLink = "https://redhat-best-practices-for-k8s.github.io/guide/#redhat-best-practices-for-k8s-high-level-cnf-expectations"
	TestPodHighAvailabilityBestPracticesDocLink        = "https://redhat-best-practices-for-k8s.github.io/guide/#redhat-best-practices-for-k8s-high-level-cnf-expectations"
	TestPodDeploymentBestPracticesIdentifierDocLink    = "https://redhat-best-practices-for-k8s.github.io/guide/#redhat-best-practices-for-k8s-no-naked-pods"
	TestDeploymentScalingIdentifierDocLink             = "https://redhat-best-practices-for-k8s.github.io/guide/#redhat-best-practices-for-k8s-high-level-cnf-expectations"
	TestStateFulSetScalingIdentifierDocLink            = "https://redhat-best-practices-for-k8s.github.io/guide/#redhat-best-practices-for-k8s-high-level-cnf-expectations"
	TestImagePullPolicyIdentifierDocLink               = "https://redhat-best-practices-for-k8s.github.io/guide/#redhat-best-practices-for-k8s-use-imagepullpolicy-if-not-present"
	TestPodRecreationIdentifierDocLink                 = "https://redhat-best-practices-for-k8s.github.io/guide/#redhat-best-practices-for-k8s-upgrade-expectations"
	TestLivenessProbeIdentifierDocLink                 = "https://redhat-best-practices-for-k8s.github.io/guide/#redhat-best-practices-for-k8s-high-level-cnf-expectations"
	TestReadinessProbeIdentifierDocLink                = "https://redhat-best-practices-for-k8s.github.io/guide/#redhat-best-practices-for-k8s-high-level-cnf-expectations"
	TestStartupProbeIdentifierDocLink                  = "https://redhat-best-practices-for-k8s.github.io/guide/#redhat-best-practices-for-k8s-pod-exit-status"
	//nolint:gosec
	TestPodTolerationBypassIdentifierDocLink           = "https://redhat-best-practices-for-k8s.github.io/guide/#redhat-best-practices-for-k8s-taints-and-tolerations"
	TestPersistentVolumeReclaimPolicyIdentifierDocLink = "https://redhat-best-practices-for-k8s.github.io/guide/#redhat-best-practices-for-k8s-csi"
	TestCPUIsolationIdentifierDocLink                  = "https://redhat-best-practices-for-k8s.github.io/guide/#redhat-best-practices-for-k8s-cpu-isolation"
	TestCrdScalingIdentifierDocLink                    = "https://redhat-best-practices-for-k8s.github.io/guide/#redhat-best-practices-for-k8s-high-level-cnf-expectations"

	// Performance Test Suite
	TestExclusiveCPUPoolIdentifierDocLink       = NoDocLinkFarEdge
	TestSharedCPUPoolSchedulingPolicyDocLink    = NoDocLinkFarEdge
	TestExclusiveCPUPoolSchedulingPolicyDocLink = NoDocLinkFarEdge
	TestIsolatedCPUPoolSchedulingPolicyDocLink  = NoDocLinkFarEdge
	TestRtAppNoExecProbesDocLink                = NoDocLinkFarEdge

	// Operator Test Suite
	DocOperatorRequirement                              = "https://redhat-best-practices-for-k8s.github.io/guide/#redhat-best-practices-for-k8s-cnf-operator-requirements"
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
	TestLoggingIdentifierDocLink                            = "https://redhat-best-practices-for-k8s.github.io/guide/#redhat-best-practices-for-k8s-logging"
	TestTerminationMessagePolicyIdentifierDocLink           = "https://redhat-best-practices-for-k8s.github.io/guide/#redhat-best-practices-for-k8s-pod-exit-status"
	TestCrdsStatusSubresourceIdentifierDocLink              = "https://redhat-best-practices-for-k8s.github.io/guide/#redhat-best-practices-for-k8s-cnf-operator-requirements"
	TestPodDisruptionBudgetIdentifierDocLink                = "https://redhat-best-practices-for-k8s.github.io/guide/#redhat-best-practices-for-k8s-upgrade-expectations"
	TestAPICompatibilityWithNextOCPReleaseIdentifierDocLink = "https://redhat-best-practices-for-k8s.github.io/guide/#redhat-best-practices-for-k8s-to-be-removed-apis"

	// Manageability Test Suite
	TestContainersImageTagDocLink      = "https://redhat-best-practices-for-k8s.github.io/guide/#redhat-best-practices-for-k8s-image-tagging"
	TestContainerPortNameFormatDocLink = "https://redhat-best-practices-for-k8s.github.io/guide/#redhat-best-practices-for-k8s-requirements-cnf-reqs"
)
