# 413 — `[ball] size` 滑块：Go 侧读+写接口契约（交外部 agent 接线用）

- 票：`.scratch/wisp/issues/413-ball-size-slider-has-config-and-validation-but-no-reader.md`
- 落地腿：本腿（2026-11-10）补齐 AC#2/AC#4/AC#6；AC#1 喂线由前腿 c8dbd2288 落（本腿复跑确认在盘）。
- ★本文只写契约，**不写任何前端码**（`frontend/**` 归外部 agent，本编队零写面）。

## 1. 今天已在盘的 Go 侧接口（真件、真测试、可直接引用）

### 读接口

```go
// internal/config/settings_413.go:75
func (m *Manager) BallSizeOptions() BallSizeView

// internal/config/settings_413.go:60
type BallSizeView struct {
    Size int // 当前生效直径 px；0 仅在无 Manager 支撑读取时出现（调用方答"未配置"，绝不造数）
    Min  int // 44（schema band，validate.go 判加载、SetBallSize 判写入用的是同一组常量）
    Max  int // 160（上限 72→160 经机主批准 A1636；默认 56 不变）
}
```

- 语义：读**活内存**（`Manager.Config()`），守卫写路径保证文件与内存相等——页面拿到的就是本进程下次会作用的数。
- nil Manager 安全：`Size 0`，band 照答。

### 写接口

```go
// internal/config/settings_413.go:45
func (m *Manager) SetBallSize(size int) ([]string, error)
```

- 入参：`size int`（**整数 px**，44..160）。
- 返回：`[]string`＝收据（文件实际改动的键路径，成功时 `["ball.size"]`）；`error`＝拒绝原因。
- 校验：越界**先拒后写**（落盘前），拒绝句与加载闸共用一个单源 helper（`internal/config/validate.go` 的 `ballSizeRangeError`，band 常量 `ballSizeMinPx`/`ballSizeMaxPx`），所以写路径拒绝句＝加载失败句，逐字同一串，永不分叉。

### 错误形状（错误码＝`observe.ClassConfig`，`*observe.Error`）

| 情形 | Detail（逐字） | 出处 |
|---|---|---|
| 越界（如 43） | `config.toml: ball.size 43 out of range [44, 160]` | `internal/config/validate.go`（ballSizeRangeError，单源） |
| 越界（如 161） | `config.toml: ball.size 161 out of range [44, 160]` | 同上 |
| 守卫被绕过时兜底 | `config: refused to write ball.size; the value would not survive validation…第 3 种拒因：校验不过…`（整验包装句） | `internal/config/sectionwrite.go` WriteSection 第 (3) 步 |

- 拒绝写**一个字节都不动**：config.toml 与内存逐字节保持原样（`TestSetBallSizeRefusesOutsideRange` 字节级钉住）。

### 生效档位（具名，⛔ 含糊）

- **内存档＝热**（`[ball]` 整节 hot-tier，D36；写完 `Manager.Config()` 立即可见）。
- **可见档＝重启**：球窗几何（`WindowEdgePx` 依赖它建窗，`internal/ball/ball_windows.go:64` SizePx、`:240` WindowEdgePx）在**建窗时**从开机种子读定（`cmd/wisp/ball_size_413_windows.go:40` ballSizeSeed413 → `cmd/wisp/resident_ball_windows.go:423` SizePx）——没有进程内 re-size，所以页面收据要答 **restart**，措辞照 `cmd/wisp/config_cmd_65.go:187` 那句（"生效档位：重启档…"），⛔ 谎称"立即生效"。

## 2. 前端接线要走的两条路由（★尚未接线，本腿写面不含 `internal/panel`——具名留给下一腿）

### 2.1 读：快照键 `ballSize`

- 形状照 `internal/panel/ball_style_65.go` 的 `BallStyleSection` 同形：`BallSizeSection{ Size int `json:"size"`; Min int `json:"min"`; Max int `json:"max"` }`。
- 接线点（四件，全部照 65 的既存链复刻）：
  1. `internal/panel/composer.go`：`BallSize *BallSizeSection \`json:"ballSize,omitempty"\``（照 :143 BallStyle 字段）；
  2. `internal/panel/pump.go`：`PumpSources.BallSize func() BallSizeSection` reader＋fill（照 BallStyle 的 reader-gated 形状，nil reader 不发键）；
  3. `cmd/wisp/run.go`：装配根注入 `BallSize: func() panel.BallSizeSection { v := mgr.BallSizeOptions(); … }`（照 run.go:1150 的 BallStyle 注入，同一个 mgr）；
  4. 载体件 `internal/panel/ball_size_413.go`（照 ball_style_65.go）。
