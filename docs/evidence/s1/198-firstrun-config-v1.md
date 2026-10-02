# 票 198 · 对抗验收腿 `198-v1` · 裁决表（首建 config.toml 只挂 `wisp run`）

> 本件是**验收腿**（非实现者）的裁决。被测对象＝产码 `613606c0`（`cmd/wisp/run.go` ＋15/0 一处插入 ＋ 新增
> `cmd/wisp/firstrun.go` 98 行／`firstrun_198_test.go` 283 行／`firstrun_acl_198_windows_test.go` 36 行）
> ＋尾部修正 `9edeba1f` ＋实现者证据件 `.scratch/wisp/probes/198/r1/verdict.md`（本腿现量 **251 行／39,384 字节**）。
> ⛔ 本腿未改任何产码或测试；唯一写面＝本文件＋票面末尾一行进度。
> ⚠ 写作纪律：本文件不逐字落"占位符自查尺"的拼法（落面即自匹配，实现者 §8 刚复踩过一次）。

---

## §0 起手锚与写面

| 项 | 读数（起手逐字） | 尺 |
|---|---|---|
| 起手时刻 | `2026-10-02T09:43:17+08:00` | `date -Iseconds` |
| HEAD（起手） | `5fad9a000f1cf0a6a817488d01bf86824c8b94f4`　`2026-10-02T09:41:23+08:00`　标题 `docs(orchestrator)：界面侧带话落成一段可直接复制的 prompt（owner 10-02 明确要求"我丢给它就行"）` | `git log -1 --date=iso-strict --format='%H %ad %s'` |
| 分支 | `dev` | `git rev-parse --abbrev-ref HEAD` |
| 起手写面 | **0 行**（空输出；`09:47` 复跑同尺仍 0 行） | `git status --porcelain cmd/wisp internal/config` |
| 终态写面 | **0 行＝与起手名册逐枚差集为空**（起手空名册，故差集两侧都是零枚） | 同尺，见 §4 末 |
| 骨架 commit | `5ccd507b`（显式 pathspec＝本文件一枚） | `git log -1 --format='%h %s'` |
| 被测两枚在不在史上 | `git merge-base --is-ancestor` ⇒ `613606c0` **是** HEAD 的祖先、`9edeba1f` **是** | 同尺 |
| 本腿写面（声明） | 仅 `docs/evidence/s1/198-firstrun-config-v1.md` ＋票面 `198-*.md` 末尾追加一行进度（⛔ 未改原文一字、⛔ 未碰任何勾选框） | `git show --name-only` |
| 本腿跑过的尺 | `go build`／`go vet`（四枚包）／`gofumpt -l`（四枚文件）／`sh scripts/d22scan.sh`／`go test ./cmd/wisp`（定向 5 发＋整包 1 发）／`go build -o /tmp/198v1/wisp.exe`＋真件 3 发＋不可写数据根 1 发／七枚变异 | 逐条见 §2／§3 |
| 未碰 | 两枚常驻钉、`internal/config/**`、`internal/risk`／`internal/tools` 的测试（只读不跑）、`scripts/**`（只执行不改）、阈值／golden／`thresholds.go`／`allowlist.txt`、票面 AC 框、台账 | 起手声明，非事后补 |
| 环境硬尺 | 每发 `go test` 都带 `export PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"`（缺它以 `0xc0000135`＋无 `--- FAIL` 收场）；**判红绿只认 `--- FAIL` 行** | 题面 §允许跑的尺 |

---

## §1 逐格 AC 判语

