# 票 85 预做（pre-flight）：staticcheck 这道门到底存在过没有

- 代理：`audit-85-preflight`（**只读**预做，不改码、不改票面、不 commit、不 push）
- 时间：2026-09-21 20:35 起（`date` 实测 `Mon Sep 21 20:35:39 CST 2026`）
- 基线：**只用 `git archive HEAD` 解出的纯净树** `/tmp/a85-swq85`，HEAD = `3c5d1c3`。
  工作树的 `.github/workflows/ci.yml` 与 `scripts/` 归票 111，**本次未读、未改**。
- 本机负载声明：**同时有 3 个写码代理在编译/跑测** ⇒ 本文所有耗时数字**仅供量级参考，不是空载机器时**。

---

## 0. 结论速览（先给判据，再给过程）

| 待验命题 | 判定 | 档 |
|---|---|---|
| "staticcheck 从来没有产出过任何一条 **finding**" | **成立** —— 全历史 121 枚有 job 的 run：`success` **0 次**；11 枚逐字日志 finding 行 **0 条**；本机复现 0 条 | 〔日志＋归档，我抽验〕+〔独立复现〕 |
| "它一直在工具链不匹配的状态下**跑过去了**" | **不成立 —— 请你按下面纠正** | 〔日志＋归档，我抽验〕 |
| "这道门从未存在过" | **成立，但机制不是你以为的那个**：93 次被上面红步骤**跳过**、28 次**响亮地崩红**，**不是"假绿"** | 〔日志＋归档，我抽验〕 |

**要纠正你的那一句**：这一步的 `conclusion` 在**每一枚**采样到的 run 上都是 `skipped` 或 `failure`，**从来没 `success` 过**，
也从来没"跑过去"。日志末尾是 `##[error]Process completed with exit code 1.`。
⇒ 这不是"门开着但没人看见"，是"门一直红着、而 CI 本来就一路红（`lint` job 121/121 全红），于是这个红被当成背景噪声略过了"。
这两者的修法不同：前者要先证明门**能**绿，后者只要修工具链。（详见 §1.3、§3.7）

---

## 1. 证伪优先：staticcheck 到底报了什么

### 1.1 取数口径与"正向形状判据"

- 端点：`gh api repos/CarlosShao/wisp/actions/runs/<run>/jobs` → 取 `job.name=="lint"` 的**逐步骤** `conclusion`。
  这一步是**字段真解析**（`[.steps[]|select(.name=="staticcheck")][0].conclusion`），
  不是"输出里找 error 字样"。取不到字段时记 `STEP-ABSENT`，取不到响应时记 `FETCH-FAILED` 并重试 4 次。
- 日志：`gh api repos/CarlosShao/wisp/actions/jobs/<job>/logs`。
  成功判据用**正向形状**：`rc==0` **且** 落盘字节数 > 0（实测 4 枚：52624 / 52411 / 37386 / 34596 bytes）。
  ⚠ 未按"没看到 error 就算成功"判——那条我在 TLS handshake timeout 上栽过。
- 采样 run 跨 2026-09-19 23:44 → 2026-09-21 12:33（run 1 → 167，全仓 CI 历史全集就这么多）。

### 1.2 逐字输出（run 167 / job 106334066361 / 静态检查步骤）

```
##[group]Run go install honnef.co/go/tools/cmd/staticcheck@2025.1.1
go install honnef.co/go/tools/cmd/staticcheck@2025.1.1
staticcheck ./...
shell: /usr/bin/bash -e {0}
##[endgroup]
go: downloading honnef.co/go/tools v0.6.1
go: downloading golang.org/x/tools v0.30.0
go: downloading golang.org/x/exp/typeparams v0.0.0-20231108232855-2478ac86f678
go: downloading github.com/BurntSushi/toml v1.4.1-0.20240526193622-a339e1f7089c
go: downloading golang.org/x/sync v0.11.0
go: downloading golang.org/x/mod v0.23.0
go: downloading golang.org/x/exp v0.0.0-20231110203233-9a3e6036ecaa
-: internal error in importing "internal/byteorder" (cannot decode "internal/byteorder", export data version 4 is greater than maximum supported version 2); please report an issue (compile)
-: internal error in importing "internal/cpu" (...) ... (compile)
-: internal error in importing "internal/goarch" (...) ... (compile)
-: internal error in importing "internal/goos" (...) ... (compile)
-: internal error in importing "math/bits" (...) ... (compile)
-: internal error in importing "unicode/utf8" (...) ... (compile)
##[error]Process completed with exit code 1.
```
（`(...)` 处为省略重复的同形文字，**非省略判据**；原文整行见 `/tmp/a85-logs/job-106334066361.log` 第 515–533 行，
其中 515 = 步骤起、519 = endgroup、520–526 = `go: downloading`、527–532 = 6 行崩溃、533 = `##[error]…exit code 1`。）

### 1.3 "有没有任何一次真的给过 finding" —— 逐条回答

**没有。一次都没有。** 判据是**正向**的，不是"没看到 error"：

- finding 的形状是 `文件:行:列: 说明 (CHECKER码)`，码形如 `SA1019`/`U1000`/`ST1000`/`QF1000`/`S1011`。
  对 4 枚跨 run 日志按 `\(SA|ST|QF|U)[0-9]+` 正则以字段级扫描 ⇒ **命中 0 条**。
