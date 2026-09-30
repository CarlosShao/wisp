# 197-v1 验收件 — 票 197 `AC#5`「子代理永不自批」那枚正控的对抗验收

> 被验收对象：**唯一一格＝`AC#5`**（工单 `.scratch/wisp/issues/197-…subagent-entity.md:44`，
> 判据原文「子代理不自带"允许"出口：与 `Q-49` 丙那批判据同族；正控＝造一枚"子代理自批自己"的假腿 ⇒ 要红」）。
> 被测文件＝`cmd/wisp/subagent_selfapproval_197_test.go`（710 行，产码零改动，commit `a818df46`，写腿 `197-r1`）。
> 编排者 09-30 12:13:34 首跑两枚 PASS——**那不是判语**，本件是本腿（非实现者）自己跑的。
> 本腿代号 `197-v1`，起钟 `2026-09-30 14:48:30 +0800`（`date` 现读），骨架落笔 `2026-09-30 14:50:10 +0800`（同一把尺现读）。
> ⚠ 本文件**只写盘上可核的事实**；每个读数都在下方某处带着它自己的命令。

---

## 0. 逐格 1:1 表（判语；每一格的读数都在 §1-§3 里带着自己的命令）

**总裁＝`AC#5` 成立，附四条具名条件**（四条都不推翻"子代理拿不到允许"这件事本身，见 §0.2）。
⛔ 本腿**没有**碰票面任何 `- [ ]` 框（终态尺：`git status --porcelain .scratch/wisp/issues/` 空，§6）。

| 格 | 工单／普查要求的东西 | 本腿现量落点 | 判语 |
|---|---|---|---|
| **正控**（"造一枚自批假腿 ⇒ 要红"） | 普查 §① 那句"今天没有任何尺会红"的欠账；接缝＝`Native().Allow(childCorr, capturedGrant)` 必须失败 | `:428/:429/:450/:466/:472` 五发花令牌的尝试逐发 `ErrBadGrant`／`ErrUnknownCorrelation`；**M6**（`allowScoped` 不校验令牌）⇒ `:250` 两发"竟然没有被拒"当场红 | ✅ **成立**——普查点名的那枚接缝被逐字执行了，且它有牙 |
| 行为半（每一发都被拒**且理由可核**） | 拒因必须是"令牌不对"，不是"这条路由没在跑" | `:238-258` 六枚用例；**但** D1（把 `ErrBadGrant` 换成 `ErrPanelAllow`）**全绿** | ⚠ **附条件 C2**：六枚共用 `{ErrBadGrant, ErrPanelAllow}` 一枚集合⇒**尺分不清门与权限**（§2.2） |
| 结构半（没有出口可够） | 装配可达性＋`tools.Gate` 方法集恰好两枚 | 反射腿绿；**M3** 指名 `tools.SubagentDeps.Door197: interface tools.answerDoor197 declares Allow()`；**M4** 指名 `tools.Options.Gate: interface tools.Gate declares Allow()`＋方法集变三枚；自带种门正控在场 | ✅ **成立**（两发突变都指名，且 M4 需要 `-overlay`，口径见 §1.0） |
| 令牌不外泄（七枚载具） | 快照／名册（工具侧）／流式正文／**spawn 正文**／stdout／stderr／持久日志 | P1 读数：`stdout=2163 stderr=6879 log=12619 snap=3261 roster=160 chunk=137 `**`spawn=0`**；M1 真形红快照、M2 红 stdout、M2b 红 stderr＋日志 | ⚠ **附条件 C1**：第七枚载具（`rec.spawnText`，`:149` 声明、`:308` 读、**全文无赋值**）恒空＝**宣称了但没测**（§1.2） |
| **正控之正控**：宿主那一发真的落盘 | "如果宿主那一发其实没落盘，'每一发都被拒'就可能是整条路都是死的" | 基线 `宿主允许=<nil>`＋`落盘=true/false`＋审计 `ANSWER-ALLOW`＋**D3**（原生允许改投 `AnswerReject`）⇒ 红**只在** `:277/:280` 正控那两行、**六发拒因一字未变** | ✅ **是真的**，不是死路造成的假绿（§1.1 D3 行、§2.3） |
| 拒绝是谁给的（门 vs 权限） | 本仓栽过 `NoGate` 那枚形状 | 七枚"被拒"逐枚判（§2.1）：**没有一枚**来自 `NoGate`；`waitForChildCard197:131` 那枚 `"L2"` 是承重的（**D4** 换成 `"L1"` ⇒ 红在具名前置读数上）；生产侧三句文案逐字互不相同 | ✅ **本文件里没有那枚形状**；但尺侧的混用成立＝条件 C2 |
| 甲／乙复判 | `internal/streamkey` 那一形满不满足裁定 | §3：编排者那把尺 0 命中（＋本腿加的 `go list -deps` 两向 0）**只量到 R1**；R2（panel 为真相源）与 R3（装配根注入格式化函数）**盘上都没做** | ⚠ **甲＝方向对、理由错**，建议台账追加"改判"式更正（§3.5，归编排者落账） |

### 0.2 四条附条件逐名（都是**尺**的缺陷，不是产码的缺陷）

- **C1（载具虚报一格）**：`subagent_selfapproval_197_test.go` 的 `spawnText` 从未赋值 ⇒ 该格恒绿。
  修＝要么把它接上（`spawn197` 已经握着那份正文），要么把 §WHAT IT MEASURES (4) 与 `197-r1` §② 那两处"spawn 正文"删到实话。
- **C2（门/权限混用）**：`:254` 那枚两值集合改成**逐枚哨兵表**（收紧，不是放宽）。⇒ 新腿跑 **D1** 那发必须红。
- **C3（脚注指错落点）**：文件末尾 `M1` 那行指向 `ui.go` 的 `PanelItem`，实测**改了东西全绿**（它根本进不了快照字节）；
  照它跑的下一腿会误判"这格没牙"。⇒ 脚注要改成"进 `NativeVerdict`/`ApprovalCardView` 那一路"或直接删，并注明**M1 真形**（本腿已给台件）。
- **C4（用例名不符实）**：`:466` 那发不是"跨卡借证"，是"同名重放已花掉的令牌"（两张卡共用 corr）；
  真跨卡那一形**要造两枚同时活着的卡**，配合 D2（`bindDigest` 零尺）＝同一块石头，建议**单开一票**。

⇒ **这五条之外没有可退回项**：AC#5 声称要防的那个结局（孩子自己批自己）在**任何一发突变下都被当场抓住**——
M6 甚至同时抓出了"那张只被面板路线递过令牌的卡真的执行了、文件真的被写出来"这一形（`:284/:291` 两行红）。


---

## 1. M1–M6 逐发读数（前／后＋指名变红的用例名）

### 1.0 口径（先读这段再看表）

- 台件＝`.scratch/wisp/probes/197/v1/mut.py`（本腿写的驱动器）；每发的完整日志在同目录 `<发名>/run-*.txt`。
  驱动器每发做四件事：**落地前断言锚点在那枚文件里恰好命中 1 次** ⇒ 跑 ⇒ **用内存里那份原始字节还原** ⇒
  逐枚打印 `RESTORE <path> bytes=<n> exact=True` 与 `git status --porcelain -- internal cmd`。
  ⛔ 全程没有 `checkout`／`stash`／`reset`／`--amend`／`clean`，没有 commit 任何产码。
