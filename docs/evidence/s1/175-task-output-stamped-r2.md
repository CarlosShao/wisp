# 175-r2 证据件｜`task.output` 进 C25 名册＋桥级三发常驻判据（J-1／J-2／J-3）

派单：`.scratch/wisp/dispatches/2026-09-28-002x-write-175-r2-roster-entry-and-bridge-level-stamp-criteria.md`
写手腿：`175-r2`｜时刻：`2026-09-28 08:4x +08`（`date` 现跑，见 §1）｜分支：`dev`
工具调用：**第 39 次用完 §7 的全部条目，此后仅余本件写作与两枚提交；未超 45 顶**（末尾 §10 给逐枚去向）

---

## 1. 起手四件（我自己现跑的锚，逐字）

| 件 | 读数 |
|---|---|
| `date` | `Mon Sep 28 00:14:00 CST 2026` |
| `git rev-parse --abbrev-ref HEAD` | `dev` |
| `git rev-parse HEAD` | `9ca57217eac8512cb4c9d2e32cdbd2bf4af66c03` |
| `git log --oneline -3` | `9ca57217 派单 175-r2（写手腿·一格一单，硬顶 45 次）…`／`0f840ac2 ledger(A357)…`／`b52de78b evidence(176-a1b 只读设计核·重派枚)…` |
| `git status --porcelain -- internal/ cmd/` | 空（`rc=0 status-lines-above`） |

### §1 基线那枚争用红（复跑两向，一字节未动）

命令：`go test -count=1 -run TestResolvePerCallBudget ./internal/risk/`

| 次 | 读数 |
|---|---|
| 第 1 遍 | `ok github.com/CarlosShao/wisp/internal/risk 1.430s` |
| 第 2 遍 | `FAIL` — `C26 Resolve: 1029978 ns/op = 1.030 ms/op (budget 1.000 ms, 1011 samples)`／`C26 budget breach: Resolve averages 1.029978ms per call, budget 1ms` |
| 第 3 遍（落完码之后，门禁那一轮） | `ok … 2.089s` |

⇒ 判为**争用／复跑绿**（三枚读数两绿一红，红那枚 `1.030 ms/op` 与派单转述的 `1.098`／`1.199` 同形）。
`internal/risk/thresholds.go`、`pathresolver_budget_norace_test.go`、`provenance_test.go`、`task_output_leg_test.go`
本程改动枚数＝**0**：命令 `git diff --numstat HEAD~2 -- internal/tools/task_output_leg_test.go internal/risk/provenance_test.go internal/risk/thresholds.go`
⇒ 空输出。

---

## 2. 本程没测什么（先说这一条）

- **端到端没测，也测不了**：`TaskRoster.Record` 生产零写者、`cmd/wisp` 只建空名册（票 176 的起跑口），
  真 `wisp run` 今日产不出带 `ArtifactPath` 的 `task.output` 桩。每枚判据自己经许可的接缝注入 record。
- **没测 `runtests.sh` 的名册两向 `comm -3`**：输出里没有名为 `comm -3`／差集的段落
  （命令 `grep -nE 'comm -3|diff|名册|roster|undeclared|unexecuted' .scratch/wisp/probes/175/r2/runtests-after.txt`
  ⇒ 只命中两行测试名 `.../not-a-roster`、`.../not-a-roster`）。本程只报仪器自证：`rc=0`、
  `runtests.sh: OK - packages=[./...] top-level: PASS=34 FAIL=0 SKIP=0, === RUN=76, '[no tests to run]'=0`。
  ⇒ **只许多不许少这一条我没独立复算，属未验。**
- **没重跑 `provenance_test.go:136`／`:139` 的敏感度**（派单已现量、且本单禁止我把它们当哨）；两枚一字未动。
- **没测陪聊／语音、面板、`internal/panel/**`**（别家常红，一行未碰）。
- **没跑** `./...` 全仓、`flip-declaration.sh`（派单禁）。
- **没量豁免窗口的字节坐标本身**（那是 risk 级 `shape_a_exemption_test.go` 的射程）；本程只量它到没到真桥。