| 格 | 判语 | 凭据（本腿现跑读数；件内引用标行） |
|---|---|---|
| **AC#1 前提**（零配置 ⇒ 真产出，二次不重建，删了重建） | **成立** | 本腿两级读数。①**包内**：`-run '^TestTicket198' -count=1 -v` ⇒ `--- PASS: TestTicket198AC1FirstRunCreatesConfigThenLeavesItAlone (0.03s)`，`firstrun_198_test.go:102` 逐字 `stage 1 config.toml: 2571 bytes, first line "schema_version = 2"`，退码 2／stdout 空由 `:87-92` 断言且未触发。②**真件（实现者没有这一发，本腿自跑 `/tmp/198v1/wisp.exe`）**：`EXIT1=2`／`EXIT2=2`／`EXIT3=2`，第一发建出 `config.toml` 2,571 字节、`md5 5c2ccf78207db680d11603055feb24a5`、data 根只长 `config.toml logs secrets`；第二发 `grep -c 新建默认配置`＝**0**、mtime 停在 `09:56:56.245837300` 未动、md5 相同、`.wisp-config-*.tmp` 残留 0；删了再跑 `grep -c 新建默认配置`＝**1**、2,571 字节、md5 相同 ⇒ 三发**逐字复认**实现者 §3.1（其 `verdict.md:100-109`） |
| **AC#2 默认值同源** | **取值同源成立（本腿独立读数坐实）；"判据能区分"那一半在仓内仍是 部分** | ①本腿自己的交叉尺（**不用 `NewDefaults()` 作期望侧**，正是普查件 §5.2 要的形）：把真件那发产出的 `config.toml` 里每条**非零标量**逐枚回 `internal/config/schema.go` 找 `toml:"k" … default:"v"` ⇒ **45 枚命中**（`language zh-CN`／`theme dark`／`size 56`／`opacity_idle 0.35`／`permission_mode ask_every_step`／`confirm_timeout_sec 300`／`temperature 0.7`×4／`thinking_intensity off`×4…），未命中只有三枚且各有名分：`schema_version = 2`（`loader.go:243` 盖章，`schema.go:104-107` 注释逐字"never defaulted"）、`repeat_thresholds = [3, 5, 8]`（`schema.go:425` `default:"3,5,8"` 逗号拼法）、`veto_words = ['取消','停下','别']`（`schema.go:203` `default:"取消,停下,别"`）⇒ **零枚取值是写手自填**。②**仪器无牙**（本腿变异 V3a 现量）：把 `NewDefaults()` 的结果改一个合法非默认值 `theme="light"` ⇒ 六枚用例**全绿**、文件变 **2,572 字节**——连 `firstrun.go:93-96` 那句"全部取值来自内置默认表"在那一发下是假话，仪器不响。③有牙的那一支（V3b）：整张表换成 `&config.Config{}`（完全不套 default 标签）⇒ `--- FAIL: TestTicket198AC1CreatedFileHoldsEveryStaticSectionAndNoInventedTables`，红句逐字 `firstrun_198_test.go:228: a config first-run wrote cannot be loaded: config: config.toml: app.theme "" must be one of dark|light|auto`。④实现者件 `:11-13`／`:200`（Y6）／`:209`（P3）已**自报**"这一枚不构成对 AC#2 的判定、归 198-r2"⇒ 不是越界报功。⇒ **本腿判：产码侧同源成立；AC#2 那一格按票面"判据要能区分"仍缺仪器，勾它归 198-r2** |
| **AC#3 路径与权限** | **成立**（附一枚字面冲突，归编排者裁） | ①落点：`cmd/wisp/run.go:199-214` 已由现成判定者 `resolveDataDir`（`doctor.go:247`）给根，`firstrun.go:79` 再调 J1 末段点名的同一枚 `secret.NewStore(dataDir)`（`internal/secret/store.go:49` → `winsec.PrivateDirAll`）；`firstrun.go` 代码里 `filepath.` **只有一枚 `Join`**（`:73`），`Clean|Abs` 命中 1 处是 `:70` 的注释名词，剥注释后 0 枚。②权限：本腿跑 `--- PASS: TestTicket198AC3CreatedFileLandsPrivate (0.04s)`，`icacls` 逐字三枚 ACE `NT AUTHORITY\SYSTEM:(F)`／`BUILTIN\Administrators:(F)`／`DESKTOP-LVS7839\swq:(F)`，`Everyone`／`S-1-1-0`／`WD` 零命中（尺＝`logsink_windows_test.go:110-128`）；文件侧来自 `SaveFile`→`atomicWrite`→`winsec.SealFile(tmpName)`（`parse.go:215`，写第一字节之前）。③⚠ **字面冲突登记**：票面 AC#3 括注"不许 `filepath.Join` 手拼"，与 `AGENTS.md §1.2`／`PLAN.md:1288`（禁令射程＝`risk.PathResolver` 之外的 `Clean|Abs`）不一致；照票面字面判，`firstrun.go:73` **和**现存的 `run.go:388` 一起违规——那是把既有生产形状判错，本腿不据此扣分也不赦（普查件 J6 同一格，实现者件 `:212` P6 已上报）。④⛔ 无第二套序列化器：`firstrun.go:82` 唯一一句写盘＝`config.SaveFile(cfgPath, config.NewDefaults())`，import 块（`:38-48`）只有标准库＋`internal/config`＋`internal/secret`。 |
| **AC#4 缺 key 时说实话** | **部分（缺什么 ✓／去哪儿补 ✗）**——本腿不勾它 | 回执逐字（本腿真件那发）：`wisp run: 已在 C:\Users\swq\AppData\Local\Temp\198v1\data1\config.toml 新建默认配置：全部取值来自内置默认表（schema 的 default 标签），未替你选任何模型；[llm] 的模型与 api_key_ref 仍缺，本轮仍按未配置失败退码`，紧接 `wisp run: 文本角色未配置（Unconfigured）：config: llm: role "chat" is unset and text_chain is empty (Unconfigured)` ⇒ 两句给了"缺什么"和文件绝对路径，**没给补救入口**：`wisp secret set` 只出现在 `cmd/wisp/secret.go:162`／`:334` 的 usage 文本里（`grep -rn "secret set" cmd/wisp internal/config` 排除 `_test` 后零枚在 run 那条路上），而 `api_key_ref` 指的是 DPAPI blob 的**名字**、用户手改 toml 补不上 key ⇒ 按票面"指名缺什么、**去哪儿补**"这一格**未闭合**。实现者件 `:44` 末段与 `:209`（P3）**已具名划界**"AC#4 归 198-r2"，且 §6-P3 明写"若验收把四要素当 AC#1 读会判成缺半"——本腿按票面读，判缺半，不判谎报。 |
| **AC#5 不静默** | **行为侧成立（本腿活量读数）；仪器侧无牙** | ①**成功不静默有牙**（变异 V5＝删掉回执）⇒ 两枚红：`firstrun_198_test.go:99: stage 1: the run created the file but never said so with its path (AC: 不许静默建完就当没事发生). stderr:` 与 `firstrun_198_test.go:241: no creation receipt to test against; stderr:`。②**失败侧行为有活量读数**（⛔ 非变异，本腿拿真 exe 撞不可写数据根：`WISP_TEST_DATA_DIR=C:\Windows\System32\wisp198v1-denied`）⇒ `EXIT=2`，stderr 第 4 行逐字 `wisp run: 首次配置建不出来（数据根不可建（C:\Windows\System32\wisp198v1-denied）：secret: create C:\Windows\System32\wisp198v1-denied\secrets: mkdir C:\Windows\System32\wisp198v1-denied: Access is denied.）：本轮仍按缺配置响亮失败，不会拿半份配置顶上`，第 3 行另有 `持久日志未启用（… Access is denied.）`、第 5 行 `凭据存储不可用：…`；盘上核验 `ls /c/Windows/System32/wisp198v1-denied` ⇒ `No such file or directory`＝**没建出半份文件、也没报成功**（票面 AC#5 的"不许建出一份空文件还报成功"这一形本腿量到）。③**但仪器无牙**（变异 V4＝把 `run.go:245-247` 那句吞成 `_ = frErr`）⇒ `--- PASS` ×6（198 全枚）＋票 101 三形全绿；`grep -rn "首次配置建不出来" cmd/wisp internal docs` 在任何 `_test.go` 里**零命中** ⇒ "建了没建成必响亮"今天**没有一枚用例守着**。④"只读盘／磁盘满"另两支未量（本腿只量了"目录不可建"一支）⇒ 见 §6 甲-2。⇒ **判语：行为成立、仪器缺一枚钉；这一格本腿不勾** |
| **J1 落点裁定**（首建只挂 `wisp run`，⛔ 不进 `assembleRuntime`） | **成立，且是被执行的、不是被注释的** | ①**基线色（本腿现跑）**：`-run '^(TestAC246ShippedResidentLegTakesTheInjectionAndAttemptsThePipeline\|TestAC246ShippedResidentLegWithoutConsoleRefusesItsTaskEntry\|TestAC2SealNoticeLandsInTheRunLegLogFile\|TestAC2AuditTrailLandsInTheRunLegLogFile\|TestAC2RunLegInstallsTheSinkBeforeItsFirstSealingSite\|TestAC2EveryLegRefusesTheSameShapeAndWritesNothing128\|TestTicket101UnreadableModeFailsLoudlyAndStrict\|TestMissingBlobFailsUnconfiguredNeverSilently)$' -count=1 -v` ⇒ **8 枚全 `--- PASS`、零 `--- FAIL`、零 SKIP，`ok … 11.243s`**；其中 `TestAC246ShippedResidentLegTakesTheInjectionAndAttemptsThePipeline (3.44s)` 跑的是 `buildWispForTest` 的**真 exe**（`resident_task_source_246_windows_test.go:359`），`:389-391` 那枚"常驻腿不许造 config.toml"的断言**今天就是绿的** ⇒ 常驻那条路（`wisp` 起来不跑 `run`）没造出 `config.toml`＝**活量读数，不是读注释**。②**本腿亲自把它搬进 `assembleRuntime`（变异 V1）**⇒ 正是这两枚红（逐字见 §3-V1），而 198 自己的 AC 用例（三段式）与 logsink 那枚**仍绿** ⇒ 落点错了不会被本票用例抓到、只被 J1 那两枚钉抓到 ⇒ 裁定的牙在位。③生产调用边：`grep -rn "runTextTask(" cmd/wisp --include=*.go \| grep -v _test` ⇒ 只有 `main.go:159`（cmdRun）与定义 `run.go:181`；`ensureFirstRunConfig` 的调用者是 `run.go:245`（`runTextTask` 体内、`:249` 那句 `assembleRuntime` 之前）。④插入次序两处都对：严格晚于 `resolveDataDir` 的成功分支（`:208-212` 先打印先退码），严格晚于 `installLogSink`（`:222`）⇒ 票 128 的腿表那枚钉（`TestAC2EveryLegRefusesTheSameShapeAndWritesNothing128` 五枚子用例 `runTextTask`／`cmdModels`／`cmdProviders`／`cmdDoctor`／`resolveSecretLayout` 全 PASS）与 117 的"sink 早于首个封权点"都没被挤动。 |
| **两枚常驻钉未被放宽／未触碰** | **成立** | `git show --name-only 613606c0` 五枚文件名逐枚＝`.scratch/wisp/probes/198/r1/verdict.md`、`cmd/wisp/firstrun.go`、`cmd/wisp/firstrun_198_test.go`、`cmd/wisp/firstrun_acl_198_windows_test.go`、`cmd/wisp/run.go` ⇒ **不含那两枚钉**；`git log --since=2026-09-30 -- cmd/wisp/resident_task_source_246_windows_test.go cmd/wisp/logsink_windows_test.go` 只有 `2071f59e`（2026-09-30，两枚钉自己的交付腿，`merge-base --is-ancestor` 判其**早于** `613606c0`）⇒ 自交付以来无人再动，本腿与实现者都没动。`grep -c "t.Skip"` 三枚新文件＝**0**。本腿全程未碰那两枚钉一字（§0 写面声明＋终态差集为空坐实）。 |
| **退码 2 未被放宽** | **成立** | ①`git show --numstat 613606c0 -- cmd/wisp/run.go` ⇒ **15/0 纯插入**（删零行＝没改任何判据与退码表）；创建失败那支 `run.go:245-247` **只打印不 return**，继续 `rt, code := assembleRuntime(s)`。②本腿真件三发退码全 `2`（含失败那发 `EXIT=2`）；包内三段式在 `:87`／`:111`／`:141` 三处各断 `code != 2 ⇒ Fatalf`，全绿。③V2 变异（装载失败就建）下票 101 的 `code == 2`／`reached == false` 断言在 存储损坏／版本不认识 两形红 ⇒ 说明"首建悄悄把归因盖掉"才是这一路的真危险，**今天的产码没发生**。 |
| **无第二真相源**（⛔ 新默认值表／模板串／第二套序列化器） | **成立** | `firstrun.go` 全文取值只有一枚来源：`:82` `config.SaveFile(cfgPath, config.NewDefaults())`（`loader.go:238` 与 `defaults.go:58` 都已导出）；`grep -nE '"[A-Za-z0-9_.\-]+ *=' cmd/wisp/firstrun.go` ⇒ **0 命中**（无 TOML 模板字符串）；`grep -vE "^\s*//" cmd/wisp/firstrun.go \| grep -n "filepath\."` ⇒ 整份产码只有 `:73` 那一枚 `Join`；import 块无新依赖。⛔ 未用未导出的 `fileMissing`：本腿读到 `parse.go:233-235` 逐字 `return errors.Is(err, fs.ErrNotExist)` ⇒ `firstrun.go:74` 是它的**等价式**不是第二套判定。（仪器有没有牙是另一格，见 §1-AC#2 与 §3-V3a。） |
| **实现者三发自报读数** | **复认**（两级都量） | §1-AC#1 那行的 ①②；字节数 2,571、首行 `schema_version = 2`、stderr 回执逐字、stdout 空、退码 2、mtime/md5 未动、`.wisp-config-*.tmp` 零残留、删了再建字节一致——**九项全复现**。本腿另加两级实现者没有的读数：真 exe 三发＋不可写数据根一发。 |
| **实现者 §5-Y1（陷阱 #1 因果链）** | **复认，并再推进一步**（详见 §3-V2／V2b） | 复认：V2 下目录形子用例 `    --- PASS: TestTicket198AC1DirectoryShapedConfigIsNeverOverwritten` 对应的那枚票 101 子用例 `--- PASS: …/权限读不到 (0.02s)`，红的是 `…/存储损坏`＋`…/版本不认识`（六行逐字见 §3-V2）。**再推进一步**：V2b（"任何 stat 失败就建"这一形）下**六枚 198＋票 101 三形全绿** ⇒ 目录那枚用例对"判据选哪一枚"**完全不敏感**（目录 stat 得动，既不是 not-exist 也不是 stat 失败），它守的是写盘那一侧（`atomicWrite` 的 temp+rename 撞目录必失败，`parse.go:226`）。⇒ **票面 §4 第一枚陷阱的因果（"目录形 ⇒ 判据必须用 fileMissing"）双重不成立**，本腿与实现者各贡献一层；产码选择本身仍正确（V2 咬归因链那一发独立成立）。 |
| **实现者 §5-Y5（"18 section"实为 17＋1）** | **复认** | 本腿现量：`grep -cE "type [A-Za-z]*Section struct" internal/config/schema.go` ⇒ **18**；带 `toml:"…"` 且类型是 Section 的字段 ⇒ **17**；`Plugins` 于 `schema.go:131` 逐字 `toml:"-"`（另 `:149` 的 `Portable` 也是 `toml:"-"` 但那是叶子）；真件产出的 `config.toml` 里 `[plugins]` 头恰 **1** 次、`[llm.providers`／`[models.local_override`／`[plugins.` 三枚动态表头 **0** 次 ⇒ "静态 18 枚各出现一次"这句**不精确**，实为 17 枚带标签的头＋`[plugins]` 由 `marshalPlugins` 自己落一次。 |

