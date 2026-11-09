# 148-verify-r1 — 票 148 收口复核（AC#1/#3/#5 复核翻框＋AC#4 契约轴核，判 `-done`）

- 腿：`148-verify-r1`（Wisp 编队写腿，2026-11-09 13:1x–13:2x +08）。工作树 HEAD `71ff0fc87`。
- 派单格：票 148 收口复核——AC#1/AC#3/AC#5 复核翻框、AC#4 契约轴核、判能否 `-done`。本腿产码零字节，只碰票面（翻勾＋Status 追加＋Progress log 追加）＋本证据件。
- **结论先行**：AC#1/#3/#4/#5 四格逐格复核全部成立（＋AC#2 已由编排者按甲案翻勾）⇒ 票 148 改名 `-done`。`148-r3`（`024681b8c`，wave48）落的三格机制账全部独立复跑对上；无一格复跑对不上。

## 0. 起手三核（现跑，不抄转述）

1. **票面现读**：`.scratch/wisp/issues/148-…unwritten.md` 全文现量。AC#1/#3/#4/#5 原文逐字在；AC#2 框确认已 `[x]`（编排者甲案裁定段在票面，派单前提「AC#2 已翻」成立）。
2. **引文出处核**：`docs/evidence/s1/148-ac2-r1.md` 现读（三把尺读数、甲/乙两案、三处过期读数清单都在）；`148-r3` 证据件 `docs/evidence/s1/148-r3.md` 在。
3. **交件真在**：`git show --stat 024681b8c` 共 10 件——`internal/agent/approval/` 9 枚＋`docs/evidence/s1/148-r3.md` 1 枚；后两笔 148 相关 commit `340df9e73`（148-ac2-r1 docs）与 `c8a7beb49`（AC#2 翻勾）只碰票面＋docs。

**前提核**：票头 `Status: blocked` 确为过期读数——台账 `A1591`（机主 2026-11-07 原话「第一个可以验收，148选B，然后别的都按推荐来」）＝Ⓐ（放开 `queue.go` 那一行）与 Ⓑ（面板载荷取白名单）两枚批准均已落；编排者甲案裁定段已在票面。本腿按惯例在 Status 行**追加**终态（原文保留），未改写任何 AC 正文。

## 1. AC#1 复核（replay 第二出口改走 clone）⇒ 翻

- **产码现量**：`internal/agent/approval/queue.go:624`＝`d := cloneDecision(it.Dec)`；`replay` 全函数（`:609-640`）零裸 `d := it.Dec`。
- **全仓扫**：`grep -rn ':= it\.Dec' --include='*.go' internal cmd | grep -v _test.go` 现跑＝命中 2 处，均为标量字段读（`queue.go:406 settled := it.Dec.Tool`、`:414 tool := it.Dec.Tool`），不携带任何引用槽位背衬——票面立的那味（按值交出整个 `tools.Decision`）已清零。
- **会响的检**：`ticket148_replay_backing_test.go` 三向形状（值保真＋背衬隔离 fail-closed census＋字面就地写断言：只经返回值 clobber 每一槽位、读回 history 记录要求不变），头部自陈摘掉 `cloneDecision` 那一行全红（突变读数登记在 `148-r3.md`）。
- **票 146 仪器**：改名后仍在——`TestCloneDecisionSharesNoReferenceSlotWithTheQueue`（`ticket146_liveapprovals_backing_test.go:222`）、`TestCloneDecisionInPlaceWriteCannotReachTheQueueRecord`（`:276`），反射走 `tools.Decision` 具名类型，第 9 枚引用字段进不来。
- 以上全部在包跑绿（见 §3 测试读数）。

## 2. AC#3 复核（往外送什么有边界句）⇒ 翻

- **有边界的句子真在产码注释**：`pending_read.go` 的 `panelVerdict` 头注释具名白名单＝「the fields the panel needs」（Ⓑ 批准方向），并逐字列 6 字段 census 依据。
- **白名单本体**：`panelVerdict(d tools.Decision)` 现量＝`Tool`/`Args`（`cloneBacking`）/`Level`/`RulesHit`（`cloneBacking`）/`Reason`/`SessionOverrideBlocked`，其余字段一律不出口。
- **census 对上生产读者**：`cmd/wisp/panel_pump.go:87-92` 现量逐字读的就是这 6 字段（`Args` 经 `panelArgs`），全文零 `Params` 读；`CorrelationID`/`Position` 走 `LiveApproval` wrapper（非 Decision 字段）。
- **`Params` 二层处置**（AC#3 硬性要求「必须一起处置」）：`Params` 不在白名单 ⇒ 面板载荷根本不携带它 ⇒ 二层容器（`[]any`/嵌套 map）从 `LiveApprovals` 这扇门**结构性出不了队列**——「closes by exclusion」，注释具名＋`cloneDecision` 保留的一层拷贝只服务 `Queue.replay` 的进程内消费＋台账 `A261` 登记未抹。
- **会红的对账**：`ticket148_panel_whitelist_test.go` 双尺——ruler 1 反射 `tools.Decision` 全字段要求＝allowed∪dropped 完整划分（fail-closed，新增未分类字段即红）；ruler 2 对真队列 enforcement（allowed 字段值到达且引用槽位自带背衬、dropped 字段零值；换回 `cloneDecision` 或删任一白名单行即红）。全绿（见 §3）。
- `LiveApprovals` 现量逐行走 `panelVerdict(it.Dec)`（`pending_read.go:177` 附近）。

