# 35-a6 — 六格可开工性普查（verdict）

腿：`35-a6`，只读普查。锚 HEAD `7a367a08292ef94908d86656608974e81c42f0d4`（`dev`），起手锚件 `00-anchor.md`（commit `dcbb93d2`）。
全程 **零 `go` 命令**（⛔ 未跑 `go build`/`go vet`/`go test`；`go env`/`go list` 也**一条没跑**——本腿所有结论只需静态读盘，跑它们不会多取任何凭据，故按"自报＝零条"登记）。
⛔ 未开任何真窗、未 `git push`、AC 框一字未动。

引框规矩沿用 `A684` §2：**行号＋名字**，不用"AC#N"。

---

## ① `:42` — Whitelist fuzz（50 枚随机/禁名 ⇒ 全拒＋记日志；在册每枚都派发）

### 1. 被测物：有（前半格），且分两半

拒＋日志那一半今天真在产码里：

- `internal/panel/bridge.go:132` 逐字 `	if !knownComposerMethod(r.Method) {`
- `internal/panel/bridge.go:133` 逐字 `		return r, fmt.Errorf("%w: 方法 %q 不是面板 composer 通路的能力入口", ErrComposerRequest, r.Method)`
- `internal/panel/composer_dispatch.go:155-159`（`Handle` 里 parse 拒了以后）逐字 `	req, err := ParseComposerRequest(raw)` → `d.record(req, err)` → `return RefusedEnvelopeForUser(req, err), err`
- `internal/panel/composer_dispatch.go:244` 逐字 `	d.Audit("panel: INBOUND-DISPATCH request=%q method=%q source=%q origin=%q err=%v detail=%q",` ⇒ **"记日志"有载体**，且这条串今天被两枚既有尺认：`cmd/wisp/panel_inbound_33_test.go:141`、`cmd/wisp/panel_config_248_test.go:452`。
- 名册真身（方法名册现量，非派单转述）：`internal/panel/bridge.go:42-45` 四枚 composer ＋ `:66`/`:67` `config.get`/`config.set` ＝ **6 枚**；`knownComposerMethod` 的 case 表在 `:148`（一行列全六枚）。
- 派发表：`internal/panel/composer_dispatch.go:176-209` 六枚 case ＋ `default: d.rosterMismatch(req)`。
- 生产装配（谁真接了处理器）：`cmd/wisp/panel_inbound.go:271-281` —— `Mode: modeWrites`、`Config: configWrites`、`Workspace/Attachment/Message: nil` ⇒ 那三枚走 `composer_dispatch.go:230-235 unattached`（`ErrNoHandlerAttached`，按名拒＋审计，不是静默丢）。

### 2. 最少落点

**前半格（50 枚随机名 ⇒ 拒＋逐枚审计行）**：1 枚新 `*_test.go`（`cmd/wisp/` 或 `internal/panel/`），估 **120–200 行**，⛔ 零产码改动。
**后半格（"each listed method dispatches"）**：见第 3 条——按 SPEC-08 那张表读**不是落点问题而是名册问题**，本腿不给落点。

### 3. ⛔ 要不要新增 C17 契约面 —— **要，且这半格不归本票**

这格的"listed"有两种读法，结论完全相反，必须摊开：

| 读法 | 名册出处 | 今天派发枚数 | 要不要新名 |
|---|---|---|---|
| 甲：**按 Go 今天的名册** | `bridge.go:42-45/:66-:67`（6 枚） | 6/6 都进 `dispatch` 的 case | **零枚新名** |
| 乙：**按票面 `:18` 那句 "Method whitelist per SPEC-08 §5.2 table"** | `docs/specs/SPEC-08-ui-ball-panel.md:163-174`（18 枚 invoke/推送名） | **只有 2 枚**在册（`config.get`/`config.set`）；那 4 枚 composer 名**不在 SPEC-08 表里** | **16 枚新名** |

⇒ 读法乙＝**改 C17 白名单＝人工批准**（`AGENTS.md` §0 第 2 句；票 `:163` 编排者自己写的落地禁区逐字"不许新增点分法名（N1 会在加名的那一刻红，而名册属 C17 面＝人工批准）"）。
⇒ 而且这一格**已经有专门的票**：`.scratch/wisp/issues/194-...-align-the-code-to-the-spec.md` 标题就是"两块名册互不相识，owner 裁：代码向规格对齐"，它 **AC#1–AC#7 七枚框全未勾**，其 AC#3 逐字"白名单扩张要**逐枚**落 `A##`（文件／行／理由／边界／撤销口令），不许拿"owner 批了名册"当全域通行证一次洗平"。⇒ **本腿判：`:42` 的后半格归口票 194，⛔ 不该在票 35 里重复列一份落点**（同一机制链只记一次）。