---

## 3. 两处落点（逐枚）

### 3.1 `internal/risk/provenance.go`（名册那一处，`+4 −0`）

`git diff --numstat` 现量：`4 0 internal/risk/provenance.go`（删除列 0）。

- 现读锚：`sensitiveSourceTools` 在**起手**是 `:96-99` 八枚（派单写 `:96-101`，行号我复算过＝八枚名字两行＋`var`行＋收尾 `}`）。
- 落地后（`grep -n 'sensitiveSourceTools = \[\]string{' -A 4 internal/risk/provenance.go`）：
  ```
  99:var sensitiveSourceTools = []string{
  100-	SrcFSRead, SrcSearchContent, SrcClipboardRead, SrcSystemGet,
  101-	SrcWebFetch, SrcDocRead, SrcScreenCapture, SrcTranscript,
  102-	SrcTaskOutput,
  103-}
  ```
  ⇒ 名册 8 枚 → **9 枚**（数法：`sed -n '99,103p'` 里逗号分隔的名字，逐行 4＋4＋1）。
- 新增常量 `SrcTaskOutput = "task.output"`（`:93-95` 处，附两行说明），既有八枚的名字、顺序、注释措辞一字未改；
  `Mark` 的冻结文档块（起手 `:468-473`）未碰。
- 为什么落这里：派单 §2.1 给的修法，也是 `175-c1` 裁过的形（名册＝举例定义 ⇒ 实现缺陷修复，非契约变更）。

### 3.2 `internal/tools/bridge.go`（**只有说明过期，产码一行未动**）

`git diff -U0 internal/tools/bridge.go | grep -E '^[+-]' | grep -vE '^(\+\+\+|---)' | grep -vE '^[+-][[:space:]]*//'` ⇒ **空输出**
⇒ 变动的 6 行全是注释（`numstat` 为 `6 4`，成对来自那段说明改写）。
- 过期说明原文（起手 `:558-561`）写着：`task.output is not in the SPEC-06 §5 roster yet … no production mark carries a non-empty hostPath today`。
  名册一补，这句当场变假 ⇒ 就地改成：`task.output joined that roster with ticket 175-r2, which is the first day a non-empty hostPath is reachable here at all`，
  并点名判据文件。
- ⚠ **派单 §2.2 的前提我复算不符，报回**：「177 只把载具修好了、没人往里放东西」——现量**有人放**：
  `internal/tools/task.go:253-255` 逐字 `if box := hostPathBoxFromCtx(ctx); box != nil { box.set(rec.ArtifactPath) }`，
  `bridge.go:457` 起盒、`:533` 交盒、`:571` 送 `MarkWithHostPath`。惰性只在**一处闸门**：`bridge.go:565`
  （起手 `:563`）`if b.prov == nil || dec.TaskID == "" || !risk.IsSensitiveSource(dec.Tool) { return }` 早退在盖戳之前。
  ⇒ 名册一补，载具与戳同时开始真跑；**本枚不需要任何 bridge.go 产码改动**，我也就没造（造它是多余的第二处修法）。

---

## 4. 三发常驻判据（＋两枚控条）逐枚

文件：`internal/tools/ticket175r2_stamp_live_test.go`（新建；未追加进 `task_output_leg_test.go`／`provenance_test.go`）
枚数命令：`grep -c '^func Test.*175r2' internal/tools/ticket175r2_stamp_live_test.go` ⇒ `5`
夹具：真 `Bridge`＋真 `BuiltinFSEntries`／`BuiltinTaskEntries`＋真 `risk.NewProvenance`＋真临时目录＋真文件；
注入面只有 `TaskRoster.Record`（票 176 起跑口缺席 ⇒ 每腿自己注入），**零 mock 顶真件**。

