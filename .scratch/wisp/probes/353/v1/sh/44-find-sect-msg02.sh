#!/bin/sh
# 353-v1 R44 = 定位那枚字面分节码位在哪一行哪一处（量的）
S7=$(printf '\xc2\xa7')
F=/d/tmp/wisp/20261020/353v1/msg-02.txt
grep -n "$S7" "$F" | while IFS= read -r l; do
  echo "LINE_HEAD: ${l%%:*}"
  echo "$l" | grep -o ".\{0,40\}$S7.\{0,40\}"
done
echo "== 同尺打在判语主件⇄读数⇄草案（应 0）=="
cd '/d/work/workspace/projects plans/Wisp' || exit 9
for f in docs/evidence/s1/353-v1-verdict.md .scratch/wisp/probes/353/v1/01-independent-readings.txt .scratch/wisp/probes/353/v1/02-ticket-draft-secC.md .scratch/wisp/probes/353/v1/00-ruler-register.md; do
  printf '  %-52s SECT=%s\n' "$(basename "$f")" "$(grep -o "$S7" "$f" | wc -l)"
done
echo "== 已入库笔②的 commit message（%B）里那枚（尺＝git log --format=%B -1 <笔②>）=="
git log --format=%B -1 4fb6c088 > /d/tmp/wisp/20261020/353v1/_msg02-from-repo.txt
printf '  REPO_MSG02 SECT=%s CR=%s\n' "$(grep -o "$S7" /d/tmp/wisp/20261020/353v1/_msg02-from-repo.txt | wc -l)" "$(tr -cd '\r' < /d/tmp/wisp/20261020/353v1/_msg02-from-repo.txt | wc -c)"
echo "== 笔①的 message 同尺 =="
git log --format=%B -1 3ac9d761 > /d/tmp/wisp/20261020/353v1/_msg01-from-repo.txt
printf '  REPO_MSG01 SECT=%s\n' "$(grep -o "$S7" /d/tmp/wisp/20261020/353v1/_msg01-from-repo.txt | wc -l)"
echo "== HEAD ⇄ 353-r1 有没有新笔（锚定尺）=="
git log -1 --format='%h %ad %s' --date=format:'%H:%M' | cut -c1-100
git log --all --format='%s' > /d/tmp/wisp/20261020/353v1/subj-now.txt
printf '  anchored 353-r1=%s  anchored 353-v1=%s\n' "$(grep -c '^票 353 落地腿 353-r1' /d/tmp/wisp/20261020/353v1/subj-now.txt || true)" "$(grep -c '^票 353 终裁腿 353-v1' /d/tmp/wisp/20261020/353v1/subj-now.txt || true)"
grep '^票 353 落地腿 353-r1' /d/tmp/wisp/20261020/353v1/subj-now.txt | cut -c1-40
