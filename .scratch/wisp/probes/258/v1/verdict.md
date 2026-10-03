# 票 258 · 验收腿 `258-v1` 判决书（非实现者，逐格攻）

验收腿：`258-v1`；起手锚 `dffd9456`（HEAD，dev）；实现件锚 `5e8748b3`（编排者代提收尾 commit）。
起手时刻 2026-10-03 21:19 +0800（自取 `date`）；骨架 commit `de82ec76`（第 3 轮内落盘）。
本文件只建不删；读数逐字；判语归格；**勾归编排者**。

**角色声明**：本腿是 D22 双角色的"另一个 agent"。实现者两任（`258-r1` 死于 150 轮帽、产码零 commit；收尾代提＝编排者本人 A544 预授权）。本腿不碰 AC 勾框、不改产码一字（探针 test 文件是本腿自己的新文件，非产码；跑完保留在树上待编排者处置——见 §8-2）、只 commit 证据件（显式 pathspec）、不 push。

---

## §0 起手锚与基线红名集合

- 起手 HEAD：`dffd9456d25ac11559f9f476b5afa848a9b0526b`；实现件锚 `5e8748b3`（2026-10-03 21:07:08 +0800）。
- **基线整包**（本腿自跑，21:19–21:29，档 `.scratch/wisp/probes/258/v1/baseline-gotest-v1.txt`，rc=1）：
  - `cmd/wisp`：红 1 枚＝`TestTicket223HandEditedFsLooseningCostsAnL2Card`（config_reload_223_test.go:312/315 "the console did not render the reload card"）。
    **solo 复跑**（21:34，`-run TestTicket223HandEditedFsLooseningCostsAnL2Card`）＝ **PASS (2.33s)** ⇒ 整包并发挤压间歇形，与派单点名的 `TestAlwaysBranch...`（C18 窗挤压 41s）**同族不同名**（死腿交件自报"C18 形"，实际这一发倒下的是 223——同为挤压间歇，归属判法相同，但名册要对上：本腿基线红名是 223 不是 C18）。非 258 回归（判据只涉 D36 reload 卡渲染，与热键链零交集）。
  - `internal/ball`：红 1 枚＝`TestC21TableColourRowsMatchTokensCSS`（design/assets/tokens.css 不存在）——派单点名的既有形；树内根因可见：`git status` 显示 `design/assets/*` 整批 `D`（未提交删除）。与本票无关，本腿不修（派单豁免）。
  - **7 枚 Test258\* 在基线整包全绿**；本腿 2 枚 probe（§2）亦绿。
- SKIP 读数：`grep -i skip` 仅命中日志文本；无 `--- SKIP` 行；无 SKIP 被当通过。

### §0-α 跟踪状态重大事实（先于一切判定）

`git ls-tree 5e8748b3 -- cmd/wisp/` 里 258 相关只有 **2 枚**：`resident_hotkey_258_seam_windows.go`（21 行缝）＋`resident_hotkey_258_test.go`（1 枚 AST walk）。
**6 枚主判据测试全部 untracked（`??`）**：`resident_hotkey_258_windows_test.go`（6 枚 Test258*）＋`resident_hotkey_live_258_windows_test.go`（4 枚 TestLive258*，winlive tag）。`git log --all -- <两文件>`＝空（从未进过任何 commit）。
⇒ 派单口径"实现件……＋7 枚 Test258*"与 git 事实不符：**实现件锚里只有 1 枚 Test258\***。go test 按目录编译、不看 git 状态，所以一切测试读数仍真实有效，但"判据已交付"在 git 里今天不成立 ⇒ §7 R1。

---

## §1 恒真两问

### 问① 摘掉桥的 rebind 跳 ⇒ 哪几枚红（-overlay 突变台，树零改动，md5 前后复量一致）

四枚突变（每枚 hunk 数＝1，票 250 定式；读数档 `mut{0,1,2,3}-run.txt`）：

