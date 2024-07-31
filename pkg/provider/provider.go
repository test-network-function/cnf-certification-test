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

package provider

import (
	"context"
	"regexp"
	"time"

	"fmt"
	"strings"

	"encoding/json"

	mcv1 "github.com/openshift/api/machineconfiguration/v1"
	olmv1Alpha "github.com/operator-framework/api/pkg/operators/v1alpha1"
	"github.com/redhat-best-practices-for-k8s/certsuite/internal/clientsholder"
	"github.com/redhat-best-practices-for-k8s/certsuite/internal/log"
	"github.com/redhat-best-practices-for-k8s/certsuite/pkg/autodiscover"
	"github.com/redhat-best-practices-for-k8s/certsuite/pkg/configuration"
	k8sPrivilegedDs "github.com/redhat-best-practices-for-k8s/privileged-daemonset"
	plibRuntime "github.com/redhat-openshift-ecosystem/openshift-preflight/certification"
	"helm.sh/helm/v3/pkg/release"
	scalingv1 "k8s.io/api/autoscaling/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	policyv1 "k8s.io/api/policy/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	storagev1 "k8s.io/api/storage/v1"
	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
)

// CentOS Stream CoreOS starts being used instead of rhcos from OCP 4.13 latest.
const (
	AffinityRequiredKey              = "AffinityRequired"
	containerName                    = "container-00"
	DaemonSetName                    = "tnf-debug"
	debugPodsTimeout                 = 5 * time.Minute
	CniNetworksStatusKey             = "k8s.v1.cni.cncf.io/network-status"
	skipConnectivityTestsLabel       = "redhat-best-practices-for-k8s.com/skip_connectivity_tests"
	skipMultusConnectivityTestsLabel = "redhat-best-practices-for-k8s.com/skip_multus_connectivity_tests"
	rhcosName                        = "Red Hat Enterprise Linux CoreOS"
	cscosName                        = "CentOS Stream CoreOS"
	rhelName                         = "Red Hat Enterprise Linux"
)

// Node's roles labels. Node is role R if it has **any** of the labels of each list.
// Master's role label "master" is deprecated since k8s 1.20.
var (
	WorkerLabels = []string{"node-role.kubernetes.io/worker"}
	MasterLabels = []string{"node-role.kubernetes.io/master", "node-role.kubernetes.io/control-plane"}
)

type TestEnvironment struct { // rename this with testTarget
	Namespaces     []string `json:"testNamespaces"`
	AbnormalEvents []*Event

	// Pod Groupings
	Pods            []*Pod                 `json:"testPods"`
	DebugPods       map[string]*corev1.Pod // map from nodename to debugPod
	AllPods         []*Pod                 `json:"AllPods"`
	CSVToPodListMap map[string][]*Pod      `json:"CSVToPodListMap"`

	// Deployment Groupings
	Deployments []*Deployment `json:"testDeployments"`
	// StatefulSet Groupings
	StatefulSets []*StatefulSet `json:"testStatefulSets"`

	// Note: Containers is a filtered list of objects based on a block list of disallowed container names.
	Containers             []*Container `json:"testContainers"`
	Operators              []*Operator  `json:"testOperators"`
	AllOperators           []*Operator  `json:"AllOperators"`
	AllOperatorsSummary    []string     `json:"AllOperatorsSummary"`
	PersistentVolumes      []corev1.PersistentVolume
	PersistentVolumeClaims []corev1.PersistentVolumeClaim
	ClusterRoleBindings    []rbacv1.ClusterRoleBinding
	RoleBindings           []rbacv1.RoleBinding
	Roles                  []rbacv1.Role

	Config  configuration.TestConfiguration
	params  configuration.TestParameters
	Crds    []*apiextv1.CustomResourceDefinition `json:"testCrds"`
	AllCrds []*apiextv1.CustomResourceDefinition

	HorizontalScaler       []*scalingv1.HorizontalPodAutoscaler `json:"testHorizontalScaler"`
	Services               []*corev1.Service                    `json:"testServices"`
	Nodes                  map[string]Node                      `json:"-"`
	K8sVersion             string                               `json:"-"`
	OpenshiftVersion       string                               `json:"-"`
	OCPStatus              string                               `json:"-"`
	HelmChartReleases      []*release.Release                   `json:"testHelmChartReleases"`
	ResourceQuotas         []corev1.ResourceQuota
	PodDisruptionBudgets   []policyv1.PodDisruptionBudget
	NetworkPolicies        []networkingv1.NetworkPolicy
	AllInstallPlans        []*olmv1Alpha.InstallPlan   `json:"-"`
	AllCatalogSources      []*olmv1Alpha.CatalogSource `json:"-"`
	IstioServiceMeshFound  bool
	ValidProtocolNames     []string
	DaemonsetFailedToSpawn bool
	ScaleCrUnderTest       []ScaleObject
	StorageClassList       []storagev1.StorageClass
	ExecutedBy             string
	PartnerName            string
	CollectorAppPassword   string
	CollectorAppEndpoint   string
	SkipPreflight          bool
}

