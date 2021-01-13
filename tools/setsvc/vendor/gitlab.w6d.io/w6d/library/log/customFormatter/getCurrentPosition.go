package runtime

import (
	"runtime"
	"strconv"
	"strings"

	logrus "gitlab.w6d.io/w6d/library/log/logrus"
)

// Find the correct function name using "runtime.Caller(skip)"
func (f *Formatter) GetCurrentPosition(entry *logrus.Entry) (string, string, string) {
	skip := StackJump
	if len(entry.Data) == 0 {
		skip = FieldLessStackJump
	}
	for {
		pc, file, line, ok := runtime.Caller(skip)
		if !ok {
			file = "<???>"
			line = 1
		}
		lineNumber := strconv.Itoa(line)
		i := strings.LastIndex(file, "/")
		j := strings.Index(file[i+1:], ".")
		if j < 1 {
			file = "<???>"
		}
		file = file[i+1:] + ":" + lineNumber
		function := runtime.FuncForPC(pc).Name()

		if strings.Contains(function, "gitlab.w6d.io/w6d/library/log") {
			skip++
			continue
		}

		return function, file, lineNumber
	}
}
