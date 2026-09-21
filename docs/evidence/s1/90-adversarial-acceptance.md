# 票 90 —— 独立对抗验收（`acceptor-ticket90`）

**被验收对象**：`.scratch/wisp/issues/90-user-facing-permission-modes.md`（六框含 AC#3b，`Status: ready-for-review`）。
代码 SHA：`1d4f289`（落码）+ `582b9a9`（落测试）+ `d5564c2`（AC#3b 拆条 + 台账改口）；票面收尾 `467040b`。
**本代理不写生产码、不碰工作树生产文件、不 push。** 唯一产出=本文件。

**特殊性前提（接单时被告知）**：前任 `agent-ticket90` 撞轮数上限死时 AC#4 那枚框是**勾着的但没有一次变异读数**；
接续的 `agent-ticket90b` 补做 5 轮变异并把读数写进票面。⇒ **本验收的首要任务就是打它补写的那批数字**：
**下面每一行都是本代理自己在仓外纯净快照里重跑的原文读数，没有一处照抄它的红名或它的数。**

---

## 0. 装置与"我没在工作树里动手"的证据

| 项 | 实测 |
|---|---|
| 纯净快照 | `/tmp/wisp90acc-90c`（会话后缀 `90c`），`git archive d5564c2 \| tar -x -C /tmp/wisp90acc-90c` |
| 快照非 git 仓 | `git rev-parse --is-inside-work-tree` ⇒ `fatal: not a git repository`，**rc=128** |
| 变异全部只在快照里 | 是；见第 2、3 节每一轮的 `grep` 落地自证与还原 |
| 快照最终还原干净 | `cp` 回原件后 `diff -q` 一致，并与**第二次独立解出的** `/tmp/wisp90acc-90c-ref` 做 `diff -rq` 整树比较 ⇒ **零差异，rc=0**（含我临时加的 4 个探针文件也已删除） |
| 工作树未被本代理改动 | 开工前后各一次 `git diff --quiet -- internal/risk/mode.go` ⇒ **rc=0**。验收期间工作树里出现的 `M cmd/wisp/run.go`、`M .scratch/wisp/issues/101-*.md` 是 `agent-ticket101` 在写，**本代理一个字没碰** |
| 冻结面自查（D22） | `git show --stat` 三枚提交只动：`internal/{config/{manager,permmode,schema,unwired,validate}.go, perm/*, risk/mode.go, tools/{bridge,gate,mode}.go}`；`git diff --name-only 1d4f289^..d5564c2` 里 **`assessor.go`/`pathresolver*`/`rules_gateway.go`/`docs/PLAN.md`/`docs/specs/*` 零命中** ⇒ 没踩冻结面，不需要上报 |

---

## 1. 基线：四数 + SKIP 逐条点名 + `=== RUN` 乘法核对

命令（在纯净快照里）：
`go test -count=2 -v ./internal/risk/ ./internal/config/ ./internal/perm/ ./internal/tools/`

**四数（本代理原文读数，日志 `/tmp/wisp90acc-90c-baseline.log`，99460 字节）：**

| 口径 | 数 |
|---|---|
| `=== RUN` | **724** |
| `--- PASS`（含缩进的子测试） | **721** |
| `--- FAIL` | **1** |
| `--- SKIP` | **2** |
| 顶层 `^--- PASS`（前任用的口径，见下） | 443 |
| 包级 `^ok` | 3（config / perm / tools），`internal/risk` 是 `FAIL` |

**`=== RUN` == 2 × 不同名 核对**：`grep '^=== RUN' \| awk '{print $3}' \| sort -u \| wc -l` = **362**；**362 × 2 = 724 = 总 RUN** ⇒ **成立**（没有"只跑了一遍冒充两遍"）。

**SKIP 逐条点名（2 条，同一个名字跑了两遍）**：
- `--- SKIP: TestSyncRegistryProbeLive (0.00s)`（日志第 268、585 行，`internal/risk`）
  原因读了源码（`internal/risk/syncdirs_test.go:346-348`）：本机 HKCU 没有 registry-grade 同步盘记录 ⇒ `t.Skip("no registry-grade sync record on this machine …")`。**与票 90 无关，属票 19 家族的环境跳过。**
  ⇒ **假绿排查**：这 2 条 SKIP 在两次样本里都是同一条，**没有任何一条票 90 的用例被 SKIP 冒充**（票 90 的 15 条用例在下面每轮都以 `--- PASS` 或 `--- FAIL` 现形）。

**唯一的那条 FAIL 是本仓已登记的仪器噪声，本代理没拿它当借口、逐样本量过它**：
`TestResolvePerCallBudget`（`internal/risk/pathresolver_budget_norace_test.go:21`，C26 性能门，预算 1ms 墙钟）
- 基线两样本：**样本 1 PASS（2.33s）/ 样本 2 FAIL**，打印 `C26 Resolve: 2003748 ns/op = 2.004 ms/op (budget 1.000 ms, 758 samples)`。
- 静默单跑 3 次（`-run TestResolvePerCallBudget`）：**0.551 ms/op PASS、0.917 ms/op PASS、1.105 ms/op FAIL**（三次全报，不挑运气那次）。
- 变异轮里它也跑：M1 PASS / M3 **FAIL** / M4 PASS / M5 **FAIL**。
- **合计 9 个样本 = 5 PASS / 4 FAIL** ⇒ 它不是"负载下偶尔假失败"，是**边界性 flake，失败率接近 44%**。
  判据核对：票 90 的三枚提交**没有一处改动 `Resolve`/`pathresolver*`/`resolveBudget`**（第 0 节 diffstat + 下面第 2 节每轮变异前后它的行为一致）⇒ **记在票 86 的账上，不是票 90 的新红**。
  ⚠ 但**票 94 的验收读数（静默 7 样本 6 ok/1 FAIL）在本机当前负载下复现不出来**，本代理 9 样本 4 红 ⇒ "争用所致"这句归因**比实况乐观**，留给票 86，不影响本票判决。

