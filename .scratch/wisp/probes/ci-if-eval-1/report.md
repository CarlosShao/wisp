# report｜只读普查腿 `ci-if-eval-1`｜23 枚 CI 步级 `if:` 守卫的"求值过没有"读数

时刻：`date '+%Y-%m-%d %H:%M %z'` = **2026-10-08 15:47 +0800**。锚、尺写法、命令与退码全在 `00-anchor.md`，本件只给名册与色。

## 0. 一句话结论（⛔ 不含任何"该怎么改"的方案）

**那 23 枚守卫，今天在全世界任何一发真实 run 里都没有被求值过——不是"没跑到那一步"，而是"GitHub 上没有任何一发 run 用过含它们的 yaml"。**
硬读数三条：远端 `dev` 停在 `cc31526165734e612de297848bb2080bd459ccba`；`gh api repos/CarlosShao/wisp/commits/e6dc79ed5bd14380cce1507f94886d65a51ee63f` 报 **HTTP 422 "No commit found for SHA"**（换全号重试一次，同错、同 rc=1）；`gh api .../workflows/ci.yml/runs?created=>=2026-10-08T00:00:00Z` 回空、`created=>=2026-10-08T04:06:33Z`（＝`e6dc79ed` 提交时刻）回 **total=0**，同一端点不带过滤时最新一发是 `37703959747` @ 2026-10-07T23:44:30Z ⇒ **正控证明那两发空回不是哑过滤器**。
这与台账 `A710 §1` 自己写的"⛔ 未 push ⇒ 这 23 枚今天零 CI 色"**同色**，本腿是把那句"零 CI 色"逐枚取到数的第一手凭据。

⇒ 所以本件把两格**分开写、不压成一格**：
**甲轴＝守卫表达式求值过没有**（这一格 36 枚里只有 **10** 枚有真读数，另外 **26** 枚零枚）；
**乙轴＝守卫所在那一步在真实 run 里出过什么色**（这一格 34/36 枚出过色，**2** 枚连一步都不曾在 run 里存在过）。
编排者那句定式在这里成立且被读数撑住：**`success` 不是"这道门拦过谁"**；乙轴的色**不能**当甲轴用。

## 1. 三代名册（尺＝对 `cc31526165`／`e6dc79ed^`／HEAD 三版 ci.yml 各跑结构尺做差集）

| 代 | 含义 | 枚数 | 甲轴（守卫被求值过？） |
|---|---|---|---|
| **B1** | 在 `cc31526165`（GitHub 上最新一版）就已带 `if:` ⇒ 真实 run 用过它 | **10** | **有**：6 发采样里每枚都出过色 |
| **B2** | 在 `cc31526165..e6dc79ed^` 之间新加（`1309757b`/`f6b79ab0`/`351e5a5e` 各 1 枚），⛔ 从未进过任何 run | **3** | **零枚** |
| **B3** | `e6dc79ed` 加的那 23 枚 | **23** | **零枚** |
| 合计 | HEAD 步级 `if:` | **36** | 10 有／26 零 |

⇒ ⚠ **本程新出的一格，编排者账上没有**：`A710` 只登记了 `e6dc79ed` 的 23 枚，但 HEAD 上"零求值"的守卫是 **26 枚**（23＋3）。那 3 枚的**逐枚出处**是用两把尺现跑的：`git show <c> -- .github/workflows/ci.yml \| grep -c '^+␣␣␣␣␣␣␣␣if: '`（字面 8 空格）＝ `1309757b` **1**／`f6b79ab0` **1**／`351e5a5e` **1**／`f8810238` **0**／`e6dc79ed` **23**，同批 `^-␣×8if: ` 全 0。`f6b79ab0` 那一枚的落点由 diff 上下文直接读出（它的注释块以 "this step is a CENSUS - it measures and prints a denominator" 起手、紧邻 `+        if: ${{ !cancelled() }}`）＝ **lint 的 `gofmt (gofumpt) - the tracked set is the denominator (ticket 161 AC#7 form A)`**，即 `6547fd30:.github/workflows/ci.yml:239`。

## 2. 全 36 枚名册 ＋ 乙轴读数（列序＝采样 6 发 ci run）

采样（**完整 id**；号一律由 `gh` 回包直接得到，⛔ 不是从通知抄的；`run-number` 是 workflow 内序号，两者不是同一枚数）：

