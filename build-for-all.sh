#!/usr/bin/env bash

TARGETS=" \
android/386 \
android/amd64 \
android/arm \
android/arm64 \
darwin/386 \
darwin/amd64 \
darwin/arm \
darwin/arm64 \
dragonfly/amd64 \
freebsd/386 \
freebsd/amd64 \
freebsd/arm \
linux/386 \
linux/amd64 \
linux/arm \
linux/arm64 \
linux/mips \
linux/mips64 \
linux/mips64le \
linux/mipsle \
linux/ppc64 \
linux/ppc64le \
linux/s390x \
nacl/386 \
nacl/amd64p32 \
nacl/arm \
netbsd/386 \
netbsd/amd64 \
netbsd/arm \
openbsd/386 \
openbsd/amd64 \
openbsd/arm \
plan9/386 \
plan9/amd64 \
plan9/arm \
solaris/amd64 \
windows/386/.exe \
windows/amd64/.exe \
"

N=network-list-parser

for TARGET in ${TARGETS} ; do
    IFS='/' read -a ARR <<< "${TARGET}"
    GOOS=${ARR[0]}
    GOARCH=${ARR[1]}
    EXT=${ARR[2]}
    if [[ -z "${EXT}" ]] ; then
        EXT=".bin"
    fi
    echo $GOOS $GOARCH
    go build -o ${N}-${GOOS}-${GOARCH}${EXT} github.com/unsacrificed/network-list-parser
done
