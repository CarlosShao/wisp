#!/bin/sh
# 120-r3 gate driver: the two container rosters and the four host rulers, each with its
# own rc= on disk. Snapshot = `git archive HEAD` on the D: volume (the container cannot
# see Git Bash' /tmp, which is C: Temp and arrives empty); temp inside the runs is the
# image's own overlay /tmp, the same table tooth.sh used, so the paired readings are
# comparable. Host rulers run with GOTMPDIR on the C: volume, created first: pointing
# it at a directory that does not exist produces "creating work dir ..." AND rc=0.
set -e
REPO='D:/work/workspace/projects plans/Wisp'
S=/d/tmp/probe120r3
SNAP=$S/fixed
LOGS=$S/logs
export GOTMPDIR='C:/Users/swq/AppData/Local/wisp-test-gotmp/120r3gotmp'
mkdir -p "$GOTMPDIR" "$LOGS"
cd "$REPO"

echo "== snapshot: rebuilt from HEAD $(git rev-parse --short HEAD), .go counted before any measurement =="
rm -rf "$SNAP"
mkdir -p "$SNAP"
git archive HEAD | tar -x -C "$SNAP"
echo "go-files=$(find "$SNAP" -name '*.go' | wc -l) winsec-go=$(ls "$SNAP"/internal/winsec/*.go | wc -l)"

roster() {
	name=$1
	pat=$2
	body=$LOGS/$name-body.log
	rc=0
	# MSYS_NO_PATHCONV belongs to the docker call ONLY: exported into the host rulers
	# below it mangles the path scripts/d22scan.sh derives from its own location
	# ("D:\d\work\..." is not a repository root), which is a table failure, not a scan.
	MSYS_NO_PATHCONV=1 docker run --rm -e CGO_ENABLED=0 -v "$SNAP":/work -v probe120r3build:/root/.cache/go-build \
		-w /work golang:1.27 sh -c "go test -count=1 -v -run '$pat' ./internal/winsec/ > /work/roster.log 2>&1; rc=\$?; echo \"go-test rc=\$rc\"; exit \$rc" \
		> "$LOGS/$name.log" 2>&1 || rc=$?
	cp "$SNAP/roster.log" "$body"
	{
		echo "=== $name: run-filter='$pat' docker-and-test rc=$rc"
		printf 'RUN=%s (=== RUN rows, t.Run subtests included)\n' "$(grep -c '=== RUN' "$body" || true)"
		printf 'top-level result rows (anchored ^---): PASS=%s FAIL=%s SKIP=%s\n' \
			"$(grep -c -- '^--- PASS' "$body" || true)" \
			"$(grep -c -- '^--- FAIL' "$body" || true)" \
			"$(grep -c -- '^--- SKIP' "$body" || true)"
		printf 'all result rows (unanchored): PASS=%s FAIL=%s SKIP=%s\n' \
			"$(grep -c -- '--- PASS' "$body" || true)" \
			"$(grep -c -- '--- FAIL' "$body" || true)" \
			"$(grep -c -- '--- SKIP' "$body" || true)"
	} | tee -a "$LOGS/$name.log"
}

echo; echo "== container roster 1: the ticket's own AC#1-4 POSIX filter, counted as it stands today =="
roster roster-ac1234 'TestAC[1-4]POSIX'
echo "top-level names in the filter, from the source not the log: $(grep -h '^func TestAC[1-4]POSIX' internal/winsec/*_test.go | wc -l)"
echo "names this leg added: $(grep -h '^func TestAC1POSIX.*120r3' internal/winsec/*_test.go | wc -l)"

echo; echo "== container roster 2: this leg's seven cases alone (names ending in 120r3) =="
roster roster-120r3 '120r3$'

echo; echo "== host rulers =="
# From here on a non-zero exit is a READING to be recorded, not a reason to stop the
# run: every step captures its own rc= onto disk (dispatch 3.4: no `cmd | tail` then
# PIPESTATUS).
set +e
FILES='internal/winsec/sealhandle_window_120_other_test.go internal/winsec/sealhandle_identity_120_other_test.go internal/winsec/sealhandle_privateroute_120_other_test.go'
gofmt -l $FILES > "$LOGS/gofmt.log" 2>&1
echo "gofmt -l rc=$? body=$(wc -c < "$LOGS/gofmt.log") bytes" >> "$LOGS/gofmt.log"
cat "$LOGS/gofmt.log"
/d/work/base/gopath/bin/gofumpt -l $FILES > "$LOGS/gofumpt.log" 2>&1
echo "gofumpt -l (absolute path) rc=$? body=$(wc -c < "$LOGS/gofumpt.log") bytes" >> "$LOGS/gofumpt.log"
cat "$LOGS/gofumpt.log"
go build ./... > "$LOGS/build.log" 2>&1
echo "go build ./... rc=$?" >> "$LOGS/build.log"
GOOS=linux GOARCH=amd64 go vet ./internal/winsec/ > "$LOGS/vet-linux.log" 2>&1
echo "GOOS=linux go vet rc=$?" >> "$LOGS/vet-linux.log"
GOOS=darwin GOARCH=arm64 go vet ./internal/winsec/ > "$LOGS/vet-darwin.log" 2>&1
echo "GOOS=darwin GOARCH=arm64 go vet rc=$?" >> "$LOGS/vet-darwin.log"
for l in build vet-linux vet-darwin; do echo "-- $l"; tail -n 1 "$LOGS/$l.log"; done
sh scripts/d22scan.sh > "$LOGS/d22scan.log" 2>&1
echo "sh scripts/d22scan.sh rc=$?" >> "$LOGS/d22scan.log"
tail -n 12 "$LOGS/d22scan.log"
sh scripts/check-path-length-budget.sh > "$LOGS/path-length.log" 2>&1
echo "sh scripts/check-path-length-budget.sh rc=$?" >> "$LOGS/path-length.log"
tail -n 6 "$LOGS/path-length.log"

echo; echo "== the product file this leg must NOT have touched =="
echo "repo winsec_other.go sha256 = $(sha256sum internal/winsec/winsec_other.go | cut -d' ' -f1)"
echo "expect 137442cdecbb77701c7029c42960593016a152fc965679613314b764c13e2313"
git status --porcelain -- internal/winsec
