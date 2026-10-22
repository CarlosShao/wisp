# 223-v1 终裁腿：`AC#2` ⛔翻（射程裁定不成立）／`AC#7` 能翻

- 身份＝非实现者终裁腿。起手锚 `e94a5b17`；被裁对象 `223-r3` 四笔 `394cc212`／`d989a217`／`4824c267`／`b408b262`；件 `docs/evidence/s1/223-r3-impl.md` 现读 **107 行**。
- 我全部读数在**仓外**（`D:\tmp\wisp\20261022\223v1\`）：导出树 `exp/`（`git archive 5c284b3a`，`.go` 与 `7825d31e` 逐字节相同）、基线树 `base/`（`e94a5b17`）、233 基线树 `base233/`（`a8f3c020`）；仓内只落本件。
- ★**一句话**：`AC#7` 的出口**在生产里、有牙、why 逐字**⇒能翻；`AC#2` 的 config 层三档读数**全部成立且五发突变都有牙**，但**编排者的射程裁定不成立**（它把 09-29 的"233 落地之后再谈"读成了"射程＝config 层"）⇒**本轮 ⛔ 翻**，缺的是 reload 的生产消费者与 hot 的"行为"（见 §4）。

## §1 票面前提新鲜度（对 HEAD blob 复量）

| 前提 | 今天（HEAD `7825d31e`，文件与 `5c284b3a` 同） | 判 |
|---|---|---|
| `cmd/wisp/config_reload.go:119` | 逐字 `rt.mgr.OnRestartPending = rt.reportRestartPending`（`func startConfigReload` 在 `:109`） | ✅ |
| `:325-330` 打到 `rt.stdout` | `fmt.Fprintf(rt.stdout,` 在 `:325`，格式串 `:326-329`，`sections, strings.Join(restartTierKeys, " / ")` 在 `:330` | ✅ |
| `:327-328` 那句「为什么」 | 逐字 `原因：这几枚键在进程启动时被读一次就交给平台层（开机自启注册、单实例锁、界面语言），本次运行没有重新注册它们的路径。涉及：%s。` | ✅ |
| `:319-324` 两行审计 | `state=restart-pending … tier=restart`＋`RESTART-PENDING detail=%q keys=%q` | ✅ |
| `:342-344` `restartTierKeys` | `app.language / app.autostart / app.single_instance` | ✅ |
| `run.go:861` 装配根 | `rt.startConfigReload()` | ✅ |
| `panel_config_store.go:92/:98` | `cfg := s.mgr.Config()`／`ChatModel: cfg.LLM.Roles.Chat.Model` | ✅ |
| `tiers.go:48/:47` | `app.language:"restart"`／`app.theme:"hot"` | ✅ |
| 票 255 名册 `config_readers_255.go:110/:144` | `llm = consumed:`／`ball = no-reader:` | ✅ |
| `config_reload.go:183-184` | 唯一提 `rep.Reload` 的操作员面＝审计行（我 `git grep 'rep.Reload' cmd/wisp` 只 1 行 `:184`；**无 stdout 句**） | ✅ |

- ⚠ **票面自身的行号腐坏（不是腿的）**：票 223 现量写 `manager.go:113`，今天 `CheckAndReload` 在 `:142`；票 233 第 5 行写 `restartTierKeys` 在 `config_reload.go:299`，今天 `:342-344`；票 223 09-29 节写 `:284→:288`，今天 `:315-330`。票面自己写了"别信这里的行号"。
- ⚠ **"零产码"那把尺的射程**：`git diff --name-only 394cc212 HEAD -- '*.go'` 在**腿收工点 `b408b262`＝空**（我复跑）；在**今天 HEAD＝3 枚**（`internal/tools/task_164_ac4_cancel_tail_test.go`←`821302c9`、`internal/winsec/ancestor*_348_windows_test.go`←`195ce6d0`/`e9caddff`），**全属别的腿**。⇒ 能长期复跑的零产码尺是 `git diff e94a5b17 HEAD -- internal/config`（只有那枚用例）＋逐文件 blob 对 `hash-object`（§5 全 SAME）。台账 A1155 在 23:3x 仍写"空"，而 `195ce6d0` 23:18、`821302c9` 23:19 已在 HEAD ⇒ 那句在它自己的时刻就过期了；这是**尺的射程**问题，不是腿造假。

## §2 `AC#7` 核心判点：出口真在生产里，且真能答"为什么" ⇒ **能翻**

