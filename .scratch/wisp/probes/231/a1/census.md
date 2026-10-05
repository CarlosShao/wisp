# 票 231 · 只读普查腿 `231-a1` —— 落地腿开工前的撞钉预检（交料，不裁决）

- 腿＝`231-a1`（只读普查）；写点＝本文件一处，⛔ 未动任何代码／测试／工单／台账／票面 `- [ ]` 框。
- ⛔ **本腿一枚 `go` 命令都没跑**（无 `go test`／`go build`／`go vet`／`go env`），尺＝`grep`／`git show`／`git log` 的现读。
- 起手锚：HEAD `b70cf069`；现量时刻 `2026-10-05 15:0x +08`。
- 起手树态（⚠ 不是本腿改的）：`git status --porcelain -- cmd internal tools scripts docs` ＝ **1 行** ` M cmd/wisp/resident_approval_windows.go` ＝在飞写腿 `268-r1` 的地界；本腿交件时这一行必须仍是唯一一行。
- 票面：`.scratch/wisp/issues/231-a-config-written-by-a-newer-build-is-reported-to-the-operator-as-a-validation-failure.md`（34 行／待[填] 字节，5 格 AC 全未勾）。
- 我这轮**不答**的：票 231 AC#1（种一份新 `config.toml` 抄回操作员原句）＝要跑程序 ⇒ 归编排者安静窗口；本件里凡标"生产分类器今天会说哪句"的，都是**读码得到的形状**，⛔ 不是实跑读数。

## 节状态

| 节 | 内容 | 状态 |
|---|---|---|
| §0 | 起手复认：票面锚 vs 现读 | 待[填] |
| §1 | 问一：分类器今天有哪几条出口 | 待[填] |
| §2 | 问二："版本更高"那支的真身（＋对照支） | 待[填] |
| §3 | 问三：钉住今天这三形的用例 | 待[填] |
| §4 | 问四：新增一支会撞谁（四类尺穷尽） | 待[填] |
| §5 | 问五：既有定式对照＋票 232 撞行核查 | 待[填] |
| §6 | 必须交给落地腿的那一页名册 | 待[填] |

---

## §0 起手复认：票面／派单锚 vs 现读

尺＝逐枚 `Read` 现读＋`git show <commit>:<file> | grep -n` 复算历史行号；票面的行号我一律不抄，全部重新认。工单件尺＝`wc -l -c` ＝ **33 行／5,856 字节**（末行无换行，故 Read 报 34 行；5,856 与台账 `A444` 立票时写的"5,856 字节／5 格"**逐字相同** ⇒ 票面自 09-29 起一字未动，下面这些漂移全是**代码侧**动的）。

| # | 锚 | 票面／派单说 | 我现读（HEAD `b70cf069`） | 判 | 漂因（逐枚 `git show` 复算，不是猜） |
|---|---|---|---|---|---|
| 1 | "版本更高"那支 | `internal/config/loader.go:106-110` | **`:120-124`**（`if ver > SchemaVersionCurrent` 在 `:120`，那句英文文本在 `:122`） | **漂 ＋14** | `8e443ed7`（2026-10-04 18:59，257-r1b，`--numstat`＝**14 增 0 删**，它在 `:67-80` 插了 `fileMissing` 那一支）。复算链：`1fd19d4f`:47 → `78988901`:70 → `248095d1`:90 → `5a755c3c`:108 → `8e443ed7`:122＝HEAD:122 |
| 2 | 吞掉真话的那条 case | `cmd/wisp/config_reload.go:357`（票面还写了尺＝`sed -n '353,362p'`） | **`:396`**（`case strings.HasPrefix(d, "config.toml:"):`；它回的 `return` 在 **`:397-399`**） | **漂 ＋39**；票面那条锚在 `5a755c3c` 那发**逐字成立过**（复算＝`:357`），此后过期 | `ae60a87c`（10-03 09:36，票 255 r2 AC#1，`:357`→`:390`＝＋33）＋`67ab595d`（10-03 10:37，255-r2 `restartTierKeys` 那枚测试的注释，`:390`→`:396`＝＋6）。`a4906b6c` 未动这一行（复算＝`:390`） |
| 3 | 钉句的用例 | `cmd/wisp/config_sentences_223r2_test.go:67-69`（派单标"会漂"） | **未漂**：`:66-69` 是一整枚表条目，`:67`＝用例名＋body、`:68`＝三根 needle、`:69`＝收尾 | 行号对上，但**票面对这三行的描述错了**（说它钉 (b)(c)(f) 三形，现读它钉的是 (a) 那一形一枚）⇒ 见 §3 | 该文件自创建起**零枚**后续 commit（`git log --oneline -- <file>` ＝只有 `5a755c3c` 一发） |
| 4 | "迁移那一支"的锚 | 票面写 `:355`；**派单把它读成 `internal/config/loader.go:355`** | ⛔ **越界**：`wc -l internal/config/loader.go`＝**263**（最大可读行 264）⇒ `:355` 不可能是 loader.go。票面引的那句中文（「文件被原样留着，不会被重置…升级 Wisp 或恢复备份才会读它」）逐字在 **`cmd/wisp/config_reload.go:394-395`**（`cause=migration` 那支的 `return`，223-v2 当时写 `:354-356`，同支漂到 `:393-395`） | **派单/票面这一格指错文件**；两支我都抄全：分类器侧 `config_reload.go:377-395`＋配置层侧 `loader.go:125-130`→`internal/config/migrate.go:53-56`／`:60-63` | 同 2（＋39） |
| 5 | AC#2 的"既有四句" | 「既有四句（缺失／语法错／权限不够／热加载被禁用）一字不改」 | 复认：`describeReloadFailure` 今天有 **7 条 `cause=` 出口**（尺＝`grep -c "cause=" cmd/wisp/config_reload.go`＝7）＝missing/permission/syntax/unknown-key/migration/invalid/unclassified；**"热加载被禁用"不在这个函数里**，它是常量 `hotReloadDisabledPanelInbound`（`cmd/wisp/config_reload.go:97-99`）。⇒ 票面那四支之外**还有 3 支（unknown-key／migration／invalid）没被"一字不改"保护**，其中 `invalid` 恰是本票要动的那条 | **半对半错**：四句各有出口＝对；"四句"这个枚数＝错（该写"四句＋三支"） | — |
| 6 | AC#3 的两半 | 「`:67-69` 那三形必须**改成断新出口**，并**保留 (a) 那形**」 | ⛔ **自相矛盾**（现读 §3）：那三行今天**就是 (a)**；按字面改它＝删掉 (a) 唯一的常驻钉。而今天**没有任何一枚常驻用例**种得出"解析得开的未来版"那一形 ⇒ 拿掉新分支**不会红**任何现有行，AC#3 的"正控"今天无处可跑 | **推翻票面**，材料见 §3/§4/§6 | 台账 `A444`（`docs/reports/pending-and-issues.md:9484`）写的是「那三行是我那句裁定的化身，**动它＝动裁点**」＝与票面 AC#3 前半句指同三行、要相反的动作 |
| 7 | 223-v2 的匹配顺序 | 「`config.toml parse` → `unknown key` → `cannot migrate` → `config.toml:`」 | 复认**成立**，现读四行分别在 `:369`／`:373`／`:377-378`／`:396`（顺序就是 switch 的书写顺序，Go 的 `switch { case … }` 无 break 穿透 ⇒ 先命中者得） | 对 | — |
| 8 | 派单里"编排者另外让你做 X" | — | **零枚**：本腿收到的只有这一条派单，正文里没有任何"另写文件／另改代码／另跑门"的指示；我也**未**在任何文件里读到自称"派单增量"的内容。若后续出现 ⇒ 按规矩停手回报，不照做 | 无异常 | — |

