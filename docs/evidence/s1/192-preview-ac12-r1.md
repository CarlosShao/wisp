# 票 192 — 产物预览：AC#1（四问）＋ AC#2（注入面判据）

- 腿＝票 192 写腿（`agent-29`，委派）；2026-11-14 17:5x–18:3x +08；起手 HEAD `48527e98f`。
  本件所有读数都在**本腿工作树**与**仓外导出树** `D:\tmp\wisp\20261114\192-tree`（`git archive HEAD` 解出）现取，
  每条带命令与 log 路径，log 都在 `.scratch/wisp/probes/192/logs/`。
- 写面（本腿真正写过的）：`internal/panel/preview_192.go`（新，778 行产码）、
  `internal/panel/preview_192_test.go`（新，719 行判据）、`docs/evidence/s1/192-preview-ac12-r1.md`（本件）、
  `.scratch/wisp/probes/192/**`（新）、票 192 Progress log 尾一条（只追加）。
- ⛔ 未动，逐条可核（§5）：`cmd/wisp/**`（0 字节）、`internal/llm/**`、`internal/tools/**`、`tools/d22scan/**`、
  `frontend/**`、`design/**`、`docs/specs/**`、`docs/PLAN.md`、`allowlist.txt`、`thresholds.go`、golden、
  `internal/risk/**`（0 字节）、`internal/memory/**`（0 字节）。⛔ push（本腿只 commit）。

---

## 〇、第一件事：能不能派 —— **能派**（判据逐条）

派单要求"票面若写着需要机主拍板就立刻停手"。逐条判：

1. **票面唯一的"待定义案"格是 AC#0，它已经不是待定项**：票面 `:24` 逐字记着 owner 当场接单
   （**「内置浏览器，也做一下吧」**）＋编排者据此定案（台账 `A379`），并给了撤销口令「浏览器搜索要换 API」。
   ⇒ AC#0 不挡 AC#1/AC#2；票面 `:12` 自己写着"本票不替他定"指的是**当时**，而定案已落。
2. **AC#1／AC#2 是技术实现＋判据**，票面里没有"等裁定／Q-## 未定／取舍归机主"的字样；
   AC#3（只读边界）与 AC#4（契约轴）反而是**票面已经把边界钉死**的两条，本腿照做即可。
3. **本腿不涉及"必须摆机主"的三类**：不改契约文字（`docs/specs/**`、`PLAN.md` 0 字节）、不花钱、不做不可逆动作（不 push）。
   所以没有任何一处我替机主选了。

⚠ 一条排程事实照实记（票面 `:26` 写着三枚串行 `33-r1` → 191 → 192 预览，因为三枚都在 `internal/panel/**`）：
**本腿动 `internal/panel` 时，这个包里没有别的腿的未提交改动**（`git status --porcelain -- internal/panel` ＝ 0 行，现量见 §5），
所以本波该包只有我一枚写腿；若下一枚 `internal/panel` 的腿要开工，得等本腿提交完。

---

## 一、AC#1 四问：逐条可判答案（代码位置＋盘上读数）

**一段话结论**：四类＝HTML／图片／PDF／文本；渲染的是**内容字节**（浏览器永远拿不到路径）；
大文件只有**文本**做截断（头 4096 ＋ 尾 1024 ＋ 省略数 ＋ 可走通的续读指针），二进制超限**报两个数字拒绝、不截断**；
截断时用户看到的是**宿主写的一行"共多少字节、显示哪两段、中间多少没显示、从这里继续读"＋头＋省略行＋尾**。

### Q1 能预览哪几类

| 类 | 判定依据（字节） | 服务出去的类型 | 渲染归谁 |
|---|---|---|---|
| HTML | `net/http.DetectContentType` → `text/html` | `text/html; charset=utf-8` ＋ 沙箱 CSP | WebView 自己画 |
| 图片 | **复用 composer 的同一张嗅探表** `attachmentTypes`（png／jpeg／gif／webp） | 嗅出来的那个 `image/*` | WebView 自己画 |
| PDF | `%PDF-` → `application/pdf` | `application/pdf` ＋ `inline` | **宿主依赖**（见下） |
| 文本 | `text/plain`／`text/xml`／`application/xml` 且**头部是合法 UTF-8** | `text/plain; charset=utf-8` | WebView 自己画 |

