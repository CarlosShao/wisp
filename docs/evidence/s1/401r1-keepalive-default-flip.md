# 票 401 落地腿 `401-r1` —— keepalive 默认翻转（写腿·实施件）

日期 2026-11-01 · 腿 `401-r1`（写、实现：翻转 keepalive 生产默认 false→true＋测试断言联动）·
起手 HEAD `372c62a01`（dev，开工三核过：①HEAD 在 dev＝派单门槛本身（A1376 拍板提交）
②索引跟踪面零改动（porcelain 全为未跟踪 probes 杂物）③写面两枚零他人脏）·
派单存档 `.scratch/wisp/probes/orch/20261101/txt/401-r1-dispatch.txt` 在案 ·
现跑 2026-11-01 20:5x–21:5x +08。

> **环境**：Git Bash runner（`D:\work\soft\Git\bin\bash.exe`）；`unset GOTMPDIR` 每跑；
> rc 一律重定向进 log；gofumpt 绝对路径 `/d/work/base/gopath/bin/gofumpt.exe`。
> **staged 形状**（397-v2/392-v3 同款、A1326 卡）：`go test -c -tags webrtcap
> -o D:\tmp\wisp\20261101\401r1\bin\*.test.exe ./cmd/wisp/`（CGO 两旗标 `D:/` 盘符形
> ＝393-v1 处方）＋三 DLL 贴 binary（sherpa-onnx-c-api/-cxx-api/onnxruntime）＋
> `WISP_SILERO_MODEL`/`WISP_SPIKE_MODELS`/`WISP_TTS_DIR`/`WISP_KWS_DIR`/`WISP_15E_TMP`/
> `WISP_43_REPO_ROOT` 六 env＋cwd=包目录＋`-test.v -test.timeout 30m`＋每跑前错峰环
> （tasklist 无 test.exe 方开跑）。**裸跑/`./...` 形**：POSIX `/d/` 形 PATH 注入
> `third_party/sherpa-onnx`（398-v3 绞碎陷阱处方；本腿第一发 `./...` 扫没注入即撞
> `0xc0000135` 三包假红，具名见 §6④）。**临时件（只建不删）**：
> `D:\tmp\wisp\20261101\401r1\{bin,logs,out,mut,cases}\**`。**⛔ push；⛔ 票面/台账/
> specs/PLAN/thresholds/golden/testdata 仓根/frontend/design 零字节；⛔ internal/**
> 零字节；⛔ 三处窗间清逻辑零字节。**

## 0. 一句话

**翻转落地、AC#1–#5 全过、停手三条全未触发。** ★:316 `keepalive: false→true`
一字翻转（:87/:253 两处注释随翻联动、行数中性 10/10；三处窗间清逐块零触碰——
突变件 diff 即为证）；★断言联动＝`WARM REVIVAL MOVED` 守卫上默认路径（持引擎窗
r1/r2 无开关直判）＋默认钉翻转＋冷姿势降为开关面；★AC#2 staged 全链复读 **K=3 发
全绿**：持引擎窗 episodes=1/triggers=1/engine_warm=true/engine_built=false/**ready=9728
逐位同 397-v2**/rms -23.59..-27.86@-30 同带、持引擎 (b) 窗 triggers=0 无假 start 漏网
（398-v3 曲线在默认下成立）、392-r1 finding 形不回归（第二暖窗重新会判 ×3）——
**397-v2 残余①（M3 摆动族多发稳定性）就此闭合**；★AC#3 反向三发：摘 det
Finalize＋冷段 ⇒ **判死形默认下定向复现**（rc=1、r2 episodes=0/triggers=0、红句逐字
——红名断言红侧也证到在默认下活）；摘 APM Reset／摘 resampler Flush ⇒ 行为面无
稳定失败形（1 摆动红＋2 绿／rc=0，397-v2 跨面记录同向、承重读数仍在 aecspike
rig 残余深度面）；还原核 pristine rc=0；★AC#4 零回归：26p keyed 4/4、tuning/
dispose rc=0、六族 34/1/0（唯一红＝在册负载红 `Test384ResidentWatchdogEndToEnd`，
单跑 rc=0）、**`LATENCY openDelay=262ms chunk=640ms ballLatency=902ms` 逐位**、
四数 **rc=0 RUN=98/PASS=75/FAIL=0/SKIP=11 逐位**（SKIP 集逐名同形）、398 面
**rc=0 71/65/0/6 逐位同 398-v3**（6 SKIP＝punc 缺席族逐名）；`./...` 扫（PATH 注入
修正后）红集＝cmd/wisp 一枚（**HEAD 既有整包死形、双树逐位归因闭合、零翻转归
属**，§4①）＋risk C26 预算计时摆动（在册族）；d22scan 正规门 **clean rc=0**
（ban#8 internal/=735、cmd/=184、bans#1-5 internal/=288、cmd/=52；本腿零新文件）；
gofumpt 0 行、vet tagged/notag rc=0。★纪律格全过（§5）。

