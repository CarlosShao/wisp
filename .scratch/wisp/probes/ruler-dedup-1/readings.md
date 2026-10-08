# ruler-dedup-1 读数件（命令＋输出逐条留档，可重跑）

- 时刻：`2026-10-08 13:54:42+0800`
- **CWD＝仓库根 `D:\work\workspace\projects plans\Wisp`**，下面每条命令都在这枚 CWD 跑的；相对路径按仓库根解析。
- 零 `go` 命令。零改动、零 commit、零 push。

## §0 我引用的每一枚文件都单枚复量过"工作树＝HEAD"

```sh
git status --porcelain -- internal/agent/approval/ticket259_panel_capability_rulers_test.go \
  internal/agent/approval/ticket259_denial_rulers_test.go internal/agent/approval/ticket242_panelface_test.go \
  internal/agent/approval/ticket242_binding_test.go internal/panel/inbound_roster_253_test.go \
  internal/panel/composer_test.go internal/panel/git_test.go internal/panel/bridge.go \
  internal/panel/composer_dispatch.go cmd/wisp/panel_geometry_255_test.go \
  cmd/wisp/panel_locked_naming_33r11_windows_test.go cmd/wisp/panel_transport_35r1_test.go \
  cmd/wisp/panel_inbound_guards_35r3_test.go cmd/wisp/panel_host_gate_test.go \
  cmd/wisp/subagent_selfapproval_197_test.go internal/agent/approval/ui.go
```
输出＝**空**（整棵工作树的 porcelain 是 759 行＝别人的改动，`git status --porcelain | wc -l`＝**759**，我一枚没碰）。
⇒ 上面这 16 枚文件的盘上内容＝`git show HEAD:<path>` 的内容，读盘不会读到被突变污染的那一份。

## §1 锚点

```
$ git rev-parse --abbrev-ref HEAD ; git rev-parse --short HEAD
dev
53d73db4
```

## §2 三张票的行数与框数（带 `[[:space:]]*` 的那把）

```sh
for t in 242-the-grant-binding-layer-and-the-panel-facing-read-surface-both-have-zero-rulers.md \
         253-panel-inbound-three-ruler-holes.md 259-grant-binding-identity-ruler-holes.md; do
  p=".scratch/wisp/issues/$t"; printf "%s open=%s done=%s\n" "$t" \
    "$(grep -cE '^[[:space:]]*- \[ \]' "$p")" "$(grep -cE '^[[:space:]]*- \[x\]' "$p")"; done
```
```
242-the-grant-binding-layer-and-the-panel-facing-read-surface-both-have-zero-rulers.md open=3 done=0
253-panel-inbound-three-ruler-holes.md open=4 done=0
259-grant-binding-identity-ruler-holes.md open=5 done=1
```
行数：`wc -l` 分别＝**63／53／45**。未勾框所在行：242 `:18`/`:21`/`:28`；253 `:16`/`:17`/`:18`/`:19`；259 `:20`/`:21`/`:22`/`:23`/`:24`（已勾＝`:19`）。

## §3 票 259 `AC#3` 原文（那一枚格，逐字）

```
.scratch/wisp/issues/259-grant-binding-identity-ruler-holes.md:22
- [ ] **AC#3 三枚能力尺（补票 242 AC#2 不成立那一格，随本票一起做，⛔ 不另开一票）**：① `PanelAPI` 的**方法名封闭集**尺（新增一枚能改变卡状态的方法名 ⇒ 指名红；M-E 那一形必须被它抓到）；② `PanelItem` 的**字段封名单**尺要**按能力**判（能改变状态／能携带凭据形状的字段名怎么改都算——M-D2 那一形必须被抓到，⛔ 只认字符串名的尺不许当凭据）；③ 语义侧尺（出向读面拿到的东西**不可**回填成答复）。⛔ 这三枚只加尺，**不许为此解冻乙形**（`242-v1` 已判"缺的零件可在甲内补"）。
```

## §4 票 242 的两处（原判据＋10-08 射程注）