- 上述 6 行**不是 finding**，是工具自己的崩溃报告：前缀 `-:`（无文件无行号）、后缀 `(compile)`、
  内文 `internal error in importing` + `please report an issue`。
  ⇒ staticcheck 在**导入阶段**就死在标准库上，**一行仓库代码都没分析到**，所以它连"能报 finding 的机会"都没拿到。
- 本机纯净树复现同一崩溃：`rc=1`、5 行、同形（〔独立复现〕，见 §2.3）。
- ⚠ 一个诚实的细节：**报错行数在 5 与 6 之间抖动**（run 161/163/164 那批 5 行，run 167 6 行；本机 5 行）。
  `internal/goos` 时有时无 ⇒ 这是并行导入的**竞态**，不是"代码变了"。
  ⇒ 任何"固定行数"的判据都会飘；**要钉就钉 rc + 关键子串，别钉条数**。

### 1.4 门从未存在的**两种**机制（这是比"0 finding"更要紧的一条）
采样到的步骤状态只有两种形态，**没有第三种**（全历史 121 枚跑出来过 job 的 run，见 §7）：

1. **`skipped` × 93（更早的一段）**：步骤号 7，conclusion `skipped` —— 上面的 `gofmt`/`go vet` 红，
   Actions 一步失败即跳后续 ⇒ **这一步根本没跑**。
   run 1/2/…/18 等均是此形态（run 160 = job 106319703752 今天也回到了这个形态，因为 `gofmt (gofumpt)` 步骤 6 红）。
2. **`failure` × 28（上方的红被修掉之后）**：步骤号 9，conclusion `failure` —— 跑了、崩在导入、rc=1。
3. **`success` × 0** —— 这一步**历史上没有一次绿过**。
   附带一条更难看的事实：**`lint` job 本身 121/121 全是 `failure`**，即本仓 lint 门从未整体绿过。

⇒ **"这道门从未存在过"成立**，但根因是两条叠加，且**第 1 条票 85 单独修不掉**：
即使把 staticcheck 版本钉对，只要上方任何一步红，它仍然拿不到判据。
⚠ 注意：本仓库已三次被同一形状咬过（`ci.yml` 注释自述 A54/A61/票 71），
**这次是同一个病复发在 staticcheck 身上**，而且 `ci.yml:57-67` 的注释里已经写了这个机制。

---

## 2. 定因：`export data version 4 > maximum supported version 2` 是哪一侧

### 2.1 结论

**是"读的那一侧"太旧**：`staticcheck@2025.1.1` 里**vendor 进来的** `golang.org/x/tools v0.30.0` 的
pkgbits 解码器只认到版本 2；而 **Go 1.27 编译器写出的导出数据库是版本 4**。
不是"Go 工具链坏了"，也不是"我们的代码坏了"。**要动的是钉死在 `ci.yml` 里的那枚 staticcheck 版本**
（或同时降 go.mod 的 Go 线，但那是另一笔账）。

### 2.2 两侧的版本号各从哪一行读来（不引转述）

**消费侧（读不了的那侧）= x/tools v0.30.0，上限 2：**
- 载有它的二进制：`staticcheck.exe --version` → `staticcheck.exe 2025.1.1 (0.6.1)`
- 它带的 x/tools：`go version -m /tmp/a85-bin/staticcheck.exe` →
  `dep golang.org/x/tools v0.30.0 h1:BgcpHewrV5AUp2G9MebG4XPFI1E2W41zU1SaqVA9vJY=`
  （同一枚 v0.30.0 也逐字出现在 CI 日志的 `go: downloading golang.org/x/tools v0.30.0`，两侧对上）
- 上限的字面：`$GOMODCACHE/golang.org/x/tools@v0.30.0/internal/pkgbits/version.go` 常量块只有
  `V0 Version = iota` / `V1` / `V2` / `numVersions = iota` ⇒ numVersions=3，
  而报错行打的是 `numVersions-1` = **2**。
- 报错串出处（同一行号，两侧源码同形）：
  `golang.org/x/tools@v0.30.0/internal/pkgbits/decoder.go:86`
  ```go
  panic(fmt.Errorf("cannot decode %q, export data version %d is greater than maximum supported version %d", pkgPath, pr.version, numVersions-1))
  ```

**生产侧（写了 v4 的那侧）= Go 1.27：**
- `go version` → `go1.27.1 windows/amd64`
- `D:/work/base/go/src/internal/pkgbits/version.go:37-40`
  ```go
  // V4: encodes generic methods as standalone function objects
  V4

  numVersions = iota
  ```
  ⇒ Go 1.27 自己解码器上限是 4，并且它**写出** 4。
- CI 侧的 Go 版本来源：`ci.yml:37-39` `uses: actions/setup-go@v5` + `go-version-file: go.mod`；
  纯净树 `go.mod` 前 5 行 = `module github.com/CarlosShao/wisp` / `go 1.27` / `toolchain go1.27.1`。
  ⇒ **runner 上的 Go 与本机同线（1.27），所以这个 bug 本机可复现**（事实复现了，见 §2.3）。

### 2.3 本机独立复现（〔独立复现〕）

