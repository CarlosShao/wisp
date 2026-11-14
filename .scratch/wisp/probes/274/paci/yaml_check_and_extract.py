"""Ticket 274 paci probe: (a) is ci.yml still valid YAML, (b) does the step I
committed sit where I say it does, (c) extract its `run:` block byte for byte so
the very bytes that get committed are the bytes the control runs."""

import sys
import pathlib
import yaml

REPO = pathlib.Path(r"D:\work\workspace\projects plans\Wisp")
WF = REPO / ".github" / "workflows" / "ci.yml"
OUT = REPO / ".scratch" / "wisp" / "probes" / "274" / "paci" / "extracted-run.ps1"

text = WF.read_text(encoding="utf-8")
try:
    doc = yaml.safe_load(text)
except Exception as exc:  # noqa: BLE001 - the rc is the verdict
    print(f"YAML_PARSE_FAILED: {exc}")
    sys.exit(2)

print("YAML_PARSE_OK: yaml.safe_load returned", type(doc).__name__)
jobs = doc["jobs"]
print("JOBS:", ", ".join(jobs))

steps = jobs["slo-smoke"]["steps"]
names = [s.get("name", s.get("uses", "<unnamed>")) for s in steps]
print("SLO_SMOKE_STEPS:")
for i, n in enumerate(names):
    print(f"  {i}: {n}")

idx = next(i for i, n in enumerate(names) if n.startswith("Shipped wisp.exe"))
step = steps[idx]
print(f"STEP_INDEX={idx} (0-based) of {len(steps)} step(s)")
print(f"PREVIOUS_STEP={names[idx - 1]!r}")
print(f"NEXT_STEP={names[idx + 1]!r}")
print("STEP_KEYS:", sorted(step))
print("STEP_IF:", step.get("if"))
print("STEP_SHELL:", step.get("shell"))
run = step["run"]
print("RUN_BLOCK_LINES:", len(run.splitlines()), "chars:", len(run))
OUT.write_text(run, encoding="utf-8", newline="")
print("WROTE:", OUT)
sys.exit(0)