| 突变 | 内容 | 默认档 `-run Test258` 结果 |
|---|---|---|
| **M1** | `resident_ball_windows.go` 桥装配块 `if hotReload != nil` → `if false &&` | **红 2 枚：`Test258BridgeRebindsLiveKeysFromConfigEdit`（0.31s 处 'R' never arrived）＋`Test258OccupiedCombinationNamesTheNewValue`（1.08s）**；其余 5 枚绿 |
| **M0** | `resident_windows.go` 装配根第四参 `hotReload258` → `nil` | **7 枚全绿（rc=0）** |
| **M2** | seam `hotkeyBridgeCheck258` 的 `check258()` 调用 → `_ = rb.hotkeyBridge`（no-op） | 7 枚全绿；`BridgeRebinds` 从 0.31s 变 1.33s（靠桥自身 1s ticker 兜底，seam 被旁路仍过） |
| **M3** | `hotCfg258` 闭包 `if err != nil || c == nil` → `if true ||`（恒回空 ⇒ 恒退 defaults） | **7 枚全绿（rc=0）** |

**读法**：
- 默认档对 **M1（桥装配块被禁用）有牙**：两枚桥判据红。恒真问①的直接答案＝**摘掉桥的 rebind 跳，红的是 `Test258BridgeRebindsLiveKeysFromConfigEdit` 与 `Test258OccupiedCombinationNamesTheNewValue` 两枚**（M1 形）；若摘法是"装配根第四参置 nil"（M0 形），**默认档零红**——因为 6 枚主判据全部自带闭包注入（`startResidentBall(observe.NewRegistry(), nil, src, src)`），装配根传什么它们根本不看；唯一看装配根的那枚是 AST walk `Test258AssemblyRootWiresTheChainAndTheBridge`，而它只断言"`len(call.Args) >= 4`"（`nil` 也是第 4 参）＋"resident_windows.go 里有 `config.LoadFile` 字样"——**`nil` 形与恒空闭包形都挡不住**。M0/M3 全绿＝实测，不是推断。
- M2 全绿＋时长变化（0.31s→1.33s）＝seam 只是加速器，真正执行 rebind 的是桥自己的 ticker；seam 被旁路判据仍有牙（M1 下 seam no-op 等效于桥死，两枚照红——M1 同时把 `hotkeyBridge` 置 nil，两形叠加）。
- **shipped 档（winlive）对 M3 有牙**：overlay 不穿透 `buildWispForTest` 的子 `go build`（第一次 M3-winlive 实验 rc=0 全绿，**作废**——子 exe 编的是 pristine 树，见 §8-3）。改用**真树突变**（backup→mutate→run→restore，md5 `04e50ea3…` 前后一致，`git diff` 空）：M3 树突变下 `TestLive258ResidentBootTakesConfigHotkeys` **FAIL**（resident_hotkey_live_258_windows_test.go:95 "a boot with planted [hotkey] values did not name config as the source"，4.13s，真红非 SKIP）⇒ **shipped 档对装配根的"恒退 defaults"腐坏有牙**。M0（nil）形未在树突变下复跑（时间预算），其 shipped 档牙性**量不到**，具名进 §8。

### 问② 反形自查（把"改配置生效"换成"改配置不生效"会不会也绿）

**不会**。反向形已由突变台实证：判据要求的正是"新值上车"这一事件本身——M1/M0/M3 三种"改了不生效"的腐坏形里，判据要么红（M1 两枚、M3 winlive boot 正控），要么其"绿"依赖的机制与正向判据不同物（M0 默认档绿是因为判据**根本不读装配根**，是判据的覆盖洞，不是"反向也绿"）。具体反证：
- `Test258BridgeMutationNoSrcKeepsOldBinding` 是**显式的反向对照**（rebind 跳被拆的死腿形，断言"新值永远不上车、旧值活下来"），它在 M1（桥死）下绿、在正控（桥活）下也绿——它钉的是"没有桥就没有 rebind"这条物理，不是"改了生效"；若把它的语义倒过来用（当作"改配置不生效"的判据），它会在真实好件上红（正控形下新值上车、VK='Q' 断言失败）⇒ 反形不绿。
- 正控探针（§2 本腿自写）在 M0/M3 形下必然红：其断言目标是 `RegisteredHotkeys()[hkSummon].VK == 0x37`（Win32 实持），装配根断供 ⇒ 桥 src 回空 ⇒ `RebindHotkeys(空集)` ⇒ summon 槽消失 ⇒ VK=0x0 ≠ 0x37 ⇒ 红。M1 突变下的 `BridgeRebinds` 失败消息（"summon VK = 0x0 … want 'R'"）就是这一物理的实测样张。

