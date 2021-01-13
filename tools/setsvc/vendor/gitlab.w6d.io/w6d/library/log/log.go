package log

import (
	"io/ioutil"
	"os"

	logrus "gitlab.w6d.io/w6d/library/log/logrus"

	runtime "gitlab.w6d.io/w6d/library/log/customFormatter"
	"gitlab.w6d.io/w6d/library/log/logrus/hooks/writer"
)

const (
	Reset     = 0
	Black     = 30
	Red       = 31
	Green     = 32
	Yellow    = 33
	Blue      = 34
	Magenta   = 35
	Cyan      = 36
	White     = 37
	Blink     = 5
	BGblack   = 40
	BGred     = 41
	BGgreen   = 42
	BGyellow  = 43
	BGblue    = 44
	BGmagenta = 45
	BGcyan    = 46
	BGwhite   = 47
)

const (
	// UNSPECIFIED logs nothing
	UNSPECIFIED int = iota // 0 :
	// INFO logs Info, Warnings and Errors
	INFO // 1
	// DEBUG logs INFO and Debug
	DEBUG // 2
	// TRACE logs INFO and Debug and Trace
	TRACE // 3
	// ERROR just logs Errors
	ERROR // 4
)

var (
	log  = logrus.New()
	opts = Options{}
)

//Options ...
type Options struct {
	Color     bool
	Output    string
	FuncDepth int
	PathFile  bool
	Function  bool
	File      bool
	Level     int
}

func init() {
	log.SetOutput(ioutil.Discard) // Send all logs to nowhere by default

	log.AddHook(&writer.Hook{ // Send logs with level higher than warning to stderr
		Writer: os.Stderr,
		LogLevels: []logrus.Level{
			logrus.PanicLevel,
			logrus.FatalLevel,
			logrus.ErrorLevel,
			logrus.WarnLevel,
		},
	})
	log.AddHook(&writer.Hook{ // Send info and debug logs to stdout
		Writer: os.Stdout,
		LogLevels: []logrus.Level{
			logrus.InfoLevel,
			logrus.DebugLevel,
			logrus.TraceLevel,
		},
	})

	// default config
	opts = Options{
		FuncDepth: 0,
		PathFile:  false,
		Function:  true,
		File:      true,
		Color:     true,
		Level:     INFO,
	}
	logrus.ErrorColor = Red
	logrus.FatalColor = Red
	logrus.PanicColor = Red
	logrus.WarningColor = Yellow
	logrus.InfoColor = Cyan
	logrus.DebugColor = White
	logrus.TraceColor = White
	applyLogConfig()
}

func applyLogConfig() {
	if opts.Output == "json" {
		customFormatter := logrus.JSONFormatter{}
		runtimeFormatter := &runtime.Formatter{CustomFormatter: &customFormatter}
		runtimeFormatter.File = opts.File
		runtimeFormatter.Function = opts.Function
		runtimeFormatter.FuncDepth = opts.FuncDepth
		runtimeFormatter.PathFile = opts.PathFile
		customFormatter.TimestampFormat = "2006-01-02T15:04:05.000"
		customFormatter.FieldMap = logrus.FieldMap{
			logrus.FieldKeyTime: "timestamp",
		}
		log.Formatter = runtimeFormatter
	} else {
		customFormatter := logrus.TextFormatter{}
		runtimeFormatter := &runtime.Formatter{CustomFormatter: &customFormatter}
		runtimeFormatter.File = opts.File
		runtimeFormatter.Function = opts.Function
		runtimeFormatter.FuncDepth = opts.FuncDepth
		runtimeFormatter.PathFile = opts.PathFile
		customFormatter.TimestampFormat = "02-01-2006 15:04:05.0000"
		customFormatter.FullTimestamp = true
		customFormatter.ForceColors = true
		if !opts.Color {
			customFormatter.DisableColors = true
		}
		log.Formatter = runtimeFormatter
	}
	switch opts.Level {
	case INFO:
		log.Level = logrus.InfoLevel
	case DEBUG:
		log.Level = logrus.DebugLevel
	case ERROR:
		log.Level = logrus.ErrorLevel
	case TRACE:
		log.Level = logrus.TraceLevel
	default:
		log.Level = logrus.InfoLevel
	}
}

