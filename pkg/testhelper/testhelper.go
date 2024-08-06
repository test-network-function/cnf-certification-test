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

package testhelper

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/test-network-function/cnf-certification-test/pkg/provider"
)

const (
	SUCCESS = iota
	FAILURE
	ERROR
)

type ReportObject struct {
	ObjectType         string
	ObjectFieldsKeys   []string
	ObjectFieldsValues []string
}

type FailureReasonOut struct {
	CompliantObjectsOut    []*ReportObject
	NonCompliantObjectsOut []*ReportObject
}

func Equal(p, other []*ReportObject) bool {
	if len(p) != len(other) {
		return false
	}
	for i := 0; i < len(p); i++ {
		if p[i] == nil && other[i] == nil {
			continue
		}
		if p[i] == nil || other[i] == nil {
			return false
		}
		if !reflect.DeepEqual(*p[i], *other[i]) {
			return false
		}
	}
	return true
}

// FailureReasonOutTestString returns a string representation of the FailureReasonOut struct.
func FailureReasonOutTestString(p FailureReasonOut) (out string) {
	out = "testhelper.FailureReasonOut{"
	out += fmt.Sprintf("CompliantObjectsOut: %s,", ReportObjectTestStringPointer(p.CompliantObjectsOut))
	out += fmt.Sprintf("NonCompliantObjectsOut: %s,", ReportObjectTestStringPointer(p.NonCompliantObjectsOut))
	out += "}"
	return out
}

// ReportObjectTestStringPointer takes a slice of pointers to ReportObject and returns a string representation of the objects.
// The returned string is in the format "[]*testhelper.ReportObject{&{...}, &{...}, ...}".
func ReportObjectTestStringPointer(p []*ReportObject) (out string) {
	out = "[]*testhelper.ReportObject{"
	for _, p := range p {
		out += fmt.Sprintf("&%#v,", *p)
	}
	out += "}"
	return out
}

// ReportObjectTestString returns a string representation of the given slice of ReportObject.
// Each ReportObject is formatted using the %#v format specifier and appended to the output string.
// The resulting string is enclosed in square brackets and prefixed with "[]testhelper.ReportObject{".
func ReportObjectTestString(p []*ReportObject) (out string) {
	out = "[]testhelper.ReportObject{"
	for _, p := range p {
		out += fmt.Sprintf("%#v,", *p)
	}
	out += "}"
	return out
}

// Equal checks if the current FailureReasonOut is equal to the other FailureReasonOut.
// It compares the CompliantObjectsOut and NonCompliantObjectsOut fields of both structs.
// Returns true if they are equal, false otherwise.
func (p FailureReasonOut) Equal(other FailureReasonOut) bool {
	return Equal(p.CompliantObjectsOut, other.CompliantObjectsOut) &&
		Equal(p.NonCompliantObjectsOut, other.NonCompliantObjectsOut)
}

// When adding new field types, please update the following:

const (
	Namespace                       = "Namespace"
	Name                            = "Name"
	PodName                         = "Pod Name"
	ContainerName                   = "Container Name"
	ProcessID                       = "Process ID"
	ProcessCommandLine              = "Process CommandLine"
	SchedulingPolicy                = "Scheduling Policy"
	SchedulingPriority              = "Scheduling Priority"
	ReasonForNonCompliance          = "Reason For Non Compliance"
	ReasonForCompliance             = "Reason For Compliance"
	Category                        = "Category"
	RoleBindingName                 = "Role Binding Name"
	ClusterRoleName                 = "Cluster Role Reference Name"
	RoleBindingNamespace            = "Role Binding Namespace"
	ServiceAccountName              = "Service Account Name"
	ServiceMode                     = "Service Type"
	ServiceName                     = "Service Name"
	ServiceIPVersion                = "Service IP Version"
	DeploymentName                  = "Deployment Name"
	StatefulSetName                 = "StatefulSet Name"
	PodDisruptionBudgetReference    = "Pod Disruption Budget Reference"
	CustomResourceDefinitionName    = "Custom Resource Definition Name"
	CustomResourceDefinitionVersion = "Custom Resource Definition Version"
	SCCCapability                   = "SCC Capability"
	Path                            = "Path"
	Repository                      = "Repository"
	ImageName                       = "Image Name"
	Version                         = "Version"
	OpenAPIV3Schema                 = "OpenAPIV3Schema"
	CrdVersion                      = "Operator CRD Version"
	OCPVersion                      = "OCP Version"
	OCPChannel                      = "OCP Channel"
	NodeSelector                    = "Node Selector"
	PersistentVolumeName            = "Persistent Volume Name"
	PersistentVolumeClaimName       = "Persistent Volume Claim Name"
	TolerationKey                   = "Toleration Key"
	TolerationEffect                = "Toleration Effect"
	StorageClassName                = "Storage Class Name"
	StorageClassProvisioner         = "Storage Class Provisioner"
	ChangedFolders                  = "Changed Folders"
	DeletedFolders                  = "Deleted Folders"
	TaintBit                        = "Taint Bit"
	TaintBitDescription             = "Taint Bit Description"
	TaintMask                       = "Taint Mask"
	ModuleName                      = "Module Name"
	Taints                          = "Taints"
	SysctlKey                       = "Sysctl Key"
	SysctlValue                     = "Sysctl Value"
	OSImage                         = "OS Image"
	DebugPodName                    = "Debug Pod Name"

	// ICMP tests.
	NetworkName              = "Network Name"
	DestinationNamespace     = "Destination Namespace"
	DestinationPodName       = "Destination Pod Name"
	DestinationContainerName = "Destination Container Name"
	DestinationIP            = "Destination IP"
	SourceIP                 = "Source IP"

	// Rbac roles.
	RoleName     = "Role Name"
	Group        = "Group"
	ResourceName = "Resource Name"
	Verb         = "Verb"

	// Listening ports.
	PortNumber   = "Port Number"
	PortProtocol = "Port Protocol"

	// OLM.
	SubscriptionName = "Subscription Name"
	OperatorPhase    = "Operator Phase"
	OperatorName     = "Operator Name"

	// Lists.
	OperatorList = "Operator List"
)

// When adding new object types, please update the following.
const (
	UndefinedType                = "Undefined Type"
	CnfType                      = "Cnf"
	PodType                      = "Pod"
	HelmType                     = "Helm"
	OperatorType                 = "Operator"
	ContainerType                = "Container"
	ContainerImageType           = "Container Image"
	NodeType                     = "Node"
	OCPClusterType               = "OCP Cluster"
	OCPClusterVersionType        = "OCP Cluster Version"
	ContainerProcessType         = "ContainerProcess"
	ContainerCategory            = "ContainerCategory"
	ServiceType                  = "Service"
	DeploymentType               = "Deployment"
	StatefulSetType              = "StatefulSet"
	ICMPResultType               = "ICMP result"
	NetworkType                  = "Network"
	CustomResourceDefinitionType = "Custom Resource Definition"
	RoleRuleType                 = "Role Rule"
	RoleType                     = "Role"
	ListeningPortType            = "Listening Port"
	DeclaredPortType             = "Declared Port"
	ContainerPort                = "Container Port"
	HostPortType                 = "Host Port"
	HostPathType                 = "Host Path"
	HelmVersionType              = "Helm Version"
	Error                        = "Error"
	OperatorPermission           = "Operator Cluster Permission"
	TaintType                    = "Taint"
	ImageDigest                  = "Image Digest"
	ImageRepo                    = "Image Repo"
	ImageTag                     = "Image Tag"
	ImageRegistry                = "Image Registry"
	PodRoleBinding               = "Pods with RoleBindings details"
)

