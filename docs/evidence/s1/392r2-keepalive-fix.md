# 票 396 落地腿 `392-r2` —— 保温修复链收口（毒源半隔离＋Reset 接入保温）（写腿）

日期 2026-11-01 · 腿 `392-r2`（写、`internal/audio/**` 独占＋本件）·
起手 HEAD `6bc014439`（dev，开工三核过：①锚点在 ②索引净——唯一 tracked 脏件
`.scratch/wisp/issues/30-result-routing-d10-done.md` 在他人面、零触碰，写面
`internal/audio/**` 零脏 ③进程表现核无 rig 在飞）· 派单存档
`.scratch/wisp/probes/orch/20261101/txt/392-r2-dispatch.txt` 在案。

## 0. 一句话

**毒源钉死＝detector 半**（不是 resampler）：392-r1 §6.1 的逐件隔离矩阵在
monitor 判官链的确定性复刻上跑完，保温判死＝**中断把 endpointer 的开口话态带进
下一暖窗**（开口话态不关闭就不可能再发 start ⇒ episodes=0），保温 resampler 的
远端输出被逐位洗清（残余 -62.11 与 cold 全同）；392-r1 finding 的「残余骑近端
电平」一半归属**改判到 APM 跨窗携带形**（裸携带 APM 的下一窗残余 -42.78 ≈ 近端
-38.17、只消 ~5dB，与 392-r1 candidate-2 的「彻底失准」同形；candidate-3 的每窗
新 APM 形不存在该现象）——392-r1 的两半在原 evidence 里是并成一句写的，本腿把
它们拆开各归各。**AC#2 复活＝「持有＋窗间清」纪律**：保温池每件在窗间走自己的
公开清面——APM `Reset()`（393-r2 的面）＋resampler `Flush()`（documented 面）＋
detector `Finalize(nil)`＋**1.2s 静音冷段**（新发现：silero RNN 语境跨空窗携带、
hangover 能在静音上开火，冷段把它吸收）——W2/W3 全愈（triggers=1@talker、残余
cancelled 深度）；反向（摘 Finalize ⇒ finding 复现；摘 Reset ⇒ 残余死 19.24dB）
＋还原核＋两跑逐位全过。零回归：四数 `-v` 口径 **RUN=98 PASS=75 FAIL=0 SKIP=11**
与 392-r1 自记逐位同、SKIP 集同形；`keepalive=false` 默认面 **0 字节**
（`cmd/wisp/**` 与 `internal/audio/aec/**` 双零 diff——393 Reset 语义未碰）。

## 1. AC#1 毒源半隔离（392-r1 §6.1 具名残余的执行）

**仪器**：monitor 判官链的包内确定性复刻（`internal/audio/aecspike/
keepalive392r2_windows_test.go`，门 `WISP_AEC392R2=1`；393-r2 先例：原测量在
`cmd/wisp`、本腿写面进不去，故在包内自足复刻同一测量形状）。链＝真 matcha
语声公告（合成 chunking 过流式 resampler 远端 tap）→ 26p 回声（延迟 480、耦合
0.2）＋ -20dBFS talker@400ms 的 512 样本 seam 帧 → anchor/armJudge（包络相关
测 lag）→ APM3 render/capture lockstep（delay hint＝lag/16 ms）→ 收紧档
BargeinGate.Feed → 残余 noiseGate（floor/step/hold 常数照抄 monitor）→ silero
分集 → Evaluate 判火。全样本计数驱动、无墙钟、无 goroutine、无音频硬件。

**矩阵**（8 臂 × W1/W2/W3 × 两跑；W1 建引擎、W2 判别、W3 持久）：

| 臂 | det | resampler | APM | W2 读数（R1=R2 逐位） |
|---|---|---|---|---|
| cold-control | 每窗新 | 每窗新 | 每窗新 | episodes=1 triggers=1@12320，残余 -62.11（近端 -38.17） |
| warm-392（392-r1 candidate-3 原形） | 保温裸 | 保温 | 每窗新 | **episodes=0 triggers=0**（finding 复现），**残余 -62.11＝cold 全同** |
| iso-res-fresh（仅 det 保温） | 保温裸 | 每窗新 | 每窗新 | **episodes=0 triggers=0**——毒与 resampler 无关 |
| iso-det-fresh（仅 resampler 保温） | 每窗新 | 保温 | 每窗新 | triggers=1，**残余 -62.11＝cold 全同**——保温远端输出干净 |
| revival（Finalize/Flush/Reset） | 保温+清 | 保温+清 | 保温+Reset | episodes=1（判官复活）但 triggers=0（§2 hangover 瞬态），残余 -62.02 |
| revival-cool（＋冷段） | 保温+清+冷段 | 保温+清 | 保温+Reset | **episodes=1 triggers=1@12320**，残余 -62.02 |
| revival-no-reset | 保温+清+冷段 | 保温+清 | 保温裸携 | triggers=0，**残余 -42.78**（≈近端 -38.17） |
| revival-no-finalize | 保温裸 | 保温+清 | 保温+Reset | **episodes=0 triggers=0**（finding 复现） |

