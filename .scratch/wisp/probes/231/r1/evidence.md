# 票 231 · 写码腿 `231-r1` 交件与读数（落笔前锚 + 每格判语 + 尺读数）

- 腿＝`231-r1`（写码腿）；本件＝过程记录与读数，**不是裁决件**（裁决＝非实现者另做，SPEC-12 §4.3 #1/#3）。
- ⛔ 票面五枚 `- [ ]` 框一枚未翻（翻框归编排者）。
- 起手锚见下一节；落笔顺序＝①产品分支 ②表里加一行 ③再证它会红（派单指定的顺序，因为票面 AC#3 前半被普查件推翻后，正控今天无处可跑）。

## §0 起手锚（落笔前现量）

- HEAD＝`b5418b6c886484caee1d706c14d68b660c230929`（`git rev-parse HEAD`）。
- 现量时刻＝`2026-10-05 23:0x +08`；起手 CPU／MEM＝`31.4% / 68.0%`（`Get-Counter`＋`Win32_OperatingSystem` 现读，均 <70%）。
- 起手 scoped porcelain（`git status --porcelain -- cmd/wisp internal tools`）＝**0 行**（⛔ 无别人的未提交改动撞 `config_reload.go`，票 232 亦零落地）。
- 三处锚的现读行号（落笔前，全部 `grep -n` 现坐实）：
  1. 误吞的那条 case＝`cmd/wisp/config_reload.go:396` `case strings.HasPrefix(d, "config.toml:"):`，它回的 `return` 在 `:397-399`（票面说 `:357`＝漂 ＋39，与普查件 `probes/231/a1/census.md` §0 表 #2 对上）。
  2. 钉句的那枚表条目＝`cmd/wisp/config_sentences_223r2_test.go:66-69`（`"声明未来版_正文语法坏_J1"` 在 `:67`），**一枚 J1＝(a) 形**：`cause=syntax`／notCause `cause=invalid`／notFragment `语法没问题`（票面 AC#3 前半说它钉 (b)(c)(f) 三形＝不实，普查件 §3.2/§3.3 已推翻）。
  3. loader 那句原文＝`internal/config/loader.go:120-124`，产码英文逐字（`:122`）：
     `config.toml: schema_version %d was written by a newer build (this build understands %d); upgrade Wisp or restore a backup`
     ⛔ 本腿一字未改（钉手 `internal/config/loader_test.go:273` needle `"99"`、`internal/config/migrate_test.go:160` needle `"newer build"`）。
- 分类器出口枚数尺（落笔前）＝`grep -c "cause=" cmd/wisp/config_reload.go`＝**7**（票面说"四句"＝枚数错）。
- `cause=invalid` 中文全文的测试侧钉＝**零枚**（尺＝`grep -rn "内容被校验拒绝" cmd internal tools --include=*.go`＝**1**，唯一命中＝产码自己 `:398`）⇒ 既有四句"一字不改"今天靠票面不靠门 ⇒ 本腿另用逐字节尺自证（§6）。

## 节状态

| 节 | 内容 | 状态 |
|---|---|---|
| §0 | 起手锚 | 已写满 |
| §1 | AC#1 现复现（操作员实际看到的那句原文） | 已判（实跑两发：改前／改后） |
| §2 | AC#2 新出口落点与内容三项 | 已判 |
| §3 | AC#3 常驻钉（加行）＋"拿掉分支必红"的正控 | 已判 |
| §4 | AC#4 前缀身兼两职的登记（只登记不实现） | 已判 |
| §5 | AC#5 整包终态与逐名红名册 | 待读数落地后补齐 |
| §6 | 四枚既有句的逐字节零变化尺 | 已判 |
| §7 | 门禁四数 | 已判 |
| §8 | 判不动的地方／本腿没做的 | 已判 |

---

## §1 AC#1 现复现（实跑读数，不是票面那张表）

尺＝overlay 注入一枚**只读数、零判等**的台件 `zz_231r1_ac1_probe_test.go`（物理件在本目录，
经 `go test -overlay .scratch/wisp/probes/231/r1/overlay-ac1.json` 映射成
`cmd/wisp/zz_231r1_ac1_probe_test.go`；工作树里从未存在这枚文件，`git status -- cmd` 全程 0 行）。
它走的是票 223 那套真宿主：`newReloadRun223` → `live()` → `runTextTask` → `cmd/wisp` 装的 1s tick，
种 `schema_version = 99\n\n[ball]\nsize = 64\n`（正文合法且解析得开＝唯一能种到 `loader.go:120` 的形状，
正文坏则被 `:104` 抢先归 `cause=syntax`＝票面要保留的 (a) 形）。

