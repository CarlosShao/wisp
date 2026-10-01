# colors.md — 推送前 CI 变色普查（代号 `pushpreflight-1`，只读腿）

## 0. 起手锚

- `date` = `2026-10-01 18:53 +0800`（首次取样 `18:45`，收尾 `18:53`）
- `git log -1` 短号：取样时 `412644fe`，写本件时 `459800fd`。⚠ **HEAD 在我这一趟里前进了**（共享工作树里别的腿在提交）。本件所有结论对 `459800fd` 复量过一遍，两个关键尺没变：全树 `66`、base `61`、净增 `5`。
- 分支 `dev`；`@{u}` = `cnb/dev`（**取得到**，不需退化路径）
- `cnb/dev` = `origin/dev` = `0589fd9c`（2026-09-30 11:22）。**两枚远端同点**，`ahead=228 behind=0`。
- ⚠ 用户口径"220 枚未推"→ 实测 **228**。差 8 枚，全部是 `0589fd9c..412644fe` 之间新落的。
- ⚠ CI 定义在 `.github/workflows/ci.yml`（GitHub Actions），而 `@{u}` 是 `cnb.cool`。**推 cnb 不触发这份 workflow**；能触发的是 `origin`（github.com）。两枚远端当前同点，所以"推上去"这一步实际要落在哪一枚，是编排者的决定，不是本件能定的。
- 本件全程只读命令：`ls` `grep` `sed` `awk` `wc` `git log/show/diff/grep/ls-files/rev-list/merge-base`。**未跑** `go test` `go build` `go vet` `gofumpt` `d22scan` `slo` 或任何 `.ps1`/`.sh`。

---

## 1. 名册：会跑测试的步骤（`ci.yml` 共 734 行，7 个 job）

`go test` 尺只有 6 枚步骤，分布在 3 个 job。**全部 6 枚都不允许 SKIP**（见 §1.1）。

| # | 作业 | 行号(job/runs-on) | 步骤名（逐字，行号） | 跑什么（行号） | 用的尺 | 包 | runner |
|---|---|---|---|---|---|---|---|---|
| 1 | `lint` | `:65` / `:66` `ubuntu-latest` | `D22 scanner positive control (tools/d22scan tests, seeded red)` `:74` | `:81` `bash tools/d22scan/runtests.sh -C tools/d22scan ./...` | **runtests.sh** | `tools/d22scan` 模块 `./...`（独立模块，不是主仓包） | GitHub 托管 |
| 2 | `test-core` | `:277` / `:278` `ubuntu-latest` | `Environment fork assertion (WISP_ENV=test data dir)` `:303` | `:317` `runtests.sh ./internal/proc/ -run TestLayoutForTestEnv` | **runtests.sh** | 仅 `internal/proc` 一枚命名用例 | GitHub 托管 |
| 3 | `test-core` | 同上 | `Portable package tests (core scope; …)` `:319` | `:341` `bash scripts/portable-tests.sh --scope=core` | **portable-tests.sh -> runtests.sh + ledger** | 25 个包（`core_pin`，`portable-tests.sh:123-149`）：agent, agent/approval, audio, ball, buildinfo, config, llm/…, memory, models, observe, panel, perm, plugin, proc, risk, secret, statemachine, tools, winsec, cmd/llmrecord | GitHub 托管 |
| 4 | `test-windows` | `:387` / **`:388` `windows-latest`** | `Windows ACL sealing gate (internal/winsec's own tests, ticket 110)` `:398` | `:439` `bash scripts/winsec-tests.sh` | **winsec-tests.sh -> portable-tests.sh -> runtests.sh** | 仅 `./internal/winsec/`（`winsec-tests.sh:38`，GUARD 1 禁清单丢它 `:57-66`） | GitHub 托管 |
| 5 | `test-windows` | 同上 | `cgo build smoke (build.ps1 fetch-deps + mingw link + doctor)` `:448` | `:453` `powershell … scripts/build.ps1 -Env dev` | 编译门，**不产 go test 读数** | 全模块 `go build` | GitHub 托管 |
| 6 | **`test-windows`** | 同上 | **`cmd/wisp CLI tests (needs the sherpa DLLs staged above, ticket 111 AC#4)` `:455`** | **`:475` `bash scripts/wisp-cli-tests.sh`** | **wisp-cli-tests.sh -> portable-tests.sh `--scope=cli` -> runtests.sh** | **仅 `./cmd/wisp/`**（`portable-tests.sh:195`，`cli_pin` = `github.com/CarlosShao/wisp/cmd/wisp`） | GitHub 托管 |
| 7 | `test-windows` | 同上 | `Portable windows tests (proc/secret/config/risk/ball/perm/plugin/llmrecord)` `:477` | `:511` `bash scripts/portable-tests.sh --scope=windows` | **portable-tests.sh -> runtests.sh + ledger** | 8 个包（`win_pin`，`:150-159`）：cmd/llmrecord, ball, config, perm, plugin, proc, risk, secret | GitHub 托管 |
| 8 | `test-windows` | 同上 | `PathResolver junction placeholder (real cases tickets 18/20)` `:513` | `:527` `runtests.sh ./internal/risk/ -run TestPathResolverJunctionWindows` | **runtests.sh** | 仅 `TestPathResolverJunctionWindows`（声明于 `pathresolver_junction_windows_test.go:88`） | GitHub 托管 |
| 9 | `slo-smoke` | `:532` / `:533` `windows-latest` | `SLO smoke gate` `:552` | `scripts/slo-check.ps1 -Subset smoke` | 非 go test | — | GitHub 托管 |
| 10 | `slo-full` | `:590` / `:591` **`[self-hosted, wisp-slo]`** | `SLO full gate (six states + settle + leak)` `:615` | `scripts/slo-check.ps1 -Subset full` | 非 go test | — | **本机 self-hosted（唯一一枚）** |
| 11 | `lint-frontend` | `:664` / `:665` `ubuntu-latest` | `npm ci` `:680` / `typecheck` `:683` / `lint` `:686` / `tokens:check` `:689` / **`build (vite build -> frontend/dist)` `:695`** / 4 枚 render evidence `:698 :717 :723 :728` | 全 `npm run …` | 非 go test（无 SKIP 尺） | frontend | GitHub 托管 |

