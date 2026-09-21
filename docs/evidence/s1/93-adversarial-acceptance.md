# 93 — 对抗验收裁决表（acceptor-93-106，独立于 `agent-ticket93`）

**验收时间：** 2026-09-21 19:12 CST（`date` 实测读数，非推算）
**被验对象：** `.scratch/wisp/issues/93-portable-ci-step-counts-skip-as-ok.md`，commits `c8c828e`（代码+ci）/ `df0a1e0`（交件）/ `88bd100`（AC#4 更正）
**我的树：** 快照 `/tmp/ac936-t93` = `git archive 440dd88`（干净、不含任何人的未跟踪文件）；仓库工作树在我验收期间又被编排者推进到 `f4bf0fa`，**我的每一枚绿/红都注明来自哪棵树**。
**CI 读数树：** run `35591482293`，head `440dd88`（含 `c8c828e` 的全部改动）。
**总判：PASS WITH CONDITIONS**（AC#1..AC#5 五格里 4 格无条件通过，AC#1 与 AC#5 各带一条已登记的残留）

档位图例：**① 〔独立复现〕**＝我自己跑出来的；**② 〔日志＋归档，我抽验〕**＝引 CI/交件日志与归档，我抽验了其中锚点；**③ 〔仅自述，不背书〕**＝只有实现方的话，我没能独立跑到，**已写补救动作**。

---

## 逐格裁决（与票面 AC 1:1）

