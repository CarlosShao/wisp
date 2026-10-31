# 票 43 终裁腿 `43-v1` —— power lifecycle events 裁决表（非实现者）

> 尺：`5056a365c`（43-r1 交件，11 枚 1943 行）**导出树**现跑 2026-10-31 14:03–14:25 +08。导出树 `D:\tmp\wisp\20261031\43v1\export`（`git archive 5056a365c`；junction 只读穿透未跟踪模型资产 third_party/{sherpa-onnx,spike-models,model-fetch}＋build/）。变异与正控全在仓外；**工作树零触碰**（387-r1/386-r1 在飞面零接触、零归因）。在读数前补在册 PATH（A1307补）；本腿自己也漏过一次（见 §5.3）。

## 1. 必跑清单逐条读数（全部本腿现跑，⛔ 采信实现件字面）

| # | 项 | rc | 读数 |
|---|---|---|---|
| 1 | WM_QUERYENDSESSION winlive 真子进程 | 0 | `TestLive43QueryEndSessionExitsWithin5s` PASS：**exit 94.8576ms**（实现件 90.8ms 同量级；≤5s 红线余量约 52 倍）；trail 机器断言 step1<4(stop-audio)<5(release-speech-sessions)<9(close-job-object)<10(exit) 过＋fast-path 具名行在；exit 0＝走自身关停链。子进程直 boot 读数：`电源通知=已注册`＋`WM_QUERYENDSESSION 门已建（lifecycle-sta 线程自泵）` |
| 2 | L2 voided＋D43 #24 行核 | 0 | headless 挂起案 PASS：真 queue 卡（approval.New 真门＋真队、ask 停真 answer 通道）→ voided=2、parked ask 收 AnswerReject 且理由含「挂起」、Pending()=0、console 逐字「错过 2 个待确认操作（系统挂起，已全部作废并按拒绝处理；恢复后不会重放）」、审计 `SUSPEND-L2-VOIDED why="43-suspend" refused_l2=1 voided_l1=1`。**表核**：PLAN.md:3094＝#24 原文（AwaitingApproval｜队头被拒绝/超时 300s→Acting）；编码行现跑 `internal/statemachine/table.go:168`（EvApprovalDenied→Acting）＋`:172`（EvApprovalTimeout→Acting）；`git diff 5056a365c HEAD -- internal/statemachine/table.go`＝0 行、交件未触该文件 ⇒ **D43 表零改动（无大红）**。走行判定：`voidForSuspend`→`cards.Reject`＝队列自身拒绝动词、落点即 #24 拒绝行；M3 变异证明该漏斗真被调用（拔除即红） |
| 3 | kill -9 | 0 | `TestLive43KillNineMidTaskMarksInterruptedOnNextBoot` PASS：真 TerminateProcess → marker 存活 → 第二次 boot stdout「已标记为中断」→ DB `state=interrupted` → `ListTaskLogs` 恰 1 行（**零重跑**）→ 干净退出后 marker 消失 → sweep 0/nil |
| 4 | autostart | 0 | 两枚 PASS：真 HKCU 往返（参数化值名 Wisp43r1Test，机主真键 `Wisp` 零触碰）写→读回带引号路径→幂等「已存在且一致」→删→诚实「无本程序条目，无动作」；OFF default＝schema 零值＋fresh boot receipt `autostart=false…无动作`（子进程直 boot 同证） |
| 8 | mutation 三发（仓外导出树） | 各 1 | M1 摘 dispose 挂接→`lifecycle_43_windows_test.go:248: EnginesHot after suspend = 1, want 0`；M2 摘 SetStringValue→`autostart_43_windows_test.go:58: Run value "Wisp43r1Test" missing after apply(true)`（收据谎报被注册表抓）；M3 摘 L2 拒绝漏斗→`:225: voided = 1, want 2`。**还原核**：三发各自还原哈希核过＋本腿两枚文件 MUTANT 残留 0（全树 28 命中均他票既有字样）＋还原后 Test43 9/9 rc=0 |
| 9 | 零回归 | 0 | 四族合跑 **19 PASS/0 FAIL**（15e 6＋26R3 3＋41a2r1 5＋381r2 5），`LATENCY openDelay=262ms chunk=640ms ballLatency=902ms …seamWaitToFirstText=5.51s finalLatency=4.096s speechMS=3392` **逐位同值**；`internal/session` **42/42** rc=0；`internal/config` ok；255 名册 20/20（含 `RosterStillMatchesTheActualReadSites`）rc=0；`go vet` 默认＋winlive 双档 rc=0。★233 §8 形状核：config_readers_255.go diff＝只加 sectionReadSites `"App": {"cmd/wisp/autostart_windows.go"}`＋app.theme 行文字同步，无第三类改动 |
| 10 | 两坑 | — | **dumpbin 实勘**：powrprof.dll 导出 #85＝`PowerRegisterSuspendResumeNotification`；user32 只有 `Register/UnregisterSuspendResumeNotification`（无 Power 前缀、签名不同）⇒「user32 无此导出、猜错即 panic」坑成立，`Find()` 防呆在码（lifecycle_windows.go:312）。**泵线程形状**：createWindow 在 pump() 内（:347）、pump 在 lifecycle-sta Spawn（:191）⇒ 窗口与消息队列同线程；94.9ms 即此修法直接产物。⚠ 在册 WARN `goroutine outside the D38 roster (leak symptom) goroutine=lifecycle-sta`＝名册外先例形状（panel-sta 同款、33 钉过「加名册=人工批准」），非泄漏 |
| 11 | d22scan | 0 | 导出树 **rc=0 clean**：bans#1-5 internal/=280 cmd/=48、ban#8 internal/=715 cmd/=169（与实现件逐位同值）。**正控**：种 `≤`（U+2264）→ rc=1 逐字点名 `…zz_d22_positive_control_43v1.go:5: [emoji] ban #8 glyph …string literals included; comments are exempt per Q-46(c)` → 拔除 → 复扫 rc=0 ⇒ 尺活 |

