# 票 65 · 步：CLI 切球风格路径（`wisp config get/set/list`）· 收件证据

- 写腿：65-ball-style-cli（cmd/wisp 独占波；`internal/ball/**` wave73 在飞未碰；`frontend/**`、`design/**` 零写面）
- 日期：2026-11-07 20:2x +08
- commits＝本 commit；只 commit 不 push；AC 勾与 Status 未动
- 机主原话锚点：「两个都要，但是设置项里面外观可以配置切换风格，**默认就是玻璃球**，切换 Jarvis 才换球的风格」

## 1. 步1 取证（全部现跑，非转述）

1. **子命令分发形状**（`cmd/wisp/main.go`）：`func main` 顶部 `switch args[0]`，每腿一
   `case` → `attachParentConsole()` + `os.Exit(cmdXxx(args[1:], xxxIO{...}))`。非冻结契
   约：133 门禁（`leg_dispatch_gate_133_test.go`）只要求三件事——分发表枚举不出错、
   `minLegs133=12` 下限（现 13）、usage 块与分发表双向对账（`censusVsUsage133`）＋每腿
   「有钉或有旁证裁定」（`coverageReds133`）。⇒ 加子命令合法，代价＝usage 加一行＋
   `legCovers133` 注册一枚真调入口的测试。已照办。
2. **`SetBallStyle` 签名与守卫**（`internal/config/settings_65.go:46`）：
   `func (m *Manager) SetBallStyle(style string) ([]string, error)`——枚举守卫在
   setter 内部（`isBallStyle` 问 schema 的 `ballStyles`，拒绝句点名字段＋词表），守卫
   在写之前，拒绝时文件字节不动。可直接从 CLI 调，CLI ⛔ 自造第二套校验。
3. **配置写盘既有路径**：`writeOneKey`（settings.go:269）→ `WriteSection`
   （sectionwrite.go）——预读快照、apply、`validate(candidate)` 全量校验、D36 方向裁
   定、守卫合并写（票 226 merge writer）。task 提示的 `panelwrite.go:168` 不存在；真
   形状即上面这条链，无需新写盘路径。
4. **主流 CLI 配置命令样子**：`git config user.name` / `kubectl config set` ⇒
   `wisp config get <key>` / `wisp config set <key> <value>` / `wisp config list`，
   退出码沿用本仓约定（2 用法/环境、1 操作性拒绝、0 成功，同 `wisp slo`/`wisp secret`
   的注）。
5. **生效档位实测（现跑核清，收据据此答）**：`[ball]` 整段热档（schema.go:177 注），
   reload 看门狗 1s 收文件改动进内存（config_reload.go），但**换球只挂两个钩**：面板
   写路径的 `applyBallStyle65`（进程内，panel_config_store.go:300）＋开机种子
   （resident_windows.go:419/432）。CLI 是另一个进程 ⇒ CLI 写盘**答不了 now**，诚实
   档位＝重启档。⇒ 另证 provenance：ball.style 非锁定段、方向 neutral，neutral 分支
   不落任何 SECTION-WRITE audit 行（sectionwrite.go 只在 loosen/tighten 落）⇒ 直用
   `SetBallStyle`（内部标 SourcePanel/Origin panel.settings）**不会有任何假账落盘或
   落日志**，CLI 侧零新增 internal/config 面。

## 2. 裁定与落地形状

- **形状**（按铁律 ⓪⓪⓪，主流 harness 样子，自拍并说明理由）：
  - `wisp config get ball.style` → 打印生效值（文件缺 ⇒ schema 默认值，stderr 明说
    这是默认值不是文件里的值）；rc 0。
  - `wisp config set ball.style jarvis` → 文件缺则先走 `ensureFirstRunConfig`（票 198
    的唯一首建口——显式 `set` 就是那个「asked it for」）→ `config.NewManager` →
    `SetBallStyle` → 收据两行：文件路径＋键值、档位句（重启档，见 §1.5）。rc 0。
  - `wisp config set ball.style glass` → 切回，同上。
  - `wisp config list` → 名册一行：`ball.style = <当前值>    取值: glass|jarvis`，词
    表走 `BallStyleOptions`（schema 自己的名单，零拷贝）。
- **键名册封闭**：CLI 只开放 `ball.style`（`configCliKey65`），未知键 get/set 一律
  rc 2 并点名名册——与面板 `knownWritableField` 四表并集同形（封闭名册，不做猜键）。
- **守卫照用**：非法值拒绝句＝`SetBallStyle` 自己的具名句（「不在枚举内（glass|
  jarvis）」），CLI 原样透传，⛔ 复写。纵深防御实测：`WriteSection` 第(3)步
  `validate(candidate)`（validate.go:70 同一张 ballStyles 表）也会拒——两层守卫，文
  件永远安全；牙④因此钉住**哪层在说话**（见 §3.2）。
- **默认值走 schema**：CLI 全文件零 `glass` 字面量。缺文件 ⇒ `config.NewDefaults()`
  （defaults.go:58「exactly the `default:"..."` tag values」＝schema 默认表的导出
  口）；文件在 ⇒ loader 已叠加默认（parse.go decodeStrict onto cfg carrying
  defaults）⇒ 打印的就是生效值。
