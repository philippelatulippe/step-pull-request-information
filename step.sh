#!/bin/bash

THIS_SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

go run ${THIS_SCRIPT_DIR}/step.go
exit $?
