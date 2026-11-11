# 票 297 落地腿 `297-l2` 读数档 — AC#2 形①机内收窄 ＋ AC#5 执行面现跑 ＋ AC#3 默认档核验

- 腿：`297-l2`（派单：读全票 297 剩余 AC，机内能做的做完；写面＝internal/audio/** 仅新增测试/probes；⛔ 本波别碰 cmd/wisp——423 收口占槽）。
- 锚：起手 HEAD＝`ddb759340`（基线两发时点；过夜后树内 423 遗留件被槽主改过，见 §3）。
- 写面：`internal/audio/aecspike/formone_297_windows_test.go`（新增 1 枚 windows-tagged 仪器用例）＋ `probes/297/l2/**`。**⛔ 产码 0 字节改动**（`git show --stat` 可核；AC#0..AC#4 框全未勾）。

## 1. 加了什么

`TestAC297FormOneRenderReachProbe`（aecspike 包内，`WISP_AEC_SPIKE=1` 门，同 388 rig 家族形状）：

- 相 A：生产麦克风（真 `WASAPIMicrophone`＋LevelSink）6s 安静窗；
- 相 B：5×1s 已知幅度件（`SynthSpeechLike`，gain=0.35）经 `RunLoopback` 渲染默认渲染端点；每次渲染自带端点环回 tap（**数字真值**：端点混音里到底有没有声）；
- 相 C：2s 恢复窗；
- `[volume]` 只读快照：渲染会话 `ISimpleAudioVolume::GetMasterVolume/GetMute`（slot 4/6）＋默认渲染端点 `IAudioEndpointVolume::GetMasterVolumeLevelScalar/GetMute`（slot 9/15）。**全出参指针读、零写**——浮点按指针过 `SyscallN` 安全（388 立过的 XMM 雷只在写值路径），故 CGO_ENABLED=0 下可跑；
- 断言只钉仪器自身前提（5/5 tap 携带、麦克风 ~31 levels/s 流速）；判别句只打印不判案（AC#2 归编排者）。

## 2. 形①读数（两次独立真机，verbatim 在 `formone-gated-run{2,3}.log`）

run2（3 渲染后被 rig panic 打断）＋ run3（完整 **PASS 18.67s**）：

| 面 | run2 | run3 |
|---|---|---|
| `[volume]` 会话 ISimpleAudioVolume | （同） | **volume=1.000 muted=false** |
| `[volume]` 端点 IAudioEndpointVolume | （同） | **volume=0.200 muted=false** |
| 安静窗 mean（n≈188，31/s ✓） | 0.001088（max 0.034） | 0.002108（max **0.121＝环境突发声**——**麦听得见房间**） |
| 端点 tap（每次渲染） | 3/3 rms=640.06 peakrms=1459.29 | **5/5 rms=640.06 peakrms=1459.29** pushed=captured=48000 |
| 麦渲染覆盖窗 mean | 0.000759 / 0.000856 / 0.001046 | 0.001866 / 0.001735 / 0.001959 / 0.001288 / 0.001596 |
| 判别 | 全部 ≤ 安静底 | **best 比 0.93x**，全部 ≤ 安静底 |

`[verdict]`（run3 逐字）：the microphone did NOT hear the render (best 0.001959 = 0.93x quiet 0.002108, after-face mean 0.001467): **form ① SUPPORTED this run** -- the sound did not reach the capture device while the endpoint feed carried it (taps above); the fact question is WHY: session volume=1.000 muted=false, endpoint volume=0.200 muted=false, or physical speakers/headphones off -- the owner's ears settle it.

⇒ 形①的证据链补齐成三段：**端点喂了声（tap 5/5）→ 软件可读轴全部开放（会话 1.0/端点 0.2，均未静音）→ 麦没听到（渲染覆盖窗全部 ≤ 安静底）**。剩余嫌疑收敛到物理层（扬声器物理静音/没接、外放路由到麦听不到的物理路径），机主耳朵定。旁证：麦听得见房间突发声（0.121）却听不见渲染（≤0.005）——**麦路径本身是活的**。与 l1 六窗（sound 0.001358 < quiet-A 0.003003）同形，本腿用已知幅度件＋命名音量事实复现。

## 3. 附带发现（记给编排者，不是本腿的活）

- **rig 既有健壮性缺口**（run2 撞到）：`aecspike/rig_windows.go:707` `fill()` 在 GetBuffer 返回 S_OK 而 data=nil 时 `unsafe.Slice` 恐慌；registry 回收干净、仪器 skip 干净（log 逐字在 run2）。run3 未复现（5/5 过）。rig 是既有产码件，本波写面外。
- **423 遗留（槽主在飞）**：基线两发（昨晚 23:08）`go test ./cmd/wisp/` ＝ **build failed**（未跟踪件 `cmd/wisp/config_wire_423_windows_test.go` 引用不存在的 `h.rt`＋`EnableConversation` 返回值数不符，逐字在 `baseline-6runs.log`）；过夜后该件被槽主改过，今晨 cmd/wisp 可编译但红名册一片（run1 8 FAIL、run2 25 FAIL，两发均撞 **600s 包超时**——本腿两发漏了 `-timeout 30m`，记过）⇒ **cmd/wisp 门禁本波不可判**（树在变、红名册与 l1 的三枚在册集不同）。结构事实：`go list -deps -test ./cmd/wisp` 不含 aecspike ⇒ 本腿写面对 cmd/wisp 零影响（按构造）。
- **l1 遗留 CI 暴露（上报）**：见 §4 末条。

## 4. AC#5 执行面现跑（裁定 ⓔ 三跳，尺＝编译二进制＋符号表）

四份 list 现量在 `ac5-lists/`：windows `go test -list`：`./internal/audio/` 62 行、`./internal/audio/aecspike/` 34 行（**含 `TestAC297FormOneRenderReachProbe`**）。GOOS=linux 交叉 `-list` 无法在本机执行（fork/exec not a valid Win32）⇒ 改用 `go test -c` 交叉编译两枚 linux 测试二进制（rc=0）＋ `go tool nm` 列 Test 符号：**audio linux 48 枚、aecspike linux 18 枚，两枚 linux 二进制均无 297 家族符号**（windows build tag 挡住）。

CI 档 scope（推送树 `ddb759340` 现量）：core(ubuntu)＝`./internal/audio/`＋`./internal/audio/aecspike/`（portable-tests.sh:443）；windows＝`./internal/audio/`（:457，win_pin :392）；cli＝cmd/wisp。

⇒ **本腿新用例的执行面答案：零足迹（按构造）**——落点 aecspike：core 档的 linux 二进制不含它（tag），windows 档 scope 不含 aecspike ⇒ 与整个 gated rig 家族同形（rig_windows_test.go 头注的既有约定），只归机主/腿手动跑。

★ **l1 遗留暴露（三跳现跑坐实，本腿碰不得名册）**：l1 的 `TestAC297LivePacketByteFace` 落在 `./internal/audio/`（**windows 档 scope 内**）且 `WISP_LIVE_MIC` 门＝顶层 SKIP；推送树 `scripts/portable-tests.sh` grep `297`＝**0 命中**（名册只有 :807 `TestLiveWasapiSmoke` 一行）；`git branch -r --contains 0b5a19f34`＝origin/dev（已推）⇒ CI windows 档会对它**无账 SKIP＝strict runner fatal**（tools/d22scan/runtests.sh:24 规矩）。需要名册行（scripts/** 不在本腿写面）或 CI env，编排者定。

## 5. AC#3 默认档核验（零改动）

- `internal/config/schema.go`：`MicMutedDefault bool toml:"mic_muted_default" default:"true"`（现量 **:299**；票面引 :277——行号漂移、**值逐字未动**）。
- 本机 `%APPDATA%\wisp\` **无 config.toml**（只有 logs/、deps.toml）⇒ 默认档即生效档；票面「本机 config.toml 该节第 5 行 mic_muted_default = true」是立票时点状态，现机器该文件不存在（如实记录）。
- `level.go` "WHAT THIS TICKET DELIBERATELY DOES NOT DECIDE" 块原样（raw scale 定义，无加性底、无增益/压缩）。
- `internal/ball/liquid.go:30` `SilenceLevelGate = 0.06` 原样。
- `git diff --stat` 四文件＝空。

## 6. 门禁读数（每把门一行 rc；日志在 probes/297/l2/）

- `go test ./internal/audio/ -count=1`：改前 ok 18.080s / ok 18.140s（rc=0）；改后 ok 18.781s / ok 18.292s（rc=0）⇒ **新增红 0**
- `go test ./internal/audio/aecspike/ -count=1`：改前 0.188s / 0.187s（rc=0）；改后 0.323s / 0.248s（rc=0）⇒ **新增红 0**
- `go test ./cmd/wisp/`：基线两发 build failed（423 遗留）；改后两发撞 600s 超时＋红名册一片（槽主树在变）⇒ 本波不可判，见 §3
- `sh scripts/d22scan.sh` rc=0，`d22scan: clean - no D22 ban violations`（462 产码文件）
- `sh scripts/check-path-length-budget.sh` rc=0，`VERDICT GREEN`
- gofumpt -l cmd internal tools＝空 rc=0（`D:\work\base\gopath\bin\gofumpt.exe`）；gofmt -l＝空 rc=0
- `go vet ./internal/audio/...` rc=0
- scoped build `./cmd/... ./internal/... ./tools/...` **rc=0**；整树 `go build ./...` **rc=1**（gitignored `build/mut_rvw.go` 未定义符号——与 l1 同款预存遗留，两读数都落 `gate-build.log`）
- `tasklist /FI "IMAGENAME eq wisp.exe"` gated 真机两发起手＝0

## 7. 提交

pathspec＝`internal/audio/aecspike/formone_297_windows_test.go`＋`.scratch/wisp/probes/297/l2/**`＋`.scratch/wisp/issues/297-*.md`（Progress log 追加，AC 原文零改动）。只 commit 不 push。
