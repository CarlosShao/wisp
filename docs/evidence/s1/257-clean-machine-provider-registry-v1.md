# 257-v1 终裁件 — 票 257（干净机器上设置写不进去／形 ⓒ 首启回执）的对抗验收

> 腿：`257-v1`（非实现者终裁，`SPEC-12 §4.3` #1／`AGENTS.md` §0.3）。
> 票：`.scratch/wisp/issues/257-clean-machine-provider-registry-nil-blocks-writes.md`（AC 框一枚未碰）。
> 被裁的产码：`cmd/wisp/firstrun.go`（＋61 行，只进 stderr）；被裁的测试：
> `cmd/wisp/firstrun_257_test.go`（六枚）与 `cmd/wisp/firstrun_257_nonpreset_test.go`（两枚）。
> 本件只出判语，翻勾与删不删都归编排者。

## 0. 起手锚

| 项 | 读数（本腿现跑，逐字） |
|---|---|
| 进场时刻 | `2026-10-05 12:20:31 +0800`（`date` 自取） |
| 进场 HEAD | `5b638498`（派单给的锚是 `21bec8a1`→HEAD；这 15 分钟内链头已被别的腿推走，见 §6 用绝对 range） |
| 分支 | `dev` |
| `git status --porcelain -- cmd internal` | **0 行**＝`cmd/**`＋`internal/**` 此刻无别人的活在树上（12:20:51 现量） |
| `go env GOFLAGS` | 空串（两次现量 12:20:31／12:20:51 相同） |
| `git status --porcelain -- internal/agent/approval` | **0 行**＝`259-r2` 此刻在该包无未提交活（12:20:51） |
| 并发腿 | 派单告知 `259-r2` 在 `internal/agent/approval` 做注释级改动，会被 `cmd/wisp` 编译带上 ⇒ 本腿两发整包前后各取一次上面两行（§9、§10） |
| 本腿写面（授权） | `docs/evidence/s1/257-clean-machine-provider-registry-v1.md`＋`.scratch/wisp/probes/257/v1/**`；突变只许种进 `cmd/wisp/**`（⛔ 不种 `internal/agent/approval`／`internal/config`／`internal/panel`） |

## 1. 编队容量（起计时敏感腿前的三发，⛔ 任一 ≥70 % 就不抢机器）

`powershell -NoProfile -ExecutionPolicy Bypass -File .scratch/fleet-load.ps1` 连取三发（逐字）：

| 发 | 时刻 | 读数 |
|---|---|---|
| 1 | 12:20:31 | `CPU=58% MEM=56.4%` |
| 2 | 12:20:49 | `CPU=67% MEM=56.5%` |
| 3 | 12:20:51 | `CPU=60% MEM=56.6%` |

⇒ 三发均 <70 %，本腿放行去跑计时敏感那两枚；CPU 67 % 那一枚贴边，故 §9 那两发整包**各自起腿前再取一发**，任一 ≥70 % 就停下登记。

## 2. ★裁定一：两测试套是不是同一枚判据做了两遍、哪一套算交件

判语：**不是同一枚判据做了两遍——两枚各有独占射程，两套都算交件、必须同批留；本腿不删任何一枚**（删不删是编排者拿着这一格的动作）。另有一枚**只存在于证据件里的第三套设计**已被盘上推翻，具名见末段。

**尺一（重叠面，正向）**：本腿 12:36:48→12:36:56 现跑 `GOFLAGS= go test -count=1 -v -run 'TestTicket257R2|TestTicket198' ./cmd/wisp/`＝**rc=0，八枚 257 用例逐名 `--- PASS`**（`firstrun_257_test.go` 六枚＋`firstrun_257_nonpreset_test.go` 两枚），名册原件 `.scratch/wisp/probes/257/v1/baseline-257-family.log`。两枚共用的只有"干净机起步"这一枚**前提**（`firstrun_257_test.go:72-92` 与 `firstrun_257_nonpreset_test.go:48-63`，两份刻意不互相引用），断言体不重叠。

**尺二（决定性那一发，本腿自己种的突变 M2）**：12:35:50→12:36:01，只把 `firstrun.go:156` 那半句 `非预设名必须自己写 protocol，否则这份文件加载不过` 抹掉（needle 命中数断言＝恰好 1，脚本会拒绝"没改到东西的突变"）。读数：**只红一枚**＝`TestTicket257R2AC1ReceiptStatesTheNonPresetCondition`（`firstrun_257_nonpreset_test.go:77` 两行逐字），**六枚那一套全绿**。
⇒ 结论硬：那一枚 `firstrun_257_test.go` 在"回执不再教非预设名"的世界里**不会红**（它的 `t257TeachProviderRow` 只钉到 `[llm.providers.<名>]` 就停）。反发 M4／M5 只红六枚那一套里的 AC#2／AC#3 用例、nonpreset 两枚全绿 ⇒ **两套各自兜着对方兜不住的那一格**。

**尺三（真重叠的那几枚，逐名点名，⛔ 不删）**：
- 盘上明文反尺两枚同源：`firstrun_257_test.go:425-428` 与 `firstrun_257_nonpreset_test.go:145-148`（同一条 `api_key =`／`sk-` 反控）。
- 占位串断言两枚同源：`firstrun_257_test.go:406-418`（`dpapi:`／`env:` 计数相等）与 `firstrun_257_nonpreset_test.go:151-153`（同两串必在）。
⇒ 这三处是**同一判据做了两遍**，共 4 枚断言、占八枚用例里的两枚的一部分。可执行处置＝**两套都留**（合并的收益只有省 4 枚断言，代价是把两套唯一共用的 helper 变成跨文件依赖，正是 `A617` 那两条事故想要避开的形状）；要瘦身只瘦上面点名的那 4 枚，⛔ 不动 helper、不动用例枚数。

