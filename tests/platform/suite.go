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

package platform

import (
	"fmt"
	"strconv"
	"strings"

	clientsholder "github.com/test-network-function/cnf-certification-test/internal/clientsholder"
	"github.com/test-network-function/cnf-certification-test/internal/log"
	"github.com/test-network-function/cnf-certification-test/pkg/checksdb"
	"github.com/test-network-function/cnf-certification-test/pkg/compatibility"
	"github.com/test-network-function/cnf-certification-test/pkg/provider"
	"github.com/test-network-function/cnf-certification-test/pkg/testhelper"
	"github.com/test-network-function/cnf-certification-test/tests/common"
	"github.com/test-network-function/cnf-certification-test/tests/identifiers"
	"github.com/test-network-function/cnf-certification-test/tests/platform/bootparams"
	"github.com/test-network-function/cnf-certification-test/tests/platform/cnffsdiff"
	"github.com/test-network-function/cnf-certification-test/tests/platform/hugepages"
	"github.com/test-network-function/cnf-certification-test/tests/platform/isredhat"

	"github.com/test-network-function/cnf-certification-test/tests/platform/operatingsystem"
	"github.com/test-network-function/cnf-certification-test/tests/platform/sysctlconfig"

	"github.com/test-network-function/cnf-certification-test/tests/platform/nodetainted"
)

var (
	env provider.TestEnvironment

	beforeEachFn = func(check *checksdb.Check) error {
		env = provider.GetTestEnvironment()
		return nil
	}
)

//nolint:funlen
func LoadChecks() {
	log.Debug("Loading %s suite checks", common.PlatformAlterationTestKey)

	checksGroup := checksdb.NewChecksGroup(common.PlatformAlterationTestKey).
		WithBeforeEachFn(beforeEachFn)

	checksGroup.Add(checksdb.NewCheck(identifiers.GetTestIDAndLabels(identifiers.TestHyperThreadEnable)).
		WithSkipCheckFn(testhelper.GetNoBareMetalNodesSkipFn(&env)).
		WithCheckFn(func(c *checksdb.Check) error {
			testHyperThreadingEnabled(c, &env)
			return nil
		}))

	checksGroup.Add(checksdb.NewCheck(identifiers.GetTestIDAndLabels(identifiers.TestUnalteredBaseImageIdentifier)).
		WithSkipCheckFn(
			testhelper.GetNonOCPClusterSkipFn(),
			testhelper.GetDaemonSetFailedToSpawnSkipFn(&env),
			testhelper.GetNoContainersUnderTestSkipFn(&env)).
		WithCheckFn(func(c *checksdb.Check) error {
			testContainersFsDiff(c, &env)
			return nil
		}))

	checksGroup.Add(checksdb.NewCheck(identifiers.GetTestIDAndLabels(identifiers.TestNonTaintedNodeKernelsIdentifier)).
		WithSkipCheckFn(testhelper.GetDaemonSetFailedToSpawnSkipFn(&env)).
		WithCheckFn(func(c *checksdb.Check) error {
			testTainted(c, &env)
			return nil
		}))

	checksGroup.Add(checksdb.NewCheck(identifiers.GetTestIDAndLabels(identifiers.TestIsRedHatReleaseIdentifier)).
		WithSkipCheckFn(testhelper.GetNoContainersUnderTestSkipFn(&env)).
		WithCheckFn(func(c *checksdb.Check) error {
			testIsRedHatRelease(c, &env)
			return nil
		}))

	checksGroup.Add(checksdb.NewCheck(identifiers.GetTestIDAndLabels(identifiers.TestIsSELinuxEnforcingIdentifier)).
		WithSkipCheckFn(
			testhelper.GetNonOCPClusterSkipFn(),
			testhelper.GetDaemonSetFailedToSpawnSkipFn(&env)).
		WithCheckFn(func(c *checksdb.Check) error {
			testIsSELinuxEnforcing(c, &env)
			return nil
		}))

	checksGroup.Add(checksdb.NewCheck(identifiers.GetTestIDAndLabels(identifiers.TestHugepagesNotManuallyManipulated)).
		WithSkipCheckFn(
			testhelper.GetNonOCPClusterSkipFn(),
			testhelper.GetDaemonSetFailedToSpawnSkipFn(&env)).
		WithCheckFn(func(c *checksdb.Check) error {
			testHugepages(c, &env)
			return nil
		}))

	checksGroup.Add(checksdb.NewCheck(identifiers.GetTestIDAndLabels(identifiers.TestUnalteredStartupBootParamsIdentifier)).
		WithSkipCheckFn(
			testhelper.GetNonOCPClusterSkipFn(),
			testhelper.GetDaemonSetFailedToSpawnSkipFn(&env)).
		WithCheckFn(func(c *checksdb.Check) error {
			testUnalteredBootParams(c, &env)
			return nil
		}))

	checksGroup.Add(checksdb.NewCheck(identifiers.GetTestIDAndLabels(identifiers.TestSysctlConfigsIdentifier)).
		WithSkipCheckFn(
			testhelper.GetNonOCPClusterSkipFn(),
			testhelper.GetDaemonSetFailedToSpawnSkipFn(&env)).
		WithCheckFn(func(c *checksdb.Check) error {
			testSysctlConfigs(c, &env)
			return nil
		}))

	checksGroup.Add(checksdb.NewCheck(identifiers.GetTestIDAndLabels(identifiers.TestServiceMeshIdentifier)).
		WithSkipCheckFn(
			testhelper.GetNoIstioSkipFn(&env),
			testhelper.GetNoPodsUnderTestSkipFn(&env)).
		WithCheckFn(func(c *checksdb.Check) error {
			testServiceMesh(c, &env)
			return nil
		}))

	checksGroup.Add(checksdb.NewCheck(identifiers.GetTestIDAndLabels(identifiers.TestOCPLifecycleIdentifier)).
		WithSkipCheckFn(testhelper.GetNonOCPClusterSkipFn()).
		WithCheckFn(func(c *checksdb.Check) error {
			testOCPStatus(c, &env)
			return nil
		}))

	checksGroup.Add(checksdb.NewCheck(identifiers.GetTestIDAndLabels(identifiers.TestNodeOperatingSystemIdentifier)).
		WithSkipCheckFn(testhelper.GetNonOCPClusterSkipFn()).
		WithCheckFn(func(c *checksdb.Check) error {
			testNodeOperatingSystemStatus(c, &env)
			return nil
		}))

	checksGroup.Add(checksdb.NewCheck(identifiers.GetTestIDAndLabels(identifiers.TestPodHugePages2M)).
		WithSkipCheckFn(
			testhelper.GetNonOCPClusterSkipFn(),
			testhelper.GetNoHugepagesPodsSkipFn(&env)).
		WithCheckFn(func(c *checksdb.Check) error {
			testPodHugePagesSize(c, &env, provider.HugePages2Mi)
			return nil
		}))

	checksGroup.Add(checksdb.NewCheck(identifiers.GetTestIDAndLabels(identifiers.TestPodHugePages1G)).
		WithSkipCheckFn(
			testhelper.GetNonOCPClusterSkipFn(),
			testhelper.GetNoHugepagesPodsSkipFn(&env)).
		WithCheckFn(func(c *checksdb.Check) error {
			testPodHugePagesSize(c, &env, provider.HugePages1Gi)
			return nil
		}))
}

