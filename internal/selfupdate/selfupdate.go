// Package selfupdate provides function to update binary
package selfupdate

import (
	semverv3 "github.com/blang/semver"
	"github.com/blang/semver/v4"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
	"golang.org/x/xerrors"
)

const Version = "0.1.0"
const slug = "mpppk/everest"

// Do execute updating binary
func Do() (bool, error) {
	v := toV3PRVersion(semver.MustParse(Version))
	latest, err := selfupdate.UpdateSelf(v, slug)
	if err != nil {
		return false, xerrors.Errorf("Binary update failed: %w", err)
	}
	return !latest.Version.Equals(v), nil
}

func toV3PRVersion(v semver.Version) semverv3.Version {
	var v3PRVersions []semverv3.PRVersion
	for _, version := range v.Pre {
		v3PRVersions = append(v3PRVersions, semverv3.PRVersion(version))
	}

	v3v := semverv3.Version{
		Major: v.Major,
		Minor: v.Minor,
		Patch: v.Patch,
		Pre:   v3PRVersions,
		Build: v.Build,
	}
	return v3v
}
