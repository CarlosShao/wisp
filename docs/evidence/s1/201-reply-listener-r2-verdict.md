# 票 201 · 续腿 `201-r2` 独立验收裁决（裁决者 ≠ 实现者，D22 双角色）

- 裁决时刻起手：`2026-09-29 11:57 +0800`（`date` 现跑），收尾各节各自带自己的 `date`。
- 锚点：`git rev-parse --short HEAD` 现跑＝**`7096470d`**。
  ⚠ 被审的产码链尾是 `36b46125`，HEAD 已因别的腿（223-c1 / 220-c1 / 224-c1）前移；
  我用 `git diff --name-only 36b46125 HEAD -- cmd/wisp internal/agent/approval internal/config` 现跑＝**空输出**
  ⇒ 我下面所有行号读的就是这批产码在盘上的现身，不是归档。
- 被审件（本腿五枚 commit）：`58bf0158`（编号 211→201）→ `49eb440b`（接缝＋写路径＋`always`）→
  `5a251ece`（判据＋AC#6 落点）→ `6b026854`＋`36b46125`（证据件）。
  新件：`internal/agent/approval/replies.go`、`replies_201_test.go`、`internal/config/allowdirs.go`、
  `cmd/wisp/approval_always.go`、`approval_seam_201_test.go`、`approval_always_201_test.go`。
- 我的尺只有：Read / Grep / Glob / `sed -n` / `ls` / `date` / `git log|show|diff --stat --name-only|cat-file`。
  ⛔ **未跑 `go test`／`go build`／`go vet`／`d22scan`／任何编译**（写腿 `222-r1` 在动 `internal/tools/**`）。
  ⛔ **未读、未引、未转述 `frontend/**` 与 `design/**`**。
- 门禁读数：一律引编排者的 `.scratch/wisp/probes/201/r2/gates-full.txt`（5,859 字节、72 行，
  文件 mtime `2026-09-29 11:52`，内里日志时间戳 11:50:53–11:51:02 与编排者报的窗口一致）。
  **这枚不是我现跑的**，按硬约束我也不自产；但我把它逐行读了，**读出的数与编排者转述给我的数不一样**，见 §① 第 1 条。

---

## ⓪ 门禁读数核对（先把地基钉住，再谈逐格）

| 我在这枚文件里现读到的 | 读数 | 行号 |
|---|---|---|
| `cmd/wisp` | ok **93.362s** | `gates-full.txt:1` |
| `internal/agent/approval` | ok **0.576s** | `:3` |
| `internal/config` | ok **1.237s**（编排者转述给我的是 1.186s） | `:19` |
| `internal/memory` | ok **14.525s**（转述 14.655s） | `:26` |
| `internal/perm` | ok **0.671s**（转述 0.512s） | `:53` |
| `internal/tools` | ok **17.972s** | `:69` |
| `^ok` 行总数 | **23 枚**（与转述一致） | `grep -c "^ok  "` ＝ 23 |
| FAIL 的**包** | **三枚**：`internal/ball`（1 例，`:14-17`）、`internal/panel`（4 例，`:30-52`）、**`internal/risk`（1 例，`:58-63`）** | — |

`internal/risk` 那枚逐字：

```
--- FAIL: TestResolvePerCallBudget (2.24s)
    pathresolver_budget_norace_test.go:34: C26 Resolve: 1222225 ns/op = 1.222 ms/op (budget 1.000 ms, 1689 samples)
    pathresolver_budget_norace_test.go:37: C26 budget breach: Resolve averages 1.222225ms per call, budget 1ms
```

⇒ **编排者那句"FAIL 只有 `internal/ball` 1 例与 `internal/panel` 4 例"与盘上不符**。三枚红里：
`internal/ball`/`internal/panel` 的逐字因由确实指向 `design/assets/tokens.css` 缺失与面板契约字段（`tokens_table_test.go:1468`
原文我在门禁文件里读到，不是我去了 `design/**`）；**`internal/risk` 那枚是新出现的、没人归因的一枚**。
本批产码与此的关系我核过：**五枚 commit 的文件清单一行都没有 `internal/risk` 与 `internal/tools`**
（`git show --stat 58bf0158 49eb440b 5a251ece` 现读），且 `TestResolvePerCallBudget` 量的是 C26 每调用墙钟预算，
不读 `[fs] allowed_dirs` 的程序化写面 ⇒ **我判这枚红不归本批**；但我也**不能**判它是 flaky——
我没有任何空载读数。它的归属是编排者的活，我这一格只登记"转述失真＋一枚无人归因的红"。

**这一节对本批产码的判语：不因此批失败，但"本批无新增红"这句结论在 `internal/risk` 那枚上今天没有证据。**

---

## ① 逐格裁决（票面 AC 表 `.scratch/wisp/issues/201-…md:41-47` 逐枚对文字，编排者 09-29 改写段 `:36`/`:37` 一并对照）

起手时间戳 `2026-09-29 12:00 +0800`。每格末列**恒真检查**＝"把实现拿掉这格会不会红"，不会红即不算满足（票面硬约束 7）。

### AC#1「有人能答：三枚入口至少两枚有生产调用者」

- 实现者声称：`replies.go:313 Allow`／`:336 Reject`／`:419 Veto`／`:342 PanelReject`／`:395 PanelAllow`／`:167 Record`
  已成"可交给宿主的导出接缝"，生产构造点 `cmd/wisp/run.go:458 rt.liveCards.bind(rt.gate, "", "", "")`，
  console 是第一个调用者（`approval_reply.go:210 / :249 / :251 / :278 / :310`）；**并自陈"入口枚数仍＝1"**（证据件 `:31`、`:78`）。
- 我现跑的尺与读数：
  - `grep -rn "NewReplies|approval.Replies|\*Replies" --include=*.go cmd/ internal/` ⇒ **非测试只有两枚命中**：
    `cmd/wisp/approval_reply.go:111 type nativeCards struct{ h *approval.Replies }` 与 `:113 newNativeCards()`。
    ⇒ 全仓今天**只有一个持有者**，而那一个的取用路径全在 console 语法里（`approval_reply.go:469-489` 的动词表）。
  - `grep -rn "DecideFromNative|DecideFromPanel|\.Veto(" --include=*.go internal/ cmd/ | grep -v _test.go` ⇒
    生产调用点全部落在 `replies.go:325/:372/:377/:401/:427`，**再往上一层的入口只有 `approval_reply.go:310`（console）**。
  - `grep -rn "internal/ball" --include=*.go cmd/ internal/ | grep -v _test.go` ⇒ 二进制引用者只有 `cmd/balldebug/main.go:27`；
    `cmd/wisp/**` 里 `internal/ball` 只出现在注释（`approval_always.go:163`、`notify_windows.go:13`）。
  - `sed -n '80,90p' internal/ball/tray_windows.go` ⇒ 托盘菜单项只有 `appendItem(menuOpenPanel, "打开面板", false)`（`:86`）
    加分隔线，**没有批准/拒绝条目**，且"打开面板"自陈是 no-op stub（`:5`）。
  - `grep -n "no WebView2 host" internal/panel/pump.go cmd/wisp/panel_pump.go` ⇒ `pump.go:16`、`panel_pump.go:13` 双处自陈无宿主。