**要机主的那一句（只在他要走读法乙时才需要）**：
> "票 35 `:42` 的『每枚在册方法都能派发』要不要按 `SPEC-08:163-174` 那张表兑现？要，就得给 C17 入向名册补 16 枚——其中 4 枚今天**连后端能力都没有**（`panel.resync`／`transcript.get`／`models.list`／`models.delete`，票 194 已具名），3 枚撞票 194 AC#2 的三条禁区；不补的话，`:42` 只能按 Go 今天的 6 枚名册读。"

⚠ 补一句盘上事实给编排者：任何新增名会当场红 `internal/panel/inbound_roster_253_test.go:59` 那枚 `wantFullInboundRosterSize253 = 6`（该文件 `:55-56` 自己写明"这是**两枚必须由 diff 故意移动**的数"，实际断言在 `:445` 与 `:787`），并动 `internal/panel/git_test.go:387`/`:519` 的四枚 `panel.*` pin。**这不是"顺手改测试"**。

### 4. 有没有仪器能钉住它

- 今天**没有**一枚尺打"50 枚随机名＋逐枚审计行"这一格（`grep -rniE 'fuzz|randomMethod' --include=*_test.go internal/panel cmd/wisp` 无相关命中）。⇒ 交完**有牙**，不是自述。
- 最近的既有尺（不许拿来抵账，但要知道它们存在）：
  - `internal/panel/l2_grant_boundary_test.go`（`bridge.go:148` 的 case 表与之对拉；票 `:186` 记它 `grantRouteWords` 10 枚、`:1251-1267` 要求 `knownComposerMethod` 拒 11 枚候选名）
  - `internal/panel/inbound_roster_253_test.go`（四信号名册尺）
  - `internal/panel/composer_dispatch_test.go:249` `TestEveryListedMethodWithoutAHandlerIsRefusedNotDropped`（"在册但没处理器＝按名拒"已有钉）
  - `cmd/wisp/panel_inbound_guards_35r3_test.go:150` `TestInboundRosterGuardRefusesUnregisteredNameOnPageEdge`（经页面那条边的名册拒）
- 本腿具名仪器（四件全在，blob 见 `00-anchor.md`）：`cmd/wisp/panel_transport_35r1_test.go:195`／`panel_transport_35r2_test.go`（9 枚 `func Test`，`:1736-:2097`）／`panel_inbound_guards_35r3_test.go:150/:158/:167`／`panel_transport_live_35v2_windows_test.go:172`。

---

## ② `:44` — Concurrency routing（10 枚并发 invoke、回执乱序 ⇒ 配对正确）

### 1. 被测物：**一块没写**（不是"写了零调用者"）

- 回执今天**只能同步返回**：`cmd/wisp/panel_host_windows.go:713-716` 逐字
  `	if err := w.Bind(panelDispatchBinding, func(raw string) string {` → `		reply, _ := m.dispatchRaw(ctx, raw)` → `		return reply`。
- Go 侧**没有按 id 路由回执的挂起表**：全仓非 test 找 pending/map-by-request-id 的形状，`internal/panel/composer_dispatch.go` 的 `ComposerDispatch` 结构体（`:121-139`）只有五枚处理器字段＋`Audit`，**没有任何请求态容器**；`Handle`/`dispatch` 每次调用自带返回。
- 页面侧的配对由**库自己**的 `window._rpc[seq]` 做（票 `:198` 逐字引 `webview.go:462-478`），⛔ 那不是本仓产码，本票不能拿它抵账。
- 有一枚"写了但零调用者"的邻居，账要分开记：`internal/panel/bridge.go:80 func NewRequestID() string` 非 test 调用者 **0**（尺＝`grep -rn "NewRequestID" --include=*.go .` 剔 `.scratch` ⇒ 只剩定义 `:80` 与 `internal/panel/bridge_test.go:48`）。页面自己铸号（票 `:70` 记 `panel.ts:211` 的 `requestId: pc-<seq>-<uuid>`），所以 Go 这枚是**死码**。
- 结构性事实（这条决定"乱序"能不能发生）：门是在库的 `msgcb` 里被调，跑在面板线程上（`cmd/wisp/panel_resident_windows.go:148` 那条 `tasks chan func()` 与 `w.Dispatch(fn)` 都是单线程泵），⇒ **真页面上今天不存在"10 枚并发 invoke 乱序回执"这条路径**。

### 2. 最少落点

- 读法甲（按"10 枚 goroutine 并发打 `ComposerDispatch.Handle`，处理器故意乱序答，每枚调用方拿回自己那条回执"）：**1 枚新 `internal/panel/*_test.go`，约 100–160 行**，零产码。⚠ 但本腿要具名警告这一读的牙很薄：`Handle` 今天**没有可被串台污染的状态**，测出来会像"恒真"。真正的辨别力要靠定向突变（把 `raw` 与 `req.RequestID` 的配对打乱）。
- 读法乙（票面字面：回执**异步**乱序 ⇒ 按 correlationId 配对）：必须先有"Go 主动推给页面"那一跳＝`:47`/`:49` 同一枚前置。落点 **3–4 枚文件、300–600 行**：`panel_host_windows.go`（把回执从同步返回改成异步推）、`internal/panel/pump.go` 或新建一枚挂起表、`cmd/wisp/panel_pump.go`（接字节）、加判据。

