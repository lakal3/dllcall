// +build !windows

package main

import (
	"fmt"
	"runtime"
)

func sys_init() error {
	return fmt.Errorf("System %s is not supported", runtime.GOOS)
}