仪器健康：8 臂 W1 全部逐位同（每臂第一窗＝全新引擎同料）；cold 三窗逐位同
（剔除窗口序号后 JSON 全等）；两跑逐位同；还原核（全矩阵跑完后再读一冷窗）
与开场冷窗逐位同。

**机制单元探针**（`Test392R2DetectorOpenUtteranceCarries`）：talker 开口
（MinSpeech 300ms 经 VAD）→ 400ms 静音（< 700ms 关门阈）→ 第二 talker ⇒
**零新 start**、`Finalize` ⇒ EndOfStream 且 Total 为跨两段累积——开口话态跨窗
携带被直接读出并直接摘除；对照（Finalize 后）start 复现。**附带实测发现**：
对照臂的 start 落在**静音段**（silence=1 talker=0）——silero RNN 语境在空窗期
不衰减（无输入即无更新），hangover 能在 post-speech 静音上喂出 speech run
开火 ⇒ §2 冷段的由来。

**判读（AC#1 收口）**：毒源＝**held detector 的跨窗态（开口话态）**；
resampler 半洗清。392-r1 finding 拆判：①「episodes=0、不判」＝detector 半
（本矩阵复现）；②「残余剖面骑近端电平 -22~-25dBFS」＝**APM 跨窗携带形**
（revival-no-reset 在链上复现：-42.78 ≈ 近端 -38.17；candidate-3 的每窗新
APM 形残余全程 cancelled 深度、不出现该现象）——392-r1 该句的残余读数在其
工件里本就是空壳（round 结构体字段未导出、JSON 序列化为 `{}`、剖面从未落盘），
本腿读数把这两半拆开各归各。

## 2. AC#2 Reset 接入保温（keepalive 显式开启读数）

**复活形状＝「持有＋窗间清」**：393-v1「Reset 即统一器」从 APM 半推广到整池——
每件保温引擎在窗间走自己的公开清面，无一处动生产路径：

- **APM**：保温池持有＋窗间 `Reset()`（393-r2 落的面，reset==rebuild 逐位）。
  读数：残余 -62.02 vs cold -62.11（差 0.09dB＝393-v1 §6 两引擎形在链上的
  显影，无害）；**Reset 承重 19.24dB**（§3 反向）。
- **resampler**：保温池持有＋窗间 `Flush()`（`resample.go` 自documented 的
  清态面）。读数：与 fresh resampler 逐位同带（iso-det-fresh 残余全同 cold；
  warm 未 Flush 的相位延续把 lag strength 从 1.000 挪到 0.999、lag 仍 480＝
  无害）。
- **detector**：保温池持有＋窗间 `Finalize(nil)`（公开面：关掉中断留下的开口
  话态、endpointer 回 idle、silero session 仍保温）**＋ 1.2s 静音冷段**——
  RNN 语境跨空窗携带、hangover 在静音上开火（探针实测 silence=1），Finalize
  只清 endpointer 清不掉 RNN；冷段让语境衰减、把第一暖窗的假 start 吸收掉。
  冷段中产生的边界事件一律弃置（非窗读数）。

**读数（keepalive 显式开启）**：revival-cool W2 与 W3 全愈——episodes=1、
triggers=1、trigger_at=12320（talker 本体，与 cold 同位）、残余 -62.02
（cancelled 深度）。**transient（具名、未饰）**：只 Finalize 不冷段的
revival W2＝episodes=1、triggers=0——判官活了但 hangover 的假 start 抢先
吃掉 endpointer、talker 的 onset 开不了口；该假 start 在本材料的收紧档下
不过阈（残余骑在收敛头，≤ -30）＝无假打断、但该窗的 talker 打断会漏。冷段
补齐后瞬态消失。

## 3. AC#3 反向判据（三读数＋还原核）

1. **摘 detector 清面**：revival-no-finalize W2＝episodes=0、triggers=0——
   finding 回归，`Finalize` 承重。
2. **摘 APM Reset**：revival-no-reset W2 残余 -42.78 vs 复活臂 -62.02
   （退化 19.24dB）、近端 -38.17——裸携带 APM 消到近零（392-r1 candidate-2
   「残余≈近端电平」形状在链上复现），`Reset` 承重。
3. **冷姿势对照**：cold-control 三窗全判、echo-only 残余 -62.11（比近端低
   24dB）＝仪器基准；warm-392 W3 自愈（W2 跑满全窗 ⇒ 尾部静音超 700ms ⇒
   携带话态关闭）＝毒源是开口话态的**寿命**而非永久性腐坏——与机制自洽，
   也与 392-r1 测试只测到 warm2 的形状一致。
4. **还原核**：全矩阵＋两跑之后全新冷窗复读开场 cold W1——**逐位全同**。

## 4. AC#4 零回归