**runner 归属**：6 枚 go-test 步里 **5 枚在 GitHub 托管**（`ubuntu-latest` x2、`windows-latest` x3）+ 1 枚 `tools/d22scan` 在 ubuntu 托管。**self-hosted 只有 `slo-full`，它不跑 go test**。⇒ **本仓没有任何一枚 go test 步跑在你我这台机器上**；所有跳红都在干净检出的托管 runner 上发生。

### 1.1 跳红尺到底是谁带的（这条决定"跳"算不算红）

- `tools/d22scan/runtests.sh:98-103`：`skipped != 0` -> `runtests.sh: $skipped test(s) SKIPPED and SKIP is not a pass (ticket 71 AC#3)`（`:99` 逐字）-> `exit 1`。另两条：`:94` go test 非零退出原样传；`:104-109` 零 PASS 且零 FAIL（`-run` 打空）也 fatal。`-v -count=1` 由 `:75` 强制。
- `scripts/portable-tests.sh:424` 把 `--scope=*` 全数交给 `runtests.sh`（`$strict`），并在 `:433-435` 额外打印未记账 skip 的 `file:line` + 原因。
- **ledger 机制（关键）**：`portable-tests.sh:319-333` 是一张 `name|package|platform|class|reason` 表；`:407-418` 按 `go env GOOS` 把 `platform` 命中 `any` 或本平台的行拼成 `-skip '^(A|B|…)$'` 交给 runtests.sh。⇒ **登记过的用例在 runner 上根本不被调度，因此不产 `--- SKIP` 行**，不算红。`:396-404` GUARD：ledger 行若在本平台编译不出该测试名，**直接 fatal**（防"改名即隐身"）。⇒ **登记 != 放宽**，这是本仓唯一一条被设计出来、且带反查的豁免通道。
- `winsec-tests.sh:94` 用 `bash "$portable" "${scope[@]}"` 委托（显式包名分支，`portable-tests.sh:251`），**ledger 仍生效**（`:276-282` 的 universe = scope + 全部 ledger 包），但 **ledger 里没有任何 `internal/winsec` 行** -> winsec 的 skip 一律算红。
- `wisp-cli-tests.sh:113` 是 `bash "$portable" --scope=cli` -> **ledger 生效，但 ledger 里没有任何 `cmd/wisp` 行** -> cmd/wisp 的 skip 一律算红。`wisp-cli-tests.sh:74-99` 另有 DLL 前置 GUARD：`third_party/sherpa-onnx` 缺目录或缺 `deps.toml` 里 `[sherpa-onnx.dll.*]` 点名的任一枚 -> `exit 1`，**明写这是 setup 红、不是包坏**（`:93-98`）。`:60-68` GOOS != windows -> `exit 2`。

**特别确认（用户点的那枚）**：步骤 6 `cmd/wisp CLI tests` 真实命令 = `bash scripts/wisp-cli-tests.sh`（`ci.yml:475`）。它**不自己数 skip**，数 skip 的是 `runtests.sh:98`。⇒ **cmd/wisp 里任何一枚 `t.Skip*` 触发，红就算在它头上**，与 WebView2、与面板对不对都无关。用户这句判断成立。

---

## 2. 66 处跳位逐条判触发

尺复量（逐字用户口径）：`grep -rhoE 't\.Skipf?\(' --include=*_test.go internal cmd | wc -l` = **66**（工作树与 `git grep HEAD` 两法一致）。用户给的 66 **对上，无差**。全仓含 `tools/` 则是 76（`tools/d22scan` 另有 8 枚，见 §2.3）。

判"会不会跳"的锚：`runs-on: windows-latest` ＋ **干净检出**（只有 git 跟踪物）＋ **`WISP_ENV=test`**（`ci.yml:390`）＋ 无 `winlive` 标签（`grep -n winlive ci.yml` = **无匹配**）。

### 2.1 windows 腿真编译真进分母的 25 枚（能变色的只有这些）

