# 133-v（非实现者复核腿）：锚、环境、窗口、快照与漂移

- **身份**：非实现者（实现者＝133-r1）。本文件只登记读数与形状，判语在
  `docs/evidence/s1/ticket133-133-v-accept-v1.md`。
- **锚**：起手 20:51:42 `git rev-parse --short HEAD` = `51e2f5ea`；20:52:57 现读 = **`b75089cc`**
  （编排者"133-r1 收件"提交在这一分钟内落地 ⇒ 漂移已逐报）。本腿一切读数取自 `b75089cc` 的
  `git archive` 快照；两枚仪器文件在 `da389ca5` 之后零改动（`git diff --name-only da389ca5 HEAD -- cmd/wisp/`
  为空）。131 的枚举门文件最后动于 `6f702aeb`，不在本票三枚提交的任一文件表里。
- **仓工作树状态（只读自证）**：开工时 `git status --porcelain cmd/wisp/` 空；21:2x 复核时出现
  `?? cmd/wisp/task_output_pointer_notice_174r8_test.go`（本波另一条腿的在飞笔，非本腿产物，零触碰）。
  ⛔ 本腿**未在仓内做任何变异**：五发＋AC#3/AC#4 全部落 `D:/tmp/wisp133v/base`（仓外快照）。
- **环境**：`GOTMPDIR='C:\tmp\wisp133vgotmp'`（先自建，只反斜杠形）；`PATH` 前置
  `/d/tmp/wisp133v/base/third_party/sherpa-onnx`（MSYS 绝对形——相对形会让子进程 0xc0000135，
  第 1 趟全包读数据此归因过，见 03）；`go1.27.1 windows/amd64`；快照＝`git archive b75089cc`，
  全日志在 `D:/tmp/wisp133v/logs/`（只建不删）。
- **tasklist 窗口自证**（cmd/wisp 重包串行的自证，`go.exe|test.exe|compile.exe` 零命中；
  `cargo.exe` 会因子串 "go.exe" 误配，逐条已肉眼归因；`wisp.exe`＝机主常驻球，热键占位者＝宿主底噪源）：
  - **W0 20:57:11**：无 go/test/compile；`wisp.exe` PID 40296 ＋ 2 枚 `cargo.exe`（宿主持续项）。
  - **W1 21:11:06 → 21:21:19**：字面 `go test -count=2 -v ./cmd/wisp/`（第 1 趟，见 03）。
  - **W2 21:25:58**：无 go/test/compile；2 枚 `cargo.exe` 仍在。
  - **W3**（字面命令第 2 趟，补 `third_party` 后）窗口见 `D:/tmp/wisp133v/logs/14b-literal-rerun-window.txt`。
- **仓外快照树**：`base`＝`b75089cc`；`pred`＝`d01f4b71`（＝`da389ca5^`，仪器落地前）。
  快照里 `third_party/` 不入 git（`.gitignore:12`）⇒ 第 2 趟起手工拷入三枚 DLL，posture 修正的凭据：
  修前 `TestAC3EarlyRecordLandsOnDiskBeforeTheInstallRecord` 红 `no native DLLs in ..\..\third_party\sherpa-onnx`
  ＋子进程 `0xc0000135`；修后同一发 **PASS（10.69s）**。
- **纯净 sha256（还原锚）**：
  - `cmd/wisp/main.go` = `24efbffbee6eeb2719b114dbb9241da1c738b4e32544480760f2069b5076e0a3`
  - `cmd/wisp/leg_dispatch_gate_133_test.go` = `1273fdcc53098f78fcd62c397ad2830163692b2562fbbe29d4a65d9f18dbf1c5`
  - `cmd/wisp/dispatch_hop_selfproof_133_test.go` = `cc19ac942565845d593cd2ecdaf4390dbde624926754de7d59c10e803a541c81`
  每发变异后逐发还原、sha256 复量相等（读数见 01/02）。
- **变异清单（本腿自造，形状与 r1 探针不重名不重文件）**：
  1. N-3：`main.go` 摘掉 `os.Exit(cmdModels(...))` 一行（换 `os.Exit(2)`）。
  2. X4：新增早退 `if len(args) > 0 && args[0] == "--hopv133"` ＋新件 `hopv133v.go`（`cmdHopV133` 调 `installLogSink`）。
  3. X8：`case "doctor":` → `case doctorCmdV133:`（新件 `cmdlabels133v.go` 一枚 const）。
  4. X12：新件 `hopalias133v.go`（`var hopAliasV133 = installLogSink`）＋ `case "hopalias133"`。
  5. X14 第一拍：新件 `vhold133v.go`（`holderFieldV133{openFn: installLogSink}` ＋ `vHolder133.openFn("")`）＋ `case "vhold133"`；
     第二拍：删字段初值与那次调用，腿仍在（`cmdHoldV133` 空身）＋ case 保留。
  6. AC#3：`isRunnableCase133` 去掉 `fd.Recv != nil` 半边（≠r1 的 `return true` 顶置）。
  7. AC#4：注册表腿名换靶 1 行；usage 加一枚 `phantomcmd` 行；`run`↔`models` 两枚 claim 的 test 名对调 2 行；
     删整枚 `case "doctor"` 子句（y1drop）；`args := os.Args[1:]` → `args := os.Args`（y2blind）。
