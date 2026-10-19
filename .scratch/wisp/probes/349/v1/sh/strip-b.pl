#!/usr/bin/perl
# V1-B 结构尺：整文件字符状态机（code / "str" / 'rune' / `raw` / //line / /*block*/）
# 只输出 code 态字节；随后删所有空白；sha256 由外层算。
# 与编排者那把 sed 's://.*::' 不同形：本尺不剥字符串字面量里的 //。
use strict; use warnings;
local $/; my $s = <STDIN>;
my $out = ''; my $n = length $s; my $i = 0;
while ($i < $n) {
  my $c = substr($s,$i,1);
  my $d = ($i+1 < $n) ? substr($s,$i+1,1) : '';
  if ($c eq '/' && $d eq '/') { while ($i < $n && substr($s,$i,1) ne "\n") { $i++ } next; }
  if ($c eq '/' && $d eq '*') { $i += 2; while ($i < $n) { if (substr($s,$i,1) eq '*' && substr($s,$i+1,1) eq '/') { $i += 2; last } $i++ } next; }
  if ($c eq '"') { $out .= $c; $i++;
    while ($i < $n) { my $x = substr($s,$i,1); $out .= $x; if ($x eq "\\") { $out .= substr($s,$i+1,1); $i += 2; next } $i++; last if $x eq '"'; }
    next; }
  if ($c eq "'") { $out .= $c; $i++;
    while ($i < $n) { my $x = substr($s,$i,1); $out .= $x; if ($x eq "\\") { $out .= substr($s,$i+1,1); $i += 2; next } $i++; last if $x eq "'"; }
    next; }
  if ($c eq '`') { $out .= $c; $i++;
    while ($i < $n) { my $x = substr($s,$i,1); $out .= $x; $i++; last if $x eq '`'; }
    next; }
  $out .= $c; $i++;
}
$out =~ s/\s+//g;
print $out;