type MachineConfig struct {
	*mcv1.MachineConfig
	Config struct {
		Systemd struct {
			Units []struct {
				Contents string `json:"contents"`
				Name     string `json:"name"`
			} `json:"units"`
		} `json:"systemd"`
	} `json:"config"`
}

type CniNetworkInterface struct {
	Name       string                 `json:"name"`
	Interface  string                 `json:"interface"`
	IPs        []string               `json:"ips"`
	Default    bool                   `json:"default"`
	DNS        map[string]interface{} `json:"dns"`
	DeviceInfo deviceInfo             `json:"device-info"`
}

type ScaleObject struct {
	Scale               CrScale
	GroupResourceSchema schema.GroupResource
}

type deviceInfo struct {
	Type    string `json:"type"`
	Version string `json:"version"`
	PCI     pci    `json:"pci"`
}

type pci struct {
	PciAddress string `json:"pci-address"`
}
type PreflightTest struct {
	Name        string
	Description string
	Remediation string
	Error       error
}

type PreflightResultsDB struct {
	Passed []PreflightTest
	Failed []PreflightTest
	Errors []PreflightTest
}

type APIRequestCount struct {
	Metadata struct {
		Name string `json:"name"`
	} `json:"metadata"`
	Status struct {
		RemovedInRelease string `json:"removedInRelease"`
		Last24h          []struct {
			ByNode []struct {
				ByUser []struct {
					UserName  string `json:"username"`
					UserAgent string `json:"userAgent"`
				} `json:"byUser"`
			} `json:"byNode"`
		} `json:"last24h"`
	} `json:"status"`
}

type ToBeRemovedWorkloadAPI struct {
	APIName          string
	RemovedInRelease string
	UserInfo         map[string]struct{} // Using a map to store unique tuples of username/useragent
}

var (
	env    = TestEnvironment{}
	loaded = false
)

func deployDaemonSet(namespace string) error {
	k8sPrivilegedDs.SetDaemonSetClient(clientsholder.GetClientsHolder().K8sClient)

	dsImage := env.params.TnfImageRepo + "/" + env.params.TnfDebugImage
	if k8sPrivilegedDs.IsDaemonSetReady(DaemonSetName, namespace, dsImage) {
		return nil
	}

	matchLabels := make(map[string]string)
	matchLabels["name"] = DaemonSetName
	matchLabels["redhat-best-practices-for-k8s.com/app"] = DaemonSetName
	_, err := k8sPrivilegedDs.CreateDaemonSet(DaemonSetName, namespace, containerName, dsImage, matchLabels, debugPodsTimeout,
		configuration.GetTestParameters().DaemonsetCPUReq,
		configuration.GetTestParameters().DaemonsetCPULim,
		configuration.GetTestParameters().DaemonsetMemReq,
		configuration.GetTestParameters().DaemonsetMemLim,
	)
	if err != nil {
		return fmt.Errorf("could not deploy tnf daemonset, err=%v", err)
	}
	err = k8sPrivilegedDs.WaitDaemonsetReady(namespace, DaemonSetName, debugPodsTimeout)
	if err != nil {
		return fmt.Errorf("timed out waiting for tnf daemonset, err=%v", err)
	}

	return nil
}