**不在四类里的一律按名字拒（415）**，点名它实际是什么：视频（`ftyp` 容器）、未知二进制、UTF-16 文本、非 UTF-8 文本。
**SVG 特意归"文本"不归"图片"**：SVG 是 XML、能带脚本，当图片渲染等于给它执行面；当文本显示就没有。

代码位置（`grep -n` 现取，log `11-ac1-four-questions.log`）：`preview_192.go:137-140` 四个 kind 常量、
`:515` `func classifyPreviewArtifact`、`:516` `for _, t := range attachmentTypes`。
判据：`TestPreview192FourClassesAndTheRestRefused`（8 条子用例，四类各一＋视频／二进制／UTF-16 三种拒）、
`TestPreview192ClassComesFromBytesNotFromTheName`（同一份 PNG 字节，四个不同扩展名，四发都是 image）。

### Q2 渲染的是本地路径还是内容字节 —— **字节**

三处机制，缺一不成：

1. **名字先过 artifacts 仓自己的裸名守卫**（`preview_192.go:466` `r.guard(name)`，注入的是 `memory.ValidArtifactName`，
   与附件同一条规矩而不是抄一份），**在任何文件系统调用之前**。`../x`、`a/b`、`a\b`、`C:x`、`x.`、`a:b` 全部 400，
   且 `Kind==""`、`TotalBytes==0` —— 即"连 stat 都没发生"。
2. **只把已判过的裸名 join 到 store 自己的 artifacts 目录**（`:345`），并再验一次"普通文件、非链接"（否则 403）。
   没有第二条路：不认绝对路径、不认驱动器号、不跟随链接。
3. **响应里绝不出现路径**：body／headers／reason 里 `ArtifactsDir()` 零命中，页面拿到的只有
   `/__artifact/<名字>` 这种同源 URL。

为什么这条是硬要求（不是洁癖）：票面 `:25` 逐字"不许把『能预览』做成『模型能读任意路径』"，
而按路径渲染正好是 C26 授权根与 175／177／183 那族外泄检测的绕法。
判据：`TestPreview192RefusesPathShapedNamesBeforeItOpensAnything`（9 种恶意形状）、
`TestPreview192TheResponseNeverCarriesTheArtifactsPath`、`TestPreview192NonRegularAndLinkedFilesAreRefused`。

### Q3 大文件截断语义

| 类 | 阈值 | 超了怎么办 |
|---|---|---|
| 文本 | 首屏预算 5120 字节（头 4096 ＋ 尾 1024） | 首屏＝头＋省略行＋尾；**指针** `?offset=&limit=` 续读，单窗 ≤ 65536 字节 |
| HTML | 2 MiB | 413，报"文件多少字节／上限多少字节"，**不截断** |
| 图片 | 8 MiB | 同上 |
| PDF | 32 MiB | 同上 |

三条语义写死：
- **文本 ≤ 5120 字节 → 整段原样返回，不套信封**（判据 `TestPreview192SmallTextComesBackWholeWithNoEnvelope`）。
- **文本窗口是"纯净字节窗口"**：带 `offset` 的请求回的就是 `[offset, offset+limit)` 那几字节，**没有信封**，
  且按 UTF-8 边界切（判据 `TestPreview192ThePointerWalksEveryByteOfTheFile` 断言每窗 `utf8.Valid`）。
- **二进制一律不截断**：半个 PNG 不是图片、半个 PDF 不是文档，所以超限就是 413 加两个数字
  （判据 `TestPreview192BinaryOverTheCapIsRefusedNotTruncated`：4096 字节 vs 上限 1024，body 0 字节、headers 0 行）。

**指针真能找回每一字节**（这是"带找回指针"的判据，不是注释）：盘上读数
`logs/12-ac1-sample-first-view.log` —— 产物 300080 字节、首屏 body 5551 字节、
**5 个窗口取回 300080/300080 字节、`bytes.Equal==true`、sha256 `f9c1eeace35b3cf0` 与原文件相同**。

