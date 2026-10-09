# 247-a3 — 票 247 问③（降级：设备被占／无权限／无设备）单问普查件（只读腿）

## §0 锚点、射程与漂移记录

- 起手锚 `0091b1b9`（`dev`，11:17 +0800）；落笔锚 `9204f858`（`dev`，11:24 +0800，写本节时 `date` 现量）。
- **漂移记录**：写件期间 `HEAD` 从 `0091b1b9` 漂到 `9204f858`。`git diff --name-only 0091b1b9..9204f858` 落点**只有票 281／285 的 issues／probes／`pending-and-issues.md`**，**一枚都不在我引用的产码路径**（`internal/audio|ball|observe|statemachine`、`cmd/wisp`、`docs/specs/SPEC-05`）。⇒ 本件全部读数在两枚锚上**字节相同**，行号有效；下面每处 `〔尺〕` 均为 `git show HEAD:<path>` 对象层现读。
- **射程**：**只答问③**。问①（落点甲／乙）／问②（隐私闸门与默认档）／问④（`Armed` 那一格）＝ `a1`／`a2` 已答，本件**不重做**，仅在具名对照处引它们。
- 零设备访问：本腿不碰麦克风、不录音。
- 票面问③逐字复核（`sed -n '23p'` 于 `247-the-capture-stack-...md`，现读到）：
  > ③ **降级**＝设备被占／无权限／无设备时球的形状与文案说什么（SPEC-05 §3.4：分类＋可见，不许静默）。
  —— 与编排者转述一致，**无更正**。

---

## §1 三形各自"今天代码走到哪一步"（`internal/audio` 现读；⛔ 不照抄 `a1`，行号自取）

三形在 `internal/audio` 里**都有活分岔、都会返回一枚具名分类错误 `*observe.Error`（Class=`audio_device`）**；无一是"零分岔"。错误从 `Start` **同步**返回（`wasapimic_windows.go:89 return <-started`）。

| 形 | HRESULT 常量（`device.go`） | 抛出点（打开动作，`wasapi_windows.go`） | 分类映射（`DeviceError`，`device.go`） | 给出的文案（人看句／日志句） | 用例钉（`hotplug_test.go`） |
|---|---|---|---|---|---|
| **设备被占**（别程序独占） | `hrAUDCLNTDeviceInUse=0x8889000A` `device.go:77` | `Initialize(shared mode)` 失败 `wasapi_windows.go:285` | `case device.go:100-101` → `observe.New(ClassAudioDevice, detail+" (hr): "+inUseGuidance)` | `inUseGuidance` `device.go:89`（逐字"…held in exclusive mode by another application; close that application or pick another device"） | `hotplug_test.go:388` `DeviceError(hrAUDCLNTDeviceInUse, dev, "Initialize(shared mode)")` |
| **无权限**（隐私拒入） | `hrEAccessDenied=0x80070005` `device.go:72` | `Initialize(shared mode)` 失败 `wasapi_windows.go:285` | `case device.go:98-99` → `…+privacyGuidance` | `privacyGuidance` `device.go:86`（逐字给出 Windows "Settings > Privacy & security > Microphone" 路径） | `hotplug_test.go:393` `DeviceError(hrEAccessDenied, dev, "Initialize(shared mode)")` |
| **无设备**（默认捕获端点举不到） | 走枚举层，**不经 `DeviceError`**：`GetDefaultAudioEndpoint` 非 `S_OK` | 枚举返错 `mmdevice_windows.go:86-87`（逐字"GetDefaultAudioEndpoint(flow=…) failed (HRESULT)"） | `openCurrent` 用 **`observe.Wrap`** 二次包：`wasapimic_windows.go:261` `observe.Wrap(ClassAudioDevice, err, "capture device enumeration failed")` | ⚠**只有类＋HRESULT，无建议句**（见 §6-更正2） | 无专属用例（`a1` 亦如实说"没有专属用例"） |