| 编号 | 钉什么 | 未修码（锚 `9ca57217`）| 落地后 | 把它防的做宽法做进去 ⇒ 红 |
|---|---|---|---|---|
| **J-1** `TestTaskOutputStampedOnRealBridge175r2` | `task.output` 的成功答复必须盖 C25 戳、`RuneLen>0`、`Inspect(notify)` 命中且 `src="task.output"`、盖了戳就得开活作用域 | **红**（前提也红：`污点表=[]`） | **绿** | ＝起手锚本身：名册摘掉 `SrcTaskOutput` ⇒ 同一枚红（读数见 §4.1） |
| **J-1b** `TestNonContentResultStaysUnmarked175r2` | 非外部内容（`fs.list`）不得盖戳、不得开账——钉住凡成功结果都盖那形 | 绿（今日谁都不盖） | 绿 | **变异 M3**：删 `bridge.go:565` 的 `!risk.IsSensitiveSource(dec.Tool)` ⇒ `--- FAIL` ×2（`盖了 1 枚来源标记（want 0）`＋`非外部内容源也开了 C25 污点 scope`），另两枚仍绿 |
| **J-2** `TestForeignMentionOfHostPathStillHitsR4OnRealBridge175r2` | 宿主桩已盖戳的前提之下，外来正文逐字点同一条 artifacts 路径、模型随后真 `fs.read` 它 ⇒ 仍判 L2、`RulesHit` 含 R4、`SessionOverrideBlocked=true`、卡上来源不许是 `task.output`；对照支（只有宿主自己那枚 mark）须判 L0 | **响不了**：`precondition: task.output 没盖戳，这一发就在测别的东西` ⇒ 红在前提上（真桥今日无这条码路） | **绿** | **变异 M1b**：`Inspect` 里按**值形状**豁免（候选串含 `/` 就不匹配）⇒ `--- FAIL` ×3：`判成 "L0", want L2`／`rules_hit = [], want R4`／`R4 命中必须标成…`。⇒ 这发就是洗戳警报器 |
| **J-3** `TestHostMintedPointerRereadStaysCleanOnRealBridge175r2` | 宿主自己在桩里写下的那条路径：先证它**盖了戳**，再证 `Inspect(fs.read, path=P)` 不命中、真 `fs.read` 判 **L0**、原文逐字节回全 | **响不了**：`precondition: task.output 没盖戳，豁免这条腿根本没跑：[]` | **绿** | **变异 M2**：删 `task.go:254` 的 `box.set(rec.ArtifactPath)` ⇒ `--- FAIL`：`Inspect(fs.read, path=…) 命中，豁免没送到真桥`＋`续读宿主自己的指针被判成 "L2", want L0（含 [R4]）`；**同时票 164 那两枚续读腿当场红**（§5） |
| 普查腿 `TestEveryRegisteredToolIsClassifiedForMarking175r2` | 每枚注册内置工具必须有一句分类，且那句话得与名册同向（票 175 AC#3 判据要防下一程自动通） | 红：`task.output 不在名册里，分类表却写着 marked…` | 绿（8 枚名字逐枚 `t.Logf`） | 名册与分类表任一方向单独变动 ⇒ 红（起手锚那发就是它） |

### 4.1 未修码那两向读数（真桥锚，逐字）

命令：`go test -count=1 -run '175r2|TestLongOutputPointerRecoversEveryByte|TestPointerPast256KiBIsNotFullyReadable' ./internal/tools/`

