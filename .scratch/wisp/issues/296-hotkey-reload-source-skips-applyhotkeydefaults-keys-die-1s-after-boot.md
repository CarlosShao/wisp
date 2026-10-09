# 票 296 — 热加载那一半**没过 `ApplyHotkeyDefaults`**：开机 1 秒后 `live=3` 变 `live=0`，静音/召唤/面板三枚键在本机全死（真机现量，非推断）

**立票**：2026-10-09 17:5x 编排者（来路＝**票 290 `AC#3` 第②形那次真机窗口**：机主在场、我按他的托盘操作取读数，跑的就是他机器上那份 `config.toml`。⇒ 这不是"读码发现的形状"，是**真机上打印出来的两行自相矛盾的日志**）
**性质**：★**回归**，而且**是我自己落的那一笔**。引入点＝`5e8748b3`（2026-10-03，258-r1 由编排者代笔收尾，台账 A544 预授权那批）。票 258 选形 A 之前，常驻腿**没有** bridge ⇒ 开机绑上的三枚键会一直活着；之后 bridge 每 1 秒把"没过默认值合并"的裸配置读回来的四格**重新绑一遍**，空串＝"配置里关掉了"，于是三枚键在开机后约 1 秒被自己拆掉。
**后果（对用户怎么说）**：**契约许诺的那枚"一键静音"热键（`docs/PLAN.md:1312` 的 `Muted` 行、`SPEC-08`）在他这台机器上按下去没有任何反应**；同理 `Ctrl+Alt+Q`（召唤）与 `Ctrl+Alt+P`（面板）。今天唯一还能用的静音入口是**托盘那一项**（那条是票 290 `290-r1` 刚接上的，本次真机已证能开门关门）。

## 现量（每条都带尺；⛔ 引用前先重跑，行号是快照）

- **机主的配置真相**（尺＝`sed -n '/^\[hotkey\]/,/^\[/p' C:/Users/swq/AppData/Roaming/wisp-dev/config.toml`，17:5x 现跑）：`[hotkey]` 节逐字四行＝`summon = ''`／`mute = ''`／`cancel = 'Esc'`／`panel = ''`。⇒ **三格是空串**，而这正是"没配"的样子（不是他主动关的；`cancel` 那格有值所以它活下来了）。
- **★两半的来源闭包，一枚过默认值、一枚没过**（同一枚装配根，逐字）：
  - 构造那一半 `cmd/wisp/resident_windows.go:181 hotCfg258 := func() ball.HotkeyConfig {` … `:195 return ball.HotkeyConfig{Summon: h.Summon, ...}` ⇒ 裸映射，**但它的消费者补了默认值**：`cmd/wisp/resident_ball_windows.go:235` 逐字 `	cfg := ball.ApplyHotkeyDefaults(hotCfg())`（在 `residentBallHotkeyChain258` 里）。
  - 热加载那一半 `cmd/wisp/resident_windows.go:205` 逐字 `	hotReload258 := func() ball.HotkeyConfig {` … `:214` 逐字 `		return ball.HotkeyConfig{Summon: c.Hotkey.Summon, Mute: c.Hotkey.Mute, Cancel: c.Hotkey.Cancel, Panel: c.Hotkey.Panel}` ⇒ **裸映射，而它的消费者不补**：`cmd/wisp/resident_ball_windows.go:353` 逐字 `		bridge := ball.NewHotkeyReloader(b, b.ConfiguredHotkeys(), hotReload)`——第三个实参直接就是那枚裸闭包，**链路里没有任何 `ApplyHotkeyDefaults`**（尺＝`grep -rn "ApplyHotkeyDefaults" cmd/wisp/*.go | grep -v _test` ⇒ 命中只有 `resident_ball_windows.go:235` 与 `resident_windows.go:193` 那句注释，**`hotReload258` 体内 0 命中**）。
