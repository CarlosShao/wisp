# 216-a2 只读普查——票 216 `AC#5` 今天的形 ＋ `AC#3` 没被打过的那半（可达性）

腿 `216-a2`｜角色＝只读普查｜⛔ 产码 ⛔ 测试件 ⛔ 判语 ⛔ 翻框 ⛔ 改票面判据文字 ⛔ `go test/build/vet`
锚＝**具名提交 `44f51fc8`**（`git log --format=%H -1 44f51fc8` ⇒ `44f51fc8745c3ca551f76dc2c0fe78f93ca2b86d`；⛔ `HEAD~n`）
口径＝**全部尺读提交树 `44f51fc8`**（`git grep <pat> 44f51fc8 -- <paths>`），⛔ 混在飞工作树（并发在飞 `326-v1`/`165-v2`/`327-a1`/`328-a1`/`270-r2`）。
"生产"＝排除 `*_test.go`。逐枚尺的原始 stdout/stderr＋`rc=` 见 `logs/reprobe-216a2.log`（尺本体 `logs/reprobe-216a2.sh`，可整跑）；下文把那份原始读数简写为 **`out`**，`out:5-11` 这种引用指它的**行号读数**（⛔ 当判据，只当坐标）。落名用 `.log` 而非 `.out`：`.gitignore:8` 是 `*.out`，我这枚件要能进提交树，⛔ 动 `.gitignore`（机主外部前端 agent 的在飞件）。
派单给的起点路径 `internal/cmdutil` **在本仓不存在**（`git ls-tree --name-only 44f51fc8:internal | grep -x cmdutil` ⇒ rc=1）⇒ 甲① 按派单原射程目录 `internal/panel cmd/wisp` 量，另补 `internal cmd tools`。

## 甲① 错误提示今天怎么选句：三类枚数

分母两把尺（**量的**，不是算的）：
- `A1` `Error()` 站点·含测试＝**125** 行；`A2` 同尺去 `_test`＝**23** 行（逐枚见 `logs/…out:4-26`，全部 `file:line` 已在场）。
- `A3` `errors.Is(|errors.As(` 含测试＝**324** 行；`A4` 同尺生产＝**107** 行。
⇒ **生产分母＝23＋107＝130 枚**（两把尺各有射程，不相减、不并集去重：A2 射程 `internal/panel cmd/wisp`，A4 射程 `internal cmd tools`）。

| 类 | 枚数 | 分母 | 逐枚落点 |
|---|---|---|---|
| (a) 按码/类型分支 | **107** | 130 | `A4` 全量；再加 `A9` 两枚显式 `switch`：`cmd/wisp/run.go:83`（`switch observe.ErrorClass(errClass)`→退出码）、`internal/perm/store.go:225`（`switch err`→`context.Canceled/DeadlineExceeded`）。分类侧＝`internal/llm/errors.go:94 ClassifyTransportError`（12 枚 `errors.Is/As`，无一枚读文案）＋`internal/llm/errors.go:176 ClassifyHTTPStatus`（`switch d.Status` 整型）＋`internal/llm/anthropic/stream.go:153 classify()`（`switch e.Type` 线上类型字段，`:165 status = 413` 是**赋整型**、不是比串） |
| (b) 按字符串子串/正则分支 | **0** | 130 | `A5` 尺＝`strings.(Contains|Index|HasPrefix|HasSuffix)(…Error()…|…Err…)` ＋ `switch …Error()` ＋ `Error() ==` 三形合扫，生产 `internal cmd tools` ⇒ **rc=1（0 命中）**（`out:30-31`） |
| (c) 原样把错误串透出去 | **23** | 130 | `A2` 那 23 枚逐枚（`out:4-26`）。形态＝拼前缀／塞字段／进日志／塞 JSON／截断，**没有一枚按文案选句**：如 `internal/panel/git.go:579`（`"git 这一维读取失败：" + err.Error()`）、`internal/panel/config_handlers.go:641`（`payloadError %q`）、`cmd/wisp/panel_models_320.go:111`（`boundedFirstLine320(err.Error(), …)`＝截断不是分支）、`cmd/wisp/panel_control_windows.go:744/748`（`ctlJsString(err.Error())` 递给 JS） |

