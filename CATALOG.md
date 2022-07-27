# cnf-certification-test catalog

The catalog for cnf-certification-test contains a variety of `Test Cases`, which are traditional JUnit testcases specified internally using `Ginkgo.It`.
## Test Case Catalog

Test Cases are the specifications used to perform a meaningful test.  Test cases may run once, or several times against several targets.  CNF Certification includes a number of normative and informative tests to ensure CNFs follow best practices.  Here is the list of available Test Cases:

### access-control

#### cluster-role-bindings

Property|Description
---|---
Test Case Name|cluster-role-bindings
Test Case Label|access-control-cluster-role-bindings
Unique ID|http://test-network-function.com/testcases/access-control/cluster-role-bindings
Version|v1.0.0
Description|http://test-network-function.com/testcases/access-control/cluster-role-bindings tests that a Pod does not specify ClusterRoleBindings.
Result Type|normative
Suggested Remediation|In most cases, Pod's should not have ClusterRoleBindings.  The suggested remediation is to remove the need for 	ClusterRoleBindings, if possible.
Best Practice Reference|https://connect.redhat.com/sites/default/files/2022-05/Cloud%20Native%20Network%20Function%20Requirements%201-3.pdf Section 5.2.10 and 5.3.6
Exception Process|There is no documented exception process for this.
#### container-host-port

Property|Description
---|---
Test Case Name|container-host-port
Test Case Label|access-control-container-host-port
Unique ID|http://test-network-function.com/testcases/access-control/container-host-port
Version|v1.0.0
Description|http://test-network-function.com/testcases/access-control/container-host-port Verifies if containers define a hostPort.
Result Type|informative
Suggested Remediation|Remove hostPort configuration from the container
Best Practice Reference|https://connect.redhat.com/sites/default/files/2022-05/Cloud%20Native%20Network%20Function%20Requirements%201-3.pdf Section 5.3.6
Exception Process|There is no documented exception process for this.
#### namespace

Property|Description
---|---
Test Case Name|namespace
Test Case Label|access-control-namespace
Unique ID|http://test-network-function.com/testcases/access-control/namespace
Version|v1.0.0
Description|http://test-network-function.com/testcases/access-control/namespace tests that all CNF's resources (PUTs and CRs) belong to valid namespaces. A valid namespace meets the following conditions: (1) It was declared in the yaml config file under the targetNameSpaces tag. (2) It doesn't have any of the following prefixes: default, openshift-, istio- and aspenmesh-
Result Type|normative
Suggested Remediation|Ensure that your CNF utilizes namespaces declared in the yaml config file. Additionally, 	the namespaces should not start with "default, openshift-, istio- or aspenmesh-", except in rare cases.
Best Practice Reference|https://connect.redhat.com/sites/default/files/2022-05/Cloud%20Native%20Network%20Function%20Requirements%201-3.pdf Section 5.2, 16.3.8 and 16.3.9
Exception Process|There is no documented exception process for this.
#### namespace-limit-range

Property|Description
---|---
Test Case Name|namespace-limit-range
Test Case Label|access-control-namespace-limit-range
Unique ID|http://test-network-function.com/testcases/access-control/namespace-limit-range
Version|v1.0.0
Description|http://test-network-function.com/testcases/access-control/namespace-limit-range Checks to see if CNF workload pods are running in namespaces that have limit ranges applied.
Result Type|informative
Suggested Remediation|Apply a LimitRange to the namespace your CNF is running in
Best Practice Reference|https://connect.redhat.com/sites/default/files/2022-05/Cloud%20Native%20Network%20Function%20Requirements%201-3.pdf Section 4.6.8
Exception Process|
#### namespace-resource-quota

Property|Description
---|---
Test Case Name|namespace-resource-quota
Test Case Label|access-control-namespace-resource-quota
Unique ID|http://test-network-function.com/testcases/access-control/namespace-resource-quota
Version|v1.0.0
Description|http://test-network-function.com/testcases/access-control/namespace-resource-quota Checks to see if CNF workload pods are running in namespaces that have resource quotas applied.
Result Type|informative
Suggested Remediation|Apply a ResourceQuota to the namespace your CNF is running in
Best Practice Reference|https://connect.redhat.com/sites/default/files/2022-05/Cloud%20Native%20Network%20Function%20Requirements%201-3.pdf Section 4.6.8
Exception Process|
#### one-process-per-container

Property|Description
---|---
Test Case Name|one-process-per-container
Test Case Label|access-control-one-process-per-container
Unique ID|http://test-network-function.com/testcases/access-control/one-process-per-container
Version|v1.0.0
Description|http://test-network-function.com/testcases/access-control/one-process-per-container check that all containers under test 		have only one process running
Result Type|informative
Suggested Remediation|Launch only one process per container
Best Practice Reference|https://connect.redhat.com/sites/default/files/2022-05/Cloud%20Native%20Network%20Function%20Requirements%201-3.pdf Section 10.8.3
Exception Process|There is no documented exception process for this.
#### pod-automount-service-account-token

Property|Description
---|---
Test Case Name|pod-automount-service-account-token
Test Case Label|access-control-pod-automount-service-account-token
Unique ID|http://test-network-function.com/testcases/access-control/pod-automount-service-account-token
Version|v1.0.0
Description|http://test-network-function.com/testcases/access-control/pod-automount-service-account-token check that all pods under test have automountServiceAccountToken set to false
Result Type|normative
Suggested Remediation|Check that pod has automountServiceAccountToken set to false or pod is attached to service account which has automountServiceAccountToken set to false
Best Practice Reference|https://connect.redhat.com/sites/default/files/2022-05/Cloud%20Native%20Network%20Function%20Requirements%201-3.pdf Section 12.7
Exception Process|Identify which Kubernetes APIs are required if you need to utilize automount service tokens.  Depending on 												which APIs are utilized, Red Hat possibly might make those APIs available to use via OpenShift.
#### pod-host-ipc

