# A30 证据引用修复：复现输出汇总（2026-09-20）

授权：编排者（A30②④「②④ 归证据更正代理——只追加、不覆写原文，每条带可复现命令」）。
本文只是**四份已归档验收报告更正小节的复现输出留档**，裁决以那四个小节为准；本文不新增裁决。

被更正文件（均为**文末追加** `## 更正（A30，2026-09-20，编排者授权）`，原文未改一字）：

- `docs/evidence/s1/08-adversarial-acceptance.md` 第 10 行「d22scan 实跑全仓 clean」
- `docs/evidence/s1/18-adversarial-acceptance.md` 第 11 行「tools/d22scan 全仓 clean（allowlist 1 行豁免…）」
- `docs/evidence/s1/17-adversarial-acceptance.md` 第 12 行「越界/D22 | PASS」（三无 PASS：无命令/无 file:line/无 SHA）
- `docs/evidence/s1/09-adversarial-acceptance.md` 第 42 行 + 第 254 行「`TestGoldenCancellationMidStream`（adapter_test.go:389-427）」

## 0. 测量条件（四条更正共用）

- 时间：2026-09-20T15:54Z（本地 23:54 前后）。
- 树：`dev` 分支，工作树含**另一代理在途的全仓 gofumpt 未提交改动**（`git status --porcelain | grep -c '\.go$'` = 67）。
- 本修复**未执行任何格式化命令**、未 `git add` 任何 `.go`、未 `git checkout`/`restore`/`stash`。
- 结论侧：**没有出现任何意外 findings / 意外失败**。

## 1. d22scan：错命令 vs 对命令（票 08 / 17 / 18 共用）

### 1a. 当时那条错命令（仓库根）——证明它确实不扫任何东西

```
$ cd "D:\work\workspace\projects plans\Wisp" && go run ./tools/d22scan -root .
main module (github.com/CarlosShao/wisp) does not contain package github.com/CarlosShao/wisp/tools/d22scan
EXIT=1
```

事实：`tools/d22scan` 自 `64d083d` 起就是独立 module（`git show 64d083d:tools/d22scan/go.mod` →
`module github.com/CarlosShao/wisp/tools/d22scan`）。该命令**不打印任何 findings、也不打印 examined**，
exit code 为 1。用它得到的"clean"没有任何含义；失败信号存在但被读成了干净。

### 1b. 正确命令（模块内）——自报工作量 + exit code

```
$ cd "D:\work\workspace\projects plans\Wisp\tools\d22scan" && go run . -root ../..
d22scan: examined 194 production Go files under internal/ and cmd/ of D:/work/workspace/projects plans/Wisp
d22scan: clean - no D22 ban violations, no emoji in design/ or frontend/
EXIT=0
```

`examined 194` 与票 67 记录一致。`examined` 行是 1a 与 1b 的唯一区分凭据。

### 1c. "clean" 不是空跑的阳性对照（扫描器自测）

```
$ cd tools/d22scan && go test . -count=1 -v
--- PASS: TestScanDetectsAllSeededViolations (0.02s)
--- PASS: TestScanCleanRepoIsGreen (0.01s)
--- PASS: TestAllowlistSuppressesOnlyListedPaths (0.01s)
--- PASS: TestCheckRootRejectsBlindRoots (0.02s)
    --- PASS: TestCheckRootRejectsBlindRoots/empty_dir
    --- PASS: TestCheckRootRejectsBlindRoots/go.mod_only
    --- PASS: TestCheckRootRejectsBlindRoots/the_d22scan_module_itself
    --- PASS: TestCheckRootRejectsBlindRoots/repo_skeleton_with_two_go_files
--- PASS: TestScanAloneIsNotAFalsifier (0.00s)
--- PASS: TestCheckRootAcceptsRealRepo (0.00s)   scan_test.go:234: real repo production Go files in scope: 194
--- PASS: TestScannerSelfScanOfRealRepoIsGreen (0.12s)
PASS
ok  	github.com/CarlosShao/wisp/tools/d22scan	0.207s
EXIT=0
```

## 2. allowlist 计数（票 18 行内嵌的第二处失实）

```
$ git show 3df0218:tools/d22scan/allowlist.txt | grep -vc '^#\|^$'
4
```

`3df0218` = 票 18 实现提交（2026-09-20 08:23:14 +0800；报告戳记 00:27:30Z = 08:27:30）。
四条豁免（今天工作树同样是这四条）：`internal/memory/open.go`（票 06）、
`internal/models/manifest.go`（票 14）、`internal/models/downloader.go`（票 14）、
`internal/risk/pathresolver.go`（票 18）。⇒「全仓 1 行豁免」不成立；属票 18 的只有 1 行。

## 3. 票 17 第 5 行的另两个子主张

```
$ git show --stat 67ffbd8        # 票面 Progress log:46 给的 SHA
9 files changed, 1222 insertions(+), 1 deletion(-)
 .scratch/wisp/issues/17-risk-assessor-c19.md, internal/risk/{assessor,assessor_test,rules_gateway,
 rules_irreversible,rules_network,rules_scale,rules_shell,rules_test}.go        EXIT=0
```
⇒ 越界子主张成立（全在 `internal/risk/`，T18 文件零触碰）。

