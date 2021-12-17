#!/bin/bash

# files=$(find . -name "*.go")
for file in $(find . -name "*.go"); do
  echo "$file"
  go vet "$file"
done
