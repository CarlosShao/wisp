#!/usr/bin/env bash
# 327-v1 verdict leg: targeted-mutation driver.
#
#   usage: bash v1-mutate.sh <m1|m2|m3|m4|m5|m6|m7> <tree-dir> <roster-base-file>
#
# Everything happens in an OUT-OF-REPO export tree. The repository's tracked
# bytes are never written. The mutated file is restored in a trap (the finally),
# so an interrupted run cannot leave a mutation behind; the restored sha256 is
# printed and compared against the pre-mutation sha256.
#
#   m1 = writer side only: delete the `case isDisplayWashCf(r):` branch of
#        scrubForDisplay's switch (the judge keeps its Cf term)
#   m2 = judge side only: delete the `|| isDisplayWashCf(r)` term of
#        needsDisplayWash (the writer branch stays)
#   m3 = both rulers blind: isDisplayWashCf returns false
#   m4 = roster too wide: isDisplayWashCf = unicode.Is(unicode.Cf, r)
#   m5 = ticket AC#4's own attack: the verdict passes anything whose path name
#        says 已允许 (this edits internal/risk inside the tree copy)
#   m6 = a silent widening nobody ruled: add U+206B to the roster
#   m7 = AC#3's inverse shape: the Cf branch DELETES the rune instead of escaping it
set -u
MUT="$1"; TREE="$2"; BASE_ROSTER="$3"
SRC="$TREE/internal/panel/display_scrub_216.go"
RISK="$TREE/internal/risk/assessor.go"
LOGDIR="$(dirname "$BASE_ROSTER")"
OUT="$LOGDIR/mut-$MUT.log"

restore() {
  rc=$?
  cp "$BACKUP" "$SRC"
  [ -n "${RISKBACKUP:-}" ] && cp "$RISKBACKUP" "$RISK"
  NOW="$(sha256sum "$SRC" | cut -d' ' -f1)"
  echo "restored_sha256=$NOW"
  echo "pre_mutation_sha256=$BEFORE"
  OK=1
  [ "$NOW" = "$BEFORE" ] || OK=0
  if [ -n "${RISKBACKUP:-}" ]; then
    RNOW="$(sha256sum "$RISK" | cut -d' ' -f1)"
    echo "restored_risk_sha256=$RNOW"
    echo "pre_mutation_risk_sha256=$RISKBEFORE"
    [ "$RNOW" = "$RISKBEFORE" ] || OK=0
  fi
  if [ "$OK" = 1 ]; then echo "RESTORE_OK"; else echo "RESTORE_FAIL"; fi
  trap - EXIT
  exit $rc
}

BEFORE="$(sha256sum "$SRC" | cut -d' ' -f1)"
BACKUP="$(mktemp /d/tmp/gotmp327v1/backup.XXXXXX)"
cp "$SRC" "$BACKUP"
if [ "$MUT" = "m5" ]; then
  RISKBEFORE="$(sha256sum "$RISK" | cut -d' ' -f1)"
  RISKBACKUP="$(mktemp /d/tmp/gotmp327v1/backup.XXXXXX)"
  cp "$RISK" "$RISKBACKUP"
fi
trap restore EXIT INT TERM

python - "$MUT" "$SRC" "$RISK" <<'PY'
import hashlib, io, sys
mut, src, risk = sys.argv[1], sys.argv[2], sys.argv[3]

def load(p):
    return io.open(p, encoding="utf-8", newline="").read()

def store(p, s):
    io.open(p, "w", encoding="utf-8", newline="").write(s)

