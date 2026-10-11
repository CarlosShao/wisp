#!/bin/sh
# 72-r1 two-surface runner (out-of-repo exported tree; no worktree, no clone).
# usage: sh run-two-surface.sh <commit-ref> <tree-dir> <summary-log>
set -u
REPO="/d/work/workspace/projects plans/Wisp"
REF="$1"
TREE="$2"
OUT="$3"
TAG="$4"
LOGS="$REPO/.scratch/wisp/probes/72/r1/logs"
PROBE="/c/Users/swq/tmp/w72r1/probe"

R7='TestClassifyAnchorSpellingIsNotVerdict|TestCanonicalInputGainsNoSecondForm|TestAListWinsWhereBothTablesHit|TestPathResolverShortNameAListDenied|TestPathResolverUNCAListDenied|TestPathResolverExtendedLengthPrefixAListDenied|TestBListDefaultDenyAndOverride'
R16='TestPathResolver|TestReparse|TestAList|TestBList|TestResolve|TestClassify|TestCanonical|TestOverride'
ALIAS='C:\Users\swq\tmp\RISKAT~1\RUNNER~1'
MACH='C:\Users\swq\AppData\Local\Temp'

emit() { # label tempvalue runregex rawlog
  label="$1"; tv="$2"; roster="$3"; log="$4"
  TEMP="$tv" TMP="$tv" TMPDIR="$tv" go test ./internal/risk/ -count=1 -timeout 600s -v -run "$roster" > "$log" 2>&1
  rc=$?
  echo "---- $label rc=$rc temp=$tv raw=$(basename "$log")"
  echo "     RUN=$(grep -c '=== RUN' "$log") PASS=$(grep -c '^    --- PASS\|^--- PASS' "$log") FAIL=$(grep -c '^    --- FAIL\|^--- FAIL' "$log") SKIP=$(grep -c '^    --- SKIP\|^--- SKIP' "$log")"
  grep -E '^ *--- FAIL' "$log" | sed 's/^/     FAILNAME: /'
  echo "     REDPOS: $(grep -oE '[A-Za-z0-9_]+\.go:[0-9]+:' "$log" | sort -u | tr '\n' ' ')"
}

{
  echo "# 72-r1 two-surface start $(date '+%Y-%m-%d %H:%M:%S %z') ref=$REF tree=$TREE"
  rm -rf "$TREE" 2>/dev/null; echo "rm_rf_rc=$? (own out-of-repo tree)"
  mkdir -p "$TREE" || { echo "tree_mkdir_rc=$?"; exit 1; }
  cd "$REPO" || exit 1
  git archive "$REF" go.mod go.sum internal cmd | tar -x -C "$TREE"
  echo "archive_rc=$?  tree_top=$(ls "$TREE" | tr '\n' ' ')  tree_go_files=$(find "$TREE" -name '*.go' | wc -l)"
  echo "repo_dirty_internal_cmd=$(git status --porcelain -- internal cmd | wc -l)"
  echo "=== surface probe (same API t.TempDir uses, out-of-repo module) ==="
  mkdir -p "$PROBE"
  printf 'package main\n\nimport (\n\t"fmt"\n\t"os"\n\t"strings"\n)\n\nfunc main() {\n\td, err := os.MkdirTemp("", "W72R1PROBE")\n\tfmt.Printf("tempdir_probe env TEMP=%%s TMP=%%s TMPDIR=%%s os.TempDir()=%%s created=%%s hasTilde=%%v err=%%v\\n", os.Getenv("TEMP"), os.Getenv("TMP"), os.Getenv("TMPDIR"), os.TempDir(), d, strings.Contains(d, "~1"), err)\n}\n' > "$PROBE/main.go"
  [ -f "$PROBE/go.mod" ] || (cd "$PROBE" && go mod init w72r1probe >/dev/null 2>&1)
  echo "  ALIASED: $(cd "$PROBE" && TEMP="$ALIAS" TMP="$ALIAS" TMPDIR="$ALIAS" go run . 2>&1 | tr -d '\r')"
  echo "  MACHINE: $(cd "$PROBE" && TEMP="$MACH" TMP="$MACH" TMPDIR="$MACH" go run . 2>&1 | tr -d '\r')"
  cd "$TREE" || exit 1
  emit "alias-7named" "$ALIAS" "$R7" "$LOGS/tmp-$TAG-alias-7.txt"
  emit "alias-16all" "$ALIAS" "$R16" "$LOGS/tmp-$TAG-alias-16.txt"
  emit "machine-7named" "$MACH" "$R7" "$LOGS/tmp-$TAG-machine-7.txt"
  emit "machine-16all" "$MACH" "$R16" "$LOGS/tmp-$TAG-machine-16.txt"
  echo "# done $(date '+%Y-%m-%d %H:%M:%S %z')"
} >> "$OUT" 2>&1
cat "$OUT"
