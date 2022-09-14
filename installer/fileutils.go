package installer

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

// makeDirectoryIfNotExists if not exist make dir
func makeDirectoryIfNotExists(path string, perm os.FileMode) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return mkdirAll(path, os.ModeDir|perm)
	}
	return nil
}

// mkdirAll all dirs
func mkdirAll(path string, perm os.FileMode) error {
	err := os.MkdirAll(path, perm)
	if err != nil {
		return errors.New(fmt.Sprintf("path %s, err: %s", path, err.Error()))
	}
	return nil
}

func emptyPath(path string) error {
	if path == "" {
		return errors.New("path can not be empty")
	}
	return nil
}

func checkDir(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return err
	}
	return nil
}

func checkVersionBool(version string) bool {
	var hasV bool
	var correctLen bool
	if version[0:1] == "v" { // make sure have a v at the start v0.1.1
		hasV = true
	}
	p := strings.Split(version, ".")
	if len(p) == 3 {
		correctLen = true
	}
	if hasV && correctLen {
		return true
	}
	return false
}

func checkVersion(version string) error {
	if version[0:1] != "v" { // make sure have a v at the start v0.1.1
		return errors.New(fmt.Sprintf("incorrect provided: %s version number try: v1.2.3", version))
	}
	p := strings.Split(version, ".")
	if len(p) >= 2 && len(p) < 4 {
	} else {
		return errors.New(fmt.Sprintf("incorrect lenght provided: %s version number try: v1.2.3", version))
	}
	return nil
}