- **判语：不成立（票面文字＝"至少两枚"，读数＝1 枚）。** 但这是**诚实申报的不成立**，不是虚报：
  实现者在证据件 `:31` 与 `:78` 两处白纸黑字写"仍是 1 枚"，并在 `:70-78` 列出缺的三枚具体件，没有替它打勾。
- 恒真检查：**会红**。把 `replies.go` 删掉 ⇒ `approval_seam_201_test.go:89` 编译不过；
  把账本与路由的绑定拆掉（`run.go:458` 那行删）⇒ 同一枚测试的 `Allow` 拿到 `ErrNoGateAttached`。
  ⇒ 这一格"接缝能力"是真测到了的；**没测到的是"第二枚入口"，而票面要的正是枚数**。
- 翻勾前置具名：**球／托盘宿主装配腿**（把 `*approval.Replies` 交给点击/菜单处理器，票面 §要建什么 第 1 段 `:25` 原句"先接球与托盘这两枚"）；
  面板那一枚另需 `internal/panel/bridge.go` 入向名册补 `panel.approval.request` **且先落一条 `A##`**（票面禁区 `:51`）＋票 33 的 WebView2 宿主。

### AC#2（改写句：L1 等待期间能被否决＋等待态可见；到点无答复仍按 SPEC-06 §2 执行）

- 实现者声称（证据件 `:32`）：否决走 `replies.go:419 Veto`→`Gate.Veto`；到点无否决仍执行，
  `internal/tools/bridge.go:384-386` 的冻结语义**一字未动**，极性未翻。
- 我现跑的尺与读数：
  - `git show --stat` 三枚产码 commit ⇒ **文件清单无 `internal/tools`、无 `internal/risk`**（桥那三行注释不在本批射程）。
  - `git show 58bf0158 -- internal/agent/approval/gate.go internal/agent/approval/queue.go` 的 `^[+-]` 行 ⇒
    **4 行全是注释里的 "211"→"201"**，`gate.go`/`queue.go` 语义零改动（实现者"只改编号"的说法复核成立）。
  - 否决这一支的真实装配：`run.go:146 replyVeto approval.Channel`（字段）＋`run.go:601 attachReplyListener(s.reply, s.replyVeto)`；
    `grep -rn "replyVeto" --include=*.go cmd/ | grep -v _test` ⇒ **只剩字段声明与那一处读取，没有任何生产赋值**；
    `grep -rn "replyVeto" --include=*_test.go cmd/` ⇒ 唯一赋值在 `approval_reply_201_test.go:574 h.replyVeto = approval.ChannelEsc`。
  - `sed -n '387,409p' internal/agent/approval/gate.go` ⇒ `Veto` 先 `g.channels.check(v.Channel)`，通道不可达时对**活着的窗口**发
    `EventWarning` 并 `return err`（拒，且说得出声）。
  - `run.go:436-440` 逐字自陈："What it still cannot do is veto an L1 window … so `NewChannels()` stays empty and
    `runSpec.replyVeto` stays unset"；`approval_reply_201_test.go:565-568` 逐字自陈那半支测试是
    "the seam **standing in** for the ball/Esc-hook leg (tickets 07/77/92)"。
- **判语：一半成立、一半"能测到但不是生产"。**
  ①「到点仍执行、极性未翻」＝**成立**，且有反控（`approval_reply_201_test.go:492` 第一支断言
  `out.IsError` 必须为假＋审计含 `ANSWER-EXPIRED … decision=timeout->execute`，`:543-545` 断"没落地的否决不得被记成落地的"）。
  ②「L1 等待期间能被否决」＝**机制成立、生产不可达**： shipped 的 `wisp run` 里那一发 veto 必然被 `channels.check` 拒掉。
- 恒真检查：**两半都会红**。删 `replies.go:419 Veto` ⇒ `replies_201_test.go:96` 与 console 那支不编译；
  把 `bridge`/窗口的到点极性改成落拒绝 ⇒ `:538-541` 的 `t.Fatalf("the L1 timeout polarity changed …")` 红。
- 对翻勾的意见：**带注可翻，但注必须写成逐字那两句**（"仓内无任何 host 声明 veto 通道，
  `replyVeto` 生产零赋值 ⇒ 今天人能否决这一发在二进制里不可达"）。编排者的初判"可翻（带注）"我**不推翻**，
  但**证据件 `:32` 那一行没有写这句自陈**（写在了测试文件里），下一位读表的人会把"能被否决"读成已经能用——这是这格唯一的真缺陷。
- 翻勾前置具名：**球/Esc 宿主腿**（票 07/77/92 那族）把 `replyVeto` 真填上。

### AC#3「L2 批得动；同一发从界面方法名直接送进来的允许要被判红」

- 实现者声称（证据件 `:33`）：批得动＝`approval_seam_201_test.go` 第一枚（整场 `h.reply=nil`、宿主调 `Replies.Allow`、`fs.write` 真落盘）；
  判红＝第二枚（`PanelAllow` 返回 `ErrPanelAllow`、被烧的令牌再拿去原生 allow 得 `ErrBadGrant`、`PanelReject` 仍能拒）。
