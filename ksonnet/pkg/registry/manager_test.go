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

package registry

import (
	"testing"

	"github.com/ksonnet/ksonnet/pkg/app"
	"github.com/ksonnet/ksonnet/pkg/app/mocks"
	"github.com/ksonnet/ksonnet/pkg/util/github"
	ghmocks "github.com/ksonnet/ksonnet/pkg/util/github/mocks"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_List(t *testing.T) {
	withApp(t, func(a *mocks.App, fs afero.Fs) {
		c := &ghmocks.GitHub{}
		c.On("ValidateURL", mock.Anything).Return(nil)
		c.On("CommitSHA1", mock.Anything, github.Repo{Org: "ksonnet", Repo: "parts"}, mock.AnythingOfType("string")).
			Return("12345", nil)

		ghcOpt := GitHubClient(c)
		githubFactory = func(a app.App, spec *app.RegistryConfig, opts ...GitHubOpt) (*GitHub, error) {
			return NewGitHub(a, spec, ghcOpt)
		}

		specs := app.RegistryConfigs{
			"incubator": &app.RegistryConfig{
				Protocol: string(ProtocolGitHub),
				URI:      "github.com/ksonnet/parts/tree/master/incubator",
			},
		}

		appMock := &mocks.App{}
		appMock.On("Registries").Return(specs, nil)

		registries, err := List(appMock, nil)
		require.NoError(t, err)

		require.Len(t, registries, 1)

	})
}
