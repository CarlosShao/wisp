#!/bin/sh
set -e
cd "/d/work/workspace/projects plans/Wisp"
T=.scratch/wisp/probes/orch/20261018/txt
cat "$T/A1028.txt" >> docs/reports/pending-and-issues.md
cat "$T/HANDOVER-40em.md" >> docs/reports/HANDOVER.md
printf '%s\n' "收 166-r1⇒票 166 AC#2 仍不成立（腿自判＋我两棵树复跑：HEAD 四枚 PASS ⇄ 基线 no tests to run）｜第二波补两枚已放＝227-r1（写·internal/config 测试面）∥340-a1（只读·票 340 三形表，AC#2 人工批准门未开）｜形状＝在飞 4≤8、写腿 3≥3、cmd-wisp/panel 零写腿｜★派单公共块补两条定式：跑门前 unset GOTMPDIR/TMPDIR＋add/commit 两处 pathspec 写死字面路径｜账 A1028＋停车点 4.0em｜推送已授权仍未执行" > "$T/40em-msg.txt"
git add -- docs/reports/pending-and-issues.md docs/reports/HANDOVER.md "$T/A1028.txt" "$T/HANDOVER-40em.md" "$T/40em-msg.txt" "$T/227r1-head.txt" "$T/227r1-specific.txt" "$T/227r1-dispatch.txt" "$T/227r1-send.txt" "$T/340a1-head.txt" "$T/340a1-specific.txt" "$T/340a1-dispatch.txt" "$T/340a1-send.txt" "$T/assemble-wave2.sh" "$T/common-20261018.txt" "$T/166-verify.sh"
echo "=== staged audit ==="
git diff --cached --name-only | grep -v -E "^(docs/reports/(pending-and-issues|HANDOVER)\.md|\.scratch/wisp/probes/orch/20261018/)" || echo "  only-mine"
git commit -F "$T/40em-msg.txt" -- docs/reports/pending-and-issues.md docs/reports/HANDOVER.md "$T/A1028.txt" "$T/HANDOVER-40em.md" "$T/40em-msg.txt" "$T/227r1-head.txt" "$T/227r1-specific.txt" "$T/227r1-dispatch.txt" "$T/227r1-send.txt" "$T/340a1-head.txt" "$T/340a1-specific.txt" "$T/340a1-dispatch.txt" "$T/340a1-send.txt" "$T/assemble-wave2.sh" "$T/common-20261018.txt" "$T/166-verify.sh"
echo "=== after ==="
git log -1 --format='%h %s' | cut -c1-70
echo "unpushed=$(git rev-list --count origin/dev..HEAD)"
git show --stat --oneline HEAD | tail -4
