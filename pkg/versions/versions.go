package versions

import (
	"regexp"

	"github.com/Masterminds/semver"
)

var (
	// GitCommit is the latest commit in the current git branch.
	GitCommit string
	// GitRelease is the list of tags (if any) applied to the latest commit
	// in the current branch.
	GitRelease string
	// GitPreviousRelease is the last release at the date of the latest commit
	// in the current branch.
	GitPreviousRelease string
	// GitDisplayRelease is a string used to hold the text to display
	// the version on screen and in the claim file.
	GitDisplayRelease string
	// ClaimFormat is the current version for the claim file format to be produced by the TNF test suite.
	// A client decoding this claim file must support decoding its specific version.
	ClaimFormatVersion string
)

// getGitVersion returns the git display version: the latest previously released
// build in case this build is not released. Otherwise display the build version.
func GitVersion() string {
	if GitRelease == "" {
		GitDisplayRelease = "Unreleased build post " + GitPreviousRelease
	} else {
		GitDisplayRelease = GitRelease
	}

	return GitDisplayRelease + " (" + GitCommit + ")"
}

func IsValidSemanticVersion(version string) bool {
	_, err := semver.NewVersion(version)
	return err == nil
}

func IsValidK8sVersion(version string) bool {
	r := regexp.MustCompile(`^(v)([1-9]\d*)+((alpha|beta)([1-9]\d*)+){0,2}$`)
	return r.MatchString(version)
}
