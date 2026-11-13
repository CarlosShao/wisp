# 437 — `Test313WritableUnionEqualsTheTwoRosters` 红：门报告 26 个可写字段，两本名册只有 20 个——多出 6 个 ball/panel 面 + `role_chat_thinking`

- Status: **done**（写面＝`internal/panel/**` 的 `WritableFields()`／`AdditionalWritableFields()` 两张表，或 `cmd/wisp/panel_config_313_test.go` 的断言——⚠ **先只读取证，判清"哪边漏了"，再定改法**）
- 来源：round 89 整包（`D:\tmp\wisp\20261112\round89\full.log`，`PASS=583 / FAIL=11`）＋台账 `:32512` 存量真红
- 关联：票 313（A881：union == table1 + table2）· 票 68／77／216（ball/panel 面）· 台账 `A1806`

## 这是什么（人话）

面板配置有一道"门"（`panel` 的命令分发器），它对外**报告自己能写哪些字段**。
这个报告必须是**两张表的并集**：
- 表1 `WritableFields()` —— 基础可写字段
- 表2 `AdditionalWritableFields()` —— 追加可写字段

有测试盯着这件事：**门报出来的，必须一字不差等于两表的并集**（A881 的规矩）。
现在：**门报了 26 个，两表并集只有 20 个** ⇒ 门**多报了 6 个**。

## 现量（逐字，可复跑）

**红句**（`cmd/wisp/panel_config_313_test.go:252`）：

门报（**26 个**）：
```
audio_input_device ball_click_through ball_opacity_idle ball_size ball_style
model_context_window model_price_in model_price_out model_row_create model_row_delete
panel_font_size panel_width provider_api_key_ref provider_base_url provider_billing
provider_credential provider_row_create provider_row_create_custom provider_row_delete
role_chat_model role_chat_pair role_chat_thinking
voice_aec_enabled voice_cloud_asr_chain voice_cloud_tts_chain
voice_realtime_model voice_realtime_provider
```

两表并集（**20 个**）：
```
audio_input_device model_context_window model_price_in model_price_out
model_row_create model_row_delete provider_api_key_ref provider_base_url
provider_billing provider_credential provider_row_create provider_row_create_custom
provider_row_delete role_chat_model role_chat_pair
voice_aec_enabled voice_cloud_asr_chain voice_cloud_tts_chain
voice_realtime_model voice_realtime_provider
```

**★多出的 6 个（逐名，我核对过两边）**：
1. `ball_click_through`
2. `ball_opacity_idle`
3. `ball_size`
4. `ball_style`
5. `panel_font_size`
6. `panel_width`
7. `role_chat_thinking`

⚠ **实际是 7 个**（我数了两遍：26 − 20 = 6，但逐名比对上列了 7 个——**腿必须自己重数**，这是我在票面里就暴露的不确定，⛔ 照抄我的数）。

**耗时 0.00s** ⇒ **稳定，与负载无关。**

## 机制链（逐字，⛔ 凭印象）

- 测试：`cmd/wisp/panel_config_313_test.go:246-253`
  ```go
  reported := strings.Split(rest, "/")
  want := append(append([]string{}, panel.WritableFields()...), panel.AdditionalWritableFields()...)
  sort.Strings(want)
  if !reflect.DeepEqual(reported, want) {
      t.Errorf("the door reports the union %v; the two rosters say %v (A881: union == table1 + table2)", reported, want)
  }
  ```
- ⇒ 它**跑真实的门**（`disp.Handle`），拿到报告字符串，与两表并集比。
- **⚠ 仪器**：这是**运行时行为**检查 ⇒ `go test -overlay`（拷贝放**仓外**）**有效**。（对照台账 `A1804`：读磁盘源码的静态检查才无效。）

## ⚠ 必须自己判清的核心岔路（⛔ 不许我替你定）

**"门多报了" vs "表漏了"——两种答案，修法不同：**

- **甲：表漏了（门是对的）** ⇒ 那 7 个字段**确实可写**，只是没登记进 `WritableFields()`／`AdditionalWritableFields()`。⇒ **加表**。
  **判据**：逐个去门的分发逻辑里查，这 7 个字段**是不是真的有处理分支**。有 ⇒ 甲。
