package hugepages

import (
	"errors"
	"testing"

	mcv1 "github.com/openshift/api/machineconfiguration/v1"
	"github.com/stretchr/testify/assert"

	"github.com/test-network-function/cnf-certification-test/internal/clientsholder"
	"github.com/test-network-function/cnf-certification-test/pkg/provider"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	// No hugepages params.
	testKernelArgsHpNoParams = []string{"systemd.cpu_affinity=0,1,40,41,20,21,60,61", "nmi_watchdog=0"}

	// Single param.
	testKernelArgsHpSingleParam1 = []string{"systemd.cpu_affinity=0,1,40,41,20,21,60,61", "hugepages=16", "nmi_watchdog=0"}
	testKernelArgsHpSingleParam2 = []string{"systemd.cpu_affinity=0,1,40,41,20,21,60,61", "default_hugepagesz=1G", "nmi_watchdog=0"}
	testKernelArgsHpSingleParam3 = []string{"systemd.cpu_affinity=0,1,40,41,20,21,60,61", "default_hugepagesz=2M", "nmi_watchdog=0"}
	testKernelArgsHpSingleParam4 = []string{"systemd.cpu_affinity=0,1,40,41,20,21,60,61", "hugepagesz=1G", "nmi_watchdog=0"}

	// Default size + size only.
	testKernelArgsHpDefParamsOnly = []string{"systemd.cpu_affinity=0,1,40,41,20,21,60,61", "default_hugepagesz=1G", "hugepagesz=1G", "nmi_watchdog=0"}

	// size + count pairs.
	testKernelArgsHpPair1 = []string{"systemd.cpu_affinity=0,1,40,41,20,21,60,61", "hugepagesz=1G", "hugepages=16", "nmi_watchdog=0"}
	testKernelArgsHpPair2 = []string{"systemd.cpu_affinity=0,1,40,41,20,21,60,61", "hugepagesz=2M", "hugepages=256", "nmi_watchdog=0"}
	testKernelArgsHpPair3 = []string{"systemd.cpu_affinity=0,1,40,41,20,21,60,61", "hugepagesz=1G", "hugepages=16", "hugepagesz=2M", "hugepages=256", "nmi_watchdog=0"}

	// default size + (size+count) pairs.
	testKernelArgsHpDefSizePlusPairs1 = []string{"systemd.cpu_affinity=0,1,40,41,20,21,60,61", "default_hugepagesz=2M", "hugepagesz=1G", "hugepages=16", "nmi_watchdog=0"}
	testKernelArgsHpDefSizePlusPairs2 = []string{"systemd.cpu_affinity=0,1,40,41,20,21,60,61", "default_hugepagesz=1G", "hugepagesz=2M", "hugepages=256", "nmi_watchdog=0"}
	testKernelArgsHpDefSizePlusPairs3 = []string{"systemd.cpu_affinity=0,1,40,41,20,21,60,61", "default_hugepagesz=1G", "hugepagesz=1G", "hugepages=16", "hugepagesz=2M", "hugepages=256", "nmi_watchdog=0"}
)

const (
	// Sizes, in KBs.
	oneGB = 1024 * 1024 // 1G
	twoMB = 2 * 1024    // 2M: also RHEL's default hugepages size
)

