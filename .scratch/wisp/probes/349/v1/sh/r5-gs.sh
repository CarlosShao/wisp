#!/bin/sh
cd '/d/work/workspace/projects plans/Wisp' || exit 9
BASE=070de7fa1b38b5c5d52445564c7ae19d1a73d6bd
LEG=ed8db1ea27eb50d3057dd0cd7d5c955a13e978fc
HDR=9f4e4f693c25ac92331a822469cc3cac7d90bdb7
L=/d/tmp/wisp/20261019/349v1/logs2
# ---- GS-B: 整文件 perl 状态机（保留字符串字面量内容）→ 删全部空白 → sha256
gsB() { git cat-file blob "$1:$2" | perl /d/tmp/wisp/20261019/349v1/sh/strip-b.pl | tr -d '\n' | sha256sum | cut -c1-64; }
# ---- GS-C: 先删整行 ^[[:space:]]*// 的行 → 剩下的行用 perl 状态机剥 /* */ 与行尾 // → 删空白 → sha256
gsC() { git cat-file blob "$1:$2" | grep -vE '^[[:space:]]*//' | perl /d/tmp/wisp/20261019/349v1/sh/strip-b.pl | tr -d '\n' | sha256sum | cut -c1-64; }
# 非注释行数（整行注释形，两形各自的"行数"读数）
ncount() { git cat-file blob "$1:$2" | grep -cvE '^[[:space:]]*//|^[[:space:]]*$' || true; }
for f in internal/risk/blacklist.go internal/risk/provenance.go internal/risk/syncdirs.go; do
  s=$(basename $f)
  bb=$(gsB $BASE $f); lb=$(gsB $LEG $f); hb=$(gsB $HDR $f)
  bc=$(gsC $BASE $f); lc=$(gsC $LEG $f); hc=$(gsC $HDR $f)
  nb=$(ncount $BASE $f); nl=$(ncount $LEG $f)
  echo "FILE=$s"
  echo "  GS-B base=$bb"
  echo "  GS-B del=$lb"
  echo "  GS-B HEADnow=$hb"
  echo "  GS-B base==del : $([ "$bb" = "$lb" ] && echo IDENTICAL || echo DIFFER)"
  echo "  GS-B del==HEAD : $([ "$lb" = "$hb" ] && echo IDENTICAL || echo DIFFER)"
  echo "  GS-C base=$bc"
  echo "  GS-C del=$lc"
  echo "  GS-C HEADnow=$hc"
  echo "  GS-C base==del : $([ "$bc" = "$lc" ] && echo IDENTICAL || echo DIFFER)"
  echo "  GS-C del==HEAD : $([ "$lc" = "$hc" ] && echo IDENTICAL || echo DIFFER)"
  echo "  noncomment-nonblank-lines base=$nb del=$nl $([ "$nb" = "$nl" ] && echo SAME || echo DIFF)"
done
echo '--- positive control: GS-B must DIFFER if I inject a code change (sanity of the ruler) ---'
git cat-file blob $BASE:internal/risk/blacklist.go | sed 's/^package /package XX/' | perl /d/tmp/wisp/20261019/349v1/sh/strip-b.pl | tr -d '\n' | sha256sum | cut -c1-16
echo '--- and must NOT differ if I only change a comment (proof ruler is comment-blind) ---'
git cat-file blob $BASE:internal/risk/blacklist.go | sed 's:// AC#1):// AC#1)ZZZ:' | perl /d/tmp/wisp/20261019/349v1/sh/strip-b.pl | tr -d '\n' | sha256sum | cut -c1-16
echo RC=$?
