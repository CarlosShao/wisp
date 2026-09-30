# 62 — 票 62「液态玻璃悬浮球视觉重做」**1:1 对抗验收表**

**验收人＝编排者之外的非实现代理**（腿号 `62-v1`）。**票 62 的实现码我一行没写**：本腿不产 `internal/**`、
`cmd/**` 的任何功能代码；我只跑尺、跑变异（`go test -overlay`，不改工作树）、读台件，然后逐格下判语。
本表的存在本身＝票面 `:95`（AC#8）所判据的对象——AC#8 要的正是"这份表"，所以它**由本表落盘而起**，
但**AC#8 是否因此成立仍要单独裁**（见 §1 第 8 行），本腿不许给自己签字。

**被验对象锚点**：`git rev-parse --short HEAD` ＝ **`d5a59d66`**（本腿起手自取；编排者派单时报的 `d5a59d66`，一致）。
**票面**：`.scratch/wisp/issues/62-liquid-glass-ball-visuals.md` —— **不带 `-done` 后缀＝未结案**；
Status 行（`:3`）＝ `review（改回 review：本票从未被验收）`。
**起手名册**：`git status --porcelain` ＝ **151 行**（存 `.scratch/wisp/probes/62/v1/_roster_start.txt`）。
本腿终态判据**不是"名册为空"**，而是"与起手名册逐字相同、只多出本腿自己的件"。

**前置对齐件（只读腿 `62-a1`，本腿继承其现量、不重跑它的活）**：
`.scratch/wisp/probes/62/a1/ac-evidence.md`（23,232 B，起手锚点 `98df640a`）。
它的合计：7 枚未勾格里 **A 已有非实现者裁决＝0** ／ **B 仅实现方自述＝4**（#1/#2/#4/#6）／
**C 无任何凭据＝1**（#8）／ **D 判据今天不可满足＝2**（#3/#5）。本腿逐格复核并在下面明确写
**哪些是复认、哪些被推翻**。

**锚定尺（本腿现跑，非照抄）**：`grep -nE '^- \[[ x]\]'` 于票面 ＝ **8 格** ——
未勾 **7**：`:82 :83 :84 :85 :86 :87 :95`；已勾 **1**：`:88`（AC#7）。
⚠ 票面 `:10` 那句「`grep -c "^- [ ]"` = 8、已勾 0」**别照抄**：本腿现跑该尺返回 **0**（`[ ]` 是字符类，
不是数框尺），它是 09-20 的过期自述。

**判语取值域**（每格只能落在其一，允许附条件）：
`成立` ／ `不成立` ／ `判据今天不可满足（附具名缺的东西）` ／ `附条件`。
**凭据列**必须区分：**〔我现跑〕**＝本腿本机本轮亲自跑出的读数；**〔读台件，未复跑〕**＝只引用已入库件。
两者不许混。实现方自述（`62-ball-visual-prototype.md`、票面 Progress log）**永远不算裁决凭据**，
只算"待证的说法"。

---

## 1. 逐格裁决表（8 行，一格一行，不合并、不"整体通过"）

| 格 | 票面行号 | 判据（票面原句截断） | 判语 | 凭据 |
|---|---|---|---|---|
| AC#1 | `:82` | 可运行原型 `build\balldebug.exe -stay` 即玻璃液态球，**浅色与深色桌布都清晰可见**（差分截屏为证） | 待裁 | （待填） |
| AC#2 | `:83` | Sleeping：静态单帧、**零活动定时器句柄**（实测断言，非策略表）、CPU 采样 ≤0.5% | 待裁 | （待填） |
| AC#3 | `:84` | 会话态动效五态液体颜色/流速/边框差异**肉眼可辨**，逐态截屏留证 | 待裁 | （待填） |
| AC#4 | `:85` | 喂已知电平（C8 seam wav）时液面旋转幅度随电平**单调变化**；静音收敛进边框过渡态 | 待裁 | （待填） |
| AC#5 | `:86` | 靠边吸附：左/右/上收缩半隐留可命中区；悬停/单击弹回；**跨 DPI 与副屏**正确 | 待裁 | （待填） |
| AC#6 | `:87` | 空闲与说话两口径私有工作集＋CPU 实测回填 `docs/SLO.md`，**阈值不得下调** | 待裁 | （待填） |
| AC#7 | `:88`（已勾） | 契约草案 20 态视觉映射表成文，**等 owner 签字后才回填 SPEC-08** | 待裁（复核已勾是否成立） | （待填） |
| AC#8 | `:95` | 对抗验收由非实现者执行，报告含与上述 AC **1:1** 的裁决表 | 待裁 | （待填） |

