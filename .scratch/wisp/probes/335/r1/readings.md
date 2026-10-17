# 票 335 落地腿 335-r1·起手复跑与改形读数（⛔ 抄编排者的数，全部现跑）

锚＝具名提交 `a06a73d8`（＝`859026d2` 之后的 HEAD，`git log -1` 现认，2026-10-17 11:44 +0800）。
台面＝仓外导出树 `D:\V335R1\T1`（配方＝八件套 `go.mod go.sum internal cmd tools scripts frontend models`
＋仓根 `docs/`＋手工三件 `build/*.dll`、`third_party/sherpa-onnx/*.dll`、`frontend/dist`；
D: 卷、目录名 `V335R1` 产生 8.3 短名；取出后先数 `.go`＝**757** 证非空）。
本腿⛔ 在共享工作树跑过 `./cmd/wisp`、⛔ `go test ./...`、⛔ 动 PID 21892。

## 1. AC#0 起手四行逐枚自己重跑（枚枚带 `rc=`）

| 表里读数 | 我这把尺的读数 | `rc=` | 件 |
|---|---|---|---|
| 那枚钉子红、红句 14 行＝`:448`＋`:451`×7 枚，`U+206B` 不在串里 | **逐字复认**：`--- FAIL: TestTheRosterIsClosedAgainstTheRunesTheRulingLeftOut`，红句 14 行、行号 `display_scrub_216_test.go:448`/`:451`、七枚＝`U+2061 U+2062 U+2063 U+2064 U+061C U+180E U+206A`，`U+206B` 0 命中 | rc=1 | `logs/ac0-nail-anchor.log` |
| 起手名册 5 ⇄ 交件 6，gained 只有那一枚、lost 0 | **与表里不同＝台面差别，具名更正**：我这把尺的导出树没有 `.git`，同一条命令读到起手 **8** 枚＝表的 5 枚＋那一枚钉子＋两枚台面红（`TestComposerRenderFixtureTellsTheTruth`、`TestReadGitOnThisRepositoryIsSelfConsistent`，红因＝读盘的 git 面在导出树里不存在）。带 `.git` 那一发（共享工作树，终态）读到 **5** 枚＝与表里"起手 5"逐字同名，见 `logs/final-shared-tree.log` | rc=1 | `logs/pre-T1-c.log`、`logs/final-shared-tree.log` |
| `post-c: RUN=257 PASS_TOP=161 FAIL_TOP=6 SKIP=0` ⇄ `post-d` 同四数 | **⛔ 抄**：那是腿报、编排者未复跑的数，且发面＝共享工作树。我自己的成对发（同一条命令只挪临时根）＝改前 `257/159/8/0`⇄`257/159/8/0`、改后 `257/160/7/0`⇄`257/160/7/0`；枚数差 2 全部由上一行的 `.git` 台面差解释，同一台面内改前⇄改后只差那一枚 | rc=1（四发皆包红，红名册已知） | `logs/pre-T1-*.log`、`logs/post-T3-*.log` |
| 三处过期文字 `:90`／`:92`／`:193` | **按内容锚现认，锚 blob 行号与表里一致**：① `//	format rune above U+00FF (the ticket 327 Cf roster) -> \uXXXX, lowercase hex`；② `// The last shape is the one the invisible roster lands in: U+202E prints as`；③ `// C1 U+007F-U+009F), the invisible formatting roster in isDisplayWashCf, plus any`。三枚都是整行注释，行号只当读数 | rc=0 | `git show a06a73d8:internal/panel/display_scrub_216.go` |

★再加的那把尺（那 7 枚今天是否都在 `Cf` 范畴内）＝仓外纯 stdlib 探针 `D:\V335R1\cfcheck\main.go`：
`GO_UNICODE_VERSION=17.0.0`、`CF_COUNT=170`、正控 `CC_COUNT=65`（对上 333-r1 的数）；
七枚逐枚 `cf=true zs=false`，`U+206B` 同形 `cf=true`，旧十五枚抽查 5 枚全 `cf=true`，
反控 `U+0041` cf/cc 皆 false ⇒ **7/7 成立且⛔ 一枚不多**（rc=0，`logs/ac0-cf-category-probe.log`）。
arm 3 的十枚样本另跑一把前提尺（`D:\V335R1\samples\main.go`）＝逐枚 cf/cc 皆 false，rc=0（`logs/ac0-arm3-sample-premise.log`）。