Property|Description
---|---
Test Case Name|pod-host-ipc
Test Case Label|access-control-pod-host-ipc
Unique ID|http://test-network-function.com/testcases/access-control/pod-host-ipc
Version|v1.0.0
Description|http://test-network-function.com/testcases/access-control/pod-host-ipc Verifies that the spec.HostIpc parameter is set to false
Result Type|informative
Suggested Remediation|Set the spec.HostIpc parameter to false in the pod configuration
Best Practice Reference|https://connect.redhat.com/sites/default/files/2022-05/Cloud%20Native%20Network%20Function%20Requirements%201-3.pdf Section 5.3.6
Exception Process|There is no documented exception process for this.
#### pod-host-network

Property|Description
---|---
Test Case Name|pod-host-network
Test Case Label|access-control-pod-host-network
Unique ID|http://test-network-function.com/testcases/access-control/pod-host-network
Version|v1.0.0
Description|http://test-network-function.com/testcases/access-control/pod-host-network Verifies that the spec.HostNetwork parameter is set to false
Result Type|informative
Suggested Remediation|Set the spec.HostNetwork parameter to false in the pod configuration
Best Practice Reference|https://connect.redhat.com/sites/default/files/2022-05/Cloud%20Native%20Network%20Function%20Requirements%201-3.pdf Section 5.3.6
Exception Process|There is no documented exception process for this.
#### pod-host-path

Property|Description
---|---
Test Case Name|pod-host-path
Test Case Label|access-control-pod-host-path
Unique ID|http://test-network-function.com/testcases/access-control/pod-host-path
Version|v1.0.0
Description|http://test-network-function.com/testcases/access-control/pod-host-path Verifies that the spec.HostPath parameter is not set (not present)
Result Type|informative
Suggested Remediation|Set the spec.HostNetwork parameter to false in the pod configuration
Best Practice Reference|https://connect.redhat.com/sites/default/files/2022-05/Cloud%20Native%20Network%20Function%20Requirements%201-3.pdf Section 5.3.6
Exception Process|There is no documented exception process for this.
#### pod-host-pid

Property|Description
---|---
Test Case Name|pod-host-pid
Test Case Label|access-control-pod-host-pid
Unique ID|http://test-network-function.com/testcases/access-control/pod-host-pid
Version|v1.0.0
Description|http://test-network-function.com/testcases/access-control/pod-host-pid Verifies that the spec.HostPid parameter is set to false
Result Type|informative
Suggested Remediation|Set the spec.HostPid parameter to false in the pod configuration
Best Practice Reference|https://connect.redhat.com/sites/default/files/2022-05/Cloud%20Native%20Network%20Function%20Requirements%201-3.pdf Section 5.3.6
Exception Process|There is no documented exception process for this.
#### pod-role-bindings

Property|Description
---|---
Test Case Name|pod-role-bindings
Test Case Label|access-control-pod-role-bindings
Unique ID|http://test-network-function.com/testcases/access-control/pod-role-bindings
Version|v1.0.0
Description|http://test-network-function.com/testcases/access-control/pod-role-bindings ensures that a CNF does not utilize RoleBinding(s) in a non-CNF Namespace.
Result Type|normative
Suggested Remediation|Ensure the CNF is not configured to use RoleBinding(s) in a non-CNF Namespace.
Best Practice Reference|https://connect.redhat.com/sites/default/files/2022-05/Cloud%20Native%20Network%20Function%20Requirements%201-3.pdf Section 5.3.3 and 5.3.5
Exception Process|There is no documented exception process for this.
#### pod-service-account

Property|Description
---|---
Test Case Name|pod-service-account
Test Case Label|access-control-pod-service-account
Unique ID|http://test-network-function.com/testcases/access-control/pod-service-account
Version|v1.0.0
Description|http://test-network-function.com/testcases/access-control/pod-service-account tests that each CNF Pod utilizes a valid Service Account.
Result Type|normative
Suggested Remediation|Ensure that the each CNF Pod is configured to use a valid Service Account
Best Practice Reference|https://connect.redhat.com/sites/default/files/2022-05/Cloud%20Native%20Network%20Function%20Requirements%201-3.pdf Section 5.2.3 and 5.2.7
Exception Process|There is no documented exception process for this.
#### requests-and-limits

Property|Description
---|---
Test Case Name|requests-and-limits
Test Case Label|access-control-requests-and-limits
Unique ID|http://test-network-function.com/testcases/access-control/requests-and-limits
Version|v1.0.0
Description|http://test-network-function.com/testcases/access-control/requests-and-limits Check that containers have resource requests and limits specified in their spec.
Result Type|informative
Suggested Remediation|Add requests and limits to your container spec.  See: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/#requests-and-limits
Best Practice Reference|https://TODO Section 4.6.11
Exception Process|
#### security-context-capabilities-check

Property|Description
---|---
Test Case Name|security-context-capabilities-check
Test Case Label|access-control-security-context-capabilities-check
Unique ID|http://test-network-function.com/testcases/access-control/security-context-capabilities-check
Version|v1.0.0
Description|http://test-network-function.com/testcases/access-control/security-context-capabilities-check Tests that the following capabilities are not granted: 			- NET_ADMIN 			- SYS_ADMIN  			- NET_RAW 			- IPC_LOCK 
Result Type|normative
Suggested Remediation|Remove the following capabilities from the container/pod definitions: NET_ADMIN SCC, SYS_ADMIN SCC, NET_RAW SCC, IPC_LOCK SCC
Best Practice Reference|https://connect.redhat.com/sites/default/files/2022-05/Cloud%20Native%20Network%20Function%20Requirements%201-3.pdf Section 5.2
Exception Process|Identify the pod that is needing special capabilities and document why  
#### security-context-non-root-user-check

