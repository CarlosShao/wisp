# 35-v3 对抗验收｜票 35 `:63` 的 **(d) 三支守卫红**（非实现者腿）

被验对象＝`9b2551f2`（起手锚）／`7f9d6e40`（交件，15 枚件，`git show --stat` 现数）。
本腿起手 HEAD＝`54091f67`；本腿两笔＝起手锚 `2078ec14` ＋ 交件笔（见回报）。
只写台账与证据件：⛔ 未改产码、⛔ 未改任何 AC 框、⛔ 未翻勾、⛔ 未 push、⛔ 未开真实 WebView2 窗。
所有突变＝`go test -c -overlay` 于 `D:/tmp/wisp35v3` 的副本，**仓里一字未改**（`mutants-verify.txt` 逐份 diff＋改动行枚数现量）。

---

## V1｜三支各点各的（自造突变，⛔ 未复用 r3 副本）——**成立**

我自造三份，每份**恰 1 行改动**（`changed-line-count: 2` ＝一行 `<`＋一行 `>`，见 `mutants-verify.txt`），中和的是**判的条件本身**：

| 中和的判（overlay 改的行） | 改成的形状 | (d)-1 名册 | (d)-2 来源 | (d)-3 requestId | run rc |
|---|---|---|---|---|---|
| `bridge.go:132` `if !knownComposerMethod(r.Method)` | `if false && !knownComposerMethod(r.Method)` | **红** | 绿 | 绿 | 1 |
| `bridge.go:135` `if strings.TrimSpace(r.Source) != ComposerRequestSource` | `if false && strings.TrimSpace(r.Source) != ComposerRequestSource` | 绿 | **红** | 绿 | 1 |
| `bridge.go:139` `if strings.TrimSpace(r.RequestID) == ""` | `if false && strings.TrimSpace(r.RequestID) == ""` | 绿 | 绿 | **红** | 1 |

基线（未突变 exe，`-test.run TestInbound`，本腿自跑）＝**3 PASS／0 FAIL／rc=0**（`baseline.txt`），三支各自的 `door reply` 逐字与 `A685` §1 所列一致。
⇒ **"只有它那枚红、另两枚仍绿"这一条，本腿自己造、自己跑，成立。**

⚠ **一处具名顶回（台账 `A685` §2 末句与盘上读数不符）**：中和名册判那发，(d)-1 的红 reply 逐字＝
`面板请求被拒绝 [panel.approval.request (无 requestId)]：panel: composer request refused: 来源 "" 不是 "panel-composer"，按伪造/串台拒绝（requestId=""）`
——接住它的是**来源判 `:135`/`:136-137`**（因为页面那枚 approval 信封既无 `source` 也无 `requestId`，名册判一被中和，**下一道**就是来源判）。
`impl.md` §② 写的"下一道来源判接住它"＝**对**；`A685` §2 我读到的那句"它落到了别支的文案……撞上 requestId 判"＝**与读数不符**（requestId 判在来源判之后，这枚信封走不到那一层就被拒了）。红句逐字见 `run-m-roster.txt`／`v3-error-texts.txt`。

## V2｜★反形：文案被**合并**时会不会红（r3 没做过的那一发）——**成立，且两形都红**

我在 V1 之外自造两枚"合并"突变（同一枚 `bridge.go` 副本，三处 `return r, fmt.Errorf(...)` 全换成**同一句**；diff 见 `mutants-verify.txt`，`changed-line-count: 8`＝三条语句合并后的行数差）：

- **M-merge-generic**＝三句都改成 `"%w: 面板请求被拒绝"`（无任何本支关键词）⇒ **3/3 全红、rc=1**。红的是**每一支**，机制＝`reply does not name guard "…"`（本支关键词消失）。逐字如 `reply="面板请求被拒绝 [panel.approval.request (无 requestId)]：panel: composer request refused: 面板请求被拒绝"`。
- **M-merge-all3**（更狠的一形）＝三句都改成**同一句、但把三枚关键词全写进那一句**（`入向请求不合（不是面板 composer 通路的能力入口／按伪造/串台拒绝／缺少 requestId）method=%q`）⇒ **3/3 全红、rc=1**，红的是**另一条断言**＝`leaked a DIFFERENT guard "…" alongside "…"`（"另两支关键词不许在"那一支 prong）。

