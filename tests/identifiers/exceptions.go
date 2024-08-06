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

package identifiers

/*
	Use this file to store the strings for the various exception processes for each test in our suite.
	There are various tests that have a level of flexibility to their results depending on the situation and some that do not have
	an exception process.
*/

const (

	// Tests with exception processes.
	// TODO: Add more exception processes if/when we encounter more opportunities with partners.
	IsRedHatReleaseExceptionProcess = `Document which containers are not able to meet the RHEL-based container requirement and if/when the base image can be updated.`

	SecConNonRootUserExceptionProcess = `No exceptions - will only be considered under special circumstances. Must identify which container needs access and document why with details.`

	SecConCapabilitiesExceptionProcess = `Identify the pod that is needing special capabilities and document why`

	// Tests that do not have an exception process but have additional insight.
	UnalteredBaseImageExceptionProcess = `Images should not be changed during runtime. There is no exception process for this.`

	// Generic Exception Process Message.
	NoDocumentedProcess = `There is no documented exception process for this.`

	// Generic provide info message.
	ElaborateOnWhyItIsNeeded = `Please elaborate why it's needed and explain how it's used.`

	// No exception process for extended tests.
	NoExceptionProcessForExtendedTests = `No exception needed for optional/extended tests.`

	// No exception process.
	NoExceptions = `No exceptions`

	// affiliated certification exception process.
	AffiliatedCert = NoDocumentedProcess + " " + `A partner can run the Red Hat Best Practices Test Suite before passing other certifications (Container/Operator/HelmChart) but the affiliated certification test cases in the Red Hat Best Practices Test Suite must be re-run once the other certifications have been granted.` //nolint:lll
)
