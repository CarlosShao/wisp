# 票 396 终裁腿 `392-v3` —— 保温修复链收口（非实现者裁决表）

日期 2026-11-01 · 裁决对象 `49ce619f1`（2 枚 +1173：`internal/audio/aecspike/keepalive392r2_windows_test.go` 999 行＋证据件 174 行；真父 `6bc014439` diff 坐实＝申报写面恰好 2 枚）·
越界面本腿现核：`git diff 49ce619f1^ 49ce619f1 -- cmd/wisp/ internal/audio/aec/`＝**0 行**（keepalive=false 默认面双零成立；AEC3 Reset/392 Tuning 语义原样）·
材料现核：票面 5 格＋`392r2-keepalive-fix.md`＋`392r1-session-integration.md`＋`393r2-flush-reset.md`＋`393-v1-verdict.md`＋392-r1 工件逐字节 ·
本腿全部实验在导出树 `D:\tmp\wisp\20261101\392v3\`（`git archive 49ce619f1`＋补 sherpa DLL/spike-models；关键 Go 件 SHA256 逐位对上，`.scratch` 四枚非法字节名 tar 报错与构建面无关）——工作树零触碰、⛔ push。

## 0. 一句话

**五格全过、可收口。** ★★★三项全部本腿门下复跑：毒源矩阵在**两个全新进程逐臂全量**（非抽臂）
复现——48/48 格类别读数（episodes/triggers/trigger 位置/残余深度带）对实现件矩阵零错配；
根因链改判成立且 392-r1 空壳**实锤**（其工件 round 对象逐字段现读＝字面 `{}`）；
复活纪律（持有＋三处窗间清＋1.2s 冷段）W2/W3 全愈复现；采纳路径三步拍定（§6）。

## 1. ★392-r1 空壳核（根因链改判的前半，实锤）

`D:\tmp\wisp\20261031\392r1\out\392-keepalive.json`（1008 字节）逐字段现读：
`warm_round1`/`warm_round2`/`cold_rounds[0,1]`/`cold_echo`/`cold_restore` 全为 `{}`，
唯一有字段的 `cold_control` 是姿势状态结构（Armed/EnginesHot 等），无任何窗读数。
⇒ **392-r1 finding 的「残余剖面骑在近端电平 -22~-25dBFS」半句确无盘上支撑**
（round 字段未导出、剖面从未落盘）；其**类别半**（第二暖窗 episodes=0、triggers=0）成立。
改判（判死＝detector 话态／残余骑近端＝APM 跨窗携带）两半各归各**成立**：
②半在本矩阵有全导出工件（40KB、剖面在案）并跨进程复现（§2）。

## 2. ★毒源矩阵复核（8 臂×3 窗×两跑，门 `WISP_AEC392R2=1`，两全新进程）

本腿自建（CGO 旗标＝aec3r1 build.py COMMON 全集＋`-include cstdint`，`D:/` 盘符形，build rc=0）
两发 rc=0 全 PASS（矩阵 9.90/9.77s＋探针 7.69/7.60s）。**48/48 格（8 臂×3 窗×2 发）类别读数
对实现件矩阵零错配**：episodes/triggers 全同、trigger 恒落 talker 本体（11840/12320 一帧游动，
480 样本＝回声延迟帧距；无一次静音落点）、残余深度：健康臂跨进程差 ≤0.07dB、near ≤1.2dB、
最敏感的 no-reset 臂带内 4.5dB（-42.78/-47.25/-43.47）。逐臂判：

- **毒源＝detector 半**：warm-392 W2＝ep=0（finding 复现）而残余＝cold 全同（-62.11/-62.18/-62.12
  三进程）；iso-res-fresh（仅 det 保温）W2 同死；iso-det-fresh（仅 resampler 保温）W2 trig=1 且
  残余＝cold 全同——**resampler 洗清坐实**。W3 自愈（开口话态寿命）同形。
- **机制单元探针**三进程同：开口话态跨 400ms 静音吞掉第二 talker 的 start、`Finalize` 以跨两段
  累积 Total 摘除；对照（Finalize 后）start 复现且落点 **silence=1 talker=0**（RNN hangover 在
  post-speech 静音上开火——新发现复现）。
- **仪器健康**：每臂 W1 逐位同（三进程）；进程内两跑逐位断言三进程全过（Fatalf 未触发）；
  还原核（全矩阵后全新冷窗）对开场 cold W1 逐位（三进程）。
- 精度注记：证据件 §2「strength 挪到 0.999」对自家工件应为 warm-392 W2 lag=479/str=0.996
  （本腿 0.997 同带）——类别（≈1.000、lag 本体不变、无害）成立，数字以工件为准。

## 3. ★复活纪律复跑＋反向三对＋还原核

- **revival-cool（APM Reset＋resampler Flush＋detector Finalize＋1.2s 冷段，边界事件弃置）**：
  W2/W3 全愈——ep=1、trig=1@talker、残余 -62.09/-62.01 vs cold -62.18/-62.12（实现件
  -62.02/-62.11，同深度带）——三进程复现。**transient 具名复现**：只 Finalize 不冷段的
  revival W2＝ep=1、trig=0（hangover 假 start 吃掉 endpointer、talker 开不了口），冷段吸收后
  消失。
- **反向**：摘 Finalize ⇒ W2 ep=0（finding 回归，三进程）；摘 Reset ⇒ 残余 -47.25/-43.47 vs
  复活臂 -62.09/-62.01（退化 **14.84/18.54dB**；实现件 19.24dB 同类）且 ≈near（-37.88/-36.95）
  ——裸携带 APM 只消 4.6~9.4dB vs 健康 24dB，「残余骑近端」形在链上复现，Reset 承重成立
  （跨进程读类别量 ≥14.8dB，单点数字随收敛态摆）；cold 基准三窗全判、残余 -62 深度＝仪器基准。

## 4. ★零回归（导出树现跑）

- **四数**（无 tag、`PATH=$PWD/third_party/sherpa-onnx`、`-count=1 -v`、错峰核过无 rig 在飞）：
  **rc=0，`-v` 口径 RUN=98 PASS=75 FAIL=0 SKIP=11——与 392-r1:108/392-r2 自记逐位同**；
  SKIP 集与 392-r2 日志逐名同形（TestLiveWasapiSmoke＋aecspike 十枚硬件/工件门，一枚门名
  被密掩占位符替换、枚数与集形与前两腿记录吻合）。新测量件三重门在无 tag 编译域外。
- **d22scan**（Git Bash，导出树；「gitignore 未应用＝全路径扫描」loud 注记在案）：实扫 clean
  rc=0，分母逐枚同（bans#1-5 internal/=287、cmd/=52；ban#8 internal/=731＝393-v1 的 730＋1、
  cmd/=183）；selftest **PASS=35 FAIL=0 SKIP=0（RUN=77）**。ban#8 不扫 docs/，本件不入分母。
- build（webrtcap `-c`）rc=0；vet 三形状（tagged/无 tag/GOOS=linux）全 rc=0；gofumpt
  `-l internal/audio/`＝0 行。

## 5. 五格逐条判＋可否收口

- **AC#1 毒源半隔离 ✓**（§2：detector 半＝毒源、resampler 半洗清，双向隔离＋单元探针直接读出）。
- **AC#2 Reset 接入保温 ✓**（§3：复活纪律 W2/W3 全愈三进程复现；393-v1「Reset 即统一器」推广
  到整池的形状在链上成立）。
- **AC#3 反向判据 ✓**（§3 三读数＋还原核逐位）。
- **AC#4 零回归 ✓**（§4 四数逐位＋SKIP 集同形＋默认面双零 diff＋d22scan/selftest 逐枚同）。
- **AC#5 纪律 ✓**（写面恰 2 枚；越界面 0 行现核；票面/台账/AEC3 源树/thresholds/golden 未动；
  无 push）。
- **可收口：是**（票 396 候选；AC 框勾选与票面 Status 归编排者收口动作，本腿⛔改）。

## 6. ★采纳路径裁定（给编排者；392-r2 读数＝裁定材料）

1. **第一步·monitor 复活＝三处窗间清一体落地**（一腿一个 commit，⛔ 拆开）：`release()`
   keepalive 分支＝APM 持有＋`Reset()`（393-v1 统一器形状：读数与每窗重建同健康且省构造；
   不选维持重建）、resampler 持有＋`Flush()`、detector 持有＋`Finalize(nil)`＋1.2s 冷段、
   边界事件弃置。理由：矩阵证明三件**联合承重**——摘任一件即三失败形之一（no-Finalize＝
   判死、no-Reset＝残余死 ≥14.8dB、no-cooldown＝首暖窗漏打断）；全部公开面、零新 API。
2. **第二步·internal/speech 开票**：silero session 无 RNN reset 面（具名缺口），冷段是 owner
   侧纪律不是 session 清面；hangover 假 start 在收紧档不过阈＝无假打断、但无冷段则首暖窗
   talker 打断会漏。开票评估 session 级清面让 owner 免掉冷段；**非 monitor 采纳的阻塞**。
3. **第三步·keepalive 默认翻转＝后续裁定**：默认保持 false（0 字节已核）。翻转前置＝①落地
   ＋cmd/wisp staged 全链复读（392-r1 的 finding 是 staged rig 读数，翻转前须在真姿势上读
   到全愈）；机械上 392-r1 的 `FINDING MOVED` 红名断言必须连默认一起翻，属那一腿的验收。

## 7. 残余与边界（如实记账）

1. **实现腿 run1-3 是开发史**：其日志含三次失败迭代（仪器确定性断言与探针初版），
   终态 run4＝提交件（行号逐一对上）＋终版工件；本腿两全新进程复现的即终态读数——
   断言咬合（run1 的失败恰证确定性检查非装饰）。
2. **跨进程逐位不成立也不需要**（392-r2 §6.3 预告）：真 TTS 材料，本腿深度微差全量在案
   （§2）；一切判据＝类别/相对判据，「两跑逐位」指进程内（断言三进程全过）。
3. **no-reset 臂是全矩阵最敏感读数**（跨进程 4.5dB 带）——「Reset 承重」应读作 ≥14.8dB 的
   类别量，证据件的单点 19.24dB 是其进程读数。
4. **冷段＝吸收不是清除**：RNN 语境衰减 1.2s 后 W2/W3 全愈，但「衰减多少够」未单独扫参
   （1.2s 一点读数）；更硬的清面归 §6 第 2 步的 internal/speech 票。
5. 392-r1 的 staged-rig finding 本腿未在 cmd/wisp 复读（写面外；§6 第 3 步翻转前置里做）。

## 附：本腿自跑尺（全部可重放）

`D:\tmp\wisp\20261101\392v3\`：tree＝`git archive 49ce619f1` 导出＋DLL/模型补件；
`bin/aecspike392v3.test.exe`（build rc=0）；logs/392v3-run{1,2}.log（门下两发全 PASS）；
out/run{1,2}/392r2-keepalive-matrix.json（全导出工件，与实现件 48/48 格类别对照在案）；
logs/{four-numbers,vet-tagged,vet-notag,vet-linux,d22scan-real,d22scan-selftest}.log；
logs/gofumpt.out（0 行）。CGO_CXXFLAGS 用 `D:/` 盘符形（393-v1 踩坑处方）。