// SetContainerProcessValues sets the values for a container process in the report object.
// It takes the scheduling policy, scheduling priority, and command line as input parameters.
// It adds the process command line, scheduling policy, and scheduling priority fields to the report object.
// Finally, it sets the object type to ContainerProcessType.
func (obj *ReportObject) SetContainerProcessValues(aPolicy, aPriority, aCommandLine string) *ReportObject {
	obj.AddField(ProcessCommandLine, aCommandLine)
	obj.AddField(SchedulingPolicy, aPolicy)
	obj.AddField(SchedulingPriority, aPriority)
	obj.ObjectType = ContainerProcessType
	return obj
}

// NewContainerReportObject creates a new ReportObject for a container.
// It takes the namespace, pod name, container name, reason, and compliance status as parameters.
// It returns a pointer to the created ReportObject.
func NewContainerReportObject(aNamespace, aPodName, aContainerName, aReason string, isCompliant bool) (out *ReportObject) {
	out = NewReportObject(aReason, ContainerType, isCompliant)
	out.AddField(Namespace, aNamespace)
	out.AddField(PodName, aPodName)
	out.AddField(ContainerName, aContainerName)
	return out
}

// NewCertifiedContainerReportObject creates a new ReportObject for a certified container.
// It takes a ContainerImageIdentifier, aReason string, and a boolean indicating whether the container is compliant.
// It returns a pointer to the created ReportObject.
func NewCertifiedContainerReportObject(cii provider.ContainerImageIdentifier, aReason string, isCompliant bool) (out *ReportObject) {
	out = NewReportObject(aReason, ContainerImageType, isCompliant)
	out.AddField(ImageDigest, cii.Digest)
	out.AddField(ImageRepo, cii.Repository)
	out.AddField(ImageTag, cii.Tag)
	out.AddField(ImageRegistry, cii.Registry)
	return out
}

// NewNodeReportObject creates a new ReportObject for a node with the given name, reason, and compliance status.
// It returns the created ReportObject.
func NewNodeReportObject(aNodeName, aReason string, isCompliant bool) (out *ReportObject) {
	out = NewReportObject(aReason, NodeType, isCompliant)
	out.AddField(Name, aNodeName)
	return out
}

// NewClusterVersionReportObject creates a new ReportObject for a cluster version.
// It takes the version, aReason, and isCompliant as input parameters and returns the created ReportObject.
func NewClusterVersionReportObject(version, aReason string, isCompliant bool) (out *ReportObject) {
	out = NewReportObject(aReason, OCPClusterType, isCompliant)
	out.AddField(OCPClusterVersionType, version)
	return out
}

// NewTaintReportObject creates a new ReportObject with taint-related information.
// It takes in the taintBit, nodeName, aReason, and isCompliant parameters and returns a pointer to the created ReportObject.
func NewTaintReportObject(taintBit, nodeName, aReason string, isCompliant bool) (out *ReportObject) {
	out = NewReportObject(aReason, TaintType, isCompliant)
	out.AddField(NodeType, nodeName)
	out.AddField(TaintBit, taintBit)
	return out
}

// NewPodReportObject creates a new ReportObject for a pod.
// It takes the namespace, pod name, reason, and compliance status as input parameters.
// It returns a pointer to the created ReportObject.
func NewPodReportObject(aNamespace, aPodName, aReason string, isCompliant bool) (out *ReportObject) {
	out = NewReportObject(aReason, PodType, isCompliant)
	out.AddField(Namespace, aNamespace)
	out.AddField(PodName, aPodName)
	return out
}

// NewHelmChartReportObject creates a new ReportObject for a Helm chart.
// It takes the namespace, Helm chart name, reason, and compliance status as input parameters.
// It returns the created ReportObject.
func NewHelmChartReportObject(aNamespace, aHelmChartName, aReason string, isCompliant bool) (out *ReportObject) {
	out = NewReportObject(aReason, HelmType, isCompliant)
	out.AddField(Namespace, aNamespace)
	out.AddField(Name, aHelmChartName)
	return out
}

// NewOperatorReportObject creates a new ReportObject for an operator.
// It takes the namespace, operator name, reason, and compliance status as input parameters.
// It returns the created ReportObject.
func NewOperatorReportObject(aNamespace, aOperatorName, aReason string, isCompliant bool) (out *ReportObject) {
	out = NewReportObject(aReason, OperatorType, isCompliant)
	out.AddField(Namespace, aNamespace)
	out.AddField(Name, aOperatorName)
	return out
}

