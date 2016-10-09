#!/bin/bash

DEPENDENCIES=(npm go)
declare -A NPM_DEPENDENCIES=([bower]=bower [grunt]=grunt-cli)
REQUIRED_GO_VERSION=1.7

function is_installed {
  command -v $1 >/dev/null 2>&1
}

function npm_modules_is_installed {
  npm ls --global=true --parseable=true | grep $1 > /dev/null 2>&1
}

function check_dependencies {
  for dependency in "${DEPENDENCIES[@]}"
  do
    is_installed $dependency || { echo >&2 "ERROR: $dependency is not installed!"; exit 1; }
  done
}

function is_required_go_version {
  go version | grep $REQUIRED_GO_VERSION >/dev/null 2>&1
}

function install_npm_dependencies {
  for dependency in "${!NPM_DEPENDENCIES[@]}"
  do
    npm_modules_is_installed $dependency || { echo >&2 "$dependency not installed. Now installing...";
    npm install -g "${NPM_DEPENDENCIES[$dependency]}";
  }
  done
}

check_dependencies
is_required_go_version || { echo >&2 "ERROR: You have the wrong version of Go installed! Please install version $REQUIRED_GO_VERSION"; exit 1; }
install_npm_dependencies

#Error catching is now complete. Below is package installation
npm install
bower install
go get
grunt build

echo "Installation has successfully completed! Use 'go run main.go' to run the website"