---

## 2. 逐格细账（每裁完一格追加一节）

### 2.2 AC#2 —— `:83` Sleeping：静态单帧、零活动定时器句柄（实测断言，非策略表）、CPU ≤0.5%

**判语：附条件。** 三个子句里，"零活动定时器句柄"这一子句**已由我真删一次证成落地**；
"CPU ≤0.5%"这一子句差**一条 D32 原口径（1 分钟窗口）的树外带球读数**。

#### 我现跑的读数〔我现跑〕

**基准（先钉今天绿的用例名）**
```
go test ./internal/ball/ -run '…' -count=1                （默认套，不含 winlive）
  PASS TestDockSquashIsMonotonicAndClamped        0.00s
  PASS TestDockHoverPopBackWalksTheRampHome       0.00s
  PASS TestLiquidRotationIsMonotonicInLevel       0.00s
  PASS TestLiquidConvergesAndBorderFadesInOnSilence 0.00s
  PASS TestAnimationPolicyZeroTimerInSleeping     0.00s
  ok  github.com/CarlosShao/wisp/internal/ball    0.026s

go test -tags winlive ./internal/ball/ -run 'TestLiveSleepingZeroTimerHandles|TestBallLiveAudioLiquidGate'
  PASS TestLiveSleepingZeroTimerHandles           1.78s
  PASS TestBallLiveAudioLiquidGate                1.00s
  ok  2.851s
```
起跑前 `tasklist //FI "IMAGENAME eq go.exe"` ＝ `No tasks are running`（无并发争用）。
`PATH` 已带 `third_party/sherpa-onnx` ＋ `build`。

**发现（先于变异）：唯一的"实测断言"不在默认套里。**
`internal/ball/hotkey_live_test.go:1` ＝ `//go:build windows && winlive`；
`internal/ball/live_windows_test.go:1` 同。⇒ 本票 AC#2 要的**实测**零定时器用例
`TestLiveSleepingZeroTimerHandles`（`hotkey_live_test.go:295`）**被 CI 默认套排除**，
只有显式 `-tags winlive` 才跑。票面 `:83` 写的"实测断言，非策略表"，
在"能跑"这一层成立、在"每次 merge 前自动咬"这一层**不成立**。

#### 我真删了一次（M1，`go test -overlay`，⛔ 未改工作树）

要摘的分支先 grep 证明今天真在：`internal/ball/ball_windows.go:348-353` 逐字
`policy := AnimationPolicy(s)` / `b.stopAnimTimer()` / `b.animTimerActive = policy.PeriodMs > 0` /
`if b.animTimerActive { pSetTimer.Call(… timerAnimID, uintptr(policy.PeriodMs) …) }`。
**`anim.go` 的策略表一个字没动**——这正是这一发的全部要点。

变异体：`internal/ball/ball_windows.go` 的副本里，在 `stopAnimTimer()` 之前插入
`if policy.PeriodMs == 0 { policy = AnimPolicy{AnimFrame, 100} }`
（＝让窗口在策略表不 grant 的态里真的 `SetTimer`）。台件：`.scratch/wisp/probes/62/v1/m1/`。

