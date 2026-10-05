# 268-a2 独立复核：对 `268-a1/census.md` 复跑 AC#0 四问那把尺——逐问判"复现／不复现／口径不同"

编队：`268-a2`（非实现者腿，裁决腿）。本件只回答一件事：census 的 §1 四问与 §6 判语能不能被另一个 agent 独立复跑成立，
从而 AC#0 能不能翻。本件**不实现任何东西、不碰票面 AC 框、不碰台账**；⛔ 全程零 `go` 命令（硬闸，与 `267-r2` 等三枚写腿防撞）。

## §0 起手锚

- 本编队起手：HEAD `c6cf66e6`，`2026-10-05 11:12 +0800`，分支 `dev`，同一枚共享工作树。
- 参照：census 自报复跑时刻 10:24–10:29、HEAD `14735046`。两枚锚点不同 ⇒ census 的行号在本件里一律当**待验断言**，
  下引行号全部是本编队亲读（读的是**工作树**，树里同时有其它在飞腿的未提交行，`cmd/wisp/firstrun.go` 起手即 `M` 脏，
  脏文件命中读数时本件会点名）。
- 本编队读面：`cmd/wisp/{resident_approval_windows,resident_windows,config_reload,logsink,console_windows,main,doctor}.go`、
  `internal/config/{parse,validate,loader}.go`、`internal/observe/{redact,logging,errors,errors_test}.go`、
  `cmd/wisp/resident_approval_risk_256_windows_test.go`（全文）、`scripts/build.ps1:100-173`、票 268、census 全文。
- 禁区自陈：`frontend/**`／`design/**` 零读零引；工单票面 AC 框与 `docs/reports/pending-and-issues.md` 零改动。
- ★**commit 事故（如实自曝）**：本件初版落盘后的 `3f0c4fff` 一笔，`git add` 我按纪律用了显式 pathspec，但 `git commit -F` 提交的是**整个索引**——
  共享树里 `257-r2` 此刻已 staged 的 8 枚文件（`probes/257/r2/` 五枚 log＋`evidence.md`＋`cmd/wisp/firstrun.go`＋`firstrun_257_test.go`，
  逐枚清单见 `git show --stat 3f0c4fff`）被我这笔一并带上盘。内容零丢失、他们的活没有被改写，但**提交署名错了人**。
  按纪律不改写历史，以本段与追加 commit 更正；本腿其后的 commit 一律改走 `git commit -F msg -- <pathspec>` 形。
  教训具名：**显式 pathspec 只守 `add` 不守 `commit`，在共享工作树里等于没守**。

## §1 复跑的尺逐把（命令原文＋读数＋时刻；编号 A＝a2，对账 census 的 C 系）