Property|Description
---|---
Test Case Name|security-context-non-root-user-check
Test Case Label|access-control-security-context-non-root-user-check
Unique ID|http://test-network-function.com/testcases/access-control/security-context-non-root-user-check
Version|v1.0.0
Description|http://test-network-function.com/testcases/access-control/security-context-non-root-user-check Checks the security context runAsUser parameter in pods and containers to make sure it is not set to uid root(0)
Result Type|normative
Suggested Remediation|Change the pod and containers "runAsUser" uid to something other than root(0)
Best Practice Reference|https://connect.redhat.com/sites/default/files/2022-05/Cloud%20Native%20Network%20Function%20Requirements%201-3.pdf Section 5.2
Exception Process|If your application needs root user access, please document why your application cannot be ran as 											non-root and supply the reasoning for exception.
#### security-context-privilege-escalation

Property|Description
---|---
Test Case Name|security-context-privilege-escalation
Test Case Label|access-control-security-context-privilege-escalation
Unique ID|http://test-network-function.com/testcases/access-control/security-context-privilege-escalation
Version|v1.0.0
Description|http://test-network-function.com/testcases/access-control/security-context-privilege-escalation Checks if privileged escalation is enabled (AllowPrivilegeEscalation=true)
Result Type|normative
Suggested Remediation|Configure privilege escalation to false
Best Practice Reference|https://connect.redhat.com/sites/default/files/2022-05/Cloud%20Native%20Network%20Function%20Requirements%201-3.pdf Section 5.2
Exception Process|There is no documented exception process for this.
#### sys-nice-realtime-capability

Property|Description
---|---
Test Case Name|sys-nice-realtime-capability
Test Case Label|access-control-sys-nice-realtime-capability
Unique ID|http://test-network-function.com/testcases/access-control/sys-nice-realtime-capability
Version|v1.0.0
Description|http://test-network-function.com/testcases/access-control/sys-nice-realtime-capability Check that pods running on nodes with realtime kernel enabled have the SYS_NICE capability enabled in their spec.
Result Type|informative
Suggested Remediation|If pods are scheduled to realtime kernel nodes, they must add SYS_NICE capability to their spec.
Best Practice Reference|https://connect.redhat.com/sites/default/files/2022-05/Cloud%20Native%20Network%20Function%20Requirements%201-3.pdf Section 2.7.4
Exception Process|There is no documented exception process for this.
#### sys-ptrace-capability

Property|Description
---|---
Test Case Name|sys-ptrace-capability
Test Case Label|access-control-sys-ptrace-capability
Unique ID|http://test-network-function.com/testcases/access-control/sys-ptrace-capability
Version|v1.0.0
Description|http://test-network-function.com/testcases/access-control/sys-ptrace-capability check that if process namespace sharing is enabled for a Pod then the SYS_PTRACE capability is allowed
Result Type|informative
Suggested Remediation|Allow the SYS_PTRACE capability when enabling process namespace sharing for a Pod
Best Practice Reference|https://connect.redhat.com/sites/default/files/2022-05/Cloud%20Native%20Network%20Function%20Requirements%201-3.pdf Section 2.7.5
Exception Process|There is no documented exception process for this.

### affiliated-certification

#### container-is-certified

Property|Description
---|---
Test Case Name|container-is-certified
Test Case Label|affiliated-certification-container-is-certified
Unique ID|http://test-network-function.com/testcases/affiliated-certification/container-is-certified
Version|v1.0.0
Description|http://test-network-function.com/testcases/affiliated-certification/container-is-certified tests whether container images listed in the configuration file have passed the Red Hat Container Certification Program (CCP).
Result Type|normative
Suggested Remediation|Ensure that your container has passed the Red Hat Container Certification Program (CCP).
Best Practice Reference|https://connect.redhat.com/sites/default/files/2022-05/Cloud%20Native%20Network%20Function%20Requirements%201-3.pdf Section 5.3.7
Exception Process|There is no documented exception process for this.
#### helmchart-is-certified

Property|Description
---|---
Test Case Name|helmchart-is-certified
Test Case Label|affiliated-certification-helmchart-is-certified
Unique ID|http://test-network-function.com/testcases/affiliated-certification/helmchart-is-certified
Version|v1.0.0
Description|http://test-network-function.com/testcases/affiliated-certification/helmchart-is-certified tests whether helm charts listed in the cluster passed the Red Hat Helm Certification Program.
Result Type|normative
Suggested Remediation|Ensure that the helm charts under test passed the Red Hat's helm Certification Program (e.g. listed in https://charts.openshift.io/index.yaml).
Best Practice Reference|https://connect.redhat.com/sites/default/files/2022-05/Cloud%20Native%20Network%20Function%20Requirements%201-3.pdf Section 5.2.12 and 5.3.3
Exception Process|There is no documented exception process for this.
#### operator-is-certified

Property|Description
---|---
Test Case Name|operator-is-certified
Test Case Label|affiliated-certification-operator-is-certified
Unique ID|http://test-network-function.com/testcases/affiliated-certification/operator-is-certified
Version|v1.0.0
Description|http://test-network-function.com/testcases/affiliated-certification/operator-is-certified tests whether CNF Operators listed in the configuration file have passed the Red Hat Operator Certification Program (OCP).
Result Type|normative
Suggested Remediation|Ensure that your Operator has passed Red Hat's Operator Certification Program (OCP).
Best Practice Reference|https://connect.redhat.com/sites/default/files/2022-05/Cloud%20Native%20Network%20Function%20Requirements%201-3.pdf Section 5.2.12 and 5.3.3
Exception Process|There is no documented exception process for this.

### lifecycle

#### container-shutdown

