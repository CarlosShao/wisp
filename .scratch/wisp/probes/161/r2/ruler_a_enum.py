#!/usr/bin/env python3
"""
161-r2 RULER A -- type-level enumeration of "ask the human" verdict producers.

Structurally different from a word grep in three ways:
  1. Go source goes through a comment/string stripper FIRST, so a token that
     only appears inside a comment or a string literal CANNOT be counted as a
     production site. (A word grep counts it.)
  2. Matching is on IDENTIFIERS (word boundaries after stripping), so
     `PendingWindow` does not match `pendingWindowHelper`.
  3. Every hit is classified by syntactic position (call / return / case /
     decl / assign / selector / other) and bucketed by file kind.

This ruler alone makes NO absence claim. Ruler B is the independent shape.

Usage:
  ruler_a_enum.py <root> [--ext .go] [--emit] [--only tok1,tok2]
  --emit : also write ruler_a_emit.json  (consumed by ruler_c_reachability.py)
"""
import json
import os
import re
import sys

SKIP_DIRS = {".git", ".scratch", "node_modules", "third_party", "build",
             "tmp-none", "models", "dist", "vendor", ".gocache"}

# The family whose DEFINITION sites were read off disk by hand (ruler B
# re-derives them from the parser, so a typo here gets caught, not inherited).
FAMILY = {
    "risk.L1":                "risk.Level value: ask via L1 window",
    "risk.L2":                "risk.Level value: ask via L2 card",
    "PendingWindow":          "tools.Gate method: opens the L1 ask",
    "PendingApproval":        "tools.Gate method: opens the L2 ask",
    "Prompt":                 "approval.UI method: puts the question on screen",
    "push":                   "approval.Queue.push: creates a pending item",
    "EvApprovalNeeded":       "statemachine event: the ball's ask doorbell",
    "StateAwaitingApproval":  "statemachine state: shown while asking",
    "StateConfirming":        "statemachine state: shown while asking (L1)",
    "DecideFromNative":       "answer router (NOT an ask) - brief's candidate",
    "DecideFromPanel":        "answer router (NOT an ask) - brief's candidate",
    "AdmitTextTask":          "D47 registration (NOT an ask) - brief's candidate",
}
ORDER = ["risk.L1", "risk.L2", "PendingWindow", "PendingApproval", "Prompt",
         "push", "EvApprovalNeeded", "StateAwaitingApproval", "StateConfirming",
         "DecideFromNative", "DecideFromPanel", "AdmitTextTask"]

IDENT_RE = re.compile(r"\b([A-Za-z_][A-Za-z0-9_]*\.)?([A-Za-z_][A-Za-z0-9_]*)\b")


def strip_go(src):
    """Blank out comments and string/rune literals, keep line/column shape."""
    out = []
    i, n = 0, len(src)
    while i < n:
        c = src[i]
        nxt = src[i + 1] if i + 1 < n else ""
        if c == "/" and nxt == "/":
            j = src.find("\n", i)
            i = n if j < 0 else j
            continue
        if c == "/" and nxt == "*":
            j = src.find("*/", i + 2)
            j = n if j < 0 else j + 2
            out.append("".join(ch if ch == "\n" else " " for ch in src[i:j]))
            i = j
            continue
        if c in "\"'`":
            q, j = c, i + 1
            while j < n:
                if src[j] == "\\" and q != "`":
                    j += 2
                    continue
                if src[j] == q:
                    j += 1
                    break
                j += 1
            else:
                j = n
            out.append("".join(ch if ch == "\n" else " " for ch in src[i:j]))
            i = j
            continue
        out.append(c)
        i += 1
    return "".join(out)


def classify(line, tok):
    s = line.strip()
    if re.search(r"\breturn\b[^;]*\b" + re.escape(tok) + r"\b", line):
        return "return"
    if re.match(r"^(func|const|var|type)\b", s):
        return "decl"
    if re.match(r"^case\b", s):
        return "case"
    if re.search(r"\b" + re.escape(tok) + r"\s*\(", line):
        return "call"
    if re.search(r"\b" + re.escape(tok) + r"\b\s*[=:]", line):
        return "assign"
    return "other"


