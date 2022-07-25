// Copyright (C) 2020-2022 Red Hat, Inc.
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

package podrecreation

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/onsi/ginkgo/v2"
	"github.com/sirupsen/logrus"
	"github.com/test-network-function/cnf-certification-test/internal/clientsholder"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	retry "k8s.io/client-go/util/retry"
)

const (
	ReplicaSetString            = "ReplicaSet"
	DeploymentString            = "Deployment"
	StatefulsetString           = "StatefulSet"
	DaemonSetString             = "DaemonSet"
	DefaultGracePeriodInSeconds = 30
	Cordon                      = "cordon"
	Uncordon                    = "uncordon"
	DeleteBackground            = "deleteBackground"
	DeleteForeground            = "deleteForeground"
	NoDelete                    = "noDelete"
)

func CordonHelper(name, operation string) error {
	clients := clientsholder.GetClientsHolder()

	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		// Fetch node object
		node, err := clients.K8sClient.CoreV1().Nodes().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return err
		}
		switch operation {
		case Cordon:
			node.Spec.Unschedulable = true
		case Uncordon:
			node.Spec.Unschedulable = false
		default:
			return fmt.Errorf("cordonHelper: Unsupported operation:%s", operation)
		}
		// Update the node
		_, err = clients.K8sClient.CoreV1().Nodes().Update(context.TODO(), node, metav1.UpdateOptions{})
		return err
	})
	if retryErr != nil {
		logrus.Error("can't ", operation, " node: ", name, " error=", retryErr)
	}
	return retryErr
}

func CountPodsWithDelete(nodeName, mode string) (count int, err error) {
	clients := clientsholder.GetClientsHolder()
	pods, err := clients.K8sClient.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{
		FieldSelector: "spec.nodeName=" + nodeName, LabelSelector: "pod-template-hash",
	})
	if err != nil {
		logrus.Errorf("error getting list of pods err: %s", err)
		return 0, err
	}
	count = 0
	var wg sync.WaitGroup
	for idx := range pods.Items {
		if skipDaemonPod(&pods.Items[idx]) {
			continue
		}
		count++
		if mode == NoDelete {
			continue
		}
		err := deletePod(&pods.Items[idx], mode, &wg)
		if err != nil {
			logrus.Errorf("error deleting %s", &pods.Items[idx])
		}
	}

	wg.Wait()
	return count, nil
}
func skipDaemonPod(pod *corev1.Pod) bool {
	for _, or := range pod.OwnerReferences {
		if or.Kind == DaemonSetString {
			return true
		}
	}
	return false
}

func deletePod(pod *corev1.Pod, mode string, wg *sync.WaitGroup) error {
	clients := clientsholder.GetClientsHolder()
	logrus.Debugf("deleting ns=%s pod=%s with %s mode", pod.Namespace, pod.Name, mode)
	gracePeriodSeconds := *pod.Spec.TerminationGracePeriodSeconds
	// Create watcher before deleting pod
	watcher, err := clients.K8sClient.CoreV1().Pods(pod.Namespace).Watch(context.TODO(), metav1.ListOptions{
		FieldSelector: "metadata.name=" + pod.Name + ",metadata.namespace=" + pod.Namespace,
	})
	if err != nil {
		return fmt.Errorf("waitPodDeleted ns=%s pod=%s, err=%s", pod.Namespace, pod.Name, err)
	}
	// Actually deleting pod
	err = clients.K8sClient.CoreV1().Pods(pod.Namespace).Delete(context.TODO(), pod.Name, metav1.DeleteOptions{
		GracePeriodSeconds: &gracePeriodSeconds,
	})
	if err != nil {
		logrus.Errorf("error deleting %s err: %v", pod.String(), err)
		return err
	}
	if mode == DeleteBackground {
		return nil
	}
	wg.Add(1)
	podName := pod.Name
	namespace := pod.Namespace
	go func() {
		err = waitPodDeleted(namespace, podName, gracePeriodSeconds, watcher)
		if err != nil {
			logrus.Errorf("waitPodDeleted failed with err=%s", err)
		}
		wg.Done()
	}()
	return nil
}

func CordonCleanup(node string) {
	err := CordonHelper(node, Uncordon)
	if err != nil {
		ginkgo.AbortSuite(fmt.Sprintf("cleanup: error uncordoning the node: %s", node))
	}
}

func waitPodDeleted(ns, podName string, timeout int64, watcher watch.Interface) error {
	logrus.Tracef("Entering waitPodDeleted ns=%s pod=%s", ns, podName)
	defer watcher.Stop()

	for {
		select {
		case event := <-watcher.ResultChan():
			if event.Type == watch.Deleted || event.Type == "" {
				logrus.Debugf("ns=%s pod=%s deleted", ns, podName)
				return nil
			}
		case <-time.After(time.Duration(timeout) * time.Second):
			logrus.Infof("watch for pod deletion timedout after %d seconds", timeout)
			return nil
		}
	}
}