## 3. AC#5 复核（逐包门禁）⇒ 翻

门禁一律在**导出树**跑：`git clone --shared` @ `71ff0fc87` 于 `/d/tmp/wisp-148v-gate`（保留备查）。测试环境按派单配方：`unset GOTMPDIR TMPDIR TEMP TMP`、`export PATH=$PWD/build:$PWD/third_party/sherpa-onnx:$PATH`、未设 `CGO_ENABLED=0`。

| 门禁 | 读数 |
|---|---|
| gofumpt CI 形状 `-l internal cmd tools scripts frontend`（`/d/work/base/gopath/bin/gofumpt.exe`） | 空，rc=0 |
| gofumpt 票面形状 `-l . tools/d22scan tools/mockllm` | 命中 20 枚，**全在 `.scratch/wisp/probes/**` 历史探针件**（`148-r3.md` §6 / `148-ac2-r1.md` §5 已登记的结构性撞，CI 形状为准） |
| `go vet ./internal/agent/approval/` | rc=0 |
| `go build ./...` | rc=0（导出树） |
| `go test ./internal/agent/approval/ -count=1` | `ok 5.673s`（含 5 秒 subject 检 `Test148WorkPeakSubjectRunsToolCallsLiveApprovalsFires`） |
| **名册两向差集**（`024681b8c^` vs HEAD，`go test -list` 后 sort，100→105 枚） | 消失 2＝票 146 仪器**改名**（`TestLiveApprovalsSharesNoReferenceSlotWithTheQueue`/`TestLiveApprovalsInPlaceWriteCannotReachTheQueueRecord` → `TestCloneDecisionSharesNoReferenceSlotWithTheQueue:222`/`TestCloneDecisionInPlaceWriteCannotReachTheQueueRecord:276`，断言同形仍在）；新增 5＝`Benchmark148LiveApprovalsArms`、`Test148LiveApprovalsCarriesOnlyWhitelistedFields`、`Test148PanelVerdictWhitelistRosterIsClosed`、`Test148ReplayReturnsACloneNotTheHistoryRecord`、`Test148WorkPeakSubjectRunsToolCallsLiveApprovalsFires`。100+5=105 对平，**零无名消失** |
| `sh scripts/d22scan.sh` | rc=0 clean；examined 全非零：bans#1-5 internal/=387、cmd/=64、ban#6 frontend/=121、ban#7 tools/=48、ban#10 risk/=21、ban#8 design/=14、frontend/=121、internal/=966、cmd/=242 |

## 4. AC#4 契约轴核 ⇒ 成立（报读数）

1. `git show --name-only 024681b8c` 10 件**零点名**：`internal/risk/**`、`thresholds.go`、golden、`allowlist.txt`、`tools/d22scan/**`、`scripts/slo-check.ps1`、`docs/PLAN.md`、`docs/specs/**`、`frontend/**`、`design/**`、`internal/panel/l2_grant_boundary_test.go`。
2. 后两笔 148 相关 commit（`340df9e73`、`c8a7beb49`）`--name-only` 现跑：只碰票面文件＋docs 证据件，契约轴零命中。
3. 工作树契约轴 `git status --porcelain`：零 M 件；`design/**` 仅有 `??` 未跟踪件（机主外部 agent 面，非本票改动）。
4. `queue.go` 在 `024681b8c` 的改动＝Ⓐ 批的那一行（replay，A1591）；`approval.go` 零字节。两批批准在先，不违「未被批准前零字节」条款。

## 5. 判定

AC#1/#3/#4/#5 四格成立＋AC#2 已翻 ⇒ **票 148 改名 `-done`**（`git mv` 后直接 `git commit -- 新名`，落账双名核旧=0/新=1；旧名删除若留脏单独补一笔）。Status 行追加终态、Progress log 追加一条，AC 正文零改写。

**诚实口径**（沿用 `148-ac2-r1.md` §2，不重复主张）：满深度一次 `LiveApprovals` 比裸值拷多 **3968B/72allocs**（bench 机制账上观察到的差异）；生产 WorkPeak 档整表读数结构性不可达（subject sit-still＋`internal/proc` 零 approval 引用），缺口维持票 146 §4.3 永久登记，不是本票欠账。本腿未写「无代价」，未合成任何数字。

**只 commit 不 push。**
