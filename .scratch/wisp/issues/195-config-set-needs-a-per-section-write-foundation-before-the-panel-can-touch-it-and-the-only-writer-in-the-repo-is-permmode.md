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
- [ ] **AC#1 分节写入面存在且被生产调用者用**：新入口在 `internal/config/` 有真实现，`grep` 到**至少一枚非测试调用方**（装配根 `cmd/wisp/` 或 `internal/` 内某条真回路）；**只有测试听众＝不算**（`A205`／[[verification-blind-spots]] 第 29ⓑ 那一族：注释里声称的钉、尺本身是空的）。
- [ ] **AC#2 安全节"放宽"不落盘、进队列**：造一发卡"请求放宽安全节某项"⇒ 断言①`config.toml` 里那一项**值未变**，②审批队列多一条待决项，③审计写了归因。反向那发（收紧）必须**直接写成功**——否则这一枚会把 fail-closed 方向一起禁掉，撞 owner 那条「别把批准和次生收紧绑一体」。
- [ ] **AC#3 面板来源不许冒充原生**：安全节的**放宽**路径对"来源＝面板"结构性拒绝，判据形状沿用 `Q-49` 丙那批（`internal/panel/l2_grant_boundary_test.go` 的钉形），**不许新增"面板本次允许"这类控件或出口**。
- [ ] **AC#4 `SetPermissionMode` 没被冒充分节写**（票面第一句那条捷径＝本格的反向判据）：若实现里 `config.set`／安全节写入最终落到 `permmode.go:64` 那枚函数上 ⇒ 判据红。尺：`grep -rn "SetPermissionMode" --include=*.go internal/ cmd/ | grep -v _test.go` 的命中集**起手名册＝终态名册**（新增命中逐枚点名理由）。
- [ ] **AC#5 门禁四数＋`gofumpt`**：`sh scripts/d22scan.sh` clean；`bash .scratch/wisp/probes/154/gate-clauses.sh` 的 **BAD 腿名册只 `G6neg`**（在册常红）；`go test -count=1 ./internal/config/` 全绿；`"$(go env GOPATH)/bin/gofumpt" -l` 空。⚠ **别拿 `ban #8 internal/` 那格的文件枚数当不变量**（`A384` 已更正：那是分母，不是违规数）。
- [ ] **AC#6 本票没顺手接面板口**：`bridge.go:42-45` 的**四枚顺序与命名一字不动**，且本票不新增任何 `panel.*` 或 `config.*` 白名单常量——**名册那两枚要等本票地基绿了再按票 194 堆1 单独入账（每枚一枚 `A##`）**。

## 禁区（动到即判失败）

- `frontend/**`／`design/**` **零写面**（读可以，写不行，也不许计入任何"零命中"宣称）。
- `allowlist.txt`／`thresholds.go`／golden／C18 审批超时常量／C17 白名单**既有名字**一字节不动。
- Git 纪律：只 commit 不 push；commit 必带**显式 pathspec**；禁 `git add -A`／`.`／`-a`；禁 `--amend`／`reset`／`rebase`／`stash`／`checkout .`／`restore`／`clean`／`worktree`；仓内零删除命令；临时件只建不删。
- 工作树里躺着别人的脏件（`design/**` 那批未提交删除、`probes/152`、`probes/161/r6/logs/flip-*`、在飞程的文件）——**不还原、不提交、不删**。
- 只读闸门（若派普查程）：**终态必须等于起手名册（逐枚具名差集为空）**，不写"必须为空"（共享树里做不到，`A374` 已入账）。

## Progress log

- 09-28 16:4x 编排者立票：来路＝`194-c1` 的堆3 那一格＋`A384` 裁定②。**Owner 本轮没有新增批准**：他批的是"名册按规格补"，本票只把 `config.set` 那一支**排到地基之后**，并把"不许用 `SetPermissionMode` 冒充分节写"写成常驻判据。**若他要求现在就接面板口**，口令＝「config.set 直接接」——届时必须先摊出接完之后用户能点出来的那一项会改到哪些安全节的哪个键，并逐枚落 `A##`。
