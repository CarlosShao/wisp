# 35-a6 raw readings — 逐尺现量＋每尺自带 rc

跑法＝`out=$(<ruler>); rc=$?` 先取**尺本体**的退码，再截断显示（⛔ 不许把 `head` 的 0 当尺的 0）。
工作目录＝仓库根；⛔ 未在仓内落任何临时件（读数直接进本 `.md`）。

⚠ **HEAD 漂移具名**：起手锚 `00-anchor.md` 记的是 `7a367a08292ef94908d86656608974e81c42f0d4`；本件复跑时 `git log -1 --format=%H` 已前进到另一枚（见 R1）。
共享工作树随时前进，属预期；本腿所有行号引用锚在**锚时刻 blob**（`bridge.go bebe8e70`／`panel_host_windows.go 1f9060df`，见 `00-anchor.md` §3），下面逐尺复量这些 blob 是否仍为真身（R14–R16）。

---

## R1 — HEAD 与分支

```
$ git log -1 --format=%H          rc=0
$ git rev-parse --abbrev-ref HEAD rc=0
```

## R2 — 脏面（四棵写面）

```
$ git status --porcelain -- cmd internal docs .github scripts   rc=0
（输出 0 行）
```

## R3 — 票面框数尺（未勾）

```
$ grep -cE '^[[:space:]]*- \[ \]' .scratch/wisp/issues/35-panel-bridge-c17.md   rc=0
6
$ grep -nE '^[[:space:]]*- \[ \]' .scratch/wisp/issues/35-panel-bridge-c17.md    rc=0
42:44:45:47:49:51 行（逐字行号见 verdict §开头）
```

## R4 — 票面框数尺（已勾）

```
$ grep -cE '^[[:space:]]*- \[x\]' .scratch/wisp/issues/35-panel-bridge-c17.md    rc=0
4
$ grep -nE '^[[:space:]]*- \[x\]' ...                                            rc=0
52 / 63 / 75 / 77
```

## R5 — `panel.resync` 在 Go 产码里（词边界尺，⛔ 不用子串尺）

```
$ grep -rniE '\bresync\b' --include=*.go cmd internal tools   rc=1（1＝无命中）
命中 0 行
```
> 尺意：票 `:95` 记过一把子串尺把 `fixtureSy**reSync**Root` 当成命中 ⇒ 本腿改用 `\bresync\b`。

## R6 — `resync` 全跟踪树（含页面）

```
$ git grep -nE '\bresync\b' HEAD -- '*.go' '*.ts' '*.tsx'   rc=0
命中 1 行 ＝ HEAD:frontend/src/lib/panel-views.ts:90
```
逐字（整行读，未截断）：
```
 * PanelSnapshot yet - ticket 35's pump owns it, and `panel.resync` (the push
```

## R7 — 出向原语（**非 test** 调用者）

```
$ grep -rn --include=*.go '\.Eval(' cmd internal | grep -v '_test.go'   rc=1
命中 0 行（非 test）
$ grep -rn --include=*.go '\.Eval(' cmd internal                        rc=0
命中 5 行，逐枚核：全部落在 _test.go（含注释 3 行 ＋ 测试体调用 2 行：
  cmd/wisp/panel_resident_windows_test.go:175  /  cmd/wisp/panel_transport_live_35v2_windows_test.go:159）
```

## R8 — 出向 COM 拼写（票 `:119` 那把"拼错名字尺"的正面版）

```
$ grep -rn --include=*.go -E 'PostWebMessage|EvaluateScript|CreateWebMessageAsJson' cmd internal | grep -v '_test.go'   rc=1
命中 0 行
```

## R9 — 泵字节的读者（`lastPanelSnapshot()`）

```
$ grep -rn --include=*.go 'lastPanelSnapshot()' cmd | grep -v '_test.go'   rc=0
唯一命中＝cmd/wisp/panel_pump.go:383:func (rt *agentRuntime) lastPanelSnapshot() (panel.Snapshot, []byte, bool) {
⇒ 那是**定义**不是调用 ⇒ 非 test 调用者＝0（"写了但零调用者"那一族）
$ grep -rn --include=*.go 'lastPanelSnapshot()' cmd | wc -l                 rc=0
12（1 定义 ＋ 11 处 _test.go 调用）
```

## R10 — `pageTransport` 接口面（出向方法有没有位置）

```
$ grep -n -A4 'type pageTransport interface' cmd/wisp/panel_host_windows.go   rc=0
652:type pageTransport interface {
653-	Bind(name string, f interface{}) error
654-	Init(js string)
655-}
⇒ 无出向方法
```

## R11 — `failedPost`（背压计数）有没有尺

```
$ grep -rn "failedPost" --include=*_test.go cmd    rc=1
命中 0 行（计数在产码 panel_resident_windows.go:130/:369，零测试引用）
$ grep -n "make(chan func" cmd/wisp/panel_resident_windows.go   rc=0
148:		tasks:    make(chan func(), 16),
```

