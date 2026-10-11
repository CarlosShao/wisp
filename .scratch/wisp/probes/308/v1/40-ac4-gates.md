# 308-v1 · 40-ac4-gates（`AC#4` 三门禁＋格式名册两把射程＋逐笔越界；分母那半截在 10 件）

单位＝每发自己的 `rc=`；件⛔ 空（build/vet 绿的时候 stdout 本来就是 0 字节 ⇒ 本腿把 `rc=N` 追进**同一枚件**，⛔ 让下一位读者拿 0 字节当"没交"）。

## 1. 三把门禁（本腿现跑，⛔ 转抄 `probes/308/r1/logs/gate-*`）

| 命令 | rc | 件 |
|---|---|---|
| `GOFLAGS= go build ./...` | **0** | `logs/v30-gate-build.txt`（末行 `rc_build=0 (0 bytes of stdout IS the pass shape; this line is the reading)`） |
| `go vet ./internal/risk/` | **0** | `logs/v31-gate-vet.txt`（同款 rc 行） |
| `sh scripts/d22scan.sh` | **0** | `logs/v32-gate-d22scan.txt` 23 084 字节，末行逐字＝`d22scan: clean - no D22 ban violations; live scope work: bans #1-5 internal/=229, bans #1-5 cmd/=39, ban #6 frontend/=85, ban #7 internal/tools/=23, ban #8 design/=39, ban #8 frontend/=85, ban #8 internal/=531, ban #8 cmd/=119; …` |
| （补）`sh scripts/check-path-length-budget.sh --with-self-test` | **0 GREEN** | `logs/v34-pathlen.txt`：`denominator: tracked paths=10245 over-budget=57 covered by roster=57 not in roster=0`；`worst full path=228 chars (hat budget 165)`；`VERDICT GREEN - every over-budget tracked path is rostered by name with a reason, and the roster equals the tree` ⇒ 本腿新件（最长名 `v26-docker-linux-full.txt`＝25 字符）⛔ 顶它 |

⚠ `go build ./...` 是全仓尺，本发⛔ 命中任何 `internal/tools/**` 红句（309-r1 在飞面）；若下一发命中，按票规矩＝具名记那发红句＋稍后重试一发，⛔ 归给本票。
`d22scan` 射程⛔ 覆盖 `.scratch/**` ⇒ 本腿写面⛔ 在仪器射程内（⛔ 意为⛔ 受它管：emoji 尺＝本腿三枚件＋本目录⛔ emoji）。

## 2. ★格式名册两把并排，各写射程与单位（CRLF 假枚这一族）

`git config core.autocrlf` ＝ **true**（件 `logs/v33-format-rosters.txt` 同页）⇒ 工作树里的 `w/crlf` 是** checkout 产物**，⛔ 等于 blob 脏。尺＝对每枚命中先数 **CR 字节**（`tr -dc '\r' | wc -c`）再读 `git ls-files --eol`。

| 射程（尺） | 命中 | 逐枚 | 单位／性质 |
|---|---|---|---|
| 窄＝`gofmt -l internal/risk`（工作树） | **1** | `internal\risk\provenance.go` | CR 字节＝**1115**，`eol=i/lf w/crlf attr/text eol=lf` ⇒ **CRLF 假枚**（索引侧是 lf），⛔ 本票引入、⛔ 本腿修 |
| 宽＝`git ls-files -z '*.go' \| xargs -0 gofmt -l`（全仓 tracked） | **26** | `.scratch/**` 22 枚（⛔ 干净是它们的职责＝变异拷贝）＋非 `.scratch` **4** 枚：`cmd/wisp/models.go`、`internal/agent/approval/pending_read.go`、`internal/agent/tools.go`、`internal/risk/provenance.go` | 同形 CRLF／刻意脏，⛔ 本票账 |
| 三枚交付件＝`gofmt -l <anchor\|junction\|volumeless>`（工作树） | **0** | — | CR 字节均＝**0**、`w/lf` ⇒ 交付自身格式干净 |

⚠ **枚数与上一腿⛔ 相等，且本腿知道为什么**：`308-r1` 报宽尺 **27**（23＋4），本腿报 **26**（22＋4），差在 `.scratch` 那一桶。非-`.scratch` 那一桶⛔ 动（4⇄4）⇒ 少的那一枚是**别的腿在两发之间动了 `.scratch/**` 里的 tracked `.go`**（本腿 13:5x ⇄ r1 13:1x）。⇒ 这条正是"分母随写面类型移动"的实例：**只有新增 `.go` 才挪这把尺**，本腿写面 `.md/.txt/.log` 对它＝0 影响（⛔ 在 `probes/**` 造 `.go`，本腿守住了）。
票 `A839` §1 末行那句"窄＝1 枚／放宽到 `internal cmd`＝5 枚"本腿未复跑该第三把射程，⛔ 用它顶上面两把，⛔ 说它错。

## 3. 逐笔越界（尺＝`git show --name-only --format= <sha>` 滤掉自家写面）

五笔滤后各＝**0 行**（`grep rc=1`＝无命中；件 `logs/v05-numstat.txt` 逐笔带 `rc_outside=1`）。
写面清单：`510c95ca`＝`probes/308/r1/00-anchor.md`；`1ad48333`＝3 枚 `internal/risk/*_test.go`；`ee0a93e4`＝`probes/308/r1/{10,20,30,40}*.md`＋29 枚 `logs/*.txt`；`79242df7`＝`90-final.md`；`1270e36e`＝`40-tags-and-gates.md`。
★**⛔ 越界但⛔ 完全没碰票文件**：票 308 末尾那枚 r1 bullet 由**编排者 13:3x 代提**（票面 §"agent=编排者（代提，正文归 308-r1）"），r1 自己在 `90-final.md` §①/§③ Q-308-r1-D 里具名报了为什么⛔ 提交（同一枚文件此刻躺着编排者在飞态）。本腿复核：`AC#0` 那枚 `[x]` 的归属＝编排者 13:1x 现跑翻的（`git log` 里那三笔 `c023e8e0`/`dd638f98`/`24929768` 归 a1/裁定），⛔ r1 翻框。⇒ 那一笔停手上报＝**做对了**，本腿给它背书。

## 4. 本格判语

**成立**。三门禁各 rc=0 且各件自落 `rc=N`（⛔ 0 字节空件）；格式名册两把并排、各写射程与 CR 单位、⛔ 一枚被误报成脏；逐笔越界＝0；长度门⛔ 绿着，本腿写面⛔ 顶它。
（票 `AC#4` 里"⛔ push／commit 显式 pathspec／写面＝自家 probes＋裁决表＋票尾 1 bullet"这一族是对**落地腿**的约束，r1 逐条守住，见上；本腿同样按此走，自证见 `90-final.md` §4。）