- ★**我自己现跑** `grep -rn 'OnRestartPending\s*=' --include=*.go cmd internal | grep -v _test.go` ⇒ **恰 1 行**：`cmd/wisp/config_reload.go:119`。`cmd/wisp` 单包同尺也＝1。对照：`OnReload\s*=` 非测试 `cmd internal` ＝ **1 行**（`cmd/balldebug/main.go:244`），`cmd/wisp` ＝ **0**（reload 档消费者仍 0，归 233）。
- ★**"用户在哪看到"是字节级核过的**：`:325` 的 `Fprintf` 第一参数就是 `rt.stdout`；`:327-328` 的 why 与腿报**逐字相同**（含 `原因：`／`交给平台层`／`涉及：%s。`）。
- ★**两枚既有用例是不是真钉住这一跳（我不看 PASS，我摘跳）**：
  - **VF**：把 `:119` 赋值注释掉 ⇒ `TestTicket223RestartTierSaysItWillNotApply` **FAIL 42.44s**，红行逐字 `config_reload_223_test.go:680: the audit trail never carried "config: HOT-RELOAD state=restart-pending sections=[app]" AFTER the plant within 40s`。⇒ 这枚用例**真的钉住生产赋值点**，不是只钉句子。
  - **VG**：把 `:327-328` 的 why 整段删掉（并去掉对应 `%s` 实参）⇒ 同用例 **FAIL 2.92s**，红行逐字四处 `:718: the operator's restart sentence omits "app.autostart"`／`"开机自启"`／`"原因："`／`"交给平台层"`。⇒ **why 在 stdout 上、且被钉住**（不是只靠审计行）。
  - **VH**：把 `tiers.go:48` 的 `app.language` 改 `hot` ⇒ `TestTicket255RestartTierKeysAreBackedByATest` **FAIL 0.21s**，红行逐字 `restart_tier_keys_255r2_test.go:61: restartTierKeys names "app.language" but config.TierRegistry says tier="hot" ok=true`；三个行为子测试仍 PASS。⇒ 255 钉的是 **why 的数据（名册＋分档）**，不是句子；句子由 223 那枚钉。两半合起来，AC#7 的"为什么"有牙。
- ★**AC#7 的"不许与 AC#2 第一档合成一枚用例"**：腿加了 `OneEditGivesTwoDifferentObservables`（同一笔改 `theme`＋`language`）作**结构证据**，但它同时有**各自独立**的 hot／restart／reload 用例；我的 VA／VB 证明两半可分开打红（§5）。⇒ 字面有一处张力（确实存在一枚"合读"的用例），但**实质不是用一枚绿冒充两格**；这一条**不挡翻勾**，记在这里。

## §3 `AC#2` 三档逐档：我的独立核（读数全部成立，但"行为"不是"读面"）

| 档 | 我的读数（自己跑） | 真值判 |
|---|---|---|
| **hot** | 生产读面 `configStore.ReadSettings`（`newConfigStore`＋真 `Manager`）在**一次轮询后同一进程**答 `m1→m2`：`zz_probe_223v1_test.go:87: production reader: "m1" -> "m2" after one poll, same process`；`Manager.Config()` 是**深拷贝后的新读**（跨轮询持有的旧指针仍 `dark/zh-CN`，新 `Config()` `light/zh-CN`） | 读面**真变了**；但 255 名册自陈 `llm` 的"已立即生效"＝"live reader answers differently, **not the in-flight model swapped**"（`config_readers_255.go:101-109`），且 `run.go:448`（`:436 cfg := mgr.Config()` 之后）的 provider chain 仍用 boot snapshot ⇒ **不是运行中的"行为变了"** |
| **restart** | 旧值 `zh-CN` 不动；`Report.Restart=[app]`；`OnRestartPending` 在**一次轮询内恰一次**；why＝`TierOf("app.language")=="restart"`（同段 `app.theme=="hot"` 对照） | 生产出口＋why 已由 §2 钉住；⚠ 腿的"恰好一次"只是**一轮内**：我同一 Manager 在"同字节重写＋第三次轮询"后读到 **2 次**通知（`notices after three polls = 2`），因为 `cur` 保留旧 `language` 而 `fresh` 仍是新值 ⇒ **每次真重读都会再告知一次**。这不是错，但"恰好一次"必须带"每轮"限定 |
| **reload** | 值应用＋`OnReload` 恰一次（`[voice]`）；**生产消费者 0**；纯 reload 改动**stdout 一句都没有**（`reportReload` 只对 `rep.Hot`／`rep.Locked` 写 stdout，`rep.Reload` 只进 `:183-184` 审计行） | 腿**没有偷偷声称"生产也生效"**：§1 写"只钉 config 层那一半"，§10 归 233，`tiers.go` 头部也写明 consumer 不在本腿 ⇒ **无越界声称**；但这一档的产品面今天**没有任何读数** |

