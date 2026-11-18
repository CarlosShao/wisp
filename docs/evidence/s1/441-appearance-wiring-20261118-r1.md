# 441 — 外观波孤儿件认领＋两条消费接线（写腿交件，2026-11-18 r1）

写腿 agent-132，按派单 `.scratch/wisp/probes/orch/20261118/txt/441-appearance-wiring-dispatch.txt` 做票 441 全 4 格。
本件是**实现者侧交件**；AC 框留给非实现者勾（裁决表归对抗验收）。

## AC#1 — panel.font_size 真接到页面字号

**现成路勘定**（派单要求"走哪条现成路你查"）：装配根读真 Manager 喂快照的路**已经在**，不是 panel_host_windows.go——
写入/读取/快照/装配四跳在 7f3177f93 已通：

- `internal/config/schema.go:603`（panel.font_size，默认 12）＋ `internal/config/settings_appearance_wave.go:179`（PanelFontSizeOptions 生产读面）
- `cmd/wisp/panel_resident_pump.go:141`（常驻无管线形装配点）＋ `cmd/wisp/run.go:1338-1339`（管线形装配点）——两处都读 `mgr.PanelFontSizeOptions()` 喂 `PumpSources.PanelFontSize`
- `internal/panel/pump.go:556-558` 装包 → `internal/panel/composer.go:190` `panelFontSize` 键（resync-on-show 推给页面）

census-unwired 2026-11-16 点名的缺口＝**页面 0 命中（没人读）**。本腿补上的装配点（新增）：

- `frontend/src/lib/panel-font.ts`（新件）：`applyPanelFontSize` 把快照的 fontSize 钉上文档根（`--wisp-panel-font-size` inline ＋ `data-wisp-font-size`）；缺节/非正撤回出厂默认
- `frontend/src/App.tsx`（useEffect，消费 `snapshot.panelFontSize?.fontSize`）：每次快照推送幂等重跑
- `frontend/src/styles/theme.css`（html 规则）：声明出厂默认 12px，rem 基 = 配置值 × 4/3（12×4/3=16px＝浏览器默认）——**默认 12 与接线前排字逐位相同**（schema 注释的"wiring changes nothing at the default"兑现），写入值才移动整页 rem 尺度

**行为牙＝`render:panelfont441`**（`frontend/scripts/render-panel-font-441.tsx`，npm script `render:panelfont441`）：喂 15 ⇒ 文档根 `--wisp-panel-font-size: 15px` ＋ `data-wisp-font-size="15"`；喂 12 ⇒ 12px/12；缺节/0/负/NaN/undefined ⇒ 双双撤回；同值二推幂等。**反向判据（变异现跑）**：

- M1（把 App.tsx 的消费调用真删）：仪器 RED 1 条（"App must consume the snapshot's panelFontSize"）rc=1
- M2（把 theme.css 的 html 规则整块删）：仪器 RED 2 条（DECLARE/CONSUME 两钉）rc=1
- 两枚正向牙各配自己的变异红；M1 第一次试是把调用注释进注释里——仪器没红，这正是仓里"plant-in-comment"的老坑（probes/263 那族），已按真删重跑
- 配置屏徽章重审（`config-screen.tsx`＋`render-config-screen` 重判）：font_size 行撤 `consumerMissing`、徽章回到档位词 `关窗/重启生效`（fed 计数 关窗/重启生效 2→3、缺消费引擎 1→0）；`缺消费引擎` 的 rowBadge 机制与 BADGE_WORDS 留作通用机制（今天零行使用，两个 render 都钉 0）

## AC#2 — app.theme 接到球渲染器

**SetTheme 的生产调用者（file:line）＝`cmd/wisp/resident_windows.go`**（装配根，ballStyleSeed65 → `ball.SetStyle` 的同一座位）：

```
ballTheme441, ballThemeOK441 := ballThemeSeed441(rt.Layout.DataDir)   // 读 [app].theme（每次 boot 一发 LoadFile）
if ballThemeOK441 { ball.SetTheme(ballTheme441) }                     // ← 生产调用者
```

