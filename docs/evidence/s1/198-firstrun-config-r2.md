# 票 198 · 落地写腿 `198-r2` · 证据件（四小格：AC#2 仪器／AC#5 失败支用例／AC#4 "去哪儿补"／落点仪器收回本票）

> 本件是**实现腿**（写腿）的证据。射程＝票面 §6「10-02 10:1x 验收后裁定」第 5 条排程的四小格：
> ① 给"改一枚默认值"装上牙（AC#2）；② 给"首建失败那句"补一枚用例（AC#5）；
> ③ AC#4 的"去哪儿补"要先回答 DPAPI 那一层；④ 把"落点错了"从别人家的钉收回到本票自己的仪器。
> 上一枚非实现者裁决表＝`docs/evidence/s1/198-firstrun-config-v1.md`（本件要补的洞由它 §1／§3／§6-甲点名）。
> ⛔ 本腿不碰任何 AC 勾选框；⛔ 一字未动 `docs/PLAN.md`／`docs/specs/**`／阈值／golden／`allowlist.txt`／
> `internal/panel` 两枚钉／`internal/perm/ticket90_persist_test.go`／`.github/workflows/ci.yml`。
> ⛔ `frontend/**`／`design/**` 不读、不写、不引。

---

## §0 起手锚（同发取 date／HEAD／porcelain）

**〔本节由编排者代提——腿死于 15:17 前后的模型服务当日额度（通知 15:4x 到、只当未验证；盘上现量为凭），产码它已提交（`fa5593c0`）但证据件还是骨架；下面每一行是编排者 15:46–15:48 现跑，非本腿自跑〕**

- `date`：`2026-10-02 15:46:57 +0800`。
- HEAD：**`f2e5f72a`**（＝腿自己的骨架发，15:18）。
- `git status --porcelain -- cmd internal`＝**0 行**⇒ 本腿产码已全部入库、无突变残留（腿死前把 `mutate.sh` 落在台件目录但**没有执行残留**：写面干净）。
- 本腿两发产码提交：`1d8106d8`（证据件骨架，15:18 前落）＋ **`fa5593c0`**（15:12，`cmd/wisp/firstrun.go` **23/1**＋新增 `cmd/wisp/firstrun_198r2_test.go` **473/0**，去重路径＝2）。

---

## §1 四格各自的修法与判据（每格带 file:line）

**〔本节同上，编排者从 `git show fa5593c0` 的 diff 里逐行读出来代提；判语一格不写〕**

1. **AC#2（改一枚默认值必须响）**：新增 `cmd/wisp/firstrun_198r2_test.go:221` `TestTicket198R2AC2CreatedFileCarriesNothingButTheSchemaDefaultTags` 与 `:274` `TestTicket198R2AC2ExportedDefaultTableStillRendersAsItsTags`。形状＝**期望侧不抄任何字面默认值**，只读 `schema.go` 各字段的 `default` 标签、与首建写出的文件**逐字节比**（提交标题原话"自带八枚种下必红的正控"——哪八枚、种在哪，属 §3 变异自证表，本节不代填）。
2. **AC#5（失败句要有用例守）**：新增 `:299` `TestTicket198R2AC5CreationFailureIsLoudAndLeavesNoHalfFile`。形状＝把 secrets 存储落点用**普通文件占位**（winsec 拒密封那族）⇒ 钉住三样：失败句响亮、盘上不留半份文件、进程仍退码 2。
3. **AC#4（"去哪儿补"说真话）**：产码只动 `cmd/wisp/firstrun.go:92-115`（diff `+23/-1`）——回执从一句涨到三段，**只写现读出来的两条真走得通的路**：key 走 `wisp secret set <blob 名>`（隐藏输入／`--from-stdin`，明文进 DPAPI 存储、文件里补的是 `api_key_ref = "dpapi:<blob 名>"` 那个**名字**，不是明文；不用 DPAPI 可写 `env:<环境变量名>`）；模型没有 CLI 写入者（`wisp providers discover/probe` 只读该文件去问端点、不代写）⇒ 文案照实说"只有改这份文件一条路"。⚠ 三段里只有命令名／字段名／占位名，**无凭据值**（编排者只读 diff 复核过）。
4. **J1 落点仪器收回本票**：新增 `:432` `TestTicket198R2J1TheAssemblyRootStillCreatesNothing`（两枚同样干净的空目录：装配根不建、`wisp run` 入口建）⇒ "落点错了"从此不再只靠别人票的两枚钉。

