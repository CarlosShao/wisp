"""codepoint-scan.py - refuse any banned code point in a delivery file or commit message.

The ban list is the one this repo's D22 scanner enforces, plus the wider set the ticket
texts name (U+00A7 section sign, U+2713 check mark, U+2264/2265/2260 relations, circled
digits, arrows, the math-operator band, emoji, variation selector). A file that mentions
one of these must spell it by NAME ("U+00A7"), never as the character.

Usage: python codepoint-scan.py <file> [<file> ...]
Exit:  0 = clean, 1 = at least one hit (each hit printed with line, codepoint, name)
"""
import sys
import unicodedata


def banned(cp):
    return (
        cp in (0x00A7, 0x2713, 0x2264, 0x2265, 0x2260, 0xFE0F)
        or 0x2460 <= cp <= 0x24FF  # circled digits
        or 0x2190 <= cp <= 0x21FF  # arrows
        or 0x2200 <= cp <= 0x22FF  # math operators
        or 0x2600 <= cp <= 0x27BF  # misc symbols / dingbats
        or 0x2B00 <= cp <= 0x2BFF
        or 0x1F000 <= cp <= 0x1FAFF  # emoji + the wider band the scanner uses
        or 0x1F1E6 <= cp <= 0x1F1FF  # flags
    )


total = 0
for path in sys.argv[1:]:
    hits = []
    with open(path, encoding="utf-8") as fh:
        for i, line in enumerate(fh, 1):
            for ch in line:
                cp = ord(ch)
                if banned(cp):
                    hits.append((i, "U+%04X" % cp, unicodedata.name(ch, "?")))
    total += len(hits)
    print("%-60s hits=%d" % (path, len(hits)))
    for h in hits[:20]:
        print("   line %d  %s  %s" % h)

print("codepoint-scan: total hits = %d (want 0)" % total)
sys.exit(0 if total == 0 else 1)
