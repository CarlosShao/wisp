#!/usr/bin/env python3
"""313-v1: mechanical three-way key comparison.

Sources (all read-only):
  A. Go payload: lines V1-EMPTY-PAYLOAD / V1-FULL-PAYLOAD in the probe log,
     produced by running config.get through the production assembly.
  B. Spec: the key table in docs/evidence/s1/313-panel-ask.md section 1.
  C. Page: the ConfigView / ConfigViewProvider / ConfigViewVoice interfaces in
     frontend/src/lib/config-view.ts.

Prints the three key sets per level (top / provider row / voice) and every
symmetric difference, so a divergence is named key by key.
"""
import json
import re
import sys

ROOT = "D:/work/workspace/projects plans/Wisp"


def go_keys(log_path):
    empty = full = None
    for line in open(log_path, encoding="utf-8", errors="replace"):
        if line.startswith("V1-EMPTY-PAYLOAD "):
            empty = json.loads(line.split(" ", 1)[1])
        elif line.startswith("V1-FULL-PAYLOAD "):
            full = json.loads(line.split(" ", 1)[1])
    if empty is None or full is None:
        sys.exit("probe log missing payload lines")
    def row_keys(p):
        return sorted(p["providers"][0].keys()) if p.get("providers") else []
    return {
        "top": (sorted(empty.keys()), sorted(full.keys())),
        "row": (row_keys(empty), row_keys(full)),
        "voice": (sorted(empty["voice"].keys()), sorted(full["voice"].keys())),
    }


def spec_keys(path):
    text = open(path, encoding="utf-8").read()
    sec = text.split("## 1.", 1)[1].split("## 2.", 1)[0]
    top, row, voice = [], [], []
    for raw in sec.splitlines():
        if not raw.startswith("|") or raw.startswith("|---") or "键" in raw.split("|")[1]:
            continue
        # split on real cell separators only: the table escapes literal pipes
        # inside cells as \| (e.g. "string[] \| null")
        cells = [c.strip() for c in re.split(r"(?<!\\)\|", raw.strip("|"))]
        if len(cells) < 3:
            continue
        first = cells[0]
        # top-level: first cell holds one or more `name` tokens separated by /
        tops = re.findall(r"`([A-Za-z][A-Za-z0-9]*)`", first)
        if first.startswith("`providers[]`"):
            top.append("providers")
            row += re.findall(r"`([A-Za-z][A-Za-z0-9]*)`", cells[2])
        elif first.startswith("`voice`"):
            top.append("voice")
            voice += re.findall(r"`([A-Za-z][A-Za-z0-9]*)`", cells[2])
        else:
            top += tops
    return sorted(set(top)), sorted(set(row)), sorted(set(voice))


TS_IFACE = re.compile(r"export interface (ConfigView|ConfigViewProvider|ConfigViewVoice) \{(.*?)\r?\n\}", re.S)


def page_keys(path):
    text = open(path, encoding="utf-8").read()
    got = {}
    for m in TS_IFACE.finditer(text):
        name = m.group(1)
        body = m.group(2)
        keys = re.findall(r"^\s*([a-zA-Z][a-zA-Z0-9]*)\??\s*:", body, re.M)
        got[name] = sorted(keys)
    return got["ConfigView"], got["ConfigViewProvider"], got["ConfigViewVoice"]


def show(label, go, spec, page):
    print(f"== {label} ==")
    print(f"  go  ({len(go)}): {go}")
    print(f"  spec({len(spec)}): {spec}")
    print(f"  page({len(page)}): {page}")
    gs, sg = sorted(set(go) - set(spec)), sorted(set(spec) - set(go))
    gp, pg = sorted(set(go) - set(page)), sorted(set(page) - set(go))
    sp, ps = sorted(set(spec) - set(page)), sorted(set(page) - set(spec))
    print(f"  go-vs-spec : go-only {gs} | spec-only {sg}")
    print(f"  go-vs-page : go-only {gp} | page-only {pg}")
    print(f"  spec-vs-page: spec-only {sp} | page-only {ps}")


def main():
    go = go_keys(f"{ROOT}/.scratch/wisp/probes/313/v1/logs/20-probe-v1-raw.txt")
    spec_top, spec_row, spec_voice = spec_keys(f"{ROOT}/docs/evidence/s1/313-panel-ask.md")
    page_top, page_row, page_voice = page_keys(f"{ROOT}/frontend/src/lib/config-view.ts")

    # The empty/full payload must agree with each other at the top level.
    assert go["top"][0] == go["top"][1], "top-level keys differ between empty and full payloads"
    assert go["voice"][0] == go["voice"][1], "voice keys differ between payloads"
    show("top-level", go["top"][1], spec_top, page_top)
    show("provider row", go["row"][1], spec_row, page_row)
    show("voice", go["voice"][1], spec_voice, page_voice)
    print("empty payload providers value:", json.loads(
        [l.split(" ", 1)[1] for l in open(f"{ROOT}/.scratch/wisp/probes/313/v1/logs/20-probe-v1-raw.txt",
                                         encoding="utf-8", errors="replace")
         if l.startswith("V1-EMPTY-PAYLOAD ")][0])["providers"])


if __name__ == "__main__":
    main()
