# 280-r1b 三格交付：7 枚名册的归属普查（尺＝gofmt go1.27.1 ＋ gofumpt v0.12.0，byte 源＝HEAD 快照）

腿：`280-r1b`。起手锚 `00-anchor.md`（同目录，commit `f1a43277`）。锚点 `2bdf385d`，分支 `dev`，取数时刻 `2026-10-09 10:1x +0800`。

## 0. 射程（两把尺、两个 byte 源，必须分开读）

| 维度 | 本程的尺 | 射程 |
|---|---|---|
| 名册**枚数** | `gofmt -l internal cmd tools`，**工作树** | 现跑＝**7 枚**，与票面 `cmd/wisp/models.go`／`panel_inbound_guards_35r3_test.go`／`panel_transport_35r2_test.go`／`internal/agent/approval/pending_read.go`／`internal/agent/tools.go`／`internal/risk/provenance.go`／`internal/tools/bridge.go` **同名同数**；票面原射程 `gofmt -l cmd internal` 亦 7 枚 ⇒ `tools/` 目录贡献 **0 枚**（与 `280-r1` 独立复现一致） |
| 差异**内容** | `git show HEAD:<path> \| gofmt -d`／`\| gofumpt -d`，**HEAD blob** | 7 枚逐条，见 §1。**零**工作树差异字节被用作判据 |
| gofmt | `/d/work/base/go/bin/gofmt`（go1.27.1 自带） | `-l`／`-d` 只读 |
| gofumpt | `D:\work\base\gopath\bin\gofumpt.exe`，`--version` = `v0.12.0 (go1.27.1)`；**不在 PATH** | `-d` 只读；零 `go install`、零下载、零 `-w` |
| 工具射程 | 两把尺都**只吃单枚文件**（stdin／具名路径），本程**没有**对全仓跑 `gofumpt -l .` | 带不带 `_test.go`＝**带**（名册里 2 枚是 `_test.go`）；`.scratch/**` 本程**不入判据** |

辅助读数（同一把 `sed -n '1,12p'`／`tr -dc '\r'` 取数）：`core.autocrlf` = `true`；`HEAD:.gitattributes` 存在且含 `* text=auto`／`*.go text eol=lf`；`git status --porcelain -- <这 7 枚>` = **空**（取数那一刻 7 枚工作树相对 HEAD 无字节差，含正被另一枚腿占着的 `internal/agent/approval/pending_read.go`，⚠ 时点读数）。

## 1. AC#1 逐枚判归属（差异形状逐处写明；gofumpt 是否同声具名）

**头条发现：票面那句"7 枚文件不过裸 `gofmt`"在 HEAD 上只成立 2 枚。** 其余 5 枚在 HEAD blob 上两把尺都**完全干净**（`gofmt -d` 0 行、`gofumpt -d` 0 行），它们被工作树点红的原因是**行尾**，不是风格：

| 枚 | HEAD 行数 | HEAD gofmt -d 行数 | HEAD gofumpt -d 行数 | 工作树 CR 字节 | HEAD blob CR 字节 |
|---|---|---|---|---|---|
| `internal/agent/approval/pending_read.go` | 131 | **0** | **0** | 131 | 0 |
| `internal/agent/tools.go` | 225 | **0** | **0** | 225 | 0 |
| `internal/risk/provenance.go` | 1115 | **0** | **0** | 1115 | 0 |
| `internal/tools/bridge.go` | 1284 | **0** | **0** | 1284 | 0 |
| `cmd/wisp/models.go` | 334 | **0** | **0** | 334 | 0 |
| `cmd/wisp/panel_inbound_guards_35r3_test.go` | 171 | 14 | 14 | 0 | 0 |
| `cmd/wisp/panel_transport_35r2_test.go` | 2142 | 14 | **101** | 0 | 0 |