```
go test -tags winlive -overlay …/m1/overlay.json ./internal/ball/
      -run 'TestLiveSleepingZeroTimerHandles|TestAnimationPolicyZeroTimerInSleeping' -count=1 -v

  hotkey_live_test.go:303: Sleeping received 1 WM_TIMER messages in 600ms (zero-timer discipline broken)
  --- FAIL: TestLiveSleepingZeroTimerHandles          (0.86s)
  --- PASS: TestAnimationPolicyZeroTimerInSleeping     (0.00s)      ← 策略表这发仍然绿
```

**这一发红一绿是 AC#2 措辞的直接凭据**：把"静态度多挂一个真 Win32 定时器"这个缺陷塞进产品代码后，
读策略表的用例**看不见**（仍然 PASS），只有数 Win32 真投递的 `WM_TIMER` 的实测用例**咬住**。
所以票面那句"实测断言，非策略表"不是修辞——这台仪器今天确实分辨得开，
而且 `TestLiveSleepingZeroTimerHandles` 自带正控（`hotkey_live_test.go:308-316`：探针若一条消息都没看见就
`Fatal(…the Sleeping verdict above proved nothing)`，且必须先在 `Warm` 看见过一枚真定时器）。
⇒ **"零活动定时器句柄"子句判 成立（凭我现跑的定向突变）**，附一条登记：它不在 CI 默认套内。

#### CPU ≤0.5%：正面裁那处在册冲突（在册的两枚数各是哪一发）

- **`docs/SLO.md:177` 逐字**「D32 `Sleeping` 行 CPU ≤0.5% 一项 FAIL（0.5450% / 0.6285%），阈值未改、样本未挑」，
  样本行 `:198`（`sleeping-3` 10s/250ms → **0.5450%**）、`:199`（`sleeping-60s` 60s/250ms → **0.6285%**，
  单样本峰值 11.24%）。〔读台件，未复跑〕
  **这一发的主体是 `wisp.exe`（posture=skeleton）＋ 采样器跑在被测树内部**（`SLO.md:206-208`：
  观测者自己 `ReadTree` 的成本被计进 CPU；比例证据 118 reads/60s→0.629%、40 reads/10s→0.545%、
  15 reads/30s→0.087%）。**它量的不是球。**
- **实现侧自述的 `cpu=0.000`** 出自 `docs/evidence/s1/62-diff-signoff/diff-table.txt` 的 `Sleeping` 行
  （`cpu_pct_1core=0.000 / cpu_pct_allcores=0.000`，`timers=no`，`privateWorkingSet_MB=11.76`，`px_delta_ge8=2120`）。
  跑法在 `62-ball-visual-prototype.md:105` 逐字：`build\balldebug.exe -diff … -diff-sample 4s`。〔读台件，未复跑〕
  **这一发是树外读、主体是球本体（`balldebug` 子进程，父进程树外采样）。**
  ⚠ 我的 grep 顺手量到一件对这一格不利的事：那行 `-diff-sample` 是 **4s**，
  而 D32 的 `Sleeping` 行原口径是 **1min 窗口**（`SLO.md:199` 自己写着"= D32「1min 窗口」原长度"）。
  **仓里今天没有一条"球本体 × 树外 × 60s 窗口"的读数**；唯一的 60s 样本恰是那发树内自污染的 `wisp` 跑。

