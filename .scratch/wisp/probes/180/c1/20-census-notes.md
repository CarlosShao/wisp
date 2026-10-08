# 180-c1 格 2 · AC#4 同族普查（`internal/config/schema.go` 的 `default:` 字段 × 生产读者）

起手锚 `601c2b18`；读数取于 **`70b00885`**（工作树被别的腿推进，见 `10-ac1-verdict.md` §5）。
产码零字节改动；`internal/config/schema.go` 在 `601c2b18..HEAD` 的 diff ＝ **0 行**。

---

## 1. 母体先定数（⛔ 不引任何人的转述）

三把独立尺，各落 rc：

```
grep -c 'default:' internal/config/schema.go                       => 70    rc=0
grep -n "default:" internal/config/schema.go | grep -v 'toml:'       => 1 行: schema.go:8   rc=0
python census.py: default_tags_total                                 => 70    rc=0
python census.py: lines_with_two_plus_default_tags                   => 0     rc=0
```

⇒ **70 是"含 `default:` 这个子串的行数"，其中 `schema.go:8` 是文档注释在引用该写法本身**
（`// D36 rule 3: defaults live ONLY in the `default:"..."` struct tags. The TOML`），不是字段标签。
⇒ **本程母体＝69 枚带 `default:` 标签的字段**；去重字段名 **62** 枚。
⇒ 派单的 70 与票内 09-28/10-02 两程的 69 **都对**，差的就是 `schema.go:8` 那一行。
⇒ 没有任何一行含两枚以上 `default:` 标签，所以"标签数＝行数"这一条是现量确认的，不是假设。

`Height`（`schema.go:544`）与 `Scale` 等**没有 `default:` 标签**，因此**不在 69 枚母体内** ——
`[panel] height` 今天其实有读者（`panel_resident_windows.go:207` 同一行返回），却不会出现在这张表里。
**这是母体定义的必然结果，不是漏项；要看 `[panel]` 全五枚请读 `10-ac1-verdict.md` §6。**

---

## 2. 三把尺，不是一把（这是本程与票面 AC#4 最大的方法差别）

票面 AC#4 只给了**一把**尺：`grep -rn "\.<字段名>" --include=*.go internal/ cmd/ | grep -v _test.go`。
本程照跑，但同时跑了两把更严的，因为它**答不了它被拿来答的问题**：

| 尺 | 定义 | 零命中枚数 | 这数是什么意思 |
|---|---|---|---|
| **A 名字尺**（票面原尺） | `\.Field\b`，范围 `internal/ cmd/`，去 `_test.go` | **17 / 69** | 全树连提都没提过这名字 |
| **B 去自身包** | A 再去掉 `internal/config/**` | **36 / 69** | schema 自己包外没人提这名字 |
| **C 归属尺** | B ＋**同一行还得出现所属类型名**（`Panel.Width`、`Ball.Size`…） | **55 / 69** | **没有任何一处被证实是在读这一枚字段** |

⇒ **确证有生产读者的只有 14 枚；55 枚确证不了。**
⇒ 票面那把尺给的"零读者 17 枚"是**下界中的下界**：它把 `internal/config/manager.go` 的整块拷贝、
`defaults.go` 的反射填标签、`catalog.go` 的名册字符串**全算成了读者**，还把 21 枚别的结构体的同名 `Enabled`
算成 `PanelSection.Enabled` 的读者。**同一族病复发在尺上，不在字段上。**

三档完整分解（69 ＝ 17＋19＋19＋14）：

- **17 枚 A 档零**：全树非测试码不提此名（清单见 `reader-classes.txt` CLASS A）。
- **19 枚 只在 `internal/config` 内被提**：schema 自己碰自己，外部无读者。
- **19 枚 外部有同名命中但归属未定**：见 §4，⛔ 不算"有读者"。
- **14 枚 归属确证**：见下。

### 14 枚归属确证的读者（本程唯一敢写"有读者"的一档）

```
Cancel(HotkeySection)     att=1  internal/agent/approval/approval.go:117
Max(RetryConfig)          att=1  internal/llm/retry.go:60
BackoffMS(RetryConfig)    att=1  cmd/wisp/run.go:1226
Loose(Compat)             att=4  internal/llm/anthropic/adapter.go:314
AllowMissingUsage(Compat) att=7  internal/llm/anthropic/adapter.go:265
Enabled(ModelSpec)        att=1  internal/llm/resolver.go:128
TimeoutMS(LLMSection)     att=1  cmd/wisp/run.go:1086
SteeringEnabled(Agent)    att=1  internal/agent/loop.go:283
ProjectInstructionsEnabled(Agent) att=2 cmd/wisp/run.go:1055
ConfirmTimeoutSec(Risk)   att=2  cmd/wisp/resident_approval_windows.go:516
L1WindowSec(Risk)         att=2  cmd/wisp/resident_approval_windows.go:515
DeleteEnabled(FSSection)  att=1  internal/tools/fs_write.go:596
Mode(ProxyConfig)         att=2  internal/agent/prompt.go:261
Width(PanelSection)       att=2  cmd/wisp/panel_resident_windows.go:207  ← 本票那一枚
```

