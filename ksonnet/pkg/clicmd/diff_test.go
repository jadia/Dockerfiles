// Copyright 2018 The ksonnet authors
//
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

package clicmd

import (
	"testing"

	"github.com/ksonnet/ksonnet/pkg/actions"
)

func Test_diffCmd(t *testing.T) {
	cases := []cmdTestCase{
		{
			name:   "diff two environments",
			args:   []string{"diff", "env1", "env2"},
			action: actionDiff,
			expected: map[string]interface{}{
				actions.OptionApp:            nil,
				actions.OptionClientConfig:   nil,
				actions.OptionSrc1:           "env1",
				actions.OptionSrc2:           "env2",
				actions.OptionComponentNames: []string{},
			},
		},
		{
			name:  "no args",
			args:  []string{"diff"},
			isErr: true,
		},
		{
			name:  "too many args",
			args:  []string{"diff", "env1", "env2", "env3"},
			isErr: true,
		},
		{
			name:  "invalid jsonnet flag",
			args:  []string{"diff", "default", "--ext-str", "foo"},
			isErr: true,
		},
	}

	runTestCmd(t, cases)
}
