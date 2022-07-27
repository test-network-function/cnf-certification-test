// Copyright (C) 2022 Red Hat, Inc.
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
	AutomountServiceTokenRemediation = `Check that pod has automountServiceAccountToken set to false or pod is attached to service account which has automountServiceAccountToken set to false`

	ShutdownRemediation = `
	It's considered best-practices to define prestop for proper management of container lifecycle.
	The prestop can be used to gracefully stop the container and clean resources (e.g., DB connection).
	
	The prestop can be configured using :
	 1) Exec : executes the supplied command inside the container
	 2) HTTP : executes HTTP request against the specified endpoint.
	
	When defined. K8s will handle shutdown of the container using the following:
	1) K8s first execute the preStop hook inside the container.
	2) K8s will wait for a grace period.
	3) K8s will clean the remaining processes using KILL signal.		
		`

	IsRedHatReleaseRemediation = `Build a new container image that is based on UBI (Red Hat Universal Base Image).`

	NodeOperatingSystemRemediation = `Please update your workers to a version that is supported by your version of OpenShift`

	SecConNonRootUserRemediation = `Change the pod and containers "runAsUser" uid to something other than root(0)`

	UnalteredBaseImageRemediation = `Ensure that Container applications do not modify the Container Base Image.  In particular, ensure that the following
	directories are not modified:
	1) /var/lib/rpm
	2) /var/lib/dpkg
	3) /bin
	4) /sbin
	5) /lib
	6) /lib64
	7) /usr/bin
	8) /usr/sbin
	9) /usr/lib
	10) /usr/lib64
	Ensure that all required binaries are built directly into the container image, and are not installed post startup.`

	OCPLifecycleRemediation = `Please update your cluster to a version that is generally available.`

	DeploymentScalingRemediation = `Ensure CNF deployments/replica sets can scale in/out successfully.`

	StatefulSetScalingRemediation = `Ensure CNF statefulsets/replica sets can scale in/out successfully.`

	SecConCapabilitiesRemediation = `Remove the following capabilities from the container/pod definitions: NET_ADMIN SCC, SYS_ADMIN SCC, NET_RAW SCC, IPC_LOCK SCC`

	SecConPrivilegeRemediation = `Configure privilege escalation to false`

	ContainerIsCertifiedRemediation = `Ensure that your container has passed the Red Hat Container Certification Program (CCP).`

	ContainerHostPortRemediation = `Remove hostPort configuration from the container`

	PodHostNetworkRemediation = `Set the spec.HostNetwork parameter to false in the pod configuration`

	PodHostPathRemediation = `Set the spec.HostPath parameter to false in the pod configuration`

	PodHostIPCRemediation = `Set the spec.HostIpc parameter to false in the pod configuration`

	PodHostPIDRemediation = `Set the spec.HostPid parameter to false in the pod configuration`

	HugepagesNotManuallyManipulatedRemediation = `HugePage settings should be configured either directly through the MachineConfigOperator or indirectly using the
	PerformanceAddonOperator.  This ensures that OpenShift is aware of the special MachineConfig requirements, and can
	provision your CNF on a Node that is part of the corresponding MachineConfigSet.  Avoid making changes directly to an
	underlying Node, and let OpenShift handle the heavy lifting of configuring advanced settings.`

	ICMPv4ConnectivityRemediation = `Ensure that the CNF is able to communicate via the Default OpenShift network. In some rare cases,
	CNFs may require routing table changes in order to communicate over the Default network. To exclude a particular pod
	from ICMPv4 connectivity tests, add the test-network-function.com/skip_connectivity_tests label to it. The label value is not important, only its presence.`

	ICMPv6ConnectivityRemediation = `Ensure that the CNF is able to communicate via the Default OpenShift network. In some rare cases,
	CNFs may require routing table changes in order to communicate over the Default network. To exclude a particular pod
	from ICMPv6 connectivity tests, add the test-network-function.com/skip_connectivity_tests label to it. The label value is not important, only its presence.`

	ICMPv4ConnectivityMultusRemediation = `Ensure that the CNF is able to communicate via the Multus network(s). In some rare cases,
	CNFs may require routing table changes in order to communicate over the Multus network(s). To exclude a particular pod
	from ICMPv4 connectivity tests, add the test-network-function.com/skip_connectivity_tests label to it. The label value is not important, only its presence.`

	ICMPv6ConnectivityMultusRemediation = `Ensure that the CNF is able to communicate via the Multus network(s). In some rare cases,
	CNFs may require routing table changes in order to communicate over the Multus network(s). To exclude a particular pod
	from ICMPv6 connectivity tests, add the test-network-function.com/skip_connectivity_tests label to it.The label value is not important, only its presence.`

	NamespaceBestPracticesRemediation = `Ensure that your CNF utilizes namespaces declared in the yaml config file. Additionally,
	the namespaces should not start with "default, openshift-, istio- or aspenmesh-", except in rare cases.`

	NonTaintedNodeKernelsRemediation = `Test failure indicates that the underlying Node's kernel is tainted.  Ensure that you have not altered underlying
	Node(s) kernels in order to run the CNF.`

	OperatorInstallStatusSucceededRemediation = `Ensure all the CNF operators have been successfully installed by OLM.`

	OperatorNoPrivilegesRemediation = `Ensure all the CNF operators have no privileges on cluster resources.`

	OperatorIsCertifiedRemediation = `Ensure that your Operator has passed Red Hat's Operator Certification Program (OCP).`

	HelmIsCertifiedRemediation = `Ensure that the helm charts under test passed the Red Hat's helm Certification Program (e.g. listed in https://charts.openshift.io/index.yaml).`

	OperatorIsInstalledViaOLMRemediation = `Ensure that your Operator is installed via OLM.`

	PodNodeSelectorAndAffinityBestPracticesRemediation = `In most cases, Pod's should not specify their host Nodes through nodeSelector or nodeAffinity.  However, there are
	cases in which CNFs require specialized hardware specific to a particular class of Node.  As such, this test is purely
	informative, and will not prevent a CNF from being certified. However, one should have an appropriate justification as
	to why nodeSelector and/or nodeAffinity is utilized by a CNF.`

	PodHighAvailabilityBestPracticesRemediation = `In high availability cases, Pod podAntiAffinity rule should be specified for pod scheduling and pod replica value is set to more than 1 .`

	PodClusterRoleBindingsBestPracticesRemediation = `In most cases, Pod's should not have ClusterRoleBindings.  The suggested remediation is to remove the need for
	ClusterRoleBindings, if possible.`

	PodDeploymentBestPracticesRemediation = `Deploy the CNF using ReplicaSet/StatefulSet.`

	ImagePullPolicyRemediation = `Ensure that the containers under test are using IfNotPresent as Image Pull Policy.`

	PodRoleBindingsBestPracticesRemediation = `Ensure the CNF is not configured to use RoleBinding(s) in a non-CNF Namespace.`

	PodServiceAccountBestPracticesRemediation = `Ensure that the each CNF Pod is configured to use a valid Service Account`

	ServicesDoNotUseNodeportsRemediation = `Ensure Services are not configured to use NodePort(s).`

	UnalteredStartupBootParamsRemediation = `Ensure that boot parameters are set directly through the MachineConfigOperator, or indirectly through the PerformanceAddonOperator.  
	Boot parameters should not be changed directly through the Node, as OpenShift should manage the changes for you.`

	PodRecreationRemediation = `Ensure that CNF Pod(s) utilize a configuration that supports High Availability.  
	Additionally, ensure that there are available Nodes in the OpenShift cluster that can be utilized in the event that a host Node fails.`

	SysctlConfigsRemediation = `You should recreate the node or change the sysctls, recreating is recommended because there might be other unknown changes`

	ServiceMeshRemediation = `Ensure all the CNF pods are using service mesh if the cluster provides it.`

	ScalingRemediation = `Ensure CNF deployments/replica sets can scale in/out successfully.`

	IsSELinuxEnforcingRemediation = `Configure selinux and enable enforcing mode.`

	UndeclaredContainerPortsRemediation = `Ensure the CNF apps do not listen on undeclared containers' ports`

	CrdsStatusSubresourceRemediation = `Ensure that all the CRDs have a meaningful status specification (Spec.versions[].Schema.OpenAPIV3Schema.Properties["status"]).`

	LoggingRemediation = `Ensure containers are not redirecting stdout/stderr`

	TerminationMessagePolicyRemediation = `Ensure containers are all using FallbackToLogsOnError in terminationMessagePolicy`

	LivenessProbeRemediation = `Add liveness probe to deployed containers`

	ReadinessProbeRemediation = `Add readiness probe to deployed containers`

	OneProcessPerContainerRemediation = `Launch only one process per container`

	SysPtraceCapabilityRemediation = `Allow the SYS_PTRACE capability when enabling process namespace sharing for a Pod`

	SYSNiceRealtimeCapabilityRemediation = `If pods are scheduled to realtime kernel nodes, they must add SYS_NICE capability to their spec.`

	OCPReservedPortsUsageRemediation = `Ensure that CNF apps do not listen on ports that are reserved by Openshift`

	RequestsAndLimitsRemediation = `Add requests and limits to your container spec.  See: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/#requests-and-limits`

	NamespaceResourceQuotaRemediation = `Apply a ResourceQuota to the namespace your CNF is running in`

	NamespaceLimitRangeRemediation = `Apply a LimitRange to the namespace your CNF is running in`
)
