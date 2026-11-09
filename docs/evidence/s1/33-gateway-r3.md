# 票 33 · panel-host-c27 网关 — r3 证据（AC#10 落地 + AC#7 判据实测推翻）

- 腿：33-r3（写腿；只 commit 不 push）
- 日期：2026-11-09
- 基线 HEAD：`31668c2d7`（`dev`）
- 出口树：`/d/tmp/wisp33w11`（`git clone --shared --no-hardlinks`，`core.autocrlf=false`、`core.eol=lf`）
- 本轮落地面：`cmd/wisp/panel_host_windows.go`、`cmd/wisp/panel_resident_windows.go`、`cmd/wisp/panel_config_ac10_windows_test.go`（新）、`internal/panel/ac7_tear_experiment_test.go`（新）

---

## 0. 一句话

AC#10 的五格里，**两格真正落地并各自有一枚能变红的判据**（`enabled`、`keep_alive_in_session`）；
第三格 `scale` **被 fork 结构堵死**，如实记原因、不假装。
AC#7 我按票面判据实测，**判据本身有一处不成立**：撕裂真实存在，但票面规定的"把 `Snapshot()` 改成持锁一次性快照"**消不掉它**——
这一发读数就是本轮的交付物，它把一张会变假绿的判据钉在了实测上。

---

## 1. 起手三核（先量，不采信转述）

| 核 | 做法 | 读数 |
|---|---|---|
| ① 派单点名的串/行还在不在 | `grep` 现跑票面、`sed` 取原文 | AC#7 原文在 `33-panel-host-c27.md:79-86`；AC#10 原文在 `:245`；13 个未勾格位置未变 |
| ② 引用出处 | 逐条回源到文件行 | `SnapshotPump.mu` 只护 `Publishes()`/`count()`（`internal/panel/pump.go` 仅两处 `p.mu.` 调用点）；`Marshal()` 生产调用者 **0**（重测）；`PanelSection` 五字段在 `internal/config/schema.go` |
| ③ 红绿是不是别人的在飞文件 | 出口树 + 主树双跑比对 | 见 §5：`internal/panel` 失败集合两树**逐行相同**，只有 `TestC21DesignTokensFourWayAgree`（`frontend/**` 面，非本腿） |

### 1.1 先修测量工具（否则后面全是假读数）

第一版出口树用默认 `core.autocrlf` 克隆，`internal/ball` 两发红：
`TestOrbServingConstants`（sha 不符）、`TestJarvisAssetShippedAsBuilt`（报 CRLF）。
我没有立刻当缺陷上报，而是量了 **blob**：

```
git cat-file -p HEAD:internal/ball/assets/orb-siri.html | tr -dc "\r" | wc -c   → 0
git cat-file -p HEAD:internal/ball/assets/orb-siri.html | sha256sum            → d93da649d11bc4bc…（= 钉值）
```

⇒ **blob 是干净的，红是检出物转 CRLF 造成的**（`.gitattributes` 里 `* text=auto` 在 Windows 检出会转）。
重建出口树时加 `-c core.autocrlf=false -c core.eol=lf`，两发转绿。

> 同一条也给下一腿留个坑标记：**在 Windows 上做出口树门禁，必须显式关掉 autocrlf**，否则 `internal/ball` 会红两发假钉。

---

## 2. AC#10：`[panel]` 五字段的生产者

票面判据（逐字要点）：五字段都要有生产者；`width`/`height`/`scale` 要**真设窗口边界**；
⚠ 判据要能区分"读值并真设边界"与"只在测试里读一下"；⛔ 不许靠改 `default:"640"` 交差。

### 2.1 先量清哪几格本来就有生产者

| 字段 | 本轮之前的生产者 | 判定 |
|---|---|---|
| `width` | `panelGeometrySource` 每次建窗重读 `config.toml`，`windowOptions()` 交给唯一建窗点（票 255 AC#4 已落，`cmd/wisp/panel_geometry_255_test.go` 是那票的判据） | **已有**，本轮不重做 |
| `height` | 同上 | **已有**，本轮不重做 |
| `enabled` | **零**（`grep` 生产代码 0 处） | 本轮落地 |
| `keep_alive_in_session` | **零** | 本轮落地 |
| `scale` | **零**，且**被 fork 堵死** | 本轮不落地，记原因（§2.4） |