**别名链也量了**（怕子串尺只认字面 `.Error()`）：`A7` 生产里把 `err.Error()` 绑进变量/字段＝**10 枚**（`out:34-43`：`CredentialError`／`detail`／`lastErr`×2／`res.Message`／`res.Detail`×3／`line`／`reason`）；`A8` 对这些绑定变量再做子串判定＝**rc=1（0 命中）**（`out:45-46`）。正控 `A6`＝同族尺（`strings.(Contains|ContainsAny|HasPrefix|HasSuffix|EqualFold|Index)`）在生产包 `internal/risk` 上＝**37 行**（`out:32`）⇒ 这把子串尺会出数，只是 0 处作用在错误文案上。

**★ 关键形（AC#5 的"必须"那半今天没有落地物）**：票面假设存在一层"按码选句"，仓里确有那层的**钩子**——`internal/observe/errors.go:163 func (c ErrorClass) MessageKey() string { return "error." + string(c) }`，注释自己写着"copy text itself is deferred (D23 error-copy work, S6)"。但 `A10`＝**MessageKey 的生产消费者 0 枚**（`out:52-54` 只剩声明＋它自己的注释）；正控 `A11`＝同尺在测试里**有**消费者（`internal/observe/errors_test.go:191-192`、`internal/llm/matrix_14_2_test.go:175/249`，`out:56-63`）⇒ 尺不是恒 0。同类先例：`.scratch/wisp/probes/198/a1/census.md:191` 已具名记过 `MappedStates` 也是"生产侧 0 枚调用者"。

## 甲② 有没有错误串带数字被下游子串判定（`413`／`too long`／`exit status`）

生产 `internal cmd tools`（⛔ `_test`）逐 token 扫，**每把都带正控**（正控＝`B2` 同尺 `observe.ClassProvider` ⇒ 13 枚文件、`out:79-91`）：

- `413`＝4 命中，逐枚看过、**无一枚是子串判定**：`cmd/wisp/notify_windows.go:126`（注释里的 Win32 码 1413）、`cmd/wisp/panel_host_windows.go:644`（注释里的行号区间 408-413）、`internal/llm/anthropic/stream.go:165`（`status = 413` 赋整型）、`internal/llm/errors.go:204`（`case 413:` 在 `switch d.Status` 整型分支里＝**这就是"按码分支"本身**）。
- `too long`＝1 命中＝`internal/tools/task.go:102` 注释。
- `exit status`＝1 命中＝`tools/d22scan/runtests.sh:103` 注释（CI 脚本，非产品路径）。
- 退出码走的**是**类型路：`B3` ⇒ `cmd/balldebug/diff_windows.go:561/564`（`st.ExitCode()`）、`cmd/wisp/slo_windows.go:879`（`ps.Exited(), ps.ExitCode()`）；`grep 'Contains(…err…,"exit'` 全仓（含测试）＝rc=1。
⇒ **今天没有任何一条生产路径靠"英文错误串里带不带某个数字"决定说什么。**

⚠ 但**数字确实随错误串出向**（这是事实、不是违规）：`observe.Error.Error()`（`internal/observe/errors.go:188`）把 `class (provider_code): detail: cause` 拼成一句，`NewHTTPError` 在无上游码时 `e.ProviderCode = strconv.Itoa(d.Status)`（`internal/llm/errors.go:241`），`httpErrorDetailString` 印 `"HTTP %d"`（`:254`）⇒ 甲①的 (c) 类 23 枚把带数字的串原样送进 stdout/stderr/JSON/JS。今天下游无人按它分支；**这枚事实是"落地腿该钉什么"的坐标**（钉＝不许后来者 `Contains` 它）。测试侧确实有拿带数字的错误串做断言的（尺 `grep -nE 'Contains\((err\.Error\(\)|[^,]*[Ee]rr[^,]*), *"[^"]*[0-9][^"]*"'` ⇒ 16 命中，如 `internal/config/loader_test.go:99` `"line 6"`、`:273` `"99"`、`cmd/wisp/resident_approval_risk_268_windows_test.go:110` `"out of range [31, 3600]"`）＝**断言文案、不是选句**，与 216-a1 名册 §128 那 159 枚同族。