**改前（HEAD `b5418b6c`，取数 `2026-10-05 23:14 +08`，日志 `logs-ac1-before-fix.txt`）操作员实际看到的那句原文逐字**：

```
[audit] config: HOT-RELOAD state=not-applied cause=invalid detail="config.toml 语法没问题，但内容被校验拒绝（值不合法或引用解不开）。本次运行继续用内存里的旧配置"
```

⇒ 票面现象**复认成立**，且"第二句还是假的"这一条有物证：那一形根本没走到 `decodeStrict`/`validate`
（`loader.go:120` 就返回了），"值不合法或引用解不开"对它不成立。
同一条 stderr 上 stdout 侧**零枚**相关行（七句全部只走 `:155` 的 `auditf`＝`[audit] ` 前缀＋JSONL，
`cmd/wisp/run.go:931-935`），与普查件 §1 的前提逐字对上 ⇒ 本腿⛔未新造任何 stdout 打印。

**改后（`290dca87`＋`a16d1ff7`，取数 `2026-10-06 09:0x +08`，日志 `logs-ac1-after-fix.txt`）同一枚台件的读数**：

```
[audit] config: HOT-RELOAD state=not-applied cause=newer-build detail="config.toml 是由一个更新的 Wisp 写出来的（它声明的 schema_version 比这份程序懂得的高；这一条不说语法错，也不说值不合法，因为它还没走到校验）。本次运行继续用内存里的旧配置；升级 Wisp 或恢复备份才会读它"
```

判语：**AC#1 成立**（改前那句逐字抄回，非票面表；改后同一形状走自己的出口）。

## §2 AC#2 给它一条自己的出口

- 落点＝`cmd/wisp/config_reload.go:396-423`（新 `case` 现读 `:396`，`return` 在 `:420-423`；原 `:396` 那条 HasPrefix 顺移到 `:424`）。
  ⚠ **位置承重**：必须写在 HasPrefix 那条**之前**——写晚＝死代码，且新那一行会红在 `:110` 那把前缀尺上
  （这条不是推测，§3 的正控红句原文就是它）。
- 认领形状＝`strings.Contains(d, "was written by a newer build")`＝**短语而不是前缀**。尺：
  - 该短语在 `internal/` 非测产码里＝**唯一一枚**（`grep -rn "was written by a newer build" internal/ --include=*.go | grep -v _test.go`＝`loader.go:122`），
    cause 词根与该句原文同源（编排者已裁 `cause=newer-build`，理由在派单第 4 条）。
  - 它抢的那根前缀 `config.toml: ` 在 `internal/config` 非测产码里身兼 **48 枚** detail
    （尺＝`grep -rn '"config\.toml: ' internal/config/ --include=*.go | grep -v _test.go | wc -l`＝48）
    ⇒ 新分支严格窄于 HasPrefix 那条，**不可能**把 unknown-key／值不合法／引用解不开／迁移产物不合法任一形抢过来。
- 句子三项（票面 AC#2 逐字要求）：①由更新版本写出（"是由一个更新的 Wisp 写出来的"）
  ②本次运行继续用内存里的旧配置 ③升级或恢复备份才会读它；
  另照 `:357`／`:361` 两支的既有形状加了一句"说清自己不是什么"（"这一条不说语法错，也不说值不合法，因为它还没走到校验"）。
- ⛔ 未改 `loader.go:120-124` 一字；⛔ 未把 `:424` 那条 `HasPrefix` 改宽；⛔ 未动既有四句与 unknown-key／migration／unclassified 三支。
  尺＝`git diff -U0 -- cmd/wisp/config_reload.go` **只有一个 hunk** `@@ -395,0 +396,28 @@`＝纯插入、零删零改；`git diff --numstat`＝`28 0`。
