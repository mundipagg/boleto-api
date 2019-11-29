package tracer

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"text/template"
	"time"
)

var pointerPattern = regexp.MustCompile(`<\*>\([^)]+\)`)

type Entry struct {
	Owner         string
	Level         uint8
	Message       string
	Args          []interface{}
	Time          time.Time
	TransactionId string
	StackTrace    StackTrace
}

type Writer interface {
	Write(entry Entry)
}

func DefaultStringer(value interface{}) string {
	if stringer, ok := value.(fmt.Stringer); ok {
		return stringer.String()
	} else if err, ok := value.(error); ok {
		return err.Error()
	}
	rep := fmt.Sprintf("%+v", value)
	return pointerPattern.ReplaceAllStringFunc(rep, func(str string) string {
		return ""
	})
}

type Formatter func(entry Entry) string

type FileWriter struct {
	formatter Formatter
	writer    io.Writer
}

func (fw *FileWriter) Write(entry Entry) {
	_, err := fw.writer.Write([]byte(fw.formatter(entry)))
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "FAIL ON LOGGING ENTRY BECAUSE: %v\n", err.Error())
	}
}

func NewFileWriter(file io.Writer, formatter Formatter) *FileWriter {
	return &FileWriter{
		writer:    file,
		formatter: formatter,
	}
}

var variablePattern = regexp.MustCompile(`@\w+`)

func SimpleFormatter(format string, stringer func(interface{}) string) Formatter {
	if stringer == nil {
		stringer = DefaultStringer
	}
	format = variablePattern.ReplaceAllStringFunc(format, func(match string) string {
		return fmt.Sprintf("{{.%s}}", strings.Title(match[1:]))
	})
	t, err := template.New("log").Parse(format)
	if err != nil {
		panic(err)
	}
	return func(entry Entry) string {
		var buf bytes.Buffer
		reps := make([]interface{}, len(entry.Args))
		for i := range entry.Args {
			reps[i] = stringer(entry.Args[i])
		}
		entry.Args = reps
		err := t.Execute(&buf, struct {
			Entry
			Time      string
			LevelName string
		}{
			Entry:     entry,
			Time:      entry.Time.Format(time.RFC3339),
			LevelName: LevelNames[entry.Level],
		})
		if err != nil {
			return fmt.Sprintf("FAILED TO FORMAT ENTRY BECAUSE: %v\n", err.Error())
		}
		return buf.String()
	}
}
