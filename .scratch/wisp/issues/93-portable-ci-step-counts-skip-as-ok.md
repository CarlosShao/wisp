# 93 — portable 测试步是裸 `go test` ⇒ **SKIP 记成 ok**：`TestSyncRegistryProbeLive` 两侧都跳，CI 却全绿

**Status:** open（2026-09-21 17:1x 编排者建，来源=票 82 交件时**自己点名**的残留；我没让它顺手修，因为那超出票 82 的界）
**Status (agent-ticket93, 2026-09-21 18:5x):** ready-for-review —— AC#1/#2/#3/#5 见文末"交件段"；
**AC#4 未勾**（本代理不 push，`c8c828e` 还没进任何 run ⇒ 欠编排者 push 后复跑，取数命令已写在文末）。
`allowlist.txt` 与 `tools/d22scan/**` 一字未动；没删步骤、没改阈值、没放宽任何断言、没用 build tag 排除整包。
**Type:** 门禁完整性（票 71 家族："门没跑"和"跑了没问题"在 CI 输出上长得一模一样）
**Blocks:** nothing（但它在**削弱** `test-core` 全绿这句话的含义）· **Blocked by:** nothing
**Packages:** `.github/workflows/ci.yml`（`test-core` 的 "Portable package tests" 步）、
              `scripts/` 下一个可复用 runner（**先读 `tools/d22scan/runtests.sh`，它已经会拒 SKIP**）。
              **禁改**：任何测试的**断言**、任何阈值、`docs/PLAN.md`、`docs/specs/*.md`、
              `tools/d22scan/**` 与 `allowlist.txt`（那是编排者的门禁面）。

## 现场（可直接复现）

`internal/risk` 的 `TestSyncRegistryProbeLive` 在 **Windows 与 ubuntu 两侧都 `--- SKIP`**
（Windows：这台机器 HKCU 无 `UserFolder`；POSIX：压根没有注册表）。CI 的 portable 步是
`go test <16 个包> -count=1`（`ci.yml` 的 "Portable package tests"），**裸 go test 对 SKIP 的记账是 `ok`**
⇒ 这条用例在 CI 上**从来没有产出过一个结论**，而输出看起来和"它通过了"完全相同。
票 82 的代理在 `-count=2` 下量到 **2 行 SKIP，同名 1 个**。

⚠ 同一仓里已经有**两套仪器对同一件事一严一松**：`tools/d22scan/runtests.sh` **拒 SKIP**，
portable 步不拒。这就是本票的缺陷本体——**不是那条用例该不该跳，是"跳了没人记账"**。

## AC（1:1，裁决表 `docs/evidence/s1/93-*.md` 由验收方出）

- [x] **AC#1** 先做**全量点名**：把 CI portable 清单里 16 个包在 **Windows 与 ubuntu 两侧**当前所有
      `--- SKIP` 逐条列出（包 / 测试名 / file:line / 跳过原因）。**判据是"清单是完整的"**：
      你的取数命令必须能证明它抓到了全部输出（贴 `grep -c -- "--- SKIP"` 与总行数的对照），
      不许只贴 `internal/risk` 那一条。
- [x] **AC#2** 落一个**真会红**的机制：portable 步遇到任何未被显式记账的 SKIP ⇒ 该步 rc≠0。
      ⚠ **不许**用"往 allowlist 里加一行"把它糊过去——allowlist 只会把问题从"红"变成"没人看的白名单"。
      若某条 SKIP 是**合法的**（平台 API 就是不存在），修法必须是把那条判据**搬到平台层**
      （`//go:build`，票 82/81 已经做过两次，照做），让它在不适用的平台上**不参评**，
      而不是在适用的平台上"跳过并算过"。
- [x] **AC#3** 双向变异：(i) 把机制改成 no-op ⇒ **必须有既有用例红**（不是本票新写的）；
      (ii) 往快照里**种一条 `t.Skip`** ⇒ portable 步必须红并点名它。锚点=承载行为那一行，
      同链 grep 证落地，**编译失败不算变异**，还原后 `git diff --quiet` 证干净。
- [ ] **AC#4** CI 上有**真实 run id + 步级结论**（`gh run view --job` 读**step**，不读 job status；
      并区分 failure / cancelled（`total_count=0` 不是样本）/ 未跑完）。⚠ 代理**不 push** ⇒
      这条若要 push 才能拿到，就在票面明写"欠编排者 push 后复跑"，**不许用"本地跑过了"替代**。
