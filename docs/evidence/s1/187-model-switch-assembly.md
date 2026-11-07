# 187 — 装配根最后一枚缺口：`PumpSources.ModelSwitch` 注入（145 腿 run.go 写面 · 本腿交付）

- 日期：2026-11-07 21:2x–21:4x +08
- 腿：Wisp 编队写腿（写面 `cmd/wisp/**` 本波独占）
- 起点 HEAD：`7d1cd7139`（A1621 wave76b）

## 0. 缺口是什么（起手三核现跑，锚 `7d1cd7139`）

| # | 核 | 命令（现跑） | 读数 |
|---|---|---|---|
| 1 | 装配根没接 | `grep -rn ModelSwitch --include=*.go cmd/wisp/` | **零命中**（缺口仍在） |
| 2 | 读口在位 | `internal/panel/pump.go:231` | `ModelSwitch func() ModelSwitchSection`；填段 `:432-446`（reader-gated，nil 表归一 `[]`） |
| 3 | 载件在位 | `internal/panel/model_switch_187.go` | 五字段：`Current`/`CurrentThinking`/`Models`/`ThinkingLevels`/`Reason` |
| 4 | 同款两根已接 | `cmd/wisp/run.go:1049/:1059` | `BallStyle: func() panel.BallStyleSection {` 与 `Run: rt.runFeed.Section,` 都在 ⇒ 只有 ModelSwitch 缺 |
| 5 | 无在飞件 | `git status --porcelain \| grep -v '^??'` ＋ `git diff --stat` | **空** ⇒ 红绿非别人在飞件造成 |

## 1. reader 出处（读码定，非猜）

- `internal/panel/pump.go:224` 字段注释逐字点名：*"read off the live config by the assembly root (internal/config's **ChatSwitchOptions** is the production reader)"*。
- `internal/config/settings_187.go:166` `func (m *Manager) ChatSwitchOptions() ChatSwitchOptions` —— 返回 `CurrentProvider/CurrentModel/CurrentThinking/Models/ThinkingLevels/Reason`，与 `panel.ModelSwitchSection` 一一对应。
- 与 `BallStyleOptions`（`internal/config/settings_65.go:86`）同 Manager 同款 ⇒ run.go 字面量照 `BallStyle: func() ...` 同形状落闭包，读**本 run 自己的 `mgr`**（= settings 腿 `rt.settings` 写经的同一 Manager，picker 显示的＝下一次 ApplySetting 会作用的）。

## 2. 落地

- **`cmd/wisp/run.go`**（唯一生产码改动，插在 `Run: rt.runFeed.Section,` 与 `Out` 之间）：
  `ModelSwitch: func() panel.ModelSwitchSection { v := mgr.ChatSwitchOptions(); ... }` —— 逐字段 verbatim 映射（Models 逐行转 `panel.ModelOption`），本字面量不造 provider 也不造 tier。摘掉这根线＝wire 无 `modelSwitch` 键（213 reader-gated 同款）。
- **新牙 `cmd/wisp/panel_model_switch_187_test.go`** 三发：
  1. `Test187AssemblyRootCarriesTheModelSwitchReader` —— 65 腿字节钉同款：run.go 字面量不再带那行 ⇒ 当场红；
  2. `Test187WiredModelSwitchReaderSendsTheLiveValues` —— 真 Manager（`thinking187Base` 同一枚最小 187 配置）经真 `ChatSwitchOptions` → 生产泵 marshal → wire bytes 里有 `modelSwitch` 键且五字段**硬编码**对（`deepseek/deepseek-chat`、`off`、models 恰一枚、levels 恰 `off|low`、Reason 空）；并核对 wire 与 `ChatSwitchOptions()` 直读一致（picker 不说谎）；
  3. `Test187MutationReaderRemovedSendsNoModelSwitchKeyAtAll` —— **突变**：`PumpSources{}`（摘注入）⇒ wire bytes 无 `"modelSwitch"` 键（145 hop-2 突变同款，断的是面板收到的字节）。

## 3. 门禁读数（全为 HEAD `7d1cd7139` 工作树真跑；env＝`unset GOTMPDIR TMPDIR TEMP TMP && export PATH=$PWD/build:$PWD/third_party/sherpa-onnx:$PATH`，未设 `CGO_ENABLED=0`）

| 门 | 读数 |
|---|---|
| `/d/work/base/gopath/bin/gofumpt.exe -l cmd/wisp/run.go cmd/wisp/panel_model_switch_187_test.go` | **空**，rc=0 |
| `go vet ./cmd/wisp` | **rc=0** |
| `go build ./cmd/wisp` | **rc=0** |
| `go test ./cmd/wisp -run 'Test187\|TestModelSwitch' -count=1 -timeout 30m` | **ok 0.117s**，5 发全 PASS（含上一腿两发 store 牙＋本腿三发） |
| `go test ./internal/panel -count=1` | **FAIL（rc=1）**，红名册恰 3 枚＝在册既有红：`TestApprovalCardViewJSONKeysMatchFrontendTypes` / `TestComposerContractTypesMatchFrontend` / `TestC21DesignTokensFourWayAgree`（前端 TS 声明＋C21，**非本腿面**；本腿 `internal/panel/**` 零改动） |
| `sh scripts/d22scan.sh`（CI 同形） | **rc=0 clean**：`bans #1-5 internal/=386, cmd/=63, #6 frontend/=116, #7 internal/tools/=48, #10 internal/risk/=21, #8 design/=39, frontend/=116, internal/=958, cmd/=239` |
| LF | 两件 `grep -c -P "\r"` ＝ **0 / 0** |

## 4. 边界

- 只 commit 不 push；`git add`/`git commit` 两处显式 pathspec 点名四件。
- ⛔ 未碰 `internal/panel/**`、`internal/ball/**`、`internal/tools/**`、`internal/config/**`（只读）、`frontend/**`、`design/**`、thresholds/golden；⛔ 未翻票面 AC 勾或 Status（只追加 Progress log 一行）。
- 起手即核红绿归属：开工时 tracked 改动为零（§0 #5），三枚 panel 红在改动前后同集，非在飞件造成。

## 5. 至此票 187 的 Go 侧三根线

1. 写腿 `SetRoleChatThinking`＋store case（187 腿＋A1621 欠账腿）；
2. 快照载件 `modelSwitch`（187 腿，reader-gated）；
3. **装配根注入（本腿）** —— 活机 `modelSwitch` 快照键自此而出。
