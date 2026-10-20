#!/bin/bash
# 352-r1 续跑（cont）· 收件尺：名册两把对撞 ⇄ 逐枚写面自证 ⇄ blob 分段 ⇄ 票框终量 ⇄ 票尾 numstat
# 尺型＝量的（读 git 对象库＋工作树文件），⛔ 算的、⛔ 采信任何正文里的枚数。
# 射程＝锚定行首「票 352 落地腿 352-r1」的全部笔 ⇄ .scratch/wisp/probes/352/r1 ⇄ 票 352 那一枚文件。
cd '/d/work/workspace/projects plans/Wisp' || exit 9
T=.scratch/wisp/issues/352-wrappath-errno-half-stringified-so-fs-ernotexist-cannot-punch-through-private-route.md
P=.scratch/wisp/probes/352/r1

echo '== R1 名册尺①＝锚定行首 --grep（⛔ 假设其后任何字符）=='
git log --grep='^票 352 落地腿 352-r1' --format='%h|%ad|%s' --date='format:%H:%M' | cat
echo "ANCHORED_COUNT=$(git log --grep='^票 352 落地腿 352-r1' --format='%h' | grep -c . || true)"
echo "NON_ANCHORED_COUNT=$(git log --grep='票 352 落地腿 352-r1' --format='%h' | grep -c . || true)"

echo '== R2 名册尺②＝单次 %H%x09%s 导出＋内存内行首锚（独立实现，⛔ 同 --grep 机制）=='
git log --format='%H%x09%s' > /d/tmp/wisp/20261020/352r1cont/tmp/head-subjects.txt
printf '%s' '票 352 落地腿 352-r1' > /d/tmp/wisp/20261020/352r1cont/tmp/needle.txt
perl '/d/tmp/wisp/20261020/352r1cont/sh/12-roster-inmemory.pl' \
  /d/tmp/wisp/20261020/352r1cont/tmp/head-subjects.txt /d/tmp/wisp/20261020/352r1cont/tmp/needle.txt | cat
echo '  （★正控＝同尺换 needle，证明这把尺咬得动别的名而不⛔ 本名）'
printf '%s' 'A1063' > /d/tmp/wisp/20261020/352r1cont/tmp/needle-ctl.txt
echo '  CONTROL_NEEDLE_A1063:'
perl '/d/tmp/wisp/20261020/352r1cont/sh/12-roster-inmemory.pl' \
  /d/tmp/wisp/20261020/352r1cont/tmp/head-subjects.txt /d/tmp/wisp/20261020/352r1cont/tmp/needle-ctl.txt | head -3 | sed 's/^/    /'
printf '%s' '票 352 落地腿 352-r1X' > /d/tmp/wisp/20261020/352r1cont/tmp/needle-neg.txt
echo "  NEGATIVE_NEEDLE_COUNT=$(perl '/d/tmp/wisp/20261020/352r1cont/sh/12-roster-inmemory.pl' /d/tmp/wisp/20261020/352r1cont/tmp/head-subjects.txt /d/tmp/wisp/20261020/352r1cont/tmp/needle-neg.txt | sed -n '2p')"

echo '== R3 逐枚 name-only ⇄ 枚数 ⇄ winsec.go 命中 ⇄ 越界 =='
for c in $(git log --grep='^票 352 落地腿 352-r1' --format='%h' | tac); do
  files=$(git show --name-only --format='' "$c" | grep -c '[^[:space:]]' || true)
  go=$(git show --name-only --format='' "$c" | grep -c '\.go$' || true)
  ws=$(git show --name-only --format='' "$c" | grep -c 'internal/winsec/winsec\.go' || true)
  tk=$(git show --name-only --format='' "$c" | grep -c 'issues/352-' || true)
  echo "$c FILES=$files GO=$go WINSEC_GO=$ws TICKET352=$tk"
done
echo '== R4 全部笔路径并集 ⇄ 允许写面之外枚数 ⇄ 产码 .go 枚数 =='
for c in $(git log --grep='^票 352 落地腿 352-r1' --format='%h'); do git show --name-only --format='' "$c"; done \
  | grep '[^[:space:]]' | sort -u > /d/tmp/wisp/20261020/352r1cont/tmp/paths-union.txt
echo "PATHS_TOTAL=$(grep -c '' /d/tmp/wisp/20261020/352r1cont/tmp/paths-union.txt || true)"
echo "OUTSIDE_ALLOWED_WRITEFACE=$(grep -vc '^\.scratch/wisp/probes/352/r1/\|^internal/winsec/winsec\.go$\|^\.scratch/wisp/issues/352-' /d/tmp/wisp/20261020/352r1cont/tmp/paths-union.txt || true)"
echo "GO_FILES_TOUCHED=$(grep -c '\.go$' /d/tmp/wisp/20261020/352r1cont/tmp/paths-union.txt || true)"
for pat in pinnedleaf_342_windows.go winsec_other.go winsec_windows.go internal/risk internal/config cmd/wisp internal/agent internal/panel frontend/ design/ go.mod go.sum .gitignore ci.yml; do
  echo "  FORBIDDEN[$pat]=$(grep -c -- "$pat" /d/tmp/wisp/20261020/352r1cont/tmp/paths-union.txt || true)"
done

echo '== R5 winsec.go blob 逐笔分段（⛔ 区间尺）=='
for c in $(git log --grep='^票 352 落地腿 352-r1' --format='%h' | tac); do
  echo "  $c $(git rev-parse "$c:internal/winsec/winsec.go")"
