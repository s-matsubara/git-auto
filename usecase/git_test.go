package usecase

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	versionBase       = "1.1.1"
	versionBasePrefix = "v1.1.1"
)

func TestTagVersionUp(t *testing.T) {
	var u gitUsecase
	var tagVersionUp = (u).tagVersionUp

	var tag string
	var target string
	var version string
	var err error

	tag = versionBase
	target = VersionMajor
	version, err = tagVersionUp(tag, target)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "2.0.0", version)

	tag = versionBase
	target = VersionMinor
	version, err = tagVersionUp(tag, target)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "1.2.0", version)

	tag = versionBase
	target = VersionPatch
	version, err = tagVersionUp(tag, target)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "1.1.2", version)

	tag = versionBasePrefix
	target = VersionPatch
	version, err = tagVersionUp(tag, target)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "v1.1.2", version)

	tag = versionBasePrefix + "aaa"
	target = VersionPatch
	_, err = tagVersionUp(tag, target)
	if err == nil {
		t.Fatal("error")
	}
}

func TestCheckVersionPrefix(t *testing.T) {
	var u gitUsecase
	var checkVersionPrefix = (u).checkVersionPrefix

	version := versionBasePrefix
	check := checkVersionPrefix(version)
	assert.Equal(t, true, check)

	version = versionBase
	check = checkVersionPrefix(version)
	assert.Equal(t, false, check)
}

func TestIgnoreVersionPrefix(t *testing.T) {
	var u gitUsecase

	version := versionBasePrefix
	var getIgnoreVersionPrefix = (u).getIgnoreVersionPrefix
	version = getIgnoreVersionPrefix(version)
	assert.Equal(t, "1.1.1", version)
}

func TestIncrementVersion(t *testing.T) {
	var u gitUsecase
	var incrementVersion = (u).incrementVersion

	var version string
	var target string
	var err error

	version = versionBase
	target = VersionMajor
	version, err = incrementVersion(version, target)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "2.0.0", version)

	version = versionBase
	target = VersionMinor
	version, err = incrementVersion(version, target)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "1.2.0", version)

	version = versionBase
	target = VersionPatch
	version, err = incrementVersion(version, target)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "1.1.2", version)
}