### 3. 要不要新增 C17 契约面 —— **不要新方法名**

C17 的契约本体（`docs/PLAN.md:1367` 逐字"**回复必须按 correlationId 路由**"）**已经要求这件事**，读法乙是**兑现契约**、不是改契约；它需要的是**出向通道**，不是新入向名。⛔ 例外：若异步回执要借一个**点分事件名**（`task.delta` 之类）作为帧类型并把它写进 `bridge.go` 的 `Method*` 常量块，就会落进 `inbound_roster_253_test.go` 的 DECLARED 尺（枚数 6→7）⇒ 那一支才要人工批准。**要机主的那一句＝暂无**（除非落地时选形选到"必须入册"那一支）。

### 4. 有没有仪器

今天**零枚**尺打在这条边上：`internal/panel/` 里唯一的并发尺是 `subagent_stream_197_test.go:153` 的 `sync.WaitGroup`，它打的是 `StreamLog`、不是 `Handle`。四枚本票具名仪器（r1/r2/r3/v2）都不测并发配对。⇒ 交完**会变色**（新尺自己就是牙），但读法甲那版要先过"恒真自查"（发一枚中和配对的突变）。

---

## ③ `:45` — Forged panel allow（网页直发 `approval.decide{allow:true}` ⇒ 服务端拒）

### ★查重结论：**同一枚物理缺陷，本框归口，⛔ 不另列落点**

派单问"这格和 `35-r3`/`35-r4` 那三枚入向守卫是不是同一枚物理缺陷"。现量答：**是，而且是两层各自都已有牙的那种同一枚**——

**(A) 路由层＝与 `:63` (d)-1 完全同一枚。** 决定性凭据（本腿现量，`cmd/wisp/panel_inbound_guards_35r3_test.go:80` 逐字）：
```
	envelope35r3UnregisteredName = `{"method":"panel.approval.request","correlationId":"pc-35r3-approval-1","outcome":"allow"}`
```
⇒ (d)-1 那枚用例**今天就在拿一枚带 `outcome:"allow"` 的审批形状信封，从页面真用的那条 postMessage 边打进去，断它在 `bridge.go:132` 被名册判拒**（用例 `:150`，断语含票 `:71` 钉的逐字红句"不是面板 composer 通路的能力入口"，且 `spy.reqs==0`）。这正是 `:45` 那句"direct bridge call … from the webview context → server rejects"的**路由层**。
⚠ 同时本腿要**顶回派单的一处行号归属**（见 §③末）：那三枚守卫不在 `cmd/wisp/panel_host_windows.go`，而在 `internal/panel/bridge.go`。

**(B) 能力层＝早于本票就有钉。** 被测物在，且有生产调用者：
- `internal/agent/approval/gate.go:736 func (g *Gate) DecideFromPanel(ctx context.Context, r Request) error`，`:737 if r.Allow {` → `:738` 打 `PANEL-ALLOW-REJECTED` → `:741-743 g.q.revokeGrants(r.CorrelationID)`（烧 nonce）→ `:746 return fmtw(ErrPanelAllow, "面板来源的「允许」被服务端 API 直接拒绝（F2 第三层），可改为拒绝或查看完整参数")`
- `internal/agent/approval/ui.go:133 ErrPanelAllow = errors.New("approval: 面板来源不得允许（F2 第三层：允许只接受原生侧）")`
- 生产链路（不是零调用者）：`internal/agent/approval/replies.go:431 func (r *Replies) PanelAllow(...)` → `:437 err := g.DecideFromPanel(ctx, Request{`；宿主侧调用者 `cmd/wisp/approval_reply.go:331 offered, err := s.live.h.PanelAllow(s.ctx, corr)`。