- 每发都跑同一把定向尺（含 PATH 两枚，见硬约束 4）：
  `go test ./cmd/wisp/ -run "Test197SubagentSelfApprovalIsRefusedAndTheHostAnswerLands|Test197NoAllowDoorIsReachableFromASubagentsAssembly" -count=1 -v`。
  **每一发都出现了 `--- PASS` 或 `--- FAIL` 行**，所以没有一枚是 `0xc0000135` 那种"根本没跑"。
- 起手基线（本腿自己跑，`2026-09-30 14:50:23` 收尾钟，台件 `baseline-targeted.txt`）：
  **两枚 PASS（11.17s／0.00s，`ok 11.244s`）**——与编排者 12:13:34 那发（11.05s／0.00s）同形，本腿不采信它当凭据。
- 唯一用了 `-overlay` 的那发是 **M4**：具名声明买到哪一半覆盖——
  overlay 只换**编译期**的源，盘上字节一字未动；M4 打的这枚反射腿（`Test197NoAllowDoorIsReachableFromASubagentsAssembly`）
  **不读任何源文件**，所以对这一发 overlay 是 faithful 的。⚠ 但它照不到"读盘的测试"（本包的
  `TestSubagentStreamKeyHasOneMintSite`、`leg_sink_gate_131` 那族走 AST/读盘），那几枚今天**没有**被 M4 覆盖到。

### 1.1 六发脚注突变＋本腿自加的三发：逐发读数

| 发 | 改了什么（落点） | 前 | 后 | 指名变红的用例＋逐字红句（截短） |
|---|---|---|---|---|
| **基线** | 无 | — | **两枚 PASS** | `ok 11.244s`；`落盘=true/false`、六发拒因逐字在 `baseline-targeted.txt:13` |
| **M1**（脚注原文：`ui.go` 给 `PanelItem` 加 `Grant`，并在 `queue.go:524 viewLocked` 填上真令牌） | 3 枚文件 | PASS | ⛔ **全绿，一枚都不红** | **没有任何用例指名**——`--- PASS` 两枚照旧（`M1a/run-cmd_wisp_targeted.txt`）⇒ **脚注这一发指不到它声称的那一格**，见 §1.3 |
| **M1 真形**（本腿补：把同一枚活令牌一路带到**线上包**里——`Queue.Leak197`→`panel_pump.go` 拼进已有的 `reason` 字符串键） | 4 枚文件 | PASS | **红** | `--- FAIL: Test197SubagentSelfApprovalIsRefusedAndTheHostAnswerLands (11.08s)`，红句在 `subagent_selfapproval_197_test.go:319`：**"这一程自己发布的快照字节（卡还挂在队上时取的）里出现了原生一次性令牌本身"** ⇒ **载具扫描本身是活的**，瞎的是脚注指的那枚落点 |
| **M2**（`cmd/wisp/run.go:1235` 让控制台把 `p.Grant` 印出来） | 1 枚文件 | PASS | **红** | 同一枚用例 FAIL，`:319` 两行："**stdout 里出现了原生一次性令牌本身**"（grant1、grant2 各一行） |
| **M2b**（本腿自加：`queue.go:326 grantNonce`（`:336` issue 那一行之后） 增发一枚只走审计漏斗的 `GRANT-DEBUG` 行） | 1 枚文件 | PASS | **红** | `:319` 打了**四**行：stderr 两行＋**"持久日志目录 里出现了原生一次性令牌本身"两行** ⇒ 载具 #6 与 #7 都是活的，`<data>/logs` 那格不是空气 |
| **M3**（`internal/tools/subagent_197.go`：`SubagentDeps` 多一枚字段 `Door197 answerDoor197`，那枚本地接口声明 `Allow`） | 1 枚文件 | PASS | **红，两枚具名** | ① `--- FAIL: Test197NoAllowDoorIsReachableFromASubagentsAssembly`，`:563` 指名 **`tools.SubagentDeps.Door197: interface tools.answerDoor197 declares Allow()`**；② `internal/tools/subagent_197_test.go:880 --- FAIL: Test197SubagentHasNoSelfApprovalOutlet`（"SubagentDeps 字段集合变了：[Roster Door197 …]"）⇒ 脚注承诺的"＋那一枚字段数钉"**兑现** |
| **M4**（`tools.Gate` 接口多出 `Allow(ctx, corr, grant) error`；生产实现 `*approval.Gate` 补上同名方法） | overlay 2 枚 | PASS | **红，两枚具名** | `:563` 指名 **`tools.Options.Gate: interface tools.Gate declares Allow()`** ＋ `:570` "**tools.Gate 的方法集是 [Allow PendingApproval PendingWindow]**，想要 {PendingWindow, PendingApproval}" ⇒ 可达性腿与"方法集恰好两枚"那枚钉都有牙 |
| **M5**（`gate.go:718 DecideFromNative`（`:723` 那一支） 让 `Source=="subagent-self"` 直接放行＝把 `Request.Source` 变成权威） | 1 枚文件 | PASS | **红** | `:250` "**自称 Source=subagent-self 的原生路由：这一发竟然没有被拒**"，`t.Logf` 里那一枚从 `原生令牌无效…` 变成 `<nil>`，其余五发不变 |
| **M6**（`queue.go:376 allowScoped` 不再校验令牌：`spent := true`） | 1 枚文件 | PASS | **红，十条具名** | `:250` 空令牌／借证两发"竟然没有被拒"；`:255` 三发理由退化成 `ErrUnknownCorrelation`；**`:270` 正控塌了**（"宿主侧允许…竟然失败了：本机账上没有这张卡"）；**`:284`＋`:291`：那张"只被面板路线递过令牌"的卡真的执行了、文件真的被写出来了**；`:295` 收口失败；`:341` 审计里找不到 `FORGED-OR-STALE` ⇒ **M6 造出了 AC#5 声称要防的那个结局，而尺当场红**（这是这格里最值钱的一发） |
| **D1**（本腿自加，攻"拒绝是谁给的"：把 `queue.go:380` 令牌不valid 那一支的返回值从 `ErrBadGrant` 换成 `ErrPanelAllow`） | 1 枚文件 | PASS | ⛔ **全绿** | **无一名红**，而 `t.Logf` 里六发拒因**全部**变成"approval: 面板来源不得允许（F2 第三层…）"⇒ 见 §2，这是**尺的精度缺陷**（安全性质本身没被推翻：仍然是拒） |
| **D2**（本腿自加：`approval.go:292 grantStore.spend` 不再比对 binding 摘要，活令牌可花在本机任何一张卡上） | 1 枚文件 | PASS | ⛔ **全绿，两把尺都绿** | 定向尺绿＋**整包 `go test ./internal/agent/approval/ -count=1` 也绿（`ok 0.392s`）**；全仓 `_test.go` 里 `bindDigest`／`不绑定` 只在这枚被测文件里出现过 ⇒ **bind 这一层今天零尺**，见 §2.4 |
| **D3**（本腿自加，攻"正控是不是真的"：`queue.go:382` 原生允许送达的那一枚 `AnswerAllow` 改成 `AnswerReject`） | 1 枚文件 | PASS | **红，红在正控那一格** | `:277` "被允许的写回了错误结果"＋**`:280` "正控失败：…child197-under-host-allow.txt 没有被写出来"**，`落盘=false/false`；**六发拒因一字未变** ⇒ "每一发都被拒"不是"整条路都是死的"造成的假绿 |
| **P1**（临时往被测文件加一行载具长度读数，跑完还原） | 1 枚文件 | PASS | PASS（读数新增） | `CARRIER-LEN stdout=2163 stderr=6879 log=12619 snap=3261 roster=160 chunk=137 `**`spawn=0`** ⇒ 七枚载具里**六枚有字节、一枚恒空**，见 §1.2 |

