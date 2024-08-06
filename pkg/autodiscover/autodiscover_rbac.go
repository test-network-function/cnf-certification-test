// Copyright (C) 2023-2024 Red Hat, Inc.
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

	"github.com/test-network-function/cnf-certification-test/internal/log"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	rbacv1typed "k8s.io/client-go/kubernetes/typed/rbac/v1"
)

// getRoleBindings returns all of the rolebindings in the cluster.
func getRoleBindings(client rbacv1typed.RbacV1Interface) ([]rbacv1.RoleBinding, error) {
	// Get all of the rolebindings from all namespaces.
	roleList, roleErr := client.RoleBindings("").List(context.TODO(), metav1.ListOptions{})
	if roleErr != nil {
		log.Error("Executing rolebinding command failed with error: %v", roleErr)
		return nil, roleErr
	}
	return roleList.Items, nil
}

// getClusterRoleBindings returns all of the clusterrolebindings in the cluster.
func getClusterRoleBindings(client rbacv1typed.RbacV1Interface) ([]rbacv1.ClusterRoleBinding, error) {
	// Get all of the clusterrolebindings from the cluster.
	// These are not namespaced so we want all of them.
	crbList, crbErr := client.ClusterRoleBindings().List(context.TODO(), metav1.ListOptions{})
	if crbErr != nil {
		log.Error("Executing clusterrolebinding command failed with error: %v", crbErr)
		return nil, crbErr
	}
	return crbList.Items, nil
}

// getRoles returns all of the roles in the cluster.
func getRoles(client rbacv1typed.RbacV1Interface) ([]rbacv1.Role, error) {
	// Get all of the roles from all namespaces.
	roleList, roleErr := client.Roles("").List(context.TODO(), metav1.ListOptions{})
	if roleErr != nil {
		log.Error("Executing roles command failed with error: %v", roleErr)
		return nil, roleErr
	}
	return roleList.Items, nil
}
