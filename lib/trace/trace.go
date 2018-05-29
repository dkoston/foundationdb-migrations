package trace

import (
	"runtime"
	"fmt"
)

// Trace
// Returns function name, with one additional function scope above
func Trace() string {
	slice := make([]uintptr, 10)
	runtime.Callers(2, slice)
	function := runtime.FuncForPC(slice[0])
	return fmt.Sprintf("%s", function.Name())
}