```
$ /tmp/a85-bin/staticcheck.exe ./...      # 在 /tmp/a85-swq85 纯净树上
rc=1  elapsed=18s  lines=5
-: internal error in importing "internal/byteorder" (cannot decode ... export data version 4 is greater than maximum supported version 2); please report an issue (compile)
-: internal error in importing "internal/cpu" (...)
-: internal error in importing "internal/goarch" (...)
-: internal error in importing "math/bits" (...)
-: internal error in importing "unicode/utf8" (...)
```

### 2.4 CI 装的是什么版本（引字面）

纯净树 `.github/workflows/ci.yml`：
```
87	      - name: staticcheck
88	        run: |
89	          go install honnef.co/go/tools/cmd/staticcheck@2025.1.1
90	          staticcheck ./...
```
⇒ 票面若写"第 89 行"→ 对得上；写别的行号即为腐坏。
⚠ 这一版 tag 会解析成 module `v0.6.1`（本机 `go list -m honnef.co/go/tools@2025.1.1` → `honnef.co/go/tools v0.6.1`），
与 CI 日志 `go: downloading honnef.co/go/tools v0.6.1` 一致。

### 2.5 本机环境的一条与票面预期不符之处（要点名）

任务书写"本机可跑：`staticcheck --version`"——**本机根本没有 staticcheck**：
`which staticcheck` → `no staticcheck in (...)`，`staticcheck: command not found`。
⇒ 为做 §3 的爆炸半径，我把两枚版本 `go install` 到 **`GOBIN=/tmp/a85-bin`**，
**没有**写进共享的 `D:\work\base\gopath\bin`（那是 3 个在跑的写码代理共用的 PATH），也没动仓库任何文件。

---

## 3. 爆炸半径：把 staticcheck 换成能真工作的版本，全仓到底报多少条

### 3.1 口径（先把话说死，否则数字没意义）

- 树：`/tmp/a85-swq85` = `git archive HEAD`（HEAD `3c5d1c3`）**纯净解出**，未在仓库目录内建 worktree/checkout。
- 工具：`staticcheck 2026.2.1 (0.8.1)`，vendor `golang.org/x/tools v0.44.1-0.20260420230617-19499e7caabc`
  （`go install honnef.co/go/tools/cmd/staticcheck@2026.2.1` → 解析为 module `v0.8.1`）。
- **一条都没修**，只数、只分类。
- ⚠ **`./...` 只覆盖 root module**：`tools/d22scan` 与 `tools/mockllm` 是**各自独立的 go.mod**，
  root 的 `./...` 结构性看不见它们 ⇒ 我另跑了一遍（见 3.4）。
- `frontend/` 里有 Go（`frontend/embed.go`），**计入**；`frontend/` 的 TS/Vue 不是 staticcheck 的活，不涉及。
- **含 `_test.go`**：staticcheck 默认连测试文件一起分析，我把含/不含拆开了（见下表）。
- ⚠ **本机 Windows ≠ CI ubuntu**：lint job 是 `runs-on: ubuntu-latest`，
  而本仓大量代码是 `*_windows.go` ⇒ 我**两平台都量了**，两个数都要看。

### 3.2 主表（root module，`staticcheck ./...`）

| 口径 | 条数 | U1000 | SA1019 | S1011 | SA9009 | SA4006 | SA4000 | SA4010 | ST1005 |
|---|---|---|---|---|---|---|---|---|---|
| **GOOS=windows（本机原生）** | **78** | 66 | 4 | 3 | 1 | 1 | 1 | 1 | 1 |
| **GOOS=linux（CI 平台形状）** | **34** | 26 | 3 | 2 | 1 | 1 | 1 | 0 | 0 |

- `_test.go` 占比：windows 24/78，linux 16/34。非测试代码 windows 54、linux 18。
- 包数：windows 21 个包；linux 17 个包。
- **耗时**：windows 全仓 **23s**、linux 交叉 **28s**、`go install` 编译两枚二进制 5s + 40s。
  ⚠ **本机此刻有 3 个写码代理在编译/跑测，非空载 ⇒ 以上用时仅供量级参考，不可用作 CI 时长预算。**
- **没有出现任何崩溃行**：78/78 与 34/34 都能被 `(CODE)` 形状解析出来。
  linux 唯一多出来的一行是
  `-: build constraints exclude all Go files in .../sherpa-onnx-go-linux@v1.13.8 (compile)`
  ⇒ 高度疑似**我交叉编译时 CGO_ENABLED=0 的仪器假象**，不是 CI 会看到的finding（CI 有 gcc）。
  这条我标〔独立复现，但**判为仪器产物**〕，**不要**把它算进爆炸半径，也不要替我删掉这个怀疑——
  真值要等第一枚真 CI run 才能定。

### 3.3 linux 侧按包分布（34 条，即"钉上版本那一刻 CI 会看到的红"）

```
5  internal/tools      4  internal/audio      4  internal/agent/approval
3  internal/agent      2  internal/observe    2  internal/memory
2  internal/llm/openaichat   2  internal/llm/adaptertest   2  internal/ball
1  each: internal/proc, internal/perm, internal/panel, internal/models,
        internal/llm/anthropic, internal/llm, internal/config, frontend
```
`internal/risk` **0 条** ⇒ 与 A56① 的"其中 internal/risk 0 条"**对上了**（子命题成立）。

