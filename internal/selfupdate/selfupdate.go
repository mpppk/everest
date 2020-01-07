// Package selfupdate provides function to update binary
package selfupdate

import (
	"github.com/blang/semver"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
	"golang.org/x/xerrors"
)

const Version = "0.1.0"
const slug = "mpppk/everest"

// Do execute updating binary
func Do() (bool, error) {
	v := semver.MustParse(Version)
	latest, err := selfupdate.UpdateSelf(v, slug)
	if err != nil {
		return false, xerrors.Errorf("Binary update failed: %w", err)
	}
	return !latest.Version.Equals(v), nil
}
