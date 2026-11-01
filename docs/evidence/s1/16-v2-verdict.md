# 票 16 对抗验收腿 `16-v2` —— 格⑤终裁（非实现者裁腿）

> **身份**：票 16 格⑤「Adversarial acceptance pass by a different agent (stub scan +
> contract conformance)」的执行腿（重派——16-v1 判死 11h 零产物）。**被裁件**＝16-r1
> （`83c08b2f6` 三枚：`cmd/wisp/voice_e2e_demo_16r1_windows_test.go` 541 行＋
> `docs/evidence/s1/16r1-impl.md` 153 行＋`docs/SLO.md` 附录 D 纯追加）。
> **尺**：`HEAD=eb1c42c26`（dev）。**方法**：读核＋导出树复跑——`git archive HEAD` 解树至
> `D:\tmp\wisp\20261101\16v2\export`（junction `build`/`third_party`＋三 DLL 贴身＋`WISP_*`
> 五 env＋cwd=导出根，A1326 卡同款 staging；测试二进制 Git Bash 直调、rc 直取、
> `unset GOTMPDIR`）。主仓零写面（本文件除外）、零 push；在飞 `26p-v2`/`382-w-v2` 只读
> 同行零归因；音频错峰（每发真设备跑前 tasklist 现核）。
> **写面**：本文件＋`D:\tmp\wisp\20261101\16v2\**`（logs：demo／sixfamilies-run2／
> warmgate-solo／corrupt／cer-ci／livemic×2／kill-run1..verify／resume-range／d22scan* 等）。

## 0 一句话判语

**16-r1 三枚全部够格，票 16 可收口**（对照 15-done 先例＝第 48 票带 10 条具名残余收口）：
格①②③ 翻（各带具名边界）、格④ 维持残余（真 CI 工件待 push）、格⑤ 本件即其本体。
唯一升级发现：**D43#15 工具臂与 Acting 家族（#17-20）在今日产码全仓零生产派发位**
（不止语音路）——产品级缺口，非 16-r1 伪造，归口 26 集成面。

## 1 格⑤-1 stub scan（S2 面；白名单＝C8/C5/C17 接缝）

| 扫描面 | 读数 | 判 |
|---|---|---|
| `todo!()`/`unimplemented!()` | internal/speech＋internal/audio＋cmd/wisp **0 命中** | 清 |
| S2 产码 stub/mock/fake 标记 | 逐命中皆注释性接缝说明或明确否认（`device.go:3` 枚举缝＝热拔注入假**设备清单**；`wasapimic_windows.go:15`「audio-data seam is never mocked」；`run.go:141`「not a mock standing in」；panel Bind stub＝go-webview2 架构词） | 清 |
| demo 的每个注入点 | golden SSE＝C5 缝内（`golden.ServerFromBytes`）✓；wav＝C8 缝（`audio.NewWavInjector`→BoundedFrames）✓；`demo16Bridge.submitTask`＝保持异步形状、直调**真组合根** `runTextTask`（run.go:202，`wisp run` 同一函数），非替真；`onRuntime`＝产线 boot 逐字复刻（resident_windows.go:580/:610）；speaker/TTS 工厂**未覆盖**（adjust=nil）＝真 matcha 真 waveOut | 白名单内 |
| 替真引擎 | KWS／Paraformer／matcha／waveOut 全真（复跑转写逐字同、88828 样本、sha256 验过的 store） | 零替真 |
| DEFERRED 对账（SPEC-10 §8-2） | 16-r1 零产品码改动⇒无新标记；`internal/speech/doc.go:19` 为票 15 已登记形 | 清 |

## 2 格⑤-2 contract conformance（D43/D10/D32/C 系）

