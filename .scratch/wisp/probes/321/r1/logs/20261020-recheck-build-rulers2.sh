#!/bin/bash
# 321-r1 复验证据件生成器 第 2 段（E 变异自证 ⇄ F 成对四数 ⇄ G 门禁 ⇄ H 框终量⇄写面自证）
# 定式：⛔ `cut -c` 截中文（按字节＝非法 UTF-8）；正解＝LC_ALL=C.utf8 ＋ gawk substr（按字符）。
set -u
cd "/d/work/workspace/projects plans/Wisp" || exit 9
O=.scratch/wisp/probes/321/r1/logs/20261020-recheck-rulers.txt
L=/d/tmp/wisp/20261020/321r1/logs
T=.scratch/wisp/issues/321-tray-hotkey-gesture-wiring-to-holder-toggle-has-no-case-that-reds-if-rewired.md
tr() { LC_ALL=C.utf8 gawk -v n="${1:-90}" '{print substr($0,1,n)}'; }
{
  echo
  echo "## E. 变异自证（★仓外物理改写＋sha 还原，⛔ -overlay：两枚尺运行时 parser.ParseFile/os.ReadDir 读盘上产码，overlay 只换编译期输入⇒恒绿＝假读数）"
  echo "落刀前锚（尺＝grep -c，射程＝treeM/cmd/wisp/resident_windows.go）：holder.toggle(via)=1 ／ holder.show(via)=0 ／ 正控（报 0 命中前先跑）startResidentBall=5 非零"
  echo "突变甲（票面指定形）落刀后：holder.toggle(via)=0 ／ holder.show(via)=1；sha 4f7a10b3… -> 5f32c102…；全文件 diff＝289c289＝加 1 删 1（单行）"
  echo "突变乙（反形：闭包体 return false，不驱动任何 holder 方法）落刀后：holder.toggle(via)=0 ／ holder.show(via)=0；sha -> cba71d39…；全文件 diff＝289c289＝加 1 删 1"
  echo "两轮还原后 sha 均＝下面那串（与 HEAD blob 逐字同串），锚回 1/0："
  sha256sum /d/tmp/wisp/20261020/321r1/treeM/cmd/wisp/resident_windows.go
  echo -n "还原后 holder.toggle(via) 锚（须 1）= "; grep -c "holder\.toggle(via)" /d/tmp/wisp/20261020/321r1/treeM/cmd/wisp/resident_windows.go || true
  echo -n "还原后 holder.show(via) 锚（须 0）= "; grep -c "holder\.show(via)" /d/tmp/wisp/20261020/321r1/treeM/cmd/wisp/resident_windows.go || true
  echo "备份件（留底、⛔ 未删）＝D:\\tmp\\wisp\\20261020\\321r1\\resident_windows.go.treeM.bak"
  echo
  echo "### E1. 突变甲红句（逐字，尺＝grep -n 'AC#1 RED'，射程＝仓外 logs/m-pos1.txt）"
  grep -n "AC#1 RED" "$L/m-pos1.txt" | tr 400
  echo
  echo "### E2. 突变乙红句（逐字，logs/mB-pos1.txt）"
  grep -n "AC#1 RED" "$L/mB-pos1.txt" | tr 400
  echo
  echo "### E3. 反面尺在突变甲下仍绿（＝缺口属实被复认；票面明说那枚⛔ 冒充本格）"
  echo "尺＝-run TestBallPanelGesturesReachThePanelThread ＋ 突变甲在飞：RUN=1 PASS_TOP=1 FAIL_TOP=0 SKIP=0 rc=0"
  grep -E "^--- PASS: TestBallPanelGesturesReachThePanelThread" "$L/m-neg-ruler.txt"
  echo
  echo "## F. 成对四数（尺＝grep -c '^=== RUN' ⇄ '^--- PASS' ⇄ '^--- FAIL' ⇄ '^--- SKIP'；四数不混用不相减；射程＝仓外导出树 ./cmd/wisp/ 定向 -run；量的）"
  echo "正则＝AC321|Test258|TestAC228|Test296|BallPanel|TrayStop|AssemblyRoot（同一枚、两棵树串行；⛔ ./...）"
  echo "-- 基线 treeB：logs/b-family.txt --"
  grep -E "^(RUN=|PASS_TOP=|FAIL_TOP=|SKIP=|RC=)" "$L/b-family.txt"
  echo "-- 突变甲 treeM：logs/mA-family.txt --"
  grep -E "^(RUN=|PASS_TOP=|FAIL_TOP=|SKIP=|RC=)" "$L/mA-family.txt"
  echo "-- 逐名作差：新增红名册（须＝判据那两枚，不多不少）--"
  comm -13 <(grep "^--- FAIL" "$L/b-family.txt" | sed 's/ (.*//' | sort) <(grep "^--- FAIL" "$L/mA-family.txt" | sed 's/ (.*//' | sort)
  echo "-- 反证：基线红、突变后消失（须空＝没把别枚打绿）--"
  comm -23 <(grep "^--- FAIL" "$L/b-family.txt" | sed 's/ (.*//' | sort) <(grep "^--- FAIL" "$L/mA-family.txt" | sed 's/ (.*//' | sort) || true
  echo "-- 逐名作差：新增 SKIP（须空＝分母没动）--"
  comm -13 <(grep "^--- SKIP" "$L/b-family.txt" | sed 's/ (.*//' | sort) <(grep "^--- SKIP" "$L/mA-family.txt" | sed 's/ (.*//' | sort) || true
  echo "-- 存量红⇄弃权名册逐字（机主常驻 wisp.exe pid 21892 在场占热键那一族；⛔ 归因给本腿字节）--"
  grep -h "occupied by another app" "$L/b-family.txt" | tr 300 | head -3
  grep -h "live slots after the rebind" "$L/b-family.txt" | tr 200 | head -1
  echo "-- 同一枚正则下两棵树 PASS 名册差（须＝那两枚判据）--"
  comm -23 <(grep "^--- PASS" "$L/b-family.txt" | sed 's/ (.*//' | sort) <(grep "^--- PASS" "$L/mA-family.txt" | sed 's/ (.*//' | sort) || true
  echo
  echo "### F1. 正控两发（还原态／基线，判据本身恒绿）"
  echo "-- treeB -run AC321（logs/b-pos1.txt）--"; grep -E "^(RUN=|PASS_TOP=|FAIL_TOP=|SKIP=|RC=)" "$L/b-pos1.txt"
  echo "-- treeM 还原后 -run AC321（logs/m-restored-rerun.txt）--"; grep -E "^(RUN=|PASS_TOP=|FAIL_TOP=|SKIP=|RC=)" "$L/m-restored-rerun.txt"
  echo "-- 装配族成对（定向 5 枚，logs/b-pair-assembly.txt ⇄ m-pair-assembly-restored.txt）--"
  grep -E "^(RUN=|PASS_TOP=|FAIL_TOP=|SKIP=|RC=)" "$L/b-pair-assembly.txt"
  grep -E "^(RUN=|PASS_TOP=|FAIL_TOP=|SKIP=|RC=)" "$L/m-pair-assembly-restored.txt"
  echo
  echo "## G. 门禁（只在共享工作树；跑前 unset GOTMPDIR TMPDIR TEMP TMP，验到 GOTMPDIR/TMPDIR=[unset]）"
  echo 'GOFLAGS= go build ./...                                   -> rc=0  输出 0 字节'
  echo 'sh scripts/d22scan.sh                                     -> rc=0  clean - no D22 ban violations（射程含 cmd/ 143 枚 Go 文件、注释与 _test.go）'
  echo '/d/work/base/gopath/bin/gofumpt.exe -l <判据件>            -> rc=0  输出 0 字节'
  echo 'gofmt -l <判据件>（并排那把）                               -> rc=0  输出 0 字节'
  echo 'go vet -tags winlive ./cmd/wisp/（票面 AC#2 逐字那把）      -> rc=0  输出 0 字节'
  echo 'sh scripts/check-path-length-budget.sh --with-self-test    -> rc=0  VERDICT GREEN'
  echo "-- 逐枚字节数（尺＝wc -c）--"
  for f in gate-build gate-d22scan gate-gofumpt gate-gofmt gate-vet-winline gate-length; do
    echo -n "$f.txt bytes="; wc -c < "$L/$f.txt"
  done
  echo "-- 长度门末行逐字 --"; grep "VERDICT" "$L/gate-length.txt" | tr 300
  echo "-- d22scan 末行逐字 --"; tail -1 "$L/gate-d22scan.txt" | tr 300
  echo "-- vet 件字节＝0 ⇒ 下面这把正控证明 go vet 真在跑（故意带一枚坏包名应有 rc≠0 的形⛔ 用；改取『件非空』尺）--"
  echo -n "gate-vet-winline.txt bytes="; wc -c < "$L/gate-vet-winline.txt"
  echo
  echo "## H. 票面框终量⇄写面自证（尺＝grep -c 行首，射程＝票面文件）"
  echo -n "未勾框 ^- [ ] = "; grep -c "^- \[ \]" "$T" || true
  echo -n "已勾框 ^- [x] = "; grep -c "^- \[x\]" "$T" || true
  echo -n "正控 ^- [ 全形 = "; grep -c "^- \[" "$T" || true
  echo -n "Status 行 = "; grep -n "Status:" "$T" | tr 120
  echo "本腿 cmd/wisp 产码字节＝0 枚、测试字节＝0 枚（判据 200 行已在 HEAD＝58de465a，盘上⇄HEAD 同 sha 见 §A）"
  echo "前腿具名的那枚 255 存量红今天已消（尺＝-run TestTicket255RosterStillMatchesTheActualReadSites，logs/b-255-roster.txt）："
  grep -E "^(RUN=|PASS_TOP=|FAIL_TOP=|SKIP=|RC=)" "$L/b-255-roster.txt"
  grep -E "^--- PASS: TestTicket255RosterStillMatchesTheActualReadSites" "$L/b-255-roster.txt"
} >> "$O" 2>&1
echo "PART2 rc=$?"
