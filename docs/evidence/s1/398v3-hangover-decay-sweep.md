# 票 398 残余③执行腿 `398-v3` —— silero RNN hangover 跨空窗自然衰减扫参（写腿·纯测量）

日期 2026-11-01 · 腿 `398-v3`（写、**纯测量**：零产品码、零既有断言改动；写面＝测量件＋本件）·
起手 HEAD `afe93c5a7`（dev，开工三核过：①HEAD 在 dev＝派单门槛本身 ②索引 0 staged、
tracked 零改动 ③`internal/speech/vad_decay398v3_test.go` 不存在、`docs/evidence/s1/398v3-*`
为空＝写面无他人脏）· 派单存档 `.scratch/wisp/probes/orch/20261101/txt/398-v3-dispatch.txt`
在案 · 现跑 2026-11-01 20:1x–20:5x +08。

> **环境**：Git Bash runner（`D:\work\soft\Git\bin\bash.exe`）；`unset GOTMPDIR` 每跑；
> rc 一律重定向进 log（`[rc=N]` 尾行）；gofumpt 绝对路径 `D:\work\base\gopath\bin\gofumpt.exe`。
> **staged 形状**（387-r3/398-r1 §头部同款、A1326 卡）：`go test -c -o <tmp>\staging\*.test.exe
> ./internal/speech` ＋ 三 DLL 贴身（SHA256 与 third_party 逐位对上，`dll-sha256.txt`）＋
> `WISP_SILERO_MODEL`/`WISP_SPIKE_MODELS` 两 env ＋ `-test.v -test.timeout 30m`。
> **裸跑形状**：`go test -count=1 -v ./internal/speech`、PATH 注入 `third_party/sherpa-onnx`、
> 无 DLL 贴身、无 `WISP_*`。★**环境陷阱（本腿第一跑即撞、已具名修复）**：MSYS2 给原生子进程
> 转 PATH 时把 `D:/...` 形条目绞碎（实测 cmd 读到 `D;D:\work\soft\Git\work\workspace\...`），
> 测试进程拿不到 sherpa DLL ＝ `exit status 0xc0000135`；**PATH 注入必须用 POSIX `/d/...` 形**
> （探针 `probe-path.sh`/`probe-path2.sh` 在案：D:/ 形两式全败、/d/ 形 PASS rc=0）。
> **音频错峰**：跑前现核无 test/wisp/rig 进程在飞。**临时件（只建不删）**：
> `D:\tmp\wisp\20261101\398v3\{staging,bare,logs}\**`。**⛔ push；⛔ 票面/台账/specs/PLAN/
> thresholds/golden/testdata 仓根/frontend/design 零字节；⛔ cmd/wisp、internal/{risk,session,
> audio} 零字节。**

## 0. 一句话

**衰减扫参落地、停手三条全未触发。** 自然臂 N=0..16＋每空窗 Reset 归零对照臂 N=1..16，
双跑逐位稳定（去时间戳后零差）；曲线三点位：**①假 start 恒只落边界后第 1 个空窗**（全部
深度 `decayGapStarts=[1 0 ...]`，第 2 空窗起零开火）；**②RNN 开火自然衰减点＝深度 1**
（400ms 馈入静音后探针即 `(silence=0 talker=1)`，398r1 控制形 `(1,0)` 仅在深度 0）；**③假
话语自然关门＝第 4 空窗**（1600ms；`leftOpen` 在 n=3→4 翻 0）；归零对照臂全深度平零
`(0,1)`。喂 392-v3 采纳③的量化裁定：**免冷段（Finalize+Reset）＝零空窗需求；自然衰减
RNN 半最省 400ms、但假话语半需 1600ms 或一次 Finalize**——「衰减多少够」（392-v3 §7.4
原样）就此关闭。零回归：本腿自跑基线 staged **69/63/0/6**＋裸跑 **69/25/0/44**（与
398-r1/398-v2 在案逐位同）；落件后 staged **71/65/0/6**＋裸跑 **71/25/0/46**（恰两枚具名
新增/skip、既有 69 枚名册逐位同序同判）；gofumpt/vet/d22scan clean（selftest 35/0，
RUN=77）。

