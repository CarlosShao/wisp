# 74 — Make the C21 token table mirror the constants that actually drive the ball

**Status:** in-progress (agent ticket74 implementing; family 2 first, per the ticket's own priority)
**Type:** defect-fix (contract fidelity, cosmetic surface but governance-loaded: C21 is a frozen contract)
**Blocks:** nothing · **Blocked by:** nothing — ticket 69 shipped the machine check that discovered this
**Packages:** `internal/ball` + `docs/evidence/s1/c21-native-tokens.md` only. Do **not** touch
`internal/tools`, `internal/observe`, `internal/secret`, `tools/d22scan`, `internal/llm` (other agents in flight).
**Spec refs:** C21 DesignTokens, registry **A24-D4**, ticket 69 report D7–D12

## Proven premise (measured by `4ca2c6c` / `9c29555`, do not re-derive)
`internal/ball/tokens_table_test.go` now cross-checks the C21 table against `tokens.go` both ways.
It emitted **29 named divergences** the check deliberately does not fail on. Three families:
1. **Code uses constants the table never mentions** — `hit.go` `RingMarginPx`/`ClickTolerancePx`,
   and **9 animation constants in `liquid.go`**.
2. **Semantic drift (the dangerous one)**: `tokens.go` declares `SwimLevelGain` / `SpinLevelGain`
   with **zero consumers**, while the real driver is `liquid.go`'s `Spin*RadPerS` group. So the table
   can document a number that is wired to nothing while the number that moves pixels on screen is
   undocumented. **This is the same shape as A33 ("声明✓ / 实测✗") and must be treated as the priority.**
3. **Table rows nobody consumes**: 18 `Palette` fields, `looks[*].Deep`, `Dur*Ms`, `CountdownFontPx`.

## What to build
Kill family 2 first, then family 1, then decide family 3 explicitly — do not silently do nothing
about any of them.
- **Family 2**: delete the zero-consumer `SwimLevelGain`/`SpinLevelGain` **in the same commit** that
  points the table at `liquid.go`'s real drivers (removing a documented token and adding its
  replacement in one step; never leave the table pointing at a dead constant). If you instead prove
  the dead constants are the intended future contract, write that proof in the commit message and
  make the *test* enforce it — do not leave both halves floating.
- **Family 1**: put the 11 unconsumed-but-live constants into the table with their real units.
- **Family 3**: for each of the 29→N remaining table-only rows, either (a) add the missing consumer
  assertion, or (b) list it in an explicit `documented-but-unconsumed` exemption table **with a
  one-line reason per row**. A blanket "not used yet" is not a reason.
- Tighten `tokens_table_test.go` only where you can name the exact bug class it newly catches: the
  geometry match is currently *containment*, not bijection (two rows can mask each other's number),
  and `FontFamily` has no value assertion. Fixing containment→bijection is in scope **only if you
  show a mutation that today passes and must fail**; otherwise record it as a gap.

## AC (1:1 verdict table required, one row per box)
- [x] **AC#1** `SwimLevelGain`/`SpinLevelGain`: no longer both "declared" and "unconsumed" — either
  deleted-with-replacement-in-one-commit, or the test pins them to a named consumer. Show the grep
  that proves the consumer exists (`git grep -n`).
- [x] **AC#2** Every constant that *changes pixels* under `internal/ball` (hit.go, liquid.go,
  renderer_windows.go) is either in the C21 table or in the exemption list with a reason; the
  machine check itself proves it (the report count goes to 0 unattributed).
- [x] **AC#3** Mutation proof for the new/changed assertions: pick at least 2 rows, change the code
  value, show the specific test go red by name, grep-prove the mutation landed, restore, grep-prove
  it is gone (`git diff --quiet` on both files).
- [x] **AC#4** Table↔`tokens.css` gap: either add the cross-check, or move this AC to a written
  hand-off (name the ticket). **"Deferred" without a ticket number is not an option** — the owner's
  standing rule is that postponed work must be registered with its completion criteria.
