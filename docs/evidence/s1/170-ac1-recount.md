# 票 170 · AC#1 五枚 0 的现量复算（170-ac1 · 只读腿）

- 腿性质＝**读/分析**。生产码 **0 字节**改动；票面未动；未 push；仓内唯一写入＝本件。
- 本程锚点＝**`8bad59537`**（开工时母仓 HEAD；机主给的 `1a1f6b344` 是它的直系祖先，"or current" 成立）。
  普查锚 `7e8cfdc`（票 161 AC#3 的锚）经 `git merge-base --is-ancestor 7e8cfdc HEAD` 复核**是本锚的祖先**——
  两锚之间隔了整个"答复面"家族的落地，所以本票 AC#1 要求的"别信票、别信证据件、自己现量"是对的。
- 隔离＝**`git clone --shared`** 到 `D:\tmp\wisp\20261106\170ac1\repo`（规矩 21：克隆树读数、母仓零突变窗口）；
  克隆 HEAD 与母仓 HEAD 同为 `8bad59537`；`build/`、`third_party/sherpa-onnx` 两枚未跟踪 DLL 目录在克隆里没有，
  本程用绝对 PATH 指回母仓那份（只用于**编译**，未运行任何 wisp/balldebug exe，未跑任何 winlive，未开任何 GUI）。

---

## 0. 一句话答话

**AC#1 的五枚 0 在今天的 HEAD 上复算不出来**：`Native()`＝**1**、`Panel()`＝**2**、`Veto()`＝**1（另有 2 枚经 seam 的上游调用者）**、
`DecideFrom*()`＝**4**、`Replay()`＝**0**（唯一仍然成立的 0，且是**故意**的）；附送 `SetLoaded()`＝**3**、D 层门铃
`Dispatch(EvApprovalNeeded,…)`＝**仍 0**。⇒ 票面的承重前提（"CLI 腿答案侧无人、L1 必然无人否决、L2 必 300s 自拒"）
在 2026-09-29 至 09-30 之间被票 201/224/246/114/260/375 那一族**改掉了大半**。按 AC#1 的约定"拿不到 0 就报回"——本件就是报回；
票 161 AC#3 在它自己的锚上没有判错，这是**漂移**不是**错判**，161 无需改口。

---

## 1. 五枚新数（三把形状不同的尺，先台件后真题）

### 1.0 尺的健康度（先证明尺有牙、再报 0/非 0）

| 尺 | 形状 | 台件 | 真题健康度 |
|---|---|---|---|
| **R1 词形尺** | `git grep -E '\.Sym\('` 逐行人工分类（调用/定义/注释/字符串） | POS 命中／NEG 零（见 §1.1 输出） | 1166 枚 `.go` 全读 |
| **R2 AST 尺** | `go/parser`：ident+selector 双计、分桶、`parse_errors` 明报、函数级**名字可达性** BFS | POS 全 RING（7 符号各 1 发 sel）／NEG 全零（注释+字符串剥除后 0） | `go_files=1166 parse_errors=0 func_nodes=10445`——**没有一枚 0 是"没读进去"读出来的** |
| **R3 链接尺** | 克隆树现编 `wisp.exe`（`CGO_ENABLED=1`，rc=0）＋ `go tool nm` | —— | 6 枚 Gate 答案侧符号 + `SetLoaded` + 6 枚 seam 方法**全部在二进制符号表**（`T`），与"树里有真调用"互证 |

