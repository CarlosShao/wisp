# 24-v1-verdict — 票 24 验收腿 24-v1（非实现者裁决；裁决者≠实现者）

- agent＝24-v1（验收腿／只读＋本证据件）· 日期＝2026-11-02 22:5x +08 · **只 commit 本件，不 push**
- 被验收交付＝`8523211a0`（24-r2a 实现腿，票 24 零依赖半：docx＋search.content；PDF＝stub）
- 尺的边界（如实）：本腿⛔改码⇒**未做 mutation 反转**（不摘保护看测试变红）；反向判据以「测试逐枚真断言 error＋零正文＋零泄露」＋「保护调用点在 Execute 主路径（非 mock 替身）」双读代替，逐条列 §4。
- 读数日志＝`D:\tmp\wisp\20261102\24v1\test-dispatch-regex.log`＋`test-full-19.log`（只建不删）

## 0. 起手三核（盘上现读）

| 尺 | 读数 | 判定 |
|---|---|---|
| HEAD 现盘最新 | `a11e63bee`（44-r4）＞`f2cc9f73f`；`8523211a0` 为其祖先（`git log` 序 a11e63bee→65768beca→8523211a0→f2cc9f73f） | **PASS** |
| 索引干净 | `git status --porcelain` 中 tracked 改动＝**0**（1204 行全为 `.scratch` 历史未跟踪噪声）；`-- internal/tools`＝**空** | **PASS** |
| 票面在场 | `.scratch/wisp/issues/24-doc-search-tools.md` 无 `-done` 后缀、AC 复选框未勾（与「半格交付」一致）；前置 `20-host-bridge-fs-tools-done.md` 已 done | **PASS** |

## 1. 真跑读数（★只 scope `./internal/tools`，⛔ ./...）

| 命令 | 读数 | 判定 |
|---|---|---|
| `go test ./internal/tools -run 'Document\|Doc\|SearchContent\|TextExtractor' -count=1 -timeout 20m -v` | 命中 17 枚＝r2a 15 枚（正则漏掉 3 枚，见下行）＋外来 1 枚（`Test172DocCopyThatLiesReddensTheRuler`，同包他票测试，**也 PASS**）＋表驱动子项 3；**FAIL=0**，`ok … 0.210s`，exit 0 | **PASS** |
| 补跑 `go test ./internal/tools -run 'Test24R2A' -count=1 -timeout 20m -v`（派单正则不覆盖 `ExtractorSetDispatch`／`PDFStub…`／`ClipsSnippet…` 三枚，为凑齐题面 19 枚所加，仍只 scope internal/tools） | **19/19 PASS、FAIL=0**，exit 0；枚名与 24r2a-impl §7 逐一对得上 | **PASS** |

## 2. 还原核

- `git log 8523211a0..HEAD -- <三枚代码件>`＝**空**（后两枚编排者/他腿提交未碰本交付）。
- 盘面对提交：**blob hash 三枚全同**（document_text.go `3425ed127002`＝disk；document_text_test.go `e12de4287822`＝disk；search_content.go `7f8bd5865ade`＝disk，`git hash-object` 现算）；worktree 无脏。⇒ §1 跑的就是 `8523211a0` 交付的码。

## 3. 逐格裁决

### 3.1 docx 读取（文本抽出、坏档报错）

| 格 | 读数（源码＋测试双读） | 判定 |
|---|---|---|
| 文本抽出 | `extractDocxXML` 逐 token：`w:t` CharData 累加、`</w:p>`/`w:br`→`\n`、`w:tab`→`\t`；命名空间三式（`w` 前缀／完整 WML URI／无前缀）外一律跳过；`Test24R2ADocxExtractsRunsTabsBreaks` 断言 5 runs＋顺序＋`w:jc` 属性值不入文本 | **PASS** |
| 坏档报错（不静默） | 非 zip／缺 `word/document.xml`／body XML 语法错 ⇒ 三式各 error＋**零正文**（测试逐枚断 `err!=nil`＋`wantSub` 文案＋`got.Text==""`）；零文本走 `doc.read`「无可提取文本」IsError 支（票面 :20-21），测试 `…NoTextIsNotASilentSuccess` 钉住两半边 | **PASS** |
| 预算/解压炸弹 | 文本上限 `MaxTextBytes`（默认 256 KiB）；解压流 `limit*8+4096` 硬预算，撞顶⇒`Truncated=true`（`…HonoursByteBudget` 断 12 字节预算不越界且置位）；spill 归 D15(3) 宿主循环层、工具不自行落盘＝票面 :18 分工 | **PASS** |
| 路径纪律 | 只经 C26 `Canonicalize`；无 resolver 拒（fail-closed）；`??` 未解析拒；A/B `risk.Classify` 在**打开之前**判（`Execute` 主路径 :458 先于任何 open），命中拒读且测试验「leak 零正文」 | **PASS** |
| 格式分发 | `NewExtractorSet` 重复格式注册即拒（接线 bug 不能静默遮蔽）；xlsx 无 extractor＝拒绝而非返回 shared-strings 裸 XML | **PASS** |

