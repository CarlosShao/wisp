# 309-r2 · 50 门禁与"行为逐字节不变"（本腿的 AC#4 面）

时刻 2026-10-11 14:5x–15:0x +08。起手 HEAD＝`202ed4a5`（派单给的锚，与我现跑一致），交回 HEAD＝`d75db9ee`（本腿代码那一笔）；
其间 `dev` 上落过 308-r2 的笔（`6f389389`→`daa2f9be`），具名在下。

## 1. 四道门禁，每道自己落一行 rc（件⛔ 0 字节）

| 门禁 | rc | 件 |
|---|---|---|
| `GOFLAGS= go build ./...` | **0** | `logs/gate-build.log`（0 字节＝无诊断，rc 另落行） |
| `GOFLAGS= go vet ./internal/tools/` | **0** | `logs/gate-vet.log`（同上） |
| `sh scripts/d22scan.sh` | **0** | `logs/gate-d22scan.log`（**45 KB**，live 那句 `clean - no D22 ban violations` 命中 2 次） |
| `sh scripts/check-path-length-budget.sh --with-self-test` | **0** | `logs/gate-pathlen.log`／起手另有一发 `logs/pathlen-start.log`，VERDICT GREEN |

⚠ **d22scan 本腿两轮红＝全是我造成的，⛔ 别腿瞬态**：首轮 **rc=1** 两枚命中、次轮 **rc=1** 一枚命中，三枚都是我自己写的注释里那种
"把三个 root 拼成一枚路径"的引用——`bridge.go:892` 与 `:960` 的 `cmd/internal/tools/scripts`、`bridge.go:975` 的
`.scratch/wisp/probes/309/r2/30-comments.md`（写注释时那枚件还没落盘），ban #9 phantom-citation（票 212）。
处置＝root 那两处改写成"the four roots cmd, internal, tools and scripts"（⛔ 斜杠拼路径），证据件先落盘；
⛔ 动 `tools/d22scan/`、⛔ 动 allowlist。改后同一把尺 **rc=0**（三轮读数都留在 `logs/gate-d22scan.log` 里，红的那两段⛔ 删）。
另：`internal/risk` 那两枚文件起手就在飞（308-r2），本腿**没有**任何一发 `go build`／`go test` 红是指名 `internal/risk` 的，
⛔ 需要按派单那条"疑为别腿瞬态"重跑（读数＝`logs/gate-build.log`、`logs/gate-vet.log`、`logs/prov-rigs.stdout`）。

★**字符串面自证**（派单那条 `⛔` 落 ban #8 射程的规矩）：尺＝`grep -nE 't\.(Fatalf|Errorf)\(' internal/tools/toolrequest_nail_309_test.go`
逐枚读过，格式串里**零枚** `⛔ ≈ ✓ ≤ →`（要表达一律中文词，例："接近零"、"⛔ 报绿"只出现在注释里）；
本腿⛔ 撞过 ban #8（首轮红的是 ban #9），⛔ 动 allowlist。

## 2. 格式名册，两把尺并排，各写射程

| 尺 | 射程 | 枚数 | 本腿新增 |
|---|---|---|---|
| `gofmt -l internal/tools` | 工作树面 | **0** | 0 |
| `gofumpt -l internal/tools` | 工作树面 | **0** | 0 |
| `gofmt -l .` | 工作树面，整仓（含 `.scratch` 里别家故意脏的突变台件） | 43 | 0 |
| `gofumpt -l .` | 工作树面，整仓 | 50 | 0 |
| `gofmt -l`／`gofumpt -l` | **HEAD blob 面**：两枚被改文件从 `git show HEAD:` 导到 `/tmp` 里量（⛔ 落 `probes/**`，那会挪别家 gofumpt 的 tracked 分母） | **0／0** | 0 |

