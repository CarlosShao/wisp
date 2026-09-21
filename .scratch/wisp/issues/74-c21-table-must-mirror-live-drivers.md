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
- [ ] **AC#3** Mutation proof for the new/changed assertions: pick at least 2 rows, change the code
  value, show the specific test go red by name, grep-prove the mutation landed, restore, grep-prove
  it is gone (`git diff --quiet` on both files).
- [ ] **AC#4** Table↔`tokens.css` gap: either add the cross-check, or move this AC to a written
  hand-off (name the ticket). **"Deferred" without a ticket number is not an option** — the owner's
  standing rule is that postponed work must be registered with its completion criteria.
- [ ] **AC#5** Gates: `gofmt -l internal/ball` empty, `go vet ./internal/ball/...` rc=0,
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
