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

package podsets

import (
	"fmt"
	"time"

	"github.com/test-network-function/cnf-certification-test/internal/clientsholder"
	"github.com/test-network-function/cnf-certification-test/internal/log"
	"github.com/test-network-function/cnf-certification-test/pkg/provider"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	ReplicaSetString  = "ReplicaSet"
	StatefulsetString = "StatefulSet"
)

var WaitForDeploymentSetReady = func(ns, name string, timeout time.Duration, logger *log.Logger) bool {
	logger.Info("Check if Deployment %s:%s is ready", ns, name)
	clients := clientsholder.GetClientsHolder()
	start := time.Now()
	for time.Since(start) < timeout {
		dp, err := provider.GetUpdatedDeployment(clients.K8sClient.AppsV1(), ns, name)
		if err != nil {
			logger.Error("Error while getting Deployment %q, err: %v", name, err)
		} else if !dp.IsDeploymentReady() {
			logger.Warn("Deployment %q is not ready yet", dp.ToString())
		} else {
			logger.Info("Deployment %q is ready!", dp.ToString())
			return true
		}

		time.Sleep(time.Second)
	}
	logger.Error("Deployment %s:%s is not ready", ns, name)
	return false
}

var WaitForScalingToComplete = func(ns, name string, timeout time.Duration, groupResourceSchema schema.GroupResource, logger *log.Logger) bool {
	logger.Info("Check if scale object for CRs %s:%s is ready", ns, name)
	clients := clientsholder.GetClientsHolder()
	start := time.Now()
	for time.Since(start) < timeout {
		crScale, err := provider.GetUpdatedCrObject(clients.ScalingClient, ns, name, groupResourceSchema)
		if err != nil {
			logger.Error("Error while getting the scaling fields %v", err)
		} else if !crScale.IsScaleObjectReady() {
			logger.Warn("%s is not ready yet", crScale.ToString())
		} else {
			logger.Info("%s is ready!", crScale.ToString())
			return true
		}

		time.Sleep(time.Second)
	}
	logger.Error("Timeout waiting for CR %s:%s scaling to be complete", ns, name)
	return false
}

func WaitForStatefulSetReady(ns, name string, timeout time.Duration, logger *log.Logger) bool {
	logger.Debug("Check if statefulset %s:%s is ready", ns, name)
	clients := clientsholder.GetClientsHolder()
	start := time.Now()
	for time.Since(start) < timeout {
		ss, err := provider.GetUpdatedStatefulset(clients.K8sClient.AppsV1(), ns, name)
		if err != nil {
			logger.Error("Error while getting the %s, err: %v", ss.ToString(), err)
		} else if ss.IsStatefulSetReady() {
			logger.Info("%s is ready", ss.ToString())
			return true
		}
		time.Sleep(time.Second)
	}
	logger.Error("Statefulset %s:%s is not ready", ns, name)
	return false
}

func isDeploymentReady(name, namespace string) (bool, error) {
	appsV1Api := clientsholder.GetClientsHolder().K8sClient.AppsV1()

	dep, err := provider.GetUpdatedDeployment(appsV1Api, namespace, name)
	if err != nil {
		return false, err
	}

	return dep.IsDeploymentReady(), nil
}

func isStatefulSetReady(name, namespace string) (bool, error) {
	appsV1Api := clientsholder.GetClientsHolder().K8sClient.AppsV1()

	sts, err := provider.GetUpdatedStatefulset(appsV1Api, namespace, name)
	if err != nil {
		return false, err
	}

	return sts.IsStatefulSetReady(), nil
}

// Helper function to get a slice of namespace:name strings from a slice of *provider.Deployments.
// E.g: [tnf:test tnf:hazelcast-platform-controller-manager].
func getDeploymentsInfo(deployments []*provider.Deployment) []string {
	deps := []string{}
	for _, dep := range deployments {
		deps = append(deps, fmt.Sprintf("%s:%s", dep.Namespace, dep.Name))
	}

	return deps
}

