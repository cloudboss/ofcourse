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

type testResource struct{}

func (r *testResource) Check(source Source, version Version, logger *Logger) ([]Version, error) {
	versions := []Version{
		{
			"c": "d",
		},
	}
	return versions, nil
}

func (r *testResource) In(dir string, source Source, params Params,
	version Version, logger *Logger) (Version, Metadata, error) {
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

func (r *testResource) Out(dir string, source Source, params Params,
	logger *Logger) (Version, Metadata, error) {
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

func Test_check(t *testing.T) {
	resource := &testResource{}

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
}

func Test_in(t *testing.T) {
	resource := &testResource{}

	var tests = []struct {
		dir    string
		input  []byte
		output []byte
	}{
		{
			"foo",
			[]byte(`{"source":{},"params":{},"version":null}`),
			[]byte(`{"version":{"c":"d"},"metadata":[{"name":"e","value":"f"}]}`),
		},
		{
			"foo",
			[]byte(`{"source":{},"params":{},"version":{"a":"b"}}`),
			[]byte(`{"version":{"c":"d"},"metadata":[{"name":"e","value":"f"}]}`),
		},
		{
			"foo",
			[]byte(`{"source":{"a":"b"},"params":{"x":"y"},"version":{"a":"b"}}`),
			[]byte(`{"version":{"c":"d"},"metadata":[{"name":"e","value":"f"}]}`),
		},
		{
			"foo",
			[]byte(`{"source":{"a":"b"},"params":{"x":"y"},"version":null}`),
			[]byte(`{"version":{"c":"d"},"metadata":[{"name":"e","value":"f"}]}`),
		},
	}
	for _, test := range tests {
		output, _ := in(resource, test.dir, test.input)
		assert.Equal(t, output, test.output)
	}
}

func Test_out(t *testing.T) {
	resource := &testResource{}

	var tests = []struct {
		dir    string
		input  []byte
		output []byte
	}{
		{
			"foo",
			[]byte(`{"source":{},"params":{},"version":null}`),
			[]byte(`{"version":{"c":"d"},"metadata":[{"name":"e","value":"f"}]}`),
		},
		{
			"foo",
			[]byte(`{"source":{},"params":{},"version":{"a":"b"}}`),
			[]byte(`{"version":{"c":"d"},"metadata":[{"name":"e","value":"f"}]}`),
		},
		{
			"foo",
			[]byte(`{"source":{"a":"b"},"params":{"x":"y"},"version":{"a":"b"}}`),
			[]byte(`{"version":{"c":"d"},"metadata":[{"name":"e","value":"f"}]}`),
		},
		{
			"foo",
			[]byte(`{"source":{"a":"b"},"params":{"x":"y"},"version":null}`),
			[]byte(`{"version":{"c":"d"},"metadata":[{"name":"e","value":"f"}]}`),
		},
	}
	for _, test := range tests {
		output, _ := out(resource, test.dir, test.input)
		assert.Equal(t, output, test.output)
	}
}
