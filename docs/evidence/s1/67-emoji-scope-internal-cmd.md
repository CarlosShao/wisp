# 票 67 AC#3 判据②：ban #8 覆盖面真正扩到 `internal/` + `cmd/`，并证明它咬得住

生成：2026-09-21T01:29Z（`date -u`）· agent=agent-ticket67b · 分支 `dev` · 前置 HEAD `e5f4901`
registry：**A23 判据①**（"emoji 门对 Go 源码全盲"）

## 0. 结论先说

| 事项 | 结果 |
| --- | --- |
| ban #8 覆盖面 | `design/`（全部文本文件）**+ `internal/` + `cmd/` 的 `.go`，含 `_test.go`、含注释** |
| 覆盖面代码位置 | `tools/d22scan/main.go` 的 `emojiScopes()`（唯一真相源）+ `walkEmoji()` |
| `frontend/` 死作用域 | **已删除**（不存在的树 ⇒ 恒 0 文件）；并加了"声明了但走 0 文件 ⇒ exit 2"的致命检查 |
| 末行话术 | 由 `describeEmojiScopes()` **生成**，打印实际作用域 + 实际文件数；不再出现 `frontend/` |
| 第二批字形 | `internal/llm/probe_health.go:17/:120/:203` 注释已改 ASCII `PASS`/`FAIL`，**仅注释** |
| `allowlist.txt` | **5 行非注释，未变、未加豁免**（`git diff --name-only -- tools/d22scan/allowlist.txt` 为空） |
| 门现在能不能全绿 | **不能，且不止本地**：`internal/tools/bridge_junction_windows_test.go:444` 的 `⚠` 已在 `d63bc49`（票 20，01:27:37Z）**入库**，HEAD 里就有 ⇒ CI 检出 HEAD 也会红在这一行（§8.4 实测 exit 1）。我没碰该文件（硬约束：票 20 在写它），**没加豁免、没改断言、没 skip** |

## 1. 调用姿势（这个坑今天骗过编排者一次，先钉在证据里）

`tools/d22scan` 是**独立 Go module**。从仓根 `go run ./tools/d22scan -root .` 只会打印一条模块错误，
扫描器从未启动 —— 那是假绿，不是 clean。本次所有输出**全部**来自下面这条正确调用的真实 stdout：

```
cd tools/d22scan && go run . -root ../..
```

⚠ `go run` 会把子进程任意非零退出压成 `1`（票 67 前任已登记）。为区分 exit 1（有违规）与
exit 2（仪器坏了），种子阳性那一步**另跑了编译出的二进制**并把两个退出码都贴出来。

## 2. 扩面前：门对 Go 源码全盲（`git show HEAD:tools/d22scan/main.go` 实测）

```
156:	if err := s.walkEmoji(filepath.Join(root, "design")); err != nil {
159:	if err := s.walkEmoji(filepath.Join(root, "frontend")); err != nil {
427:func (s *scanner) walkEmoji(dir string) error {
544:	fmt.Println("d22scan: clean - no D22 ban violations, no emoji in design/ or frontend/")
90:	emojiRe = regexp.MustCompile(`[\x{1F000}-\x{1FAFF}\x{2600}-\x{27BF}\x{2B00}-\x{2BFF}\x{FE0F}\x{1F1E6}-\x{1F1FF}]`)
449:			if emojiRe.MatchString(line) {
```

- `walkEmoji` 只有 `design/` 与 **不存在的** `frontend/` 两个调用点 ⇒ 产品 310 个 `.go`（`find internal cmd -name '*.go' | wc -l`）
  对 ban #8 全部不可见，而末行仍打印 "no emoji in design/ or frontend/"。
- `:449` 逐行匹配**原文**、不剥注释 ⇒ 注释里的字形在门语义里本来就是违规；
  票 67 AC#3 原文"注释不算用户可见、扩展时排除注释"因此**不是**少改几行，而是要**新增**剥注释的代码 = 缩小覆盖面，
  按 R16#4 不许做。已按此实现，理由写在 `walkEmoji` 的函数注释里。

## 3. 扩面后、清字形前：第一次真跑就是红的（4 处 = 预测的爆炸半径）

