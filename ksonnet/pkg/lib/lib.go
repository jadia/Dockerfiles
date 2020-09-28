// Copyright 2017 The ksonnet authors
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

package lib

import (
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"

	"github.com/ksonnet/ksonnet/pkg/util/kslib"
)

const (
	schemaFilename = "swagger.json"
	k8sLibFilename = "k8s.libsonnet"

	// ExtensionsLibFilename is the file name with the contents of the
	// generated ksonnet-lib
	ExtensionsLibFilename = "k.libsonnet"

	// KsonnetLibHome is name of the directory ksonnet lib files will be generated in.
	KsonnetLibHome = "ksonnet-lib"
)

// KsLibGenerator generates ksonnet-lib.
type KsLibGenerator interface {
	Generate() (*kslib.KsonnetLib, error)
}

type defaultKsLibGenerator struct {
	spec ClusterSpec
}

func (g *defaultKsLibGenerator) Generate() (*kslib.KsonnetLib, error) {
	if g.spec == nil {
		return nil, errors.Errorf("uninitialized ClusterSpec")
	}

	b, err := g.spec.OpenAPI()
	if err != nil {
		return nil, err
	}

	return kslib.Ksonnet(b)
}

// Manager operates on the files in the lib directory of a ksonnet project.
// This included generating the ksonnet-lib files needed to compile the
// ksonnet (jsonnet) code.
type Manager struct {
	// K8sVersion is the Kubernetes version of the Open API spec.
	K8sVersion string

	libPath string
	fs      afero.Fs

	generator KsLibGenerator
}

// NewManager creates a new instance of lib.Manager
func NewManager(k8sSpecFlag string, fs afero.Fs, libPath string, httpClient *http.Client) (*Manager, error) {
	//
	// Generate the program text for ksonnet-lib.
	//
	spec, err := ParseClusterSpec(k8sSpecFlag, fs, httpClient)
	if err != nil {
		return nil, err
	}

	version, err := spec.Version()
	if err != nil {
		return nil, err
	}

	return &Manager{
		K8sVersion: version,
		fs:         fs,
		libPath:    libPath,
		generator:  &defaultKsLibGenerator{spec: spec},
	}, nil
}

// GenerateLibData will generate the swagger and ksonnet-lib files in the lib
// directory of a ksonnet project. The swagger and ksonnet-lib files are
// unique to each Kubernetes API version. If the files already exist for a
// specific Kubernetes API version, they won't be re-generated here.
func (m *Manager) GenerateLibData() error {
	kl, err := m.generator.Generate()
	if err != nil {
		return err
	}

	genPath := filepath.Join(m.ksLibDir(), m.K8sVersion)

	ok, err := afero.DirExists(m.fs, genPath)
	if err != nil {
		return err
	}
	if ok {
		// Already have lib data for this k8s api version
		return nil
	}

	err = m.fs.MkdirAll(genPath, os.FileMode(0755))
	if err != nil {
		return err
	}

	files := []struct {
		path string
		data []byte
	}{
		{
			// schema file
			filepath.Join(genPath, schemaFilename),
			kl.Swagger,
		},
		{
			// k8s file
			filepath.Join(genPath, k8sLibFilename),
			kl.K8s,
		},
		{
			// extensions file
			filepath.Join(genPath, ExtensionsLibFilename),
			kl.K,
		},
	}

	log.Infof("Generating ksonnet-lib data at path '%s'", genPath)

	for _, a := range files {
		fileName := path.Base(string(a.path))
		if err = afero.WriteFile(m.fs, string(a.path), a.data, os.FileMode(0644)); err != nil {
			log.Debugf("Failed to write '%s'", fileName)
			return err
		}
	}

	return nil
}

// GetLibPath returns the absolute path pointing to the directory with the
// metadata files for the provided k8sVersion.
func (m *Manager) GetLibPath() (string, error) {
	basePath := m.ksLibDir()

	ok, err := afero.DirExists(m.fs, basePath)
	if err != nil {
		return "", err
	}

	if !ok {
		log.Debugf("Expected lib directory '%s' but was not found", m.K8sVersion)

		// create the directory
		if err = m.GenerateLibData(); err != nil {
			return "", err
		}

	}
	return filepath.Join(basePath, m.K8sVersion), nil
}

func (m *Manager) ksLibDir() string {
	ksLibPath := filepath.Join(m.libPath, m.K8sVersion)
	exists, _ := afero.IsDir(m.fs, ksLibPath)

	if exists {
		return m.libPath
	}

	return filepath.Join(m.libPath, "ksonnet-lib")
}