| # | 尺（命令原文） | 读数 | 时刻 |
|---|---|---|---|
| A1 | `grep -rn riskProvenance cmd/wisp/*.go \| grep -v _test` | 12 行（＝census C1 的枚数）。**成分按我的口径**：声明 4（`:109/:430/:434/:438`）／注释 4（`:103/:419/:431/:435`）／写入 1（`:367`）／return 3（`:449/:457/:461`） | 11:18:35 |
| A2 | `grep -rn riskProvenance --include=*.go .`（剔 frontend/design） | 只命中 `cmd/wisp/resident_approval_windows.go`、`resident_approval_risk_256_windows_test.go`、`.scratch/wisp/probes/260/v1/logs/head-copy-*`（快照件，非产码）。⇒ **面板／托盘／doctor／internal/panel 全树零读者**，比 census C1 的包内尺更硬 | 11:18:34 |
| A3 | `grep -c 'pass("\|fail("\|info("' cmd/wisp/doctor.go` | **28**（fail 15／info 5／pass 8）＝census C6 调用点口径复现 | 11:18:36 |
| A4 | `grep -o '\(pass\|fail\|info\)("[^"]*' cmd/wisp/doctor.go \| sort -u \| wc -l` | **23** ＝census C6 去重具名字面口径复现 | 11:18:36 |
| A5 | `grep -rn "fs.ErrNotExist" cmd/wisp internal/config --include=*.go` | 非测试命中 **5**：`config_reload.go:356`／`firstrun.go:53`（注释）／`firstrun.go:74`／`loader.go:74`（注释）／`parse.go:234`（＋测试与 `firstrun*.go` 若含在飞脏）。census C4 报 3——**口径不同**（见 §2②） | 11:18:36 |
| A6 | `grep -rn "DefaultL1Window" ...` | `queue.go:116` 逐字 `DefaultL1Window = 3 * time.Second`；`MinL1Window = 2s`（`:120`）；零值回落支在 `gate.go:146` | 11:18:37 |
| A7 | `grep -n "fmt.Printf" resident_windows.go resident_approval_windows.go` | resident 腿 stdout 面共 14 行 Printf（`resident_windows.go` 11 行含 `:187`；`resident_approval_windows.go` 3 行＝`:482/:829/:865`）；**`residentRiskGateValues`／`newResidentApprovalWithConfig` 两支零 Printf** | 11:18:39 |
| A8 | 全文精读（非 grep）：`redact.go`、`errors.go`、`logging.go`、`parse.go`、`validate.go`、`loader.go`、`config_reload.go`、`logsink.go`、`console_windows.go`、`resident_windows.go`、`resident_approval_windows.go`、`main.go`、`doctor.go`、`256_windows_test.go`、`build.ps1:100-173` | 逐字行号见 §2／§3 各条 | 11:12–11:18 |
| A9 | `find . -path ./.git -prune -o -iname "*doctor*" -print` | 全仓 **1 枚**＝`cmd/wisp/doctor.go`——census C7"doctor 枚数钉零枚"的尺**复跑复现** | 11:24:05 |
| A10 | `grep -ln "os\.Stdout" cmd/wisp/*_test.go`＋`grep -ln "captureStdout128(" cmd/wisp/*_test.go` | 换 stdout 的在捕面＝`dataroot_128_test.go`（harness 本体）／`leg_sink_nail_131`／`panel_assets_143`；调用 harness 的＝`260r3`／`260r4`——共 5 枚 | 11:24:04–11:26:15 |
| A11 | `grep -ln "bootResidentLeg(" cmd/wisp/*_test.go` | 起真常驻子进程的用例＝ **6 枚文件**：`resident_approval_246`／`resident_ball_228`／`resident_ball_live_228`／`resident_hotkey_live_258`／`resident_sink_nail_127`／`resident_task_source_246` | 11:27:12 |
| A12 | `git show --stat 3f0c4fff` | commit 事故逐枚清单（9 files changed＝本腿 1 枚＋`257-r2` 在飞 8 枚），披露见 §0 末条 | 11:28:12 |

未跑清单（本编队全数具名）：`go build`／`go test`／`go vet`／`wisp doctor` 真跑／任何实机双击。硬闸与 census §3 同因。

## §2 四问逐问复现判定（逐条指回 census 行）

### ① 复现（两处成分口径不同＋一处漏点名）

- census `:15-24` 的锚逐枚我在 HEAD `c6cf66e6` 亲读到，**行号未漂**：`:434` 声明、`:457` 唯一产码点（在 `:453` 的 `err != nil || c == nil` 支里）、
  `:449` NoView 支、`:367` 写字段、`:379-384` 那条 `slog.Info` 读的是**局部变量 `provenance`**（`:380`）——全部逐字成立。
- "产码唯一读者＝那条日志；字段 `ra.riskProvenance` 产码零读者"——**复现且加重**：我的 A2 是全仓 Go 树尺（census C1 只量了 `cmd/wisp/*.go`），
  面板／托盘／doctor／`internal/panel` 零命中。到不到用户眼睛：产码侧到不了，只落盘。
- 口径差一（census `:23-24`／C1）：它把 12 行拆成"声明 2／…／读 1（`:380`）"——`:380` 根本不在这把 grep 的命中行里（命中的"读"是 `:461` 的 return）。
  枚数 12 复现，**成分拆法不复现**（我的 4/4/1/3，见 A1）。结论句"唯一读者是那枚 slog.Info、读局部变量"复现。
- 口径差二（census `:25`，〔仅腿报〕）：测试读者名册**不全**——全文精读 `256_windows_test.go` 后实测为七组：
  `:124/:125`（形状表）、`:149-151`、`:165-167`、`:174-176`、`:247-249`、`:317-319`。census 只点了前四组。
  这直接关乎 AC#2 的写法：字段断言不止在回落用例里，`:247`（种子合法值支）与 `:317`（construction-time-only 支）也在读这枚字段。

### ② 复现（"半分得开、另半分只剩词面"成立；三处口径细节）