- [x] **AC#5** 门禁（只跑自己碰的范围）：`gofmt -l` 空、改动脚本 `bash -n` / `go vet` rc=0、
      本地按新 runner 跑一遍并贴出**四数**（`=== RUN` / PASS / FAIL / **SKIP 必须为 0 或被逐条点名**）。

## Rules（本仓固定）

`git commit -q -F - -- <显式路径> <<'MSGEOF'`（引号 heredoc）；提交前 `git diff --cached --name-only`；
禁 `git add -A`/`.`、`--amend`/`reset`/`rebase`/`stash`/`checkout .`（A34）；**不 push**；
不在仓内建 worktree（A38④，快照 `git archive <sha> | tar -x -C /tmp/<带会话后缀>`）；
不跑整仓门禁（共树：`internal/winsec/**` 票 89、`frontend/`+`internal/panel/` 票 77 有人在写，别碰）；
票面 append-only，改行前先读，在标题前插段落要把标题重抄进 `new_string` 且 `git diff --numstat` 删除列为 0；
四种假绿逐条点名；数字不达标就 FAIL 附数字，不许调阈值、不许挑运气那次、多样本全报。
15 次工具调用内交回第一枚 checkpoint；每次提交同步 Status + 勾框 + `next=`；接近轮数上限主动收尾留断点。

## Progress log（append-only）

- 2026-09-21 17:1x（编排者）：建票。来源是票 82 报告的两条残留之一："`TestSyncRegistryProbeLive` 两侧都 SKIP，
  而 CI 的 portable 步是裸 `go test` ⇒ SKIP 记成 ok（只有 `tools/d22scan/runtests.sh` 拒 SKIP）"。
  我把它单独开票而**不塞进票 82**，因为票 82 的界是那 8 条红；顺手改它=扩界（它自己也是这么判的）。
  ⚠ 判据里我最在意的一条是 AC#2 那句"**不许用 allowlist 糊过去**"：这票的修法如果让白名单变长，
  那它就是把"红"改成了"没人读的清单"，等于没修。
  next= 派单（它不依赖任何在飞的票，可以插队；但**先等 `internal/panel`/`frontend/` 那两路收敛**，
  因为它要动 `ci.yml`，而票 77 也在等 CI 结论 ⇒ 同文件风险）。

### agent-ticket93 交件段（append-only，2026-09-21 17:4x-18:1x）

- [x] **AC#1** 全量点名已完成（数据在下面第一枚 checkpoint 里，快照 `a1613d9`，两平台各一次 `go test -v -count=1`，
      命令与抓取证明见 Progress log 的 "AC#1 checkpoint" 段）。

#### AC#1 checkpoint：portable 16 包两侧全量 SKIP 台账（**全部是 `-v` 取数**）

取数命令（在 `git archive a1613d9 | tar -x -C /tmp/t93-agent-ticket93` 的快照里跑，不在仓库目录内跑）：

```
go test -v -count=1 <16 包>            # Windows，本机 go1.27.1 windows/amd64
docker run ... golang:1.27 sh -c 'go test -v -count=1 <16 包>'   # ubuntu 侧（linux/amd64 容器）
```

抓取完整性判据（同一条命令产出的两个数对照，证明 grep 没漏）：

| 侧 | 日志总行数 | `=== RUN` | `^--- PASS` | `^--- FAIL` | `^--- SKIP` | rc |
|---|---|---|---|---|---|---|
| Windows | 3265 | 986 | 617 | 0 | **7** | 0 |
| linux（golang:1.27 容器） | 2957 | 867 | 536 | 1 | **8** | 1（见下方 F-1） |

两侧 126 行差、119 个 `=== RUN` 差 = 平台后缀/`//go:build windows` 的测试文件在 linux 侧不参评，
`go list` 对照：`GOOS=linux ./internal/risk` 仍编译 12 个测试文件、`GOOS=windows` 15 个（包在两侧都在门禁里）。

逐条（包 / 测试名 / file:line / 原因 / 哪侧跳）：