- 种子 `cmd/wisp/ball_theme_441_windows.go`（新件）：light→`ball.ThemeLight`、dark→`ball.ThemeDark`（都 ok=true，默认值也算文件的回答）；auto 与读不了的文件 ok=false **不冒充**（Warn 具名：auto 的 system-theme follow 是 D29 另一笔）
- 时序钉死：seed → SetTheme → `startResidentBall(rt.Registry…)`（ball.New）严格递增（`Test441ThemeAssemblySitePin` parse 钉）——都在首帧之前，"caller must trigger a repaint afterwards" 由预建时序免掉
- **行为牙＝`TestSetTheme441SwitchesTheActivePalette`**（`internal/ball/theme_palette_441_windows_test.go`）：`SetTheme(ThemeLight)` ⇒ 包级 ACTIVE `pal` == `PaletteFor(ThemeLight)` 且 ≠ DarkPalette；反向 dark 同证（渲染器每一笔都读 `pal`，palette 真变）
- 链路：cfg（light）→ 种子映射 → SetTheme → pal，三跳各钉（种子映射＝`Test441BallThemeSeedMapsTheConfigKey`）
- `Options.Theme`（ball_windows.go:93）保持未填不动——本票的机制只有 SetTheme 一条路，不造第二条

## AC#3 — 四枚无名件补归属＋登记票池索引

四枚（出自外观波本体 `7f3177f93`；此前只在 180/437 被提及）：

| 件 | 具名方式 |
|---|---|
| `cmd/wisp/appearance_wave_windows_test.go` | **文件头补名段**（票 441 AC#3 归属注释，写面内） |
| `internal/panel/config_write_wave.go` | **台账具名**——该腿写面不含 `internal/panel/**`，按票面判据"文件头注释**或台账具名**"走台账 |
| `internal/panel/panel_tuning_wave.go` | 同上 |
| `internal/panel/ball_idle_appearance_wave.go` | 同上 |

台账＝`.scratch/wisp/issues/README.md` 2026-11-18 追加的「孤儿件认领」条目（四枚逐一列出＋归属票 441）；同条目即"外观波"的**票池索引登记**。另按机主 2026-11-16 全局令（铁律 4）在 README 规矩 7 顶部加了作废标注（「前端归外部 agent」原话只留档不再执行）。

## AC#4 — 三枚牙互不串

- 牙一（font）：`render:panelfont441`（页面根载荷）；牙二（theme）：`TestSetTheme441SwitchesTheActivePalette`＋种子映射测试；互不串证明＝`Test441TeethDoNotCross`（`cmd/wisp/appearance_theme_441_windows_test.go`）：
  - 只动 font_size 的两份文件（12 vs 15）⇒ 种子 theme 答案逐位相同（ThemeDark==ThemeDark）
  - 只动 theme 的两份文件（dark vs light，font 同 14）⇒ `panelFontSizeSectionFromConfig` 双双答 14
  - 热写半：走真门 `mgr.SetPanelFontSize(15)` 落盘（written=[panel.font_size]）后重跑种子 ⇒ theme 仍 (ThemeLight, true)
- 变异各红各的：前端 M1/M2 只红 font 牙；Go 侧 SetTheme 调用拔除由 `Test441ThemeAssemblySitePin`（静态钉）红、`pal` 断言由 `TestSetTheme441SwitchesTheActivePalette` 红——没有同一枚假阳性牙充数
- ⛔ 未动 theme.css 的 `[data-theme]` 面、未动页面主题行机制（页面跟随宿主主题的引擎仍未接线，行内如实具名——球侧已接线、页半未接，两半分开说）

## 门禁读数（全部现跑）