- 缺失支＝有机读形状——**复现**：`parse.go:233-234` `fileMissing` 逐字 `errors.Is(err, fs.ErrNotExist)`；`loader.go:67` 用它分支；
  `loader.go:76` `observe.Wrap(ClassConfig, err, "…")` 把原 err 放进 `Error.Err`（`errors.go:184-186`），`Unwrap()` 在 `errors.go:207-212`
  （census 引的 207-212 一字不差）⇒ `LoadFile` 返回值上 `errors.Is(err, fs.ErrNotExist)` 为真。
  生产先例 `config_reload.go:356` 逐字复现；**census 没点名的第二枚先例**：`firstrun.go:74`（该文件起手带外脏，行号按工作树口径）。
  正向钉另有 `settings_257_test.go:286-287`（断"缺失链仍带着 `fs.ErrNotExist`，否则 cause=missing 打不中"）。
- 拒载支＝无正向机读形状——**复现**：`validate.go:145-148` 逐字 `observe.New(ClassConfig, fmt.Sprintf("config.toml: risk.confirm_timeout_sec %d out of range [%d, %d]"…))`；
  `New` 置 `Err` 为 nil（`errors.go:179-181`）⇒ 无 cause；两支 class 同值 `ClassConfig` ⇒ `ClassOf` 不区分（census `:38` 成立）；
  现成分辨路只有 `config_reload.go:396` 的 `case strings.HasPrefix(d, "config.toml:")`——**词面不是类型**，我亲读该 case 逐字在 `:396`。
  ⇒ 具名结论照 census：**"缺失"可正向机读判别；"被拒"今天只能靠词面（前缀或 `out of range` 子串），量到不能机读分开就说不可以**。
  落地可用形状也复现：在 `residentRiskGateValues` 里 `errors.Is(err, fs.ErrNotExist)` 分"缺失／非缺失"是类型级、有生产先例；
  但"非缺失"还含着语法／权限／迁移诸支（票 223 的"四因不许折"纪律），**把它直接改名成"被拒"就是新谎言**——census `:39` 的三条词面脆弱点我全部认可。
- 钉料复现：`errors_test.go:16` 逐字 `if len(classes) != 17 {`、`:76` 逐字 `if len(cases) != 17 {`——两枚 want-17 都在。
  小口径差：census 正文写 `:76`、§2 表 C′ 行写 `:77`，同一件内自相矛盾，我读到的是 **`:76`**。
- A5 口径差：census C4 报"3 命中"，我同目录同词尺读到 5 枚非测试命中——它漏了 `firstrun.go:53/:74`（其一在飞脏）；枚数差**不改②的任何判定**。

### ③ 复现（机制全链成立；结论句一处说得过死）

- tee 与两分支复现：`resident_windows.go:66` 装 `installLogSink(rt.Layout.DataDir)`；`logsink.go:153-161` tee，
  `:157` primary＝`p.Handler()`（`logging.go:87-92`＝`redactHandler` 包标准 `slog.NewJSONHandler`），`:160` mirror＝
  `slog.NewTextHandler(os.Stderr, LevelInfo)` **零脱敏**（`:193-198` 自陈逐字在）。
- 双击不可见复现：`build.ps1:115` `-H=windowsgui` 逐字；无参入口 `main.go:65` 调 `attachParentConsole()`（census 写 `:64`，漂 1 行）；
  `console_windows.go:38` 注释逐字「Fails harmlessly when … no parent console exists」、`rebindStdHandle` `:49-51` `CreateFile("CONOUT$")` 失败即 return ⇒ std 不改绑；
  `resident_windows.go:55-57` 逐字「double click the icon, no terminal attached, stderr going nowhere」。⇒ 双击起法下 stdout/stderr 都无落点，**成立**。
- 形状差复现（census `:46` 那格，本件对落地腿最有用的一格，我逐字亲读两侧）：`[hotkey]` 支在 `resident_windows.go:185-187`
  同一闭包里 `slog.Warn` **加** 一行 `fmt.Printf`；`[risk]` 支在 `resident_approval_windows.go:454-456` 只有 `slog.Warn`、零 Printf（A7 实测：
  含构造链的 `newResidentApprovalWithConfig`/`residentRiskGateValues` 两函数体内 Printf 计数为 0）。
