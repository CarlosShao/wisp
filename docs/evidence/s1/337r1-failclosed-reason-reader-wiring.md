# 票 337 落地腿 `337-r1` · 证据件 — `BlacklistNote.Reason` 的 fail-closed 实话接上读者（甲形·本机审计日志）

- 腿：`337-r1`（写腿／实现者）。类型：接线（先只读普查，再决定接哪一层并接线）。
- 来源＝非实现者终裁腿 `334-v1` 必答①（判语件 `docs/evidence/s1/334-v1-verdict.md`）。
- 派单存档：`.scratch/wisp/probes/orch/20261102/txt/337-r1-dispatch.txt`。
- 主题：`337-r1`。

> 所有行号只当**读数**，每条附可重跑的尺。尺一律取**具名提交树**（`git grep <pat> HEAD`），并**剔 `.scratch/**` 影子件**（旧腿在 `.scratch` 里存了大量 `bridge.go`/坏样本快照，不剔＝假命中堆）；正控见 §1。

---

## 0. 起手三核

| 核 | 读数 | 判 |
|---|---|---|
| ① HEAD 在 dev ≥ `f0e42dace` | 起手 `git rev-parse HEAD`＝`07a6ca7b4f9e0d701f7818f2f19e1da9fb50ee26`（dev；`f0e42dace` 为其祖先，`git rev-list --left-right --count HEAD...origin/dev`＝`0 0`） | ✅ |
| ② 索引干净 | `git status --porcelain` 的**跟踪件**面＝0 行（余下全是 `.scratch/**`、`design/**`、`logs/**` 等别腿/环境的未跟踪件，非本腿产出） | ✅ |
| ③ 写面无他人脏 | `git status --porcelain -- internal/risk internal/agent internal/panel docs/evidence/s1` 起手＝空 | ✅ |

⚠ 共享工作树：本腿开工期间 HEAD 由别腿从 `f0e42dace` 推进到 `07a6ca7b4`（另腿的 `.scratch` 票面提交）。本腿读数一律现取 `HEAD`，提交时再取一次（见 §6）。

---

## 1. AC#0 只读普查：`BlacklistNote.Reason` 消费者＝0 枚（尺可重跑＋自带正控）

`BlacklistNote` 定义在 `internal/tools/mode.go:66`；`Reason` 字段在 `:76-77`；唯一写点 `internal/tools/mode.go:111`（`n.Reason = d.Reason`，`n` 在 `:87` 声明为 `BlacklistNote`）。

尺（`HEAD`＝`07a6ca7b4`，剔 `.scratch/**`；逐字结果存 `D:\tmp\wisp\20261102\337r1\census-rulers.log`）：

| 尺 | 命令（`HOME`＝仓根） | 读数 |
|---|---|---|
| **R1 显式字段读** | `git grep -nE 'Blacklist\.Reason' HEAD -- '*.go' \| grep -v ':\.scratch/'` | **0 枚** |
| **R2 笔记变量读形** | `git grep -nE '(^\|[^A-Za-z])n\.Reason($\|[^A-Za-z])' HEAD -- '*.go' \| grep -v ':\.scratch/' \| grep -v '_test.go:'` | **1 枚＝`internal/tools/mode.go:111  n.Reason = d.Reason`（★写侧，非读）⇒ 读＝0** |
| R3 类型级引用 | `git grep -nE 'BlacklistNote' HEAD -- '*.go'`（剔 `.scratch`） | 定义 `mode.go:66`／字段声明 `gate.go:54`／`readBlacklist` 签名 `mode.go:86`／唯一测试 `ticket146..._test.go:82`／注释。**`internal/panel` 对 `Blacklist` 非测试 0 命中**（panel 只吃 `risk.Decision`，不带 `Blacklist`）。 |
| **正控（同尺有牙）** | `git grep -nE '(dec\|decision\|it\.Decision)\.Reason($\|[^A-Za-z])' HEAD -- '*.go'`（剔 `.scratch`、剔 `_test.go`） | **8 枚**：`cmd/wisp/panel_pump.go:75`、`internal/panel/approval.go:83/84`、`internal/tools/bridge.go:317/468/1122/1123/1205` ⇒ 尺认得真读者，上面的 0 是真 0。 |

**旁证（接线落点相关）**：