def bucket(rel):
    fn = os.path.basename(rel)
    if fn.endswith("_test.go"):
        return "go-test"
    if rel.startswith("internal/panel/") or rel.startswith("frontend/"):
        return "panel-or-frontend"
    if rel.startswith("cmd/"):
        return "go-prod-cmd"
    if rel.startswith("internal/"):
        return "go-prod-internal"
    if rel.startswith("tools/"):
        return "go-tooling"
    return "other"


def main():
    args = [a for a in sys.argv[1:]]
    root = args[0]
    ext = ".go"
    emit = "--emit" in args
    only = None
    if "--ext" in args:
        ext = args[args.index("--ext") + 1]
    if "--only" in args:
        only = set(args[args.index("--only") + 1].split(","))
    fixtures = "--fixtures" in args

    hits, scanned = {}, 0
    for dirpath, dirnames, filenames in os.walk(root):
        dirnames[:] = [d for d in dirnames if d not in SKIP_DIRS]
        for fn in sorted(filenames):
            if not fn.endswith(ext):
                continue
            p = os.path.join(dirpath, fn)
            rel = os.path.relpath(p, root).replace("\\", "/")
            try:
                raw = open(p, encoding="utf-8", errors="replace").read()
            except OSError:
                continue
            scanned += 1
            code = strip_go(raw)
            clines, rlines = code.split("\n"), raw.split("\n")
            for i, ln in enumerate(clines):
                if not ln.strip():
                    continue
                found = set()
                for m in IDENT_RE.finditer(ln):
                    sel, name = (m.group(1) or "").rstrip("."), m.group(2)
                    if sel in ("risk", "statemachine") and name in ("L1", "L2"):
                        found.add("risk." + name)
                    if name in FAMILY and (sel or name in ("Prompt", "push")):
                        found.add(name)
                    if name in FAMILY and not sel:
                        found.add(name)
                for tok in found:
                    hits.setdefault(tok, []).append({
                        "file": rel, "line": i + 1, "bucket": bucket(rel),
                        "pos": classify(ln, tok),
                        "text": rlines[i].strip()[:120] if i < len(rlines) else "",
                    })
    toks = only or set(ORDER)
    print("ruler A | root=%s ext=%s files_scanned=%d family=%d" %
          (root, ext, scanned, len(FAMILY)))
    summary = {}
    for tok in ORDER:
        if tok not in toks:
            continue
        rows = hits.get(tok, [])
        c = {}
        for r in rows:
            c[r["bucket"] + "/" + r["pos"]] = c.get(r["bucket"] + "/" + r["pos"], 0) + 1
        prod = [r for r in rows if r["bucket"].startswith("go-prod")
                or r["bucket"] == "panel-or-frontend"]
        summary[tok] = {"total": len(rows),
                        "prod": len(prod),
                        "test": sum(1 for r in rows if r["bucket"] == "go-test"),
                        "by": c}
        print("\n== %-22s total=%-4d prod=%-4d test=%-4d  def=%s" %
              (tok, len(rows), len(prod),
               sum(1 for r in rows if r["bucket"] == "go-test"),
               FAMILY[tok]))
        for k in sorted(c):
            print("     %-40s %d" % (k, c[k]))
        if fixtures:
            for r in prod:
                print("     >> %s:%d [%s] %s" % (r["file"], r["line"], r["pos"], r["text"]))
    if emit:
        with open(os.path.join(os.path.dirname(os.path.abspath(__file__)),
                               "ruler_a_emit.json"), "w", encoding="utf-8") as f:
            json.dump({"root": root, "hits": hits, "summary": summary}, f, indent=1)
        print("\nwrote ruler_a_emit.json")


if __name__ == "__main__":
    main()