**既有尺具名（6 枚 Go 边界尺 ＋ 2 枚静态尺）：**
- `internal/agent/approval/ticket97_alias_direction_test.go:97`/`:100`（`t.Errorf("SECURITY: DecideFromPanel(%q, allow=true, 真 grant) 返回 nil")`）
- `internal/agent/approval/queue_test.go:107`/`:110`（期望 `ErrPanelAllow`）
- `internal/agent/approval/ticket259_panel_capability_rulers_test.go:316-319`（AC#3 RED：从面板可读材料铸 allow 必须被拒）
- `cmd/wisp/approval_seam_201_test.go:169`（`errors.Is(err, approval.ErrPanelAllow)`）
- `cmd/wisp/subagent_selfapproval_197_test.go:469`
- `cmd/wisp/panel_inbound_guards_35r3_test.go:150`（=上面 (A)）
- `internal/panel/l2_grant_boundary_test.go:1251-1256`（候选名册第一枚逐字 `"panel.approval.request", // verbatim from frontend/src/lib/panel.ts:181`）与 `:1962-1966`（**线形清单里逐枚写着** `{"method":"approval.decide","requestId":"r-1","source":"panel-composer","outcome":"grant"}` ＋ `{"method":"panel.l2.allow",...,"allow":true}`）⇒ **字面 `approval.decide` 的 allow 形状今天就在被拒**；⚠ 区别（也是票 `:71` 早写过的那句）：这枚尺是**包内直调名册函数**、不经传输边，所以它抵不了 `:63` 的"经页面那条边到达＋点名哪一道拒"，但**抵得了 `:45` 的"这枚名字被不被拒"**。
- 静态 CI：**d22scan ban #6 `panel-approval`**——`tools/d22scan/main.go:23` 定义"`approval.decide` in frontend/ (D33/F2: allow decisions are native-side only)"、`:271` 真走 `frontend/`、空面致命（`:418`）
- `internal/panel/frontend_hygiene_test.go:281` `TestFrontendNeverNamesAnApprovalDecision`（走 `frontend/` 全树逐行匹 `panelDecisionIdentifierRe`）

⇒ **落点＝0 枚文件。** 按"同一机制链只记一次"，本框应记**归口 `:63` (d)-1（路由层）＋ `ErrPanelAllow` 那族既有钉（能力层）**。

### 唯一没被覆盖的残差（具名，⛔ 不许顺手做）

票面 `:45-46` 那句里的 **"＋ browser-context e2e"** 那一支：
1. 要真窗 ⇒ 那一维今天只有 `-tags winlive`（`cmd/wisp/panel_transport_live_35v2_windows_test.go`）取得到，票 `:274`/`:375` 已裁成〔仅本机可量〕＋逐波复跑；
2. 要真把 `approval.decide` 发进桥，就得先把它**入册** ⇒ **改 C17 白名单＝人工批准**；⚠ 且票 194 **AC#2①** 逐字禁这条："`approval.decide` 的「allow」侧永不允许从面板发起（`SPEC-08:169` 逐字…出处 D33／`Q-49` 那一族）"。
⇒ 这一支今天**不可派**：它的实现前置不是"少写一枚测试"，是"少一枚名册条目＋少一段后端能力"。票 `:73` 已把这族（"面板上点拒绝到底能不能真拒绝"）并进 `A671`〔待人拍板〕，本腿复认这个归口，不另开账。

**要机主的那一句**：只有当他想让**面板侧的 `reject` 真的能拒**时才需要（`gate.go:748` 那半句 `g.q.reject(…)` 能力在、路由不在）；`:45` 这格本身（防伪造 allow）**不需要**他批任何东西，因为它已经有牙。

---

## ④ `:47` — Resync（hide → 状态变 → show 显示新状态；SDK 禁跨 show 缓存）

### 1. 被测物：行为半格**一块没写**；lint 半格**在、且有牙**

词边界尺（⛔ 不用票 `:95` 那把大小写不敏感子串尺）：
```
grep -rniE '\bresync\b' --include=*.go cmd internal tools | wc -l   →  0        rc=0
git grep -nE '\bresync\b' HEAD -- '*.go' '*.ts' '*.tsx'  | wc -l    →  1        rc=0
```
那 1 枚＝`frontend/src/lib/panel-views.ts:90`，逐字："PanelSnapshot yet - ticket 35's pump owns it, and `panel.resync` (the push"…`:91-92` 续"that would re-establish it on every show) is still **zero-hit repo-wide**" ⇒ **页面在等本票，不是页面欠着**。

出向那一跳今天为什么不可能（四把尺，全非 test）：
- `grep -rn --include=*.go '\.Eval(' cmd internal | grep -v _test.go | wc -l` ＝ **0**
- `grep -rn --include=*.go -E 'PostWebMessage|EvaluateScript|CreateWebMessageAsJson' cmd internal | grep -v _test.go | wc -l` ＝ **0**
- 传输接口只两枚方法：`cmd/wisp/panel_host_windows.go:652-655` 逐字 `type pageTransport interface {`／`	Bind(name string, f interface{}) error`／`	Init(js string)`／`}` ⇒ **没有出向方法可叫**
- `Show()` 不推任何东西：`cmd/wisp/panel_host_windows.go:499-526`（体内只有 `bringUp`／`pnlShowWindow.Call`／`pnlUpdateWindow.Call`／`pnlSetForeground.Call`／`m.shown = true`）
- 泵出来的字节没人读：`cmd/wisp/panel_pump.go:383 func (rt *agentRuntime) lastPanelSnapshot() (panel.Snapshot, []byte, bool)` 是**定义**；同名符号带括号的调用共 12 行命中，其中非 test **0 处**（`grep -rn "lastPanelSnapshot()" --include=*.go cmd | grep -v _test.go` 只回定义那一行）⇒ **写了但零调用者**那一族。

