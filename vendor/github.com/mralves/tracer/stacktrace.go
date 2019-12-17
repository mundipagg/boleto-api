package tracer

import (
	"fmt"
	"path"
	"runtime"
	"strings"
)

const DefaultDepth = 32

type Caller struct {
	File     string
	Function string
	Line     int
}

func (c Caller) String() string {
	_, file := path.Split(c.File)
	return fmt.Sprintf("at %s(%s:%d)", c.Function, file, c.Line)
}

type StackTrace []Caller

func GetStackTrace(skip int, maxDepth ...int) StackTrace {
	maxDepth = append(maxDepth, DefaultDepth)
	depth := maxDepth[0] + skip
	var stack StackTrace
	for i := skip; i < depth; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		funcName := "unknown"
		f := runtime.FuncForPC(pc)
		if f != nil {
			funcName = f.Name()
		}
		stack = append(stack, Caller{
			File:     file,
			Function: funcName,
			Line:     line,
		})
		if funcName == "main.main" {
			break
		}
	}
	return stack
}

func (st StackTrace) String() string {
	var formatted strings.Builder
	for i := range st {
		caller := st[len(st)-1-i]
		prefix := strings.Repeat("  ", i)
		// fmt.Fprintf(&formatted, "%s%s at %s:%d\n", prefix, caller.Function, caller.File, caller.Line)
		fmt.Fprintf(&formatted, "%s%s\n", prefix, caller.String())
	}
	return formatted.String()
}
