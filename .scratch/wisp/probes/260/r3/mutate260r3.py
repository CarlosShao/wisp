import hashlib
import os
import subprocess
import sys

APPROVAL = "internal/agent/approval/approval.go"
RESIDENT = "cmd/wisp/resident_approval_windows.go"

# (id, file, old bytes, new bytes, -run pattern, package)
MUTATIONS = [
    (
        "M1",
        RESIDENT,
        'return fmt.Sprintf("按 %s 否决了卡片 %s（%s / %s）：该调用未执行，答案已入审计",\n\t\tkey, card.CorrelationID, card.Level, card.Tool)'.encode(),
        'return fmt.Sprintf("按 Esc 否决了卡片 %s（%s / %s）：该调用未执行，答案已入审计",\n\t\tcard.CorrelationID, card.Level, card.Tool)'.encode(),
        "TestTicket260R3SeededKeyReplacesEsc",
        "./cmd/wisp/",
    ),
    (
        "M2",
        RESIDENT,
        'ra.residentAuditf("resident-approval: 审批门已装配进常驻进程，取消通道 %s 已加载（本票只落 %s 一条通道；"+\n\t\t"单击球 / KWS 否决词 / 面板拒绝三条仍按各自归口未接入）", key, key)'.encode(),
        '_ = key\n\tra.residentAuditf("resident-approval: 审批门已装配进常驻进程，取消通道 Esc 已加载（本票只落 Esc 一条通道；" +\n\t\t"单击球 / KWS 否决词 / 面板拒绝三条仍按各自归口未接入）")'.encode(),
        "TestTicket260R3SeededKeyReplacesEsc",
        "./cmd/wisp/",
    ),
    (
        "M3",
        RESIDENT,
        '"四条否决通道全部保持未加载（取消键无处可借，卡片无处可呈）")'.encode(),
        '"四条否决通道全部保持未加载（Esc 无处可借，卡片无处可呈）")'.encode(),
        "TestTicket260R3UnloadBranchesNameNoKey",
        "./cmd/wisp/",
    ),
    (
        "M4",
        RESIDENT,
        '\t\t\tif bd.Name == "cancel" && bd.Binding != "" {'.encode(),
        '\t\t\tif false && bd.Name == "cancel" && bd.Binding != "" {'.encode(),
        "TestTicket260R3",
        "./cmd/wisp/",
    ),
    (
        "M5",
        APPROVAL,
        'return "快捷键取消不可用"'.encode(),
        'return "Esc 取消不可用"'.encode(),
        "TestTicket260R3",
        "./internal/agent/approval/",
    ),
]

if __name__ == "__main__":
    which = sys.argv[1]
    for mid, path, old, new, pattern, pkg in MUTATIONS:
        if mid != which:
            continue
        pristine = open(path, "rb").read()
        md5_before = hashlib.md5(pristine).hexdigest()
        if pristine.count(old) != 1:
            print("%s ABORT: marker matched %d times" % (mid, pristine.count(old)))
            sys.exit(2)
        open(path, "wb").write(pristine.replace(old, new))
        if len(open(path, "rb").read()) == 0:
            print("FATAL truncated")
            sys.exit(3)
        env = dict(os.environ)
        env["PATH"] = os.path.join(os.getcwd(), "third_party", "sherpa-onnx") + os.pathsep + env["PATH"]
        proc = subprocess.run(["go", "test", "-count=1", "-v", "-run", pattern, pkg],
                              capture_output=True, text=True, encoding="utf-8",
                              errors="replace", env=env)
        open(path, "wb").write(pristine)
        md5_after = hashlib.md5(open(path, "rb").read()).hexdigest()
        print("=== %s (%s) ===" % (mid, path))
        print("rc=%s restore_md5_equal=%s" % (proc.returncode, md5_before == md5_after))
        for ln in (proc.stdout + proc.stderr).splitlines():
            s = ln.strip()
            if s.startswith("--- FAIL") or s.startswith("--- PASS") or "_test.go:" in s \
               or s.startswith("FAIL") or s.startswith("ok ") or "cannot use" in s or "declared and not used" in s or "error:" in s:
                print(ln)
        print("md5_before=%s md5_after=%s" % (md5_before, md5_after))
