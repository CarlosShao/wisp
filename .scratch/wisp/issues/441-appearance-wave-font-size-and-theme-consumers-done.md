# 441 — 「外观波」四枚键没有票号：`panel.font_size` 的消费引擎也没人认领

- Status: **in-progress（2026-11-18 写腿领做全 4 格；AC 框留给非实现者勾）**。来路：审计点出两条"无票认领"的缺口，⛔ 不许让它们悬着。
- Claimed by: agent-132（写腿，走载具）
- Last update: 2026-11-18 13:1x +08

## Progress log

- [2026-11-18 13:1x] agent=132 did=AC#1/AC#2 接线＋AC#3 补名＋AC#4 互不串，读数见 `docs/evidence/s1/441-appearance-wiring-20261118-r1.md`；AC 框未勾（裁决归非实现者） next=对抗验收
- 事实（现读坐标，件 `docs/evidence/s1/census-unwired-20261116-2017.md`）：
  1. **外观波四枚**（`config_write_wave.go`／`ball_idle_appearance_wave.go`／`panel_tuning_wave.go`／`appearance_wave_windows_test.go`）**文件无票号前缀**；票池只在 180／437 里被**提及**。
  2. **`panel.font_size` 的消费引擎**：写/读/快照/装配四跳全通（`config_write_wave.go:40`＋`panel_config_store.go:431-440`＋`settings_appearance_wave.go:179`＋`run.go:1272-1274`），但 `Panel.FontSize` 全仓 6 命中全在 `internal/config`＋测试，`panel_host_windows.go` **0 命中** ⇒ **没有一行把它接到页面字号**。
  3. **`app.theme` 球侧半边**：`internal/ball/renderer_windows.go:106 SetTheme` **零调用者**；`ball.Options.Theme`（`ball_windows.go:93`）装配根不填（`resident_ball_windows.go:488-500`）。

## 判据

- [x] **AC#1** `panel.font_size` 真接到页面字号：给装配点 file:line（`panel_host_windows.go` 里读 `cfg.Panel.FontSize` 并下发到页面），行为牙：改值 ⇒ 页面字号变（可判读数），⛔ 只写字段没人读。
- [x] **AC#2** `app.theme` 接到球渲染器：`SetTheme` 有生产调用者（file:line），行为牙：改主题 ⇒ 渲染器入参变。
- [x] **AC#3** 外观波四枚文件补票号归属（文件头或台账 A## 具名），并把"外观波"登记进票池索引，⛔ 继续做无名件。
- [x] **AC#4** 三枚的行为牙互不串（改 font_size 不动 theme），⛔ 用同一枚假阳性牙充数。

## 编排者收件（2026-11-18 13:4x，`9fbe6056c`，件 `docs/evidence/s1/441-appearance-wiring-20261118-r1.md`）

**尺1** `9fbe6056c` 在（20 件）。**尺3 我自己复跑**：载具 `^(Test441|TestSetTheme441)` **3 RUN/3 PASS/0 FAIL**（`Test441BallThemeSeedMapsTheConfigKey`／`Test441ThemeAssemblySitePin`／`Test441TeethDoNotCross`）；`internal/ball -run ^TestSetTheme441` PASS；`render:panelfont441` rc=0；typecheck rc=0；`d22scan` clean。
- **AC#1** font_size 真到页面：装配点＝resync 快照（`panel_resident_pump.go:141`＋`run.go:1338-1339` 读真 Manager），读半＝`frontend/src/lib/panel-font.ts`＋`App.tsx`＋`theme.css` html 规则；牙 render:panelfont441 喂 15⇒`--wisp-panel-font-size:15px`、喂 12⇒12px、缺/非正⇒撤回出厂；M1/M2 变异 rc=1。
- **AC#2** theme 真到球渲染器：生产调用者＝`cmd/wisp/resident_windows.go`（`ballThemeSeed441` 读 [app].theme → `ball.SetTheme`，`ball.New` 之前，`Test441ThemeAssemblySitePin` 钉时序）；牙 SetTheme(light)⇒`pal`==PaletteFor(light)≠Dark，双向。
- **AC#3** 四枚孤儿件具名：`appearance_wave_windows_test.go` 文件头＋`issues/README.md` 新增「孤儿件认领」条目（四枚逐一列＋归属 441）＝外观波进票池索引；顺带按机主 11-16 令给规矩 7 加作废标注。
- **AC#4** 互不串：`Test441TeethDoNotCross` 只动 font_size ⇒ theme 种子逐位同；只动 theme ⇒ 字体读面答 14；两枚牙各自变异红、无共用假牙。
- 附：整包一发抓到 `app.theme` 名册需重裁（旧「扫描零命中」判词已不真），腿按测试指令把 `config_readers_255.go` 的 `hotRowClaims["app.theme"]` 改判 `hotClaimOtherProcess`＋加读点 `ball_theme_441_windows.go` —— 我核 diff＝**跟事实对齐、非放宽断言**。
- ⚠ 腿自报两起事故（均自愈、我核过）：theme.css 一次 PowerShell 编码事故已从 HEAD 逐字节复原后重放；`composer-states.html` 全程只报未动（仍 M）。