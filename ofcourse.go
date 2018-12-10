// Copyright © 2018 Joseph Wright <joseph@cloudboss.co>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

// Package ofcourse reduces boilerplate for implementing Concourse resources.
package ofcourse

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/fatih/color"
)

const (
	// ErrorLevel for logging.
	ErrorLevel = "error"
	// WarnLevel for logging.
	WarnLevel = "warn"
	// InfoLevel for logging.
	InfoLevel = "info"
	// DebugLevel for logging.
	DebugLevel = "debug"
)

const (
	errorLevel = iota
	warnLevel
	infoLevel
	debugLevel
)

var (
	internalLogger = NewLogger(ErrorLevel)
)

// Logger is passed to resource functions so that they can log to the Concourse UI without
// printing to stdout, as doing so would corrupt the JSON output expected by Concourse.
// Resources using this library may set their log level in the source configuration of
// the pipeline using the parameter `log_level`, which may have a case insensitive value
// of "error", "warn", "info", or "debug". Respective to those log levels, the output
// colors are red, yellow, green, and blue.
type Logger struct {
	Level int
}

// NewLogger returns a logger instance with the given log level, defaulting to "info" if
// the given level is not recognized.
func NewLogger(level string) *Logger {
	var intLevel int
	switch strings.ToLower(level) {
	case ErrorLevel:
		intLevel = errorLevel
	case WarnLevel:
		intLevel = warnLevel
	case InfoLevel:
		intLevel = infoLevel
	case DebugLevel:
		intLevel = debugLevel
	default:
		intLevel = infoLevel
	}
	return &Logger{Level: intLevel}
}

// Errorf logs a red message to the Concourse UI.
func (l *Logger) Errorf(message string, args ...interface{}) {
	if l.Level >= errorLevel {
		color.New(color.FgHiRed).Fprintf(color.Error, message, args...)
	}
}

// Warnf logs a yellow message to the Concourse UI.
func (l *Logger) Warnf(message string, args ...interface{}) {
	if l.Level >= warnLevel {
		color.New(color.FgHiYellow).Fprintf(color.Error, message, args...)
	}
}

// Infof logs a green message to the Concourse UI.
func (l *Logger) Infof(message string, args ...interface{}) {
	if l.Level >= infoLevel {
		color.New(color.FgHiGreen).Fprintf(color.Error, message, args...)
	}
}

// Debugf logs a blue message to the Concourse UI.
func (l *Logger) Debugf(message string, args ...interface{}) {
	if l.Level >= debugLevel {
		color.New(color.FgHiBlue).Fprintf(color.Error, message, args...)
	}
}

// Version represents Concourse a version, which is a set of key/value pairs.
type Version map[string]string

// Metadata represents Concourse metadata, which is a set of key/value pairs that
// are used for printing extra information to the Concourse UI on `get` or `put`.
type Metadata []NameVal

// NameVal is one item of a Metadata array.
type NameVal struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// Source is the pipeline source configuration for the resource.
type Source map[string]interface{}

// Params is the pipeline `get` or `put` parameters for the resource.
type Params map[string]interface{}

// CheckInput represents the stdin to an /opt/resource/check command.
type CheckInput struct {
	Source  Source  `json:"source"`
	Version Version `json:"version"`
}

// InInput represents the stdin to an /opt/resource/in command.
type InInput struct {
	Source  Source  `json:"source"`
	Params  Params  `json:"params"`
	Version Version `json:"version"`
}

// OutInput represents the stdin to an /opt/resource/out command.
type OutInput struct {
	Source Source `json:"source"`
	Params Params `json:"params"`
}

type inOutOutput struct {
	Version  Version  `json:"version"`
	Metadata Metadata `json:"metadata"`
}

// Resource is a type that contains Check, In, and Out methods. The user of this
// library must implement this interface.
type Resource interface {
	Check(src Source, ver Version, log *Logger) ([]Version, error)
	In(dir string, src Source, par Params, ver Version, log *Logger) (Version, Metadata, error)
	Out(dir string, src Source, par Params, log *Logger) (Version, Metadata, error)
}

func check(resource Resource, input []byte) ([]byte, error) {
	var checkInput CheckInput
	err := json.Unmarshal(input, &checkInput)
	if err != nil {
		return nil, err
	}

	var logger *Logger
	if logLevel, ok := checkInput.Source["log_level"].(string); ok {
		logger = NewLogger(logLevel)
	} else {
		logger = NewLogger("info")
	}

	versions, err := resource.Check(checkInput.Source, checkInput.Version, logger)
	if err != nil {
		return nil, err
	}

	versionBytes, err := json.Marshal(versions)
	if err != nil {
		return nil, err
	}

	return versionBytes, nil
}

func in(resource Resource, dir string, input []byte) ([]byte, error) {
	var inInput InInput
	err := json.Unmarshal(input, &inInput)
	if err != nil {
		return nil, err
	}

	var logger *Logger
	if logLevel, ok := inInput.Source["log_level"].(string); ok {
		logger = NewLogger(logLevel)
	} else {
		logger = NewLogger("info")
	}

	version, metadata, err := resource.In(dir, inInput.Source, inInput.Params, inInput.Version, logger)
	if err != nil {
		return nil, err
	}

	output := inOutOutput{
		Version:  version,
		Metadata: metadata,
	}
	return json.Marshal(output)
}

func out(resource Resource, dir string, input []byte) ([]byte, error) {
	var outInput OutInput
	err := json.Unmarshal(input, &outInput)
	if err != nil {
		return nil, err
	}

	var logger *Logger
	if logLevel, ok := outInput.Source["log_level"].(string); ok {
		logger = NewLogger(logLevel)
	} else {
		logger = NewLogger("info")
	}

	version, metadata, err := resource.Out(dir, outInput.Source, outInput.Params, logger)
	if err != nil {
		return nil, err
	}

	output := inOutOutput{
		Version:  version,
		Metadata: metadata,
	}
	return json.Marshal(output)
}

// Check takes an implementation of Resource as its input. The Main function
// for the /opt/resource/check command that is run by Concourse should create
// an instance of the resource and pass it to this function.
func Check(resource Resource) {
	input, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		internalLogger.Errorf("%s\n", err)
		os.Exit(1)
	}

	output, err := check(resource, input)
	if err != nil {
		internalLogger.Errorf("%s\n", err)
		os.Exit(1)
	}

	fmt.Printf(string(output))
}

// In takes an implementation of Resource as its input. The Main function
// for the /opt/resource/in command that is run by Concourse should create
// an instance of the resource and pass it to this function.
func In(resource Resource) {
	if len(os.Args) < 2 {
		internalLogger.Errorf("missing directory argument\n")
		os.Exit(1)
	}
	dir := os.Args[1]

	input, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		internalLogger.Errorf("%s\n", err)
		os.Exit(1)
	}

	output, err := in(resource, dir, input)
	if err != nil {
		internalLogger.Errorf("%s\n", err)
		os.Exit(1)
	}

	fmt.Printf(string(output))
}

// Out takes an implementation of Resource as its input. The Main function
// for the /opt/resource/out command that is run by Concourse should create
// an instance of the resource and pass it to this function.
func Out(resource Resource) {
	if len(os.Args) < 2 {
		internalLogger.Errorf("missing directory argument\n")
		os.Exit(1)
	}
	dir := os.Args[1]

	input, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		internalLogger.Errorf("%s\n", err)
		os.Exit(1)
	}

	output, err := out(resource, dir, input)
	if err != nil {
		internalLogger.Errorf("%s\n", err)
		os.Exit(1)
	}

	fmt.Printf(string(output))
}