// NewDeploymentReportObject creates a new ReportObject for a deployment.
// It takes the namespace, deployment name, reason, and compliance status as input parameters.
// It returns a pointer to the created ReportObject.
func NewDeploymentReportObject(aNamespace, aDeploymentName, aReason string, isCompliant bool) (out *ReportObject) {
	out = NewReportObject(aReason, DeploymentType, isCompliant)
	out.AddField(Namespace, aNamespace)
	out.AddField(DeploymentName, aDeploymentName)
	return out
}

// NewStatefulSetReportObject creates a new ReportObject for a StatefulSet.
// It takes the namespace, statefulSetName, reason, and compliance status as parameters.
// It returns the created ReportObject.
func NewStatefulSetReportObject(aNamespace, aStatefulSetName, aReason string, isCompliant bool) (out *ReportObject) {
	out = NewReportObject(aReason, StatefulSetType, isCompliant)
	out.AddField(Namespace, aNamespace)
	out.AddField(StatefulSetName, aStatefulSetName)
	return out
}

// NewCrdReportObject creates a new ReportObject for a custom resource definition (CRD).
// It takes the name, version, reason, and compliance status as parameters and returns the created ReportObject.
func NewCrdReportObject(aName, aVersion, aReason string, isCompliant bool) (out *ReportObject) {
	out = NewReportObject(aReason, CustomResourceDefinitionType, isCompliant)
	out.AddField(CustomResourceDefinitionName, aName)
	out.AddField(CustomResourceDefinitionVersion, aVersion)
	return out
}

// NewReportObject creates a new ReportObject with the specified reason, type, and compliance status.
// If isCompliant is true, the reason is added as a field with the key ReasonForCompliance.
// If isCompliant is false, the reason is added as a field with the key ReasonForNonCompliance.
// Returns a pointer to the created ReportObject.
func NewReportObject(aReason, aType string, isCompliant bool) (out *ReportObject) {
	out = &ReportObject{}
	out.ObjectType = aType
	if isCompliant {
		out.AddField(ReasonForCompliance, aReason)
	} else {
		out.AddField(ReasonForNonCompliance, aReason)
	}
	return out
}

// AddField adds a key-value pair to the ReportObject.
// It appends the given key to the ObjectFieldsKeys slice and the given value to the ObjectFieldsValues slice.
// It returns the modified ReportObject.
func (obj *ReportObject) AddField(aKey, aValue string) (out *ReportObject) {
	obj.ObjectFieldsKeys = append(obj.ObjectFieldsKeys, aKey)
	obj.ObjectFieldsValues = append(obj.ObjectFieldsValues, aValue)
	return obj
}

// NewNamespacedReportObject creates a new ReportObject with the specified reason, type, compliance status, and namespace.
// It adds the namespace field to the ReportObject.
func NewNamespacedReportObject(aReason, aType string, isCompliant bool, aNamespace string) (out *ReportObject) {
	return NewReportObject(aReason, aType, isCompliant).AddField(Namespace, aNamespace)
}

// NewNamespacedNamedReportObject creates a new namespaced named report object with the given parameters.
// It returns a pointer to the created ReportObject.
// The report object contains the specified reason, type, compliance status, namespace, and name.
func NewNamespacedNamedReportObject(aReason, aType string, isCompliant bool, aNamespace, aName string) (out *ReportObject) {
	return NewReportObject(aReason, aType, isCompliant).AddField(Namespace, aNamespace).AddField(Name, aName)
}

// SetType sets the type of the ReportObject.
// It takes aType as a parameter and updates the ObjectType field of the ReportObject.
// It returns a pointer to the updated ReportObject.
func (obj *ReportObject) SetType(aType string) (out *ReportObject) {
	obj.ObjectType = aType
	return obj
}

