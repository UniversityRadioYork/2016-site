#!/bin/bash

DEPENDENCIES=(npm go)
REQUIRED_GO_VERSION=1.6

function is_installed() {
  command -v $1 >/dev/null 2>&1
}

function check_dependencies() {
  for dependency in "${DEPENDENCIES[@]}"
  do
    is_installed $dependency || { echo >&2 "ERROR: $dependency is not installed!"; exit 1; }
  done
}

function is_required_go_version {
  go version | grep $REQUIRED_GO_VERSION >/dev/null 2>&1
}

check_dependencies
is_required_go_version || { echo >&2 "ERROR: You have the wrong version of Go installed! Please install version $REQUIRED_GO_VERSION"; exit 1; }

#Error catching is now complete. Below is package installation
npm install -g grunt-cli bower
npm install
bower install
go get
grunt build

echo "Installation has successfully completed! Use 'go run' to run the website"