- **不同意 census `:47` 的结论句**（"今天真能过得了的只有 primary 那一支"）：按它自己 `:44-45` 的机制，从终端起常驻腿时
  `attachParentConsole` 成功 ⇒ mirror 那行 TextHandler 输出**当场可见**（含未脱敏原句，见 §3）。
  准确说法是：**双击起法＝只有 JSONL 落盘；终端起法＝JSONL 之外 stderr 还有一遍转瞬可见的镜像**。"用户要自己翻日志"只对双击起法成立。

### ④ 复现 28／23 两口径；"实印"口径不复现〔仅腿报〕；退码链全复现

- 调用点 **28**（fail 15／info 5／pass 8，A3）与去重具名字面 **23**（A4）——census C6 两枚数复现，census `:52` "三个数不是三个版本，是同一棵树的三把尺"的方法论我照收。
- 但第三口径"实印 14～16 枚"（census `:52`，〔仅腿报〕）**我静态枚举复现不了**：`cmdDoctor` 的每条 append 都进 `results` 一次、`:118-123` 每元素打一行，
  全路径枚举只有两个结局——`deps.toml` 不在＝**11 行**、在＝**13 行**（三枚 pin 各 1 行）。枚数构成：3 info + gcc 1 + sherpa 1 + ort 1 + ort 版本 info 1 + DLL 2 + data dir 1 + deps 段（1 info 或 3）。
  ⇒ 判 **不复现**（除非腿那两把尺另有定义，census 未给命令原文）。
- "有没有一处会读配置加载错误"＝**没有**，复现：`doctor.go:3-16` import 块无 `internal/config`（我逐行读：errors/fmt/os/os-exec/path/filepath/runtime/strings/buildinfo/proc/sherpa/go-toml），
  全文只碰 `deps.toml`/`portable.txt`/`doctor-write-probe.tmp`，不触 `config.toml`。
- 构建门链复现（逐字亲读）：`build.ps1:167-169` smoke 跑 `wisp.exe doctor`、非零即 `Fail`；`doctor.go:34-36` `fail()` 置 `critical:true`、
  `:120-121` 任一 critical FAIL ⇒ `:125-128` 返 false ⇒ `main.go:97-98` `os.Exit(1)`。
  ⇒ **给 doctor 加一条 fail 级检查项，会把"数据根里躺着被门拒的配置"的那台机器的构建门打红**——条件真（红在拒载值存在时；CI 无配置＝missing，不红），census `:55` 的警告成立，C 形只能走 `info()`。
- 名册钉口径修正：census `:56` 从 `256_windows_test.go:560/:563`（`fed != 1`／`bare != 0`，逐字复现）退出"不许在 `runResident` 里加第二枚配置读取点"——
  **射程画宽了**。那两枚是 AST 尺（`:537-558`），只数 `newResidentApprovalWithConfig`／`newResidentApproval` **这两个构造名**在 `runResident` 体内的调用次数；
  在 `runResident` 里新加一枚裸 `config.LoadFile` 不撞它。对本票真正的钉是：**别多调一枚构造函数**。
- 枚数钉零枚＝**复现**（A9：`find` 全仓只命中 `cmd/wisp/doctor.go` 本身，doctor 无测试文件、无 `len(results)` 型断言）。
  但 census §④ 没说的一个边界本件补上：**doctor 的全量 stdout 已被一枚消费者捕获**——`dataroot_128_test.go:132`
  以 `captureStdout128(t, func(){ ok = cmdDoctor() })` 驱动 refusal-legs 表的 cmdDoctor 那一腿。该腿跑在"用户配置目录不可得"
  的重绑机上（`failConfigDir128`），doctor 改动前本就返 1，断言为 rescueMarkers 的 contains 型 ⇒ 加 info 行不撞、加 fail 行也不打红**这一发**；
  构建门那发的红判（§2④ 上文）不变。

## §3 那句 `err` 透传的独立证真——**证真（四支链全码读通；⛔ 未跑起来，跑即违规）**

census `:29-31` 判：回落前先 `slog.Warn(…, "err", err, …)`（`resident_approval_windows.go:454-456`），而 `redact.go` 的 `Attr` 对 `key="err"`
落 `KindAny` 原样透传 ⇒ 拒载原句今天已在盘上 JSONL。我没有照它的判语，把链拆成四支逐支独立核：

1. **attr 怎么生成**：`slog.Warn` 的键值对经 `log/slog` 折叠成 `Attr`，value 走 `AnyValue(err)`。
   `err` 的动态类型是 `*observe.Error`（`loader.go:76/:81/:102/:118/:121` 与 `validate.go:146` 都返回它）——
   我读了 `errors.go` 全文：该类型只有 `Error()/Unwrap()/Retryable()` 三方法，**没有 `MarshalText`、没有 `LogValue`、不是 `[]byte`**
   ⇒ 不入 slog 的任何特判支 ⇒ `Kind() == KindAny`。
