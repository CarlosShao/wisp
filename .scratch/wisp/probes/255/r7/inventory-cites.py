#!/usr/bin/env python3
"""Ticket 255-r2 probe: inventory every path:LINE([token]) reference in the roster
and test files, and measure token uniqueness in the referenced file."""
import re
import os
import sys

ROOT = r"D:\work\workspace\projects plans\Wisp"

CITE = re.compile(r"([\w./-]+\.go):(\d+) \[([^\]]+)\]")
ANYPATH = re.compile(r"([\w./-]+\.go):(\d+)")

FILES = [
    "cmd/wisp/config_readers_255.go",
    "cmd/wisp/config_receipt_255_test.go",
]


def read_lines(rel):
    with open(os.path.join(ROOT, rel), "r", encoding="utf-8") as f:
        return f.read().split("\n")


def main():
    for rel in FILES:
        lines = read_lines(rel)
        print("=" * 78)
        print("FILE", rel, "lines=", len(lines))
        print("-" * 78)
        # full cite form
        seen = set()
        for i, line in enumerate(lines, 1):
            for m in CITE.finditer(line):
                key = (i, m.group(0))
                if key in seen:
                    continue
                seen.add(key)
                target, lineno, token = m.group(1), int(m.group(2)), m.group(3)
                tpath = os.path.join(ROOT, target.replace("/", os.sep))
                if not os.path.exists(tpath):
                    print(f"L{i}: CITE {target}:{lineno} [{token[:60]}] -> FILE MISSING")
                    continue
                tl = read_lines(target)
                hits = [j for j, t in enumerate(tl, 1) if token in t]
                rec_has = bool(1 <= lineno <= len(tl) and token in tl[lineno - 1])
                uniq = len(hits) == 1
                tag = "UNIQUE" if uniq else ("NONE" if not hits else f"NONUNIQ({len(hits)})")
                print(f"L{i}: CITE {target}:{lineno} [{token[:70]}] -> hits={hits} recorded_has_token={rec_has} {tag}")
        # bare refs (no [token] after)
        print("-" * 78)
        for i, line in enumerate(lines, 1):
            for m in ANYPATH.finditer(line):
                # skip if part of a full cite match
                if any(m.start() >= c.start() and m.end() <= c.end() for c in CITE.finditer(line)):
                    continue
                print(f"L{i}: BARE {m.group(0)} ... context: {line.strip()[:110]}")


if __name__ == "__main__":
    main()
