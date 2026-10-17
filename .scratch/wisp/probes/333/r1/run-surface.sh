#!/bin/sh
# Ticket 333 leg 333-r1: one colour reading of internal/panel on ONE named surface.
# The command line is IDENTICAL for every surface; only TMP/TEMP/GOTMPDIR move.
#   usage: sh .scratch/wisp/probes/333/r1/run-surface.sh <label> <c|d>
label="$1"
surface="$2"
cd "/d/work/workspace/projects plans/Wisp" || exit 99
log=".scratch/wisp/probes/333/r1/logs/$label.log"

if [ "$surface" = "d" ]; then
	# D-volume temp root, both dirs created up front with mkdir -p.
	export TMP="D:/wisp-333r1-tmp/tmp"
	export TEMP="D:/wisp-333r1-tmp/tmp"
	export GOTMPDIR="D:/wisp-333r1-tmp/gotmp"
else
	# C-volume temp root = the shell default (%LOCALAPPDATA%\Temp).
	export TMP="/tmp"
	export TEMP="/tmp"
	unset GOTMPDIR
fi

go test -v ./internal/panel/ -count=1 -timeout 30m >"$log" 2>&1
rc=$?
	printf 'label=%s surface=%s rc=%s\n' "$label" "$surface" "$rc"
	printf 'command=%s\n' 'go test -v ./internal/panel/ -count=1 -timeout 30m'
sha256sum "$log" >"$log.sha256"
exit $rc
