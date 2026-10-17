# 216-v2 判语件——票 216 `AC#5` 两半的终裁形状 ＋ `AC#3` 没被打那半的可达性与那一发的攻法

腿 `216-v2`｜角色＝**非实现者裁决者（只读）**｜⛔ 产码 ⛔ 测试件 ⛔ 翻框 ⛔ 改判据文字 ⛔ `go test/build/vet/list/env`（一枚都没跑）
★起手锚逐字＝`git rev-parse --short HEAD` ⇒ **`feab45e3`**（分支 dev）。⚠ 具名：本腿台件运行那一刻活 HEAD 已前进到 `6197dceb`（别腿在落 commit）⇒ 本件所有读数一律取**具名提交树**，⛔ 取工作树在飞字节。
复现锚＝`44f51fc8`（＝`216-a2` 自己的锚）。两枚锚上下面每条读数**逐枚同数**处会具名写"两枚锚同"。
尺名册＝`.scratch/wisp/probes/216/v2/logs/replay-216v2.sh`（157 行）；原始读数＝同目录 `replay-216v2.log`（534 行／49001 字节，逐把带 `pat=`/`rng=`/`rc=`，stderr 并流、⛔ `2>/dev/null` 吞）。
"生产"＝排除 `*_test.go`。引用记号 `F*/M*/C*/D*/G*` 都指那枚台件里的尺名。

## 1 四件判语（只用 成立／不成立／不适用／不可判）

**① `AC#5`ⓐ（⛔ 按英文错误串里有没有某个数字去分支）＝成立**
尺：`F1`＝`216-a2` 的 `A5` **字面逐字重放**，射程 `internal cmd tools`（⛔`_test`）⇒ **rc=1、0 命中、stderr 空，两枚锚同**；本腿自加的宽形尺 `F3`（把 `MatchString/EqualFold/==/!=/switch` 与 `Error()` 同线都算上）⇒ **同样 rc=1**。
分母本腿重算为**并集 230**（`F 系列` 前置的 `N` 段：生产 `Error()` 站点 **123** ＋ `errors.Is/As` 站点 **107**，同一把尺同行去重后并集 **230**，两枚锚同）⇒ ⓐ 的 0 不是靠窄分母撑的。
别名链：`F5` 生产把 `err.Error()` 绑进变量/字段 **10 枚**（逐枚在场：`CredentialError`／`detail`／`lastErr`×2／`res.Message`／`res.Detail`×3／`line`／`reason`）；`F6`＝a2 的 `A8` 字面（射程只它枚举的 6 枚文件）⇒ rc=1；本腿把同族尺扩到**全生产树**（`F7`）⇒ **28 枚**命中，但**逐枚眼看**都是同名不同源的 `line`/`text`（SSE 帧、golden 夹具、gitignore 解析、`[think]` 标记），⛔ 一枚来自 `err.Error()` 的绑定。⇒ **判 成立**。⚠ 这条"成立"没有牙，见 §2；**不许把它读成"已钉"**。

