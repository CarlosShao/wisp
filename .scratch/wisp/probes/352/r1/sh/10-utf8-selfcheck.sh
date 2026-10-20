#!/bin/sh
# 352-r1 UTF-8 尺的自校：先造一枚已知非法 UTF-8（按字节截断中文＝公共块 cut -c 那味）当正控，
# 再数我写面那六枚净化后的件。★⛔ iconv：那命令在这台 shell 里没有（命令找不到时的"BAD"是假警）。
set -u
cd "/d/work/workspace/projects plans/Wisp/.scratch/wisp/probes/352/r1/logs" || exit 9
RULER="/d/work/workspace/projects plans/Wisp/.scratch/wisp/probes/352/r1/sh/09-utf8-ruler.pl"
CTRL="/d/tmp/wisp/20261020/352r1/posctrl-bad-utf8.txt"
perl -e 'open(my $o, ">:raw", $ARGV[0]) or die; print $o "\xe4\xb8"; close $o;' "$CTRL"
echo "== 正控（按字节截断的中文＝该报 INVALID＝尺灵）=="
perl "$RULER" "$CTRL"
echo "== 我写面净化后的六枚（该全 UTF8=VALID）=="
echo "★CR 那几枚不是本腿新造字节＝本机件里 cmd /c mklink 自己吐的 CRLF（逐字 8 行 'Successfully processed 1 files;'），"
echo "  照实保留＝原始读数⛔ 净化；★那条 LF 规矩管的是产码与报告件的新字节（本腿 internal/winsec/winsec.go 入库 blob CR＝0）。"
perl "$RULER" ac2-col-s1-treeB.raw.txt ac2-col-s2-treeL.raw.txt ac2-host-fiv-s1-treeB.raw.txt \
  ac2-host-fiv-s2-treeL.raw.txt ac2-host-pair-s1-treeB.raw.txt ac2-host-pair-s2-treeL.raw.txt
