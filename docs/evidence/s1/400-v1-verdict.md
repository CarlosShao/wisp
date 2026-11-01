# 票 400 终裁腿 400-v1 · 裁决表（非实现者）

- 腿＝`400-v1`（终裁、非实现者）· 2026-11-01 · 对象＝`e44acb69a`（10 枚；HEAD 相对本审计面 diff 为空另证）。全部读数**本腿现跑**，主尺在导出树 `D:\tmp\wisp\20261101\400v1\tree`（`git archive e44acb69a` 纯读抽取，`git hash-object` 四枚逐字节＝commit blob：`f8933aa4a/c86e9b9e6/cbafe712e/b8aaddb026d4`），落地树（主工作树）配对复核；突变三发各在腿外一次性副本 m1/m2/m3。主树全程零已跟踪改动。

## ① panel 契约对齐 —— 成立

- 两把对账尺（导出树现跑 rc=0，与 `400r1-impl.md` §AC#1 逐枚同号）：`ApprovalCardView 10／ResultChunk 3／Snapshot<->PanelSnapshot 6`；`Snapshot 6／ComposerState 11／ModeView 3／WorkspaceView 6／AttachmentRef 9／ResultChunk 3`。
- 嵌套形状（两把尺只对顶层键、无仪器覆盖，本腿逐键手核）：`ProjectInstructionFile 8／InstructionsSection 3／GitWorktree 5／GitView 10／TaskRowView 12／TaskRosterSection 6` 全部与 Go json tag 一字不差；`CredentialState` 五枚举值（unknown/config_unreadable/no_ref_declared/all_recorded/partly_missing）逐值一致。
- `tsc -b` rc=0；oxlint rc=0（**35 warnings 0 errors**，与实现件同数＝存量）。纯类型增量核：+118 行全部为 interface/type/注释，零行为面；`frontend_hygiene` 五枚形状断言不受影响（落地树配对 PASS）。
- 修法定案前提现核：C17 入站守卫 `knownComposerMethod`（internal/panel/bridge.go:154）无任务入口；`PLAN.md` C17 行逐字含「**任务状态**」推送 ⇒「Go 侧收敛（少发键）」被冻结契约否决；`missing-features v3 §3.3`（:202-212）逐字「转绿只有一跳：界面那侧 frontend/src/lib/panel.ts 补 tasks?:」。同病同治属实：ComposerState 五键（git/currentModel/modelKnown/credentialState/credentialKnown）Go 侧自 181/187/248 起恒发、TS 名册一直未跟。

## ★frontend 编辑合法性裁定 —— **追认**（附三条件）

- 事实链（全部现读）：400-r1 派单存档逐字「**⛔ frontend/design（机主外部 agent）**」；A273（台账 :6768）＝「批了 **internal/panel** 局部解冻」，解冻存档（2026-09-26-093x）更逐字「**frontend/\*\* 继续完全隔离：不读、不写**」⇒ **A273 引用错位成立、编辑越界成立、且属明知**（实现件自认「具名判断」，并预明「审查不认则 revert 该一枚即回红」）。
- 两选项代价（本腿现量）：**回退**＝两把对账尺即刻回红（M1 实测红句逐字＝ci-r1 原句）、票 400 格①不成立、16 收口格①③残余扩大、30-w-v1/384-v1 指名要的声明落空。**追认**＝类型增量与 Go 侧逐键一致（双尺＋tsc/oxlint＋手核）；代价＝外部 agent 的树被本编队动过一枚、须知悉（A1363② 已在册「外部 agent 下次动 frontend 会覆盖」风险）。
- 裁定＝**追认**。理由依机主铁律（判据＝最有价值、非风险最小）：内容经全仪器核真且为唯一诚实修法；回退是「动作更小」不是「更有价值」；治理修复走**披露＋记账**，不销毁正确修复。
- 附条件：①外部 agent 必须知悉（措辞见下）；②编排者补 A 账记「越界＋追认＋引用错位教训」（A1363② 已立底）；③此后此类跳点一律先转界面支或先取机主预批，**不许再以「具名判断」自我越权**——本追认不构成先例授权。
- 给机主的措辞（可直接转外部 agent）：
  > frontend/src/lib/panel.ts 于 11-01 由 wisp 编队 400-r1 腿补了 PanelSnapshot/ComposerState 的纯类型声明（+118 行、零行为改动）：Go 侧自票 181/187/197/200/248 起就在线上发送这些键，TS 名册一直未跟、契约对账尺因此常红。该编辑未经你们的树权限，已由终裁腿逐键核对与 Go json tag 一致（tsc/oxlint 全绿）。请知悉，并在下次改动 frontend 时以当前 git 历史为基线合并；若有并行未提交改动，请先 git pull 后核对 panel.ts 再续作。

