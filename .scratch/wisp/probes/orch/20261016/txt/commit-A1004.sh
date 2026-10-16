#!/bin/sh
# 编排者收档机械件（A1004）：逐枚显式 pathspec，⛔ add -A／⛔ 目录形／⛔ amend／⛔ push。
set -eu
cd "$(git rev-parse --show-toplevel)"
I=.scratch/wisp/issues
P=.scratch/wisp/probes/orch/20261016

F317OLD=$I/317-localslash-fold-in-gotmpdir-root-makes-windows-junction-tests-falsely-red.md
F317NEW=$I/317-localslash-fold-in-gotmpdir-root-makes-windows-junction-tests-falsely-red-done.md
F325OLD=$I/325-runtests-sh-gotmp-tail-keeps-forward-slashes-so-d22scan-runs-on-a-mixed-shape-root.md
F325NEW=$I/325-runtests-sh-gotmp-tail-keeps-forward-slashes-so-d22scan-runs-on-a-mixed-shape-root-done.md
F287OLD=$I/287-depth-guard-has-no-fixture-so-the-subagent-spawn-hop-is-unreachable-and-untested.md
F287NEW=$I/287-depth-guard-has-no-fixture-so-the-subagent-spawn-hop-is-unreachable-and-untested-done.md
F165=$I/165-plan-first-then-act-a-brake-the-voice-entry-needs.md

git add -- \
  docs/reports/pending-and-issues.md \
  "$F317NEW" "$F325NEW" "$F287NEW" "$F165" \
  "$P/txt/flip-1004.sh" "$P/txt/commit-A1004.sh" \
  "$P/txt/bullet-317-done.md" "$P/txt/bullet-325-verdicts.md" "$P/txt/bullet-325-done.md" \
  "$P/txt/bullet-165-flip.md" "$P/txt/bullet-287-done.md" \
  "$P/txt/ledger-A1004.md" "$P/txt/msg-A1004.txt" \
  "$P/logs/orch-325v1-ac3-spot4.txt"

echo "=== staged 名册（应只含上面这些）==="
git diff --cached --name-status | cat

git commit -q -F "$P/txt/msg-A1004.txt" -- \
  docs/reports/pending-and-issues.md \
  "$F317OLD" "$F317NEW" "$F325OLD" "$F325NEW" "$F287OLD" "$F287NEW" "$F165" \
  "$P/txt/flip-1004.sh" "$P/txt/commit-A1004.sh" \
  "$P/txt/bullet-317-done.md" "$P/txt/bullet-325-verdicts.md" "$P/txt/bullet-325-done.md" \
  "$P/txt/bullet-165-flip.md" "$P/txt/bullet-287-done.md" \
  "$P/txt/ledger-A1004.md" "$P/txt/msg-A1004.txt" \
  "$P/logs/orch-325v1-ac3-spot4.txt"

echo "=== HEAD 与路径清单 ==="
git log -1 --stat --format='%H%n%s' | cat
