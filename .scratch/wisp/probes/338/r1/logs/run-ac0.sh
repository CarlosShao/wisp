#!/bin/sh
# 票 338 落地腿 338-r1 · AC#0 起手现跑 v2（v1 的 blob 循环被 Windows 反斜杠路径打断，本枚修正）
# 尺：/d/work/base/gopath/bin/gofumpt.exe v0.12.0 (go1.27.1)；gofmt＝同版 go
set -u
GOFUMPT="/d/work/base/gopath/bin/gofumpt.exe"
BLOB=/d/tmp/wisp/20261018/338r1/blobcheck
W=/d/tmp/wisp/20261018/338r1/wtest
mkdir -p "$BLOB" "$W"
cd "/d/work/workspace/projects plans/Wisp" || exit 9
echo "HEAD=$(git rev-parse --short HEAD)"

echo '===== [0] env at start ====='
echo "GOTMPDIR=[${GOTMPDIR:-unset}] TMPDIR=[${TMPDIR:-unset}]"

echo '===== [1] 工作树名册：gofumpt -l internal cmd ====='
"$GOFUMPT" -l internal cmd > "$W/roster-worktree-raw.txt"
sed -i 's/\\/\//g' "$W/roster-worktree-raw.txt"
sort "$W/roster-worktree-raw.txt" > "$W/roster-worktree.txt"
cat "$W/roster-worktree.txt"
echo "rows_worktree=$(wc -l < "$W/roster-worktree.txt") rc=$?"

echo '===== [2] 同名册按跟踪集过滤（CI 的分母＝跟踪件） ====='
git ls-files internal cmd > "$W/tracked.txt"
: > "$W/roster-tracked.txt"
while IFS= read -r f; do
  if grep -qxF "$f" "$W/tracked.txt"; then echo "$f" >> "$W/roster-tracked.txt"; else echo "$f  <== UNTRACKED(不在分母)" >> "$W/roster-untracked.txt"; fi
done < "$W/roster-worktree.txt"
cat "$W/roster-tracked.txt"
echo "rows_tracked=$(wc -l < "$W/roster-tracked.txt")"
cat "$W/roster-untracked.txt" 2>/dev/null

echo '===== [3] 逐枚 blob 定性：git show HEAD:<f> -> 仓外文件 -> gofumpt -l ====='
while IFS= read -r f; do
  flat=$(echo "$f" | tr '/\' '__')
  out="$BLOB/$flat"
  git show "HEAD:$f" > "$out"
  hascrlf=$(awk '/\r/{c++} END{print c+0}' "$out")
  rows=$("$GOFUMPT" -l "$out" | wc -l)
  if [ "$rows" -eq 0 ]; then verdict=BLOB_CLEAN_fake_dirty; else verdict=BLOB_DIRTY_real_dirty; fi
  echo "$f | bytes=$(wc -c < "$out") | crlf_lines=$hascrlf | gofumpt_l_rows=$rows | $verdict"
done < "$W/roster-tracked.txt"

echo '===== [4] 两枚目标：仓外副本 gofumpt -w 后跑非空白尺（⛔ 不碰工作树） ====='
for f in internal/tools/task_output_pointer_notice_errtext_174r7_test.go cmd/wisp/panel_transport_35r2_test.go; do
  flat=$(echo "$f" | tr '/\' '__')
  before="$W/$flat.before"; after="$W/$flat.after"
  cp "$f" "$before"
  cp "$before" "$after"
  "$GOFUMPT" -w "$after"; rcw=$?
  nb_before=$(tr -d ' \t' < "$before" | wc -c)
  nb_after=$(tr -d ' \t' < "$after" | wc -c)
  out=$("$GOFUMPT" -l "$after"); rows_after=$(printf '%s' "$out" | grep -c . || true)
  crlf_after=$(awk '/\r/{c++} END{print c+0}' "$after")
  echo "--- $f"
  echo "    rc_w=$rcw before_bytes=$(wc -c < "$before") after_bytes=$(wc -c < "$after")"
  echo "    nonblank_bytes before=$nb_before after=$nb_after"
  echo "    gofumpt_l_after_rows=$rows_after crlf_lines_after=$crlf_after"
  echo "    RULER(diff <(tr -d ' \\t' <before>) <(tr -d ' \\t' <after>)):"
  if diff <(tr -d ' \t' < "$before") <(tr -d ' \t' < "$after") > "$W/ruler-$flat.txt"; then
    echo "    RULER_RC=0 EMPTY_OUTPUT=YES  <-- 纯空白改动，放行"
  else
    echo "    RULER_RC=$? EMPTY_OUTPUT=NO  <-- 非空白字节变了（本腿禁区）"
  fi
  echo "    numstat_shape: added=$(diff "$before" "$after" | grep -c '^>' || true) removed=$(diff "$before" "$after" | grep -c '^<' || true)"
done
echo '===== [5] 另四枚 CRLF 件：逐字节自证未动（hash of HEAD blob vs worktree content） ====='
for f in internal/agent/approval/pending_read.go internal/agent/tools.go internal/risk/provenance.go cmd/wisp/models.go; do
  flat=$(echo "$f" | tr '/\' '__')
  git show "HEAD:$f" > "$BLOB/$flat"
  h_blob=$(sha1sum < "$BLOB/$flat" | cut -d' ' -f1)
  h_wt=$(sha1sum < "$f" | cut -d' ' -f1)
  # 归一化换行后比对：CRLF 假脏＝工作树仅换行符不同
  tr -d '\r' < "$f" > "$BLOB/$flat.lf"
  h_wt_lf=$(sha1sum < "$BLOB/$flat.lf" | cut -d' ' -f1)
  echo "$f | blob_sha1=$h_blob | wt_lf_normalized_sha1=$h_wt_lf | match=$([ "$h_blob" = "$h_wt_lf" ] && echo YES || echo NO) | raw_wt_sha1=$h_wt"
done