- **D43 逐行走位**：五行 Sleeping+kws.enabled→Armed(#5)／Armed+kws.wake-word→Listening(#7)／
  Listening+asr.vad-stop→Thinking(#11)／Thinking+llm.first-token→Speaking(#15 纯文本臂)／
  Speaking+tts.speak-done→Warm(#26)——本腿复跑逐行同、全部产线派发（demo 测试文件无一处
  手工 machine.Dispatch，现核）；每行与冻结表（PLAN.md:3065 起）逐字对上，无非法转移。
- **D10**：run.go:1511-1537 完成结果＝notify＋voiceAnnounce **同一句**；demo
  `NOTIFY 1 posted ("Wisp\0已经读完了…")`＋spoken=1 双落。
- **D32 16.3.4**：预算行（PLAN.md:2330-2346 现核）vs SLO 附录 D 12 行——见 §4 格②。
- **C 系**：C5 golden 注入✓；C8 wav 注入✓；C29 sha256（真翻转→mismatch→拒交未验证字节）✓；
  C31 SCOPE-WARM→SETTLING(3s)→DISPOSED✓；C19/D34：`fs.read` 授权目录内＝L0/allow（D34 表
  逐字「L0（授权目录内）/L2（越界）」），tool_row level=L0 decision=allow outcome=success✓。

## 3 三枚诚实边界裁量（16-r1 照录的三件）

1. **Acting 无球态载体＝缺口（产品级，非证据级）**。D43#15 行语义＝「首 token＋有 tool
   call→Acting」（table.go:101 guard `first-token-has-tool-call`；machine 测试钉住双臂）；
   demo 回合 SSE 首响应确实带 tool_calls，而产线在 :979-981 以 `Facts{HasToolCall:false}`
   走纯文本臂。**升级**：全仓 grep——`HasToolCall:true` 的生产派发位＝0；
   `EvApprovalNeeded/EvToolsCompleted/EvRepeatThreshold/EvPathConflict` 生产派发位＝0
   （只存在于冻结表＋表测试＋ball 交互测试）⇒ Acting 家族**机器完备、产线零接线**。
   demo 逐行走位每行仍合法（无非法转移；「该句为纯文本」对所播报的句属实）——demo 不假
   报，但票面 AC#1 链字中的 Acting 一环在 HEAD 无可观测载体：**判缺口**，归 26 集成面挂账。
2. **settle token→Settling 8984ms＝口径差／如实形状**。该腿含真 TTS 装载 3.7s＋播报 4.1s
   ＋staging 暖窗 1.2s，与 12-r2 桩文本路 1201-1206ms 两形状并钉（测试体明言「recorded,
   not bounded here」）；D32 对应预算在别的腿：3s 静置窗 3039/3074ms PASS、RSS 回落
   261-263ms PASS——无预算被躲。本腿复跑 10691ms 同形状（真机时长天然漂移）。
3. **wake 路 openDelay 6.95s vs 15e 262ms＝口径差**。demo 头注与测试体明言 wake 形状
   （唤醒后才开句）非 15e 单句形状、两者⛔互当；262/640/902 零回归三数由同二进制
   Test15e 另行钉住——本腿复跑逐位同（§5）。两口径各自自洽，如实并列。

## 4 格①-④证据复核（本腿复跑读数，全在导出树）

- **格① demo 复跑**：`Test16r1VoiceDemo…RealHardware` **PASS 34.76s**——球走位五行、
  TRANSCRIPT 逐字同（真 Paraformer 对 4.wav 确定性）、tool_row fs.read L0 allow success、
  TASK_LOG done、NOTIFY 1、spoken=1 aborted=0 samples=88828（16-r1 记 88856＝真 waveOut
  抖动）；**LATENCY 行 openDelay=6.95s chunk=640ms finalLatency=3.52s speechMS=2816
  closedBy=end seamWaitToFirstText=9.99s 与 16-r1 逐位同**；settling_to_fallback=3074ms
  （3s 冻结窗内）。
- **格② SLO 附录 D 12 行逐行**：行 1/2/3/9/10/11 PASS（行 2＝704ms 两形状同值；行 10＝
  3s 窗 3000-3039）；行 6＝TTS 首帧 3.715s **FAIL 照录**（诚实 FAIL 是门在响，非缺陷）；
  行 4/5/8 未测/deferred 如实具名；行 7 首用 5.51s 超 P50（含首装）照录——超出已写明、
  单样本不成分布，无遮掩；行 12 skeleton 口径引 C.5。`git diff 83c08b2f6^..83c08b2f6 --
  docs/SLO.md`＝**0 删／32 加，纯追加**；预算行一字未动。
- **格③ 六场景映射核**：票面六＝SPEC-10 §7 三项（被占/满盘/下载中断）＋D42 #12/#2＋C29
  损坏半——逐项可追溯。抽两发复跑：**场景 5 真翻字节 PASS 32.86s**（`sha256 mismatch
  (pins 81a70226a893…, got ab9962e4c1a2…)` 与 16-r1 逐字同＝同位翻转确定性；自愈 walk
  24.07s、healed store verifies clean、`took≥5s` 断言防缓存假绿）；**场景 6 真镜像杀进程
  三发**——真 URL、`Stop-Process -Force` 杀于进程存活期、partial 存活、二跑 RC=0
  「已交还」、`models verify` RC=0「与已验签清单一致」；本机镜像过速（100ms 轮询首见即
  全量）使 CLI 级 mid-byte 窗不可得——Range 续传协议腿由 `TestResumeContinuesAtByteOffset`
  重跑 **PASS** 钉住（服务端 Range 可观测的唯一层）。场景 1/3 seam 形状与场景 2 OS 半负
  结果照 16-r1 存证核过（映射与 D37 类逐格对上）。
- **格④ CER**：CI 同款 runner（`bash tools/d22scan/runtests.sh ./internal/speech/cer/`，
  与 ci.yml:732 逐字同命令）**24/24/0、rc=0**；`near_clean` 门集（60+40 逐 wav 哈希＋
  perfect=0/empty=100% 红断言）在案；`noisy_far` EXERCISED-not-MET（15-v1 R1 具名）照录；
  **真 CI 工件今日不存在**（仓内存档 CI 日志全早于 10-26 接线）——维持 16-r1 残余定性：
  本地记录成立、真工件待 push，⛔ 以本地记录冒充 CI 工件。

## 5 零回归＋d22scan（本腿复跑）

- **六族**：首轮 27/28——`Test41a2r1WarmGateFeedsKwsNotTheAsr` wav 排空超时（192ms 有界
  通道 slow-consumer 丢帧 WARN；当时机器 ~68% 载、同行腿在飞）；**单跑 RC=0**；**全量复跑
  28/28 PASS／0 FAIL／0 SKIP、rc=0**（16v2-sixfamilies-run2.log）。15-v1 §2 同款漂移处理：
  丢帧⛔当常量、⛔当回归。
- **LATENCY 三数**：`openDelay=262ms chunk=640ms ballLatency=902ms` **逐位同**（16-r1／
  12-v3／15-v1 同值）；seamWaitToFirstText 5.51s＝墙钟漂移，不属钉死三数。
- **d22scan**：导出树 **rc=0 clean**，分母与 16-r1 逐字同（bans#1-5 internal/=287 cmd/=52、
  ban#8 internal/=730 cmd/=183）；**正控**：种真 U+2713 字符串（字节 e2 9c 93 验明）→
  **rc=1** 逐字点名 `zz_16v2_poscontrol.go:3 [emoji] …string literals included`，删后复跑
  clean（可逆）。ban#8 未扫 docs/（仪器只走 internal/cmd/design/frontend，本腿遵守）。

## 6 终判（五格逐条）

| 格 | 判 | 依据 |
|---|---|---|
| ① E2E demo | **翻（带残余：Acting 载体）** | §4 复跑 PASS＋§2 D43 逐行＋§3-1 |
| ② 延迟表 | **翻** | 附录 D 12 行逐行核＋纯追加＋FAIL 照录 |
| ③ 六场景 | **翻（边界照录）** | 映射可追溯＋场景 5/6 复跑＋场景 4 具名归票 45 |
| ④ CER | **⛔（残余维持）** | 本地 24/24 绿但真 CI 工件不存在，待 push |
| ⑤ 对抗验收 | **本件即本体** | §1-§5＝SPEC-10 §8 七项义务全落 |

**票 16 可收口**——对照 15-done 先例（第 48 票带 10 条具名残余收口）：格⑤过＝票面 AC#5
本体落地；①②③翻带具名边界；④残余挂账。翻框/改 Status＝编排者收尾动作（15-v1 同款）。

## 7 残余清单（16-r1 §7 五枚＋本腿两枚）

1. **Acting 球态载体**（升级：D43#15 工具臂＋#17-20 全仓零生产派发位）——归 26 集成面，
   建议编排者裁量开票或并入在飞 26p 范围。
2. **真 CI 工件待 push**（格④）——下一次编排者 push 后 CER 步 run 页留存即可翻。
3. **TTS 首帧 ≤1600ms（需加载）＝3.715s FAIL**——与 S0 §4 4.7s 串行 FAIL 同源；Warm 已加载
   形状归票 28 SessionScope。
4. **OS 级三触发不可达**（mic-permission ConsentStore／热拔 pnputil／真满盘虚拟盘＝admin）——
   16-r1 负结果读数在案，未冒充。
5. **场景 4 磁盘满写前预检守卫**——本体未建，票 45（`ready-for-agent`）。
6. （新）**live smoke 声学相位两发红**——本腿复跑
   `TestAC247LiveMicrophoneLevelsReachTheBallSeam` FAIL×2（sound 相 max 0.0022-0.0049 不可
   分；第二发 tasklist 见在飞 `wisp-notag.test.exe`＝同行音频并发，错峰不可得）；捕获端点
   两发皆真（Realtek 采集 376 真电平）；16-r1 当日两跑 PASS＋demo 真 waveOut 独立成立 ⇒
   环境性、不判回归、不阻收口。另：demo 头注所称 `TestLiveWasapiSmoke` 为旧昵称、真名如上
   （文档小疵，不裁）。
7. （新，承 16-r1 §6 照录）**`internal/speech` 整包在导出树构建位的 kws 路径字符串红**
   （walk-up 锚＝编译期源码路径的协议差）——本腿未重触（六族/demo 面不受累），维持照录。
