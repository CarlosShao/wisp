# 309-v1 · `AC#4` 门禁与越界（三枚 rc＝本腿自己跑；格式名册两档并排；eol 复量）

台面＝HEAD `0d2260eb`，工作树（母仓），`GOFLAGS=` 显式清空。每枚门禁件自己落一行 `rc=N`，⛔ 0 字节件。

## 1. 三枚门禁（原始＝`logs/gates-rc.txt`＋`logs/gate-build.log`／`gate-vet.log`／`gate-d22scan.log`）
```
GOFLAGS= go build ./...        -> rc_build=0      （件 0 字节＝无诊断，已就地补 rc 行）
GOFLAGS= go vet ./internal/tools/ -> rc_vet=0     （同上）
sh scripts/d22scan.sh          -> rc_d22scan=0    （件 23078 字节，末行逐字 clean）
```
`d22scan` 末行逐字（射程自报）＝`d22scan: clean - no D22 ban violations; live scope work: bans #1-5 internal/=229, bans #1-5 cmd/=39, ban #6 frontend/=85, ban #7 internal/tools/=23, ban #8 design/=39, ban #8 frontend/=85, ban #8 internal/=531, ban #8 cmd/=119 …`
⇒ 与腿报的 `internal/=229`／`cmd/=39`／`internal/tools/=23`／`ban#8 internal/=531` **逐枚对得上**（本腿自己那发，⛔ 抄它）。
⛔ 一枚门禁⛔ 跑过：`go test ./cmd/wisp/`⛔ 本格判据（那是票 298 `AC#4` 的账，编排者自己按住了）。
并发面自报：`git status --porcelain -- internal cmd`＝**0 行** ⇒ 我这几发没有踩到 `308-v1` 在 `internal/risk` 的在飞态；三枚 rc 各**只跑了一发就=0**（⛔ 需要第二发反证，⛔ 出现"别人文件名"的红）。

## 2. eol 那段自抓自纠，本腿复量（三把尺并排）
| 尺 | 读数 |
|---|---|
| `tr -cd '\r' < internal/tools/bridge.go \| wc -c` | **0**（CR 枚数；腿自报曾 `CR=1313`＝全文每行都带 CR，行数 1313 对得上） |
| `git ls-files --eol -- internal/tools/bridge.go` | `i/lf    w/lf    attr/text eol=lf` |
| `.gitattributes` 规则 | 第 4 行逐字 `*.go text eol=lf` |
⇒ **自抓自纠坐实**：今日工作树文件＝LF、索引＝LF、属性＝eol=lf，三者一致。
旁证（同一把尺打到别的文件上，说明这台机器确实造得出 CRLF 假枚）＝下面名册里那 4 枚全为 `w/crlf`。

## 3. 格式名册：两把并排＋射程具名（⛔ 混着数）
`gofmt`（在 PATH）／`gofumpt`（`$(go env GOPATH)/bin/gofumpt.exe`，本腿**找到了**，⛔ 拿"⛔ 工具"当免检）：

| 尺 | 射程 | 工作树 | HEAD blob 导出树 | 归 309-r1 名下新增 |
|---|---|---|---|---|
| `gofmt -l` | `internal/tools` | **0** | **0** | **0** |
| `gofmt -l` | `internal cmd`（宽域） | **4** | **0** | **0** |
| `gofumpt -l` | `internal/tools` | **0** | （blob 侧那 1 枚⛔ 在本包） | **0** |
| `gofumpt -l` | `internal cmd`（宽域） | **5** | **1**＝`cmd/wisp/panel_transport_35r2_test.go` | **0** |
| `git ls-files '*.go' \| xargs gofmt -l` | 全仓 tracked（分母 **995**） | **26 ＝ 22 `.scratch` ＋ 4 非** | — | **0** |

宽域工作树那 4/5 枚逐枚属性＝`i/lf w/crlf attr/text eol=lf`：`cmd/wisp/models.go`／`internal/agent/approval/pending_read.go`／`internal/agent/tools.go`／`internal/risk/provenance.go`（第 5 枚 gopumpt 才多抓的＝`panel_transport_35r2_test.go`）。
⇒ 四条⛔ 一枚住在 `internal/tools/**`，⛔ 一枚由 `16e86bae` 引入（该笔只碰 `bridge.go`＋两枚新 `_test.go`，而那三枚文件在两把尺下都是 0）。
blob 侧那 1 枚（`panel_transport_35r2_test.go`）＝**票 298 `AC#5` 已承接的既有债**（台账 `A846` §(b) 末条逐字"已由我在票 298 新立 AC#5 承接"）⇒ ⛔ 记到 309，⛔ 因它判红。
最后一把全仓尺与 `A846` §(c) 那发（`26＝22＋4`，分母 995）**逐枚对得上**，且 `bridge.go` 已⛔ 在那 4 枚里＝CRLF 归一**留住了**。

## 4. 越界（写面）
`git show --name-only --format=` 逐笔：`1cd57ec1`＝1 枚（`probes/309/r1/00-anchor.md`）；`16e86bae`＝3 枚（`internal/tools/bridge.go`＋两枚 `internal/tools/*_309_*_test.go`）；`b2dd7fb2`＝25 枚（`probes/309/r1/**` ＋票 309 末尾 1 行）；`8ffa61c6`＝1 枚；`8bbffc50`＝1 枚。
⇒ 全部落在票面许可的三张写面（`internal/tools/**` ⇄ `probes/309/**` ⇄ 本票末尾 append-only bullet）内，**越写面＝0 枚**。
尺＝`for c in …; do git show --name-only --format= $c; done` 滤 `^internal/tools/|^\.scratch/wisp/probes/309/|^\.scratch/wisp/issues/309-` 取反计数＝0。
票 bullet 那一笔的 numstat＝`1 0`（1 加／0 删）已在 `b2dd7fb2` 的 `--numstat` 里核过。

## 5. 本腿⛔ 碰的东西（自查）
`tools/d22scan/allowlist.txt`／`internal/observe/thresholds.go`／`testdata/golden`／D43 转移表／C1–C32／D1–D47／三枚冻结件：⛔ 读⛔ 写⛔ 判；
本腿写的字节只有 `.scratch/wisp/probes/309/v1/**`（＋其 `logs/`）与 `docs/evidence/s1/309-v1-verdict.md` 与票 309 末尾那一行。
提交前 `git -c core.quotePath=false diff --cached --name-only` 逐名核，⛔⛔ 名册里出现别人名下的件；提交一律显式 pathspec。
