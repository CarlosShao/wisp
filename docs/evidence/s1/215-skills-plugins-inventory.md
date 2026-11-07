# 票 215 — composer 加号菜单第三类：技能／插件清单（215a 落地 · 215b 登记未接）· 裁决表

- 交付腿：215-plugins-leg（写面 `internal/panel/**`＋新测试件＋本件＋票 215 Progress log 一行）
- 时刻：2026-11-07 23:1x +08 · 只 commit 不 push · AC 框与 Status 零触碰

## ① 起手三核（现跑读数，非票面转述）

| 核 | 命令 | 读数 |
|---|---|---|
| unwired 旗 | `grep -rn "unwired" --include=*.go internal/` ＋ `wc -l internal/config/unwired.go` | `unwired.go` 146 行在位；`unwiredKeys` guard 表 **5 行**（risk.shell_enabled / risk.allow_shell_string / risk.shell_allowlist / net.allowlist / net.block_private_ranges）；`lockedKeyDisposition` 内 `unwired:` 前缀 **5 条**；`[plugins]` 六行标注 `not built: plugin engine is ticket 50/51`（unwired.go:140-145） |
| 技能/插件真源 | `grep -c Plugin internal/config/schema.go` / `grep -c Skill internal/config/schema.go` | Plugin **8**（`PluginEntry` schema.go:591＋`PluginsSection` schema.go:609，四字段 enabled/capabilities/net_allowlist/host_api）；Skill **0** —— 与票面 09-28 22:3x 更正逐字一致 |
| 谁在飞 | `git status --porcelain` 过滤 tracked 修改 | **0 枚 tracked 修改**（只有 .scratch／design／logs 等 untracked）；本腿基线红里没有任何在飞件成分；`internal/ball/**`、`frontend/**`、`cmd/wisp/**` 本腿零字节 |

## ② 关键判断（铁律 ⓪⓪⓪，本腿自拍）

- **215a（插件）今天有真源** ⇒ 接。真源＝`internal/config.PluginsSection`（8 处 schema 命中），照 213/187/65 的 reader-gated 范式上 `plugins` 顶层键，**不造第五个范式**。
- **215b（技能）今天无真源** ⇒ **如实登记，⛔ 造假数据**。连配置键都不存在（Skill=0）；补键＝动 D36 配置模型＝契约面（AGENTS.md §0.2 人工批准，批只落 `A##`）。线上的形状是 `SkillState` 一句话（登记「无键、待批、不列空表」），**不是**一枚会读成「没有技能」的空列表。
- **「生效了没有」的诚实答法**：插件引擎（票 50/51）未建，`[plugins]` 各键今天**零读者**（unwired.go:140-145 具名）⇒ 线上**没有也不许有** "loaded/effective" 位；`enabled=true` 只能读成「已配置」，同一包里 `EngineState` 句子把这件事说死。AC#2 的正控语义**长在形状里**：开关的事实与「没人消费它」的事实**同一包、永不分离**（TestPluginInventory215EnabledRowAndNotWiredSentenceRideOnePacket 钉死）。
- **交付 2（从菜单开关）未做，登记不化装**：需要 config 写路径＋锁定节三档审计，且判据 #3（跨档一致）对着一个不存在的引擎无法度量 ⇒ 不塞写字段、不加 C17 方法。第三条「未接线维度」记在本件。

## ③ 落盘清单

| 件 | 内容 |
|---|---|
| `internal/panel/plugins_inventory_215.go`（新） | 载体：`PluginRow{ID,Enabled,Capabilities,NetAllowlist,HostAPI}`＋`PluginInventorySection{Tier2Enabled,Rows,EngineState,SkillState}`＋两句 canonical 句子常量（`PluginEngineNotWired215`/`SkillsNotInConfigModel215`） |
| `internal/panel/composer.go` | `Snapshot.Plugins *PluginInventorySection json:"plugins,omitempty"`（第五枚 reader-gated 键；注释同款声明 Q-51 程序红） |
| `internal/panel/pump.go` | `PumpSources.Plugins func() PluginInventorySection`＋装配块（**拷贝后排序** by ID——map 顺序不许上线上，且泵不许写回读者的 slice；nil→`[]`；两句 verbatim）＋`sort` import |
| `internal/panel/plugins_inventory_215_test.go`（新） | 三枚测试＋三发突变牙 |
| 本件＋票 215 Progress log 一行 | 交付登记 |

## ④ 牙与突变读数（每发真红后还原，`git diff` 验零）