---

## 3. 三处控制的逐枚读数（派单点名的反造假格）

### ① 正控（尺必须给出非 0）——4 枚有效 ＋ 1 枚无效自报

```
Width    nontest=2      Mirror  nontest=3
Enabled  nontest=21     Scale   nontest=1
RetainTurns nontest=0   ← ⛔ 无效控制：`grep -n "RetainTurns|retain_turns" internal/config/schema.go` => rc=1
                            它根本不是一枚 `default:` 字段，报 0 不是尺的失灵，是挑错了样本。具名记着。
```

`Width`/`Mirror` 的命中逐行读过：`Mirror` 的真读者是 `cmd/wisp/models.go:163 cfg.Models.Mirror`
（另两枚是 `config_readers_255.go` 里的引文）；`Width` 见 §5。⇒ **尺对"确实在读"的字段给非 0，正控通过。**

### ② 反控（尺报 0 的，逐枚单独复跑、去掉管道看原始）——16 枚

`raw` ＝ 含测试的全部命中，`nontest` ＝ 票面尺读数：

```
ClickThrough raw=1  nontest=0      InputDevice       raw=1  nontest=0
HideOnFullscreen raw=1 nontest=0   SampleRate        raw=1  nontest=0
WarmTimeoutSec raw=3 nontest=0     DiagnosticsOptIn  raw=1  nontest=0
SettlingSec raw=1 nontest=0        RetentionDays     raw=1  nontest=0
ConversationIdleSec raw=1 nontest=0 L1Enabled        raw=1  nontest=0
EchoRef raw=0 nontest=0            L1Max             raw=1  nontest=0
L3Enabled raw=0 nontest=0（⛔ 此名非本表字段，是一枚无效探针）
SizeMB raw=0 nontest=0             Days raw=0        nontest=0
SLOSampleIntervalSec raw=1 nontest=0
```

**关键结论（这是 `grep | wc -l` 的 0 和"命中了但被过滤掉"的分界）**：16 枚里 **12 枚 raw>0**，
命中全部落在 `internal/config/boundary_test.go`（`WarmTimeoutSec` 另有 `loader_test.go`）。
⇒ 这些"0 生产读者"的准确说法是 **"只有测试读者"**，不是"这名字不存在"。
⇒ 剩下 4 枚（`EchoRef` / `SizeMB` / `Days` 以及票面的 `L3Enabled` 无效探针）raw＝0，才是真·全树无提及。
⇒ **`grep -v _test.go` 确实吞掉了命中，但没有吞掉任何"生产命中"**：本程未对这 12 枚做任何"其实有读者"的追认。

### ③ 同名撞车（⛔ 不许直接算成"有读者"）——13 枚已标 `collision_flag`

最危险的四枚，用归属尺单独数了一遍：

```
Ball.Size      名字尺=17  →  "Ball\.Size"   = 0
Observe.Level  名字尺=68  →  "Observe\.Level" = 0
App.Theme      名字尺=4   →  "App\.Theme"   = 0
Proxy.Mode     名字尺=59  →  "Proxy\.Mode"  = 2（这一枚是真的）
```

`PanelSection.Enabled` 单独立尺：
```
grep -rn "Panel\.Enabled\|panel\.Enabled" --include=*.go internal/ cmd/ | grep -v _test.go  => rc=1 空输出
```
⇒ 那 21 枚 `.Enabled` **没有一枚属于 `[panel] enabled`**。
另举一例坐实撞车：`internal/ball/hotkey_windows.go:104` 的 `cfg.Panel` 是**一枚热键字符串字段**
（`cfg.Panel == ""` / `cfg.Panel = d.Panel`），和 `config.Config.Panel`（`PanelSection`）**同名不同物**。

---

## 4. 归属未定的 19 枚（本程**不**判它们"有读者"，也**不**判它们"死键"）

