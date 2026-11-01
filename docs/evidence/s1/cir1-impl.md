# CI 修复腿 ci-r1 — 实现腿申报（2026-11-01）

写腿：ci-r1（实现者）。派单：`.scratch/wisp/probes/orch/20261101/txt/ci-r1-dispatch.txt`。
开工三核：HEAD `2d4ac8460`（dev，≥ 派单要求的 `2d4ac8460`）✓；索引干净 ✓；写面无他人脏文件 ✓。
日志源：`gh run view 41274267651 --repo CarlosShao/wisp --log-failed` 全量 1,555,678 字节拉回仓外 `D:\tmp\wisp\20261101\cir1\run-41274267651-logfailed.txt`。

---

## 0. 结论速览

| CI job | 派单时的定位 | 实测真因 | 本腿结果 |
|---|---|---|---|
| lint | staticcheck U1000 清单＋d22scan self-test「type unparseable」 | **三个失败步骤**：①path-length 预算门（roster 61 vs 树 62）②`go vet`＝aecspike 的 `aecsmoke.c` ③staticcheck 78 条。**self-test 是误报**（40/40 全过） | ①②③全部处置；staticcheck 仅剩 `frontend/embed.go:4` SA9009（本腿 ⛔ frontend/**，见 §6） |
| test-core | `internal/proc` TestLayoutForTestE exit 1 | **误判**。TestLayoutForTestEnv PASS=1 FAIL=0；真因＝`portable-tests.sh --scope=core` 的 `go list` 被 aecspike 拖死（与 lint 的 vet 同源） | aecspike 修复＋三包认领后 scope 可解析、GUARD A/C 本地过；**panel 四条既有红仍在**（§5），core 仍会红但已从「无法解析」变为「真红可见」 |
| test-windows | 待拉 | 五个失败步骤：winsec 342r1 竞态（真产品发现）、CLI 两红（autostart 已修／223 宿主环境）、winlive vet（过）、census GUARD D 四包、portable windows 一红（环境反常） | autostart 修复＋census 4→1（speech 结构项上报）；winsec/223/ReadFailure 为上报项（§5） |

---

## 1. 派单误判的三处更正（附证据）

1. **d22scan self-test 没有 exit 1。** CI 日志第 489/491 行的两条 `##[error]` 是 `actions/setup-go` 的 problem-matcher 把 self-test 输出里**内嵌的 Go 解析器错误文本**（`.../broken.go:5:2: expression in go must be function call`）注解成了错误行；self-test 自己的结论行是 `clean - all 40 direction checks passed (20 expect-ring, 20 expect-silent)`，随后步骤正常进入下一步。本地复跑（Git Bash，`go run . -self-test`）＝`clean - all 40, rc=0`。**停手条件②不触发**：`[unparseable]` 指向的 broken.go 本来就是 self-test 自己种的「解析必败」探针，不是种子缺陷。
2. **TestLayoutForTestEnv 是绿的。** CI 日志：`runtests.sh: OK - packages=[./internal/proc/ -run TestLayoutForTestEnv] top-level: PASS=1 FAIL=0`。本地复跑同样 PASS。**停手条件③不触发**。test-core 的 exit 1 来自其后的 portable 步骤：`go list` 对 23 包清单报 `package github.com/CarlosShao/wisp/internal/audio/aecspike: C source files not allowed when not using cgo or SWIG: aecsmoke.c`。
3. **settle_trace 12-r2 的字段确实在被使用**（`resident_voice_engines_windows.go:457/481` 生产调用＋两个 `_windows_test.go`），但全部调用方在 windows-tagged 文件里——GOOS=linux 的 staticcheck 看不见。**这不是 staticcheck 版本配置差，是平台形状**；修法不是删（删了 windows 构建断），是仓内已有的「注册」惯例（先例：`cmd/wisp/leg_sink_nail_131_windows_test.go:636` 的 `var _ = []any{...}`）。停手条件①不触发（它针对「被用而 staticcheck 仍报＝配置差」的情形；此处是真平台盲，处置方式仍如实申报于此）。

---

## 2. aecspike：一处修复，三处收益

**根因**：`internal/audio/aecspike/aecsmoke.c` 无 build tag。ubuntu（CGO_ENABLED=1，有 gcc）上该包没有任何 import "C" 的 Go 文件参与构建 → Go 工具链拒绝包内裸 .c 文件 → `go list ./...` 整模块 rc=1 → lint 的 `go vet ./...` 与 test-core 的 portable `go list` 同时红。本地 `CGO_ENABLED=1 GOOS=linux go list` 精确复现（go1.27.1 在 CGO=0 时静默丢 .c，所以此前本地看不出来；CI 的 go1.27.3+cgo 才炸）。

**修法**：`aecsmoke.c` 顶部加 `//go:build windows && cgo`（build 约束对 .c 文件原生有效；aecsmoke.h 无需 tag，实测 H 文件不触发该检查）。行为矩阵：

| 形状 | 修前 | 修后 |
|---|---|---|
| linux + cgo（CI lint/test-core） | rc=1 C 文件拒绝 | rc=0（本腿本地实测） |
| linux + nocgo | rc=0（静默丢 .c） | rc=0 |
| windows + cgo（本地全套） | rc=0 | rc=0（`go build` 实测） |
| windows + nocgo | rc=1（同样炸，CI 未测过） | rc=0（顺带修好） |

**收益**：lint 的 vet 步骤、test-core 的 portable `go list`、以及静态检查首次能分析 aecspike 的 linux 面（顺带暴露并修掉其 `aecspike_test.go:93` 的 SA4006，见 §3）。

---

## 3. staticcheck 78 条 → 剩 1 条（frontend，本腿禁区）

CI（GOOS=linux，staticcheck 2026.2.2，三模块）78 条逐条处置。仓内惯例核查：无任何 `lint:ignore` 先例；`//nolint` 仅用于 gosec/errcheck；「注册」先例为 `var _ = []any{...}`（leg_sink_nail_131_windows_test.go:636）与大量接口/实现的 `var _ X = ...` 行。

**处置分类**（全表见提交）：

- **真死码→删（34 处）**：`maxTrackedCards`、`probeSink`、`logf`、`calls`、`withSummarizer`、`sleepArgs`、`done`、`views`、`setSteps`、`isDeviceInvalidated`、`evSignature`、`serveGolden`、`reserved`、`inv76Abs`、`purgeAllDomains`、`equalIDs`、`entry`(archive.go)、`openCalls`、`lastReqSeen`、`composerMethodWhitelist`、`quoteOrEmpty`、`loadsOf`、`feedWindow`、`validateSourceSource`、`buildNormalizedTaintString`、`fragmentMatch`、`blacklistedSensitive`、`fsBridge`、`capB`、`failAddOn`、`secretNameRe`、`estimateTokens` 等。删除前逐一 grep 全仓（含 `_windows` 文件与全部 `_test.go`）确认零引用；孤儿化连锁逐一排除（`hrAUDCLNTDeviceInv`／`getWindowHashes`／`IsSensitiveSource`／`risk.` import 等都有独立存活读者）。
- **平台盲（linux 不可见 windows 独用）→注册（9 组）**：settle_trace 三戳方法＋承载字段（注册行在 `settle_trace_12r2.go` 尾，注释声明 12-r2 三戳逻辑零触碰）；`registerLegNail131`；`meter.reopened`；`streamOpener`/`deviceStream`；`Apm{rate,frame,sess}`；`stagingNameOf`/`stagingAttributable`；`volumeRootOf`；`hostile334`（该文件自己的注释声称「windows-tagged teeth 不是唯一读者」已过期——实测两个读者全带 windows tag，注册行注释如实记录）；`punctuationAlphabet`（注释里「positive controls 用它」一句已过期，一并更正）。
- **thresholds.go ⛔ 一字节不动**：`captureBufferLimitConversation` 的注册行放在 `internal/observe/sampler_test.go`（同包测试文件），`thresholds.go` 本体零改动。
- **风格/缺陷类机械修复（15 处）**：ST1012 `planted128Cause→errPlanted128Cause`；S1011 ×5（run.go、197 测试、ball ×2、panel）；SA4006 ×3（anthropic 空读 stop、present 死赋值、aecspike 新暴露）；SA4000（observe 双重同判）；SA5008（`json:"rec,-"` 假 option→`json:"rec,omitempty"`，探针断言的 key 不变）；ST1018 ×3（裸控制字节→`\u0085/\u0090/\u009f` 转义，字符串值逐位不变）；S1016（`subagentTextSink(t)` 转换，编译器验证）；SA9009（d22scan gitignore.go 注释重排，行首不再形似指令）；SA1019 ×3（`runtime.GOROOT` 定位器弃用→按 Go 官方指引只走 PATH；mockllm 两处保持 skip 语义、crossvet 保持 loud-fail 语义）。
- **ST1005（schema.go:449）**：错误文案去句点；同句断言在 `taskstate_196r2_test.go:179` 同步。文案语义（拒绝 vs 放行）零变化，仅结尾标点。
- **frontend/embed.go:4 SA9009 未修**：`// go:embed below is the ONLY way...` 是散文注释形似指令。修法＝注释首词重排（一字不动代码）。但 `frontend/**` 归机主外部 agent（AGENTS.md 铁律＋本腿派单 ⛔），本腿零写面。**staticcheck 步骤因此仍会 exit 1**——这是本腿唯一无法关闭的 lint 红，需外部 agent 或机主一次性授权处理。

---

## 4. 脚本与仪器改动（含派单写面之外的必要项，逐项申报）

1. **`scripts/check-path-length-budget.sh`**：roster 登记 `.scratch/wisp/issues/233-…-done.md`（126 字符相对路径，超 121 帽）。理由按既有条目措辞：票 233 与其余 61 条同代（rule 9 之前立票），漏入原始 roster 编纂，被 count guard（树 62 vs roster 61）点名。**帽值/阈值/常量零改动**；本地 VERDICT GREEN（62=62）。这是 lint job 第三个失败步骤的修复，派单写面未列、不修则 lint 不绿。
2. **`scripts/portable-tests.sh`**：core scope 认领 `./internal/present/`、`./internal/watchdog/`、`./internal/audio/aecspike/`；`./internal/audio/...` 收窄为 `./internal/audio/` ＋ `./internal/audio/aecspike/`（`internal/audio/aec` 在 linux 零测试文件，glob 声明它会被 GUARD A 拒绝）；core_pin 同 commit 同步（GUARD C 预检：linux+cgo 解析 31 包，pin==resolution）。理由段按脚本惯例写入（present/watchdog 零 tag/零 skip/零 env/零 exec，本地绿；aecspike 的 linux 面不链 sherpa）。**census 判据（GUARD D 逻辑）零改动**——三包是按 census 自己给出的第一补救路径认领的。
3. **`internal/tools/spawn_roster_census_229_test.go`**：补登 `"lifecycle-sta"` 入 `spawn229OffRosterKnown`。该站点是 43-r1（5056a365c，2026-10-31）落地的既有红——站点注释逐字写明「有意名册外」，但漏做 229 测试自己规定的登记补救。不补则 test-core 的 tools 包必红。登记理由逐字引用站点注释并注明补登者。

---

## 5. 本腿不修、如实上报的 CI 红与产品发现

| # | 位置 | 定性 | 证据 |
|---|---|---|---|
| 1 | `internal/winsec` TestAC1…342r1（CI: symlink-dir 2/400 截断；本机：junction+symlink-dir 双双复现命中） | **真产品竞态发现**（祖先换位窗）。winsec 封链语义 ⛔ 本腿不碰 | CI 日志 2007-2020 行含逐字拒绝链与 racer 错误 |
| 2 | cmd/wisp TestTicket223PermissionDenied（CI红/本机绿） | hosted runner（提权会话）上 icacls deny 咬不住，测试按设计拒绝空转假绿 | 本机 Medium 完整性下必绿；修法需 owner 拍（ledger 类别或提权探测） |
| 3 | internal/config TestReadFailureDoesNotAdopt（CI红/本机绿） | dwShareMode=0 冲突在 hosted runner 上未产生 sharing violation——windows-2025 镜像（20261021.618）行为反常，需 hosted 会话探针定位层级 | 机制合同见 manager_readfail_42r1 文件头 |
| 4 | cmd/wisp Test43Autostart（CI红） | hosted runner 的 HKCU 缺 `CurrentVersion\Run` 子键；**产品侧发现**：缺键时 apply(true) 只回「需重启后重试」，实际应创建键 | 本腿已做测试侧播种修复（§6 第 8 条）；产品行为改动留 owner |
| 5 | internal/panel TestApprovalCardView/TestComposerContract（本机红，HEAD 既有） | Go Snapshot 发出 `[instructions tasks]` 而 PanelSnapshot TS 接口未声明——Go/前端契约漂移，frontend/** 归外部 agent | panel 本机全套输出 |
| 6 | internal/panel TestStreamLogFlood/DroppedNamingLedger 35r8（本机红，HEAD 既有） | **产品发现**：key bound 之下流日志无界（4x flood→4.0x 驻留），ticket 35 :49 的无界断言在下层无对应物 | 35r8 测试输出原文 |
| 7 | census：`internal/speech`（13 测试文件）不可认领 | 结构矛盾：GUARD D exit 1 与脚本自带的「speech 刻意不认领」（加载期死）互相矛盾——census 在 speech 上结构性红，需机主裁决（DLL-staging carrier 或票面覆盖声明） | portable-tests.sh:250-260 vs GUARD D |
| 8 | `internal/audio` TestAC247、`internal/risk` TestResolvePerCallBudget（并行负载下偶发，各自单跑×3 全绿） | 既有负载 flake，非本腿引入；记录在案 | 本地批量复现记录 |
| 9 | winsec 本地闸（本机）与 CI 同数：RUN=127 PASS=75 FAIL=1（342r1/symlink-dir） | #1 的跨主机复现读数，坐实产品竞态而非 hosted 伪影 | winsec-local.txt |

---

## 6. 本腿全部写面文件（commit pathspec 与之一致）

**产品/测试码（staticcheck 清单＋新暴露项）**：`internal/audio/aecspike/aecsmoke.c`、`internal/audio/aecspike/aecspike_test.go`、`cmd/wisp/{approval_reply,run,run_test,settle_trace_12r2,leg_sink_gate_131_test,subagent_selfapproval_197_test,dataroot_128_test,autostart_43_windows_test}.go`、`internal/agent/approval/fakes_test.go`、`internal/agent/{harness_test,testtools_test}.go`、`internal/audio/{audio,device}.go`、`internal/audio/aec/aec.go`、`internal/ball/tokens_table_test.go`、`internal/llm/adaptertest/{mockllm,registry_test}.go`（registry_test.go 仅 adaptertest 下）、`internal/llm/anthropic/adapter.go`、`internal/llm/openaichat/{adapter_test,mockllm_integ_test}.go`、`internal/llm/ratelimit.go`、`internal/memory/{privacy,schema}.go`、`internal/memory/{artifacts_path_invariant_test,ticket29_e2e_test,taskstate_196r2_test}.go`、`internal/models/archive.go`、`internal/observe/{goroutine_test,sampler_test}.go`、`internal/panel/{attachments_test,config_route_248_test,git_test,l2_grant_boundary_test,workspace_handlers_test,display_scrub_216_test}.go`、`internal/perm/ticket90_persist_test.go`、`internal/present/e2e_present_30_test.go`、`internal/proc/crossvet_test.go`、`internal/risk/{tts_gate,anchor_world_334_test}.go`、`internal/session/scope_test.go`、`internal/speech/{asr_test,punctuate}.go`、`internal/tools/{fs_staging,mode,platform_other,subagent_197}.go`、`internal/tools/{helpers_test,spawn_roster_census_229_test}.go`、`internal/watchdog/watchdog_test.go`、`tools/d22scan/{gitignore,main}.go`、`tools/mockllm/chat.go`。
**脚本**：`scripts/check-path-length-budget.sh`、`scripts/portable-tests.sh`。
**证据**：本文件。

## 7. 本地验证读数（全部 rc 现跑）

- `go build ./...` rc=0；`go vet ./...` rc=0（windows 本机）。
- gofumpt v0.12.0 `-l internal cmd tools scripts frontend`：本腿文件零标记（`internal/agent/approval/pending_read.go`、`internal/agent/tools.go` 两处为本地 autocrlf 工作树假阳性——CI 同 HEAD 通过，仓内 blob 干净，本腿未触碰）。
- d22scan：`go run . -self-test` = clean 40/40 rc=0；模块测试 OK。
- path-length 预算门：VERDICT GREEN（over-budget 62 / roster 62 / not-in-roster 0）。
- census（windows 本机）：unclaimed-with-tests 4→1（剩 speech，§5#7）。
- 本地 core scope（DLL staged）两轮：第一轮 panel/risk/tools/winsec 红；**补登 229 后第二轮 tools 转绿**，剩 panel（契约漂移＋35r8，§5#5/6）、risk（TestResolvePerCallBudget 负载 flake，§5#8）、winsec（342r1，§5#1/#9）＝全部既有红。全套四数 RUN=2459 PASS=1662 FAIL=6 SKIP=11（SKIP 全部有 ledger 行记账）。
- autostart 修复：本地 PASS；223：本地 PASS（对照 CI 红＝环境差）。
- **linux 容器权威读数（golang:1.27，GOOS=linux，CGO_ENABLED=1，真 gcc＝CI ubuntu 形状）**：staticcheck 根模块**只剩 `frontend/embed.go:4` SA9009**（ROOT-SC-RC=1 仅为它）；tools/d22scan rc=0 零 findings；tools/mockllm rc=0；`go vet ./...` rc=0（修前 CI 红）；`go list ./...` rc=0（修前 CI 红）；`go vet ./internal/audio/...` rc=0；**三认领包首次 ubuntu 读数全绿**：present ok 22.95s、watchdog ok 0.33s、aecspike ok 0.12s；GUARD A 切片 present 4/0、watchdog 3/0、aecspike 3/0、audio 8/0、aec 0/0（aec 正确地不被声明）。日志存 `D:\tmp\wisp\20261101\cir1\linux-verify-out.txt`。
- CLI 闸本机（DLL staged）：第一轮 RUN=249 PASS=186 FAIL=0 但包级 FAIL（`TestAC13BringUpSurvivesAReusedThreadQuit` 窗体测试在 docker 满载编译同期猝死，单跑×2 绿）；空载复跑见本节尾部追记。

## 8. 提交

主题逐字：`CI 修复腿 ci-r1`（无 `# ` 前缀）；显式 pathspec；只 commit 不 push（推送归编排者）。
