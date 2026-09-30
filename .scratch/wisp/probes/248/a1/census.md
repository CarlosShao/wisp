# 票 248 · 只读普查件 `248-a1` — 「设置」那条路的三问 + C17 文档漂移一问

> 腿：`248-a1`（只读普查） · 票面：`.scratch/wisp/issues/248-...md` AC#0 · 归口：面板
> 起手锚点：`2b4e2602`（分支 dev）· 起手钟：`2026-09-30 23:41:16 +08`（`date` stdout）
> 写点只有本文件一处；⛔ 本腿未改任何产码／测试／配置／`go.mod`／`go.sum`／票面。
> ⛔ 本件不含任何凭据值，只写字段名与变量名。
> ⚠ 本件里所有对 `cmd/wisp` 的引用一律是〔起手读数，可能被票 33-r1 改动〕。

状态：**骨架已立，各节填写中。**

## ① 要新增哪几枚方法名：参数形状 / 错误形状 / 生效级别 / 重启档 / posture 敏感字段

填写中。

## ② 写路径落在谁手里：唯一写者 / 依赖边代价 / 两处推广具名

填写中。

## ③ 读侧怎么回答「配了没给／配好了」而不泄露值 + 快照那一维

填写中。

## ④ C17 那张方法表在文档里有几处拷贝、漂移会出现在哪几行、有没有仪器看得见

填写中。

## ⑤ 我跑了哪些尺、每条真实读数

> 纪律：本腿⛔ 未跑 `go build`／`go vet`／`go test`／任何 `./...`；只跑过 `grep`／`find`／`ls`／`sed -n`／
> `git status`／`git log`／指名单包的 `go list -deps`。所有当下状态断言与取数同发，带 `date` 钟点。
> 共 **34 把尺**（R01-R34，含 R34 那把我跑空了的尺）。