- **契约写在被调那一方的文档注释里，是明文要求的**（⇒ 本票不是"我偏好另一种写法"）：`internal/ball/hotkey_reload.go:48-49` 逐字 `// HotkeySource returns the currently effective bindings (the host's mapping` / `// of [hotkey], already defaulted via ApplyHotkeyDefaults).`；同文件 `:16-19` 那段用法示例逐字含 `return ball.ApplyHotkeyDefaults(ball.HotkeyConfig{`。⇒ **`cmd/wisp` 这个宿主违反了它自己引用的前置条件**。
- **空串在绑定的那一端怎么读，也是明文**：`internal/ball/hotkey_windows.go:494-496` 逐字 `		case b.Binding == "":` / `			b.Status = HotkeyDisabled // disabled by config: never attempted` / `			slog.Info("hotkey disabled (unset in config)", "hotkey", b.Name)`。⇒ 桥拿到裸四格后**照规矩**判定"配置里关了"，它没错；错在喂给它的人。
- **`ApplyHotkeyDefaults` 自己的注释已经把这一形点名预言了**（`internal/ball/hotkey_windows.go:78-85`，逐字起头 `// ApplyHotkeyDefaults fills the empty fields of a config-sourced binding set` / `// with the product defaults. [hotkey] in config.toml distinguishes nothing` / `// between "unset" and "explicitly disabled" - both are an empty string - so a` / `// host that maps config.Hotkey straight in would silently end up with NO` / `// summon key at all on a fresh install`）：函数逐字 `:93 func ApplyHotkeyDefaults(cfg HotkeyConfig) HotkeyConfig {` 起，`:95 if cfg.Summon == "" {`／`:98 if cfg.Mute == "" {` 两支填默认。⇒ **"straight in"那一支就是本票的 `hotReload258`**。
- ★**真机读数（本机、本次窗口，逐字来自 `.scratch/wisp/probes/orch/2026-10-09-c1c2-resident.raw.md`）**：
  - `:18` `time=2026-10-09T17:27:04.936+08:00 ... msg="ball: the resident leg created the floating ball window" hotkeys_provenance=defaults hotkeys_live=3 hotkeys="summon=Ctrl+Alt+Q live, mute=Ctrl+Alt+M live, cancel=Esc not bound while idle (cancel is borrowed only during Confirming), panel=Ctrl+Alt+P live"`
  - `:27`/`:28`/`:30` `time=2026-10-09T17:27:05.939+08:00 level=INFO msg="hotkey disabled (unset in config)" hotkey=summon`／`hotkey=mute`／`hotkey=panel`（**比上一行晚 1.003 秒**）
  - `:31` `time=2026-10-09T17:27:05.939+08:00 level=INFO msg="ball: hotkeys rebound after config change" summon="" mute="" cancel=Esc panel="" live=0`
  ⇒ `live=3 → live=0`、**无人改过配置文件**、且 `summon="" mute="" panel=""` 逐字复现了上面那三行空串。这一发是**确定性**的（桥的 tick＝`resident_ball_windows.go:355 bridge.Run(ctx, time.Second)`）。
- **为什么 258 的两枚在程测试全是绿的（本格是"缺哪个世界"，不是"哪枚断言写错了"）**：`cmd/wisp/resident_hotkey_258_windows_test.go:266-281` 那枚 AC#2 正控 `Test258BridgeRebindsLiveKeysFromConfigEdit` 写的夹具是**四格全显式赋值**（逐字 `summon = "Ctrl+Alt+Z"` / `mute = "Ctrl+Alt+M"` / `cancel = "Esc"` / `panel = "Ctrl+Alt+P"`，改后的那一发也只换 `summon = "Ctrl+Alt+R"`），并且它把**同一枚闭包同时当构造源与热加载源**（`:283 rb := startResidentBall(observe.NewRegistry(), nil, src, src)`）⇒ 桥读回来的永远是"四格都有值"。**"文件里三格是空串"这个形状，在测试过的世界里从来没出现过**（尺＝`grep -rn 'mute = ..' --include=*_test.go cmd/wisp` ⇒ 没有任何一枚桥用例写空串）。同文件 `:153` 那格反而把"源返回全空"钉成了期望（`if c := loadHotkey258(t, dir); c != (ball.HotkeyConfig{}) {`）⇒ **源本身全空是对的**（默认值不属于它），**错的是没有任何一枚用例把"源全空 ⇒ 合并后仍须 live=3"这条组合关系跑一遍**。
- ⛔ **不许把这一格读成"内部包有洞"**：`internal/ball` 的两端（`ApplyHotkeyDefaults` 与 `registerAllWith`）各自行为都正确、且互相一致；洞只在 `cmd/wisp` 那**一枚闭包的组合**上。尺＝上面 `grep -rn "ApplyHotkeyDefaults" cmd/wisp/*.go` 的命中分布。
- ⚠ **本票与票 290 的关系要说清**：票 290 `290-r1` 把**手势回调**（`OnMuteHotkey` / `OnTrayMute`）接到了门上，那一接本身是对的；本次真机两形读数（开门/关门各一次）都是**走托盘**取的。本票坏的是**热键那一支能不能被按到**——它今天按不到，所以票 290 `AC#3` 若要把"热键形"也记成已读，**必须按在本票落地之后**（⛔ 不许拿"代码里 `OnMuteHotkey` 已接上"充当"用户按得动"）。

