package log

import (
	"strings"

	"github.com/devspace-cloud/devspace/pkg/util/survey"
	"github.com/mgutz/ansi"
	"github.com/sirupsen/logrus"
)

var defaultLog Logger = &stdoutLogger{
	survey: survey.NewSurvey(),
	level:  logrus.DebugLevel,
}

// Discard is a logger implementation that just discards every log statement
var Discard = &DiscardLogger{}

// PrintLogo prints the devspace logo
func PrintLogo() {
	logo := `
     ____              ____                       
    |  _ \  _____   __/ ___| _ __   __ _  ___ ___ 
    | | | |/ _ \ \ / /\___ \| '_ \ / _` + "`" + ` |/ __/ _ \
    | |_| |  __/\ V /  ___) | |_) | (_| | (_|  __/
    |____/ \___| \_/  |____/| .__/ \__,_|\___\___|
                            |_|`

	stdout.Write([]byte(ansi.Color(logo+"\r\n\r\n", "cyan+b")))
}

// StartFileLogging logs the output of the global logger to the file default.log
func StartFileLogging() {
	defaultLogStdout, ok := defaultLog.(*stdoutLogger)
	if ok {
		defaultLogStdout.fileLogger = GetFileLogger("default")
	}

	OverrideRuntimeErrorHandler(false)
}

// GetInstance returns the Logger instance
func GetInstance() Logger {
	return defaultLog
}

// SetInstance sets the default logger instance
func SetInstance(logger Logger) {
	defaultLog = logger
}

// WriteColored writes a message in color
func writeColored(message string, color string) {
	defaultLog.Write([]byte(ansi.Color(message, color)))
}

//SetFakePrintTable is a testing tool that allows overwriting the function PrintTable
func SetFakePrintTable(fake func(s Logger, header []string, values [][]string)) {
	fakePrintTable = fake
}

var fakePrintTable func(s Logger, header []string, values [][]string)

// PrintTable prints a table with header columns and string values
func PrintTable(s Logger, header []string, values [][]string) {
	if fakePrintTable != nil {
		fakePrintTable(s, header, values)
		return
	}

	columnLengths := make([]int, len(header))

	for k, v := range header {
		columnLengths[k] = len(v)
	}

	// Get maximum column length
	for _, v := range values {
		for key, value := range v {
			if len(value) > columnLengths[key] {
				columnLengths[key] = len(value)
			}
		}
	}

	s.Write([]byte("\n"))

	// Print Header
	for key, value := range header {
		writeColored(" "+value+"  ", "green+b")

		padding := columnLengths[key] - len(value)

		if padding > 0 {
			s.Write([]byte(strings.Repeat(" ", padding)))
		}
	}

	s.Write([]byte("\n"))

	if len(values) == 0 {
		s.Write([]byte(" No entries found\n"))
	}

	// Print Values
	for _, v := range values {
		for key, value := range v {
			s.Write([]byte(" " + value + "  "))

			padding := columnLengths[key] - len(value)

			if padding > 0 {
				s.Write([]byte(strings.Repeat(" ", padding)))
			}
		}

		s.Write([]byte("\n"))
	}

	s.Write([]byte("\n"))
}