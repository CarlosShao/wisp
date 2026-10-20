#!/bin/sh
# 353-v1 提交尺＝逐枚字面 pathspec 的 add⇄commit（⛔ 目录形／add -A／.／从共享 index 反推）
# 用法: sh 50-commit.sh "<file1> <file2> ..." "<D:\tmp msg path>"
cd '/d/work/workspace/projects plans/Wisp' || exit 9
LIST="$1"; MSG="$2"
# shellcheck disable=SC2086
set -- $LIST
git add -- "$@"
echo "ADD_RC=$?"
# shellcheck disable=SC2086
git commit -F "$MSG" -- "$@" > /d/tmp/wisp/20261020/353v1/commit-out.txt 2>&1
echo "COMMIT_RC=$?"
tail -5 /d/tmp/wisp/20261020/353v1/commit-out.txt
echo "=== git log -1 ==="
git log -1 --format='%H%n%s'
echo "=== git show --name-only ==="
git show --name-only --format='%h %s' HEAD | sed -n '1,60p'
echo "=== .go 枚数（尺＝grep -c '\.go$'）==="
git show --name-only --format= HEAD | grep -c '\.go$' || true
echo "=== 越界尺：本笔件名⛔ 属本腿写面者 ==="
git show --name-only --format= HEAD | grep -vcE '^(\.scratch/wisp/probes/353/v1/|docs/evidence/s1/353-v1-verdict\.md)' || true