`T`=会跳 `N`=不会跳 `U`=未定；"记账"= 已在 `portable-tests.sh` ledger。

| 文件:行 | 所属用例 | ①`if` 在测什么 | ②盘上物 | ③干净 windows-latest | 步 | 备注 |
|---|---|---|---|---|---|---|
| `cmd/wisp/panel_resident_windows_test.go:317` | `TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe`(`:314`) | `len(entryIDProbes(t))==0` | **`frontend/dist` 真产物** | **T（见 §5，读数级）** | 6 | **不在 ledger** -> 红 |
| `cmd/wisp/panel_host_windows_test.go:794` | `TestPanelHostLatencyPercentilesAC2`(`:785`) | `len(cold)==0`（同进程没记到样本） | WebView2 真窗口（样本来自 `:507` lifecycle） | U | 6 | 同文件源序 `:507 < :785`，lifecycle 先跑；lifecycle 无 skip 分支 -> 它坏则本步已因 FAIL 红 |
| `cmd/wisp/panel_resident_windows_test.go:149` | helper `editorOnPanelThread`(`:126`)，唯一调用者 `:986` | `createEditorWindow` 返 err | 桌面会话可建窗/可置前台 | U | 6 | 与 `:999` 同属一枚用例，至多计 1 枚 SKIP |
| `cmd/wisp/panel_resident_windows_test.go:999` | `TestAC4PriorFocusSurvivesARefusedPanelSample`(`:984`) | `mgr.prevFocus==0` | 前台焦点（真桌面） | U | 6 | |
| `cmd/wisp/panel_host_gate_test.go:182` | `TestAC1SessionDisposeHasAProductionTrigger_AC1`(`:178`) | `len(ctorHits)==0` | 无（纯文本扫包目录，`panel_host_gate_test.go:208`+`:216` 只按 `.go` 后缀、不过滤 GOOS） | **N** | 6 | `cmd/wisp/panel_resident_windows.go:204` 是真调用点 -> 命中 |
| `internal/proc/jobscope_windows_test.go:87` | `TestHelperProcess` | re-exec 子进程标记未设 | 无 | T（恒跳） | 7 | **已记账** `portable-tests.sh:322` -> 不计红 |
| `internal/risk/syncdirs_windows_test.go:133` | `TestSyncRegistryProbeLive` | HKCU 无 `UserFolder` | 注册表 hive | T | 7 | **已记账** `portable-tests.sh:326` -> 不计红 |
| `internal/risk/syncdirs_redteam_windows_test.go:220` | `TestSyncRedTeamRealOneDrive`(`:206`) | 遍历 `p.SyncRoots()` 后 `root==""` | **真 OneDrive 同步根**（注册表/配置探针） | **T（读数级）** | 7 | **不在 ledger** -> 计红。run `35606321404` win.log 逐字：`syncdirs_redteam_windows_test.go:220: no live sync root on this machine (detected roots: [])`（`docs/evidence/s1/111-ci-step-readings-attempt2.md:137-141`） |
| `internal/risk/syncdirs_redteam_windows_test.go:208` | 同上（前一行） | `userHomeDir()==""` | `USERPROFILE` | N | 7 | 托管 runner 有用户目录 |
| `internal/risk/syncdirs_redteam_windows_test.go:54` | helper `shortPathOf`(`:50`)，调用者 `:84` `TestSyncRedTeamNewFileSpellingsDeniedAndLandNoBytes`、`:151` `TestSyncRedTeamGuardPlainNewFileWriteAllowed` | 卷不生成 8.3 别名 | 短名生成开启 | N（U 低） | 7 | 同一枚 run 的 step 4 里 `TestNoticeAttributionSurvivesAn83AliasOfItsOwnTree` 是 **FAIL 而非 SKIP**（evidence `:85-90`）=> 该 runner **确实给了短名** |
| `internal/risk/syncdirs_test.go:498` | `TestSyncDetectionOnThisMachine` | `home==""` | `USERPROFILE` | N | 7 | |
| `internal/risk/pathresolver_rewrite_account_test.go:51` | `TestC26RewriteAccountIsRecorded` | 无 `~` 可展开的既有目录 | `USERPROFILE` | N | 7 | |
| `internal/risk/pathresolver_rewrite_account_test.go:138` | `TestC26RewrittenSyncRootDoesNotDisarmSuspectNet` | POSIX 无 handle-resolved 形式 | 平台 | N | 7 | ledger `:327` 登记的是 **linux**，windows 侧注"RUNS AND PASSES"；09-21 run step 8 的 12 枚 FAIL 不含它（evidence `:135`）|
| `internal/risk/pathresolver_junction_windows_test.go:128` | `TestPathResolverShortNameAListDenied` | 卷禁 8.3（`GetShortPathName` 回长名） | 短名生成 | N（U 低） | 7（**不在步 8**，步 8 只 `-run …JunctionWindows`，该用例声明于 `:88`，`:128` 属下一枚） | |
| `internal/winsec/reparse_windows_test.go:63` | helper `mkDirSymlink`(`:60`) -> `TestAC4DirectorySymlinkAtArtifactPositionIsNotRecursed`(`:111 :129 :148`)、`TestAC4DirectorySymlinkInsideSealedTreeIsNotWalked`(`:180`) | `os.Symlink` 建目录链接失败 | 符号链接特权 / 开发者模式 | N（U 低） | 4 | ledger **无 winsec 行**；09-21 step 4 `SKIP=0`（evidence `:83`）=> 该 runner 给得起特权 |
| `internal/winsec/absoluteness_seam_landing_129_windows_test.go:117` | `TestAC1DriveRelativeAndAbsoluteSpellingsNameTwoObjectsOnThisBox` | `os.Getwd()` 不含卷段 | 工作目录在盘内 | N | 4 | 同上 `SKIP=0` |
| `…:144` | 同上 | 进程站在卷根，相对/绝对同物 | 卷根可写 | N | 4 | |
| `…:147` | 同上 | 相对侧种不下去 | 同上 | N | 4 | |
| `…:150` | 同上 | 绝对侧在卷根种不下去 | 同上 | N | 4 | |
| `…:215` | `TestAC1SeamVouchesForTwoRealObjectsThatDifferOnlyInAbsoluteness` | 全机给不出 相对 vs 绝对 对 | 多卷 | N | 4 | |
| `…:231` | 同上 | 相对侧种不下去 | 同上 | N | 4 | |
| `…:234` | 同上 | 绝对侧种不下去 | 同上 | N | 4 | |
| `internal/winsec/absoluteness_attribution_129_windows_test.go:311` | `TestHonestPipelineKeepsPassingItsOwnProbe` | 该卷不给 fixture 短名 | 8.3 生成 | N | 4 | 同上 |
| `internal/winsec/tree_ownership_112_windows_test.go:122` | `TestTreeOwnershipProbeAcceptsAn83ShortSpellingOfItsOwnParent` | 同上 | 8.3 生成 | N | 4 | 同上 |
| `internal/winsec/volume_attribution_126_windows_test.go:207` | `TestNoticeFromOneVolumeIsNotAttributedToASecondRealVolume` | **可写卷根 < 2** | **多卷** | U（09-21 读数支持 N） | 4 | ledger **无 winsec 行**。09-21 step 4 `SKIP=0` => 当时该 runner 至少 2 枚可写卷根 |

