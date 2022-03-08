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

package lifecycle

import (
	"fmt"
	"time"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
	"github.com/test-network-function/cnf-certification-test/cnf-certification-test/common"
	"github.com/test-network-function/cnf-certification-test/cnf-certification-test/identifiers"
	"github.com/test-network-function/cnf-certification-test/cnf-certification-test/lifecycle/graceperiod"
	"github.com/test-network-function/cnf-certification-test/cnf-certification-test/lifecycle/ownerreference"
	"github.com/test-network-function/cnf-certification-test/cnf-certification-test/lifecycle/podrecreation"
	"github.com/test-network-function/cnf-certification-test/cnf-certification-test/lifecycle/podsets"
	"github.com/test-network-function/cnf-certification-test/cnf-certification-test/lifecycle/scaling"
	"github.com/test-network-function/cnf-certification-test/pkg/provider"
	"github.com/test-network-function/cnf-certification-test/pkg/testhelper"
	"github.com/test-network-function/cnf-certification-test/pkg/tnf"

	v1 "k8s.io/api/core/v1"
)

const (
	timeout                    = 60 * time.Second
	timeoutPodRecreationPerPod = time.Minute
	timeoutPodSetReady         = 7 * time.Minute
)

//
// All actual test code belongs below here.  Utilities belong above.
//
var _ = ginkgo.Describe(common.LifecycleTestKey, func() {
	var env provider.TestEnvironment
	ginkgo.BeforeEach(func() {
		env = provider.GetTestEnvironment()
	})
	testContainersPreStop(&env)
	testContainersImagePolicy(&env)
	testContainersReadinessProbe(&env)
	testContainersLivenessProbe(&env)
	testPodsOwnerReference(&env)
	testHighAvailability(&env)

	testID := identifiers.XformToGinkgoItIdentifier(identifiers.TestPodNodeSelectorAndAffinityBestPractices)
	ginkgo.It(testID, ginkgo.Label(testID), func() {
		testPodNodeSelectorAndAffinityBestPractices(&env)
	})

	testID = identifiers.XformToGinkgoItIdentifier(identifiers.TestNonDefaultGracePeriodIdentifier)
	ginkgo.It(testID, ginkgo.Label(testID), func() {
		testGracePeriod(&env)
	})
	testID = identifiers.XformToGinkgoItIdentifier(identifiers.TestPodRecreationIdentifier)
	ginkgo.It(testID, ginkgo.Label(testID), func() {
		// Testing pod re-creation for deployments
		testPodsRecreationDeployment(&env)

		// Testing pod re-creation for statefulsets
		testPodsRecreationStatefulset(&env)
	})

	if env.IsIntrusive() {
		testScaling(&env, timeout)
	}
})

func testContainersPreStop(env *provider.TestEnvironment) {
	testID := identifiers.XformToGinkgoItIdentifier(identifiers.TestShudtownIdentifier)
	ginkgo.It(testID, ginkgo.Label(testID), func() {
		badcontainers := []string{}
		for _, cut := range env.Containers {
			logrus.Debugln("check container ", cut.Namespace, " ", cut.Podname, " ", cut.Data.Name, " pre stop lifecycle ")

			if cut.Data.Lifecycle == nil || (cut.Data.Lifecycle != nil && cut.Data.Lifecycle.PreStop == nil) {
				badcontainers = append(badcontainers, cut.Data.Name)
				tnf.ClaimFilePrintf("container %s does not have preStop defined", cut.StringShort())
			}
		}
		if len(badcontainers) > 0 {
			tnf.ClaimFilePrintf("bad containers %v", badcontainers)
		}
		gomega.Expect(0).To(gomega.Equal(len(badcontainers)))
	})
}

func testContainersImagePolicy(env *provider.TestEnvironment) {
	testID := identifiers.XformToGinkgoItIdentifier(identifiers.TestImagePullPolicyIdentifier)
	ginkgo.It(testID, ginkgo.Label(testID), func() {
		badcontainers := []string{}
		for _, cut := range env.Containers {
			logrus.Debugln("check container ", cut.Namespace, " ", cut.Podname, " ", cut.Data.Name, " pull policy, should be ", v1.PullIfNotPresent)
			if cut.Data.ImagePullPolicy != v1.PullIfNotPresent {
				s := cut.Namespace + ":" + cut.Podname + ":" + cut.Data.Name
				badcontainers = append(badcontainers, s)
				logrus.Errorln("container ", cut.Data.Name, " is using ", cut.Data.ImagePullPolicy, " as image policy")
			}
		}
		if len(badcontainers) > 0 {
			tnf.ClaimFilePrintf("bad containers %v", badcontainers)
		}
		gomega.Expect(0).To(gomega.Equal(len(badcontainers)))
	})
}