// Helper function to get a slice of namespace: name strings from a slice of *provider.Statefulsets.
func getStatefulSetsInfo(statefulSets []*provider.StatefulSet) []string {
	stsInfo := []string{}
	for _, sts := range statefulSets {
		stsInfo = append(stsInfo, fmt.Sprintf("%s:%s", sts.Namespace, sts.Name))
	}

	return stsInfo
}

// Helper function that checks the status of each deployment in the slice and returns
// a slice with the not-ready ones.
func getNotReadyDeployments(deployments []*provider.Deployment) []*provider.Deployment {
	notReadyDeployments := []*provider.Deployment{}
	for _, dep := range deployments {
		ready, err := isDeploymentReady(dep.Name, dep.Namespace)
		if err != nil {
			log.Error("Failed to get %s: %v", dep.ToString(), err)
			// We'll mark it as not ready, anyways.
			notReadyDeployments = append(notReadyDeployments, dep)
			continue
		}

		if ready {
			log.Debug("%s is ready.", dep.ToString())
		} else {
			notReadyDeployments = append(notReadyDeployments, dep)
		}
	}

	return notReadyDeployments
}

// Helper function that checks the status of each statefulSet in the slice and returns
// a slice with the not-ready ones.
func getNotReadyStatefulSets(statefulSets []*provider.StatefulSet) []*provider.StatefulSet {
	notReadyStatefulSets := []*provider.StatefulSet{}
	for _, sts := range statefulSets {
		ready, err := isStatefulSetReady(sts.Name, sts.Namespace)
		if err != nil {
			log.Error("Failed to get %s: %v", sts.ToString(), err)
			// We'll mark it as not ready, anyways.
			notReadyStatefulSets = append(notReadyStatefulSets, sts)
			continue
		}

		if ready {
			log.Debug("%s is ready.", sts.ToString())
		} else {
			notReadyStatefulSets = append(notReadyStatefulSets, sts)
		}
	}

	return notReadyStatefulSets
}

func WaitForAllPodSetsReady(env *provider.TestEnvironment, timeout time.Duration, logger *log.Logger) (
	notReadyDeployments []*provider.Deployment,
	notReadyStatefulSets []*provider.StatefulSet) {
	const queryInterval = 15 * time.Second

	deploymentsToCheck := env.Deployments
	statefulSetsToCheck := env.StatefulSets

	logger.Info("Waiting %s for %d podsets to be ready.", timeout, len(deploymentsToCheck)+len(statefulSetsToCheck))
	for startTime := time.Now(); time.Since(startTime) < timeout; {
		logger.Info("Checking Deployments readiness of Deployments %v", getDeploymentsInfo(deploymentsToCheck))
		notReadyDeployments = getNotReadyDeployments(deploymentsToCheck)

		logger.Info("Checking StatefulSets readiness of StatefulSets %v", getStatefulSetsInfo(statefulSetsToCheck))
		notReadyStatefulSets = getNotReadyStatefulSets(statefulSetsToCheck)

		logger.Info("Not ready Deployments: %v", getDeploymentsInfo(notReadyDeployments))
		logger.Info("Not ready StatefulSets: %v", getStatefulSetsInfo(notReadyStatefulSets))

		deploymentsToCheck = notReadyDeployments
		statefulSetsToCheck = notReadyStatefulSets

		if len(deploymentsToCheck) == 0 && len(statefulSetsToCheck) == 0 {
			// No more podsets to check.
			break
		}

		time.Sleep(queryInterval)
	}

	// Here, either we reached the timeout or there's no more not-ready deployments or statefulsets.
	logger.Error("Not ready Deployments: %v", getDeploymentsInfo(deploymentsToCheck))
	logger.Error("Not ready StatefulSets: %v", getStatefulSetsInfo(statefulSetsToCheck))

	return deploymentsToCheck, statefulSetsToCheck
}

func GetAllNodesForAllPodSets(pods []*provider.Pod) (nodes map[string]bool) {
	nodes = make(map[string]bool)
	for _, put := range pods {
		for _, or := range put.OwnerReferences {
			if or.Kind != ReplicaSetString && or.Kind != StatefulsetString {
				continue
			}
			nodes[put.Spec.NodeName] = true
			break
		}
	}
	return nodes
}