- ★**"两格不合并"那枚证据我核了**：VA（注释 `:358` theme apply）⇒ `OneEdit` 红 `:218`、`EveryAppKey/app.theme` 红 `:278`，restart 三个子测试仍 PASS；VB（注释 `:202` OnRestartPending）⇒ `OneEdit` 红 `:225`、restart 三子测试红 `:291`，`app.theme` 仍 PASS。⇒ A／B **真可分开**。
- ⚠ 但 **A 半只是 config 层**：`app.theme` 在 255 名册是 `no-reader:`（`config_readers_255.go:205`），所以 `Config().App.Theme=="light"` 不是生产行为变化，只是内存活值变化。腿的 §3 措辞是"活值变了"，没冒充生产；但这也说明**这枚结构证据撑不起"生产行为变了"**。

## §4 射程裁定：⛔ **不成立**；`AC#2` ⇒ **⛔ 翻**

- ★**理由①（编排者的 ①）是反向的**：09-29 原句是"AC#2 仍不勾（它字面要'三档各有读数'，reload 那档今天确实没有），洞归口到票 233；⛔ 不许拿'已归口 233'去翻 AC#2 的勾——**233 落地之后再谈**"。这是**把 AC#2 的翻勾条件钉在 233（生产消费者）落地之后**；它没有说"射程＝config 层、⛔ 含生产消费者"。拿它当"射程＝config 层"的证据，恰好把它的条件句读反了。
- ★**理由②只证明"config 层有读数"，没证明这些读数满足 AC 的字面**：AC#2 第一档写的是"**行为变了**"，不是"值变了"；第二档写的是"**明确告知未生效**"——后者是用户面的话，腿的 config 用例只钉了 hook 被调用，用户面那句是 AC#7 的。若 AC#2 真是纯 config 层，第二档的"告知"就没有主体。
- ★**理由③（没拿 `[ball]` 冒充）为真**，但它是**诚实性**，不是**射程**、也不是**充分性**。三点合起来**推不出**"AC#2 按 config 层算成立"。
- ★**腿自己也拒绝闭合**：§0"reload 档生产消费者仍＝0……本腿具名交回、⛔ 冒充已闭合"；§1 第三行"⛔ 没有消费者"；§10.1 归 233。⇒ 今天翻 `AC#2` 会与实现腿自己的边界相撞。
- ★**缺什么（按 AC#2 字面）**：
  1. **reload 档的生产读数**：`cmd/wisp` 里给 `Manager.OnReload` 一枚真赋值点（今天 0；唯一非测试在旁支 `cmd/balldebug`），并让操作员**看到**这一档的结果（或按票 233 AC#3 的"不接"支给一句具名实话）。今天纯 reload 改动只在审计行出现。
  2. **hot 档的"行为"**：今天生产面只有 `configStore.ReadSettings` 这一个**读面**答新值（我已跑通），运行中的模型链仍是 boot snapshot（255 名册自己的 caveat）。要按字面翻，要么**显式改写票面**把"行为"收窄成"活读面答新值"（这是票面文字动作，⛔ 由终裁腿替做），要么接一条**动作面**的生产变化。
  3. **09-29 那句"233 落地之后再谈"必须先被显式撤销**（若编排者坚持 config 层射程）；用"它本身就是判决"来支撑，与它的条件句相反。

## §5 我自造的反向判据（逐字读数；仓外导出树，工作树零突变）

- 尺＝`go test -count=1 -run 'TestTicket223R3' -v ./internal/config/`；pristine `manager.go` `sha256=4bb091d99f3c62925558be03933217c9fbb10d74696a8c57d16c251a05b2e8d0`、`tiers.go` `sha256=1f168efe1d3e7a562f973cd48f0709d6bac18ee47b30dc9aa4bacae5260c4de0`（**SHA-256、64 位**；与 `git show 394cc212:<p> | sha256sum` 逐字相同）；每发改后 `rc=1`、还原后 `rc=0`、还原哈希逐字回到上面两串。