//nolint:dupl
func testContainersReadinessProbe(env *provider.TestEnvironment) {
	testID := identifiers.XformToGinkgoItIdentifier(identifiers.TestReadinessProbeIdentifier)
	ginkgo.It(testID, ginkgo.Label(testID), func() {
		badcontainers := []string{}
		for _, cut := range env.Containers {
			logrus.Debugln("check container ", cut.Namespace, " ", cut.Podname, " ", cut.Data.Name, " readiness probe ")
			if cut.Data.ReadinessProbe == nil {
				s := cut.Namespace + ":" + cut.Podname + ":" + cut.Data.Name
				badcontainers = append(badcontainers, s)
				logrus.Errorln("container ", cut.Data.Name, " does not have ReadinessProbe defined")
			}
		}
		if len(badcontainers) > 0 {
			tnf.ClaimFilePrintf("bad containers %v", badcontainers)
		}
		gomega.Expect(0).To(gomega.Equal(len(badcontainers)))
	})
}

//nolint:dupl
func testContainersLivenessProbe(env *provider.TestEnvironment) {
	testID := identifiers.XformToGinkgoItIdentifier(identifiers.TestLivenessProbeIdentifier)
	ginkgo.It(testID, ginkgo.Label(testID), func() {
		badcontainers := []string{}
		for _, cut := range env.Containers {
			logrus.Debugln("check container ", cut.Namespace, " ", cut.Podname, " ", cut.Data.Name, " liveness probe ")
			if cut.Data.LivenessProbe == nil {
				s := cut.Namespace + ":" + cut.Podname + ":" + cut.Data.Name
				badcontainers = append(badcontainers, s)
				logrus.Errorln("container ", cut.Data.Name, " does not have livenessProbe defined")
			}
		}
		if len(badcontainers) > 0 {
			tnf.ClaimFilePrintf("bad containers %v", badcontainers)
		}
		gomega.Expect(0).To(gomega.Equal(len(badcontainers)))
	})
}

func testPodsOwnerReference(env *provider.TestEnvironment) {
	testID := identifiers.XformToGinkgoItIdentifier(identifiers.TestPodDeploymentBestPracticesIdentifier)
	ginkgo.It(testID, ginkgo.Label(testID), func() {
		ginkgo.By("Testing owners of CNF pod, should be replicas Set")
		badPods := []string{}
		for _, put := range env.Pods {
			logrus.Debugln("check pod ", put.Namespace, " ", put.Name, " owner reference")
			o := ownerreference.NewOwnerReference(put)
			o.RunTest()
			if o.GetResults() != testhelper.SUCCESS {
				s := put.Namespace + ":" + put.Name
				badPods = append(badPods, s)
			}
		}
		if len(badPods) > 0 {
			tnf.ClaimFilePrintf("bad containers %v", badPods)
		}
		gomega.Expect(0).To(gomega.Equal(len(badPods)))
	})
}

func testPodNodeSelectorAndAffinityBestPractices(env *provider.TestEnvironment) {
	var badPods []*v1.Pod
	for _, put := range env.Pods {
		if len(put.Spec.NodeSelector) != 0 {
			tnf.ClaimFilePrintf("ERROR: Pod: %s has a node selector clause. Node selector: %v", provider.PodToString(put), &put.Spec.NodeSelector)
			badPods = append(badPods, put)
		}
		if put.Spec.Affinity != nil && put.Spec.Affinity.NodeAffinity != nil {
			tnf.ClaimFilePrintf("ERROR: Pod: %s has a node affinity clause. Node affinity: %v", provider.PodToString(put), put.Spec.Affinity.NodeAffinity)
			badPods = append(badPods, put)
		}
	}
	if n := len(badPods); n > 0 {
		logrus.Debugf("Pods with nodeSelector/nodeAffinity: %+v", badPods)
		ginkgo.Fail(fmt.Sprintf("%d pods found with nodeSelector/nodeAffinity rules", n))
	}
}

