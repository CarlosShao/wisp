#!/bin/sh
# 编排者 A846 复尺：298-v1 的承重读数逐枚自己再跑一遍（静态＋格式尺；⛔ 跑 cmd/wisp 整包，那发归编队安静窗口）
cd "D:/work/workspace/projects plans/Wisp" || exit 9
D=.scratch/wisp/probes/orch/20261011
OUT=$D/logs/orch-298v1-verify.txt
F3=.scratch/wisp/probes/308/a1/logs/orch-298v1-verify-35r3.diff
F2=.scratch/wisp/probes/308/a1/logs/orch-298v1-verify-35r2.diff
mkdir -p $D/logs
P35R2=cmd/wisp/panel_transport_35r2_test.go
P35R3=cmd/wisp/panel_inbound_guards_35r3_test.go
PRE=48e7b0a4      # 18d43f8e^（改前那棵树）
POST=18d43f8e     # 落地笔
FIRST3=7f9d6e40   # 35r3 首笔
FIRST2=286a7f30   # 35r2 首笔
TMPW=$(mktemp -d)
{
  echo "date=$(date '+%Y-%m-%d %H:%M:%S %z')"
  echo "HEAD=$(git rev-parse --short HEAD)"

  echo "### 1) 同一性：腿锚 6414a4bb ⇄ 我这锚 $PRE 两枚件同 blob？"
  for p in $P35R3 $P35R2; do
    if cmp -s <(git show 6414a4bb:$p) <(git show $PRE:$p); then echo "  BLOB_IDENTICAL $p rc=0"; else echo "  ★BLOB_DIFF $p"; fi
  done

  echo "### 2) AC#1 三把尺（我⛔ 转抄腿的）"
  echo "  numstat(pre->post)="; git diff --numstat $PRE $POST -- $P35R3 $P35R2 | sed 's/^/    /'
  echo "  git diff -w --numstat 行数=$(git diff -w --numstat $PRE $POST -- $P35R3 $P35R2 | grep -c '')"
  for p in $P35R3 $P35R2; do
    a=$(git show $PRE:$p | tr -d '[:space:]' | wc -c | tr -d ' ')
    b=$(git show $POST:$p | tr -d '[:space:]' | wc -c | tr -d ' ')
    if git show $PRE:$p | tr -d '[:space:]' | cmp -s - <(git show $POST:$p | tr -d '[:space:]'); then echo "  NONWS_BYTES_IDENTICAL $p $a==$b"; else echo "  ★NONWS_DIFF $p $a vs $b"; fi
  done
  echo "  名册 $POST=$(git show --name-only --format= $POST | grep -c .) 其中非 .go=$(git show --name-only --format= $POST | grep -vc '\.go$')"

  echo "### 3) AC#0 改前 blob 的脏形（gofmt -d 的 hunk 头＋退码）"
  mkdir -p $TMPW/first
  git archive tar $PRE | tar -x -C $TMPW >/dev/null 2>&1
  for p in $P35R3 $P35R2; do
    b=$(basename $p)
    gofmt -d "$TMPW/$p" > "$D/logs/tmp-298-$(echo $b | sed 's/_test\.go$//').gofmt-d.txt" 2>&1; rc=$?
    echo "  $b gofmt -d rc=$rc hunks=$(grep -c '^@@' "$D/logs/tmp-298-$(echo $b | sed 's/_test\.go$//').gofmt-d.txt") 头=$(grep -m1 '^@@' "$D/logs/tmp-298-$(echo $b | sed 's/_test\.go$//').gofmt-d.txt") 行数=$(grep -c '' "$D/logs/tmp-298-$(echo $b | sed 's/_test\.go$//').gofmt-d.txt")"
    gofmt -l "$TMPW/$(dirname $p)" > /dev/null 2>&1; echo "  同目录 gofmt -l rc=$? 命中=$(gofmt -l $TMPW/cmd/wisp 2>/dev/null | grep -c .)"
  done

  echo "### 4) AC#2 CRLF 假枚（models.go）＋扩射程"
  git ls-files --eol cmd/wisp/models.go | sed 's/^/  /'
  echo "  CR=$(tr -cd '\r' < cmd/wisp/models.go | wc -c | tr -d ' ') 行数=$(grep -c '' cmd/wisp/models.go)"
  echo "  tracked 里 w/crlf 枚数=$(git ls-files --eol cmd internal tools | grep -c 'w/crlf')"
  echo "  非 .scratch 的那些=$(git ls-files --eol | grep 'w/crlf' | grep -v '^\.scratch' | awk '{print $NF}' | tr '\n' ' ')"

  echo "### 5) 当前 HEAD 那两枚件的格式状态（三把尺并排＋射程）"
  git archive tar HEAD | tar -x -C $TMPW 2>/dev/null
  echo "  blob 尺 cmd/wisp（HEAD 导出树）gofmt -l 命中=$(gofmt -l $TMPW/cmd/wisp 2>/dev/null | grep -c .)  件名=$(gofmt -l $TMPW/cmd/wisp 2>/dev/null | xargs -r -n1 basename | tr '\n' ' ')"
  echo "  工作树尺 gofmt -l cmd/wisp 命中=$(gofmt -l cmd/wisp 2>/dev/null | grep -c .)  件名=$(gofmt -l cmd/wisp 2>/dev/null | tr '\n' ' ')"
  GF=$(go env GOPATH)/bin/gofumpt.exe
  echo "  gofumpt 工作树 cmd/wisp 命中=$( $GF -l cmd/wisp 2>/dev/null | grep -c . ) 件名=$( $GF -l cmd/wisp 2>/dev/null | tr '\n' ' ')"
  echo "  gofumpt 导出树(HEAD blob) cmd/wisp 命中=$( $GF -l $TMPW/cmd/wisp 2>/dev/null | grep -c . ) 件名=$( $GF -l $TMPW/cmd/wisp 2>/dev/null | sed "s|$TMPW/||" | tr '\n' ' ')"
  echo "  全仓 tracked .go 的 gofmt -l=$(git ls-files '*.go' | xargs gofmt -l 2>/dev/null | grep -c .) 其中 .scratch=$(git ls-files '*.go' | xargs gofmt -l 2>/dev/null | grep -c '^\.scratch')"

  echo "### 6) 归因链：两枚件的首笔与 $POST 的关系"
  for pair in "$P35R3:$FIRST3" "$P35R2:$FIRST2"; do
    p=$(echo $pair | cut -d: -f1); f=$(echo $pair | cut -d: -f2)
    echo "  oldest($p)=$(git log --format=%h -- $p | tail -1) 腿声称首笔=$f 该笔存在=$(git rev-parse --verify --quiet $f >/dev/null && echo YES || echo NO)"
    echo "    首笔 blob=$(git rev-parse --short $f) 与 $PRE 同内容？$(cmp -s <(git show $f:$p) <(git show $PRE:$p) && echo YES || echo NO) 行数差=$(echo "$(git show $PRE:$p | grep -c '') - $(git show $f:$p | grep -c '')" | bc)"
  done
  echo "### 7) 298-v1 名册尺的口径（它报：顶层 PASS=285 FAIL=2 SKIP=1／名册 288 行）＝⛔ 我复跑整包，只核它件上那两枚红名在册⛔ 在 305 家"
  echo "  那两枚红的定义处=$(git grep -l 'TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe\|TestAC14GoSideEvalPushReachesThePage' HEAD -- cmd/wisp | tr '\n' ' ')"
  echo "  本票射程文件（两枚 35 件）里⛔ 含它们=$(git grep -c 'TestAC13ColdStart\|TestAC14GoSideEval' HEAD -- $P35R2 $P35R3 | grep -c .)"
  echo "TMPW=$TMPW"
} > $OUT 2>&1
cat $OUT
