#!/usr/bin/env python3
"""Ticket 222 adversarial acceptance leg (222-v1): mutation builder.

Mutations NEVER touch the working tree. Each mutation is produced from
`git cat-file blob HEAD:<path>` (the committed producer code), rewritten in
memory, written under .scratch/wisp/probes/222/v1/mutations/<id>/, and handed
to `go test -overlay=<json>` which substitutes the file at compile time only.

usage: python mut222.py <mutation-id>
"""
import json
import os
import subprocess
import sys

REPO = os.path.abspath(os.getcwd())  # run from the repo root, always
MUT_ROOT = os.path.join(REPO, ".scratch/wisp/probes/222/v1/mutations")


def blob(rel):
    out = subprocess.run(["git", "cat-file", "blob", "HEAD:" + rel],
                         cwd=REPO, capture_output=True, check=True)
    return out.stdout.decode("utf-8")


def write_overlay(mid, entries):
    d = os.path.join(MUT_ROOT, mid)
    os.makedirs(d, exist_ok=True)
    replace = {}
    for rel, text in entries.items():
        p = os.path.join(d, os.path.basename(rel))
        with open(p, "w", encoding="utf-8", newline="\n") as f:
            f.write(text)
        replace[(REPO + "/" + rel).replace("\\", "/")] = p.replace("\\", "/")
    oj = os.path.join(d, "overlay.json")
    with open(oj, "w", encoding="utf-8", newline="\n") as f:
        json.dump({"Replace": replace}, f, indent=1)
    print("OVERLAY=" + oj.replace("\\", "/"))
    for k, v in replace.items():
        print("  replaces " + k + " <- " + v)


def sub1(text, old, new, mid, tag):
    if text.count(old) < 1:
        raise SystemExit("MUTATION %s anchors missing: %r" % (mid, old[:60]))
    return text.replace(old, new, 1)


