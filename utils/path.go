package utils

import (
	"fmt"
	"path/filepath"
	"runtime"
)

var (
	_, b, _, _ = runtime.Caller(0)
	basePath   = filepath.Dir(b)
)

// Base gives the current caller base path.
// It does not depend on where it is called from.
func Base() string {
	return fmt.Sprintf("%s/", filepath.Dir(basePath))
}

// FromBase concat the base path with the given path.
func FromBase(addPath string) string {
	return fmt.Sprintf("%s%s", Base(), addPath)
}
