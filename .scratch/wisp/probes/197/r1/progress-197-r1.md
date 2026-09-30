
---

## 197-r1 进度（后端写码腿，起钟 2026-09-30 11:37:58 +0800／落笔 12:10:0x +0800）

⛔ **本腿未跑任何 `go test`**（整包／单包／`-run` 定向都没跑）——门禁窗口那段时间归独占桌面采样的验收腿。
本文件一枚 AC 复选框都没碰（勾框归编排者）。台账与 HANDOVER 未写（归编排者）。
交付＝一枚 commit `a818df46`：`cmd/wisp/subagent_selfapproval_197_test.go`（710 行，**产码零改动**；
`git show HEAD:<path>` 与工作树逐字节 diff rc=0 验过）。

### ① 起手复跑票面「现量」表：五行里三行已过期（具名更正；锚 `bf26dfc4`，收尾锚见 §⑤）

| 票面那行（锚 `fccaf3e3`） | 票面读数 | 本腿现跑 | 尺 |
|---|---|---|---|
| Go 侧「子代理」实体 | 0 命中 | **9 个文件**（`internal/tools/subagent_197.go`、`internal/panel/subagent_roster_197.go`、`cmd/wisp/panel_pump.go`、`internal/streamkey/streamkey.go` 等） | 票面原尺逐字复跑 |
| 中文「子代理」在 Go 侧 | 0 命中 | **3 个文件** | 票面原尺逐字复跑 |
| `TaskRoster.byTask` 的 `Record:123`／`Look:139`／`Count:152` | — | **漂到 215／268／281**，且名册已有 13 枚方法（新增 `Descendants:305`、`PublishSubagent:342`、`MarkRoot:355`、`TryAcquireSubagentSlot:363`、`InFlightSubagents:390`、`RunningSubagentIDs:402`、`AttachCancel:420`、`DetachCancel:434`、`Cancel:447`、`WatchRow:248`） | `grep -n "func (r \*TaskRoster)" internal/tools/task.go` |
| 流式载体 `pump.go:301`／`:339`，上限尺 `:275`／`:287-290` | 键上限 32、溢出「合并而不是丢」 | **上限常量在 `pump.go:406`（`DefaultStreamKeys = 32`）＋ `:408` `StreamKeyHardCeilingMultiple`**；溢出那一支已经改成「显式声明截断、不再并键」：`:509` 注释逐字 `no key is ever folded into another`，判定体在 `:514`／`:521` ⇒ **票面这一行对子代理的那条反对意见已被现形解决** | `grep -n "DefaultStreamKeys\|maxKeys" internal/panel/pump.go` |
| 面板今天能读到的：`Snapshot` 4 枚键 | 4 | **6 枚**（`internal/panel/composer.go:57-`：`pending`／`results`／`composer`／`generatedAt` ＋ `instructions`（票 200）＋ `tasks`（本票））；`Results: rt.stream.Chunks` 从 `run.go:442-448` **漂到 `run.go:615`** | `sed -n '56,66p' internal/panel/composer.go`；`grep -n "Results:   rt.stream.Chunks" cmd/wisp/run.go` |
| ⛔ 撞上的两枚既有定案 | `PLAN.md:1544`／`:1447` | **原文在场、一字未改**：`1544` 行首逐字 `\| **REJECTED** \| 多 Agent 协作 \| 草案 6.6 明确不做 \|`，`1447` 行内含「不得开放多任务并发」；且 `git status --porcelain docs/PLAN.md docs/specs tools/d22scan/allowlist.txt` **全空** | `sed -n '1544p' docs/PLAN.md`；`sed -n '1447p' docs/PLAN.md` |

⇒ 票名尾巴那句「go-side-has-zero-subagent-entity」连同上面三行现量都是**历史状态**（与本票 §「要 owner 带给界面那支的一跳」第一条同一口径）；上表是本腿在 `bf26dfc4` 上重新读到的。

### ② 本格交付＝AC#5 的那枚正控（票面 §09-29 那节留下的唯一 Go 侧格）

09-29 普查把 AC#5 判成〔仅文档〕，并留下了一句可执行接缝：「最便宜的合法接缝＝`cmd/wisp` 的 `package main`
测试……记下 `grant` 后断言 `Native().Allow(childCorr, capturedGrant)` 必须失败」。本腿把那一句变成读数。

**(1) `Test197SubagentSelfApprovalIsRefusedAndTheHostAnswerLands`（真 run 装配，产码零改动）**
- `task.spawn` 派一枚**实名册子代理**（`rt.bridge.Execute` 走真桥；名册行由 `PublishSubagent` 与子环路自己的准入钩子发布，本文件不造行），
  随后**用孩子自己的 task id** 通过真桥撞出两张真 L2 卡（`fs.write` 打到 `t.TempDir()`，即 R2⇒L2 那一形）。
