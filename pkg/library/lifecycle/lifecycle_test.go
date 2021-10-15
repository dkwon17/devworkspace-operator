//
// Copyright (c) 2019-2021 Red Hat, Inc.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package lifecycle

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"

	dw "github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
	"github.com/stretchr/testify/assert"
	"sigs.k8s.io/yaml"
)

type testCase struct {
	Name   string                             `json:"name,omitempty"`
	Input  dw.DevWorkspaceTemplateSpecContent `json:"input,omitempty"`
	Output testOutput                         `json:"output,omitempty"`
}

type testOutput struct {
	InitContainers []dw.Component `json:"initContainers,omitempty"`
	MainContainers []dw.Component `json:"mainContainers,omitempty"`
	ErrRegexp      *string        `json:"errRegexp,omitempty"`
}

func loadTestCaseOrPanic(t *testing.T, testFilename string) testCase {
	testPath := filepath.Join("./testdata/lifecycle", testFilename)
	bytes, err := ioutil.ReadFile(testPath)
	if err != nil {
		t.Fatal(err)
	}
	var test testCase
	if err := yaml.Unmarshal(bytes, &test); err != nil {
		t.Fatal(err)
	}
	t.Log(fmt.Sprintf("Read file:\n%+v\n\n", test))
	return test
}

func TestGetInitContainers(t *testing.T) {
	tests := []testCase{
		loadTestCaseOrPanic(t, "no_events.yaml"),
		loadTestCaseOrPanic(t, "prestart_exec_command.yaml"),
		loadTestCaseOrPanic(t, "prestart_apply_command.yaml"),
		loadTestCaseOrPanic(t, "init_and_main_container.yaml"),
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			// sanity check that file reads correctly.
			assert.True(t, len(tt.Input.Components) > 0, "Input defines no components")
			gotInitContainers, gotMainComponents, err := GetInitContainers(tt.Input)
			if tt.Output.ErrRegexp != nil && assert.Error(t, err) {
				assert.Regexp(t, *tt.Output.ErrRegexp, err.Error(), "Error message should match")
			} else {
				if !assert.NoError(t, err, "Should not return error") {
					return
				}
				assert.Equal(t, tt.Output.InitContainers, gotInitContainers, "Init containers should match expected")
				assert.Equal(t, tt.Output.MainContainers, gotMainComponents, "Main containers should match expected")
			}
		})
	}
}
