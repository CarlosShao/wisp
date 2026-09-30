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

| 格 | 票面行号 | 判据（票面原句截断） | **判语** | 凭据（〔我现跑〕／〔读台件，未复跑〕） |
|---|---|---|---|---|
| AC#1 | `:82` | 可运行原型即玻璃液态球，**浅色与深色桌布都清晰可见**（差分截屏为证） | **附条件** | 浅色半侧＝〔我现跑〕60s 树外差分 `px>=8=2101`、框 40×45、`teardown=clean`；深色半侧＝**今天不可满足**（具名：深色桌布下同法差分零枚；本腿不改 owner 桌面设置） |
| AC#2 | `:83` | Sleeping 静态单帧、**零活动定时器句柄（实测断言，非策略表）**、CPU ≤0.5% | **成立**（附一条纪律） | 〔我现跑〕M1 定向突变：实测用例**红**、策略表用例仍**绿**；CPU＝〔我现跑〕D32 原口径 60s 树外 `cpuAll=0.000%`。⚠ 纪律：那枚实测用例在 `//go:build windows && winlive` 后面，**CI 默认套不跑** |
| AC#3 | `:84` | 五态液体颜色/流速/边框差异**肉眼可辨**，逐态截屏留证 | **判据今天不可满足** | 缺三件（全具名）：owner 两张参考图从未入库／五态逐态差分图现只有 **2/5**（〔我现跑〕`ls`）／按参照页判要读 `design/**`＝禁读。〔我现跑〕可分辨性用例 `grep` **零命中**＝无机器代理 |
| AC#4 | `:85` | 喂已知电平（**wav 注入票 13 的 C8 seam**）时旋转幅度随电平**单调**；静音收敛进边框 | **附条件** | 〔我现跑〕M2/M3 两形门突变（红在 `:640`／`:649`）＋ M4 单调律突变（红在 `liquid_test.go:30`）。缺件具名：**C8 wav→球 的端到端注入全仓无生产者**（`internal/ball` 唯一非测试 importer＝`cmd/balldebug`，其电平源自述 synthetic）；肉眼幅度归 owner |
| AC#5 | `:86` | 左/右/上收缩半隐**留可命中区**；悬停/单击**弹回完整球**；**跨 DPI 与副屏**正确 | **不成立** | 律这一半〔我现跑〕成立：M5 红在 `dock_test.go:374`，11 枚 dock 律 `-count=3` 全绿。但真机差分顶反框面：`DockOverlapFrac=0.42`（〔我现跑〕`tokens.go:386`）对上台件 35/44≈**79.5%**（〔我现算〕）；A29③ `r=21<可见 26`。副屏＝〔我现跑〕`SLO.md:181` 单屏 ⇒ 不可满足 |
| AC#6 | `:87` | 空闲与说话两口径私有工作集＋CPU **实测回填 `docs/SLO.md`**，阈值不得下调 | **不成立** | 复测〔我现跑〕两口径齐（60s `Sleeping` 11.51MB/0.000%；4s `Speaking` 15.99MB/0.007%）；阈值未下调〔我现量〕`SLO.md:203/229/230/419/400`；**回填＝〔我现跑〕`grep` 票 62 名下数字在 `docs/SLO.md` 里零命中**（里面那张两口径表署名票 12 的 `12-ball-walk`） |
| AC#7 | `:88`（已勾） | 契约草案 20 态表成文，**等 owner 签字后才回填 SPEC-08** | **成立**（只在框面自述收窄的"成文"半） | 〔我现跑〕结构尺逐枚数过：193 行、§1 表体 **20 行**、§2 冲突 **22 条**、§3 〔待 owner 定〕**15 条**。签字半仍欠；另有一枚待定性：〔我现跑元数据〕`5866c6fc` 在草案后动过 SPEC-08 **1 增 1 删**，与票面 `:74-77` 字面禁令冲突 |
| AC#8 | `:95` | 对抗验收由**非实现者**执行，报告含 **1:1** 裁决表 | **成立**（附自证偏利声明：本格成立≠本票可结案） | 〔我现跑〕本腿 9 枚提交只碰本文件、`internal`/`cmd` 零字节；表 8 行对 8 框；**新造 5 枚红**（M1-M5，每枚先 grep 证分支在）。⛔ 本腿未翻任何框 |

**合计（8 格）**：成立 **3**（AC#2／AC#7／AC#8）｜附条件 **2**（AC#1／AC#4）｜不成立 **2**（AC#5／AC#6）｜
判据今天不可满足 **1**（AC#3）。⇒ **票 62 今天不可结案，`Status: review` 与该保持的未勾框都别动。**

**对前置件 `62-a1` 的账**：**改判 1 格**（AC#5：D 判据不可满足 → **不成立**，理由是盘上已有的真机读数把框面顶反，
不是量不了）；**加强 2 格**（AC#3 补出"逐态截屏只覆盖 2/5"与"无可分辨性用例"；AC#2 补出"唯一的实测断言在
winlive tag 后、CI 默认套不跑"）；**升级凭据 2 格**（AC#1 浅色半侧、AC#6 两口径，从〔读台件〕升为〔我现跑〕）；
**复认不推翻 2 格**（AC#4 的"B 仅自述"本腿把它推进到"律已证、注入路径缺件"；AC#8 的 C 判定＝本表已补上那份缺失的表）。

---

## 2. 逐格细账（每裁完一格追加一节）

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

