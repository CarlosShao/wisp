#!/usr/bin/env bash
# Ticket 152 end-to-end probe, v2. v1 (e2e-subject-kill.sh) found NO VICTIM:
# its PowerShell one-liner was built inside bash double quotes, so the
# `-Filter "Name='wisp.exe'"` quoting arrived mangled and enumeration came back
# empty. That is a broken instrument, not a reading, so v2 keeps the raw v1 log,
# fixes the call (single quotes, no bash expansion, retries) and prints every
# wisp.exe it sees with its full command line, so the victim's identity is
# auditable rather than assumed.
#
# What v2 asks: when the REAL `wisp slo` observer is left waiting on a REAL
# subject child that a REAL TerminateProcess killed, what does the operator
# actually get on stderr, and after how long?
set -u
repo="/d/work/workspace/projects plans/Wisp"
scratch=/d/work/tmp/wisp152/scratch
mkdir -p "$scratch/e2e2"
exe="$scratch/wisp.exe"
log="$scratch/e2e2/e2e.log"
err="$scratch/e2e2/observer-stderr.txt"
: >"$log"
say() { echo "$@" | tee -a "$log"; }

cd "$repo" || exit 1
export PATH="/d/work/workspace/projects plans/Wisp/third_party/sherpa-onnx:$PATH"
say "# anchor=$(git rev-parse --short HEAD) date=$(date '+%Y-%m-%d %H:%M %z')"
say "# exe=$exe (built by v1 step 1, rc=0)"

ps_list='Get-CimInstance Win32_Process | Where-Object { $_.Name -eq "wisp.exe" } | ForEach-Object { "{0} {1}" -f $_.ProcessId, $_.CommandLine }'

say ""
say "## step 1: start the real observer (window 6s, then it enters collectReport)"
WISP_ENV=test WISP_TEST_DATA_DIR="$scratch/e2e2/data" \
	"$exe" slo -state Sleeping -seconds 6 -settle-ms 0 -interval-ms 250 \
	-out "$scratch/e2e2/kill-report.json" >>"$log" 2>"$err" &
observer=$!
say "observer shell pid=$observer at $(date '+%H:%M:%S')"

say ""
say "## step 2: enumerate wisp.exe children until the -self-sample one is seen"
victim=""
victim_line=""
for attempt in 1 2 3 4 5 6 7 8 9 10 11 12; do
	sleep 1
	kids=$(powershell.exe -NoProfile -Command "$ps_list" 2>>"$log" | tr -d '\r')
	if [ -n "$kids" ]; then
		say "--- attempt $attempt sees: ---"
		say "$kids"
		line=$(echo "$kids" | grep -- '-self-sample' | head -1)
		if [ -n "$line" ]; then
			victim=$(echo "$line" | awk '{print $1}')
			victim_line=$line
			say "victim pid=$victim"
			break
		fi
	fi
done
if [ -z "$victim" ]; then
	say "STILL NO VICTIM after 12 attempts - reporting the instrument as broken, not a reading"
	kill -9 "$observer" 2>/dev/null
	wait "$observer" 2>/dev/null
	exit 3
fi

say ""
say "## step 3: kill it for real (TerminateProcess via Stop-Process)"
say "victim command line: $victim_line"
t0=$(date +%s)
powershell.exe -NoProfile -Command "Stop-Process -Id $victim -Force" >>"$log" 2>&1
say "Stop-Process rc=$?"
powershell.exe -NoProfile -Command \
	"if (Get-Process -Id $victim -ErrorAction SilentlyContinue) { 'victim STILL RUNNING' } else { 'victim GONE (no process object)' }" |
	tee -a "$log"

say ""
say "## step 4: wait for the observer to give up and print its verdict"
wait "$observer"
rc=$?
say "observer rc=$rc elapsed_since_kill=$(( $(date +%s) - t0 ))s"
say "--- observer stderr verbatim (what a human running \`wisp slo\` sees) ---"
say "stderr bytes=$(wc -c <"$err")"
say "$(cat "$err")"
say "--- out file ---"
say "kill report bytes=$(wc -c <"$scratch/e2e2/kill-report.json" 2>/dev/null || echo 'NO FILE (the run failed closed)')"
say "# end $(date '+%Y-%m-%d %H:%M %z')"