| 事实 | 读数 |
|---|---|
| `risk.Gate(` 生产调用点 | **恰 1 枚**＝`internal/tools/mode.go:98`（`d := risk.Gate(c, overrides)`） |
| `risk.Logf` 生产注入点 | **0 枚**（`git grep -nE '(^\|[^A-Za-z])Logf\s*=' HEAD -- 'internal/risk/*.go'` 剔 `_test.go`＝0；9 处赋值全在 `_test.go`） |
| `risk.Logf` 是唯一携带 fail-closed 实话的出口 | `internal/risk/blacklist.go:83/87/90`：`logf("risk: A-list DENY path=%s rule=%s ...", p, rule)`，其中 `rule`＝`classifyForms` 返回的 `aRuleReason(...)`＝**正是** `PathDecision.Reason`（`BlacklistNote.Reason`）去掉 `"A-list sensitive path (non-overridable): "` 前缀后的那句实话。sink 为 nil ⇒ 整条丢。 |

**结论**：`BlacklistNote.Reason` 判定型／打印型／出向读者＝**0 枚**；那句 fail-closed 实话只活在返回值 + 一枚**生产未注入**的 `risk.Logf` 里 ⇒ 票 334 买到的"说实话"没有读者，票面前提**成立**（未触发"未定义即停①：消费者其实 >0"）。

---

## 2. 层决策：接 **甲（本机审计日志）** —— 按主流 harness 自己拍，理由如下

> 机主铁律（AGENTS.md ⓪）要求：技术选型 agent 自己拍、写理由、⛔ 丢回机主；判据＝**最有价值＋最符合主流 harness**，⛔ 用"风险最小／动作最小"当判据。本腿据此拍板，未把选择题丢回机主。

三形的现量代价（承 `337-a1` 普查，本腿复核未变）：

- **甲＝本机审计日志（`risk.Logf` 补注入点）**：出向面**零**变化（日志非出向）；不碰 `frontend/**`；无新导出名；绝对/折叠路径**不进模型**。落点＝给那枚今天为 nil 的 sink 补一个生产注入。**价值**：一接线即让整条 `risk` 安全日志（A DENY／B deny／B override）第一次在生产里落盘，其中 `rule=` 就是票 334 的实话。**主流 harness 归位**：诊断级、带折叠真路径的拒绝理由属**取证/审计**，落本地日志正是主流做法；面向人/模型的是**清洗后**的短句（面板/模型正文已有 R3 短句）。
- **乙＝面板拒绝文案**：面板字节面今天**完全不带** `BlacklistNote`（panel 只吃 `risk.Decision`）；复用现有键＝改的是"卡片 reason 这句话本身"，新键＝AC#5 禁区。且 `display_scrub_216.go:58` 自陈**卡片 Reason 不进 216 清洗名册** ⇒ 把带真路径的实话塞进卡片＝**无清洗出向**。
- **丙＝模型正文**：落点在 `internal/tools/bridge.go:468`（`'该目标被绝对禁止访问（R3 A 档…）：' + dec.Reason`）——**不是"新增一句"，是改/加长既有出向句**，会把绝对路径＋控制字符带进模型上下文（票 174／264／216 一族先例）。

**为什么最终是甲（三条硬理由，均非"图省事"）**：
1. **AC#2 硬门（编排者设）**：未经台账具名批准（`A##` 记文件/行/理由/边界/撤销口令），本票任何腿**⛔**把这句接进模型正文或面板出向内容；默认落点＝甲。乙/丙属"让机主或模型看到"的**功能级改动**，须编排者摆机主拍板后才能动。盘上**没有**这样的 `A##` 批准 ⇒ 乙/丙此刻不可落地。本腿**不绕过**这枚硬门。
2. **写面**：本腿写面＝`internal/risk/**`＋`internal/agent/**`＋`internal/panel/**` 的**读侧接线**。丙的落点 `internal/tools/bridge.go` **不在写面内**；乙的落点 `internal/panel` 出向改动＝越出"读侧"且撞 AC#5。唯一同时满足"在写面内 + 甲 + 能给 A 档 fail-closed 实话找读者"的落点是 `internal/agent/approval`（本仓唯一能看见 `tools.Decision.Blacklist.Reason`/其源头 `risk.Gate` 且带生产审计 sink 的包）。
3. **最有价值且不制造新泄漏**：甲把实话接进**审计面**（正是票 334 说的"也不进日志"那一面），零新出向、零新导出、零模型上下文污染；乙/丙都要求把**无 216 清洗的真路径句**推向人/模型，是净负价值。故甲是"最有价值且合规"的层，不是"最省事"的层。

---

## 3. 接线（写侧实现）

**改动**：`internal/agent/approval/gate.go` `New()` 中，把宿主审计 sink 交给 `risk` 自己的安全 sink：

```go
if o.Logf != nil {
    risk.Logf = o.Logf
}
```