```
$ cd tools/d22scan && go run . -root ../..
d22scan: examined 197 production Go files under internal/ and cmd/ of D:/work/workspace/projects plans/Wisp
d22scan: ban #8 scope design/    examined 16 file(s)
d22scan: ban #8 scope internal/  examined 289 file(s)
d22scan: ban #8 scope cmd/       examined 21 file(s)
internal/llm/probe_health.go:17: [emoji] ban #8 glyph in scope internal/ is banned (D23): covers comments and _test.go, not only string literals
internal/llm/probe_health.go:120: [emoji] ban #8 glyph in scope internal/ is banned (D23): covers comments and _test.go, not only string literals
internal/llm/probe_health.go:203: [emoji] ban #8 glyph in scope internal/ is banned (D23): covers comments and _test.go, not only string literals
internal/tools/bridge_junction_windows_test.go:444: [emoji] ban #8 glyph in scope internal/ is banned (D23): covers comments and _test.go, not only string literals
d22scan: 4 finding(s); D22 bans are not negotiable (see PLAN.md D22, tools/d22scan/allowlist.txt)
exit status 1          # go run；二进制直跑同输出，EXIT=1
```

- `289 + 21 = 310`，与 `find internal cmd -name '*.go'` 逐数吻合 ⇒ 作用域自报的数字不是装饰。
- 3 行生产命中 + 1 行测试命中，与编排者 09:17 的实测**逐行一致**，没有挖出第 5 处。

## 4. 第二批字形（仅注释，零逻辑改动）

`internal/llm/probe_health.go` 三处 `「声明 ✓ 实测 ✗」` → `「声明 PASS 实测 FAIL」`，
落点措辞与同族文件 `cmd/wisp/providers.go:7/:41`（`fff4cad`）一致。
SPEC-05 §3.1 / SPEC-03 原文仍用对勾与叉（`docs/specs/**` 冻结、且不在 ban #8 作用域内），
故 `:17` 注释里补了一句"此处是 ASCII、spec 保留原字形"，防止后人当成误引去"改正"。

`git diff -- internal/llm/probe_health.go` 的三个 hunk 全部以 `//` 开头；常量名、函数签名、逻辑零改动
（`gofmt -l internal/llm` 空、`go vet ./internal/llm/` exit 0）。清理后该文件 ban 区间字形计数：

```
$ grep -cP '[\x{1F000}-\x{1FAFF}\x{2600}-\x{27BF}\x{2B00}-\x{2BFF}\x{FE0F}\x{1F1E6}-\x{1F1FF}]' internal/llm/probe_health.go
0
```

## 5. 种子阳性（证明门对新覆盖面真的会红，不是"看起来会变红"）

种子 = 一行注释里的 U+2713，插进**已跟踪**文件 `internal/llm/probe_health.go`（本票唯一允许我写的生产文件），
这样撤销才能用 `git diff` 证伪——用新建临时文件的话 `git diff` 永远为空，那是假证。

### 5.1 先证明种子真的落盘

```
$ grep -nP '\x{2713}' internal/llm/probe_health.go
3:// zz_ticket67b_seed: ban #8 positive control - a U+2713 here must turn the gate red ✓
$ grep -cP '\x{2713}' internal/llm/probe_health.go
1
$ git diff --stat -- internal/llm/probe_health.go
 internal/llm/probe_health.go | 2 ++
 1 file changed, 2 insertions(+)
```

### 5.2 扫描器变红（两种调用，退出码都贴）

```
$ cd tools/d22scan && go run . -root ../.. ; echo EXIT=$?
d22scan: examined 197 production Go files under internal/ and cmd/ of D:/work/workspace/projects plans/Wisp
d22scan: ban #8 scope design/    examined 16 file(s)
d22scan: ban #8 scope internal/  examined 289 file(s)
d22scan: ban #8 scope cmd/       examined 21 file(s)
internal/llm/probe_health.go:3: [emoji] ban #8 glyph in scope internal/ is banned (D23): covers comments and _test.go, not only string literals
internal/tools/bridge_junction_windows_test.go:444: [emoji] ban #8 glyph in scope internal/ is banned (D23): covers comments and _test.go, not only string literals
d22scan: 2 finding(s); D22 bans are not negotiable (see PLAN.md D22, tools/d22scan/allowlist.txt)
exit status 1
EXIT(go run)=1

$ /tmp/d22scan67b.exe -root ../.. ; echo EXIT=$?      # 同一棵树，编译出的二进制（go run 会把退出码压成 1）
d22scan: examined 197 production Go files under internal/ and cmd/ of D:/work/workspace/projects plans/Wisp
d22scan: ban #8 scope design/    examined 16 file(s)
d22scan: ban #8 scope internal/  examined 289 file(s)
d22scan: ban #8 scope cmd/       examined 21 file(s)
internal/llm/probe_health.go:3: [emoji] ban #8 glyph in scope internal/ is banned (D23): covers comments and _test.go, not only string literals
internal/tools/bridge_junction_windows_test.go:444: [emoji] ban #8 glyph in scope internal/ is banned (D23): covers comments and _test.go, not only string literals
d22scan: 2 finding(s); D22 bans are not negotiable (see PLAN.md D22, tools/d22scan/allowlist.txt)
EXIT(binary)=1
```

