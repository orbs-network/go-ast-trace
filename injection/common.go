package injection

import (
	"runtime"
	"fmt"
	"strings"
	"time"
)

func getTimestamp() (string, int64) {
	t := time.Now()
	return t.Format("15:04:05.000000"), t.UnixNano()
}

func getCallStack() (string, string) {
	firstFunction := ""
	stackParts := []string{}
	callers := make([]uintptr, 4)
	runtime.Callers(3, callers)
	frames := runtime.CallersFrames(callers)
	var more bool = true
	var frame runtime.Frame
	for more {
		frame, more = frames.Next()
		if firstFunction == "" {
			firstFunction = frame.Function
		}
		file := frame.File
		if file == "" {
			break
		}
		pathParts := strings.Split(file, "/")
		if len(pathParts) >= 2 {
			file = strings.Join(pathParts[len(pathParts)-2:], "/")
		}
		stackParts = append(stackParts, fmt.Sprintf("%s:%d", file, frame.Line))
	}
	return firstFunction, strings.Join(stackParts, " ")
}