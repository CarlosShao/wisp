# 机主一次 3 分钟窗口 — 操作单（票 290／291／449／247 同窗收口，314 具名单开）

**这张纸给谁**：机主本人。他不写代码，只看屏、按键、签收。
**派单存档**：`.scratch/wisp/probes/orch/20261115/txt/owner-3min-runbook-dispatch.txt`
**本件性质**：只读腿产出。零产码、零票面改动、零台账改动、零 push；写面只有本文件。
**跑本（步骤的机器面）**：`.scratch/wisp/probes/runbook/capture-once.ps1` ＋ `.scratch/wisp/probes/runbook/levelprobe/main.go`（11-14 22:35／22:27 的现成件，本窗口**原样用，不改一行**）
**采到哪**：`D:\tmp\wisp\20261115\realonce\`（`summary.md`／`cells.json`／`facts.json`／`probe\levelprobe.json`／`probe\levelprobe.frames.jsonl`）

---

## ① 一句话开场

**这次你只做 9 个动作（真正按下去的只有 5 次按键组），约 3 分钟，一次收口 4 张票：290 `AC#3②`、291 `AC#2`、449 `AC#8`、247 `AC#2`。**
**第 5 张票 314 这次不做，必须单独约一次 —— 它要的是「你离开机器 23 分钟」，和这 3 分钟里「你正坐在这儿按键、说话」互相打架（数字见 ⑤ 末）。**

---

## ② 动作表（一行一个动作；秒数是"从你坐下算起"的约数）

