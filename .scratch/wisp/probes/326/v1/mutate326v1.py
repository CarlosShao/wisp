# 326-v1 mutation driver (adjudication leg, read-only on production bytes).
#
# Usage:  python mutate326v1.py <export-tree> <mode>
#   mode = pristine | A | B | D | restore | show
#
# Every mode writes ONLY into the out-of-repo `git archive` tree given as argv[1].
# The shared working tree is never a mutation target (dispatch section 0.3).
# No -overlay is used: mutations hit physical bytes on disk.
#
#  A  de-tooth        : delete the single `UPDATE session SET updated_at=?` statement
#                       line from AppendSessionTurn (AC#1 criterion 2).
#  B  anti-form        : keep the turn INSERT committed alone, move the refresh into
#                       its OWN s.write after the commit (AC#1 criterion 3, the shape
#                       the delivered final-state cases are supposed to be blind to).
#  D  die-mid-tx       : keep the refresh inside the transaction but exit the process
#                       between the two statements, gated on WISP_326V1_EXIT_MID_TX=1
#                       so the same binary can be run with and without it (AC#3
#                       "is this bench ever-passing?" probe).
#  restore             : byte-for-byte restore from the pristine snapshot saved by the
#                       first call; prints the restored sha256 and compares.
import hashlib
import os
import sys

TARGET_REL = os.path.join("internal", "memory", "dao_session.go")

UPDATE_LINE = ("\t\t_, err = tx.ExecContext(ctx, `UPDATE session SET updated_at=? WHERE id=?`, "
               "t.CreatedAt, t.SessionID)\n")
IMPORT_OLD = '\t"errors"\n\t"fmt"\n\t"strings"\n'
IMPORT_NEW = '\t"errors"\n\t"fmt"\n\t"os"\n\t"strings"\n'
TAIL_OLD = ('\tif err != nil {\n\t\treturn 0, fmt.Errorf("memory: append session turn: %w", err)\n\t}\n'
            '\treturn id, nil\n}\n')
TAIL_NEW = ('\tif err != nil {\n\t\treturn 0, fmt.Errorf("memory: append session turn: %w", err)\n\t}\n'
            '\t// 326-v1 MU-B: the touch moved OUT of the turn transaction.\n'
            '\tif werr := s.write(ctx, func(c context.Context, touch *sql.Tx) error {\n'
            '\t\t_, e := touch.ExecContext(c, `UPDATE session SET updated_at=? WHERE id=?`, '
            't.CreatedAt, t.SessionID)\n\t\treturn e\n'
            '\t}); werr != nil {\n'
            '\t\treturn 0, fmt.Errorf("memory: append session turn touch: %w", werr)\n\t}\n'
            '\treturn id, nil\n}\n')
GATE = ('\t\tif os.Getenv("WISP_326V1_EXIT_MID_TX") == "1" {\n'
        '\t\t\t// 326-v1 MU-D: die between the two statements of the SAME transaction.\n'
        '\t\t\tos.Exit(7)\n\t\t}\n')


def sha(b):
    return hashlib.sha256(b).hexdigest()


def main():
    tree, mode = sys.argv[1], sys.argv[2]
    path = os.path.join(tree, TARGET_REL)
    snap = os.path.join(tree, "..", "pristine.dao_session.go")
    current = open(path, "rb").read()

    if mode == "show":
        print("sha256=%s bytes=%d" % (sha(current), len(current)))
        return
    if mode == "restore":
        pristine = open(snap, "rb").read()
        open(path, "wb").write(pristine)
        after = open(path, "rb").read()
        print("restored sha256=%s cmp_pristine=%s" % (sha(after), "MATCH" if after == pristine else "DIFFER"))
        return

    body = current.decode("utf-8")
    base = sha(current)
    if not os.path.exists(snap):
        open(snap, "wb").write(current)
    assert UPDATE_LINE in body, "UPDATE statement line not found verbatim -- wrong tree?"
    if mode == "A":
        body = body.replace(UPDATE_LINE, "", 1)
    elif mode == "B":
        body = body.replace(UPDATE_LINE + "\t\treturn err\n", "\t\treturn err\n", 1)
        assert TAIL_OLD in body, "AppendSessionTurn tail not found"
        body = body.replace(TAIL_OLD, TAIL_NEW, 1)
    elif mode == "D":
        body = body.replace(IMPORT_OLD, IMPORT_NEW, 1)
        body = body.replace(UPDATE_LINE, GATE + UPDATE_LINE, 1)
    else:
        raise SystemExit("unknown mode " + mode)
    data = body.encode("utf-8")
    open(path, "wb").write(data)
    print("mode=%s sha_before=%s sha_after=%s bytes=%d" % (mode, base, sha(data), len(data)))
    # re-read from disk: the ruler must be about bytes, not about this script's memory
    disk = open(path, "rb").read()
    print("re-read-from-disk sha256=%s equal=%s" % (sha(disk), sha(disk) == sha(data)))


main()