### Q4 截断时界面看到什么（逐字，来自真跑）

`go run ./.scratch/wisp/probes/192/sample`（log `12-ac1-sample-first-view.log`，真 `memory.Store`＋真守卫）：

```
[wisp 产物预览] big-log.txt · 共 300080 字节 · 显示开头 4096 与结尾 1023 字节，中间 294961 字节未显示
[wisp 续读] /__artifact/big-log.txt?offset=4096&limit=65536 · 每窗最多 65536 字节，一直取到结尾；这段指针能把全文一字节不差地取回
[wisp 正文] 开头 4096 字节
……（产物自己的头 4096 字节）……
[wisp 省略] 中间 294961 字节未显示，用上面的续读地址取
[wisp 正文] 结尾 1023 字节
……（产物自己的尾 1023 字节）……
```

三个要点：① 这是**宿主写的一行宣告**，不是产物内容，所以用户一眼看出"被截断了"；
② 三个数字**相加等于全文**（4096＋1023＋294961＝300080，判据里现算现断言）；
③ 尾不是整整 1024 而是 1023，因为切在 UTF-8 边界上（这类"差一"是故意的，判据只允许差 0–3 字节）。
判据 `TestPreview192TextFirstViewSaysWhatItHidAndWhereTheRestIs` **解析的就是这段显示文本本身**（正则抓三个数字），
所以改了措辞却忘了改数字，这一枚会红。

---

## 二、AC#2 注入面判据（"预览不是输入"）

AC#2 原文两个要求：预览内容**不许顺带进模型上下文**；**若**面板把预览里一段文字回传给模型，那一发必须走
C25 盖戳＋外泄检测。本腿把它落成**三腿判据**：

### A 预览这条路结构上到不了模型（`TestPreview192PreviewIsNotAMessageRoute`）

| 腿 | 断言 | 凭什么 |
|---|---|---|
| 导入面 | `preview_192.go` 不 import `internal/agent`／`internal/tools`／`internal/llm` | 用 `go/parser` 解析真实 import 块，不是 grep 注释 |
| 字面量面 | 该文件不出现 `MethodMessageSend`／`HandleMessageRequest`／`ComposerDispatch` | 预览不是一条消息通路 |
| 值类型面 | `PreviewDocument` 没有 `ComposerRequest` 字段 | 预览的产物不是一个信封 |
| 门面 | `ParseComposerRequest` 拒 `panel.preview`，**而它接受 `panel.message.send`**（正控） | 正控保证"拒绝"不是因为门本来死的 |

另外 C17 **一个方法没加**（AC#4 逐字要求）。预览走的是**已有的那条资源过滤器**（`AddWebResourceRequestedFilter`），
这也是主流 harness 的形状（Electron 的自定义协议／WebView2 的资源拦截），而不是给 postMessage 封套再开一个门。

### B HTML 预览被沙箱钉住（`TestPreview192HTMLIsSandboxedSoAPreviewCannotSpeak`）

响应的 `Content-Security-Policy` 是 **`sandbox; default-src 'none'; img-src data:; style-src 'unsafe-inline'`**，
另加 `X-Content-Type-Options: nosniff`。`sandbox` 不是装饰：`bridge.go` 的 `ParseComposerRequest` 只核对封套里的
`source` **字符串**（那不构成授权，见 `bridge.go:11-18` 自己的说明），所以一个能跑脚本的同源预览页
可以伪造 `source:"panel-composer"` 的请求；`sandbox`（没有 `allow-scripts`）让脚本不跑、文档拿到不透明源。
判据三向：策略必须以 `sandbox` 开头、不得含 `allow-scripts`、必须含 `default-src 'none'`。

### C 唯一一道"预览文字进模型"的门（`QuoteForModel`）