- **正控先行**：第一张卡由宿主账上那枚真令牌经 `Replies.Allow` 允许 ⇒ 调用真的执行、文件真的落盘。
  没有这一发，后面每一发「被拒」都只是「这条路由今天没在跑」。
- **自批的每一发**：空令牌／假令牌（运行时从真令牌改尾 4 字，源码里没有令牌字面量）／
  自称 `Source="subagent-self"` 的 `DecideFromNative` ⇒ 逐发 `ErrBadGrant`；已结卡的重放 ⇒ 只允许 `ErrUnknownCorrelation`。
- **第二张卡**：跨卡借证 `ErrBadGrant`；真令牌在 `DecideFromPanel` 上露一次面即被按泄露烧掉（`ErrPanelAllow`），
  此后连原生侧花它也 `ErrBadGrant`；收口是宿主的 `Reject`，文件始终不存在。
  审计侧同时要求 `FORGED-OR-STALE`／`PANEL-ALLOW-REJECTED`／`ANSWER-ALLOW` 三行在场（被拒与放行是两种读数，不许靠推断）。
- **载具扫描（卡还挂着时取样，不是事后）**：这一程自己发布的快照字节／工具侧 `TaskOutput` 全字段／
  `StreamLog` 全字段／`task.spawn` 回给父模型的正文／stdout／stderr／`<data>/logs` 持久日志
  ⇒ **令牌本身一处都不许出现**（`bindDigest` 绑的是事不是人，「拿不到」是这一层唯一防线）。
- 设计裁定落点复跑：**子级永不自批**＝本格；**结论复用现成源名 `task.output`** 的尺
  `git grep -n "subagent\.output" -- 'internal/**/*.go' 'cmd/**/*.go'` **＝ 1 命中，且那一行逐字是
  「No "subagent.output" name is invented」**（`internal/tools/subagent_197.go:35`）；
  **父取消不级联且写给模型看**＝票面 §09-29 ② 已判为做到，本腿未动（孩子那一侧的停止出口仍在票 220／221，本腿不回收）。
- **⚠ 本腿自己发现并当场改掉的空转判据（具名，不藏）**：起初写了「取到的那枚快照字节要出现在持久台账里」——
  这条**今天恒红**：`publishPanelSnapshot` 把最新包留在 `rt.lastSnap`／`lastSnapBytes`，收尾那一发会覆盖挂卡时那枚，
  而本用例**故意不调** `publishPanelSnapshot`（197-r4 的 `case 1 never calls publishPanelSnapshot` 同一纪律），
  所以永远读不到挂卡那一枚的 sha ⇒ 判据与产码事实互斥。改成两条真话（`snap.Pending` 里认得出这张卡 ＋ 这一程至少落账过一枚快照），
  并把做不到的那半（sha-to-ledger tie）**具名留在注释里**交给整包复跑，不留一条恒红判据给别人。

**(2) `Test197NoAllowDoorIsReachableFromASubagentsAssembly`（纯反射，不起 run）**
- `tools.SubagentDeps`／`tools.TaskDeps`／`tools.Options`／`agent.Options` 的**静态类型图**里没有任何类型
  声明 allow 类方法（`Allow`／`Native`／`DecideFromNative`／`DecideFromPanel`／`GrantNonce`），
  且 `tools.Gate` 的方法集仍恰好是 `{PendingWindow, PendingApproval}`——这才是「孩子自带允许出口」的可达性那一半。
  名字集刻意不含 `Reject`／`Veto`：那两个方向混进来会把两件事判成一件事（拒绝方向那一格是票 220 的地界）。
- **自带种门的正控**（第 64 条；负向尺必配「种 X 必响」）：文件里植了一枚带 `Allow`／`Native` 的假门载体，
  扫描器读不出那两枚名字就 `t.Fatalf`——**照不见门的仪器不许签发「没有门」的判据**。

### ③ 期望编排者跑哪几发、期望看到什么（⛔ 未跑测试，等门禁窗口；本腿不写「应该能过」）

