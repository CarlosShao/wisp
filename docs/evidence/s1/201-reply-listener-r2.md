# 票 201 · 续腿 `201-r2` 证据（答复入口的第二枚 / 等待态 / 长期那一支）

- 起手锚点：`93ea9f6b`（编排者代落的 201-r1 半成品）。本腿 commits：`58bf0158`（改编号）→ `49eb440b`（接缝＋写面＋always 支）→ `5a251ece`（判据＋AC#6 调用位置更正）。
- 本文件只写**这一腿新增的事实与判据**；201-r1 的读数在 `201-reply-listener-r1.md`，不重抄。
- 派单原文＝`.scratch/wisp/issues/201-…md` 顶部 09-29 12:3x 那一段（编排者已就地改写 AC#2/AC#5，本腿按改写后的口径交）。

---

## ① 起手现量：这枚二进制里到底有哪些"人能按的面"

派单要我"先花 ≤5 次调用把宿主面数清楚"。尺与读数（都在本腿自己跑，不引 201-r1）：

| 问的 | 读数 | 尺 |
|---|---|---|
| 球／托盘编进生产二进制了吗 | **没有**：`internal/ball` 的非测试引用者只有 `cmd/balldebug/main.go` 一枚 | `grep -rln "wisp/internal/ball" --include=*.go .` → `./cmd/balldebug/main.go` |
| `cmd/wisp` 里有没有 `ball.` 的调用 | **零** | `grep -rn "ball\." --include=*.go cmd/wisp/ \| grep -v _test.go` → 空 |
| 面板 WebView 的入向在不在 | **不在**：本树无 postMessage 接收器、无 router、无宿主 | `internal/panel/pump.go:16`「no WebView2 host (ticket 33 is unclaimed)」、`internal/panel/git.go:82` 同句、`go.mod` 无任何 webview/web 依赖 |
| 状态机在生产路径有没有生产者 | **没有**：`statemachine` 的引用者是 `cmd/balldebug` 与 `internal/tools`/`internal/panel` 的读侧 | `grep -rn "statemachine" --include=*.go cmd/ internal/ \| grep -v _test \| grep -v "^internal/statemachine"` |
| 面板快照能不能读"在等人" | **能**：`Snapshot.Pending []ApprovalCardView` 每轮真造、卡片事件后重发 | `internal/panel/composer.go:58`、`cmd/wisp/panel_pump.go:62`（`Queue().LiveApprovals()`）、`cmd/wisp/run.go` 的 `u.publish()` |

⇒ 按派单的第二条支：**面存在但这个二进制够不着**。所以本腿交的是 Go 侧接缝＋判据，"最后一寸"白纸黑字写在 §⑤，不替它打勾。

⚠ 一条派单没说但影响判断的事实：**接缝本身过去也不是可交给宿主的**。`DecideFromNative/DecideFromPanel` 是导出的，但"卡片↔一次性令牌"的账本当时住在 `package main`（`nativeCards`），别的包 import 不到 ⇒ 任何非终端的宿主面要么自己再抄一份账本（＝令牌管理分叉），要么接不上。本腿把账本与路由选择搬进 `internal/agent/approval/replies.go`，cmd/wisp 的 console 改为**委托**它。

---

## ② 逐格判定（met / met-with-seam / blocked）

