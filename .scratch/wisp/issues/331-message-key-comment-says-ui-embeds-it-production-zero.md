# 331 — `internal/observe/errors.go` 那句注释说"UI 和日志已经把这个键嵌进去了"，而**生产嵌键 0 枚**（过期文字族；本票只改注释、⛔ 动产码一字）

- Status: **待派**（料齐，凭据在下面两枚尺，本票起手仍要落地腿自己复跑一遍）。
- 来路＝非实现者裁决腿 `216-v2` §5 第 5 条（`docs/evidence/s1/216-v2-verdict.md:65`）＋编排者 10-17 自己复跑的那把尺；账 `A1010`。
- ⚠ **本票不是"错误文案体系"那一层**：那层已登记推迟＝`docs/specs/SPEC-12-roadmap-governance.md:77` 逐字 `| DEFERRED | 完整错误文案体系 | D23 | 所有用户可见错误有人话文案+建议（含 D37 全 17 类） | S6 | 部分错误暴露技术细节 |`。**推迟的是"做那层文案"，不是"注释可以替它说假话"**。本票只把那句不实注释改成说实话；⛔ 在本票里实现文案层、⛔ 新加导出名、⛔ 新增 C17 方法名。

## 现量（起手逐枚复跑；锚＝`0e944d93`）

| 事实 | 读数 | 尺 |
|---|---|---|
| 注释自称 | `internal/observe/errors.go:160-162` 三行逐字："MessageKey returns the stable user-visible copy key for the class. The copy text itself is deferred (D23 error-copy work, S6); **UI and logs embed the key** so copy can be added later without re-touching error sites." | 现读 |
| 生产嵌键 | **0 枚**——`MessageKey` 在射程 `internal cmd tools`（⛔`*_test.go`）只剩 `:160` 自家注释＋`:163` 声明两行 | `git grep -n "MessageKey" HEAD -- internal cmd tools ':(exclude)*_test.go'` |
| 消费者在哪 | 只在测试里：`internal/observe/errors_test.go:191-192`、`internal/llm/matrix_14_2_test.go:175`／`:249` | 同尺去掉 `:(exclude)` |
| 正控（尺不是恒 0） | 同尺含测试＝**4 处在场**（上面那四行）；同族尺在 `internal/risk` 出 **37** 枚 | 见 `216-v2` 台件 `M1`／`M2` |
| 那句"将来加上"的前提成立吗 | **不成立**：`MessageKey()` 返回 `"error." + string(c)`，全仓没有任何 UI/日志调用它 ⇒ 注释里"so copy can be added later without re-touching error sites"许诺的是一条**今天没有消费者、因此也无所谓"不重碰"**的路 | 现读码 |

## 判据

- [ ] **AC#0（料格）**：起手把上面四把尺**逐枚重跑**（一律取具名提交树，⛔ 读工作树在飞字节），枚枚带 `rc=`；与表里读数不同 ⇒ 当场具名更正并把差异写进证据件，⛔ 照抄票面。
- [ ] **AC#1（改口）**：把 `internal/observe/errors.go:160-162` 那三行注释改成**说实话**：① 键是为那层文案预备的、② **今天生产没有任何地方嵌它**（消费者只有测试，逐枚具名）、③ 那层文案仍归 `D23`／`S6` 的 `DEFERRED`（`SPEC-12 §5` 第 77 行）。⛔ 写成"将来也永远拦得住"式许诺；⛔ 顺手加"我们已经有 UI 用键"的任何同族句。
- [ ] **AC#2（产码字节零改动那把尺）**：滤掉**整行注释**后，锚 blob ⇄ 工作树两侧**逐字节同串**（`sha256` 两侧同值＋`diff` 空、`rc=0`）；★**正控必做**＝只在副本的**非注释行**改一枚字面量 ⇒ 同一把尺必须回**非空、rc=1**（否则这把尺两种形状都放行＝恒真，不许拿它的"空"当证据）。逐行核＝改动枚数、其中非整行注释枚数**必为 0**。
- [ ] **AC#3（行为不变）**：`go test ./internal/observe/ -count=1 -timeout 30m -v` 改前／改后各一发，台面＝`git archive` 的**仓外导出树**（⛔ 共享工作树跑"改前那发"），四数逐枚点名（`^=== RUN` 含子测试 ⇄ `^--- PASS` 只数顶层、⛔ 混用）＋ `^panic` 枚数；终态红名册与起手**作差＝0**。
- [ ] **AC#4（门禁四数，只在共享工作树跑）**：`gofumpt.exe -l`（真身 `/d/work/base/gopath/bin/gofumpt.exe`，PATH 里可能 rc=127）／`sh scripts/d22scan.sh`（★它自己声明 `design/`＋`frontend/` 射程，指到缺 `design/` 的树会 `exit 2`＝GUARD 拒绝，⛔ 读成红）／`sh scripts/check-path-length-budget.sh`（要 `.git`）／`go vet ./internal/observe/`，逐枚落 `rc=`；红句含 `walk`／`build failed`／`no test files`／`setup failed`／`0xc0000135` ⇒ 判**台面故障**、具名写是哪一行，⛔ 计入判据。
- [ ] **AC#5（纪律格）**：⛔ 翻任何框（本票六枚 `- [ ]` 由**非实现者**裁）、⛔ 改 `Status:`／票名／加 `-done`；写面只五处＝`internal/observe/errors.go` 的**整行注释**面＋`.scratch/wisp/probes/331/r1/**`＋`docs/evidence/s1/**` 一枚件（短名、⛔ 带 `§`）＋本票尾**一条** bullet；⛔ `git add -A`／`.`、`git add` 与 `git commit` 两处逐枚显式 pathspec；只 commit、⛔ push；⛔ `--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`；临时件只建不删。

## 禁区

⛔ 动 `MessageKey()` 函数体或签名一字（本票零产码改动，AC#2 就是钉这件事）；⛔ 动 `internal/observe/thresholds.go`、golden、`tools/d22scan/allowlist.txt`、`go.mod`／`go.sum`／`.gitignore`／`ci.yml`／`.github/**`、三枚冻结件、`frontend/**`、`design/**`、`PLAN.md`、`docs/specs/**`；⛔ 为"让注释说满"去补一条生产调用（那是 `D23`／`S6` 那层 DEFERRED 的活，不归本票）；⛔ 跑整包 `./cmd/wisp`、⛔ `go test ./...`。