**弱针一枚（具名，不算重复但算脆）**：`firstrun_257_nonpreset_test.go:72` 那枚 needle `名字对上内置预设的，protocol 与 base_url 可以留空` **单独使用时没有牙**——票 198 的旧文案 `firstrun.go:116-117` 里有同一串。本腿 M2 只删后两句照样红，靠的就是后两句（并行腿在合并件 §9-c 里自己先招了这条）。⇒ 别把这一枚当牙用，删它不必要；要删由编排者裁。

**★第三套设计只在证据件里活着，盘上零枚（这条对翻勾者有后果）**：`git show 9c0d4c9a:.scratch/wisp/probes/257/r2/evidence.md` §3 逐字描述的四枚用例名——
`TestTicket257R2AC1ReceiptNamesTheThreeHandAddedThings`／`...CleanMachineRefusesAllSevenFields`／`...FollowingTheReceiptUnlocksAllSevenFields`／`...SecondRunSaysTheGuidanceOnlyOnce`——
本腿现量：`grep -r` 工作树 **0 命中**、`git grep HEAD` **0 命中**（四枚名逐枚数过，另有 §3 末行那枚省略号前缀 `TestTicket257R2AC0` 同样 0 命中）。那份设计还声称"⛔ 不硬编码 tag、用 `第 [123] 种拒因：` 正则从生产拒句反推"（§4.3），而盘上交付的那枚用的是**字面量**（`firstrun_257_test.go:50-54`）。
⇒ 判语：**这四枚名字一枚都不在树上**（工作树与 HEAD 各 `grep` 0 命中；它 §3 末行还以省略号写出第五个前缀 `TestTicket257R2AC0`，同样 0 命中）。
**谁按 `9c0d4c9a` 那份件翻 AC 框，就会给四枚不存在的用例记账。**本腿的处置＝不改它一个字（⛔ 临时件只建不删、也不改别人的件），只在这里具名；原件在 commit 里逐字可取，这条不需要谁来复核我的转述。

## 3. AC#1 干净机真跑（无 `config.toml` 起步＋ ⓒ 三样指引兑现）

判语：**成立（附两条具名边界，见本节末）。⛔ 没有一枚用例拿"手工先塞一份完整 config.toml"当前提。**

凭据（本腿现跑，两形）：

1. **包内那一形**：12:36:48→12:36:56 `GOFLAGS= go test -count=1 -v -run 'TestTicket257R2|TestTicket198' ./cmd/wisp/`＝**rc=0，八枚全 `--- PASS`**，其中 `firstrun_257_test.go:216` 那行 `t.Logf` 逐字复现＝
   `AC#1 clean machine: 7 refused writes by reason: map[第 2 种拒因：行不存在:7]`（七枚逐枚拒、每枚只报一种因、且逐枚点名 `[llm.providers.`）。
   起步态是**断言出来的**不是假设的：`firstrun_257_test.go:74-91` 与 `firstrun_257_nonpreset_test.go:50-62` 都在 `t.TempDir()` 上先 `os.Stat` 必须 `fs.ErrNotExist`、再走真入口 `run198`（＝`runTextTask`，`main.go:159` 那条 CLI 缝）、退码必须仍是 2、且 stderr 里必须有 `"新建默认配置"`，四道任一不满足就 `t.Fatalf`。
   数量两格由本腿复跑：`...ReceiptChainWalkUnlocksAllSevenFields`（照回执教的三样手加 ⇒ **7/7 接受且各报键路径**）与 `...ReceiptOnlyFirstItemUnlocksThreeOfSeven`（只加第一样 ⇒ **恰好 3/7**，多一枚也红）。
2. **真进程那一形（本腿自己搭的干净机，⛔ 不入库的读数不作数，原件全留在 `.scratch/wisp/probes/257/v1/probe/`）**：
   `go build -o …/probe/wisp.exe ./cmd/wisp`（12:37:05→12:37:12 rc=0）→ 同目录放一枚 `portable.txt`（`resolveDataDir` 认它，`doctor.go:247-255`）→ 该目录原本无 `data-dev/` ⇒ 真·干净机 →
   `./wisp.exe run "总结一下 这份笔记"`：**rc=2**、stderr **5,674 字节**、stdout 只有版本行 119 字节。stderr 里逐串命中次数＝`第一样＝一节 [llm.providers.<名>]` 1、`第 2 种拒因` 1、`第 3 种拒因` 1、`第三样＝就地填已有的 [llm.roles.chat]` 1（`real-run-stderr.txt`）。
   生产建出的 `data-dev/config.toml`＝**2,571 字节**、`grep -c "llm.providers"` ＝ **0**（与票 198 那枚逐字节尺的 `want 2571` 同数）。
   第二发（12:38:13，`real-run2-stderr.txt`）＝rc=2、`新建默认配置` 命中 0、指引命中 0 ⇒ 首启专属这一格在真进程里也复现。

**边界两条（⛔ 不是本腿划掉的，是交付面本来就没覆盖的）**：
- 名册第 7 枚 `provider_credential` 的**值腿**（`StoreCredential`→DPAPI 真密封）两枚测试都没走：面板那条路本身就不给它值（`cmd/wisp/panel_config_store.go:218-220` 逐字回"凭据值不走这条写入路径…必须走 StoreCredential"），测试改走它的**引用侧** `SetProviderAPIKeyRef`（`firstrun_257_test.go:126-131` 注释自己写明"枚数 7、腿数 6"）。⇒ 真机 DPAPI 那一格仍欠 `docs/evidence/s1/` 一行，归编排者（与 `A617` §1 残余条同一格）。
- 七枚**在设置页面上的长相**（信封键、回执串）不在本票写面（票 §8 边界③：要从界面建行＝单开票）。本腿只在 `internal/panel/config_handlers.go:501`／`:381-394` 读到"错误串会被写进审计行／回执"这一形状，⛔ 不据此判页面那一格。

## 4. AC#2 三种拒因三句话（含定向突变）

判语：**成立**（三句各带各自的补救，折成一句必红）。

