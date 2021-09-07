package cmdutil

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"k8s.io/klog/v2"
)

// GetSysSig register exit signals
func GetSysSig() <-chan os.Signal {
	// handle system signal
	ch := make(chan os.Signal, 1)
	signal.Notify(
		ch,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGTERM,
	)
	return ch
}

// ReallyCrash For testing, bypass HandleCrash.
var ReallyCrash bool

// HandleCrash simply catches a crash and logs an error. Meant to be called via defer.
func HandleCrash() {
	if ReallyCrash {
		return
	}

	r := recover()
	if r != nil {
		callers := ""
		for i := 1; true; i++ {
			_, file, line, ok := runtime.Caller(i)
			if !ok {
				break
			}
			callers = callers + fmt.Sprintf("%v:%v\n", file, line)
		}
		klog.Warningf("Recovered from panic: %#v (%v)\n%v", r, r, callers)
	}
}

// LogTraceback prints traceback to given logger
func LogTraceback(r interface{}, depth int, logger interface {
	Infof(fmt string, arg ...interface{})
}) {
	var callers []string
	for i := depth; true; i++ {
		_, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		callers = append(callers, fmt.Sprintf("%v:%v", file, line))
	}
	callers = callers[0 : len(callers)-1]

	logger.Infof("panic: %#v", r)
	for i := range callers {
		logger.Infof("tb| %s", callers[i])
	}
}
