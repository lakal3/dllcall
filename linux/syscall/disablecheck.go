package syscall

import (
	"errors"
)

// DisableCgocheck disables go runtime feature to check embeded pointers to go memory. This function links private variable from go runtime
// to change variable value and may break in any future versions for Go.
// Use new flag -pin to generate Pinned memory guards available in go1.21
var DisableCgocheck = func() error {
	for _, dv := range dbgVars {
		if dv.name == "cgocheck" {
			*dv.value = 0
			return nil
		}
	}
	return errors.New("can't disable cgocheck. Set environment GODEBUG=cgocheck=0 or use -pin switch to generate call library")
}