func buildTestEnvironment() { //nolint:funlen
	start := time.Now()
	env = TestEnvironment{}

	env.params = *configuration.GetTestParameters()
	config, err := configuration.LoadConfiguration(env.params.ConfigFile)
	if err != nil {
		log.Fatal("Cannot load configuration file: %v", err)
	}
	log.Debug("CERTSUITE configuration: %+v", config)

	// Wait for the debug pods to be ready before the autodiscovery starts.
	if err := deployDaemonSet(config.DebugDaemonSetNamespace); err != nil {
		log.Error("The TNF daemonset could not be deployed, err: %v", err)
		// Because of this failure, we are only able to run a certain amount of tests that do not rely
		// on the existence of the daemonset debug pods.
		env.DaemonsetFailedToSpawn = true
	}

	data := autodiscover.DoAutoDiscover(&config)
	// OpenshiftVersion needs to be set asap, as other helper functions will use it here.
	env.OpenshiftVersion = data.OpenshiftVersion
	env.Config = config
	env.Crds = data.Crds
	env.AllInstallPlans = data.AllInstallPlans
	env.AllCatalogSources = data.AllCatalogSources
	env.AllOperators = createOperators(data.AllCsvs, data.AllSubscriptions, data.AllInstallPlans, data.AllCatalogSources, false, false)
	env.AllOperatorsSummary = getSummaryAllOperators(env.AllOperators)
	env.AllCrds = data.AllCrds
	env.Namespaces = data.Namespaces
	env.Nodes = createNodes(data.Nodes.Items)
	env.IstioServiceMeshFound = data.IstioServiceMeshFound
	env.ValidProtocolNames = append(env.ValidProtocolNames, data.ValidProtocolNames...)
	for i := range data.AbnormalEvents {
		aEvent := NewEvent(&data.AbnormalEvents[i])
		env.AbnormalEvents = append(env.AbnormalEvents, &aEvent)
	}
	pods := data.Pods
	for i := 0; i < len(pods); i++ {
		aNewPod := NewPod(&pods[i])
		env.Pods = append(env.Pods, &aNewPod)
		// Note: 'getPodContainers' is returning a filtered list of Container objects.
		env.Containers = append(env.Containers, getPodContainers(&pods[i], true)...)
	}
	pods = data.AllPods
	for i := 0; i < len(pods); i++ {
		aNewPod := NewPod(&pods[i])
		env.AllPods = append(env.AllPods, &aNewPod)
	}
	env.DebugPods = make(map[string]*corev1.Pod)
	for i := 0; i < len(data.DebugPods); i++ {
		nodeName := data.DebugPods[i].Spec.NodeName
		env.DebugPods[nodeName] = &data.DebugPods[i]
	}

	env.CSVToPodListMap = make(map[string][]*Pod)
	for k, podList := range data.CSVToPodListMap {
		var pods []*Pod
		for i := 0; i < len(podList); i++ {
			aNewPod := NewPod(podList[i])
			pods = append(pods, &aNewPod)
		}
		env.CSVToPodListMap[k] = pods
	}

	env.OCPStatus = data.OCPStatus
	env.K8sVersion = data.K8sVersion
	env.ResourceQuotas = data.ResourceQuotaItems
	env.PodDisruptionBudgets = data.PodDisruptionBudgets
	env.PersistentVolumes = data.PersistentVolumes
	env.PersistentVolumeClaims = data.PersistentVolumeClaims
	env.ClusterRoleBindings = data.ClusterRoleBindings
	env.RoleBindings = data.RoleBindings
	env.Roles = data.Roles
	env.Services = data.Services
	env.NetworkPolicies = data.NetworkPolicies
	for _, nsHelmChartReleases := range data.HelmChartReleases {
		for _, helmChartRelease := range nsHelmChartReleases {
			if !isSkipHelmChart(helmChartRelease.Name, config.SkipHelmChartList) {
				env.HelmChartReleases = append(env.HelmChartReleases, helmChartRelease)
			}
		}
	}
	for i := range data.Deployments {
		aNewDeployment := &Deployment{
			&data.Deployments[i],
		}
		env.Deployments = append(env.Deployments, aNewDeployment)
	}
	for i := range data.StatefulSet {
		aNewStatefulSet := &StatefulSet{
			&data.StatefulSet[i],
		}
		env.StatefulSets = append(env.StatefulSets, aNewStatefulSet)
	}

	env.ScaleCrUnderTest = updateCrUnderTest(data.ScaleCrUnderTest)
	env.HorizontalScaler = data.Hpas
	env.StorageClassList = data.StorageClasses

	env.ExecutedBy = data.ExecutedBy
	env.PartnerName = data.PartnerName
	env.CollectorAppPassword = data.CollectorAppPassword
	env.CollectorAppEndpoint = data.CollectorAppEndpoint

	operators := createOperators(data.Csvs, data.Subscriptions, data.AllInstallPlans, data.AllCatalogSources, false, true)
	env.Operators = operators
	log.Info("Operators found: %d", len(env.Operators))
	for _, pod := range env.Pods {
		isCreatedByDeploymentConfig, err := pod.CreatedByDeploymentConfig()
		if err != nil {
			log.Warn("Pod %q failed to get parent resource: %v", pod, err)
			continue
		}

		if isCreatedByDeploymentConfig {
			log.Warn("Pod %q has been deployed using a DeploymentConfig, please use Deployment or StatefulSet instead.", pod.String())
		}
	}
	log.Info("Completed the test environment build process in %.2f seconds", time.Since(start).Seconds())
}

