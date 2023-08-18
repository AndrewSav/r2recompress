#!/usr/bin/env pwsh

docker build --progress plain -t "build" .

$pwd = (get-item $PSScriptRoot ).Parent.FullName
docker run -it --rm --name build -v $pwd`:/src -e TERM=dumb build pwsh -c "git config --global --add safe.directory /src;cd /src;./build.ps1"
