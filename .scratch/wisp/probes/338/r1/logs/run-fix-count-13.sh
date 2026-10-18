#!/bin/sh
# 票 338 落地腿 338-r1 · 更正尺：全集去重枚数的射程＝13 笔（含 6cfbd5db），先前那份 40 枚只覆盖 12 笔
set -u
cd "/d/work/workspace/projects plans/Wisp" || exit 9
ALLOW='^(\.scratch/wisp/probes/338/r1/|internal/tools/task_output_pointer_notice_errtext_174r7_test\.go|cmd/wisp/panel_transport_35r2_test\.go|\.scratch/wisp/issues/338-)'
CS="0b1d1f0b 87521cb4 3d26240c 93d0f39c ee75c23f fa5f008f b84dedd0 787dbfb9 b4844295 91d08c3d a3994e28 5e18e1e2 6cfbd5db"
n=$(printf '%s\n' $CS | grep -c .)
for c in $CS; do git show --name-only --format= "$c" | grep .; done | sort -u > /d/tmp/wisp/20261018/338r1/all-paths-13.txt
echo "commits_counted=$n"
echo "unique_paths_13=$(wc -l < /d/tmp/wisp/20261018/338r1/all-paths-13.txt)"
echo "overreach_13=$(grep -vc -E "$ALLOW" /d/tmp/wisp/20261018/338r1/all-paths-13.txt || true)"
echo "unique_paths_12=$(wc -l < /d/tmp/wisp/20261018/338r1/all-paths-12.txt)"
echo '--- 6cfbd5db 带入的两枚（＝40→42 的那两枚） ---'
git show --name-only --format= 6cfbd5db | grep .
