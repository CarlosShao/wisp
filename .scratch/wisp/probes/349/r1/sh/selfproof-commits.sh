#!/bin/sh
# 收件自证①：逐枚提交号＋逐枚 git show --name-only ⇒ ⊆ 写面；顺带查明 HEAD..base 之间
# 那一枚不是我落的件（`.scratch/wisp/probes/orch/…`）是谁的笔＝⛔ 记在我账上。
set -u
cd '/d/work/workspace/projects plans/Wisp' || exit 90
A=7b633fa9713b0514806c8b808ed4db91b1485cdc
echo "== who added the orch send file (want: NOT one of my five) =="
git log --format='%h %s' "$A..HEAD" -- .scratch/wisp/probes/orch/20261019/349r1-send.txt | cut -c1-90
echo "== every commit in base..HEAD (author + subject head) =="
git log --format='%h | %an | %s' "$A..HEAD" | cut -c1-80
echo "== my five commits: file counts =="
for c in bf7f13bc 07f178df db2f9eab 7d56962b ed8db1ea; do
  n=$(git show --name-only --format='' "$c" | grep -c .)
  g=$(git show --name-only --format='' "$c" | grep -c '[.]go$')
  printf '%s files=%s go=%s\n' "$c" "$n" "$g"
done
echo "== out-of-scope scan over MY five commits only =="
for c in bf7f13bc 07f178df db2f9eab 7d56962b ed8db1ea; do
  bad=$(git show --name-only --format='' "$c" | grep . | grep -vE '^\.scratch/wisp/probes/349/r1/|^\.scratch/wisp/issues/349-|^internal/risk/(blacklist|provenance|syncdirs)[.]go$')
  if [ -n "$bad" ]; then printf '  %s OUT-OF-SCOPE: %s\n' "$c" "$bad"; else printf '  %s in-scope (all)\n' "$c"; fi
done
