# v2 Prototype

## Build Instructions
To compile the CLI and plugin:

```bash
make build
```

## Running

```bash
# Prepare the plugin directory
# Each plugin needs to be under the plugin directory and the
# name much equal "validation component title" from the component definition
# with the type of plugin. In this case test-plugin implements both.
mkdir ./bin/plugin-dir
cp ./bin/test-plugin ./bin/plugin-dir/example

export COMPDEF_PATH=./testdata/component-definition.template.json
./bin/test-cli generate
./bin/test-cli scan
```