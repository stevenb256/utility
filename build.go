package utl

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"

	l "github.com/stevenb256/log"
)

// PlatformWindows - build for windows platform
var PlatformWindows = "windows"

// PlatformMacOS - build for macOS platform
var PlatformMacOS = "macos"

// PlatformLinux - build for linux platform
var PlatformLinux = "linux"

// WindowsEXE -

// Build - builds a go directory for a specific platform
func Build(path, platform string) error {

	// locals
	var err error
	var extension string

	// save current env
	oldGOOS := os.Getenv("GOOS")
	oldGOARCH := os.Getenv("GOARCH")
	oldGOBIN := os.Getenv("GOBIN")
	oldCC := os.Getenv("CC")
	oldCXX := os.Getenv("CXX")
	os.Unsetenv("GOBIN") // for some reason this has to be unset for cross compile

	// http://crossgcc.rts-software.org/doku.php?id=compiling_for_linux
	if PlatformWindows == platform {
		os.Setenv("GOOS", "windows")
		os.Setenv("GOARCH", "amd64")
		os.Setenv("CGO_ENABLED", "1")
		os.Setenv("CC", "/usr/local/bin/x86_64-w64-mingw32-gcc")
		os.Setenv("CXX", "/usr/local/bin/x86_64-w64-mingw32-g++")
		extension = ".exe"
	} else if PlatformMacOS == platform {
		os.Setenv("GOOS", "darwin")
		os.Setenv("GOARCH", "amd64")
		os.Setenv("CGO_ENABLED", "1")
		os.Setenv("CC", "clang")
		os.Setenv("CXX", "clang++")
		extension = ""
	} else if PlatformLinux == platform {
		os.Setenv("GOOS", "linux")
		os.Setenv("GOARCH", "amd64")
		os.Setenv("CGO_ENABLED", "1")
		os.Setenv("CC", "/usr/local/gcc-4.8.1-for-linux64/bin/x86_64-pc-linux-gcc")
		os.Setenv("CXX", "/usr/local/gcc-4.8.1-for-linux64/bin/x86_64-pc-linux-g++")
		extension = ".linux"
	} else {
		return errors.New("invalid platform")
	}

	// find go exe path
	goExe, err := exec.LookPath("go")
	if l.Check(err) {
		return err
	}

	// run go install
	err = Execute(true, path, goExe, "build", "-o", "gbexe", "-gcflags", "-trimpath="+path, "-asmflags", "-trimpath="+path)
	if l.Check(err) {
		return err
	}

	// copies the built file to the filename with right extension
	err = MoveFile(Join(path, "gbexe"), Join(path, filepath.Base(path)+extension))
	if l.Check(err) {
		return err
	}

	// restore environment
	os.Setenv("GOOS", oldGOOS)
	os.Setenv("GOARCH", oldGOARCH)
	os.Setenv("CC", oldCC)
	os.Setenv("CXX", oldCXX)
	os.Setenv("GOBIN", oldGOBIN)

	// done
	return nil
}