---

## §2 AC#1 正控实测（自写一发，非转述）

**探针文件**：`cmd/wisp/resident_hotkey_v1probe_test.go`（本腿自写新文件；形状同 `Test258BridgeRebindsLiveKeysFromConfigEdit` 但判据数值全新：boot 种 `Ctrl+Alt+Z`、编辑成 `Ctrl+Alt+7`（VK 0x37）、win32 实持断言、verdict 两句断言。⛔ 不引用实现者读数）。

**第一次跑（21:44）＝真红，红因本身是一枚有效读数**：
```
resident_hotkey_v1probe_test.go:87: 258-v1 PROBE RED: the [hotkey] edit never reached Win32: summon VK = 0x0 (map[2:{Mods:16387 VK:77} 4:{Mods:16387 VK:80}]), want 'W'
--- FAIL: Test258V1ProbeSummonEditRebindsLiveBall (1.47s)
```
诊断（同发日志逐字）：`level=ERROR msg="hotkey occupied by another program, not registered" hotkey=summon binding=Ctrl+Alt+W err="Hot key is already registered."` ⇒ **rebind 发生了、新值被 Win32 拒了，因为 Ctrl+Alt+W 在本机被第三方程序持有**——这正是 `DefaultHotkeys()` 注释里 R10 裁定的同一枚键（"a sweep of 84 candidate combinations … found Ctrl+Alt+W occupied by a third-party app"）。**桥的占用语义工作正常**（HotkeyTaken＋Problems 行＋不静默），是探针选键不巧。**这枚红不是产码红**，改键复跑。

**第二次跑（21:45 后，`Ctrl+Alt+7`）＝两枚全绿**（档 `v1probe-run2.txt`，rc=0）：
```
--- PASS: Test258V1ProbeSummonEditRebindsLiveBall (0.29s)
--- PASS: Test258V1ProbeConstructionMissingFileNamesTheFallback (0.01s)
```
**读数逐字（本腿自己的）**：boot 后 `RegisteredHotkeys()[hkSummon258].VK == 'Z'`（种下的值，非编译字面）；盘上把 summon 改成 `Ctrl+Alt+7` 后，经 seam 驱动桥的 Check hop，`RegisteredHotkeys()[hkSummon258].VK == 0x37` 上车；`HotkeyReport().Live()[hkSummon].VK == 0x37`；`ConfiguredHotkeys().Summon == "Ctrl+Alt+7"`；verdict 含 `hotkeys from config` 与 `hotkeys live 3/4`。**⇒ AC#1 的"改 `[hotkey]` 里 summon ⇒ 下一次建球注册的就是新值"在本进程全长产接线（`startResidentBall` 真实形状＋真 config.toml＋真球窗）下实测成立。**

**合并终态**（档 `family-run-final.txt`，21:47，rc=0）：9 枚全绿＝7 枚 Test258*＋2 枚 v1 probe。

## §3 缺省退回＋provenance 句读数

三档 provenance 的实读：
1. **nil view**（`hotCfg == nil`）⇒ `DefaultHotkeys()`＋`hotkeyProvenanceNone`（="no host config view"）——`Test258ConstructionChainMissingFileFallsBackAndSaysIt` 第一分支钉住。
2. **missing file**（LoadFile 报错 ⇒ 闭包回空视图）⇒ `DefaultHotkeys()`＋`hotkeyProvenanceDefaults`（="defaults"）——同测第二分支＋本腿 probe `Test258V1ProbeConstructionMissingFileNamesTheFallback` 第一分支独立复测（绿）。**verdict 句形状**：shipped 档 `TestLive258ResidentBootNamesProvenanceAndDefaults`（winlive，真机 3.96s PASS）实读 verdict 逐字片段：
   `hotkeys from defaults, hotkeys live 3/4: summon=Ctrl+Alt+Q live, mute=Ctrl+Alt+M live, cancel=Esc not bound while idle (`
   ⇒ **"说得出这句话"在 shipped 二进制 stdout 上实测成立**（"hotkeys from defaults"六个字在句里）。
