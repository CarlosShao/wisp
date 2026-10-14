#!/usr/bin/env python3
"""255-v3 acceptance leg: independent re-measurement of ticket 255's roster.

Re-implements (not imports) two things the ticket's instruments do:
  1. every `path.go:LINE [token]` cite inside hotRowClaims -> token occurrence
     count in the named file, and whether the token is a substring of the
     recorded line (AC#6 criterion = token unique hit; :LINE = reading);
  2. the reader scan's `readsFieldAccess` rule for a Config field, walked over
     cmd/ internal/ tools/ scripts/ with the same exclusions the Go scan uses,
     so "production reader set == sectionReadSites set" can be checked without
     trusting the Go test's own error output.

Read-only. Prints one JSON blob.
"""
import json, os, re, subprocess, sys

ROOT = os.path.abspath(os.path.join(os.path.dirname(__file__), "..", "..", "..", "..", ".."))
ROSTER = os.path.join(ROOT, "cmd", "wisp", "config_readers_255.go")
TEST = os.path.join(ROOT, "cmd", "wisp", "config_receipt_255_test.go")

cite_re = re.compile(r"([\w./-]+\.go):(\d+) \[([^\]]+)\]")

def read_lines(rel):
    with open(os.path.join(ROOT, rel), "rb") as f:
        raw = f.read()
    return raw.decode("utf-8", "replace").split("\n")

def roster_verdict_cites():
    """Cites that appear in hotRowClaims map VALUES (the roster ruler's scope)."""
    lines = read_lines("cmd/wisp/config_readers_255.go")
    in_map = False
    out = []
    for i, line in enumerate(lines, 1):
        if line.startswith("var hotRowClaims = map[string]string{"):
            in_map = True
            continue
        if in_map and line.startswith("}"):
            in_map = False
        if not in_map:
            continue
        m = re.match(r'\s*"([^"]+)":\s*(.*)$', line)
        if not m:
            continue
        row, value = m.group(1), m.group(2)
        for c in cite_re.finditer(value):
            out.append({"row": row, "roster_line": i, "file": c.group(1),
                        "line": int(c.group(2)), "token": c.group(3)})
    return out

def check_cites():
    res = []
    for c in roster_verdict_cites():
        lines = read_lines(c["file"])
        hits = [i + 1 for i, l in enumerate(lines) if c["token"] in l]
        recorded = lines[c["line"] - 1].strip() if 1 <= c["line"] <= len(lines) else "<out of range>"
        res.append({
            "row": c["row"], "cite": f'{c["file"]}:{c["line"]}',
            "token": c["token"], "hits": hits, "unique": len(hits) == 1,
            "recorded_line_now": recorded,
            "token_on_recorded_line": c["token"] in recorded,
            "recorded_equals_hit": len(hits) == 1 and hits[0] == c["line"],
        })
    return res

# ---- reader scan replication (scanReadFiles / readsFieldAccess) --------------

def is_ident_byte(ch):
    return ch == "_" or ch.isdigit() or ("a" <= ch <= "z") or ("A" <= ch <= "Z")

def reads_field_access(line, field):
    dotted = "." + field + "."
    i = 0
    while i + len(dotted) < len(line):
        if line[i:i + len(dotted)] != dotted:
            i += 1
            continue
        nxt = line[i + len(dotted)]
        if "A" <= nxt <= "Z":
            return True
        i += 1
    for prefix in ("Config()." + field, "Config." + field):
        idx = line.find(prefix)
        if idx < 0:
            continue
        end = idx + len(prefix)
        if end >= len(line) or not is_ident_byte(line[end]):
            return True
    return False

def scan_read_files(field, skip={"config_readers_255.go"}):
    hits = set()
    for d in ("cmd", "internal", "tools", "scripts"):
        base = os.path.join(ROOT, d)
        for dirpath, dirnames, filenames in os.walk(base):
            if "testdata" in dirnames:
                dirnames.remove("testdata")
            for fn in filenames:
                if not fn.endswith(".go") or fn.endswith("_test.go"):
                    continue
                p = os.path.join(dirpath, fn)
                rel = os.path.relpath(p, ROOT).replace(os.sep, "/")
                if rel.startswith("internal/config/") or fn in skip:
                    continue
                try:
                    lines = read_lines(rel)
                except OSError:
                    continue
                for l in lines:
                    t = l.strip()
                    if t.startswith("//") or t.startswith("*"):
                        continue
                    if reads_field_access(t, field):
                        hits.add(rel)
                        break
    return sorted(hits)

def roster_section_sites():
    lines = read_lines("cmd/wisp/config_readers_255.go")
    in_map = False
    out = {}
    for line in lines:
        if line.startswith("var sectionReadSites = map[string][]string{"):
            in_map = True
            continue
        if in_map and line.startswith("}"):
            break
        if not in_map:
            continue
        m = re.match(r'\s*"([^"]+)":\s*(.*)$', line)
        if not m:
            continue
        val = m.group(2).strip()
        if val == "nil":
            out[m.group(1)] = []
        else:
            out[m.group(1)] = re.findall(r'"([^"]+)"', val)
    return out

def counts():
    with open(ROSTER, "rb") as f:
        roster = f.read().decode()
    with open(TEST, "rb") as f:
        test = f.read().decode()
    r1 = re.findall(r"[A-Za-z0-9_./]*\.go:[0-9]+", roster)
    r2 = re.findall(r":[0-9]+ \[", roster)
    test_refs = [l for l in test.split("\n") if re.search(r"[A-Za-z0-9_./]*\.go:[0-9]+", l)]
    return {
        "ruler1_total": len(r1), "ruler1_unique": len(set(r1)),
        "ruler2_items": len(r2), "ruler2_rows": sum(1 for l in roster.split("\n") if re.search(r":[0-9]+ \[", l)),
        "test_side_lines": len(test_refs),
        "test_side_texts": [l.strip()[:100] for l in test_refs],
    }

out = {
    "cites": check_cites(),
    "counts": counts(),
    "voice_scan_files": scan_read_files("Voice"),
    "section_read_sites": roster_section_sites(),
}
site_map = out["section_read_sites"]
out["voice_scan_matches_roster"] = out["voice_scan_files"] == sorted(site_map.get("Voice", []))

# every section the roster adjudicates: production scan vs sectionReadSites
out["all_sections"] = {
    "field": {f: {"scan": scan_read_files(f), "roster": sorted(v),
                  "match": scan_read_files(f) == sorted(v)}
              for f, v in sorted(site_map.items())}
}

# bare (pre-sharpening) tokens: are they still non-unique now? Proves the
# sharpening is what buys uniqueness, not an unrelated file edit.
bare = {
    "llm": ("cmd/wisp/panel_config_store.go", "cfg.LLM.Roles.Chat.Model"),
    "ball": ("internal/ball/ball_windows.go", "SizePx"),
    "audio": ("cmd/wisp/resident_audio_windows.go", "c.Audio.MicMutedDefault"),
    "privacy": ("internal/observe/logging.go", "RedactPaths"),
    "voice": ("cmd/wisp/resident_audio_windows.go", "c.Voice.Enabled"),
    "agent_comment": ("cmd/wisp/run.go", "cfg.Agent.PerToolTimeoutMS"),
}
out["bare_tokens_now"] = {
    k: {"file": f, "hits": [i + 1 for i, l in enumerate(read_lines(f)) if t in l]}
    for k, (f, t) in bare.items()
}
print(json.dumps(out, indent=1, ensure_ascii=False))
