# 224 — **"本次会话内允许"这一格今天零执行者**：`SessionID` 只活在 DAO 里、生产不铸造会话身份、`approval_grant` 没有生产写手

- Status: **待派（前置已解＝09-29 13:13 批准记录 `A435`：会话＝什么 已定死、只读 seam 已批，撤销口令「撤 224 会话定义」；`PLAN.md` 与 `docs/specs/**` 一字未动）**。来源：非实现者设计复核 `219-v0`（`docs/evidence/s1/219-approval-reply-design-adversarial.md`，27,833 字节）第 3 条＋编排者复跑。
  ⚠ **下面"落点两支"那一节的取舍已按 `A435` 定案**：**走乙形、收窄成"宿主随机铸造、不许派生、结束点今天＝进程退出"**（⇒ 与 `A424` 第 2 条边界"随进程死"一字不差）；**"面板给一枚'提前结束会话'的按钮"与 `grants.list`／`grants.revoke` 两枚方法都不在本票射程**（前者属界面侧、由 owner 自己带；后者会被 `l2_grant_boundary_test.go` 的禁名词表直接判死）。原两支分析保留供追溯。
- ⚠ **它挡住的不止票 219 的第二枚按钮**：`docs/PLAN.md:1536` 把"本会话内允许"当作**今天已有的三种硬编码确认动作之一**来写——**而这一格今天根本没有执行者**。也就是说：**冻结文字里"已存在"那一档，实建状态是〔建了但没接〕**（表在、DAO 在、没人写没人读、身份没人铸）。

> ⚠ **09-29 11:5x 普查腿 `224-c1`（`.scratch/wisp/probes/224/c1/census.md`，86,894 字节／607 行，提交 `7096470d`）交回，编排者逐条现跑复认（台账 `A426`）。四格新事实改变了本票的落点与代价：**
> ⑥ **铸造者的家其实早就划好了，只是没建**：`internal/session/` **只有 `doc.go`**，`:10` 写着要管「per-session grants ledger (D45-2)」、`:16` 挂着 `DEFERRED(SessionScope): implemented by ticket 28`。⇒ 本票乙形那句"要新增一枚契约面"**要收窄成**：面已存在（C31／`internal/session`），缺的是实现与接线；但**"会话＝什么"这一问仍属待人定案**（`doc.go` 只划边界、没定义结束点）。⚠ 顺带撞出的"标记与 `SPEC-12 §5` 不 1:1"不归本票＝票 225。
> ⑦ **写手／读者／匹配器三件今天全为零**：`grep -rn "InsertGrant\|RevokeGrant\|ListGrantsBySession" --include=*.go internal/ cmd/ | grep -v _test` 只剩定义行与注释；另有三件**根本不存在**——`Answer` 词表里没有会话档、`ToolRequest` 里没有会话位、**全仓零枚路径模式匹配器**（`approval_grant.pattern` 没有任何代码比对它）。⇒ **"存哪／按什么匹配／什么时候查"三件套里，"匹配"那一枚是从零起。**
> ⑧ **判定链最小落点今天拿不到 store**（写腿的硬墙）：要加"命中 grant 就不问"这一步，位置在 `internal/tools/bridge.go` 里，而 `memory.Store` 只在 `cmd/wisp/run.go` 两处存在（`perm`／`approval`／`tools` 三包**零引用**）⇒ **必须新增一枚只读 seam（照现成 `ModeSource` 那形），这是新契约面、要先落 `A##`**；普查同时排除了"塞进 `route` 的形参"那写法（`bridge.go` 里那句 `sil is a parameter, not a lookup` 逐字挡着）。
> ⑨ **冻结钉侧两枚硬墙**（我复跑过词表与守卫，都在 `internal/panel/l2_grant_boundary_test.go`，**那枚文件一字不许动**）：`SPEC-08:171` 的 `grants.list`／`grants.revoke` 若实现成**Go 应答的入向方法**会被禁名词表**直接判死**；且 registry 腐蚀守卫使"在 `internal/panel` 新增可解码结构体"这条路**在不动冻结件的前提下走不通** ⇒ **落点必须避开 `internal/panel`**，答复字段只叫 `reason` 安全。
> ⑩ **AC#4 今天一件都判不了红**：普查给出的三件事我已复认其必要性——**①调生产铸造函数跨两次启动断两枚 id 不等／②用生产 id 写、再用第二次生产 id 查断 0 行／③走真 `Bridge.Execute` 断"问与不问"**。⇒ **AC#4 的措辞按这三件重写，缺一件就不算判据**（本仓"恒真判据"的老形，别再交一次）。

## 现量（起手逐条复算，别信这里的行号）

