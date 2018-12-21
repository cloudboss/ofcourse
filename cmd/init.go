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

package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/cloudboss/ofcourse/ofcourse"
	"github.com/spf13/cobra"
)

var (
	resource       string
	dockerRegistry string
	importPath     string
	initCmd        = &cobra.Command{
		Use:   "init",
		Short: "Initialize Concourse resource project",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(resource, dockerRegistry, importPath)
		},
	}
)

func init() {
	RootCmd.AddCommand(initCmd)
	initCmd.Flags().StringVarP(&resource, "resource", "r",
		"", "Name of Concourse resource")
	initCmd.Flags().StringVarP(&dockerRegistry, "docker-registry", "R",
		"", "Registry where resource docker image will be pushed")
	initCmd.Flags().StringVarP(&importPath, "importPath", "i",
		"", "Go import path where code for resource will be located")
}

func run(resource, dockerRegistry, importPath string) error {
	if err := os.Mkdir(resource, 0777); err != nil {
		return err
	}
	data := map[string]string{
		"Resource":       resource,
		"DockerRegistry": dockerRegistry,
		"ImportPath":     importPath,
	}
	for _, asset := range ofcourse.AssetNames() {
		bytes, err := ofcourse.Asset(asset)
		if err != nil {
			return err
		}
		dir := fmt.Sprintf("%s/%s", resource, filepath.Dir(asset))
		if err := os.MkdirAll(dir, 0777); err != nil {
			return err
		}
		t := template.New(resource)
		t, err = t.Parse(string(bytes))
		if err != nil {
			return err
		}
		fd, err := os.Create(fmt.Sprintf("%s/%s", resource, asset))
		if err != nil {
			return err
		}
		defer fd.Close()
		err = t.Execute(fd, data)
		if err != nil {
			return err
		}
	}
	return nil
}
