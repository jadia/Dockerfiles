#!/bin/sh
# Shell script to convert YAML files to ksonnet
# Author: nitish@jadia.dev
# Version: v0.1
# Date: 2020-08-25

if [[ $# -ne 2 ]]; then
    echo """Usage: docker run --rm -v "\${PWD}":/workdir yaml2ksonnet <name of yaml file> <name of component in ksonnet> """
    exit 1
fi

inputfile=$1
workdir="/workdir"
# Source: https://stackoverflow.com/a/50710152
outputfile=$(echo ${inputfile%.*}.jsonnet)
outputfile="$workdir/$outputfile"
echo $outputfile
echo """local env = std.extVar(\"__ksonnet/environments\");
local params = std.extVar(\"__ksonnet/params\").components[\"$2\"];
[""" > $outputfile
# Sed command source: https://gist.github.com/sv99/6852cc2e2a09bd3a68ed
yq r -d'*' --prettyPrint -j $inputfile | sed -e ':a;N;$!ba;s/}\n{/},\n{/g' >> $outputfile
echo "]" >> $outputfile
