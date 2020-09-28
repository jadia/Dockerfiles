## ks prototype search

Search for a prototype

### Synopsis


The `prototype search` command allows you to search for specific prototypes by name.
Specifically, it matches any prototypes with names that contain the string <name-substring>.

### Related Commands

* `ks prototype describe` — See more info about a prototype's output and usage
* `ks prototype list` — List all locally available ksonnet prototypes

### Syntax


```
ks prototype search <name-substring> [flags]
```

### Examples

```

# Search for prototypes with names that contain the string 'service'.
ks prototype search service
```

### Options

```
  -h, --help            help for search
  -o, --output string   Output format. Valid options: table|json
```

### Options inherited from parent commands

```
      --dir string        Ksonnet application root to use; Defaults to CWD
      --tls-skip-verify   Skip verification of TLS server certificates
  -v, --verbose count     Increase verbosity. May be given multiple times.
```

### SEE ALSO

* [ks prototype](ks_prototype.md)	 - Instantiate, inspect, and get examples for ksonnet prototypes