**② `AC#5`ⓑ（要按错误码选一句用户看得懂的文案）＝不成立**
尺：`M1`＝`MessageKey` 的生产消费者，两枚锚都只剩 **2 行**＝`internal/observe/errors.go:160`（自家注释）＋`:163`（声明）⇒ **消费者 0 枚**；正控 `M2`/`M3`＝同尺含测试与扩到 `frontend docs` ⇒ 消费者确有，逐枚＝`internal/observe/errors_test.go:191-192`、`internal/llm/matrix_14_2_test.go:175/249` ⇒ 尺不瞎。
`M4` 生产里所有"按 Class 分支"的行 **7 枚**逐枚看过：`cmd/wisp/run.go:83` 分的是**退出码**（该函数签名 `runExitCode(...) int`，逐字在场）；`internal/audio/device.go:131`／`internal/llm/anthropic/adapter.go:211`／`internal/llm/chain.go:75` 是行为判定（是不是取消）；`rules_gateway.go:77`／`internal/tools/mode.go:99`／`internal/tools/paths.go:381` 里的 `Classify/Class` 是**风险档**、不是错误类。CLI 那句给人看的话（`run.go:210/228/246/397/420/441`）＝**写死的前缀 ＋ `%v` 原样串**，按调用点分、不按码分。⇒ **今天没有落地物，判 不成立**（⛔ 成立）。
★形状（**建议、⛔ 本腿决定；开票与翻框归编排者**）：推荐读成**未闭合**（要另立一枚接线票），两读的凭据各自具名——
- 支持"未闭合"：① ⓑ 是票面**必做**语气、且钩子 `MessageKey()` 就是为它造的；② 窄读（" sites 至少把键嵌进去"）也被树否证：`errors.go:161-162` 注释自称 "UI and logs embed the key"，而生产嵌键 **0 枚** ⇒ 注释与树相反；③ 缺的东西可指名：`internal/observe/errors.go:163`（钩子）＋ `cmd/wisp/run.go`（CLI 那句）＋ `internal/panel/pump.go`／`cmd/wisp/panel_models_320.go`／`internal/panel/config_handlers.go:641`（面板/JSON 那几处原样透出）。
- 支持"不适用"：这整层**已登记推迟**＝`docs/specs/SPEC-12-roadmap-governance.md:77` 逐字 `| DEFERRED | 完整错误文案体系 | D23 | 所有用户可见错误有人话文案+建议（含 D37 全 17 类） | S6 | 部分错误暴露技术细节 |`（尺 `G3`），另 `PLAN.md:1528` 同行、`PLAN.md:182` 把"人话化错误文案体系"列进 D23 **不做**清单（尺 `G4`＝`错误文案` 在 `docs/PLAN.md docs/specs` 命中 **9** 枚）。⇒ 若走"不适用"，本票 AC#5ⓑ 要读成那行的 S6 子项，⛔ 在本票翻勾。
- ⚠ 顺带一条**不是新账**的复认：`G1`＝生产里 `DEFERRED(` 标记 **26 枚**，`G2`＝`internal/observe/errors.go` 里**一枚都没有**（rc=1，它只在注释里用英文写 "deferred (D23 … S6)"），而 `docs/evidence/s1/225-deferred-registry-audit.md:158/163` **早已**把 §5 第 77 行记为"无标记／代码就在手边却没标"⇒ 归**票 225** 的射程，⛔ 为本腿这条读数二次开票。

**③ 票面那句 `ChatInput.tsx:2205` ＝不适用（该形属外部仓 openchamber，⛔ 作本仓判据）**
尺：`C1` 提交树名单里找 `chatinput` ⇒ **两枚锚都 rc=1**；正控 `C3` 同尺 `chat-composer` ⇒ `frontend/src/components/ai-native/chat-composer.tsx` 在场。★新现量（本腿补，派单与 a2 都没写）：票面那句"未逐字复核"**在本仓就能闭合**——`docs/reports/survey-2026-09-28-composer-plus-menu.md:378-381`（尺 `C7`）逐字抄了那形：`ChatInput.tsx:2205-2206` ＋ ``if (normalized.includes('payload too large') || normalized.includes('413') || normalized.includes('entity too large')) { toast.error(...) }``。⇒ 不必去外部仓；但**抄件不是原件**，且它是**别人家的代码**，所以判语仍是不适用。旁证 `C8`＝本仓 `212-a3` 那把假引用尺早先已把这枚件名归为 `true-absent`（3 行在场）；`C9`＝`d22scan` 的 `phantom-citation` 禁令射程＝`internal/`＋`cmd/` 生产 Go 注释（⛔ 管 `docs/**`），所以那枚件名出现在调研件里不算违规。
界面侧本腿也只读提交树复量（`C6`）：`413` **rc=1**、按文案子串选句族 **rc=1**、正控 `[.]includes(` **32 行** ⇒ "尺会抓、本仓没有"。★⛔ 改票面原文（本仓规矩：只在票尾与台账记现量，冻结文字一字不动；本腿已按此在票 216 末尾追加一条，`numstat` 只加不减）。