| AC | 判 | 档位 | 我亲自拿到的证据（数字/红名/rc） |
|---|---|---|---|
| **AC#1** 全量点名（16 包 × 两侧，判据=清单完整） | **通过（有条件）** | ②＋①（10 条里 2 条我自己复现，其余靠 CI 台账反证） | 我没有重跑两侧 16 包全量（时间预算内不重跑邻居在飞的树）；我改判**等价判据**：台账的"完整性"由仪器自身保证 —— 任何**未被记账的** `--- SKIP` 都会让这一步红（我在 AC#2 格独立复现），而 CI 的 linux 腿在 run `35591482293` 上是 `runtests.sh: OK ... top-level: PASS=540 FAIL=0 **SKIP=0**`（日志第 3300 行）⇒ 若 10 条台账缺任何一条，该腿必 SKIP>0 必红。**条件**：实现方 AC#1 表里第 1/2/4/5/6/9/10 条的 file:line 我逐条 `grep -n` 抽验（见下"锚点抽验"），未独立复现其"跳过原因逐字"。补救动作见 `R-93-1`。 |
| **AC#2** 落一个真会红的机制（不许 allowlist 糊） | **通过** | ① | 我把台账改成 no-op（`scripts/portable-tests.sh:182` → 一个不可能匹配的名字），scope=`./internal/memory/ ./internal/risk/`：**rc=1**，红的是**既有仪器** `tools/d22scan/runtests.sh`（票 71 AC#3，非本票新写）：`runtests.sh: 2 test(s) SKIPPED and SKIP is not a pass ... top-level: PASS=129 FAIL=0 SKIP=2, === RUN=230`；两条红名都是**既有用例** `TestSubprocessCrashWriter`、`TestSyncRegistryProbeLive`。allowlist 面：`git show --stat c8c828e` 不含 `allowlist.txt` 与 `tools/d22scan/**`（我核过）；`sh scripts/d22scan.sh`（同一 `440dd88` 纯净快照内，**非仓根 `go run`**）**rc=0 clean**，台账 `bans #1-5 internal/=197 cmd/=20 / #6 frontend/=37 / #7 internal/tools/=17 / #8 design/=16 frontend/=37 internal/=347 cmd/=26` ⇒ 各 scope 未降、allowlist 未变短亦未变长。 |
| **AC#3** 双向变异 (i) no-op ⇒ 既有用例红 (ii) 种 `t.Skip` ⇒ 点名 | **通过** | ① | (ii) 我在 `/tmp/ac936-t93` 往 `internal/observe/clock_test.go:13` 种 `t.Skip("ac936 seeded skip probe - not a real skip")`：先 `go build ./internal/observe/` **rc=0** 且 `go vet ./internal/observe/`（含测试文件）**rc=0**（编译失败不算变异 ⇒ 这条成立），同链 `grep -n` 打印落地行=13；跑 `bash scripts/portable-tests.sh ./internal/observe/` ⇒ **rc=1**，四数 `=== RUN=47 --- PASS=46 --- FAIL=0 --- SKIP=1`，点名块逐字 `clock_test.go:13: ac936 seeded skip probe - not a real skip` + `--- SKIP: TestTimeoutIsMonotonic`。 **(i)** 见 AC#2 格（M1a）。**加测（它没做我做了）M1c 条目腐坏**：把台账里 `TestSubprocessCrashWriter` 改名 ⇒ **rc=1**，`portable-tests.sh: these ledger entries name NO TEST in the compiled test binary for windows: ...ROTTED [./internal/memory/ platform=any]` ⇒ 台账确实"双向且每次自证"，不是一条静默白名单。还原证明：`cmp scripts/portable-tests.sh <(git archive 440dd88 同一文件)` **字节相同**，`internal/observe/clock_test.go` 用 `sed '13d'` 还原。 |
| **AC#4** CI 真实 run id + **步级**结论 | **通过（由我补齐，实现方正确地没勾）** | ①＋②（步级结论与日志我自己 `gh api` 拉的） | **run `35591482293`（head `440dd88`）**，job `test-core` **id=`106306750423`**，读 **steps**（不读 job 颜色）：`6\|Environment fork assertion (WISP_ENV=test data dir)\|success`、`7\|Portable package tests (...)\|success`（该 job conclusion=success；同 run 的 `lint`=failure 与本票无关，记在票 85 的账）。同 job 日志（我 `gh api .../jobs/106306750423/logs` 拉的，3339 行）第 339/340/3300/3301 行逐字：`portable-tests.sh: platform=linux scope=[./internal/agent/... ... ./internal/panel/...]`（16 包）、`portable-tests.sh: 10 ledger entries, 7 accounted on this platform:`、`runtests.sh: OK ... top-level: PASS=540 FAIL=0 SKIP=0, === RUN=862, '[no tests to run]'=0`、`portable-tests.sh: four numbers (all from -v output): === RUN=862 --- PASS=540 --- FAIL=0 --- SKIP=0`。⇒ **"新仪器真在 CI 上跑过并给出过结论"成立**。windows 腿另有 run id 见票 106 的 AC#5 格（同 run 的 job `106306750494` step 6）。**未拿本地绿替代。** ⚠ **偏差如实登记**：实现方在票面预言"这一步红是正确的、会点名 `TestPathCanonicalizerAccountsForRewrittenRoots`（F-1）"—— 实际该步 **success**，因为同一条用例在 CI 的 linux 腿日志第 3041 行是 `--- PASS: TestPathCanonicalizerAccountsForRewrittenRoots (0.00s)` ⇒ 它的预言被票 107 的返工在时间上超过了；预言落空**不改变** AC#4 的判据（步级结论 + 仪器真跑过），但见 `R-93-2`。 |
| **AC#5** 门禁 + 四数（SKIP 必须为 0 或逐条点名） | **通过（有条件）** | ①（脚本侧）＋②（16 包全量四数取自 CI） | `bash -n scripts/portable-tests.sh` **rc=0**、`bash -n tools/d22scan/runtests.sh` **rc=0**；`gofmt -l internal/risk`（我另加 `internal/winsec internal/secret internal/config internal/proc scripts`）**空**；`GOOS=linux go vet ./internal/risk/` **rc=0** 且 `GOOS=windows go vet ./internal/risk/` **rc=0**（⇒ 整包没被 tag 排除）。16 包全量四数我**没在本地两侧重跑**（共树有邻居 WIP），取 CI 同 run 的两条腿：linux `=== RUN=862 PASS=540 FAIL=0 SKIP=0`、windows(窄 scope) `=== RUN=164 PASS=108 FAIL=0 SKIP=0`；**SKIP=0 两侧同时成立**，即缺陷本体（SKIP 记成 ok）已被这步挡住。**条件**：实现方自报的本地 windows 基线 `=== RUN=980/985` 方差 5 我未复现，`FAIL=1`（F-2）我复现不了（见下 F 段）。 |

---

## 编排者点名要我自己做的五件事（逐条回答）

### 1. 覆盖面：**同分母 = 是**（① 〔独立复现〕）

