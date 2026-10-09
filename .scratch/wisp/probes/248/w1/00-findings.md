# 248-w1 · AC#8 文案名册（只读审计腿）

- 取数时刻 `2026-10-09 11:38 +0800`（同一条 bash 里 `date '+%F %H:%M %z'` 现量）。
- 锚点：分支 `dev`，`git rev-parse --short HEAD` = **`c2b422a4`**。
- 射程口径（⛔ 不是"全仓"）：读数一律 `git show HEAD:<path>` / `git grep … HEAD -- <pathspec>`。
  判据对象按派单＝**用户会读到的字面**＝产码里的字符串字面量（`cmd/`、`internal/`，排除 `_test.go`）＋页面文案（`frontend/src`、`frontend/index.html`、`frontend/fixtures`，只读对象层文案）。
  `.scratch/wisp/**`（工单／台账／裁决／probes 拷贝）与 `docs/**` **只计数、只登记，不归罪**；probes 里那些 `mutation/`、`backup/`、`baseline-*` 拷贝是别票的变异夹具，不是发布物，**不进口径①**。
- 零改动：本件是本腿唯一写点。未跑任何 go／scripts／gh／网络。

---

## 0. 票面那一格的复核（派单给的行号我先验过）

`.scratch/wisp/issues/248-…-methods.md:56` **确实在 `:56`**，原文（一次取数）：

> `- [ ] **AC#8（归口"重建那一跳"）**：`[llm]`／endpoint 类字段在同一次运行内不重建（J9 现量），故 AC#4 只按"重启后生效"判；本格判据＝票面与回执文案**不许出现"保存即生效"**，且"要重启"必须到达**页面可见面**。⛔ 本票不许新增重建跳（票 223 地界）。`

⚠ 但同一枚票 **`:126`** 已经写着「**AC#8＝附条件成立**」，且 `docs/evidence/s1/248-settings-write-path-v1.md:146`／`:152`／`:254` 已把这格裁过一次并把"两套档位词表互不相认"那一问交回编排者。⇒ 这一格没勾的原因**不是"没人看过文案"**，是那一问待裁。见 §5 顶回第 1 条。

---

## 1. 假话名册（"保存即生效"这一族）

### 1a. 字样扫描（整词族，一次取数，计数＝**全仓含 probes/docs**，别当产码数）

| 字样 | HEAD 命中行数 | 处置 |
|---|---|---|
| `立即生效` | 465 | 按目录分：`.scratch/wisp` 138 枚文件／`cmd/wisp` 7／`docs/evidence` 6／`docs/reports` 3／`internal/panel` 2／`docs/specs` 1／`docs/PLAN.md` 1／`.scratch/ci-logs` 1 |
| `保存即生效` | 22 | 产码／页面**零枚**（尺：`git grep -F "保存即生效" HEAD -- cmd internal frontend tools docs/specs` ⇒ 零命中）；命中全在工单、裁决、台账里**引用禁句** |
| `不用重启` | 21 | 产码里只有 `cmd/wisp/panel_host_windows.go:544` 一枚**注释**（引号里的"不用重启就见效"，说的是热键重绑那条腿），登记不归罪 |
| `takes effect` | 17 | 全为注释（`internal/panel/composer_dispatch.go:152`、`internal/panel/config_handlers.go:380`、`internal/config/validate.go:106` 等） |
| `applies immediately` | 14 | 含 `internal/config/schema.go:34,37`（注释，第二套词表的措辞来源） |
| `即时生效` | 4 | 产码／页面零命中 |
| `hot reload` | 11 | 全为注释／`internal/config/doc.go:2` 一行的模块自述 |
| `无需重启` | 1 | 产码／页面零命中 |
| `no restart` | 1 | 产码／页面零命中 |
| `马上生效` | 0 | — |
| `live reload` | 0 | — |

### 1b. ① 档（说的是 `[llm]`／endpoint 类字段＝假话，要改）—— **2 枚**