---

## §2 门禁四数（本腿现跑，带 HEAD）

HEAD 说明：本腿那几把尺跑在 `5ccd507b`（骨架那发）与 `1e7debde`（并发腿推进后的终态）两枚锚上；两枚锚下 `cmd/wisp` 与 `internal/config` 的工作树内容**逐字节等于 `9edeba1f`**——尺＝`git diff 9edeba1f..HEAD -- cmd/wisp internal/config` ⇒ **空输出**（本腿七枚变异全部回滚后复量）。

| 门禁 | 尺（完整命令） | 现值 |
|---|---|---|
| `go build` | `go build ./cmd/wisp ./internal/config ./internal/secret ./internal/winsec` | **rc=0**（起手锚与终态锚各一发，`FINAL-BUILD-RC=0`） |
| `go vet` | `go vet` 同四枚包 | **rc=0**（`FINAL-VET-RC=0`，静默） |
| `gofumpt -l` | `"$(go env GOPATH)/bin/gofumpt.exe" -l cmd/wisp/firstrun.go cmd/wisp/firstrun_198_test.go cmd/wisp/firstrun_acl_198_windows_test.go cmd/wisp/run.go \| wc -l` | **0**（四枚包/文件零命中；起手与终态各一发同为 0） |
| `sh scripts/d22scan.sh` | `sh scripts/d22scan.sh` | **`d22scan: clean - no D22 ban violations`**；`examined 262 production Go files`；分母逐字：**ban#8 `internal/`=482、`cmd/`=86**（⚠ 编排者那发是 481/86 ⇒ `internal/` 的分母在本腿两发之间**自己涨了 1 枚**，是并发腿往 `internal/` 落了文件，与本票无关；`cmd/`=86 逐字相同）；另仪器自报边界 `skipped as git-ignored: 1 file(s) … frontend/dist/assets/` |
| `staticcheck` | ⛔ 本腿未跑（本机版解不开 go1.27 的 export data，会产"干净的绿"） | **这一格无凭据，明写不做**（见 §6 乙-1） |
| `cmd/wisp` **整包** | `go test ./cmd/wisp -count=1 -v`（带 sherpa PATH，`-count=1` 一发，未反复跑） | **174 枚顶层 `--- PASS`／1 枚 `--- FAIL`／0 枚 SKIP，`FAIL … 186.783s`**。那枚红＝`TestTicket223HandEditedFsLooseningCostsAnL2Card (2.08s)`，红句 `config_reload_223_test.go:303: the console did not render the reload card: …`＋`:306: the card does not name the key it is asking about: …` ⇒ **三发行文都在，是 stdout 单读读早了**。**隔离复量**：`-run '^TestTicket223HandEditedFsLooseningCostsAnL2Card$' -count=2` ⇒ **两发全绿（2.12s／2.09s）**。**归因于票 223 而非本票**凭据三枚：①那枚用例的 host `newReplyHost` 在 `approval_reply_201_test.go:118-145` **先把 `config.toml` 写进 dataDir** ⇒ `firstrun.go:74` 走早退分支、既不写也不打印；②整包日志里 `新建默认配置` **只命中 1 次**（来自 198 用例自己的 `t.Logf`）；③台账在册同形：`docs/reports/pending-and-issues.md:9626`（该枚 `-count=3` 3/3 绿，判争用红不记账）、`:10481`（同码两发一红一绿，判决行 `--- FAIL … (2.43s)` vs `--- PASS … (2.13s)`）、`:10553`（A518 结档：根因＝`config_reload_223_test.go:301` 那枚**未加界的 stdout 单读**，归票 223 裁决侧）。⇒ **本腿判：`cmd/wisp` 整包今天不是全绿，但那枚红与票 198 无关、且不是新红** |