**谁覆盖谁（我的裁定）**：**两个数不互相覆盖，它们答的是两个问题，而且 `docs/SLO.md` 自己已经这么裁过。**
`SLO.md:399-400`（附录 B.4）逐字："不判 PASS、也不判 FAIL，判「仪器未定义」……
CPU 那一行的产品侧证据以 A.2（树外、0.000%）为准，`wisp slo` 的树内 CPU 数字在票 66 修好前只作记录、不作门。
D32 的 0.5% 阈值**一个字没动**"；附录 C 更用同窗口对照把差因坐实（`:415` 区：树外 0.0173% / 树内 0.7629%，
"只有树外那条当门 ⇒ 差的确实就是观测者自己"）。
⇒ 我的结论：**`:177` 那条 FAIL 至今有效、不许被 `0.000` 抹掉**，但它记的是**仪器缺陷（A14/A15）**，
不是"球超预算"；球本体的 CPU 证据只有**4 秒窗口**这一条，**不足以判 D32 那行成立**，
所以 AC#2 的 CPU 子句＝**附条件**（缺的那枚读数具名：`balldebug -diff -diff-sample 1m` 带球、树外、独占安静桌面）。
⛔ 我**没有**改 `docs/SLO.md` 的任何数字或阈值（本腿对它只读；终态名册核对见 §5）。

#### 静态单帧

今天没有一枚用例直接断言"Sleeping 只画了一帧"。它由两条间接撑着：
`applyStateLocked` 在 Sleeping 不 grant 定时器（M1 已证该分支真被咬），且 `renderFrame` 只在
`pushLevel/stepBorder` 返回 `moved` 时发生（`liquid_windows.go:66-69`、`:103-107`）。
⇒ 记为**成立（间接）**，不单列成一枚红。

**AC#2 总判语：附条件。** 凭据构成：定时器子句＝我现跑的定向突变（1 红 1 绿）；
CPU 子句＝读台件（且我现量到窗口口径缺一条 60s 树外带球读数）；静态单帧＝读代码＋间接。

### 2.4 AC#4 —— `:85` 音频驱动：喂已知电平（wav 注入票 13 的 C8 seam）时液面旋转幅度随电平单调变化；静音时收敛进"边框"过渡态

**判语：附条件。** 单调律与"麦克风叫不醒 Sleeping"这道门**都由我真删过、都咬住**；
但**这一格判据点名的注入路径（票 13 的 C8 wav seam）今天在全仓没有生产者**，
现有读数全部来自一枚 debug 二进制里的**合成**电平。

#### 两形突变（硬规矩 #2：不许透过一道会自己拒绝的门去判"拒绝"）

`Ball.applyLevel` 的拒绝不是一个门，而是**两道**（`internal/ball/liquid.go`）：
`liquidDriven(state)`（`:62`，会话族才收电平）与 `liquidMotion.busy(state)`（`:299`，内部先问
`transitionDriven(state)`）。⚠ 只拆第一道会得出**假结论**——`syncLiquidTimer` 仍被 `busy()` 拦住，
计时器根本不会挂，用例照样绿，看着像"仪器不咬"，其实是我只拆了一半。所以 M2 一次拆两道。

台件：`.scratch/wisp/probes/62/v1/m2|m3|m4/`（全部 `go test -overlay`，工作树未动）。

**形 A＝M2（把门那支改成放行：两个 state gate 一律 return true）**
```
go test -tags winlive -overlay …/m2/overlay.json ./internal/ball/ -run TestBallLiveAudioLiquidGate
  live_windows_test.go:640: audio level armed a timer in Sleeping
  --- FAIL: TestBallLiveAudioLiquidGate (0.28s)
```
⇒ 基准（同 tag、无 overlay）该用例 **PASS 1.00s**。Sleeping 里那枚"电平被丢掉"的**拒绝，
确实是这两道门给的**，不是别处给的——这条结论现在才允许入账。

**形 B＝M3（把判定改成恒不成立：`busy()` 永远 false，液体计时器永不挂）**
```
go test -tags winlive -overlay …/m3/overlay.json ./internal/ball/ -run TestBallLiveAudioLiquidGate
  live_windows_test.go:649: a summoned orb in Listening must run the liquid burst timer
  --- FAIL: TestBallLiveAudioLiquidGate (0.29s)
```
⇒ 这一发专治"假绿"：如果 `TestBallLiveAudioLiquidGate` 里 Sleeping 那几条"没有计时器"是因为
**全仓任何态都挂不上计时器**而白捡的，那 M3 之后它会继续绿。它没绿，红在正向断言上
（Listening 召唤态**必须**挂上 burst 计时器）。⇒ **那枚用例有正向对照、量得到东西**，
形 A 的红不是运气。