（`+31` 行，含说明注释；`git diff --stat`＝`1 file changed, 31 insertions(+)`。）

**为什么这一处是"宿主注入点"**：`approval.New(Options{Logf: rt.auditf})` 是生产装配里**已经在收审计 sink** 的构造点（`cmd/wisp/run.go:728` 控制台门、`cmd/wisp/resident_approval_windows.go:404` 常驻门都把 `Logf` 接到本进程审计 sink）。把同一 sink 交给 `risk.Logf`，与在 `cmd/wisp` 写 `risk.Logf = rt.auditf` **生产效果等同**，但落在本腿写面内。`readBlacklist`（`internal/tools/mode.go:86`）在 `risk.Gate` 里对每条路径调 `logf`，故 A 档 fail-closed 实话在**路由之前**即落审计日志（不依赖该决策是否走到审批门，覆盖 A 档 Deny）。

**不是**权威、不是放宽、不是出向面：只喂本机审计轨迹（同 MODE-READ / PATH-ACCOUNT 那条）；不产面板字节、不进模型正文、不能允许/拒绝任何东西。

**残余（具名，不隐）**：宿主若**注入自己的门**（`run.go` 的 `s.gate != nil` 分支）不走本构造器，`risk.Logf` 在该宿主里仍不接线——由该宿主自理 risk 安全 sink。两个**产品**装配（控制台 `wisp run` 门、常驻门）都走本构造器。此残余属"注入式宿主"边界，非本票核心路径。

---

## 4. 反向判据（摘接线 ⇒ 理由不可见 ⇒ 红）＋还原核

- **台面**：`internal/agent/approval`（本腿写面）。
- **突变 M-337**：把 `if o.Logf != nil { risk.Logf = o.Logf }` 改成 `if false && o.Logf != nil { … }`（接线死），跑同一 `-run TestTicket337FailClosedReasonReachesAuditLog`。
- **读数**：`rc=1`，红句逐字＝
  `ticket337_reason_reader_test.go:90: approval.New did not arm risk.Logf: the fail-closed reason has no reader again (ticket 337)` ⇒ **理由无读者，失败形复现**。（日志 `D:\tmp\wisp\20261102\337r1\test-337-mutant.log`）
- **还原核**：`cp` 回原始字节，`sha256(gate.go)`＝`03f899cc072b7c613558c5729b6cffcd534081e393db5bd3be94167e7ce4a9c1`（起手＝还原**逐字节同串**，`RESTORE MATCHES`）；还原后同测试 `rc=0` PASS（日志 `.../test-337.log`）。
- **正控（尺有牙，不是恒绿）**：接线在时同一测试 `rc=0`，并**真的**在审计 sink 上打出那句实话（日志逐字，已折叠/转义，安全面）：
  `risk: A-list DENY path=…\profile\work\id_ed25519  rule=~/.aws/credentials [fail-closed, NOT a match on this rule]: the candidate "…" is not proven handle-expanded (its own walk stopped at "…\profile\work") and the anchor "…\profile\.aws\credentials" is not installed here (…), so R17 protects the ancestry chain the two share instead of proving a miss (non-overridable)`

仪器 `internal/agent/approval/ticket337_reason_reader_test.go` 两半同时钉：①`New(Options{Logf: sink})` 后 `risk.Logf` 非 nil 且直接投递到 sink（接线）；②自造 world（public API：`os`/`t.Setenv`，非 `newWorld334`）跑真 `risk.Gate`，A 档保证行与 fail-closed 实话都落 sink（流向）。`//go:build windows`（fail-closed 形只在 Windows 造得出，见 `anchor_miss_reason_334_test.go` 头注）。

---

## 5. 零回归（受影响四包四数＋门禁四数；cmd/wisp 整包按 AC#4 具名跳过）

**包级四数**（`go test ./<pkg>/ -v -count=1`，`HEAD`＝`07a6ca7b4`＋本腿改动；`RUN`＝`PASS+FAIL+SKIP`；日志 `D:\tmp\wisp\20261102\337r1\pkg-*.log`）：

| 包 | rc | RUN | PASS | FAIL | SKIP |
|---|---|---|---|---|---|
| `internal/agent/approval`（本腿改动面） | 0 | 86 | 85 | **0** | 1（`TestDefaultDeadlineWallClockMeasurement`，存量平台 skip） |
| `internal/risk`（sink 生产方） | 0 | 172 | 171 | **0** | 1（`TestSyncRegistryProbeLive`，存量平台 skip） |
| `internal/tools`（`readBlacklist` 调用方） | 0 | 253 | 253 | **0** | 0 |
| `internal/panel`（甲形不碰面，负控） | 0 | 173 | 173 | **0** | 0 |

