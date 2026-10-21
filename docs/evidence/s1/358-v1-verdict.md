# 票 358 终裁腿 358-v1 - 判语件 part 1/3（非实现者，裁决者 NOT 实现者）

被裁两实例＝`358-r1`（`05f76ad6` 5 件 / `.go` 4 ＋`3ed3fdce` 1 件＝产码面）＋`358-r1b`
（`43f44f59` 3 件全 `.go` ＋其后四笔交付件 `a8ae09ce` / `10214eec` / `45f0b72f` / `cfa8bff2`＝出口面）。
起手锚＝`2053d3bf`；收线时 HEAD 已被并发腿推到 `9697a9bd`（只动 `docs/reports/pending-and-issues.md`；
`git diff --stat 2053d3bf..HEAD -- <358 七枚件>` 为空，故本腿读数在新旧两锚同字节）。
本腿零产码字节、零翻框、零 push；写面＝本两件＋`D:\tmp\wisp\20261021\358v1\**`。
凡报 0 之前先跑正控；临时件只建不删。

## 0 终版档位（一行一件）

| 件 | 判语 | 一句话凭据（本腿独立尺） |
|---|---|---|
| (1) 票面前提新鲜度 | **部分过期**（一处数错、一处结论对） | 票面 `:180-183` 实为 `b23c7f77` 那棵树的行号；票锚 `c50218bf` 已漂到 `:191-194`；HEAD `:400-403` |
| (2) `AC#0` 只读现读 | **成立，可翻** | 本腿自造三发历史＋正负控；形(1) vs 形(3) 改前**完全同形**经 pre-358 blob 实测复现 |
| (3) `AC#2` 落地＋牙 | **成立，可翻** | 本腿自造三发变异（摘 attempted / 让 skipped 印 / 改 compacted 一字）逐发转红、逐发还原核 sha256 |
| (4) `AC#5` 触发基数 | **成立，可翻**（附一枚残余面） | 两支各一枚用例今天响；`Source` 带到报告＋两条痕；但调用方传空 `Source` 时报告落空串（sec 4.3） |
| (5) 票 355 那颗钉子 | **改对了，且确比原来强** | 原句一字未删；新两半我各自造反例均转红；NOT 删断言换绿 |
| (6) 门禁六把 | **全 rc=0**（整包 rc=1 NOT 当判据） | 本腿现跑；整包按名作差两向皆空（sec 6.3） |
| (7) 红名册 | **新增红 0 枚** | 同尺 pre-358 基线 152 枚 vs 本树 159 枚，按名 join |
| (8) 写面与纪律 | **零越界**（按派单白名单逐笔） | 七笔逐枚 `--name-only`；正控 `b676092d` `.go`=3、负控 `3ed3fdce` `.go`=0 |

## 1 票面前提新鲜度（硬闸）- 部分过期

**结论对、行号错。** 票面第 1 条引 `compress.go:180-183` 为“循环自己 `break`”的真实发生位置。
本腿把 `compress.go` 全史逐枚现取（`git log --all -- internal/agent/compress.go`，六枚有该文件）：

```
commit 05f76ad6 : if-len=402 break=403
commit b676092d : if-len=193 break=194
commit b23c7f77 : if-len=182 break=183   <- 票面 :180-183 只有这棵树对得上
commit 1acdd032 : if-len=136 break=137
commit d25f139b : if-len=96  break=97
commit ae23da67 : if-len=96  break=97
```

=> `:180-183` 是 `b23c7f77`（09-26 票 153 那笔）的坐标；**票的落锚 `c50218bf` 上已在 `:191-194`**
（票面把 `for` 起点写 `:180`、`break` 写 `:183`，而 `b23c7f77` 里 `for`=180、`break`=183 逐字相符，
故不是笔误、是引了旧树），相对票锚已过期 11 行。

