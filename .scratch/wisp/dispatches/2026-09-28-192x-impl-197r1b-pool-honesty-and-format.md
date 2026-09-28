# 派单 2026-09-28 19:2x — 写腿 `197-r1b`：把票 197 实体层留下的三处小账收平（池诚实＋格式＋一枚假注释）

## 0. 起手必做（不许跳过，跳了就算没交件）

1. **自取锚点**：`date "+%Y-%m-%d %H:%M %z"`＋`git log --oneline -8`＋`git status --short`。**别引用本单里任何行号当现量**，本单的行号是我 19:2x 读的，会漂。
2. **撞钉预检（这次要求比平时严）**：动手前先把这两枚包跑一遍并**抄下今天的绿用例名**——
   `PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test -count=1 ./internal/tools/ ./cmd/wisp/`
   ⚠ `cmd/wisp` 不加这枚 PATH 会报 `exit status 0xc0000135`（DLL 没找到，**一枚用例都没跑**，看着像包坏了其实是仪器瞎）。
   把这批绿名抄进你的证据件；你改完之后**同一批名必须还在**（少了哪一枚要具名解释）。
3. 读本票：`.scratch/wisp/issues/197-*.md`、`.scratch/wisp/issues/211-*.md`，以及上一轮证据件 `docs/evidence/s1/197-subagent-entity-r1.md`。

## 1. 硬规矩（违背一条即退回）

- **只 commit，绝不 push**；commit 必带**显式 pathspec**；禁 `git add -A` / `git add .` / `-a`；禁 `--amend` / `reset` / `rebase` / `stash` / `checkout .` / `restore` / `clean` / `merge` / `worktree`。
- 新文件先 `git add -- <那个路径>`（**逐枚点名**，不许目录）。
- **临时件只建不删**；绝不删仓库内任何东西；要还原被自己写坏的文件只用 `git cat-file blob HEAD:<path> > <path>`，**别用 checkout**。
- **写面**只有下面第 2 节列出的那几枚文件。别人在飞的dirty 文件（`.gitignore`、`.scratch/wisp/probes/152/**`、`.scratch/wisp/probes/161/r6/logs/**`、`design/**` 那 16 枚删除、`internal/panel/instructions_200.go` 那枚未跟踪文件）**一律不碰、不提交、不还原**。
- `frontend/**`、`design/**` **不写不读不转述**；`internal/panel/tokens_fourway_test.go` 一字不动；`internal/tools/bridge.go` 与 `internal/agent/budgets.go` **只能读、不能写**（本腿不动契约）。
- 写中文长段一律走 Edit/Write 工具，**别用 shell heredoc**（反引号会被真执行）。

## 2. 这一腿做什么（三件，按顺序）

### (a) 池的诚实数：`MaxConcurrentSubagents` 与桥的天花板对齐（票 211 的**甲**）

- 现量（我 19:2x）：`internal/tools/subagent_197.go` 里 `MaxConcurrentSubagents = 8`，占位口 `TryAcquireSubagentSlot(MaxConcurrentSubagents)`；
  `internal/tools/bridge.go` 里 `const MaxToolConcurrency = 4`，注释写死"不可上调"（D38d，契约）。
  上一轮实测：8 枚并发桥调用被排成队并撞 C22 的 30s ⇒ **7/8 红**。
- **裁定（不要再问，也不要往上改）**：把池改成**等于桥的天花板**（今天＝4），语义形状不动（超了一枚＝硬拒＋可读理由，不排队、不留名册行）。
- 新增**一枚常驻判据**：池常量 **大于** `MaxToolConcurrency` 时必须红。
  ⚠ **正控必做**：人为把池写成 8（临时改，别提交）→ 那枚判据**必须真的红**，把改前/改后读数都抄进证据件，然后改回 4。
  只写"两枚相等"这种今天恒真的断言**不算判据**。
- 池的测量**改走真桥**（父任务真次派生第 5 枚会被拒），把 `spawnDirect` 那枚绕行要么去掉、要么在它旁边**写清它量的是池语义而不是生产并发**（二选一，选哪个在证据件里说理由）。
- ⚠ 会打红的现成钉：`internal/tools/subagent_197_test.go` 末尾那枚"常量漂了"断言写死 `want 8 / 1`；`Test197SubagentPoolCapsAtEight` 的循环/期望也按 8 写。**改断言要连着改用例名里那句"Eight"**（名字和行为不一致比红更难查）。深度 1 不动。

### (b) 格式门：把本仓跟踪集里那 3 枚源件跑成 gofumpt 干净