**门禁四数**（共享工作树；日志 `.../gate-*.log`）：

| 门 | 命令 | 读数 |
|---|---|---|
| G1 gofmt(gofumpt) | `gofumpt -l internal cmd`（CI 同 scope，`.scratch` 除外） | 见下注：**跟踪集（CI 实际读的字节）rc=0 空**；工作树对 `internal cmd` 只命中 `internal/agent/approval/pending_read.go`、`internal/agent/tools.go` 两枚——**非本腿改动**，且 `git ls-files --eol` 现量 `i/lf w/crlf`＝**工作树 CRLF 检出伪影**、`.gitattributes attr/text eol=lf`；把这两枚的**跟踪 blob**（`git show HEAD:…`）落盘再跑 `gofumpt -l`＝**rc=0 空**（`.../tracked-blobs/`）。本腿两枚文件 `gofumpt -l`＝**rc=0 空**。 |
| G2 d22scan | `bash scripts/d22scan.sh` | **rc=0 `clean`**（scope：bans#1-5 internal=288 cmd=52；ban#6 frontend=116；ban#7 internal/tools=25；ban#8 design=39 frontend=116 internal=736 cmd=184） |
| G3 路径长度预算 | `bash scripts/check-path-length-budget.sh` | **rc=0 `VERDICT GREEN`**（over-budget 62 全具名在册、roster＝tree、双单位选同 62 路径） |
| G4 go vet | `go vet ./internal/risk/ ./internal/agent/... ./internal/panel/` | **rc=0** |

**具名跳过/残余（诚实记，⛔ 当绿）**：
- **`./cmd/wisp` 整包**：按票 337 AC#4 与派单（写面互斥）**不整包跑**，具名跳过。本腿改动在 `internal/agent/approval`，其生产装配点 `run.go:728`／`resident_approval_windows.go:404` 未改一个字节（`git diff --name-only` 里 `cmd/**`＝0）。
- **「六族＋262/640/902」**：这族是**依赖真桌面/真音频/球与像素**的 Windows-live 面（`openDelay=262ms chunk=640ms ballLatency=902ms` 为 `12-v3` 的实时延迟读数）；本腿**无真桌面会话**，⛔ 妄报该族逐位同。本腿给的是**包级四数（上表）＋门禁四数（上表）＋反向判据**作为可复跑零回归证据；六族/延迟族留编排者桌面轮复跑。

---

## 6. 纪律格（逐条自证）

- **⛔ 翻框／⛔ 改票面任何一行（含 `Status:`）／⛔ 加 `-done`**：本腿**未改** `.scratch/wisp/issues/337-*.md` 一字节。
- **⛔ 台账／⛔ push**：未碰 `docs/reports/pending-and-issues.md`；只 `commit` 不 `push`。
- **写面**：产出仅三处＝`internal/agent/approval/gate.go`（改）＋`internal/agent/approval/ticket337_reason_reader_test.go`（新）＋本证据件。`git diff --name-only`／新增件逐枚在 `internal/risk|internal/agent|internal/panel` 与 `docs/evidence/s1/337r1-*` 之内。
- **⛔ 禁区**：未动 `thresholds.go`／golden／`allowlist.txt`／`go.mod`／`go.sum`／`.gitignore`／`ci.yml`／D1–D47／C1–C32／R1–R9／D43；未动 `frontend/**`／`design/**`／`internal/speech/**`／`internal/audio/**`／`testdata`（`git diff --name-only` 逐枚核对）。
- **⛔ 为变绿放宽断言／删用例／`t.Skip`**：无。
- 提交：**显式 pathspec**（三枚具名），⛔ `add -A`／⛔ `amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`。

---

## 7. 残余与后续（交编排者）

1. **注入式宿主**（`run.go` `s.gate != nil` 分支）不走 `approval.New`，`risk.Logf` 仍不接线——需该宿主自理，或在 `cmd/wisp` 装配根补一行 `risk.Logf = rt.auditf`（**越本腿写面**，交下一条写腿/编排者）。
2. **乙／丙形**（面板出向／模型正文）**未落地**：AC#2 硬门要求 `A##` 具名批准（记文件/行/理由/边界/撤销口令）后方可动；且二者会把**无 216 清洗的真路径句**带向人/模型。是否上、上哪形，交编排者摆机主拍板。
3. **六族／262-640-902 延迟族**：需真桌面会话，留编排者桌面轮复跑（本腿真跑不到，⛔ 充当绿）。
