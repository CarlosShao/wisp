#!/bin/sh
# Ticket 192 AC#2 -- "the criterion can go red, and it names the line".
#
# The main repository is NEVER mutated here. Every mutation is applied to
# $TREE, a git-archive export of HEAD (out of the repo, per the dispatch),
# with this ticket's two files copied in byte-for-byte from the working tree
# so the exported copy judges the same code the commit ships.
#
# WHAT IT SHOWS, in order:
#   01 baseline           the export tree's copy is green (same bytes as the
#                         repo), so what follows is about the mutation.
#   M1 quote-not-stamped  the C25 Mark on the quote door is short-circuited
#                         => the "the stamp is live" assertion must red.
#   M2 laundering allowed the "the quote must be inside the artifact" check is
#                         disabled => the laundering assertion must red.
#   M3 sandbox dropped    the HTML preview's CSP loses its sandbox directive
#                         => the "a preview cannot speak" assertion must red.
#   M4 pointer dropped    the truncation pointer is zeroed => the "notice
#                         shows head N but the pointer resumes at 0"
#                         assertion must red.
#   M5 name guard dropped the artifacts bare-name guard is skipped => the
#                         "refused before it opens anything" assertion must
#                         red.
#   every restore         byte-for-byte back: sha256 printed and compared to
#                         the pristine one, then the whole set is green again.
#
# Usage: sh .scratch/wisp/probes/192/redproof-192.sh [tree]
set -u

tree=${1:-/d/tmp/wisp/20261114/192-tree}
here=$(cd "$(dirname "$0")" && pwd)
repo=$(cd "$here/../../../.." && pwd)
logs="$here/logs"
[ -d "$tree" ] || { echo "redproof-192: no export tree at $tree" >&2; exit 2; }
[ -f "$tree/go.mod" ] || { echo "redproof-192: $tree has no go.mod" >&2; exit 2; }
mkdir -p "$logs"

pkg="$tree/internal/panel"
file="$pkg/preview_192.go"
cp "$repo/internal/panel/preview_192.go" "$file" || exit 2
cp "$repo/internal/panel/preview_192_test.go" "$pkg/preview_192_test.go" || exit 2

sha() { sha256sum "$1" | cut -d' ' -f1; }

run() { # run <label> <test-regex>
	label=$1
	regex=$2
	log="$logs/$label.log"
	(cd "$tree" && go test -count=1 -run "$regex" ./internal/panel/) >"$log" 2>&1
	rc=$?
	echo "== $label rc=$rc log=$log"
	grep -E '^(--- FAIL|--- PASS|ok |FAIL|PASS)|preview_192_test\.go:[0-9]+:' "$log" | head -8 | sed 's/^/   /'
	return $rc
}

pristine=$(sha "$file")
echo "pristine sha256 $pristine  ($file)"
run 01-baseline 'TestPreview192' || { echo "redproof-192: the export tree is not green before any mutation; nothing below would mean anything" >&2; exit 1; }

bad=0
mutate() { # mutate <name> <sed-expr> <label> <test-regex>
	name=$1
	expr=$2
	label=$3
	regex=$4
	cp "$file" "$logs/mut-$name.before"
	sed -i "$expr" "$file"
	echo "--- mutant $name: sed -i '$expr' (sha256 $(sha "$file"))"
	diff "$logs/mut-$name.before" "$file" | sed 's/^/   /' | head -6
	if run "$label" "$regex"; then
		echo "   MUTANT $name DID NOT GO RED"
		bad=1
	fi
	cp "$logs/mut-$name.before" "$file"
	restored=$(sha "$file")
	if [ "$restored" != "$pristine" ]; then
		echo "   RESTORE OF $name FAILED: sha256 $restored != $pristine" >&2
		exit 2
	fi
	echo "   restored sha256 $restored (equals pristine)"
}

mutate m1-quote-not-stamped \
	's/if !marker\.Mark(scopeID, PreviewSourceTool/if false \&\& marker.Mark(scopeID, PreviewSourceTool/' \
	02-m1-quote-not-stamped 'TestPreview192QuoteDoorStampsAndTheEgressGateThenCatches'

mutate m2-laundering-allowed \
	's/if !bytes\.Contains(body, \[\]byte(quoted))/if !bytes.Contains(body, []byte(""))/' \
	03-m2-laundering-allowed 'TestPreview192QuoteDoorRefusesTextThatIsNotInTheArtifact'

mutate m3-sandbox-dropped \
	"s/sandbox; default-src 'none'; img-src data:/default-src 'none'; img-src data:/" \
	04-m3-sandbox-dropped 'TestPreview192HTMLIsSandboxedSoAPreviewCannotSpeak'

mutate m4-pointer-dropped \
	's/doc\.NextOffset = int64(len(head))/doc.NextOffset = 0/' \
	05-m4-pointer-dropped 'TestPreview192TextFirstViewSaysWhatItHidAndWhereTheRestIs'

mutate m5-name-guard-dropped \
	's/if err := r\.guard(name); err != nil {/if err := error(nil); err != nil {/' \
	06-m5-name-guard-dropped 'TestPreview192RefusesPathShapedNamesBeforeItOpensAnything'

echo "--- final green run on the restored bytes"
run 07-final-green 'TestPreview192' || exit 1
final=$(sha "$file")
echo "final sha256 $final (pristine $pristine)"
[ "$final" = "$pristine" ] || exit 2
[ "$bad" = 0 ] || { echo "redproof-192: at least one mutant stayed green" >&2; exit 1; }
echo "redproof-192: 5 mutants, each went red on its own named assertion, each restored to the pristine bytes"