**本节结论（一句话）**：票面的**现象**与**匹配顺序**现读全部复认成立；票面的**行号**两枚已漂（#1 漂 14、#2 漂 39）、一枚指错文件（#4）、票面 **AC#3 与 AC#2 的枚数口径各有一处要推翻**（#5、#6）。

## §1 问一：`cmd/wisp/config_reload.go` 的分类器今天有哪几条出口

尺＝`Read cmd/wisp/config_reload.go:342-409` 现读＋`grep -n "cause=" cmd/wisp/config_reload.go`＝**7**。函数＝`describeReloadFailure(err error) string`（`:354`）。它是**外一层 switch ＋ `errors.As` 后内层 switch ＋ 兜底 return** 三段，书写顺序＝命中顺序（Go 的 `switch { case … }` 无穿透），所以**先命中者得句**；票面/223-v2 说的顺序复认成立。

**先记一条全节都吃的前提（具名，防落地腿照错形状写）**：这 7 条出口的字符串**只走审计面，不走 stdout**。发射点是 `:155` `rt.auditf("config: HOT-RELOAD state=not-applied %s", describeReloadFailure(err))`（`reloadOnce` 的错误支在 `:154-157` 直接 `return`，⛔ 没有任何 `Fprintf(rt.stdout, …)`）；`auditf` 的落点＝`cmd/wisp/run.go:931-935`＝**stderr 前缀 `[audit] ` ＋ JSONL sink 一条 `audit: ` 记录**。stdout 那几行（`:187`／`:190-193`／`:201-203`／`:223-225`／`:228-230`／`:234-236`／`:321-326`）全部属于 **D36 三档回执**（`reportReload`／`reportRestartPending`），与"读不到"这七句**不是同一条路**。票 231 AC#2 要的"出口"按现形＝**再给一条 `cause=` 审计句**；要不要另配一行 stdout 是**新决策**，票面没写（本腿不替它定）。