⛔ 上面**没有一发**需要放宽任何断言、没有一枚 `SKIP`、阈值／golden／`thresholds.go`／`allowlist.txt` 一字节未动。
每发还原后的字节数与起手记录逐枚相等（`run.go 58916`／`queue.go 19660`／`gate.go 31010`／
`approval.go 11277`／`ui.go 7739`／`pending_read.go 5517`／`tools/subagent_197.go 27205`／
`panel_pump.go 18853`／被测文件 `30890`），终态 `git status --porcelain -- internal cmd` 与起手名册相同（**都空**）。

### 1.2 载具扫描：七格里有一格是装饰（`spawn=0`）

`subagent_selfapproval_197_test.go:149` 声明 `spawnText string`，`:308` 把它当"task.spawn 回给父任务模型的正文"扫，
**但整个文件里没有任何一处给它赋值**（`grep -n "spawnText" cmd/wisp/subagent_selfapproval_197_test.go` 只有那两行：
149 声明、308 读）。P1 那一发的读数当场证实：**`spawn=0`**，`strings.Contains("", 令牌)` 永远为假。
⇒ 脚注与 `197-r1` 进度件 §② 那句"载具扫描…`task.spawn` 回给父模型的正文"**是写了但没测**的那一格。
危害判读：**不构成 AC#5 的洞**（同一条正文的另一半——`roster`/`chunk`/快照/stdout/stderr/日志——都在扫且非空；
`task.spawn` 的返回文本今天是测试里的常量而非被测物），但**这格必须具名记为"宣称的第七枚载具不存在"**。
另外两枚读数口径：`roster=160`、`chunk=137` 字节很窄，但两枚都是 `%+v` 整结构转储
（`:423`、`:424`），**新增字段一定会被印出来**，所以它们对"多加一枚会漏令牌的字段"这一形是敏感的（比按 key 过滤的 JSON 扫更宽），
本腿**不**判它们为装饰，只记口径窄。

### 1.3 脚注（M1–M6）与盘面的三处出入（逐条具名，不替它圆）

1. **M1 的落点指不到它声称的那一格**：`Snapshot.Pending` 是 `[]panel.ApprovalCardView`（`internal/panel/composer.go:58`），
   内容由 `cmd/wisp/panel_pump.go:58 liveVerdicts` 从 `Queue.LiveApprovals()` 逐枚**显式字段拷贝**成 `panel.NativeVerdict`（`pump.go:67-81`）再生成卡片；
   `approval.PanelItem` 只服务 `PanelAPI.Head()/View()`（`gate.go:697-699`），**没有任何一条路把它的字段带进快照字节**。
   ⇒ 照脚注那样改 `ui.go` 只会得到"改了东西但全绿"（本腿实测），
   而"面板可见载具带着活令牌"这一形**今天全仓无尺**（扫全仓的读数见 §1.4）。⚠ 这不是被测文件的判据坏了——
   **M1 真形（令牌真的进了线上包）当场红**——是**脚注把下一腿带偏**：照它跑会误判"这格没牙"。
2. **"跨卡借证"那一发其实不是跨卡**：`startChildWrite197`（`:109`）把 `CorrelationID` 写成 `taskID`，
   于是两次 `fs.write` 挂出来的两张卡**共用同一个 correlation id**——**所有打出 `t.Logf` 那行的 13 发读数里，三枚 id 逐次全等**（尺：对 `.scratch/wisp/probes/197/v1/*/run-*.txt` 与 `baseline-targeted.txt` 反查 `child=<id> corr1=<同 id> corr2=<同 id>`，命中 **13**，与含 `child=` 的文件数 **13** 相等 ⇒ 没有一发出现「两张卡不同名」）
   （例：`child=corr1=corr2=45ce09d5-6621-40c9-88ca-858f7b6b29d9`）。所以 `:466` 那一发实际是
   "在同一个名字下重放一枚已花掉的令牌"，**测的是新鲜度，不是跨卡绑定**；真正的"两张同时活着的卡、A 的令牌花到 B 上"
   这一形**本文件今天构造不出来**（加上 D2 的读数，绑定那一层全仓无尺）。
   ⇒ 用例名与注释里的"别张卡的令牌（跨卡借证）"是**名不符实**，要改的是名字（或让它真的造出两枚不同 corr 的卡）。
3. **`M1、M2 会同时打红冻结件 `internal/panel/l2_grant_boundary_test.go``（进度件 §③ 第 4 行那句）本腿复核为不成立**：
   那枚冻结件的射程是**入向**envelope（它找的是本包自己的 JSON **decode** 调用与 `panel.*` 路由，见其文件头 facet 1/2/3），
   M2 是 `cmd/wisp` 往 stdout 印东西，与它无交集；M1 真形本腿刻意走**已有的 `reason` 字符串键**，
   也不需要新增键 ⇒ **两发都不会打红那枚冻结件**。
   （支持这一句的两把现量尺：`grep -n "cmd/wisp\|run\.go\|consoleApprovalUI" internal/panel/l2_grant_boundary_test.go` ＝ **零命中**；
   它的读盘面是本包 AST，不越界到 `cmd/wisp`。）
   反过来说更值得记：**"用已有字符串键 smuggle 一枚令牌"这一形今天没有任何尺看得见**
   （新增 json 键那一形才有 `TestApprovalCardViewJSONKeysMatchFrontendTypes` 那把尺，而那把尺今天**本来就红**，
   在册名册见 §1.5，不能拿来当互证）。

### 1.4 全仓扫那一发（M1a-sweep）

台件：`.scratch/wisp/probes/197/v1/M1a-sweep/run-whole_repository.txt`（`go test ./... -count=1 -timeout 10m`，
跑在 **M1a 那棵树上**＝`approval.PanelItem` 带着**活的原生一次性令牌**）。读数：

- **红 6 枚，逐枚都能归因到与本发无关的东西**：
  - `internal/panel` 4 枚＝§1.5 那四枚在册常红（界面键名册／色表）；
  - `internal/ball` `TestC21TableColourRowsMatchTokensCSS`＝**未突变的基线也红**（本腿 15:1x 单独复跑 `ok`/`FAIL` 逐字：
    `read design/assets/tokens.css: … cannot find the path specified`，成因是**别人在工作树里删了 `design/assets/**`**，
    那枚文件此刻不在盘上）；⛔ `internal/ball` 不是本腿地界，本腿**只记归因、不判它**；
  - `internal/risk` `TestResolvePerCallBudget`＝**全仓并跑时的饥饿**，不是这发改出来的：本腿随后单独复跑
    `go test ./internal/risk/ -count=1` ⇒ **`ok 3.761s`**，同发日志里那一行逐字
    `C26 Resolve: 277143 ns/op = 0.277 ms/op (budget 1.000 ms, 4527 samples)`＝预算 1.000ms，单独跑 0.277ms 通过。
    ⇒ **仪器口径教训（写进 §4 给编排者）：`go test ./...` 里那枚 ns/op 预算尺会被同批并跑的包饿到红，
    凡用全仓扫做突变读数，性能类判据必须单包复跑才能入账。**（同一发里 `cmd/wisp` 从 126.594s 涨到 143.673s 是同一成因。）
