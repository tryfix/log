package log

import (
	"github.com/logrusorgru/aurora/v4"
)

type Level string
type Output string

const (
	FATAL Level = `FATAL`
	ERROR Level = `ERROR`
	WARN  Level = `WARN`
	INFO  Level = `INFO`
	DEBUG Level = `DEBUG`
	TRACE Level = `TRACE`
)

const (
	OutText Output = `text`
	OutJson Output = `json`
)

var logColors = map[Level]string{
	FATAL: aurora.Red(`FATAL`).String(),
	ERROR: aurora.Red(`ERROR`).String(),
	WARN:  aurora.Yellow(`WARN `).String(),
	INFO:  aurora.Blue(`INFO `).String(),
	DEBUG: aurora.Cyan(`DEBUG`).String(),
	TRACE: aurora.Magenta(`TRACE`).String(),
}

var logTypes = map[Level]int{
	FATAL: 0,
	ERROR: 1,
	WARN:  2,
	INFO:  3,
	DEBUG: 4,
	TRACE: 5,
}