| # | 出口 | 靠什么形状认领（★＝机读形状，☆＝纯词面） | 它回的原句（逐字抄，含 `cause=` 前缀） | `file:line` |
|---|---|---|---|---|
| 1 | `cause=missing` | ★ `errors.Is(err, fs.ErrNotExist)`（sentinel，走 `*observe.Error.Unwrap`，`internal/observe/errors.go:207-212`） | `cause=missing detail="config.toml 读不到：文件不存在（这一条只说缺失，不说语法、不说权限）。本次运行继续用内存里的旧配置；文件回来之后要它的 mtime 或大小变过才会被下一 tick 重读"` | 认领 `cmd/wisp/config_reload.go:356`，原句 `:357-359` |
| 2 | `cause=permission` | ★ `errors.Is(err, fs.ErrPermission)` | `cause=permission detail="config.toml 读不到：这个进程没有读它的权限（文件在，也读得开名字，只是不让读）。本次运行继续用内存里的旧配置；这一条不是语法错，也不是文件缺失"` | 认领 `:360`，原句 `:361-363` |
| — | （进入内层的闸门） | ☆ `errors.As(err, &oe)` 取 `*observe.Error`，然后 `d := oe.Detail`——⚠ **分类码 `oe.Class` 全程不参与归句**（`observe.ClassConfig` 被 `internal/config` 的 85 枚错误点共用，尺＝`grep -rn "observe.New(\|observe.Wrap(" internal/config/ --include=*.go \| grep -v _test.go \| wc -l`＝85） | — | `:365-367` |
| 3 | `cause=syntax` | ☆ `strings.HasPrefix(d, "config.toml parse")` | `cause=syntax detail="config.toml 读到了但解析不了：这一行不是合法 TOML 语法（不是权限、不是缺失）。本次运行继续用内存里的旧配置；修好之后要再出现一次新的 mtime/大小才会被重读"` | 认领 `:369`，原句 `:370-372` |
| 4 | `cause=unknown-key` | ☆ `strings.Contains(d, "unknown key")` | `cause=unknown-key detail="config.toml 语法没问题，但里面有这份 schema 不认的键（拼错的键会被这样拒绝，而不是被忽略）。本次运行继续用内存里的旧配置"` | 认领 `:373`，原句 `:374-376` |
| 5 | `cause=migration` | ☆ `strings.Contains(d, "cannot migrate from schema version") \|\| strings.Contains(d, "no migration registered")`（两串或） | `cause=migration detail="config.toml 声明了一个这份 Wisp 不会迁移的 schema_version（文件被原样留着，不会被重置）。本次运行继续用内存里的旧配置；升级 Wisp 或恢复备份才会读它"` | 认领 `:377-378`，原句 `:393-395`（注释 `:379-392` 是 223-r2 那段边界说明） |
| 6 | **`cause=invalid`＝本票的吞点** | ☆ `strings.HasPrefix(d, "config.toml:")`（**词面前缀，既当身份又当归类**） | `cause=invalid detail="config.toml 语法没问题，但内容被校验拒绝（值不合法或引用解不开）。本次运行继续用内存里的旧配置"` | 认领 **:396**（票面说 `:357`，漂 ＋39，见 §0 表 #2），原句 **:397-399** |
| 7 | `cause=unclassified`（兜底，非分支） | 什么都不认领：把 `err.Error()` **原文拼进句尾**（"引原因而不是假装归过类"） | `cause=unclassified detail="config.toml 没有被重读成功，原因没有归入已知四类（缺失/语法/权限/schema）；本次运行继续用内存里的旧配置: " + err.Error()` | `:406-408` |
| 8 | `state=disabled`（⛔ 不在这个函数里） | ☆ 常量整串原样印 | `config: HOT-RELOAD state=disabled host=panel-inbound detail="本宿主没有接 config.toml 轮询，这次运行期间手改配置不会生效；它也没有可弹卡的面，所以放宽本来就无法被确认（fail-closed）。要生效请重启进程并用 wisp run."`（源码里 `detail=` 用的是转义引号） | 常量 `:97-99`，发射点在 `cmd/wisp/panel_inbound*.go`（本腿不展开） |

**票面那条断言（"靠 `config.toml:` 前缀认领校验拒绝"）＝成立，行号现读 `:396`**；它认领的这枚前缀**今天身兼 6 类互不相同的事实**，因为 `internal/config` 非测源码里以 `"config.toml: ` 开头的 detail 共 **48 枚**（尺＝`grep -rn '"config\.toml: ' internal/config/ --include=*.go \| grep -v _test.go \| wc -l`），逐文件＝loader 3／parse 8／validate 21／migrate 4／catalog 11／unwired 1，且这六枚文件**全部在 `readConfigFile` 的可达路径上**（`validate.go:38` 调 `validateCatalog(c)`⇒catalog 那 11 枚也可达）。其中真正会被 `cause=invalid` 认领的读路径形状至少含：

1. **版本更高**（`loader.go:122`）＝本票要给的出口——⛔ 今天它零出口，且那句"值不合法或引用解不开"对它**假**（根本没走到 `decodeStrict`/`validate`）。
2. `api_key_ref` 解不开（`loader.go:233-234`／`:241-242`）＝**这才是"引用解不开"的真身**。
3. `decodeStrict` 的类型/位置错（`parse.go:154-155` `config.toml: line %d, col %d: …`）＝**"值不合法"的真身**。
4. `[plugins]` 形状错（`parse.go:100/106/110/115/120/126`）。
5. **迁移成功但产物不合法**（`migrate.go:71-73`／`:78-80`，串＝`config.toml: migration to schema version %d produced an invalid config (%v); the file was left untouched`）——它**不含** `cannot migrate from schema version`／`no migration registered` 任一串 ⇒ **今天也被 `cause=invalid` 吞**，与"版本更高"同族、同一根前缀。**票 231 字面没要求处理它**（本腿只具名，不裁决）。
6. **写了但不起作用**（`unwired.go:106` `config.toml: %s is written but does nothing: …`）。

⚠ **`unknown key` 那枚为什么今天没被前缀抢走**：`parse.go:148` 的 detail **同时**满足"以 `config.toml: ` 开头"与"含 `unknown key`"，靠的是**内层 switch 里 `:373` 写在 `:396` 之前**。⇒ 落地腿新增一支时，**位置**（必须早于 `:396`）与**串不重叠**两件事任缺其一就会被吞回去。

## §2 问二：`internal/config/loader.go` 的"版本更高"那支的真身

### 2.1 那支本体（现读 `:120-124`，票面说 `:106-110` ⇒ 漂 ＋14，漂因见 §0 表 #1）

```go
	if ver > SchemaVersionCurrent {                                     // :120
		return nil, observe.New(observe.ClassConfig, fmt.Sprintf(       // :121
			"config.toml: schema_version %d was written by a newer build (this build understands %d); upgrade Wisp or restore a backup",   // :122
			ver, SchemaVersionCurrent))                                 // :123
	}                                                                   // :124
```

- **产给日志的原句（逐字，含格式动词）**：`config.toml: schema_version %d was written by a newer build (this build understands %d); upgrade Wisp or restore a backup`。`SchemaVersionCurrent`＝**2**（`internal/config/schema.go:28`）⇒ 种 `schema_version = 99` 时渲染出的 detail＝`config.toml: schema_version 99 was written by a newer build (this build understands 2); upgrade Wisp or restore a backup`，`err.Error()` 再过一层 `internal/observe/errors.go:188-204` 变成 `config: config.toml: schema_version 99 was written by a newer build …`。⚠ 这两串我**不是实跑读的**，是从 `:122`＋`errors.go:192-203` 的拼装规则推的；`docs/evidence/s1/223-hot-reload-wiring-v2.md:188` 的原始件读数与此**逐字相同**（那是非实现者腿实跑的，可当对照）。
- **类型／分类码／cause 链**（这决定落地腿能不能不靠字符串分它）：

