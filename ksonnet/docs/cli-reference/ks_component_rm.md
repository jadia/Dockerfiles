## ks component rm

Delete a component from the ksonnet application

### Synopsis

Delete a component from the ksonnet application. This is equivalent to deleting the
component file in the components directory and cleaning up all component
references throughout the project.

```
ks component rm <component-name> [flags]
```

### Examples

```

# List all components
ks component list
```

### Options

```
  -h, --help   help for rm
```

### Options inherited from parent commands

```
      --dir string        Ksonnet application root to use; Defaults to CWD
      --tls-skip-verify   Skip verification of TLS server certificates
  -v, --verbose count     Increase verbosity. May be given multiple times.
```

### SEE ALSO

* [ks component](ks_component.md)	 - Manage ksonnet components