| 格 | 判定 | 生产执行点（不是"测试能碰到"，是二进制里真跑的那行） |
|---|---|---|
| **AC#1 有人能答** | **met-with-seam** | 可交给宿主的导出入口：`internal/agent/approval/replies.go:313 Allow` / `:336 Reject` / `:419 Veto`（原生路由）与 `:342 PanelReject` / `:395 PanelAllow`（面板路由）；账本入账 `:167 Record`。生产构造点：`cmd/wisp/run.go:458 rt.liveCards.bind(rt.gate, …)`（每次 run 都绑，不依赖监听器）；console 作为第一个调用者：`cmd/wisp/approval_reply.go:210 / :278 / :310`。**入口数仍是 1 枚（console）**，球/托盘/WebView 零 —— 见 §⑤ |
| **AC#2 L1 能被否决＋等待态可见** | **met（按编排者改写后的口径，极性未翻）** | 否决走 `replies.go:419 Veto`→`Gate.Veto`；到点无否决仍执行（`internal/tools/bridge.go:384-386` 的冻结语义，本腿一字未动，`gate.go` 的 `ANSWER-EXPIRED` 行是 201-r1 已落的有声记录） |
| **AC#3 L2 批得动 / 面板来源的允许被判红** | **met** | 批得动：`approval_seam_201_test.go` 第一枚——整场 `h.reply=nil`（监听器根本没起），宿主调 `Replies.Allow` 后 `fs.write` 真落盘。判红：同文件第二枚——`Replies.PanelAllow` 返回 `ErrPanelAllow`、被烧掉的令牌再拿去原生 allow 得 `ErrBadGrant`、`PanelReject` 仍能拒且目标文件不存在。拒绝豁免的仍然是**路由**（`gate.DecideFromPanel`），不是便利 |
| **AC#4 "一直"印出规则＋二次确认** | **met** | `cmd/wisp/approval_always.go:70 always` → 先印"要存的规则：…`allowed_dirs += \"…\"`"（文本由 `replies.go` 的 `WideningRule` 从卡片自己的 Paths 推，不唯一/裸盘根一律拒），再 `:99` 起监听器的续体 → `:112 runWidening` 用 `gate.PendingApproval`（:115）真起第二张 L2 卡；批准才 `:134 rt.mgr.AddAllowedDir(dir)`，落 `:140 WIDEN-APPLIED`，否则 `WIDEN-REFUSED` |
| **AC#5 长期那一支真落库** | **met（落的是 config.toml 的 [fs] allowed_dirs；session 那一支照旧不活过重启）** | `internal/config/allowdirs.go:56 AddAllowedDir` / `:87 SetAllowedDirs`：应用→`SaveFile` 原子写→失败回滚→认领自家 mtime（沿用 `permmode.go:64 SetPermissionMode` 的先例）。`approval_grant` 表本腿**没有**加生产写手，`internal/perm/ticket90_persist_test.go` 三枚钉子未碰 |
| **AC#6 等待态可见** | **met-with-seam** | 读侧：`replies.go:249 AwaitingHuman` / `:283 WaitingState`（只回 D43 冻结名 `Confirming`/`AwaitingApproval`，不新造名）。生产调用者：`cmd/wisp/run.go:1103` 入账后紧邻的 `:1110 u.run.bookWaitingState("ui-prompt")` → `approval_always.go:171` 落 `approval: WAITING-STATE state=… corr=…`；事件侧 `run.go:1156` 把等待结束也记下来。UI 侧沿用快照既有 `pending` 键，**没加新 JSON 键** |
| **AC#7 无界面兜底** | **met（201-r1 已成立，本腿未削）** | `cmd/wisp/main.go:147 interactiveStdin()`→`run.go:601 attachReplyListener`；没有面时监听器不起，卡片按各自极性到点处理 |

---

## ③ 这轮动的文件（八枚改动＋五枚新件）