- marker 互斥自证：新串 `cause=newer-build` 不含 `cause=missing`／`syntax`／`unknown-key`／`invalid`／
  `permission`／`migration`／`unclassified`／`state=disabled` 任一枚作子串（普查件 §4.4-B 的命名禁忌）；
  也不含被禁的折叠语 `配置未生效`／`重启就好`（`firstrun_257_test.go:342-352` 那两枚尺管的面，本句亦避开）。

判语：**AC#2 成立**（有自己的出口、走同一条审计通道、三项齐全、禁改项零触碰）。

## §3 AC#3 同时改那一处钉句的断言（实为加行）＋正控

- 票面 AC#3 前半句"那三形必须改成断新出口"**不成立**（普查件 §3.2/§3.3 现读推翻）：`:66-69` 今天只有
  一枚 J1＝(a) 形，按字面改它＝删掉 (a) 唯一的常驻钉、违票面 AC#3 自己的后半句。
  ⇒ 本腿按 223-v2 自己给的改法（该件 `:215`＝"**在表里加一行**"）执行：⛔ J1 一字未动，
  在 `config_sentences_223r2_test.go:87` 之后**加一行**：
  `"声明未来版_正文解析得开_归更高版本自己那句", "schema_version = 99\n\n[ball]\nsize = 64\n", "cause=newer-build", "cause=invalid", "语法没问题"`
  （body 建议取自 223-v2 的 (f) 形，派单第 3 条指定）。三把尺（加行后现读 `:124` 前缀／`:128` 反向 cause／`:132` 反向半句，
  原址 `:110`/`:114`/`:118` 因本腿那一行整体上移 14 行）全用现成的，未新增尺、未放宽任何断言。
- 加行后的正向读数（`2026-10-06 09:0x +08`）：`go test -run TestTicket223R2FailureSentenceRouting -count=1 -v ./cmd/wisp`
  ＝ `--- PASS` 父用例＋**9 枚子形全 PASS**（原 8 枚＋新 1 枚）。
- **"拿掉 AC#2 那分支必须红"的正控（`-overlay` 注入删干件，工作树从未被写）**：
  脚本 `mutation/make-no-branch.ps1` 从**当前** `config_reload.go` 现切
  `mutation/config_reload_no_branch.go`（切断后自检：文件里已无 `"was written by a newer build"` 那枚 case），
  `mutation/overlay-no-branch.json` 把它映射回 `cmd/wisp/config_reload.go`，再跑同一枚用例（无覆盖标志，
  避开本仓实测的 `-overlay`×`-cover*` 静默忽略坑）。
  取数 `2026-10-06 09:0x +08`，日志 `logs-ac3-mutation-red-final.txt`，**红句原文逐字**：

  ```
  config_sentences_223r2_test.go:125: shape 声明未来版_正文解析得开_归更高版本自己那句 is booked "cause=invalid detail=\"config.toml 语法没问题，但内容被校验拒绝（值不合法或引用解不开）。本次运行继续用内存里的旧配置\"", want it to open with "cause=newer-build" - the two sentences swapped again
  ```

  同一跑里另外 8 枚子形**全部仍 PASS**（尺＝`--- FAIL` 枚数 **1**、子形 `--- PASS` 枚数 **8**）
  ⇒ 这颗牙只咬新那一形。同一跑还并发出另外两把尺（`:129` borrows `cause=invalid`、`:133` says `语法没问题`），
  三把尺全响＝派单说的"不许拿现有 (a) 那枚断言冒充这一格的牙"这一条有独立正控。
- 派单第 3 条的推论也自证了一遍：**改前**（无新分支、无新行）拿掉分支不会红任何一枚——今天这张表里
  `cause=invalid` 只作 `notCause` 出现 6 次、从未作"必须是"被钉（普查件 §3.2 ⇒ 本腿复算＝
  `grep -c "cause=invalid" cmd/wisp/config_sentences_223r2_test.go` 命中全在 `notCause` 位）。
