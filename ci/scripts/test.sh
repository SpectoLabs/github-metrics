#!/bin/bash

set -e -x

go get github.com/stretchr/testify/assert

cd $(dirname $0)/../..

go test