```
242-…-zero-rulers.md:21（AC#2 原判据，逐字）
- [ ] **AC#2 出向读面不得带 grant**：反射扫 `approval.PanelItem` 的字段名枚数与名字 ⇒ 出现 `Grant`／可答复类字段**就红**；再加一发**能力侧**判据（不是词面）：把 `PanelAPI` 的出向读面扩到"能改变卡的状态"那一形也要红。
```
```
242-…-zero-rulers.md:18 内嵌射程注（⚠ 不是判据正文，是编排者 10-08 就地追加）
〔10-08 15:6x 编排者**就地加射程注**、⛔ 未翻勾、原判据一字未改：… ⇒ 那一半**归口票 259**：它 `:12` 逐字记着"对外只有一个 `ErrBadGrant`、三种原因并成一句"，`internal/agent/approval/approval.go:477` 逐字写着 **"ticket 259 does not move it"**，而**选边（ⓐ具名降级／ⓑ补牙）是票 259 `AC#1` 明写了归编排者的那一格**。⛔ 任何腿不许为了闭合本格去自己改合并文案。…〕
```
注的出处 commit（`git log --oneline -3 -- <242 票>` 首行）：
```
cade79b3 242 射程注（编排者自写、⛔ 未翻勾、原判据一字未改）：票 242 AC#1 就地加一行带时刻的射程收窄注…
```
```
$ git merge-base --is-ancestor cade79b3 HEAD && echo "cade79b3 IS ancestor of HEAD"
cade79b3 IS ancestor of HEAD
```
242:56（把 `AC#2` 推给 259 `AC#3` 的那句，逐字）：
```
归**票 259 AC#3** 三枚能力尺（`PanelAPI` 方法名封闭集／`PanelItem` 按能力判的字段封名单／语义侧"读面不可回填成答复"），⛔ 不为此解冻乙形（腿判"缺的零件可在甲内补"）。
```

## §5 票 253 的"只能落 cmd/wisp"那几格（逐字）

```
253-panel-inbound-three-ruler-holes.md:16
- [ ] **AC#1 可达性那枚钉改成问能力**：判据不许再认某个 JS 调用的**词面**，要问**"页面发起的一次调用，Go 侧有没有收到并回执"**…⚠ 正控必配：造一发"只 `postMessage`、不 `wispDispatch`"的假腿 ⇒ 指名用例**必须红**。
253-panel-inbound-three-ruler-holes.md:42（§8 里的窄义改判）
⇒ 真洞只在窄义：**没有任何词面/AST 尺断言产码里 `w.Bind("wispDispatch")` 存在**。所以 AC#1＝"把可达性裁决从词面尺迁到能力形＋补一发只 `postMessage` 不 `wispDispatch` 的正控"，⛔ 不是从零造一台仪器。
253-panel-inbound-three-ruler-holes.md:43／:53（落点＝cmd/wisp）
**排程更新**：AC#1／AC#3 写面＝`cmd/wisp`（`255-r2` 在飞）⇒ 按住到它交件；AC#2 写面＝`internal/panel` 新增测试…
★**AC#1／AC#3 的排程更新**：`internal/panel` 此刻空了，但这两格写面＝**`cmd/wisp`**（`255-r2` 在飞）⇒ 照旧按住；
```

## §6 259 两枚尺文件的在册用例名（grep 现跑）

