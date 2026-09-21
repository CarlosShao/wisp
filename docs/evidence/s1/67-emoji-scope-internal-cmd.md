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
| 门现在能不能全绿 | **本地不能**：`internal/tools/bridge_junction_windows_test.go:444` 还剩 1 个 `⚠`，该文件此刻由票 20 的代理在写，我没碰（见 §6） |

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

见本文件末尾 §8（commit B 追加：先 grep 证明种子落盘 → 扫描器红输出 + 两个 exit code → 撤销 → `git diff` 空 + grep 0）。

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

- 文件：`internal/tools/bridge_junction_windows_test.go:444`（注释里的 `⚠`）
- 状态：**票 20 的代理此刻正在写这个文件**，我在 `git status` 里看到它是 `??`（未跟踪），本票按硬约束**一个字没碰**。
- 影响面：因为它未跟踪，**CI 检出 HEAD 时看不到它** ⇒ CI 的 lint job 在我这两个 commit 上应当仍是绿的（§8 用 `git archive HEAD` 独立验证，不看未跟踪文件）。
  本地工作树跑 `sh scripts/d22scan.sh` 会红在这一行。
- 请求：票 20 落地后，由编排者（或该代理自己被门拦到时）把这 1 个 `⚠` 改成 ASCII（推荐 `[!]` 或 `WARN:`），
  然后 `cd tools/d22scan && go run . -root ../..` 应当 exit 0。**不要**通过给 `internal/tools/` 加 allowlist 来"清"它。

## 8. （commit B 追加）种子阳性的完整证据链 + HEAD 树等价验证
