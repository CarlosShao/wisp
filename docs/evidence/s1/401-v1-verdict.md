# 票 401 终裁腿 `401-v1` —— keepalive 默认翻转（非实现者裁决表）

日期 2026-11-01 · 裁决对象 `c96c0a59c`（3 枚：产品码＋测试件＋证据件；真父 `372c62a01` 现核一致；401 实现归属全库唯一此提交）·
本腿全部实验在导出树 `D:\tmp\wisp\20261101\401v1\tree\`（`git archive c96c0a59c`；两枚码面 SHA256 **commit=导出树逐位同** `9fc7c5ae…`/`9e4afdc6…`）·
补 gitignored 面（397-v2 §6⑤ 处方）：`build/models/` manifest 对 12KB＋`third_party/model-fetch/` 845MB＋`third_party/spike-models/` 420MB；staging＝A1326 全套（staged 二进制＋三 DLL 贴身＋WISP_* env＋cwd=包目录＋`unset GOTMPDIR`＋Git Bash runner＋每跑前 rig 错峰环）·
工作树零触碰、⛔ push/amend/reset/rebase/stash/checkout、⛔ `.go` 写面。

## 0. 一句话

**五格全过、可收口（附定性残余）。** ★AC#1 读码核：`:316` 一字翻＋两注释联动行数中性（diff 恰三 hunk；`4c8b8ea60→372c62a01` 产品码零漂移 ⇒ 三处窗间清 `:557-564/:565-567/:568-573` 与 397-v2 行段逐位同＝零触碰）；★WARM REVIVAL MOVED 守卫上默认路径核＝连默认翻达成（持引擎窗无 setKeepalive 直判，红句逐字）；★AC#2 本腿自跑 K=3 全绿、`episodes=1/triggers=1/engine_warm=true/engine_built=false/ready=9728` ×3 发逐位同 397-v2＝**397-v2 残余①闭合**；★AC#3 反向三发：M3 判死形默认路径**首发定向复现**（rc=1、round2 `episodes=0 triggers=0`、红句逐字）、M1 三发行为绿＋机制核（arm `:389-391` 冷臂 closeEngines 全清重建＝冷面引擎集与 pristine 全同）、M2 rc=0、还原核 rc=0 双树零染；★AC#4 零回归：26p 4/4＋六族 34/1/0（唯一红＝在册 Test384、单跑 rc=0）＋`262/640/902` 逐位＋四数 **98/75/0/11 逐位**＋398 面 **71/65/0/6 逐位**＋d22scan 正规门 clean（735/184/288/52 逐位）；★三条新发现核全落：./... 死形**本腿双树逐位归因闭合**（572=572、唯一 PASS/FAIL 差一枚在册摆动族且方向利好翻转）、MSYS2 陷阱形态**同树同包双证**（无注入 `0xc0000135` vs 注入 rc=0）、PerCallBudget 在册＋本腿一红一绿摆动实锤。

## 1. AC#1 默认翻转核 ✓（读码核，非采信）

1. **`:316` 一字**：`keepalive: false, // the measured-good default: the exact pre-392 shape` → `keepalive: true, // the 401 flip: the cleared warm posture IS the default`——功能恰一字、注释同句随翻。票面 AC#1 写 `:294`＝397 落地前行号（397-v2 §1 同款口径），现位 `:316`。
2. **两注释联动、行数中性**：文件头 `:84-89`（`Production runs the posture ON since the 401 flip…`）＋struct 字段注释 `:251-259`（`runs it ON since the 401 flip…`）。产品 diff 恰三 hunk、10/10 行中性；`git diff 4c8b8ea60 372c62a01` 产品码**零漂移** ⇒ 三处窗间清块现读 `:557-564`（APM Reset）/`:565-567`（Flush）/`:568-573`（det Finalize＋1.2s 冷段）与 397-v2 行段完全一致；本腿三突变件 sed 行段 diff＝恰目标块逐行（`557,564d`/`565,567d`/`568,573d`）＝**零触碰块级证**。
3. **断言联动＝连默认翻**（392-v3 §6③「红名断言连默认一起翻」执行核）：默认钉翻转（`:405-407` readings().Keepalive 必真、红句「the 401 flip made the cleared posture the default」）；★**WARM REVIVAL MOVED 守卫上默认路径**：Part A 两扇持引擎窗 r1/r2 **无任何 setKeepalive 开关直判** `episodes=1/triggers=1`（`:438-453`，红句逐字含 `live under the default`）——绿侧 ×3 发（§2）＋红侧 M3 一发（§3）双证；持引擎 (b) r3 钉 `triggers=0`＋无 EnergyPass（`:460-472`）＝398-v3 曲线默认路径探针；冷姿势降为开关面（`:477` setKeepalive(false)，r4/r5 pre-392 形）；restore `:508-522`。
4. **七窗重构核**：control／持 (a) r1／持 (a) r2／持 (b) r3／冷 (a) r4／冷 (b) r5／restore r6；工件键 `control_first_build/warm_rounds/warm_echo/cold_round/cold_echo/restore_echo`；测试名保持 `Test392KeepalivePostureFindingsAndColdDefault`。397-r4 JSON 哨兵（-999）与 `Test392JSONSentinelClampsNonFiniteReadings` 在件（两树同有，ac5 双树 RUN 集含它）。
5. **写面外**：`resident_voice_windows.go:751-757` 只描述 dispose 级联、无默认值声称；`armJudge` `:755-759` 注释翻后仍真（measure-then-arm 零动、ready=9728 全窗同值为证）。

