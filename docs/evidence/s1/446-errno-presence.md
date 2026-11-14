# 票 446 — errno 存在性由「英文句子」判改为「码」判：证据

**票**：`.scratch/wisp/issues/446-errno-presence-is-decided-by-an-english-error-sentence.md`
**裁决人（本文件作者）＝实现者**（写腿）。对抗验收/缺口审计须由**另一个** agent 做（`SPEC-12 §4.3` #1/#3），本文件不代替它。
**修复 commit**：见本文件所在 commit（`git log --format=%H -1 -- docs/evidence/s1/446-errno-presence.md`）；基线 = `447929bfe`。
**改动文件**：`internal/tools/system_backend_windows.go`（唯一产码文件）
**修复后 sha256**：`8CAD968D67AE4F3D4795CDBDB136A6A28E48AB93B44A76FEF716D61026738845`（AC#3 的还原判据就是这个串）
**时刻**：2026-11-14 18:1x–18:3x +08

一句话：**这个 errno 到底有没有，现在是拿码判的（`errors.Is(e, syscall.Errno(0))`），不再拿一句英文判**；两处调用点走同一个判定函数，句子在这个文件里已经不存在。

---

## AC#1 现量复核：全量 132 处，其中判存在性 2 处

尺（可逐字重跑，`HEAD` = 修复前那棵树）：

```
git grep -n  'operation completed successfully' HEAD          # 全部命中（含测试件）
git grep -c  'operation completed successfully' HEAD -- '*.go'   # 1 个文件 / 2 行
git grep -c  'operation completed successfully' HEAD -- ':!*.go' # 130 行
git grep -n -E 'Err(or)?\(\)\s*(==|!=)\s*"' HEAD -- '*.go'      # 同形族，见下表
```

| 类 | 命中行数 | 文件数 | 判「判存在性」还是「只打印」 |
|---|---|---|---|
| **A. 产码 `*.go`** | **2** | 1 | **判存在性，2 处，全部改了** |
| B. probe/CI 日志落盘（`.scratch/**/logs`、`.scratch/ci-logs`） | 121 | 112 | 只打印（`lastErr=…`、`err=…` 回显，无分支） |
| C. 已落档证据正文（`docs/evidence/s1/216-ac5-branch-by-code.md`、`33-panel-host-c27-r9.md`） | 5 | 2 | 只引用（散文） |
| D. 票面正文（`issues/216-…`、`issues/446-…`） | 4 | 2 | 只引用（散文） |
| **合计** | **132** | 117 | 判存在性 **2**，只打印 **130** |

**A 类逐条**（修复前）：

| 行 | 原判据 | 判定 |
|---|---|---|
| `internal/tools/system_backend_windows.go:286` | `if r == 0 && e != nil && e.Error() != "The operation completed successfully." {` | 判存在性（EnumWindows 失败分支） |
| `internal/tools/system_backend_windows.go:426` | `if e != nil && e.Error() != "The operation completed successfully." {` | 判存在性（SendInput 少送分支，决定要不要把 errno 包进错误） |

**同形族普查**（`Err(or)?().*(==\|!=)\s*"`，全仓 `*.go`，`HEAD`）：**只有上面这 2 处**，没有第三处；`*_test.go` 命中 **0**。
⇒ M=2，改了 2，没有「同形漏改」。

---

## 前提核查：票面那半句「非英文 Windows 上永远判错」——**机制成立，但本机走不到**

票面写的是「`syscall.Errno(0).Error()` 恰好是这句英文 ⇒ 非英文 Windows 上比较永远不成立」。这条**不是靠猜**，本腿把它追到了源码行并量了文本表：

`$(go env GOROOT)/src/syscall/syscall_windows.go:165-180`：

```go
func (e Errno) error() string {
	// ask windows for the remaining errors
	var flags uint32 = FORMAT_MESSAGE_FROM_SYSTEM | FORMAT_MESSAGE_ARGUMENT_ARRAY | FORMAT_MESSAGE_IGNORE_INSERTS
	b := make([]uint16, 300)
	n, err := formatMessage(flags, 0, uint32(e), langid(LANG_ENGLISH, SUBLANG_ENGLISH_US), b, nil)   // :169 先要 en-US
	if err != nil {
		n, err = formatMessage(flags, 0, uint32(e), 0, b, nil)                                     // :171 失败才退回系统默认语言
		...
```

**读法**：这句英文是「向系统要 en-US 表」要来的，**不是** Go 里的静态串。只有 en-US 表拿不到时才会退回系统默认语言 ⇒ 那时 `Error()` 返回本地化文本，字符串比较**永远不成立** ⇒ 判「有 errno」⇒ 走错分支。这就是票面说的那个 bug。

**本机实测（`go run .scratch/wisp/probes/446/r1/errno-locale-probe.go`，日志 `probes/446/r1/04-locale-reading.log`）**：

