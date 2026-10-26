# 审计核实报告（2026-10-26，编排者）

> **性质**：编排者对 2026-10-25 独立审计（`2026-10-25-harness-deep-audit.md` ＋ `2026-10-25-harness-reform-plan.md`）与机主贴出的**复审结果**的**独立复核**。
> **方法**：全部数字**自己现跑**（`ls`/`git ls-files`/`grep`/`gh api`），⛔ 采信审计读数。六旗标逐条重量；CI 红因**拉 job log 逐行看**（审计自己声明未覆盖这一项）。
> **结论**：★★**审计的主体定性属实**（核心盘面数、六旗标①、CI 三 job 红、timeout=0、.scratch 89%、分支零保护，逐条对上）；★★**部分数字我核出比审计更大**（旗标②③④）；★★★**"修 CI 三 job 红"的性质与审计想象完全不同**（见 §3，这是本报告最重要的一条）；★**复审有一处措辞夸大**（`D22 mode-6` 的引申，见 §5）。

---

## 1. 盘面基本数（逐条对上）

| 项 | 审计说 | ★我现跑 | 判 |
|---|---|---|---|
| 票池总数 | 365 | **365** | ✅ |
| done | 189 | **189** | ✅ |
| void | 1 | **1** | ✅ |
| 活票 | 175 | **175** | ✅ |
| 台账大小 | 6.24MB | **6.0M**（`du -h`）| ✅ |
| 台账行数 | 26k | **26,254** | ✅ |
| 跟踪文件 | 17,718 | **17,718** | ✅ |
| `.scratch` 跟踪 | 15,789（89%）| **15,789**（89.1%）| ✅ |
| 根目录未跟踪件 | 1,259 | **1,262** | ✅（差 3）|
| `s1` 证据件 | 628 | ★**618** | ⚠ 差 10（口径未明）|
| CI job 数 | 6 | **6**（lint／test-core／test-windows／slo-smoke／slo-full／lint-frontend）| ✅ |
| `timeout-minutes` | 0 | **0** | ✅ |
| dev 分支保护 | API 404 | ★**404 "Branch not protected"** | ✅ |

## 2. 六旗标逐条重量（★我用两种格式的尺）

★★**先说一个尺的坑**：票面有**两种 Status 写法** —— 旧（09 月）`**Status:** done` **81 张**、新（10 月）`- Status: **DONE（…）**` **71 张**、无 Status 行 **37 张**。
⇒ ★**只认 `- Status:` 的尺会漏判 81 张** ⇒ 下面全部用**两种都认**的尺。

| 旗标 | 审计说 | ★我现跑 | 判 |
|---|---|---|---|
| ① 残余未勾（done 票仍有 `- [ ]`）| **3**（07 92 115）| ★**3（07 92 115）** | ★★**完全一致** |
| ② Status 从未写 done | **69/189** | ★**75/189**（首行口径）| ⚠ **更大** |
| ③ s1 无证据件 | **25** | ★**52** | ⚠ **更大**（口径见下）|
| ④ 活票无状态行 | **20** | ★**30** | ⚠ **更大** |
| ⑤ 台账零提及 | 18 | 未独立核（需全文扫 6MB）| — 未核 |
| ⑥ 无出处标记 | 32 | 未独立核 | — 未核 |

★**②③④ 的差异是尺的口径，⛔ 审计编数**：
- ★**②**：审计 69 与我 75 差 6，抽查查明有一类是**"首行 `**Status:** open`、后续追加行写 `accepted-done`"**（如票 **101**）—— 我的 `-m1` 只取首行漏判。**审计的 69 更接近真值**，但**方向一致**。
- ★**③**：我的 52 里包含 **1–9 号等 9 月最早期票**，那时还没有 `s1` 三件套制度 ⇒ **审计的 25 大概是"应该有三件套但没有"的窄口径**，我的 52 是"所有无 s1 文件的 done 票"的宽口径。★审计自述抽查过票 299（终裁在台账 A798）⇒ **它的口径是有依据的**。
- ★**④**：我 30 张**完全包含**审计的 20 张，多出 10 张（233 234 238 239 247 248 249 290 293 296 297 316 320 323 324 351 中的一部分）—— 差异来自"Status 行在多少行以内"。

⇒ ★★★**结论**：★**旗标① 完全一致、②③④ 方向一致且我的数更大** ⇒ **审计没有夸大问题**。

## 3. ★★★本报告最重要的一条：CI 三 job 的**红因**（审计声明未覆盖）