2. **词表匹配穷举（三种 key 词表逐支）**：`redact.go:55-69` 三枚词表——`secretWords` 13 枚（key/token/secret/password/passwd/pwd/auth/authorization/credential/credentials/bearer/passphrase）、
   `audioWords` 5 枚（audio/pcm/samples/waveform/wav）、`contentWords` 5 枚（body/html/content/document/payload）。
   `keyWords("err")` 切出唯一词 `err`，**三表皆无**（`hasWord` 按全词匹配、非子串；反例核过："hotkey" 一枚词、不中 key）。⇒ 不进 `:109-116` 任何改写支。
3. **分派支穷举**：`Attr` `:118-128` 四支——`KindString`（不中，err 非串）／`KindAny`+`[]byte`（不中）／`KindAny` 其余 ⇒ `:125` 逐字 `return slog.Attr{Key: key, Value: v}` **原样返回**／`default` 同形 `:127`。
   上游 `redactHandler.Handle`（`logging.go:132-139`）对**每一枚 attr** 过 `red.Attr`，无旁路 ⇒ 透传点成立。
4. **落盘面会渲染它**：primary 的 next 是**标准** `slog.NewJSONHandler`（`logging.go:89`），对 `KindAny` 走 `encoding/json`；
   `*observe.Error` 字段全导出且带 json tag：`class`、`detail`（`Detail`＝那句拒载原话）、`Err` 标 `json:"-"`（cause 不外泄）。
   ⇒ 盘上一枚 JSONL 里 `"err":{"class":"config","detail":"config.toml: risk.confirm_timeout_sec 20 out of range [31, 3600]"}`。
   mirror 侧 TextHandler 同值渲染（`logsink.go:160`，不过脱敏）。

**我复算支持那句透传**，并两处加严：
(a) census 验到 1-3 支（它只自报读了 redact.go 的 107-127），第 4 支（JSONHandler 序列化 + `Detail` 可导出落盘）它没有点名——我补核，链闭合；
(b) 透传**绕开了 rule-4 的 512 截断**：`KindString` 支才走 `r.String()`（`redact.go:119-120`、`MaxLoggedString=512` `:37`），
KindAny 原样返回后由 `json.Marshal` 全量序列化 ⇒ 长 Detail 不截断。对本题是好消息（句子完整在盘），对"日志行数有界"那条规则是既有缺口（具名登记建议，不由本件裁）。
由此 census `:31` 对票面"且不说"的更正——**盘上有、名字错（unreadable）、用户手上没有**——三短句我全部独立复认。

## §4 U4：给常驻腿加一行 stdout 会不会打红既有 stdout 形状断言——静态撞钉名册（已收口）

尺＝枚举一切捕获 stdout 的用例（在程换管／子进程管道），逐枚读它断什么、是否在"配置读失败态"构造常驻门（census U4 未穷尽那格归本件）。
⛔ 零跑测：以下判定全部静态出身，颜色一枚不背书（具名在 §6）。

**名册 11 枚文件（A10＋A11），分四族：**