- 同批（不改不红、但牙会变钝）：`cause=newer-build` 已进两张互斥名册
  ＝`config_reload_223_test.go:570-574` 的 `all`（8 枚→**9 枚**）与
  `config_reload_perm_223_windows_test.go:112-119` 的名册（6 枚→**7 枚**）。
  这两枚负钉名册跑在整条 trail 上，所以新 marker 进册后，"版本更高"那形将来误吞别的形状会有人响。
  ⚠ 派单第 5(i) 条要求的正是这一处；5(ii) 那条暗坑（perm 用例 `:142-143` 把 `state=not-applied cause=invalid`
  当**控制流探测器**）本腿自查**未撞**：新分支认领的是短语而不是前缀，perm 用例种的 body 是 `ver = 2`
  合法体（`config_reload_perm_223_windows_test.go:81-82`），里面没有 "was written by a newer build"，
  ⇒ 探测器看到的仍是 `cause=invalid`／`state=applied` 那两条既有形状，实测该用例改后连跑 3 发
  （`logs-ticket223-repeat-6/7.txt` 及 `-run TestTicket223PermissionDeniedSitsInItsOwnSentence`）各 2.3–2.8s PASS，
  ⛔ 没有以 `three plants all landed…` 或 `neither the permission sentence nor an adoption arrived…` 红。

判语：**AC#3 成立**（新增那一行是常驻钉；拿掉分支必红，红句逐字在上；J1 与既有断言零改动）。

## §4 AC#4 前缀别再身兼两职——**只登记，不实现**

**今天这一形靠什么避免再次误吞**：`config.toml:` 这根前缀在 `internal/config` 非测产码里身兼 **48 枚**
互不相同的 detail（loader 3／parse 8／validate 21／migrate 4／catalog 11／unwired 1，逐枚尺见普查件 §1 末表），
而分类器今天有 **五** 支在内层 switch 里按词面认领（`:369` HasPrefix "config.toml parse"／`:373` Contains
"unknown key"／`:377-378` 两串或／**本票新增 `:396` Contains "was written by a newer build"**／`:424` HasPrefix
"config.toml:"）。⇒ 本票的加法**没有消除**前缀身兼两职，只是给 48 枚里的一枚抢回了自己的句子；
"新增一支必须早于 `:424`、且认领串不得与既有支重叠"这件事，今天仍然只写在注释与名册里，**没有门**。

**登记（缺尺清单，逐条具名）**：

1. **零枚尺数得住"同一个失败被印了两条审计行"**：`describeReloadFailure` 的返回只被 `:155` 一枚 `auditf`
   消费；两张名册只查"缺席的 marker"，不查"重复的行"。⇒ 新支若与 `:424` 认领串重叠且写在其**之后**，
   操作员得到**零条**（被吞）；若写在其**之前**但认领过宽，别的形状会被抢走。**方向不可判**＝今天没有尺。
2. **`cause=invalid`／`missing`／`unknown-key`／`migration` 四句的中文全文在测试侧零枚被正向钉**
   （尺＝`grep -rn "内容被校验拒绝" cmd internal tools --include=*.go`＝1，命中即产码自己 `:426`）。
   ⇒ "一字不改"今天靠票面、不靠门；本腿用 §6 的逐字节尺自证，但那只对本票这一笔有效。
3. **同族第二吞点今天仍未处理**：`internal/config/migrate.go:71-73`／`:78-80`
   （`config.toml: migration to schema version %d produced an invalid config (%v); the file was left untouched`）
   不含 `cannot migrate from schema version`／`no migration registered` 任一串 ⇒ **今天照样被 `:424` 吞成
   `cause=invalid`**，与"版本更高"同一根前缀、同样零常驻钉。**票 231 字面没要求处理它 ⇒ 本腿不动。**
4. **可选的机读路（本腿未采纳，具名交回）**：`observe.Error.ProviderCode` 是已导出字段、全仓非测产码
   零枚赋值 ⇒ 填它不新增导出名；但一旦非空，`Error()` 句首会渲染成 `config (<code>): …`
   （`internal/observe/errors.go:194-196`）⇒ 任何按 `err.Error()` 全文比对的既有断言形状会变。
   另一条路（按 `observe.ClassConfig` 分类码归句）**不可用**：`internal/config` 非测源码 85 枚错误点共用这一枚分类码。
5. **本仓已记录在案的取向**（登记时引它比新造说法稳）：票 268 在 `cmd/wisp/resident_approval_windows.go:447-456`
   写下"因为不能读 error prose 所以故意不细分"，并立了一枚禁 `strings.*` 分类的 AST 尺
   （`cmd/wisp/resident_approval_risk_268_windows_test.go:224-293`，射程只到那一枚文件）。

