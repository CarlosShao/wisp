#!/bin/sh
# 353-v1 提交尺（单枚⇄少数枚件）＝逐枚字面 pathspec（⛔ 目录形／add -A／.／从共享 index 反推）
# 用法: sh 52-commit-file.sh "<commit-msg-file Windows形>" <repo-relative-path> [<more paths>...]
cd '/d/work/workspace/projects plans/Wisp' || exit 9
MSG="$1"; shift
git add -- "$@"
echo "ADD_RC=$?"
git commit -F "$MSG" -- "$@" > /d/tmp/wisp/20261020/353v1/commit-out-last.txt 2>&1
echo "COMMIT_RC=$?"
tail -3 /d/tmp/wisp/20261020/353v1/commit-out-last.txt
echo "=== git log -1 ==="
git log -1 --format='%H %s'
echo "=== 本笔自证三把尺（尺＝git show --name-only）==="
git show --name-only --format='' HEAD | grep -c . | sed 's/^/  件数=/'
git show --name-only --format='' HEAD | grep -c '\.go$' | sed 's/^/  .go 枚数=/'
git show --name-only --format='' HEAD | grep -vcE '^(\.scratch/wisp/probes/353/v1/|docs/evidence/s1/353-v1-verdict\.md)' | sed 's/^/  越界枚数=/'
echo "=== 本笔逐枚件名 ==="
git show --name-only --format='' HEAD | sed 's/^/  /'