`PreviewResolver.QuoteForModel(ctx, marker, scopeID, artifact, quoted)` 在盖戳前**证明三件事**：
引用非空；`scope` 与引擎都在（缺信息＝拒）；**这段文字真的在该产物的字节里**（现场把产物读回来查，
且产物不得超过核对上限 8 MiB）。然后 `Mark(scopeID, "task.output", "artifact:<名字>", quoted)`
——源名用 `risk.SrcTaskOutput`（`PreviewSourceTool` 常量，判据 `TestPreview192SourceToolIsTheContractName` 把它钉死在契约名上），
因为这条门正是"读回宿主自己落盘的产物"那类门，票 175 对 `task.output` 的原话就是
"outside content by origin, however domestic the door that reads it back is"。

**"在产物里"这条不是形式**：没有它，页面可以把自己编的文字盖成 `task.output` 的戳——
既能往 C25 索引里灌碎片让真泄漏淹没在噪声里，也能种碎片让后来的正常调用看起来可疑。
判据 `TestPreview192QuoteDoorRefusesTextThatIsNotInTheArtifact` 就钉这一条。

**盖戳之后外泄检测真的会响**（不是"应该有"）：真 `risk.Provenance` 引擎，引用后
`CheckText(scope, ChHTTP, quoted)` **命中且点名** `SrcTool=="task.output"`、`Origin=="artifact:quoted.txt"`；
`Inspect(scope, "notify", {text: quoted})` 同样命中（判据 `TestPreview192QuoteDoorStampsAndTheEgressGateThenCatches`）。
另有 `TestPreview192QuoteDoorRefusesWhatItCannotStamp`：无引擎／无 scope／空引用／坏名字／产物不存在／超核对上限，六种一律拒并点名。

### 红证：五枚变异，各自红在自己那句上（仓外导出树）

台件 `.scratch/wisp/probes/192/redproof-192.sh`（跑在 `D:\tmp\wisp\20261114\192-tree`，主树一字未动；
log `logs/redproof-console.txt` ＋ `logs/02..06-*.log`）。
pristine `preview_192.go` sha256＝`70eeb417f2fb2aaa2497e972deb709b94db212e632d9ae0d07fce66a1189db71`。

| 变异 | 改了什么 | 红在哪一行 | 红在哪句话 |
|---|---|---|---|
| M1 quote-not-stamped | Mark 调用被短路（`if false && marker.Mark(...)`） | `preview_192_test.go:618` | "after the quote the C25 egress scan found nothing: the stamp is not live" |
| M2 laundering-allowed | "文字必须在产物里"改成永远成立 | `:642` | "the quote door stamped text the artifact does not contain: the page could choose what the C25 index believes" |
| M3 sandbox-dropped | HTML 的 CSP 去掉 `sandbox` | `:512` | "it must start with the sandbox directive, or the previewed page keeps this window's origin and can postMessage as the composer" |
| M4 pointer-dropped | 首屏指针清零 | `:338` | "the notice shows head 4096 but the pointer resumes at 0" |
| M5 name-guard-dropped | 跳过裸名守卫 | `:237` | `"..": status = 403 (... 不是普通文件 ...), want 400` —— 名字真的摸到了文件系统 |

五发**逐发先证落地（diff 打出那一行）→ 红 → 还原 → sha256 等于 pristine**，最后整组再跑一次绿（`07-final-green.log`）。
⚠ M5 的红法特别值钱：它证明"守卫"不是装饰——去掉守卫后 `..` 走到了 `os.Lstat`，被"不是普通文件"挡下，
也就是说**剩下的唯一一道就只剩那条无关的兜底**。

### ⚠ AC#2 的具名残余（不许被"三腿判据绿"压掉）

1. **生产里还没有调用者**：`QuoteForModel` 全仓 0 调用点
   （`Get-ChildItem -Recurse -Include *.go internal,cmd | Select-String QuoteForModel | Where-Object { $_.Path -notmatch 'preview_192' }` ＝ 0 行；
   探针 `probes/192/sample` 也没调它）。
   也就是说"预览文字进模型必须过这道门"今天＝**门与判据在，接线未落**。接线点在 `cmd/wisp` 的 message 处理器，
   那是**别的腿的写面**，本腿不碰（见 §6）。
