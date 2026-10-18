#!/bin/sh
# 票 338 落地腿 338-r1 · 越界尺的**形状更正**（具名：先前几笔我用的是 grep -c "^\\."＝只数点开头的路径，
# ⇒ 若真越界落在 cmd/**／internal/**（不以点开头）会漏数＝假 0。本枚改用"数全部非空行"＝正确分母。）
set -u
cd "/d/work/workspace/projects plans/Wisp" || exit 9
ALLOW='^(\.scratch/wisp/probes/338/r1/|internal/tools/task_output_pointer_notice_errtext_174r7_test\.go|cmd/wisp/panel_transport_35r2_test\.go|\.scratch/wisp/issues/338-)'
for c in 0b1d1f0b 87521cb4 3d26240c 93d0f39c ee75c23f fa5f008f b84dedd0 787dbfb9 b4844295 91d08c3d a3994e28; do
  all=$(git show --name-only --format= "$c" | grep -c .)
  over=$(git show --name-only --format= "$c" | grep -v -E "$ALLOW" | grep -c .)
  printf '%s files=%s OVERREACH_corrected=%s\n' "$c" "$all" "$over"
  git show --name-only --format= "$c" | grep -v -E "$ALLOW" | grep . | sed 's/^/      !! /'
done
echo '===== 全集自证：11 笔的全部路径去重 ⊆ 三处写面 ====='
for c in 0b1d1f0b 87521cb4 3d26240c 93d0f39c ee75c23f fa5f008f b84dedd0 787dbfb9 b4844295 91d08c3d a3994e28; do
  git show --name-only --format= "$c" | grep .
done | sort -u | tee /d/tmp/wisp/20261018/338r1/all-paths.txt | grep -v -E "$ALLOW" | grep . | sed 's/^/!! OVERREACH: /'
echo "total_unique_paths=$(wc -l < /d/tmp/wisp/20261018/338r1/all-paths.txt)"
echo "overreach_unique=$(grep -vc -E "$ALLOW" /d/tmp/wisp/20261018/338r1/all-paths.txt || true)"
echo '===== 票 338 框（读-only，⛔ 翻） ====='
F=".scratch/wisp/issues/338-two-tracked-files-are-gofumpt-dirty-in-blob-form-so-cis-format-step-goes-red-on-stock.md"
echo "unchecked=$(grep -c '^- \[ \]' "$F")"
echo "checked=$(grep -c '^- \[x\]' "$F" || true)"