| 号 | 尺（命令或读法） | 真实读数 | 钟点 |
|---|---|---|---|
| R01 | 编排者原尺：`grep -rn "panel\.[a-z]*\.[a-z]*" --include=*.go internal/panel/ \| grep -v _test \| grep -o … \| sort -u` | **恰四枚**：`panel.attachment.add`／`panel.message.send`／`panel.mode.request`／`panel.workspace.request` | 23:42:32 |
| R02 | 同一把尺扩到 `internal/ cmd/ tools/`（非测试） | 同样**四枚**，一枚不多 | 23:42:32 |
| R03 | 扩到测试全量计数（`uniq -c`） | 去重 **23 枚**；最高 `panel.approval.request` 17 次、`panel.review.allow` 11、`panel.mode.set` 7 | 23:42:32 |
| R04 | `grep -rniE "panel\.(config\|secret\|credential\|settings)[a-z._]*" --include=*.go .` | **零命中** ⇒ 票面现量 #1 成立 | 23:42:32 |
| R05 | 每枚字面量的 `file:line` 归属 | 生产码只有 `internal/panel/bridge.go:42-45`（声明）、`knownComposerMethod`（派发）、`composer_handlers.go:85,109`、`cmd/wisp/panel_inbound.go:196`；其余 19 枚**全在 `*_test.go`**，其中绝大多数在冻结件 `l2_grant_boundary_test.go`（植入负控） | 23:42:32 |
| R06 | `sed -n '95,135p'`／`'490,560p'` `internal/panel/composer_test.go` | `panel.mode.set` 是**被禁形状**（`composerModeWriteRe`），只出现在植入件里；`TestTheRendererHoldsExactlyOneDoorToTheHost` 要求「渲染侧每一枚 `panel.*` 字面量都是 Go 侧答的那枚」 | 23:42:47 |
| R07 | `grep -rn "宁缺毋造" --include=*.go --include=*.md .` | 16 处，**全在 `.scratch/**` 的票与派单、`docs/evidence/**`**；Go 侧成文的同规矩在 `internal/panel/pump.go:25-46`（"NO NEW KEYS"）与 `composer.go:44-56` | 23:42:55 |
| R08 | `grep -n "json:\"" internal/panel/composer.go` | `Snapshot` 结构体 **6 枚带 tag 字段**：`pending`/`results`/`composer`/`generatedAt` 恒发 + `instructions,omitempty`(:74) + `tasks,omitempty`(:91)；`ComposerState`（`:235-256`）**9 枚键**：mode/workspace/attachments/acceptedAttachmentMimes/maxAttachmentBytes/attachmentError/git/currentModel/modelKnown | 23:43:08 |
| R09 | `sed -n '340,430p' cmd/wisp/run.go` | 单写者那句在 **`run.go:381-390`**；`secret.NewStore` 唯一装配点在 `run.go:371` | 23:43:08 |
| R10 | `grep -nE "^func (m \*Manager)…"` `internal/config/manager.go` | Manager 导出面只有读与重载：`NewManager`/`Config`/`Resolved`/`CheckAndReload`（+`plan*` 私有），**没有任何通用 Set** | 23:43:15 |
| R11 | `grep -nE "^func (s \*Store)…"` 四枚文件 + `sed -n '18,30p' internal/secret/dpapi_windows.go` | `Store.Store/Resolve/Exists/Delete/Dir/Blobs`；`protect` 在 `:18-28`，`CryptProtectData` 那行＝**`:24`**（票面引用成立） | 23:43:15 |
| R12 | 写面检索：`grep -rn "^func " writeguard.go permmode.go` + `grep -rn "SetPermissionMode" \| grep -v _test` | 受保护的写原语只有 `(*Manager).mergeWrite(ownedKey, setOn)`（`writeguard.go:111`）；其上只有**两枚导出写者**：`SetPermissionMode`（`permmode.go:70`）、`AddAllowedDir`/`SetAllowedDirs`（`allowdirs.go:62,109`）。生产调用者：`internal/perm/store.go:208` | 23:43:21 |
| R13 | `go list -deps ./internal/panel`（指名单包） | wisp-internal 起手名册 **8 枚**：`frontend`(embed)/`observe`/`panel`/`projctx`/`risk`/`secret`/`streamkey`/`winsec` | 23:43:36 |
| R14 | `grep -rn "CarlosShao/wisp/internal/" internal/panel/*.go \| grep -v _test` | 直连边只有三枚：`risk`(5 文件)、`projctx`(2)、`streamkey`(1)；**panel 不直连 `config`，也不直连 `secret`** | 23:43:46 |
| R15 | 谁把 `internal/secret` 拉进 panel 的闭包 | `internal/observe -> internal/secret`（redact 那侧），是**传递**边不是面板的直连边 | 23:43:52 |
| R16 | 提议接法的依赖增量：`comm -13` 两份 `-deps` 读数（wisp 口径 + 全量口径） | `panel -> config` 多出 **恰好 1 枚**（两种口径同值：只有 `internal/config` 本身；第三方零增量） | 23:44:00 / 23:47:34 |
| R17 | `grep -rn "C17" docs/PLAN.md docs/specs/` | **13 行 / 5 个文件**：`PLAN.md:1367,1747,1808,2164,2434,2968`；`SPEC-01:59`；`SPEC-08:3,156,235`；`SPEC-10:3,18`；`SPEC-12:48` | 23:43:36 |
| R18 | `sed -n '150,190p' docs/specs/SPEC-08-ui-ball-panel.md`（§5.2 那张表） | 文档表列 **20 余名**（`panel.resync`/`tasks.list`/`history.query`/`approval.decide`/**`config.get` / `config.set`**/`grants.*`/`privacy.*`/`cost.summary`/`models.list`/`diagnostics.export` + 五枚事件推送），标题自署**【SPEC 提案，S5 定稿走契约批准】** | 23:45:06 |
| R19 | `sed -n '2960,2985p'`／`'2430,2440p'` `docs/PLAN.md` | PLAN 侧对 C17 只提**四项要求**（白名单/capability+是否二次授权/correlationId 背压/`panel.resync` 无状态），**全篇未枚一枚方法名** | 23:45:06 |
| R20 | `sed -n '20,95p' docs/specs/SPEC-03-config-secrets-envs.md` | §3 表 `:24-43`（逐段生效级别）；§4.2 `:73-79`（hot/reload/restart/🔒 四支）；§4.3 `:81-85`：**"GUI 改配置 = 写回 TOML → 同一检测路径生效（单一真相源始终是文件）"** | 23:46:10 |
| R21 | `grep -n "permission_mode" SPEC-03` | **0 命中** ⇒ 文档 §3 的 `[risk]` 行（`:34`）里没有 `permission_mode`，代码 `schema.go:470` 有 ⇒ 文档↔代码**已在漂** | 23:46:24 |
| R22 | `sed -n '243,300p' internal/config/manager.go` | 分档实现：`planLocked` 四枚锁定段 → L2；`planApp`＝restart（theme 除外）；`planVoice`＝reload；其余 12 段整体 hot（`llm` 在 `:281`） | 23:46:33 |
| R23 | 冻结件 `internal/panel/l2_grant_boundary_test.go` 五处读段（头段 1-60、1240-1290、1360-1395、1408-1510、2021-2160、2290-2345） | 四条 facet + 植入件族；`git log` 该文件工作树**干净**、2380 行（23:48:08 同发） | 23:44:11-23:49:36 |
| R24 | `grep -n "grantRouteWords = " -A 20` + `carriesGrantWord` 实现 | 裁决词表 11 枚：`approve/approval/grant/allow/permit/ratify/authorize/authorised/decide/decision/verdict`；归一化**抹掉 `._-` 与空格再子串匹配**（`:242-250`） | 23:44:23 |
| R25 | 尺寸／拼写硬钉位置检索（`wantAnswered`、`!= 4`、`guardAnchor`、`fieldAnchor`） | `:1469 wantAnsweredInGrid=4`、`:1496 wantAnswered{四枚常量}`、`:1635`/`:1983` 逐枚循环、`:2043 fieldAnchor="\tText string \`json:\"text,omitempty\"\`"`、`:2051 guardAnchor="case MethodModeRequest, MethodWorkspaceRequest, MethodAttachmentAdd, MethodMessageSend:"`、`:2139 got != 4` | 23:44:30 / 23:48:56 / 23:49:15 |
| R26 | `sed -n '350,400p'`／`'505,525p'` `internal/panel/git_test.go` | 票 181 的反漂移钉 AC#3 门 3：`wantMethods` = 四枚常量（`:381-390`，排序全等比较）；同一函数里还有正控（植入第五枚 `panel.worktree.switch` 必数到 5，`:514`）与实尺（`len(real) != 4`，`:517`） | 23:45:19 / 23:48:35 |
| R27 | 「有没有仪器把代码方法名核文档那张表」：`grep -rln "SPEC-08\|docs/PLAN.md" --include=*.go` + `grep -rn "docs/specs\|ReadFile" ` + `grep -rniE "config\.get\|config\.set\|models\.list" --include=*.go` | **没有。** 全仓只有一枚仪器把文档当预言机：`git_test.go:365 os.ReadFile(docs/PLAN.md)`，核的是 **D34 工具表**里不许有 `git.*` 行；Go 侧提到文档方法名的只有 `internal/session/grants.go:40` 一句注释 | 23:45:19 / 23:45:35 |
| R28 | `grep -nE "1 bare-goroutine…8 emoji" tools/d22scan/main.go` | 八条禁令射程内**没有任何一条**看 C17 方法名册；#3 plaintext-key 只看「key/token/secret 命名的标识符被赋字符串字面量」 | 23:45:35 |
| R29 | `grep -n "confirmModeSwitch\|rt.pump = " cmd/wisp/run.go` + 函数体 | `confirmModeSwitch` 在 `run.go:805-827`（一张 C18 卡、`host:mode-switch`、非 Allow 即错）；只有 `wisp run` 装它（`:626`）；泵装配在 `run.go:673` | 23:46:06 |
| R30 | 凭据读侧原语语义：`sed` `store.go:120-175`／`refs.go:1-80`／`configrefs.go`／`redact.go` | `Exists` **只答 `dpapi:`**，`env:` 直接错（`store.go:129-145`）；`ConfigRefs(configPath) map[字段]ref`（`configrefs.go:23`）；`RefFieldNames(configPath, ref)`（`:62`）；`RedactSecret`＝除末 4 字符全 `*`（`redact.go:6-14`）；`NewRef()` 随机 hex（`refs.go:72`） | 23:46:47 |
| R31 | `grep -nE "func \|Err \|api_key" internal/llm/resolver.go` | 解不开 key 的两个响亮出口：`:139`（有 ref 无 resolver，Unconfigured）、`:141-144`（`Keys.Resolve` 失败 → Unconfigured）；错误文案带 ref **形状**不带值（`secretRefShape`，`:154`） | 23:46:54 |
| R32 | `grep -n "func routeLabelsInFunc" -A 40` 冻结件 | 枚举行走的是**函数内每个 `ast.SwitchStmt` 的每条 `CaseClause`**（`:887-913`）⇒ 新增独立 case 行也进名册，`:2139` 那枚「4」照样变 5 | 23:49:51 |
| R33 | `git status --short` + `git rev-parse --abbrev-ref HEAD` + `git log --oneline -3` | 起手锚点 `2b4e2602`（dev），工作树有大量他人未提交改动（`design/**` 删除等）⇒ 我全程只用显式 pathspec | 23:41:16 |
| R34 | ⚠ 一把跑空的尺：`grep docs/specs/SPEC-03-config-secrets.md` | `No such file or directory`（实名 `SPEC-03-config-secrets-envs.md`）⇒ 23:46:06 那次读空，23:46:10 改用通配重跑才有 R20 的读数 | 23:46:06 |