- 现量（我 19:2x，用 `$GOPATH/bin/gofumpt.exe`，版本 `v0.12.0 (go1.27.1)`）：跟踪的 `.go` 里 `gofumpt -l` 点出 **7 枚**——
  `cmd/wisp/run.go`、`internal/tools/subagent_197.go`、`internal/tools/subagent_197_test.go` ＋ 4 枚 `.scratch/wisp/probes/**` 的变异样本（后 4 枚**不归你**，别动）。
- `gofumpt` **不在 PATH 但在 `$GOPATH/bin`**：Windows 下 `"$GOPATH/bin/gofumpt.exe"`。上一轮报"本机没有 gofumpt"是不实（`gofmt -l` 空 ≠ gofumpt 空，两者判的不是一回事）。
- 我核过：`cmd/wisp/run.go` 在 `b9fa815b` 之前是**干净**的，这枚回归是 197 那一腿带进来的。
- ⚠ **写面冲突保护**：动 `cmd/wisp/run.go` 之前**当场**跑 `git status --short -- cmd/wisp/run.go`；
  如果那一刻它是脏的（说明有别的腿正在改它），**跳过这枚文件**、只处理另外两枚，并在证据件里具名报"run.go 那一刻被别人占着，我没动"。别硬合并、别替别人格式化。
- 交件判据：`git ls-files -z '*.go' | xargs -0 "$GOPATH/bin/gofumpt.exe" -l` 的输出里**不再出现** `internal/`、`cmd/` 下的任何文件（`.scratch/wisp/probes/**` 那 4 枚允许还在，逐名列出）。

### (c) 一枚**假注释**：`SubagentStreamKeyPrefix` 的那枚"钉住两枚拷贝相等"的测试**盘上不存在**

- 现量（我 19:2x）：`internal/tools/subagent_197.go` 顶部注释声称 `cmd/wisp/subagent_stream_key_197_test.go` 把 `internal/tools` 与 `internal/panel` 里同名的两枚字面量钉成相等；
  但 `grep -rln "SubagentStreamKeyPrefix" --include=*.go .` 只点出 4 枚文件（两包各一枚源件＋各一枚测试件），**`cmd/wisp/` 下一枚都没有**，那枚文件不存在。
- **修法二选一，选一支并写理由**：
  - **丙（推荐）**：把注释说的事**真做出来**——新建 `cmd/wisp/subagent_stream_key_197_test.go`，比较两包的那枚前缀常量相等，并配**正控**（临时把 `internal/tools` 那枚改成别的字面量 ⇒ 判据必须红；改回）。
    ⚠ 别在这条里去"合并成单一真源"——那是票 197 的**载体层 r3** 的活（我已经在账 `A406` 裁过：panel 侧是真源，`cmd/wisp` 注入）。这一枚先把注释不再撒谎。
  - **丁**：不写测试，只把注释那句"谁钉住了"删掉/改成实话。选丁要在证据件里说明为什么不补钉。

## 3. 交件（三样，缺一不可）

1. 代码＋测试，**一次或两次 commit**，message 里带上"197-r1b"和票号 197/211。
2. 证据件 `docs/evidence/s1/197-subagent-entity-r1b.md`，逐节写：
   ① 起手锚点（sha＋时间）；② (a)/(b)/(c) 各自的**改前读数→改后读数**；
   ③ **两发变异**：至少 (a) 的正控（池写成 8 ⇒ 那枚新钉红）＋ (c) 的正控（改一枚字面量 ⇒ 钉红），
   写清"哪一枚 commit 才算未修码"，别只说"会红"；④ 门禁四数＋全仓 `gofumpt` 名单（含 `.scratch` 那 4 枚的逐名归属）；
   ⑤ **没测什么**（具名，不许留空）；⑥ 对编排者（我）派单里任何一处**不服**——我这次给你的行号、常量、结论有错就当面写出来。
3. 报告回来时**只报盘上可核的事实**：commit 号、文件字节数（`wc -c`）、判据逐名的读数。我不认"应该已经"。

## 4. 明确不做（避免你顺手扩大射程）

- 不动 `MaxToolConcurrency`、`DefaultToolTimeout`、C22 预算、`thresholds.go`、golden、`allowlist.txt`、任何 SLO 阈值。
- 不新建第二座桥、不给子代理单开桥实例（那是票 211 的**乙**，契约级，要 owner 一句话）。
- 不碰 `internal/panel/**`（那一包的写面此刻可能被票 200 占着）。
- 不做票 197 载体层（r3）的任何一格。