| 维度 | 现读 | 后果 |
|---|---|---|
| Go 类型 | `*observe.Error`（`internal/observe/errors.go:170-176`） | 分类器 `errors.As(err, &oe)` 取得到（`config_reload.go:365-366`） |
| 构造子 | `observe.New(class, detail)`（`errors.go:179-181`） | ⛔ **`Err` 字段没赋值 ⇒ cause 链为空**（与 `observe.Wrap` 的区别就在这一条；`Unwrap` 返回 nil，`errors.go:207-212`） |
| 分类码 | `observe.ClassConfig` ＝ `"config"`（`errors.go:19`） | ⛔ **不可用于归句**：`internal/config` 非测源码里 85 枚错误点全用这一枚（尺见 §1 表内） |
| `ProviderCode` | `""`（未赋值） | **全仓产码从未给 `ProviderCode:` 赋过值**（尺＝`grep -rn "ProviderCode:" cmd/ internal/ --include=*.go \| grep -v _test.go \| wc -l`＝**0**；`internal/config` 内亦 0）⇒ 它是**现成空闲的机读位，且是已导出字段**（`errors.go:172`），填它**不新增导出名**（票面禁区＝"不新增导出名／要新增先落 `A##`"） |
| `RetryAfter` | `0`；`RetryPolicy(ClassConfig)`＝`RetryNever`（`errors.go:114-117`） | 与本票无关，仅记下别再打它的主意 |

⇒ **今天唯一的判别材料就是 detail 字符串本身**（这正是 bug 的形状：前缀既当身份又当归类）。落地腿若想"用机读形状而不靠字符串"，`ProviderCode` 是本腿找到的**唯一一枚不新增导出名的现成位**；⚠ **本腿具名一处冲突交回**：`ProviderCode` 一旦非空，`Error()` 会把句首渲染成 `config (<code>): …`（`errors.go:194-196`）⇒ 分类器读的 `oe.Detail` 不受影响（AC#2 那句"不许改 loader 原文"仍守得住），但**任何按 `err.Error()` 全文比对的既有断言**形状会变（本腿在 `internal/config` 侧只找到 `Contains` 型尺，见 2.4，未见全等尺）。本腿**不裁决**用不用它——AC#4 明写"这一步允许只登记不实现"。

- **今天怎么被吞的（复现链，逐行现读）**：`schema_version = 99` ＋ 正文解析得开 ⇒ `peekSchemaVersion` 返回 `(99, nil)` ⇒ `:84`（`ver == 0 && peekErr != nil`）不进、`:104`（`peekErr != nil && ver >= SchemaVersionCurrent`）不进 ⇒ **`:120` 命中** ⇒ detail 以 `config.toml: ` 开头 ⇒ `config_reload.go:369/373/377` 三支全不命中 ⇒ **`:396` 命中** ⇒ 操作员得到 `cause=invalid`「语法没问题，但内容被校验拒绝（值不合法或引用解不开）」。
- ⚠ **种子形状必须"解析得开"**：若正文坏（如 `schema_version = 99\nbroken [[[\n`），`:104` 抢在 `:120` 之前返回 `config.toml parse` ⇒ 归 `cause=syntax`＝票面 AC#3 要保留的 (a) 形。**这条短路是 loader 的书写顺序给的，不是分类器给的**——编排者跑 AC#1 时按这个形状种。

### 2.2 对照支："版本更低＋不会迁移"——它有出口，且出口里带行动指引

- 路由：`loader.go:125-130`（`if ver != SchemaVersionCurrent { raw, err = applyMigrations(path, raw, ver) … }`）→ 错误产在 `internal/config/migrate.go`：
  - `:53-56`：`config.toml: no migration registered from schema version %d (this build understands up to %d); the file was left untouched - fix or restore it manually, it will never be silently reset`
  - `:60-63`：`config.toml: cannot migrate from schema version %d: %v; the file was left untouched - fix or restore it manually, it will never be silently reset`
- 分类器侧的专属出口＝**`config_reload.go:377-378` 认领 → `:393-395` 原句**（逐字见 §1 表 #5）：`config.toml 声明了一个这份 Wisp 不会迁移的 schema_version（文件被原样留着，不会被重置）。本次运行继续用内存里的旧配置；升级 Wisp 或恢复备份才会读它`。
- ⚠ **派单 Q2 说的"`internal/config/loader.go:355` 附近"这枚锚指错了文件**（loader.go 只有 263/264 行）：那句中文真身在 `cmd/wisp/config_reload.go:394-395`（223-v2 当时记 `:354-356`）。⇒ "有出口／没出口"的对照**两边都在分类器里**，配置层只是给它喂串。
- 钉住这条对照的既有断言＝`internal/config/migrate_test.go:123 TestMigrateCorruptFileUntouched`（尺＝`strings.Contains(err.Error(), "migrat")`，`:131`），票 223 与本票都写死它一字不许动。**本腿独立复认**：`git log --oneline -3 -- internal/config/migrate_test.go` 最新＝`a95ee3a8`（票 83 那批）⇒ 与票 223 件里那句"`migrate_test.go` 最后一次被改是 `a95ee3a8`"对上，`:123` 那枚断言今天确实没被动过。

### 2.3 "版本更高"在操作员面上的出口枚数＝**零**；配置层原句的常驻钉＝两枚

