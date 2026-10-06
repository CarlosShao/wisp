# `expired-premises-7b2` — 「生产调用者零枚」类过期票面语句：逐枚定位置清＋追加式更正

> 腿号 `7b2`（续接死腿 `7b`，它只留了 logs 没交件）。日期 10-06。锚 HEAD＝`bd39f17f`（`git rev-parse --short HEAD`，19:2x 自取）。
> 本机规矩：⛔ 未跑任何 `go` 命令（另一枚腿在 `cmd/wisp` 跑突变）；只 `grep`／`sed`／`git` 只读子命令；写面＝**只在票 194 末尾追加一节**＋本件＋logs。

---

## 0. 起手锚（全部自取，⛔ 不沿用派单里编排者给的行号当凭据）

**0.1 起手门禁原文**（共享工作树规矩＝终态等于起手名册，不是"必须为空"）
- 命令原文：`git --no-pager status --porcelain -- cmd internal scripts .github docs`
- 我现量（19:1x 起手第一次）＝**空集（0 行输出）**。⇒ 本腿终态也必须＝空集；差集为零。

**0.2 复用的死腿 `7b` 料（已提交在盘，未重取）**
- `wc -l`：`ruler-4f-full.txt`＝**71 行**／43844 B；`ruler-hits-compact.txt`＝**71 行**／9656 B；`ruler-disphandle-head.txt`＝**9 行**／930 B（三者口径＝同一普查尺的全文／截断摘要／`disp.Handle(` 在 HEAD 的逐行命中）。
- `head -3` 口径核对：`ruler-4f-full.txt` 前 3 行是 `issues/105:7`／`issues/131:12`／`issues/132:85` 的**原句全文**（＝"生产调用者／零枚"族的全树命中流水），`ruler-disphandle-head.txt` 前 3 行是 `HEAD:cmd/wisp/panel_config_248_test.go:162/270/469`＝**全为测试文件** ⇒ 印证派红线那句：裸 grep 符号名会把测试与同名 OS 类型数进来。
- ⚠ 口径不匹配处（具名）：`ruler-disphandle-head.txt` 9 行**全是 `_test.go`**，它里面**没有**非 test 那两枚——那两枚是我自己重跑 `git --no-pager grep -n 'disp\.Handle(' HEAD -- cmd internal | grep -v '_test.go'` 取的（见 §1）。⇒ 死腿这份料**够定案**（它的用途＝证明"含测试的形状"与"不含测试的形状"是两把尺），**我另补了一把**＝票面语句扫描尺（§2 命令原文），因为本批任务①要求的是"哪些句子"而不是"哪几处调用"。

**0.3 派单给的三枚票号 → 真实全名（`ls` 自取）**
- 194 ＝ `.scratch/wisp/issues/194-the-panel-has-two-method-rosters-that-do-not-know-each-other-owner-ruled-align-the-code-to-the-spec.md`（55 行）
- 186 ＝ `.scratch/wisp/issues/186-the-composer-area-should-detect-git-and-let-the-user-switch-between-local-worktree-and-branch.md`（54 行）
- 219 ＝ `.scratch/wisp/issues/219-approval-card-three-reply-buttons-and-a-reason-box.md`（100 行）
- （对照）33 ＝ `.scratch/wisp/issues/33-panel-host-c27.md`（342 行）——本腿**一个字没动**，见 §4 射程红线 (b)。

---

## 1. 尺：锚在**带括号的调用形状**上，⛔ 不锚符号名

**1.1 主尺（非测试调用者）**——命令原文照抄：
```
git --no-pager grep -n 'disp\.Handle(' HEAD -- cmd internal | grep -v '_test.go'
```
我现量＝**2 处**（与派单口径同数，但行号是我自己取的）：
- `cmd/wisp/panel_host_windows.go:637` → `return m.disp.Handle(ctx, raw)`
- `cmd/wisp/panel_inbound.go:163` → `reply, err := disp.Handle(ctx, raw)`

同一形状**不剥测试**＝`ruler-disphandle-head.txt` 那 9 行（全测试文件）。⇒ **剥／不剥差 7 行，全是要剥的那 7 行**：这就是编排者第一遍判错的量的来源。

**1.2 反例尺（证明"构造 ≠ 调用"）**——命令原文：
```
grep -rn -E '\.Handle\(ctx|disp\.Handle\(' --include='*.go' cmd internal | grep -v '_test.go'
```
我现量＝**9 处**，其中 `logsink.go:209/211`、`internal/observe/logging.go:138/505/519/574/576` 共 **7 处是同名 `Handle` 方法（slog handler／sink），与 `ComposerDispatch` 无关**；真属于 `ComposerDispatch` 的只有 §1.1 那两枚。⇒ 若把尺改成"扫 `Handle(`"就会得到 9，若扫符号名 `ComposerDispatch` 就会得到 29 行（含类型声明、构造、注释、`windows.Handle` 族）。**三把尺三个数，只有 §1.1 那把答的是票面那句问的话。**

