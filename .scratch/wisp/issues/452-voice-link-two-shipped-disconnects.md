# 票 452 — 语音链路两个**出厂构建**里的断点：barge-in 根本不在 `go build` 里；拔麦不推状态

**状态**：待派（两处都经验收腿现读现跑发现，编排者已复核）

## 断点一：barge-in 不在出厂构建里
`scripts/build.ps1:270` 的 `go build` **不带 `-tags webrtcap`** ⇒ 进的是 `aec_stub`：`NewApm` 恒返回 `ErrUnavailable`、monitor 不武装。
⇒ 票面/证据里那些 **≤400ms barge-in 读数全部来自专用构型**，出厂用户拿到的那份**没有这个能力**。
⇒ 而且 26p-v2 复跑 6 窗只 **1 火** ⇒ 「每轮必火」也不成立。

## 断点二：拔麦克风不推状态，球可能停在 Listening
`EvAudioDeviceLost` 在产码里的全部命中 = **4 处**，其中 **2 处是注释**：
- `internal/statemachine/events.go:30` 常量声明
- `internal/statemachine/table.go:97` D43 转移表条目（`From: StateListening → To: StateError`）
- `cmd/wisp/resident_audio_windows.go:39` **注释**
- `cmd/wisp/resident_audio_windows.go:358` **注释**，逐字「这一条不推状态、也不拒绝启动」

⇒ **零派发**：语音侧从不读 `mic.Err()`，因此**表里有那条边、码里没人走那条边**。用户拔掉麦克风，球可以一直停在 Listening。这正是票 32 AC#5 说的 stuck state。

## AC

- [ ] **AC#1 先把两个断点都变成可判事实**：(a) 跑 `go build -n`／读 `build.ps1` 现证**没有** `webrtcap` 标签，并给"打了标签时 `NewApm` 返回什么、没打时返回什么"两条实测；(b) 写一枚测试：把 `mic.Err()` 置成设备丢失，断言状态机**收到** `EvAudioDeviceLost` 并走到 `StateError`。⛔ 先改产码 —— **它必须先红**。
- [ ] **AC#2 决定 barge-in 的出厂形态并落成事实**（⓪⓪⓪：按最有价值＋主流 harness 选，写理由）：三选一并写明代价 —— ①**出厂即带 `webrtcap`**（代价：体积/构建依赖）②**出厂不带，但产品明说"此构建无声学回声消除"**并让界面可判 ③**换成不依赖该标签的默认路径**。⛔ 不选"什么都不做"；⛔ 不许一边不带标签一边继续引用 ≤400ms 那些读数当成绩。
- [ ] **AC#3 把 `EvAudioDeviceLost` 接上派发**：语音侧检测到设备丢失就派发该事件（走 `observe.Registry.Spawn` 若要新协程），球与面板都离开 Listening；拔麦测试由红转绿。
- [ ] **AC#4 两处注释同步**：`resident_audio_windows.go:358` 那句「这一条不推状态」必须改成事实句（接上之后就变成假的）。⇒ 这是**同一条改动里的必做项**，⛔ 留一句过期注释。
- [ ] **AC#5 门禁**：`gofumpt` 空／`go vet` rc=0／`d22scan` clean／path-length GREEN／`go test ./internal/audio ./internal/statemachine ./cmd/wisp` 全绿（整包跑不完如实写）。

## 边界
⛔ 动 `thresholds.go`／golden／`allowlist.txt`／`docs/PLAN.md`／`docs/specs/**`／D43 转移表（**表里那条边已经对，别改表**）／`frontend/**`／`design/**`；⛔ 放宽既有断言；⛔ 把 `ErrUnavailable` 改成"假装成功"；⛔ add -A/--amend/reset/rebase/stash/checkout .；⛔push。
写面：`internal/audio/**`、`internal/speech/**`、`internal/statemachine/**`（只读该表）、`cmd/wisp/**`（每波一枚写腿）、`scripts/build.ps1`、`docs/evidence/s1/452-*`。

## 撤销口令
「452 撤」＝删本票与其证据件（不 reset/amend）。

## Progress
- 2026-11-15 写腿（agent-51，`cmd/wisp` 本波唯一写腿）：
  - AC#1(a) 现跑：`scripts/build.ps1:270` 的 `go build -trimpath -ldflags …` 无 `-tags`，全脚本 `webrtcap` 命中 0；
    `go list ./internal/audio/aec` = `GoFiles=[aec.go aec_stub.go]` ⇒ `NewApm` 恒 `nil, ErrUnavailable`（已由 `internal/audio/aec/aec_shape_452_test.go` 钉住）。
    带标签：`go build -tags webrtcap ./internal/audio/aec` **编不过**（`api/scoped_refptr.h: No such file or directory`，本机没有 webrtc 树/meson）；
    `CGO_ENABLED=0 go build -tags webrtcap` rc=0（标签单独不给能力）。读数原文 `.scratch/wisp/probes/452/ac1-a-measure.txt`。
  - AC#1(b) 先红：两枚测试写在 `cmd/wisp/resident_audio_device_lost_452_windows_test.go`，改产码**之前**跑出
    「姿态是空串 `""`」与「拔麦 3s 后仍 Listening、球投影 `[]`」两条红句（逐字抄进证据 §3）。
  - AC#2 裁定 **②**（出厂不带标签 + 产品明说 + 界面判得了）：`aec.BindingBuilt()` 能力位、`bargeinShippedPosture452` 常量在
    **构造期**就写进 `reason`、启动印一次、`posture()` 增 `插话=`。①的代价＝整条原生构建流水线（本机无树、无 meson/ninja、`fetch-deps.ps1` 无此路），
    ③的代价＝重做 389/392/401 三轮读数；面板能力字段出写面 ⇒ 已写明落点移交（证据 §2.3）。
  - AC#3 已接：`residentAudio.captureErr()`（粘性，源头 `Err()` 是破坏性一次读）+ `rv.checkAudioDevice()` +
    `rv.run()` 两个 select 各一条 ticker case（**不新增协程**，骑既有 owner/`asr-infer` 线程）+ 仅从 `StateListening` 派发 D43#14（表未动）。
    接线 `rv.deviceErr = ra.captureErr` 在 `rv.start()` 之前；测试对该生产行做存在性断言。
  - AC#4 已改：`resident_audio_windows.go` 启动分支的「这一条不推状态…」与文件头 "It pushes NO state" 改成「启动期/运行期」两句。
  - 突变：M1 关掉派发 ⇒ 行为测试红；M2 删掉生产接线行 ⇒ 接线断言红。两次都还原。
  - 证据：`docs/evidence/s1/452-voice-link-two-shipped-disconnects.md`。
  - 未做：重新插回/换设备后的恢复（球停 Error 等 `EvErrorAck`，D43#38）；面板能力字段（出写面，已列落点）。