凭据：
- 生产侧三枚 tag 现读：`internal/config/settings.go:78-80`（`refusalFileMissing`／`refusalRowMissing`／`refusalInvalid`），它们各自的真身分支 `:242`（role 缺行改判第 2 种）／`:293`（预写 `validate()` 门＝第 3 种）／`:367`＋`:379`（第 2 种配 `guidanceModelRow`／`guidanceProviderRow`），文件没建那一支在 `internal/config/loader.go`。回执三句（`cmd/wisp/firstrun.go:162-171`）与三枚 tag **同串**，且第 1 句把本机真路径 `Fprintf` 进去（不是模板话）。
- **本腿自己种的突变 M4**（12:36:08→12:36:17）：把回执里第 2、第 3 句**折成同一句** `"  配置未生效——去这份文件里看看，"`（两枚 needle 各命中 1 次）。指名用例必须红——红句逐字（`mut-M4-…log`）：
  `firstrun_257_test.go:328: AC#2 RED: "第 2 种拒因：行不存在" appears on 0 receipt lines, want exactly 1 (three reasons, three sentences):`
  `firstrun_257_test.go:328: AC#2 RED: "第 3 种拒因：校验不过" appears on 0 receipt lines, want exactly 1 (three reasons, three sentences):`
  `firstrun_257_test.go:345: AC#2 RED: the folded verdict "配置未生效" appears 3 times, want exactly 1 (its own prohibition): …`
  `firstrun_257_test.go:349: AC#2 RED: "配置未生效" is used as a verdict rather than named as the forbidden collapse:`（两行，逐枚被折的那两句各一行）
  rc=1、`--- FAIL: TestTicket257R2AC2ThreeRefusalsStayThreeSentences`。还原后 md5 回基线。
- 第二发（同一 AC 的另一格，票 §8 边界①＋`A560`）：`TestTicket257R2AC2EffectTimingSaysThreeProcessShapes` 本腿复跑绿，它钉的是"什么时候算用上"也分三形状（`每 1s 重读一次 config.toml 的看门狗`／`根本没有会重读盘的东西`／`自己明说不带轮询`），⛔ 折成"重启就好"同样被 `配置未生效|重启就好` 只许出现一次的计数尺拦着（并行腿的 M4 那发已在合并件 §5 记过同形红）。

## 5. AC#3 凭据面（不新增回显、不落明文）

判语：**成立**（附一条**准确性 nit**，不构成本格失败，具名见下）。

凭据（本腿现跑的四把尺）：
- `StoreCredential` 那一侧**有没有被扩**：`git diff -G"StoreCredential" --name-only 21bec8a1..HEAD` 只点回 `.scratch/wisp/probes/257/r2/{evidence.md,snapshot-firstrun_257_test.go}`、`cmd/wisp/firstrun_257_test.go`、`docs/reports/pending-and-issues.md`——**产码零枚**；同范围内 `cmd/wisp/panel_config_store.go`＋`internal/panel/config_handlers.go`＋`internal/panel/bridge.go` 三枚文件 `git diff --name-only` 命中 **0**。⇒ 没有新增任何"把值回显给页面"的方法。
- 本票那 61 行的形状：`grep -c '^func ' cmd/wisp/firstrun.go` ＝ **1**（仍只 `ensureFirstRunConfig` 一枚，签名 `(dataDir string, stderr io.Writer) (bool, error)` 未漂，静态钉在 `firstrun_257_test.go:434-439`）；新增行里 `go func(|filepath.Clean|filepath.Abs|time.Since|time.Now` 命中 **0**；`DEFERRED` 标记新增 **0**（⛔ 不动 `SPEC-12 §5` 那张表的 1:1）。
- 回执里**有没有键值形状**：`firstrun_257_test.go:404-428` 那把"总数==占位数"尺本腿复跑绿；`grep` 本腿自数 `dpapi:`／`env:` 只以 `dpapi:<blob 名>`／`env:<环境变量名>` 出现。种下去会红：本腿 **M5**（12:36:17→12:36:26，往凭据那句里塞 `api_key = "sk-t257v1mutantvalue"`）⇒ 唯一红枚 `TestTicket257R2AC3CredentialSurfaceUntouched`，红句逐字：
  `firstrun_257_test.go:420: AC#3 RED: the receipt shows a plaintext key assignment shape:`
  `firstrun_257_test.go:423: AC#3 RED: the receipt carries a vendor key-shaped literal:`
- 对话与输出零凭据值：真进程那一形（§3 第 2 条）的 stdout 只有版本行 119 字节、stderr 5,674 字节里 `api_key =`／`sk-` 命中 **0**；生成的 `config.toml` 里 `api_key =` 命中 **0**；本件与本腿 commit message 内无任何值形状串（只有 `env:T257_*` 这类永不置值的变量名与 `dpapi:<blob 名>` 占位）。

**准确性 nit（新文案里的一句与盘上面板读数不完全一致，⛔ 不是值泄漏）**：
`cmd/wisp/firstrun.go:183` 写"这一页任何时候都不回显密钥的值，**只说已录入还是没录入**"；而 `internal/panel/config_handlers.go:459-461`（`renderSettingsView`）在设置页读数里逐字印 `引用 %s`（`refOrNone(p.APIKeyRef)`）。
⇒ 前半句（不回显**值**）为真：印的是引用名；后半句（"只说已录入/未录入"）与盘上那句不完全一致——那一页还会把引用名念出来。这是**票 198/248 既有形状**（该文件在本范围内 0 改动），新文案把它说窄了。判语：AC#3 仍成立；这一句要不要改成"只念引用名、不念值"归编排者（改文案＝动那 61 行，⛔ 本腿不改产码一字）。

## 6. AC#4 越界（`21bec8a1..HEAD` 禁区文件名命中数）

判语：**成立**（本票的落在范围内＝零越界；范围内那五枚 `internal/config/**` 命中逐笔归票 267，不是本票）。

凭据（本腿现跑，`git diff --name-only 21bec8a1..HEAD` 全量存 `.scratch/wisp/probes/257/v1/ac4-diff-names.txt`）：

