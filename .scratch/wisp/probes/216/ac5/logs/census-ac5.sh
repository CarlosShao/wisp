#!/bin/sh
# Ticket 216 AC#5 - the census rulers, all re-runnable.
#
#   sh .scratch/wisp/probes/216/ac5/logs/census-ac5.sh \
#     > .scratch/wisp/probes/216/ac5/logs/census-ac5-before.log 2>&1
#
# Optionally point it at a plain export of an older tree (git archive HEAD,
# no .git): the rulers then run through `git grep --no-index`:
#
#   sh .scratch/wisp/probes/216/ac5/logs/census-ac5.sh /d/wisp-216ac5-base2 \
#     > .scratch/wisp/probes/216/ac5/logs/census-ac5-before.log 2>&1
#
# AC#5 (verbatim): "错误提示按错误码分支，不许按英文错误串里有没有某个数字".
# The negative example is openchamber ChatInput.tsx:2205-2206, transcribed
# verbatim in ../control/openchamber-chatinput-2205.txt (the upstream repo is
# not in this tree).
#
# Every ruler prints id / pat / rng / rc / hits, then an explicit "expect=".
# The zero rulers carry controls: a 0 that comes from a blind pattern is the
# failure this project keeps catching (216-v2's D13 false zero, ledger A1010).
set -u

