# 305-a2 · 10 表①：窗口 `416d9d56..f718e9b6` 候选笔名册（37 枚）

## 0. 尺与射程（每把都自落 rc；raw 在 `logs/`）

| 尺 | 命令（逐字） | rc | 件 | 枚数 |
|---|---|---|---|---|
| S1 窗口总枚数 | `git rev-list --count 416d9d56..f718e9b6` | 0 | `logs/10-revlist-count-total.txt` | **963**（与任务书一致） |
| S2 名册尺 | `git log --oneline 416d9d56..f718e9b6 -- cmd/wisp internal/panel` | 0 | `logs/10-roster-oneline.txt` | **37 行＝37 枚**（射程＝仅 `cmd/wisp`＋`internal/panel` 两目录，含测试文件，不含注释行计数） |
| S3 逐笔文件尺 | `git log --format='C %h' --name-only 416d9d56..f718e9b6 -- cmd/wisp internal/panel` | 0 | `logs/11-roster-names.txt` | 37 枚的 name-only 全量 |
| S4 ★承重尺（定罪／脱罪都靠它） | `git grep -l -E 'SetHtml\|wispDispatch\|currentWindow\|bringUp\|firstRoundTrip\|serveEntry' 416d9d56 -- cmd internal` ＋ HEAD 同尺 `grep -rln --include=*.go …\| coldStartPageHandover …' cmd internal` | 0 | `logs/19-symbol-census-two-ends.txt`＋`logs/15-extra-rulers.txt` | 绿点：**5 枚文件，全在 `cmd/wisp`**（`panel_host_windows.go`／`panel_resident_windows.go`／`panel_resident_windows_test.go`／`panel_host_windows_test.go`／`panel_host_windows_live_test.go`）；HEAD：同样只有 `cmd/wisp/**`（多出的都是 `_test.go`，其中 `config_readers_255.go` 命中的是三行**注释**，见 §3 附注） |
| S5 窗口内碰这 5 枚文件的笔 | `git log --format='C %h %ad %s' --date=short --name-only 416d9d56..f718e9b6 -- <那 5 枚路径>` | 0 | `logs/28-final-rulers.txt` | **4 枚**：`cf95c799`／`4658dbb6`／`7a0236b3`／`f7d28ef0` |
| S6 窗口内碰 `frontend/**` 的笔 | `git log --oneline 416d9d56..f718e9b6 -- frontend \| wc -l` | 0 | `logs/28-final-rulers.txt` | **0 枚** ⇒ 入口页与 embed 图案在 git 这一侧⛔ 移动过（未跟踪的 `frontend/dist/index.html` 另算，见 30 号件的台面轴） |
| S7 `internal/panel/assets.go`（`serveEntry` 的 `Built()`/`Resolve()` 的家） | S3 的 name-only 逐枚读 | 0 | `logs/11-roster-names.txt` | 窗口**⛔ 任何笔碰过它**（`internal/panel` 被碰的文件全集＝git.go／git_test.go／instructions_200.go／workspace.go／workspace_test.go／workspace_account_181r3_test.go／bridge.go／composer.go／composer_dispatch.go／composer_dispatch_test.go／config_handlers.go／config_route_248_test.go／l2_grant_boundary_test.go／pump.go／subagent_blocked_220_test.go／subagent_roster_197.go／inbound_roster_253_test.go） |
| S8 名册尺⛔ 覆盖的目录 | `git log --oneline 416d9d56..f718e9b6 -- internal/agent internal/tools internal/ball internal/config internal/llm internal/risk \| wc -l` | 0 | `logs/28-final-rulers.txt` | **46 枚**在 S2 射程**外**。⚠ 诚实记：S2 那把尺⛔ 是全窗口尺；这 46 枚靠 S4 脱罪（S4 在**绿点树**上跑，符号只在 `cmd/wisp` 那 5 枚文件里 ⇒ 不碰这 5 枚文件的笔⛔ 能改这条链的产码行为）。另 S6/S7 覆盖 `frontend/**` 与 assets。 |

