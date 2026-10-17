#!/usr/bin/env bash
# AC#6 arm A present-state reading: for each of the three ALLOW-side functions,
# take the body window BY FUNCTION NAME out of the anchor tree and count the
# five substring operators. Independent gauge from the Go pin: this one is a
# shell/python ruler, so the two readings can be compared.
# Usage: bash ac6a-present-state.sh <tree-root>
set -u
T="${1:?need exported tree root}"
python - "$T" <<'PY'
import io, re, sys
t = sys.argv[1]
targets = [("rules_shell.go", "shellAllowlisted"),
           ("rules_network.go", "hostAllowlisted"),
           ("blacklist.go", "overrideApplies")]
ops = ("Contains", "ContainsAny", "ContainsRune", "Index", "IndexAny")
for fname, fn in targets:
    p = t + "/internal/risk/" + fname
    lines = io.open(p, encoding="utf-8").read().split("\n")
    start = None
    for i, ln in enumerate(lines):
        if ln.startswith("func " + fn + "("):
            start = i
            break
    if start is None:
        print("%s:%s FUNCTION-NOT-FOUND" % (fname, fn))
        continue
    depth, opened, end = 0, False, None
    for i in range(start, len(lines)):
        for ch in lines[i]:
            if ch == "{":
                depth += 1; opened = True
            elif ch == "}":
                depth -= 1
        if opened and depth <= 0:
            end = i
            break
    hits = []
    for i in range(start, end + 1):
        for m in re.finditer(r"strings\.([A-Za-z0-9_]+)", lines[i]):
            if m.group(1) in ops:
                hits.append("%s:%d: strings.%s" % (fname, i + 1, m.group(1)))
    print("%s:%s window=lines %d..%d (%d) substring_operators=%d %s"
          % (fname, fn, start + 1, end + 1, end - start + 1, len(hits), hits))
PY
