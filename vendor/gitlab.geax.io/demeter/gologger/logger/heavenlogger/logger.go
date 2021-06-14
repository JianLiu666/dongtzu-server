package heavenlogger

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"

	"gitlab.geax.io/demeter/gologger/constants"
	"gitlab.geax.io/demeter/gologger/logger/heavenlogger/stack"
	"gitlab.geax.io/demeter/gologger/util/color"
)

// Logger implement of paradise
type Logger struct {
	fatalCallback func(msg string)
	level         int
	serviceCode   constants.ServiceCode
}

// New instance
func New(level string) *Logger {
	return &Logger{
		level: levelNameMap[level],
	}
}

// level list
const (
	debug1Level = -4
	debug2Level = -3
	debug3Level = -2
	debug4Level = -1
	debugLevel  = 0
	infoLevel   = 1
	warnLevel   = 2
	errorLevel  = 3
	fatalLevel  = 4
)

var levelNameMap = map[string]int{
	"debug1": -4,
	"debug2": -3,
	"debug3": -2,
	"debug4": -1,
	"debug":  0,
	"info":   1,
	"warn":   2,
	"error":  3,
	"fatal":  4,
}

var levelIDMap = map[int]string{
	-4: "debug1",
	-3: "debug2",
	-2: "debug3",
	-1: "debug4",
	0:  "debug",
	1:  "info",
	2:  "warn",
	3:  "error",
	4:  "fatal",
}

var colorMap = map[int]string{
	-4: color.Blue1.Add("DEBUG"),
	-3: color.Blue2.Add("DEBUG"),
	-2: color.Blue3.Add("DEBUG"),
	-1: color.Blue4.Add("DEBUG"),
	0:  color.Blue.Add("DEBUG"),
	1:  color.Cyan.Add("INFO"),
	2:  color.Yellow.Add("WARN"),
	3:  color.Red.Add("ERROR"),
	4:  color.Magenta.Add("FATAL"),
}

var (
	workingDir string
)

const (
	callerSkipOffset = 3 // logger.InfoCallStack & heavenlogger.InfoCallStack & heavenlogger.formatMsg
)

func init() {
	if dir, err := os.Getwd(); err == nil {
		workingDir = dir
	}
}

func now() string {
	return color.Green.Add(time.Now().UTC().Format(time.RFC3339))
}

func (l *Logger) getServiceCode() string {
	if l.serviceCode > 0 {
		return fmt.Sprintf(" service-code: %v", l.serviceCode)
	}

	return ""
}

func fieldsToString(fields []interface{}) string {
	fs, split := "", " "
	for idx := range fields {
		switch val := fields[idx].(type) {
		case constants.ServiceCode:
			continue
		case error:
			fs += fmt.Sprint(val.Error(), split)
		case string:
			fs += fmt.Sprint(val, split)
		default:
			js, _ := json.Marshal(val)
			fs += fmt.Sprint(string(js), split)
		}
	}

	return strings.Trim(fs, split)
}

func formatArguments(args []interface{}) {
	for idx := range args {
		val := reflect.ValueOf(args[idx])

		if val.Kind() == reflect.Ptr {
			val = reflect.Indirect(val)
		}

		if !val.IsValid() {
			continue
		}

		switch args[idx].(type) {
		case error:
			args[idx] = args[idx].(error).Error()
			continue
		}

		switch val.Interface().(type) {
		case error:
			args[idx] = val.Interface().(error).Error()
			continue
		}

		switch val.Kind() {
		case reflect.Struct, reflect.Array, reflect.Map, reflect.Slice:
			jstr, _ := json.Marshal(args[idx])
			args[idx] = string(jstr)
		}
	}
}

func (l *Logger) formatMsg(level int, msg string, isStack bool) string {
	code := l.getServiceCode()
	if !isStack {
		return fmt.Sprintf("%v [%v]%v %v\n", now(), colorMap[level], code, msg)
	}
	stack := stack.TakeStacktrace(callerSkipOffset)
	return fmt.Sprintf("%v [%v]%v %v\n%v\n", now(), colorMap[level], code, msg, stack)
}

func (l *Logger) log(level int, args ...interface{}) {
	if l.level > level {
		return
	}
	fields := fieldsToString(args)
	fmt.Print(l.formatMsg(level, fields, false))
}

func (l *Logger) logf(level int, format string, args ...interface{}) {
	if l.level > level {
		return
	}
	formatArguments(args)
	msg := fmt.Sprintf(format, args...)
	fmt.Print(l.formatMsg(level, msg, false))
}
