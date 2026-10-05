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

判语：未判。凭据：未判。

## 8. ★裁定二：正向断言（"提示真到达 stderr/stdout"）是不是恒真

判语：未判。凭据：未判。

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

## 10. 门禁四门（本腿自跑，不引用腿的读数）

判语：未判。凭据：未判。

## 11. 突变名册（每发起止时刻到秒，编排者拿它作差）

未判。

## 12. 七项对抗检查（`SPEC-10 §8`）

未判。

## 13. 我攻不动的地方／判不动的格子（具名归口，⛔ 一枚不写成"应该没问题"）

未判。

## 14. 交件判语与 commit 链

未判。
