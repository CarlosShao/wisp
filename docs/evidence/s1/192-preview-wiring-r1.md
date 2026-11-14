# 192 - 产物预览的接线腿（`cmd/wisp` 三步）：读数表

- 票：`192-the-browser-and-artifact-preview-stack-is-unclaimed-and-sits-on-the-undefined-web-search-ruling.md`
- 前一腿（栈本体）：`docs/evidence/s1/192-preview-ac12-r1.md`，产物 `internal/panel/preview_192.go`（778 行）＋判据 719 行
- 本腿写面：`cmd/wisp/**`（装配＋处理器）＋本件＋`.scratch/wisp/probes/192/**`＋票面 Progress log（只追加）
- 本腿**一字节未动**：`internal/panel/**`、`internal/llm/**`、`internal/tools/**`、`tools/d22scan/**`、`frontend/**`、`design/**`；`docs/PLAN.md`／`docs/specs/**`／`allowlist.txt`／`thresholds.go`／golden 0 字节；C17 方法一个没加（`internal/panel/bridge.go` 未改）

---

## 一、三步各落在哪一行

| 步 | 落点 | 改了什么 |
|---|---|---|
| ① 预览分支＋`doc.Headers` 拼进响应 | `cmd/wisp/panel_control_windows.go:571`（`panelAssetLookup` 新签名）、`:603`（分支）、`:504-512`（响应头的取舍）、`cmd/wisp/panel_preview_192.go:216`（`panelPreviewAssetResponse`） | 查找函数从 4 元组改成 `panelAssetResponse` struct，多了 `Headers` 一个字段；`/__artifact/...` 走预览栈**并把 query 一起交给它**（`?offset=&limit=` 是截断指针，旧的 `IndexAny("?#")` 会把它切掉）；预览响应的头一律用 `PreviewDocument.Headers` 拼（CRLF 连接），只有它为空时（嵌入包／miss／预览拒绝）才回落到面板那条兜底 CSP |
| ② 组合根装配 `PreviewResolver` | `cmd/wisp/panel_resident_windows.go:342`（`newResidentPanelManager` 造 `panelPreviewSource`）＋`:390`（`withArtifactPreview`）＋`panel_host_windows.go:323`（选项与字段）＋`panel_host_windows.go:671`（`installAssetFilter(m.assets, m.preview)`）＋`panel_inbound.go:281`（链把同一个源交给 message 门） | 生产装配路径上一个 `panelPreviewSource`：懒开 `memory.Store`（首个预览请求才开库，与票 214 的 ingest 同姿态）→ `panel.NewPreviewResolver(store, memory.ValidArtifactName)`。同一个对象同时交给窗口的过滤器与 message 门 ⇒ 一个 artifacts 目录、一个名字守卫、一套预算 |
| ③ message 处理器接 `QuoteForModel` | `cmd/wisp/panel_preview_192.go:294`（`panelPreviewQuoteDoor.HandleMessageRequest`）＋`cmd/wisp/panel_inbound.go:372`（`Message:` 从 nil 改成这枚门）＋`cmd/wisp/resident_windows.go:576-586`（把进程的 `src.run.c25` 绑进门） | 从 0 个调用点变成 1 个，且这枚在 message 处理器上、不在测试里 |

### 门读的关联是怎么来的（没有新增字段、没有新增方法）

页面原本没有任何办法说"这段文字来自产物 X"。加字段不行：composer 信封是所有已答方法共用的（票 248 对凭据写过同一条裁定），加方法更不行（AC#4 逐字"C17 不加方法"）。所以门读的是 message 信封上**今天没有任何人读**的那个键：`path` 装**同一发预览响应用的同一个同源地址**（`panel.PreviewURL(name)`，可带 `?offset=&limit=`，query 被忽略、整份产物是主体），`text` 装选中文字。这不是第二套地址空间——就是页面取那几字节时用的那条 URL，于是"用户选中的文字"和"那几字节出自哪份产物"是同一个地址。

---

## 二、三条判据的可判事实

