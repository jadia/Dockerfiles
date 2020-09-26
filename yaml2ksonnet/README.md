# Yaml2Ksonnet

Convert a single file with multiple YAML documents to a single .jsonnet file which is supported by ksonnet.

[ks import](https://ksonnet.io/docs/examples/import-yaml/) does not support YAML file with multiple YAML documents. That is why this method is needed.

## Usage

```bash
docker run --rm -v "${PWD}":/workdir ghcr.io/jadia/Dockerfiles/yaml2ksonnet:latest <name of yaml file> <name of component in ksonnet> 
```

Example:

```bash
docker run --rm -v "${PWD}":/workdir ghcr.io/jadia/Dockerfiles/yaml2ksonnet:latest ./jenkins.yaml jenkins
```

Then move the `jenkins.jsonnet` file to `components/jenkins.jsonnet`.   

Put entry in `components/params.libsonnet`

The entry must look something like this:

```json
{
    global: {},
    componenets: {
        "jenkins": {
        // some data
        },
    },
}
```
