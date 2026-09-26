#!/usr/bin/env bash
# 票 154 AC#4（续程 r2）—— **逐支**零字节尺：票面 AC#4 点名的每一支禁改面 × 逐枚 commit 各答一次。
# 用法： bash .scratch/wisp/probes/154/ac4-per-face.sh <sha> [<sha> …]
#
# 与前一枚脚本（zero-byte-per-commit.sh，154 实现程 r1 留）的差别，三条，都是它那一把尺的不足：
#   ① 它只印"冻结面命中 = 0"这一枚总数，不逐支答 ⇒ 看不出一支漏答；本脚本一行一支。
#   ② 它用 `git show --name-only`（默认开改名探测）⇒ 一枚"把冻结面文件改名到面外"的 commit 只会露出新名，
#      旧名被吞。本脚本对每枚 commit 同时跑两把名册尺：N1 = --name-only（改名折叠），
#      N2 = --name-status --no-renames（改名拆成 A+D，两侧都在册），并核 N1≡N2。
#   ③ 它的正则把 `tools/d22scan/**` 写成 `allowlist.txt|.*\.go` ⇒ 该目录里的 go.mod / runtests.sh 不在尺上；
#      "任何 golden" 写成 `.*/testdata/golden/.*` ⇒ `internal/llm/golden/` 那一支不在尺上。
#      本脚本两支都按票面措辞重画，并把宽窄两把尺都跑（见 GOLDEN-WIDE）。
#
# 尺面（写死，供"零命中"这句话负责）：
#   · 名册 = 单枚 commit 的 diff 文件路径（`git show`，非区间；区间尺会把别人的 commit 算进来，见输出末段 REV）
#   · 匹配 = `grep -E`，大小写敏感，路径前缀锚定（^ 或 (^|/)…$），**不**整词
#   · 分母 = `git ls-files` 在 HEAD 上命中该支的 tracked 件枚数（分母为 0 的尺＝死尺）
#   · 含不含 `_test.go`：一律含（本尺盘的是文件，不是符号）
#   · 每一发 grep 都跑两遍并记 rc（本仓 09-26 实发：单发会吞命中）
set -u
cd "$(git rev-parse --show-toplevel)" || exit 9

FACE_LABEL='docs/PLAN.md
docs/specs/**
internal/risk/**
internal/panel/**
internal/agent/**
internal/agent/approval/**
internal/observe/**
thresholds.go
任何 golden（窄：*/testdata/golden/*）
任何 golden（宽：路径里含 golden，大小写敏感）
allowlist.txt
scripts/slo-check.ps1
tools/d22scan/**
cmd/wisp/slo_windows.go'

FACE_RE='^docs/PLAN\.md$
^docs/specs/
^internal/risk/
^internal/panel/
^internal/agent/
^internal/agent/approval/
^internal/observe/
(^|/)thresholds\.go$
.*/testdata/golden/.*
golden
(^|/)allowlist\.txt$
^scripts/slo-check\.ps1$
^tools/d22scan/
^cmd/wisp/slo_windows\.go$'

