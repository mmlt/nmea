#!/bin/bash
# 
# Generate sentences.go
#
# Command:
#   cd pkg/parser/ && ./generate.sh ; cd -

# TODO make this a for loop with intermediate file, on error the last 'good' yaml is in the file for trouble-shooting

cat ../../spec/spec.yaml \
| yq -y -f 01-spec-add-type.jq \
| yq -y -f 02-spec-add-xarg.jq >_spec.yaml

#more _spec.yaml
gomplate -d spec=_spec.yaml -f sentences.tmpl >sentences.go