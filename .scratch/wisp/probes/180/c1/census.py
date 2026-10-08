#!/usr/bin/env python3
"""180-c1 AC#4 census. Read-only: parses internal/config/schema.go, greps for readers.

Two rulers per field, because they answer different questions:
  hits_all          = the ticket's own ruler, verbatim scope internal/ cmd/ minus _test.go
  hits_outside_cfg  = same, but dropping internal/config/** (the schema's own package:
                      defaults.go's reflection and manager.go's block copies hit there
                      without any component acting on the value)
"""
import re, subprocess, sys, os

REPO = r"D:/work/workspace/projects plans/Wisp"
SCHEMA = os.path.join(REPO, "internal", "config", "schema.go")

struct_re = re.compile(r"^type\s+(\w+)\s+struct\b")
field_re = re.compile(r"^\s+(\w+)\s+(\S+)\s+`([^`]*)`")
default_re = re.compile(r'default:"([^"]*)"')
toml_re = re.compile(r'toml:"([^"]*)"')

rows = []
cur_struct = "<none>"
in_struct = False
with open(SCHEMA, encoding="utf-8") as fh:
    for lineno, line in enumerate(fh, 1):
        m = struct_re.match(line)
        if m:
            cur_struct = m.group(1); in_struct = True; continue
        if in_struct and re.match(r"^\}", line):
            in_struct = False; continue
        if "default:" not in line:
            continue
        fm = field_re.match(line)
        if not fm:
            rows.append((lineno, cur_struct, "<UNPARSED>", "", "", line.strip(), "", ""))
            continue
        name, ftype, tag = fm.group(1), fm.group(2), fm.group(3)
        d = default_re.search(tag)
        t = toml_re.search(tag)
        rows.append((lineno, cur_struct, name, ftype,
                     d.group(1) if d else "", t.group(1) if t else "", "", ""))

default_lines = sum(1 for _ in open(SCHEMA, encoding="utf-8") if "default:" in _.replace("\t", " "))
raw_default_lines = 0
with open(SCHEMA, encoding="utf-8") as fh:
    for line in fh:
        if "default:" in line:
            raw_default_lines += 1
with open(SCHEMA, encoding="utf-8") as fh:
    all_lines = fh.readlines()
total_tags = sum(l.count('default:"') for l in all_lines)
comment_default = sum(1 for l in all_lines if "default:" in l and l.strip().startswith("//"))
multi_tag = sum(1 for l in all_lines if l.count('default:"') > 1)

def grep(pattern):
    """Run grep and decode as UTF-8 with replacement.

    Windows' default cp936/GBK codec raises UnicodeDecodeError on the Chinese
    comment lines this repo's Go files contain, which silently turned 17 fields
    into fake zeros on the first run. bytes + explicit decode is the fix.
    """
    try:
        p = subprocess.run(["grep", "-rn", pattern, "--include=*.go", "internal/", "cmd/"],
                           cwd=REPO, capture_output=True, timeout=180)
        txt = p.stdout.decode("utf-8", "replace")
        out = [l for l in txt.splitlines() if "_test.go" not in l]
        return out, p.returncode
    except Exception as e:
        sys.stderr.write("grep fail %s %s\n" % (pattern, e)); return [], -1

out = []
recorded = []
out.append("field\tstruct\ttoml_key\tdefault\tkind\thits_all\thits_outside_cfg\thits_attributed\tsample_all\tsample_outside_cfg\tcollision_flag")
for lineno, st, name, ftype, dflt, toml, _a, _b in rows:
    if name == "<UNPARSED>":
        out.append("UNPARSED\t%s\tschema.go:%d\t\t\t-1\t-1\t\t\t" % (st, lineno)); continue
    pat = r"\." + name + r"\b"
    hits_all, grc = grep(pat)
    # Ruler C (attribution): the hit LINE must also name the owning type, so a hit on
    # some other struct's identically-spelled field does not count. e.g. Ball.Size
    # matches "cfg.Ball.Size" but not "opts.Size" / "s.Size".
    token = re.sub(r"(Section|Config)$", "", st) or st
    a_pat = r"\b" + re.escape(token) + r"\." + name + r"\b"
    hits_att, grc2 = grep(a_pat) if token != name else (hits_all, grc)
    # keep it a strict subset of ruler B: a real production reader is outside the
    # schema package AND names the owning type on the same line
    hits_att = [h for h in hits_att if not h.startswith("internal/config/")]
    if grc == -1 or grc2 == -1:
        # a crashed grep is NOT a zero; mark it so no summary counts it as unread
        recorded.append((name, st, -2, -2))
        out.append("\t".join([name, st, toml, dflt, ftype, "GREP_FAIL", "GREP_FAIL", "GREP_FAIL", "", "", ""]))
        continue
    hits_out = [h for h in hits_all if not h.startswith("internal/config/")]
    # collision flag: hit whose receiver is a different struct's field of same name
    collide = ""
    if hits_all:
        others = [h for h in hits_all if (name + " ") not in h and "." + name not in h]
        collide = "generic-name" if name in ("Width", "Height", "Mode", "Depth", "Level",
                                            "Size", "Type", "Name", "Enabled", "Scale",
                                            "Speed", "Threshold", "Path", "Format",
                                            "Timeout", "Interval", "Model", "Voice",
                                            "Provider", "Role", "Rate", "Device") else ""
    recorded.append((name, st, len(hits_all), len(hits_out)))
    out.append("\t".join([
        name, st, toml, dflt, ftype,
        str(len(hits_all)), str(len(hits_out)), str(len(hits_att)),
        (hits_all[0].split(":")[0] + ":" + hits_all[0].split(":")[1]) if hits_all else "",
        (hits_out[0].split(":")[0] + ":" + hits_out[0].split(":")[1]) if hits_out else "",
        collide,
    ]))

dest = os.path.join(REPO, ".scratch", "wisp", "probes", "180", "c1", "census.tsv")
with open(dest, "w", encoding="utf-8", newline="\n") as fh:
    fh.write("\n".join(out) + "\n")

print("default_lines_in_file=%d" % raw_default_lines)
print("default_tags_total=%d" % total_tags)
print("default_in_comment_lines=%d" % comment_default)
print("lines_with_two_plus_default_tags=%d" % multi_tag)
print("parsed_field_rows=%d" % len(rows))
print("unparsed=%d" % sum(1 for r in rows if r[2] == "<UNPARSED>"))
print("distinct_field_names=%d" % len(set(r[2] for r in rows if r[2] != "<UNPARSED>")))
print("zero_hits_all=%d" % sum(1 for r in recorded if r[2] == 0))
print("zero_hits_outside_cfg=%d" % sum(1 for r in recorded if r[3] == 0))
print("nonzero_hits_all=%d" % sum(1 for r in recorded if r[2] > 0))
print("grep_fail=%d" % sum(1 for r in recorded if r[2] == -2))
print("generic_name_flagged=%d" % sum(1 for r in recorded if r[0] in
      ("Width", "Height", "Mode", "Depth", "Level", "Size", "Type", "Name", "Enabled",
       "Scale", "Speed", "Threshold", "Path", "Format", "Timeout", "Interval", "Model",
       "Voice", "Provider", "Role", "Rate", "Device")))
print("rc=0")