3. **`[hotkey]` 节缺失但文件存在**（canonical config：`LoadFile` 经 `readConfigFile` 的 `cfg := NewDefaults()` 先填 schema tag `cancel default:"Esc"` ⇒ 视图 `Cancel="Esc"` 非全空）⇒ chain 走 `ApplyHotkeyDefaults` 补空三枚 ⇒ 终值＝`DefaultHotkeys()` 但 **provenance＝"config"**（"the FILE is what produced the set"）。
   ⚠ **措辞裁定点**：票面 AC#1 原话"配置文件缺失／那一节缺失时退回 `DefaultHotkeys()` 并说得出这句话"。逐字对照：**终值满足**（节缺失 ⇒ 终值确是 DefaultHotkeys 表，因 ApplyHotkeyDefaults 补的就是它），**但档位词印的是 "config" 不是 "defaults"**——即"说得出这句话"被实现解释为"退回值是默认表＋空槽补齐有日志句"（resident_windows.go 的 `empty_slots_note` slog 逐字："empty slots are filled from the product defaults by ball.ApplyHotkeyDefaults"），而非"verdict 档位词印 defaults"。**这一句的判读归编排者**：若以"退回 DefaultHotkeys 且不静默"为准 ⇒ 成立；若以"节缺失也必须印 defaults 档位词"为准 ⇒ 带条件。本腿不代裁。

**链路证据**（静态核对，全带 file:line）：`configFileName = "config.toml"`（secret.go:63）；`LoadFile`→`readConfigFile`→`NewDefaults()`（loader.go:117）→`decodeStrict`；`ApplyHotkeyDefaults`（hotkey_windows.go:87-102）逐字段补空；`DefaultHotkeys()`（hotkey_windows.go:68）＝Summon `Ctrl+Alt+Q`/Mute `Ctrl+Alt+M`/Cancel `Esc`/Panel `Ctrl+Alt+P`；`Hotkeys: cfg`（resident_ball_windows.go:269，255 收据行号引用 `:269` 实测指向该行，`TestTicket255RosterEvidenceLinesStillSayWhatTheyClaim` 等 14 枚 255 族全绿）；hot 档应用表 `{"hotkey", …}`（manager.go:281）仍在；`pRegisterHotKey` 产码调用点仍唯一（hotkey_windows.go:391）。

## §4 AC#3 越界核对（`git show 5e8748b3 --stat` 全名单逐一过）

21 文件（已在前文 §0 引 stat 原文）。**禁区核对（file 级，`--name-only` 精确匹配）**：
- `docs/PLAN.md`／`docs/specs/**`：**0 文件**。
- `internal/observe/thresholds.go`／golden／`tools/d22scan/allowlist.txt`：**0 文件**（commit 里 tools/d22scan 的三文件是 main.go/selftest.go/selftestsamples.go——票 212 的 ban #9 判据与名册，非 allowlist）。
- `frontend/**`／`design/**`：**0 文件**。
- 三枚冻结件（tokens_fourway_test / l2_grant_boundary_test / ticket90_persist_test）：**0 文件**。
⇒ `git show 5e8748b3 --name-only | grep -E '^(docs/|frontend/|design/|internal/observe/thresholds|tools/d22scan/allowlist|…三冻结件)'`＝空（rc=1）。

**非票 258 文件的性质**：commit 是 212＋258 合体。212 部分（d22scan 三文件＋11 枚注释行）本腿按"非实现者"只核**越界维度**（路径禁区＋结构），其内容判据归 212-v1；列名：models.go、resident_approval_246_windows_test.go、resident_approval_live_246_windows_test.go（三枚仅是 `startResidentBall` 新四参签名的**机械跟改**：`nil, nil` 两参插入，判据字未动）、pending_read.go、queue.go、tools.go、statevisual.go、pathresolver.go、provenance.go、bridge.go（**全部是注释行更正**，212 交件自述"10 枚真②改注释为诚实形"与本腿 diff 复读一致：4 处注释＋0 代码语义变更——`git diff` 逐行复读全部 `//` 行）。
**`agentRuntime` 结构**：`type agentRuntime`（run.go:266）在 commit 中 **0 触碰**；internal/agent/ 三文件仅 4 行注释。⇒ **没动**。
**Esc borrow（票 245 在册）**：`ball_windows.go` +19 行全部是两枚 Debug squat seam（`DebugRegisterHotkeySquat`/`DebugUnregisterHotkeySquat`，:916-930）；`escTakenOver`/`takeEscWith`/`RebindHotkeys` 的 borrow 语义 **0 行改动**；`RebindHotkeys` 注释仍逐字保留"将丢弃进行中的 Esc 借用……该残留已在票 245 中记录，而非掩盖" ⇒ **没有被顺手修**（符合 A538 边界"落地腿不许顺手做"）。
**门禁**：`scripts/d22scan.sh` 全仓 clean（bans #1-8 全过，本腿 21:5x 自跑，档 d22scan-full.txt；cmd/=97 含新增 test 文件）；`tools/d22scan` 包测试 ok 16.3s。
⇒ **AC#3 越界核对：零触碰，成立**。

