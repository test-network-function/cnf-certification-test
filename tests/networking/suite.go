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

package networking

import (
	"fmt"
	"strconv"

	"github.com/test-network-function/cnf-certification-test/internal/log"
	"github.com/test-network-function/cnf-certification-test/pkg/checksdb"
	"github.com/test-network-function/cnf-certification-test/pkg/provider"
	"github.com/test-network-function/cnf-certification-test/pkg/testhelper"
	"github.com/test-network-function/cnf-certification-test/tests/common"
	"github.com/test-network-function/cnf-certification-test/tests/identifiers"
	"github.com/test-network-function/cnf-certification-test/tests/networking/icmp"
	"github.com/test-network-function/cnf-certification-test/tests/networking/netcommons"
	"github.com/test-network-function/cnf-certification-test/tests/networking/netutil"
	"github.com/test-network-function/cnf-certification-test/tests/networking/policies"
	"github.com/test-network-function/cnf-certification-test/tests/networking/services"
	networkingv1 "k8s.io/api/networking/v1"
)

const (
	defaultNumPings = 5
	nodePort        = "NodePort"
)

type Port []struct {
	ContainerPort int
	Name          string
	Protocol      string
}

var (
	env provider.TestEnvironment

	beforeEachFn = func(check *checksdb.Check) error {
		env = provider.GetTestEnvironment()
		return nil
	}
)