### 3.4 两个独立 module（root `./...` 照不到的地方）

```
tools/d22scan: rc=1, 2 条   main.go:96:2 field failAddOn is unused (U1000)
                            main.go:103:2 var secretNameRe is unused (U1000)
tools/mockllm: rc=1, 1 条   chat.go:219:6 func estimateTokens is unused (U1000)
```
⇒ 若票 85 只想管现在这一步，爆炸半径是 **78（win）/ 34（linux）**；
若顺手把两把 module 也纳入，再加 **3** 条。

### 3.5 对旧账（A56①）的当场核对 —— 两处对不上，要点名

A56①（`docs/reports/pending-and-issues.md:2058` 附近）原文：
> "升到 2026.2.1 能跑，立刻 **35 条 finding** 散在 15 个包（`27×U1000` 未使用、`3×SA1019` 弃用 API、
> **`3×S1011` nil 解引用**、`1×SA4006`、`1×SA4000` 自己和自己比），其中 **`internal/risk` 0 条**"

核对结果：
1. **量级成立，逐字已腐坏**：我在 HEAD `3c5d1c3` 量到 linux **34 条 / 17 包**，他记 **35 条 / 15 包**；
   U1000 26 vs 27、S1011 2 vs 3，且我多出 1 条 SA9009 他没记。
   ⇒ 差异**可解释**：A56 写于今天 14:4x，之后 108/109/112/113 等票连续进树。
   **结论：票 85 里不要抄 35 这个数字，抄了就是在引用一小时内就会烂掉的数**；要引就引"命令 + HEAD sha"。
2. **`S1011` 被标错成语义了**：S1011 的实际报文逐字是
   `should replace loop with tokens = append(tokens, c21SortedNames(code)...)`，
   工具自己的清单一字写作 **`S1011 Use a single append to concatenate two slices`**
   （`staticcheck-2026.2.1.exe -list-checks`，共 149 条）。
   ⇒ **"nil 解引用"是错的**。⚠ 而且我顺手核了一圈：**2026.2.1 的 149 条清单里根本没有"nil 指针可能解引用"这一类检查**
   （`SA5011` 不在清单中；含 nil 字样的只有 `S1009/S1020/S1031/SA1012/SA4022/SA4023/SA4031/SA5000`，
   都不是"可能解引用空指针"）。
   ⚠ 反证不成立的一处：`staticcheck 2025.1.1` 不支持 `-list-checks`（输出 0 条），
   所以**它那一代有没有 SA5011 我测不出来**，只能断定 2026.2.1 没有。
   ⇒ 这条错标如果进了票 85 的 AC，会让人去 audit 一批**这道门根本不检查**的问题。**必须在开工前改掉。**

### 3.6 但这 34 条里**有真货**，别当噪音（决定要不要修的判断材料）

- ✅ **一条真实测试缺陷** —— `internal/observe/goroutine_test.go:106`：
  ```go
  if !errors.As(err, &oe) || !errors.As(err, &oe) {
  ```
  同一个调用写了两遍（SA4000 "identical expressions on the left and right side of the '||'"）。
  ⇒ 这条断言**有一半永远不起作用**。这道门一旦钉上版本，**当场就能抓到一条真的写错的断言**。
- ✅ 一条值得看的死代码 —— `internal/audio/device.go:126 func isDeviceInvalidated is unused`、
  `internal/agent/approval/approval.go:316 func (*grantStore).live is unused`
  ⇒ "写出来打算做安全判定、但没人调"的形状，正是该被门照出来的东西。
- ⚠ **一条是仪器的假阳性，且修它会撞到文档** —— `frontend/embed.go:4:1 SA9009`：
  真正的指令在第 19 行 `//go:embed all:dist`，**是好的**；第 4 行只是散文注释里
  "`go:embed` below is the ONLY way…" 恰好以 `// go:embed` 开头 ⇒ 被判"多余空格使指令无效"。
  ⇒ **不要**为了让这一步变绿而顺手改写这段注释：它是票 77 AC#1 的载荷文档。
  要么按 checker 的正当姿势处理，要么这条单独记为已知例外**由你裁定**——我不替你把门拆开。
- 26 条 U1000 里有 **16 条落在 `_test.go`**：测试里的 `unused` 常常是**故意的**（留给断言/反射/将来的钩子），
  删之前要逐条看，不能一把梭。

### 3.7 关键的一条排期事实（和直觉相反）

**这一步今天已经是红的。** §1 已证：`conclusion=failure`、`rc=1`。
⇒ "钉上版本"这个动作**不会把绿变成红**——lint job 现在就是红的，钉完还是红的。
它改变的是**这个红能不能读**：从 6 行"工具自己崩了"变成 34 行"你的代码有这些问题"。
⇒ 因此"先清完 34 条再钉版本"这个看似稳妥的顺序，**代价是把真判据往后推**，
而且清账过程中每次改动都没有门在看着。这条与 A56① 原本的"先升 pin 让 lint 红在真问题上再逐条清账"**一致**，
我实测后**支持原方案**，不支持反过来。

### 3.8 二次锚定复测（〔独立复现〕）