| 我造的突变 | 锚 | 变异红（逐字行） | 变异后 sha256（我的 sed） |
|---|---|---|---|
| **M4** `app.language` 登记改 `hot` | `tiers.go:48` | `:156 TierOf(app.language) = "hot"/true, want restart`／`:262 … want "restart"`；**其余 hot／reload／OneEdit 全 PASS** | `e7aede2f0397ee6c8ae23213aac4e21cd2074fce7d227cfdf66d5020308277c7` |
| **VA** 摘 theme 的 hot apply | `manager.go:358` | `:218 app.theme = "dark", want light`／`:278 app.theme is registered hot but the live value did not move: dark`；restart 三子测试 PASS | `8342273558db2ba3fbb64a6ffaaaceca82f6ffbc4c427a9ed2297ddc4cd53143` |
| **VB** 摘 restart 派发 | `manager.go:202` | `:151 OnRestartPending = []`／`:225 … want one [app] notice`／`:291 ×3 … told 0 times`；`app.theme` 子测试 PASS | `addf05c9057dd8214d8a26c655b7695265728182879fc0a33ecfd1e669a78d09` |
| **VC** 摘 hot 段 apply（`[llm]`） | `manager.go:302` | `:116 after one poll the live reader still answers "m1", want m2` | `3c84ff8d4406e7a71b342e9eeb3858778dac3922b39628d0b420084e2403b9ff` |
| **VD** 摘 reload 派发 | `manager.go:199` | `:185 OnReload = [], want exactly one [voice] event` | `d807ed3b5dc55056dd0ff4a74d1a1e76b192e45bc27d564198242b13b5be32ff` |

- ★**M4 是不是真的专打"为什么"**：是。它是**唯一**一个只红 `:156`／`:262`、而 **HotTier／ReloadTier／OneEdit 纹丝不动**的突变；行为面（值／通知／事件）一行没红 ⇒ 它打的是"档位表＝why 的真值"，不是行为。
- ★**别混 hash 算法**：上表全是 `sha256sum`（64 位）；`git hash-object` 是 SHA-1（40 位）——§5 的 blob 对账用后者，**两把尺不互比**。

## §6 门禁 rc 与形状核对（串行，变异全部结束后）

- `GOFLAGS= go build ./...` **rc=0**；`go vet ./internal/config/` **rc=0**；`go test -count=1 ./internal/config/` **rc=0**（`ok 2.527s`）；`sh scripts/d22scan.sh` **rc=0**（`clean`，ban #8 `internal/` **621** 枚、`cmd/` **148** 枚）；`sh scripts/check-path-length-budget.sh` **rc=0**（`VERDICT GREEN`：over-budget 61、rostered 61、not in roster 0、worst 229）。★**腿的 §7 没报第五道门（path-length）——它漏了一道，不是红**。
- ★**形状核对（两发同数）**：`-v` 两发均 **RUN=205／PASS 顶层=120／PASS 缩进子测试=85／FAIL=0／SKIP=0**；`205 = 120 + 85`；`grep '^    === RUN'`＝**0**（子测试 RUN 不缩进），故 205 含子测试 RUN。改前基线 `e94a5b17` 同尺 **196／115／81／0／0**，`196 = 115 + 81` ⇒ 差 **+9 RUN／+5 顶层／+4 子测试**，恰等于"5 枚顶层用例＋4 子测试"，**无多报少报**。
- ★**格式甲乙两尺**（`gofumpt v0.12.0`）：**甲（判据）** `exp/` 导出树 `gofumpt -l .` ＝ **20 行、全 `.scratch[\\/]`**，非 `.scratch` ＝ **0**，本件命中 **0**；**乙（工作树，只诊断）** ＝ **54 行**（腿报 49，工作树后来长了）、非 `.scratch` ＝ **3**，逐名 `cmd\wisp\models.go`／`internal\agent\approval\pending_read.go`／`internal\agent\tools.go`，`git diff --quiet HEAD` 三枚均 **rc=0**（不是腿的改动）；CR 逐枚 `389／131／225`＝腿读数。乙的 `gofumpt -l` 进程 rc=2（某枚 `.scratch` 未跟踪探针不可解析），甲 rc=0；本件两尺都 0 命中。
- ★**CRLF**：`tr -cd "\r" | wc -c` ⇒ `internal/config/tiers_223r3_test.go` **0**、`docs/evidence/s1/223-r3-impl.md` **0**；`gofmt -l` 本件两尺均空。

## §7 新增红 0（我的核）

