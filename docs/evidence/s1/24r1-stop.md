# 24r1-stop — 票 24（doc.read + 本地检索）落地腿 24-r1 停手上报

- agent=24-r1（实现腿／写腿）· 日期=2026-11-02 22:1x +08 · ⛔push（本腿只 commit）
- 结论：**未定义即停（AGENTS §0／派单闸③）——PDF 解析依赖未列。** 零代码落地、零文件改动，本件为唯一交付。

## 1. 起手三核（盘上读数，全部现跑）

| 尺 | 命令/位置 | 读数 |
|---|---|---|
| HEAD | `git log --oneline -1` | `2bf1dcde2` 编排者A1434：派24-r1(文档检索工具,internal/tools);50暂park… |
| 索引/工作树 | `git status --porcelain -- internal/tools` | **空**（web*.go 无他脏；全树仅 `??` 存档杂物，无他腿脏写面） |
| 票面未落地 | `ls internal/tools` + 票文件名 | 无 `doc*/pdf*/search*` 新件；票为 `24-doc-search-tools.md`（**无 `-done` 后缀**、Status: ready-for-agent）⇒ 不属「已落地/已-done」停手支，本腿因**依赖闸**停 |

## 2. 停手闸因＝新第三方依赖未列（派单★条款原文命中）

> 派单：「★若读取 PDF/docx 需新第三方依赖而 go.mod 未列 ⇒ 未定义即停（⛔自选库/发明路径）」

三处盘上证据：

1. `go.mod` `require` 全列（:6-17）＝ go-webview2 / sherpa-onnx / go-toml v2 / x/crypto·sys·term / modernc sqlite（+indirect）。**无任何 PDF/docx 解析库**。
2. `docs/specs/SPEC-01-architecture.md:96-101` 依赖白名单（D22）**逐字标「新增须人批准」**，在册者＝sherpa-onnx·goja·go-webview2·x/sys·SQLite·TOML·goreleaser·webrtc-audio-processing —— **文档解析器不在册**。
3. 票面自证 `24-doc-search-tools.md:16-17`：「PDF text extraction + docx XML text extraction with pure-Go **or whitelisted-dep parsers** (propose + justify any new dependency — **whitelist change needs approval**)」。

**技术判读（agent 自拍、不丢回）**：docx＝OOXML（zip+xml），标准库 `archive/zip`+`encoding/xml` 零新依赖可做；search 家族（files/content/apps）可做（stdlib+既有 `risk.PathResolver`/C26）。**唯 PDF 文本层抽取绕不开第三方**：内容流解压只是第一步，ToUnicode CMap/字体编码/多对象交叉表纯 stdlib 手写＝自造解析器（违「⛔发明路径」，且主流 harness——含 SPEC 参照系——一律用命名库如 pdfcpu/unipdf）。自选库＝白名单变更＝人工批准，属摆机主三类之一 ⇒ **本腿无权选，停**。

## 3. 引文出处核对（第19/20处）——票面引用全部在场

| 票面引用 | 盘上现读 | 判定 |
|---|---|---|
| SPEC-07 §3 | `docs/specs/SPEC-07-tools-and-plugins.md:32` 「## 3. D34 内置工具权威表」 | ✅ 在（票写路径 SPEC-07 实际文件名 `SPEC-07-tools-and-plugins.md`，同物） |
| D34 表行 | `docs/PLAN.md:2549/2550/2551/2555` search.files/content/apps + doc.read（L0/fs.read/S3/taint/DEFERRED） | ✅ 逐字在 |
| SPEC-07 §3 表行 | `SPEC-07…md:49-51/55` 同形四行 | ✅ 在 |
| DEFERRED 登记 | `SPEC-12-roadmap-governance.md:78`「DEFERRED｜doc.read xlsx/OCR（D34）」；`PLAN.md:1532/:2626` | ✅ 在册；⚠ **注意**：票面要求代码标记 `DEFERRED(D-34)`，而 §5 登记表该行**没有独立 `D-xx` 编号字段**、仓内既有标记形为 `D28-1`/`D11-3`（`grep "DEFERRED("` 实测）⇒ AGENTS §1.1「1:1 双向」落具体标记形需编排者定，落地腿⛔自造 |
| 前置票 20 | `.scratch/wisp/issues/20-host-bridge-fs-tools-done.md` | ✅ 已 done，前置成立 |
| 派单存档 | `.scratch/wisp/probes/orch/20261102/txt/24r1-dispatch.txt` | ⚠ **不在盘**（`Test-Path` = False）——如实记，派单正文以本会话收到的为准 |

## 4. `search.content` SPEC 面（次生观察，非本闸主因）

SPEC 只给权威表行（名字/L0/fs.read/taint），工具参数面与扫描语义只写在票面 constraints——比 web 族（票 22 有 SPEC 路径裁定）薄。因 PDF 闸已先行停手，此条**不单独成立停手**，登记交编排者裁：落地前建议补一段参数/语义 SPEC（⛔由实现腿发明）。

## 5. 处置建议（供编排者）

1. **摆机主批白名单**：PDF 文本抽取选型（主流纯 Go 候选：`pdfcpu`＝Apache-2.0/MIT 双许可件；`unidoc/unipdf`＝AGPL⛔不适用自用分发——编排者复核后按铁律「最有价值+主流」提名一枚摆批，⛔实现腿自选）。
2. 批准 → 重派 `24-r2` 逐格做 AC#1-#5；或裁「拆格」：先落地零新依赖半（docx + search 三件），PDF 格待批——**拆不拆由编排者拍**，本腿不自行做半票（票面 AC#1 含 PDF 格，硬做半票即违「⛔硬做」）。
3. `DEFERRED(D-34)` 标记形与 §5 表 1:1 的落法先行定字（见 §3 ⚠ 行）。

## 6. 本腿读数汇总

改动文件＝**仅本件**；`go test`＝未跑（零代码）；commit＝本件一枚（显式 pathspec、⛔push）；写面 `internal/tools/**` 一字未触。
