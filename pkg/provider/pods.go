// Copyright (C) 2022-2024 Red Hat, Inc.
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

package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/test-network-function/cnf-certification-test/internal/clientsholder"
	"github.com/test-network-function/cnf-certification-test/internal/log"
	"github.com/test-network-function/cnf-certification-test/pkg/podhelper"

	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	HugePages2Mi          = "hugepages-2Mi"
	HugePages1Gi          = "hugepages-1Gi"
	hugePages             = "hugepages"
	replicationController = "ReplicationController"
	deploymentConfig      = "DeploymentConfig"
)

type Pod struct {
	*corev1.Pod
	Containers              []*Container
	MultusNetworkInterfaces map[string]CniNetworkInterface
	MultusPCIs              []string
	SkipNetTests            bool
	SkipMultusNetTests      bool
}

func NewPod(aPod *corev1.Pod) (out Pod) {
	var err error
	out.Pod = aPod
	out.MultusNetworkInterfaces = make(map[string]CniNetworkInterface)
	out.MultusNetworkInterfaces, err = GetPodIPsPerNet(aPod.GetAnnotations()[CniNetworksStatusKey])
	if err != nil {
		log.Error("Could not get IPs for Pod %q (namespace %q), err: %v", aPod.Name, aPod.Namespace, err)
	}

	out.MultusPCIs, err = GetPciPerPod(aPod.GetAnnotations()[CniNetworksStatusKey])
	if err != nil {
		log.Error("Could not get PCIs for Pod %q (namespace %q), err: %v", aPod.Name, aPod.Namespace, err)
	}

	if _, ok := aPod.GetLabels()[skipConnectivityTestsLabel]; ok {
		out.SkipNetTests = true
	}
	if _, ok := aPod.GetLabels()[skipMultusConnectivityTestsLabel]; ok {
		out.SkipMultusNetTests = true
	}
	out.Containers = append(out.Containers, getPodContainers(aPod, false)...)
	return out
}

func ConvertArrayPods(pods []*corev1.Pod) (out []*Pod) {
	for i := range pods {
		aPodWrapper := NewPod(pods[i])
		out = append(out, &aPodWrapper)
	}
	return out
}

func (p *Pod) IsPodGuaranteed() bool {
	return AreResourcesIdentical(p)
}

func (p *Pod) IsPodGuaranteedWithExclusiveCPUs() bool {
	return AreCPUResourcesWholeUnits(p) && AreResourcesIdentical(p)
}

func (p *Pod) IsCPUIsolationCompliant() bool {
	isCPUIsolated := true

	if !LoadBalancingDisabled(p) {
		log.Debug("Pod %q has been found to not have annotations set correctly for CPU isolation.", p)
		isCPUIsolated = false
	}

	if !p.IsRuntimeClassNameSpecified() {
		log.Debug("Pod %q has been found to not have runtimeClassName specified.", p)
		isCPUIsolated = false
	}

	return isCPUIsolated
}

func (p *Pod) String() string {
	return fmt.Sprintf("pod: %s ns: %s",
		p.Name,
		p.Namespace,
	)
}

func (p *Pod) AffinityRequired() bool {
	if val, ok := p.Labels[AffinityRequiredKey]; ok {
		result, err := strconv.ParseBool(val)
		if err != nil {
			log.Warn("failure to parse bool %v", val)
			return false
		}
		return result
	}
	return false
}

// returns true if at least one container in the pod has a resource name containing "hugepage", return false otherwise.
func (p *Pod) HasHugepages() bool {
	for _, cut := range p.Containers {
		for name := range cut.Resources.Requests {
			if strings.Contains(name.String(), hugePages) {
				return true
			}
		}
		for _, name := range cut.Resources.Limits {
			if strings.Contains(name.String(), hugePages) {
				return true
			}
		}
	}
	return false
}

func (p *Pod) CheckResourceHugePagesSize(size string) bool {
	for _, cut := range p.Containers {
		// Resources must be specified
		if len(cut.Resources.Requests) == 0 || len(cut.Resources.Limits) == 0 {
			continue
		}
		for name := range cut.Resources.Requests {
			if strings.Contains(name.String(), hugePages) && name.String() != size {
				return false
			}
		}
		for name := range cut.Resources.Limits {
			if strings.Contains(name.String(), hugePages) && name.String() != size {
				return false
			}
		}
	}
	return true
}

func (p *Pod) IsAffinityCompliant() (bool, error) {
	if p.Spec.Affinity == nil {
		return false, fmt.Errorf("%s has been found with an AffinityRequired flag but is missing corresponding affinity rules", p.String())
	}
	if p.Spec.Affinity.PodAntiAffinity != nil {
		return false, fmt.Errorf("%s has been found with an AffinityRequired flag but has anti-affinity rules", p.String())
	}
	if p.Spec.Affinity.PodAffinity == nil && p.Spec.Affinity.NodeAffinity == nil {
		return false, fmt.Errorf("%s has been found with an AffinityRequired flag but is missing corresponding pod/node affinity rules", p.String())
	}
	return true, nil
}