### 3.2 search.content（allowlist／预算／taint）

| 格 | 读数 | 判定 |
|---|---|---|
| C26 越界拒（fail-closed） | `resolveSearch`：无 resolver 拒⇒`Canonicalize`⇒`??`拒⇒**`InAllowlist` 不在授权根即拒**（票面 :23「within allowlist」在工具层自闸，不等门）；`Decl.PathParams=["path"]` 让 C19 对同一规范形判 R2/R3（越界走门不走后门）；测试 `…RefusesOutsideAllowlist` 断 IsError＋**不回任何结果行**（无 SECRET），`…FailClosedWithoutResolver` 另断空 query 不变全盘列举 | **PASS** |
| 三重预算 | `MaxDepth`(8)／`MaxFiles`(2000)／`MaxFileBytes`(256 KiB，**超预算整文件跳过、绝不搜一半**)，外加 `MaxResults`(50) 命中闸；`…CapsAndSkips` 三式分钉＋`…ThroughResolverReportsTruncated` 断 Truncated 置位＋头部含 `truncated`＋片段 ≤200B 不切 rune | **PASS** |
| 敏感清单不泄元数据 | 每个候选（目录与文件）在 `os.Open` **之前**过 `risk.Classify`，A/B 命中⇒排除、名字与正文都不打印，只给无名计数「敏感清单排除 N 项（不列出）」；`…NeverListsBlacklistedFiles` 经 `SearchDeps.Classify` 接缝验牙齿（默认＝真 `risk.Classify`，非 mock 替真） | **PASS** |
| 取消⇒部分结果标 truncated | 显式栈遍历＋每帧/每文件查 `ctx.Err()`；取消⇒`aborted` 非空＋结果 `Truncated=true`，**不是整调用失败也不是假装扫完**（票面 :27）；`…CancelMarksPartialResult` 钉住 | **PASS** |
| taint | 工具侧 `Result.Origin`＝C26 规范形（与 fs 族同形，测试断 `Origin==canon`）；**盘上现读 `internal/risk/provenance.go:100-101`：`SrcDocRead`／`SrcSearchContent` 已在 SPEC-06 §5 名册 `sensitiveSourceTools` 内**⇒桥 C25 真打标；即便不在册，:579 未入册源亦 fail-closed 记为 sensitive。24r2a-impl §5-4「名册待裁」半项**现盘已足**，可销 | **PASS** |
| 未发明 | 字面＋大小写不敏感单子集；⛔正则⛔布尔（SPEC-07 §3 该行只给名/级/能力/taint 四点，未定义即停）；显式栈、符号链接不跟随 | **PASS** |

### 3.3 零第三方依赖

| 尺 | 读数 | 判定 |
|---|---|---|
| `git show 8523211a0 -- go.mod go.sum` | **空 diff**（题面所期一致） | **PASS** |
| imports（三件全 dump） | 仅 stdlib（archive/zip・context・encoding/json・encoding/xml・errors・fmt・io・io/fs・os・sort・strings・path/filepath[仅测试]・testing）＋同模块 `internal/risk`；**零外部模块** | **PASS** |
| pdfcpu 在场性 | `go.mod` require 无 pdfcpu；SPEC-01 :96 白名单无 pdfcpu⇒「依赖未批＝stub」前提成立 | **PASS** |

### 3.4 纪律

| 尺 | 读数 | 判定 |
|---|---|---|
| stat 与题面符 | 4 files changed / 1681 insertions / **0 deletions**；name-status 全 `A`＝document_text.go(512)＋document_text_test.go(637)＋search_content.go(416)＋24r2a-impl.md(116)＝题面「4 新件」逐符 | **PASS** |
| 写面 | 全部落 `internal/tools` 文档类新件＋证据件；**未触** `internal/tools/web*.go`／`internal/plugin`／`cmd/wisp`／thresholds·golden／SPEC-12／go.mod；票面未改、无 `-done`（半格如实）；只 commit 不 push（`8523211a0` 后无 amend/reset 痕迹，历史只增） | **PASS** |
| 形状禁令 | 新件无裸 `go func(`；无 `risk.PathResolver` 外 `filepath.Clean\|Abs`（新源件根本未 import path/filepath，扩展名切分用字符串且只读）；emoji 禁段（U+2200–22FF·2600–27BF·2B00–2BFF·FE0F·1F000–1FAFF）三件 ripgrep＝**0 命中**；`gofumpt -l` 三件归空 | **PASS** |
| 证据自述与盘符 | 24r2a-impl §7「19 枚全绿」＝本腿复跑 19/19 复现；§0 起手 HEAD 自述（ed93b44a0→f2cc9f73f 漂移）与提交序相符 | **PASS** |

### 3.5 反向核（题面两问）