1. `cmd/wisp/config_reload.go:191`
   `fmt.Fprintf(rt.stdout, "wisp run: 配置热加载：这些段已立即生效（D36 立即档）：%s\n", bracketed(split.claimable))`
   **`[llm]` 今天真会进 `split.claimable`**：尺＝`cmd/wisp/config_readers_255.go:103` 把 `llm` 行的判定前缀写成 `hotClaimConsumed`（`:79` 常量），而 `splitHotTier` 在 `:275` 只按 `strings.HasPrefix(verdict, hotClaimConsumed)` 放行 → `:280` `s.claimable = append(…)`。
   ⇒ 人手改 `config.toml` 的 `[llm]`（endpoint／模型就在这段）时，**运行中的 `wisp run` 会在 stdout 印"[llm] 已立即生效"**，而同一份名册 `:100-102` 自己写着这句 caveat：真正的读者只有面板那条"读设置"的腿（`cmd/wisp/panel_config_store.go:96`），跑任务的模型链是装配时 `run.go:435 [res := llm.NewResolver(cfg, st)]` 建的、本进程不换。caveat 只落 stderr 的 `HOT-RELOAD-READER` 审计行（`config_reload.go:188`），**印给人看的那一句里没有它**。这一枚是本腿的新料，喂进 §5 那一问。
2. `internal/panel/config_handlers.go:431`
   `return "这一项立即生效。"`（`tierSentence` 的 `case EffectiveNow:`，`:430`）
   作用对象＝**面板可写名册里的字段**，而那七枚全是 `[llm]`／endpoint 族（`:90-93` `allWritableFields = {provider_base_url, provider_api_key_ref, model_context_window, model_price_in, model_price_out, role_chat_model, provider_credential}`）⇒ 这一支一旦可达，说的就是 endpoint 字段。
   今天**不可达**：`cmd/wisp/panel_config_store.go:228`（`res.Tier = panel.EffectiveRestart`，前面零按键判断）与 `:275` 无条件填重启档；尺 `git grep -n "Effective(Now|NextTask)" HEAD -- cmd internal ':(exclude)*_test.go'` ⇒ 除定义 `:199-200` 与 `case` `:430/:432` 外**无写者**。
   ⇒ 判①的凭据是"文案形状"（AC#8 词面判据直接命中），不是"今天已经骗到人"。既有牙：`cmd/wisp/panel_config_248_test.go`、`internal/panel/config_route_248_test.go` 里 grep 禁句的钉在。

### 1c. ② 档（实话，另给"凭什么成立"的调用点尺）—— **4 枚**

3. `cmd/wisp/config_reload.go:326`（`reportRestartPending`，函数体 `:315-331`）
   `"wisp run: 这些段的改动本次运行不会生效（D36 重启档，需要重启进程）：%v。原因：这几枚键在进程启动时被读一次就交给平台层（开机自启注册、单实例锁、界面语言）…"`
   作用对象＝`restartTierKeys`，与 `internal/config/tiers.go:48-50`（`app.language`／`app.autostart`／`app.single_instance` = `restart`）同源。尺：`config_readers_255.go:236` `if tier, ok := config.TierOf(name); ok { // TierOf's first production caller`。⇒ 实话，但**射程不含 `[llm]`**（`tiers.go:30 "llm": "hot"`，永远进不了这一支）。
4. `cmd/wisp/config_reload.go:194-197`（quiet 那一半）
   `"…这些段的值已换进内存，但本宿主没有会按新值做事的读者，本次运行不会因此改变行为（票 255 AC#1：这一半不许说成「已立即生效」…）"`
   尺＝`config_readers_255.go:79-83` 五种判定前缀（consumed / snapshot-only / no-reader / other-process / debug-host-only）＋名册每行必须带 `path/file.go:LINE [token]`，由 `TestTicket255RosterEvidenceLinesStillSayWhatTheyClaim` 重读那一行验 token。⇒ 形状成立，问题只在 `[llm]` 被归到了 consumed（见①-1）。
