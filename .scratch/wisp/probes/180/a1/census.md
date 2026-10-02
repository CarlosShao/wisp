# 180-a1（10-02 第二程）—— 配置面全名册现量：哑键普查 · `unwiredKeys` 覆盖率 · 生效级别三档 · "此项不生效"的出口 · `[panel] width` 断链点

派单：只读取证腿 `180-a1`，10-02 09:5x 起手。**产码零字节改动**；票面 AC 框一枚未碰；票面正文未改（只在末尾**追加**一节）。

---

## §0 起手锚（逐字）

- `date -Iseconds` ⇒ `2026-10-02T09:57:47+08:00`
- `git log -1` ⇒ commit `00e7efeb0ca22e1ce3cbe62763d0a3088dfe4330`（"ledger(A521)：167-a1 结档…"，Author CarlosShao，Date Fri Oct 2 09:57:39 2026 +0800）
- `git rev-parse --abbrev-ref HEAD` ⇒ `dev`
- `git status --porcelain internal/config internal/panel cmd/wisp` ⇒ **空输出**（三目录起手干净）
- 全树 `git status --porcelain`：**不干净**——别人的脏件与未推临时件大量在场（`design/**` 一批 D/M、`.scratch/**` 一批 ??、`scripts/testdata/portable-tests/go` M）。本程**不还原、不提交**它们的任何东西。
- 本程读到的**所有行号均取自 HEAD `00e7efe`**（除注明）。票面锚是 `039efb47`、上一程 `180-a1` 锚是 `38fc7c0e`——**两处行号都已过期**，差集见 §8。

## §1 全名册：叶子键 → 读者 → 非测试／仅测试／无

（表头先行；枚数在 §1 正文，尺与射程逐条附。）

| 段 | 叶子键（TOML 路径） | Go 字段 | 生产读者 file:line | 判 |
|---|---|---|---|---|
| — | — | — | — | — |

## §2 `unwiredKeys` 名册覆盖率

## §3 "生效级别"三档的现量

## §4 "此项今天不生效"有没有出口 ＋ 最近可扩展点

## §5 `[panel] width` 那一环断在哪

## §6 我可能判错的条目

> 每条给"为什么会错"与"错了的后果"。这节先于表写满，是因为表填完把这节留空＝不算交件。

1. **"叶子键＝150 枚"是我定的口径，不是盘上唯一的口径。**
   尺：`internal/config/schema.go` 的 struct 走查（递归展开 `toml` 标签，`toml:"-"` 除外），
   动态子表（`llm.providers.<id>`、其下 `models.<id>`、`plugins.Entries.<id>`）的键**按模板算一枚**，不按实例算。
   ⇒ 静态 115 ＋ 动态模板 35 ＝ 150。
   **如果错了**：换口径（例如"按一份真实 `config.toml` 里出现的行数"）分母会变，
   于是 §2 的覆盖率百分比会变；**但 D 档名单的成员一枚不变**——判定用的是逐枚读者检索，不依赖分母。
   最可能的分歧是 `schema_version`：我把它算作一枚叶子键（根表里的裸键，不属任何 section），
   有人可能不算。它判 R（`loader.go`/`migrate.go` 消费），挪出去不改任何结论。

2. **"18 个 section"成立，但我用了两把独立的尺才敢这么说。**
   甲：`Config` struct 的 section 字段枚数＝18（`schema.go:110-133`，`Plugins` 虽然带 `toml:"-"` 但由
   `parse.go` 的尾块吸收器落盘，文件里确实有 `[plugins]`）。
   乙：票 198 的首建真件 `.scratch/wisp/probes/198/r1/firstrun-defaults.dump.txt`
   有 33 个表头 ＝ 18 个顶层段 ＋ 15 个子表。两把尺独立给出 18。
   **如果错了**：例如有人把 `position`/`roles`/`wake_word` 这类子表也叫"section"，那 18 会变 33；
   后果只是叙述口径，不影响逐枚判定。

