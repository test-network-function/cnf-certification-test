// Copyright (C) 2022 Red Hat, Inc.
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

package tolerations

import (
	v1 "k8s.io/api/core/v1"
)

var (
	nonCompliantTolerations  = []v1.TaintEffect{v1.TaintEffectNoExecute, v1.TaintEffectNoSchedule, v1.TaintEffectPreferNoSchedule}
	tolerationSecondsDefault = 300
)

func IsTolerationModified(t v1.Toleration) bool {
	const (
		notReadyStr    = "node.kubernetes.io/not-ready"
		unreachableStr = "node.kubernetes.io/unreachable"
	)
	// Check each of the tolerations to make sure they are the default tolerations added by k8s:
	// tolerations:
	// - effect: NoExecute
	//   key: node.kubernetes.io/not-ready
	//   operator: Exists
	//   tolerationSeconds: 300
	// - effect: NoExecute
	//   key: node.kubernetes.io/unreachable
	//   operator: Exists
	//   tolerationSeconds: 300

	if t.Effect == v1.TaintEffectNoExecute {
		if t.Key == notReadyStr || t.Key == unreachableStr &&
			(t.Operator == v1.TolerationOpExists && t.TolerationSeconds != nil && *t.TolerationSeconds == int64(tolerationSecondsDefault)) {
			return false
		}
	}

	// Check through the list of non-compliant tolerations to see if anything snuck by the above short circuit
	for _, nct := range nonCompliantTolerations {
		if t.Effect == nct {
			return true
		}
	}

	return false
}
