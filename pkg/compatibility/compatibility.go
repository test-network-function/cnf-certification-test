// Copyright (C) 2022-2023 Red Hat, Inc.
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

package compatibility

import (
	"strings"
	"time"

	gv "github.com/hashicorp/go-version"
	"github.com/test-network-function/cnf-certification-test/internal/log"
	"github.com/test-network-function/cnf-certification-test/pkg/stringhelper"
)

/* Notes for this package
* Refer to this document for more information about OCP compatibility: https://access.redhat.com/support/policy/updates/openshift

* This module will help compare the running OCP version against a matrix of end of life dates.
 */

const (
	// OCP Lifecycle Statuses.
	OCPStatusGA      = "generally-available"
	OCPStatusMS      = "maintenance-support"
	OCPStatusEOL     = "end-of-life"
	OCPStatusUnknown = "unknown"
	OCPStatusPreGA   = "pre-general-availability"
)

type VersionInfo struct {
	GADate  time.Time // General Availability Date
	FSEDate time.Time // Full Support Ends Date
	MSEDate time.Time // Maintenance Support Ends Date

	MinRHCOSVersion      string   // Minimum RHCOS Version supported
	RHELVersionsAccepted []string // Contains either specific versions or a minimum version eg. "7.9 or later" or "7.9 and 8.4"
}