5. `frontend/src/components/config-screen.tsx:99-100`＋`:114-115`（`live ? "bg-green-tint text-green" : …` 那枚"生效徽标"）
   活的两项是 `panel.opacity`／`panel.animations`，落地方式是页面自己写 DOM：`:48` `const ALPHA_VAR = "--panel-alpha"`、`:54` `MOTION_ATTR = "data-motion"`（`:61-62` 真去 `setAttribute`）。不经 Go、不经 `config.toml`。⇒ 对这两枚是实话。
   ⚠ 附一句给落地腿的备注（不改口径）：`:83` 注释逐字「The config section, spelled the way config.toml spells it」，行名前缀 `app`／`panel`／`ball` 与配置段同名 ⇒ 用户可能把"绿徽标"读成"这段配置现在生效"。其余项已用 `blocked:` 说明改不到（`:98`／`:101-105`）。
6. `cmd/wisp/config_reload.go:240`
   `"wisp run: 注意 - 本次运行的路径判定仍按启动时建好的 C26 名单，新放宽的目录要重启进程才参与判定（原因：canonicalizer 只在装配时构造一次）。"`
   对象＝`[fs]` 放宽（`d.Section == "fs"`，`:235`）。⇒ 实话，且它是**今天唯一一句印在 stdout、带"要重启进程"字样、并给了原因的话**，但它和 `[llm]` 无关。

### 1d. ③ 档（说不清作用对象／须人裁）—— **3 枚**

7. `docs/evidence/s1/248-settings-write-path-v1.md:152`＋`:254`、票面 `:126`：**两套档位词表零行互译**（`internal/config/tiers.go:30` 记 `llm`＝hot，`cmd/wisp/panel_config_store.go:228` 硬填 restart；台账 `docs/reports/pending-and-issues.md:10738`／`:11103` 也记着"AC#8／AC#10 翻勾起跑判据＝`TierOf` 真被产码吃进去"）。本腿 §1①-1 是这一问的直接证据：那句 stdout 今天**真会**为 `[llm]` 印"已立即生效"。**要人裁的不是文案，是 `[llm]` 到底归哪一档。**
8. `internal/config/schema.go:34,37`：`// TierHot applies immediately: the new value is visible through …`／`// TierReload applies immediately AND emits a reload event …` —— 纯注释，但它是第二套词表的措辞来源（`config.Tier` 与 `panel.EffectiveTier` 并存；`pending-and-issues.md:10692` 已记"第二套档位词表且有产码读者"）。只登记。
9. `frontend/src/components/showcase.tsx:85`：`desc: "检查器风格：mono 键名 + 控件 + 生效徽标，唯一活控制是面板不透明度。"` —— showcase 页描述句，不指任何字段。只登记。

**三档合计：① 2 枚／② 4 枚／③ 3 枚。**

---

## 2. "要重启"到达面（逐跳）

一句话结论：**到不了用户的眼睛（WebView 页面）；断在最后一跳——页面根本没有任何一处读回执，也没有 `config.set` 的出向调用。** 但**终端面有两处能到人眼**，其中一处（重启档那句）射程不含 endpoint。

