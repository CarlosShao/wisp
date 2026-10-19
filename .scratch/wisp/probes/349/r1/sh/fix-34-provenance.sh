#!/bin/sh
# 票 349 AC#1 第③枚＋第④枚：internal/risk/provenance.go 两枚注释（:54 taskScopeID／:427 disposalScope）。
# 两枚都先追一手（本脚本第一段＝git log -S ＋逐提交 grep ＋"它当过定义吗"的反尺），⛔ 猜"原本想指什么"：
#   taskScopeID  ＝ 自 73b60571（票 19 落 C25 那笔）起就只在注释行里出现，全仓⛔ 任何定义
#                  ⇒ 它从来不是符号＝散文占位；盘上真身＝Detector(scopeID string)（provenance.go:858）
#                    ⇒ 占位换成那枚方法自己的形参名 scopeID（包内 GREEN-bound，⛔ 新造名）
#   disposalScope＝同一来路；真身＝别包类型 internal/plugin.DisposalScope（internal/plugin/disposal.go:129）
#                  及其方法 Defer（:186）。⇒ 按票面 AC#1"别包的要写全包路径"点名，并写成**散文形**：
#                  `DisposalScope(...)` 那种紧跟左括号的形状会被调用形尺再算一枚命中（本包⛔ 这个裸名）
#                  ＝自己造同族残留；`deferring this closure` 不带括号，尺上干净。
# ⛔ 新增任何符号去让注释变真、⛔ 动可执行字节。
# ★两枚旧锚都现量 before=1（唯一），⛔ 用绝对行号落刀；落刀后旧锚＝0、新锚＝1。
set -u
cd '/d/work/workspace/projects plans/Wisp' || exit 90
F=internal/risk/provenance.go
echo "== history probe (一手凭据) =="
echo "-- taskScopeID: oldest commit that ever added the string (repo-wide) --"
git log --all --format="%h %ad %s" --date=short -S 'taskScopeID' -- internal/risk | tail -2
echo "-- was taskScopeID EVER a definition? grep the introducing commit's tree --"
git grep -n 'taskScopeID' 73b60571 -- internal || echo "ZERO-HITS in 73b60571 tree"
echo "-- 'func taskScopeID' anywhere in repo history (want 0 commits) --"
git log --all --format="%h" -S 'func taskScopeID' | wc -l
echo "-- disposalScope: oldest commit + any definition form in that tree --"
git log --all --format="%h %ad %s" --date=short -S 'disposalScope' -- internal/risk | tail -2
git grep -nE 'func .*disposalScope|var disposalScope|disposalScope :=|disposalScope =' 73b60571 -- internal || echo "ZERO definitions in 73b60571 tree"
echo "-- real seams on disk today --"
grep -n 'func (p \*Provenance) Detector' internal/risk/provenance.go
grep -n 'type DisposalScope struct' internal/plugin/disposal.go
grep -n 'func (s \*DisposalScope) Defer' internal/plugin/disposal.go
grep -n 'no \*plugin.DisposalScope reaches' internal/risk/provenance.go
echo "########## 枚③ taskScopeID -> scopeID ##########"
export N='prov.Detector(taskScopeID)'
export R='prov.Detector(scopeID)'
echo "-- old anchor hits before (want 1):"; grep -c -F "$N" "$F" || true
perl -0pi -e 'my $q=quotemeta($ENV{N}); s/$q/$ENV{R}/' "$F" || echo "PERL_RC=NONZERO"
echo "-- old anchor hits after (want 0):"; grep -c -F "$N" "$F" || true
echo "-- new anchor hits after (want 1):"; grep -c -F "$R" "$F" || true
echo "########## 枚④ disposalScope -> 别包真身（全包路径＋散文形） ##########"
# ★旧锚是**行内连续串**：`e.g.` 在 :426 行尾、调用形在 :427 行首＝跨行不连续 ⇒ 锚只取 :427 那一枚
export N2='disposalScope.Defer(func() { _ = scope.Close() })'
export R2="the task's internal/plugin.DisposalScope deferring this closure"
echo "-- old anchor hits before (want 1):"; grep -c -F "$N2" "$F" || true
perl -0pi -e 'my $q=quotemeta($ENV{N2}); s/$q/$ENV{R2}/' "$F" || echo "PERL_RC=NONZERO"
echo "-- old anchor hits after (want 0):"; grep -c -F "$N2" "$F" || true
echo "-- new anchor hits after (want 1):"; grep -c -F "$R2" "$F" || true
echo "########## after-state checks ##########"
echo "-- the two lines now --"
sed -n '54p' "$F"
sed -n '427,428p' "$F"
echo "-- whole-file max line length (anchor side = 267; no new record wanted) --"
awk '{ if (length > m) m = length } END { print m }' "$F"
echo "-- line numbers of the two edits --"
grep -n 'prov.Detector(scopeID)' "$F"
grep -n 'internal/plugin.DisposalScope deferring' "$F"
echo "-- CR count on disk (this file's working copy is CRLF = 存量, blob side is LF) --"
tr -cd '\r' < "$F" | wc -c
echo "-- diff -U0 changed lines (each must start with //) --"
git diff -U0 -- "$F" | grep -E '^[+-]' | grep -vE '^(\+\+\+|---)' || true
echo "-- non-comment changed lines (want 0) --"
git diff -U0 -- "$F" | grep -E '^[+-]' | grep -vE '^(\+\+\+|---)' | grep -vcE '^[+-][[:space:]]*//' || true
echo "-- numstat --"
git diff --numstat -- "$F"