## 1. ★AC#1 默认翻转 ✓（一字＋断言联动）

1. **`cmd/wisp/resident_voice_bargein_windows.go:316`**：`keepalive: false, // the
   measured-good default...` → `keepalive: true, // the 401 flip: the cleared warm
   posture IS the default`——功能改动恰一字（false→true），注释同句随翻。
2. **联动注释两处（行数中性，diff 10 insertions/10 deletions）**：文件头 :87-89
   （`Production runs the posture OFF...` → `ON since the 401 flip...`）＋ struct
   keepalive 字段注释 :253-259（`runs it OFF...` → `runs it ON since the 401
   flip...`）。`git diff --stat`＝恰一枚产品文件 10/10；397-v2 §1 的
   「closeEngines/dispose/readings/drain 零 hunk」在翻后照样成立（diff 仅三 hunk：
   头注释/字段注释/构造行）。
3. **三处窗间清逻辑零触碰的块级证**：AC#3 三突变件以翻转后文件为基、sed 行段
   恰 `557,564d`（APM Reset 块）／`565,567d`（resampler Flush）／`568,573d`
   （det Finalize＋1.2s 冷段）——行号与 397-v2 完全一致（行数中性改法的目的即此）；
   每件 diff 逐行核＝恰目标块；还原核后 `m.apm.Reset()` 现读在 :558。
4. **断言联动（392-v3 §6 第 3 步「红名断言必须连默认一起翻」的执行）**：
   - 默认钉翻转：`readings().Keepalive` 必真（旧钉「the measured finding keeps it
     off」反转为「the 401 flip made the cleared posture the default」）；
   - **`WARM REVIVAL MOVED` 守卫上默认路径**：重构后 Part A 两扇持引擎 talker 窗
     （r1/r2）在**无任何 setKeepalive 开关**下直判 episodes=1/triggers=1——红句
     逐字保留（`want episodes=1 triggers=1@talker (the 392-r2 revival-cool
     readings, live under the default)`），397-r2 改钉的运行路径核＝连默认翻达成
     （绿侧 ×3 发＋红侧 M3 一发双证）；
   - 持引擎 (b) 窗（r3）钉 `triggers=0`＋无 EnergyPass episode——398-v3 衰减曲线
     的「窗间清生效 ⇒ 无假 start 漏网」在默认路径上的直接探针；
   - 冷姿势降为开关面（Part B，`setKeepalive(false)`）：冷 (a)（built=true/warm=
     false/ready≥9000/triggers=1）＋冷 (b)（triggers=0/built=true）＝pre-392 形保持；
   - restore 核：开关回 true，冷轮关掉的引擎由默认首窗重建（built=true）＋(b) 形
     保持——七窗全序：control／持 (a) r1／持 (a) r2／持 (b) r3／冷 (a) r4／冷 (b)
     r5／restore (b) r6（397-v2 的七窗形状集全保留，归属重排为「默认担多窗复活、
     冷面担反向」）；测试名保持 `Test392KeepalivePostureFindingsAndColdDefault`
     （摆动族跨腿可比性优先，名中 ColdDefault 半句由件内注释声明为历史钉）。
