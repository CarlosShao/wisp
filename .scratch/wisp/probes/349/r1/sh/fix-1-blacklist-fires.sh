#!/bin/sh
# 票 349 AC#1 第①枚：internal/risk/blacklist.go:310 注释点名 `fires()`＝包内⛔ 此符号。
# 现读该句断言的＝"上面那两枚关系尺决定了 class/tier/overridability"。包内真身＝aEq/aUnder
# （二者都 delegate 到 aRelates），⇒ 改成盘上真名，⛔ 保留指向不存在符号的裸名、⛔ 新增符号。
# 内容锚（⛔ 绝对行号）；落刀前量命中＝1，落刀后旧锚＝0 且新锚＝1；只动这一行。
set -u
cd '/d/work/workspace/projects plans/Wisp' || exit 90
F=internal/risk/blacklist.go
OLD='fires()/aRelates above decided'
NEW='aEq/aUnder (via aRelates) above decided'
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
echo "== line 310 now =="
sed -n '310p' "$F"
echo "== line length (want <= 100) =="
sed -n '310p' "$F" | wc -c
echo "== diffstat: only comment lines? =="
git diff -U0 -- "$F" | grep -E '^[+-]' | grep -vE '^(\+\+\+|---)' || true
echo "== non-comment changed lines (want 0) =="
git diff -U0 -- "$F" | grep -E '^[+-]' | grep -vE '^(\+\+\+|---)' | grep -vcE '^[+-][[:space:]]*//' || true
