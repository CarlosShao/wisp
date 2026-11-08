# 票 211 五格收口＋票 354 `AC#0` 收口（closure 腿判语件）

- 日期：2026-11-08 21:3x +08。腿＝`211-cls`／`354-cls`（收口腿：凭上一腿证据件＋**自己复跑**翻勾；零产码字节）。
- 收件三把尺全过：①提交真在（`1bab25b7f`／`6b256bba1`，`git log` 现核）②件真在（两份证据件＋`readings.md`＋01/02 原始日志盘上现读全文）③**自己复跑**（下节全部读数是收口腿亲手跑出，⛔ 抄证据件）。
- 复跑锚：`git clone --shared` 导出树 `/d/tmp/wisp/20261108/leg79/tree-cls`，HEAD **`66a03a320`**；
  `git diff --stat ca790c46e..HEAD -- internal/tools/ internal/winsec/`＝**空**（两腿均零产码，区间码面同锚）＝红/绿非他腿半保存态（三核③）。
- 跑前环境：`unset GOTMPDIR TMPDIR TEMP TMP`；`PATH=$PWD/build:$PWD/third_party/sherpa-onnx:$PATH`；⛔ CGO_ENABLED=0 未设。
- 原始日志落盘：`/d/tmp/wisp/20261108/leg79/01-ac2windows-cls.log`（临时件只建不删，在导出树旁不入库）。

## 一、票 211 五格逐格判语（证据件说 ⇄ 收口腿复跑）

| 格 | 证据件说（`211-pool-ceiling-readings-r1.md`） | 收口腿复跑 | 判 |
|---|---|---|---|
| AC#1 | grep 两常量并排 4＝4（`:85`/`:24`，`budgets.go:47` 拷贝） | 同读数：`subagent_197.go:85`＝4、`bridge.go:24`＝4、`budgets.go:47`＝4，行号逐字同 | **对上 ⇒ 翻** |
| AC#2 | 牙在树＋突变正控 4→8 FAIL 红句逐字＋cp 还原 PASS | 判据本体 `subagent_197_test.go:412-421` 现读＝比较两枚独立常量、非恒真；导出树 sed 4→8 ⇒ FAIL，红句逐字同（"池 8 大于桥的 D38d 天花板 4：多出来的 4 枚会被桥排成队…"）；cp 还原重跑 `ok` | **对上 ⇒ 翻** |
| AC#3 | `Test197FullPoolRefusesNextSpawnWithReadableReason` PASS＋拒句带 RunningSubagentIDs（`:288-291`） | 四发定向跑含该枚 PASS；`subagent_197.go:287-290` 现读，`strings.Join(t.d.Roster.RunningSubagentIDs(), ", ")` 逐字在 | **对上 ⇒ 翻** |
| AC#4 | `Test197SubagentPoolCapsAtBridgeCeiling` PASS＝真桥；spawnDirect 已删非文档化 | 该枚 PASS（真桥四发之一）；`grep -rn spawnDirect internal/tools/`＝2 处命中且全在测试头自陈散文（"spawnDirect is gone rather than documented"） | **对上 ⇒ 翻** |
| AC#5 | 两枚 4 的 `git diff` 空 | `git diff --stat -- internal/tools/bridge.go internal/agent/budgets.go`＝空输出；三文件 `git status --porcelain`＝空 | **对上 ⇒ 翻** |

汇总色（导出树 `66a03a320`，四发定向 `-v`）：

```
--- PASS: Test197SubagentPoolNeverExceedsBridgeCeiling (0.00s)
--- PASS: Test197SpawnDescriptionNamesTheRealPoolCap (0.00s)
--- PASS: Test197SubagentPoolCapsAtBridgeCeiling (0.00s)
--- PASS: Test197FullPoolRefusesNextSpawnWithReadableReason (0.00s)
ok  github.com/CarlosShao/wisp/internal/tools 0.151s
```

