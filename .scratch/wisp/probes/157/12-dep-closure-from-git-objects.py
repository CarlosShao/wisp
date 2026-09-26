#!/usr/bin/env python3
"""票 157 实现程 · 第 8 笔的当轮尺：从 git 对象（不读工作树）算 internal/agent 的依赖闭包。

两把尺分开算：
  R-dep      = 非测试文件的 import 传递闭包（相当于 go list -deps ./internal/agent）
  R-testdep  = 再加 _test.go（含包内测试与 agent_test 外测）（相当于 go list -test -deps）
只用 git ls-tree / git show 取版，绝不读脏工作树。
"""
import re, subprocess, sys

ANCHOR = sys.argv[1] if len(sys.argv) > 1 else "6de3d1c5"
MOD = "github.com/CarlosShao/wisp/"
IMP = re.compile(r'^\t(?:_ )?(?:"([\w./-]+)")', re.M)
BLOCK = re.compile(r'^import\s*\(([\s\S]*?)^\)', re.M)
ONES = re.compile(r'^import\s+"([\w./-]+)"', re.M)


def show(path, rev=ANCHOR):
    r = subprocess.run(["git", "show", f"{rev}:{path}"], capture_output=True, text=True, encoding="utf-8", errors="replace")
    return r.stdout if r.returncode == 0 else ""


def rev_tree(pkg_dir):
    return f"{ANCHOR}:{pkg_dir}"


def list_blobs(pkg_dir):
    # 只取该包目录**本级**的 .go（不带 -r，否则会把 internal/agent/approval 的算进 internal/agent）
    out = subprocess.run(["git", "ls-tree", "--name-only", rev_tree(pkg_dir)],
                         capture_output=True, text=True, encoding="utf-8", errors="replace").stdout.splitlines()
    return [f"{pkg_dir}/{p}" for p in out if p.endswith(".go")]


def imports_of(pkg_dir, include_tests):
    got = set()
    for p in list_blobs(pkg_dir):
        is_test = p.endswith("_test.go")
        if is_test and not include_tests:
            continue
        src = show(p)
        for m in BLOCK.finditer(src):
            for mm in IMP.finditer(m.group(1)):
                got.add(mm.group(1))
        for mm in ONES.finditer(src):
            got.add(mm.group(1))
    return {i[len(MOD):] for i in got if i.startswith(MOD)}


def closure(root_pkg, include_tests):
    seen, queue = {root_pkg}, [root_pkg]
    while queue:
        pkg = queue.pop()
        for dep in sorted(imports_of(pkg, include_tests)):
            if dep not in seen:
                seen.add(dep)
                queue.append(dep)
    seen.discard(root_pkg)
    return seen


for label, inc in (("R-dep（非测试）", False), ("R-testdep（含 _test.go）", True)):
    c = sorted(closure("internal/agent", inc))
    print(f"### {label} 锚点 {ANCHOR}：{len(c)} 枚包")
    for p in c:
        print("   ", p)
    print("  internal/tools 在闭包里？", "YES" if "internal/tools" in c else "NO")
    print()

print("### 直接依赖（只第一层，internal/agent 目录本体，非测试）")
print(sorted(imports_of("internal/agent", False)))
print("### 直接依赖（internal/agent 目录本体，含 _test.go）")
print(sorted(imports_of("internal/agent", True)))
print("### internal/agent/approval 是不是 internal/agent 的依赖（非测试第一层）")
print("approval in deps:", "internal/agent/approval" in imports_of("internal/agent", False))
