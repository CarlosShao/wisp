# 票 424 — Linux 编译断裂修复（腿码 424-fix）

- 锚：`a25ce0a5f5fd901510959d3aab08cd16a07cdbd1`（进场现读）
- 时刻：2026-11-12 10:46–10:54 +08:00
- 腿码：**424-fix**
- CI 逐字病灶：`vet: cmd/wisp/memory_retention_180k.go:73:38: undefined: residentVoice`

## 改了什么（逐字）

只两枚代码文件，函数体零改动：

1. **新建 `cmd/wisp/memory_retention_180k_windows.go`**（首行 `//go:build windows`），全文：

```go
//go:build windows

package main

import (
	"fmt"
	"log/slog"
)

// attachMemoryL1IfEnabled180K is the resident boot's l1_enabled gate, the
// addressable seam the wiring test drives: enabled runs the ticket-29 attach,
// disabled leaves the WarmWindow hook body-less and says the key decided
// that. It runs where the assembly root used to call attachMemoryL1
// unconditionally (resident_windows.go).
//
// Windows-only by build tag: it is the one retention wire whose parameter
// list names *residentVoice, a type that only exists behind the windows
// build (resident_voice_windows.go). Everything else in the 180-census
// retention wire is cross-platform and stays in memory_retention_180k.go.
func attachMemoryL1IfEnabled180K(rv *residentVoice, run *agentRuntime, dataDir string) {
	l1, src := memoryL1Enabled180K(dataDir)
	if !l1 {
		slog.Info("memory: L1 extraction not attached",
			"reason", "[memory] l1_enabled=false", "config_source", src)
		fmt.Printf("wisp: 记忆抽取未接（[memory] l1_enabled=false，来源 %s）：Warm 窗口不挂抽取管线，画像不学习\n", src)
		return
	}
	rv.attachMemoryL1(run)
}
```

   （原注释逐字照搬；仅追加第二段注释说明为何挂 windows build 约束——不属函数体，无断言/契约改动。）

2. **`cmd/wisp/memory_retention_180k.go`**：只删 `attachMemoryL1IfEnabled180K` 一个函数（连同其 5 行专属注释，共 16 行）。diff：

```diff
--- a/cmd/wisp/memory_retention_180k.go
+++ b/cmd/wisp/memory_retention_180k.go
@@ -65,19 +65,3 @@
 	return c.Memory.L1Enabled, cfgPath
 }

-// attachMemoryL1IfEnabled180K is the resident boot's l1_enabled gate, the
-// addressable seam the wiring test drives: enabled runs the ticket-29 attach,
-// disabled leaves the WarmWindow hook body-less and says the key decided
-// that. It runs where the assembly root used to call attachMemoryL1
-// unconditionally (resident_windows.go).
-func attachMemoryL1IfEnabled180K(rv *residentVoice, run *agentRuntime, dataDir string) {
-	l1, src := memoryL1Enabled180K(dataDir)
-	if !l1 {
-		slog.Info("memory: L1 extraction not attached",
-			"reason", "[memory] l1_enabled=false", "config_source", src)
-		fmt.Printf("wisp: 记忆抽取未接（[memory] l1_enabled=false，来源 %s）：Warm 窗口不挂抽取管线，画像不学习\n", src)
-		return
-	}
-	rv.attachMemoryL1(run)
-}
-
 // retentionBootLine180K is the one audit sentence the retention start books,
 // lifted out so the feed test can assert the same numbers the log carries.
```

其余 3 个跨平台函数（`retentionWindowDays180K` / `memoryL1Enabled180K` / `retentionBootLine180K`）与 `l3RetentionDefaultDays` 原地保留。`attachMemoryL1IfEnabled180K` 体内只调了跨包可见的同 package `memoryL1Enabled180K`（留在原文件），无需搬私有 helper。

## 判据

本机是 Windows 工具链（mingw，无 linux 头文件），完整交叉编译必然死在 `runtime/cgo`（`sys/mman.h` 等——本机工具链问题，非本票）；GOOS=linux 裸 `go vet` 又会先死在 `sherpa-onnx-go-linux` 的 build constraints（CGO_ENABLED=0 时该包零文件）。所以判据用 **GOOS=linux type-check 面的两向对照**证明（`go vet -overlay` 虚拟替换 go.mod → 指向 D:\tmp 里的纯 Go sherpa 桩，只在本机 D:\tmp\wisp\20261112\424fix\，不进仓库）：

| 项 | rc 或读数 |
|---|---|
| win build `go build ./cmd/wisp` | **rc=0** |
| win vet `go vet ./cmd/wisp` | **rc=0** |
| **linux type-check（修后树）** | **不再出现 `undefined: residentVoice`**（日志 `D:\tmp\wisp\20261112\424fix\linuxvet-typed.log`，零 residentVoice 行） |
| **linux 阴性对照（overlay 塞回 HEAD 版原文件）** | **逐字复现 CI**：`vet.exe: ...memory_retention_180k.go.head:73:38: undefined: residentVoice`（`linuxvet-head-negctl.log`）——证明该仪器真能走到 type-check、真能抓这个错 |
| win test `-run 'Test180K\|Test423\|TestTicket255' -count=1` | **RUN=34 PASS=33 FAIL=0 SKIP=1**（rc=0；`win-test.log`） |

裸 GOOS=linux vet（判据第 ④ 步原样，`linuxvet.log`）在 sherpa 处 rc=1，日志无 `undefined: residentVoice`；但阴性对照证明它死在 import 处够不到 type-check，故判据以 overlay type-check 两向对照为准（更强）。

## ⚠ 顺手发现（同族第二枚，不在本票射程，报编排者）

修后树的 linux type-check 露出了**下一个**错误（vet 只报第一枚，之前被 residentVoice 排队遮住）：

```
vet.exe: cmd\wisp\logsink.go:224:80: undefined: configErrText
```

`configErrText` 定义在 `cmd/wisp/resident_audio_windows.go:475`（windows-only），却有 **2 个跨平台调用点**：
- `cmd/wisp/logsink.go:224`
- `cmd/wisp/memory_retention_180k.go:62`（本票保留的 `memoryL1Enabled180K` 体内）

即 Linux 侧修完 residentVoice 后仍会挂在这两处——需另立票处理（搬 helper 或加平台桩）。本票铁律禁动其他文件，未动。

## 越界自证

`git status --porcelain`：本腿只有 3 个落点——` M cmd/wisp/memory_retention_180k.go`、`?? cmd/wisp/memory_retention_180k_windows.go`、`?? docs/evidence/s1/424-fix-linux-build.md`（本件）。工作树里其余大量 `??` 均为别的腿既有未跟踪件，一字未动、未提交。`git diff --stat` 原文件＝16 行纯删除。

## 一句话

只把 `attachMemoryL1IfEnabled180K` 一枚函数搬进新建的 `memory_retention_180k_windows.go`（首行 `//go:build windows`），Windows build/vet/test 全绿，GOOS=linux type-check 两向对照（阴性复现 CI 逐字 / 阳性零 residentVoice）证明 `undefined: residentVoice` 在 Linux 侧消失；另发现同族第二枚 `undefined: configErrText`（logsink.go:224 等 2 处跨平台调用点）待另票，未越界。