- **四数**（无 tag，`PATH="$PWD/third_party/sherpa-onnx:$PATH"`，
  `-count=1 -v`）：rc=0，两包 `ok`＋aec 无测试文件；**-v 口径 RUN=98
  PASS=75 FAIL=0 SKIP=11**——与 `392r1-session-integration.md:108` 自记口径
  逐位同；SKIP 集＝TestLiveWasapiSmoke＋aecspike 十枚硬件/工件门（日志中一枚
  门名被密掩占位符替换，枚数与集形与 393-v1 记录吻合）。本腿测量件
  （`windows && cgo && webrtcap` 三重门）在无 tag 编译域外 ⇒ 无 tag 宇宙
  逐字节不动。
- **keepalive=false 默认**：`cmd/wisp/**` 本腿 **0 字节**（仅读）；找
  `keepalive` 语义的载体 `internal/audio/aec/**` 亦 **0 字节**（393 Reset
  语义未碰）——默认翻转与生产姿势改动都属后续裁定（392-r1 裁定在案）。
- **build/vet**：tagged build rc=0、tagged vet rc=0、无 tag vet rc=0、
  GOOS=linux vet rc=0。
- **gofumpt**（`D:\work\base\gopath\bin\gofumpt.exe`）：`-l internal/audio/`
  ＝空。
- **d22scan**（Git Bash `scripts/d22scan.sh`）：clean rc=0；selftest
  PASS=35 FAIL=0 SKIP=0（RUN=77）；分母 bans#1-5 internal/=287、cmd/=52、
  ban#8 internal/=731（=393-v1 的 730＋本腿测试件 1）、cmd/=183。

## 5. 写面申报（commit pathspec 逐枚）

仓内：
1. `internal/audio/aecspike/keepalive392r2_windows_test.go`（新：AC#1 矩阵
   ＋AC#2 复活＋AC#3 反向＋机制探针，8 臂×3 窗×2 跑＋还原核，门
   `WISP_AEC392R2=1`、工件 `WISP_AEC392R2_OUT`）
2. `docs/evidence/s1/392r2-keepalive-fix.md`（本件）

仓外 `D:\tmp\wisp\20261101\392r2b\**`（测试 exe＋三枚 sherpa DLL、logs\
{vet-tagged,build-tagged,run1..4,four-numbers,vet-notag,vet-linux,d22scan}.log、
out\392r2-keepalive-matrix.json），只建不删。

⛔ 越界核：`cmd/wisp/**` 零触碰（26p-k1 在飞）；`internal/audio/aec/**` 零
触碰（393 Reset／392 Tuning 语义原样）；`internal/speech/**` 零触碰（detector
的 Finalize/Accept 都是既有公开面）；AEC3 源树零触碰（钉定树只读）；389 门
本体（`internal/audio/bargein.go`）/thresholds/golden/D43 表/票面/台账/specs/
PLAN/deps.toml/testdata 仓根/frontend/design/scripts/.github 零触碰；⛔ push
（推送归编排者）。

## 6. 残余与边界（如实记账）

1. **monitor 采纳路径（给编排者/下一腿）**：monitor 的保温姿势复活＝三处
   窗间清（`release()` 的 keepalive 分支：APM 持有＋`Reset()`——或维持每窗
   重建、两者读数同健康，Reset 形是统一器兼省构造；resampler 持有＋`Flush()`；
   detector 持有＋`Finalize(nil)`＋冷段、边界事件弃置）。全部公开面、零新
   API。**默认仍 `keepalive=false`**——翻转属后续裁定，本腿读数是裁定材料。
2. **RNN hangover 是 detector 面的具名缺口**：silero session 无 reset 面
   （`internal/speech` 无此 API），冷段是 owner 侧纪律不是 session 清面；
   若未来要更硬的清面（RNN 重置）属 internal/speech 票。 hangover 的假
   start 在 26p 回声电平（≈-32.6dBFS）下不过收紧档＝无假打断，但该窗
   talker 打断会漏——冷段吸收。
3. **跨进程材料微差**：matcha 渲染跨进程读数差 ~0.2-0.6dB（onnx 级），
   本仪器一切判据是**类别判据或相对判据**，两跑逐位指**进程内**两跑；跨进程
   逐位复现不成立也不需要（393-r2 的全合成刺激可以跨进程逐位，本腿材料含
   真 TTS，做不到，如实记）。
4. **392-r1 工件空壳**：`392-keepalive.json` 的 round 对象全是 `{}`（字段
   未导出），其 finding 的剖面读数无盘上支撑；本腿工件全导出字段、剖面在案
   （`392r2-keepalive-matrix.json`，40KB，两跑全量）。
5. **两引擎形在链上 0.09dB**：Reset 形（Initialize-locked 分带形）与每窗
   新 wap_new 形（全带形）的 393-v1 §6 差异在本链读出 0.09dB（-62.02 vs
   -62.11）——无害，未动生产路径，统一引擎形仍须先跑 392 K=4 对照（393-v1
   处置建议原样）。
6. **停止条款核**：①毒源两件都隔离不掉——未触发（detector 单独承毒、
   resampler 单独洗清，矩阵双向判明）；②Reset 接入后保温仍判死——未触发
   （复活臂 W2/W3 全愈）；③必须动 keepalive 默认——未触发（默认 0 字节）。