**票面上那句"未变异的 `-v` 全量跑 rc=0（PASS 222 / SKIP 1 / FAIL 0）"**：本代理**没能复现成 rc=0**，
差的就是上面这条 flake（我的 `-count=2` 跑里 1 红）。**222 这个数量级对得上**（前任数的是顶层 `^--- PASS`：我基线顶层 443 / 2 ≈ **221.5–222**），
它与我用的"含子测试"口径不同 ⇒ **不是它的数造假，是两种计数口径**；本代理在下面同时报两个口径的数，避免下一个人误判。

---

## 2. AC#4：两侧变异本代理**自己重做**（5 轮红线变异 + 2 轮持久化交叉变异）

锚点全部是**承载行为那一行**；每轮**同一条 `&&` 链里先 grep 证落地**，再 `go build ./...` 单独打印 rc
（**编译失败不算变异**：本代理 7 轮变异 **build rc 全部实测 = 0**，没有一轮是靠"红"来冒充变异成功），
还原后 `diff -q` 与原件一致。**参照数**：本轮变异只跑 risk/perm/tools 三包（`-count=1 -v`）。

| # | 变异 | PROOF（同链 grep） | build rc | 测试 rc | 红名（本代理原文） | 断言原文首句 | 还原 |
|---|---|---|---|---|---|---|---|
| **M1** 去掉"不可逆"红线 | 删 `case R8:` 两行（`internal/risk/mode.go` `redLine()`） | `grep -n "case R8:"` ⇒ **rc=1（空）**；邻居 `case R5:/R9:/SessionOverrideBlocked` 计数 **4** 仍在（外科手术） | **0** | **1** | **2**：`TestTicket90ScreenTable`、`TestTicket90IrreversibleStillAsksInEveryMode` | `ticket90_test.go:335: mode=auto_approve: irreversible call raised 0 L2 cards, want 1` ／ `store_test.go:299: R8 不可逆 under auto_approve: level = L0, want L2` | `diff -q` 一致，md5 回到 `faec9e19…` |
| **M2** 去掉"污染升级"红线·证人一 | 删 `if d.SessionOverrideBlocked {…}` 整支 | `grep -n "d.SessionOverrideBlocked"` ⇒ **rc=1**；`case R4/R5/R8/R9` 计数 **4** 仍在 | **0** | **1** | **1**：`TestTicket90TaintFlagAndDenyAreNeverSilenced` | `store_test.go:319: mode=auto_approve: a tainted verdict came out L0 (silenced=true); PLAN.md:1640 says no authorization - including a permission mode - covers a C25 escalation` | 一致 |
| **M2′**（本代理**加测**的镜像刀）只删第二个证人 | 删 `case R4:` 两行 | `grep -n "case R4:"` ⇒ **rc=1**；flag 支仍在（邻居计数 **3**） | **0** | **1** | **1**：`TestTicket90ScreenTable` | `store_test.go:299`（"R4 污染"那一行）—— 与 M2 红的是**不同**用例 | 一致 |
| **M3** 两个证人同灭 | 上面两刀一起 | `grep -n "case R4:\|d.SessionOverrideBlocked"` ⇒ **rc=1**，其余 6 个 case 仍在 | **0** | **1** | **3 票 90 + 1 flake**：`TestTicket90TaintEscalationNeverSilenced`、`TestTicket90ScreenTable`、`TestTicket90TaintFlagAndDenyAreNeverSilenced`（+ `TestResolvePerCallBudget`） | `ticket90_test.go:371: mode=auto_approve: tainted call raised 0 L2 cards, want 1` | 一致 |
| **M4** 去掉"Deny 永不被静默"（= `PLAN.md:1588`/`ban #6` 的机器形状） | `mode.go:161` 的 `Silenced{Level: Deny, Kept: "拒绝级…"}` 改成 `Silenced{Level: L0, Silenced: true} // MUT-C` | `grep -n "Level: Deny, Kept"` ⇒ **rc=1**；落地行打印在 **`mode.go:161`** | **0** | **1** | **4**：`TestTicket90TierADenySurvivesEveryMode`、`TestTicket90TaintFlagAndDenyAreNeverSilenced` + **两条非本票的旧卫兵** `TestBridgeRefusesTheRealShortNameOfAnAListFile`、`TestSensitiveFileIsDeniedNotEscalated` | `ticket90_test.go:417: mode=ask_every_step: A-tier path ran the tool 1 times, want 0` ／ `bridge_junction_windows_test.go:502: 8.3 拼法绕过了 A 档拒绝` ／ `fs_test.go:173: A-list read must be denied` | 一致 |
| **M5** 红线改成"模式优先" | `Screen()` **函数体第一句**插入 `if m == ModeApprove… { return Silenced{Level: L0, Silenced: true} } // MUT-D` | `sed -n '157,163p'` 打印它确为函数体首句；`grep -c MUT-D` = **1** | **0** | **1** | **5 票 90 + 1 flake**：`TestTicket90IrreversibleStillAsksInEveryMode`、`TestTicket90TaintEscalationNeverSilenced`、`TestTicket90TierADenySurvivesEveryMode`、`TestTicket90ScreenTable`、`TestTicket90TaintFlagAndDenyAreNeverSilenced` | 五条首句逐字为 `ticket90_test.go:335`、`:371`、`:417: mode=auto_approve: A-tier path ran the tool 1 times, want 0`、`store_test.go:299`、`store_test.go:319` | 一致；末轮 `grep -c "MUT-"` = **0** 且 `go build ./...` rc=0 |