**AC#2 总判语（初判）：附条件。** 凭据构成：定时器子句＝我现跑的定向突变（1 红 1 绿）；
CPU 子句＝读台件（且我现量到窗口口径缺一条 60s 树外带球读数）；静态单帧＝读代码＋间接。

#### 更正（同日稍后，本腿自跑补齐了上面具名缺的那枚读数）——**AC#2 总判语改：成立（附一条纪律）**

初判写完之后，本腿把 harness 建到 scratch 并**自己跑了那条缺的口径**（命令与全行读数见 §2.1）：
```
diff: place=(1720,720) margin=120px cpu-sample=1m0s cores=12
diff: Sleeping cpu1c= 0.002% cpuAll= 0.000% privWS= 11.51MB … handles= 349 gdi=4 user=7 timers=no … clean
```
⇒ **D32 原口径（1 分钟窗口）× 树外读 × 球本体**这一格今天有了〔我现跑〕的数：
**全核 0.000%（1 核 0.002%）**，判据是 ≤0.5%，**PASS**；同一行 `timers=no` 是读活对象，
与 M1 那发突变互相咬得上。这条补齐后，"CPU 子句缺 60s 树外读数"这个理由不再成立，
所以本腿把 AC#2 从 `附条件` 改判 `成立`。

**两条必须跟着走的边界（不是客套，是这一格下次被读时会被问到的）**：
1. **`:177` 那条 FAIL 没有被这一发覆盖、也不该被覆盖**。它量的是 `wisp.exe` 树内采样器（A14/A15 仪器缺陷），
   `docs/SLO.md` 自己在 B.4（`:399-400`）与附录 C 已裁过"树外那条当门、树内只作记录"。
   本腿**没有**动 `docs/SLO.md` 任何字节（终态名册核对见 §5）。
2. **测量环境不是"独占安静桌面"**：跑这一发的 09:06–09:07 本机 self-hosted CI 在跑（编排者 00:0x 的推送，含 `slo-full`），
   且截图里桌面上有**别人弹出的一个模态对话框**压在采样区上。⇒ 对**每进程**口径（CPU%/私有 WS/句柄/GDI）
   争用只会把数字往不利方向推，读数仍是 0.000%，所以这一发可用；
   但它**不满足** `SLO.md` A.6 那套"桌面独占＋前后 `tasklist` 无残留"的仪式，
   所以它**不能顶替票 12/66 那条全量门**，只能顶替"AC#2 的 CPU 子句今天没有 60s 树外读数"这个具体缺口。

**仍然挂着的那一条纪律（判"成立"不抹掉它）**：唯一的实测零定时器用例
`TestLiveSleepingZeroTimerHandles` 在 `//go:build windows && winlive` 后面，
**CI 默认套不跑**（本腿现量：不带 tag 时 `-run` 它 ＝ 一条用例都不出现）。
⇒ 判语 `成立` 的范围是"今天有一枚会真咬的实测断言，且我用突变证明它咬得住"，
**不是**"每次 merge 前自动咬"。这一条要不要补进 CI，是编排者的决定，不是我能替它勾的。



### 2.3 AC#3 —— `:84` 会话态动效：五态液体颜色/流速/边框差异肉眼可辨，逐态截屏留证

**判语：判据今天不可满足。**（与 `62-a1` 的 D 同向，但本腿把"缺什么"量得更死，见下面第 3 条。）

**1. "肉眼可辨"这一半：仓里没有任何机器代理。**〔我现跑〕
`grep -rniE "distinguish|distinct|可辨" --include=*_test.go internal/ball` ＝ **零命中**。
最接近的 `TestVisualForCoversAllTwentyStates`（`internal/ball/tokens_test.go:119`）只断言
每态 `SizePx>0 && 0<Opacity<=1`（**覆盖性**，不是**可分辨性**），而且它读的是**库默认档**
（`:147-149` 注释逐字："read the mode the library DEFAULT is in today (prototypeVisuals off)"，
并当场断言 `Sleeping == 12px @ 0.35`）。⇒ 没有任何一枚用例会因为"两态长得一样"而红，
所以这一半**不是我能用突变去证的**，它按票面处置表 `:25` 归票 65＋owner 眼睛。

**2. 参考图与参照页两条路本腿都走不通（具名）。**
① owner 那两张玻璃参考图（vivo 蓝心小V 语音球）**从未入库**——`62-a1` §2.3 已核，票 65 头部现读
`Status: blocked-on-owner（缺参考图）`；② 按规格参照页判要读 `design/**`（SPEC-08 §2.1 引的那页），
**`design/**` 在本编队的禁读清单里，连文件名都不许列**，所以本腿**没有去看**，直接具名标不可满足。

**3. "逐态截屏留证"这一半：票 62 名下只覆盖了 5 态里的 2 态。**〔我现跑〕
```
ls docs/evidence/s1/62-diff-{baseline,glass,border,signoff}/*-alive.png
  baseline / glass : Sleeping, Listening, Speaking          （3 态）
  border           : Sleeping                                 （1 态）
  signoff          : Sleeping, Sleeping-dock-right,
                     Listening, Listening-dock-right,
                     Speaking, Speaking-dock-right            （6 枚 alive 图）
```
⇒ 判据点名的五态 **Listening / Thinking / Acting / Warm / Speaking** 里，
票 62 名下有图的只有 **Listening、Speaking**；**Thinking、Acting、Warm ＝ 零枚**。
旁路我核过并**不采信**：`docs/evidence/s1/12-ball-walk/` 那 7 态（Sleeping/Listening/Thinking/Acting/
Speaking/Warm/Settling）确实齐，但 `grep -inE "frozen|prototype|liquid|glass" 12-ball-walk/diff-table.txt`
＝ **零命中**——那张表**没记自己是哪一档渲染**，所以不能算票 62 新液态视觉的逐态凭据（它是票 12 的
资源差分）。`docs/evidence/s1/ball-states/` 那 20 枚本腿 `ls -la` 复认：mtime **全部 09-19 21:38**
＝旧微点期档案（`scripts/dev/ball-cycle.ps1:19` 的默认 OutDir 正指着它，⛔ 谁要重跑必须带 `-OutDir`，
这条 `62-a1` §3 已钉，本腿复认）。