func testGracePeriod(env *provider.TestEnvironment) {
	badDeployments, deploymentLogs := graceperiod.TestTerminationGracePeriodOnDeployments(env)
	badStatefulsets, statefulsetLogs := graceperiod.TestTerminationGracePeriodOnStatefulsets(env)
	badPods, podLogs := graceperiod.TestTerminationGracePeriodOnPods(env)

	numDeps := len(badDeployments)
	if numDeps > 0 {
		logrus.Debugf("Deployments found without terminationGracePeriodSeconds param set: %+v", badDeployments)
	}
	numSts := len(badStatefulsets)
	if numSts > 0 {
		logrus.Debugf("Statefulsets found without terminationGracePeriodSeconds param set: %+v", badStatefulsets)
	}
	numPods := len(badPods)
	if numPods > 0 {
		logrus.Debugf("Pods found without terminationGracePeriodSeconds param set: %+v", badPods)
	}
	ginkgo.By("Test results for grace period on deployments")
	tnf.ClaimFilePrintf("%s", deploymentLogs)
	ginkgo.By("Test results for grace period on statefulsets")
	tnf.ClaimFilePrintf("%s", statefulsetLogs)
	ginkgo.By("Test results for grace period on unmanaged pods")
	tnf.ClaimFilePrintf("%s", podLogs)

	if numDeps > 0 || numSts > 0 || numPods > 0 {
		ginkgo.Fail(fmt.Sprintf("Found %d deployments, %d statefulsets and %d pods without terminationGracePeriodSeconds param set.", numDeps, numSts, numPods))
	}
}

func testScaling(env *provider.TestEnvironment, timeout time.Duration) {
	testID := identifiers.XformToGinkgoItIdentifier(identifiers.TestDeploymentScalingIdentifier)
	ginkgo.It(testID, ginkgo.Label(testID), func() {
		ginkgo.By("Testing deployment scaling")
		defer env.SetNeedsRefresh()

		if len(env.Deployments) == 0 {
			ginkgo.Skip("No test deployments found.")
		}
		failedDeployments := []string{}
		skippedDeployments := []string{}
		for i := range env.Deployments {
			// TestDeploymentScaling test scaling of deployment
			// This is the entry point for deployment scaling tests
			deployment := env.Deployments[i]
			ns, name := deployment.Namespace, deployment.Name
			key := ns + name
			if hpa, ok := env.HorizontalScaler[key]; ok {
				// if the deployment is controller by
				// horizontal scaler, then test that scaler
				// can scale the deployment
				if !scaling.TestScaleHpaDeployment(deployment, hpa, timeout) {
					failedDeployments = append(failedDeployments, name)
				}
				continue
			}
			// if the deployment is not controller by HPA
			// scale it directly
			if !scaling.TestScaleDeployment(deployment, timeout) {
				failedDeployments = append(failedDeployments, name)
			}
		}

		if len(skippedDeployments) > 0 {
			tnf.ClaimFilePrintf("not ready deployments : %v", skippedDeployments)
		}
		if len(failedDeployments) > 0 {
			tnf.ClaimFilePrintf(" failed deployments: %v", failedDeployments)
		}
		gomega.Expect(0).To(gomega.Equal(len(failedDeployments)))
		gomega.Expect(0).To(gomega.Equal(len(skippedDeployments)))
	})
}

// testHighAvailability
func testHighAvailability(env *provider.TestEnvironment) {
	testID := identifiers.XformToGinkgoItIdentifier(identifiers.TestPodHighAvailabilityBestPractices)
	ginkgo.It(testID, ginkgo.Label(testID), func() {
		ginkgo.By("Should set pod replica number greater than 1")
		if len(env.Deployments) == 0 && len(env.SatetfulSets) == 0 {
			ginkgo.Skip("No test deployments/statefulset found.")
		}

		badDeployments := []string{}
		badStatefulSet := []string{}
		for _, dp := range env.Deployments {
			if dp.Spec.Replicas == nil || *(dp.Spec.Replicas) == 1 {
				badDeployments = append(badDeployments, provider.DeploymentToString(dp))
			}
		}
		for _, st := range env.SatetfulSets {
			if st.Spec.Replicas == nil || *(st.Spec.Replicas) == 1 {
				badStatefulSet = append(badStatefulSet, provider.StatefulsetToString(st))
			}
		}

		if n := len(badDeployments); n > 0 {
			logrus.Errorf("Deployments without a valid high availability : %+v", badDeployments)
			tnf.ClaimFilePrintf("Deployments without a valid high availability : %+v", badDeployments)
		}
		if n := len(badStatefulSet); n > 0 {
			logrus.Errorf("Statefulset without a valid podAntiAffinity rule: %+v", badStatefulSet)
			tnf.ClaimFilePrintf("Statefulset without a valid podAntiAffinity rule: %+v", badStatefulSet)
		}
		gomega.Expect(0).To(gomega.Equal(len(badDeployments)))
		gomega.Expect(0).To(gomega.Equal(len(badStatefulSet)))
	})
}