**与前任票面的对账（逐条，不平均）**：它报的 **(i)-a 红 2 / (i)-b1 红 1 / (i)-b2 红 3 / (i)-c 红 4 / (ii) 红 5**
**五轮的名单、文件行号、断言文本本代理全部独立复现，无一相符不上**。它的总数 222 是**顶层口径**（见第 1 节），不是假数。
⇒ **AC#4 那枚"补作业"写下的数字站得住。**

### 2A. 复核它那条最有价值的细节：`SessionOverrideBlocked` 删掉后真实链为什么还绿

它的归因是"第二证人 `case R4:` 还在"。**本代理不信叙述，装了探针去量"为什么绿"**：
在快照 `internal/tools/` 临时加一条探针，跑**真实链**（`tools.Bridge.Execute` + `risk.NewRiskAssessor().WithTaintDetector(...)`，`mode=auto_approve`），
打印路由那一刻 `Decision` 的形状。**探针只在快照里，收尾时删除，整树 `diff -rq` 已证清干净。**

| 状态 | 探针原文读数 |
|---|---|
| **未变异** | `PROBE mode=auto_approve level=L2 rules=[R4] flag=true kept="R4 污染升级（C25，任何授权都不可覆盖，含权限模式）" silenced=false` |
| **M2**（删 flag 支，留 `case R4:`） | `… level=L2 rules=[R4] flag=true kept="R4 污染升级（C25）" silenced=false` ⇒ 兜住它的那句话**就是 `case R4:` 的字符串**（`mode.go:210`），不是第三个原因 |
| **M2′**（只删 `case R4:`，留 flag 支） | `… kept="R4 污染升级（C25，任何授权都不可覆盖，含权限模式）"` ⇒ 镜像方向也成立，兜住的是 `mode.go:198` 那一支 |
| **M3**（两个都删） | `… level=L0 kept="" silenced=false`，且真实链用例**当场红**（`ticket90_test.go:371`） |

**判决：它的归因 CONFIRMED（本代理用 `Kept` 字符串本身把它钉死，而不是靠"绿所以是它"）。**
并且**覆盖度不是"较弱"而是"两把独立的锁"**：真实链的污染判定同时带 `RulesHit=[R4]` 与 `SessionOverrideBlocked=true`
（未变异那一行是本代理第一次把它量出来），删任一把都还有另一把，两把全拆真实链才红 ⇒ 这是**冗余防御**，不是覆盖洞。
⚠ 按 A42 的规矩说一句：**这条冗余必须被说出来**，票面已经说了（"这条绿是有原因的绿，不是假绿"），**说对了**。
⇒ **不需要"补一条用例"**；本代理反而把它加的那条 M2′ 镜像刀**留在本节当判据**，建议下一张碰 `redLine()` 的票照抄这个双向拆法。

---

## 3. AC#3b 的边界（owner 的定案 R20 里那条被点名的边界）

### 3.1 `PLAN.md:1640` 原文（本代理自己 `awk 'NR>=1626 && NR<=1645'` 读，不信票面概括）

```
1639:    - **（新增）D45 边界**：批量聚合确认**不得**把 L2 聚进去；
1640:      会话授权**不得**覆盖 C25 污染升级；会话结束后授权**必须**失效；
1641:      单次影响 ≥50 文件 → 必须升 L2（C19/R7）
```
顺带把票面另引的一条也对了原文：`PLAN.md:1628-1629` = "**（新增）A 档黑名单**：读 `~/.git-credentials` / `.git/config` /
`%APPDATA%\wisp\config.toml` → **必须拒绝，且任何授权与会话授权都无法覆盖**"。
⇒ **票面对 `:1640` 的引用是准的**（那一行确实同时载着"不得覆盖污染升级"与"会话结束必须失效"两半）。
⚠ **一处口径要说清**：契约写的是"**会话结束后**失效"，owner 的 M3 要的是"**跨重启**不回默认"。
**这两句不是一件事**，票 90 的实现把前者留给票 49（按 session id 判），把后者只做在"模式"这一个键上 —— 见 3.3。

### 3.2 三条用例拆开之后**没有一条同时依赖两件事**（前任合并过，`d5564c2` 才拆）—— 本代理用**交叉变异**验，不用 grep 敷衍

静态那半边先量：`grep -n "Mode\b" internal/perm/ticket90_persist_test.go` 限定在 **AC#3b 函数体（220–281 行）→ 0 命中**（票面那句"函数体内再无任何 `Mode` 引用"属实）。
但**"没有引用"不等于"没有依赖"** ⇒ 本代理补了两把方向相反的刀：

| 刀 | 改在哪 | PROOF | build rc | 结果（逐条点名） |
|---|---|---|---|---|
| **MUT-E**：模式**不再落盘** | `internal/config/permmode.go:73` `SaveFile(m.path, m.cur)` → `error(nil) // MUT-E no persist` | grep 打印落地行、`SaveFile(m.path` 在函数内**已不存在** | **0** | **只有 `TestTicket90ManualSwitchSurvivesRestart` 红**（`ticket90_persist_test.go:161: after restart mode = ask_every_step, want the manually chosen ask_high_risk` / `… chosen auto_approve`）；`UntouchedConfigStartsAtTheDefault`、**`SessionGrantDoesNotSurviveRestart`**、`ConfigKeyChangesWhatTheChainAsks`、`HandEditLooseningGoesThroughD36` **全绿** |
| **MUT-F**：会话授权**变成跨会话可见** | `internal/memory/dao_misc.go:81` `WHERE session_id=?` → `WHERE (session_id=? OR 1=1)` | grep 打印落地行（第 81 行） | **0** | **只有 `TestTicket90SessionGrantDoesNotSurviveRestart` 红**（`ticket90_persist_test.go:263: a session grant survived the restart (1 rows for the new session)`）；**AC#3(a)、AC#3(b) 全绿** |