| # | 包 | 测试名 | file:line | 跳过原因（逐字摘自日志） | Win | Linux |
|---|---|---|---|---|---|---|
| 1 | internal/agent/approval | TestDefaultDeadlineWallClockMeasurement | ticket84_no_owner_test.go:224 | 有意慢：300s 墙钟计量，只在 WISP_84_MEASURE=1 时跑 | Y | Y |
| 2 | internal/memory | TestSubprocessCrashWriter | concurrent_test.go:186 | crash-writer subprocess; runs under TestCrashRecoveryKillMidWrite | Y | Y |
| 3 | internal/risk | TestSyncRegistryProbeLive | syncdirs_test.go:347 | no registry-grade sync record on this machine (HKCU Accounts without UserFolder...) | Y | Y |
| 4 | internal/models | TestRealDownloadVadThroughPipeline | manifest_real_test.go:120 | real-network spot check; set WISP_IT_REAL_MIRROR=1 to run | Y | Y |
| 5 | internal/models | TestRealDownloadPuncArchiveThroughPipeline | manifest_real_test.go:147 | real-network spot check; set WISP_IT_REAL_MIRROR=1 to run | Y | Y |
| 6 | internal/audio | TestLiveWasapiSmoke | hotplug_test.go:527（文件头 `//go:build windows`） | live WASAPI smoke requires WISP_LIVE_MIC=1 and a real microphone | Y | 不参评 |
| 7 | internal/proc | TestHelperProcess | jobscope_windows_test.go:87（`//go:build windows`） | helper process mode not set | Y | 不参评 |
| 8 | internal/risk | TestC26RewrittenSyncRootDoesNotDisarmSuspectNet | pathresolver_rewrite_account_test.go:138 | platform gives no handle-resolved form for an existing directory | N（跑过并 PASS） | Y |
| 9 | internal/tools | TestD34WriteMatrix | fs_write_test.go:172（helper `otherVolumeDir` fs_write_test.go:140/151） | no volume identity for the reference path on this platform | N | Y |
| 10 | internal/tools | TestCrossVolumeMoveStopsWithTwoCopiesOnLateStop | fs_write_test.go:677（同一 helper） | no volume identity for the reference path on this platform | N | Y |

F-1（顺手发现的**真红**，不是 SKIP，属票 102 的账，登记不修）：linux 侧
`--- FAIL: TestPathCanonicalizerAccountsForRewrittenRoots`，`paths_rewrite_ticket102_test.go:64`：
`InAllowlist("/tmp/TestPathCanonicalizerAccountsForRewrittenRoots.../proj/a.txt") = false although the
expanded root is a real tree`。Windows 侧同一条 PASS ⇒ 这是一条只在 POSIX 上红的用例，
今天 CI 的 portable 步（裸 `go test`）会报 FAIL 而不是假绿，但它**从来没被本票之外的任何本地门禁看到过**。
我没有改它（不是我的界，且禁改断言）。next= 编排者决定是否回派票 102。

#### AC#2/AC#3/AC#5 交件段（同一快照谱系：`c8c828e`，父 `360efdf`；所有运行都在 /tmp 快照，仓库目录内没有建 worktree）

落了什么：

1. `scripts/portable-tests.sh`（新）。判据**全部**来自既有严仪器 `tools/d22scan/runtests.sh`
   （票 71 AC#3：非零 go test 退出码原样传播 / 任何 `--- SKIP` 即 fatal / 零 top-level 结果即 fatal，
   并强制 `-v -count=1`）。本脚本在其上只加一层"显式记账"：`name|pkg|platform|class|reason` 十条台账，
   每次运行都用 `go test -list '^<name>$' <pkg>` 对**本平台已编译的测试二进制**复核条目仍在，
   不在 ⇒ 这一步红（改名/删除/被 tag 挪走都会当场暴露）。这是它与 `allowlist.txt` 的唯一实质差别：
   白名单是单向且沉默的，台账是双向且每次自证的。`tools/d22scan/**` 与 `allowlist.txt` **一字未动**
   （`git show --stat c8c828e` 的文件清单里没有它们）。
2. `ci.yml`：`test-core` 的 portable 步改成 `bash scripts/portable-tests.sh`（16 包清单收进脚本，单一真相源）；
   `test-windows` 的 "Portable windows tests (proc/secret/config)" 原来也是裸 `go test`（proc 的
   `TestHelperProcess` 每次必 SKIP）⇒ 同治。加严一道并登记：把 `Environment fork assertion` 提到 portable 步**之前**
   —— 它原来排在一道可能红的步骤后面，而 GitHub 在首个失败步后停止执行 ⇒ portable 步一红它就从未跑过（A44-1 形状，
   票 71 在 lint job 里修过同一个）。没删步骤、没改阈值、没放宽任何断言。