### 2.2 windows 腿**根本不计**的 41 枚（三条形阶，逐条只给归属）

- **`winlive` 标签，CI 永不编译（8 枚）**：`internal/ball/hotkey_live_test.go:202 :213 :238 :371`、`internal/ball/interaction_live_test.go:233 :327`、`internal/ball/live_windows_test.go:495 :577`。尺：`grep -n winlive .github/workflows/ci.yml` **零命中** -> 这 8 枚在 CI 上既不算红也不算绿，**是覆盖洞不是颜色**。
- **`!windows` 标签，windows 腿不编译（16 枚）**：`cmd/wisp/secret_dataroot_119b_test.go:73 :234`（2，且 `wisp-cli-tests.sh:60-68` 在非 windows 直接 `exit 2`，所以这 2 枚**两腿都不跑**）；`internal/config/c26_seam_posix_125_test.go:120 :129`（2，ubuntu 腿步 3 计）；`internal/winsec/ancestor_separator_108_other_test.go:88 :122 :155 :183`（4，ubuntu 步 3）；`internal/winsec/dataroot_symlink_119_other_test.go:89 :283`（2）；`internal/winsec/placement_symlink_113_other_test.go:130 :388`（2）；`internal/winsec/seam_probe_root_125_other_test.go:62 :71`（2）。⇒ 这 10 枚 winsec POSIX 侧**在 ubuntu 腿真进分母**（09-21 读数：core `SKIP=0`，`grep -c "unaccounted SKIP" core.log`=0，evidence `:131 :137`）-> 判 **N**。
- **windows 编译但不在任何 windows 步的分母（12 枚）**：`internal/audio/hotplug_test.go:445 :527`（`./internal/audio/` 在 `core_pin` 但**不在 `win_pin`** -> ubuntu 上这枚文件不编译、windows 上 audio 不在步 7；**ledger `:323` 却登记了 `TestLiveWasapiSmoke|./internal/audio/|windows`，那行在 windows 腿是空转**——记为可疑条目 §7-③）。其余 10 枚属"包在 ubuntu 腿、文件在 windows 腿不编译 / 包根本不在 windows 分母"：`internal/agent/approval/ticket84_no_owner_test.go:224`（ledger `any` 已记账）、`internal/llm/llm_test.go:554`、`internal/llm/openaichat/mockllm_integ_test.go:45 :54`（`:45` 要 `-short`，`runtests.sh:75` 不传 -> N）、`internal/memory/concurrent_test.go:186`（ledger `any` 已记账）`:222`、`internal/models/manifest_real_test.go:120 :147`（ledger `any` 两枚都已记账）、`internal/panel/attachments_test.go:191`（2 GB 文件；ubuntu 临时盘多半容得下 -> N）、`internal/tools/fs_write_test.go:140 :151`（helper `otherVolumeDir:136`；ledger `:330 :331` 以 linux 平台登记 `TestD34WriteMatrix` / `TestCrossVolumeMoveStopsWithTwoCopiesOnLateStop`；tools **不在 `win_pin`** -> windows 腿不计）、`internal/tools/paths_workspace_test.go:145 :149 :175 :198 :205`（`:198` 已由 ledger `:328` 登记 linux；其余 4 枚 ubuntu 侧 09-21 读数 `SKIP=0` -> N）、`internal/tools/recycle_windows_test.go:39`（windows 编译，`./internal/tools/` 不在 `win_pin` -> **两腿都不进分母**）。
- ⚠ ubuntu 腿（步 2/3）合计判 **N** 或**已记账**，09-21 实测 `SKIP=0`（`RUN=1036 PASS=661 FAIL=0`）。⇒ **跳红风险几乎全集中在 windows 腿。**