| 事实 | 读数 | 尺 |
|---|---|---|
| 表与 DAO 在 | `internal/memory/dao_misc.go`：插入时校验 `g.Tool == "" \|\| g.Pattern == "" \|\| g.SessionID == ""`（`:21` 一带）、写参 `:35`、扫描 `:120`；模型 `internal/memory/models.go:98` `SessionID`，`:104` 注释逐字 `GrantScopeSession` 是 SPEC-02 §3 **唯一**的 scope 值 | `sed -n` 逐行 |
| ⛔ **生产零写手／零读者** | 那几枚 DAO 方法的**非测试调用者＝0**；`internal/agent/journal.go` 的 `DecisionAllowGrant` 决策类型**零写者** | `grep -rn 'InsertGrant\|ListGrantsBySession\|DecisionAllowGrant' --include='*.go' internal cmd \| grep -v _test` |
| ⛔ **生产不铸造会话身份** | `grep -rn 'SessionID' --include='*.go' internal cmd \| grep -v _test` 的命中**全在 `internal/memory` 的 DAO／模型自己**（`:21`／`:35`／`:120`／`models.go:98`），**没有任何一处生成或传递 session id** | 现跑 |
| 判定链只认最严 | `internal/risk/mode.go`（`resolveMode` 一带）把"档"判成上界，**没有"这一条 grant 命中了吗"这一步** | 现读；枚数与行号请续腿自己重取 |
| ⚠ **守卫不是三枚，是四枚**（09-29 由 `224-c1` 顶回、编排者现跑复认） | `internal/perm/ticket90_persist_test.go` 里 `TestTicket90SessionGrantDoesNotSurviveRestart` 用的是**测试自己写死的两枚字符串**（`session-before-restart`／`session-after-restart`，断言在 `:234/:244/:258`）：断"新 session 读到 0 行"，同一函数又断那一行**还活着**（留给审计）。**但还有一枚同族守卫在 `cmd/wisp/run_mode101_test.go:389` `TestTicket101SessionGrantDoesNotCrossRestart`**（`:411` 用 `memory.GrantScopeSession` 写行、`:436` 先断"fixture 那行必须是活的，否则以下断言全为空"、重启后按**另一枚写死的 session 字面量**查）⇒ **准确说法＝"四枚全绿而永久通行证实际生效"是可能的**，因为**四枚都不看生产怎么铸 id**。⚠ 但也**不许把这四枚说成全瞎**：`run_mode101_test.go:432-441` 对"匹配器把 session 条件丢了"这一形**是有射程的**（它先断同一 session 读到 1 行、再断新 session 读到 0 行）。文件头逐字 `They are three test functions on purpose`／`a single "persistence" case is how a permanent免审通行证 gets in` | 编排者 09-29 11:5x **自己 `sed` 逐枚读过这四处断言**（`A426`）；续腿动手前仍要自己重跑取现行号 |

## 为什么不能顺手"生成一个 id 就行"

**这正是上面 ⚠ 那一枚坑**：如果 session id 是按进程／按任务／按时间派生出来的，"本会话内允许"就会**事实上变成永久免审通行证**——那撞 `internal/perm/store.go:19-27` 的边界、`docs/PLAN.md:1642`（会话结束后授权**必须**失效），也撞那三枚故意分开的钉。**所以本票的第一问题是"什么算同一次会话"，不是"怎么写库"。** 这一问法**属于待人定案**（见"落点两支"），未定义即停。

## 落点两支（**都要写判据，我倾向乙，但这一支要 owner 点头**）

- **甲：把"会话"定义成一次 `wisp run` 进程**（最省事，但见上面的坑——它让"重启后仍在"与"进程活着就仍在"变成同一件事，安全轨会被稀释）。
- **乙（我倾向）**：**会话身份由宿主显式铸造并有明确的结束点**（例如球／面板的一次交互会话＝从打开到关闭，或一个有界时长），**结束即撤销**；授权匹配只在"身份相同且未撤销且未过期"三条件齐全时才成立。⚠ 这一支**要新增一枚"会话身份从哪来"的契约面**（C17／SPEC-02 §3 的 `approval_grant` 形状都只写了 `session` 这个 scope 值，**没写身份从哪来**）⇒ **先落 `A##` 请 owner 批"会话＝什么"**，再动码。

## 判据（草稿，逐格要 `file:line` 与正控）

- [ ] **AC#1 身份有唯一的铸造点**：`grep` 点数：改前生产铸造者＝0，改后 ≥1，且具名写出"会话结束"是哪一行代码在做什么。
- [ ] **AC#2 三件套齐全**：**写**（答复"本会话内允许" ⇒ `approval_grant` 真有一行，带 session 身份＋工具＋模式＋创建时间）、**读**（下一次同类请求命中它 ⇒ **不再弹卡**）、**失效**（会话结束 ⇒ 同一身份查不到、行为回到"要问"）。三格各自有独立用例，⛔ **不许合并成一枚"持久化"用例**（`ticket90_persist_test.go` 文件头逐字"They are three test functions on purpose"）。
- [ ] **AC#3 重启必失效**：进程重启 ⇒ 同类请求**重新弹卡**（这条是"永久免审通行证"的命）。
- [ ] **AC#4 反控（正控！）**：种一条"派生式 session id 让授权跨会话仍生效"的假腿 ⇒ **AC#2／AC#3 至少一枚必须红**；如果三枚钉全绿而授权实际生效 ⇒ **本票判失败**，因为那正是 `store.go:19-27` 那段注释点名要挡的形状。
- [ ] **AC#5 不许把"长期"混进这一票**："长期／永久允许"落 `[fs] allowed_dirs`，见票 219 与批准记录 `A424`，**不是这张票的射程**。

