#!/usr/bin/env bash
# Ticket 152 end-to-end probe, v3. Identical to v2 except WHOSE BYTES the
# observer runs from: v2 used /d/work/tmp/wisp152/scratch/wisp.exe, which was
# built at 09:26 from the live working tree - and that tree had already picked up
# ticket 151's uncommitted cmd/wisp/run.go edit (mtime 09:25:30). v3 builds from
# the clean snapshot snap-pre (git archive 5365cb2, cmd/wisp blobs md5-checked),
# so the reading is attributable to the anchor and to nobody else's in-flight work.
set -u
snap=/d/work/tmp/wisp152/snap-pre
scratch=/d/work/tmp/wisp152/scratch
mkdir -p "$scratch/e2e3"
exe="$scratch/wisp-anchor.exe"
log="$scratch/e2e3/e2e.log"
err="$scratch/e2e3/observer-stderr.txt"
: >"$log"
say() { echo "$@" | tee -a "$log"; }

export PATH="$snap/third_party/sherpa-onnx:$PATH"

say "# date=$(date '+%Y-%m-%d %H:%M %z')"
say "## step 0: prove the bytes the exe is built from"
say "\$ git -C $snap rev-parse HEAD   (a git-archive snapshot has no .git)"
for f in cmd/wisp/slo_windows.go cmd/wisp/slo_report_144_windows_test.go cmd/wisp/run.go; do
	a=$(cd "/d/work/workspace/projects plans/Wisp" && git show 5365cb2:$f | md5sum | cut -d' ' -f1)
	b=$(md5sum < "$snap/$f" | cut -d' ' -f1)
	say "$f anchor=$a snap=$b same=$([ "$a" = "$b" ] && echo YES || echo NO)"
done

say ""
say "## step 1: build the anchor CLI"
say "\$ (cd $snap && go build -o $exe ./cmd/wisp)"
(cd "$snap" && go build -o "$exe" ./cmd/wisp) >>"$log" 2>&1
say "build rc=$? exe=$exe size=$(wc -c <"$exe" 2>/dev/null)"

ps_list='Get-CimInstance Win32_Process | Where-Object { $_.Name -eq "wisp-anchor.exe" -or $_.Name -eq "wisp.exe" } | ForEach-Object { "{0} {1}" -f $_.ProcessId, $_.CommandLine }'

say ""
say "## step 2: start the real observer (window 6s), then find its -self-sample child"
cd "$scratch" || exit 1
WISP_ENV=test WISP_TEST_DATA_DIR="$scratch/e2e3/data" \
	"$exe" slo -state Sleeping -seconds 6 -settle-ms 0 -interval-ms 250 \
	-out "$scratch/e2e3/kill-report.json" >>"$log" 2>"$err" &
observer=$!
say "observer shell pid=$observer at $(date '+%H:%M:%S')"

victim=""
victim_line=""
for attempt in 1 2 3 4 5 6 7 8; do
	sleep 1
	kids=$(powershell.exe -NoProfile -Command "$ps_list" 2>>"$log" | tr -d '\r')
	if [ -n "$kids" ]; then
		say "--- attempt $attempt sees: ---"
		say "$kids"
		line=$(echo "$kids" | grep -- 'wisp-anchor.exe' | grep -- '-self-sample' | head -1)
		if [ -n "$line" ]; then
			victim=$(echo "$line" | awk '{print $1}')
			victim_line=$line
			break
		fi
	fi
done
if [ -z "$victim" ]; then
	say "NO VICTIM after 8 attempts - instrument reported broken, no reading claimed"
	kill -9 "$observer" 2>/dev/null
	wait "$observer" 2>/dev/null
	exit 3
fi
case $victim in
*[!0-9]*) say "victim token '$victim' is not a bare pid - instrument broken, refusing to claim a reading"; exit 4 ;;
esac
say "victim pid=$victim"

say ""
say "## step 3: kill it for real"
say "victim command line: $victim_line"
t0=$(date +%s)
powershell.exe -NoProfile -Command "Stop-Process -Id $victim -Force" >>"$log" 2>&1
say "Stop-Process rc=$?"
powershell.exe -NoProfile -Command \
	"if (Get-Process -Id $victim -ErrorAction SilentlyContinue) { 'victim STILL RUNNING' } else { 'victim GONE (no process object)' }" |
	tee -a "$log"
say "victim subject report file at kill time: $(ls -la "$(echo "$victim_line" | sed -n 's/.*-out \(.*\) -self-sample.*/\1/p')" 2>&1 | head -1)"

say ""
say "## step 4: the observer's verdict"
wait "$observer"
rc=$?
say "observer rc=$rc elapsed_since_kill=$(( $(date +%s) - t0 ))s"
say "--- observer stderr verbatim ---"
say "stderr bytes=$(wc -c <"$err")"
say "$(cat "$err")"
say "kill report bytes=$(wc -c <"$scratch/e2e3/kill-report.json" 2>/dev/null || echo 'NO FILE (failed closed)')"
say "# end $(date '+%Y-%m-%d %H:%M %z')"