（CR 那两列是**行尾诊断**，不是 gofmt 差异内容；本程没有读过那 5 枚的工作树 diff 字节。"这 5 枚只差行尾"是由「HEAD blob 两把尺 0 行」＋「工作树 CR＝行数、HEAD blob CR＝0」＋「`git status` 判这 7 枚未改」三条**合推**出来的，具名写成推论。）

逐枚归属：

1-5. 上述 5 枚：这一处的差异**既不属于 `gofmt` 风格也不属于 `gofumpt` 加严**＝**行尾 artifact**（Windows 工作树字节带 CRLF，`gofmt`/`gofumpt` 输出恒为 LF，故必点红；`core.autocrlf=true` 让 git 用 clean 过滤器归一化，所以 `status` 反而判"未改"）。⇒ 票面把这 5 枚算作"未格式化"是**尺的射程错**，不是文件的债。

6. `cmd/wisp/panel_inbound_guards_35r3_test.go` —— **gofmt 风格（对齐）**，一处、1 个 hunk：`@@ -63,8 +63,8 @@`，`const ( … )` 组里三枚常量的**尾随 `//` 注释对齐列**。gofmt 要把 `guard35r3Source`／`guard35r3RequestID` 两行 `//` 前的空格各减一格（值本身是含 CJK 的反引号 raw string，`guard35r3Roster` 那行不动）。load-bearing 原文（本腿 `gofmt -d` 现量）：

```
-	guard35r3Source    = `按伪造/串台拒绝`               // bridge.go:136-137 (…)
-	guard35r3RequestID = `缺少 requestId`                // bridge.go:140-141 (…)
+	guard35r3Source    = `按伪造/串台拒绝`              // bridge.go:136-137 (…)
+	guard35r3RequestID = `缺少 requestId`          // bridge.go:140-141 (…)
```

`gofumpt -d` 对它：**同声**——也是 14 行、同一个 hunk、逐字相同（对齐是 gofmt 的债，gofumpt 只是继承）。空行／`//` 后空格／分组三类加严形状**一枚都没有**。

7. `cmd/wisp/panel_transport_35r2_test.go` —— **混合**，且 gofmt 那把尺**低估**它：
   - `gofmt -d` 只有 14 行＝1 个 hunk `@@ -804,8 +804,8 @@`：两枚单行方法声明的**单行体前多余对齐空格**——`func (p *jsParser) peek() jsTok    { … }`／`next()` 要各压成一格（下面 `atPunct` 是多行声明，不再成对齐组）。⇒ 这一处＝**gofmt 风格**，gofumpt 同声。
   - `gofumpt -d` 是 **101 行**（新增行 34 枚），除上面那一 hunk（在 `@@ -804,16` 里逐字包含）之外**多出三类 gofumpt 独有加严**，gofmt 对这些**完全 0 发声**：
     - `@@ -400,18 +400,22`、`@@ -421,12 +425,14`：连续单行 `type` 声明并入 `type ( … )` **分组**（`jsIdent`／`jsLit`／`jsMember` 与 `jsIndex`／`jsCallNode`；`jsAssignNode`／`jsPostInc`／`jsFuncExpr`）。
     - `@@ -698,8 +704,10`：跨行写的 composite literal 重排为**每行一枚元素**（`var jsPuncts = []string{ … }`）。
     - `@@ -804,16`／`@@ -821,6` 一带的 `+` 空行：相邻**顶层声明之间补空行**（`atPunct`／`atID`／`eatPunct`／`expectPunct` 之间）。

## 2. AC#2 与 CI 的关系（判据＝本腿读到的 CI 定义原文；行号与短语同一次取数）

尺：`git show HEAD:.github/workflows/ci.yml | grep -nE "gofumpt|gofmt|lint|staticcheck|vet"` ＋ `sed -n '160,245p'`；`wc -l` = **1015**。

