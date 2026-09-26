"""Ticket 153 AC#3 mutation harness (out-of-repo snapshot only).

Every op prints a landing proof, so no reading is ever taken on a mutation that
did not actually land. Restore is a byte-exact copy back from the pristine set
captured by run-ac3-matrix.sh.

ops:
  m5              compress.go: `if rep.Ran {`        -> `if c.Need(hist) {`
  drop-attr       remove the `task` append block from the trace
  always-attr     append `task` unconditionally (placeholder laundering)
  fake-id         traceTaskID returns one constant uuid-shaped string
  unloop          loop.go stops tagging the ctx (withTraceTask(ctx, taskID) -> ctx)
  drop-test:NAME  delete one top-level test function from a _test.go file
  restore         copy every file back from the pristine dir
"""
import shutil
import sys
from pathlib import Path

AGENT = "internal/agent"
COMPRESS = f"{AGENT}/compress.go"
LOOP = f"{AGENT}/loop.go"
TESTFILE = f"{AGENT}/compress_trace_test.go"

ATTR_BLOCK = (
    '\t\tif task := traceTaskID(ctx); task != "" {\n'
    '\t\t\tattrs = append(attrs, "task", task)\n'
    "\t\t}\n"
)
FAKE_BODY = (
    "func traceTaskID(ctx context.Context) string {\n"
    "\t_ = ctx\n"
    '\treturn "00000000-0000-4000-8000-000000000000" // MUT fabricated id\n'
    "}\n"
)


def read(p: Path) -> str:
    with p.open(encoding="utf-8", newline="") as f:
        return f.read()


def write(p: Path, s: str) -> None:
    with p.open("w", encoding="utf-8", newline="") as f:
        f.write(s)


def sub(root: Path, rel: str, old: str, new: str, proof: str) -> int:
    p = root / rel
    src = read(p)
    if src.count(old) != 1:
        print(f"REFUSED {rel}: pattern occurs {src.count(old)} times, want 1")
        return 2
    write(p, src.replace(old, new, 1))
    print(f"LANDED ({rel}): {proof}")
    return 0


def drop_test(root: Path, name: str) -> int:
    for rel in (TESTFILE, f"{AGENT}/compress_test.go"):
        p = root / rel
        if not p.exists():
            continue
        lines = read(p).splitlines(keepends=True)
        start = next((i for i, ln in enumerate(lines) if ln.startswith(f"func {name}(")), None)
        if start is None:
            continue
        end = next((i for i in range(start, len(lines)) if lines[i] == "}\n"), None)
        if end is None:
            print(f"REFUSED: no closing brace for {name}")
            return 2
        write(p, "".join(lines[:start] + lines[end + 1:]))
        print(f"LANDED ({rel}): removed {name} (lines {start + 1}-{end + 1})")
        return 0
    print(f"REFUSED: {name} not found")
    return 2


def main() -> int:
    op, root, pristine = sys.argv[1], Path(sys.argv[2]), Path(sys.argv[3])
    if op == "restore":
        for f in pristine.rglob("*"):
            if f.is_file():
                dst = root / f.relative_to(pristine)
                dst.parent.mkdir(parents=True, exist_ok=True)
                shutil.copyfile(f, dst)
        print(f"restored {sum(1 for f in pristine.rglob('*') if f.is_file())} files from {pristine}")
        return 0
    if op == "m5":
        return sub(root, COMPRESS, "\tif rep.Ran {\n", "\tif c.Need(hist) {\n",
                   "guard `if rep.Ran` -> `if c.Need(hist)`")
    if op == "drop-attr":
        return sub(root, COMPRESS, ATTR_BLOCK, "", "task attribute append removed")
    if op == "always-attr":
        return sub(root, COMPRESS, ATTR_BLOCK,
                   '\t\tattrs = append(attrs, "task", traceTaskID(ctx)) // MUT no empty check\n',
                   "task appended even when the caller tagged nothing")
    if op == "fake-id":
        p = root / COMPRESS
        src = read(p)
        head = src.index("func traceTaskID(")
        tail = src.index("}\n", head) + 2
        write(p, src[:head] + FAKE_BODY + src[tail:])
        print("LANDED (compress.go): traceTaskID now returns a constant id")
        return 0
    if op == "unloop":
        return sub(root, LOOP, "l.comp.Compress(withTraceTask(ctx, taskID), hist)",
                   "l.comp.Compress(ctx, hist)", "loop stopped tagging the ctx")
    if op.startswith("drop-test:"):
        return drop_test(root, op.split(":", 1)[1])
    print("unknown op")
    return 2


if __name__ == "__main__":
    sys.exit(main())