---

## §2 门禁四数（build／`gofumpt -l` 本腿动过的目录／`go test` 相关包 `-count=1`／`d22scan` exe）

**〔本节同上，编排者 15:46–15:48 现跑代提，HEAD `f2e5f72a`、写面 porcelain＝0；数字与署名各归各，⛔ 不是本腿自跑〕**

- **尺一 `GOFLAGS= go build ./...`** ⇒ **rc=0、合并输出 0 字节**（15:46:57）。
- **尺二 `gofumpt -l cmd/wisp`** ⇒ **rc=0、list-count=0**（15:47:08；真身 `$(go env GOPATH)/bin/gofumpt.exe`，不在本 shell PATH）。⚠ 口径＝本腿动过的目录（`cmd/wisp`）。
- **尺三 `go test ./cmd/wisp ./internal/config -count=1`**（带 sherpa PATH 前缀）⇒ **合并 rc=0**：`ok cmd/wisp 336.478s`／`ok internal/config 3.759s`——**整包零红**，含本腿新增那五枚用例；那枚既有间歇红（`TestTicket223ModeLooseningChangesTheRunningModeAfterAllow`）这一发**没赶上**（它在册读数现为 5 发整包里 1 红，见台账 `A533` §3）。⚠ 336.478s 按【同机可能有争用】读。
- **尺四 `./tools/d22scan/d22scan.exe`** ⇒ **rc=0**，逐字 `clean - no D22 ban violations`＋`examined 262 production Go files under internal/ and cmd/`；分母变化具名：`ban #8 cmd/` 从前值 86 涨到 **87**（＝本腿新增的那枚测试件）。

---

## §3 变异自证表（每格一行：种什么形 → 哪枚必须红 → 复跑终值）

**〔本节由验收腿 `198-v2` 自己取数：六发变异，2026-10-02 16:16–16:23 现跑，HEAD `315c6f5f`→`fa9a1fc6`（两锚下三枚突变标的文件 md5 逐字一致，`git diff HEAD -- cmd internal`＝0）。每发前 `git status --porcelain -- cmd internal`＝0；突变前原档落 `.scratch/wisp/probes/198/v2/backup/`（md5：firstrun `1a75f24d`／run `1d95cfaf`／defaults `ad3e7da4`）；测试输出先落 `probes/198/v2/mut-M*.txt`，rc 取 go test 本体；每发后用备份覆写还原，porcelain 复量＝0 才进下一发。判红绿只认 `--- FAIL` 行〕**

**基线（还原后复绿凭据）**：`go test ./cmd/wisp -run '^TestTicket198R2' -count=1 -v` ⇒ **五枚全 `--- PASS`、rc=0**（`baseline-verbose.txt`，16:23:54，锚 `fa9a1fc6`）：