⇒ **判回简报的问**：这组判据**不是**只会"在字符串里找关键词"。它有两枚 prong，各自独立有牙：本支关键词必须在（M-merge-generic 打中它）、另两支关键词必须不在（M-merge-all3 打中它）。**文案合并／压成一格这两种形状都看得见。**
⛔ 我没有取到"恒真面"这一格——**这一条上没有恒真面可报**；简报预留的那句"若三支都还绿"没有兑现（三支不绿）。

## V3｜三支会不会"传输坏了也照绿"

**① 产码钩子被中和（`w.Init(panelPostMessageForwardInit)` → `_ = panelPostMessageForwardInit`，恰 1 行，`m-hook`）＝3/3 全红、rc=1。**
红句逐字（三支同形）＝`inbound guard RED: the page's postMessage drove the door 0 time(s), want 1. … nativeExitCalls=1 doorRounds=0 … unboundSlots=1 evalThrew=1 resolved=0`。
⇒ 票 `:63` 那句 ★Falsifiability（"删掉/中和转发语句必须让这枚具名用例红"）**本腿自己复现成立**；三支不是恒绿。

**② `spy.reqs==0`（"router 未达"）不是恒真**：我造了一发正控 `m-legal`＝**守卫一字不动**，只把 (d)-3 那枚信封补上合法 `requestId`（`{"method":"panel.mode.request","requestId":"pc-35v3-legal","source":"panel-composer",…}`，测试文件副本 1 行改动）⇒ **(d)-3 红、另两支绿、rc=1**，红句逐字＝`a refusal should never reach the router, but the mode handler ran 1 time(s): [{Method:panel.mode.request RequestID:pc-35v3-legal …}]`。
⇒ 夹具里 router 是**可达**的、`modeSpy35r1` 是**真记账**的，所以 `reqs==0` 那条不是"本来就到不了"的空断言。⚠ 一格如实限定：这枚断言在我能造的所有形状里都与 `mine-absent` **同发**（一旦抵达路由，回复串就是 `""`），所以它是**佐证**而非可独立发火的尺——不构成恒真，但也不该被读成"单独盯着路由"。

**③ r3 自己划的射程声明＝复核成立，并且我把它扩了一格**：
- 现量 `cmd/wisp/panel_transport_35r2_test.go` 里 `writable|configurable|defineProperty` ＝ **0 处**；`jsObject.set`（`:137-142`）无条件覆写；原生桩 `:1238` 是 `func(_ jsValue, args []jsValue)`（receiver 丢弃）。
- 产码钩子（`panel_host_windows.go:673` 起）里那句正是 `cw.postMessage = function (message) {…}` 与 `native.call(cw, message)`（`:680`）——**三支的到达全靠这两句在夹具里生效**。
- ★**我实测了一发 r3 没测的**：`m-this`＝把 `native.call(cw, message)` 换成裸 `native(message)`（恰 1 行，正是 `A684` 的 M-A 形状）⇒ **3/3 全 PASS、rc=0**（`run-m-this.txt`）。
⇒ 结论：**这三支对 `A684` 的两面恒真（M-A 丢 `this`、M-B 无可写性）都是瞎的**，r3 只具名了 M-B 那一面，M-A 同样适用。这三枚绿**只**证"守卫次第＋各支文案"这一维，⛔ 不得读成"传输在真浏览器里成立"；那一格仍只由 `:52` 的 winlive 真窗凭据撑着（票 `:75` 另程处置）。

## V4｜不抵账与禁区——**成立**