## §5 AC#2 诚实格判定

**死腿自报**："AC#2 的 live 重绑只有非真窗读数覆盖"——**这句话半真半假，按读数更正**：
- 默认档 `Test258BridgeRebindsLiveKeysFromConfigEdit` 是**真窗读数**：它走 `startResidentBall` 全产码装配、`rb.b != nil` 才继续、断言 `RegisteredHotkeys()`（= Win32 实持注册表）——本机实测它 PASS 且非 SKIP（0.31s），即本机的默认档读数**有真窗**。死腿说"非真窗"不对；它没有的是 **shipped 二进制**这一层。
- winlive 层（真 shipped exe）4 枚本腿真机实测（21:48，档 winlive-run.txt）：
  - `TestLive258ResidentBootNamesProvenanceAndDefaults` **PASS**（3.96s）；
  - `TestLive258ResidentBootTakesConfigHotkeys` **PASS**（3.33s）；
  - `TestLive258ResidentRebindsWithoutRestart` **FAIL**（15.59s）——**尺子 bug，不是产码 bug**：失败消息贴出的 sink tail 原文里**就有**重绑行 `{"msg":"ball: hotkeys rebound after config change","summon":"Ctrl+Alt+R",…,"live":3}`（21:48:52.272，编辑后约 1s）；但判据 grep `summon=Ctrl+Alt+R`（**控制台等号格式**），sink 是 **JSON 冒号格式**（`"summon":"Ctrl+Alt+R"`）⇒ `strings.Contains` 永不命中。**产线行为（重绑发生、新值上车）实际成立**，读数在失败消息自己的附件里。同形红可由 M1 突变旁证其有牙（M1 下默认档同名判据红）。
  - `TestLive258ResidentOccupiedCombinationNamesTheNewValue` **FAIL**（3.41s）——**前提不成立，非产码 bug**：verdict 逐字 `panel=Ctrl+Alt+U live` ⇒ Ctrl+Alt+U 本次跑时本机**没有被占**，panel 注册成功、`Problems()` 无行可说、判据要求的 occupied 行自然缺席。该键的"machine's own risk"指的是 internal/ball 的 live 族（hotkey_live_test.go:228-238）里"**如果**被占则 SKIP-LOUD"的同款约定——那族自己也承认这键可能空着；占用前提在实现文件里没有先手验证（没有 squat 动作，纯赌环境）⇒ 判据形状对环境是脆的（winlive 同 tag 的 246 族"故意容忍无桌面"的先例在此未复用）。**占用语义本身已由默认档 `Test258OccupiedCombinationNamesTheNewValue`（进程内 squat，自搭占用，PASS）与 probe 的 Ctrl+Alt+W 意外红（真 Win32 拒绝＋Problems 行逐字出现）双向实测。**

**判定（对照票面 AC#2 原话"种一发改 summon ⇒ 下一次建球注册的就是新值"）**：
- "建球"是否必须真窗：**不必**。判据问的是"注册的就是新值"＝Win32 注册表实持；默认档 `startResidentBall` 形与本腿 probe 形都建了真球窗、断言了真注册表。真窗要求真正不可替代的一层是 **shipped 进程的装配根**（`hotCfg258`/`hotReload258` 闭包读 `rt.Layout.DataDir`）——这一层 winlive boot 正控（PASS）＋本腿 M3 树突变 winlive 红（有牙实测）已覆盖；**rebind 那一格的 shipped 读数被尺子 bug 挡住**（上面 3），但重绑行的证据已在失败附件里、且默认档同名判据＋probe 双绿。
- **判语：AC#2 成立但带条件**——条件＝(i) winlive rebind 判据的 JSON/控制台格式失配要修（一行：匹配 `"summon":"Ctrl+Alt+R"` 或两格式都吃）；(ii) 修后重跑 winlive 该枚拿一枚干净 PASS（或编排者接受"失败附件里的重绑行"作为读数并记录尺子缺陷）。**非真窗读数不足**这个死腿自报的理由**不成立**；真正的缺口是 shipped 档 rebind 的**读数仪器**，不是覆盖本身。

