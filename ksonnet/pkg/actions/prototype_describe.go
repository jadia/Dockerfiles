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
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/ksonnet/ksonnet/pkg/app"
	"github.com/ksonnet/ksonnet/pkg/pkg"
	"github.com/ksonnet/ksonnet/pkg/prototype"
	"github.com/ksonnet/ksonnet/pkg/registry"
	"github.com/pkg/errors"
)

// RunPrototypeDescribe runs `prototype describe`
func RunPrototypeDescribe(m map[string]interface{}) error {
	pd, err := NewPrototypeDescribe(m)
	if err != nil {
		return err
	}

	return pd.Run()
}

// PrototypeDescribe describes a prototype.
type PrototypeDescribe struct {
	app            app.App
	out            io.Writer
	query          string
	packageManager registry.PackageManager
}

// NewPrototypeDescribe creates an instance of PrototypeDescribe
func NewPrototypeDescribe(m map[string]interface{}) (*PrototypeDescribe, error) {
	ol := newOptionLoader(m)

	app := ol.LoadApp()
	httpClientOpt := registry.HTTPClientOpt(ol.LoadHTTPClient())

	pd := &PrototypeDescribe{
		app:   app,
		query: ol.LoadString(OptionQuery),

		out:            os.Stdout,
		packageManager: registry.NewPackageManager(app, httpClientOpt),
	}

	if ol.err != nil {
		return nil, ol.err
	}

	return pd, nil
}

// Run runs the env list action.
func (pd *PrototypeDescribe) Run() error {
	prototypes, err := pd.packageManager.Prototypes()
	if err != nil {
		return err
	}

	index, err := prototype.NewIndex(prototypes, prototype.DefaultBuilder)
	if err != nil {
		return errors.Wrap(err, "create prototype index")
	}

	prototypes, err = index.List()
	if err != nil {
		return err
	}

	p, err := findUniquePrototype(pd.query, prototypes)
	if err != nil {
		return err
	}

	fmt.Fprintln(pd.out, `PROTOTYPE NAME:`)
	fmt.Fprintln(pd.out, p.Name)
	fmt.Fprintln(pd.out)
	fmt.Fprintln(pd.out, `DESCRIPTION:`)
	fmt.Fprintln(pd.out, p.Template.Description)
	fmt.Fprintln(pd.out)
	fmt.Fprintln(pd.out, `REQUIRED PARAMETERS:`)
	fmt.Fprintln(pd.out, p.RequiredParams().PrettyString("  "))
	fmt.Fprintln(pd.out)
	fmt.Fprintln(pd.out, `OPTIONAL PARAMETERS:`)
	fmt.Fprintln(pd.out, p.OptionalParams().PrettyString("  "))
	fmt.Fprintln(pd.out)
	fmt.Fprintln(pd.out, `TEMPLATE TYPES AVAILABLE:`)
	fmt.Fprintf(pd.out, "  %s\n", p.Template.AvailableTemplates())

	return nil
}

type prototypeFn func(app.App, pkg.Descriptor) (prototype.Prototypes, error)

func findUniquePrototype(query string, prototypes prototype.Prototypes) (*prototype.Prototype, error) {
	index, err := prototype.NewIndex(prototypes, prototype.DefaultBuilder)
	if err != nil {
		return nil, err
	}

	sameSuffix, err := index.SearchNames(query, prototype.Suffix)
	if err != nil {
		return nil, err
	}

	if len(sameSuffix) == 1 {
		// Success.
		return sameSuffix[0], nil
	} else if len(sameSuffix) > 1 {
		// Ambiguous match.
		names := specNames(sameSuffix)
		return nil, errors.Errorf("ambiguous match for '%s': %s", query, strings.Join(names, ", "))
	} else {
		// No matches.
		substring, err := index.SearchNames(query, prototype.Substring)
		if err != nil || len(substring) == 0 {
			return nil, errors.Errorf("no prototype names matched '%s'", query)
		}

		partialMatches := specNames(substring)
		partials := strings.Join(partialMatches, "\n")
		return nil, errors.Errorf("no prototype names matched '%s'; a list of partial matches: %s", query, partials)
	}
}

func specNames(prototypes []*prototype.Prototype) []string {
	partialMatches := []string{}
	for _, proto := range prototypes {
		partialMatches = append(partialMatches, proto.Name)
	}

	return partialMatches
}