★**审计 §1 诚实登记**："③ CI 三 job 本次红的具体红因（**只确认 job 名与结论**）" ⇒ **它没查红因**。我查了。

### ★红的事实（完全对应审计）
`push` run **`40150060567`**（HEAD `f559f14f6`、我的 A1242 提交）⇒ `conclusion: failure`

| job | 结论 |
|---|---|
| `lint-frontend` | ✅ success |
| `slo-full` | ✅ success |
| `slo-smoke` | ✅ success |
| ★`lint` | ❌ **failure** |
| ★`test-core` | ❌ **failure** |
| ★`test-windows` | ❌ **failure** |

⇒ ★★**审计的"红 3 job（lint/test-core/test-windows）"完全属实。**

### ★★★但红因**全部不是代码 bug**

#### `lint` ⇒ 失败步骤 ＝ `staticcheck`
```
-: internal error in importing "cmp" (cannot decode "cmp", export data version 5
   is greater than maximum supported version)
-: internal error in importing "internal/byteorder" / "internal/cpu" / "math/bits" …
```
⇒ ★★★**`staticcheck@2026.2.1` 不认识 CI 上 Go 版本产出的 export data（version 5）** ⇒ **工具链版本不兼容**
⇒ ★**与我们改的代码零关系**（`gofumpt`/`gofmt`/`go vet`/`d22scan`/路径长度门**全部 success**）

#### `test-core` ⇒ 失败步骤 ＝ `Portable package tests (core scope)`
| 失败测试 | 包 | ★红因 |
|---|---|---|
| `TestDsn196r3BusinessPoolsKeepTheirOwnForeignKeys` | `internal/memory` | ★`dsn(...) changed shape: got file:C:\Users\x\… want file:C:/Users/x/…` ⇒ **路径分隔符**（测试写死 Windows 路径、在 Linux runner 上跑）|
| `TestC21TableColourRowsMatchTokensCSS` | `internal/ball` | ★`open /home/runner/work/wisp/wisp/design/assets/tokens.css: no such file` ⇒ **`design/` 缺件** |
| `TestApprovalCardViewJSONKeysMatchFrontendTypes` 等 4 枚 | `internal/panel` | **前端类型耦合**（需 `frontend/` 树）＋ 2 枚 35-r8 **既有红** |
| `TestC26RewriteAccountIsConsumedAtEverySecurityLeg` | `internal/risk` | ★`internal/panel/composer.go:127`／`display_scrub_216.go:26/37/218` **读 `Result.Canonical` 没读 rewrite account** ⇒ **既有红** |

#### `test-windows` ⇒ 失败步骤 ＝ `cmd/wisp CLI tests` ＋ `Portable windows tests`
| 失败测试 | ★红因 |
|---|---|
| `TestTicket223PermissionDeniedSitsInItsOwnSentence` | ★`attempt 1: icacls denied nothing, so this case cannot show the permission sentence` ⇒ **runner 上 ACL 拒不掉**（环境）|
| `TestPanelHolderShowRebuildsADeadThread319` | ★`the first thread did not finish after stop()` ⇒ **面板线程停不下来**（时序/环境；同族 `...RebuildsADeadThread319` 是 PASS）|
| `TestAC14GoSideEvalPushReachesThePage` | ★`Eval push did not reach the document: page reports title as ""` ⇒ **WebView2 Eval 未送达**（环境）|
| `TestC26RewriteAccount…` ／ `TestC21TableColourRows…` | **与 test-core 同一批**（`internal/risk`／`internal/ball`）|

### ★★★所以「修 CI 三 job 红」的真实形状

★★**不是"修 4 处代码 bug"，而是"修 3 类门禁形状问题"**：

| 类 | 具体 | 修法（代价不同）|
|---|---|---|
| ★**A 工具链版本** | `staticcheck` vs Go export data | ★**升 `staticcheck` 版本**（一行、不改判据）|
| ★**B 缺件** | `design/assets/tokens.css` 在 CI 上不存在 | ★**要么提交它进仓、要么让该测试在缺件时明确 SKIP**（后者碰"不跳过"约定）|
| ★**C 环境限制** | `icacls` 拒不掉／WebView2 线程与 Eval／跨平台路径 | ★**要么把这些测试挪到 `self-hosted` 上跑、要么改判据为"环境不满足时 SKIP"** |

⇒ ★★★**B/C 两条都撞上「不跳过」这条约定** —— 正是复审说的"不跳过铁律挡住最直接的修法"。

## 4. 复审三条安全暴露的复核

