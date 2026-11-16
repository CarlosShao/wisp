# 441 — 「外观波」四枚键没有票号：`panel.font_size` 的消费引擎也没人认领

- Status: **open（2026-11-16 20:3x 编排者按只读审计腿 `census-unwired` 立票）**。来路：审计点出两条"无票认领"的缺口，⛔ 不许让它们悬着。
- 事实（现读坐标，件 `docs/evidence/s1/census-unwired-20261116-2017.md`）：
  1. **外观波四枚**（`config_write_wave.go`／`ball_idle_appearance_wave.go`／`panel_tuning_wave.go`／`appearance_wave_windows_test.go`）**文件无票号前缀**；票池只在 180／437 里被**提及**。
  2. **`panel.font_size` 的消费引擎**：写/读/快照/装配四跳全通（`config_write_wave.go:40`＋`panel_config_store.go:431-440`＋`settings_appearance_wave.go:179`＋`run.go:1272-1274`），但 `Panel.FontSize` 全仓 6 命中全在 `internal/config`＋测试，`panel_host_windows.go` **0 命中** ⇒ **没有一行把它接到页面字号**。
  3. **`app.theme` 球侧半边**：`internal/ball/renderer_windows.go:106 SetTheme` **零调用者**；`ball.Options.Theme`（`ball_windows.go:93`）装配根不填（`resident_ball_windows.go:488-500`）。

## 判据

- [ ] **AC#1** `panel.font_size` 真接到页面字号：给装配点 file:line（`panel_host_windows.go` 里读 `cfg.Panel.FontSize` 并下发到页面），行为牙：改值 ⇒ 页面字号变（可判读数），⛔ 只写字段没人读。
- [ ] **AC#2** `app.theme` 接到球渲染器：`SetTheme` 有生产调用者（file:line），行为牙：改主题 ⇒ 渲染器入参变。
- [ ] **AC#3** 外观波四枚文件补票号归属（文件头或台账 A## 具名），并把"外观波"登记进票池索引，⛔ 继续做无名件。
- [ ] **AC#4** 三枚的行为牙互不串（改 font_size 不动 theme），⛔ 用同一枚假阳性牙充数。