- `:172` `- name: gofmt (gofumpt)` → `:173` `if: ${{ !cancelled() }}` → `:175` `go install mvdan.cc/gofumpt@latest` → `:176` `OUT="$(gofumpt -l . tools/d22scan tools/mockllm)"`，`[ -n "$OUT" ]` 则打印名册并 `exit 1`。⇒ **scope＝仓根 `.` 全树**（含 `.scratch/**`、含 `_test.go`），**带 `if:` 守卫**。
- `:189` `- name: "gofmt (gofumpt) - the tracked set is the denominator (ticket 161 AC#7 form A)"` → `:239` `if: ${{ !cancelled() }}` → `:240` `run: sh .scratch/wisp/probes/161/r5/attrib.sh --tracked-only`。该脚本（`HEAD:` 现量注释 `:37`）的形 (A)＝`git ls-files '*.go' | xargs gofumpt -l`＝**tracked 全集分母**，零 tracked `.go` 时 `exit 2`（不许空过）。⇒ **同样带 `if:` 守卫**（这一枚是 `111-r5` 补的，票 275 具名）。
- ⚠ **"这一步在 yaml 里" ≠ "它跑过"，本程的射程就到此为止**：本程不许跑 `gh`、无 run 观测取数 ⇒ 我对锚点 `2bdf385d` 这两步**从未求值**，绝不能读成"拦得住"。仓内既有记录（本腿 `git grep HEAD` 现量，都是别人已落账的观测，不是我看的 run）：票 111 `:310` 记 job `112053738745` step8 `gofmt (gofumpt)` **failure**；票 111 `:394` 记某 run 它红**只有 1 条命中**＝`.scratch/wisp/probes/185/c1/mut/fs_broken.go:4:1`（tracked 的故意坏样本，票 269 的头病），且普查步与两枚 `go vet` 当时"12 发里全是 `skipped`"；票 276 记普查步在当前 dev 树上**必然 rc=1／19 枚、全部住 `.scratch/**`、产码命中 0**（10-07 的账）。⇒ 结论口径写成"**会不会进那一步的 OUT 名册**"，不写"门会不会为它变红"：那扇门今天**已因别的字节常红**。
- ⚠ 版本没钉：CI 是 `gofumpt@latest`（`:175`；票 124 `:251` 早具名同一处），本程的尺是 `v0.12.0`。名册只会随新版本变长不会变短 ⇒ 下面"会进名册"是**下限**；"不进名册"的 5 枚则带一条具名残余风险（更新的 gofumpt 若加严行尾以外的形状，本程判不到）。
- 顺带核对编排者那句转述（本腿现量一致）：`:606` `- name: "winlive compile gate (go vet -tags winlive, ticket 111 AC#11)"`、`:655` `run: go vet -tags winlive ./cmd/wisp/ ./internal/ball/` ⇒ **只编译／类型检查，不执行测试**，与格式步无关；本 7 枚的债**不经**这一步。

逐枚：

| 枚 | 会不会被 CI 点红（进 OUT／分母名册） | 凭据 |
|---|---|---|
| `pending_read.go`／`tools.go`／`provenance.go`／`bridge.go`／`models.go` | **不会**。两步都不点：HEAD blob 上 gofumpt v0.12.0 对它 0 发声；ubuntu 检出按 `.gitattributes` `*.go text eol=lf` 给 LF ⇒ CRLF artifact 在 runner 上根本不存在 | `:176` scope、`:240` 分母、CR 计数、`.gitattributes` |
| `panel_inbound_guards_35r3_test.go` | **会**（tracked、住 `cmd/**` ⇒ 同时进 `:176` 的 `.` 与 `:240` 的 tracked 分母；gofumpt 对它发声） | §1 第 6 条 hunk |
| `panel_transport_35r2_test.go` | **会**，且**只有 gofumpt 那把才看得见它大半**：CI 里没有"裸 gofmt 那一档"，`:176`/`:240` 都是 gofumpt ⇒ gofmt 的 1 处之外，type 分组／literal 重排／顶层空行三类加严**同样被 CI 点** | §1 第 7 条 hunk |

