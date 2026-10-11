# 308-v1 · 00-anchor（非实现者终裁腿的起手尺；读数一律本腿现跑，⛔ 转抄 `308-r1`/`308-a1`）

- 腿＝`308-v1`（验收腿，⛔ 我写过的码我来判）；票＝`.scratch/wisp/issues/308-spelling-invariance-test-panics-when-temp-has-no-volume.md`
- 起手 HEAD＝`0d2260eb10243df10fc925d766edbc64bedf98e8`，分支 `dev`，共享工作树 `D:/work/workspace/projects plans/Wisp`
- 被裁交付＝五笔 `510c95ca`→`1ad48333`→`ee0a93e4`→`79242df7`→`1270e36e`（作者 CarlosShao，13:02–13:23 +0800）
- 起手 `git status --porcelain -- internal/risk`＝**空**（工作树⇄HEAD 对 `internal/risk` 等价 ⇒ 本腿所有工作树读数可直接归给 HEAD）
- 起手票 308 框尺＝`grep -cE '^[[:space:]]*- \[x\]'`＝**1**／`- [ ]`＝**5**／行数＝**82**（＝编排者 13:1x 裁语后的态；本腿全程⛔ 动框，只许末尾追加 1 行）

## 射程声明（每格都认这三层，⛔ 混说）

| 层 | 是什么 | 本腿用它答什么 |
|---|---|---|
| 工作树 | `D:/work/workspace/projects plans/Wisp`（共享，别的腿在飞 `internal/tools`） | 门禁、名册、突变（每次当场还原＋逐字节证） |
| HEAD blob | `0d2260eb` 的 `internal/risk/*` | 产码⛔ 动的证据（hash 比对），⛔ 任何"工作树看着没改"的口径 |
| 仓外导出树 | `git archive --format=tar <sha> \| tar -x` 到 `mktemp -d`，**一树一枚新目录**，量之前先数 `.go` 证非空 | 改前⇄改后对照（⛔ 在共享工作树里跑"改前那发"） |

导出树两枚（本腿自建、⛔ 仓内）：
- PRE＝`1ad48333^`（＝码落地前）路径 `go.mod go.sum internal/{risk,observe,plugin,secret,winsec,tools}` → `.go`＝**193**，`internal/risk/*.go`＝**42**
- POST＝`HEAD` 同路径 → `.go`＝**196**，`internal/risk/*.go`＝**43**（＋1＝新件）
- ⚠ 第一发本腿只导了 5 个内部包，`TestC26RewriteAccountIsConsumedAtEverySecurityLeg` 因读不到 `..\tools\paths.go` 各报一枚红（件 `logs/v10-pre-tree-full.txt`/`v11-post-tree-full.txt`＝**射程假象，⛔ 交付缺陷**）⇒ 补 `internal/tools` 重导两枚干净树（`v12`/`v13`，两发 rc 均＝0）。**先证分母再信"0 问题"这条本腿自己踩了一遍，逐字留在件里。**

## 单位

- 名册枚数＝`grep -c` 行数（顶层尺＝`^--- (PASS|FAIL|SKIP): `，子层级尺＝`^    --- PASS`，两把⛔ 混）；
- CR 字节＝`tr -dc '\r' | wc -c`（工作树侧）⇄ `git ls-files --eol` 的 `i/? w/?`（索引侧）；
- 路径长度＝字符数（门 `check-path-length-budget.sh` 自报 byte/char 同选 57 枚）；
- 时间＝+0800。

## 本腿编译面命令自报（⛔ 瞒）

`go test ./internal/risk/` **21** 发＝工作树 12（`v-target-case`/`v0-full`/`teeth-pristine`/`teeth-mut-target`/`teeth-mut-full`/`teeth-restored-full`/`teeth-restored-target`/`teeth2-M2`/`teeth2-M3`/`teeth2-restored-again`/`teeth3-M3b`/`teeth3-restored`）＋仓外导出树 7（`v10`/`v11` 窄导那两发、`v12`/`v13`、`v40`/`v41` 别名台面、`v45` 分母探针）＋容器内 2（`v26`/`v27`，GOOS=linux）；其余＝`go test -list` 3 发（⛔ 执行任何码）。／`go build ./...` 1／`go vet ./internal/risk/` 1／`go vet -tags winlive ./internal/risk/` 1／`GOOS=linux go vet` 1／`GOOS=linux go test -c` 1／`sh scripts/d22scan.sh` 1／`sh scripts/check-path-length-budget.sh --with-self-test` 1／`go test -list` 3 发（⛔ 执行任何码）／`docker run golang:1.27 go test ./internal/risk/` **2** 发（linux 真色 pre⇄post）／仓外探针程序 1 枚（`mktemp -d` 里 `go build`，只读 env，自删自建；产出的 7 MB linux 测试二进制已迁出仓，见 `logs/v23-binary-relocated.txt`，⛔ 留 `.go`/二进制在 `probes/**`）。

`cmd/wisp` 测试包＝**⛔ 跑过 0 发**（309-v1 在飞＋sherpa DLL PATH 归编排者）。
