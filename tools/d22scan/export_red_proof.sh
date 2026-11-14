#!/bin/sh
# tools/d22scan/export_red_proof.sh - ticket 71 AC#5's "one real red", made
# re-runnable instead of retold.
#
# WHAT THIS ANSWERS, AND WHY IT IS A SCRIPT
# AC#5 says: hook the gate into CI and produce ONE REAL RED - cause a violation,
# show the red line naming the file, roll back, show green again. A sentence in a
# report saying "the gate can go red, see ticket 71's unit tests" is exactly the
# evidence shape this ticket was opened against: a locally reproducible red proves
# the SCANNER works, not that the CI STEP produces a red a human can read. And an
# agent here may not push (the ticket's constraint), so "wait for a GitHub run" is
# unavailable by rule, not by mood. This script closes that gap the only way the
# disk allows: it exports one commit to a directory OUTSIDE the repository, runs
# the EXACT command CI's lint job runs (`sh scripts/d22scan.sh`), plants one
# violation, and requires the run's own output to name the planted file and to
# move that scope's self-reported number by exactly one. Then it removes the file
# and requires the tree hash to return to the pre-plant value.
#
# WHY THE EXPORT, AND NOT THIS WORKING TREE
# The working tree is shared with other agents mid-flight, so a red here could be
# somebody else's uncommitted edit - the claim would be unfalsifiable. A commit
# exported to a scratch directory is the byte set CI checks out. It has no .git,
# so the scanner announces on stderr that gitignore rules cannot be applied
# (A207's loud direction) and scans every path; that is a property of the export,
# printed in the log, not hidden here.
#
# WHY A PINNED REVISION (-r), AND NOT "HEAD"
# HEAD moves under this leg: three commits by other agents landed during one
# sitting. An export hashed against one revision and compared against another
# produces a "missing files" reading that looks like an extraction bug and is
# actually a moving repository. So the revision is resolved once, printed, used
# for both the export and the completeness check, and can be pinned by the caller.
#
# WHAT IT DOES NOT CLAIM: no GitHub Actions run is claimed or implied (that needs
# a push this leg may not do). The claim is exactly "the command CI's lint step
# runs, executed on an out-of-repo export of revision X, exits 1 and names the
# planted file in its own output, and returns to exit 0 with an identical tree
# hash once the file is removed".
#
# Usage:
#   sh tools/d22scan/export_red_proof.sh [-r <rev>] [-C <export-dir>] [-k]
#     -r <rev>  revision to export (default HEAD; pass a full/short sha to pin)
#     -C <dir>  extract there instead of a fresh temp dir (must be OUTSIDE the
#               repository and empty/absent)
#     -k        keep the export directory and logs on success (they are kept on
#               failure regardless, so a red can be inspected)
set -eu

here=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
root=$(CDPATH= cd -- "$here/../.." && pwd)

rev_arg=HEAD
export_dir=
keep=
while [ $# -gt 0 ]; do
    case "$1" in
        -r)
            [ $# -ge 2 ] || { echo "export_red_proof.sh: -r needs a revision" >&2; exit 2; }
            rev_arg=$2
            shift 2
            ;;
        -C)
            [ $# -ge 2 ] || { echo "export_red_proof.sh: -C needs a directory" >&2; exit 2; }
            export_dir=$2
            shift 2
            ;;
        -k)
            keep=1
            shift
            ;;
        -h | --help)
            sed -n '2,55p' "$0"
            exit 0
            ;;
        *)
            echo "export_red_proof.sh: unknown argument: $1" >&2
            exit 2
            ;;
    esac
done

die() {
    echo "export_red_proof.sh: $*" >&2
    exit 2
}

claim_failed() {
    echo "export-red-proof: FAIL - $*" >&2
    exit 1
}

# The tools whose absence would make this script lie in a different way (an unrun
# tar still leaves a directory that looks like an export).
for tool in git tar sha256sum go find sort xargs grep sed awk; do
    command -v "$tool" >/dev/null 2>&1 || die "needs $tool on PATH"
done

[ -f "$root/go.mod" ] || die "no go.mod at the derived repo root $root - refusing"
[ -f "$root/scripts/d22scan.sh" ] || die "no scripts/d22scan.sh at $root - that IS the CI step"

rev=$(git -C "$root" rev-parse "$rev_arg") || die "cannot resolve revision $rev_arg in $root"
rev_short=$(git -C "$root" rev-parse --short "$rev")

if [ -z "$export_dir" ]; then
    export_dir=${TMPDIR:-/tmp}/wisp-71-export-$$
