#!/usr/bin/env bash
# 票 157 验收程 · 取数器（全部只读；输出落本目录 *.txt，凭据＝真回显）
# 纪律：锚点一律 --no-walk 现读；数枚数一律锚定 grep；引行号一律按"被引者的锚点"取版。
R="D:/work/workspace/projects plans/Wisp"
cd "$R" || exit 9
OUT=".scratch/wisp/probes/157-accept"
A16="16-commit name list read at start-of-run"

run() { printf '$ %s\n' "$*"; eval "$@" 2>&1; printf 'rc=%s\n---\n' "$?"; }

# ---------- 00 preflight ----------
{
  echo "### P. 开工现量（时刻/分支/HEAD/名册尺）"
  run "date '+%Y-%m-%d %H:%M:%S %z'"
  run "go version"
  run "git rev-parse --abbrev-ref HEAD"
  echo "### P1. 开工时 16 枚 no-walk 名册（第一发，未转述）"
  git log --no-walk --pretty='%h %ad %s' --date=format:'%m-%d %H:%M' -16 | cut -c1-90
  echo "### P2. 同一条尺的枚数（锚定 grep 数 157 名下 commit）"
  run "git log --pretty=%h --grep='157 第' | wc -l"
  echo "### P3. 反面教材复算：不加 --no-walk 的同一条尺会数出多少枚"
  run "git log --pretty=%h -16 | wc -l"
  echo "### P4. 工作树脏行数（只登记不碰）"
  run "git status --porcelain | wc -l"
} > "$OUT/00-preflight.txt" 2>&1

# ---------- 10 items 4 & 5 by AST ----------
{
  echo "### S45. 第 4/5 笔换尺：go/ast 语法判据（不是 grep/sed 文本尺）"
  echo "取版尺＝git cat-file blob 5365cb22:<路径> 落 D:/tmp/157accept/{loopblob.txt}"
  run "git cat-file blob 5365cb22:internal/agent/loop.go | wc -l -c"
  run "git rev-parse 5365cb22:internal/agent/loop.go"
  echo "--- 分析器本体在 D:/tmp/157accept/astfacts2.go（零依赖，go/parser＋go/token）---"
  cd "D:/tmp/157accept" || { echo "CD-FAIL"; exit 8; }
  echo '$ go run astfacts2.go loopblob.txt current Loop run'
  go run astfacts2.go loopblob.txt current Loop run 2>&1
  printf 'rc=%s\n' "$?"
} > "$R/.scratch/wisp/probes/157-accept/10-ast-for-items-4-and-5.txt" 2>&1

# ---------- 11 item 1 by column split ----------
cd "$R" || exit 9
{
  echo "### S1. 第 1 笔换尺：git cat-file blob ＋ awk -F'|' 逐格切列（不是 grep -c）"
  echo "--- 155 证据件 §4 那张表在 155 自己那枚锚点 bb3a7a1 上的列形状 ---"
  git cat-file blob bb3a7a1:docs/evidence/s1/155-three-unjudged-cells-r1.md \
    | awk 'NR>=257 && NR<=268 {n=gsub(/\|/,"|"); if(n>0) printf "N%d cells=%d :: %.80s\n", NR, n+1, $0}'
  echo "--- §4 那十行里 (a)-(f) 与三枚档字的逐行抽取（awk 逐格，不用 grep -o 数总）---"
  git cat-file blob bb3a7a1:docs/evidence/s1/155-three-unjudged-cells-r1.md | sed -n '259,263p' \
    | while IFS= read -r l; do
        printf '%s\n' "  档字=[$(printf '%s' "$l" | grep -o -E '附条件|未裁|成立' | tr '\n' ',')] 字母=[$(printf '%s' "$l" | grep -o -E '\([a-f]\)' | tr '\n' ',')]"
      done
  echo "### S2. 第 1 笔的地基：153 验收件 §12 那八行的『最后一格』＝真档（锚点 b319bab＝该文件最后一次改动）"
  git cat-file blob b319bab:docs/evidence/s1/153-trace-lies-unguarded-r1-accept-r1.md \
    | awk 'NR>=831 && NR<=836 || NR>=842 && NR<=843 {n=split($0,p,"|"); g=""; for(i=n;i>=2;i--){if(length(p[i])>1){g=p[i];break}} printf "N%d 支=%.20s 真档(最后一格)=[%s]\n", NR, p[2], substr(g,1,60)}'
  echo "### S3. 同一批行拿它那把尺（index 第一个 '| **'）会读到什么"
  git cat-file blob b319bab:docs/evidence/s1/153-trace-lies-unguarded-r1-accept-r1.md \
    | awk 'NR>=831 && NR<=836 || NR>=842 && NR<=843 {i=index($0,"| **"); printf "N%d 它那把尺=[%s]\n", NR, (i>0?substr($0,i,34):"NONE")}'
  echo "### S4. 157 补进去那四行原文（锚点 ec81582）"
  git cat-file blob ec81582:docs/evidence/s1/155-three-unjudged-cells-r1.md | grep -n '^| (b) 问①\|^| (c) 问②\|^| AC#3 句②\|^| (a) 交付物本身' | cut -c1-120
} > "$OUT/11-item1-columns-and-grades.txt" 2>&1