```
--- FAIL: TestTaskOutputStampedOnRealBridge175r2 (0.00s)
    ticket175r2_stamp_live_test.go:159: task.output 读回来的外部内容没有盖 C25 来源标记：scope "175r2-j1" 的污点表=[]。…
    ticket175r2_stamp_live_test.go:168: 外来的任务输出被读回来之后 R4 一侧不再看得见它：Inspect(notify) 未命中，污点表=[]
    ticket175r2_stamp_live_test.go:174: 盖了戳却没有开 C25 污点 scope，task="175r2-j1"：…
--- FAIL: TestForeignMentionOfHostPathStillHitsR4OnRealBridge175r2 (0.01s)
    ticket175r2_stamp_live_test.go:232: precondition: task.output 没盖戳，这一发就在测别的东西（票 175 的破口本身）：[]
--- FAIL: TestHostMintedPointerRereadStaysCleanOnRealBridge175r2 (0.01s)
    ticket175r2_stamp_live_test.go:306: precondition: task.output 没盖戳，豁免这条腿根本没跑：[]
--- FAIL: TestEveryRegisteredToolIsClassifiedForMarking175r2 (0.00s)
    ticket175r2_stamp_live_test.go:367: task.output 不在名册里，分类表却写着 "marked: C25 roster since ticket 175-r2, …"
```
⇒ 派单 §3 那句「J-1 今天不响」我按**判据颜色**复算成反向：J-1 今日**红**（红＝缺戳本身在响）；
今日**不响的是那道戳**（`污点表=[]`）。两枚反向腿（J-2／J-3）今日红在**前提**上，正是「载具没 traffic」的形状。

### 4.2 落地后那两向读数

命令：`go test -count=1 -v -run '175r2|TestLongOutputPointerRecoversEveryByte|TestPointerPast256KiBIsNotFullyReadable' ./internal/tools/`

```
--- PASS: TestLongOutputPointerRecoversEveryByte (0.01s)
--- PASS: TestPointerPast256KiBIsNotFullyReadable (0.10s)
--- PASS: TestTaskOutputStampedOnRealBridge175r2 (0.00s)
--- PASS: TestNonContentResultStaysUnmarked175r2 (0.01s)
--- PASS: TestForeignMentionOfHostPathStillHitsR4OnRealBridge175r2 (0.03s)
--- PASS: TestHostMintedPointerRereadStaysCleanOnRealBridge175r2 (0.02s)
--- PASS: TestEveryRegisteredToolIsClassifiedForMarking175r2 (0.00s)
ok  	github.com/CarlosShao/wisp/internal/tools	0.235s
```

### 4.3 变异一发的落地证明与还原（三枚都用 `grep -n` 先证、`grep -c` 后证）

| 变异 | 落地锚（先 `grep -n`）| 改后现场 | 还原后 |
|---|---|---|---|
| M1a（废） | `709: if !gateOpen && !isPathKey(k) {` → `if !gateOpen {` | **全绿** ⇒ 这一支不是 `fs.read` 的 `path` 候选走的码路（原因未深挖，属〔我现量、未定性〕），M1a 拿不到读数 ⇒ **不作数**，改跑 M1b | `grep -c '!gateOpen && !isPathKey(k)'` ⇒ `1` |
| M1b | 在 `Inspect` 候选匹配处插入按值形状豁免（含 `/` 的候选不匹配） | J-2 红 ×3，J-1／J-3／164 腿绿 | `grep -c 'ContainsRune' internal/risk/provenance.go` ⇒ 无命中 |
| M2 | `254: box.set(rec.ArtifactPath)` → `_ = box` | J-3 红 ×2、J-2 红 ×1（卡上来源变 `task.output`：`R4: 包含来自 task.output 的内容`）、**票 164 两枚红** | `grep -c 'box.set(rec.ArtifactPath)'` ⇒ `1` |
| M3 | `565: if b.prov == nil \|\| dec.TaskID == "" \|\| !risk.IsSensitiveSource(dec.Tool) {` 删去第三项 | J-1b 红 ×2，J-1／J-3 绿 | `grep -c '!risk.IsSensitiveSource(dec.Tool)'` ⇒ `1` |
- 备份件（只建不删）：`.scratch/wisp/probes/175/r2/bak1-{provenance,task,bridge}.go`；终态差分见 §7 末（`git diff --numstat` 只有名册那 4 行与那段注释）。

---

## 5. 票 164 那两枚续读腿（豁免真生效的正面凭据）

命令：`go test -count=1 -v -run 'TestLongOutputPointerRecoversEveryByte|TestPointerPast256KiBIsNotFullyReadable' ./internal/tools/`

