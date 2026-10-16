# mutate270v1.py - the mutation bench for the 270-v1 adjudication leg.
#
# Why python and not -overlay: overlay replaces only the bytes the COMPILER
# sees, while several rulers in internal/tools read their source at RUNTIME
# (paths_twocontainments_252_r2_test.go:194 os.ReadDir, task_cancel_221_legs_test
# .go:224 os.ReadFile). Mutations therefore have to be written onto the physical
# file inside the git-archive export tree, and restored by sha256.
#
# No line numbers are used: every op locates its target by the exact condition
# text and REFUSES (rc=3 / rc=4) unless it finds exactly one occurrence, so a
# drifted file cannot silently mutate the wrong branch.
#
# usage: python mutate270v1.py <target.go> <op>
#        python mutate270v1.py <target.go> restore <pristine.go>
import hashlib
import sys

SENT = "宿主没有给出派生用的装配（BaseOptions/ParentTools 未接线，fail-closed：不起第二套运行时）"
ROSTER_SENT = "任务名册未接线（fail-closed：拒绝派生子代理，派生了也没有地方登记它）"

GUARD_BLOCK = (
    "\tif t.d.BaseOptions == nil || t.d.ParentTools == nil {\n"
    "\t\treturn Result{Text: \"" + SENT + "\", IsError: true}, nil\n"
    "\t}\n"
)
JIA_BLOCK = (
    "\tif t.d.BaseOptions == nil {\n"
    "\t\trelease()\n"
    "\t\treturn Result{Text: \"" + SENT + "\", IsError: true}, nil\n"
    "\t}\n"
)
# anti-forms: what the call site would look like if 甲 were written wrong.
JIA_PANIC = (
    "\tif t.d.BaseOptions == nil {\n"
    "\t\trelease()\n"
    "\t\tpanic(\"270v1 anti-form: nil BaseOptions reaches the call and panics\")\n"
    "\t}\n"
)
JIA_SILENT = (
    "\tif t.d.BaseOptions == nil {\n"
    "\t\trelease()\n"
    "\t\treturn Result{Text: \"这一轮没有子代理被派生。\", IsError: false}, nil\n"
    "\t}\n"
)
JIA_NO_RELEASE = (
    "\tif t.d.BaseOptions == nil {\n"
    "\t\treturn Result{Text: \"" + SENT + "\", IsError: true}, nil\n"
    "\t}\n"
)
JIA_BORROWED = (
    "\tif t.d.BaseOptions == nil {\n"
    "\t\trelease()\n"
    "\t\treturn Result{Text: \"" + ROSTER_SENT + "\", IsError: true}, nil\n"
    "\t}\n"
)
# the || split one half at a time (AC#2's shape, measured even though form 甲
# leaves the branch unsplit - "半边摘" is still a question the tree can answer).
GUARD_BASE_ONLY = (
    "\tif t.d.BaseOptions == nil {\n"
    "\t\treturn Result{Text: \"" + SENT + "\", IsError: true}, nil\n"
    "\t}\n"
)
GUARD_PARENT_ONLY = (
    "\tif t.d.ParentTools == nil {\n"
    "\t\treturn Result{Text: \"" + SENT + "\", IsError: true}, nil\n"
    "\t}\n"
)

# op -> list of (find, replace) applied in order; "strip-guard" prefixes are
# combined so the anti-forms are REACHABLE (with the guard in place they never
# are, which is exactly the shape form 甲 leaves behind).
OPS = {
    "strip-guard": [(GUARD_BLOCK, "")],
    "strip-jia": [(JIA_BLOCK, "")],
    "strip-both": [(GUARD_BLOCK, ""), (JIA_BLOCK, "")],
    "guard-off+jia-panic": [(GUARD_BLOCK, ""), (JIA_BLOCK, JIA_PANIC)],
    "guard-off+jia-silent": [(GUARD_BLOCK, ""), (JIA_BLOCK, JIA_SILENT)],
    "guard-off+jia-no-release": [(GUARD_BLOCK, ""), (JIA_BLOCK, JIA_NO_RELEASE)],
    "guard-off+jia-borrowed-sentence": [(GUARD_BLOCK, ""), (JIA_BLOCK, JIA_BORROWED)],
    "guard-keeps-baseoptions-half": [(GUARD_BLOCK, GUARD_BASE_ONLY)],
    "guard-keeps-parenttools-half": [(GUARD_BLOCK, GUARD_PARENT_ONLY)],
}


def main() -> int:
    target, op = sys.argv[1], sys.argv[2]
    data = open(target, encoding="utf-8", newline="").read()
    if op == "restore":
        src = open(sys.argv[3], encoding="utf-8", newline="").read()
        open(target, "w", encoding="utf-8", newline="").write(src)
        print("restored-from=%s sha256=%s" % (sys.argv[3], hashlib.sha256(src.encode()).hexdigest()))
        return 0
    if op not in OPS:
        print("unknown op %r" % op)
        return 2
    for find, repl in OPS[op]:
        n = data.count(find)
        if n != 1:
            print("REFUSED: op %s expects exactly 1 occurrence, found %d" % (op, n))
            return 3
        data = data.replace(find, repl, 1)
    if data.count(SENT) not in (0, 1, 2):
        print("REFUSED: sentence occurrences unreadable")
        return 4
    open(target, "w", encoding="utf-8", newline="").write(data)
    print("op=%s sha256=%s sent-occurrences-left=%d" % (op, hashlib.sha256(data.encode()).hexdigest(), data.count(SENT)))
    return 0


if __name__ == "__main__":
    sys.exit(main())