## 要建什么（⛔ 先只读把射程量完，再落地）

- [ ] **AC#0 射程普查（只读，⛔ 不许只答"没找到"）**：把"宿主把 `[hotkey]` 裸映射进 `ball`"的**所有**落点数出来，逐枚判它有没有过 `ApplyHotkeyDefaults`。名册至少含：`cmd/wisp/resident_windows.go:181/:205` 两枚闭包、`cmd/balldebug/main.go:237-241`（尺＝现读那三行：balldebug 的 src **有** `ApplyHotkeyDefaults` ⇒ **同一个坑 balldebug 没有、只有 wisp 有**，本格必须具名确认这一枚反例）、`internal/ball/hotkey_reload.go:13-20` 那段示例注释、以及 `cmd/wisp` 里任何别处 `ball.HotkeyConfig{` 复合字面量（尺＝`grep -rn "ball.HotkeyConfig{" --include=*.go cmd internal | grep -v _test`，⚠ 名册要**排除** `.scratch/wisp/probes/**` 下的变异拷贝并在件里写明排除了几枚）。判据＝每枚给"读回来的四格是否可能为空串／谁补默认"两行＋尺。
- [ ] **AC#1 缺的那个世界先进夹具（⛔ 这一格先于产码改动）**：新增一枚在程用例，夹具**逐字复现机主那四行**（`summon=''`/`mute=''`/`cancel='Esc'`/`panel=''`），跑**真装配**（`startResidentBall(reg, nil, hotCfg, hotReload)` 用**生产那两枚闭包的同形**，⛔ 不许 `src, src` 复用同一枚——那正是逃掉这个 bug 的形状），断言**两件事**：ⓐ 桥第一次 `Check()` 之后 `report.Live()` 仍含 summon/mute/panel 三枚（`live=3`），ⓑ 且 `ConfiguredHotkeys()` 读回的是合并后的值而非空串。完成判据＝**改产码前这枚用例必须红**，红句要能指出 `live=0`（⛔ 不许"改前也绿"；改前也绿＝它验的不是这件事）。
- [ ] **AC#2 落地那一发（修法二选一，本格先写代价再裁）**：
  - **甲（改宿主，射程最小）**：`hotReload258` 体内套 `ball.ApplyHotkeyDefaults(...)`，与 `resident_ball_windows.go:235` 同形。
  - **乙（改被调方，把前置条件变成不变式）**：`NewHotkeyReloader`／`Check()` 内部对 src 的返回值过一次 `ApplyHotkeyDefaults`，并**同时**把 `ApplyHotkeyDefaults` 的"空＝没配"语义在 `internal/ball` 一侧写成文档＋用例。
  ⇒ 判据＝两形各写"要动哪几枚文件＋会不会让 balldebug 变＋`[hotkey] mute = 'none'`（用户真想关掉一枚键）这一形将来怎么表达"。⚠ **这一问是甲案的实质风险**：`hotkey_windows.go:83-84` 逐字写着 `a config with summon = "none"/"off" is the host's business` / `(it arrives as "" and re-enables the default, which is the safe direction: ` ⇒ 按现契约**"关不掉"是故意的**。本格必须把"用户能不能主动关一枚键"写成**已定案**还是**契约空白**（空白 ⇒ 上机主清单，⛔ 不许写码腿自己挑）。