## 2. ★★★ 三个关键裁点

**裁 5｜替代形合法性（格 1）——合法，是「替代形」不是「没做」。** (a) 真睡眠不可自证是结构性事实：无人值守会话随宿主一起睡死、唤醒需装置＋BIOS 支持，测试会话无法自证——本判定不受「怕麻烦」驱动。(b) 双门是真组件：PowerRegisterSuspendResumeNotification 真注册（子进程直读「电源通知=已注册」＋dumpbin 实勘导出在 powrprof）＋WM_POWERBROADCAST 顶层窗口真建（live 日志）。(c) 序列是生产体：通知回调体与 `deliverPowerEvent` 同函数（:305-308），漏斗单线程串行。(d) headless 读数是真组件：真 approval 门/队、真 audio HalfDuplexGate、真 session.NewScope；fake 仅在 C8/C5 认可接缝（ASR 引擎、采集源）；EnginesHot 1→0＋引擎 Close 恰 1＋`SCOPE-DISPOSED reason=suspend:43-suspend cascade=models-unload,panel-hide,free-os-memory,grants-expire`（28 的级联原样跑）＋采集 stops=1/starts=2、恢复后冷装载 loaderCalls=2。(e) 诚实具名齐：替代形标注、VAD 微检测器不属 SessionScope 保持驻留（不虚报「全部模型离场」）、真断电一跳留给机主在场 psshutdown（harness FindWindow＋SendMessage 形状可复用）。

**裁 6｜定向注入合法性（格 2）——合法，「唯『消息从系统来』一跳是注入」的诚实声明成立。** (a) HWND_BROADCAST 打全桌面全部顶层窗口＝外向动作（他人应用会真应答关机、殃及用户未存工作）；定向 SendMessageTimeoutW 到子进程**自己的**窗口（class＋pid 两事实身份、228 live tier 同规）零外溢。(b) 泵/wndproc 即答 TRUE/fastExit/RequestStop/boot 自身 defer→rt.Shutdown(fast)→D38(e) 十步全生产；注入只占「消息到达」一跳——真实关机即系统发同一消息进同一窗口；与 C8/C5/C17 接缝同形（OS 边界即接口），生产代码零测试分支。(c) 非伪授权形状：无任何授权/用户手势被伪造（injection-timeline 判据不触发）。(d) AC 的「broadcast→exit ≤5s」预算 govern 到达后的退出路径；实测 send→exit 94.9ms 是保守读数。

