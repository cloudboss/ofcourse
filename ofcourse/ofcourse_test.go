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
package ofcourse

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type resource struct{}

func (r *resource) Check(source Source, version Version, env Environment,
	logger *Logger) ([]Version, error) {
	versions := []Version{
		{
			"c": "d",
		},
	}
	return versions, nil
}

func (r *resource) In(outDir string, source Source, params Params,
	version Version, env Environment, logger *Logger) (Version, Metadata, error) {
	newVersion := Version{
		"c": "d",
	}
	metadata := Metadata{
		{
			Name:  "e",
			Value: "f",
		},
	}
	return newVersion, metadata, nil
}

func (r *resource) Out(inDir string, source Source, params Params,
	env Environment, logger *Logger) (Version, Metadata, error) {
	version := Version{
		"c": "d",
	}
	metadata := Metadata{
		{
			Name:  "e",
			Value: "f",
		},
	}
	return version, metadata, nil
}

type emptyResource struct{}

func (r *emptyResource) Check(source Source, version Version, env Environment,
	logger *Logger) ([]Version, error) {
	versions := []Version{}
	return versions, nil
}

func (r *emptyResource) In(outDir string, source Source, params Params,
	version Version, env Environment, logger *Logger) (Version, Metadata, error) {
	newVersion := Version{}
	metadata := Metadata{}
	return newVersion, metadata, nil
}

func (r *emptyResource) Out(inDir string, source Source, params Params,
	env Environment, logger *Logger) (Version, Metadata, error) {
	version := Version{}
	metadata := Metadata{}
	return version, metadata, nil
}

func Test_check(t *testing.T) {
	resource := &resource{}
	eResource := &emptyResource{}

	var tests = []struct {
		input  []byte
		output []byte
	}{
		{
			[]byte(`{"source":{},"version":null}`),
			[]byte(`[{"c":"d"}]`),
		},
		{
			[]byte(`{"source":{},"version":{"a":"b"}}`),
			[]byte(`[{"c":"d"}]`),
		},
		{
			[]byte(`{"source":{"a":"b"},"version":{"a":"b"}}`),
			[]byte(`[{"c":"d"}]`),
		},
		{
			[]byte(`{"source":{"a":"b"},"version":null}`),
			[]byte(`[{"c":"d"}]`),
		},
	}
	for _, test := range tests {
		output, _ := check(resource, test.input)
		assert.Equal(t, output, test.output)
	}

	tests = []struct {
		input  []byte
		output []byte
	}{
		{
			[]byte(`{"source":{},"version":null}`),
			[]byte(`[]`),
		},
	}
	for _, test := range tests {
		output, _ := check(eResource, test.input)
		assert.Equal(t, output, test.output)
	}
}

func Test_in(t *testing.T) {
	resource := &resource{}
	eResource := &emptyResource{}

	var tests = []struct {
		input  []byte
		output []byte
	}{
		{
			[]byte(`{"source":{},"params":{},"version":null}`),
			[]byte(`{"version":{"c":"d"},"metadata":[{"name":"e","value":"f"}]}`),
		},
		{
			[]byte(`{"source":{},"params":{},"version":{"a":"b"}}`),
			[]byte(`{"version":{"c":"d"},"metadata":[{"name":"e","value":"f"}]}`),
		},
		{
			[]byte(`{"source":{"a":"b"},"params":{"x":"y"},"version":{"a":"b"}}`),
			[]byte(`{"version":{"c":"d"},"metadata":[{"name":"e","value":"f"}]}`),
		},
		{
			[]byte(`{"source":{"a":"b"},"params":{"x":"y"},"version":null}`),
			[]byte(`{"version":{"c":"d"},"metadata":[{"name":"e","value":"f"}]}`),
		},
	}
	for _, test := range tests {
		output, _ := in(resource, "foo", test.input)
		assert.Equal(t, output, test.output)
	}

	tests = []struct {
		input  []byte
		output []byte
	}{
		{
			[]byte(`{"source":{},"version":null}`),
			[]byte(`{"version":{},"metadata":[]}`),
		},
	}
	for _, test := range tests {
		output, _ := in(eResource, "foo", test.input)
		assert.Equal(t, output, test.output)
	}
}

