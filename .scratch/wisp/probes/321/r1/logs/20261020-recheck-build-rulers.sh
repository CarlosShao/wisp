#!/bin/bash
# 321-r1 复验证据件生成器（仓内写面＝.scratch/wisp/probes/321/r1/logs/20261020-recheck-rulers.txt）
# 定式一条：⛔ `cut -c` 截中文（按字节＝非法 UTF-8，公共块 台面四味 ④）。
# 正解＝LC_ALL=C.utf8 ＋ gawk substr（按字符）＝已现验有效 UTF-8。
set -u
cd "/d/work/workspace/projects plans/Wisp" || exit 9
O=.scratch/wisp/probes/321/r1/logs/20261020-recheck-rulers.txt
L=/d/tmp/wisp/20261020/321r1/logs
T=.scratch/wisp/issues/321-tray-hotkey-gesture-wiring-to-holder-toggle-has-no-case-that-reds-if-rewired.md
tr() { LC_ALL=C.utf8 gawk -v n="${1:-90}" '{print substr($0,1,n)}'; }
{
  echo "# 321-r1 复验 2026-10-20 逐字读数（每段带尺名＋射程＋量的还是算的；本件由 logs/20261020-recheck-build-rulers.sh 现生成）"
  echo
  echo "## A. 起手三把尺（盘上物证，⛔ 采信票尾 2026-10-15 那条 bullet 的任何读数）"
  echo '$ git rev-parse HEAD'; git rev-parse HEAD
  echo '$ git branch --show-current'; git branch --show-current
  echo '$ git cat-file -e HEAD:cmd/wisp/exec_holder_toggle_321_test.go && echo EXISTS'
  git cat-file -e HEAD:cmd/wisp/exec_holder_toggle_321_test.go && echo EXISTS
  echo '$ git show HEAD:cmd/wisp/exec_holder_toggle_321_test.go | wc -l'; git show HEAD:cmd/wisp/exec_holder_toggle_321_test.go | wc -l
  echo '$ sha256sum cmd/wisp/exec_holder_toggle_321_test.go   (盘上)'
  sha256sum cmd/wisp/exec_holder_toggle_321_test.go
  echo '$ git show HEAD:cmd/wisp/exec_holder_toggle_321_test.go | sha256sum   (HEAD blob)'
  git show HEAD:cmd/wisp/exec_holder_toggle_321_test.go | sha256sum
  echo
  echo "## B. 名册尺（防双派那把为什么漏＝锚定⇄非锚定对照；尺＝git log --grep，射程＝--all）"
  echo '$ git log --all --grep="^票 321" --format="%h %ad %s" | 按字符截 90'
  git log --all --grep="^票 321" --format="%h %ad %s" --date=format:"%m-%d %H:%M" | tr 90
  echo '$ git log --all --grep="票 321" --format="%h %ad %s" | 按字符截 90   (非锚定＝4 笔)'
  git log --all --grep="票 321" --format="%h %ad %s" --date=format:"%m-%d %H:%M" | tr 90
  echo '$ git log --all --grep="321-r1" --format="%h %ad %s" | 按字符截 70'
  git log --all --grep="321-r1" --format="%h %ad %s" --date=format:"%m-%d %H:%M" | tr 70
  echo '$ git merge-base --is-ancestor 58de465a HEAD'; git merge-base --is-ancestor 58de465a HEAD && echo YES
  echo
  echo "## C. 内容锚（⛔ 绝对行号；装配处今天＝resident_windows.go:288-290＝票面写的 :268-270 已漂 20 行）"
  echo '$ git grep -n "holder.toggle" HEAD -- cmd/wisp/resident_windows.go'
  git grep -n "holder\.toggle" HEAD -- cmd/wisp/resident_windows.go
  echo '$ git grep -n "withPanelHost(" HEAD -- cmd/wisp/'; git grep -n "withPanelHost(" HEAD -- cmd/wisp/
  echo '$ git grep -n "holder := newPanelHolder"   (holder 绑定＝尺认构造器不认变量名)'
  git grep -n "holder := newPanelHolder" HEAD -- cmd/wisp/resident_windows.go
  echo '$ git grep -n "func (h \*panelHolder) toggle"   (能力尺要去解析的那枚声明)'
  git grep -n "func (h \*panelHolder) toggle" HEAD -- cmd/wisp/
  echo
  echo "## D. pristine（跑前⇄跑后各验；尺＝sha256sum／git status --porcelain，射程＝共享树＋两棵仓外导出树）"
  echo '$ git status --porcelain -- cmd/wisp/ | wc -l   (共享树整趟＝0)'
  git status --porcelain -- cmd/wisp/ | wc -l
  echo '$ git show HEAD:cmd/wisp/resident_windows.go | sha256sum   (HEAD blob)'
  git show HEAD:cmd/wisp/resident_windows.go | sha256sum
  echo '$ 两棵导出树起步 sha（须与 HEAD blob 同串）'
  sha256sum /d/tmp/wisp/20261020/321r1/treeB/cmd/wisp/resident_windows.go
  echo '$ 取出后先数 .go 证非空（"0"⛔ 是读数）'
  echo -n "treeB .go="; find /d/tmp/wisp/20261020/321r1/treeB -name '*.go' | wc -l
  echo -n "treeM .go="; find /d/tmp/wisp/20261020/321r1/treeM -name '*.go' | wc -l
} > "$O" 2>&1
echo "PART1 rc=$?"
