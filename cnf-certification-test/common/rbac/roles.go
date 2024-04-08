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
	"strings"

	rbacv1 "k8s.io/api/rbac/v1"
	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

type RoleRule struct {
	Resource RoleResource
	Verb     string
}
type RoleResource struct {
	Group, Name string
}

type CrdResource struct {
	Group, SingularName, PluralName string
	ShortNames                      []string
}

// GetCrdResources converts a list of apiextv1.CustomResourceDefinition structs into a list of list of CrdResource structs.
// Returns:
//   - []CrdResource : a slice of CrdResource objects.
func GetCrdResources(crds []*apiextv1.CustomResourceDefinition) (resourceList []CrdResource) {
	for _, crd := range crds {
		var aResource CrdResource
		aResource.Group = crd.Spec.Group
		aResource.SingularName = crd.Spec.Names.Singular
		aResource.PluralName = crd.Spec.Names.Plural
		aResource.ShortNames = crd.Spec.Names.ShortNames
		resourceList = append(resourceList, aResource)
	}
	return resourceList
}

// GetAllRules retrieves a list all of rules defined by the role passed in input.
// Returns:
//   - []RoleRule : a slice of RoleRule objects.
func GetAllRules(aRole *rbacv1.Role) (ruleList []RoleRule) {
	for _, aRule := range aRole.Rules {
		for _, aGroup := range aRule.APIGroups {
			for _, aResource := range aRule.Resources {
				for _, aVerb := range aRule.Verbs {
					var aRoleRule RoleRule
					aRoleRule.Resource.Group = aGroup
					aRoleRule.Resource.Name = aResource
					aRoleRule.Verb = aVerb
					ruleList = append(ruleList, aRoleRule)
				}
			}
		}
	}
	return ruleList
}

// isResourceInRoleRule Checks if a CRD resource is matched by a rule by comparing its group and plural name.
// Returns:
//   - bool : if a CrdResource matches a RoleRule based on their properties return true , otherwise return false.
func isResourceInRoleRule(crd CrdResource, roleRule RoleRule) bool {
	// remove subresources to keep only resource (plural) name
	ruleResourcePluralName := strings.Split(roleRule.Resource.Name, "/")[0]

	return crd.Group == roleRule.Resource.Group && crd.PluralName == ruleResourcePluralName
}

// FilterRulesNonMatchingResources filters RoleRules based on whether they match any CrdResource in the resourceList.
// Returns :
//   - Matching: a slice of RoleRule that contains all rules where a CrdResource matches a RoleRule based on their properties.
//   - NonMatching: a slice of RoleRule that contains all rules not matching the CRD resource.
func FilterRulesNonMatchingResources(ruleList []RoleRule, resourceList []CrdResource) (matching, nonMatching []RoleRule) {
	for _, aRule := range ruleList {
		for _, aResource := range resourceList {
			if isResourceInRoleRule(aResource, aRule) {
				matching = append(matching, aRule)
			}
		}
	}
	nonMatching = SliceDifference(ruleList, matching)
	return matching, nonMatching
}

// SliceDifference checks if there is a difference between s1 and s2 RoleRule slices.
// Returns :
//   - []RoleRule : the elements that are exist in s1 but not in s2.
func SliceDifference(s1, s2 []RoleRule) (diff []RoleRule) {
	var temp []RoleRule
	if len(s2) > len(s1) {
		temp = s1
		s1 = s2
		s2 = temp
	}
	for _, v1 := range s1 {
		missing := true
		for _, v2 := range s2 {
			if v1 == v2 {
				missing = false
				break
			}
		}
		if missing {
			diff = append(diff, v1)
		}
	}
	return diff
}