func testHyperThreadingEnabled(check *checksdb.Check, env *provider.TestEnvironment) {
	var compliantObjects []*testhelper.ReportObject
	var nonCompliantObjects []*testhelper.ReportObject
	baremetalNodes := env.GetBaremetalNodes()
	for _, node := range baremetalNodes {
		nodeName := node.Data.Name
		check.LogInfo("Testing node %q", nodeName)
		enable, err := node.IsHyperThreadNode(env)
		//nolint:gocritic
		if enable {
			check.LogInfo("Node %q has hyperthreading enabled", nodeName)
			compliantObjects = append(compliantObjects, testhelper.NewNodeReportObject(nodeName, "Node has hyperthreading enabled", true))
		} else if err != nil {
			check.LogError("Hyperthreading check fail for node %q, err: %v", nodeName, err)
			nonCompliantObjects = append(nonCompliantObjects, testhelper.NewNodeReportObject(nodeName, "Error with executing the check for hyperthreading: "+err.Error(), false))
		} else {
			check.LogError("Node %q has hyperthreading disabled", nodeName)
			nonCompliantObjects = append(nonCompliantObjects, testhelper.NewNodeReportObject(nodeName, "Node has hyperthreading disabled ", false))
		}
	}
	check.SetResult(compliantObjects, nonCompliantObjects)
}

