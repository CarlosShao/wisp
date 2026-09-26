#!/usr/bin/env bash
# ticket 158 r2, AC#4 zero-byte census.
# Ruler shape is fixed by dispatch 2026-09-26-194x section 3:
#   - the commit SET of this round, selected by subject prefix, then walked
#     with `git log --no-walk` (NEVER a range diff as the census, NEVER
#     `git log --since`; a time window sweeps in the orchestrator's own commits
#     and this repo has already been fooled by that once today).
#   - --no-walk is load-bearing: without it git walks the ancestor chain and
#     swallows the whole repo history (measured in section 2, this round).
# Read-only: this script writes nothing into the repo except stdout.
#
# Revision note (mine, same round): the first draft of section 5 matched numstat
# lines with a pattern containing a LITERAL TAB. Whether that tab survives the
# Write tool is not something a reader should have to trust, so section 5 and its
# positive control in section 8 now share ONE tab-free ruler: `cut -f3` (numstat
# is tab-separated) fed into an anchored filename regex.
set -u
cd "$(git rev-parse --show-toplevel)" || exit 9

ANCHOR="${1:-bb08d1f}"                     # this round's step-0 anchor
PFX='票 158 r2'                            # subject prefix that names THIS round

# The frozen faces of AC#4, plus the three the dispatch added (docs/reports, this
# ticket's own face, ticket 154's gate body). One regex, shared by section 5 (the
# reading) and section 8 (its control), so the two can never drift apart.
# internal/tools/bridge.go is IN the pattern on purpose: the authorized comment
# face has to show up as a HIT instead of being quietly excluded by the ruler.
FIRE='^(docs/PLAN\.md|docs/specs/|internal/(risk|panel|agent|observe)/|internal/tools/bridge\.go$|tools/d22scan/|frontend/|design/|.*golden|.*allowlist\.txt$|scripts/slo-check\.ps1$|docs/evidence/s1/15[146]-|\.scratch/wisp/issues/15[146]-)'

echo "# AC#4 census - anchor for commit selection: $ANCHOR ($(git log -1 --format='%h %ad %s' --date=format:'%H:%M:%S' "$ANCHOR"))"
echo "# tip now: $(git rev-parse --short HEAD) ($(git log -1 --format='%ad %s' --date=format:'%H:%M:%S'))"
echo

# ---- 1. the commit set, and proof of what each selection rule yields -------
echo "## 1. 本程 commit 号集合（按 subject 前缀选，不是按时间窗）"
SET=$(git log --format='%H' "$ANCHOR..HEAD" --grep="^$PFX")
echo "\$ git log --format='%H' $ANCHOR..HEAD --grep='^$PFX'"
for s in $SET; do git log -1 --format='  %h %ad %s' --date=format:'%H:%M:%S' "$s"; done
echo "# 枚数=$(echo $SET | wc -w)"
echo "# 反向自证：同一区间里【不是】本程的提交（选规则没有多吞/漏吞）："
git log --format='  %h %ad %s' --date=format:'%H:%M:%S' "$ANCHOR..HEAD" --invert-grep --grep="^$PFX"
echo "# 每一枚的 --no-walk 自证（都是 commit，不是 tree/blob）："
for s in $SET; do echo "  $(git cat-file -t "$s") $s"; done
echo

# ---- 2. the touched-path roster, per commit, with --no-walk ---------------
echo "## 2. 名册尺（主尺）：git log --no-walk --pretty=tformat: --name-only <逐枚号>"
echo "\$ git log --no-walk --pretty=tformat: --name-only $SET | sed '/^$/d' | LC_ALL=C sort | uniq -c"
git log --no-walk --pretty=tformat: --name-only $SET | sed '/^$/d' | LC_ALL=C sort | uniq -c
echo "# 名册去重后枚数=$(git log --no-walk --pretty=tformat: --name-only $SET | sed '/^$/d' | LC_ALL=C sort -u | wc -l)"
echo "# 派单点名的那枚坑，本程自己量一遍（同三枚号，只差 --no-walk）："
echo "  带 --no-walk   = $(git log --no-walk --pretty=tformat: --name-only $SET | sed '/^$/d' | wc -l) 行"
echo "  不带 --no-walk = $(git log --pretty=tformat: --name-only $SET | sed '/^$/d' | wc -l) 行"
echo

# ---- 3. per-commit name-status + numstat ----------------------------------
echo "## 3. 逐枚 name-status + numstat（谁碰了什么、动了几行）"
for s in $SET; do
  git log -1 --format='### %h %s' "$s"
  git show --format='' --name-status "$s" | sed '/^$/d'
  echo "  -- numstat 合计: added=$(git show --format='' --numstat "$s" | awk '{a+=$1} END {print a+0}') deleted=$(git show --format='' --numstat "$s" | awk '{d+=$2} END {print d+0}')"
done
echo