**lint 半格已有仪器**（这格唯一的"在"）：`internal/panel/frontend_hygiene_test.go:169 TestPanelFrontendIsStateless`，扫 7 枚禁存面（`:49-55` `localStorage`/`sessionStorage`/`indexedDB`/`document.cookie`/`caches.(open|match|keys)`/`navigator.serviceWorker`/`fs`），且 `:171-174` 空文件集 `t.Fatal("no files under frontend/src - this check would pass by seeing nothing")` ⇒ 不是恒真尺。

### 2. 最少落点（真要做）

**≥4 枚文件**（Go 侧 200–500 行 ＋ 页面侧）：
1. `cmd/wisp/panel_host_windows.go`：`pageTransport` 加一枚出向方法（`Eval` 或消息发送），`Show`/`HotShow` 里发起"每次 show 全量推"。估 60–140 行。
2. `internal/panel/pump.go` ＋ `cmd/wisp/panel_pump.go`：把已装配好的 `Snapshot`/`Marshal` 字节从 ledger 那一路接到 show 那一跳（今天 `Publish` 非 test 出口只到 `run.go` 的 booking）。估 80–180 行。
3. `frontend/src/lib/panel.ts` ＋ `frontend/src/App.tsx`：**页面耳朵**（`addEventListener("message")` 今天 0 处；`panel.ts:141` 只有 `postMessage(message: string): void` 这一枚发送声明）。⛔ **页面写权未放开**——票 `:97`/`:140`/`:163` 三处都记着"机主放开的只有 `frontend/**` **只读**"、"⛔ 不动 `frontend/**` 一个字节"。
4. 判据：hide→mutate→show 的 stale-state 测试（载体只有两形：`panel_transport_35r2_test.go` 那把行为尺，或 `-tags winlive` 真窗＝〔仅本机可量〕）。

### 3. 要不要新增 C17 契约面 —— **有一个陷阱，必须让编排者知道**

`panel.resync` 在 `SPEC-08:163` 已挂名、方向逐字标 "**Go→前端推送**"，所以它**不是**入向白名单成员、原则上不需要进 `bridge.go`。
⚠ 但**如果**落地时把它铸成 `bridge.go` 里的 `Method…` 常量（最顺手写法），就会落进 `internal/panel/inbound_roster_253_test.go` 的 DECLARED 信号（该文件 `:25-26`："every package-level string constant of bridge.go whose identifier begins with "Method" **OR** whose VALUE is a dotted name"），枚数 pin `:59 wantFullInboundRosterSize253 = 6` 当场红（断言体在 `:445`/`:787`）。
⇒ **要机主的那一句（只在选形落到"必须进 bridge.go 常量块"时才需要）**：
> "resync 的推送名 `panel.resync` 要不要成为 C17 入向名册的一员（那就是 6→7，`inbound_roster_253` 那枚刻意钉住的数就得动）？还是只做**出向帧类型**、不进名册？"
本腿不替他选。

### 4. 有没有仪器

行为半格：**今天零枚尺会因此变色**——四枚具名仪器（r1/r2/r3/v2）全测**入向**，没有一枚断言"Go 推了东西给页面"。⇒ 若实现只交"判据跑绿"而无出向调用点，就是票 `:99` 警告的那枚**自证通过**形状（该节编排者自己写的："任何声称覆盖出向的判据，必须能指名它断言了哪个 Go→页调用点"）。lint 半格已有 `TestPanelFrontendIsStateless` 在钉。

---

## ⑤ `:49` — Backpressure（blocked consumer 下灌事件 ⇒ 合并、无无上界内存、内容完整）

### 1. 被测物：**部分在，且钉的方向与 AC 相反**

在的那部分（有界队列＋背压计数）：
- `internal/panel/pump.go:397` 逐字 `// Overflow TRUNCATES and never merges (ticket 197 leg B re-cut this rule, which`、`:398-399` 续"previously folded the two oldest chunks into one row"
- 三枚计数出口：`internal/panel/pump.go:632 Truncated()`、`:644 ElidedRunes()`、`:657 DroppedKeys()`（`:655-656` 逐字"and the host can say WHICH ones and say they were never merged into a row that is / still on screen"）
- 面板线程队列有界：`cmd/wisp/panel_resident_windows.go:148 tasks:    make(chan func(), 16)`，满了不排队而是**计数并丢**：`:369 rp.failedPost.Add(1)`，字段声明 `:130 failedPost atomic.Int64`。
⚠ `failedPost` 今天**零测试引用**（`grep -rn "failedPost" --include=*_test.go cmd` ＝ 0 命中）＝"写了但没尺"那一族。

不在的那部分：票面 `:27` 那句 "Event push: **bounded queue per panel**; overflow MERGES increments" 里的那枚**每面板事件推送队列**＝**一块没写**（出向整条零产码，见 `:47` 那四把尺）。

