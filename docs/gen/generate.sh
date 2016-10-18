#!/bin/bash
#
# Generate man pages for docker/docker
#

set -eu

mkdir -p ./man/man1  docs/dist

# Generate man pages from cobra commands
go build -o /tmp/gen-manpages ./docs/gen/generate.go
/tmp/gen-manpages --root . --man ./man/man1 --cli docs/dist

# Generate legacy pages from markdown
./man/md2man-all.sh -q