## 2. AC#2 staged 全链复读 ✓（本腿自跑 K=3＝权威读数）

- **3 发全 rc=0**（`ac2-run{1,2,3}.log`＋`out/392-keepalive-run{k}.json`）：持引擎 r1/r2 `episodes=1/triggers=1/engine_warm=true/engine_built=false/ready=9728` **×3 发逐位同 397-v2 warm2 权威形**（比 397-v2 单窗更强：每扇持引擎窗成立）；rms **-19.55..-25.43@-30**（397-v2 warm2 -25.71、401r1 -23.59..-27.86——真 TTS 材料跨跑微差带，392-v3 §7 类别判据口径在册）；持 (b) r3 `triggers=0` ×3（below-cut -42.1..-51.5）＝**无假 start 漏网、398-v3 曲线默认下成立**；冷 (a) r4 `triggers=1/built=true/warm=false/ready=9728`；冷 (b) r5 `triggers=0/built=true`；restore r6 `triggers=0/built=true` ×3。
- **392-r1 finding 不回归**：判死形（episodes=0/triggers=0）pristine 链 ×3 发零出现——**397-v2 残余①（M3 摆动族多发稳定性）就此闭合**（本腿 3＋实现件 3＝连续 6 发全绿）。

## 3. AC#3 反向三发 ✓（本腿自跑；判死形首发定向复现）

| 发 | 摘除 | 本腿读数 | 判 |
|---|---|---|---|
| M3 | det Finalize＋冷段（`568,573d`） | **首发 rc=1＝round2 `episodes=0 triggers=0`（WARM REVIVAL MOVED 红句逐字）** | ✓ **判死形默认路径定向复现**（392-r1 finding 形） |
| M1 | APM Reset（`557,564d`） | **3/3 发 rc=0 行为绿**（复活 1/1、冷面形全保持、restore 净）；复跑1 冷 (a) lag=0 一次（对齐器零读数、触发照常＝在册微差） | ✗ 失败形未复现（与 397-v2/401r1 同向） |
| M2 | resampler Flush（`565,567d`） | **rc=0 全 PASS** | ✗ 无失败形（同向） |

- **M1 机制核**：`arm()` `:389-391` 冷臂 `(m.apm!=nil||m.det!=nil) && !m.keepalive → closeEngines()` 全清重建 ⇒ M1 突变下**冷面引擎集与 pristine 全同**（裸携 APM 只暴露在持引擎路径；其失败形＝392-r2 矩阵 no-reset W2，在 staged 链 397-v2 ×1＋401r1 ×3＋本腿 ×3 行为绿未迁移）。
- **还原核**：三发后 pristine 复读 rc=0；工作树 `git status --porcelain -- cmd/wisp/` 全程空、导出树 pristine SHA 不变（`9fc7c5ae…`）；四二进制 SHA 互异（pristine `0289da09…`/M1 `c40a76ba…`/M2 `85beadd6…`/M3 `8f1af7cf…`）。

