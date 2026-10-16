#!/bin/sh
# run-mut.sh <op> - one named mutation of the EXPORT TREE (仓外), full internal/tools
# run against the physical bytes, bounded extraction into logs/, restore + sha256 proof.
# Never touches the working tree. Raw logs stay in the export tree.
set -eu
REPO="/d/work/workspace/projects plans/Wisp"
BENCH="/d/work/workspace/wisp-270r2-comment-bench/repo"
LOGS="$REPO/.scratch/wisp/probes/270/r2/logs"
TARGET="$BENCH/internal/tools/subagent_197.go"
PRISTINE="$TARGET.pristine"
PYMUT="$REPO/.scratch/wisp/probes/270/r2/mutate_270r2.py"
op=${1:?usage: run-mut.sh <half-off|half-off-no-jia|restore>}

cp "$PRISTINE" "$TARGET"
[ "$op" = restore ] || python "$PYMUT" "$TARGET" "$PRISTINE" "$op"
mut_sha=$(sha256sum "$TARGET" | cut -d' ' -f1)

cd "$BENCH"
export PATH="$BENCH/third_party/sherpa-onnx:$PATH"
rc=0
go test ./internal/tools/ -count=1 -v > "raw-mut-$op.txt" 2>&1 || rc=$?

ext="$LOGS/mut-$op.txt"
{
	echo "# op=$op 台面=$BENCH（git archive HEAD 2e706255 的仓外导出树）raw=raw-mut-$op.txt"
	echo "# rc=$rc 突变后 sha256=$mut_sha"
	echo "# 尺: grep -cE '^=== RUN'(含子测试) / '^--- PASS' / '^--- FAIL' / '^--- SKIP'(只数顶层) / '^panic'；射程=internal/tools 单包；量的"
	echo "RUN=$(grep -cE '^=== RUN' "raw-mut-$op.txt")"
	echo "PASS_TOP=$(grep -cE '^--- PASS' "raw-mut-$op.txt")"
	echo "FAIL_TOP=$(grep -cE '^--- FAIL' "raw-mut-$op.txt")"
	echo "SKIP=$(grep -cE '^--- SKIP' "raw-mut-$op.txt")"
	echo "PANIC_LINES=$(grep -cE '^panic' "raw-mut-$op.txt")"
	echo "# 两枚点名的 nil-BaseOptions 腿（谁答的那一句看它们的状态）"
	grep -E '^--- (PASS|FAIL): (Test236R2TaskSpawnRefusesWhenBaseOptionsIsUnwired|Test270R1NilBaseOptionsRefusesWithTheAssemblySentence|Test236R2TaskSpawnRefusesWhenParentToolsIsUnwired) ' "raw-mut-$op.txt" || echo "(点名腿未出现)"
	echo "# 红名册逐枚"
	grep -E '^--- FAIL' "raw-mut-$op.txt" | sed 's/^--- FAIL: //; s/ .*//' | while read -r name; do
		echo "----- $name"
		awk -v n="$name" '$0 ~ "^=== RUN   " n "([[:space:]]|$)" { grab=1 } grab { print } $0 ~ "^--- (FAIL|PASS): " n " " { if (grab) exit }' "raw-mut-$op.txt" \
			| grep -E '^[[:space:]]+[a-z_0-9]+\.go:[0-9]+:' || echo "  (无 message 行)"
	done
} >"$ext" 2>&1

python "$PYMUT" "$TARGET" "$PRISTINE" restore >/dev/null
headsha=$( (cd "$REPO" && git cat-file -p 2e70625521dca0eef138a514145043194b54ffca:internal/tools/subagent_197.go) | sha256sum | cut -d' ' -f1)
restsha=$(sha256sum "$TARGET" | cut -d' ' -f1)
{
	echo "# 还原自证（导出树锚＝archive 的 2e706255；本腿 087936bc 只动整行注释，产码字节与它逐字同形）"
	echo "restored_tree_sha256=$restsha"
	echo "archive_commit_blob_sha256=$headsha"
	[ "$restsha" = "$headsha" ] && echo "MATCH=yes" || echo "MATCH=no"
	echo "# 工作树那枚件本腿从未拿它当台面：git status 对该路径应为空"
	(cd "$REPO" && git status --porcelain -- internal/tools/subagent_197.go | wc -l)
} >>"$ext"
cat "$ext"
