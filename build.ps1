#!/usr/bin/env pwsh

$appname = "r2recompress"
$targetRepo = "AndrewSav/$appname"

$ErrorActionPreference = "Stop"
"Getting build information..." | Write-Host
$r = [regex]"[0-9]+\.[0-9]+(\.[0-9]+)?"
$m =  $r.Match((go version))
if ($LASTEXITCODE) { exit 1 }
if (!$m) {
  "Could not get go version" | Write-Error
  exit 1
}
$goVersion = $m.Value
$package = "github.com/$targetRepo/internal/version"
$timestamp = (Get-Date -Format "dddd, d MMMM yyyy HH:mm:ss K")
$version = (git describe --tags --match=v[0-9]*\.[0-9]*\.[0-9]*)
if ($LASTEXITCODE) { exit 1 }
$version = $version.Substring(1)
$changed = (git status --porcelain)
if ($changed) {
  $version = "dev"
}

$ldflags=@(
  "-X '$package.Version=$version'"
  "-X '$package.GoVersion=$goVersion'"
  "-X '$package.BuildTime=$timestamp'"
)

"Version $version built on $timestamp (go $goVersion)" | Write-Host

go mod tidy
if ($LASTEXITCODE) { exit 1 }

$haveGoimports = Get-Command "goimports" -ErrorAction SilentlyContinue
if ($haveGoimports) {
  "Running goimports..." | Write-Host
  goimports -w .
} else {
  "Running go fmt..." | Write-Host
  go fmt ./...
}
if ($LASTEXITCODE) { exit 1 }

"Running go build..." | Write-Host
$env:GOOS="windows"
$env:GOARCH="amd64"
go build -ldflags="$ldflags" cmd/$appname/$appname.go
if ($LASTEXITCODE) { exit 1 }
"Executable is built successfully" | Write-Host