| 门 | 读数 |
|---|---|
| 载具一发①（`'^Test441'`，跑前 `Get-Process wisp.test`=空） | rc=0，Test441BallThemeSeedMapsTheConfigKey / Test441ThemeAssemblySitePin / Test441TeethDoNotCross 全 PASS |
| 载具整包 `go test ./cmd/wisp`（一次一发，1217s） | 4 红：①TestTicket255RosterStillMatchesTheActualReadSites（**本腿造成，扫描当场抓到新读点**，重裁后绿，见下）②TestAC1ResidentLegBooksItsShutdownBeforeClosingTheSink（子进程 shutdown 处 cgo 0xc0000005 崩——2026-11-16 三发载具全 PASS，本日首发红；单独重跑 PASS 9.99s ⇒ 环境偶发，非本腿）③Test313WritableUnionEqualsTheTwoRosters（存量红，7f3177f93 交件已具名"union pin predates 187/65 rosters"，今为票 39 的 app_language 继续）④TestAC4FocusReturnToPriorWindowGap33r5（存量环境红，43bf68b51 交件具名） |
| 载具一发②（`'^TestTicket255RosterStillMatchesTheActualReadSites|^Test441|^TestAC1ResidentLegBooksItsShutdownBeforeClosingTheSink$'`） | rc=0 全 PASS——含重裁后的 roster 测试与 AC1 的单测复跑 |
| `go build ./...` | rc=0 |
| `go vet ./cmd/wisp ./internal/ball` | rc=0 |
| `gofmt -l`（五枚改动 Go 件） | 空 |
| `go test ./internal/ball -run TestSetTheme441` | rc=0 PASS（0.00s，纯函数不开窗） |
| `sh scripts/d22scan.sh` | clean，0 findings（见回话首行末行） |
| 前端 build / typecheck / lint | rc=0 / rc=0 / rc=0（45 warnings 全为存量，新件 0） |
| `render:panelfont441` | rc=0（M1 真删 rc=1、M2 整块删 rc=1，均恢复后复绿） |
| `render-config-screen`（重判后） | rc=0 |
| 回归：render:receiver / shell77 / cfggolden / composer / palette38 / nav | 全 rc=0 |

## 整包载具抓到的第 5 件事：255 名册重裁（本腿造成、按测试指令断案）

整包 run 里 `TestTicket255RosterStillMatchesTheActualReadSites` 红——**正确的红**：`[app]` 节的读者扫描扫到了新件
`cmd/wisp/ball_theme_441_windows.go`（本腿的 app.theme 种子），而 `config_readers_255.go` 的名册没重裁。按测试自己的
判词（"something moved and cmd/wisp/config_readers_255.go has not been re-adjudicated"）重裁两处：

- `sectionReadSites["App"]` 增补 `cmd/wisp/ball_theme_441_windows.go`（autostart/single_instance 排序保持）；
- `hotRowClaims["app.theme"]` 从 `hotClaimNoReader`（"扫描零命中"）重裁为 `hotClaimOtherProcess`：球侧消费已接
  （boot 种子读 [app].theme → 装配根调 ball.SetTheme，ball.New 之前），种子是每次 boot 一发 LoadFile、不读 live
  Manager ⇒ 已保存主题随下一次进程启动落到球配色，hot 档生效节奏如实答 other-process；引证 token
  `switch c.App.Theme`（该文件唯一命中，EvidenceLines 尺要求）；
- 重跑：roster 测试 PASS（一发② rc=0）。"扫描零命中"旧判词在本行内随接线作废，句子里写明了作废原因。

## ⚠ 交底（诚实半）

1. **theme.css 一次事故与复原**：变异 M2 用 PowerShell 的 `Get-Content -Raw`/`Set-Content`（默认 GBK）种进去，把文件编码打坏（首字节从 `/* ==` 变成乱码）。该文件在本腿动手前与 HEAD 逐位一致（开工时 `git status` 里只有他人在飞的 `frontend/fixtures/composer-states.html` 是 M），故经 `git show HEAD:frontend/src/styles/theme.css > 同路径` 复原到 HEAD 字节（`od -c` 核过首字节），再用编辑工具重放本腿那 13 行——最终 `git diff` 只含那 13 行（本回话上方已贴）。⛔ 没碰 `frontend/fixtures/composer-states.html`（全程保持 M，只报不动）。
2. `internal/panel/**` 三枚文件只走台账具名、一行未改（写面纪律）；`.scratch/wisp/issues/README.md` 与票面文件的改动是 AC#3 与票池规矩明确要求的落点，随本腿一并提交。
3. `panel.font_size` 的生效节奏＝下一次 resync 推送（下一次面板显示），与写回执 Note（"关窗再开即跟上"）同一时态，没有谎报"立即"。
