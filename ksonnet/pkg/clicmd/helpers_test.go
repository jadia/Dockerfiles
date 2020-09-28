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
	"github.com/ksonnet/ksonnet/pkg/app"
	"github.com/ksonnet/ksonnet/pkg/client"
	"github.com/ksonnet/ksonnet/pkg/util/test"
	"github.com/spf13/afero"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func withCmd(id initName, override actionFn, fn func()) {
	if override != nil {
		ogFn := actionFns[id]
		actionFns[id] = override

		defer func() {
			actionFns[id] = ogFn
		}()
	}

	fn()
}

type cmdTestCase struct {
	name     string
	args     []string
	action   initName
	isErr    bool
	expected map[string]interface{}
}

type stubCmdOverride struct {
	got map[string]interface{}
}

func (s *stubCmdOverride) override(m map[string]interface{}) error {
	s.got = m
	return nil
}

func runTestCmd(t *testing.T, cases []cmdTestCase) {
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := stubCmdOverride{}

			withCmd(tc.action, s.override, func() {
				fs := afero.NewMemMapFs()

				wd := "/"
				if tc.action != actionInit {
					wd = "/app"
					test.StageFile(t, fs, "app.yaml", "/app/app.yaml")
				}

				root, err := NewRoot(fs, wd, tc.args)
				require.NoError(t, err)

				err = root.Execute()
				if tc.isErr {
					require.Error(t, err)
					return
				}

				require.NoError(t, err)

				for k, v := range s.got {
					switch k {
					case actions.OptionApp:
						var expected = (*app.App)(nil)
						assert.Implements(t, expected, v)
					case actions.OptionClientConfig:
						var expected *client.Config
						assert.IsType(t, expected, v)
					case actions.OptionFs:
						var expected *afero.MemMapFs
						assert.IsType(t, expected, v)
					case actions.OptionAppRoot, actions.OptionTLSSkipVerify:
						if tc.expected[k] != nil {
							assert.Equal(t, tc.expected[k], v, "unexpected value for %q", k)
						}
					default:
						assert.Equal(t, tc.expected[k], v, "unexpected value for %q", k)
					}
				}
			})
		})
	}
}