开工期间 HEAD 自己前进了（`3c5d1c3` → `a8f9459`），我按纪律**重新 `git archive HEAD` 解第二棵纯净树**
（`/tmp/a85-swq85b`）再量一遍，不拿旧数当新数：

| 工具链形状 | HEAD `3c5d1c3` | HEAD `a8f9459` | finding 集合 diff |
|---|---|---|---|
| GOOS=linux | 34 | **34** | **空（逐字节相同）** |
| GOOS=windows | 78 | **78** | 空 |

- 分类分布在两枚 HEAD 上完全一致（linux：`26 U1000 / 3 SA1019 / 2 S1011 / 1 SA9009 / 1 SA4006 / 1 SA4000`）。
- 耗时：linux 28s→30s、windows 23s→24s。⚠ **本机仍非空载**（3 个写码代理在跑），仍只作量级参考。
- 第二棵树同样**一条未修**。

---

## 4. 排期判断：**你给的第 4 步的前提没有成立**，所以先纠正再给拆法

### 4.1 前提纠正

你写："如果第 3 步量出来是**几百条**，那钉版本和清 finding 必须拆成两件事"。
⇒ **实测不是几百条**：CI 平台形状 **34 条**（本机 windows 形状 78 条）。
所以"必须拆"这条判据**没有被触发**。但——**我仍然建议拆，理由换成下面两条更硬的**：

1. **拆的理由不是量，是"混在一件事里验收判不出来"。** 钉版本改 `ci.yml`（票 111 地界），
   清 finding 要动 **21 个包**（`internal/ball` 一家就占 windows 侧 37 条）。
   同一张票里既有"工具链修复"又有"跨全仓代码改写"，对抗验收时**分不出哪条红是谁造的**。
2. **这个数字会漂 —— 但我这次的猜测被自己的复测否掉了，如实记下。**
   我先从"A56① 记 35 条 / 我今天记 34 条"推出"条数一小时内就会烂"。
   ⇒ 复测结果**不支持这个推广**：本次开工期间 HEAD 从 `3c5d1c3` 前进到 `a8f9459`
   （3 个写码代理真的又落了 commit），我**重新解了一棵纯净树再量一遍**：
   **linux 34 条 / windows 78 条，与前一枚 HEAD 逐字节相同（finding 集合 diff 为空）**。
   ⇒ 正确说法是：**这个数字会随进树的代码变，但在一小段窗口内是稳的**；
   A56① 的 35 → 34 差在 6 小时 + 108/109/112/113 多枚 commit 的尺度上，不在 40 分钟的尺度上。
   ⇒ 对票 85 的操作含义**不变，且更有利**：**判据可以写成"条数下降"这种可核对的形状**，
   不必因为怕漂就只能写 rc。但**仍然不许把 34/78 抄进票面当固定契约** ——
   开工当天请先重跑一次取当日值（命令见 §3.1）。
   〔独立复现：两枚 HEAD 各跑一遍，见 §3.8〕

⚠ 另一个只有拆才能容纳的事实：**CI 只看得到 34，但本机 windows 开发机永远看得到 78。**
差出来的 44 条在 `*_windows.go` / `*_windows_test.go` 里。
⇒ 只把 linux 侧的 34 条清干净，CI 绿了，**而开发机上 `staticcheck ./...` 仍然一片红**；
这个"绿/红不一致"要么明写进票面并接受，要么就得连 windows 形状一起清（那是另一笔工作量，
且**只有真 windows 机器能量**，与票 106/108 那条"dev 机不复现 runner"是同族反过来用）。**需要你拍。**

### 4.2 可执行的拆法（给选项 + 推荐，**不动手**）

**候选 A｜一票到底**：钉版本 + 清完 34（或 78）条，一次转绿。
- 代价：一次 commit 横跨 21 包 + `ci.yml`，与 3 个在飞的写码代理抢地界；验收无法归因。
- **不推荐。**

**候选 B｜两票串行（推荐）**
- **85a「钉死工具链版本」**：`staticcheck@2025.1.1` → `2026.2.1`、`gofumpt@latest` → 具体版本
  （`ci.yml:72` 现在就是 `@latest`，A56② 已点名）。**只改 `ci.yml`，不改一行业务码。**
  判据（形状，不写条数）：
  ① 该步输出里**不再出现** `internal error in importing` / `export data version`；
  ② 该步**真的产出了 finding 报文**（正向判据：至少一行匹配 `\(U1000\)$` 或别的 checker 码），
     或整步 rc=0；
  ③ 明写"**这一步预期是红的，红在真问题上；本次交件不许为了让它绿而删任何一条判据**"。
  ⚠ **需要 owner 批准的点**：85a 落地后 `lint` job **继续红**（从今天这种读不懂的红，变成 34 行读得懂的红）。
  也就是说"CI 转绿"这件事**不会因为 85a 而前进**，账面上还是同一格红。你若不能接受"多留几轮可读的红"，
  那只能选 C，而 C 我不建议。
