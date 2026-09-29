# AC#5 of ticket 234: relocate 10 duplicate denominators that sit in CLOSED tickets
# but are owned by STILL-OPEN tickets, by turning each cell's first line into a
# pointer line in the shape the repo already set today (ticket 80 :54, 89 :348, 118 :59).
#
# Rule honoured: the original sentence is kept VERBATIM after the marker, so the only
# changed characters are the line's checkbox prefix. Continuation lines of each cell
# are never touched. No `-done` file is renamed. No checkbox that is not in this list
# is flipped.
import glob
import io

PREFIX = ("- ⛔ **（2026-09-29 18:5x `gate-rerun-1` 处置＝账在仍开放的票名下，本格不再占结案票的分母；"
          "原句逐字保留、续行未动。判性来源＝只读腿 `undone-28-1` 的\"丙\"行＋派单票 234 AC#5，"
          "台账 `A446`；读数与 sha 请引接收票那一格，别引本行。）** ")

# (issue number, line, receiving open ticket + where the account lives)
CELLS = [
    ("07", 54, "票 64 `:72`（ball 交互档，64 开放）"),
    ("07", 59, "票 64 `:44`/`:59`/`:60`（热键三枚，64 开放）"),
    ("07", 61, "票 64 `:83`（副屏那一格至今未勾，64 开放）"),
    ("11", 53, "票 12 `:120`/`:123`/`:128`（正式移交代号 A11，12 开放）"),
    ("63", 46, "票 12 `:109`/`:112`/`:116-119`（正式移交代号 A8，12 开放）"),
    ("80", 43, "票 230 `:22`（AC#3 交接要么落地要么作废，230 开放）"),
    ("80", 50, "票 230 `:22`（同一格逐字写\"把票 80 那两格转进本票\"，230 开放）"),
    ("97", 53, "票 230 `:21`（AC#2 那句注释改成实话，230 开放）"),
    ("115", 44, "票 230 `:19`（AC#1 具名要 115 名下的裁决表，230 开放）"),
    ("92", 52, "票 114 `:50`（AC#5 不越界，114 开放）"),
]


def find(num):
    hits = sorted(glob.glob(".scratch/wisp/issues/%s-*.md" % num))
    if len(hits) != 1:
        raise SystemExit("ABORT: %s matches %d files" % (num, len(hits)))
    return hits[0]


def unchecked(lines):
    return sum(1 for x in lines if x.startswith("- [ ] "))


files = {}
for num, ln, target in CELLS:
    f = find(num)
    if f not in files:
        files[f] = io.open(f, encoding="utf-8").read().split("\n")
    lines = files[f]
    old = lines[ln - 1]
    if not old.startswith("- [ ] "):
        raise SystemExit("ABORT %s:%d is not an unchecked cell first line: %r" % (f, ln, old[:40]))
    lines[ln - 1] = PREFIX + old[len("- [ ] "):] + " ⇒ **指向**＝" + target

for f, lines in files.items():
    io.open(f, "w", encoding="utf-8", newline="\n").write("\n".join(lines))
    print("REWROTE %s lines=%d" % (f, len(lines)))

print("=== receiver side, counted AFTER (ruler: grep -c -- '^- \\[ \\]') ===")
for num in ("64", "12", "230", "114"):
    f = find(num)
    ls = io.open(f, encoding="utf-8").read().split("\n")
    print("%s un=%d chk=%d" % (num, unchecked(ls), sum(1 for x in ls if x.startswith("- [x] "))))
