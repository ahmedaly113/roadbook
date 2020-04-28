// Package safe provides functions that enforce call patterns for potentially
// dangerous operations.
package safe

import (
	"fmt"
	"runtime/debug"

	"github.com/sirupsen/logrus"
)

// Do executes f with a protection against panic()s. Any goroutine that could
// potentially panic should be called via Do.
func Do(scope string, f func()) {
	defer RecoverAndLog(scope)
	f()
}

// DoOrClose executes f with a protection against panic()s. In case of a panic,
// close the provided channel -- typically a "quit channel". Can be used to signal
// that an unexpected condition has occurred and the application should shut down
// gracefully.
func DoOrClose(scope string, ch chan struct{}, f func()) {
	defer RecoverAndLog(scope, func(_ string) {
		close(ch)
	})
	f()
}

// RecoverAndLog provides consistent panic recovery and logging. For cases
// involving a goroutine, prefer safe.Do.
func RecoverAndLog(scope string, acts ...func(panicMsg string)) {
	if r := recover(); r != nil {
		logrus.WithField("panic", r).Errorf("panic in %s", scope)
		logrus.Errorf("Stack trace:\n%s", debug.Stack())
		for _, act := range acts {
			act(fmt.Sprint(r))
		}
	}
}