- **`cmd/wisp` 在 M1a 下 `ok 143.673s`、`internal/agent/approval` `ok 0.598s`** ⇒
  被验收的那两枚用例**以及 approval 包自己全部 0 枚红**。

⇒ **结论（这一发是本腿自己加的那枚"面板可见读面带着活令牌"的形状）**：
`internal/agent/approval/approval.go:223` 逐字 "The panel-facing surface (PanelItem) has no grant"、
`pending_read.go:9` 逐字 "PanelItem deliberately carries no Params and no grant"、
`cmd/wisp/approval_reply.go:28` 第三次重复同一句——**三条注释、零枚仪器**。
谁哪天给 `PanelItem` 加一枚 grant 字段并在 `viewLocked` 填上真令牌，
**今天全仓没有一把尺会红**（本腿实测：`git grep -c` 见 §2.4 的那条同类读数）。
⇒ 这条**不是** AC#5 判语的否决条件（被测文件没有宣称扫这一面，它扫的是快照/名册/流/stdout/stderr/日志），
但它与票 197 §"规格早就预见过这两个坑"里那句"**面板只能拒绝／查看**"是同一根神经，
本腿把它写成 §4 的一条**建议新开格**（判据很便宜：反射扫 `PanelItem` 字段名不得含 `Grant`，
＋让 `l2_grant_boundary_test.go` 那族把**出向读面**也纳入射程）。


### 1.5 整包复跑（本腿自己跑，不是抄编排者）

`2026-09-30 15:09:07 → 15:11:17`，同一把尺（五枚包、带 PATH 两枚）、**全部突变已还原之后**：

```
ok  	github.com/CarlosShao/wisp/cmd/wisp	         126.594s
ok  	github.com/CarlosShao/wisp/internal/tools	          13.174s
ok  	github.com/CarlosShao/wisp/internal/agent/approval	  0.349s
FAIL	github.com/CarlosShao/wisp/internal/panel	          1.261s   ← 在册四枚，逐名如下
?     github.com/CarlosShao/wisp/internal/streamkey	         [no test files]
```

`internal/panel` 那四枚红逐名＝`TestApprovalCardViewJSONKeysMatchFrontendTypes`／`TestComposerContractTypesMatchFrontend`／
`TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme`／`TestC21DesignTokensFourWayAgree`，
**与 `197-r4` §① 登记的起手在册红同一枚数同一名字**（界面侧键名册与色表，最后一枚的红句逐字
`read design/assets/tokens.css: … cannot find the path specified`）⇒ **相对 `subagent_blocked_197_test.go` 那一族没有多出一枚红**，
本腿那 14 发突变没在树上留下任何东西。⚠ `internal/streamkey` **没有自己的测试**（`[no test files]`），
它的守卫全在 `cmd/wisp` 那枚单一铸造点钉上（§3 用得到这一条）。


---

## 2. "拒绝是谁给的"逐枚判定：门拦下 vs 权限拦下

> 本仓栽过的那一枚形状：一枚 L1 工具穿过 `approval.NoGate`（它的 `PendingWindow` 自己返回 `AnswerReject`），
> 观察到的"被拒绝"是**那道门**给的，不是**权限判定**给的。下面逐枚读这 710 行里每一处"断言被拒"，
> 判它拒的是门还是权限；每行都带**代码落点**与**本腿实测的那一发**。

### 2.1 逐枚判定表

| # | 被测文件那一处断言（行号） | 实际走的路（现读落点） | 拒因是门还是权限 | 本腿用哪一发证明它有牙 |
|---|---|---|---|---|
| 1 | `:428` 空令牌 → `:249` 要红 | `nativeAPI.Allow`（`gate.go:624`）→ `q.allow` → `allowScoped:376 grants.spend("")` → `equalSecret` 对空串先假（`approval.go:269-273`）→ `:380` | **权限**（令牌校验本体） | **M6** 把这一支整体放行 ⇒ `:250` "空令牌…竟然没有被拒" |
| 2 | `:429` 假令牌（真令牌改尾 4 字） | 同一支：`spend` 循环里 `equalSecret(v, nonce)` 不中 → `:380` | **权限** | M6；另：令牌值由运行时算出（`bogusProof197:486`），源码里没有令牌字面量，d22scan ban #3 那条形状没被踩 |
| 3 | `:431` 自称 `Source="subagent-self"` 的原生路由 | `DecideFromNative:718` → `:723 q.allow`（**`Source` 只在 `:720` 的日志里出现，全仓无任何分支读它**：`Request` 的注释逐字 "logged and never consulted"，`approval.go:229-231`） | **权限**（因为 `Grant:""` 也同时不合法，见下一栏） | **M5** 让 `Source=="subagent-self"` 直接放行 ⇒ `:250` 指名这一发"竟然没有被拒"；⚠ **精度缺陷**：这一发传的令牌是空串，所以"来源被忽略"与"这次没带令牌"两件事在**基线**上是同一个读数——要把它做成纯的来源判据，得带一枚**活**令牌自称 subagent-self 并断言它被**允许**（那才是 `bindDigest` 绑事不绑人的原话），本文件没做，也做不了（孩子拿不到活令牌）。本腿判它**够用但不纯**，不是假绿 |
| 4 | `:466` "跨卡借证" | `q.allow(card2.Corr, card1.Grant)` ⇒ 两枚 corr **同名**（`startChildWrite197:109` 把 CorrelationID 写成 taskID）⇒ 实际是"同名重放一枚已花的令牌" → `spend` 找不到 → `:380` | **权限（新鲜度）**，**不是绑定校验** | M6 红；⛔ **D2 全绿**⇒"绑定"那一层（`bindDigest`，`approval.go:253`）在这一格里没有任何读数（见 §2.4）；用例名与注释要改（§1.3.2） |
| 5 | `:469` 面板路线递出真令牌 | `DecideFromPanel:731-747`：`r.Allow` 即**按路由拒**（代码注释逐字 "a grant presented here is refused **even if it is real and live**"），并且**只有**在 `Grant!=""` 时顺带 `q.revokeGrants` | **门**（路由本身，且是设计要求的门：`AGENTS.md` §1.2 铁律／`Q-49 丙`／SPEC-06 §9 layer 3） | 它不需要"权限级"牙——**但它不能替权限那一支作证**；本腿判：这一格算"被拒"成立，**不许**在判语里把它当"令牌校验在工作"的证据 |
| 6 | `:472` 被烧掉的令牌再走原生侧 | `q.allow` → `revokeGrants:430` 已把 `values` 清空 → `spend` 空表 → `:380` | **权限** | **M6** 时这一发退化成 `ErrUnknownCorrelation`（卡已被第 1 发的伪造令牌答掉）⇒ 红；⇒ "烧掉"这一形**是**被第 6 发独立测到的，不是靠第 5 发推断 |
| 7 | `:450` 已结卡的重放 | `allowScoped:365 lookupForAllowLocked` 回 nil → `ErrUnknownCorrelation` | **队列状态**（既不是门也不是权限） | `:261` 只接受 `ErrUnknownCorrelation`，且文件自己在 `:259-260` 写明"卡片离开队列后 byID 里就没有它了"⇒ **诚实标注**，不算混用 |
| 8 | `:283` 第二张卡没执行＋`:290` 文件没落盘 | 直接成因是 `:473 rt.liveCards.h.Reject` → `q.reject:460 deliver(AnswerReject)` | **人（宿主）给的拒绝** | **M6** 那发红过（`:284`＋`:291` 两行）⇒ 有牙；⚠ 但它**不能**被引用成"自批被拒"的证据：那张卡不是被"自批失败"关掉的，是被宿主手动关掉的。措辞建议见 §2.3 |