| 腿 | 加名册前（锚） | **加名册后** | 变异 M2（删宿主声明）后 |
|---|---|---|---|
| `task_output_leg_test.go:165 TestLongOutputPointerRecoversEveryByte` | 绿（整包 `ok`，见 §6 基线） | **绿** | `--- FAIL` `:220: the pointer must be re-readable with fs.read, got {Text:L2 审批通道尚未接入（票 21），已拒绝执行 IsError:true RiskLevel:L2 ErrorClass:user_rejected}` |
| `task_output_leg_test.go:346 TestPointerPast256KiBIsNotFullyReadable` | 绿 | **绿** | `--- FAIL` `:381: the first segment must read back, got {… L2 … user_rejected}` |

⇒ **正面凭据成立且不是恒真**：这两枚加完名册仍绿＝豁免送到了真桥；M2 把它们打红＝它们真能看见豁免消失。
派单 §3 要求「若红＝停手上报」——**没红**，两枚一字未动（§1 的 numstat 命令）。
⚠ 一处**加强**说明：这两枚走 `NoGate{}`，红的方式是 L2 把结果顶成 `已拒绝执行`（不是判级断言）；
所以我另立 J-3 直接断言 `RiskLevel`＋`Inspect` 不命中，避免干净是没盖戳骗来的（J-3 里那枚 `mark175r2(...)==nil 即 fatal` 就是干这件事的）。

---

## 6. 两条禁令逐条自证

- **(i) 不给 `internal/agent.Loop` 加收收 `taskID` 的导出方法**
  `git diff --name-only HEAD~2 -- internal/agent/` ⇒ **空**（本程一枚 `internal/agent/**` 文件都没打开）。
  `grep -rn 'func (l \*Loop) [A-Z][A-Za-z]*(' internal/agent/*.go | grep -iE 'taskid'` ⇒ **空**；同形导出方法总数
  `grep -rn 'func (l \*Loop) [A-Z][A-Za-z]*(' internal/agent/*.go | wc -l` ⇒ `6`（六枚全在起手锚就存在，本程没加第七枚）。
  ⚠ **「G3 安静 ≠ 合规」怎么自证**：我不拿 `gate-clauses.sh` 的 G3 rc 当通行证，改用**我名下差分枚数为零**这条独立尺——
  G3 那条腿按字面 `taskID` 抓，改名或包进 struct 就纹丝不动；而 `git diff --name-only … -- internal/agent/` 为空
  与参数名无关，改名／包 struct／改签名都藏不住。另附一句：本程唯一新增的可调用东西是一枚包内未导出测试夹具
  （`stamp175r2Bridge`），不在 `agent` 包里，也不在 C1 面上。
- **(ii) 不把宿主内部 artifacts 写入做成受门控的 Tool**
  `git diff --name-only HEAD~2 -- internal/tools/tool.go internal/agent/spill.go` ⇒ **空**（C1 面与 spill 本体未碰）；
  `git diff HEAD~2 -- internal/tools/bridge.go | grep -E '^\+' | grep -cE 'Register|Decl\{|Entry\{'` ⇒ `0`（没注册任何工具、没造 `Decl`）；
  `grep -n 'gated tool' internal/agent/spill.go` ⇒ `30:// a gated tool, and it lands in the memory store's artifacts directory so the`
  （⚠ 派单转述的行段 `:29-31` 我复算成 `:30` 那一行，措辞是 `a gated tool` 不是 `deliberately not a gated tool` 整句；
  判据未受影响，报回备查）。**我没觉得"非做不可"**：名册那一枚就是 C25 认下后台产物的全部所需，无需新门。
- 附带：`b.prov == nil || dec.TaskID == "" || !risk.IsSensitiveSource(dec.Tool)` 那枚闸门**留着没拆**（拆它＝M3 那形，已被 J-1b 判红）。

---

## 7. 门禁五枚＋名册差集＋`gofumpt`