func testServiceMesh(check *checksdb.Check, env *provider.TestEnvironment) {
	var compliantObjects []*testhelper.ReportObject
	var nonCompliantObjects []*testhelper.ReportObject
	for _, put := range env.Pods {
		check.LogInfo("Testing Pod %q", put)
		istioProxyFound := false
		for _, cut := range put.Containers {
			if cut.IsIstioProxy() {
				check.LogInfo("Istio proxy container found on Pod %q (Container %q)", put, cut)
				istioProxyFound = true
				break
			}
		}
		if !istioProxyFound {
			check.LogError("Pod %q found without service mesh", put)
			nonCompliantObjects = append(nonCompliantObjects, testhelper.NewPodReportObject(put.Namespace, put.Name, "Pod found without service mesh container", false))
		} else {
			check.LogInfo("Pod %q found with service mesh", put)
			compliantObjects = append(compliantObjects, testhelper.NewPodReportObject(put.Namespace, put.Name, "Pod found with service mesh container", true))
		}
	}

	check.SetResult(compliantObjects, nonCompliantObjects)
}

// testContainersFsDiff test that all CUT did not install new packages are starting.
func testContainersFsDiff(check *checksdb.Check, env *provider.TestEnvironment) {
	var compliantObjects []*testhelper.ReportObject
	var nonCompliantObjects []*testhelper.ReportObject
	for _, cut := range env.Containers {
		check.LogInfo("Testing Container %q", cut)
		debugPod := env.DebugPods[cut.NodeName]

		ctxt := clientsholder.NewContext(debugPod.Namespace, debugPod.Name, debugPod.Spec.Containers[0].Name)
		fsDiffTester := cnffsdiff.NewFsDiffTester(check, clientsholder.GetClientsHolder(), ctxt, env.OpenshiftVersion)
		fsDiffTester.RunTest(cut.UID)
		switch fsDiffTester.GetResults() {
		case testhelper.SUCCESS:
			check.LogInfo("Container %q is not modified", cut)
			compliantObjects = append(compliantObjects, testhelper.NewContainerReportObject(cut.Namespace, cut.Podname, cut.Name, "Container is not modified", true))
			continue
		case testhelper.FAILURE:
			check.LogError("Container %q modified (changed folders: %v, deleted folders: %v", cut, fsDiffTester.ChangedFolders, fsDiffTester.DeletedFolders)
			nonCompliantObjects = append(nonCompliantObjects, testhelper.NewContainerReportObject(cut.Namespace, cut.Podname, cut.Name, "Container is modified", false).
				AddField("ChangedFolders", strings.Join(fsDiffTester.ChangedFolders, ",")).
				AddField("DeletedFolders", strings.Join(fsDiffTester.DeletedFolders, ",")))

		case testhelper.ERROR:
			check.LogError("Could not run fs-diff in Container %q, err: %v", cut, fsDiffTester.Error)
			nonCompliantObjects = append(nonCompliantObjects, testhelper.NewContainerReportObject(cut.Namespace, cut.Podname, cut.Name, "Error while running fs-diff", false).AddField(testhelper.Error, fsDiffTester.Error.Error()))
		}
	}

	check.SetResult(compliantObjects, nonCompliantObjects)
}