Property|Description
---|---
Test Case Name|container-shutdown
Test Case Label|lifecycle-container-shutdown
Unique ID|http://test-network-function.com/testcases/lifecycle/container-shutdown
Version|v1.0.0
Description|http://test-network-function.com/testcases/lifecycle/container-shutdown Ensure that the containers lifecycle pre-stop management feature is configured.
Result Type|normative
Suggested Remediation| 	It's considered best-practices to define prestop for proper management of container lifecycle. 	The prestop can be used to gracefully stop the container and clean resources (e.g., DB connection). 	 	The prestop can be configured using : 	 1) Exec : executes the supplied command inside the container 	 2) HTTP : executes HTTP request against the specified endpoint. 	 	When defined. K8s will handle shutdown of the container using the following: 	1) K8s first execute the preStop hook inside the container. 	2) K8s will wait for a grace period. 	3) K8s will clean the remaining processes using KILL signal.		 		
Best Practice Reference|https://connect.redhat.com/sites/default/files/2022-05/Cloud%20Native%20Network%20Function%20Requirements%201-3.pdf Section 5.1.3, 12.2 and 12.5
Exception Process|Identify which pod is not conforming to the process and submit information as to why it cannot  									use a prestop shutdown specification.
#### deployment-scaling

Property|Description
---|---
Test Case Name|deployment-scaling
Test Case Label|lifecycle-deployment-scaling
Unique ID|http://test-network-function.com/testcases/lifecycle/deployment-scaling
Version|v1.0.0
Description|http://test-network-function.com/testcases/lifecycle/deployment-scaling tests that CNF deployments support scale in/out operations.  			First, The test starts getting the current replicaCount (N) of the deployment/s with the Pod Under Test. Then, it executes the  			scale-in oc command for (N-1) replicas. Lastly, it executes the scale-out oc command, restoring the original replicaCount of the deployment/s. 		    In case of deployments that are managed by HPA the test is changing the min and max value to deployment Replica - 1 during scale-in and the  			original replicaCount again for both min/max during the scale-out stage. lastly its restoring the original min/max replica of the deployment/s
Result Type|normative
Suggested Remediation|Ensure CNF deployments/replica sets can scale in/out successfully.
Best Practice Reference|https://connect.redhat.com/sites/default/files/2022-05/Cloud%20Native%20Network%20Function%20Requirements%201-3.pdf Section 5.2
Exception Process|There is no documented exception process for this.
#### image-pull-policy

Property|Description
---|---
Test Case Name|image-pull-policy
Test Case Label|lifecycle-image-pull-policy
Unique ID|http://test-network-function.com/testcases/lifecycle/image-pull-policy
Version|v1.0.0
Description|http://test-network-function.com/testcases/lifecycle/image-pull-policy Ensure that the containers under test are using IfNotPresent as Image Pull Policy..
Result Type|normative
Suggested Remediation|Ensure that the containers under test are using IfNotPresent as Image Pull Policy.
Best Practice Reference|https://connect.redhat.com/sites/default/files/2022-05/Cloud%20Native%20Network%20Function%20Requirements%201-3.pdf  Section 12.6
Exception Process|There is no documented exception process for this.
#### liveness

Property|Description
---|---
Test Case Name|liveness
Test Case Label|lifecycle-liveness
Unique ID|http://test-network-function.com/testcases/lifecycle/liveness
Version|v1.0.0
Description|http://test-network-function.com/testcases/lifecycle/liveness check that all containers under test 		have liveness probe defined
Result Type|normative
Suggested Remediation|Add liveness probe to deployed containers
Best Practice Reference|https://connect.redhat.com/sites/default/files/2022-05/Cloud%20Native%20Network%20Function%20Requirements%201-3.pdf Section 5.2.16, 12.1 and 12.5
Exception Process|There is no documented exception process for this.
#### pod-high-availability

Property|Description
---|---
Test Case Name|pod-high-availability
Test Case Label|lifecycle-pod-high-availability
Unique ID|http://test-network-function.com/testcases/lifecycle/pod-high-availability
Version|v1.0.0
Description|http://test-network-function.com/testcases/lifecycle/pod-high-availability ensures that CNF Pods specify podAntiAffinity rules and replica value is set to more than 1.
Result Type|informative
Suggested Remediation|In high availability cases, Pod podAntiAffinity rule should be specified for pod scheduling and pod replica value is set to more than 1 .
Best Practice Reference|https://connect.redhat.com/sites/default/files/2022-05/Cloud%20Native%20Network%20Function%20Requirements%201-3.pdf Section 5.2
Exception Process|There is no documented exception process for this.
#### pod-owner-type

Property|Description
---|---
Test Case Name|pod-owner-type
Test Case Label|lifecycle-pod-owner-type
Unique ID|http://test-network-function.com/testcases/lifecycle/pod-owner-type
Version|v1.0.0
Description|http://test-network-function.com/testcases/lifecycle/pod-owner-type tests that CNF Pod(s) are deployed as part of a ReplicaSet(s)/StatefulSet(s).
Result Type|normative
Suggested Remediation|Deploy the CNF using ReplicaSet/StatefulSet.
Best Practice Reference|https://connect.redhat.com/sites/default/files/2022-05/Cloud%20Native%20Network%20Function%20Requirements%201-3.pdf Section 5.3.3 and 5.3.8
Exception Process|There is no documented exception process for this.
#### pod-recreation

Property|Description
---|---
Test Case Name|pod-recreation
Test Case Label|lifecycle-pod-recreation
Unique ID|http://test-network-function.com/testcases/lifecycle/pod-recreation
Version|v1.0.0
Description|http://test-network-function.com/testcases/lifecycle/pod-recreation tests that a CNF is configured to support High Availability.   			First, this test cordons and drains a Node that hosts the CNF Pod.   			Next, the test ensures that OpenShift can re-instantiate the Pod on another Node,  			and that the actual replica count matches the desired replica count.
Result Type|normative
Suggested Remediation|Ensure that CNF Pod(s) utilize a configuration that supports High Availability.   	Additionally, ensure that there are available Nodes in the OpenShift cluster that can be utilized in the event that a host Node fails.
Best Practice Reference|https://connect.redhat.com/sites/default/files/2022-05/Cloud%20Native%20Network%20Function%20Requirements%201-3.pdf Section 5.2
Exception Process|There is no documented exception process for this.
#### pod-scheduling