- 附带第四形（会话中掉设备，热插拔失效，非"启动降级"）：`hrAUDCLNTDeviceInv=0x88890004` `device.go:75` → `case device.go:102-103`，文案"device invalidated (unplugged or disabled)"。列出仅供并排，不属问③三形本体。
- **`DeviceError` 的 switch 无 `NotInit` 那一枚 case**：`hrAUDCLNTNotInit=0x88890001` `device.go:74`，`switch` `device.go:97-110` 未列它 ⇒ 走 `default` `device.go:108-109`（只给码不给建议）。**但无设备这一形今天根本不进 `DeviceError`**（它进 `openCurrent` 的 `observe.Wrap`），所以这条缺 case 是"将来若改走 `DeviceError` 才会撞上"的隐患，非现路径缺陷——这点与 `a1 §③` 的推测不同，见 §6-更正2。

---

## §2 错误往上传到哪一层就断了：断在**第 0 跳**（`internal/audio` 之上无人接）

- **包内两面都保留分类**（不丢成泛化文本）：
  - 同步面：`Start` 返回的就是 `*observe.Error`（`wasapimic_windows.go:89`；`run` 在 `:178 started <- err` 把 `openCurrent` 的分类错误原样递回）。
  - 异步面：`Err()` `wasapimic_windows.go:119-122` 取回的仍是同一枚 `*observe.Error`。
  - 计数面：`Stats().LastError` `audio.go:88/172` 记其字符串；`FramesDropped` `audio.go:84` 与 `Dropped()` `audio.go:92` 是"静默丢帧禁止"的既有量具（D38d）。
  - 机制面：`observe.Wrap` 造 `&Error{Class, Detail, Err}`（`errors.go:184-186`）**保住 Class 与 cause**；`ClassOf` `errors.go:230` 能把 `ClassAudioDevice` 取回。⇒ 分类是**具名类型＋`Class` 字段**，不是被 `fmt` 拍扁的一句文本。
- **往上追就断在第 0 跳**：`internal/audio` 全仓 **非测试 importer ＝ 0**，`AudioSource` 生产消费者 ＝ 0。
  - 尺（调用形状锚，⛔ 不锚符号名）：`git grep 'github.com/CarlosShao/wisp/internal/audio' HEAD -- '*.go'` → **exit 1（零命中）**；`git grep 'AudioSource'` 去掉 `internal/audio/`、`.scratch/`、`_test.go` 后 → **exit 1**；`WASAPIMicrophone`／`mic.Start`／`.Start(ctx` 在 `internal/audio/` 之外 → **No matches**。
  - 结论：分类信息**没机会**往 `cmd/wisp`／`internal/ball` 那一跳上被保留或被丢——因为那一跳的边**不存在**（票面标题"zero importers"的现量复述）。`cmd/wisp/resident_ball_windows.go:22` 逐字自陈"this file starts no state machine of its own…no microphone"、`:350-351`／`:408`（`ballGestureWhy`）同族自陈"this process has no task pipeline and no microphone"——是**该常驻腿声明自己不拥有采集**的范围句，**不是**"麦克风不可用"的降级文案。
- **D37 归位（只报形状不判合规）**：错误分类与跨边界传播归 `PLAN.md` 的 D37；`observe.Error`＋`ErrorClass`＋`Wrap`/`ClassOf` 就是那套机械。本腿只报"包内保住分类、包外零调用边"这一形状，⛔ 不判它合不合 D37。

---

## §3 球那一侧到底有没有对象可画：状态名册＋文案名册（除 `Armed` 外**没有降级位**）

### 3.1 状态枚数与枚名（`internal/statemachine/states.go` 现读，`states.go:11-30`，逐枚）

共 **20 枚**：`FirstRun / Sleeping / Armed / Muted / Listening / Thinking / Acting / Speaking / Warm / Conversation / Confirming / AwaitingApproval / Settling / Downloading / Unconfigured / NoNetwork / Error / Queued / Stuck / WatchdogAlert`。

- **能表达"采不到音"的：零枚。**
  - `Muted` 由 `EvMuteKey`（`hotkey.mute`，`table.go:65/75/79`）驱动＝**用户主动静音**，不是设备降级，别混。
  - `NoNetwork` 由 `EvNetworkDown`（`table.go:252`，`AnyState → NoNetwork`）驱动＝**一条既有"专用失败态＋专用图标"的先例**——网络这一族失败有专属可见态；**音频设备没有它的同族兄弟**。
  - `Error`（`IconX`）是唯一泛化兜底态。
