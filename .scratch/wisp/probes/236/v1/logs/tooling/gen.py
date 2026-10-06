import os, sys, json, subprocess, hashlib

ROOT = "D:/work/workspace/projects plans/Wisp"
AB = "D:/work/workspace/projects plans/Wisp"
OUT = "D:/tmp/wisp236v1/mut"
OV = "D:/tmp/wisp236v1/overlay"
PR = "D:/work/workspace/projects plans/Wisp/.scratch/wisp/probes/236/v1/logs"
for d in (OUT, OV, PR + "/diff", PR + "/overlay", PR + "/copies"):
    os.makedirs(d, exist_ok=True)

TASK, SUB = "internal/tools/task.go", "internal/tools/subagent_197.go"
LEGS = "internal/tools/task_cancel_221_legs_test.go"
FAIL = "internal/tools/failclosed_236_teeth_test.go"

def read(rel):
    with open(os.path.join(ROOT, rel), encoding="utf-8") as f:
        return f.read().split("\n")

def w(basename, name, lines):
    d = os.path.join(OUT, name); os.makedirs(d, exist_ok=True)
    p = os.path.join(d, basename)
    open(p, "w", encoding="utf-8", newline="\n").write("\n".join(lines))
    return p.replace(os.sep, "/")

def must(L, n, exp):
    if L[n-1].strip() != exp.strip():
        print("MISMATCH line %d want=%r got=%r" % (n, exp, L[n-1])); sys.exit(3)

def jp(name, mapping):
    p = os.path.join(OV, name + ".json")
    json.dump({"Replace": mapping}, open(p, "w", encoding="utf-8"), indent=1)
    return p

def R(rel): return AB + "/" + rel
def iffalse(name, rel, n, cond):
    L = read(rel); must(L, n, cond)
    L[n-1] = L[n-1].replace(cond, "if false {")
    return {R(rel): w(os.path.basename(rel), name, L)}
def deleter(name, rel, k1, k2, checks):
    L = read(rel)
    for n, e in checks: must(L, n, e)
    return {R(rel): w(os.path.basename(rel), name, L[:k1-1] + L[k2:])}
def repl(name, rel, n, old, new):
    L = read(rel); must(L, n, old)
    L[n-1] = L[n-1].replace(old, new)
    return {R(rel): w(os.path.basename(rel), name, L)}

runs = {}
L = read(TASK); S = read(SUB); F = read(FAIL); G = read(LEGS)

# ---------- positive control: byte-identical copies of both production files
runs["pos-ctrl"] = {R(TASK): w("task.go", "pos-ctrl", L), R(SUB): w("subagent_197.go", "pos-ctrl", S)}

# ---------- A-* : "落到相邻分支" (if false) -- the attack the ticket predicted green
runs["A-1-cancel-caller-false"] = iffalse("A-1-cancel-caller-false", TASK, 707, 'if caller == "" {')
runs["A-2-cancel-roster-false"] = iffalse("A-2-cancel-roster-false", TASK, 699, 'if t.d.Roster == nil {')
runs["A-3-spawn-roster-false"] = iffalse("A-3-spawn-roster-false", SUB, 253, 'if t.d.Roster == nil {')
runs["A-4-spawn-assembly-false"] = iffalse("A-4-spawn-assembly-false", SUB, 256,
    'if t.d.BaseOptions == nil || t.d.ParentTools == nil {')
runs["A-5-spawn-parentid-false"] = iffalse("A-5-spawn-parentid-false", SUB, 260, 'if parentID == "" {')

# ---------- A-6 : right sentence, wrong behaviour (IsError flipped to a pass)
runs["A-6-cancel-caller-noterror"] = repl("A-6-cancel-caller-noterror", TASK, 708,
    'return Result{Text: "这条调用没有宿主给的任务 id（fail-closed：不知道是谁要停，只能拒）", IsError: true}, nil',
    'return Result{Text: "这条调用没有宿主给的任务 id（fail-closed：不知道是谁要停，只能拒）", IsError: false}, nil')

# ---------- A-7 : the two cancel literals swapped between the two doors
LitR = '"任务名册未接线（fail-closed：拒绝停掉任何任务——名册才是唯一的停法）"'
LitC = '"这条调用没有宿主给的任务 id（fail-closed：不知道是谁要停，只能拒）"'
LT = read(TASK); must(LT, 700, 'return Result{Text: ' + LitR + ', IsError: true}, nil')
must(LT, 708, 'return Result{Text: ' + LitC + ', IsError: true}, nil')
LT[699] = LT[699].replace(LitR, LitC); LT[707] = LT[707].replace(LitC, LitR)
runs["A-7-swap-cancel-literals"] = {R(TASK): w("task.go", "A-7-swap-cancel-literals", LT)}

