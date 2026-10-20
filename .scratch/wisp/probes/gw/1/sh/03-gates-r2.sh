#!/bin/sh
# gw-1 味① 复跑 r2＝修两处：① 我首发的 sed 写坏了（`s/^/X=` 少了收尾斜杠⇒锚没落进件＝⛔ 拿"没打印"当没发生）；
#   ② gofumpt/gofmt 的参数形：`./...` 是 Go **包模式**、⛔ 是文件路径，这两枚工具按**路径**走盘⇒
#     首发 rc=2 且 offender 逐字 `open ...\.git`＝**用法坏**、⛔ 读成"仓脏"。r2 一律打在目录 `.` 上。
set -u
cd '/d/work/workspace/projects plans/Wisp' || exit 9
R=/d/tmp/wisp/20261020/gw1
LOGS="$R/logs"
mkdir -p "$LOGS"
OUT="$LOGS/gates-r2.txt"
: > "$OUT"
{
  echo "== gw-1 gates-r2 START anchor =="
  git rev-parse HEAD | sed 's|^|HEAD_FULL=|'
  echo "WHEN=$(date '+%F %T %z')"
  echo
  unset GOTMPDIR TMPDIR TEMP TMP
  echo "GOTMPDIR=[${GOTMPDIR-<unset>}] TMPDIR=[${TMPDIR-<unset>}] TEMP=[${TEMP-<unset>}] TMP=[${TMP-<unset>}] GOFLAGS=[${GOFLAGS-<unset>}]"
  echo
  echo "== preflight: tracked dirty（⛔ .go 脏＝半保存态假红前兆；逐枚列名）=="
  git status --porcelain -uno | sed 's|^|  DIRTY |'
  echo "DIRTY_LINES=$(git status --porcelain -uno | wc -l)"
  echo "DIRTY_GO=$(git status --porcelain -uno -- '*.go' | wc -l)"
  echo
  echo "== gate 4a r2: /d/work/base/gopath/bin/gofumpt.exe -l .（路径形，⛔ -w）=="
  /d/work/base/gopath/bin/gofumpt.exe -l . > "$LOGS/gate4a-r2.txt" 2>&1
  echo "GATE4A_RC=$?"
  echo "GATE4A_LISTED=$(grep -c . "$LOGS/gate4a-r2.txt" || true)"
  sed 's|^|  FMT |' "$LOGS/gate4a-r2.txt"
  echo
  echo "== gate 4b r2: gofmt -l .（并排的第二把独立尺，同一路径形）=="
  gofmt -l . > "$LOGS/gate4b-r2.txt" 2>&1
  echo "GATE4B_RC=$?"
  echo "GATE4B_LISTED=$(grep -c . "$LOGS/gate4b-r2.txt" || true)"
  sed 's|^|  FMT |' "$LOGS/gate4b-r2.txt"
  echo
  echo "== gate 4c r2: 两把尺逐名作差（comm，⛔ 拿枚数顶名册）=="
  grep '\.go$' "$LOGS/gate4a-r2.txt" | sort > "$LOGS/g4a.go.sorted"
  grep '\.go$' "$LOGS/gate4b-r2.txt" | sort > "$LOGS/g4b.go.sorted"
  echo "ONLY_GOFUMPT=$(comm -23 "$LOGS/g4a.go.sorted" "$LOGS/g4b.go.sorted" | grep -c . || true)"
  comm -23 "$LOGS/g4a.go.sorted" "$LOGS/g4b.go.sorted" | sed 's|^|  GONLY |'
  echo "ONLY_GOFMT=$(comm -13 "$LOGS/g4a.go.sorted" "$LOGS/g4b.go.sorted" | grep -c . || true)"
  comm -13 "$LOGS/g4a.go.sorted" "$LOGS/g4b.go.sorted" | sed 's|^|  MONLY |'
  echo "BOTH=$(comm -12 "$LOGS/g4a.go.sorted" "$LOGS/g4b.go.sorted" | grep -c . || true)"
  echo
  echo "== gate 4d: 首发放坏的用法留档复现（⛔ 读成红＝证明 rc=2 是尺形坏）=="
  /d/work/base/gopath/bin/gofumpt.exe -l ./... > "$LOGS/gate4a-wrongshape.txt" 2>&1
  echo "WRONGSHAPE_RC=$?"
  sed 's|^|  WS |' "$LOGS/gate4a-wrongshape.txt"
  echo
  echo "== END anchor =="
  git rev-parse HEAD | sed 's|^|HEAD_FULL_END=|'
  echo "WHEN_END=$(date '+%F %T %z')"
} >> "$OUT" 2>&1
cat "$OUT"