func Test_hugepagesFromKernelArgsFunc(t *testing.T) {
	testCases := []struct {
		expectedHugepagesDefSize int
		expectedHugepagesPerSize map[int]int
		kernelArgs               []string
	}{
		// No params.
		{
			expectedHugepagesDefSize: twoMB,
			expectedHugepagesPerSize: map[int]int{twoMB: 0},
			kernelArgs:               testKernelArgsHpNoParams,
		},

		// Single params TCs.
		{
			expectedHugepagesDefSize: twoMB,
			expectedHugepagesPerSize: map[int]int{twoMB: 16},
			kernelArgs:               testKernelArgsHpSingleParam1,
		},
		{
			expectedHugepagesDefSize: oneGB,
			expectedHugepagesPerSize: map[int]int{oneGB: 0},
			kernelArgs:               testKernelArgsHpSingleParam2,
		},
		{
			expectedHugepagesDefSize: twoMB,
			expectedHugepagesPerSize: map[int]int{twoMB: 0},
			kernelArgs:               testKernelArgsHpSingleParam3,
		},
		{
			expectedHugepagesDefSize: twoMB,
			expectedHugepagesPerSize: map[int]int{oneGB: 0},
			kernelArgs:               testKernelArgsHpSingleParam4,
		},
		{
			expectedHugepagesDefSize: twoMB,
			expectedHugepagesPerSize: map[int]int{oneGB: 16},
			kernelArgs:               testKernelArgsHpPair1,
		},

		// Default sizes Tc.
		{
			expectedHugepagesDefSize: oneGB,
			expectedHugepagesPerSize: map[int]int{oneGB: 0},
			kernelArgs:               testKernelArgsHpDefParamsOnly,
		},

		// size+count pairs.
		{
			expectedHugepagesDefSize: twoMB,
			expectedHugepagesPerSize: map[int]int{oneGB: 16},
			kernelArgs:               testKernelArgsHpPair1,
		},
		{
			expectedHugepagesDefSize: twoMB,
			expectedHugepagesPerSize: map[int]int{twoMB: 256},
			kernelArgs:               testKernelArgsHpPair2,
		},
		{
			expectedHugepagesDefSize: twoMB,
			expectedHugepagesPerSize: map[int]int{oneGB: 16, twoMB: 256},
			kernelArgs:               testKernelArgsHpPair3,
		},

		// default size + (size+count) pairs.
		{
			expectedHugepagesDefSize: twoMB,
			expectedHugepagesPerSize: map[int]int{twoMB: 0, oneGB: 16},
			kernelArgs:               testKernelArgsHpDefSizePlusPairs1,
		},
		{
			expectedHugepagesDefSize: oneGB,
			expectedHugepagesPerSize: map[int]int{oneGB: 0, twoMB: 256},
			kernelArgs:               testKernelArgsHpDefSizePlusPairs2,
		},
		{
			expectedHugepagesDefSize: oneGB,
			expectedHugepagesPerSize: map[int]int{oneGB: 16, twoMB: 256},
			kernelArgs:               testKernelArgsHpDefSizePlusPairs3,
		},
	}

	for _, tc := range testCases {
		mc := provider.MachineConfig{MachineConfig: &mcv1.MachineConfig{}}

		// Prepare fake MC object: only KernelArguments is needed.
		mc.Spec.KernelArguments = tc.kernelArgs

		// Call the function under test.
		hugepagesPerSize, defSize := getMcHugepagesFromMcKernelArguments(&mc)

		assert.Equal(t, defSize, tc.expectedHugepagesDefSize)
		assert.Equal(t, hugepagesPerSize, tc.expectedHugepagesPerSize)
	}
}

type fakeK8sClient struct {
	execCommandFunctionMocker func() (stdout string, stderr string, err error)
}

func (client *fakeK8sClient) ExecCommandContainer(_ clientsholder.Context, _ string) (stdout, stderr string, err error) {
	return client.execCommandFunctionMocker()
}

