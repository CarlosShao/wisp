#!/bin/sh
# 322-r1 gates. Every ruler writes its own rc= into its own log; gofmt/gofumpt are read by
# LISTED FILES, not by exit code (their rc is 0 even when a file is broken - the empty-ruler
# accident is booked as A996 and 120-r3 hit it itself).
#
# MSYS_NO_PATHCONV is NOT exported here: leaking it to the host rulers made d22scan read
# `D:\d\work\...missing go.mod` and exit 2 for 120-r3. It is set per docker call only.
REPO='D:/work/workspace/projects plans/Wisp'
LOGS="$REPO/.scratch/wisp/probes/322/r1/logs"
mkdir -p "$LOGS"
cd "$REPO" || exit 1
export GOTMPDIR='C:/Users/swq/AppData/Local/wisp-test-gotmp/322r1gotmp'
mkdir -p "$GOTMPDIR"

MINE='internal/winsec/winsec.go internal/winsec/winsec_other.go internal/winsec/winsec_windows.go internal/winsec/openleaf_322_other_test.go internal/winsec/openleaf_overwrite_322_test.go internal/winsec/openleaf_guard_322_windows_test.go'
GOFUMPT=/d/work/base/gopath/bin/gofumpt.exe
: >"$LOGS/gates-summary.txt"

say() { echo "$@" | tee -a "$LOGS/gates-summary.txt"; }

ruler() {
	name=$1
	shift
	rc=0
	"$@" >"$LOGS/$name.log" 2>&1 || rc=$?
	echo "$name rc=$rc" | tee -a "$LOGS/gates-summary.txt"
}

# ---------- host rulers ----------
ruler build-go-all go build ./...
say "  build ./... output bytes=$(wc -c <"$LOGS/build-go-all.log" | tr -d ' ') (0 = silent success)"

ruler vet-linux go vet ./internal/winsec/
ruler vet-darwin-env env GOOS=darwin GOARCH=arm64 go vet ./internal/winsec/
ruler vet-windows-env env GOOS=windows GOARCH=amd64 go vet ./internal/winsec/

# gofmt / gofumpt: the READING is the list of files, printed in full both scopes.
for scope in mine all-winsec; do
	if [ "$scope" = mine ]; then
		files=$MINE
	else
		files=internal/winsec
	fi
	rc=0
	gofmt -l $files >"$LOGS/gofmt-$scope.log" 2>&1 || rc=$?
	n=$(grep -c '[^[:space:]]' "$LOGS/gofmt-$scope.log" || true)
	say "gofmt -l [$scope] rc=$rc listed-files=$n"
	say "  names: $(tr '\n' ' ' <"$LOGS/gofmt-$scope.log")"
	rc=0
	"$GOFUMPT" -l $files >"$LOGS/gofumpt-$scope.log" 2>&1 || rc=$?
	n=$(grep -c '[^[:space:]]' "$LOGS/gofumpt-$scope.log" || true)
	say "gofumpt -l [$scope] rc=$rc listed-files=$n"
	say "  names: $(tr '\n' ' ' <"$LOGS/gofumpt-$scope.log")"
done

ruler d22scan sh scripts/d22scan.sh
tail -3 "$LOGS/d22scan.log" | tee -a "$LOGS/gates-summary.txt" >/dev/null
say "  d22scan tail: $(tail -1 "$LOGS/d22scan.log")"

ruler path-length-budget sh scripts/check-path-length-budget.sh --with-self-test
say "  budget tail: $(tail -1 "$LOGS/path-length-budget.log")"

# ---------- container: consumer colour, base vs fixed ----------
S=/d/tmp/probe322r1
FIXED=$(ls -d "$S"/fixed.* | head -1)
BASE=$(ls -d "$S"/base.* | head -1)
CONS='./internal/risk ./internal/secret ./internal/config ./internal/agent ./internal/memory'
consumerrun() {
	label=$1
	tree=$2
	pkgs=$3
	rc=0
	MSYS_NO_PATHCONV=1 docker run --rm -e CGO_ENABLED=0 -e GOPROXY=off \
		-v "$tree":/work -v probe322r1build:/root/.cache/go-build -v probe322mod:/gomod:ro \
		-w /work golang:1.27 \
		sh -c 'export GOMODCACHE=/gomod GOTMPDIR=/tmp/gotmp322r1; mkdir -p "$GOTMPDIR"; go test -count=1 '"$pkgs"' > /work/'"$label"'.log 2>&1; rc=$?; echo "go-test rc=$rc" >> /work/'"$label"'.log; exit $rc' \
		>"$LOGS/$label.docker.log" 2>&1 || rc=$?
	cp "$tree/$label.log" "$LOGS/$label.body.log" 2>/dev/null || echo "NO BODY (see $LOGS/$label.docker.log)"
	say "=== $label rc=$rc tree=$tree"
	say "$(tail -n +1 "$LOGS/$label.body.log" 2>/dev/null | grep -E '^(ok|FAIL|---|\s)' | head -20 | tr '\n' '|')"
}
consumerrun cons-base "$BASE" "$CONS"
consumerrun cons-fixed "$FIXED" "$CONS"

# cmd/wisp: directed -run only (ticket 314 booked a whole-process death zone on that face,
# and 165-r2 is live there in the working tree, so this runs in the SNAPSHOT, not the repo).
consumerrun cmdwisp-fixed "$FIXED" '-run "TestAC3POSIXSecretRouteLinkInsideItsDataRootStillRefused119|TestAC1POSIXSecretRouteSymlinked" ./cmd/wisp'
consumerrun cmdwisp-base "$BASE" '-run "TestAC3POSIXSecretRouteLinkInsideItsDataRootStillRefused119|TestAC1POSIXSecretRouteSymlinked" ./cmd/wisp'

say "### gates done $(date '+%F %T %z')"