| 禁区模式 | 范围内命中枚数 |
|---|---|
| `^frontend/` | 0 |
| `^design/` | 0 |
| `docs/PLAN.md` | 0 |
| `^docs/specs/` | 0 |
| `internal/observe/thresholds.go` | 0 |
| `golden` | 0 |
| `allowlist.txt` | 0 |
| `^internal/config/` | 5（`schema.go`／`unwired.go`／`validate.go`／`validate_267_test.go`／`validate_test.go`） |

- 那五枚逐笔 `git log 21bec8a1..HEAD -- <file>` 全指向 `37f8e5c6`／`ec6a47a8`／`873c3063`（subject 逐字起于 `ticket 267 r1: band [risk].confirm_timeout_sec at load (shape a)`）＝**票 267 的写面**，与本票无关。
- **本票各腿自己的 commit 过同一把尺＝零命中**（逐笔 `git show --name-only --format=''`，命中数枚枚为 0）：
  `9c0d4c9a`(5 枚)／`3f0c4fff`(9 枚)／`06b66e66`(1 枚)／`7c644bb0`(38 枚)／`ee9cf2eb`(3 枚)／`e940334c`(6 枚)／`62325dff`(16 枚)。
- ⚠ 一处**范围外但必须具名**的落点差（不是本腿新发现，是它自己登记的）：`8e443ed7`（10-04 18:59，`257-r1b`）改过
  `internal/config/loader.go`＋`internal/config/settings.go`＋新增 `internal/config/settings_257_test.go`，
  而票 §8-4 给 `257-r1` 的写面写的是「⛔ 不碰 `internal/config`」。那一笔早于本范围锚 `21bec8a1`（10-04 20:41），
  故不计入上面那张表；它自己的证据件 §1 标题逐字就是「写面授权与落点差」，编排者已在 `A606` 收账。
  ⇒ 判语：AC#4 在**给定范围**内成立；那一格落点差归编排者确认（⛔ 本腿不替它判作废，也不据它翻红）。
- 尺的形状：⛔ 后两枚（`golden`／`allowlist.txt` 所在的 `tools/d22scan/`）与 `internal/config/**` 三枚包本腿**连文件内容都没读**，只过文件名。

## 7. ★裁定三：那 61 行有没有把指引漏进生成的 `config.toml`（实弹突变）

判语：**没有漏**；而且**"漏了就响亮红"这一雷是真会响的**——本腿自己种下去证过一遍（不是引用并行腿的 M4 读数）。

凭据（两把尺同发）：
- **今天盘上那一形**：真进程干净机建出的 `data-dev/config.toml`＝2,571 字节、`grep -c "llm.providers"`＝**0**；`grep -n "^\["` 尾部节名＝`[cost]`／`[models]`／`[observe]`／`[observe.roll]`／`[plugins]`，没有任何 provider 节被"顺手建出来"（原件 `.scratch/wisp/probes/257/v1/probe/`）。包内那一形同向：`firstrun_257_test.go:165-171` 逐枚读 `os.ReadFile(cfgPath)` 反控 `[llm.providers`（⛔ 它读的是文件，不是 stderr，这正是"教"与"建"之间那道界）。
- **本腿种的实弹 M3**（12:36:01→12:36:08）：在 `SaveFile` 成功后加三行，把一段 `[llm.providers.t257v1leak]` 追加进**生成的文件**（`os.ReadFile`＋`os.WriteFile`＋`append(raw, …)`）。结果 **rc=1、红 6 枚**，逐字红句（`mut-M3-guidance-leaks-into-generated-file.log`）：
  `firstrun_198_test.go:217: the created file carries the dynamic table "[llm.providers" anyway - that would be an invented entry, not a default`
  `firstrun_198r2_test.go:239: AC#2 RED: the config.toml first-run wrote is not the schema's `default` tags rendered: first difference at line 179 (byte 2571): want "", got ""; total bytes want 2571 got 2599`
  `firstrun_257_test.go:170: AC#1 RED: the guidance leaked into the generated config.toml, which is invented-default territory and ticket 198's forbidden shape:`
  `firstrun_257_test.go:249: AC#1 RED: the hand-add shape the receipt teaches does not even load: config: config.toml: llm.providers.t257v1leak.protocol must be set explicitly for a non-preset provider …`
  另两枚同发红＝`TestTicket257R2AC1ReceiptOnlyFirstItemUnlocksThreeOfSeven`／`TestTicket257R2AC1NonPresetRowNeedsProtocolBeforeItLoads`（都因那份文件自己加载不过）。
  ⇒ 三把尺同时响：票 198 的**名册反控**、票 198 的**逐字节尺**（2571→2599）、本票 AC#1 自己的**泄漏反控**。`A543` 记的那条最硬的雷（"ⓐ 空壳行会让首建文件自己加载不过、进程起不来"）在这一发里被复现成红句，不是推测。
- 还原证明：五发突变各归 `md5sum cmd/wisp/firstrun.go`＝`0491339282492f2cabdbf5be576c8a57`（＝`git cat-file blob HEAD:cmd/wisp/firstrun.go` 同一值），脚本每发跑完立刻还原、末尾再双读（`FINAL md5 … equals baseline: True`）。

## 8. ★裁定二：正向断言（"提示真到达 stderr/stdout"）是不是恒真

判语：**不恒真——它有牙**；⚠ 另要更正一处口径：**那一格真到达的通道是 stderr，不是 stdout**，两枚测试测的是 `runTextTask` 注入的那枚 `io.Writer`，"生产把它绑到 `os.Stderr`"这一步**没有任何用例钉着**（本腿用真进程读数顶了这一格，⛔ 它不该被当成入库的钉）。

