# 24-r6 impl — 替换 24r5 上报的两枚旧 stub 测试（PDF 真抽取后的正确判据）

- 腿：票 24 修补腿 24-r6（实现腿）· 编排者派单 2026-11-03（存档 `.scratch/wisp/probes/orch/20261103/txt/20261103-24r6-dispatch.txt`）
- 前序：24-r5（`docs/evidence/s1/24r5-impl.md` §4.1 上报这两枚过期测试不在其写面）
- ⛔push：本腿只 commit，推送归编排者

## 1. 写了什么（显式 pathspec 逐枚）

| 文件 | 内容 |
|---|---|
| [internal/tools/document_text_test.go](../../internal/tools/document_text_test.go) | 仅替换派单点名的两枚 + 同步两处注释文案；其余测试与断言零改动 |
| docs/evidence/s1/24r6-impl.md | 本文件 |

未动：`document_text.go`（无需极小修，实现行为已满足新判据）、`document_text_pdf_24r5_test.go`、`search_content.go`、`bridge.go`、`cmd/wisp/**`、plugin/thresholds、SPEC/PLAN 契约文字、`go.mod`。

## 2. 两枚测试的替换（删 stub 断言 ⇒ 真判据，非放宽）

1. `Test24R2APDFStubIsErrUnsupportedAndCarriesTheMarker` ⇒ **`Test24R6PDFBranchIsARealReaderNotAStub`**
   - 有效 PDF（复用 24r5 的字节写入器 fixture `r5PDF`）：seam 抽出正文、Runs>0、⛔Truncated；
   - 打不开的 PDF：具名 fail-closed 错误（`无法打开文件`）、`errors.Is(err, ErrUnsupported)` 必须为 false、错误文本⛔含 `DEFERRED(D34-1)`、失败必带空 `Extracted`。
2. `Test24R2ADocReadPDFIsARefusalNotASuccess` ⇒ **`Test24R6DocReadPDFSucceedsAndBrokenPDFFailsClosed`**
   - 有效 PDF 经 `doc.read` 端到端：`IsError=false`、正文含绘制文本、`Origin`＝C26 规范路径（C25 打标缝）；
   - 坏 PDF（`%PDF-1.7 fake`）：`IsError=true` 且点名 `解析失败`＋`不是有效的 PDF`；⛔`DEFERRED(D34-1)`、⛔`该格式暂未实现`（坏文件不得洗成推迟）；错误回答同样带 `Origin`。
   - `DEFERRED(D34-1)` 断言全部删除；`DEFERRED(xlsx, D34)`/`DEFERRED(OCR, D34)` 语义未碰（`Test24R2AExtractorSetDispatch` 的 xlsx 拒绝、24r5 的扫描件/标记还原核原样在跑且绿）。

## 3. 判据尺（本腿现跑）

- `go test ./internal/tools -run 'PDF|Pdf|Doc|Extract|Document' -count=1 -timeout 6m -v` ＝ **全 PASS**（17 枚顶层＋4 子测试，含两枚新 24R6；`ok github.com/CarlosShao/wisp/internal/tools 0.273s`，TESTEXIT=0；⛔命中 Search/网络测试）。
- 整包此前唯红的正是这两枚（24r5 §4.1），替换后该包该域无 FAIL。
- gofumpt（绝对路径 `D:\work\base\gopath\bin\gofumpt.exe -w`）：`document_text_test.go` 已格式化，跑测即其后内容。

## 4. 备注（转呈编排者，本腿未动）

- 24r5 §4.2（go.mod `// indirect` 过期）与 §4.3（SPEC-12 §5 的 D34-1 尺向）仍开放，不属本腿写面。
- 票 24 AC 的「scanned PDF error case／spill」两格由 24r5 测试覆盖；真实中文 CJK/Identity-H 样本复跑（24r5 §4.4）仍待验收腿。
