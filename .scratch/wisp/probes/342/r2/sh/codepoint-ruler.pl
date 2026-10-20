#!/bin/perl
# 342-r2 codepoint ruler: counts U+00A7 (SECTION SIGN) occurrences per file, with every file also
# checked for the neighbouring middle dot U+00B7 so the two are never confused (dispatch's warning).
use strict; use warnings; use utf8;
binmode(STDOUT, ':encoding(UTF-8)');
for my $f (@ARGV) {
    open(my $in, '<:encoding(UTF-8)', $f) or die "open $f: $!";
    my ($a7, $b7, $lines) = (0, 0, 0);
    while (my $l = <$in>) {
        $lines++;
        $a7 += () = $l =~ /\x{00A7}/g;
        $b7 += () = $l =~ /\x{00B7}/g;
    }
    close($in);
    printf("%s lines=%d U+00A7=%d U+00B7=%d\n", $f, $lines, $a7, $b7);
}