**4. 还有一层会让"逐态可辨"今天根本看不到的事**：`62-visual-spec-draft.md` §0.2 自己写着
**库默认值＝冻结档**（`prototypeVisuals=false`，全仓唯一开启者是 `cmd/balldebug/main.go:122`），
并引 registry **A24-D2**"owner 签收的球不在默认构建里"。⇒ 五态液体差异今天只在 debug harness 里存在；
在默认构建里画的是旧 §2.1 那 20 行。这一条不扣 AC#1（AC#1 点名 balldebug），
但它是 AC#3 这句"会话态动效……肉眼可辨"在产品侧**无法被看到**的直接原因。

**要判这一格，今天缺的东西（逐件具名）**：
(a) owner 手上那两张玻璃参考图入库（或 owner 当场口头判"像不像"）；
(b) 五态各一枚**新液态视觉**的差分截屏（现只有 2/5），且落**新目录**；
(c) 若要按参照页判：给本编队解禁 `design/**` 的读权限，或由 owner 那侧的界面助手代判。
⛔ 本腿**没有**读 `frontend/**`／`design/**` 任何字节来凑这一格。

**AC#3 总判语：判据今天不可满足**（缺 (a)(b)(c) 三件，全部具名如上）。对 `62-a1`：**不推翻，加强**——
它记的"62 名下无裁决表／参考图未落盘／要读 design"三条本腿全部复现，
本腿另量出"逐态截屏只覆盖 2/5 态"与"没有任何可分辨性用例"这两条它没写的。

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

### 2.5 AC#5 —— `:86` 靠边吸附：拖到左/右/上边缘收缩半隐并留可命中区；悬停/单击弹回完整球；跨 DPI 与副屏行为正确

**判语：不成立（作为一票之框）。** 这一格要拆成三块才说得清，而**只有第一块今天可判成立**；
`62-a1` 给这一格打的 **D（判据今天不可满足）本腿改判为「不成立」**——
"不可满足"说的是量不了，而这一格的问题是**盘上已有的量测直接把框面顶反了**。这是本腿对 `62-a1` 的**一处改判**。

#### 块一：弹回／收缩的**步进律**——今天可判成立，且我用突变证明它咬得住〔我现跑〕

**M5**（台件 `.scratch/wisp/probes/62/v1/m5/`，`go test -overlay`）：
先 grep 证明要摘的分支真在——`internal/ball/dock.go:44` 逐字
`next = clamp01(approach(p, target, float32(dt)/float32(DockAnimMs*time.Millisecond)))`。
把它换成 `next = clamp01(target)`（＝台阶直接跳到终点、不走路）：
```
go test -overlay …/m5/overlay.json ./internal/ball/ -run TestDockHoverPopBackWalksTheRampHome
  dock_test.go:374: primary/left: the pop-out took 1 frames of ~33 ms, budget is DockAnimMs=160
  --- FAIL: TestDockHoverPopBackWalksTheRampHome (0.00s)
```
基准（无 overlay）该用例 **PASS 0.00s**。⇒ "丝滑弹回"不是一句形容词，
是一枚**数帧预算**的断言，把"瞬移"塞进产品代码它就红。
整组 dock 律我另跑了 `-count=3` 隔离复量（硬规矩：计时类红要两发齐）：**11 枚 × 3 次全 PASS**、
零翻转 ⇒ 这一组**不是**计时偶红：
`TestDockGeometryShape / TestDockTangentMeansGapZero / TestDockPosKeepsTheOrbTangentAtEveryRampLevel /
TestDockSquashIsMonotonicAndClamped / TestDockProgressReadsThePush / TestNearestDockEdgePicksTheEdgeYouPushed /
TestDockedLeavesAHittableStrip / TestDockFreePosKeepsTheBallReachable /
TestDockOnSecondaryMonitorWithNegativeOrigin / TestDockHoverPopBackWalksTheRampHome / TestDockRampWalksItsBudget`。
DPI 那一半在**律**上也有覆盖：`TestHitTestAndDPIInjection`（`internal/ball/tokens_test.go:351`）
把 96/120/144/192 四档各测中心命中／角部穿透／"窗随 DPI 变大"，整包 `-count=1` 跑出它 PASS。

#### 块二：真机差分把框面的"留可命中区／收缩半隐"顶反（在册两处，本腿独立复算）

- **A29②（停靠露出 ≈82% 而非 42%）**：契约值我现量在 `internal/ball/tokens.go:386`
  ＝ `DockOverlapFrac = 0.42 // fraction of the orb left visible when docked`。
  拿台件 `62-diff-signoff/diff-table.txt` 自己的两行做除法：未停靠 `Sleeping` 成像框 **44** 宽，
  停靠 `Sleeping-dock-right` 成像框 **35** 宽 ⇒ **35/44 ≈ 79.5%**，与在册的"≈82%"同量级，
  与契约的 42% 差一倍。〔我现算；输入是台件数字，未复跑真机〕
  ⚠ 口径边界照写：`bbox_w` 是"变化像素外接框"，不是"球体露出宽度"的严格同义词，
  所以这一发是**复认在册读数**、不是把它测成定案。
