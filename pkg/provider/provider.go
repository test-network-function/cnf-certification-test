// Copyright (C) 2020-2021 Red Hat, Inc.
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
	"os"
	"time"

	"errors"
	"fmt"
	"strings"

	"encoding/json"

	mcv1 "github.com/openshift/machine-config-operator/pkg/apis/machineconfiguration.openshift.io/v1"
	olmv1Alpha "github.com/operator-framework/api/pkg/operators/v1alpha1"
	"github.com/sirupsen/logrus"
	"github.com/test-network-function/cnf-certification-test/cnf-certification-test/platform/operatingsystem"
	"github.com/test-network-function/cnf-certification-test/internal/clientsholder"
	"github.com/test-network-function/cnf-certification-test/pkg/autodiscover"
	"github.com/test-network-function/cnf-certification-test/pkg/configuration"
	"github.com/test-network-function/cnf-certification-test/pkg/stringhelper"
	"helm.sh/helm/v3/pkg/release"
	appsv1 "k8s.io/api/apps/v1"
	scalingv1 "k8s.io/api/autoscaling/v1"
	corev1 "k8s.io/api/core/v1"
	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	appv1client "k8s.io/client-go/kubernetes/typed/apps/v1"
)

const (
	DaemonSetNamespace               = "default"
	DaemonSetName                    = "debug"
	debugPodsTimeout                 = 5 * time.Minute
	debugPodsRetryInterval           = 5 * time.Second
	CniNetworksStatusKey             = "k8s.v1.cni.cncf.io/networks-status"
	skipConnectivityTestsLabel       = "test-network-function.com/skip_connectivity_tests"
	skipMultusConnectivityTestsLabel = "test-network-function.com/skip_multus_connectivity_tests"
	rhcosName                        = "Red Hat Enterprise Linux CoreOS"
	rhelName                         = "Red Hat Enterprise Linux"
)

// Node's roles labels. Node is role R if it has **any** of the labels of each list.
// Master's role label "master" is deprecated since k8s 1.20.
var (
	WorkerLabels      = []string{"node-role.kubernetes.io/worker"}
	MasterLabels      = []string{"node-role.kubernetes.io/master", "node-role.kubernetes.io/control-plane"}
	rhcosRelativePath = "%s/platform/operatingsystem/files/rhcos_version_map"
)

type Pod struct {
	Data               *corev1.Pod
	Containers         []*Container
	MultusIPs          map[string][]string
	SkipNetTests       bool
	SkipMultusNetTests bool
}

func NewPod(aPod *corev1.Pod) (out Pod) {
	var err error
	out.Data = aPod
	out.MultusIPs = make(map[string][]string)
	out.MultusIPs, err = GetPodIPsPerNet(aPod.GetAnnotations()[CniNetworksStatusKey])
	if err != nil {
		logrus.Errorf("Could not decode networks-status annotation")
	}

	if _, ok := aPod.GetLabels()[skipConnectivityTestsLabel]; ok {
		out.SkipNetTests = true
	}
	if _, ok := aPod.GetLabels()[skipMultusConnectivityTestsLabel]; ok {
		out.SkipMultusNetTests = true
	}
	out.Containers = append(out.Containers, getPodContainers(aPod)...)
	return out
}

func ConvertArrayPods(pods []*corev1.Pod) (out []*Pod) {
	for i := range pods {
		aPodWrapper := NewPod(pods[i])
		out = append(out, &aPodWrapper)
	}
	return out
}

type TestEnvironment struct { // rename this with testTarget
	Namespaces        []string               `json:"testNamespaces"`
	Pods              []*Pod                 `json:"testPods"`
	Containers        []*Container           `json:"testContainers"`
	Operators         []Operator             `json:"testOperators"`
	DebugPods         map[string]*corev1.Pod // map from nodename to debugPod
	Config            configuration.TestConfiguration
	variables         configuration.TestParameters
	Crds              []*apiextv1.CustomResourceDefinition          `json:"testCrds"`
	Deployments       []*appsv1.Deployment                          `json:"testDeployments"`
	StatetfulSets     []*appsv1.StatefulSet                         `json:"testStatetfulSets"`
	HorizontalScaler  map[string]*scalingv1.HorizontalPodAutoscaler `json:"testHorizontalScaler"`
	Nodes             map[string]Node                               `json:"-"`
	K8sVersion        string                                        `json:"-"`
	OpenshiftVersion  string                                        `json:"-"`
	OCPStatus         string                                        `json:"-"`
	HelmChartReleases []*release.Release                            `json:"testHelmChartReleases"`
	LimitRanges       []corev1.LimitRange
	ResourceQuotas    []corev1.ResourceQuota
	IstioServiceMesh  bool
}

