#!/usr/bin/env bash
# 328-r1 去牙突变（硬框 ⓒ）：只打在仓外导出树的物理字节上，内容锚（函数体内那枚 token），
# ⛔ 绝对行号、⛔ -overlay、⛔ 改跟踪文件。每发串行：取 pristine → 断言 pristine 哈希 →
# 突变 → 跑用例 → 当场还原（trap EXIT 再兜一次）→ 再断言哈希同串。
# 用法：mutate.sh <M1|M2>
set -u
TREE="D:/work/workspace/wisp-328r1-after-tree"
F="$TREE/internal/panel/subagent_roster_197.go"
LOGDIR="D:/work/workspace/projects plans/Wisp/.scratch/wisp/probes/328/r1/logs"
REPO="D:/work/workspace/projects plans/Wisp"
export GOPATH='D:\work\base\gopath' GOTMPDIR='D:\tmp\328r1-gotmp'
mkdir -p "$GOTMPDIR"
SHOT="${1:?M1 或 M2}"
PRISTINE="$LOGDIR/subagent_roster_197.go.pristine-$SHOT"
LOG="$LOGDIR/mut-$SHOT.log"
BLOB=$(git -C "$REPO" rev-parse f705b082:internal/panel/subagent_roster_197.go)
cp "$F" "$PRISTINE"
restore() { cp "$PRISTINE" "$F"; }
trap restore EXIT
sha_before=$(sha256sum "$F" | cut -d' ' -f1)
h_before=$(git hash-object "$F")
{
  echo "date=$(date '+%F %T %z') shot=$SHOT tree=$TREE"
  echo "git_blob_at_anchor=$BLOB"
  echo "sha256_pristine=$sha_before"
  echo "git_hash_object_pristine=$h_before"
  echo "pristine_equals_blob_check=$( [ "$h_before" = "$BLOB" ] && echo YES || echo NO )"
} > "$LOG"
case "$SHOT" in
  M1) # 抹掉分隔符那一半 ⇒ 纯前缀匹配
    sed -i 's/strings.HasPrefix(key, task+"#")/strings.HasPrefix(key, task)/' "$F"
    grep -n 'strings.HasPrefix(key, task' "$F" >> "$LOG"
    ( cd "$TREE" && go test -v ./internal/panel/ -count=1 -timeout 30m \
        -run 'TestTheClashFormIsRecognisedOnlyAsTheQueueWritesIt|TestProductionCorrelationShapesLightTheirRows' ) >> "$LOG" 2>&1
    echo "go_test_rc=$?" >> "$LOG"
    ;;
  M2) # 把新臂改回今天的形状：第三臂不登记任何东西
    sed -i 's/set\[task\] = true/_ = task/' "$F"
    grep -n '_ = task' "$F" >> "$LOG"
    ( cd "$TREE" && go test -v ./internal/agent/ -count=1 -timeout 30m \
        -run 'TestACorrelationTheLoopActuallyMintedLightsItsRosterRow|TestTheLoopMintsTheShapesTheRosterAccepts' ) >> "$LOG" 2>&1
    echo "go_test_agent_rc=$?" >> "$LOG"
    ( cd "$TREE" && go test -v ./internal/panel/ -count=1 -timeout 30m \
        -run 'TestProductionCorrelationShapesLightTheirRows|TestTheClashFormIsRecognisedOnlyAsTheQueueWritesIt' ) >> "$LOG" 2>&1
    echo "go_test_panel_rc=$?" >> "$LOG"
    ;;
  *) echo "unknown shot"; exit 2 ;;
esac
restore
trap - EXIT
sha_after=$(sha256sum "$F" | cut -d' ' -f1)
h_after=$(git hash-object "$F")
{
  echo "sha256_restored=$sha_after"
  echo "git_hash_object_restored=$h_after"
  echo "restored_same_as_pristine=$( [ "$sha_after" = "$sha_before" ] && echo YES || echo NO )"
  echo "--- red lines (verbatim) ---"
  grep -E '^(--- FAIL|    [a-z_0-9./]+_test\.go:[0-9]+:)' "$LOG" | head -30 || echo NO_FAIL_LINES
  echo "RERUN_RECONCILE: $(cd "$TREE" 2>/dev/null && echo ok)"
} >> "$LOG"
tail -40 "$LOG"
