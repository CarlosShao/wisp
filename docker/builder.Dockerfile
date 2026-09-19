# Wisp Linux cross-build convenience image (SPEC-11 §3.1)
# NOTE: cross-linking sherpa Windows prebuilt libs with mingw is validated in ticket 02 (spike);
# the PRIMARY supported build path is native Windows (scripts/build.ps1, docs/BUILD.md).
# Deviation from SPEC-11 §3.1 (golang:1.24 -> 1.27, 2026-09-19): go.mod pins go 1.27;
# a 1.24 base image would download a newer toolchain at build time, breaking pinning.
FROM golang:1.27-bookworm
RUN apt-get update && apt-get install -y --no-install-recommends \
    gcc-mingw-w64-x86-64 g++-mingw-w64-x86-64 git ca-certificates && rm -rf /var/lib/apt/lists/*
ENV CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc GOOS=windows GOARCH=amd64