★**分级判据（写死，⛔ 靠 commit message 叙事）**：一枚笔要能翻转这两枚症状，必须碰 S4 那 5 枚文件之一（产码面）或碰 `internal/panel` 的入向门（`bridge.go`/`composer*.go`，报告信封面）。⚠ 一条独立脱罪内容锚：`startPanelForTest` 交回来的 `mgr.disp.Mode` 被测试自己的 `recordingModeHandler` 占着（`panel_resident_windows_test.go:320` 与 `:861` 两处逐字 `h := mgr.disp.Mode.(*recordingModeHandler)`）⇒ `internal/panel` 里那些**模式处理器**（workspace.go／instructions_200.go／git.go）在这两枚用例里⛔ 在场。

## 1. 名册（37 枚，`git log` 图序，新→旧；日期列＝S2 的 `%ad --date=short` 作者日，关键四枚另给 `%cd` 时刻）

路径缩写：`cw/`＝`cmd/wisp/`，`ip/`＝`internal/panel/`。

| # | sha | 日期 | 该笔在这两枚症状相关的文件 | 它到底改了什么语义（内容锚，⛔ 行号当真相） | 分级 |
|---|---|---|---|---|---|
| 1 | `8d30a862` | 10-07 | `cw/config_reload_223_test.go` | 台账/验收文案笔，测试文件内注释与名册；S4 符号⛔ 命中该文件 | 与这两条路径无关 |
| 2 | `3d9b8374` | 10-06 | `cw/config_reload_223_test.go` | 232-r2 代提；同上 | 无关 |
| 3 | `42e89896` | 10-06 | `cw/config_sentences_223r2_test.go` | 231 AC#3 注释措辞；同上 | 无关 |
| 4 | `c5040ea7` | 10-06 | `cw/config_reload.go` | 注释⛔ 抄既有原句；文件不在 S4 五枚里 | 无关 |
| 5 | `a16d1ff7` | 10-06 | `cw/config_reload_223_test.go`＋2 枚 223 测试 | 常驻钉加一行、`cause=newer-build` 进两张名册；⛔ 面板文档面 | 无关 |
| 6 | `290dca87` | 10-06 | `cw/config_reload.go` | 给"版本更高"一条出口（`cause=newer-build`）；配置重载路，⛔ SetHtml/Eval | 无关 |
| 7 | `9a941965` | 10-05 | `cw/resident_approval_windows.go`＋1 枚 268 测试 | 常驻腿把"配置缺失/被拒"分两枚 provenance；审批门路 | 无关 |
| 8 | `16901acb` | 10-05 | `ip/git.go`、`ip/workspace.go`、`ip/instructions_200.go`＋3 枚测试 | 改写账户生产者换成 C26 读数——**全在模式处理器里**，而这两枚用例的 Mode 被测试自装的 `recordingModeHandler` 占着（`:320`/`:861` 逐字） | 无关 |
| 9 | `7c644bb0` | 10-05 | `cw/firstrun_257_nonpreset_test.go` | 257-r2b 回执文案＋同包测试 | 无关 |
| 10 | `3f0c4fff` | 10-05 | `cw/firstrun.go`＋1 枚测试 | 268-a2 复核：err 透传；首跑路 | 无关 |
| 11 | `32e74479` | 10-05 | `cw/panel_pump_test.go`＋4 枚审批/装配测试 | 267-r2 种子迁移：把带外 `confirm_timeout_sec` 抬进带内；`panel_pump_test.go:112` 那枚 `executeOn145` 的 30s ctx 抬到 45s。**只动测试种子与超时**，⛔ 产码 | 可能不能（**仅整包序**：它改的是同包另一枚用例的时长；三枚具名隔离跑时它⛔ 在场） |
| 12 | `7883f855` | 10-04 | `cw/resident_approval_windows.go` | 265-r1c 门禁读数笔（先修 `internal/panel/pump.go` 的 d22scan 注释——注释，⛔ 语义） | 无关 |
| 13 | `a04a095f` | 10-04 | `cw/resident_approval_windows.go`、`cw/resident_task_source_windows.go`＋2 枚测试 | 265-r1 遗留产码收编；审批/任务源 | 无关 |
| 14 | `79c579e2` | 10-04 | `cw/resident_approval_windows.go`＋1 枚测试 | 260-r4 卡片通道标签读取时求值 | 无关 |
| 15 | `170e0459` | 10-04 | `cw/resident_approval_windows.go` | 260-r3 slog 键名放回句子 | 无关 |
| 16 | `2ae018be` | 10-04 | `cw/resident_approval_windows.go`＋2 枚测试 | 260-r3 Esc 文案说实话 | 无关 |
| 17 | `2a10134e` | 10-04 | `cw/resident_approval_windows.go`、`cw/resident_windows.go`＋1 枚测试 | 256-r1 审批门改吃 [risk] 两枚配置 | 无关 |
| 18 | `cbece45e` | 10-04 | `cw/config_readers_255.go`、`cw/resident_ball_windows.go`、`cw/resident_windows.go`＋3 枚测试 | 258-r2 档位词跟集合说真话＋两枚 winlive 尺修形状。`config_readers_255.go` 里命中 S4 符号的只有 `:18`/`:22`/`:158` 三行**注释**（逐字含 "create call inside bringUp spelled its geometry as the literals"） | 无关 |
| 19 | `d9bce7ce` | 10-03 | `cw/resident_hotkey_v1probe_test.go` | 258-v1 残件入库 | 无关 |
| 20 | `bd5c049f` | 10-03 | 2 枚 258 winlive 测试 | 258-v1 验收补 commit（测试侧） | 无关 |
| 21 | `5e8748b3` | 10-03 | `cw/config_readers_255.go`、`cw/models.go`、`cw/resident_ball_windows.go`、`cw/resident_hotkey_258_seam_windows.go`、`cw/resident_windows.go`＋4 枚测试 | 212-r1＋258-r1 代笔收尾；热键缝与档位读数，⛔ 文档交接面 | 无关 |
| 22 | `9bab468f` | 10-03 | `cw/firstrun.go`、`cw/run.go`系 7 枚测试 | 261-r2 fixture 补 `enabled` 键 | 无关 |
| 23 | `cf95c799` | 10-03 | `cw/panel_host_windows_test.go`（S4 五枚之一） | **⛔ 产码**，但它换了同包一枚真窗用例的形状：`TestPanelHostRealWindowHopAndLifecycle` 里 `before := readWebviewTree(...)` 换成 `before, beforeExtra := settleTreeReading(t, self)`（新增 `treeSettleWait = 3 * time.Second`/`treeSettleStep = 100 * time.Millisecond` 的**采样等待**），并把"树内计数等值"判据换成 `browserHostPids` 宿主 pid 同一尺（新增 `mineHosts`/`baselineHosts` 豁免）。⇒ 同一次 `go test` 进程里真窗用例的**时长与浏览器复用窗口**变了 | 可能不能（**仅整包序**；三枚具名隔离跑时该用例⛔ 在场） |
| 24 | `4658dbb6` | 10-03 13:58（`%cd`） | `cw/panel_host_windows.go`＋`cw/panel_resident_windows.go`＋`cw/config_readers_255.go`＋3 枚几何测试 | **⛔ `bringUp` 的建窗那一跳被改**：`webview2.NewWithOptions` 的 `WindowOptions: webview2.WindowOptions{Title: panelTitle, Width: 420, Height: 260}` 被删，换成 `WindowOptions: m.windowOptions()`；新增 `panelHostOption`/`withGeometrySource`/`PanelManager.geometry func() (width, height int)`，`windowOptions()` 里**每次建窗现调 `m.geometry()`**；装配根 `newResidentPanelManager` 交下来的 `panelGeometrySource(dataDir)` 逐字 `cfg, _, err := config.LoadFile(cfgPath, nil)` ⇒ **面板线程上、`bringUp` 之内多了一次磁盘读＋TOML 解析**（默认档答 0/0 ⇒ 回落到常量 420×260）。笔自身标注"编排者代提·死腿半成品·未验证" | ★**能**（唯一在窗口里动过"建窗那一刻"的产码笔；几何数值可读码判为等值，**新增的磁盘/解析那一跳⛔ 等值**） |
| 25 | `67ab595d` | 10-03 | `cw/config_reload.go`＋1 枚测试 | 255-r2 `restartTierKeys` 那句有测试兜住 | 无关 |
| 26 | `480b970d` | 10-03 | `ip/pump.go`、`ip/subagent_roster_197.go`＋1 枚测试 | 220-r1 名册改成 corr 与 taskID 双查。`ip/pump.go` 是**快照泵**（`type SnapshotPump`/`func (p *SnapshotPump) Snapshot()/Marshal()/Publish()`），⛔ 引用 WebView/SetHtml/Eval；S4 两枚端点都⛔ 命中它 | 可能不能 ⇒ 读码后**判为无关**（理由＝它是出向快照组包，这两枚用例的判据走入向门＋`Eval`；分级保守列此，供裁决者核我这条脱罪链） |
| 27 | `13dd60b4` | 10-03 | `cw/config_receipt_255_test.go` | 255 r2 证据件终值 | 无关 |
| 28 | `17ff058d` | 10-03 | `cw/config_receipt_255_test.go` | gofumpt 格式化 | 无关 |
| 29 | `87bc6aca` | 10-03 | `ip/inbound_roster_253_test.go` | 253-r5 形ⓐ 另立名册尺，⛔ 改 `bridge.go` 一字 | 无关 |
| 30 | `a4906b6c` | 10-03 | `cw/config_readers_255.go`、`cw/config_receipt_255_test.go`、`cw/config_reload.go` | 回执加方括号＋句子装配确定性钉 | 无关 |
| 31 | `ae60a87c` | 10-03 | 同 30 ＋`cw/config_reload_223_test.go` | AC#1 回执只说真话＋[app] 键级补尺 | 无关 |
| 32 | `fa5593c0` | 10-02 | `cw/firstrun.go`＋1 枚测试 | 198-r2 四格落地 | 无关 |
| 33 | `613606c0` | 10-02 | `cw/firstrun.go`、`cw/run.go`＋2 枚测试 | 198-r1 首建 config.toml | 无关 |
| 34 | `b644d310` | 10-02 | `cw/panel_config_248_test.go` | 248-r1 死腿收尾代提（测试路径两枚） | 无关 |
| 35 | `0d87a681` | 10-02 09:03（`%cd`） | `cw/panel_inbound.go`、`cw/panel_config_store.go`、`cw/run.go`、`ip/bridge.go`、`ip/composer.go`、`ip/composer_dispatch.go`、`ip/config_handlers.go`、`ip/pump.go`＋3 枚测试 | **动了入向门的白名单**：`ip/bridge.go` 在 const 块里新增 `MethodConfigGet`/`MethodConfigSet`（逐字注释"They are the two spellings the product spec already carries"），`newComposerDispatchChain` 新增 `Config: configWrites` 槽，`cw/panel_inbound.go` 的 CLI 支新增 `if reply != "" { fmt.Fprintf(s.stdout, …) }`。⇒ 这三枚都⛔ 在 `wispDispatch` 绑定闭包**之下**、报告信封**之内**；若守卫把测试的 `reportJSEnv` 信封拒了，症状应是 `awaitReport` 超时（"no report … within 15s"），⛔ 今天这两枚的红句 | 可能不能（唯一在窗口里动过"页面报回那一跳"的产码笔；读码推它**⛔ 造成本红**，凭据＝nail1 今天绿着 ⇒ 门是通的） |
| 36 | `7a0236b3` | 10-01 18:04（`%cd`） | `cw/panel_host_windows.go`＋`cw/panel_resident_windows.go`＋`cw/panel_resident_windows_test.go`（S4 五枚里占 3 枚） | 两截：ⓐ产码＝`var errPanelRefusedThread = errors.New(...)`，`bringUp` 的拒绝改成包这枚哨兵（`return fmt.Errorf("%w: …", errPanelRefusedThread, closedHwnd)`）；`panel_resident_windows.go` 把 `RequestShow` 的闭包换成 `rp.post(func() { rp.showOnThread(via) })`，`showOnThread` **认哨兵 ⇒ 把常驻面板线程判死**（走 `setStartUp` 那枚已有通道）。ⓑ测试＝在 `panel_resident_windows_test.go` **插进 179 行新用例**（hunk 头 `@@ -589,0 +591,179 @@`，位置在 AC13（`:313`）与 nail2（`:855`）之间）：`TestAC13NamedRefusalEndsThePanelThreadInsteadOfRefusingForever` 等，⛔ 开真窗、`PostMessageW` 种 WM_CLOSE、`refusalSettleBudget = 5 * time.Second` | ★**能**（产码那截改了"拒绝之后线程活⛔ 活"；测试那截在同文件里插进 179 枚会开真窗的行，直接改整包序） |
| 37 | `f7d28ef0` | 10-01 16:20（`%cd`） | `cw/panel_host_windows.go`＋`cw/panel_resident_windows_test.go` | 产码＝`bringUp` 在建窗之前新增**第二道闸**：`drainStaleQuitBeforeCreate()` 之后加 `if dirty, closedHwnd := staleCloseQueued(); dirty { return fmt.Errorf("panel host: refusing to create the panel window: …undispatched WM_CLOSE…", closedHwnd) }`，新增 `staleCloseQueued()`（`PeekMessageW(…, pnlWmClose, pnlWmClose, 0)` PM_NOREMOVE）与 `pnlWmClose = 0x10`。测试＝同文件加 175 行 owner 侧清理尺（`pumpThreadToQuiet`/`threadQueueHead`/`releaseThreadClean`/`threadWindowCount`），并把 `TestAC13BringUpSurvivesAReusedThreadQuit` 改成"还池前必须把线程派空" | ★**能**（绿点 `416d9d56` 之后**第一枚**动 `bringUp` 的笔；它给建窗加了"这道线程⛔ 干净就整枚拒绝"的新出口——拒绝那支会令 `Show` 失败，从而让"最后一份文档"停在探针上） |