- [ ] **AC#3 改完后的真机复跑（⛔ 这一格归编排者跑，派单就要写）**：同一台机器、同一份 `config.toml`（⛔ 一个字都不改它），起进程后**必须看到** `hotkeys_live=3` 且**不再出现** `live=0` 那行；随后**按一次 `Ctrl+Alt+M`**，要求 stdout 出现票 290 那两句开门/关门读数之一（`resident_ball_windows.go:536` 的 `wisp: ball %s: %s\n` 形）。跑前三把尺：`tasklist //FI "IMAGENAME eq balldebug.exe"`＝0、`tasklist //FI "IMAGENAME eq wisp.exe"`＝0、`date` 记锚。
- [ ] **AC#4 门禁＋越界**：`GOFLAGS= go build ./...` rc=0；`sh scripts/d22scan.sh` rc=0；`gofmt -l` 空；`$(go env GOPATH)/bin/gofumpt.exe -l` 空（⛔ 裸 `gofumpt` rc=127 不算"跳过"）；`PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test ./cmd/wisp/ -count=1` 改前改后各**≥2 发取交集**、逐名作差＝新增红 0 枚（⚠⛔ **改前那两发起手前先跑 `tasklist //FI "IMAGENAME eq wisp.exe"`＝0**：本仓实测过一枚活着的开发进程会把整包逐名名册搅出 ±1 枚假红，而这枚 bug 的验证恰好要起真进程 ⇒ **AC#3 那一发必须排在两发改前读数之后**，否则我把自己的读数弄脏）；每一把门禁件自己落一行 `rc=N`（⛔ 0 字节的件＝那格没交）；`git show --stat` 名册只含 `cmd/wisp/`（甲案）＋`probes/296/**`，⛔ `frontend/**`／`design/**`／三枚冻结件／golden／`thresholds.go`／`allowlist.txt`／`internal/config/schema.go`／D43 表零字节；⛔ 零 push、commit 必带显式 pathspec。

## 禁区

- ⛔ **不改 `config.toml` 的任何内容**（那是机主的配置，不是夹具）；本票要的是"代码读懂空串的方式"，不是"把盘上的空串填上值"。⛔ 尤其不许把这一格修成"启动时往配置文件写默认值"（那会把 `D36` 的三档生效级别与"用户写的文件"混成一谈）。
- ⛔ **不动 `internal/config/schema.go` 里 `[hotkey]` 的四枚默认值**（票 290 同族的规矩：默认值一字不动）。本票只碰 `cmd/wisp`（甲案）或 `cmd/wisp`＋`internal/ball`（乙案）。
- ⛔ 不许为了变绿放宽票 258 既有的任何断言（尤其 `resident_hotkey_258_windows_test.go:153` 那句"源可以返回全空"——它是对的，别动它；要加的是**组合**那一层）。
- ⛔ 不许顺手改票 290 的 `AC#3` 勾（那一格按 AC#3 之后、由非实现者验收腿判）；⛔ 不许顺手把 `hotCfg258`/`hotReload258` 合并成一枚闭包（两半的**用途**不同：构造期读一次、热加载每 tick 读一次，合并会改掉 `:188` 那条 Warn 的出现时机）。
- git：只 commit 不 push；⛔ `add -A`／`amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`；仓内不删文件、不建 worktree；**临时件只建不删**（README 规则 8——本票的证据件与任何临时日志都落 `.scratch/wisp/probes/296/**`，⛔ 不落仓根、⛔ 跑完不 `rm`）。

## 排程

`296-a1`（只读＝AC#0＋AC#2 的代价表，⛔ 零产码）→ 编排者裁甲／乙（若 AC#2 判"none 那一形是契约空白"⇒ 先摆机主一句话，含"不做"栏）→ `296-r1`（落地＝AC#1 夹具先行＋AC#2 那一发＋AC#4）→ `296-v1`（非实现者验收翻勾）→ **AC#3 由编排者跑真机**（与票 290 `AC#3` 的"热键形"补读**同一次窗口**，机主只被占用一次）。
⛔ 与票 293/294/295 **同属地界**（都碰 `cmd/wisp`，⛔ 不许混批、同包只一枚写腿在飞）；`296-r1` 排在 `292-r1` 交回并复跑之后。

**Status:** **未开工**。⛔ 零翻框、零 push。真机凭据已在盘上（`.scratch/wisp/probes/orch/2026-10-09-c1c2-resident.raw.md`，`:18/:27/:28/:30/:31` 五行），本票是那次窗口的**第一手结论**，不是待验假设。

## Progress log

（追加式；每条 `- [YYYY-MM-DD HH:MM:SS +08] agent=… did=… next=…`，钟点由 `date` 的 stdout 插值，不手打。）
- [2026-10-09 18:01:46 +08] agent=编排者 did=立票 296（真机现量：票 290 AC#3 那次窗口的第一手结论）；引入点 5e8748b3 与 balldebug 那枚反例都现跑核过 next=派 296-a1 只读（AC#0 射程＋AC#2 甲乙代价表）
