#!/bin/bash
# 353-v1 R21 = 补齐 R20 因引号破损而未执行到的几段（量的）
cd '/d/work/workspace/projects plans/Wisp' || exit 9
L=/d/tmp/wisp/20261020/353r1/logs
V=/d/tmp/wisp/20261020/353v1
echo "###### 3-正控) 同一把尺打在 fa395b33 与 a59f5630 之间（应非 0＝乙本体确实改了可执行行）"
strip_exec() { git cat-file blob "$1:internal/winsec/$2" | grep -vE '^[[:space:]]*//' | grep -vE '^[[:space:]]*$'; }
n=$(diff -u0 <(strip_exec fa395b33 pinnedleaf_342_windows.go) <(strip_exec a59f5630 pinnedleaf_342_windows.go) | grep -Ec '^[+-][^+-]')
echo "  EXEC_DIFF_LINES_BASE_vs_2=$n"
echo "###### 4) deleted 非零（正控先行）"
echo "  正控串存在性 deleted= 总命中=$(grep -rho 'deleted=' "$L" | wc -l)"
echo "  deleted 非零命中=$(grep -rhoE 'deleted=[1-9][0-9]*' "$L" | wc -l)"
grep -rhoE 'deleted=[1-9][0-9]*' "$L" | sort | uniq -c | sed 's/^/   /'
echo "###### 5) 是否跑过 cmd/wisp（尺＝尺脚本本体＋底片包名）"
grep -rl 'cmd/wisp' /d/tmp/wisp/20261020/353r1/sh | sed 's/^/   尺脚本点名: /'
echo "  底片 cmd\/wisp 包名命中=$(grep -rho 'CarlosShao/wisp/cmd/wisp' "$L" | wc -l)  正控 winsec 包名命中=$(grep -rho 'CarlosShao/wisp/internal/winsec' "$L" | wc -l)"
echo "###### 6) 名册三把尺"
git log --all --format='%s' > $V/subj-v1.txt
git log --all --format='%H%x09%s' > $V/subjhash-v1.txt
echo "  ruler1 anchored=%s"
grep -c '^票 353 落地腿 353-r1' $V/subj-v1.txt | sed 's/^/   /'
echo "  ruler2 nonanchored"
grep -c '票 353 落地腿 353-r1' $V/subj-v1.txt | sed 's/^/   /'
echo "  ruler3 single-export index==1"
awk -F'\t' -v P='票 353 落地腿 353-r1' 'index($2,P)==1{c++} END{print c+0}' $V/subjhash-v1.txt | sed 's/^/   /'
echo "  ruler4 腿码任意位（别形尺）"
awk -F'\t' 'index($2,"353-r1")>0{c++} END{print c+0}' $V/subjhash-v1.txt | sed 's/^/   /'
grep '^票 353 落地腿 353-r1' $V/subj-v1.txt | cut -c1-46 | sed 's/^/   /'
echo "  正控 异锚：^票 353 普查腿 353-a1 / ^票 353 终裁腿 353-v1"
grep -c '^票 353 普查腿 353-a1' $V/subj-v1.txt | sed 's/^/   a1=/;s/$//;t;break'
grep -c '^票 353 普查腿 353-a1' $V/subj-v1.txt | sed 's/^/   a1=/ '
grep -c '^票 353 终裁腿 353-v1' $V/subj-v1.txt | sed 's/^/   v1=/ '
echo "###### 7) tracked 脏逐枚列名＋mtime"
git status --porcelain > $V/status-v1.txt
while read -r st p; do
  [ -n "$p" ] || continue
  f=${p# }
  printf '   [%s] %s mtime=%s\n' "${st:0:2}" "$f" "$(stat -c '%y' "$f" 2>/dev/null | cut -c1-19)"
done < $V/status-v1.txt
echo "###### 8) 11 枚族件在位尺（逐枚 cat-file -e）"
git ls-tree --name-only HEAD internal/winsec/ | sed 's|internal/winsec/||' | grep -E 'ancestor|leafrace|openleaf|posix_seal|sealhandle|world' | sed 's/^/   在册: /'
echo "  族件总数=$(git ls-tree --name-only HEAD internal/winsec/ | grep -cE 'ancestor|leafrace|openleaf|posix_seal|sealhandle')"