type CsvInstallPlan struct {
	// Operator's installPlan name
	Name string `yaml:"name" json:"name"`
	// BundleImage is the URL referencing the bundle image
	BundleImage string `yaml:"bundleImage" json:"bundleImage"`
	// IndexImage is the URL referencing the index image
	IndexImage string `yaml:"indexImage" json:"indexImage"`
}

type Operator struct {
	Name             string                            `yaml:"name" json:"name"`
	Namespace        string                            `yaml:"namespace" json:"namespace"`
	Csv              *olmv1Alpha.ClusterServiceVersion `yaml:"csv" json:"csv"`
	SubscriptionName string                            `yaml:"subscriptionName" json:"subscriptionName"`
	InstallPlans     []CsvInstallPlan                  `yaml:"installPlans,omitempty" json:"installPlans,omitempty"`
	Package          string                            `yaml:"package" json:"package"`
	Org              string                            `yaml:"org" json:"org"`
	Version          string                            `yaml:"version" json:"version"`
	Channel          string                            `yaml:"channel" json:"channel"`
}
type Container struct {
	Data                     *corev1.Container
	Status                   corev1.ContainerStatus
	Namespace                string
	Podname                  string
	NodeName                 string
	Runtime                  string
	UID                      string
	ContainerImageIdentifier configuration.ContainerImageIdentifier
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

type Node struct {
	Data *corev1.Node
	Mc   MachineConfig `json:"-"`
}

func (node Node) MarshalJSON() ([]byte, error) {
	return json.Marshal(&node.Data)
}

func (node *Node) IsWorkerNode() bool {
	for nodeLabel := range node.Data.Labels {
		if stringhelper.StringInSlice(WorkerLabels, nodeLabel, true) {
			return true
		}
	}
	return false
}

func (node *Node) IsMasterNode() bool {
	for nodeLabel := range node.Data.Labels {
		if stringhelper.StringInSlice(MasterLabels, nodeLabel, true) {
			return true
		}
	}
	return false
}

func (node *Node) IsRHCOS() bool {
	return strings.Contains(strings.TrimSpace(node.Data.Status.NodeInfo.OSImage), rhcosName)
}

func (node *Node) IsRHEL() bool {
	return strings.Contains(strings.TrimSpace(node.Data.Status.NodeInfo.OSImage), rhelName)
}

func (node *Node) IsRTKernel() bool {
	// More information: https://www.redhat.com/sysadmin/real-time-kernel
	return strings.Contains(strings.TrimSpace(node.Data.Status.NodeInfo.KernelVersion), "rt")
}

func (node *Node) GetRHCOSVersion() (string, error) {
	// Check if the node is running CoreOS or not
	if !node.IsRHCOS() {
		return "", fmt.Errorf("invalid OS type: %s", node.Data.Status.NodeInfo.OSImage)
	}

	path, err := os.Getwd()
	if err != nil {
		return "", err
	}

	filePath := fmt.Sprintf(rhcosRelativePath, path)

	// Red Hat Enterprise Linux CoreOS 410.84.202205031645-0 (Ootpa) --> 410.84.202205031645-0
	splitStr := strings.Split(node.Data.Status.NodeInfo.OSImage, rhcosName)
	longVersionSplit := strings.Split(strings.TrimSpace(splitStr[1]), " ")

	// Get the short version string from the long version string
	shortVersion, err := operatingsystem.GetShortVersionFromLong(longVersionSplit[0], filePath)
	if err != nil {
		return "", err
	}

	return shortVersion, nil
}

func (node *Node) GetRHELVersion() (string, error) {
	// Check if the node is running RHEL or not
	if !node.IsRHEL() {
		return "", fmt.Errorf("invalid OS type: %s", node.Data.Status.NodeInfo.OSImage)
	}

	// Red Hat Enterprise Linux 8.5 (Ootpa) --> 8.5
	splitStr := strings.Split(node.Data.Status.NodeInfo.OSImage, rhelName)
	longVersionSplit := strings.Split(strings.TrimSpace(splitStr[1]), " ")

	return longVersionSplit[0], nil
}

type cniNetworkInterface struct {
	Name      string                 `json:"name"`
	Interface string                 `json:"interface"`
	IPs       []string               `json:"ips"`
	Default   bool                   `json:"default"`
	DNS       map[string]interface{} `json:"dns"`
}

var (
	env    = TestEnvironment{}
	loaded = false
)

func GetContainer() *Container {
	return &Container{}
}

func GetUpdatedDeployment(ac appv1client.AppsV1Interface, namespace, podName string) (*appsv1.Deployment, error) {
	return autodiscover.FindDeploymentByNameByNamespace(ac, namespace, podName)
}
func GetUpdatedStatefulset(ac appv1client.AppsV1Interface, namespace, podName string) (*appsv1.StatefulSet, error) {
	return autodiscover.FindStatefulsetByNameByNamespace(ac, namespace, podName)
}

func buildTestEnvironment() { //nolint:funlen
	// Wait for the debug pods to be ready before the autodiscovery starts.
	err := WaitDebugPodsReady()
	if err != nil {
		logrus.Errorf("Debug daemonset failure: %s", err)
	}
	env = TestEnvironment{}

	data := autodiscover.DoAutoDiscover()
	env.Config = data.TestData
	env.Crds = data.Crds
	env.Namespaces = data.Namespaces
	env.variables = data.Env
	env.Nodes = createNodes(data.Nodes.Items)
	env.IstioServiceMesh = data.Istio
	pods := data.Pods

	for i := 0; i < len(pods); i++ {
		aNewPod := NewPod(&pods[i])
		env.Pods = append(env.Pods, &aNewPod)
		env.Containers = append(env.Containers, getPodContainers(&pods[i])...)
	}

	env.DebugPods = make(map[string]*corev1.Pod)
	for i := 0; i < len(data.DebugPods); i++ {
		nodeName := data.DebugPods[i].Spec.NodeName
		env.DebugPods[nodeName] = &data.DebugPods[i]
	}

	env.OpenshiftVersion = data.OpenshiftVersion
	env.OCPStatus = data.OCPStatus
	env.K8sVersion = data.K8sVersion
	env.LimitRanges = data.LimitRangeItems
	env.ResourceQuotas = data.ResourceQuotaItems
	for _, nsHelmChartReleases := range data.HelmChartReleases {
		for _, helmChartRelease := range nsHelmChartReleases {
			if !isSkipHelmChart(helmChartRelease.Name, data.TestData.SkipHelmChartList) {
				env.HelmChartReleases = append(env.HelmChartReleases, helmChartRelease)
			}
		}
	}
	for i := range data.Deployments {
		env.Deployments = append(env.Deployments, &data.Deployments[i])
	}
	for i := range data.StatefulSet {
		env.StatetfulSets = append(env.StatetfulSets, &data.StatefulSet[i])
	}
	env.HorizontalScaler = data.Hpas

	operators, err := createOperators(data.Csvs, data.Subscriptions)
	if err != nil {
		logrus.Errorf("Failed to get cluster operators: %s", err)
	}
	env.Operators = operators
	logrus.Infof("Operators found: %d", len(env.Operators))
}

func getPodContainers(aPod *corev1.Pod) (containerList []*Container) {
	for j := 0; j < len(aPod.Spec.Containers); j++ {
		cut := &(aPod.Spec.Containers[j])
		var state corev1.ContainerStatus
		if len(aPod.Status.ContainerStatuses) > 0 {
			state = aPod.Status.ContainerStatuses[j]
		} else {
			logrus.Errorf("%s is not ready, skipping status collection", aPod.String())
		}
		aRuntime, uid := GetRuntimeUID(&state)
		container := Container{Podname: aPod.Name, Namespace: aPod.Namespace,
			NodeName: aPod.Spec.NodeName, Data: cut, Status: state, Runtime: aRuntime, UID: uid,
			ContainerImageIdentifier: buildContainerImageSource(aPod.Spec.Containers[j].Image)}
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
			logrus.Infof("Helm chart with name %s was skipped", helmName)
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
	return !env.variables.NonOcpCluster
}

func IsinstalledCsv(csv *olmv1Alpha.ClusterServiceVersion, subscriptions []olmv1Alpha.Subscription) (bool, olmv1Alpha.Subscription) {
	var returnsub olmv1Alpha.Subscription
	for i := range subscriptions {
		if subscriptions[i].Status.InstalledCSV == csv.Name {
			returnsub = subscriptions[i]
			return true, returnsub
		}
	}
	return false, returnsub
}
func WaitDebugPodsReady() error {
	oc := clientsholder.GetClientsHolder()
	nodes, err := oc.K8sClient.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to get node list, err:%s", err)
	}

	nodesCount := int32(len(nodes.Items))
	isReady := false
	for start := time.Now(); !isReady && time.Since(start) < debugPodsTimeout; {
		daemonSet, err := oc.K8sClient.AppsV1().DaemonSets(DaemonSetNamespace).Get(context.TODO(), DaemonSetName, metav1.GetOptions{})
		if err != nil {
			return fmt.Errorf("failed to get daemonset, err: %s", err)
		}

		logrus.Infof("Waiting for (%d) debug pods to be ready: %+v", nodesCount, daemonSet.Status)

		if daemonSet.Status.DesiredNumberScheduled != nodesCount {
			logrus.Warnf("daemonset DesiredNumberScheduled (%d) not equal to number of nodes:%d, please instantiate debug pods on all nodes",
				daemonSet.Status.DesiredNumberScheduled, nodesCount)
		} else if isDaemonSetReady(&daemonSet.Status) {
			isReady = true
			break
		}

		time.Sleep(debugPodsRetryInterval)
	}

	if !isReady {
		return errors.New("daemonset debug pods not ready")
	}

	logrus.Infof("All the debug pods are ready.")
	return nil
}

func isDaemonSetReady(status *appsv1.DaemonSetStatus) bool {
	//nolint:gocritic
	return status.DesiredNumberScheduled == status.CurrentNumberScheduled &&
		status.DesiredNumberScheduled == status.NumberAvailable &&
		status.DesiredNumberScheduled == status.NumberReady &&
		status.NumberMisscheduled == 0
}

func (c *Container) GetUID() (string, error) {
	split := strings.Split(c.Status.ContainerID, "://")
	uid := ""
	if len(split) > 0 {
		uid = split[len(split)-1]
	}
	if uid == "" {
		logrus.Debugln(fmt.Sprintf("could not find uid of %s/%s/%s\n", c.Namespace, c.Podname, c.Data.Name))
		return "", errors.New("cannot determine container UID")
	}
	logrus.Debugln(fmt.Sprintf("uid of %s/%s/%s=%s\n", c.Namespace, c.Podname, c.Data.Name, uid))
	return uid, nil
}

func buildContainerImageSource(url string) configuration.ContainerImageIdentifier {
	source := configuration.ContainerImageIdentifier{}
	urlSegments := strings.Split(url, "/")
	n := len(urlSegments)
	if n > 1 {
		source.Repository = urlSegments[n-2]
	}
	colonIndex := strings.Index(urlSegments[n-1], ":")
	atIndex := strings.Index(urlSegments[n-1], "@")
	if atIndex == -1 {
		if colonIndex == -1 {
			source.Name = urlSegments[n-1]
		} else {
			source.Name = urlSegments[n-1][:colonIndex]
			source.Tag = urlSegments[n-1][colonIndex+1:]
		}
	} else {
		source.Name = urlSegments[n-1][:atIndex]
		source.Digest = urlSegments[n-1][atIndex+1:]
	}
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

func (c *Container) StringLong() string {
	return fmt.Sprintf("node: %s ns: %s podName: %s containerName: %s containerUID: %s containerRuntime: %s",
		c.NodeName,
		c.Namespace,
		c.Podname,
		c.Data.Name,
		c.Status.ContainerID,
		c.Runtime,
	)
}
func (c *Container) String() string {
	return fmt.Sprintf("container: %s pod: %s ns: %s",
		c.Data.Name,
		c.Podname,
		c.Namespace,
	)
}

func (p *Pod) String() string {
	return fmt.Sprintf("pod: %s ns: %s",
		p.Data.Name,
		p.Data.Namespace,
	)
}

func DeploymentToString(d *appsv1.Deployment) string {
	return fmt.Sprintf("deployment: %s ns: %s",
		d.Name,
		d.Namespace,
	)
}

func StatefulsetToString(s *appsv1.StatefulSet) string {
	return fmt.Sprintf("statefulset: %s ns: %s",
		s.Name,
		s.Namespace,
	)
}

func CsvToString(csv *olmv1Alpha.ClusterServiceVersion) string {
	return fmt.Sprintf("operator csv: %s ns: %s",
		csv.Name,
		csv.Namespace,
	)
}

func (op *Operator) String() string {
	return fmt.Sprintf("csv: %s ns:%s subscription:%s", op.Name, op.Namespace, op.SubscriptionName)
}

// GetPodIPsPerNet gets the IPs of a pod.
// CNI annotation "k8s.v1.cni.cncf.io/networks-status".
// Returns (ips, error).
func GetPodIPsPerNet(annotation string) (ips map[string][]string, err error) {
	// This is a map indexed with the network name (network attachment) and
	// listing all the IPs created in this subnet and belonging to the pod namespace
	// The list of ips pr net is parsed from the content of the "k8s.v1.cni.cncf.io/networks-status" annotation.
	ips = make(map[string][]string)

	var cniInfo []cniNetworkInterface
	err = json.Unmarshal([]byte(annotation), &cniInfo)
	if err != nil {
		return nil, errors.New("could not unmarshal network-status annotation")
	}
	// If this is the default interface, skip it as it is tested separately
	// Otherwise add all non default interfaces
	for _, cniInterface := range cniInfo {
		if !cniInterface.Default {
			ips[cniInterface.Name] = cniInterface.IPs
		}
	}
	return ips, nil
}

func (env *TestEnvironment) SetNeedsRefresh() {
	loaded = false
}

func (env *TestEnvironment) IsIntrusive() bool {
	return !env.variables.NonIntrusiveOnly
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
		if e.IsMasterNode() {
			masterCount++
		}
	}
	return masterCount
}

// getInstallPlansInNamespace is a helper function to get the installPlans in a namespace. The
// map installPlans is used to store them in order to avoid repeating http requests for a namespace
// whose installPlans were already obtained.
func getInstallPlansInNamespace(namespace string, clusterInstallPlans map[string][]olmv1Alpha.InstallPlan) ([]olmv1Alpha.InstallPlan, error) {
	// Check if installplans were stored before.
	nsInstallPlans, exist := clusterInstallPlans[namespace]
	if exist {
		return nsInstallPlans, nil
	}

	clients := clientsholder.GetClientsHolder()
	installPlanList, err := clients.OlmClient.OperatorsV1alpha1().InstallPlans(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("unable get installplans in namespace %s, err: %s", namespace, err)
	}

	nsInstallPlans = installPlanList.Items
	clusterInstallPlans[namespace] = nsInstallPlans

	return nsInstallPlans, nil
}

// getCsvInstallPlans is a helper function that returns the installPlans for a CSV in a namespace.
// The map clusterInstallPlans is used to store previously retrieved installPlans, in order to save
// http requests.
func getCsvInstallPlans(namespace, csv string, clusterInstallPlans map[string][]olmv1Alpha.InstallPlan) ([]*olmv1Alpha.InstallPlan, error) {
	nsInstallPlans, err := getInstallPlansInNamespace(namespace, clusterInstallPlans)
	if err != nil {
		return nil, err
	}

	installPlans := []*olmv1Alpha.InstallPlan{}
	for i := range nsInstallPlans {
		nsInstallPlan := &nsInstallPlans[i]
		for _, csvName := range nsInstallPlan.Spec.ClusterServiceVersionNames {
			if csv != csvName {
				continue
			}

			if nsInstallPlan.Status.BundleLookups == nil {
				logrus.Warnf("InstallPlan %s for csv %s (ns %s) does not have bundle lookups. It will be skipped.", nsInstallPlan.Name, csv, namespace)
				continue
			}

			installPlans = append(installPlans, nsInstallPlan)
		}
	}

	if len(installPlans) == 0 {
		return nil, fmt.Errorf("no installplans found for csv %s (ns %s)", csv, namespace)
	}

	return installPlans, nil
}

func getCatalogSourceImageIndexFromInstallPlan(installPlan *olmv1Alpha.InstallPlan) (string, error) {
	// ToDo/Technical debt: what to do if installPlan has more than one BundleLookups entries.
	catalogSourceName := installPlan.Status.BundleLookups[0].CatalogSourceRef.Name
	catalogSourceNamespace := installPlan.Status.BundleLookups[0].CatalogSourceRef.Namespace

	clients := clientsholder.GetClientsHolder()
	catalogSource, err := clients.OlmClient.OperatorsV1alpha1().CatalogSources(catalogSourceNamespace).Get(context.TODO(), catalogSourceName, metav1.GetOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to get catalogsource: %s", err)
	}

	return catalogSource.Spec.Image, nil
}

//nolint:funlen
func createOperators(csvs []olmv1Alpha.ClusterServiceVersion, subscriptions []olmv1Alpha.Subscription) ([]Operator, error) {
	installPlans := map[string][]olmv1Alpha.InstallPlan{} // Helper: maps a namespace name to all its installplans.
	operators := []Operator{}
	for i := range csvs {
		csv := &csvs[i]
		op := Operator{Name: csv.Name, Namespace: csv.Namespace, Csv: csv}

		packageAndVersion := strings.SplitN(csv.Name, ".", 2) //nolint:gomnd // ok
		op.Version = packageAndVersion[1]

		for s := range subscriptions {
			subscription := &subscriptions[s]
			if subscription.Status.InstalledCSV != csv.Name || subscription.Namespace != csv.Namespace {
				continue
			}

			op.SubscriptionName = subscription.Name
			op.Package = subscription.Spec.Package
			op.Org = subscription.Spec.CatalogSource
			op.Channel = subscription.Spec.Channel
			break
		}

		csvInstallPlans, err := getCsvInstallPlans(csv.Namespace, csv.Name, installPlans)
		if err != nil {
			return nil, fmt.Errorf("failed to get installPlans for csv %s (ns %s), err: %s", csv.Name, csv.Namespace, err)
		}

		for _, installPlan := range csvInstallPlans {
			indexImage, err := getCatalogSourceImageIndexFromInstallPlan(installPlan)
			if err != nil {
				return nil, fmt.Errorf("failed to get installPlan image index for csv %s (ns %s) installPlan %s, err: %s",
					csv.Name, csv.Namespace, installPlan.Name, err)
			}

			op.InstallPlans = append(op.InstallPlans, CsvInstallPlan{
				Name:        installPlan.Name,
				BundleImage: installPlan.Status.BundleLookups[0].Path,
				IndexImage:  indexImage,
			})
		}
		operators = append(operators, op)
	}

	return operators, nil
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
			logrus.Warnf("Non OCP cluster detected. MachineConfig retrieval for node %s skipped.", node.Name)
			continue
		}

		// Get Node's machineConfig name
		mcName, exists := node.Annotations["machineconfiguration.openshift.io/currentConfig"]
		if !exists {
			logrus.Errorf("Failed to get machineConfig name for node %s", node.Name)
			continue
		}
		logrus.Infof("Node %s - mc name: %s", node.Name, mcName)
		mc, err := getMachineConfig(mcName, machineConfigs)
		if err != nil {
			logrus.Errorf("Failed to get machineConfig %s, err: %v", mcName, err)
			continue
		}

		wrapperNodes[node.Name] = Node{
			Data: node,
			Mc:   mc,
		}
	}

	return wrapperNodes
}