| 发 | 刀 | 红 |
|---|---|---|
| M1 | `sed` 摘掉 pump `if p.src.Plugins` 装配块（＋sort import） | `TestThePluginsKeyTravelsOnlyFromItsReader215`「a pump with a plugins reader sent no plugins key」＋另两枚 FAIL —— **摘掉 reader，键必须缺席**这半（有读者⇒键在）被钉 |
| M2 第一发 | sed 打挂编译（build failed） | **作废重打**（假红不算数，如实记） |
| M2 重打 | 排序比较器改恒 `false` | 「wire rows not sorted by ID: [zeta, alpha]」＋Publishes 测「not sorted」 —— 排序牙真红 |
| M3 | 无读者时泵自造 `&PluginInventorySection{}`（拆 reader-gating） | 「a pump with no plugins reader sent the key」＋「top-level keys」缺测 —— 键必须**只因读者而在**这半被钉 |
| 无读者基线（常驻，非突变） | `NewSnapshotPump(PumpSources{})` | 无 `"plugins"` 键、顶层恰 4 键 —— AC#2 正控的缺席半：**配置开着但没人接读者⇒线上看不见 plugins 键**，与 A408「文案在、控件不在」同一面镜子 |

AC#2 正控两向合计：塞「配置 true」而装配里没有读者 ⇒ 无键（M1＋常驻缺席半）；读者在 ⇒ `enabled:true` 与「零读者」句子**同包同行**（常驻正控枚）。

## ⑤ 门禁读数

| 门 | 读数 |
|---|---|
| gofumpt -l（点名 4 件，绝对路径 `/d/work/base/gopath/bin/gofumpt.exe`） | **空**，rc0 |
| `go vet ./internal/panel` | rc0（曾抓 `PluginRow` 含 `[]string` 不可 `!=` 比 ⇒ 测试改 `reflect.DeepEqual`，如实记） |
| `go build ./internal/panel` | rc0 |
| `go test ./internal/panel -count=1` | **恰 3 红，全数非本面**：①②`TestApprovalCardViewJSONKeysMatchFrontend`＋`TestComposerContractTypesMatchFrontend` 逐字 `Go Snapshot emits [plugins] that interface PanelSnapshot does not declare`（213 立的既定程序：键上货、panel.ts 声明归界面腿 Q-51）——★**起手后异动**：本腿起手时基线红为 `[commands modelSwitch ballStyle run]`＋`ComposerState files`，**收尾前界面腿的在飞件落地**（`frontend/src/lib/panel.ts` 未提交修改，Snapshot 与 PanelSnapshot **11 键 reconciled**），那批红转绿，recon 残红只剩本腿的 `[plugins]` 一词＝**本腿增量即这枚预定红**；③`TestC21DesignTokensFourWayAgree` 既存红，改前＝改后同点位同族（`c21-native-tokens.md:58-60/103-105` 点名 GlassRing/GlassHi/GlassHiSoft 不在 `internal/ball/tokens.go` 的 Light/Dark 调色板；`internal/ball/**` 另一波腿在飞、本腿零字节） |
| d22scan | **clean**（曾抓 1 枚 phantom-citation：注释里证据件路径被换行截断 ⇒ 路径解到一行＋本件落盘后复跑 0 findings；如实记） |
| LF | 4 件改动全 LF |
| unwired 旗 AC#5 | `internal/config/**` **零字节改动**（`git status` 无该目录条目）；guard 表 5 行、disposition 5 条**改前＝改后** |

## ⑥ 诚实边界（具名，不许被读成别的）

0. **commit 后读数漂移补录（2026-11-07 23:2x）**：commit `5341ead7a` 落盘后自跑，`internal/panel` 红从 3 枚收敛为 **1 枚**（只剩 C21，ball 族）——两枚 recon 转绿，原因＝**界面腿在飞件给 `plugins` 键补了 TS 声明**（`frontend/src/lib/panel.ts:372-384` `plugins?: PluginsSection`，未提交、引用本腿载体件）。即 recon 的红是**程序红、随界面腿落地即转绿**：程序两头都对了一次。本表 ⑤ 的 3 红读数是 commit 时刻的现量，此处为最终现量。

1. **界面点不动**：网页→Go 那一跳今天只有命令行路（账 `A410`）⇒ 本票终态＝**Go 侧可枚举＋命令行可执行＋界面看不到**；`plugins` 键的 TS 声明与渲染归界面腿（Q-51），recon 的红是既定程序不是缺陷。
2. **生产装配根一行未接**：`cmd/wisp/run.go` 的 `PumpSources{...}` 加 `Plugins: <从 config 读 [plugins] 节>` 一行归下一腿（本波 `cmd/wisp/**` 被 219 腿独占，本腿零字节）；在那行落地前，运行中的进程不发 `plugins` 键——这是范式钉的形状，不是遗漏。
3. **生效位没有**：`EngineState` 说的是「今天零读者」；引擎（50/51）落地的那条腿必须同步改 unwired.go disposition 行与这句句子——诚实句子不得活得比它描述的缺席更久（载体注释与句子常量注释均已具名此义务）。
4. **交付 2（开关）与 215b（技能）未做**：登记于 ②，等各自的票/批。
