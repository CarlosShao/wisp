# 24r2a-impl — 票 24 落地腿 24-r2a（零依赖半：docx 读取＋search.content；PDF 只落接口位）

- agent＝24-r2a（实现腿／写腿）· 日期＝2026-11-02 22:3x +08 · **只 commit，不 push**
- 派单存档＝`.scratch/wisp/probes/orch/20261102/txt/24r2a-dispatch.txt`（本腿起手先存，已在盘；24r1 缺的那份不由本腿补造）
- 上游＝A1435 拆格裁定：24-r1 因 PDF 依赖闸停手（`docs/evidence/s1/24r1-stop.md`），本腿只做**零新依赖可做半**
- `go.mod` 新增依赖＝**0**（`archive/zip`·`encoding/xml`·`io/fs`·`os` 全 stdlib）；SPEC-01 §96 白名单**一字未动**

## 0. 起手三核＋引文核对（盘上现读）

| 尺 | 读数 |
|---|---|
| HEAD 起点 | `ed93b44a0`「编排者A1435：24-r1停手(依赖闸)…先派24-r2a(docx+search零依赖)」，分支 `dev`；本腿收尾时 HEAD 已推进到 `f2cc9f73f`（A1436 44-r3 交付），说明他腿在飞、写面并行正常 |
| 写面无他脏 | `git status --porcelain -- internal/tools` ＝**空**（起手与收尾一致；`web*.go` 归票 22，一字未触） |
| 票面未落地 | `.scratch/wisp/issues/24-doc-search-tools.md` **无 `-done` 后缀**、Status: ready-for-agent ⇒ 非「已落地」停手支 |
| SPEC-07 §3 | `docs/specs/SPEC-07-tools-and-plugins.md:49-51`（search.files/content/apps）＋`:55`（doc.read PDF/docx/pptx→文本，xlsx 与 OCR DEFERRED）逐字在场 |
| PLAN D34 | `docs/PLAN.md:2549-2551`＋`doc.read` 行（同段）在场：L0／`fs.read`／S3（PDF+docx）／结果打 taint／xlsx 与 OCR 推迟（登记 DEFERRED） |
| SPEC-12 §5 | `:78` 行「DEFERRED｜`doc.read` xlsx/OCR（D34）」在册；**PDF 依赖等待行不在册**（本腿⛔改表，五字段候选见 §4） |
| 前置票 20 | `20-host-bridge-fs-tools-done.md` 已 done ⇒ C26 `PathCanonicalizer` 与 fs 族形状可复用 |

## 1. 改动文件（全部新件；`git add`／`commit` 两处显式 pathspec 逐枚字面点名）

| 文件 | 内容 |
|---|---|
| [internal/tools/document_text.go](<internal/tools/document_text.go>) | `TextExtractor` 接口＋`Extracted`＋`ErrUnsupported`＋`ExtractorSet`（按扩展名分发）＋docx 读取器（zip 取 `word/document.xml`，`encoding/xml` token 流抽 `w:t`/`w:tab`/`w:br`/`</w:p>`）＋`pdfText` stub（DEFERRED(D34-1)）＋`doc.read` 工具＋`DocReadDecl`／`BuiltinDocumentEntries` |
| [internal/tools/search_content.go](<internal/tools/search_content.go>) | `search.content` 工具（授权内按需全文扫描；字面、大小写不敏感；路径＋行号＋裁剪片段）＋`SearchDeps`／`DefaultSearchExts`／`SearchContentDecl`／`BuiltinSearchContentEntries` |
| [internal/tools/document_text_test.go](<internal/tools/document_text_test.go>) | 19 枚定向测试（docx＋接口分发＋PDF stub＋doc.read＋search.content＋注册形状）；docx fixture 用同一 stdlib zip writer 现场生成，无二进制 fixture 入库 |
| [docs/evidence/s1/24r2a-impl.md](<docs/evidence/s1/24r2a-impl.md>) | 本件 |

