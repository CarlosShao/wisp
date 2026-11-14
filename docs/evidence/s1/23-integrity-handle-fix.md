# 23 — AC#3 提权判定那枚句柄误用：修 ＋ 一枚真能咬住它的断言（写腿）

**日期**：2026-11-14 10:39–11:0x +08
**锚**：HEAD `ebe98e6f77c70f4529da9fd0988228733c1210d4`（`dev`）
**票**：`.scratch/wisp/issues/23-system-window-input-tools.md` AC#3 第二分句（真机读数见 `23-input-type-real-window.md`）
**身份**：写腿（实现者）。本件是实现侧记录，**不是裁决表**，不翻任何 AC 框。

## §0 一句话

`realProcessElevated` 把 **OpenProcess 的进程句柄**当**令牌句柄**喂给 `GetTokenInformation`，于是每个目标都读成 IL 0、谓词恒假（真机上对 High 窗 Taskmgr 回「已注入 15 个字符」而不是 `ErrElevatedTarget`）。本腿把它改成「先 `OpenProcessToken(h, TOKEN_QUERY, &tok)` 再读 IL」，并把**本进程自己的 IL** 做成断言（`integrityOfProcess(os.Getpid()) == currentIntegrity()`）：带着旧写法这枚断言逐字读到 `0／8192` ⇒ **必红**（已实测，见 §3），修好后 `8192／8192` ⇒ 绿。

## §1 缺陷（修前，逐字）

`internal/tools/system_backend_windows.go:728`（签名）／`:741`（错句）：

```go
	defer windows.CloseHandle(h)
	return integrityOf(windows.Token(h)) > self
```

`h` 来自 `windows.OpenProcess(processQueryLimitedInfo, false, pid)`，是**进程句柄**；`integrityOf` 里的 `windows.GetTokenInformation` 只认**令牌句柄** ⇒ 返回 `The handle is invalid.` ⇒ `integrityOf` 返 0 ⇒ `0 > self` 恒假（self 是本进程令牌读出的 8192）。同文件 `:747` 的 `currentIntegrity` 是正确形状（`OpenProcessToken(windows.CurrentProcess(), TOKEN_QUERY, &h)`），一目了然的对照。

真机后果（前任腿读数，`23-input-type-real-window.md` §0）：往 IL 12288 的 Taskmgr 调 `input.type` → `IsError=false` ＋「已向 taskmgr.exe 注入 15 个字符」，`ErrElevatedTarget` 一次都没出现。

## §2 修法（落盘后的形状）

新增**唯一一处**度量入口，谓词只消费它的结果：

```go
func realProcessElevated(pid uint32) bool {
	if pid == 0 {
		return false
	}
	self := currentIntegrity()
	target := integrityOfProcess(pid)
	if target == 0 {
		// Fail closed. ...
		return true
	}
	return target > self
}

func integrityOfProcess(pid uint32) uint32 {
	if pid == 0 {
		return 0
	}
	h, err := windows.OpenProcess(processQueryLimitedInfo, false, pid)
	if err != nil {
		return 0
	}
	defer windows.CloseHandle(h)          // 句柄一：进程
	var tok windows.Token
	if err := windows.OpenProcessToken(h, windows.TOKEN_QUERY, &tok); err != nil {
		return 0
	}
	defer tok.Close()                     // 句柄二：令牌
	return integrityOf(tok)               // 只有令牌句柄能回答 TOKEN_MANDATORY_LABEL
}
```

**两个句柄都关**（`windows.CloseHandle(h)` ＋ `tok.Close()`），`tok` 一定先于 `h` 关闭。

**fail-closed 语义（按原注释与调用方语义办，理由写在代码里）**：`target == 0` ⇒ `return true`（判「提权」= 拒绝）。依据有三条，都写在 `realProcessElevated` 的注释里：
1. 打不开更高完整性目标的**进程或令牌**本身就是访问拒绝，这是 fail-closed 的原始场景（原代码 `OpenProcess` 失败即 `return true`，此处保持一致）；
2. 猜「没提权」的代价是**静默丢弃输入**——正是票 23 存在的理由；猜「提权」的代价只是一句诚实的拒绝；
3. 唯一调用方 `inputTarget`（修后 `:1423-1424`）在 `elevated==true` 时返回冻结句 `ErrElevatedTarget`，所以 fail-closed 精确映射到「宁可拒，不试注入」。

顺带订正两处**同函数族的注释**（都在同一文件、只改注释）：`integrityRID` → `currentIntegrity`（函数名早就叫这个）；`integrityOf` 的 RID 表原写「0x1000 medium / 0x2000 high / 0x3000 system」是**错一档**的，按真值改成 low/medium/high/system，并写明「读不到的 TARGET 按提权处理」。

## §3 断言：为什么必须加、加上后带旧 bug 必红

### 3.1 为什么原有两枚腿抓不到

