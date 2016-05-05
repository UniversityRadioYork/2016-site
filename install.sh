#!/bin/bash

DEPENDENCIES=(npm go)

function is_installed() {
  program_location="$(command -v $1 >/dev/null 2>&1)"
}

function check_dependencies() {
  for dependency in "${DEPENDENCIES[@]}"
  do
    is_installed $dependency || { echo >&2 "ERROR: $dependency is not installed!"; exit 1;}
  done
}

check_dependencies