- 我现跑的尺与读数：
  - `sed -n '638,649p' internal/agent/approval/gate.go` ⇒ `DecideFromPanel` 在 `r.Allow` 分支**先拒再谈令牌**，
    且 `r.Grant != ""` 时 `g.q.revokeGrants(...)`（按泄露烧掉），返回 `ErrPanelAllow`。
    ⇒ 拒的是**路由**，不是便利，符合票面 `:21-22` 那句硬禁与 `PLAN.md:2027` 定案的"严格版"。
  - `approval_seam_201_test.go:89` 真的走 `rt.liveCards.h.Allow(...)` 且 `h.reply` 留 nil（`:52` 注释逐字）；
    `:113-116` 断目标文件**真的存在且内容对**；`:128-130` 反向断"审计里不得出现 console 的 `approval: REPLY ` 前缀行"。
  - `:168-186` 三步反控（拒 allow → 烧令牌 → 仍能 reject）＋ `:199-204` 断"这一发必须没写成文件"。
  - `replies.go:395-415 PanelAllow` 的返回形状（`(bool, error)`，`nil error` 视为安全故障）＋ `approval_reply.go:281-286`
    把"竟然没被拒"写成显式故障句——**这条 happy path 不存在**，不是有一条宽松分支。
- **判语：成立。** 与编排者初判一致，我没有找到可以放水的地方。
- 恒真检查：**会红**。把 `gate.go` 的面板 allow 拒枝拿掉 ⇒ `:169 errors.Is(err, ErrPanelAllow)` 红、
  且 `approval_reply.go:285` 那句故障会被触发；把 `Replies.Allow` 的 grant 传递拿掉 ⇒ `:89` 红。
- 翻勾前置具名：**无**（这一格可以直翻）。⚠ 但见附 A-2：这格的"面板侧不得允许"今天靠的是
  `internal/panel` 里那把冻结尺的**射程之外**的一枚 Go 方法，不是那把尺本身。

### AC#4「"一直"看得见存了什么：选"一直"后卡上出现规则文本＋二次确认；正控＝去掉那行文本 ⇒ 红」

- 实现者声称（证据件 `:34`）：`approval_always.go:70 always` 先印规则文本（`:103-106`，文本由 `replies.go:467 WideningRule`
  从卡片自己的 Paths 推），再 `:112 runWidening` 用 `gate.PendingApproval`（`:115`）真起第二张 L2 卡，批准才 `:134 AddAllowedDir`。
  并在 `:100 不同意见 4` 主动写："**如果编排者认为长期必须当场生效，那这格应退回而不是算 met**"。
- 我现跑的尺与读数（票面逐字＋代码＋测试三处对齐）：
  - 票面 `:44` 逐字＝"选"一直"后**卡上出现规则文本**＋二次确认；正控＝去掉那行文本 ⇒ 红"。**这一句里没有"当场生效"四个字**。
    票面 §要建什么 第 3 段（`:28`）逐字＝"卡片上有这次/一直/拒绝三枚语义，且选'一直'时**把会存成哪条规则印出来＋二次确认**；
    '一直'落到 D45 已设计的那张 `approval_grant` 表"——**也没有生效时机的要求**，落点那一半已被编排者 `:37` 改写走。
  - 我 grep 了整张票与 `PLAN.md`：`grep -n "热加载放宽|不得静默生效" docs/PLAN.md` ⇒ 只有 `:1644-1645` 一处，逐字
    "**D33 配置提权**：热加载放宽 `[risk]`/`[fs]`/`[net]`/`[plugins]` → **必须触发 L2 级重新确认，不得静默生效**"。
    **这句话约束的方向是"不许未经确认就生效"，不是"必须当场生效"**；本实现两个方向都满足：
    (a) 写之前真有一张 L2 卡（`approval_always.go:115-124`，`ans != AnswerAllow` 就走 `:128 WIDEN-REFUSED`）；
    (b) 本次运行仍按旧名单（`:143` 逐字"下一次启动生效，本次运行仍按旧名单"，
    代码依据我复核为真：唯一的运行期读者 `run.go:386-390` 是在装配期把 `cfg.FS.AllowedDirs` 逐枚拷进**新切片**再交给
    `tools.NewPathCanonicalizer`，`AddAllowedDir` 换掉的是 `m.cur` 的切片头（`allowdirs.go:110`），碰不到那份快照）。
  - 测试三处硬钉（`approval_always_201_test.go`）：`:208-209` 必须先在屏幕上等到 `要存的规则：` 与 `allowed_dirs += "`，
    否则整枚用例 fail；`:211-215` 第二张卡必须真的出现（`waitWidenCard187` 只认 `c.Tool == allowWidenTool`，`:250-264`）；
    `:218` 第二张卡的 Reason 必须含 `这是「一直」要求的 L2 级重新确认`；`:181-189` 反控——第二张卡答 `no` ⇒
    `config.toml` 不含那一行＋目标文件不存在＋审计出现 `WIDEN-REFUSED`。
  - 规则文本的推导不是自由发挥：`replies.go:467-483` 只在**去重后恰好一个目录**时产出，目录不唯一（`:478`）、
    裸盘根（`:473 isBareDrive`）、空 Paths 一律拒，`replies_201_test.go:114-143` 五种形逐枚钉住。
- **判语：成立（按票面文字）。退回不成立。**
  我对编排者这条判定的复核结论＝**同意**，但同意的**依据**要摆明："立即生效"既不在 AC#4 的票面文字里、也不在 §要建什么 第 3 段里，
  而 `PLAN.md:1644-1645` 那句反而**只禁止"未经确认就生效"**；把 D33 那句读成"必须当场放宽"是**给冻结文字加它没写的内容**，
  那种读法如果成立，就该由 owner 落一条 `A##`/改票面，而不是由验收侧在翻勾时补。**我不接受"因为代码自己承认没做就退回"这条推路**——
  代码自陈没做的东西，只有票面要求了才构成退回；票面没要求、且明确属于另一格（热加载落点）的功能，正确处置是**登记**，不是退回。
- 恒真检查：**会红**，而且是这格里最硬的一枚。删掉 `approval_always.go:103` 那句返回文本 ⇒ `:208` 红；
  删掉第二张卡（`:112-123`）⇒ `:211` 红；把"批准才写"改成"直接写"⇒ `:181` 红。
- 翻勾前置具名（不是本格的退回条件，是本格之外的两笔要登记）：
  **①"当场生效"＝票 223 那一族（`CheckAndReload` 的落点普查已由 `172bda57` 量清）＋ D36 生效级别的定案**——
  必须落成 `DEFERRED(D-xx)` 或一条 `A##`，不能只写在证据件 `:86` 的一行里（票面硬约束：SPEC-12 §5 五字段登记）；
  **②"可撤销"那一半在本格没门**，见下一格。

### AC#5（改写句：长期那一支真落库＝往 `[fs] allowed_dirs` 加一行，持久、**可撤销**，并按 `PLAN.md:1645-1646` 触发一次 L2 重新确认；`approval_grant` 的 session 那一支必须不活过重启）