func TestPositiveMachineConfigSystemdHugepages(t *testing.T) {
	// helper pod, so the hugepages struct doesnt crash when accessing the debug container.
	fakeDebugPod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{Name: "pod1", Namespace: "ns1"},
		Spec:       corev1.PodSpec{Containers: []corev1.Container{{Name: "container1"}}},
	}

	// helper struct to hold the mc units for each tc
	type mcUnitConfig struct {
		name     string
		contents string
	}

	// helper function to get a provider.MachineConfig object from an slice of units to
	// be used on each tc.
	getMcFromUnits := func(units []mcUnitConfig) provider.MachineConfig {
		mc := provider.MachineConfig{MachineConfig: &mcv1.MachineConfig{}}

		for _, unit := range units {
			mc.Config.Systemd.Units = append(mc.Config.Systemd.Units, struct {
				Contents string "json:\"contents\""
				Name     string "json:\"name\""
			}{
				Name:     unit.name,
				Contents: unit.contents,
			})
		}
		return mc
	}

	testCases := []struct {
		nodeHugePagesCmdOutput string
		mcUnits                []mcUnitConfig
	}{
		// One numa with one size only (2MB).
		{
			nodeHugePagesCmdOutput: `/host/sys/devices/system/node/node0/hugepages/hugepages-2048kB/nr_hugepages count:4`,
			mcUnits: []mcUnitConfig{
				{
					name: "hugepages-allocation-2048kB-NUMA0.service",
					contents: `Environment=HUGEPAGES_COUNT=4
							   Environment=HUGEPAGES_SIZE=2048kB
							   Environment=NUMA_NODE=0`,
				},
			},
		},
		// One numa with two sizes (2MB and 1GB).
		{
			nodeHugePagesCmdOutput: `/host/sys/devices/system/node/node0/hugepages/hugepages-2048kB/nr_hugepages count:4
									 /host/sys/devices/system/node/node0/hugepages/hugepages-1048576kB/nr_hugepages count:8`,
			mcUnits: []mcUnitConfig{
				{
					name: "hugepages-allocation-2048kB-NUMA0.service",
					contents: `Environment=HUGEPAGES_COUNT=4
							   Environment=HUGEPAGES_SIZE=2048kB
							   Environment=NUMA_NODE=0`,
				},
				{
					name: "hugepages-allocation-1048576kB-NUMA0.service",
					contents: `Environment=HUGEPAGES_COUNT=8
							   Environment=HUGEPAGES_SIZE=1048576kB
							   Environment=NUMA_NODE=0`,
				},
			},
		},
		// Two numas, one size (2MB).
		{
			nodeHugePagesCmdOutput: `/host/sys/devices/system/node/node0/hugepages/hugepages-2048kB/nr_hugepages count:4
									 /host/sys/devices/system/node/node1/hugepages/hugepages-2048kB/nr_hugepages count:8`,
			mcUnits: []mcUnitConfig{
				{
					name: "hugepages-allocation-2048kB-NUMA0.service",
					contents: `Environment=HUGEPAGES_COUNT=4
							   Environment=HUGEPAGES_SIZE=2048kB
							   Environment=NUMA_NODE=0`,
				},
				{
					name: "hugepages-allocation-2048kB-NUMA1.service",
					contents: `Environment=HUGEPAGES_COUNT=8
							   Environment=HUGEPAGES_SIZE=2048kB
							   Environment=NUMA_NODE=1`,
				},
			},
		},
		// Two numas, two sizes (2MB and 1GB)
		{
			nodeHugePagesCmdOutput: `/host/sys/devices/system/node/node0/hugepages/hugepages-2048kB/nr_hugepages count:256
									 /host/sys/devices/system/node/node0/hugepages/hugepages-1048576kB/nr_hugepages count:8
									 /host/sys/devices/system/node/node1/hugepages/hugepages-2048kB/nr_hugepages count:256
									 /host/sys/devices/system/node/node1/hugepages/hugepages-1048576kB/nr_hugepages count:8`,
			mcUnits: []mcUnitConfig{
				{
					name: "hugepages-allocation-2048kB-NUMA0.service",
					contents: `Environment=HUGEPAGES_COUNT=256
							   Environment=HUGEPAGES_SIZE=2048kB
							   Environment=NUMA_NODE=0`,
				},
				{
					name: "hugepages-allocation-2048kB-NUMA1.service",
					contents: `Environment=HUGEPAGES_COUNT=256
							   Environment=HUGEPAGES_SIZE=2048kB
							   Environment=NUMA_NODE=1`,
				},
				{
					name: "hugepages-allocation-1048576kB-NUMA0.service",
					contents: `Environment=HUGEPAGES_COUNT=8
							   Environment=HUGEPAGES_SIZE=1048576kB
							   Environment=NUMA_NODE=0`,
				},
				{
					name: "hugepages-allocation-1048576kB-NUMA1.service",
					contents: `Environment=HUGEPAGES_COUNT=8
							   Environment=HUGEPAGES_SIZE=1048576kB
							   Environment=NUMA_NODE=1`,
				},
			},
		},
		// Size mismatch: size 1GB does not appear in mc, but there are 0 hugepages for that size.
		{
			nodeHugePagesCmdOutput: `/host/sys/devices/system/node/node0/hugepages/hugepages-2048kB/nr_hugepages count:4
									 /host/sys/devices/system/node/node0/hugepages/hugepages-1048576kB/nr_hugepages count:0`,
			mcUnits: []mcUnitConfig{
				{
					name: "hugepages-allocation-2048kB-NUMA0.service",
					contents: `Environment=HUGEPAGES_COUNT=4
									   Environment=HUGEPAGES_SIZE=2048kB
									   Environment=NUMA_NODE=0`,
				},
			},
		},
		// Numas mismatch: numa 1 does not exist in mc, but node has no hugepages.
		{
			nodeHugePagesCmdOutput: `/host/sys/devices/system/node/node0/hugepages/hugepages-2048kB/nr_hugepages count:4
									 /host/sys/devices/system/node/node1/hugepages/hugepages-1048576kB/nr_hugepages count:0`,
			mcUnits: []mcUnitConfig{
				{
					name: "hugepages-allocation-2048kB-NUMA0.service",
					contents: `Environment=HUGEPAGES_COUNT=4
									   Environment=HUGEPAGES_SIZE=2048kB
									   Environment=NUMA_NODE=0`,
				},
			},
		},
	}

	// instantiate the fakeClient so we can mock the output from each command to get the node's hugepages files.
	client := fakeK8sClient{}

	for _, tc := range testCases {
		client.execCommandFunctionMocker = func() (string, string, error) {
			return tc.nodeHugePagesCmdOutput, "", nil
		}

		hpTester, _ := NewTester(
			&provider.Node{
				Data: &corev1.Node{ObjectMeta: metav1.ObjectMeta{Name: "node1", Namespace: "ns1"}},
				Mc:   getMcFromUnits(tc.mcUnits),
			},
			fakeDebugPod,
			&client,
		)

		assert.Nil(t, hpTester.Run())
	}
}