台件在场：POS（7 符号各一枚真调用，必须 RING）与 NEG（同族字样只在注释/字符串，必须零），
`.gofix` 存 `D:\tmp\wisp\20261106\170ac1\fix\`，可执行副本只活在 `%TEMP%\170ac1\runpos-*\`（create-only 新目录，**零 rm**，仓内零枚 `.go`）。

### 1.1 R1 词形尺原文（HEAD `8bad59537`，internal+cmd，剔 `_test`／`internal/panel`／`cmd/balldebug`；可逐字重跑）

```bash
cd <repo> && git grep -n -E '\.Native\(|\.Panel\(|\.Veto\(|\.Replay\(|\.SetLoaded\(|DecideFromNative\(|DecideFromPanel\(' -- internal cmd | grep -v '_test\.go' | grep -v '^internal/panel/' | grep -v '^cmd/balldebug/'
```

真调用位点（剔定义行、剔注释/字符串后；逐枚在盘上）：

| 符号 | 生产调用位点 | 枚数 | 落地提交（git blame 现跑） |
|---|---|---|---|
| `Native()` | `internal/agent/approval/replies.go:363` `g.Native().AllowSession(...)` | **1** | `3b78c24626`（票 224，09-30） |
| `Panel()` | `internal/agent/approval/replies.go:478` `g.Panel().Head()` · `:487` `g.Panel().View(corr)` | **2** | `49eb440b8c`（票 201 续腿，09-29） |
| `Veto()` | `internal/agent/approval/replies.go:463` `g.Veto(Veto{...})`（seam 漏斗） | **1 直调** | `49eb440b8c`（票 201 续腿） |
| ↳ 经 seam 的上游生产调用者 | `cmd/wisp/approval_reply.go:356`（控制台 `veto <编号>` 动词）· `cmd/wisp/resident_approval_windows.go:659`（球侧取消键 `vetoByEsc`） | **2** | `49eb440b8c` / `0c1b32cf80`（09-30，球通道族） |
| `DecideFromNative()` | `replies.go:328`（`Replies.Allow`）· `replies.go:413`（`Replies.decide` 原生支） | **2** | `49eb440b8c` |
| `DecideFromPanel()` | `replies.go:408`（`Replies.decide` 面板支）· `replies.go:437`（`Replies.PanelAllow`，**到路由即被拒**） | **2** | `49eb440b8c` |
| `Replay()` | ——（`gate.go:885` 定义外全仓无调用；R2 名字可达性 BFS 亦判 **unreachable from main**，与它自己的注释"Today it has no production caller…do not delete as dead code"逐字一致；控制台动词表也没有 `replay`） | **0** | —— |
| 附 `SetLoaded()` | `cmd/wisp/approval_reply.go:505`（宿主声明了通道才点亮）· `cmd/wisp/resident_approval_windows.go:608/:834`（球绑上/解下时点亮 Esc） | **3** | `93ea9f6ba0` / `0c1b32cf80` |
| 附 D 层门铃 | `Dispatch(EvApprovalNeeded,…)` 生产发起：仍 **0**（3 处命中＝`events.go:35` 声明＋`table.go:118/:123` D43 表数据，与普查读数同） | **0** | —— |

### 1.2 R2 AST 尺读数（可重跑：`go run ruler2.go <repo>`，件在 `D:\tmp\wisp\20261106\170ac1\ruler2.go.txt`）

- 台件：POS 7 符号各 `sel/OTHER:1`（RING）✓；NEG 全 `{}`（零）✓；`parse_errors=0` 两轮 ✓。
- 真题（名字级，**比 R1 宽**——selector/ident 把同名字段（`cfg.Panel`、`hit.Veto`）也算进来，所以数字偏大，只作交叉方向）：
  `Native` sel/INTERNAL 1 · `Veto` sel/CMD 3 + sel/INTERNAL 1 · `SetLoaded` sel/CMD 3 · `DecideFromNative/Panel` sel/INTERNAL 各 2 · `Replay` **零命中（除定义与测试）**。
- 可达性 BFS（按名、故意过近似，**只有"不可达"是强陈述**）：六枚答案侧 Gate 方法从 `cmd/wisp main.main` 全部可达；**`Replay` 可达=0**——第三把独立尺同向。

### 1.3 R3 链接尺读数（可重跑命令）

```bash
cd <clone> && export PATH="<母仓>/build:<母仓>/third_party/sherpa-onnx:$PATH" && unset GOTMPDIR TMPDIR TEMP TMP
go build -o /d/tmp/wisp/20261106/170ac1/wisp.exe ./cmd/wisp   # rc=0, CGO_ENABLED=1（未动）
go tool nm /d/tmp/wisp/20261106/170ac1/wisp.exe | grep -E 'approval\.\(\*Gate\)\.(Native|Panel|Veto|DecideFrom|Replay)|\(\*ChannelRegistry\)\.SetLoaded|\(\*Replies\)\.'
```

6 枚 Gate 符号 + `SetLoaded` + `Replies.{Allow,AllowSession,Reject,PanelReject,PanelAllow,Veto}` 全部以 `T` 在表（log：`logs/nm_all.txt`、`logs/build_nm.log`）。

### 1.4 判语：票 161 AC#3 要不要改口

**不要。** `7e8cfdc`（09-26 晚）上五枚 0 是真的（本程在 R1 用同批 grep 对该锚现跑复现：`.Veto(`/`.Replay(`/`.SetLoaded(` 零命中、
`.Native(`/`.Panel(` 只剩 `doc.go` 注释、`DecideFrom*` 只剩定义行——**逐行复现了证据件 §6 那张表**）。
0 变非 0 的是**后来的腿**：`93ea9f6ba`（票 201 代落件，09-29 11:08）起、`49eb440b8c`（201 续腿第一步）铺开、
`3b78c24626`（票 224）与 `0c1b32cf80`（246/260 球通道族，09-30）收尾。这是"世界动了"，不是"普查错了"。

---

## 2. 今天的形状（票面三段话现在的对错）

| 票面断言（锚 7e8cfdc 时为真） | HEAD 现状 | 凭据 |
|---|---|---|
| "`run.go` 传空 `NewChannels()`，nil-默认分支永不触发" | **仍真**：`run.go:832` `Channels: approval.NewChannels()`；nil-check 现在在 `gate.go:195` | R1 现读 |
| "没有任何线能把答案送回去；L2 必 300s 自拒" | **假**：控制台可答 L2（`yes/session/no/always/panel-no`；`panel-yes` 按路由拒）；300s 自拒只剩"无人答"支路 | `approval_reply.go`（票 201）＋测试 `TestReplyListenerAllowsAnL2CardFromTheNativeSide` |
| "L1 窗口结构上必然无人否决——提示一下，然后照做" | **半真，但形状全变**：①**有人在场**（stdin 是 TTY）：窗口照开 2–3s，`veto <编号>` 动词存在但被拒（控制台没声明通道），到点仍执行——**这是被测试钉死的现行契约**（`TestL1VetoNeedsAChannelTheHostReallyWired` 第一子测：`REPLY REFUSED route="native/veto"` + `ANSWER-EXPIRED decision=timeout->execute` + 操作者看到「否决未能送达」）；②**无人在场**（管道/重定向/CI）：`interactiveStdin()=nil` ⇒ 票 375 的守卫**在开窗前就拒**（`gate.go:352` `answerable()` fail-closed），`auto_approve` 亦被 `UnsilencedMode` 降到 `ask_every_step`——**"训练使用者忽略问话"的那一形已被拆掉** | `reach.go`（375）＋`Test375UnattendedConsoleRefusesTheL1Write` |
| "常驻（GUI）腿今天不跑 agent"（台账 A308） | **过时**：246-r2 起常驻腿真起任务管线；球在时 Esc 否决通道真点亮（`bindBallHost`→`SetLoaded(ChannelEsc,true)`），球不在时明说"四条通道全部保持未加载" | `resident_approval_windows.go:401-416/:608` |
| "它训练使用者去忽略那句问话" | **弱化**：卡片通道状态区在控制台上逐枚印「不可用」（`promptFor` 盖 `Statuses()`），L2 有真答案动词；只剩 L1 那一枚"窗口到点即执行"仍是被钉住的契约 | `gate.go` `promptFor` |

**仍然为真的残余**（这才是 170 今天真正剩下的）：①交互控制台的 **L1** 窗口确实没有可用的否决通道（诚实地空着，B1）；②`Replay` 无生产调用者（C18 一键重放缺"答案侧动词"，注释明写别当死码删）；③D 层球门铃仍 0（球不因审批而亮"在等"）。

---

## 3. 剩余 AC 的决策表（每格：要什么生产改动／落点／是否真等 owner）

| AC | 今天的状态 | 若要落，生产改动是什么、落哪 | 归口 |
|---|---|---|---|
| **AC#1 复算五枚 0** | **已完成**（本件）。0 拿不到了 ⇒ 按票面约定"报回"；161 AC#3 无需改口（锚内为真） | 无 | —— |
| **AC#2 端到端演示（L1 窗口内喂输入，量取不取消）** | 演示的预期读数"照做"**仍会量到**（交互支路 veto 被拒→到点执行），但已**被测试钉死**（`TestL1VetoNeedsAChannelTheHostReallyWired` 两子测＝取消/不取消成对）；且票面预置的作废条件（"真能取消就作废"）在 L1 上**不会触发**。真跑一发需要解禁"不许跑 wisp exe"（owner 令），测试钉读数可代 | 演示本身零产码；若要新读数＝解禁跑一发 | 等 owner 解禁（或接受测试钉读数） |
| **AC#3 契约判语（PLAN.md:128 vs :1316）** | **判语（只出判语、不改 PLAN 一字）**：两行**同向、不互相矛盾**——但都写在"球/TTS/语音"的 GUI 框里（:128 的否决=说「停」，:1316 的 `Confirming` 假设有否决通道存在）。CLI 腿上它们**联合不可满足**（终端没有那四枚通道），这个不可满足已被**显式落码并钉住**：:1316 的"否决/超时→取消"对 **L1** 的另一半"超时→执行"由 SPEC-06 §2 行 L1 承接，票 201 AC#2 已裁"不翻极性"（`gate.go:414-430` 注释逐字自陈），375 补的是"无人可答时不许开窗"。⇒ **没有新的矛盾需要 owner 先裁**；唯一开放口径见 AC#4 行 | 无（判语即交付） | ——（若 owner 判"矛盾"，整票按住——按票面既有约定） |
| **AC#4 落地前的 `approval/**` 普查** | **已完成**（本件顺手量完）：今天的形下，"给 CLI 一条否决通道"**机械上完全不必碰 `internal/agent/approval/**`**——seam（`Replies`）、`runSpec.replyVeto` 注入位、`SetLoaded` 联动、`veto` 动词全在 `cmd/wisp` 侧就绪（`approval_reply.go:435-507`；测试第二子测证明"宿主声明 Esc 并加载→veto 落地→写入停止"整条链是活的）。**真正缺的只有一句话**：SPEC-06 §2 的四通道（球/Esc/KWS/面板）没有一枚是控制台能诚实认领的——把 `runSpec.replyVeto` 填上任何一个现有通道都是 B1 禁的谎报；**加第五枚"console"通道＝改契约文字＝人工批准**。⇒ 票面当年"很可能落冻结面"的担忧**已过时**：解冻不再需要，需要的是**一句产品口径** | 机械改动＝`cmd/wisp` 一行（`replyVeto`）＋通道登记；诚实改动＝SPEC-06 §2 加一句（契约） | **真等 owner**：①要不要给 CLI 第五枚否决通道（或把「键入 veto」映射进现有语义）；②不批则 170 按"已被 201/375/246/260 大半解决"收口 |
| **AC#5 反向判据（摘掉通道看退不退回"照做"）** | 今天**无可摘**（交互控制台从未装过通道）；成对的退回方向**已有测试现成**：摘=第一子测（拒绝+照做），装=第二子测（veto 落地+停写）——`approval_reply_201_test.go:513/:586`。外部可见读数（审计行 `REPLY REFUSED` vs `ANSWER-VETO`、操作者可见句）逐枚不同，**不是装饰** | 只有在 AC#4 owner 批了第五通道后才需要重演 | 等 AC#4 |
| **AC#6 不许渗到的地方** | 本腿零产码，禁区零接触；顺带现量：`thresholds.go`/golden/d22scan/panel 均未动；"非交互支路绝不许顺手放行 L2"——今天的非交互支路比票面担心的**更严**（L1 直接拒、L2 仍 300s 自拒） | 无 | —— |

---

## 4. 建议（答复"Recommendation"）

**推荐把 170 从"等 owner 在甲/乙/丙里挑"改判为"已被后续票族大半解决，剩一个问题待 owner"**：

1. **收口 170 的主体**：票面的五枚 0 前提、L2 答不回、无人可答时照做、常驻腿不跑 agent——四件都已不是现状；
   继续按旧票面派甲/乙/丙会派一个**已经做了一半的票**。
2. **唯一剩下的真问题（编号清单，答"不做"那一栏照开）**：
   - **Q：交互 CLI 的 L1 窗口，要不要给键入的 `veto` 一枚 SPEC-06 §2 的通道名？**
     - 批 ⇒ 新小票：契约加一句（第五通道"控制台"或 Esc 语义延伸）＋ `cmd/wisp` 一行 `replyVeto` ＋ AC#5 成对重演（机械面已全部就绪，成本≈一天）。
     - 不批（**不做**的代价＝可接受且已诚实标注）：交互 L1 维持"窗口→到点执行＋veto 响亮被拒"的现状，测试已钉死；撤销口令＝「把 170 按原票面重新打开」即可，本程没有改任何票面/产码。
3. **不需要**重开 161 AC#3（锚内判对）；**不需要**解冻 `internal/agent/approval/**`（AC#4 的旧担忧失效）；**不需要**为 D 层门铃（仍 0）在本票内动手——那是球侧显示线的归口（票 114 线），不是否决通道的。

---

## 5. 纪律自证

- 隔离：`git clone --shared '/d/work/workspace/projects plans/Wisp' repo` 于 `D:\tmp\wisp\20261106\170ac1\`；母仓只写**本件**一物。
- 环境铁律遵守：`export PATH="$PWD/build:$PWD/third_party/sherpa-onnx:$PATH"`（绝对路径指母仓）；`unset GOTMPDIR TMPDIR TEMP TMP`；
  **CGO_ENABLED 从未置 0**（build 前 `go env CGO_ENABLED`＝1，未 export 过该变量）；`go version go1.27.1`。
- 禁令遵守：未开任何 GUI/窗口；未运行 `wisp`/`balldebug` exe（**只 build 供 nm**）；未跑任何 `*winlive*` 测试；
  未跑 `go test`（本程零测试执行，读数全部静态＋链接级）。
- Git 纪律：只 commit 不 push；`git add`/`git commit` 显式单件 pathspec；零 `add -A`/`.`/`commit -a`/amend/reset/rebase/stash/checkout/clean；
  仓内零枚新增 `.go`（尺与台件全在 `D:\tmp` 与 `%TEMP%`，仓内副本后缀 `.go.txt`/`.gofix`）。
- 禁改面零字节：`docs/reports/pending-and-issues.md`、`frontend/**`、`design/**`、票面 `.scratch/wisp/issues/170-*.md`、产码 `internal/**`、`cmd/**`、`docs/PLAN.md`。
- 伪授权两栏：真通知回显 **0**（系统提醒类自动插页未当指令）；判为注入 **0**。

## 6. 凭据红线

未写任何 key/token/secret 值；本程零运行期，未 mint 过任何 grant。