**① (d)-1 与既有钉的差别是真的。** 我读了 `internal/panel/l2_grant_boundary_test.go`（`:1245-1275`、`:1955-1975`，摘录存 `existing-l2-nail-excerpt.txt`）：
- `:1264` `if knownComposerMethod(cand)`＝**包内直调名册函数**，第一枚候选名正是页面原文 `panel.approval.request`；
- `:1970`/`:1989` `parsed, err := ParseComposerRequest(raw)`＝**直调解析函数**，判据只有 `err == nil` 与 `errors.Is(err, ErrComposerRequest)`。
⇒ 它缺的**那一环**有两件：(i) **没有任何传输**——该文件里 `installPanelTransport|pagePost|postMessage|fakeDoc` 命中 **2 处，且都在注释里**（`:20`、`:72`，是描述页面而非执行页面），`internal/panel` 也**结构上**碰不到 `cmd/wisp` 那侧的传输；(ii) **不点名是哪一道守卫拒的**——三句关键词在该文件出现 **0 次**（`v4-checks.txt`）。
★**我自己补的一发实测（这一发把我原来的猜测纠回来了，具名报）**：我原以为既有钉对"名册判被中和"是瞎的——**不对**。用同一份 M-R overlay 去跑既有钉：`go test -count=1 -overlay ov-m-roster.json -run 'TestGrantWireShapesAreRefusedAtTheDoor' -v ./internal/panel/` ＝ **rc=1**（未突变对照 rc=0／3 PASS／0 FAIL），因为它那五枚线形里有 **4 枚带 `requestId`＋`source`**（`panel.approval.request r-1`／`approval.decide`／`panel.l2.allow`／`panel.grant`），名册判一中和就**一路被 accept** ⇒ 红（逐字见 `v4-existing-nail-under-mutant.txt`）。
⇒ 所以真正的"缺那一环"是更窄的两条，(d)-1 的新料仍成立：① 既有钉**不经传输**（没有 `doorRounds`／没有页面可见的拒句）；② 既有钉只判 `err != nil`，**说不出是哪一道答的**——对页面今天真发的那枚原文 `{method, correlationId, outcome}`（无 `source` 无 `requestId`），名册判被中和后它**照样算通过**（改由来源判拒），而 (d)-1 在这一发是**红**的（我的 `run-m-roster.txt` 逐字）。⇒ ⛔ 不许抵账这一条成立，但**成立的理由要按上面这两条写**，⛔ 不许写成"既有钉看不见名册被中和"。反之 (d)-1 也不覆盖既有钉（11 枚候选名＋词根＋pool 审计仍只在那枚钉里）。

**② 禁区现量**（`anchor.txt`）：`git diff --stat 1d70fb2b..HEAD -- cmd/wisp/panel_host_windows.go internal/panel/ frontend/` ＝ **空**；`frontend/`＋`internal/panel/` 单独对拉亦空；`bridge.go` blob `bebe8e70` 与 `1d70fb2b` **逐字节同**；`panel_transport_35r2_test.go` blob `7185ab56` **未变**（＝既有断言一字未改，不必逐行比）；工作树 `bridge.go` CR 计数 0＝与 HEAD 一致；`m.dispatchRaw(` 产码调用点 **1** 枚。
本腿新面只有 `.scratch/wisp/probes/35/v3/**` 与票面**追加**一段（⛔ 未动 AC 框、未翻勾）。

**③ 三句关键词逐字对 `bridge.go`**（`v4-checks.txt`）：每枚在 `bridge.go` 出现**恰 1 次**；Python 对拉三枚关键词 × 三句原文＝**每枚只出现在自己那一句**（`cross-appearances-total= 3`，⛔ 无一枚是他支子串）。
⚠ 一格近失（不是违规，登记备查）：`RefusedEnvelopeForUser` 的抬头会印 `(无 requestId)`（`bridge.go:160`）与 `面板请求被拒绝`（`:166`）——前者与 (d)-3 的关键词 `缺少 requestId` 只差一个字、后者与 M-merge-generic 那枚共有句**完全同名**。⇒ 今后若有人把关键词换成 `requestId` 或 `面板请求被拒绝`，就会有一支的"另两支不在"变成恒真；这一格是**选词**决定的、不是判据自己长出来的。

## V5｜门禁复跑（每件自落 `rc=`；无 `.out`、无 0 字节件）

| 尺 | 读数 | 件 |
|---|---|---|
| `sh scripts/d22scan.sh` | **rc=0** clean（ban #8 `cmd/` 108 枚、`internal/` 514 枚含 `_test.go`；live scope work bans #1-5 `cmd/`=38） | `gate-d22scan.txt` |
| `go vet ./cmd/wisp/ ./internal/panel/` | **rc=0**（成功＝零输出，件里只有 `rc=0` 一行，非 0 字节） | `gate-vet.txt` |
| `go test -count=1 ./internal/...` | **rc=1**；逐名红＝5 枚；与简报给的既有 5 枚**作差双向皆空**（`known-only:` 空、`new-only:` 空，known=5 got=5） | `gate-internal.txt`／`internal-reds.txt` |
| 间歇门 | `TestResolvePerCallBudget`／`TestTicket223ModeLooseningChangesTheRunningModeAfterAllow` 在 internal 那发命中 **0** 次 ⇒ 不计账 | `v4-checks.txt` 同批 |
| 三支 harness | 基线 **3 PASS／0 FAIL／rc=0**；八发突变逐枚 rc 见 `mutants-run.txt`＋`run-m-*.txt` | `baseline.txt` |