fi
case "$export_dir" in
    "$root" | "$root"/*)
        die "export dir $export_dir is inside the repository ($root): the whole point is a tree outside the agent's working copy"
        ;;
esac
if [ -e "$export_dir" ]; then
    if [ -n "$(ls -A "$export_dir" 2>/dev/null)" ]; then
        die "export dir $export_dir already exists and is not empty - refusing to mix a commit into someone else's tree"
    fi
else
    mkdir -p "$export_dir"
fi

base_log=$export_dir.base.log
red_log=$export_dir.red.log
restore_log=$export_dir.restore.log

# ---- the export --------------------------------------------------------------
# Two routes, because "git archive" is fine but the EXTRACTOR is where Windows
# loses non-ASCII names: msys GNU tar refuses some of .scratch/'s filenames
# ("Invalid empty pathname"), while bsdtar reads the zip form correctly. The zip
# route is tried first and the tar stream is the fallback; whichever ran, the
# completeness check below decides whether the export is allowed to be used.
echo "export-red-proof: exporting $rev_short ($rev) to $export_dir"
extracted=no
zip_file=$export_dir.zip
extractor=
if command -v unzip >/dev/null 2>&1; then
    extractor=unzip
elif tar --version 2>/dev/null | grep -q bsdtar; then
    extractor=bsdtar
elif [ -x "$SYSTEMROOT/System32/tar.exe" ]; then
    extractor=$SYSTEMROOT/System32/tar.exe
fi
if [ -n "$extractor" ]; then
    if git -C "$root" archive --format=zip -o "$zip_file" "$rev"; then
        case "$extractor" in
            unzip) (cd "$export_dir" && unzip -q "$zip_file") ;;
            *) "$extractor" -xf "$zip_file" -C "$export_dir" ;;
        esac && extracted=yes
        rm -f "$zip_file"
    fi
fi
if [ "$extracted" != yes ]; then
    echo "export-red-proof: zip route unavailable or failed; falling back to the tar stream"
    git -C "$root" archive --format=tar "$rev" | tar -x -C "$export_dir"
fi

# Completeness: EVERY path that revision tracks must exist in the export, or the
# gate is scanning a tree with holes in it - and ban #9 in particular convicts a
# real citation as phantom when its target is missing, which would make the
# baseline red look like the repository's fault.
git -C "$root" ls-tree -r --name-only "$rev" >"$export_dir.tracked"
tracked=$(wc -l <"$export_dir.tracked")
missing=0
while IFS= read -r rel; do
    if [ ! -e "$export_dir/$rel" ]; then
        missing=$((missing + 1))
        if [ "$missing" -le 5 ]; then
            echo "export-red-proof: MISSING from export: $rel" >&2
        fi
    fi
done <"$export_dir.tracked"
echo "export-red-proof: export completeness: $tracked tracked paths, $missing missing"
[ "$missing" -eq 0 ] ||
    claim_failed "the export is missing $missing of $tracked tracked paths; an incomplete tree is not a CI checkout and the proofs below would be about a tree nobody checks out"

manifest() { # $1 tree, $2 manifest file; echoes one hash for (path, content)
    (cd "$1" && find . -type f -print0 | LC_ALL=C sort -z | xargs -0 sha256sum) >"$2"
    sha256sum <"$2" | awk '{print $1}'
}

# The scan's own output is a section of the log, not the whole log: the wrapper
# runs `go test ./...` first and those tests PRINT self-report lines too, so a
# grep over the whole log would happily read a fixture's numbers as ours.
scan_section() { # $1 log, $2 out
    awk '/^d22scan\.sh: scan of /{on=1} on{print}' "$1" >"$2"
}

scope_n() { # $1 section, $2 exact scope label
    line=$(grep -F "d22scan: scope $2 " "$1" | head -1) || true
    [ -n "$line" ] || return 1
    printf '%s\n' "$line" | sed -e 's/.*examined *//' -e 's/ .*//'
}

run_ci_step() { # $1 log, $2 description; sets rc
    echo "export-red-proof: running the CI lint step ($2): sh scripts/d22scan.sh"
    rc=0
    (cd "$export_dir" && sh scripts/d22scan.sh) >"$1" 2>&1 || rc=$?
    echo "export-red-proof: $2 exit code = $rc (log: $1)"
    scan_section "$1" "$1.scan"
}