⇒ **两条刀各只红一条用例，方向互不污染**：AC#3(a) 咬的是"模式必须 durable"，AC#3b 咬的是"授权必须 NOT durable"，
**谁也不靠对方的机制拿绿**。`d5564c2` 的拆分是**真拆**，不是把对照段删了个文字。**这一格本代理判 PASS，且证据强度高于票面自述。**

### 3.3 安全底线：有没有任何一条路径能让"会话授权"跨重启？（本代理把三条通道逐个量过）

| 通道 | 实测（探针在快照，已删） |
|---|---|
| **写盘**（配置通道） | `risk.blacklist_overrides` 手工写进 config.toml 后 `config.LoadFile` ⇒ **报错、拒绝加载**：`config.toml: risk.blacklist_overrides is written but does nothing … because a file on disk cannot have clicked anything. Remove the key or set it back to its default - accepting it silently would be a lying config option.` ⇒ **磁盘上不存在任何静态预授权通道**（票 83 的规矩在位）。`internal/config/schema.go` 里也**没有任何 grant/授权类键**（grep `grant` 只命中注释与 `plugins` 的能力授予，那是 D36 的另一件事）。 |
| **写盘**（数据库通道） | `approval_grant` 那一行**确实落盘**：探针 `PROBE2-C surviving row: scope=session session="dead" expires_in=3600s tool=fs.write pattern=/tmp/*`。本代理判它**不是**跨重启生效，三条实测理由：① `InsertGrant` 只认 `scope=="session"`，`scope="permanent"` ⇒ **`memory: invalid grant scope "permanent" (want session)`**（探针 `PROBE2-C`）；② `expires_at` **必填**（`=0` 直接拒）；③ 取用**只按 session id**（`ListGrantsBySession`），新会话读到 0 行（MUT-F 那把刀正是"改掉这个 WHERE 才会红"）。残留的**行**由 `retention.go` 按"失效后 30 天"清（SPEC-02 §4）。 |
| **缓存 / 内存** | `perm.Store` 只缓存**切档记录**（`Switch` 历史，`Snapshot()` 只读，`TestTicket90SnapshotIsReadOnly` 盯），**没有任何授权位**；`tools.Bridge` 的 B 档免审位来自 `Options.Confirmations`，生产**根本没注入**（`git show d5564c2:cmd/wisp/run.go:262-278` 的 Options 里没有 `Confirmations`），且 `bridge.go:178-183` 对 nil 显式返回 `nil // no file was ever confirmed: B tier stays "ask first"`。 |
| **谁能消费授权** | `grep -rn "InsertGrant\|ListGrants\|RevokeGrant\|GrantScope" --include=*.go .` 排除 `internal/memory/` 与测试 ⇒ **零命中**。**今天生产里根本没有授权消费者**（票 49 未接），所以连"误用陈旧授权"这件事都还没有出口。 |

**⇒ 结论：没有任何路径能让"会话授权"跨重启生效。** 唯一"跨重启还在"的是**审计行**（数据库里的死会话授权记录），
它取不到、作用不了、且被票 90 自己的用例**断言必须保留**（`ticket90_persist_test.go:270` "want the one old-session row kept for audit"）。
⚠ **允许残留（不判 FAIL，理由如上，且无消费者）**：AC#3b 钉的是"**按 session id 查不到**"这一形状。
若票 49 将来写一个"按 tool+pattern 匹配"的消费者（走 `ListGrants()` 全表视图），**AC#3b 不会红**（MUT-F 只改了 `ListGrantsBySession`）。
⇒ **建议补一条用例**：断言"任何消费者拿到的 grant 集合必须与 `session_id` 绑定 / 未过期"，或把这条判据写进票 49 的 AC。
**这是本文件里唯一"补一条用例"的建议，且它是票 49 的地界，不是票 90 的缺陷。**

---

## 4. 三档的默认与持久化语义：**坏输入实测落在哪一档**

票面说"fail-closed 到最严档"。本代理造了 **6 种坏输入**，全部在快照里用探针量（`PROBE-A…I`，探针已删）：

| 坏输入 | 实测落点（原文） |
|---|---|
| 档位文件**根本不存在** | `NewManager err=config: config.toml read: open …: The system cannot find the path specified.` **mgrNil=true** ⇒ 配置层**不产出 Config**，没有"悄悄回默认" |
| 文件**损坏**（非 TOML 字节） | `NewManager err=config: config.toml: cannot migrate from schema version 1: not valid TOML … the file was left untouched - fix or restore it manually, **it will never be silently reset**` mgrNil=true |
| **版本不认识**（`schema_version` 比本机大 3） | `NewManager err=config: config.toml: schema_version 5 was written by a newer build (this build understands 2)` mgrNil=true |
| 值不在词表（`permission_mode = "everything_off"`） | 加载**报错并点名键路径**：`config.toml: risk.permission_mode: permission mode "everything_off" is not one of ask_every_step\|ask_high_risk\|auto_approve (risk.permission_mode): refusing to guess`；**热路径**（绕过 validate 手搓的 `*Config`）实测 `Config.PermissionMode() = ask_every_step` |
| 类型不对（`permission_mode = 7`） | `NewManager err=… line 4, col 19: cannot decode TOML integer into struct field config.RiskSection.PermissionMode of type string` mgrNil=true |
| 键写了但是空串 | `PROBE2-B empty permission_mode: raw="" mode=ask_every_step`；`store mode = ask_every_step`（`ParseMode` 把 `""` 归到默认，见 `mode.go:84`） |
| 下游拿到 nil 的三种形状 | `perm.New(nil manager)` ⇒ `perm: New requires a ConfigManager (the mode has no other source)`；`(*Store)(nil).PermissionMode() = ask_every_step`；`(*config.Config)(nil).PermissionMode() = ask_every_step` |
| 默认档真的什么都不放过 | `risk.DefaultMode() = ask_every_step`，`Screen(L1) = {Level:L1 Silenced:false Kept: Mode:ask_every_step}`（L1 也没被静默） |