| # | 种什么形（文件:改动） | 哪枚必须红 | 复跑终值（逐字） |
|---|---|---|---|
| **M1** | `cmd/wisp/firstrun.go:82`：`config.NewDefaults()` 后种一枚**合法非默认值** `c.App.Theme = config.ThemeLight`（＝v1 的 V3a 形，写侧漂） | `:221 AC2CreatedFileCarriesNothingButTheSchemaDefaultTags` | **rc=1**：`--- FAIL: TestTicket198R2AC2CreatedFileCarriesNothingButTheSchemaDefaultTags (0.02s)`；红句逐字 `firstrun_198r2_test.go:239: AC#2 RED: the config.toml first-run wrote is not the schema's `default` tags rendered: first difference at line 5 (byte 54): want "theme = 'dark'", got "theme = 'light'"; total bytes want 2571 got 2572`；**`:274` 那枚同发 PASS**（0.00s）⇒ 两枚用例正好分清"写侧漂"与"源侧漂" |
| **M2** | `internal/config/defaults.go:60`：`NewDefaults()` 在 `applyDefaults` 之后种同一枚 `c.App.Theme = "light"`（源侧漂：默认表自己被改） | `:274 AC2ExportedDefaultTableStillRendersAsItsTags` | **rc=1，双红**：`--- FAIL: …AC2ExportedDefaultTableStillRendersAsItsTags (0.00s)`，红句逐字 `firstrun_198r2_test.go:279: AC#2 RED: internal/config's default table is no longer just its `default` tags: first difference at line 5 (byte 54): want "theme = 'dark'", got "theme = 'light'"`；**`:221` 也红**（写侧继承同一漂移，方向正确）；v1 V3a 那发"六枚全绿"的洞今日闭合 |
| **M3** | `cmd/wisp/run.go:245-247`：失败句整块吞成 `_, _ = ensureFirstRunConfig(…)`（＝v1 的 V4 形） | `:299 AC5CreationFailureIsLoudAndLeavesNoHalfFile` | **rc=1**：`--- FAIL: TestTicket198R2AC5CreationFailureIsLoudAndLeavesNoHalfFile (0.01s)`；红句逐字 `firstrun_198r2_test.go:316: AC#5 RED: the run could not create its first config and never said so (run.go:245-247 is the sentence this case exists for). stderr:`；v1 V4"全绿"的洞今日闭合 |
| **M4** | `cmd/wisp/firstrun.go:111-118`：**"去哪儿补"两段整段删除**（回执只剩"缺什么"） | `:363 AC4ReceiptNamesTheRealEntryPoints` | **rc=1**：`--- FAIL: TestTicket198R2AC4ReceiptNamesTheRealEntryPoints (0.12s)`；红句逐字 `firstrun_198r2_test.go:384: AC#4: the first-run receipt never names "wisp secret set", so it still says what is missing without saying where to fix it:`（同形再报 `dpapi:`／`env:`／`[llm.providers.`／`roles.chat` 五行） |
| **M5** | `cmd/wisp/firstrun.go`：回执里 `wisp providers discover` 改成**不存在的** `wisp catalog sync`（编造入口点） | `:363` 的 AST 能力检查半 | **rc=1**：`--- FAIL: TestTicket198R2AC4ReceiptNamesTheRealEntryPoints (0.09s)`；红句逐字 `firstrun_198r2_test.go:407: AC#4 RED: the guidance sends the user to entry points this build does not dispatch: [wisp catalog (no cmdCatalog in package main)]` |
| **M6** | `cmd/wisp/run.go`：把 `ensureFirstRunConfig` 那跳从 `runTextTask` 搬进 `assembleRuntime`（`:390` 之前，＝v1 的 V1 形） | `:432 J1TheAssemblyRootStillCreatesNothing`＋既有 AST 钉 | **rc=1，双红**：`--- FAIL: TestTicket198R2J1TheAssemblyRootStillCreatesNothing (0.03s)`，红句逐字 `firstrun_198r2_test.go:449: J1 RED: the assembly root created a config.toml nobody asked it for (stat err <nil>) - the resident leg reaches this function directly…`＋`:453` 回执半也红；**既有 AST 钉** `--- FAIL: TestTicket198FirstRunCallerIsTheRunEntryOnly (0.07s)`，逐字 `firstrun_198_test.go:281: AC J1 RED: ensureFirstRunConfig is called by [assembleRuntime], want exactly [runTextTask]` |

**六发全部还原后**：porcelain=0（逐发复量）＋基线五枚全绿（上锚）＋`git diff HEAD -- cmd internal` 空。原档只覆写自己动过的三枚文件，未用 `git checkout`/`restore`。**注**：死腿台件 `.scratch/wisp/probes/198/r2/mutate.sh` 的 M1–M8 形状表只当线索读，本表六发是本腿自己种的、读数自己取的；死腿没跑过任何一发，不存在"照抄读数"这回事。

---

