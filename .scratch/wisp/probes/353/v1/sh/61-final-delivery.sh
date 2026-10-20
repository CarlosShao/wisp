#!/bin/sh
# 353-v1 R61 = 交件前终尺（名册三把⇄逐笔 name-only 的 .go⇄越界⇄主件 blob 件身⇄框终量⇄tracked 脏 mtime）
cd '/d/work/workspace/projects plans/Wisp' || exit 9
git log --all --format='%s' > /d/tmp/wisp/20261020/353v1/subj-end.txt
git log --all --format='%H%x09%s' > /d/tmp/wisp/20261020/353v1/subjhash-end.txt
echo "== 名册三把（射程逐字：尺1|2 打 %s 导出、尺3 打 %H%x09%s 导出＋内存 index==1）=="
printf '  ruler1 anchored-353v1=%s\n' "$(grep -c '^票 353 终裁腿 353-v1' /d/tmp/wisp/20261020/353v1/subj-end.txt || true)"
printf '  ruler2 index1-353v1=%s\n' "$(grep -F '票 353 终裁腿 353-v1' /d/tmp/wisp/20261020/353v1/subjhash-end.txt | awk -F'\t' '{if(index($2,"票 353 终裁腿 353-v1")==1)c++} END{print c+0}')"
printf '  ruler3 nonanchored=%s\n' "$(grep -c '票 353 终裁腿 353-v1' /d/tmp/wisp/20261020/353v1/subj-end.txt || true)"
printf '  ctrl r1-anchored=%s  ctrl a1-anchored=%s\n' "$(grep -c '^票 353 落地腿 353-r1' /d/tmp/wisp/20261020/353v1/subj-end.txt || true)" "$(grep -c '^票 353 普查腿 353-a1' /d/tmp/wisp/20261020/353v1/subj-end.txt || true)"
echo "== 本腿逐枚提交号＋名数＋.go＋越界（逐枚现跑）=="
for c in $(grep -n '^票 353 终裁腿 353-v1' /d/tmp/wisp/20261020/353v1/subj-end.txt >/dev/null; git log --all --format='%H' --grep='^票 353 终裁腿 353-v1'); do
  h=$(git rev-parse --short "$c")
  n=$(git show --name-only --format='' "$c" | grep -c .)
  g=$(git show --name-only --format='' "$c" | grep -c '\.go$' || true)
  o=$(git show --name-only --format='' "$c" | grep -vcE '^(\.scratch/wisp/probes/353/v1/|docs/evidence/s1/353-v1-verdict\.md)' || true)
  t=$(git log -1 --format='%ad' --date=format:'%H:%M' "$c")
  printf '  %s %s files=%s go=%s outofrange=%s\n' "$h" "$t" "$n" "$g" "$o"
done
echo "== 主件入库 blob 件身 =="
S7=$(printf '\xc2\xa7')
git cat-file blob 'HEAD:docs/evidence/s1/353-v1-verdict.md' > /d/tmp/wisp/20261020/353v1/_main.bin
printf '  BYTES=%s LINES=%s CR=%s SECT=%s LAST=%s\n' "$(wc -c < /d/tmp/wisp/20261020/353v1/_main.bin)" "$(wc -l < /d/tmp/wisp/20261020/353v1/_main.bin)" "$(tr -cd '\r' < /d/tmp/wisp/20261020/353v1/_main.bin | wc -c)" "$(grep -o "$S7" /d/tmp/wisp/20261020/353v1/_main.bin | wc -l)" "$(tail -c 1 /d/tmp/wisp/20261020/353v1/_main.bin | od -An -tx1 | tr -d ' \n')"
echo "== 票 353 框终量（入库 blob）＋Status 行在位 =="
T=.scratch/wisp/issues/353-pinned-name-swap-after-seal-deletes-an-object-the-caller-never-named.md
printf '  UNCHECKED=%s CHECKED=%s 票名-ddone枚数=%s\n' "$(git cat-file blob "HEAD:$T" | grep -cE '^- \[ \]' || true)" "$(git cat-file blob "HEAD:$T" | grep -cE '^- \[x\]' || true)" "$(git ls-tree --name-only HEAD .scratch/wisp/issues/ | grep -c 'issues/353-.*-done' || true)"
echo "== tracked 脏逐枚列名＋mtime =="
git status --porcelain -uno | while IFS= read -r line; do
  st=$(printf '%s' "$line" | cut -c1-2); f=$(printf '%s' "$line" | cut -c4-)
  printf '  [%s] %s mtime=%s\n' "$st" "$f" "$(stat -c '%y' "$f" 2>/dev/null | cut -c1-19)"
done
echo "== 本腿写面残留？（status 里 v1/主件的行数）=="
git status --porcelain -uall -- docs/evidence/s1/353-v1-verdict.md .scratch/wisp/probes/353/v1 | sed 's/^/  /'
echo "== 只读自证（★尺形＝排除本把尺自身：它⛔ 跑任何 go 命令、只是把那三枚串写在自己注释里＝不自我指涉）=="
HITS=$(grep -rl 'go test' .scratch/wisp/probes/353/v1/sh | grep -v '61-final-delivery.sh' | wc -l)
echo "  本腿其余 sh/ 里出现 'go test' 的枚数=$HITS（应 0＝本腿⛔ 写任何跑测试的尺；⛔ 那是被裁腿的尺、本腿只读它的底片）"
echo "  正控＝同尺打被裁腿的尺包 sh/：$(grep -rl 'go test' /d/tmp/wisp/20261020/353r1/sh 2>/dev/null | wc -l) 枚（>0＝尺灵＝那才是跑过 go test 的一族）"
