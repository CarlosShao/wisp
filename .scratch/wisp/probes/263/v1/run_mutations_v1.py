#!/usr/bin/env python3
# 263-v1 突变台件（Python 编排，⛔ 不用 PowerShell 编排：本腿连撞三把 PS 5.1 的尺——
# stderr 进管道变 NativeCommandError、1>/2> 与 Start-Process 取不到退码、ProcessStartInfo 属性赋值报错）。
# 每枚副本都由 tracked 字节逐字节生成，锚不命中即抛；跟踪文件全程只读，跑完再 md5 复尺。
import hashlib
import os
import subprocess
import sys
from pathlib import Path

HERE = Path(r"D:\work\workspace\projects plans\Wisp\.scratch\wisp\probes\263\v1")
MUT = HERE / "mutations"
TRACKED = Path(r"D:\work\workspace\projects plans\Wisp\scripts\slo-check.ps1")
REPO = TRACKED.parent.parent                      # 仓根
# base 直接由 git 取 HEAD 字节（不把跟踪文件的副本入库，免得留第二真相源）
BASE_BYTES = subprocess.run(["git", "-C", str(REPO), "cat-file", "blob", "HEAD:scripts/slo-check.ps1"],
                            capture_output=True, check=True).stdout
GUI = str(HERE / "bin" / "probe263v1-gui.exe")
CODE_STATE = b"    $code = $run.code\n"
PLANT_BARE = b"    $code = $LASTEXITCODE\n"
PLANT_GETVAR = (
    b"    & $WispExe slo -state $state -seconds $SecondsPerState -interval-ms 250 -out $outFile\n"
    b"    $code = (Get-Variable -Name 'LASTEXITCODE' -ValueOnly)\n"
)
PLANT_OTHERFILE = (
    b"    . (Join-Path $PSScriptRoot 'helper-v1.ps1')\n"
    b"    $code = Get-V1ExitCode -Exe $WispExe -State $state -Seconds $SecondsPerState -Out $outFile\n"
)
NAIL_CALLS = b"Test-ScriptShapeNail\nTest-ExitCodeInstrument\n"
JUDGE_LINE = b"            $pass = [bool]$reportJson.pass\n"

MUTATIONS = [
    ("copy-noop", []),                                            # 副本机制自身的对照
    ("drop-wait", [(b"        $proc.WaitForExit()\n", b"        # WAIT-REMOVED-BY-263-V1\n")]),
    ("plant-bare", [(CODE_STATE, PLANT_BARE)]),
    ("plant-getvariable", [(CODE_STATE, PLANT_GETVAR)]),
    ("delete-nails", [(NAIL_CALLS, b"# NAIL-CALLS-REMOVED-BY-263-V1\n")]),
    ("read-in-other-file", [(CODE_STATE, PLANT_OTHERFILE)]),
    ("force-state-fail", [(JUDGE_LINE, b"            $pass = $false # FORCE-RED-BY-263-V1\n")]),
    ("sample-fail", []),             # 交付字节不动判据，只把样本换成真超标形状（V1_VERDICT=fail）
    ("delete-nails-and-wait", [
        (NAIL_CALLS, b"# NAIL-CALLS-REMOVED-BY-263-V1\n"),
        (b"        $proc.WaitForExit()\n", b"        # WAIT-REMOVED-BY-263-V1\n"),
    ]),
    ("plant-in-comment", [
        (CODE_STATE, b"    # note: never fall back to $LASTEXITCODE in this loop\n" + CODE_STATE),
    ]),
]


def md5_bytes(b: bytes) -> str:
    return hashlib.md5(b).hexdigest()


def main() -> int:
    tracked_md5_before = md5_bytes(TRACKED.read_bytes())
    base_bytes = BASE_BYTES
    if md5_bytes(base_bytes) != tracked_md5_before:
        print("FATAL: HEAD 字节与工作树跟踪文件不同（工作树不干净 ⇒ 本尺作废）", file=sys.stderr)
        return 2
    rows = []
    wanted = set(sys.argv[1:])
    for name, edits in MUTATIONS:
        if wanted and name not in wanted:
            continue
        blob = base_bytes
        for anchor, repl in edits:
            cnt = blob.count(anchor)
            if cnt != 1:
                print(f"FATAL: mutation {name} 锚命中 {cnt} 次（要求恰好 1）: {anchor!r}", file=sys.stderr)
                return 2
            blob = blob.replace(anchor, repl, 1)
        copy = MUT / f"slo-check-{name}.ps1"
        copy.write_bytes(blob)
        outdir = MUT / f"out-{name}"
        outdir.mkdir(exist_ok=True)
        env = dict(os.environ)
        env.pop("V1_VERDICT", None)
        if name == "sample-fail":
            env["V1_VERDICT"] = "fail"
        proc = subprocess.run(
            ["powershell", "-NoProfile", "-ExecutionPolicy", "Bypass", "-File", str(copy),
             "-Subset", "smoke", "-SecondsPerState", "1", "-WispExe", GUI, "-OutDir", str(outdir)],
            capture_output=True, env=env, cwd=str(HERE), timeout=600)
        text = proc.stdout.decode("utf-8", "replace")
        err = proc.stderr.decode("utf-8", "replace")
        log = MUT / f"log-{name}.txt"
        log.write_text(f"=== MUTATION {name} rc={proc.returncode} ===\n{text}"
                       + (f"\n--- stderr ---\n{err}" if err.strip() else "") + "\n", encoding="utf-8")
        rows.append((name, proc.returncode, len(text.splitlines()), str(log)))
    tracked_md5_after = md5_bytes(TRACKED.read_bytes())
    print("name                  rc   console-lines  log")
    for r in rows:
        print(f"{r[0]:<20} {r[1]:>3}  {r[2]:>10}   {r[3]}")
    print(f"tracked md5 before = {tracked_md5_before}")
    print(f"tracked md5 after  = {tracked_md5_after}")
    print(f"tracked UNCHANGED  = {tracked_md5_before == tracked_md5_after}")
    return 0


if __name__ == "__main__":
    sys.exit(main())