//nolint:funlen
func LoadChecks() {
	log.Debug("Loading %s suite checks", common.NetworkingTestKey)

	checksGroup := checksdb.NewChecksGroup(common.NetworkingTestKey).
		WithBeforeEachFn(beforeEachFn)

	// Default interface ICMP IPv4 test case
	checksGroup.Add(checksdb.NewCheck(identifiers.GetTestIDAndLabels(identifiers.TestICMPv4ConnectivityIdentifier)).
		WithSkipCheckFn(testhelper.GetNoContainersUnderTestSkipFn(&env), testhelper.GetDaemonSetFailedToSpawnSkipFn(&env), testhelper.GetNoPodsUnderTestSkipFn(&env)).
		WithCheckFn(func(c *checksdb.Check) error {
			testNetworkConnectivity(&env, netcommons.IPv4, netcommons.DEFAULT, c)
			return nil
		}))

	// Multus interfaces ICMP IPv4 test case
	checksGroup.Add(checksdb.NewCheck(identifiers.GetTestIDAndLabels(identifiers.TestICMPv4ConnectivityMultusIdentifier)).
		WithSkipCheckFn(testhelper.GetNoContainersUnderTestSkipFn(&env), testhelper.GetDaemonSetFailedToSpawnSkipFn(&env), testhelper.GetNoPodsUnderTestSkipFn(&env)).
		WithCheckFn(func(c *checksdb.Check) error {
			testNetworkConnectivity(&env, netcommons.IPv4, netcommons.MULTUS, c)
			return nil
		}))

	// Default interface ICMP IPv6 test case
	checksGroup.Add(checksdb.NewCheck(identifiers.GetTestIDAndLabels(identifiers.TestICMPv6ConnectivityIdentifier)).
		WithSkipCheckFn(testhelper.GetNoContainersUnderTestSkipFn(&env), testhelper.GetDaemonSetFailedToSpawnSkipFn(&env), testhelper.GetNoPodsUnderTestSkipFn(&env)).
		WithCheckFn(func(c *checksdb.Check) error {
			testNetworkConnectivity(&env, netcommons.IPv6, netcommons.DEFAULT, c)
			return nil
		}))

	// Multus interfaces ICMP IPv6 test case
	checksGroup.Add(checksdb.NewCheck(identifiers.GetTestIDAndLabels(identifiers.TestICMPv6ConnectivityMultusIdentifier)).
		WithSkipCheckFn(testhelper.GetNoContainersUnderTestSkipFn(&env), testhelper.GetDaemonSetFailedToSpawnSkipFn(&env), testhelper.GetNoPodsUnderTestSkipFn(&env)).
		WithCheckFn(func(c *checksdb.Check) error {
			testNetworkConnectivity(&env, netcommons.IPv6, netcommons.MULTUS, c)
			return nil
		}))

	// Undeclared container ports usage test case
	checksGroup.Add(checksdb.NewCheck(identifiers.GetTestIDAndLabels(identifiers.TestUndeclaredContainerPortsUsage)).
		WithSkipCheckFn(testhelper.GetNoContainersUnderTestSkipFn(&env), testhelper.GetDaemonSetFailedToSpawnSkipFn(&env), testhelper.GetNoPodsUnderTestSkipFn(&env)).
		WithCheckFn(func(c *checksdb.Check) error {
			testUndeclaredContainerPortsUsage(c, &env)
			return nil
		}))

	// OCP reserved ports usage test case
	checksGroup.Add(checksdb.NewCheck(identifiers.GetTestIDAndLabels(identifiers.TestOCPReservedPortsUsage)).
		WithSkipCheckFn(testhelper.GetNoContainersUnderTestSkipFn(&env), testhelper.GetDaemonSetFailedToSpawnSkipFn(&env), testhelper.GetNoPodsUnderTestSkipFn(&env)).
		WithCheckFn(func(c *checksdb.Check) error {
			testOCPReservedPortsUsage(c, &env)
			return nil
		}))

	// Dual stack services test case
	checksGroup.Add(checksdb.NewCheck(identifiers.GetTestIDAndLabels(identifiers.TestServiceDualStackIdentifier)).
		WithSkipCheckFn(testhelper.GetNoServicesUnderTestSkipFn(&env)).
		WithCheckFn(func(c *checksdb.Check) error {
			testDualStackServices(c, &env)
			return nil
		}))

	// Network policy deny all test case
	checksGroup.Add(checksdb.NewCheck(identifiers.GetTestIDAndLabels(identifiers.TestNetworkPolicyDenyAllIdentifier)).
		WithSkipCheckFn(testhelper.GetNoPodsUnderTestSkipFn(&env)).
		WithCheckFn(func(c *checksdb.Check) error {
			testNetworkPolicyDenyAll(c, &env)
			return nil
		}))

	// Extended partner ports test case
	checksGroup.Add(checksdb.NewCheck(identifiers.GetTestIDAndLabels(identifiers.TestReservedExtendedPartnerPorts)).
		WithSkipCheckFn(testhelper.GetNoPodsUnderTestSkipFn(&env), testhelper.GetDaemonSetFailedToSpawnSkipFn(&env)).
		WithCheckFn(func(c *checksdb.Check) error {
			testPartnerSpecificTCPPorts(c, &env)
			return nil
		}))

	// DPDK CPU pinning exec probe test case
	checksGroup.Add(checksdb.NewCheck(identifiers.GetTestIDAndLabels(identifiers.TestDpdkCPUPinningExecProbe)).
		WithSkipCheckFn(testhelper.GetNoCPUPinningPodsSkipFn(&env)).
		WithCheckFn(func(c *checksdb.Check) error {
			dpdkPods := env.GetCPUPinningPodsWithDpdk()
			testExecProbDenyAtCPUPinning(c, dpdkPods)
			return nil
		}))

	// Restart on reboot label test case
	checksGroup.Add(checksdb.NewCheck(identifiers.GetTestIDAndLabels(identifiers.TestRestartOnRebootLabelOnPodsUsingSRIOV)).
		WithSkipCheckFn(testhelper.GetNoSRIOVPodsSkipFn(&env)).
		WithCheckFn(func(c *checksdb.Check) error {
			sriovPods, err := env.GetPodsUsingSRIOV()
			if err != nil {
				return fmt.Errorf("failure getting pods using SRIOV: %v", err)
			}
			testRestartOnRebootLabelOnPodsUsingSriov(c, sriovPods)
			return nil
		}))
}

