# `304-r1` 30 — `AC#3` 门禁与越界（每把尺：命令逐字＋射程＋量的还是算的＋件里一行 `rc=N`）

每把门禁件的**最后一行就是它自己的 `rc=N`**（⛔ 0 字节件＝那格没交：`logs/g-build.txt`/`g-vet.txt` 只有 1 行＝那行 rc，命令本身零输出＝绿；`wc -l` 逐枚在表里）。
台件＝`gates.sh`（一次跑齐，逐字可重跑）。锚：跑的时候 `HEAD` 已从起手 `3974ab4d` 走到 `a1f84d39`（别的腿在共享 `dev` 上落了几笔），⛔ 我碰过别人的件。

| 尺 | 命令逐字 | 射程 | rc | 分母／读数 | 件 |
|---|---|---|---|---|---|
| 编译 | `GOFLAGS= go build ./...` | 整个主模块（⛔ 只我改的包） | **0** | stdout＋stderr 共 0 行（件 1 行＝rc） | `logs/g-build.txt`（`wc -l`=1） |
| vet | `GOFLAGS= go vet ./internal/llm/... ./internal/proc/` | 我动的两包＋复用的 `internal/proc` | **0** | 0 行输出；含 `_test.go` 编译面（两枚 spawn 点一枚是 `mockllm.go`＝⛔ `_test.go`，一把尺⛔ 看得见两枚，票面 `:15`＋`304-a1` 表①同判） | `logs/g-vet.txt`（=1） |
| D22 七禁扫描 | `sh scripts/d22scan.sh` | 全仓（件里自报分母） | **0** | 逐字末行：`d22scan: clean - no D22 ban violations; live scope work: bans #1-5 internal/=231, bans #1-5 cmd/=39, ban #6 frontend/=85, ban #7 internal/tools/=23`；ban#8 名册分母 `internal/` **533 Go files, comments and _test.go included**、`cmd/` 119、`design/` 39、`frontend/` 85 ⇒ 我新加的 2 枚 `.go` 落在 533/231 的分母里（⛔ 豁免、⛔ 动 allowlist） | `logs/g-d22scan.txt`（=252 行，末行 rc） |
| 路径长度预算 | `bash scripts/check-path-length-budget.sh --with-self-test` | 全仓 tracked | **0** | `denominator: tracked paths=10828 over-budget=57 covered by roster=57 not in roster=0`、`VERDICT GREEN`、`worst full path on the self-hosted runner=228 chars`（⇄ 编排者 `93de4b4a` 那笔快照 `tracked=10779`；⚠ 这个差**归因⛔ 成立**：其间别的腿有 commit，且我这轮新落的名册⛔ 进分母＝未 commit 的件⛔ 被数） | `logs/g-pathlen.txt`（=20 行） |
| 格式名册·工作树那把 | `$(go env GOPATH)/bin/gofumpt.exe -l <四枚路径>`（逐枚，⛔ 目录扫） | 只我碰的 4 枚：`internal/llm/adaptertest/{childjob_windows.go,childjob_other.go,mockllm.go}`、`internal/llm/openaichat/mockllm_integ_test.go` | **0** | **名册 0 枚**（件里 `gofumpt -l` 那几行**空的**＝0 枚，⛔ 把"没打印"当"尺没打中"）；`wc -l`=2＝rc 行＋scope 行 | `logs/g-format-worktree.txt` |
| 格式名册·HEAD blob 那把 | 同一把尺，件先 `git show HEAD:<path>` 取到仓外 `/d/tmp/wisp-304-r1/headblobs/` 再量 | 同上 4 枚 | **0** | **名册 0 枚**；`wc -l`=2。两把并排＝**同一集合、两个台面**（工作树⇄HEAD blob）⇒ 本轮⛔ 出现 CRLF 假枚（`gofumpt.exe` 本机⛔ 在 PATH 上，必须走 `$(go env GOPATH)/bin/gofumpt.exe` 全路径） | `logs/g-format-headblob.txt` |
| census（票面 `AC#3` 明写那把） | `bash scripts/portable-tests.sh --scope=census` | 主模块包名册 | **0**（两发） | 逐字 totals 行：跑前 `packages=35 with-zero-compiled-tests=7 claimed-by-no-scope=7 unclaimed-with-tests=0`＝跑后同一行 ⇒ `census_totals_delta=IDENTICAL`，且⇄票面 `AC#3` 写的现量基线**逐字相等** | `logs/g-census-before.txt`(39 行)／`logs/g-census-after.txt`(39 行)／`logs/g-summary.txt`(5 行) |

## 三处⛔ 算绿的说明（写给判语那程）
1. **census 的"跑之前"那一发是在我落地之后跑的**（唯一一次树呈"⛔ 带 job"原形的窗口＝17:50–18:05，全花在同一条窗口的 `AC#2` 反形上）。⇒ "逐字未变"这里成立的两把是：①⇄票面 `AC#3` 那条现量基线逐字相等，②⇄我这次门禁跑前/跑后两次作差为 0。**"我自己在改前跑的那一发"欠着**＝`Q-304-r1-3`。结构上为什么⛔ 会变（⛔ 新增 package、⛔ 动测试文件名册）是〔读码推〕级，⛔ 当读数。
2. `format_headblob_rc` / `format_worktree_rc` 这两行取的是 `for` 循环**最后一枚**的 rc（bash 语义），⛔ 是"四枚合起来的 rc"；合起来的证据是名册为空（`-l` 只要有一枚要重排就会打印名字）。
3. ⛔ 跑 `go test ./...` 整包、⛔ 跑 SLO／golden 判定：本轮我只跑了 `AC#2` 需要的两枚用例（`TestMockllmGoldenByteIdenticalEvents`、`TestMockllmSharedFaultSuite`，逐字在 `logs/paired-*.txt`），golden 文件与 `thresholds.go` 零字节改动（尺＝`git diff --name-only b9f37892 -- internal/llm/testdata internal/observe` 现量 **0 行**）。

## 越界尺（逐笔，⛔ 区间）
`git show --name-only --format= <每枚我自己的 commit>` 逐笔名册＝`logs/over-reach.txt`（本腿 commit：`b9f37892` 1 枚件／`e875f00f` 10 枚／后续两笔见该件），写面只落 ①`.scratch/wisp/probes/304/r1/**` ②`internal/llm/**` 那 4 枚 ③票 304 末尾 1 行 bullet；`foreign=` 逐笔名册里⛔ 出现别家路径（判语⛔ 归我）。
