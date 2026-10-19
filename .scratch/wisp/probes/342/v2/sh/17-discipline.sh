#!/bin/sh
# 342-v2 纪律格尺：① 17 笔逐枚 name-only ⊆ 声明写面；② `:96`/`:106` 旗标与"名字已存在⛔ 构成拒绝理由"
# 两段⛔ 动一字；③ 两枚删除的错误照旧被丢弃；④ 票名/Status 未动。
cd '/d/work/workspace/projects plans/Wisp' || exit 9
OUT=/d/tmp/wisp/20261019/342v2
LOG=$OUT/logs/x01-discipline.txt
{
echo '### A. 逐枚 name-only ⊆ 声明写面（尺＝名册 17 笔 × git show --name-only）'
grep -v '^COMMIT ' "$OUT/per-commit-nameonly.txt" | sort -u > "$OUT/leg-unique-files-sorted.txt"
TOTAL=$(wc -l < "$OUT/leg-unique-files-sorted.txt")
BAD=$(grep -vc -E '^(internal/winsec/|\.scratch/wisp/probes/34(2|8)/r1/|\.scratch/wisp/issues/34(2|8)-)' "$OUT/leg-unique-files-sorted.txt" || true)
echo "UNIQUE_FILES=$TOTAL OUT_OF_DECLARED_WRITEFACE=$BAD"
grep -v -E '^(internal/winsec/|\.scratch/wisp/probes/34(2|8)/r1/|\.scratch/wisp/issues/34(2|8)-)' "$OUT/leg-unique-files-sorted.txt" || true
echo
echo '### B. 覆盖写旗标⇄两段承诺文字在两端逐字同（尺＝git cat-file blob + grep -F 精确串）'
for A in fde1df02 9581a48c; do
  git cat-file blob "$A:internal/winsec/winsec.go" > "$OUT/winsec-$A.go"
  printf '%s flag96=%s flag106=%s\n' "$A" \
    "$(grep -c 'os.O_WRONLY|os.O_CREATE|os.O_TRUNC' "$OUT/winsec-$A.go")" \
    "$(grep -c 'os.O_WRONLY|os.O_CREATE|os.O_EXCL' "$OUT/winsec-$A.go")"
  printf '%s promise=%s\n' "$A" "$(grep -c 'Turning a pre-existing name' "$OUT/winsec-$A.go")"
done
echo '--- 那两行的逐字对照（sha16 of the exact lines）---'
for A in fde1df02 9581a48c; do
  grep -E 'return privateFile\(path, data, os.O_WRONLY' "$OUT/winsec-$A.go" | sha256sum | cut -c1-16 | sed "s/^/$A privateFile_flags_line_sha16=/"
  sed -n '87,94p' "$OUT/winsec-$A.go" | sha256sum | cut -c1-16 | sed "s/^/$A promise_block_sha16=/"
  sed -n '118,140p' "$OUT/winsec-$A.go" | sha256sum | cut -c1-16 | sed "s/^/$A selfstatement_118_140_sha16=/"
done
echo
echo '### C. 两处删除的返回错误是否仍被丢弃（落地形）'
grep -n 'unlinkPinnedLeaf' internal/winsec/winsec.go
grep -n '_ = unix.Unlinkat' internal/winsec/pinnedleaf_342_unix.go
grep -n '_ = p.root.Remove(base)' internal/winsec/pinnedleaf_342_windows.go
echo
echo '### D. 票名⇄Status⇄框（起手＝终量应逐字同）'
git ls-tree HEAD --name-only -- .scratch/wisp/issues/ | grep -E 'issues/(342|348)-'
head -4 .scratch/wisp/issues/342-privatefile-open-follows-a-swapped-ancestor-on-both-platforms-no-ticket-owns-it.md | tail -1 | cut -c1-40
head -4 .scratch/wisp/issues/348-refusal-path-remove-follows-an-ancestor-link-and-deletes-a-foreign-file.md | tail -1 | cut -c1-40
echo "T342_DONE=$(grep -c '^- \[x\]' .scratch/wisp/issues/342-privatefile-open-follows-a-swapped-ancestor-on-both-platforms-no-ticket-owns-it.md)"
echo "T342_OPEN=$(grep -c '^- \[ \]' .scratch/wisp/issues/342-privatefile-open-follows-a-swapped-ancestor-on-both-platforms-no-ticket-owns-it.md)"
echo "T348_DONE=$(grep -c '^- \[x\]' .scratch/wisp/issues/348-refusal-path-remove-follows-an-ancestor-link-and-deletes-a-foreign-file.md)"
echo "T348_OPEN=$(grep -c '^- \[ \]' .scratch/wisp/issues/348-refusal-path-remove-follows-an-ancestor-link-and-deletes-a-foreign-file.md)"
} > "$LOG" 2>&1
echo "RC=$? LINES=$(wc -l < "$LOG")"; cat "$LOG"
