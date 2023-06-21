#!/usr/bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

set -e
"${DIR}"/install_go.sh
"${DIR}"/setup_pmdk.sh
