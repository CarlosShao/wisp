# 票 253 — 面板入向那三枚尺洞：可达性钉只认 `postMessage(`、`config.get`/`config.set` 不带 `panel.` 前缀掉在两把名册尺之外、旧封套被宿主库**静默 `resolve(null)` 且 Go 侧零审计**

**立票时刻**：2026-10-02 10:0x +08，锚点 HEAD `f25f9572`（`dev`）
**来路**：只读腿 `114-a2`（`.scratch/wisp/probes/114/a2/census.md`，241 行 / 57,999 字节，`8d9f4102`→`f5f9cc34`→`981f08b2`）的 R9／R10 两条＋它新增的一条缺陷；台账 `A519`。
**性质**：⚠ **第三枚算产品行为**（页面以为话发出去了、Go 侧什么都没收到、也没有任何审计行），前两枚是仪器射程。
**⛔ 本票不接管票 114 的 9 枚未勾格**——那 9 枚是"差集"（档位／附件／工作区／消息），归票 114 自己的落地腿。

## 现量（腿的读数；行号一律〔待验〕，落地腿起手先复跑）

1. **可达性钉看不见生产那扇门**。钉只认 `.postMessage(`（`internal/panel/composer_test.go:381`〔待验〕），而生产入口是 `w.Bind("wispDispatch")` → `window.external.invoke` 那条 RPC（`cmd/wisp/panel_host_windows.go:551` 一侧）。⇒ **"这条链路可达"那一格，今天是用一把看不见门的尺在判**。
2. **名册两把尺的射程之外多出一批名字**。票 248 新增的 `config.get`／`config.set` **不带 `panel.` 前缀**（`internal/panel/bridge.go:66-67`，我现读为真），因此掉在 `routeLiteralRe`（`internal/panel/composer_test.go:394`〔待验〕）与 `panelMethodRe`（`cmd/wisp/git_test.go:394`〔待验〕）之外 ⇒ 票 114 的 **AC#8 那个洞被削宽**：新增方法名可以在**不触发任何名册尺**的情况下进树。
3. **旧封套静默丢弃**。页面若还按旧形状直接 `postMessage`，宿主库会**把这次调用的 Promise 解析成 `null` 并返回**、**Go 侧零审计**（腿报位置在依赖模块 `webview.go:162-168`〔待验，且属第三方模块，本仓不修上游〕）。⇒ 用户视角的症状＝**点了按钮，什么都没发生，日志里也什么都没有**。

## 要建什么

- [ ] **AC#1 可达性那枚钉改成问能力**：判据不许再认某个 JS 调用的**词面**，要问**"页面发起的一次调用，Go 侧有没有收到并回执"**（载体沿用票 33 那族"由页面自己报回"的定式：`REPLIED` 那形）。⚠ 正控必配：造一发"只 `postMessage`、不 `wispDispatch`"的假腿 ⇒ 指名用例**必须红**。⛔ 不许把两种拼法都列进白名单了事（那是把尺调宽，不是把门修好）。
- [ ] **AC#2 名册尺的射程要覆盖"不带前缀的那批"**：`config.get`/`config.set` 这类非 `panel.` 前缀的方法名必须被某把尺数到（要么两把尺各自扩射程，要么单立一枚"全量入向名册"尺）。判据＝种一枚"新增方法名没进名册"的假腿 ⇒ 指名那一步**必须红**；⛔ 现在它今天**是绿的**（这就是缺陷本体）。
- [ ] **AC#3 静默丢弃那一支要有声**：⛔ 不许改上游库；本仓射程内能做的是**在 Go 侧给"收到但形状不认识"落一条审计/告警**（或至少让面板宿主那条腿**具名拒收并回一个错误**）。判据＝一发假腿按旧封套发一次 ⇒ ①Go 侧有可读的现场（审计行或错误文本），②页面拿到的不是无声 `null`。
- [ ] **AC#4 三枚各自结线，不许合并成一句"入向修好了"**：AC#1 是钉的射程、AC#2 是名册的射程、AC#3 是产品行为——结案时逐格给读数。

## 禁区

- ⛔ 不许动 `frontend/**`／`design/**`（**既不读也不引**）；界面侧那半（改 `postMessage` 调用形）属 owner 自己用的那枚 agent，本票**只写 Go 侧**，需要它配合的地方写进"归口"由 owner 带话。
- ⛔ 不许改 `internal/panel/l2_grant_boundary_test.go` 那批冻结件射程之外的口径；⛔ 不许放宽任何名册尺来"让它绿"。
- ⛔ 阈值／SLO／golden／`internal/observe/thresholds.go`／`tools/d22scan/allowlist.txt` 一字节不动；⛔ 不许 `t.Skip`；⛔ 不许把 SKIP 读成通过。
- Git 纪律：只 commit 不 push；显式 pathspec；⛔ `git add -A`／`.`／`-a`／`--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`；临时件只建不删。

## 归口与排程

- 写面＝`internal/panel/**`＋`cmd/wisp/**` ⇒ **两枚包此刻都有验收腿在飞**（`198-v1` 正在整包跑 `cmd/wisp`）⇒ **本票按住**，起跑判据：`198-v1` 退出且 `git status --porcelain internal/panel cmd/wisp` 为空。
- 与票 114（差集九格）、票 242（`bindDigest` 绑定层与出向读面零仪器）**互不吞并**：那两票各自的名字与格子不许并进本票。
- 需要界面侧配合的只有一件：**把面板里的调用统一走 `wispDispatch`**（不带 `panel.` 前缀的那两名除外，它们就是 `config.get`/`config.set`）。这一句要 owner 带给界面侧 agent，⛔ 本编队永不转达。
