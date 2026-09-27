#!/bin/sh
# 163/a1 gotest 台件包装：logdir 取本脚本自身目录下的 logs/，不继承 CWD。
set -u
PROBE_DIR="$(cd "$(dirname "$0")" && pwd)"
REPO="$(git rev-parse --show-toplevel)"
cd "$REPO" || exit 9
go run "$PROBE_DIR/main.go" -logdir "$PROBE_DIR/logs" 2>&1 | tee "$PROBE_DIR/logs/gotest.out"