5. **`resident_voice_windows.go:751` 现核**＝只描述 dispose 级联、无默认值声称，
   写面外零触碰；`armJudge` :755-759 注释（「pre-392 shape every announcement
   runs」）翻后仍真（measure-then-arm 零动，ready=9728 全窗同值为证）。

## 2. ★AC#2 staged 全链复读 ✓（K=3 发全绿＝权威读数）

**七窗全序 ×3 发、全 rc=0**（`logs/ac2-run{1,2,3}.log`＋`out/392-keepalive-run{k}.json`
快照；工件键＝control_first_build/warm_rounds/warm_echo/cold_round/cold_echo/
restore_echo）：

| 窗 | 形状 | run1 | run2 | run3 |
|---|---|---|---|---|
| control | 默认首建：keepalive=true/built=true/warm=false/**EnginesHot=true** | ✓ | ✓ | ✓ |
| 持 (a) r1 | ep=1/trig=1/warm=true/built=false/ready=9728 | rms -23.59 lag 480 | -26.00/621 | -25.83/458 |
| 持 (a) r2 | **同上（397-v2 warm2 权威形）** | -23.60/536 | -26.01/560 | -27.86/564 |
| 持 (b) r3 | trig=0/below-cut/无假 start 漏网 | ep=2 rms -49.56 | ep=1/-43.20 | ep=1/-50.89 |
| 冷 (a) r4 | trig=1/built=true/warm=false/ready=9728 | -25.06 lag 0 | -25.66/490 | -25.27/527 |
| 冷 (b) r5 | trig=0/built=true/below-cut | ep=1/-46.36 | ep=1/-57.59 | ep=1/-38.94 |
| restore r6 | trig=0/built=true/warm=false | ep=1/-46.21 | ep=0/**-999 哨兵** | ep=1/-39.31 |

- **对照 397-v2 七窗全序（397-v2 §2 权威读数）**：warm2 权威形（episodes=1/
  triggers=1/engine_warm=true/engine_built=false/**ready=9728 逐位**/rms@-30 带）
  在本腿＝**每一扇持引擎窗都成立 ×3 发**（比 397-v2 的单窗更强：默认把复活变成
  多窗常态）；ready=9728 与 397-v2 冷/暖全窗逐位同值；rms -23.59..-27.86@-30 对
  397-v2 warm2 -25.71 同带（392-v3 §7 在册口径：真 TTS 材料判类别/相对，单点
  微差非移动）；持 (a) 剖面（res 桶 0 消量深度 -27.7..-34.9、桶 1 talker 骑近端
  -23.3..-26）与 397-v2 warm2（-33.9/-25.1/-21.6）同类。
- **对照 398-v3 衰减曲线**：三处窗间清（Finalize 关口开话态＋1.2s 冷段吸收
  hangover）在默认路径上零假 start 漏网——持 (b) 窗 triggers=0 ×3、restore/冷 (b)
  triggers=0 ×3（398-v3 点位：假 start 恒落边界后第 1 空窗且事件弃置、RNN 半
  400ms 死、话语半 Finalize 即收——三清组合正控）。
- **392-r1 finding 形不回归**：第二暖窗重新会判 ×3（r2 全绿）——判死形
  （episodes=0/triggers=0）在 pristine 默认链零出现。
- **397-v2 残余①闭合**：M3 摆动族「staged 链判死形非逐发确定」的多发稳定性读数
  ＝本腿 K=3 发全绿（同 rig、同 staging、错峰窗）。
- **口径注记**：r2 lag 值 536/560/564 vs 397-v2 warm2 485、AC2 run1 冷 (a) lag=0
  一次（对齐器零读数、触发照常）——真 TTS 材料跨跑微差带，类别判据（triggers=1）
  不动。
- **397-r4 哨兵首次在默认链接活**：run2 restore episodes=0（echo 完全门净＝比
  常形更干净的健康形）→ rms 走 `-999` 哨兵、工件照常落盘、rc=0——397-v2 残余②
  的修复在翻转链上实战验证。