//nolint:funlen
func testTainted(check *checksdb.Check, env *provider.TestEnvironment) {
	var compliantObjects []*testhelper.ReportObject
	var nonCompliantObjects []*testhelper.ReportObject

	// errNodes has nodes that failed some operation while checking kernel taints.
	errNodes := []string{}

	type badModuleTaints struct {
		name   string
		taints []string
	}

	// badModules maps node names to list of "bad"/offending modules.
	badModules := map[string][]badModuleTaints{}
	// otherTaints maps a node to a list of taint bits that haven't been set by any module.
	otherTaints := map[string][]int{}

	check.LogInfo("Modules allowlist: %+v", env.Config.AcceptedKernelTaints)
	// helper map to make the checks easier.
	allowListedModules := map[string]bool{}
	for _, module := range env.Config.AcceptedKernelTaints {
		allowListedModules[module.Module] = true
	}

	// Loop through the debug pods that are tied to each node.
	for _, n := range env.Nodes {
		nodeName := n.Data.Name
		check.LogInfo("Testing node %q", nodeName)

		// Ensure we are only testing nodes that have CNF workload deployed on them.
		if !n.HasWorkloadDeployed(env.Pods) {
			check.LogInfo("Node %q has no workload deployed on it. Skipping tainted kernel check.", nodeName)
			continue
		}

		dp := env.DebugPods[nodeName]

		ocpContext := clientsholder.NewContext(dp.Namespace, dp.Name, dp.Spec.Containers[0].Name)
		tf := nodetainted.NewNodeTaintedTester(&ocpContext, nodeName)

		// Get the taints mask from the node kernel
		taintsMask, err := tf.GetKernelTaintsMask()
		if err != nil {
			check.LogError("Failed to retrieve kernel taint information from node %q, err: %v", nodeName, err)
			nonCompliantObjects = append(nonCompliantObjects, testhelper.NewNodeReportObject(nodeName, "Failed to retrieve kernel taint information from node", false).
				AddField(testhelper.Error, err.Error()))
			continue
		}

		if taintsMask == 0 {
			check.LogInfo("Node %q has no non-approved kernel taints.", nodeName)
			compliantObjects = append(compliantObjects, testhelper.NewNodeReportObject(nodeName, "Node has no non-approved kernel taints", true))
			continue
		}

		check.LogInfo("Node %q kernel is tainted. Taints mask=%d - Decoded taints: %v",
			nodeName, taintsMask, nodetainted.DecodeKernelTaintsFromBitMask(taintsMask))

		// Check the allow list. If empty, mark this node as failed.
		if len(allowListedModules) == 0 {
			taintsMaskStr := strconv.FormatUint(taintsMask, 10)
			taintsStr := strings.Join(nodetainted.DecodeKernelTaintsFromBitMask(taintsMask), ",")
			check.LogError("Node %q contains taints not covered by module allowlist. Taints: %q (mask=%q)", nodeName, taintsStr, taintsMaskStr)
			nonCompliantObjects = append(nonCompliantObjects, testhelper.NewNodeReportObject(nodeName, "Node contains taints not covered by module allowlist", false).
				AddField(testhelper.TaintMask, taintsMaskStr).
				AddField(testhelper.Taints, taintsStr))
			continue
		}

		// allow list check.
		// Get the list of modules (tainters) that have set a taint bit.
		//   1. Each module should appear in the allow list.
		//   2. All kernel taint bits (one bit <-> one letter) should have been set by at least
		//      one tainter module.
		tainters, taintBitsByAllModules, err := tf.GetTainterModules(allowListedModules)
		if err != nil {
			check.LogError("Could not get tainter modules from node %q, err: %v", nodeName, err)
			errNodes = append(errNodes, nodeName)
			nonCompliantObjects = append(nonCompliantObjects, testhelper.NewNodeReportObject(nodeName, "Failed to get tainter modules", false).
				AddField(testhelper.Error, err.Error()))
			continue
		}

		// Keep track of whether or not this node is compliant with module allow list.
		compliantNode := true

		// Save modules' names only.
		for moduleName, taintsLetters := range tainters {
			moduleTaints := nodetainted.DecodeKernelTaintsFromLetters(taintsLetters)
			badModules[nodeName] = append(badModules[nodeName], badModuleTaints{name: moduleName, taints: moduleTaints})

			// Create non-compliant taint objects for each of the taints
			for _, taint := range moduleTaints {
				check.LogError("Node %q - module %q taints kernel: %q", nodeName, moduleName, taint)
				nonCompliantObjects = append(nonCompliantObjects, testhelper.NewTaintReportObject(nodetainted.RemoveAllExceptNumbers(taint), nodeName, taint, false).AddField(testhelper.ModuleName, moduleName))

				// Set the node as non-compliant for future reporting
				compliantNode = false
			}
		}

		// Lastly, check that all kernel taint bits come from modules.
		otherKernelTaints := nodetainted.GetOtherTaintedBits(taintsMask, taintBitsByAllModules)
		for _, taintedBit := range otherKernelTaints {
			check.LogError("Node %q - taint bit %d is set but it is not caused by any module.", nodeName, taintedBit)
			nonCompliantObjects = append(nonCompliantObjects, testhelper.NewTaintReportObject(strconv.Itoa(taintedBit), nodeName, nodetainted.GetTaintMsg(taintedBit), false).
				AddField(testhelper.ModuleName, "N/A"))
			otherTaints[nodeName] = append(otherTaints[nodeName], taintedBit)

			// Set the node as non-compliant for future reporting
			compliantNode = false
		}

		if compliantNode {
			check.LogInfo("Node %q passed the tainted kernel check", nodeName)
			compliantObjects = append(compliantObjects, testhelper.NewNodeReportObject(nodeName, "Passed the tainted kernel check", true))
		}
	}

	if len(errNodes) > 0 {
		check.LogError("Failed to get kernel taints from some nodes: %+v", errNodes)
	}

	if len(badModules) > 0 || len(otherTaints) > 0 {
		check.LogError("Nodes have been found to be tainted. Tainted modules: %+v", badModules)
	}

	if len(otherTaints) > 0 {
		check.LogError("Taints not related to any module: %+v", otherTaints)
	}

	check.SetResult(compliantObjects, nonCompliantObjects)
}