```sh
grep -n "^func Test" internal/agent/approval/ticket259_panel_capability_rulers_test.go
grep -n "^func Test" internal/agent/approval/ticket259_denial_rulers_test.go
```
```
66:func TestTicket259R1PanelAPIMethodSetIsClosed(t *testing.T) {
93:func TestTicket259R1PanelCarrierCarriesNoAnswerVerb(t *testing.T) {
169:func TestTicket259R1PanelItemFieldCapabilityFence(t *testing.T) {
184:func TestTicket259R1PanelItemCarriesNoCredentialValue(t *testing.T) {
216:func TestTicket259R1PanelItemProjectionIgnoresGrantState(t *testing.T) {
283:func TestTicket259R1PanelReadFaceCannotBeLaunderedIntoAnAnswer(t *testing.T) {
115:func TestTicket259R1DenialNamesMissingNonce(t *testing.T) {
134:func TestTicket259R1DenialNamesSpentOrNeverLiveNonce(t *testing.T) {
165:func TestTicket259R1DenialNamesMisbound(t *testing.T) {
192:func TestTicket259R1DenialNamesCardNotPending(t *testing.T) {
226:func TestTicket259R1FourDenialLabelsArePairwiseDistinct(t *testing.T) {
249:func TestTicket259R1OutwardAnswerStaysMergedAcrossDenials(t *testing.T) {
```
行数＝**349／288**。交件出处（两枚同一次）：
```
$ git log --oneline -1 -- internal/agent/approval/ticket259_panel_capability_rulers_test.go
b6b1d6a4 票 259 腿 259-r1 第 2 步：AC#2 拒因可指名＋AC#3 三枚能力尺的码与新尺
```
在册名册与关键红句（同文件行号）：`:44` `var panelAPIAllowedMethods = []string{"Reject", "Head", "View"}`；`:73` 红句含 `the panel surface is closed at Reject/Head/View`；`:81-84` 计数等式；`:108`/`:113`/`:118`/`:123` 四枚形状断言；`:138-164` Kind 判定表；`:199`/`:202`/`:205` 值级红句；`:229` 不变量红句；`:294-296` 取料为零即 `t.Fatal`；`:300`/`:307`/`:316`/`:323` 回填试验的四个点。

## §7 同一枚石头今天已有几把尺（PanelItem／PanelAPI 面）

```sh
grep -rn "PanelAPI\|PanelItem" --include=*_test.go internal cmd | grep -v "ticket259_panel_capability\|ticket242_panelface"
```
```
internal/agent/approval/fakes_test.go:104:	views    []approval.PanelItem
internal/agent/approval/queue_test.go:136:	// The PanelAPI itself has no way to ask for an allow: there is no method
internal/agent/approval/ticket146_liveapprovals_backing_test.go:338:	// PanelItem already copies Paths (queue.go's viewLocked), so this is the
internal/agent/approval/ticket146_liveapprovals_backing_test.go:346:		t.Errorf("AC#2 RED: PanelItem.Paths=%v，期望队列投影响原样（面板卡片的越界提示读这一枚）", item.Paths)
internal/agent/approval/ticket146_liveapprovals_backing_test.go:349:		t.Errorf("AC#2 RED: PanelItem 非引用字段也被牵动：tool=%q reason=%q", item.Tool, item.Reason)
cmd/wisp/subagent_selfapproval_197_test.go:696://	M1 internal/agent/approval/ui.go: add Grant to PanelItem -> (4) red on the
```
（注：`queue_test.go:136` 是一段行为用例里的注释，`:140` 才调 `Panel().Reject`，不是封名单尺。）

被测物当前形状（`internal/agent/approval/ui.go`）：
```
ui.go:50-57   type PanelItem struct{ CorrelationID string; Tool string; Level string; Reason string; Paths []string; Depth int }  // 6 枚字段
ui.go:167-171 type PanelAPI interface{ Reject(correlationID, reason string) error; Head() (PanelItem, bool); View(correlationID string) (PanelItem, bool) }  // 3 枚方法
ui.go:143      type NativeAPI interface{ … }（AllowSession 只在这枚上）
ui.go:155-156  注释逐字："It exists on this interface and on no other surface: PanelAPI has no Allow method at all"
ui.go:158  	AllowSession(ctx context.Context, correlationID, grant string) error
```
票 242 那把词面尺（仍在盘；259 文件头 `:27-29` 逐字声明 "the two 242 rulers stay exactly as they were"）：
```
ticket242_panelface_test.go:21-28 panelItemReadFace = 6 名（CorrelationID/Tool/Level/Reason/Paths/Depth）
:30 TestTicket242PanelItemReadFaceIsFullyDeclared（双向）
:53 TestTicket242PanelItemStaysGrantFree（:59 禁词表 = grant/allow/approve/nonce/token，名字与类型串各过一遍 Contains）
```

## §8 "按形状/能力判"的反射尺名册（不读盘，overlay 可见）

