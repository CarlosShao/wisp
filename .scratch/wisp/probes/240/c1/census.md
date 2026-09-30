# 240-c1 只读普查：麦克风采集侧 / KWS 侧 / 接线落点 / 半成品名册

- 腿：`240-c1`（只读，不产码、不跑编译与测试）
- 取钟：`2026-09-30 11:26:55 +0800`（`date` 与 `mkdir` 同发，落笔在本节之前）
- 锚点：起手 `git log --oneline -1` = `0589fd9c`（球被窗口切成方形那枚）
- 红线自查：本腿零 `go build`／零 `go test`／零 `gofumpt`／零 `go vet`／零 `d22scan`；
  零读、零写、零转述 `frontend/**` 与 `design/**`；不改 `PLAN.md`／`docs/specs/**`／
  `docs/SLO.md`／`internal/ball/thresholds.go`／golden／`allowlist.txt`／台账／HANDOVER。

---

## §0 现量口径（我这轮跑了哪些尺）

> 每条尺都先在本机真跑，读数照抄。尺失效（rc=1 且期望非空）本身就是读数。

| # | 尺（命令要点） | 本机真实读数 | 口径说明 |
|---|---|---|---|
| R1 | `grep -rn 'CarlosShao/wisp/internal/audio' --include=*.go .`（排 `.scratch`） | **零命中，rc=1** | 全仓模块内没有任何文件 import `internal/audio` |
| R2 | 同上，`internal/speech` | **零命中，rc=1** | `internal/speech` 只有 `doc.go`，无人 import |
| R3 | `grep -rn 'NewWASAPIMicrophone\|WavInjector' --include=*.go .`（排 `.scratch`） | 命中全部落在 `internal/audio/` 自身：定义 + `hotplug_test.go:529` 等测试；`cmd/`、其它 `internal/` 包 **零命中** | 真设备栈有实现、无生产调用者 |
| R4 | `grep -rn 'SetAudioLevel' --include=*.go .`（排 `.scratch`） | 非测试生产者仅 `cmd/balldebug/main.go:418`、`:477`；定义 `internal/ball/liquid_windows.go:42`；测试 `internal/ball/live_windows_test.go:626/637/645/674` | 与编排者 09-30 当面认的缺口一致 |
| R5 | `grep -rn 'sherpa\.' --include=*.go .`（排 `.scratch`） | `cmd/wisp/` 只用 `GetVersion`/`GetOnnxruntimeVersion`（`doctor.go:60/81`、`main.go:165`）；`NewKeywordSpotter` **只出现在 `scripts/spike/`** | 生产码链了 sherpa 但只用来报版本号 |
| R6 | `grep -rniE 'keywordspotter\|kws' --include=*.go .`（排 `.scratch`） | 命中三类：`scripts/spike/{model-residency,speech-baseline}/main.go`（真建 spotter）、`internal/agent/approval/*`（否决通道占位）、`internal/statemachine/*`（事件与守卫已成型） | 待 §2 逐枚拆 |
| R7 | `cat scripts/spike/go.mod` | `module github.com/CarlosShao/wisp/scripts/spike` —— **独立模块**，不在主模块构建图里 | 所以 R5 的"生产零命中"与 spike 里有真代码并存，不矛盾 |
| R8 | `cat go.mod`（主模块） | 直依赖含 `github.com/k2-fsa/sherpa-onnx-go v1.13.8`；`sherpa-onnx-go-windows` 等三平台包是 indirect | cgo 那层＝上游模块，不是本仓 vendored 源码 |
| R9 | `ls third_party/sherpa-onnx/` | 三枚 DLL：`onnxruntime.dll`、`sherpa-onnx-c-api.dll`、`sherpa-onnx-cxx-api.dll` + `.cache-manifest.json`；**无 .go、无 .h** | 本仓不持有 sherpa 源码，只持有运行期 DLL |
| R10 | `sed -n '515,535p' docs/PLAN.md`（只读） | 空闲 ≤25MB/≤0.5%、待唤醒 ≤90MB/≤2% 两档原文确认在位 | 用于 §2 找代码里的对应常量 |

（§0 未完：后续节追加尺子编号从 R11 起。）

---

## §1 采集侧：`internal/audio` 与 `internal/speech` 各暴露什么

（待写。骨架要点已就位：R1/R2 两条零命中是本节的起点。）

---

## §2 KWS 侧：绑定怎么引的、谁建过 spotter、两档预算的落点、opt-in 与指示

（待写。）

---

## §3 两条链的落点清单

### 3.1 采集 → 电平 → 球

（待写，一张表。）

### 3.2 采集 → KWS → 唤起会话

（待写，一张表。）

---

## §4 半成品名册

（待写：每枚 `文件:行` + 「补一行就通」/「缺一整块」。）

---

## §5 我答不上来的

（待写，本节不许为空。）
