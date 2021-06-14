package heavenlogger

import (
	"fmt"

	"gitlab.geax.io/demeter/gologger/constants"
	"gitlab.geax.io/demeter/gologger/gracefulshutdown"
	"gitlab.geax.io/demeter/gologger/slackAlert"
)

// SetFatalCallback config fatal callback
func (l *Logger) SetFatalCallback(fn func(msg string)) {
	l.fatalCallback = fn
}

// SetServiceCode config code
func (l *Logger) SetServiceCode(code constants.ServiceCode) {
	l.serviceCode = code
}

// OpenFile output to file
func (l *Logger) OpenFile(fileName string) {
	// develop mode skip
}

// FatalOnError ...
func (l *Logger) FatalOnError(err error, msg string, args ...interface{}) {
	if l.level > fatalLevel {
		return
	}
	if err != nil {
		m := fmt.Sprintf("%v %v %v", err, msg, fieldsToString(args))
		e := l.formatMsg(fatalLevel, m, true)
		fmt.Print(e)

		if slackAlert.IsSlackEnable() {
			slackAlert.SendError(e)
		}
		if l.fatalCallback != nil {
			l.fatalCallback(e)
		}

		gracefulshutdown.Shutdown()
	}
}

// FatalOnFail ...
func (l *Logger) FatalOnFail(success bool, msg string, args ...interface{}) {
	if l.level > fatalLevel {
		return
	}
	if !success {
		m := fmt.Sprintf("%v %v", msg, fieldsToString(args))
		e := l.formatMsg(fatalLevel, m, true)
		fmt.Print(e)

		if slackAlert.IsSlackEnable() {
			slackAlert.SendError(e)
		}

		if l.fatalCallback != nil {
			l.fatalCallback(e)
		}

		gracefulshutdown.Shutdown()
	}
}

// Error ...
func (l *Logger) Error(args ...interface{}) {
	if l.level > errorLevel {
		return
	}
	fields := fieldsToString(args)
	e := l.formatMsg(errorLevel, fields, true)
	fmt.Print(e)
	if slackAlert.IsSlackEnable() {
		slackAlert.SendWarns(e)
	}
}

// WarnCallStack ...
func (l *Logger) WarnCallStack(args ...interface{}) {
	if l.level > warnLevel {
		return
	}
	fields := fieldsToString(args)
	e := l.formatMsg(warnLevel, fields, true)
	fmt.Print(e)
	if slackAlert.IsSlackEnable() {
		slackAlert.SendWarns(e)
	}
}

// Warn ...
func (l *Logger) Warn(args ...interface{}) {
	if l.level > warnLevel {
		return
	}
	fields := fieldsToString(args)
	e := l.formatMsg(warnLevel, fields, false)
	fmt.Print(e)
	if slackAlert.IsSlackEnable() {
		slackAlert.SendWarns(e)
	}
}

// InfoCallStack ...
func (l *Logger) InfoCallStack(args ...interface{}) {
	if l.level > infoLevel {
		return
	}
	fields := fieldsToString(args)
	fmt.Print(l.formatMsg(infoLevel, fields, true))
}

// Info ...
func (l *Logger) Info(args ...interface{}) {
	l.log(infoLevel, args...)
}

// Debug ...
func (l *Logger) Debug(args ...interface{}) {
	l.log(debugLevel, args...)
}

// DebugCallStack console with stack trace
func (l *Logger) DebugCallStack(args ...interface{}) {
	if l.level > debugLevel {
		return
	}
	fields := fieldsToString(args)
	fmt.Print(l.formatMsg(debugLevel, fields, true))
}

// Fatalf equal to fmt.Printf but auto given log prefix
func (l *Logger) Fatalf(format string, args ...interface{}) {
	if l.level > fatalLevel {
		return
	}
	msg := fmt.Sprintf(format, args...)
	e := l.formatMsg(fatalLevel, msg, true)
	fmt.Print(e)
	if slackAlert.IsSlackEnable() {
		slackAlert.SendError(e)
	}

	if l.fatalCallback != nil {
		l.fatalCallback(e)
	}

	gracefulshutdown.Shutdown()
}

// Errorf equal to fmt.Printf but auto given log prefix
func (l *Logger) Errorf(format string, args ...interface{}) {
	if l.level > errorLevel {
		return
	}
	formatArguments(args)
	msg := fmt.Sprintf(format, args...)
	e := l.formatMsg(errorLevel, msg, true)
	fmt.Print(e)
	if slackAlert.IsSlackEnable() {
		slackAlert.SendWarns(e)
	}
}

// Warnf equal to fmt.Printf but auto given log prefix
func (l *Logger) Warnf(format string, args ...interface{}) {
	if l.level > warnLevel {
		return
	}
	formatArguments(args)
	msg := fmt.Sprintf(format, args...)
	e := l.formatMsg(warnLevel, msg, false)
	fmt.Print(e)
	if slackAlert.IsSlackEnable() {
		slackAlert.SendWarns(e)
	}
}

// Infof equal to fmt.Printf but auto given log prefix
func (l *Logger) Infof(format string, args ...interface{}) {
	l.logf(infoLevel, format, args...)
}

// Debugf equal to fmt.Printf but auto given log prefix
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.logf(debugLevel, format, args...)
}

// Debug ...
func (l *Logger) Debug1(args ...interface{}) {
	l.log(debug1Level, args...)
}

// Debugf equal to fmt.Printf but auto given log prefix
func (l *Logger) Debug1f(format string, args ...interface{}) {
	l.logf(debug1Level, format, args...)
}

// Debug ...
func (l *Logger) Debug2(args ...interface{}) {
	l.log(debug2Level, args...)
}

// Debugf equal to fmt.Printf but auto given log prefix
func (l *Logger) Debug2f(format string, args ...interface{}) {
	l.logf(debug2Level, format, args...)
}

// Debug ...
func (l *Logger) Debug3(args ...interface{}) {
	l.log(debug3Level, args...)
}

// Debugf equal to fmt.Printf but auto given log prefix
func (l *Logger) Debug3f(format string, args ...interface{}) {
	l.logf(debug3Level, format, args...)
}

// Debug ...
func (l *Logger) Debug4(args ...interface{}) {
	l.log(debug4Level, args...)
}

// Debugf equal to fmt.Printf but auto given log prefix
func (l *Logger) Debug4f(format string, args ...interface{}) {
	l.logf(debug4Level, format, args...)
}
