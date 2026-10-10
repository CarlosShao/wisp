# 299-a1 · 00-anchor — 只读普查腿射程声明

## 钟点与锚

- 钟（`date "+%Y-%m-%d %H:%M:%S +08"` 的 stdout）＝`2026-10-10 09:43:39 +08`
- 分支＝`dev`
- HEAD 对象层锚＝`6ef14788a0d1bc0c71962ba7342479069369d97f`（`git rev-parse HEAD`）
- `git log --oneline -1` ＝ `6ef14788 票 296 · 验收腿 296-v1 交回三格判语（AC#0 成立／AC#1 成立-夹具有牙／AC#2 部分成立）`
- 起手树态（`git status --porcelain -- cmd/wisp internal/panel frontend`）＝**空输出**（0 行）。
  ⚠ 这只是 09:43:39 那一瞬的快照；同包写腿 `293-r1` 据称在动 `cmd/wisp`，故本腿**全程不以工作树为真相源**。

## 射程

- 只读普查：不改任何跟踪文件、不新建任何非 `.md` 件、不翻票 299 任何 `- [ ]`。
- 本腿只新建 `.scratch/wisp/probes/299/a1/*.md` ＋ 对票面**追加**一条 Progress log。
- ⛔ 禁跑 go 编译面（`go build`／`go vet`／`go test`）。允许 `go env`／`go list`；跑过哪几条在本件末「自报」节逐条列名。
- ⛔ 不起 localhost HTTP／不开监听（D29 逐字禁）。
- ⛔ 读机主 config 值、⛔ 读 `%APPDATA%\wisp-dev\secrets\`、⛔ 任何凭据进件。
- ⛔ 读 `frontend/src/**` 内容；`frontend/dist/index.html` 只走 `git show HEAD:…`。
- ⛔ 零 push；commit 显式 pathspec；⛔ `add -A`／`amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`。

## 与派单的偏差声明（定则：读数只来自现跑的尺）

编排者在派单 §0 自述曾凭空写过"这枚腿已落三笔 commit、带回四条读数"。本腿**不采纳**那些历史读数：
下面每一枚数、每一枚行号、每一枚片段都由本腿自己在 HEAD 上现跑得到。凡与票面原文冲突，以原文为准并具名报回。

## 五把尺（逐字，本腿后文所有枚数/名册的可复现依据）

1. 运行时资源名册（件 `10-runtime-asset-roster.md`）：
   - `git show HEAD:frontend/dist/index.html` ＋ 对 `internal/panel/assets.go` 里 `entryRefRe` 的**同形正则**（正则本体逐字抄进件 10）
   - 对 js 产物再扫一层 `import(` / `new URL(...import.meta.url)` 形状
   - `git ls-tree -r --name-only HEAD -- frontend/dist`
2. 依赖侧注册 API 名册（件 `20-registration-api-inventory.md`）：
   - `go env GOMODCACHE` ＋ 在该模块目录内 grep（只引 file:line＋逐字片段，⛔ 大段源码进件）
   - `git show HEAD:cmd/wisp/panel_host_windows.go`（宿主手里握的是哪枚对象）
3. 两形代价表（件 `30-two-shape-cost-table.md`）：
   - `git grep -in 'WebResourceRequested' HEAD -- ':!.scratch'`（注册方现状）
   - `git grep -n 'go:embed' HEAD -- cmd/wisp internal/panel`
   - `git show HEAD:cmd/wisp/panel_assets.go`／`git show HEAD:internal/panel/assets.go`
4. CSP `<meta>` 名册（件 `40-open-questions.md` Q2）：
   - `git grep -n -i 'Content-Security-Policy\|CSP' HEAD -- ':!.scratch'`
5. 用例进不进 CI（件 `40-open-questions.md` Q3）：
   - `git show HEAD:<file>` 首行 `//go:build`
   - `git show HEAD:.github/workflows/ci.yml` 里对应那一步的名册

## 自报：本腿跑过的 go 命令（终态）

- **跑过、且只跑过这一条**：`go env GOMODCACHE GOFLAGS GOVERSION` ⇒ `D:\work\base\gopath\pkg\mod`／（`GOFLAGS` 空行）／`go1.27.1`。用途＝定位依赖模块目录（件 `20` 的尺①）。
- ⛔ **没跑 `go list`**（原计划用得上，最终靠 `git show HEAD:go.mod`＋`git grep -n 'jchv' HEAD` 就够，不必再派生一层读数）。
- ⛔ **没跑任何编译面**：`go build`／`go vet`／`go test` 全 0 发（同包写腿 `293-r1` 独占，本腿未碰）。⇒ 件 `40` Q3 里那句"AC#13 在今天的 CI 上会 skip"是**预测**，已具名标注，⛔ 不是读数。
- 其余全部读数来自 `git show HEAD:<path>`／`git grep -n … HEAD`／`git ls-tree`／`git check-ignore`／`wc -c`，以及对**工作树构建产物** `frontend/dist/**` 的 grep（该件不在 HEAD 对象层，理由与后果见件 `10` §0）。
- 依赖模块目录内的 grep＝对 `GOMODCACHE` 里 `github.com/jchv/go-webview2@v0.0.0-20260205173254-56598839c808` 的只读 grep，⛔ 未在该目录内新建／修改任何文件；大输出（49 行 API 名册／248 行 webview2 命中／37 行 CSP 命中）先落 `$TEMP/wisp-299-a1/`（仓外），⛔ 未落仓内、⛔ 未生成 `.out`。