func Test_out(t *testing.T) {
	resource := &resource{}
	eResource := &emptyResource{}

	var tests = []struct {
		input  []byte
		output []byte
	}{
		{
			[]byte(`{"source":{},"params":{},"version":null}`),
			[]byte(`{"version":{"c":"d"},"metadata":[{"name":"e","value":"f"}]}`),
		},
		{
			[]byte(`{"source":{},"params":{},"version":{"a":"b"}}`),
			[]byte(`{"version":{"c":"d"},"metadata":[{"name":"e","value":"f"}]}`),
		},
		{
			[]byte(`{"source":{"a":"b"},"params":{"x":"y"},"version":{"a":"b"}}`),
			[]byte(`{"version":{"c":"d"},"metadata":[{"name":"e","value":"f"}]}`),
		},
		{
			[]byte(`{"source":{"a":"b"},"params":{"x":"y"},"version":null}`),
			[]byte(`{"version":{"c":"d"},"metadata":[{"name":"e","value":"f"}]}`),
		},
	}
	for _, test := range tests {
		output, _ := out(resource, "foo", test.input)
		assert.Equal(t, output, test.output)
	}

	tests = []struct {
		input  []byte
		output []byte
	}{
		{
			[]byte(`{"source":{},"version":null}`),
			[]byte(`{"version":{},"metadata":[]}`),
		},
	}
	for _, test := range tests {
		output, _ := out(eResource, "foo", test.input)
		assert.Equal(t, output, test.output)
	}
}

func Test_environmentGet(t *testing.T) {
	testsDefault := []struct {
		envVars      map[string]string
		variable     string
		defaultValue string
		value        string
	}{
		{
			envVars:      map[string]string{},
			variable:     "PATH",
			defaultValue: "/bin",
			value:        "/bin",
		},
		{
			envVars: map[string]string{
				"PATH": "/bin:/usr/bin",
			},
			variable:     "PATH",
			defaultValue: "",
			value:        "/bin:/usr/bin",
		},
		{
			envVars: map[string]string{
				"PATH": "/bin:/usr/bin",
			},
			variable:     "PATH",
			defaultValue: "/sbin:/usr/sbin",
			value:        "/bin:/usr/bin",
		},
	}
	testsNoDefault := []struct {
		envVars  map[string]string
		variable string
		value    string
	}{
		{
			envVars:  map[string]string{},
			variable: "PATH",
			value:    "",
		},
		{
			envVars: map[string]string{
				"PATH": "/bin:/usr/bin",
			},
			variable: "PATH",
			value:    "/bin:/usr/bin",
		},
	}
	for _, test := range testsDefault {
		env := NewEnvironment(test.envVars)
		value := env.Get(test.variable, test.defaultValue)
		assert.Equal(t, test.value, value)
	}
	for _, test := range testsNoDefault {
		env := NewEnvironment(test.envVars)
		value := env.Get(test.variable)
		assert.Equal(t, test.value, value)
	}
}

func Test_environmentGetAll(t *testing.T) {
	tests := []struct {
		envVars  map[string]string
		variable string
		ok       bool
	}{
		{
			envVars: map[string]string{
				"PATH": "/bin:/usr/bin",
			},
			variable: "HOME",
			ok:       false,
		},
		{
			envVars: map[string]string{
				"PATH": "/bin:/usr/bin",
			},
			variable: "PATH",
			ok:       true,
		},
	}
	for _, test := range tests {
		env := NewEnvironment(test.envVars)
		variables := env.GetAll()
		assert.Equal(t, test.envVars, variables)
		_, ok := variables[test.variable]
		assert.Equal(t, test.ok, ok)
	}
}