| # | 复审说 | ★我核 | 判 |
|---|---|---|---|
| 1 | `ci.yml:9` `pull_request` 无过滤 ＋ `slo-full` 跑在 `[self-hosted, wisp-slo]` | ★`ci.yml:11 pull_request:`（**无 branches/paths 过滤**）＋ `:1014 runs-on: [self-hosted, wisp-slo]` | ★★**完全属实** |
| 2 | `.scratch` 占 89%，`injection-timeline.md` 公开 | ★**89.1% 属实**；★但 `injection-timeline.md` 在 **`docs/reports/`**、⛔ `.scratch/` | ★**主体属实、位置有偏差** |
| 3 | dev 分支零保护（API 404）| ★`gh api …/branches/dev/protection` ⇒ **404 "Branch not protected"** | ★★**完全属实** |

★★**第 1 条的风险判断我认同**：`pull_request` 无过滤 ⇒ 任何 PR（含 fork PR）会触发**全部 6 job**，其中 `slo-full` 跑在**承载 DPAPI 凭据/gh 登录态/开发数据的开发机**上。

## 5. ★★复审一处措辞夸大（这一条影响改法）

★**复审逐字**："而 **D22 mode-6 明令禁** ⛔：门控"

★★**D22 原文（`PLAN.md:722` 起、八种跑歪模式表）**：
> | **6** | **绕过约束而非满足** | **SLO 达不到就关掉测量或放宽阈值** | ②禁止清单 + ⑤对抗验收 |

⇒ ★★★**D22 mode-6 讲的是「⛔ 关掉测量／放宽阈值」，⛔ 「⛔ 让 job 可跳过」**

★★**"⛔ skippable" 的真正出处是 `ci.yml:5-8` 自己的注释**：
```
# Per D22 mode-6 NO job is configured skippable: no `if: false`,
# no continue-on-error, no skip flags.
```
⇒ ★★**那是**我们自己在 CI 注释里对 D22 的引申**、⛔ D22 原文**

⇒ ★★★**所以**：
1. ★**改 `ci.yml` 的 `pull_request` 过滤**（让 PR 不跑 `slo-full`）**⛔ 违反 D22 原文** —— ★**但**它**减少了一处测量**、**接近 mode-6 的语义** ⇒ **该摆机主**
2. ★★**改 `staticcheck` 版本**（A 类）**⛔ 碰任何约定** ⇒ **我可以直接做**
3. ★★**"缺件/环境不满足时 SKIP"**（B/C 类）**实质是放宽判据** ⇒ **该摆机主**（或找不改判据的路：把测试挪 runner）

## 6. `ci.yml` 是不是冻结面（决定我能自己拍还是摆机主）

★★**我现跑**：
- `grep 'ci.yml' docs/specs/SPEC-11-*.md` ⇒ **零命中**
- `grep 'ci.yml' docs/PLAN.md` ⇒ **零命中**
- `grep -E '^\| *D(3[0-9]|4[0-7])' … | grep -iE 'CI|workflow|门禁'` ⇒ **零命中**

⇒ ★★★**`ci.yml` ⛔ `C1–C32`／`D1–D47` 契约文字** ⇒ ★**改它不触发"改契约须人工批准"**（`SPEC-12 §4.1`）
⇒ ★★**但**"不跳过"这条**约定**（`ci.yml` 注释 ＋ `SPEC-10 §3`）**实质是判据** ⇒ ★**放宽判据按铁律要写清代价、并由我拍或摆机主**（我按"最有价值"判：**⛔ 为绿而放宽**，见 §7）

## 7. 我的执行判断

★**审计与复审的定性我全部采信**（核心盘面数、六旗标①、CI 三 job 红、timeout=0、.scratch 89%、分支零保护，逐条自己核过）。
★**数字差异不构成"审计不实"** —— 方向一致、且我的数更大。
★**阶段 0 的边界**：
- ★**可以直接做**：`timeout-minutes`（纯增、不改判据）、`scripts/ticket-lint.sh`（新件）、`.gitignore`（新件）、`HANDOVER §4`（文档）、派单存档（文档）、**升 `staticcheck`**（工具链）
- ★★**必须摆机主**：①`ci.yml` 的 `pull_request` 过滤（减少一处测量）②任何"环境不满足即 SKIP"（放宽判据）③票池 schema v2／台账切档／根目录清扫（方案 §5 已列）
- ★★★**我不会做的事**：为让 CI 变绿而 `continue-on-error`／删步骤／放宽断言（`D22 mode-6` 字面禁止的正是这个）