**哪枚文件哪一行为什么今天不动**：

| 具名位置 | 今天不动的理由 |
|---|---|
| `cmd/wisp/config_reload.go:424`（`case strings.HasPrefix(d, "config.toml:")`） | 派单/票面 AC#2 逐字禁令："⛔ 不许把 `HasPrefix` 那条改宽去'顺手兼容'"。改窄＝动它的认领面，而 `config_reload_perm_223_windows_test.go:142-143` 正把 `state=not-applied cause=invalid` 当**控制流探测器**用；收窄它需要先看那枚探测器还认不认得出"tick 读到了文件"，那是票 232/新票的地界，不是本票 AC#2 射程。 |
| `internal/config/loader.go:120-124`（`observe.New(ClassConfig, …)` 那支） | AC#2 明令不改那句给日志的原文；且 `ProviderCode`／分类码路线落在这一枚文件会改 `err.Error()` 的渲染（见上面登记 #4），钉手＝`loader_test.go:273`、`migrate_test.go:160`。 |
| `internal/config/migrate.go:71-73`／`:78-80`（迁移产物不合法那两枚吞点） | 票 231 字面只覆盖"版本更高"那一形（AC#2 的三项内容对它才成立；"迁移产物不合法"的修法既不是升级也不是"由更新版本写出"），按登记 #3 移交后续票。 |
| `cmd/wisp/config_reload.go` 的函数头注 `:342-353` | 那段注释现在写的是"four different things … so the four stay four"，而函数实际已有 **8** 条 `cause=` 出口＋1 枚常量。这是**注释过期**，不是行为错误；改它＝改票 223 立的定式文字（`SPEC-12 §4.1` 口径下属"动别人钉过的措辞"），留给票 232（它正是"句子存在≠句子带货"那族测试面加固）。⚠ 本腿**未**改这段注释一字。 |
| `DEFERRED(...)` 代码标记 | 本票**没有**加代码标记：AGENTS.md §1.1 要求 `DEFERRED(D-xx)` 与 `SPEC-12 §5` 登记表 **1:1 双向**对得上，而本腿禁改 `docs/specs/**` 与台账 ⇒ 登记落在本件（§4），由编排者决定是否升格为 `A##`/`DEFERRED`。 |

判语：**AC#4 以"只登记不实现"交付**，缺尺 5 条＋不动理由 5 行全部具名到文件与行。

## §5 AC#5 整包终态与逐名红名册

命令逐字＝`PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" GOFLAGS= go test -count=1 ./cmd/wisp ./internal/...`
（⛔ 无 `-race`、⛔ 无 `-cover*`、⛔ 无 `-short`）。跑了**两发**：

| 发 | 取数时刻（+08） | HEAD | rc | 日志 |
|---|---|---|---|---|
| 甲（收尾笔之前） | `2026-10-06 09:09:59 → 09:17:23` | `c85a9b65`（含 `290dca87`＋`a16d1ff7`，不含 `c5040ea7`） | 1 | `logs-ac5-full.txt`（env 件 `logs-ac5-full.env`） |
| 乙（**终态，判语以此发为准**） | `2026-10-06 09:19:58 → 09:26:48` | `c85a9b65`（含本腿三笔，末笔 `c5040ea7`） | 1 | `logs-ac5-final.txt`（env 件 `logs-ac5-final.env`） |

**终态（乙发）逐名红名册**——`FAIL` 包 **2 枚**、红用例 **5 枚**，逐名与票面 AC#5 的历史在册名册**逐字对上、零枚新增**：

| 包 | 红用例（逐名） | 它红在哪一句（本腿现读） | 归属 |
|---|---|---|---|
| `internal/ball` 1 枚 | `TestC21TableColourRowsMatchTokensCSS` | `tokens_table_test.go:1468: read design/assets/tokens.css: … cannot find the path specified - the CSS leg of this check must never skip` | 别人地界（`design/assets/tokens.css` 在工作树里是 ` D` 老脏面，AGENTS.md 派单第 3 条明令本腿不碰 `design/**`） |
| `internal/panel` 4 枚 | `TestApprovalCardViewJSONKeysMatchFrontendTypes`、`TestComposerContractTypesMatchFrontend`、`TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme`、`TestC21DesignTokensFourWayAgree` | 前两枚＝Go 侧 emit 的 JSON 键与 `frontend/**` 的接口声明不符（`Go Snapshot emits [instructions tasks] that interface PanelSnapshot does not declare`）；第三枚＝`frontend/src/components/harness/right-rail.tsx:89` 出现第二处颜色源；第四枚＝同一枚缺失的 `design/assets/tokens.css` | 别人地界（`frontend/**`／`design/**` 本腿零读零写） |