---

## §3 本腿打过的变异（七枚，逐枚已回滚；回滚凭据＝`git status --porcelain cmd/wisp internal/config` 终态 **0 行**＋`grep -c "198v1-V" cmd/wisp/*.go` **0 枚非零**＋回滚后 `-run '^TestTicket198' ⇒ ok 0.210s`）

| # | 改了什么（文件:行） | 结果 | 红句逐字 |
|---|---|---|---|
| **V1** | `cmd/wisp/run.go`：删掉 `:245-247` 那一跳，原样搬进 `assembleRuntime`（`:384` 的 notify 默认之后、`:387` 那行之前）——**本腿亲自违反 J1 裁定一次** | **两枚红**（正是要它红的那两枚）；`TestTicket198AC1FirstRunCreatesConfigThenLeavesItAlone` 与 `TestAC2SealNoticeLandsInTheRunLegLogFile` **没红**（那是信息：本票的 AC 用例对落点不敏感） | `firstrun_198_test.go:281: AC J1 RED: ensureFirstRunConfig is called by [assembleRuntime], want exactly [runTextTask] - the resident leg reaches the assembly root directly, and first-run creation shared there writes a config.toml its pins forbid`　／　`resident_task_source_246_windows_test.go:390: AC#7 RED: the leg created a config.toml it was never asked for (err <nil>)`（真 exe 那发，5.08s） |
| **V2** | `firstrun.go:73-77`：判据换成实现者的 M1 形 `if _, _, lerr := config.LoadFile(cfgPath, nil); lerr == nil { return false, nil }` | **目录形没红**（`    --- PASS: TestTicket101UnreadableModeFailsLoudlyAndStrict/权限读不到 (0.02s)`），**红的是两枚归因子用例**＝复认实现者 §3.2-M1；`TestTicket198AC1*` 三枚全绿 | `run_mode101_test.go:661: no MODE-READ-FAILED audit line for 存储损坏:`／`:664: the failure did not name the strictest档 it falls back to:`／`:667: the user-visible line is missing for 存储损坏:`／同三行换 `版本不认识` ⇒ 六行 |
| **V2b** | `firstrun.go:74`：判据换成最偷懒的 `if _, err := os.Stat(cfgPath); err == nil { return false, nil }`（＝"任何 stat 失败都建"，⛔ 不是 `fs.ErrNotExist`） | **全绿（七枚用例＋三形子用例无一红）**：`TestTicket198` 六枚 `--- PASS` ＋ `TestTicket101…` 顶层 PASS ＋ 子用例 存储损坏／版本不认识／权限读不到 三枚 PASS，`ok 1.636s` ⇒ **今天的仪器分不开 shipped 形与这一形**；目录那形对判据选择**不敏感**（它 stat 得动）。⇒ 票面 §4 陷阱 #1 的因果被**再推翻一层**（见 §1 倒数第三行、§5-3） | 无（"没红"也是读数） |
| **V3a** | `firstrun.go:82`：`c := config.NewDefaults(); c.App.Theme = "light"` 再 `SaveFile`（合法值、非默认标签） | **六枚全绿**、文件从 2,571 变 **2,572 字节** ⇒ 仪器**区分不了"照抄默认"与"改一个合法值"**＝票面 AC#2 要的那把"能区分"的尺**在仓内不存在**；回执那句"全部取值来自内置默认表"在此下发假而无人响 | 无（"没红"也是读数） |
| **V3b** | `firstrun.go:82`：整张表换成 `&config.Config{}`（完全不套 `default` 标签） | **一枚红** ⇒ "整张表是手搓的"这一形**有牙**（牙长在可加载性那一格） | `firstrun_198_test.go:228: a config first-run wrote cannot be loaded: config: config.toml: app.theme "" must be one of dark|light|auto` |
| **V4** | `run.go:245-247`：把"首次配置建不出来…"那句吞成 `_ = frErr`（创建失败静默） | **全绿**（`TestTicket198` 六枚＋票 101 三形，`ok 2.544s`）＋`grep -rn "首次配置建不出来" cmd/wisp internal docs` 在任何 `_test.go` 零命中 ⇒ **AC#5 的失败支没有仪器**（行为本身有活量读数，见 §1-AC#5 ②，那是另一发⛔非变异的真件量） | 无（"没红"也是读数） |
| **V5** | `firstrun.go:93-96`：删掉成功回执那一句 | **两枚红** ⇒ "建了但没说"有牙 | `firstrun_198_test.go:99: stage 1: the run created the file but never said so with its path (AC: 不许静默建完就当没事发生). stderr:` ／ `firstrun_198_test.go:241: no creation receipt to test against; stderr:` |