func testIsRedHatRelease(check *checksdb.Check, env *provider.TestEnvironment) {
	var compliantObjects []*testhelper.ReportObject
	var nonCompliantObjects []*testhelper.ReportObject
	for _, cut := range env.Containers {
		check.LogInfo("Testing Container %q", cut)
		baseImageTester := isredhat.NewBaseImageTester(clientsholder.GetClientsHolder(), clientsholder.NewContext(cut.Namespace, cut.Podname, cut.Name))

		result, err := baseImageTester.TestContainerIsRedHatRelease()
		if err != nil {
			check.LogError("Could not collect release information from Container %q, err=%v", cut, err)
		}
		if !result {
			check.LogError("Container %q has failed the RHEL release check", cut)
			nonCompliantObjects = append(nonCompliantObjects, testhelper.NewContainerReportObject(cut.Namespace, cut.Podname, cut.Name, "Failed the RHEL release check", false))
		} else {
			check.LogInfo("Container %q has passed the RHEL release check", cut)
			compliantObjects = append(compliantObjects, testhelper.NewContainerReportObject(cut.Namespace, cut.Podname, cut.Name, "Passed the RHEL release check", true))
		}
	}

	check.SetResult(compliantObjects, nonCompliantObjects)
}

func testIsSELinuxEnforcing(check *checksdb.Check, env *provider.TestEnvironment) {
	const (
		getenforceCommand = `chroot /host getenforce`
		enforcingString   = "Enforcing\n"
	)
	var compliantObjects []*testhelper.ReportObject
	var nonCompliantObjects []*testhelper.ReportObject
	o := clientsholder.GetClientsHolder()
	nodesFailed := 0
	nodesError := 0
	for _, debugPod := range env.DebugPods {
		ctx := clientsholder.NewContext(debugPod.Namespace, debugPod.Name, debugPod.Spec.Containers[0].Name)
		outStr, errStr, err := o.ExecCommandContainer(ctx, getenforceCommand)
		if err != nil || errStr != "" {
			check.LogError("Could not execute command %q in Debug Pod %q, errStr: %q, err: %v", getenforceCommand, debugPod, errStr, err)
			nonCompliantObjects = append(nonCompliantObjects, testhelper.NewPodReportObject(debugPod.Namespace, debugPod.Name, "Failed to execute command", false))
			nodesError++
			continue
		}
		if outStr != enforcingString {
			check.LogError("Node %q is not running SELinux, %s command returned: %s", debugPod.Spec.NodeName, getenforceCommand, outStr)
			nonCompliantObjects = append(nonCompliantObjects, testhelper.NewNodeReportObject(debugPod.Spec.NodeName, "SELinux is not enforced", false))
			nodesFailed++
		} else {
			check.LogInfo("Node %q is running SELinux", debugPod.Spec.NodeName)
			compliantObjects = append(compliantObjects, testhelper.NewNodeReportObject(debugPod.Spec.NodeName, "SELinux is enforced", true))
		}
	}

	check.SetResult(compliantObjects, nonCompliantObjects)
}

func testHugepages(check *checksdb.Check, env *provider.TestEnvironment) {
	var compliantObjects []*testhelper.ReportObject
	var nonCompliantObjects []*testhelper.ReportObject
	for i := range env.Nodes {
		node := env.Nodes[i]
		nodeName := node.Data.Name
		check.LogInfo("Testing node %q", nodeName)
		if !node.IsWorkerNode() {
			check.LogInfo("Node %q is not a worker node", nodeName)
			compliantObjects = append(compliantObjects, testhelper.NewNodeReportObject(nodeName, "Not a worker node", true))
			continue
		}

		debugPod, exist := env.DebugPods[nodeName]
		if !exist {
			check.LogError("Could not find a Debug Pod in node %q.", nodeName)
			nonCompliantObjects = append(nonCompliantObjects, testhelper.NewNodeReportObject(nodeName, "tnf debug pod not found", false))
			continue
		}

		hpTester, err := hugepages.NewTester(&node, debugPod, clientsholder.GetClientsHolder())
		if err != nil {
			check.LogError("Unable to get node hugepages tester for node %q, err: %v", nodeName, err)
			nonCompliantObjects = append(nonCompliantObjects, testhelper.NewNodeReportObject(nodeName, "Unable to get node hugepages tester", false))
		}

		if err := hpTester.Run(); err != nil {
			check.LogError("Hugepages check failed for node %q, err: %v", nodeName, err)
			nonCompliantObjects = append(nonCompliantObjects, testhelper.NewNodeReportObject(nodeName, err.Error(), false))
		} else {
			check.LogInfo("Node %q passed the hugepages check", nodeName)
			compliantObjects = append(compliantObjects, testhelper.NewNodeReportObject(nodeName, "Passed the hugepages check", true))
		}
	}

	check.SetResult(compliantObjects, nonCompliantObjects)
}