3. **"零读者"的尺是选择器链正则，不是数据流分析。**
   覆盖面我补了三把反查才敢下判：
   ① 反射按名取值的口子：`FieldByName|NumField()|reflect.TypeOf` 在 `cmd internal tools` 的非测试码里
   **零命中**（只有 `internal/config` 自己用反射写默认值）；
   ② 按点分路径字符串寻址：全路径字面量在 `internal/config` 之外只有
   `cmd/wisp/config_reload.go:300`（重启档三枚键名）与两枚同名巧合
   （`internal/statemachine/events.go:21` 的 `"hotkey.mute"` 是事件名、`internal/memory/schema.go:149`
   的 `"schema_version"` 是 SQLite meta 键）；
   ③ 逐枚键名的非测试提及枚数（`cmd internal tools`，排除 `internal/config`）。
   **残余风险**：值经**接口方法**返回后再无人接手（例：`Loop.Budgets()` 返回派生量而非原始键值）；
   以及我按字段名匹配、**同名异主**的命中被我判成"非读者"时判错方向的风险——
   这类我逐枚读过命中行（见 §1 表里具名的 file:line），但只读了命中行本身，没读调用链全文。
   **如果错了**：某一枚我判 D 的键其实有人读 ⇒ 那一格从"哑键"降为"已接"，§2 的哑键总数减一枚。

4. **`[voice]` 21 枚全判 D，是本报告最重的一格，也是我最怕判错的一格。**
   尺：`config.Voice|cfg.Voice|VoiceSection` 在 `internal/config` 之外的非测试码里**零命中**；
   段级令牌 `.Voice` 同样零命中。
   但 `internal/audio/gate.go:56` 的注释逐字写着 `(config AudioSection.MicMutedDefault) here at boot wiring`
   ——**注释指名了一条本该存在的接线**，而按本仓纪律注释指名的东西一律当待验断言。
   **如果错了**（语音链路其实经另一条形如 `audio.Boot(wav, muted)` 的通道取值）：
   `[voice]`/`[audio]` 共 25 枚从 D 改判，本报告的哑键总数从 72 掉到 47 量级，
   而"owner 能配的模型参数里六成套了锁"这句话会弱化成"只有语音那一族没接"。
   ⇒ 这一格的终判该由**语音现量腿（票 228/241 系）**复核，不是这条只读腿能钉死的。

5. **`[hotkey]` 4 枚我判 S（仅旁路调试程序读者），不是 R 也不是 D。**
   依据：唯一生产读者是 `cmd/balldebug/main.go:238` 的 `mgr.Config().Hotkey`；
   出货路径 `cmd/wisp/resident_ball_windows.go:171` 逐字写 `Hotkeys: ball.DefaultHotkeys()`，
   而 `internal/ball/hotkey_windows.go:69` 逐字返回 `"Ctrl+Alt+Q"/"Ctrl+Alt+M"/"Esc"/"Ctrl+Alt+P"` 四枚字面量。
   `mgr.OnReload = bridge.OnReload()` 也**只在 balldebug:244** 出现。
   **如果错了**（balldebug 被算作一个真宿主）：这 4 枚改判 R，但 `wisp run` 里改 `[hotkey]` 仍然没反应——
   结论"改了就无效"不变，变的只是"有没有人读到过它"。

6. **`price.{in,out,cached,audio_in,audio_out}` 5 枚我单独判 W（类型在场·生产零赋值），是最可能判错的一档。**
   依据：`internal/agent/loop.go:130` 有 `Price config.Price` 字段、`internal/agent/cost.go:38` 收它算账，
   但**给它赋值的非测试语句零枚**（尺：`Price:` 全根零命中 rc=1；生产里 `agent.Config{` 字面量只有
   `cmd/wisp/run.go:1010` 一枚且不写 Price）。⇒ 成本账今天按零价卡算。
   **如果错了**（另有装配点把价卡塞进去而我没认出来）：
   票 248 设置页那两枚 `model_price_in/out` 的写入就是**端到端通**的，我这条"写了没人读"的指控当场作废，
   而 owner「自己配模型」这条路上最疼的一格（价卡）其实已经接上。**这一格务必让落地腿复跑。**

7. **`[panel] width` 断点判在"建窗那一行"，尺是反向的：我先证明没有别处能改尺寸。**
   `SetBounds|MoveWindow|SetWindowPos|Resize(` 在 `cmd/wisp internal/panel` 非测试码里
   只命中 `panel_host_windows.go:45` 的**注释**（讲 `(*edge.Chromium).Resize()`，不是调用）。
   **如果错了**：存在一条我找不到的 DPI/scale 重设路径，则断点后移到"那条路径读不读 `cfg.Panel.Scale`"，
   而 §5 那句"断在建窗"要改成"断在配置没递进构造函数"——两句话的修法不同（前者改一行、后者要加参数）。

8. **行号时效。** 本报告所有行号取自起手 HEAD `00e7efe`。
   票面引的 `schema.go:527-528`、派单引的 `unwired.go:57-101` 都已漂（现量见 §8）。
   **如果错了**（我读文件与落笔之间树又前进）：§0 的终态锚与 §8 的差集节就是给人抓 this 的；
   本程没有第二把终态尺之前，任何"某行存在某字段"的引用都应带 HEAD 号读。