- 分类器里：**0 条** `cause=` 与它对应（§1 那 7 条逐条排除）。
- 配置层两枚钉住那句英文的常驻用例（⛔ 都不许被本票改动，AC#2 已写明不改 loader 原文）：
  - `internal/config/loader_test.go:267-276 TestLoadFileNewerSchemaVersionRejected`——种 `"schema_version = 99\n"`，needle＝`strings.Contains(err.Error(), "99")`（`:273`）。
  - `internal/config/migrate_test.go:155-163 TestMigrateHigherIntermediateVersionRejected`——同种子，needle＝`strings.Contains(err.Error(), "newer build")`（`:160`）。
  ⇒ 这两枚**只钉 `err.Error()`，钉不到操作员那一句**；这正是票面说的"那句真话在代码里活着"的物证。

### 2.4 本节判语

现象、匹配顺序、"零出口"、"同一句中文只在 migration 支里齐全"——**四条全部现读复认成立**。票面/派单的两枚行号锚（`loader.go:106-110`、`loader.go:355`）一枚漂 14 行、一枚**指错文件**；`ProviderCode` 那条机读路存在但带一处渲染副作用，**具名交回不裁决**。

## §3 问三：`cmd/wisp/config_sentences_223r2_test.go` 今天钉着哪几形

文件尺＝`wc -l` **124 行**；枚数尺＝`grep -c "^func Test"`＝**1**（唯一用例 `TestTicket223R2FailureSentenceRouting`，`:39`），表里 **8 枚子形**（尺＝数 `{` 条目，`:50`–`:86`，逐枚现读；子用例全名＝`TestTicket223R2FailureSentenceRouting/<表里的中文名>`）。该文件自创建起**零枚后续 commit**（`git log --oneline -- <file>`＝只有 `5a755c3c`）⇒ 票面说的 `:67-69` **未漂**。

### 3.1 三把尺（`:110`／`:114`／`:118`，逐字）

| 尺 | 代码（现读） | 语义 |
|---|---|---|
| ①前缀尺 | `:110` `if !strings.HasPrefix(line, tc.wantCause) {` → 红句 `:111-112` `"shape %s is booked %q, want it to open with %q - the two sentences swapped again"` | 操作员句**必须以 `cause=xxx` 开头** ⇒ ⚠ **新增那一支的返回串也必须以 `cause=` 起头**，否则这里红 |
| ②反向 cause 尺 | `:114` `if strings.Contains(line, tc.notCause) {` → 红句 `:115-116` `"shape %s borrows %q, which is a different pipeline's answer:\n%s"` | 不许借用另一支的 `cause=` |
| ③反向半句尺 | `:118` `if tc.notFragment != "" && strings.Contains(line, tc.notFragment) {` → 红句 `:119-120` `"shape %s says %q, which contradicts its own cause:\n%s"` | 不许说出与自己归句矛盾的那半句中文 |

读数路径＝生产分类器本体（`:104` `mgr.CheckAndReload()` → `:108` `describeReloadFailure(err)`），种子装配＝`:92` `config.SaveFile(path, config.NewDefaults())` → `:95` `config.NewManager(path, nil)` → `:99` `Sleep(30ms)` → `:100` `os.WriteFile(path, tc.body, 0o600)` → `:103` `Sleep(30ms)` → 重读。⛔ 全表**没有一行种得出"解析得开的未来版"**。

### 3.2 八枚子形逐枚（needle 逐字抄回，行号现读）

| 行 | 子形名（＝表里的 `name`） | 种的 body | `wantCause` | `notCause` | `notFragment` |
|---|---|---|---|---|---|
| `:51-53` | `声明当前版_注释以方括号开头_C1` | `"# [fs] 这不是表头\nschema_version = 2\nbroken [[[\n"` | `cause=syntax` | `cause=invalid` | `语法没问题` |
| `:55-56` | `声明当前版_CRLF_E1` | `"schema_version = 2\r\nbroken [[[\r\n"` | `cause=syntax` | `cause=invalid` | `语法没问题` |
| `:59-60` | `声明当前版_无空格_G1` | `"schema_version=2\nbroken [[[\n"` | `cause=syntax` | `cause=invalid` | `语法没问题` |
| `:63-64` | `声明当前版_缩进版本行_L1` | `"   schema_version = 2\nbroken [[[\n"` | `cause=syntax` | `cause=invalid` | `语法没问题` |
| **`:66-69`** | **`声明未来版_正文语法坏_J1`** | **`"schema_version = 99\nbroken [[[\n"`** | **`cause=syntax`** | **`cause=invalid`** | **`语法没问题`** |
| `:71-72` | `读不出版本_A1_基线不动` | `"this is not toml [[[\n"` | `cause=syntax` | `cause=invalid` | `语法没问题` |
| `:78-79` | `声明旧版_坏表头_A2_仍归迁移` | `"schema_version = 1\n[llm\nbroken ===\n"` | `cause=migration` | `cause=syntax` | `""` |
| `:84-85` | `解析得开_未知键_不抢语法错` | `"schema_version = 2\n\nthis_key_does_not_exist = 1\n"` | `cause=unknown-key` | `cause=syntax` | `""` |

⇒ **今天这张表里 `cause=invalid` 只作为"不许成为"出现（6 次，全是 `notCause`），一次都没作为"必须是"被钉**；`语法没问题` 这半句中文在 `cmd/wisp` 测试侧的全部 7 枚命中（尺＝`grep -rn "语法没问题" cmd/ --include=*_test.go | wc -l`＝7）**全在本文件**，且**全是否定尺**（`:14` 是注释）。

### 3.3 票面 AC#3 复认：前半推翻、后半成立

