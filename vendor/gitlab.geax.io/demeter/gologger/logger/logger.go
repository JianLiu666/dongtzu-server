package logger

import (
	"gitlab.geax.io/demeter/gologger/constants"
	"gitlab.geax.io/demeter/gologger/env"

	"gitlab.geax.io/demeter/gologger/logger/heavenlogger"
	"gitlab.geax.io/demeter/gologger/logger/zaphdr"
)

// Prodution mode list
const (
	ErrorLevel     = "error"
	ProdutionLevel = "info"
)

// Development mode list
const (
	DevelopmentErrorLevel = "dev-error"
	DevelopmentWarnLevel  = "dev-warn"
	DevelopmentInfoLevel  = "dev-info"
	DevelopmentLevel      = "debug" // high -> low = debug -> debug4 -> debug3 -> debug2 -> debug1
	Debug4Level           = "debug4"
	Debug3Level           = "debug3"
	Debug2Level           = "debug2"
	Debug1Level           = "debug1"
)

// level for handlers
const (
	errorLevel = "error"
	warnLevel  = "warn"
	infoLevel  = "info"
	debugLevel = "debug"
)

// Logger interface
type Logger interface {
	FatalOnError(err error, msg string, args ...interface{})
	FatalOnFail(success bool, msg string, args ...interface{})
	Error(args ...interface{})
	WarnCallStack(args ...interface{})
	Warn(args ...interface{})
	InfoCallStack(msg ...interface{})
	Info(msg ...interface{})
	Debug(msg ...interface{})
	DebugCallStack(msg ...interface{})
	Fatalf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Debugf(format string, args ...interface{})

	Debug1(msg ...interface{})
	Debug2(msg ...interface{})
	Debug3(msg ...interface{})
	Debug4(msg ...interface{})
	Debug1f(format string, args ...interface{})
	Debug2f(format string, args ...interface{})
	Debug3f(format string, args ...interface{})
	Debug4f(format string, args ...interface{})

	OpenFile(fileName string)
	SetFatalCallback(fn func(msg string))
	SetServiceCode(code constants.ServiceCode)
}

var instance Logger

// Init logger by level
func Init(level string) {
	switch level {
	case DevelopmentLevel, Debug1Level, Debug2Level, Debug3Level, Debug4Level:
		instance = heavenlogger.New(level)

	case DevelopmentInfoLevel:
		instance = heavenlogger.New(infoLevel)

	case DevelopmentWarnLevel:
		instance = heavenlogger.New(warnLevel)

	case DevelopmentErrorLevel:
		instance = heavenlogger.New(errorLevel)

	case ProdutionLevel:
		instance = zaphdr.New(infoLevel)

	case ErrorLevel:
		instance = zaphdr.New(errorLevel)

	default:
		instance = zaphdr.New(infoLevel)
	}
	instance.SetServiceCode(env.Setting.Code)
}

// OpenFile output to file
func OpenFile(fileName string) {
	instance.OpenFile(fileName)
}

func init() {
	Init(env.Setting.Level) // default mode
	instance.SetServiceCode(env.Setting.Code)
}

// SetServiceCode config code
func SetServiceCode(code constants.ServiceCode) {
	instance.SetServiceCode(code)
}

// SetFatalCallback config callback task, which run before fatal.panic
func SetFatalCallback(fn func(msg string)) {
	instance.SetFatalCallback(fn)
}

// FatalOnError console & panic when err isn't null, auto filled service-code field when arg was constants.ServiceCode type
func FatalOnError(err error, msg string, args ...interface{}) {
	instance.FatalOnError(err, msg, args...)
}

// FatalOnFail console & panic when not success, auto filled service-code field when arg was constants.ServiceCode type
func FatalOnFail(success bool, msg string, args ...interface{}) {
	instance.FatalOnFail(success, msg, args...)
}

// Error console when err isn't null, auto filled service-code field when arg was constants.ServiceCode type
func Error(args ...interface{}) {
	instance.Error(args...)
}

// WarnCallStack console warning & stack-trace, auto filled service-code field when arg was constants.ServiceCode type
func WarnCallStack(args ...interface{}) {
	instance.WarnCallStack(args...)
}

// Warn console warning, auto filled service-code field when arg was constants.ServiceCode type
func Warn(args ...interface{}) {
	instance.Warn(args...)
}

// InfoCallStack console info & stack-trace, auto filled service-code field when arg was constants.ServiceCode type
func InfoCallStack(args ...interface{}) {
	instance.InfoCallStack(args...)
}

// Info console info, auto filled service-code field when arg was constants.ServiceCode type
func Info(args ...interface{}) {
	instance.Info(args...)
}

// Debug console debug log, auto filled service-code field when arg was constants.ServiceCode type
func Debug(args ...interface{}) {
	instance.Debug(args...)
}

// DebugCallStack console debug log & stack-trace, auto filled service-code field when arg was constants.ServiceCode type
func DebugCallStack(args ...interface{}) {
	instance.DebugCallStack(args...)
}

// Fatalf equal to fmt.Printf but auto given log struct
func Fatalf(format string, args ...interface{}) {
	instance.Fatalf(format, args...)
}

// Errorf equal to fmt.Printf but auto given log struct
func Errorf(format string, args ...interface{}) {
	instance.Errorf(format, args...)
}

// Warnf equal to fmt.Printf but auto given log struct
func Warnf(format string, args ...interface{}) {
	instance.Warnf(format, args...)
}

// Infof equal to fmt.Printf but auto given log struct
func Infof(format string, args ...interface{}) {
	instance.Infof(format, args...)
}

// Debugf equal to fmt.Printf but auto given log struct
func Debugf(format string, args ...interface{}) {
	instance.Debugf(format, args...)
}

// Debug console debug log, auto filled service-code field when arg was constants.ServiceCode type
func Debug1(args ...interface{}) {
	instance.Debug1(args...)
}

// Debug console debug log, auto filled service-code field when arg was constants.ServiceCode type
func Debug2(args ...interface{}) {
	instance.Debug2(args...)
}

// Debug console debug log, auto filled service-code field when arg was constants.ServiceCode type
func Debug3(args ...interface{}) {
	instance.Debug3(args...)
}

// Debug console debug log, auto filled service-code field when arg was constants.ServiceCode type
func Debug4(args ...interface{}) {
	instance.Debug4(args...)
}

// Debugf equal to fmt.Printf but auto given log struct
func Debug1f(format string, args ...interface{}) {
	instance.Debug1f(format, args...)
}

// Debugf equal to fmt.Printf but auto given log struct
func Debug2f(format string, args ...interface{}) {
	instance.Debug2f(format, args...)
}

// Debugf equal to fmt.Printf but auto given log struct
func Debug3f(format string, args ...interface{}) {
	instance.Debug3f(format, args...)
}

// Debugf equal to fmt.Printf but auto given log struct
func Debug4f(format string, args ...interface{}) {
	instance.Debug4f(format, args...)
}