```
GetUserDefaultUILanguage   = 0x0409        （界面语言：en-US）
GetSystemDefaultLangID     = 0x0804        （系统区域：zh-CN）
FormatMessage(errno=0, en-US 0x0409)      = "The operation completed successfully."
FormatMessage(errno=0, zh-CN 0x0804)      = "操作成功完成。"
FormatMessage(errno=0, langid 0 = default)= "The operation completed successfully."
syscall.Errno(0).Error()                  = "The operation completed successfully."
syscall.Errno(5).Error()                  = "Access is denied."
```

**诚实结论（三点，分开说）**：

1. **本机上旧代码是「碰巧对」**：这台机器的界面语言是 en-US（0x0409），`:169` 那句要 en-US 要到了，所以 `Errno(0).Error()` 就是那句英文，字符串比较在本机成立。⇒ **没有人在这台机器上看得到这个 bug**，这就是它能活到现在的原因。
2. **换一台机器就错**：只要 en-US 消息表拿不到（本地化 Windows 未装 en-US 语言包就是这种），`:171` 退回系统默认语言，同一个 errno 0 的文本变成 `操作成功完成。`（上面第三/四行已经把这台机器上的本地化文本量出来了）⇒ 旧比较判「有错」。
3. **本腿没有真机复现那一支**：把 `:169` 弄失败需要一台没有 en-US 表的机器，本机做不到。所以本文件里「退回分支被走到」是**源码行 + 本地化文本**两条读数推出来的，**不是**一台非英文 Windows 上的实跑。这一条留给 AC#5 的 CI（ubuntu/windows runner 的界面语言也不是本机这台）或对抗验收去坐实，本腿**不声称**已复现。

**不管机制在哪台机器上触发，判据本身都不该是一句话**：`errors.Is` 除了不看文本，还会**拆包**（`fmt.Errorf("…: %w", …)`）并认 `Is` 方法，字符串比较只看得见最外层的文本。这一条是**本机能跑红**的（见 AC#3 B 段）。

---

## AC#2 改法：两处都改，判定收成一个函数

```go
// errnoPresent reports whether e carries a real Windows error code.
func errnoPresent(e error) bool {
	return e != nil && !errors.Is(e, syscall.Errno(0))
}
```

- `:301`（原 286）`if r == 0 && errnoPresent(e) {`
- `:441`（原 426）`if errnoPresent(e) {`
- 新增 import：`errors`、`syscall`（`staging_live_windows.go:36` 早就在用 `errors.Is(err, syscall.Errno(...))`，不是新风格）

**顺手记下的一条**：旧写法的第一个子句 `e != nil` 其实**恒真**——`Proc.Call` 的 lastErr 是 `syscall.Errno(0)`，装进 `error` 接口后不是 nil 接口。所以旧判据整条的分量全压在那一句英文上；现在压在码上。

---

## AC#3 能变红的检查（一发，两个突变方向，都点名）

器具（**不在 `internal/tools/` 里留测试文件**，用 `go test -overlay` 把探针挂进 `package tools`）：

- 探针：`.scratch/wisp/probes/446/r1/probe446_errno_presence_windows_test.go`
- 叠加：`.scratch/wisp/probes/446/r1/overlay-446.json`
- 静态闸：`.scratch/wisp/probes/446/r1/gate-446-static.ps1`

```
go test ./internal/tools -run TestProbe446 -count=1 -v -overlay .scratch/wisp/probes/446/r1/overlay-446.json
```

探针 10 格，其中 **5 格在 errno=0 这条路上**：`nil`、`errno-zero-raw`、`errno-zero-via-windows-alias`、`errno-zero-localized-text`（**文本是 `操作成功完成。` 的同一个 errno 0**，就是票面那台机器的形状）、`errno-zero-wrapped`、`errno-zero-wrapped-twice`；再加真错误 3 格（`5`/`87`/`87 包一层`，都要判「有」）与 1 格「读起来像成功文本但不是 errno 0」（旧比较会把它吞成「没 errno」）。

**三段实测（日志 `probes/446/r1/32-ac3-sequence.log`，逐字可重跑）**：

| 段 | 树 | sha256 | 结果 |
|---|---|---|---|
| **A 修复后** | 判据 = `errors.Is(e, syscall.Errno(0))` | `8CAD96…38845` | **PASS**（`probe446 GREEN: 10 cases, 5 of them on the errno=0 path`） |
| **B 把判定临时改回字符串比较** | `return e != nil && e.Error() != "The operation completed successfully."` | `AE4F32…E178` | **FAIL，且点名**：`PROBE 446 RED at case "errno-zero-localized-text": errnoPresent("操作成功完成。") = true, want false` |
| **C 还原** | 与 A 同 | `8CAD96…38845`（**与 A 同串**） | **PASS** |

**突变二（调用点不走向判定）**：把 `:441` 改回内联的句子比较（判定函数本身不动）⇒ 行为探针**仍然绿**（它测的是判定，不是调用点）。所以另配了静态闸，**这一发专门抓这个方向**：

```
pwsh -NoProfile -File .scratch/wisp/probes/446/r1/gate-446-static.ps1
```