| # | 跑法（都要带 `PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"`，否则会以 `exit status 0xc0000135` 失败且不打 `--- FAIL`） | 期望 |
|---|---|---|
| 1 | `go test ./cmd/wisp/ -run 'Test197NoAllowDoorIsReachableFromASubagentsAssembly' -count=1 -v` | `--- PASS`；秒级（这一发不起 run）。若红：**先读 findings 里指名的「类型.字段.方法」**，那是判「是不是新门」的原始读数，不是先删断言 |
| 2 | `go test ./cmd/wisp/ -run 'Test197SubagentSelfApprovalIsRefusedAndTheHostAnswerLands' -count=1 -v` | `--- PASS`。本用例把 mockllm latency 设成 2500ms（`runTextTask` 在根环路结束后就 `rt.close()`，那 2.5s 就是探针的工作窗口），健康读数个位数秒；红路径最坏约 50s（写调用 ctx 30s 上限）。`t.Logf` 那行期望：`child=<id> corr1=<id> corr2=<id>`、六发 error 是 `ErrBadGrant`×4 ＋ `ErrPanelAllow` ＋ `ErrBadGrant`、`replay=ErrUnknownCorrelation`、`落盘=true/false` |
| 3 | 整包 `go test ./cmd/wisp/ -count=1` | 本腿两枚 PASS；红名册相对 `subagent_blocked_197_test.go` 那一族**不应多出一枚**。⛔ 不许为变绿放宽任何断言 |
| 4 | 变异对照 **M1–M6**（清单逐条写在该测试文件末尾注释里） | 每发要看到**具名那一格**红：M1 面板载具带 Grant／M2 控制台打印令牌／M3 `SubagentDeps` 多一枚能答的字段／M4 `tools.Gate` 多出答复方法／M5 让 `Request.Source` 变成权威／M6 `Replies.Allow` 不要令牌。⚠ M1、M2 会**同时**打红冻结件 `internal/panel/l2_grant_boundary_test.go`——那是同族互证，不是冲突，不许为放行其中一枚去改另一枚 |

### ④ 撞钉预检（第 64 条：读断言，不是只 grep 新符号名；本腿未跑任何包，凡未读者一律标〔未跑包，此条来自读断言〕）

| 今天绿着的钉子 | 它的射程为什么盖不到／盖得到本腿，判词 |
|---|---|
| `cmd/wisp/leg_dispatch_gate_133_test.go` 的 `runRosterReds135`（`:1632` 起，实现逐条读过） | 只查「账上写的 `covered=test` 名字在这轮 binary 里必须 startable」，**不查反向** ⇒ 新增测试枚不欠账、不顶红。〔未跑包〕 |
| `cmd/wisp/leg_sink_gate_131_test.go` | 走 `func main` 的分派形状与 `installLogSink` 可达性，读源码不读测试名册 ⇒ 不受新增测试文件影响。〔未跑包〕 |
| `cmd/wisp/subagent_carrier_197_test.go:430 TestSubagentStreamKeyHasOneMintSite` | `stringLitSites197` 显式跳过 `_test.go`（`:462`），判的是「非测试文件里恰好一处字面量」；本腿既没在产码里、也没在测试里写那枚字面量 ⇒ 不动它。〔未跑包〕 |
| `internal/tools/subagent_197_test.go:846 Test197SubagentHasNoSelfApprovalOutlet` | `len(SubagentDeps 字段) != 5` 即红 ⇒ 本腿零字段改动；它的 `AdmitTask == nil` 那一支与本腿反射腿相邻互补、不重钉。〔未跑包〕 |
| 冻结件 `internal/panel/l2_grant_boundary_test.go`（**未动一字**） | 它钉「面板路线不得承载裁决」，本腿钉「子代理读得到的一切不得载有原生令牌」——互补；M1 那发同时打红两边属预期（§③ 第 4 行）。 |
| `internal/agent/approval/{queue,ticket87,ticket97,ticket146}_*_test.go` | 都在包内自造 Gate，不读 `cmd/wisp` 源码 ⇒ 不受影响。⚠ 本腿是 `Native()` 在 **`cmd/wisp` 里的第一个调用者**（普查 §① 的预告），但全在 `_test.go` ⇒ 账 `A418 ③`「生产零消费者」那句话今天仍然成立，没有被本腿偷偷改掉。〔未跑包〕 |
| 全仓仪器 `tools/d22scan/d22scan.exe` | **本腿自跑过＝clean**（这条不是读断言）：ban #1（裸 `go func(`）与 ban #3（明文密钥）**不扫 `_test.go`**（`main.go:665`／`:858` 跳过），ban #8（emoji）扫；本文件的 ⚠／「」全在注释里，字符串字面量里没有 U+2600–U+27BF／U+2B00–U+2BFF／U+FE0F。 |

### ⑤ 本腿跑过的门（收尾现量；最后一次复跑与落笔同钟 **12:09:55 +0800**，跑的就是 `a818df46` 那枚提交里的字节）

- `go build ./...` rc=0；`go vet ./cmd/wisp/` rc=0；`go vet ./internal/panel/` rc=0。
- ⛔ **`go vet ./internal/tools/` rc=1，非本腿所动**：别人未入库的 `internal/tools/grant_test.go`
  （`grantsOf redeclared`／`mustCanonical redeclared`／`undefined: grantSourceType`）——票 224 的活在飞。