```
$ git show 67ffbd8 | grep -cP "[\x{1F300}-\x{1FAFF}\x{2600}-\x{27BF}\x{FE0F}]"
0        （无命中 ⇒ grep EXIT=1）
$ git show 5aa3258 | grep -cP "[\x{1F300}-\x{1FAFF}\x{2600}-\x{27BF}\x{FE0F}]"
9        （同管线阳性对照 ⇒ EXIT=0）
```
⇒ 零 emoji 子主张成立。**注意工具口径**：d22scan 的 ban #8 只扫 `design/` 与 `frontend/`
（`tools/d22scan/main.go:21、150-153`），不覆盖 `internal/`，故票 17 的"零 emoji"从不由 d22scan 背书。

## 4. 票 09 的失效引用 → 现名

```
$ grep -rn "TestGoldenCancellationMidStream" --include=*.go . ; echo EXIT=$?      EXIT=1
$ git grep -n "TestGoldenCancellationMidStream" HEAD -- '*.go' ; echo EXIT=$?     EXIT=1
$ git log --oneline -S "TestGoldenCancellationMidStream" -- internal/llm/openaichat/
5ddedf7 feat(llm): add anthropic + openai-responses adapters on one shared golden harness (ticket 11 AC#1, AC#3)
3b7bf24 feat(llm): golden-driven adapter tests, mockllm integration (... )
$ wc -l internal/llm/openaichat/adapter_test.go
291        （⇒ 被引的 389-427 现在落在文件之外；写下时 3b7bf24 精确命中 389/427）
```

现名 **`TestSharedGoldenSuite/cancel`**（断言体 `internal/llm/adaptertest/harness.go:579-623` `runCancel`，
入口 `internal/llm/openaichat/harness_golden_test.go:61`）：

```
$ go test ./internal/llm/openaichat/ -run 'TestSharedGoldenSuite/cancel' -count=1 -v
=== RUN   TestSharedGoldenSuite
=== RUN   TestSharedGoldenSuite/cancel
--- PASS: TestSharedGoldenSuite (1.38s)
    --- PASS: TestSharedGoldenSuite/cancel (1.38s)
PASS
ok  	github.com/CarlosShao/wisp/internal/llm/openaichat	1.414s
EXIT=0

$ go test ./internal/llm/anthropic/ ./internal/llm/openairesponses/ ./internal/llm/openaichat/ -run 'TestSharedGoldenSuite/cancel' -count=1
ok  	github.com/CarlosShao/wisp/internal/llm/anthropic        1.419s
ok  	github.com/CarlosShao/wisp/internal/llm/openairesponses  1.422s
ok  	github.com/CarlosShao/wisp/internal/llm/openaichat       1.423s
EXIT=0
```

`-run` 的阳性对照（不匹配的 `-run` 也打 `ok`，所以"ok"单独不构成凭据）：

```
$ go test ./internal/llm/openaichat/ -run 'TestSharedGoldenSuite/no-such-subtest-xyz' -count=1 -v
=== RUN   TestSharedGoldenSuite
--- PASS: TestSharedGoldenSuite (0.00s)
testing: warning: no tests to run
PASS
ok  	github.com/CarlosShao/wisp/internal/llm/openaichat	0.034s [no tests to run]
EXIT=0
```

四条原断言逐条对上新断言体（`harness.go:601` err==nil、`:604` 末事件 done、`:616` 含 `Stop{cancelled}`、
`:622` 事件数 <32），MINOR-3 的零 `EvError` 在 `:611-612`/`:619`；fixture 同为 paced `long-text`
（`harness_golden_test.go:53-54`）。⇒ **等价物存在，且现对三个协议腿各跑一遍（覆盖面变大）**。

## 5. 四条更正的落点

| 报告 | 被更正句 | 重判 |
|---|---|---|
| 08 第 10 行 | d22scan 实跑全仓 clean | **作废/无凭据**；今日 clean 不追认当时；该行其余子句归 A27/A23，本节不裁 |
| 18 第 11 行 | tools/d22scan 全仓 clean | **作废/无凭据**；另"allowlist 1 行豁免"是**写下当时即失实**的计数（实为 4） |
| 18 第 11 行后半 | 真 junction 测试 | **不受影响**（第 2 行红队四连有独立凭据） |
| 17 第 12 行 | 越界 / D22 / 零 emoji | 越界 **成立**（补 `git show --stat 67ffbd8`）；D22 **无凭据**（不追认）；零 emoji **成立**（补真命令 + 阳性对照） |
| 09 第 42 / 254 行 | `TestGoldenCancellationMidStream`（`adapter_test.go:389-427`） | **引用失效 ⇒ 现名 `TestSharedGoldenSuite/cancel`**；断言实质成立、未放宽 |

**没有任何一条被我用"今天的干净"追认成"当时的干净"**；票 08/17/18 的其余 AC 各有独立证据、不受影响。