## ⑥ 我可能写错的条目（对抗我自己）

1. **我最先要推翻我自己的一条**：写这节前我曾把「白名单扩一枚 ⇒ 冻结件红」的证据记成 `l2_grant_boundary_test.go:2144-2153` 一段写着 `want 4: the whitelist changed size` 的文本。
   23:48:28 与 23:48:35 两次复查：`grep -rn "changed size" --include=*.go .` 与 `grep -rn "did not move with it" .` **全仓零命中**。那段文字不属于本仓任何文件。
   ⇒ 本件只承认查得的三枚真钉：`:2051 guardAnchor`（精确拼写）、`:2139 got != 4`（植入件的实尺快照）、`git_test.go:385/:514/:517`。任何后来人按我原来那句去找 `:2144` 会找不到。
2. **我一个测试都没跑**（票面禁令：246-r2 正在跑 `cmd/wisp`）。所以「加第五枚会让 `:2139` 红」「facet 2 不拦 config 族键名」这类结论**都是从码面推出来的结构判断，不是读数**。裁决者应把它们当预演，不该当已验。
3. R08 数的是**结构体 tag**，不是**发射态**。票面现量 #5 说「快照今天只有四键」，`pump.go:25-46`／`composer.go:44-56` 也说四键，但 `Snapshot` 有 6 枚带 tag 字段，其中两枚带 `omitempty`。⇒ 真实一发到底送 4 枚还是 6 枚，只有跑泵才知道；本件只报「4 恒发 + 2 视装配」，**没跑过就不写死**。
4. owner 给的 `run.go:356` 与我读到的 `:381-390` 不是同一处 ⇒ 行号已被 246-r2 落地推移。我引的是〔起手读数，可能被 33-r1 改动〕，**别按 356 去对质**。
5. 我把「SPEC-08 §5.2 有 `config.get`/`config.set` 而代码没有」写成「文档↔代码漂移」。可能判错：那张表标题自署**【SPEC 提案，S5 定稿走契约批准】**（`SPEC-08:156`），它可能是一张**尚未定稿的提案**，不是「批了没实现」。两种读法对 248-r1 的命名自由度影响不同 ⇒ 交裁（见⑦ J8）。
6. facet (2)（信封不得绑定裁决形键名）我只读了头段、`carriesGrantWord`、`grantRouteWords` 与调用点，**没逐行读 `scanGrantBoundary` 的键名收集实现**。如果它对「值／secret 形键名」还有额外约束，我①里「config/secret 族名安全」的结论就只覆盖了裁决词表那一半。
7. 冻结件行号可能在我只读期间被别的腿改动。同发核过的那一发是 23:48:08：`git status --short` 对 `internal/panel/l2_grant_boundary_test.go` 与 `bridge.go` **干净**、`wc -l`=2380。此后引用不带同发读数的一律按「某时刻读数」读。
8. R16「panel→config 只多 1 枚」是**新增直连边**的代价；如果实现腿改成从装配根（`cmd/wisp`）把 config 读写腿以接口注进 panel（像 `composer_handlers.go:73-78` 的 `ModeWriter` 那样**本地声明接口、不 import**），代价＝**0 枚新依赖边**。我把两种接法都列进②，没替实现腿选。
9. ③里「消息文本进上下文预算与日志」我只给了**判据形状**，没实测 logsink 的落盘目录与轮转后可 grep 面（`installLogSink` 我在 `secret.go:226` 读到调用，没读到实现细节）。哨兵正控真正要做到的覆盖面，得由跑的人定。
10. 我把 `panel.mode.set` 读成「被禁的渲染侧拼写」（R06）。它是**植入件**里的字符串，不是名册成员；如果实现腿误以为它是既有能力去「复用」，会把 `composer_test.go:111-120` 那枚正控钉踩响。本件按 R05 的归属读数写它＝只存在于测试。