- [x] **AC#5** Gates: `gofmt -l internal/ball` empty, `go vet ./internal/ball/...` rc=0,
  `go test ./internal/ball/... -count=2`, and `go test ./internal/ball/ -v` with **RUN count == 2×
  distinct test names and zero `SKIP` occurrences** in the log (SKIP counts as ok — that trap has
  burned this project twice).

## Explicit prohibitions for this ticket
- No `git commit --amend`, `reset`, `rebase`, `stash`, `checkout .` (A34). Use
  `git commit -F - -- <explicit paths>` with a **quoted** heredoc; an unquoted heredoc let the shell
  execute backticks in a commit message on `9bf05cb` (A31) and also swept a foreign staged rename
  into that commit.
- **Do not edit `docs/PLAN.md` or `docs/specs/*`.** C21's *contract text* is frozen; this ticket only
  touches the evidence-layer table and code. If you believe the contract itself is wrong, stop and
  report it as a D22 item instead of editing.
- Do not delete a guard in one commit and add its replacement in the next (AC#1 family 2 is the one
  place a deletion is allowed, and only when paired).
- No `t.Skip`, no lowered thresholds, no rewriting golden values to make a check pass.
- If an AC cannot be proven, **leave the box empty and say why**.
- First checkpoint commit **within your first 15 tool calls**; sync Status + boxes + a `next=`
  Progress log line on every commit (a ticket whose face lags HEAD gets its work redone by the next
  agent).

## Progress log（append-only；每个 commit 一行 `- [ISO-UTC] agent=... did=...`）

- [2026-09-21T02:05:00Z] agent=ticket74 did=**基线复测 + 落盘 checkpoint（未改任何码）**。
  `go test ./internal/ball/ -run TestC21 -v` 四条全绿，报告口径复核为 **29 of 105 零消费者**、
  **15 条表外豁免**（hit.go 5 + hotkey_windows.go 1 + liquid.go 9），与票 69 交付时一致。
  family 2 的**前提已独立验真**：`git grep -n "SwimLevelGain\|SpinLevelGain" -- internal cmd` 只有
  `tokens.go:351-352` 的声明和 `tokens_table_test.go:380-381` 的金标准条目 ⇒ **零消费者成立**；
  真正的驱动在 `liquid.go:224`（`SpinBaseRadPerS + SpinLevelRadPerS*level + SpinSummonRadPerS*summon`）。
  ⚠ 并且**"偏心增益 0.55"那一半也不是没被实现，而是被写成字面量**：`renderer_windows.go:274` 的
  `or := liqOffset[bi]*R` 根本不读电平，电平改的是 `:543 shrink := 1 - 0.12*v.LiquidLevel` /
  `:544 flow := 1 + 0.10*v.SummonFlow` —— 两个**未命名的裸字面量**在改像素，机器检查按 `const` 声明
  枚举，因此结构上看不见它们（记为缺口，见 AC#2 判定）。
  next=一个 commit 内做完 family 2 的"删除 + 换指向"：`tokens.go` 删 `SwimLevelGain`/`SpinLevelGain`，
  表 :153 行改指 `liquid.go` 的 `Spin*RadPerS` 真实驱动，同 commit 把它们从 `c21OutTableExempt` 摘掉，
  并让几何用例在**整包**范围内解析表名（否则 liquid.go 的名字会被判"tokens.go 无此外常量"）。