> 七枚变异逐枚回滚后本腿复跑过 `-run '^TestTicket198' ⇒ ok 0.210s` 与 `git diff 9edeba1f..HEAD -- cmd/wisp internal/config ⇒ 空`，`t.Skip` 零枚、任何断言一字未放宽（本腿唯一写面＝§0 声明那一枚文件＋票面末尾一行）。

---

## §4 退回与否

**不退回到"重写"那一档：AC#1／AC#3／J1 落点／两枚钉未被碰／退码 2 未放宽／无第二真相源——六格本腿判成立且凭据都是本腿现跑读数（含一票真件读数与两枚"违反裁定才红"的变异）；实现者三发读数全部复认，两处自报推翻（Y1 因果、Y5 枚数）本腿复认并各再推进一层。**
**要退回的是两格缺凭据的活，指名：AC#2（第 2 格）缺"逐值同源"的仪器——V3a 那发全绿就是证据；AC#4（第 4 格）缺"去哪儿补"那句话；AC#5（第 5 格）缺"建了没建成必响亮"的用例——V4 那发全绿就是证据。三格实现者都已具名划界／自报（件内 `:44`、`:200`、`:209`），本腿按票面判它们未闭合，⛔ 编排者今天不许勾这三格。**
另需编排者落手的两笔更正（都不是本腿能改的）：①票面 §4 第一枚陷阱形的**因果**要按 V2＋V2b 两发改写（记我这一发为第二枚凭据）；②AC#3 括注"不许 `filepath.Join`"与 `AGENTS.md §1.2`／`PLAN.md:1288` 的射程冲突（J6）要定一句验收口径，否则下一腿还会撞。