//SetLevel ...
func SetLevel(level int) {
	opts.Level = level
	applyLogConfig()
}

//SetColor ...
func SetColor(flag bool) {
	opts.Color = flag
	applyLogConfig()
}

//SetFunction ...
func SetFunction(flag bool) {
	opts.Function = flag
	applyLogConfig()
}

//SetFile ...
func SetFile(flag bool) {
	opts.File = flag
	applyLogConfig()
}

//SetFuncDepth ...
func SetFuncDepth(depth int) {
	opts.FuncDepth = depth
	applyLogConfig()
}

//SetPathFile ...
func SetPathFile(flag bool) {
	opts.PathFile = flag
	applyLogConfig()
}

//SetOutput ...
func SetOutput(output string) {
	opts.Output = output
	applyLogConfig()
}

//SetFatalColor ...
func SetFatalColor(color int) {
	logrus.FatalColor = color
}

//SetPanicColor ...
func SetPanicColor(color int) {
	logrus.PanicColor = color
}

//SetErrorColor ...
func SetErrorColor(color int) {
	logrus.ErrorColor = color
}

//SetWarningColor ...
func SetWarningColor(color int) {
	logrus.WarningColor = color
}

//SetInfoColor ...
func SetInfoColor(color int) {
	logrus.InfoColor = color
}

//SetDebugColor ...
func SetDebugColor(color int) {
	logrus.DebugColor = color
}

//SetTraceColor ...
func SetTraceColor(color int) {
	logrus.TraceColor = color
}

//Fields ...
type Fields = logrus.Fields

//Entry ...
type Entry = logrus.Entry

//WithFields ...
func WithFields(f Fields) *Entry {
	return logrus.NewEntry(log).WithFields(f)
}

// Fatal is equivalent to l.Critical(fmt.Sprint()) followed by a call to os.Exit(1).
func Fatal(args ...interface{}) {
	log.Fatal(args...)
}

// Fatalf is equivalent to l.Critical followed by a call to os.Exit(1).
func Fatalf(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}

// Panic is equivalent to l.Critical(fmt.Sprint()) followed by a call to panic().
func Panic(args ...interface{}) {
	log.Panic(args...)
}

// Panicf is equivalent to l.Critical followed by a call to panic().
func Panicf(format string, args ...interface{}) {
	log.Panicf(format, args...)
}

// Error logs a message using ERROR as log level.
func Error(args ...interface{}) {
	log.Error(args...)
}

// Errorf logs a message using ERROR as log level.
func Errorf(format string, args ...interface{}) {
	log.Errorf(format, args...)
}

// Warning logs a message using WARNING as log level.
func Warning(args ...interface{}) {
	log.Warning(args...)
}

// Warningf logs a message using WARNING as log level.
func Warningf(format string, args ...interface{}) {
	log.Warningf(format, args...)
}

// Info logs a message using INFO as log level.
func Info(args ...interface{}) {
	log.Info(args...)
}

// Infof logs a message using INFO as log level.
func Infof(format string, args ...interface{}) {
	log.Infof(format, args...)
}

// Debug logs a message using DEBUG as log level.
func Debug(args ...interface{}) {
	log.Debug(args...)
}

// Debugf logs a message using DEBUG as log level.
func Debugf(format string, args ...interface{}) {
	log.Debugf(format, args...)
}

// Trace logs a message using DEBUG as log level.
func Trace(args ...interface{}) {
	log.Trace(args...)
}

// Tracef logs a message using DEBUG as log level.
func Tracef(format string, args ...interface{}) {
	log.Tracef(format, args...)
}
