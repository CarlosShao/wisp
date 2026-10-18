#!/usr/bin/env bash
# 328-v1 起手尺 v2（非实现者终裁腿，第二程入场）：逐条现跑，⛔ 抄 c59d0354 的读数、⛔ 抄派单里的号。
# 输出 tee 进 logs/start-rulers-v2.txt＝入库件。
set -u
REPO='D:/work/workspace/projects plans/Wisp'
cd "$REPO" || exit 1
TICKET='.scratch/wisp/issues/328-l2-card-corr-shape-never-matches-the-roster-join-key.md'

echo "date=$(date '+%F %T %z') leg=328-v1 (非实现者终裁腿·第二程)"
echo
echo '## 0. ★落点改判（机主 2026-10-18 10:1x 亲口新铁律＝盘符根目录零污染，覆盖派单台面配方两行；lead 10:1x 转达，我已回执）'
echo '本腿新建仓外导出树只落 D:\tmp\wisp\20261018\328v1\{before,after}；GOTMPDIR 同根下 gotmp。'
echo '⛔ 复用/⛔ 删除：wisp-328r1-before-tree、wisp-328r1-after-tree（上一程落地腿旧形）、w328v1a、w328v1b（第一程我这条腿建在别处的两棵）——一律只读不动。'
for t in 'D:/work/workspace/wisp-328r1-before-tree' 'D:/work/workspace/wisp-328r1-after-tree' 'D:/work/workspace/w328v1a' 'D:/work/workspace/w328v1b'; do
  echo "  untouched_old_tree $t exists=$([ -d "$t" ] && echo yes || echo NO)"
done
echo "  new_base=$( [ -d 'D:/tmp/wisp/20261018' ] && echo 'D:\tmp\wisp\20261018 存在（机主既有工作根，非我新建）' || echo MISSING )"
echo

echo '## 1. 起手锚（现跑）'
git log -1 --format='%H %ad %s' --date=iso
echo "branch=$(git rev-parse --abbrev-ref HEAD)"
echo '★派单里那枚 d175e483 与第一程的 78d3b46a 都已被后续盖过 ⇒ 我的锚＝上面这枚（编排者补的那句 c11163c8 现跑对上）。'
echo "c59d0354_is_ancestor=$(git merge-base --is-ancestor c59d0354 HEAD && echo YES || echo NO)"
echo "46d815d9_is_ancestor=$(git merge-base --is-ancestor 46d815d9 HEAD && echo YES || echo NO)"
echo "f705b082_is_ancestor=$(git merge-base --is-ancestor f705b082 HEAD && echo YES || echo NO)"
echo "648eda04_is_ancestor=$(git merge-base --is-ancestor 648eda04 HEAD && echo YES || echo NO)"
echo

echo '## 2. 框数尺（尺＝派单指定的那条 grep，射程＝票 328 一枚文件，量的）'
grep -n '^[[:space:]]*- \[' "$TICKET" | sed 's/\(.\{64\}\).*/\1…/'
echo "boxes_total=$(grep -c '^[[:space:]]*- \[' "$TICKET")"
echo "boxes_unchecked=$(grep -c '^[[:space:]]*- \[ \]' "$TICKET")"
echo "boxes_checked=$(grep -c '^[[:space:]]*- \[x\]' "$TICKET")"
echo '⇒ 盘上现量＝五枚框（AC#0…AC#4），未勾 3＝AC#2/AC#3/AC#4，已勾 2＝AC#0/AC#1（裁形那轮编排者翻的）。派单那句"五枚"成立；腿那句"七格"作废。'
echo

echo '## 3. 共享编译面 status 尺（本波写面互斥的入场复核）'
echo '$ git status --porcelain -- internal/panel cmd/wisp'
out=$(git status --porcelain -- internal/panel cmd/wisp); rc=$?
echo "  输出枚数=$(printf '%s' "$out" | grep -c . || true) rc=$rc"
[ -n "$out" ] && echo "$out"
echo '⇒ 空＝internal/panel 与 cmd/wisp 起手零在飞字节（编排者派单那句"零写腿"在我这侧对盘上成立）。'
echo

echo '## 4. 代码代表性尺（我要裁的是盘上今天这份字节）'
echo '$ git diff --name-only f705b082..HEAD -- internal/panel internal/agent'
git diff --name-only f705b082..HEAD -- internal/panel internal/agent
echo "diff_lines=$(git diff f705b082..HEAD -- internal cmd | wc -l)"
echo '⇒ 0 行：落地腿交件侧 f705b082 的那两包字节 == 我的锚（编排者 10-18 10:0x 现跑的 78d3b46a..HEAD=0 我入场重验同读数）。after 侧树以我的锚取即代腿的字节。'
echo '$ git diff --name-only 648eda04..HEAD -- internal（改前⇄改后的归属尺）'
git diff --name-only 648eda04..HEAD -- internal
echo '⇒ 那 648eda04..HEAD 之间 internal 面只有腿自己四枚文件 ⇒ 成对两把临时根台面的差异＝本票那一改，无别腿字节混进来。'
echo

echo '## 5. 落地腿提交尺（严格前缀，派单说 11 笔）'
echo "$ git log --all --grep='^328-r1 ' --oneline | wc -l"
git log --all --grep='^328-r1 ' --oneline | wc -l
git log --all --grep='^328-r1 ' --format='  %h %s' | sed 's/\(.\{90\}\).*/\1…/'
echo

echo '## 6. 台面件可用性正控（先证存在再信任何 0 读数）'
echo "design_dir=$([ -d design ] && echo present || echo MISSING)"
echo "gofumpt_real=$(ls -l /d/work/base/gopath/bin/gofumpt.exe 2>/dev/null | awk '{print $5" 字节"}')"
echo "docs_plan_shared=$([ -f docs/PLAN.md ] && echo present || echo MISSING)"
echo "scripts_d22scan=$([ -f scripts/d22scan.sh ] && echo present || echo MISSING) scripts_lenbudget=$([ -f scripts/check-path-length-budget.sh ] && echo present || echo MISSING)"
echo '⇒ 门禁四数那把 gofumpt 真身与 d22scan 射程的 design/ 都在；凡红句含 walk/build failed/no test files/setup failed/0xc0000135 ⇒ 判台面故障，不计入判据。'
echo

echo '## 7. ★具名上报（复认第一程 c59d0354 那处，我自己现跑，⛔ 抄它）'
echo '$ git show 46d815d9:.scratch/wisp/probes/328/r1/logs/mutate.sh | grep -n -A2 "M1)"'
git show 46d815d9:.scratch/wisp/probes/328/r1/logs/mutate.sh | grep -n -A3 'M1)'
echo '$ 派单 AC#2③ 要我补的那一发（逐字）：strings.HasPrefix(key, task+"#") → strings.HasPrefix(key, task)'
echo '⇒ 盘上腿 M1 的替换串与派单点名要我补的那一发逐字节同形（同一枚函数体内 token、同一枚去处）。'
echo '   ⇒ 我仍自己现跑这一发（完成通知一律当未验证；腿的日志不是我的读数），判语按我的读数写。'
echo '   ⇒ 对编排者的具名发现：派单把这一发描述成"腿没发的反形"与盘上不符；⛔ 据此退腿那一发（腿件①第三节发一确实打了 :142 那一枚）。'
echo '   ⇒ 派单真正的判据我照打：第三臂单独敏感性（去边界＝纯前缀）必须让 :142 child-1x 红，这一发的读数只有我自己跑出来才算。'