## 禁区

不动 `docs/PLAN.md`／`docs/specs/**`／`thresholds.go`／golden／`allowlist.txt`；不动 `internal/perm/ticket90_persist_test.go`（三枚钉各管一头）／`internal/panel/l2_grant_boundary_test.go`（冻结件，且它的词根表禁 `allow`／`grant`／`decision`／`verdict` 等十二词×两拼 ⇒ **新字段只叫 `reason` 安全**）／`tokens_fourway_test.go`；不新增 C17 方法名；`frontend/**`／`design/**` 零写零转述；⚠ 与票 201／222／223 同撞 `cmd/wisp`／`internal/agent` ⇒ **串行**。

### 2026-09-29 23:5x 编排者增量 — 撞钉预检（腿 `nail-1`，台件 `.scratch/wisp/probes/nail/1/precheck.md`＝**38,240 字节**）抓到**一枚与本票 AC#2 极性相反、今天还绿着的常驻钉**：⛔ 不写清这条就派 224 写腿，必然造出一枚假绿或一枚真红

- **那枚钉，我现读复认（锚点 `98df640a`）**：`cmd/wisp/run_mode101_test.go`——`:400` 有字面量 `const session = "session-before-restart"`，`:427-441` 的注释与断言逐字为「Refused with the grant live IN ITS OWN SESSION, too: the assembly hands the bridge a mode, **never a grant source** (`Options.Confirmations` nil)」＋ `if !firstErr { t.Fatalf("boot 1 already let the B-tier .env write through …") }`。⇒ 它把"**带着生效中的本会话授权，同类写入仍必须被拒**"钉成了期望；而本票 **AC#2（`:36`）要的是相反的方向**（答复"本会话内允许" ⇒ 命中授权就不再弹卡）。
- ⚠ **对腿那句"今天绿纯属侥幸"的更正（我把它写成两件事，别合并）**：它绿**不是因为字面量撞不上**，而是因为**生产那一侧今天根本没有授权来源**——我现跑 `grep -rn "Confirmations:" --include=*.go internal cmd`＝**1 命中，且在测试里**（`internal/tools/ticket90_test.go:524`），生产装配零枚 ⇒ `:427-441` 那句"never a grant source"**今天是真话**。本票 AC#1 现读也写着"生产铸造者＝0"。**⇒ 两件事叠加（没有铸造者＋没有授权来源）使这枚钉成立；本票一旦落地，恰好就是把这两件都打掉**，所以这枚钉**必须与本票 AC#2 同批改期望**，否则：不接⇒AC#2 假绿；接了⇒那枚钉红。
- **我给的派单口径（写腿照做，验收腿照此判；三条一枚都不许少）**：
  1. **授权改写 `cmd/wisp/run_mode101_test.go:427-441` 的期望**（该件**不是**冻结件，可改），但那两句注释必须改成**带条件的事实句**（不许留"never a grant source"这种在改后会变成假话的绝对句），且**逐字保留它原本要防的东西**（B 档 `.env` 在无授权时仍被拒）。
  2. **加一枚锁**：断言"生产铸造的 session id **永远不等于** `session-before-restart` 这类测试字面量"——⛔ 没这条，AC#4（`:38`，反控／正控）会**假绿**：字面量若哪天被生产码复用，反控那发根本不红。
  3. **交回时两发读数**：改前该钉绿的具名读数＋改后同一枚按新期望跑的结果（红必须能还原成绿）；并 `git status --porcelain -- internal cmd` 为空。
- **相邻两枚禁字段核（本票不许顺手绕）**：`internal/tools/ticket90_test.go:431` 用反射钉死 `tools.Decision` **不许长 grant/allow 字段**、`agent.ToolRequest` **不许长 mode/allow 字段**。⇒ 本票要加的"会话身份"落在哪一枚结构上**是设计决定，不是实现细节**：写腿只能**提出落点**（`SessionID` 挂哪、谁铸、谁读），由**非实现者验收腿**判"这算不算那枚禁字段核的射程内"。⛔ 不许写腿自己边写边判、更不许动三枚冻结件（`internal/panel/tokens_fourway_test.go`／`internal/panel/l2_grant_boundary_test.go`／`internal/perm/ticket90_persist_test.go`，本增量一句也不引其内容）。
- **排程**：本票写腿落 `internal/tools/bridge.go`＋`cmd/wisp/run.go` ⇒ ⛔ **与在飞裁决腿 `235-v1` 串行**（它此刻正占 `internal/tools` 的测试与突变窗口；同包并发＝互相洗读数）。票 213／214／220 三张卡的结论是"有钉但都可用、不必改票面"，各自的界线已写进 `precheck.md`，派单时逐张抄过去。
