# 票 26 Path-C 集成面终判腿 `26p-v2` —— 非实现者裁决表（v1 静默死亡判死后重派）

日期 2026-11-01 10:0x-11:0x +08 · 腿 `26p-v2`（终判、非实现者）· 被裁件＝26p-r1 `4c21ac42c`
（11 枚 +2330/-10）· 量测锚＝导出树 `git archive eb1c42c26`（392/393/16-r1/382-w-v2 全在树内；
26p 写面五文件 diff 4c21ac42c..HEAD＝空，逐字节 26p-r1 原形）。派单存档
`.scratch/wisp/probes/orch/20261031/txt/26p-v1-dispatch.txt`（v1 13h 零产物判死，本腿重派）。

**方法**：仓外导出树 `D:\tmp\wisp\20261101\26pv2\tree`（git archive＋junction third_party＋
build/models 拷贝；工作树零碰）；staged 全链（A1326 三 DLL 贴测试二进制＋WISP_SILERO_MODEL/
WISP_SPIKE_MODELS/WISP_15E_TMP/WISP_43_REPO_ROOT＋cwd=包目录＋unset GOTMPDIR）；staged
`libwebrtc_audio_processing.a`＝aec3r1 成品（零重建零源改；CGO 宏集照 aec3r1 build.py 的
root-meson windows 形补齐全量——证据件 aec3r1-remeasure §7 只记两枚 -I，缺宏必假红，踩坑入
件）；工具链 go1.27.1＋mingw gcc 16.2.0。音频错峰＝每跑前进程表现核（在飞腿 wisp.test 四次
在场全让窗）。全部读数本腿现跑；日志 `D:\tmp\wisp\20261101\26pv2\logs\**`。

## 0. 一句话

**(b) 面、首格等价、反向两对（M1/M2＋还原核）全部复刻成立；(a) 面五发 6 窗 1 火（唯一火
24ms，远内 400ms）、其余 5 窗被引擎深压在 -39.35..-43.44**——392-v2 主发现「(a) 是摆动脆弱
面」在 HEAD 独立复证（红非 392 回归：26p 写面零字节；绿非稳定性质：同材料跨窗摆幅约 19.7dB
与 AC#2 双讲摆动同源）。**终判＝有条件 PASS**：(a) 验收语义按本件 §3 裁定重 key（5 窗至少
1 火＋凡火必 400ms 内＋深压分布入档），摆动残余具名移交 28-PathC／393 轨。

## 1. 三真数复刻（全部本腿现跑）

### (a) 播报中发声 ⇒ barge-in 停播 —— `Test26PBargeinStopsWithin400ms` 五发复跑

| 发 | 窗 | episode rms (dBFS) | 阈 | tighten | 判 | VAD起点→停播 |
|---|---|---|---|---|---|---|
| 1 | R1 | -42.01 | -30 | false | 不触发⇒红 :354 | — |
| 2 | R1 | -41.75 | -30 | false | 不触发⇒红 :354 | — |
| 3 | R1 | -40.06 | -30 | true | 不触发⇒红 :354 | — |
| 4 | R1 | -23.79 | -30 | false | **触发** | **24ms** |
| 4 | R2 | -43.44 | -30 | false | 不触发⇒红 :354 | — |
| 5 | R1 | -39.35 | -30 | true | 不触发⇒红 :354 | — |

- 对照：26p-r1 五轮 99/56/67/65/40ms（5/5 火）、392-r1 五轮 58/76/195/57/46（5/5）、392-v2
  三方 0/3＋0/2＋1/3；本腿 1/6。唯一火 24ms 带内——**「凡火必 400ms 内」五腿零反例；
  「每轮必火」不成立**。
- 深压带本腿 -39.35..-43.44（392-v2 记 -39.75..-45.16，同带）；对齐面每发健康（lag 26-39ms、
  strength 0.61-0.86、ready=9728）——失败全部落在引擎双讲抑制段，VAD/分段/触发/abort 链零
  异常。tighten 标志是判决时刻读数（测试件自带注记）：tighten=true 的两窗照样被压死。
- 红形与 392-v2 逐字节同形（`26p_windows_test.go:354` 超时红）；26p 写面对 4c21ac42c 零字节
  ⇒ 非 392/393/387/382/16 任何一腿回归。

### 首格 Feed 在线 vs 窗级重放 —— `Test26PFeedOnlineMatchesWindowReplay` PASS

6/6 窗在线 Feed 与离线 ObserveWindow 判决逐位同、零漂移（lag=366、strength=0.57；26p-r1
lag=644/0.89、392-v2 lag=509/0.70——会话量摆、判决等价稳）。工件 `26p-feed-equivalence.json`。

### (b) 回声整窗注入 ⇒ triggers=0、转写 Δ=0 —— `Test26PEchoNeverTriggersAndNeverTranscribes` PASS

