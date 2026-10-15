"""guard_strip.py - delete ONE branch of internal/tools/subagent_197.go's assembly guard.

The branch is located by its condition text, never by a line number, because this
leg's own edit shifts the file. Exactly one match is required; anything else is a
refusal to mutate (a mutant that quietly deleted two branches, or none, would make
AC#1's reading meaningless). The call site's own check added by form jia is
matched by a DIFFERENT text and is deliberately left in place here.

Usage: python guard_strip.py <path-to-subject>
"""

import io
import sys

CONDITION = "if t.d.BaseOptions == nil || t.d.ParentTools == nil {"


def main() -> int:
    if len(sys.argv) != 2:
        print("usage: python guard_strip.py <subject>", file=sys.stderr)
        return 2
    path = sys.argv[1]
    with io.open(path, encoding="utf-8") as handle:
        lines = handle.read().splitlines(keepends=True)

    hits = [i for i, line in enumerate(lines) if CONDITION in line]
    if len(hits) != 1:
        print("guard_strip: found %d matches for the guard condition, want 1" % len(hits),
              file=sys.stderr)
        return 3
    start = hits[0]
    body = lines[start + 1]
    brace = lines[start + 2]
    if "宿主没有给出派生用的装配" not in body or brace.strip() != "}":
        print("guard_strip: the line after the guard is not the refusal this leg knows: "
              "%r / %r" % (body, brace), file=sys.stderr)
        return 4
    del lines[start:start + 3]
    with io.open(path, "w", encoding="utf-8", newline="") as handle:
        handle.write("".join(lines))
    print("guard_strip: removed the guard branch at line %d (3 lines)" % (start + 1))
    return 0


if __name__ == "__main__":
    sys.exit(main())
