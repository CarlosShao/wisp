#!/bin/sh
# 342-r2 evidence commit. add and commit BOTH carry the literal per-file pathspecs (no directory form,
# no -A, no reverse-derivation from the shared index - the public block's rule about index poisoning).
cd "/d/work/workspace/projects plans/Wisp" || exit 1
FILES=".scratch/wisp/probes/342/r2/report.md
.scratch/wisp/probes/342/r2/logs/bench-probes-verbatim.txt
.scratch/wisp/probes/342/r2/logs/gates-and-roots.txt
.scratch/wisp/probes/342/r2/logs/appended-lines.txt
.scratch/wisp/probes/342/r2/logs/base-det-r1.log
.scratch/wisp/probes/342/r2/logs/base-det-r2.log
.scratch/wisp/probes/342/r2/logs/base-det-r3.log
.scratch/wisp/probes/342/r2/logs/base-det-shortroot.log
.scratch/wisp/probes/342/r2/logs/base-full-2.log
.scratch/wisp/probes/342/r2/logs/base-full-3.log
.scratch/wisp/probes/342/r2/logs/base-targeted-1.log
.scratch/wisp/probes/342/r2/logs/d22scan-final.log
.scratch/wisp/probes/342/r2/logs/land-det-r1.log
.scratch/wisp/probes/342/r2/logs/land-det-r2.log
.scratch/wisp/probes/342/r2/logs/land-det-r3.log
.scratch/wisp/probes/342/r2/logs/land-det-shortroot.log
.scratch/wisp/probes/342/r2/logs/land-full-2.log
.scratch/wisp/probes/342/r2/logs/land-full-3.log
.scratch/wisp/probes/342/r2/logs/land-targeted-1.log
.scratch/wisp/probes/342/r2/logs/landminus-full-1.log
.scratch/wisp/probes/342/r2/logs/leaf-base-r1.log
.scratch/wisp/probes/342/r2/logs/leaf-base-r2.log
.scratch/wisp/probes/342/r2/logs/leaf-base-r3.log
.scratch/wisp/probes/342/r2/logs/leaf-land-v2-r1.log
.scratch/wisp/probes/342/r2/logs/leaf-land-v2-r2.log
.scratch/wisp/probes/342/r2/logs/leaf-land-v2-r3.log
.scratch/wisp/probes/342/r2/logs/lenbudget-final.log
.scratch/wisp/probes/342/r2/logs/mutA-targeted-r1.log
.scratch/wisp/probes/342/r2/logs/mutA-targeted-r2.log
.scratch/wisp/probes/342/r2/logs/mutB-targeted-r1.log
.scratch/wisp/probes/342/r2/logs/sid-sharedtree-nogotmp.log
.scratch/wisp/probes/342/r2/sh/codepoint-ruler.pl
.scratch/wisp/probes/342/r2/sh/commit-evidence.sh
.scratch/wisp/probes/342/r2/sh/gates-and-roots.sh
.scratch/wisp/probes/342/r2/sh/roots-and-faults.sh
.scratch/wisp/probes/342/r2/sh/self-proof-rulers.sh"

# shellcheck disable=SC2086
git add -- $FILES
echo "ADD_RC=$?"
git diff --cached --name-only | wc -l | sed 's/^/STAGED_COUNT=/'
git diff --cached --name-only | grep -v '^\.scratch/wisp/probes/342/r2/' | sed 's/^/OUTSIDE-FACE: /'
git commit -F /d/tmp/wisp/20261019/342r2/commit7.txt -- $FILES
echo "COMMIT_RC=$?"
git log -1 --format='%H %s'
git show --name-only --format= | head -40
git show --name-only --format= | wc -l | sed 's/^/NAMES_IN_COMMIT=/'
