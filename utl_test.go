package utl

import (
	"testing"

	l "github.com/stevenb256/log"
)

// TestBuild - tests building go code via code
func TestBuild(t *testing.T) {

	// start log
	err := l.StartLog("", "1.0", true)
	if nil != err {
		panic(err.Error())
	}
	defer l.CloseLog()

	// build windows version of utility
	err = Build("/Users/stevenbailey/go/src/github.com/stevenb256/nuclide", PlatformLinux)
	if l.Check(err) {
		t.Fail()
		return
	}
}