var (
	ocpBetaVersions   = []string{"4.13", "4.14", "4.15", "4.16", "4.17"}
	ocpLifeCycleDates = map[string]VersionInfo{
		// TODO: Adjust all of these periodically to make sure they are up to date with the lifecycle
		// update documentation.

		// Full Support
		"4.15": {
			GADate:  time.Date(2024, 2, 27, 0, 0, 0, 0, time.UTC), // February 27, 2024
			FSEDate: time.Date(2025, 4, 27, 0, 0, 0, 0, time.UTC), // April 27, 2025
			MSEDate: time.Date(2025, 8, 27, 0, 0, 0, 0, time.UTC), // August 27, 2025
			// Note: FSEDate (Release of 4.16 + 3 months) is currently a "guess".  Update when available.

			// OS Compatibility
			MinRHCOSVersion:      "4.15",
			RHELVersionsAccepted: []string{"8.4", "8.5"},
		},
		"4.14": {
			GADate:  time.Date(2023, 10, 31, 0, 0, 0, 0, time.UTC), // October 31, 2023
			FSEDate: time.Date(2024, 5, 27, 0, 0, 0, 0, time.UTC),  // May 27, 2024
			MSEDate: time.Date(2025, 5, 1, 0, 0, 0, 0, time.UTC),   // May 1, 2025

			// OS Compatibility
			MinRHCOSVersion:      "4.14",
			RHELVersionsAccepted: []string{"8.4", "8.5"},
		},

		// Maintenance Support
		"4.13": {
			GADate:  time.Date(2023, 5, 17, 0, 0, 0, 0, time.UTC),  // May 17, 2023
			FSEDate: time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC),  // January 31, 2024
			MSEDate: time.Date(2024, 11, 17, 0, 0, 0, 0, time.UTC), // November 17, 2024

			// OS Compatibility
			MinRHCOSVersion:      "4.13",
			RHELVersionsAccepted: []string{"8.4", "8.5"},
		},
		"4.12": {
			GADate:  time.Date(2023, 1, 17, 0, 0, 0, 0, time.UTC), // January 17, 2023
			FSEDate: time.Date(2023, 8, 17, 0, 0, 0, 0, time.UTC), // August 17, 2023
			MSEDate: time.Date(2024, 7, 17, 0, 0, 0, 0, time.UTC), // July 17, 2024

			// OS Compatibility
			MinRHCOSVersion:      "4.12",
			RHELVersionsAccepted: []string{"8.4", "8.5"},
		},
		"4.11": {
			GADate:  time.Date(2022, 8, 10, 0, 0, 0, 0, time.UTC), // August 10, 2022
			FSEDate: time.Date(2023, 4, 17, 0, 0, 0, 0, time.UTC), // April 17, 2023
			MSEDate: time.Date(2024, 2, 10, 0, 0, 0, 0, time.UTC), // February 10, 2024

			// OS Compatibility
			MinRHCOSVersion:      "4.11",
			RHELVersionsAccepted: []string{"8.4", "8.5"},
		},
		"4.10": {
			GADate:  time.Date(2022, 3, 10, 0, 0, 0, 0, time.UTC),  // March 10, 2022
			FSEDate: time.Date(2022, 11, 10, 0, 0, 0, 0, time.UTC), // November 10, 2022
			MSEDate: time.Date(2023, 9, 10, 0, 0, 0, 0, time.UTC),  // September 10, 2023

			// OS Compatibility
			MinRHCOSVersion:      "4.10",
			RHELVersionsAccepted: []string{"8.4", "8.5"},
		},

		// End of life
		"4.9": {
			GADate:  time.Date(2021, 10, 18, 0, 0, 0, 0, time.UTC), // October 18, 2021
			FSEDate: time.Date(2022, 6, 10, 0, 0, 0, 0, time.UTC),  // June 10, 2022
			MSEDate: time.Date(2023, 4, 18, 0, 0, 0, 0, time.UTC),  // April 18, 2023

			// OS Compatibility
			MinRHCOSVersion:      "4.9",
			RHELVersionsAccepted: []string{"7.9", "8.4"},
		},
		"4.8": {
			GADate:  time.Date(2021, 7, 27, 0, 0, 0, 0, time.UTC), // July 27, 2021
			FSEDate: time.Date(2022, 1, 27, 0, 0, 0, 0, time.UTC), // January 27, 2022
			MSEDate: time.Date(2023, 1, 27, 0, 0, 0, 0, time.UTC), // January 27, 2023

			// OS Compatibility
			MinRHCOSVersion:      "4.8",
			RHELVersionsAccepted: []string{"7.9"},
		},
		"4.7": {
			GADate:  time.Date(2021, 2, 24, 0, 0, 0, 0, time.UTC),  // February 24, 2021
			FSEDate: time.Date(2021, 10, 27, 0, 0, 0, 0, time.UTC), // October 27, 2021
			MSEDate: time.Date(2022, 8, 24, 0, 0, 0, 0, time.UTC),  // August 24, 2022

			// OS Compatibility
			MinRHCOSVersion:      "4.7",
			RHELVersionsAccepted: []string{"7.9"},
		},
		"4.6": {
			GADate:  time.Date(2020, 10, 27, 0, 0, 0, 0, time.UTC), // October 27, 2020
			FSEDate: time.Date(2021, 3, 24, 0, 0, 0, 0, time.UTC),  // March 24, 2021
			MSEDate: time.Date(2021, 10, 18, 0, 0, 0, 0, time.UTC), // October 18, 2022

			// OS Compatibility
			MinRHCOSVersion:      "4.6",
			RHELVersionsAccepted: []string{"7.9"},
		},
		"4.5": {
			GADate:  time.Date(2020, 7, 13, 0, 0, 0, 0, time.UTC),  // July 13, 2020
			FSEDate: time.Date(2020, 11, 27, 0, 0, 0, 0, time.UTC), // November 27, 2020
			MSEDate: time.Date(2021, 7, 27, 0, 0, 0, 0, time.UTC),  // July 27, 2021

			// OS Compatibility
			MinRHCOSVersion:      "4.5",
			RHELVersionsAccepted: []string{"7.8", "7.9"},
		},
		"4.4": {
			GADate:  time.Date(2020, 5, 5, 0, 0, 0, 0, time.UTC),  // May 5, 2020
			FSEDate: time.Date(2020, 8, 13, 0, 0, 0, 0, time.UTC), // August 13, 2020
			MSEDate: time.Date(2021, 2, 24, 0, 0, 0, 0, time.UTC), // February 24, 2021

			// OS Compatibility
			MinRHCOSVersion:      "4.4",
			RHELVersionsAccepted: []string{"7.6"},
		},
		"4.3": {
			GADate:  time.Date(2020, 1, 23, 0, 0, 0, 0, time.UTC),  // January 23, 2020
			FSEDate: time.Date(2020, 6, 5, 0, 0, 0, 0, time.UTC),   // June 5, 2020
			MSEDate: time.Date(2020, 10, 27, 0, 0, 0, 0, time.UTC), // October 27, 2020

			// OS Compatibility
			MinRHCOSVersion:      "4.3",
			RHELVersionsAccepted: []string{"7.6"},
		},
		"4.2": {
			GADate:  time.Date(2019, 10, 16, 0, 0, 0, 0, time.UTC), // October 16, 2019
			FSEDate: time.Date(2020, 2, 23, 0, 0, 0, 0, time.UTC),  // February 23, 2020
			MSEDate: time.Date(2020, 7, 13, 0, 0, 0, 0, time.UTC),  // July 13, 2020

			// OS Compatibility
			MinRHCOSVersion:      "4.2",
			RHELVersionsAccepted: []string{"7.6"},
		},
		"4.1": {
			GADate:  time.Date(2019, 6, 4, 0, 0, 0, 0, time.UTC),   // June 4, 2019
			FSEDate: time.Date(2019, 11, 16, 0, 0, 0, 0, time.UTC), // November 16, 2019
			MSEDate: time.Date(2020, 5, 5, 0, 0, 0, 0, time.UTC),   // May 5, 2020

			// OS Compatibility
			MinRHCOSVersion:      "4.1",
			RHELVersionsAccepted: []string{"7.6"},
		},
	}
)