# ---------- v-* : independent redo of the write legs' whole-branch deletions
runs["v-m1a-del-cancel-caller"] = deleter("v-m1a-del-cancel-caller", TASK, 707, 709,
    [(707, 'if caller == "" {'), (709, '}')])
runs["v-m1b-del-cancel-roster"] = deleter("v-m1b-del-cancel-roster", TASK, 699, 701,
    [(699, 'if t.d.Roster == nil {'), (701, '}')])
runs["v-m1c-del-spawn-roster"] = deleter("v-m1c-del-spawn-roster", SUB, 253, 255,
    [(253, 'if t.d.Roster == nil {'), (255, '}')])
runs["v-m1e-del-spawn-parentid"] = deleter("v-m1e-del-spawn-parentid", SUB, 260, 263,
    [(259, 'parentID := CorrelationID(ctx)'), (260, 'if parentID == "" {'), (263, '}')])
O256 = 'if t.d.BaseOptions == nil || t.d.ParentTools == nil {'
runs["v-m1d1-keep-baseonly"] = repl("v-m1d1-keep-baseonly", SUB, 256, O256, 'if t.d.BaseOptions == nil {')
runs["v-m1d2-keep-parentonly"] = repl("v-m1d2-keep-parentonly", SUB, 256, O256, 'if t.d.ParentTools == nil {')

# ---------- no-recover pairing (the one-ticket-two-shapes hole)
must(F, 232, 'defer func() {'); must(F, 233, 'if rec := recover(); rec != nil {')
NR = F[:231] + F[236:]
must(NR, 232, 'return h.spawn(t, t.Context(), label, prompt)')
runs["ctrl-no-recover-only"] = {R(FAIL): w("failclosed_236_teeth_test.go", "ctrl-no-recover-only", NR)}
runs["A-4-without-recover"] = dict(runs["A-4-spawn-assembly-false"])
runs["A-4-without-recover"][R(FAIL)] = w("failclosed_236_teeth_test.go", "A-4-without-recover", NR)

# ---------- B-* : AC#2
# B-1 teeth-m13: delete only the marker line :23
must(L, 23, "//\ttask.list    -   DEFERRED with five fields, PLAN.md \u00a77 :1531")
runs["B-1-teeth-m13"] = {R(TASK): w("task.go", "B-1-teeth-m13", L[:22] + L[23:])}
# injected copy carrying DEFERRED + task.cancel -> the (a) branch's input
INJ = "//\ttask.cancel  -   DEFERRED injected by 236-v1 (a)-branch probe"
INJcopy = w("task.go", "B-2-inject-cancel-marker", L[:23] + [INJ] + L[23:])
# B-2a : archived ruler's read path rewritten to the injected copy (task.go NOT overlaid)
must(G, 224, 'src, err := os.ReadFile("task.go")')
Ga = list(G); Ga[223] = Ga[223].replace('"task.go"', '"D:/tmp/wisp236v1/mut/B-2-inject-cancel-marker/task.go"')
runs["B-2a-a-branch-readpath"] = {R(LEGS): w("task_cancel_221_legs_test.go", "B-2a-a-branch-readpath", Ga)}
# B-2c : pure overlay of the injected task.go (read path untouched) -> (a) also blind?
runs["B-2c-a-branch-blind"] = {R(TASK): INJcopy}
# B-2b : positive control, read path rewritten to a byte-identical copy
Gb = list(G); Gb[223] = Gb[223].replace('"task.go"', '"D:/tmp/wisp236v1/mut/B-2b-ctrl-copy/task.go"')
ctrlcopy = w("task.go", "B-2b-ctrl-copy", L)
runs["B-2b-readpath-ctrl"] = {R(LEGS): w("task_cancel_221_legs_test.go", "B-2b-readpath-ctrl", Gb),
                              R(TASK): ctrlcopy}