| 门禁 | 命令 | 读数 |
|---|---|---|
| 1 分包测试 | `go test -count=1 ./internal/risk/ ./internal/tools/` | `ok internal/risk 7.430s`／`ok internal/tools 18.159s`；另 `go test -count=1 -v ./internal/tools/` ⇒ `=== RUN` 计数 `178`、`--- FAIL` 计数 `0` |
| 1b 基线（改前同命令） | 同上，起手锚 | `ok internal/risk 4.036s`／`ok internal/tools 15.120s` |
| 2 d22scan | `sh scripts/d22scan.sh`（存 `.scratch/wisp/probes/175/r2/d22scan-after.txt`） | `rc=0`；`d22scan: clean - no D22 ban violations`；扫描面 `internal/ 428`／`cmd/ 45` 枚 Go 文件（含注释与 `_test.go`） |
| 3 成对普查尺 | `sh .scratch/wisp/probes/154/gate-clauses.sh`（存 `gate-clauses-after.txt`） | `rc=0`；逐字：`# 腿数＝14 声明与实测不符＝0`＋`# 腿数断言：名册=14 声明=14 记账=14 缺腿=0 空头声明=0`＋`# 腿数断言＝相符（名册上每一腿都记了账）`。**G2／G5／G6 一枚没被我顶红**（本程未动 `cmd/**`；新注释不含 `OpenTask|CloseTask` 字面；测试文件在 G2 的排除名单内） |
| 4 d22scan 自审 | `bash tools/d22scan/runtests.sh -C tools/d22scan ./...`（存 `runtests-after.txt`） | `rc=0`；`runtests.sh: OK - packages=[./...] top-level: PASS=34 FAIL=0 SKIP=0, === RUN=76, '[no tests to run]'=0`（名册两向 `comm -3` 段落在输出里不存在 ⇒ 见 §2，未独立复算） |
| 5 `gofumpt` | `"$(go env GOPATH)/bin/gofumpt.exe" --version` ＋ 对三名 `-l` | `v0.12.0 (go1.27.1)`；`-l` 对 `internal/risk/provenance.go`／`internal/tools/bridge.go`／`internal/tools/ticket175r2_stamp_live_test.go` ⇒ **空输出，`rc=0`**（PATH 里没有 `gofumpt`，`command -v` ⇒ `command not found`，我用 `$(go env GOPATH)/bin` 现取的） |
| 名册差集（C25 那一枚） | `grep -n 'sensitiveSourceTools = \[\]string{' -A 4 internal/risk/provenance.go` | 8 枚 → **9 枚**（只许多、不许少：摘掉的枚数＝0） |
| 写面差分 | `git diff --numstat HEAD~2` | `4 0 internal/risk/provenance.go`／`6 4 internal/tools/bridge.go`（全注释）／新增 `internal/tools/ticket175r2_stamp_live_test.go` |

---

## 8. 被拒／没成功的调用

- **取数后**（写面内、被工具拒的只有 1 枚）：第 25 次 `Edit` 对 `ticket175r2_stamp_live_test.go` 报 `0 occurrences`——
  我先发的那枚 Edit 已把目标串改掉，第二次是重复编辑；随后用正确的 `old_string` 成功（第 26 次）。**未越权、未伤文件**。
- **取数前**（仪器没给我要的数，非被拒）：`gofumpt` 不在 PATH（`command not found`，`rc=127`）⇒ 改从 `$(go env GOPATH)/bin` 取到；
  `runtests-after.txt` 里找不到 `comm -3` 段（`grep` 只命中两行测试名）⇒ 记为没测（§2）。
- 无权限拒绝、无门被拦、无 `git` 操作被拒。

## 9. 有没有跑过删除命令

**没有。** 零 `rm`／`rmdir`／`unlink`／`git clean`／`git restore`／`checkout .`／`stash`／`--amend`／`rebase`／`reset`／`merge`／`worktree`。
三枚变异全用 `cp` 备份＋`cp` 还原（`.scratch/wisp/probes/175/r2/bak1-*.go` 只建不删）。未 push。