- **不做的**（理由）：不加 `--set-ball-style` 旗标（子命令不是冻结契约，主流就是子
  命令）；不给 internal/config 加 Native-source 新 setter（neutral 方向零 audit 行，
  加面无收益）；CLI 不试图探活常驻进程（无 IPC，探了也只能说「下次启动生效」）。

## 3. 五牙读数（全绿，rc 见 §4）

- 牙① `TestConfigLegGetSetGuarded65` 前半：空数据根 `get ball.style` → `glass`
  rc 0，stderr 明说 schema 默认值。**PASS**
- 牙② 同测试中段：`set ball.style jarvis` rc 0，收据含 `ball.style = jarvis`＋重启
  档位句，`config.toml` 真在数据根；`get` 读回 `jarvis`。**PASS**
- 牙③ 同测试后段：`set ball.style glass` rc 0；`get` 读回 `glass`。**PASS**
- 牙④ `TestConfigSetBallStyleInvalidValueRefused65`：种子 glass 后
  `set ball.style 紫色` → rc 1，stderr＝前置守卫具名句（含 `不在枚举内`＋
  `glass|jarvis`），**config.toml 前后逐字节相等**（os.ReadFile + bytes.Equal），
  `get` 仍答 glass。**PASS**
- 牙⑤ **突变演练**：把 `settings_65.go` 的枚举守卫改成
  `if false && !isBallStyle(v) {`（守卫旁路）后重跑牙④——**如设计转红**，红行逐字：
  ```
  config_cmd_65_test.go:143: refusal must be the setter's own named sentence (field, enum wording, schema vocabulary); missing "不在枚举内" in stderr: wisp config: 拒绝写入：config: config: refused to write ball.style; the value would not survive validation, so config.toml is unchanged。第 3 种拒因：校验不过：值本身过不了这份 schema 的校验，这一行要改的是值；文件一个字节都没动。: config: config.toml: ball.style "紫色" must be one of glass|jarvis
  FAIL
  FAIL	github.com/CarlosShao/wisp/cmd/wisp	0.131s
  ```
  红的语义如实：validate 第二层仍拒写（红行句内自带「文件一个字节都没动」），但
  CLI 不再回前置守卫的具名句 ⇒ 断言「哪层在说话」的红点正确。还原后
  `sha256sum internal/config/settings_65.go` =
  `3d87f82430c8f30bcbdce1d8ef0f7ded1dc856bf229d03c85d282774d705b06f`（与突变前一致，
  逐字节还原），复绿。附牙：`TestConfigKeyRosterRefused65`（未知键/未知子命令/缺值
  全 rc 2 且不建文件）＋`TestConfigListShowsRoster65`（list 双态）**PASS**。

## 4. 门禁读数（2026-11-07 20:2x +08，环境照纪律 unset GOTMPDIR/TMPDIR/TEMP/TMP +
PATH 前置 build、third_party/sherpa-onnx）

| 门禁 | 读数 |
|---|---|
| `gofumpt.exe -l`（四件点名：config_cmd_65.go、config_cmd_65_test.go、main.go、leg_dispatch_gate_133_test.go） | 空 |
| `go vet ./cmd/wisp ./internal/config` | rc 0 |
| `go build ./cmd/wisp` | rc 0 |
| `go test ./cmd/wisp -run 'Test65\|TestConfig' -count=1 -timeout 30m` | ok 0.394s（含新四牙＋既有 Test65 四牙） |
| `go test ./cmd/wisp -run 'TestAC1AC2DispatchHopGate133\|TestAC4EveryLegIsNailedOrRuled\|TestAC3SelfProof133' -count=1` | PASS（133 对账含新 `config` 腿：usage 行＋legCovers133 钉 TestConfigLegGetSetGuarded65→cmdConfig 双向核过） |
| `go test ./internal/config -count=1` | ok 4.457s |
| `tools/d22scan/d22scan.exe .` | clean rc 0（internal 956＋cmd 236） |
| LF | 四件 Go 件 `file` 读数无 CRLF（UTF-8 text） |

## 5. 边界如实（⛔ 不许报成「完事」的部分）

1. **设置页里那行还是不存在**：`frontend/src/components/config-screen.tsx` 的
   APPEARANCE_ROWS 与 `frontend/src/lib/panel.ts` 的 `ballStyle` 键归外部 agent
   （Q-51），本编队零写面。机主今天的切换路径＝**本 CLI**，设置页行等 Q-51。
2. **CLI 档位＝重启档**：`wisp config set` 写盘后，正在跑的 wisp 不会立刻换球（换球
   挂面板写路径与开机种子两钩）；收据句已如实说，要立刻看到就重启 wisp。无 IPC 面
   前这是诚实上限，不是偷懒。
3. `cmd/wisp` 盘上有两枚别人的未跟踪测试件
   （panel_resident_323_windows_test.go、panel_resident_assembly_323_test.go）：本腿
   未碰、未收入 commit（pathspec 点名制）。
4. `internal/ball/**`（wave73 在飞）、`internal/panel/**` 既有件、
   `internal/ball/assets/orb-siri.html`：零触碰（`git status` 读数旁证：本 commit 后
   这些路径无本腿改动）。
