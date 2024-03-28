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

package identifiers

const (
	//nolint:gosec
	AutomountServiceTokenRemediation = `Check that pod has automountServiceAccountToken set to false or pod is attached to service account which has automountServiceAccountToken set to false, unless the pod needs access to the kubernetes API server. Pods which do not need API access should set automountServiceAccountToken to false in pod spec.`

	IsRedHatReleaseRemediation = `Build a new container image that is based on UBI (Red Hat Universal Base Image).`

	NodeOperatingSystemRemediation = `Please update your workers to a version that is supported by your version of OpenShift`

	SecConNonRootUserRemediation = `Change the pod and containers "runAsUser" uid to something other than root(0)`

	SecConRemediation = `Exception possible if CNF uses mlock(), mlockall(), shmctl(), mmap(); exception will be considered for DPDK applications. Must identify which container requires the capability and detail why.`

	UnalteredBaseImageRemediation = `Ensure that Container applications do not modify the Container Base Image. In particular, ensure that the following directories are not modified: 1) /var/lib/rpm 2) /var/lib/dpkg 3) /bin 4) /sbin 5) /lib 6) /lib64 7) /usr/bin 8) /usr/sbin 9) /usr/lib 10) /usr/lib64 Ensure that all required binaries are built directly into the container image, and are not installed post startup.`

	OCPLifecycleRemediation = `Please update your cluster to a version that is generally available.`

	DeploymentScalingRemediation = `Ensure CNF deployments/replica sets can scale in/out successfully.`
	CrdScalingRemediation        = `Ensure CNF crd/replica sets can scale in/out successfully.`

	StatefulSetScalingRemediation = `Ensure CNF statefulsets/replica sets can scale in/out successfully.`

	SecConCapabilitiesRemediation = `Remove the following capabilities from the container/pod definitions: NET_ADMIN SCC, SYS_ADMIN SCC, NET_RAW SCC, IPC_LOCK SCC`

	BpfCapabilityRemediation = `Remove the following capability from the container/pod definitions: BPF`

	SecConPrivilegeRemediation = `Configure privilege escalation to false. Privileged escalation should not be allowed (AllowPrivilegeEscalation=false).`

	ContainerHostPortRemediation = `Remove hostPort configuration from the container. CNF should avoid accessing host resources - containers should not configure HostPort.`

	PodHostNetworkRemediation = `Set the spec.HostNetwork parameter to false in the pod configuration. CNF should avoid accessing host resources - spec.HostNetwork should be false.`

	PodHostPathRemediation = `Set the spec.HostPath parameter to false in the pod configuration. CNF should avoid accessing host resources - spec.HostPath should be false.`

	PodHostIPCRemediation = `Set the spec.HostIpc parameter to false in the pod configuration. CNF should avoid accessing host resources - spec.HostIpc should be false.`

	PodHostPIDRemediation = `Set the spec.HostPid parameter to false in the pod configuration. CNF should avoid accessing host resources - spec.HostPid should be false.`

	HugepagesNotManuallyManipulatedRemediation = `HugePage settings should be configured either directly through the MachineConfigOperator or indirectly using the PerformanceAddonOperator. This ensures that OpenShift is aware of the special MachineConfig requirements, and can provision your CNF on a Node that is part of the corresponding MachineConfigSet. Avoid making changes directly to an underlying Node, and let OpenShift handle the heavy lifting of configuring advanced settings. This test case applies only to Nodes that are configured with the "worker" MachineConfigSet.`

	ICMPv4ConnectivityRemediation = `Ensure that the CNF is able to communicate via the Default OpenShift network. In some rare cases, CNFs may require routing table changes in order to communicate over the Default network. To exclude a particular pod from ICMPv4 connectivity tests, add the test-network-function.com/skip_connectivity_tests label to it. The label value is trivial, only its presence.`

	ICMPv6ConnectivityRemediation = `Ensure that the CNF is able to communicate via the Default OpenShift network. In some rare cases, CNFs may require routing table changes in order to communicate over the Default network. To exclude a particular pod from ICMPv6 connectivity tests, add the test-network-function.com/skip_connectivity_tests label to it. The label value is trivial, only its presence. Not applicable if IPv6 is not supported.`

	ICMPv4ConnectivityMultusRemediation = `Ensure that the CNF is able to communicate via the Multus network(s). In some rare cases, CNFs may require routing table changes in order to communicate over the Multus network(s). To exclude a particular pod from ICMPv4 connectivity tests, add the test-network-function.com/skip_connectivity_tests label to it. The label value is trivial, only its presence. Not applicable if MULTUS is not supported.`

	ICMPv6ConnectivityMultusRemediation = `Ensure that the CNF is able to communicate via the Multus network(s). In some rare cases, CNFs may require routing table changes in order to communicate over the Multus network(s). To exclude a particular pod from ICMPv6 connectivity tests, add the test-network-function.com/skip_connectivity_tests label to it.The label value is trivial, only its presence.`

	TestServiceDualStackRemediation = `Configure every CNF services with either a single stack ipv6 or dual stack (ipv4/ipv6) load balancer.`

	NamespaceBestPracticesRemediation = `Ensure that your CNF utilizes namespaces declared in the yaml config file. Additionally, the namespaces should not start with "default, openshift-, istio- or aspenmesh-".`

	NonTaintedNodeKernelsRemediation = `Test failure indicates that the underlying Node's kernel is tainted. Ensure that you have not altered underlying Node(s) kernels in order to run the CNF.`

	OperatorInstallStatusSucceededRemediation = `Ensure all the CNF operators have been successfully installed by OLM.`

	OperatorNoPrivilegesRemediation = `Ensure all the CNF operators have no privileges on cluster resources.`

	OperatorIsCertifiedRemediation = `Ensure that your Operator has passed Red Hat's Operator Certification Program (OCP).`

	HelmIsCertifiedRemediation = `Ensure that the helm charts under test passed the Red Hat's helm Certification Program (e.g. listed in https://charts.openshift.io/index.yaml).`

	OperatorIsInstalledViaOLMRemediation = `Ensure that your Operator is installed via OLM.`

	OperatorHasSemanticVersioningRemediation = `Ensure that the installed Operator has semantic versioning.`

	PodNodeSelectorAndAffinityBestPracticesRemediation = `In most cases, Pod's should not specify their host Nodes through nodeSelector or nodeAffinity. However, there are cases in which CNFs require specialized hardware specific to a particular class of Node.`

	PodHighAvailabilityBestPracticesRemediation = `In high availability cases, Pod podAntiAffinity rule should be specified for pod scheduling and pod replica value is set to more than 1 .`

	PodClusterRoleBindingsBestPracticesRemediation = `In most cases, Pod's should not have ClusterRoleBindings. The suggested remediation is to remove the need for ClusterRoleBindings, if possible. Cluster roles and cluster role bindings discouraged unless absolutely needed by CNF (often reserved for cluster admin only).`

	PodDeploymentBestPracticesRemediation = `Deploy the CNF using ReplicaSet/StatefulSet.`

	ImagePullPolicyRemediation = `Ensure that the containers under test are using IfNotPresent as Image Pull Policy.`

	PodRoleBindingsBestPracticesRemediation = `Ensure the CNF is not configured to use RoleBinding(s) in a non-CNF Namespace. Scope of role must <= scope of creator of role.`

	PodServiceAccountBestPracticesRemediation = `Ensure that the each CNF Pod is configured to use a valid Service Account`

	ServicesDoNotUseNodeportsRemediation = `Ensure Services are not configured to use NodePort(s).CNF should avoid accessing host resources - tests that each CNF Service does not utilize NodePort(s).`

	UnalteredStartupBootParamsRemediation = `Ensure that boot parameters are set directly through the MachineConfigOperator, or indirectly through the PerformanceAddonOperator. Boot parameters should not be changed directly through the Node, as OpenShift should manage the changes for you.`

	PodRecreationRemediation = `Ensure that CNF Pod(s) utilize a configuration that supports High Availability. Additionally, ensure that there are available Nodes in the OpenShift cluster that can be utilized in the event that a host Node fails.`

	SysctlConfigsRemediation = `You should recreate the node or change the sysctls, recreating is recommended because there might be other unknown changes`

	ServiceMeshRemediation = `Ensure all the CNF pods are using service mesh if the cluster provides it.`

	ScalingRemediation = `Ensure CNF deployments/replica sets can scale in/out successfully.`

	IsSELinuxEnforcingRemediation = `Configure selinux and enable enforcing mode.`

	UndeclaredContainerPortsRemediation = `Ensure the CNF apps do not listen on undeclared containers' ports.`

	CrdsStatusSubresourceRemediation = `Ensure that all the CRDs have a meaningful status specification (Spec.versions[].Schema.OpenAPIV3Schema.Properties[“status”]).`

	LoggingRemediation = `Ensure containers are not redirecting stdout/stderr`

	TerminationMessagePolicyRemediation = `Ensure containers are all using FallbackToLogsOnError in terminationMessagePolicy`

	LivenessProbeRemediation = `Add a liveness probe to deployed containers.`

	ReadinessProbeRemediation = `Add a readiness probe to deployed containers`

	StartupProbeRemediation = `Add a startup probe to deployed containers`

	OneProcessPerContainerRemediation = `Launch only one process per container. Should adhere to 1 process per container best practice wherever possible.`

	SysPtraceCapabilityRemediation = `Allow the SYS_PTRACE capability when enabling process namespace sharing for a Pod`

	SYSNiceRealtimeCapabilityRemediation = `If pods are scheduled to realtime kernel nodes, they must add SYS_NICE capability to their spec.`

	OCPReservedPortsUsageRemediation = `Ensure that CNF apps do not listen on ports that are reserved by OpenShift. The following ports are reserved by OpenShift and must NOT be used by any application: 22623, 22624.`

	RequestsAndLimitsRemediation = `Add requests and limits to your container spec. See: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/#requests-and-limits`

	NamespaceResourceQuotaRemediation = `Apply a ResourceQuota to the namespace your CNF is running in. The CNF namespace should have resource quota defined.`

	PodDisruptionBudgetRemediation = `Ensure minAvailable is not zero and maxUnavailable does not equal the number of pods in the replica`

	//nolint:gosec
	PodTolerationBypassRemediation = `Do not allow pods to bypass the NoExecute, PreferNoSchedule, or NoSchedule tolerations that are default applied by Kubernetes.`

	PersistentVolumeReclaimPolicyRemediation = `Ensure that all persistent volumes are using the reclaim policy: delete`

	ContainersImageTagRemediation = `Ensure that all the container images are tagged. Checks containers have image tags (e.g. latest, stable, dev).`

	NoSSHDaemonsAllowedRemediation = `Ensure that no SSH daemons are running inside a pod. Pods should not run as SSH Daemons (replicaset or statefulset only).`

	NetworkPolicyDenyAllRemediation = `Ensure that a NetworkPolicy with a default deny-all is applied. After the default is applied, apply a network policy to allow the traffic your application requires.`

	CPUIsolationRemediation = `CPU isolation testing is enabled. Please ensure that all pods adhere to the CPU isolation requirements.`

	UID1337Remediation = `Use another process UID that is not 1337.`

	LimitedUseOfExecProbesRemediation = `Reduce the number of exec probes in the cluster for this CNF to less than 10. Increase the update period of the exec probe to be superior or equal to 10 seconds.`

	ReservedPartnerPortsRemediation = `Ensure ports are not being used that are reserved by our partner`

	AffinityRequiredRemediation = `Pods which need to be co-located on the same node need Affinity rules. If a pod/statefulset/deployment is required to use affinity rules, please add AffinityRequired: 'true' as a label.`

	ContainerPortNameFormatRemediation = `Ensure that the container's ports name follow our partner naming conventions`

	DpdkCPUPinningExecProbeRemediation = "If the CNF is doing CPU pinning and running a DPDK process do not use exec probes (executing a command within the container) as it may pile up and block the node eventually."

	CheckStorageProvisionerRemediation = `Use a non-local storage (e.g. no kubernetes.io/no-provisioner and no topolvm.io provisioners) in multinode clusters. Local storage are recommended for single node clusters only, but a single local provisioner should be installed.`

	ExclusiveCPUPoolRemediation = `Ensure that if one container in a Pod selects an exclusive CPU pool the rest also select this type of CPU pool`

	SharedCPUPoolSchedulingPolicyRemediation = `Ensure that the workload running in Application shared CPU pool should choose non-RT CPU schedule policy, like SCHED _OTHER to always share the CPU with other applications and kernel threads.`

	ExclusiveCPUPoolSchedulingPolicyRemediation = `Ensure that the workload running in Application exclusive CPU pool can choose RT CPU scheduling policy, but should set priority less than 10`

	IsolatedCPUPoolSchedulingPolicyRemediation = `Ensure that the workload running in an application-isolated exclusive CPU pool selects a RT CPU scheduling policy (such as SCHED_FIFO/SCHED_RR) with High priority.`

	RtAppNoExecProbesRemediation = `Ensure that if one container runs a real time application exec probes are not used`

	SRIOVPodsRestartOnRebootLabelRemediation = `Ensure that the label restart-on-reboot exists on pods that use SRIOV network interfaces.`

	HelmVersionV3Remediation = `Check Helm Chart is v3 and not v2 which is not supported due to security risks associated with Tiller.`

	ContainerIsCertifiedDigestRemediation = "Ensure that your container has passed the Red Hat Container Certification Program (CCP)."

	PodHugePages2MRemediation = "Modify pod to consume 2Mi hugepages only"

	PodHugePages1GRemediation = "Modify pod to consume 1Gi hugepages only"

	HyperThreadEnable = "Check that baremetal workers have hyperthreading enabled"

	ContainerPostStartIdentifierRemediation = `Identify which pod is not conforming to the process and submit information as to why it cannot use a postStart startup specification.`

	ContainerPrestopIdentifierRemediation = `Identify which pod is not conforming to the process and submit information as to why it cannot use a preStop shutdown specification.`
)