- **乙：门多报了（表是对的）** ⇒ 那 7 个字段**不该可写**，门却放行了。⇒ **改门的报告逻辑**（或改分发，让它们真的不可写）。
  **判据**：若门只是"报告"了但**分发时并没有真的处理**它们，那是**报告撒谎**——⛔ 加表盖章。

**★这是本票的关键：⛔ 用"把表补齐"来消红。**
若那 7 个字段确实可写（有真实处理分支）⇒ 补表（甲）；若它们只是被**列出来**但**写不进去** ⇒ 门在撒谎，改门（乙）。

**⚠ 特别注意 `role_chat_thinking`**：它与其他 6 个（ball/panel 面）**不同族**，可能是**另一次改动带进来的**。腿要单独查它。

## 要做的（按序）

- [x] **AC#0 只读取证（⛔ 本格不许改任何文件）**：
  1. 用 carrier 复现，贴逐字红句：
     ```
     D:\work\soft\Git\bin\bash.exe -c "scripts/wisp-cli-staged-tests.sh '-test.run' '^Test313WritableUnionEqualsTheTwoRosters$'"
     ```
     ⚠ **必须用 `D:\work\soft\Git\bin\bash.exe`**（本机 `bash` 是 WSL 转发器，一枚测试都不跑）。
  2. **自己重数**两边集合的差（⛔ 照抄我票面的 6／7——我在票面里已经暴露了这个不确定）。**逐名列出差集。**
  3. 对**每一个**差集字段，去门的分发逻辑里查：**它有真实处理分支吗？** 逐字段给出 `文件:行`。
  4. `git log --oneline -10 -- 'internal/panel/**'`，找最近谁动过这两张表或门的报告。
  5. **单独查 `role_chat_thinking`**：它是不是另一次改动带进来的（`git log -S 'role_chat_thinking'`）。
  6. 判甲／乙（**可能逐个字段不同！**），给一句话结论＋逐字段依据。

- [x] **AC#1 按 AC#0 的判定修（只修被判定那一处）**：
  - **甲**：把确实可写的字段加进对应的表（⛔ 全部塞进一张表——按语义分表）。
  - **乙**：改门的报告逻辑，让它只报真的能写的。
  - **可能混合**：有的该加表、有的该从门里去掉。**若混合，逐个字段说明**。
  - **必须配突变控制真跑**：把修好的那处改回去 ⇒ 必须回到红。

- [x] **AC#2 越界检查**：`git diff --name-only` 逐行贴。出现 `docs/PLAN.md`／`docs/specs/**`／`internal/observe/thresholds.go`／golden／`tools/d22scan/allowlist.txt`／`frontend/**`／`design/**` ⇒ **直接退回**。

- [x] **AC#3 门禁**：`go build ./cmd/wisp/` · `go vet ./cmd/wisp/` · `gofumpt -l` 零输出 · **`scripts/d22scan.sh`**（⚠ ⛔ 用 `go run ./tools/d22scan`）。四条读数全贴。
  ⚠ `build/mut_rvw.go` 是未跟踪残留件，会让 `go build ./...` 红。**⛔ 不许删**。

- [x] **AC#4 与票 313／68／216 的关系写清**：本票只处理这一格，⛔ 不翻那些票任何框。

## 本票**不**解决

- 不碰 `TestTicket198*`（票 434 在飞）／`TestLockedSuffix*`（票 435）／`TestAC247*`（票 436）。
- 不改 `docs/PLAN.md` 与 `docs/specs/**` 一字。

## 派单约束

