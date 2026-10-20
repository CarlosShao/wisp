#!/bin/sh
# gw-1 味① 第 4 把的**第三形**＝CI-shape 复量：共享树的工作树字节带本机 CRLF⇒本地 (A) 尺读到 35 行；
# 而 `git archive` 出的导出树写的是**入库 blob 字节**（LF）＝CI checkout 看到的同一形。
# ⇒ 在同一把尺（gofumpt -l，⛔ -w）上量"CI 那形"＝把生产面四枚具名件的字节色取出来。
# ★⛔ 洗任何件、★射程＝仓外导出树 gwL（锚＝a19d5bf2）。
set -u
R=/d/tmp/wisp/20261020/gw1
LOGS="$R/logs"
cd '/d/work/workspace/projects plans/Wisp' || exit 9
unset GOTMPDIR TMPDIR TEMP TMP
echo "== START anchor =="
git rev-parse HEAD | sed 's|^|HEAD_FULL=|'
T="$R/gwL"
echo "TREE=$T  GO_FILES=$(find "$T" -name '*.go' | wc -l)"
echo
echo "== 1) 生产面四枚：blob 字节 ⇄ 工作树字节 的 CR 尺（对**入库 blob**量＝tr -cd）=="
for p in cmd/wisp/models.go internal/agent/approval/pending_read.go internal/agent/tools.go internal/risk/provenance.go; do
  echo "  $p worktree_cr=$(tr -cd '\r' < "$p" | wc -c) blob_cr=$(git cat-file blob "HEAD:$p" | tr -cd '\r' | wc -c)"
done
echo "  （正控＝自造 2 枚 CR 件同尺读出 $(tr -cd '\r' < "$R/posctrl/cr2.txt" | wc -c)）"
echo
echo "== 2) 导出树里同名的四枚（＝blob 字节形）打 gofumpt -l ⇄ gofmt -l =="
( cd "$T" && /d/work/base/gopath/bin/gofumpt.exe -l cmd/wisp/models.go internal/agent/approval/pending_read.go internal/agent/tools.go internal/risk/provenance.go ) > "$LOGS/g4-export-four.txt" 2>&1
echo "GOFUMPT_EXPORT_RC=$?"
echo "GOFUMPT_EXPORT_LISTED=$(grep -c . "$LOGS/g4-export-four.txt" || true)"
sed 's|^|  EX |' "$LOGS/g4-export-four.txt"
( cd "$T" && gofmt -l cmd/wisp/models.go internal/agent/approval/pending_read.go internal/agent/tools.go internal/risk/provenance.go ) > "$LOGS/g4-export-four-gofmt.txt" 2>&1
echo "GOFMT_EXPORT_RC=$?"
echo "GOFMT_EXPORT_LISTED=$(grep -c . "$LOGS/g4-export-four-gofmt.txt" || true)"
sed 's|^|  EX |' "$LOGS/g4-export-four-gofmt.txt"
echo
echo "== 3) 共享树工作字节同四枚（对照发＝证明 35 行那把尺里这 4 枚＝CRLF 形而非语义脏）=="
/d/work/base/gopath/bin/gofumpt.exe -l cmd/wisp/models.go internal/agent/approval/pending_read.go internal/agent/tools.go internal/risk/provenance.go > "$LOGS/g4-worktree-four.txt" 2>&1
echo "GOFUMPT_WORKTREE_RC=$?"
echo "GOFUMPT_WORKTREE_LISTED=$(grep -c . "$LOGS/g4-worktree-four.txt" || true)"
sed 's|^|  WT |' "$LOGS/g4-worktree-four.txt"
echo
echo "== 4) 导出树整面尺（＝CI 那形打在整棵八件套树上：gofumpt -l . tools/d22scan tools/mockllm）=="
echo "   ⚠ 射程具名＝导出树**不含** .scratch/**⛔ 未跟踪件（archive 只取锚里那批八件套）⇒ 该把尺的 31 枚 bench 脏行在这里结构性不存在＝⛔ 拿这发的 0 顶替 (A) 那把"
( cd "$T" && /d/work/base/gopath/bin/gofumpt.exe -l . tools/d22scan tools/mockllm ) > "$LOGS/g4-export-all.txt" 2>&1
echo "GOFUMPT_ALL_EXPORT_RC=$?"
echo "GOFUMPT_ALL_EXPORT_LISTED=$(grep -c . "$LOGS/g4-export-all.txt" || true)"
sed 's|^|  EA |' "$LOGS/g4-export-all.txt"
echo
echo "== END anchor =="
git rev-parse HEAD | sed 's|^|HEAD_FULL_END=|'