## 1. ★AC#1 扫参设计（398-r1 探针形状复用）

**单 trial 形状**（深度 n，两臂同形、只差 Reset 插入点）：

```
talker 1500ms → Finalize(nil) → [归零臂：Reset] gap 400ms × n → Finalize(nil)
  → probe gap 400ms → probe talker 1500ms
```

- **读数四列**（392-r2 具名形 `(silence=N talker=M)` 为主读数）：`decayGapStarts[k]`/
  `decayGapEnds[k]`＝衰减馈入期第 k 个空窗内的 start/关门数（假 start 落点、假话语自然
  关门点）；`leftOpen`＝探针前那次 Finalize 还须关几个话语（**衰减馈入结束时假话语是否还
  开着**）；`probe (silence/talker)`＝探针空窗上的 start（hangover 仍活着）＋探针 talker
  上的 start（下一口语声正常开口）。
- **探针前的 Finalize 只关话语、不碰 RNN**（vad.go `Finalize` 零字节于 vad——398-r1 钉过
  的正交面），所以 probe 读数隔离出「N 个空窗后 RNN 还剩多少语境」这一单变量。
- **派单「每空窗后跑探针」的实现读法（如实申报）**：探针按**每深度 N 的空窗序列末尾**各跑
  一次（每 trial 全新 detector、单变量归因），不是把探针插进同一条空窗流——探针语声若落在
  空窗之间会把 RNN 重新喂上语声语境，被测的衰减本身即被销毁（设计取舍在案；曲线仍按深度
  逐点给出，与「逐空窗读数」同表同形）。
- **复用点**：`count398Starts`（398-r1 计数形）、1500ms 代理语声＋400ms gap（vad_test.go
  实测校准件，一字未改）、`(silence=N talker=M)` 具名读数、`newTestDetector`（守卫家族）。
- **断言政策（纯测量腿）**：仅两枚测量有效性守卫——①subject 守卫（trial 语声开不了口＝
  无被测对象，398-r1 同名形）；②归零对照守卫（归零臂任一深度开火或探针 talker 哑＝对照失
  效/Reset 触分类器，`MECHANISM MOVED` 具名句）。**曲线本体零断言**——衰减点位是测量产出
  不是机制钉子，只入 log 与本件。

## 2. ★AC#2 衰减曲线（读数）

自然臂（无 Reset、log `sweep-run1.log`；n=5..16 与 n=4 同形，全表在 log）：

| n | decayGapStarts | decayGapEnds | leftOpen | probe (silence/talker) |
|---|---|---|---|---|
| 0 | （无空窗） | （无空窗） | 0 | **(1, 0)** ＝398-r1 控制形逐字 |
| 1 | `[1]` | `[0]` | 1 | **(0, 1)** |
| 2 | `[1 0]` | `[0 0]` | 1 | (0, 1) |
| 3 | `[1 0 0]` | `[0 0 0]` | 1 | (0, 1) |
| 4 | `[1 0 0 0]` | `[0 0 0 1]` | **0** | (0, 1) |
| 5..16 | `[1 0 …]` | `[0 0 0 1 0 …]` | 0 | (0, 1) |

`398v3-DECAY-POINT n=1 (400ms of fed silence)`（测试自派生行）。归零对照臂（每空窗前
Reset）：**n=1..16 全深度 `decayGapStarts` 全零、`leftOpen=0`、probe (0,1)**——平零曲线
（log 同件）。**双跑逐位稳定**：`sweep-run1.log` vs `sweep-run2.log` 去时间戳后零差。

**三个点位（本腿曲线的读法）**：

1. **假 start 恒落第 1 个空窗**：n≥1 全部深度 `decayGapStarts[0]=1`、其余恒 0——不 Reset
   时，边界后第一条馈入静音**必然**吃掉一次假 start；自然衰减只阻止**后续**开火，不阻止
   第一次。