- **85b「逐条清账」**：按 checker 分批，**每批都要做一次变异证明**（A56① 已要求："退回旧形状 ⇒ 必须重新报它"）。
  建议批次顺序 = 先做**能证明门真的在看这个包**的那几条，而不是最少的条数：
  1. `SA4000`（`observe/goroutine_test.go:106` 那条写重复的断言）+ `SA4006` + `SA4010` —— **3 条真逻辑味**，
     先修这三条，因为它们修完最能证明"这扇门有价值"；
  2. `SA1019 × 4`（全是 `runtime.GOROOT` 弃用）—— 同一处方，机械；
  3. `S1011 × 3` + `ST1005 × 1` —— 机械；
  4. `U1000 × 66`（windows 形状）—— **最贵且最容易误伤**，测试里的 unused 要逐条判"是不是故意留的"，
     不许一把梭删；
  5. `SA9009 × 1`（`frontend/embed.go:4`）—— **建议单独挂账、不在本票顺手改**：
     它要动的注释是票 77 AC#1 的载荷文档（§3.6）。**需要 owner 拍：改注释 / 记已知例外 / 挪给票 77 的地界。**

**候选 C｜钉版本 + 该步"先只报不挡"** —— **我不推荐，也不替你做。**
且必须点名：本仓 `ci.yml:5-8` 与票 111 明令 **D22 mode-6：no `continue-on-error`、no skip flag**。
"只报不挡"在这一家里**就是被禁的那个假修法**，写出来只为了说明我没顺手把它塞进推荐位。
⇒ 若你仍要"一段观察期"，**唯一与本仓纪律相容的形状是 `if: always()`**（加严：让后面的步骤也跑出判据），
它**不改这一步自己是红是绿**，因此**不构成拆门**。这条与票 85 已有的 AC#6 / R-4 是同一个动作，见下。

### 4.3 票 85 身上已经堆的非-staticcheck 账（拆票时必须一起判，别让它继续长）

| 条目 | 出处 | 我的核对结果 |
|---|---|---|
| AC#1 两把工具钉版本、禁 `@latest` | A56② (`pending-and-issues.md:2139`) | ✅ 成立，`ci.yml:72` 确为 `gofumpt@latest` |
| AC#2 "第一次产出真实判据" | A7x (`:1418` 区) | ⚠ **措辞要改**：不是"从未产出判据"，是"从未产出**可读的**判据"（它产出了 rc=1） |
| AC#4 改 `test-windows` 那个带 "placeholder" 的步骤名 | A56⑤ (`:2156`) | 与 staticcheck 无关，**可并到 85a**（都只动 `ci.yml`） |
| AC#6 `Environment fork assertion` 补 `if: always()` | A59④ (`:2050`) | ❌ **前提已失效，见 §4.4**：该步现在两枚 run 均为 `success`（票 93 重排治的）。⚠ 且它引的 `ci.yml:136-142` 已腐烂，HEAD 上是 **127-141** |
| R-4 `mockllm module vet` 被 staticcheck 连带 skip | A69② (`:1705`) | ✅ **本次拿到活数据**：run 167 步骤 9 failure ⇒ 步骤 10 `mockllm module vet` = `skipped` |
| R-99-4 `ci.yml:48` 与 `:68` 重复跑同一台仪器 | A7x (`:1328`) | ✅ 成立：`scripts/d22scan.sh:51` 内部**又**跑了一遍 `runtests.sh -C tools/d22scan ./...` |
| A73④ 抄来的作用域清单腐烂（`ci.yml:25`） | A73④ (`:1559`) | 只动 `ci.yml`，**可并到 85a** |

⇒ 我的排法建议：**"只动 `ci.yml` 的"合成一张（85a + AC#4 + R-4 + R-99-4 + A73④；**AC#6 拿掉**，理由见 §4.4）**，
**"动业务码的"另起一张（85b 清账）**。前者一次改完、一次验收、一次冲突；后者按 checker 分批。
**但这两张都仍排在 110/111 之后**（同文件），且 111 正在改 `ci.yml` 和 `scripts/`
⇒ **上面这些行号在 111 交件后还会再烂一次**，开工前必须重新解一遍 `git archive HEAD`，别拿本文行号当凭据。

### 4.4 ⚠ 当场核出来：**AC#6 的前提已经不成立了**（票 85 别再照原样做这件事）

A59④（`pending-and-issues.md:2050`）原话：
> `ci.yml:136-142` 的 `Environment fork assertion` 排在 `Portable package tests` **之后且没有 `if: always()`**
> ⇒ 两次 run 里它都是 `skipped` ⇒ 归票 85 一并落。

实测 HEAD `a8f9459`（取 `test-core` 的**步骤级** conclusion，两枚 run 一致）：

| run | 步骤 6 `Environment fork assertion` | 步骤 7 `Portable package tests` |
|---|---|---|
| 35600043583 (job 见 §1.1) | **`success`** | `failure` |
| 35599458439 | **`success`** | `failure` |

⇒ **它现在每一轮都真跑到、且绿了** —— 病因已被**票 93 的重排**治掉
（纯净树 `ci.yml:134-140` 的注释自己就写着 "MOVED ABOVE the portable step (ticket 93…)"；
现在 env 断言在 127-141 行、portable 在 143 行之后）。
⇒ **AC#6 按原措辞已经无可做**。三种处置，**请你挑一种**（我不替你结票）：
- (a) **勾掉 AC#6**，结案理由写"由票 93 重排达成，非本票所赐"，并附上面这枚 run id + step 6；
- (b) **改写成更强的形式**：重排只保证"`compose 启动`(步骤4) 与 `Probe mock-llm`(步骤5) 都绿"时才跑到；
  那两步任一红，env 断言仍会回到 `skipped`。⇒ 若要的是"任何情况下都出判据"，`if: always()` **仍是加严**，
  可以做，但要明写它补的是 4/5 这两步的洞，**不是** A59④ 说的那个洞；