- ⚠ 1、2、4 在 `internal/panel/**`——**超出本腿写面**（本腿写面＝`cmd/wisp/**`＋`internal/config/**`＋`internal/ball/**`），下腿补。

### 2.2 写：`config.set` 字段 `ball_size`

- 值形状：**整数的字符串**（`"90"`），照 `ball_style` 的 plain-value 形（不携带 provider/model）。
- 接线点（两件）：
  1. `internal/panel/config_write_65.go`（或新 roster 件）：`const FieldBallSize = "ball_size"`，进 `WritableFields`/`knownWritableField`/`acceptedWritableFields`/`takesValueFor`，且**不进** needsProvider/needsModel（照 :39-59 FieldBallStyle 的四枚钉子）；
  2. `cmd/wisp/panel_config_store.go` ApplySetting 加 case（在 `case panel.FieldBallStyle:` :284 旁）：

```go
case panel.FieldBallSize:
    // 票 413: 值＝整数 px 字符串；越界由 config 侧 SetBallSize 先拒后写
    // （拒绝句＝validate.go 单源句，见 §1 错误形状）。
    v, perr := strconv.Atoi(strings.TrimSpace(w.Value))
    if perr != nil {
        return panel.SettingWriteResult{Field: w.Field, Tier: panel.EffectiveNotApplied},
            fmt.Errorf("字段 %q 要的是 44..160 的整数（px），收到的值解析不了", w.Field)
    }
    written, err = s.mgr.SetBallSize(v)
```

- 档位：**静态答 `EffectiveRestart`**（size 没有 65 那种进程内 applier，读者行为固定，照 187 的静态标法；不像 ball_style 要实测）。
- ⚠ 件 1 在 `internal/panel/**`——同样超出本腿写面，下腿补；件 2 在 `cmd/wisp/**`，下腿顺手。

### 2.3 今天就能用的替代写口（无需等面板接线）

CLI 现只有 `ball.style`（`cmd/wisp/config_cmd_65.go:61` 键名册闭合）。若外部 agent 想先自证链路，Go 测试即用例：`Test413RoundTripSettingsWrite90ReachesTheBallWindow`（cmd/wisp）演示了写 90→读回 90→种子 90→球持 90 全程。CLI 扩键属下腿可选件，本腿未动（65 的键名册/用量文案有整组测试钉着，动它是另一片的活）。

## 3. 证据读数（本腿现跑，2026-11-10）

- 正向：`go test ./internal/config ./internal/ball -count=1` → ok（config 2.661s / ball 0.690s）；`go test ./cmd/wisp -run Test413 -count=1 -timeout 30m` → ok 3.198s，五枚全 PASS（含两枚 live boot，球窗真起、非 skip）。
- 往返：`Test413RoundTripSettingsWrite90ReachesTheBallWindow` PASS（1.15s）——写 90 ⇒ `BallSizeOptions().Size==90` ⇒ `ballSizeSeed413==90` ⇒ `ball.Options.SizePx==90`。
- 越界：`Test413RoundTripRefuses43And161WithValidateSentence` PASS（0.01s）——43/161 拒绝，错误串逐字含 `config.toml: ball.size 43 out of range [44, 160]`（161 同形），文件字节不变，种子仍读 90。
- 反向①（删喂线）：摘 `SizePx: hs.sizePx` ⇒ `Test413FeedReachesBallOptionsSizePx` **FAIL**（`ball.Options.SizePx = 56, want cfg.Ball.Size = 88`）；还原 ⇒ PASS（0.91s）。
- 反向②（删写校验）：摘 `SetBallSize` 越界守卫 ⇒ `TestSetBallSizeRefusesOutsideRange` **FAIL**（拒绝句退化成整验包装句，逐字断言红；注：`writeOneKey` 的 validate() 闸仍拒写，守卫不是唯一防线，但具名句契约破即红）；还原 ⇒ PASS。
- 门禁：gofumpt 空／`go vet ./internal/config ./cmd/wisp ./internal/ball` rc0／`go build` 同三包 rc0／d22scan clean／path-length VERDICT GREEN。

## 4. 本腿没做到的（诚实账）

- **面板两条路由未接线**（§2.1 的 composer/pump/carrier 三件＋§2.2 的 roster 件在 `internal/panel/**`，超本腿写面）——今天前端页面还拿不到 `ballSize` 快照键、`config.set ball_size` 还会被门拒。契约已具名（§2），下腿照 65 链复刻即可。
- **CLI 未扩 `ball.size` 键**（`config_cmd_65.go` 名册仍是闭合的 `ball.style` 单键）——可选项，留给下腿与面板接线同片做，避免本腿踩 65 的整组钉子测试。
- `frontend/src/components/config-screen.tsx:105` 那行的两处过期（范围 `44-72`、`blocked` 假前提）归外部 agent 改，本腿零触碰。