| 第几秒 | 按什么／说什么 | 屏幕该出现什么 | 出现＝哪张票哪一格翻 | 没出现＝记什么 |
|---|---|---|---|---|
| 0:00 | 打开一个 PowerShell 窗口（`Win+X` → 点「终端」），把这一行整个粘进去、按回车 | 先出现 7 行：`wisp real-machine once - capture run` / `repo : …` / `data dir : C:\Users\swq\AppData\Roaming\wisp-dev` / `run dir : D:\tmp\wisp\20261115\realonce` / **`ball : pid 49348 (wisp)`** / `config : mic_muted_default=true voice.enabled=true` / `sink : …\logs\wisp-20261115-001.jsonl (… lines)`；紧接着出现黄字 `【第 1 步 / Step 1】请按一次「静音热键」（ticket 290）… 按完再按回车。` | **290 `AC#3①`（`290-1a`）**：`capture armed but muted at boot; device not opened`＋`config_source=…\wisp-dev\config.toml`（这格我 17:5x 已经从盘上读到了，读数在 ④ 第 1 条；这里出现只是把它绑到本次 run） | `ball : ABSENT - no run-marker.json` ⇒ **整窗作废**：18 格全记 `ABSENT`，note 逐字写 `no run-marker.json`（`capture-once.ps1:281`）。⛔ 不要自己去找球、不要双击 `wisp.exe` 重开一个（那是第二个球，会把热键抢走）。 |
| 0:20 | **按一次「静音热键」**（就是你平时叫它开麦／关麦那个组合；键位记不清就见本行末），**按完心里数到一（约半秒到一秒）就立刻按回车** | 屏上 `  按回车继续（Enter when done）` 之后出现 `正在编译电平探针（go build，约 2-5 秒）…`，再出现 `【第 2 步 / Step 2】接下来 17 秒…`；**同时**日志里多一行 `msg="mute gesture turned this process's capture gate"`、`gesture="mute-hotkey"`、`outcome="已取消静音：设备已交接，采集线程在跑；…"`（日志不是屏幕，我核） | **290 `AC#3②` 的 `290-2a`（那行 outcome）＋`290-2b`（`gesture=mute-hotkey`）** | 三种成因分开记，⛔ 不合并：**①** 日志窗内没有任何 `mute gesture` 行 ⇒ `290-2a`=`ABSENT`、note 逐字 `no 'mute gesture' record in <窗> of <sink>`（`capture-once.ps1:397-400`）；**②** 有行但 `outcome` 以 `已静音：采集已关闭，设备未打开（再按一次取消静音）` 开头，说明门本来是开的 ⇒ `290-2a` 记 `ABSENT`＋`the gate was already open at step 1`；**③** 有行但 `gesture="tray-mute"` ⇒ `290-2a` 可有、`290-2b` 记 `ABSENT`（票面点名要 `mute-hotkey`）。⚠ 键位忘了**不要换托盘菜单**：托盘那支手势词是 `tray-mute`，`290-2b` 就采不到。 |
| 0:21 | 按回车（就一下） | 同上那一行 `  按回车继续（Enter when done）` 消失，进入 `【第 2 步 / Step 2】` | 无（推进） | 卡住不动 ⇒ 记 `290-2a`／`290-2b`／`291-*`／`449-*`／`247-*` 全 `ABSENT`，note `script stalled before step 2`，并把我这一发的 `run dir` 也记下（空目录也是读数）。 |
| 0:30 | **不出声、不碰键盘鼠标**（4 秒） | `=== quiet-1: 完全安静：别说话、别碰键盘（4 秒） ===`，然后 4 秒后自动出现 `=== speak: 用正常音量说一句话，一直说到窗口结束（5 秒） ===` | `291-1`（`291 AC#2` 形①）、`449-a`／`449-d`（`449 AC#8`）；`449-b`／`449-c` 逐帧也算在窗内 | `probe\levelprobe.json` 里该窗 `samples=0` ⇒ `291-1`／`449-a`／`449-d` 记 `ABSENT`，note 逐字 `levelprobe.json window 'quiet-1': 0 samples`（`capture-once.ps1:467-469`）。若是 `probe_exit_code=2` ＋ 屏上 `MICROPHONE DID NOT OPEN` ⇒ **12 格全 `ABSENT`**，note 记那一句原文＋「设备被独占／无权限」。 |
| 0:34 | **用平时说话的音量说一句话，一直说到下一行出现**（5 秒，别喊、别贴麦） | `=== speak: ……（5 秒） ===`；5 秒后自动出现 `=== quiet-2: 再安静一次，别出声（4 秒，对照组） ===` | `291-2`（形②，`291 AC#2`）、`291-5`（形①能不能被既有 2× 门槛分开）、`247-a`（`247 AC#2`：说话那一发与安静那一发不同） | 该窗 0 帧 ⇒ `291-2`／`291-5`／`247-a` 记 `ABSENT`（`291-5` 缺前件，⛔ 不许凭空算倍数）。若是"说了但和安静一样" ⇒ **不是 `ABSENT`**，`291-5` 记 `FILLED-NOT-SEPARABLE`（票面明确允许），⛔ 不许调门槛凑绿。 |
| 0:39 | **再不出声 4 秒**（对照组） | `=== quiet-2: ……（4 秒，对照组） ===`；4 秒后出现 `=== tone: 什么都别说：本程序现在自己放一段 1 kHz 音（4 秒） ===` | `291-1b`（`291 AC#2` 对照窗） | 该窗 0 帧 ⇒ `291-1b` 记 `ABSENT`。若它比 `quiet-1` 明显更响 ⇒ 照样记读数，note 逐字 `room got noisy between quiet-1 and quiet-2`（`capture-once.ps1:475` 只判 span，这句要手写）。 |
| 0:43 | **什么都不做 4 秒**（程序自己用扬声器放 1 kHz 音；⛔ 别说话） | `=== tone: ……（4 秒） ===`，然后出现 4 行读数 `levelprobe: quiet-1 samples=… new(mean=… max=…)` … `levelprobe: tone …`，再出现 `【第 3 步 / Step 3】再按一次同一个静音热键…` | `291-3`（形③，`291 AC#2`）、`247-c`（端点名／帧率／放音幅度，`247 AC#2`） | 屏上 `levelprobe: tone playback failed: …` ⇒ `291-3` 记 `ABSENT` ＋ 那句原文；若放成功但麦听不到（`tone` 的 `max` 与 `quiet-1` 同级）⇒ `291-3` 记 `READ`、note 逐字 `外放这枚麦听不到`（这是票面允许的答案，见 ④ 第 8 条）。 |
| 0:52 | **再按一次同一个静音热键**，按完数到一就按回车 | 日志多一行 `msg="mute gesture turned this process's capture gate"`、`gesture="mute-hotkey"`、`outcome="已静音：采集已关闭，设备未打开（再按一次取消静音）"`；屏上出现 `captured -> D:\tmp\wisp\20261115\realonce` ＋ `  cells: N READ, M ABSENT, 18 total` | `290-2d`（`290 AC#3②` 的另一半：两个方向都动过） | 记 `290-2d`=`ABSENT`，note 逐字 `no 'mute gesture' record in <窗>` ＋ `WARNING: the mute gesture back was not observed`（`capture-once.ps1:603-606`）⇒ **麦可能还开着**，这一条必须原样进结论，⛔ 不许写成"关回去了"。 |
| 0:56 | **收工，什么都不用按** | 屏上最后三行：`captured -> <目录>`、`  cells: N READ, M ABSENT, 18 total`、`  summary.md / cells.json / facts.json` | 全部 18 格落盘（含 `ABSENT`） | 屏上没有 `captured ->` ⇒ 脚本中途死了；把屏上最后几行**拍照或抄下来**，记 `run_rc` 与最后一行原文，别的格全 `ABSENT`。 |

