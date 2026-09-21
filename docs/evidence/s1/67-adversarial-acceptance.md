# 票 67 对抗验收 —— 让扫描器可信：`mockllm` 走受管出口 + 门必须自报覆盖面

**验收人**：编排者（**非实现者**；实现方是 `agent-ticket67` 与续做的 `67b`）
**验收时间**：2026-09-21 09:48
**被验收的 commit**：`fff4cad`（AC#3 第一步字形 ASCII 化）、`bcf44d6`（AC#3 覆盖面扩展）、
`a8ae9ad`（编排者收口：新门上线后抓到的第一个真命中）
**证据档位**：〔独立复现〕／〔日志＋归档，我抽验〕／〔仅自述，不背书〕

## 裁决表（与票面 4 个 AC 框 1:1）

| AC | 裁决 | 我跑的那条判据 | 档位 |
|---|---|---|---|
| **AC#1** `mockllm.go:68` 裸 `go func()` 改走 `observe.Registry.Spawn` | **PASS** | 静态面：`cd tools/d22scan && go run . -root ../..` 对 `internal/`+`cmd/` 全量扫 ban #1 **exit 0** ⇒ 那行不再是裸 spawn；机制面：`go test -count=2 ./internal/llm/adaptertest/` → `ok 3.401s`。**我明确不信"改完就好了"这类自述**：本票立场是"**不许用 allowlist 豁免把它变绿**"，所以我核的是豁免文件本身（见下 AC#4 行）。 | 〔独立复现〕 |
| **AC#2** 门自己必须可证伪（含"把仓根误调用变成不可能"） | **PASS** | `cd tools/d22scan && go test ./... -v` → `=== RUN`/`--- PASS` 合计 **24** 条、`SKIP\|FAIL` **0** 条、`ok 0.586s`。这条就是我 A16 那次事故的解药：**扫描器被 seeded-violation 自测钉住，"没有命中输出"不再等于"跑过了"**。误调用面：`main.go:509` 的守卫把"作用域内 0 文件"从**静默空转**改成**致命退出**（`67b` 报告：构造 `cmd/` 无 `.go` 的假仓 ⇒ `examined 0 file(s)` + **EXIT=2**），比原裁定（只要求 `N==0` 报错）更严。 | 〔独立复现〕（EXIT=2 那条为〔日志＋归档，我抽验〕：我核对了报文格式与 `emojiScopes`/`emptyEmojiScope` 代码位置，未重跑它的假仓构造） |
| **AC#3** emoji 门覆盖面扩到 `internal/`+`cmd/` | **PASS** | 两条我自己敲的：①`go test ./... -run SelfScan -v` → `--- PASS: TestScannerSelfScanOfRealRepoIsGreen (0.24s)`；②纯净树 `git archive HEAD` 解到仓库外再扫 → **exit 0**，末行自报 `internal/ 289 Go files, comments and _test.go included; cmd/ 21`。⚠ **这条框一直红到 09:40**，元凶是比 `bcf44d6` 晚 5 分钟入库的 `d63bc49` 注释里一个 `U+26A0` ⇒ **HEAD 的 lint 红约 13 分钟**，由我用 1 行纯注释 ASCII 收掉（`a8ae9ad`）。判据已写进**票 71 第 3 条**：**扩覆盖面必须同批修完它新照到的存量违规**。 | 〔独立复现〕 |
| **AC#4** 门禁（票面自订："**AC#3 落地后必须重跑**"） | **PASS（09:48 重跑，这才是这一框现在的凭据）** | `gofmt -l internal/llm/adaptertest internal/llm cmd/wisp` → **空**；`go vet ./internal/llm/... ./cmd/wisp/` → **rc=0**；`go test -count=2 ./internal/llm/adaptertest/` → **ok 3.401s**；`cd tools/d22scan && go test ./...` → **ok 0.586s**；扫描器纯净树 **exit 0**；**`allowlist.txt` 非注释行 5 行，且该文件最后一次改动是 `38b3715`（票 70 的 R16#1），不是本票** ⇒ **票 67 全程没有为凑绿加过豁免**。 | 〔独立复现〕 |

## 两处必须写下来的诚实限定（不许被"PASS"吃掉）

1. **主模块整仓 `go vet ./...` 我这次**故意没跑**。票面原判据要求它，但此刻 `internal/ball/` 正被票 74 的代理改着
   （`tokens.go`/`renderer_windows.go`/`tokens_table_test.go` 在工作树里未提交），跑全仓等于**把邻居的 WIP 当成 HEAD 来判**。
   我的替代：①按包 scope 复跑；②用 `git archive HEAD` 纯净树覆盖"扫描器"那一格。
   **残余风险**：`internal/ball` 自己是否有 vet 回归，本表**不背书**；那一格由票 74 落地后我复跑，并由 CI 的 lint job 兜底。
   ⚠ 这条限定同时是 A39 那条判据的又一例：**共享树里"我看到的绿"必须说明它是 HEAD 的绿还是工作树的绿**。
2. **`allowlist.txt` 从 4 行变 5 行**，AC#4 的旧文本写的是"仍是 4 行"。差的那一行是**票 70** 在 `38b3715`
   按 R16#1 加的**按文件**豁免，不是票 67 加的。旧文本我**不覆盖**，在这里追加更正（第 13 条规矩）。

## 结论

**PASS，4/4** ⇒ 归档 `-done`。本票真正的交付不在这四个框，而在**它之后所有人调用门禁的方式变了**：
`frontend/` 那个恒 0 文件的死作用域被删掉、"clean" 那行改成由代码生成的自报覆盖面、
以及"0 文件"从静默空转变成致命退出。这三件事都是**判据仪器的加固**，价值高于任何一条功能实现。

## 移交（带票号，不许沉底）

- ban #6 仍指向 `frontend/`（禁令文本，属 D22）→ **票 71**
- `ci.yml:23` 注释与现覆盖面不符 → **票 71**（票 70 已收尾，该文件现归票 71 的代理）
- "扩覆盖面须同批修存量违规" + 纯净树判据 → **票 71**