**M4（单调律本身）**：`liquid.go:224-225` 的 `omega` 去掉 `SpinLevelRadPerS*m.level` 一项。
```
go test -overlay …/m4/overlay.json ./internal/ball/ -run TestLiquidRotationIsMonotonicInLevel
  liquid_test.go:30: rotation travel must grow with the level: level 0.15 gave 0.4455 rad, previous level gave 0.4455
  --- FAIL: TestLiquidRotationIsMonotonicInLevel (0.00s)
```
⇒ 基准 PASS 0.00s。**"随电平单调"这条不是文档里的形容词，是一枚会红的用例**，
且电平一被从角速度里摘掉就立刻平掉（0.4455 vs 0.4455，同一数）。
静音→边框那一半：`TestLiquidConvergesAndBorderFadesInOnSilence` 基准 **PASS**（纯模型、仿真时间，
断言 `border` 在 `BorderOpenMs` 内到位、且 `border` 必须**单调推进**，`liquid_test.go:179`）。

#### 本格今天拿不出的那件东西（具名）

判据逐字写着"喂已知电平（**wav 注入票 13 的 C8 seam**）"。我现查生产者：
```
grep -rln "CarlosShao/wisp/internal/ball" --include=*.go（剥掉 internal/ball 自身与 *_test.go）
  → ./cmd/balldebug/main.go        ← 全仓唯一一个非测试 importer
grep -rn "SetAudioLevel"（非测试）
  → 只有 cmd/balldebug/main.go:418、:477
```
`internal/audio` 侧 C8 的测试骨干在（`internal/audio/wavinjector.go:15` `WavInjector is the C8 test backbone`），
**但仓里没有任何一行把它接到球上**。而 balldebug 用的电平源，flag 帮助文本逐字：
`cmd/balldebug/main.go:106` — `"synthetic audio envelope 0..1 pushed to the ball at ~30fps (0 = feed nothing)"`；
实现 `feedLevels`（`:404-421`）＝ 33ms ticker，每第 12 拍把 `amp` 除以 4 当作"句间呼吸"。
⇒ **今天所有 AC#4 的电平读数都是合成标量，不是 wav 经 C8 seam。** 产品二进制 `wisp.exe` 里
**根本没有球**（零 importer），所以这条端到端路径连"产品侧宿主"都还不存在。
⚠ 这不是"实现方撒谎"——`62-ball-visual-prototype.md` 与票面 Progress log 从没声称跑过 wav；
是**判据点名的凭据类型与盘上有的凭据类型不同名**。⇒ 记 `附条件`，缺的东西具名：
**一条 `WavInjector → RMS/包络 → Ball.SetAudioLevel` 的端到端注入读数**（票面 `:106` 里 62-C 自己
也把 "AC#4 measured with wav injection" 写在 next 栏，即当时未量）。

另有一半仍归人：**"液面旋转幅度"作为肉眼看得见的幅度**（模型里的 `angle` 单调 ≠ 屏幕上的旋涡看得出快慢），
按票 62 处置表 `:26`"未验收 ⇒ 票 65 同批"，与 AC#3 一样要 owner 两发电量对比的眼睛。

**AC#4 总判语：附条件。** 已证：单调律（M4）、静音收敛与边框（用例绿，纯模型）、
状态门是拒绝的真来源（M2 红）且用例不是假绿（M3 红）＝〔我现跑〕。
未证：C8 wav 端到端注入（全仓无生产者）＝具名缺件；肉眼幅度可辨＝归 owner。

### 2.1 AC#1 —— `:82` 可运行原型：`build\balldebug.exe -stay` 启来即是玻璃液态球，浅色与深色桌布上都清晰可见

**判语：附条件（浅色／当前桌布那一半成立，深色那一半今天不可满足）。**