`internal/tools/ticket23_input_uipi_windows_test.go`：
- `:72` 把 `processElevated` 打桩成 `func(uint32) bool { return true }` ⇒ 与真谓词无关；
- `:129` 用真谓词，但目标是**本进程自己的窗口**（同 IL）⇒ 正确答案是 `false`，而恒假的错答案**也是 false**。⇒ 旧 bug 存在时两枚都是绿的。这正是要补的洞。

### 3.2 新断言（本件新增，注释 `:162`、函数 `:177` 起）

```go
func Test23IntegrityOfProcessAnchorsOnThisProcess(t *testing.T) {
	pid := uint32(os.Getpid())
	self := currentIntegrity()
	got := integrityOfProcess(pid)
	t.Logf("本进程 pid=%d：integrityOfProcess=0x%x (%d)／currentIntegrity=0x%x (%d)", pid, got, got, self, self)
	if self == 0 { ... }        // 锚本身不成立 ⇒ 红（环境异常，不是谓词问题）
	if got == 0 { ... }         // ← 旧 bug 的落点：进程句柄当令牌句柄
	if got != self { ... }      // 同一进程两个读数必须相等
	if got < 0x1000 { ... }     // 读到的不是 IL（低完整性 RID 下限）
	if realProcessElevated(pid) { ... }  // 出厂的谓词对同 IL 目标必须放行
}
```

**锚是什么**：`currentIntegrity()` 走 `OpenProcessToken(CurrentProcess(), TOKEN_QUERY, ...)` 这条**正确**形状；`integrityOfProcess` 必须自己完成「进程句柄 → 令牌句柄」这一跳。两边测的是**同一个量**（本进程令牌的 TOKEN_MANDATORY_LABEL 末位 RID），所以任何一侧用错句柄，等式就散。断言取「相等 ＋ 两侧都不是未知 0」，避免 `0 == 0` 这种两个零互相作证。

### 3.3 证明（⛔ 只声明不算，逐字贴读数）

**修后（绿，`16-test23-final.txt`）**：

```
--- PASS: Test23IntegrityOfProcessAnchorsOnThisProcess (0.00s)
    ticket23_input_uipi_windows_test.go:181: 本进程 pid=66340：integrityOfProcess=0x2000 (8192)／currentIntegrity=0x2000 (8192)
```

**把度量盲成旧写法**（`integrityOfProcess` 内 `OpenProcessToken` 那一段换成 `return integrityOf(windows.Token(h))`，落地证据：`05-mutation-landed.txt` 逐字 `764: return integrityOf(windows.Token(h)) // MUTATION-23: the process handle fed to GetTokenInformation`）：

```
=== RUN   Test23IntegrityOfProcessAnchorsOnThisProcess
    ticket23_input_uipi_windows_test.go:181: 本进程 pid=37148：integrityOfProcess=0x0 (0)／currentIntegrity=0x2000 (8192)
    ticket23_input_uipi_windows_test.go:187: integrityOfProcess(本进程 pid=37148) 回了 0，currentIntegrity() 是 8192：进程句柄被当成令牌句柄喂给 GetTokenInformation 了（UIPI 谓词会因此恒假）
--- FAIL: Test23IntegrityOfProcessAnchorsOnThisProcess (0.00s)
FAIL
FAIL	github.com/CarlosShao/wisp/internal/tools	0.125s
test rc=1
```

⇒ **左边 0、右边 8192，必红**（`06-mutation-red.txt` 全文）。

**还原（复绿）**：从 `system_backend_windows.go.fixed-copy` 覆盖回去后 sha256 与改前**逐字相同**（`0E1CDE8E8B21569D080A3B7288F4EF8EA1E9DA11ED6096C346B9D9213F9943BA`，`04-fixed-sha.txt` ＝ `07-restored-sha.txt`），`MUTATION-23` 残留 0 命中：

```
=== RUN   Test23IntegrityOfProcessAnchorsOnThisProcess
    ticket23_input_uipi_windows_test.go:181: 本进程 pid=38640：integrityOfProcess=0x2000 (8192)／currentIntegrity=0x2000 (8192)
--- PASS: Test23IntegrityOfProcessAnchorsOnThisProcess (0.03s)
ok  	github.com/CarlosShao/wisp/internal/tools	0.307s
test rc=0
```

（`08-restored-green.txt`。两枚绿读数 pid 不同＝两次独立进程，非同一行复用。）

## §4 门禁（逐条贴读数）