3. 平台层搬迁（AC#2  prescribed 的正解，仅一条）：`TestSyncRegistryProbeLive` 函数体**逐字**移入既有
   `//go:build windows` 的 `internal/risk/syncdirs_windows_test.go`，`syncdirs_test.go` 留指针注释。
   POSIX 侧从此**不参评**而不是"跳过并算过"；Windows 侧断言原样在。
   登记手法与合法性判定（重要；同一手法在不同包合法性相反）：这里用的是 build tag（该用例的对象是注册表句柄，
   POSIX 上没有该对象，属"平台 API 天生不存在"），**且只 tag 一个测试文件**；
   整包没有被排除 —— `go vet ./internal/risk/` 在 GOOS=windows 与 GOOS=linux 两侧都 rc=0，
   `go list` 计得 linux 编译 12 个测试文件 / windows 15 个（两侧都在门禁里）。
   剩下两条同类（POSIX 无对象）SKIP **没有**走这条路，因为够不到：`TestC26RewrittenSyncRootDoesNotDisarmSuspectNet`
   在 `internal/risk/pathresolver_rewrite_account_test.go`，命中本票禁改模式 `pathresolver*.go`；
   `TestD34WriteMatrix` / `TestCrossVolumeMoveStopsWithTwoCopiesOnLateStop` 需要第二个卷
   （`otherVolumeDir`，`internal/tools/fs_write_test.go:140/151`），CI ubuntu 容器里 /tmp 与 / 同一个文件系统。
   这三条进台账并各带 remedy 指针，不静默吸收。

AC#3 双向变异（快照 `/tmp/t93m-agent-ticket93` = `git archive c8c828e`；每发都在**同一条 `&&` 链**里
`grep -n` 打印被改后的整行；先 `go build` / `go vet`(编译测试文件) rc=0 再跑；全部 `-v`）：

| 变异 | 改后落地锚点（grep -n 实测行号） | 编译证明 | 结果 |
|---|---|---|---|
| M1a 台账 no-op（`skip_pattern` 换成不可能匹配的名字，严仪器保留） | `scripts/portable-tests.sh:182` | n/a（未动 Go） | **linux rc=1**，`top-level: PASS=183 FAIL=1 SKIP=6` ⇒ **既有仪器 `tools/d22scan/runtests.sh`（票 71，非本票新写）当场红**并逐条点名 6 条 SKIP；本脚本另打印 "unaccounted SKIP lines" 带 file:line |
| M1b 严仪器 no-op（`sh "$strict" …` 换成票 93 之前的裸 `go test "${scope[@]}" -count=1`） | `scripts/portable-tests.sh:190` | n/a | 病回来了：`ok …/internal/memory`、`ok …/internal/models`、`ok …/internal/risk` 三行 **ok 里含着 6 条 SKIP**；四数量到 `=== RUN=0 PASS=0 SKIP=0`（裸 go test 不打印 per-test 行 ⇒ 证明"没有 `-v` 就没有台账"）。同一 scope 跑既有仪器 = `SKIP=6 rc=1`。（这次 M1b 的 rc=1 只来自 `internal/tools` 里那条与本票无关的 F-1 FAIL，不是 SKIP —— 正是本票要修的"SKIP 藏在 ok 里"） |
| M2 种一条 `t.Skip`（AC#3(ii)） | `internal/observe/clock_test.go:13`（`t.Skip("ticket93 seeded skip probe - not a real skip")`） | `go build ./internal/observe/` rc=0 + `go vet ./internal/observe/`（含测试文件）rc=0 | **windows rc=1**：四数 `=== RUN=985 PASS=621 FAIL=1 SKIP=1`，点名块 `clock_test.go:13: ticket93 seeded skip probe` + `--- SKIP: TestTimeoutIsMonotonic`；**linux rc=1**：四数 `=== RUN=860 PASS=536 FAIL=1 SKIP=1`，同一名 |
| 还原 | — | — | `diff` 快照文件 vs `git archive c8c828e` 同一文件 = 空；`cmp` 脚本 vs 原件 rc=0；仓库工作树里从未出现变异（`git status --porcelain` 只剩我自己的票面 + 邻居的 `internal/winsec`） |

AC#5 本地四数（新 runner，默认 16 包全量，**两侧都带 `-v`**；同一棵树 `c8c828e`）：

| 侧 | `=== RUN` | `--- PASS` | `--- FAIL` | `--- SKIP` | rc | 台账生效条数 |
|---|---|---|---|---|---|---|
| windows（本机 go1.27.1） | 980 | 617 | 1（F-2） | **0** | 1 | 7/10 |
| linux（golang:1.27 容器） | 859 | 536 | 1（F-1） | **0** | 1 | 7/10 |