### ★⛔ 必须先报的冲突（未定义即停，不是可派格）

AC `:49` 要 "merges"；盘上钉的是 "**TRUNCATES and never merges**"，而且这枚规则是**票 197 leg B 重新裁过的**，并被命名测试钉死：
- `internal/panel/pump_test.go:165` `TestTheStreamLogTruncatesInsteadOfMerging`
- `internal/panel/subagent_stream_197_test.go:209` `TestOverflowTruncatesEachStreamsOwnHeadAndTail`
- `internal/panel/subagent_stream_197_test.go:275` `TestHardCeilingDropsWholeKeysAndNamesThem`
理由逐字在 `pump.go:400-406`：子代理的流不是可互换文本，把两行折一行会让用户看到**关于来源的假话**（同一枚病＝D30）。

⇒ 照 `:49` 字面实现**只有两条路**：① 放宽/改掉上面那三枚既有钉 ⇒ ⛔ 违反 `AGENTS.md` §1.1"不许为了变绿放宽任何断言"；② 把 merge 只用在**新的事件推**上、不碰 StreamLog ⇒ 那需要先有出向通道（＝`:47` 的前置）。
⇒ **本腿判：`:49` 今天不可直接派**。要么先把 AC 措辞按票 197 的裁定改写（措辞改动＝动票面文字，⛔ 本腿不动，交编排者/机主），要么按②把它降为"`:47` 之后的第二跳"。
**要机主的那一句**：
> "票 35 `:49` 说溢出要『合并』，票 197 leg B 已经把同一族规则改成『截断、永不合并』并用 `TestTheStreamLogTruncatesInsteadOfMerging` 钉住。新的事件推送要 merge 还是 truncate？两枚 AC 现在互相矛盾，我不替你裁。"

### 2. 最少落点（若按上面②那条走）
2–3 枚文件、**250–450 行**：一枚每面板有界队列（新，`internal/panel/` 或 `cmd/wisp/`）＋出向那一跳（＝`:47` 落点 1/2 复用）＋判据（blocked consumer 注入、断 `Truncated/Elided/Dropped` 三计数与堆上界）。⛔ 前置是 `:47`，单独派会造出"没有被测物的判据"。

### 3. 要不要新增 C17 契约面
**不要**（事件名那一行 `SPEC-08:174` 已列，但"合并/背压"是**载荷之上的队列策略**，不需要新枚方法名）。⚠ 唯一例外同 `:47`：若把 `task.delta`/`tool.chip`/`ball.state`/`cost.tick` 铸进 `bridge.go` 常量块就会动那枚名册 pin。

### 4. 有没有仪器
今天唯一有牙的三枚尺（上面具名）**咬的方向与 AC 相反**：谁把策略改回 merge，它们当场红。`failedPost` 那枚背压计数**没有任何尺**。⇒ 这格若实现，需要**新**判据；旧尺不会为它变绿、也不会为它变红。

---

## ⑥ `:51` — No secret leakage scan across bridge payloads（桥载荷不泄密）

### ★查重结论：**已有成套同类尺，本框大部分归口票 248 那族**

按"这句话在说什么"扫（⛔ 不只搜符号名）后现量到的尺：

- `cmd/wisp/panel_config_248_test.go:36-39 const canary248 = "sk-canary248notarealkey0f2a9b7c"`；`:41-46 var secretShape248 = []*regexp.Regexp{ regexp.MustCompile("sk-[A-Za-z0-9]{12,}"), regexp.MustCompile(canary248) }` ＝ **形状尺不是名字尺**（文件头 `:12-13` 自己写明"a made-up string… must not become any real credential"，`:20-21` 写明失败路径**不回显**命中文本，否则测试日志自己就是被检面）
- `:154 TestAC2CredentialSentinelAppearsInNoArtifact`，扫面（`:186-191` 逐行）＝ **receipt（桥回执）／audit／slog-sink／snapshot（桥载荷）／ledger-line／data-root** 六面＋写成功正控
- `:218 TestAC2LeakRulerFiresWhenTheCanaryIsReallyThere` ＝ **恒真自查**（`:228`/`:246`/`:256` 三处真把 canary 种进 `ResultChunk`/`NativeVerdict`/审计行，断尺必须响）
- `:329 TestAC2SharedEnvelopeCannotCarryTheCredentialValue`（结构面：共享封套不许带凭私值键）／`:510 TestAC2SnapshotReportsRefWithoutBlobAsAPartState`
- `internal/panel/config_route_248_test.go:349 const canary = "canary-248-not-a-real-key-0f2a"` → `:375 if strings.Contains(reply, canary) {`
- 另一族同义尺：`cmd/wisp/instructions_200r2_test.go:73 t.Errorf("instruction body leaked onto the wire: %s", raw)`
- 静态 CI：d22scan **ban #3 plaintext-key**（`tools/d22scan/main.go:17-18` 定义；`:164-165` `secretNameWords`/`secretNameRe`；`:754-755`/`:772-773` 定罪句"identifier … holds a plaintext key literal (D22/C28: SecretStore refs only)"）——⚠ 它扫的是**标识符赋值字面量**，不覆盖运行时载荷，不能抵账那六面，但是同一意图的门。