| 跳 | 有／无 | 凭据（`HEAD:<file>:<line>`，同一次取数） |
|---|---|---|
| ① Go 侧回执构造 | **有** | `internal/panel/config_handlers.go:435` `return "这一项要重启进程并重新运行才算用上（同一次运行里模型通路不会重建）。"`；拼装处 `:381-396` `renderSettingReceipt`（`:383` "设置已写入。"＋`:387-389` 键路径＋`:391` `" " + tierSentence(r.Tier)`）；档位来源 `cmd/wisp/panel_config_store.go:228`／`:275`（无条件 `EffectiveRestart`）。禁句自述在 `:377-380` 注释 |
| ② 过 `internal/panel` 方法名册 | **有** | `internal/panel/bridge.go:66-67` `MethodConfigGet = "config.get"`／`MethodConfigSet = "config.set"`（已在白名单）；`internal/panel/composer_dispatch.go:155-166` `Handle` 把 settings 那条的 receipt 原样返回；`:147-154` 注释逐字给出理由「ticket 248 AC#8 requires "when this takes effect" to reach a page-visible surface」；`:174-176` "only the settings doors ever fill it"；宿主侧 `cmd/wisp/panel_host_windows.go:802-804` `w.Bind(panelDispatchBinding, func(raw string) string { reply, _ := m.dispatchRaw(ctx, raw); return reply })`，binding 名 `:80` = `"wispDispatch"`，`dispatchRaw` `:814-821` |
| ③ 页面侧读取处 | **无（断点）** | `frontend/src/lib/panel.ts:140-142` 桥接口只有 `postMessage(message: string): void`；`:163-165` 逐字「Send one bridge request and return nothing: responses arrive as a fresh PanelSnapshot push, never as a return value」；`:211` `function sendRequest(method: string, payload: …): void`；出向调用只有四枚 `:246` `panel.mode.request`／`:255` `panel.workspace.request`／`:272` `panel.attachment.add`／`:300` `panel.message.send` ⇒ **无 `config.set`／`config.get`**；快照字段表 `internal/panel/composer.go:58-91`（`pending`／`results`／`composer`／`generatedAt`＋两枚可选 `instructions` `:74`／`tasks` `:91`）**没有任何回执／tier 字段** |
| 补：诊断命令面（人手起的 CLI，不是产品面） | **有** | `cmd/wisp/panel_inbound.go:182` `fmt.Fprintf(s.stdout, "wisp panel-inbound: 第 %d 行回执：%s\n", lineNo+1, reply)`——`:177-180` 注释明写这是为 AC#8 才印 Handle 原话；`-data` 缺失时 `:136` 那句拒绝也是给人看的 |
| 补：手改 `config.toml` 的 reload 面（终端，人可读） | **有，但不覆盖 endpoint** | `cmd/wisp/config_reload.go:326`（重启档那句，见 §1c-3）＋`:240`（`[fs]` 那句）；`[llm]` 走的是 `:191` 那句"已立即生效"（§1①-1） |

`design/**` 一个字节都没读、没转述（越界自报见 §5 第 3 条）。

---

## 3. 给落地腿的最小形状（不写代码、不动字节）

**4 枚文件，⛔ 不需要新入向方法名，⛔ 零重建跳。**

1. `cmd/wisp/panel_config_store.go`（`:228`／`:275` 两行硬填）——把档位改成问**已有**那张登记表：`internal/config.TierOf`（`internal/config/tiers.go:93`），产码先例＝`cmd/wisp/config_readers_255.go:236`。这一改让"回执说什么档"和"热加载引擎怎么归档"同源，正面解掉 §1d-7 那一问的一半。
2. `internal/panel/config_handlers.go`（`tierSentence` `:426-441`）——`EffectiveNow` 那一支（`:431`）在名册里今天没有任何字段撑得起：要么删，要么只允许由 ① 的登记表证明才印。⛔ 不许"为了让两枚零写者的值看起来活着"而顺手启用它们（台账 `pending-and-issues.md:10750` 已把这写成禁句）。
3. `cmd/wisp/config_readers_255.go`（`:103` 的 `llm` 行）＋ `cmd/wisp/config_reload.go`（`:191` 那句）——这是 §1①-1 的落点，**但它撞票 255 AC#1／AC#5 的地界**（`pending-and-issues.md:10750` 给票 255 新补的 AC#5 判据正是"回执那句'什么时候生效'必须由登记表产出"）。⇒ 本腿建议：**由编排者先裁 `[llm]` 归哪一档，再动这两枚文件**，不许落地腿自己降格。
4. 页面可见面（AC#8 的"到得了眼睛"那一半）：缺的不是名字、是**出向投递形状**。两条现成先例，都不新增入向名：
   - 把回执并进快照的可选段（`internal/panel/composer.go:74`／`:91` 已是"可选段"先例；`票 248 J3` 裁过"不加顶层第五键"），页面读取处随之在 `frontend/src/lib/panel.ts` 的 `PanelSnapshot`（`:130-137`）里多一个字段；
   - 或等票 33 的 Go→页 push 那一形落定（`panel_host_windows.go:828-835` 逐字写明 `firstRoundTrip` 只证"页面到得了 Go"，**不是** AC#14 的回执那一跳的证据）。
   ⛔ 本编队不写 `frontend/**`（票 248 AC#9 同口径），这一枚只登记形状与地界。