func testUnalteredBootParams(check *checksdb.Check, env *provider.TestEnvironment) {
	var compliantObjects []*testhelper.ReportObject
	var nonCompliantObjects []*testhelper.ReportObject
	alreadyCheckedNodes := map[string]bool{}
	for _, cut := range env.Containers {
		check.LogInfo("Testing Container %q", cut)
		if alreadyCheckedNodes[cut.NodeName] {
			check.LogInfo("Skipping node %q: already checked.", cut.NodeName)
			continue
		}
		alreadyCheckedNodes[cut.NodeName] = true

		err := bootparams.TestBootParamsHelper(env, cut, check.GetLogger())
		if err != nil {
			check.LogError("Node %q failed the boot params check", cut.NodeName)
			nonCompliantObjects = append(nonCompliantObjects, testhelper.NewNodeReportObject(cut.NodeName, "Failed the boot params check", false).
				AddField(testhelper.DebugPodName, env.DebugPods[cut.NodeName].Name))
		} else {
			check.LogInfo("Node %q passed the boot params check", cut.NodeName)
			compliantObjects = append(compliantObjects, testhelper.NewNodeReportObject(cut.NodeName, "Passed the boot params check", true).
				AddField(testhelper.DebugPodName, env.DebugPods[cut.NodeName].Name))
		}
	}

	check.SetResult(compliantObjects, nonCompliantObjects)
}

func testSysctlConfigs(check *checksdb.Check, env *provider.TestEnvironment) {
	var compliantObjects []*testhelper.ReportObject
	var nonCompliantObjects []*testhelper.ReportObject

	alreadyCheckedNodes := map[string]bool{}
	for _, cut := range env.Containers {
		check.LogInfo("Testing Container %q", cut)
		if alreadyCheckedNodes[cut.NodeName] {
			continue
		}
		alreadyCheckedNodes[cut.NodeName] = true
		debugPod := env.DebugPods[cut.NodeName]
		if debugPod == nil {
			check.LogError("Debug Pod not found for node %q", cut.NodeName)
			nonCompliantObjects = append(nonCompliantObjects, testhelper.NewNodeReportObject(cut.NodeName, "tnf debug pod not found", false))
			continue
		}

		sysctlSettings, err := sysctlconfig.GetSysctlSettings(env, cut.NodeName)
		if err != nil {
			check.LogError("Could not get sysctl settings for node %q, error: %v", cut.NodeName, err)
			nonCompliantObjects = append(nonCompliantObjects, testhelper.NewNodeReportObject(cut.NodeName, "Could not get sysctl settings", false))
			continue
		}

		mcKernelArgumentsMap := bootparams.GetMcKernelArguments(env, cut.NodeName)
		validSettings := true
		for key, sysctlConfigVal := range sysctlSettings {
			if mcVal, ok := mcKernelArgumentsMap[key]; ok {
				if mcVal != sysctlConfigVal {
					check.LogError("Kernel config mismatch in node %q for %q (sysctl value: %q, machine config value: %q)",
						cut.NodeName, key, sysctlConfigVal, mcVal)
					nonCompliantObjects = append(nonCompliantObjects, testhelper.NewNodeReportObject(cut.NodeName, fmt.Sprintf("Kernel config mismatch for %s", key), false))
					validSettings = false
				}
			}
		}
		if validSettings {
			check.LogInfo("Node %q passed the sysctl config check", cut.NodeName)
			compliantObjects = append(compliantObjects, testhelper.NewNodeReportObject(cut.NodeName, "Passed the sysctl config check", true))
		}
	}

	check.SetResult(compliantObjects, nonCompliantObjects)
}

