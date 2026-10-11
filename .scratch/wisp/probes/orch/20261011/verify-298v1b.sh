#!/bin/sh
# A846 复尺 v2：修我这把坏尺（v1 写的是 `git archive tar $SHA`＝fatal: not a valid object name: tar
# ⇒ 导出树是空的 ⇒ blob 那两把 gofmt/gofumpt 读数＝★假清洁）。v1 的原样读数留在 logs/orch-298v1-verify.txt⛔ 改。
cd "D:/work/workspace/projects plans/Wisp" || exit 9
D=.scratch/wisp/probes/orch/20261011
OUT=$D/logs/orch-298v1-verify-b.txt
P35R2=cmd/wisp/panel_transport_35r2_test.go
P35R3=cmd/wisp/panel_inbound_guards_35r3_test.go
PRE=48e7b0a4
POST=18d43f8e
FIRST3=7f9d6e40
FIRST2=286a7f30
T1=$(mktemp -d); T2=$(mktemp -d); T3=$(mktemp -d); T4=$(mktemp -d)
{
  echo "date=$(date '+%Y-%m-%d %H:%M:%S %z')  HEAD=$(git rev-parse --short HEAD)"
  echo "ruler_form_fixed=git archive --format=tar <sha> | tar -x -C <仓外目录>（v1 那把⛔ 语法对，导出树空＝假清洁）"

  echo "### 3-fix) 改前两枚 blob 的 gofmt -d 脏形（hunk 头／行数／退码）"
  git archive --format=tar $PRE | tar -x -C $T1 >/dev/null 2>&1; echo "  导出 $PRE 件数=$(find $T1 -name '*.go' | grep -c .)"
  for p in $P35R3 $P35R2; do
    b=$(basename $p | sed 's/_test\.go$//')
    gofmt -d "$T1/$p" > $D/logs/tmp-298b-$b.gofmt-d.txt 2>&1; rc=$?
    echo "  $b rc=$rc hunks=$(grep -c '^@@' $D/logs/tmp-298b-$b.gofmt-d.txt) 头=$(grep -m1 '^@@' $D/logs/tmp-298b-$b.gofmt-d.txt) 总行=$(grep -c '' $D/logs/tmp-298b-$b.gofmt-d.txt) 内容行=$(grep -cE '^[+-][^+]' $D/logs/tmp-298b-$b.gofmt-d.txt)"
  done
  echo "  同目录 gofmt -l（改前那棵树 cmd/wisp）命中=$(gofmt -l $T1/cmd/wisp 2>/dev/null | grep -c .) 件名=$(gofmt -l $T1/cmd/wisp 2>/dev/null | xargs -r -n1 basename | tr '\n' ' ')"

  echo "### 4-fix) HEAD blob 那两枚件＋整个 cmd/wisp 的格式三把尺"
  git archive --format=tar HEAD | tar -x -C $T2 >/dev/null 2>&1
  echo "  HEAD 导出件数=$(find $T2 -name '*.go' | grep -c .)  gofmt -l cmd/wisp 命中=$(gofmt -l $T2/cmd/wisp 2>/dev/null | grep -c .) 件名=$(gofmt -l $T2/cmd/wisp 2>/dev/null | xargs -r -n1 basename | tr '\n' ' ')"
  GF=$(go env GOPATH)/bin/gofumpt.exe
  echo "  HEAD 导出树 gofumpt -l cmd/wisp 命中=$( $GF -l $T2/cmd/wisp 2>/dev/null | grep -c . ) 件名=$( $GF -l $T2/cmd/wisp 2>/dev/null | sed "s|$T2/||" | tr '\n' ' ')"
  echo "  改后笔 $POST 导出树 gofmt -l 命中=$( git archive --format=tar $POST | tar -x -C $T3 >/dev/null 2>&1; gofmt -l $T3/cmd/wisp 2>/dev/null | grep -c .)"

  echo "### 6-fix) 归因链两枚首笔（v1 缺的枚数与 hunk，这发补齐；⛔ bc 用 \\$(( ))）"
  for pair in "$P35R3:$FIRST3" "$P35R2:$FIRST2"; do
    p=$(echo $pair | cut -d: -f1); f=$(echo $pair | cut -d: -f2)
    TF=$(mktemp -d)
    git archive --format=tar $f | tar -x -C $TF >/dev/null 2>&1
    pre_l=$(git show $PRE:$p | grep -c ''); f_l=$(git show $f:$p | grep -c '')
    pre_b=$(git show $PRE:$p | wc -c | tr -d ' '); f_b=$(git show $f:$p | wc -c | tr -d ' ')
    echo "  $p 首笔=$f 行数 首笔$f_l ⇄ 改前$pre_l（差 $((pre_l - f_l))）字节 首笔$f_b ⇄ 改前$pre_b"
    echo "    首笔那版 gofmt -d hunk 头=$(gofmt -d "$TF/$p" 2>&1 | grep -m1 '^@@') rc=$(gofmt -d "$TF/$p" >/dev/null 2>&1; echo $?)"
  done
  echo "### 7-fix) 两枚红的名字在 HEAD 的定义处（具名）＋本票射程文件枚数"
  for n in TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe TestAC14GoSideEvalPushReachesThePage; do
    echo "  $n 定义处=$(git grep -ln "func $n(" HEAD -- cmd/wisp | tr '\n' ' ')"
  done
  echo "### 8) 298-v1 那把 '27＝22＋5' 的尺我复跑（stdout/stderr 分离）"
  git ls-files '*.go' | xargs gofmt -l 2>/dev/null > $D/logs/tmp-298b-gofmt-all.txt; rc=$?
  echo "  rc=$rc 总命中=$(grep -c '' $D/logs/tmp-298b-gofmt-all.txt) .scratch=$(grep -c '^\.scratch' $D/logs/tmp-298b-gofmt-all.txt) 非.scratch=$(grep -v '^\.scratch' $D/logs/tmp-298b-gofmt-all.txt | grep -c .)"
  echo "  非 .scratch 的件名=$(grep -v '^\.scratch' $D/logs/tmp-298b-gofmt-all.txt | tr '\n' ' ')"
} > $OUT 2>&1
cat $OUT
