package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_SemverUpdater_prepatch(t *testing.T) {
	updater, err := SemverUpdater("pre-patch")
	assert.NoError(t, err)
	v0 := "v0.0.0"

	// update with 'pre-patch' will increase the patch version if there's no
	// prerelease on the version and add a -beta.0
	v1, err := updater(v0)
	assert.NoError(t, err)
	assert.Equal(t, "v0.0.1-beta.0", v1)

	// update with 'pre-patch' will increase the version of the prerelease if there's
	// a prerelease present (-beta.1)
	v2, err := updater(v1)
	assert.NoError(t, err)
	assert.Equal(t, "v0.0.1-beta.1", v2)

	// update with 'pre-patch' will fail if the prerelease version is not formatter properly
	_, err = updater("v0.0.1-rc1")
	assert.Error(t, err)
}