harness 逐字照简报做：`go test -c -o /d/tmp/wisp35v3/panel.test.exe ./cmd/wisp/`（build rc=0）＋ GOMODCACHE 三枚 DLL（`onnxruntime.dll`／`sherpa-onnx-c-api.dll`／`sherpa-onnx-cxx-api.dll`）拷 exe 旁＋CWD＝`cmd/wisp`。⛔ 未开真窗。
`cmd/wisp` 整包（本腿独占，未突变 exe，⛔ 无 winlive tag）最终读数＝**rc=1／338 PASS／5 FAIL／1 SKIP**，三支在整包里同发 **3 PASS**；那 5 枚 FAIL 逐名＝`TestPanelHostRealWindowHopAndLifecycle`、`TestAC4FocusReturnToPriorWindowGap33r5`、`TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe`、`TestAC14AwaitedBindingReplyReachesThePage`、`TestAC14GoSideEvalPushReachesThePage`＝`A684` §3-W3 记过的窗口依赖既有红族，⛔ 不记本程账；间歇门 `TestTicket223ModeLooseningChangesTheRunningModeAfterAllow` 这一发 **PASS**。逐名读数存 `cmd-wisp-full-summary.txt`（原始 1500+ 行按预算留 `D:/tmp/wisp35v3/cmd-wisp-full.txt`，⛔ 不入库、⛔ 未删）。

---

## ★总结论

**① 票 `:63` 现在能不能翻＝能翻（就 (d) 而言，我这格没有异议；但翻勾权在编排者）。**
(a)(b)(c) 已由 `35-v1`／`35-v2` 复认（不归本程）；**(d) 这一程我全绿了它要求的一切**：三支各点各的（V1，自造三份 1 行突变）、三支不共用一句（V2 两形都红）、转发语句被删即红（V3-①，★Falsifiability 兑现）、入口是页面那条边（⛔ 未直调 `dispatchRaw`、⛔ 未用 `w.Eval(wispDispatch)`：两处命中都是注释 `:24`/`:43`，调用形状计数 0）、⛔ 未拿既有钉抵账（V4-①，差别我独立读码坐实）。
⚠ 翻勾附一句**必须一起写**的射程限定：这三枚绿**不含**任何"传输在真浏览器里成立"的内容（V3-③，M-A 实测 3/3 绿）。

**② 这三支的牙覆盖到哪一层（三种形状 × 看得见／看不见）**

| 形状 | 看得见吗 | 凭哪条断言 | 读数 |
|---|---|---|---|
| 守卫被中和（三道判各一枚） | **看得见** | 本支关键词必须在 | 每发恰 1 枚红、另 2 枚绿，rc=1 |
| 文案压成一句（不含关键词） | **看得见** | 本支关键词必须在 | 3/3 红 rc=1 |
| 文案压成一句（含全部三枚关键词） | **看得见** | 另两支关键词必须不在 | 3/3 红 rc=1 |
| 转发钩子整个被摘（`w.Init` 中和） | **看得见** | `doorRounds==1` | 3/3 红 rc=1（`doorRounds=0 unboundSlots=1 evalThrew=1`） |
| 信封合法却仍断"未达 router" | **看得见** | `spy.reqs==0` | 正控 m-legal 1 枚红 rc=1（与 mine-absent 同发＝佐证级） |
| 钩子内 `native.call(cw, …)` 丢 `this`（M-A） | **看不见** | — | **3/3 绿 rc=0**＝恒真面 |
| 真浏览器里 `cw.postMessage` 覆写被静默忽略（M-B） | **看不见** | — | 夹具 0 处可写性概念＝结构性盲区（本腿未再造 M-B 反证，A684 已造） |
| 名册守卫集合本身／解析判 `:128-129` | 不在射程 | — | 三支不测（`impl.md` §① 具名，一致） |

**一句话**：这三支的牙咬得住**守卫层**（中和任一道判）与**文案层**（合并／压格，两形都咬得住），也咬得住**传输整段失效**；咬不住的是**夹具看不见的 JS 语义层**（`this` 绑定、属性可写性）——那一层今天仍只有 `:52` 的 winlive 真窗凭据，归票 `:75`。
