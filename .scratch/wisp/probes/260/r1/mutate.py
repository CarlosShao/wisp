# 260-r1 mutation applier (used by mutation-driver.sh).
#   python mutate.py M1 apply   -> plant the defect this ticket is about
#   python mutate.py M1 undo    -> not used; the driver restores from a byte-identical backup
import io
import sys

PATH = "internal/ball/hotkey_windows.go"

M1_OLD = "\treturn cancelBorrow{Binding: bind, Acc: acc}\n"
M1_NEW = ("\t_ = acc // MUTATION-M1: the configured accelerator is dropped, the pre-260 hard-coded borrow\n"
          "\treturn cancelBorrow{Binding: escBorrowBinding, Acc: escBorrowAcc()}\n")

M2_OLD = "\t\tAcc:     b.Acc,\n"
M2_NEW = "\t\tAcc:     escBorrowAcc(), // MUTATION-M2: receipt and registration spelled twice\n"

M3_OLD = ('\t\t\tFallback: fmt.Errorf("the configured cancel binding %q is not a valid key combination (%w); "+\n'
          '\t\t\t\t"the default Esc was borrowed for this card instead - fix it in [hotkey]", bind, err),\n')
M3_NEW = "\t\t\t// MUTATION-M3: the parse failure is swallowed, no notice is carried (P6 not landed)\n"

CASES = {"M1": (M1_OLD, M1_NEW), "M2": (M2_OLD, M2_NEW), "M3": (M3_OLD, M3_NEW)}

tag = sys.argv[1]
old, new = CASES[tag]
src = io.open(PATH, encoding="utf-8").read()
if src.count(old) != 1:
    raise SystemExit("%s: anchor found %d times, expected 1" % (tag, src.count(old)))
io.open(PATH, "w", encoding="utf-8", newline="\n").write(src.replace(old, new))
print("%s applied" % tag)