## ⑦ 判不动的地方（逐条甲／乙／不做 + 现量）

> 以下八条**本腿一条都没动**，全部交编排者裁。J1 是这票真正的闸门。

**J1 冻结件与 AC#1 正面相撞：白名单每扩一枚，`l2_grant_boundary_test.go` 必红。**
现量：①`:2051` 要求 pristine `bridge.go` 里**逐字**含 `case MethodModeRequest, MethodWorkspaceRequest, MethodAttachmentAdd, MethodMessageSend:`，不含即 `t.Fatalf("plant B has no anchor")`；②`:2139` 在「复制真实包 + 植入」的快照上断 `len(pkg.answered) != 4` 即红，而 R32 证明 `routeLabelsInFunc` 走**全部** case 子句；③`:2043` 还逐字钉着信封里那行 `Text string`。该文件是 AGENTS.md §1.5／票 248 AC#5 的**三枚冻结件之一，实现腿禁改**。
- 甲：把「冻结件里那两枚**尺寸/拼写**钉」按人工批准改成**从 guard 派生**（`len(answered)` 与常量名册一致，不再写死 4），与 `A484` 同批登记；`AGENTS.md §1.2` 的「不许为变绿放宽断言」在这条形不成文——它是尺寸快照不是安全断言，但这话不归我说。
- 乙：白名单**不扩进 `knownComposerMethod`**，新能力另开一条管道——但 AC#1 明确要求「进 `bridge.go` 白名单与 `ParseComposerRequest`」，且另开门正是该冻结件 facet 1 的「second switch/if chain」要判红的形。
- 不做：本票 248 的 AC#1 在冻结件重排之前不可落地；先做②③（读侧那一维）也不通，因为没有入口。
⇒ **要我拍**：甲/乙/不做，选一个。