**我现跑**：本腿把 harness 重建成 `.scratch/wisp/probes/62/v1/balldebug.exe`
（`go build -o` 指到 scratch，⛔ 不写 `build/`、不改工作树），
跑 `…balldebug.exe -diff .scratch/wisp/probes/62/v1/diff-selfrun -diff-states Sleeping -diff-sample 60s -diff-margin 120 -x 1720 -y 720`：
```
diff: place=(1720,720) margin=120px cpu-sample=1m0s cores=12 amp=6x out=…/diff-selfrun
diff: Sleeping cpu1c= 0.002% cpuAll= 0.000% privWS= 11.51MB commit=109.91MB ws= 39.43MB
      handles= 349 gdi=4 user=7 timers=no px>=3:2156/px>=8:2101/px>=24:1232 of 97344
      max=193 mean=1.188 box=(136,131 40x45) clean
```
⇒ **差分截屏这条判据我今天由我自己复现成立**：`px>=8` ＝ **2101** 枚变化像素、成像框 **40×45**、
`teardown=clean`。与台件 `62-diff-signoff/diff-table.txt` 的 `Sleeping` 行（`px>=8=2120`、框 44×44、
`privWS=11.76MB`）同量级 ⇒ 不是只有实现代理看得见这个数。
我另用肉眼看了自己刚拍的 `01-Sleeping-alive.png`：球是一颗带亮边、内部蓝紫青渐变的玻璃体，
**背景是"一个浅色对话框压在深色桌面上"的混合底**，球在明暗两块底上都分得开。
⚠ 但**这不算深色桌布那一半**：它不是受控的深色壁纸条件（是别人弹出的窗口），
本腿**不去改 owner 的桌面背景设置**来造这个条件。

**默认档复认**：`cmd/balldebug/main.go:122` 逐字 `ball.EnablePrototypeVisuals(!*frozen)`，
flag 帮助 `:97` 写着 `-frozen` 才回到冻结档 ⇒ 判据那句"启来即是一个玻璃液态球"在 harness 里成立。
⚠ 同一件事的另一面（对 AC#3/AC#7 有用，不是对 AC#1 的扣分，因为 AC#1 点名的就是 balldebug）：
**库默认值是关的**——`internal/ball/tokens_test.go:147-149` 仍断言
`VisualFor(…, Sleeping) == 12px @ 0.35`，注释逐字说那是"the mode the library DEFAULT is in today
(prototypeVisuals off)"，翻这个默认值是票 68 AC#2 的活。⇒ 产品侧今天画的仍是旧微点档。

**今天不可满足的那一半（具名）**：**深色桌布下同法差分一枚都没有**。
我现查：`grep -rniE "dark|wallpaper" docs/evidence/s1/62-diff-*/*.txt` ＝ **零命中**；
`docs/evidence/s1/ball-contrast/` 只有两枚 09-20 11:53 的 PNG
（`same-region-no-ball.png` 4,701B、`sleeping-over-dark-taskbar.png` 4,854B），
文件名自述是**旧微点期对任务栏**的对照，不是本球对深色桌布的差分。
复认 `62-a1` §2.1（它引票面 `:146` 把 dark wallpaper 列进 unverified、台账 A452 同判）。
⇒ **本腿不推翻 `62-a1` 的 AC#1 结论，只把"浅色那一半"从〔读台件〕升成〔我现跑〕。**

**AC#1 总判语：附条件。** 缺的东西具名：**把桌面背景换成深色后重跑一次
`balldebug -diff -diff-states Sleeping`（差分表落新目录，⛔ 别写进 `ball-states/` 档案）**。
这一步要么 owner 换壁纸、要么编排者授权本腿改系统设置——我没有这个授权，所以停在这里。

<!-- NEXT-CELL -->

---

## 3. 我攻不动的地方（本节不许为空）

（待补）

---

## 4. 没做完／留给编排者（本节不许为空）

（待补）