## §6 AC 格判语（勾归编排者）

| 格 | 判语 | 一句话依据 |
|---|---|---|
| **AC#1** | **成立但带条件** | 正控（本腿自写 probe＋winlive boot 正控＋缺省三档）全实测绿；条件＝§3 措辞裁定点（节缺失时档位词印 "config" 而 AC#1 句子字面是"退回 DefaultHotkeys 并说得出这句话"——终值满足、"这句话"的落点要编排者裁）＋§7 R1（判据文件未 commit，"交付"在 git 里不成立）＋§8-1（装配根接线在默认档零 runtime 判据、shipped 档有牙但 M0-nil 形的 shipped 牙性量不到）。 |
| **AC#2** | **成立但带条件** | 正控实测绿（默认档真窗＋probe＋winlive 附件内重绑行证据）；条件＝§5 (i)(ii)（winlive rebind 尺子 JSON 失配修复＋复跑；占用 live 判据的环境脆性登记）。死腿"非真窗"自报**判错**，本腿以读数更正。 |
| **AC#3** | **成立** | 禁区 file 级零触碰；agentRuntime 零结构改动；Esc borrow 零顺手修（+19 行全是 Debug squat seam）；d22scan 全仓 clean；255 族 14 枚全绿（收据行号同步 `:269` 实测指向真行）。 |

## §7 推翻清单（票面 §7、A539/A576/A578、派单逐句待验）

- **R1（推翻派单句"实现件……＋cmd/wisp/resident_hotkey_258_test.go（7 枚 Test258*）"）**：实现件锚 `5e8748b3` 里只有 **1 枚** Test258*（AST walk）＋21 行缝；6 枚主判据＋4 枚 winlive 全部 untracked（`git ls-tree`＋`git log --all` 双证，§0-α）。**A578 的"258 残局全绿"在树上真、在 git 里半**——收尾 commit 漏了两个判据文件没 add（A532 式 pathspec 事故的镜像：这次不是 add 多而是 add 少）。处置建议归编排者：以显式 pathspec 补 commit 两文件（含本腿 probe 文件的去留一并裁）。
- **R2（修正死腿红名名册）**：死腿交件自报整包红＝"C18 形（TestAlwaysBranch... 41s）"；本腿基线实测红＝**`TestTicket223HandEditedFsLooseningCostsAnL2Card`**（solo PASS 2.33s＝挤压间歇）。同族不同名；两枚 258 红名（M1 下）与死腿自报的"Test258 全 PASS"一致。
- **R3（推翻死腿自报"AC#2 的 live 重绑只有非真窗读数覆盖"）**：默认档 `Test258BridgeRebindsLiveKeysFromConfigEdit` 就是真窗读数（真球窗＋Win32 注册表断言，非 SKIP，PASS）。死腿把"非 shipped"误述成"非真窗"。真实缺口在 shipped 档的读数仪器（§5）。
- **R4（推翻派单句"7 枚 Test258* 里哪几枚红"的隐含前提"恒真问只有一形"）**：摘桥有两形——M1（装配块禁用）默认档红 2 枚；M0（装配根 nil）默认档**零红**。恒真问①的答案必须两形分报（§1）。
- **R5（票 258 票面 AC#1 句子的"两条入口"早在 A538 已自我修正为"run 腿空集"）**：本腿复核确认——`wisp run` 产码零 `ball.` 引用、不建球；票面句子与实现的偏差已由票面 §7 第 1 条自己登记过，非新推翻，列此存照。
- **R6（A578 引用的门禁读数 stale 一格）**：A578 写"bans #8 … cmd/=96"；本腿复跑 cmd/=**97**（多出的 1 = 258 判据文件未 commit 时按目录扫入的 `resident_hotkey_258_windows_test.go`）。无违规、纯数字漂移，存照。
- **R7（winlive occupied 判据的键选择）**：`Ctrl+Alt+U` 的"machine's own risk"出处是 internal/ball live 族同键＋SKIP-LOUD 约定；实现文件直接赌它被占、无 squat 兜底 ⇒ 本机实测未占 ⇒ 判据前提落空。这不是越界，是**判据形状缺陷**（与 §5 条件 (iii) 同源），归实现文件作者（编排者）裁改不改。
- **R8（未推翻）**：票面 §7（A538）选形 A 的四样（文件、理由、边界、撤销口令）逐句核过：装配位在 resident_windows.go（装配根）、桥机在 resident_ball_windows.go（balldebug 全套同形）、零 agentRuntime 改动、rebind 走 uiRun post-and-wait（ball_windows.go:741-752 uiRun＋:813 RebindHotkeys b.uiRun）——**A538 无一句被推翻**。A539/A576 与 258 相关部分（253-r1 占测试面、258 补位）与本腿读数不冲突。