⇒ **落点结论：本框不应按"从零造一枚泄密扫描尺"来派。** 残差只有两处没进那把尺的面清单：
1. **入向原始串**（`dispatchRaw` 的 `raw`、`panel_host_windows.go:725`）没被列为被检面；
2. **未来的事件推载荷**（不存在，随 `:47`/`:49` 一起做）。

### 2. 最少落点（若只做残差 1）
**1 枚文件、60–150 行**：`cmd/wisp/` 包内直接复用同包已有的 `secretShape248`/`hits248`（`panel_config_248_test.go:47`），把面清单加"入向 raw"。⛔ 不需要新产码，⛔ 不需要改那三枚既有 248 断言。
（残差 2 ＝ 跟 `:47` 合并做，不单列。）

### 3. 要不要新增 C17 契约面
**不要。** 一格都不碰名册。

### 4. 有没有仪器
**有，而且这格本身就是仪器**：`TestAC2CredentialSentinelAppearsInNoArtifact` 若快照/回执里出现 secret 形状会当场红（`pump.go:60` 那一段还具名记着票 248 加的两枚键 `credentialState`/`credentialKnown` 是被双向对账尺 `TestComposerContractTypesMatchFrontend` 盯着的）。⇒ 交完**有牙**；但⚠ 若新落的出向载荷不扩面清单，就又是"无牙自述"。

---

## ⑦ 两处查重的总结论（编排者特别在意的那两问）

1. **`:45` vs `35-r3`/`35-r4` 那三枚入向守卫＝同一枚物理缺陷（路由层），已归口。** 铁证＝`panel_inbound_guards_35r3_test.go:80` 那枚信封里逐字带着 `"outcome":"allow"`，`:150` 断它在 `bridge.go:132` 被拒；能力层另有 5 枚 `ErrPanelAllow` 尺＋2 枚静态尺（d22scan ban #6、`TestFrontendNeverNamesAnApprovalDecision`）。⇒ `:45` **零落点**，唯一残差是票面那句"browser-context e2e"，那一支要真窗（〔仅本机可量〕）＋要入册（＝C17 面＝人工批准，且票 194 AC#2① 明令禁 allow 从面板发起）⇒ **不可派**。
2. **`:51` 已有同类尺（票 248 那族成套，含恒真正控；票 200 那族一枚同义尺；d22scan ban #3 静态）。** ⇒ 本框只做"把入向 raw 加进面清单"＝1 枚文件 60–150 行，⛔ 不重复造尺。

---

## ⑧ 具名顶回编排者的句子（本腿逐枚复跑）

1. ⛔ **"守卫在 `cmd/wisp/panel_host_windows.go` 的 `:132/:135/:139`/:680 一带"——三枚入向守卫不在那枚文件，在 `internal/panel/bridge.go`。** 现量：`bridge.go:132`＝`if !knownComposerMethod(r.Method) {`、`:135`＝`if strings.TrimSpace(r.Source) != ComposerRequestSource {`、`:139`＝`if strings.TrimSpace(r.RequestID) == "" {`。而 `cmd/wisp/panel_host_windows.go:132-137` 逐字是一枚 MSG 结构体字段（`message uint32`／`wParam uintptr`／`lParam uintptr`／`time uint32`／`pt struct{ x, y int32 }`／`_ uint32 // x64 padding DWORD`）。
2. ⛔ **同句里 "`:680` 一带"也偏了。** `panel_host_windows.go:680` 是注释行（`:678-681` 讲的是 stub 同步注册 `_rpc[seq]`），JS 那道 typeof 门逐字在 **`:699`**：`    if (inside || typeof window.%[1]s !== "function") { return native.call(cw, message); }`。（票 `:78`/`:292` 历史上把它写成 `:680`，与此刻 HEAD 的真身差 19 行——本锚 blob 是 `1f9060df`，不是票 `:205`/`:308` 时期的 `26b5de83`。）
3. ⚠ **"不是刚补的那三枚守卫框"——盘上是 4 枚已勾，不是 3 枚。** 现量 `- [x]` ＝ `:52`（Transport agreement）、`:63`（Inbound judge）、`:75`（Behavioural yard，含 `:75` 那一整节）、`:77`（Behavioural yard 第三面／AC#8）＝ **4 勾／6 未勾／共 10 框**（与票 `:370` 自己那句"翻后 6 未勾／4 已勾"同色）。
4. ⚠ **"我已经把这条链上『传输/夹具行为尺』那一族做完了（验收框已翻）"——框翻了，票面自己还挂着两笔未清的编排者欠账，"做完了"要带上它们**：票 `:374` 逐字"★**同波控制③（把产码改坏再跑真窗）我没做**，所以『它今天能转红』没有本波凭据…**下一波碰传输必须补这一支**"；票 `:357` 逐字"编排者欠一发 `-tags winlive` 真窗＋控制组读数"。本腿⛔无开窗权限，取不到这两发，只具名转述票面原文。
5. ✅ **派单转述里对上的（不顶回，只记同色）**：6 枚未勾框的行号 `:42/:44/:45/:47/:49/:51`（尺复跑＝6 枚）；"这 6 枚是最初的功能性 AC"；"某格要求的 resync 通道"确实零命中（本腿换成词边界尺复跑＝`*.go` 0 命中）。