- 实现者声称（证据件 `:35`）：落的是 `config.toml` 的 `[fs] allowed_dirs`（`allowdirs.go:56 AddAllowedDir`／`:87 SetAllowedDirs`）；
  `approval_grant` **本腿没有加生产写手**；`ticket90_persist_test.go` 三枚钉子未碰。
- 我现跑的尺与读数（三段分开判）：
  - 写手在、且真落盘：`approval_always.go:134` 是 `AddAllowedDir` 的**唯一生产调用点**（全仓 grep 只有这一处＋定义处），
    测试 `approval_always_201_test.go:135-142` 断 `config.toml` 里出现那一行。⇒ **"长期那一支真落库"成立**。
  - 重新确认在：上一格已钉。⇒ **成立**。
  - **"可撤销"不成立**：`grep -rn "AddAllowedDir|SetAllowedDirs" --include=*.go .` ⇒
    **`SetAllowedDirs`（作者自己写明的那枚"撤销半"，`allowdirs.go:79-87`）今天零调用者、且零测试**；
    答复语法里也没有撤销动词（`approval_reply.go:469-489` 的动词表＝`yes/no/veto/panel-no/panel-yes/always/head/view/help/quit`，
    `grep -n "never|revoke|撤销" cmd/wisp/approval_reply.go cmd/wisp/approval_always.go` 命中的全是注释散文）。
    ⇒ 用户今天存了规则之后，**在这个二进制里没有任何一条能把它撤回去的路**；改写句的"可撤销"三个字没有落点。
  - **session 那一支"不活过重启"＝恒真**：`grep` 现读 `approval_grant` 的 INSERT 只在 `internal/memory/dao_misc.go:33`（DAO 层），
    而这枚 DAO 方法（`InsertGrant`）的非测试调用者＝**零**（`grep -rn "InsertGrant|RevokeGrant" --include=*.go cmd/ internal/`
    去掉 DAO 自身后**全部命中 `_test.go`**：`run_mode101_test.go:411`、`dao_test.go:383/391`、`ticket90_persist_test.go:230`）。
    ⇒ "没写所以不会跨重启"是被"根本没写"满足的，**按票面硬约束 7 这一半不算满足**。编排者这条我完全同意，且我的读数支持他的用词。
- **判语：不翻**（结论与编排者一致，**依据不同**：编排者用的是票面原句的尺；按编排者自己 `:37` 改写后的句子，
  这格真正的缺口是**"可撤销"缺一枚门**＋session 那半是恒真，而不是"表里没有写手"这一条本身）。
- 恒真检查：长期那一支**会红**（删 `approval_always.go:134` ⇒ `:140` 红）；
  "可撤销"那一半**今天没有任何尺会红**——这正是它不能算满足的原因。
- 翻勾前置具名：一枚**撤销门**（最省＝在答复语法里加一枚与 `always` 配对的动词，走同一张 L2 重新确认＋`SetAllowedDirs` 收紧；
  或明确登记为 DEFERRED 并具名转给下一腿）。**这枚门属于票面 §要建什么 第 3 段（`:29`）的"落到 D45 那张表"那一支的替代形状**，
  不是我新加的额外要求。

### AC#6「等待态可见：等人答复时球进"等人"那一态；把这一态的驱动拿掉要能被判红」

- 实现者声称（证据件 `:36`）：读侧 `replies.go:249 AwaitingHuman`／`:283 WaitingState`（只回 D43 冻结名），
  生产调用者 `run.go:1103` 入账后紧邻的 `:1110 u.run.bookWaitingState("ui-prompt")` → `approval_always.go:171`，
  事件侧 `run.go:1156`；**并自陈"met-with-seam"、UI 侧没加新 JSON 键**。
- 我现跑的尺与读数：
  - `sed -n '1095,1115p' cmd/wisp/run.go` ⇒ `consoleApprovalUI.Prompt` 里 `u.live.record(p)` 之后**紧邻**
    `u.run.bookWaitingState("ui-prompt")`，注释逐字"the instant a card exists is the instant someone is being waited on,
  so that fact is booked HERE"。事件侧 `:1148-1157`：`EventDismissed/EventStarted` 先 `forget` 再 `bookWaitingState("ui-"+Kind)`。
  - `sed -n '171,184p' cmd/wisp/approval_always.go` ⇒ 两态都落审计行，未绑 gate 时 `state=none awaiting=false`，
    有人等时把 `corr/tool/level/paths` 一起报名；名字来源 `statemachine`（`:190 waitingStateName` 返回 `statemachine.State`）。
  - `WaitingState` 只可能回两枚名（`replies.go:291 StateAwaitingApproval`／`:295 StateConfirming`），
    `replies_201_test.go:40-69` 用**五点点判**钉住"两态分开＋撤掉会退回来＋撤干净不报名"。
  - 球那一侧我复核"面存在但没编进来"：`internal/ball/statevisual.go` 我不读（不在禁区，但本格只需知道引用面），
    需要的是 `cmd/wisp` 不 import 它（`grep -rn "internal/ball" cmd/ | grep -v _test` ⇒ 只有 `cmd/balldebug/main.go:27`）。
- **判语：成立但只成立一半——"驱动"第一次有了生产调用者，"球真的变那一态"没有。**
  票面动词是"球进'等人'那一态"：**严格照文字这格不算全绿**；照实现者交的形状（同一枚 producer 报名字、
  由任何宿主面去画）它是这格里除 AC#3 之外最扎实的一枚。⇒ **可翻，但注必须是逐字这句**：
  "今天红的是'名字有人报'，不是'球真的变'；球接入后需要的调用（`WaitingState()`）已在 `replies.go:283` 就位。"
- 恒真检查：**会红**。(a) 删 `approval_always.go:171 bookWaitingState` ⇒ `approval_seam_201_test.go:120`
  断的那条审计行 `approval: WAITING-STATE state="AwaitingApproval" … corr="t201-seam-corr"` 缺失；
  (b) 删 `run.go:1110` 那一处调用 ⇒ 同一枚断言红（这枚测试整场 `h.reply=nil，`监听器不起，所以只能靠生产 Prompt 路径报名）；
  (c) 让两态塌成一枚名 ⇒ `replies_201_test.go:47/52/62` 任一点红。
- 翻勾前置具名：**球宿主腿**（把 `WaitingState()` 接到 `Ball.SetState`）；本腿之外，**不需要新的 Go 侧 API**。

### AC#7「无界面兜底：行为明确且响亮，不许静默执行」