**⇒ 判 PASS。** 语义是"**两层都最严**"：装配层**响亮地拒绝对话**（进程拿不到 Config，`cmd/wisp/run.go:216` 走 `return rt, 2`），
热路径**万一被手搓绕过 validate** 也只落 `ask_every_step`。方向自始至终是"多问一次"，没有一次落向"沉默"。

---

## 5. "模式不改变 OS 能力"（裁定 R21 第 3 条点名要求票 90 明写的那一句）

- **要求**：R21③"**并要求票 90 明写一句"权限模式不改变 OS 能力"**——防的是将来有人把"全自动"实现成"顺手降个权"，那是**假承诺**"。
- **代码侧（本代理量的）**：`grep -rn "CreateProcessWithToken\|ImpersonateLoggedOnUser\|SetTokenInformation\|AdjustTokenPrivileges\|CreateRestrictedToken\|DROP_PRIVILEGES\|SeDebug" internal/risk/mode.go internal/perm/ internal/tools/mode.go internal/config/permmode.go internal/tools/bridge.go` ⇒ **rc=1，零命中**。
  票 90 的三枚提交也**没有一处新增子进程 / 令牌 / ACL 操作**（diffstat 见第 0 节）⇒ **"全自动"确实只是"少问一次"，没被实现成任何形式的降权或放宽 OS 权限**。R21 担心的那个假承诺**没有被写出来**。
- **文档侧（票面）**：`grep -n "OS 能力\|降权\|令牌\|privilege\|AppContainer\|沙箱"` 打在票 90 全文 ⇒ **只有第 25、73-74 行的"现状/建票"叙述**（说 OS 级隔离**没有**），
  **"权限模式不改变 OS 能力"这一句一个字都没写**。
- **时序核对（不许拿"R21 在后头"当挡箭牌）**：R21 落在 `6114e3d`，票面收尾落在 `467040b`，**`git log` 顺序显示 `467040b` 是 `6114e3d` 的孩子** ⇒ 收尾那一刻 R21 已在仓内，**追得上**。
- **判决：本票欠（不是"允许残留 + 转下票"）。** 理由：① R21③ 是**点名给票 90 的判据**，不是给下一张票的；
  ② 它**入库早于本票最后一次票面提交**，没有客观障碍；
  ③ 但它**不影响任何一枚 AC 框的真假**（代码里根本没有可被误读成降权的东西），所以它**不推翻验收结论**，
  只把 `-done` 卡在一行文档上：**挂 `-done` 之前请在票面补一句"权限模式不改变 OS 能力（降权属票 100 RESERVED，R21③）"**。
  ⚠ 顺带**给它记一句好话**：`mode.go:32` 的注释块把 `PLAN.md:3143` 那条红线解释成 "**this layer does not own that clock, it only decides whether a card exists at all**"——
  这句实际上已经**从反面**说了"模式不碰别的机制"，只是没把 R21 要的那句写成正面陈述。

---

## 6. 台账 AC#5：`risk.Gate` 到底有没有真实生产调用者（这一格最像成功、最容易空）

**先数，再跑，最后才对账。**

1. **grep 非测试命中，本代理自己数**：`grep -rn "risk\.Gate(" --include=*.go . \| grep -v "_test.go"` ⇒ **恰好 1 条**：
   `internal/tools/mode.go:98`（`readBlacklist`）。**"零调用点 → 有调用点"这句在调用点这一级为真。**
2. **那道守卫在生产里满不满足**（这是"能不能走到"的关键，不是宣布的）：
   - 守卫：`bridge.go:291` `if len(rawPaths) > 0 && verdict.Level >= risk.L1 { dec.Blacklist = b.readBlacklist(rawPaths) }`
   - `rawPaths = pathArgs(params, entry.Decl.PathParams)`，**生产工具名册真的声明了 `PathParams`**：`fs.go:299,312`、`fs_write.go:711,726,739,753`；
   - `readBlacklist` 只要求 `b.paths != nil`，而 **`git show d5564c2:cmd/wisp/run.go:264` 就是 `Paths: rt.paths`**（不是测试自己塞的）。
3. **再把它跑出来（本代理加的探针 3，探针只在快照、已删）**：用**生产名册 + 生产那一份 Options 形状**
   （`BuiltinFSEntries(FSDeps{Paths:…})` + `Gate` + **故意不给 `Modes`、不给 `Confirmations`**，逐字照 `run.go`），
   对一次真 `fs.read` 打 B 档文件，看只有 `risk.Gate` 本体能打的那行日志有没有出现：
   `PROBE3 err=<nil> isErr=false risk=L2 gatesCalled=true blist=true window=0 approval=1 text=material`
   ⇒ **`risk.Gate` 在生产形状下真的被执行**（`risk: B-list default DENY -> L2` 只在 `internal/risk/blacklist.go:89` 打印），并且**升了 L2、开了一次卡**。
4. **它没有偷偷变成权威**：`Confirmations` 为 nil ⇒ `overrideApplies` 永假 ⇒ B 档**永远"先问"**；
   票面 `unwired.go:81` 的措辞（"Gate 有了真实调用点，缺的是填 `bOverrides` 的确认记录，那归票 21"）**与本代理读数一致**，是**实测形状**，不是文案。