| 文件 | 动了什么 |
|---|---|
| `internal/agent/approval/replies.go` **新** | `Replies`/`ReplyCard`/`HostBinding`/`MaxTrackedCards`＋三枚接缝级哨兵错（`ErrNoTrackedCard`/`ErrRouteHasNoAllow`/`ErrNoGateAttached`）。它只调 gate 既有两台路由器，不新开第三道权威、不签发令牌 |
| `internal/agent/approval/replies_201_test.go` **新** | 包内白盒四枚（WaitingState 两态分开、未绑定时任何答复都是拒、未知编号、WideningRule 五种形） |
| `internal/config/allowdirs.go` **新** | `[fs] allowed_dirs` 的程序化写面（AddAllowedDir/SetAllowedDirs/共用的 writeAllowedDirs＋statOwnWrite） |
| `cmd/wisp/approval_always.go` **新** | `always` 动词的实现：规则文本→监听器续体→第二张 L2 卡→批准才写；`bookWaitingState`/`waitingStateName` |
| `cmd/wisp/approval_always_201_test.go` **新** | AC#4/AC#5 两支（只差第二张卡的答案） |
| `cmd/wisp/approval_seam_201_test.go` **新** | AC#1/AC#3/AC#6 的"没有 console"两枚 |
| `cmd/wisp/approval_reply.go` | 账本改为**委托**接缝（`liveCard = approval.ReplyCard`、`nativeCards` 成 wrapper）；路由调用换成 `s.live.h.*`；`always` 进动词表与 help；`bind` 在监听器接入时带上 vetoChannel 与两个来源标签 |
| `cmd/wisp/run.go` | 账本 `record(p)`（整张 Prompt 入账，Paths 不再漏）；gate 装配后立即 `bind`；`consoleApprovalUI` 带 `run` 回指并在入账后记等待态；"ticket 211"→201 |
| `cmd/wisp/approval_reply_211_test.go` → `…_201_test.go` | `git mv`＋注释/编号/`t211-`→`t201-`（正文与断言语义一字未改） |
| `internal/agent/approval/gate.go`、`queue.go`、`cmd/wisp/main.go`、`approval_reply_stdin_{windows,other}.go` | 只把错编号 211→201（这些注释是 201-r1 为自己那批产码写的） |

**没动的**（禁区核对）：`internal/perm/ticket90_persist_test.go`、`internal/panel/tokens_fourway_test.go`、`internal/panel/l2_grant_boundary_test.go`、`docs/PLAN.md`、`docs/specs/**`、`thresholds.go`、golden、`allowlist.txt`、`frontend/**`、`design/**`；C17/面板入向名册一字未加（`always` 是终端动词，不是线协议方法名）。

---

## ④ 判据与正控/反控

1. **AC#1 的"能力"判据不是"测试能碰到"**：`TestNativeHostSeamAnswersAnL2CardWithoutAnyConsole` 把 `h.reply` 留 nil（监听器不装配，`attachReplyListener` 根本不执行），宿主直接调导出方法 ⇒ 若把 `replies.go` 删掉这枚编译不过；若把账本↔路由的绑定拆掉，答语就回不了。同场还反证 console 没参与：审计里出现 `approval: ANSWER-ALLOW` 而**没有**任何 `approval: REPLY ` 前缀行。
2. **AC#3 反控**：`TestNativeHostSeamRefusesAPanelSourcedAllow` 三步（拒 allow→烧令牌→仍能 reject），并把"目标文件不存在"作为落盘级判据。
3. **AC#4 正控**："去掉那行文本即红"：`driveAlways` 必须先等到屏幕出现 `要存的规则：` 与 `allowed_dirs += "`，否则整枚用例 fail；第二张卡的 Reason 也要出现 `这是「一直」要求的 L2 级重新确认`。
4. **AC#4/AC#5 反控**：第二张卡答 `no` ⇒ `config.toml` 不含那一行、目标文件不存在、`WIDEN-REFUSED` 有声。
5. **AC#6 正控（拿掉生产者即红）**：`TestReplySeamWaitingStateNamesTheTwoD43Rows` 逐点判：空账本不报名 ⇒ `L1` 只报 `Confirming` ⇒ 加一枚 `L2` 变 `AwaitingApproval` ⇒ 撤掉 `L2` 退回 `Confirming` ⇒ 撤干净又回到不报名。`Forget`/`Record` 任何一枚被删，第二点就红。运行侧另有 `approval: WAITING-STATE state="AwaitingApproval" … corr="t201-seam-corr"` 在整场二进制里出现。

---

## ⑤ 最后一寸（AC#1 需要的宿主面）具体缺什么 —— 本腿没有替它打勾

接缝已经就位，宿主面缺的是**三枚具体件**，都不在票 201 的写面里：