- 实现者声称（证据件 `:37`）：`main.go:147 interactiveStdin()` → `run.go:601 attachReplyListener`；
  没有面时监听器不起，卡片按各自极性到点处理（201-r1 已成立，本腿未削）。
- 我现跑的尺与读数：
  - `sed -n '145,153p' cmd/wisp/main.go` ⇒ `reply == nil` 时向 stderr 逐字印
    "本机没有可交互控制台，本轮没有人能答复卡片：L2 卡会等到超时后按拒绝处理，L1 窗口没有人能否决（要能当场答复，请在终端里跑）"。
  - `approval_reply_stdin_windows.go:41`／`approval_reply_stdin_other.go:23` ⇒ 管道、重定向、非 Windows 一律返回 nil（不会把"有面"当成"有面"）。
  - 行为侧的钉：`approval_reply_201_test.go:429 TestUnansweredL2CardTimesOutIntoRejectNeverExecution`，
    `:487` 断 `tool_call` 行的 decision 必须是 C18 auto-reject 而不是 allow。
- **判语：成立（行为），但"响亮"那一半今天没有尺**。
  `grep -rn "本机没有可交互控制台" --include=*.go cmd/` ⇒ **只命中 `main.go:150` 自己，测试文件零引用**
  ⇒ 把那两行兜底提示删掉，**今天不会有任何一枚用例红**。按票面硬约束 7，这半只能算"做到了但没被钉住"。
- 恒真检查：**行为会红、句子不会红**（见上一条读数）。
- 处置：不是退回项（本腿没削它，201-r1 就是这形状），但我登记为**要补一枚尺**：
  断"`reply==nil` 时 stderr 必须出现那句"，落点自然在 `cmd/wisp` 的既有测试族里。归属＝编排者派下一腿时顺手带。

---

## 附 A｜编排者点名的两枚待裁点（我裁完的答复）

起手 `2026-09-29 12:05 +0800`。

### A-1 `internal/config/allowdirs.go` 的三条会不会吞掉用户的真实手改？

**逐条裁：**

1. **含 `..` 一律 fail-closed 拒（`:61-67`、`SetAllowedDirs` 的 `:94-97`）⇒ 不吞手改。**
   它只作用在**程序化写**这一支；`strings.Contains` 的射程确实过宽（目录名字面上含 `..` 的合法路径会被拒），
   但后果是**报错＋什么都不改**（`:69-77` 在 `m.mu` 里先去重再构造 `next`，`AddAllowedDir` 的 `:72` 命中已存在就是 no-op），
   用户的手改文件一个字节都不动。方向保守，判**可接受**。
2. **写失败回滚内存（`:111-115`）⇒ 不吞手改，且这是正确的一半。**
   `old := m.cur.FS.AllowedDirs` → `SaveFile` 失败 → 复位 → `observe.Wrap` 上报；
   我核过它与 `permmode.go:70-77` 的先例逐字同形。内存与文件不可能分叉。
3. **`statOwnWrite`（`:116`、`:125-131`）单独看不吞手改**：它认领的是**自己刚写完之后**的 mtime+size；
   用户在 `t2 > t1` 的手改仍然会让轮询看到 mtime/size 变化。**stat 本身无罪。**

**但三条叠起来有一条真的会吞——罪名不在这三条里，在 `SaveFile` 的写法上：**

- `writeAllowedDirs` 持久化的是**整个 `m.cur` 快照**（`:111 SaveFile(m.path, m.cur)` → `loader.go:135-145`：
  `deepCopyConfig(c)` → `MarshalCanonical` → `atomicWrite`），**写前不 re-read 磁盘**。
- 于是：进程起来后、`always` 落盘前，用户手改了 `config.toml` 的**任何**一处（另一个 section、或另一行 `allowed_dirs`），
  那一发程序化写会**把整份内存快照盖回文件**，手改就此消失；紧接着 `statOwnWrite` 认领的正是**覆盖之后**的那份 stat
  ⇒ 在会轮询的宿主（`cmd/balldebug` 那支）里，`CheckAndReload` 之后**永远看不见那次覆盖**。
  在 `wisp run` 里今天没有轮询（`grep -rn "CheckAndReload" | grep -v _test` ⇒ 唯一非测试调用者是 `cmd/balldebug/main.go:243`，
  与票 223 的现量一致），所以今天的实际暴露＝"运行期间的手改被覆盖、且没有提示"。
- **继承还是新造**：`SetPermissionMode`（`permmode.go:63-77`）是同一形状 ⇒ **不是本腿新造的病**，
  但本腿把暴露面从 1 枚扩到 2 枚，而票面、证据件 §⑥ 与代码注释**都没有登记这一条**。
  判：**票面没覆盖的风险，具名上报**（不构成本格退回；修法归 D36 那一族——写前 re-read 或改 patch-merge）。

### A-2 `SetAllowedDirs` 是不是可以被绕过的口子？仓里今天有没有尺会响？

**是口子；今天零尺。**

- 作者自己的话就在门上（`allowdirs.go:82-86`）："这枚方法**分不出收紧与放宽两种意图**，所以名字叫 SET"。
  **名字不是门。** 拿它加一行而不签卡，编译器不拦、注释不拦。
- 我逐枚查了仓里可能响的尺，读数如下：
  | 可能的尺 | 现跑读数 | 会响吗 |
  |---|---|---|
  | `internal/config` 自己的测试 | `grep allowed_dirs --include=*_test.go internal/config/` ⇒ 命中的是 `manager_test.go:162-244`（**reload/ConfirmLocked** 那支）、`migrate_test.go:92`、`unwired_test.go:192`；**没有一枚调用 `AddAllowedDir`/`SetAllowedDirs`** | ❌ |
  | `cmd/wisp` 的两枚 `always` 测试 | 只覆盖"走卡＋`AddAllowedDir`"那一条正路 | ❌（绕开卡片的那条不在射程） |
  | `tools/d22scan` 的禁形状 | 禁的是裸 `go func`、`filepath.*` 决策、明文密钥、墙钟超时、镜像取哈希、面板侧 L2 允许、artifacts 当受门控 Tool、emoji（`main.go:5-52` 的 ban 清单） | ❌ 无一条管"锁定 section 的程序化放宽未经确认" |
  | `internal/panel/l2_grant_boundary_test.go`（冻结尺） | 走查根现读 `:1231`、`:1574`、`:1848`、`:2023` **四处全是 `filepath.Join(root, "internal", "panel")`** | ❌ 看不见 `internal/config` 与 `internal/agent/approval` |
  | 本包内的文档性防线 | `allowdirs.go:16-24`「THE ONE RULE THIS FILE MUST NOT SOFTEN」＋"pair pinned by a test" | ⚠ 那句"pinned by a test"钉的是 `cmd/wisp/approval_always_201_test.go`（**调用方那一侧**），不是 `SetAllowedDirs` 本身 |
