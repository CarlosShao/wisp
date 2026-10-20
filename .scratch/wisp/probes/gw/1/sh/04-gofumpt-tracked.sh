#!/bin/sh
# gw-1 味① 第 4⇄5 把的**权威形**复跑 r3。
# 尺的出处＝本仓 CI 逐字（.github/workflows/ci.yml:172-181）＝
#   OUT="$(gofumpt -l . tools/d22scan tools/mockllm)"; [ -n "$OUT" ] ⇒ exit 1
# 另枚＝ci.yml:189-230 那步把**分母钉成 tracked 集**（票 161 AC#7 形 A＝`attrib.sh --tracked-only`），
#   理由逐字写在 ci.yml:196-200＝`.scratch/**` 按定义带着故意写坏的样本⇒全树尺⛔ 可满足。
# ⇒ 本腿并排打三把：① CI 逐字形（全树）② 同输出对 tracked 集过滤＝判据那一把 ③ gofmt 独立第二把。
# ★⛔ 洗已知脏件、⛔ -w；存量逐枚具名登记（含 CR 尺：对**工作树字节**量＝gofumpt 看到的正是工作树字节，
#   对**入库 blob**另量一枚＝说明本地 CRLF 与 CI checkout 的 LF 差＝那一形⛔ 是真红）。
set -u
cd '/d/work/workspace/projects plans/Wisp' || exit 9
R=/d/tmp/wisp/20261020/gw1
LOGS="$R/logs"
mkdir -p "$LOGS"
OUT="$LOGS/gates-gofumpt-r3.txt"
: > "$OUT"
{
  echo "== gates-r3 START anchor =="
  git rev-parse HEAD | sed 's|^|HEAD_FULL=|'
  echo "WHEN=$(date '+%F %T %z')"
  unset GOTMPDIR TMPDIR TEMP TMP
  echo "env GOTMPDIR=[${GOTMPDIR-<unset>}] TMPDIR=[${TMPDIR-<unset>}] TEMP=[${TEMP-<unset>}] TMP=[${TMP-<unset>}]"
  echo
  echo "== ci.yml:176 逐字形＝gofumpt -l . tools/d22scan tools/mockllm =="
  /d/work/base/gopath/bin/gofumpt.exe -l . tools/d22scan tools/mockllm > "$LOGS/g4-ci-shape.txt" 2>&1
  echo "G4_CI_RC=$?"
  echo "G4_CI_LINES=$(wc -l < "$LOGS/g4-ci-shape.txt")"
  echo
  echo "== 同一把尺对 tracked 集过滤（＝判据那一把：分母＝tracked）=="
  : > "$LOGS/g4-tracked.txt"
  : > "$LOGS/g4-untracked.txt"
  : > "$LOGS/g4-notgo.txt"
  while IFS= read -r line; do
    [ -n "$line" ] || continue
    p=$(printf '%s\n' "$line" | cut -d: -f1 | tr '\\' '/')
    case "$p" in
      *.go) ;;
      *) echo "$line" >> "$LOGS/g4-notgo.txt"; continue ;;
    esac
    if git ls-files --error-unmatch "$p" >/dev/null 2>&1; then
      echo "$p" >> "$LOGS/g4-tracked.txt"
    else
      echo "$p" >> "$LOGS/g4-untracked.txt"
    fi
  done < "$LOGS/g4-ci-shape.txt"
  echo "TRACKED_DIRTY=$(grep -c . "$LOGS/g4-tracked.txt" || true)"
  grep . "$LOGS/g4-tracked.txt" | sed 's|^|  TRACKED |' || true
  echo "UNTRACKED_DIRTY=$(grep -c . "$LOGS/g4-untracked.txt" || true)"
  grep . "$LOGS/g4-untracked.txt" | sed 's|^|  UNT |' || true
  echo "NOTGO_LINES=$(grep -c . "$LOGS/g4-notgo.txt" || true)"
  grep . "$LOGS/g4-notgo.txt" | sed 's|^|  PARSE |' || true
  echo
  echo "== 逐枚 CR 尺（工作树字节＝工具看到的形 ⇄ 入库 blob＝CI checkout 的形）=="
  while IFS= read -r p; do
    [ -n "$p" ] || continue
    WT=$(tr -cd '\r' < "$p" | wc -c)
    BL=$(git cat-file blob "HEAD:$p" 2>/dev/null | tr -cd '\r' | wc -c)
    echo "  CRFILE $p worktree_cr=$WT blob_cr=$BL"
  done < "$LOGS/g4-tracked.txt"
  echo
  echo "== 第二把独立尺＝gofmt -l（同参数形）=="
  gofmt -l . tools/d22scan tools/mockllm > "$LOGS/g4b-ci-shape.txt" 2>&1
  echo "G4B_RC=$?"
  echo "G4B_LINES=$(wc -l < "$LOGS/g4b-ci-shape.txt")"
  echo "TRACKED_GOFMT_DIRTY=$(gofmt -l . tools/d22scan tools/mockllm 2>/dev/null | tr '\\' '/' | grep '\.go$' | while IFS= read -r p; do git ls-files --error-unmatch "$p" >/dev/null 2>&1 && echo x; done | grep -c x || true)"
  echo
  echo "== 正控（证明这把尺⛔ 瞎）＝自造一枚脏 .go 落仓外，量它被列出来 =="
  mkdir -p "$R/posctrl"
  printf 'package p\n\nfunc  F(){ }\n' > "$R/posctrl/dirty.go"
  /d/work/base/gopath/bin/gofumpt.exe -l "$R/posctrl/dirty.go" > "$LOGS/g4-posctrl.txt" 2>&1
  echo "POSCTRL_RC=$?"
  echo "POSCTRL_LISTED=$(grep -c . "$LOGS/g4-posctrl.txt" || true)"
  sed 's|^|  POS |' "$LOGS/g4-posctrl.txt"
  echo
  echo "== END anchor =="
  git rev-parse HEAD | sed 's|^|HEAD_FULL_END=|'
  echo "WHEN_END=$(date '+%F %T %z')"
} >> "$OUT" 2>&1
cat "$OUT"
