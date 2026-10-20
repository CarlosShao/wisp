#!/usr/bin/perl
# 352-r1 UTF-8 合法性尺（＝⛔ iconv 不在这台 shell 里：命令找不到时 BAD_UTF8 是假警，本尺⛔ 那味）
use strict; use warnings;
use Encode ();
for my $f (@ARGV) {
  open(my $h, "<:raw", $f) or do { print "OPENFAIL $f\n"; next };
  local $/; my $raw = <$h>; close $h;
  # ★Encode::decode 会**就地吃掉**入参 ⇒ 必须先拷贝再解，否则解完 length/CR 全＝0＝假读数（本腿自抓：
  #   第一把尺就报出过"VALID CR=0 BYTES=0"这种六枚全 0 的怪相）。
  my $copy = $raw;
  my $ok = eval { Encode::decode("UTF-8", $copy, Encode::FB_CROAK); 1 };
  my $bytes = length($raw);
  my $cr = () = $raw =~ /\r/g;
  printf "%-60s UTF8=%s CR=%d BYTES=%d\n", $f, ($ok ? "VALID" : "INVALID"), $cr, $bytes;
}