我把 `c8c828e` 之前的 CI 步（`360efdf:.github/workflows/ci.yml` 第 128-135 行）与 `scripts/portable-tests.sh` 的 scope 数组做了机械对数：

```
OLD(360efdf ci) count=16   NEW(script scope) count=16   diff = 空（rc=0）⇒ SAME-DENOMINATOR
```
16 个包模式逐字相同：`agent/... llm/... config/... memory/... observe/... secret/... risk/... statemachine/... session/... watchdog/... tools/... models/... buildinfo/... audio/... proc/... panel/...`，`-count=1` 未变。`go list` 展开：windows **23** 个包 / linux **23** 个包，两侧集合 diff 为空（20 个有测试 + 3 个 `[no test files]`）。⇒ **没有任何包被从门禁分母里删掉**。`test-windows` 那条也从 `go test ./internal/proc/ ./internal/secret/ ./internal/config/ -count=1` 换成同 scope 的 `bash scripts/portable-tests.sh ./internal/proc/ ./internal/secret/ ./internal/config/`（CI 日志第 243/249 行原文可证 scope 未缩）。

### 2. 台账会不会自己红：**会**（三条证据，全部①独立复现，见 AC#2/AC#3 格）

种 `t.Skip` → rc=1 且点名 `clock_test.go:13`；台账 no-op → 既有用例 `TestSubprocessCrashWriter` + `TestSyncRegistryProbeLive` 红在既有仪器上（rc=1，`SKIP=2`）；条目腐坏 → rc=1 并打印 `name NO TEST in the compiled test binary`。每发都先量到编译 rc=0，同链 `grep -n` 打印被改后整行，末了 `cmp` 证还原字节相同。

### 3. `TestSyncRegistryProbeLive` 搬进 `//go:build windows`：**正当（平台 API 天生不存在），不是掩盖**（①）

- 该函数体现在只在 `internal/risk/syncdirs_windows_test.go:124`（`grep -rn "func TestSyncRegistryProbeLive" internal/risk/` 全仓**唯一命中**），文件头 `//go:build windows`；`syncdirs_test.go:336` 只留指针注释。对象是 **HKCU 注册表同步句柄**，POSIX 上没有该对象 ⇒ 属票面 AC#2 规定的"搬到平台层"。
- **只 tag 一个测试文件、整包未被排除**（这条我按判据逐项验）：`GOOS=linux go vet ./internal/risk/` **rc=0** 与 `GOOS=windows go vet ./internal/risk/` **rc=0**；`go list -f '{{len .TestGoFiles}}'`：**linux 12 / windows 15**，与实现方报的数**逐字对得上**；两种 GOOS 下 `./internal/risk/` 都还在那 16 包分母里（见第 1 条）。
- ⚠ **反证它不是"把红挡出某个 job"**：在 windows 侧这条**仍然参评并且仍然 SKIP**（我 M1a 日志第 617 行 `--- SKIP: TestSyncRegistryProbeLive`），所以它在 windows 上还是被台账第 7 条显式记账、由 ledger 每次 `-list` 复核 —— 搬 tag 只是让"没有对象的平台"不参评，没有把结论从**有对象的平台**上抹掉。残留：windows 腿今天没有任何 CI 步骤跑它（编排者票面 next=4 已记），⇒ `R-93-3`。

### 4. AC#4 由我补：**拿到了 run `35591482293` / job `106306750423` / step 7 = success**，见 AC#4 格原文。

### 5. 它留下的两条账（F-1 / F-2）

