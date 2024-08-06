// Copyright (C) 2020-2022 Red Hat, Inc.
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

package sysctlconfig

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/test-network-function/cnf-certification-test/internal/clientsholder"
	"github.com/test-network-function/cnf-certification-test/pkg/provider"
)

// Creates a map describing the final sysctl key-value pair out of the results of "sysctl --system".
func parseSysctlSystemOutput(sysctlSystemOutput string) map[string]string {
	retval := make(map[string]string)
	splitConfig := strings.Split(sysctlSystemOutput, "\n")
	for _, line := range splitConfig {
		if strings.HasPrefix(line, "*") {
			continue
		}

		keyValRegexp := regexp.MustCompile(`(\S+)\s*=\s*(\S+)`) // A line is of the form "kernel.yama.ptrace_scope = 0"
		if !keyValRegexp.MatchString(line) {
			continue
		}
		regexResults := keyValRegexp.FindStringSubmatch(line)
		key := regexResults[1]
		val := regexResults[2]
		retval[key] = val
	}
	return retval
}

func GetSysctlSettings(env *provider.TestEnvironment, nodeName string) (map[string]string, error) {
	const (
		sysctlCommand = "chroot /host sysctl --system"
	)

	o := clientsholder.GetClientsHolder()
	ctx := clientsholder.NewContext(env.DebugPods[nodeName].Namespace, env.DebugPods[nodeName].Name, env.DebugPods[nodeName].Spec.Containers[0].Name)

	outStr, errStr, err := o.ExecCommandContainer(ctx, sysctlCommand)
	if err != nil || errStr != "" {
		return nil, fmt.Errorf("failed to execute command %s in debug pod %s, err=%s, stderr=%s", sysctlCommand,
			env.DebugPods[nodeName], err, errStr)
	}

	return parseSysctlSystemOutput(outStr), nil
}