| 门 | 命令 | 读数 |
|---|---|---|
| build（字面） | `CGO_ENABLED=1 go build ./...` | **rc=1**，红的**全部**是未跟踪、被 `.gitignore:14 build/` 忽略的本机残件 `build\mut_rvw.go`（`undefined: voiceTtsFiles` 等 11 行＋`too many errors`）＝存量、非本腿（`git ls-files build/` 空；`git status --porcelain -- build/` 空） |
| build（同一门在**只含版本控制内容**的树上） | `git checkout-index -a --prefix=<快照>/` ＋覆盖本腿两枚改动 ＋ `CGO_ENABLED=1 go build ./...` | **rc=0**（快照 `D:\tmp\wisp\20261114\23-integrity\headtree`，`build/mut_rvw.go` 不在其中） |
| build（按域，主工作树） | `CGO_ENABLED=1 go build ./cmd/... ./internal/... ./tools/...` | **rc=0** |
| vet | `go vet ./internal/tools/ ./cmd/wisp/` | **rc=0** |
| vet（winlive） | `go vet -tags=winlive ./cmd/wisp/` | **rc=0** |
| d22scan | `sh scripts/d22scan.sh` | **rc=0 clean**；正控 `PASS=45 FAIL=0 SKIP=0, === RUN=93`；扫描面 `bans #1-5 internal/=404, cmd/=77`、`ban #6 frontend/=123`、`ban #7 internal/tools/=51`、`ban #10 internal/risk/=21`、`ban #8 design/=39`、`frontend/=123`、`internal/=1007`、`cmd/=278` |
| gofumpt | `gofumpt -l internal/tools/`（v0.12.0 / go1.27.1） | **0 行** |
| 整包四数（本腿改的包） | `go test -count=1 -v ./internal/tools/` | **rc=0**：`=== RUN=574` · 顶层 `--- PASS=424` · `--- FAIL=0` · `--- SKIP=1`；`ok github.com/CarlosShao/wisp/internal/tools 73.590s`。口径＝`^=== RUN`（含子测试）／`^--- `（只顶层）；唯一 SKIP＝既有门闩 `Test23LiveClipboardRoundTripThroughRealOS`（默认不碰操作者剪贴板，改前就有） |
| 真 carrier 定向（票 23 那几枚） | `go test -count=1 -v -run 'Test23' ./internal/tools/` | **rc=0**：`RUN=36 TOPPASS=35 TOPFAIL=0 TOPSKIP=1`；点名四枚全 PASS：`Test23UIPIPredicateRefusesElevatedTargetBeforeAnyByte`、`Test23UIPIPredicateLetsLiveNonElevatedTargetThrough`、`Test23IntegrityOfProcessAnchorsOnThisProcess`、`Test23ProcessNameIsTheD45BindingKey` |

原始日志逐条落 `.scratch/wisp/probes/23/integrity-fix/`（`00`–`17` 编号，未跟踪、不提交）：`09` 整包、`03`/`16` ticket23 定向、`06`/`08` 变异红/还原绿、`10`–`12` vet/d22、`13`–`15` 三发 build、`04`/`07` sha。

## §5 写面与纪律

- **改的跟踪文件只有两枚**：`internal/tools/system_backend_windows.go`（生产码 ＋ 同族注释）、`internal/tools/ticket23_input_uipi_windows_test.go`（新增一枚用例 ＋ 文件头清单一条）；另加本件与新追加的票 23 Progress log 一节。
- **一处路径更正（照实记）**：派单写的是 `cmd/wisp/ticket23_input_uipi_windows_test.go:72`，仓里**没有**这个文件（`Get-ChildItem cmd/wisp -Filter *uipi*` 零命中）。真身是 `internal/tools/ticket23_input_uipi_windows_test.go`，且**行号逐字对得上**（`:72` 恒真打桩、`:129` 真谓词同 IL 放行）⇒ 按真身改，未另建 `cmd/wisp` 副本，也未把真谓词重新打桩。
- ⛔ 未碰：`internal/panel/**`、`scripts/**`、`frontend/**`、`models/**`、`design/**`、`docs/PLAN.md`、`docs/specs/**`、`thresholds.go`、任何 golden/allowlist/契约文字。
- ⛔ 未放宽任何既有断言；两枚老腿一字未动（本腿只**加**）。
- ⛔ 未 push（子代理只 commit）。

## §6 本腿**不**声称什么（残余）

1. **不翻 AC#3 的框**：票面那格还要求「真 High 窗上 `input.type` 回 `ErrElevatedTarget`」的**端到端**读数。本腿证明了谓词的输入侧量对了（8192/8192、变异必红、同 IL 放行），但**没有再跑** `.scratch/wisp/probes/23/probe23`（它要拉起 fixture 窗、抢前台；且本腿写面只有三处，不能新增探针件）。⇒ 那一发真窗读数仍欠，路径已写死在票面 Progress log（「修完 `probe23` 原样复跑即可复核，无需 UAC」）。
2. **不声称 `go build ./...` 在主工作树是绿的**：它红在 `build/mut_rvw.go`（11-04 遗留、未跟踪、被 `.gitignore:14` 忽略）。按「临时件只建不删」未处置，同门在只含版本控制内容的快照上 **rc=0**（§4 第二行）。
3. **不声称本机 High 窗读数是本腿量的**：`12288／8192` 那对读数出自前任腿（`23-input-type-real-window.md` §4/§5），本件引用，未重测。