episodes=2（残余 **-57.17/-62.72** 全被收紧档吃；392-v2 记 -51.63 同形）、triggers=0、finals
Δ=0、handed Δ=0、6 窗全 untrusted。工件 `26p-b-echo-zero.json`。

### (c) -30 裕量加严 —— `Test26PMarginRetestAndGateReverse` PASS

L1 control-26：post-AEC **-31.22（margin -1.22dB，本发被吃）**——四腿读数带 26p-r1
+0.58／392-r1 -28.69..-41.61／392-v2 -28.94..-33.60／本腿 -1.22 ⇒「offline 3.3dB 裕量不保证
过链」再证、摆动是主导项；L2 quiet-36.58：-42.33 被吃、标称轴可达；L3 失配回声：-34.91
（margin -4.91dB）动态门 held。工件 `26p-c-margins.json`。

## 2. 反向两发复刻（全带还原核）

- **M1 摘组合门（dynamicOff ⇒ 回 -45 松轴）⇒ 自触发复现**：同一失配回声材料 -36.91 过松轴
  ＋播报自我中断（无说话者）＝11.5/10min 族在集成面复现；**还原核**：门回 ON 同材料 -35.26
  被 -30 收紧档 held、triggers=0。成对闭合。
- **M2 摘 AEC3（aecBypass）⇒ 逐字转写复现**：裸回声开 episode（-37.4）＋同材料过真 ASR 链
  （silero＋流式 Paraformer）逐字转写「不已经完成明天下午三点的会议改到四点勤准时参加谢谢」
  （distinctive 命中；与 392-v2「(不)已经完成…」同为 ASR 抖动级）；**还原核**：AEC 回来
  episodes=0、triggers=0。成对闭合。工件 `26p-m2-verbatim.json`。

## 3. (a) 加稳语义裁定（派单新增；铁律④裁量：最有价值，不用风险/动作当判据）

| 形 | 内容 | 代价 | 为什么（不）最有价值 |
|---|---|---|---|
| **重 key** | (a) 断言改 episode keyed：K 窗（建议 5）至少 N 火（建议 1）＋凡火必 400ms 内＋每窗深压读数（rms/阈值/tighten）必录 | 只改测试形状（26p 测试件），零产品码、零契约文字；须照 392 的 FINDING MOVED 牙式把旧判语钉死，防回退 | **裁定选它**：两条 S4 字面（凡火 400ms 内＋零自激）逐 episode 原样保住；N 地板防空转死（全压死仍红）；摆动从掷币变显式受测维度。主流 harness 对环境敏感验收面的做法＝事件 keyed＋有界尝试＋分布入档 |
| **加稳** | 塌峰检测（389 三件套已在跑）＋动态退避：公告窗内塌峰报警时把收紧档向标称轴回退（或 VAD 佐证重判） | 产品码改 monitor（＋可能动 389 门常量面）；M1 即直接危害——本腿 M1 实发 -36.91 证明松轴带被回声伪迹占着，退避必须配自己的反向对；根因在 AEC3 suppressor 状态机（392-r1 三旋钮无一拉平、393-r2 Reset 移除 carry 退化），monitor 侧退避治标 | 真产品价值但**不是验收腿的活**：应随 28-PathC 集成（26p-r1 残余①已把双讲调参划给 28）且踩在 393 Reset API 上做；验收腿内改产品＝裁/实现混角色 |
| **放宽** | 验收语义改分布读数（如 P50 400ms 内或仅记录） | 削 S4 行为字面（PLAN:1688「400ms 内 TTS 停播转听」）；改验收文字＝契约变更须人工批（SPEC-12 §4.1） | **驳回**：重 key 已买到同样的稳定且不丢硬界；放宽是纯减保 |

**裁定＝重 key**：26p (a) 的验收语义按「5 窗至少 1 火＋凡火必 400ms 内＋深压分布入档」收口
（K=5/N=1 起步；393 Reset 落地后重测若分布收窄可上调 N）。加稳移 28 轨；放宽驳回。

## 4. 零回归复核

- **六族 staged 合跑 33 PASS / 2 FAIL / 0 SKIP**——两红 `Test15eSpeechEndToEndHandsThe…`
  （10.99s，球投影 [Listening] vs [Listening Thinking] 取样竞速）与
  `Test384ResidentWatchdogEndToEnd`（3.30s 合跑负载）**单跑全 PASS**（8.23s / 1.16s）＝负载
  时窗归因（A1329 先例形状）；总数 35 对齐基线、红名集合对 34/1/0 零新增。**`262/640/902`
  逐位复现**（`LATENCY openDelay=262ms chunk=640ms ballLatency=902ms`，单跑读数）。