func testExecProbDenyAtCPUPinning(check *checksdb.Check, dpdkPods []*provider.Pod) {
	var compliantObjects []*testhelper.ReportObject
	var nonCompliantObjects []*testhelper.ReportObject

	for _, cpuPinnedPod := range dpdkPods {
		execProbeFound := false
		for _, cut := range cpuPinnedPod.Containers {
			check.LogInfo("Testing Container %q", cut)
			if cut.HasExecProbes() {
				check.LogError("Container %q defines an exec probe", cut)
				nonCompliantObjects = append(nonCompliantObjects, testhelper.NewPodReportObject(cpuPinnedPod.Namespace, cpuPinnedPod.Name, "Exec prob is not allowed", false))
				execProbeFound = true
			}
		}

		if !execProbeFound {
			check.LogInfo("Pod %q does not define any exec probe", cpuPinnedPod)
			compliantObjects = append(compliantObjects, testhelper.NewPodReportObject(cpuPinnedPod.Namespace, cpuPinnedPod.Name, "Exec prob is allowed", true))
		}
	}
	check.SetResult(compliantObjects, nonCompliantObjects)
}

//nolint:funlen
func testUndeclaredContainerPortsUsage(check *checksdb.Check, env *provider.TestEnvironment) {
	var compliantObjects []*testhelper.ReportObject
	var nonCompliantObjects []*testhelper.ReportObject
	var portInfo netutil.PortInfo
	for _, put := range env.Pods {
		// First get the ports declared in the Pod's containers spec
		declaredPorts := make(map[netutil.PortInfo]bool)
		for _, cut := range put.Containers {
			check.LogInfo("Testing Container %q", cut)
			for _, port := range cut.Ports {
				portInfo.PortNumber = port.ContainerPort
				portInfo.Protocol = string(port.Protocol)
				declaredPorts[portInfo] = true
			}
		}

		// Then check the actual ports that the containers are listening on
		firstPodContainer := put.Containers[0]
		listeningPorts, err := netutil.GetListeningPorts(firstPodContainer)
		if err != nil {
			check.LogError("Failed to get container %q listening ports, err: %v", firstPodContainer, err)
			nonCompliantObjects = append(nonCompliantObjects,
				testhelper.NewPodReportObject(put.Namespace, put.Name, fmt.Sprintf("Failed to get the container's listening ports, err: %v", err), false))
			continue
		}
		if len(listeningPorts) == 0 {
			check.LogInfo("None of the containers of %q have any listening port.", put)
			continue
		}

		// Verify that all the listening ports have been declared in the container spec
		failedPod := false
		for listeningPort := range listeningPorts {
			if put.ContainsIstioProxy() && netcommons.ReservedIstioPorts[listeningPort.PortNumber] {
				check.LogInfo("%q is listening on port %d protocol %q, but the pod also contains istio-proxy. Ignoring.",
					put, listeningPort.PortNumber, listeningPort.Protocol)
				continue
			}

			if ok := declaredPorts[listeningPort]; !ok {
				check.LogError("%q is listening on port %d protocol %q, but that port was not declared in any container spec.",
					put, listeningPort.PortNumber, listeningPort.Protocol)
				failedPod = true
				nonCompliantObjects = append(nonCompliantObjects,
					testhelper.NewPodReportObject(put.Namespace, put.Name,
						"Listening port was declared in no container spec", false).
						SetType(testhelper.ListeningPortType).
						AddField(testhelper.PortNumber, strconv.Itoa(int(listeningPort.PortNumber))).
						AddField(testhelper.PortProtocol, listeningPort.Protocol))
			} else {
				check.LogInfo("%q is listening on declared port %d protocol %q", put, listeningPort.PortNumber, listeningPort.Protocol)
				compliantObjects = append(compliantObjects,
					testhelper.NewPodReportObject(put.Namespace, put.Name,
						"Listening port was declared in container spec", true).
						SetType(testhelper.ListeningPortType).
						AddField(testhelper.PortNumber, strconv.Itoa(int(listeningPort.PortNumber))).
						AddField(testhelper.PortProtocol, listeningPort.Protocol))
			}
		}
		if failedPod {
			nonCompliantObjects = append(nonCompliantObjects,
				testhelper.NewPodReportObject(put.Namespace, put.Name, "At least one port was listening but not declared in any container specs", false))
		} else {
			compliantObjects = append(compliantObjects,
				testhelper.NewPodReportObject(put.Namespace, put.Name, "All listening were declared in containers specs", true))
		}
	}
	check.SetResult(compliantObjects, nonCompliantObjects)
}

