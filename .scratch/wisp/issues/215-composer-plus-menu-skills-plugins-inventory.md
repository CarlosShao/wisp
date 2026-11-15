# 215 — **加号菜单第三类："技能／插件／MCP 要能列出来、能开关、能看出这次生效了没有"**（配置里那一堆键今天没人读，别家做成三枚分包）

- Status: **待派（owner 09-28 22:1x 口头接单，原话见票 213 第一节）**。
- 为什么单独一票：这一类的**后端不是新能力，是"已配置的东西没人读"**——和票 214（已建没接）同族，但射程不同（214 管一次对话的输入，215 管**长期挂着的东西**）。

## 现量（起手逐条复算）

> ⚠ **09-28 22:3x 更正本票起手第一格**（我原来那句"配置里那一堆键今天没人读"**对插件成立、对技能不成立**；调研腿 `survey-plusmenu` 顶回来、我复跑确认）：
> `grep -c Skill internal/config/schema.go` ⇒ **0**；`grep -c Plugin` ⇒ **8**。
> ⇒ **技能这一维连配置键都不存在**，不是"键在没人读"。**本票据此拆成两半**：
> **215a＝插件**（有 8 处配置、缺"读者＋清单＋生效与否"，属〔建了但没接〕，便宜）；
> **215b＝技能**（要先补配置面 ⇒ **动 D36 那棵 section 树＝配置模型契约面**，按 `AGENTS.md` §0.2 要人工批准，**批准只落 `A##`、不改 `PLAN.md` 一字**）。
> 另：别家在菜单里**都不显示**技能/插件的开关联调状态（五家现量一致），所以"显示开关"不是及格线、是**我们自己的选择**——
> 真正该抄的是 minimax 那枚 `SkillSnapshot`（`winners`／`losers(reason+winnerLocationUri)`／`diagnostics`／`refresh_metrics` 全是**可数字段**，`skills/src/types.ts:49-115`）：**"配了没生效"要能报出是谁赢了、为什么**。

| 事实 | 读数 | 尺 |
|---|---|---|
| 配置里有这些开关 | **插件 8 处命中／技能 0 处**（技能连键都没有） | `grep -c Plugin internal/config/schema.go`＝8；`grep -c Skill internal/config/schema.go`＝0 |
| **"写了不管用"的旗还在** | `internal/config/unwired.go` 那套"未接线的键直接响亮拒收"是本仓的强项（账 `A408`），**本票不许把它拆了** | `internal/config/unwired.go` |
| 别家做成独立分包 | `ui-skill`、`ui-plugin-manager`、`ui-settings-plugin-inventory`、`ui-settings-plugins`（＋服务端 `packages/skill`、`packages/mcp`） | `docs/reports/survey-2026-09-28-dsh-ui-packages.md` |
| 我们的取舍 | MCP 是 **D13 明确取舍过**的（架构冲突论未被推翻，但别家有一份可读的 MCP 运行时：连接池 911 行＋持久名字治理 258 行） | 账 `A413`/`A415` 那一轮＋`PLAN.md` D13 |

## 这一票要交付什么

1. **一枚可枚举的清单**（技能／插件各一份）：名字、版本、开没开、**这次这一轮它到底生效了没有**（生效＝被装配进去，不是"配置里写着 true"）。
2. **开关的落点**：从菜单里开关 ⇒ 走现成配置写路径（**hot-tier 的立即生效、需要重启的要明说**），并且**开关动作本身要进审计**（三档都写）。
3. **来源与安全**：外部插件/技能带进来的文本一律按**不可信外部内容**处理（票 200 那条纪律：读说明书≠授权，**不许因为它在菜单里就改权限档位或 `allowed_dirs`**）。
4. **MCP 不在本票**：要不要接 MCP 是 **D13 的取舍**，本票只保证"如果哪天接，菜单里有位置"——**这一条不许顺手做**。

## 判据