## 3. ★AC#3 反向三发 ✓（判死形默认下复现；两发无失败形同 397-v2 跨面记录）

突变件＝release() keepalive 分支逐块摘除（**三件 diff 逐行核：每件恰目标块**；
三 overlay 构建二进制 SHA 互异且均 ≠ pristine；overlay JSON 全正斜杠键；仓外
`-overlay` 构建期虚拟替换、工作树文件零触碰——还原核 diff 仍 10/10）。同一
pristine 测试（重构后件）、同一 staging：

| 发 | 摘除 | 本腿读数 | 判 |
|---|---|---|---|
| M3 | det Finalize＋冷段 | **首发 rc=1＝r2 `episodes=0 triggers=0`（`WARM REVIVAL MOVED` 红句逐字）** | ✓ **判死形定向复现**（默认路径、392-r1 finding 形） |
| M1 | APM Reset（裸携） | 首发 rc=1 落**冷面 (a) 触发**（`triggers=0 episodes=1`、r4 剖面 7 桶＝未打断播满）；**复跑两发 rc=0 全 PASS**（冷面 trig=1、r4 桶 0 -28.8/-33.6 同 pristine 带） | ✗ 失败形未复现（首发＝(a) 触发边际摆动族，397-v2 残余①同族；与实现件「行为绿」同向） |
| M2 | resampler Flush | **rc=0 全 PASS** | ✗ 无失败形（与实现件同向） |

**还原核**：三发后 pristine 复读 rc=0（`ac3-restore-pristine.log`）；工作树零染
（`git diff --stat`＝10/10、Reset 行 :558 现读在）。**跨面差异维持具名（397-v2
三发在案）**：no-reset 的承重失败形（残余骑近端 ≥14.8dB）在 staged 链行为断言面
四发（397-v2 一发行为绿＋本腿首发摆动红（冷面触发）＋两发绿）仍未迁移——Reset
承重读数仍以 aecspike rig 19.24dB 残余深度为最敏感证据（不在本测试断言面）。
**失败发不写工件 ⇒ 快照残留上一过跑件**（`392-keepalive-m1-no-reset.json`＝
run3 逐位同、`...-m3-...json` 同踩——证据一律以 log `--- PASS/FAIL` 行与
ROUND-*-PROFILE 为准；m1-rerun1/m2/restore-pristine 快照均鲜）。

## 4. ★AC#4 零回归 ✓（翻转骑默认进每一族）

| # | 跑 | 二进制 | rc | 读数 |
|---|---|---|---|---|
| 1 | 26p keyed 四枚（(a) 400ms/(b)＋M2/等价网格） | staged | **0** | **PASS=4 FAIL=0**（翻转后 26p 窗骑持引擎跑、全绿） |
| 2 | 392 调参族 `Test392AecTuningDoubleTalkSwing` | staged | **0** | PASS=1（调参窗骑持引擎跑＝翻转后生产形状；apmTuning 守卫保一参一发） |
| 3 | dispose 面 `Test392DisposeReleasesTheEngines` | staged | **0** | PASS=1（C11 尾把持引擎摘下） |
| 4 | 六族合跑 `Test43\|Test15e\|Test41a2r1\|Test381r2\|Test390r1\|Test26R3\|Test384` | no-tag | 1 | **34/1/0——唯一红 `Test384ResidentWatchdogEndToEnd`＝合跑负载时窗在册形状**；单跑 **rc=0**（397-v2 §4 同形） |
| 5 | LATENCY 尺 | — | — | **`openDelay=262ms chunk=640ms ballLatency=902ms` 逐位**（sixfamilies.log :428） |
| 6 | 四数（`./internal/audio/... -count=1 -v`、无 tag、POSIX PATH 注入、静窗） | — | **0** | **RUN=98 / PASS=75 / FAIL=0 / SKIP=11 顶层口径与 392-r1:108/392-r2/392-v3/397-v2 逐位同**；SKIP 集＝TestLiveWasapiSmoke＋aecspike 硬件/工件门逐名同形（含密掩占位门名原样） |
| 7 | 398 面（`./internal/speech` staged、三 DLL＋两 env） | speech | **0** | **RUN=71 / PASS=65 / FAIL=0 / SKIP=6 逐位同 398-v3 落件后**；6 SKIP＝punc 缺席族逐名（300Char/PerExample/PositiveControls/LatencyBudget/CrossesSeam/RealAsrFinal） |
| 8 | `go test ./...` 全仓扫（无 WISP_* env、POSIX PATH 注入、-count=1） | — | 1 | 32 包 ok；**cmd/wisp 一枚红（§4① 具名归因）＋internal/risk 一枚红（§4②）**；speech/aecspike 在 PATH 注入后绿（首发扫的 0xc0000135 三包假红＝本腿环境缺口，§6④） |
| 9 | d22scan 正规门（`sh scripts/d22scan.sh`：正控＋实扫） | — | **0** | **clean**；ban#8 internal/=735、cmd/=184、bans#1-5 internal/=288、cmd/=52、frontend/=116、design/=39（本腿零新文件；与 398-v3 在册 735/183 的 cmd/+1、bans#1-5 internal/ 287→288 归 HEAD 中途票 16/A1376 两提交，非本腿） |
| 10 | gofumpt `-l cmd/wisp/`；vet tagged＋notag | — | **0** | 0 行；双形状 rc=0 |

