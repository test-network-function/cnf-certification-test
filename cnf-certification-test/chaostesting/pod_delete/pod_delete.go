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
package poddelete

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/test-network-function/cnf-certification-test/pkg/provider"
	"github.com/test-network-function/cnf-certification-test/pkg/tnf"

	"github.com/sirupsen/logrus"
	"github.com/test-network-function/cnf-certification-test/internal/clientsholder"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	yamlutil "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/restmapper"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

const (
	// costum rescources file to create them
	serviceAccountFile = "chaostesting/chaos-test-files/service-account.yaml"
	experimentFile     = "chaostesting/chaos-test-files/experiment-delete.yaml"
	chaosEngineFile    = "chaostesting/chaos-test-files/chaos-engine.yaml"
	chaosTestName      = "pod-delete" // test name
	completedResult    = "completed"
	pass               = "Pass"
	chaosresultName    = "engine-test-pod-delete"
)

// a function to search the right label for the deployment that we wanr to test pod delete on it
func GetLabelDeploymetValue(env *provider.TestEnvironment, labelsMap map[string]string) (string, error) {
	var key string
	for _, label := range env.Config.TargetPodLabels {
		if label.Prefix != "" {
			key = fmt.Sprintf("%s/%s", label.Prefix, label.Name)
		} else {
			key = label.Name
		}
		if l, ok := labelsMap[key]; ok && l == label.Value {
			return fmt.Sprintf("%s=%s", key, label.Value), nil
		}
	}
	return "", errors.New("didnt find a key and value that matching the deployment")
}

func ApplyAndCreatePodDeleteRecources(appLabel, appKind, namespace string) error {
	// create the chaos experiment recource
	if err := ApplyAndCreateFile(appLabel, appKind, namespace, experimentFile); err != nil {
		logrus.Errorf("cant create the experiment of the test: %s", err)
		return err
	}
	// create the chaos serviceAccount recource
	if err := ApplyAndCreateFile(appLabel, appKind, namespace, serviceAccountFile); err != nil {
		logrus.Errorf("cant create the serviceAccount of the test: %s", err)
		return err
	}
	// create the chaos chaosEngine recource
	if err := ApplyAndCreateFile(appLabel, appKind, namespace, chaosEngineFile); err != nil {
		logrus.Errorf("cant create the chaosEngine of the test: %s", err)
		return err
	}
	return nil
}

func ApplyAndCreateFile(appLabel, appKind, namespace, filename string) error {
	fileDecoder, err := applyTemplate(appLabel, appKind, namespace, filename)
	if err != nil {
		logrus.Errorf("cant create the decoderfile of the test: %s", err)
		return err
	}
	if err = createResource(fileDecoder); err != nil {
		logrus.Errorf("%s error create the resources for the test.", err)
		return err
	}
	return nil
}

func DeleteAllResources(namespace string) {
	oc := clientsholder.GetClientsHolder()
	deletePolicy := metav1.DeletePropagationForeground
	deleteOptions := metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}
	gvr := schema.GroupVersionResource{Group: "litmuschaos.io", Version: "v1alpha1", Resource: "chaosengines"}
	if err := oc.DynamicClient.Resource(gvr).Namespace(namespace).Delete(context.TODO(), "engine-test", deleteOptions); err != nil {
		logrus.Errorf("error while removing the chaos engine resources %e", err)
	}
	err := oc.K8sClient.CoreV1().ServiceAccounts(namespace).Delete(context.TODO(), "test-sa", deleteOptions)
	if err != nil {
		logrus.Errorf("error while removing the ServiceAccountsresources %e", err)
	}
	if err = oc.K8sClient.RbacV1().Roles(namespace).Delete(context.TODO(), "test-sa", deleteOptions); err != nil {
		logrus.Errorf("error while removing the chaos engine resources %e", err)
	}
	if err = oc.K8sClient.RbacV1().RoleBindings(namespace).Delete(context.TODO(), "test-sa", deleteOptions); err != nil {
		logrus.Errorf("error while removing the chaos engine resources %e", err)
	}
	gvr = schema.GroupVersionResource{Group: "litmuschaos.io", Version: "v1alpha1", Resource: "chaosexperiments"}
	if err := oc.DynamicClient.Resource(gvr).Namespace(namespace).Delete(context.TODO(), chaosTestName, deleteOptions); err != nil {
		logrus.Errorf("error while removing the chaos engine resources %e", err)
	}
	gvr = schema.GroupVersionResource{Group: "litmuschaos.io", Version: "v1alpha1", Resource: "chaosresults"}
	if err := oc.DynamicClient.Resource(gvr).Namespace(namespace).Delete(context.TODO(), chaosresultName, deleteOptions); err != nil {
		logrus.Errorf("error while removing the chaos results resources %e", err)
	}

}