- [ ] **AC#1**：清单可枚举，**每一条都有真装配读者**（现跑 grep 点数：改前"某键零读者"⇒ 改后 ≥1）。
- [ ] **AC#2**：**"配置开着但没生效"必须能被看出来**——正控：塞一条"配置 true、装配里没人读"的假项 ⇒ 判据必红（这一形就是 `A408` 点名的"文案在、控件不在"的镜像）。
- [ ] **AC#3**：开关动作**跨档位一致**：默认最严档里关一个技能要能立刻生效、不需要重启；做不到就**在清单上写明"这一枚要重启"**，不许静默假装生效。
- [ ] **AC#4**：外部内容带进来的指令**改不动权限**（复用票 200 那枚 hostile-file 判据的形状：投一枚写着"以后都允许"的插件说明 ⇒ 档位与 `allowed_dirs` 一字不变）。
- [x] **AC#5**：`unwired.go` 那面旗**一枚不少**（现跑点数：改前＝改后）。

## 禁区

不接 MCP（D13 取舍）；不新增 D34 工具行；不新增 C17 方法名；不拆 `internal/config/unwired.go` 的拒收旗；
`frontend/**`、`design/**` 零写；`PLAN.md`／`docs/specs/**`／`thresholds.go`／golden／`allowlist.txt` 一字不动；
只 commit 不 push、显式 pathspec、禁 `--amend`/`reset`/`rebase`/`stash`/`checkout .`/`clean`；仓内不删东西。

## 前置（三票共用，写死在这里）

**界面点不动**：网页→Go 那一跳今天只有命令行路（账 `A410`）。⇒ **票 181／33 的入向线接通之前，213/214/215 三票的终态都是"Go 侧可枚举＋命令行可执行＋界面看不到"**，
每一票交件时必须具名写这一句，**不许把"菜单做好了"当成 owner 能用的功能**（他今天已经因为我这种说法返工过一次）。
— 2026-10-09 12:0x 编排者腿 boxes-1 搬框：本节原为散句 bullet，现改成 - [ ] 形以进全勾尺；⛔ 判据文字一字未改、⛔ 未翻任何框。

## Progress log