## §4 我可能写错的条目

（由验收腿 `198-v2` 填）

1. **M1/M2 我只种了 `theme` 这一枚字段**，没逐枚种满八类 kind 的生产侧（用例内部的正控 `:248-266` 已覆盖八类 kind 的"渲染会动"，但那是期望侧自证，不是我种进生产码的）。如果那八类里有某一类的**生产写侧**另有旁路（比如某段被 `marshalPlugins` 之外的特例手写），我这套六发抓不到它。**后果**：AC#2 的牙对那枚旁路是瞎的——但用例正控的形状（种下必红否则 `:264` 自红）意味着任何**新增**的旁路会在下一次跑用例时被正控网住，不依赖我的六发。
2. **M3 我只量了"吞掉失败句"这一形**，没量"失败句改词"（比如把 `首次配置建不出来` 打成别的字）。`:315` 的断言是子串匹配，改词会红——这一点是**读码推断**，不是我跑出来的读数；如果错了（断言其实匹配别的行），AC#5 的牙比 §3 写的窄。
3. **M4/M5 我只量了"整段删"与"编造命令"两形**，没量"只删一半"（比如只删 key 段留模型段）。`:374-386` 那组是逐名断言，读码推断任何一名缺失都会红，但我只跑过全删这一发。
4. **M6 我搬的位置是 `:390` 之前、`notify` 默认之后**（与 v1 V1 同位）；如果搬到 `assembleRuntime` 更深处（如 `NewStore` 成功之后），`:432` (a) 半是否还红我没量——读码推断会红（创建发生在任何 return 之前），但这是推断不是读数。
5. **`exit status 0xc0000135` 假象**：六发全部带 sherpa PATH，且每发都核过 `--- FAIL` 行与 rc 一致（rc=1 必有 `--- FAIL`、基线 rc=0 必全 PASS），无一发是"根本没跑"的形。如果我的判读错了，错在把某发 rc=1 但红句属于别票的情形算进来——复核：六发红句全部落在 `firstrun_198r2_test.go`／`firstrun_198_test.go`，无一例外。
6. **HEAD 漂移**：取数期间 HEAD 从 `315c6f5f` 挪到 `fa9a1fc6`（258-a1 骨架发，只动 `.scratch/wisp/issues/**`）。我核过三枚突变标的文件 md5 两锚逐字一致、`git diff HEAD -- cmd internal`＝0，所以六发读数在同码上取的。如果 `fa9a1fc6` 实际动了 `cmd/`（我没逐 commit 复核它的全部 pathspec），后果是 M1–M6 的读数混锚——复尺命令在 §3 头部，谁都能复跑。

## §5 判不动的地方

（由验收腿 `198-v2` 填）

1. **AC#4 回执内容的"真走得通"半**：`:363` 那枚守的是"点名真实存在的入口"（AST 能力检查）＋七个名词在场。但"`wisp secret set` 之后把回显的 blob 名补进 toml 这条路**用户照着走真能走通**"——这要真跑一轮 `wisp secret set`＋手改 toml＋再跑 `wisp run` 的端到端，本腿没跑（被验四格的射程是仪器有没有牙，不是端到端演示；端到端凭据 198-v1 的真件三发已给过 AC#1 侧）。**补得上**：任何腿拿真 exe 走一轮即可；本腿不判它。
2. **AC#5 的"只读盘／磁盘满"两形**：`:299` 用 `secrets` 位置放普通文件这一形挡 store（winsec "is not a directory"），与 v1 的 `C:\Windows\System32` 活量读数同族；磁盘满/只读卷两形本腿造不出（无盘配额手段），v1 §6-甲-4 也判了同一格。**不新增判语**。
3. **`:274` 那枚对 `applyDefaults` 反射器本身的变异**（比如把反射走查改成跳过某类 kind）：本腿没种——那要改 `internal/config/defaults.go` 的循环体本身，属于源侧更深的一形。读码推断 `:274` 会红（tag 期望侧与 applyDefaults 走查是两套独立走法），但没跑。
4. **时序读数**：六发里红用例的秒数（0.02s–0.12s）都是同机争用下的读数，只当"红了"的证据用，不当性能引。