**⛔ 没有一枚红落在本票射程内**：`cmd/wisp`＝**ok 405.208s**（全包最大枚数的一枚套件，含票 223 全族＋本票新那一行）、
`internal/config`＝**ok 5.298s**（`loader.go` 那句原文与它的两枚钉手 `loader_test.go:273`／`migrate_test.go:160` 都在这发里绿）。

**`internal/risk TestResolvePerCallBudget`（争用型假红那一枚）**：甲发（`09:10`）红＝
`pathresolver_budget_norace_test.go:37: C26 budget breach: Resolve averages 1.792657ms per call, budget 1ms`；
乙发（`09:2x`，同机安静下来后）**绿**＝`ok internal/risk 7.318s`（日志里 `TestResolvePerCallBudget` 零命中＝没红）。
按票面口径（"以安静 `-count=3` 为准"）另跑一发独立尺：
`go test -run TestResolvePerCallBudget -count=3 -v ./internal/risk`（取数 `2026-10-06 09:3x +08`，起手 CPU/MEM＝`31.7/73.7`）
＝**3 枚全 PASS**（2.17s／2.33s／1.25s，`ok 5.779s`）⇒ 判定＝争用型假红坐实，**非本票新增**。

**两把防"把没跑读成绿"的尺**（派单末条）：
- `grep -c "SKIP"` 两发都＝**0**；全仓 `[no test files]`＝4 枚（`internal/speech`／`internal/streamkey`／
  `internal/watchdog`／`internal/agent/scheduler`），⛔ 没有一枚是本腿新增，⛔ 本腿未写过一行 `t.Skip`。
- `cmd/wisp` 那发的**包级** `ok` 出现＝测试二进制真跑过（若缺 sherpa DLL 会是 `exit status 0xc0000135` 且零 `--- FAIL`），
  起手 PATH 里带了 `$PWD/third_party/sherpa-onnx:$PWD/build`（现量两目录各含 `onnxruntime.dll`＋`sherpa-onnx-c-api.dll`）。

**放宽断言的自证**＝本票对测试面只做了两件事（`git diff --numstat` 全程可查）：表里 `+14 0`（加注释＋加一行）、
两张名册各 `+4 1`／`+4 0`（marker 进册）；⛔ 未删任何断言、⛔ 未改任何 needle、⛔ 未改 SLO／`thresholds.go`／golden。
`--- FAIL` 枚数在正控那发＝1、其余各发＝0。

判语：**AC#5 成立**——整包到终态、逐名比红名册零新增（别人地界 5 枚逐名对上、risk 那枚以安静 `-count=3` 判为假红），
⛔ 一枚都没顺手修。

## §6 四枚既有句的逐字节零变化尺（⛔ 不拿"用例全绿"代替）

尺＝把**受保护的 12 行组**（`config_reload.go` 的 `:97-99` 常量＋missing/permission/syntax/unknown-key/
migration/invalid 六支 `return`，`loader.go` 的 `:120-124`）用 `sed -n` 抽成单一文件，落笔前后各一次，
比 `diff`＋`sha256sum`：

| 项 | 落笔前 | 落笔后 |
|---|---|---|
| 文件 | `baseline-protected-lines-before.txt`（2,073 字节） | `protected-lines-after.txt`（2,073 字节） |
| sha256 | `4ddddbc44724ed16cd62c18968be7703d818b5614ea904ce61507131c8b5c57b` | **同一枚** `4ddddbc4…c57b` |
| `diff` | — | **0 行**（`DIFF=0 identical (12 protected line groups)`） |