- **成因在草案里已经写死**（`62-visual-spec-draft.md:42` 逐字）：
  `DockSquash(1)=DockOverlapFrac=0.42`，但**"只压液斑/高光/rim 的 X 向；shell/halo/caustic 仍是整圆"**，
  且 `DockPos` 的 `orbR` 取的是**窗口**半径。⇒ 这不是量错，是**挤压没作用到玻璃壳**，
  所以"半隐"在屏幕上不成立。
- **A29③（命中 21px < 可见 26px）**：同一张草案表 `:45` 逐字
  "`r=32 vs 可见 36/42` 与 `r=21 vs 可见 26` 是**同一个缺口的两个读数**"，并说 registry 只记了 Sleeping 那格、
  "同形排查应扩展到全球包"（⇒ 列为 Q-2，**不自行改值**）。
  ⇒ **律与真机在这里分家**：`TestDockedLeavesAHittableStrip`（`dock_test.go:183`，断言
  命中条 ≥ 可见球宽）今天 PASS，而真机差分说它反了。所以"律全绿"不能顶这一格。
- 修不修**等 owner 在 R15#4/#5 拍板**（`62-a1` §2.5 已记，本腿复认）；本腿**没有**自行改任何值。

#### 块三：跨 DPI 与副屏——今天不可满足（具名）

本机**物理单屏**：`docs/SLO.md:181` 逐字"32GB RAM · 单屏 3440×1440"。
⇒ 判据"副屏行为正确"缺的东西＝**一块第二显示器（或 owner 授权用虚拟屏）上的真拖拽**。
`TestDockOnSecondaryMonitorWithNegativeOrigin` 是**负原点几何 mock**，按票 64 AC#6 当年的判例
（其验收表明令"没拿双屏 mock 冒充"），本腿**不把它当硬件凭据**。
"跨 DPI"那一半同样缺**真高 DPI 桌面上的差分**：现有全部真机证据都是 96 DPI
（`62-visual-spec-draft.md` §0.3 表头逐字"默认配置 56 / 96 DPI"）。

#### 顺手钉一枚与本票无关的红（防下一位归因错）

整包 `go test ./internal/ball/ -count=1` 今天 **54 PASS / 1 FAIL**，那枚 FAIL 是
`TestC21TableColourRowsMatchTokensCSS`，报错逐字：
`read design/assets/tokens.css: … The system cannot find the path specified. - the CSS leg of this check must never skip`。
本腿查起手名册第 14 行：` D design/assets/tokens.css`（该文件在 HEAD 里存在，`git cat-file -e` 现测），
**是别人未提交的工作树删除**，与本票无关、与本腿的 overlay 突变无关（这一发我没带 overlay）。
⇒ 别把它记成"票 62 名下一枚红"。

**AC#5 总判语：不成立。** 块一成立（M5 红＋11×3 绿）；块二被在册两处读数顶反（本腿复算其一）；
块三今天不可满足（缺第二块屏／缺真高 DPI 桌面）。⛔ 本腿未改任何值、未翻任何框。

### 2.6 AC#6 —— `:87` 资源复测：空闲与说话两种口径的私有工作集＋CPU 实测回填 `docs/SLO.md`，阈值不得下调

**判语：不成立（"回填"那一半从未发生；另两半成立）。** 这一格拆三截：
**复测＝成立（本腿现在两种口径都有我自己的数）**、**阈值未下调＝成立（我现量）**、
**回填 `docs/SLO.md`＝不成立（票 62 名下的数不在那文件里）**。

#### 两种口径我各跑了一发〔我现跑〕

同一枚 scratch 二进制、同一台机器、同一法（球在子进程、父进程树外读）：

| 口径 | 采样窗 | CPU 1核 / 全核 | 私有工作集 | 句柄 | GDI/USER | timers | px≥8 | 成像框 | teardown |
|---|---|---|---|---|---|---|---|---|---|
| 空闲（`Sleeping`） | **60s**（D32 原窗长） | 0.002% / **0.000%** | **11.51MB** | 349 | 4/7 | **no** | 2101 | 40×45 | clean |
| 说话（`Speaking`，`-diff-level 0.9`） | 4s | 0.078% / **0.007%** | **15.99MB** | 380 | 4/10 | yes | 4848 | 72×72 | clean |

台件：`.scratch/wisp/probes/62/v1/diff-selfrun{,-speaking}/diff-table.txt`。
对照实现方台件 `62-diff-signoff/diff-table.txt`（`Sleeping` 11.76MB/0.000%/px≥8=2120；
`Speaking` 13.18MB/0.001%/px≥8=4848）：`px≥8` 那一列我复跑到 **4848 与台件逐字相同**，
私有 WS 我这次高 2.8MB（09:06 本机 self-hosted CI 在跑，含 `slo-full`，且两口径都仍 ≤25MB）。
⇒ **"空闲与说话两种口径的私有工作集＋CPU 实测"这一截今天不再只有实现方自述**。

#### 回填那一半：不成立（我现查）