种子在 `internal/` 的**注释行**被抓 ⇒ §2 里"注释算不算"这条不是解释出来的，是跑出来的。

### 5.3 撤销并证明撤销干净

```
$ git diff -- internal/llm/probe_health.go | wc -l
0
$ git status --short -- internal/llm/          # 空输出
$ grep -cP '\x{2713}' internal/llm/probe_health.go
0
$ grep -rn 'zz_ticket67b_seed' --include='*.go' internal cmd tools | wc -l
0
$ grep -rnP '[\x{1F000}-\x{1FAFF}\x{2600}-\x{27BF}\x{2B00}-\x{2BFF}\x{FE0F}\x{1F1E6}-\x{1F1FF}]' --include='*.go' internal cmd
internal/tools/bridge_junction_windows_test.go:444:	// ⚠ 本用例发现的真实软处（原样钉住，不顺手改）：junction 挡在 A 档目标前面时，
```

撤销后 `internal/` + `cmd/` 的 `.go` 里 ban 区间字形**只剩票 20 那一行**（§7）。

## 6. d22scan 自身测试：10 条，9 绿 1 红，红因点名

```
$ cd tools/d22scan && go test -v ./...      # REAL_EXIT=1
--- PASS: TestScanDetectsAllSeededViolations
--- PASS: TestScanCleanRepoIsGreen
--- PASS: TestAllowlistSuppressesOnlyListedPaths
--- PASS: TestCheckRootRejectsBlindRoots
--- PASS: TestScanAloneIsNotAFalsifier
--- PASS: TestCheckRootAcceptsRealRepo
--- FAIL: TestScannerSelfScanOfRealRepoIsGreen
    scan_test.go:264: repo HEAD violates: internal/tools/bridge_junction_windows_test.go:444: [emoji] ...
--- PASS: TestEmojiBanCoversGoSourcesNotJustDesign        # 本次新增
--- PASS: TestDeclaredEmojiScopeCannotWalkZeroFiles       # 本次新增
--- PASS: TestScopeReportMatchesRealCoverage              # 本次新增
```

唯一红因 = `internal/tools/bridge_junction_windows_test.go:444` 的那个 `⚠`（U+26A0）。
**未改断言、未加豁免、未 skip**（票面硬规矩 + R16#4）。

## 7. 已知待清 1 处 —— 点名移交（同文件并发＝假并行，这条优先于"CI 必须绿"）

- 文件：`internal/tools/bridge_junction_windows_test.go:444`（注释里的 `⚠`，U+26A0）
- 我**一个字没碰**，也没加 allowlist 豁免、没改断言、没 skip。

> **追加更正（本文件 commit A 的判断已被随后发生的事实推翻，保留原文以便审计）**
>
> commit A（`bcf44d6`，2026-09-21T01:32Z）前后我两次 `git status` 看到它是 `??` 未跟踪，据此写下
> "CI 检出 HEAD 时看不到它 ⇒ CI lint job 仍绿"。**错了**：票 20 的代理在 `d63bc49`
> （2026-09-21T09:27:37+08:00 = 01:27:37Z，比我 commit 早约 5 分钟）把这个文件**连字形一起入库**了。
> 于是 HEAD 里就带着这 1 处违规 ⇒ **CI 的 lint job 在 HEAD 上是红的**，不是绿的。
> 实测见 §8.4/§8.5（`git archive HEAD` 导出的纯净树扫描 = 1 finding / exit 1；`sh scripts/d22scan.sh` = exit 1）。
> 我没有为了让它绿去改那个文件（硬约束优先于"CI 必须绿"），也没有用豁免绕（`allowlist.txt` 仍 5 行、`git diff` 为空）。

## 8. 仪器与门状态实测（commit B）

### 8.1 死作用域不再可能"假装在扫"：新增 `emptyEmojiScope()` 真的 exit 2

构造一棵 `cmd/` 下**没有任何 `.go`** 的假仓（11 个 `internal/*.go` 过 `checkRoot` 阈值，`cmd/onlydocs/readme.md` 撑出目录）：