---

## §5 我可能判错的条目（每条附"如果错了后果"）

1. **`cmd/wisp` 整包那枚红我判成"与票 198 无关"**。凭据是三段：`approval_reply_201_test.go:118-145` 预写 `config.toml` ⇒ `firstrun.go:74` 早退；整包日志 `新建默认配置` 只命中 1 次；台账 `:9626`／`:10481`／`:10553` 在册同形。⚠ 我**没有**跑过"把 `613606c0` revert 掉再整包"那一发（⛔ 共享工作树里不许 revert/reset，题面也禁）⇒ 归因靠的是代码路径＋在册记录，不是 A/B。**如果错了**（其实 198 的插入把那枚 stdout 的时序挤早/挤晚了一拍）**后果**：我把一枚真回归放过去、编排者勾了格，下一轮 CI 又红⇒ 白烧一发整包（约 190s）＋台账多记一枚 R。补救：见 §6 乙-2（谁都能跑的那一把）。
2. **V2b 我读成"目录形对判据不敏感"**。我只在**这台 Windows、这些用例**上量；`fs.ErrNotExist` 与 `err == nil` 的真分歧在"stat 得动之外的那一形"（ACL 拒读、reparse、只读父目录），今天全仓没有一枚用例造得出那一形（票 101 的注释 `run_mode101_test.go:628-633` 逐字说目录是"权限读不到"的**便携替身**，正是为了让 `t.TempDir` 删得掉）。**如果错了**（其实存在我漏跑的用例能区分）**后果**：我给"判据选型"记了一枚比实际更强的"无牙"，下一腿就不补那枚用例 ⇒ AC#5／票 101 归因链多一条无人守的路。
3. **我复认了"陷阱 #1 因果推翻"并再加一层**。若编排者坚持票面 §4 那句因果（目录形 ⇒ 必须 `fileMissing`），**后果**：我的 §4 更正建议是多余的，且 V2b 那发读数会被当成"用例没覆盖到的形"而非"因果写错"——两种处置会决定下一腿补哪枚用例，别混。
4. **AC#2 我给"取值同源成立"用的是我自己的交叉尺**（把真件产出的 toml 逐值回 `schema.go` 找 `default` 标签）。它的漏网形是"值恰好等于另一个字段的默认标签"（同名不同段的取值我按 `toml:"k"` 匹配，未做段路径区分）与"零值/空字符串那一档"（我按 Go 零值跳过：`''`／`0`／`false`／`[]`／`0.0`，共 40 余条没逐一回查标签）。**如果错了**（某枚跳过的值其实是写手填的）**后果**：AC#2 的"同源"被我判成成立而实际漏一枚——但 V3a 已经证明**仓内仪器本来也漏**，所以后果不会比 §1 那格写的更糟。
5. **AC#5 我把"目录不可建"那一发算作行为侧成立**。票面点名的三形是"只读盘／磁盘满／目录不存在"，我只量到第三形；前两形没量（本腿不许写用例，也拿不到磁盘满）。**如果错了**（前两形里 `SaveFile` 失败但 `atomicWrite` 留下半截文件、或失败句没打出来）**后果**：AC#5 被我高判一格，真实缺口归零失败。已按未闭合处理，⛔ 这一格没勾。
6. **`staticcheck` 那格我按题面直接写"无凭据"**，没试图换版本跑。如果编排者手里有能解 go1.27 export data 的版本，**后果**：那格本该有第四把尺而今天空着，本腿的 §2 表少一行凭据（不影响其余判语）。

