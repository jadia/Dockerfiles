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

package actions

import (
	"testing"

	"github.com/ksonnet/ksonnet/pkg/app"
	amocks "github.com/ksonnet/ksonnet/pkg/app/mocks"
	"github.com/ksonnet/ksonnet/pkg/component"
	cmocks "github.com/ksonnet/ksonnet/pkg/component/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParamSet(t *testing.T) {
	withApp(t, func(appMock *amocks.App) {
		componentName := "deployment"
		path := "replicas"
		value := "3"

		c := &cmocks.Component{}
		c.On("SetParam", []string{"replicas"}, 3).Return(nil)

		in := map[string]interface{}{
			OptionApp:   appMock,
			OptionName:  componentName,
			OptionPath:  path,
			OptionValue: value,
		}

		a, err := NewParamSet(in)
		require.NoError(t, err)

		a.resolvePathFn = func(app.App, string) (component.Module, component.Component, error) {
			return nil, c, nil
		}

		err = a.Run()
		require.NoError(t, err)
	})
}

func TestParamSet_invalid_component(t *testing.T) {
	withApp(t, func(appMock *amocks.App) {
		componentName := "deployment"
		path := "replicas"
		value := "3"

		in := map[string]interface{}{
			OptionApp:   appMock,
			OptionName:  componentName,
			OptionPath:  path,
			OptionValue: value,
		}

		a, err := NewParamSet(in)
		require.NoError(t, err)

		a.resolvePathFn = func(app.App, string) (component.Module, component.Component, error) {
			return nil, nil, nil
		}

		err = a.Run()
		require.Error(t, err)
	})
}

func TestParamSet_asString(t *testing.T) {
	withApp(t, func(appMock *amocks.App) {
		componentName := "deployment"
		path := "label"
		value := "3"

		c := &cmocks.Component{}
		c.On("SetParam", []string{"label"}, "3").Return(nil)

		in := map[string]interface{}{
			OptionApp:      appMock,
			OptionName:     componentName,
			OptionPath:     path,
			OptionValue:    value,
			OptionAsString: true,
		}

		a, err := NewParamSet(in)
		require.NoError(t, err)

		a.resolvePathFn = func(app.App, string) (component.Module, component.Component, error) {
			return nil, c, nil
		}

		err = a.Run()
		require.NoError(t, err)
	})
}

func TestParamSet_resolveImage(t *testing.T) {
	withApp(t, func(appMock *amocks.App) {
		componentName := "deployment"
		path := "image"
		value := "foo/bar:latest"

		c := &cmocks.Component{}
		c.On("SetParam", []string{"image"}, "foo/bar@sha256:abcde").Return(nil)

		in := map[string]interface{}{
			OptionApp:          appMock,
			OptionName:         componentName,
			OptionPath:         path,
			OptionValue:        value,
			OptionResolveImage: true,
		}

		a, err := NewParamSet(in)
		require.NoError(t, err)

		a.resolvePathFn = func(app.App, string) (component.Module, component.Component, error) {
			return nil, c, nil
		}

		a.resolveImageFn = func(string) (string, error) {
			return "foo/bar@sha256:abcde", nil
		}

		err = a.Run()
		require.NoError(t, err)
	})
}

func TestParamSet_global(t *testing.T) {
	withApp(t, func(appMock *amocks.App) {
		module := "/"
		path := "replicas"
		value := "3"

		m := &cmocks.Module{}
		m.On("SetParam", []string{"replicas"}, 3).Return(nil)

		in := map[string]interface{}{
			OptionApp:    appMock,
			OptionName:   module,
			OptionPath:   path,
			OptionValue:  value,
			OptionGlobal: true,
		}

		a, err := NewParamSet(in)
		require.NoError(t, err)

		a.getModuleFn = func(app.App, string) (component.Module, error) {
			return m, nil
		}

		err = a.Run()
		require.NoError(t, err)
	})
}

func TestParamSet_env(t *testing.T) {
	withApp(t, func(appMock *amocks.App) {
		name := "deployment"
		path := "replicas"
		value := "3"

		in := map[string]interface{}{
			OptionApp:     appMock,
			OptionName:    name,
			OptionPath:    path,
			OptionValue:   value,
			OptionEnvName: "default",
		}

		a, err := NewParamSet(in)
		require.NoError(t, err)

		envSetter := func(ksApp app.App, envName, name, pName, value string) error {
			assert.Equal(t, "default", envName)
			assert.Equal(t, "deployment", name)
			assert.Equal(t, "replicas", pName)
			assert.Equal(t, "3", value)
			return nil
		}
		a.setEnvFn = envSetter

		err = a.Run()
		require.NoError(t, err)
	})
}

func TestParamSet_env_resolveImage(t *testing.T) {
	withApp(t, func(appMock *amocks.App) {
		name := "deployment"
		path := "image"
		value := "foo/bar:latest"

		in := map[string]interface{}{
			OptionApp:          appMock,
			OptionName:         name,
			OptionPath:         path,
			OptionValue:        value,
			OptionEnvName:      "default",
			OptionResolveImage: true,
		}

		a, err := NewParamSet(in)
		require.NoError(t, err)

		envSetter := func(ksApp app.App, envName, name, pName, value string) error {
			assert.Equal(t, "default", envName)
			assert.Equal(t, "deployment", name)
			assert.Equal(t, "image", pName)
			assert.Equal(t, "foo/bar@sha256:abcde", value)
			return nil
		}
		a.setEnvFn = envSetter
		a.resolveImageFn = func(string) (string, error) {
			return "foo/bar@sha256:abcde", nil
		}

		err = a.Run()
		require.NoError(t, err)
	})
}

func TestParamSet_envGlobal(t *testing.T) {
	withApp(t, func(appMock *amocks.App) {
		path := "replicas"
		value := "3"

		in := map[string]interface{}{
			OptionApp:     appMock,
			OptionPath:    path,
			OptionValue:   value,
			OptionEnvName: "default",
		}

		a, err := NewParamSet(in)
		require.NoError(t, err)

		envSetter := func(ksApp app.App, envName, pName, value string) error {
			assert.Equal(t, "default", envName)
			assert.Equal(t, "replicas", pName)
			assert.Equal(t, "3", value)
			return nil
		}
		a.setGlobalEnvFn = envSetter

		err = a.Run()
		require.NoError(t, err)
	})
}

func TestParamSet_requires_app(t *testing.T) {
	in := make(map[string]interface{})
	_, err := NewParamSet(in)
	require.Error(t, err)
}
