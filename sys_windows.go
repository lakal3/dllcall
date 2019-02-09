package main

import (
	"fmt"
	"runtime"
)

func sys_init() error {
	if runtime.GOARCH != "amd64" {
		return fmt.Errorf("Unsupported architecture %s", runtime.GOARCH)
	}
	return nil
}