> 键位：`config.toml` 的 `[hotkey] mute = ''`（＝用编译默认值）；球启动那一行逐字印 `hotkeys:"summon=Ctrl+Alt+Q live, mute=… live, cancel=…, panel=Ctrl+Alt+P live"`（11-14 16:23:09 那行，`wisp-20261114-001.jsonl:72`）。**本单不印键位本身**；忘了就按 `Ctrl+Alt+Q` 把球叫出来看一眼托盘。

---

## ③ 每张票一格「他做完我怎么核」

> 一律**可跑命令 + 读数判据**；⛔ "看日志"这种话不写。
> 统一变量：`$run = 'D:\tmp\wisp\20261115\realonce'`（若目录已存在会是 `realonce-001`，以屏上 `captured ->` 那行为准）。
> **翻框归非实现者**（`issues/README`／`SPEC-12 §4.3`）：下面只给读数与出处，翻不翻由编排者裁。

### 290 — 翻 `AC#3②`（`290-2a`＋`290-2b`＋`290-2d`）

```powershell
$run = 'D:\tmp\wisp\20261115\realonce'
Get-Content "$run\cells.json" -Raw | ConvertFrom-Json |
  Select-Object -ExpandProperty cells |
  Where-Object { $_.id -like '290-*' } |
  Format-Table id,status,value,source -AutoSize
# 权威复核（脚本那格只有 1 秒宽，见 ④ 第 10 条）：按本次 run 的时间窗扫 sink
$f = 'C:\Users\swq\AppData\Roaming\wisp-dev\logs\wisp-20261115-001.jsonl'
$facts = Get-Content "$run\facts.json" -Raw | ConvertFrom-Json
$t0 = [DateTimeOffset]::Parse($facts.started_at); $t1 = [DateTimeOffset]::Parse($facts.finished_at)
Select-String -Path $f -Pattern 'mute gesture' | ForEach-Object {
  $o = $_.Line | ConvertFrom-Json
  $ts = [DateTimeOffset]::Parse($o.time)
  if ($ts -ge $t0 -and $ts -le $t1) {
    [pscustomobject]@{ line = $_.LineNumber; time = $o.time; gesture = $o.gesture; outcome = $o.outcome }
  }
} | Format-Table -AutoSize
```

- **判过**（三条同时成立）：① 上面那张表在本次 run 的时间窗内**恰好 2 行**；② 一行 `outcome` 以 `已取消静音：设备已交接，采集线程在跑` 开头且 `gesture=mute-hotkey`；③ 另一行 `outcome` 以 `已静音：采集已关闭，设备未打开` 开头。
- **判不过**：只有 1 行（缺一个方向）、或两行 `gesture` 是 `tray-mute`、或 0 行。**0 行时还要看并发**：若本次窗口期间任何时刻 `Get-Process wisp.test,wisp-winlive.test` 不为空，这一发记 `ABSENT-CONTENDED`（键位被别的测试进程占着，`290-2a` 不作数），不是"他没按"。
- 附：`290-1a`（`AC#3①`）**不靠本次 run** —— 现成读数在 ④ 第 1 条（`wisp-20261114-001.jsonl:82`）。

### 291 — 翻 `AC#2`

```powershell
$run = 'D:\tmp\wisp\20261115\realonce'
$j = Get-Content "$run\probe\levelprobe.json" -Raw | ConvertFrom-Json
$j.windows | ForEach-Object { [pscustomobject]@{win=$_.name; n=$_.samples; mean=$_.new_level_dc_stripped.mean; max=$_.new_level_dc_stripped.max; span_s=$_.span_s; per_s=$_.levels_per_s} } | Format-Table -AutoSize
$j.synthetic | Format-Table name,measured,expected,within_tolerance -AutoSize
$q = ($j.windows | Where-Object name -eq 'quiet-1').new_level_dc_stripped
$s = ($j.windows | Where-Object name -eq 'speak').new_level_dc_stripped
"ratio={0:N2}x  speak_max>quiet_max={1}" -f ($s.mean/$q.mean), ($s.max -gt $q.max)
```