- [2026-11-07 23:1x +08] agent=215-plugins-leg did=**215a（插件半）Go 侧落地＋215b（技能半）如实登记**：`panel.Snapshot` 新增 reader-gated 顶层键 `plugins`（载体 `internal/panel/plugins_inventory_215.go`＋`composer.go` 字段＋`pump.go` 读口与装配块，照 213/187/65 范式、未造第五范式）；行源自 `internal/config.PluginsSection`（真源现跑核过：Plugin 8 命中/Skill 0，与 09-28 更正一致），泵拷贝后按 ID 排序（map 顺序不上线、泵不写回读者 slice）、nil rows 渲染 `[]`；引擎（50/51）未建 ⇒ 线上**无生效位**，`enabled=true` 与「[plugins] 各键今天零读者」的 `EngineState` 句**同包永不分家**（AC#2 正控常驻钉死）——「配置开着但没生效」在形状里；215b 连配置键都不存在 ⇒ ⛔ 造假数据，线上只携 `SkillState` 登记句（补键＝动 D36＝契约面待批）。三发突变牙真红后还原（M1 摘 pump 装配块/M2 排序比较器恒 false/M3 拆 reader-gating；M2 第一发 sed 打挂编译作废重打，如实记）；无读者基线＝无 `plugins` 键＋顶层恰 4 键（AC#5 侧写：`internal/config/**` 零字节，unwired guard 表 5 行改前＝改后）。门禁：gofumpt 点名 4 枚**空**／vet·build rc0／`go test ./internal/panel` 恰 3 红＝**非本面**（两枚 recon 逐字 `Go Snapshot emits [plugins] that interface PanelSnapshot does not declare`＝213 立的既定程序、Q-51 界面腿；`TestC21DesignTokensFourWayAgree` 既存红＝c21 证据件点名 GlassRing/GlassHi/GlassHiSoft 不在 ball tokens.go、`internal/ball/**` 在飞零触碰。★起手后异动如实记：起手基线 recon 红为 `[commands modelSwitch ballStyle run]`＋`ComposerState files`，收尾前**界面腿在飞件落地**（panel.ts 未提交修改、Snapshot 11 键 reconciled），那批红转绿，本腿增量即 `[plugins]` 一枚预定红）／d22scan clean（曾抓 1 枚 phantom-citation＝注释路径断行，修复复跑 0，如实记）／LF。**诚实边界**：生产装配根一行（`cmd/wisp/run.go` 加 `Plugins:` 读 config）未接——本波 `cmd/wisp/**` 被 219 腿独占，归下一腿；交付 2（开关）与 215b 未做，登记于证据件；终态照票面＝**Go 侧可枚举＋命令行可执行＋界面看不到**（账 A410）。证据件 `docs/evidence/s1/215-skills-plugins-inventory.md`；⛔ push。
- [2026-11-07 23:2x +08] agent=215-plugins-leg did=**commit 后读数漂移补录**：commit `5341ead7a` 自跑复验，`internal/panel` 红 3→**1**（只剩 C21＝ball 族既存）——两枚 recon 转绿，原因＝界面腿在飞件给 `plugins` 键补了 TS 声明（panel.ts:372-384 `plugins?: PluginsSection`，未提交、引用本腿载体件），recon 程序红随界面腿落地即转绿，程序两头一次对上。本腿写面与 AC 读数不变；读数漂移全录证据件 §⑥.0。
- [2026-11-08 09:5x +08] agent=215-assembly-leg did=**装配根接线落地（上一腿欠的那一行）**：`cmd/wisp/run.go` PumpSources 字面量加 `Plugins: func() panel.PluginInventorySection` 读口（run.go:1104，照 65/187/213 三枚同款；净 +37 行零删除），读者＝`mgr.Config()`（Manager 锁内深拷）映射 `[plugins]` 节进载体＋两句 canonical 常量原样随行——internal/config **无** plugins 专用 reader（Options 族仅 BallStyle/ChatSwitch，现跑核过），装配根 pre-shape 正是 pump.go:252-264 字段文档点名的范式，未造第五范式。新测试件 `cmd/wisp/panel_plugins_215_test.go` 照 187 三牙：字面钉＋真 Manager 接线读数（plugins215Base 喂真 NewManager，期望硬编码：rows 恰 [clip zed] 泵排序、字段 DeepEqual 逐字、tier2Enabled 逐字、enabled=true 与两句同包）＋无读者 wire bytes 零键突变；**M1 真突变＝sed 摘 run.go 接线块 → 字面钉 FAIL → cmp 字节还原 → 复绿**（先试 overlay 换 pristine 副本没红＝作废如实记：字面钉读盘上文件、不读 overlay 编译视图）。门禁：gofumpt 点名 2 件空／vet rc0·build rc0（★中途 raw vet 两度红＝219 腿在飞件 approval_always_219ac2_test.go 正被实时编辑：undefined: sync→repoRoot219、mtime 晚于本腿基线 vet rc0；照 62-v1 先例 overlay 隔离自证，工作树零触碰，提交时 raw 已自愈 rc0——归属取证全录证据件 ⑤）／`go test ./cmd/wisp -run 'Test215|TestPlugins' -count=1` 3/3 PASS 0.087s／`go test ./internal/panel -count=1` 改前=改后恰 1 红 TestC21DesignTokensFourWayAgree（ball 族既存非本面）／d22scan clean（internal 963＋cmd 241＋frontend 117＋design 39）／改动件 CR=0。**运行进程从此发 `plugins` 键**（215a 线上面收口）；界面点不动照票面前置（账 A410），TS 声明已由 65 界面腿落（a3fdb50dd）。证据件 `docs/evidence/s1/215-plugins-assembly.md`；⛔ push。
---

## 收件（2026-11-15 14:0x，`9ed201f5f` ＋ 本件由编排者提交）：**满足 1／部分 2／不满足 2（总 5）**

