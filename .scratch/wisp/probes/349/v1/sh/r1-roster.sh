#!/bin/sh
# 名册尺：单次 %H%x09%s 导出 + 内存 index($2,P)==1，两把锚并集去重
cd '/d/work/workspace/projects plans/Wisp' || exit 9
git log --all --format='%H%x09%s' > /d/tmp/wisp/20261019/349v1/logs/roster-export.txt 2>&1
echo EXPORT_LINES=$(wc -l < /d/tmp/wisp/20261019/349v1/logs/roster-export.txt)
A1='票 349 落地腿 349-r1'
A2='349-r1'
echo '--- ANCHOR1 (index==1) ---'
awk -F'\t' -v P="$A1" 'index($2,P)==1 {print $1}' /d/tmp/wisp/20261019/349v1/logs/roster-export.txt | tee /d/tmp/wisp/20261019/349v1/logs/roster-a1.txt | wc -l
echo '--- ANCHOR2 (index==1) ---'
awk -F'\t' -v P="$A2" 'index($2,P)==1 {print $1}' /d/tmp/wisp/20261019/349v1/logs/roster-export.txt | tee /d/tmp/wisp/20261019/349v1/logs/roster-a2.txt | wc -l
echo '--- UNION DEDUP ---'
cat /d/tmp/wisp/20261019/349v1/logs/roster-a1.txt /d/tmp/wisp/20261019/349v1/logs/roster-a2.txt | sort -u > /d/tmp/wisp/20261019/349v1/logs/roster-union.txt
echo UNION=$(wc -l < /d/tmp/wisp/20261019/349v1/logs/roster-union.txt)
echo '--- UNION SUBJECTS ---'
while read h; do git log -1 --format='%h %ad %s' --date=format:'%H:%M' "$h"; done < /d/tmp/wisp/20261019/349v1/logs/roster-union.txt
echo '--- ANCESTRY CHECK vs HEAD ---'
while read h; do git merge-base --is-ancestor "$h" HEAD && echo "ANCYES $h" || echo "ANCNO $h"; done < /d/tmp/wisp/20261019/349v1/logs/roster-union.txt
echo '--- POSITIVE CONTROL: 349-v1 anchor should be 0 ---'
awk -F'\t' 'index($2,"票 349 终裁腿 349-v1")==1 {print $1}' /d/tmp/wisp/20261019/349v1/logs/roster-export.txt | wc -l
echo '--- NEGATIVE SHAPE CONTROL: bare --grep non-anchored count ---'
git log --all --grep='349-r1' --format='%H' | wc -l
echo RC=$?