def main():
    mid = sys.argv[1]

    if mid == "asis":  # control overlay: byte-identical copies, proves the harness itself is inert
        src = {rel: blob(rel) for rel in
               ("internal/tools/bridge.go", "internal/tools/subagent_197.go")}
        write_overlay(mid, src)
        return

    if mid == "m1-drop-giveback":
        # AC#1/AC#2 positive control = the pre-fix shape: the spawn holds its
        # D38d token across the whole wait. subagent_197.go:356.
        rel = "internal/tools/subagent_197.go"
        t = blob(rel)
        t = sub1(t, "\tgiveBackWhileWaiting(ctx)\n", "", mid, "m1")
        write_overlay(mid, {rel: t})
        return

    if mid == "m2-drop-sync-once":
        # r1's own self-reported mutation 2 (ticket face, 收 222-r1 §给下一位的待办 ②):
        # remove the sync.Once from inFlightSlot.giveBack (bridge.go:670-680).
        rel = "internal/tools/bridge.go"
        t = blob(rel)
        old = """func (s *inFlightSlot) giveBack() bool {
	gave := false
	s.once.Do(func() {
		if s.sem != nil {
			<-s.sem
			s.sem = nil
			gave = true
		}
	})
	return gave
}"""
        new = """func (s *inFlightSlot) giveBack() bool {
	gave := false
	if s.sem != nil {
		<-s.sem
		s.sem = nil
		gave = true
	}
	return gave
}"""
        t = sub1(t, old, new, mid, "m2")
        write_overlay(mid, {rel: t})
        return

    if mid == "m3-child-on-fake-dir":
        # AC#1's own hollow-point attack: put the harness back to
        # subagent_197_test.go's shape (children handed h.dir, the fake197Dir).
        # Test file is mutated only inside the overlay, never in the worktree.
        rel = "internal/tools/subagent_222_test.go"
        t = blob(rel)
        t = sub1(t, "\t\tParentTools: h.bridge,", "\t\tParentTools: &fake197Dir{},", mid, "m3a")
        t = sub1(t, "\t\t\t\tTools:     h.bridge,", "\t\t\t\tTools:     &fake197Dir{},", mid, "m3a")
        write_overlay(mid, {rel: t})
        return

    if mid == "m4-child-bypasses-bridge":
        # The subtler AC#1 attack: the child's tool surface still RUNS the probe,
        # but NOT through the bridge (a directory that calls Execute directly).
        # If these legs stayed green, AC#1 would be hollow.
        rel = "internal/tools/subagent_222_test.go"
        t = blob(rel)
        t = sub1(t, "\t\tParentTools: h.bridge,", "\t\tParentTools: h.bypass(),", mid, "m4")
        t = sub1(t, "\t\t\t\tTools:     h.bridge,", "\t\t\t\tTools:     h.bypass(),", mid, "m4")
        t = t + """
// bypass222 is a tool directory that runs the SAME probe tool without ever
// taking a bridge slot - the shape ticket 222 says production does not have.
type bypass222 struct{ h *h222 }

func (b bypass222) Tools(_ context.Context) ([]agent.ToolInfo, error) {
	return []agent.ToolInfo{{Name: probe222Tag, RiskLevel: "L0"}}, nil
}

func (b bypass222) Execute(ctx context.Context, req agent.ToolRequest) (agent.ToolOutcome, error) {
	res, err := b.h.probe.Execute(ctx, req.Args, nil)
	return agent.ToolOutcome{Text: res.Text, IsError: res.IsError}, err
}

func (h *h222) bypass() agent.ToolProvider { return bypass222{h: h} }
"""
        write_overlay(mid, {rel: t})
        return

    if mid == "m5-ceiling-to-8":
        # AC#2's reverse-control vacuity attack: raise the frozen ceiling itself.
        rel = "internal/tools/bridge.go"
        t = blob(rel)
        t = sub1(t, "const MaxToolConcurrency = 4", "const MaxToolConcurrency = 8", mid, "m5")
        write_overlay(mid, {rel: t})
        return

    if mid == "m6-pool-to-8":
        # AC#5 attack: raise the subagent POOL to 8 while the ceiling stays 4.
        rel = "internal/tools/subagent_197.go"
        t = blob(rel)
        t = sub1(t, "MaxConcurrentSubagents = 4", "MaxConcurrentSubagents = 8", mid, "m6")
        write_overlay(mid, {rel: t})
        return

    if mid == "m7-giveback-after-delta":
        # AC#1's "no deadline needed" attack: move the handoff AFTER the two
        # "已派生" lines, i.e. the position that would make the parked-delta
        # marker stop meaning "slot already returned".
        rel = "internal/tools/subagent_197.go"
        t = blob(rel)
        old = """	giveBackWhileWaiting(ctx)
	t.feed(bg.ID, "已派生："+label)"""
        new = """	t.feed(bg.ID, "已派生："+label)"""
        t = sub1(t, old, new, mid, "m7")
        anchor = """	done := make(chan agent.Result, 1)
	var releaseOnce sync.Once"""
        t = sub1(t, anchor, "\tgiveBackWhileWaiting(ctx)\n" + anchor, mid, "m7")
        write_overlay(mid, {rel: t})
        return

    if mid == "m8-double-drain-for-real":
        # The shape r1's mutation-2 claim actually names ("同一枚许可被扣两次"):
        # BOTH guards removed, so giveBack drains once per call site - the tool's
        # explicit call and run's deferred call each take a token.
        rel = "internal/tools/bridge.go"
        t = blob(rel)
        old = """func (s *inFlightSlot) giveBack() bool {
	gave := false
	s.once.Do(func() {
		if s.sem != nil {
			<-s.sem
			s.sem = nil
			gave = true
		}
	})
	return gave
}"""
        new = """func (s *inFlightSlot) giveBack() bool {
	<-s.sem
	return true
}"""
        t = sub1(t, old, new, mid, "m8")
        write_overlay(mid, {rel: t})
        return

    if mid == "m9-drop-nil-guard":
        # Isolate which guard is load-bearing: keep sync.Once, drop `s.sem = nil`.
        rel = "internal/tools/bridge.go"
        t = blob(rel)
        old = """		if s.sem != nil {
			<-s.sem
			s.sem = nil
			gave = true
		}"""
        new = """		if s.sem != nil {
			<-s.sem
			gave = true
		}"""
        t = sub1(t, old, new, mid, "m9")
        write_overlay(mid, {rel: t})
        return

    if mid == "m10-budget-1ms":
        # AC#3's remaining half, measured ON THE SHIPPED producer code: shrink
        # only the test-local C22 knob so the child cannot possibly finish inside
        # its parent's per-tool budget. If the parent still gets the bridge's
        # timeout wording instead of its child's conclusion, "父侧结论通道要么真"
        # is not closed by ticket 222-r1's fix.
        rel = "internal/tools/subagent_222_test.go"
        t = blob(rel)
        t = sub1(t, "h222PreFixBudget = 3 * time.Second",
                 "h222PreFixBudget = time.Millisecond", mid, "m10")
        write_overlay(mid, {rel: t})
        return

    if mid == "m11-budget-50ms-gated":
        # AC#3's remaining half, deterministically: keep the fix in place, shrink
        # the test-local C22 knob to 50ms and hold the children inside their
        # (gated) bridge call for 200ms, so the parent's budget provably expires
        # while a healthy child is still working. Red = the parent still cannot
        # get its child's conclusion on the SHIPPED producer code.
        rel = "internal/tools/subagent_222_test.go"
        t = blob(rel)
        t = sub1(t, "h222PreFixBudget = 3 * time.Second",
                 "h222PreFixBudget = 50 * time.Millisecond", mid, "m11")
        t = sub1(t, "\tclose(h.probe.gate222)\n\tresults := await222Tokens(t, \"父任务拿到结论\", h.parents, n)\n\tfor i, res := range results {\n\t\tif res.IsError {\n\t\t\tt.Errorf(\"父任务 %d 收到的是错误而不是结论：%q\", i, res.Text)",
                 "\ttime.Sleep(200 * time.Millisecond)\n\tclose(h.probe.gate222)\n\tresults := await222Tokens(t, \"父任务拿到结论\", h.parents, n)\n\tfor i, res := range results {\n\t\tif res.IsError {\n\t\t\tt.Errorf(\"父任务 %d 收到的是错误而不是结论：%q\", i, res.Text)",
                 mid, "m11")
        write_overlay(mid, {rel: t})
        return

    raise SystemExit("unknown mutation id: " + mid)


main()