- **后半"保留 (a) 那形（未来版＋正文坏 ⇒ 仍归语法错）"＝真被这样断着**，物证就是 `:66-69` 这一枚（①②③三把尺双向齐全），223-v2 件 `:208` 也记了"本腿 ②2.3 那发把它打回改前形状时 J1 恰好转红 ⇒ 这颗钉有牙"。⇒ 票面这半句**复认成立**。
- **前半"`:67-69` 那三形必须改成断新出口"＝推翻**：那三行今天只有**一枚**子形（J1＝(a)），不是三形；(b)(c)(f) 那三形（现读未来版＋解析得开）在常驻表里**根本不存在**，它们只活在 223-v2 的 overlay 台件里（`docs/evidence/s1/223-hot-reload-wiring-v2.md:177` 具名：那枚件"**只 `t.Logf` 读数、零判等**"，物理件 `.scratch/wisp/probes/223/v2/overlay/zz_v2probe_cmdwisp_test.go`）。223-v2 自己给的改法（`:215`）是"**并在 `config_sentences_223r2_test.go` 的表里加一行**"＝**加行**，不是**改那三行**。
- ⇒ **连带后果（本条是 §4/§6 的承重）**：按票面 AC#3 的"判据形状＝改完跑该用例全绿；**把 AC#2 那分支拿掉必须红**"——今天**拿掉新分支不会红任何一枚现有子形**（8 枚里没有任何一枚种得出那一形），正控**必须先新增一行才谈得上红不红**。票面 现量表第 6 行"要修这句就得同时动那一处断言"因此也**只说对了一半**：动它≠红，**加行**才红。

### 3.4 要修那句，必须同时改／必须不改哪几行（逐枚具名）

| 处置 | 行 | 为什么 |
|---|---|---|
| ⛔ **一字不许动** | `config_sentences_223r2_test.go:66-69`（J1） | (a) 形唯一的常驻钉；改了它＝票面 AC#3 后半自己违自己 |
| ⛔ 不许动 | `:78-79`（A2，`cause=migration` 唯一钉）／`:83-85`（未知键唯一钉） | 新分支不该抢这两形 |
| ✅ **必须新增**（否则 AC#3 的正控无处跑） | `:50`–`:86` 那张表里加**一行**，body 要"解析得开的未来版"（如 `"schema_version = 99\n\n[ball]\nsize = 64\n"`），`wantCause`＝新出口、`notCause`＝`cause=invalid`、`notFragment`＝`语法没问题` | 只有它种得到 `loader.go:120`；子形名建议沿用表里的中文风格 |
| ⚠ 建议同批（不改不会红，但牙会变钝） | `config_reload_223_test.go:570-573` 的 `all` 名册（现 8 枚）、`config_reload_perm_223_windows_test.go:112-115` 的名册（现 6 枚） | 这两枚是"各句互斥"的**负钉名册**：新 `cause=` 不进名册 ⇒ 新那一形将来误吞别的形状时没人响 |
| ⛔ 产品侧唯一位置约束 | 新 `case` 必须写在 `cmd/wisp/config_reload.go:396` **之前** | 写在它后面＝死代码（前缀先命中），新增行会红在①那把前缀尺上、红句是 `the two sentences swapped again` |

## §4 问四：新增一支会撞谁（`cmd/wisp/**_test.go` 四类尺穷尽）

**尺法**：目录＝`cmd/wisp`，测试文件枚数＝`ls cmd/wisp/*_test.go | wc -l`＝**65**；四类各一把 `grep -rn … cmd/wisp/ --include=*_test.go | wc -l`（⛔ 不拿 `head` 的样例当枚数）。判定口径＝"新增一支 `case` ＋ 新增一行表条目"之后，**这一枚会不会红**。

### 4.1 ①类：对 stdout／stderr 文案的 contains 断言（命中 9 枚文件，逐枚具名）