Property|Description
---|---
Test Case Name|pod-scheduling
Test Case Label|lifecycle-pod-scheduling
Unique ID|http://test-network-function.com/testcases/lifecycle/pod-scheduling
Version|v1.0.0
Description|http://test-network-function.com/testcases/lifecycle/pod-scheduling ensures that CNF Pods do not specify nodeSelector or nodeAffinity.  In most cases, Pods should allow for instantiation on any underlying Node.
Result Type|informative
Suggested Remediation|In most cases, Pod's should not specify their host Nodes through nodeSelector or nodeAffinity.  However, there are 	cases in which CNFs require specialized hardware specific to a particular class of Node.  As such, this test is purely 	informative, and will not prevent a CNF from being certified. However, one should have an appropriate justification as 	to why nodeSelector and/or nodeAffinity is utilized by a CNF.
Best Practice Reference|https://connect.redhat.com/sites/default/files/2022-05/Cloud%20Native%20Network%20Function%20Requirements%201-3.pdf Section 5.2
Exception Process|There is no documented exception process for this.
#### readiness

Property|Description
---|---
Test Case Name|readiness
Test Case Label|lifecycle-readiness
Unique ID|http://test-network-function.com/testcases/lifecycle/readiness
Version|v1.0.0
Description|http://test-network-function.com/testcases/lifecycle/readiness check that all containers under test 		have readiness probe defined
Result Type|normative
Suggested Remediation|Add readiness probe to deployed containers
Best Practice Reference|https://connect.redhat.com/sites/default/files/2022-05/Cloud%20Native%20Network%20Function%20Requirements%201-3.pdf Section 5.2.16, 12.1 and 12.5
Exception Process|There is no documented exception process for this.
#### scaling

Property|Description
---|---
Test Case Name|scaling
Test Case Label|lifecycle-scaling
Unique ID|http://test-network-function.com/testcases/lifecycle/scaling
Version|v1.0.0
Description|http://test-network-function.com/testcases/lifecycle/scaling tests that CNF deployments support scale in/out operations.  			First, The test starts getting the current replicaCount (N) of the deployment/s with the Pod Under Test. Then, it executes the  			scale-in oc command for (N-1) replicas. Lastly, it executes the scale-out oc command, restoring the original replicaCount of the deployment/s.
Result Type|normative
Suggested Remediation|Ensure CNF deployments/replica sets can scale in/out successfully.
Best Practice Reference|https://connect.redhat.com/sites/default/files/2022-05/Cloud%20Native%20Network%20Function%20Requirements%201-3.pdf Section 5.2
Exception Process|There is no documented exception process for this.
#### statefulset-scaling

Property|Description
---|---
Test Case Name|statefulset-scaling
Test Case Label|lifecycle-statefulset-scaling
Unique ID|http://test-network-function.com/testcases/lifecycle/statefulset-scaling
Version|v1.0.0
Description|http://test-network-function.com/testcases/lifecycle/statefulset-scaling tests that CNF statefulsets support scale in/out operations.  			First, The test starts getting the current replicaCount (N) of the statefulset/s with the Pod Under Test. Then, it executes the  			scale-in oc command for (N-1) replicas. Lastly, it executes the scale-out oc command, restoring the original replicaCount of the statefulset/s. 			In case of statefulsets that are managed by HPA the test is changing the min and max value to statefulset Replica - 1 during scale-in and the  			original replicaCount again for both min/max during the scale-out stage. lastly its restoring the original min/max replica of the statefulset/s
Result Type|normative
Suggested Remediation|Ensure CNF statefulsets/replica sets can scale in/out successfully.
Best Practice Reference|https://connect.redhat.com/sites/default/files/2022-05/Cloud%20Native%20Network%20Function%20Requirements%201-3.pdf Section 5.2
Exception Process|There is no documented exception process for this.

### networking

#### icmpv4-connectivity

Property|Description
---|---
Test Case Name|icmpv4-connectivity
Test Case Label|networking-icmpv4-connectivity
Unique ID|http://test-network-function.com/testcases/networking/icmpv4-connectivity
Version|v1.0.0
Description|http://test-network-function.com/testcases/networking/icmpv4-connectivity checks that each CNF Container is able to communicate via ICMPv4 on the Default OpenShift network.  This test case requires the Deployment of the debug daemonset.
Result Type|normative
Suggested Remediation|Ensure that the CNF is able to communicate via the Default OpenShift network. In some rare cases, 	CNFs may require routing table changes in order to communicate over the Default network. To exclude a particular pod 	from ICMPv4 connectivity tests, add the test-network-function.com/skip_connectivity_tests label to it. The label value is not important, only its presence.
Best Practice Reference|https://connect.redhat.com/sites/default/files/2022-05/Cloud%20Native%20Network%20Function%20Requirements%201-3.pdf Section 5.2
Exception Process|There is no documented exception process for this.
#### icmpv4-connectivity-multus

Property|Description
---|---
Test Case Name|icmpv4-connectivity-multus
Test Case Label|networking-icmpv4-connectivity-multus
Unique ID|http://test-network-function.com/testcases/networking/icmpv4-connectivity-multus
Version|v1.0.0
Description|http://test-network-function.com/testcases/networking/icmpv4-connectivity-multus checks that each CNF Container is able to communicate via ICMPv4 on the Multus network(s).  This test case requires the Deployment of the debug daemonset.
Result Type|normative
Suggested Remediation|Ensure that the CNF is able to communicate via the Multus network(s). In some rare cases, 	CNFs may require routing table changes in order to communicate over the Multus network(s). To exclude a particular pod 	from ICMPv4 connectivity tests, add the test-network-function.com/skip_connectivity_tests label to it. The label value is not important, only its presence.
Best Practice Reference|https://connect.redhat.com/sites/default/files/2022-05/Cloud%20Native%20Network%20Function%20Requirements%201-3.pdf Section 5.2
Exception Process|There is no documented exception process for this.
#### icmpv6-connectivity

