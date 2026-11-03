# 24-r5 impl — PDF 文本层落地（ledongthuc/pdf），stub 与 DEFERRED(D34-1) 删除

- 腿：票 24 落地腿 24-r5（实现腿）· 编排者派单 2026-11-03
- 前序：24-r4（pdfcpu，无文本抽取 API，停手 `docs/evidence/s1/24r4-stop.md`）；24-r2a（stub 与接口位）
- 依赖：`github.com/ledongthuc/pdf v0.0.0-20261020125853-5e46805f718e`（BSD-3，纯 Go，go.mod 已在 require 块）
- ⛔push：本腿只 commit，推送归编排者

## 1. 写了什么（显式 pathspec 逐枚）

| 文件 | 内容 |
|---|---|
| [internal/tools/document_text.go](../../internal/tools/document_text.go) | `pdfText.ExtractText` 由 stub 改为真实读取器；删除文件头 `DEFERRED(D34-1)` 推迟条目与 stub 注释；`DefaultExtractors`／`docRead.Description()` 文案同步（不再写「PDF 待依赖批准」）；新增 `appendTextWithin` 局部填充助手；新增 import `os`/`unicode/utf8`/`github.com/ledongthuc/pdf` |
| [internal/tools/document_text_pdf_24r5_test.go](../../internal/tools/document_text_pdf_24r5_test.go) | 新测试文件（同包）。只写不读的最小 PDF fixture 生成器 `r5PDF`（对象偏移真实计算，非手填魔数）＋ 6 枚定向测试 |

未动：`document_text_test.go`（见 §4 上报）、`search_content.go`、`bridge.go`、`cmd/wisp/**`、web/search/observe/memory/agent/plugin/thresholds、SPEC/PLAN 契约文字、`go.mod`（见 §4）。

## 2. 实现形状（对应票面四条）

1. **逐页拼接**：`pdf.NewReader(f, size)` → `NumPage()` → `Page(i)`（ledongthuc 页号 1-based）→ `Page.GetPlainText(nil)`（fonts 传 nil，库内 `p.fontCache()` 自建字体缓存）。每页 `strings.TrimSpace` 剥掉库为每个 `BT` 块插入的排版换行，页与页之间以单个 `\n` 连接；`Runs`＝有正文的页数。
2. **上下文预算（C15／D15(3)）**：本读取器只**报** `Truncated`，不自造 spill（spill 层在 loop，`DocDeps.MaxTextBytes` 注释原样保留）。`appendTextWithin` 按剩余预算**部分填充**并在 rune 边界回退（`utf8.ValidString`），因此截断结果永远带正文头部而非空串——这条是实现中被测试抓出来的真缺陷（见 §3）。
3. **扫描件／无文本层→明示、不猜**：`runs==0` 时区分两种零——
   - 有页级错误被记下（`pageErr`）⇒ `pdf: 内容流解析失败：<err>`（不冒充扫描件）；
   - 干净但零正文 ⇒ `pdf: 无可提取文本（扫描件暂不支持 OCR）（疑似扫描件，无文本层）`，措辞沿用 `scannedPDFReason`（OCR 腿与验收腿钉同一句）。
   - `NumPage()==0` 亦走 `scannedPDFReason`。
4. **坏 PDF→error fail-closed**：打不开⇒`无法打开文件`；非法结构⇒`不是有效的 PDF（解析失败）`（含 `invalid header`）；两者都带 `Extracted{}` 空文本、都⛔复用 `ErrUnsupported`、⛔出现 DEFERRED 字样。单页内容流坏不毁整档（跳页继续），与「不静默成功」不冲突：全零时仍按上一条报错。

**Unicode 核查（停手条件判定）**：库在 `page.go:250` 读 `/ToUnicode` CMap，`ps.go` 走 PostScript CMap 解析，并带 `ps_malformed_cmap_test.go`（Type0/Identity-H＋子集字体真实崩溃样本）。⇒ 拿 Unicode 正文的 API 存在，**24-r5 停手条件不触发**，未硬造解析器。

## 3. 测试（`go test ./internal/tools -run Test24R5` 6/6 绿）