```sh
grep -rln "NumMethod()" --include=*_test.go cmd internal
```
```
cmd/wisp/subagent_selfapproval_197_test.go
internal/agent/approval/ticket220_l1_window_read_test.go
internal/agent/approval/ticket259_panel_capability_rulers_test.go
internal/panel/composer_test.go
internal/tools/grant_test.go
internal/tools/ticket224_setter_scope_test.go
internal/tools/ticket90_test.go
```
各自用例与枚数等式（`awk` 把每处 `NumMethod()` 归到它所在的 `func Test`）：
```
internal/agent/approval/ticket220_l1_window_read_test.go :343 typ.NumMethod() != 0 ∈ :338 TestL1WindowIsPlainDataWithNoRouteToAnAnswer（:361-362 另扫 *Gate 方法集）
internal/panel/composer_test.go                          :172-173 ∈ :171 TestModeViewHasNoWriteSurface
internal/tools/grant_test.go                             :521-531 ∈ :508 TestTicket224GrantSeamCarriesNoCallerSuppliedAuthority（==1）
internal/tools/ticket224_setter_scope_test.go            :81 / :102-105 ∈ :41 TestTicket224BridgeGrantChannelIsOneSeamWithNoSetter
internal/tools/ticket90_test.go                          :433-435 ∈ :431 TestTicket90ModeCarriesNoAllowAuthority（红句 "want exactly 1 (read-only)"）
cmd/wisp/subagent_selfapproval_197_test.go               :663-664 ∈ :557 Test197NoAllowDoorIsReachableFromASubagentsAssembly
```
197 那枚的两张名册（逐字）：
```
:541 var allowDoorMethods197 = []string{"Allow", "Native", "DecideFromNative", "DecideFromPanel", "GrantNonce"}
:547-554 toolAssembledRoots197 = {tools.SubagentDeps, tools.TaskDeps, tools.Options, agent.Options}
```

## §9 读盘型尺：go/parser 名册＋逐枚射程

```sh
grep -rln "go/parser" cmd internal
```
22 条命中；去掉 `internal/config/config.test`（`git ls-files internal/config/config.test`＝**空**＝未跟踪二进制）与 `tools/d22scan/main.go`（CI 扫描器，非尺）⇒ **20 枚测试文件**。逐枚读取形状（`grep -n "os\.ReadDir(\|os\.ReadFile(\|parser\.ParseFile\|filepath\.Glob\|filepath\.WalkDir\|runtime\.Caller\|go:build" <file>` 现跑）：