Property|Description
---|---
Test Case Name|icmpv6-connectivity
Test Case Label|networking-icmpv6-connectivity
Unique ID|http://test-network-function.com/testcases/networking/icmpv6-connectivity
Version|v1.0.0
Description|http://test-network-function.com/testcases/networking/icmpv6-connectivity checks that each CNF Container is able to communicate via ICMPv6 on the Default OpenShift network.  This test case requires the Deployment of the debug daemonset.
Result Type|normative
Suggested Remediation|Ensure that the CNF is able to communicate via the Default OpenShift network. In some rare cases, 	CNFs may require routing table changes in order to communicate over the Default network. To exclude a particular pod 	from ICMPv6 connectivity tests, add the test-network-function.com/skip_connectivity_tests label to it. The label value is not important, only its presence.
Best Practice Reference|https://connect.redhat.com/sites/default/files/2022-05/Cloud%20Native%20Network%20Function%20Requirements%201-3.pdf Section 5.2
Exception Process|There is no documented exception process for this.
#### icmpv6-connectivity-multus

Property|Description
---|---
Test Case Name|icmpv6-connectivity-multus
Test Case Label|networking-icmpv6-connectivity-multus
Unique ID|http://test-network-function.com/testcases/networking/icmpv6-connectivity-multus
Version|v1.0.0
Description|http://test-network-function.com/testcases/networking/icmpv6-connectivity-multus checks that each CNF Container is able to communicate via ICMPv6 on the Multus network(s).  This test case requires the Deployment of the debug daemonset.
Result Type|normative
Suggested Remediation|Ensure that the CNF is able to communicate via the Multus network(s). In some rare cases, 	CNFs may require routing table changes in order to communicate over the Multus network(s). To exclude a particular pod 	from ICMPv6 connectivity tests, add the test-network-function.com/skip_connectivity_tests label to it.The label value is not important, only its presence.
Best Practice Reference|https://connect.redhat.com/sites/default/files/2022-05/Cloud%20Native%20Network%20Function%20Requirements%201-3.pdf Section 5.2
Exception Process|There is no documented exception process for this.
#### ocp-reserved-ports-usage

Property|Description
---|---
Test Case Name|ocp-reserved-ports-usage
Test Case Label|networking-ocp-reserved-ports-usage
Unique ID|http://test-network-function.com/testcases/networking/ocp-reserved-ports-usage
Version|v1.0.0
Description|http://test-network-function.com/testcases/networking/ocp-reserved-ports-usage check that containers do not listen on ports that are reserved by Openshift
Result Type|normative
Suggested Remediation|Ensure that CNF apps do not listen on ports that are reserved by Openshift
Best Practice Reference|https://TODO Section 3.5.9
Exception Process|There is no documented exception process for this.
#### service-type

Property|Description
---|---
Test Case Name|service-type
Test Case Label|networking-service-type
Unique ID|http://test-network-function.com/testcases/networking/service-type
Version|v1.0.0
Description|http://test-network-function.com/testcases/networking/service-type tests that each CNF Service does not utilize NodePort(s).
Result Type|normative
Suggested Remediation|Ensure Services are not configured to use NodePort(s).
Best Practice Reference|https://connect.redhat.com/sites/default/files/2022-05/Cloud%20Native%20Network%20Function%20Requirements%201-3.pdf Section 5.3.1
Exception Process|There is no documented exception process for this.
#### undeclared-container-ports-usage

Property|Description
---|---
Test Case Name|undeclared-container-ports-usage
Test Case Label|networking-undeclared-container-ports-usage
Unique ID|http://test-network-function.com/testcases/networking/undeclared-container-ports-usage
Version|v1.0.0
Description|http://test-network-function.com/testcases/networking/undeclared-container-ports-usage check that containers do not listen on ports that weren't declared in their specification
Result Type|normative
Suggested Remediation|Ensure the CNF apps do not listen on undeclared containers' ports
Best Practice Reference|https://connect.redhat.com/sites/default/files/2022-05/Cloud%20Native%20Network%20Function%20Requirements%201-3.pdf Section 16.3.1.1
Exception Process|There is no documented exception process for this.

### observability

#### container-logging

Property|Description
---|---
Test Case Name|container-logging
Test Case Label|observability-container-logging
Unique ID|http://test-network-function.com/testcases/observability/container-logging
Version|v1.0.0
Description|http://test-network-function.com/testcases/observability/container-logging check that all containers under test use standard input output and standard error when logging
Result Type|informative
Suggested Remediation|Ensure containers are not redirecting stdout/stderr
Best Practice Reference|https://connect.redhat.com/sites/default/files/2022-05/Cloud%20Native%20Network%20Function%20Requirements%201-3.pdf Section 10.1
Exception Process|There is no documented exception process for this.
#### crd-status

Property|Description
---|---
Test Case Name|crd-status
Test Case Label|observability-crd-status
Unique ID|http://test-network-function.com/testcases/observability/crd-status
Version|v1.0.0
Description|http://test-network-function.com/testcases/observability/crd-status checks that all CRDs have a status subresource specification (Spec.versions[].Schema.OpenAPIV3Schema.Properties["status"]).
Result Type|informative
Suggested Remediation|Ensure that all the CRDs have a meaningful status specification (Spec.versions[].Schema.OpenAPIV3Schema.Properties["status"]).
Best Practice Reference|https://connect.redhat.com/sites/default/files/2022-05/Cloud%20Native%20Network%20Function%20Requirements%201-3.pdf Section 5.2
Exception Process|There is no documented exception process for this.
#### termination-policy

