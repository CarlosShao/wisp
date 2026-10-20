#!/bin/sh
# 353-v1 R27 = 票 322 AC#1 逐字 ⇄ 仪器件弃权句 ⇄ "报窗口"三字全仓命中（量的；只读）
cd '/d/work/workspace/projects plans/Wisp' || exit 9
echo "===== 1) 票 322 票面 AC#1 那一格逐字（先例形＝票面 bullet 里写没写那味）====="
git cat-file blob 'HEAD:.scratch/wisp/issues/322-privatefile-open-truncates-through-a-swapped-leaf-done.md' | sed -n '30p' | cut -c1-400
echo "===== 2) 那味真身＝仪器件里的弃权句（逐字行号）====="
git cat-file blob 'HEAD:internal/winsec/openleaf_322_other_test.go' | grep -nE 'Nothing here claims|not evidence about a window|abstain' | cut -c1-200
echo "===== 3) '报窗口'三字在票 322⇄票 353 面命中枚数（两把尺）====="
T322='.scratch/wisp/issues/322-privatefile-open-truncates-through-a-swapped-leaf-done.md'
T353='.scratch/wisp/issues/353-pinned-name-swap-after-seal-deletes-an-object-the-caller-never-named.md'
printf '  票322面=%s 票353面=%s  353a1草案=%s\n' \
  "$(git cat-file blob "HEAD:$T322" | grep -c '报窗口' || true)" \
  "$(git cat-file blob "HEAD:$T353" | grep -c '报窗口' || true)" \
  "$(git cat-file blob 'HEAD:.scratch/wisp/probes/353/a1/logs/04-ac2-instrument-shape-draft.txt' | grep -c '报窗口' || true)"
