# 票 289 / 腿 289-r1 — AC#3 改前/改后三数并排＋红名册逐名作差

harness（两发同一支；下面这一行就是实际执行的那一行，逐字）：

```
cd "D:\work\workspace\projects plans\Wisp" && \
PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test -v ./internal/panel/ ./cmd/wisp/ -count=1
```

- CWD＝仓根工作树（`D:\work\workspace\projects plans\Wisp`），⛔ 不是包目录。
- PATH 前缀逐字照派单给的 harness；⚠ 具名一处偏离：派单写"harness 逐字用 `... go test ./internal/panel/ ./cmd/wisp/ -count=1`"，
  本腿在它前面加了 `-v`，**为的是 AC#3 要的三数**（非 `-v` 时通过的包只出 `ok`，拿不到 PASS/SKIP 名册，也就拿不到逐名红名册）。
  `-v` 只改输出冗长度，不改被验物集合、不改 `-count=1`、不改包列表、不改 PATH。
- 环境红排查：两发原始输出里 `grep -c "0xc0000135"` 均＝**0** ⇒ `cmd/wisp` 的用例真的跑了
  （包级耗时 460.508s／458.876s 即证；⛔ 缺 DLL 时会是 `exit status 0xc0000135` 且零 `--- FAIL`＝用例根本没跑）。
- 原始输出：改前＝本目录 `raw-pre.md`（2427 行），改后＝本目录 `raw-post.md`。
- 基线那一段工作树：改前两枚件均＝HEAD 内容（`git diff --stat HEAD -- <4 枚件>` 空）；
  别人在飞的改动（`design/**` 若干删除、`.gitignore` 修改、别人的 probes/票面件）两发里都原样存在，本腿未还原、未提交。
- 归因尺（本腿自加，因为这是共享工作树、HEAD 在两发之间被别人的 commit 推进过）：
  `git diff --name-only fd692f7e..HEAD | grep "\.go$"` ＝ **0 行**
  ⇒ 两发编译的是同一批产码；两发之差只可能来自本腿那 4 枚件，而那 4 枚件的非注释行＝0（`10-four-comments.md` 尺 B）。

## 三数并排（顶层 `^--- X:`；子测试另记一档）

| 包 | 改前 PASS | 改前 FAIL | 改前 SKIP | 改后 PASS | 改后 FAIL | 改后 SKIP |
|---|---|---|---|---|---|---|
| `internal/panel` | 124 | 6 | 0 | 124 | 6 | 0 |
| `cmd/wisp` | 265 | 6 | 2 | 266 | **5** | 2 |
| 合计（顶层） | **389** | **12** | **2** | **390** | **11** | **2** |
| 合计（含子测试的 PASS） | 573 | 12 | 2 | 574 | 11 | 2 |

包级行与整发 rc（由被量命令自己的 `$?` 落，⛔ 不是管道尾元素的退码）：

- 改前：`FAIL github.com/CarlosShao/wisp/internal/panel 7.019s` ＋ `FAIL github.com/CarlosShao/wisp/cmd/wisp 460.508s`，**rc=1**
- 改后：`FAIL github.com/CarlosShao/wisp/internal/panel 7.156s` ＋ `FAIL github.com/CarlosShao/wisp/cmd/wisp 458.876s`，**rc=1**

## 红名册逐名作差（尺＝名册逐名，不比计数）

| # | 用例名 | 改前 | 改后 |
|---|---|---|---|
| 1 | `TestApprovalCardViewJSONKeysMatchFrontendTypes` | FAIL | FAIL |
| 2 | `TestStreamLogFloodBelowKeyBoundIsNotBounded35r8` | FAIL | FAIL |
| 3 | `TestStreamLogDroppedNamingLedgerIsNotBounded35r8` | FAIL | FAIL |
| 4 | `TestComposerContractTypesMatchFrontend` | FAIL | FAIL |
| 5 | `TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme` | FAIL | FAIL |
| 6 | `TestC21DesignTokensFourWayAgree` | FAIL | FAIL |
| 7 | `TestTicket223ModeLooseningChangesTheRunningModeAfterAllow` | FAIL | **PASS**（见下一节） |
| 8 | `TestPanelHostRealWindowHopAndLifecycle` | FAIL | FAIL |
| 9 | `TestAC4FocusReturnToPriorWindowGap33r5` | FAIL | FAIL |
| 10 | `TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe` | FAIL | FAIL |
| 11 | `TestAC14AwaitedBindingReplyReachesThePage` | FAIL | FAIL |
| 12 | `TestAC14GoSideEvalPushReachesThePage` | FAIL | FAIL |

