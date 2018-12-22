package resource

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	oc "github.com/cloudboss/ofcourse/ofcourse"
	"github.com/stretchr/testify/assert"
)

var (
	testLogger = oc.NewLogger(oc.SilentLevel)
)

func TestCheck(t *testing.T) {
	env := oc.NewEnvironment()
	var testCases = []struct {
		sourceIn    oc.Source
		versionIn   oc.Version
		versionsOut []oc.Version
		err         error
	}{
		{
			oc.Source{},
			nil,
			[]oc.Version{oc.Version{"count": "1"}},
			nil,
		},
		{
			oc.Source{},
			oc.Version{"count": "1000"},
			[]oc.Version{oc.Version{"count": "1001"}},
			nil,
		},
		{
			oc.Source{},
			oc.Version{"numero": "1"},
			nil,
			ErrVersion,
		},
	}

	for _, tc := range testCases {
		r := Resource{}
		versions, err := r.Check(tc.sourceIn, tc.versionIn, env, testLogger)
		assert.Equal(t, tc.versionsOut, versions)
		assert.Equal(t, tc.err, err)
	}
}

func TestIn(t *testing.T) {
	env := oc.NewEnvironment()
	var testCases = []struct {
		sourceIn    oc.Source
		paramsIn    oc.Params
		versionIn   oc.Version
		metadataOut oc.Metadata
		err         error
	}{
		{
			oc.Source{},
			oc.Params{},
			oc.Version{"count": "1"},
			oc.Metadata{
				{Name: "a", Value: "b"}, {Name: "c", Value: "d"},
			},
			nil,
		},
		{
			oc.Source{},
			oc.Params{},
			oc.Version{"count": "1234"},
			oc.Metadata{
				{Name: "a", Value: "b"}, {Name: "c", Value: "d"},
			},
			nil,
		},
	}

	for _, tc := range testCases {
		td, err := ioutil.TempDir("", "resource-")
		assert.Nil(t, err)
		defer os.RemoveAll(td)

		r := Resource{}
		version, metadata, err := r.In(td, tc.sourceIn, tc.paramsIn, tc.versionIn, env, testLogger)

		// First test the return values
		assert.Equal(t, tc.versionIn, version)
		assert.Equal(t, tc.metadataOut, metadata)
		assert.Equal(t, tc.err, err)

		// Ensure the output file has been created
		path := fmt.Sprintf("%s/version", td)
		_, fileErr := os.Stat(path)
		assert.False(t, os.IsNotExist(fileErr))

		// Read the output file
		bytes, err := ioutil.ReadFile(path)
		assert.Nil(t, err)

		// Ensure the contents are as expected
		var readVersion oc.Version
		err = json.Unmarshal(bytes, &readVersion)
		assert.Equal(t, tc.versionIn, readVersion)
	}
}

func TestOut(t *testing.T) {
	env := oc.NewEnvironment()
	var testCases = []struct {
		sourceIn    oc.Source
		paramsIn    oc.Params
		versionOut  oc.Version
		metadataOut oc.Metadata
		err         error
	}{
		{
			oc.Source{},
			oc.Params{},
			nil,
			nil,
			ErrParam,
		},
		{
			oc.Source{},
			oc.Params{"version_path": "hey/version"},
			oc.Version{"count": "4567"},
			oc.Metadata{},
			nil,
		},
	}

	for _, tc := range testCases {
		// Make a temp dir to be the input directory
		td, err := ioutil.TempDir("", "resource-")
		assert.Nil(t, err)
		defer os.RemoveAll(td)

		// If version_path passed in params, create the version file there
		versionPath, ok := tc.paramsIn["version_path"].(string)
		if ok {
			fullVersionPath := fmt.Sprintf("%s/%s", td, versionPath)

			// Create any leading subdirectories in version path
			err = os.MkdirAll(filepath.Dir(fullVersionPath), 0777)
			assert.Nil(t, err)

			// Marshal the version map to JSON
			contents, err := json.Marshal(tc.versionOut)
			assert.Nil(t, err)

			// Write JSON to file
			err = ioutil.WriteFile(fullVersionPath, []byte(contents), 0666)
			assert.Nil(t, err)
		}

		r := Resource{}
		version, metadata, err := r.Out(td, tc.sourceIn, tc.paramsIn, env, testLogger)
		assert.Equal(t, tc.versionOut, version)
		assert.Equal(t, tc.metadataOut, metadata)
		assert.Equal(t, tc.err, err)
	}
}