func updateCrUnderTest(scaleCrUnderTest []autodiscover.ScaleObject) []ScaleObject {
	var scaleCrUndeTestTemp []ScaleObject
	for i := range scaleCrUnderTest {
		aNewScaleCrUnderTest := ScaleObject{Scale: CrScale{scaleCrUnderTest[i].Scale},
			GroupResourceSchema: scaleCrUnderTest[i].GroupResourceSchema}
		scaleCrUndeTestTemp = append(scaleCrUndeTestTemp, aNewScaleCrUnderTest)
	}
	return scaleCrUndeTestTemp
}

func getPodContainers(aPod *corev1.Pod, useIgnoreList bool) (containerList []*Container) {
	for j := 0; j < len(aPod.Spec.Containers); j++ {
		cut := &(aPod.Spec.Containers[j])

		var cutStatus corev1.ContainerStatus
		// get Status for current container
		for index := range aPod.Status.ContainerStatuses {
			if aPod.Status.ContainerStatuses[index].Name == cut.Name {
				cutStatus = aPod.Status.ContainerStatuses[index]
				break
			}
		}
		aRuntime, uid := GetRuntimeUID(&cutStatus)
		container := Container{Podname: aPod.Name, Namespace: aPod.Namespace,
			NodeName: aPod.Spec.NodeName, Container: cut, Status: cutStatus, Runtime: aRuntime, UID: uid,
			ContainerImageIdentifier: buildContainerImageSource(aPod.Spec.Containers[j].Image, cutStatus.ImageID)}

		// Warn if readiness probe did not succeeded yet.
		if !cutStatus.Ready {
			log.Warn("Container %q is not ready yet.", &container)
		}

		// Warn if container state is not running.
		if state := &cutStatus.State; state.Running == nil {
			reason := ""
			switch {
			case state.Waiting != nil:
				reason = "waiting - " + state.Waiting.Reason
			case state.Terminated != nil:
				reason = "terminated - " + state.Terminated.Reason
			default:
				// When no state was explicitly set, it's assumed to be in "waiting state".
				reason = "waiting state reason unknown"
			}

			log.Warn("Container %q is not running (reason: %s, restarts %d): some test cases might fail.",
				&container, reason, cutStatus.RestartCount)
		}

		// Build slices of containers based on whether or not we are "ignoring" them or not.
		if useIgnoreList && container.HasIgnoredContainerName() {
			continue
		}
		containerList = append(containerList, &container)
	}
	return containerList
}

