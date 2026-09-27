#!/bin/sh
# 161-v2 AC#7 delete-experiments, recorded verbatim so a later run can re-do them
# without guessing my sed expressions. NOTHING here writes inside the repository:
# every mutated copy lives in $TMPDIR (or, for the classifier mutation, in this
# directory as an explicitly-labelled derived file).
#
# Two experiments, both run at anchor 30121fef17a42791eadab463c4773325a98838e2:
#
# (1) delete the "the ticket file really exists" flavour of the attribution rule
#     orig : sh .scratch/wisp/probes/161/r5/attrib.sh --classify-line '.scratch/wisp/probes/999/bad-sample.go'
#            -> UNATTRIBUTABLE ... <== REAL INJURY          rc=1
#     mut  : sh attrib-noknown.sh --classify-line 同上
#            -> ticket-999 (untracked, parse_err=0)          rc=0     <- self-proving green
#     and   : sh attrib-noknown.sh --self-test -> cases=8 failures=4  rc=1  <- caught by its own pair
#
# (2) delete the hollow guard of ruler (A) in an empty synthetic git tree
#     orig : attrib.sh --tracked-only over 0 tracked .go -> "refusing to report 'empty' ..." rc=2
#     mut  : same tree, guard line removed             -> "lines=0 files=0" + GREEN        rc=0
#
# Commands used (kept as comments so nobody has to re-derive them):
#
#   cp .scratch/wisp/probes/161/r5/attrib.sh .scratch/wisp/probes/161/v2/attrib-noknown.sh
#   perl -0pi -e "s/(ticket_known\\(\\) \\{)/\$1\\n    printf '1'; return 0 # existence flavour removed/" attrib-noknown.sh
#
#   E=$(mktemp -d)/emptyrepo; mkdir -p "$E/probe"; cd "$E"; git init -q .
#   printf 'module empty.test/v2\n\ngo 1.24\n' > go.mod
#   cp .../probes/161/r5/attrib.sh "$E/probe/attrib.sh"
#   cp .../probes/161/r5/attrib.sh "$E/probe/attrib-noguard.sh"
#   perl -0pi -e 's/\[ "\$TRACKED_GO" -gt 0 \] \|\| \{ echo "attrib\.sh: \(A\)[^"]*"[^\n]*\n/# hollow guard removed\n/' "$E/probe/attrib-noguard.sh"
#   sh "$E/probe/attrib.sh"        --tracked-only   # rc=2
#   sh "$E/probe/attrib-noguard.sh" --tracked-only  # rc=0  <-- vacuous green