| 族 | 文件 | 断的是什么 | 拒载/缺失态会构造门吗 | 加一行 stdout 撞吗 |
|---|---|---|---|---|
| 在程换 stdout | `dataroot_128_test.go` | `captureStdout128` harness 本体＋cmdDoctor 腿（contains 型 rescueMarkers） | 否（不构造门） | 不撞 |
| 在程换 stdout | `leg_sink_nail_131_windows_test.go` :163-196 | 驱动的是**有返回码的控制台腿**（run/models/providers/panel-inbound）；`:240 if n != 1` 数的是 sink 目录 **jsonl 文件枚数**、不是行数 | 否（那些腿不建常驻门） | 不撞 |
| 在程换 stdout | `panel_assets_143_test.go` :85-91 | `cmdPanelAssets` 写临时文件 | 否 | 不撞 |
| 在程捕获调用方 | `resident_cancel_key_wording_260r3_windows_test.go` :98/:140/:164/:178 | bindBallHost 审计句；:103-105 的逐字 want 走的是 **`strings.Contains`**、非相等 | **否**——全部用裸 `newResidentApproval()`＝NoView 支（构造在捕获块**之外**） | 不撞 |
| 在程捕获调用方 | `resident_cancel_key_label_260r4_windows_test.go` :86/:128 | 同上族，contains 型 | 同上 | 不撞 |
| 常驻子进程（CI 跑） | `resident_sink_nail_127_windows_test.go`（3 发） | stdout 全 `lockedBuf.has`（contains）；stderr 枚数钉数的是 **winsec 封缝句**那一枚串（`Count(...)!=1` :493）；JSONL 的 record 0/1 位置钉——`[risk]` 两行日志本就排在 install 记录之后（门在 `resident_windows.go:132` 构造、sink 在 `:66` 先装），stdout 多一行**不动 JSONL 索引** | **会**（`t.TempDir()` 无 config.toml ⇒ fallback 支触发） | 不撞（附带条件见下） |
| 常驻子进程（CI 跑） | `resident_approval_246`／`resident_task_source_246`／`resident_ball_228` | 卡片挂起句／任务来源姿态／球姿态 claims，均 contains 型（如 `:119 residentReachedLoop`、`:66 cardRaisedClaim`） | 会 | 不撞（同条件） |
| 常驻子进程（winlive，默认 CI 不跑） | `resident_ball_live_228`／`resident_hotkey_live_258` | hotkey verdict 片段；census 引的 `:69` 断的是 `"provenance=defaults"`——那是 `[hotkey]` 桥接行的 console 形状，**与 `[risk]` 无关**（〔仅腿报〕那枚我亲读定性） | 会 | 不撞 |

**关键负空间**：在程触发 fallback 支的用例只有 `256_windows_test.go` ruler① 的 "unreadable" case（`newResidentApprovalWithConfig(t.TempDir())`），
而该文件**根本不在捕获名册里**——加行后它只会往 `go test` 的真 stdout 淌一行无人捕获的文本，判定无影响。

**枚数与结论**：捕获名册 11 枚文件／常驻子进程 6 枚／其中在拒载-或缺失态构造门的 6 枚；断言形状逐枚读尽为 contains、或针对**其它字符串／文件枚数**的精确计数——
**静态判定："只加一行 stdout Printf" 撞 0 枚钉**（census U4 的"落地腿开工前必做一次撞钉预检"这一格，本件以静态尺交回大半）。

具名留给写腿的两条残余（本件量不到）：
1. **词面条件**：新句措辞不得**吞任何既有 claim 子串**（`resident event loop running`／`exited through the D38(e) shutdown order`／
   `wisp: 卡片挂起：`／`hotkeys from`／`wisp: persistent log sink installed`／winsec 封缝句）——contains 族会被误命中；这是词面约束，静态尺只能警告、不能替它钉绿。
2. **一发真跑**：`go test ./cmd/wisp` 整包颜色（尤其 127 族的 record 位置断言在新 slog 记录加入时的表现）——本腿禁跑，交回写腿，与本件 §2③ 的"抄 hotkey 形状"最小组合正交：那只加 stdout 行、不加 slog 记录时，残余 2 亦应绿。

## §5 我不同意 census 的地方

1. **`:33` 那枚"同源瑕疵"是误读，不该进任何票面**：它说 `fallback` 字面 `DefaultL1Window=3s` 与"钳位下界 `MinL1Window`（现读 2s）"对不上。
   实测回落支走的是 `gate.go:146` 的**零值→`DefaultL1Window`**支，而 `queue.go:116` 逐字 `DefaultL1Window = 3 * time.Second` ⇒ 那句字面**数值与常量名双双正确**；
   它真正可挑剔的只有"不在常量引用链上＝改常量不红"（散文钉），与"2s/3s 之错"无关。census 把 3s 的 schema-tag（`l1_window_sec` 文件沉默态）和 300s fallback 的常量链搅在一起了。
2. **`:56` 的 `fed/bare` 射程**：构造名计数≠"配置读取点"禁令（见 §2④ 末条）。
3. **`:47` 的"只有 primary 一支今天真能过得了"说得过死**：终端起法下 stderr mirror 是一枚转瞬可见面（见 §2③）。
4. **`:25` 测试读者名册不全**：漏 `:247-249/:317-319` 两组字段读数（见 §2①）。
5. **"实印 14～16"静态复现不了**（11/13，见 §2④）；该数在 census 里本就标〔仅腿报〕，我的判语是**不复现**而非"口径不同"。
6. 小项：`loader.go` 那句注释在 `:72-73`（census 引 `:74`，漂 1 行；`fs.ErrNotExist` 字面在 `:74` 我读到——**该句跨两行，引哪枚都要带口径**）；
   `main.go` attach 在 `:65`（census `:64`）；`errors_test` 第二枚钉在 `:76`（census 表内 `:77`）。