func GetLifeCycleDates() map[string]VersionInfo {
	return ocpLifeCycleDates
}

func BetaRHCOSVersionsFoundToMatch(machineVersion, ocpVersion string) bool {
	ocpVersion = FindMajorMinor(ocpVersion)
	machineVersion = FindMajorMinor(machineVersion)

	// Check if the versions exist in the beta list
	if !stringhelper.StringInSlice(ocpBetaVersions, ocpVersion, false) || !stringhelper.StringInSlice(ocpBetaVersions, machineVersion, false) {
		return false
	}

	// Check if the versions match
	return ocpVersion == machineVersion
}

func IsRHELCompatible(machineVersion, ocpVersion string) bool {
	if machineVersion == "" || ocpVersion == "" {
		return false
	}

	lifecycleInfo := GetLifeCycleDates()
	if entry, ok := lifecycleInfo[ocpVersion]; ok {
		if len(entry.RHELVersionsAccepted) >= 2 { //nolint:mnd
			// Need to be a specific major.minor version
			for _, v := range entry.RHELVersionsAccepted {
				if v == machineVersion {
					return true
				}
			}
		} else {
			// Collect the machine version and the entry version
			mv, _ := gv.NewVersion(machineVersion)
			ev, _ := gv.NewVersion(entry.RHELVersionsAccepted[0])

			// If the machine version >= the entry version
			return mv.GreaterThanOrEqual(ev)
		}
	}

	return false
}

func FindMajorMinor(version string) string {
	splitVersion := strings.Split(version, ".")
	return splitVersion[0] + "." + splitVersion[1]
}

func IsRHCOSCompatible(machineVersion, ocpVersion string) bool {
	if machineVersion == "" || ocpVersion == "" {
		return false
	}

	// Exception for beta versions
	if BetaRHCOSVersionsFoundToMatch(machineVersion, ocpVersion) {
		return true
	}

	// Split the incoming version on the "." and make sure we are only looking at major.minor.
	ocpVersion = FindMajorMinor(ocpVersion)

	lifecycleInfo := GetLifeCycleDates()
	if entry, ok := lifecycleInfo[ocpVersion]; ok {
		// Collect the machine version and the entry version
		mv, err := gv.NewVersion(machineVersion)
		if err != nil {
			log.Error("Error parsing machineVersion: %s err: %v", machineVersion, err)
			return false
		}
		ev, err := gv.NewVersion(entry.MinRHCOSVersion)
		if err != nil {
			log.Error("Error parsing MinRHCOSVersion: %s err: %v", entry.MinRHCOSVersion, err)
			return false
		}

		// If the machine version >= the entry version
		return mv.GreaterThanOrEqual(ev)
	}

	return false
}