## 4. AC#4 零回归 ✓（导出树批量）

| # | 跑 | rc | 读数 |
|---|---|---|---|
| 1 | 26p keyed 四枚（staged） | **0** | **PASS=4 FAIL=0** |
| 2 | 392 调参族＋dispose 面（staged） | **0** | PASS=1 ＋ PASS=1 |
| 3 | 六族合跑（no-tag） | 1 | **34/1/0**——唯一红 `Test384ResidentWatchdogEndToEnd`＝合跑负载时窗在册形状，**单跑 rc=0** |
| 4 | LATENCY 尺 | — | **`openDelay=262ms chunk=640ms ballLatency=902ms` 逐位** |
| 5 | 四数（internal/audio 无 tag、POSIX PATH 注入） | **0** | **RUN=98/PASS=75/FAIL=0/SKIP=11 逐位**；SKIP 集＝TestLiveWasapiSmoke＋aecspike 门逐名同形（含 `TestSpikeAec3LongWindow` 占位名原样） |
| 6 | 398 面（speech staged、树指 env） | **0** | **RUN=71/PASS=65/FAIL=0/SKIP=6 逐位同 398-v3/401r1**；6 SKIP＝punc 缺席族逐名同 |
| 7 | d22scan 正规门（`sh scripts/d22scan.sh`） | **0** | **clean**；bans#1-5 internal/=288 cmd/=52、ban#6 frontend/=116、ban#8 internal/=735 cmd/=184 逐位同在册；selftest（模块内）rc=0 **PASS=35/RUN=77 逐位** |
| 8 | gofumpt＋vet tagged/notag | **0** | 0 行；双形状 rc=0 |

- **中间红两枚如实记（均非回归）**：①补面前树（spike-models 缺）398 面 71/64/1/6——`TestDefaultKwsConfigResolution` 走查红＝397-v2 §6⑤ 导出树 gitignored 面缺口同类；②树建二进制＋工作库指 env 时同枚再红——该测试拿 **env 锚 real 对 `runtime.Caller` 锚 walk-up**（`kws_test.go:270-285`），树建二进制的 Caller 恒指树、两锚几何必然分叉；树指 env 后 rc=0 71/65/0/6。internal/speech 不可能 import package main＝零翻转编译域。
- root 模块 `go run ./tools/d22scan` rc=1＝**子模块 go.mod 边界**（`tools/d22scan/go.mod` 树内核实）；在册形＝正规门脚本＋模块内 selftest，两者全绿。design/=14（导出树仅 tracked 文件；工作库 39 含未跟踪）非违规。

## 5. 三条新发现核

1. **`./...` 整包 cmd/wisp 死形＝双树逐位归因闭合、零翻转归属（本腿独立闭合）**：本腿导出树 rig 上**静默硬死形不复现**——翻转树与 pre-flip 基线（两枚写面文件 `-overlay` 换 `372c62a01` 版）verbose 全包跑**双双完整跑完**（顶层 RUN 集 **572=572 逐位同**、包汇总正常落 `FAIL …cmd/wisp 395.7s/388.7s`）；AC14 在两侧同为**正常断言红**（33r5 push 维度、panel 窗域）。逐位比对唯一 PASS/FAIL 差＝`TestAC4FocusReturnToPriorWindowGap33r5` **pre-flip 红/flipped 绿**（方向利好翻转）；同一翻转树 dotdotdot3/4 红、ac5-verbose 绿＝**摆动族在同一树态内实锤**。实现件工作库 rig 的静默死形（190 枚完成集）作为其 rig 证据在案；两 rig 两形同归一词：AC14/panel 窗域损伤 HEAD 既有、与翻转零归属（overlay 换码后 AC14 两侧同红）。panel/WebView2 域修复票建议维持（两形并记）。
2. **MSYS2 `./...` 陷阱形态＝同树同包双证**：`internal/speech` 无 PATH 注入裸跑 rc=1 **`exit status 0xc0000135`**；POSIX `/d/` 形注入后 rc=0 ok。398-v3/401r1 遗产成立：**任何 `./...` 扫必须显式注入、不得依赖前步 export 残留**。
3. **`TestResolvePerCallBudget` 在册**：本腿 dotdotdot2 红（3.70s）→dotdotdot3/4 绿＝计时摆动族行为实锤；在册出处＝台账 A1358「上报不修」名单（**AC247+PerCallBudget flake**）；与翻转零编译域交集（package main 不可 import）。