行号位移的对照（证明"抽的是同一段"）：missing `:357-359`／permission `:361-363`／syntax `:370-372`／
unknown-key `:374-376`／migration `:393-395`／常量 `:97-99`／loader `:120-124` **落笔前后同址**；
只有 `cause=invalid` 那三支从 `:397-399` 顺移到 `:425-427`（＋28，正是插入点之后的整体位移）。
再叠一把独立尺＝`git diff --numstat -- cmd/wisp/config_reload.go`＝**`28 0`**（零删除）＋
`git diff -U0` 只有一枚 hunk `@@ -395,0 +396,28 @@` ⇒ 既有句所在行**不可能**被改（改一行会出删除计数）。
`internal/config/loader.go`＝`git diff --numstat` **零输出**（未进 diff）。

判语：**AC#2 的"既有四句一字不改"以逐字节尺自证**（＋那四句今天无测试正向钉，见 §4 登记 #2）。

## §7 门禁四数

| 门 | 读数 | 时刻 |
|---|---|---|
| `sh scripts/d22scan.sh` | 正控 `runtests.sh: OK … PASS=35 FAIL=0 SKIP=0`；实扫 `clean - no D22 ban violations`（bans #1-5 internal=228／cmd=38，#6 frontend=85，#7 tools=23，#8 design=39／frontend=85／internal=512／cmd=104） | `2026-10-05 23:1x +08`（改动后复跑） |
| `sh scripts/check-path-length-budget.sh --with-self-test` | `positive control PASSED`；`denominator: tracked paths=6089 over-budget=57 covered by roster=57 not in roster=0`；`VERDICT GREEN` | 同上 |
| `GOFLAGS= go vet ./cmd/wisp/` | **rc=0**（零输出） | `2026-10-06 09:1x +08`（最终文件态复跑） |
| `"D:/work/base/gopath/bin/gofumpt.exe" -l cmd/wisp` | 只剩 **1 枚**＝`cmd\wisp\models.go`（**预存** CRLF 历史脏枚，⛔ 非本票、本腿未格式化它，具名即可）；`cmd\wisp\config_reload.go` 在加入本腿分支后被列进过一**次**（`:``+` 粘连写法），已按 gofumpt 现形修正 ⇒ 现不列出 | 同上 |

CPU／MEM 现量（每发起腿前）：起手 `31.4/68.0`；写码后 `64.6/75.8` → 等待 45s `61.7/75.9`；
整包前 `39.2/75.5`。⚠ MEM 那一条**全天贴着 70% 以上**＝本机常驻（`Memory Compression` 3.0GB／
`vmmemWSL` 1.5GB／两枚 Qoder CN 各 ~0.9GB／`DeepSeek Harness` 0.9GB／`MsMpEng` 0.5GB，`Get-Process` 现读），
⛔ 不是本腿的在飞测试（`Get-Process go,wisp` 在整包前＝零枚）。⇒ 本腿把 MEM 判为"环境基线、非争用"，
按 CPU 门放行长跑，并把这一处**具名交回**（派单说的"机器内存贴着 70% 帽"就是这一形状）。

## §8 判不动的地方／本腿没做的

1. **票面 AC#5 的"逐名比红名册"里，别人的地界本腿一律未修**（`internal/ball` 1＋`internal/panel` 4
   属历史在册、`internal/risk TestResolvePerCallBudget` 属争用型假红）——读数见 §5（待补）。
2. **没有把 `cause=newer-build` 做成 stdout 行**：票面没写，普查件 §1 的前提（七句全走审计面）成立 ⇒
   新分支照同一条通道，加 stdout 属"票外动作"，未做。
3. **没有把 `:424` 改窄、没有引入机读分类**：见 §4 表格的四条具名理由。
4. **`git status --porcelain -- cmd/wisp internal tools`**：本腿交件时＝**0 行**（两笔改动已 commit），
   工作树里从未出现探针件（`-overlay` 只在内存里生效）；物理件全在本目录与 `mutation/` 下（临时件只建不删）。
5. 本腿**未读**也**未引用** `frontend/**`、`design/**`、`.scratch/wisp/probes/{268,111/r3,evidence-close,pool-validity}`；
   ⛔ 未碰 `.github/workflows/ci.yml`、`scripts/**`（并行腿 `111-r3` 地界，只读执行了两把门禁脚本）、
   `docs/reports/HANDOVER.md`、`docs/reports/pending-and-issues.md`、`docs/evidence/s1/**`、`thresholds.go`、golden、
   `allowlist.txt`、`PLAN.md`、`docs/specs/**`；未 push；仓内零删除。
