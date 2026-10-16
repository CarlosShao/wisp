# 165 — 「先出方案、你点头再动手」那一档：两家都有、都带批准门，我们表上**零枚**（**契约级·立而不派**）

- Status: **立而不派**（新增一档＝动 `D31`/`D43`/`C12` 那条状态与审批线，**人工批准**；且 owner 已定过"改状态机要谨慎"）
- Status 更正（09-27 18:4x 编排者现量，**只追加不抹原句**）：上面那句"新增一档＝动 `D31`/`D43`/`C12` 那条状态与审批线，人工批准"要按 `A323` 落地的那一行**收窄**：owner 批准并落进冻结文本的是 **`plan.present` 这枚工具名（`PLAN.md:2566`，L2；`SPEC-07:73` 镜像；提交 `824bb9c`）**，**不是**新增一档审批状态、也**不是**改 `D43` 转移表 ⇒ **`D43`／`C12` 那两条线对本票仍然冻结**（碰它＝停手上报）。⇒ 本票可做的是"把 `plan.present` 这枚工具按已批准的甲路实现出来"，**不可做**的是"顺手把状态机扩一档"。
- 来源：09-26 对标——dsh 有 plan-mode 且**带批准门**（`packages/plan/plan-mode/src/index.ts:278`）、Step-Code 也有（`src/step/plan-mode-tools.ts:117,136`）；我们 `D34` 抽出来的名册里**没有任何 plan/propose 那一形**。台账 `A308`。
- 关联：`D43`/`C12`（冻结状态机，**本票默认不动它**）· `D31`/`D45`（并发与批量授权）· 票 162/163/164（同一批能力面）· **R20**（三档权限模型，已定案）

## 这是什么（人话）

用户说一句模糊的话（"帮我把这个项目收拾一下"），今天模型可以直接开始改文件。
要补的是中间那一档：**先只读、先出一份要干什么的清单、人点头之后才允许写**。
它对我们尤其重要，因为**入口是语音**——听错了没有回车键可撤回，**这一档就是那个刹车**。

## 为什么值得做／不做的后果

不做：我们所有"写"类风险（`fs.write`/`fs.trash`/`fs.move`、以后的 `shell.exec`）都只能**逐次弹卡**，一次一句、500 次就是半小时（`PLAN.md:2118` 早算过这笔账）。有了"先方案"这一档，**一批改动只问一次**，而且问的是意图不是每个字节。

## 硬约束

1. **不许新造一套状态机**。`D43`/`C12` 那张表是冻结的；本票**优先复用现有的 `AwaitingApproval` 那一族**，若量到"必须新增状态才成立"，**停手上报**（那是契约级，不是实现细节）。
2. **不许把"方案被批准"变成"后续所有写都被批准"**——批准只批这一次/这一批（与 `Q-49` 已批方向、以及三家"一次性放行"一致）。
3. **不许动阈值/golden/`allowlist.txt`**；不许新增 `t.Skip`。

## AC

- [x] **AC#1 未修码读数**：今天"只读地出一份方案"这一发能不能做到？（做不到＝本票成立；做得到＝我判断错，直接报回）⚠ 不许把"应该做不到"写成交件。
- [x] **AC#2 复用形状**：方案态**不新增状态**（拿 `--name-only` 名册自证没动 `internal/risk/**` 与状态机文件）；若必须新增，**本票退回**、转成一条待人拍板的 `Q##`。
- [x] **AC#3 批准只批这一批**：拿两发变异证——①批完方案后再来一发**没在方案里**的写，必须响；②摘掉"批过的范围"这一味，是否存在一发越界从此打不红。
- [x] **AC#4 面板不许代答**：⚠ 绝对禁令在案（面板侧来源的 L2"允许"不算允许）。本票只加"发起请求"的出口，**不许**加"给出允许"的出口。
- [x] **AC#5 契约轴**：`docs/PLAN.md`、`docs/specs/**`、`internal/risk/**`、`internal/panel/**`、`internal/agent/approval/**`、`thresholds.go`、golden、`allowlist.txt`、`frontend/**`、`design/**` 零字节。