// testDefaultNetworkConnectivity test the connectivity between the default interfaces of containers under test.
func testNetworkConnectivity(env *provider.TestEnvironment, aIPVersion netcommons.IPVersion, aType netcommons.IFType, check *checksdb.Check) {
	netsUnderTest := icmp.BuildNetTestContext(env.Pods, aIPVersion, aType, check.GetLogger())
	report, skip := icmp.RunNetworkingTests(netsUnderTest, defaultNumPings, aIPVersion, check.GetLogger())
	if skip {
		check.LogInfo("There are no %q networks to test with at least 2 pods, skipping test", aIPVersion)
	}
	check.SetResult(report.CompliantObjectsOut, report.NonCompliantObjectsOut)
}

func testOCPReservedPortsUsage(check *checksdb.Check, env *provider.TestEnvironment) {
	// List of all ports reserved by OpenShift.
	OCPReservedPorts := map[int32]bool{
		22623: true,
		22624: true}
	compliantObjects, nonCompliantObjects := netcommons.TestReservedPortsUsage(env, OCPReservedPorts, "OCP", check.GetLogger())
	check.SetResult(compliantObjects, nonCompliantObjects)
}

func testPartnerSpecificTCPPorts(check *checksdb.Check, env *provider.TestEnvironment) {
	// List of all of the ports reserved by partner.
	ReservedPorts := map[int32]bool{
		15443: true,
		15090: true,
		15021: true,
		15020: true,
		15014: true,
		15008: true,
		15006: true,
		15001: true,
		15000: true,
	}
	compliantObjects, nonCompliantObjects := netcommons.TestReservedPortsUsage(env, ReservedPorts, "Partner", check.GetLogger())
	check.SetResult(compliantObjects, nonCompliantObjects)
}

func testDualStackServices(check *checksdb.Check, env *provider.TestEnvironment) {
	var compliantObjects []*testhelper.ReportObject
	var nonCompliantObjects []*testhelper.ReportObject
	for _, s := range env.Services {
		check.LogInfo("Testing Service %q", s.Name)
		serviceIPVersion, err := services.GetServiceIPVersion(s)
		if err != nil {
			check.LogError("Could not get IP version from Service %q, err=%v", s.Name, err)
			nonCompliantObjects = append(nonCompliantObjects, testhelper.NewReportObject("Could not get IP Version from service", testhelper.ServiceType, false).
				AddField(testhelper.Namespace, s.Namespace).
				AddField(testhelper.ServiceName, s.Name))
		}
		if serviceIPVersion == netcommons.Undefined || serviceIPVersion == netcommons.IPv4 {
			check.LogError("Service %q (ns: %q) only supports IPv4", s.Name, s.Namespace)
			nonCompliantObjects = append(nonCompliantObjects, testhelper.NewReportObject("Service supports only IPv4", testhelper.ServiceType, false).
				AddField(testhelper.Namespace, s.Namespace).
				AddField(testhelper.ServiceName, s.Name).
				AddField(testhelper.ServiceIPVersion, serviceIPVersion.String()))
		} else {
			check.LogInfo("Service %q (ns: %q) supports IPv6 or is dual stack", s.Name, s.Namespace)
			compliantObjects = append(compliantObjects, testhelper.NewReportObject("Service supports IPv6 or is dual stack", testhelper.ServiceType, true).
				AddField(testhelper.Namespace, s.Namespace).
				AddField(testhelper.ServiceName, s.Name).
				AddField(testhelper.ServiceIPVersion, serviceIPVersion.String()))
		}
	}

	check.SetResult(compliantObjects, nonCompliantObjects)
}

