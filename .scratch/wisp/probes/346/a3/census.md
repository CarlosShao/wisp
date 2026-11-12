# 346-a3 AC#4 颜色与门禁

- 腿码 **346-a3**（执行腿：只跑测试取读数）。日期 **2026-11-12 +08**。
- 起手锚（第 0 步现跑）：`a7580348a70d62a7b836212ab96b602bc6626a89`；跑数期间编排者于 **09:57:04** 落了新锚 `3561f37a8063500424fc17f86fef7c86fc5c655e`（只加 `docs/evidence/s1/346-ac4-four-runs.md`）。
  ★**两锚之间 `internal/winsec` 零字节差**（尺＝`git diff --stat a7580348 3561f37a8 -- internal/winsec`＝空；两枚仪器 blob＝`dfaddac8…`/`0b0e5173…`、`winsec.go`＝`cbafe712…`，`git rev-parse`⇄`git hash-object` 逐串核对）⇒ 本腿全部读数对两枚锚同样成立。
- 时刻（读数落盘区间，本机 Asia/Shanghai）：Windows 档 win1＝09:44:39–09:46:05、rep1–3＝09:48:37–09:50:37；Linux 容器档 lin1–3＋ole1–3＋psw1–3＝09:55:52–09:56:23。**全部先于新锚落盘**（除 lin 链尾 09:56:23 早于 09:57:04，同样成立）。
- PATH 带 DLL：`third_party\sherpa-onnx` 与 `build` 两目录前置且存在（尺＝`Test-Path` 双 True）；全程未出现 `0xC0000135`。
- Docker：**可用**，`docker --version`＝`Docker version 29.6.2, build dfc4efb`（rc=0）；镜像 `golang:1.27-bookworm`（容器内 `go1.27.3 linux/amd64`；宿主 `go1.27.1 windows/amd64`；`go.mod` 要 `go 1.27` + toolchain `go1.27.1` ⇒ 容器版高于 toolchain 声明下限，toolchain 切换被 GOFLAGS=-mod=mod 下正常走完，见下）。
- 门禁前提逐条照办：跑前 **`Remove-Item Env:GOTMPDIR / Env:TMPDIR`**（宿主）＋容器内 **`env -u GOTMPDIR -u TMPDIR`**（票 311 AC#1 shape B GUARD，带 env 时 rc=2 射程拒绝⛔ 读成红——本腿每发都剥干净，0 发 rc=2）。
- ⛔ 改任何 `.go`：**0 字节**（尺＝`git status --porcelain -- internal/winsec` 全程空）。
- ⛔ 跑 `./cmd/wisp` 整包：**没跑**（读数＝本件只含 `./internal/winsec/` 九发＋容器预热 `-run '^$'`）。
- ⛔ `go test ./...`：**没跑**（同上，逐发命令都带显式包名）。
- ⛔ push：没推（只 commit）。

## 三句弃权红句逐字（顺带记；AC#2 的判据是它们还在——本腿现读源码逐字抄，三句都在）

1. `internal/winsec/openleaf_322_other_test.go:206`（源码为多行拼接，206–208 行）：
   `AC#1 RED instrument: across %d raced attempts the leaf leg (%q) answered %d and the racer landed %d flips (%d errors, first %q), so the damaged count above is not evidence about a window that was never reached`
   （票面 2026-10-18 CI 实例：`...across 3000 raced attempts the leaf leg ("too many levels of symbolic links") answered 0 and the racer landed 6 flips (0 errors, first "")...`）
2. `internal/winsec/posix_seal_window_120_other_test.go:233`（233–234 行）：
   `AC#1 RED instrument is dead shape=%s: refused=%d, sealed-own=%d - a zero window reading here proves nothing`
3. `internal/winsec/posix_seal_window_120_other_test.go:240`（`t.Logf`）：
   `AC#1 shape=%s: every refusal came from the check before the act (refused-at-act=0), so this run does NOT exercise the pinned identity gate`

顺带：本腿 winsec 树上另见大量 `t.Skipf` 形弃权（符号链接权／短名／卷根等台面原因，346-a1 普查已列 12 处/8 文件），不在本票三句射程内，本腿未动、未抄全。

