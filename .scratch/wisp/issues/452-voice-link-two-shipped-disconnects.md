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
- （待派）