- `internal/audio/...` 双形全绿：no-tag（audio 17.5s＋aecspike 0.2s）＋`-tags webrtcap`
  （audio 17.4s＋aec 0.053s 含 393-r2 Reset 套件＋aecspike 0.2s）双 rc=0。aecspike 需 DLL
  搜索序：PATH 前置三 DLL 目录；Git Bash 下 cgo 无 gcc＝环境假红（setup failed），pwsh 复跑
  全绿——仪器事项非码面。
- 两尺：gofumpt `-l`（cmd/wisp＋internal/audio 全树）＝0（离线 v0.12.0，GOPROXY=off 直跑
  模块缓存目录）；vet 无 tag／webrtcap 双 rc=0。
- **d22scan clean rc=0**：bans#1-5 internal/=287、cmd/=52（对上）；ban#8 internal/=**730**、
  cmd/=**183**——对派单 729/180 的增量逐枚归因：internal +1＝`aec393r2_reset_windows_test.go`
  （393-r2）；cmd +3＝`resident_voice_bargein_392_windows_test.go`（392-r1）＋
  `resident_voice_ttsinherit_382wr2_windows_test.go`（382-w-r2）＋
  `voice_e2e_demo_16r1_windows_test.go`（16-r1）。**零未解释增量**。selftest 35/0/0（RUN=77）。

## 5. 终判：Path-C 集成面＝**有条件 PASS**

对照 59-v3 §6 重开条件字面：组合门＋动态门控落 internal/audio（389-r1）成立；-30/-36 裕量
重测加严（26p-r1＋本腿＋392 两腿读数带）成立；400ms 停播链实测（26p-r1 五轮＋本腿唯一火
24ms）成立——「差什么：什么都不差」兑现。S4 双字面（PLAN:1688）：「不自激」稳定成立（本腿
＋三前任四腿零反例＋M1/M2 反向对闭合）；「400ms 内停播」凡火必达，唯「必火」是引擎态敏感
读数（本腿 1/6 实证、392-v2 主发现独立复证）。

**条件两条**（缺一不作无条件 PASS）：
1. **(a) 验收语义按 §3 裁定重 key**（小写腿改测试形状＋FINDING MOVED 牙；不改产品、不改
   契约文字）。
2. **摆动残余具名移交**：AEC3 双讲深压（跨腿带约 -28.1..-45.16）归 28-PathC 集成（调参/
   退避踩 393 Reset API）＋393 flush/reset 轨——26 面不背、也不阻塞。

NO-GO 不成立：(b)／首格等价／反向对零回归且稳定，失败模式是具名引擎态属性、修复路径已在
轨；无条件 PASS 不成立：「五轮全 400ms 内」在四腿三方读数下不是可复现保证。

## 6. 26-PathC 收官与 26/28 接力裁定

- **26 收官**：按条件 PASS 结票——(a) 重 key 小写腿是最后一步（或并入 28 首切片同文件）；
  其余面已收官。本判后 26 不再有独立测量需求（59-v3「不再需要独立测量腿」在本面同样成立）。
- **接力**：59-v3 给 28 的维持条件＝「26 的门模块落地＋首格验收过」——两者已成事实（389-r1
  在树、首格等价＋(b) 四腿稳定）。28-PathC（票 392 会话集成切片）设计输入已消费（保温测死
  默认关、调参面无免费解、dispose 级联）⇒ **28 维持第二顺位、正式进集成**；26 转移件＝双讲
  摆动（§3 加稳形）＋保温深度，落点 28＋393 Reset。393-v1 已裁「Reset 即统一器、保温修复
  无需动生产路径」⇒ 28 集成保温用 Reset，不必再动 monitor 生产路径。

## 7. 写面申报与工件

- 仓内唯一写面＝本件。禁 push/amend/reset/rebase/stash/checkout；工作树零改（实验全在仓外
  导出树）；在飞 16-v2／382-w-v2 写面零触碰、零归因。
- 仓外 `D:\tmp\wisp\20261101\26pv2\**`（tree＝eb1c42c26 导出＋third_party junction、bin/
  bin-notag 双二进制＋三 DLL、logs 全量、msg、out-backup-26pr1、export.tar），只建不删。
- **工件溯源**：26p 测试硬码 out 根 `D:\tmp\wisp\20261031\26pr1\out`——本腿四件 JSON 现为
  **本腿读数**；实现件/392-v2 原读数备份 `D:\tmp\wisp\20261101\26pv2\out-backup-26pr1\`
  （392-v2 同款处置）。
- 锚点申明：量测锚 eb1c42c26；判后 HEAD 已走 56f94624b（16-v2 落件）——26p 写面零字节核对
  含 56f94624b（diff 4c21ac42c..56f94624b 对五文件仍空），本判不重锚。