func TestNegativeMachineConfigSystemdHugepages(t *testing.T) {
	// helper pod, so the hugepages struct doesnt crash when accessing the debug container.
	fakeDebugPod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{Name: "pod1", Namespace: "ns1"},
		Spec:       corev1.PodSpec{Containers: []corev1.Container{{Name: "container1"}}},
	}

	// helper struct to hold the mc units for each tc
	type mcUnitConfig struct {
		name     string
		contents string
	}

	// helper function to get a provider.MachineConfig object from an slice of units to
	// be used on each tc.
	getMcFromUnits := func(units []mcUnitConfig) provider.MachineConfig {
		mc := provider.MachineConfig{MachineConfig: &mcv1.MachineConfig{}}

		for _, unit := range units {
			mc.Config.Systemd.Units = append(mc.Config.Systemd.Units, struct {
				Contents string "json:\"contents\""
				Name     string "json:\"name\""
			}{
				Name:     unit.name,
				Contents: unit.contents,
			})
		}
		return mc
	}

	testCases := []struct {
		nodeHugePagesCmdOutput string
		mcUnits                []mcUnitConfig
		expectedErrorMsg       string
	}{
		// Numas mismatch: numa 1 has hugepages > 0 but mc does not have any config for that numa.
		{
			nodeHugePagesCmdOutput: `/host/sys/devices/system/node/node0/hugepages/hugepages-2048kB/nr_hugepages count:4
									 /host/sys/devices/system/node/node1/hugepages/hugepages-2048kB/nr_hugepages count:8`,
			mcUnits: []mcUnitConfig{
				{
					name: "hugepages-allocation-2048kB-NUMA0.service",
					contents: `Environment=HUGEPAGES_COUNT=4
							   Environment=HUGEPAGES_SIZE=2048kB
							   Environment=NUMA_NODE=0`,
				},
			},
			expectedErrorMsg: "failed to compare machineConfig systemd's unit hugepages config with node values, err: node's numa 1 hugepages config does not exist in node's machineconfig",
		},
		// Numas mismatch: mc numa id does not exist in the node.
		{
			nodeHugePagesCmdOutput: `/host/sys/devices/system/node/node0/hugepages/hugepages-2048kB/nr_hugepages count:4`,
			mcUnits: []mcUnitConfig{
				{
					name: "hugepages-allocation-2048kB-NUMA0.service",
					contents: `Environment=HUGEPAGES_COUNT=4
							   Environment=HUGEPAGES_SIZE=2048kB
							   Environment=NUMA_NODE=0`,
				},
				{
					name: "hugepages-allocation-1048576kB-NUMA1.service",
					contents: `Environment=HUGEPAGES_COUNT=8
							   Environment=HUGEPAGES_SIZE=1048576kB
							   Environment=NUMA_NODE=1`,
				},
			},
			expectedErrorMsg: "failed to compare machineConfig systemd's unit hugepages config with node values, err: node does not have numa id 1 found in the machine config",
		},
		// Size mismatch: Node's hp size (1GB) does not exist in mc info, but the count is not zero.
		{
			nodeHugePagesCmdOutput: `/host/sys/devices/system/node/node0/hugepages/hugepages-2048kB/nr_hugepages count:4
									 /host/sys/devices/system/node/node0/hugepages/hugepages-1048576kB/nr_hugepages count:8`,
			mcUnits: []mcUnitConfig{
				{
					name: "hugepages-allocation-2048kB-NUMA0.service",
					contents: `Environment=HUGEPAGES_COUNT=4
							   Environment=HUGEPAGES_SIZE=2048kB
							   Environment=NUMA_NODE=0`,
				},
			},
			expectedErrorMsg: "failed to compare machineConfig systemd's unit hugepages config with node values, err: node's numa 0 hugepages size=1048576 does not appear in MC, but the count is not zero (8)",
		},
		// Size mismatch: mc size (1GB) does not exist in node.
		{
			nodeHugePagesCmdOutput: `/host/sys/devices/system/node/node0/hugepages/hugepages-2048kB/nr_hugepages count:4`,
			mcUnits: []mcUnitConfig{
				{
					name: "hugepages-allocation-2048kB-NUMA0.service",
					contents: `Environment=HUGEPAGES_COUNT=4
							   Environment=HUGEPAGES_SIZE=2048kB
							   Environment=NUMA_NODE=0`,
				},
				{
					name: "hugepages-allocation-1048576kB-NUMA0.service",
					contents: `Environment=HUGEPAGES_COUNT=8
							   Environment=HUGEPAGES_SIZE=1048576kB
							   Environment=NUMA_NODE=0`,
				},
			},
			expectedErrorMsg: "failed to compare machineConfig systemd's unit hugepages config with node values, err: node's numa id 0 does not have size 1048576 found in the machine config",
		},
		// Count mismatch: one numa, one size only.
		{
			nodeHugePagesCmdOutput: `/host/sys/devices/system/node/node0/hugepages/hugepages-2048kB/nr_hugepages count:4`,
			mcUnits: []mcUnitConfig{
				{
					name: "hugepages-allocation-2048kB-NUMA0.service",
					contents: `Environment=HUGEPAGES_COUNT=8
							   Environment=HUGEPAGES_SIZE=2048kB
							   Environment=NUMA_NODE=0`,
				},
			},
			expectedErrorMsg: "failed to compare machineConfig systemd's unit hugepages config with node values, err: mc numa=0, hugepages count:8, size:2048 does not match node ones=4",
		},
		// Count mismatch: two numas two sizes. The count for size 1GB does not match in numa 1.
		{
			nodeHugePagesCmdOutput: `/host/sys/devices/system/node/node0/hugepages/hugepages-2048kB/nr_hugepages count:256
									 /host/sys/devices/system/node/node0/hugepages/hugepages-1048576kB/nr_hugepages count:8
									 /host/sys/devices/system/node/node1/hugepages/hugepages-2048kB/nr_hugepages count:256
									 /host/sys/devices/system/node/node1/hugepages/hugepages-1048576kB/nr_hugepages count:9`,
			mcUnits: []mcUnitConfig{
				{
					name: "hugepages-allocation-2048kB-NUMA0.service",
					contents: `Environment=HUGEPAGES_COUNT=256
							   Environment=HUGEPAGES_SIZE=2048kB
							   Environment=NUMA_NODE=0`,
				},
				{
					name: "hugepages-allocation-2048kB-NUMA1.service",
					contents: `Environment=HUGEPAGES_COUNT=256
							   Environment=HUGEPAGES_SIZE=2048kB
							   Environment=NUMA_NODE=1`,
				},
				{
					name: "hugepages-allocation-1048576kB-NUMA0.service",
					contents: `Environment=HUGEPAGES_COUNT=8
							   Environment=HUGEPAGES_SIZE=1048576kB
							   Environment=NUMA_NODE=0`,
				},
				{
					name: "hugepages-allocation-1048576kB-NUMA1.service",
					contents: `Environment=HUGEPAGES_COUNT=8
							   Environment=HUGEPAGES_SIZE=1048576kB
							   Environment=NUMA_NODE=1`,
				},
			},
			expectedErrorMsg: "failed to compare machineConfig systemd's unit hugepages config with node values, err: mc numa=1, hugepages count:8, size:1048576 does not match node ones=9",
		},
	}

	// instantiate the fakeClient so we can mock the output from each command to get the node's hugepages files.
	client := fakeK8sClient{}

	for _, tc := range testCases {
		client.execCommandFunctionMocker = func() (string, string, error) {
			return tc.nodeHugePagesCmdOutput, "", nil
		}

		hpTester, _ := NewTester(
			&provider.Node{
				Data: &corev1.Node{ObjectMeta: metav1.ObjectMeta{Name: "node1", Namespace: "ns1"}},
				Mc:   getMcFromUnits(tc.mcUnits),
			},
			fakeDebugPod,
			&client,
		)

		assert.Equal(t, errors.New(tc.expectedErrorMsg), hpTester.Run())
	}
}