| 树 | 读数 | 结果 |
|---|---|---|
| 修复后 | `sentence-hits=0 text-compares=0 errnoPresent-defs=1 errnoPresent-calls=2 code-judgements=1`，两个调用点实名 `line 301` / `line 441` | **GREEN rc=0** |
| `:441` 回退成内联比较 | `sentence-hits=1 text-compares=1 errnoPresent-calls=1` | **RED rc=1**，三行点名：`sentence still decided on at line(s): 441` / `error TEXT compared to a literal at line(s): 441` / `expected exactly 2 call sites … found 1 at line(s) 301` |

⇒ 「只写改了」在这里是拿不出来的：判定回退→行为探针红并点名；调用点回退→静态闸红并点名；两次还原都验了 sha256 同串（`8CAD96…38845`）。

**没做到的一格（说清楚）**：两处调用点（`realEnumWindows` 的 `r == 0`、`sendInputs` 的 `n != len(ins)`）**没有**用集成手法把真实系统调用推到「少送事件且 `GetLastError()==0`」那一支——那要么往活的桌面注入输入、要么假造 `LazyProc`，两者都会把被测的东西换掉。所以「分支必须被走到」是在**唯一判定点**（`errnoPresent`）上构造 errno=0 走到的，不是在 Syscall 返回处走到的。

---

## AC#4 门禁

| 门 | 命令 | 读数 |
|---|---|---|
| `go build` | `go build ./internal/tools` | **rc=0** |
| `go vet` | `go vet ./internal/tools` | **rc=0**（日志 `probes/446/r1/40-vet.log`） |
| `gofumpt -l` | `gofumpt -l internal` | **没装，本腿没跑**（`Get-Command gofumpt` 空；`$(go env GOPATH)\bin` 空；`$(go env GOMODCACHE)\mvdan.cc` 不存在 ⇒ 离线装不上，CI 那一步是 `go install mvdan.cc/gofumpt@v0.12.0` 联网装的）。**这一格留空给 CI/编排者，不拿别的冒充。** |
| （**不是** gofumpt 的旁证） | `gofmt -l internal/tools` | 空输出 ⇒ 本包的 `gofmt`（只是 gofmt）干净。**替代不了 gofumpt，别当 AC#4 过。**（日志 `probes/446/r1/41-gofmt-NOT-gofumpt.log`） |
| 额外（不在 AC#4 里） | `go test ./internal/tools -count=1` | **rc=0**（`ok … 76.4s`，日志 `probes/446/r1/60-package-tests.log`） |
| `d22scan` | `sh scripts/d22scan.sh` | **clean，rc=0**：`examined 483 production Go files under internal/ and cmd/`；`clean - no D22 ban violations`（日志 `probes/446/r1/50-d22scan.log`） |
| path-length | `sh scripts/check-path-length-budget.sh --with-self-test` | **GREEN rc=0**：`denominator: tracked paths=20122 over-budget=69 covered by roster=69 not in roster=0`，正控 3/3 PASSED（日志 `probes/446/r1/51-path-length.log`） |
| 写面自检 | `git status --short -- internal/tools/` | 只有 `internal/tools/system_backend_windows.go` 一个 M |

`d22scan` 跑的那棵树里，本文件 sha256 = `8CAD96…38845`，与最终串**相同**（突变两次都已还原，见 AC#3 C 段）。

---

## AC#5 CI

**留白。** 本腿没有 CI 权限（不推、不触发、读不到 run），**不声称** `test-core`/`test-windows` 绿。
要填的三样（`run 号` + `git merge-base --is-ancestor <fix> <head>` + 完成时刻）**全部留空**，由编排者在带上本 commit 的推送后填。

---

## 没覆盖 / 风险（留给裁决者）

1. **真机复现缺口**：见上文「本腿没有真机复现那一支」——en-US 表缺失的那台机器不在本腿上。AC#3 的 B 段用**本地化文本 + `Is` 模型**把同一失败形状跑红了，但它证明的是「判据看文本会错」，不是「这台机器上会错」。
2. **句子只剩一处 `*.go` 命中**：`probes/446/r1/probe446_errno_presence_windows_test.go` 里那句英文是**构造的测试输入**（造一个「读起来像成功文本」的错误），不是判据。`internal/tools/` 内 `git grep 'operation completed successfully' -- '*.go'` = **0 命中**。
3. **未动的同形**：全仓 `*.go` 没有第三处 `Err(or)?().*(==|!=)\s*"`；`*_test.go` 里也没有——如果哪天测试里出现同类断言，这条尺要在 CI 里常驻（本腿只留了 `gate-446-static.ps1`，没有往 `tools/d22scan` 里加规则，那是别的腿的写面）。

---

## 撤销口令

「**票 446 撤**」＝删本票 + 对修复 commit 发一个新 commit 反改（把 `errnoPresent` 两处调用点与函数回退成句子比较）＋删 `docs/evidence/s1/446-errno-presence.md` 与 `.scratch/wisp/probes/446/`。⛔ 不 reset/amend。