- **判过**（四形齐 + 一个允许的"分不开"）：`quiet-1`／`speak`／`quiet-2`／`tone` 四窗 `samples>0` 且 `span_s ≥ 3`（`capture-once.ps1:475` 的判据）；`synthetic` 六枚 `within_tolerance=true`；`ratio > 2.0` **且** `speak_max > quiet_max`（＝`FILLED-SEPARABLE`）。**`ratio ≤ 2.0` 时写 `FILLED-NOT-SEPARABLE` 也算这一格填完**（票面 `291:17` 逐字允许"分开不了就具名说分开不了"），⛔ 不许挑一个能过的门槛。
- **判不过**：任一窗 `samples=0`、或 `span_s < 3`、或 `synthetic` 有一枚 `within_tolerance=false`（后者记 `READ-DISAGREES`，直接进 ⑤ 的失败登记）。⚠ `quiet-1.mean` 若恰好是 `0`（尺去直流后恒定帧读 0），那行 `$s.mean/$q.mean` 会**除零**报错 ⇒ 这本身就是答案：`FILLED-NOT-SEPARABLE` 判不出来，改记 `ratio=undefined (quiet_mean=0)`，⛔ 不许把 0 换成一个小数去算倍数。

### 449 — 翻 `AC#8`（判据已更正为「同一窗 `new ≤ old`」＋「`mean < 0.06`」，票面 `449:86`）

```powershell
$run = 'D:\tmp\wisp\20261115\realonce'
$fr = Get-Content "$run\probe\levelprobe.frames.jsonl" | Where-Object { $_.Trim() } | ForEach-Object { $_ | ConvertFrom-Json }
$d  = $fr | ForEach-Object { [double]$_.new - [double]$_.old }
"max(new-old)={0:E3}  frames_new_gt_old={1} of {2}" -f ($d | Measure-Object -Maximum).Maximum, @($fr | Where-Object { [double]$_.new -gt [double]$_.old }).Count, $fr.Count
$dc = $fr | ForEach-Object { [double]$_.dc_counts }
"dc_mean={0:N1}  dc_max_abs={1:N1} counts" -f ($dc | Measure-Object -Average).Average, ($dc | ForEach-Object { [math]::Abs($_) } | Measure-Object -Maximum).Maximum
$best=0;$n=0; foreach($x in $fr){ if($x.phase -match '^quiet'){ if([double]$x.new -lt 0.06){$n++}else{$n=0}; if($n -gt $best){$best=$n} } else {$n=0} }
"longest_run_below_gate={0} frames ({1} ms)" -f $best, ($best*32)
# mean<0.06：见 ③ 291 那段的 quiet-1 mean
```

- **判过**：`max(new-old) ≤ 0` 且 `frames_new_gt_old=0`；`quiet-1` 的 `mean < 0.06`；`longest_run_below_gate ≥ 5` 帧（`SilenceHoldMs=150` ÷ 32 ms ＝ 4.69 ⇒ 第 5 帧够，`internal/ball/liquid.go:31`）。
- **判不过**：出现 `frames_new_gt_old>0`（数学上不该有；有＝尺坏了，直接进 ⑤）；或 `mean ≥ 0.06`（球会永远判"有人在说话"，`liquid.go:272`）；或连续帧数 < 5。
- ⚠ 安静窗 `max` 单帧越过 `0.06` **不算判不过**（一声碰键就够），判决看 `mean` ＋ 连续帧数；尖峰要在 note 里点名，⛔ 不许抹平。⛔ 票面那枚 `0.368` 今天复现不出来（见 ④ 第 4 条），只引用、不当判据。

### 247 — 翻 `AC#2`（只能翻一半；`247-b` 结构性缺）

```powershell
$run = 'D:\tmp\wisp\20261115\realonce'
$j = Get-Content "$run\probe\levelprobe.json" -Raw | ConvertFrom-Json
$j.endpoints | Format-List
$j.tone | Format-List frequency_hz,amplitude_counts,digital_rms_fullscale
($j.windows | Where-Object name -eq 'quiet-1').new_level_dc_stripped | Format-List mean,max
($j.windows | Where-Object name -eq 'speak').new_level_dc_stripped  | Format-List mean,max
Get-Content "$run\cells.json" -Raw | ConvertFrom-Json | Select-Object -ExpandProperty cells | Where-Object { $_.id -like '247-*' } | Format-Table id,status,value -AutoSize
```