// testPodsRecreationDeployment tests that pods belonging to deployments are re-created and ready in case a node is lost
func testPodsRecreationDeployment(env *provider.TestEnvironment) { //nolint:dupl,funlen // not duplicate
	ginkgo.By("Testing node draining effect of deployment")
	for _, dut := range env.Deployments {
		ginkgo.By(fmt.Sprintf("Testing pod-recreation for deployment: %s", provider.DeploymentToString(dut)))
		isReady := podsets.WaitForDeploymentSetReady(dut.Namespace, dut.Name, timeoutPodSetReady)
		if !isReady {
			tnf.ClaimFilePrintf("deployment: %s is not in a good starting state", provider.DeploymentToString(dut))
			continue
		}
		nodes := podrecreation.GetDeploymentNodes(env.Pods, dut.Name)
		for _, n := range nodes {
			err := podrecreation.CordonNode(n)
			if err != nil {
				logrus.Errorf("error cordoning the node: %s", n)
				err = podrecreation.UncordonNode(n)
				if err != nil {
					logrus.Fatalf("error uncordoning the node: %s", n)
				}
			}
			logrus.Debugf("deployment: %s, node: %s cordoned", provider.DeploymentToString(dut), n)
			count := podrecreation.CountPods(n)
			nodeTimeout := timeoutPodRecreationPerPod * time.Duration(count)
			logrus.Debugf("deployment %s, draining node: %s with timeout: %s", provider.DeploymentToString(dut), n, nodeTimeout.String())
			podrecreation.DeletePods(n)
			isReady := podsets.WaitForDeploymentSetReady(dut.Namespace, dut.Name, nodeTimeout)
			if !isReady {
				tnf.ClaimFilePrintf("deployment: %s recovery, NOK after loosing node: %s", provider.DeploymentToString(dut), n)
			} else {
				tnf.ClaimFilePrintf("deployment: %s recovery, OK after loosing node: %s", provider.DeploymentToString(dut), n)
			}
			err = podrecreation.UncordonNode(n)
			if err != nil {
				logrus.Fatalf("error uncordoning the node: %s", n)
			}
		}
	}
}

// testPodsRecreationDeployment tests that pods belonging to statefulsets are re-created and ready in case a node is lost
func testPodsRecreationStatefulset(env *provider.TestEnvironment) { //nolint:dupl,funlen // not duplicate
	ginkgo.By("Testing node draining effect of statefulset")
	for _, sut := range env.SatetfulSets {
		ginkgo.By(fmt.Sprintf("Testing pod-recreation for statefulset %s", provider.StatefulsetToString(sut)))
		isReady := podsets.WaitForStatefulSetReady(sut.Namespace, sut.Name, timeoutPodSetReady)
		if !isReady {
			tnf.ClaimFilePrintf("statefulset %s is not in a good starting state", provider.StatefulsetToString(sut))
			continue
		}
		nodes := podrecreation.GetStatefulsetNodes(env.Pods, sut.Name)
		for _, n := range nodes {
			err := podrecreation.CordonNode(n)
			if err != nil {
				logrus.Errorf("error cordoning the node: %s", n)
				err = podrecreation.UncordonNode(n)
				if err != nil {
					logrus.Fatalf("error uncordoning the node: %s", n)
				}
			}
			logrus.Debugf("statefulset: %s, node: %s cordoned", provider.StatefulsetToString(sut), n)
			count := podrecreation.CountPods(n)
			nodeTimeout := timeoutPodRecreationPerPod * time.Duration(count)
			logrus.Debugf("statefulset %s, draining node: %s with timeout: %s", provider.StatefulsetToString(sut), n, nodeTimeout.String())
			podrecreation.DeletePods(n)
			isReady := podsets.WaitForStatefulSetReady(sut.Namespace, sut.Name, nodeTimeout)
			if !isReady {
				tnf.ClaimFilePrintf("statefulset %s, recovery NOK loosing node: %s", provider.StatefulsetToString(sut), n)
			} else {
				tnf.ClaimFilePrintf("statefulset %s, recovery OK after loosing node: %s", provider.StatefulsetToString(sut), n)
			}
			err = podrecreation.UncordonNode(n)
			if err != nil {
				logrus.Fatalf("error uncordoning the node: %s", n)
			}
		}
	}
}