## §6 交件判语

（由验收腿 `198-v2` 填；勾选框一枚未碰，翻勾归编排者）

| 格 | 牙有没有 | 凭据 |
|---|---|---|
| **① AC#2"改一枚默认值"的仪器** | **有牙** | M1（写侧种合法非默认值）⇒ `:221` 红、红句点名 `theme = 'dark' vs 'light'` 且字节差 2571→2572；M2（源侧改表）⇒ `:274` 红＋`:221` 继承红。两枚用例把"写侧漂"与"源侧漂"分开了（M1 下 `:274` 仍绿）。v1 §3-V3a"六枚全绿"的那枚洞，今天这发**实测闭合**。用例自带八类 kind 正控（`firstrun_198r2_test.go:248-266`），网住"期望侧瞎眼"的形。 |
| **② AC#5 失败支用例** | **有牙** | M3（吞掉失败句＝v1 V4 原形）⇒ `:299` 红，红句逐字点名 run.go:245-247。v1 §3-V4"全绿"的那枚洞，今天**实测闭合**。用例还钉了半份文件/临时残留/退码 2/成功回执不许误发四样（读码核过断言在位，本腿跑过的是失败句那一半）。 |
| **③ AC#4"去哪儿补"** | **有牙（仪器半）；内容半本腿不判** | M4（删两段）⇒ `:363` 红、五行逐名点名；M5（编造 `wisp catalog sync`）⇒ AST 能力检查红、逐字 `no cmdCatalog in package main`——**这条牙的形状是对的**（"编造入口点"这个失败模式第一次有仪器守）。回执内容是否"真走得通"见 §5-1，本腿不判。 |
| **④ J1 落点仪器收回本票** | **有牙** | M6（搬进 `assembleRuntime`＝v1 V1 原形）⇒ 本票自己的 `:432` 红（`:449` 落点半＋`:453` 回执半）＋既有 AST 钉 `firstrun_198_test.go:281` 同发红。**v1 §1-J1 暴露的"落点错了只被别人家的钉抓"已不复成立**：今天本票用例先红、且行为断言（不是只有 AST 读法）在位。 |
| **复核 §0–§2（编排者代提的三节）** | **复核通过，两笔小注** | §0：porcelain=0、两发产码 `1d8106d8`/`fa5593c0` 在史（`git log` 可证），与本腿起手锚 `315c6f5f` 一致——**对**。§1：四格修法与 `git show fa5593c0` diff 逐行核对一致（行号 `:221/:274/:299/:363/:432` 全部对得上）——**对**。§2：gofumpt 本腿复跑 `cmd/wisp` list-count=0（16:24:32，锚 `fa9a1fc6`）；d22scan 本腿复跑 rc=0、`clean - no D22 ban violations`、ban#8 `internal/`=482／`cmd/`=87，与编排者读数**逐字一致**；`go test` 三包合并 rc=0 那发本腿未整包重跑（336s 一发，非本腿职责），但本腿六发定向＋基线五枚全绿与之不矛盾——**对，test 一格采信不重跑**。小注一：§2 尺三分母写"ban#8 cmd/=87"没错但没写 `internal/`=482，本腿补齐。小注二：§0 说"台件 mutate.sh 没有执行残留"，本腿复量 porcelain=0 同读——**对**。 |

**总判**：四格全部实测有牙（六发种下必红、终值逐字在 §3；基线还原后五枚全绿）。v1 点名的三个洞（V3a 改值无牙／AC#4 没说去哪儿补／V4 吞句全绿）在这份产码上**逐一实测闭合**，J1 第四格从"别人家的钉"收回了本票。**本腿不翻任何勾**：AC#2／AC#4／AC#5 的勾选框留给编排者按本件 §3/§6 裁；本腿唯一的保留意见是 §5-1（AC#4 内容半的端到端没人走过）与 §5-2（磁盘满/只读两形无人量过）——这两条不挡翻勾，但该记在谁头上由编排者定。