### 2.2 混用有没有发生：发生在**尺**上，不在生产里

`subagent_selfapproval_197_test.go:254` 那一支把 **两枚不同的拒因当成一个集合**接受：

```go
if !errors.Is(tc.err, approval.ErrBadGrant) && !errors.Is(tc.err, approval.ErrPanelAllow) { … }
```

⇒ 本腿写了 **D1** 专门攻这一处：把 `queue.go:380` 令牌校验失败返回的 `ErrBadGrant`
原地换成 `ErrPanelAllow`（**门**的句子，挂在**原生**路由上）。读数：**两枚用例全绿**，
而 `t.Logf` 里那一行六发拒因**全部**变成 `approval: 面板来源不得允许（F2 第三层：允许只接受原生侧）`。

- **安全性质没有被推翻**：那一发仍然是"拒"，落盘仍然是 `true/false`，正控仍然绿。
- **被推翻的是这一格的措辞**：`:236` 那句"每一发的理由都是「令牌不对」，不是「这条路由没在跑」"
  在这把尺下**不成立**——尺无法区分"令牌不对"与"路由不给允许"。同一次改动里审计行
  `FORGED-OR-STALE`（`queue.go:379`）照旧在场，**审计说的是令牌、返回值说的是路由**，
  两者从此互相矛盾，而这枚用例**看不见**矛盾。
⇒ 判词：**门拦下的文案不许复用权限拦下的文案**这一条，**生产侧成立**（三条句子逐字互不相同：
`approval: 原生令牌无效（缺失/已用/与本次请求不绑定）`／
`approval: 面板来源不得允许（F2 第三层：允许只接受原生侧）`＋包装句
`面板来源的「允许」被服务端 API 直接拒绝（F2 第三层），可改为拒绝或查看完整参数`／
`L2 审批通道尚未接入（票 21），已拒绝执行`＝`internal/tools/gate.go:144` 的 `NoGate`），
**尺侧不成立**（六枚用例共用一个两枚集合）。修法很小且必须写清：那一圈 `for` 里给每枚用例配上**它自己那枚**哨兵
（1/2/3/4/6 要 `ErrBadGrant`，5 要 `ErrPanelAllow`，7 已经单独钉住），
并把这一格从"两枚集合"改成"逐枚常量表"——**这是收紧，不是放宽**。

### 2.3 `NoGate` 那一形在本文件里有没有复现：**没有**

`tools.Gate` 的 `NoGate` 缺省（`internal/tools/gate.go:135-145`）确实会"自己返回 `AnswerReject`"，
但本文件测的**不是工具侧看到的答复**，而是 `*approval.Gate` 的**答复入口返回值**
（`Native()`／`DecideFromNative`／`DecideFromPanel`），而这三条路都进 `allowScoped`。
两枚结构性防呆本腿核过：
- `waitForChildCard197:131` 要求 `c.Level == "L2"`（文件自己在 `:124-126` 写明"an L1 window is booked into the same ledger with an empty Grant"）；
- `Replies.Allow`（正控那一发）在到令牌校验之前要过三道**门**：`ErrNoTrackedCard`／`ErrRouteHasNoAllow`／`ErrNoGateAttached`
  （`replies.go:318-326`）——**正控绿 ⇒ 三道门当场全开**，所以六发拒因不可能是"门没接"给的。
（本腿另外跑了一发把 `:131` 那枚 L2 换成 L1 的突变，读数见 §2.5。）

### 2.4 `bindDigest` 那一层今天零尺（D2，跨条发现）

`grantStore.spend` 的绑定比对（`approval.go:292-306`，比较 `equalSecret(stored, bind)`）
本腿整枚删掉后（**任何一枚活的未花令牌都能花在本机任意一张卡上**）：
- 定向尺：**两枚 PASS**；
- **整包 `go test ./internal/agent/approval/ -count=1`：`ok 0.392s`（0 枚红）**；
- 全仓 `_test.go` 里能指到这一层的词（`bindDigest`／`不绑定`／`grants.spend`）**只出现在被验收的这枚文件里**。

⇒ **`internal/agent/approval/queue.go:44-47` 那句"bind 覆盖 corr/task/tool/level/args/seq，所以一枚 nonce 不能被重放到别的请求上"
今天是一条没有仪器的断言**。⚠ 射程要说准：这不是 AC#5 的洞（孩子仍然拿不到**活的未花**令牌，
AC#5 的防线是"拿不到"），这是**纵深防御的那一层无人看守**＋**本文件的头部注释（`:14-15`、`:320`）恰好把 `bindDigest` 当成本格成立的前提来引用**。
建议具名新开一票（形状像票 220/221 那种拆法），判据要造**两枚同时活着的卡**——
这正是 §1.3.2 那枚"同名两张卡"缺的东西，两笔账是同一块石头。

### 2.5 L2 那枚承重防呆的读数（D4）

把 `waitForChildCard197:131` 的 `c.Level == "L2"` 改成 `"L1"`（＝故意让孩子可能落在"倒计时窗口"那枚门上）：

```
subagent_selfapproval_197_test.go:228: the self-approval probe never completed its setup:
  孩子在 78cdf3ab-… 名下始终没有挂出一张 L2 卡（本用例的全部判据都挂在它上面）；
  这一程已经显示了 1 枚确认卡
--- FAIL: Test197SubagentSelfApprovalIsRefusedAndTheHostAnswerLands (31.12s)
```

⇒ 这一发**当场红在判据本体上**（`t.Fatalf` 走 `:228` 那条具名前置读数），
而且红句自己说出"已经显示了 1 枚确认卡"——**不会有人把它读成"没有卡"**。
所以 §2.3 的结构性判断拿到了一发实测支持：本文件不会因为一枚 L1 倒计时窗口的"门给的拒绝"而签绿。

### 2.6 本节判语

- 七枚"被拒"断言逐枚读过：**没有一枚是 `NoGate` 那种"门自己返回 AnswerReject"造成的假绿**，
  正控（`:438` 宿主 `Replies.Allow`＋`:440` `os.Stat`＋审计 `ANSWER-ALLOW`＋D3/M6 两发）把"路是活的"钉住了。
- **但是**：六枚用例共用一枚 `{ErrBadGrant, ErrPanelAllow}` 集合，**D1 实测证明这把尺分不清门与权限**；
  第 5 枚（面板路线）本来就是**门**给的，不能拿它给"令牌校验"作证——替它作证的是第 6 枚。
- ⇒ 修法两条，都归写手腿，本腿不动被测文件：①`:238-258` 那一圈改成**逐枚哨兵表**；
  ②`:283-285` 那两行的归因文案改成实话（那张卡是被宿主 `Reject` 关掉的，不是被"自批失败"关掉的）。



---

## 3. 甲／乙裁定复判（编排者裁了甲，本腿判裁得对不对）

### 3.0 出处先钉正（这一处腿顶得对，本腿复核支持它）

权威文本**只在台账**：`docs/reports/pending-and-issues.md:8786-8791`（票 197 那一段），
`docs/specs/**` 与 `.scratch/wisp/issues/README.md` 里 `SubagentStreamKey`／这一段裁定 **各 0 命中**
（`grep -rn "SubagentStreamKey|真相源.*panel|装配根注入" docs/specs/`＝0，`grep -n "SubagentStreamKey|streamkey" .scratch/wisp/issues/README.md`＝0）。
⇒ 派单里"权威文本在台账 `:8786-8792`"这一句**成立**，而"引到 specs 或 issues/README"是引错出处；
写腿当场顶正这件事，本腿判**正确**（票面 `:185-188` 那三行是它按同一份文本复述的，不是第二处权威）。

