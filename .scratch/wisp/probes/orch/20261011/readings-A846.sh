#!/bin/sh
# A846 批读数（编排者自己复跑；⛔ 转抄腿件；含我上一发那两把坏尺的更正）
cd "$(git rev-parse --show-toplevel)" || exit 9
D=.scratch/wisp/probes/orch/20261011
LOG=$D/logs/orch-A846-readings.txt
F1=cmd/wisp/panel_inbound_guards_35r3_test.go
F2=cmd/wisp/panel_transport_35r2_test.go
GF=$(go env GOPATH)/bin/gofumpt.exe
PRE=48e7b0a4
POST=18d43f8e
NOW=$(git rev-parse --short HEAD)

{ echo "date=$(date '+%Y-%m-%d %H:%M:%S %z')  HEAD=$NOW"
  echo "⛔ 本件⛔ 跑 go build/vet/test（Go 编译面此刻由 309-r1＋308-r1 两枚写腿持有）"
} > $LOG 2>&1

# ---- 1) 298：那两枚件自 18d43f8e 之后无人再动（三把尺） ----
{ echo "### 1) 两枚件 18d43f8e ⇄ HEAD 同一性＋工作树脏度"
  for p in $F1 $F2; do
    if git show $POST:$p | cmp -s - <(git show HEAD:$p); then echo "  UNMOVED $p rc=0 ($POST == HEAD)"; else echo "  MOVED $p"; fi
  done
  echo "  git status --porcelain -- cmd/wisp 行数=$(git status --porcelain -- cmd/wisp | grep -c '')"
} >> $LOG 2>&1

# ---- 2) 298 AC#0：diff 内容行的三把明尺（腿件写"内容行 12"） ----
{ echo "### 2) 腿件自带的改前 blob gofmt -d 原件，三把明尺并排"
  for p in $F1 $F2; do
    n=$(basename $p)
    L=.scratch/wisp/probes/298/v1/logs/ac0-gofmt-d-${n%_test.go}-preblob.txt
    [ -f "$L" ] || { echo "  件⛔ 在 $L"; continue; }
    echo "  $n 件=$L wc_l=$(grep -c '' $L) 尺A'^[+-]'=$(grep -c '^[+-]' $L) 尺B'^[+-][^+-]'=$(grep -cE '^[+-][^+-]' $L) 尺C'-'=$(grep -c '^-' $L) 尺D'+'=$(grep -c '^+' $L)"
  done
  echo "  腿声称的'内容行 12'⛔ 可从上述任一把复现；12 那枚数出自它自己那件文件名写歪的失败尝试件（1,11d0 ＋ 11 行 < 前缀）"
} >> $LOG 2>&1

# ---- 3) 298 AC#2：CRLF 名册的三档射程（更正我上一发那把贴错标签的尺） ----
{ echo "### 3) CRLF 假枚名册三档射程（上一发我把'非 .scratch'那把尺打成了全仓名册＝标签坏，这发改对）"
  echo "  档1 窄域（cmd internal tools）：$(git ls-files --eol cmd internal tools | grep -c 'w/crlf') 枚"
  echo "  档1 件名：$(git ls-files --eol cmd internal tools | grep 'w/crlf' | awk '{print $NF}' | tr '\n' ' ')"
  ALL=$(git ls-files --eol | grep 'w/crlf' | awk '{print $NF}')
  echo "  档2 全仓 tracked 命中=$(printf '%s\n' "$ALL" | grep -c '')"
  NON=$(printf '%s\n' "$ALL" | grep -v '^\.scratch/')
  echo "  档3 全仓剥 .scratch 后=$(printf '%s\n' "$NON" | grep -c '') 枚：$(printf '%s\n' "$NON" | tr '\n' ' ')"
  echo "  core.autocrlf（只读）=$(git config --get core.autocrlf)"
} >> $LOG 2>&1

# ---- 4) 298 AC#3 那枚"三件事⇄四件事"的文案矛盾（逐字两处） ----
{ echo "### 4) 口径句在两处的逐字（票面 ⇄ 腿凭据件）"
  echo "  票面：$(grep -o '这[三四]件事.\{0,60\}' .scratch/wisp/issues/298-*.md | head -2 | tr '\n' '|')"
  echo "  腿件：$(grep -o '这[三四]件事.\{0,60\}' .scratch/wisp/probes/298/v1/20-*.md | head -2 | tr '\n' '|')"
} >> $LOG 2>&1

# ---- 5) 298 ⓑ gofumpt 残留：减前⇄减后两把（HEAD blob 导出树，⛔ 工作树） ----
T1=$(mktemp -d); T2=$(mktemp -d)
{ echo "### 5) gofumpt 那枚残留的 减前⇄减后（两棵导出树，同一把尺）"
  for pair in "pre=$PRE" "post=$POST" "head=$NOW"; do
    tag=${pair%%=*}; sha=${pair##*=}
    dst=$T1; [ "$tag" = "post" ] && dst=$T2
    git archive --format=tar $sha cmd/wisp | tar -x -C $dst; rc=$?
    hits=$(cd $dst && $GF -l cmd/wisp 2>/dev/null | grep -c '')
    names=$(cd $dst && $GF -l cmd/wisp 2>/dev/null | tr '\n' ' ')
    echo "  $tag($sha) archive_rc=$rc gofumpt_blob_hits=$hits 件名=$names"
  done
  echo "  同两棵树的 gofmt 尺：pre=$(cd $T1 && gofmt -l cmd/wisp | grep -c '') post=$(cd $T2 && gofmt -l cmd/wisp | grep -c '')"
} >> $LOG 2>&1
rm -rf $T1 $T2 2>/dev/null

# ---- 6) 308：我那"两把台面"其实同一把（自纠）＋在场尺逐字 ----
{ echo "### 6) temp env 台面在场尺（Go 只认 TMP/TEMP；TMPDIR 在 windows 档⛔ 读）"
  cat > $D/tmp-volprobe.go <<'EOF'
package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	d, err := os.MkdirTemp("", "VOLPROBE")
	fmt.Printf("TMP=%q TEMP=%q TMPDIR=%q\n", os.Getenv("TMP"), os.Getenv("TEMP"), os.Getenv("TMPDIR"))
	fmt.Printf("MkdirTemp=%q err=%v VolumeName=%q len=%d\n", d, err, filepath.VolumeName(d), len(filepath.VolumeName(d)))
}
EOF
  echo "  （这枚⛔ 编译＝本件只列在场尺读过的 env 形；实跑读数见 logs/orch-308a1-verify.txt 与 probes/308/a1/logs/probe-*.txt）"
  echo "  编排者自纠：verify 那两发（TMPDIR=/tmp；TMP=/tmp TEMP=/tmp）都在'带盘符'一侧 ⇒ 我那是【两写法一把台面】，⛔ 等于票面要的'无卷标台面'那一发"
} >> $LOG 2>&1

cat $LOG