★**头号嫌疑（各一枚，带 diff 证据）**
- **症状①（AC13 入口没交接上来）＝`f7d28ef0`**。凭据（读到的 diff 行，⛔ 叙事）：`bringUp` 里 `drainStaleQuitBeforeCreate()` 之后**新增**（逐字两行）：`+ if dirty, closedHwnd := staleCloseQueued(); dirty {` / `+ return fmt.Errorf("panel host: refusing to create the panel window: …", closedHwnd)`，再加函数 `staleCloseQueued()`（体内 `pnlPeekMessageW.Call(…, pnlWmClose, pnlWmClose, 0)`＝PM_NOREMOVE）与常量 `pnlWmClose = 0x10`。这条出口在绿点**⛔ 存在**；它一旦命中，`bringUp` 在建窗前就 return，`coldStartPageHandover` 那两枚 `SetHtml`（`:486`/`:467`）永远⛔ 跑，活文档就⛔ 是探针——但红句要求"探针在场**且**线程此前建过窗"，这只在整包序里成立 ⇒ **读码推到，⛔ 判死**。
- **症状②（nail2 `title=""`）＝`7a0236b3`**（`showOnThread` 认哨兵判死线程）与 `f7d28ef0` 的闸同源；另一枚**形状不同**的候选＝`4658dbb6`（`bringUp` 里新增磁盘读）。⇒ 这一枚的读数最硬（nail2 的 `""` 逐字＝探针那份 136 B 文档**没有 `<title>`**，见 20 号件 §3），所以判死那一发**优先在 clone 面上跑 nail2**。

