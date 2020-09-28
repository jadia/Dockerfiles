// Copyright 2018 The kubecfg authors
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

package e2e

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/blang/semver"
	ksstrings "github.com/ksonnet/ksonnet/pkg/util/strings"
	"github.com/pkg/errors"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/discovery"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	// test helpers
	// nolint: golint
	. "github.com/onsi/gomega"

	// client go auth plugins
	_ "k8s.io/client-go/plugin/pkg/client/auth/azure"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	_ "k8s.io/client-go/plugin/pkg/client/auth/openstack"
)

var kubeconfig = flag.String("kubeconfig", "", "absolute path to kubeconfig file")
var kubectx = flag.String("context", "", "kubernetes context")

type e2e struct {
	root       string
	restConfig *rest.Config

	corev1 corev1.CoreV1Client
}

func newE2e() *e2e {
	dir, err := ioutil.TempDir("", "")
	Expect(err).ToNot(HaveOccurred())

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	Expect(err).ToNot(HaveOccurred(), "build client config from flags")

	e := &e2e{
		root:       dir,
		restConfig: config,
	}

	return e
}

func (e *e2e) serverVersion() string {
	discoveryClient, err := discovery.NewDiscoveryClientForConfig(e.restConfig)
	Expect(err).ToNot(HaveOccurred(), "create discovery client")

	versionInfo, err := discoveryClient.ServerVersion()
	Expect(err).ToNot(HaveOccurred(), "retrieve server version info")

	v, err := semver.Make(strings.TrimPrefix(versionInfo.String(), "v"))
	Expect(err).ToNot(HaveOccurred(), "parse server version")

	return fmt.Sprintf("v%d.%d.%d",
		v.Major, v.Minor, v.Patch)
}

func (e *e2e) createNamespace() string {
	name := fmt.Sprintf("ks-e2e-%s", ksstrings.LowerRand(6))

	c, err := corev1.NewForConfig(e.restConfig)
	Expect(err).ToNot(HaveOccurred())

	result, err := c.Namespaces().Create(
		&v1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: name,
			},
		},
	)
	Expect(err).ToNot(HaveOccurred())

	return result.GetName()
}

func (e *e2e) removeNamespace(name string) {
	c, err := corev1.NewForConfig(e.restConfig)
	Expect(err).ToNot(HaveOccurred())

	err = c.Namespaces().Delete(name, &metav1.DeleteOptions{})
	Expect(err).ToNot(HaveOccurred())
}

func (e *e2e) close() {
	err := os.RemoveAll(e.root)
	Expect(err).ToNot(HaveOccurred())
}

func (e *e2e) wd() string {
	wd, err := os.Getwd()
	ExpectWithOffset(1, err).ToNot(HaveOccurred())

	return wd
}

func (e *e2e) ksBin() string {
	return filepath.Join(e.root, "ks")
}

func (e *e2e) ks(args ...string) *output {
	cmd := exec.Command(e.ksBin(), args...)
	o, err := runWithOutput(cmd)
	ExpectWithOffset(1, err).ToNot(HaveOccurred())
	return o
}

func (e *e2e) ksInApp(appDir string, args ...string) *output {
	ExpectWithOffset(1, appDir).To(BeADirectory())
	cmd := exec.Command(e.ksBin(), args...)
	cmd.Dir = appDir
	o, err := runWithOutput(cmd)
	ExpectWithOffset(1, err).ToNot(HaveOccurred())
	return o
}

func (e *e2e) buildKs() {
	args := []string{
		"build",
		"-o",
		e.ksBin(),
		// TODO: this needs to come from elsewhere
		`-ldflags="-X main.apimachineryVersion=kubernetes-1.10.4"`,
		`github.com/ksonnet/ksonnet/cmd/ks`,
	}

	cmd := exec.Command("go", args...)

	o, err := runWithOutput(cmd)
	ExpectWithOffset(1, err).ToNot(HaveOccurred())
	assertExitStatus(o, 0)
}

func (e *e2e) initApp(options *initOptions) app {
	appID := ksstrings.Rand(6)
	appDir := filepath.Join(e.root, appID)

	opts := []string{
		"init",
		appID,
		"--dir",
		appDir,
	}

	if options == nil {
		options = &initOptions{}
	}

	io, err := options.toSlice()
	ExpectWithOffset(1, err).ToNot(HaveOccurred())
	opts = append(opts, io...)

	o := e.ks(opts...)
	assertExitStatus(o, 0)
	return app{dir: appDir, e2e: e}
}

type initOptions struct {
	context        string
	envName        string
	namespace      string
	server         string
	skipRegistries bool
}

func (o *initOptions) toSlice() ([]string, error) {
	if o.server != "" && o.context != "" {
		return nil, errors.Errorf("can't specify server and context")
	}

	var options []string

	switch {
	case o.server == "" && o.context == "":
		options = append(options, "--server", "http://example.com")
	case o.server != "":
		options = append(options, "--server", o.server)
	case o.context != "":
		options = append(options, "--context", o.context)
	}

	if o.namespace != "" {
		options = append(options, "--namespace", o.namespace)
	}

	if o.skipRegistries {
		options = append(options, "--skip-default-registries")
	}

	if o.envName != "" {
		options = append(options, "--env", o.envName)
	}

	return options, nil
}