---

## §7 判不动的地方

> 甲＝别人补得上（给命令与期望读数）；乙＝补不上，明写不做。

**甲#1 —— "这些哑键里，哪几枚是规格要求今天就必须生效的"**
我只量了码上有没有读者，没裁"该不该有"。规格文字逐段核对（D36 表＋SPEC-03 §3 表＋SPEC-08）需要另一程。
命令：`sed -n '2726,2748p' docs/PLAN.md` 与 `sed -n '24,43p' docs/specs/SPEC-03-config-secrets-envs.md`
逐段对照本表 §1 的 D 档名单，逐枚标"规格要求生效／规格未要求／规格自己写漏"。
期望读数：`[panel]`/`[ball]`/`[hotkey]`/`[session]`/`[voice]`/`[audio]`/`[memory]`/`[privacy]`/`[cost]`/
`[observe]`/`[agent]` 三段里除 §1 已判 R 的三枚 `[agent]` 键外**全部标着 `hot`**
⇒ 若照规格字面读，D 档每一枚都是"规格要求生效但码上没接"，本票的射程会从"一枚 width"扩到几十枚。
**这一格不该由只读腿裁**（它决定要不要把票 180 拆成一族票）。

**甲#2 —— "语音/音频那 25 枚到底是不是真没人读"**
见 §6 第 4 条。命令与期望读数：
`grep -rn "MicMutedDefault\|InputDevice\|SampleRate\|HalfDuplex\|VetoWords\|Thresholds" --include=*.go cmd internal tools | grep -v _test.go | grep -v '^internal/config/'`
⇒ 本程读数：**只剩 `internal/audio/gate.go:56` 一行注释**（非代码）。
语音腿要答的是"这条注释指的 boot wiring 在哪个文件、接线票号是谁"。

**甲#3 —— "`OnReload` 在 `wisp run` 里没挂，是不是票 42（watchdog）的既定分工"**
现量：`OnReload` 的非测试赋值全仓**两枚**——`cmd/balldebug/main.go:244` 与
`internal/config/manager.go:198` 的**读取**；`cmd/wisp/config_reload.go` 只挂了 `ConfirmLocked`（:114）
和 `OnRestartPending`（:115）。⇒ reload 档在出货进程里"有引擎、没听众"。
命令：读票 42 的票面与 `internal/watchdog/doc.go`，判"reload 事件没人接"是**分工未到期**还是**漏**。
期望读数：`internal/config/manager.go:198` 那三行若在生产跑，`OnReload` 必为非 nil；
`cmd/wisp` 里零枚赋值 ⇒ 它就是 nil ⇒ 换 ASR 模型不会触发模型加载/卸载。

**甲#4 —— "面板快照该不该加一栏承载'此项不生效'"**
见 §4。这属 C17／契约面（快照键集是有钉的名册），**只读腿不裁、不许自己决定加不加**。
编排者裁；裁完若要验"加了有没有人读"，参照
`internal/panel/composer.go:40-60` 注释里那条"新段在 pump 装配根没有读者就还是常量"的判据形状。

**乙#1 —— 不做：把 92 枚无消费者键逐枚写出"它该由哪个票号接"**
票面 AC#4 只要"出名单并逐枚归口"。归口要读 `.scratch/wisp/issues/**` 全池与 SPEC 切片表，
射程远超一条只读腿，且 `issues/83-...md:155` 那张归口表（票面 09-28 记录称 18 枚全在其中）
我**没有复认**——引用它之前得自己读到那行为止，而这一程我没读到。
⇒ 本程只交"名册＋覆盖率＋断点"，归口一律标为未做，不写"以后加固"。

**乙#2 —— 不做：跑任何编译/测试类尺**
`go test ./internal/config/`、`go list -deps` 全都没跑（`198-v1` 正在整包跑 `cmd/wisp`）。
⇒ 本报告的枚数全部出自文本尺与 Python 走查，**没有一把是编译器给的**。
若某枚判定最终依赖"这个字段确实被链接进来"，本程给不出这个级别的凭据。

**乙#3 —— 不做：`frontend/**` 与 `design/**` 两层的任何核对**
禁令覆盖。⇒ 若面板侧另有配置编辑界面（哪怕它已经渲染 `panel.width` 的输入框），
本程看不见、也不引它行号；"此项不生效没有出口"这句的射程因此**只到 Go 侧为止**。

## §8 我推翻票面与派单哪一句
