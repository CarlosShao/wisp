# 136 — AC#14 `Gate` 落地件（`136-ac14b-impl.md`）的非实现者终裁 r2 —— **acceptance**

> **这份文件是什么**：票 136 `AC#14` 的**第二张**裁决表。前一张 `136-ac14-r1-acceptance.md`（429 行／7 枚
> commit）裁的是**前一程** `worker-ticket136-ac14`（`aef82f5`）并给出附条件＝编排者批准翻布尔；本表裁的是
> **那个批准落地之后**的盘上终态（落地程 `worker-ticket136-ac14b`，证据 `136-ac14b-impl.md`，码 `52191ce`）。
> **本程身份**：`acceptor-ticket136-ac14b-r2`，**只读裁决者**，不是落地程、不是前一程实现者、不是编排者。
> 本程**不翻勾、不改码、不 revert**；发现的一切缺陷一律报回，不就地修。
> 本文件只引用变量名与文件名，不含任何凭据值（本程没有读到过任何凭据）。

---

## §0 锚点与口径

| 项 | 读数（本程自己量） |
| --- | --- |
| `git rev-parse HEAD` | **`4ecc284d1dd7b2b39f35bf54ae2222ec8d837548`**（简报说的 `4ecc284` 或其后代 ⇒ 实测正是 `4ecc284`，未漂） |
| 被验码 | `52191ce feat(136 AC#14 Gate): flip settleCoverageRowGates true + rewrite 2 coverage legs …`；`52191ce^` = `ade897c` |
| 取数范围 | `git diff 5a946d3..HEAD`（`5a946d3` ＝ 落地程 §0 自量锚点，也是编排者落批准那枚 commit） |
| 工具链 | `go version go1.27.1 windows/amd64`；`CGO_ENABLED=1`；本程自建二进制见 §5 |
| 工作树里别人的东西 | `git status --porcelain` ＝ 16 枚 `design/**` 删除 ＋ 2 枚未跟踪目录（`design/old/`、`design/doubao/`）⇒ owner 的，**本程一枚没碰、没还原、没代提交、没 stage** |
| 临时件落点 | 全部在 `D:\tmp\wisp141r2-ac14b-002e0e3e\`（`head/` `snapA/` `snapC/` `bin/`）与 `D:\tmp\ac14b-r2-136\`；**仓库目录内没有新建任何东西**；只建不删 |
| 快照怎么来的 | `git archive HEAD \| tar -x`（仓外）⇒ 三形快照：`head`＝终态、`snapA`＝只还原测试改写（gate 仍 true）、`snapC`＝`head` 再摘掉门行构造点（M1）。**全部在仓外，未在仓库内建 worktree/checkout** |

---

## §1 判据①：改动范围 —— 〔独立复现〕**成立**

### 1.1 `internal/observe/` 内的 diff 逐字节

```
$ git diff --name-only 5a946d3..HEAD -- internal/observe/
internal/observe/sampler.go
internal/observe/sampler_settle_coverage_136_test.go
$ git diff --stat 5a946d3..HEAD -- internal/observe/
 internal/observe/sampler.go                        |  2 +-
 .../observe/sampler_settle_coverage_136_test.go    | 56 ++++++++++++++++++++--
 2 files changed, 53 insertions(+), 5 deletions(-)