★**两枚弃权仪器的编译面**（诚实边界，读数自取）：两枚仪器首行都是 **`//go:build !windows`**（本腿现读），⇒ **本机 Windows 档根本不编译它们**（既不绿也不弃权，与编排者 `A1053` 口径一致）⇒ **「两枚弃权仪器各自 3 发独立复跑」只能在 Linux 容器侧做**，本腿照此办（见下）。

## Windows 本机 一发

命令：`go test ./internal/winsec/ -count=1 -timeout 30m -v`（`D:\work\base\go\bin\go.exe`，`Start-Process -PassThru` 读 `$p.ExitCode`；读数件 `win1.out`/`win1.err`，err＝0 字节）。四数按派单尺（顶格、剥 `/` 子测试、`Sort-Object -Unique`）。

| RUN | PASS | FAIL | SKIP | rc | 红名 | SKIP 名 |
|----|----|----|----|----|----|----|
| 78 | 78 | 0 | 0 | 0 | 无 | 无 |

（Windows 档不编译两枚弃权仪器 ⇒ 「红⇄绿⇄弃权」这一维在 Windows 档不适用，如实记。）

## 三发独立复跑（同一台面＝同一宿主、同一命令、逐发独立 go test 进程、`-count=1`）

| 发 | RUN | PASS | FAIL | SKIP | rc | 红/绿/弃权 | 红名 |
|---|---|---|---|---|---|---|---|
| rep1 | 78 | 78 | 0 | 0 | 0 | 全绿（无红、无弃权行） | 无 |
| rep2 | 78 | 78 | 0 | 0 | 0 | 全绿（无红、无弃权行） | 无 |
| rep3 | 78 | 78 | 0 | 0 | 0 | 全绿（无红、无弃权行） | 无 |

散布：**四数与 rc 逐发全同＝散布 0**；输出行数 1609/1608/1607（差＝go 的时延行，非测试行）。逐名色变＝无（78 枚逐发同绿）。
旁证：编排者 `3561f37a8` 里上一条腿的四发 Windows 读数（78/78/0/0 ×4）与本腿独立同向。

## Linux 容器那半（成对读数）

台面：`git archive HEAD`（起手锚 `a7580348`）→ `D:\tmp\wisp\20261112\346a3\tree` → `docker run -v … golang:1.27-bookworm` 容器内 `-w /tree`。
完整性自证：导出树 `internal/winsec` **57/57 枚文件齐**（缺 0）；两枚仪器与 `winsec.go` 的 `git hash-object` 逐串等于 HEAD blob（`dfaddac8…`/`0b0e5173…`/`cbafe712…`）⇒ 容器跑的就是 HEAD 的字节。容器内跑前 `env -u GOTMPDIR -u TMPDIR`，逐发独立 `go test` 进程、`-count=1 -timeout 30m`，读数件 `lin*.out`/`ole*.out`/`psw*.out`（err 全部 0 字节）。

**① Linux 档整包 3 发**（与 Windows 档同一包、同一 `-count=1 -timeout 30m`）：

| 发 | RUN | PASS | FAIL | SKIP | rc | 红名 | SKIP 名 |
|---|---|---|---|---|---|---|---|
| lin1 | 48 | 48 | 0 | 0 | 0 | 无 | 无 |
| lin2 | 48 | 48 | 0 | 0 | 0 | 无 | 无 |
| lin3 | 48 | 48 | 0 | 0 | 0 | 无 | 无 |

**② 两平台成对名册作差（逐名点名，四数逐枚、⛔ 混用⛔ 相减）**：