**④ `AC#3` 后半（已勾格，本格只问可达性）＝可达性 0 复现成立；那一发的攻法见 §4**
`D5`＝`attachment` 一词进 `internal/risk internal/perm internal/tools internal/agent`（⛔`_test`）⇒ **rc=1、0 个文件**；正控 `D6`＝同尺 `params` ⇒ `assessor.go:8`／`provenance.go:23`／`rules_gateway.go:1` 三枚在场；`D7`＝`internal/perm` 目录**非空**（`store.go`＋3 枚测试）⇒ 这个 0 不是空转。机制逐字复核**在我锚上仍对得上**（`D8`，两枚锚同）：`internal/panel/attachments.go:69` `// DisplayName is what the user sees; it is never used to build a path.` ＋ `:260` `name := "attachment-" + shortHash(data) + "." + accepted.ext`。`D4`＝生产里对 name/base 做 `允许|已允许|approved|allowed` 字面判定 ⇒ **rc=1**；正控同尺去字面 ⇒ `blacklist.go:240-252` **7 枚**在场。⇒ **判：那一格"⛔ 因名字里有'允许'而改变判定"今天没有被任何生产路径违反，也没有被任何用例攻过。**

## 2 必答："它被违反长什么样"＋现有尺看不看得见那一形（⛔ 恒绿尺）

形状矩阵（`fired=0`＝尺**抓到**、`fired=1`＝尺**看不见**；逐字读数在 `F2`/`F2b`/`F12`）：

| 违规形（会写在哪） | `A5`(a2) | `A5W`(本腿宽形) | `F11`(本腿字段形) |
|---|---|---|---|
| `strings.Contains(err.Error(), "413")` | **抓到** | 抓到 | 看不见 |
| `strings.Contains(strings.ToLower(oe.Error()), "too long")` | **抓到** | 抓到 | 看不见 |
| `switch err.Error() {` | **抓到** | 抓到 | 看不见 |
| `if msg := err.Error(); strings.Contains(msg, "413")`（同线别名） | **看不见** | **看不见** | **看不见** |
| `regexp.MustCompile("4[0-9][0-9]").MatchString(err.Error())` | **看不见** | 抓到 | 看不见 |
| `strings.Contains(e.Message, "413")`／`strings.ToLower(msg)` 嵌套别名 | **看不见** | **看不见** | 抓到（字段那半） |

**点名"怎么写就会红"**（今天这几处是最自然的落点，逐枚具名）：① `internal/panel/git.go:579` 现在是 `"git 这一维读取失败：" + err.Error()`，若有人改成 `if strings.Contains(err.Error(), "not a git repository")` 换一句 ⇒ **`A5` 当场红**（矩阵第一行已实证 fired=0）；② `cmd/wisp/run.go:420` 那句 `配置未就绪（Unconfigured）：%v` 若加一发 `strings.Contains(strings.ToLower(err.Error()), "413")` ⇒ **`A5` 红**；③ `internal/llm/probe.go:229-254` 已有 `res.Detail = err.Error()`（别名绑定），若后来者在下游写 `strings.Contains(res.Detail, "413")` ⇒ **`A5`/`A5W` 都看不见**，只有本腿新加的 `F11` 看得见（fired=0）；④ `internal/panel/config_handlers.go:687` 的 `line = err.Error()` ⇒ a2 的 `A8` 恰好覆盖（该文件在它的枚举里），**换一枚新文件就不覆盖**。
⇒ **结论（按派单 ② 的口径写）**：ⓐ 那一半**主形有尺、且实测今天 0**，所以判语写"成立"；但**别名形与正则形属"缺世界"**——现有尺族（含本腿补的两把）**没有任何一枚能看全"绑进变量再判"这一族**，要真钉住只能做**定义-使用（def-use）级**检查，而仓里唯一具备该形状的仪器是 `tools/d22scan`（`go/ast` 基，`main.go` 逐字在场 `ast`/`pos.Line`；现有 **9 枚禁令**＝`bare-goroutine`/`pathresolver-bypass`/`plaintext-key`/`wallclock-timeout`/`mirror-hash`/`unparseable`/`phantom-citation`/`emoji` 等，`F10`＝CI 禁令/脚本里**没有一个**管错误串分支，rc=1）。⇒ **不许把 ⓐ 读成"有牙"**：`AC#5`ⓐ 缺的那枚钉（禁令本体）今天不存在，本腿⛔ 造它（写面不属本腿）。

## 3 补 `216-a2` 自报的两处不足