- **降级位（除 `Armed` 之外）**：**没有**。`Armed` 那一格 `a2` 已答；`a2` 之外再无采音降级位（本轮把 20 枚全枚举过一遍确认）。

### 3.2 状态机有没有合法边把"采不到音"画进某个态

- 事件存在：`EvAudioDeviceLost Event = "audio.device-lost"` `events.go:30`（D43 #14）。
- 但全表**唯一一条**用它：`table.go:97 From: StateListening, Event: EvAudioDeviceLost, To: StateError`。**`From` 只有 `Listening`**，**无 `Sleeping`／启动期边**。⇒ 启动期采集失败**今天没有合法态可进**（与 `a1 §⑦`/`⛔不能走状态机那条边` 同源，本腿 `sed` 复量 `table.go:97` 确认）。硬造一枚 `Sleeping→…` 起点＝改 D43 转移表＝契约面（本腿不裁，交编排者）。

### 3.3 图标层可画什么（`internal/ball/statevisual.go` 现读）

- `Icon` 枚名 `statevisual.go:10-21`：`None / AudioLines(Listening) / Volume(Speaking) / Mic(Conversation) / Slash(Muted) / X(Error) / WifiOff(NoNetwork) / Key(Unconfigured) / ArrowUpRight(FirstRun) / AlertTriangle(WatchdogAlert) / Refresh(Stuck)`。
- `Visual` 结构体（`statevisual.go` 内）字段只有 `SizePx / Opacity / CoreColor / CoreAlpha / GlowColor / RingWidth / Icon / Progress`——**没有文字／prose 字段**：球＝球＋图标，不画句子。
- **没有任何一枚图标表示"麦克风不可用"**（`IconMic` 归 `Conversation` 会话态，非降级态）。⇒ 即便接线，"采不到音"在球上**能借的对象只有 `Error`/`IconX`**（还得先解决 §3.2 无合法边）。

### 3.4 "降级文案"名册扫（⛔ 只搜一个词就下零命中＝本仓抓过的假话形；下面是我用过的语义变体名册）

在 `internal/audio/` **之外**的产码里，按语义变体名扫"麦克风不可用"这一类句子（尺＝`git grep -i`，逐形）：
`microphone unavailable`／`mic unavailable`／`no microphone`／`no mic`／`mic not`／`capture device`／`audio device`／`device busy`／`no input`／`mic.*denied`／`input device`。

- 命中**全部**落在两类，**没有一类是给用户看的降级提示**：
  1. `internal/audio/device.go:86/89` 那两句 guidance——其上方注释 `device.go:83-85` 逐字"Log-safe: device names only, never audio content (D16)"＝**为日志而写**，非球／面板文案。
  2. `cmd/wisp/resident_ball_windows.go:22/350-351/408` 的"no microphone"＝**该腿自陈不拥有采集**（§2 已述），不是"设备不可用"。
- ⇒ **结论**：`internal/audio` 之外，**"麦克风不可用"这一类人看文案＝零命中**。分类文案只在 `device.go` 里、只喂日志，永远够不到球。
- 诚实边界：球／面板的"文字面"若存在，应在 `frontend/**`（本腿冻结件，零读），故 §3.4 只覆盖 Go 产码侧；面板文字另见 §7-(f)。

---

## §4 SPEC-05 §3.4 要什么＋三形逐形判（⛔ 不用"应该会"填空）

### 4.1 §3.4 原文（`docs/specs/SPEC-05-agent-core.md`，`sed -n` 现读；标题行＝`:76`）

- 标题逐字（`:76`）：**「### 3.4 重试与失败语义（§14.2，失败必须可见且可区分）」**
- 要"可见＋可区分"的两句话，逐字各抄一行：
  1. 标题内：「**失败必须可见且可区分**」（`:76`）
  2. 表格"重试耗尽"行：「| 重试耗尽 | 明确报错，**不得静默降级** | `Error` + 面板给原因与建议（换 provider/查 Key/查额度） |」（`:82`）

