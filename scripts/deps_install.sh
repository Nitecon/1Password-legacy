#!/usr/bin/env bash
SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ] ; do SOURCE="$(readlink "$SOURCE")"; done
appDir="$( cd -P "$( dirname "$SOURCE" )/../" && pwd )"
cd $appDir
npm install