func TestPositiveMachineConfigKernelArgsHugepages(t *testing.T) {
	// helper pod, so the hugepages test won't crash when accessing the debug container.
	fakeDebugPod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{Name: "pod1", Namespace: "ns1"},
		Spec:       corev1.PodSpec{Containers: []corev1.Container{{Name: "container1"}}},
	}

	// helper function to get a provider.MachineConfig object from a kernelArguments slice
	getMcFromKernelArgs := func(kernelArgs []string) provider.MachineConfig {
		return provider.MachineConfig{MachineConfig: &mcv1.MachineConfig{Spec: mcv1.MachineConfigSpec{KernelArguments: kernelArgs}}}
	}

	testCases := []struct {
		nodeHugePagesCmdOutput string
		mcKernelArgs           []string
	}{
		// No hugepages info found in kernelArgs, but the node has a file for the default hugepages size for RHEL (2MB) with count 0
		{
			nodeHugePagesCmdOutput: `/host/sys/devices/system/node/node0/hugepages/hugepages-2048kB/nr_hugepages count:0`,
			mcKernelArgs:           []string{},
		},
		// kernelArgs has only default size. Node has a file for the default hugepages size for RHEL (2MB) with count 0
		{
			nodeHugePagesCmdOutput: `/host/sys/devices/system/node/node0/hugepages/hugepages-2048kB/nr_hugepages count:0`,
			mcKernelArgs:           []string{"default_hugepagesz=2M"},
		},
		// kernelArgs has only default size. Node has a file for size 2MB with count 0
		{
			nodeHugePagesCmdOutput: `/host/sys/devices/system/node/node0/hugepages/hugepages-2048kB/nr_hugepages count:0`,
			mcKernelArgs:           []string{"default_hugepagesz=2M"},
		},
		// kernelArgs has 16 hugepages of the default size for RHEL (2MB)
		{
			nodeHugePagesCmdOutput: `/host/sys/devices/system/node/node0/hugepages/hugepages-2048kB/nr_hugepages count:16`,
			mcKernelArgs:           []string{"hugepages=16"},
		},
		// kernelArgs has only a value, which is assumed to be the RHEL's default size (2MB). Node has two numas whose total
		// matches the kernelArgs' one.
		{
			nodeHugePagesCmdOutput: `/host/sys/devices/system/node/node0/hugepages/hugepages-2048kB/nr_hugepages count:8
									 /host/sys/devices/system/node/node1/hugepages/hugepages-2048kB/nr_hugepages count:8`,
			mcKernelArgs: []string{"hugepages=16"},
		},
		// kernelArgs has only a value, which is assumed to be the RHEL's default size (2MB). Node has two numas but only
		// hugepages on the first one, which matches kernelArgs' one.
		{
			nodeHugePagesCmdOutput: `/host/sys/devices/system/node/node0/hugepages/hugepages-2048kB/nr_hugepages count:16
									 /host/sys/devices/system/node/node1/hugepages/hugepages-2048kB/nr_hugepages count:0`,
			mcKernelArgs: []string{"hugepages=16"},
		},
		// kernelArgs has only a value, which is assumed to be the RHEL's default size (2MB). Node has two numas but only
		// hugepages on the second one, which matches kernelArgs' one.
		{
			nodeHugePagesCmdOutput: `/host/sys/devices/system/node/node0/hugepages/hugepages-2048kB/nr_hugepages count:0
									 /host/sys/devices/system/node/node1/hugepages/hugepages-2048kB/nr_hugepages count:16`,
			mcKernelArgs: []string{"hugepages=16"},
		},
		// kernelArgs has one size with value, and node has two numas and two sizes, but only hugepages count for size 1GB.
		// For the other size (2M), the node has no hugepages.
		{
			nodeHugePagesCmdOutput: `/host/sys/devices/system/node/node0/hugepages/hugepages-2048kB/nr_hugepages count:0
									 /host/sys/devices/system/node/node0/hugepages/hugepages-1048576kB/nr_hugepages count:4
									 /host/sys/devices/system/node/node1/hugepages/hugepages-2048kB/nr_hugepages count:0
									 /host/sys/devices/system/node/node1/hugepages/hugepages-1048576kB/nr_hugepages count:4`,
			mcKernelArgs: []string{"hugepagesz=1G", "hugepages=8"},
		},
		// Node has two numas and two sizes, with hugepages count on both. KernelArgs has both sizes with their values.
		{
			nodeHugePagesCmdOutput: `/host/sys/devices/system/node/node0/hugepages/hugepages-2048kB/nr_hugepages count:128
									 /host/sys/devices/system/node/node0/hugepages/hugepages-1048576kB/nr_hugepages count:8
									 /host/sys/devices/system/node/node1/hugepages/hugepages-2048kB/nr_hugepages count:128
									 /host/sys/devices/system/node/node1/hugepages/hugepages-1048576kB/nr_hugepages count:8`,
			mcKernelArgs: []string{"hugepagesz=1G", "hugepages=16", "hugepagesz=2M", "hugepages=256"},
		},
		// Same as before, but kernelArgs' params order inverted
		{
			nodeHugePagesCmdOutput: `/host/sys/devices/system/node/node0/hugepages/hugepages-2048kB/nr_hugepages count:128
									 /host/sys/devices/system/node/node0/hugepages/hugepages-1048576kB/nr_hugepages count:8
									 /host/sys/devices/system/node/node1/hugepages/hugepages-2048kB/nr_hugepages count:128
									 /host/sys/devices/system/node/node1/hugepages/hugepages-1048576kB/nr_hugepages count:8`,
			mcKernelArgs: []string{"hugepagesz=2M", "hugepages=256", "hugepagesz=1G", "hugepages=16"},
		},
		// Node has two numas and two sizes, with hugepages count on the first numa only. The second numa does not have any
		// hugepages for any size.
		{
			nodeHugePagesCmdOutput: `/host/sys/devices/system/node/node0/hugepages/hugepages-2048kB/nr_hugepages count:256
									 /host/sys/devices/system/node/node0/hugepages/hugepages-1048576kB/nr_hugepages count:16
									 /host/sys/devices/system/node/node1/hugepages/hugepages-2048kB/nr_hugepages count:0
									 /host/sys/devices/system/node/node1/hugepages/hugepages-1048576kB/nr_hugepages count:0`,
			mcKernelArgs: []string{"hugepagesz=1G", "hugepages=16", "hugepagesz=2M", "hugepages=256"},
		},
		// Node has two numas and two sizes, with hugepages count on the second numa only. The first numa does not have any
		// hugepages for any size.
		{
			nodeHugePagesCmdOutput: `/host/sys/devices/system/node/node0/hugepages/hugepages-2048kB/nr_hugepages count:0
									 /host/sys/devices/system/node/node0/hugepages/hugepages-1048576kB/nr_hugepages count:0
									 /host/sys/devices/system/node/node1/hugepages/hugepages-2048kB/nr_hugepages count:256
									 /host/sys/devices/system/node/node1/hugepages/hugepages-1048576kB/nr_hugepages count:16`,
			mcKernelArgs: []string{"hugepagesz=1G", "hugepages=16", "hugepagesz=2M", "hugepages=256"},
		},
	}

	// instantiate the fakeClient so we can mock the output from each command to get the node's hugepages files.
	client := fakeK8sClient{}

	for _, tc := range testCases {
		client.execCommandFunctionMocker = func() (string, string, error) {
			return tc.nodeHugePagesCmdOutput, "", nil
		}

		hpTester, _ := NewTester(
			&provider.Node{
				Data: &corev1.Node{ObjectMeta: metav1.ObjectMeta{Name: "node1", Namespace: "ns1"}},
				Mc:   getMcFromKernelArgs(tc.mcKernelArgs)},
			fakeDebugPod,
			&client,
		)

		assert.Nil(t, hpTester.Run())
	}
}