## 2. 排除名册（**30 枚**＝上面 §1 名册里分级列为"无关"的全部行，行数与"排除 30 枚"逐字相等）

`8d30a862`、`3d9b8374`、`42e89896`、`c5040ea7`、`a16d1ff7`、`290dca87`、`9a941965`、`16901acb`、`7c644bb0`、`3f0c4fff`、`7883f855`、`a04a095f`、`79c579e2`、`170e0459`、`2ae018be`、`2a10134e`、`cbece45e`、`d9bce7ce`、`bd5c049f`、`5e8748b3`、`9bab468f`、`67ab595d`、`13dd60b4`、`17ff058d`、`87bc6aca`、`a4906b6c`、`ae60a87c`、`fa5593c0`、`613606c0`、`b644d310`。

**这一格"无关"的尺＝S4＋S5 两把并排**（⛔ 我逐枚读它们的 diff）：这 30 枚的 `--name-only` 清单里⛔ 一枚命中 S4 那 5 枚文件（尺＝S3，件 `logs/11-roster-names.txt`），而 S4 在**绿点树**上证明这条链的全部符号只住在那 5 枚文件里。⇒ **30＝37−4（S5）−3（另列嫌疑/可能：`480b970d`、`32e74479`、`0d87a681` 走 S4 之外的门与时长面，单列"可能不能"）**。⚠ 这句里的 3 枚与 4 枚不重叠：`cf95c799` 属 S5 的 4 枚（列"可能不能"），`4658dbb6`/`7a0236b3`/`f7d28ef0` 属 S5 的 4 枚且列"★能"。

**⛔ 这张表⛔ 等于"已排除"**：它只排除"产码面直接改这两条路径"这一种可能，⛔ 排除"整包序/线程污染/时序"这一族（那族只能靠跑，见 30 号件）。