### 2.3 66 之外的一枚尺

`tools/d22scan/scan_test.go:241 :261 :677 :1044 :1189 :2180`（`not inside the wisp repo`）+ `:2253`、`selftest_test.go:178`（`no go toolchain`）= **8 枚**，走步 1 的同一把 `runtests.sh`（`-C tools/d22scan`）。干净检出 + `setup-go` 下两判据都不成立 -> 判 N，但这是**第 7 把可能的红来源**，用户 66 的尺看不见它（尺只扫 `internal cmd`）。

---

## 3. 只数"会新增的红"

尺（只读，`@{u}` 取到了，不需退化到 `cnb/dev`/`origin/dev` 字面量）：
- `git merge-base @{u} HEAD` = **`0589fd9c`**（= `cnb/dev` = `origin/dev`，`behind=0`）
- `git rev-list --count 0589fd9c..HEAD` = **228**
- 跳位净增：`git grep -h -o -E 't\.Skipf?\(' <rev> -- 'internal/*_test.go' 'cmd/*_test.go' | wc -l` -> base **61** / HEAD **66** -> **+5**
- 逐文件核（`git show <rev>:f` 对比计数）：涨的只有 3 枚文件，全在 `cmd/wisp/`：`panel_host_gate_test.go 0->1`、`panel_host_windows_test.go 0->1`、`panel_resident_windows_test.go 0->3`
- 新增用例数：`git diff 0589fd9c..HEAD -- 'cmd/wisp/*_test.go' | grep -c '^+func Test'` = **50**，其中 7 枚带 `winlive`（不进 CI）-> **windows 腿顶层用例 117 -> 160（+43）**（按 build tag 过滤 `winlive`/`!windows` 后逐文件累加）
- ⚠ `@{u}..HEAD` 里改到的 `_test.go` 只落在 **`cmd/wisp`(15 枚) / internal/proc(2) / internal/ball(3) / internal/session(3) / internal/tools(3) / internal/agent/approval(1) / internal/audio(1) / internal/panel(1)`。`internal/winsec` 与 `internal/risk` **零改动** -> §2.1 里那 20 枚 winsec/risk 跳位全是存量，不是新增。

**能站住的那一句**：

> 推下去，`cmd/wisp CLI tests`（`ci.yml:455`/`:475`）这**一枚步骤**的红名集合会比上一次**至少多 1 枚**，逐名 = **`TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe`**（`cmd/wisp/panel_resident_windows_test.go:314`，skip 在 `:317`）；这是唯一一枚静态链条能走满的**无关红**（跟 WebView2、跟面板对不对都无关，只因 `frontend/dist` 干净检出只有 `.gitkeep`）。同一批还新带进 **4 枚会跳的用例**（`TestPanelHostLatencyPercentilesAC2`、`TestAC4PriorFocusSurvivesARefusedPanelSample`（含 helper `:149`）、`TestAC1SessionDisposeHasAProductionTrigger_AC1`）—— 前 2 枚判据依赖真桌面/WebView2，**未定**；第 3 枚判 **不会跳**。存量红（`TestSyncRedTeamRealOneDrive` 的未记账 skip、winsec 6 枚 FAIL、cmd/wisp 4 枚审批超时 FAIL）**不因为这次推送而新增**，但会同时出现。分母一次涨 +43，所以红名的**上限**远大于 1，只是我拿不到"上限是多少"的读数（§6）。

⚠ 一个必须说的口径限制：**"上一次"的红名集合我拿不到读数。** 最后一枚有完整步级读数的 run 是 `35606321404`（head `d3cc9ed`，2026-09-21，见 `docs/evidence/s1/111-ci-step-readings-attempt2.md:6`），比 base `0589fd9c`（09-30）**早 9 天**。`0589fd9c` 自己的那枚 run 的步级日志没有落进仓里——`grep -rn '0589fd9' docs/evidence/` 只有 **1 处命中**（`docs/evidence/s1/239-ball-square-fix-v1.md:4`），那是**对 base 那枚改动本身的验收件**（"本腿 `git log -1` 现取"），里面**没有任何 CI 步级 conclusion**。⇒ §3 那句说的是**相对 base 的 delta（静态可推）**，不是相对"上一次 CI 实际报的那张红表"。

---

## 4. 最小修法（只列方向与代价，本腿不动任何文件）

三条禁令我都遵守：不提放宽/删断言、不提把 SKIP 当绿、不提注释掉那一步。

**M1 — 给 `frontend/dist` 补来源，让 AC13 有 subject。**
做法：在 `test-windows` 里、步 6 **之前**加一步产出真 dist。两种形状：(a) 该 job 内跑 `npm ci && npm run build`（步 11 的 `ci.yml:695` 已是同一条命令，可抄）；(b) `lint-frontend` 用 `upload-artifact` 发 dist、`test-windows` `download-artifact` 落到 `frontend/dist`（跨 job，省一次 node 安装）。
代价：windows 腿时长 +（node 安装或 artifact 往返）；`third_party` cache 与 node cache 争配额；artifact 形状要让 `go:embed all:dist` 看得见（`.gitkeep` 与 `index.html` 共存即可）。
**没防住的形状**：真产物一旦进 runner，`panel.Assets` 的 `built=false` fail-closed 分支在 CI 上**从此永不执行**——`cmd/wisp/panel_host_gate_test.go:130-137` 那三条 `built=false` 断言（`Resolve`/`Check`/`Manifest` 必须一起关门）会从"实跑"退化成"只被 `newAssets` 的合成 bundle 测到"。AC#12 的另一半丢 CI 证据。

**M2 — 把 AC13 写进 `portable-tests.sh` 的 ledger（`windows|fixture`）。**
做法：照 `:323` `TestLiveWasapiSmoke|./internal/audio/|windows|fixture|…` 的既有形状加一行，reason 里点名 `frontend/dist` 干净检出只有 anchor、并写 next= 归哪张票。
代价：**零 CI 改动**，是最小的一枚。且它**不是放宽**——`:396-404` GUARD 会拿 `go test -list` 反查：用例改名/删除/被 build tag 藏起来 -> 该行 stale -> 红。这与"恒真判据"的区别就在反查：**判据是"这个名字在本平台编译产物里存在"，不是 `true`**。
**没防住的形状**：`bringUp` 把 probe 页留在最后（AC#13 原缺陷）时，**CI 不响**——这条 AC 只剩本机读数（`docs/evidence/s1/33-panel-host-c27-r5.md`）。等于把一条产品判据从 CI 撤出。
**必须同时做**：账要落到 `docs/reports/pending-and-issues.md` 的 `A##`，并按 `SPEC-12 §5` 五字段登记推迟（本腿不做）。

