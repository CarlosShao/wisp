#!/usr/bin/env bash
# 162-v3 变异台件：复算"把 fs.edit 声明档位从 L2 改成 L1，哪些 AC#5 用例翻红"
# 只裁不改——盘上文件一律不动；变异经 -overlay 落仓外 $TMPDIR 丢弃树。
# 锚点=fd8201e7（162-r4 交件枚）。⚠ -overlay 与 -cover* 不许同用。
set -euo pipefail
cd "$(git rev-parse --show-toplevel)"
ANCHOR=fd8201e7
REPOABS="$(pwd -W 2>/dev/null || pwd)"
DISC="$(mktemp -d)"                      # 临时样本只落仓外丢弃树，只建不删
git show "$ANCHOR":internal/tools/fs_edit.go > "$DISC/fs_edit_L1.go"

# 承重行唯一性自检（防静默 no-op）：目标行必须恰 1 命中，否则拒绝继续
python - "$DISC/fs_edit_L1.go" <<'PY'
import sys
p=sys.argv[1]; s=open(p,encoding='utf-8').read()
line="\t\tDeclared:     risk.L2,"
assert s.count(line)==1, ("target line count!=1", s.count(line))
open(p,"w",encoding="utf-8").write(s.replace(line,"\t\tDeclared:     risk.L1,"))
PY

cat > "$DISC/overlay.json" <<EOF
{"Replace": {"$REPOABS/internal/tools/fs_edit.go": "$DISC/fs_edit_L1.go"}}
EOF

echo "DISCARD=$DISC"
echo "=== 基线（无 overlay，四枚应全绿） ==="
go test -count=1 -v -run 'TestFSEditAndFSWriteShareTheOutOfScopeVerdict|TestFSEditRoutesToApprovalNotTheL1Window|TestFSEditOutOfScopeIsJudgedFromTheDeclaredPathNotTheToolWord|TestFSEditRefusesAPathC26CannotCanonicalize' ./internal/tools/ 2>&1 \
  | grep -E '^--- (PASS|FAIL)' || true
echo "=== 叠加 Declared L2->L1（预期：并排枚仍 PASS=R2 白送；路由枚 FAIL、空白 path 枚 FAIL） ==="
go test -count=1 -v -overlay "$DISC/overlay.json" -run 'TestFSEditAndFSWriteShareTheOutOfScopeVerdict|TestFSEditRoutesToApprovalNotTheL1Window|TestFSEditOutOfScopeIsJudgedFromTheDeclaredPathNotTheToolWord|TestFSEditRefusesAPathC26CannotCanonicalize' ./internal/tools/ 2>&1 \
  | grep -E '^--- (PASS|FAIL)' || true
echo "=== 冻结文本档位挂钩是否被任何用例钉：无 Go 码运行时读 docs ⇒ 改 PLAN.md/SPEC-07 的 L2->L1 零用例红 ==="
git -c core.quotePath=false grep -nE 'ReadFile\(|os\.Open\(|go:embed|ParseFS|ReadDir' "$ANCHOR" -- 'internal/tools/**_test.go' 'internal/risk/**' | grep -iE 'plan|spec|docs' && echo "有读取" || echo "0 命中＝契约文本无运行时读取"
