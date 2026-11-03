# 24-v3 裁决表（编排者 Lead 出；实现腿＝24r2a/24r3/24r5/24r6 子代理，Lead 非实现者）

> 原 24-v3 子代理腿（qwen）跑 35+min 无输出，疑卡在 `-run '…Search…'` 走网络（本机 DNS 被劫 198.18/15 ⇒ hang），已 `job_kill`；改由编排者窄过滤自验并出本表。

## 证据
- **24-r5 `0950ddd7c`**：`internal/tools` PDF 分支改用 `github.com/ledongthuc/pdf`（BSD-3 纯 Go）真抽取——逐页 `Page(i).GetPlainText`、预算 spill 只报 `Truncated`（`appendTextWithin` 回退 rune 边界）、扫描件(runs==0)→「疑似扫描件，无文本层」marker 不猜、坏 PDF→具名 error fail-closed；**删 `DEFERRED(D34-1)` 与 stub**。库在 `page.go:250` 读 `/ToUnicode` CMap，真解 Unicode。
- **24-r6 `52ed59cff`／`f6fd8270a`**：把 2 枚旧 stub 测试（`Test24R2APDFStubIsErrUnsupported…`、`Test24R2ADocReadPDFIsARefusal…`）重写为真抽取判据（有效 PDF→抽出正文/端到端成功；坏 PDF→fail-closed 且 `errors.Is(err,ErrUnsupported)==false`），未删测试、未放宽其它断言。
- **编排者现跑** `go test ./internal/tools -run 'PDF|Pdf|Doc|Extract|Document' -count=1 -timeout 6m` ⇒ **`ok … 0.338s` rc=0 全绿**（避开 Search/网络）。
- **DEFERRED 登记**：`document_text.go:39/42` 有 `DEFERRED(xlsx, D34)`／`DEFERRED(OCR, D34)`；`SPEC-12 §5` 行 78「`doc.read` xlsx/OCR（D34）」在 ⇒ 1:1 对上。

## 逐格（票面 5 格）
| AC（票面行） | 判决 | 读数 |
|---|---|---|
| line34 `doc.read`：text PDF / 多页 >4k spill / docx / 扫描件 error | **PASS** | 24r5 新测 + 24r6 复跑 `ok 0.338s`；Scan→marker、Broken→fail-closed、Docx→PASS、budget spill 在 `appendTextWithin` |
| line35 `search.content`：fixture 查询→带 taint 结果、黑名单文件不出现 | **PASS** | 24-r2a（`8523211a0`）已验 |
| line37 `search.files`：glob + allowlist 作用域（越界→R2 拒） | **PASS** | 24-r3（`e1ee3b357`）已验，24-v2 同意 |
| line38 `search.apps`：枚举 fixture app；结果喂 `app.launch`（22） | **PASS** | 24-r3（`d59093153`/`dbd82a6b4`）已验，24-v2 同意 |
| line39 `DEFERRED(D-34)` xlsx/OCR 标记 + registry 交叉核 | **PASS** | 标记在 `document_text.go:39/42`；§5 row 78 在 |

## 结论
**票 24 五格全 PASS ⇒ 同意结案 `-done`。** PDF 文本层由「推迟」改为「已实现」（机主批 `A1445` 允许 ledongthuc/pdf），是本次结案的关键推进。
