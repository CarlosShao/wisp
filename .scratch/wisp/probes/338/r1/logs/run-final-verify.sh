#!/bin/sh
# 票 338 落地腿 338-r1 · 交件前最后一轮自证（三把尺＋必答四问的盘上复量）
set -u
cd "/d/work/workspace/projects plans/Wisp" || exit 9
GOFUMPT="/d/work/base/gopath/bin/gofumpt.exe"
F=.scratch/wisp/issues/338-two-tracked-files-are-gofumpt-dirty-in-blob-form-so-cis-format-step-goes-red-on-stock.md
unset GOTMPDIR TMPDIR

echo '===== [1] 逐枚提交号＋git show --name-only（我这一串 9 笔） ====='
for c in 0b1d1f0b 87521cb4 3d26240c 93d0f39c ee75c23f fa5f008f b84dedd0 787dbfb9 b4844295; do
  echo "----- $c"
  git show --name-only --format="%h %s" $c | head -1 | cut -c1-100
  git show --name-only --format= $c | grep "^" | sed 's/^/      /'
  over=$(git show --name-only --format= $c | grep -v -E "^(\.scratch/wisp/probes/338/r1/|internal/tools/task_output_pointer_notice_errtext_174r7_test\.go|cmd/wisp/panel_transport_35r2_test\.go|\.scratch/wisp/issues/338-)" | grep -c "^\.")
  echo "      OVERREACH=$over"
done

echo '===== [2] 必答①：两枚的非空白尺逐字（用 logs 里留档的改前⇄改后重跑） ====='
D1=/d/tmp/wisp/20261018/338r1/ac1
D2=/d/tmp/wisp/20261018/338r1/ac1b
diff <(tr -d ' \t' < "$D1/tools.before") <(tr -d ' \t' < "$D1/tools.after")
echo "TOOLS_RULER_RC=$? (期望 0＝空输出)"
echo "tools nonblank before=$(tr -d ' \t' < "$D1/tools.before" | wc -c) after=$(tr -d ' \t' < "$D1/tools.after" | wc -c)"
diff <(tr -d ' \t' < "$D2/panel.before") <(tr -d ' \t' < "$D2/panel.after") > /dev/null
echo "PANEL_RULER_RC=$? (甲形＝预期 1＝非空白有变，票面已让路)"
echo "panel nonblank before=$(tr -d ' \t' < "$D2/panel.before" | wc -c) after=$(tr -d ' \t' < "$D2/panel.after" | wc -c)"

echo '===== [3] 必答②：名册交件侧现量＋与改前名册作差 ====='
"$GOFUMPT" -l internal cmd | sed 's/\\/\//g' | sort > /d/tmp/wisp/20261018/338r1/post-now.txt
cat /d/tmp/wisp/20261018/338r1/post-now.txt
echo "rows_now=$(wc -l < /d/tmp/wisp/20261018/338r1/post-now.txt)"
echo "--- 改前(6) 逐名 ---"; cat "$D1/pre.tracked.txt" 2>/dev/null || cat /d/tmp/wisp/20261018/338r1/wtest/roster-tracked.txt
echo "--- 作差：改前有、现无 ---"
comm -23 <(sort "$D1/pre.tracked.txt") /d/tmp/wisp/20261018/338r1/post-now.txt
echo "--- 作差：现有、改前无（必须空） ---"
comm -13 <(sort "$D1/pre.tracked.txt") /d/tmp/wisp/20261018/338r1/post-now.txt
echo "--- 逐枚复量＝0 枚 ---"
echo "tools_l=$("$GOFUMPT" -l internal/tools/task_output_pointer_notice_errtext_174r7_test.go | grep -c . || true) panel_l=$("$GOFUMPT" -l cmd/wisp/panel_transport_35r2_test.go | grep -c . || true)"

echo '===== [4] 必答③：另四枚 CRLF 件一字未动的尺 ====='
for f in internal/agent/approval/pending_read.go internal/agent/tools.go internal/risk/provenance.go cmd/wisp/models.go; do
  echo "$f | wt_sha256=$(sha256sum < "$f" | cut -c1-16) | git_status=$(git status --porcelain -- "$f" | wc -l) | eol=$(git ls-files --eol -- "$f" | awk '{print $1" "$2}')"
done
echo "--- blob 侧仍 CLEAN（＝假脏、票面禁区，未动） ---"
for f in internal/agent/approval/pending_read.go internal/agent/tools.go internal/risk/provenance.go cmd/wisp/models.go; do
  flat=$(printf '%s' "$f" | tr '/' '_')
  git show "HEAD:$f" > "/d/tmp/wisp/20261018/338r1/blobcheck/$flat"
  echo "$f blob_gofumpt_l_rows=$("$GOFUMPT" -l "/d/tmp/wisp/20261018/338r1/blobcheck/$flat" | grep -c . || true)"
done

echo '===== [5] 交件侧 eol＋numstat＋CR（两枚） ====='
git ls-files --eol -- internal/tools/task_output_pointer_notice_errtext_174r7_test.go cmd/wisp/panel_transport_35r2_test.go | cat
git show --numstat --format= 87521cb4 -- internal/tools/task_output_pointer_notice_errtext_174r7_test.go | cat
git show --numstat --format= 93d0f39c -- cmd/wisp/panel_transport_35r2_test.go | cat
echo "tools_CR=$(awk '/\r/{c++} END{print c+0}' internal/tools/task_output_pointer_notice_errtext_174r7_test.go) panel_CR=$(awk '/\r/{c++} END{print c+0}' cmd/wisp/panel_transport_35r2_test.go)"

echo '===== [6] report＋logs 尺 ====='
wc -l -c .scratch/wisp/probes/338/r1/report.md | cat
echo "logs_count=$(ls .scratch/wisp/probes/338/r1/logs | wc -l)"
SEC=$(printf '\302\247')   # 分节号字符＝UTF-8 U+00A7；本脚本⛔ 写出它的字面（⛔ 描述它时把它打出来）
echo "report_section_char=$(grep -c "$SEC" .scratch/wisp/probes/338/r1/report.md || true)"
echo "logs_section_char_total=$(cat .scratch/wisp/probes/338/r1/logs/* | grep -c "$SEC" || true)"
echo "this_script_own_literal_count=$(grep -c "$SEC" .scratch/wisp/probes/338/r1/logs/run-final-verify.sh || true) (修后期望 0＝连本脚本也不留字面)"
zero=0; for f in .scratch/wisp/probes/338/r1/logs/* .scratch/wisp/probes/338/r1/report.md; do [ -s "$f" ] || { echo "ZERO_BYTE=$f"; zero=1; }; done; echo "any_zero_byte=$zero"
maxline=$(wc -l .scratch/wisp/probes/338/r1/logs/* | grep -v total | awk '{print $1}' | sort -n | tail -1)
echo "max_log_lines=$maxline (帽 160)"

echo '===== [7] 票 338＝四格一枚没翻＋只追加 ====='
echo "boxes_unchecked=$(grep -c '^- \[ \]' "$F") boxes_checked=$(grep -c '^- \[x\]' "$F" || true)"
git log --oneline -3 -- "$F" | cat
git show --numstat --format= b4844295 -- "$F" | cat
echo '===== [8] 我没在 docs/evidence/s1/ 落任何件 ====='
git status --porcelain -- docs/evidence/s1 | cat
echo "rows=$(git status --porcelain -- docs/evidence/s1 | grep -c . || true)"