# ---- 4. the frozen faces, one by one, from the per-commit ruler -----------
# Each entry: label | pathspec passed to `git log --no-walk ... -- <pathspec>`
FROZEN='docs/PLAN.md|docs/PLAN.md
docs/specs/**|docs/specs
internal/risk/**|internal/risk
internal/panel/**|internal/panel
internal/agent/** (incl approval/**)|internal/agent
internal/observe/**|internal/observe
thresholds.go|internal/observe/thresholds.go
golden: internal/llm/golden/**|internal/llm/golden
golden: any tracked path with golden in it|:(glob)**/*golden*
golden: any testdata/golden/**|:(glob)**/testdata/golden/**
allowlist.txt|tools/d22scan/allowlist.txt
scripts/slo-check.ps1|scripts/slo-check.ps1
tools/d22scan/**|tools/d22scan
frontend/**|frontend
design/**|design
evidence 151 r1|docs/evidence/s1/151-task-scope-never-closed-r1.md
evidence 151 accept|docs/evidence/s1/151-task-scope-never-closed-r1-accept-r1.md
evidence 154 host-id|docs/evidence/s1/154-host-id-never-closed-r1.md
evidence 154 accept|docs/evidence/s1/154-close-gate-never-rings-r1-accept-r1.md
evidence 156 r1|docs/evidence/s1/156-exited-asks-os-r1.md
evidence 156 accept|docs/evidence/s1/156-close-and-load-bearing-r1-accept-r1.md
evidence 156 r4|docs/evidence/s1/156-r4-gates-and-cleanup-r1.md
ticket face 151|:(glob).scratch/wisp/issues/151-*
ticket face 154|:(glob).scratch/wisp/issues/154-*
ticket face 156|:(glob).scratch/wisp/issues/156-*
本票票面（编排者写面，非 AC#4 名单但同样应零枚）|:(glob).scratch/wisp/issues/158-*
台账/停车点（编排者写面）|docs/reports
门本体（r1 写面，本程应零枚）|.scratch/wisp/probes/154
WRITEFACE 本票证据件（应有命中）|docs/evidence/s1/158-gate-scope-blind-spot-r1.md
WRITEFACE 本程探针件（应有命中）|.scratch/wisp/probes/158
WRITEFACE bridge.go（授权注释面，应有命中且只注释）|internal/tools/bridge.go'

echo "## 4. 逐支冻结面：本程 commit 号集合上再套一条 pathspec（命中即越界）"
printf '%s\n' "$FROZEN" | while IFS='|' read -r label spec; do
  [ -n "$label" ] || continue
  hits=$(git log --no-walk --pretty=tformat: --name-only $SET -- "$spec" | sed '/^$/d' | LC_ALL=C sort -u)
  n=$(printf '%s' "$hits" | grep -c . || true)
  tracked=$(git ls-files -- "$spec" | wc -l | tr -d ' ')
  printf '  越界枚数=%-3s (该面已跟踪文件数=%-4s)  %s\n' "$n" "$tracked" "$label"
  [ "$n" != "0" ] && printf '%s\n' "$hits" | sed 's/^/      HIT /'
done
echo "# 冻结面每一支都＝0 ⇒ AC#4 的逐支账；三支 WRITEFACE 有命中，那是本程写面本体"
echo

# ---- 5. second caliber: byte-level, same ruler as the control in section 8 --
echo "## 5. 第二口径（字节级：把 numstat 的文件名列打上面那串正则）"
echo "\$ git log --no-walk --pretty=tformat: --numstat $SET | cut -f3 | grep -E \"$FIRE\""
git log --no-walk --pretty=tformat: --numstat $SET | sed '/^$/d' | cut -f3 \
  | grep -E "$FIRE" | LC_ALL=C sort -u | sed 's/^/  HIT /'
echo "  # 预期：只有 bridge.go 那一枚命中（授权面）。其余冻结面全零 ⇒ 与 §4 同判"
echo

# ---- 6. the working-tree side --------------------------------------------
echo "## 6. 工作树侧（本程还没 commit 的改动）"
git status --porcelain -- docs/evidence/s1/158-gate-scope-blind-spot-r1.md .scratch/wisp/probes/158 internal/tools cmd/wisp
echo "  # 这里出现的都属本程写面；别人的脏文件不在列（口径见证据件 §C.4）"
echo

# ---- 7. third caliber: range diff, explicitly NOT the census ruler --------
echo "## 7. 第三把尺：区间 diff（派单 §3 明令它【不能】当普查主尺，这里只当反向对照）"
echo "\$ git diff --name-only $ANCHOR HEAD | grep -v '^\.scratch/wisp/probes/158/r2b/'"
git diff --name-only "$ANCHOR" HEAD | grep -v '^\.scratch/wisp/probes/158/r2b/' | sed 's/^/  /'
echo "  # 区间尺把【编排者同期提交】也扫进来了。下面逐枚归名（这些号都不在 §1 的本程集合里）："
for s in $(git log --format=%H "$ANCHOR..HEAD" --invert-grep --grep="^$PFX"); do
  git show --format='  --- %h %s' --name-only "$s" | sed '/^$/d'
done
echo "  # ⇒ 这就是普查必须按 commit 号集合算的理由：区间尺会把别人写成我的。"
echo

# ---- 8. positive control: prove the section-5 ruler can fire at all --------
echo "## 8. §5 那把尺的正控（同一串正则，打在【已知动过冻结面】的三枚真提交上，必须响）"
for c in 12181923 00bbb76f cbbbdf85; do
  short=$(git rev-parse --short "$c" 2>/dev/null || echo MISSING)
  hits=$(git log -1 --pretty=tformat: --numstat "$c" | sed '/^$/d' | cut -f3 | grep -E "$FIRE" | LC_ALL=C sort -u)
  n=$(printf '%s' "$hits" | grep -c . || true)
  printf '  %s 命中=%s  %s\n' "$short" "$n" "$(git log -1 --format='%s' "$c" | cut -c1-60)"
  [ "$n" != "0" ] && printf '%s\n' "$hits" | sed 's/^/      FIRE /'
done
echo "  # 三枚全命中 ⇒ §5 那个『只剩 bridge.go』是『码上真的没有』，不是尺没通电。"