```

`sampler.go` 全量 diff ＝ **恰好一枚字面量、1 增 1 删**，站点 `:556`：

```
-const settleCoverageRowGates = false
+const settleCoverageRowGates = true
```

⚠ 简报里那句"`rep.Pass = foldSettlePass(...)` 那一行**若已在**就不该再动"—— **盘上确认它本来就在**：
`sampler.go:537`，`git blame -L 537,537` ⇒ `aef82f5d`（前一程 09-24 19:08）、`git log -S "foldSettlePass(memOK"`
⇒ 唯一引入枚 `aef82f5`；本范围 `5a946d3..HEAD` 对 `sampler.go` 的 diff 只有上面那一行（1 增 1 删）⇒ 落地程**没动它**，
与授权面（"只 `:556` 一枚布尔"）逐格一致。同理 `:529` `rep.Verdicts = buildSettleVerdicts(*rep)` 也 blame 到 `aef82f5d`、
本范围未触碰——这一行是 §4(C) 那一发变异的靶子。

### 1.2 冻结面 0 行

```
$ git diff 5a946d3..HEAD -- internal/observe/thresholds.go | wc -l
0
$ git diff --name-only 5a946d3..HEAD | grep -icE 'golden|testdata|thresholds'
0
```

全仓 golden/testdata 清单在本范围**逐枚 0 行**（`internal/llm/testdata/golden/*.sse`、
`internal/agent/testdata/golden/*.sse`、`internal/llm/golden/**` 一字节未动）。

### 1.3 没有别的包动过

```
$ git diff --name-only 5a946d3..HEAD | grep '\.go$' | xargs -n1 dirname | sort -u
internal/observe            ← 全范围唯一动过的 Go 包
$ git diff --name-only 5a946d3..HEAD
.scratch/wisp/issues/136-…-54-tests-stay-green-without-it.md
docs/evidence/s1/136-ac14b-impl.md
docs/evidence/s1/141-ac2-stock-inventory-r1.md
internal/observe/sampler.go
internal/observe/sampler_settle_coverage_136_test.go
```

### 1.4 逐枚 commit 的 `git show --name-only`（本范围 17 枚，每枚**只带一枚路径**）

| sha | 带的路径 | 在授权面内？ |
| --- | --- | --- |
| `52191ce` | `internal/observe/sampler.go` ＋ `sampler_settle_coverage_136_test.go` | **是**（两枚，唯一码 commit） |
| `b9ca0b0` `1e620d6` `1657135` `115173b` `92dd40f` | `docs/evidence/s1/136-ac14b-impl.md`（各一枚） | **是** |
| `a4b5deb` | `​.scratch/wisp/issues/136-….md`（票面 delivery note） | **是** |
| `ade897c` `5dfb8ba` `e0c4ada` `64664d6` `d036929` `a2d84ac` `94f4e33` `68ff486` `877f979` `4ecc284` | `docs/evidence/s1/141-ac2-stock-inventory-r1.md`（各一枚） | **否 ⇒ 见下** |

**越界旗（按简报字面清单如实报）**：本范围 17 枚里有 **10 枚**（`ade897c`、`5dfb8ba`、`e0c4ada`、`64664d6`、
`d036929`、`a2d84ac`、`94f4e33`、`68ff486`、`877f979`、`4ecc284`）落在简报给的三面授权清单**之外**——它们全是
`docs/evidence/s1/141-ac2-stock-inventory-r1.md`（票 141 盘点程 r1 的第 1—10 节）。
**判**：**不是本程要裁的越界**。简报正文自己就预告了"`docs/evidence/s1/141-*.md` 与 `design/**` 的 churn 是*别的*活"，
且这 10 枚**没有一枚碰过 `internal/`、`cmd/`、阈值、golden、票 136 的任何一面**；票 136 的写者只有
`52191ce`/`b9ca0b0`/`1e620d6`/`1657135`/`115173b`/`92dd40f`/`a4b5deb` 七枚。**结论：落地程写集与批准面逐格吻合；范围判据成立。**

**档位：〔独立复现〕**（`git diff`/`git show --name-only` 本程自己跑，不抄 §1.1—§1.4 任何转述）。

---

## §2 判据②：独立复现全绿 —— 〔独立复现〕**成立，四数与名册逐格对上**

本程在**跟踪树 `4ecc284`** 上自己跑 `go test -count=2 -v ./internal/observe/`（原始件
`D:\tmp\ac14b-r2-136\test-head.txt`）：

| 口径 | 本程实测 | 落地件 §1.3 自报 | 判 |
| --- | --- | --- | --- |
| `^=== RUN` | **142** | 142 | 对上 |
| `^--- PASS`（顶层） | **142** | 142 | 对上 |
| `^--- FAIL` | **0** | 0 | 对上 |
| `^--- SKIP` | **0** | 0 | 对上 |
| 子测试 `^    --- PASS` | **0**（顶层即全集） | — | 无嵌套稀释 |
| 去重顶层名数 | **71**（每枚恰好 ×2，`uniq -c` 只有 `2` 这一档） | 71 | 对上 |
| `grep -ci panic`（名带 Panic） | **8** | 8 | 对上 |
| 真 `^panic:` | **0** | 0 | 对上 |
| 包级行 | `ok github.com/CarlosShao/wisp/internal/observe 6.626s`，`rc=0` | — | — |

**名册是集合差集，不是只看四数**（简报明令不许只收包级 rc）：本程把自己 71 枚名册与落地件留在盘上的两枚
原始名册各做一次 `comm -3`——

```
$ comm -3 <本程 71 枚> D:\tmp\wisp141gate\logs\base.names   → 0 行
$ comm -3 <本程 71 枚> D:\tmp\wisp141gate\logs\post.names   → 0 行
```

（那两枚文件是 `PASS: TestX` 形，本程先剥前缀再比；剥后各 71 枚。）
⇒ **71 枚逐名三方守恒**：基线（gate=false、旧披露腿）、落地后（本程自己量的 `4ecc284`）、落地程自报名册，**同一份**。
无一名消失、无一名转 SKIP、无改名。

**"PASS 掉而 FAIL 不涨"这一维本程专门看了**：`--- SKIP` 两味（顶层与 `=== SKIP`）皆 **0**、`=== PAUSE` **0**、
RUN=PASS=142 三者相等 ⇒ **没有任何一枚用 SKIP 换色，也没有任何一枚被包级 `rc` 吞掉**。
本程**不以任何人的包级 rc 作判**：上表四数逐味是本程从自己的 `-v` 输出里 `grep -c` 出来的。

**档位：〔独立复现〕**。

---