```
grep -n "11.76\|13.18\|14.95\|13.16\|62-diff" docs/SLO.md   →  零命中
```
`docs/SLO.md`（我现量 45,546B／583 行）里与球有关的附录只有三本：**附录 A（票 12 的 S1 空闲跑）**、
**附录 B（编排者复跑 A.4 两起仪器缺陷）**、**附录 C（票 66 修好后的复跑）**。
A.2 那张表（`:217-224`）确实**同时有 `Sleeping` 11.79MB/0.000% 与 `Speaking` 13.11MB/0.001%**
——两口径齐，但**出处是 `docs/evidence/s1/12-ball-walk/diff-table.txt`，署名票 12**，
不是票 62 的复测，也不是那四套 `62-diff-*`。
⇒ 按框面字面（"实测回填 `docs/SLO.md`"），**票 62 名下从未回填过一行**。
⛔ 本腿**没有**动 `docs/SLO.md` 一个字节（终态名册核对见 §5），也不建议由我这条腿来补这一行——
它要的是"谁的数、哪一发、什么口径"署得清，那是编排者的活。

#### 阈值不得下调：成立（我现量文本，且与票面自述互相对得上）

`docs/SLO.md:203` 七门写 `内存 ≤25MB`、`:229` 私有工作集门 `≤25MB`、`:230` CPU 门 `≤0.5% 全核`、
`:419` 逐字"`cpu_percent_all_core` 的 limit 在下列每张表里都写作 `<=0.5%`，与 D32 16.3.2 原文一致，未改"、
`:400` 逐字"D32 的 0.5% 阈值**一个字没动**"。
拿票面自己的数对：`:44` "CPU ≤0.5%"、`:53` "预算 ≤25MB" ⇒ **文档侧没有任何一处被下调**。
⚠ 口径边界（写进 §3，别让它被当成我没做的事）：**仪器里那枚常量我没验**——
`thresholds.go` 在本腿的禁读清单上，所以"文档写的阈值＝代码里的阈值"这一跳**我没核、也无权核**。

#### 还有一条对本格不利的在册事实（带走，别丢）

台账 `A451` 逐字：今晚 CI 的 `slo-full` 步正文 0 行、工件缺失 ⇒ "球常驻的资源达标"这份 CI **给不了依据**，
且"绝不许写 'slo-full 绿＝D32 达标'"。〔读台件，未复跑——CI 是编排者那侧的，本腿不重跑别人的台件〕
⇒ 本格今天能靠的只有上面那两发树外差分＋`SLO.md` 附录 C 的树外当门口径。

**AC#6 总判语：不成立。** 缺的东西具名：**把本腿（或任一非实现腿）这两行带署名的读数
回填进 `docs/SLO.md` 一个新附录**，并补一条**独占安静桌面**下的同法跑（本腿这一发有 CI 争用、
桌面还有别人的模态对话框，够判 AC#2 的 CPU 子句、不够顶 D32 全量门）。

### 2.7 AC#7 —— `:88`（**唯一一枚已勾 `[x]`**）契约草案：20 态视觉映射表成文，等 owner 签字后才回填 SPEC-08

**判语：成立——但只在框面自己收窄的那个口径上成立**（框下注释逐字："本框只勾'成文'，
不暗示票 62 已验收"）。⛔ **本腿没有翻、也没有改这一枚框，一行都没动**；下面是给那一勾做的外部核。

#### 我现量它的内部结构（不是复述编排者在票面 `:89-92` 说的"我核过"）〔我现跑〕

`docs/evidence/s1/62-visual-spec-draft.md` ＝ **45,087B／193 行**（`wc -l` 现量，与票面 `:29`/`:89`
说的"193 行"逐字对上）。入库锚点 `48f0cfd2`（`git log -1 --format=%ci` ＝ 2026-09-21 08:45:39），
subject 逐字："docs(62): 补写 AC#7 那份从未存在的契约草案 62-visual-spec-draft.md（20 态视觉映射，等 owner 签字）"。

| 框面声称 | 我用的尺 | 现量 |
|---|---|---|
| §1 逐态一行 | `sed -n '65,94p' \| grep -c "^\|"` 再扣表头＋分隔行 | 22 ＝ 2＋**20 行**，**不多不少** |
| §2 逐条冲突 | `sed -n '96,125p' \| grep -c "^\|"` | 24 ＝ 2＋**22 条**（`**C-01**`…`C-22`） |
| §3 15 个〔待 owner 定〕 | `grep -c "待 owner 定"` | **15** |
| §4 明列"证据里找不到数字的格子" | 标题行现读 | 在 `:148`，标题逐字带"直说找不到，不编" |
| §6 不动票面 Status 与勾选框 | 标题行现读 | 在 `:180`，标题逐字带"（**不动票面 Status 与勾选框**）" |

⇒ **"成文"这一半我认，且认得比 `62-a1` 硬**：`62-a1` §2.8 只复认了"193 行"这一个数，
本腿把 20／22／15 三枚结构数各自数过一遍，全部对上，没有"标题在、内容空"的形状。

#### 但本格的后半句牵出一条要编排者定性的事（本腿不自己判它违规）

