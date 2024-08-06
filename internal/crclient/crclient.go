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

package crclient

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/test-network-function/cnf-certification-test/internal/clientsholder"
	"github.com/test-network-function/cnf-certification-test/internal/log"
	"github.com/test-network-function/cnf-certification-test/pkg/provider"
)

const PsRegex = `(?m)^(\d+?)\s+?(\d+?)\s+?(\d+?)\s+?(.*?)$`

type Process struct {
	PidNs, Pid, PPid int
	Args             string
}

const (
	DevNull          = " 2>/dev/null"
	DockerInspectPID = "chroot /host docker inspect -f '{{.State.Pid}}' "
)

func (p *Process) String() string {
	return fmt.Sprintf("cmd: %s, pid: %d, ppid: %d, pidNs: %d", p.Args, p.Pid, p.PPid, p.PidNs)
}

// Helper function to create the clientsholder.Context of the first container of the debug pod
// that runs in the give node. This context is usually needed to run shell commands that get
// information from a node where a pod/container under test is running.
func GetNodeDebugPodContext(node string, env *provider.TestEnvironment) (clientsholder.Context, error) {
	debugPod := env.DebugPods[node]
	if debugPod == nil {
		return clientsholder.Context{}, fmt.Errorf("debug pod not found on node %s", node)
	}

	return clientsholder.NewContext(debugPod.Namespace, debugPod.Name, debugPod.Spec.Containers[0].Name), nil
}

func GetPidFromContainer(cut *provider.Container, ctx clientsholder.Context) (int, error) {
	var pidCmd string

	switch cut.Runtime {
	case "docker":
		pidCmd = DockerInspectPID + cut.UID + DevNull
	case "docker-pullable":
		pidCmd = DockerInspectPID + cut.UID + DevNull
	case "cri-o", "containerd":
		pidCmd = "chroot /host crictl inspect --output go-template --template '{{.info.pid}}' " + cut.UID + DevNull
	default:
		log.Debug("Container runtime %s not supported yet for this test, skipping", cut.Runtime)
		return 0, fmt.Errorf("container runtime %s not supported", cut.Runtime)
	}

	ch := clientsholder.GetClientsHolder()
	outStr, errStr, err := ch.ExecCommandContainer(ctx, pidCmd)
	if err != nil {
		return 0, fmt.Errorf("cannot execute command: \" %s \"  on %s err:%s", pidCmd, cut, err)
	}
	if errStr != "" {
		return 0, fmt.Errorf("cmd: \" %s \" on %s returned %s", pidCmd, cut, errStr)
	}

	return strconv.Atoi(strings.TrimSuffix(outStr, "\n"))
}

// To get the pid namespace of the container.
func GetContainerPidNamespace(testContainer *provider.Container, env *provider.TestEnvironment) (string, error) {
	// Get the container pid.
	ocpContext, err := GetNodeDebugPodContext(testContainer.NodeName, env)
	if err != nil {
		return "", fmt.Errorf("failed to get debug pod's context for container %s: %v", testContainer, err)
	}

	pid, err := GetPidFromContainer(testContainer, ocpContext)
	if err != nil {
		return "", fmt.Errorf("unable to get container process id due to: %v", err)
	}
	log.Debug("Obtained process id for %s is %d", testContainer, pid)

	command := fmt.Sprintf("lsns -p %d -t pid -n", pid)
	stdout, stderr, err := clientsholder.GetClientsHolder().ExecCommandContainer(ocpContext, command)
	if err != nil || stderr != "" {
		return "", fmt.Errorf("unable to run nsenter due to : %v", err)
	}

	return strings.Fields(stdout)[0], nil
}

func GetContainerProcesses(container *provider.Container, env *provider.TestEnvironment) ([]*Process, error) {
	pidNs, err := GetContainerPidNamespace(container, env)
	if err != nil {
		return nil, fmt.Errorf("could not get the containers' pid namespace, err: %v", err)
	}

	return GetPidsFromPidNamespace(pidNs, container)
}

// ExecCommandContainerNSEnter executes a command in the specified container namespace using nsenter.
func ExecCommandContainerNSEnter(command string,
	aContainer *provider.Container) (outStr, errStr string, err error) {
	env := provider.GetTestEnvironment()
	ctx, err := GetNodeDebugPodContext(aContainer.NodeName, &env)
	if err != nil {
		return "", "", fmt.Errorf("failed to get debug pod's context for container %s: %v", aContainer, err)
	}

	ch := clientsholder.GetClientsHolder()

	// Get the container PID to build the nsenter command.
	containerPid, err := GetPidFromContainer(aContainer, ctx)
	if err != nil {
		return "", "", fmt.Errorf("cannot get PID from: %s, err: %v", aContainer, err)
	}

	// Add the container PID and the specific command to run with nsenter.
	nsenterCommand := "nsenter -t " + strconv.Itoa(containerPid) + " -n " + command

	// Run the nsenter command on the debug pod.
	outStr, errStr, err = ch.ExecCommandContainer(ctx, nsenterCommand)
	if err != nil {
		return "", "", fmt.Errorf("cannot execute command: \" %s \"  on %s err:%s", command, aContainer, err)
	}

	return outStr, errStr, err
}

func GetPidsFromPidNamespace(pidNamespace string, container *provider.Container) (p []*Process, err error) {
	const command = "trap \"\" SIGURG ; ps -e -o pidns,pid,ppid,args"
	env := provider.GetTestEnvironment()
	ctx, err := GetNodeDebugPodContext(container.NodeName, &env)
	if err != nil {
		return nil, fmt.Errorf("failed to get debug pod's context for container %s: %v", container, err)
	}

	stdout, stderr, err := clientsholder.GetClientsHolder().ExecCommandContainer(ctx, command)
	if err != nil || stderr != "" {
		return nil, fmt.Errorf("command %q failed to run in debug pod=%s (node=%s): %v", command, ctx.GetPodName(), container.NodeName, err)
	}

	re := regexp.MustCompile(PsRegex)
	matches := re.FindAllStringSubmatch(stdout, -1)
	// If we do not find a successful log, we fail.
	for _, v := range matches {
		// Matching only the right PidNs.
		if pidNamespace != v[1] {
			continue
		}
		aPidNs, err := strconv.Atoi(v[1])
		if err != nil {
			log.Error("could not convert string %s to integer, err=%s", v[1], err)
			continue
		}
		aPid, err := strconv.Atoi(v[2])
		if err != nil {
			log.Error("could not convert string %s to integer, err=%s", v[2], err)
			continue
		}
		aPPid, err := strconv.Atoi(v[3])
		if err != nil {
			log.Error("could not convert string %s to integer, err=%s", v[3], err)
			continue
		}
		p = append(p, &Process{PidNs: aPidNs, Pid: aPid, Args: v[4], PPid: aPPid})
	}
	return p, nil
}