凭据（攻法＝把"指引"从**到达**降级为**存在于二进制但谁也没收到**）：
- 本腿 **M1**（12:35:39→12:35:50）：只把 `firstrun.go:151` 那枚 `fmt.Fprint(stderr, …)` 的目标改成 `io.Discard`——三样指引的**字面串一个都没动**（仍在源码与二进制里），只是不写给操作者。⇒ **rc=1、红 3 枚**：
  `TestTicket257R2AC1CleanMachineReceiptTeachesTheWalkableChain`（逐条 `AC#1 RED: the first-run receipt never says "第一样＝一节 [llm.providers.<名>]"` 等五条）
  `TestTicket257R2AC1ReceiptChainWalkUnlocksAllSevenFields`（`the receipt stopped teaching the in-place fill`）
  `TestTicket257R2AC1ReceiptStatesTheNonPresetCondition`（两行 `never states "非预设名必须自己写 protocol"`／`"否则这份文件加载不过"`）
  ⇒ 在"回执里根本没有指引"的世界里这几枚**必红**，恒真性攻不动。
- 前提也不是恒真：`t257CleanMachine`／`t257npCleanMachine` 都在拿回执之前先 `t.Fatalf` 断言"文件此前不存在＋退码仍 2＋stderr 里有 `新建默认配置`"（`firstrun_257_test.go:76-91`／`firstrun_257_nonpreset_test.go:52-61`）——若 `ensureFirstRunConfig` 一声不响地建了文件，用例在**读到回执之前**就红，不会把"没有回执"读成"回执通过"。
- **通道那一格的真读数（本腿现跑，⛔ 不入库的读数不作数）**：`main.go:159-163` 逐字把 `stdout: os.Stdout, stderr: os.Stderr` 交给 `runTextTask`；真进程那一发 stdout＝**119 字节**（只有版本行，无指引），stderr＝**5,674 字节**（三样指引＋三因＋三形状全在这里）。
  ⇒ 若按票面/派单那句"真到达 stdout"逐字裁：**那一格不成立，且本来就不该成立**——文案落的是 stderr（票 198 结案的同族形状：回执走 stderr、退码 2 一字未放宽）。本腿按"到达用户可见通道"裁＝**成立**，并具名登记"stdout 是措辞不符，不是交付缺失"，翻勾时别把它写成 stdout。
- ⚠ 残余（把它写成结论就是假话）：**没有任何用例从真子进程口径测过这枚通道**（两枚测试都在同进程里传 buffer）。本腿那发真进程读数顶了今天，⛔ 不可重放为门禁；要长期钉住需一发行文＋一枚 `exec` 级用例，归编排者决定要不要新开一格。

## 9. 那枚新增红＋那枚命名跳的复跑两形与归因

判语：**票 257 的改动造成的新增红＝0 枚**；实现者那一发的红与跳＝**机器争用**（本腿低载复跑不复现）；本腿自己这一发另出**一枚新的带载红**（`TestAC14GoSideEvalPushReachesThePage`），同样与本票写面无因果路径，具名见下。⛔ 阈值一字未动、⛔ 无 `t.Skip`、⛔ 没把 SKIP 读成通过。

**容量（每发起腿前现取，owner 那条 ≥70 % 不抢机器）**：起手三发（§1）＝12:20:31 `58%`／12:20:49 `67%`／12:20:51 `60%`；整包起腿前 12:21:49 `43%`；整包收腿后 12:30:36 `32%`；单发前 12:30:58 `35%`；定向隔离前 12:31:18 `52%`；count=3 前 12:31:43 `63%`。⛔ 全部 <70 %，无一发起腿被本腿自己拒掉。

| 发 | 时刻（起止到秒） | 命令形状 | 读数 |
|---|---|---|---|
| 整包（形 2） | 12:22:3x→12:29:07（`PKG1 END rc=1`） | `GOFLAGS= go test -count=1 -v ./cmd/wisp/`（带 dll 注入） | `RUN 331／top PASS 234／top FAIL 1／top SKIP 0／子 PASS 96`（234＋1＋96＝331 闭合）；红＝`TestAC14GoSideEvalPushReachesThePage (0.61s)`；**lifecycle PASS 1.32s（cold 917.218 ms）／latency PASS（n=1，cold P50=P95=917.218 <1500、hot 129.544 <200）** |
| `-count=1` 单发（形 1） | 12:30:58→12:31:09 rc=0 | `-run 'TestPanelHostRealWindowHopAndLifecycle\|TestPanelHostLatencyPercentilesAC2'` | 两枚全 PASS：cold `980.227 ms`／hot `71.059 ms`，latency `n=1` 有样本、**零 SKIP** |
| `-count=3` 附加形 | 12:31:43→12:31:54 rc=0 | 同上两枚 ×3 | 三迭代全 PASS，cold `960.457 / 757.303 / 802.172 ms`，**没有一枚越过 1500 ms**（实现者那发是 1031→2107.9 ms） |
| 定向隔离（归因本腿那枚红） | 12:31:18→12:31:32 rc=0 | `-run 'TestAC14GoSideEvalPushReachesThePage\|TestAC1ResidentLegInstallsItsLogListenerOnDisk'` | 两枚都 **PASS**（AC#14 页面上自己的字 `"PUSHED-33R5-OK"`）⇒ 本腿整包里那枚红是**包内顺序／带载**那一形，不是产码 |

**逐名作差**（本腿整包 vs 实现者整包 `.scratch/wisp/probes/257/r2/full-cmd-wisp.log`，名册文件＝`roster-theirs.txt`／`roster-mine.txt`）：

- 顶层结果枚数 **235 vs 235**，名集 `diff -q` **完全相同**（无丢名、无新增名）。
- 只差三枚的**颜色**：`TestPanelHostRealWindowHopAndLifecycle`（他 FAIL→我 PASS）、`TestPanelHostLatencyPercentilesAC2`（他 SKIP→我 PASS 且真带 `n=1` 样本）、`TestAC14GoSideEvalPushReachesThePage`（他 PASS→我 FAIL，隔离复跑又 PASS）。