2. **RNN 开火自然衰减点＝深度 1**：第一个 400ms 空窗过后，新探针空窗不再开火（probe
   silence=0 从深度 1 起）。粗粒度注记：扫参分辨率＝400ms 空窗，衰减在 (400,800]ms 内完成，
   与 vad_test.go 在案「~720ms hangover」同带；更细刻度不在本腿范围。
3. **假话语自然关门＝第 4 空窗**：假 start 开出的话语在 1200ms 馈入静音时仍开着（n≤3
   `leftOpen=1`），1600ms 内自然关（n=4 `decayGapEnds[3]=1`）——它比 RNN 语境**活得长**，
   因为 700ms 关门静音要从 hangover 死透后才开始计。

**398-r1 形状锚（停手③核）**：n=0 行＝`(1,0)` 与 398-r1 PRE-RESET 逐字同；归零臂 n=1 行
＝(0,1) 与 398-r1 POST-RESET 同形（探针多一空窗、读数同）。

## 3. ★AC#3 裁定材料（喂 392-v3 采纳③「keepalive 默认翻转＝后续裁定」）

| owner 选项 | 本腿读数 | 下一口语声 | 空窗需求 |
|---|---|---|---|
| 无冷段、无 Reset（现状裸形） | n=0 probe **(1,0)** | **被吞**（假话语开着） | —— hazard |
| 自然衰减（Finalize 后只喂静音） | 深度 1 起 probe (0,1)；n≤3 `leftOpen=1`、n≥4 `leftOpen=0` | 探针形（含隔离 Finalize）深度 1 即开口；**无后置收口则需 ≥4 空窗** | RNN 半 **400ms**；假话语半 **1600ms** 或一次 Finalize |
| Reset 面（每空窗 Reset；398-r1 钉过单次 Reset≡fresh） | 全深度 (0,1)、零开火 | 正常开口 | **0**（Finalize+Reset 组合即完整窗间清） |

**口径差如实记（无矛盾、要分清）**：392-v3 在案「1.2s 冷段全愈」读数出自 aecspike rig
（`边界事件弃置` 在案）；本腿读的是 bare `Detector`——1.2s（3 空窗）在 bare 形状下假话语
**仍开着**（n=3 `leftOpen=1`），自然关门要到 1600ms。差别的承重面＝**假话语由谁关**：rig
侧的边界弃置/收口把这一半吃掉了，bare Detector 把两半拆开量化（RNN 半 400ms、话语半
1600ms 或一次 Finalize）。这不是 1.2s 读数翻案，是把「衰减多少够」从一点读数补成曲线、
并把 Finalize 在时序里的承重具名化。

**结论句（裁定材料原样）**：免冷段的量化依据成立——不 Reset 时 RNN 开火要 400ms 馈入
静音才停、假话语要 1600ms 或一次 Finalize 才收口、且第一条馈入静音必吃一次假 start；
`Finalize(nil)＋Reset()` 组合零空窗归零且正控（probe talker=1）在全部深度成立。

## 4. ★AC#4 零回归（本腿自跑）

| # | 跑 | 形状 | rc | 读数 |
|---|---|---|---|---|
| 1 | 基线 staged（落件前干净树，`speech398v3base.test.exe`） | 三 DLL 贴身＋两 WISP env | **0** | **RUN=69 / PASS=63 / FAIL=0 / SKIP=6**；6 SKIP＝punc 缺席族逐名同（`suite-staged-base.log`） |
| 2 | 基线 裸跑 | POSIX PATH 注入、无 WISP_* | **0** | **PASS=25 / FAIL=0 / SKIP=44 / RUN=69**（`suite-bare-base.log`） |
| 3 | 落件后 staged | 同 1 ＋本腿测量件 | **0** | **RUN=71 / PASS=65 / FAIL=0 / SKIP=6** |
| 4 | 落件后 裸跑 | 同 2 | **0** | **RUN=71 / PASS=25 / FAIL=0 / SKIP=46** |