## 甲③ 票面反面教材 `ChatInput.tsx:2205` 在本仓存在吗

**不存在。**两把尺＋两把正控：
- `C1` 工作树 `find . -name 'ChatInput.tsx' -not -path './.git/*'` ⇒ **0 行，rc=0（stderr 空）**；正控 `C2` 同尺 `find . -name '*.tsx'` ⇒ **90** 枚（尺真会出数）。
- `C3` 提交树 `git ls-tree -r --name-only 44f51fc8 | grep -i chatinput` ⇒ **rc=1（0 命中）**；正控 `C4` 同尺 `grep -i chat-composer` ⇒ `frontend/src/components/ai-native/chat-composer.tsx`（本仓确实有同名族件、尺不瞎）。
⇒ **本仓无此件；那半句是外部仓（openchamber）的形，⛔ 当本仓判据。**（与 `docs/evidence/s1/216-census.md:108` 与票 37 行 `216-r1` 自报同向，此处为独立复跑。）
本仓 `frontend/**` 是我零写面（⛔ 碰/⛔ 暂存），但**只读**量了（见下面"甲②′"）：那侧也没有一处按错误文案子串选句。

## 甲④ 结论（一句；⛔ 我翻框、⛔ 我改判据文字）

**"这一格今天的形状＝禁止那半已成立（`A5` 生产 0/130 ＋ `A8` 别名链 0 ＋ 甲② 三 token 无一判定，全带正控），必做那半没有落地物（`A10`：按码选句的钩子 `MessageKey()` 生产消费者 0 枚，错误文案层按 `internal/observe/errors.go:160-162` 注释自己写着推迟到 D23/S6）"** ⇒ 按派单给的三形，它最接近**"已经成立、缺的是钉"**那一形，**但"缺"的那枚东西不是只缺一枚测试**——支撑尺＝`A5`/`A8`/`A10`/`A11`（正控）；**⛔ 判它"成立/不适用/未闭合"归非实现者**（票 41 行已把这格交给 `216-v2`，我⛔ 替编排者排程）。

**必答题（照派单 36 行）**：若按今天的事实翻勾，"凭据是哪一行代码"——**没有那一行**。候选只有两枚且各自不足：`internal/observe/errors.go:163`（`MessageKey()` 存在但 0 枚生产者调用＝钩子空转）与 `cmd/wisp/run.go:83`（`switch observe.ErrorClass(errClass)` 真按码分支，但分的是**退出码**、不是"屏幕/CLI 上那句错误提示"）。⇒ **缺的那枚东西＝"消费错误码去选一句用户可见文案"的那层落地物（今天不存在，票面 `:160-162` 自陈推迟）＋一枚"⛔ 任何人不许 `Contains` 错误串里的数字"的钉**；**归编排者裁**（⛔ 我替他排程、⛔ 我设计方案）。

## 甲②′ 顺带把界面侧也量了（只读，射程＝提交树 `44f51fc8 -- frontend/**`；我零写面）

票面反面教材是 `.tsx`，所以同一族尺换射程再跑一遍（⛔ 混在飞工作树）：
- `F1` `git grep -n '413' 44f51fc8 -- frontend` ⇒ **rc=1（0 命中）**（整个前端树里没有 `413` 这个数字串）。
- `F2` 尺＝前端"按文案子串判定"族：`\.(includes|indexOf|startsWith|match|search|test)\(.{0,3}(rror|too long|exit status|413|429|4[0-9][0-9])`，射程 `frontend/src frontend/scripts`（提交树）＝**0 命中，rc=1、stderr 空**。这把尺的**判别力当场验过**：把票面那形喂给它（`if (msg.includes('413')) showTooLong()`）⇒ **1 命中**（尺会抓）；抓不到的是"本仓没有"，不是"尺瞎"。
- 正控 `F3`＝同族尺 `[.]includes[(]` 在 `frontend/src` ⇒ **32 行**（`out` 末段）⇒ 前端确实在用子串判定、只是没有一处作用在错误文案上。
⇒ **甲② 的"0"在两侧射程都成立**（Go 侧 0/130 ＋ 前端侧 0），禁止那半今天在本仓**任何**射程都找不到实例。