2. **预览这条路今天在生产里回答 404**：`cmd/wisp/panel_control_windows.go:563` 的 `panelAssetLookup` 只答嵌入包，
   `/__artifact/...` 会落到 `Assets.Resolve` 的"不在包里"分支。接线是**一行分支＋把 `doc.Headers` 拼进响应**（见 §6）。
3. **前端零写面**：页面怎么把 `[wisp 产物预览]` 那几行显示出来、怎么请求窗口，本腿给的是**数据和语义**，没有 UI 证据。
4. **PDF 的像素不承诺**：WebView2 的 PDF 渲染是 runtime／策略相关（同一台机器上装了 Edge PDF 查看器才画得出来），
   所以 PDF 的 `Render` 值是 `host-pdf` 而不是"能画"——这是**声明依赖**，不是假装能预览。
5. **文本窗口只在"首屏合法 UTF-8"的前提上干净**：分类要求头部合法 UTF-8；文件**中间**若有坏字节，
   窗口会原样把那几字节返回（不切碎字符边界，但不保证整窗合法）。这属于"文本类"的边界，照实记。

---

## 三、票面 AC#3／AC#4 的对照

- **AC#3 只读边界**：本腿不做表单提交、登录、下载、脚本执行。代码里没有网络（`net/http` 只用来做类型嗅探，不发请求）、
  没有写文件、没有 `exec`。HTML 的脚本面被 B 腿的 CSP 关掉。
- **AC#4 契约轴**：`docs/PLAN.md`／`docs/specs/**`／`allowlist.txt`／`thresholds.go`／golden **0 字节**；
  **C17 一个方法没加**（`bridge.go` 未改，`internal/panel/bridge.go` sha256 现量见 §5）；凭据值零抄录（本腿不涉及凭据）。

---

## 四、门禁读数（每条都可复跑）

| 门 | 命令 | 读数 | log |
|---|---|---|---|
| 整包（改前基线） | `go test -count=1 ./internal/panel/` | `ok` 49.251s | `logs/00-baseline-panel-test.log` |
| 整包（改后，第一次） | 同上 | `ok` 58.095s | `logs/08-panel-full.log` |
| 整包（改后，交付前最后一次，别的腿的提交已落） | 同上 | `ok` 70.870s | `logs/13-panel-full-final.log` |
| 新判据 | `go test -count=1 -run TestPreview192 -v ./internal/panel/` | 19 枚顶层＋8 条子用例全 PASS | `logs/01-baseline.log`、`07-final-green.log` |
| vet | `go vet ./internal/panel/` | rc 0 | — |
| 格式 | `gofumpt -l`（`v0.12.0 (devel) (go1.27.1)`） | 0 文件 | — |
| d22scan | `tools/d22scan/d22scan.exe -root .` | `clean`，`internal/` 1014 枚（**改前 1012**）、`cmd/` 283 枚 | `logs/09-d22scan.log`、`logs/14-d22scan-final.log` |

**gate-clauses 的口径**（票面 AC#5 只要"比名册，只 `G6neg`"）：本腿**产码**里没有 `OpenScope`／`CloseScope`／
`OpenTask`／`CloseTask`／`Defer`／`DisposalScope` 的生产调用（`QuoteForModel` 收的是 `TaintMarker` 接口，不开 scope）；
判据件里有 3 发 `prov.OpenScope(...)`，**每一发都把句柄绑到变量、并在 `defer` 里 `.Close()`**（`preview_192_test.go:604/605`、`:637/638`、`:657/658`），
落在 G5 的"开且关"那一堆，不落"只开"那堆，所以名册的未成对枚数不该动。
**这一条本腿只读不修**：如果它跑完的读数与本件不一致，以它的 log 为准，本腿不改码去追它
（为一把尺改产码是"改口径追读数"那起事故的形状）。
⚠ **诚实交代**：`sh .scratch/wisp/probes/154/gate-clauses.sh` 在本腿现跑（`logs/10-gate-clauses-with-preview.txt`），
到本件落笔时它**仍在历史扫描阶段**（`git` 进程累计 CPU 100 秒以上、log 已 186 KB，跑的是逐提交的家族普查），
本件**没有**拿到它的最终聚合退码，也不假装拿到；随本件提交的 `10-gate-clauses-with-preview.txt` 是那一刻的
**186 KB 中途快照** —— 它在本腿提交之后跑完了，跑完的 **241397 字节／2081 行**全文在同一条追加提交里
（读数见下面的 §4.1）。