| 测试 | 钉住 |
|---|---|
| `Test24R5PDFTextLayerExtracts` | 票面 AC#1「text PDF」；非零 Runs；默认预算下⛔ Truncated |
| `Test24R5PDFMultiPageConcatenatesAndHonoursTheBudget` | 三页按页序出现（索引单调）；`limit=40` 下 ①不超预算 ②必报 Truncated ③**头部仍是第一页正文**（首版实现返回空串，此测试把它判死 → 改为 `appendTextWithin` 部分填充后转绿） |
| `Test24R5PDFScanIsMarkedNotGuessed` | 两页零正文⇒err 含 `scannedPDFReason`＋「扫描件」；⛔DEFERRED；经 `doc.read` 时 `IsError=true`（无静默空结果进模型） |
| `Test24R5PDFBrokenFailsClosed` | garbage／仅头／缺文件三例逐字面点名失败原因；⛔DEFERRED；⛔ErrUnsupported |
| `Test24R5DocReadEndToEndOnPDF` | 分发＋`Result.Origin`＝规范路径（C25 打标缝）；`Lookup` 仍解析为 pdf |
| `Test24R5DeferralMarkersAreGone` | **反向判据＋还原核**：读 `document_text.go` 源文本，`D34-1` 零命中；`DEFERRED(xlsx, D34)`／`DEFERRED(OCR, D34)` 必须仍在（SPEC-12 §5 仍登记，1:1 不新开环） |

Fixture 说明：`r5PDF` 是**字节写入器**（catalog/pages/page×N/content×N/font＋真实偏移 xref＋trailer＋startxref），不是解析器，不读回任何东西；内容流不压缩，故本 fixture 集证明的是「结构＋简单 Type1 文本层」通路。

gofumpt（绝对路径 `D:\work\base\gopath\bin\gofumpt.exe -l`）：两枚自有文件**零输出**（已格式化）。

## 4. 上报编排者（超出本腿写面，未自行处理）

1. **两枚过期测试仍红，且不在本腿写面**：`internal/tools/document_text_test.go:257 Test24R2APDFStubIsErrUnsupportedAndCarriesTheMarker` 与 `:270 Test24R2ADocReadPDFIsARefusalNotASuccess` 断言的正是票面命令删掉的 stub（`ErrUnsupported`＋`DEFERRED(D34-1)`）。删 stub⇒这两枚必然失效；改 `document_text_test.go` 属「其他 tools 文件」，⛔本腿写面。现 `go test ./internal/tools` 整包＝这两枚 FAIL，其余全绿。**需另腿/裁决者把这两枚替换为「PDF 已是真实读取器」的形状**（本腿的 6 枚已顶上覆盖）。⇒ 本腿 commit 的是「新增能力＋新测试绿」，整包转绿依赖这一步，⛔当成整包已完成。
2. **`go.mod` 的 `// indirect` 标记已过期**：ledongthuc/pdf 现为直接依赖（`go.mod:24`）。`go mod tidy` 会去掉注释；`go.mod` 不在本腿写面，未动。构建与测试不受影响。
3. **SPEC-12 §5 的 1:1 尺方向变了**：现读 §5 **无** `DEFERRED(D34-1)` 行（`24-v2-verdict.md:88` 同此读数）。代码侧标记删除后，`d22scan` 的「标记↔登记」双向尺**不再有 D34-1 开环**；若编排者先前已按 A1435/机主批准追加过该行，则该行现在须**撤销或改写为「PDF 已实现，OCR 仍推迟」**——契约文字一步归编排者/机主具名，本腿⛔改。`DEFERRED(OCR, D34)`／`DEFERRED(xlsx, D34)` 保持指向 §5 现有「doc.read xlsx/OCR（D34）」行，未新开编号式标记。
4. **未覆盖的真实样本**：CJK／Identity-H 子集字体的正文抽取仅由库能力与 §2 的源码核查支撑，本腿 fixture 未含此类 PDF。建议验收腿拿一枚真实中文 PDF 复跑 `doc.read`（票面 AC#1 的「text PDF」若要求非拉丁正文，这一格需要外部样本才能勾）。
5. **bridge 接线仍未落**（`BuiltinDocumentEntries` 写面外，沿用 `24r2a-impl.md §4` 记录）。