票面 `:74-77`「契约变更流程（本票特有）」逐字：**"在签字之前，任何提交都不得修改
`docs/specs/SPEC-08-ui-ball-panel.md` 或 `docs/PLAN.md`"**。我现查（**只用文件级元数据，
`docs/specs/**` 的内容一个字节都没读**）：
```
git log --numstat 48f0cfd..HEAD -- docs/specs/SPEC-08-ui-ball-panel.md
  5866c6fc | 2026-09-24 09:58:06 +0800 | docs(Q-39..Q-43 执行): 缺口文档入仓带作废声明＋两处过期计数改对＋…
  1  1  docs/specs/SPEC-08-ui-ball-panel.md
```
⇒ 草案入库**之后、owner 签字之前**，SPEC-08 被改过**恰好一次、1 增 1 删**。
两种可能，本腿分不出来（分它要读 SPEC-08，那是禁读）：
**甲**＝那枚提交的 subject 自述的"两处过期计数改对"，与 §2.1 的 20 态视觉表无关
⇒ 框面后半句（未回填）仍然干净，但**字面禁令被越过了一次**，该补一条批准记录；
**乙**＝它动的正是 §2.1 那张表 ⇒ **这就是未经签字的回填＝契约违规**，AC#7 那一勾要降格。
⛔ 本腿**不替编排者选甲乙**，也不去读 SPEC-08 凑一个答案。这一条同时进 §4「留给编排者」。

**AC#7 总判语：成立（限"成文"半，框面自己声明的口径）**；签字半仍欠（R15 第 10 项 Q-1…Q-15）；
另附一枚待定性：`5866c6fc` 对 SPEC-08 的那 1/1 行改动是甲还是乙。

### 2.8 AC#8 —— `:95` 对抗验收由非实现者执行，验收报告含与上述 AC 1:1 的裁决表（README 规则 6）

**判语：成立（附一条自证偏利的声明）。**

**这一格判据的对象就是本文件本身**，所以它天然自证——本腿先把"我既是被判者又是判者"这件事写在明面上：
**本格成立不等于票 62 可结案**。AC#3 今天不可满足、AC#5 不成立、AC#6 不成立，
那三格不因本表存在而变绿；本票的 `Status: review` 与本框的未勾状态**都该保持**。
⛔ **本腿没有翻这一枚框，也没有改任何一行框面文字**（框归编排者，见 §5 名册核对）。

#### 判据的三个词，逐个给我现跑的凭据

**① "由非实现者执行"** —— 可核，不靠我自我声明：
```
git log --name-only --format="%h" d5a59d66..HEAD | （剥掉 commit 行后按路径计数）
   9  docs/evidence/s1/62-adversarial-acceptance.md
   1  docs/reports/HANDOVER.md          ← 不是本腿的
   1  docs/reports/pending-and-issues.md ← 不是本腿的
```
本腿自锚点 `d5a59d66` 起的**全部提交只碰这一枚证据文件**（口径与推导式：
`git log --name-only --format="%h" d5a59d66..HEAD` 里属于本腿的那些行，路径去重后只有本表一条；
上面那次现量的时点是 9 枚，本表每再追加一枚提交此数 ＋1，**别把这个数当常量读**）；
`internal/**`、`cmd/**`
**一个字节都没写**（每一发突变之后 `git status --porcelain -- internal cmd` 均为空，逐次读数见 §5）。
⇒ 票 62 的实现码本腿零行，"非实现者"这一条成立。
⚠ 顺带记一句给编排者：这是**共享工作树**，锚点之后另外那两枚提交（HANDOVER／台账）不是我提交的，
别把它们算进本腿的产出。

**② "1:1 的裁决表"** —— §1 那张表 **8 行对 8 枚框**（`:82 :83 :84 :85 :86 :87 :88 :95`），
一格一行、**没有合并、没有"整体通过"**；每行都有判语与凭据列，凭据一律标〔我现跑〕／〔读台件，未复跑〕。

**③ "对抗"** —— 本腿**新造 5 枚红**，每枚都落在指名的那一条断言上，且每枚之前都先 grep 证明要摘的分支今天真在：

| 突变 | 摘掉的东西（先 grep 证其在） | 指名的用例 | 结果 |
|---|---|---|---|
| M1 | `ball_windows.go:348-353` 给策略表不 grant 的态硬塞真 `SetTimer`（`anim.go` 不动） | `TestLiveSleepingZeroTimerHandles` | **红**（`hotkey_live_test.go:303`，1 枚 WM_TIMER）；同跑的**策略表**用例仍**绿** |
| M2 | `liquidDriven` ＋ `transitionDriven` 两道门一起放行 | `TestBallLiveAudioLiquidGate` | **红**（`live_windows_test.go:640`） |
| M3 | `liquidMotion.busy()` 恒 false（计时器永不挂） | 同一枚用例 | **红**（`:649` 正向断言）⇒ 反证那几条"没有计时器"不是假绿 |
| M4 | `omega` 里去掉 `SpinLevelRadPerS*m.level` | `TestLiquidRotationIsMonotonicInLevel` | **红**（`liquid_test.go:30`，0.4455 vs 0.4455） |
| M5 | `dockRampFrame` 直接 `clamp01(target)`（台阶瞬移） | `TestDockHoverPopBackWalksTheRampHome` | **红**（`dock_test.go:374`，1 帧 vs 预算 160ms） |

另加两发自跑测量（§2.1／§2.6：60s 与 4s 两口径树外差分）。
⛔ 本腿**没有**为了让任何一格变绿而放宽过断言或阈值；5 枚突变全在 `go test -overlay` 里做，
工作树未改（`.scratch/wisp/probes/62/v1/m{1..5}/` 是临时件，按规矩只建不删）。

**AC#8 总判语：成立**——三个词逐个有凭据，但**它的成立只覆盖"这份报告存在且形状对"这一件事**，
不覆盖本票其余七格的结论。

---

## 3. 我攻不动的地方（本节不许为空）

下面每一条都是**本腿试过或明确判过"我这条腿做不了"**的，写清做不了的原因；
它们**不是**"实现方说做不了"，也不是"我没去做"。

