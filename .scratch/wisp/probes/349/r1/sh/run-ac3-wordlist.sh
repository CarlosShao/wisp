#!/bin/sh
# 票 349 AC#3 第二把尺：有界词表形（343-v1 提的 `^a[A-Z]…`＝票 350 那一族）我自己独立量一遍。
# ★⛔ 采信任何腿给的枚数：名册从【注释文本】现抽（只数 `//` 开头行里的驼峰 token），
#   定义侧用三形对撞（func 自由形／method 形／非 func 的声明形 type|var|const 与结构体字段）。
# 两棵树都我自己的 git archive 导出：treeBase＝起手锚、treeDel＝交件锚。
set -u
REPO='/d/work/workspace/projects plans/Wisp'
cd "$REPO" || exit 90
ROOT=/d/tmp/wisp/20261019/349r1
L=.scratch/wisp/probes/349/r1/logs/ac3
mkdir -p "$L"
A_BASE=7b633fa9713b0514806c8b808ed4db91b1485cdc
A_DEL=db2f9eab7bcffd827240763699aa0f4ad7ff3266
{
for side in base del; do
  if [ "$side" = base ]; then ANCHOR="$A_BASE"; T="$ROOT/treeBase"; else ANCHOR="$A_DEL"; T="$ROOT/treeDel"; fi
  echo "############ side=$side tree=$T anchor=$ANCHOR ############"
  TMPF="$ROOT/wl-$side.txt"
  : > "$TMPF"
  for f in internal/risk/blacklist.go internal/risk/provenance.go internal/risk/syncdirs.go; do
    git show "$ANCHOR:$f" >> "$TMPF"
  done
  # 名册：三枚产码文件【注释行】里的 a[A-Z]… token（去重）
  grep -E '^[[:space:]]*//' "$TMPF" | grep -oE '\ba[A-Z][A-Za-z0-9_]*\b' | sort -u > "$ROOT/wl-$side-names.txt"
  printf 'comment a[A-Z] names in the 3 production files = %s\n' "$(wc -l < "$ROOT/wl-$side-names.txt")"
  sed 's/^/    /' "$ROOT/wl-$side-names.txt"
  echo "  -- 逐枚问'包里有定义吗'（三形对撞，射程＝该树 internal/risk 全 49 枚 .go）--"
  while read -r n; do
    [ -n "$n" ] || continue
    d=$(grep -rhoE "^func $n\(|^func \([^)]*\) $n\(|^type $n |^var $n |^const $n |^[[:space:]]+$n +" "$T/internal/risk" 2>/dev/null | sort -u | head -2)
    if [ -n "$d" ]; then
      printf '  GREEN  %-14s %s\n' "$n" "$(printf '%s' "$d" | tr '\n' ';' | tr -s ' ')"
    else
      printf '  RED    %-14s 包内 0 定义形\n' "$n"
    fi
  done < "$ROOT/wl-$side-names.txt"
done
echo "############ 对照：全 package 生产件（去 _test.go）的 a[A-Z] 名册枚数（编排者现跑 HEAD=5 那把尺的口径）############"
for side in base del; do
  if [ "$side" = base ]; then ANCHOR="$A_BASE"; else ANCHOR="$A_DEL"; fi
  git ls-tree -r --name-only "$ANCHOR" -- internal/risk | grep '\.go$' | grep -v '_test\.go$' > "$ROOT/prodfiles-$side.txt"
  : > "$ROOT/prodcomments-$side.txt"
  while read -r f; do git show "$ANCHOR:$f" >> "$ROOT/prodcomments-$side.txt"; done < "$ROOT/prodfiles-$side.txt"
  grep -E '^[[:space:]]*//' "$ROOT/prodcomments-$side.txt" | grep -oE '\ba[A-Z][A-Za-z0-9_]*\b' | sort -u > "$ROOT/prodnames-$side.txt"
  printf 'side=%s prodfiles=%s whole-package-production a[A-Z] names=%s\n' "$side" "$(wc -l < "$ROOT/prodfiles-$side.txt")" "$(wc -l < "$ROOT/prodnames-$side.txt")"
  sed 's/^/    /' "$ROOT/prodnames-$side.txt"
done
} > "$L/ac3-wordlist.log" 2>&1
echo "wrote=$L/ac3-wordlist.log rc=$?"
cat "$L/ac3-wordlist.log"