- **判过**：`quiet-1` 与 `speak` 的 `mean` **不同**（`247-a`）＋ `endpoints.capture` 是**真设备名**（不是占位／虚拟）＋ `windows[tone].levels_per_s ≈ 31.25`（帧率）＋ `tone.digital_rms_fullscale = 0.353553390593274`（`247-c`）。
- **这一格只能翻一半**：`247-b`（球侧 `levels_delivered > 0`）**本窗采不到** —— 那行只在进程退出时打（`cmd/wisp/resident_audio_windows.go:479` 现读；票面旧引 `:429` 已漂移，具名）。⛔ 不许拿 `290-2c` 的历史行充当本次读数。

### 314 — 本次**不翻**，单独一次窗口（见 ⑤ 末「时长账」）

```powershell
# 到时候要跑的那一发（现在⛔ 不要跑：会占热键、会抢前台）
Get-Content 'docs\evidence\s1\314-clean-single-run.md' | Select-String 'CLEAN-FULL-'
```
- **判过**：同一发里 `CLEAN-FULL-RC`／`PASS/FAIL/SKIP` 三数 ＋ 终点锚 `TestZZ315A1ProjectInstructionsSeedCannotStrandTheFirstTask` 有 `PASS` 行，**并且**这一发跑的时候 `Get-Process LogonUI` 存在（＝真锁屏，不是"以为锁了"）。
- **为什么现在不能算**：11-14 那发 `LogonUI` 前后都是 0（`rundll32` 锁屏在后台进程里静默失败，`314-clean-single-run.md:55`）⇒ 那一发的"安静桌面"根本没成立。

---

## ④ 提前排除（我自己就能判的，逐条给读数 —— 他不用为这些动一下）

> 下面每条都是 17:5x 现跑／现读，出处在每行末尾。**这些格外面已经能填，窗口只是复核。**

1. **`290-1a`＝`290 AC#3①`（默认档不偷开麦）已经读到，他不用做任何动作。**
   读：`C:\Users\swq\AppData\Roaming\wisp-dev\logs\wisp-20261114-001.jsonl:82` 逐字
   `"time":"2026-11-14T16:23:12.0130392+08:00" … "msg":"audio: capture armed but muted at boot; device not opened","config_source":"C:\\Users\\swq\\AppData\\Roaming\\wisp-dev\\config.toml"`；
   在跑的进程 `PID 49348` 启动 `2026-11-14 16:23:02`（`Get-Process 49348`）；`config.toml` 里 `[audio] mic_muted_default = true` 逐字在。（出处：`wisp-20261114-001.jsonl:82`、`%APPDATA%\wisp-dev\config.toml` 的 `[audio]` 节）