### 3.1 编排者那把尺，本腿独立复跑（逐字）

```
$ grep -rl '"github.com/CarlosShao/wisp/internal/panel"' internal/tools/*.go
（无输出，rc=1；去掉 _test.go 后 0 枚）
$ grep -rl '"github.com/CarlosShao/wisp/internal/tools"' internal/panel/*.go
0 枚
```

本腿再加一把**更严**的尺（`grep import` 只看直接边，传递边它看不见）：

```
$ go list -deps ./internal/tools | grep -c wisp/internal/panel   → 0
$ go list -deps ./internal/panel  | grep -c wisp/internal/tools   → 0
$ go list -deps ./internal/streamkey | grep wisp/                 → 只有它自己（叶子包，除 strings 之外不 import 任何东西）
```

⇒ **"不新开 tools→panel 依赖边"这一条：成立，两把尺都成立**，且 `internal/streamkey` 是干净的叶子，
不构成"绕一层把 panel 拖进 tools"的暗边。这把握得住。

### 3.2 但裁定原文里有**三条**要求，盘上只满足一条

| 那一段裁定要求的（逐字出处） | 盘上现形 | 满足？ |
|---|---|---|
| **R1**「**不新开 tools→panel 的依赖边**」（`:8791`） | 两包互不 import，两向 `go list -deps` 都 0 | ✅ **满足**（§3.1） |
| **R2**「**裁定＝`internal/panel` 那一枚是真相源**（流键形状属于"流的载体"，谁装袋子谁定名）」（`:8789`） | 真相源是**第三包** `internal/streamkey/streamkey.go:29`（`const SubagentPrefix = "subagent:"`，全仓非测试码里唯一那枚字面量，本腿 `grep -rn '"subagent:"' --include=*.go internal cmd \| grep -v _test.go`＝只有这一处命中＋一枚注释）；`internal/panel/pump.go:430,441` 与 `internal/tools/subagent_197.go:53,122` **都是别名、平级** | ❌ **不满足**（panel 不再"定名"，它只是照着念） |
| **R3**「`internal/tools` 那份**改成由装配根注入**（`TaskDeps` 带一枚格式化函数，`cmd/wisp/run.go` 里把 `panel.SubagentStreamKey` 递进去）」（`:8790`） | `TaskDeps`（`internal/tools/task.go:35`）的字段只有 `Roster`／`Paths`／`Spill*Tokens` 三枚计数——**没有任何格式化函数**；tools 侧靠**直接静态 import** `internal/streamkey` 拿到拼法；装配根唯一那一次调用是**读侧**的 `cmd/wisp/panel_pump.go:162 key = panel.SubagentStreamKey(rec.TaskID)` | ❌ **不满足**（注入那一跳整支没做，也不该按原样做，见 §3.4） |
| 裁定自带的那枚判据「**同一段文本在两处拼同一枚键＝红**」（`:8792`） | `cmd/wisp/subagent_carrier_197_test.go:431 TestSubagentStreamKeyHasOneMintSite`：要求非测试字面量落点**恰好 1 枚**、且**必须是 `internal/streamkey/streamkey.go`**，再把两包的 `SubagentStreamKeyPrefix` 与两枚 id（含一枚中文 id）的铸造结果逐枚比对 | ✅ **满足，且比裁定那一形更强**（本腿未跑它，它在 §1.5 的 `cmd/wisp ok 126.594s` 里是绿的；⚠ 它用 AST **读盘**，所以对 `-overlay` 失明——M4 那一发没有覆盖它） |

### 3.3 「形状满足」与「裁定要求」是不是一回事——**直说：不是**

1. **尺的射程只盖住 R1**。拿一把只量得出 R1 的尺去宣布"**满足这条裁定**"，
   等于**用仪器照不到的那两支签了字**。这正是 `AGENTS.md` §1.2 那段"规格与仪器的射程不是一回事"的形状，
   只不过这次是反的：**裁定比尺要求更多，尺更窄**。本腿在 §1 里给同一形状记了三发读数（M1a／D1／D2 都是"改了东西但全绿"）。
2. **R2 被"移交"了，这不是"形状不同、精神已满足"的同义反复**。裁定的理由是"谁装袋子谁定名"——现形里
   **袋子（panel）不再定名**，名字住在第三家。判"这样更好"完全成立（§3.4），但那是**改判**，不是"满足"。
3. **`:8788` 那句"装配根 `cmd/wisp` 是唯一的接缝"是当时的现状描述，不是要求**。
   它在那一段里的作用是**前提**（"正因为两包互不 import，只有装配根能比对它们 ⇒ 才看不见漂移"），
   而 `internal/streamkey` 落地之后**这句话作为事实已经不成立**（接缝移到了 streamkey，比对不再需要装配根同时 import 两包）。
   把一句**已被现形作废的描述**当作裁定的内容去宣布"满足"，是本腿要纠的第二处。
   ⚠ 且这个推广动作**编排者自己在同一本台账 `:9761` 已经立过规矩**：把这一句推广到别的依赖边裁定时
   "**推广动作本身在此具名，不写成'仓里本来就有这条规矩'**"。`:9776` 那一行没有具名推广，
   它把推广后的句子直接写成了"这条裁定"。

### 3.4 甲这一**决定**本身：对，而且比裁定那一形更强

本腿独立核过的三条理由：
- **漂移被机械关掉了**：单一铸造点那枚钉（§3.2 第四行）比"由装配根注入"更硬——
  注入函数**可以传 nil**（现形里 `TaskDeps` 若多一枚 `func(string) string`，装配根漏传就退化成空键或包内 fallback），
  一枚"必须恰好 1 个落点且必须是它"的扫描不会。
- **两包不再互抄，同时也没把方向装反**：`go list -deps` 两向 0（§3.1），裁定担心的"视图层变成工具层的下家"没有发生。
- **裁定没预见过的那半也被现形顺手关掉了**：`internal/streamkey/streamkey.go:32-40` 逐字记录
  "空/纯空白 id 那一支**过去只在读侧有、写侧没有**（＝台账 A406 指的那处分歧），别名化一次把两侧同关"。
  ⇒ 现形做到的比裁定的两支**多**，不只是形状不同。

### 3.5 复判结论（一句话）

**裁甲＝方向对、理由错**：那一形**满足 R1、满足裁定要求的判据、并额外关掉空 id 分歧**，
但**不满足 R2（panel 为真相源）与 R3（装配根注入格式化函数）**，而编排者用来宣布"满足这条裁定"的那把尺
**只量得到 R1**。⇒ 本腿建议台账**追加**一条更正（append-only，不改 `:9776` 原句），
把"满足"改成「**以第三包 `internal/streamkey` 改判：R1 满足；R2、R3 两支作废，作废理由＝现形更强（附单一铸造点那枚钉＋空 id 同关）**」，
并把「装配根是唯一的接缝」那一句标成**已被现形作废的历史描述**、后续引用须具名推广。
**要不要落这一条更正归编排者**，本腿不动台账。


---

## 4. 没做完／留给编排者

⛔ 不留半成品：这一节在骨架提交（`321110da`）时就已写满，以下是**交件时刻**的版本（欠账变了就改它，不删旧事实）。