### 4.2 §3.4 的**主语其实是 LLM provider，不是音频设备**（更正编排者的隐含前提）

`:78-83` 那张表五行逐枚是：`429/5xx/网络抖动` / `重试耗尽` / `401/403` / `额度耗尽` / `流式中途断开`——**没有一行讲音频设备**。⇒ 票面把 §3.4 挂到"设备被占／无权限／无设备"上，是**借它的"失败必须可见且可区分／不得静默降级"原则做类比**（`a1 §③` 末段已标"借来的"，本腿 `sed` 复量确认那张表的行）。真·音频降级口径的**现成硬锚**应是 `D42#2/#12` ＋ `observe.ClassAudioDevice` ＋ `device.go` 那五处 guidance，而非 §3.4。

### 4.3 按 §1–§3 现量逐形判（对"分类／可见／可区分"三轴）

| 形 | 分类（`ClassAudioDevice` 具名错误＋建议句？） | 可见（用户／球看得见？） | 可区分（用户分得出是哪一形？） | 判定 |
|---|---|---|---|---|
| **设备被占** | **合格**：`DeviceError case device.go:100-101`＋`inUseGuidance device.go:89` | **不合格**：§2 零 importer ⇒ 球静默 | 文本层分得出，呈现层丢（文案进不了球） | **今天不合格**；可见轴"今天不可能合格"（前置＝AC#1 接线） |
| **无权限** | **合格**：`case device.go:98-99`＋`privacyGuidance device.go:86` | **不合格**：同上 | 同上 | **今天不合格**；同上 |
| **无设备** | **部分**：有 `ClassAudioDevice`（`mmdevice_windows.go:86-87`→`Wrap wasapimic_windows.go:261`），**无建议句**（只有 HRESULT） | **不合格**：同上 | **不合格**：连人看句都没有 | **今天不合格**；分类轴也仅半合格 |

⇒ 三形对 §3.4"分类＋可见、不许静默"的借来标准，**今天全不合格**，且"可见"这一轴在接线（票 247 AC#1）之前**不可能合格**——不是实现得差，是**那条边还没接**（本腿据实报形，⛔ 不判该不该）。

---

## §5 给写腿一句可执行的最小形状（谁返回／谁保分类／谁画；⛔ 不写码不动产码）

- **谁返回**：`internal/audio` **已经做完**——`WASAPIMicrophone.Start` 同步返回 `*observe.Error(ClassAudioDevice)`（`wasapimic_windows.go:89`），`Err()`／`Stats().LastError` 也已持其。**返回侧零新代码**（被占／无权限两形无需新增；无设备若要补人看句，改 `device.go` 一处即可，见下）。
- **谁保分类**：`cmd/wisp` 的常驻腿（`a1` 形甲那枚进程）拿到 `mic.Start` 的返回错误后，**必须用 `observe.ClassOf` 取 `ClassAudioDevice` 往下带**，⛔ 不许 `err.Error()` 一句拍扁成泛化文本。复用现成 verdict/slog 形状（`resident_ball_windows.go:345-352` 那套"记 verdict＋`slog`＋不 `os.Exit`"），把被占／无权限／无设备**具名说出来**。
- **谁画**：球。**D43 禁路**——设备丢失进 `Error` 只有 `Listening→Error`（`table.go:97`），启动期／`Sleeping` **无合法边**。所以：
  - **形①最小（不碰状态机、零 D43 风险）**：降级只落 **日志／tray／verdict**，球态不动。→ 要动产码**≈2 枚**（`cmd/wisp` 的调用腿 ＋ `resident_ball_windows.go` 的 verdict 文案）。
  - **形②要在球上看得见、又不改状态机**：给球加一枚**非状态机的降级位**（`liquidMotion`/`Ball` 上一枚 flag ＋ `statevisual.go` 一枚覆盖／暗化图标）→ 再 +2~3 枚（`internal/ball` 内）。⚠ **一旦想做成"新增一枚 State"＝改 D43／契约面＝未定义即停**，须人工批准，本腿不替写腿拍。
  - 可选 +1：`internal/audio/device.go` 给"无设备"补一句 guidance（或把 `openCurrent` 的枚举失败从 `observe.New` 泛句改指到带建议的 `DeviceError`），让三形文案对称。