func TestNegativeMachineConfigKernelArgsHugepages(t *testing.T) {
	// helper pod, so the hugepages test won't crash when accessing the debug container.
	fakeDebugPod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{Name: "pod1", Namespace: "ns1"},
		Spec:       corev1.PodSpec{Containers: []corev1.Container{{Name: "container1"}}},
	}

	// helper function to get a provider.MachineConfig object from a kernelArguments slice
	getMcFromKernelArgs := func(kernelArgs []string) provider.MachineConfig {
		return provider.MachineConfig{MachineConfig: &mcv1.MachineConfig{Spec: mcv1.MachineConfigSpec{KernelArguments: kernelArgs}}}
	}

	testCases := []struct {
		nodeHugePagesCmdOutput string
		mcKernelArgs           []string
		expectedErrorMsg       string
	}{
		// No hugepages config in kernelArgs, but the node has non-zero value for 2MB size.
		{
			nodeHugePagesCmdOutput: `/host/sys/devices/system/node/node0/hugepages/hugepages-2048kB/nr_hugepages count:2`,
			mcKernelArgs:           []string{},
			expectedErrorMsg:       "failed to compare machineConfig KernelArguments with node ones, err: total hugepages of size 2048 won't match (node count=2, expected=0)",
		},
		// Count mismatch: kernelArgs value is non-zero, while node's one is zero.
		{
			nodeHugePagesCmdOutput: `/host/sys/devices/system/node/node0/hugepages/hugepages-2048kB/nr_hugepages count:0`,
			mcKernelArgs:           []string{"hugepages=16"},
			expectedErrorMsg:       "failed to compare machineConfig KernelArguments with node ones, err: total hugepages of size 2048 won't match (node count=0, expected=16)",
		},
		// Count mismatch: Node has two numas with non-zero hugepages each but the total won't match.
		{
			nodeHugePagesCmdOutput: `/host/sys/devices/system/node/node0/hugepages/hugepages-2048kB/nr_hugepages count:8
									 /host/sys/devices/system/node/node1/hugepages/hugepages-2048kB/nr_hugepages count:7`,
			mcKernelArgs:     []string{"hugepages=16"},
			expectedErrorMsg: "failed to compare machineConfig KernelArguments with node ones, err: total hugepages of size 2048 won't match (node count=15, expected=16)",
		},
		// Size mismatch: kernelArgs size won't match node's one (count > 0).
		{
			nodeHugePagesCmdOutput: `/host/sys/devices/system/node/node0/hugepages/hugepages-2048kB/nr_hugepages count:8`,
			mcKernelArgs:           []string{"hugepagesz=1G", "hugepages=8"},
			expectedErrorMsg:       "failed to compare machineConfig KernelArguments with node ones, err: node's numa 0 hugepages size=2048 does not appear in kernelArgs, but the count is not zero (8)",
		},
		// Size mismatch: kernelArgs size won't match node's one (count == 0).
		{
			nodeHugePagesCmdOutput: `/host/sys/devices/system/node/node0/hugepages/hugepages-2048kB/nr_hugepages count:0`,
			mcKernelArgs:           []string{"hugepagesz=1G", "hugepages=8"},
			expectedErrorMsg:       "failed to compare machineConfig KernelArguments with node ones, err: node's numa 0 has no hugepages of kernelArgs' size 1048576",
		},
		// Count mismatch: node has two numas and two sizes, with hugepages count on both, but MC's only defines de default size.
		{
			nodeHugePagesCmdOutput: `/host/sys/devices/system/node/node0/hugepages/hugepages-2048kB/nr_hugepages count:0
									 /host/sys/devices/system/node/node0/hugepages/hugepages-1048576kB/nr_hugepages count:4
									 /host/sys/devices/system/node/node1/hugepages/hugepages-2048kB/nr_hugepages count:0
									 /host/sys/devices/system/node/node1/hugepages/hugepages-1048576kB/nr_hugepages count:4`,
			mcKernelArgs:     []string{"default_hugepagesz=1G"},
			expectedErrorMsg: "failed to compare machineConfig KernelArguments with node ones, err: total hugepages of size 1048576 won't match (node count=8, expected=0)",
		},
		// Count mismatch: node has two numas and two sizes, with hugepages count on both. Total for size 2MB won't match kernelArgs.
		{
			nodeHugePagesCmdOutput: `/host/sys/devices/system/node/node0/hugepages/hugepages-2048kB/nr_hugepages count:100
									 /host/sys/devices/system/node/node0/hugepages/hugepages-1048576kB/nr_hugepages count:4
									 /host/sys/devices/system/node/node1/hugepages/hugepages-2048kB/nr_hugepages count:200
									 /host/sys/devices/system/node/node1/hugepages/hugepages-1048576kB/nr_hugepages count:4`,
			mcKernelArgs:     []string{"hugepagesz=1G", "hugepages=8", "hugepagesz=2M", "hugepages=256"},
			expectedErrorMsg: "failed to compare machineConfig KernelArguments with node ones, err: total hugepages of size 2048 won't match (node count=300, expected=256)",
		},
		// Count mismatch: node has two numas and two sizes, with hugepages count on both. Total for size 1GB won't match kernelArgs.
		{
			nodeHugePagesCmdOutput: `/host/sys/devices/system/node/node0/hugepages/hugepages-2048kB/nr_hugepages count:256
									 /host/sys/devices/system/node/node0/hugepages/hugepages-1048576kB/nr_hugepages count:16
									 /host/sys/devices/system/node/node1/hugepages/hugepages-2048kB/nr_hugepages count:0
									 /host/sys/devices/system/node/node1/hugepages/hugepages-1048576kB/nr_hugepages count:0`,
			mcKernelArgs:     []string{"hugepagesz=1G", "hugepages=8", "hugepagesz=2M", "hugepages=256"},
			expectedErrorMsg: "failed to compare machineConfig KernelArguments with node ones, err: total hugepages of size 1048576 won't match (node count=16, expected=8)",
		},
	}

	// instantiate the fakeClient so we can mock the output from each command to get the node's hugepages files.
	client := fakeK8sClient{}

	for _, tc := range testCases {
		client.execCommandFunctionMocker = func() (string, string, error) {
			return tc.nodeHugePagesCmdOutput, "", nil
		}

		hpTester, _ := NewTester(
			&provider.Node{
				Data: &corev1.Node{ObjectMeta: metav1.ObjectMeta{Name: "node1", Namespace: "ns1"}},
				Mc:   getMcFromKernelArgs(tc.mcKernelArgs)},
			fakeDebugPod,
			&client,
		)

		assert.Equal(t, errors.New(tc.expectedErrorMsg), hpTester.Run())
	}
}
