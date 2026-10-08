#!/usr/bin/env python3
"""111-ciif1 ruler: structural scan of .github/workflows/ci.yml.

Two independent passes, then they must agree:
  PASS 1 (yaml)  : safe_load + walk jobs -> steps, report `if` presence per step.
  PASS 2 (lines) : indent-based scan, skipping comment lines, classifying each
                   `if:` line as job-level vs step-level and counting it.

This is parse-only: nothing from ci.yml is executed here.
"""
import re
import sys
import yaml

PATH = ".github/workflows/ci.yml"

# ---------- PASS 1: yaml ----------
with open(PATH, "r", encoding="utf-8") as fh:
    raw = fh.read()
lines = raw.splitlines()

doc = yaml.safe_load(raw)          # the syntax self-check
jobs = doc["jobs"]

yaml_rows = []
yaml_step_if = 0                   # step-level `if` keys found by the parser
yaml_job_if = 0
for job_id, job in jobs.items():
    if "if" in job:
        yaml_job_if += 1
    for idx, step in enumerate(job.get("steps") or [], start=1):
        name = step.get("name") or step.get("id") or step.get("uses") or step.get("run", "")
        name = str(name).replace("\n", " ")[:70]
        has = "if" in step
        if has:
            yaml_step_if += 1
        yaml_rows.append((job_id, idx, name, has, step.get("if", "")))

# ---------- PASS 2: line/indent scan ----------
top_key_re = re.compile(r"^([A-Za-z_][A-Za-z0-9_.-]*):\s*(#.*)?$")
step_if_lines = []
job_if_lines = []
comment_if_lines = []
other_if_lines = []
in_jobs = False
cur_job = None
cur_step = None
step_start_re = re.compile(r"^\s*-\s+(\S[^:]*):\s*(.*)$")

for ln, text in enumerate(lines, start=1):
    if not text.strip():
        continue
    stripped = text.lstrip(" ")
    indent = len(text) - len(stripped)

    if indent == 0:
        in_jobs = bool(top_key_re.match(text) and text.startswith("jobs:"))
        cur_job = None
        cur_step = None
        continue
    if not in_jobs:
        continue
    if stripped.startswith("#"):
        if re.match(r"^#\s*(-\s+)?if:", stripped):
            comment_if_lines.append((ln, indent, stripped.strip()[:70]))
        continue
    if indent == 2:
        m = re.match(r"^([A-Za-z0-9_.-]+):\s*$", stripped)
        if m:
            cur_job = m.group(1)
            cur_step = None
        continue
    if indent == 4 and re.match(r"^if:\s*\S", stripped):
        job_if_lines.append((ln, cur_job, stripped.strip()[:70]))
        continue
    m = step_start_re.match(text)
    if m and indent >= 6:
        cur_step = f"{m.group(1)}: {m.group(2)}"[:60]
    if re.match(r"^if:\s*\S", stripped) and indent >= 6:
        step_if_lines.append((ln, cur_job, cur_step, stripped.strip()[:70]))
    elif re.match(r"^-\s+if:\s*\S", stripped):
        step_if_lines.append((ln, cur_job, "(step starts with if)", stripped.strip()[:70]))

print("== PASS 1 (yaml) ==")
print(f"jobs: {len(jobs)}   steps: {len(yaml_rows)}")
print(f"job-level `if` keys   : {yaml_job_if}")
print(f"step-level `if` keys  : {yaml_step_if}")
print()
print("== PASS 2 (indent scan) ==")
print(f"step-level `if:` lines: {len(step_if_lines)}")
for ln, job, step, txt in step_if_lines:
    print(f"  {ln:>4}  {job:<22} {txt}")
print(f"job-level `if:` lines : {len(job_if_lines)}")
for ln, job, txt in job_if_lines:
    print(f"  {ln:>4}  {job:<22} {txt}")
print(f"comment-only `if:` hits (NOT steps): {len(comment_if_lines)}")
for ln, ind, txt in comment_if_lines:
    print(f"  {ln:>4}  indent={ind}  {txt}")
print()
print("== AGREE? ==")
ok = (yaml_step_if == len(step_if_lines)) and (yaml_job_if == len(job_if_lines))
print(f"yaml_step_if={yaml_step_if} lines_step_if={len(step_if_lines)} "
      f"yaml_job_if={yaml_job_if} lines_job_if={len(job_if_lines)} -> {ok}")

print()
print("== ROSTER (job / step / step-level if?) ==")
for job_id, idx, name, has, expr in yaml_rows:
    print(f"{job_id}|{idx}|{name}|{'IF' if has else 'NO-IF'}|{expr}")

sys.exit(0)