func isSkipHelmChart(helmName string, skipHelmChartList []configuration.SkipHelmChartList) bool {
	if len(skipHelmChartList) == 0 {
		return false
	}
	for _, helm := range skipHelmChartList {
		if helmName == helm.Name {
			log.Info("Helm chart with name %s was skipped", helmName)
			return true
		}
	}
	return false
}

func GetTestEnvironment() TestEnvironment {
	if !loaded {
		buildTestEnvironment()
		loaded = true
	}
	return env
}

func IsOCPCluster() bool {
	return env.OpenshiftVersion != autodiscover.NonOpenshiftClusterVersion
}

func buildContainerImageSource(urlImage, urlImageID string) (source ContainerImageIdentifier) {
	const regexImageWithTag = `^([^/]*)/*([^@]*):(.*)`
	const regexImageDigest = `^([^/]*)/(.*)@(.*:.*)`

	// get image repository, Name and tag if present
	re := regexp.MustCompile(regexImageWithTag)
	match := re.FindStringSubmatch(urlImage)

	if match != nil {
		if match[2] != "" {
			source.Registry = match[1]
			source.Repository = match[2]
			source.Tag = match[3]
		} else {
			source.Repository = match[1]
			source.Tag = match[3]
		}
	}

	// get image Digest based on imageID only
	re = regexp.MustCompile(regexImageDigest)
	match = re.FindStringSubmatch(urlImageID)

	if match != nil {
		source.Digest = match[3]
	}

	log.Debug("Parsed image, repo: %s, name:%s, tag: %s, digest: %s",
		source.Registry,
		source.Repository,
		source.Tag,
		source.Digest)

	return source
}

func GetRuntimeUID(cs *corev1.ContainerStatus) (runtime, uid string) {
	split := strings.Split(cs.ContainerID, "://")
	if len(split) > 0 {
		uid = split[len(split)-1]
		runtime = split[0]
	}
	return runtime, uid
}

// GetPodIPsPerNet gets the IPs of a pod.
// CNI annotation "k8s.v1.cni.cncf.io/networks-status".
// Returns (ips, error).
func GetPodIPsPerNet(annotation string) (ips map[string]CniNetworkInterface, err error) {
	// This is a map indexed with the network name (network attachment) and
	// listing all the IPs created in this subnet and belonging to the pod namespace
	// The list of ips pr net is parsed from the content of the "k8s.v1.cni.cncf.io/networks-status" annotation.
	ips = make(map[string]CniNetworkInterface)

	var cniInfo []CniNetworkInterface
	err = json.Unmarshal([]byte(annotation), &cniInfo)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal network-status annotation, err: %v", err)
	}
	// If this is the default interface, skip it as it is tested separately
	// Otherwise add all non default interfaces
	for _, cniInterface := range cniInfo {
		if !cniInterface.Default {
			ips[cniInterface.Name] = cniInterface
		}
	}
	return ips, nil
}

func GetPciPerPod(annotation string) (pciAddr []string, err error) {
	var cniInfo []CniNetworkInterface
	err = json.Unmarshal([]byte(annotation), &cniInfo)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal network-status annotation, err: %v", err)
	}
	for _, cniInterface := range cniInfo {
		if cniInterface.DeviceInfo.PCI.PciAddress != "" {
			pciAddr = append(pciAddr, cniInterface.DeviceInfo.PCI.PciAddress)
		}
	}
	return pciAddr, nil
}

func (env *TestEnvironment) SetNeedsRefresh() {
	loaded = false
}

func (env *TestEnvironment) IsIntrusive() bool {
	return !env.params.NonIntrusiveOnly
}

func (env *TestEnvironment) IsPreflightInsecureAllowed() bool {
	return env.params.AllowPreflightInsecure
}

func (env *TestEnvironment) GetDockerConfigFile() string {
	return env.params.PfltDockerconfig
}

func (env *TestEnvironment) GetOfflineDBPath() string {
	return env.params.OfflineDB
}

func (env *TestEnvironment) GetWorkerCount() int {
	workerCount := 0
	for _, e := range env.Nodes {
		if e.IsWorkerNode() {
			workerCount++
		}
	}
	return workerCount
}

