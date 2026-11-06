# 182 r4 leg — script 03: AC#5 per-pile minefield rulers (panel-initiated approve/write)
set -u
cd "/d/work/workspace/projects plans/Wisp" || exit 9
echo "== exact guikou rule phrase search =="
grep -rn "闭格" .scratch/wisp/issues/README.md AGENTS.md docs/reports/orchestrator-handoff-20261013.md 2>/dev/null | head -6
grep -rln "归口不是闭格手段" .scratch/wisp/issues/ docs/ 2>/dev/null | head -6
echo "== inbound C17 roster: panel.* methods in bridge.go =="
grep -n '"panel\.' internal/panel/bridge.go
echo "== bridge.go config methods =="
grep -n '"config\.' internal/panel/bridge.go
echo "== write-verb ruler on bridge.go (expect only ban-comment line) =="
grep -nE "\"(decide|allow|approve|commit|push|revert|stop|remove|apply|toggle|terminal)" internal/panel/bridge.go
echo "== lockedFieldFamilies (plugins family locked) =="
grep -n "lockedFieldFamilies" internal/panel/config_handlers.go | head -3
sed -n '105,115p' internal/panel/config_handlers.go
echo "== M1 K1-B attachment remove: demo control + Go inbound =="
grep -c "data-attach-del" design/doubao/demo/rightbar.js || true
grep -c "attachment.remove" internal/panel/bridge.go || true
echo "== M2 K6 diff apply: demo control + Go inbound =="
grep -c "rb-apply-btn" design/doubao/demo/rb-review.js || true
grep -cnE "apply|accept|revert" internal/panel/bridge.go || true
echo "== M3 K5 plugin toggle: demo control + Go inbound + locked family =="
grep -c "data-plg-toggle" design/doubao/demo/rb-plugins.js || true
grep -cin "plugin" internal/panel/bridge.go || true
echo "== M4 K7 terminal exec: demo control + Go inbound =="
grep -c "data-terminal-add" design/doubao/demo/rb-terminal.js || true
grep -cin "terminal" internal/panel/bridge.go || true
echo "== M5 K3 panel allow: demo controls + Go inbound =="
grep -nE "允许|拒绝|批准|approve|allow" design/doubao/demo/rightbar.js | head -5
grep -cnE "decide|allow|approve" internal/panel/bridge.go || true
echo "== K4 file tree Go-side write verbs (expect 0; positive control os.Open/os.Lstat) =="
grep -cnE "os\.Remove|os\.Rename|os\.WriteFile|os\.Create|RemoveAll" internal/panel/tree_reader_190.go || true
grep -cE "os\.Open|os\.Lstat" internal/panel/tree_reader_190.go || true
echo "== approval display side is read-only: ApprovalCardView fields present =="
grep -c "json:" internal/panel/approval.go || true
echo "== K8/K9/K10 Go inbound write verbs (expect 0; positive control json count) =="
grep -cnE "\"(decide|allow|approve|stop|kill|remove)" internal/panel/subagent_roster_197.go || true
grep -c "json:" internal/panel/subagent_roster_197.go || true
exit 0