| `file:line` | 断的原句／needle（逐字） | 加一支会不会红它 |
|---|---|---|
| `cmd/wisp/config_sentences_223r2_test.go:110`／`:114`／`:118` × 8 子形 | 三把尺见 §3.1；红句原文 `the two sentences swapped again`／`which is a different pipeline's answer`／`which contradicts its own cause` | ⛔ **不红**（8 枚子形没有一枚种得出"解析得开的未来版"，`loader.go:104` 的短路把 (a) 挡在 syntax 支）；✅ 新增那一行**受这三把尺管**（必须以 `cause=` 开头、不得含 `cause=invalid`） |
| `cmd/wisp/config_reload_223_test.go:581` | `"config: HOT-RELOAD state=not-applied " + tc.wantCause`（四子形：missing／syntax／migration／unknown-key，种的 body 逐枚见 `:528-566`） | ⛔ 不红（四枚 body 全 `ver <= 2` ⇒ 新支不命中） |
| `cmd/wisp/config_reload_223_test.go:570-573` ＋ `:582-589` | `all := []string{"cause=missing", "cause=syntax", "cause=unknown-key", "cause=invalid", "cause=permission", "cause=unclassified", "cause=migration", "state=disabled"}` ＝**8 枚互斥名册**，负钉跑在 **整条 trail** 上（`awaitAudit` 返回 `r.h.err.String()` 全文，`:116`） | ⛔ 不红（新 `cause=` 不在名册里，缺席无人管）；⚠ **变钝**：新那一形不受互斥保护 ⇒ §3.4 建议同批进名册；⚠ **命名禁忌见 4.4-B** |
| `cmd/wisp/config_reload_223_test.go:249/253/263-264/280/284/311/314/346-347/350-351/357/362/379-380/383-384/387/416-417/421/430/433/453/460/478-479/482/490-494/499/504` | D36 三档与卡面文案（`配置热加载已接管`／`值已换进本进程内存`／负钉 `这些段已立即生效`／`[确认 L2 config.reload]`／`fs.allowed_dirs`／`result=allow`／`effect=applied-after-L2`／`放宽已经过 L2 重新确认`／`仍按启动时建好的 C26 名单`／`result=deny`／`effect=kept-old-values`／`放宽本次没有生效`／`direction=tighten`＋负钉 `config: D36-CONFIRM`／`state=denied`／`risk.permission_mode`／`state=restart-pending sections=[app]`／`effect=next-process-start`／`config: RESTART-PENDING detail=`／`本次运行不会生效`／needle 组 `app.autostart`·`开机自启`·`重启进程后生效`／负钉 `这些段已立即生效：[app]`） | ⛔ 不红——这些支的入口是 `rep`（读成功之后的分档），本票那一支的入口是 `err != nil`（`reloadOnce:153-157`），两条路互斥 |
| `cmd/wisp/config_reload_223_test.go:614` ＋ `:617-620` | `"config: HOT-RELOAD state=disabled host=panel-inbound"` ＋ 4 枚负钉（`cause=missing`／`cause=syntax`／`cause=permission`／`state=armed`） | ⛔ 不红（panel-inbound 不 arm tick，`describeReloadFailure` 在那条腿上根本不被调） |
| `cmd/wisp/config_reload_perm_223_windows_test.go:122`／`:112-119` | 正钉 `没有读它的权限`；6 枚互斥名册（`cause=missing`／`syntax`／`unknown-key`／`invalid`／`unclassified`／`migration`） | ⛔ 不红（种的是 ver=2 合法体＋真 ACL）；⚠ **暗坑 C**：`:142-143` 把 `state=not-applied cause=invalid` 当**控制流**信号用（见 4.4-C） |
| `cmd/wisp/config_receipt_255_test.go:402/408/435/470/473/487/509/515/545/566/600/608` | `state=armed`／`state=applied`／`config: HOT-RELOAD-READER section=panel|session|llm|voice` | ⛔ 不红（种的改动全读得开） |
| `cmd/wisp/firstrun_198_test.go:240`＋`:243-251` | 正钉 `新建默认配置`；负钉 `cause=missing`／`config.toml 读不到：文件不存在`／`本次运行继续用内存里的旧配置`（红句 `the first-run receipt reuses hot-reload wording %q (票 223 keeps the four causes four sentences)`） | ⛔ 不红；⚠ 这枚负钉**只禁首启回执借词**，不禁分类器继续用 `本次运行继续用内存里的旧配置` ⇒ 新句子沿用同一尾巴是安全的、也是一致的 |
| `cmd/wisp/firstrun_257_test.go:199-200`／`:316-338`／`:342-352`／`:356-359` | 三枚 reason×remedy（`替你办完`·`没有任何旧配置可言`／`改不了服务商与模型的存在性`·`手加上面那三样`／`过不了这份 schema 的校验`·`文件一个字节都没动`）；折叠钉 `配置未生效`·`重启就好` 各 `!= 1` 即红；借词负钉 `cause=missing`·`本次运行继续用内存里的旧配置` | ⛔ 不红（这是 settings 写入回执那条链，不经 `describeReloadFailure`）；⚠ **命名提醒**：`配置未生效` 这半句在本仓是**被禁的折叠语** ⇒ 新句子不许出现它 |
| `cmd/wisp/resident_approval_risk_268_windows_test.go:200/205/207/216` | `strings.Count(out, "wisp: resident [risk]:") != 1` 即红；`contains` loader 的答复与 `DefaultApprovalTimeout=300s`；missing 形⛔不得带该 marker | ⛔ 不红（走 `newResidentApprovalWithConfig`，不经分类器）；但见 4.4-E（它把 `config_reload.go:396` 写成了"反面教材"） |

### 4.2 ②类：对句子**枚数**的计数尺

| `file:line` | 尺 | 加一支会不会红 |
|---|---|---|
| `cmd/wisp/firstrun_257_test.go:199-200` | `if seen != 1 { t.Errorf("%s: refusal names %d of the three reasons, want exactly 1 …") }` | ⛔ 不红（settings 拒绝路） |
| `cmd/wisp/firstrun_257_test.go:327-329` | `if len(hits) != 1 { "AC#2 RED: %q appears on %d receipt lines, want exactly 1 (three reasons, three sentences)" }` | ⛔ 不红（同上，且分母是 receipt 不是 trail） |
| `cmd/wisp/firstrun_257_test.go:342-346` | `strings.Count(receipt, "配置未生效"／"重启就好") != 1` | ⛔ 不红；⚠ 新句子若写进 receipt 才会撞上（本票不碰 receipt） |
| `cmd/wisp/firstrun_198_test.go:204/212/220` | `strings.Count(text, "\n["+head+"]\n") != 1` 等三枚（文件正文尺） | ⛔ 不红 |
| `cmd/wisp/config_reload_223_test.go:288-290/341-343/427-429/507-509` | `r.rt.windowCount()` 卡枚数（0／1） | ⛔ 不红——**但**：读失败那一支今天**不发卡**，若落地腿手滑让新形走 `ConfirmLocked`，会红在 `:341`／`:427` |
| `cmd/wisp/resident_ball_228_windows_test.go:66-67`／`cmd/wisp/resident_sink_nail_127_windows_test.go:493`／`cmd/wisp/early_log_nail_130_windows_test.go:257/266`／`cmd/wisp/resident_approval_risk_268_windows_test.go:200` | 各枚 `strings.Count(…, needle) != 1` | ⛔ 不红（全是 resident/child 进程面，与分类器不同一条链） |

**②类的空档（具名交回）**：全仓**没有任何一枚尺数得住"同一个失败被印了两条审计行"**——`describeReloadFailure` 的返回只被 `:155` 一枚 `auditf` 消费，`all` 名册只查"缺席的 marker"，不查"重复的行"。⇒ 若新支与 `:396` 同时命中（顺序写错之外还有一种：新支认领串与旧支重叠且写在旧支**之后**），操作员会得到**两条**或**零条**，**门不响**。这条是 AC#4 该登记的"缺尺"，本腿只具名。

### 4.3 ③类：对 `config.toml:` 前缀的词面尺