func (env *TestEnvironment) GetMasterCount() int {
	masterCount := 0
	for _, e := range env.Nodes {
		if e.IsControlPlaneNode() {
			masterCount++
		}
	}
	return masterCount
}

func (env *TestEnvironment) IsSNO() bool {
	return len(env.Nodes) == 1
}

func getMachineConfig(mcName string, machineConfigs map[string]MachineConfig) (MachineConfig, error) {
	client := clientsholder.GetClientsHolder()

	// Check whether we had already downloaded and parsed that machineConfig resource.
	if mc, exists := machineConfigs[mcName]; exists {
		return mc, nil
	}

	nodeMc, err := client.MachineCfg.MachineconfigurationV1().MachineConfigs().Get(context.TODO(), mcName, metav1.GetOptions{})
	if err != nil {
		return MachineConfig{}, err
	}

	mc := MachineConfig{
		MachineConfig: nodeMc,
	}

	err = json.Unmarshal(nodeMc.Spec.Config.Raw, &mc.Config)
	if err != nil {
		return MachineConfig{}, fmt.Errorf("failed to unmarshal mc's Config field, err: %v", err)
	}

	return mc, nil
}

func createNodes(nodes []corev1.Node) map[string]Node {
	wrapperNodes := map[string]Node{}

	// machineConfigs is a helper map to avoid download & process the same mc twice.
	machineConfigs := map[string]MachineConfig{}
	for i := range nodes {
		node := &nodes[i]

		if !IsOCPCluster() {
			// Avoid getting Mc info for non ocp clusters.
			wrapperNodes[node.Name] = Node{Data: node}
			log.Warn("Non-OCP cluster detected. MachineConfig retrieval for node %q skipped.", node.Name)
			continue
		}

		// Get Node's machineConfig name
		mcName, exists := node.Annotations["machineconfiguration.openshift.io/currentConfig"]
		if !exists {
			log.Error("Failed to get machineConfig name for node %q", node.Name)
			continue
		}
		log.Info("Node %q - mc name %q", node.Name, mcName)
		mc, err := getMachineConfig(mcName, machineConfigs)
		if err != nil {
			log.Error("Failed to get machineConfig %q, err: %v", mcName, err)
			continue
		}

		wrapperNodes[node.Name] = Node{
			Data: node,
			Mc:   mc,
		}
	}

	return wrapperNodes
}
func (env *TestEnvironment) GetBaremetalNodes() []Node {
	var baremetalNodes []Node
	for _, node := range env.Nodes {
		if strings.HasPrefix(node.Data.Spec.ProviderID, "baremetalhost://") {
			baremetalNodes = append(baremetalNodes, node)
		}
	}
	return baremetalNodes
}

func GetPreflightResultsDB(results *plibRuntime.Results) PreflightResultsDB {
	resultsDB := PreflightResultsDB{}
	for _, res := range results.Passed {
		test := PreflightTest{Name: res.Name(), Description: res.Metadata().Description, Remediation: res.Help().Suggestion}
		resultsDB.Passed = append(resultsDB.Passed, test)
	}
	for _, res := range results.Failed {
		test := PreflightTest{Name: res.Name(), Description: res.Metadata().Description, Remediation: res.Help().Suggestion}
		resultsDB.Failed = append(resultsDB.Failed, test)
	}
	for _, res := range results.Errors {
		test := PreflightTest{Name: res.Name(), Description: res.Metadata().Description, Remediation: res.Help().Suggestion, Error: res.Error()}
		resultsDB.Errors = append(resultsDB.Errors, test)
	}

	return resultsDB
}

// Retrieve APIRequestCounts for the last 24 hours and unmarshal them
func RetrieveAPIRequestCounts(ocpClient rest.Interface) ([]APIRequestCount, error) {
	apiRequestCountURL := "/apis/apiserver.openshift.io/v1/apirequestcounts"
	req := ocpClient.Get().RequestURI(apiRequestCountURL)
	resp, err := req.DoRaw(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve APIRequestCount objects: %w", err)
	}

	var result struct {
		Items []APIRequestCount `json:"items"`
	}

	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	return result.Items, nil
}
