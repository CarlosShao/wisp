#!/bin/bash
# 304-v1b helper spawn point shots (forward T1, reverse T2), single-process watcher + guard-attributed stop.
SD=$(cd "$(dirname "$0")" && pwd)
LABS="$SD/logs"
T1=$(cat "$LABS/t1.path")
T2=$(cat "$LABS/t2.path")
WT1=$(cygpath -w "$T1")
WT2=$(cygpath -w "$T2")
WK=$(cygpath -w "$SD/wk.ps1")
ROW=$(cygpath -w "$SD/ro.ps1")
RUN='-test.run=^TestFailoverAfterPrimaryExhaustion$'
echo "== fwd helper (T1, guard present) =="
(cd "$T1" && "$T1/ac2t.exe" "$RUN" -test.v -test.timeout 600s > "$LABS/ac2-fwd-helper3-c-body.txt" 2>&1) &
TP=$!
powershell -NoProfile -ExecutionPolicy Bypass -File "$WK" -Tree "$WT1" > "$LABS/ac2-fwd-helper3-c-watchkill.txt" 2>&1
echo "wk_rc=$?" >> "$LABS/ac2-fwd-helper3-c-watchkill.txt"
wait $TP
echo "go_test_rc=$?" >> "$LABS/ac2-fwd-helper3-c-watchkill.txt"
sleep 1
powershell -NoProfile -ExecutionPolicy Bypass -File "$ROW" -Tree "$WT1" >> "$LABS/ac2-fwd-helper3-c-watchkill.txt" 2>&1
echo "== rev helper (T2, guard removed) =="
(cd "$T2" && GOFLAGS= go test -c -o "$T2/ac2t.exe" ./internal/llm > "$LABS/ac2-rev-helper-c-compile.txt" 2>&1)
echo "compile_rc=$?" >> "$LABS/ac2-rev-helper-c-compile.txt"
(cd "$T2" && "$T2/ac2t.exe" "$RUN" -test.v -test.timeout 600s > "$LABS/ac2-rev-helper-c-body.txt" 2>&1) &
TP=$!
powershell -NoProfile -ExecutionPolicy Bypass -File "$WK" -Tree "$WT2" > "$LABS/ac2-rev-helper-c-watchkill.txt" 2>&1
echo "wk_rc=$?" >> "$LABS/ac2-rev-helper-c-watchkill.txt"
wait $TP
echo "go_test_rc=$?" >> "$LABS/ac2-rev-helper-c-watchkill.txt"
sleep 1
powershell -NoProfile -ExecutionPolicy Bypass -File "$ROW" -Tree "$WT2" -Kill >> "$LABS/ac2-rev-helper-c-watchkill.txt" 2>&1
echo "== done =="