方差如实全报：windows 侧另两次同树复跑量到 `=== RUN=985 PASS=621`（M2 那次）与基线 `980/617`，
差 5 个 top-level `=== RUN` / 4 个 `--- PASS`；**变异前后包集合完全相同**（两次都是 20 个测试包 + 3 个 `[no test files]`），
所以这是这台机器上的用例数量方差（条件性 `t.Run` 或 fixture 枚举），不是覆盖变化 —— 我没有把它抹平，
两条数都在表里，SKIP 轴在两棵树上都是 0（基线）/1（M2），不受它影响。
其余门禁：`gofmt -l internal/risk/` 空；`bash -n scripts/portable-tests.sh` rc=0；
`sh scripts/d22scan.sh`（快照内）clean，八项 scope 计数全在输出里；我的新文件与追加段零非 ASCII 字符（ban #8）。

AC#4 **未勾**：本代理不 push，`c8c828e` 与 `df0a1e0` 都还没进任何 run
（`gh run list --limit 50 --json headSha` 里两个 sha 各 `grep -c` = 0）⇒ **欠编排者 push 后复跑**。

拿得到的**修前**步级读数（`gh run view 35590599782 --json jobs` 读 steps 数组，不读 job 颜色；
该 run headSha=`a64d06fb8cb6d3b1055a17deb5ac423f6bbc2853`，job conclusion=failure）：

    test-core :: Portable package tests (...)      :: success
    test-core :: Environment fork assertion        :: success

这就是本票要钉的形状在 CI 台账上的证据：那一步自己判 `success`，而它内部两侧各含着 7-8 条 `--- SKIP`
（AC#1 的台账）——"步级 success"与"结论产出过"不是同一件事。注意该 run 的 headSha 比 360efdf 旧，
所以 **F-1 进树之后的 CI 还没有任何 run**：我不能说那道后门"已经从未执行过"，只能说
**一旦 portable 步红（F-1 让它在 linux 上必红），GitHub 在首个失败步后停止执行，后面的门就不会跑**；
本票已把 `Environment fork assertion` 提到该步之前，属加严（未删步骤、未改阈值）。
push 后的期望读数与勾框条件：日志含 `portable-tests.sh: platform=linux ... four numbers
(all from -v output): === RUN=… --- PASS=… --- FAIL=… --- SKIP=0`，且该步 conclusion 为 failure 并点名
`TestPathCanonicalizerAccountsForRewrittenRoots`（这一步红是**正确的**，是该步第一次产出真结论）；
把 run id + job id + 该步 conclusion 原文回填本段才勾框。另核 `test-windows` 的
"Portable windows tests"（现在也走同一个 runner）。

顺手发现（都不在我界内，登记不修）：
- **F-1**（linux 真红，包 = `internal/tools`，不是我上面初判的 risk）：`--- FAIL: TestPathCanonicalizerAccountsForRewrittenRoots`，
  `internal/tools/paths_rewrite_ticket102_test.go:64`，`InAllowlist("/tmp/…/proj/a.txt") = false although the expanded root is a real tree`。
  windows 侧 PASS ⇒ 只在 POSIX 红。裸 `go test` 也会报这条（FAIL 不像 SKIP 那样被记成 ok），所以它不是假绿；
  但它让 `test-core` 的 portable 步在 linux 上必红，并把后面的门挡掉（见上）。next= 回派票 102 的 owner。
- **F-2**（windows 真红）：`internal/risk/pathresolver_budget_norace_test.go:37`
  `TestResolvePerCallBudget: C26 Resolve 2.202 ms/op, budget 1.000 ms, 1207 samples`。
  `a1613d9` 上同一批包是 FAIL=0（`=== RUN=986 PASS=617 SKIP=7`），到 `360efdf`（父链含票 103 的 winsec 改动）
  变红。阈值属禁改面，我没动。next= 票 103/106 或编排者判它是不是真回归。

next=（本票留给下一手的具体命令）：
1. 编排者：push `c8c828e` ⇒ 复跑 AC#4（命令见上），把步级结论回填本段并勾框。
2. `internal/risk/pathresolver_rewrite_account_test.go` 的 owner：把
   `TestC26RewrittenSyncRootDoesNotDisarmSuspectNet` 的 POSIX 分支搬到 `//go:build`（AC#2 正解），搬成之后
   从 `scripts/portable-tests.sh` 删掉那条台账 —— 删了不会红（红的是"条目还在但测试不在"），
   留着也不会掩盖新 SKIP。
3. `internal/tools` 的 owner：`otherVolumeDir` 在 CI ubuntu 里能否用 tmpfs 造第二卷；造得出就删两条台账。
4. 谁碰 `test-windows`：把 `-run TestSyncRegistryProbeLive` 接进 windows job（用 `runtests.sh`），
   那才是"这条只在 Windows 上参评"的真结论；今天它在 windows job 里没有任何步骤跑它。