1. **球／托盘进二进制**：`cmd/wisp/run.go` 里没有 `ball.New(...)`（现量见 §①）；`internal/ball/tray_windows.go` 的 `menuOpenPanel` 那项还是空 stub，且托盘菜单里没有"批准/拒绝"条目。要落的是：GUI 装配腿在 STA 线程上起球，把 `*approval.Replies`（run 已经造好并绑过 gate）交给点击/菜单处理器，处理器调 `Allow/Reject/Veto` 并按 `WaitingState()` 驱动 `Ball.SetState`。
2. **面板入向**：`internal/panel/bridge.go` 的入向名册只有 4 枚、没有 `panel.approval.request`；票面写明这枚对不上**只登记不改**（界面侧由 owner 带过去），且新增方法名要先落一条 `A##`（A273/A389 口径）。⇒ 面板那"一枚入口"在两件事落地之前不可能有生产调用者，本腿不报它。
3. **唤醒词否决（KWS）**属语音那族票，本票不动。

一句总结：**"至少两枚入口"今天仍不成立；成立的是"入口的形状已经可以是导出接缝，第二枚只差宿主面，不再差 Go 侧的 API"。**

---

## ⑥ 没做、做不动、需要人拍板的

| 项 | 状态 |
|---|---|
| 长期规则对**本次运行**生效 | 没做。`cmd/wisp/run.go:386-389` 的 `tools.NewPathCanonicalizer` 在装配时读一次 `[fs]`，本腿写的是文件，卡片与审计都明说"下一次启动生效"。要做"当场放宽"得再造成本器重建，属新功能，不混进本票 |
| `config.Manager.ConfirmLocked` 仍为 `nil`（`run.go:341 config.NewManager(cfgPath, nil)`） | 没改。含义：手改放宽 `[fs]` 走的是 **deny（不静默生效）**，本腿的程序化写走"调用方先起 L2 卡"这一支（`permmode.go:54-63` 已把这套先例写死）。要不要把 nil 换成真起卡的钩子，属 D36/热加载那一族的定案，本腿不自作 |
| 认领 mtime 这件事在这枚二进制里今天其实用不上 | 写下来免得下一位误读：`allowdirs.go` 的 `statOwnWrite` 只在 `CheckAndReload` 被调时才有意义，而编排者票 223 的普查（`172bda57`）量到 `CheckAndReload` 的唯一非测试调用者是 `cmd/balldebug/main.go:243`。⇒ 留着是因为它与 `SetPermissionMode` 的既有形状一致、且在常驻/球那侧一旦接上热加载就立刻需要；**不是**本腿用来支撑"不静默生效"的论据。本腿的论据只有一条：写之前那张 L2 卡必须已经被答，判据是 `approval: REPLY WIDEN-APPLIED` 只在 `PendingApproval` 返回 allow 之后出现 |
| D38b roster 少一枚名 | 长期支的续体沿用 `approval-waiter` 这个名字（`internal/observe/goroutine.go:52`）。新造一枚 roster 名是契约面，本腿不造，只在代码里写明原因 ⇒ 如果要让"widening worker"单独计数，需要 owner 动 D38b |
| `approval_grant` 生产写手 | 仍为零（票面 AC#5 原句已被编排者改写为 allowed_dirs 那一支，本腿按改写后的口径交） |
| `internal/panel` 的 4 枚历史红 | 未碰、未修绿（见 §⑧ 读数） |

---

## ⑦ 我与派单的不同意见（照例优先于同意）