## 2. AC#1 改形＝把"这 7 枚⛔ 洗"翻成"这 7 枚必须被洗"，同时保住那道闸

旧断言（逐字两句，`U+2061` 为例）＝
`U+2061 is not on the ruled list but needsDisplayWash fired on it` ⇄
`U+2061 is not on the ruled list but the wash changed it: "photo\u2061.png" -> "photo\\u2061.png"`。
红的成因⛔ 在产码（`isDisplayWashCf` 按范畴判定是对的，票 333 裁甲-全＝台账 `A1015`，撤销口令「333 改裁乙」），
成因在这段文字描述的是被那一纸撤销的旧形状 ⇒ 处方＝改形（⛔ 删用例、⛔ `t.Skip`、⛔ 留恒红），方向**更严**。

用例名 `TestTheRosterIsClosedAgainstTheRunesTheRulingLeftOut` 一字未改（红名册作差要靠它对齐），函数体＝三臂：

- **arm 1 名册已翻**：逐枚那 7 枚（`U+2061 U+2062 U+2063 U+2064 U+061C U+180E U+206A`）现在**必须**被洗——
  范畴前提（`unicode.Is(unicode.Cf, r)`，不成立就报"这臂失去存在理由"）、快路径必须响、洗出的形＝转义⛔ 删除、
  再洗一次必须是恒等。命名只是**读数**（那批是 327 名册漏的七枚，漏掉的第八枚 `U+206B` 由 327-r2 的 nail 2 逐枚点名），⛔ 是洗集的定义。
- **arm 2 闸的正身＝"洗谁"⛔ 再由手抄清单回答**：用本包既有的 `washSetOver`（`display_scrub_327r2_test.go` 的助手）
  把**产品谓词**与 `unicode.Cf` 各自在全标量上收集成集，三句断言＝
  ① 手抄名册复原检测器 `len(product) <= len(cfRoster)` 必须不成立（本文件那份 15 枚 `cfRoster` 是盘上真存在的清单，
  作"最多能手抄多少枚"的地板）；② 范畴成员⛔ 有一枚没被洗（missed 双向）；③ 范畴外的⛔ 多洗一枚（extra 双向）。
  空集检测器：两集任一为 0 直接 `t.Fatalf`（一条没扫过的臂⛔ 算过）。
- **arm 3 反方向**：十枚既⛔ 在 `Cf` 也⛔ 在 `Cc` 的样本（`U+0020 U+00A0 U+2028 U+2029 U+0300 U+2192 U+0061 U+4E2D U+1F600 U+FE0F`）
  必须逐字节原样回来，快路径不许响；样本前提自检测器（任一枚落进被洗范畴 ⇒ 报"这臂失去前提"）⇒ 后来者⛔ 悄悄往清单里加一枚。

改形后的绿（交件侧，同一条尺、带 `file:line`）：
`--- PASS: TestTheRosterIsClosedAgainstTheRunesTheRulingLeftOut (0.03s)`，
`display_scrub_216_test.go:519: arm 2: product-washed=170, unicode.Cf=170, missed=0, extra=0, this file's cfRoster=15`，rc=0。
起手侧红 ⇄ 交件侧绿都出自**同一枚 `-run` 名册**（`TestTheRosterIsClosedAgainstTheRunesTheRulingLeftOut`），
差别只有本腿那一笔测试件改动（台面对齐见 `logs/ac1-reshaped-green.log`＝T3 里同一命令、`rc=0`，改前那发＝`logs/ac0-nail-anchor.log`、`rc=1`）。

⛔ 放宽任何断言、⛔ 为变绿改产码、⛔ 动 SLO/golden/thresholds：本腿产码字节＝0（AC#2 那把尺逐字节自证，见 `delivery.md` §1）。
