## ks registry

Manage registries for current project

### Synopsis


A ksonnet registry is basically a repository for *packages*. (Registry here is
used in the same sense as a container image registry). Registries are identified
by a `registry.yaml` in their root that declares which packages they contain.

Specifically, registries contain a set of versioned packages that the user can
install and manage in a given ksonnet app, using the CLI. A typical package contains:

1. **A library definining a set of "parts"**. These are pre-fabricated API objects
which can be combined together to configure a Kubernetes application for some task.
(e.g. a Deployment, a Service, and a Secret, specifically tailored for Redis).

2. **A set of "prototypes"**, which are pre-fabricated combinations of parts, as
described above. (See `ks prototype --help` for more information.)

----


### Options

```
  -h, --help   help for registry
```

### Options inherited from parent commands

```
      --dir string        Ksonnet application root to use; Defaults to CWD
      --tls-skip-verify   Skip verification of TLS server certificates
  -v, --verbose count     Increase verbosity. May be given multiple times.
```

### SEE ALSO

* [ks](ks.md)	 - Configure your application to deploy to a Kubernetes cluster
* [ks registry add](ks_registry_add.md)	 - Add a registry to the current ksonnet app
* [ks registry describe](ks_registry_describe.md)	 - Describe a ksonnet registry and the packages it contains
* [ks registry list](ks_registry_list.md)	 - List all registries known to the current ksonnet app
* [ks registry set](ks_registry_set.md)	 - Set configuration options for registry