**归因（三格分别说，⛔ 不含糊）**：
1. **票 257 造成的新增红＝0 枚**。本票整条交付只有 `cmd/wisp/firstrun.go` 的 stderr 文案＋两枚同包测试件；上面三枚变色全在真窗体族（`panel_host_*`／`panel_resident_windows_*`），与本票写面无因果路径。
2. **机器争用＝这三枚**：实现者 11:30→11:39 那发的红与跳（冷启 619→1031→2107.9 ms、`machine-wide msedgewebview2` 19~24），本腿在 32~63 % 水位复跑同一对（单发＋count=3 共 4 个冷启样本 `917/980/960/757/802 ms`）**无一越过 1500 ms**，且那枚命名 SKIP 不再出现 ⇒ 判**带载**，不是回归。⚠ 那一格的关键读法我照实现者的口径保留：**SKIP＝什么都没测**（它自己的句子逐字 `no cold/hot sample recorded in this process`），所以我这发必须是"有样本的 PASS"才算数——`n=1`、cold/hot 两数都在日志里。
3. **无法归因／具名说**：本腿整包那枚 `TestAC14GoSideEvalPushReachesThePage` 红（`title=""`，want `"PUSHED-33R5-OK"`）在隔离复跑里绿 ⇒ 归"包内顺序＋带载"那一族，与 `A617` 里已在册的同名红同一形；本腿不新立案、不划掉，交编排者按票 33 那侧记账。

**⚠ 污染面申报（那两发整包的前后读数，逐字）**：
- 起腿前 12:21:49：`git status --porcelain -- internal/agent/approval`＝**0 行**、`go env GOFLAGS`＝空、HEAD `65142d9a`。
- 收腿后 12:30:36：`internal/agent/approval`＝` M approval.go`＋` M queue.go`（两枚，`git diff --stat`＝109 增／13 删，读片段确为注释级）、`GOFLAGS`＝空、HEAD `549cb218`。
- ⇒ **`259-r2` 在我整包窗口中段落了笔**（正是派单预告的那种撞法）。本腿仍收下这一发读数的理由具名：那两枚文件在 `cmd/wisp` 的编译闭包里，但变色那三枚的判据不读 `approval`；三枚变色里 lifecycle／latency 两枚另有**低载隔离复跑**（12:30:58 起腿前后各取一次 status，见 `iso1-status-approval-before.txt`）顶着，AC#14 那枚另有 12:31:18 的隔离复跑顶着。**若编排者认为整包那一发因污染不作数，需要复跑的是"包内顺序能否复现 AC#14 的红"这一格，不是票 257 的六格。**

## 10. 门禁四门（本腿自跑，⛔ 不引用实现腿的读数）

| 门 | 时刻（+08） | 命令（逐字） | 读数 |
|---|---|---|---|
| D22 静态扫描 | 12:33:14 起→12:33:44 完 | `sh scripts/d22scan.sh` | **rc=0**，末行逐字 `d22scan: clean - no D22 ban violations`；覆盖面 `bans #1-5 internal/=228、cmd/=38、ban #6 frontend/=85、ban #7 internal/tools/=23、ban #8 design/=39、frontend/=85、internal/=512、cmd/=103`（`gate-d22scan.log`） |
| 路径长度帽 | 12:33:44 | `sh scripts/check-path-length-budget.sh --with-self-test` | **rc=0**，`VERDICT GREEN`；`tracked=5866／over-budget=57／covered by roster=57／not in roster=0`；`longest=180 chars relative`（`252-…` 那枚）；自检那一行逐字 `control 2/3 ok - the planted over-budget tracked path is rejected and named verbatim`（`gate-pathlen.log`） |
| `go vet` | 12:34:00 | `go vet ./cmd/wisp/` | **rc=0**，日志 **0 字节**（`gate-vet.log`）。⚠ 空读数这一发不作"跑到了"的证据，所以同发跑了 `go test` 与下面那把 gofumpt 的版本尺；本腿全程没有遇到过 11:26:51 那一形的编译红 |
| gofumpt 点名 | 12:34:02 | `"D:/work/base/gopath/bin/gofumpt.exe" -l cmd/wisp` | **rc=0，只列一枚 `cmd\wisp\models.go`**（预存 CRLF，不归本票，⛔ 未顺手格式化）；本票三枚（`firstrun.go`／`firstrun_257_test.go`／`firstrun_257_nonpreset_test.go`）**零命中**。尺本身活着的证明＝同把 `-version` 回 `v0.12.0 (go1.27.1)` |

⛔ 没用 `~/GOPATH/bin/gofumpt` 那把瞎尺；`go env GOPATH` 本腿现量＝`D:\work\base\gopath`（与派单一致）。

## 11. 突变名册（每发起止时刻到秒，编排者拿它作差）

⛔ 五发全种在 `cmd/wisp/**` 唯一一枚 `firstrun.go`；⛔ 未种 `internal/agent/approval`／`internal/config`／`internal/panel`。台件＝`.scratch/wisp/probes/257/v1/mutate257v1.py`（每枚 needle **要求恰好命中 1 次**，不符就 `ABORT` 不跑测试＝不会产"没改到东西的假牙"；每发跑完立刻还原并 `md5` 双读）。

| 号 | 起点→终点 | 种法（一句话） | 应红格 | 实读 |
|---|---|---|---|---|
| M1 | 12:35:39→12:35:50 | 指引那枚 `Fprint` 的目标改 `io.Discard`（字串不动） | §8 恒真性 | rc=1 红 3 枚（两套各红，见 §8 逐字红句） |
| M2 | 12:35:50→12:36:01 | 抹掉 `非预设名必须自己写 protocol，否则这份文件加载不过` | §2 两套射程 | rc=1 红 **1 枚**＝`…ReceiptStatesTheNonPresetCondition`，六枚那一套全绿 |
| M3 | 12:36:01→12:36:08 | 把 `[llm.providers.t257v1leak]` 追加进**生成的 config.toml** | §7 最硬那枚雷 | rc=1 红 **6 枚**（含票 198 名册反控＋逐字节尺 2571→2599） |
| M4 | 12:36:08→12:36:17 | 回执第 2、3 句折成同一句「配置未生效——去这份文件里看看」 | §4 AC#2 | rc=1 红 1 枚＝`…AC2ThreeRefusalsStayThreeSentences`（四条红句） |
| M5 | 12:36:17→12:36:26 | 凭据那句里回显 `api_key = "sk-t257v1mutantvalue"` | §5 AC#3 | rc=1 红 1 枚＝`…AC3CredentialSurfaceUntouched`（两条红句） |

