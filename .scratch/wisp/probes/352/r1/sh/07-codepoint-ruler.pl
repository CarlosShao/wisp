#!/usr/bin/perl
# 352-r1 码位尺：数被禁字面字符（码位 U+00A7）的命中＋给出命中行的可见上下文（换成码位名打印＝⛔ 再写那枚字符本身）。
use strict; use warnings;
my $band = "\x{00A7}";
for my $f (@ARGV) {
  open(my $fh, "<:encoding(UTF-8)", $f) or do { print "OPENFAIL $f: $!\n"; next };
  my $n = 0; my $i = 0;
  while (my $line = <$fh>) { $i++;
    my $c = () = $line =~ /$band/g;
    if ($c) { $n += $c; my $vis = $line; $vis =~ s/$band/<<U00A7>>/g; $vis =~ s/\s+$//;
      print "HIT $f:$i ($c) :: ", substr($vis,0,180), "\n"; }
  }
  close $fh;
  print "TOTAL $f = $n\n";
}