- **要不要新依赖边**：**要，且恰好一枚**——`cmd/wisp → internal/audio`（这正是票 247 AC#1 要接的那条边，今天不存在）。除此之外**不需新边**：`internal/audio` 已 import `observe`；`cmd/wisp` 已 import `observe`＋`internal/ball`；`internal/ball` 已 import `statemachine`＋`observe`。⇒ 无论形①形②，新边都只有 `cmd/wisp → internal/audio` 这一条。

---

## §6 我推翻／更正上面哪几句转述（含编排者、含 `a1`）

1. **更正编排者"③唯一没被答的一问"**：**不准确**。`a1/census.md` 已含一整节 `## ③ 降级…`（`a1:202-229`），它**答了本件问③的第①子问**（三形各返回什么、落在哪行）并**已标 SPEC-05 §3.4 主语是 LLM**。`a1` **没答**的是本件的**第②（往上传到哪断）／第③（球状态名册＋文案名册）／第④（逐形判）**。⇒ 本件与 `a1 §③` 是**互补**，非重复；编排者可并排引用，⛔ 别把 `a1 §③` 当"没有"。
2. **更正 `a1 §③` 对"无设备形"的推测**：`a1` 写"实际串形要以真机读数定、且它走 `Wrap` 而非 `DeviceError`"。本腿静态读定的具体串形＝`ClassAudioDevice` + `"capture device enumeration failed"`（`wasapimic_windows.go:261` Wrap）＋底层 `"GetDefaultAudioEndpoint(flow=…) failed (HRESULT)"`（`mmdevice_windows.go:86-87` `observe.New`）——**不经 `DeviceError`，故没有 guidance 建议句**。`a1` 担心的"`DeviceError` 缺 `NotInit` case"只在**将来若把该形改走 `DeviceError`** 才成立，**非现路径**。此条**仍非真机读数**（本腿没碰设备），但静态形状已定。
3. **无设备＝有分岔、非零分岔**：编排者留了"若某形零分岔就明写"的口子。本腿现量：**三形都有返回分类错误的活分支**，无零分岔形；只是"无设备"分岔给的是**泛化 HRESULT 文本**（分类有、建议无）。

---

## §7 这一问里仍然量不到的格子＋缺什么读数

**要真跑／要设备（本腿零设备，做不到）**
- (a) 真机被占／被拒／拔干净时，`Start` 是否**逐字**吐 `inUseGuidance`／`privacyGuidance`／枚举泛句——需一次真机采集会话。
- (b) 共享模式撞车在本机驱动上到底回 `0x8889000A`（InUse）还是 `0x88890004`（DeviceInv）——HRESULT 随 Windows／驱动变，静态读定不了。
- (c) "无设备"时 `GetDefaultAudioEndpoint` 返的是 `AUDCLNT_E_DEVICE_INVALIDATED` 还是 `NOT_INITIALIZED`——需一台零捕获端点的机器；且这决定 §1 那条"`DeviceError` 缺 `NotInit` case"是否会**真**被走到。

**要 owner／要编排者派单（非技术事实，腿不能自决）**
- (d) 写腿许不许加"非状态机的球降级位"（形②），还是只准走日志／verdict（形①）——`a1 §⑦-P5` 已把"不能走状态机那条边"停车等裁，本件不重开。
- (e) 把票面问③的 SPEC-05 §3.4 引用**改指**到 D42/D37 音频锚（而非借 LLM 重试表）＝文档面变更，需人工批准，本腿只标不改。

**本腿尺够不着的形状**
- (f) 球／面板到底**有没有文字面**可承载降级句：`internal/ball/statevisual.go` 的 `Visual` 无 prose 字段（§3.3）；面板文字在 `frontend/**`＝本腿冻结件零读。故"文案能不能画在球上"的另一半（面板侧）本腿**没量**，须由能读 `frontend` 的腿补，或据 `internal/ball` 判"球体本身只有图标没有字"。