# B-6 : read path -> teeth-m13 copy (marker row gone, :279 still props it)
G6 = list(G); G6[223] = G6[223].replace('"task.go"', '"D:/tmp/wisp236v1/mut/B-1-teeth-m13/task.go"')
runs["B-6-readpath-teeth-m13"] = {R(LEGS): w("task_cancel_221_legs_test.go", "B-6-readpath-teeth-m13", G6)}
# B-7 : read path -> all-markers-gone copy (:23 deleted AND :279 token downcased)
must(L, 279, "// diagnostics, not as a tool surface: task.list is DEFERRED (\u00a77 :1531), and")
NO = L[:22] + L[23:]
must(NO, 278, "// diagnostics, not as a tool surface: task.list is DEFERRED (\u00a77 :1531), and")
NO[277] = NO[277].replace("task.list is DEFERRED", "task.list is deferred")
NOcopy = w("task.go", "B-7-all-markers-gone", NO)
G7 = list(G); G7[223] = G7[223].replace('"task.go"', '"D:/tmp/wisp236v1/mut/B-7-all-markers-gone/task.go"')
runs["B-7-readpath-all-markers-gone"] = {R(LEGS): w("task_cancel_221_legs_test.go", "B-7-readpath-all-markers-gone", G7)}
# B-3 m-2c : register an UNRELATED new task.* row (the missing reading)
must(L, 603, '{Tool: taskCancel{d: d}, Decl: taskCancelDecl()},')
B3 = L[:603] + ['\t\t{Tool: taskNoteRow236v1{taskCancel{d: d}}, Decl: taskCancelDecl()},'] + L[603:]
B3 = B3 + ["", "type taskNoteRow236v1 struct{ taskCancel }", "",
           'func (taskNoteRow236v1) Name() string { return "task.note" }']
runs["B-3-m2c-register-note"] = {R(TASK): w("task.go", "B-3-m2c-register-note", B3)}
# B-4 m-2b : register task.list (independent redo)
B4 = L[:603] + ['\t\t{Tool: taskListRow236v1{taskCancel{d: d}}, Decl: taskCancelDecl()},'] + L[603:]
B4 = B4 + ["", "type taskListRow236v1 struct{ taskCancel }", "",
           'func (taskListRow236v1) Name() string { return "task.list" }']
runs["B-4-m2b-register-list"] = {R(TASK): w("task.go", "B-4-m2b-register-list", B4)}
# B-5 m-2e : unregister task.cancel (independent redo)
runs["B-5-m2e-unregister-cancel"] = {R(TASK): w("task.go", "B-5-m2e-unregister-cancel", L[:602] + L[603:])}

# ---------- emit jsons, diffs, md5s, arithmetic
idx = []
for name, mapping in sorted(runs.items()):
    p = jp(name, mapping)
    json.dump(mapping, open(os.path.join(PR, "overlay", name + ".json"), "w", encoding="utf-8"), indent=1)
    for repo, cp in mapping.items():
        rel = repo.replace(AB + "/", "")
        orig = read(rel)
        mut = open(cp, encoding="utf-8").read().split("\n")
        tmpo = os.path.join(OUT, "_head_" + os.path.basename(rel))
        open(tmpo, "w", encoding="utf-8", newline="\n").write("\n".join(orig))
        dif = subprocess.run(["diff", "-u", "--label", "HEAD/" + rel, "--label", "overlay/" + cp,
                              tmpo, cp], capture_output=True)
        open(os.path.join(PR, "diff", name + "__" + os.path.basename(rel) + ".diff"), "wb").write(dif.stdout)
        h = hashlib.md5(open(cp, "rb").read()).hexdigest()
        ho = hashlib.md5(("\n".join(orig)).encode("utf-8")).hexdigest()
        hd = hashlib.md5(open(os.path.join(ROOT, rel), "rb").read()).hexdigest()
        idx.append((name, rel, cp, h, "head-and-disk=" + ho + "/" + hd, "SAME" if h == hd else "DIFF",
                    len(orig), len(mut)))
open(os.path.join(PR, "copies", "index.txt"), "w", encoding="utf-8").write(
    "mutant\tfile\tcopy\tsrc-md5\tdisk-md5-pair\tsame-as-disk\tn-lines\tn+1-lines\n" +
    "\n".join("\t".join(str(x) for x in r) for r in idx) + "\n")

def cnt(path, tok):
    return sum(1 for ln in open(path, encoding="utf-8").read().split("\n")
               if "DEFERRED" in ln and tok in ln)
arith = []
for cp_name in ["pos-ctrl", "B-1-teeth-m13", "B-2-inject-cancel-marker",
                "B-7-all-markers-gone", "B-3-m2c-register-note", "B-4-m2b-register-list",
                "B-5-m2e-unregister-cancel"]:
    p = os.path.join(OUT, cp_name, "task.go")
    arith.append("%s: DEFERRED-lines=%d co-carry-task.list=%d co-carry-task.cancel=%d" % (
        cp_name, sum(1 for ln in open(p, encoding="utf-8").read().split("\n") if "DEFERRED" in ln),
        cnt(p, "task.list"), cnt(p, "task.cancel")))
open(os.path.join(PR, "copies", "lexical-arithmetic.txt"), "w", encoding="utf-8").write(
    "\n".join(arith) + "\n")
print("\n".join(arith))
print("TOTAL mutants:", len(runs))
for k in sorted(runs): print(" ", k)