- `internal/config` 整包 **FAIL=0** ⇒ 红名册空、**新增红 0**。
- `cmd/wisp` 我**只按名单跑两枚**：`TestTicket223RestartTierSaysItWillNotApply` **PASS 2.72s**、`TestTicket255RestartTierKeysAreBackedByATest` **PASS 0.11s（3 子测试）**；⛔ **不跑、不报整包四数**（`A1126`）。★我核了腿件：§8 明写"不跑、也不报"整包，全文没有 `cmd/wisp` 整包四数 ⇒ **它没偷偷用整包数**。
- 归因注：我跑门时工作树有**别的腿**一枚未提交 `.go`＝`internal/tools/task.go`（diff 逐字 `_ = cancel // MUT164v1b-V1`），在 `internal/config` 之外；我的 `internal/config` 读数在该子树干净的前提下取得。

## §8 写面纪律（四笔逐笔＋正控）

| 提交 | 内容 | 该笔 `.go` | 越界 |
|---|---|---|---|
| `394cc212` | `internal/config/tiers_223r3_test.go`（**+298**） | **1** | **0** |
| `d989a217` | `docs/evidence/s1/223-r3-impl.md` | **0** | **0** |
| `4824c267` | 同件（补号） | **0** | **0** |
| `b408b262` | 同件（`:287`→`:285`＋计数） | **0** | **0** |

- 尺＝`git show --name-only --format= <c> | grep -c '\.go$'`；四笔的路径**逐名都在写面**（`tiers_223r3_test.go`／`223-r3-impl.md`），无 `.scratch`／`cmd`／`scripts`／`.github`。
- **正控**（同尺打含 `.go` 的真提交）：`195ce6d0` ⇒ `.go`＝**2**；腿用的 `17cfe8cf` 我也复跑＝**2**（`internal/memory/dao_session.go`／`internal/session/ticket166r1_ac2_gaps_test.go`）。**负控**：`46e61076`＝**0**、`b4102f1e`＝**0**。⇒ 这把尺**看得见 `.go`**，上表 0/1 不是恒 0。
- **零产码逐枚**（尺＝`git rev-parse HEAD:<p>` 对 `git hash-object <p>`，**同为 SHA-1**）：`internal/config/*.go` **41 枚全 SAME**（17 枚产码＋24 枚测试），`diff=0`；`git status --porcelain -- internal/config`＝**空**。

## §9 票 233 那处"一真一半腐"的数

- **真**：`TierHot/TierReload/TierRestart` 非测试命中各 **2**（`schema.go:34/36`、`37/39`、`40/43`＝注释＋声明）⇒ **赋值点 0**。
- **半腐**：`TierRegistry` 的 `"reload"` 行 **＝20**（`tiers.go:53-72`，全 `voice.*`）；`a8f3c020` 时 `tiers.go` **不存在** ⇒ "自票 255 起有 20 枚"对**登记表**成立。
- ★**我的加一层更正**：233 第 4 行说"没有一条键路能填进 `rep.Reload`"，**在它自己的起手锚 `a8f3c020` 就已假**——`planVoice` 在 `:361`、`rep.Reload = append(…"voice")` 在 `:412`，我在那棵树上跑 `TestManagerReloadTierEmitsEvent` **PASS**。⇒ 真正腐的是**"生产者键"与"登记行"混尺**：生产者键一直在，255 补的是**逐键登记**；"缺的只剩消费者"对**登记/消费**框架成立，对"生产者"不成立。
- "第 11 处"是**台账自己的序数**（`A1155③`，`docs/reports/pending-and-issues.md:22313` 逐字"这是'票面腐坏'的第 11 处"）；台账里"腐坏"共 84 处提及、无编号总表 ⇒ **序数本身不可由树独立复核**，但它所依据的 **20** 我核了。

## §10 判语（交回）

- **`AC#2`：⛔ 翻**。config 层三档读数**全部成立、五发突变都有牙**（§3／§5），但**射程裁定不成立**（§4）；缺 **reload 生产消费者＋操作员读数**、**hot 的"行为"（今天只是读面）**，以及**对 09-29"233 落地之后再谈"的显式撤销**（若坚持 config 层射程）。⇒ 这不是"做得不好"，是**格还差两半**；233 落地＋票面收窄之后再谈。
- **`AC#7`：能翻**。生产赋值点恰 1（`:119`）、用户面在 `rt.stdout`（`:325-330`）、why 逐字（`:327-328`）；VF 摘赋值 **42.44s 红**、VG 删 why **2.92s 红四处**、VH 改档 **:61 红** ⇒ 两枚既有用例**真的钉住这一跳**，不是只看 PASS。唯一记项：AC#7 字面"不许与 AC#2 第一档合成一枚用例"与腿那枚 `OneEdit…` 有**字面张力**，但独立用例齐备、VA／VB 证明两半可分 ⇒ 不挡翻。