**尺1**：腿**未 commit**（照 195/351/314 的惯例交收件批次），件与 probes 由**我**提交：`docs/evidence/s1/215-receipt.md` **42155 B**、`.scratch/wisp/probes/215/**`。
**尺2**：件在盘，五列逐格＋三件读数＋探针索引（42 KB／LF／0 CR）。
**尺3 我核的三件，全部成立**（其中一件是我自己的搜法漏了）：
- `cmd/wisp/run.go:1307` = `Plugins: func() panel.PluginInventorySection {` ✓（结构体字面量字段）
- `cmd/wisp/panel_resident_pump.go:144` = **`src.Plugins = func() panel.PluginInventorySection { … }`** ✓ —— ★我第一次搜 `Plugins:` 漏了它，因为这是**赋值语句形式**不是字段形式；腿说的"两处"**对**
- 那枚"只钉字面"的测试定义在 `panel_plugins_215_test.go:82` ✓
- `go vet ./cmd/wisp` 现在**无输出（rc=0）** ⇒ 452 造成的临时编译错已不在

### 五格判语（⛔ 只翻 AC#5）
| 格 | 判语 | 依据 |
|---|---|---|
| AC#1 | **部分** | 插件半成立（载体 `internal/panel/plugins_inventory_215.go:107`＋读口 `pump.go:296`＋**两处**生产赋值点；改前 `internal/config` 包外消费者 **0** ⇒ 改后两处）；**技能半不存在**（`Skill` 命中 0，只有登记句）。缺：技能清单（215b＝动 D36 待批）＋`run.go` 那格的行为钉 |
| AC#2 | **部分** | 可见性成立且**实测有牙**（摘 `panel_resident_pump.go:252` 的未装配句 ⇒ `Test65…ConfigTruth` 红于 `:141`）；但 AC 字面「塞假项 ⇒ 判据红」**无承载体** —— 今天没有任何一行可能是"真生效"，假项与真项同形 |
| AC#3 | **不满足** | `Manager` **23 枚 setter 无一枚 plugins**，交付 2 整个没开始 |
| AC#4 | **不满足** | 生产侧零读入口（`internal/plugin/manifestslot` **无任何生产 import**）⇒ "改不动权限"是**空真**（vacuously true），且无 hostile-plugin 判据 |
| **AC#5** | **满足** ✅ | 两个 215 commit 对 `internal/config` **零字节**；overlay 摘一行处置 ⇒ `unwired_test.go:372` **红**（实测） |

### ★这一份收件最值钱的一条读数
**摘掉 `run.go:1307` ⇒ `0 枚红`（3/3 PASS）**。那格只有一枚**盘上字面钉**（`Test215AssemblyRootCarriesThePluginsReader`，它断言源码里有那串字），**没有行为牙** ⇒ 删掉它没有任何行为测试会响。
对照：另一处 `panel_resident_pump.go:144` **摘掉就红**（`Test65ResidentPumpCarriesConfigTruth`）⇒ 那处是**真承重**的。
⇒ 结论：**`run.go:1307` 那一格是装饰性的**，今天"装配根接线落地"这句话**对了一半**。

### 票面 ≠ 盘上（腿逐条对上/点名，7 条）
`I1` 票面只记 `run.go` **一处**赋值点，盘上**两处**（常驻泵那处由 `a5f84886f` 后落，且**唯一有行为牙**）／`I4` 票面 `grep -c Plugin schema.go`＝**8 从未复现**（立票提交 18／HEAD 19）／`I3` `run.go:1104`→**1307**／`I2` `panel.ts:372-384`→436-448/516／`I5` C21 今天绿／`I6` `unwiredKeys` **5→3**（`7bcc2ae08`／票 163）／`I7` Options 族 **2→7** 枚。
⇒ I1／I4 记进本票（票面数错了两处，改法归我，已在下节）。

### 一处我记录的口径（⛔ 不是腿的问题）
它用「沙箱 cwd ＋ mutant run.go」把 `Test215AssemblyRootCarriesThePluginsReader` 复现成红（`FAIL panel_plugins_215_test.go:88`）—— 这是**为了证明那枚测试存在**，⛔ 不是"接线有牙"的证据。两件事必须分开读，我这边已分开。
