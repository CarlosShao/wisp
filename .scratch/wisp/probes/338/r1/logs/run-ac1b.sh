#!/bin/sh
# 票 338 落地腿 338-r1 · AC#1b 归位第二枚（panel 枚）＝lead 12:1x 票面裁语「按甲形放行 token 级改动」
# 补偿条件逐条落实：① 改动只来自 gofumpt -w 本身（⛔ 手改）② 逐发 -l 直到空输出＋每发 rc=
#                  ⑤ 交件侧 gofmt -d 仍 0 字节
set -u
GOFUMPT="/d/work/base/gopath/bin/gofumpt.exe"
P=cmd/wisp/panel_transport_35r2_test.go
D=/d/tmp/wisp/20261018/338r1/ac1b
mkdir -p "$D"
cd "/d/work/workspace/projects plans/Wisp" || exit 9
unset GOTMPDIR TMPDIR

echo "HEAD=$(git rev-parse --short HEAD)"
echo '===== [0] 起手：该枚相对 HEAD 无本地脏 ====='
echo "status_rows=$(git status --porcelain -- "$P" | grep -c . || true)"
git diff --stat HEAD -- "$P" | cat
cp "$P" "$D/panel.before"
echo "before_sha=$(sha256sum < "$D/panel.before" | cut -d' ' -f1) bytes=$(wc -c < "$D/panel.before")"
echo "before_gofumpt_l_rows=$("$GOFUMPT" -l "$P" | grep -c . || true)"

echo '===== [1] 逐发 gofumpt -w 直到 -l 空输出（条件①＋②） ====='
i=0
while : ; do
  i=$((i+1))
  "$GOFUMPT" -w "$P"; rcw=$?
  lrows=$("$GOFUMPT" -l "$P" | grep -c . || true)
  echo "round$i rc_w=$rcw gofumpt_l_rows=$lrows sha=$(sha256sum < "$P" | cut -c1-16) bytes=$(wc -c < "$P")"
  [ "$lrows" -eq 0 ] && break
  [ "$i" -ge 6 ] && { echo "NO_CONVERGE_ABORT"; break; }
done
cp "$P" "$D/panel.after"

echo '===== [2] 条件⑤：交件侧 gofmt -d 必须 0 字节 ＋ gofmt -l 0 行 ====='
gofmt -d "$P" > "$D/gofmt-d-after.txt" 2>&1
echo "gofmt_d_bytes=$(wc -c < "$D/gofmt-d-after.txt")"
echo "gofmt_l_rows=$(gofmt -l "$P" | grep -c . || true)"

echo '===== [3] 非空白尺（甲形下预期不空＝记录逐字，票面已让路） ====='
diff <(tr -d ' \t' < "$D/panel.before") <(tr -d ' \t' < "$D/panel.after") > "$D/ruler.txt" 2>&1
echo "RULER_RC=$? rows=$(wc -l < "$D/ruler.txt")"
echo "nonblank_bytes before=$(tr -d ' \t' < "$D/panel.before" | wc -c) after=$(tr -d ' \t' < "$D/panel.after" | wc -c)"

echo '===== [4] 改动形状自查：除 gofumpt -w 产物外无别的字节（diff 全文逐字入库） ====='
git diff --numstat -- "$P" | cat
git diff -- "$P" > "$D/git-diff-panel.txt"
echo "diff_bytes=$(wc -c < "$D/git-diff-panel.txt")"
echo '--- 新增行的 token 类别清点（只应见 type ( ) 逗号 空行） ---'
grep '^+' "$D/git-diff-panel.txt" | grep -v '^+++' | tr -d ' \t' | sort | uniq -c | sort -rn | head -40
echo '--- 删除行清点 ---'
grep '^-' "$D/git-diff-panel.txt" | grep -v '^---' | tr -d ' \t' | sort | uniq -c | sort -rn | head -40

echo '===== [5] eol 自证：新字节无 CR、交件侧 i/lf ====='
echo "CR_lines=$(awk '/\r/{c++} END{print c+0}' "$P")"
git ls-files --eol -- "$P" | cat

echo '===== [6] 名册复量（交件侧，逐名） ====='
"$GOFUMPT" -l internal cmd | sed 's/\\/\//g' | sort
echo "rc_list=$?"

echo '===== [7] 四枚 CRLF ＋ tools 枚：本笔写操作前后一字未动 ====='
for f in internal/agent/approval/pending_read.go internal/agent/tools.go internal/risk/provenance.go cmd/wisp/models.go internal/tools/task_output_pointer_notice_errtext_174r7_test.go; do
  echo "$f | wt_sha=$(sha256sum < "$f" | cut -d' ' -f1)"
done
