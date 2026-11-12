# CI `lint` job 实读数 — staticcheck 第一次产出真判据（票 85 AC#1/AC#2 ＋ 票 122 AC#6 的公共读数）

- 采集：**2026-11-12 10:2x +08**，编排者用 `gh` 现取（⛔ 子代理转述）
- 数据源：`gh run view 43139602748` ＋ `gh api repos/CarlosShao/wisp/actions/jobs/130923107457`
- 起手 HEAD：**`a6d203613`**（该 run 的 `headSha`）
- 原始日志落点：`D:\tmp\wisp\20261112\ci\failed.log`（5,990,996B）

---

## 1. ★该 run 的 job 一览

| job | 结论 | job id |
|---|---|---|
| **lint** | **failure** | `130923107457` |
| lint-frontend | success | `130923107719` |
| slo-smoke | success | `130923107742` |
| slo-full | success | `130923107793` |
| test-windows | failure | `130923107815` |
| test-core | failure | `130923107840` |

## 2. ★★`lint` job 的 16 个 step（**step 号是票 122 AC#6 要的"step 号"**）

| # | step 名 | 结论 |
|---|---|---|
| 1 | Set up job | success |
| 2 | actions/checkout@v4 | success |
| 3 | actions/setup-go@v5 | success |
| 4 | D22 scanner positive control (tools/d22scan tests, seeded red) | success |
| 5 | D22 scanner self-test (ticket 161 AC#2 - every ban, both directions) | success |
| 6 | D22 seven-ban + emoji scan (tools/d22scan) | success |
| 7 | Tracked path-length budget (ticket 262) | success |
| **8** | **SLO nail census (ticket 266 - the outer tooth over slo-check.ps1's nails)** | **failure** |
| 9 | gofmt (gofumpt) | success |
| 10 | gofmt (gofumpt) - the tracked set is the denominator (ticket 161 AC#7 form A) | success |
| **11** | **go vet (module)** | **failure** |
| 12 | go vet (tools/d22scan module) | success |
| **13** | **staticcheck** | **failure** |
| 14 | mockllm module vet | success |
| 15 | Portable tests carrier self-test (ticket 111 AC#3 - GUARD D's positive control) | success |
| 16 | Attribution instrument self-test with its roster nails (ticket 171 AC#3 - the mutation goes red) | success |

## 3. ★★★staticcheck step（#13）的逐字输出

**版本钉死（票 85 AC#1 的前半 + AC#2 的前提）**：日志第一行逐字
```
##[group]Run go install honnef.co/go/tools/cmd/staticcheck@2026.2.2
go install honnef.co/go/tools/cmd/staticcheck@2026.2.2
```

**真判据产出（票 85 AC#2 的正面证据）**：三个 module 各自一行 + 总结一行，逐字：
```
=== staticcheck module: .
--- module .: exit=1 packages=40 findings=54 toolchain-crash-lines=0
=== staticcheck module: tools/d22scan
--- module tools/d22scan: exit=0 packages=1 findings=0 toolchain-crash-lines=0
=== staticcheck module: tools/mockllm
--- module tools/mockllm: exit=0 packages=1 findings=0 toolchain-crash-lines=0
staticcheck self-report: version=staticcheck 2026.2.2 (0.8.2) modules=3 packages=42 findings=54 toolchain-crash-lines=0 (step exit is 1)
```

## 4. 这些读数对三张票的意义

### 票 85 — AC#1（两把工具钉版本 + 一次真 run 日志行）

| 票面要求 | 现读数 | 判 |
|---|---|---|
| 两把工具都钉成具体版本 | `ci.yml:307 go install mvdan.cc/gofumpt@v0.12.0`；`ci.yml:465 go install ...staticcheck@2026.2.2`（现读源码） | **满足** |
| **一次真 run 的日志行证明它下载的就是那个版本** | `##[group]Run go install ...staticcheck@2026.2.2` ＋ `staticcheck self-report: version=staticcheck 2026.2.2 (0.8.2)` ⇒ **下载的与自报的版本一致** | **满足** |

### 票 85 — AC#2（staticcheck 第一次产出真实判据：按"包 × 检查码 × file:line"列表）

- **第一次产出**：`findings=54`、`toolchain-crash-lines=0` ⇒ **计数器真跑起来了**（此前是"跑不动/零判据"）。
- ⚠ **"按 包×检查码×file:line 列表"那一半需从 `m_out` 全文抽**（日志里 `printf '%s\n' "$m_out"` 后有 54 行 finding）。**本件只证"判据产出了"，具体 54 行的清单要另取**。
- ⚠ **`findings=54` 与票 122 自报的 `linux 20→14` 不一致**——两发跑在**不同提交**上（122 的读数是批 2 后、本 run 是 `a6d203613`）。**这个差必须核**，别当成同一件事。

### 票 122 — AC#6（步级 CI 读数：run id + job id + step 号 + 结论）

| 票面要的四件 | 现读数 |
|---|---|
| **run id** | `43139602748` |
| **job id** | `130923107457`（job 名 `lint`） |
| **step 号** | **`13`**（step 名 `staticcheck`） |
| **结论** | **failure**（因 `findings=54` 存在 ⇒ staticcheck 非零退出，是门禁在**正常工作**，不是工具坏了） |

⚠ **但 122 的 AC#6 要的是"清账到绿"之后的读数**（`lint` 步必须 success）。本 run 是 `a6d203613`，**在 122 的 SA1019 修复（`dc8a2d536`）之前** ⇒ **这张 run 只能证"门在工作"，⛔ 证"已清绿"**。要取 122 AC#6 的终读数必须看**含 `dc8a2d536` 及之后的 run**。

## 5. ★其余 in_progress 的 run（含最新提交，取 122 终读数用）

`gh run list` 现读：
| run id | headSha | 状态 |
|---|---|---|
| `43140706797` | `9e1a8dde5`（最新） | **in_progress** |
| `43140608900` | `b03e20e44` | in_progress |
| `43140335973` | `d77e5b650` | in_progress |
| `43139602748` | `a6d203613` | completed / failure（本件数据源） |

⇒ **`9e1a8dde5` 的 run `43140706797` 跑完后再取一次，才是 122 AC#6 的终读数**。

## 6. 撤销口令

删本件即复原（本件零产码，仅采集 CI 读数）。