func testOCPStatus(check *checksdb.Check, env *provider.TestEnvironment) {
	clusterIsInEOL := false
	switch env.OCPStatus {
	case compatibility.OCPStatusEOL:
		check.LogError("OCP Version %q has been found to be in end of life", env.OpenshiftVersion)
		clusterIsInEOL = true
	case compatibility.OCPStatusMS:
		check.LogInfo("OCP Version %q has been found to be in maintenance support", env.OpenshiftVersion)
	case compatibility.OCPStatusGA:
		check.LogInfo("OCP Version %q has been found to be in general availability", env.OpenshiftVersion)
	case compatibility.OCPStatusPreGA:
		check.LogInfo("OCP Version %q has been found to be in pre-general availability", env.OpenshiftVersion)
	default:
		check.LogInfo("OCP Version %q was unable to be found in the lifecycle compatibility matrix", env.OpenshiftVersion)
	}

	var compliantObjects []*testhelper.ReportObject
	var nonCompliantObjects []*testhelper.ReportObject

	if clusterIsInEOL {
		nonCompliantObjects = []*testhelper.ReportObject{testhelper.NewClusterVersionReportObject(env.OpenshiftVersion, "Openshift Version is in End Of Life (EOL)", false)}
	} else {
		compliantObjects = []*testhelper.ReportObject{testhelper.NewClusterVersionReportObject(env.OpenshiftVersion, "Openshift Version is not in End Of Life (EOL)", true)}
	}

	check.SetResult(compliantObjects, nonCompliantObjects)
}