## ② 35r8 流日志有界 —— 成立

- 码核：行级 `streamRowKeepRunes=4096` 挂 `streamKept.window`，键数之下也 clamp（head/tail＋计数 elision，shrink 计入同一本账）；截断档 256 行为路径逐字保留（`clampLocked(key, streamTruncateKeepRunes)` 于 truncated 形、`TestTruncationFactsRideThePacket` 在包内绿）。具名台账 `streamDroppedNamesCap = DefaultStreamKeys(32)×StreamKeyHardCeilingMultiple(2) = 64`（Ring）＋`DroppedCount()` 总数 accessor。
- fan-out 同号核（自跑）：`1x text=35072 rows=64 names=64 | 4x text=35072 rows=64 names=64 | cap text=36864 rows=64` —— 留存与 flood 规模无关（heap 读数随 GC 抖动，谓词余量大）。
- below-bound：`1x 24690 / 4x 24693`（差 3＝marker 位数）；行级 elided 恰＝flood−8192 断言在场且绿；flood 形状 3×200/800×2000 与 4000/16000×1×1 一字未动。
- 正控在场且打红：`TestFloodRulerFiresOnAnUnboundedSink35r8`（unbounded sink 800000 runes>cap 36864、rows 100>64）＋绿尺判语区分有界/无界 sink。

## ③ winsec 342r1 竞态修 —— 成立（封链语义内核核过）

- 码核：`privateFile` 交 `pinnedParent(path, resolved)`（winsec.go:194，产码唯一调用点）；Windows walk 改下探 `pathComponents(caller)`＋非绝对拼写具名拒；`resolved` 保留三职＝traversal verdict＋全部拒答文本锚＋体积根（`pinVolumeRoot(resolved)` 不变）；POSIX `_ = caller` 走 `resolved` 逐字不变（floor/C26 直通不改写）；fallback 同签名拒。risk/C26/入口签名/旗标/封链形状零触碰。
- 6400 attempts（自跑 8 发×2 形×400）：**truncated 累计＝0**；liveness 逐发在场（provedReparse 14–98/形/发、formBLegRefusals 3–50/形/发）＝342 AC#1 口径两半齐（修前分子>0、修后 0 且 liveness 非零）。

## ④ 反向三发 —— 全部有牙（本腿自跑）

- M1 panel.ts 删 `instructions?:`/`tasks?:` 两键：红句逐字 `Go Snapshot emits [instructions tasks] that interface PanelSnapshot does not declare`（＝ci-r1 原句）。
- M2 pump.go 两调用点回 truncated-only 守（=修前行级无界形）：below-bound 绿尺红＝`pinned text 4800000 runes > the cap 36864`＋`retention scaled with the flood: 1x pinned 1200000, 4x pinned 4800000`（heap 数随 GC 抖动，方向与量级同实现件）。
- M3 winsec.go `pinnedParent(resolved, resolved)`（走回答案拼写＝修前形）：TestAC1 **3/3 发红**（symlink-dir 形 truncated=3/3/1、liveness 同发在场 provedReparse 13–92/refusals 3–40；junction 形本轮 0——窗在 identity 面、与修前签名一致）。
- 还原核：突变只发生在腿外副本 m1/m2/m3；导出树四枚审计面件 `git hash-object` 逐字节＝commit blob；主树 `git status` 已跟踪改动＝**NONE**。

## ⑤ 零回归＋门禁 —— 成立

