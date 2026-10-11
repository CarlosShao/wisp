#!/bin/sh
# A846 批读数 第二发（改对我自己 a 发里三处坏尺：② 腿件名我写歪了、④ 指错目录、⑤ pre/post 两棵树共用同一枚临时目录＝"pre=0" 读的其实是 HEAD 树）
cd "$(git rev-parse --show-toplevel)" || exit 9
D=.scratch/wisp/probes/orch/20261011
LOG=$D/logs/orch-A846-readings-b.txt
GF=$(go env GOPATH)/bin/gofumpt.exe
PRE=48e7b0a4
POST=18d43f8e
NOW=$(git rev-parse --short HEAD)

{ echo "date=$(date '+%Y-%m-%d %H:%M:%S %z')  HEAD=$NOW"
  echo "本发改的三处坏尺（a 发读数⛔ 删，逐字留在 logs/orch-A846-readings.txt）"
} > $LOG 2>&1

# ---- 2b) 腿件自带的改前 blob gofmt -d 原件（真名＝ac0-gofmt-d-35r{3,2}-preblob.txt）三把明尺 ----
{ echo "### 2b) 298-v1 件上'内容行 12'那把尺：三把明尺并排复量它自己的原件"
  for n in 35r3 35r2; do
    L=.scratch/wisp/probes/298/v1/logs/ac0-gofmt-d-$n-preblob.txt
    echo "  $n 件=$L 存在=$([ -f $L ] && echo YES || echo NO) wc_l=$(grep -c '' $L) 尺A'[+-]'=$(grep -c '^[+-]' $L) 尺B'[+-][^+-]'=$(grep -cE '^[+-][^+-]' $L) 尺C'^-'=$(grep -c '^-' $L) 尺D'^+'=$(grep -c '^+' $L)"
  done
  echo "  结论＝14 行拆成【3 枚路径头 ＋ 1 枚 hunk 头 ＋ 10 枚正文】，真变化行＝2 删 ＋ 2 加 ＝ 4 枚；'12'⛔ 出自这三把任何一把"
} >> $LOG 2>&1

# ---- 4b) "三件事 ⇄ 四件事" 的逐字两处（真身＝实现腿 r1 的 20-…，⛔ 是验收腿自己的件） ----
{ echo "### 4b) 口径句两处逐字"
  echo "  实现腿件（r1/20-…）：$(grep -o '这[三四]件事[^。]\{0,40\}' .scratch/wisp/probes/298/r1/20-*.md | head -1)"
  echo "  实现腿件逐字行：$(grep -h '三件事' .scratch/wisp/probes/298/r1/20-*.md | cut -c1-150)"
  echo "  票面逐字行：      $(grep -h '四件事' .scratch/wisp/issues/298-*.md | cut -c1-150)"
} >> $LOG 2>&1

# ---- 5b) gofmt/gofumpt 的 减前⇄减后（三棵⛔ 同一目录） ----
A=$(mktemp -d); B=$(mktemp -d); C=$(mktemp -d)
{ echo "### 5b) 三棵各自独立的导出树（尺＝git archive --format=tar <sha> | tar -x -C <仓外目录>，同一枚 gofumpt/gofmt）"
  for trio in "pre=$PRE:$A" "post=$POST:$B" "head=$NOW:$C"; do
    tag=${trio%%=*}; rest=${trio#*=}; sha=${rest%%:*}; dst=${rest##*:}
    git archive --format=tar $sha cmd/wisp | tar -x -C $dst
    echo "  $tag($sha) 件数=$(cd $dst && find cmd/wisp -name '*.go' | grep -c '') gofmt_hits=$(cd $dst && gofmt -l cmd/wisp | grep -c '') gofumpt_hits=$(cd $dst && $GF -l cmd/wisp 2>/dev/null | grep -c '') gofumpt件名=$(cd $dst && $GF -l cmd/wisp 2>/dev/null | tr '\n' ' ')"
  done
} >> $LOG 2>&1
rm -rf $A $B $C 2>/dev/null

# ---- 7) CRLF 那把尺的单位名（腿把 .go 子集写成"全仓非 .scratch"） ----
{ echo "### 7) CRLF 假枚名册的四种尺（单位名逐把写明）"
  echo "  尺1 窄域 tracked（cmd internal tools，全扩展名）=$(git ls-files --eol cmd internal tools | grep -c 'w/crlf')"
  echo "  尺2 窄域 tracked `.go`＝$(git ls-files --eol cmd internal tools | grep 'w/crlf' | grep -c '\.go$')"
  echo "  尺3 全仓 tracked（⛔ 剥 .scratch，全扩展名）=$(git ls-files --eol | grep -c 'w/crlf')"
  echo "  尺4 全仓 tracked 剥 .scratch 后=$(git ls-files --eol | grep 'w/crlf' | awk '{print $NF}' | grep -vc '^\.scratch/')"
  echo "  尺4 里 `.go`＝$(git ls-files --eol | grep 'w/crlf' | awk '{print $NF}' | grep -v '^\.scratch/' | grep -c '\.go$')（件名＝$(git ls-files --eol | grep 'w/crlf' | awk '{print $NF}' | grep -v '^\.scratch/' | grep '\.go$' | tr '\n' ' ')）"
} >> $LOG 2>&1

cat $LOG
