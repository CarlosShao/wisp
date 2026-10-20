#!/bin/bash
# 353-v1 R21b = 补齐段（逐段落文件、⛔ 可疑 sed）
cd '/d/work/workspace/projects plans/Wisp' || exit 9
L=/d/tmp/wisp/20261020/353r1/logs
V=/d/tmp/wisp/20261020/353v1
OUT=$V/logs-R21.txt
mkdir -p $V
: > "$OUT"
strip_exec() { git cat-file blob "$1:internal/winsec/$2" | grep -vE '^[[:space:]]*//' | grep -vE '^[[:space:]]*$'; }
{
echo "###### 3-CTRL exec diff lines BASE vs fa395b33 (expect >0)"
diff -u0 <(strip_exec fa395b33 pinnedleaf_342_windows.go) <(strip_exec a59f5630 pinnedleaf_342_windows.go) | grep -Ec '^[+-][^+-]'
echo "###### 4 deleted nonzero"
grep -rho 'deleted=' "$L" | wc -l
grep -rhoE 'deleted=[1-9][0-9]*' "$L" | wc -l
echo "###### 5 cmd/wisp in leg rulers"
grep -rl 'cmd/wisp' /d/tmp/wisp/20261020/353r1/sh | wc -l
echo "###### 5b negatives pkg-name counts"
grep -rho 'CarlosShao/wisp/cmd/wisp' "$L" | wc -l
grep -rho 'CarlosShao/wisp/internal/winsec' "$L" | wc -l
echo "###### 6 roster rulers"
git log --all --format='%s' > $V/subj-v1.txt
git log --all --format='%H%x09%s' > $V/subjhash-v1.txt
echo -n "anchored="; grep -c '^票 353 落地腿 353-r1' $V/subj-v1.txt
echo -n "nonanchored="; grep -c '票 353 落地腿 353-r1' $V/subj-v1.txt
echo -n "index1="; grep -F '票 353 落地腿 353-r1' $V/subjhash-v1.txt | awk -F'\t' '{if(index($2,"票 353 落地腿 353-r1")==1)c++} END{print c+0}'
echo -n "legcode-anywhere="; grep -c '353-r1' $V/subj-v1.txt
echo -n "ctrl_a1="; grep -c '^票 353 普查腿 353-a1' $V/subj-v1.txt
echo -n "ctrl_v1="; grep -c '^票 353 终裁腿 353-v1' $V/subj-v1.txt
echo "###### 7 tracked dirty"
git status --porcelain
echo "###### 8 family files on register"
git ls-tree --name-only HEAD internal/winsec/ | sed 's|internal/winsec/||' | grep -E 'ancestor|leafrace|openleaf|posix_seal|sealhandle'
} >> "$OUT" 2>&1
echo RC=$?
cat "$OUT"