- ⚠ **未归因（第 74／75 条）**：派单给的起跑读数「`internal/tools` 161 顶层声明／161 PASS」本腿**不复现、也不改口径去凑**——
  源码侧 `^func Test` 枚数在 12:05→12:07 之间从 141 涨到 177（`grep -rhE '^func Test' internal/tools --include='*_test.go' | wc -l`；
  期间 `3b78c246 feat(224)` 落库、`grant_test.go` 仍是未入库新件），`internal/panel` 99 枚、`cmd/wisp` 122 枚（含本腿 2 枚）。
  ⇒ 「161 PASS」要复现就得跑包，本腿不跑 ⇒ **不当成自己的读数，也不判谁漂了**。
- `gofumpt -l cmd/wisp/subagent_selfapproval_197_test.go` 空。
- `tools/d22scan/d22scan.exe`：`clean - no D22 ban violations`。口径提示：`ban #8 internal/=466`、`cmd/=64` 是**被扫文件数**（含 `_test.go`），不是违规数。
- 冻结件与别人地界零改动：`git status --porcelain docs/PLAN.md docs/specs tools/d22scan/allowlist.txt internal/panel/tokens_fourway_test.go internal/panel/l2_grant_boundary_test.go internal/perm/ticket90_persist_test.go` **全空**；
  `frontend/**`／`design/**` **零读零写**（本节与测试文件里没有来自那两层的任何结论）。
- Git：`--only` ＋ 显式 pathspec 一枚文件；未 push；未 amend／reset／rebase／stash／`checkout .`／clean；未建 worktree。
  提交前后别人的脏文件（`internal/tools/grant*.go`、`cmd/wisp/run.go` 的会话令牌活、`.gitignore`、`design/**` 那 16 枚删除、
  `docs/evidence/s1/152-*`、`.scratch/wisp/probes/161/r6/logs/**`）仍在原位，本腿一枚未提交。
- 台件只建在 `.scratch/wisp/probes/197/r1/`（本节那一枚，**只建不删**；这份文本随后原样追加进本票面，票面是权威）；
  仓库其它目录没有本腿留下的东西。

### ⑥ 没做完／留给编排者（⛔ 不留半成品不声明）

1. **AC#5 仍不该被勾**：本腿交付的是那枚**正控**（普查 §① 欠的那格），本格是否成立要**非实现者**裁
   （裁决者≠实现者：`SPEC-12 §4.3` #1/#3、`issues/README` 硬约束末条、D22 双角色）。本腿一句不替它判。
2. **AC#0／#1／#2／#3／#6 一枚未动**：AC#0 的答案在票面 §AC#0 那一节（S5 为家、S7 为闸）；
   AC#6 的两格早已移出去到票 220／票 221，本腿没有顺手回收，也没有把「孩子没有停止出口」这格藏起来。
3. **一支需要编排者判的「两支」（本腿按保守那支做完并具名上报，没有自行假设结案）**：
   票面 §0 与本腿派单都写着「流式键今天在两个包里各拼了一份 ⇒ panel 为真相源、由装配根注入函数、不新开 tools→panel 依赖边」，
   而盘面事实是**这一格已由载体层那一程用第三包落掉**：字面量只在 `internal/streamkey/streamkey.go` 一处，
   `internal/panel/pump.go:430` 与 `internal/tools/subagent_197.go:53` 各留一枚 alias，
   外加 `TestSubagentStreamKeyHasOneMintSite` 那枚钉＋`cmd/wisp/subagent_stream_key_197_test.go` 那对比对。
   两支：**(甲)** 承认现形＝裁定精神已满足（真相源只有一处、`tools` 不 import `panel`，那条依赖边确实没开），
   记成「已执行、形状与裁定文字不同」；**(乙)** 要求改回「panel 导出＋装配根注入格式化函数」那一形
   ＝重写已入库且带钉的东西。**本腿按甲执行（那三处一行未动）并在此上报**；
   台账里那句「197-r3 要执行」是否结清，归编排者落账，不归本腿判。
4. **⚠ 一枚可能撞红的预告（不预判、不改判据）**：票 224 的会话令牌那一档若把带 `Allow` 方法的类型挂进
   `tools.Options`／`tools.TaskDeps`，本腿第 (2) 枚反射腿**会红**。那时要裁的是
   「它是判级输入还是答复出口」（`allowed_dirs` 那条教训同族：判级输入不是执行时硬边界，反之也不该被当成出口），
   **不是**放宽这枚尺。
5. 界面那一跳（`PanelSnapshot` 补 `tasks?:` 一枚对象键）与本腿无关：本腿没动 marshal 形状，
   红句今天仍是 `Go Snapshot emits [instructions tasks] that interface PanelSnapshot does not declare` 那一族。
6. 票名里那句「go-side-has-zero-subagent-entity」依旧是历史状态（现量见 §① 那张表）。
