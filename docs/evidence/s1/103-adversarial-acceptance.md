# 票 103 独立对抗验收裁决表 —— 密封缝守卫（R-c，今天可利用）与 `RemoveUnlinked` 祖先检查 + tripwire（R-b，陷阱档）

**验收代理：** `acceptor-ticket103`（独立对抗，未参与实现）
**验收日期：** 2026-09-21
**被验收的码：** 修前红用例 `184af22`；生产码 `0717bf2`（`internal/winsec/resolve.go`、`winsec.go`、`winsec_windows.go`、`winsec_other.go`）；顶端文档 commit `8b6f691`
**票面：** `.scratch/wisp/issues/103-seal-seam-is-unguarded-and-removeunlinked-follows-junctions.md`
**仓库状态：** branch `dev` @ `8b6f691`；`git status --porcelain` 显示 `internal/config/parse.go`、`migrate.go`、`internal/risk/syncdirs*_test.go`、`.github/workflows/ci.yml` 是**别人的未提交活** ⇒ 本表所有 Go 命令一律**按包**或在 `/tmp` 快照里跑（快照 `C:/Users/swq/AppData/Local/Temp/wisp-ac103` = `git archive HEAD`，目录名带会话后缀，仓内未建 worktree）。

**来源档位图例**
- 〔独立复现〕= 本代理亲自敲的命令、亲自读到的 rc 与原文行。
- 〔日志＋归档，我抽验〕= 实现方日志给了读数，本代理抽验了其中可核的部分（commit 存在性、文件内容、断言行是否真在树里）。
- 〔仅自述，不背书〕= 只有实现方日志，本代理未能复现；**补救动作**逐条写明。

---

## 裁决总表（AC 1:1）

| AC | 判据（票面原句要点） | 本代理主证据 | 结论 | 档位 |
|---|---|---|---|---|
| AC#1（R-c） | 缝只能装一次 + 装的必须是真解析器；两条腿；"什么都没发生"要 SID 级读数 | 见下 §AC1 | **通过但有条件**（R-103-1） | 〔独立复现〕 |
| AC#2（R-b） | junction 输入 ⇒ 拒 + 不删目标；memory 侧"下降即红"tripwire | 见下 §AC2 | **通过但有条件**（R-103-2） | 〔独立复现〕 |
| AC#3 | 变异三向（去守卫红／被拒仍有一条绿／tripwire 改下降红） | 见下 §AC3 | 待补（进行中） | 〔独立复现〕 |
| AC#4 | 四包回归 + gofmt/gofumpt + 按包 vet + `sh scripts/d22scan.sh` | 见下 §AC4 | 待补（进行中） | 〔独立复现〕＋〔仅自述，不背书〕 |

---

## §AC1 —— R-c：`SetPathResolver` 的三道守卫

### A. 判据原句
> seam 只能被**装一次**、且装的必须是**真解析器**……判据用例两条腿：**装第二个 ⇒ 红**；**装一个恒说 OK 的 ⇒ winsec 必须拒**，且"什么都没发生"不算绿（要能证明**外来 DACL 没被改**，取 SID 级读数）。

### B. 形状保留（不许倒回 `filepath.Abs`）
- 命令：`grep -rn "SetPathResolver" --include=*.go .` ⇒ 生产侧唯一安装点 `internal/risk/winsec_c26.go:21`（`init()`）。
- `internal/winsec/resolve.go:198-227` `ResolvePath` 仍是唯一 mint；`privateDirAll` 仍收 `ResolvedPath`（`winsec.go:160`）。
- `grep -n "filepath.Abs" internal/winsec/*.go` ⇒ 无（见 §命令记录）。
- 结论：**形状在、守卫加了**，票 94 的账没清零。〔独立复现〕

### C. 两条腿 + 守卫幂等性
（逐条读数见下节，随验收进行补全。）

### D. P1 探针：守卫有没有留别的门（本代理自己造的攻击，非实现方用例）
（待填：恒改写型伪造。）

---

## §AC2 —— R-b：`RemoveUnlinked` 祖先检查 + `internal/memory` tripwire

（待填。）

---

## §AC3 —— 变异三向

（待填。）

---

## §AC4 —— 门禁复跑

### 已量到的 CI 步级读数（run id 必报）
`gh api repos/CarlosShao/wisp/actions/runs?per_page=6` ⇒ 顶端 run **`35587986855`**（run_number **148**，`head=b2fa2ed`，**含 `0717bf2` 的 103 生产码**，created 2026-09-21T10:17:51Z，`status=completed conclusion=failure`）。
`gh api repos/CarlosShao/wisp/actions/runs/35587986855/jobs --jq ...` ⇒ 步级：
- `JOB test-core => failure`：`Portable package tests (agent/llm/config/memory/observe/secret/risk/statemachine/...) => **failure**` ⇒ **POSIX 侧（含 `internal/risk` 与 `internal/memory`）确实跑了，且是红的**（细节待填）。
- `JOB test-windows => failure`：`Portable windows tests (proc/secret/config) => failure`。
- `JOB lint => failure`：`D22 scanner positive control => success`、`D22 seven-ban + emoji scan => success`、`gofmt (gofumpt) => success`、`go vet (module) => success`、`staticcheck => failure`。
- `slo-smoke`/`slo-full`/`lint-frontend` 全 success。

---

## 对实现方点名的三处复核要求（票面 `next=`）的直答

### ① `%VAR%`／前导 `~` 的**改写型**伪造仍归票 102 的 `Actable()`
（待填。）

### ② POSIX 侧只有编译期验证 ⇒ 本票最大未验面
（待填：见 §AC4 的 run 148 步级读数。）

### ③ 祖先检查 fail-closed ⇒ 契约从严还是可用性回归
（待填。）

---

## 缺陷登记 `R-103-x`（验收期不修）

（待填。）

---

## 命令记录（原文输出，按调用顺序）

（待填。）