**裁 7｜§16.11#10 定案判语核——判语成立、归口正确。** PLAN.md:3244-3247 原文：倾向「作废并判拒绝」（陈旧批准是安全事故）、「本轮未定案、须在 S7 切片卡明确」；DEFERRED 行（:1536）同倾向。实现给出的最终决定＝该 lean 的可执行化、逐点一致；授权链成立：票面 Key constraints 逐字「follow the §16.11#10 lean: void + auto-reject + '错过 N 个待确认操作' notice (decision logged; final wording at 54 review)」＝票级明示；落点只骑**已冻结**的 D43 #24 拒绝行、未发明新边（293 先例）。**归口正确**：PLAN.md 的 D43 note/DEFERRED 行/§16.11 文字＝契约文字（SPEC-12 §4.1＋铁律三类摆机主），实现腿零触碰、摆编排者/机主批——既未擅改契约，也未拿「要批准」当不实现的借口。

## 3. 五格逐条判

| 格 | AC 原文缩写 | 判 | 缺什么 |
|---|---|---|---|
| 1 | suspend/resume→pre-sleep dispose、audio restored、resume-e2e | **可翻** | 无行为面缺口；残余＝「真断电一跳」需机主在场 psshutdown（替代形已如实具名） |
| 2 | WM_QUERYENDSESSION→exit ≤5s、4/5/9 顺序、Job children | **可翻** | 无（94.9ms＋trail 机器断言＋step9 job close＝OS 杀子进程本体） |
| 3 | pending L2→voided＋计数（D43 note 更新最终决定） | **可翻** | 决策面全落地；残余＝PLAN.md note 文字待机主批（判语见 §4） |
| 4 | kill -9→中断可见、never auto-rerun | **可翻** | 无（真杀全链复刻；无重跑是结构性事实＋行数断言） |
| 5 | autostart toggle Run key＋OFF default | **可翻** | 残余＝托盘「开机自启」菜单项归 internal/ball 冻结菜单表（非本腿写面）；config 半边全通 |

**结论：可带残余结案**（15/28/41 先例：AC 勾保持未勾）。残余四条定性：①托盘菜单项＝真残余、归下一张能写 internal/ball 的票（ball 冻结面）；②wall schedules＝非缺口（本构建无墙钟派生日程装置、无可重算项，resume 审计行具名）；③update-pending 消费端＝事件面已落（UPDATE-APPLY-HANDOFF 登记、文件不移动不删除＝56 协议），消费面归票 56；④D43 note 文字＝契约面摆机主（唯一需机主动作的残余）。

## 4. 摆机主：§16.11#10 定案判语（建议文字，批准后由编排者落 PLAN.md 三处）

「`AwaitingApproval` 遇系统挂起：**一律作废**全部待确认项——L2 走 D43 #24 拒绝漏斗（队头被拒绝→Acting、取消该调用）、L1 倒计时窗口作废＋审计行；按**拒绝**处理并以 console 告知「错过 N 个待确认操作（系统挂起，已全部作废并按拒绝处理；恢复后不会重放）」；恢复后**不自动重放**（重跑/重放只能由操作者重新发起，与崩溃恢复同一姿态）。本决定落 D43 转移表 note、§16.11 第 10 条、DEFERRED 行三处；最终措辞归票 54 review 定稿。」

## 5. 瑕疵与本腿自陈

1. 实现件 §4「session 27 案全过」＝28-r1 改前旧数；现分母 42 案（42/42 绿 rc=0）——数字陈旧瑕疵，不折损结论。
2. M1/M3 红句行号差一行（实现件引 :247/:224、现跑 :248/:225）——测试文件行号漂移，断言文字逐字相同。
3. 本腿四族首跑 6 FAIL、再跑 2 FAIL，红句均为导出树缺**未跟踪**模型资产（`build/models/manifest.json`、`third_party/model-fetch/x-punc/...`）——junction 补齐后 19/19；另有一次 `0xc0000135`＝本腿自己漏补在册 PATH（A1307补 默认动作罚单自领，补后即绿）。
4. 工作树零触碰：本腿仓内唯一写面＝本文件；临时件全在 `D:\tmp\wisp\20261031\43v1\`（export/bin/logs/pristine/cases15e，只建不删）。
