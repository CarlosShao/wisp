100	3	**Status:** **RESERVED（不是 open，也不是 blocked-on-owner）**：AC#0 那条触发门**没满足之前，本票不许开工**，
103	3	**Status:** **rejected-needs-fix**（验收判 `PASS WITH CONDITIONS`，但**编排者据其探针把它退回**：
107	3	**Status:** **rejected-needs-fix**（2026-09-21 18:5x 编排者退回。验收判 **AC#3 不通过**（AC#1/2/4/5 通过），
108	3	**Status:** **rejected-needs-fix**（2026-09-21 20:3x 编排者标注，依据=`acceptor-ticket108` 的
109	3	**Status:** **rejected-needs-fix**（2026-09-21 21:2x 编排者按 `acceptor-ticket109b` 的裁决表改判：**总判 FAIL-退回**。
111	3	**Status:** open（2026-09-21 20:0x 编排者建；来源=`agent-ticket110` 的 AC#1 全仓对账，run `35591482293`）
112	3	**Status:** ready-for-review（2026-09-21 20:3x 编排者建；来源=**run `35595651898` / job `106319703680`**，票 110 新加的 "Windows ACL sealing gate" 步）
114	9	**Status:** in-progress（2026-09-23 10:4x `agent-ticket114-ac2` 交 AC#2 原生侧门 + AC#1 现状表入库；
12	3	**Status:** review
120	3	**Status:** open（2026-09-21 21:2x 编排者建；来源=`acceptor-ticket113` 的
122	11	**Status:** open（2026-09-21 21:3x 编排者建；来源=`audit-85-preflight` 的
123	3	**Status:** open（2026-09-21 22:0x 编排者建；来源=`ci-read-111b` 的
124	3	**Status:** open（2026-09-22 16:41 编排者建；来源 `agent-ticket119` 交回待裁第 4 条 + `acceptor-ticket119` 的 `R-119-8`）
128	3	**Status:** in progress（2026-09-22 16:41 编排者建；来源 `acceptor-ticket117` 的 `R-117-B`。AC#1 量毕；**AC#2 落地 + AC#3 变异自证已完成**（`agent=T128-ac23`，报告 `docs/evidence/s1/128-ac2-refusal-and-ac3-mutation.md`）；**AC#4 门禁已收**（非实现者终裁表 `docs/evidence/s1/128-ac4-r1-acceptance.md`，成立无附条件，回执在 AC#4 格下方 
130	3	**Status:** open（2026-09-22 22:3x 编排者建；来源 `acceptor-ticket117` 的 `R-117-2` + `agent-ticket125` 的 `R-125-3`，`acceptor-ticket127` 判定"不推翻票 117、是另账一张票"）
132	3	**Status:** ready-for-agent（2026-09-23 09:4x 编排者建；来源＝owner 对 `Q-31` 的批复「算、要锁」+ 本次实测的零调用者读数）
133	3	**Status:** ready-for-agent（2026-09-23 09:5x 编排者建；来源＝`acceptor-ticket121` 的
135	3	**Status:** ready-for-agent · **Type:** 判据仪器缺陷（后果断言**没有自证腿**）
136	3	**Status:** ready-for-agent（2026-09-23 18:2x 编排者建；来源＝`acceptor-ticket134-r2` 的总判
138	3	**Status:** **`blocked`（09-26 00:5x 编排者按 138 验收程的第 4 节判；等 `Q-55` 那一次批准一起落）**
139	3	**Status:** ready-for-agent（2026-09-24 10:0x 编排者建；来源＝owner 批准的 `Q-43` 后半，
140	3	**Status:** open（2026-09-24 17:4x 编排者建；**尚未派发**）
145	3	- Status: ready-for-agent（**未派**，见末节"为什么现在按住"）
148	3	- Status: **blocked（两枚前置都是人工批准，不许实现方自决）** —— 见下面「批准清单」，回两个字母即可解锁
15	3	**Status:** ready-for-agent
150	3	- Status: **blocked（三选一的前置是人工批准——三条出路都要动冻结面或动契约句；见 `Q-55`）**
158	3	- Status: ready-for-agent（**立而不派**：此刻 157 的 AC#5 重判程在飞（只读 docs），156 的续程按住等 `cmd/wisp` 的读者散开——本票要跑 `internal/tools` 与 `cmd/wisp` 两包，**等编队空**）
159	3	- Status: ready-for-agent（**立而不派**，排在票 156 结案之后、票 158 之前或之后由编排者定）
16	3	**Status:** ready-for-agent
160	3	- Status: **ready-for-agent（缩窄版）**——owner 09-27 09:30 那句「那六张票全做，按你的推荐可以」**批准单位＝开这块地**（台账 `A322`：放开 `internal/risk/provenance.go`（含测试）＋两枚**文件**的必要改动；⚠ 调用点是**三处**＝`bridge.go:642` 开、`bridge.go:695` 关、`panel_assets.go:232` 只开不合，`A327` 现量纠正）。撤销口令「160 别动」。⚠ **两条不在批准范围内、也不许动**：`internal
163	3	- Status: **立而不派**（要扩 `shell.exec` 的语义＝动 `D34`/`D46` 那条线，**人工批准**；且它和票 161 抢同一批审批判据）
164	3	- Status: **立而不派**（`D34` 表里那两枚名存实无，先要一次名册更正＝**人工批准**，再谈实现）
165	3	- Status: **立而不派**（新增一档＝动 `D31`/`D43`/`C12` 那条状态与审批线，**人工批准**；且 owner 已定过"改状态机要谨慎"）
166	3	- Status: **ready-for-agent，但立而不派**（排在票 158 验收与 161 之后；本票只做 **Go 侧**，`skipped=frontend(owner-delegated)`）
167	3	- Status: **ready-for-agent，但立而不派**（排在票 145 之后——它俩抢同一批快照字段；`skipped=frontend(owner-delegated)`）
168	3	- Status: **立而不派**（前置＝票 50（Tier-1 清单插件）＋票 51（Tier-2 goja）**都没做完**，插件地基不存在；且它要动 `C17`/`C24` 那两张**未定案的白名单**＝**人工批准**）
169	3	- Status: **owner／那一支写着 `frontend/**` 的会话的接手件**——本编队（编排者＋子代理）**结构性写不到那儿**（`frontend/**` 自 09-25 起在每一张票的禁改名单里，票 158 AC#4 逐字列着它）。本票**不派、不 claim、不加 `-done`**，存在的唯一目的＝把那枚残余从"两个人都自报了、没人接"变成**带票号的书面接手**（本仓规矩：不收的格必须写成有名字的接手，"暂缓"无票号不合格）。
170	3	- Status: **立而不派**（推荐那一支很可能要碰 `internal/agent/approval/**`＝冻结面；且"要不要让 CLI 变成会停下来等答案的东西"是产品口径，不是我该拍的）
171	3	- Status: **ready-for-agent**（写面＝`probes/**` 的尺本体与 `.github/workflows/ci.yml` 的**新增步骤**；不碰 `internal/**`、`cmd/**`、`tools/d22scan/**` 的禁令射程。派单前置：编队里 `internal/tools` 那两程交完，因为本票要跑全仓级那三把尺）
172	3	- Status: **ready-for-agent**（写面＝一枚新尺＋它的用例；**读** `docs/PLAN.md`／`docs/specs/**` 可以，**改它们一字都不许**；若要进 CI 只能新增步骤，动触发表＝契约级 ⇒ 停手上报）
173	3	- Status: **ready-for-agent**（写面＝`internal/tools/bridge.go` 那一味 + 它的用例；⚠ 若判"审计行格式属冻结契约"⇒ 停手上报，见 AC#3）
174	3	- Status: **ready-for-agent**（写面＝`internal/tools/**` 与 `internal/agent/spill.go` 的**回执文案/断言**侧；⚠ 本票**不碰** `[fs]` 的默认值、不碰 C26 的判定语义、不改 `docs/PLAN.md` 与 `docs/specs/**` 一字。派单前置：**无**（测量腿现在就能跑），但 AC#2 的"怎么修"要等 `Q-60` 有 owner 的一句话）
175	3	- Status: **ready-for-agent**（写面＝`internal/risk/provenance.go` 的名册**或** `internal/tools/bridge.go` 的 `mark` 判定形，**两枚都属已批准射程之外**：改名册＝契约追加，须 owner 批；改判定形＝动 C25 语义，亦须批。⇒ **本票第一腿只裁不改**，裁完才知道写面落哪一枚。派单前置：票 164 的 `164-v1` 验收交件——它正在读同一批文件）
176	3	- Status: **立而不派**（编排者排程：**排在票 175 的落地腿之后**——那枚写者一动手就必须与本票同批，见 AC#2；且它与票 163／165 抢同一批 `internal/tools/**` 文件，串行）。
177	3	- Status: **ready-for-agent（只读设计核腿）**——第一腿是裁"该走哪一支＋那一支的最小形状"，**零产码**；写手腿**按住**等 owner 对 `Q-61` 的一句话（三支候选各动不同的冻结面）。派单前置：**无**（只读腿现在就能跑；此刻也没有别的程在读那两族尺）。
178	3	- Status: **ready-for-agent（只读＋改仪器）**——量清射程＋给修法，改动只落在 `.scratch/wisp/probes/154/gate-clauses.sh` 这一枚仪器（**不是判据、不是阈值**，但要按"改仪器要留凭据"的规矩办）。
179	3	- Status: **ready-for-agent（写码＋判据）**。本票是**缺陷修复**，方向由已定案的 D4 决定，**不是契约变更**（理由见下面"裁定轴"）。
180	3	- Status: **ready-for-agent（只读普查＋一支小落地）**。
181	3	- Status: **ready-for-agent（只读普查＋一支只读落地）**；⚠ 本票**不含**任何改变 git 状态的写侧动作，那半在 `Q-64`、默认不做。
182	3	- Status: **ready-for-agent（只读普查，零产码）**。
185	3	- Status: **已立、按住不派**（等 `183-v1` 交完并由我裁归属：本票与票 183 AC#2 的措辞是同一场官司）。
186	3	- Status: **ready-for-agent（先只读普查，再落地）**。⚠ 本票是**功能票**，不是缺陷票。
187	3	- Status: **ready-for-agent（先只读普查，再落地）**。⚠ **功能票**。
188	3	- Status: **ready-for-agent（先只读设计核，再落地）**。⚠ **功能票**。
189	3	- Status: **ready-for-agent（先只读设计核，再落地）**。⚠ **功能票**，且**与票 181 同一条读面**（复用，别重造）。
190	3	- Status: **ready-for-agent（先只读设计核，再落地）**。⚠ **功能票**。
191	3	- Status: **ready-for-agent（AC#0 已由编排者按甲裁结·09-28 16:4x，见 Progress log 末段与台账 `A379`）**。⚠ **功能票**（原"批准项"那一半已解除：owner 给的是**做法授权**，不是又一次待批）。
192	3	- Status: **ready-for-agent（AC#0 由编排者按"零 key 抓结果页"定案·09-28 16:4x，owner 原话与撤销口令在 Progress log 末段／台账 `A379`）**。 **功能票**（原"在册待定项"那一半由 owner 的"也做一下吧"接掉，我不拿它继续挡功能）。。
193	3	- Status: **ready-for-agent（先只读列门禁清单，再逐枚放开）**。⚠ **口径变更票**（owner 09-28 原话入 `A373`）。
194	3	- Status: **ready-for-agent（先只读核代价，再分批落地）**。⚠ **功能票＋契约变更已获点名批准**（owner 09-28 16:5x 原话逐字在台账 `A381`）。
195	3	- Status: **blocked（排在票 33 片 B 交件之后；地基未建之前不许接面板口）**。⚠ 本票**不是**新造功能请求：来路＝owner 09-28 已裁的「方法名册按当前规格补，别搞债务了」（逐字 `A381`），`config.get`／`config.set` 逐字在 `docs/specs/SPEC-08-ui-ball-panel.md:168`；**本票只是把那一支的落地顺序改成"先地基、后面板"**。
196	3	- Status: **未派（先摆 owner 一次：要不要给 schema 加约束＝可能触 `D35` 数据模型）**。⚠ **本票不是"缺字段"**（那是票 188），本票是**已经并存的二致性缺陷**。
197	3	- Status: **在派（批准已落，见下面那一节）**。⚠ **功能票，owner 09-28 18:0x 原话逐字**：「**子代理必须看到状态，而且点击某个子代理，能看到它们各自的流式工作页面，能明白吗？这是主流 harness 必做的，不要偷懒**」⇒ 这一句同时是 `Q-71` 的裁定＝**乙（子代理上界面）**，票 188 的 `AC#1` 已按此勾。
20	3	**Status:** review（两段 fs 工具族已落地并由编排者验收：`docs/evidence/s1/20-needs-assertion-restored.md`；
200	3	- Status: **在派（owner 09-28 18:1x 放权后自行排入，不等批准）**。
201	3	- Status: **部分结案（09-29 12:0x：七格里已勾 2＝AC#4／AC#7；未闭 5＝AC#1／AC#2／AC#3／AC#5／AC#6，前置具名见下面那节）**。⚠ **本票剩余的 AC#1／AC#2／AC#3a／AC#6b 全部归口同一枚前置——"宿主原生入口（球／托盘）在 `wisp run` 里没接"**，Go 侧的答复逻辑今天第一次全齐；⛔ **不翻 `-done`**（`-done` 是防重领唯一键，须 0 枚未勾才结）。早先那行"在派队列里"的口径已被上面这行取代，原文不再单独保留。
21	3	**Status:** segment-1-landed-decision-layer (L1 window + veto channels + L2 queue + source check + batch + applied-steps type); segment 2 (native card rendering, ball visuals, loop rewiring) still open
211	3	- Status: **在派队列里（本票不是"漏做"，是"做到今天诚实的那条线为止"；往上走一格要 owner 一句话）**。
213	3	- Status: **待派（owner 09-28 22:1x 口头接单，原话逐字见下面第一节）**。这是**功能请求，不是待拍板项**——按他今天的口径：**做，且做完整**；风险只翻成下面的实现约束。
214	3	- Status: **待派（owner 09-28 22:1x 口头接单，原话见票 213 第一节）**。这一票是**接线＋补最小缺口**，不是从零做——**这是今天查出来的最便宜的一枚**。
215	3	- Status: **待派（owner 09-28 22:1x 口头接单，原话见票 213 第一节）**。
216	3	- Status: **待派**。横跨票 213／214／215 的**每一枚字符串**，所以单独立一枚（owner 22:1x 的原话是「**聊天框这块的内容不少的**」，这一枚就是那块地基）。
219	3	- Status: **接单＝做**（owner 原话见下）。**这一票不是"要不要做"，是"按哪条通道做"**：三枚按钮的形状已经定在下面 §2，
22	3	**Status:** ready-for-agent
220	3	- Status: **待派**。来源＝只读腿 `197-c2` 的 ④ 节＋编排者顺着 `queue.go` 读出的第二形；现读证据件＝`.scratch/wisp/probes/197/r5c/census.md` ④（**每一行都由编排者自己在现树重跑过尺**）。
222	3	- Status: **产码已交并静态核过（09-29 12:46，见文末"收 `222-r1`"那一节），⛔ 六格今天一枚没勾**——两枚前置未到：① **我自己的整包复跑**（被在飞写腿 `226-r1` 挡住：它正在改 `internal/config/allowdirs.go`／`loader.go`／`permmode.go`，而 `cmd/wisp` 与 `internal/tools` 都吃 `internal/config`⇒ 现在跑编译读到的是它半截的码）；② **非实现者终裁 `222-v1`**（`SPEC-12 §4.3`／`AG
223	3	- Status（**09-29 16:5x 收写腿 `223-r2`，表＝`docs/evidence/s1/223-hot-reload-wiring-r2.md`（骨架 `28475620`→代码＋读数 `5a755c3c`→收工读数若干 commit），台账待编排者落 `A##`）：⛔ 仍七格不翻勾、不翻 `-done`（写腿不自裁，D22 双角色）。r2 只做派单两件事、零扩射程，交回的是读数**：①**AC#5 那枚本票红已改确定性**——`TestTicket223RestartTierSaysItWillNotApply` 的 stdou
224	3	- Status: **待派（前置已解＝09-29 13:13 批准记录 `A435`：会话＝什么 已定死、只读 seam 已批，撤销口令「撤 224 会话定义」；`PLAN.md` 与 `docs/specs/**` 一字未动）**。来源：非实现者设计复核 `219-v0`（`docs/evidence/s1/219-approval-reply-design-adversarial.md`，27,833 字节）第 3 条＋编排者复跑。
225	3	- Status: **待派（只读审计腿即可，不产码）**。来源：票 223 的普查腿 `223-c1` 顺带撞出、编排者现跑复认（台账 `A426`）。
226	3	- Status: **产码已交（09-29 13:04，四枚 commit `78988901`／`361b314a`／`d90b9aec`／`c31239a9`），⛔ 六格今天一枚没勾**——非实现者对抗验收腿 **`226-v1`** 正在飞（表＝`docs/evidence/s1/226-config-write-no-clobber-v1.md`；实现腿自报 AC#1/2/3/4/6 已闭、AC#5 半格，**我一个字都不预认**，理由＝`SPEC-12 §4.3` 与本仓 `A428` 定死的口径"裁决腿的判语不自动等于翻勾、字面判据优先"）
227	3	- Status: **待派，排在票 223 之后**。来源：票 226 的非实现者对抗验收腿 `226-v2`（`docs/evidence/s1/226-config-write-no-clobber-v2.md`，31,097 字节，六枚 commit `ee40bbd0`→`172715d6`）第 **§⑤**（删行一问的结论）与第 **§⑥ 条 1/3**（推翻实现腿两句），加**编排者自己现读复认**（台账 `A437`）。
228	3	- Status（**09-29 16:1x 收 `228-a1`，台账 `A439`）：AC#0 成立并翻勾**（代价表＝`.scratch/wisp/probes/228/a1/cost-table.md`，69,717 字节／354 行／10 枚提交，编排者自己复跑六把尺全部对上）；**方向由编排者按 owner 长期原话裁死＝走甲**（把审批门与任务通路接进**常驻**那条腿），⛔ 不是"待他选"——他的原话是「**别人怎么做的，你就怎么做**……**必须特么做完整功能**」，而 D2（`PLAN.md:74/83-88/:103`）本来就把球＋
229	3	- Status: **待派，小活**（队尾，排在票 228 之后或与票 227 同批）。⚠ **两支修法里只有一支不用问人**：把 spawn 名换成名册内的形状＝可直接做；往名册里加 `subagent-finish-`＝**动 D38b 那张表＝人工批准**，须先落 `A##` 并摆给 owner。⛔ 第三支"把那行警告调轻／删掉"＝**削弱仪器**，本票不许。
23	3	**Status:** ready-for-agent
230	3	- Status: **待派，小活**（排在票 229 之后；其中 AC#2 动 `internal/agent/approval/queue.go` ⇒ 与票 220 同包**串行**）。来源＝只读盘点腿 `pool-2`（`.scratch/wisp/probes/pool-2/undone-in-done.md`，80,574 字节／531 行／10 枚提交）＋**编排者 09-29 16:4x 自己复跑四把尺复认**（台账 `A440`）。⛔ **不许靠把 `-done` 改名回来"解决"——那个后缀是防重领的唯一键**（`issues/READ
232	3	- Status: **待派，测试面加固**（与票 231 同撞 `cmd/wisp/config_reload*.go` ⇒ 同批或串行）。来源＝非实现者验收腿 `223-v2` 第②节（它**造出了那一发突变并跑出绿**，不是我推的）＋台账 `A444`。⚠ **形状不是 r2 引入的**（`223-v2` 具名：归 r1 那发形状），但今天还在。
233	3	- **Status**：**待派**（编排者 09-29 17:5x 立，起手锚点 `a8f3c020`）。来源＝非实现者验收腿 `223-v2` 交回的"要人拍的第 3 笔"（`docs/evidence/s1/223-hot-reload-wiring-v2.md` 的「没做完／判不了」第 1 条）＋票 223 **AC#2 那半格**（台账 `A445`）。
234	3	- **Status**：**待派，且必须按住到票 222 的验收腿交完**（本票的腿要跑突变＝编译面，⛔ 与任何突变腿并发都会互相洗读数）。编排者 09-29 18:3x 立，起手锚点见 `git log -1`。来源＝只读盘点腿 `undone-28-1`（件 `.scratch/wisp/probes/undone-28-1/table.md`＝**74,413 字节／687 行／13 枚提交**，附录 D 的四分类小结），台账 `A446`。
236	3	- **Status**：**待派**（编排者 09-29 20:2x 立，起手锚点 `c5d88a7f`＝本票现读 `git rev-parse --short HEAD`）。
237	3	- **Status**：**待派**（编排者 09-29 23:5x 立，起手锚点＝本票现读 `git rev-parse --short HEAD`）。来源＝**非实现者裁决腿 `235-v1`**（`docs/evidence/s1/235-pool-nail-and-bridge-preread-v1.md`，**27,557 字节**，六枚 commit `a207d129`／`48e09742`／`8f59a180`／`b490f8ae`／`141395ae`／`c28858fc`）。台账 `A458`。
238	3	- **Status**：**待派**（编排者 09-29 23:5x 立，起手锚点＝本票现读 `git rev-parse --short HEAD`）。来源＝独立复核腿 **`132-c3`**（`.scratch/wisp/probes/132/c3/verify.md`，**64,812 字节／473 行**）第 5 条顶正；台账 `A458`。
239	3	- **Status**：**待派，小活**（编排者 09-30 11:1x 立，起手锚点＝本票落地的那枚 commit）。
24	3	**Status:** ready-for-agent
244	3	- Status: **已立，未派**（09-30 16:1x，编排者立）。来源：只读普查腿 `243-c1` 交的表里第 22／23 行（裁决件 `.scratch/wisp/probes/243/c1/census.md` §3 把这两行列为"最危险"），**指认由我自己逐把尺复跑复认**（下面现量表全是我 16:0x 自己跑的，不是转抄腿）。
245	3	- Status: **已立，未派，排在队列最前**（09-30 16:3x，编排者立；触发＝实现腿 `228-r1` 把球装进了常驻那条腿，**这件副作用从今天起在正常跑的产品进程里生效**，不再只在调试旁支程序里）。
246	3	- Status: **已立，未派**（09-30 18:2x，编排者立；账 `A480`）。触发＝验收腿 `245-v1` 交回的一条事实，我自己现跑复认过。
247	3	Status: OPEN（编排者 09-30 23:0x 立，料全部出自 `240-c1`／票 241 与本编排者 22:5x–23:0x 的现跑尺）
248	3	Status: OPEN（编排者 09-30 23:1x 立，来路＝owner 当场提的功能要求，原话见下面第一行）
249	3	Status: OPEN（编排者 10-01 **13:13** 立，料全部出自探针腿 `33-p1` 与本编排者 13:1x 的现跑尺；台账 `A502`）→ **WITHDRAWN（10-01 14:2x 编排者自撤；出处＝台账 `A503`；⛔ 任何腿不许据此开工）** — 撤回原因：**上面标题与「现量」表第 2 行的那句前提是错的**，我立票时没去读决定它的那三行码；原文一字不抹，判定与复认尺全在文末「撤回理由」一节。撤回前它未派出、零枚提交、零格被勾。
25	3	**Status:** ready-for-agent
252	252-the-two-spellings-of-one-path-split-inside-the-allowlist-check-short-name-stays-short-long-name-folds-long-so-in-root-writes-read-as-out-of-bounds-r2-l2.md	NO_STATUS_LINE
253	253-panel-inbound-three-ruler-holes.md	NO_STATUS_LINE
255	255-config-says-these-sections-took-effect-immediately-panel-for-a-width-nobody-reads-while-the-panel-host-never-receives-config-and-tier-has-zero-readers.md	NO_STATUS_LINE
256	256-resident-gate-built-before-session-grants.md	NO_STATUS_LINE
26	3	**Status:** ready-for-agent
260	260-configured-cancel-hotkey-never-registers.md	NO_STATUS_LINE
262	262-tracked-path-length-gate-for-ci-checkout.md	NO_STATUS_LINE
264	264-error-text-carries-absolute-paths-into-model-visible-text.md	NO_STATUS_LINE
266	266-no-external-tooth-over-the-slo-check-self-locking-nails.md	NO_STATUS_LINE
269	3	- Status: **已立，未派落地**（10-05 11:5x 编排者立；来源＝取证腿 `ci-attr-1`（`.scratch/wisp/probes/ci/1005-a/verdict.md`，174 行／28,313 字节，commit `d7236abc`）的"该立案的红名清单"第①枚，**每一把我自己复跑过**）。撤销口令「**269 撤**」。→ ★**WITHDRAWN／已并案（10-05 12:3x 编排者自撤，出处＝台账 `A620`；并案回票 236 AC#6）**——⛔ **任何腿不许据此票开工，也不许再派腿**：只读腿 `2
27	3	**Status:** ready-for-agent
270	3	- Status: **已立，未派**（10-06 18:3x 编排者立；来源＝非实现者终裁腿 `236-v1` §5 第 2 行／§2 第 6 行，裁决表 `.scratch/wisp/probes/236/v1/verdict.md`；台账 `A643`）
271	3	- Status: **已立，未派**（10-06 18:3x 编排者立；来源＝非实现者终裁腿 `236-v1` §2 第 8 行的附带读数，裁决表 `.scratch/wisp/probes/236/v1/verdict.md`；台账 `A643`）
272	3	- Status: **已立，未派**（10-06 19:4x 编排者立；来源＝非实现者终裁腿 `232-v1` §5 第 3 条，裁决表 `.scratch/wisp/probes/232/v1/verdict.md`；台账 `A647`）
273	273-shipping-process-builds-state-machine-without-a-sink-so-every-d43-effect-falls-into-no-op.md	NO_STATUS_LINE
274	274-no-nail-requires-shipped-exe-to-carry-page-build-ps-step-2-comment-is-false-in-both-halves.md	NO_STATUS_LINE
275	3	- Status: **已立，未派**（10-07 12:1x 编排者立；来源＝写码腿 `111-r5` §4b 顺手量到的真发现，件 `.scratch/wisp/probes/111/r5/ci-guards.md`；台账 `A658`）
276	276-gofumpt-census-still-exits-1-on-seventeen-probe-artifacts-from-fourteen-other-tickets.md	NO_STATUS_LINE
279	279-p39-comment-and-the-255-roster-citation-ride-together.md	NO_STATUS_LINE
28	3	**Status:** ready-for-agent
284	284-gate-replays-comment-omits-the-zero-production-caller-claim.md	NO_STATUS_LINE
287	25	**Status:** **AC#1 已闭合，余三格未开工**（2026-10-09 12:2x 编排者：只读腿 `287-a1` 交件、我复量关键锚后翻 `AC#1`，前提三把尺全部复现 ⇒ 本票不作废）。排程＝⛔ 仍按在 Go 编译/测试面空（此刻编排者自己占着 `cmd/wisp` 跑基线），且 `AC#2` 是**写腿**（要造夹具＋跑 `go test`），⛔ 不许与 `247-r1` 同批——`247-r1` 独占 `cmd/wisp` 与整个导入图。`AC#3` 那张代价表是只读格，可随下一枚只读腿补派。⛔ 零 push。
288	23	**Status:** **未开工**。排程＝随下一枚真动 `cmd/wisp/**` 的写腿**同批捎带**最省（单独一发要独占 Go 面、不值）；若没有那样的腿在飞，本票可以单独排一小发。⛔ 零翻框、零 push。
29	3	**Status:** ready-for-agent
290	84	**Status:** **`AC#0`/`AC#1`/`AC#2`/`AC#4`/`AC#5` 已裁并翻勾（16:5x，凭据全出自非实现者验收腿 `290-v1`）；`AC#3` 未勾**——欠的是**机主在场那一发真机读数**（见上面处置节首条），⛔ 不据此判失败、也不算腿的欠账；残余四问已归口票 293/294/295。旧 Status 句逐字＝「**`AC#0`/`AC#1` 已裁并翻勾（14:2x）；落地腿未开工**——`AC#2`..`AC#5` 四格未勾，排程见上面裁定末条。」**按 16:5x 处置作废（落地腿 `290-r1` 已交、验
291	25	**Status:** **`AC#1` 已裁并翻勾（14:2x）；落地腿未开工**。⛔ `AC#2`/`AC#3`/`AC#4` 三格未勾——`AC#2`＝机主在场那发四形真读数（与票 290 落地腿、票 247 三格同一次窗口发），`AC#3`/`AC#4` 的落地那一发⛔ 按在 `AC#2` 读数交回**并且**曲线那一问摆过机主之后（理由见上面裁定节第一条：`A466` 早已裁"必须等真声音实测读数再摆"）。⛔ 零 push。（本行是归位；旧行原文逐字留此备查：「**未开工**。排程＝**只读 AC#1/AC#2 可先派**……AC#3/AC#4
293	71	**Status:** **`AC#0`／`AC#1` 已由只读腿交付并经编排者裁完（20:2x，见上节：甲形、同批改写那段"故意不镜像"的注释、三条禁区）；`AC#2`..`AC#5` 一枚未翻，等 `293-r1` 与非实现者 `293-v1`；真机那形归合并窗口。**⛔ 零 push。
294	53	**Status:** **`AC#0` 已由只读腿交付并经编排者 HEAD 复跑翻勾（21:1x，见上节）；`AC#1`..`AC#3` 只裁到"落点形状＋缺哪件仪器"，读数（红/绿、真窗 `IsLive(hkMute)`、`executed==true`＋M2 反形）一枚未翻，等 `294-r1` 与非实现者 `294-v1`；`AC#4`/`AC#5` 归落地腿。★排程新增硬约束＝`295-r1` 先于 `294-r1`（同一段 `toggleMute`）。**⛔ 零 push。
296	85	**Status:** **`AC#3`／`AC#4` 已由编排者亲跑＋自跑门禁翻勾（22:2x，见上节：真机 `live=3` 存活 3 分、`live=0` 0 命中、`Ctrl+Alt+M` 开→关两行读数、整包两发 275/385/2/4 逐字相同）；`AC#0`/`AC#1`/`AC#2` 三格仍 `- [ ]` 等非实现者 `296-v1`；`AC#3` 的 ⓒ（托盘「退出」＋`levels_delivered=` 优雅收尾）具名留给合并真机窗口。**⛔ 零 push。
297	55	**Status:** **未开工**，但**现象已定实**（六窗读数在盘上、`LevelOfSamples` 零形有钉、`convertPacket` 零仪器三条都现跑核过）。⛔ 零翻框、零 push。
298	37	**Status:** **未开工**。低利害内务，⛔ 不上机主清单。⛔ 零翻框、零 push。
30	3	**Status:** ready-for-agent
300	35	**Status:** **未开工**（本票 10:5x 立，六格 `AC#0`..`AC#5` 全未勾；`AC#0`→`AC#1` 有**次序硬门**：布局权威没定，用例的红什么也证不了）。⛔ 零翻框、⛔ 零产码、⛔ 零 push。
301	301-windows-tagged-tests-in-a-core-claimed-package-never-run-in-ci-and-guard-d-is-package-level.md	NO_STATUS_LINE
302	302-page-reply-instruments-run-in-the-ci-cli-tier-and-can-never-get-a-reply.md	NO_STATUS_LINE
303	303-page-to-go-reply-hop-died-and-fails-on-owner-machine.md	NO_STATUS_LINE
304	304-external-killed-test-binaries-orphan-mockllm-exe.md	NO_STATUS_LINE
305	305-after-the-channel-fix-the-live-document-is-still-the-probe-and-eval-push-never-reached-it.md	NO_STATUS_LINE
309	5	**Status:** 未开工（六格全未勾）。⛔ 零翻框（翻框归非实现者验收程）、⛔ 零 push。
31	3	**Status:** ready-for-agent
310	310-test-binary-deaths-leave-temp-test-dirs-accumulating.md	NO_STATUS_LINE
311	5	**Status（2026-10-12 16:3x 更新，⛔ 旧行不删）：** 本票现尺＝**6 勾 / 0 未勾**（尺＝`grep -c '^[[:space:]]*- \[ ]'`＝**0** ⇄ `grep -c '^[[:space:]]*- \[x\]'`＝**6**）⇒ ★**六格齐但仍⛔ `-done`**，残余三笔逐枚具名＝①**新登 `Q-311-r1b-1`**＝`wisp-cli` 那档的**形①（默认根）／形②（同卷换目录）／摘除侧至今没有任何人跑过**（两枚墙：那两形会 `bash portable --scope=cli` 
312	312-scan-local-for-existing-model-configs-nice-to-have-parked.md	NO_STATUS_LINE
313	313-default-provider-stepfun-and-model-config-stays-in-the-ui.md	NO_STATUS_LINE
314	314-cmdwisp-suite-process-death-in-window-zone.md	NO_STATUS_LINE
32	3	**Status:** ready-for-agent
33	3	**Status:** ready-for-agent
34	3	**Status:** ready-for-agent
35	3	**Status:** ready-for-agent
36	3	**Status:** ready-for-agent
37	3	**Status:** ready-for-agent
38	3	**Status:** ready-for-agent
39	3	**Status:** ready-for-agent
40	3	**Status:** ready-for-agent
41	3	**Status:** ready-for-agent
42	3	**Status:** ready-for-agent
43	3	**Status:** ready-for-agent
44	3	**Status:** ready-for-agent
45	3	**Status:** ready-for-agent
46	3	**Status:** ready-for-agent
47	3	**Status:** ready-for-agent
48	3	**Status:** ready-for-agent
49	3	**Status:** ready-for-agent
50	3	**Status:** ready-for-agent
51	3	**Status:** ready-for-agent
52	3	**Status:** ready-for-agent
53	3	**Status:** ready-for-agent
54	3	**Status:** ready-for-agent
55	3	**Status:** ready-for-agent
56	3	**Status:** ready-for-agent
57	3	**Status:** ready-for-agent
58	3	**Status:** ready-for-agent
59	3	**Status:** ready-for-agent
60	3	**Status:** ready-for-agent
61	3	**Status:** ready-for-agent
62	3	**Status:** review（**改回 review：本票从未被验收**，见下方更正块）
64	3	**Status:** review（票面工作全部落地或已显式移交；剩 4 框各归其主，见 Progress log 2026-09-20T14:19Z）
65	3	**Status:** blocked-on-owner（缺参考图，见 §第一步）
68	3	**Status:** review（AC#1 静态三列表 + AC#4 已交付并由编排者复现勾框；**AC#2/AC#3 blocked-on-owner**，见 R15 #3/#4/#5）
70	3	**Status:** ~~in-progress~~ → **blocked-on-owner**（agent-ticket70-d 本轮承诺的 4 项已全部交付；剩下每一件事的判据都不在代理手里：`internal/risk` 的 ubuntu 归属 / `staticcheck` 的 pin-vs-35-findings 顺序 / 放行 `internal/config` 那一行注释 / push 之后才能取到的 AC#6 五 job 全 pass）
71	3	**Status:** in progress（我所有路径内的部分已落地：台账 + runner + CI 三处接入 + A43；AC#1/#2 的
72	3	**Status:** review（修复与用例已入库 `f1033e1` + 插单回执 commit；AC#1/#2/#3/#5 本机可证已勾，**AC#4 的 runner 半边与 AC#6 的改名都不在代理手里**）
77	10	**Status:** **blocked（09-23，owner 把前端收走了）** / 上一状态 **in progress**（AC#2 / AC#5 / AC#7 已交 `63ef895`；AC#1 PARTIAL）/ 原 claimed（2026-09-21 14:0x，`agent-ticket77` 已认领并开工，
85	3	**Status:** open（**排队**：本票要改 `ci.yml` 与**全仓**多个包，必须等写码槽位空出来，
86	3	**Status:** open（排队）
91	3	**Status:** open（**只读评估票**：不写生产码、不改配置，只产出一份决策备忘录）
98	3	**Status:** open（2026-09-21 15:4x 编排者建；根因由 `acceptor-ticket87` 挖出，它把它当"环境事实"登记了，
