//go:build !go1.21
// +build !go1.21

package syscall

import (
	_ "unsafe"
)

type dbgVar struct {
	name  string
	value *int32 // for variables that can only be set at startup
}

//go:linkname dbgVars runtime.dbgvars
var dbgVars []dbgVar
