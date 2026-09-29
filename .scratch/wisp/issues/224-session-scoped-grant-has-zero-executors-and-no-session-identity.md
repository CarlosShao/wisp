# 224 — **"本次会话内允许"这一格今天零执行者**：`SessionID` 只活在 DAO 里、生产不铸造会话身份、`approval_grant` 没有生产写手

- Status: **待派**。来源：非实现者设计复核 `219-v0`（`docs/evidence/s1/219-approval-reply-design-adversarial.md`，27,833 字节）第 3 条＋编排者复跑。
- ⚠ **它挡住的不止票 219 的第二枚按钮**：`docs/PLAN.md:1536` 把"本会话内允许"当作**今天已有的三种硬编码确认动作之一**来写——**而这一格今天根本没有执行者**。也就是说：**冻结文字里"已存在"那一档，实建状态是〔建了但没接〕**（表在、DAO 在、没人写没人读、身份没人铸）。

## 现量（起手逐条复算，别信这里的行号）

| 事实 | 读数 | 尺 |
|---|---|---|
| 表与 DAO 在 | `internal/memory/dao_misc.go`：插入时校验 `g.Tool == "" \|\| g.Pattern == "" \|\| g.SessionID == ""`（`:21` 一带）、写参 `:35`、扫描 `:120`；模型 `internal/memory/models.go:98` `SessionID`，`:104` 注释逐字 `GrantScopeSession` 是 SPEC-02 §3 **唯一**的 scope 值 | `sed -n` 逐行 |
| ⛔ **生产零写手／零读者** | 那几枚 DAO 方法的**非测试调用者＝0**；`internal/agent/journal.go` 的 `DecisionAllowGrant` 决策类型**零写者** | `grep -rn 'InsertGrant\|ListGrantsBySession\|DecisionAllowGrant' --include='*.go' internal cmd \| grep -v _test` |
| ⛔ **生产不铸造会话身份** | `grep -rn 'SessionID' --include='*.go' internal cmd \| grep -v _test` 的命中**全在 `internal/memory` 的 DAO／模型自己**（`:21`／`:35`／`:120`／`models.go:98`），**没有任何一处生成或传递 session id** | 现跑 |
| 判定链只认最严 | `internal/risk/mode.go`（`resolveMode` 一带）把"档"判成上界，**没有"这一条 grant 命中了吗"这一步** | 现读；枚数与行号请续腿自己重取 |
| ⚠ 三枚钉量的是**身份**不是**持久化** | `internal/perm/ticket90_persist_test.go` 里 `TestTicket90SessionGrantDoesNotSurviveRestart` 用的是**测试自己写死的两枚字符串**（`session-before-restart`／`session-after-restart`）：断"新 session 读到 0 行"，同一函数又断那一行**还活着**（留给审计）。文件头逐字 `They are three test functions on purpose`／`a single "persistence" case is how a permanent免审通行证 gets in` ⇒ **一枚"每次派生"的生产 session id 可以让三枚钉全绿而通行证实际生效** | **编排者 09-29 11:2x 自己 `sed` 逐枚读过这三段断言**（原〔腿报，未复核〕已被现读推翻为"腿说得对"）；续腿动手前仍要自己重跑一遍取现行号 |

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