**§4① cmd/wisp `./...` 红的归因（双树对照闭合＝HEAD 既有死形、零翻转归属）**：
无 `-v` 扫（dotdotdot2）与定向 verbose 跑（cmdwisp-verbose）两发在**翻转树**同位
死：`=== RUN TestAC14GoSideEvalPushReachesThePage` 之后零输出（无 `--- FAIL`、无
panic、无 exit status 行——测试二进制在 panel 窗创建中硬暴毙，`33r9-nail` 面板
harness 窗翻动域），顶层完成集 190 枚。**HEAD pre-flip 基线**（`git show
372c62a01:` 两枚写面文件经仓外 `-overlay` 构建期替换＝同工作树同环境、零工作树
触碰）**逐位复现同死形**：同枚测试、同位置、顶层完成集 190 枚与翻转树 **零差**
（Compare-Object 空）。⇒ 整包 run 的 panel 域死形为 HEAD 既有形状、与翻转零
归属；在册批量跑法（六族等定向 `-test.run` regex）从不运行整包、从不踩它；
单发修复属 panel/WebView2 域、写面外，本腿⛔越（如实记账，归编排者派单）。
**§4② internal/risk 红**：`TestResolvePerCallBudget`（C26 预算 1.011ms/op vs
1.000ms、1316 样本）＝计时预算摆动族（397-v2 `./...` 红集含 risk 族「env/负载
形状差」在案）；本腿改动面（cmd/wisp 一枚 flag 默认）与 risk/pathresolver 零
编译域交集（package main 不可 import）。

## 5. ★AC#5 纪律格 ✓

写面恰 3 枚（commit pathspec 逐枚）：`cmd/wisp/resident_voice_bargein_windows.go`
＋`cmd/wisp/resident_voice_bargein_392_windows_test.go`＋
`docs/evidence/s1/401r1-keepalive-default-flip.md`；`resident_voice_windows.go`/
`internal/**`/AEC3 源树/thresholds/golden/deps.toml/票面/台账/specs/PLAN/
testdata 仓根/frontend/design 零触碰；**三处窗间清逻辑零字节**（§1.3 块级证）；
临时件只在 `D:\tmp\wisp\20261101\401r1\**`（只建不删）；**⛔ push**；子代理角色
（写腿：无翻框/无台账/无票面改动）。

## 6. 残余与边界（如实记账）

1. **M1 首发冷面 (a) 触发摆动（新读数）**：摘 APM Reset 的首发在冷面 (a)
   triggers=0（r4 未打断播满 7 桶），复跑两发全绿——(a) 触发边际摆动族
   （26p (c)/392 调参/397-v2 M3 首发、在册同族）在突变体上的又一点名；机制上
   M1 冷面引擎集与 pristine 全同（裸携 APM 在开关面被 closeEngines 关掉、冷臂
   全新重建），归属摆动而非承重失败形。
