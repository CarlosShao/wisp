# 票 211 补读数腿 `211-r1`（readings-only：五格判据的读数首次以 211 名义交付）

- 日期：2026-11-08 21:1x +08。腿＝`211-r1`（写腿·只交读数·**零产码字节**）。
- 读数锚：共享工作树起手 HEAD `75ff9c4eb`（AC#1/AC#5 两条 grep 在此取）；导出树 `git clone --shared` ＝
  `/d/tmp/wisp/20261108/leg78/tree-leg78`，HEAD **`ca790c46e`**（A/B/C/E 四发在导出树取，起手三核③：
  色/突变一律在导出树跑，⛔ 不把别人的半保存态算进本腿）。
- 写面：本件＋票 211 尾部一条 Progress bullet。⛔ 翻框（五格归编排者凭本件裁）；⛔ 动 `bridge.go`／`budgets.go`／任何 `.go`。

## 起手三核

1. 票面五格判据原文仍在（`.scratch/wisp/issues/211-subagent-pool-cannot-exceed-the-d38d-tool-ceiling.md` 判据节，逐字现读）。
2. 甲支代码已在树（票面 09-29 增量节自陈"已随票 197 的修复腿落地，见 A411"）：
   `MaxConcurrentSubagents = 4`（`internal/tools/subagent_197.go:85`）＋常驻判据 `Test197SubagentPoolNeverExceedsBridgeCeiling`
   （`internal/tools/subagent_197_test.go:412`）都在。**但 `docs/evidence/s1/` 下 211 名义证据件＝0 枚**
   （尺＝`ls docs/evidence/s1/ | grep -c '^211'`＝0，本腿起手现跑）⇒ 五格"落地腿必须逐条给读数"至今无人交付，本腿补上。
3. 红/绿非他腿在飞造成：读数全部取自独立导出树（HEAD `ca790c46e`）；共享树 `internal/tools` 起手 `git status --porcelain` 空。

## AC#1：池常量实测等于桥天花板（现跑 grep 并排）

```
internal/tools/subagent_197.go:85:	MaxConcurrentSubagents = 4
internal/tools/bridge.go:24:const MaxToolConcurrency = 4
internal/agent/budgets.go:47:const MaxToolConcurrency = 4
```

两枚独立常量（票面现量表点名 `subagent_197.go` 把 4 写成字面、不做别名，便于本比对）实测同值 4＝4。

## AC#2：常驻判据在树，且"种 8 必响"正控本腿现跑复现

- 判据本体＝`Test197SubagentPoolNeverExceedsBridgeCeiling`（`subagent_197_test.go:412-420`）：比较两枚独立常量，
  `MaxConcurrentSubagents > MaxToolConcurrency` 即报错——非恒真（判据自己比较两个读数，不是"相等"断言）。
- 正控（导出树，`sed` 把池常量 4→8 后只跑这一枚）：

```
=== [B] mutation pool 4 to 8 ===
85:	MaxConcurrentSubagents = 8
=== RUN   Test197SubagentPoolNeverExceedsBridgeCeiling
    subagent_197_test.go:414: 池 8 大于桥的 D38d 天花板 4：多出来的 4 枚会被桥排成队，名册却说它们在跑；
    诚实数 = 天花板（票 211 甲），要更多并发得先动契约（票 211 乙），不是改这枚常量
--- FAIL: Test197SubagentPoolNeverExceedsBridgeCeiling (0.00s)
FAIL
B_RC=1
```

- 还原核（用 `cp` 备份还原，⛔ checkout；还原后同枚重跑）：

```
=== [C] restored rerun ===
ok  	github.com/CarlosShao/wisp/internal/tools	0.044s
C_RC=0
```

- 交付时点该正控的首发读数在 `docs/evidence/s1/197-subagent-entity-r1b.md`（判据注释自陈）；本腿是 211 名义的独立复跑，两读数同向。

## AC#3：池满第 5 枚硬拒＋可读理由，不占名册行、不占池位

- 判据＝`Test197FullPoolRefusesNextSpawnWithReadableReason`（`subagent_197_test.go:535` 起），PASS（见下方 [A] 全绿读数）。
- 理由的生产形状＝`internal/tools/subagent_197.go:288-291`：拒句带
  `strings.Join(t.d.Roster.RunningSubagentIDs(), ", ")`——"当前在跑的是哪几枚 taskID"在理由里逐字在。

## AC#4：池=4 的测量走真桥（spawnDirect 绕行已删，不是文档化保留）

- 判据＝`Test197SubagentPoolCapsAtBridgeCeiling`（`subagent_197_test.go:450` 起），PASS。
- 该测试头注释逐字自陈：`This replaces the pre-211 leg that drove the tool directly (spawnDirect), and
  spawnDirect is gone rather than documented`——四枚并发 `task.spawn` 经父任务真桥同跑、四枚桥位真忙、
  名册行数＝池数、全部收尾且占位回收，全在用例体内断言（`:479-:533`）。

## AC#5：两枚"4"一字未动

```
=== [E] AC5 constants untouched ===
git diff --stat -- internal/tools/bridge.go internal/agent/budgets.go   ⇒ 空
E_RC=0 (empty diff expected)
共享树 git status --porcelain -- internal/tools/bridge.go internal/agent/budgets.go internal/tools/subagent_197.go ⇒ 空
```

本腿全程零 `.go` 字节（突变只发生在导出树、用 `cp` 备份还原并留 `.leg78bak` 孤儿件在导出树，共享树无痕）。

## 汇总色（导出树，HEAD `ca790c46e`）

```
=== [A] tools 197 pool tests as-landed ===
--- PASS: Test197SubagentPoolNeverExceedsBridgeCeiling (0.00s)
--- PASS: Test197SpawnDescriptionNamesTheRealPoolCap (0.00s)
--- PASS: Test197SubagentPoolCapsAtBridgeCeiling (0.00s)
--- PASS: Test197FullPoolRefusesNextSpawnWithReadableReason (0.00s)
PASS  ok github.com/CarlosShao/wisp/internal/tools 0.076s   A_RC=0
```

（四发含 `Test197SpawnDescriptionNamesTheRealPoolCap`＝"给模型读的说明文本与常量不分家"那枚伴生钉，一并绿。）

## 边界与不冒领

- 本腿**零产码、零框翻**：甲支代码归票 197 修复腿（A411），五格翻勾归编排者凭本件裁。
- 判据注释里那段"为什么池仍不许超过天花板"的自陈（票 222 giveBackWhileWaiting 之后语义已变）属注释散文、无仪器看管
  （注释自己逐字承认）——本腿照抄现状，不另立尺（那是票 222/237 的账）。
- 撤销口令：本腿只追加两个 `.md`，`git revert` 两笔即回原状。
