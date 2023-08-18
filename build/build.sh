#!/bin/bash

set -e

docker build --progress plain -t "build" .

SCRIPT_DIR="$( cd -- "$( dirname -- "${BASH_SOURCE[0]:-$0}"; )" &> /dev/null && pwd 2> /dev/null; )";
pwd="$(dirname "$SCRIPT_DIR")"
docker run -t --rm --name build -v $pwd:/src -e TERM=dumb build pwsh -c "git config --global --add safe.directory /src; cd /src;./build.ps1"
