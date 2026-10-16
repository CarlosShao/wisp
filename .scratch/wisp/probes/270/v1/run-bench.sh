#!/bin/sh
# run-bench.sh - 270-v1 adjudication bench driver.
#
# One argument per invocation: the op name (see mutate270v1.py OPS) or "base"
# for the unmutated export tree. It restores the pristine bytes first, applies
# exactly one named mutation, runs the FULL internal/tools package (-count=1 -v)
# against the physical file, extracts the four roster numbers plus every red
# name with its verbatim message lines, then restores and proves the restore by
# sha256 against the repo's HEAD blob (git cat-file), outside the tree.
#
# Raw logs stay in the export tree (仓外); the bounded extraction lands in the
# leg's logs/ dir. Nothing here deletes anything.
set -eu

REPO=$(git rev-parse --show-toplevel)
BENCH=$(cat "$REPO/.scratch/wisp/probes/270/v1/logs/bench-path.txt")
LOGS="$REPO/.scratch/wisp/probes/270/v1/logs"
TARGET="$BENCH/internal/tools/subagent_197.go"
PRISTINE="$BENCH/internal/tools/subagent_197.go.pristine"
PYMUT="$REPO/.scratch/wisp/probes/270/v1/mutate270v1.py"
GOTMPROOT="${GOTMPDIR:-}"

export PATH="$REPO/third_party/sherpa-onnx:$PATH"

op=${1:?usage: run-bench.sh <op|base>}

# always start from pristine
cp "$PRISTINE" "$TARGET"
if [ "$op" != "base" ]; then
	python "$PYMUT" "$TARGET" "$op"
fi
mut_sha=$(sha256sum "$TARGET" | cut -d' ' -f1)

raw="$BENCH/raw-$op.txt"
cd "$BENCH"
rc=0
go test ./internal/tools/ -count=1 -v >"$raw" 2>&1 || rc=$?

ext="$LOGS/read-$op.txt"
{
	echo "# op=$op  台面=git archive HEAD 导出的仓外树 $BENCH  raw=$raw"
	echo "# rc=$rc  突变后 sha256=$mut_sha  GOTMPDIR=$GOTMPROOT"
	echo "# 尺: grep -c '^=== RUN' / '^--- PASS' / '^--- FAIL' / '^--- SKIP' / '^panic' 对 raw -v 日志(量的, 射程=internal/tools 单包)"
	echo "RUN=$(grep -c -E '^=== RUN' "$raw")"
	echo "PASS_TOP=$(grep -c -E '^--- PASS' "$raw")"
	echo "FAIL_TOP=$(grep -c -E '^--- FAIL' "$raw")"
	echo "SKIP=$(grep -c -E '^--- SKIP' "$raw")"
	echo "PANIC_LINES=$(grep -c -E '^panic' "$raw")"
	echo "TAIL=$(tail -2 "$raw" | tr '\n' '|')"
	echo "# 红名册逐枚 + 红句原文(从 === RUN 到 --- FAIL 之间的 t.Errorf/Fatalf 行, 逐字)"
	grep -E '^--- FAIL' "$raw" | sed 's/^--- FAIL: //; s/ .*//' | while read -r name; do
		echo "----- $name"
		awk -v n="$name" '
			$0 ~ "^=== RUN   " n "([[:space:]]|$)" { grab = 1 }
			grab { print }
			$0 ~ "^--- (FAIL|PASS): " n " " { if (grab) { exit } }
		' "$raw" | grep -E '^\s+[a-z_0-9]+\.go:[0-9]+:' || echo "  (无 message 行)"
	done
} >"$ext" 2>&1

# restore + prove
python "$PYMUT" "$TARGET" restore "$PRISTINE" >/dev/null
headsha=$( (cd "$REPO" && git cat-file -p HEAD:internal/tools/subagent_197.go) | sha256sum | cut -d' ' -f1)
restsha=$(sha256sum "$TARGET" | cut -d' ' -f1)
{
	echo "# 还原自证"
	echo "restored_tree_sha256=$restsha"
	echo "repo_HEAD_blob_sha256=$headsha"
	[ "$restsha" = "$headsha" ] && echo "MATCH=yes" || echo "MATCH=no"
} >>"$ext"
cat "$ext"