- (c) 保持原样留给 85 但把行号与病因改对。
- 我推荐 **(a) + 把 (b) 单独立一条低优先级**：因为 (b) 的受益面比 A59④ 当时描述的小得多，
  混在 85 里会让这张票背上它其实没有的收益。
- ⚠ 顺带：A59④ 引的 **`ci.yml:136-142` 已腐烂**，HEAD 上该步是 **127-141**。
- ⚠ 与之对照，**R-4 那条仍然活着**：`mockllm module vet`（步骤 10）在 run 35600043583 里确实因
  staticcheck（步骤 9）红而 `skipped` ⇒ **R-4 保留在 85**，AC#6 不该和它同进退。

---

## 5. 下一张该派什么

0. **先派一张"票 85 票面修订"（只读代理即可，不必写码）**：
   要一次性改掉四处**我当场撞实**的坏引用，否则它们会变成开工代理的假前提（A43/A54 同族）：
   - `S1011 = nil 解引用` ⇒ 错（§3.5）；
   - `ci.yml:136-142` ⇒ 腐烂，现为 **127-141**（§4.4）；
   - **"staticcheck 从未产出判据"** ⇒ 错措辞，它产出了 `rc=1`，只是不可读（§0/§1.4）；
   - **"35 条"** ⇒ 当日值是 linux **34** / windows **78**（§3.2），且要写成"开工先重测"而不是抄数；
   - **AC#6 的前提** ⇒ 已被票 93 治掉，按 §4.4 的 (a)/(b)/(c) 由你裁定，不要让它混在 85 里蒙混结案。
   ⇒ **这一张不需要碰 `ci.yml`，也就不占 111 的地界，可以和 110/111 并行派。**这是我回来以后认为**性价比最高的一张**。
1. **然后派 85a**（只动 `ci.yml`，含 AC#4/R-4/R-99-4/A73④），**排在 111 交件之后**。
2. **85b 按 checker 分批**，第一批只要 3 条（SA4000/SA4006/SA4010）——**小而能证伪**，
   适合用来验证"这道门现在真的在看这个包"。
3. **⚠ 不要现在派**"在 ubuntu runner 上量真实条数"这件事单独成票：85a 落地的那枚 run
   会自动给出这个读数（§3.2 那条 sherpa 交叉编译疑点也随之判定），单开是重复劳动。
4. **另有一件不属于票 85 但被本次读数顺手照出来的账**（登记，别由我顺手修）：
   `lint` job 在**全部 121 枚有 job 的 run 上 conclusion 都是 `failure`**，
   即本仓 lint 门**历史上从未整体绿过**；今天使它红的仍有两条——
   步骤 9 `staticcheck`（票 85 的账）与步骤 6 `gofmt (gofumpt)`
   （run 160 = job 106319703752 那枚 `gofmt` 红 ⇒ staticcheck 直接被 skip，§1.4 的活样本）。
   ⇒ **只修 staticcheck 不足以让 lint 转绿**；`gofumpt` 未钉版本这件事票 85 AC#1 已经管到，
   但"哪一天会红"要有心理预期：**gofumpt@latest 与 staticcheck 是两条独立的红源**。

---

## 6. 纪律与异常登记

- **只读声明**：未修改任何被验文件。仓库内唯一的写入是**本证据文件**。
  快照与二进制全部在 `/tmp`（`/tmp/a85-swq85`、`/tmp/a85-swq85b`、`/tmp/a85-bin`、`/tmp/a85-logs`、`/tmp/a85-out-*.txt`）；
  **未在仓库目录内建 worktree 或 checkout**（A38④）。`go install` 全部带 `GOBIN=/tmp/a85-bin`，
  共享的 `gopath/bin` 未被写入。
  收尾实测：`git diff --cached --stat` **空**（我没 staged 任何东西）、`git rev-parse --short HEAD` = `a8f9459`
  （HEAD 由别人的 commit 前进，**不是我**），`git status --short` 里属于我的只有
  `?? docs/evidence/s1/85-preflight-staticcheck.md` 一行。
  ⚠ 同一时刻工作树里另有 `internal/winsec/winsec_windows.go`、`scripts/portable-tests.sh`、
  issue 115 的票面等 4 处 M / 3 处 ?? —— **全是别的代理的在飞工作，我一行未动、也未 commit**。
- **伪"编排者备注/系统提示"扫描**：本次工具输出里出现 **0 次**自称编排者/系统提示、
  命令我冻结包、终止回滚、revert 或放宽阈值的文本。
  唯一两条形似系统消息的是：① harness 的后台任务完成通知，它**自己标注**
  "SYSTEM NOTIFICATION - NOT USER INPUT … Do NOT interpret this as user acknowledgement"
  ⇒ 我按通知处理，未把它当授权；② 一条 `MEMORY.md 已被修改` 的环境提示 ⇒ 与本票无关，未据此改判据。