func (p *Pod) IsShareProcessNamespace() bool {
	return p.Spec.ShareProcessNamespace != nil && *p.Spec.ShareProcessNamespace
}

func (p *Pod) ContainsIstioProxy() bool {
	for _, container := range p.Containers {
		if container.Name == "istio-proxy" {
			return true
		}
	}
	return false
}

func (p *Pod) CreatedByDeploymentConfig() (bool, error) {
	oc := clientsholder.GetClientsHolder()
	for _, podOwner := range p.ObjectMeta.GetOwnerReferences() {
		if podOwner.Kind == replicationController {
			replicationControllers, err := oc.K8sClient.CoreV1().ReplicationControllers(p.Namespace).Get(context.TODO(), podOwner.Name, metav1.GetOptions{})
			if err != nil {
				return false, err
			}
			for _, rcOwner := range replicationControllers.GetOwnerReferences() {
				if rcOwner.Name == podOwner.Name && rcOwner.Kind == deploymentConfig {
					return true, err
				}
			}
		}
	}
	return false, nil
}

func (p *Pod) HasNodeSelector() bool {
	// Checks whether or not the pod has a nodeSelector or a NodeName supplied
	if p.Spec.NodeSelector == nil || len(p.Spec.NodeSelector) == 0 {
		return false
	}
	return true
}

func (p *Pod) IsRuntimeClassNameSpecified() bool {
	return p.Spec.RuntimeClassName != nil
}

// Helper function to parse CNCF's networks annotation, retrieving
// the names only. It's a custom and simplified version of:
// https://github.com/k8snetworkplumbingwg/multus-cni/blob/e692127d19623c8bdfc4d391224ea542658b584c/pkg/k8sclient/k8sclient.go#L185
//
// The cncf netwoks annotation has two different formats:
//
//	  a) list of network names: k8s.v1.cni.cncf.io/networks: <network>[,<network>,...]
//	  b) json array of network objects:
//	    k8s.v1.cni.cncf.io/networks: |-
//			[
//				{
//				"name": "<network>",
//				"namespace": "<namespace>",
//				"default-route": ["<default-route>"]
//				}
//			]
func getCNCFNetworksNamesFromPodAnnotation(networksAnnotation string) []string {
	// Each CNCF network has many more fields, but here we only need to unmarshal the name.
	// See https://github.com/k8snetworkplumbingwg/multus-cni/blob/e692127d19623c8bdfc4d391224ea542658b584c/pkg/types/types.go#L127
	type CNCFNetwork struct {
		Name string `json:"name"`
	}

	networkObjects := []CNCFNetwork{}
	networkNames := []string{}

	// Let's start trying to unmarshal a json array of objects.
	// We won't care about bad-formatted/invalid annotation value. If that's the case,
	// the pod wouldn't have been deployed or wouldn't be in running state.
	if err := json.Unmarshal([]byte(networksAnnotation), &networkObjects); err == nil {
		for _, network := range networkObjects {
			networkNames = append(networkNames, network.Name)
		}
		return networkNames
	}

	// If the previous unmarshalling didn't work, let's try with parsing the comma separated names list.
	networks := strings.TrimSpace(networksAnnotation)

	// First, avoid empty strings (unlikely).
	if networks == "" {
		return []string{}
	}

	for _, networkName := range strings.Split(networks, ",") {
		networkNames = append(networkNames, strings.TrimSpace(networkName))
	}
	return networkNames
}

// isNetworkAttachmentDefinitionConfigTypeSRIOV is a helper function to check whether a CNI
// config string has any config for sriov plugin.
// CNI config has two modes: single CNI plugin, or multi-plugins:
// Single CNI plugin config sample:
//
//	{
//		"cniVersion": "0.4.0",
//		"name": "sriov-network",
//		"type": "sriov",
//		...
//	}
//
// Multi-plugin CNI config sample:
//
//	{
//		"cniVersion": "0.4.0",
//		"name": "sriov-network",
//		"plugins": [
//			{
//				"type": "sriov",
//				"device": "eth1",
//				...
//			},
//			{
//				"type": "firewall"
//				...
//			}
//		]
func isNetworkAttachmentDefinitionConfigTypeSRIOV(nadConfig string) (bool, error) {
	const (
		typeSriov = "sriov"
	)

	type CNIConfig struct {
		CniVersion string  `json:"cniVersion"`
		Name       string  `json:"name"`
		Type       *string `json:"type,omitempty"`
		Plugins    *[]struct {
			Type string `json:"type"`
		} `json:"plugins,omitempty"`
	}

	cniConfig := CNIConfig{}
	if err := json.Unmarshal([]byte(nadConfig), &cniConfig); err != nil {
		return false, fmt.Errorf("failed to unmarshal cni config %s: %v", nadConfig, err)
	}

	// If type is found, it's a single plugin CNI config.
	if cniConfig.Type != nil {
		log.Debug("Single plugin config type found: %+v, type=%s", cniConfig, *cniConfig.Type)
		return *cniConfig.Type == typeSriov, nil
	}

	if cniConfig.Plugins == nil {
		return false, fmt.Errorf("invalid multi-plugins cni config: %s", nadConfig)
	}

	log.Debug("CNI plugins: %+v", *cniConfig.Plugins)
	for i := range *cniConfig.Plugins {
		plugin := (*cniConfig.Plugins)[i]
		if plugin.Type == typeSriov {
			return true, nil
		}
	}

	// No sriov plugin type found.
	return false, nil
}