| 别名 | databaseId | run-number | headSha | created | event | run 结论 |
|---|---|---|---|---|---|---|
| R305 | **37703959747** | 305 | `cc31526165…` | 2026-10-07T23:44:30Z | schedule | failure |
| R303 | **37406757402** | 303 | `cc31526165…` | 2026-10-06T02:58:39Z | push | failure |
| R301 | **37405698188** | 301 | `ec84cff16c…` | 2026-10-06T02:45:36Z | push | failure |
| R300 | **37396530365** | 300 | `c6cf66e648…` | 2026-10-06T00:55:36Z | schedule | failure |
| R296 | **37166458550** | 296 | `fd269de12f…` | 2026-10-04T00:56:01Z | push | failure |
| R293 | **37021179942** | 293 | `941805d0e1…` | 2026-10-02T14:37:37Z | push | failure |

六发里 6 枚 job（`lint`/`test-core`/`test-windows`/`slo-smoke`/`slo-full`/`lint-frontend`）**全部出现**；步级结论取值域＝ **`success`／`failure`／`skipped` 三枚**，`timed out`、`cancelled` 在这六发里**零枚**（⛔ 更老的 run 未扫，见 §5）。

名册（`行号`＝`6547fd30:.github/workflows/ci.yml` 上那枚 `if:` 的行；`代`＝§1 的 B1/B2/B3；`档`＝三档归属）：