// ResultToString converts an integer result code into a corresponding string representation.
// It takes an integer result as input and returns the corresponding string representation.
// The possible result codes are SUCCESS, FAILURE, and ERROR.
// If the input result code is not recognized, an empty string is returned.
func ResultToString(result int) (str string) {
	switch result {
	case SUCCESS:
		return "SUCCESS"
	case FAILURE:
		return "FAILURE"
	case ERROR:
		return "ERROR"
	}
	return ""
}

func GetNonOCPClusterSkipFn() func() (bool, string) {
	return func() (bool, string) {
		if !provider.IsOCPCluster() {
			return true, "non-OCP cluster detected"
		}
		return false, ""
	}
}

func GetNoServicesUnderTestSkipFn(env *provider.TestEnvironment) func() (bool, string) {
	return func() (bool, string) {
		if len(env.Services) == 0 {
			return true, "no services to check found"
		}

		return false, ""
	}
}

func GetDaemonSetFailedToSpawnSkipFn(env *provider.TestEnvironment) func() (bool, string) {
	return func() (bool, string) {
		if env.DaemonsetFailedToSpawn {
			return true, "no daemonSets to check found"
		}

		return false, ""
	}
}

func GetNoCPUPinningPodsSkipFn(env *provider.TestEnvironment) func() (bool, string) {
	return func() (bool, string) {
		if len(env.GetCPUPinningPodsWithDpdk()) == 0 {
			return true, "no CPU pinning pods to check found"
		}

		return false, ""
	}
}

func GetNoSRIOVPodsSkipFn(env *provider.TestEnvironment) func() (bool, string) {
	return func() (bool, string) {
		pods, err := env.GetPodsUsingSRIOV()
		if err != nil {
			return true, fmt.Sprintf("failed to get SRIOV pods: %v", err)
		}

		if len(pods) == 0 {
			return true, "no SRIOV pods to check found"
		}

		return false, ""
	}
}

func GetNoContainersUnderTestSkipFn(env *provider.TestEnvironment) func() (bool, string) {
	return func() (bool, string) {
		if len(env.Containers) == 0 {
			return true, "no containers to check found"
		}

		return false, ""
	}
}

func GetNoPodsUnderTestSkipFn(env *provider.TestEnvironment) func() (bool, string) {
	return func() (bool, string) {
		if len(env.Pods) == 0 {
			return true, "no pods to check found"
		}

		return false, ""
	}
}

func GetNoDeploymentsUnderTestSkipFn(env *provider.TestEnvironment) func() (bool, string) {
	return func() (bool, string) {
		if len(env.Deployments) == 0 {
			return true, "no deployments to check found"
		}

		return false, ""
	}
}

func GetNoStatefulSetsUnderTestSkipFn(env *provider.TestEnvironment) func() (bool, string) {
	return func() (bool, string) {
		if len(env.StatefulSets) == 0 {
			return true, "no statefulSets to check found"
		}

		return false, ""
	}
}

func GetNoCrdsUnderTestSkipFn(env *provider.TestEnvironment) func() (bool, string) {
	return func() (bool, string) {
		if len(env.Crds) == 0 {
			return true, "no roles to check"
		}

		return false, ""
	}
}

func GetNoNamespacesSkipFn(env *provider.TestEnvironment) func() (bool, string) {
	return func() (bool, string) {
		if len(env.Namespaces) == 0 {
			return true, "There are no namespaces to check. Please check config."
		}

		return false, ""
	}
}

func GetNoRolesSkipFn(env *provider.TestEnvironment) func() (bool, string) {
	return func() (bool, string) {
		if len(env.Roles) == 0 {
			return true, "There are no roles to check. Please check config."
		}

		return false, ""
	}
}

func GetSharedProcessNamespacePodsSkipFn(env *provider.TestEnvironment) func() (bool, string) {
	return func() (bool, string) {
		if len(env.GetShareProcessNamespacePods()) == 0 {
			return true, "Shared process namespace pods found."
		}

		return false, ""
	}
}