func DetermineOCPStatus(version string, date time.Time) string {
	// Safeguard against empty values being passed in
	if version == "" || date.IsZero() {
		return OCPStatusUnknown
	}

	// Split the incoming version on the "." and make sure we are only looking at major.minor.
	splitVersion := strings.Split(version, ".")
	version = splitVersion[0] + "." + splitVersion[1]

	// Check if the version exists in our local map
	lifecycleDates := GetLifeCycleDates()
	if entry, ok := lifecycleDates[version]; ok {
		// Safeguard against the latest versions not having a date set for FSEDate set.
		// See the OpenShift lifecycle website link (above) for more details on this.
		if entry.FSEDate.IsZero() {
			entry.FSEDate = entry.MSEDate
		}

		// Pre-GA
		if date.Before(entry.GADate) {
			return OCPStatusPreGA
		}
		// Generally Available
		if date.Equal(entry.GADate) || date.After(entry.GADate) && date.Before(entry.FSEDate) {
			return OCPStatusGA
		}
		// Maintenance Support
		if date.Equal(entry.FSEDate) || (date.After(entry.FSEDate) && date.Before(entry.MSEDate)) {
			return OCPStatusMS
		}
		// End of Life
		if date.Equal(entry.MSEDate) || date.After(entry.MSEDate) {
			return OCPStatusEOL
		}
	}

	return OCPStatusUnknown
}

/*

Note:
You must use RHCOS machines for the control plane, and you can use either RHCOS or RHEL for compute machines.

Compatibility information gathered from the release note pages of each release:

4.11
https://docs.openshift.com/container-platform/4.11/release_notes/ocp-4-11-release-notes.html
OpenShift Container Platform 4.11 is supported on Red Hat Enterprise Linux (RHEL) 8.4 and 8.5, as well as on Red Hat Enterprise Linux CoreOS (RHCOS) 4.11.

4.10
https://docs.openshift.com/container-platform/4.10/release_notes/ocp-4-10-release-notes.html
OpenShift Container Platform 4.10 is supported on Red Hat Enterprise Linux (RHEL) 8.4 and 8.5, as well as on Red Hat Enterprise Linux CoreOS (RHCOS) 4.10.

4.9
https://docs.openshift.com/container-platform/4.9/release_notes/ocp-4-9-release-notes.html
OpenShift Container Platform 4.9 is supported on Red Hat Enterprise Linux (RHEL) 7.9 and 8.4, as well as on Red Hat Enterprise Linux CoreOS (RHCOS) 4.9.

4.8
https://docs.openshift.com/container-platform/4.8/release_notes/ocp-4-8-release-notes.html
OpenShift Container Platform 4.8 is supported on Red Hat Enterprise Linux (RHEL) 7.9 or later, as well as on Red Hat Enterprise Linux CoreOS (RHCOS) 4.8.

4.7
https://docs.openshift.com/container-platform/4.7/release_notes/ocp-4-7-release-notes.html
OpenShift Container Platform 4.7 is supported on Red Hat Enterprise Linux (RHEL) 7.9 or later, as well as Red Hat Enterprise Linux CoreOS (RHCOS) 4.7.

4.6
https://docs.openshift.com/container-platform/4.6/release_notes/ocp-4-6-release-notes.html
OpenShift Container Platform 4.6 is supported on Red Hat Enterprise Linux 7.9 or later, as well as Red Hat Enterprise Linux CoreOS (RHCOS) 4.6.

4.5
https://docs.openshift.com/container-platform/4.5/release_notes/ocp-4-5-release-notes.html
OpenShift Container Platform 4.5 is supported on RHEL 7, version 7.7 or 7.8, as well as Red Hat Enterprise Linux CoreOS (RHCOS) 4.5.

4.4
https://docs.openshift.com/container-platform/4.4/release_notes/ocp-4-4-release-notes.html
OpenShift Container Platform 4.4 is supported on Red Hat Enterprise Linux 7.6 or later, as well as Red Hat Enterprise Linux CoreOS (RHCOS) 4.4.

4.3
https://docs.openshift.com/container-platform/4.3/release_notes/ocp-4-3-release-notes.html
OpenShift Container Platform 4.3 is supported on Red Hat Enterprise Linux 7.6 or later, as well as Red Hat Enterprise Linux CoreOS 4.3.

4.2
https://docs.openshift.com/container-platform/4.2/release_notes/ocp-4-2-release-notes.html
OpenShift Container Platform 4.2 is supported on Red Hat Enterprise Linux 7.6 and later, as well as Red Hat Enterprise Linux CoreOS 4.2.

4.1
https://docs.openshift.com/container-platform/4.1/release_notes/ocp-4-1-release-notes.html
OpenShift Container Platform 4.1 is supported on Red Hat Enterprise Linux 7.6 and later, as well as Red Hat Enterprise Linux CoreOS 4.1.

*/
