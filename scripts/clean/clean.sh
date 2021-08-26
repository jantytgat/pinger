#!/usr/bin/env bash

SOURCE="${BASH_SOURCE[0]}"
echo $SOURCE
while [ -h "$SOURCE" ] ; do SOURCE="$(readlink "$SOURCE")"; done
DIR="$( cd -P "$( dirname "$SOURCE" )/../.." && pwd )"

cd "$DIR"

rm -rf bin
rm -rf pkg
sleep 5