## §8 判不动／量不到

- **§8-1 M0（装配根第四参 nil）在 shipped 档的牙性**：未在树突变下复跑（时间与测试面预算；M3 树突变一发已花一次真机 winlive）。默认档 M0 零红是实测；shipped 档 M0 是否红**量不到**。推断（不作为读数）：winlive boot 正控断言 `summon=Ctrl+Alt+J` 上 summary，M0-nil 不改 hotCfg258 ⇒ boot 正控大概率仍绿 ⇒ M0-nil 在 shipped 档可能同样无牙——**但桥武装句（`bridge armed`）在 nil 下不印**，`TestLive258ResidentBootNamesProvenanceAndDefaults:61-66` 断言该句 ⇒ 那枚大概率红。留给下一程，本腿不推断成结论。
- **§8-2 本腿 probe 文件的处置**：`cmd/wisp/resident_hotkey_v1probe_test.go` 是本腿自写的正控探针（非产码、非判据挪用），跑完留在树上。它让"9 枚"变成"9 枚"的账面（7＋2）；编排者补 commit 判据两文件时须连带裁它的去留（留下＝以后每发整包多 0.3s＋多一枚真窗依赖测试；删去＝本读数仍由本判决书存照）。
- **§8-3 overlay 不穿透 `buildWispForTest`**：winlive 族经 `go build` 子进程编译 shipped exe，`-overlay` 只影响父 go 进程的包图 ⇒ **winlive 判据的突变台必须用真树突变＋md5 还原**（本腿 M3 第二发即此法）。这个坑本腿已踩过一发（mut3-winlive-run.txt rc=0 全绿＝无效实验，档存照），记入防止复踩。
- **§8-4 物理按键验证**：`Ctrl+Alt+7` 真按下去是否触发 summon 手势（OnSummonHotkey）不在本腿射程（internal/ball live 族的 inject 真输入管那半）；本腿只断言 Win32 注册表与报告/verdict 句。
- **§8-5 审计面**：`bridge.Check()` 的重绑 Problems 行进 slog sink（hotkey_reload.go:98-100 实读）与 `HotkeyReport().Problems()`——票面 258-a1 §2 第③问的"失败句链条完整"本腿只复核了 Problems()/verdict/summary 三处读数（默认档 Occupied 测试＋probe 的 W 键意外红），**重绑行的 operator 可见性在 shipped 档的日志文件里有没有落行**量到一半（winlive rebind 失败附件里有 JSON 重绑行＝落行成立；但该读数取自 FAIL 消息，未做独立的干净 PASS 发）。

---

**终态 status**：判语三格全落（§6），无"待填"；证据件＝本判决书＋9 份读数档（baseline/v1probe-run/v1probe-run2/family-run-final/mut0-3-run/mut3-winlive-tree-run/winlive-run/d22scan-full/pristine-md5）。**未 commit 的盘上件**：本判决书增量、读数档 9 枚、probe test 文件、mut 目录与 overlay json（临时件只建不删）。commit 由本腿下一步以显式 pathspec 落（只 commit 证据件；probe 文件去留已列 §8-2 归编排者裁）。