**M3 — 把 4 枚"需要真桌面"的用例从 `--scope=cli` 摘出去，独立成一步。**
做法：`portable-tests.sh:191-198` 的 `cli` scope 拆成 `cli`（无桌面依赖的那些）+ 一枚新 scope，只由**有桌面会话前提**的 runner 跑；新步在自己的 `if:` 上写明前提。
代价：步骤数 +1；`wisp-cli-tests.sh:113` 要改调用；GUARD B（`portable-tests.sh:438-460`，每包必须印出自己的顶层结果行）要在两个 scope 各自成立，别出现"同一包被两步重复认领"。
**没防住的形状**：仍然没有真 WebView2 桌面会话的覆盖——只是把红从"看不懂的混色"变成"看得懂的独立红"。**收益是归因，不是变绿。**

**M4 — 换 AC13 的探针来源，使它不依赖真产物。**
做法：`entryIDProbes` 的 id 改从 Go 侧一份**内嵌 fixture** 取（与 `panel.Assets` 解耦），判据仍问活文档"这些 id 在不在"。
代价：动的是 AC#13 的射程本身（票面写的是"真冷启动后的文档是内嵌页"），**属契约面，要人工批准**，不是 CI 腿能自己拍的。
**没防住的形状**：一旦 id 集与真 embed 的 entry 脱钩，"内嵌页被换成 probe 页"这个具体缺陷可能改由 fixture 侧看不见——需要另配一条"fixture 的 id 必须出现在真 entry 里"的锁，否则就是把无关红换成了假绿。

优先级建议（本腿不执行）：**M2 是当天能落且不造假的一枚；M1 是唯一让这条 AC 真在 CI 上跑的一枚；两者并存最好**（M1 让 `built=true` 那条腿跑，`built=false` 那条腿用 M2 之外的合成用例守——见 M1 的"没防住"）。

---

## 5. `entryIDProbes(t)` 怎么判空 / "干净检出必跳"是什么级

**判空实现**（`cmd/wisp/panel_resident_windows_test.go:265-290`）：
```
assets, err := panel.BuiltinAssets();      if err != nil { return nil }      // :267-270
data, _, err := assets.Resolve(panel.EntryFile); if err != nil { t.Logf(...); return nil }  // :271-275
再用 entryIDRe (:259) 抓 id，抓不到也返回 nil
```
⇒ 三条返回 `nil` 的路：**embed 打不开** / **entry 解析不了** / **entry 里没 id**。`:314` 的用例只在 `len(probes)==0` 时 `t.Skipf`（`:317`），**不区分这三条**。