2. **`291-4`（形④六枚合成帧）已经读到——纯代码，不需要麦、不需要他。**
   我 17:5x 现跑（`D:\tmp\wisp\20261115\preflight-probe-001\`，probe `RC=0`）：`constant-zero=0`／`constant-dc-12345=0`／`square-1lsb=3.0517578125E-05`／`square-fullscale=0.999969482421875`／`hottest-legal=0.999984741210938`／`dc12345-sine6000=0.12947651924196`（期望 `0.129475118625466`，容差 `SineLevelTolerance=5e-5`，`internal/audio/level.go:89`）—— **六枚全部 `within=True`**。
3. **`449-b`（新尺 ≤ 旧尺）已经读到**：同一发 `levelprobe.frames.jsonl` **106 帧**，`max(new-old) = 0.000E+000`、`frames_new_gt_old = 0`。
4. **`449-c`（今天有没有那枚直流）已经读到，结论是"没有"**：`dc_mean = -0.7 counts`、`dc_max_abs = 161.7 counts`；而票面那枚 `0.368462` 折算是 **`12073.76 counts`**（`0.368462 × 32768`）⇒ **`0.368` 今天复现不出来**，`AC#8` 只按"同窗 `new ≤ old` ＋ `mean < 0.06`"判，这必须在结论里具名。
5. **`449-d`（连续低于门限的帧数）已经读到**：我的安静窗 **30 帧连续 `< 0.06`＝960 ms**，门槛只要 **5 帧＝160 ms**（`liquid.go:31` `SilenceHoldMs=150`）。
6. **`449-a` 只算半格**：我的 `quiet-2` `mean=0.001233`（`< 0.06`）、`quiet-1` `mean=0.005520` 但 `max=0.065971`（我正坐在键盘前敲命令，一帧尖峰）。⇒ 正式那一格以他 4 秒不动那段的 `quiet-1.mean` 为准，判据不变（`< 0.06`）；我这发只用来证明"链路通、麦能开、数字落在量级上"，不据它判 `449-a`。
7. **`247-c`（端点／帧率／放音幅度）已经读到**：`capture=麦克风 (Realtek(R) Audio)`、`render=扬声器 (Realtek(R) Audio)`、`capture_thread=52272`、`levels_per_s` 实测 `31.00`（帧长 32 ms ⇒ 理论 `31.25`，`internal/audio/audio.go:47-53`）、`tone.digital_rms_fullscale = 0.353553390593274`。
8. **`291-3`（形③外放）我今天试了一发小声的**：`tone` 窗 `mean=0.001226 / max=0.001297`，与 `quiet-2`（`0.001233/0.001443`）同级 ⇒ **外放这枚麦听不到**。⚠ 我那发只放了 **0.3 秒（14 帧）**，不据此判；正式那发是 4 秒。若正式那发仍同级 ⇒ `291-3` 记 `READ` ＋ note「外放听不到」，这就是票面允许的答案。
9. **`290-2c` 与 `247-b` 是结构性缺，不是失败**：两格要的 `levels_delivered` 只在 `stop()` 里打（`cmd/wisp/resident_audio_windows.go:479` 现读逐字 `audio capture leg stopped: levels_delivered=%d frames_sent=%d frames_dropped=%d reopens=%d`），而本窗口**明令不许停产品进程**（`real-machine-once.md` §4）。补它的代价＝**一次有序退出**（托盘 Exit／Ctrl+C），另开一发 30 秒，不在本窗。
10. **脚本自身那格手势窗只有 1 秒宽 —— 具名的缺陷，不是他按错。**
    `capture-once.ps1:213-221` 的 `Wait-Human` 把 `PromptAt` **在 `Read-Host` 返回之后**才取，于是接受窗＝`[回车−1 s, 回车+2 s]`；而日志刷盘周期 `500 ms`（`internal/observe/logging.go:52` `flushInterval = 500 * time.Millisecond`）⇒ 按完键到回车之间必须 **≥0.5 s 且 ≤1 s**，那行才既落盘又落在窗内。**这就是本单在 ② 写"数到一就回车"的来路**；并且 ③ 里给了**按 run 时间窗的权威复核**：脚本那格写 `ABSENT` 而 sink 有行时，以复核为准（有逐字 `文件:行号`）。
11. **第 2 步"产品门开着时探针还能开麦"这条前提，现读成立**：产品采集走共享模式（`internal/audio/wasapi_windows.go:285` `audclntShareModeShared`＝`0`）⇒ 探针与产品可同时持有同一枚采集端点。若哪天改成独占，第 2 步会 `MICROPHONE DID NOT OPEN`——那时 12 格全 `ABSENT`，不是他说话太小。
12. **今天日志落到哪个文件，现读钉死**：writer 按 **UTC 日**切（`internal/observe/logging.go:314-334`，`logDayLayout="20060102"`），UTC 日切＝北京 **08:00** ⇒ 现在写进 `wisp-20261115-001.jsonl`（现 2857 行，末行 `17:11:44`，里面**没有**任何 `mute gesture` 行、**没有**今天的开机行）⇒ 今天的 `mute gesture` 行只可能来自本次窗口。
    ⚠ 同时读到一个**污染面**：这个 sink 是**所有进程共用**的（同日第二个进程 `O_APPEND` 续写同一个 `-001`，`logging.go:344-348`），今天 2857 行里 **2604 行是测试台件 `MARK299`** 写的 ⇒ 窗口期间若别的测试进程在跑，它的 `mute gesture` 行会伪装成他那一按。**所以开窗前我自己先清场**（见下一条）。
