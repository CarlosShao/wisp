#!/bin/sh
# 票 155 对抗验收程 · AC#3 契约轴：名册现算 + 逐支禁改面扫 + 正控
# 规矩（派单）：名册不许抄它表里的数；零命中一律重跑并记 rc。
R="D:/work/workspace/projects plans/Wisp"
cd "$R" || exit 9
RANGE="829abe51^..bb3a7a1e"

echo '### 1) 本程一枚尺的坑：git log A B C（不带 --not/区间）会把整条历史吞进来'
cm=$(git log --format=%h $RANGE -- docs/evidence/s1/155-three-unjudged-cells-r1.md | tr '\n' ' ')
echo "-- 只按证据件路径筛出的 commit（$cm）再 naive 展开，名册枚数："
git log --pretty=tformat: --name-only $cm | grep -v '^$' | sort -u | wc -l
echo "   ⇒ 上面那枚数明显大于一切合理名册＝区间没界住（本程记成一枚仪器坑，不用它当读数）"

echo
echo "### 2) 逐枚 commit 判定它属不属于 155（尺＝只碰这两枚写面路径之一）"
: > /tmp/one55.txt
for c in $(git log --format=%h $RANGE | sort -r); do
  files=$(git show --pretty=tformat: --name-only "$c" | grep -v '^$')
  if printf '%s\n' "$files" | grep -qE '^(docs/evidence/s1/155-|\.scratch/wisp/probes/155/)'; then
    echo "155   $(git log -1 --format='%h %ad %s' --date=format:'%H:%M:%S' $c | cut -c1-60)"
    printf '%s\n' "$c" >> /tmp/one55.txt
  else
    echo "非155 $(git log -1 --format='%h %ad %s' --date=format:'%H:%M:%S' $c | cut -c1-60)"
  fi
done
echo "-- 属于 155 的枚数 --"; wc -l < /tmp/one55.txt

echo
echo "### 3) 名册并集（逐枚 git show 相加再 sort -u，不用 git log 的多 rev 形式）"
: > /tmp/roster55.txt
while read -r c; do
  git show --pretty=tformat: --name-only "$c" | grep -v '^$' >> /tmp/roster55.txt
done < /tmp/one55.txt
sort -u /tmp/roster55.txt > /tmp/roster55-uniq.txt
echo "并集枚数=$(wc -l < /tmp/roster55-uniq.txt)"
cat /tmp/roster55-uniq.txt
echo "-- 名册里有没有落在两枚写面之外的路径（应只有它自己的两枚） --"
grep -vE '^(docs/evidence/s1/155-three-unjudged-cells-r1\.md|\.scratch/wisp/probes/155/)' /tmp/roster55-uniq.txt
echo "rc=$?（1＝没有越界路径）"

echo
echo "### 4) 逐支禁改面扫（尺＝名册 grep -E 行首锚定；每支两跑记 rc）"
for face in '^docs/PLAN\.md' '^docs/specs/' '^internal/risk/' '^internal/panel/' \
            '^internal/agent/approval/' '^internal/observe/' 'thresholds\.go' 'golden' \
            'allowlist\.txt' '^scripts/slo-check\.ps1' '^tools/d22scan/' '^frontend/' '^design/' \
            '^internal/' '^cmd/' '^third_party/'; do
  n=$(grep -cE "$face" /tmp/roster55-uniq.txt); rc1=$?
  n2=$(grep -cE "$face" /tmp/roster55-uniq.txt); rc2=$?
  printf "%-32s 命中=%s (rc两跑=%s/%s)\n" "$face" "$n" "$rc1" "$rc2"
done
echo "注：'golden' 那支在本程名册里也 0——它没碰任何 fixture；正控在下一节。"

echo
echo "### 5) 正控：同一批尺打在 b23c7f7（153 的生产改动）上，必须命中"
git show --pretty=tformat: --name-only b23c7f7 > /tmp/pc55.txt
for face in '^internal/' '^internal/agent/' '^docs/PLAN\.md' '^frontend/'; do
  printf "%-24s 命中=%s\n" "$face" "$(grep -cE "$face" /tmp/pc55.txt)"
done
echo "-- b23c7f7 全名册 --"; cat /tmp/pc55.txt

echo
echo "### 6) 票面与 -done：本程与被验程都没动（现读）"
for f in .scratch/wisp/issues/153-*.md .scratch/wisp/issues/155-*.md; do
  printf "%-140s 未勾框=%s 已勾框=%s\n" "$f" "$(grep -c '^- \[ \]' "$f")" "$(grep -c '^- \[x\]' "$f")"
done
echo "-done 枚数(153/155)=$(ls .scratch/wisp/issues/ | grep -cE '^(153|155)-.*-done')"
echo "-- 153 票面最后一次改动 / 155 票面最后一次改动 --"
git log -1 --format="%h %ad %s" --date=format:"%m-%d %H:%M" -- .scratch/wisp/issues/153-*.md | cut -c1-80
git log -1 --format="%h %ad %s" --date=format:"%m-%d %H:%M" -- .scratch/wisp/issues/155-*.md | cut -c1-80
