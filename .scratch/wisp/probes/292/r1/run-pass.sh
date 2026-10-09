#!/bin/sh
# 票 292 腿 292-r1 AC#3 单发整包 runner
# 用法: sh .scratch/wisp/probes/292/r1/run-pass.sh <tag>
# 只打印三数 + tasklist 读数 + FAIL 名册，整包日志留在 logs/<tag>.txt
set -u
tag="$1"
cd "$(git rev-parse --show-toplevel)" || exit 1
mkdir -p .scratch/wisp/probes/292/r1/logs

tl() {
  # 现跑 tasklist：某镜像名的枚数（无进程时 tasklist 仍回表头，grep -c 计数为 0）
  n=$(tasklist //FI "IMAGENAME eq $1" //NH 2>/dev/null | grep -ci "^$2" || true)
  printf '%s=%s' "$1" "$n"
}

printf 'PASS %s start ' "$tag"
date '+%Y-%m-%d %H:%M:%S %z'
printf '  tasklist-before: %s %s\n' "$(tl wisp.exe wisp.exe)" "$(tl balldebug.exe balldebug.exe)"

PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test ./cmd/wisp/ -count=1 -v \
  > ".scratch/wisp/probes/292/r1/logs/$tag.txt" 2>&1
grc=$?

printf '  tasklist-after: %s %s\n' "$(tl wisp.exe wisp.exe)" "$(tl balldebug.exe balldebug.exe)"
log=".scratch/wisp/probes/292/r1/logs/$tag.txt"
total=$(grep -c '^=== RUN' "$log" || true)
pass=$(grep -c '^--- PASS' "$log" || true)
fail=$(grep -c '^--- FAIL' "$log" || true)
skip=$(grep -c '^--- SKIP' "$log" || true)
printf '  go-test-rc=%s RUN=%s PASS=%s FAIL=%s SKIP=%s bytes=%s\n' "$grc" "$total" "$pass" "$fail" "$skip" "$(wc -c < "$log")"
printf '  FAIL-names:\n'
grep '^--- FAIL' "$log" | sed 's/^/    /' || true
printf '  exit-status-hex (环境红族): '
grep -o 'exit status 0xc000013[5a]' "$log" | sort | uniq -c | tr '\n' ' ' || true
printf '\n'