## R12 — 名册枚数（真符号面）

```
$ grep -n 'Method[A-Z]' internal/panel/bridge.go   rc=0
:42 MethodModeRequest / :43 MethodWorkspaceRequest / :44 MethodAttachmentAdd / :45 MethodMessageSend / :66 MethodConfigGet / :67 MethodConfigSet
$ sed -n '148p' internal/panel/bridge.go           rc=0
	case MethodModeRequest, MethodWorkspaceRequest, MethodAttachmentAdd, MethodMessageSend, MethodConfigGet, MethodConfigSet:
⇒ 六枚，全在一行（case 表在 :148）
```

## R13 — 名册枚数被哪些 pin 钉住（新增名会当场红谁）

```
$ grep -n 'wantFullInboundRosterSize253' internal/panel/inbound_roster_253_test.go   rc=0
55-56: 注释 "…are the two / numbers a diff has to move on purpose."
59:	wantFullInboundRosterSize253 = 6
445:	if len(wantFullInboundRoster253) != wantFullInboundRosterSize253 {
787:	if len(miss) != wantFullInboundRosterSize253 {
$ grep -n 'wantNonPanelPrefixedInbound253 = ' 同件                                     rc=0
64:	wantNonPanelPrefixedInbound253 = 2
$ grep -n 'whitelistMethodsFromSource' internal/panel/git_test.go                     rc=0
387（断言）/ 409（提取器定义）/ 516（台件自身 5 枚）/ 519（真身 len(real) != 4）
$ grep -n 'func TestTheRendererHoldsExactlyOneDoorToTheHost\|func TestPlantedRendererDoorShapesGoRed' internal/panel/composer_test.go rc=0
502 / 533
```

## R14 — 三枚入向守卫的**真文件**（顶回派单第 1 条的凭据）

```
$ grep -n 'if !knownComposerMethod' internal/panel/bridge.go        rc=0  → 132
$ grep -n 'if strings.TrimSpace(r.Source)' internal/panel/bridge.go rc=0  → 135
$ grep -n 'if strings.TrimSpace(r.RequestID)' ...                   rc=0  → 139
$ sed -n '132,137p' cmd/wisp/panel_host_windows.go                  rc=0
	message uint32
	wParam  uintptr
	lParam  uintptr
	time    uint32
	pt      struct{ x, y int32 }
	_       uint32 // x64 padding DWORD
⇒ 派单说这三枚守卫在 panel_host_windows.go 的 :132/:135/:139 ＝ 错（那是 MSG 结构体字段）
```

## R15 — JS `typeof` 门的真行号（顶回派单第 2 条）

```
$ grep -n 'typeof' cmd/wisp/panel_host_windows.go   rc=0
699:    if (inside || typeof window.%[1]s !== "function") { return native.call(cw, message); }
$ sed -n '680p' cmd/wisp/panel_host_windows.go      rc=0
（该行是注释，不是门）
```

## R16 — 禁区 blob 是否仍是锚时刻真身（本腿零写入的自证）

```
$ git rev-parse HEAD:internal/panel/bridge.go        rc=0  → bebe8e70…（＝锚值，未变）
$ git rev-parse HEAD:cmd/wisp/panel_host_windows.go  rc=0  → 1f9060df…（＝锚值，未变）
$ git status --porcelain -- internal/panel/bridge.go cmd/wisp/panel_host_windows.go  rc=0  → 0 行
```

## R17 — `:45` 归口的铁证（那枚信封带 allow）

```
$ grep -n 'envelope35r3UnregisteredName' cmd/wisp/panel_inbound_guards_35r3_test.go  rc=0
80:	envelope35r3UnregisteredName = `{"method":"panel.approval.request","correlationId":"pc-35r3-approval-1","outcome":"allow"}`
$ grep -n 'func TestInboundRosterGuardRefusesUnregisteredNameOnPageEdge' 同件          rc=0
150
```

## R18 — `:45` 能力层的生产链（不是零调用者）

```
$ grep -n 'func (g \*Gate) DecideFromPanel' internal/agent/approval/gate.go        rc=0 → 736
$ grep -n 'ErrPanelAllow, "面板来源' internal/agent/approval/gate.go                rc=0 → 746
$ grep -n 'func (r \*Replies) PanelAllow' internal/agent/approval/replies.go        rc=0 → 431
$ grep -n 'g.DecideFromPanel(ctx, Request{' internal/agent/approval/replies.go      rc=0 → 437
$ grep -n 's.live.h.PanelAllow' cmd/wisp/approval_reply.go                          rc=0 → 331
$ grep -rn "NewRequestID" --include=*.go . | grep -v '^\./\.scratch'                rc=0
   → 只剩 bridge.go:80（定义）＋ bridge_test.go:48 ⇒ 非 test 调用者 0
```