尺＝`grep -rn '"config\.toml:' cmd/wisp/ --include=*_test.go | wc -l`＝**0**；`grep -rn "HasPrefix" cmd/wisp/ --include=*_test.go | wc -l`＝**43**（其中与归句有关的只有一枚：`config_sentences_223r2_test.go:110`，它前缀钉的是 `cause=` 不是 `config.toml:`）。
⇒ **测试面没有一枚把 `config.toml:` 当词面 needle**：这个前缀**只活在产码 `:396` 和两处注释**（`resident_approval_risk_268_windows_test.go:37`、`:281`）。⇒ 落地腿收窄/改名 `:396` 的认领形状 ⇒ **不会因"词面尺"红**；会红的只有 4.4-C 那枚控制流。

### 4.4 ④类：日志行数／位置钉 ＋ 四处暗坑

| 暗坑 | `file:line` | 内容与判 |
|---|---|---|
| ④ 位置钉 | `cmd/wisp/resident_sink_nail_127_windows_test.go:608`（record 0 必须是 early resolver 句）／`:619`（`installIdx != 1` 即红）／`:629`（shutdown trail 非空）／`:633-635`（**sink 最后一条**必须是 shutdown 步） | ⛔ 不红——但**前提本腿查出来了**：`startConfigReload()` 的调用点在 `cmd/wisp/run.go:813`，而那行位于 `assembleRuntime`（`run.go:382` 起），**resident 腿共用这枚装配根**（`cmd/wisp/resident_task_source_windows.go:278`）⇒ 常驻腿**今天也 armed 了 tick**。这些用例种的都是 canonical config（不失败）⇒ 新支不发句 ⇒ 不红。⚠ 一旦落地腿给新支**另加一行 stdout**（票面没要求），resident 面的 stdout 捕获会变宽——`268:200`/`228:66-67` 只数各自 marker 故仍不红，但这是"票外动作"，本腿具名不背书 |
| 暗坑 A（位置） | 产码 `cmd/wisp/config_reload.go:396` | 新 `case` 必须早于 `:396`，否则**死代码**；红法＝新增表条目红在 `:110` 那把前缀尺，红句 `the two sentences swapped again` |
| 暗坑 B（命名） | `config_reload_223_test.go:570-573`＋`config_reload_perm_223_windows_test.go:112-119` | 新 `cause=` 串**不得包含**任一枚既有 marker 作子串（`cause=invalid`／`cause=missing`／`cause=syntax`／`cause=unknown-key`／`cause=permission`／`cause=migration`／`cause=unclassified`／`state=disabled`）。反例：`cause=invalid-version` 会让新那一形自己撞 ②反向尺；`cause=newer-build`／`cause=newer-version` 两枚候选都安全。⚠ **两名不一致具名**：票面 AC#2 写 `cause=newer-build`，223-v2 件 `:215` 建议 `cause=newer-version` ⇒ 择一归编排者/落地腿，本腿不裁 |
| 暗坑 C（控制流借词） | `config_reload_perm_223_windows_test.go:142-143` | `if strings.Contains(got, "config: HOT-RELOAD state=applied") \|\| strings.Contains(got, "state=not-applied cause=invalid") { return "", true }`——`cause=invalid` 在这里**不是断言，是"tick 读到了文件"的探测器**。它种的 body 是 ver=2 且解析得开（`:81-82`），靠 `validate`/catalog 那批 `config.toml: ` 串落到 invalid ⇒ ⚠ **只要新支的认领串宽到把那种 body 抢走**（例如误用 `strings.Contains(d, "config.toml: ")` 之类），这枚用例会以 `three plants all landed in the tick's read window; this case cannot be decided on this machine` 或 `neither the permission sentence nor an adoption arrived; stderr:…` 红。**本腿未实跑，这是形状级风险，具名交回落地腿自查** |
| 暗坑 D（同族第二吞点无人钉） | 产码 `internal/config/migrate.go:71-73`／`:78-80` | 尺＝`grep -rn "produced an invalid config" cmd/ internal/ --include=*.go`——命中**全在产码**（`internal/config/migrate.go`），测试侧另有一枚 `internal/config/unwired_test.go` 只钉 "unknown key" 那条（`:254`）。⇒ "迁移产物不合法"今天也被 `:396` 吞成 `cause=invalid`，与"版本更高"同一根前缀、**同样零常驻钉**；票 231 字面没要求处理它 ⇒ 归 AC#4 的登记材料 |
| 暗坑 E（治理口径已被人立过） | `cmd/wisp/resident_approval_risk_268_windows_test.go:32-38`＋AST 尺 `:224-293` | 268 件头逐字：`the rejected alternative (config_reload.go's strings.HasPrefix on the detail string) would redden it`；红句 `:281`：`matching an error message is the shape config_reload.go:396 already carries and ticket 268 was told not to copy`。⚠ **射程本腿查清了**：那枚 AST 尺只 `parser.ParseFile` 一个文件（`:225`＝`resident_approval_windows.go`）且只量 `residentRiskGateValues`（`:233`）⇒ **在 `describeReloadFailure` 里加 `strings.*` 不会红它**；但它是本仓**已记录在案的"不要靠 prose 分类"取向**，AC#4 登记时应引它 |

### 4.5 本节判语

**穷尽读数：加一支"新出口＋新表条目"在今天的 `cmd/wisp` 测试面不会红任何一枚**（①②③④四类逐枚判完，红＝0 枚）。真正会红的是**写歪的三种方式**：写晚于 `:396`（暗坑 A）、命名含既有 marker 子串（暗坑 B）、认领串抢走 perm 用例的 body（暗坑 C）。三处"不会红但会变钝／无人管"＝互斥名册没进新形、重复审计行无数得住、同族第二吞点零常驻钉。

## §5 问五：既有定式对照＋票 232 撞行核查

未判。

## §6 必须交给落地腿的那一页名册

未判。