- **AC#3 的判据＝新增红 0 枚：成立** ⇒ 改后名册（11 枚）⊖ 改前名册（12 枚）＝ ∅。
- 反向作差不空：改前 ⊖ 改后 ＝ {`TestTicket223ModeLooseningChangesTheRunningModeAfterAllow`}＝**1 枚红转绿**。
  ⛔ 本腿没有为了变绿放宽过任何一枚断言（非注释行＝0 就是这一句的尺），也没有把这一枚读成"本腿修的"——
  下面单独量它。

## 那一枚红转绿是谁干的：具名归因（结论＝不是本腿，是整包串跑里的时序）

改前那一发的失败原文（`raw-pre.md:751` 往上截，连表头一起截）：

```
    config_reload_223_test.go:545: the card does not name risk.permission_mode: wisp run: 答复监听已接入（卡片上给了编号）。…
time=2026-10-09T14:09:04.234+08:00 level=WARN msg="config: locked loosening approved via L2 re-confirmation (D36 rule 1)" section=risk keys=[risk.permission_mode]
--- FAIL: TestTicket223ModeLooseningChangesTheRunningModeAfterAllow (2.18s)
```

- 失败点是 `cmd/wisp/config_reload_223_test.go:545`——**枚文件不在本腿写面里**（本腿只碰那 4 枚件），
  形状是"卡片还没写上 `risk.permission_mode` 就去读了"，紧接着一行日志里那次 L2 复核才落地 ⇒ 读比写早，典型的串跑时序。
- 靶向复跑（同一 harness 的 PATH，只加 `-run`，逐发落 rc）：

```
for i in 1 2 3; do PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" \
  go test ./cmd/wisp/ -count=1 -run '^TestTicket223ModeLooseningChangesTheRunningModeAfterAllow$' -v; done
run1 rc=0 ⇒ --- PASS: TestTicket223ModeLooseningChangesTheRunningModeAfterAllow (2.06s)  ok  github.com/CarlosShao/wisp/cmd/wisp  2.121s
run2 rc=0 ⇒ --- PASS: TestTicket223ModeLooseningChangesTheRunningModeAfterAllow (2.02s)  ok  github.com/CarlosShao/wisp/cmd/wisp  2.089s
run3 rc=0 ⇒ --- PASS: TestTicket223ModeLooseningChangesTheRunningModeAfterAllow (2.19s)  ok  github.com/CarlosShao/wisp/cmd/wisp  2.264s
```

- 读数摆开：单独跑 3/3 绿；改前整包红；改后整包绿。**同一枚用例在"整包串跑"里两个方向都出现过，
  而两包之间本腿唯一的变化是 4 句注释（非注释行＝0、两发之间零 `.go` 文件入仓）⇒ 它是与同包其它用例串跑相关的红，
  不由本腿造成，也不能记成本腿修的。** ⛔ 本腿不据它判绿、不据它翻 AC 框，只把它作为"改前基线不是绿的、
  且那 12 枚里至少这一枚还会自己变"的具名事实交裁决者。

## 与票面现量的两处不符（具名报回，⛔ 不改票面一字）

1. 票 289 `AC#3` 写 `internal/panel` "现量 **5 枚具名红**"并列了 5 个名字；本腿两发现跑都是 **6 枚**，
   多出的一枚＝`TestApprovalCardViewJSONKeysMatchFrontendTypes`（与票面已点名的 `TestComposerContractTypesMatchFrontend`
   同属前端契约那一族）。这一处编排者已在 `712f8d1a`（台账 A779）自认"票 289 AC#3 的'5 枚'偏低(真值 6)"，
   并写明"此刻⛔ 不改票面——等 289-r1 交回再一次性追加更正"、"它的判据是逐名作差＝新增红 0 枚，不受我这枚错数影响"。
   ⇒ 本腿按现跑数记账，与那条更正同向。
2. `cmd/wisp` 那 6 枚（含上面会变的那枚）票面**完全没列**。本腿照样逐名入册——作差的尺是"名册逐名"，
   不是"票面列了几枚"；否则那一族红就没有任何见证力了。
3. ⚠ 另有一处**派单与票面都落空的前置**：派单要求"在票面 `## Progress log` 末尾追加一条"，
   但票 289 的票面**没有 `## Progress log` 这一节**（现跑 `grep -n "^## " <票面>` ＝ `## 现量`／`## 要建什么`／`## 禁区` 三节，
   同族先例票 288 也只有那三节）。按本仓"未定义即停"与 `A780` 里对票 291 同款缺节的裁法
   （"腿按未定义即停没自建、不判越格"），本腿**没有自造那一节**，把这一处具名交回编排者。
