# 票 154 AC#4 —— 零字节尺：**按 commit 量，不按区间量**（派单明写的坑）。
# 用法：bash .scratch/wisp/probes/154/zero-byte-per-commit.sh <sha> [<sha> …]
set -u
cd "$(git rev-parse --show-toplevel)" || exit 9
# 冻结/零字节面＝派单"硬约束 AC#4"那串 ∪ 票面 AC#4 那串（去重后逐条列在下面）
FREEZE='^(internal/(risk|panel|agent|observe|speech|secret)/|internal/observe/thresholds\.go$|tools/d22scan/(allowlist\.txt|.*\.go)$|scripts/slo-check\.ps1$|docs/PLAN\.md$|docs/specs/|cmd/wisp/slo_windows\.go$|\.gitattributes$|go\.sum$|.*testdata/golden/.*$|frontend/|design/)'
for sha in "$@"; do
  echo "########## $sha  $(git log -1 --format='%s' "$sha" | cut -c1-40)"
  if [ "$(git cat-file -t "$sha" 2>/dev/null)" != commit ]; then echo "  !! 不是 commit，跳过"; continue; fi
  files="$(git show --pretty=format: --name-only "$sha" | sed '/^$/d')"
  echo "  本枚碰到的文件（全量）:"
  echo "$files" | sed 's/^/    /'
  hits="$(echo "$files" | grep -E "$FREEZE" || true)"
  if [ -z "$hits" ]; then
    echo "  冻结面命中 = 0  ✓"
  else
    echo "  冻结面命中（违规候选）:"
    echo "$hits" | sed 's/^/    !! /'
  fi
done
echo
echo "########## 反向自核：本票写面允许的三处，各碰了几枚"
echo "  internal/tools/bridge.go 的注释面 → 由 AC#5 那条『增删行全部以 // 开头』的机器核担保"
