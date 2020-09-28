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
	"io"
	"os"
	"strings"

	param "github.com/ksonnet/ksonnet/metadata/params"
	"github.com/ksonnet/ksonnet/pkg/app"
	"github.com/ksonnet/ksonnet/pkg/component"
	"github.com/ksonnet/ksonnet/pkg/prototype"
	"github.com/ksonnet/ksonnet/pkg/registry"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/pflag"
)

// RunPrototypeUse runs `prototype use`
func RunPrototypeUse(m map[string]interface{}) error {
	pl, err := NewPrototypeUse(m)
	if err != nil {
		return err
	}

	return pl.Run()
}

// PrototypeUse generates a component from a prototype.
type PrototypeUse struct {
	app                 app.App
	args                []string
	out                 io.Writer
	packageManager      registry.PackageManager
	createComponentFn   func(app.App, string, string, string, param.Params, prototype.TemplateType) (string, error)
	bindFlagsFn         func(p *prototype.Prototype) (*pflag.FlagSet, error)
	extractParametersFn func(fs afero.Fs, p *prototype.Prototype, f *pflag.FlagSet) (map[string]string, error)
}

// NewPrototypeUse creates an instance of PrototypeUse
func NewPrototypeUse(m map[string]interface{}) (*PrototypeUse, error) {
	ol := newOptionLoader(m)

	app := ol.LoadApp()
	httpClientOpt := registry.HTTPClientOpt(ol.LoadHTTPClient())

	pl := &PrototypeUse{
		app:  app,
		args: ol.LoadStringSlice(OptionArguments),

		out:                 os.Stdout,
		packageManager:      registry.NewPackageManager(app, httpClientOpt),
		createComponentFn:   component.Create,
		bindFlagsFn:         prototype.BindFlags,
		extractParametersFn: prototype.ExtractParameters,
	}

	if ol.err != nil {
		return nil, ol.err
	}

	return pl, nil
}

// Run runs the env list action.
func (pl *PrototypeUse) Run() error {
	prototypes, err := pl.packageManager.Prototypes()
	if err != nil {
		return err
	}

	index, err := prototype.NewIndex(prototypes, prototype.DefaultBuilder)
	if err != nil {
		return err
	}

	prototypes, err = index.List()
	if err != nil {
		return err
	}

	if len(pl.args) == 0 {
		return errors.New("prototype name was not supplied as an argument")
	}

	query := pl.args[0]

	p, err := findUniquePrototype(query, prototypes)
	if err != nil {
		return err
	}

	flags, err := pl.bindFlagsFn(p)
	if err != nil {
		return errors.Wrap(err, "binding prototype flags")
	}

	if err = flags.Parse(pl.args); err != nil {
		if strings.Contains(err.Error(), "help requested") {
			return nil
		}
		return errors.Wrap(err, "parse preview args")
	}

	// Try to find the template type (if it is supplied) after the args are
	// parsed. Note that the case that `len(args) == 0` is handled at the
	// beginning of this command.
	var componentName string
	var templateType prototype.TemplateType
	if args := flags.Args(); len(args) == 1 {
		return errors.Errorf("Command is missing argument 'componentName'")
	} else if len(args) == 2 {
		componentName = args[1]
		templateType = prototype.Jsonnet
	} else if len(args) == 3 {
		componentName = args[1]
		templateType, err = prototype.ParseTemplateType(args[1])
		if err != nil {
			return err
		}
	} else {
		return errors.Errorf("Command has too many arguments (takes a prototype name and a component name)")
	}

	moduleName, err := flags.GetString("module")
	if err != nil {
		return err
	}

	mn, prototypeName := component.FromName(componentName)
	if mn != "" {
		logrus.WithField("module-name", mn).Info("Using module from component name instead of flag")
		moduleName = mn
	}

	name, err := flags.GetString("name")
	if err != nil {
		return err
	}

	if name == "" {
		if err = flags.Set("name", prototypeName); err != nil {
			return err
		}
	}

	rawParams, err := pl.extractParametersFn(pl.app.Fs(), p, flags)
	if err != nil {
		return err
	}

	text, err := expandPrototype(p, templateType, rawParams, prototypeName)
	if err != nil {
		return err
	}

	ps := param.Params{}
	for k, v := range rawParams {
		ps[k] = v
	}

	_, err = pl.createComponentFn(pl.app, moduleName, prototypeName, text, ps, templateType)
	if err != nil {
		return errors.Wrap(err, "create component")
	}

	return nil
}