- 先只读（AC#0 零写），判定后再定写面。
- ⛔ `git add -A`／`add .`／`--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`；`add`／`commit` 两处显式 pathspec。
- 只 commit，⛔ push。
- ⛔ 翻任何 AC 框。
- 临时件只建不删，落 `D:\tmp\wisp\<YYYYMMDD>\437-<腿码>\`。
- **⚠ 本票与真窗族测试（`Test258*`／`TestAC4Focus*`）⛔ 并行**。

## Progress log（append-only，newest last）

- [2026-11-13T01:35Z] agent=437-313 read-only-first leg did=AC#0 只读取证完成，**判定为「断言侧陈旧」，零产品改动**。① 复现（carrier，`D:\work\soft\Git\bin\bash.exe`）逐字红句已存；② **自己重数：门报 27、两表 20、差 7**（票面写的 26／6 有误，逐名清单本身是对的；门列表无重名，roster-only=0）；③ 差集七名逐一查分发器：`ball_click_through`(panel_config_store.go:354)／`ball_opacity_idle`(:340)／`ball_size`(:326)／`ball_style`(:321)／`panel_font_size`(:377)／`panel_width`(:365)／`role_chat_thinking`(:262)——**七枚全有真实 case 且各调真 setter**，无一枚落 default 拒支；全量核对 27 报名 ↔ 27 个 `case panel.Field…` ↔ 0 缺支；④ 机制：门报的是**五册并集** `acceptedWritableFields`（config_write_313.go:176-183），断言只读前两册（7+13=20 vs 7+13+1+2+4=27），差集恰＝313 之后落地的三册（187=1＋65=2＋wave=4）；⑤ `role_chat_thinking` 单独查：生于 `007de34f0`(2026-11-07，票187)，写腿支在 `7817ad09a`(同日至次)，**比 ball/panel 六枚早三天、另一票带入**，非游离件（有 fieldTier187 与三枚钉）；⑥ 判：**七枚逐枚皆「断言陈旧」，无一是乙（门撒谎）**，甲也不是——补表会重名。AC#1 结论＝**本波写面内无正确改法**：M2（票面诱惑的甲：七名塞 fields313）⇒ 门报涨到 34 且七名**各重名两次**、红；M6（把三册并进 AdditionalWritableFields）⇒ 同样重名、红；M3（把报告缩回两册）⇒ 断言绿但 internal/panel 四枚钉红（appearance_wave_test.go:67／ball_size_422_test.go:54／ball_style_65_test.go:55／config_write_187_test.go:134）⇒ 正是票面禁止的"报告撒谎"。**正确修法在断言侧**（`cmd/wisp/panel_config_313_test.go:249` 三行并上 187/65/wave），实测 M1/M4 绿、M5 仍有牙（截一册即红并逐名点名），**但该包本波另一腿独占 ⇒ 按派单停下报告，未改**。门禁四读数：build rc0／vet rc0／gofumpt 空（⚠ 本机 gofumpt 不在 PATH，真身 D:\work\base\gopath\bin\gofumpt.exe）／d22scan clean（480 件，PASS=45 FAIL=0 SKIP=0）。证据件 `docs/evidence/s1/437-313-door-union.md`。`git diff --name-only` 无任何产品文件。⛔ 未 push／未翻 AC 框／未删 build/mut_rvw.go。next=断言侧三行修复请编排者另派一波（该包写面互斥）；本腿无残留待办。
---

## ★取证腿结论（2026-11-13 09:3x，腿 `ticket-437-313-door-union`，commit `018b3112d`）

**判定：门是对的，过时的是断言侧。** 那 7 个名字**全部有真实写分支**，一个都不是"门在撒谎"。

- **★我票面的数错了**：实际 **门 27 / 两册 20 / 差 7**（我写 26／6）。腿**机械解析红句**得出、非目测；并核对门列表无重复。**我的名字清单是对的，数错了。**（第 10 次同根错）
- **★差异的来源（算术闭合）**：门的报告**⛔ 两册并集，是五个名册的并集**（`acceptedWritableFields`，`internal/panel/config_write_313.go:176-183`），而断言的 `want` 只拼了**两个**。`7 + 13 = 20` vs `7 + 13 + 1 + 2 + 4 = 27` ⇒ 那 7 个正是**票 313 之后落地**的名册：票 187（1 个）＋票 65（2 个）＋某波（4 个）。
- **★逐字段真实分支（全有，`panel_config_store.go`）**：`ball_click_through` `:354`→`SetBallClickThrough:364`／`ball_opacity_idle` `:340`→`:353`／`ball_size` `:326`→`:339`／`ball_style` `:321`→`:325`／`panel_font_size` `:377`→`:386`／`panel_width` `:365`→`:376`／`role_chat_thinking` `:262`→`:267`。**全扫：27 个报告名 ↔ 27 个 case ↔ 0 个无分支。**
- **★不止"列出来"，是真的能写进去**：把测试自己 half-2 的验收循环扩到全部 27 个 ⇒ **PASS**（27 个都被接受且各自到达腿一次）。
- **★`role_chat_thinking` 单独查清了**：**不同族、不同票、早三天**。生于 `007de34f0`（2026-11-07，票 187），其腿分支在 `7817ad09a` 加——**那笔提交记录该字段此前是 fail-closed**。ball/panel 那 6 个是 2026-11-10/11 来的。**⛔ 野的**：自己的名册、自己的实测档（`fieldTier187`→`EffectiveNextTask`）、三个钉子。
- **★为什么 `internal/panel` 侧无正确修法（三条实测）**：**M2**（我票面诱人的甲：把 7 个加进 `fields313`）⇒ 门报 **34 且 7 个重复**，**红掉三个不相交钉子**；**M6**（把后落名册并进 `AdditionalWritableFields`）⇒ 同样重复、同样三钉红；**M3**（把报告缩到两册）⇒ **断言绿了，但四个 `internal/panel` 钉子红**，每个都在说**拒绝文案现在漏掉了一个已接受的名字**——**那正是我票面禁止的"报告撒谎"**。
- **⇒ 修法必须在断言侧：`cmd/wisp/panel_config_313_test.go`**。本波写面不含它 ⇒ **腿零改动、只交证据件，测试仍红**（我亲跑复核：`FAIL 0.00s`）。**这是诚实交付。**

## 下一步（编排者已派）

断言侧单独立票并独占 `cmd/wisp/**` 一波。

## ★断言侧写腿（2026-11-13 09:4x，腿 `437b-313assert`，commit 见下）

**判定沿用上一腿，不重做。** 本腿只做单点：把断言的 `want` 改成与门同一口径——**五册**，⛔ 两册。

- **五个名册的真实函数名（逐字，读的是定义处而非调用点）**：`WritableFields()`（`config_handlers.go:100`，7 枚）· `AdditionalWritableFields()`（`config_write_313.go:154`，13 枚）· `WritableFields187()`（`config_write_187.go:66`，1 枚）· `WritableFields65()`（`config_write_65.go:69`，2 枚）· `WritableFieldsWave()`（`config_write_wave.go:62`，4 枚）。`7+13+1+2+4=27`，与门报数闭合。
- **改动（唯一一处，`cmd/wisp/panel_config_313_test.go:249-252`）**：`want` 在三行里并上 187／65／wave。改前＝`want := append(append([]string{}, panel.WritableFields()...), panel.AdditionalWritableFields()...)`；改后＝该行原样保留，再追加 `want = append(want, panel.WritableFields187()...)`／`…65()`／`…Wave()`。**`sort.Strings(want)` 与 `!reflect.DeepEqual(reported, want)` 一字未动**——仍是**等集**比较，⛔ 未放宽成 Contains／子集。
- **`^Test313` 全族复跑（carrier `D:\work\soft\Git\bin\bash.exe`）**：改前 **PASS=2 / FAIL=1**（红的是 `Test313WritableUnionEqualsTheTwoRosters` 0.00s）；改后 **PASS=3 / FAIL=0**。全族恰三枚，另两枚一直绿。
- **half-2 现在验 27 个**：`want` 即 half-2 的循环源，故由 20 涨到 27。**实测读数**（仓外 overlay 加一行 `t.Logf`）：`PROBE437B half2 covers 27 union members` ＋ `--- PASS`，零条 "refused"/"reached the leg N time(s)" ⇒ **27 个全被门接受且各到达腿一次**。
- **突变控制（真跑，两份拷贝都在仓外 `D:\tmp\wisp\20261113\437b-313assert\`）**：**M1**＝`want` 退回两册 ⇒ **FAIL**（`panel_config_313_test.go:253`，逐名点名那 7 枚）；**恢复**（无 overlay）⇒ **PASS 3/3**；**M5**（额外加牙）＝五册 `want` 但把门的报告截掉 wave 册 ⇒ **FAIL**，恰好点名四个 wave 字段（`ball_click_through`/`ball_opacity_idle`/`panel_font_size`/`panel_width`）。⇒ 改动是绿的原因（M1 红），且等集比较两个方向都有牙（M5 红）。
- **门禁四读数**：`go build ./cmd/wisp/` rc0 空 · `go vet ./cmd/wisp/` rc0 空 · `gofumpt -l`（真身 `D:\work\base\gopath\bin\gofumpt.exe`）**空** · `scripts/d22scan.sh` **clean，rc0**。`build/mut_rvw.go` 未删。
- **边界**：`internal/panel/**` **零改动**（`git diff --name-only -- internal/panel/` 空）。本腿自己的代码改动只有 `cmd/wisp/panel_config_313_test.go` 一处三行；`git diff --name-only` 里另有 11 份 `.scratch/wisp/issues/43x` 删除、`probes/291`、`probes/orch/.../346-a2-dispatch.txt`、`frontend/fixtures/composer-states.html`——**都是别的腿的存量脏**，按 mtime 排除（本会话 09:43 起，那几件最后写入是 11-09／11-11／11-12）。证据件 `docs/evidence/s1/437b-313-assertion-five-rosters.md`。⛔ 未 push／⛔ 未翻任何 AC 框／⛔ 未改 half-2／⛔ 未改别的测试。
- **留一条小账**：失败句里 "the two rosters say …" 的字样现已过时（两侧都是五册）。本腿被派单限定"只改期望集的构造"，故**未动文案**——若要改，另开一格，属措辞不属判据。
- [2026-11-13 09:5x +08] **编排者：两波收件＋三把尺全过＋AC#0–AC#4 五格翻勾＋`-done`**。**第一波**取证腿 `018b3112d`（只交证据件 15063B，零改动，判"门对、断言过时"）；**第二波**写腿 `3fe4f7365`（3 files／+199；证据件 `437b-313-assertion-five-rosters.md` 9193B）。**三把尺**：①两笔 commit 均在；②两件均在盘；③**我自己复跑 `^Test313` 全族 ⇒ PASS=3 / FAIL=0**（含 `Test313WritableUnionEqualsTheTwoRosters` PASS 0.00s）。**越界**：`git diff --name-only -- internal/panel/` **空** ⇒ **面板零改动**，只改 `cmd/wisp/panel_config_313_test.go` 的 `want` 构造。
- **★五个名册的真实函数名（腿从定义处读，⛔ 从门调用点猜）**：`WritableFields()`（`config_handlers.go:100`，7 个）＋`AdditionalWritableFields()`（`config_write_313.go:154`，13 个）＋`WritableFields187()`（`config_write_187.go:66`，1 个）＋`WritableFields65()`（`config_write_65.go:69`，2 个）＋`WritableFieldsWave()`（`config_write_wave.go:62`，4 个）。**`7+13+1+2+4 = 27`** ⇒ **与门报的数闭合**。
- **★改动实质（只加三行 append）**：`want = append(want, panel.WritableFields187()...)` / `…Fields65()...` / `…FieldsWave()...`。**`sort.Strings` 与 `!reflect.DeepEqual` 一字未动** ⇒ **仍是等集比较，⛔ 放宽成包含**。
- **★half-2 从 20 扩到 27（实测，⛔ 推断）**：腿用仓外 overlay 加一行 `t.Logf` 量出 `PROBE437B half2 covers 27 union members`，且 **PASS 无任何 "refused" / "reached the leg N time(s)" 行** ⇒ **27 个全部被接受、各自到达腿一次**。
- **★突变控制两发**：**M1（必做）** `want` 改回两册 ⇒ **FAIL `:253`**，逐名点出全部 7 个多余项；**M5（额外牙口）** 五个名册 `want` ＋ 门报告截断（去掉 wave 拼接）⇒ **FAIL**，逐名点出**正好那 4 个 wave 字段** ⇒ **证明比较两向都有牙，⛔ 子集捷径**。
- **★腿留的一条未办（我认可它的克制）**：失败句仍写 `the two rosters say …`，**措辞已过时**；但它⛔ 比较的一部分，且我派单把它限定在"期望集构造"，**它没顺手改**。**这条留给下一个碰这个文件的人**。
- **撤销口令：「437 回两册」**——删掉那三行 append 即复原，`internal/panel` 全程未动。

