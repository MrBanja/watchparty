#!/bin/sh
set -e

docker run \
    --rm -i \
    -u "$(id -u):$(id -g)" \
    -e USER="$USER" \
    -v /etc/passwd:/etc/passwd:ro \
    -v "$PWD:$PWD":rw \
    -v "$HOME/.cache:$HOME/.cache":rw \
    --workdir "$PWD" \
    docker.io/themakers/protoc:v1.0.11 \
    /bin/bash <<-"EOF"


target=protocol/gen-go
tmp=tmp/gen

function find_generated {
    dir=$1
    shift 1
    find $dir -type f -regextype posix-egrep -regex '.+(\.pb\.go|\.connect\.go)$' $@
}

set -e -x

rm -rf $tmp
mkdir -p $tmp

protoc \
    --proto_path=protocol \
    \
    --go_opt=paths=source_relative \
    --connect-go_opt=paths=source_relative,require_unimplemented_servers=false \
    \
    --go_out=$tmp \
    --connect-go_out=$tmp \
    \
    $(find ./protocol -type f -name '*.proto'  |  sed -e 's/^\.\/protocol\///')

find_generated $target -exec rm {} +

cp -rf tmp/gen/* $target/
rm -rf tmp


EOF