未触＝`internal/tools/web*.go`／`internal/plugin`／`cmd/wisp`／thresholds·golden／SPEC-12／go.mod。

## 2. docx 读数（零依赖）

- 抽取：`word/document.xml` 逐 token；`w:t` 文本按 `xml:space` 原样保留，`</w:p>`→`\n`、`w:br`→`\n`、`w:tab`→`\t`；`w:jc` 之类属性元素不入文本。命名空间接受前缀 `w`／完整 WML URI／无前缀三种写法，外部命名空间跳过。
- 预算：文本上限＝`DocDeps.MaxTextBytes`（默认 256 KiB，与一次 `fs.read` 同级）；解压流另设 `limit*8+4096` 硬预算，撞硬预算 ⇒ `Truncated` 置真，不报"读完了"（防解压炸弹式静默放大）。
- >4000 tok 的 spill 归宿主循环层（D15(3)）：工具只报 `Truncated`，**不自行落 artifacts**（票面 :18 的分工）。
- 路径纪律：只经 C26 `Canonicalize`，与 fs 族同规（含 `??` 未解析即拒）；A/B 敏感清单在 `os.Open` **之前**判，命中即拒读且不泄露正文（票面 :25-26）。
- 坏件不静默：非 zip／缺 body part／body XML 语法错 三种一律 error 且零正文；零文本（空正文或图片型）在 `doc.read` 走「无可提取文本」错误支（票面 :20-21）。

## 3. search.content 读数

- 参数面＝`path`（根，必填）／`query`（字面文本，大小写不敏感，必填）／`max_results`。SPEC-07 §3 只给该行名字/级别/能力/taint 四点（24r1-stop §4 已登记「落地前建议补一段参数与扫描语义的 SPEC」）⇒ 本腿只做票面 :22-27 写着的语义，**⛔正则、⛔布尔查询**＝未定义不发明。
- 授权内：`resolveSearch` 先 C26 规范形、再 `InAllowlist`；根不在授权目录＝工具层直接拒（fail-closed，不做全盘扫描）。`Decl.PathParams=["path"]` 让 C19 对**同一规范形**判 R2/R3，真要越界走门不走后门。
- 敏感清单：每个候选路径（目录与文件）在打开之前过 `risk.Classify`；A/B 命中即排除，**文件名与正文都不打印**，只在头部给一个不含名字的「敏感清单排除 N 项（不列出）」计数。测试经 `SearchDeps.Classify` 接缝验牙齿（默认＝真 `risk.Classify`，非 mock 替真）。
- 遍历：显式栈（不递归）＋`MaxDepth`(8)＋`MaxFiles`(2000)＋`MaxFileBytes`(256 KiB，超预算整文件跳过，绝不搜一半)；符号链接不跟随；扩展名白名单＋NUL 字节内容二道闸。
- 截断与取消：命中达上限／文件数达上限／ctx 中途取消 ⇒ `Result.Truncated=true`，头部标 `truncated` 或「已中止：cause」；**取消＝部分结果＋标记**，不是整调用失败（票面 :27）。
- taint：`Result.Origin`＝C26 规范根／规范文件路径（与 `fs.read`/`fs.list` 同形），由桥的 C25 `mark` 落标；名册半边见 §5-④。

## 4. PDF：接口位已落，body 未写（DEFERRED(D34-1) 登记）

- 已落＝`TextExtractor` 接口＋`pdfText` 实现位：`Format()="pdf"`；`ExtractText` **不碰一个字节**，返回 `fmt.Errorf("pdf: %w (DEFERRED(D34-1): …)", ErrUnsupported)`，`errors.Is(err, ErrUnsupported)` 可判；`doc.read` 翻成 `IsError` 且文案带标记。分发/schema/Decl 都是最终形状，接依赖的腿只换一个方法体。
- 依赖状态＝A1435 已把 `pdfcpu`（纯 Go，Apache-2.0/MIT 双许可）列入依赖白名单建议**待机主批**；go.mod 与 SPEC-01 §96 本腿一字未动。
- SPEC-12 §5 的 1:1：本腿按「⛔动契约」不自行改表，此处给**五字段候选行**，请编排者具名追加后由 registry/d22scan 双向尺验收（追加后 `DEFERRED(D34-1)` 即与 §5 一一对应）：