- [2026-09-21T02:40:00Z] agent=ticket74 did=**family 2 + family 1 一个 commit 落地**（AC#1 勾）。
  删除与换指向同 commit，表没有一秒钟指向死常量：
  ① `tokens.go` 删 `SwimLevelGain`(0.55)/`SpinLevelGain`(1.0)，同一处新增
  `LiquidGatherPerLevel`=0.12 / `SummonFlowSpread`=0.10 —— 它们就是 `renderer_windows.go:543-544`
  里那两个裸字面量，提出来才有名字可进表；② 表 :153 一行改两行（几何 / 转速），另加包络门限、
  静止判据、`hit.go` 边界共 5 行，把 D7/D8 两族 **11 条**从 `c21OutTableExempt` 提进表并钉住值
  （新增 `c21MotionGolden`，豁免表只剩 4 条真不是 token 的 Win32/热键名）；
  ③ `tokens_table_test.go` 的名字解析从"只认 tokens.go"扩到**整包**（`c21PackageConsts`），
  并且新加一条反空跑断言：表名了 tokens.go 之外的常量而 `c21MotionGolden` 没钉值 ⇒ 红。
  AC#1 的 grep 证据（HEAD 工作树实测）：`git grep -n "SwimLevelGain\|SpinLevelGain" -- internal`
  ⇒ 只剩注释与本表自述，**零 `const` 声明**；`git grep -n "LiquidGatherPerLevel\|SummonFlowSpread"`
  ⇒ `renderer_windows.go:545-546` 两处真实像素路径；`git grep -n "Spin.*RadPerS" -- internal/ball/liquid.go`
  ⇒ `liquid.go:224-225` 的 `omega`。
  新行**已验真会红**（临时改码，跑完即还原）：`liquid.go:28 SpinLevelRadPerS 4.20 -> 4.35` ⇒
  `TestC21GeometryRowsMatchCodeConstants` FAIL 并点名 `c21-native-tokens.md:159`
  （"SpinLevelRadPerS = 4.35, but this row states no such number"）；还原后 `git diff --quiet` 空、67 条声明复绿。
  零消费者报告 29 → **27**（少掉的就是那两条死常量），表外豁免 15 → **4**。
  门禁：`gofmt -l internal/ball` 空、`go vet ./internal/ball/...` rc=0、
  `go test ./internal/ball/... -count=2` ok、`go test ./internal/ball/ -v` RUN=54 / SKIP=0 / PASS=54。
  next=family 3：在表里开「未接线豁免」小节，27 条零消费者**逐条一句理由**，并把消费者报告从
  "只打印"改成"未归因即红"（理由放表里、由测试解析，避免理由清单本身成为第二份会漂移的真相源）。
- [2026-09-21T03:05:00Z] agent=ticket74 did=**family 3：零消费者从"打印"变成"必判"，理由写在表里由测试解析**（勾 AC#2）。
  新增表小节「未接线豁免（documented but unconsumed，票 74 family 3）」= **24 行覆盖 27 个 token**
  （11 个 `Tint*` 死色端点、`FgPrimary`/`FgTertiary`/`GlassRing`/`GlassHi`/`GlassHiSoft`/`InfoLine`/
  `OrbShadow` 7 个面板侧色、`looks[*].Deep`、`BallSizeSmall/LargePx`、`Dur*Ms` 三档、`FontSizeMonoPx`、
  `CountdownFontPx`、`FirstRunGuidePulseMs`），每行一句**指向具体去处**的理由。
  为什么理由放在 md 里而不是 Go map：豁免清单若另存一份，它就是第二份会漂移的真相源——本票要修的正是这个。
  `TestC21TokenConsumerReport` 现在双向比对：**漏一行 ⇒ 红**、**已重新接线仍挂豁免 ⇒ 红（stale）**、
  **理由 < 20 rune 或不含任何代码名/数字 ⇒ 红**（`c21MinExemptReason`，挡"暂未使用"这类占位话）。
  口径同时收紧两处：① 被检 token 集合从 105 → **116**（`c21MotionGolden` 的 11 条 hit.go/liquid.go
  常量纳入同一条规则）；② `c21CollectRefs` 不再把**声明本身**当使用（tokens.go 本来就跳，但票 74 起
  扫描覆盖 hit.go/liquid.go，那里声明与使用同文件——不修就会把"只有 `const` 行提到它"读成已接线）。
  变异证据（临时改码，跑完还原）：`statevisual.go:132 * SleepRestRatio -> * 0.62`（`grep -c MUTATION-74D`=1
  且行内容打印为证）⇒ `TestC21TokenConsumerReport` FAIL 点名 `SleepRestRatio is tabled but nothing in the
  ball package reads it … 未接线豁免 table has no row for it`；还原 `git diff --quiet` 空。
  报告口径：ZERO-CONSUMER **27 of 116**、`未接线豁免: 24 row(s) attribute all 27` ⇒ **未归因 0 条**。
  ⚠ AC#2 的**已知缺口（写进表末段）**：`renderer_windows.go` 仍有**未命名裸字面量**在改像素
  （`liqRate`/`liqPhase`、渐变 stop 的 0.42/0.78、`permille` 的 775+225…），机器检查按 `const` 声明枚举，
  结构上看不见它们；把它们提成 token 属改绘制代码，不在本票范围。
  next=收紧两处判据：几何值匹配 containment → **单射**（已实测今天的 containment **会放过**
  `tokens.go:361 DockTriggerPx 16 -> 160`：整包 `go test ./internal/ball/` 仍 `ok`、几何用例仍 PASS，
  只在日志里多打一条 "MATCHED WITHOUT THE PROMISED UNIT (5): … DockTriggerPx=160 @ …:155" ——
  它借走了同属 `DockAnimMs` 的那个 160ms），以及给 `FontFamily` 补**值断言**
  （`--font-sans` 里确有 "Microsoft YaHei UI"，所以值断言可以做成真的）。