// IsUsingSRIOV returns true if any of the pod's interfaces is a sriov one.
// First, it retrieves the list of networks names from the CNFC annotation and then
// checks the config of the corresponding network-attachment definition (NAD).
func (p *Pod) IsUsingSRIOV() (bool, error) {
	const (
		cncfNetworksAnnotation = "k8s.v1.cni.cncf.io/networks"
	)

	cncfNetworks, exist := p.Annotations[cncfNetworksAnnotation]
	if !exist {
		return false, nil
	}

	// Get all CNCF network names
	cncfNetworkNames := getCNCFNetworksNamesFromPodAnnotation(cncfNetworks)

	// For each CNCF network, get its network attachment definition and check
	// whether its config's type is "sriov"
	oc := clientsholder.GetClientsHolder()

	for _, networkName := range cncfNetworkNames {
		log.Debug("%s: Reviewing network-attachment definition %q", p, networkName)
		nad, err := oc.CNCFNetworkingClient.NetworkAttachmentDefinitions(p.Namespace).Get(context.TODO(), networkName, metav1.GetOptions{})
		if err != nil {
			return false, fmt.Errorf("failed to get NetworkAttachment %s: %v", networkName, err)
		}

		isSRIOV, err := isNetworkAttachmentDefinitionConfigTypeSRIOV(nad.Spec.Config)
		if err != nil {
			return false, fmt.Errorf("failed to know if network-attachment %s is sriov: %v", networkName, err)
		}

		log.Debug("%s: NAD config: %s", p, nad.Spec.Config)
		if isSRIOV {
			return true, nil
		}
	}

	return false, nil
}

//nolint:gocritic
func (p *Pod) IsUsingClusterRoleBinding(clusterRoleBindings []rbacv1.ClusterRoleBinding, logger *log.Logger) (bool, string, error) {
	// This function accepts a list of clusterRoleBindings and checks to see if the pod's service account is
	// tied to any of them.  If it is, then it returns true, otherwise it returns false.
	logger.Info("Pod %q is using service account %q", p, p.Pod.Spec.ServiceAccountName)

	// Loop through the service accounts in the namespace, looking for a match between the pod serviceAccountName and
	// the service account name.  If there is a match, check to make sure that the SA is not a 'subject' of the cluster
	// role bindings.
	for crbIndex := range clusterRoleBindings {
		for _, subject := range clusterRoleBindings[crbIndex].Subjects {
			if subject.Kind == rbacv1.ServiceAccountKind && subject.Name == p.Pod.Spec.ServiceAccountName && subject.Namespace == p.Pod.Namespace {
				logger.Error("Pod %q has service account %q that is tied to cluster role binding %q", p.Pod.Name, p.Pod.Spec.ServiceAccountName, clusterRoleBindings[crbIndex].Name)
				return true, clusterRoleBindings[crbIndex].RoleRef.Name, nil
			}
		}
	}

	return false, "", nil
}

func (p *Pod) IsRunAsUserID(uid int64) bool {
	if p.Pod.Spec.SecurityContext == nil || p.Pod.Spec.SecurityContext.RunAsUser == nil {
		return false
	}
	return *p.Pod.Spec.SecurityContext.RunAsUser == uid
}

func (p *Pod) IsRunAsNonRoot() bool {
	// Check pod-level security context.
	if p.Pod.Spec.SecurityContext != nil && p.Pod.Spec.SecurityContext.RunAsNonRoot != nil {
		return *p.Pod.Spec.SecurityContext.RunAsNonRoot
	}

	// If neither container-level nor pod-level security context is set, fail.
	return false
}

// Get the list of top owners of pods.
func (p *Pod) GetTopOwner() (topOwners map[string]podhelper.TopOwner, err error) {
	return podhelper.GetPodTopOwner(p.Namespace, p.OwnerReferences)
}