### 2.2 `enabled` 的生产者

**落点**：`cmd/wisp/panel_host_windows.go`

- `PanelManager` 新增 `enabled func() bool`、`disabledByConfig bool`、`disabledReason string`
- `withPanelEnabled(src func() bool)` 选项；`nil` 源 = 启用（schema 默认 `true` ⇒ 所有既有构造点行为不变）
- `panelDisabled()` 在 **`bringUp` 最顶部、`probePanelRuntime()` 之前**做闸
- `errPanelDisabledByConfig` 是**具名**错误，与 `errPanelRuntimeMissing` **分开**
- `Disabled() (bool, string)` 是对外读数

**为什么闸必须在探测之前**（这是判据的核）：策略拒绝与机器拒绝是两件不同的事。
票 37 的原生降级卡要读 `Unavailable()` 决定"能不能用原生面顶上"；
一个只是把面板关掉的用户，**不该**让那张卡弹出来说"这台机器没有 WebView2"。

**装配根接线**：`cmd/wisp/panel_resident_windows.go:322` 一带

```go
withPanelEnabled(panelEnabledSource(dataDir)),
withPanelKeepAlive(panelKeepAliveSource(dataDir)),
```

`panelEnabledSource` 与既有 `panelGeometrySource` 同形（闭包 + 每次现读 `config.LoadFile`），
失败方向刻意是 **true**：文件读不出来不能静默关掉一个默认开启的功能。

### 2.3 `keep_alive_in_session` 的生产者

`Hide()` 末尾新增 `teardownAfterHideUnlessKept(keepAlive)`：源答 `false` ⇒ `Destroy()` 真拆窗；
默认/`true` ⇒ 保持 C27 的"隐藏不销毁"逐字节不变。
拆窗走 `Destroy()` 既有那条库自带的 WM_CLOSE 链，**没有造第二条拆除路径**。

### 2.4 `scale`：本轮不落地，原因是结构性的（不是排期）

三条独立读数，任一条都足以堵死：

1. `forks/go-webview2/pkg/edge/ICoreWebView2Controller2.go:16` —— `PutZoomFactor` **只以原始 `ComProc` vtbl 槽存在**，没有 Go 包装；
2. 同文件 `:38-39` —— `ICoreWebView2Controller2` 只暴露 `AddRef`/`GetDefaultBackgroundColor`/`PutDefaultBackgroundColor`；
3. `forks/go-webview2/pkg/edge/chromium.go` —— `Chromium.controller` 是**未导出**字段，`cmd/wisp` 够不到。

⇒ 要落 `scale` 必须改 `forks/**`，而 `forks/**` 本腿禁写。
**如实记：这一格还欠着，欠的原因是"要动 fork"，不是"没时间"。**
并且钉了一枚会自己失效的判据 `TestAC10ScaleHasNoProducerAndTheReasonIsNamed`：
将来谁补了 fork 的 Go 包装，这发会**主动红**并告诉他"借口失效了，去接线并删掉这发"。

### 2.5 AC#10 的判据（每发都写了红形）

文件：`cmd/wisp/panel_config_ac10_windows_test.go`

