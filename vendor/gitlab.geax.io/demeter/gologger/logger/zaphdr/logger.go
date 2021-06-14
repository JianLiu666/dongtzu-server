package zaphdr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"

	"gitlab.geax.io/demeter/gologger/constants"
	"gitlab.geax.io/demeter/gologger/gracefulshutdown"
	"gitlab.geax.io/demeter/gologger/logger/zaphdr/stack"
	"gitlab.geax.io/demeter/gologger/slackAlert"
)

const (
	callerSkipOffset = 3 // zaphdr.fatal & zaphdr.FatalOnError & logger.FatalOnError
)

// key list
const (
	timeKey        = "time"
	stackstraceKey = "stacktrace"
	codekey        = "serviceCode"
	msgKey         = "msg"
	levelKey       = "level"
)

// level list
const (
	debugLevel = 0
	infoLevel  = 1
	warnLevel  = 2
	errorLevel = 3
	fatalLevel = 4
)

var levelNameMap = map[string]int{
	"debug": 0,
	"info":  1,
	"warn":  2,
	"error": 3,
	"fatal": 4,
}

var levelIDMap = map[int]string{
	0: "debug",
	1: "info",
	2: "warn",
	3: "error",
	4: "fatal",
}

var (
	workingDir string
)

// Logger of zap implement
type Logger struct {
	level         int
	serviceCode   constants.ServiceCode
	fatalCallback func(msg string)
	outputFile    *os.File
}

func init() {
	if dir, err := os.Getwd(); err == nil {
		workingDir = dir
	}
}

// New init logger
func New(level string) *Logger {
	return &Logger{
		level: levelNameMap[level],
	}
}

// Fatal console log err
func (l *Logger) fatal(fields ...interface{}) {
	output := l.jsonOutput(fatalLevel, true, fields)
	if l.fatalCallback != nil {
		l.fatalCallback(output)
	}
	if slackAlert.IsSlackEnable() {
		slackAlert.SendError(output)
	}

	gracefulshutdown.Shutdown()
}

// Error console log err
func (l *Logger) error(fields ...interface{}) string {
	e := l.jsonOutput(errorLevel, true, fields)
	if slackAlert.IsSlackEnable() {
		slackAlert.SendWarns(e)
	}
	return e
}

// Warn console log warn
func (l *Logger) warn(addStack bool, fields ...interface{}) {
	e := l.jsonOutput(warnLevel, addStack, fields)
	if slackAlert.IsSlackEnable() {
		slackAlert.SendWarns(e)
	}
}

// Info console log info
func (l *Logger) info(addStack bool, fields ...interface{}) {
	l.jsonOutput(infoLevel, addStack, fields)
}

// Debug console log debug
func (l *Logger) debug(addStack bool, fields ...interface{}) {
	l.jsonOutput(debugLevel, addStack, fields)
}

func now() string {
	return time.Now().UTC().Format(time.RFC3339)
}

func fieldsToString(fields []interface{}) string {
	fs, split := "", " "
	for idx := range fields {
		switch val := fields[idx].(type) {
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

func (l *Logger) jsonOutput(level int, addStack bool, fields []interface{}) string {
	if l.level > level {
		return ""
	}

	msg := fieldsToString(fields)
	outmap := map[string]interface{}{
		levelKey: levelIDMap[level],
		timeKey:  now(),
		msgKey:   msg,
	}

	if addStack {
		outmap[stackstraceKey] = stack.TakeStacktrace(callerSkipOffset + 1)
	}
	if l.serviceCode > 0 {
		outmap[codekey] = l.serviceCode
	}

	js, _ := jsonMarshal(outmap)
	output := string(js)
	fmt.Println(output)
	if l.outputFile != nil {
		line := output + "\n"
		if _, err := l.outputFile.WriteString(line); err != nil {
			panic(err)
		}
	}
	return output
}

func jsonMarshal(obj interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(obj)
	return buffer.Bytes(), err
}