# 前程脚本自带、票面与两枚派单都没点名的四支（加严，不是本票义务）——一并答，好让 13 支与 18 支对得上账
EXTRA_LABEL='〔加严〕internal/speech/**
〔加严〕internal/secret/**
〔加严〕.gitattributes
〔加严〕go.sum'
EXTRA_RE='^internal/speech/
^internal/secret/
^\.gitattributes$
^go\.sum$'

n_face=$(printf '%s\n' "$FACE_LABEL" | sed '/^$/d' | wc -l | tr -d ' ')
n_extra=$(printf '%s\n' "$EXTRA_LABEL" | sed '/^$/d' | wc -l | tr -d ' ')

cnt() { # $1 = pattern, $2 = text  → "命中枚数/grep rc"
  local n rc
  n=$(printf '%s\n' "$2" | grep -cE "$1"); rc=$?
  printf '%s/%s' "$n" "$rc"
}

ALLSHAS=("$@")
if [ ${#ALLSHAS[@]} -eq 0 ]; then echo "用法：$0 <sha> …"; exit 2; fi

echo "########## 0. 输入枚数与类型现读（引 sha 前先 git cat-file -t）"
for s in "${ALLSHAS[@]}"; do
  printf '  %-10s type=%-6s parents=%s  %s\n' "$s" "$(git cat-file -t "$s" 2>/dev/null)" \
    "$(git rev-list --no-walk --parents "$s" | awk 'NR==1{print NF-1}')" \
    "$(git log -1 --format='%ad %s' --date=format:'%m-%d %H:%M:%S' "$s" 2>/dev/null | cut -c1-46)"
done

echo
echo "########## 1. 名册现算（不抄任何前程表）"
for s in "${ALLSHAS[@]}"; do
  n1=$(git show --pretty=format: --name-only "$s" | sed '/^$/d')
  n2=$(git show --pretty=format: --name-status --no-renames "$s" | sed '/^$/d' | cut -f2-)
  c1=$(printf '%s\n' "$n1" | sed '/^$/d' | wc -l | tr -d ' ')
  c2=$(printf '%s\n' "$n2" | sed '/^$/d' | wc -l | tr -d ' ')
  same=YES; [ "$n1" != "$n2" ] && same=NO
  printf '  %-10s N1(name-only)=%2s  N2(no-renames)=%2s  N1≡N2: %s\n' "$s" "$c1" "$c2" "$same"
done
UNION=$(for s in "${ALLSHAS[@]}"; do
          git show --pretty=format: --name-only "$s" | sed '/^$/d'
          git show --pretty=format: --name-status --no-renames "$s" | sed '/^$/d' | cut -f2-
        done | sed '/^$/d' | sort -u)
printf '%s\n' "$UNION" > /tmp/ac4-union.$$.txt 2>/dev/null || true
U=$(printf '%s\n' "$UNION" | sed '/^$/d' | grep -c . )
E=$(for s in "${ALLSHAS[@]}"; do git show --pretty=format: --name-only "$s" | sed '/^$/d'; done | grep -c .)
printf '  并集去重 = %s 枚（逐枚 commit 条目合计 %s 枚）\n' "$U" "$E"
printf '  名册全文：\n'; printf '%s\n' "$UNION" | sed 's/^/    /'
echo '  ⚠ 对照：派单里那条命令照字面跑（git log --pretty=tformat: --name-only <四枚>，不带 --no-walk）'
printf '    会把四枚的全部祖先一起算进来 ⇒ 枚数 = %s\n' \
  "$(git log --pretty=tformat: --name-only "${ALLSHAS[@]}" | sed '/^$/d' | sort -u | grep -c . )"

echo
echo "########## 2. 逐支 × 逐枚（数字 = 命中枚数/grep rc；两发一致才写）"
printf '  %-42s denom  ' "支（票面 AC#4 逐字）"
for s in "${ALLSHAS[@]}"; do printf '%-9s' "${s:0:7}"; done
printf '%s\n' "正控"
i=1
while [ "$i" -le "$n_face" ]; do
  label=$(printf '%s\n' "$FACE_LABEL" | sed -n "${i}p")
  re=$(printf '%s\n' "$FACE_RE" | sed -n "${i}p")
  denom=$(git ls-files | grep -cE "$re"); drc=$?
  printf '  %-42s %-6s  ' "$label" "${denom}/$drc"
  for s in "${ALLSHAS[@]}"; do
    n1=$(git show --pretty=format: --name-only "$s" | sed '/^$/d')
    n2=$(git show --pretty=format: --name-status --no-renames "$s" | sed '/^$/d' | cut -f2-)
    a=$(printf '%s\n' "$n1" | grep -cE "$re"); ra=$?
    b=$(printf '%s\n' "$n2" | grep -cE "$re"); rb=$?
    tag="$a/$ra"
    [ "$a" != "$b" ] && tag="$a!=$b"
    printf '%-9s' "$tag"
  done
  probe=$(git ls-files | grep -E "$re" | head -1)
  if [ -n "$probe" ]; then
    pc=$(printf '%s\n' "$probe" | grep -cE "$re")
    [ "$pc" = 1 ] && printf '正控=1 (%s)\n' "$probe" || printf "正控=$pc <== 尺死\n"
  else
    printf '正控=n/a 分母 0 <== 尺死\n'
  fi
  i=$((i+1))
done

echo
echo "########## 2b. 前程脚本自带、票面未点名的四支（同一把尺）"
j=1
while [ "$j" -le "$n_extra" ]; do
  label=$(printf '%s\n' "$EXTRA_LABEL" | sed -n "${j}p")
  re=$(printf '%s\n' "$EXTRA_RE" | sed -n "${j}p")
  denom=$(git ls-files | grep -cE "$re")
  hits=""
  for s in "${ALLSHAS[@]}"; do
    n=$(git show --pretty=format: --name-only "$s" | sed '/^$/d' | grep -cE "$re"); rc=$?
    hits="$hits$(printf '%-9s' "$n/$rc")"
  done
  printf '  %-42s denom=%-6s %s\n' "$label" "$denom" "$hits"
  j=$((j+1))
done

echo
echo "########## 3. 第二口径：面级 tree/blob hash 逐枚相等（字节级，不看文件名）"
BLOBFACES='docs/PLAN.md
internal/observe/thresholds.go
tools/d22scan/allowlist.txt
scripts/slo-check.ps1
cmd/wisp/slo_windows.go
.gitattributes'
DIRFACES='docs/specs
internal/risk
internal/panel
internal/agent
internal/agent/approval
internal/observe
tools/d22scan
internal/llm/testdata/golden
internal/agent/testdata/golden
tools/mockllm/testdata/golden
internal/llm/golden'
FIRST=${ALLSHAS[0]}
BASE="${FIRST}^"
printf '  %-34s %s\n' "面" "$(printf '%-10s' "base"; for s in "${ALLSHAS[@]}"; do printf '%-10s' "${s:0:7}"; done) HEAD  全等?"
for f in $DIRFACES; do
  row=$(printf '  %-34s ' "$f")
  vals=""
  for ref in "$BASE" "${ALLSHAS[@]}" HEAD; do
    h=$(git rev-parse --quiet --verify "$ref:$f" 2>/dev/null)
    [ -z "$h" ] && h="(缺)"
    row="$row$(printf '%-10s' "${h:0:8}")"
    vals="$vals $h"
  done
  uniq=$(printf '%s\n' $vals | sed '/^$/d' | sort -u | grep -c .)
  printf '%s  全等=%s\n' "$row" "$([ "$uniq" = 1 ] && echo YES || echo "NO（$uniq 枚不同 ⇒ 看 REV 段归因）")"
done
for f in $BLOBFACES; do
  row=$(printf '  %-34s ' "$f")
  vals=""
  for ref in "$BASE" "${ALLSHAS[@]}" HEAD; do
    h=$(git rev-parse --quiet --verify "$ref:$f" 2>/dev/null)
    [ -z "$h" ] && h="(缺)"
    row="$row$(printf '%-10s' "${h:0:8}")"
    vals="$vals $h"
  done
  uniq=$(printf '%s\n' $vals | sed '/^$/d' | sort -u | grep -c .)
  printf '%s  全等=%s\n' "$row" "$([ "$uniq" = 1 ] && echo YES || echo "NO（$uniq 枚不同）")"
done

echo
echo "########## REV. 区间口径为什么不能用（本程现量，给下一位留证据而不是留规矩）"
printf '  区间 %s..%s 的 diff 条目 = %s 枚\n' "$BASE" "HEAD" "$(git diff --name-only "$BASE" HEAD | grep -c .)"
for f in docs/specs internal/risk internal/panel internal/agent internal/observe tools/d22scan scripts/slo-check.ps1 cmd/wisp/slo_windows.go docs/PLAN.md; do
  n=$(git diff --name-only "$BASE" HEAD | grep -cE "^$f"); printf '    %-30s 区间命中=%s\n' "$f" "$n"
done
printf '  其中 frontend/design 命中（**票面明写不算进任何零命中宣称**）：\n'
git diff --name-status "$BASE" HEAD | grep -E '(frontend|design)/' | sed 's/^/    /'
printf '  这些命中的归属（区间里别人的 commit）：\n'
git log --format='    %h %s' "$BASE"..HEAD -- frontend design | cut -c1-90