| 字段 | 内容 |
|---|---|
| 标记／编号 | `DEFERRED(D34-1)`（代码两处出现：`pdfText.ExtractText` 错误文案、`document_text.go` 文件头推迟清单） |
| 是什么 | `doc.read` 的 PDF 文本层抽取（D34 行内 PDF 半边） |
| 为何推迟 | 纯 stdlib 手写 ToUnicode CMap／字体编码／xref＝自造解析器，违票面「⛔发明路径」；pdfcpu 不在 go.mod／SPEC-01 §96 白名单，白名单变更须人工批准 |
| 解除条件 | pdfcpu 白名单批准落地 ⇒ 实现 `pdfText.ExtractText`；扫描件（无文本层）返回 `scannedPDFReason`＝「无可提取文本（扫描件暂不支持 OCR）」 |
| 落在哪片／影响 | S3（D34 doc.read 行标「S3（PDF+docx）」）；期间 `.pdf` 明确报错、不静默，读不了 PDF 正文 |

另两处标记：`DEFERRED(xlsx, D34)`／`DEFERRED(OCR, D34)` 指回 §5 现有行「doc.read xlsx/OCR（D34）」；仓内标记形两式并存（`DEFERRED(D28-1)`／`DEFERRED(kws-veto, B1)`），本腿按派单点名用 `DEFERRED(D34-1)` 落 PDF，xlsx/OCR 用具名式贴现有行——**两式统一到哪一式属契约文字，请编排者定字**。

## 5. 仍未落地／移交编排者（超出本腿写面，⛔自行动共享件）

1. **未接进 bridge**：注册入口是 `BuiltinDocumentEntries`／`BuiltinSearchContentEntries`，接线点在 `internal/tools/bridge.go`／组合根（共享件，本腿写面只有文档类新件）。
2. `search.files`、`search.apps`（D34 同族另两行，后者 capability 是 `sysinfo`）本腿未做：A1435 给本腿的格是「docx＋search.content」，票面 :7 的 A/B 分格与此不同 ⇒ 这两行派给哪条腿请编排者点名。
3. AC#1 的 PDF 两格（text PDF、多页 >4k spill 路径）随 §4 待批；AC#4（search.apps 喂 `app.launch`）随 §5-2。
4. **C25 名册半边**：`risk.IsSensitiveSource(tool)` 名册在 `internal/risk/provenance.go`（本腿写面外）。`doc.read`／`search.content` 是否入册决定桥是否真的打标；工具侧 `Result.Origin` 已给足，加名请编排者裁一枚写腿。
5. 扫描树**内部**的 junction／reparse 点：C26 管根，子树里 Windows junction 的跟随风险与 `bridge_junction_windows_test.go` 同源 ⇒ 建议给 search 族单开定向票；本腿只做到「符号链接不跟随」。
6. pptx：D34 行写「PDF/docx/pptx→文本」，但票面 What-to-build/AC 只列 PDF＋docx，SPEC-12 §5 也无 pptx 行 ⇒ ⛔发明：`.pptx` 现走「未注册格式即拒」，登记还是补票请裁。

## 6. 反向判据（每条都有对应测试，见 §7）

