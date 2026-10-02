#!/usr/bin/env python3
"""33-r8b mutation harness (ticket 33, internal/ball only).

For each named mutation: patch one exact string in one production file, run the
two release nails (plus the cap nail once it exists), save the verbatim reading,
restore the pristine bytes and verify md5. Nothing here ever commits; a mutation
that fails to compile is recorded as such rather than being smoothed over.
"""
import hashlib
import os
import shutil
import subprocess
import sys

REPO = r"D:\work\workspace\projects plans\Wisp"
STA = os.path.join(REPO, "internal", "ball", "sta_windows.go")
BALL = os.path.join(REPO, "internal", "ball", "ball_windows.go")
OUT = os.path.join(REPO, ".scratch", "wisp", "probes", "33", "r8b", "logs")

MUTATIONS = [
    ("M0-not-wired", STA,
     "\tdefer s.releaseThread()",
     "\tdefer s.releaseCOM() // MUT-M0, the shape ebe3bd57 shipped in"),
    ("M1-no-destroy", STA,
     "\tif hwnd != 0 {\n\t\tif r, _, err := pDestroyWindow.Call(uintptr(hwnd)); r == 0 {",
     "\tif false && hwnd != 0 { // MUT-M1\n\t\tif r, _, err := pDestroyWindow.Call(uintptr(hwnd)); r == 0 {"),
    ("M2-no-queue-pump", STA,
     "\treleaseOwnQueueToQuiet()\n\ts.releaseCOM()",
     "\t// MUT-M2: releaseOwnQueueToQuiet() dropped\n\ts.releaseCOM()"),
    ("M3-hwnd-not-forgotten", STA,
     "\thwnd := s.hwnd\n\ts.hwnd = 0\n\ts.mu.Unlock()",
     "\thwnd := s.hwnd\n\t// MUT-M3: s.hwnd not cleared\n\ts.mu.Unlock()"),
    ("M3b-close-does-not-forget", BALL,
     "\t\t\tb.sta.forgetWindow()\n",
     # Kept callable-but-skipped on purpose: deleting the statement outright also
     # deletes the file's only use of "fmt", so the mutation would show up as a
     # build failure instead of as this nail biting.
     "\t\t\tif false { b.sta.forgetWindow() } // MUT-M3b: forgetWindow() dropped\n"),
    ("M4-pm-noremove", STA,
     "const pmRemove = 1",
     "const pmRemove = 0 // MUT-M4"),
]

RUN = os.environ.get("R8B_RUN", "")
TAGS = os.environ.get("R8B_TAGS", "")
TIMEOUT = os.environ.get("R8B_TIMEOUT", "300s")
PREFIX = "winlive-" if "winlive" in TAGS else "plain-"
ONLY = os.environ.get("R8B_ONLY", "")


def md5(path):
    with open(path, "rb") as fh:
        return hashlib.md5(fh.read()).hexdigest()


def run_test(tag):
    env = dict(os.environ)
    third = os.path.join(REPO, "third_party", "sherpa-onnx")
    build = os.path.join(REPO, "build")
    env["PATH"] = third + os.pathsep + build + os.pathsep + env["PATH"]
    cmd = ["go", "test", "./internal/ball/", "-count=1", "-v", "-timeout", TIMEOUT]
    if TAGS:
        cmd[4:4] = ["-tags=" + TAGS]
    if RUN:
        cmd += ["-run", RUN]
    p = subprocess.run(
        cmd,
        cwd=REPO, capture_output=True, text=True, encoding="utf-8", errors="replace",
        env=env, shell=False)
    out = p.stdout + p.stderr
    with open(os.path.join(OUT, "mutation-%s%s.txt" % (PREFIX, tag)), "w", encoding="utf-8", newline="\n") as fh:
        fh.write(out)
    fails = [ln for ln in out.splitlines() if ln.startswith("--- FAIL")]
    passes = [ln for ln in out.splitlines() if ln.startswith("--- PASS")]
    return p.returncode, len(fails), len(passes), fails, ("build failed" in out or "FAIL\t[build failed]" in out)


pristine = {f: md5(f) for f in (STA, BALL)}
print("pristine md5:", pristine)
for src_path, digest in pristine.items():
    shutil.copyfile(src_path, os.path.join(OUT, "pristine-" + os.path.basename(src_path)))
    if md5(os.path.join(OUT, "pristine-" + os.path.basename(src_path))) != digest:
        sys.exit("pristine copy of %s does not hash the same - refusing to run" % src_path)
report = []
for tag, path, old, new in MUTATIONS:
    if ONLY and not tag.startswith(ONLY):
        continue
    src = open(path, encoding="utf-8", newline="").read()
    n = src.count(old)
    if n != 1:
        report.append("%s SKIP anchor-not-unique count=%d" % (tag, n))
        continue
    if md5(path) != pristine[path]:
        report.append("%s ABORT file-not-pristine-before-patch %s != %s" % (tag, md5(path), pristine[path]))
        break
    open(path, "w", encoding="utf-8", newline="").write(src.replace(old, new, 1))
    try:
        rc, nf, np, fails, broke = run_test(tag)
    finally:
        open(path, "w", encoding="utf-8", newline="").write(src)
    ok = md5(path) == pristine[path]
    report.append("%s rc=%d FAIL=%d PASS=%d restored_md5=%s build_failed=%s | %s"
                  % (tag, rc, nf, np, "yes" if ok else "NO", broke, " ; ".join(fails) if fails else "(no red)"))
    print(report[-1])
    if not ok:
        shutil.copyfile(os.path.join(OUT, "pristine-" + os.path.basename(path)), path)
        report.append("%s RESTORED-FROM-COPY md5=%s" % (tag, md5(path) == pristine[path]))
        print(report[-1])
for line in report:
    print(line)
with open(os.path.join(OUT, "mutation-report.txt"), "w", encoding="utf-8", newline="\n") as fh:
    fh.write("\n".join(report) + "\n")
print("final md5:", {f: md5(f) for f in (STA, BALL)})