- Windows 名册 78 枚、Linux 名册 48 枚；**交集 6 枚**（`TestAC2AncestorPrefixesForEverySeparatorShape`、`TestAC2ComponentsAndTraversalPerSeparatorShape`、`TestAC2OverwritePromiseHoldsForEveryProductionCaller322r1`、`TestAC3WindowsTheRouteCallsTheLeafGuardBeforeThePin322r1`、`TestAC5FailedSealRefusesTheWrite`、`TestAC5FailureIsNotSwallowedByTheHappyPath`）——**两侧同绿**。
- 仅 Windows 72 枚：**全部**落在 `*_windows_test.go`（逐名归位到 24 个文件，最大 `seam_bypass_108_windows_test.go` 6 枚）；仅 Linux 42 枚：**全部**落在 `//go:build !windows` 侧文件。
- ⇒ **名册差 100% 由 build tag 解释**（Go 的跨平台测试面本来就不对称），**无一枚"同名测试两平台翻色"**。成对判读＝两侧各自全绿＋交集 6 枚同绿＋两枚仪器在 Linux 侧全绿。

**③ 两枚弃权仪器各自 3 发独立复跑**（容器内单独 `-run` 逐发独立进程；逐名点名红⇄绿⇄弃权）：

| 仪器 | 发 | 判定 | rc | 仪器读数（散布如实报） |
|---|---|---|---|---|
| `TestAC1POSIXOpenThroughASwappedLeafLinkNeverTruncatesAForeignFile322r1`（openleaf_322） | ole1 | **PASS（绿）** | 0 | attempts=3000 flips=908 flipErrs=0 damagedWitnesses=0 leafLegRefusals=224 |
| 同上 | ole2 | **PASS（绿）** | 0 | attempts=3000 flips=648 flipErrs=0 damagedWitnesses=0 leafLegRefusals=143 |
| 同上 | ole3 | **PASS（绿）** | 0 | attempts=3000 flips=1100 flipErrs=0 damagedWitnesses=0 leafLegRefusals=241 |
| `TestPOSIX120SealWindowLeafSwapIsClosed`（posix_seal_window_120） | psw1 | **PASS（绿）** | 0 | attempts=2000 refused-at-check=912 refused-at-act=12 sealed-own=965 hit-foreign=0 |
| 同上 | psw2 | **PASS（绿）** | 0 | attempts=2000 refused-at-check=940 refused-at-act=17 sealed-own=936 hit-foreign=0 |
| 同上 | psw3 | **PASS（绿）** | 0 | attempts=2000 refused-at-check=955 refused-at-act=15 sealed-own=923 hit-foreign=0 |

- 6 发 **红 0、弃权 0（无 `--- SKIP`、无弃权红句触发）、绿 6**；判定（色）散布＝0。
- 仪器读数的**数字**散布如实报：ole flips 908/648/1100、psw refused 912/940/955——竞速仪器命中数本就随机，判定全绿即"窗口到达且门全部拦住"。
- psw 三发均见 `other-shape example`（同屏 `winsec: cannot apply a private security descriptor … openat blob.bin: path escapes from parent`）＝仪器第三味在跑；`refused-at-act` 12/17/15 ＞0 ⇒ "每一发都真到达了 act 那一刻"（红句③那句条件 `refused-at-act=0` 未触发）。

## 红名册作差逐名归位

Windows 档 4 发（win1+rep1–3）∪ Linux 档 3 发（lin1–3）∪ 仪器 6 发（ole1–3+psw1–3）＝**13 发逐发 `--- FAIL` 名册全空** ⇒ **红名册＝∅，无一枚需要归位**（票 346 描述的 PASS⇄FAIL 翻脸在本腿台面上 0 发复现）。

## 一句话总判

AC#4＝**成对**（Windows 本机 ⇄ Linux 容器：同一包、同一 `-count=1 -timeout 30m`、四数逐枚点名、交集 6 枚两侧同绿、名册差全部由 build tag 解释；两枚弃权仪器各 3 发独立复跑——因仪器是 `//go:build !windows` 只能在容器侧——全绿、数字散布如实报）；红名册＝**空（0 枚红）**。

## 原始读数落点（只建不删）

- 宿主发：`.scratch/wisp/probes/346/a3/win1.{out,err}`、`rep1..3.{out,err}`（err 全 0 字节）。
- 容器发：自 `D:\tmp\wisp\20261112\346a3\out\` 原样复制到 `.scratch/wisp/probes/346/a3/lin/`（lin1–3、ole1–3、psw1–3 各 `.out/.err`，err 全 0 字节）。
- 容器链日志：`D:\tmp\wisp\20261112\346a3\chain.log`（仓外）。
