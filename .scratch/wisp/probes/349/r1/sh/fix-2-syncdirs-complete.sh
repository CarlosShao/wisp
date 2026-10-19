#!/bin/sh
# 票 349 AC#1 第②枚：internal/risk/syncdirs.go:453 注释写 `Complete()==false`＝包内⛔ 此方法。
# 真名＝SyncDetectionComplete（同文件 :472 定义、:442 那句注释与全包 15 处测试都在用它）
# ⇒ 纯小写"调用形"＝只有调用形尺抓得到（票面 AC#0 那把尺的现场样本）。改成真名，⛔ 动可执行字节。
# 锚形注意：新名把旧名当子串包着 ⇒ 旧锚取带前导空格的 " via Complete()==false"，
# 落刀后该锚必＝0（前面多了 SyncDetection），另用新锚证明落了刀。
set -u
cd '/d/work/workspace/projects plans/Wisp' || exit 90
F=internal/risk/syncdirs.go
OLD=' via Complete()==false'
NEW=' via SyncDetectionComplete()==false'
echo "== before: old-anchor hits (want 1) =="
grep -c -F "$OLD" "$F" || true
echo "== CR count before (want 0＝LF 文件) =="
tr -cd '\r' < "$F" | wc -c
export N="$OLD" R="$NEW"
perl -0pi -e 'my $q=quotemeta($ENV{N}); s/$q/$ENV{R}/' "$F" || echo "PERL_RC=NONZERO"
echo "== after: old-anchor hits (want 0) / new-anchor hits (want 1) =="
grep -c -F "$OLD" "$F" || true
grep -c -F "$NEW" "$F" || true
echo "== CR count after (want 0) =="
tr -cd '\r' < "$F" | wc -c
echo "== the two lines now (453 = the fixed one, 472 = the real definition) =="
sed -n '452,454p' "$F"
sed -n '472p' "$F"
echo "== line 453 length (want <= 100) =="
sed -n '453p' "$F" | wc -c
echo "== diff -U0 (only comment lines) =="
git diff -U0 -- "$F" | grep -E '^[+-]' | grep -vE '^(\+\+\+|---)' || true
echo "== non-comment changed lines (want 0) =="
git diff -U0 -- "$F" | grep -E '^[+-]' | grep -vE '^(\+\+\+|---)' | grep -vcE '^[+-][[:space:]]*//' || true
