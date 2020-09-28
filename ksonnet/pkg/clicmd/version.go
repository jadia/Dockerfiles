// Copyright 2017 The kubecfg authors
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
	"fmt"

	jsonnet "github.com/google/go-jsonnet"
	"github.com/spf13/cobra"
)

var (
	versionLong = `
The ` + "`version`" + ` command prints out version info about the current ksonnet CLI,
as well as for any of its helper libraries (e.g. ` + "`client-go`" + `).

### Syntax
`
)

// Version is overridden by main
var Version = "(dev build)"

// APImachineryVersion is overriden by main
var APImachineryVersion = ""

func newVersionCmd() *cobra.Command {
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print version information for this ksonnet binary",
		Long:  versionLong,
		Run: func(cmd *cobra.Command, args []string) {
			out := cmd.OutOrStdout()

			fmt.Fprintln(out, "ksonnet version:", Version)
			fmt.Fprintln(out, "jsonnet version:", jsonnet.Version())
			fmt.Fprintln(out, "client-go version:", APImachineryVersion)
		},
	}

	return versionCmd

}
