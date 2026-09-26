#!/usr/bin/env bash
# Ticket 152 end-to-end probe: the OPERATOR-facing red sentence, produced by the
# real `wisp slo` observer with a real `wisp` subject child killed for real.
# Everything here runs outside the repository; the raw log is copied into
# .scratch/wisp/probes/152/ afterwards.
set -u
repo="/d/work/workspace/projects plans/Wisp"
scratch=/d/work/tmp/wisp152/scratch
mkdir -p "$scratch/e2e"
exe="$scratch/wisp.exe"
log="$scratch/e2e/e2e.log"
: >"$log"

say() { echo "$@" | tee -a "$log"; }

say "# anchor=$(cd "$repo" && git rev-parse --short HEAD) date=$(date '+%Y-%m-%d %H:%M %z')"
say "# cwd=$repo  exe=$exe"

cd "$repo" || exit 1
export PATH="/d/work/workspace/projects plans/Wisp/third_party/sherpa-onnx:$PATH"

say "## step 1: build the real CLI outside the repo"
say "\$ go build -o $exe ./cmd/wisp"
go build -o "$exe" ./cmd/wisp >>"$log" 2>&1
say "build rc=$?"
ls -la "$exe" | tee -a "$log"

say ""
say "## step 2: BASELINE - a healthy real run, same flags, no kill (control)"
say "\$ $exe slo -state Sleeping -seconds 3 -settle-ms 0 -interval-ms 250 -out \$out"
t0=$(date +%s)
WISP_ENV=test WISP_TEST_DATA_DIR="$scratch/e2e/data-base" \
	"$exe" slo -state Sleeping -seconds 3 -settle-ms 0 -interval-ms 250 \
	-out "$scratch/e2e/base-report.json" >>"$log" 2>"$scratch/e2e/base-stderr.txt"
rc=$?
say "baseline rc=$rc elapsed=$(( $(date +%s) - t0 ))s"
say "baseline stderr: $(wc -c <"$scratch/e2e/base-stderr.txt") bytes"
say "baseline report bytes=$(wc -c <"$scratch/e2e/base-report.json" 2>/dev/null || echo 0)"

say ""
say "## step 3: KILL the instrumented (-self-sample) subject child mid-run"
say "# the observer enters collectReport at ~t+3s; the subject would write its"
say "# report at ~t+6s (sampling 3s + subjectGrace 3s). Kill at ~t+4.5s, so the"
say "# observer waits on a corpse."
WISP_ENV=test WISP_TEST_DATA_DIR="$scratch/e2e/data-kill" \
	"$exe" slo -state Sleeping -seconds 3 -settle-ms 0 -interval-ms 250 \
	-out "$scratch/e2e/kill-report.json" >>"$log" 2>"$scratch/e2e/observer-stderr.txt" &
observer=$!
t0=$(date +%s)
sleep 4.5
say "\$ Get-CimInstance Win32_Process -Filter \"Name='wisp.exe'\" | where CommandLine -match '-self-sample'"
kids=$(powershell.exe -NoProfile -Command \
	"Get-CimInstance Win32_Process -Filter \"Name='wisp.exe'\" | ForEach-Object { \"\$(\$_.ProcessId)\|\$(\$_.CommandLine)\" }" \
	2>>"$log" | tr -d '\r')
say "--- every wisp.exe process the OS reports ---"
say "$kids"
victim=$(echo "$kids" | grep -- '-self-sample' | head -1 | cut -d'|' -f1)
say "victim pid=$victim"
if [ -z "$victim" ]; then
	say "NO VICTIM FOUND - the self-sample child was not visible; this run measures nothing"
	kill -9 "$observer" 2>/dev/null
else
	say "\$ Stop-Process -Id $victim -Force   (a real TerminateProcess on a real subject)"
	powershell.exe -NoProfile -Command "Stop-Process -Id $victim -Force" >>"$log" 2>&1
	say "stop rc=$?"
	say "\$ GetExitCodeProcess-equivalent check that pid $victim is gone:"
	powershell.exe -NoProfile -Command \
		"if (Get-Process -Id $victim -ErrorAction SilentlyContinue) { 'STILL RUNNING' } else { 'GONE (no process object)' }" >>"$log" 2>&1
fi
wait "$observer"
rc=$?
say "observer rc=$rc elapsed=$(( $(date +%s) - t0 ))s"
say "--- the kill run's out file ---"
say "kill report bytes=$(wc -c <"$scratch/e2e/kill-report.json" 2>/dev/null || echo 'no file')"
say "--- observer stderr verbatim (what a human running \`wisp slo\` actually sees) ---"
say "stderr bytes=$(wc -c <"$scratch/e2e/observer-stderr.txt")"
say "$(cat "$scratch/e2e/observer-stderr.txt")"
say "--- baseline stderr for comparison ---"
say "$(cat "$scratch/e2e/base-stderr.txt")"
say "# end $(date '+%Y-%m-%d %H:%M %z')"
