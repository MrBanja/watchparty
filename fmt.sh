#!/usr/bin/env bash
set -e -x

docker run \
    --rm -i \
    -u "$(id -u):$(id -g)" \
    -e USER="$USER" \
    -v /etc/passwd:/etc/passwd:ro \
    -v "$(pwd):$(pwd)":rw \
    -v "$HOME/.cache:$HOME/.cache":rw \
    --workdir "$(pwd)" \
    docker.io/themakers/protoc:v1.0.10 \
    /bin/bash <<- EOF

find . -type f -name "*.proto" | xargs clang-format -i

EOF