//nolint:funlen
func testNodeOperatingSystemStatus(check *checksdb.Check, env *provider.TestEnvironment) {
	failedControlPlaneNodes := []string{}
	failedWorkerNodes := []string{}
	var compliantObjects []*testhelper.ReportObject
	var nonCompliantObjects []*testhelper.ReportObject
	for _, node := range env.Nodes {
		nodeName := node.Data.Name
		check.LogInfo("Testing node %q", nodeName)
		// Get the OSImage which should tell us what version of operating system the node is running.
		check.LogInfo("Node %q is running operating system %q", nodeName, node.Data.Status.NodeInfo.OSImage)

		// Control plane nodes must be RHCOS (also CentOS Stream starting in OCP 4.13)
		// Per the release notes from OCP documentation:
		// "You must use RHCOS machines for the control plane, and you can use either RHCOS or RHEL for compute machines."
		if node.IsControlPlaneNode() && !node.IsRHCOS() && !node.IsCSCOS() {
			check.LogError("Control plane node %q has been found to be running an incompatible operating system %q", nodeName, node.Data.Status.NodeInfo.OSImage)
			failedControlPlaneNodes = append(failedControlPlaneNodes, nodeName)
			nonCompliantObjects = append(nonCompliantObjects, testhelper.NewNodeReportObject(nodeName, "Control plane node has been found to be running an incompatible OS", false).AddField(testhelper.OSImage, node.Data.Status.NodeInfo.OSImage))
			continue
		}

		// Worker nodes can either be RHEL or RHCOS
		if node.IsWorkerNode() {
			//nolint:gocritic
			if node.IsRHCOS() {
				// Get the short version from the node
				shortVersion, err := node.GetRHCOSVersion()
				if err != nil {
					check.LogError("Node %q failed to gather RHCOS version, err: %v", nodeName, err)
					failedWorkerNodes = append(failedWorkerNodes, nodeName)
					nonCompliantObjects = append(nonCompliantObjects, testhelper.NewNodeReportObject(nodeName, "Failed to gather RHCOS version", false))
					continue
				}

				if shortVersion == operatingsystem.NotFoundStr {
					check.LogInfo("Node %q has an RHCOS operating system that is not found in our internal database. Skipping as to not cause failures due to database mismatch.", nodeName)
					continue
				}

				// If the node's RHCOS version and the OpenShift version are not compatible, the node fails.
				check.LogDebug("Comparing RHCOS shortVersion %q to openshiftVersion %q", shortVersion, env.OpenshiftVersion)
				if !compatibility.IsRHCOSCompatible(shortVersion, env.OpenshiftVersion) {
					check.LogError("Worker node %q has been found to be running an incompatible version of RHCOS %q", nodeName, shortVersion)
					failedWorkerNodes = append(failedWorkerNodes, nodeName)
					nonCompliantObjects = append(nonCompliantObjects, testhelper.NewNodeReportObject(nodeName, "Worker node has been found to be running an incompatible OS", false).
						AddField(testhelper.OSImage, node.Data.Status.NodeInfo.OSImage))
					continue
				}
				check.LogInfo("Worker node %q has been found to be running a compatible version of RHCOS %q", nodeName, shortVersion)
				compliantObjects = append(compliantObjects, testhelper.NewNodeReportObject(nodeName, "Worker node has been found to be running a compatible OS", true).
					AddField(testhelper.OSImage, node.Data.Status.NodeInfo.OSImage))
			} else if node.IsCSCOS() {
				// Get the short version from the node
				shortVersion, err := node.GetCSCOSVersion()
				if err != nil {
					check.LogError("Node %q failed to gather CentOS Stream CoreOS version, err: %v", nodeName, err)
					failedWorkerNodes = append(failedWorkerNodes, nodeName)
					nonCompliantObjects = append(nonCompliantObjects, testhelper.NewNodeReportObject(nodeName, "Failed to gather CentOS Stream CoreOS version", false))
					continue
				}

				// Warning: CentOS Stream CoreOS has not been released yet in any
				// OCP RC/GA versions, so for the moment, we cannot compare the
				// version with the OCP one, or retrieve it on the internal database
				msg := `
					Node %s is using CentOS Stream CoreOS %s, which is not being used yet in any
					OCP RC/GA version. Relaxing the conditions to check the OS as a result.
					`
				check.LogDebug(msg, nodeName, shortVersion)
			} else if node.IsRHEL() {
				// Get the short version from the node
				shortVersion, err := node.GetRHELVersion()
				if err != nil {
					check.LogError("Node %q failed to gather RHEL version, err: %v", nodeName, err)
					failedWorkerNodes = append(failedWorkerNodes, nodeName)
					nonCompliantObjects = append(nonCompliantObjects, testhelper.NewNodeReportObject(nodeName, "Failed to gather RHEL version", false))
					continue
				}

				// If the node's RHEL version and the OpenShift version are not compatible, the node fails.
				check.LogDebug("Comparing RHEL shortVersion %q to openshiftVersion %q", shortVersion, env.OpenshiftVersion)
				if !compatibility.IsRHELCompatible(shortVersion, env.OpenshiftVersion) {
					check.LogError("Worker node %q has been found to be running an incompatible version of RHEL %q", nodeName, shortVersion)
					failedWorkerNodes = append(failedWorkerNodes, nodeName)
					nonCompliantObjects = append(nonCompliantObjects, testhelper.NewNodeReportObject(nodeName, "Worker node has been found to be running an incompatible OS", false).AddField(testhelper.OSImage, node.Data.Status.NodeInfo.OSImage))
				} else {
					check.LogInfo("Worker node %q has been found to be running a compatible version of RHEL %q", nodeName, shortVersion)
					compliantObjects = append(compliantObjects, testhelper.NewNodeReportObject(nodeName, "Worker node has been found to be running a compatible OS", true).AddField(testhelper.OSImage, node.Data.Status.NodeInfo.OSImage))
				}
			} else {
				check.LogError("Worker node %q has been found to be running an incompatible operating system %q", nodeName, node.Data.Status.NodeInfo.OSImage)
				failedWorkerNodes = append(failedWorkerNodes, nodeName)
				nonCompliantObjects = append(nonCompliantObjects, testhelper.NewNodeReportObject(nodeName, "Worker node has been found to be running an incompatible OS", false).AddField(testhelper.OSImage, node.Data.Status.NodeInfo.OSImage))
			}
		}
	}

	if n := len(failedControlPlaneNodes); n > 0 {
		check.LogError("Number of control plane nodes running non-RHCOS based operating systems: %d", n)
	}

	if n := len(failedWorkerNodes); n > 0 {
		check.LogError("Number of worker nodes running non-RHCOS or non-RHEL based operating systems: %d", n)
	}

	check.SetResult(compliantObjects, nonCompliantObjects)
}

func testPodHugePagesSize(check *checksdb.Check, env *provider.TestEnvironment, size string) {
	var compliantObjects []*testhelper.ReportObject
	var nonCompliantObjects []*testhelper.ReportObject
	for _, put := range env.GetHugepagesPods() {
		check.LogInfo("Testing Pod %q", put)
		result := put.CheckResourceHugePagesSize(size)
		if !result {
			check.LogError("Pod %q has been found to be running with an incorrect hugepages size (expected size %q)", put, size)
			nonCompliantObjects = append(nonCompliantObjects, testhelper.NewPodReportObject(put.Namespace, put.Name, "Pod has been found to be running with an incorrect hugepages size", false))
		} else {
			check.LogInfo("Pod %q has been found to be running with a correct hugepages size %q", put, size)
			compliantObjects = append(compliantObjects, testhelper.NewPodReportObject(put.Namespace, put.Name, "Pod has been found to be running with a correct hugepages size", true))
		}
	}
	check.SetResult(compliantObjects, nonCompliantObjects)
}
