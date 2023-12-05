// Copyright (C) 2020-2023 Red Hat, Inc.
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

package flags

import (
	"flag"
	"time"
)

const (
	NoLabelsExpr           = "none"
	labelsFlagName         = "label-filter"
	labelsFlagDefaultValue = "common"

	labelsFlagUsage = "--label-filter <expression>  e.g. --label-filter 'access-control && !access-control-sys-admin-capability'"

	timeoutFlagName         = "timeout"
	TimeoutFlagDefaultvalue = 24 * time.Hour

	timeoutFlagUsage = "--timeout <time>  e.g. --timeout 30m  or -timeout 1h30m"

	listFlagName         = "list"
	listFlagDefaultValue = false

	listFlagUsage = "--list Shows all the available checks/tests. Can be filtered with --label-filter."

	serverModeFlagName         = "serverMode"
	serverModeFlagDefaultValue = false

	serverModeFlagUsage = "--serverMode or -serverMode runs in web server mode."
)

var (
	// labelsFlag holds the labels expression to filter the checks to run.
	LabelsFlag     *string
	TimeoutFlag    *string
	ListFlag       *bool
	ServerModeFlag *bool
)

func InitFlags() {
	LabelsFlag = flag.String(labelsFlagName, labelsFlagDefaultValue, labelsFlagUsage)
	TimeoutFlag = flag.String(timeoutFlagName, TimeoutFlagDefaultvalue.String(), timeoutFlagUsage)
	ListFlag = flag.Bool(listFlagName, listFlagDefaultValue, listFlagUsage)
	ServerModeFlag = flag.Bool(serverModeFlagName, serverModeFlagDefaultValue, serverModeFlagUsage)

	flag.Parse()
	if *LabelsFlag == "" {
		*LabelsFlag = NoLabelsExpr
	}
}
