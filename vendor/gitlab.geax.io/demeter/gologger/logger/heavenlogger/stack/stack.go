package stack

import (
	"runtime"
	"sync"

	"gitlab.geax.io/demeter/gologger/logger/zaphdr/stack/bufferpool"
)

var (
	_stacktracePool = sync.Pool{
		New: func() interface{} {
			return newProgramCounters(64)
		},
	}
)

func getCallerFrame(skip int) (frame runtime.Frame, ok bool) {
	const skipOffset = 2 // skip getCallerFrame and Callers

	pc := make([]uintptr, 1)
	numFrames := runtime.Callers(skip+skipOffset, pc[:])
	if numFrames < 1 {
		return
	}

	frame, _ = runtime.CallersFrames(pc).Next()
	return frame, frame.PC != 0
}

// TakeStacktrace ...
func TakeStacktrace(skip int) string {
	buffer := bufferpool.Get()
	defer buffer.Free()
	programCounters := _stacktracePool.Get().(*programCounters)
	defer _stacktracePool.Put(programCounters)

	var numFrames int
	for {
		// Skip the call to runtime.Callers and takeStacktrace so that the
		// program counters start at the caller of takeStacktrace.
		numFrames = runtime.Callers(skip+2, programCounters.pcs)
		if numFrames < len(programCounters.pcs) {
			break
		}
		// Don't put the too-short counter slice back into the pool; this lets
		// the pool adjust if we consistently take deep stacktraces.
		programCounters = newProgramCounters(len(programCounters.pcs) * 2)
	}

	i := 0
	frames := runtime.CallersFrames(programCounters.pcs[:numFrames])

	// Note: On the last iteration, frames.Next() returns false, with a valid
	// frame, but we ignore this frame. The last frame is a a runtime frame which
	// adds noise, since it's only either runtime.main or runtime.goexit.
	for frame, more := frames.Next(); more; frame, more = frames.Next() {
		if i != 0 {
			buffer.AppendByte('\n')
		}
		i++
		buffer.AppendString(frame.Function)
		buffer.AppendByte('\n')
		buffer.AppendByte('\t')
		buffer.AppendString(frame.File)
		buffer.AppendByte(':')
		buffer.AppendInt(int64(frame.Line))
	}

	return buffer.String()
}

type programCounters struct {
	pcs []uintptr
}

func newProgramCounters(size int) *programCounters {
	return &programCounters{make([]uintptr, size)}
}
