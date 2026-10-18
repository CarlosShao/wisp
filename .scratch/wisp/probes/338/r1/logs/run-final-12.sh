#!/bin/sh
set -u
cd "/d/work/workspace/projects plans/Wisp" || exit 9
ALLOW='^(\.scratch/wisp/probes/338/r1/|internal/tools/task_output_pointer_notice_errtext_174r7_test\.go|cmd/wisp/panel_transport_35r2_test\.go|\.scratch/wisp/issues/338-)'
echo "HEAD=$(git rev-parse --short HEAD)"
for c in 0b1d1f0b 87521cb4 3d26240c 93d0f39c ee75c23f fa5f008f b84dedd0 787dbfb9 b4844295 91d08c3d a3994e28 5e18e1e2; do
  all=$(git show --name-only --format= "$c" | grep -c .)
  over=$(git show --name-only --format= "$c" | grep -v -E "$ALLOW" | grep -c .)
  printf '%s files=%s OVERREACH=%s\n' "$c" "$all" "$over"
done
echo '--- 全集（12 笔）去重 ⊆ 三处写面 ---'
for c in 0b1d1f0b 87521cb4 3d26240c 93d0f39c ee75c23f fa5f008f b84dedd0 787dbfb9 b4844295 91d08c3d a3994e28 5e18e1e2; do
  git show --name-only --format= "$c" | grep .
done | sort -u > /d/tmp/wisp/20261018/338r1/all-paths-12.txt
echo "unique=$(wc -l < /d/tmp/wisp/20261018/338r1/all-paths-12.txt) overreach=$(grep -vc -E "$ALLOW" /d/tmp/wisp/20261018/338r1/all-paths-12.txt || true)"
echo '--- 写面残留／票框 ---'
echo "my_uncommitted_rows=$(git status --porcelain -- .scratch/wisp/probes/338/r1 internal/tools/task_output_pointer_notice_errtext_174r7_test.go cmd/wisp/panel_transport_35r2_test.go .scratch/wisp/issues | grep -c '^\s*[MAD?]')"
F=".scratch/wisp/issues/338-two-tracked-files-are-gofumpt-dirty-in-blob-form-so-cis-format-step-goes-red-on-stock.md"
echo "boxes_unchecked=$(grep -c '^- \[ \]' "$F") boxes_checked=$(grep -c '^- \[x\]' "$F" || true)"
echo "report_lines=$(wc -l < .scratch/wisp/probes/338/r1/report.md) logs=$(ls .scratch/wisp/probes/338/r1/logs | wc -l)"
gofumpt_rows=$(/d/work/base/gopath/bin/gofumpt.exe -l internal cmd | wc -l)
echo "gofumpt_l_now=$gofumpt_rows"