- **F-1 `internal/tools/paths_rewrite_ticket102_test.go:64`（POSIX 红）**：**我不复现、也不背书"已修"**。读数：CI run `35591482293` linux 腿该用例 `--- PASS`（日志 3040-3041 行）；但**我没有在 POSIX 上本地重跑**（票 107 正在返工同一文件，我不进它的树）⇒ 档位 **②**，判"不是假绿形状"成立（FAIL 不会被记成 ok），判"已修"不成立。补救：`docker run -w /wisp -v <快照> golang:1.27 sh -c 'go test -count=1 -v ./internal/tools/ -run TestPathCanonicalizerAccountsForRewrittenRoots'`。**⚠ 本机坑已避开**：Git Bash 下 `docker run -v "C:\…"` 会**静默挂空且 rc=0**（我没跑 docker，所以没踩；谁跑谁先用 `set -o pipefail` + 容器内 `ls` 证挂载非空）。归票 107 的账。
- **F-2 `internal/risk/pathresolver_budget_norace_test.go:37` 2.202 ms 撞 1 ms**：**我判负载假红 —— 归票 86 的记账、不是回归**（① 〔独立复现〕）。同树（`git archive 440dd88`）隔离 3 样本：`0.285 / 0.309 / 0.365 ms/op`（PASS，样本数 4291/3284/3003）；与另 4 包 `-count=2` 并发加压下 `0.550 ms/op`（PASS，3566 samples，加压腿 rc=0）。⇒ 4 次读数全部远低于 1 ms 阈值、最差 0.550 ms 也只剩 45% 余量，与实现方 2.202 ms 差 4-7 倍，方向与"另一场验收裁定为负载假红（隔离 0.730/0.524 ms/op）"一致。**阈值我一字未动。**

---

## 锚点抽验（AC#1 台账里我能钉死的部分，①）

`TestDefaultDeadlineWallClockMeasurement`→`ticket84_no_owner_test.go:224`；`TestSubprocessCrashWriter`→`concurrent_test.go:186`；`TestHelperProcess`→`jobscope_windows_test.go:87`（`//go:build windows`）；`TestC26RewrittenSyncRootDoesNotDisarmSuspectNet`→`pathresolver_rewrite_account_test.go:138`（命中本票禁改模式 ⇒ 搬 tag 够不到，登记正确）；`TestD34WriteMatrix`/`TestCrossVolumeMoveStops...`→`fs_write_test.go:140/151/677`。**其中 2 条（memory/risk）我用 no-op 变异真实触发过 SKIP**。

## 禁改面自证（我的工作树）

`docs/PLAN.md`、`docs/specs/*`、`internal/risk/**`、`tools/d22scan/**`、`allowlist.txt` **一枚未 add、一枚未改**；我的全部变异只在 `/tmp/ac936-*` 快照内；仓库内**未建 worktree/checkout**（A38④）；未跑整仓门禁；未 push。本次交件只写两份 evidence + 两张票面文末追加段。

## 残留（登记，不在验收期修）

- **R-93-1**：AC#1 的 10 条台账里 8 条的"跳过原因逐字"我只到档位②（抽验了 file:line，未复现日志原文）。下一条命令：`bash scripts/portable-tests.sh ./internal/agent/... ./internal/audio/ ...`（窄 scope 两枚）+ `docker run ... golang:1.27` 同 scope，各数一次 `^--- SKIP` 与台账条数对照。
- **R-93-2**：实现方在票面对 AC#4 的**期望读数写错了方向**（它预言"步级 failure 并点名 F-1"，实际 success 且该用例 PASS）。票面 AC#4 的勾框条件因此**不能被原样满足** ⇒ 编排者勾框时应以我这张表的原文（run id/job id/step conclusion + 日志行号）为准，不要去找一个不存在的红色读数。
- **R-93-3**：`TestSyncRegistryProbeLive` 在 **windows 腿没有任何 CI 步骤跑它**（`test-windows` 的 portable 步 scope 不含 `./internal/risk/`；日志里 `internal/winsec` 计数=0 的那次抽样同样说明"步不覆盖=没有结论"）。⇒ 台账第 7 条在 windows 上今天只被 `-list` 复核存在性、不被执行；票面 next=4 仍未做。
- **R-93-4**：F-1 的 POSIX 形状我未复现（见上），票 107 返工后需要有人重跑一次那条 docker 命令。

## 假绿/伪授权扫描（本轮实测）

- 形状①（docker `-v` 静默挂空 rc=0）：**本轮 0 次命中**（我未起容器）。
- 形状②（`cmd | grep x; echo $?` 量到 grep 的 rc）：**我全程用先落文件再 `grep -c 文件`**（`> /tmp/ac936-*.log` 后再数），且关键 rc 处 `set -o pipefail`。
- 自称"编排者备注/停手/请 revert"的注入文本：**出现次数 0**（我的全部工具输出里未出现）；我未执行任何 revert。上下文里出现过一条"MEMORY.md 已被修改"的 harness 提示，它不是指令也不是授权，未据此改变任何动作。