# ---------- 12 item 8 by tree hash + real go list ----------
{
  echo "### S8a. 第 8 笔换尺①：目录树对象等值（一枚 sha 顶替逐文件 diff）"
  for r in b23c7f7 6de3d1c5 bb3a7a1 HEAD; do
    printf '  %-9s internal/agent tree=%s\n' "$r" "$(git rev-parse $r:internal/agent)"
  done
  echo "### S8b. 第 8 笔换尺②：逐枚 git show --name-only 并集（区间 b23c7f7..6de3d1c5）"
  echo "$ git rev-list b23c7f7..6de3d1c5 | wc -l"
  git rev-list b23c7f7..6de3d1c5 | wc -l
  for c in $(git rev-list b23c7f7..6de3d1c5); do git show --name-only --pretty= "$c"; done | sort -u \
    | sed 's/^/  /'
  echo "### S8c. 第 8 笔换尺③：真跑 go list 的依赖闭包（HEAD；先证 HEAD 与 6de3d1c5 的 internal/agent 同树）"
  echo "  非测试闭包（internal 部分）:"
  go list -deps ./internal/agent 2>/dev/null | grep '^github.com/CarlosShao/wisp/internal' | sed 's/^/    /'
  echo "  test 闭包（internal 部分，去重）:"
  go list -test -deps ./internal/agent 2>/dev/null | grep '^github.com/CarlosShao/wisp/internal' | sort -u | sed 's/^/    /'
  echo "  internal/tools 在非测试闭包？ $(go list -deps ./internal/agent 2>/dev/null | grep -c 'wisp/internal/tools$')  在 test 闭包？ $(go list -test -deps ./internal/agent 2>/dev/null | grep -c 'wisp/internal/tools$')"
  echo "  internal/agent/approval 在 test 闭包？ $(go list -test -deps ./internal/agent 2>/dev/null | grep -c 'wisp/internal/agent/approval$')"
  echo "### S8d. 6de3d1c5 与 HEAD 逐包树等值（闭包成员有没有在别处漂）"
  for p in $(git ls-tree --name-only 6de3d1c5:internal); do
    a=$(git rev-parse 6de3d1c5:internal/$p); b=$(git rev-parse HEAD:internal/$p)
    [ "$a" = "$b" ] && printf '  %-14s SAME\n' "$p" || printf '  %-14s DIFF\n' "$p"
  done
} > "$OUT/12-item8-treehash-golist-and-union.txt" 2>&1

# ---------- 13 task-2 contradiction search ----------
{
  echo "### C. 那件'自相矛盾'的原文到底在不在盘上"
  echo '### C1. 全仓（docs ＋ issues ＋ dispatches）逐字搜那两句'
  run "grep -rn '锚点反查' docs .scratch/wisp/issues .scratch/wisp/dispatches | wc -l"
  grep -rn '锚点反查' docs .scratch/wisp/issues .scratch/wisp/dispatches 2>/dev/null | sed 's/^/  /' | cut -c1-190
  echo "### C2. 155 证据件在它自己锚点 bb3a7a1 上的 §0 标题与反查表（取版＝git cat-file blob）"
  git cat-file blob bb3a7a1:docs/evidence/s1/155-three-unjudged-cells-r1.md | sed -n '9,17p' | cut -c1-150 | sed 's/^/  /'
  echo "### C3. 实现程那枚反查凭据是不是它自己的、并是否入库"
  run "git cat-file -t bb3a7a1:.scratch/wisp/probes/155/00-anchor-reverse-check.txt"
  git cat-file blob bb3a7a1:.scratch/wisp/probes/155/00-anchor-reverse-check.txt | sed -n '1,20p' | sed 's/^/  /' | cut -c1-150
  echo "### C4. 本程自己 cat-file -t：155 §0 表里那六枚 ＋ 验收件那两枚"
  for a in ff550f3 5365cb22 86b0161 b23c7f7 6de3d1c5 777d6cc ac7fb00 bb3a7a1 b319bab 3274e4d; do
    printf '  %-11s %s\n' "$a" "$(git cat-file -t $a 2>&1)"
  done
  echo "### C5. 155 验收件在它自己锚点 3274e4d 上的 §7.1 那一行（'AC#1 成立'那句的真出处）"
  git cat-file blob 3274e4d:docs/evidence/s1/155-three-unjudged-cells-r1-accept-r1.md | sed -n '254,260p' | cut -c1-150 | sed 's/^/  /'
  echo "### C6. 两枚 155 文件里所有含'反查'的行（锚点取版，逐枚列）"
  git cat-file blob bb3a7a1:docs/evidence/s1/155-three-unjudged-cells-r1.md | grep -n '反查' | cut -c1-150 | sed 's/^/  证据件 /'
  git cat-file blob 3274e4d:docs/evidence/s1/155-three-unjudged-cells-r1-accept-r1.md | grep -n '反查' | cut -c1-150 | sed 's/^/  验收件 /'
  echo "### C7. 两句引文在仓里的全部命中面（除台账与派单之外还有没有）"
  grep -rn '没做锚点反查\|三枚锚点号先反查' docs .scratch 2>/dev/null | cut -c1-150 | sed 's/^/  /'
} > "$OUT/13-contradiction-original-text.txt" 2>&1