## 本票**不**解决

- 不做界面（谁显示方案、点哪儿算点头＝`frontend/**`，owner 委托会话的写面）。
- 不做"计划/日程"那个产品功能（`D12` 那条线，`GAP-25` 重提已被拒）。

## 派单前置条件

owner 一句话（要不要这一档）；票 161 结案（先有"有没有人按门铃"那把尺，否则这一档也会变成建好没人喂的门）。

- [2026-09-26 23:3x] 编排者代落：只读普查程 `16x-c1` 交回本票的**冻结文本射程**七问（全文见 `docs/evidence/s1/16x-contract-lines-c1.md`，§3 那张表里是 owner 一句可批的话）。**六枚里唯一一枚「票面把自己锁死」**：它可能的三条路（`docs/PLAN.md` 的 D34／`internal/risk/mode.go:65`／`internal/agent/approval/batch.go:15`）全在自己票面 `:29` 的零字节名单里 ⇒ 不动冻结面它一格也落不了。另附带量到：它复用的 `AwaitingApproval` 那一族有一格行为**至今未定义**（`PLAN.md:1530`，也是 `AGENTS.md §2` 的具名停手项）。

- [2026-10-15 16:27:33 +0800] agent=165-r1 did=落件 plan.present（甲路，写面 internal/tools/**）：`internal/tools/plan_present.go`（新）＋`plan_present_165_test.go`（五格用例）＋探针 `probes/165/`（七枚）＝提交 `4a49480d`（起手锚 b74b87d0 → 落件时锚 9b0f7d55，漂移已记）。五格读数：AC#1 未修码名册 10 枚、零枚 plan 形（票成立）；AC#2 复用形状（Declared=L2 → R1 → 现成 Gate.PendingApproval，即 D43 #17/#23/#24/#25；Gate 仍 2 回调、GrantSource 仍 1 只读方法，零新增状态/通道/模式）；AC#3 清单批过后清单外那一发 fs.write 仍响（L1 窗口被问 1 次、零授权行、decision 非 allow_session_grant）＋「批过的范围」按批大小写死（文本必含「只限上面这 N 步」「清单之外」）；AC#4 只加发起请求的出口（Decl 与参数面零 allow/approve/grant 名，AppliedSteps=nil）；AC#5 本票字节面只 internal/tools 两枚新件＋探针。三发变异（仓外备份＋sha256 还原证）：M2 摘掉范围句 ⇒ Test165ScopeClauseNamesTheBatch 唯一红；M3 L1 支不再问 ⇒ Test165PlanApprovalRingsAgainForAnUnlistedWrite 唯一红；M1 Declared L2→L1 ⇒ AC#2 那条红（附证）；三发还原 sha256 与基线逐位相等。门禁：go vet ./internal/tools rc=0、go test -run Test165 五条全 PASS。⛔ 未翻 AC 框/Status、未动台账/HANDOVER/冻结面；只 commit 不 push。next=① 非实现者对抗验收（判语归另一 agent）② cmd/wisp 装配点接线（本票写面外，在飞 135-r1）③ 「方案条数上限」契约未说，本件未发明数字。

- [2026-10-16 10:54:54 +0800] agent=165-v1（非实现者终裁·只读）did=五格判语全「成立」，⛔ 翻框（翻法归编排者）：起手现量 HEAD `71002d00`、票框 5 未勾／0 已勾；台件＝`git archive` 三棵仓外导出树（`t0`=起树锚 `b74b87d0`／`t1`=HEAD 1091 枚 `.go`／`t2`=HEAD 留突变态），⛔ `-overlay`、⛔ 共享工作树。判语与凭据全文＝`docs/evidence/s1/165-v1-acceptance.md`（119 行／17810 字节），原始读数与逐发哈希＝`.scratch/wisp/probes/165/v1/00-readings.md`＋`01-mutations.md`。要点：① **AC#1 成立**——未修码那发是真跑的（`t0` 名册 10 枚、零枚 plan 形、`registry.Lookup=false`；把 r2 那枚调用点尺盖进 `t0` 靶向跑 rc=1、红句「调用点 0 枚…命中清单＝[]」），⛔ 写成"应该做不到"；② **AC#2 成立**（只裁"没新增状态"）——26 枚路径里 `internal/risk`／`internal/statemachine`／`internal/agent/approval`／`internal/panel` 各 **0 枚**，`Gate` 仍 **2** 枚回调、`GrantSource` 仍 **1** 枚，D43 `#17/#23/#24/#25` 现读＝`PLAN.md:3077/3083/3084/3085` 对 `table.go` 八条既存转移，反形我自己造（`Gate` 加第三条回调⇒唯一红「Gate 有 3 个方法」）；③ **AC#3 成立**——两发我自己重放且复现「唯一红」（M2 摘范围句／M3 L1 支强制 allow），另补四形都有牙（M2b 边界句在而数不跟批／MA 摘 `run.go:822-827`⇒RUN=3／FAIL_TOP=3／M9 假账本种一行授权⇒两条子判据同红／M8 真 `InsertGrant` 种一行⇒红在授权行计数那枚，⚠ 那一发卡照样升起：`Covering` 只短路 L1 支），并跑阳性对照证 `ListGrantsBySession` 0→1 读得出；④ **AC#4 成立**（按"能力类判据必问生产调用者"裁）——`Decl` 8 枚字段零 `allow/approve/grant`、剥注释全文 0 命中、顶层参数 2 枚键；产码调用点 **1 枚**（`cmd/wisp/run.go:822`，AST 尺读 `range=true register=true`）、`internal/panel` 面 **0 命中**（分母 41 枚跟踪 `.go`），面板四枚既有边界钉我在导出树自己复跑 rc=0／RUN=6（4 顶层＋2 子测试）／PASS_TOP=4／FAIL_TOP=0，⛔ 拿 `go build` 的 rc=0 当"接上了"；反形 M6b（顶层 `approved`）与 M7（`Decl` 加 `AutoApprove`）各唯一红，M6（同键塞在 `steps` 内部）rc=0＝具名**尺的射程边界**（`planPresentArgs` 无字段接它＋顶层 `additionalProperties:false` 兜住）⇒ 记射程、⛔ 记漏洞；⑤ **AC#5 成立**——两腿全部四枚提交（`4a49480da`＋`a754a390`＋`17ca21ee`＋`988ce413`）的 **26 枚路径**逐枚 pathspec 越界尺＝命中 **0 枚**，同尺正控喂 `824bb9c` 那 2 枚命中 2 枚；区间 93 枚里动过禁区的 8 枚逐枚具名归属、零枚署名 165。四处交付缺陷具名（全在描述面／尺射程，⛔ 影响任何一格成立）：ⓐ `plan_present_165_test.go:14-15` 头注释写「两发变异（-overlay 注入、树零改）」与同腿 `06` 读数件及提交信息（仓外备份＋sha256）自相矛盾、且 `-overlay` 冒充物理突变是本编队禁用形；ⓑ r1 票尾那句「M1 ⇒ AC#2 那条红」少报一枚顶层红（真量两枚，r1 自己的 `05` 件也记两枚）；ⓒ 编排者点名要本腿答的「r2 缺 10-* 那一格落过盘否」＝**判据落过盘**（件 08＋01b＋11），缺的是编号那一枚，⛔ 追认 10 号件曾存在；ⓓ 模型侧真调用**零读数**（全仓 56 枚跟踪件命中里夹具面 0 枚）⇒ 归票 161 那一格、建议挂 `Q##`。本腿不裁并往上交六项：翻框／`AwaitingApproval` 遇挂起（DEFERRED `PLAN.md:1530`，归 S7）／方案条数上限／面板那一屏（`frontend`、`design` 归机主外部 agent）／整包 `go test ./cmd/wisp`（票 314 死亡区，本波那一档归 `220-r2`）／`-tags winlive`。另具名一处别腿账：共享工作树 `go vet ./cmd/wisp` rc=1 红在**未跟踪**件 `cmd/wisp/panel_l1_blocked_220_test.go:236`（220-r2 在飞），同一把尺在 HEAD 导出树 rc=0。突变台件与 22 枚日志只建不删（`D:\tmp\165v1\{t0,t1,t2,kept,logs}`），六枚被改文件逐枚对 HEAD blob sha256＝全 MATCH。只 commit、⛔ push。

- [2026-10-16 11:0x +08] 编排者★**翻本票五格（AC#1–AC#5 全翻），⛔ 本票不改 `-done`**：判语来源＝非实现者腿 `165-v1`（件 `docs/evidence/s1/165-v1-acceptance.md`＝17810 字节／119 行，起手现量 HEAD `71002d0`，台面＝`git archive b74b87d0`／`git archive HEAD` 仓外导出树、⛔ `-overlay`）。AC#1 **成立**（未修码那一发真跑：名册 10 枚、零枚 plan/propose 形、`registry.Lookup("plan.present")=false`；反形加一味 `BuiltinPlanEntries()` ⇒ 11 枚／plan 形 1 枚＝尺不瞎）；AC#2 **成立**（只裁"没新增状态"：`internal/risk`／`internal/statemachine`／`internal/agent/approval`／`internal/panel` 各 0 枚，`Gate`＝2 枚回调、`GrantSource`＝1 枚，`plan_present.go` 里 `Answer/Mode/State` 常量引用 0 命中；反形给 `Gate` 加第三条回调 ⇒ **唯一红**；★本格⛔ 读成"功能做全了"，且 `docs/PLAN.md:2567` 自带那句"本票因此不减少确认次数"＝「一批只问一次」**不是**本票 AC）；AC#3 **成立**（两发我自己重放复现"唯一红"，另补四形都有牙：M2b／MA 摘 `run.go:822-827` ⇒ `RUN=3 FAIL_TOP=3`／M9 往假账本种一行授权 ⇒ 两条子判据同红／M8 真 `InsertGrant` 种一行 ⇒ 红在授权行计数那枚）；AC#4 **成立**（按"能力类判据必问生产调用者"裁：`Decl` 8 枚字段零 allow/approve/grant、产码调用点 **1 枚**＝`cmd/wisp/run.go:822`、`internal/panel` 面对 `plan.present`／`BuiltinPlanEntries` **0 命中**（分母 41 枚跟踪 `.go`）；反形 M6b/M7 各唯一红，**M6（同键塞进 `steps` 内部）rc=0 ⇒ 具名为尺的射程边界**、记射程⛔ 记漏洞）；AC#5 **成立**（两腿四枚提交的 26 枚路径逐枚 pathspec 越界尺命中 **0**，正控同尺 2 枚命中；区间 93 枚里动过禁区的 8 枚逐枚具名归属、零枚署名 165）。★我自己复跑：`git grep -c 'plan\.present' HEAD -- cmd/wisp` ⇒ `run.go` 1 枚＋两枚测试件各 2；`git grep -n BuiltinPlanEntries HEAD` ⇒ 9 枚；仓外 HEAD 树 `go test ./internal/tools/ -run 'Test165|Test287'` ⇒ `RUN=8 PASS_TOP=8 FAIL_TOP=0` rc=0。★**Status 更正**（⛔ 抹原句）：抬头那句"**立而不派**（新增一档＝动 `D31`/`D43`/`C12`，人工批准）"按 `A323`＋`docs/PLAN.md:2566` 的落地事实要读成＝**owner 已批准并落进冻结文本的是 `plan.present` 这枚工具名（L2），状态机与转移表一字未动**（AC#2 已由非实现者按名册自证）。★**本票⛔ 改 `-done` 的三笔残余（具名）**：① ★**门铃零读数**＝全仓跟踪件对 `plan.present` 命中 56 枚，同尺喂夹具面（`testdata/**`、`*.golden`、`*.sse`、`internal/llm/**`）＝**0 枚** ⇒ "模型今天真的会调这枚刹车"没有任何凭据，属接线／仪器缺口、⛔ 本票五格之内，归口见台账 `A1004`；② `internal/tools/plan_present_165_test.go:14-15` 的头注释写「两发变异（**-overlay 注入**、树零改）」，与同腿读数件／提交信息（仓外备份＋sha256 还原）**自相矛盾**，且 `-overlay` 冒充物理突变是本编队明令禁止的形状 ⇒ 待一枚 `internal/tools` 面空的微腿把注释改说实话（真量＝三发、非 overlay，该腿已逐枚对 HEAD blob MATCH 复证）；③ `cmd/wisp` 整包四数（票 314 死亡区）＝本波归写腿 `220-r2`。另具名一处**我派单的过期前提**：我给 `165-v1` 的第 5 行写"`216-r1` 在动 `internal/tools`"——现量＝`216-r1` 在 `b74b87d0..HEAD` 这 93 枚里**零提交**，`internal/tools` 真被 `03da4941b`（270-r1）／`e1973191b`（287-r1）动过（两枚都没碰 165 任何一件）；该腿顶回、我复跑认账（记我，台账 `A1004`）。

---

- [ ] **AC#6（编排者 2026-10-16 12:0x 新开，来源＝只读普查腿 `165-c2` 的三处断口；⛔ 改已冻结判据一字，本格只加不减）**：**"模型真的会调它"这一维今天零仪器**——票面既有五格全勾，但五格量的是"条款在不在、装配接没接、名册出不出得来"，⛔ 有一枚常驻判据管"请求里带没带这枚工具＋`tool_call` 那一跳真被走过"。三处断口逐枚具名（凭据＝`docs/evidence/s1/165-doorbell-census.md`，105 行／18148 字节）：① `internal/agent/loop.go:637-640` 的 `tool_call` sink 在，但**零枚用例**把 `plan.present` 打到这两行；② `cmd/wisp/plan_assembled_165r2_test.go:123` 只断到名册出口，跳 6/7/8 零覆盖；③ `llm.Request.ToolChoice` 在产码里 **0 枚**写入 ⇒ `tools/mockllm/chat.go:172-185` 那枚"必发工具调用"的合成从装配面按不到（且它按位取 `Tools[0]`、⛔ 按名）。★普查分母（全取自 HEAD 导出树，带正控）：`.sse` 52／`*/testdata/*` 56／`internal/llm/**` 90／`tools/mockllm/**` 14／`*.golden` **真 0**（正控 `-iname '*golden*'`＝12）⇒ **本仓注册的 11 枚工具名在夹具面逐枚 0 命中**＝缺口是整族、⛔ `plan.present` 独有。★选形已裁＝**甲**（新建 `internal/agent/testdata/golden/plan-present.sse` ＋ 1 枚用例；该目录**无孤儿钉**：`ReadDir` 三处全在 `t.TempDir()`；已知钉 0 枚；⛔ 碰共享编译面）。★**必须先接住的一枚硬约束**：`internal/agent` 根 21 枚测试件逐枚 `package agent`、`agent_test` 全树 **0 枚**，而 `bridge.go:13` 已 import `agent` ⇒ **真桥进不去＝导入环**，落地腿只能走包内形状。落地腿＝`165-r3`（写面只 `internal/agent/**`）。⛔ 丙形（真机手测）当仪器——它 0 枚 `.go`、属〔仅本机可量〕＋需 key＋winlive（该腿严尺现量 14 枚＝`cmd/wisp` 9＋`internal/ball` 5；台账旧枚数"8"复现不出）。