### 4.1 等编排者做的台账／派单动作（本腿无权落 `A##`／`Q##`／`R##`，也不动票面框）

1. **`AC#5` 那一格要不要翻勾**：本腿总裁＝**成立，附四条具名条件**（§0.2）。
   翻勾文案请引本件 §1.1 的逐发读数（尤其 **M6** 与 **D3** 两发），**不要引 `197-r1` 的自述**——
   它 §② 那句"载具扫描…`task.spawn` 回给父模型的正文"经 P1 实测**是虚的**（C1），照它写会把一格判语建立在一枚不存在的读数上。
2. **写腿腿（下一程）的三枚修**，都是尺的修法、不是产码的修：
   ① 被测文件末尾脚注 `M1` 那行改落点或删除（C3）；
   ② `:254` 那枚两值集合改**逐枚哨兵表**（C2），改完必须让 **D1** 那发红；
   ③ `:466` 那一发的名字／注释从"跨卡借证"改成实话（C4）。
   ⚠ 这三处若由写腿动被测文件，**必须同时把 P1/M1b/D1 三发当正控跑一遍**（本腿台件可直接复用）。
3. **建议新开一票（本腿不塞回票 197）**，判据形状两条、同一块石头：
   - 造**两枚同时活着、不同 correlation id** 的卡，把 A 的**活**令牌花到 B 上 ⇒ 要红（今天 D2 全绿，`bindDigest` 那层零尺）；
   - `approval.PanelItem`（`PanelAPI.Head()/View()` 的出向读面）**不得带 grant**——
     `approval.go:223`／`pending_read.go:9`／`cmd/wisp/approval_reply.go:28` **三处注释都这么写，全仓零仪器**（M1a 实测，§1.4）。
     最便宜的判据＝反射扫 `PanelItem` 字段名＋把 `l2_grant_boundary_test.go` 那族的射程从"入向 envelope"扩到"出向读面"。
4. **§3 的甲／乙复判要落一次台账更正**（append-only，追加，不改 `:9776` 原句）：
   建议措辞见 §3.5；同时把 `:8788` 那句"装配根 `cmd/wisp` 是唯一的接缝"标成**已被现形作废的历史描述**，
   并提醒后续引用它做依赖边裁定时按 `:9761` 自己立的规矩**具名写成推广**。
5. **`197-r1` 进度件 §⑥.3 那句"本腿按甲执行并上报"** 本腿判它**上报姿势正确**（它没有自行结案、没有动那三处、
   也没有把"197-r3 要执行"当已完成）——这条是给它记的一枚正面事实，别在更正时顺手抹掉。

### 4.2 仪器口径教训（写出来给下一腿，不属判语）

- ⚠ **`go test ./...` 里的 ns/op 预算尺会被同批并跑的包饿**：M1a-sweep 那发里 `internal/risk TestResolvePerCallBudget` 红，
  本腿随后**单独**复跑同包 ⇒ `ok 3.761s`（同发日志逐字 `277143 ns/op = 0.277 ms/op (budget 1.000 ms, 4527 samples)`）。
  ⇒ 凡用全仓扫做突变读数，性能类判据**必须单包复跑才入账**；同一发里 `cmd/wisp` 从 126.594s 涨到 143.673s 是同一成因。
- ⚠ **M4 这一形做不成"单点真文件突变"**：`tools.Gate` 是接口，加方法会连带要求 `internal/perm/ticket90_persist_test.go`
  里那枚 `t90gate` 实现它——**那是冻结件，本腿不能动**。本腿改用 `-overlay` 买编译期替换（买到什么／买不到什么已写在 §1.0）。
- ⚠ 在册常红引用时要带归因：`internal/panel` 4 枚是登记的界面键名册／色表常红；
  `internal/ball TestC21TableColourRowsMatchTokensCSS` 的成因是**别人在工作树里删了 `design/assets/tokens.css`**
  （未突变的基线同样红），⛔ `internal/ball`／`internal/audio` 不是本腿地界，本腿**只记归因、不判、不改**。
- ⚠ 本腿**没有**跑过未突变的 `go test ./...` 全仓基线（只在 M1a 下跑过一次＋单独复跑 risk/ball）⇒
  全仓基线名册请引 `197-r4` §① 与 `241-v1` 的在册读数，别当本腿的读数用。

### 4.3 本腿明确没做的事

- 没动票面 6 枚未勾框（AC#0/#1/#2/#3/#6 不由本腿裁）、没动任何产码、没动三枚冻结件、没动阈值／golden／`allowlist.txt`。
- 没读 `frontend/**`／`design/**`（连结论都没引；§1.4 那句 `tokens.css` 缺失是**测试红句里的路径名**，不是本腿去读了那两层）。
- 没 push；只 commit；每枚 commit 都带显式 pathspec（`git commit --only <那一枚文件>`）。


---

## 5. 本腿攻不动的地方

判不动就写判不动。每一格都注明**为什么**攻不动，以及本腿用别的方式旁证到了哪里。

1. **"模型真的发出 `tool_call`"那一支攻不动**：`mockllm` 只会回显、脚本不出 tool_call，
   所以孩子的调用是从装配好的真桥以孩子身份发出的（文件 `:39-43` 自己声明，本腿同意那是
   **mockllm 的能力问题、不是 AC#5 的格**）。⇒ 本腿**无法**判"模型侧是否也可能自批"。
   旁证到的部分：`internal/tools/subagent_197.go:526` 那条 sink 与 `task.spawn` 那道深度检查走同一枚桥，
   所以"身份"这一半是真的；"谁按的按钮"那一半不是。
2. **真·跨卡绑定那一形造不出来**（不是不想造）：要把 A 的**活**令牌花到 B 上，需要
   **两枚同时 pending、不同 correlation id** 的卡；本文件两次 `fs.write` 因 `:109` 把
   `CorrelationID` 写成 taskID 而**共用一个名字**（13 发读数逐字如此）。
   改那一枚常量能得到两枚不同 corr，但要用它做绑定判据就得**重排整段用例**
   （先挂两枚都不答、再互相花令牌、再各带自己的正控）——那是**重写这枚用例的一半**，越过验收腿的射程。
   ⇒ 判"造不出"，判据形状写进 §4.1.3 的新票；D2 那发只证明"今天没有任何东西看守绑定"。
3. **反射腿的动态类型盲区**：`findAllowDoors197` 读的是**静态类型**（文件 `:44-49` 自己声明）。
   "某枚字段的静态类型是 ask-only 接口、真正塞进去的具体类型带 `Allow`"这一形**这枚尺结构性照不到**；
   M4 只覆盖"接口**声明**里出现答复动词"那一半。⇒ 另一半要另一种判据
   （装配完成后对 `opts.Gate` 的 `reflect.TypeOf` 具体类型再扫一遍），本腿**不做、也不替它签绿**。
4. **`sha`-to-ledger 那一跳**：本用例**故意不**调 `publishPanelSnapshot`（与 197-r4 同一纪律），
   "挂卡那一枚快照字节 == 落账那枚 sha"在本腿可见范围外。本腿只核到替代的两条实话：
   `snap.Pending` 里认得出这张卡（`snapClearedCar`，基线绿）**且**这一程至少落账过一枚快照
   （`:356` 那条 `ledgerSummaries145` 非空，基线绿）。
5. **界面那两层攻不动**：`frontend/**`／`design/**` 两层禁（不读、结论不引）。
   工单 `:62-66` 那枚 `tasks?:` 键转绿的一跳归界面那支；本腿唯一与它们沾边的一条读数是一个**测试红句里的路径名**
   （`design/assets/tokens.css` 不存在），那不是本腿去读了那两层。
