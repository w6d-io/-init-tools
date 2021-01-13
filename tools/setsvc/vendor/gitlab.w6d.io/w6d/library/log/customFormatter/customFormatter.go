package runtime

import (
	logrus "gitlab.w6d.io/w6d/library/log/logrus"
)

// FunctionKey holds the function field
const FunctionKey = "function"

// FileKey holds the file field
const FileKey = "file"

const (
	//StackJump ...
	StackJump = 7
	//FieldLessStackJump ...
	FieldLessStackJump = 9
)

// Formatter decorates log entries with Function name (optional) Line number (optional) GoFuncDepth (optional)
type Formatter struct {
	CustomFormatter logrus.Formatter
	// When true, function name will be tagged to fields as well
	Function bool
	// When true, file name will be tagged to fields as well
	File bool
	// When true, only base name of the file will be tagged to fields
	BaseNameOnly bool
	// print path depth func ()
	PathFile bool
	// print path depth of go func ()
	FuncDepth int
}

// Format the current log entry by adding the function name and line number of the caller.
func (f *Formatter) Format(entry *logrus.Entry) ([]byte, error) {
	pathFunc, file, _ := f.GetCurrentPosition(entry)

	pathFunc, file = f.GetDepthFunction(pathFunc, file)

	data := logrus.Fields{}

	if f.Function {
		data[FunctionKey] = pathFunc
	}

	if f.File {
		if f.BaseNameOnly {
			data[FileKey] = file
		} else {
			data[FileKey] = file
		}
	}
	for k, v := range entry.Data {
		data[k] = v
	}
	entry.Data = data

	return f.CustomFormatter.Format(entry)
}