**⇒ AC#5 判 PASS：这句话是量出来的，不是宣布的。**
⚠ **两处诚实记账（都不推翻 AC#5，但必须说出来）**：
- (a) 本代理**没法在本机跑真二进制**去数一次真实命中（`go test ./cmd/wisp/` 加载期 `0xc0000135`，缺 sherpa dll，**票 98 的账**）。
  我拿到的最强证据是"生产名册 + 生产 Options 形状"级别的复现，比票面自己给的（测试自装 Options）强一档，但仍不是生产进程本体。
- (b) **同一份台账里有一句是宣布的**：`internal/config/unwired.go:127` 把 `risk.permission_mode` 记成
  `"consumed: Config.PermissionMode is read per tool call by the bridge (internal/tools)"`——
  **在 `d5564c2` 那一刻生产根本没注入 `Modes`**（`run.go:262-278` 里没有那一行），桥因此**从没读过 Config 里的档位**，
  `permissionMode()` 一直返回 `risk.DefaultMode()`（这条方向是**更严**，所以无安全后果，且已被票 101 立案）。
  ⇒ 但**"防说谎"的那张表自己有一格用将来时态写了现在时**，这属于"改判据让它绿"的**同类倾向**，记给票 101 落地时同批改。

---

## 7. 门禁与本票其余数字

| 门禁 | 命令（都在纯净快照 `/tmp/wisp90acc-90c`） | 实测 |
|---|---|---|
| `gofmt -l` | `gofmt -l internal/risk internal/config internal/perm internal/tools` | **输出为空，rc=0** |
| `gofumpt -l` | `D:/work/base/gopath/bin/gofumpt.exe -l <同四包>` | **输出为空，rc=0** |
| `go vet` | `go vet ./internal/risk/ ./internal/config/ ./internal/perm/ ./internal/tools/` | **rc=0** |
| **D22 门** | `sh scripts/d22scan.sh`（**走脚本，不在仓根 `go run ./tools/d22scan`**，那是独立 module，从仓根跑**根本没跑**） | **rc=0，23.30s**；正向对照先跑 `go test ./...` on `tools/d22scan` ⇒ `ok … 4.468s`（**证明这道门能红，"clean" 不是眼瞎**） |

**`d22scan` 台账逐作用域贴出（本代理原文，非引用票面）：**
```
d22scan: examined 217 production Go files under internal/ and cmd/ of C:/Users/swq/AppData/Local/Temp/wisp90acc-90c
d22scan: scope bans #1-5 internal/      examined 197 production Go files
d22scan: scope bans #1-5 cmd/           examined  20 production Go files
d22scan: scope ban #6 frontend/         examined  37 text files
d22scan: scope ban #7 internal/tools/   examined  17 production Go files
d22scan: scope ban #8 design/           examined  16 text files
d22scan: scope ban #8 frontend/         examined  37 text files
d22scan: scope ban #8 internal/         examined 335 Go files, comments and _test.go included
d22scan: scope ban #8 cmd/              examined  25 Go files, comments and _test.go included
d22scan: clean - no D22 ban violations
```
**对账**：`ban #6 frontend/=37`、`ban #8 frontend/=37` 与票 94 验收时翻上来的那个数（35→37）**一致，没有降**；
`217 生产文件` 与票 90b 自报的 `217` 一致。

**其余被点名的两条仪器账（本代理的态度）**：
- `go test ./cmd/wisp/` ⇒ 本机加载期 `0xc0000135`（缺 sherpa dll，**票 98 的账**）。本票**没跑它**，也**没把它算成票 90 的红**。
- `TestResolvePerCallBudget` ⇒ 第 1 节 9 样本全报（5 PASS / 4 FAIL，数字逐条列在表里）。
  **它是 flake 这一条成立**（票 90 没动 `Resolve` 一个字），但"只在负载下假失败"这句**在本机当前状态下不成立**，
  本代理**没有拿它当借口**掩盖任何红：所有 6 处票 90 的用例红都是**变异当场产生的**，逐条有名字、行号、断言文本。

**AC#1/AC#2/AC#3 的静态与语义补充读数**：
- 没有任何包级可变档位变量：`grep -rn "^var .*Mode" internal/{risk,perm,tools,config}/` 排除测试 ⇒ **0 命中**；
  模式经 `b.permissionMode()`（`tools/mode.go:35`，nil ⇒ `risk.DefaultMode()`，还带 `recover()` 兜底）
  → `mode.Screen(verdict)` → **作为参数**传进 `bridge.go:369 route(ctx, dec, sil)`。
- `ModeAskEveryStep Mode = iota` 是**零值**（`risk/mode.go:49-56`）⇒ 未初始化即最严档，这条是**类型层面**的，改不动。
- 三档确认/审计的语义用例在位并全绿：`TestTicket90ThreeModesAndNoFourth`、
  `TestTicket90OnlyAutoCostsAConfirmAndAllThreeAreAudited`（M4：只有全自动花一次 L2，三档切换全写审计）、
  `TestTicket90NoConfirmChannelMeansNoAuto`（没有确认通道 ⇒ 全自动不可达）、`TestTicket90ConfirmRefusalLeavesModeAlone`、
  `TestTicket90PersistFailureKeepsMemory`（写盘失败内存回滚）。

