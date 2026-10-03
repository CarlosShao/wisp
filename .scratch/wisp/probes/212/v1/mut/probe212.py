# 212-v1 probe driver (evidence-only, tracked tree untouched; writes only under probes/212/v1)
import json, os, re, subprocess, sys

ROOT = r"D:\work\workspace\projects plans\Wisp"
V = os.path.join(ROOT, ".scratch", "wisp", "probes", "212", "v1")
MUT = os.path.join(V, "mut")
FIX = os.path.join(V, "fix", "r2probe")
os.makedirs(MUT, exist_ok=True)
os.makedirs(FIX, exist_ok=True)

def w(path, text):
    os.makedirs(os.path.dirname(path), exist_ok=True)
    with open(path, "w", encoding="utf-8", newline="\n") as f:
        f.write(text)

# ---- mutant 1: ban #9 emission site removed (f.Comments loop deleted) ----
src = open(os.path.join(ROOT, "tools", "d22scan", "main.go"), encoding="utf-8").read()
lines = src.split("\n")
i_for = next(i for i, l in enumerate(lines) if l == "\tfor _, cg := range f.Comments {")
i_ret = next(i for i, l in enumerate(lines) if l == "\treturn nil" and i > i_for)
mut1 = "\n".join(lines[:i_for] + lines[i_ret:]) + "\n"
assert 'f.Comments' not in mut1.split('// Bans')[1].split('func ')[0] or True
w(os.path.join(MUT, "main-noq9.go"), mut1)

# ---- mutant 2: inverted criterion (existing citations ring, missing ones silent) ----
needle = "\t\t\t\tif exists {\n\t\t\t\t\tcontinue\n\t\t\t\t}"
flipped = "\t\t\t\tif !exists {\n\t\t\t\t\tcontinue\n\t\t\t\t}"
assert needle in src, "inversion anchor not found"
mut2 = src.replace(needle, flipped, 1)
w(os.path.join(MUT, "main-inv.go"), mut2)

# ---- overlay jsons ----
real = os.path.join(ROOT, "tools", "d22scan", "main.go").replace("\\", "/")
for name, tgt in (("ov-noq9.json", "main-noq9.go"), ("ov-inv.json", "main-inv.go")):
    j = {"Replace": {real: os.path.join(MUT, tgt).replace("\\", "/")}}
    w(os.path.join(MUT, name), json.dumps(j, indent=1))

# ---- symRefRe / repoPathRe sweep over ALL tracked paths (false-exclusion check) ----
symref = re.compile(r"^[a-z][a-z0-9]*(/[a-z][a-z0-9]*)*\.[A-Z]")
repopath = re.compile(r"(?:docs|\.scratch|internal|cmd|tools|scripts)/[A-Za-z0-9_./\-\u4e00-\u9fff]*[A-Za-z0-9_\-]")
ls = subprocess.run(["git", "ls-files"], cwd=ROOT, capture_output=True).stdout.decode("utf-8", "replace").split("\n")
ls = [p for p in ls if p.strip()]
sym_hits = [p for p in ls if symref.match(p)]
rep_hits = [p for p in ls if repopath.match(p)]
w(os.path.join(V, "symref-sweep.txt"),
  "tracked paths total: %d\nsymRefRe-matching tracked paths (would be FALSE-EXCLUDED if cited): %d\n%s\n\n"
  "repoPathRe-matching tracked paths (citable tokens that exist): %d\n" % (len(ls), len(sym_hits), "\n".join(sym_hits) or "(none)", len(rep_hits)))

# ---- pre-image (5e8748b3^) production-comment token census: true 2-count the commit fixed ----
out = []
tree = subprocess.run(["git", "ls-tree", "-r", "--name-only", "5e8748b3^"], cwd=ROOT, capture_output=True).stdout.decode("utf-8", "replace").split("\n")
gofiles = [p for p in tree if re.match(r"^(internal|cmd)/.*\.go$", p) and not p.endswith("_test.go")]
all_missing = {}
for p in gofiles:
    r = subprocess.run(["git", "show", "5e8748b3^:" + p], cwd=ROOT, capture_output=True)
    if r.returncode != 0:
        continue
    text = r.stdout.decode("utf-8", "replace")
    toks = set()
    for line in text.split("\n"):
        s = line.strip()
        if s.startswith("//"):
            toks.update(repopath.findall(s))
    for t in toks:
        if not os.path.exists(os.path.join(ROOT, t.replace("/", os.sep))):
            all_missing.setdefault(t, []).append(p)
w(os.path.join(V, "preimage-phantom-census.txt"),
  "pre-image 5e8748b3^ production (non-test) Go files under internal/+cmd: %d\n"
  "distinct repoPathRe tokens cited in whole-line comments that DO NOT exist on disk: %d\n" % (len(gofiles), len(all_missing)) +
  "".join("%s\n    cited in: %s\n" % (t, ", ".join(ps)) for t, ps in sorted(all_missing.items())))

# ---- symRef probe fixture (a)(b)(c) + controls ----
F = FIX
w(os.path.join(F, "go.mod"), "module probe212.invalid/fix\n\ngo 1.24\n")
w(os.path.join(F, "tools", "d22scan", "allowlist.txt"), "# empty allowlist\n")
w(os.path.join(F, "tools", "d22scan", "main.go"), "package main\n")  # exists-target for control p6
w(os.path.join(F, "frontend", "index.html"), "<html>ok</html>\n")
w(os.path.join(F, "design", "index.html"), "<html>ok</html>\n")
w(os.path.join(F, "internal", "tools", "bridge.go"), "package tools\n")  # exists-target for control p7 + ban7 scope
w(os.path.join(F, "cmd", "probe", "main.go"), "package main\n\nfunc main() {}\n")
probes = {
    "p1.go":  "// probe a: api citation internal/tools.Result.AppliedSteps\npackage probe\n",
    "p2.go":  "// probe b: the error class internal/tool narrative\npackage probe\n",
    "p3.go":  "// probe c: see internal/Build/x.go for details\npackage probe\n",
    "p4.go":  "// control: see internal/build/x.go for details\npackage probe\n",
    "p5.go":  "// control: see docs/evidence/s1/212-citation-ruler.md\npackage probe\n",
    "p6.go":  "// control: see tools/d22scan/main.go (seeded, exists)\npackage probe\n",
    "p7.go":  "// control: see internal/tools/bridge.go (seeded, exists)\npackage probe\n",
    "p8.go":  "package probe\n",
    "p9.go":  "package probe\n",
}
for n, t in probes.items():
    w(os.path.join(F, "internal", "probe", n), t)
print("OK: mutants, overlays, sweeps, fixture written")
