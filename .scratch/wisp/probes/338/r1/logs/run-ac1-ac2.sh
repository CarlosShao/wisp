#!/bin/sh
# 票 338 落地腿 338-r1 · AC#1 归位（只一枚）＋ AC#2 成对门（改前⇄改后同台面）
# ★panel 枚按 lead 裁定前一字不动（其 gofumpt -w 非空白不达标＝AC#0 已具名上报）
set -u
GOFUMPT="/d/work/base/gopath/bin/gofumpt.exe"
T=internal/tools/task_output_pointer_notice_errtext_174r7_test.go
P=cmd/wisp/panel_transport_35r2_test.go
D=/d/tmp/wisp/20261018/338r1/ac1
mkdir -p "$D"
cd "/d/work/workspace/projects plans/Wisp" || exit 9
filter_tracked() {
  sed 's/\\/\//g' | sort | while IFS= read -r f; do
    git ls-files --error-unmatch "$f" >/dev/null 2>&1 && echo "$f"
  done
}

echo "HEAD_at_start=$(git rev-parse --short HEAD)"
echo '===== [0] 起手：两枚目标相对 HEAD 无本地脏（否则归位会把别腿字节卷进来） ====='
git status --porcelain -- "$T" "$P" | cat
echo "status_rows=$(git status --porcelain -- "$T" "$P" | grep -c . || true)"
echo "GOTMPDIR=[${GOTMPDIR:-unset}] TMPDIR=[${TMPDIR:-unset}]"

echo '===== [1] 改前名册（同台面，逐名） ====='
"$GOFUMPT" -l internal cmd > "$D/pre.raw"
sed 's/\\/\//g' "$D/pre.raw" | sort > "$D/pre.txt"
cat "$D/pre.txt"
echo "pre_rows_total=$(wc -l < "$D/pre.txt")"
filter_tracked < "$D/pre.txt" > "$D/pre.tracked.txt"
cat "$D/pre.tracked.txt"
echo "pre_rows_tracked=$(wc -l < "$D/pre.tracked.txt")"

echo '===== [2] 改前字节留档＋哈希（仓外） ====='
cp "$T" "$D/tools.before"
cp "$P" "$D/panel.before"
echo "tools_before_sha=$(sha256sum < "$D/tools.before" | cut -d' ' -f1) bytes=$(wc -c < "$D/tools.before")"
echo "panel_before_sha=$(sha256sum < "$D/panel.before" | cut -d' ' -f1) bytes=$(wc -c < "$D/panel.before")"
for f in internal/agent/approval/pending_read.go internal/agent/tools.go internal/risk/provenance.go cmd/wisp/models.go; do
  echo "crlf4|$f|wt_sha=$(sha256sum < "$f" | cut -d' ' -f1)|head_lf_sha=$(git show "HEAD:$f" | sha256sum | cut -d' ' -f1)"
done

echo '===== [3] 唯一写操作：gofumpt -w 只打这一枚字面路径 ====='
"$GOFUMPT" -w "$T"
echo "rc_w=$?"

echo '===== [4] ★非空白零变化尺（改前⇄改后逐字，必须空输出） ====='
cp "$T" "$D/tools.after"
diff <(tr -d ' \t' < "$D/tools.before") <(tr -d ' \t' < "$D/tools.after") > "$D/ruler.out" 2>&1
echo "RULER_RC=$?"
echo "RULER_BYTES=$(wc -c < "$D/ruler.out") RULER_ROWS=$(wc -l < "$D/ruler.out")"
echo '--- ruler 逐字输出（应为空） ---'
cat "$D/ruler.out"
echo "--- 逐字结束 ---"
echo "nonblank_bytes before=$(tr -d ' \t' < "$D/tools.before" | wc -c) after=$(tr -d ' \t' < "$D/tools.after" | wc -c)"

echo '===== [5] 逐枚 gofumpt -l 复量（目标枚必须 0 行）＋幂等二发 ====='
echo "tools_l_rows=$("$GOFUMPT" -l "$T" | grep -c . || true)"
echo "panel_l_rows=$("$GOFUMPT" -l "$P" | grep -c . || true)"
cp "$D/tools.after" "$D/tools.after2"; "$GOFUMPT" -w "$D/tools.after2"
echo "idempotent_second_w=$([ "$(sha256sum < "$D/tools.after"|cut -d' ' -f1)" = "$(sha256sum < "$D/tools.after2"|cut -d' ' -f1)" ] && echo YES || echo NO)"

echo '===== [6] git 面自证：numstat / eol / diff 形状 ====='
git diff --numstat -- "$T" | cat
git diff --numstat -- "$P" | cat
echo "panel_numstat_rows=$(git diff --numstat -- "$P" | grep -c . || true) (必须 0＝一字未动)"
echo '--- ls-files --eol 交件侧 ---'
git ls-files --eol -- "$T" "$P" | cat
echo '--- 新字节 CR 计数（必须 0） ---'
echo "tools_CR_lines=$(awk '/\r/{c++} END{print c+0}' "$T")"
echo '--- 改后逐字 diff（只此一枚） ---'
git diff -- "$T" | cat

echo '===== [7] 改后名册（同台面，逐名）＋作差 ====='
"$GOFUMPT" -l internal cmd > "$D/post.raw"
sed 's/\\/\//g' "$D/post.raw" | sort > "$D/post.txt"
cat "$D/post.txt"
echo "post_rows_total=$(wc -l < "$D/post.txt")"
filter_tracked < "$D/post.txt" > "$D/post.tracked.txt"
cat "$D/post.tracked.txt"
echo "post_rows_tracked=$(wc -l < "$D/post.tracked.txt")"
echo '--- 作差：改前有改后无（＝被归位） ---'
comm -23 "$D/pre.tracked.txt" "$D/post.tracked.txt"
echo '--- 作差：改后有改前无（＝多动的，必须空） ---'
comm -13 "$D/pre.tracked.txt" "$D/post.tracked.txt"

echo '===== [8] CRLF 四枚：写操作后复量（一字未动的尺） ====='
for f in internal/agent/approval/pending_read.go internal/agent/tools.go internal/risk/provenance.go cmd/wisp/models.go; do
  wt=$(sha256sum < "$f" | cut -d' ' -f1)
  echo "crlf4|$f|wt_sha=$wt"
done
echo '--- git status（四枚相对 HEAD 的形状必须与改前同一） ---'
git status --porcelain -- internal/agent/approval/pending_read.go internal/agent/tools.go internal/risk/provenance.go cmd/wisp/models.go | cat
