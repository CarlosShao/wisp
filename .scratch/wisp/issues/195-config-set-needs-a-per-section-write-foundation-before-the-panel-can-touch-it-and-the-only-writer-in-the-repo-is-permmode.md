# 195 — 面板要能改配置，但"按节写入"这层地基根本不存在：先建它，才轮到 `config.get`／`config.set` 进白名单

- Status: **blocked（排在票 33 片 B 交件之后；地基未建之前不许接面板口）**。⚠ 本票**不是**新造功能请求：来路＝owner 09-28 已裁的「方法名册按当前规格补，别搞债务了」（逐字 `A381`），`config.get`／`config.set` 逐字在 `docs/specs/SPEC-08-ui-ball-panel.md:168`；**本票只是把那一支的落地顺序改成"先地基、后面板"**。
- 来路：`194-c1`（只读代价表，`docs/evidence/s1/194-method-roster-census-c1.md`，台账 `A384` 裁定②）现量出——**`internal/config/` 里今天唯一现成的写手是 `Manager.SetPermissionMode`（`permmode.go:64`）**，没有"按 section＋key 写入"的通用面。
- 为什么必须单立一票：把那枚函数当捷径接进面板的 `config.set`，**做出来的正好是本仓明令禁止的那一形**——"面板直接改权限状态"（`docs/specs/SPEC-06-security-gatekeeping.md:19`「「允许」只接受原生侧来源……面板只能「拒绝/查看」」＋ R20 那条口径「面板里的档位与工作区是**权限输入口**：只许显示＋发起请求」）。
- 关联：**票 194**（名册对齐，堆3 那一格由本票接管）· **票 33**（入向那一跳，`internal/panel/composer_dispatch.go`）· **票 128／132**（配置与安全节的既有门）· **`Q-49` 丙**（面板侧来源的结构判据形状，可复用）· **票 101**（档位持久化）。

## 现量（锚 `ba65b29f`，尺可复制）

| 事实 | 读数 | 尺 |
|---|---|---|
| `internal/config/` 里的导出写手 | **只有 `SetPermissionMode` 一枚** | `grep -rn "^func (m \*Manager) Set" --include=*.go internal/config/ \| grep -v _test.go` |
| "按节写入"地基 | **不存在**（无任何 `Set(section, key, value)` 形状） | `grep -rnE "func .*(SetConfig\|WriteSection\|SetSection)" --include=*.go internal/ \| grep -v _test.go` |
| 规格要求 | `SPEC-08:168` 逐字：「安全节放宽走 L2 重新确认（SPEC-03 §4.2）」 | `sed -n '168p' docs/specs/SPEC-08-ui-ball-panel.md` |
| 白名单那份代码 | `bridge.go:42-45` 恰好 4 枚，**无 `config.*`** | `grep -nE '=\s*"panel\.' internal/panel/bridge.go` |

## 要建什么（只做地基，不碰面板口）

给 `internal/config/` 一枚**按节写入**的入口：能落 `config.toml`、能区分**安全节**与**其余节**两条回路——其余节直接写并记审计，**安全节的"放宽"一律转进 C18 审批队列（不落盘）**，安全节的"收紧"直接写。⚠ 生效级别（`PLAN.md` D36 那三档）由既有配置模型决定，本票不改 D36 一字。

## 验收判据（每格须由**非实现者**裁）