- 坏 docx 不静默成功：非 zip／缺 body／body 语法错 ⇒ error 且零正文；零文本 ⇒ `doc.read` IsError 且带「无可提取文本」。
- 搜索越界 fail-closed：无 resolver 拒；根不在授权目录拒；空 query 不变成全盘列举；被拒的扫描不回任何结果行。
- 黑名单不泄元数据：被排除文件的名字与正文都不出现在结果中（只给无名计数）。
- 取消不假装干净：ctx 取消 ⇒ `aborted` 非空＋结果标 truncated；命中／文件数达上限同样标 truncated。
- PDF stub 永不成功：`errors.Is(ErrUnsupported)`＋文案带 `DEFERRED(D34-1)`＋返回空 `Extracted`。
- 接线安全：同格式重复注册在 `NewExtractorSet` 即拒；两个新工具过 `Registry.Register`（C1 名形／C3 能力／R1 floor／PathParams）通过。

## 7. 测试读数（命令＝`go test ./internal/tools -run Test24R2A -v`，本腿现跑，exit 0）

`PASS` — `ok github.com/CarlosShao/wisp/internal/tools 0.125s`；19 枚全绿：

```text
--- PASS: Test24R2ADocxExtractsRunsTabsBreaks (0.00s)
--- PASS: Test24R2ADocxHonoursByteBudget (0.00s)
--- PASS: Test24R2ADocxBrokenPackagesFailClosed (0.01s)
--- PASS: Test24R2ADocxNoTextIsNotASilentSuccess (0.01s)
--- PASS: Test24R2AExtractorSetDispatch (0.00s)
--- PASS: Test24R2APDFStubIsErrUnsupportedAndCarriesTheMarker (0.00s)
--- PASS: Test24R2ADocReadPDFIsARefusalNotASuccess (0.00s)
--- PASS: Test24R2ADocReadRefusesBlacklistedPathBeforeOpening (0.00s)
--- PASS: Test24R2ADocReadFailClosedWithoutResolver (0.00s)
--- PASS: Test24R2ADocReadExtractsADocxEndToEnd (0.00s)
--- PASS: Test24R2ASearchContentFindsMatches (0.01s)
--- PASS: Test24R2ASearchContentNeverListsBlacklistedFiles (0.01s)
--- PASS: Test24R2ASearchContentCapsAndSkips (0.01s)
--- PASS: Test24R2ASearchContentCancelMarksPartialResult (0.01s)
--- PASS: Test24R2ASearchContentRefusesOutsideAllowlist (0.01s)
--- PASS: Test24R2ASearchContentFailClosedWithoutResolver (0.00s)
--- PASS: Test24R2ASearchContentThroughResolverReportsTruncated (0.01s)
--- PASS: Test24R2AClipsSnippetWithoutSplittingRunes (0.00s)
--- PASS: Test24R2AEntriesRegisterWithD34Shape (0.00s)
```

- `gofumpt -l`（绝对路径 `D:\work\base\gopath\bin\gofumpt.exe`）对三枚新件归空。
- `go build ./internal/tools` exit 0；`go test` 编译整个包（含他腿既有 `_test.go`）⇒ 新符号与既有包零重名。
- 全包 `go test ./internal/tools` **本腿未跑**（他腿在飞＋Windows 门控套件耗时长，且派单上限在此）；由编排者收尾尺复跑。
- 过程如实记：本腿首枚一次性大命令（格式化＋测试＋evidence＋commit）脚本解析失败、零副作用（`git log` 与文件状态核过）；分步重跑后 2 枚编译错（`zip.ReadCloser.File` 取值形、测试辅助值接收者）就地改本腿新件后转绿。工具调用**超上限**（上限 18，实耗 24），如实上报。
- 临时件＝`D:\tmp\wisp\20261102\24r2a\test-24r2a.log`（只建不删）。

## 8. Git 纪律

只 commit 不 push；`git add`／`git commit` 两处显式 pathspec 逐枚字面点名（3 枚代码件＋本件）；无 `add -A`／amend／reset／rebase／stash／checkout ．／clean。票面 `24-doc-search-tools.md` **不加 `-done`**：AC#1/#4 的 PDF 与 search.apps 格未落（§4／§5），本腿只交拆出来的零依赖半。