func GetNotIntrusiveSkipFn(env *provider.TestEnvironment) func() (bool, string) {
	return func() (bool, string) {
		if !env.IsIntrusive() {
			return true, "not intrusive test"
		}

		return false, ""
	}
}

func GetNoPersistentVolumesSkipFn(env *provider.TestEnvironment) func() (bool, string) {
	return func() (bool, string) {
		if len(env.PersistentVolumes) == 0 {
			return true, "no persistent volumes to check found"
		}

		return false, ""
	}
}

func GetNotEnoughWorkersSkipFn(env *provider.TestEnvironment, minWorkerNodes int) func() (bool, string) {
	return func() (bool, string) {
		if env.GetWorkerCount() < minWorkerNodes {
			return true, "not enough nodes to check found"
		}

		return false, ""
	}
}

func GetPodsWithoutAffinityRequiredLabelSkipFn(env *provider.TestEnvironment) func() (bool, string) {
	return func() (bool, string) {
		if len(env.GetPodsWithoutAffinityRequiredLabel()) == 0 {
			return true, "no pods with required affinity label found"
		}

		return false, ""
	}
}

func GetNoGuaranteedPodsWithExclusiveCPUsSkipFn(env *provider.TestEnvironment) func() (bool, string) {
	return func() (bool, string) {
		if len(env.GetGuaranteedPodsWithExclusiveCPUs()) == 0 {
			return true, "no pods with exclusive CPUs found"
		}

		return false, ""
	}
}

func GetNoAffinityRequiredPodsSkipFn(env *provider.TestEnvironment) func() (bool, string) {
	return func() (bool, string) {
		if len(env.GetAffinityRequiredPods()) == 0 {
			return true, "no pods with required affinity found"
		}

		return false, ""
	}
}

func GetNoStorageClassesSkipFn(env *provider.TestEnvironment) func() (bool, string) {
	return func() (bool, string) {
		if len(env.StorageClassList) == 0 {
			return true, "no storage classes found"
		}
		return false, ""
	}
}

func GetNoPersistentVolumeClaimsSkipFn(env *provider.TestEnvironment) func() (bool, string) {
	return func() (bool, string) {
		if len(env.PersistentVolumeClaims) == 0 {
			return true, "no persistent volume claims found"
		}
		return false, ""
	}
}

func GetNoBareMetalNodesSkipFn(env *provider.TestEnvironment) func() (bool, string) {
	return func() (bool, string) {
		if len(env.GetBaremetalNodes()) == 0 {
			return true, "no baremetal nodes found"
		}
		return false, ""
	}
}

func GetNoIstioSkipFn(env *provider.TestEnvironment) func() (bool, string) {
	return func() (bool, string) {
		if !env.IstioServiceMeshFound {
			return true, "no istio service mesh found"
		}
		return false, ""
	}
}

func GetNoHugepagesPodsSkipFn(env *provider.TestEnvironment) func() (bool, string) {
	return func() (bool, string) {
		if len(env.GetHugepagesPods()) == 0 {
			return true, "no pods requesting hugepages found"
		}
		return false, ""
	}
}

func GetNoOperatorsSkipFn(env *provider.TestEnvironment) func() (bool, string) {
	return func() (bool, string) {
		if len(env.Operators) == 0 {
			return true, "no operators found"
		}
		return false, ""
	}
}

func GetNoOperatorCrdsSkipFn(env *provider.TestEnvironment) func() (bool, string) {
	return func() (bool, string) {
		if len(env.Crds) == 0 {
			return true, "no operator crds found"
		}
		return false, ""
	}
}

func ResultObjectsToString(compliantObject, nonCompliantObject []*ReportObject) (string, error) {
	reason := FailureReasonOut{
		CompliantObjectsOut:    compliantObject,
		NonCompliantObjectsOut: nonCompliantObject,
	}

	bytes, err := json.Marshal(reason)
	if err != nil {
		return "", fmt.Errorf("could not marshall FailureReasonOut object: %v", err)
	}

	return string(bytes), nil
}

var AbortTrigger string