```
Size(BallSection)  Enabled(WakeWord)  Provider(ASRConfig)  Provider(TTSConfig)
Enabled(AECConfig) Enabled(RealtimeConfig) Enabled(VoiceSection) ConversationMode(VoiceSection)
MicMutedDefault(AudioSection) Temperature(Role) ThinkingIntensity(Role) RepeatThresholds(LoopGuard)
MaxRounds(AgentSection) TokenBudget(AgentSection) PermissionMode(RiskSection) RedactPaths(PrivacySection)
Enabled(PanelSection) VerifySignature(ModelsSection) Level(ObserveSection)
```

其中 `MaxRounds` / `TokenBudget` / `Temperature` / `ThinkingIntensity` 大概率**确有读者**
（`internal/agent/guard.go` / `budgets.go` / `prompt.go` 就在干这些事），只是取值形状是
**先把值取进本地变量或别的 struct 再传**，所以"同一行同时出现段名和字段名"这条归属尺照不到它们。
⇒ 这是**尺的射程下界**，不是字段的死刑。**这 19 枚一律标"归属未定"，交落地腿逐枚读码定案。**

---

## 5. 我这把尺的缺陷（具名，不遮掩）

1. **文本尺不是编译器**。全部 69 枚读数出自 `grep`，没有一把是类型检查给的。
   ⛔ 本程零 go 命令（硬约束 1：`cmd/wisp` 测试面归 `198-v1`／`33-r10` 在飞腿）。
   ⇒ 无法区分"读了一个字段"与"读了一个同名的 map key / JSON tag / 别的类型的字段"，除非归属尺恰好命中。
2. **归属尺（C）要求段名与字段名同行**，对"取进局部变量再传下去"的读者**系统性漏报**（见 §4 那 19 枚）。
   所以 **55 是"无读者"的上界，14 是"有读者"的下界**，真值落在中间；本程**不裁中间那 19 枚**。
3. **名字尺（A）系统性高报读者**：把 `internal/config` 自己包内的拷贝/反射/名册都算成读者，
   把别家的同名字段算成本字段读者。**票面 AC#4 若只跑 A，会得出"只有 17 枚哑键"，比真实情况乐观 3 倍。**
4. **命中里含字符串字面量**。`Width` 的 2 枚命中之一是 `cmd/wisp/config_readers_255.go:161` 的**台账句子**
   （一句引用 `cfg.Panel.Width` 的字符串），**不是读者**。真读者只有 `panel_resident_windows.go:207` 那一枚。
   ⇒ `Width att=2` 这个数**虚高一枚**；§2 表里已注明真读者行。
5. **注释也算命中**。`cmd/wisp/panel_resident_windows.go:207` 之外，
   `grep "\.Height\b"` 的命中之一来自 `panel_host_windows.go:248` 的一句注释。本程未剔除注释。
6. **母体只含带 `default:` 标签的字段**，`[panel] height`/`scale` 这类无 default 标签的字段天然不在表内（见 §1）。
7. **`Scale`/`RetainTurns`/`L3Enabled` 三枚控制探针本身挑错**（前两枚不是 default 字段，后一枚名字不存在），
   已各自具名标为无效，不参与任何结论。
8. **射程只到 Go 侧的 `internal/` ＋ `cmd/`**（票面尺的射程）。
   `frontend/**` 与 `design/**` 未读未引（硬约束 2 与票面"不改任何前端文件"）⇒
   **"零读者"这句话今天仍然只对 Go 侧成立。**
9. **行号会漂**。本程引用的 `panel_host_windows.go` 行号在 `70b00885` 复核过（该 commit 刚动过这文件）；
   `schema.go` 行号在 `601c2b18` 与 `70b00885` 上同立（diff 为 0 行）。

---

## 6. 与本票前几程的数字对账

| 来源 | 报的数 | 本程复跑 | 差在哪 |
|---|---|---|---|
| 派单 | `default:` 70 处 | **70 行**，其中 1 行是注释 ⇒ **69 枚字段** | 一致（口径差 `schema.go:8`） |
| 09-28 `180-c1` | 69 行 / 去重 61 名 / 零读者 **18** | 69 枚 / 62 名 / A 档零 **17** | 名＋1、零＋−1：那 1 枚差的是 `[panel] width` —— 它在 09-28 零读者，**今天有了**（票 255 AC#4） |
| 10-02 `180-a1` 第二程 | 69 枚 default 标签；叶子键全量 150；名册外哑键 **76** | 母体一致 | **不矛盾**：a1 按**全部叶子键 150** 数，本程按**带 `default:` 的 69** 数。两个分母。 |
| 本程（第三把尺 C） | — | 确证有读者 **14**／无读者 **55** | 比票面 AC#4 单尺口径严一档；⛔ 与 a1 的 76 不直接可比（不同母体） |