---

## ⑨ 六格一句话结论表（正文用；凭据在上面各节）

| 框 | 被测物 | 最少落点 | 新 C17 名 | 仪器 |
|---|---|---|---|---|
| `:42` Whitelist fuzz | **有**（拒＋`INBOUND-DISPATCH` 审计行都在产码） | 前半 1 枚测试 120–200 行 | 前半**不要**；后半（按 SPEC-08 读）要 16 枚＝**归口票 194** | 今天无专属尺⇒交完有牙 |
| `:44` Concurrency routing | **一块没写**（回执同步返回；`NewRequestID` 写了零调用者） | 甲：1 枚测试 100–160 行；乙：3–4 枚文件 300–600 行（前置＝`:47`） | **不要**（除非事件名进 `bridge.go` 常量块） | 今天 0 枚；甲版须过恒真自查 |
| `:45` Forged panel allow | **有，两层都在**，且**与 `:63` (d)-1 同一枚物理缺陷** | **0**（归口；残差＝browser e2e＝不可派） | 残差那支要＝C17 面＝人工批准，且票 194 AC#2① 明令禁 | 已有 6 枚 Go 边界尺＋2 枚静态尺 |
| `:47` Resync | 行为半格**一块没写**（推送原语非 test 调用者 0；`resync` 在 `*.go` 0 命中）；lint 半格**有** | ≥4 枚文件 200–500 行＋页面（⛔ 页面写权未放开） | 有陷阱：`panel.resync` 进 `bridge.go`＝动 `wantFullInboundRosterSize253=6` | 今天 0 枚会变色＝最易自证通过 |
| `:49` Backpressure | **部分在，且钉的方向相反**（StreamLog 有界＋三计数，规则是 truncate-not-merge；每面板推队列没写；`failedPost` 有计数零尺） | **先裁冲突**；按②走＝2–3 枚 250–450 行，前置＝`:47` | 不要 | 3 枚既有尺会因 merge 红（反咬）；`failedPost` 无尺 |
| `:51` No secret leakage | **有，成套**（票 248 那族含正控） | 残差 1 枚文件 60–150 行（复用 `hits248`） | 不要 | 尺本身就在，会红 |

---

## ⑩ 纪律自报

- `git status --porcelain -- cmd internal docs .github scripts`：起手 **0 行**（`00-anchor.md` §2）；交件时复跑仍＝**0 行**（尺＝`01-readings.md` R16/R24；本腿写面只有 `probes/35/a6/**` 与票 35 文末那一节）。⚠ 仓内**其它路径**（`.gitignore`、`probes/152`、`probes/161/r6/logs/flip-*.txt` 等）此刻由别的腿挂着 ` M`／`??`，与本腿无关，本腿一枚没 stage。
- ⛔ 零 `go build`/`go vet`/`go test`/`go run`；`go env`/`go list` **零条**（§开头）。
- ⛔ 零 `git add -A`／`git add .`；每次 commit 带显式 pathspec；⛔ `--amend`/`reset`/`rebase`/`stash`/`checkout .`/`clean`；⛔ push。
- ⛔ 未开任何真窗（WebView2/浏览器）。
- ⛔ 未新建 `.sh`/`.ps1`/`.txt`/`.out`：本腿写面全为 `.md`。
- ⛔ 未在仓库目录内建 worktree／checkout；未在仓外建工作件（所有尺都在短命令里现跑，无须 `D:/tmp/wisp35a6/`）。
- 票面 AC 框：⛔ 一字未勾、未动、措辞未改；只在文末追加一节。追加前后同数＝**6 未勾／4 已勾**、未勾行号逐枚仍是 `42 44 45 47 49 51`，且 `git diff --numstat` 读出 **29 增／0 删**＝纯追加（尺与 rc＝`01-readings.md` R3/R4/R24；"已勾是 4 枚不是 3 枚"那条顶回见 §⑧ 第 3 条）。⚠ 落笔过程中我曾在一枚**原有行**里吃掉一个空格（`它**不覆盖**` 与 `` `frontend/**` `` 之间），已按原文补回，终态 0 删——失误与还原凭据都在 R24 里具名记着。
- 本件行数＝正文唯一交付凭据件；票面那一节是它的摘要。