- **坏 docx 不静默成功**：`Test24R2ADocxBrokenPackagesFailClosed` 三式各断「必 error＋文案点名＋零正文」；`…NoTextIsNotASilentSuccess` 断零文本必 IsError。真跑绿。⛔ 未做摘保护 mutation 反转（越权改码），以主路径调用点读码代替——拒绝调用在 `Execute`/`ExtractText` 主干而非测试脚手架内，无「mock 替真」形状。
- **search 越界 fail-closed**：`…RefusesOutsideAllowlist`（根不在授权＝拒且零回显）＋`…FailClosedWithoutResolver`（无 resolver 拒＋空 query 不列举）。真跑绿。

**小观察（不判 FAIL，随接线票处置）**：① `run()` 对单个不可读目录 `continue` 静默跳过（权限边界非正确性，建议接线票决定是否给 aborted/告警位）；② docx 抽取中途 ctx 取消返回 error 而非 truncated 部分结果——票面 :27 只约束 search 族，合规；③ `doc.read` 的 `max_chars` 语义实为字节数（名/述微差，无害）；④ config 驱动的 per-tool timeout 与 >4k spill 落 artifacts 两格属宿主层，随 bridge 接线票验收。

## 4. PDF 处置裁决（明确）

**① docx＋search 部分是否 PASS＝是。** §1–§3 全部格 PASS：19/19 真跑 FAIL=0、还原核全同、越界/坏档/泄露三反向有牙、零新依赖、写面纪律无违。

**② 24 能否在 PDF 单列/待批下结案＝「零依赖半」可结，整票不可结。**
- 同意结的格＝24-r2a 交付的 `doc.read`(docx 实)＋`search.content` 半（含 PDF 接口位＋stub 拒答形状），验收即本件。
- **整票 24 不得加 `-done`**：票面 AC#1 的 PDF 两格（text PDF／多页 >4k spill）、AC#3 `search.files`、AC#4 `search.apps` 喂 `app.launch`、以及 §5-1 的 bridge/组合根接线均未落——这些与 pdfcpu 白名单无关（search.files/apps、接线本就零依赖），**不得混进「等 PDF」一并挂起**；依 A1435 拆格裁定由编排者点名续腿。
- PDF 半格处置维持 24-r1/24-r2a 的既有裁定：依赖闸属「白名单变更须人工批准」，pdfcpu 提案已摆机主（A1435）；stub 永不成功（`errors.Is(ErrUnsupported)`＋文案带标记＋空 `Extracted`，两枚测试钉死），`.pdf` 在读路径是明确报错而非静默——符合铁律「不为省事自造解析器」与票面「⛔发明路径」。**此格单列待批，不阻塞其余半格结案。**

**③ PDF 残留登记（结案前置，缺一即 1:1 双向尺不闭）：**
1. **SPEC-12 §5 追加五字段行**＝逐字采用 24r2a-impl §4 候选行：标记＝`DEFERRED(D34-1)`（代码两处出现：错误文案 :337＋文件头 :29/:325）；是什么＝doc.read 的 PDF 文本层抽取（D34 行内 PDF 半）；为何推迟＝pdfcpu 不在 go.mod/SPEC-01 §96 白名单，白名单变更须人工批准，自造解析器违票面；解除条件＝pdfcpu 白名单批准落地⇒实现 `pdfText.ExtractText`，扫描件回 `scannedPDFReason`；落哪片/影响＝S3，期间 `.pdf` 明确报错不静默。**§5 是契约文档，追加须人工批准⇒摆机主/编排者具名落表，本验收腿⛔自改。**
2. **新开续票（登记在票面池，非本腿写面）**：建议一枚「24b：pdfcpu 白名单批准落地后实现 PDF 文本层＋text-PDF/多页 spill/扫描件错误三支 AC」，解除条件＝白名单批准；另两枚零依赖续票位＝`search.files`、`search.apps`（＋`app.launch` 喂格）；bridge/组合根接线一枚（含 §3.5-④ 两观察随票验）。
3. **标记两式统一**（`DEFERRED(D34-1)` 编号式 vs `DEFERRED(xlsx, D34)` 具名式）属契约文字，维持 impl §4 上报，由编排者定字；本腿不裁。

## 5. 结论

- 逐格：**§3.1 docx＝PASS｜§3.2 search.content＝PASS｜§3.3 零三方依赖＝PASS｜§3.4 纪律＝PASS｜§2 还原核＝PASS｜§1 真跑 FAIL=0**。
- 结案面：**同意结 24（docx＋search 部分，即 24-r2a 零依赖半）**，条件＝§4-③ 的 §5 五字段行由编排者具名追加（追加后 `DEFERRED(D34-1)` 与 §5 一一对应，1:1 尺闭）。**不同意票 24 整体结案**：PDF 两格＋search.files＋search.apps＋bridge 接线在盘上未落，须拆续票点名，⛔以「等 PDF」为名整体挂起零依赖已可做的续格。
