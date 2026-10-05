import json
import pathlib
import sys

ROOT = pathlib.Path("D:/work/workspace/projects plans/Wisp")
SRC = ROOT / "cmd/wisp/firstrun.go"
OUT = ROOT / ".scratch/wisp/probes/257/r2/mut"
text = SRC.read_text(encoding="utf-8")
lines = text.split("\n")


def sub(content, needle, repl, name):
    n = content.count(needle)
    if n != 1:
        sys.exit("MUTATION SETUP FAIL %s: needle found %d times, want 1: %r" % (name, n, needle))
    return content.replace(needle, repl, 1)


def block_replace(start_needle, end_needle, repl_lines, name):
    si = None
    for i, ln in enumerate(lines):
        if start_needle in ln:
            si = i
            break
    if si is None:
        sys.exit("MUTATION SETUP FAIL %s: start anchor missing" % name)
    # walk back to the fmt.Fprint* line that opens this statement
    while "fmt.Fprint" not in lines[si]:
        si -= 1
    ei = None
    for j in range(si, len(lines)):
        if end_needle in lines[j]:
            ei = j
            break
    if ei is None:
        sys.exit("MUTATION SETUP FAIL %s: end anchor missing" % name)
    return lines[:si] + repl_lines + lines[ei + 1:]


mutants = {}

# M1: the first hand-added item's section name becomes the shape A543 recorded
# as fatal (a top-level [providers.x] hits parse.go's DisallowUnknownFields).
mutants["M1"] = sub(text, "第一样＝一节 [llm.providers.<名>]", "第一样＝一节 [providers.<名>]", "M1")

# M2: the three refusal reasons folded into one shared verdict (AC#2's ban).
mutants["M2"] = "\n".join(block_replace(
    "写不进去的时候有三种原因", "\t\tcfgPath)",
    ['\tfmt.Fprint(stderr, "wisp run: 写不进去的时候就是配置未生效。\\n")'], "M2"))

# M3: the three process shapes folded into one restart sentence (boundary 1 / A560).
mutants["M3"] = "\n".join(block_replace(
    "改完什么时候才算用上", "只说已录入还是没录入。\\n\")",
    ['\tfmt.Fprint(stderr, "wisp run: 改完重启就好。\\n")'], "M3"))

# M4: the guidance leaks out of the receipt and into the generated config.toml
# (ticket 198's pin at firstrun_198_test.go:215 is the tooth here).
m4 = list(lines)
anchor = None
for i, ln in enumerate(lines):
    if "config.toml 首建未写成" in ln:
        anchor = i
        break
if anchor is None:
    sys.exit("MUTATION SETUP FAIL M4: anchor missing")
close = anchor + 1
if lines[close].strip() != "}":
    sys.exit("MUTATION SETUP FAIL M4: expected the if-block to close right after the anchor, got %r" % lines[close])
m4[close + 1:close + 1] = [
    "\tif raw, rerr := os.ReadFile(cfgPath); rerr == nil {",
    "\t\t_ = os.WriteFile(cfgPath, append(raw, []byte(\"\\n[llm.providers.deepseek]\\napi_key_ref = 'env:M4LEAK'\\n\")...), 0o600)",
    "\t}",
]
mutants["M4"] = "\n".join(m4)

# M5: the role item taught as "append a second section" - a TOML duplicate table.
mutants["M5"] = sub(
    text,
    "第三样＝就地填已有的 [llm.roles.chat] 那一节",
    "第三样＝在这份文件末尾再追加一节 [llm.roles.chat]",
    "M5")

# M6: the credential paragraph starts quoting a value-shaped literal (AC#3).
mutants["M6"] = sub(
    text,
    '这一页任何时候都不回显密钥的值，只说已录入还是没录入。\\n")',
    '这一页任何时候都不回显密钥的值，写法形如 api_key = \\"sk-MutationCanaryNotARealKey\\"。\\n")',
    "M6")

for name, content in mutants.items():
    if content == text:
        sys.exit("MUTATION SETUP FAIL %s: mutant identical to source" % name)
    d = OUT / name
    d.mkdir(parents=True, exist_ok=True)
    mp = d / "firstrun.go"
    mp.write_text(content, encoding="utf-8")
    overlay = {"Replace": {str(SRC): str(mp)}}
    (d / "overlay.json").write_text(json.dumps(overlay), encoding="utf-8")
    print("built %s -> %s (%d bytes vs %d)" % (name, mp, len(content), len(text)))