## §6 量不到（具名；禁跑 go 的射程损失）

- 全部既有钉**今天**的红绿颜色——零 `go test`，本节任何"钉在／射程"判断都出自码读，不背书颜色（与 census U2 同格）。
- `err` 落盘的**现场样本**：链是读通的（§3 四支），但我没在真 JSONL 里引一枚现成行——要一枚有一枚就得跑或翻在飞腿的台件。
- `wisp doctor` 实印行数——静态枚举 11/13 是**路径覆盖枚举**，非实跑读数。
- `attachParentConsole` 实机行为（census U3 同格：双击 stderr 真零落点是静态推理）。
- 页面侧有没有承接位（census U1）——**本编队量不到**（⛔ 不读 frontend/design）。
- U4 名册里"今天跑起来会不会真红"——只能静态判形状，最后一发留给写腿（本节具名交回）。

## §7 自我对抗（≥5 条真改）

1. 我初稿准备照 census 收编"DefaultL1Window=3s 瑕疵"——读了 `queue.go:116`/`gate.go:146` 后**整条翻反**（§5.1）。
2. 我初稿把 `fed/bare` 当"配置读取点钉"抄——读全 `:535-572` 的 AST 尺后收窄为"构造名计数钉"（§5.2）。
3. 测试读者名册从 census 的四组补到七组（§2①）。
4. 拒载 err 的 Kind：我先假设 KindString（因为 err"像句子"），被 `errors.go` 全文证伪——无 `MarshalText`/`LogValue`，是 KindAny（§3 第 1 支）。
5. "实印 14–16"我一度按〔仅腿报〕带过——自己枚举 `results` append 路径后判**不复现**并给出 11/13＋构成（§2④）。
6. C4"3 命中"我以同词同目录尺复跑成 5 命中——把它降为"口径不同、不改判定"，并点名 `firstrun.go` 在飞脏对读数的污染（§2②）。
7. census §3 表把它自己的 `:76` 写成 `:77`——我不合并两种写法，以亲读为准（§5.6）。
8. 131 的 `:240 if n != 1` 我初读当作"行数精确钉"（若是，U4 结论要翻成"高危"）——重读 `:225-262` 确认数的是 **jsonl 文件枚数**，把判定从"高危"改回"不撞"并写进名册（§4）。
9. `dataroot_128_test.go:132` 那枚 `cmdDoctor` 全量捕获是 census §④ 没点的洞——不推翻它的"枚数钉零枚"，但把"doctor 输出无人消费"的隐含说法收窄了（§2④）。
10. 首笔 commit 的 `git commit -F` 整索引事故（§0 末条）是真改动而非修辞：本件因此新增"commit 也须 pathspec 形"一条纪律自陈，A12 附原命令清单不涂改。

## §8 判语：AC#0 四问能不能翻

**能翻**，条件是翻到 census 的**修正后版本**、不是票面原文那把尺：

- ①③④的机制与行号我在 `c6cf66e6` 逐枚独立复现（②半分离结论独立复现），census §6 三行判语的**重心更正成立**：
  缺的不是记录（拒载原句今天已带 `err` 落进 JSONL，§3 四支链证真），缺的是（a）一枚分得开的名字、（b）一处双击用户真会读到的面。
- **U4（派给本件的那格）已由 §4 静态尺收口：撞钉 0 枚、词面残余与一发整包真跑具名交回写腿**；doctor 枚数钉那把 find（A9）也已复跑复现。
  census §6 列的其余判不动格归口不变：**U1 本编队量不到**（不读 frontend/design）、**U3 等实机**、**U5 不翻 AC**。
- 落地形最小组合仍是 census §6.2 的 A＋B（抄 `resident_windows.go:187` 那枚形状），C 只能 `info()` 级；
  `DefaultL1Window=3s` 那枚"瑕疵"**不得随件落地**（§5.1，它不成立）。
- ⛔ 本件不代翻 AC#0 的框——框归编排者按 §2 逐问核勾；本件只回答"独立腿复跑，它撑不撑得住"。撑得住，附 §5 六条修正。
