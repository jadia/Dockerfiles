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

package component

import (
	"path/filepath"

	param "github.com/ksonnet/ksonnet/metadata/params"
	"github.com/ksonnet/ksonnet/pkg/app"
	"github.com/ksonnet/ksonnet/pkg/params"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

// Delete deletes the component file and all references.
// Write operations will happen at the end to minimal-ize failures that leave
// the directory structure in a half-finished state.
func Delete(a app.App, name string) error {
	log.Debugf("deleting component %s", name)

	moduleName, componentName, err := extractPathParts(a, name)
	if err != nil {
		return err
	}

	m := NewModule(a, moduleName)

	var c Component
	components, err := m.Components()
	if err != nil {
		return err
	}
	for i := range components {
		if components[i].Name(false) == componentName {
			c = components[i]
		}
	}

	if c == nil {
		return errors.Errorf("unable to find component %q", name)
	}

	// Build the new component params.libsonnet file.
	componentParamsFile, err := afero.ReadFile(a.Fs(), m.ParamsPath())
	if err != nil {
		return err
	}
	componentJsonnet, err := param.DeleteComponent(c.Name(false), string(componentParamsFile))
	if err != nil {
		return err
	}

	// Build the new environment/<env>/params.libsonnet files.
	// environment name -> jsonnet
	envParams := make(map[string]string)
	envs, err := a.Environments()
	if err != nil {
		return err
	}
	for envName, env := range envs {
		var updated string
		updated, err = collectEnvParams(a, env, name, envName)
		if err != nil {
			return err
		}

		envParams[envName] = updated
	}

	//
	// Delete the component references.
	//
	log.Infof("Removing component parameter references ...")

	// Remove the references in component/params.libsonnet.
	log.Debugf("... deleting references in %s", m.ParamsPath())
	err = afero.WriteFile(a.Fs(), m.ParamsPath(), []byte(componentJsonnet), defaultFilePermissions)
	if err != nil {
		return err
	}

	if err = updateEnvParam(a, envs, envParams); err != nil {
		return errors.Wrap(err, "writing environment params")
	}

	//
	// Delete the component file in components/.
	//
	log.Infof("Deleting component %q", name)
	if err := c.Remove(); err != nil {
		return err
	}

	// TODO: Remove,
	// references in main.jsonnet.
	// component references in other component files (feature does not yet exist).
	log.Infof("Successfully deleted component '%s'", name)
	return nil
}

// collectEnvParams collects environment params in
func collectEnvParams(a app.App, env *app.EnvironmentConfig, componentName, envName string) (string, error) {
	log.Debugf("collecting params for environment %s", envName)
	path := filepath.Join(a.Root(), "environments", envName, "params.libsonnet")
	var envParamsFile []byte
	envParamsFile, err := afero.ReadFile(a.Fs(), path)
	if err != nil {
		return "", err
	}

	ecr := params.NewEnvComponentRemover()
	return ecr.Remove(componentName, string(envParamsFile))
}

// updateEnvParam removes the component references in each environment's
// params.libsonnet.
func updateEnvParam(a app.App, envs app.EnvironmentConfigs, envParams map[string]string) error {
	for envName := range envs {
		path := filepath.Join(a.Root(), "environments", envName, "params.libsonnet")
		log.Debugf("... deleting references in %s", path)
		if err := afero.WriteFile(a.Fs(), path, []byte(envParams[envName]), app.DefaultFilePermissions); err != nil {
			return errors.Wrapf(err, "writing params for environment %q", envName)
		}
	}

	return nil
}
