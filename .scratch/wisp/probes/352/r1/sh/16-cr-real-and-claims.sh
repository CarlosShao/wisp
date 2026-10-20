#!/bin/bash
# 尺⑯：① 六枚 raw 件的 CR 真值（利尺＝tr，钝尺＝grep -c $'\r'／awk /\r/，已自门坐实）
#       ② 逐枚笔的按类枚数（readings⇄logs⇄sh⇄票）⇄该笔正文自述枚数 对账
cd '/d/work/workspace/projects plans/Wisp' || exit 9
P=.scratch/wisp/probes/352/r1
echo '===== A) CR 三把尺对撞（利尺 tr ⇄两把钝尺）====='
for f in ac2-col-s1-treeB ac2-col-s2-treeL ac2-host-fiv-s1-treeB ac2-host-fiv-s2-treeL ac2-host-pair-s1-treeB ac2-host-pair-s2-treeL; do
  fp="$P/logs/$f.raw.txt"
  tr_=$(tr -cd '\r' < "$fp" | wc -c)
  awk_=$(perl -0777 -ne 'my $n = () = /\r/g; print $n' "$fp")
  lines=$(wc -l < "$fp")
  printf '%s CRBYTES_tr=%s CRPCT_perl=%s LINES=%s MKLINK_CRLF_LINES=%s\n' "$f" "$tr_" "$awk_" "$lines" "$(grep -c 'Successfully processed 1 files' "$fp" || true)"
done
echo '  （正控＝造一枚含 2 个 CR 的件：tr 报 2、perl 报 2、grep -c 报 3＝行计数⇒grep 钝）'
D=/d/tmp/wisp/20261020/352r1cont/tmp
mkdir -p "$D"
printf 'a\r\nb\r\nc\n' > "$D/ctl-crlf.txt"
CTR=$(tr -cd '\r' < "$D/ctl-crlf.txt" | wc -c)
CPR=$(perl -0777 -ne '$n = () = /\r/g; print $n' "$D/ctl-crlf.txt")
CGR=$(grep -c $'\r' "$D/ctl-crlf.txt" || true)
CLN=$(wc -l < "$D/ctl-crlf.txt")
echo "  CONTROL tr=$CTR perl=$CPR grep=$CGR lines=$CLN"
echo '===== B) 逐枚笔按类枚数⇄正文自述 ====='
for c in 0fb399c3 ecffb93b aacde632 c16ff47e 06aa2abf 3ecdea30 5c9402f8 1c787282 8a10a7e3; do
  git show --name-only --format='' "$c" | grep '[^[:space:]]' > "$D/f-$c.txt"
  tot=$(grep -c '' "$D/f-$c.txt")
  rd=$(grep -c 'readings-.*\.md' "$D/f-$c.txt" || true)
  lg=$(grep -c '/logs/' "$D/f-$c.txt" || true)
  sh=$(grep -c '/sh/' "$D/f-$c.txt" || true)
  pk=$(grep -c 'issues/352-' "$D/f-$c.txt" || true)
  echo "COMMIT $c TOTAL=$tot READINGS=$rd LOGS=$lg SH=$sh TICKET=$pk"
  git show -s --format='%b' "$c" | grep -oE '枚数＝本笔[^。]{0,80}|本笔 [0-9]+ 枚文件|1 枚报告\+[0-9]+ 枚读数散件|[0-9]+ 枚净化读数' | sed 's/^/    CLAIM: /' | cat
done
echo '===== C) 票尾 numstat 逐枚（真值⇄readings-ac3 第 45 行那句）====='
git show --numstat --format='' 1c787282 | cat
sed -n '45p' "$P/readings-ac3-discipline.md"
echo '===== D) 名册枚数随 HEAD 漂（笔⑨正文那两把尺的数 vs 现量）====='
echo "  NOW_ANCHORED=$(git log --grep='^票 352 落地腿 352-r1' --format='%h' | grep -c . || true)"
sed -n 's/^.*现量逐字：\(ANCHORED_COUNT=8\).*/  CLAIM_IN_BODY: \1/p' /dev/null
git show -s --format='%b' 8a10a7e3 | grep -oE 'ANCHORED_COUNT=8|INMEMORY_ANCHOR0_COUNT=8|UNPUSHED_COMMITS_MINE=8' | sed 's/^/  CLAIM: /' | cat
echo '===== E) probes/352/r1 之下工作树⇄HEAD 差异（除我那两枚新件外须 0）====='
git status --porcelain -- "$P" | cat
echo DONE