1. **深色桌布下的差分（AC#1 的后半句、AC#3 的一半前提）**——做不了。
   `balldebug -diff` 的成像判据是"同一块桌布上 alive 减 dead"，所以它**只能量当下那张壁纸**。
   要拿深色那一发，得把 owner 的桌面背景换成深色。**本腿没有这个授权，也没有动系统设置**，
   所以停在这里。我这一发的现场背景顺带记一下：采样区里压着**别人弹出的一个模态对话框**
   （浅底）＋深色桌面边缘，是混合底，**不能冒充受控深色条件**。
2. **"肉眼可辨"（AC#3 全部、AC#4 的后半句）**——今天没有机器代理，我攻不动。
   〔我现跑〕`grep -rniE "distinguish|distinct|可辨" --include=*_test.go internal/ball` ＝ **零命中**。
   没有任何一枚用例会因"两态长得一样"而红，所以它不在突变能覆盖的范围里；
   判它要么 owner 的眼睛，要么按 `design/**` 那张参照页——**后者在本编队的禁读清单内，
   本腿连文件名都没列**，所以直接具名标不可满足，没有绕过去看。
3. **owner 手上那两张玻璃参考图（AC#3）**——不在盘上，我造不出来。票 65 头部现读
   `blocked-on-owner（缺参考图）`；`62-a1` §2.3 同判，本腿复认。
4. **副屏与真高 DPI（AC#5 后半句）**——物理不可满足。〔我现跑〕`docs/SLO.md:181` 逐字"单屏 3440×1440"。
   `TestDockOnSecondaryMonitorWithNegativeOrigin` 是**负原点几何 mock**，按票 64 AC#6 的判例
   （当年明令"没拿双屏 mock 冒充"）本腿**不把它当硬件凭据**；现有全部真机差分都是 96 DPI。
5. **`thresholds.go` 在禁读清单上（AC#6 的"阈值不得下调"只能核到文档层）**——
   我核的是 `docs/SLO.md` 写的门与票面自述互相对得上；**"文档写的阈值＝仪器里的常量"这一跳我没核、也无权核**。
   同理 `golden`、`allowlist.txt`、三枚冻结件（`internal/panel/tokens_fourway_test.go`／
   `internal/panel/l2_grant_boundary_test.go`／`internal/perm/ticket90_persist_test.go`）本腿一个字节没读。
6. **`docs/specs/**` 禁读 ⇒ AC#7 那枚 SPEC-08 改动我定不了性**——我只拿到文件级元数据
   （`git log --numstat`：`5866c6fc`，1 增 1 删）。是"过期计数改对"还是"未签回填"，
   **分它必须读 SPEC-08 的内容**，所以本腿把它作为待定性交给编排者，没有猜。
7. **CI 那侧的读数（AC#6 的 `slo-full`）**——`A451` 记的"今晚 `slo-full` 正文 0 行、工件缺失"
   属〔读台件，未复跑〕：本腿**不重跑别人台件**（会洗掉被引用的读数），只把它当约束带走。
8. **`-race` 那一族我没有用作证据**——本腿全部读数是 `-count=1`／`-count=3`，没跑 `-race`。
   台账 `A458`／票 237 登记的那个坑（`internal/tools/subagent_197.go:347` 起孩子排在
   `:368` 父发派生回执之前）在 `internal/tools` 包，与本票无关；
   本腿撞到的**唯一一枚整包红**是 `TestC21TableColourRowsMatchTokensCSS`，
   成因是**工作树里 `design/assets/tokens.css` 被别人删了**（起手名册第 14 行 ` D design/assets/tokens.css`，
   该文件在 HEAD 里存在），**不是计时偶红、也不是票 62 的缺陷**（详见 §2.5 末节）。
   按硬规矩"整包名册红＋单包 `-count=3` 隔离复量"两发我都做了：整包 54 PASS/1 FAIL →
   隔离到那 1 枚具名红 → dock 组 `-count=3` 全绿零翻转。
9. **测量环境不干净这件事我无法排除**——本机 self-hosted CI 正在跑（编排者 00:0x 推送，含 `slo-full`），
   所以我那两发树外差分**不满足** `SLO.md` A.6 的"桌面独占＋前后 `tasklist` 无残留"仪式。
   它们够判 AC#2 的 CPU 子句与 AC#6 的"两口径有数"，**不够顶 D32 全量门**。

---

## 4. 没做完／留给编排者（本节不许为空）

**本腿的产出边界**：只交这份表＋台件；**一枚 AC 框都没碰**（不翻勾、不改框行、不写"已勾"），
`docs/SLO.md` 一字未动，`internal`／`cmd` 零字节。下面每一格要往下走，都需要框面之外的动作。

1. **AC#6（不成立）——回填那一行没人写**。票 62 名下两口径读数至今不在 `docs/SLO.md` 里
   （〔我现跑〕`grep "11.76|13.18|14.95|13.16|62-diff" docs/SLO.md` 零命中）。
   建议：把 §2.6 那两行（或一条独占安静桌面下的新跑）**带署名**追加成 `SLO.md` 一个新附录。
   ⚠ 这条**不该由验收腿自己写**——写它等于自己把判据补成成立。
2. **AC#7——请定性 `5866c6fc` 对 SPEC-08 的那 1 增 1 删**（甲＝过期计数修正，补一条人工批准记录；
   乙＝动了 §2.1 的 20 态表 ⇒ 未签回填＝契约违规，AC#7 那一勾要降格）。本腿禁读 `docs/specs/**`，不猜。
