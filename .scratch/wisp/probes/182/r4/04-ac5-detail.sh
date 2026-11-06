# 182 r4 leg — script 04: AC#5 detail lines + K10/K7 backing + K6/K4 demo control widening check
set -u
cd "/d/work/workspace/projects plans/Wisp" || exit 9
echo "== bridge.go apply/accept/revert lines (context) =="
grep -nE "apply|accept|revert" internal/panel/bridge.go
echo "== bridge.go decide/allow/approve lines (context) =="
grep -nE "decide|allow|approve" internal/panel/bridge.go
echo "== K6 demo: any commit/push control in rb-review.js =="
grep -nE "commit|push|提交|推送" design/doubao/demo/rb-review.js || echo "(none)"
echo "== K4 demo: any delete/move control in rb-files.js =="
grep -nE "删除|移动|delete|move|rename" design/doubao/demo/rb-files.js || echo "(none)"
echo "== K10 Go-side browser surface in panel (0 expected; positive control WebView in cmd/wisp) =="
grep -rli "browser" internal/panel/ --include=*.go | grep -v _test.go | wc -l
grep -rli "WebView" cmd/wisp/ --include=*.go | grep -v _test.go | wc -l
echo "== K7 backing: registration loops in run.go (no Shell expected) =="
grep -n "tools.Builtin.*Entries(" cmd/wisp/run.go
echo "== K1 piles: composer.go workspace/attachments anchors =="
grep -n "json:\"workspace\"\|json:\"attachments\"" internal/panel/composer.go
echo "== K3 stop: any stop method in panel inbound =="
grep -cin "\"panel.task.stop\"\|task.stop\|methodStop" internal/panel/bridge.go || true
echo "== approval.go ApprovalCardView name check =="
grep -n "type ApprovalCardView struct" internal/panel/approval.go
exit 0
