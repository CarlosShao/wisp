# config 收尾腿 `cfg125-r1` — 修 `TestAC2POSIXSeamInstallSurvivesASymlinkSpelledTemp125`（门禁新红＝347 落地默认拒的正确后果）

- 锚＝HEAD `d36059b44`（2026-10-27 18:4x +08；起手时 HEAD 曾为 `70f020a92`，期间编排者 A1283 落 `pending-and-issues.md` 一件，与代码面无关——`git show --name-only d36059b44` 现量只有那一件）。
- 起手核：`git diff --cached --name-status` **空**；`git status --porcelain -- internal/config internal/risk cmd/wisp internal/winsec` **空**。共享工作树有 `.scratch/**` 等他人在飞件＝非本腿写面、零触碰。
- 本腿写面＝`internal/config/c26_seam_posix_125_test.go`（一枚件）＋`docs/evidence/s1/cfg125-r1-impl.md`。**`internal/risk`、`cmd/**` 零字节**（§6）。日志全在 `D:\tmp\wisp\20261026\cfg125r1\logs\*.log`（只建不删）。

## 1. 修法选型（读数 1）：**② 的强化形**；① 落不下去（读数有据）

**红因**（与 347-r1 判定一致，机理未重查、只取读数）：`reparseComponents` 落地后，AC#1 钉子的**行为探针**输入 `<子进程TMPDIR>/../wisp-125-probe` 的祖先 `varlink125` 是 symlink ⇒ `Resolve` 正确返回 `ErrReparseDenied`（`D:\tmp\...\347r1\linux-gate-core.log:1275` 逐字）⇒ 钉子要求"非错" ⇒ 红。**安装本身没失败**：`resolver installed` 在两个子进程里都打出来了。

- **①（注入 `[fs] reparse_point_exceptions`）⇒ 落不下去**：解析端在（`internal/config/schema.go:489`），但 **seam 读取端不存在**——`risk.c26Pipeline.Resolve` **逐字硬编码 `Resolve(input, nil)`**（`internal/risk/winsec_c26.go:40`），该文件头注释逐字："**Reparse exceptions are deliberately nil** … For placement, **default-deny is the only policy**"。给子进程喂这个键＝要么改 `internal/risk`（**硬边界 1**）要么推翻 seam 刻意的默认拒（＝编排者说的"须另裁"）。配置的生产消费者是另一条缝（`cmd/wisp/run.go:544` → `tools.NewPathCanonicalizer`），且测试子进程根本不装载 config。
- **② 选用的形**：**被测子进程 TMPDIR 保持穿链**（主题不动），**只把 AC#1 探针 base 锚定到真拼写**（`filepath.EvalSymlinks`；t.TempDir 结果必存在 ⇒ 它就是 realpath 的等价物；无链 base 上是恒等）。AC#2 fixture 的 base 同样锚定。
- **对另一条的代价**：③（认"被拒"为证据）编排者已否（最弱、放宽断言）；朴素②（子进程换 plain temp）会把被测子进程变成第二个 control——用例自己的 fixture 守卫（`:139-141`"no link ⇒ this is the control shape, not the measured one"）会当场 Fatal，等于为绿架空用例。强化形②避开两者。

## 2. "主题仍在量"（读数 3，逐字）

- **原本量什么**：AC#2 的失败句逐字（`:160-162`）——"in a process whose temp dir is spelled … through the link … the sealing seam does not hold C26. **The claim under test is winsec.PathResolverInstalled()'s type**"。即：**当 TMPDIR 拼写穿过 symlink 时，安装决策（risk init → seam 装 C26）依然成立**——ticket119 实测过的"装不上"形状。
- **改后量什么**：**同一件事**。被测子进程的 `TMPDIR` 仍逐字设成 `…/varlink125/tmproot125`（`:191` 读完数逐字"child TMPDIR=…/varlink125/tmproot125 rc=0"）；唯一变化是钉子**行为探针的输入拼写**不再顺路穿过那条祖先链（`base` 先 `EvalSymlinks` 锚定，`hostile = base + "/../wisp-125-probe"`）。
- **为什么等价**：改前这条探针的红是**祖先链的默认拒**（对任何经此链的输入都会发生、与 C26 在不在场无关——floor 与 C26 都拒它），不是"主题"的读数；改后探针回到它原本问的形状（`..` 穿越必须被折改、不得拒绝、不得残留 `..`）——**三条断言（非错／不原样传回／不残留 `..`）一条未删**。ticket119 那个"装不上"的世界在改前改后都被同一台仪器盯着（`PathResolverInstalled()` 为 nil ⇒ AC#1 第一发 Fatal）。此外 fixture 锚定还把 darwin `/var→/private/var` 的既有 skip 世界（`:120`）变成真测量（skip 条件收紧、测量面变大）。