13. **开窗前我（不是他）要做的清场**：`Get-Process wisp.test,wisp-winlive.test` 必须为空、`Get-Process wisp` 恰好 **1 枚且 pid=49348**（17:49 现读有一枚 `wisp.test` PID 64600，命令行 `…\.scratch\wisp\stage-cli\wisp.test.exe -test.count=1 -test.timeout=90m -test.v -test.run=Test46`＝别的腿的单枚定向跑；窗口前必须等它退，因为它一样会占热键、一样往同一个 sink 写字）。
14. **在跑的那枚球是旧镜像，449 的读数不靠它**：盘上 `wisp.exe` 已在 `2026-11-15 15:46` 重编（内嵌 `vcs.revision=e108b0945`，含 `f9b6d60f7`（`449` 落地，`2026-11-14 22:14:29`）——`git merge-base --is-ancestor f9b6d60f7 e108b0945` 返回 0）；而**在跑的 PID 49348 起于 11-14 16:23**，早于那笔落地 ⇒ 它进程里那枚尺还是旧的。**449 的四格由探针（HEAD 代码）给**，我不声称"在跑的那枚产品已经用上新尺"；要把那一条也变成事实，得重启球＝另一窗（本窗⛔ 禁）。
15. **本次我自己复核过的钉**：`go test ./internal/audio -count=1 -run 'TestLevelDc|TestLevelTranslationInvariant|TestLevelFullScale|TestLevelConstant|TestAC297|TestLevelSine' -v` ⇒ **6 PASS / 1 SKIP / 0 FAIL**（SKIP＝`TestAC297LivePacketByteFace`，要真包），`ok … 0.632s`；行号现读（每枚都是现跑 `Select-String` 出来的，⛔ 不是抄旧件）：`LevelOfSamples`=`internal/audio/level.go:119`、`LevelFullScale=:67`、`LevelLSB=:72`、`SineLevelTolerance=:89`、`SilenceLevelGate=internal/ball/liquid.go:30`、`SilenceHoldMs=:31`、`LevelReleaseTauMs=:33`、判"安静"那一行 `liquid.go:272`、帧长 `internal/audio/audio.go:47`（`FrameSamples=512`）与 `:53`（`32ms`）、2× 门槛 `cmd/wisp/resident_audio_247_live_windows_test.go:194` 逐字 `if loud.mean() <= quietA.mean()*2 || loud.max() <= quietA.max()`。

---

## ⑤ 失败也记（每个动作"没出现"时，记什么、票面怎么留）

**总规矩**：`ABSENT` 是"没采到"的**正式答案**，⛔ 不许改成 0、⛔ 不许补一个看起来合理的数（`real-machine-once.md` §1 填表铁律）。**⛔ 本单没有"再试一次"这一栏** —— 每个失败都写清"下一发要先把哪枚前提拿掉"，拿不掉就具名挂着。

| 动作 | 没出现时的记法（逐字） | 票面怎么留 | 下一发要先拿掉的前提（不是"再按一次"） |
|---|---|---|---|
| 0:00 脚本起不来 | `290-1a`＝`ABSENT`，note `no run-marker.json`（或 `run-marker names a pid that is not running`，`capture-once.ps1:275`） | `290 AC#3②` 与本窗全部读格**留空**；`AC#3①` 用 ④ 第 1 条那枚既有读数，⛔ 不因本窗失败而退回 | 先把球起回来（**编排者的动作**，`⛔` 不是让机主双击第二个） |
| 0:20 按了热键但日志没有 `mute gesture` | `290-2a`=`ABSENT`，note `no 'mute gesture' record in <窗> of <sink>`；**另记一行** `wisp.test 竞争进程：<pid 列表或 none>` | `290 AC#3②` 留空，票面写「`290-2a` 欠一格；本发 note＝<逐字>」 | 清场（`wisp.test`/`wisp-winlive.test`＝0）；⛔ 不换托盘菜单凑（那会变 `tray-mute`，`290-2b` 照样缺） |
| 0:20 有行但 `outcome` 是 `已静音…` | `290-2a`=`ABSENT`，note `the gate was already open at step 1` | 同上留空；**另记**「本发现门起手是开的」——这本身是读数，进结论 | 窗口前先用托盘菜单确认门在静音位（编排者做，⛔ 不占他的手） |
| 0:30／0:39 安静窗 0 帧 | `291-1`／`291-1b`／`449-a`／`449-d`＝`ABSENT`，note `levelprobe.json window '<name>': 0 samples` | `291 AC#2`／`449 AC#8` 留空（`449 AC#8` 可**部分**成立：`449-b`／`449-c` 来自逐帧，仍可填，须具名"只有两格成立"） | 麦被独占（会议软件）或权限；`probe_exit_code=2` 原文进 note |
| 0:34 说了但分不开 | `291-5`＝`FILLED-NOT-SEPARABLE`（**不是** `ABSENT`），并抄 `report ratio=<n>x speak_max=… quiet_max=…` | `291 AC#2` **可以翻**（票面允许）；`247 AC#2` 只记 `247-a` 的实数（两个数不同就算到），⛔ 不据此判 247 失败 | 麦增益／说话距离；`291` 编排者裁定节引 `247-v1` 件 `:44-45` 已量过"增量 `0.09917` < 两窗漂移 `0.10178`"，同形再出现就照录（⛔ 不许挑门槛） |
| 0:43 外放麦听不到 | `291-3`＝`READ` ＋ note `外放这枚麦听不到（tone max 与 quiet 同级）` | `291 AC#2` 照样可翻（票面允许"具名说分不开"）；⛔ 不许写"外放确实能被听到" | 扬声器音量／路由（`real-machine-once.md:25` 的 P3 一条已写明"形③量的是外放到底能不能被这枚麦听到"，11-08 那次就是这一形分不开——`291` Progress log 11-08 17:04 那条） |
| 0:52 第二按没有 `已静音` | `290-2d`=`ABSENT`，note `no 'mute gesture' record in <窗>` ＋ `WARNING: the mute gesture back was not observed` | `290 AC#3②` **只翻一半**（开门有、关门无）⇒ 留空；**结论里必须写"麦可能还开着"** | 门态核对（编排者用托盘菜单看勾）后再择窗 |
| 任何格 `READ-DISAGREES` | 原样抄 `measured` vs `expected`，**不改任何期望值** | 该票对应格**不翻**，并在票尾追加一行「<时刻> 一发真机读数与码不符：<逐字>」 | 先查是不是尺真的坏了（`449-b` 越界就是这一形），⛔ 先动断言 |