### 判据 1：预览不再 404（产品路径上的状态与字节数）

同一个函数、同一个 URL，只差"有没有把预览源交给它"：

| 情形 | 状态 | 字节 | content-type | 头字节 |
|---|---|---|---|---|
| **before**（没交预览源＝本票修之前的形状） | **404** | **0** | — | 0 |
| **after**（生产装配交给它的那个源） | **200** | **129** | `text/html; charset=utf-8` | 192 |

URL = `https://wisp.panel/__artifact/page.html`（`panelOrigin + panel.PreviewURL("page.html")`），产物是测试写进真 artifacts 树的 129 字节 HTML。

复跑命令（`logs/10-cmdwisp-full.txt` 里 grep `step 1`）：

```
scripts/wisp-cli-staged-tests.sh -test.run Test192PreviewLookupAnswersTheProductPath -test.v
```
读数原文：
```
step 1 before: panelAssetLookup(no preview source, "https://wisp.panel/__artifact/page.html") = status 404 "Not Found", 0 byte(s), 0 header byte(s)
step 1 after:  panelAssetLookup(preview source, "https://wisp.panel/__artifact/page.html") = status 200 "OK", 129 byte(s), content-type "text/html; charset=utf-8"
```

### 判据 2：CSP 真的来自 `doc.Headers`（逐字）

预览响应的整个头块，逐字（`\r\n` 已写成转义）：

```
Content-Type: text/html; charset=utf-8\r\nX-Content-Type-Options: nosniff\r\nCache-Control: no-store\r\nContent-Security-Policy: sandbox; default-src 'none'; img-src data:; style-src 'unsafe-inline'
```

判"它来自 `doc.Headers`、不是面板那条兜底"的四把尺，全在 `Test192PreviewLookupAnswersTheProductPath` 里：

1. **含** `Content-Security-Policy: sandbox; default-src 'none'`（判据逐字要求的那一段）；
2. **不含** `panelCSPHeaderValue`（面板兜底那条策略逐字，出现在这里就说明类别策略没送到）；
3. 头块 **不等于** `panelResponseHeaders(200, ct)`（兜底构造器的输出）；
4. 与 `PreviewResolver.Lookup(...).Headers` 用 `\r\n` 连接后**逐字节相等**（唯一能证明"是那一份头、不是长得像"的尺）。

外加：响应体与响应头里都**不出现** artifacts 目录路径（浏览器永远不知道字节住在哪）。

预览**拒绝**那一侧：`/__artifact/no-such-artifact.txt` = 状态 404、正文 119 字节（解析器自己那句中文）、头块 0 字节 ⇒ 回落到面板兜底 CSP（`panelResponseHeaders` 仍带 `panelCSPHeaderValue`）。即"能服务的预览用自己的策略，不能服务的预览仍受面板策略管"。

### 判据 3：`QuoteForModel` 有产码调用点（git grep，改前 → 改后）

复跑：`sh .scratch/wisp/probes/192/wire-r1/01-callsites-before-after.sh`（全文 `logs/01-callsites.txt`）

| 尺 | 改前（`HEAD`） | 改后（工作树，`--untracked`） |
|---|---|---|
| `'.QuoteForModel('`，非测试 `.go` | **0** | **1** — `cmd/wisp/panel_preview_192.go:294` |
| `'panel.NewPreviewResolver('`，`cmd/`＋`internal/`，非测试 | **0** | **1** — `cmd/wisp/panel_preview_192.go:164` |
| `'PreviewPathPrefix'`，`cmd/`，非测试 | **0** | **4**（1 行注释＋3 处代码：`panel_control_windows.go:603` 的分支、`panel_preview_192.go:272/274`） |

⚠ 那把尺必须带 `--untracked`：`git grep` 跳过未跟踪文件，而不带它时本腿加的那一整个产码文件不被计，改前改后都读 0（"什么都没变"的假读数）。脚本里写死了这个开关。

