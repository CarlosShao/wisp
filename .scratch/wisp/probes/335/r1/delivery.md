# 票 335 落地腿 335-r1·交件自证（三把尺＋两题必答）

写面＝`internal/panel/**` 的测试面＋整行注释面、`.scratch/wisp/probes/335/r1/**`、票尾一条 bullet。
⛔ push、⛔ 翻框、⛔ 改票名／`Status:`／`-done`、⛔ `add -A`／`add .`／目录形、⛔ amend／reset／rebase／stash／checkout .／clean、⛔ 仓内 worktree。

## 1. 尺① 逐枚提交号＋逐枚 `git show --name-only`（真跑，⛔ 自报）

| 提交 | 内容 | 文件（`--name-only` 全集） | `--numstat` |
|---|---|---|---|
| `9567b094` | AC#1 钉子改形 | `internal/panel/display_scrub_216_test.go` 一枚 | 97 增 / 5 删 |
| `77d98354` | AC#2 三处过期注释 | `internal/panel/display_scrub_216.go` 一枚 | 3 增 / 3 删 |
| 本笔 | 证据件＋两把脚本＋`logs/` | `.scratch/wisp/probes/335/r1/**` | 见 `git show --numstat` |
| 末笔 | 票尾一条 bullet | `.scratch/wisp/issues/335-...md` | 1 增 / 0 删 |

**AC#2 那笔＝非注释字节零改动**（两条独立尺）：
`git diff -U0 77d98354^ 77d98354 -- internal/panel/display_scrub_216.go | grep '^[+-][^+-]'`＝6 行（3 删 3 加），
其中 `grep -vc '^[+-][[:space:]]*//'`＝**0**（实体 0 枚非整行注释行）；
整行注释滤掉后的逐字节尺（照票 331 的 `zero-change-check.sh` 形状，本腿件＝`.scratch/wisp/probes/335/r1/zc335r1.sh`）：
`anchor-stripped-sha256=cand-stripped-sha256=5274713dead6a3f5b4b2c06970697bbaff50eff4c9dda2a4e33410cb7178b305`、`diff-nonempty=no`、`diff-rc=0`；
**正控**＝副本里把 `isDisplayWashCf` 的 `Cf` 换成 `Cc`（一枚非注释字节）⇒ 同一把尺 `same_sha=no`、`diff-rc=1`、`bdb58d7b…`；
**负控**＝锚副本里只动一枚整行注释 ⇒ 同一把尺 `rc=0`（本尺按构造对整行注释失明，射程如实写在脚本头）。逐字输出＝`logs/ac2-zero-change-ruler.log`。

禁区自证（AC#5）：四笔 `--name-only` 并集⛔ 含 `internal/risk/**`（`334-r1` 在飞，本腿读代码一律取具名提交树 `a06a73d8`／`77d98354`）、
⛔ 含 `cmd/**`、`docs/specs/**`、`PLAN.md`、台账、`HANDOVER`、`frontend/**`、`design/**`、三枚冻结件、`thresholds.go`、golden、
`allowlist.txt`、`go.mod`／`go.sum`／`.gitignore`／`ci.yml`；产码⛔ 一字未动（判定层 `internal/risk/taintmatch.go` 的 `isIgnorableTaintRune`
⛔ 拖进来——`R-333-1` 明确推迟，两轴取向相反，本腿只在注释里具名它的存在）；洗集范围⛔ 借本票改（仍＝票 333 裁甲-全的 `Cf` 范畴）。

## 2. 尺② 证据件（绝对路径＋`wc -c`＋行数＋空槽尺 0 命中）

件名册与逐件字节数、行数、空槽尺命中数见 `logs/placeholders.log`（尺的精确命令行记在该日志第 1 行，
射程＝两枚 `.md`＋两枚脚本＋两枚 Go 件，逐件 0 命中；`logs/` 原始体⛔ 扫、因尺句会自指）。⛔ 0 字节件。