- [2026-09-21T03:35:00Z] agent=ticket74 did=**收紧两处判据 + 五次变异全绿反证**（勾 AC#3）。
  ① 几何值匹配 containment → **单射**（`c21TakeClaim`：匹配到的数字被"用掉"，同行第二枚常量不得复用同一个
  写出来的数字；`c21MatchClaim` 删除，不留死码）。表的对应改动：`描边宽` 一行把 1.5px 写三遍、2px 写两遍
  ——**同值几枚常量就写几遍**成为表的书写规则，这条规则本身就是可读性收益。
  ② `FontFamily` 补**值断言**：字体行如实写出 `Microsoft YaHei UI`，测试要求行的 `值/用途` 文本包含
  `c21StringTokens` 钉住的字符串 ⇒ 换任何字体名都会红（此前只查名字双向，值随意）。
  ③ 顺带把 `c21StringTokens` 缺钉值、以及"钉了值但行里没写"两种空档也判红。
- [2026-09-21T03:36:00Z] agent=ticket74 did=**AC#3 变异台账（5 次，逐条：改→grep 证落盘→点名红→还原→`git diff --quiet` 证干净）**。
  | # | 改动 | 落盘证据 | 结果 |
  |---|---|---|---|
  | M-A | `liquid.go:28 SpinLevelRadPerS 4.20 -> 4.35` | `grep -c MUTATION-74A`=1 | **FAIL** `TestC21GeometryRowsMatchCodeConstants`："…:159: SpinLevelRadPerS = 4.35, but this row states no such number" |
  | M-B | `hit.go:27 ClickTolerancePx 4.0 -> 6.0` | `grep -n MUTATION-74B` 打印该行 | **FAIL** 同一用例："…:162: ClickTolerancePx = 6, but this row states no UNCLAIMED number equal to it (its numbers are: 8px, 4px, 96dpi)" |
  | M-C | `tokens.go DockTriggerPx 16 -> 160`（**收紧前整包 `ok`、用例 PASS**，只在日志多打一条 "MATCHED WITHOUT THE PROMISED UNIT: DockTriggerPx=160"） | `grep -n MUTATION-74C` | 收紧后 **FAIL**："…:164: … two constants are sharing one written digit" ⇒ containment→单射的必要性由此坐实 |
  | M-D | `statevisual.go:132 * SleepRestRatio -> * 0.62` | `grep -c`=1 且 `sed -n 132p` 打印 | **FAIL** `TestC21TokenConsumerReport`："SleepRestRatio is tabled but nothing … has no row for it" |
  | M-E | `tokens.go FontFamily "Microsoft YaHei UI" -> "Segoe UI Variable Text"` | `git diff --stat` 1 file/1 insertion | **FAIL** 几何用例："…:155: FontFamily = "Segoe UI Variable Text", but this row never writes that string" |
  还原复核：`grep -c MUTATION-74B` = 0、`git diff --quiet -- hit.go tokens.go statevisual.go liquid.go` **四文件皆空**、
  复跑 `go test ./internal/ball/... -count=2` → `ok 0.355s`。
  ⚠ 一次**假变异**被纪律本身抓到：M-B 第一次 sed 模式带了中国缩进的 `\t`，替换没落盘 ⇒ 测试"绿"，
  `grep -n MUTATION-74B` 零命中当场暴露——这就是"先 grep 证明落盘再跑测试"存在的理由（A24 系列第 N 次）。
  next=AC#4：`表 ↔ design/assets/tokens.css` 的第三腿要不要今天机器化（先量：CSS 129 条声明、表内 78 配色行
  声称与 CSS 同源）——能加就加，加不动就写成有名有姓的移交票；然后 AC#5 收尾门禁。
