#!/bin/sh
# gw-1 量色腿 味①＝票 321 AC#2 五把门对今天 HEAD 重跑。
# ★只在共享工作树、跑前 unset GOTMPDIR TMPDIR TEMP TMP（带 env＝票 311 AC#1 shape B 那道 GUARD rc=2
#   射程拒绝＝⛔ 读成红，台账 A1027⑤／A1072 并入的 109-r1 那味：台根⛔ 落 D:\tmp\wisp）。
# ★退码＝重定向文件后 echo RC=$? 紧跟裸命令（⛔ 管进 tail＝取到的是 tail 的退码；⛔ PowerShell 内联 $?）。
# ★gofumpt⇄gofmt 并排两把、⛔ 洗已知脏件、⛔ -w：存量逐枚具名登记。
set -u
cd '/d/work/workspace/projects plans/Wisp' || exit 9
R=/d/tmp/wisp/20261020/gw1
LOGS="$R/logs"
mkdir -p "$LOGS"
OUT="$LOGS/gates.txt"
: > "$OUT"
{
  echo "== gw-1 gates START anchor =="
  git rev-parse HEAD | sed 's/^/HEAD_FULL='
  date '+%F %T %z' | sed 's/^/WHEN='
  echo
  echo "== env after unset =="
  unset GOTMPDIR TMPDIR TEMP TMP
  echo "GOTMPDIR=[${GOTMPDIR-<unset>}] TMPDIR=[${TMPDIR-<unset>}] TEMP=[${TEMP-<unset>}] TMP=[${TMP-<unset>}]"
  echo "GOFLAGS=[${GOFLAGS-<unset>}]"
  echo
  echo "== gate 1: GOFLAGS= go build ./...（共享树）=="
  GOFLAGS= go build ./... > "$LOGS/gate1-build.txt" 2>&1
  echo "GATE1_BUILD_RC=$?"
  echo "GATE1_BYTES=$(wc -c < "$LOGS/gate1-build.txt")"
  sed 's/^/  B1| /' "$LOGS/gate1-build.txt"
  echo
  echo "== gate 2: sh scripts/d22scan.sh =="
  sh scripts/d22scan.sh > "$LOGS/gate2-d22scan.txt" 2>&1
  echo "GATE2_D22_RC=$?"
  echo "GATE2_BYTES=$(wc -c < "$LOGS/gate2-d22scan.txt")"
  tail -5 "$LOGS/gate2-d22scan.txt" | sed 's/^/  B2| /'
  echo
  echo "== gate 3: sh scripts/check-path-length-budget.sh（要 .git）=="
  sh scripts/check-path-length-budget.sh > "$LOGS/gate3-length.txt" 2>&1
  echo "GATE3_LEN_RC=$?"
  echo "GATE3_BYTES=$(wc -c < "$LOGS/gate3-length.txt")"
  tail -6 "$LOGS/gate3-length.txt" | sed 's/^/  B3| /'
  echo
  echo "== gate 4a: /d/work/base/gopath/bin/gofumpt.exe -l ./...（全仓尺，⛔ -w）=="
  /d/work/base/gopath/bin/gofumpt.exe -l ./... > "$LOGS/gate4a-gofumpt.txt" 2>&1
  echo "GATE4A_GOFUMPT_RC=$?"
  echo "GATE4A_LISTED=$(grep -c . "$LOGS/gate4a-gofumpt.txt" || true)"
  sed 's/^/  FMT| /' "$LOGS/gate4a-gofumpt.txt"
  echo
  echo "== gate 4b: gofmt -l ./...（并排的第二把独立尺）=="
  gofmt -l ./... > "$LOGS/gate4b-gofmt.txt" 2>&1
  echo "GATE4B_GOFMT_RC=$?"
  echo "GATE4B_LISTED=$(grep -c . "$LOGS/gate4b-gofmt.txt" || true)"
  sed 's/^/  FMT| /' "$LOGS/gate4b-gofmt.txt"
  echo
  echo "== gate 4c: 两把尺逐名作差（comm，⛔ 拿枚数当名册）=="
  sort "$LOGS/gate4a-gofumpt.txt" > "$LOGS/gate4a.sorted"
  sort "$LOGS/gate4b-gofmt.txt" > "$LOGS/gate4b.sorted"
  echo "-- only-in-gofumpt --"
  comm -23 "$LOGS/gate4a.sorted" "$LOGS/gate4b.sorted" | sed 's/^/  GONLY /'
  echo "-- only-in-gofmt --"
  comm -13 "$LOGS/gate4a.sorted" "$LOGS/gate4b.sorted" | sed 's/^/  MONLY /'
  echo "-- both --"
  comm -12 "$LOGS/gate4a.sorted" "$LOGS/gate4b.sorted" | wc -l | sed 's/^/  BOTH_COUNT=/'
  echo
  echo "== gate 5: go vet -tags winlive ./cmd/wisp/（共享树）=="
  GOFLAGS= go vet -tags winlive ./cmd/wisp/ > "$LOGS/gate5-vet.txt" 2>&1
  echo "GATE5_VET_RC=$?"
  echo "GATE5_BYTES=$(wc -c < "$LOGS/gate5-vet.txt")"
  sed 's/^/  V5| /' "$LOGS/gate5-vet.txt"
  echo
  echo "== gate 5 preflight: tracked dirty（⛔ .go 脏＝半保存态假红前兆）=="
  git status --porcelain -uno | sed 's/^/  DIRTY /'
  echo "DIRTY_LINES=$(git status --porcelain -uno | wc -l)"
  echo "DIRTY_GO=$(git status --porcelain -uno -- '*.go' | wc -l)"
  echo
  echo "== END anchor =="
  git rev-parse HEAD | sed 's/^/HEAD_FULL_END='
  date '+%F %T %z' | sed 's/^/WHEN_END='
} >> "$OUT" 2>&1
echo "WROTE $OUT"
wc -l "$OUT"
cat "$OUT"