- panel（导出树）`RUN=271/PASS=170/FAIL=3/SKIP=0`——3 红逐一自证导出伪影（无 `.git`／junction node_modules walk `Incorrect function`／渲染夹具缺失），落地树配对复跑三枚＋两尺全 PASS ⇒ **落地口径 271/173/0/0 成立**。winsec `RUN=127/PASS=76/FAIL=0/SKIP=0`（逐字同实现件）。
- cmd/wisp（导出树，PATH 前置仓根与实现件同法）`RUN=568/PASS=375/FAIL=55/SKIP=4`（463s）。55 红普查分类：**36** 语音模型库缺（build/models manifest／third_party/spike-models／WISP_SPIKE_MODELS 未设——红名族 15e/41a2r1/26P/392/381r2/390r1/26R3/382WR2/WireR1 与实现件红名册同族）＋**14** sherpa DLL 缺（`third_party/sherpa-onnx` 全 gitignore 未跟踪、不在 git archive 内＝导出伪影，且正落「涉 sherpa＝staged（A1326 卡）」域）＋**1** 导出树无 `.git`（AC11 rev-parse 128）＋**4** 窗/驻留族（33r5 focus／AC14 push／256 landing reading／41a2r1 note——四名皆在实现件红名册同族）。**零枚触 panel/winsec/panel.ts 审计面**；红计数 55 对实现件 41 的差＝14+1 导出伪影，逐枚有因。
- 本腿真实消费面（cmd/wisp `-run 'Test197|TestSubagentStreamKey|TestRunPacket|TestBookedSnapshot|TestRunBooks|TestSnapshotWorkspace|TestThePumpIsDriven'`）：**RUN=15 PASS=15 FAIL=0 SKIP=0 rc=0**（roster/streamlog/instructions/197/329 全族绿；实现件口径 17 与本腿 15 差在 regex 组成，实质同为 0 FAIL）。
- d22scan（导出树 `scripts/d22scan.sh`）**rc=0 clean**，分母与实现件逐枚同：bans#1-5 internal/=287、cmd/=52、ban#7 tools=25、ban#8 internal/=732、cmd/=183；正控（runtests.sh 强制 count=1，含 `TestBuiltBinaryGoesRedEndToEnd`）＋selftest 全绿。`go vet` rc=0；gofumpt（panel+winsec）0 命中。⚠ 未按 ban#8 扫 docs/（派单禁令）。
- 环境注：本腿导出树无未跟踪原生件（sherpa DLL/模型库/dist），由此产生的红已逐枚归因；时序敏感件（TestAC1 8 发、M3 3 发）均在盒子独占窗跑。

## ⑥ 五格收口裁定＋残余定性

- AC#1 panel 契约对齐＝**成立**（以上面「追认」为收口条件）；AC#2 35r8 有界＝**成立**；AC#3 winsec 竞态修＝**成立**；AC#4 反向判据×3＝**成立**（三发有牙＋还原核）；AC#5 纪律格＝**成立**（禁改面一字未动、只 commit 未 push、临时件全在仓外）。
- ⇒ **票 400 可收口（加 `-done`）**，条件＝机主圈阅本表「追认」段（那是改树域边界的一次性追认，属摆机主三类；A1366④ 摆机主清单已列 panel.ts 追认）。
- 残余（不阻收口，须登记）：①wire 携 `DroppedCount` 需动 cmd/wisp 装配点 **panel_pump.go:250**（现读坐实 `state.DroppedStreamKeys = rt.stream.DroppedKeys()`；本票写面外，立票归编排者）；②票 342 AC#1/AC#4 两格仍开（Windows 半区终态判读归 342 自己的终裁线，本表不代勾；本腿交付的 6400+M3 读数可供其引用）；③frontend 追认若机主不圈阅，回退路径＝revert `frontend/src/lib/panel.ts` 一枚即回红（代价已写明）。

## 纪律格

- 本腿写面仅 `docs/evidence/s1/400-v1-verdict.md` ＋仓外 `D:\tmp\wisp\20261101\400v1\**`（导出树/m1-m3/logs/脚本）；⛔ 落 `.go`、⛔ 翻框、⛔ 改票面/台账、⛔ push/amend/reset/rebase/stash/checkout；commit 显式 pathspec；ban#8 未扫 docs/。
- 未定义即停检查：全程未遇派单未覆盖情形；导出树三处伪影以**配对复跑**定因，未做未授权假设。