- 基线读数（同缝、未突变）：12:36:48→12:36:56 rc=0，八枚 257 全 PASS＋十一枚票 198 全 PASS。
- 还原凭据：`BASELINE md5 0491339282492f2cabdbf5be576c8a57 worktree-equals-head=True` 起手即证；每发 `restored_md5_ok=True`；收尾 `FINAL md5 … equals baseline: True`。
- 窗口内的污染面：这五发跑在 12:35:39→12:36:26，`internal/agent/approval` 在 12:30:36 已脏（`259-r2` 注释级），本腿没碰它；五发的 `-run` 只挑 `TestTicket257R2|TestTicket198`，不含计时敏感用例。

## 12. 七项对抗检查（`SPEC-10 §8`，逐格本腿自跑／自读）

1. **stub 扫描**：⛔ 无 mock 顶真件。两枚测试都走**真入口** `runTextTask`（`main.go:159` 那条缝，README「seams」列的合法注入面之一）与**真加载器** `config.NewManager(path, nil)`（`nil` 正是产线三处的形状，`A560`）；`firstrun.go` 的 61 行里 `grep` 无 `TODO`／返回 nil 假成功形状，`ensureFirstRunConfig` 仍只在真建成功时打印。
2. **DEFERRED 双向核对**：本票新增行里 `DEFERRED` 标记 **0 枚**（`git diff -U0 21bec8a1..HEAD -- cmd/wisp/firstrun.go | grep -c "^+.*DEFERRED"`＝0）⇒ 没有新增登记项，`SPEC-12 §5` 那张表不被本票牵动。⚠ 顺带复认一处**既有**在册缺陷（不属本票、不新立）：AGENTS.md 与记忆里都写着"表 13 行／代码 0 枚标记＝空转规矩"，本腿没有尺可核。
3. **SLO 测量完整性**：`internal/observe/thresholds.go` 在范围内 **0 命中**（§6 那张表），D32 那两枚预算（冷启 1500 ms／热 200 ms）本腿**一字未动**、读数取自用例自己的 `t.Log`（§9 那四发冷启 917/980/960/757/802 ms）；⛔ 无"放宽断言换绿"，⛔ 无 `t.Skip`，⛔ 未把 SKIP 读成通过（§9 明写"跳＝什么都没测"）。
4. **契约偏离检查**（D34 工具表／C19 R1–R9／D43 转移表）：本票没碰任何一张——`internal/panel` 三枚契约面文件（`config_handlers.go`／`bridge.go`／`panel_config_store.go`）范围内 0 改动，名册仍七枚（`config_handlers.go:58-68` 与 `firstrun_257_test.go:107-132` 那七枚**逐名相等**，本腿对过）。ⓒ 那形的"不新增 C17 面字段"兑现＝ⓑ 那形没被偷偷做。
5. **D22 七条新增禁令静态扫描**：见 §10 第一行（rc=0 clean）＋§5 第二把尺（新增行里 `go func(|filepath.Clean|filepath.Abs|time.Since|time.Now` 命中 0）＋`firstrun.go` 内函数枚数仍 1。⛔ 明文 API 密钥：五发之外全程零值形状（§5）。
6. **场景↔切片双向核对**：本格超出单票射程——`SPEC-00 §4` 那张矩阵要按切片核，归编排者（§13 第 4 行具名）。本腿只补一句：票 257 的**用户可见后果**（"机主第一次配模型能不能走通"）今天只有**静态可走通**这一种证明（手加三样 ⇒ 7/7），没有人替机主真在编辑器里打过那三样。
7. **失败预演＋缺口审计**：审计者＝本腿（≠ 实现腿 `257-r1*`／`257-r2`／`257-r2b`）。缺口逐条落在 §3 边界两条、§5 nit 一条、§7/§8 各一条残余、§13 全表；⛔ 没有一格被写成"应该没问题"。

## 13. 我攻不动的地方／判不动的格子（具名归口）

| # | 格子 | 为什么够不着 | 归口 |
|---|---|---|---|
| N1 | `provider_credential` 的**值腿**（`StoreCredential`→DPAPI 真密封＋"只写不回显"在真机上的长相） | ⛔ AC#3 不许为测试去写真 blob；两枚测试都只走引用侧；本腿同样不写（没种这一发） | 编排者：`docs/evidence/s1/` 单行真机裁决（同 `A617` §1 残余第一条） |
| N2 | 三句拒因与三样指引**在设置页上的完整长相** | ⛔ `internal/panel/**` 不是本票写面，且本腿派单禁改；本腿只读到 `config_handlers.go:381-394`／`:501` 的形状，量不到渲染终串 | 编排者（票 §8 边界③已预告"要界面建行＝单开票"） |
| N3 | 回执通道（stderr）被真子进程口径**长期钉住** | 今天没有任何 `exec` 级用例；本腿那发真进程读数是手工的、不可当门禁 | 编排者：要不要新开一格（§8 末行） |
| N4 | 两枚测试里三枚 tag 用**字面量**（`firstrun_257_test.go:50-54`）这件事本身好不好 | 有意为之（并行腿 §7A-N7 给了理由：换成读常量就落回 257-r1b 的 M1 盲区），本腿认同其取舍但⛔ 无权改它一根针 | 编排者知悉即可；本腿判"脆性可接受" |
| N5 | `TestTicket257R2AC1…` 那两枚 `AC1…` 用例与票面 AC 编号的**一对一映射表**（README 规则 6 要求的裁决表按 AC 编号） | 票面 AC 框只有 4 枚（AC#0 已勾），本件已按 AC#1–AC#4 出判语；`AC#0` 的形状裁在 `A543`，⛔ 本腿不重裁契约 | 编排者翻勾时按 §3／§4／§5／§6 四格即可 |
| N6 | **机主人因层面**："教了三样"不等于"机主会照三样打" | 只能靠 e2e／自用期读数，静态与包内尺顶不了 | 编排者（与 §12 第 6 项同一格） |
| N7 | 我攻不动的最后一枚：五发突变之后，`A617` 那条"同一轮派两枚写腿"的**根因**不在这棵树的任何一格判据里 | 那是派单纪律，不是本票交付物 | 编排者自己已记账，本腿不代劳 |