1. **"面板 WebView 三处仍零"这句在票里把两件事当一件事**。球/托盘是"面存在但没编进这枚二进制"；面板是"宿主面根本不存在（ticket 33 未领）＋入向名册缺一枚且缺名字要先落 `A##`"。所以 AC#1 的"第二枚入口"最省的落点是**球**，不是面板；派单"面板算一枚"的口径如果照打，就是在给一条今天接不通的线打勾。我按"球或托盘"这一支交接缝，没报面板。
2. **派单暗示"缺的是 API"，实测缺的不是 API**：`Gate.DecideFromNative/DecideFromPanel` 一直是导出的。真正挡住宿主面的是"卡片↔一次性令牌"的账本住在 `package main`。这条差别决定了本腿为什么写 `replies.go` 而不是再包一层新接口——**并且我没有新增第三道权威**，`Replies` 只有转发与读。
3. **AC#2 的"等待态可见"我顺手做到了审计行里，但派单把它算进 AC#6**。两格其实同一枚 producer；我在 AC#2 只声称"能被否决＋到点按 §2 执行"，没声称任何极性变化。
4. **派单说 AC#4 若超尺寸就停下来报缺**。我判断不超：写面 3 枚函数＋监听器续体，都复用现成机制（`PendingApproval`/`SaveFile`/`AdmitTextTask`，全是 `confirmModeSwitch` 用过的）。但**"当场生效"我明确没做**（§⑥ 第一行），如果编排者认为 长期 必须当场生效，那这格应退回而不是算 met。
5. 一处我自己踩到的实现坑，记下来给下一腿：**AC#6 的生产调用者位置不是"能跑就行"**。第一版我放在 `ui.Prompt` 末尾（publish 之后），`TestNativeHostSeam…` 当场读到 `tracked=0` ——不是 bug，是**真并发**：钩子一探到账本就答 allow，`EventDismissed→forget` 与 Prompt 尾部并发。挪到"入账之后立刻记"之后语义才对（那张卡存在的那一瞬间有人在等）。这条读不出来的人会以为接缝在漏记。

---

## ⑧ 门禁读数（本腿自己跑，终态）

命令与结果，2026-09-29 11:4x +08，起手锚点 `93ea9f6b`、HEAD 为 `5a251ece`＋本证据件：

```
PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test ./cmd/wisp ./internal/... -count=1
  → EXIT=1（终态，包含全部包，没有用 -run 单测冒充包判）
  ok   github.com/CarlosShao/wisp/cmd/wisp                      85.753s
  ok   github.com/CarlosShao/wisp/internal/agent/approval       0.424s
  ok   github.com/CarlosShao/wisp/internal/config               0.906s
  ok   internal/{agent,llm,llm/*,memory,models,observe,perm,plugin,proc,projctx,risk,tools,…}   全 ok
  FAIL github.com/CarlosShao/wisp/internal/ball   0.190s   —— 1 枚：TestC21TableColourRowsMatchTokensCSS
  FAIL github.com/CarlosShao/wisp/internal/panel  2.566s   —— 4 枚：TestApprovalCardViewJSONKeysMatchFrontendTypes /
       TestComposerContractTypesMatchFrontend / TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme /
       TestC21DesignTokensFourWayAgree
```

**两枚红的归属，逐条查过而不是推断**：

- `internal/panel` 那 4 枚＝编排者派单里点名的历史在册红（HANDOVER §4.0v），本腿没修绿也没新增。
- `internal/ball` 那 1 枚的失败原因逐字是 `tokens_table_test.go:1468: read design/assets/tokens.css: open …design\assets\tokens.css: The system cannot find the path specified` —— **别人正在飞的 `design/**` 删除**（本腿第一条命令 `git status --porcelain` 起手就记着 ` D design/assets/tokens.css` 等 13 枚删除，早于我任何一次写入）。派单禁我碰 `design/**`，所以这颗红我不动、只具名上报；同理 `TestC21DesignTokensFourWayAgree` 这一枚现在也带上了同一条缺文件的因由（它原本的红是名册对不上）。
- 除这两包外 `./internal/...` 全 ok；`cmd/wisp` 全绿（带 sherpa PATH；这条是票 98 那枚 `0xc0000135` 假绿的反面）。

其它两把尺：

```
gofumpt v0.12.0 (go1.27.1) -l <本腿 12 枚改动文件>   → 空输出（净），exit 0
tools/d22scan/d22scan.exe                             → clean - no D22 ban violations
     examined 246 production Go files under internal/ and cmd/
     ban #8 internal/=456 / cmd/=58 Go files（注释与 _test.go 也进射程）
go build ./...                                        → 净
go vet ./cmd/wisp ./internal/agent/approval ./internal/config → 净
```

新增用例数（本腿）：包内白盒 4 枚（`replies_201_test.go`）＋二进制判据 4 枚（`approval_seam_201_test.go` 2 枚、`approval_always_201_test.go` 2 枚）＝8 枚；201-r1 那 5 枚一字未改地仍然绿。