## 乙 `AC#3` 那半：今天有没有生产路径把"名字里的字"当判级输入

先划清：`AC#3` 的前半（控制字符必须洗过）已由 `216-r1` 落地＋编排者复跑（票 40-41 行）；这里只量后半的**可达性**，⛔ 设计方案。

**尺 `D1`＝判级包内所有字符串匹配生产枚数＝37**（射程 `internal/risk internal/perm`，⛔ `_test`）。逐枚答"到不到 perm／risk／门控"（`out:109-145` 全量在场）：

| 枚 | 名字从哪来 | 到不到判定 | 方向（谁受益） |
|---|---|---|---|
| `internal/risk/blacklist.go:242/244/246/248/250/252` | **文件基名** `base := strings.ToLower(baseName(s))`（`:240`） | **到 risk**：`Classify`/`Gate` 的 A/B 判级 | `.env*`/`*.pem`/`*.p12`/`id_*`/`secrets.*`/`Contains(base,"credentials")&&*.json` ⇒ 全是 **deny／要 L2**（更严），**没有一枚能给 ALLOW** |
| `internal/risk/blacklist.go:170` `isUnder`（＋`:227 isBrowserCredentialStore(baseName)`＝`switch base` 四枚**全拼写相等**） | 路径/基名 | 到 risk | deny 方向 |
| `internal/risk/blacklist.go:115-125 overrideApplies` | 单文件覆盖**键** | 到 risk 的 **ALLOW 侧** | `bOverrides[cand.real]` ＝**map 键全等**，非子串 ⇒ 名字里加字不解锁 |
| `internal/risk/rules_shell.go:99 shellAllowlisted`（caller `:41`） | `argv[0]` | 到 risk 的 **ALLOW 侧** | `strings.EqualFold(argv0, entry) \|\| EqualFold(base, entry)` ＝**全拼写/基名相等**，非子串 |
| `internal/risk/rules_network.go:123 hostAllowlisted`（caller `:74`） | host | 到 risk 的 **ALLOW 侧** | `host == entry \|\| HasSuffix(host, "."+entry)` ＝点锚定子域，非任意子串 |
| `internal/risk/provenance.go:996 isRemoteSinkKey` / `:1018 looksRemoteSink` / `:1010 ContainsAny` | **参数键名**／**任意字符串参数值**（含文件名的值） | 到 risk/C25 的 `writeGate`（扫描开/关） | 只能把门**打开（多扫）**：`writeChannelKeys`/`isPathKey` 豁免都要求 `!gateOpen`（`:765/:789`），而 `gateOpen` 由路径形状决定；源码注释 `:753` 自己写"**never from the name the caller gave its payload**（M-7）"、`:1003-1004`"renaming the key cannot close the gate" |
| `internal/risk/taintmatch.go:245 Contains(f.src, w)` | 参数值（可含名字） | 到 risk 污染匹配 | 命中＝**判脏**（更严） |
| `internal/risk/pathresolver.go:175/176/181/195/249/255/257/274/277`、`pathresolver_windows.go:69/77`、`syncdirs.go:282/326/333/338/347/362/368/516` | 路径**书写形式**（`%VAR%`/`~`/`\\?\`/UNC/盘符） | 到 C26 归一与 `IsSyncPath` | 定形状，不定权限；`syncdirs.go:516` 匹配的是注册表 `Personal`/`Business` 键名（非用户文件名字） |
| `internal/perm` | — | **0 枚字符串匹配**（`D1` 射程含它）；`store.go:225` 判的是 `switch err` 哨兵、`:155 PermissionMode() risk.Mode`＝枚举 | **perm 侧根本不读名字** |
| 门控调用点 `D8`（`out:173-181`） | — | 生产入口共 5 枚：`internal/panel/approval.go:65`（Assess）、`internal/tools/mode.go:98`＋`internal/tools/paths.go:381`（Gate/Classify）、`internal/risk/rules_gateway.go:77`（Classify）、`internal/risk/provenance.go:868`（Inspect） | 逐枚见上行归属；`ast.Inspect`（`tools/d22scan/main.go:717`）＝CI 工具，非产品 |

**票 216 那枚真附件名的可达性＝0（这才是 `AC#3` 正题）**：
- `D5` 尺＝`git grep -ci 'attachment' 44f51fc8 -- internal/risk internal/perm internal/tools internal/agent`（⛔ `_test`）⇒ **0 个文件**（`out:162-163`）；正控 `D6`＝同尺 `params` ⇒ `assessor.go:8`/`provenance.go:23`/`rules_gateway.go:1`（`out:165-167`）⇒ 尺会出数。**`internal/perm` 非空目录**（`D9`：`store.go`＋3 枚测试，`out:184-187`），这个 0 不是空转。
- 机制上有交代：`internal/panel/attachments.go:69` 注释 "**DisplayName is what the user sees; it is never used to build a path**"；落盘名是内容哈希：`D7` ⇒ `attachments.go:260` `name := "attachment-" + shortHash(data) + "." + accepted.ext`（扩展名取自闭集表），再过 `memory.ValidArtifactName`（`internal/memory/artifacts.go:91`）。⇒ 带 `已允许`/控制字符的名字**没有一条路走进 `Assess`/`Gate`/`Classify`/`Inspect`**。
- `D3` 尺＝生产里"对 name/base 做 `允许|已允许|approved|allowed` 字面子串判定"⇒ **rc=1（0 命中）**（`out:151-152`）；正控 `D4`＝同尺去掉字面 ⇒ `blacklist.go:240-252` 7 枚在场（`out:154-160`）。全仓 `已允许` 生产命中 4 枚，逐枚＝`cmd/wisp/approval_reply.go:235`（**出向文案**）＋`internal/panel/display_scrub_216.go:12/76/77`（注释），⛔ 一枚作判据。