- 结论：**这枚导出的放宽门今天是"文档约束＋零仪器"**，且 `SetAllowedDirs` 连一个调用者都没有（＝纯粹为"撤销半"预留的导出面）。
  要么补一枚尺（我建议的最小形状：`internal/config` 里一枚断言"`AddAllowedDir` 与 `SetAllowedDirs` 放宽方向上
  必须让调用方先给出一条已定的 L2 答复，或干脆把 `SetAllowedDirs` 收成 non-exported／改名 `ReplaceAllowedDirsForRevocation`），
  要么把这枚导出不交付（本腿不删，我只具名）。**归属：票 201 的续腿或下一腿的登记项，不是编排者可以就地忽略的一条注释。**

### A-3 `ConfirmLocked` 生产零赋值 ⇒ 这算不算把 D33 的"reload 时复确认"偷换成"答复时复确认"？

- 先复核事实（编排者这条是对的，我重跑了一遍）：
  `grep -rn "ConfirmLocked = |ConfirmLocked:" --include=*.go cmd/ internal/ | grep -v _test` ⇒ **空**；
  唯一的生产赋值面是 `manager.go:78-88 NewManager(path, res)`——**它的第二枚参数是 `SecretResolver`，不是这枚钩子**。
  ⚠ 因此证据件 §⑥ 第 2 行那句"`config.Manager.ConfirmLocked` 仍为 nil（`run.go:341 config.NewManager(cfgPath, nil)`）"
  **括号里的因果是错的**（那枚 nil 是密钥解析器）：钩子为 nil 的真正原因是**全仓没有任何生产代码给这个公开字段赋过值**
  （`manager.go:240 approved := m.ConfirmLocked != nil && m.ConfirmLocked(section, loosen)` ⇒ nil 即 deny，fail-closed 成立）。
  这条是证据件的**文字缺陷**，不改判语，但会误导下一位，必须更正。
- **两种读法我都裁一遍：**
  - **读法一（"复确认"指的是 reload 那一刻）**：那么 D33 那句对这批改动**根本没有被触发**——
    `wisp run` 不跑 `CheckAndReload`，程序化写也不是热加载；那一句的守门人（nil 钩子→deny）在本腿前后行为完全一致，
    我核过 `manager_test.go:158-244` 那批 reload 钉子没被碰（`git diff` 无 `internal/config/manager*`）。
    ⇒ 在这一读法下**"偷换"不成立**，因为被换掉的那条路今天在这枚二进制里并不存在。
  - **读法二（"复确认"是对"放宽"这件事的意图要求）**：那本腿交的形状**恰恰是这句话要的形状**——
    放宽发生前，用户亲手答了一张 L2 卡（`approval_always.go:115-124`），且**不静默生效**（`:143`）。
    把确认挪到"写之前"而不是"读之后"，比 reload 侧更强：reload 那一支即使将来接上钩子，也只覆盖"别人改的文件"。
- **我的裁定：不算偷换，但必须补一句边界，否则"答复时复确认"会被下一腿当成万能豁免。**
  逐字依据：`PLAN.md:1644-1645` 的原句主语是"**热加载**放宽"，它管的是**运行中读到的新值**；
  本腿动的是"**答复后自己写的**新值"，两个世代不重叠。真正还开着的缺口是 §A-1 那条：
  `writeAllowedDirs` 当场把 `m.cur` 换成放宽后的名单（`allowdirs.go:110`），
  **今天**没有第二个运行期读者所以承诺为真；**一旦有第二枚 per-call 读者**（面板 workspace 判定、常驻重建 canonicalizer），
  那行"本次运行仍按旧名单"立刻变成假话，而那时候生效的就是"未经确认的当场放宽"——那才是 D33 那句话的形状。
  ⇒ 要求（具名）：本腿不退回；但**票面或 `A##` 必须写清这条边界**：
  "`[fs]` 的程序化写只保证下一次启动生效；任何运行期读者若改为直读 `Manager.Config()`，必须先接 `ConfirmLocked` 或走 C26 快照重建"。

---

## §B｜对编排者初判的更正（交付格式要求的"另加三节"＝§B/§C/§D；两枚待裁点＝附 A。不同意的摆前面）

起手 `2026-09-29 12:10 +0800`。

1. **不同意——"FAIL 只有 `internal/ball` 1 例与 `internal/panel` 4 例"这句与盘上不符。**
   尺＝`grep -n "^FAIL\s" .scratch/wisp/probes/201/r2/gates-full.txt` ⇒ **三枚**：`:17` ball、`:52` panel、**`:63` `internal/risk` 7.825s**。
   逐字因由：`pathresolver_budget_norace_test.go:37 C26 budget breach: Resolve averages 1.222225ms per call, budget 1ms`。
   另外三枚耗时读数与文件不符（`config` 1.237s 非 1.186s、`perm` 0.671s 非 0.512s、`memory` 14.525s 非 14.655s），
   而 `cmd/wisp 93.362s`／`approval 0.576s`／`tools 17.972s`／`^ok` 枚数 23 与转述一致 ⇒ **转述像是另一场跑数**。
   **后果**：拿"只有别人地界红"去支撑"本批无新增红"，在 `internal/risk` 那枚上今天没有证据。
   我能替本批辩护的只有这一点，且是我现跑的：五枚 commit 的文件清单（`git show --stat`）**一行都没有 `internal/risk`/`internal/tools`**，
   而那枚尺量的是 C26 每调用墙钟、不读程序化写面。
   **建议处置**：空载复跑一次或把 risk 那枚具名挂到"与写腿 `222-r1` 同场跑的负载"，**别让它默认无人归属**；
   台账里那段"FAIL 只有两包"需要按 `A##` 的规矩追加更正，不是抹掉。