## 3. 两平台颜色（读数 2）

- **Linux**（golang:1.27 容器、导出树、`unset GOTMPDIR`，复刻 347-r1 跑法）：**改前** `04` 前基线：AC#2 **红**（measured child rc=1，`00-baseline-125-target.log`）→ **改后 rc=0 双子全绿**（`04-postfix-125-target.log`：control rc=0、measured rc=0、两个子进程都逐字打出 `winsec: sealing path resolver installed`）。**Linux 必须绿：到手。**
- **Windows**：该文件 `//go:build !windows`，不编译进 Windows 测试二进制 ⇒ **改前 `go test ./internal/config/` rc=0（`02-win-config-before.log`）＝改后 rc=0（`03-win-config-after.log`）**，"改前改后都绿"到手（本就绿、零退化）。

## 4. 反向判据（读数 4）：突变三读数 ＋ 还原核（全在仓外导出树）

- **变异红（M1＝把修法摘掉：HEAD 原版文件放回 export-mut）**：`08-mut-M1-fix-reverted.log` rc=1、FAIL 恰 1 枚＝`TestAC2POSIXSeamInstallSurvivesASymlinkSpelledTemp125`，红句逐字（AC#2 RED…rc=1）⇒ **修法承重**。
- **还原**：fix 放回，SHA256 逐字节核 `mut=DF80B283…` ＝ `worktree=DF80B283…` **SAME=True**（mut.ps1 输出）。
- **还原核（时机＝最后一发变异 M1 之后）**：`08-restore.log` rc=0、`--- FAIL/SKIP` 计数 **0**、两枚 PASS。
- **原样绿**＝`04-postfix-125-target.log`（**从未变异**的同一棵 export-mut 上同联合筛 rc=0）。

## 5. 门禁级 `--scope=core`（Linux 容器、无 git 导出树）改前/改后对照（读数 5）

| 树 | 四数 | rc | 红集 |
|---|---|---|---|
| 基线（HEAD 逐字节＝工作树非 doc 面） | `PASS=1393 FAIL=10 SKIP=1`（`=== RUN=2052`） | 1 | `AC2…Symlink…125`＋`NothingWritten…`＋`HandDeleted…`＋`Dsn196r3…`＋`C26RewriteAccount…`＋panel 4 枚＋`ReadGit…` |
| 本腿 | `PASS=1394 FAIL=9 SKIP=1`（`=== RUN=2052`） | 1 | 上面 **−1 枚**（config 那枚消失）＋`HandDeleted…` 轮换回来 |

- **集合差＝恰 −1 枚**＝`TestAC2POSIXSeamInstallSurvivesASymlinkSpelledTemp125`（499 处 diff 逐名核对，无第二处差异）。`SKIP=1` 两发同枚（`TestCanonicalizeErrorNotice…174r3` 账本行）＝既有。
- **`TestHandDeletedSchemaVersion…` 是容器内已知轮换**：本腿故障探针（`05-flakiness-probe.log`：export-base 2 发中 1 发、export-mut 2 发中 0 发）＋**347-r1 自己的基线门禁日志 `D:\tmp\…\347r1\linux-gate-core-base.log` 逐字有它**、而其本腿发没有 ⇒ 本腿之前就在轮换；本腿单文件测试改动不触及 writeguard 机理（零产码改动）。归"容器既有"，⛔ 当新红。
- 包级 `./internal/config`（容器）：改前 `3 FAIL`（AC#2＋writeguard 2 枚，`00-baseline-config-pkg.log`）→ 改后 `1 FAIL`（`NothingWritten…`＝存量，`04-postfix-config-pkg.log`）⇒ **config 红 −1、−1**。