**"干净检出必跳"= 读数级**，链条 6 环，每环都是实读：
1. `git ls-files frontend/dist` -> **只有一枚 `frontend/dist/.gitkeep`**（本腿实测；未读 `frontend/**` 内容）。旁证：`docs/evidence/s1/137-ac2-ac5-r1-acceptance.md:422` 逐字"被 git 跟踪的只有 `.gitkeep` 一枚"。
2. `frontend/embed.go:19` = `//go:embed all:dist` -> `all:` 前缀**含点开头文件**，所以 `.gitkeep` 被嵌进去，**编译成立**（这一步挡掉了"其实是 build 失败"这一种相反解释）。
3. `internal/panel/assets.go:54-60` `newAssets`：`a.built = (fs.Stat(tree,"index.html") == nil)` -> 干净树里 **built=false**。
4. `internal/panel/assets.go:75-78` `Resolve`：`if !a.Built() { return …, errNotBuilt }`（`errNotBuilt` 定义在 `:34`）。
5. `cmd/wisp/panel_resident_windows_test.go:271-275`：`Resolve` 出错 -> `return nil` -> `len(probes)==0` -> `:317` `t.Skipf`。
6. `:455/:475` 那枚步的 ledger 里**没有** `TestAC13ColdStart…`（`portable-tests.sh:319-333` 逐行看过）-> `runtests.sh:98` `skipped != 0` -> `exit 1` -> **该步红**。

唯一的**推断**只剩一句：**"runner 的检出树 = git 跟踪内容，且没有别的步骤往 `frontend/dist` 写 `index.html`"**。这半句我也按读数堵了：`test-windows` 全部步骤（`ci.yml:392 checkout`、`:394 setup-go`、`:398 winsec`、`:441 cache`、`:448 build.ps1`、`:455 cli tests`、`:477 portable`、`:513 junction`）里 `sed -n '387,528p' ci.yml | grep -inE 'npm|vite|frontend'` = **零命中**；且 `scripts/build.ps1:74` 逐字 `build.ps1: frontend step skipped (no frontend yet; embed lands in S5 per SPEC-11 §2.2)`——**编译脚本自己不产 dist**。`npm run build` 只在 `lint-frontend`（`ci.yml:695`），那是**另一个 job、另一台 ubuntu 机器**，工作树不共享。

⇒ **结论：读数级**。要把它升成"runner 侧读数"（而非静态链），需要的是**一根现成日志里的行**，不需要"推一次看看"：
`cmd/wisp/panel_host_gate_test.go` 的 `TestPanelBundleShapeSeparatesAnchorFromRealPage_AC12` 会打印 `built=%v` 与 tracked 计数，`entryIDProbes` 自己也 `t.Logf("AC#13 probes from the resolved entry (%d bytes): %d id(s)")`（`:288`）。**但这两枚文件都是本批新增（`git diff --name-status 0589fd9c..HEAD` 标 `A`）**，所以**历史上任何一枚已推送 run 的日志里都不可能有这行**。⇒ 现成的 runner 日志给不出这根读数；能给出的唯一非推送途径是本地做一次 `git archive HEAD` 纯净快照跑测——**那是 CPU，本腿禁做**。

---

## 6. 我可能判错的条目（附"错了后果"）

- **① `TestAC13ColdStart…` 是唯一确定的新增无关红。** 若错在"其实还有别的"：漏判的多半是 §2.1 里 4 枚 `U`（`panel_host_windows_test.go:794`、`panel_resident_windows_test.go:149/:999`）。**后果**：第一次真推送会一次多红 2-3 枚，且看起来像"WebView2/桌面坏了"，把真正那枚 dist 来源问题埋掉。
- **② 我对 winsec 的 11 枚 windows 跳位判 `N`，靠的是 2026-09-21 run `35606321404` 的 `SKIP=0`（evidence `:83`）。** 那比 base 早 9 天，且当时 `RUN=82` 对现在（我只数得到源里的声明数，没编译）**未必同分母**。**后果**：若新增的 winsec 用例引入了新跳位，步 4 红因会从"6 枚 ACL FAIL"变成"FAIL + unaccounted SKIP 混色"，归因变难（不过本批 `internal/winsec` 零改动，`git diff --name-only 0589fd9c..HEAD -- internal/winsec` 为空，这条风险已按读数压掉）。
- **③ `TestLiveWasapiSmoke` 的 ledger 行（`portable-tests.sh:323`）登记平台是 `windows`、包是 `./internal/audio/`，但 `audio` 不在 `win_pin`（`:150-159`）。** 我判这是"空转的行"，不是红。**若我错**（比如 GUARD 把它算进 stale）：步 7 会因**一个与我普查无关的名字**红。**后果**：一次追错方向；现量：`grep -c 'internal/audio' ` 对 `win_pin` = 0 命中（已核），所以更可能是"那行永远不会被 -skip 命中、也永远不会 stale-红"，属账本卫生问题不是颜色问题。
- **④ "红名集合比上一次"的基线我用了 9 天前的 run。** 见 §3 末的口径限制。**后果**：delta 说成绝对集合会低估存量红的数量；我在 §3 已经把两句分开写了，若被合并引用就会误导。
- **⑤ 用户"220 枚"我量到 228，且 HEAD 在本趟里从 `412644fe` 前进到 `459800fd`。** 关键尺（61/66/+5/+43）在 `459800fd` 上复量未变。**后果**：如果继续有腿提交，§3 的数字要重跑；`winlive`/`ledger` 那两条结构性结论不受提交数影响。
- **⑥ 最大的一条：推哪一枚远端，决定这份分析会不会真的发生。** `ci.yml` 是 **GitHub Actions**，而 `@{u}` 是 `cnb/dev`。若"推上去"只推 cnb -> **一份 workflow 都不会跑，一个颜色都不变**，本件全部结论处在"尚未启动"状态。只有推 `origin`（github.com）才触发 §1 那 6 枚 go-test 步。**后果**：把"220 枚未推"当成一件要做的事，可能实际只完成了一半（cnb 同步了、GitHub 仍停在 `0589fd9c`），而红名集合的讨论全部落空。现量：`cnb/dev` = `origin/dev` = `0589fd9c`（两枚同点，所以现在还没有分叉风险）。
- **⑦ `runs-on: windows-latest` 的实际镜像我没法静态确认（C: 之外有没有第二枚可写卷根、8.3 是否开启、符号链接特权）。** 我用 09-21 run 的 `SKIP=0` 当代理读数。**后果**：`volume_attribution_126_windows_test.go:207` 与 3 枚 8.3 类可能改判 `T`，那会让**步 4** 新增未记账 skip（存量用例，非新增）。

