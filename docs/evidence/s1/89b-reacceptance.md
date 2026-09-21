# 票 89 复验（独立对抗验收 · 第二任）—— `acceptor-ticket89b`

**被验对象**：`.scratch/wisp/issues/89-0600-is-decorative-on-windows-acl-for-private-data.md`
（`Status: fix-complete-awaiting-reacceptance`）· 退回单四条 = 本文件同目录
`89-adversarial-acceptance.md`（上一轮 `acceptor-ticket89` 的面，本代理**不写它**）的"编排者·验收退回单"。
**被验的三枚 commit**：`c8d5c94`（退回单 #1 #2 判据+接线）/ `01e7007`（#3 #4 + 前三次变异读数）/
`f801d34`（#4 变异读数 + 全量门禁，纯 docs）。本轮 HEAD = `63fc82a`，
`git diff --name-only f801d34 HEAD -- internal/winsec internal/secret internal/memory internal/agent` **空**
⇒ 复验在 `f801d34`/HEAD 的码上取数即等价于被验三枚。

**环境**：Windows 11 Pro，`go version go1.27.1 windows/amd64`，当前账户 `DESKTOP-LVS7839\swq`
（SID 见下），`icacls` = `C:\Windows\System32\icacls.exe`。
**修复者的任何数字都不作为本代理的证据**：下面每一条都是本代理自己跑出来的，红名自己贴。

**快照纪律**：全部在仓外 `%TEMP%`（Git Bash 的 `/tmp`），仓内**未建 worktree**、工作树未被本代理改动。

| 快照 | SHA | 用途 |
| --- | --- | --- |
| `/tmp/wisp89bacc-89b2` | `c8d5c94` | 复现修复者变异 #1（删 `propagatePrivate`） |
| `/tmp/wisp89bacc-hd` | `63fc82a`(HEAD) | 四包 `-count=2` 回归 + 本代理在 HEAD 上重放变异 |

---

## 第 1 条：AC#2 覆盖面判据搬进包内 —— **复现成立（红名唯一、其余全绿），但票面 AC#2 的表述要降一档**

### 1.1 我自己重做那发变异（`c8d5c94`，与修复者同一 SHA、同一锚点）

- 锚点：`internal/winsec/winsec_windows.go:293` 的 `return propagatePrivate(path)` → `return nil // MUTATION-89BACC-WALK`
 （`applyDescriptor` 一行未动，walk 整体留着但没人调它）。
- **同链证明**：`grep -n "MUTATION-89BACC-WALK"` ⇒ `293:	return nil // MUTATION-89BACC-WALK`（就是锚点那一行）。
- **先证可编译**（编译失败不算变异）：变异前 `go build ./...` **rc=0**，变异后 `go build ./...` **rc=0**。
- 跑：`go test -count=1 -v ./internal/winsec/` ⇒ **TEST_RC=1**。

| 读数 | 本代理实测 |
| --- | --- |
| `=== RUN` | **25** |
| 顶层 `--- PASS` | **14** |
| 顶层 `--- FAIL` | **1** |
| 缩进 `    --- PASS`（子测试） | 10 |
| `SKIP`（全量输出 grep -c） | **0** |

**唯一红名**：`--- FAIL: TestAC2SealDirNarrowsChildrenThatCarryTheirOwnExplicitACEs (0.45s)`
⇒ 与修复者自称的"rc=1、`=== RUN` 25、14 绿 / 1 红、0 SKIP、红名就是这一条"**逐字对上**。

红因（OS 自己的原文，本代理从日志逐字取）：

```
acl_windows_test.go:454: icacls text of operator-widened-file.txt still names the foreign principal (Everyone):
    …\TestAC2SealDirNarrowsChildrenThatCarryTheirOwnExplicitACEs…\data\operator-widened-file.txt Everyone:(RX)
        NT AUTHORITY\SYSTEM:(I)(F)
        BUILTIN\Administrators:(I)(F)
        DESKTOP-LVS7839\swq:(I)(F)
acl_windows_test.go:455: NOT PRIVATE operator-widened-file.txt: foreign SID(s) S-1-1-0
    (principals: [Everyone NT AUTHORITY\SYSTEM BUILTIN\Administrators DESKTOP-LVS7839\swq])
```

`operator-widened-dir` 与 `operator-widened-dir/inside.txt` 同形红（三条都带显式 ACE）。

**还原自证**：把快照的该文件用 `git show c8d5c94:…` 覆盖回去后
`grep -rn "MUTATION-89BACC" --include=*.go .` **0 命中**、`diff` 与 `c8d5c94` 逐字相同、
仓内 `git diff --quiet -- internal/winsec/winsec_windows.go` **rc=0**（工作树未做变异）。

### 1.2 变异在 HEAD 上打到哪几格（本代理补测，修复者没报这一格）

（待补：HEAD 快照同锚点重放，见 §1.2 结尾）

### 1.3 判"注释里那句'纯继承的孩子不区分任何事'"—— **是诚实，不是自我缩窄**

（判定见 §1.3 结尾）
