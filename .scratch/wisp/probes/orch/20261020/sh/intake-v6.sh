#!/usr/bin/env bash
# 编排者收件尺 v6：353-v1 交件复量 + 我那枚"三笔同串 2235e21e"的口径对质。
cd "/d/work/workspace/projects plans/Wisp" || exit 9

echo "=== ① 名册三把尺 ==="
for P in "票 353 终裁腿 353-v1" "票 353 落地腿 353-r1" "票 353 普查腿 353-a1" "票 139 终裁腿 139-v1" "gw-1 量色腿 gw-1"; do
  a=$(git log --format='%s' | grep -c "^$P")
  b=$(git log --all --format='%s' | grep -c "^$P")
  c=$(git log --format='%H%x09%s' | awk -F'\t' -v P="$P" 'index($2,P)==1{n++} END{print n+0}')
  echo "[$P] 锚定HEAD=$a 锚定--all=$b 内存=$c"
done

echo
echo "=== ② 逐枚 name-only（.go／越界／件数）==="
for c in 3ac9d761 4fb6c088 bf6aeec2 43f0742e 35fffa29 faea245d 15b4ee83 1bdf94c9 1055a835 47ac2ec0; do
  files=$(git show --name-only --format= $c | grep -c .)
  go=$(git show --name-only --format= $c | grep -c '\.go$')
  out=$(git show --name-only --format= $c | grep -vcE '^(\.scratch/wisp/probes/353/v1/|docs/evidence/s1/353-v1-verdict\.md)')
  echo "   $c 件=$files .go=$go 越界=$out"
done

echo
echo "=== ③ 判语主件 blob 量（声称 24071 B／69 行／CR 0）==="
git cat-file blob HEAD:docs/evidence/s1/353-v1-verdict.md > /tmp/v353.txt
echo "   bytes=$(wc -c </tmp/v353.txt) lines=$(wc -l </tmp/v353.txt) CR=$(tr -cd '\r' </tmp/v353.txt | wc -c) 末字节=$(tail -c1 /tmp/v353.txt | od -An -tx1 | tr -d ' ')"
echo -n "   被禁码位 U+00A7 枚数（正控＝同尺打自造串）："; grep -o '§' /tmp/v353.txt | wc -c
printf 'AA\302\247BB\302\247' > /tmp/pc.txt; echo "   正控自造 2 枚 => 尺读 $(grep -o '§' /tmp/pc.txt | wc -c)"

echo
echo "=== ④ ★我那枚 2235e21e 的口径对质（逐枚锚 × 逐枚文件 × 两把尺）==="
for f in internal/winsec/winsec.go internal/winsec/pinnedleaf_342_windows.go; do
  echo "  -- $f"
  for a in ecffb93b fa395b33 110ed00b ffc3c62a HEAD; do
    g=$(git rev-parse --short=8 $a:$f 2>/dev/null)
    s=$(git cat-file blob $a:$f 2>/dev/null | sha256sum | cut -c1-8)
    echo "     $a gitshort=$g sha256_8=$s"
  done
done
echo "   本腿（编排者）先前自报＝三笔同串 2235e21e（⛔ 记成口径前未具名尺形）"

echo
echo "=== ⑤ 乙本体的删除原语照旧 + :207-211 段四锚同串（353-v1 第 1 件）==="
echo -n "   Remove(base) 命中 HEAD／fa395b33／BASE(0fb399c3)："
for a in 0fb399c3 fa395b33 HEAD; do
  echo -n " $(git cat-file blob $a:internal/winsec/pinnedleaf_342_windows.go 2>/dev/null | grep -c 'p\.root\.Remove(base)')"
done; echo
echo -n "   :207-211 段 sha256（四锚）："
for a in 0fb399c3 a59f5630 fa395b33 HEAD; do
  echo -n "$(git cat-file blob $a:internal/winsec/pinnedleaf_342_windows.go 2>/dev/null | sed -n '207,211p' | sha256sum | cut -c1-8) "
done; echo
echo "   353-v1 声称＝a59f5630/fa395b33/ffc3c62a/HEAD 四枚同串 6c2448515378（各 5 行）"
echo -n "   同样取 HEAD 的 :207-211 行数校验："; git cat-file blob HEAD:internal/winsec/pinnedleaf_342_windows.go | sed -n '207,211p' | wc -l

echo
echo "=== ⑥ 353 AC#2/AC#3 翻框凭据现量（gw-1 已交色，⛔ 用区间当枚数）==="
echo -n "   gw-1 味③ winsec 新增红行："; grep -m1 -oE '新增红[^|]{0,24}' .scratch/wisp/probes/gw/1/03-winsec.md
echo -n "   gw-1 味① 五把门 rc 行："; grep -m2 -oE 'GOFLAGS= go build \./\.\.\.[^|]{0,18}' .scratch/wisp/probes/gw/1/01-gates.md
echo "   五把门 rc 逐把："; grep -oE '\*\*[01] ⇄ [01]\*\*' .scratch/wisp/probes/gw/1/01-gates.md | head -6
echo -n "   353 票面现量：未勾="; grep -c '^- \[ \]' .scratch/wisp/issues/353-*.md
echo -n "   已勾="; grep -c '^- \[x\]' .scratch/wisp/issues/353-*.md

echo
echo "=== ⑦ 139 翻框凭据：gw-1 的 agent 成对色 ==="
grep -nE 'internal/agent' .scratch/wisp/probes/gw/1/02-five-packages.md | cut -c1-125

echo
echo "=== ⑧ ★两腿立票草案的号位撞车（都建议 354）==="
echo -n "   盘面最大号="; ls .scratch/wisp/issues/ | grep -E '^[0-9]+-' | sed -E 's/-.*//' | sort -n | tail -1
for n in 354 355 356; do echo -n "   $n 在册文件="; ls .scratch/wisp/issues/${n}-*.md 2>/dev/null | wc -l; done
echo "   354 草案头（353-v1）："; sed -n '22,24p' .scratch/wisp/probes/353/v1/02-ticket-draft-secC.md | cut -c1-150

echo
echo "=== ⑨ 台账号位现量（追加前先验）==="
for a in A1071 A1072 A1073 A1074 A1075; do echo -n "   $a="; grep -c "^## $a" docs/reports/pending-and-issues.md; done