## 14. 交件判语与 commit 链

**四条 AC 的判语（逐格一句，⛔ 框归编排者翻）**：
- **AC#1 成立**（附 §3 那两条具名边界）：干净机两形都真跑过——包内八枚 rc=0 复现 `map[第 2 种拒因：行不存在:7]`、真进程那发 stderr 5,674 字节里三样指引各到一次、生成的文件 2,571 字节零 `[llm.providers`。
- **AC#2 成立**：三因三句各带互不重叠的补救，本腿 M4 折两句 ⇒ 指名用例红，四条红句逐字在 §4。
- **AC#3 成立**（附一条准确性 nit）：范围内 `StoreCredential` 侧与三枚面板契约文件零改动、`firstrun.go` 仍一枚函数，M5 回显值 ⇒ 指名用例红。
- **AC#4 成立**：`21bec8a1..HEAD` 八枚禁区模式里，除 `internal/config/**` 那 5 枚**全部归票 267** 外其余全 0；本票各腿自己的 commit 枚枚 0 命中；一处**范围外**的 `257-r1b` 落点差具名在 §6 末段。
- **★三格加裁**：两测试套＝**两套都留、同批算交件**（M2 决定性）；正向断言＝**不恒真、有牙**（M1），通道措辞需从 stdout 更正为 stderr；61 行漏进文件＝**没漏，且漏了必响**（M3，三把尺同响）。
- **那枚红＋那枚跳＝本腿复跑判为机器争用**：票 257 新增红 **0 枚**；本腿整包另出一枚在册同名红（AC#14），隔离复跑绿，归票 33 那一族；⚠ 整包那一发的窗口里 `259-r2` 落了笔（前后 status 已报，§9 末段）。

**本腿 commit 链（只 commit、⛔ 零 push、笔笔 `git commit -F msg -- 枚枚点名 pathspec`）**：
| 笔 | 时刻 | 内容 |
|---|---|---|
| `65142d9a` | 12:21:29 | 骨架（§0 锚＋§1 容量三发现量） |
| `0793e35e` | 12:42:07 | §2 两测试套裁定＋§3/§4/§5 三格判语＋两枚 msg 台件 |
| `3862e0e2` | 12:45:43 | §6/§7/§8/§9/§10/§11/§12/§13/§14 收满＋本腿全部文本台件（37 枚，`--name-status` 只有 `A`／`M`，**零枚 `D`**） |
| 本笔 | 见 `git log -1` | 只把上一行那枚号补进这张表（交件后追写，⛔ 不改写已落的那三笔） |

收尾三把尺（本腿现量，⛔ 不是转述）：`git status --porcelain` 对票面文件与 `docs/reports/pending-and-issues.md` 各 **0 行**；`git diff` 对 `cmd/wisp/firstrun.go`＋两枚测试件 **0 行**（`md5sum`＝`0491339282492f2cabdbf5be576c8a57`＝起手基线）；本件 `wc -l`／`wc -c`／占位尺见本节末。

**纪律自证**：⛔ 票面 AC 框（`- [ ]` 四枚）与 `docs/reports/pending-and-issues.md` 一字未改（两枚文件都不在本腿任何 pathspec 里）；⛔ 产码一字未改（五发突变全部还原、`md5` 五取相等）；⛔ 未删任何一枚测试件（两套件与 `9c0d4c9a` 原件都在原位，覆盖过的、没覆盖的都没动）；⛔ 无 `add -A`／`add .`／`--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`；临时件只建不删（本腿台件全留在 `.scratch/wisp/probes/257/v1/`，含 `probe/` 那枚 31 MB 的 `wisp.exe`，⛔ 未 commit 进仓，清点见本节）；未在仓库内建 worktree。

**本腿台件清单（README 规则 8：路径＋体积＋有没有被本件引用，⛔ 一条没删）**：整目录 `du -sh .scratch/wisp/probes/257/v1`＝**31 MB**，其中 30 MB 是一枚**未入库**的可执行件 `probe/wisp.exe`（§8 那发真进程读数的宿主，⛔ 别删，删了那格读数就不可重放）。入库的 35 枚全为文本：`pkg1-cmd-wisp.log`（最大，整包全文）／`mut-M1..M5-*.log`（五发红句原件）／`baseline-257-family.log`／`iso1-count1.log`／`iso2-ac14.log`／`iso3-count3.log`／`gate-{d22scan,pathlen,vet,gofumpt}.log`＋`gates.log`／`ac4-diff-names.txt`／`roster-{mine,theirs}.txt`＋`names-{mine,theirs}.txt`／`pkg1-status-approval-{before,after}.txt`＋`pkg1-goflags-*`＋`pkg1-head-*`／`iso1-status-approval-before.txt`／`firstrun.pristine`（突变还原用的 HEAD 快照）／`mutate257v1.py`（台件本体）／`probe/{real-run,real-run2,run}-{stdout,stderr}.txt`／`msg-skeleton.txt`＋`msg-s2.txt`。除 `probe/` 那六枚 stdout/stderr 直接顶 §3/§8 两格、五枚 `mut-*.log` 顶 §2/§4/§5/§7/§8 五格外，其余每枚都被 §6/§9/§10/§11 逐名引用；**没有一枚是"建了没人引用"的**。⛔ `probe/data-dev/**`（真进程那发的 `config.toml`＋logs＋secrets 目录形状）**刻意未入库**——那一格里有 DPAPI store 的目录形状，本票凭据只需要它的文本读数。