2. **不同意 AC#2 的注可以温和。** 编排者判"可翻（带注）"我同意结论，但注的**内容**必须写成逐字这两句：
   `run.go:436-440`"…What it still cannot do is veto an L1 window … `NewChannels()` stays empty and `runSpec.replyVeto` stays unset"＋
   `approval_reply_201_test.go:565-568`"the seam **standing in** for the ball/Esc-hook leg"。
   **尺**：`grep -rn "replyVeto" --include=*.go cmd/ | grep -v _test` ⇒ 只有字段声明 `run.go:146` 与读取 `:601`，**生产零赋值**。
   证据件 `:32` 那一行只写"否决走 `replies.go:419 Veto`→`Gate.Veto`"，**没有这句自陈**——
   读表的人会以为"L1 能被否决"已经能用。这枚是证据件的真实缺陷（措辞过度），我按"退回不成立、注必须补硬"裁。
3. **AC#4 我同意"退回不成立、前置具名转票 223"，但要求两笔登记，缺一笔我这格就不算裁完**：
   ①"当场生效"要落成 `DEFERRED(D-xx)` 或一条 `A##`（票面硬约束 5：每片完成后"推迟项登记更新，新增推迟必须五字段齐全"），
   **不能只留在证据件 `:86` 那一行**；②"可撤销"那一半是**门都没有**，比"没做当场生效"更接近票面 §要建什么 第 3 段（`:29`）的落点要求，
   见 AC#5 格。**依据摆明**：票面 `:44` 与 `:28` 两句里都没有"当场生效"，`PLAN.md:1644-1645` 那句的原话方向是"**不得静默生效**"
   （禁止未经确认就生效），把它读成"必须当场生效"是给冻结文字加它没写的内容；真要要求当场生效，那是 owner 的一句话，不是验收侧的补白。
4. **AC#5 结论（不翻）我同意，依据要换。** 编排者用的是票面**原句**（`approval_grant` 写手为 0）；
   但票面 `:37` 已经把这句改写成三个分句，改写后真正的缺口是
   **"可撤销"零门**（`SetAllowedDirs` 零调用者零测试）＋ session 那半**恒真**（`InsertGrant` 非测试调用者为 0，我现跑）。
   拿原句判"不翻"会在下一腿被反驳（"改写句的长期支已经落了"），**换成改写句的口径这格才站得住**。
5. **AC#1 不翻——完全同意，且我独立跑到同一读数**（`approval_reply.go:111/:113` 是全仓唯一非测试持有者；
   球只有 `cmd/balldebug/main.go:27` 引；托盘 `internal/ball/tray_windows.go:86` 只有"打开面板"一项）。
   这一格的诚实申报我认，但**别把"接缝已就位"写进翻勾栏**：票面要的是枚数，不是形状。
6. **AC#6 带注可翻——同意**，注措辞："今天红的是'名字有人报'，不是'球真的变'；球接入所需的调用已在 `replies.go:283` 就位"。
7. **AC#7 保住——同意，但补一条：那句"响亮"今天零尺。**
   `grep -rn "本机没有可交互控制台" --include=*.go cmd/` ⇒ 只命中 `main.go:150` 自己。
   **把那两行删掉不会有任何用例红** ⇒ 我登记"要补一枚尺"，不构成本批退回（201-r1 就是这形状，本腿没削）。
8. **一处具名更正（小、但会误导下一位）**：证据件 §⑥ 第 2 行写
   "`config.Manager.ConfirmLocked` 仍为 nil（`run.go:341 config.NewManager(cfgPath, nil)`）"——
   **括号里的因果是错的**：`manager.go:78 NewManager(path string, res SecretResolver)` 的第二枚参数是**密钥解析器**，不是这枚钩子。
   钩子为 nil 的真正原因是**全仓没有任何生产代码给这个公开字段赋值**（我现跑 `grep "ConfirmLocked = |ConfirmLocked:" | grep -v _test` ⇒ 空）。
   fail-closed 的结论不变（`manager.go:240`），但这句话不能照抄进 HANDOVER。

---

## §C｜这批产码里我看到的、票面没覆盖的风险

1. **新增导出面的枚数（现读计数，票面从没说要交这么多）**：
   `internal/agent/approval` 新增 **3 型（`ReplyCard:70`/`Replies:101`/`HostBinding:124`）＋1 常量（`MaxTrackedCards:49`）
   ＋3 哨兵错（`:57/:60/:64`）＋1 构造函数（`NewReplies:143`）＋15 方法（`Replies` 上 14 枚：`Attach/Record/Look/Forget/Pending/`
   `AwaitingHuman/WaitingState/Allow/Reject/PanelReject/PanelAllow/Veto/Head/View`，加 `ReplyCard.WideningRule`）
   ＋12 导出字段**＝**35 枚导出标识符；`internal/config` 另加 3 枚导出方法（`AllowedDirs:39`/`AddAllowedDir:56`/`SetAllowedDirs:87`）。
   票面 §要建什么 只写了"把三枚接上真实入口"，没写"新建一型 35 枚标识符的对外面"。
   我的判断：**这 35 枚不是宽松面**（`Replies` 只有转发与读，我逐方法读过；两台路由器仍是唯一门），
   但**它是票面没授权的增长**，应在本票的缺口审计里具名，别在下一票被当成"本来就有"。
2. **命名不触那张禁名词表——但不是因为避开了，是因为那把尺看不见。**
   现读 `internal/panel/l2_grant_boundary_test.go:192-195`（`grantRouteWords` 含 `"allow"`、`"grant"`、`"approve"`…）
   与 `:1882-1886`（入向判据候选名含 `"allow","Allow","grant","Grant","verdict"…`），
   再核它的走查根 `:1231/:1574/:1848/:2023` **四处全是 `filepath.Join(root, "internal", "panel")`** ⇒
   `Replies.PanelAllow`、`Replies.Allow`、`ReplyCard.Grant` 全部在射程外，**今天不可能红**。
   那枚文件我**一字未动**（冻结件）。⚠ 由此带出一条真风险：
   `ReplyCard` 带着**活的 nonce**（`replies.go:78 Grant`），而 `Replies.Pending()`（`:224`）把整张卡交给被 handed 的同一枚类型；
   票 33 的 WebView 宿主如果拿的是同一枚 `*Replies`，它读一次 `Pending()` 就把 grant 抄进了页面可达的对象图。
   路由侧仍拒（`gate.go:638-649`），泄露的 nonce 也会被 `revokeGrants` 烧掉——**但那是"事后烧"，不是"不给你看"**。
   建议（不写码）：接页面宿主那天，交给它的面应当只有 `PanelReject/Head/View` 三枚，而不是整枚 `*Replies`。
