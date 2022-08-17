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

package isolation

import (
	"github.com/sirupsen/logrus"
	"github.com/test-network-function/cnf-certification-test/pkg/provider"
)

type TesterFuncs interface {
	AreResourcesIdentical() bool
	AreCPUResourcesWholeUnits() bool
	IsRuntimeClassNameSpecified() bool
	LoadBalancingDisabled() bool
}

type CPUIsolationTester struct {
	pod *provider.Pod
}

func NewCPUIsolationTester(p *provider.Pod) *CPUIsolationTester {
	return &CPUIsolationTester{
		pod: p,
	}
}

func (it *CPUIsolationTester) AreResourcesIdentical() bool {
	// Pods may contain more than one container.  All containers must conform to the CPU isolation requirements.
	for _, cut := range it.pod.Containers {
		// Resources must be specified
		if len(cut.Data.Resources.Requests) == 0 || len(cut.Data.Resources.Limits) == 0 {
			logrus.Debugf("%s has been found with undefined requests or limits.", cut.String())
			return false
		}

		// Gather the values
		cpuRequests := cut.Data.Resources.Requests.Cpu()
		cpuLimits := cut.Data.Resources.Limits.Cpu()
		memoryRequests := cut.Data.Resources.Requests.Memory()
		memoryLimits := cut.Data.Resources.Limits.Memory()

		// Check for mismatches
		if cpuRequests.Value() != cpuLimits.Value() {
			logrus.Debugf("%s has CPU requests %d and limits %d that do not match.", cut.String(), cpuRequests.Value(), cpuLimits.Value())
			return false
		}

		if memoryRequests.Value() != memoryLimits.Value() {
			logrus.Debugf("%s has memory requests %d and limits %d that do not match.", cut.String(), memoryRequests.Value(), memoryLimits.Value())
			return false
		}
	}

	return true
}

func (it *CPUIsolationTester) AreCPUResourcesWholeUnits() bool {
	isInteger := func(val float64) bool {
		return val == float64(int(val))
	}

	// Pods may contain more than one container.  All containers must conform to the CPU isolation requirements.
	for _, cut := range it.pod.Containers {
		// Resources must be specified
		if len(cut.Data.Resources.Requests) == 0 || len(cut.Data.Resources.Limits) == 0 {
			logrus.Debugf("%s has been found with undefined requests or limits.", cut.String())
			return false
		}

		// Gather the values
		cpuRequests := cut.Data.Resources.Requests.Cpu().AsApproximateFloat64()
		cpuLimits := cut.Data.Resources.Limits.Cpu().AsApproximateFloat64()

		if !isInteger(cpuRequests) {
			logrus.Debugf("%s has CPU requests %f that has to be a whole unit.", cut.String(), cpuRequests)
			return false
		}
		if !isInteger(cpuLimits) {
			logrus.Debugf("%s has CPU limits %f that has to be a whole unit.", cut.String(), cpuLimits)
			return false
		}
	}

	return true
}

func (it *CPUIsolationTester) IsRuntimeClassNameSpecified() bool {
	return it.pod.Data.Spec.RuntimeClassName != nil
}

func (it *CPUIsolationTester) LoadBalancingDisabled() bool {
	const (
		disableVar = "disable"
	)

	cpuLoadBalancingDisabled := false
	irqLoadBalancingDisabled := false

	if v, ok := it.pod.Data.ObjectMeta.Annotations["cpu-load-balancing.crio.io"]; ok {
		if v == disableVar {
			cpuLoadBalancingDisabled = true
		} else {
			logrus.Debugf("Annotation cpu-load-balancing.crio.io has a invalid value for CPU isolation.  Must be 'disable'.")
		}
	} else {
		logrus.Debugf("Annotation cpu-load-balancing.crio.io is missing.")
	}

	if v, ok := it.pod.Data.ObjectMeta.Annotations["irq-load-balancing.crio.io"]; ok {
		if v == disableVar {
			irqLoadBalancingDisabled = true
		} else {
			logrus.Debugf("Annotation irq-load-balancing.crio.io has a invalid value for CPU isolation.  Must be 'disable'.")
		}
	} else {
		logrus.Debugf("Annotation irq-load-balancing.crio.io is missing.")
	}

	// Both conditions have to be set to 'disable'
	if cpuLoadBalancingDisabled && irqLoadBalancingDisabled {
		return true
	}

	return false
}