```
cmd/wisp/panel_geometry_255_test.go                 :1 //go:build windows；读盘点 :208/:316/:343/:376；取器 :449→:451 os.ReadFile；目录 :441 runtime.Caller(0)；func Test＝5
cmd/wisp/panel_geometry_255r6_range_windows_test.go :1 windows；:42 os.ReadDir(dir)；:78 os.ReadFile；:82 ParseFile；func Test＝0
cmd/wisp/panel_locked_naming_33r11_windows_test.go  :1 windows；:176 os.ReadDir(".")；:211/:443 ParseFile(带 src)；func Test＝5；:186 零枚即 t.Fatal；文件头 :38-41 声明正控从本文件自带源码串解析
cmd/wisp/panel_host_gate_test.go                    :33 ParseFile(单枚，ImportsOnly)；:210 os.ReadDir(".")＋:221；func Test＝4；:182 有 t.Skipf 支
cmd/wisp/panel_inbound_33_test.go                   :207/:285 filepath.Glob("*.go")；:216/:296 ParseFile；func Test＝5
cmd/wisp/leg_dispatch_gate_133_test.go              :365 os.ReadDir(dir)；:375 ParseFile；func Test＝1
cmd/wisp/leg_sink_gate_131_test.go                  :561 os.ReadDir(dir)；:581 ParseFile；:593 另 os.ReadFile 原文；func Test＝1
cmd/wisp/panel_assets_143_test.go                   :267 ParseFile("panel_assets.go")；func Test＝6
cmd/wisp/panel_reshow_255r1_windows_test.go         :1 windows；:432 ParseFile(hostFile255r1)；:427 注释点名 runtime.Caller；func Test＝6
cmd/wisp/resident_approval_risk_256_windows_test.go :1 windows；:78 runtime.Caller；:340/:386/:521 三处 ParseFile；func Test＝5
cmd/wisp/resident_approval_risk_268_windows_test.go :1 windows；:227 单枚；func Test＝3
cmd/wisp/resident_ball_228_test.go                  :66 os.ReadDir(".")；:77 ParseFile；:154 扫 //go:build 注释；func Test＝2
cmd/wisp/resident_grant_writer_265_windows_test.go  :1 windows；:630 runtime.Caller；:540 os.ReadDir(dir)；:468/:644 ParseFile；func Test＝9
cmd/wisp/resident_hotkey_258_test.go                :28 单枚（resident_windows.go）；func Test＝2
cmd/wisp/subagent_carrier_197_test.go               :461 WalkDir(root+sub)；:466 ParseFile；func Test＝3
internal/agent/spill_path_invariant_test.go         :436 ParseFile("spill.go")；:142 os.ReadDir(运行时目录)；func Test＝4
internal/ball/tokens_table_test.go                  :101/:114 runtime.Caller；:373/:1231 ParseFile；func Test＝5
internal/memory/artifacts_path_invariant_test.go    :771 os.ReadDir(".")；:760 ParseFile；func Test＝5
internal/panel/composer_dispatch_test.go            :604-606 SkipDir 名单；:615 os.ReadFile；:618 ParseFile ⇒ 全仓生产 .go；func Test＝12
internal/panel/inbound_roster_253_test.go           :136/:252/:543/:600/:698/:811 os.ReadFile＋:145/:257/:548 ParseFile；root＝:420 panelRepoRoot；目标 :422 bridge.go／:423 composer_dispatch.go；分母⑤ :484 productionPanelFiles(t, dir)＋:511-513 循环；TempDir 种形 :632/:658/:712/:738/:771；func Test＝5
internal/panel/l2_grant_boundary_test.go            ⚠ 冻结件；:738 走盘；:2032 解析测试自带 bridgeSrc 串（这一支 overlay 可见）；func Test＝7
internal/projctx/projctx_test.go                    :263 os.ReadDir(dir)；:273 ParseFile；:283 解析自带串；func Test＝10
internal/tools/paths_twocontainments_252_r2_test.go :194 os.ReadDir(".")；:206 ParseFile；func Test＝2
internal/tools/task_state_188_test.go               :254 WalkDir(root)；:261 ParseFile；:248-249 注释逐字（零分母的扫描不许读成 clean）；func Test＝4
```
`panelRepoRoot` 的基准（cwd 上溯找 `go.mod`，定义在**冻结件**里）：
```
internal/panel/tokens_fourway_test.go:57 func panelRepoRoot(t *testing.T) string {
:59   dir, err := os.Getwd()
:64   if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil { return dir }
```
overlay 失明这条定式的原话（`git log` 里 `3ee9b3ff` 那条编排者口径，逐字）：
```
读盘与走 AST 的测试对 overlay 结构性失明，活例 TestSubagentStreamKeyHasOneMintSite 与 leg_sink_gate_131 那族
```

## §10 两把"词面钉"与转发那支（253 现量的复认）