if [ $# -ge 1 ]; then
	cd "$1" || exit 1
	GG="git grep --no-index"
	ANCHOR="export-tree=$1 (no git metadata; the export was git archive HEAD at the anchor named in the evidence file)"
else
	cd "$(git rev-parse --show-toplevel)" || exit 1
	GG="git grep"
	ANCHOR="HEAD=$(git rev-parse --short HEAD)"
fi
echo "anchor: $ANCHOR date=$(date -u +%Y-%m-%dT%H:%M:%SZ)"
echo "Every ruler prints id / pat / rng / rc / hits, then an explicit expect=;"
echo "the zero rulers carry controls so a 0 is never read as a blind zero."
echo

ruler() { # id, pattern, expectation, pathspec...
	id=$1; pat=$2; expect=$3; shift 3
	echo "=== $id"
	echo "pat=$pat"
	echo "rng=$*"
	out=$($GG -n -E "$pat" -- "$@" 2>&1); rc=$?
	echo "rc=$rc"
	[ -n "$out" ] && echo "$out"
	echo "expect=$expect"
	echo
}

count() { # id, pattern, pathspec...
	id=$1; pat=$2; shift 2
	echo "=== $id"
	echo "pat=$pat"
	n=$($GG -c -E "$pat" -- "$@" 2>/dev/null | awk -F: '{s+=$NF} END{print s+0}')
	echo "count=$n"
	echo "rng=$*"
	echo
}

echo "##### A. the forbidden shape: a user-visible sentence selected by error TEXT"
echo

ruler A1-error-text-predicate-lines \
	'(Contains|HasPrefix|HasSuffix|Index|MatchString|EqualFold)\([^)]*\.Error\(\)' \
	'the predicate forms on production .go: one code hit before the fix (resident_voice_engines:1240); the after log is the pair' \
	'internal/**/*.go' \
	'cmd/**/*.go' \
	'tools/**/*.go' \
	':(exclude)*_test.go'

ruler A2-literal-english-equality \
	'\.Error\(\) *[!=]= *"' \
	'exact English text equality (not copy selection; judged in the evidence file)' \
	'internal/**/*.go' \
	'cmd/**/*.go' \
	'tools/**/*.go' \
	':(exclude)*_test.go'

ruler A3-switch-on-error-text \
	'switch [a-zA-Z_.]*\.Error\(\)' \
	'0: nobody switches on the whole rendered error text either' \
	'internal/**/*.go' \
	'cmd/**/*.go' \
	'tools/**/*.go' \
	':(exclude)*_test.go'

echo "##### B. the forbidden shape one layer down: a copy kind selected from a string FIELD"
echo

ruler B1-detail-alias \
	':= *[a-zA-Z_.]+\.(Detail|Message|ProviderCode)\b' \
	'1 hit: cmd/wisp/config_reload.go:431 d := oe.Detail (the only consumer that branches on it is B2)' \
	'internal/**/*.go' \
	'cmd/**/*.go' \
	'tools/**/*.go' \
	':(exclude)*_test.go'

ruler B2-detail-branch \
	'(HasPrefix|Contains)\(d, *"' \
	'5 branches in cmd/wisp/config_reload.go describeReloadFailure selecting 6 operator sentences' \
	'internal/**/*.go' \
	'cmd/**/*.go' \
	':(exclude)*_test.go'

ruler B3-upstream-message-branch \
	'Contains\(strings\.ToLower\(d\.Message\)' \
	'internal/llm/errors.go: upstream PROSE (d.Message) used to pick an error class -> picks the copy' \
	'internal/**/*.go' \
	':(exclude)*_test.go'

ruler B4-upstream-code-branch \
	'Contains\(code, *"' \
	'internal/llm/errors.go: upstream CODE (machine token) matched by substring - judged a code branch, see evidence' \
	'internal/**/*.go' \
	':(exclude)*_test.go'

echo "##### C. the exact negative shape: a numeric literal matched as a STRING"
echo

ruler C1-status-literal-as-string \
	'"(413|429|400|401|403|404|500|502|503)"' \
	'0 production hits; the hits below are the two comments that NAME the forbidden shape' \
	'internal/**/*.go' \
	'cmd/**/*.go' \
	':(exclude)*_test.go'

ruler C2-status-literal-as-int \
	'case (413|429|400|401|403|404|500|502|503)[,:]' \
	'positive control: the same numbers as INT case labels DO exist (that is the code branch AC#5 wants)' \
	'internal/**/*.go' \
	'cmd/**/*.go' \
	':(exclude)*_test.go'

echo "##### D. frontend (READ-ONLY: frontend/** is the owner external agent)"
echo

ruler D1-frontend-negative-shape \
	"includes\(['\"]413['\"]\)|payload too large|entity too large" \
	'0 hits: the page does not sniff the error prose either' \
	'frontend/src'

count D2-frontend-includes-control \
	'[.]includes\(' \
	'frontend/src'

ruler D3-frontend-status-branch \
	"statusText|status *[!=]== *['\"]" \
	'hits compare OUR enum strings (task/tool status), not error prose; judged per line in evidence' \
	'frontend/src'

echo "##### E. the PASS side: copy selected by code (the branch AC#5 requires)"
echo

count E1-panel-copy-table \
	'observe\.Class[A-Za-z]+:' \
	'internal/panel/error_copy_216.go'

ruler E2-panel-copy-head \
	'func errorHead216' \
	'the one branch point: class in, head out' \
	'internal/panel/*.go'

ruler E3-classof-consumers \
	'observe\.ClassOf\(' \
	'production consumers of the class extractor (the code source the copy branch rides on)' \
	'internal/**/*.go' \
	'cmd/**/*.go' \
	':(exclude)*_test.go'

echo "##### F. the two-sentence discipline: 没有 vs 读不到 must not share one sentence"
echo

ruler F1-instructions-five-states \
	'InstructionsStatus[A-Za-z]+ *= *"' \
	'five named states; none_found and not_run are separate constants' \
	'internal/panel/instructions_200.go'

ruler F2-workspace-two-sentences \
	'workspace(NotChosen|NoReading)Reason *=' \
	'two constants: "not chosen" vs "no reading" are two sentences' \
	'internal/panel/workspace.go'

echo "##### G. controls: prove the rulers above are not blind"
echo

echo "=== G1-control-verbatim-negative-example"
echo "pat=(the C1 pattern run against the stored verbatim transcription)"
if [ -f .scratch/wisp/probes/216/ac5/control/openchamber-chatinput-2205.txt ]; then
	grep -n -E '"(413|429|400|401|403|404|500|502|503)"' \
		.scratch/wisp/probes/216/ac5/control/openchamber-chatinput-2205.txt
	echo "rc=$?"
else
	echo "skipped: this tree has no .scratch/ (plain git-archive export); the control was run in the repo and is in the repo log"
fi
echo "expect=>=1 hit: the C1 pattern DOES fire on the shape it is meant to catch"
echo

echo "=== G2-control-frontend-includes"
echo "pat=[.]includes\\("
n=$(grep -rn -E '[.]includes\(' frontend/src | wc -l)
echo "hits=$n"
echo "expect=>0: the frontend ruler family is not blind"
echo

echo "=== G3-control-census-dir"
if [ -d .scratch/wisp/probes/216/ac5 ]; then
	find .scratch/wisp/probes/216/ac5 -type f | sed 's/^/file=/' | sort
else
	echo "skipped: .scratch/ is not part of the plain export"
fi
echo

echo "##### H. production denominators (not a judgment: sizes, so a zero can be read against them)"
count H1-error-call-sites '\.Error\(\)' 'internal/**/*.go' 'cmd/**/*.go' 'tools/**/*.go' ':(exclude)*_test.go'
count H2-errors-is-as '(errors\.Is|errors\.As)\(' 'internal/**/*.go' 'cmd/**/*.go' 'tools/**/*.go' ':(exclude)*_test.go'