3. **"一直"的粒度从文件放大到目录。** `WideningRule`（`replies.go:467-483`）存的是被写文件的**父目录**；
   D45 那张 `approval_grant` 本来有 `pattern` 列可以写更窄（`internal/memory/schema.go:89-99`），
   改写后的 AC#5 把落点搬进 `allowed_dirs` 之后，粒度只剩目录级。
   卡上印的是目录原文 ⇒ AC#4 因此成立；但用户点的是"这一个文件"那一发，**存下来的是整个目录**。
   这条票面没写，属**必须登记的风险**，不是必须退回的缺陷。
4. **第二张卡的孤儿形状。** `always` 把重新确认挂在监听器续体上（`approval_always.go:97-102`）；
   操作者答完原卡、任务收敛、进程退出时，第二张卡会随 ctx 结束落到 `WIDEN-REFUSED`（`:128`，方向＝什么都不写，fail-closed 正确）。
   但"我一直"的直觉会落空。同步回复里写了"这会再起一张 L2 卡"（`:103-105`），所以是 UX 风险、不是安全洞。
5. **D38b 常驻名册复用 `approval-waiter`**（`:99`；证据件 §⑥ 第 4 行自陈）。我复核这枚自陈属实且方向保守
   （新造 roster 名是契约面），但**widening worker 与 waiter 在 `RosterReport` 里会并成一名**，
   真发生泄漏时读不出是哪一支。归 owner 动 D38b，本腿不该自己做——我不判退回。
6. **`statOwnWrite` 在这枚二进制里今天唯一的作用是"让未来的轮询看不见这次覆盖"**（与 §A-1 的整份覆盖串成一条链）。
   实现者在 `36b46125` 已自陈"认领今天用不上、不许下一腿当成 D36 的实现"，这句我核过与盘上一致
   （`CheckAndReload` 唯一非测试调用者＝`cmd/balldebug/main.go:243`），**但它没有把"覆盖＋认领＝看不见手改损失"这一层写出来**。
7. **卡片名没有触禁区**：`allowWidenTool = "config.allow_dir"`（`:65`）与 `allowWidenTaskID = "host:fs-allow-widen"`（`:64`）
   只是卡片与审计上的名字，不进 C17/入向名册（我核 `cmd/wisp/panel_inbound.go` 的入向链只有 composer/`panel.mode.request`），
   且与 `run.go:605-607 modeSwitchToolName = "permission.mode"` 同先例 ⇒ 票面 `:51`"新增方法名要先落 A##"**没被触发**。
   这条我列出来是为了让下一位不必再查一遍。

---

## §D｜我没核动的（诚实列，不含糊）

1. **没跑任何测试／编译／`go vet`／`d22scan`**（硬约束 1）⇒ 实现者自述的"8 枚新用例全绿""gofumpt 净""d22scan clean"
   我一枚都没有独立读数，只能引编排者的 `gates-full.txt`，而那枚读数**我已经发现与转述不符**（§B 第 1 条）。
2. **`internal/risk` 那枚 budget 红的归因没核**：它是否只在有并发写腿时红，我没有空载读数；
   我只核到"本批不碰那两包"。**这不是我判它无关的证据，是登记。**
3. **`internal/ball` 1 例与 `internal/panel` 4 例"属别人地界"我按编排者的逐字归因接受**——
   那句 `read design/assets/tokens.css … cannot find the path specified` 是我在 **probes 文件**里读到的，
   我没有（也不会）去读 `design/**`；删除方是谁我没查。
4. **`dirOf`（`batch.go:91`）对 UNC / `\\?\` / 8.3 短名的实际输出我没读**；
   因此 `replies.go:460-466`"Prompt.Paths 已是 C26 输出、`dirOf` 是纯字符串切分"这句我按自陈接受，未反证。
5. **`SaveFile` 之下我没读**：`MarshalCanonical`／`atomicWrite` 对 schema 之外的未知键与注释是否保留，
   我只读了 `loader.go:135-145` 那 11 行。**若未知键会丢，§A-1 那条覆盖的损失比我又写的更大。**
6. **`queue.expire` 到点判拒绝那一条我没现读实现**：只在 `git show 58bf0158` 的 diff 里确认 `queue.go` 本腿只改注释，
   语义依赖票面 `:12` 与 201-r1 的读数。
7. **并发时序没推演**（`-race` 要跑测试，属禁区）：`Replies` 的锁我读过（`routes()` 持锁读、`Record/Forget` 写），
   实现者 §⑦ 第 5 条踩到的"钩子一探到账本就答"我把落点（`run.go:1103`→`1110`）核成实，
   但"同一瞬间两枚卡被两个面分别答复"的时序我没有独立证据。
8. **票面 §要建什么 第 1 段原句"先接球与托盘这两枚"（`:25`）今天仍没做**：我按 AC 表逐枚裁，
   没去裁"派单四段与 AC 表冲突时以谁为准"——**那是编排者的口径，不是我的一票**。

---

## 结论（三行）

1. **AC#4 我判：成立，不退回。** 逐字依据＝票面 `:44`/`:28` 两句都只要求"规则文本＋二次确认"，
   `PLAN.md:1644-1645` 的原句方向是"不得静默生效"（＝不许未确认就生效），而不是"必须当场生效"；
   实现两处都满足且**有可判红的正控**（`approval_always_201_test.go:208-209/:211-215/:218/:181-189`）。
   前置具名＝**票 223 那一族（`CheckAndReload`/热加载落点）＋ D36 生效级别定案**，且必须落成 `DEFERRED(D-xx)`/`A##`，不能只留在证据件 §⑥。
2. **可翻＝AC#3、AC#4、AC#6（带注）、AC#7（带注：响亮那半今天零尺）；不翻＝AC#1（入口枚数仍 1）、AC#5（"可撤销"零门＋session 支恒真）；
   AC#2 带注可翻但注必须写硬（`replyVeto` 生产零赋值 ⇒ 否决今天不可达）。** 没有一格需要"退回这批产码"。
3. **两处必须编排者亲自处置、我不能替**：
   ①门禁转述失真——`gates-full.txt` 里**红的是三包不是两包**（`internal/risk` 的 C26 预算红无人归因），
   且三枚耗时读数与文件不符 ⇒ "本批无新增红"这句在 risk 上还没被证成；
   ②`SetAllowedDirs` 是**导出、零调用者、零尺**的锁定 section 放宽门（附 A-2），补尺或收面，别让它跟着"名字叫 SET"进下一票。

**交付判据**：本文件的 `wc -c` 由编排者收到后现量（我不自证字节数）；锚点 `7096470d`
（产码尾 `36b46125`，两者之间这批产码文件零改动，`git diff --name-only` 现跑＝空）。