件＝`logs/fmt-final-tools.txt`、`logs/fmt-gofmt-*.txt`、`logs/fmt-gofumpt-*.txt`。
整仓那两档⛔ 归我：`cmd/wisp/models.go`、`internal/agent/tools.go`、`internal/agent/approval/pending_read.go`、`internal/risk/provenance.go`
与 `.scratch/wisp/probes/**` 那批台件，起手就在册（⛔ 我动，⛔ 我修）。
⚠ 一次自抓自纠：`behaviour-and-gates.txt` 起初把 HEAD blob 导进 `logs/headblob/*.go`，那＝在 `probes/**` 里新建 `.go`，
会挪别家 gofumpt/d22scan 的分母 ⇒ 改成 `mktemp -d` 仓外，且把先前落下的一枚 `logs/bridge-head.go` 改名 `.go.txt`；
现 `find .scratch/wisp/probes/309/r2 -name '*.go'`＝**0**。

## 3. "行为逐字节不变"这一把我跑了什么（`logs/behaviour-vs-anchor.stdout`）

判据打在**派单锚 ⇄ 本腿代码笔**之间（⛔ 打在"工作树 ⇄ HEAD"上：代码一笔落了之后那把尺退化成 0 枚变化，等于没量——这是本腿自抓的第二枚口径洞）。

- `git diff --name-only 202ed4a5..d75db9ee -- internal cmd tools scripts`＝4 行，其中两行 `internal/risk/*_test.go`
  ＝**308-r2 的笔**（`4c63d3b2`/`ba4c8d9e` 落在这段区间里）。本腿那一笔自己＝`git show --name-only d75db9ee`＝**2 枚**（下表）。
- **V2** `bridge.go` 逐行判据：`git diff -U0 202ed4a5..d75db9ee` 增 **44** 行／删 **9** 行，**44/9 全部以 `//` 开头**，
  非注释改动行＝**0**。
- **V3** 两枚 blob 各 `sed '/^[[:space:]]*\/\//d'` 剥整行注释后 `cmp`＝**IDENTICAL**，
  22233 字节／854 行 ⇄ 22233／854，sha256 同一枚 `496ed6d683e8e070c92410dde1c6a6d39e14a069d1773a9b35980a2d040e32c2`。
- **V4** 那枚守卫 `if b.injected || b.prov == nil || taskID == "" {` 在两枚 blob 里逐字同一，各 **1** 命中。
- **V5** 被触的第二枚文件是 `_test.go`（仪器面，无产品语义）；`internal/risk/**`⛔ 动一字、`lit236*`⛔ 动、
  `agent.ToolRequest`⛔ 加字段（`git show --name-only` 那两枚就是全部写面）。
- **B4/eol**：`tr -dc '\r' | wc -c`＝**0**（两枚文件），`git ls-files --eol`＝`i/lf w/lf attr/text eol=lf`（⛔ 重演 309-r1 那枚 CRLF 自抓）。

## 4. 整包名册：改前⇄改后各一发（逐名作差）

| 面 | 尺 | 顶层 PASS/FAIL/SKIP | 含子测试 | 名册行 |
|---|---|---|---|---|
| 改前（导出树 `202ed4a5`） | `^--- ` | 210／**1**／0 | 278／1／0 | 211 |
| 改后（导出树 `d75db9ee`） | `^--- ` | 210／**1**／0 | 278／1／0 | 211 |
| 工作树旁证（带别腿 `go.mod` 在飞态，⛔ 可比面） | `^--- ` | 211／0／0 | — | — |

- 两棵树各自那枚红＝同一枚 `TestTicket252R1AlignmentAddsNoAuthorization`（导出树缺件那一族，票 298 `AC#5` 已承接，⛔ 归 309）。
- 逐名作差：`comm` 两个 passname 文件 ⇒ **ADDED 0／REMOVED 0**（件＝`logs/suite-added.txt`、`logs/suite-removed.txt`）。
  本腿⛔ 新增测试名（六枚 control 与三形判定都装在同一枚 `Test309ToolRequestProductionConstructionNail` 里），
  所以"新增 0"是**设计**，⛔ 是漏跑；⛔ 一枚 `--- SKIP`。
- 两枚真喂空 TaskID 的 236 载体在两棵树上**都 PASS**（`logs/r359-12-13-carriers.log`，R359-12）。
- 导出树模板照派单：`git archive --format=tar <sha> | tar -x -C <mktemp -d>`，`go.mod`＋`go.sum` 各自验在场，
  量之前先数 `.go` 证非空（**995** 枚，`internal/tools` 下 **78** 枚，四 root 非测试源 **293** 枚）——三处读数在两棵树上都打了。
