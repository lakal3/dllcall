//go:build go1.21
// +build go1.21

package syscall

import (
	"sync/atomic"
	_ "unsafe"
)

type dbgVar struct {
	name   string
	value  *int32        // for variables that can only be set at startup
	atomic *atomic.Int32 // for variables that can be changed during execution
	def    int32         // default value (ideally zero)
}

// Debugvars definition changed in go1.21!

//go:linkname dbgVars runtime.dbgvars
var dbgVars []*dbgVar