```
internal/panel/composer_test.go:381 var hostBridgeCallRe = regexp.MustCompile(`\.\s*postMessage\s*\(`)
internal/panel/composer_test.go:394 var routeLiteralRe = regexp.MustCompile(`"panel\.[A-Za-z.]+"`)
internal/panel/composer_test.go:409-419 composerRouteLiterals()＝5 名（4 枚 Method* 常量＋`"panel.approval.request"`）
internal/panel/git_test.go:394 panelMethodRe = regexp.MustCompile(`"(panel\.[a-z0-9_.-]+)"`)
internal/panel/git_test.go:385-390 whitelistMethodsFromSource 只点名单文件 bridge.go，红句自述 "the four methods this ticket may not extend"
```
⇒ 票 253 现量 2 那句"不带 `panel.` 前缀就掉在两把尺之外"在**形状层复认为真**（两把正则都把 `panel.` 写死在模式里）。

门那一侧（253 `AC#1` 的目标物）：
```
cmd/wisp/panel_host_windows.go:80   panelDispatchBinding = "wispDispatch"
cmd/wisp/panel_host_windows.go:801  func (m *PanelManager) installPanelTransport(w pageTransport, ctx context.Context) error {
cmd/wisp/panel_host_windows.go:802  	if err := w.Bind(panelDispatchBinding, func(raw string) string {
```
```sh
grep -rn "panelDispatchBinding" --include=*.go cmd internal   # 产码命中只有 :80/:792/:802 三处；其余全在测试假件里
grep -rn "installPanelTransport" --include=*.go cmd           # 产码唯一调用点 :408；测试直调点 35r1:203／35r3:97
```
```
cmd/wisp/panel_transport_35r1_test.go:1 //go:build windows；:195 TestPagePostMessageEnvelopeReachesDispatchRawViaTransport
:203 mgr.installPanelTransport(fc, ctx)；:210 红句 "AC#6 RED: installPanelTransport registered no page->wispDispatch forwarding hook"
:217 if fc.doorRounds != 1 {…}；:176-187 fc.msgcb 的未绑定方法死亡支（rpcDeadSlot++）
cmd/wisp/panel_inbound_guards_35r3_test.go:97 mgr.installPanelTransport(bf, context.Background())；:65-66 把 Go 侧两句拒收文案逐字钉住
cmd/wisp/panel_resident_windows_test.go:314 TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe；:812 TestAC14AwaitedBindingReplyReachesThePage
```
Go 侧具名拒收（253 `AC#3` 的相邻件）：
```
internal/panel/bridge.go:131-133 if !knownComposerMethod(r.Method) { return r, fmt.Errorf("%w: 方法 %q 不是面板 composer 通路的能力入口", ErrComposerRequest, r.Method) }
internal/panel/bridge.go:136-137 来源那支（按伪造/串台拒绝）；:140-141 requestId 那支
internal/panel/bridge.go:145-147 knownComposerMethod 的 case＝MethodModeRequest, MethodWorkspaceRequest, MethodAttachmentAdd, MethodMessageSend, MethodConfigGet, MethodConfigSet（6 名）
internal/panel/composer_dispatch.go:208 return "", d.rosterMismatch(req)；:215-222 函数体（d.record(req, err) 落审计）
internal/panel/inbound_roster_253_test.go:494 把"default 必调 rosterMismatch"当分母④钉住
```

## §11 载具那一格（259 `AC#4`／242:60 指的同一处）

```
cmd/wisp/subagent_selfapproval_197_test.go:108-109（逐字）
			out, e := rt.bridge.Execute(ctx, agent.ToolRequest{
				TaskID: taskID, CorrelationID: taskID,
```
⇒ 票面那句"两枚同时活着、correlation id 不同"的载具前提**在盘上仍不成立**。

## §12 README 尺的归口句＋"别重复"定式

```
.scratch/wisp/issues/README.md:70 ⚠ 两处已知的不彻底，都归**票 262**，别把这行当"已经装好牙"：① 这条帽**只有词面**——今天没有任何仪器
.scratch/wisp/issues/README.md:79 剩下的缺口只有一个，而且是**排程**缺口不是仪器缺口：**没任何东西要求"开票者在提第一笔之前跑它一发"**。
.scratch/wisp/issues/README.md:80 ⛔ 不许为此新造脚本——那会与票 262 重复，机主原话"不要重复劳动哈"）。
.scratch/wisp/issues/README.md:83 归 `262-v1`（它此前被按住的原因＝会种超长文件名、把别人正在跑的门禁读数洗脏）。
```
`grep -n "都归票" .scratch/wisp/issues/README.md`＝**零命中**（原文里"都归"与"票 262"之间夹了加粗标记，逐字形状＝上面 `:70` 那行）。⇒ 派单给的起手串"都归票 X"在本仓的实际写法是"都归**票 262**"，用裸串会整枚漏掉（同形于票 62 那次数框漏计）。

## §13 收尾时刻

```
$ date
Thu Oct  8 13:50:49 CST 2026
2026-10-08 13:54:42+0800   （本件写入时刻那次）
```
