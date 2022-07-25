// Copyright (C) 2020-2021 Red Hat, Inc.
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

package declaredandlistening

import (
	"context"
	"fmt"

	"github.com/onsi/ginkgo/v2"
	"github.com/sirupsen/logrus"
	"github.com/test-network-function/cnf-certification-test/cnf-certification-test/common"
	"github.com/test-network-function/cnf-certification-test/cnf-certification-test/identifiers"
	"github.com/test-network-function/cnf-certification-test/cnf-certification-test/networking/declaredandlistening"
	"github.com/test-network-function/cnf-certification-test/cnf-certification-test/networking/icmp"
	"github.com/test-network-function/cnf-certification-test/cnf-certification-test/networking/netcommons"
	"github.com/test-network-function/cnf-certification-test/cnf-certification-test/results"
	"github.com/test-network-function/cnf-certification-test/internal/clientsholder"
	"github.com/test-network-function/cnf-certification-test/internal/crclient"
	"github.com/test-network-function/cnf-certification-test/pkg/provider"
	"github.com/test-network-function/cnf-certification-test/pkg/testhelper"
	"github.com/test-network-function/cnf-certification-test/pkg/tnf"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	defaultNumPings = 5
	cmd             = `ss -tulwnH`
	nodePort        = "NodePort"
)

type Port []struct {
	ContainerPort int
	Name          string
	Protocol      string
}

//
// All actual test code belongs below here.  Utilities belong above.
//
var _ = ginkgo.Describe(common.NetworkingTestKey, func() {
	logrus.Debugf("Entering %s suite", common.NetworkingTestKey)

	var env provider.TestEnvironment
	ginkgo.BeforeEach(func() {
		env = provider.GetTestEnvironment()
	})
	ginkgo.ReportAfterEach(results.RecordResult)
	// Default interface ICMP IPv4 test case
	testID := identifiers.XformToGinkgoItIdentifier(identifiers.TestICMPv4ConnectivityIdentifier)
	ginkgo.It(testID, ginkgo.Label(testID), func() {
		testhelper.SkipIfEmptyAny(ginkgo.Skip, env.Containers, env.Pods)
		testNetworkConnectivity(&env, defaultNumPings, netcommons.IPv4, netcommons.DEFAULT)
	})
	// Multus interfaces ICMP IPv4 test case
	testID = identifiers.XformToGinkgoItIdentifier(identifiers.TestICMPv4ConnectivityMultusIdentifier)
	ginkgo.It(testID, ginkgo.Label(testID), func() {
		testhelper.SkipIfEmptyAny(ginkgo.Skip, env.Containers, env.Pods)
		testNetworkConnectivity(&env, defaultNumPings, netcommons.IPv4, netcommons.MULTUS)
	})
	// Default interface ICMP IPv6 test case
	testID = identifiers.XformToGinkgoItIdentifier(identifiers.TestICMPv6ConnectivityIdentifier)
	ginkgo.It(testID, ginkgo.Label(testID), func() {
		testhelper.SkipIfEmptyAny(ginkgo.Skip, env.Containers, env.Pods)
		testNetworkConnectivity(&env, defaultNumPings, netcommons.IPv6, netcommons.DEFAULT)
	})
	// Multus interfaces ICMP IPv6 test case
	testID = identifiers.XformToGinkgoItIdentifier(identifiers.TestICMPv6ConnectivityMultusIdentifier)
	ginkgo.It(testID, ginkgo.Label(testID), func() {
		testhelper.SkipIfEmptyAny(ginkgo.Skip, env.Containers, env.Pods)
		testNetworkConnectivity(&env, defaultNumPings, netcommons.IPv6, netcommons.MULTUS)
	})
	// Default interface ICMP IPv6 test case
	testID = identifiers.XformToGinkgoItIdentifier(identifiers.TestUndeclaredContainerPortsUsage)
	ginkgo.It(testID, ginkgo.Label(testID), func() {
		testhelper.SkipIfEmptyAny(ginkgo.Skip, env.Containers, env.Pods)
		testListenAndDeclared(&env)
	})
	testID = identifiers.XformToGinkgoItIdentifier(identifiers.TestServicesDoNotUseNodeportsIdentifier)
	ginkgo.It(testID, ginkgo.Label(testID), func() {
		testhelper.SkipIfEmptyAny(ginkgo.Skip, env.Containers, env.Pods)
		testNodePort(&env)
	})
})

