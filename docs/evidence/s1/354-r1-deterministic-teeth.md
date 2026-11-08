# 票 354 `354-r1` 写腿判据件：AC#2 确定性仪器（`leftStandingForeignName`）＋ AC#3 成对去牙

- 日期：2026-11-08 22:2x +08。腿＝`354-r1`（写腿：AC#2＋AC#3；写面＝`internal/winsec/leftstanding_354_windows_test.go`（新建）＋本件＋票 354 Progress log 一行；**产码 `internal/winsec` 非 test 字节＝0**）。
- 锚：主树 commit **`393430dc9`**（仪器本体）；全部读数跑在 `git clone --shared` 导出树 `/d/tmp/wisp/20261108/r354/tree-r1`（HEAD＝同 `393430dc9`；三核③＝红绿取自导出树）。
- 凭据（临时件只建不删，全在导出树旁不入库）：`01-gates.log`（gofumpt/vet/build）· `02-winsec-full.log`（整包）· `03-d22scan.log` · `04-instrument-r1.log`（仪器首发逐发）· `05-instrument-count3.log`（`-count=3` 旁证）· `06-toothA-recheck-pulled.log`（牙 A）· `07-toothB-removal-pulled.log`（牙 B）。

## 一、AC#2 确定性仪器（交付形状）

**新件**：`internal/winsec/leftstanding_354_windows_test.go`，两枚用例：

1. `TestAC2WindowsAfterTheRefusalTheCallersNameStandsCarryingTheForeignObject354r1` —— 分子 **`leftStandingForeignName`**（票面 AC#2 建议名，逐字采用；⛔ 未借 `entered`／`damage`／`deleted`／`truncated` 任何一枚既有分子作判据——`entered` 只按本族形状作活性配对）。
2. `TestAC2WindowsRemovalStillRunsWhenNothingTopsUpTheName354r1` —— 负控（无 top-up ⇒ 名字必须 absent＝移除照跑），双向非钝的另一半，也是牙 B 的受腿。

**确定性＝按构造，不按复读**（派单点名的三条逐条交）：

- **固定前置态**：世界＝`newWinDetSwapWorld342r2`／`installWindow` **逐字复用**（不新造接缝、不加产码名）——真封印先跑、保其判词，再抽名、再把无名字节改名顶入、最后注入拒答；无 racer、无时序、无 sleep。外来体 body＝**固定 64 字节常量** `leftStandingForeignBody354r1`，用例开头**断言其长度＝64**（改常量即红仪器，读数不漂）。
- **轮内序号**（`353-v1` 具名的尺要求）：每发两枚夹具名各带发号（`left-standing-354r1-shot<N>.bin`／`the-unnamed-object-354r1-shot<N>.bin`），每条读数行带同一发号 ⇒ 原始日志逐行可归发。
- **读数逐字节可预期**：幸存名必须读回**逐字** `size 64 body [an object the caller never named, standing in the same directory]` 且顶入对象自身名 `absent`（⇒ 调用者名＝外来对象**唯一名**）才计入分子；其余任一不符＝逐发红。
- **旁证**（非判据）：`-count=3` 三轮 × 3 发＝**9 行读数逐字节全同**（`05` 号日志 grep 计数＝9）。

**首发读数**（`04`，逐字）：

```
354r1 AC#2 shot 1..3 LEFT-STANDING : the caller's name=present (size 64 body [an object the caller never named, standing in the same directory]) | the swapped-in object's own name=absent (...)   ← ×3，逐字节同
354r1 AC#2 over 3 shots: entered=3/3, leftStandingForeignName=3 (the caller's name standing after the refusal, carrying the foreign object as its ONLY name)
--- PASS: TestAC2WindowsAfterTheRefusalTheCallersNameStandsCarryingTheForeignObject354r1
354r1 AC#2 CONTROL ...: isNotSealable=true | leg=seal refusal on this route (...) | name=absent (...)
--- PASS: TestAC2WindowsRemovalStillRunsWhenNothingTopsUpTheName354r1
```

GREEN SCOPE 逐字在场（照抄 353 不冒认）：复验之后、`Deleteat` 的 `NtOpenFile` 绑名之前落地的那一窗在平台内部，本包无钩，本仪器不覆盖也不许被读作覆盖。

## 二、AC#3 成对去牙（只在导出树，两牙两落点，⛔ 单发当判据——全部三发整跑）

**牙 A ＝ 摘复验那步**（落点：`pinnedleaf_342_windows.go` `unlinkPinnedLeaf` 的 `Lstat`⇄`SameFile` 块，原 ：277-283 整块删；`Remove` 原样）：