2. **冷 (a) lag=0 读数两次**（AC2 run1 冷面、M1 复跑2 r2）：对齐器包络零点读数、
   触发照常——真材料时序微差带，类别判据不动。
3. **持 (b) 窗 episodes 值 1..2 摆**（全 below-cut、triggers 恒 0）：episode 计数
   非断言面、profile 在工件。
4. **`./...` 首发扫的 0xc0000135 三包假红＝本腿环境缺口（398-v3 陷阱的又一形态）**：
   无 PATH 注入的 `go test ./...` 让 cmd/wisp/aecspike/speech 三包全撞 DLL 加载
   失败；修正（POSIX `/d/` 形注入）后 speech/aecspike 绿。397-v2 的同位扫未撞＝
   其脚本前一步的 PATH export 恰好续到了 ./... 步（隐性依赖）；**遗产：任何
   `./...` 扫必须显式注入，不得依赖前步 export 残留**。
5. **398-v3 分母微差**：ban#8 cmd/ 183→184、bans#1-5 internal/ 287→288——本腿
   零新文件，+1/+1 归 HEAD 中途提交（票 16 格④/A1376）；逐源可数。
6. **停手条款三件全未触发**：①pristine 复读 K=3 全绿（无判死形、复活常在）；
   ②M3 摆动族非全红（3/3 绿）；③三处窗间清零触碰（sed 块级证）。
7. **整包 cmd/wisp run 的 panel 域死形（HEAD 既有、双树归因闭合，新入册）**：
   `go test ./cmd/wisp/`（全包、无定向 regex）在 TestAC14GoSideEvalPushReachesThePage
   中途硬暴毙（无 panic／--- FAIL／exit status 行、panel 窗创建域）；HEAD pre-flip
   基线（overlay 双文件替换）逐位复现、顶层完成集零差（§4①）。修复属 panel/
   WebView2 域（写面外）——建议编排者单独开票；在册批量跑法不受影响。

## 7. 五格判＋可否收口

- **AC#1 默认翻转 ✓**（§1：一字＋联动注释＋断言联动双证）。
- **AC#2 staged 全链复读 ✓**（§2：K=3 全绿、ready=9728 逐位、无假 start 漏网、
  finding 不回归、397-v2 残余①闭合）。
- **AC#3 反向判据 ✓**（§3：M3 判死形默认下定向复现、M1/M2 无失败形如实记＋
  还原核）。
- **AC#4 零回归 ✓**（§4：定向 4+1+1 全绿、六族 34/1/0 在册形状、四数/398 面逐位、
  262/640/902 逐位、d22scan clean；`./...` 红集具名归因）。
- **AC#5 纪律 ✓**（§5）。
- **可收口：是**（票 401 候选；AC 框勾选与票面 Status 归编排者收口动作，本腿
  ⛔改）。

## 附：本腿自跑尺（全部可重放）

`D:\tmp\wisp\20261101\401r1\`：`build401r1.sh`（gofumpt/vet×2/三二进制构建）、
`ac2-run.sh`（K=3 权威复读）、`muts401r1.sh`＋`ac3-run.sh`＋`m1-reruns.sh`
（三突变＋四发＋还原核）、`ac4-batch.sh`（定向＋六族＋四数＋398 面）、
`fixups.sh`（./... 修正扫＋d22scan 正规门）、`isolate-cmdwisp.sh`（红名隔离）、
`head-baseline.sh`（HEAD pre-flip 基线、overlay 双文件替换）、`commit-msg.txt`；bin＝pristine/M1/M2/M3/no-tag/speech 六二进制＋三 DLL；
logs＝全部 rc 与读数；out＝七发工件快照（两枚陈旧已具名）。Git Bash runner、
`unset GOTMPDIR` 每跑、rc 重定向现取、音频错峰跑前现核。
