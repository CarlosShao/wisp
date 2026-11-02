# 票 337 非实现者验证腿 `337-v1` · 验收裁决表

- 腿：`337-v1`＝**非实现者验证腿**（裁决者≠实现者；两段落笔人＝`337-r1`／`337-r2`，本腿零产码）。
- 派单存档：`.scratch/wisp/probes/orch/20261102/txt/337v1b-dispatch.txt`（重跑单，config 已收口）。
- 票面：`.scratch/wisp/issues/337-fail-closed-reason-has-no-reader-wiring-census.md`（blob＝`88ae0d0052b3dabda84d13e0d388bbcdd21071da`，**本腿未动一字**）。
- 被验交付：`66d528120`（337-r1＝`internal/agent/approval/gate.go` 的 `New()` 接线）＋ `7fd2f37b5`（337-r2＝`cmd/wisp/run.go` 注入式宿主装配点补接）。
- 边界：⛔ 改任何代码／⛔ 改票面／⛔ 台账／⛔ push；**`cmd/wisp` 整包与单发归 `395-v1`**（凡必须跑 `cmd/wisp` 的格，本腿**具名跳过**并留可复跑尺）。
- 读数一律取**具名提交树**（`HEAD`）或**仓外 `-overlay`**；尺一律带射程并**剔 `.scratch/**`**。临时件＝`D:\tmp\wisp\20261102\337v1\`（只建不删）。

---

## 0. 起手三核

| 核 | 尺（逐字可重跑） | 读数 | 判 |
|---|---|---|---|
| ① HEAD ≥ `ea41fcab7` | `git rev-parse HEAD`；`git merge-base --is-ancestor ea41fcab7 HEAD` | `ea41fcab756d7fa8d6ba7719e67bdf278bf9f723`；`rc=0`（相等）；分支 `dev` | **PASS** |
| ② 索引干净 | `git status --porcelain --untracked-files=no` 行数 | **0 行**（跟踪面零脏；余下全是 `.scratch/**` 等别腿未跟踪件） | **PASS** |
| ③ 两枚交付真在 HEAD 里 | `git merge-base --is-ancestor 66d528120 HEAD`／`… 7fd2f37b5 HEAD` | `rc=0`／`rc=0` | **PASS** |

台面：起手与收笔两次 `git status --porcelain --untracked-files=no` 均＝**0**；本腿只新增本件一枚（见 §6）。

---

## 1. AC#0 — `BlacklistNote.Reason` 消费者＝0（正控必须有牙）

`BlacklistNote` 定义在 `internal/tools/mode.go:66`，`Reason` 字段在 `:76-77`，唯一写点 `:111`（`n.Reason = d.Reason`）。

| 尺 | 命令（逐字，`HEAD` 树，剔 `.scratch`） | 读数 | 判 |
|---|---|---|---|
| R1 显式字段读 | `git grep -nE 'Blacklist\.Reason' HEAD -- '*.go' \| grep -v ':\.scratch/'` | **0 枚** | — |
| R1b 类型点形 | `git grep -nE 'BlacklistNote\.Reason' HEAD -- '*.go' \| grep -v ':\.scratch/'` | **0 枚** | — |
| R2 变量读形（前缀锚＋行尾支） | `git grep -nE '(^\|[^A-Za-z])n\.Reason($\|[^A-Za-z])' HEAD -- '*.go' \| grep -v ':\.scratch/' \| grep -v '_test\.go:'` | **1 枚＝`internal/tools/mode.go:111 n.Reason = d.Reason`＝写侧 ⇒ 读＝0** | — |
| **正控**（同尺换一句真有人读的理由串） | `git grep -nE '(dec\|decision\|it\.Decision)\.Reason($\|[^A-Za-z])' HEAD -- '*.go'`（剔 `.scratch`、剔 `_test.go`） | **8 枚**＝`cmd/wisp/panel_pump.go:75`、`internal/panel/approval.go:83/84`、`internal/tools/bridge.go:317/468/1122/1123/1205` | 尺认得真读者 ⇒ 上面那个 0 是**真 0** |
| 反证「尺坏」（刻意用宽松无锚形） | `git grep -nE '\.Reason' HEAD -- '*.go'` 全量 | **386 枚**（不收紧就会放行） | 证明收紧是必要的、不是摆设 |

**接线后新增读数（＝本票交付的落点）**：`git grep -nE 'risk\.Logf[[:space:]]*=' HEAD -- '*.go'`（剔 `.scratch`、剔 `_test.go`）＝**恰 2 枚**＝`internal/agent/approval/gate.go:181`＋`cmd/wisp/run.go:758`（接线前普查＝**0**）。

> 具名说明：`BlacklistNote.Reason` 的读者**今天仍是 0**——这不是回归缺陷，而是甲形的定义本身。票 334 那句实话的唯一携带出口是 `risk.Logf` 的 `rule=` 字段（`internal/risk/blacklist.go:83`）；接线把**那句实话**接进本机审计日志，⛔ 给那个结构体字段新增读者。票面 `:13` 的前提、以及票面 AC#2 那枚硬门（未经 `A##` 批准⛔ 接进模型正文／面板），都因此继续成立。

**判定：AC#0 ＝ PASS。** 逐字读数存 `D:\tmp\wisp\20261102\337v1\ac0-rulers.log`。

---

## 2. AC#1 — 两处接线是否真在、且真在声称的位置

| 格 | 位置（`HEAD` 现量） | 逐字 | 判 |
|---|---|---|---|
| 接线① | `internal/agent/approval/gate.go` `New()` 体内 `:180-182`（说明注释 `:152-179`） | `if o.Logf != nil {` / `risk.Logf = o.Logf` / `}` | **PASS** |
| 接线② | `cmd/wisp/run.go` `assembleRuntime` 的**注入式宿主分支**（分支自 `rt.gate = s.gate` `:725` 起）末段 `:758`（说明注释 `:728-757`），紧邻 `} else {`（console 支 `:759`）之前 | `risk.Logf = rt.auditf` | **PASS** |
| 喂线方（证明①不是死码） | `git grep -n "Logf:" HEAD -- cmd/wisp/run.go`＋`git grep -n "Logf" HEAD -- cmd/wisp/resident_approval_windows.go` | `run.go:551`／`:767`／`:812`／`:941` 均＝`rt.auditf`；`resident_approval_windows.go:404`＝`ra.residentAuditf` ⇒ 控制台门与常驻门都经 `approval.New` 把自家审计 sink 递进 ① | **PASS** |
| 两支同源 | ②赋的值＝console 支递进 `Options.Logf` 的同一个值（`run.go:767`） | 逐字同＝`rt.auditf` ⇒「一条进程一个 sink 族」成立 | **PASS** |

可重跑尺：`git show HEAD:internal/agent/approval/gate.go`（取 175–185 行）；`git show HEAD:cmd/wisp/run.go`（取 755–770 行）；`git grep -nE 'risk\.Logf[[:space:]]*=' HEAD -- '*.go'`。

**判定：AC#1 ＝ PASS。**

---

## 3. AC#2 — 反向判据（仓外 `-overlay` 摘掉这两处赋值 ⇒ 必须红）＋还原核

**台面**：突变件全在仓外（`D:\tmp\wisp\20261102\337v1\overlay\`），⛔ 未向仓内写一字节。
**突变**：`overlay.json` 把两处**同时**指回「无赋值」副本——
- `internal/agent/approval/gate.go` 摘掉 `risk.Logf = o.Logf`（874→873 行；摘后该文件 `risk.Logf` 赋值命中＝0）；
- `cmd/wisp/run.go` 摘掉 `risk.Logf = rt.auditf`（1927→1926 行）。

| 格 | 尺（逐字） | 读数 | 判 |
|---|---|---|---|
| 绿基线（无 overlay） | `go test ./internal/agent/approval/ -run TestTicket337FailClosedReasonReachesAuditLog -count=1 -v` | `rc=0`、`--- PASS (0.01s)`，且仪器**真的**打出落 sink 的实话行：`risk: A-list DENY path=… rule=~/.aws/credentials [fail-closed, NOT a match on this rule]: …` | **PASS** |
| 突变红（两处同摘） | `go test -overlay D:\tmp\wisp\20261102\337v1\overlay.json ./internal/agent/approval/ -run TestTicket337FailClosedReasonReachesAuditLog -count=1 -v` | **`rc=1`**，红句逐字＝`ticket337_reason_reader_test.go:90: approval.New did not arm risk.Logf: the fail-closed reason has no reader again (ticket 337)` | **PASS**（尺有牙，不是恒绿） |
| 还原核 | `git status --porcelain --untracked-files=no`＝**0**；`git hash-object` 工作树 vs `HEAD` blob＝`43ed3212739c84661c7c5ff3d3a360dbee984040`（gate.go）／`82e052e03606a6a0ea09dfe5e8a108633c2ed078`（run.go）**逐字节同串**；重跑同测试（不带 overlay）＝`rc=0` PASS | 还原成立 | **PASS** |
| **具名跳过** | **注入式宿主（`run.go:758`）「不落 sink ⇒ 红」的经验跑**：唯一仪器＝`cmd/wisp/risk_sink_injected_337r2_windows_test.go`（`package main`）两枚 `TestTicket337R2InjectedGateArmsRiskSinkOnRuntimeAudit`／`…WhenHostNeverArmed` ⇒ 跑它必跑 `cmd/wisp`。**按派单边界（`cmd/wisp` 单发留给 `395-v1`）本腿不跑。** | 未跑 | **SKIP（具名）** |
| 该格本腿的替代证据（只读） | ① 摘这 2 枚之后生产码**不存在任何其它 `risk.Logf` 注入点**（`git grep -nE 'risk\.Logf[[:space:]]*=' HEAD -- '*.go'` 剔 `_test.go` 恰＝这 2 枚；接线前普查＝0）⇒ 注入支上 `risk.Logf` 必回 nil；② 读该两枚用例断言原文：`:174`／`:243` 起手即 `if risk.Logf == nil { t.Fatal(…) }`，`:183-185` 断言 stderr 含 `[audit] risk: A-list DENY`，`:187-189` 断言 foreign sink 收行数必须＝0 ⇒ 摘线后至少两枚断言必红。 | 静态闭合，形状与 approval 支已经验证的红同 | 记为**未经验证**，⛔ 当绿 |

**判定：AC#2 ＝ PASS**（`cmd/wisp` 那一半具名跳过、交 `395-v1` 复跑；复现尺＝上面「突变红」那行，把包名换成 `./cmd/wisp/`、`-run` 换成 `'TestTicket337R2InjectedGateArmsRiskSink'` 即可，overlay 沿用本腿同一件）。

---

## 4. AC#3 — 零回归（四包逐包四数）

尺：`go test ./<pkg>/ -count=1 -v`（`HEAD`＝`ea41fcab7`，跟踪面干净）；日志 `D:\tmp\wisp\20261102\337v1\ac3-pkgs.log`。`RUN`＝`PASS+FAIL+SKIP`。

| 包 | rc | RUN | PASS | **FAIL** | SKIP | 判 |
|---|---|---|---|---|---|---|
| `internal/agent/approval`（两段落点之一） | 0 | 86 | 85 | **0** | 1（`TestDefaultDeadlineWallClockMeasurement`，存量平台 skip） | **PASS** |
| `internal/risk`（sink 生产方） | 0 | 172 | 171 | **0** | 1（`TestSyncRegistryProbeLive`，存量平台 skip） | **PASS** |
| `internal/tools`（`readBlacklist` 调用方） | 0 | 253 | 253 | **0** | 0 | **PASS** |
| `internal/panel`（甲形不碰面，负控） | 1 | 173 | 172 | **1** | 0 | 见下方归因 |

**`internal/panel` 那枚 FAIL 的归因尺（现量，不是「我以为」）**：
- 红的名字＝**`TestC21DesignTokensFourWayAgree`**；红条逐字＝`docs/evidence/s1/c21-native-tokens.md:58/59/60/103/104/105: the table names Palette.GlassRing / GlassHi / GlassHiSoft, but internal/ball/tokens.go holds no such colour literal in DarkPalette()/LightPalette()` ＋ `check 2: frontend/src/styles/tokens.generated.css declares 34 coloured tokens = 11 mapped + 23 waived`。
- 造红那一笔＝**`ebf33fd77`（票 344 落地腿 344-r1：「C21 真相源改指 v3 ＋删三枚死玻璃色码」）**；`git show --name-only ebf33fd77` 正是 `internal/ball/tokens.go` ＋ `docs/evidence/s1/c21-native-tokens.md`。
- 该笔**不是** 337 两枚交付的祖先：`git merge-base --is-ancestor ebf33fd77 66d528120`＝**rc=1**、`… 7fd2f37b5`＝**rc=1**（时间戳 337-r1 `09:29`／337-r2 `11:24`／344-r1 `11:39`）。
- 337 两枚交付对 `internal/ball/**`、`internal/panel/**`、`c21-native-tokens.md`＝**零字节**（`git show --name-only` 逐枚对 `ball\|panel\|c21\|frontend\|design\|golden\|thresholds\|PLAN.md\|issues/`＝0 命中）。
- 票面 **AC#4 早已具名**这一族红：「`internal/panel` 那两枚前端契约键族红与 `TestC21…` 属**别面存量**、⛔ 计入本票欠账」。

**判定：AC#3 ＝ PASS**（approval／risk／tools 三包 FAIL＝0；`internal/panel` 的 1 枚具名归因到 `ebf33fd77`／票 344 面，按票面 AC#4 不计本票欠账，本腿亦⛔ 代 `344-v1` 裁它）。

---

## 5. AC#4 — 门禁四数（全部只在共享工作树现跑）

| 门 | 尺（逐字） | rc | 读数 | 判 |
|---|---|---|---|---|
| G1 d22scan | `sh scripts/d22scan.sh`（仓根） | **0** | `d22scan: clean - no D22 ban violations`；live scope：bans#1-5 `internal/`=288、`cmd/`=52；ban#6 `frontend/`=116；ban#7 `internal/tools/`=25；ban#8 `design/`=39、`frontend/`=116、`internal/`=737、`cmd/`=186；正控 `runtests.sh -C tools/d22scan ./...`＝`PASS=35 FAIL=0 SKIP=0` | **PASS** |
| G2 pathlen | `sh scripts/check-path-length-budget.sh --with-self-test` | **0** | `VERDICT GREEN`；tracked paths＝18321、over-budget＝**62**、roster＝**62**、not in roster＝**0**；三枚 positive control 全 ok；字符／字节两单位交叉核对同选 62 枚 | **PASS** |
| G3 vet | `go vet ./internal/agent/approval/ ./internal/risk/ ./internal/tools/ ./internal/panel/` | **0** | 输出**零行** | **PASS** |
| G4 gofumpt | `D:/work/base/gopath/bin/gofumpt.exe -l internal cmd`（CI 同 scope，绝对路径） | **0** | 命中 2 枚＝`internal/agent/approval/pending_read.go`、`internal/agent/tools.go`——**均非本票交付件**；`git ls-files --eol`＝`i/lf w/crlf attr/text eol=lf`（工作树 CR 字节 131／225＝**检出 CRLF 伪影**）；把这俩的**跟踪 blob**（`git show HEAD:…` 落仓外）再跑 `gofumpt -l`＝**空**；两枚最后落笔人＝`d03d166fe`／`5413f46d4`（票 212 族），与 337 无关 | **PASS** |
| G4b 交付件逐枚 | 四枚跟踪 blob 落仓外逐枚 `gofumpt -l` | **0**×4 | `gate.go`、`ticket337_reason_reader_test.go`、`run.go`、`risk_sink_injected_337r2_windows_test.go` 全部**空输出** | **PASS** |
| 具名跳过 | **`./cmd/wisp` 整包测试**（票面 AC#4＋派单：`cmd/wisp` 单发归 `395-v1`） | 未跑 | 按边界具名跳过，⛔ 充当绿 | **SKIP（具名）** |

---

## 6. AC#5 — 纪律格

| 核 | 尺 | 读数 | 判 |
|---|---|---|---|
| 两笔与题面一致 | `git show --stat 66d528120`／`git show --stat 7fd2f37b5` | r1＝`3 files changed, 305 insertions(+)`＝`internal/agent/approval/gate.go +31`＋`internal/agent/approval/ticket337_reason_reader_test.go +133`＋证据件 `+141`（＝派单所述「337-r1 接线 `gate.go New()`」）；r2＝`3 files changed, 391 insertions(+)`＝`cmd/wisp/run.go +31`（产码仅 1 行 `:758`，余为说明注释）＋`cmd/wisp/risk_sink_injected_337r2_windows_test.go +253`＋证据件 `+107`（＝派单所述「337-r2 补注入式宿主 `risk.Logf = rt.auditf`」） | **PASS** |
| 写面未越界 | `git show --name-only` 两笔逐枚 | 命中面仅＝`internal/agent/approval/**`、`cmd/wisp/**`、`docs/evidence/s1/337r*`；对禁区词 `ball\|panel\|c21\|frontend\|design\|golden\|thresholds\|PLAN.md\|issues/`＝**0 命中** | **PASS** |
| 票面未动 | `git rev-parse 66d528120~1:<票面>` 与 `git rev-parse HEAD:<票面>` | 同一 blob `88ae0d0052b3dabda84d13e0d388bbcdd21071da`；框数＝`- [ ]` **7**／`- [x]` **0**（与交付前锚点逐枚同，**一枚没翻**）；票名／`Status:`／`-done` 未改 | **PASS** |
| D43 转移表未动 | `git log --oneline -3 -- docs/PLAN.md`＋两笔 `--name-only` | `docs/PLAN.md` 最新一笔＝`962516c48`（早于 337 两笔）；337 两笔对 `docs/PLAN.md`＝**0 字节**；`STATE_MACHINE`／`issues/README`＝0 命中 | **PASS** |
| 台账未动 | `git log --oneline -3 -- docs/reports/pending-and-issues.md` | 最新＝`ea41fcab7`（编排者派单笔），本腿零触碰；两笔对 `go.mod`／`go.sum`／`.gitignore`／`ci.yml`＝0 命中 | **PASS** |
| 本腿自身纪律 | 只 commit 不 push；`git add`／`commit` 显式 pathspec＝本件一枚；⛔ `add -A`／⛔ `--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`；临时件只建不删；只读复核＝跟踪面起手 0／收笔 0 | 成立 | **PASS** |

---

## 7. 逐格总表

| 格 | 判 | 一句话尺 |
|---|---|---|
| 起手三核 | **PASS** | HEAD＝`ea41fcab7`、跟踪面 0 行、两笔均 `--is-ancestor rc=0` |
| AC#0 消费者 0 ＋正控 | **PASS** | R1/R1b＝0、R2＝1（写侧⇒读 0）、正控 8 枚、宽松形 386 枚证明尺有牙 |
| AC#1 两处接线 | **PASS** | `gate.go:180-182` ＋ `run.go:758` 逐字在位；5 处 `Logf:` 生产喂线证明非死码 |
| AC#2 反向判据 | **PASS（内含一格具名跳过）** | overlay 同摘两处 ⇒ `rc=1`，红在 `:90`；还原＝跟踪面 0 ＋ blob 同串 ＋ 重跑 PASS；注入支经验红属 `cmd/wisp`，归 `395-v1` |
| AC#3 零回归 | **PASS** | approval/risk/tools FAIL＝0；panel 的 1 枚＝`TestC21DesignTokensFourWayAgree`，归因 `ebf33fd77`（344-r1，非 337 祖先，票面 AC#4 具名别面存量） |
| AC#4 门禁四数 | **PASS** | d22scan rc=0 clean／pathlen rc=0 GREEN／vet rc=0／gofumpt 跟踪集 rc=0；`./cmd/wisp` 整包具名跳过 |
| AC#5 纪律 | **PASS** | 两笔＝`+305`／`+391` 与题面同；票面 blob 未变、7/0 框未翻；D43 与台账 0 字节 |

**残余（交编排者，⛔ 本腿代裁）**
1. `cmd/wisp` 整包 ＋ 注入支经验红 ＝ `395-v1`（尺见 §3 末行）。
2. `gate.go:175-179` 那段 **RESIDUAL 文字今天已过期**（注入支已由 `run.go:758` 自理）——337-r2 证据件自己具名上报，归一句票面／注释更正；本腿⛔ 动产码注释。
3. 六族／`262-640-902` 桌面延迟族：本验证腿无真桌面会话，⛔ 充当绿，留编排者桌面轮复跑。
4. 乙形（面板出向）／丙形（模型正文）**未落地，本腿也不裁**：票面 AC#2 硬门要求台账 `A##` 具名批准（文件／行／理由／边界／撤销口令），盘上不存在该批准；两段落点均为甲（本机审计日志），未越门。

---

## 8. 是否同意结 337

**同意结票 337。**

理由：票面要买的东西已经拿到，并且**有牙**——那句 fail-closed 实话今天在生产路径上落进了本机审计 sink，两条装配（经 `approval.New` 的控制台门与常驻门、以及注入式宿主那条 `s.gate != nil` 支）都接上同一个 `rt.auditf`；用仓外 overlay 摘掉这两处赋值，仪器当场红（`rc=1`，红句逐字点名 “the fail-closed reason has no reader again”）。零回归、门禁四数、纪律格全绿；唯一那枚 `internal/panel` 的红可复跑地归因到 `ebf33fd77`（票 344 面），且票面 AC#4 早已把它具名排除在本票欠账之外。

**附一条边界（不是保留意见，是范围事实）**：本同意以「`cmd/wisp` 那两枚注入支用例的经验红由 `395-v1` 复跑并收件」为前置条件——该格本腿按派单具名跳过，静态已闭合但未经验证。`395-v1` 若复跑出红，本票应就地重开，⛔ 以本件为结票终据。

*验证腿：`337-v1`（非实现者）。本腿产码字节＝0；只 commit 本件，⛔ push。*