**J2 票 181 的反漂移钉同样写死四枚（非冻结，但那是别人 AC#3 的判据）。**
现量：`internal/panel/git_test.go:381-390` 排序全等比较 + `:517` `len(real) != 4` + `:514` 植入第五枚必数到 5 的正控。
- 甲：248-r1 在同一枚 commit 里把 `wantMethods` 与那句「this ticket may not extend」的报错文案一并改写，并在证据里写明是票 248 动的。
- 乙：把票 181 那条判据改成「读侧展示不新增方法」的性质判断（能力形），不再数枚数——AC#1 的判据形状本来就写着「判据一律问能力，不许做成扫注释词面」。
- 不做：让 248-r1 去撞一枚写着「本票不许扩」的尺。
⇒ 要个口径；我倾向**乙**，但那是票 181 的判据，不是我的。

**J3 快照那一维落在哪：AC#2 要第五枚键，`pump.go:25-46` 写「NO NEW KEYS」，票 145 AC#3 已裁「不在这片切片」，Q-51（谁准动 `frontend/src/lib/panel.ts`）未定案。**
现量：R08（4 恒发 + 2 可选；`ComposerState` 另有 9 枚键）；四枚钉住键集的件列在 `composer.go:44-56`（`composer_test.go:48`、`approval_test.go:105`、`pump_test.go:111-124`、`:270-276`）。
- 甲：这一维**加进 `composer` 段**（与 `currentModel`/`modelKnown` 同形，`:255-256` 已是「一枚事实 + 一枚 known 位」的先例），不动 `Snapshot` 顶层。
- 乙：加顶层第五键，按 145 的路子在**同一 commit** 里动那四枚钉（其中两枚是对 `frontend/src/lib/panel.ts` 的双向核对 ⇒ 触发 Q-51）。
- 不做：这一维先不进快照，只走**读方法**的回执（面板点开设置时才问），AC#2 那格因此勾不上。
⇒ 三条都改得到东西，我没资格选。

**J4 凭据值过境的信封形状：AC#0 只禁了「塞进 `panel.message.send`」，但没有仪器拦得住。**
现量：所有已答方法解进**同一枚** `ComposerRequest`（`bridge.go:69-79`），facet 2 的裁决词表（R24）**不含** `value`/`secret`/`config` 等名 ⇒ 一枚带值键一旦进信封，`panel.message.send` 那条路**同样能带值**且解析通过；今天没有任何尺按路由分键。
- 甲：值**不进信封**——写方法只带 `api_key_ref`（`dpapi:<id>`），值另走一次「只写不回显」的一击（形状待定，可能必须新增一枚 `panel.secret.set`，其处理器**立即**交 `Store.Store` 且不留字段）。
- 乙：进共享信封但**按路由拒键**：新增一枚机读判据＝种「`panel.message.send` 带 value」必被拒，且给这判据配正控。
- 不做：靠票 248 §5 那句「界面侧不要把密钥塞进聊天消息」当约束（那是文案，不是闸门）。
⇒ 要个形状；这一条还直接决定 AC#3「不许新增任何把值回显给页面的方法」怎么被核。