**L2 超时 300s —— 这个数字本代理逐字量过，没动**：
`internal/config/schema.go:446`（在 `d5564c2` 的快照里）= `ConfirmTimeoutSec int \`toml:"confirm_timeout_sec" default:"300"\``，
`cmd/wisp/run.go:297` 用它构造 `ApprovalTimeout`；票 90 新增的四个文件里 `grep "300"` **只命中两处注释**
（`perm/store.go:86`、`risk/mode.go:32`），**没有任何一处引入第二个超时时钟**。
`perm.Store.Set` 自己**不持钟**（只接调用方 `ctx` + `ConfirmFunc`），所以它**没有能力**把 300s 改长或改短。⇒ **无人在这个数字上 FAIL。**

---

## 8. 四种假绿逐条点名（本票范围，本代理自己的读数）

1. **"没跑"冒充"绿"** —— **未发现**。本代理每一轮 rc 都单独打印（`M*_BUILD_RC` / `M*_RC` / `D22_RC` / `CLEAN3_RC`）；
   `-count=2` 基线做了 `=== RUN` 724 = 2 × 362 不同名的乘法核对；`d22scan.sh` 自带正向对照（先跑 seeded-violation 自扫，`ok tools/d22scan 4.468s`）
   ⇒ "clean" 不是门瞎。**并且这一条正是本票开工时的原罪**：AC#4 那枚勾此前就是"没跑冒充绿"，本代理把六刀重跑，
   **没有一刀是编译失败后还计入读数的**（全部 build rc=0）。
2. **"改了判据让它绿"** —— **未发现票 90 改判据**：`tools/d22scan/**`、`allowlist.txt`、`pathresolver*.go`、`assessor.go`
   在第 0 节的 diffstat 里**零命中**；第 2 节的六刀变异证明**判据咬的是行为不是某一行**（M2 与 M2′ 各红一条不同用例）。
   ⚠ **但抓到一处同类倾向（不在 AC 上，在台账上）**：`internal/config/unwired.go:127` 把 `risk.permission_mode` 记成
   "consumed: … read per tool call by the bridge"，而 `d5564c2` 的生产根本没注入 `Modes` ⇒ **用将来时写了一张防说谎的表**（见第 6 节 (b)）。
   ⚠ 另一处**计数口径**：票面"PASS 222"与本代理"PASS 264/361"不是分歧而是**顶层 vs 含子测试**两种口径（第 1 节已把两种口径的数并列给出），
   下一个人别拿它对不上就判前任造假。
3. **"看图说话"（没测就写）** —— 本代理**刻意没做**：第 3.1 节的 `PLAN.md:1639-1641` 是自己 `awk` 出来的原文；
   第 2A 节"为什么真实链还绿"是**探针打印的 `Kept` 字符串**决定的，不是从"它绿了所以是第二个证人"推出来的；
   第 6 节 `risk.Gate` 有调用点是**跑出来**的（`gatesCalled=true blist=true`）。
   仍未证的一格：**生产二进制里一次真实命中**（本机 `0xc0000135`，票 98 的账）——本代理在结论里按"生产形状"而非"生产进程"表述，没有越级。
4. **"多样本挑运气"** —— **命中并如实全报**：`TestResolvePerCallBudget` 本代理 **9 个样本全列（5 PASS / 4 FAIL，含 0.551/0.917/1.105/1.82s-2.004ms/op 等每一次读数）**，
   没有只报那三次 ok 的；票 94 留下的"静默三次全 ok"在本机当前负载下**复现不出来**，这条**不粉饰**（记在票 86 的账上，不改票 90 的判决）。

---

## 9. 裁决表（与票面六个框 1:1，顺序一致）

| 框 | 票面声称 | 本代理独立读数（指到本文哪一节） | 判决 |
|---|---|---|---|
| **AC#1** 模式作为**显式参数**进链，不许包级全局变量 | `[x]` | `grep "^var .*Mode"`（四包，非测试）**0 命中**；`tools/mode.go:35` nil→`risk.DefaultMode()` 且带 `recover()`；`bridge.go:285-286` 读一次 → `route(ctx, dec, sil)` 以参数传（第 7 节末）；`TestTicket90ModeReachesRoutingAsParameter`、`TestTicket90TwoBridgesDoNotShareOneMode`、`TestTicket90UnwiredModeSourceIsTheStrictestMode` 基线全 PASS | **PASS** |
| **AC#2** 三条红线各一条用例，任何档都必须拒/问；`ban #6` 要有静态门 | `[x]` | 三条红线**六刀变异全红**：不可逆 M1(2 红)/M5、污染 M2/M2′/M3、Deny 不被静默 M4(4 红)/M5；MUT-D 那轮 5 条红名逐字复现（第 2 节）；静态门 `ban #6 frontend/=37` 扫到且 `d22scan` rc=0（第 7 节） | **PASS** |
| **AC#3**（含 (a)/(b)）切全自动一次 L2 强确认 + **审计三档全写**；(a) 手动选过→重启读回那档；(b) 从未选过→读回默认 | `[x]` | `TestTicket90OnlyAutoCostsAConfirmAndAllThreeAreAudited`、`TestTicket90NoConfirmChannelMeansNoAuto`、`TestTicket90ConfirmRefusalLeavesModeAlone` 基线 PASS；(a) 由 `ManualSwitchSurvivesRestart` 盯，**MUT-E 一红**（`after restart mode = ask_every_step, want the manually chosen auto_approve`）证明它真在咬落盘；(b) 由 `UntouchedConfigStartsAtTheDefault`（3 次冷启）盯，坏输入 6 种全部落最严档（第 4 节） | **PASS** |
| **AC#3b** 改模式能跨重启、**会话授权不能跨重启**，两条不许合并 | `[x]`（`d5564c2` 才拆开） | 拆成真拆：**MUT-E 只红 AC#3(a)**、**MUT-F 只红 AC#3b**（`a session grant survived the restart`），互不牵连（第 3.2 节）；三条通道（配置盘 / DB 行 / 内存）**无一能跨重启生效**：`scope="permanent"` 被拒、`blacklist_overrides` 加载即报错、生产 `Confirmations` 为 nil、**生产零授权消费者**（第 3.3 节） | **PASS + 允许残留**（残留＝AC#3b 钉的是"按 session id 查不到"这一形状；补一条"任何消费者拿到的集合必须 session 绑定/未过期"的用例，**那是票 49 的地界**） |
| **AC#4** 双向变异：接模式但去掉红线 ⇒ 红；红线改成模式优先 ⇒ 红；锚点同链 grep 自证 | `[x]`（**此前无一次读数**，由 `agent-ticket90b` 补 5 轮） | 本代理**自己重做 6 刀**（多加一刀 M2′ 镜像刀）：每刀 **同链 grep 证落地 + `go build ./...` rc=0 + 红名 + 断言原文首句 + 还原 `diff -q`/整树 `diff -rq` 一致**；它写的 2/1/3/4/5 红**名单、行号、文本全部对上**（第 2 节表格）。它那句"真实链仍绿是因为 `case R4:` 还在"**被探针证实**（`Kept` 字符串在 M2 下变成 `case R4:` 那句、M2′ 下变回 flag 那句、M3 下真实链当场红） | **PASS**（补作业的数字**经得起重跑**；归因**已升级为实测**） |
| **AC#5** 台账：`risk.Gate` 从"生产零调用点"变成有真实调用者，或明写仍归票 21 | `[x]` | 非测试命中**恰好 1 处**（`internal/tools/mode.go:98`），守卫三条件在生产里**全部满足**（`run.go:264 Paths:` + `fs*.go` 的 `PathParams`），并用**生产名册 + 生产 Options 形状**真跑出来：`gatesCalled=true blist=true approval=1`；`Confirmations=nil` ⇒ B 档永远先问，未变成权威（第 6 节） | **PASS**（残留＝**没能在生产进程本体里数一次**，本机 dll 挡着，票 98 的账；同表另有一格 `permission_mode` 的 `consumed:` 用了将来时，见第 6 节 (b)） |