Property|Description
---|---
Test Case Name|termination-policy
Test Case Label|observability-termination-policy
Unique ID|http://test-network-function.com/testcases/observability/termination-policy
Version|v1.0.0
Description|http://test-network-function.com/testcases/observability/termination-policy check that all containers are using terminationMessagePolicy: FallbackToLogsOnError
Result Type|informative
Suggested Remediation|Ensure containers are all using FallbackToLogsOnError in terminationMessagePolicy
Best Practice Reference|https://connect.redhat.com/sites/default/files/2022-05/Cloud%20Native%20Network%20Function%20Requirements%201-3.pdf Section 12.1
Exception Process|There is no documented exception process for this.

### operator

#### install-source

Property|Description
---|---
Test Case Name|install-source
Test Case Label|operator-install-source
Unique ID|http://test-network-function.com/testcases/operator/install-source
Version|v1.0.0
Description|http://test-network-function.com/testcases/operator/install-source tests whether a CNF Operator is installed via OLM.
Result Type|normative
Suggested Remediation|Ensure that your Operator is installed via OLM.
Best Practice Reference|https://connect.redhat.com/sites/default/files/2022-05/Cloud%20Native%20Network%20Function%20Requirements%201-3.pdf Section 5.2.12 and 5.3.3
Exception Process|There is no documented exception process for this.
#### install-status-no-privileges

Property|Description
---|---
Test Case Name|install-status-no-privileges
Test Case Label|operator-install-status-no-privileges
Unique ID|http://test-network-function.com/testcases/operator/install-status-no-privileges
Version|v1.0.0
Description|http://test-network-function.com/testcases/operator/install-status-no-privileges The operator is not installed with privileged rights. Test passes if clusterPermissions is not present in the CSV manifest or is present  with no resourceNames under its rules.
Result Type|normative
Suggested Remediation|Ensure all the CNF operators have no privileges on cluster resources.
Best Practice Reference|https://connect.redhat.com/sites/default/files/2022-05/Cloud%20Native%20Network%20Function%20Requirements%201-3.pdf Section 5.2.12 and 5.3.3
Exception Process|There is no documented exception process for this.
#### install-status-succeeded

Property|Description
---|---
Test Case Name|install-status-succeeded
Test Case Label|operator-install-status-succeeded
Unique ID|http://test-network-function.com/testcases/operator/install-status-succeeded
Version|v1.0.0
Description|http://test-network-function.com/testcases/operator/install-status-succeeded Ensures that the target CNF operators report "Succeeded" as their installation status.
Result Type|normative
Suggested Remediation|Ensure all the CNF operators have been successfully installed by OLM.
Best Practice Reference|https://connect.redhat.com/sites/default/files/2022-05/Cloud%20Native%20Network%20Function%20Requirements%201-3.pdf Section 5.2.12 and 5.3.3
Exception Process|There is no documented exception process for this.

### platform-alteration

#### base-image

Property|Description
---|---
Test Case Name|base-image
Test Case Label|platform-alteration-base-image
Unique ID|http://test-network-function.com/testcases/platform-alteration/base-image
Version|v1.0.0
Description|http://test-network-function.com/testcases/platform-alteration/base-image ensures that the Container Base Image is not altered post-startup.  This test is a heuristic, and ensures that there are no changes to the following directories: 1) /var/lib/rpm 2) /var/lib/dpkg 3) /bin 4) /sbin 5) /lib 6) /lib64 7) /usr/bin 8) /usr/sbin 9) /usr/lib 10) /usr/lib64
Result Type|normative
Suggested Remediation|Ensure that Container applications do not modify the Container Base Image.  In particular, ensure that the following 	directories are not modified: 	1) /var/lib/rpm 	2) /var/lib/dpkg 	3) /bin 	4) /sbin 	5) /lib 	6) /lib64 	7) /usr/bin 	8) /usr/sbin 	9) /usr/lib 	10) /usr/lib64 	Ensure that all required binaries are built directly into the container image, and are not installed post startup.
Best Practice Reference|https://connect.redhat.com/sites/default/files/2022-05/Cloud%20Native%20Network%20Function%20Requirements%201-3.pdf Section 5.1.4
Exception Process|Images should not be changed during runtime.  There is no exception process for this.
#### boot-params

Property|Description
---|---
Test Case Name|boot-params
Test Case Label|platform-alteration-boot-params
Unique ID|http://test-network-function.com/testcases/platform-alteration/boot-params
Version|v1.0.0
Description|http://test-network-function.com/testcases/platform-alteration/boot-params tests that boot parameters are set through the MachineConfigOperator, and not set manually on the Node.
Result Type|normative
Suggested Remediation|Ensure that boot parameters are set directly through the MachineConfigOperator, or indirectly through the PerformanceAddonOperator.   	Boot parameters should not be changed directly through the Node, as OpenShift should manage the changes for you.
Best Practice Reference|https://connect.redhat.com/sites/default/files/2022-05/Cloud%20Native%20Network%20Function%20Requirements%201-3.pdf Section 5.2.13 and 5.2.14
Exception Process|There is no documented exception process for this.
#### hugepages-config

Property|Description
---|---
Test Case Name|hugepages-config
Test Case Label|platform-alteration-hugepages-config
Unique ID|http://test-network-function.com/testcases/platform-alteration/hugepages-config
Version|v1.0.0
Description|http://test-network-function.com/testcases/platform-alteration/hugepages-config checks to see that HugePage settings have been configured through MachineConfig, and not manually on the underlying Node.  This test case applies only to Nodes that are configured with the "worker" MachineConfigSet.  First, the "worker" MachineConfig is polled, and the Hugepage settings are extracted.  Next, the underlying Nodes are polled for configured HugePages through inspection of /proc/meminfo.  The results are compared, and the test passes only if they are the same.
Result Type|normative
Suggested Remediation|HugePage settings should be configured either directly through the MachineConfigOperator or indirectly using the 	PerformanceAddonOperator.  This ensures that OpenShift is aware of the special MachineConfig requirements, and can 	provision your CNF on a Node that is part of the corresponding MachineConfigSet.  Avoid making changes directly to an 	underlying Node, and let OpenShift handle the heavy lifting of configuring advanced settings.
Best Practice Reference|https://connect.redhat.com/sites/default/files/2022-05/Cloud%20Native%20Network%20Function%20Requirements%201-3.pdf Section 5.2
Exception Process|There is no documented exception process for this.
#### is-selinux-enforcing