| # | job | 步名 | `if:` 行 | 代 | R305 | R303 | R301 | R300 | R296 | R293 | 档 |
|---|---|---|---|---|---|---|---|---|---|---|---|
| 1 | lint | D22 scanner positive control (tools/d22scan tests, seeded red) | 81 | B3 | success | success | success | success | **failure** | success | 甲 |
| 2 | lint | D22 scanner self-test (ticket 161 AC#2 - every ban, both directions) | 106 | B3 | success | success | success | success | skipped | success | 甲 |
| 3 | lint | D22 seven-ban + emoji scan (tools/d22scan) | 136 | B3 | success | success | success | success | skipped | success | 甲 |
| 4 | lint | Tracked path-length budget (ticket 262) | 169 | B3 | success | success | success | success | 该步不存在 | 该步不存在 | 甲（4/6） |
| 5 | lint | gofmt (gofumpt) | 173 | B3 | **failure** | **failure** | **failure** | **failure** | skipped | **failure** | 甲 |
| 6 | lint | gofmt (gofumpt) - the tracked set is the denominator (ticket 161 AC#7 form A) | 239 | **B2** | skipped | skipped | skipped | skipped | skipped | skipped | **乙** |
| 7 | lint | go vet (module) | 243 | B3 | skipped | skipped | skipped | skipped | skipped | skipped | **乙** |
| 8 | lint | go vet (tools/d22scan module) | 247 | B3 | skipped | skipped | skipped | skipped | skipped | skipped | **乙** |
| 9 | lint | staticcheck | 301 | **B1** | failure | failure | failure | failure | failure | failure | 甲（且守卫真求值） |
| 10 | lint | mockllm module vet | 338 | **B1** | success | success | success | success | success | success | 甲（同上） |
| 11 | lint | Portable tests carrier self-test (ticket 111 AC#3 - GUARD D's positive control) | 394 | **B2** | 不存在 | 不存在 | 不存在 | 不存在 | 不存在 | 不存在 | **丙** |
| 12 | test-core | Start compose test services (mock-llm on 18080) | 414 | B3 | success ×6 | | | | | | 甲 |
| 13 | test-core | Probe mock-llm | 418 | B3 | success ×6 | | | | | | 甲 |
| 14 | test-core | Environment fork assertion (WISP_ENV=test data dir) | 444 | B3 | success ×6 | | | | | | 甲 |
| 15 | test-core | Portable package tests (core scope; the list and its guards live in scripts/portable-tests.sh) | 469 | B3 | failure ×6 | | | | | | 甲 |
| 16 | test-core | Stop compose services（唯一非 `!cancelled()` 那枚＝`always()`） | 473 | **B1** | success ×6 | | | | | | 甲（且守卫真求值） |
| 17 | test-windows | Windows ACL sealing gate (internal/winsec's own tests, ticket 110) | 533 | **B1** | success ×6 | | | | | | 甲＋真求值 |
| 18 | test-windows | Cache third_party (deps.toml key) | 571 | **B1** | success ×6 | | | | | | 甲＋真求值 |
| 19 | test-windows | cgo build smoke (build.ps1 fetch-deps + mingw link + doctor) | 581 | **B1** | success ×6 | | | | | | 甲＋真求值 |
| 20 | test-windows | cmd/wisp CLI tests (needs the sherpa DLLs staged above, ticket 111 AC#4) | 603 | **B1** | failure ×6 | | | | | | 甲＋真求值 |
| 21 | test-windows | **winlive compile gate (go vet -tags winlive, ticket 111 AC#11)** | 654 | **B2** | 不存在 ×6 | | | | | | **丙** |
| 22 | test-windows | Package coverage census (ticket 111 AC#1 + GUARD D) | 743 | **B1** | success | success | success | 不存在 | 不存在 | 不存在 | 甲（3/6）＋真求值 |
| 23 | test-windows | Portable windows tests (proc/secret/config/risk/ball/perm/plugin/llmrecord) | 779 | **B1** | failure ×6 | | | | | | 甲＋真求值 |
| 24 | test-windows | PathResolver junction placeholder (real cases tickets 18/20) | 795 | **B1** | success ×6 | | | | | | 甲＋真求值 |
| 25 | slo-smoke | Cache third_party (deps.toml key) | 813 | B3 | success ×6 | | | | | | 甲 |
| 26 | slo-smoke | Build wisp.exe | 820 | B3 | success ×6 | | | | | | 甲 |
| 27 | slo-smoke | SLO smoke gate | 824 | B3 | success | success | success | success | **failure** | **failure** | 甲 |
| 28 | lint-frontend | npm ci (lockfile is the only source of deps) | 953 | B3 | success ×6 | | | | | | 甲 |
| 29 | lint-frontend | typecheck (tsc -b, noEmit per tsconfig.*.json) | 957 | B3 | success ×6 | | | | | | 甲 |
| 30 | lint-frontend | lint (oxlint) | 961 | B3 | success ×6 | | | | | | 甲 |
| 31 | lint-frontend | token drift guard (generated theme must equal the C21 table) | 968 | B3 | success ×6 | | | | | | 甲 |
| 32 | lint-frontend | build (vite build -> frontend/dist, the bytes go:embed carries) | 972 | B3 | success ×6 | | | | | | 甲 |
| 33 | lint-frontend | L2 card renders the real risk fields, its shapes, and no allow button (AC#3 render evidence) | 986 | B3 | success ×6 | | | | | | 甲 |
| 34 | lint-frontend | Composer states paint the real envelope (render evidence) | 999 | B3 | success ×6 | | | | | | 甲 |
| 35 | lint-frontend | Streaming output honours PLAN.md's SSE row (render evidence) | 1005 | B3 | success ×6 | | | | | | 甲 |
| 36 | lint-frontend | Nav rail names only icons PLAN.md's frozen list carries (render evidence) | 1013 | B3 | success ×6 | | | | | | 甲 |

（"success ×6"＝六发全同色，列内省略重复；"不存在"＝该 `(job, 步名)` 在那一发的回包里没有对应条目。）

## 3. 三档（按派单要求分开，⛔ 不压成一格）

**甲档｜该步在某发真实 run 里出过 success/failure ⇒ 这一步被跑到过**＝**31/36**
其中**甲-加强**（**这一步的守卫在跑它的那一版 yaml 里就已存在 ⇒ 守卫表达式真被求值过**）＝**10/36**：`#9 301`、`#10 338`、`#16 473`、`#17 533`、`#18 571`、`#19 581`、`#20 603`、`#22 743`、`#23 779`、`#24 795`。其余 **21 枚属 B3**（`e6dc79ed` 那 23 枚里的 21 枚）**只有步色、没有守卫读数**——两格⛔ 不许合并。
值得单独留名的一形：R296（**37166458550**）里 `lint` 先红在 `gofmt (gofumpt)`（failure），紧接的 `#2 106`、`#3 136` 当场变 `skipped`＝**"前置步骤红 ⇒ 无守卫的步被 skip"这台机器在这六发里被真实观察到咬过一次**；同一发 `#1 81` 直接 `failure`。

**乙档｜只出过 `skipped`＝加了但从未生效过**＝**3/36**
`#6 lint:239`（gofmt 分母，**B2**）、`#7 lint:243`（`go vet (module)`，**B3**）、`#8 lint:247`（`go vet (tools/d22scan module)`，**B3**）——六发**全 skipped**，无一例外。
机制的步内证据同场可见：`#5 lint:173 gofmt (gofumpt)` 在 6 发里 **5 发 failure**（R296 那发它是 skipped），skip 就压在它下面那三枚上。⇒ **乙档三枚是"这道门至今一次没拦过任何东西"的确色**。
⚠ **与 `A710 §3` 的一处代际差（具名顶回，⛔ 不裁）**：台账把这三枚一起当作"那 23 枚里的真 skip"并写"其余 20 枚"。盘上尺说：其中 `lint:239` 那枚守卫是 **`f6b79ab0`（10-07 09:23）** 加的，**不在 `e6dc79ed` 的 23 枚里**；按代拆开，**23 枚里 skipped-only 的只有 2 枚**（`243`、`247`），`e6dc79ed` 名下"出过色"的是 **21 枚**而非 20 枚。两说各自内部自洽，差的是"哪一枚算进 23"这一格，**该由带 ref 的代册判**，本腿只把两把尺的数摆出来。

**丙档｜CI 里取不到这发读数（待取数）**＝**2/36**，且两枚都是**结构性取不到**，⛔ 不是网络失败
| 步 | 为什么取不到 | 已做的重试／换写法 |
|---|---|---|
| `#11 lint:394` `Portable tests carrier self-test (ticket 111 AC#3 - GUARD D's positive control)` | 该**步**由 `1309757b`（10-06 14:03 +08）新增，晚于远端 HEAD `cc31526165`（10-06 02:58:35Z）⇒ GitHub 上不存在含它的 yaml | 采样 6 发逐发对齐＝全 `不存在`；另用两把独立查法验"零发"：`gh api .../workflows/ci.yml/runs?created=>=2026-10-08T00:00:00Z` 回空（rc=0）＋同端点正控 `=>=2026-10-07` 回 `COUNT=1`（rc=0）⇒ 过滤器非哑 |
| `#21 test-windows:654` `winlive compile gate (go vet -tags winlive, ticket 111 AC#11)` | 同上，该步由 `351e5a5e`（10-08 09:01 +08）新增 | 同上＋`git log -S 'winlive compile gate' -- .github/workflows/ci.yml` 全历史只命中 `351e5a5e` 一枚（rc=0）⇒ 该步名在某枚提交之后才存在于盘上 |

⛔ 按派单第 4 条：这两枚**不读成"从未跑过"**，正确说法＝**"这一步（连同它的守卫）自 `1309757b`／`351e5a5e` 起才在盘上，更早的任何 run 里没有对应步骤"**。本腿另有一把尺撑这句：`git diff --stat cc31526165 HEAD -- .github/workflows/ci.yml` ＝ **160 insertions(+)、0 deletions** ⇒ `cc31526165` 那一版里的任何步**都没有被改名／挪位／删除过**，凡 HEAD 上出现而 `cc31526165` 上不存在的步名，只能是**新增**。
另有 **1 枚"半丙"**：`#4 lint:169`（success 4/6，两发不存在）与 `#22 test-windows:743`（3/6）＝**步本身有年代**，早期 run 无该步，故"该步出过色"成立、但**分母只有 6 发里的 4／3 发**。⛔ 不与"零枚"混写。

**网络类"取不到"＝零枚。** 本腿共 **18** 发 `gh`（`repo view` 1／`run list` 1／`run view` 6／`api` GET 10），其中 **16 发 rc=0**、**2 发 rc=1 且是同一次 422**（`e6dc79ed` 短号一次、全号重试一次），它是**决定性负读数**，⛔ 不是取数失败。派单点名的"六发 EOF 那族"在本腿**未遇到**。

## 4. 顺带回的那一枚问号：票 111 `AC#11` 的 winlive 补编译门

**现量＝编排者此前那句"只在 HEAD 上、从未进过任何 run"成立，本腿复认，不推翻。** 但要把形状补全，因为账上的规划文字与盘上有两处不重合：

1. **今天盘上有这一步**：`6547fd30:.github/workflows/ci.yml:606` 是步项起手（`- name: "winlive compile gate (go vet -tags winlive, ticket 111 AC#11)"`），**它的步级 `if:` 在 `:654`**，条件逐字 ` ${{ !cancelled() }}`。
2. **它落在 `test-windows`，不在"托管档"这一格里含糊**：`6547fd30:.github/workflows/ci.yml:516-517` ＝ `test-windows:` / `runs-on: windows-latest` ⇒ **GitHub 托管**，与规划那句"进 `test-windows` 托管档"**同色**；同一把 `runs-on` 尺全档现量＝`lint:66 ubuntu-latest`、`test-core:403 ubuntu-latest`、`test-windows:517 windows-latest`、`slo-smoke:802 windows-latest`、`slo-full:863 [self-hosted, wisp-slo]`、`lint-frontend:937 ubuntu-latest`（`slo-full` 那一格按派单是既有结论，本腿只把 `runs-on` 抄下来，⛔ 不重证、⛔ 未为取数触发任何 run）。
3. **那"12 枚"今天不成立**：规划写"12 枚 winlive 补编译门"，盘上尺＝**winlive 名下只有 1 枚步**（上面那一步，`go vet -tags winlive`），**不是 12 枚**。若"12"指的是脚本内的 12 枚包／用例，那射程在 `scripts/**` 里而**不在 ci.yml 的步数上**——本腿⛔ 未去数脚本内的枚数（那是另一格），只把"步＝1 枚"这一枚报回来，供编排者判"12"那格是不是把脚本内分母写成了步数。
4. **它的守卫属于 B2，不属于 `A710` 那 23 枚**（`e6dc79ed` 名下 test-windows ＝ **0 枚**）。⇒ 若按 `A710` 的撤销口令「撤 A710 那 23 枚 if: 守卫」执行，**winlive 那枚 `:654` 不在被撤之列**（它来自 `351e5a5e`）。这一格对"口令射程"有影响，具名报回，⛔ 本腿不动任何文字。

## 5. 台账豁免名册的行号落点（现量，⛔ 不裁谁对谁错）

`A710 §4` 写："另有两枚 `Upload SLO report`（`ci.yml:229-233`）与该 job 的前 2 步（`:469`／`:521`）同按原文豁免"。同一把尺在五个 ref 上量那三处行号：

| ref | `:229-233` 落在什么上 | `:469` | `:521` |
|---|---|---|---|
| `0c9726f9`／`f6b79ab0`／`05992d05`（＝`e6dc79ed^`） | **正是那段豁免注释**：`:229` 首行逐字 `` # `Upload SLO report` steps are deliberately NOT given a guard even though `` … `:233` = `# on "a valid sample happened". Registered in the ticket instead.` | 一行 **indent 2 的注释** `# EVERY STEP BELOW THE SETUP ONES CARRIES \`if: ${{ !cancelled() }}\`` | 注释行 `# itself guarded so a red checkout/setup cannot make it inconclusive` |
| `cc31526165` | 无关内容（staticcheck 的注释中段） | `winsec-tests.sh` 的注释行 | `test-windows` 的注释行 |
| **`6547fd30`（HEAD）／`e6dc79ed`** | 漂到 **`:234-238`**（同一段注释；`:229-233` 现在是"WHY THIS STEP AND NOT THE GATES ABOVE IT"那半段） | **`if: ${{ !cancelled() }}` 本身**（test-core `Portable package tests (core scope…)`＝**B3 那 23 枚里的一枚**） | **`- uses: actions/checkout@v4`**（test-windows 的 setup 步＝**确实无守卫**，这一枚落点是对的） |

⇒ 三格分开写：**`A710` 的行号是按 `e6dc79ed` 之前的版式引的**，在 `e6dc79ed` 起（含 HEAD）整体漂了 5 行；`:521` 在 HEAD 仍落在一枚真·无守卫的 setup 步上；`:469` 在 HEAD 落在**一枚被 `A710` 判定"加了守卫"的步**上，与"setup 步豁免"这句**在 HEAD 上不自洽**。⚠ 按本项目定式（`A710 §5` 自己那条），**裸行号不带 ref 必错**——这不是指新错，是把 `A710` 的行号补上 ref 限定。

**盘上真实的豁免形状（不靠行号，靠尺）**：HEAD 无守卫的步 **16 枚**，逐枚＝12 枚 setup（`lint:68/70`、`test-core:407/409`、`test-windows:521/523`、`slo-smoke:806/808`、`slo-full:877/879`、`lint-frontend:942/944`，全为 `- uses:` 起手）＋ 2 枚 `Upload SLO report`（`slo-smoke:827`、`slo-full:890`）＋ `slo-full` 那两枚门（`:884 Build wisp.exe (deps cached on the runner)`、`:887 SLO full gate (six states + settle + leak)`）。⇒ **`slo-full` 整个 job 5 枚步零守卫**，与 `A710 §4` 那句"一枚没补"复认；且 HEAD 文件里就写着理由：`6547fd30:.github/workflows/ci.yml:234-238`（`scripts/slo-freshness.sh` 的 P1 钉子会因 `slo-full` job body 带 `if:` 而自红）。⛔ 本腿未读 `scripts/slo-freshness.sh`，那句是**盘上注释的转述**。

`slo-full` 那枚豁免 job 在 R305（**37703959747**）里的实际读数（供"若补守卫会不会被求值"这一格作底，⛔ 本腿不提该不该补）：job 结论 `success`，5 枚 yaml 步全 `success`（`Build wisp.exe (deps cached on the runner)`／`SLO full gate (six states + settle + leak)`／`Upload SLO report`），另有 GitHub 自加的 `Set up job`／`Post Run actions/setup-go@v5`／`Post Run actions/checkout@v4`／`Complete job` 4 枚合成步同色。⇒ **该 job 每次都跑到最后一步**；"豁免"是**有意留的欠账**，不是"跑不到"。

## 6. 我看不见的格（列全，⛔ 不外推）

1. **甲轴对 26 枚（B2＋B3）永远零读数**，直到有一发 run 用过含它们的 yaml；今天⛔ 无任何 run 用过 ⇒ 本程**无法**给出"这 23 枚里哪枚真拦过谁"，只有"哪一步出过什么色"。
2. **六发采样之外**：`gh run list --limit 40` 只回到 2026-09-24T13:11Z（40 发，全 `completed/failure`），⛔ 未翻 40 发以前的 run ⇒ "B1 那 10 枚在更早的 run 里是否也求值过／是否曾拦下过东西"未取。
3. **`timed out`／`cancelled` 两色在采样 6 发里零枚**；要判"这一族 ever timed out"需要更多发的步级回包，本腿只取 6 发（每发 JSON 14.3–14.9 KB，⛔ 未取日志正文，防 1.29 MB 那一族掐上下文）。
4. **R296/R293 两发的 yaml 更老**，`#4`（4/6）与 `#22`（3/6）的"该步不存在"未逐发核到具体是哪一枚 commit 引入（`#4 Tracked path-length budget (ticket 262)` 未做 pickaxe），只具名分母。
5. **矩阵/动态 job**：本 workflow 没有 `strategy.matrix`，但 GitHub 会把 `- uses:` 步在 API 里显成 `Run actions/checkout@v4`；本腿一律按 **yaml 里的 `name:`** 对齐，凡 yaml 无 `name:` 的 setup 步在 run 侧的对应名⛔ 未逐枚比（那 12 枚不在 36 枚名册里）。
6. **`slo-fresh.yml` 另一枚 workflow**（`gh api` 显示 `slo-fresh` 最新一发 **37736935814** @ 2026-10-08T06:19:05Z，`schedule/success`）：结构尺现量＝1 个 job、3 枚步、**步级 `if:` 0 枚**。它是否算 `A710` 射程（"整个 `slo-full` job 豁免"那句有没有把这条 workflow 一起豁免）**不在盘上**，归编排者；⛔ 本腿未对它做任何对齐。
7. **同名步跨 job**：`Cache third_party (deps.toml key)` 在 `test-windows:571` 与 `slo-smoke:813` 各一枚、`Upload SLO report` 两枚——本腿一律用 `(job, 步名)` 二元键，⛔ 未把二者压成一枚；若某发 run 里 `job.name` 与 yaml 键名不同形（今天同名，六发皆然），对齐会漂，这一格只是"若"。
8. **文件内注释里引的另一发号 `35591482293`**（`6547fd30:.github/workflows/ci.yml:523` "read off run 35591482293"）按待验断言处理并已核：**真实**（`ci` @ `440dd88765…`，2026-09-21T10:57:32Z，`completed/failure`，rc=0）；但**该发里那一步的色本腿未取**（只核了"这发存在"）。

## 7. 纪律自证（本腿射程）

⛔ 未改任何既有文件；⛔ 未动 `.github/workflows/ci.yml`（`git status --porcelain -- .github/workflows/ci.yml` 起手与收笔皆回空）；⛔ 未 push、未触发 run、未 `gh run rerun`／`gh workflow run`；⛔ 零 go 命令（`build`/`vet`/`test`/`list`/`env` 全 0 发）；读面＝`git show`/`git log`/`git diff`/`git ls-tree`/`gh run list`/`gh run view`/`gh api`（GET 端点）。工作树里别人的东西（`.gitignore` 那节、`design/**` 的删改、`probes/**` 被重写的读数件、`cmd/wisp` 在飞改动）一枚没碰、没 add、没 checkout/restore/clean/stash；`git status --porcelain` 起手即见 40+ 条既有改动，本腿只在自己那枚 pathspec 上落笔。本件⛔ 不出判据、⛔ 不提改法，只取名册与色。