那一枚调用点所在的处理器是 message 门，端到端读数（`Test192MessageDoorStampsThroughQuoteForModel`，走 `PanelManager.dispatchRaw`＝页面自己那扇门）：

| 情形 | 结果 |
|---|---|
| 引擎未绑 | 拒绝（`ErrPreviewRefused`：没有 C25 引擎），引擎里 **0 条** mark |
| 引擎已绑＋文字真在产物里 | `QuoteForModel` 盖戳后**拒绝转发**（`ErrNoHandlerAttached`：消息派发通路属票 35）——`prov.CheckText("panel-preview-quote", ChHTTP, 原文)` **命中**，`SrcTool=task.output`、`Origin=artifact:quoted.txt` |
| 页面自己编的文字（产物里没有） | 拒绝（`ErrPreviewRefused`），引擎里 **0 条** mark |
| 审计行 | 有，且**不含**被引用的原文（只有 artifact 名与字数） |

---

## 三、具名残余（本腿没解决的，逐条写清"今天什么不响"）

1. **消息转发通路仍不存在（票 35 的 socket）**。门能做的事只有：证明＋盖戳，然后**拒绝转发**并说明两半。所以今天"面板里选一段预览文字发给模型"这条路**走不通**——走不通的原因是票 35 没落地，不是这道门。
2. **盖戳没有继承到任务 scope**。mark 落在进程引擎的 `panel-preview-quote` 这个 scope 上；真正把文字送进模型的那次调用会在**任务** scope 里被检查，而"把这个 mark 带进任务 scope"是转发器（票 35）那一侧的事（票 382wr2 给 TTS 做过同款继承）。⇒ 今天这道门**不能**宣称"预览文字出不去"，它能宣称的是"预览文字只能经这一道门、且一定带着可被查到的戳"。
3. **页面侧零写面**。今天没有任何前端发这个形状（`path` 装预览地址＋`text` 装选中文字）；本腿给的是**协议与语义**，没有 UI 证据。`frontend/**` 是外部 agent 的写面，本腿 0 字节。
4. **PDF 的像素仍不承诺**（前一腿的声明依赖原样继承：`Render=host-pdf`）。
5. **本机整包 `./cmd/wisp` 的 2 枚已知红不是本腿的**（见 §四）：`TestAC4Focus…`（前台被别的真窗口占着）、`TestAC246DevLeg…`（**驻留球在跑**，编排者 16:23 起的）。**具名记下，不修、不删。**

---

## 四、门禁读数（每条可复跑）

| 门 | 命令 | 读数 | log |
|---|---|---|---|
| vet | `go vet ./cmd/wisp/` | rc **0** | `logs/03-vet-gofmt.txt` |
| 格式 | `gofmt -l`（两个新文件＋三个改过的产码文件） | **0 文件** | 同上 |
| 格式 | `go run mvdan.cc/gofumpt@v0.12.0 -l <7 个文件>` | **0 文件** | 同上 |
| d22scan | `tools/d22scan/d22scan.exe -root .` | **clean**，`cmd/` 产码 **285** 枚、`internal/` 产码 **405** 枚 | `logs/02-d22scan.txt` |
| 新判据（定向） | 载具 `-test.run Test192` | 6 枚顶层全 PASS（`Test192PreviewLookup…` / `…PreviewRefusal…` / `…ResidentAssembly…` / `…MessageDoorStamps…` / `…MessageDoorIsNotASecondRelay` / `…PreWiringReading…`） | `logs/04-test192.txt` |
| 整包 | `scripts/wisp-cli-staged-tests.sh`（⛔ 裸 `go test`，本机 DLL 会出假红 `0xc0000135`；⛔ 同时跑两发） | 见 §五 | `logs/10-cmdwisp-full.txt` |

---

## 五、整包 `./cmd/wisp` 的读数（载具，单发）

**一发，未并发**（`wisp.test.exe` 同一时刻只有一枚在跑；跑之前/之后各查一次进程表）。全文 1251454 字节在 `logs/10-cmdwisp-full.txt`。