**怎么绕开票 223 地界**：以上四形改的都是**"说什么"和"说到哪儿"**——没有一处新建 `llm.NewResolver`／endpoint、没有往 `OnReload` 挂重建、没有新增重建那一跳。票 223 的地界只在"把重建接进热加载"时才碰；本形状零沾。⚠ 也别顺手把 `frontend/dist` 里那枚 `.gitkeep` 当产物（`git ls-files frontend` 现量仍只有它，票 248 AC#9 那一格没勾 ⇒ 即便页面写了读取处，今天也发不出去）。

---

## 4. 三档计数（回报用）

- ① 2 枚：`cmd/wisp/config_reload.go:191`（`[llm]` 真会印"已立即生效"）／`internal/panel/config_handlers.go:431`（回执里那句"这一项立即生效。"，今天不可达）
- ② 4 枚：`config_reload.go:326`／`config_reload.go:194-197`／`config-screen.tsx:99-100`＋`:114-115`／`config_reload.go:240`
- ③ 3 枚：`248-settings-write-path-v1.md:152`＋`:254`／票 `:126`（两套词表待裁）／`internal/config/schema.go:34,37`（注释）／`frontend/src/components/showcase.tsx:85`
- "要重启"到达页面：**无**，断在 `frontend/src/lib/panel.ts:211`＋`:140-142`＋`internal/panel/composer.go:58-91`

---

## 5. 顶回与欠账

1. **顶回派单第 2 段那句"这一格要的是名册：今天有几处这么说"**：名册我给了（①＝2 枚），但**这一格不是没被人看过**——票面 `:126` 已写「AC#8＝附条件成立」，`docs/evidence/s1/248-settings-write-path-v1.md:146/152/254` 已裁过并把"两套词表"那一问交回编排者。没勾的真实原因是**待裁项**，不是文案没人查。本腿的增量＝①-1 那枚 stdout（`[llm]` 真会印"已立即生效"），v1 当时没量到它。
2. **顶回派单的"逐枚判三档"对 465 枚 `立即生效`**：其中 138 枚文件在 `.scratch/wisp/**`（probes 变异拷贝＋工单＋台账），不可逐枚过目、也不是"用户会读到的字面"。我把判据射程收在"产码非测试字符串＋页面文案"，其余按目录计数登记。⛔ 派单没给"probes 里的 baseline／mutation 拷贝算不算名册"的口径——我的处置是**不算**（它们是别票的夹具副本，如 `.scratch/wisp/probes/231/r1/mutation/config_reload_no_branch.go:187` 是 `cmd/wisp/config_reload.go:191` 的拷贝）。如果编排者要把它们计入，名册要另开一轮。
3. **自报越界**：一次 `git grep` 的 pathspec 里我误写了 `design`（那一条是 `… HEAD -- 'cmd' 'internal' 'design' '…'`）。输出零命中、我没读到也没转述任何 `design/**` 字节，但按"一个字节都不许读"的字面，检索面我违了一次；后续取数已把该 pathspec 去掉。
4. **量不到的格子**：
   - `[llm]` 会不会进 `rep.Restart`（从而吃到 `:326` 那句）——读码结论是不会（`internal/config/tiers.go:30` 记 `hot`），但没跑（禁 go 面归 `282-r1`）。
   - 票 255 AC#5（"回执那句必须由登记表产出"）是否已经把面板回执纳进同源守卫——我只量到 `TierOf` 的产码读者一处（`config_readers_255.go:236`），`panel_config_store.go` 里没有。
   - 页面那半（`frontend/**` 是否真会在 WebView 里显示回执）＝`dist` 今天无产物 ⇒ 端到端不可测，属票 248 AC#9／票 33 地界，不是本腿欠账。