**J5 `permission_mode` 的写：面板侧改档今天只有 `panel.mode.request`，而 `wisp panel-inbound` 的 `Confirm` 是**刻意 nil**。**
现量：`panel_inbound.go:213-231`（`Confirm: nil` + 头段 `:32-39` 的理由）、`composer_handlers.go:133-138`（变宽即拒 `ErrNoL2Confirm`）、`run.go:626,805-827`（只有 `wisp run` 装了那张卡）。
- 甲：配置写方法**显式拒写 `risk.*` 全族**（连收紧也拒），档位一律走既有 `panel.mode.request`，一行代码都不给「从设置页放宽」留门。
- 乙：允许写**收紧向**、放宽向交给 `Confirm` 腿；宿主没腿即 `ErrNoL2Confirm`（把 `ModeWriteHandler` 的方向不对称推广到配置段——⚠ 这是我的推广，见②）。
- 不做：让设置页直接写 `risk.permission_mode`（那正面撞 `PLAN.md` D33/D36 的「放宽不得静默生效」与 AGENTS §1.2「面板侧来源的 L2 允许」）。
- 现量补充：票面 §5 第 1 条把「权限档位」列进了那页表单，**这张表长什么样不归本编队**（AC#5），但「档位能不能从这页写」归本票裁。

**J6 一枚通用 `config.set(section, key, value)` 还是按字段数枚具名方法。**
现量：`A484` 的边界＝「只新增配置读/写那一族方法名」；写面既有形状是**按键**的 `mergeWrite(ownedKey, …)`（`writeguard.go:111`，`ownedKey` 是逐枚键路径，`llm.providers.<name>` 在 `:199` 有具名形）；生效级别在代码里是**按段**判的（R22）。
- 甲：读一枚 + 写一枚（键路径 + 值），级别与方向由 config 层判（复用 `plan*`/`riskDirection` 那套），白名单只多两枚名。
- 乙：按字段族数枚（provider 基址/模型名/上下文窗口/价格/`api_key_ref` 各一枚），名册大，但每枚参数字面窄、错误形状具体。
- 不做：本腿继续猜。
⇒ 这决定 `A484` 那张批准记录的射程要不要重写（它写的是「那一族」，不是枚数）。

**J7 重启档与「写完没生效」的回执形状。**
现量：restart 段今天只有 `[app]`（theme 除外）＝`SPEC-03:78` + `manager.go:338 planApp`；那句「本宿主没有接 config.toml 轮询…要生效请重启进程并用 `wisp run`」是 `config_reload.go:97-99` 的 `hotReloadDisabledPanelInbound`，在 `panel_inbound.go:212` 每启动念一次；`Report` 有 `Restart []string`（`manager.go:92-94`）但**没有任何出口把它交给面板**（H10 未建）。
- 甲：写方法的**返回值**里就带 tier 判定（`applied`/`needs-restart`/`denied-loosening`），不新增回显通道。
- 乙：只回执「已写入」，把「要不要重启」全交给下一次 `panel.resync`（依赖票 33/35）。
- 不做：让页面自己按字段名猜哪枚要重启（那是第二套真相）。
⇒ 顺带一问：设置页写完之后**谁来念那句 hotReload 的话**——`panel-inbound` 宿主现在只在 stderr 念，页面上什么都没有（`panel_inbound.go:44-46` 自陈 H10 开着）。

**J8 文档那三张表谁权威、代码扩名后要不要跟着动（本腿一字未改）。**
现量：`docs/specs/README.md` 的优先级（PLAN 定案最高）＋ R17 的 13 处 C17 ＋ R18 那张自署「提案」的 §5.2 表 ＋ R21 的 SPEC-03 缺 `permission_mode`。
- 甲：以 §5.2 那张表为**待批提案**，248 落地时**只**在该表加两行并标批准号（属 `docs/specs/**`，实现腿禁写 ⇒ 谁写？）。
- 乙：承认文档那张表已过期，另开一票做文档同步，248 只动代码。
- 不做：什么都不动 ⇒ ④ 的漂移面再加一层（代码 6 枚 / 文档 20+ 枚 / 提案未定稿三种名册并存）。
- ⛔ 我不改 `PLAN.md`/`docs/specs/**` 一字；这条只要一个裁。