Property|Description
---|---
Test Case Name|is-selinux-enforcing
Test Case Label|platform-alteration-is-selinux-enforcing
Unique ID|http://test-network-function.com/testcases/platform-alteration/is-selinux-enforcing
Version|v1.0.0
Description|http://test-network-function.com/testcases/platform-alteration/is-selinux-enforcing verifies that all openshift platform/cluster nodes have selinux in "Enforcing" mode.
Result Type|normative
Suggested Remediation|Configure selinux and enable enforcing mode.
Best Practice Reference|https://connect.redhat.com/sites/default/files/2022-05/Cloud%20Native%20Network%20Function%20Requirements%201-3.pdf Section 10.3 Pod Security
Exception Process|There is no documented exception process for this.
#### isredhat-release

Property|Description
---|---
Test Case Name|isredhat-release
Test Case Label|platform-alteration-isredhat-release
Unique ID|http://test-network-function.com/testcases/platform-alteration/isredhat-release
Version|v1.0.0
Description|http://test-network-function.com/testcases/platform-alteration/isredhat-release verifies if the container base image is redhat.
Result Type|normative
Suggested Remediation|Build a new container image that is based on UBI (Red Hat Universal Base Image).
Best Practice Reference|https://connect.redhat.com/sites/default/files/2022-05/Cloud%20Native%20Network%20Function%20Requirements%201-3.pdf Section 5.2
Exception Process|Document which containers are not able to meet the RHEL-based container  											requirement and if/when the base image can be updated.
#### ocp-lifecycle

Property|Description
---|---
Test Case Name|ocp-lifecycle
Test Case Label|platform-alteration-ocp-lifecycle
Unique ID|http://test-network-function.com/testcases/platform-alteration/ocp-lifecycle
Version|v1.0.0
Description|http://test-network-function.com/testcases/platform-alteration/ocp-lifecycle Tests that the running OCP version is not end of life.
Result Type|normative
Suggested Remediation|Please update your cluster to a version that is generally available.
Best Practice Reference|https://connect.redhat.com/sites/default/files/2022-05/Cloud%20Native%20Network%20Function%20Requirements%201-3.pdf Section 7.9
Exception Process|There is no documented exception process for this.
#### ocp-node-os-lifecycle

Property|Description
---|---
Test Case Name|ocp-node-os-lifecycle
Test Case Label|platform-alteration-ocp-node-os-lifecycle
Unique ID|http://test-network-function.com/testcases/platform-alteration/ocp-node-os-lifecycle
Version|v1.0.0
Description|http://test-network-function.com/testcases/platform-alteration/ocp-node-os-lifecycle Tests that the nodes running in the cluster have operating systems 			that are compatible with the deployed version of OpenShift.
Result Type|normative
Suggested Remediation|Please update your workers to a version that is supported by your version of OpenShift
Best Practice Reference|https://connect.redhat.com/sites/default/files/2022-05/Cloud%20Native%20Network%20Function%20Requirements%201-3.pdf Section 7.9
Exception Process|There is no documented exception process for this.
#### service-mesh-usage

Property|Description
---|---
Test Case Name|service-mesh-usage
Test Case Label|platform-alteration-service-mesh-usage
Unique ID|http://test-network-function.com/testcases/platform-alteration/service-mesh-usage
Version|v1.0.0
Description|http://test-network-function.com/testcases/platform-alteration/service-mesh-usage verifies whether, if available, service mesh is actually being used by the CNF pods
Result Type|normative
Suggested Remediation|Ensure all the CNF pods are using service mesh if the cluster provides it.
Best Practice Reference|https://connect.redhat.com/sites/default/files/2022-05/Cloud%20Native%20Network%20Function%20Requirements%201-3.pdf Section 5.2
Exception Process|There is no documented exception process for this.
#### sysctl-config

Property|Description
---|---
Test Case Name|sysctl-config
Test Case Label|platform-alteration-sysctl-config
Unique ID|http://test-network-function.com/testcases/platform-alteration/sysctl-config
Version|v1.0.0
Description|http://test-network-function.com/testcases/lifecycle/pod-recreation tests that no one has changed the node's sysctl configs after the node 			was created, the tests works by checking if the sysctl configs are consistent with the 			MachineConfig CR which defines how the node should be configured
Result Type|normative
Suggested Remediation|You should recreate the node or change the sysctls, recreating is recommended because there might be other unknown changes
Best Practice Reference|https://connect.redhat.com/sites/default/files/2022-05/Cloud%20Native%20Network%20Function%20Requirements%201-3.pdf Section 5.2
Exception Process|There is no documented exception process for this.
#### tainted-node-kernel

Property|Description
---|---
Test Case Name|tainted-node-kernel
Test Case Label|platform-alteration-tainted-node-kernel
Unique ID|http://test-network-function.com/testcases/platform-alteration/tainted-node-kernel
Version|v1.0.0
Description|http://test-network-function.com/testcases/platform-alteration/tainted-node-kernel ensures that the Node(s) hosting CNFs do not utilize tainted kernels. This test case is especially important to support Highly Available CNFs, since when a CNF is re-instantiated on a backup Node, that Node's kernel may not have the same hacks.'
Result Type|normative
Suggested Remediation|Test failure indicates that the underlying Node's kernel is tainted.  Ensure that you have not altered underlying 	Node(s) kernels in order to run the CNF.
Best Practice Reference|https://connect.redhat.com/sites/default/files/2022-05/Cloud%20Native%20Network%20Function%20Requirements%201-3.pdf Section 5.2.14
Exception Process|There is no documented exception process for this.