**结论：票 211 五格全勾 ⇒ 改名 `-done`**（`git mv` 配方＋落账后双名核 `git ls-tree HEAD | grep -c "<旧名>$"` 必须 0）。

## 二、票 354 `AC#0` 逐件判语（证据件说 ⇄ 收口腿复跑）

| 判据件 | 证据件说（`354-a1-residual-readings-r1.md`） | 收口腿复跑 | 判 |
|---|---|---|---|
| 两族读数现跑 | `TestAC2Windows` 定向 6 发 rc0 | 我手同跑 6 发全 PASS（2.106s），计数逐字同：`truncated=0 deleted=0`×2、cycles=67/81、flipErrs=11/3 | **对上** |
| 正控先行（尺非钝） | shot 级 DAMAGE 行非零＋负控读 0 | 我手 3 发 DAMAGE 行逐字 `the caller's name=present (size 64 body [an object the caller never named, standing in the same directory])`＋`the swapped-in object's own name=absent`×3；负控 `name=absent` | **对上** |
| 残留①读成读数 | `nametopup_353` 三发 `standing=3/3`；`:169` 残留句 | 我手 `refusals that left the caller's name standing=3`（1 行汇总）＋DAMAGE 3 发；`:169` 句运行期打出 1 次 | **对上** |
| 归因 | 复验拒那支在 Remove 前返回＝〔乙〕代价形非回归；GREEN SCOPE 照抄 | 机制链对码成立（`unlinkPinnedLeaf` 复验 `:278` 不匹配先于 Remove 返回）；GREEN SCOPE 未见冒认 | **对上** |
| 门禁 | winsec 整包 ok 20.148s rc0 | 我手 `./internal/winsec` 整包 `ok 20.631s` rc0 | **对上** |

**★unix 半边裁定（本格能否翻的判点，现读票面判）**：

1. **票面 AC#0 原文无「两平台」要求**——全文只写「两族读数必须现跑」；同票 **AC#4 才写「两平台成对」**。同票内部一有一无＝作者会写平台要求时写得出来，AC#0 没写。
2. **残余①在 unix 机制上不存在**——`pinnedleaf_342_unix.go:173-174` 的 `unlinkPinnedLeaf`＝`_ = unix.Unlinkat(int(p.fd()), base, 0)` 直接删、**无同一性复验**（对照 Windows 版 `:278` 有 `SameFile` 复验）；「复验拒删之后名字下站着攻击者对象」那一跳在 unix 没有复验就无从发生。`353-v1` 判语同口径（unix⇄fallback 行为照旧、「unix 是否也该复验」已上交未拍）。
3. 354-a1 腿把 unix 半边如实登记〔未跑〕、未用 Windows 读数代 unix 判断＝不冒领。

**⇒ 判：AC#0 成立 ⇒ 翻勾。**

**结论：票 354 ⛔ `-done`**——AC#0（本腿翻）＋AC#1（已勾）之外，**剩 5 格**：
AC#2 确定性仪器（分子 `leftStandingForeignName`＋轮内序号尺）、AC#3 成对去牙敏感度、AC#4 颜色与门禁（两平台成对＋逐名红册）、AC#5 与 342/348 接缝、AC#6 纪律格。
归 `354-r1` 写腿＋`354-v1` 裁，本腿不越排程。

## 三、纪律与撤销口令

- 本腿零产码字节（`internal/`、`cmd/`、`frontend/`、`scripts/` 均未碰）；写面＝两张票（翻勾＋格下判语＋各一条 Progress）＋本件＋票 211 改名。
- 翻勾授权＝编排者派单（收口腿分属独立 agent，读腿 `211-r1`/`354-a1` 与实现腿均非本腿，符合「裁决者≠实现者」）。
- 撤销口令：`git revert` 两笔（收口 commit＋改名 commit）即回原状；或手工把五格＋AC#0 的 `[x]` 退回 `[ ]`、删格下判语行与两条 Progress、票 211 改回旧名。
