// Package resource is an implementation of a Concourse resource.
package resource

import (
	"encoding/json"
	"errors"
	"fmt"
	oc "github.com/cloudboss/ofcourse/ofcourse"
	"io/ioutil"
	"strconv"
)

// This resource is only a skeleton for getting started. What it does:
//
// For `Check`, it increments its version each time it is called, starting with `{"count": "1"}`. The
// next version would be `{"count": "2"}`, and so on.
//
// For `In`, it writes a file containing the latest version to its output directory.
//
// For `Out`, it looks in the directory created by `In` for the version file and reads it into a map,
// returning it back to Concourse.
//
// In case there is any confusion about which directions `In` and `Out` refer to, you should implement
// `In` to read remotely and write locally, while `Out` should read locally and write remotely. `In`
// receives an output directory and should place its result there after retrieving it from some (usually)
// remote source, configured through `source` in the pipeline's resource definition. `Out` receives an
// input directory to read from, and should place its result back to the remote source.
//
// This skeleton resource does not deal with a remote source; it simply reads and writes local files.
//
// Keep in mind, this is only an example! Beware of always returning a unique version on every check
// the way this example does. If you have many such resources that are checked frequently across many
// pipelines, it will put a lot of load on the database's CPU.

var (
	// ErrVersion means version map is malformed
	ErrVersion = errors.New(`key "count" not found in version map`)
	// ErrParam means parameters are malformed
	ErrParam = errors.New(`missing "version_path" parameter`)
)

// Resource implements the ofcourse.Resource interface.
type Resource struct{}

// Check implements the ofcourse.Resource Check method, corresponding to the /opt/resource/check command.
// This is called when Concourse does its resource checks, or when the `fly check-resource` command is run.
func (r *Resource) Check(source oc.Source, version oc.Version, env oc.Environment,
	logger *oc.Logger) ([]oc.Version, error) {
	// Returned `versions` should be all of the versions since the one given in the `version`
	// argument. If `version` is nil, then return the first available version. In many cases there
	// will be only one version to return, depending on the type of resource being implemented.
	// For example, a git resource would return a list of commits since the one given in the
	// `version` argument, whereas that would not make sense for resources which do not have any
	// kind of linear versioning.

	count := "1"
	if version != nil {
		oldCount, ok := version["count"]
		if !ok {
			return nil, ErrVersion
		}
		i, err := strconv.Atoi(oldCount)
		if err != nil {
			return nil, err
		}
		count = strconv.Itoa(i + 1)
	}

	// In Concourse, a version is an arbitrary set of string keys and string values.
	// This particular version consists of just one key and value.
	newVersion := oc.Version{"count": count}

	versions := []oc.Version{newVersion}

	return versions, nil
}

// In implements the ofcourse.Resource In method, corresponding to the /opt/resource/in command.
// This is called when a Concourse job does `get` on the resource.
func (r *Resource) In(outputDirectory string, source oc.Source, params oc.Params, version oc.Version,
	env oc.Environment, logger *oc.Logger) (oc.Version, oc.Metadata, error) {
	// Demo of logging. Resources should never use fmt.Printf or anything that writes
	// to standard output, as it will corrupt the JSON output expected by Concourse.
	logger.Errorf("This is an error")
	logger.Warnf("This is a warning")
	logger.Infof("This is an informational message")
	logger.Debugf("This is a debug message")

	// Write the `version` argument to a file in the output directory,
	// so the `Out` function can read it.
	outputPath := fmt.Sprintf("%s/version", outputDirectory)
	bytes, err := json.Marshal(version)
	if err != nil {
		return nil, nil, err
	}
	logger.Debugf("Version: %s", string(bytes))

	err = ioutil.WriteFile(outputPath, bytes, 0644)
	if err != nil {
		return nil, nil, err
	}

	// Metadata consists of arbitrary name/value pairs for display in the Concourse UI,
	// and may be returned empty if not needed.
	metadata := oc.Metadata{
		{
			Name:  "a",
			Value: "b",
		},
		{
			Name:  "c",
			Value: "d",
		},
	}

	// Here, `version` is passed through from the argument. In most cases, it makes sense
	// to retrieve the most recent version, i.e. the one in the `version` argument, and
	// then return it back unchanged. However, it is allowed to return some other version
	// or even an empty version, depending on the implementation.
	return version, metadata, nil
}

// Out implements the ofcourse.Resource Out method, corresponding to the /opt/resource/out command.
// This is called when a Concourse job does a `put` on the resource.
func (r *Resource) Out(inputDirectory string, source oc.Source, params oc.Params,
	env oc.Environment, logger *oc.Logger) (oc.Version, oc.Metadata, error) {
	// The `Out` function does not receive a `version` argument. Instead, we
	// will read the version from the file created by the `In` function, assuming
	// the pipeline does a `get` of this resource. The path to the version file
	// must be passed in the `put` parameters.
	versionPath, ok := params["version_path"]
	if !ok {
		return nil, nil, ErrParam
	}

	// The `inputDirectory` argument is a directory containing subdirectories for
	// all resources retrieved with `get` in a job, as well as all of the job's
	// task outputs.
	path := fmt.Sprintf("%s/%s", inputDirectory, versionPath)
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, nil, err
	}

	var version oc.Version
	err = json.Unmarshal(bytes, &version)
	if err != nil {
		return nil, nil, err
	}

	// Both `version` and `metadata` may be empty. In this case, we are returning
	// `version` retrieved from the file created by `In`, while `metadata` is empty.
	metadata := oc.Metadata{}
	return version, metadata, nil
}