- [2026-09-21T04:05:00Z] agent=ticket74 did=**AC#4 走了"加机器检查"那条路，没有移交**：新增
  `TestC21TableColourRowsMatchTokensCSS`（`internal/ball/tokens_table_test.go`）把表头那句
  "真相源是 `tokens.css`" 变成第 5 条用例：每条引用了 CSS 变量的配色行必须 **三方一致**
  （`design/assets/tokens.css` 的声明 = 表的 `tokens.css 值` 列 = `tokens.go` 里的颜色）。
  实测 **78/78 行通过**（:root=129、light=65 条声明被解析；36 行 look 色按定义无 CSS 出处，如实跳过并记在表头）。
  - CSS 解析器按 `;` 分段（`--tint-success-hi` 与 `-lo` 同行打包）、剥行内注释、支持跨行声明；
    分层背景 `--bg-overlay: var(--sheen), rgba(...)` 取**颜色分量**比对（`--sheen` 属「范围」段排除的
    面板专用声明），并把这个读法**写进表里**：两行 `--bg-overlay` 的 CSS 列改为照抄原文 ——
    原来只记了 rgba 那半，是同一类"记录比源码少一层"的小 A33。
  - 顺手把 `FontFamily` 的 CSS 腿也钉上：`--font-sans` 里必须真的有 `Microsoft YaHei UI`，
    否则表上"CJK 头"那句话就是假的。
  - 变异检验 **M-F**（表侧）：`docs/.../c21-native-tokens.md:55` 的 `#F2F3F5` -> `#F2F3F6`（一位十六进制）
    ⇒ 新用例 **FAIL**，报文同时点名表行与源码行："the table states --fg-primary = "#F2F3F6", but
    design/assets/tokens.css:61 declares "#F2F3F5""。还原 `grep -c F2F3F6` = 0、复跑绿。
    ⚠ 第一次 M-F 的 sed 打在错行号上（我插了 3 行说明，行号右移），**没落盘的变异跑出"绿"**，
    `grep -c` 零命中才抓住 —— 与 M-B 同一次教训，先证落盘再跑测试。
  - 表头那张"今天有无机器检查"的三分类表**如实改写**（原文写 look 色"无机器检查"，票 69 之后已不成立），
    A24-D4 的完成判据拆成两半记账：机器断言那半已落地、"升进 `tokens.css`"那半仍待 owner 解冻 SPEC-08，
    所以"不得把这 36 行当契约级事实引用"这条限制**保留不撤**。「审计」小节同步改成可复跑的入口命令。
  - AC#5 门禁实测：`gofmt -l internal/ball` 空；`go vet ./internal/ball/...` rc=0；
    `go test ./internal/ball/... -count=2` → `ok 0.462s`；`go test ./internal/ball/ -v` →
    **RUN=55 / SKIP=0 / 顶层 PASS=49（其余 6 条是子测试）/ FAIL=0**；`-v -count=2` → **RUN=110 = 2×55 个不同名、SKIP=0**
    （`-count=2` 与"非匹配 -run 也印 ok"两个坑都按票面要求用 `=== RUN` 计数排除）。
    `cd tools/d22scan && go run . -root <abs>` → `clean - no D22 ban violations`，
    唯一 `NOT COVERED` 是既有的 ban #6 `frontend/`（树不存在，属票 67/70 台账，不是本票引入）。
  next=收尾：把 AC→commit→证据的 1:1 裁决表写进票面，确认无未归因分歧；本票不再动码。
