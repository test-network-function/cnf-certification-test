package podhelper

import (
	"context"
	"fmt"
	"strings"

	"github.com/test-network-function/cnf-certification-test/internal/clientsholder"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

// Structure to describe a top owner of a pod.
type TopOwner struct {
	Kind      string
	Name      string
	Namespace string
}

// Get the list of top owners of pods.
func GetPodTopOwner(podNamespace string, podOwnerReferences []metav1.OwnerReference) (topOwners map[string]TopOwner, err error) {
	topOwners = make(map[string]TopOwner)
	err = followOwnerReferences(clientsholder.GetClientsHolder().GroupResources, clientsholder.GetClientsHolder().DynamicClient, topOwners, podNamespace, podOwnerReferences)
	if err != nil {
		return topOwners, fmt.Errorf("could not get top owners, err=%s", err)
	}
	return topOwners, nil
}

// Recursively follow the ownership tree to find the top owners.
func followOwnerReferences(resourceList []*metav1.APIResourceList, dynamicClient dynamic.Interface, topOwners map[string]TopOwner, namespace string, ownerRefs []metav1.OwnerReference) (err error) {
	for _, ownerRef := range ownerRefs {
		// Get group resource version.
		gvr := getResourceSchema(resourceList, ownerRef.APIVersion, ownerRef.Kind)
		// Get the owner resources.
		resource, err := dynamicClient.Resource(gvr).Namespace(namespace).Get(context.Background(), ownerRef.Name, metav1.GetOptions{})
		if err != nil {
			return fmt.Errorf("could not get object indicated by owner references")
		}
		// Get owner references of the unstructured object.
		ownerReferences := resource.GetOwnerReferences()
		if err != nil {
			return fmt.Errorf("error getting owner references. err= %s", err)
		}
		// if no owner references, we have reached the top record it.
		if len(ownerReferences) == 0 {
			topOwners[ownerRef.Name] = TopOwner{Kind: ownerRef.Kind, Name: ownerRef.Name, Namespace: namespace}
		}
		// if not continue following other branches.
		err = followOwnerReferences(resourceList, dynamicClient, topOwners, namespace, ownerReferences)
		if err != nil {
			return fmt.Errorf("error following owners")
		}
	}
	return nil
}

// Get the Group Version Resource based on APIVersion and kind.
func getResourceSchema(resourceList []*metav1.APIResourceList, apiVersion, kind string) (gvr schema.GroupVersionResource) {
	const groupVersionComponentsNumber = 2
	for _, gr := range resourceList {
		for i := 0; i < len(gr.APIResources); i++ {
			if gr.APIResources[i].Kind == kind && gr.GroupVersion == apiVersion {
				groupSplit := strings.Split(gr.GroupVersion, "/")
				if len(groupSplit) == groupVersionComponentsNumber {
					gvr.Group = groupSplit[0]
					gvr.Version = groupSplit[1]
					gvr.Resource = gr.APIResources[i].Name
				}
				return gvr
			}
		}
	}
	return gvr
}