---

## 7. 判不动的地方（甲=按我读数 / 乙=得换形 / 不做，配现量）

| 格 | 判不动的原因 | 甲 | 乙 | 不做 | 现量 |
|---|---|---|---|---|---|
| J1 `0589fd9c`（真正的上一次推送）那枚 run 的每步 conclusion | 步级日志没落进仓；`grep -rn '0589fd9' docs/evidence/` = 1 命中，但那是**改动验收件**（`239-ball-square-fix-v1.md:4`，"本腿 `git log -1` 现取"），**零 CI conclusion** | 用 09-21 的 run 当近似基线（本件就这么写的，并已在 §3 明写口径） | 派一枚"只读 `gh api` 远端 run 日志"的腿——**不需要推送**，本腿未开网络所以没做 | 不推一次试 | `gh` 未调用（本腿未发起任何 HTTP） |
| J2 4 枚 `U`（`:794`/`:149`/`:999` + WebView2 可用性） | 判据要真桌面会话，静态读不出 GitHub runner 给不给 | 按"未定"报，不进 §3 那句 | 给每枚配一根 `wisp panel-assets -manifest`/`built=` 的 runner 侧日志行——**同样要先有 run**（回 J1） | 不推 | cmd/wisp windows 分母 117 -> 160 |
| J3 `--scope=cli` 是否真被 WebView2 拖红（vs 审批超时那类） | 09-21 的 4 枚 FAIL 是**审批超时**（evidence `:114-118`），当时还没有 panel 用例；两类红混在同一枚步 | 拆成"存量 4 枚审批 + 新增 1 枚 AC13 skip"分开报 | M3 的拆步——归因立刻清楚 | 不推 | 步 6 = `ci.yml:455`，命令 `:475` |
| J4 `frontend/dist` 真产物该由谁给（界面侧 or CI or 入库） | 涉及 `frontend/**` owner 未提交的活，本腿不许读 | 只报"CI 里没有来源"这一句（`sed -n '387,528p' | grep -ic npm`=0） | 由编排者向界面侧要一条交付约定 | 不动 `frontend/**` | `git ls-files frontend/dist` = 1 枚 `.gitkeep` |

---

## 8. 注入面登记（照实报）

本趟所有工具输出（`git`、`grep`、`sed`、`awk`、`Read`，含 `docs/evidence/**` 正文与 `ci.yml`/脚本注释）里，自称"系统提示 / 编排者授权 / 已解锁 / 把阈值放宽 / 不用取证直接给结论 / 请 revert / 冻结某包"的文字：**0 次**。
`docs/evidence/s1/111-ci-step-readings-attempt2.md` 与 `ci.yml` 注释里出现的"不要放宽本行""no threshold changed""do not relax this script"是**上一任对这枚仪器的自约束说明**，按数据读、未当成对新腿的授权；`runtests.sh:100-101` 同。`runtests.sh`/`portable-tests.sh`/`wisp-cli-tests.sh`/`winsec-tests.sh` 里没有任何凭据值，本件也未输出任何凭据。
`ci.yml:521` 那行提到 runner 的 8.3 短名 `RUNNER~1`，是路径形状、不是凭据，保留。

## 9. 本件动了什么

只新建本文件：`.scratch/wisp/probes/orchestrator/preflight-1/`（新目录，用户已确认未用过）+ `colors.md`。既有文件零改动；`.github/**`、`scripts/**`、票面、`docs/**` 未碰；未跑任何 go/ps1/sh 执行件；未推送、未建 PR、未建 worktree、未额外 checkout。