| 用例 | 红形（删什么会红） |
|---|---|
| `TestAC10DisabledPanelRefusesBeforeTheRuntimeProbe` | 删掉闸 / 把闸挪到探测之后 ⇒ 探针 spy 计数 ≠ 0 |
| `TestAC10DisabledRefusalIsAPolicyAnswerNotAMachineOne` | 把策略拒绝写进 `runtimeMissing` ⇒ 票 37 的卡会对关掉面板的用户误报机器缺陷 |
| `TestAC10EnabledHostStillProbesTheRuntime` | 闸变成墙（`nil` 源或 `true` 源也不探测） |
| `TestAC10EnabledSourceReadsTheRealConfigFile` | 源改成装配时快照 ⇒ 改写 `config.toml` 后仍答旧值（与 `[panel]` 登记为 hot 矛盾） |
| `TestAC10KeepAliveSourceReadsTheRealConfigFile` | 同上 |
| `TestAC10KeepAliveFalseTearsTheWindowDownOnHide` | `Hide` 无视源（该拆不拆 / 不该拆却拆） |
| `TestAC10AssemblyRootWiresBothSources` | 装配根没把源交给 `NewPanelManager` ⇒ 源存在但无人问（"文案在、控件不在"的接线版） |
| `TestAC10ScaleHasNoProducerAndTheReasonIsNamed` | fork 补了包装 ⇒ 主动红，逼后来者接线 |

### 2.6 突变实测（真的变红了）

删掉 `bringUp` 里那三行闸：

```
--- FAIL: TestAC10DisabledPanelRefusesBeforeTheRuntimeProbe (5.20s)
    panel_config_ac10_windows_test.go:89: a host told [panel] enabled = false brought a window up anyway: AC#10's enabled producer is not wired to the gate
```

还原后 `ok github.com/CarlosShao/wisp/cmd/wisp 0.179s`。

---

## 3. AC#7：判据实测——撕裂真在，但票面规定的解药不管用

### 3.1 票面判据（逐字）

> 判据：本票落地 `Marshal()` 的第一枚**生产**调用者之后，造一发"两个触发点之间状态变了"的用例，
> 它必须能红（把 `Snapshot()` 改成持锁一次性快照 ⇒ 该发由红转绿）。

### 3.2 三核结果

- `Marshal()` 生产调用者**重测仍为 0**（`grep -rn "\.Marshal()" cmd/ internal/ | grep -v _test.go` 只命中定义行）。
  ⇒ 判据的**前置条件今天仍不成立**。
- `Publish()` 自己内联 `json.Marshal(snap)`（`internal/panel/pump.go:529`），**没有**调 `Marshal()`。
  ⇒ 同一个包里有两条编码路径，这本身是"第二个真相"。
- `SnapshotPump.mu` 的全部调用点只有两处（`pump.go:553-561`），护的是 `Publishes()`/`count()`。
  四个数据源各持自己的锁（`Gate.mu` 20 处、`StreamLog.mu` 18 处）。

### 3.3 实验：`internal/panel/ac7_tear_experiment_test.go`

三发，全部实跑（`ok`）：

| 用例 | 读数 |
|---|---|
| `TestAC7TornPacketIsObservable` | **红成立**：一个包里同时带着 T1 时刻的待审卡和 T2 时刻的完成结果 |
| `TestAC7PumpLockDoesNotCloseTheTear` | **★关键读数**：把泵锁**全程持住**（票面规定的解药的最强形式），包**依然撕裂** |
| `TestAC7PumpLockGuardsOnlyTheCounters` | 把上面这条写成类型性质，防它被当成夹具巧合 |

实跑输出：

```
ac7_tear_experiment_test.go:140: AC#7 red confirmed: one packet carried a card from instant T1 beside a result from instant T2
ac7_tear_experiment_test.go:160: AC#7 finding: the pump lock is held across the entire build and the packet is STILL torn,
    because SnapshotPump.mu guards only Publishes()/count() while every source owns its own lock
```

### 3.4 这条读数意味着什么

票面说"持锁 ⇒ 由红转绿"。实测是**持锁也不绿**。
所以任何按票面写出来的"红转绿"用例，**只能靠测别的东西来变绿**——
这正是票面自己在 ⚠ 里禁止的"恒真判据"，只是升了一层：
不是"永远绿的钉"，而是"**红转绿的因果是假的**"。

**真正的解**（本腿只命名、不假装做）：跨源一致性需要一个**共享纪元**——
要么单一写入者为每个源打戳，要么在包上带序号让读方丢弃不自洽的一对。
后者是主流做法，归票 35 的 D38d 传输线。
本腿把测量留在盘上，让那个决定**基于读数而不是基于票面的假设**。

### 3.5 AC#7 本腿的诚实结论