//nolint:funlen
func testListenAndDeclared(env *provider.TestEnvironment) {
	var k declaredandlistening.Key
	var failedPods []*provider.Pod
	for _, podUnderTest := range env.Pods {
		declaredPorts := make(map[declaredandlistening.Key]bool)
		listeningPorts := make(map[declaredandlistening.Key]bool)
		for _, cut := range podUnderTest.Containers {
			ports := cut.Data.Ports
			logrus.Debugf("%s declaredPorts: %v", podUnderTest, ports)
			for j := 0; j < len(ports); j++ {
				k.Port = int(ports[j].ContainerPort)
				k.Protocol = string(ports[j].Protocol)
				declaredPorts[k] = true
			}
		}
		firstPodContainer := podUnderTest.Containers[0]
		outStr, errStr, err := crclient.ExecCommandContainerNSEnter(cmd, firstPodContainer)
		if err != nil || errStr != "" {
			tnf.ClaimFilePrintf("Failed to execute command %s on %s, err: %s, errStr: %s", cmd, firstPodContainer, err, errStr)
			failedPods = append(failedPods, podUnderTest)
			continue
		}
		declaredandlistening.ParseListening(outStr, listeningPorts)
		if len(listeningPorts) == 0 {
			tnf.ClaimFilePrintf("None of the containers of %s have any listening port.", podUnderTest)
			continue
		}
		// compare between declaredPort,listeningPort
		undeclaredPorts := declaredandlistening.CheckIfListenIsDeclared(listeningPorts, declaredPorts)
		for k := range undeclaredPorts {
			tnf.ClaimFilePrintf("%s is listening on port %d protocol %d, but that port was not declared in any container spec.", podUnderTest, k.Port, k.Protocol)
		}
		if len(undeclaredPorts) != 0 {
			failedPods = append(failedPods, podUnderTest)
		}
	}
	if nf := len(failedPods); nf > 0 {
		ginkgo.Fail(fmt.Sprintf("Found %d pods with listening ports not declared", nf))
	}
}

func testNodePort(env *provider.TestEnvironment) {
	badNamespaces := []string{}
	badServices := []string{}
	client := clientsholder.GetClientsHolder()
	for _, ns := range env.Namespaces {
		ginkgo.By(fmt.Sprintf("Testing services in namespace %s", ns))
		services, err := client.K8sClient.CoreV1().Services(ns).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			tnf.ClaimFilePrintf("Failed to list services on namespace %s, Error: %v", ns, err)
			badNamespaces = append(badNamespaces, ns)
			continue
		}
		for i := range services.Items {
			service := &services.Items[i]
			if service.Spec.Type == nodePort {
				tnf.ClaimFilePrintf("FAILURE: Service %s (ns %s) type is nodePort", service.Name, service.Namespace)
				badServices = append(badServices, fmt.Sprintf("ns: %s, name: %s", service.Namespace, service.Name))
			}
		}
	}
	if ns, bs := len(badNamespaces), len(badServices); ns > 0 || bs > 0 {
		ginkgo.Fail(fmt.Sprintf("Failed to get services on %d namespaces. %d services found of type nodePort.", ns, bs))
	}
}

// testDefaultNetworkConnectivity test the connectivity between the default interfaces of containers under test
func testNetworkConnectivity(env *provider.TestEnvironment, count int, aIPVersion netcommons.IPVersion, aType netcommons.IFType) {
	netsUnderTest, claimsLog := icmp.BuildNetTestContext(env.Pods, aIPVersion, aType)
	// Saving  curated logs to claims file
	tnf.ClaimFilePrintf("%s", claimsLog.GetLogLines())
	badNets, claimsLog, skip := icmp.RunNetworkingTests(netsUnderTest, count, aIPVersion)
	// Saving curated logs to claims file
	tnf.ClaimFilePrintf("%s", claimsLog.GetLogLines())
	if skip {
		ginkgo.Skip(fmt.Sprintf("There are no %s networks to test, skipping test", aIPVersion))
	}
	if n := len(badNets); n > 0 {
		logrus.Debugf("Failed nets: %+v", badNets)
		ginkgo.Fail(fmt.Sprintf("%d nets failed the %s network %s ping test.", n, aType, aIPVersion))
	}
}