# ---------- 14 AC rosters and shape ----------
{
  echo "### V. 票面 5 枚框 ↔ 盘上形状（本程自己那把尺）"
  echo "### V1. 157 名下 commit 名册（锚定 grep ＋ no-walk 并集）"
  run "git log --pretty=%h --grep='157 第' | wc -l"
  for c in $(git log --pretty=%h --grep='157 第'); do git show --name-only --pretty= "$c"; done | sort -u | sed 's/^/  /'
  echo "  并集枚数：$(( $(for c in $(git log --pretty=%h --grep='157 第'); do git show --name-only --pretty= $c; done | sort -u | wc -l) ))"
  echo "### V2. 禁改面命中数（同一份名册上按前缀尺）"
  for c in $(git log --pretty=%h --grep='157 第'); do git show --name-only --pretty= "$c"; done | sort -u \
    | grep -E '^(internal/|cmd/|docs/PLAN|docs/specs/|docs/reports/|frontend/|design/|tools/d22scan|scripts/|third_party/)' | wc -l | sed 's/^/  前缀尺命中＝/'
  for c in $(git log --pretty=%h --grep='157 第'); do git show --name-only --pretty= "$c"; done | sort -u \
    | grep -E 'thresholds\.go|golden|allowlist\.txt|slo-check' | wc -l | sed 's/^/  阈值/golden/allowlist 命中＝/'
  echo "### V3. 票面框的勾与 -done"
  grep -c '^\- \[x\]' .scratch/wisp/issues/157-*.md | sed 's/^/  已勾框枚数＝/'
  grep -c '^\- \[ \]' .scratch/wisp/issues/157-*.md | sed 's/^/  未勾框枚数＝/'
  ls .scratch/wisp/issues | grep -c '157.*-done' | sed 's/^/  -done 枚数＝/'
  echo "### V4. 两节附录的标题与枚数（锚定 grep '^## 附录'）"
  git cat-file blob ec81582:docs/evidence/s1/155-three-unjudged-cells-r1.md | grep -n '^## 附录\|^### B\.' | cut -c1-70 | sed 's/^/  155件 /'
  git cat-file blob ec81582:docs/evidence/s1/153-trace-lies-unguarded-r1.md | grep -n '^## 附录\|^### B\.' | cut -c1-70 | sed 's/^/  153件 /'
  echo "### V5. 件数与'枚'的口径（自报 40 枚那一笔）"
  echo "  盘上文件枚数＝$(find .scratch/wisp/probes/157 -type f | wc -l)"
  echo "  入库枚数（ec81582）＝$(git ls-tree -r --name-only ec81582 .scratch/wisp/probes/157 | wc -l)"
  echo "  读数件里 '\$ ' 命令行总数＝$(cat .scratch/wisp/probes/157/*.txt | grep -c '^[$] ')"
  echo "  读数件里 '### ' 小节总数＝$(cat .scratch/wisp/probes/157/*.txt | grep -c '^###')"
  echo "  两枚交付件里出现 '40 枚' 的枚数＝$(git cat-file blob ec81582:docs/evidence/s1/155-three-unjudged-cells-r1.md | grep -c '40 枚'; git cat-file blob ec81582:docs/evidence/s1/153-trace-lies-unguarded-r1.md | grep -c '40 枚' | tr '\n' ' ')"
  echo "### V6. 交付件里有没有 AC#5 那枚双向对账块（锚定 grep）"
  echo "  155 件附录里 '对账|双向' 命中＝$(git cat-file blob ec81582:docs/evidence/s1/155-three-unjudged-cells-r1.md | sed -n '389,760p' | grep -c '对账\|双向')"
  echo "  153 件附录里 '对账|双向' 命中＝$(git cat-file blob ec81582:docs/evidence/s1/153-trace-lies-unguarded-r1.md | sed -n '593,700p' | grep -c '对账\|双向')"
  echo "### V7. 每笔带没带当轮命令（'\$ ' 行枚数，按 B.n 小节切）"
  git cat-file blob ec81582:docs/evidence/s1/155-three-unjudged-cells-r1.md \
    | awk '/^### B\./{sec=$0; c[sec]=0; n++; order[n]=sec} sec!=""{if(/^\$ /) c[sec]++} END{for(i=1;i<=n;i++) printf "  %-46s $-命令行＝%d\n", order[i], c[order[i]]}'
  git cat-file blob ec81582:docs/evidence/s1/153-trace-lies-unguarded-r1.md \
    | awk '/^### B\./{sec=$0; c[sec]=0; n++; order[n]=sec} sec!=""{if(/^\$ /) c[sec]++} END{for(i=1;i<=n;i++) printf "  %-46s $-命令行＝%d\n", order[i], c[order[i]]}'
} > "$OUT/14-ac-rosters-and-shape.txt" 2>&1

echo done
