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
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/ksonnet/ksonnet/pkg/app"
	"github.com/ksonnet/ksonnet/pkg/parts"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

var separator = fmt.Sprintf("%c", filepath.Separator)

// Fs is a registry based on a local filesystem.
type Fs struct {
	app  app.App
	spec *app.RegistryConfig
	root string
}

// NewFs creates an instance of Fs. Assign a name to the RegistryConfig if you want
// the Fs to know it's name.
func NewFs(a app.App, registryRef *app.RegistryConfig) (*Fs, error) {
	fs := &Fs{
		app:  a,
		spec: registryRef,
	}

	u, err := url.Parse(registryRef.URI)
	if err != nil {
		return nil, err
	}

	fs.root = u.Path

	return fs, nil
}

var _ Registry = (*Fs)(nil)

// Name is the registry name.
func (fs *Fs) Name() string {
	return fs.spec.Name
}

// Protocol is the registry protocol.
func (fs *Fs) Protocol() Protocol {
	return Protocol(fs.spec.Protocol)
}

// URI is the registry URI.
func (fs *Fs) URI() string {
	return fs.spec.URI
}

// RegistrySpecDir is the registry directory.
func (fs *Fs) RegistrySpecDir() string {
	return fs.root
}

// RegistrySpecFilePath is the path for the registry.yaml
func (fs *Fs) RegistrySpecFilePath() string {
	logrus.WithFields(logrus.Fields{
		"fs-registry-root": fs.root,
	}).Debug("creating registry config file path")

	return filepath.Join(fs.root, registryYAMLFile)
}

// FetchRegistrySpec fetches the registry spec.
func (fs *Fs) FetchRegistrySpec() (*Spec, error) {
	logrus.WithFields(logrus.Fields{
		"file-path": fs.RegistrySpecFilePath(),
	}).Debug("fetching registry spec")

	configPath := fs.RegistrySpecFilePath()
	path := filepath.Join(fs.app.Root(), configPath)
	if filepath.IsAbs(configPath) {
		var err error
		path, err = filepath.Rel(fs.app.Root(), configPath)
		if err != nil {
			return nil, errors.Wrap(err, "unable to create path from relative filesystem path")
		}
		path = filepath.Join(fs.app.Root(), path)
	}

	data, err := afero.ReadFile(fs.app.Fs(), path)
	if err != nil {
		return nil, err
	}

	return Unmarshal(data)
}

// MakeRegistryConfig returns an app registry ref spec.
func (fs *Fs) MakeRegistryConfig() *app.RegistryConfig {
	return fs.spec
}

// ResolveLibrarySpec returns a resolved spec for a part. `libRefSpec` is ignored.
func (fs *Fs) ResolveLibrarySpec(partName, libRefSpec string) (*parts.Spec, error) {
	partRoot := filepath.Join(fs.RegistrySpecDir(), partName, partsYAMLFile)
	data, err := afero.ReadFile(fs.app.Fs(), partRoot)
	if err != nil {
		return nil, err
	}

	return parts.Unmarshal(data)
}

// ResolveLibrary fetches the part and creates a parts spec and library ref spec.
func (fs *Fs) ResolveLibrary(partName, partAlias, libRefSpec string, onFile ResolveFile, onDir ResolveDirectory) (*parts.Spec, *app.LibraryConfig, error) {
	if partAlias == "" {
		partAlias = partName
	}

	partRoot := filepath.Join(fs.RegistrySpecDir(), partName)

	partsSpec, err := fs.ResolveLibrarySpec(partName, libRefSpec)
	if err != nil {
		return nil, nil, err
	}

	err = afero.Walk(fs.app.Fs(), partRoot, func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Make path relative to registry root
		libPath := strings.TrimPrefix(
			strings.TrimPrefix(path, fs.URI()),
			separator,
		)
		if fi.IsDir() {
			return onDir(libPath)
		}

		data, err := afero.ReadFile(fs.app.Fs(), path)
		if err != nil {
			return err
		}
		return onFile(libPath, data)
	})

	if err != nil {
		return nil, nil, err
	}

	refSpec := &app.LibraryConfig{
		Name:     partAlias,
		Registry: fs.Name(),
	}

	return partsSpec, refSpec, nil
}

// CacheRoot combines the path with the registry name.
func (fs *Fs) CacheRoot(name, path string) (string, error) {
	return filepath.Join(name, path), nil
}

// SetURI implements registry.Setter. It sets the URI for the registry.
func (fs *Fs) SetURI(uri string) error {
	if fs == nil {
		return errors.Errorf("nil receiver")
	}
	if fs.spec == nil {
		return errors.Errorf("nil spec")
	}

	// Validate
	if ok, err := fs.ValidateURI(uri); err != nil || !ok {
		return errors.Wrap(err, "validating uri")
	}

	u, err := url.Parse(uri)
	if err != nil {
		return errors.Wrapf(err, "parsing uri: %v", uri)
	}

	fs.root = u.Path
	fs.spec.URI = uri

	return nil
}

// ValidateURI implements registry.Validator. A URI is valid if:
//   * It is a valid URI (RFC 3986)
//   * It points to an existing directory on disk
//   * That directory contains a `registry.yaml` file
func (fs *Fs) ValidateURI(uri string) (bool, error) {
	if fs == nil {
		return false, errors.Errorf("nil receiver")
	}
	if fs.app == nil {
		return false, errors.Errorf("app is nil")
	}

	aferoFs := fs.app.Fs()
	if aferoFs == nil {
		return false, errors.Errorf("afero.Fs interface is nil")
	}

	u, err := url.Parse(uri)
	if err != nil {
		return false, err
	}

	path := u.Path

	if !filepath.IsAbs(path) {
		path = filepath.Join(fs.app.Root(), path)
	}

	exists, err := afero.DirExists(aferoFs, path)
	if err != nil {
		return false, err
	}

	if !exists {
		return false, errors.Errorf("registry path %q does not exist", path)
	}

	return true, nil
}