- **名册尺**（`roster398v3.sh`，名＋判逐位 diff）：staged 基线 vs 落件后、裸跑基线 vs 落件后
  两对 diff 均**恰两行新增**（`Test398v3HangoverDecaySweep`/`Test398v3DecaySweepZeroedControl`
  按件名字序插在 `Test398SileroHangoverReset` 之前），**既有 69 枚逐位同序同判**——398-r1/
  398-v2 名册 ⛔ 破未触发（判据＝FAIL=0 且 rc=0 且既有名册逐位同，三条全过）。
- **裸跑 +2 SKIP 守卫句**逐字同守卫家族（「onnxruntime 1.17.1 is older than the 1.28…」
  System32 影子，在案形状）。
- **门禁**：gofumpt `-l internal/speech/`＝0 行 rc=0；`go vet ./internal/speech` rc=0；
  d22scan clean rc=0——**ban#8 internal/=735**（HEAD 在案 734＋本腿测量件 1，逐源可数；
  cmd/=184 同在案值）；selftest 正控 **PASS=35 / FAIL=0 / SKIP=0（RUN=77）** 与在案逐位同。
  本证据件不入分母（ban#8 不扫 docs/，392-v3 在案口径）。

## 5. ★停手三条核销（派单末三条）

1. **「16 空窗内不衰减也不开火（曲线平零）」**——未触发：假 start 恒落 gap1、探针深度 1
   即 (0,1)，曲线信息量足；如实记，无「免冷段更强证据」情形。
2. **「测量件必须动既有断言」**——未触发：新件一枚，`vad_reset398_test.go`/`vad_test.go`/
   产品码零字节（`git status --porcelain -- internal/speech/`＝恰一枚 `??`）。
3. **「扫参曲线与 398-r1 探针形状对不上」**——未触发：n=0＝(1,0) 逐字、归零臂 (0,1) 同形
   （§2 形状锚）。

## 6. 写面申报（commit pathspec 逐枚）

1. `A internal/speech/vad_decay398v3_test.go` —— 唯一代码件（纯测量：两案＋一 trial 助手；
   复用 398-r1 计数形与校准件；断言仅两枚测量有效性守卫）
2. `A docs/evidence/s1/398v3-hangover-decay-sweep.md` —— 本件

## 7. 残余与边界（如实记账）

1. **判据口径**：probe 读数含探针前隔离 Finalize（RNN/话语两半拆开读）；owner 若**无**后置
   收口，其等价读数＝`leftOpen` 列（§3 口径差记）。
2. **衰减刻度＝400ms**：RNN 半的衰减点落在 (400,800]ms 内、未再细分；fixture（代理语声/
   采样率/模型）换代须重扫。
3. **1.2s vs 1600ms 口径差**已具名（§3）——若验收腿要 rig 侧复读（cmd/wisp staged 全链），
   属 392-v3 采纳③第 3 步翻转前置那一腿的写面，本腿⛔越。
4. C++ 源码级核仍缺（网络策略拦 GitHub，398-r1/398-v2 在案原样）——本腿全部为行为级读数。
5. 环境陷阱遗产：MSYS2 `D:/` 形 PATH 条目绞碎 ⇒ 0xc0000135；后续裸跑一律 POSIX `/d/` 形
   （§头部、探针在案）。

## 附：本腿自跑尺（全部可重放）

`D:\tmp\wisp\20261101\398v3\`：`run398v3.sh`（baseline/build/sweep/suite/gates 五模式）、
`roster398v3.sh`（名册尺）、`probe-path.sh`/`probe-path2.sh`（PATH 陷阱探针）、
`commit-msg.txt`；staging＝`speech398v3base.test.exe`/`speech398v3.test.exe`＋三 DLL
（`logs/dll-sha256.txt`）；logs＝{suite-staged-base, suite-bare-base, sweep-run1, sweep-run2,
suite-staged, suite-bare, gofumpt.out, vet-speech, d22scan, d22scan-selftest,
dll-sha256, build-base, build-impl}。Git Bash runner、`unset GOTMPDIR` 每跑、rc 重定向现取、
音频错峰跑前现核。
