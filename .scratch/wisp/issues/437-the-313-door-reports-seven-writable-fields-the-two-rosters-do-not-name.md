# 437 — `Test313WritableUnionEqualsTheTwoRosters` 红：门报告 26 个可写字段，两本名册只有 20 个——多出 6 个 ball/panel 面 + `role_chat_thinking`

- Status: **ready-for-agent**（写面＝`internal/panel/**` 的 `WritableFields()`／`AdditionalWritableFields()` 两张表，或 `cmd/wisp/panel_config_313_test.go` 的断言——⚠ **先只读取证，判清"哪边漏了"，再定改法**）
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

- [ ] **AC#0 只读取证（⛔ 本格不许改任何文件）**：
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

- [ ] **AC#1 按 AC#0 的判定修（只修被判定那一处）**：
  - **甲**：把确实可写的字段加进对应的表（⛔ 全部塞进一张表——按语义分表）。
  - **乙**：改门的报告逻辑，让它只报真的能写的。
  - **可能混合**：有的该加表、有的该从门里去掉。**若混合，逐个字段说明**。
  - **必须配突变控制真跑**：把修好的那处改回去 ⇒ 必须回到红。

- [ ] **AC#2 越界检查**：`git diff --name-only` 逐行贴。出现 `docs/PLAN.md`／`docs/specs/**`／`internal/observe/thresholds.go`／golden／`tools/d22scan/allowlist.txt`／`frontend/**`／`design/**` ⇒ **直接退回**。

- [ ] **AC#3 门禁**：`go build ./cmd/wisp/` · `go vet ./cmd/wisp/` · `gofumpt -l` 零输出 · **`scripts/d22scan.sh`**（⚠ ⛔ 用 `go run ./tools/d22scan`）。四条读数全贴。
  ⚠ `build/mut_rvw.go` 是未跟踪残留件，会让 `go build ./...` 红。**⛔ 不许删**。

- [ ] **AC#4 与票 313／68／216 的关系写清**：本票只处理这一格，⛔ 不翻那些票任何框。

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