---

## §6 判不动的地方（甲＝谁补得上＋命令＋期望读数；乙＝补不上，明写不做）

**甲**

1. **AC#2 的仪器**（谁：198-r2）。命令：在 `cmd/wisp/firstrun_198_test.go` 加一枚用例，**期望侧不引 `NewDefaults()`**——反射 `config.Config` 的 `toml:`＋`default:` 标签自己重述一张期望表，逐行比对真件产出的 `config.toml`。期望读数：现产码 ⇒ 绿；本腿 §3-V3a 那种"改一个合法非默认值" ⇒ **红**（这把尺必须先把 V3a 跑红才算装上牙）。
2. **AC#5 失败支的仪器**（谁：198-r2）。命令：临时把 `run.go:246` 那句吞掉（＝本腿 V4）跑 `-run '^TestTicket198'`，期望从"全绿"变**红**；再加一枚用例把数据根指向不可写位置（本腿用 `C:\Windows\System32\…` 量到过那句逐字），断言 stderr 含"首次配置建不出来"＋退码 2＋盘上没有 `config.toml`。
3. **stat 失败但不是 not-exist 那一形的仪器**（谁：票 101／票 89 那一族的地界）。命令思路：对**父目录**下 `deny (RX)` 再 stat（`icacls` 那把尺），期望 `os.Stat` 回 `Access is denied`（⛔ 不是 `fs.ErrNotExist`）；期望 shipped 判语＝**不建、把归因交给 `LoadFile` 那一支**，V2b 那种形＝会去建 ⇒ 红。这一枚本腿造不出（不许写用例，且 `icacls /deny` 会让 `t.TempDir` 删不掉——票 101 的注释已记过这个代价）。
4. **磁盘满／只读盘那两形**（谁：需要盘配额的腿）。命令：把数据根指到 read-only 卷或 `Set-StorageTier` 满配额处，跑真 exe `wisp run "…"`。期望读数：`首次配置建不出来（…）` 一句在、退码 2、`config.toml` 不存在或仍是旧内容、`.wisp-config-*.tmp` 零残留。
5. **"整包红名册"这一格**（谁：编排者的验证窗口，已给凭据）。命令：`export PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"; go test ./cmd/wisp -count=1`。期望读数：`--- FAIL` 只该有在册那枚争用红；本腿现量＝**174 PASS／1 FAIL／0 SKIP**。

**乙（明写不做）**

1. `staticcheck`：本机版解不开 go1.27 的 export data ⇒ 跑了也是"干净的绿"，**本腿不做，那一格无凭据**（题面原文已把这一格判成无凭据，本腿复认）。
2. `internal/risk`／`internal/tools` 的任何测试：题面硬禁（`TestResolvePerCallBudget` 有腿在同机做非计时取证），**本腿一枚没跑**，两枚包只读不跑。
3. `frontend/**`／`design/**`：既不读也不引，本件里除仪器自报的 git-ignored 边界行（`frontend/dist/assets/`，d22scan 输出逐字）之外零引用。
4. 票面 AC 勾选框：一枚未碰；票面只在末尾**追加一行**进度，原文一字未改。
5. `scripts/**`／阈值／golden／`thresholds.go`／`allowlist.txt`／两枚常驻钉：零写面。`sh scripts/d22scan.sh` 是执行不是修改（`git status` 终态差集为空即凭据）。