## 10. 工具调用第几次停的

**未超，共 45 次封顶内：写完本件＋两枚提交实际用到第 43 次**（本行提交时定格，逐枚去向：起手锚 2、§1 复跑 4、读票面与夹具 8、读产码 12、
判据写作与两处 Edit 4、两向读数 2、变异 3、门禁 4、名册计数 1、提交 2、本件 1）。
第 39 次之后只做了证据件、票面追加与文档提交，**没再开新的取数**（§2 里那两枚没复算的就为此留下）。

## 11. 伪授权两栏

| 栏 | 内容 | 出处 |
|---|---|---|
| 〔别人现跑／我转述、**本程复算为真**〕 | 名册起手八枚；闸门早退在盖戳之前；`task_output_leg_test.go:165`／`:346` 是已裁续读腿且禁改；`provenance_test.go:136`／`:139` 不当哨；G2 的 `want_n 2` 里一行是 `cmd/wisp/run.go:564` 的注释（该腿排掉 `*_test.go` 与 `bridge.go`，见 `gate-clauses.sh:361`）；`TestResolvePerCallBudget` 是争用红 | 派单 §1／§3／§4／§5；我各附现量命令见 §1／§3／§5／§6／§7 |
| 〔别人现跑／我转述、**本程复算不符或未验**〕 | ①「177 没人往载具里放东西」⇒ 现量**有人放**（`task.go:253-255`），惰性只在名册闸门（§3.2）。②「J-1 今天不响」⇒ 现量 J-1 **红**、今日不响的是戳本身（§4.1）。③「canary 可直接搬回」⇒ 现量其本体引用的 `marksProvenance`／`outsideContentTools`／`bridge_mark_provenance_ticket175_test.go` 全树 `grep -rln` ⇒ **零命中**（r1 那形随裁定乙被撤，从未落 `.go`），故我按现裁形（名册）改写并换名 `…175r2…`（§4）。④「`spill.go:29-31` 逐字 deliberately not a gated tool」⇒ 现量在 `:30`、句子是 `a gated tool, and it lands…`（§6-ii）。⑤ `runtests.sh` 的名册两向 `comm -3` ⇒ 输出里没这段，**未复算**（§2） | 派单 §2／§3／§4／§5 原句；复算命令逐枚在正文 |

## 12. 凭据值

零抄录。本程只出现工具名（`task.output`／`fs.read`）、路径与探针串常量（`175r2-foreign-7c3d91` 是夹具字串，非凭据）；未读 `.env`、未读密钥存储、未打印任何 provider/baseURL 值。

## 13. next=

`next=176（后台起跑口）`。缺的那一发：`cmd/wisp` 侧 `execute()` 用现成 `RunAsync` 的 `Result` 去调 `TaskRoster.Record`（`176-a1b` 裁的甲形），
并把 `internal/agent.Spiller` 落的那条 artifacts 路径填进 `TaskOutput.ArtifactPath`——今天 `Record` 生产零写者，所以我这三发全是注入 record 的桥级判据，
端到端那条腿（真 `wisp run` 产出一枚带指针的桩并让续读走通）**还站着不动**（票 177 AC#6 明写 176 不得早于 177 合入，现序符合）。
另：取消判据挂 `tools.TaskRoster`＋`tools.Stopped(ctx)`，别挂桥、别挂 Loop；`taskID` 不得裸送（与模型侧 call id 同目录同前缀会撞 `spill.go` 的静默换字节）。

**票 177 的 AC#3 能不能因此翻勾：不能。** 本程让 177 的载具**第一次真跑**（J-3 绿＋M2 能打红，是豁免在真桥上生效的凭据），
但 AC#3 的端到端那一发要的是生产里真存在一枚带 `ArtifactPath` 的 `task.output` 桩，而写它的人（起跑口）属票 176。
⇒ 177 的 AC#3 至多能改注为「桥级已生效、端到端待 176」，翻勾归编排者裁，**本程不自勾任何框**。