**票外两格（不在六框里，但被裁定点名，故单列）**：
- **R21③"明写一句权限模式不改变 OS 能力"** ⇒ 代码零降权形状（实测），**但票面一个字都没写**，且 R21 入库早于票面收尾 ⇒ **本票欠（一行文档，非阻断性）**（第 5 节）。
- **300s 这个数字** ⇒ `default:"300"` 在位，票 90 新代码没引入第二个时钟，**没人动它** ⇒ **无 FAIL**（第 7 节末）。

---

## 10. 三句必答

**① 票 90 能否挂 `-done`？**
**能挂，但先补一行再挂。** 六框**没有一格是 FAIL**，每格都有本代理自己重跑的读数；AC#4 那枚历史上"涂了勾没跑"的框
**在本代理独立重做的 6 刀里全部复现**，并且那条归因从"叙述"被本代理升级成了"实测"。
挂 `-done` 前请做**一行**收尾：把 R21③ 要求的那句"**权限模式不改变 OS 能力**"写进票面（第 5 节判"本票欠"）。
两条允许残留（AC#3b 的授权形状判据 → 票 49；`permission_mode` 台账将来时 → 票 101）**都有立案去处，不欠在本票**。

**② "三档模式今天已经能用"这句话能对 owner 说到什么程度？**
**只能说半句，且必须带一个"但是"。** 可以说的是：**"三档的语义、持久化、切档确认与审计，今天在代码里是真的，并且被真链测住了"**——
档位写在 `config.toml` 的 `risk.permission_mode`，改过就跨重启（MUT-E 证明这条在咬），没改过就是最严档，坏输入一律落最严或响亮报错，
红线三刀拆不穿。**必须带的那个"但是"**：**今天没有任何生产进程会读这个档位**——`cmd/wisp/run.go` 在 `d5564c2` 没注入 `Modes`，
`internal/perm` 生产零 importer，所以真跑起来的 Wisp 的行为是**永远按最严档"每步都问"**（fail-closed 兜住，不会变松）。
⇒ 对 owner 的准确说法是：**"开关造好了、装错了位置还没接上；接上之前你感觉到的行为=默认最严档，安全但也不能切。"**
**票 101（`ea7913f` 已立案）判：不挡 90 结案** —— 90 的交付面是"三档语义/持久化/确认/审计/红线/Gate 台账"，
六框没有一框写"装配完成"；票面第 189-193 行已把这条**自己算成交接而不是成果**（"连那一行都没落"），
把装配缺口另立一票是**正确的切分**，不是漏。⚠ 但**M3 是 owner 亲手推翻我的那一条**，
所以"持久化今天对用户不可感知"这件事**必须在给 owner 的那句话里说出来**，不能只留在票面。

**③ 有没有任何一条路径能让"会话授权"跨重启？**
**没有，三条通道本代理逐个量过（第 3.3 节）。** 配置盘上 `blacklist_overrides` 写了就**加载报错**（"a file on disk cannot have clicked anything"）；
数据库上 `InsertGrant` 只认 `scope="session"` 且 `expires_at` 必填，取用只按 session id（MUT-F 那把刀证明改掉这个 WHERE 就会红），
死会话那行只是**审计留存**、票 90 的用例还**断言它必须留着**；内存侧 `Store` 只缓存切档记录、`Bridge.Confirmations` 生产是 nil。
更硬的一条：**今天全仓没有任何生产授权消费者**（`InsertGrant/ListGrants` 非测试命中 = 0），没有出口可被误用。
⇒ **安全底线未破；AC#3b 的形状判据留一条"补用例"给票 49，不记在票 90 头上。**
