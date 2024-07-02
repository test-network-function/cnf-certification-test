// Copyright (C) 2022-2023 Red Hat, Inc.
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
	"reflect"
	"strconv"
	"testing"

	rbacv1 "k8s.io/api/rbac/v1"
	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

func TestSliceDifference(t *testing.T) {
	type args struct {
		s1 []RoleRule
		s2 []RoleRule
	}
	tests := []struct {
		name     string
		args     args
		wantDiff []RoleRule
	}{
		{
			name: "ok",
			args: args{s1: []RoleRule{generateRoleRule("mycr.com", "name1", "verb1"),
				generateRoleRule("mycr.com", "name2", "verb1"),
				generateRoleRule("mycr.com", "name3", "verb2"),
				generateRoleRule("mycr1.com", "name4", "verb1"),
			},
				s2: []RoleRule{generateRoleRule("mycr.com", "name1", "verb1"),
					generateRoleRule("mycr1.com", "name4", "verb1")}},
			wantDiff: []RoleRule{generateRoleRule("mycr.com", "name2", "verb1"),
				generateRoleRule("mycr.com", "name3", "verb2"),
			},
		},
		{
			name: "reverse args",
			args: args{s2: []RoleRule{generateRoleRule("mycr.com", "name1", "verb1"),
				generateRoleRule("mycr.com", "name2", "verb1"),
				generateRoleRule("mycr.com", "name3", "verb2"),
				generateRoleRule("mycr1.com", "name4", "verb1"),
			},
				s1: []RoleRule{generateRoleRule("mycr.com", "name1", "verb1"),
					generateRoleRule("mycr1.com", "name4", "verb1")}},
			wantDiff: []RoleRule{generateRoleRule("mycr.com", "name2", "verb1"),
				generateRoleRule("mycr.com", "name3", "verb2"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotDiff := SliceDifference(tt.args.s1, tt.args.s2); !reflect.DeepEqual(
				gotDiff,
				tt.wantDiff,
			) {
				t.Errorf("SliceDifference() = %v, want %v", gotDiff, tt.wantDiff)
			}
		})
	}
}

func generateRoleRule(group, name, verb string) (rule RoleRule) {
	rule.Resource.Group = group
	rule.Resource.Name = name
	rule.Verb = verb
	return rule
}

func TestGetCrdResources(t *testing.T) {
	type args struct {
		crds []*apiextv1.CustomResourceDefinition
	}
	tests := []struct {
		name             string
		args             args
		wantResourceList []CrdResource
	}{
		{
			name: "ok",
			args: args{
				crds: []*apiextv1.CustomResourceDefinition{
					generateCRD("group1", "singular1", "plural1", "short1_", 3),
				},
			},
			wantResourceList: []CrdResource{
				generateCRDResource("group1", "singular1", "plural1", "short1_", 3),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotResourceList := GetCrdResources(tt.args.crds); !reflect.DeepEqual(
				gotResourceList,
				tt.wantResourceList,
			) {
				t.Errorf("GetCrdResources() = %v, want %v", gotResourceList, tt.wantResourceList)
			}
		})
	}
}

func generateCRD(
	group, singular, plural, short string,
	num int,
) *apiextv1.CustomResourceDefinition {
	var crd apiextv1.CustomResourceDefinition
	crd.Spec.Group = group
	crd.Spec.Names.Singular = singular
	crd.Spec.Names.Plural = plural
	for i := 0; i < num; i++ {
		crd.Spec.Names.ShortNames = append(crd.Spec.Names.ShortNames, short+strconv.Itoa(i))
	}
	return &crd
}

func generateCRDResource(group, singular, plural, short string, num int) CrdResource {
	var resource CrdResource
	resource.Group = group
	resource.SingularName = singular
	resource.PluralName = plural
	for i := 0; i < num; i++ {
		resource.ShortNames = append(resource.ShortNames, short+strconv.Itoa(i))
	}
	return resource
}

func TestGetAllRules(t *testing.T) {
	type args struct {
		aRole *rbacv1.Role
	}
	tests := []struct {
		name         string
		args         args
		wantRuleList []RoleRule
	}{
		{
			name: "ok",
			args: args{aRole: generateRole()},
			wantRuleList: []RoleRule{
				{Resource: RoleResource{Group: "group1", Name: "resource1"}, Verb: "verb1"},
				{Resource: RoleResource{Group: "group1", Name: "resource1"}, Verb: "verb2"},
				{Resource: RoleResource{Group: "group1", Name: "resource1"}, Verb: "verb3"},
				{Resource: RoleResource{Group: "group1", Name: "resource2"}, Verb: "verb1"},
				{Resource: RoleResource{Group: "group1", Name: "resource2"}, Verb: "verb2"},
				{Resource: RoleResource{Group: "group1", Name: "resource2"}, Verb: "verb3"},
				{Resource: RoleResource{Group: "group2", Name: "resource1"}, Verb: "verb1"},
				{Resource: RoleResource{Group: "group2", Name: "resource1"}, Verb: "verb2"},
				{Resource: RoleResource{Group: "group2", Name: "resource1"}, Verb: "verb3"},
				{Resource: RoleResource{Group: "group2", Name: "resource2"}, Verb: "verb1"},
				{Resource: RoleResource{Group: "group2", Name: "resource2"}, Verb: "verb2"},
				{Resource: RoleResource{Group: "group2", Name: "resource2"}, Verb: "verb3"},
				{Resource: RoleResource{Group: "group3", Name: "resource1"}, Verb: "verb1"},
				{Resource: RoleResource{Group: "group3", Name: "resource1"}, Verb: "verb2"},
				{Resource: RoleResource{Group: "group3", Name: "resource1"}, Verb: "verb3"},
				{Resource: RoleResource{Group: "group3", Name: "resource2"}, Verb: "verb1"},
				{Resource: RoleResource{Group: "group3", Name: "resource2"}, Verb: "verb2"},
				{Resource: RoleResource{Group: "group3", Name: "resource2"}, Verb: "verb3"},
				{Resource: RoleResource{Group: "group4", Name: "resource1"}, Verb: "verb1"},
				{Resource: RoleResource{Group: "group4", Name: "resource1"}, Verb: "verb2"},
				{Resource: RoleResource{Group: "group4", Name: "resource1"}, Verb: "verb3"},
				{Resource: RoleResource{Group: "group4", Name: "resource2"}, Verb: "verb1"},
				{Resource: RoleResource{Group: "group4", Name: "resource2"}, Verb: "verb2"},
				{Resource: RoleResource{Group: "group4", Name: "resource2"}, Verb: "verb3"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRuleList := GetAllRules(tt.args.aRole); !reflect.DeepEqual(
				gotRuleList,
				tt.wantRuleList,
			) {
				t.Errorf("GetAllRules() = %#v, want %v", gotRuleList, tt.wantRuleList)
			}
		})
	}
}

func generateRole() *rbacv1.Role {
	var role rbacv1.Role
	var rule1, rule2 rbacv1.PolicyRule

	rule1.APIGroups = append(rule1.APIGroups, "group1", "group2")
	rule1.Resources = append(rule1.Resources, "resource1", "resource2")
	rule1.Verbs = append(rule1.Verbs, "verb1", "verb2", "verb3")

	rule2.APIGroups = append(rule2.APIGroups, "group3", "group4")
	rule2.Resources = append(rule2.Resources, "resource1", "resource2")
	rule2.Verbs = append(rule2.Verbs, "verb1", "verb2", "verb3")

	role.Rules = append(role.Rules, rule1, rule2)
	return &role
}

func Test_isResourceInRoleRule(t *testing.T) {
	type args struct {
		aResource CrdResource
		aRule     RoleRule
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "ok",
			args: args{
				aResource: generateCRDResource("group1", "singular1", "plural1", "short1_", 3),
				aRule:     generateRoleRule("group1", "plural1", "verb1"),
			},
			want: true,
		},
		{
			name: "nok",
			args: args{
				aResource: generateCRDResource("group1", "singular1", "plural1", "short1_", 4),
				aRule:     generateRoleRule("group1", "plural2", "verb1"),
			},
			want: false,
		},
		{
			name: "using a short name is invalid",
			args: args{
				aResource: generateCRDResource("group1", "singular1", "plural1", "short1_", 3),
				aRule:     generateRoleRule("group1", "short1_2", "verb1"),
			},
			want: false,
		},
		{
			name: "ok1",
			args: args{
				aResource: generateCRDResource("group1", "singular1", "plural1", "short1_", 3),
				aRule:     generateRoleRule("group1", "plural1", "verb1"),
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isResourceInRoleRule(tt.args.aResource, tt.args.aRule); got != tt.want {
				t.Errorf("equalResource() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilterRulesNonMatchingResources(t *testing.T) {
	type args struct {
		ruleList     []RoleRule
		resourceList []CrdResource
	}
	tests := []struct {
		name            string
		args            args
		wantMatching    []RoleRule
		wantNonMatching []RoleRule
	}{
		{
			name: "ok",
			args: args{
				ruleList: []RoleRule{generateRoleRule("group1", "plural1", "verb1"),
					generateRoleRule("group2", "plural2", "verb1"),
					generateRoleRule("group2", "plural2", "verb2"),
					generateRoleRule("group1", "plural3/finalizer", "verb1"),
				},
				resourceList: []CrdResource{
					generateCRDResource("group1", "resource1", "plural1", "short1_", 3),
					generateCRDResource("group2", "resource2", "plural2", "short2_", 3),
					generateCRDResource("group1", "resource3", "plural3", "short3_", 3),
				},
			},
			wantMatching: []RoleRule{
				{Resource: RoleResource{Group: "group1", Name: "plural1"}, Verb: "verb1"},
				{Resource: RoleResource{Group: "group2", Name: "plural2"}, Verb: "verb1"},
				{Resource: RoleResource{Group: "group2", Name: "plural2"}, Verb: "verb2"},
				{
					Resource: RoleResource{Group: "group1", Name: "plural3/finalizer"},
					Verb:     "verb1",
				},
			},
			wantNonMatching: nil,
		},
		{
			name: "resource not matching (using singular)",
			args: args{
				ruleList: []RoleRule{generateRoleRule("group1", "plural1", "verb1"),
					generateRoleRule("group1", "plural2", "verb1"),
					generateRoleRule("group1", "resource1", "verb2"),
					generateRoleRule("group1", "plural3/finalizer", "verb1"),
				},
				resourceList: []CrdResource{
					generateCRDResource("group1", "resource1", "plural1", "short1_", 3),
					generateCRDResource("group1", "resource2", "plural2", "short2_", 3),
					generateCRDResource("group1", "resource3", "plural3", "short3_", 3),
				},
			},
			wantMatching: []RoleRule{
				{Resource: RoleResource{Group: "group1", Name: "plural1"}, Verb: "verb1"},
				{Resource: RoleResource{Group: "group1", Name: "plural2"}, Verb: "verb1"},
				{Resource: RoleResource{Group: "group1", Name: "plural3/finalizer"}, Verb: "verb1"},
			},
			wantNonMatching: []RoleRule{
				{Resource: RoleResource{Group: "group1", Name: "resource1"}, Verb: "verb2"},
			},
		},
		{
			name: "resource not matching group",
			args: args{
				ruleList: []RoleRule{generateRoleRule("group1", "plural1", "verb1"),
					generateRoleRule("group1", "plural2", "verb1"),
					generateRoleRule("group2", "plural2", "verb2"),
					generateRoleRule("group1", "plural3/finalizer", "verb1"),
				},
				resourceList: []CrdResource{
					generateCRDResource("group1", "resource1", "plural1", "short1_", 3),
					generateCRDResource("group1", "resource2", "plural2", "short2_", 3),
					generateCRDResource("group1", "resource3", "plural3", "short3_", 3),
				},
			},
			wantMatching: []RoleRule{
				{Resource: RoleResource{Group: "group1", Name: "plural1"}, Verb: "verb1"},
				{Resource: RoleResource{Group: "group1", Name: "plural2"}, Verb: "verb1"},
				{Resource: RoleResource{Group: "group1", Name: "plural3/finalizer"}, Verb: "verb1"},
			},
			wantNonMatching: []RoleRule{
				{Resource: RoleResource{Group: "group2", Name: "plural2"}, Verb: "verb2"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMatching, gotNonMatching := FilterRulesNonMatchingResources(
				tt.args.ruleList,
				tt.args.resourceList,
			)
			if !reflect.DeepEqual(gotMatching, tt.wantMatching) {
				t.Errorf(
					"FilterRulesNonMatchingResources() gotMatching = %#v, want %#v",
					gotMatching,
					tt.wantMatching,
				)
			}
			if !reflect.DeepEqual(gotNonMatching, tt.wantNonMatching) {
				t.Errorf(
					"FilterRulesNonMatchingResources() gotNonMatching = %#v, want %#v",
					gotNonMatching,
					tt.wantNonMatching,
				)
			}
		})
	}
}
