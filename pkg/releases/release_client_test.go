package releases

import (
	"github.com/hashicorp/go-version"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSelectReleaseType(t *testing.T) {
	assert.Equal(t, IncludeStable, SelectReleaseType(false))
	assert.Equal(t, IncludeAll, SelectReleaseType(true))
}

func TestListAll(t *testing.T) {
	stableReleases, err := ListAll(IncludeStable)
	assert.Nil(t, err)
	assert.NotEmpty(t, stableReleases)

	allReleases, err := ListAll(IncludeAll)
	assert.Nil(t, err)
	assert.NotEmpty(t, allReleases)

	assert.Greater(t, len(allReleases), len(stableReleases))
}

func TestGetLatest(t *testing.T) {
	latestStable, err := GetLatest(IncludeStable)
	assert.Nil(t, err)
	assert.NotNil(t, latestStable)

	latestAll, err := GetLatest(IncludeAll)
	assert.Nil(t, err)
	assert.NotNil(t, latestAll)
}

func TestGetForVersion(t *testing.T) {
	release, exists, err := GetForVersion(IncludeAll, version.Must(version.NewVersion("1.12.16")))
	assert.Nil(t, err)
	assert.True(t, exists)
	assert.NotNil(t, release)
	assert.Equal(t, "go1.12.16", release.Version)

	release, exists, err = GetForVersion(IncludeStable, version.Must(version.NewVersion("1.12.16")))
	assert.Nil(t, err)
	assert.False(t, exists)
	assert.Equal(t, emptyRelease, release)
}
