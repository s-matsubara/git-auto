package main

import (
	"os"
	"testing"
)

func TestMainRun(t *testing.T) {
	orig := os.Args
	os.Args = []string{"git-auto", "--help"}
	defer func() { os.Args = orig }()
	main()
	t.Log("executed")
}