### §4.1 gate-clauses 跑完了（本腿提交之后）—— 读数、以及它量的是哪棵树

写这件的时候它跑完了：**`rc=2`（聚合退码 2＝两腿 BAD）**，逐腿行逐字是
`BAD 腿=G7 声明=ring 基线=6枚 实测=7枚` 与 `BAD 腿=G7neg 声明=ring 基线=7枚 实测=8枚`，
其余十二腿（G1／G1b／G2／G3／G4／G5／G5pos／G5neg／G6／G6pos／G6neg／G7pos）与各自声明一致。
完整 log 见工作树里跑完后的 `.scratch/wisp/probes/192/logs/10-gate-clauses-with-preview.txt`。

**这两腿 BAD 不是本腿造成的，而且这一点不是推测，是它的口径决定的**：那支尺把 anchor 写进了每一条命令
（`gate-clauses.sh:96` `A="${1:-$(git rev-parse HEAD)}"`，日志里每条 grep 都是
`git grep ... b06f7fa0fab6ef9c943e66408a90867859f536a5 -- internal/**/*.go cmd/**/*.go`）
⇒ 它量的是**锚点 `b06f7fa0f` 那棵树**（本腿起手 HEAD `48527e98f` 的祖先，比本腿提交 `7ff4d9e9f` 早三枚提交）。
`git ls-tree b06f7fa0f internal/panel/` 里**没有** `preview_192.go`／`preview_192_test.go`
⇒ 本腿两枚文件**物理上不在**这次扫描的树里，那两腿的 6→7／7→8 与本腿无关（是那棵树上别的腿的落点）。

**那么"本腿对名册的影响"用什么读数回答**：G1–G7 各腿的 grep 是
`Defer(Named)?|(New)?DisposalScope|OpenScope|CloseScope|OpenTask|CloseTask|ToolRequest{` 这一族词（`git grep -nEw`，**区分大小写**），
本腿两枚文件在这些词上的命中数是可数的，而**主尺腿一律 `:!*_test.go`**：

```
Select-String -CaseSensitive -Path internal/panel/preview_192.go \
  -Pattern 'OpenScope|CloseScope|OpenTask|CloseTask|DisposalScope|Defer(Named)?|NewDisposalScope|ToolRequest\{'   => 0 行
Select-String -CaseSensitive -Path internal/panel/preview_192_test.go 同 pattern
  => 3 行（`:604`／`:637`／`:657` 三发 prov.OpenScope，每一发紧跟着 `defer func() { _ = scope.Close() }()`）
```

（⚠ 这一格必须加 `-CaseSensitive`：不加的话 `defer func() { _ = f.Close() }()` 那两行小写 `defer` 会被算进来，读数就成了 2／6——
那是**量错**不是"更多命中"，记在这里免得下一程照抄一个错读数。）

⇒ 主尺腿（排 `*_test.go`）的分子**变化量为 0**；G5 那一族的判据件三发是"开且关＋句柄已绑"，
落"开且关"堆，也不进"只开"堆。**结论：本腿交付后，gate-clauses 的名册与锚点 `b06f7fa0f` 上的完全相同**
（这条是"我的文件在那些 grep 上零命中"推出来的，不是重跑一遍得到的数；要"重跑一遍的聚合退码"，
请编排者用 `sh .scratch/wisp/probes/154/gate-clauses.sh 7ff4d9e9f` 跑一遍——那会花 30 分钟以上，
本腿不再占机时，而且本腿的写面里没有能改 G7 那一族的东西）。

