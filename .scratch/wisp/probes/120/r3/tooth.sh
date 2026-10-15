#!/bin/sh
# 120-r3 control driver: the paired readings for ticket 120 residual gap 1.
#
# Three container runs on ONE table (same snapshot dir, same image, same
# -count=1 -v ./internal/winsec/, temp = the image's own overlay /tmp):
#   teeth-on   -> the committed tree, expected all green
#   tooth-269  -> sealHandle's LINK face neutered alone
#   tooth-272  -> sealHandle's SAMEFILE face neutered alone
#
# Anchoring is by FUNCTION BODY, not by line number: the same two condition texts
# exist a second time inside sameObjectPinned (:188 and :191) in this same file, so a
# file-wide substitution would silently neuter the wrong leg - which is exactly how
# 120-r2's control could not speak for this one. The awk below only rewrites inside
# `func sealHandle`, only the first match, and refuses (exit 3) if it did not land.
#
# The rewrite is PHYSICAL, on the snapshot, and restored by sha256 (no -overlay: the
# leg under test reads the disk at run time, so an overlay that only swaps build
# inputs would read green without teeth). The repository working tree is never
# touched, and every .bak stays on disk.
set -e

REPO='D:/work/workspace/projects plans/Wisp'
S=/d/tmp/probe120r3
SNAP=$S/fixed
LOGS=$S/logs
F=internal/winsec/winsec_other.go
export MSYS_NO_PATHCONV=1
mkdir -p "$LOGS" "$S/readouts"

# mut <tooth> : prepend `false &&` to ONE condition inside func sealHandle of $SNAP/$F.
# The line number is FOUND, never assumed (grep -n by content, restricted to the
# function body), the edit is applied to that one line, and two self-proofs are printed:
# exactly one `false &&` in the whole file (so the sameObjectPinned leg next door is
# untouched) and the line it sits on, inside the function's own range.
mut() {
	tooth=$1
	case $tooth in
	269) pat='now.Mode()&os.ModeSymlink != 0 {' ;;
	272) pat='!os.SameFile(pinned, now) {' ;;
	*) echo "unknown tooth $tooth" >&2; exit 2 ;;
	esac
	start=$(grep -n '^func sealHandle(' "$SNAP/$F" | cut -d: -f1)
	end=$(awk -v s="$start" 'NR>s && /^\}/ {print NR; exit}' "$SNAP/$F")
	line=$(awk -v s="$start" -v e="$end" -v want="$pat" 'NR>s && NR<e && index($0, want)>0 {print NR; exit}' "$SNAP/$F")
	if [ -z "$line" ]; then
		echo "MUTATION DID NOT LAND: no line carrying the tooth $tooth condition inside func sealHandle ($start..$end)" >&2
		exit 3
	fi
	echo "LANDS: tooth $tooth targets line $line of $F, inside func sealHandle ($start..$end), before: $(sed -n "${line}p" "$SNAP/$F")"
	sed -i "${line}s/if /if false \\&\\& /" "$SNAP/$F"
	echo "LANDS: after: $(sed -n "${line}p" "$SNAP/$F")"
	n=$(grep -c 'false &&' "$SNAP/$F" || true)
	if [ "$n" != 1 ]; then
		echo "MUTATION NOT ISOLATED: $n lines carry false &&, want exactly 1 (the walk's leg would be neutered too)" >&2
		exit 3
	fi
	echo "LANDS: isolated - exactly $n line in the file carries false && (grep -n: $(grep -n 'false &&' "$SNAP/$F" | cut -d: -f1))"
}

# runfull <name> : whole package, -count=1 -v. Two files per run: <name>.log keeps the
# docker and go-test exit lines, <name>-body.log keeps the test binary's own output, so
# the four numbers are counted from bytes the container wrote and not from this script.
runfull() {
	name=$1
	out=$LOGS/$name.log
	body=$LOGS/$name-body.log
	rc=0
	docker run --rm -e CGO_ENABLED=0 -v "$SNAP":/work -v probe120r3build:/root/.cache/go-build \
		-w /work golang:1.27 sh -c 'go test -count=1 -v ./internal/winsec/ > /work/full.log 2>&1; rc=$?; echo "go-test rc=$rc"; exit $rc' \
		> "$out" 2>&1 || rc=$?
	cp "$SNAP/full.log" "$body"
	{
		echo "=== $name: docker-and-test rc=$rc (go-test rc line inside $out)"
		printf 'RUN=%s (=== RUN rows, t.Run subtests included)\n' "$(grep -c '=== RUN' "$body" || true)"
		printf 'top-level result rows (anchored ^---): PASS=%s FAIL=%s SKIP=%s\n' \
			"$(grep -c -- '^--- PASS' "$body" || true)" \
			"$(grep -c -- '^--- FAIL' "$body" || true)" \
			"$(grep -c -- '^--- SKIP' "$body" || true)"
		printf 'all result rows (unanchored, the 120-v1 口径: PASS=RUN when nothing fails): PASS=%s FAIL=%s SKIP=%s counted from %s\n' \
			"$(grep -c -- '--- PASS' "$body" || true)" \
			"$(grep -c -- '--- FAIL' "$body" || true)" \
			"$(grep -c -- '--- SKIP' "$body" || true)" "$body"
		echo '--- failing names:'
		grep -- '--- FAIL' "$body" || true
	} | tee -a "$out"
}

cd "$REPO"
echo "== rebuild the snapshot from HEAD ($(git rev-parse --short HEAD)) =="
rm -rf "$SNAP"
mkdir -p "$SNAP"
git archive HEAD | tar -x -C "$SNAP"
echo "go-files-in-snapshot=$(find "$SNAP" -name '*.go' | wc -l) winsec-go=$(ls "$SNAP"/internal/winsec/*.go | wc -l)"

CLEAN=$(sha256sum "$SNAP/$F" | cut -d' ' -f1)
echo "clean sha256 of $F = $CLEAN"
grep -n '^func sealHandle\|^func sameObjectPinned' "$SNAP/$F"

echo "== run 1: teeth ON (committed tree) =="
runfull teeth-on
cp "$SNAP/$F" "$SNAP/$F.bak-120r3-teeth-on"

for tooth in 269 272; do
	echo "== run: tooth $tooth out (sealHandle inline face neutered alone) =="
	cp "$SNAP/$F.bak-120r3-teeth-on" "$SNAP/$F.120r3.bak"
	mut "$tooth"
	echo "changed-lines=$(diff "$SNAP/$F.120r3.bak" "$SNAP/$F" | grep -c '^[<>]' || true)"
	diff -u "$SNAP/$F.120r3.bak" "$SNAP/$F" | grep '^[+-][+-]\|false &&' || true
	echo "mutated sha256 = $(sha256sum "$SNAP/$F" | cut -d' ' -f1)"
	runfull "tooth-$tooth"
	cp "$SNAP/$F.120r3.bak" "$SNAP/$F"
	RESTORED=$(sha256sum "$SNAP/$F" | cut -d' ' -f1)
	if [ "$RESTORED" = "$CLEAN" ]; then
		echo "RESTORE MATCHES: $RESTORED"
	else
		echo "RESTORE MISMATCH: $RESTORED vs $CLEAN" >&2
		exit 4
	fi
done

echo "== repo working tree must be untouched by the whole control =="
git status --porcelain -- internal/winsec
echo "repo sha256 of $F = $(sha256sum "$REPO/$F" | cut -d' ' -f1) (compare with the clean sha above)"