**① 分母 130 是两把尺相加、⛔ 跨路径去重**——**本腿补了**：并集尺 `N`（同一把、逐行去重、射程 `internal cmd tools`、⛔`_test`）＝**230**（`Error()` **123** ＋ `errors.Is/As` **107**，两枚锚同）。⇒ a2 的 130 里那枚 23 是**窄射程**（`internal/panel cmd/wisp`）的读数、⛔ 全仓的 `Error()` 站点数，所以 130 的问题**不是它担心的那枚**："交叠 3 枚"不会让分母虚高（`Error()` 与 `errors.Is/As` 没有同一行同时出现，并集恰等于相加），真正的问题是**少算了 100 枚** `Error()` 站点。ⓐ 的结论不受影响（`F1` 一直跑在宽射程上）。
**② 派单起点 `internal/cmdutil` 本仓不存在**——**判 不适用**（⛔ 补路径）：`git ls-tree feab45e3:internal` ＝23 枚子包（`agent audio ball buildinfo config llm memory models observe panel perm plugin proc projctx risk secret session speech statemachine streamkey tools watchdog winsec`），**无 `cmdutil`**；`git ls-tree feab45e3:cmd` ＝ `balldebug llmrecord wisp`。两枚锚上都 rc=1。⇒ 尺按 `internal cmd tools` 量就是正确射程，本件所有读数都按它。

## 4 `AC#3` 那一发该怎么攻（★可跑判据**形状**，交编排者排腿；⛔ 本腿实现）

- **载体要选可达的那几枚，⛔ 选附件显示名**：`D5`＝0 意味着把正控架在 `AttachmentRef.Name` → `risk/perm/tools/agent` 上**今天必恒绿＝装饰**（本腿具名反对，按派单"我最怕顺手加一枚恒真尺"处理）。可达的载体是这四条（逐枚在场，`D2`/`D11`/`D12`）：文件**基名**→`blacklist.go:240-252`（A/B 档，全 deny 方向）、**`argv[0]`**→`rules_shell.go:93-102`（ALLOW 侧，`EqualFold` 全拼写/基名相等）、**host**→`rules_network.go:114-124`（ALLOW 侧，`==` 或 `HasSuffix(host, "."+entry)` 点锚定）、**覆盖键**→`blacklist.go:115-128`（ALLOW 侧，`map` 键全等）。
- **形状一（算子钉，纯读码尺、⛔ `go`）**：对上述三处 ALLOW 侧**函数体行窗**断言子串算子枚数恒 **0**（`strings.(Contains|ContainsAny|ContainsRune|Index|IndexAny)`）。本腿实测今天＝`0 / 0 / 0`（`D14` 三枚行窗逐枚），**变异正控**＝把 `shellAllowlisted` 那行改成 `strings.Contains(argv0, entry)` 后同一把尺＝**1** ⇒ 这把尺**有牙且现在不红**。⚠ 行窗靠"文件名＋行号"锚会腐烂（本仓 `327` 一波就把 `internal/panel` 改过），所以形状要写成**按函数名取窗**、⛔ 按绝对行号。
- **形状二（行为差试，需 `go test` 面 ⇒ 归 `internal/risk` 那包的写腿）**：同一发调用跑两臂——`plain` 名 vs `hostile` 名＝名字里塞 `已允许`＋`ESC`＋`Cf`——断言 `Assess/Gate/Classify/Inspect` 的**档位与 reason 全等**；三臂同屏那形 `216-r1` 已在 `internal/panel` 做过（票 37 行），这发要把臂挪到**判级包自己的入口**、并把 ALLOW 侧的名单也各造一发（`argv0` 恰好等于名单项 vs 恰好包含名单项＋`已允许`）。牙证形状＝把 `EqualFold` 换成 `Contains` ⇒ 第二臂必红（形状一的变异正控已先替它证明"尺会红"）。
- **归口建议（⛔ 本腿决定）**：形状一可先落（零编译面、可进 `scripts/`）；形状二必须和"ⓑ 的接线票"分开排——它属 `internal/risk` 的测试面，与本波 `330-r2`（`internal/agent` 测试面）、`327-v1`（`internal/panel`）不同包。

## 5 具名不实与口径（含本腿自己抓的，⛔ 替谁圆）