### 时长账（为什么是 3 分钟，以及 314 为什么不是）

- **本窗实测量级**：脚本侧 = 预检 ~1 s ＋ 探针编译 2–5 s ＋ 第 1 步（按＋回车）~10 s ＋ 探针四窗 **17 s**（4+5+4+4）＋ 第 3 步 ~10 s ＋ 写盘 ~1 s ⇒ **约 46–60 秒机器时间**；连他坐下、读屏、按，**3 分钟内收工**。他按下去的**只有 5 次**：一行命令、两次热键、两次回车。
- **314 必须单开，三个数字**：
  1. **23 分钟量级**：`cmd/wisp` 整包一发 **`1442 s`**（`314-clean-single-run.md:20`），另一发 **`1389.8 s`**（`314-ac0-forensics.md` 第 2 发）⇒ 是 **3 分钟窗口的 7.7–8 倍**，塞不进来。要 **1–2 发** ⇒ **23–46 分钟**。
  2. **它的条件与本窗互斥**：314 `AC#2` 要的是「**安静桌面**（⛔ 人不用机器）」(`314:32`)，而本窗要的正是「你坐在机器前按键、说话」——同一个人在同一段时间里不可能既在场又离场。
  3. **它会抢走本窗的第 1 个动作**：整包会起真窗用例并占住 `Ctrl+Alt+Q/M/P` 三枚全局热键（`314-ac0-forensics.md`「驻留球占热键」那两条红），11-14 已经撞出过 `chord ownership: held by another process (ERROR_HOTKEY_ALREADY_REGISTERED)` 这一串（`290` Progress log 11-14 10:38 那条）⇒ 两件事同开，`290-2a` 会采不到而**看起来像他没按**。
- **314 那一次的形态**（给编排者排，不占本单）：他做 **1 个动作**（按 `Win+L`），然后**离开 23 分钟**；编排者轮询到 `Get-Process LogonUI` 存在后再发那一发，全程他不用回来点任何东西。11-14 试过 `rundll32` 自锁，`LogonUI` 前后都是 0（静默失败）⇒ **锁屏必须是他亲手按的**。

### 撤销口令

- 「这单作废」＝ 删 `docs/evidence/s1/owner-3min-window-20261115.md`（本件是只读产物，删了不动任何票面）。
- 「本次采集作废」＝ 删 `D:\tmp\wisp\20261115\realonce*\` 那一个目录（**只删本次**，别碰 `preflight-probe-001` 那发预跑）。
- 「314 不要了」＝ 在票 314 尾追加一行「AC#2 的安静桌面那一发不做」，⛔ 不删票面原句。