- [ ] **AC#0 前提现量对得上**：上表四条读数在实现程起手时逐字复跑；任何一条不符 ⇒ 停下报回，不许按票面硬改。
- [x] **AC#1 分节写入面存在且被生产调用者用**：新入口在 `internal/config/` 有真实现，`grep` 到**至少一枚非测试调用方**（装配根 `cmd/wisp/` 或 `internal/` 内某条真回路）；**只有测试听众＝不算**（`A205`／[[verification-blind-spots]] 第 29ⓑ 那一族：注释里声称的钉、尺本身是空的）。
- [x] **AC#2 安全节"放宽"不落盘、进队列**：造一发卡"请求放宽安全节某项"⇒ 断言①`config.toml` 里那一项**值未变**，②审批队列多一条待决项，③审计写了归因。反向那发（收紧）必须**直接写成功**——否则这一枚会把 fail-closed 方向一起禁掉，撞 owner 那条「别把批准和次生收紧绑一体」。
- [x] **AC#3 面板来源不许冒充原生**：安全节的**放宽**路径对"来源＝面板"结构性拒绝，判据形状沿用 `Q-49` 丙那批（`internal/panel/l2_grant_boundary_test.go` 的钉形），**不许新增"面板本次允许"这类控件或出口**。
- [x] **AC#4 `SetPermissionMode` 没被冒充分节写**（票面第一句那条捷径＝本格的反向判据）：若实现里 `config.set`／安全节写入最终落到 `permmode.go:64` 那枚函数上 ⇒ 判据红。尺：`grep -rn "SetPermissionMode" --include=*.go internal/ cmd/ | grep -v _test.go` 的命中集**起手名册＝终态名册**（新增命中逐枚点名理由）。
- [x] **AC#5 门禁四数＋`gofumpt`**：`sh scripts/d22scan.sh` clean；`bash .scratch/wisp/probes/154/gate-clauses.sh` 的 **BAD 腿名册只 `G6neg`**（在册常红）；`go test -count=1 ./internal/config/` 全绿；`"$(go env GOPATH)/bin/gofumpt" -l` 空。⚠ **别拿 `ban #8 internal/` 那格的文件枚数当不变量**（`A384` 已更正：那是分母，不是违规数）。
  - **[2026-11-09 11:1x 更新｜票 411 AC#6 腿落地；本行为追加，原句不改]**：gate-clauses.sh 修好 argv 溢出（`-- "$@"` 两处，判据一字未动，票 411）后，仪器跑出史上第一次全树真读数 **rc=9**（锚 `cb2287a5`，证据 `docs/evidence/s1/gate-clauses-argv-overflow.md` §3）——**BAD 名册＝G1/G1b/G2/G5/G5pos/G5neg/G6neg/G7/G7neg 九腿**；原句「BAD 腿名册只 `G6neg`」是溢出期假读数（10 腿假 `git-grep-rc-126`、聚合恒 rc=10），作废。九腿具名命中零涉 `internal/config/**` 与本票交件四文件（本腿 scoped 复量六词表全 0，`docs/evidence/s1/195-ac5-refresh.md`）⇒ **本票射程内 0 红**；本格勾选留裁决者（本腿不翻框）。
- [x] **AC#6 本票没顺手接面板口**：`bridge.go:42-45` 的**四枚顺序与命名一字不动**，且本票不新增任何 `panel.*` 或 `config.*` 白名单常量——**名册那两枚要等本票地基绿了再按票 194 堆1 单独入账（每枚一枚 `A##`）**。

## 禁区（动到即判失败）

- `frontend/**`／`design/**` **零写面**（读可以，写不行，也不许计入任何"零命中"宣称）。
- `allowlist.txt`／`thresholds.go`／golden／C18 审批超时常量／C17 白名单**既有名字**一字节不动。
- Git 纪律：只 commit 不 push；commit 必带**显式 pathspec**；禁 `git add -A`／`.`／`-a`；禁 `--amend`／`reset`／`rebase`／`stash`／`checkout .`／`restore`／`clean`／`worktree`；仓内零删除命令；临时件只建不删。
- 工作树里躺着别人的脏件（`design/**` 那批未提交删除、`probes/152`、`probes/161/r6/logs/flip-*`、在飞程的文件）——**不还原、不提交、不删**。
- 只读闸门（若派普查程）：**终态必须等于起手名册（逐枚具名差集为空）**，不写"必须为空"（共享树里做不到，`A374` 已入账）。

## Progress log

- 09-28 16:4x 编排者立票：来路＝`194-c1` 的堆3 那一格＋`A384` 裁定②。**Owner 本轮没有新增批准**：他批的是"名册按规格补"，本票只把 `config.set` 那一支**排到地基之后**，并把"不许用 `SetPermissionMode` 冒充分节写"写成常驻判据。**若他要求现在就接面板口**，口令＝「config.set 直接接」——届时必须先摊出接完之后用户能点出来的那一项会改到哪些安全节的哪个键，并逐枚落 `A##`。


## 10-02 09:3x 编排者裁定（`195-a1` 交件后落；上面那张「现量」表**原话不改**，就地打旧）

**0. 现量表两格已不符**（腿报回、我抽验复认——我现跑 `grep -rn "^func (m \*Manager) Set" --include=*.go internal/config/ | grep -v _test.go` ⇒ **8 行**）：
- 导出写手＝**8 枚不是 1 枚**（`settings.go` 六枚具名 setter＋`allowdirs.go:109`＋`permmode.go:70`）。
- "按节写入地基不存在"那把尺只否掉了 `SetConfig|WriteSection|SetSection` **三个写死的函数名** ⇒ 是**由尺的拼写造出的事实，不是盘上事实**。地基**已部分存在**（具名叶子键＋写前 `validate`＋回执报实际落盘键路径）。
⇒ **AC#0 这一格永远勾不了**（它判的是"起手对得上"，两条已不符）。它的实际作用＝把票头换成下面 §1 那枚真缺口；下一枚落地腿起手**改复跑 §1 的名册**，⛔ 不必再复跑那四条。

**1. ★ 真缺口定形（本票剩下只做这一件事）**：把**"这枚键算放宽还是算收紧"的方向裁决**从 reload 侧／`cmd` 侧搬进**程序化写入口**。凭据：两枚锁定族放宽键 `fs.delete_enabled`／`fs.reparse_point_exceptions` 不在 `unwiredKeys` 名册、有真消费者，而 `internal/config/manager.go` 已认成 loosen（我现读 `:459` 与 `:461`/`:463`；`195-a1` 报的是 `:459`/`:460-461`，**引它的行号前先自跑**）⇒ **今天唯一的防线就是"没有写入者"，而本票原本要建的正是那枚写入者**。
⇒ ⛔ 本票**不是**"建写入地基"（248 已建），⛔ 也**不是**"把放宽挡住"（那已有两层现成仪器：`ui.go` 的 `PanelAPI` 无 `Allow`＋`config_handlers.go` 整族拒）。

**2. 裁 N1（会打死排程那枚）＝⛔ 不许动面板那条整族拒绝。**
`internal/panel/config_handlers.go` 的 `lockedFieldFamilies`（我现读 `:108`，内含 `"fs."` **前缀** ⇒ 每一枚 `fs.*` 都在抵达写入腿之前被拒）与钉它的 `internal/panel/config_route_248_test.go:215` 起那张表（`ErrConfigFieldRefused`＋`leg.applyCalls != 0` 即红）**是票 248 已裁并落地的正确形状**——出处不是那枚测试，是两处更硬的：`AGENTS §1.2` 禁止清单里那一枚「由面板侧来源的 L2『允许』」，以及**本票自己的 AC#3**（安全节放宽对"来源＝面板"结构性拒绝）。
⇒ **AC#2 的判据载体就此定形**：那一发"请求放宽安全节某项"的卡片**在 `internal/config` 那层造**（程序化入口直调），⛔ **不经过面板 `config.set`**。"fs 放宽 ⇒ 进审批队列"这一格**在面板这条路上今天不做**，因为面板那侧的出口本身就是被禁的那一形。
⛔ 不许为凑 AC#2 去改 `config_route_248_test.go` 的 `mustSay` 子串或那张表（＝拔 248 的钉子）；⛔ 不许把 `lockedFieldFamilies` 里的 `"fs."` 拆窄。

**3. 裁 N9＝方向裁决必须在拿锁之前做完，入口一律不阻塞。**
凭据：`manager_223_test.go:55` 钉"钩子必须在 `mu` 之外跑"，而 `settings.go` 的 `writeOneKey` 起手就 `m.mu.Lock()` ⇒ 落地形**只能是**：`锁外判方向 → 放宽：落一条待决审批＋记审计＋立即返回"未写入"；收紧：才进 `mu` 写`。⛔ 不许"在入口里等一张卡回来"——那是**死锁不是死红**，还会把 C18 那 300s 审批超时拖进写路径。

**4. `195-a1` §6 那六枚 J 我裁完**：
- **J4**（`A487` 的具名解冻覆不覆盖 195）＝覆盖形不变、**口径我更正**：`internal/panel/l2_grant_boundary_test.go` 那两枚锚今天是**〔只许按能力形改〕**，⛔ 不是〔一字不许改〕——我在票面与派单里写过"一字不许改"，那句按现量**过期**（该文件已被 `0d87a681` 改过 179 行）。本票落地腿**不许新增对该文件的引用**；要碰就停手上报。
- **J5**（两枚读 `frontend` 的钉）＝**不许开**。本票射程零前端；任何"零命中"宣称不许把 `frontend/**` 算进去（**既不读也不引**）。
- **J6**（落点在 `config` 还是 `cmd`）＝**落 `internal/config`**：方向判定的现读逻辑已在那层（`manager.go`），`cmd/wisp` 只做装配根接线；⛔ 不许新开 `tools→panel` 那类依赖边（既有裁定：装配根是唯一接缝）。
- **J1／J2／J3**（两包绿不绿／门禁四数／AC#2 那一发实测）＝**我此刻不判**：那是落地腿的读数＋`195-v1`（非实现者）的判语。⚠ 我此刻**没有**任何一包的整包读数——`cmd/wisp`＋`internal/config` 正被 `198-r1` 占着，⛔ 谁都不许在那面上跑整包。

**5. 排程**：`195-r1` **按住**，起跑判据两条——① `198-r1` 退出且 `git status --porcelain cmd/wisp internal/config` 为空；② AC#1 的"至少一枚非测试调用方"已由 `cmd/wisp/panel_config_store.go` 的 **7 处**满足（我现量），所以落地腿**只交 §1 那枚方向裁决＋AC#2 那发载具**，⛔ 不许重造 `settings.go`（会造出第二真相源）。

- [2026-11-09 09:4x +08] agent=收件r1（编排者派·非实现者）did=复跑读数 verdict=**AC#1/AC#2/AC#3/AC#4/AC#6 翻〔成立〕；AC#5 具名不翻；AC#0 照 10-02 裁定永缺**：交件 `e6791e6bb`（11-05 14:40 死腿遗产代提，4 文件 +909/−50：`sectionwrite.go`＋`sectionwrite_195_test.go`＋`manager.go`＋`settings.go`）三把尺真在；导出树（`cda38fd52`）复跑——`go test -count=1 ./internal/config/` 整包 **ok 2.419s**，十枚 AC 命名用例逐名 PASS（`TestAC2LooseningIsQueuedAndNeverWritten`／`TestAC2TighteningIsWrittenWithoutACard`／`TestAC2LooseningWithNoQueueWiredFailsClosed`／`TestAC3PanelLooseningOfALockedSectionIsRefused`／`TestAC3PanelTighteningIsNotSwallowed`／`TestAC4PermissionModeLooseningIsQueuedNotApplied`／`TestSectionWrite*` 四枚）；**AC#1**＝`Manager.WriteSection`（`sectionwrite.go:168`）真实现＋非测试调用链现跑在（`cmd/wisp/panel_config_store.go:210/220/227/234/241/243/348` 七处 → settings 六枚 setter → `writeOneKey`（`settings.go:245` 改装为 WriteSection 的薄适配，是收拢成单一真相源、非造第二个）；**AC#2/AC#3**＝方向裁决锁外判（N9 兑现：`lockedDirectionOf` 在 `mu` 之前，§(4) 注释逐字）、放宽：面板来源 `ErrPanelLoosenRefused` 结构性拒（审计记 `refused-panel`）、原生来源 `LoosenQueue.QueueLoosen` 进队列**不落盘**（`PendingID` 回执＋审计 `queued`）、无队列 fail-closed（`ErrLoosenQueueUnwired`）；收紧 `commitSectionWrite` 直接写＋事后审计——三味断言（值未变／队列多一条／审计归因）全由测试兑现；**AC#4**＝`SetPermissionMode` 名册起手 8 ↔ 终态 9（唯一新增命中＝`sectionwrite.go:15` 纯注释，具名理由＝新文件头自述与 permmode/AddAllowedDir 的边界；实现零路由到 `permmode.go:70`）；**AC#6**＝交件 `git show --stat` 只动 internal/config 四文件、`bridge.go:42-45` 零触碰、零新增 `panel.*`/`config.*` 白名单常量。**AC#5 不翻（具名）**：四把尺里 `gate-clauses.sh` 今天全编队跑不出真读数——我亲证导出树复跑聚合退码＝**10**（`.go` 未引号展开超 Windows 命令行上限，读数件 `.scratch/wisp/probes/154/r20261108/00-argv-overflow.md`），那一条判据无法取真值；其余三把绿（`d22scan.sh` rc=0 clean／config 包 ok／`gofumpt -l` 四文件空）。⚠ 主树有一枚在飞未跟踪件 `docs/evidence/s1/195-r1-direction-verdict.md`（别人 195-r1 方向，本腿未读未动，不构成收件凭据）。

- [2026-11-09 11:1x +08] agent=411ac6-refresh腿（编排者派·195 侧写腿）did=票 411 AC#6：受 gate 修复影响的 195 AC#5 文字刷新 verdict=**AC#5 原文一字未动、就地追加具名更新一行**（真读数 rc=9、九腿 BAD 名册、溢出期原句「只 G6neg」作废）；九腿逐枚归别票族（G1/G1b＝manifestslot/plugin、G2/G5/G5pos/G5neg＝cmd/wisp resident、G6neg＝tools 测试、G7/G7neg＝resident/plugin/manifestslot），**scoped 复量 internal/config 六词表全 0 ⇒ 本票射程内 0 红、判语「AC#5 具备翻框条件」**留裁决者（本腿不翻框）；证据件 `docs/evidence/s1/195-ac5-refresh.md`（新建）；d22scan rc=0 clean（451 产文件）；只 commit 不 push。AC#0 永缺照 10-02 裁定不变。

## 编排者裁定（2026-11-09 11:4x）

**AC#5 翻勾**：凭 `docs/evidence/s1/195-ac5-refresh.md` —— gate 修好后真读数 rc=9，九腿 BAD 逐枚命中全在 `internal/plugin/manifestslot/*`、`internal/plugin/*`、`internal/tools/*`、`cmd/wisp/resident_*`，与本票交件 `e6791e6bb` 四文件（全在 `internal/config/`）**零交集**；六把腿词表 scoped 到 `internal/config/**/*.go` 复量**六尺全 0 命中** ⇒ **本票射程内 0 红**。原句「BAD 名册只 G6neg」是溢出期假读数，作废（已在 AC#5 格下追加具名更新行，原文一字未动）。

**AC#0 保持未勾（永缺，非欠账）**：AC#0 是「实现程起手时」的前置条件复核格——**时序约束**，实现已由 `e6791e6bb` 落地，无法回溯满足。按 10-02 裁定**永久留空**，⛔ 伪造一个「起手读数」去勾它。⇒ 本票为**终态（terminal）而非 `-done`**——与票 216 AC#5 同一处置形：格格有终态、无欠账、不改名（README:15 要求全勾才可改名）。

**⇒ 票 195 结论：6 格已勾 ＋ AC#0 具名永缺 ＝ 终态。**