- 判据的前置条件（`Marshal()` 生产调用者）**未满足**；
- 撕裂**真实存在**且可确定性复现（不需要 race）；
- 票面规定的解药**经实测无效**；
- ⇒ 本腿**不勾这一格**，把上面三条与实验文件留给裁决者。
  我没有为了让格子好看而去接一个装饰性的 `Marshal()` 调用方——那是"造绿"，比不交更坏。

---

## 4. 门禁读数（出口树 `/d/tmp/wisp33w11`）

| 门 | 命令 | 读数 |
|---|---|---|
| gofumpt | `/d/work/base/gopath/bin/gofumpt.exe -l cmd/wisp/ internal/panel/ internal/ball/ internal/config/` | **空**（无输出） |
| vet | `go vet ./cmd/wisp ./internal/panel ./internal/ball ./internal/config` | **rc=0** |
| build | `go build ./cmd/wisp` | **rc=0** |
| `internal/panel` | `go test ./internal/panel -count=1` | 1 发红：`TestC21DesignTokensFourWayAgree`（**非本腿**，见 §5） |
| `internal/ball` | `go test ./internal/ball -count=1` | **ok** |
| `internal/config` | `go test ./internal/config -count=1` | **ok** |
| d22scan | `bash scripts/d22scan.sh` | **clean**（内部 388/64、frontend 121、tools 48、risk 21、design 14、emoji 971/245） |
| LF | 三件改动文件 `tr -dc "\r" \| wc -c` | 全 **0** |

---

## 5. 红绿归属（三核之③）

`internal/panel` 的失败集合，**主树与出口树逐行 diff 相同**：

```
--- FAIL: TestC21DesignTokensFourWayAgree
=== diff ===
IDENTICAL (pre-existing, not mine)
```

根因在 `frontend/**`（`frontend/scripts/gen-tokens.mjs`、`frontend/src/styles/tokens.generated.css`），
那是**机主外部 agent 的写面**，本编队零写面。
⇒ 不是本腿造成，也不由本腿修。

---

## 6. 边界与还没做的（如实）

1. **AC#7 未勾**：前置条件不成立 + 判据规定的解药实测无效（§3.5）。
2. **`scale` 未落**：要改 `forks/**`，本腿禁写（§2.4）。钉了一枚会主动失效的判据。
3. **`keep_alive_in_session=false` 的屏幕半场没有实机读数**：本腿只量了"决策 + 真拆窗"（`Destroy` 被走到），
   没有量"窗口真的从屏幕上消失"。那需要 winlive 读数，本腿没有，也不假装有。
4. **AC 格子一个都没勾**：按规矩，勾格属裁决者（非实现者）动作；本腿只追加 Progress log。
5. `cmd/wisp` 整包 `rc=1` 是**已知的真窗家族**（见票面既有记录），判据是按**失败用例名**比对，不是按 rc。

---

## 7. 可复跑（给裁决者）

```bash
# 出口树（Windows 上必须关 autocrlf，否则 internal/ball 红两发假钉）
cd "/d/work/workspace/projects plans/Wisp"
git -c core.autocrlf=false -c core.eol=lf clone --shared --no-hardlinks . /d/tmp/wisp33w11
cd /d/tmp/wisp33w11 && git checkout dev
mkdir -p build third_party/sherpa-onnx
cp <main>/build/*.dll build/ && cp <main>/third_party/sherpa-onnx/*.dll third_party/sherpa-onnx/
unset GOTMPDIR TMPDIR TEMP TMP
export PATH=$PWD/build:$PWD/third_party/sherpa-onnx:$PATH

# AC#10 判据
go test ./cmd/wisp -run 'TestAC10' -count=1 -v

# AC#7 实验
go test ./internal/panel -run 'TestAC7' -count=1 -v

# 突变（把闸删掉，应当红）
# 删 cmd/wisp/panel_host_windows.go 里 bringUp 顶部那三行 if err := m.panelDisabled() ...
go test ./cmd/wisp -run 'TestAC10DisabledPanelRefusesBeforeTheRuntimeProbe' -count=1
```