⇒ **乙的答**：今天**有**"名字里的字当判级输入"的生产路径（`blacklist.go:242-252` 六枚基名图案＝判到 risk 的 A/B 档；`provenance.go:996/1018`＝判到扫描开关），但**枚枚方向都是"更严"**（deny／要 L2／多扫），**ALLOW 侧三处（`rules_shell.go:99`、`rules_network.go:123`、`blacklist.go:115-125`）一律全拼写相等或点锚定后缀，无一处任意子串**；而票 216 正题那枚**附件显示名**到判级包**可达性为 0**（`D5`＋`D7`＋`attachments.go:69`）。⇒ `AC#3` 后半"⛔ 因为名字里有'允许'两个字而改变任何判定"今天**没有被任何生产路径违反**，也**没有被任何用例攻过**（票 42 行具名："另外那一半…没跑过攻它的那一发"）——**可达性我量到此，⛔ 设计那发攻击方案**（等这格交裁决腿）。

## 未量／边界（具名）

1. `frontend/**` 我**只读了提交树 `44f51fc8`**（甲②′ 那三把尺）；**在飞工作树那一版未量**（`frontend/**` 归机主外部前端 agent 在写，⛔ 读它的在飞字节当凭据）⇒ 若前端在飞那版新长了按文案选句的形，本件不覆盖，具名交回编排者。
2. 甲① 分母是**两把尺相加**（130），非严格并集去重；`internal/panel` 与 `internal cmd tools` 射程有交叠（交叠枚＝`internal/panel` 生产 3 枚，尺 `C-ctl`），若编排者要并集口径，需重跑一把跨路径去重的尺。
3. 票 216 判据文字（含 `AC#5` 那句"未逐字复核"）**我一字未动**；改判据要人批。