if mut in ("m1", "m2", "m3", "m4", "m6", "m7"):
    s = load(src)
    if mut == "m1":
        start = s.index("\t\tcase isDisplayWashCf(r):\n")
        stop = s.index("\t\tdefault:\n", start)
        removed = s[start:stop]
        assert removed.count("\n") == 7, removed
        assert "writeControlEscape(&b, r)" in removed and "i += size" in removed
        s2 = s[:start] + s[stop:]
        print("m1_removed_lines=%d" % removed.count("\n"))
    elif mut == "m2":
        old = "if (r == utf8.RuneError && size == 1) || unicode.IsControl(r) || isDisplayWashCf(r) {"
        new = "if (r == utf8.RuneError && size == 1) || unicode.IsControl(r) {"
        assert s.count(old) == 1
        s2 = s.replace(old, new)
        print("m2_removed_term=|| isDisplayWashCf(r) (one term, one line)")
    elif mut == "m3":
        head = "func isDisplayWashCf(r rune) bool {\n"
        start = s.index(head) + len(head)
        stop = s.index("\n}\n", start)
        s2 = s[:start] + "\n\treturn false" + s[stop:]
        print("m3_replaced_body_lines=%d" % (s[start:stop].count("\n")))
    elif mut == "m4":
        head = "func isDisplayWashCf(r rune) bool {\n"
        start = s.index(head) + len(head)
        stop = s.index("\n}\n", start)
        s2 = s[:start] + "\n\treturn unicode.Is(unicode.Cf, r)" + s[stop:]
        print("m4_replaced_body_lines=%d" % (s[start:stop].count("\n")))
    elif mut == "m7":
        # AC#3's inverse shape: delete the Cf rune instead of escaping it.
        old = "\t\tcase isDisplayWashCf(r):\n"
        start = s.index(old)
        stop = s.index("\t\tdefault:\n", start)
        branch = s[start:stop]
        assert "\t\t\twriteControlEscape(&b, r)\n" in branch
        s2 = s[:start] + branch.replace("\t\t\twriteControlEscape(&b, r)\n", "") + s[stop:]
        print("m7_removed_line=writeControlEscape(&b, r) from the Cf branch only (deletion, not escape)")
    elif mut == "m6":
        old = "\t\t0x2066, 0x2067, 0x2068, 0x2069:\n"
        new = "\t\t0x2066, 0x2067, 0x2068, 0x2069, 0x206B:\n"
        assert s.count(old) == 1
        s2 = s.replace(old, new)
        print("m6_added_rune=U+206B (not ruled, not among the seven named exclusions)")
    assert s2 != s
    store(src, s2)
elif mut == "m5":
    s = load(risk)
    old = "\treturn fuse(contribs)\n}"
    new = ("\tdecision := fuse(contribs)\n"
           "\tfor _, p := range facts.Paths {\n"
           "\t\tif strings.Contains(p, \"\\u5df2\\u5141\\u8bb8\") {\n"
           "\t\t\treturn Decision{Level: L0, Reason: \"m5: the name says it is allowed\"}\n"
           "\t\t}\n"
           "\t}\n"
           "\treturn decision\n}")
    assert s.count(old) == 1
    store(risk, s.replace(old, new))
    print("m5_mutated=internal/risk/assessor.go Assess(): a path containing U+5DF2 U+5141 U+8BB8 decides L0")
else:
    raise SystemExit("unknown mutation")
PY

echo "pre_mutation_sha256=$BEFORE"
echo "mutated_sha256=$(sha256sum "$SRC" | cut -d' ' -f1)"
[ "$MUT" = "m5" ] && echo "mutated_risk_sha256=$(sha256sum "$RISK" | cut -d' ' -f1)"

cd "$TREE" || exit 1
export GOTMPDIR=/d/tmp/gotmp327v1
export PATH="$TREE/third_party/sherpa-onnx:$TREE/build:$PATH"
go test -v -count=1 -timeout 30m ./internal/panel/ > "$OUT" 2>&1
GRC=$?
L=$(grep -c '^=== RUN' "$OUT"); P=$(grep -c '^--- PASS' "$OUT"); F=$(grep -c '^--- FAIL' "$OUT"); S=$(grep -c '^--- SKIP' "$OUT")
echo "go_test_rc=$GRC RUN=$L PASS_TOP=$P FAIL_TOP=$F SKIP=$S"
grep '^--- FAIL' "$OUT" | sed 's/ (.*//' | sort > "$LOGDIR/roster-$MUT.txt"
echo "=== red roster minus base roster (mutations introduced) ==="
comm -13 "$BASE_ROSTER" "$LOGDIR/roster-$MUT.txt"
echo "=== base roster minus this roster (greens this mutation broke-and-then-... n/a) ==="
comm -23 "$BASE_ROSTER" "$LOGDIR/roster-$MUT.txt"
echo "=== verbatim first error lines of the 327 tests ==="
grep -nE '^\s+display_scrub_216_test\.go:[0-9]+:' "$OUT" | head -12