//nolint:funlen
func testNetworkPolicyDenyAll(check *checksdb.Check, env *provider.TestEnvironment) {
	var compliantObjects []*testhelper.ReportObject
	var nonCompliantObjects []*testhelper.ReportObject

	// Loop through the pods, looking for corresponding entries within a deny-all network policy (both ingress and egress).
	// This ensures that each pod is accounted for that we are tasked with testing and excludes any pods that are not marked
	// for testing (via the labels).
	for _, put := range env.Pods {
		check.LogInfo("Testing Pod %q", put)
		denyAllEgressFound := false
		denyAllIngressFound := false

		// Look through all of the network policies for a matching namespace.
		for index := range env.NetworkPolicies {
			networkPolicy := env.NetworkPolicies[index]
			check.LogInfo("Testing Network policy %q against pod %q", networkPolicy.Name, put)

			// Skip any network policies that don't match the namespace of the pod we are testing.
			if networkPolicy.Namespace != put.Namespace {
				check.LogInfo("Skipping Network policy %q (namespace %q does not match Pod namespace %q)", networkPolicy.Name, networkPolicy.Namespace, put.Namespace)
				continue
			}

			// Match the pod namespace with the network policy namespace.
			if policies.LabelsMatch(networkPolicy.Spec.PodSelector, put.Labels) {
				var reason string
				if !denyAllEgressFound {
					denyAllEgressFound, reason = policies.IsNetworkPolicyCompliant(&networkPolicy, networkingv1.PolicyTypeEgress)
					if reason != "" {
						check.LogError("Network policy %q is not compliant, reason=%q", networkPolicy.Name, reason)
					}
				}
				if !denyAllIngressFound {
					denyAllIngressFound, reason = policies.IsNetworkPolicyCompliant(&networkPolicy, networkingv1.PolicyTypeIngress)
					if reason != "" {
						check.LogError("Network policy %q is not compliant, reason=%q", networkPolicy.Name, reason)
					}
				}
			}
		}

		// Network policy has not been found that contains a deny-all rule for both ingress and egress.
		podIsCompliant := true
		if !denyAllIngressFound {
			check.LogError("Pod %q was found to not have a default ingress deny-all network policy.", put)
			nonCompliantObjects = append(nonCompliantObjects, testhelper.NewPodReportObject(put.Namespace, put.Name, "Pod was found to not have a default ingress deny-all network policy", false))
			podIsCompliant = false
		}

		if !denyAllEgressFound {
			check.LogError("Pod %q was found to not have a default egress deny-all network policy.", put)
			nonCompliantObjects = append(nonCompliantObjects, testhelper.NewPodReportObject(put.Namespace, put.Name, "Pod was found to not have a default egress deny-all network policy", false))
			podIsCompliant = false
		}

		if podIsCompliant {
			check.LogInfo("Pod %q has a default ingress/egress deny-all network policy", put)
			compliantObjects = append(compliantObjects, testhelper.NewPodReportObject(put.Namespace, put.Name, "Pod has a default ingress/egress deny-all network policy", true))
		}
	}

	check.SetResult(compliantObjects, nonCompliantObjects)
}

func testRestartOnRebootLabelOnPodsUsingSriov(check *checksdb.Check, sriovPods []*provider.Pod) {
	const (
		restartOnRebootLabel = "restart-on-reboot"
	)

	var compliantObjects []*testhelper.ReportObject
	var nonCompliantObjects []*testhelper.ReportObject
	for _, pod := range sriovPods {
		check.LogInfo("Testing SRIOV Pod %q", pod)

		labelValue, exist := pod.GetLabels()[restartOnRebootLabel]
		if !exist {
			check.LogError("Pod %q uses SRIOV but the label %q was not found.", pod, restartOnRebootLabel)
			nonCompliantObjects = append(nonCompliantObjects, testhelper.NewPodReportObject(pod.Namespace, pod.Name, fmt.Sprintf("Pod uses SRIOV but the label %s was not found", restartOnRebootLabel), false))
			continue
		}

		if labelValue != "true" {
			check.LogError("Pod %q uses SRIOV but the %q label value is not true.", pod, restartOnRebootLabel)
			nonCompliantObjects = append(nonCompliantObjects, testhelper.NewPodReportObject(pod.Namespace, pod.Name, fmt.Sprintf("Pod uses SRIOV but the label %s is not set to true", restartOnRebootLabel), false))
			continue
		}

		check.LogInfo("Pod %q uses SRIOV and the %q label is set to true", pod, restartOnRebootLabel)
		compliantObjects = append(compliantObjects, testhelper.NewPodReportObject(pod.Namespace, pod.Name, fmt.Sprintf("Pod uses SRIOV and the label %s is set to true", restartOnRebootLabel), true))
	}

	check.SetResult(compliantObjects, nonCompliantObjects)
}