func applyTemplate(appLabel, appKind, namespace, filename string) (*yamlutil.YAMLOrJSONDecoder, error) {
	input, err := os.ReadFile(filename)
	if err != nil {
		logrus.Errorf("error while reading the yaml file : %s ,%s", filename, err)
		tnf.ClaimFilePrintf("error while reading the yaml file : %s ,%s", filename, err)
		return nil, err
	}
	output := bytes.ReplaceAll(input, []byte("{{APP_NAMESPACE}}"), []byte(namespace))
	output = bytes.ReplaceAll(output, []byte("{{APP_LABEL}}"), []byte(appLabel))
	output = bytes.ReplaceAll(output, []byte("{{APP_KIND}}"), []byte(appKind))
	const oneh = 100
	fileDecoder := yamlutil.NewYAMLOrJSONDecoder(bytes.NewReader(output), oneh)
	return fileDecoder, nil
}

func WaitForTestFinish(timeout time.Duration) bool {
	const pollingPeriod = 1 * time.Second
	var elapsed time.Duration
	var result bool
	for elapsed < timeout {
		result = waitForResult()

		if result {
			break
		}
		time.Sleep(pollingPeriod)
		elapsed += pollingPeriod
	}
	return result
}

func IsChaosResultVerdictPass() bool {
	oc := clientsholder.GetClientsHolder()
	gvr := schema.GroupVersionResource{Group: "litmuschaos.io", Version: "v1alpha1", Resource: "chaosresults"}
	crs, err := oc.DynamicClient.Resource(gvr).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		logrus.Errorf("error getting : %v\n", err)
		return false
	}
	if len(crs.Items) > 1 {
		logrus.Errorf("u have %d more than one chaosresults resource that is wrong bahavior \n", len(crs.Items))
		return false
	}
	cr := crs.Items[0]
	failResult := cr.Object["status"].(map[string]interface{})["experimentStatus"].(map[string]interface{})["failStep"]
	verdictValue := cr.Object["status"].(map[string]interface{})["experimentStatus"].(map[string]interface{})["verdict"]
	expKind := cr.Object["spec"].(map[string]interface{})["experiment"]
	if expKind == chaosTestName {
		if verdictValue == pass {
			return true
		}
		logrus.Debugf("test completed but it failed with reason %s", failResult.(string))
		tnf.ClaimFilePrintf("test completed but it failed with reason %s", failResult.(string))
		return false
	}
	return false
}

func waitForResult() bool {
	oc := clientsholder.GetClientsHolder()
	gvr := schema.GroupVersionResource{Group: "litmuschaos.io", Version: "v1alpha1", Resource: "chaosengines"}
	crs, err := oc.DynamicClient.Resource(gvr).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		logrus.Errorf("error getting : %v\n", err)
		return false
	}
	if len(crs.Items) > 1 {
		logrus.Errorf("u have %d more than one chaosenging resource that is wrong bahavior \n", len(crs.Items))
		return false
	}
	cr := crs.Items[0]
	if status := cr.Object["status"]; status != nil {
		if exp := status.(map[string]interface{})["experiments"]; exp != nil {
			typ := exp.([]interface{})
			status := cr.Object["status"].(map[string]interface{})["engineStatus"]
			if typ[0].(map[string]interface{})["name"] == chaosTestName {
				return status == completedResult
			}
		}
	}
	return false
}

// createResource is a helper function that uses a yaml decoder to create in the cluster
// all the resources defined in the underlying yaml file.
//nolint:funlen //
func createResource(decoder *yamlutil.YAMLOrJSONDecoder) error {
	oc := clientsholder.GetClientsHolder()
	k8sClient := oc.K8sClient
	dynamicClient := oc.DynamicClient
	for {
		var rawObj runtime.RawExtension
		if err := decoder.Decode(&rawObj); err != nil {
			if err != io.EOF {
				return err
			}
			return nil
		}
		obj, gvk, err := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme).Decode(rawObj.Raw, nil, nil)
		if err != nil {
			return err
		}
		unstructuredMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
		if err != nil {
			return err
		}

		unstructuredObj := &unstructured.Unstructured{Object: unstructuredMap}

		gr, err := restmapper.GetAPIGroupResources(k8sClient.Discovery())
		if err != nil {
			return err
		}

		mapper := restmapper.NewDiscoveryRESTMapper(gr)
		mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
		if err != nil {
			return err
		}

		var dri dynamic.ResourceInterface
		if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
			if unstructuredObj.GetNamespace() == "" {
				unstructuredObj.SetNamespace("default")
			}
			dri = dynamicClient.Resource(mapping.Resource).Namespace(unstructuredObj.GetNamespace())
		} else {
			dri = dynamicClient.Resource(mapping.Resource)
		}

		if _, err := dri.Create(context.Background(), unstructuredObj, metav1.CreateOptions{}); err != nil {
			return err
		}
	}
}