6. **L1 那一半"被拒是谁给的"不归本腿判**：票面 §09-29 与 `A416` 记载 L1 短窗口在名册上恒读 `false`
   （`LiveApprovals` 只遍历 `q.pending`），那是票 220 的地界。本腿只把**这枚文件里**的 L2 承重防呆验了（D4），
   没有替票 220 判任何东西。
7. **两枚仪器坑本腿是绕过去的、不是攻不破**：每发都带 sherpa＋build 两枚 PATH，
   14 发读数**每一发都出现 `--- PASS` 或 `--- FAIL` 行**（⇒ 没有一枚是 `0xc0000135` 那种"根本没跑"）；
   最长一发 31.12s（D4），**没有一枚接近 300.0x s**，所以本件里没有 C18 审批超时混进来的红。
8. **注释级地雷：同一枚过期指认不止那一行（派单第 9 条要的"别的过期断言"）——本腿一行未改，具名上交**
   本腿自己确认 `07-ball-state-machine-core-done.md` **在盘上＝票 07 早已结案**，然后把
   "球／托盘图标／GUI 归票 07"这一类指认在 Go 侧逐枚抽了出来：

   | 落点 | 逐字片段 | 本腿判读 |
   |---|---|---|
   | `cmd/wisp/resident_windows.go:81` | `the floating ball arrives in ticket 07` | **已知过期**（已登记为票 228 AC#7），本腿不改 |
   | `cmd/wisp/main.go:25` | `wisp GUI resident process (boots the runtime skeleton, empty event loop; the floating ball window is ticket 07)` | ⚠ **同一类，而且它不是注释——是 `wisp -h` 打给用户的 usage 正文**；后果比 :81 重（过期指认会被 owner 当场读到） |
   | `cmd/wisp/main.go:8` | `The floating ball GUI is ticket 07.` | 同一类（Ticket 03 scope 注释里的将来时归属） |
   | `cmd/wisp/notify_windows.go:12` | `the resident process's own tray icon is ticket 07/62's, in internal/ball` | 同一类，且**同时**指向票 07 与票 62（票 62 也早已结案） |
   | `internal/proc/boot_windows.go:127` | `ball bring-to-front lands with ticket 07` | 同一类；⚠ `internal/proc` **不在本腿写面**，本腿只读了一枚函数体确认这行真会打到日志里 |
   | `cmd/wisp/console_other.go:7`／`console_windows.go:26` | `links as a GUI-subsystem binary (ticket 07)`／`the windowsgui subsystem (the final GUI build, ticket 07)` | **未判**（这两句像历史归属而非将来承诺；要判得先查那次链接改动出自哪枚票，本腿没有据此下结论） |
   | `cmd/balldebug/main.go:3`／`:590` | `the ticket 07 debug harness` | **未判**（balldebug 是已交付台件，归属看起来成立；那一带属球的地界，本腿不归因） |

   ⇒ **给编排者的一句话**：票 228 AC#7 现在的完成判据只钉 `resident_windows.go:81` **一行**，
   而同一类过期归属在 `cmd/wisp`／`internal/proc` 里**至少还有四行**（其中一行是**打给用户看的正文**）。
   照现有判据改完，这条链会在下一个文件里继续传——**这正是 `240-c1 → 241-r1` 那三程传递的形状，只是这次不是一行是一族**。
   ⛔ 本腿一行未动（既没改 `cmd/wisp` 里那几处，也没碰 `internal/proc`）。

---

## 6. 终态自证（交件时刻现跑）

钟：`2026-09-30 15:19:39 +0800` 落笔上一节，下面这组读数在 **2026-09-30 15:24:49 +0800** 前后全部现跑（同一把 `date` 现读），全部**在 14 发突变都已还原之后**：

| 尺 | 读数 |
|---|---|
| `go build ./...` | **rc=0** |
| `/d/work/base/gopath/bin/gofumpt.exe -l internal/ cmd/` | **输出空**（rc=0）＝没有一枚文件待格式化 |
| `./tools/d22scan/d22scan.exe -root .` | **`clean - no D22 ban violations`，rc=0**；口径提示：`ban #8 internal/=473`、`cmd/=65` 是**被扫文件数**（含 `_test.go`），不是违规数；那句 `skipped as git-ignored: 1 file(s) under frontend/dist/assets/` 是既有 gitignore 形状，不是本腿留的 |
| **`git status --porcelain -- internal cmd`** | **空＝与本腿起手名册相同**（起手也是空；⛔ 本件不写"必须为空"这种绝对式，只写"等于起手名册"） |
| `git status --porcelain .scratch/wisp/issues/` | **空**＝票面一枚 `- [ ]` 都没碰（翻勾归编排者） |
| `git status --porcelain docs/PLAN.md docs/specs tools/d22scan/allowlist.txt internal/observe/thresholds.go internal/panel/tokens_fourway_test.go internal/panel/l2_grant_boundary_test.go internal/perm/ticket90_persist_test.go` | **空**＝冻结件与阈值一字节未动 |
| 每发突变后的还原自证 | 14 发里**动过盘上文件的那 13 发**全部逐枚打印 `RESTORE <path> bytes=<n> exact=True`（M4 走 `-overlay`，盘上文件自始至终没被写过，因此它没有 RESTORE 行——这正是选它的原因），且每发末尾 `git status --porcelain -- internal cmd` 回 `''`（驱动器 `mut.py` 的 finally 分支，跑失败也还原） |
| 本腿自己的一发废读数（具名，不藏） | M4 **第一次**跑出来是 `[setup failed]`，红句逐字 `found packages tools (bridge.go) and approval (gate.go) in ...internal\tools`＝**本腿驱动器自己的 bug**（overlay 的虚拟文件用 basename 命名，两枚同名 `gate.go` 撞成一枚），不是被测物的读数；修完（加序号前缀）复跑才是 §1.1 表里那一发 M4 |
| 别人在工作树里的脏文件 | `.gitignore`／`design/**` 那 16 枚删除／`docs/evidence/s1/152-*`／`.scratch/wisp/probes/**` **原位未动、未提交**（本腿只往 `probes/197/v1/` 里**新建**台件，只建不删） |

本腿自己的 commit 名册（**全部只碰这一枚证据件，全部带显式 pathspec，全部未 push**）：
`321110da`（骨架，§4/§5 起手即满）→ `ccc1b1cc`（§1-§3）→ 本节所在的那一发。

**终态三跑**（"绿"这个字本腿只在跑过整包与还原后的定向各一遍之后才写）：

| 时刻 | 跑法 | 读数 |
|---|---|---|
| 14:50:23 | 定向 `-run` 两枚 `-count=1 -v`（起手，未突变） | **两枚 PASS**（11.17s／0.00s，`ok 11.244s`）＝`baseline-targeted.txt` |
| 15:09→15:11 | 整包五枚包 `-count=1`（14 发突变里的 12 发已还原之后） | `cmd/wisp ok 126.594s`（0 枚红）＝§1.5 |
| 15:26:12→15:26:25 | 定向两枚 `-count=1 -v`（**全部 14 发还原、本腿三枚 commit 落完之后**） | **两枚 PASS**（11.19s／0.00s，`ok 11.251s`）＝`final-targeted.txt`；跑完 `git status --porcelain -- internal cmd` 仍等于起手名册（空） |


