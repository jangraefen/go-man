package releases

import (
	"errors"
	"net/http"
	"testing"

	"github.com/hashicorp/go-version"
	"github.com/stretchr/testify/assert"

	"github.com/jangraefen/go-man/internal/httputil"
)

func TestSelectReleaseType(t *testing.T) {
	assert.Equal(t, IncludeStable, SelectReleaseType(false))
	assert.Equal(t, IncludeAll, SelectReleaseType(true))
}

func TestListAll(t *testing.T) {
	t.Cleanup(func() {
		httputil.Client = http.DefaultClient
	})

	stableReleases, err := ListAll(IncludeStable)
	assert.NoError(t, err)
	assert.NotEmpty(t, stableReleases)

	allReleases, err := ListAll(IncludeAll)
	assert.NoError(t, err)
	assert.NotEmpty(t, allReleases)

	assert.Greater(t, len(allReleases), len(stableReleases))

	httputil.Client = httputil.StaticResponseClient(500, nil, errors.New("failure"))
	delete(ReleaseListCache, IncludeStable)

	stableReleases, err = ListAll(IncludeStable)
	assert.Error(t, err)
	assert.Empty(t, stableReleases)

	httputil.Client = httputil.StaticResponseClient(404, []byte("not found"), nil)
	delete(ReleaseListCache, IncludeStable)

	stableReleases, err = ListAll(IncludeStable)
	assert.Error(t, err)
	assert.Empty(t, stableReleases)
}

func TestGetLatest(t *testing.T) {
	t.Cleanup(func() {
		httputil.Client = http.DefaultClient
	})

	latestStable, err := GetLatest(IncludeStable)
	assert.NoError(t, err)
	assert.NotNil(t, latestStable)

	latestAll, err := GetLatest(IncludeAll)
	assert.NoError(t, err)
	assert.NotNil(t, latestAll)

	httputil.Client = httputil.StaticResponseClient(500, nil, errors.New("failure"))
	delete(ReleaseListCache, IncludeStable)

	latestStable, err = GetLatest(IncludeStable)
	assert.Error(t, err)
	assert.Nil(t, latestStable)

	httputil.Client = httputil.StaticResponseClient(404, []byte("not found"), nil)
	delete(ReleaseListCache, IncludeStable)

	latestStable, err = GetLatest(IncludeStable)
	assert.Error(t, err)
	assert.Nil(t, latestStable)
}

func TestGetForVersion(t *testing.T) {
	t.Cleanup(func() {
		httputil.Client = http.DefaultClient
	})

	release, exists, err := GetForVersion(IncludeAll, version.Must(version.NewVersion("1.12.16")))
	assert.NoError(t, err)
	assert.True(t, exists)
	assert.NotNil(t, release)
	assert.Equal(t, "go1.12.16", release.Version)

	release, exists, err = GetForVersion(IncludeStable, version.Must(version.NewVersion("1.12.16")))
	assert.NoError(t, err)
	assert.False(t, exists)
	assert.Nil(t, release)

	httputil.Client = httputil.StaticResponseClient(500, nil, errors.New("failure"))
	delete(ReleaseListCache, IncludeAll)

	release, exists, err = GetForVersion(IncludeAll, version.Must(version.NewVersion("1.12.16")))
	assert.Error(t, err)
	assert.False(t, exists)
	assert.Nil(t, release)

	httputil.Client = httputil.StaticResponseClient(404, []byte("not found"), nil)
	delete(ReleaseListCache, IncludeAll)

	release, exists, err = GetForVersion(IncludeAll, version.Must(version.NewVersion("1.12.16")))
	assert.Error(t, err)
	assert.False(t, exists)
	assert.Nil(t, release)
}