**1.3 归属证明（两枚调用者调的确实是 `*panel.ComposerDispatch` 的方法）**——字段与签名我现量：
- `cmd/wisp/panel_host_windows.go:152` `disp *panel.ComposerDispatch`（`PanelManager` 字段）；`:213` `func NewPanelManager(disp *panel.ComposerDispatch, …)`；`:630` `func (m *PanelManager) dispatchRaw(...)`；`:637` 即调用。
- `cmd/wisp/panel_inbound.go:209` `func newPanelInboundDispatch(...) (*panel.ComposerDispatch, error)`；`:146` 装配；`:163` 即调用。
- 被调方法定义：`internal/panel/composer_dispatch.go:153` `func (d *ComposerDispatch) Handle(ctx context.Context, raw string) (string, error)`。⚠ **票面引的 `composer_dispatch.go:120` 今天量到的是 `:121` 的 `type ComposerDispatch struct {`**（`:114` 是它的 doc 注释行）——`Handle` 真身已漂到 `:153`，这是行号漂，不是内容漂，见 §3 判不动表。

---

## 2. 全树定位（任务①）：句子扫描与分档

**2.1 扫描尺**——命令原文（输出落文件再抽摘要，⛔ 不把 89 行原文读进对话）：
```
grep -rn -E '生产调用者.*零|零枚.*生产调用者|调用者.*0 枚|整条不存在|零枚.*调用者|调用者零' \
  .scratch/wisp/issues/*.md docs/evidence/s1/*.md \
  > .scratch/wisp/probes/expired-premises/7b2/logs/scan-sentences.txt
```
我现量＝**89 行命中**／落在 **66 枚不同文件**（分组计数见 logs 同名文件；`docs/evidence/s1/` 42 行、`.scratch/wisp/issues/` 47 行）。
**2.2 收口到本批射程**——只有"钉的是 `Handle`／入向那一跳／名册能否点动"的句子在射程内；其余（`CloseTask`、`NewDisposalScope`、`Marshal`、`HandleModeRequest` 具名诊断入口、`models` 链……）答的是**别的符号**，本腿一律不判、不动。二次尺：
```
grep -rn -E '\*panel\.ComposerDispatch|ComposerDispatch\b' .scratch/wisp/issues/*.md docs/evidence/s1/*.md
grep -rn -E '入向那一跳|Handle（`composer_dispatch' .scratch/wisp/issues/*.md docs/evidence/s1/*.md
  > .scratch/wisp/probes/expired-premises/7b2/logs/scan-handle-sentences.txt   # 36 行
