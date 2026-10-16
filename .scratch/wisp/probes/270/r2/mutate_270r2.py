#!/usr/bin/env python3
"""mutate_270r2.py - the mandatory question, run OUTSIDE the repo.

Answers one thing only: after the BaseOptions half is taken out of the || at
subagent_197.go:265, does the "who rings first" sentence in the comment this leg
just rewrote still stand?  It rewrites PHYSICAL bytes of a file inside an export
tree (git archive), never the working tree, and it restores by sha256 proof.

ops:
  half-off        :65 loses its left half -> "if t.d.ParentTools == nil {"
  half-off-no-jia : same, plus the :297 backstop block deleted (去牙)
  restore         : put the pristine bytes back

Every op refuses to write unless its anchor text matches EXACTLY ONCE.
"""
import hashlib
import sys
from pathlib import Path

GUARD = "\tif t.d.BaseOptions == nil || t.d.ParentTools == nil {\n"
GUARD_HALF = "\tif t.d.ParentTools == nil {\n"

# the backstop block, verbatim: if / release() / return / }
JIA_BLOCK = ("\tif t.d.BaseOptions == nil {\n"
             "\t\trelease()\n"
             "\t\treturn Result{Text: \"宿主没有给出派生用的装配（BaseOptions/ParentTools 未接线，"
             "fail-closed：不起第二套运行时）\", IsError: true}, nil\n"
             "\t}\n")


def sha(p: Path) -> str:
    return hashlib.sha256(p.read_bytes()).hexdigest()


def main() -> int:
    if len(sys.argv) != 4:
        print("usage: mutate_270r2.py <tree-file> <pristine-file> <op>", file=sys.stderr)
        return 2
    target, pristine, op = Path(sys.argv[1]), Path(sys.argv[2]), sys.argv[3]
    text = target.read_text(encoding="utf-8")

    if op == "restore":
        target.write_bytes(pristine.read_bytes())
        print(f"restored sha256={sha(target)}")
        return 0

    if op == "half-off":
        if text.count(GUARD) != 1:
            print(f"REFUSE: guard anchor matched {text.count(GUARD)} times", file=sys.stderr)
            return 3
        text = text.replace(GUARD, GUARD_HALF)
    elif op == "half-off-no-jia":
        if text.count(GUARD) != 1 or text.count(JIA_BLOCK) != 1:
            print(f"REFUSE: anchors guard={text.count(GUARD)} jia={text.count(JIA_BLOCK)}",
                  file=sys.stderr)
            return 3
        text = text.replace(GUARD, GUARD_HALF).replace(JIA_BLOCK, "")
    else:
        print(f"REFUSE: unknown op {op!r}", file=sys.stderr)
        return 2

    target.write_text(text, encoding="utf-8")
    print(f"op={op} mutated sha256={sha(target)} lines={len(text.splitlines())}")
    return 0


if __name__ == "__main__":
    sys.exit(main())