done
echo "  START_ANCHOR_bf509bb8 $(git rev-parse bf509bb8:internal/winsec/winsec.go)"
echo "  NOW_HEAD $(git rev-parse HEAD) HEAD_BLOB $(git rev-parse HEAD:internal/winsec/winsec.go) WORKTREE_BLOB $(git hash-object internal/winsec/winsec.go)"

echo '== R6 那一行那笔的 -U0 全貌 ⇄ numstat ⇄ 可执行变动行数 =='
git show -U0 --format='' ecffb93b | cat
git show --numstat --format='' ecffb93b | cat
echo "  EXEC_LINES_CHANGED=2 (1 增 1 删＝同一枚 return 语句的格式串；函数签名行⛔ 在 diff 里)"

echo '== R7 票 352 框终量（行首锚＋自校）＋正控票 348 =='
CHK=$(grep -cE '^- \[x\] ' "$T" || true); UCHK=$(grep -cE '^- \[ \] ' "$T" || true); ANY=$(grep -cE '^- \[(x| )\] ' "$T" || true)
echo "  BOX_CHECKED=$CHK BOX_UNCHECKED=$UCHK BOX_ANY=$ANY SELFCHECK_SUM=$((CHK+UCHK))"
T348=.scratch/wisp/issues/348-refusal-path-remove-follows-an-ancestor-link-and-deletes-a-foreign-file.md
echo "  CONTROL_348_CHECKED=$(grep -cE '^- \[x\] ' "$T348" || true) CONTROL_348_UNCHECKED=$(grep -cE '^- \[ \] ' "$T348" || true) CONTROL_348_ANY=$(grep -cE '^- \[(x| )\] ' "$T348" || true)"
echo "  MY_TAIL_BULLETS=$(grep -c '^- \[2026-10-20 16:5x +08\] agent=`352-r1`（票 352' "$T" || true)"
echo "  TICKET_LINES=$(grep -c '' "$T")"

echo '== R8 票尾 numstat（须纯追加＝0 deletions）⇄ 票名与 Status 对回起手锚 =='
git show --numstat --format='' 1c787282 -- "$T" | cat
echo "  MINUS_CONTENT_LINES=$(git show -U0 --format='' 1c787282 -- "$T" | grep -cE '^-[^-]' || true)"
git show bf509bb8:"$T" | sed -n '1p;22p' > /d/tmp/wisp/20261020/352r1cont/tmp/tk-then.txt
sed -n '1p;22p' "$T" > /d/tmp/wisp/20261020/352r1cont/tmp/tk-now.txt
if diff -q /d/tmp/wisp/20261020/352r1cont/tmp/tk-then.txt /d/tmp/wisp/20261020/352r1cont/tmp/tk-now.txt >/dev/null; then
  echo '  TICKET_NAME_AND_STATUS_UNCHANGED=YES'
else
  echo '  TICKET_NAME_AND_STATUS_UNCHANGED=NO'
fi

echo '== R9 写面纯度：0 字节 ⇄ .out ⇄ U+00A7 ⇄ winsec.go 入库 CR =='
git ls-tree -r -l HEAD -- "$P" | awk '$4==0{print $5}' > /d/tmp/wisp/20261020/352r1cont/tmp/zero.txt
echo "  ZERO_BYTE_FILES=$(grep -c '' /d/tmp/wisp/20261020/352r1cont/tmp/zero.txt || true)"
echo "  TRACKED_IN_PROBES=$(git ls-tree -r --name-only HEAD -- "$P" | grep -c . || true)"
echo "  OUT_EXT_FILES=$(git ls-tree -r --name-only HEAD -- "$P" | grep -c '\.out$' || true)"
perl -CSD -e 'my $n=0; for my $f (@ARGV) { open(my $H,"<:encoding(UTF-8)",$f) or next; while(<$H>){ $n += () = /\x{A7}/g } close $H } print "  MY_PROBE_TAIL_U00A7_IN_HEAD=$n\n"' $(git ls-tree -r --name-only HEAD -- "$P" | sed 's|^|/d/work/workspace/projects plans/Wisp/|')
echo "  POSITIVE_CONTROL_HANDOVER_U00A7=$(perl -CSD -ne 'BEGIN{$n=0} $n += () = /\x{A7}/g; END{print $n}' docs/reports/HANDOVER.md)"
echo "  WINSEC_GO_BLOB_CRBYTES=$(git cat-file blob HEAD:internal/winsec/winsec.go | tr -cd '\r' | wc -c)"
echo "  WINSEC_GO_EOL_CANONICAL=$(git ls-files --eol -- internal/winsec/winsec.go | awk '{print $1"  "$2}')"

echo '== R10 ⛔ push 自证 ⇄ 本腿笔未推枚数 =='
echo "  UNPUSHED_TOTAL=$(git rev-list --count origin/dev..HEAD 2>/dev/null || echo NO_ORIGIN_DEV)"
echo "  UNPUSHED_MINE=$(git rev-list --count --grep='^票 352 落地腿 352-r1' origin/dev..HEAD 2>/dev/null || true)"
echo '== R11 internal/winsec 写面活体尺（续跑起手 90 分钟窗内落笔枚数）=='
echo "  RECENT90=$(git log --since='90 minutes ago' --format='%h' -- internal/winsec | grep -c . || true)"
echo DONE