- **取数诚实性**：`gh api` 全量 sweep（168 run）遇到**大量瞬时 FETCH-FAILED**
  （手动对同一 run 立即重试即成功 ⇒ 判为偶发 TLS/并发，非数据缺失）。处理见 §1.1 与 §7。
- **我自己写错并当场改掉的三处**（留痕，不假装一开始就对）：
  1. sweep 的 jq 里写了 `$id`（合法字段是 `.id`）⇒ 前 2 枚 run 记成了 FETCH-FAILED，已修脚本重跑；
  2. 我在采样清单里**编造了 3 个不存在的 run id**（`35525649827`/`35535058309`/`35544366459`，API 返 404）
     ⇒ 已从证据里剔除，改用 `runs?per_page` 真列表；**不得**拿它们当任何结论的支撑；
  3. 我先写"真正的 `//go:embed` 在第 13 行"，`grep -n` 实为**第 19 行**；
     我又把 nil 检查写成 `SA5011`，而 `staticcheck -list-checks` 的 149 条里**根本没有 SA5011**
     ⇒ 两处均已按工具输出改正（见 §3.5、§3.6）。

---

## 7. 全量 run 历史 sweep —— **已完成，分母闭合**

数据文件 `/tmp/a85-resolved.tsv`；字段 `run_id  lint_job_id  lint_conclusion  staticcheck_step_conclusion  step_number`。

### 7.1 分母闭合的证据（这条决定"采不到"还是"采全了"）

- 全集：`gh api repos/CarlosShao/wisp/actions/runs?per_page=100` 两页 ⇒ **168 枚 run**（`ci` workflow，run 1 → 167）。
- 解出 `lint` job 的：**121 枚**。
- 未解出的 **47 枚**：**逐枚**回查 `runs/<id>` 结论 + `runs/<id>/jobs`，
  47/47 全为 `conclusion=cancelled`，且 jobs 端点原样返回
  ```json
  {"total_count":0,"jobs":[]}
  ```
  ⇒ 这 47 枚**根本没有创建过任何 job**，不是我没采到。
  **判据是正向的**（真解析出 `total_count` 字段并等值为 0），
  不是"响应里没有 error 字样就算没有"。⇒ **121 + 47 = 168，账平了。**

### 7.2 结论表（121 枚**跑出来过 job** 的 run，覆盖全部历史）

| staticcheck 步骤结论 | 枚数 | 含义 |
|---|---|---|
| `skipped` | **93** | 根本没跑到（上面某步先红，Actions 吃掉后续） |
| `failure` | **28** | 跑到就崩在 stdlib 导入，`rc=1` |
| `success` | **0** | **从未绿过** |
| `STEP-ABSENT` | 0 | 每一步都存在于 job 里，只是没跑 |

- 同一张表另一列更难看的事实：**`lint` job 自己的 conclusion，121/121 全是 `failure`**
  ⇒ 这个 repo 的 lint 门**历史上没有一次是绿的**。
- 步骤号随版本变过：早期是 **7**（D22 扫描排在 staticcheck **下面**），
  票 71 把它提到顶上后 staticcheck 变成 **9**。⇒ 我在 `skipped` 群里同时看到 7 和 9 两种形态，
  与 §1.4 讲的机制**互相印证**，不是我拼出来的故事。

### 7.3 "有没有任何一次给过 finding" —— 加大样本的复核

除 §1 的 4 枚外，我又从 28 枚 `failure` 里按序抽样 **7 枚**（共 11 枚不同 run 逐字取日志），
每枚都按两种形状计数：崩溃行 vs finding 行。

```
run=35560137482 job=106211346983  bytes=34520  crashlines=5  FINDINGLINES=0
run=35562680354 job=106218410858  bytes=34520  crashlines=5  FINDINGLINES=0
run=35566711794 job=106229891059  bytes=34596  crashlines=6  FINDINGLINES=0
run=35585821747 job=106288794874  bytes=37300  crashlines=5  FINDINGLINES=0
run=35587351912 job=106293644836  bytes=52422  crashlines=5  FINDINGLINES=0
run=35590599782 job=106304008761  bytes=52417  crashlines=5  FINDINGLINES=0
run=35599458439 job=106331840214  bytes=52489  crashlines=5  FINDINGLINES=0
```
finding 行判据：`grep -cE '\((SA|ST|S[0-9]|QF|U)[0-9]+\)[[:space:]]*$'`。
日志到手判据：`bytes > 100`（全部 34520~52489，**无一枚空日志被当成 0 finding 混进分母**）。
⇒ **11 枚逐字样本 + 121 枚步骤结论 + 本机复现，全部一致：staticcheck 从未产出过 1 条 finding。**

### 7.4 一句话答复你的第 1 问

**"staticcheck 从来没有产出过任何一条 finding" = 成立**（〔日志＋归档，我抽验〕，且核心形态另有〔独立复现〕）。
**但它后半句"它一直跑过去了" = 不成立**：它 93 次被跳过、28 次响亮地崩红，**0 次绿**。
所以这道门确实从未存在过，**机制是"从未跑到过 / 一跑到就崩"，不是"跑过了但假装没事"**。
