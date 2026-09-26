#!/usr/bin/env bash
# Ticket 153 AC#5: the gate, run per package (never ./...), pre and post.
#   PRE  = snapshot of the anchor 86b0161, out of repo
#   POST = the working tree at the delivered HEAD
set -u
REPO="/d/work/workspace/projects plans/Wisp"
HELP="/c/Users/swq/AppData/Local/Temp/wisp153"
PRESNAP="/c/Users/swq/AppData/Local/Temp/wisp153-snap-pre"
LOG="$HELP/logs"; mkdir -p "$LOG" "$PRESNAP"
GOFUMPT=/d/work/base/gopath/bin/gofumpt.exe

four() { # <file> -> the four numbers + unique names
  printf 'RUN(all)=%s PASS_all=%s FAIL=%s SKIP=%s panic=%s unique=%s\n' \
    "$(grep -c '^=== RUN' "$1")" \
    "$(grep -cE '^ *--- PASS' "$1")" \
    "$(grep -cE '^ *--- FAIL' "$1")" \
    "$(grep -cE '^ *--- SKIP' "$1")" \
    "$(grep -c 'panic:' "$1")" \
    "$(grep -E '^ *--- (PASS|FAIL|SKIP)' "$1" | awk '{print $3}' | sort -u | wc -l)"
}
roster() { grep -E '^ *--- (PASS|FAIL|SKIP)' "$1" | awk '{print $3}' | sort -u; }

cd "$REPO" || exit 1
echo "== toolchain versions read live (not recalled)"
go version
"$GOFUMPT" --version
gofmt --version 2>&1 | head -1 || true
git rev-parse --short HEAD

echo
echo "== PRE: snapshot of the anchor, out of repo"
git -c core.autocrlf=false -c core.eol=lf archive --format=tar 86b0161 | tar -x -C "$PRESNAP"
cmp -s "$PRESNAP/internal/agent/compress.go" <(git show 86b0161:internal/agent/compress.go) \
  && echo "snapshot internal/agent/compress.go == anchor"
( cd "$PRESNAP" && go test ./internal/agent/ -count=1 -v > "$LOG/ac5-pre-count1.txt" 2>&1 ); echo "pre rc=$?"
four "$LOG/ac5-pre-count1.txt"

echo
echo "== POST: working tree, -count=1"
go test ./internal/agent/ -count=1 -v > "$LOG/ac5-post-count1.txt" 2>&1; echo "post rc=$?"
four "$LOG/ac5-post-count1.txt"
echo
echo "== POST: -count=2 (arithmetic self-check: RUN should be 2x)"
go test ./internal/agent/ -count=2 -v > "$LOG/ac5-post-count2.txt" 2>&1; echo "post-count2 rc=$?"
four "$LOG/ac5-post-count2.txt"

echo
echo "== roster diff, both directions"
roster "$LOG/ac5-pre-count1.txt"   > "$LOG/ac5-roster-pre.txt"
roster "$LOG/ac5-post-count1.txt"  > "$LOG/ac5-roster-post.txt"
roster "$LOG/ac5-post-count2.txt"  > "$LOG/ac5-roster-post2.txt"
echo "-- new (in post, not in pre):"; comm -13 "$LOG/ac5-roster-pre.txt" "$LOG/ac5-roster-post.txt" | sed 's/^/   /'
echo "-- lost (in pre, not in post):"; comm -23 "$LOG/ac5-roster-pre.txt" "$LOG/ac5-roster-post.txt" | sed 's/^/   /'
echo "   (empty above = nothing lost)"
echo "-- count=1 vs count=2 name lists:"
diff -q "$LOG/ac5-roster-post.txt" "$LOG/ac5-roster-post2.txt" && echo "   NAMES IDENTICAL"
echo "-- opened vs adjudicated in the POST run: RUN=$(grep -c '^=== RUN' "$LOG/ac5-post-count1.txt") adjudicated=$(grep -cE '^ *--- (PASS|FAIL|SKIP)' "$LOG/ac5-post-count1.txt")"

echo
echo "== go vet / gofmt / gofumpt (internal/agent only)"
go vet ./internal/agent/ ; echo "vet rc=$?"
echo "gofmt -l:"; gofmt -l internal/agent; echo "  (empty = clean)"
echo "gofumpt -l:"; "$GOFUMPT" -l internal/agent; echo "  (empty = clean)"
echo "gofumpt -l . tools/d22scan tools/mockllm:"; "$GOFUMPT" -l . tools/d22scan tools/mockllm | head; echo "  (empty = clean)"

echo
echo "== -race (attribution is per-call, so this is the race surface that matters)"
go test ./internal/agent/ -count=1 -race > "$LOG/ac5-post-race.txt" 2>&1; echo "race rc=$?"; tail -3 "$LOG/ac5-post-race.txt"

echo
echo "== d22scan"
sh scripts/d22scan.sh > "$LOG/ac5-d22scan.txt" 2>&1; echo "d22scan rc=$?"
echo "-- every 'examined' line in order (the head ones are tools/d22scan's OWN fixtures):"
grep -n 'examined' "$LOG/ac5-d22scan.txt" | sed 's/^/   /'
echo "-- verdict block:"
tail -12 "$LOG/ac5-d22scan.txt" | sed 's/^/   /'