- 本仪器读数 **3/3 → 0/3，断言红**。红句逐字（每发一条＋汇总各一）：
  - `354r1 AC#2 RED: shot 1 the caller's name did NOT survive the refusal (reads absent: ...) - the identity re-check in unlinkPinnedLeaf is the leg that leaves it standing, and on these bytes it no longer does; leftStandingForeignName cannot be counted from this shot`（shot 2/3 同形）
  - `354r1 AC#2 RED: leftStandingForeignName=0 of 3 entered worlds - the ratified residual moved. Form 乙's price is that the caller's name STAYS under the refused write and stays carrying the object the caller never named; a landing that changed either half changed what stands under a caller's name without saying so`
- **353r1 主案同步回到 `damage>0`**（票面 AC#3 判据逐字命中）：`unnamed objects left with no name (DAMAGE)=3`，红句逐字：`353r1 AC#2 RED: 3 of 3 entered worlds left the swapped-in object with NO name, so the refusal-path removal ran on a name that carried an object the caller never named - ticket 353's damage, on the platform where the route's own handle makes that name vacate-able (bench row S4)`。
- 两枚负控照旧 PASS（牙 A 不碰无 top-up 路径＝方向正确）。日志＝`06`。

**牙 B ＝ 摘删除那步**（落点：同函数 `Remove` 行，原 ：284 单行删；复验块原样——**落点⛔同**成立）：

- **两枚负控齐红**（`354r1` 自带负控＋353r1 负控），逐字：`354r1 AC#2 RED: with NOTHING topped up the refused write still left the name standing (present: size 0 body []), so the refusal-path removal no longer runs - the numerator above would then be counting a removal that cannot happen, and AC#5's refused write now leaves a half-made file`（353r1 负控同形红）。
- ★**主仪器照绿（`leftStandingForeignName=3/3`），如实报、非缺陷**：该分子量的是"复验拒后名字下站着谁"——摘掉删除那步不改变复验留下的形，所以看见牙 B 的是负控、不是分子。两牙各红各的腿＝成对判据成立。日志＝`07`。

**两牙还原**：牙 B 后从 `.bak` 拷回原件，导出树 `git diff --stat` 复核 `internal/winsec/pinnedleaf_342_windows.go` 回到 HEAD（`07` 之前 `cp` 恢复已做；**落 commit 的字节零去牙**——主树从未载入任何去牙形）。

## 三、AC#4 门禁读数（导出树 `393430dc9`，跑前 `unset GOTMPDIR TMPDIR TEMP TMP`＋`PATH=$PWD/build:$PWD/third_party/sherpa-onnx:$PATH`，CGO_ENABLED 未设）

| 门禁 | 读数 |
|---|---|
| gofumpt（`internal/winsec/` 全包） | 空（rc0） |
| `go vet ./internal/winsec` | rc0 |
| `go build ./...` | rc0 |
| `go test ./internal/winsec -count=1` | `ok github.com/CarlosShao/wisp/internal/winsec 37.254s` rc0 |
| d22scan（导出树现编 `tools/d22scan/d22scan.exe`） | `clean - no D22 ban violations` rc0 |

AC#4 其余半格（两平台成对 unix 半边＋逐名红册）按票面仍归收口腿裁；unix 机制侧「复验拒删留下名字」不存在（`pinnedleaf_342_unix.go` 无复验，`211-354-closure.md` 同口径），故本仪器为 `//go:build windows` 单边＝与既有两族同形。

## 四、AC#5 接缝（那一跳在哪）与 AC#6 纪律

- **接缝**：被量的那一跳＝`pinnedleaf_342_windows.go` `unlinkPinnedLeaf`：复验 `p.root.Lstat(base)`⇄`os.SameFile(pinned, now)` 不匹配 ⇒ 在 `_ = p.root.Remove(base)`（:284）**之前 return**（:277-283）⇒ 调用者名字下留着攻击者对象。仪器驱动的世界经包内**已声明接缝** `applyDescriptor`（`winsec_windows.go:35-38`）进入该路径，零新接缝。与 342/348 的既有红⇄摆动名对接：`nametopup_353` 的 `refusals that left the caller's name standing`（计数型、不核体）与 `ancestorwindow_342r2` 的 GUARD HELD 行照旧在场；本票不重述 118/120-r3。
- **纪律**：写面三处＝新 test 件＋本件＋票尾 Progress 一行；⛔ 翻票面 AC 勾／Status；⛔ push；改动 `.go` 件 LF 现核（`grep -c $'\r'`＝0）；两笔 commit 显式 pathspec；`scripts/`、`cmd/wisp/**`、`frontend/**`、`PLAN.md`、thresholds/golden 零字节。
- **撤销口令**：主树 `git revert` 两笔（`393430dc9` 仪器 commit＋证据/票尾 commit）即回原状；导出树与 `/d/tmp/wisp/20261108/r354/` 全部为临时件、不入库。
