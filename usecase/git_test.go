package usecase

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTagVersionUp(t *testing.T) {
	var u gitUsecase
	var tagVersionUp func(tag string, target string) (string, error) = (u).tagVersionUp

	var tag string
	var target string
	var version string
	var err error

	tag = "1.1.1"
	target = "major"
	version, err = tagVersionUp(tag, target)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "2.0.0", version)

	tag = "1.1.1"
	target = "minor"
	version, err = tagVersionUp(tag, target)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "1.2.0", version)

	tag = "1.1.1"
	target = "patch"
	version, err = tagVersionUp(tag, target)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "1.1.2", version)

	tag = "v1.1.1"
	target = "patch"
	version, err = tagVersionUp(tag, target)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "v1.1.2", version)

	tag = "v1.1.1aaa"
	target = "patch"
	_, err = tagVersionUp(tag, target)
	if err == nil {
		t.Fatal("error")
	}
}

func TestCheckVersionPrefix(t *testing.T) {
	var u gitUsecase
	var checkVersionPrefix func(tag string) bool = (u).checkVersionPrefix

	version := "v1.1.1"
	check := checkVersionPrefix(version)
	assert.Equal(t, true, check)

	version = "1.1.1"
	check = checkVersionPrefix(version)
	assert.Equal(t, false, check)
}

func TestIgnoreVersionPrefix(t *testing.T) {
	var u gitUsecase

	version := "v1.1.1"
	var getIgnoreVersionPrefix func(str string) string = (u).getIgnoreVersionPrefix
	version = getIgnoreVersionPrefix(version)
	assert.Equal(t, "1.1.1", version)
}

func TestIncrementVersion(t *testing.T) {
	var u gitUsecase
	var incrementVersion func(version string, target string) (string, error) = (u).incrementVersion

	var version string
	var target string
	var err error

	version = "1.1.1"
	target = "major"
	version, err = incrementVersion(version, target)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "2.0.0", version)

	version = "1.1.1"
	target = "minor"
	version, err = incrementVersion(version, target)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "1.2.0", version)

	version = "1.1.1"
	target = "patch"
	version, err = incrementVersion(version, target)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "1.1.2", version)
}