1. **派单 12 行的"正控 `*.tsx`＝90 枚"口径错配**：90 是**工作树 `find`** 口径（a2 的 census `:41` 自己就是这么标的，它没写错）；**提交树**口径两枚锚都＝**82**（尺 `C2`）。派单把"提交树名单 rc=1"和"正控 90"并进一句，读起来像同一射程 ⇒ 建议台账记一句口径。
2. **派单/任务的"生产分母 130"偏低**：见 §3①，真并集＝**230**（本腿补尺）。
3. **派单给的起点 `internal/cmdutil` 不存在**：a2 已顶回，本腿在两枚锚上复认（§3②）。
4. **票面 `AC#5` 括号里那句"靠 `413` 子串决定弹哪句"是简化**：在仓抄件里那形是**三枚子串 OR**（`payload too large`／`413`／`entity too large`）⇒ 方向不差、枚数差一枚。⛔ 改票面，只在票尾记现量。
5. **产码注释不实（不是票面不实）**：`internal/observe/errors.go:161-162` 自称 "UI and logs embed the key"，生产嵌键 **0 枚**（`M1`）⇒ 归 ⓑ 那半的射程，改它要写腿。
6. **`a2` 甲①/甲②的数本腿全部复现一致**：`413`＝4 枚（2 注释＋`stream.go:165` 赋整型＋`internal/llm/errors.go:204 case 413:`）、`too long`＝1 注释、`exit status`＝1（`tools/d22scan/runtests.sh:103` 注释）、`D1`＝37、正控＝37/13 枚族——⛔ 一处不一致。仅一处口径差：测试面"带数字错误串断言"本腿尺读出 **15**（a2 报 16），两把尺字面不同（本腿的记在 `F9`），⛔ 相减。
7. **本腿自抓**：`D13` 第一版把一枚 `a|b` 形尺漏了 `-E` ⇒ BRE 下竖线是字面量、造出一枚**假 0**（`rc=1`）。已在台件里把"假 0"与"`-E` 正解（**24 枚在场**：`composer.go:145` 调用点＋`display_scrub_216.go` 7 枚＋测试 16 枚）"并列留痕——这正是派单 ① 防的那形，本腿差点自己踩进去。
8. **锚会动**：起手 `feab45e3`、台件跑时活 HEAD `6197dceb` ⇒ 全程读具名树，读数不受影响（记下来防"我这读数怎么对不上"）。

## 6 交件自证（本腿自己重放的尺，逐字）

```
$ bash .scratch/wisp/probes/216/v2/logs/replay-216v2.sh   # 全量读数见 replay-216v2.log（534 行）
=== V0 锚（⚠ …）===            6197dceb   T1=44f51fc8 T2=feab45e3
### F1（A5 逐字，两枚锚）       rc=1        rc=1        stderr 空     ← ⓐ 今天 0 命中
### F2 变异正控 [A5]             fired=0/0/0（主形抓到）·fired=1/1/1（别名·正则·字段＝看不见）
### F4 正控 internal/risk        count=37                                            ← 同族尺会出数
### M1 MessageKey 生产（两枚锚） 2 行＝errors.go:160 注释＋:163 声明 → 消费者 0
### M2 正控（含测试）            errors_test.go:191-192 · matrix_14_2_test.go:175/249
### C1 chatinput（提交树两枚锚） rc=1 · rc=1     C3 正控 chat-composer.tsx 在场
### C2 正控 *.tsx 枚数           提交树 82 / 82 · 工作树 find 90
### D5 attachment→判级包         rc=1（0 文件）  D6 正控 params→3 枚文件 · perm 目录 4 枚件
### D4 名字里"允许"当判据        rc=1            正控 blacklist.go:240-252 7 枚在场
### D14 ALLOW 侧三处行窗子串算子 0 / 0 / 0       变异窗＝1                            ← 形状一有牙
### G2 errors.go 里的 DEFERRED 码 rc=1           G3 SPEC-12:77 登记表那行在场
```
四件具名判语（重复一遍便于编排者翻框）：**ⓐ 成立（没牙，别名/正则形＝缺世界）**｜**ⓑ 不成立（建议按"未闭合"排接线票，走"不适用"须与 `SPEC-12 §5`:77／票 225 归口对上）**｜**票面 `ChatInput.tsx:2205` 不适用（外部仓形；本仓 `survey :378-381` 已有逐字抄件）**｜**`AC#3` 后半：可达性 0 复现成立，攻法两枚形状见 §4，恒真尺那枚（架在附件显示名上）本腿具名反对**。
本腿动过的字节＝只 §写面四处；`internal/**`／`cmd/**`／`frontend/**`／`design/**`／`PLAN.md`／`docs/specs/**` 零字节；⛔ 翻任何框、⛔ 改票名/`Status:`/`-done`、⛔ push。
