#!/usr/bin/env python3
# 260-r4 突变台件：五发改动，每发只洗一枚判据，跑完立刻还原（还原靠备份回写＋git status 核空）。
# 尺：approval 定向 5＋2 枚用例 / cmd/wisp 定向 260R4 三枚用例，红句逐字进 logs/。
import subprocess
import sys
import shutil
import os

ROOT = subprocess.run(["git", "rev-parse", "--show-toplevel"], capture_output=True,
                      text=True).stdout.strip()
APPROVAL = os.path.join(ROOT, "internal", "agent", "approval")
CMDWISP = os.path.join(ROOT, "cmd", "wisp")
LOGS = os.path.join(ROOT, ".scratch", "wisp", "probes", "260", "r4", "logs")

SPELLENDAR = """\tif src == nil {
		return defaultCancelKeySpelling
	}
	if s := src(); s != "" {
		return s
	}
	return defaultCancelKeySpelling
"""
SPELLEDNEW = """\t_ = src
	return defaultCancelKeySpelling
"""

STATUSOLD = """\tif ch == ChannelEsc && !loaded {
		return channelNames[ch]
	}
	return channelLabel(ch)
"""
STATUSNEW = "\treturn channelLabel(ch)\n"

GATEOLD = """\t\t\t\tchannelLabel(v.Channel))"""
GATENEW = """\t\t\t\tchannelNames[v.Channel])"""

REPORTOLD = "\t\tb.WriteString(channelLabel(r.Channel))"
REPORTNEW = "\t\tb.WriteString(channelNames[r.Channel])"

INSTALL_OLD = "\tapproval.SetCancelKeySpelling(ra.cancelKeySpelling)"
INSTALL_NEW = "\t_ = approval.SetCancelKeySpelling"

MUTATIONS = [
    ("M1 标签算一次（装了就等于没装）", os.path.join(APPROVAL, "approval.go"),
     SPELLENDAR, SPELLEDNEW, ["approval", "cmdwisp"]),
    ("M2 未加载那行也指枚键名", os.path.join(APPROVAL, "approval.go"),
     STATUSOLD, STATUSNEW, ["approval", "cmdwisp"]),
    ("M3 gate.go:314 那处读取退回旧表", os.path.join(APPROVAL, "gate.go"),
     GATEOLD, GATENEW, ["approval"]),
    ("M4 装配根不再装读口（生产链断在这里）", os.path.join(CMDWISP, "resident_approval_windows.go"),
     INSTALL_OLD, INSTALL_NEW, ["cmdwisp"]),
    ("M5 report.go:142 那处读取退回旧表", os.path.join(APPROVAL, "report.go"),
     REPORTOLD, REPORTNEW, ["approval"]),
]


def run(args, env=None):
    p = subprocess.run(args, cwd=ROOT, capture_output=True, text=True,
                       encoding="utf-8", errors="replace", env=env)
    return p.returncode, (p.stdout or "") + (p.stderr or "")


def approval_run():
    return run(["go", "test", "-count=1", "-v", "-run", "TestTicket260",
                "./internal/agent/approval/"])


def cmdwisp_run():
    env = dict(os.environ)
    env["PATH"] = os.pathsep.join([os.path.join(ROOT, "third_party", "sherpa-onnx"), env["PATH"]])
    return run(["go", "test", "-count=1", "-v", "-run", "TestTicket260", "./cmd/wisp/"], env=env)


def reds(out):
    keep = []
    lines = out.splitlines()
    for i, ln in enumerate(lines):
        if ln.startswith("--- FAIL") or ln.startswith("    resident_") or ln.startswith("    cancel_key"):
            keep.append(ln)
    return "\n".join(keep)


def main():
    os.makedirs(LOGS, exist_ok=True)
    report = []
    for name, path, old, new, faces in MUTATIONS:
        backup = path + ".r4bak"
        shutil.copy2(path, backup)
        with open(path, "r", encoding="utf-8") as f:
            src = f.read()
        if src.count(old) != 1:
            report.append("!! %s：靶文本命中 %d 次，跳过" % (name, src.count(old)))
            os.remove(backup)
            continue
        with open(path, "w", encoding="utf-8") as f:
            f.write(src.replace(old, new, 1))
        chunk = ["=== %s ===" % name]
        for face in faces:
            if face == "approval":
                rc, out = approval_run()
            else:
                rc, out = cmdwisp_run()
            fails = [ln for ln in out.splitlines() if ln.startswith("--- FAIL")]
            passes = [ln for ln in out.splitlines() if ln.startswith("--- PASS")]
            chunk.append("[%s] rc=%d FAIL=%d PASS=%d" % (face, rc, len(fails), len(passes)))
            chunk.extend(fails)
            chunk.append("--- 红句逐字 ---")
            chunk.append(reds(out))
        # restore
        shutil.copy2(backup, path)
        os.remove(backup)
        rel = os.path.relpath(path, ROOT).replace("\\", "/")
        rc, out = run(["git", "status", "--porcelain", "--", rel])
        chunk.append("还原后 git status 该文件：%s" % (out.strip() or "(空＝与 HEAD 一致)"))
        report.append("\n".join(chunk))
        print(chunk[-1], "|", name, flush=True)
    with open(os.path.join(LOGS, "mutations.txt"), "w", encoding="utf-8") as f:
        f.write("\n\n".join(report) + "\n")
    print("WROTE logs/mutations.txt")


if __name__ == "__main__":
    sys.exit(main())
