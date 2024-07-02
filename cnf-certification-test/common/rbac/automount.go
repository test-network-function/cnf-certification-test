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

package rbac

import (
	"context"
	"fmt"

	"github.com/test-network-function/cnf-certification-test/internal/log"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1typed "k8s.io/client-go/kubernetes/typed/core/v1"
)

// AutomountServiceAccountSetOnSA checks if the AutomountServiceAccountToken field is set on a ServiceAccount.
// Returns:
//   - A boolean pointer indicating whether the AutomountServiceAccountToken field is set.
//   - An error if any occurred during the operation.
func AutomountServiceAccountSetOnSA(client corev1typed.CoreV1Interface, serviceAccountName, podNamespace string) (*bool, error) {
	sa, err := client.ServiceAccounts(podNamespace).Get(context.TODO(), serviceAccountName, metav1.GetOptions{})
	if err != nil {
		log.Error("executing serviceaccount command failed with error: %v", err)
		return nil, err
	}
	return sa.AutomountServiceAccountToken, nil
}

// EvaluateAutomountTokens evaluates whether the automountServiceAccountToken is correctly configured for the given Pod.
// Checks if the token is explicitly set in the Pod's spec or if it is inherited from the associated ServiceAccount.
// Returns:
//   - bool: Indicates whether the Pod passed all checks. if yes- return true, otherwise return false.
//   - string: Error message if the Pod is misconfigured, otherwise an empty string.
//
//nolint:gocritic
func EvaluateAutomountTokens(client corev1typed.CoreV1Interface, put *corev1.Pod) (bool, string) {
	// The token can be specified in the pod directly
	// or it can be specified in the service account of the pod
	// if no service account is configured, then the pod will use the configuration
	// of the default service account in that namespace
	// the token defined in the pod has takes precedence
	// the test would pass iif token is explicitly set to false
	// if the token is set to true in the pod, the test would fail right away
	if put.Spec.AutomountServiceAccountToken != nil && *put.Spec.AutomountServiceAccountToken {
		return false, fmt.Sprintf("Pod %s:%s is configured with automountServiceAccountToken set to true", put.Namespace, put.Name)
	}

	// Collect information about the service account attached to the pod.
	saAutomountServiceAccountToken, err := AutomountServiceAccountSetOnSA(client, put.Spec.ServiceAccountName, put.Namespace)
	if err != nil {
		return false, ""
	}

	// The pod token is false means the pod is configured properly
	// The pod is not configured and the service account is configured with false means
	// the pod will inherit the behavior `false` and the test would pass
	if (put.Spec.AutomountServiceAccountToken != nil && !*put.Spec.AutomountServiceAccountToken) || (saAutomountServiceAccountToken != nil && !*saAutomountServiceAccountToken) {
		return true, ""
	}

	// the service account is configured with true means all the pods
	// using this service account are not configured properly, register the error
	// message and fail
	if saAutomountServiceAccountToken != nil && *saAutomountServiceAccountToken {
		return false, fmt.Sprintf("serviceaccount %s:%s is configured with automountServiceAccountToken set to true, impacting pod %s", put.Namespace, put.Spec.ServiceAccountName, put.Name)
	}

	// the token should be set explicitly to false, otherwise, it's a failure
	// register the error message and check the next pod
	if saAutomountServiceAccountToken == nil {
		return false, fmt.Sprintf("serviceaccount %s:%s is not configured with automountServiceAccountToken set to false, impacting pod %s", put.Namespace, put.Spec.ServiceAccountName, put.Name)
	}

	return true, "" // Pod has passed all checks
}
