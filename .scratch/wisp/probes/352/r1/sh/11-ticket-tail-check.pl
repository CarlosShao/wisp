#!/usr/bin/perl
use strict; use warnings;
for my $f (@ARGV) {
  open(my $h, "<:raw", $f) or do { print "OPENFAIL $f\n"; next };
  local $/; my $s = <$h>; close $h;
  my $p = () = $s =~ /落地腿 352-r1/g;
  my $cr = () = $s =~ /\r/g;
  my $b = () = $s =~ /\xc2\xa7/g;
  my $nl = () = $s =~ /\n/g;
  printf "%s BYTES=%d LINES=%d CR=%d U00A7=%d PREFIX_LUODI_TUI=%d\n", $f, length($s), $nl, $cr, $b, $p;
}