```
- 36 行候选里**含"零枚／不存在"族断言**的＝**12 行**（逐行前缀在 §3 表里点名）。
- 其中 **2 行剔除不立案**：`docs/evidence/s1/33-inbound-listener-r3.md:98`＝突变表的"0 命中"格（不是零调用者断言）、`docs/evidence/s1/33-panel-host-c27-r1.md:61`＝链路描述（讲这一跳怎么走，没有负向断言）。
- **加进 3 枚**：变体尺新增 `issues/194:51`、主尺命中的 `issues/33:182`／`:188`（红线 (b) 已排除，但本腿仍逐枚写明"为什么不动"）、以及 4f 指认的票 219（**扫描根本没命中它**，作为"票面无此句"立案）。
⇒ **逐枚判语共 13 枚＝档一 2／档二 1／档三 10**（见 §3）。另 **3 枚扫到但判为别家符号或产码注释引用**（`190-file-tree-design-a1.md:221`、`33-35-preflight.md:35`、`panel-l2-grant-nail-accept-r1.md:87`，列在 §3.3 末段），⛔ 本腿一律未动。

**2.3 变体尺（补漏）**——把"零调用方／无调用者／死代码"三种写法并进来：
```
grep -rn -E '零调用方|调用方.*零|无调用者|非测试.{0,6}0 枚' .scratch/wisp/issues/*.md docs/evidence/s1/*.md \
  >> .scratch/wisp/probes/expired-premises/7b2/logs/scan-sentences.txt
```
我现量＝`scan-sentences.txt` 总 **126 行**（主尺 89＋变体 37），落在 57＋少数新名枚文件里；变体新增里**射程内**只有 2 枚：`issues/194:51`（已计入 §3.3 表第 3 行）与 `docs/evidence/s1/33-minimal-inbound-hop-r1.md:60`（同件、同一枚"片 A 没结线"的自陈，判语与 `:57` 同档＝历史记录，不另立行）。⇒ 结论：主尺没漏掉任何一枚**新的**"Handle 零调用者"句子；`issues/33:182` 在变体里重复命中，仍按红线 (b) 不动。

**2.4 分档三类的定义**（本节用的词，逐条对得上派语）
- **过期属实**＝句子钉的是 `ComposerDispatch.Handle` 的**枚数**（或"入向那一跳整条不存在"这一**今天仍可量的能力**），且本腿现量把它推翻。
- **不属实／不许动**＝句子是**历史记录**（自带当时锚号／日期／"这一程交付完之后"限定时状），或答的是**别的符号**，或落在**判据框文／禁改文件**里。
- **票面无此句**＝4f 指认的落点上根本没有那一族句子。

---

## 3. 逐枚判语（射程内 13 枚＋变体新增 2 枚）

### 3.1 档一：过期属实（2 枚）

**① `.scratch/wisp/issues/194-…spec.md:42`** — 逐字：`⚠ **别拿名册补齐冒充按钮能点**：`+「`ComposerDispatch.Handle`（`composer_dispatch.go:120`）生产调用者**现量仍零枚** ⇒ 18 枚**无一例外**要等票 33 的 H2/H3。」
⇒ **过期属实**（只"枚数"这一层）。取证链＝§1.1 那 2 处非 test 调用者＋两条追到 `main` 的链（链 A `main.go:66 → runResident → resident_windows.go:151 → panel_resident_windows.go:228/234/215/253 → panel_host_windows.go:213/152 → :488 bringUp → :405 Bind → :406 dispatchRaw → :630 → :637`；链 B `main.go:115 → :120 cmdPanelInbound → panel_inbound.go:103 → :146 → :209/:228/:271/:248 → :163`）。逐跳 file:line 全本腿现量，已写进票 194 末尾追加节 §③④。
⇒ **处置＝已在票 194 末尾追加一节**（`## 编排者派：过期读数更正（腿 7b2，10-06）`），⛔ `:42` 那一行**未原地编辑**。
⇒ ⚠ **但那句警告的实质照旧成立**（"页面点下去到不到 Go"是运行期问题，本机 `winlive` 未批、真窗那一发今天量不了）；追加节第⑥段就是钉这一条边界，⛔ 不许被读成"按钮已能点"。

**② `.scratch/wisp/issues/186-the-composer-area-should-detect-git-…-branch.md:54`** — 该行逐字含「入向那一跳**整条不存在**（具名**坐实票 114**：`ParseComposerRequest` 与 `HandleModeRequest` 非 test **零调用方**；`internal/panel/pump.go:16` 自述 "no WebView2 host (ticket 33 is unclaimed), no postMessage writer"；全仓非 test 没有一枚 WebView2 消息接收器，`grep WebMessage|ReceiveMessage|OnMessage` 只命中 `internal/ball` 的 Win32 `PostMessageW`）」
⇒ **过期属实**，而且**四个子断言分开塌**（本腿逐枚现量）：
| 子断言（09-28 写下） | 今天（HEAD `bd39f17f`）我现量 |
|---|---|
| `ParseComposerRequest` 非 test 零调用方 | **有 1 枚非 test 调用者**＝`internal/panel/composer_dispatch.go:155`（尺＝`grep -rn "ParseComposerRequest(" --include='*.go' cmd internal \| grep -v _test.go`，另两枚命中是 `bridge.go:126` 定义与该调用） |
| `HandleModeRequest` 非 test 零调用方 | **有 1 枚非 test 调用者**＝`internal/panel/composer_dispatch.go:181` `return "", d.Mode.HandleModeRequest(ctx, req)`（当年尺 `grep -rn "\.HandleModeRequest("` 非 test ⇒ No match；今天同一把尺命中这一枚） |
| 全仓非 test 没有一枚 WebView2 消息接收器 | **仓内 Go 文件确为 0**（本腿重跑 `grep -rn -E 'WebMessage\|ReceiveMessage\|OnMessage' --include='*.go' internal cmd tools \| grep -v _test.go`＝空），**但接收器住在依赖库**：`go-webview2@…/pkg/edge/chromium.go:201` `AddWebMessageReceived.Call` ＋ `:233 MessageReceived` → `:240 MessageCallback` → `webview.go:103 chromium.MessageCallback = w.msgcb` → `:139 msgcb` → `:147 callbinding` → `:164 w.bindings[d.Method]`，命中的正是 `panel_host_windows.go:405` 绑进去的 `wispDispatch`。⇒ 那把尺**射程只到本仓 Go**，"整条不存在"是尺的射程造成的读数 |
| `pump.go:16` 的自述"ticket 33 is unclaimed" | 注释今天**逐字还在**（`internal/panel/pump.go:15-16`），但它的前提已过期（票 33 已有 `33-r*`/`33-v2` 落地件）⇒ **注释过期**，属产码，⛔ 本腿一字不改（见 §5 报单） |
⇒ **处置＝票 186 一个字没动**（派单只授权追加票 194）。这一枚的更正该由编排者落还是派腿落，缺的是"授权范围"那一行，写进 §4 判不动。
⚠ 具名分歧（编排者派单说"186 面上我只命中 `:54` 一行普查记录、**那行里没有『零枚』这句**"）：我现量的真值是——`:54` 里**没有**"生产调用者零枚"这六字逐字，**但有**同族的"非 test **零调用方**"＋"整条不存在"两句，而这正是 4f 指认的内容。⇒ 编排者的锚漂在**字面**上，4f 的指认在**内容**上成立。

### 3.2 档二：票面无此句（1 枚）

**③ `.scratch/wisp/issues/219-approval-card-three-reply-buttons-and-a-reason-box.md`**（4f 说"表第 4 行"写着 `ComposerDispatch.Handle` 生产调用者零枚）
⇒ **票面无此句**。本腿尺与读数：
- `grep -n 'ComposerDispatch' .scratch/wisp/issues/219-*.md` ＝ **零命中**（rc=1）。
- 那张表（`:36-40`）第 4 行＝`:39`「**答复通道今天通到哪**」，它答的符号是 `DecideFromPanel`／`DecideFromNative`／`Gate.Veto`，与 `Handle` **不是同一件事**；`:40` 是它的"原句留档（已作废，逐字）"行。
- 票面其余"零调用者"字样（`:15`、`:97`）各自钉的是 `DecideFrom*`/`Veto` 与 `SetAllowedDirs`，⛔ 都不是 `ComposerDispatch.Handle`。
⇒ **处置＝票 219 一个字没动**，并按派语在此**具名写**：「**4f（`pool-validity-4f` §5 第 3 条）报的票 186／219 两处落点里，219 那一处在票面上不存在**（原句：`.scratch/wisp/probes/pool-validity/4f/batch4f.md:208`「票 194 AC#5／票 186 普查注／票 219 表第 4 行都写着…」）。
⚠ **顺带一枚要报（不在本批射程，本腿没动它）**：`:39` 那三枚符号的"生产零调用者"**今天也过期**——`internal/agent/approval/replies.go:328`（`Replies.Allow` 内）／`:408`／`:413`／`:437`（`PanelAllow` 内）／`:463`（`Veto` 内）五处非 test 调用者存在，且驱动它们的正是生产文件 `cmd/wisp/approval_reply.go:215`、`:259`、`:302`、`:304`、`:331` 与 `cmd/wisp/resident_approval_windows.go:612`、`:728`。另：票面引的 `gate.go:622/:610/:371` 今天量到 `:736/:724/:421`＝**行号漂 114/114/50 行**。⇒ 这是**第二枚该追加更正的票**，但它答的是"答复那一侧"，派单红线 (a)/(b) 把它划在本腿射程外 ⇒ 只报不改。

### 3.3 档三：不属实／不许动（10 枚，逐枚具名）

| 落点 | 为什么不属本批"过期"或缺哪行读数 |
|---|---|
| `.scratch/wisp/issues/194-….md:25` | AC#5 **框文内部**（`- [ ] **AC#5 与入向那一跳的依赖要具名**…但**生产调用者今天仍是零枚**（真窗口那一环 H2/H3 未建…）`）。枚数层今天已过期，但改它＝改判据 ⇒ ⛔ 未动，缺"编排者翻 AC#5 措辞"那一行（见 §4-①）。它同句后半"不许拿名册补齐冒充按钮能点"**照旧成立**。 |
| `.scratch/wisp/issues/194-….md:48` | 编排者 09-28「**我复跑并认下的读数**：…`ComposerDispatch.Handle` 入向生产调用者＝**零枚**（只命中它自己的定义与注释）」＝**带日期的历史记录**，红线 (a) ⇒ ⛔ 不许改。⚠ 它当时那把尺（符号名）今天会量到 **25 行**非 test 命中（`grep -rn "ComposerDispatch" internal/ cmd/ --include='*.go' \| grep -v _test.go`）＝正是"构造≠调用"陷阱的现场。 |
| `.scratch/wisp/issues/194-….md:51` | 裁定③「`Handle` 零调用方期间，补齐 18 枚只会交付"门牌挂满、一枚点不动"」＝**排程决定的当时理由**（历史记录，且它引用的 `33-r2` 是别家已落地的票号）⇒ ⛔ 不许改。 |
| `.scratch/wisp/issues/33-panel-host-c27.md:182` | 派单红线 (b) 逐字排除（答"具名诊断入口"那一枚）⇒ ⛔ 一字未动。 |
| `.scratch/wisp/issues/33-panel-host-c27.md:188` | 同上，红线 (b) ⇒ ⛔ 一字未动（`AC#D 的诚实答案：生产调用者＝0 枚，本程没有那枚具名诊断入口`＝"本程"限定时状）。 |
| `docs/evidence/s1/194-method-roster-census-c1.md:134` | 普查腿"本程现量"＋它逐字给出**它自己那把尺**（`grep -rn "ComposerDispatch" internal/ cmd/ --include='*.go' \| grep -v _test.go` 只命中定义与注释 `:90/:97/:120/:137/:167/:182/:192`）。⇒ 历史记录；⛔ `docs/evidence/s1/**` 派单禁改。⚠ 那 7 个行号今天对应 `:114/:121/:153/:175/:215/:230/:240`＝**行号漂**，内容同族未漂。 |
| `docs/evidence/s1/194-method-roster-census-c1.md:216` | 同上（`composer_dispatch.go:120` 现量零枚生产调用者 ⇒ 建议先派 H2/H3）。它给出的**顺序建议已被编排者采纳成裁定③**，属世系记录 ⇒ ⛔ 禁改文件。 |
| `docs/evidence/s1/181-186-git-detection-census-c1.md:50` | 逐字「**裁定：入向那一跳今天不存在**——具名坐实票 114」＋自带锚「现跑到 HEAD `c1c96008`／本件骨架 `e9ef94d0`」＝**带锚的当时读数**（票 186 `:54` 就是引它）⇒ 历史记录；禁改文件。 |
| `docs/evidence/s1/181-186-git-detection-census-c1.md:167` | 同上（"§⑤ 已现量：入向那一跳今天整条不存在…"）⇒ 历史记录；禁改文件。 |
| `docs/evidence/s1/33-minimal-inbound-hop-r1.md:57` | 逐字「**这一程交付完之后**，生产里没有任何人调用 `ComposerDispatch.Handle`。零枚。」＝**限定时状＝"片 A 那一程"**，它本来就只承诺片 A；片 B（`33-inbound-listener-r2/r3`）随后落地了链 B ⇒ 不属过期；且同件 `:60` 已自陈"没结线"。禁改文件。 |

**另外 3 枚射程外的 evidence 句子（登记为"扫到但判为别家符号"，本腿一律未动）**：`docs/evidence/s1/190-file-tree-design-a1.md:221`（引 `internal/panel/git.go:18-20` 的注释，注释今天逐字还在＝**产码注释过期**，射程＝产码，⛔ 不动）、`docs/evidence/s1/33-35-preflight.md:35`（五枚符号并排，其中 `RequestWorkspaceSwitch` 今天仍**真的**零非 test 调用者——本腿尺 `grep -rn 'RequestWorkspaceSwitch(' --include='*.go' cmd internal tools \| grep -v _test.go`＝**只命中定义行 `internal/panel/workspace.go:111`** ⇒ 那一支**不过期**，`ParseComposerRequest`/`HandleModeRequest` 两支过期）、`docs/evidence/s1/panel-l2-grant-nail-accept-r1.md:87`（`HandleModeRequest` 只有测试调用者 ⇒ 今天过期，但它是**验收表里的一行判据陈述**、非票面前提，且落禁改文件）。

---

## 4. 判不动（每一格都写满，⛔ 不留空占位）

**① 票 194 `:25`（AC#5 框文里那句"生产调用者今天仍是零枚"）——过期属实，但要不要改判据归编排者。**
缺的读数不是技术读数（技术读数本腿已给全：2 处非 test 调用者＋两条追到 `main` 的链），缺的是**权限那一行**：`:25` 整行在 `- [ ] **AC#5 …**` 框文射程内，派单第②条逐字禁止原地编辑该行、且"判据归编排者翻"。⇒ 本腿只把更正**追加在文件末尾**（`## 编排者派：过期读数更正（腿 7b2，10-06）`），框文一字未动、`grep -c '^- \[ \]'` 追加前后都＝**7**。要编排者点的就一句：`:25` 的措辞翻不翻（翻＝把"生产调用者今天仍是零枚"改成"枚数已过期、真窗那一发仍待 H2/H3"），⛔ 本腿不替它翻。

**② 票 186 `:54`（"入向那一跳整条不存在"）——过期属实，但派单只授权追加票 194。**
缺的读数＝**授权范围**那一行：派语第②条只写"票 194 一枚"，对 186 明确写"若判定票面根本没有那句话⇒只在件里具名写…⛔ 不改这两枚票一个字"。而本腿的判定是**第三种**（比派单预设的两支都多出来）：`:54` **有**同族句子（"整条不存在"＋"非 test 零调用方"），只是**没有**"生产调用者零枚"那六字逐字。⇒ 两支预设都不覆盖这一形 ⇒ 本腿不动票 186 一个字，并把这一形**具名报**在这里。⚠ 另缺一样：`:54` 是 Progress log 里的普查记录行（带"09-28 普查程 `181-c1` 交件"来路），按红线 (a) 它也可读成**历史记录**——"算不算过期票面前提"这一支的裁量权在编排者，不在本腿。

**③ 编排者派单里的锚点漂移清单（逐枚具名，真值＝本腿现量）**
| 派单给的锚 | 本腿现量真值 | 性质 |
|---|---|---|
| `cmd/wisp/main.go:66` → `runResident()` | **未漂**（`:66` 逐字 `runResident()`） | — |
| `runResident()`（`resident_windows.go:33`） | **未漂** | — |
| `resident_windows.go:151` `newResidentPanelManager(...)` | **未漂**（逐字 `rp, rpErr := newResidentPanelManager(rt.Layout.DataDir)`） | — |
| `cmd/wisp/panel_resident_windows.go:228`/`:234` | **未漂** | — |
| `panel_resident_windows.go:253` `NewPanelManager(disp, …)` | **未漂** | — |
| `panel_host_windows.go:213` 定义 | **未漂** | — |
| `panel_host_windows.go:406` Bind 闭包 | **未漂**；`Bind` 调用本体在 `:405`（派语写"→ `:406`"，`:406` 是闭包体内那一行 `reply, _ := m.dispatchRaw(ctx, raw)`） | 口径差一行，非漂 |
| `:637` `return m.disp.Handle(ctx, raw)` | **未漂** | — |
| `main.go:115 case "panel-inbound":` | **未漂**；`cmdPanelInbound` 的调用行＝**`:120`**（`:119` 是 `attachParentConsole()`）。⚠ 本腿第一次把这一跳写成 `:119`，在票 194 追加节与本件里都已就地改成 `:120`（改的是**本腿自己刚写的那两处**，⛔ 不碰任何原句） | 本腿自纠一枚 |
| `panel_inbound.go:146` 装配 / `:163` `disp.Handle` | **均未漂** | — |
| `composer_dispatch.go:120`＝`Handle`（票面 4 枚句子都引它） | **漂**：今天 `:120` 是 `type ComposerDispatch struct {` 之前的空行、`:121` 才是该声明；`Handle` 真身＝**`:153`**（定义 `func (d *ComposerDispatch) Handle(ctx context.Context, raw string) (string, error)`） | **行号漂 33 行，内容未漂**（票面那句钉的仍是这一枚方法） |
| 4f 的 `git show HEAD:cmd/wisp/resident_windows.go \| grep -n 'first non-test call of NewPanelManager'`＝`:141` | 本腿复跑＝**`:141` 命中**（`HEAD:cmd/wisp/resident_windows.go:141`）＝4f 该锚**未漂** | — |
| `docs/evidence/s1/33-panel-host-c27-v2.md:31` 引的 `panel_host_windows.go:519` | **漂**：本腿现量 `m.disp.Handle` 在 **`:637`**，`Bind(panelDispatchBinding` 在 **`:405`**（它另引的 `:287-290` 本腿不采） | 派语早说过"别把我的行号当凭据"，此处兑现 |
| 4f 指认的"票 219 表第 4 行" | **票面无此句**：`:39` 表第 4 行答的是 `DecideFrom*`/`Veto`，`grep -n 'ComposerDispatch'` 于票 219＝零命中 | 落点不存在 |

**④ "整条不存在"这一支我判不动的边界（诚实登记）**
本腿证到的是：路由器有非 test 调用者、装配链每一跳都在、依赖库里有 `WebMessage` 接收器（`chromium.go:201/:233/:240`＋`webview.go:103/:139/:147/:164`）。本腿**没证到**的是：页面那一次 `postMessage` 真能落到 `Handle`——那需要真窗（票 33 `:42` AC#14 登记"Go→页面那一跳今天没有人投递"，模块 `webview.go:443-448` 只入队＋`PostThreadMessageW`）与 `winlive` 批准，本机今天量不了。⇒ ⛔ 本件任何一节都不写"入向那一跳已通"，只写"枚数层已非零"。若有人要把 `:54` 判成"完全过期"，缺的就是这一行运行期读数。

**⑤ 产码注释那两枚过期自述，本腿不动、只报**（射程＝产码，派单未授权本腿改 `internal/**`／`cmd/**`，且另有一枚腿在 `cmd/wisp` 跑突变）：
- `internal/panel/pump.go:15-16` 「no WebView2 host (ticket 33 is unclaimed), no postMessage writer」
- `internal/panel/git.go:17-20` 「that hop does not exist in this tree, which is why SwitchBlocked below is a constant sentence」
⇒ 票 186 `:54` 与 `190-file-tree-design-a1.md:221` 都**引用**了这两枚注释；注释过期 ⇒ 引用它们的句子看起来"有凭据"这件事本身要打折。**处置**＝本腿不碰产码一字，交给编排者决定派不派"注释更正腿"。

---

## 5. 处置清单（要编排者点的，共 3 项；本腿一律未擅自做）

1. **票 194 `:25` AC#5 框文措辞翻不翻**（见 §4-①）——本腿只追加，未改框文。
2. **票 186 `:54` 该由谁落更正**（见 §4-②，派单两支预设都不覆盖"内容在、字面不在"这一形）。
3. **票 219 的 `:39`／`:15` 那两枚"答复侧零调用者"是第二枚过期票面**（§3.2 末⚠行，含 5 处非 test 调用者与 3 枚生产驱动点、外加 `gate.go` 行号漂 114/114/50）——本批射程外，⛔ 本腿一字未动，是否另派一枚腿由编排者定。
（顺带登记，非处置项：4f §5 第 3 条对 **186/219** 的指认需要按本件 §3 打折——186 是"同族句子在、逐字句面不在"，219 是"票面无此句"。）

---

## 6. 门禁三把尺（全部本腿现量）

**① 起手＝终态名册尺**
- 起手（19:1x 第一次，原文照录）：`git --no-pager status --porcelain -- cmd internal scripts .github docs` → **空集（0 行）**。
- 交件前终态：同一把尺 → 见 §7 的终跑读数；⛔ 本腿**没跑任何 `go` 命令**、没碰 `cmd internal scripts .github docs` 里任何一字（本腿写面只有 `.scratch/wisp/issues/194-…md` 与 `.scratch/wisp/probes/expired-premises/7b2/**`）。

**② 占位符尺**（⚠ 尺里用**方括号类**写法，为的是这一行本身不自造命中；三枚占位 token 的语义与原尺逐字相同）
```
grep -c '[待]填\|（[待]\|填写[中]' .scratch/wisp/probes/expired-premises/7b2/verdict.md
```
→ 我现量＝**0**（rc=1，零命中）。同一把尺打在票 194 上＝**0**（rc=1）。logs 三枚＋本腿自取两枚扫描件＝**0**（`grep -rn` 于 `7b2/logs/` 计数 0）。⇒ §3/§4/§5 每格都是写满的正文，⛔ 没有任何一格留占位词。

**③ 零删除尺**（本腿唯一动过的既有票面）
```
git --no-pager diff --numstat -- .scratch/wisp/issues/194-the-panel-has-two-method-rosters-that-do-not-know-each-other-owner-ruled-align-the-code-to-the-spec.md
```
→ **`45  0  <全名>`**（新增 45 行／**删除列＝0**），追加前后 `grep -c '^- \[ \]'`＝**7→7**、文件行数 `55→100`。⚠ 这台仓栽过的形：追加时把下一行标题当锚点却不复带它＝删行——本腿锚的是**文件最后一行**（`- **本票状态**：…`）并在 new_string 里**逐字带回**它，读数即证。⚠ 读数随写面推进而变：本腿第一次量＝`40 0`，写完追加节 §⑦（补 4f 落点分歧与新锚复跑那 5 行）后终值＝**`45 0`**，删除列两次都＝0。

---

## 7. 终跑（交件前现量，19:5x；锚 HEAD 已在线上被人推进：起手 `bd39f17f` → 终态 `76370fb5`）

⚠ **具名一枚"共享工作树里 HEAD 会漂"**：本腿起手 `git rev-parse --short HEAD`＝`bd39f17f`（19:2x），交件前＝`76370fb5`（19:5x）⇒ 本腿把主尺**在新 HEAD 上重跑了一遍**，读数不变：
- `git --no-pager grep -n 'disp\.Handle(' HEAD -- cmd internal | grep -v '_test.go'` → 仍＝**2 处**（`cmd/wisp/panel_host_windows.go:637`＋`cmd/wisp/panel_inbound.go:163`）。
- 链 A／链 B 的 20 枚锚点逐行 `sed -n '<l>p'` 复量 → **全部仍在那一行上**（含 `resident_windows.go:151`、`panel_resident_windows.go:215/228/234/253/378/383/409/427/436`、`panel_host_windows.go:152/213/405/488/630/637`、`main.go:66/115`、`panel_inbound.go:103/146/163`）。⇒ 追加节里写的 file:line 对**新 HEAD 也成立**，不是只对起手锚成立。
- 宽形状尺 `grep -rn -E '\.Handle\(ctx' --include='*.go' cmd internal | grep -v '_test.go' | wc -l` → 仍＝**9**（其中 7 枚是 slog 同名 `Handle`）。

- **门禁① 终态**：起手（19:1x）`git --no-pager status --porcelain -- cmd internal scripts .github docs` ＝ **空集**。⚠ 本腿交件前**同尺复跑多出一枚不是本腿的条目**＝`M docs/reports/pending-and-issues.md`（`git --no-pager diff --numstat` 量它＝`25 0`，新增内容是 `## A648` 一节 ⇒ 别家正在往台账追加）⇒ 按 `A374` **具名登记为"不是我的"、不 add、不还原、不评论**（派单逐字禁本腿改台账，本腿一字未碰）。除这一枚之外，`cmd internal`（产码）那两片**仍逐字等于起手空集**，⛔ 本腿对 `cmd internal scripts .github docs` **零字节改动**、**未跑任何 `go` 命令**（无 `go test`／`go build`／`go vet`／`go run`）。⇒ 差集：本腿名下＝**空**；共享树里别人在飞＝**1 枚（已具名）**。
- **门禁② 占位符尺**：尺＝派单给的那把（`grep -c` 三枚占位 token：`[待]填`／`（[待]`／`填写[中]`，方括号写法只为让**本行不自造一次命中**，语义与原尺逐字相同），打成 `grep -c '[待]填\|（[待]\|填写[中]' .scratch/wisp/probes/expired-premises/7b2/verdict.md` → **0**（rc=1 零命中；`grep -c` 零命中返回 rc=1，故本腿用 `;` 串接不用 `&&`）。票 194 同尺 → **0**（rc=1）。logs 目录 → **0**。⚠ 本腿在写 §7 之前先用**原尺逐字形**扫过本件，抓到 §6 那行自造命中 1 枚 ⇒ 已把两处尺文改写成方括号形并复跑归零；这一发是"尺自己也会咬人"的现场，具名登记在这里而不是抹掉。⇒ §3/§4/§5 每格都是写满的正文，没有任何一格留占位词。
- **门禁③ 零删除尺**：`git --no-pager diff --numstat -- .scratch/wisp/issues/194-the-panel-has-two-method-rosters-that-do-not-know-each-other-owner-ruled-align-the-code-to-the-spec.md` → 终值＝**`45  0  <全名>`**（删除列＝**0**；本腿第一次量＝`40 0`，写完追加节 §⑦ 后＝`45 0`，两支删除列都＝0）。票 194 行数 `55 → 100`；`grep -c '^- \[ \]'` → **7 → 7**（AC 七格一格未勾、一格未改、一格未新增）。
- **本腿写面名册（逐枚点名，⛔ 无 `git add -A`）**：
  1. `.scratch/wisp/issues/194-the-panel-has-two-method-rosters-that-do-not-know-each-other-owner-ruled-align-the-code-to-the-spec.md`（末尾追加一节）
  2. `.scratch/wisp/probes/expired-premises/7b2/verdict.md`（本件）
  3. `.scratch/wisp/probes/expired-premises/7b2/logs/scan-sentences.txt`（主尺＋变体尺扫描输出，126 行）
  4. `.scratch/wisp/probes/expired-premises/7b2/logs/scan-handle-sentences.txt`（二次射程尺输出，36 行）
- **不是本腿所为、本腿未动未还原**（`A374` 形，具名登记）：`git status --porcelain -- .scratch/wisp/issues` 里另有 `M .scratch/wisp/issues/271-the-ticket-221-task-family-count-ruler-false-reds-on-a-legitimate-third-tool.md` —— 该票最后一条 commit 是 `257275be`（"立票 270/271"），**属别家在飞的活**，本腿不 add、不还原、不评论。
- **禁改面自证**：`docs/evidence/s1/**`（含另一枚搬运腿刚动过的 `closed-tickets-evidence-index.md`）、`docs/reports/pending-and-issues.md`、`AGENTS.md`、`docs/PLAN.md`、`docs/specs/**`、`internal/**`、`cmd/**`、票 33／186／219 票面 ⇒ **一字节未改**；未重命名／未撤销任何 `-done`；未建 worktree、未在仓库目录外 checkout；禁读面（`.scratch/wisp/probes/232/**`、`pool-validity/5g2/**`、`frontend/**`、`design/**`）**未读未引**（`4f/batch4f.md` 是派单指定要读的料，属 `pool-validity/4f`，不在禁读名册里）。

