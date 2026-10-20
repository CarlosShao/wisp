#!/bin/sh
# 352-r1 入库正文净化：把被禁字面字符（码位 U+00A7）换成其码位名（＝第十五波公共块 shell 节那条，
# 先例＝342-v2 sh/07 同一形）。★⛔ 改任何读数、只换那一枚字符本身；★正控＝同一枚 U+00B7（间隔号）⛔ 被动。
set -u
cd "/d/work/workspace/projects plans/Wisp/.scratch/wisp/probes/352/r1/logs" || exit 9
echo "== 净化前逐枚命中 =="
perl "/d/work/workspace/projects plans/Wisp/.scratch/wisp/probes/352/r1/sh/07-codepoint-ruler.pl" *.raw.txt *.md ../../readings-*.md 2>/dev/null | grep '^TOTAL' | grep -v ' = 0$' || echo "NO_HITS_BEFORE"
for f in *.raw.txt; do
  perl -i -pe 's/\xc2\xa7/U+00A7/g' "$f"
done
echo "== 净化后命中（应全 0）=="
perl "/d/work/workspace/projects plans/Wisp/.scratch/wisp/probes/352/r1/sh/07-codepoint-ruler.pl" *.raw.txt | grep -c '^HIT' | sed 's/^/HITS_AFTER=/'
echo "== 正控：U+00B7（间隔号）在场与否＝同一 perl 尺另数，证那枚没被顺手换 =="
perl -e '
  my $dot = "\xc2\xb7"; my $tot=0;
  for my $f (@ARGV) { open(my $h,"<",$f) or next; local $/; my $s=<$h>; close $h; my $n=()=$s=~/$dot/g; $tot+=$n; }
  print "MIDDLE_DOT_U00B7_TOTAL_IN_LOGS=$tot\n";'  *.raw.txt
echo "== 净化行数＝只那六枚有命中的件被改（其余件字节⛔ 动＝sha 对回）=="
sha256sum ac2-col-s1-treeB.raw.txt ac2-col-s2-treeL.raw.txt ac2-host-fiv-s1-treeB.raw.txt ac2-host-fiv-s2-treeL.raw.txt ac2-host-pair-s1-treeB.raw.txt ac2-host-pair-s2-treeL.raw.txt
echo "== ★diff 纯度尺＝那六枚里每个 hunk 都只改那一枚字符（增删行逐字对偶）=="
for f in ac2-col-s1-treeB.raw.txt ac2-col-s2-treeL.raw.txt ac2-host-fiv-s1-treeB.raw.txt ac2-host-fiv-s2-treeL.raw.txt ac2-host-pair-s1-treeB.raw.txt ac2-host-pair-s2-treeL.raw.txt; do
  d="/d/work/workspace/projects plans/Wisp/.scratch/wisp/probes/352/r1/logs/$f"
  printf '%-34s NUMSTAT=%s\n' "$f" "$(git -C '/d/work/workspace/projects plans/Wisp' diff --numstat -- "$d" | awk '{print $1"增 "$2"删"}')"
  # 每枚删除行把 U+00A7 换成码位名后，应与对应的增加行逐字相同＝证明⛔ 改任何读数
  perl -MEncode -e '
    my ($f)=@ARGV; open(my $h,"-|:raw","git","-C","/d/work/workspace/projects plans/Wisp","diff","-U0","--",$f) or die;
    local $/; my $t=<$h>; close $h; my @m = $t =~ /^\+[^+].*$/gm; my @n = $t =~ /^-[^-].*$/gm;
    my $same=0; for my $x (@n){ my $y=$x; $y=~s/^-/+/; $y=~s/\xc2\xa7/U+00A7/g; $same++ if grep {$_ eq $y} @m; }
    printf("   PAIRED_ROWS=%d/%d  (★每枚删除行换码位名后逐字＝对应增加行 ⇒ 只换那一枚字符、⛔ 动读数)\n", $same, scalar(@n));' "$d"
done
echo "== UTF-8 合法性＋CR 尺（⛔ iconv：那命令在这台 shell 里没有，命令找不到时的"BAD_UTF8"是假警＝公共块"报 0 命中前先跑正控"同味）=="
perl "/d/work/workspace/projects plans/Wisp/.scratch/wisp/probes/352/r1/sh/09-utf8-ruler.pl" ac2-col-s1-treeB.raw.txt ac2-col-s2-treeL.raw.txt ac2-host-fiv-s1-treeB.raw.txt ac2-host-fiv-s2-treeL.raw.txt ac2-host-pair-s1-treeB.raw.txt ac2-host-pair-s2-treeL.raw.txt