# ---- baseline: the exported commit must be GREEN -----------------------------
before_manifest=$export_dir.manifest.before
before_hash=$(manifest "$export_dir" "$before_manifest")
echo "export-red-proof: baseline tree sha256 = $before_hash"
run_ci_step "$base_log" baseline
[ "$rc" -eq 0 ] || claim_failed "the exported $rev_short is not green, so a later red would not be the plant's: see $base_log (and its .scan section)"
grep -q "d22scan: clean" "$base_log.scan" || claim_failed "the baseline scan printed no clean verdict: see $base_log.scan"
base_scopes=$(grep -c -F "d22scan: scope " "$base_log.scan" || true)
[ "$base_scopes" -gt 0 ] || claim_failed "the baseline scan printed no per-scope self-report lines at all"
echo "export-red-proof: baseline self-report lines = $base_scopes"
base_internal=$(scope_n "$base_log.scan" "bans #1-5 internal/") || claim_failed "baseline has no 'bans #1-5 internal/' line"
echo "export-red-proof: baseline bans #1-5 internal/ examined $base_internal"

# ---- plant one violation -----------------------------------------------------
plant_dir=$export_dir/internal/zz71export
plant_file=$plant_dir/leak.go
plant_rel=internal/zz71export/leak.go
mkdir -p "$plant_dir"
cat >"$plant_file" <<'PLANT'
package zz71export

func worker() {}

func leak() { go worker() }
PLANT
echo "export-red-proof: planted $plant_rel (ban #1: bare goroutine)"

run_ci_step "$red_log" planted
[ "$rc" -eq 1 ] || claim_failed "a planted ban #1 violation must exit 1, got $rc: see $red_log"
red_internal=$(scope_n "$red_log.scan" "bans #1-5 internal/") || claim_failed "the planted run printed no 'bans #1-5 internal/' line, so the self-report vanished instead of moving"
echo "export-red-proof: planted bans #1-5 internal/ examined $red_internal (baseline $base_internal)"
grep -F "$plant_rel" "$red_log.scan" >/dev/null ||
    claim_failed "the red run does not NAME the planted file $plant_rel - a red nobody can locate is not the red AC#5 asks for"
grep -F "[bare-goroutine]" "$red_log.scan" >/dev/null ||
    claim_failed "the red run does not name the ban ([bare-goroutine])"
[ "$red_internal" -eq $((base_internal + 1)) ] ||
    claim_failed "the self-reported number for bans #1-5 internal/ moved from $base_internal to $red_internal; want exactly $((base_internal + 1)) - a number that does not move with the tree is a constant (AC#4)"
red_line=$(grep -F "$plant_rel" "$red_log.scan" | head -1)
echo "export-red-proof: RED line -> $red_line"

# ---- roll back: same bytes, green again --------------------------------------
rm -f "$plant_file"
rmdir "$plant_dir"
restore_manifest=$export_dir.manifest.restore
restore_hash=$(manifest "$export_dir" "$restore_manifest")
echo "export-red-proof: restored tree sha256 = $restore_hash"
if [ "$restore_hash" != "$before_hash" ]; then
    diff "$before_manifest" "$restore_manifest" | head -20 >&2 || true
    claim_failed "the tree hash after rollback ($restore_hash) differs from before the plant ($before_hash): the rollback is not a rollback"
fi
run_ci_step "$restore_log" restored
[ "$rc" -eq 0 ] || claim_failed "the restored tree is not green again (exit $rc): see $restore_log"
grep -q "d22scan: clean" "$restore_log.scan" || claim_failed "the restored run printed no clean verdict: see $restore_log.scan"
restore_scopes=$(grep -c -F "d22scan: scope " "$restore_log.scan" || true)
restore_internal=$(scope_n "$restore_log.scan" "bans #1-5 internal/") || claim_failed "the restored run printed no 'bans #1-5 internal/' line"
[ "$restore_scopes" -eq "$base_scopes" ] || claim_failed "the restored run printed $restore_scopes scope lines, baseline printed $base_scopes"
[ "$restore_internal" -eq "$base_internal" ] || claim_failed "the restored count is $restore_internal, want the baseline $base_internal"

echo "export-red-proof: RESULT PASS (rev $rev_short)"
echo "export-red-proof:   export tree   = $export_dir ($tracked tracked paths, $missing missing)"
echo "export-red-proof:   self-report   = $base_scopes scope lines per run; baseline and restored both green"
echo "export-red-proof:   planted       = $plant_rel -> exit 1, line names it, bans #1-5 internal/ $base_internal -> $red_internal"
echo "export-red-proof:   rolled back   = tree sha256 identical ($before_hash), exit 0, count back to $restore_internal"
echo "export-red-proof:   logs          = $base_log $red_log $restore_log"

if [ -z "$keep" ]; then
    rm -rf "$export_dir" "$base_log" "$red_log" "$restore_log" \
        "$base_log.scan" "$red_log.scan" "$restore_log.scan" \
        "$before_manifest" "$restore_manifest" "$export_dir.tracked" "$zip_file"
    echo "export-red-proof: export dir removed (use -k to keep it)"
fi
