# 2026-11-08 票411：新旧射程差异逐枚定位（G7 与 G5；只建不删）
cd "$(git rev-parse --show-toplevel)" || exit 9
A=$(git rev-parse HEAD)
oldlist=$(mktemp "${TMPDIR:-/tmp}/wisp-411-old.XXXXXX")
for w in internal/**/*.go cmd/**/*.go; do printf '%s\n' "$w"; done > "$oldlist"
# G7 Defer 开方文件集合
git grep -Ew 'Defer(Named)?' "$A" -- 'internal/**/*.go' 'cmd/**/*.go' ':!*_test.go' | cut -d: -f2 | sort -u > /tmp/w411-g7-new.txt
: > /tmp/w411-g7-old.txt
for b in $(split -l 200 -d "$oldlist" /tmp/w411-b. && ls /tmp/w411-b.*); do
	git grep -Ew 'Defer(Named)?' "$A" -- $(cat "$b") ':!*_test.go' 2>/dev/null | cut -d: -f2 >> /tmp/w411-g7-old.raw
done
sort -u /tmp/w411-g7-old.raw > /tmp/w411-g7-old.txt
echo "== G7 开方文件集合 diff（<＝new 独有 >＝old 独有）=="
diff /tmp/w411-g7-old.txt /tmp/w411-g7-new.txt
echo "== new 独有文件是否在 HEAD 树上 =="
for f in $(comm -13 /tmp/w411-g7-old.txt /tmp/w411-g7-new.txt); do
	echo "$f: HEAD=$(git cat-file -e "$A:$f" 2>/dev/null && echo 有 || echo 无) worktree=$([ -f "$f" ] && echo 有 || echo 无)"
done
echo "== old 独有文件 =="
for f in $(comm -23 /tmp/w411-g7-old.txt /tmp/w411-g7-new.txt); do
	echo "$f: HEAD=$(git cat-file -e "$A:$f" 2>/dev/null && echo 有 || echo 无) worktree=$([ -f "$f" ] && echo 有 || echo 无) tracked=$(git ls-files --error-unmatch "$f" >/dev/null 2>&1 && echo 是 || echo 否)"
done
echo "== old 清单里 HEAD 树上不存在的文件枚数 =="
cnt=0; miss=0
while IFS= read -r f; do cnt=$((cnt+1)); git cat-file -e "$A:$f" 2>/dev/null || { miss=$((miss+1)); [ $miss -le 5 ] && echo "  old清单有/HEAD无: $f"; }; done < "$oldlist"
echo "  old清单=$cnt 其中HEAD树上不存在=$miss"
echo "== HEAD 树上有而 old 清单没有的 internal/cmd .go（=worktree 缺席的 tracked 文件）=="
git ls-tree -r --name-only "$A" | grep -E '^(internal|cmd)/.*\.go$' | grep -vxFx -f <(sed 's|^|/|' /dev/null; cat "$oldlist") | head -10
git ls-tree -r --name-only "$A" | grep -E '^(internal|cmd)/.*\.go$' > /tmp/w411-head.txt
comm -23 /tmp/w411-head.txt "$oldlist" | head -10
echo "  枚数=$(comm -23 /tmp/w411-head.txt "$oldlist" | wc -l)"
rm -f "$oldlist" /tmp/w411-b.* /tmp/w411-g7-old.raw