**实现腿报“今 `:343-347`、漂 +163”＝数错。** `:343-347` 在 `05f76ad6` 里是 `compressCall` 的
`trigger TriggerReading` 注释，**不是那条 `break`**；同一棵树上 `break` 在 `:403`。故“漂 +163”
（183+163=346）落在注释上，不是落在那颗钉子上。本腿复量：`05f76ad6` 的 `break` 在 `:403`
（相对 `b23c7f77` 漂 +220），HEAD 也在 `:403`（`for`=400、`if len(raw) <= c.b.KeepRawRounds`=402、`break`=403）。

**结论本身仍成立**（这是要判的重点）：`break` 那三行在 HEAD 与 `c50218bf` **逐字相同**，
且 `budgets.go` 全程零字节（`git diff --stat 5ca2472e..HEAD -- internal/agent/budgets.go` 为空；
`KeepRawRounds=3`、`refHistoryTokens=12000`、`refTokenBudget=200000`、`ReferenceContextWindow=128000`
现读在场）。=> **票面“不必动契约就能修”的技术前提成立**，错的是引用行号（属“行号会烂”那族，不属前提失效）。
**这一处 NOT 由本腿改票面**（写面互斥），只在此登记。

## 2 `AC#0`（只读现读）- 成立，可翻

尺＝本腿自写探针 `zz358v1a_probe_test.go`（`-overlay` 注入，落盘零字节），只用**改前改后都在**的面
（`BudgetsFor` / `NewCompressor` / `Need` / `Compress` / `Ran` / 三计数），故**同一份字节**可分打两棵树。
阈值＝`BudgetsFor(4096)` 现读 **384**（从预算读回，NOT 手写）。

**正控先跑**：`zz358RulerLive` 每条用例开头用**产码不产出**的记录打一发，捕获器看不见即 `t.Fatal`。
**负控**：形(1) 必须 `Need()==false`，否则 `t.Fatal`；**正控**：形(2) 必须恰留 1 条 `agent: history compressed`。

本树（HEAD）三态读数逐字：

```
(1) under threshold            rounds=2 rawAtStart=2 threshold=384 | Need=false Ran=false | 20  -> 20  | keptRaw=2 compressedMsgs=0 | trace[compressed]=0 trace[skipped]=0
(2) really compacted           rounds=6 rawAtStart=6 threshold=384 | Need=true  Ran=true  | 1236-> 780 | keptRaw=3 compressedMsgs=9 | trace[compressed]=1 trace[skipped]=0
(3) over threshold, only floor rounds=1 rawAtStart=1 threshold=384 | Need=true  Ran=false | 774 -> 774 | keptRaw=1 compressedMsgs=0 | trace[compressed]=0 trace[skipped]=1
```

形(3) 那一发的新痕逐字（`state=attempted`，12 键）：
```
agent: history compress skipped || state=attempted reason=keep_raw_rounds_floor tokens_before=774 tokens_after=774
threshold=384 msgs_before=3 msgs_after=3 compressed_msgs=0 kept_raw_rounds=1 history_changed=false
trigger_source=estimate trigger_tokens=774
```

**“形(1)/形(3) 改前是否完全同形”＝本腿实测复现（不是读 diff）**：把 `compress.go` / `loop.go` 用
`-overlay` 逐字节换回 pre-358 blob（`5ca2472e`；四枚 blob 的 `git hash-object` 与 `git rev-parse`
**逐字相等**，见 sec 8），同一份探针再打一发：

```
(1) under threshold ... Need=false Ran=false | 20 -> 20   | trace[compressed]=0 trace[skipped]=0
(3) over threshold ...  Need=true  Ran=false | 774 -> 774 | trace[compressed]=0 trace[skipped]=0
SAME-SHAPE (1)vs(3) on the pre-358 surface: Ran false/false trace 0/0
```

=> **改前 (1) 与 (3) 在可读面上完全同形**（`Ran` 都 false、痕都 0 条），病根属实；改后 (3) 多出恰 1 条
`state=attempted`，而 (1) 仍 0 条、(2) 逐字未变。**本格能判**。
附本腿自加第四形（票面未要求）：`reason=under_fold_budget` 的 attempted（provider 说超、估算没超）
也留 1 条，证明 attempted **不止“压不动”一种来源**，报告的两个数把它俩分得开（见 sec 4）。

