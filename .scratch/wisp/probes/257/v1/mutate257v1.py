import hashlib
import os
import re
import subprocess
import sys
import time

ROOT = r"D:/work/workspace/projects plans/Wisp"
TARGET = os.path.join(ROOT, "cmd", "wisp", "firstrun.go")
OUTDIR = os.path.join(ROOT, ".scratch", "wisp", "probes", "257", "v1")
RUN_RE = "TestTicket257R2|TestTicket198"

# pristine copy comes from HEAD, not from the work tree, so a stray edit cannot
# become the baseline.
head_blob = subprocess.run(
    ["git", "-C", ROOT, "cat-file", "blob", "HEAD:cmd/wisp/firstrun.go"],
    capture_output=True, check=True).stdout
with open(os.path.join(OUTDIR, "firstrun.pristine"), "wb") as fh:
    fh.write(head_blob)
BASE_MD5 = hashlib.md5(head_blob).hexdigest()

MUTANTS = [
    ("M1-guidance-never-reaches-operator", [
        ('fmt.Fprint(stderr,\n\t\t"wisp run: 上面那句模型只是第一样。',
         'fmt.Fprint(io.Discard,\n\t\t"wisp run: 上面那句模型只是第一样。'),
    ]),
    ("M2-nonpreset-clause-deleted", [
        ("非预设名必须自己写 protocol，否则这份文件加载不过",
         "（257-v1 M2 把这一支抹掉了）"),
    ]),
    ("M3-guidance-leaks-into-generated-file", [
        ('\t// 说实话的回执（票面"要建什么"2 的前半',
         '\tif raw, rErr := os.ReadFile(cfgPath); rErr == nil {\n'
         '\t\t_ = os.WriteFile(cfgPath, append(raw, []byte("\\n[llm.providers.t257v1leak]\\n")...), 0o600)\n'
         '\t}\n'
         '\t// 说实话的回执（票面"要建什么"2 的前半'),
    ]),
    ("M4-three-refusal-causes-folded-into-one", [
        ('"  第 2 种拒因：行不存在——那一页改不了服务商与模型的存在性，去这份文件里手加上面那三样，"',
         '"  配置未生效——去这份文件里看看，"'),
        ('"  第 3 种拒因：校验不过——行在，但那个值过不了这份 schema 的校验"',
         '"  配置未生效——去这份文件里看看，值"'),
    ]),
    ("M5-receipt-echoes-a-credential-value", [
        ('"wisp run: 凭据这一格只有引用会进这份文件：改 api_key_ref 换的是名字不是密钥本身，"',
         '"wisp run: 凭据 api_key = \\"sk-t257v1mutantvalue\\" 这一格只有引用会进这份文件：改 api_key_ref 换的是名字不是密钥本身，"'),
    ]),
]

env = dict(os.environ)
env["PATH"] = (os.path.join(ROOT, "third_party", "sherpa-onnx") + os.pathsep
               + os.path.join(ROOT, "build") + os.pathsep + env["PATH"])
env["GOFLAGS"] = ""

print("BASELINE md5 %s  worktree-equals-head=%s" % (
    BASE_MD5,
    hashlib.md5(open(TARGET, "rb").read()).hexdigest() == BASE_MD5))

for name, edits in MUTANTS:
    src = head_blob.decode("utf-8")
    ok = True
    for needle, repl in edits:
        n = src.count(needle)
        if n != 1:
            print("%s ABORT: needle hits %d times (want exactly 1): %r" % (name, n, needle[:40]))
            ok = False
            break
        src = src.replace(needle, repl, 1)
    if not ok:
        continue
    with open(TARGET, "w", encoding="utf-8", newline="\n") as fh:
        fh.write(src)
    start = time.strftime("%Y-%m-%d %H:%M:%S")
    proc = subprocess.run(
        ["go", "test", "-count=1", "-v", "-run", RUN_RE, "./cmd/wisp/"],
        cwd=ROOT, capture_output=True, env=env)
    end = time.strftime("%Y-%m-%d %H:%M:%S")
    out = (proc.stdout + proc.stderr).decode("utf-8", "replace")
    log = os.path.join(OUTDIR, "mut-%s.log" % name)
    with open(log, "w", encoding="utf-8", newline="\n") as fh:
        fh.write(out)
    # restore immediately, then prove the restore
    with open(TARGET, "wb") as fh:
        fh.write(head_blob)
    now = hashlib.md5(open(TARGET, "rb").read()).hexdigest()
    reds = re.findall(r"^--- FAIL: (\S+)", out, re.M)
    prints = re.findall(r"^(    \s*\S*_test\.go:\d+: .*(?:RED|FAIL|refusal|never says).*)$", out, re.M)
    print("=== %s  %s -> %s  rc=%d  reds=%d  restored_md5_ok=%s" % (
        name, start, end, proc.returncode, len(reds), now == BASE_MD5))
    print("    names: %s" % ", ".join(sorted(set(reds))))
    for line in prints[:6]:
        print("    %s" % line.strip()[:300])
    print("    log: %s" % log)

print("FINAL md5 %s equals baseline: %s" % (
    hashlib.md5(open(TARGET, "rb").read()).hexdigest(),
    hashlib.md5(open(TARGET, "rb").read()).hexdigest() == BASE_MD5))