3. **AC#2——那条纪律要不要补**：唯一的实测零定时器用例在 `winlive` tag 后、CI 默认套不跑。
   要么把它纳入某种自动门，要么在票面明写"这条今天只在被显式调用时才咬"的代价。
   本腿判它"成立"的范围是**今天有一枚会真咬的断言且我用突变证明它咬得住**，不是"每次 merge 前自动咬"。
4. **AC#4——判据点名的注入路径与盘上有的不同名**：全仓没有 `C8 WavInjector → Ball.SetAudioLevel` 的生产者，
   产品二进制里根本没有球（`internal/ball` 唯一非测试 importer 是 `cmd/balldebug`）。
   两条路：**(a)** 补那条端到端注入；**(b)** 把 AC#4 的文字改成"合成电平"——
   ⚠ (b) 是**改 AC 判据＝人工批准**，本腿不动。
5. **AC#5（不成立）——A29②／A29③ 修不修等 owner 在 R15#4/#5 拍板**；
   另外这一格把"律／真机／硬件"三件事捆在一枚框里，导致它永远只能整体不成立。
   **拆框＝改判据＝人工批准**，本腿不拆，只把三块的判语分别写在 §2.5。
6. **AC#1——深色那一发**：需要 owner 换一次深色壁纸（或书面授权本腿/别的腿改系统设置）后重跑
   `balldebug -diff`，**落新目录**。⛔ 别写进 `docs/evidence/s1/ball-states/` 档案
   （`scripts/dev/ball-cycle.ps1:19` 的默认 OutDir 正指着它，`62-a1` §3 已钉覆盖风险，本腿复认）。
7. **AC#3——三件缺物**（§2.3 末具名）：参考图入库／五态逐态差分图（现 2/5）／`design/**` 解禁或由
   owner 那侧界面助手代判。**这一格今天签不了**，别用"票 65 临时通过"顶它（票面 `:14-15` 已自证那不是裁决）。
8. **票面状态**：本票 8 格里 2 枚不成立、1 枚不可满足、2 枚附条件 ⇒ **`Status: review` 保持，
   `-done` 后缀不要加**。⛔ 别再犯 09-20 那一次 premature `-done`（票面 `:8-16` 自述、台账 A30①）。
9. **临时件（按规矩只建不删）**：全部在 `.scratch/wisp/probes/62/v1/` ——
   `m1/`…`m5/`（各含 overlay 副本＋`overlay.json`）、`diff-selfrun/`、`diff-selfrun-speaking/`、
   `balldebug.exe`（scratch 自建，未写 `build/`）、`_roster_start.txt`／`_roster_end.txt`、
   `build_mutations.py`／`fix_order.py`／`sort_sections.py`。
   ⚠ 那三份 `.py` 与本表无关，但 `m1..m5` 是 §2 里每枚"红"的可复跑台件，**别删**。

---

## 5. 名册与禁区核对（本腿终态自检）

**起手名册**（跑任何尺之前先存，`.scratch/wisp/probes/62/v1/_roster_start.txt`）＝ **151 行**。
**终态名册**（`_roster_end.txt`）＝ **151 行**，`diff` 起手与终态 ＝ **零差异**。
⇒ 本腿**没有**往工作树里留任何未跟踪/已修改文件（本腿那几枚提交已入库，故不出现在名册里；
scratch 临时件在忽略路径下）。终态判据按派单是"与起手逐字相同、只多出本腿自己的件"——
本腿这一发是**完全相同**（多出来的那几件都在忽略路径／已提交）。

**逐次突变后的定点核对**（硬规矩 #1 要求留读数）——**口径写清楚，别读成"每发各核一次"**：
M1 之后单独跑过一次 `git status --porcelain -- internal cmd`（空）；
M2/M3/M4/M5 是在**同一轮**里连跑的四发，之后核过一次（空）；此后每一枚提交前又各核过一次，
**全部为空**。⇒ 准确说法是"**每一轮之后一次**"，四发突变共用一次核对，**不是** 5 次一对一。
全程 `go test -overlay`，没有改过工作树（`.scratch/wisp/probes/62/v1/m{1..5}/` 是 overlay 副本）。

**禁区自证**：
- `git status --porcelain -- internal cmd docs/SLO.md` 终态 ＝ **空** ⇒ 实现码与 SLO 数字一字节未动。
- **AC 框**：票面 `grep -nE '^- \[[ x]\]'` 本腿收尾复跑 ＝ 仍 **8 格、未勾 7（`:82 :83 :84 :85 :86 :87 :95`）／已勾 1（`:88`）**，
  与起手逐字一致 ⇒ 没有翻勾、没有改框行。
- 零读零引零枚举 `frontend/**` 与 `design/**`：本表出现的唯一 `design/**` 字样是
  `design/assets/tokens.css` 这个**路径名**，它来自**名册元数据**（`git status` 的删除行）与被禁文件的
  **存在性**判断（`git cat-file -e`），**不是内容读取**；本腿没有 `ls design/`、没有打开过其中任何文件。
- 未读 `PLAN.md`／`docs/specs/**`／`thresholds.go`／golden／`allowlist.txt`／三枚冻结件（对 SPEC-08 只用了
  `git log --numstat` 的文件级计数）。
- 未放宽任何断言或阈值；未 `git add -A`/`.`（9 枚提交全部显式 pathspec `docs/evidence/s1/62-adversarial-acceptance.md`）；
  未 `--amend`/`reset`/`rebase`/`stash`/`checkout .`/`clean`；**只 commit、未 push**。
- 本表与台件里**没有任何 API 密钥或凭据值**。