⇒ 三格合并的一句话：**"CI 真会拦"＝0 枚有独立可证的"拦"（那门今天为别的字节常红，本程又无 run 观测）；按名册口径"会进 CI 名册"＝2 枚；"只有裸 gofmt（本机工作树）看得见"＝5 枚，且它们在 HEAD 上连裸 gofmt 都不该看见（行尾）**；"判不动"＝0 枚，唯一残余不确定性＝`@latest` 版本漂移（已具名）。

## 3. AC#3 处置建议（只写形状；本程零改动、零 `-w`）

- 那 5 枚（`pending_read.go`／`tools.go`／`provenance.go`／`bridge.go`／`models.go`）：**不改**。HEAD 字节两把尺全干净，去"修"它们只能改写行尾＝无的放矢，且 `internal/agent/approval/**` 正被另一枚腿种突变。**留给下一枚真正动该文件的票**（若那一枚的编辑器把 CRLF 写回，工作树命中会再现，届时按行尾处置，不属本票）。
- `panel_inbound_guards_35r3_test.go`：**改**（一处对齐）。理由＝两把尺同声＝真·未格式化，且它已在 CI 名册里；改动面 2 行、不碰断言文案。⛔ 本程不动。
- `panel_transport_35r2_test.go`：**留给下一枚真正动该文件的票**。理由＝要清空 gofumpt 名册必须一次做完 5 处声明区形状（type 分组／literal 每行一枚／顶层补空行）＋那 1 处 gofmt 对齐；只挑 gofmt 那 1 处改＝CI 照旧红，改两遍。⛔ 本程不做批量重排。
- 与本票**同一堵墙**但不在本票射程（具名不改）：`:176` 那步的 scope 含 `.scratch/**`（票 269／275／276 三票的既有账），本程不裁、不扩、不建议豁免。

## 4. 具名报回（以盘上原文为准）

1. **票面前提需重读**：票 280 标题与 `:8` 的"7 枚"来自 `A736` 的**工作树**快照（`gofmt -l cmd internal`）。在 HEAD 上只有 2 枚成立。本程**不改票面**，只把这一条按"5 枚＝行尾 artifact"报回编排者。
2. **对上一程 `280-r1` 的一处更正**：它 `:32-33` 记「尺 C：`gofumpt -l <这 7 枚>` ＝ 7 枚全中」并据此写"本票的问题不是 gofumpt 是否也点这 7 枚（点了）"。那把尺吃的是**工作树字节**（含 CRLF）；换成 **HEAD blob** 后 gofumpt 只点 **2 枚**。两把读数不矛盾、射程不同；照搬那句到 HEAD 就错。
3. **我顶回的任务书转述**：编排者要"逐枚答它在 CI 那一步会不会被点红"时预设那一步会求值。本程在 HEAD 上**无任何 run 观测**（禁 `gh`、无网取数），且仓内三票记录那两步**常红／曾全 skipped** ⇒ 我把判据落成"会进 OUT／分母名册"而不是"拦得住"，并具名写"这一步本程从未求值"。
4. 未顶回的：`winlive` 只编译不执行（`:606`／`:655` 本腿现量一致）；`tools/` 射程贡献 0 枚（现量一致）；gofumpt 在 GOPATH/bin 不在 PATH（本腿独立定位到同一枚二进制、同一版本）。

## 5. 禁区自查

写面只有 `.scratch/wisp/probes/280/r1b/*.md`（本件＋锚）。零 `-w`、零批量重排、零产码字节改动；票面／台账／AC 框零触碰；`frontend/**`／`design/**` 零读零写零转述；三枚冻结件本程未读未量；`PLAN.md`／`docs/specs/**`／`internal/observe/thresholds.go`／golden 一字未动；未跑 `go build`／`go test`／`go vet`／`go list`／`go doc`；未 `go install`、未下载；工作树里别人的在飞改动（`M .gitignore`、`design/**` 删除、大量 `probes/**`）一律未触碰；只 commit、零 push、commit 带显式 pathspec；文案零 emoji。
