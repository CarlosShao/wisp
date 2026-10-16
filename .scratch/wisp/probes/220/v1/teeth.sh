#!/usr/bin/env bash
# 220-v1 读数与突变台件（全部打在仓外导出树 /tmp/wisp220v1/{bench,pre,full}）
set -u
REPO="/d/work/workspace/projects plans/Wisp"
PROBE="$REPO/.scratch/wisp/probes/220/v1"
LOGS="$PROBE/logs"
ROOT=/tmp/wisp220v1
export WISP220ROOT="$(cygpath -w $ROOT/bench)"
PY=python
mkdir -p "$LOGS"

env_for() {
  local got="$1/.gotmp"; mkdir -p "$got"
  local ng; ng=$(cygpath -w "$got")
  echo "$ng"
}

go_test() { # tree pkg extra-args... > log
  local tree="$1"; shift
  local log="$1"; shift
  local ng; ng=$(env_for "$ROOT/$tree")
  date '+%F %T %z' > "$log"
  echo "cmd=go test $* (tree=$tree)" >> "$log"
  ( cd "$ROOT/$tree" && PATH="$PWD/third_party/sherpa-onnx:$PATH" GOTMPDIR="$ng" TMP="$ng" TEMP="$ng" \
      go test "$@" >> "$log" 2>&1; echo "rc=$?" >> "$log" )
  echo "--- counts: RUN=$(grep -c '^=== RUN' "$log") TOPPASS=$(grep -c '^--- PASS' "$log") FAIL=$(grep -c '^--- FAIL' "$log") SKIP=$(grep -c '^--- SKIP' "$log")" >> "$log"
  echo "--- fail names ---" >> "$log"
  grep '^--- FAIL' "$log" >> "$log" 2>/dev/null || echo "(no FAIL lines)" >> "$log"
  date '+%F %T %z' >> "$log"
}

echo "### 1 生产者计数尺（同一把尺·两发：改前锚 fd01e7962 ⇄ HEAD 导出树）###"
COUNT="$LOGS/01-count.txt"; : > "$COUNT"
date '+%F %T %z' >> "$COUNT"
echo "尺=grep -rn L1Windows --include=*.go cmd  (射程=该导出树 cmd/ 全部 .go，含测试；量的)" >> "$COUNT"
for t in pre bench; do
  echo "--- tree=$t ref=$(cd "$ROOT/$t" && echo -n "") $( [ $t = pre ] && echo fd01e7962 || echo HEAD ) ---" >> "$COUNT"
  grep -rn 'L1Windows' --include='*.go' "$ROOT/$t/cmd" >> "$COUNT" 2>&1 || echo "(零命中 rc=1)" >> "$COUNT"
  echo "  全部命中=$(grep -rn 'L1Windows' --include='*.go' "$ROOT/$t/cmd" 2>/dev/null | wc -l)" >> "$COUNT"
  echo "  非测试命中=$(grep -rn 'L1Windows' --include='*.go' "$ROOT/$t/cmd" 2>/dev/null | grep -vc '_test\.go')" >> "$COUNT"
  echo "  装配点 run.go 'L1Windows:'=$(grep -c 'L1Windows:' "$ROOT/$t/cmd/wisp/run.go" 2>/dev/null || echo 0)" >> "$COUNT"
  echo "  reader 本体 'rt.gate.LiveL1Windows()'=$(grep -c 'rt\.gate\.LiveL1Windows()' "$ROOT/$t/cmd/wisp/panel_pump.go" 2>/dev/null || echo 0)" >> "$COUNT"
done
echo "正控：同一把尺在 HEAD 树上必须非零 ⇒见上；零命中那一发下面有已知存在的正控行号" >> "$COUNT"

echo "### 2 台面基线（internal/panel 整包 / approval 220 用例 / cmd/wisp 真腿用例）###"
go_test bench "$LOGS/02-panel-base.txt" ./internal/panel/ -count=1
APPROW=$(grep -hoE '^func Test[A-Za-z0-9_]+' "$ROOT/bench/internal/agent/approval/ticket220_l1_window_read_test.go" | sed -E 's/^func //' | paste -sd'|')
echo "approval 220 用例名册=$APPROW"
go_test bench "$LOGS/03-approval-base.txt" ./internal/agent/approval/ -count=1 -timeout 20m -run "^($APPROW)$"
go_test bench "$LOGS/04-cmdwisp-base.txt" ./cmd/wisp/ -count=1 -timeout 30m -v -run '^TestRunPacketLightsTheRosterRowOfItsLiveL1Window$'

echo "### 3 突变台（每发：pre 检查⇒apply⇒跑⇒restore⇒sha 对回）###"
for m in MU-A MU-B MU-C; do
  echo "=== $m ==="
  $PY "$PROBE/mut.py" pre "$m" >> "$LOGS/10-mut-check.txt" 2>&1
  $PY "$PROBE/mut.py" apply "$m" >> "$LOGS/10-mut-check.txt" 2>&1
  go_test bench "$LOGS/20-panel-$m.txt" ./internal/panel/ -count=1
  $PY "$PROBE/mut.py" restore "$m" >> "$LOGS/10-mut-check.txt" 2>&1
done
for m in MU-D MU-E; do
  echo "=== $m ==="
  $PY "$PROBE/mut.py" pre "$m" >> "$LOGS/10-mut-check.txt" 2>&1
  $PY "$PROBE/mut.py" apply "$m" >> "$LOGS/10-mut-check.txt" 2>&1
  go_test bench "$LOGS/21-cmdwisp-$m.txt" ./cmd/wisp/ -count=1 -timeout 30m -v -run '^TestRunPacketLightsTheRosterRowOfItsLiveL1Window$'
  $PY "$PROBE/mut.py" restore "$m" >> "$LOGS/10-mut-check.txt" 2>&1
done
echo "### 4 还原后台面自证（bench 树三枚文件 sha ⇄ 仓内 HEAD blob）###"
Z="$LOGS/30-restore-verify.txt"; : > "$Z"
date '+%F %T %z' >> "$Z"
for f in internal/panel/subagent_roster_197.go cmd/wisp/panel_pump.go cmd/wisp/run.go; do
  a=$(sha256sum "$ROOT/bench/$f" | cut -d' ' -f1)
  b=$( (cd "$REPO" && git cat-file -p "HEAD:$f") | sha256sum | cut -d' ' -f1)
  echo "$f tree=$a head=$b match=$([ "$a" = "$b" ] && echo YES || echo NO)" >> "$Z"
done
echo "ALLDONE $(date '+%F %T %z')"
