import json
import os
import subprocess
import sys

ROOT = r"D:/work/workspace/projects plans/Wisp"
OUT = os.path.join(ROOT, ".scratch", "wisp", "probes", "221", "v1", "mutations")
os.makedirs(OUT, exist_ok=True)


def blob(rev_path):
    p = subprocess.run(["git", "-C", ROOT, "cat-file", "blob", rev_path],
                       capture_output=True, check=True)
    return p.stdout.decode("utf-8")


def win(abs_posix):
    return abs_posix.replace("/", "\\")


def emit(name, orig_rel, mutant_src):
    src = os.path.join(OUT, name + ".go")
    with open(src, "w", encoding="utf-8", newline="\n") as f:
        f.write(mutant_src)
    overlay = {"Replace": {win(ROOT + "/" + orig_rel): win(src)}}
    with open(os.path.join(OUT, name + ".json"), "w", encoding="utf-8") as f:
        f.write(json.dumps(overlay))
    print(name, "->", os.path.basename(src))


def rep(text, old, new, count=1):
    n = text.count(old)
    if n != count:
        sys.exit("anchor mismatch (%d != %d): %r" % (n, count, old[:70]))
    return text.replace(old, new, count)


TASK = "internal/tools/task.go"
S197 = "internal/tools/subagent_197.go"
C175 = "internal/tools/ticket175r2_stamp_live_test.go"

t = blob("HEAD:" + TASK)
s = blob("HEAD:" + S197)
c = blob("HEAD:" + C175)

# m1: drop task.cancel's registration line (AC#1's ruler must bite)
emit("m1-cancel-unregistered", TASK, rep(
    t, "\t\t{Tool: taskCancel{d: d}, Decl: taskCancelDecl()},\n", ""))

# m2: remove the "caller is not the row's parent" refusal (AC#3 control 1)
emit("m2-no-parent-check", TASK, rep(
    t, "\tif rec.ParentTaskID != caller {", "\tif false && rec.ParentTaskID != caller {"))

# m3: remove the self-stop refusal (AC#3 control 2)
emit("m3-no-selfstop-check", TASK, rep(
    t, "\tif target == caller {", "\tif false && target == caller {"))

# m4: remove the "no host-minted caller id" fail-closed branch
emit("m4-no-empty-caller-guard", TASK, rep(
    t, '\tif caller == "" {', '\tif false && caller == "" {'))

# m5: remove the "roster not wired" fail-closed branch (task.cancel's own copy;
# task.output has an identically shaped guard, so anchor on the whole line)
emit("m5-no-nil-roster-guard", TASK, rep(
    t,
    '\tif t.d.Roster == nil {\n'
    '\t\treturn Result{Text: "任务名册未接线（fail-closed：拒绝停掉任何任务——名册才是唯一的停法）", IsError: true}, nil\n',
    '\tif false && t.d.Roster == nil {\n'
    '\t\treturn Result{Text: "任务名册未接线（fail-closed：拒绝停掉任何任务——名册才是唯一的停法）", IsError: true}, nil\n'))

# m6: stop the row BEFORE judging (the refusals would touch the row)
emit("m6-cancel-before-refusal", TASK, rep(
    t, "\tif target == caller {\n",
    "\tt.d.Roster.Cancel(target) // MUTANT 221-v1: stop before judging\n\n\tif target == caller {\n"))

# m7: declare L0 instead of the frozen L1
emit("m7-level-l0", TASK, rep(t, "Declared:     risk.L1,", "Declared:     risk.L0,"))

# m8: put the old unqualified promise back in Description() (乙形 1+2)
emit("m8-old-promise-returned", S197, rep(
    s,
    '\t\t"它会在任务名册里留下一行有父子关系与状态的记录；"+\n'
    '\t\t"注意：停掉父任务不会级联停掉子代理；task.cancel 只有派生它的那枚父任务能用它单独停孩子——"+\n'
    '\t\t"子代理停兄弟、停自己都一律被拒；由用户停某一枚子代理还是另一条通道（票 181／票 220），今天没有落点；"+\n',
    '\t\t"它会在任务名册里留下一行有父子关系与状态的记录，可以用 task.cancel 单独停它；"+\n'
    '\t\t"注意：停掉父任务不会级联停掉子代理（要停它得单独停）；"+\n'))

# m9: put the old bare promise back in the parent-stopped-waiting receipt (乙形 3)
emit("m9-old-receipt-returned", S197, rep(
    s,
    '\t\t\t\t\t"它的父任务 %s 还能用 task.cancel 单独停它（只有派生它的那一枚能停）；"+\n'
    '\t\t\t\t\t"它自己的流键是 %s。",\n'
    '\t\t\t\tctx.Err(), bg.ID, parentID, SubagentStreamKey(bg.ID)),',
    '\t\t\t\t\t"可以单独停它：它的流键是 %s。",\n'
    '\t\t\t\tctx.Err(), bg.ID, SubagentStreamKey(bg.ID)),'))

# m10: let parent cancellation cascade again
emit("m10-cascade-returns", S197, rep(
    s, "context.WithCancel(context.WithoutCancel(ctx))", "context.WithCancel(ctx)"))

# m11: the stopped child's row lands on a non-terminal-looking name instead of
# the borrowed D43 Muted (attacks AC#2's "行落到 D43 现成名" claim)
emit("m12-finalize-wrong-state", S197, rep(
    s, "\tt.d.Roster.Record(childID, TaskOutput{\n",
    "\tstate = subagentStateSettled // MUTANT 221-v1: stopped row loses Muted\n"
    "\tt.d.Roster.Record(childID, TaskOutput{\n"))

# m12: delete the census answer for task.cancel (extra-B: is the row load-bearing?)
emit("m11-census-answer-deleted", C175, rep(
    c, '\t"task.cancel": "unstamped: a receipt naming which roster row this call stopped",\n', ""))

print("done")