## 6. 四数（读数 6）

| 对象 | 改前 | 改后 |
|---|---|---|
| `internal/config`（包级 Windows） | rc=0 全绿 | rc=0 全绿 |
| `internal/config`（包级 Linux 容器） | 3 FAIL（含目标 1） | 1 FAIL（目标消失） |
| `internal/risk` | Windows 包级唯一红＝存量 A380 `TestC26RewriteAccount…`（`02-win-risk-before.log`） | **逐字同一枚**（`03-win-risk-after.log`）；`git diff --name-only -- internal/risk/` **空＝零字节** |
| `cmd/wisp` | — | `git diff --name-only -- cmd/` **空＝零字节** |

## 7. 格式甲乙两尺 ＋ d22scan（读数 7）

- **甲 `gofmt -l`／乙 `gofumpt -l`**（绝对路径 `/d/work/base/gopath/bin/gofumpt.exe`，全仓 tracked `*.go`，归属判 `[\\/]`）：两条尺都只报那 **3 条既存**（`cmd/wisp/models.go`、`internal/agent/approval/pending_read.go`、`internal/agent/tools.go`，CRLF/`autocrlf` 造的、⛔ 本腿责任）＋`.scratch/**` 探针件；**本腿那枚件两尺都不在名单**（单件直跑 rc=0 双尺）。
- **d22scan**（`tools/d22scan` 独立模块 `go build` 后扫导出树，复刻 347-r1 跑法）：**rc=0 clean**（`06-d22scan.log`：bans #1-5 internal/=263、cmd/=45、ban #6 frontend/=116、#7 tools/=25、#8 design=14/frontend=116/internal=682/cmd=158）。**正控**＝塞一枚真 `≤`（字节 `e2 89 a4`，`od` 现量）进 `internal/risk/zz_d22_positive_control_cfg125.go` 字符串 ⇒ **rc=1 且逐字点名** `[emoji] ban #8 … string literals included` ⇒ 尺是活的。（第一次正控未触发＝`printf '\u2264'` 按字面 6 字节落盘，`charcheck` 字节级核实后重打——过程留日志。）

## 8. 提交与门禁 rc

- 门禁 rc：Linux `core` 基线 rc=1 ／ 本腿 rc=1（差集＝§5 那 −1 枚）；两平台包级 rc 见 §6。
- 提交：`add`/`commit` 两处**逐个文件名字面 pathspec**（`internal/config/c26_seam_posix_125_test.go` ＋ `docs/evidence/s1/cfg125-r1-impl.md`）；`-F` 写在 `--` 之前、正文用 `[System.IO.File]::WriteAllText` 真临时文件；主题逐字以「config 收尾腿 cfg125-r1」开头；⛔ push。
- 提交后复跑（按新 HEAD 再导出）：目标对 rc=0 ＋ config 包红集＝§5 同一枚存量。
- 中途上报：修法定案＋读数已发编排者（本会话 send_message，2026-10-27 18:3x）。

## 9. 未核／具名残余（⛔ 藏）

1. ⛔ 未跑 macOS 实机；fixture 锚定让 darwin 的 skip 条件更难触发但 darwin 门禁数未取。
2. 台账（`pending-and-issues.md`）与票面更新＝编排者权（本腿 ⛔ 翻框）；`A1282①` 的"①/② 任选＋主题具名"由本件 §1/§2 具名交付。
3. writeguard 两枚红的轮换归因是**旁证链**（347-r1 基线日志＋本腿探针），⛔ 机理级归因（不属本腿写面）。
4. `winsec seam probe root` 那半（ticket125 的 winsec 侧）本腿零触碰。