- `D:\work\workspace\projects plans\Wisp\.scratch\wisp\probes\335\r1\readings.md`
- `D:\work\workspace\projects plans\Wisp\.scratch\wisp\probes\335\r1\delivery.md`（本件）
- `D:\work\workspace\projects plans\Wisp\.scratch\wisp\probes\335\r1\zc335r1.sh`（AC#2 零改动尺）
- `D:\work\workspace\projects plans\Wisp\.scratch\wisp\probes\335\r1\mut335r1.sh`（AC#1 去牙脚本，trap 还原）
- `D:\work\workspace\projects plans\Wisp\.scratch\wisp\probes\335\r1\logs\`（19 枚原始体＝`ls logs | wc -l` 现数：起手红、两把探针、四发颜色、三发突变、门禁、名册作差、空槽尺）

## 3. 尺③ 决定性那把尺逐字（起手侧 ⇄ 交件侧，都带 `file:line`）

命令（同一枚 `-run` 名册、同一台面配方，只换那一笔字节）＝
`go test -v ./internal/panel/ -run 'TestTheRosterIsClosedAgainstTheRunesTheRulingLeftOut' -count=1 -timeout 30m`

起手侧（`D:\V335R1\T1`＝`git archive a06a73d8`，757 枚 `.go`，`rc=1`，红句 14 行＝`:448`＋`:451`×7 枚）：

    display_scrub_216_test.go:448: U+2061 is not on the ruled list but needsDisplayWash fired on it
    display_scrub_216_test.go:451: U+2061 is not on the ruled list but the wash changed it: "photo\u2061.png" -> "photo\\u2061.png"
    --- FAIL: TestTheRosterIsClosedAgainstTheRunesTheRulingLeftOut (0.00s)

交件侧（`D:\V335R1\T3`＝同一锚树＋本腿两枚件，`rc=0`）：

    display_scrub_216_test.go:519: arm 2: product-washed=170, unicode.Cf=170, missed=0, extra=0, this file's cfRoster=15
    --- PASS: TestTheRosterIsClosedAgainstTheRunesTheRulingLeftOut (0.03s)

成对去牙（两发落点不同、逐发当场还原、`logs/ac1-pair-summary.log` 逐字全体）：

- 发ⓐ＝`isDisplayWashCf` 体内那枚 `return unicode.Is(unicode.Cf, r)` 换回旧的 15 枚手抄名册：
  `MUT_HASH=a2d4c113…`、11 行变动、`TEST_RC=1`、本钉子红 16 句、命中断言行 `:480 :484 `**:505** `:522`，逐字两句＝
  `display_scrub_216_test.go:480: U+2061 is a member of the category the ruling chose but needsDisplayWash stayed silent on it`
  `display_scrub_216_test.go:505: the boundary washes 15 runes, no more than the 15 this file could hand-copy: the wash set is a hand-kept roster again, …`
  （`arm 2` 读数＝`product-washed=15, unicode.Cf=170, missed=155, extra=0`）
- 发ⓑ＝在范畴判定里为 `U+2061` 开特例排除，落点＝函数签名行之后（⛔ 发ⓐ那一枚 token，不同内容锚、不同行）：
  `MUT_HASH=2feeb59f…`、1 行变动、`TEST_RC=1`、本钉子红 3 句、命中断言行 `:480 :484 :522`，**⛔ 命不中 `:505`**，逐字＝
  `display_scrub_216_test.go:522: a rune the Cf category names is not washed, so the answer came from a list rather than the rule: U+2061`
  （`arm 2` 读数＝`product-washed=169, unicode.Cf=170, missed=1, extra=0`）
- 发ⓒ（加分，另一向）＝范畴判定里为 `U+00A0`（范畴⛔ 点名）开特例**纳入**：`MUT_HASH=aa53cde1…`、`TEST_RC=1`、
  命中 `:526 :540 :543`＝与发ⓐ发ⓑ**交集为空**（`comm` 尺在件里），证明臂 2/臂 3 对"清单被悄悄加长"也响。
- 还原哈希三发一致：`PRE_HASH=RESTORE_HASH=5c68a930ae33e70b92a8ff4fefd107c0c6281f2264df16bcbc13d5026813157a`
  ＝`git show 77d98354:internal/panel/display_scrub_216.go` 逐字同串（脚本 `trap` 还原＋跑完再验一次，见 §5）。

两把临时根台面成对（同一条命令、同一枚名册、四数逐枚点名、⛔ 混用⛔ 相减；`logs/colours-and-roster.log`）：

    pre-T1-c  RUN=257 PASS_TOP=159 FAIL_TOP=8 SKIP=0     pre-T1-d  RUN=257 PASS_TOP=159 FAIL_TOP=8 SKIP=0
    post-T3-c RUN=257 PASS_TOP=160 FAIL_TOP=7 SKIP=0     post-T3-d RUN=257 PASS_TOP=160 FAIL_TOP=7 SKIP=0

红名册作差＝**gained 0、lost 1**（只有那枚钉子），逐名列出在 `logs/colours-and-roster.log`；
带 `.git` 的那一发（共享工作树、终态、`internal/panel` 两枚件与 `77d98354` 逐字节同）＝`RUN=257 PASS_TOP=162 FAIL_TOP=5 SKIP=0`，
余 5 枚与票的起手 5 逐名同＝AC#3 闭合。门禁四数只在共享工作树（`logs/gates-four.log`）：
`gofumpt -l`（两枚件，真身 `/d/work/base/gopath/bin/gofumpt.exe`）rc=0 空列表／`sh scripts/d22scan.sh` rc=0（`clean - no D22 ban violations`，
ban #8 internal/=584、cmd/=143、design/=39、frontend/=116）／`sh scripts/check-path-length-budget.sh` rc=0（`VERDICT GREEN`，
tracked=14820 over-budget=60 covered=60 not-in-roster=0）／`go vet ./internal/panel/` rc=0。

## 4. 必答两题

**(i) 改形后的钉子在"有人把某一枚 `Cf` 特例排除"时会红吗？我实测了哪一发？**
会。实测＝**发ⓑ**（`U+2061` 在 `isDisplayWashCf` 里被特例排除，落点＝签名行后插一行，`MUT_HASH=2feeb59f…`）：
本钉子 `--- FAIL`、3 句红，逐字断言行 `:480`（快路径沉默）＋`:484`（洗后原样）＋`:522`（`missed=1`，点名 `U+2061`）。
另配发ⓒ证明反方向也响（特例**纳入**一枚范畴外的字符⇒ `:526 :540 :543`）。如实具名一处形状：发ⓑ的红名册与发ⓐ**同名**、
断言行是发ⓐ的子集——因发ⓐ摘掉 155 枚、发ⓑ只摘 1 枚，任何发ⓑ能失败的断言发ⓐ也失败；把两发分开的凭据是**发ⓐ独有的 `:505`**
（手抄名册检测器），完全对称的差由加分发ⓒ承担（三行与ⓐⓑ全不相交）。

**(ii) 终态红名册那 5 枚起手红，逐枚归谁**（⛔ 当本腿的债、⛔ 当"全绿"）：
1. `TestStreamLogFloodBelowKeyBoundIsNotBounded35r8`＝票 35 **有意保留**的红（`backpressure_bound_35r8_test.go`，`b5439299` 交件即红，台账 `A922`/`:13952` 族）。
2. `TestStreamLogDroppedNamingLedgerIsNotBounded35r8`＝同上，票 35 同一格有意留的红灯。
3. `TestApprovalCardViewJSONKeysMatchFrontendTypes`＝页面契约缺字段老账（票 145／180 那一族，Go 侧发的键页面 interface 未声明）；`frontend/**` 归机主外部 agent，本编队零写面。
4. `TestComposerContractTypesMatchFrontend`＝同族第二枚（`composer.tsx` 那一面），同上归属。
5. `TestC21DesignTokensFourWayAgree`＝`design/assets/tokens.css` 读不到：共享树里红因＝工作树被外部在飞删除（台账 `A773` 那 16 枚 ` D`），
   我这把导出树尺里红因＝八件套⛔ 含 `design/`＝台面射程；两形都⛔ 属本腿判据，`design/**` 亦归机主外部 agent。
另有两枚只在我的导出树红（`TestComposerRenderFixtureTellsTheTruth`、`TestReadGitOnThisRepositoryIsSelfConsistent`）＝
archive 树没有 `.git` 的台面红，带 `.git` 那一发逐枚转绿，已具名在上，⛔ 计入判据。

## 5. 未定义即停／事故自报

本腿⛔ 撞票面未覆盖的形状（三处注释都按内容锚找得到、7 枚全在范畴内、改形两发＋一发加分全响）。
一条如实上报的形状：AC#1 成对两发在这枚谓词上⛔ 能做出"对称差"（发ⓐ的错集真包含发ⓑ的错集），
我用"发ⓐ独有 `:505`"＋"加分发ⓒ三行不相交"两条补全凭据——若编排者认为这一形仍⛔ 够，请另派程，本腿⛔ 自行加删用例。
临时台件（`D:\V335R1\T1|T2|T3`、`tmp`、`gotmp`、`cfcheck`、`samples`、`mut`、`pc`）只建不删；
`T2`＝`git archive 77d98354` 的落地树（含别腿已入库字节，`.go`＝758）留作后续对拉，本腿颜色用的是 `T1`⇄`T3` 那一对。