```
$ /tmp/d22scan67b.exe -root /tmp/emptyscope67b ; echo EXIT=$?
d22scan: examined 11 production Go files under internal/ and cmd/ of C:/Users/swq/AppData/Local/Temp/emptyscope67b
d22scan: ban #8 scope design/    examined 1 file(s)
d22scan: ban #8 scope internal/  examined 11 file(s)
d22scan: ban #8 scope cmd/       examined 0 file(s)
d22scan: ban #8 scope cmd/ examined 0 files - it is declared in emojiScopes() but walks nothing. Point it at a real tree or delete the entry; never leave a scope pretending to scan (ticket 71 AC#4)
EXIT=2
```

⇒ exit 2（仪器坏了）与 exit 1（有违规）可分辨；这就是"删掉 `frontend/`"之后仍然更硬的那一层：
将来谁把条目加回一棵不存在的树，门自己会喊，不需要人复核话术。

### 8.2 撤销后工作树状态

```
$ cd tools/d22scan && go run . -root ../.. ; echo EXIT=$?
（3 行作用域自报 + ）internal/tools/bridge_junction_windows_test.go:444: [emoji] ...
d22scan: 1 finding(s); D22 bans are not negotiable (...)
EXIT=1
```

我这一支（票 67b）带来的文件 **0 命中**：`probe_health.go` 三行已清（§4），种子已撤（§5.3）。

### 8.3 `git diff --cached` 纪律

本票两个 commit 前都跑过 `git diff --cached --name-only`，暂存清单恰为
`tools/d22scan/{main.go,scan_test.go}`、`internal/llm/probe_health.go`、
`docs/evidence/s1/67-emoji-scope-internal-cmd.md`、`.scratch/wisp/issues/67-d22scan-trust-mockllm-goroutine.md`；
同一时刻工作树里**别人在飞**的 `go.mod`、`internal/observe/*`、`internal/secret/*` 全程留在未暂存区，没被吞进我的 commit。

### 8.4 CI 等价验证：`git archive HEAD` 纯净树（不看任何未跟踪文件）

```
$ git archive HEAD | tar -x -C /tmp/head67b
$ /tmp/d22scan67b.exe -root /tmp/head67b ; echo EXIT=$?
d22scan: examined 197 production Go files under internal/ and cmd/ of C:/Users/swq/AppData/Local/Temp/head67b
d22scan: ban #8 scope design/    examined 16 file(s)
d22scan: ban #8 scope internal/  examined 289 file(s)
d22scan: ban #8 scope cmd/       examined 21 file(s)
internal/tools/bridge_junction_windows_test.go:444: [emoji] ban #8 glyph in scope internal/ is banned (D23): covers comments and _test.go, not only string literals
d22scan: 1 finding(s); D22 bans are not negotiable (see PLAN.md D22, tools/d22scan/allowlist.txt)
EXIT=1
```

同一棵导出树里 `grep -cP '[\x{2600}-\x{27BF}]' internal/llm/probe_health.go` = **0**
⇒ HEAD 的生产 `.go` 已零字形；唯一红因是票 20 入库的那 1 行注释。

### 8.5 CI lint 步逐字复跑（`sh scripts/d22scan.sh`）

```
$ sh scripts/d22scan.sh ; echo SCRIPT_EXIT=$?
d22scan.sh: positive control - go test ./... (tools/d22scan)
--- FAIL: TestScannerSelfScanOfRealRepoIsGreen (0.32s)
    scan_test.go:264: repo HEAD violates: internal/tools/bridge_junction_windows_test.go:444: [emoji] ban #8 glyph in scope internal/ is banned (D23): covers comments and _test.go, not only string literals
FAIL	github.com/CarlosShao/wisp/tools/d22scan	0.737s
SCRIPT_EXIT=1
```

**⇒ 本票交付后 CI 的 lint job 仍是红的，红因不在我的改动里。** 清掉它只需要一处 ASCII 替换
（`⚠` → `[!]` 或 `WARN:`），但那是对票 20 正在写的文件下笔，归编排者在票 20 落地后处理。

## 9. 移交（我没权限或没边界动、但已被这次改动说清的三件事）

1. `.github/workflows/ci.yml:23` 的注释仍写 "zero-emoji scan over design/ and frontend/" ⇒ 现在与实际覆盖面不符。
   该文件是票 70 接续代理的领地，我**没改**，请顺手的代理改成 "design/ + the Go sources of internal/ and cmd/"。
2. ban #6 `panel-approval` 仍 `walkText` 一棵不存在的 `frontend/` ⇒ 恒 0 作用域。**不是我该缩的覆盖面**，
   归票 71 AC#4（与本票删掉的那条同形）。
3. 票 12 AC#7 的"另一半"（emoji 门看不见 Go 源码）自本票起前提已成立，勾框措辞归编排者。