## R19 — `:49` 的反向钉（AC 要 merge、盘上钉 truncate）

```
$ grep -n 'Overflow TRUNCATES and never merges' internal/panel/pump.go   rc=0 → 397
$ grep -n 'func TestTheStreamLogTruncatesInsteadOfMerging' internal/panel/pump_test.go rc=0 → 165
$ grep -n 'func TestOverflowTruncatesEachStreamsOwnHeadAndTail\|func TestHardCeilingDropsWholeKeysAndNamesThem' internal/panel/subagent_stream_197_test.go rc=0 → 209 / 275
$ grep -n 'func (s \*StreamLog) Truncated\|func (s \*StreamLog) ElidedRunes\|func (s \*StreamLog) DroppedKeys' internal/panel/pump.go rc=0 → 632 / 644 / 657
$ sed -n '397,399p;655,656p' internal/panel/pump.go   rc=0（397＝"Overflow TRUNCATES and never merges…"；655-656＝"…never merged into a row that is / still on screen"）
```

## R20 — `:51` 查重（按"这句话在说什么"扫，⛔ 不只搜符号名）

```
$ grep -rniE 'redact|leak|scrub|脱敏|泄露|泄漏' --include=*.go cmd internal tools | wc -l   rc=0 → 30+ 枚文件命中
其中打"载荷不泄密"这一意图的具名尺：
  cmd/wisp/panel_config_248_test.go:154 / :218 / :329 / :510
  internal/panel/config_route_248_test.go:349 / :375
  cmd/wisp/instructions_200r2_test.go:73
  tools/d22scan/main.go:17-18（ban #3 plaintext-key）/ :164-165 / :754-755
```

## R21 — `:42` 后半格的名册出处（SPEC-08 那张表是**提案**）

```
$ sed -n '156p' docs/specs/SPEC-08-ui-ball-panel.md   rc=0
### 5.2 C17 PanelBridge 方法白名单【SPEC 提案，S5 定稿走契约批准】
$ sed -n '163,174p' 同件                             rc=0 → 18 枚名＋事件推送那一行
$ grep -nE '^- \[ \]' .scratch/wisp/issues/194-*.md | wc -l   rc=0 → 7（票 194 全未勾）
```

## R22 — 页面有没有耳朵（`:47` 页面侧前置）

```
$ grep -rnE 'addEventListener|onmessage' frontend/src/lib/panel.ts frontend/src/App.tsx frontend/src/main.tsx   rc=0
命中仅 frontend/src/main.tsx:50 的 close.addEventListener("click", ...) ⇒ message 监听 0 处
$ grep -n 'postMessage(message: string): void' frontend/src/lib/panel.ts   rc=0 → 141（发送声明，非接收）
```

## R23 — 本腿全程的 `go` 命令

```
零条。⛔ 未跑 go build / go vet / go test / go run；go env、go list 也一条没跑。
```

## R24 — 票面追加的"只追加"自证（追加前/后各跑一次）

```
追加前（落笔前现量）：
$ grep -cE '^[[:space:]]*- \[ \]' .scratch/wisp/issues/35-panel-bridge-c17.md   rc=0 → 6
$ grep -cE '^[[:space:]]*- \[x\]'  同件                                          rc=0 → 4

追加后：
$ git diff --numstat -- .scratch/wisp/issues/35-panel-bridge-c17.md             rc=0
29	0	.scratch/wisp/issues/35-panel-bridge-c17.md      ⇒ 29 增／**0 删**＝纯追加
$ git diff -- .scratch/wisp/issues/35-panel-bridge-c17.md | grep -c '^-[^-]'    rc=1（1＝grep 无命中）
0                                ⇒ 被删/被改的原有行数＝0（⚠ 中途我曾在一枚原有行里吃掉一个空格，已按原文补回，故终态 0 删——具名记这一起居失误与还原凭据＝本尺的 0）
$ grep -cE '^[[:space:]]*- \[ \]' 同件                                            rc=0 → 6（**前后同数**）
$ grep -cE '^[[:space:]]*- \[x\]'  同件                                           rc=0 → 4（**前后同数**）
$ grep -nE '^[[:space:]]*- \[ \]'  同件 | cut -d: -f1                             rc=0 → 42 44 45 47 49 51（行号与锚时刻逐枚同）
$ wc -l 同件                                                                       rc=0 → 406（锚时刻 377 ＋ 29）
```

## R25 — 交件本身非 0 字节（0 字节的件＝那格没交）

```
$ wc -c .scratch/wisp/probes/35/a6/00-anchor.md .scratch/wisp/probes/35/a6/verdict.md .scratch/wisp/probes/35/a6/01-readings.md   rc=0
三枚全 >0；⛔ 本腿写面内 0 字节件＝0 枚
$ find .scratch/wisp/probes/35/a6 -type f -size 0 | wc -l   rc=0 → 0
```