## 6. 纪律格 ✓

本腿仓内写面恰 1 枚＝本裁决表（提交 pathspec 逐枚）；零 `.go` 写、工作树零触碰（ac3/ac5 现核 status 空＋SHA 不变）；实验全在导出树＋`-overlay` 构建期替换；临时件只在 `D:\tmp\wisp\20261101\401v1\**`（只建不删）；⛔ push/amend/reset/rebase/stash/checkout。

## 7. 残余与边界（如实记账）

1. `TestAC4FocusReturnToPriorWindowGap33r5` 窗摆动（§5①）：pre-flip 红/flipped 绿一差、同树跨跑两向——在册窗测摆动族，非翻转效应。
2. 导出树 rig 形差具名：无 `.git`（`TestCleanCheckoutBuilds_AC11`/`TestReadGitOnThisRepositoryIsSelfConsistent` 结构性红，后者属 `internal/panel/git_test.go`）；裸扫 no-tag 二进制无 AEC3 绑定 ⇒ 26p/392 族在 `./...` 裸扫红（在册批量跑法全部 staged 形、不踩）；design/ 计数 14 vs 工作库 39＝未跟踪文件不入档。
3. 静默死形（实现件 rig）本腿不复现——panel/WebView2 票若开，两形（工作库静默杀/导出树 33r5 断言红）并记。
4. M3 突变件本腿首发红（1/1；397-v2 时代 1 绿 2 红）——突变侧跨 rig 摆动仍可能；闭合证据以 pristine 链 6 发全绿为准。
5. 本腿三轮 `./...` 红集随补面状态收敛（35 speech 红→清）；终态树红集＝cmd/wisp＋panel＋tools（risk 摆动），逐位清单在 `dotdotdot4.log`。

## 8. 五格判＋可否收口

- **AC#1 默认翻转 ✓**（§1：一字＋联动＋三清零触碰＋守卫上默认路径＝连默认翻）。
- **AC#2 staged 全链复读 ✓**（§2：K=3 全绿、ready=9728 逐位、398-v3 曲线默认下成立、**397-v2 残余①闭合**）。
- **AC#3 反向判据 ✓**（§3：M3 判死形首发定向复现、M1 机制核＋三发绿、M2 绿、还原核）。
- **AC#4 零回归 ✓**（§4：定向 4+1+1＋六族 34/1/0＋262/640/902＋四数 98/75/0/11＋398 面 71/65/0/6＋d22scan clean 逐位）。
- **AC#5 纪律 ✓**（§6）；三条新发现核全落（§5）。
- **可收口：是**（票 401 候选；AC 框勾选与票面 Status 归编排者收口动作，本腿⛔改）。

## 附：本腿自跑尺（全部可重放）

`D:\tmp\wisp\20261101\401v1\`：`setup401v1.sh`（导出树＋SHA 核＋补面）、`build401v1.sh`（gofumpt/vet×2/三二进制）、`ac2-401v1.sh`（K=3）、`muts401v1.sh`＋`ac3-401v1.sh`＋`m1-reruns401v1.sh`（三突变四发＋还原核）、`ac4-401v1.sh`（26p/调参/dispose/六族/四数/398 面/扫/d22scan）、`speech-rerun3.sh`（树指 env 398 面）、`ac5-401v1.sh`（双树 verbose＋逐位比对）、`final-sweep401v1.sh`（终扫＋陷阱探针）；bin＝pristine/M1/M2/M3/notag/speech 六二进制＋三 DLL；logs＝全部 rc 与读数；out＝AC2 三发工件快照。Git Bash runner、`unset GOTMPDIR` 每跑、rc 重定向现取、音频错峰跑前现核。