| 项 | 读数 |
|---|---|
| 命令 | `sh scripts/wisp-cli-staged-tests.sh`（⛔ 不是裸 `go test`：本机 DLL 装载走脚本 stage 的那三枚，裸跑出假红 `0xc0000135`） |
| 起／止 | `19:58:41` → `20:37:56`（墙钟 **2355.00s＝39m15s**，载具自己打的那行：`FAIL github.com/CarlosShao/wisp/cmd/wisp 2355.00s`） |
| 顶层用例 | **630** 枚全跑完（`go test -list` 的名单 630 与结果行 630 相等），红 **3**、其余绿 |
| 判据位 | `-test.timeout=90m` 未触发 |

### 3 枚红，逐枚具名

| # | 用例 | 是不是本腿的 | 读数（原文） |
|---|---|---|---|
| 1 | `TestAC4FocusReturnToPriorWindowGap33r5` (5.26s) | **否**（派单点名的已知红①：前台被一台真窗口占着） | `the panel did not take the foreground on Show: foreground 0x94e0122, panel hwnd 0x7931602` |
| 2 | `TestAC246DevLegIgnoresTheTestTaskInjection` (37.91s) | **否**（派单点名的已知红②：驻留球 16:23 起一直在跑） | `AC#7 RED (desktop state, not our code): a dev Wisp is already running in this session, so this leg handed its activation over and exited. Stop it and re-run.` |
| 3 | `Test219AC2SecondRevokeOfTheSameRuleIsRefused` (55.78s) | **否**（本腿写面 0 交集；见下） | `the reply driver never returned: waited on w.reply 46s (the round's whole bound)` |

**③ 单独取证（不掩盖）**：同一发载具、`-test.run '^Test219AC2SecondRevokeOfTheSameRuleIsRefused$'` ⇒ **PASS (1.70s)**，全文 `logs/11-t219-alone.txt`。整包里它是 55.78s 超时、单跑 1.70s 过 ⇒ **满载下的时序红**（那枚用例的判据自述：40s L2 窗＋15s 第二张卡＝55s，而它只等到 46s 的回程；日志里同时写着"fs.write 那侧已返回、加宽已进 store"）。它与本腿的**结构性**无关：`newWidenRun`／`newReplyHost` 只造 `fs.write`＋审批门，**不建面板宿主**（`approval_always_201_test.go:86-124`，`approval_*.go` 里零 `PanelManager`），本腿的三步全在面板这一侧。⚠ 结论按读数的强度写：**本腿没有把它改成绿，也没有把它改名或跳过**；它是"整包满载时的时序红、单跑即时绿"，需要编排者在低载机时再取一次读数才算定案。

**本腿 6 枚新判据在那发整包里全绿**（从整包日志直取，不是另外单跑的成绩）：
```
--- PASS: Test192PreviewLookupAnswersTheProductPath (0.16s)
--- PASS: Test192PreviewRefusalKeepsTheFallbackPolicy (0.12s)
--- PASS: Test192ResidentAssemblyWiresThePreviewStack (0.01s)
--- PASS: Test192MessageDoorStampsThroughQuoteForModel (0.14s)
--- PASS: Test192MessageDoorIsNotASecondRelay (0.15s)
--- PASS: Test192PreWiringReadingIsStillTheSameFunction (0.00s)
```

⚠ 记录口径：**630 里 3 红＝2 枚点名已知红＋1 枚满载时序红（单跑绿）**；本腿没有新增红，也没有让任何一枚原有的红变绿或变名。

---

## 六、撤销口令

**「192 的接线走 message 门＋资源过滤器两条，预览文字只经 QuoteForModel 盖戳」** ——
撤掉这句要一起撤三件：① `panelAssetLookup` 的预览分支与 `doc.Headers` 直取；② 组合根那一个 `panelPreviewSource`（`withArtifactPreview`）；③ message 门的 `path`＝预览地址这个协议。撤掉之后要重新定：预览走哪条通路、HTML 预览靠什么隔离、以及"预览文字进模型"过哪一道门。
