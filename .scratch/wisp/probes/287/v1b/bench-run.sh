#!/usr/bin/env bash
# 287-v1b 突变台件（第三版：OUT 一律绝对路径、测试一律在子 shell 里 cd 到仓外树，⛔ 污染本 shell 的 cwd）
# 用法：sh bench-run.sh <tree> <out> [FULL=1] <shot> [<shot> ...]
# 四件套＝种前哈希／落种复量（anchor_hits＋hunks＋改行数）／红句逐字落盘／还原哈希等值。
set -u
TREE="${1:?tree}"; OUT="${2:?out}"
shift 2
mkdir -p "$OUT/logs" "$OUT/anchors"
OUT="$(cd "$OUT" && pwd)"          # ★绝对化：后面所有写入都以此为锚，与 cwd 无关
TREE="$(cd "$TREE" && pwd)"
FULL="${FULL:-0}"; [ "$FULL" = "0" ] && [ "${1:-}" = "FULL=1" ] && { FULL=1; shift; }
PY="${PYTHON:-python}"
A="$OUT/anchors"
ROSTER="$OUT/mutation-roster.txt"
: >"$ROSTER"

S197="internal/tools/subagent_197.go"
STASK="internal/tools/task.go"
SLOOP="internal/agent/loop.go"

bak_for() { [ -f "$TREE/$1.bak" ] || cp "$TREE/$1" "$TREE/$1.bak"; }
base_sha() { sha256sum "$TREE/$1.bak" | cut -d' ' -f1; }
seed_pristine() { cp "$TREE/$1.bak" "$TREE/$1"; }

run_shot() {
  local name="$1" file="$2" kind="$3"
  local log="$OUT/logs/v1b-$name-test.log" fulllog="$OUT/logs/v1b-$name-test-full.log"
  local meta="$OUT/logs/v1b-$name.meta.txt"
  bak_for "$file"; seed_pristine "$file"
  {
    echo "SHOT=$name FILE=$file KIND=$kind TREE=$TREE FULL=$FULL"
    echo "PRE-SHA $(base_sha "$file")"
    if [ "$kind" = "v5" ]; then
      "$PY" "$OUT/mutate_v5.py" "$TREE/$file"
    else
      "$PY" "$OUT/mutate.py" "$TREE/$file" "$A/$kind.old" "$A/$kind.new"
    fi
    local mrc=$?
    echo "mutate_rc=$mrc"
    if [ "$mrc" -ne 0 ]; then echo "REFUSED-TO-SEED"; exit 7; fi
    echo "POST-SHA $(sha256sum "$TREE/$file" | cut -d' ' -f1)"
    echo "hunks=$(diff -u "$TREE/$file.bak" "$TREE/$file" | grep -c '^@@')"
    echo "changed_lines=$(diff -u "$TREE/$file.bak" "$TREE/$file" | grep -c -E '^[+-][^+-]')"
    echo "--- diff-body ---"; diff -u "$TREE/$file.bak" "$TREE/$file" | tail -30
    ( cd "$TREE" && GOFLAGS= go test ./internal/tools/ -run '^Test287' -count=1 -v ) >"$log" 2>&1
    echo "test_rc=$?"
    echo "--- case-colors ---"; grep -E '^(--- PASS|--- FAIL): Test287' "$log"
    echo "--- red-points ---"; grep -oE 'subagent_depth_fixture_287_test\.go:[0-9]+' "$log" | sort -u | tr '\n' ' '; echo
    echo "fault-tokens=$(grep -ciE 'build failed|no tests to run|0xc0000135|test timed out' "$log")"
    echo "--- red-sentences（逐字，含文件:行）---"
    grep -nE '^\s+(subagent_depth_fixture_287_test\.go:[0-9]+:|--- FAIL)' "$log" | head -40
    if [ "$FULL" = "1" ]; then
      ( cd "$TREE" && GOFLAGS= go test ./internal/tools/ -count=1 -v ) >"$fulllog" 2>&1
      echo "full_rc=$?"
      echo "full_TOP_PASS=$(grep -c '^--- PASS' "$fulllog") full_TOP_FAIL=$(grep -c '^--- FAIL' "$fulllog") full_TOP_SKIP=$(grep -c '^--- SKIP' "$fulllog") full_RUN=$(grep -c '^=== RUN' "$fulllog")"
      echo "--- full-red-roster ---"; grep -E '^--- FAIL' "$fulllog"
    fi
  } >"$meta" 2>&1
  seed_pristine "$file"
  local now; now="$(sha256sum "$TREE/$file" | cut -d' ' -f1)"
  if [ "$now" = "$(base_sha "$file")" ]; then echo "RESTORED-EQUAL-OK $now" >>"$meta"; else echo "RESTORE-MISMATCH $now" >>"$meta"; fi
  {
    echo "===== $name ====="
    grep -E '^(SHOT|PRE-SHA|POST-SHA|mutate_rc|hunks|changed_lines|test_rc|fault-tokens|full_rc|full_TOP_PASS|RESTORED-EQUAL-OK|RESTORE-MISMATCH)' "$meta"
    sed -n '/--- case-colors ---/,/--- red-sentences/p' "$meta"
    sed -n '/--- full-red-roster ---/,$p' "$meta"
  } >>"$ROSTER"
  echo "shot $name done meta=$meta"
}

for s in "$@"; do
  case "$s" in
    V1) run_shot V1 "$S197" v1 ;;
    V2) run_shot V2 "$S197" v2 ;;
    V3) run_shot V3 "$S197" v3 ;;
    V4) run_shot V4 "$S197" v4 ;;
    V5) run_shot V5 "$S197" v5 ;;
    V6) run_shot V6 "$SLOOP" v6 ;;
    V7) run_shot V7 "$STASK" v7 ;;
    V8) run_shot V8 "$SLOOP" v8 ;;
    *) echo "unknown shot $s" ;;
  esac
done
echo "ROSTER=$ROSTER"
