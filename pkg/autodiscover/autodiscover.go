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

package autodiscover

import (
	"context"
	"errors"
	"fmt"
	"time"

	configv1 "github.com/openshift/api/config/v1"
	clientconfigv1 "github.com/openshift/client-go/config/clientset/versioned/typed/config/v1"
	olmv1Alpha "github.com/operator-framework/api/pkg/operators/v1alpha1"
	"github.com/sirupsen/logrus"
	"github.com/test-network-function/cnf-certification-test/internal/clientsholder"
	"github.com/test-network-function/cnf-certification-test/pkg/compatibility"
	"github.com/test-network-function/cnf-certification-test/pkg/configuration"
	"helm.sh/helm/v3/pkg/release"
	appsv1 "k8s.io/api/apps/v1"
	scalingv1 "k8s.io/api/autoscaling/v1"
	corev1 "k8s.io/api/core/v1"
	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	tnfCsvTargetLabelName  = "operator"
	tnfCsvTargetLabelValue = ""
	tnfLabelPrefix         = "test-network-function.com"
	labelTemplate          = "%s/%s"
	// anyLabelValue is the value that will allow any value for a label when building the label query.
	anyLabelValue = ""
)

type DiscoveredTestData struct {
	Env                configuration.TestParameters
	TestData           configuration.TestConfiguration
	Pods               []corev1.Pod
	DebugPods          []corev1.Pod
	LimitRangeItems    []corev1.LimitRange
	ResourceQuotaItems []corev1.ResourceQuota
	Crds               []*apiextv1.CustomResourceDefinition
	Namespaces         []string
	Csvs               []olmv1Alpha.ClusterServiceVersion
	Deployments        []appsv1.Deployment
	StatefulSet        []appsv1.StatefulSet
	Hpas               map[string]*scalingv1.HorizontalPodAutoscaler
	Subscriptions      []olmv1Alpha.Subscription
	HelmChartReleases  map[string][]*release.Release
	K8sVersion         string
	OpenshiftVersion   string
	OCPStatus          string
	Nodes              *corev1.NodeList
	Istio              bool
}

var data = DiscoveredTestData{}

func buildLabelName(labelPrefix, labelName string) string {
	if labelPrefix == "" {
		return labelName
	}
	return fmt.Sprintf(labelTemplate, labelPrefix, labelName)
}

func buildLabelQuery(label configuration.Label) string {
	fullLabelName := buildLabelName(label.Prefix, label.Name)
	if label.Value != anyLabelValue {
		return fmt.Sprintf("%s=%s", fullLabelName, label.Value)
	}
	return fullLabelName
}
func buildLabelKeyValue(label configuration.Label) (key, value string) {
	key = buildLabelName(label.Prefix, label.Name)
	value = label.Value
	return key, value
}

//nolint:funlen
// DoAutoDiscover finds objects under test
func DoAutoDiscover() DiscoveredTestData {
	data.Env = *configuration.GetTestParameters()

	var err error
	data.TestData, err = configuration.LoadConfiguration(data.Env.ConfigurationPath)
	if err != nil {
		logrus.Fatalln("Cannot load configuration")
	}
	oc := clientsholder.GetClientsHolder()
	data.Namespaces = namespacesListToStringList(data.TestData.TargetNameSpaces)
	data.Pods = findPodsByLabel(oc.K8sClient.CoreV1(), data.TestData.TargetPodLabels, data.Namespaces)

	debugLabel := configuration.Label{Prefix: debugLabelPrefix, Name: debugLabelName, Value: debugLabelValue}
	debugLabels := []configuration.Label{debugLabel}
	debugNS := []string{defaultNamespace}
	data.DebugPods = findPodsByLabel(oc.K8sClient.CoreV1(), debugLabels, debugNS)
	data.LimitRangeItems, err = getLimitRanges(oc.K8sClient.CoreV1())
	if err != nil {
		logrus.Fatalln("Cannot get limit ranges")
	}
	data.ResourceQuotaItems, err = getResourceQuotas(oc.K8sClient.CoreV1())
	if err != nil {
		logrus.Fatalln("Cannot get resource quotas")
	}
	data.Crds = FindTestCrdNames(data.TestData.CrdFilters)
	data.Csvs = findOperatorsByLabel(oc.OlmClient, []configuration.Label{{Name: tnfCsvTargetLabelName, Prefix: tnfLabelPrefix, Value: tnfCsvTargetLabelValue}}, data.TestData.TargetNameSpaces)
	data.Subscriptions = findSubscriptions(oc.OlmClient, data.Namespaces)
	data.HelmChartReleases = getHelmList(oc.RestConfig, data.Namespaces)
	openshiftVersion, _ := getOpenshiftVersion(oc.OcpClient)
	data.OpenshiftVersion = openshiftVersion
	k8sVersion, err := oc.K8sClient.Discovery().ServerVersion()
	if err != nil {
		logrus.Fatalln("Cannot get the K8s version")
	}
	data.Istio = findIstioNamespace(oc.K8sClient.CoreV1())

	// Find the status of the OCP version (pre-ga, end-of-life, maintenance, or generally available)
	data.OCPStatus = compatibility.DetermineOCPStatus(openshiftVersion, time.Now())

	data.K8sVersion = k8sVersion.GitVersion
	data.Deployments = findDeploymentByLabel(oc.K8sClient.AppsV1(), data.TestData.TargetPodLabels, data.Namespaces)
	data.StatefulSet = findStatefulSetByLabel(oc.K8sClient.AppsV1(), data.TestData.TargetPodLabels, data.Namespaces)
	data.Hpas = findHpaControllers(oc.K8sClient, data.Namespaces)
	data.Nodes, err = oc.K8sClient.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logrus.Fatalln("Cannot get list of nodes")
	}
	return data
}

func namespacesListToStringList(namespaceList []configuration.Namespace) (stringList []string) {
	for _, ns := range namespaceList {
		stringList = append(stringList, ns.Name)
	}
	return stringList
}

func getOpenshiftVersion(oClient clientconfigv1.ConfigV1Interface) (ver string, err error) {
	var clusterOperator *configv1.ClusterOperator
	clusterOperator, err = oClient.ClusterOperators().Get(context.TODO(), "openshift-apiserver", metav1.GetOptions{})
	// error here indicates logged in as non-admin, log and move on
	if err != nil {
		switch {
		case kerrors.IsForbidden(err), kerrors.IsNotFound(err):
			logrus.Infof("OpenShift Version not found (must be logged in to cluster as admin): %v", err)
			err = nil
		}
	}
	if clusterOperator != nil {
		for _, ver := range clusterOperator.Status.Versions {
			if ver.Name == tnfCsvTargetLabelName {
				// openshift-apiserver does not report version,
				// clusteroperator/openshift-apiserver does, and only version number
				return ver.Version, nil
			}
		}
	}
	return "", errors.New("could not get openshift version")
}
