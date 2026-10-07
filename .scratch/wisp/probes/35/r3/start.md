# 35-r3 起手锚（票 `:63` "Inbound judge must ride the page's own envelope" 的 (d) 三支守卫红）

日期 `2026-10-07 21:1x +08`。HEAD `1d70fb2b`。porcelain 行数 753（既有脏，本腿不认领）。

## 本程射程（唯一）
票 `:63` 的 **(d)-1/-2/-3**：三支入向守卫，各点各的、一句文案不许共用，且**经页面那条边**（甲③ 传输，跑 JS 行为尺那条路径）进来。⛔ 不直调 `dispatchRaw`、⛔ 不 `w.Eval("window.wispDispatch(…)")`（N6，判别力为零）。
⛔ 不改产码、⛔ 不碰任何 AC 框、⛔ 不改 `panel_transport_35r2_test.go` 既有断言。

## 三支守卫现量行号（`internal/panel/bridge.go`，CR 计数 0，与 `git show HEAD:` 逐字节一致）
- 名册判：`:132` `if !knownComposerMethod(r.Method) {` ⇒ 拒 `:133`，文案含 `方法 %q 不是面板 composer 通路的能力入口`。
- 来源判：`:135` `if strings.TrimSpace(r.Source) != ComposerRequestSource {` ⇒ 拒 `:136-137`，文案含 `来源 %q 不是 %q，按伪造/串台拒绝（requestId=%q）`。
- requestId 判：`:139` `if strings.TrimSpace(r.RequestID) == "" {` ⇒ 拒 `:140-141`，文案含 `缺少 requestId，无法与审计/卡片对齐（method=%q）`。
- （解析判 `:128-129` `不是可解析的封套`，⛔ 不在本程三支之列。）
- `ComposerRequestSource = "panel-composer"`（`bridge.go:30`）；名册六枚 `bridge.go:42/66/67`＋`knownComposerMethod` case `:148`：`panel.mode.request`/`panel.workspace.request`/`panel.attachment.add`/`panel.message.send`/`config.get`/`config.set`（`panel.approval.request` 不在册）。
- 守卫次第：解析 `:128` → 名册 `:132/:133` → 来源 `:135/:136-137` → requestId `:139/:140-141`（编排者 13:4x 复跑裁，票 `:69`）。

## 到达链（页面那条边）
`ComposerDispatch.Handle`（`internal/panel/composer_dispatch.go:153`）：`ParseComposerRequest` 出错 ⇒ `d.record` + `return RefusedEnvelopeForUser(req, err), err`（`:158`），拒的文案经 `RefusedEnvelopeForUser`（`bridge.go:166`，`面板请求被拒绝 [%s %s]：%v`）把守卫那枚 `err` 原样带进回复串。门闭包 `installPanelTransport`（`panel_host_windows.go:694-696`）`reply, _ := m.dispatchRaw(...)` ⇒ 丢弃 err、只把回复串交给 JS 层 ⇒ 行为尺的 `fakeDoc35r2.lastReply` 拿到的就是这句文案。入口复用 `fakeDoc35r2`（r2，同包 `main`）：`openDocument()` 跑 chromium.go:112 垫片 + Bind 桩 + 产码 `panelPostMessageForwardInit` 钩子原文，`pagePost(envelope)` 即页面 `window.chrome.webview.postMessage(<信封>)`。

## (d)-1 与既有钉的差别（本格新料）
既有钉＝`internal/panel/l2_grant_boundary_test.go`（`:1964`/`:1965` 拿页面原始串要求 `knownComposerMethod` 拒、`:1251-1267` 要求拒 11 枚候选名）——那是**在 internal/panel 包内直调 `knownComposerMethod`/`ParseComposerRequest`**、判的是"名册函数认不认识这个名字"。本程 (d)-1 的新面＝**同枚 `panel.approval.request` 经页面 postMessage 那条边（甲③ 传输 + JS 行为尺）进来时，是名册那一道守卫先拒它、且回复文案点名"不是面板 composer 通路的能力入口"**——既有钉走不到这条边（不经传输、不经 `Handle`、不产生页面可见的拒句）。⛔ 不拿既有钉抵账。

## A684 两面恒真的射程声明
`A684` 记的这把 1824 行尺的两面恒真＝①夹具 `:1238` 原生 `postMessage` 桩丢弃 receiver（`native.call(cw,msg)`→`native(msg)` 两发都绿）、②`:137 jsObject.set` 无可写性概念（覆写被静默忽略那支测不出）。本程三支断言只依赖"信封作为不透明串经转发钩子→门→Handle→ParseComposerRequest→回复串"这条**已经过 `A683`/`A684` 复量、确有牙的**投递面（`doorRounds`/`lastReply`），⛔ 不依赖 receiver 可写性或 this 绑定。若某支恰好落进那两面，交件里具名说清楚（初步判：三支都落在投递/文案面，不落那两面；见 impl.md §⑤）。

## 门禁基线（起手，长跑前）
- `sh scripts/d22scan.sh`：随后跑。
- `go vet ./cmd/wisp/ ./internal/panel/`：随后跑。
- `go test -count=1 ./internal/...` 起手红名册（既有 5 枚）：`TestApprovalCardViewJSONKeysMatchFrontendTypes`/`TestC21DesignTokensFourWayAgree`/`TestC21TableColourRowsMatchTokensCSS`/`TestComposerContractTypesMatchFrontend`/`TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme`；间歇门 `TestResolvePerCallBudget` 与间歇 `TestTicket223ModeLooseningChangesTheRunningModeAfterAllow` 不计账。
- `cmd/wisp` 本机 `go test` 载入即死 `0xc0000135` ⇒ harness：`go test -c -o /d/tmp/wisp35r3/panel.test.exe ./cmd/wisp/` + 拷三枚 sherpa DLL 到 exe 旁 + CWD=`cmd/wisp`。
- 冲突：无（简报与盘上一致；引框只用行号＋名字，遵 `A684` §2）。
