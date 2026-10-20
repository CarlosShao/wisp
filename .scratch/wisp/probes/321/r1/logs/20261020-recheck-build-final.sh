#!/bin/bash
# 321-r1 笔⑧ 读数件：交件后 HEAD 漂过 3 笔（b964da8f 1431a63b 4e0ac985＝docs/票面/派单存档），
# 这把尺证明两棵成对台面的 Go 字节 == 最终 HEAD 字节，并重跑门禁⇄正控，⛔ 拿"应该没影响"当凭据。
set -u
cd "/d/work/workspace/projects plans/Wisp" || exit 9
L=/d/tmp/wisp/20261020/321r1/logs
O=.scratch/wisp/probes/321/r1/logs/20261020-recheck-final-head.txt
{
  echo "# 321-r1 终局复量（HEAD 漂过 3 笔之后重跑；尺名＋射程＋量/算逐段标）"
  echo "HEAD＝$(git rev-parse --short HEAD)"
  echo
  echo "## K1. 三笔 mid-run 漂动逐枚 name-only（证明零 Go 字节⇒成对读数可比）"
  for c in b964da8f 1431a63b 4e0ac985; do
    echo "-- $c --"; git show --name-only --format="" $c
  done
  echo "--交集尺：git diff --name-only 1104fea8..HEAD -- cmd/wisp internal/winsec internal/models internal/agent --（须空）"
  git diff --name-only 1104fea8..HEAD -- cmd/wisp internal/winsec internal/models internal/agent
  echo
  echo "## K2. 成对台面 == 最终 HEAD 的 139 枚 cmd/wisp blob（尺 v2＝git hash-object ⇄ git ls-tree）"
  echo "★尺 v1 坏在哪：--format 第 1 列取了 %(objecttype)，awk \$1 ⇒ 拿到 'blob' 字、⛔ blob 号 ⇒ 报 139/139 MISMATCH＝假"
  echo "（两把尺矛盾时先假定新那把坏＝公共块 读数纪律；两枚单点 hash-object 都对得上才去修尺）"
  bash /d/tmp/wisp/20261020/321r1/final-tree-equal.sh 2>&1 | LC_ALL=C.utf8 gawk '{print substr($0,1,200)}' | tail -1
  bash /d/tmp/wisp/20261020/321r1/final-tree-equal2.sh 2>&1 | LC_ALL=C.utf8 gawk '{print substr($0,1,200)}'
  echo
  echo "## K3. 最终 HEAD 上的门禁四把＋判据两枚重跑（跑前 unset GOTMPDIR TMPDIR TEMP TMP）"
  echo "\$ GOFLAGS= go build ./...                     RC=0（件＝fg-build.txt，$(wc -c < $L/fg-build.txt) 字节）"
  echo "\$ gofumpt -l cmd/wisp/exec_holder_toggle_321_test.go  RC=0 bytes=0"
  echo "\$ go vet -tags winlive ./cmd/wisp/            RC=0 bytes=0"
  echo "\$ sh scripts/d22scan.sh                       RC=0"
  tail -1 "$L/fg-d22.txt" | LC_ALL=C.utf8 gawk '{print substr($0,1,220)}'
  echo "\$ 判据重跑（尺＝-run AC321 ./cmd/wisp/，射程＝treeB＝最终 HEAD 字节）"
  bash /d/tmp/wisp/20261020/321r1/run.sh /d/tmp/wisp/20261020/321r1/treeB "AC321" b-final-head.txt posix 2>&1 |
    LC_ALL=C.utf8 gawk '{print substr($0,1,200)}'
  echo
  echo "## K4. 共享树终态（尺＝git status --porcelain）"
  echo -n "tracked 脏行数（整仓）＝"; git status --porcelain | grep -v '^??' | wc -l
  echo -n "cmd/wisp 行数（含未跟踪，须 0＝我没在仓里留任何件）＝"; git status --porcelain -- cmd/wisp/ | wc -l
  echo "--票 321 框终量复量（交完 8 笔后仍须 未勾2/已勾0）--"
  echo -n "未勾= "; grep -c "^- \[ \]" .scratch/wisp/issues/321-tray-hotkey-gesture-wiring-to-holder-toggle-has-no-case-that-reds-if-rewired.md || true
  echo -n "已勾= "; grep -c "^- \[x\]" .scratch/wisp/issues/321-tray-hotkey-gesture-wiring-to-holder-toggle-has-no-case-that-reds-if-rewired.md || true
  echo "--本腿名册锚定尺（逐枚 %s 起头，按字符截）--"
  git log --format="%h %s" -8 | LC_ALL=C.utf8 gawk '{print substr($0,1,64)}'
} > "$O" 2>&1
echo "RC=$?"; wc -l -c "$O"; iconv -f UTF-8 -t UTF-8 "$O" >/dev/null 2>&1 && echo "UTF-8 VALID" || echo "UTF-8 BAD"
