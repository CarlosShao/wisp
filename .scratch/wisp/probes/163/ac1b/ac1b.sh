#!/bin/sh
# 163/ac1b 台件包装：logdir 取本脚本自身目录下的 logs/，不继承 CWD。
# 派单环境令：PATH 前置主树 build 与 third_party/sherpa-onnx；清空 GOTMPDIR/TMPDIR/TEMP/TMP；
# 绝不设 CGO_ENABLED=0（保持默认 1）。
set -u
PROBE_DIR="$(cd "$(dirname "$0")" && pwd)"
REPO="$(git rev-parse --show-toplevel)"
cd "$REPO" || exit 9
mkdir -p "$PROBE_DIR/logs"
export PATH="$REPO/build:$REPO/third_party/sherpa-onnx:$PATH"
unset GOTMPDIR TMPDIR TEMP TMP
env | grep -E '^CGO_ENABLED=' || echo "# CGO_ENABLED 未设置（默认 1）"
go run "$PROBE_DIR/main.go" -logdir "$PROBE_DIR/logs" 2>&1 | tee "$PROBE_DIR/logs/ac1b.out"
