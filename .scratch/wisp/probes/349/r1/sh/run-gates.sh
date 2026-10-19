#!/bin/sh
# 349-r1 门禁四数（★只在共享工作树跑；跑前 unset GOTMPDIR TMPDIR TEMP TMP＝票 311 shape B GUARD
# 要求测试临时根与 %LOCALAPPDATA% 同卷，带 env 时 rc=2 射程拒绝⛔ 读成红）。
# 逐枚落 rc=；红句含 walk／build failed／no test files／setup failed／0xc0000135 ⇒ 判台面故障。
set -u
cd '/d/work/workspace/projects plans/Wisp' || exit 90
unset GOTMPDIR TMPDIR TEMP TMP
L=.scratch/wisp/probes/349/r1/logs/gates
mkdir -p "$L"
echo "### gate 1: gofumpt -l（三枚被动文件；本机工作树那批 CRLF 存量件具名登记⛔ 顺手重排）"
/d/work/base/gopath/bin/gofumpt.exe -l internal/risk/blacklist.go internal/risk/provenance.go internal/risk/syncdirs.go > "$L/g1.txt" 2>&1
echo "rc=$?"
cat "$L/g1.txt"
echo "-- 同一把尺打两枚未被我改形的 blob（起手锚 7b633fa9 ⇄ 交件 HEAD，都是 LF 的 git 对象形）--"
G=/d/tmp/wisp/20261019/349r1/gofumpt
mkdir -p "$G/before" "$G/after"
for f in blacklist provenance syncdirs; do
  git cat-file blob "7b633fa9713b0514806c8b808ed4db91b1485cdc:internal/risk/$f.go" > "$G/before/$f.go"
  git cat-file blob "HEAD:internal/risk/$f.go" > "$G/after/$f.go"
done
/d/work/base/gopath/bin/gofumpt.exe -l "$G/before"/*.go > "$L/g1-before.txt" 2>&1; echo "rc_before=$?"
/d/work/base/gopath/bin/gofumpt.exe -l "$G/after"/*.go > "$L/g1-after.txt"  2>&1; echo "rc_after=$?"
printf 'before_dirty_lines=%s after_dirty_lines=%s\n' "$(wc -l < "$L/g1-before.txt")" "$(wc -l < "$L/g1-after.txt")"
echo "### gate 2: sh scripts/d22scan.sh"
sh scripts/d22scan.sh > "$L/g2.txt" 2>&1
echo "rc=$?"
tail -12 "$L/g2.txt"
echo "### gate 3: sh scripts/check-path-length-budget.sh"
sh scripts/check-path-length-budget.sh > "$L/g3.txt" 2>&1
echo "rc=$?"
tail -8 "$L/g3.txt"
echo "### gate 4: go vet ./internal/risk/"
/d/work/base/go/bin/go.exe vet ./internal/risk/ > "$L/g4.txt" 2>&1
echo "rc=$?"
cat "$L/g4.txt"
echo "### 台面故障名册（want none）"
grep -nE 'walk|build failed|no test files|setup failed|0xc0000135' "$L/g2.txt" "$L/g3.txt" "$L/g4.txt" || echo "none"