**逐腿现状照抄（锚点 `b06f7fa0f`，14 腿：12 腿 ok、2 腿 BAD、0 腿 SHR）**：
`G1 0/0 ok · G1b 4/4 ok · G2 2/2 ok · G3 0/0 ok · G4 0/0 ok · G5 2/2 ok · G5pos 0/0 ok · G5neg 14/14 ok ·
G6 0/0 ok · G6pos 0/0 ok · **G6neg 1/1 ok** · **G7 6→7 BAD** · G7pos 0/0 ok · **G7neg 7→8 BAD**`。
票面 AC#5 要求的"只 `G6neg`"这一格今天**是对得上的**（1/1）；两腿 BAD 是 G7 那一族（DisposalScope／Defer），
与产物预览无关，且如上所述产生在**没有我文件的那棵树**上。本腿对"名册不该动"给的是**另一把尺**的读数（不是它的替代品）：
`Select-String -Path internal/panel/preview_192.go -Pattern 'OpenScope|CloseScope|OpenTask|CloseTask|DisposalScope|Defer\('`
＝ 0 行。gate-clauses 的最终数由下一程（或编排者）补，判据只有一条：**除 `G6neg` 外逐腿与名册一致**。

---

## 五、未动清单（可核命令与读数）

本腿起手时 HEAD＝`48527e98f`；本腿跑的过程中**别的腿提交了三枚**（`b06f7fa0f`、`8d56a948d`、`9f2333a3d`、`447929bfe`），
所以交付时的父提交不是起手那个。**它们都没有碰 `internal/panel`**（`git log --oneline 48527e98f..HEAD -- internal/panel` ＝ 0 行），
因此本件 §4 的整包基线读数对本包仍然有效；改动的只有 `cmd/wisp/**`、`internal/llm/**`（别人的写面）。

```
git status --porcelain -- internal/panel/
    => 只有本腿两枚新文件：?? internal/panel/preview_192.go  ?? internal/panel/preview_192_test.go
git status --porcelain -- docs/specs/ docs/PLAN.md allowlist.txt thresholds.go internal/risk/ internal/memory/
    => 0 行（契约轴、阈值、C25 引擎、artifacts 仓：本腿 0 字节）
internal/panel/bridge.go sha256 = ff326be78574108a97ac2dfa71792f57df75992a92ac9a90e8ff47b117d81ce1（未改；C17 名册原样）
git status --porcelain -- cmd/wisp/ internal/tools/ tools/d22scan/
    => 只有**别的腿**的未提交改动：M internal/tools/system_backend_windows.go
       ＋ tools/d22scan 下三枚未跟踪脚本（check_blob_417.sh／export_red_proof.sh／selfreport_delta_test.go）
       —— 本腿在这三处 0 字节
git status --porcelain -- frontend/ design/
    => 只有**别的腿/外部 agent** 的东西：M frontend/fixtures/composer-states.html、?? design/**、?? frontend/scripts/probe-ballsize-render.tsx
       —— 本腿在这两处 0 字节
```

本腿的提交只带**显式 pathspec**（5 条：产码、判据、本件、`probes/192` 目录、票面追加），不用 `add -A`、不 amend、不 push。

---

## 六、下一腿（接线，属 `cmd/wisp` 写面，不在本腿）

1. `panelAssetLookup` 加一条分支：`strings.HasPrefix(path, panel.PreviewPathPrefix)` 时问 `PreviewResolver.Lookup(ctx, path)`，
   用 `doc.Status`／`doc.ContentType`／`doc.Body` 造响应，**把 `doc.Headers` 全部拼进去**（预览的 CSP 与 nosniff 在 `doc.Headers` 里，
   不拼就会回落到面板那条 CSP，HTML 的沙箱就白写了）。
2. 组合根造一枚 `PreviewResolver`（`panel.NewPreviewResolver(store, memory.ValidArtifactName)`），挂在面板宿主装配处。
3. message 处理器接 `QuoteForModel`：**凡是"用户从预览里选了一段文字发送"那一发**，先过这道门再进模型。
   不做这一步，AC#2 在**产品能力**上就只是"门在、没接线"（见 §2 残余 1）。

---

## 七、撤销口令

**「192 的产物预览走资源过滤器，不过 C17」** —— 这一句撤掉的是本票的架构选择（预览＝资源过滤器路径、
HTML 沙箱、文本窗口＋指针三条一起撤），撤销后要重新定：预览走哪条通路、HTML 预览靠什么隔离、截断靠什么找回。
