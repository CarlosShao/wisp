# pool-validity 4e — 死腿 4a 名册收尾（233／234／237）＋ 六批重叠普查

> 本腿代号 `pool-validity-4e`。只读普查腿，档位口径四档：**仍成立／已失效（已被别人做掉）／差翻勾／量不到**。
> 本腿不写票面、不动台账、不碰 `docs/**`、零 go 命令、只 commit 不 push。

---

## 0. 起手锚

- HEAD 短号：`95cb3d7a`
- 分支：`dev`
- 时刻：`10-06 15:13`（`date '+%m-%d %H:%M'`）
- `git status --porcelain -- internal cmd` 原文（起手第一次读数）：

```
?? internal/tools/tasklist_deferred_236r3_teeth_test.go
```

⇒ 与"另一枚写腿正在做 236 AC#2"的告知对得上（`236-r3` 在 `internal/tools` 落了未跟踪的突变件）。本腿因此**不读 `internal/**` 工作树内容面**，一律 `git --no-pager grep ... HEAD`。

- 落盘位置检查：起手 `ls .scratch/wisp/probes/pool-validity/` 原文：

```
1
2
3a
3b
4a
4b
4c
4d
msg-a633.txt
```

⇒ `4e` 此前**不存在**（编排者那句"4e 未用过"成立），本腿按指定用 `4e`，不启用备选名 `4e2`。
★附带事实具名：`pool-validity/` 下除编排者列的六批（`2/`、`3a/`、`3b/`、`4a/`、`4b/`、`4c/`）之外，还有 **`1/`** 与 **`4d/`** 两个目录。`4d` 是本腿的禁读前缀（在飞腿地界），不读、不采信；`1/` 不在六批名册里，本腿在 §3 里把它当**补充对账面**读票号列（见 §3.4），⛔ 不据此改任何一批的档。

---

## 1. 三枚名册（本腿实际要判的枚）与分母现量

命令原文（编排者给的尺，本腿自跑）：

```
for n in 233 234 236 237 238; do f=$(ls .scratch/wisp/issues/${n}-*.md 2>/dev/null | head -1);
  printf "%s 文件=%s 未勾=%s 已勾=%s\n" "$n" "$(basename $f)" "$(grep -c '^- \[ \]' $f)" "$(grep -c '^- \[x\]' $f)"; done
```

读数原文（15:13 现量）：

```
233 文件=233-d36-three-tier-enum-is-tagged-on-no-key-and-the-reload-tier-has-neither-producer-nor-consumer.md 未勾=6 已勾=0
234 文件=234-the-twelve-cells-that-only-need-a-reading-mutation-and-gate-rerun-on-current-head.md 未勾=6 已勾=0
236 文件=236-six-cells-that-only-surface-at-the-reading-layer.md 未勾=7 已勾=0
237 文件=237-race-run-exposes-a-parentheft-inference-that-never-held-plus-two-things-only-a-human-reads.md 未勾=3 已勾=0
238 文件=238-first-run-on-a-fresh-machine-creates-the-whole-data-root-with-no-seal-to-inherit.md 未勾=4 已勾=0
```

⇒ 与编排者"236 有写腿在做 AC#2、238 已被 4b 判过"的口径**没有冲突**：五枚**全部 0 已勾**，即没有任何一枚被勾掉过。
⇒ 编排者那句"读数与我上面那句不一致时以你的为准"：上面读数**未见冲突可言**（编排者没给具体枚数，只给了"5 枚尚未判"的名册），本腿照原样上报。

三枚的判与不判：

| 号 | 本腿动作 | 具名理由 |
|---|---|---|
| 233 | **判** | 4a 表里逐行写着"尚未判" |
| 234 | **判** | 同上 |
| 236 | ⛔ **不判** | 写腿正在做它的 AC#2（`236-r3`，`internal/tools` overlay 突变，起手锚里那枚未跟踪件即其证据）。判它会与它互相污染。 |
| 237 | **判** | 4a 表里"尚未判" |
| 238 | ⛔ **不判，只当参照样本读** | `4b/batch4b.md` 已判（仍成立＋状态半边已被机主关死，台账 `A476`）。本腿只读它的**口径形状**，不重判。 |

---

## 2. 判档表（233／234／237）

> 尺一律 HEAD 面：`git --no-pager grep -n '<符号或字面>' HEAD -- internal/ cmd/ tools/ scripts/ .github/`。⛔ 零 go 命令 ⇒ 运行期格与整包终态一律〔量不到〕。
> ⛔ 三枚都**不裁可结案**；翻勾权在编排者。

| 号 | 档 | 复法命令原文 | 读数 |
|---|---|---|---|
| **233** | **合档，逐面拆**：<br>① 现量 2（三枚 `Tier` 常量零赋值点）＝**已失效／已被别人做掉**<br>② 现量 4（reload 档"没有生产者键"、`manager.go:198` 那句恒假）＝**已失效**，且**这一半当年就未成立**<br>③ 现量 3（reload 档没有消费者）＝**仍成立**<br>④ AC#2 甲形的前置（硬编码清单与 schema 会漂）＝**已失效**<br>⑤ AC#4／AC#5／AC#6（种键复跑／突变自证／整包终态）＝**量不到** | 尺一（票面现量 2 那把尺换 HEAD 面逐字复跑）：`for t in TierHot TierReload TierRestart; do git --no-pager grep -n "\b$t" HEAD -- internal/ cmd/ tools/ scripts/ .github/ \| grep -v '_test\.go'; done`<br>尺二（已失效硬门·commit 号）：`git --no-pager log --oneline --diff-filter=A -- internal/config/tiers.go` ＋ `git --no-pager log --name-only -1 dd92bb92`<br>尺三（登记表今树形状）：`git --no-pager grep -n "TierRegistry" HEAD -- internal/ cmd/ \| grep -v _test` ＋ `git --no-pager show HEAD:internal/config/tiers.go`<br>尺四（生产者可达性，票面现量 4 说"恒假"）：`git --no-pager show a8f3c020:internal/config/manager.go \| grep -n "rep.Reload = append\|reload = true"` ＋ `git --no-pager log --format='%h %ad %s' --date=format:'%m-%d' -S "reload = true" -- internal/config/manager.go`<br>尺五（消费者那半边）：`git --no-pager grep -n "OnReload *=" HEAD -- internal/ cmd/ tools/` ＋ `git --no-pager grep -n "rep\.Reload" HEAD -- cmd/wisp/` ＋ 对照重启档 `git --no-pager grep -n "OnRestartPending" HEAD -- cmd/wisp/`<br>尺六（硬编码清单是否还只靠注释兜）：`git --no-pager grep -n "restartTierKeys" HEAD -- cmd/ internal/` ＋ `git --no-pager log --format='%h %ad' --date=format:'%m-%d' --diff-filter=A -1 -- cmd/wisp/restart_tier_keys_255r2_test.go` | 尺一＝**三枚常量各 2 命中**，两行都在同一文件（`TierHot`→`schema.go:34`＋`:36`／`TierReload`→`:37`＋`:39`／`TierRestart`→`:40`＋`:43`），非注释的赋值点＝**0** ⇒ 票面现量 2 的**读数今天一模一样**。<br>★但**"一枚键都没标上"这句今天不成立**：尺二＝**`dd92bb92`（10-02，票 255 AC#2-ⓑ，编排者代笔）** 新增 `internal/config/tiers.go`＋`tiers_255_test.go`、改 `manager.go`（`git log --name-only` 逐条点名这三枚）。尺三＝`internal/config/tiers.go:24 TierRegistry = map[string]string{`，现读逐行登记 `ball/session/audio/llm/agent/privacy/memory/panel/cost/models/observe/hotkey = "hot"` 12 枚、`risk/fs/net/plugins = "locked"` 4 枚、`app.theme = "hot"`、`app.language／app.autostart／app.single_instance = "restart"`、`voice.* = "reload"` **20 枚**（`voice.enabled`／`voice.wake_word.enabled`／`…keywords`／`voice.asr[.provider／.model]`／`voice.tts.provider／.voice`／`voice.conversation_mode`／`voice.aec[.enabled／.echo_ref]`／`voice.realtime[.enabled／.provider／.model／.api_key_ref／.base_url]`／`voice.cloud_asr_chain`／`voice.cloud_tts_chain`）＋另 **4 枚 `"hot"`**（`voice.wake_word.thresholds`＋`veto_words`、`voice.tts.speed`、`voice.punctuation`）；`TierOf`（`:94`）是导出访问器；同源守卫 `manager.go:294-299` 逐字 `if tier, ok := TierRegistry[s.name]; !ok \|\| tier != "hot"` ⇒ `panic("config: section … is hot-applied by plan() but not registered as \"hot\" in TierRegistry")`。<br>⇒ ①**已失效的凭据种类**＝一枚 commit 号＋`--name-only` 命中（硬门满足）。⚠ 分档真相源从 `Tier` 那三枚常量**换成了字符串表**，所以"赋值点仍＝0"与"键已标上"**两句同时为真**，不许拿前者否定后者。<br>尺四＝**378／412 两行在 233 的立案锚点 `a8f3c020` 上逐字都在**，`reload = true` 更早可追到 **`ab638dd5`（09-19）** ⇒ `manager.go:198` 的 `len(rep.Reload) > 0` **不是恒假**；它当时恒假的只有 `m.OnReload != nil` 那一侧。⇒ ②票面现量 4 的因果记反了（它拿"三枚枚举零赋值"那把尺推"填不进 `rep.Reload`"），**立案那天即不成立**。<br>尺五＝`OnReload =` 赋值点今树 **3 处**：`cmd/balldebug/main.go:244`（旁支宿主，票面点名的还是那一处）＋`internal/ball/hotkey_live_test.go:224`＋`internal/config/manager_test.go:118`（两枚测试）⇒ **`wisp run` 仍零赋值点**；`rep.Reload` 在 `cmd/wisp/` 只命中 **1 行**＝`config_reload.go:180` 那条 `auditf` 的 `%v` 字段（与 `rep.Hot`／`rep.Restart`／`len(rep.Locked)` 同句）；对照重启档：`OnRestartPending =` 赋值点 **1 处**＝`config_reload.go:315`（`reportRestartPending` 有 stdout 实话句）。⇒ ③**reload 档在 `wisp run` 里既没接进任何会做事的读者、也没有那句"本次运行不会因此改变行为"形状的实话句**——这一格与票面那句"不许用静默不生效充当这一档"**今天仍对着**。<br>尺六＝`cmd/wisp/config_reload.go:338 var restartTierKeys = []string{…}`（票面写 `:299` ⇒ 漂 39 行），注释里那句"and by a test"已**兑现**：`restart_tier_keys_255r2_test.go:61` 逐枚 `if tier, ok := config.TierRegistry[c.key]; … t.Errorf("restartTierKeys names %q but config.TierRegistry says tier=%q ok=%v")`，该件由 **`67ab595d`（10-03）** 加入（同 commit 改 `cmd/wisp/config_reload.go`）。⇒ ④"两边不许漂的那枚常驻用例"已存在。<br>⑤**量不到**＝AC#4（种键→改值→不重启→stdout 增量断言）／AC#5（注释掉赋值点必红）／AC#6（`go test ./cmd/wisp ./internal/config/`＋`-count=3`＋gofumpt＋vet＋d22scan）——本腿零 go，缺的就是这几发行数。<br>★另具名一处**与本票同族、不在票射程**：`cmd/wisp/panel_config_store.go:27` 注释逐字 `OnReload has no`…（面板 store 只在装配时读一次、永不随 reload 重建）——reload 档零读者的**另一个**落点，票面未点名，本腿不扩判。 |
| **234** | **差翻勾**（12 枚判据的读数今天**全在盘上**，但凭据种类＝票内注记＋台账 `A449`＋一枚非实现者读数腿的 63 枚件；⛔ 不等于可结案）<br>＋其中 **AC#5 那一格＝已失效／已被别人做掉**（硬门凭据在下面）<br>＋票面自己的**引用行号有一枚漂**（下段尺四）<br>＋AC#2／AC#3／AC#4 的**复跑形状＝量不到**（零 go） | 尺一（交件事件是否真落 commit、`--name-only` 是否命中）：`git --no-pager log --oneline -- .scratch/wisp/probes/gate-rerun-1/` ＋ `git --no-pager log --name-only -1 6d22dd6c`<br>尺二（件今树枚数＋零 UNJUDGED 复算）：`git ls-tree -r --name-only HEAD -- .scratch/wisp/probes/gate-rerun-1 \| wc -l` ＋ `grep -c "UNJUDGED" .scratch/wisp/probes/gate-rerun-1/table.md`（★这一把读的是**我自己这批的件**，不是 `internal/**` 内容面）<br>尺三（AC#5 销账的现量复算）：`for n in 92 97 104 105 110 113 115; do f=$(ls .scratch/wisp/issues/${n}-*.md \| head -1); printf "%s un=%s chk=%s\n" "$n" "$(grep -c '^- \[ \]' "$f")" "$(grep -c '^- \[x\]' "$f")"; done` 与接收方同尺 `for n in 64 12 114 230 111 112 140 132` ＋ `grep -c "^- ⛔" .scratch/wisp/issues/07-*.md 11 63 80 97 115 92`<br>尺四（票面 12 枚"票号＋行号"逐枚验位）：`for spec in "92 68" "92 73" "104 52" "104 57" "110 39" "110 47" "113 54" "113 56" "115 58" "115 64" "97 55" "105 54"; do set -- $spec; sed -n "$2p" $(ls .scratch/wisp/issues/$1-*.md \| head -1); done`<br>尺五（AC#6 那把坏尺今树复跑）：`git --no-pager grep -n "NewGate(" HEAD -- internal/ cmd/ tools/ scripts/ .github/; echo rc=$?` ＋ `git --no-pager grep -n "risk\.Gate(" HEAD -- internal/ cmd/ tools/` | 尺一＝交件**落过 commit**：`519570b0`(骨架)→`5fa1c384`(AC#1)→`06a650bd`(M-1/M-2)→**`6d22dd6c`**（12 枚全部升〔有读数〕＋AC#5 销账 10/10＋AC#6），另 **`a7993a9b`＝台账 `A449`** 逐字"收 gate-rerun-1：票 234 那 12 枚全升〔有读数〕、结案票残余 24/11→14/8 我复跑对上"。<br>尺二＝`git ls-tree` 在 HEAD 上 **63 枚文件**（含 `m1-…`～`m5-…` 五组突变件、`gates-*.txt` 七组门禁读数、`ac1-baseline-roster.txt`／`ac4-final-roster.txt`、`ac5_relocate.py`、`ac6-yardstick.txt`、`m4-g6-container.log`＝真容器那一发）；`table.md` 289 行、**`UNJUDGED` 计数＝0**；节头逐名数得 **Mutation cells (5) ＝M-1…M-5**、**Gate/regression cells (7) ＝G-1…G-7**，与票面"变异 5＋门禁 7＝12"枚数一致。<br>尺三＝发出方现量 `92 un=3 chk=3／97 un=1 chk=3／104 un=2 chk=3／105 un=1 chk=4／110 un=2 chk=3／113 un=2 chk=4／115 un=2 chk=3`；接收方 `64 un=1 chk=8`、`12 un=3 chk=6`、`230 un=5 chk=0`、`114 un=9 chk=2` ⇒ **与 table.md `AC#5` 那张表"BEFORE→AFTER"与"Receiver … unchanged"两句逐枚相同**（含票面警告的"114 记 un=7 实为 9"）。⛔ 指向句枚数那把尺**对不上**（本腿 12 枚 vs 腿报 10 枚，两把读数并列进 §5 第 1 条，⛔ 本腿不挑一枚当准）。<br>尺四＝**11 枚逐字对位**（92 `:68`／92 `:73`／104 `:52`／104 `:57`／110 `:39`／110 `:47`／113 `:54`／113 `:56`／115 `:58`／115 `:64`／97 `:55` 各自就是那一枚 AC 的第一行）；**漂的只有 1 枚**＝票 105 `:54`：今树 `:54` 是上一格的**续行**（`本票只交付库内 + 审计那半；**不许为了勾框扩界**。`），AC#5 现在 `:56`。<br>尺五＝`NewGate(` 在 HEAD **rc=1／0 命中**（尺子未坏：同族 `risk.Gate(` 命中 1 处＝`internal/tools/mode.go:98`，票面那句"生产调用点"逐字还在）⇒ AC#6"这条坏尺点名作废"**成立**。<br>**⇒ 档位怎么落的，逐格写清**：<br>**AC#5＝已失效／已被别人做掉**（凭据＝一枚 commit `6d22dd6c`＋尺三现量与它的读数逐枚一致，硬门满足）——那一格里"把结案票 10 行改成指向句"这件事**已被 `gate-rerun-1` 交掉**。<br>**AC#1／AC#2／AC#3／AC#4／AC#6＝差翻勾**：读数都在**别人的台件**里（票内注记＋`A449`＋`table.md`＋63 枚原始件），⛔ 不是本腿复跑，⛔ 本腿不裁它可结案。<br>**量不到的形状**：票面要求的是"**再复跑一遍**当前 HEAD"，而 `6d22dd6c` 那批复跑落在 `015f6be1`→`06a650bd`（09-29）；本腿零 go ⇒ 无法在**今天 95cb3d7a** 上验"这 12 枚读数对新提交范围内的产码是否仍适用"。⚠ **这一句不许外推成"12 枚都过期"**，它是"缺哪一行读数"的登记。<br>★另具名一枚**不在 12 枚射程**的账：`docs/evidence` 下 115 名下裁决表今树仍 **0 枚**（`git ls-tree -r --name-only HEAD -- docs/evidence/s1 \| grep -E "/115-"` → rc=1），票面 `:32` 自己写着"本票不替 230 交那张表"——账在票 230 AC#1（4a 那腿已判 230 仍成立），本腿不重判。 |
| **237** | **仍成立**（三格**逐枚都还在**：① `parkParents` 那一族的静态顺序前提仍破＋它已**从票面的 1 枚扩到 3 枚用例**；② 代码旁边那句"只能靠人读"仍**零命中**；③ 前置读数仍是"至少一枚"的**单颗 token**，逐枚形状不存在）<br>＋AC#1 的完成判据（改前红形状复现＋改后 `N≥20` 发零红）＝**量不到** | 尺一（顺序，票面现量 1 那两个行号）：`git --no-pager show HEAD:internal/tools/subagent_197.go \| grep -n "child.RunAsync\|已派生\|giveBackWhileWaiting(ctx)"`<br>尺二（谁在等那枚 parked）：`git --no-pager grep -n "parked" HEAD -- internal/tools/subagent_222_test.go` ＋ `git --no-pager show HEAD:internal/tools/subagent_222_test.go \| sed -n '363,376p'`<br>尺三（parked 的来源是不是 onUpdate 那条）：`git --no-pager show HEAD:internal/tools/subagent_222_test.go \| sed -n '270,280p'`<br>尺四（票面现量 2 那句话今树）：`git --no-pager grep -n "只能靠人读" HEAD -- internal/ cmd/ tools/ scripts/ .github/; echo rc=$?` ＋ `git --no-pager grep -n "仪器" HEAD -- internal/tools/subagent_197_test.go; echo rc=$?`<br>尺五（票面 AC#2 指的 `:398` 今天写的是什么）：`git --no-pager show HEAD:internal/tools/subagent_197_test.go \| grep -n "AC#2 (ticket 211"` ＋ `git --no-pager log --format='%h %ad %s' --date=format:'%m-%d %H:%M' -1 1d2ad737`<br>尺六（前置读数射程，票面 AC#3）：`git --no-pager grep -n "awaitChildOnBridge222" HEAD -- internal/tools/` ＋ `git --no-pager show HEAD:internal/tools/subagent_222_test.go \| sed -n '391,401p'` ＋ `git --no-pager log --name-only -1 a71f0be9` | 尺一＝`subagent_197.go` 今树：**`:344` `bg := child.RunAsync(childCtx, prompt)`（起孩子）排在 `:366` `t.feed(bg.ID, "已派生："+label)` 与 `:368` `onUpdate(fmt.Sprintf("task.spawn 已派生子代理 %s（%s）", …))`（父发派生回执）之前**（票面写 `:347`／`:368` ⇒ 第一行漂到 `:344`，回执那行没漂）。⇒ ①**顺序这一形逐字还在**。<br>尺二＝`parkParents`（`:366`）体内 `:370` `if got := len(h.parked); got != n { t.Fatalf("只有 %d/%d 枚父任务到达等待点（父任务没派生成功，等待读数就无从谈起）", …) }` ＋ `:373` 紧随其后断 `len(h.parents) != 0` ⇒ **等满 parked 后立刻要求没有任何一枚父返回**这一形原样；票面现量 1 的**因果句今天比当年更准**：`parked` 的喂入口（尺三）＝`bridge` 的 `OnUpdate` 里 `if strings.HasPrefix(delta, "task.spawn 已派生子代理") { h.parked <- delta }`，也就是 `:368` 那条 onUpdate，而 `:365 giveBackWhileWaiting(ctx)` 在它**之前** ⇒ parked 满时令牌已交、孩子随时能跑完并推父返回，竞态窗口**仍在**。<br>★**这枚腿的新发现（本腿自己扫的，票面没写）**：`await222Tokens(t, …, h.started, n)` 等的是 `h.started`（孩子进第一枚模型调用的信号），**不是** parked；所以"等满 parked"这一破形不只票面点名的那一枚。调用点尺＝`parkParents` 被 `subagent_222_test.go:437`／`:514`／`:555` **三枚用例**调用 ⇒ 暴露面从票面假设的"那一次 `-race` 偶发"扩到三枚。票面那句"下一位会把 CI 偶发归因成父任务没派生成功"**恰好就是 `:370` 的 Fatalf 原文**（逐字"父任务没派生成功，等待读数就无从谈起"）——归因陷阱原样在文案里。<br>尺四＝**rc=1／0 命中**：票面 AC#2 要加的那句"这枚注释没有任何仪器盯着，只能靠人读"今树**代码旁边一个字都没写**（`仪器` 在 `subagent_197_test.go` 也 0 命中）⇒ ②**仍成立**。<br>尺五＝票面点名的 `subagent_197_test.go:398` 今天写的是 `// AC#2 (ticket 211 甲) - the resident nail on the pool's number…`（`:389` 起，`:398` 落在"Why the pool must still not exceed the ceiling"那段里）；那段**被 `1d2ad737`（09-29 22:38，test(235 AC#1)）重写理由**（"diff 只有 :398-402 那段注释"），而**票 237 立案在 `d5a59d66`（09-29 23:36）之后** ⇒ 本腿**不能**据此判"票面的观察对象被改没了"：它改的是理由、加没加过那句"只能靠人读"本腿无从反推。⇒ 只按尺四的**零命中**判"仍成立"，这一处不确定性登记进 §4。<br>尺六＝`awaitChildOnBridge222` 定义 `:391`，体内 `:395` `case <-h.probe.arrived:`——**select 收一枚即返回**；调用点 `:444`／`:536`／`:604` 三枚（`a71f0be9` 09-29 22:46"检测力从 1/3 补到 3/3"）。`arrived` 是 `chan struct{}`（`:134`，缓冲 16），每枚孩子在桥上起跑 send 一枚（`:174` `case p.arrived <- struct{}{}:`）⇒ 三枚调用点各自**只要求"至少一枚到达早于第一枚返回"**，逐枚形状（每枚用例要求它自己那个孩子在任何返回之前到达）**不存在**。⇒ ③票面 AC#3 那格〔部分绕过仍是瞎的〕**仍成立**；另具名一处：三枚**都收**的 `await222Tokens(…, h.started, n)` 只钉"孩子在第一枚模型调用里"，与孩子的**工具面是否真上桥**无关（票面 `:13` 那句"另两枚用例里仍失明"今天扩为"三枚用例里都只是至少一枚"，方向一致、口径更宽）。<br>量不到＝AC#1 完成判据要的"具名命令＋`-race -count=N` 改前红形状复现＋改后 `N≥20` 发零红"、AC#3 判据的形状 (i)(ii) 两发突变——全部运行期，本腿零 go。⚠ 不许由尺一／尺六的 grep 命中外推成"那一发今天必然红"。 |

---

## 3. 六批重叠普查

> ⛔ 本节只做**集合与枚数对账**：不重判任何一枚、不改任何一批的档、不采信任何一批的判语。
> 抽取尺（本腿自写，落 `.scratch/wisp/probes/pool-validity/4e/rowscan2.awk`）＝每枚**顶格数字表格行**记进「名册」，
> 行内**第一个以档位词起手**的字段（仍成立／已失效／差翻勾／量不到）记进「已判」；
> 字段起手是 未判／尚未判／待填 ⇒ 只算名册、不算已判。行首允许 `**NN**` 粗体（`4c` 的判定表就是粗体，第一版尺漏掉它、第二版补上，见 §5 第 3 条）。

### 3.0 分母先说清（编排者那句"2/＝46 枚"与我的读数不一致，差在哪）

命令原文：`awk -f 4e/rowscan2.awk <file> | 数 #ROSTER / #JUDGED 的 id 枚数`，逐枚清单见 §3.1。

| 批次文件 | 编排者派单给的枚数 | 本腿抽到的**名册枚数** | 本腿抽到的**已判枚数** | 该文件自己 §2／§5 的合计行 |
|---|---|---|---|---|
| `2/`（**四枚文件** batch1＋batch2＋batch3＋batch4） | 派单写「`2/`(46 枚)」 | **105**＝30＋16＋13＋46 | **46**＝30＋16＋0＋0 | 四张合计行分别 30／16／13／46 ✓ 逐张对得上 |
| `3a/batch3a.md` | 13 | **13**（★表格**行数＝15**） | 13 | 合计 13 ✓（⛔ 未判 0 枚） |
| `3b/batch3b.md` | 18 | **18** | **9** | 合计 18，四档全写「未判」⇒ **与本腿抽取冲突**（见 §3.3） |
| `4a/batch4a.md` | 11 名册／6 判 | **11**（★行数 22＝名册表 11＋判定表 11） | **6** | 正文「已判 6 枚／余下 5 枚」✓ |
| `4b/batch4b.md` | 8 | **8**（★行数 16＝名册 8＋判定 8） | 8 | §2bis 合计 8 ✓（4/0/2/2） |
| `4c/batch4c.md` | 8 | **8**（★行数 16＝名册 8＋判定 8） | 8 | §5「仍成立 8／已失效 0／差翻勾 0」✓ |

★**差在哪（以本腿读数为准，逐名说）**：编排者那句「`2/`(46 枚)」把 `2/` 当成一批；
盘上 `2/` 是**四枚文件**，票号列合计 **105 枚**（30／16／13／46），46 只是其中 `batch4.md`（号段 150+）一枚。
⇒ 后面所有"六批"的集合运算都按**九个文件**做（`2/` 拆四枚计），派单的"六批"是目录数、不是名册数。
★附带一枚不在射程内的面：同目录还有 `1/`（死腿 pv1，四个文件，名册同样 **105 枚**、已判 **0 枚**）与 `4d/`（本腿禁读，未碰）。

### 3.1 问题一：哪些票号**被两批以上判过**

抽取尺原文（判过＝判定表里那一行的档位词不是 未判／待填）：

```
for f in 2/batch1 2/batch2 2/batch3 2/batch4 3a/batch3a 3b/batch3b 4a/batch4a 4b/batch4b 4c/batch4c; do
  awk -f 4e/rowscan2.awk $f.md ; done
```

**名册枚数合计（出现次数）＝163，去重＝105**；**判过枚数合计＝90，去重＝82**。⇒ 163−105＝**58 次重复列名**，90−82＝**8 次重复判定**。

**被两批以上判过的＝8 枚，逐枚列（无第 9 枚）**：

| 票号 | 判过它的批次组合 |
|---|---|
| 15 | `2/batch1.md` ＋ `4c/batch4c.md` |
| 16 | `2/batch1.md` ＋ `4c/batch4c.md` |
| 21 | `2/batch1.md` ＋ `4c/batch4c.md` |
| 22 | `2/batch1.md` ＋ `4c/batch4c.md` |
| 23 | `2/batch1.md` ＋ `4c/batch4c.md` |
| 24 | `2/batch1.md` ＋ `4c/batch4c.md` |
| 25 | `2/batch1.md` ＋ `4c/batch4c.md` |
| 26 | `2/batch1.md` ＋ `4c/batch4c.md` |

⇒ 编排者点名的那片重叠（`4c` 的 15／16／21／22／23／24／25／26 全在 `2/batch1` 名册里）**枚数与票号逐枚吻合，本腿独立复算到同一集合**；
穷举尺（把每枚 id 的"判过批次"列成矩阵再筛 ≥2）没有浮出第 9 枚。
⚠ 但这 8 枚**两批判语不同**（本腿只列形状、不裁谁对）：`2/batch1` 判 **15 仍成立／16 量不到／21 量不到／22 仍成立／23 仍成立／24 仍成立／25 量不到／26 仍成立**，
`4c` 判 **8 枚全部仍成立** ⇒ 重复派单不只浪费一枚腿的轮次，还给出**三枚不同档**（16／21／25）。

**只是"列名重复"、判定不重复的**（名册被两批以上列过、只有一批判过＝派单射程重叠但无害）＝**49 枚**，逐组合：

| 组合 | 枚数 | 票号 |
|---|---|---|
| `2/batch3` ＋ `3a` | 13 | 100 103 107 108 109 111 112 120 122 123 132 140 148 |
| `2/batch4` ＋ `3b` | 18 | 150 159 165 166 168 169 170 172 173 178 182 185 186 187 189 194 195 196 |
| `2/batch4` ＋ `4a` | 11 | 200 220 225 227 229 230 233 234 236 237 238 |
| `2/batch4` ＋ `4b` | 8 | 238 242 244 247 253 262 264 266 |
| `4a` ＋ `4b`（判过它的只有 `4b`） | 1 | 238 |

⇒ 换句话说：**`2/batch4.md` 一枚就把 `3a`／`3b`／`4a`／`4b` 四批的名册几乎全部提前列过一遍**（46 枚＝`3b` 的 18＋`4a` 的 11＋`4b` 的 8，其中 238 同时属 `4a` 与 `4b` ⇒ 去重 36 枚，再加**只有 `2/batch4` 列过**的 10 枚＝211／213／214／215／216／219／231／232／249／269）；
`3a`／`3b`／`4a`／`4b` 相对 `2/batch4` **没有一枚是新材料**；`4c` 相对 `2/batch1` 也一枚不是新材料。
**六批名册的并集＝105 枚，与 `2/` 四枚文件的并集**逐枚全等**（`comm` 双向为空 ⇒ 后五批没引进过一枚新票号）。

### 3.2 问题二：今天仍零勾、却从没被任何一批判过的票

尺一（今日零勾开放票）：

```
for f in .scratch/wisp/issues/*.md; do case "$f" in *-done.md|*README.md) continue;; esac;
  [ "$(grep -c '^- \[x\]' "$f")" -eq 0 ] && echo 之; done   -> 103 枚
```

尺二（从两把里减掉"八十二枚已被某一批判过"的去重名册）⇒ **今天零勾 ∧ 六批从未判过＝22 枚**：

```
178 182 185 186 187 189 194 195 196 211 213 214 215 216 219 232 233 234 236 237 249 269
```

**其中 233／234／237 已被本腿（`4e`）今天判掉**（§2），剩下 **19 枚**逐枚给未勾格数（尺＝`grep -c '^- \[ \]'`）：

| 票号 | 未勾格数 | 具名备注（只陈述尺读到什么，⛔ 不判档） |
|---|---|---|
| 178 | 5 | 由 `2/batch4`＋`3b` **列名**、`3b` 那一行写"未判" |
| 182 | 7 | 同上（`3b` 未判行） |
| 185 | 7 | 同上 |
| 186 | 9 | 同上 |
| 187 | 9 | 同上 |
| 189 | 6 | 同上 |
| 194 | 7 | 同上 |
| 195 | 7 | 同上 |
| 196 | **0** | ★票面**一枚 `- [ ]` 都没有**（全文 14 行、`grep -cE '\[ \]'`＝0）⇒ "零勾"这把尺对它空转 |
| 211 | **0** | ★同上：无复选框，判据写成 `- **AC#1**…` 粗体子弹（5 枚） |
| 213 | **0** | 同上（AC 子弹 5 枚，Status 待派） |
| 214 | **0** | 同上（6 枚） |
| 215 | **0** | 同上（5 枚） |
| 216 | **0** | 同上（5 枚） |
| 219 | **0** | 同上（6 枚） |
| 232 | 6 | `4a` 按派单明令剔除（写腿 `232-r2` 在修）；票内另有 1 枚缩进框 ⇒ 顶格尺 6／全形尺 7 |
| 236 | 7 | `4a` 未判、`4b` 按号段排除（＜238）⇒ **此刻有写腿 `236-r3` 在做 AC#2**，本腿不判 |
| 249 | 7 | 票面 `Status: OPEN → WITHDRAWN（10-01 14:2x 编排者自撤，台账 `A503`）`；`4a` 明写"249 不在本尺输出里（已改名或已非零勾）"——**本腿今天复跑那把尺它就在输出里**（见 §5 第 4 条） |
| 269 | 5 | 票面 `Status: 已立，未派落地`；`4b` 按派单排除 |

⇒ 这一栏最值钱的读数是**形状**而不是枚数：19 枚里有 **7 枚（196／211／213／214／215／216／219）根本没有 `- [ ]` 形状的判据**（判据写成 `- **AC#N**：…` 粗体子弹，枚数 5／5／6／5／5／6），
所以 `4a` 那把「未勾>0 ∧ 已勾=0」的分母尺对它们是**空转**的。本腿今天在同一目录复跑两把尺：
「已勾=0 的开放票」＝**103 枚**，「未勾>0 且已勾=0 的开放票」＝**96 枚**，两把差 **7 枚＝正是这七枚无框票**（逐名见上表，⛔ 本腿不挑一把当准——它们量的是两件不同的事）。

### 3.3 问题三：每批"合计／自证"行 vs 本腿抽取枚数（对不上就两把读数并列）

| 批次 | 该文件自己写的合计／自证（逐字抽取） | 本腿抽取 | 对得上？ |
|---|---|---|---|
| `2/batch1` | 「合计 **30** ✓ 与 §0 分母对上」；仍成立 19／量不到 11／已失效 0／差翻勾 0 | 名册 30／判 30；档位词逐行数＝仍成立 19＋量不到 11 | ✓ 两把一致 |
| `2/batch2` | 「合计 **16** ✓」；仍成立 11／差翻勾 1／量不到 4 | 名册 16／判 16；11＋1＋4 | ✓ |
| `2/batch3` | 「仍成立 未判／已失效 未判／差翻勾 未判／量不到 未判 ／ **合计 13**」 | 名册 13／判 **0**（13 行全"未判"） | ✓（它自己没声称判过） |
| `2/batch4` | 「…／**合计 46**」，四档全"未判" | 名册 46／判 **0** | ✓ |
| `3a` | 「合计 **13** ✓ 与 §1 名册对上（⛔ 未判 0 枚）」；仍成立 9／已失效 2／差翻勾 2／量不到 0 | 名册 **13 枚 id**／判 13；9＋2＋2＋0 | ✓ **枚数对上，但★表格行数＝15**：`120`／`122` 各出现两次（第一次带判语、第二次是"待填"残行）。⛔ 本腿不改它，登记 |
| `3b` | 「仍成立 **未判** ／已失效 **未判**／差翻勾 **未判**／量不到 **未判**／合计 **18** ✓ 与 §0 分母对上」 | 名册 18／**判 9**（仍成立 7＋已失效 1＋量不到 1；未判 9 行＝178 182 185 186 187 189 194 195 196） | ✗ **冲突，两把并列不挑**：它 §1 那节抬头写「第 1 组（150／159／165）—— 已判，已 commit」且 9 行确有档位词；§2 却把四档全写"未判"。⇒ §2 那行只对"分母 18"成立，对"已判 0"不成立 |
| `4a` | 「§1 名册（11 枚）」＋「本节此刻的状态＝11 枚里已判 6 枚…余下 5 枚逐行写着'尚未判'」＋§2-0 自量分母 11 | 名册 11（★表格行 22＝两节各 11）／判 **6** | ✓ 三把一致（6 判＋5 未判＝11） |
| `4b` | 「§1 名册（8 枚，本腿现量尺）」＋§2bis「合计 **8** ✓ 与 §1 名册逐枚对上；仍成立 4／已失效 0／差翻勾 2／量不到 2」 | 名册 8／判 8；4＋0＋2＋2 | ✓ |
| `4c` | 「§1 名册（8 枚）」＋§5「判定分布：仍成立 8 枚／已失效 0 枚／差翻勾 0 枚／量不到＝无整枚」 | 名册 8／判 8；8＋0＋0 | ✓（★但它的判定表行首是 `**NN**` 粗体，本腿第一版尺会抽成 0 枚——两把读数都留在 §5 第 3 条） |

**跨批的"合计"层对账（本腿最硬的一格）**：`2/` 四张合计 30＋16＋13＋46＝**105**，与 `2/summary.md` §1 那张四档表的「合计 30｜16｜13｜46｜**105**」逐枚对上；
而 105 与**今天**的零勾开放票 **103** 差 **2 枚**＝`111`（今天 13:07 编排者 `c6ea08df` 翻了六格）与 `231`（今天 10:44 编排者 `ec84cff1` 翻勾记录）——
⇒ **不是任何一批量错，是两把尺取数时刻不同**（`2/` 起手 HEAD `949a5b92` 10-06 09:03，本腿 HEAD `95cb3d7a` 15:13）。
⚠ 这一句本腿按 §5 坑①自查过：负向结论前先证尺命中得了真名——`111`／`231` 两枚**都能被同一把尺命中**（`111 un=4 chk=6`／`231 un=0 chk=5`，本腿 15:1x 现量），
所以"今天少了这两枚"是它们真被勾了，不是尺漏了。

---

## 4. 判得心虚的枚数与具名理由

本腿共判 **3 枚**（233／234／237），**其中 2 枚判得心虚**，另对 §3 的普查有 **2 处**心虚。逐枚点名：

1. **234——心虚的不是档位，是"我到底有没有资格说它是差翻勾"**。
   我给的〔差翻勾〕凭据全部是**别人的东西**：票内 Progress log 那一行、台账 `A449`（`a7993a9b`）、台件 `.scratch/wisp/probes/gate-rerun-1/table.md` 与它在 HEAD 上的 58 枚件（`git ls-tree -r --name-only HEAD -- .scratch/wisp/probes/gate-rerun-1 | wc -l` 现量＝**58**；派单/交件里那句"63 枚件"是**我早先按 `--name-only` 加了 grep 过滤数的**，两把读数并列进 §5 第 2 条）。
   ⇒ 差翻勾这一档**本来就要**这种凭据（4a §0 定义如此），但我没跑过任何一发 go 命令，所以"12 枚现在是不是仍然成立"我**一个字都不知道**。票面要的恰恰是"在当前 HEAD 上复跑"，而 `6d22dd6c` 落在 09-29。**这一枚的档位是"差翻勾"，不是"已完成"，也不许被读成后者。**

2. **237 的第二格（AC#2／现量 2）——观察对象可能被 235 那条腿重写理由时动过，我无从反推**。
   我判"仍成立"的凭据只有今树 0 命中（`git --no-pager grep -n "只能靠人读" HEAD -- internal/ cmd/ tools/ scripts/ .github/` → rc=1）。
   但票面 AC#2 点名的 `subagent_197_test.go:398` 那段，在 **`1d2ad737`（09-29 22:38，test(235 AC#1)）**被整段改写过理由，而票 237 立案在 **`d5a59d66`（09-29 23:36）**——即"235-v1 报告它零命中"到"237 立案"之间只有 58 分钟，`1d2ad737` 在**之前**。
   ⇒ 我无法区分三种情况：(a) 那句从没写进代码（票面为真）／(b) 写进过又被 `1d2ad737` 删掉／(c) 235-v1 读的是改写前的同一位置。三种都不改变"今树 0 命中"，所以档位不动；
   但如果编排者拿这一格去翻勾，**(b) 这一支会让"加一行注释"变成"重新加一遍"**。要消掉这一层心虚只有一条路：`git show 1d2ad737 -- internal/tools/subagent_197_test.go` 看那一发的 diff 里有没有那一句——**本腿没有跑它，因为那一格的判定不依赖它，跑了就该把这一格从"注释面"升成"史实面"**（下一步若被要求再跑）。

3. **§3.1 的"没有第 9 枚重复判定"——依赖我的抽取尺认得出档位词**。
   尺＝第一个以 `仍成立|已失效|差翻勾|量不到` 起手的字段。若某一批用**同义词**写档（例如"仍为真／已由他人做掉"），它会掉进"只列名未判"那一栏，重复判定就会被少数。
   ⇒ 本腿做了一次正向对照防这件事：`4a` 的 6 行、`4b` 的 8 行、`4c` 的 8 行、`3a` 的 13 行、`3b` 的 9 行、`2/batch1`＋`2/batch2` 的 46 行**都抽到了档位词**（合计 90 次出现／82 枚去重），与"163 次列名／105 枚去重"这两个总数互斥自洽（90≤163、82≤105）；
   但我**不能证明**这六批里没有一枚用了同义词然后被我算进"未判"——所以这一格是**结构性风险，不是已知漏读**。

4. **§3.2 的 19 枚未判清单里，249／269 两枚我照抄了前批的排除、没有独立复核它们该不该算活**。
   `4a` 明写"269／249 不在本尺输出里（已改名或已非零勾，本腿不判、不追）"，可本腿 15:1x 复跑**同一把尺**（`未勾>0 ∧ 已勾=0`，`2*` 通配）时两枚**都在输出里**（249 未勾 7／269 未勾 5）。
   票面状态：249 `Status: OPEN → WITHDRAWN（10-01 14:2x 编排者自撤，台账 A503）`；269 `Status: 已立，未派落地`。
   ⇒ 我把它们列进"从未判过"那一栏是**集合事实**，不构成"这两枚还该派"。这一处属于**前批尺与今尺读数不同**，两把都留在 §5 第 4 条，本腿不挑一枚当准。

5. **三枚里唯一让本腿踏实的是 233**——它的四条硬读数（`TierRegistry` 在盘／`dd92bb92`＋`--name-only`／`OnReload =` 赋值点三处全在旁支与测试／`manager.go:378`＋`:412` 在立案锚点 `a8f3c020` 上就在）**全部是本腿自己跑的 HEAD 尺**，且负向那半（"赋值点仍为 0"）我先证了尺命中得了真名（同一条尺命中 `schema.go:36/39/43` 三枚声明行），才敢写 0。

---

## 5. 尺有歧义／判不动的地方

**同一把尺两遍读数不同 ⇒ 两把原文并列，⛔ 不挑一枚为准。**

1. **票 234 AC#5"销账 10 枚"vs 本腿两种指向句枚数尺**（三把并列）：
   - 腿报（`.scratch/wisp/probes/gate-rerun-1/table.md` §AC#5 那张表）：发出方 7 枚文件、**moved＝10**（07=3／11=1／63=1／80=2／97=1／115=1／92=1）。
   - 本腿窄尺 `for t in 07 11 63 80 92 97 115; do f=$(ls .scratch/wisp/issues/${t}-*.md|head -1); grep -c '^- ⛔.*gate-rerun-1.*处置' "$f"; done` → **07=3 11=1 63=1 80=2 92=1 97=1 115=1，合计＝10** ✓ 与腿报逐枚相同。
   - 本腿宽尺 `grep -c '^- ⛔'` 同七枚文件 → **07=3 11=1 63=1 80=3 92=1 97=1 115=2，合计＝12**。
   ⇒ 差的 2 枚具名：`80:54`（09-29 18:2x **编排者**处置，先于 gate-rerun-1）与 `115:61`（`- ⛔ **AC#5** 不许用以下方式变绿…`＝禁令句不是指向句）。
   ⇒ **结论：腿报的 10 是对的，宽尺多算 2 枚异源行。** 本腿没有删掉任何一把，因为"指向句"这个词本身没有唯一形状。

2. **gate-rerun-1 台件枚数**：`git ls-tree -r --name-only HEAD -- .scratch/wisp/probes/gate-rerun-1 | wc -l` → **58**；
   本腿 §2 起手写的是"63 枚"（当时我加了 `grep` 过滤、数的是另一个集合）。⇒ 交件那一笔 `6d22dd6c` 的 `--name-only` 命中 **57 枚文件**（其中 `issues/` 8 枚、`probes/gate-rerun-1/` 49 枚）。
   **正确读数＝HEAD 上 58 枚件、该 commit 落 49 枚件**；"63"是本腿自己的尺用错了口径，就地具名更正，不回滚 §2 原文。

3. **抽取尺的第一版漏掉 `4c` 全部 8 枚判定行**：第一版正则 `^\| *[0-9]+ *\|` 对 `| **15** |` 这种**粗体票号**零命中，
   于是 `4c/batch4c.md` 第一遍抽成"名册 8／判 **0**"；第二版 `^\| *[0-9*]+ *\|` 抽成"名册 8／判 8"，与它 §5 自证「仍成立 8」一致。
   ⇒ 同一个文件的**两把读数**：`roster=8 judged=0`（第一版）／`roster=8 judged=8`（第二版）。
   ⚠ 这条按 §5 坑③自查：`4c` 名册表行首是普通 `| 15 |`、判定表行首是 `| **15** |`——**同一份文件里两种行首形状**，任何"只认一种"的尺都会把它的判数读成 0 或把行数读成 16。

4. **`4a` 那句"269／249 不在本尺输出里（已改名或已非零勾）"与本腿复跑冲突**：
   本腿 15:1x 逐字复跑 `4a` §2-0 那条 for 循环（`2*` 通配、`未勾>0 ∧ 已勾=0`、排除 `-done`）→ 原始输出 **30 行，含 `249 :: 未勾7` 与 `269 :: 未勾5`**。
   ⇒ 两把读数并列：`4a`＝12 行（按 `2[0-3][0-9]` 收窄到 200–238 后，且它当时读到的 249／269 不在输出里）／本腿＝30 行（`2*` 通配全量，含 15–29 段）。
   ⛔ 本腿不判谁错（两种可能都成立：249／269 在 `4a` 取数时刻还没落盘、或它的收窄段把它筛掉了——`249`/`269` 都落在 `2[0-3][0-9]` 之外，**第二个解释由本腿尺的收窄式子直接支持**），只登记"4a 那句对 269／249 的因果解释与本腿读数不兼容，它写'已改名或已非零勾'，本腿读到的是**未改名且仍零勾**"。

5. **票 234 票面引用行号 1 枚漂**：`105 :54` → 今树 `:54` 是上一格续行（`本票只交付库内 + 审计那半；**不许为了勾框扩界**。`），AC#5 在 `:56`。
   其余 11 枚引用逐字对位（`92:68`／`92:73`／`104:52`／`104:57`／`110:39`／`110:47`／`113:54`／`113:56`／`115:58`／`115:64`／`97:55`）。

6. **票 233 票面行号 2 处漂**：`restartTierKeys` 写 `:299` → 今树 `cmd/wisp/config_reload.go:338`；`schema.go` 三枚常量今读 `:36/:39/:43`（未漂）。
   ⛔ 本腿不动票面行号，只登记。

7. **判不动（本腿一票都没法判的形状）**：
   - 票 234 的 AC#2／AC#3／AC#4 与票 233 的 AC#4／AC#5／AC#6、票 237 的 AC#1／AC#3 → 全是 **go 运行期读数／突变**；本腿零 go ⇒ 只能〔量不到〕，⛔ 不许由 grep 命中外推。
   - 票 234 票面 `:35` 的"票 92 `:79`（本腿判不了那枚）"：整格要对着**禁读地界**再 grep ⇒ 派单明令不派给任何 Go 侧腿。本腿照登记，不猜。
   - **票 234 的界面半边＝`frontend/**`／`design/**`，本腿零读零转述 ⇒ 落〔量不到〕，这是正确归口，本腿没绕。**

8. **`2/` 的"46 vs 105"**：派单口径（六批＝目录数）与本腿口径（九文件＝名册数）不是同一件事，见 §3.0。
   附带一枚时刻差：`2/summary.md` 读"开放票 163 枚"，本腿今天同一把尺读 **162 枚**；差的 1 枚＝`231-a-config-written-by-a-newer-build-is-reported-to-the-operator-as-a-validation-failure-done.md`，
   它由 **`ec84cff1`（10-06 10:44，晚于 `2/` 起手 HEAD `949a5b92` 09:03）**改名成 `-done`。⇒ 与 `2/summary` 自己那句"268 因时刻不同"同一形状，不是量错。

9. **派单两条款在本腿身上冲突，具名登记（照 `4a` §前言同一处理，不擅自改规矩）**：
   「每做完一节立刻 commit」＋「骨架期不许写待填，要写'尚未算'这种实话」⇒ 骨架那 4 行占位句**必须**被真实内容替换；
   而「commit 前 `git diff --numstat` 删除列必须为 0」在这一刻字面不可达（本腿这一笔＝**210 增／4 删**，4 行删除逐行是本腿自己 §2／§3／§4／§5 那四行"尚未算／尚未写"占位，⛔ 不含任何他人文字、不含任何票面与台账）。
   ⇒ 本腿选择"交件必须可判"这一边，并在此具名；不回滚、不改写已提交的历史（更正一律追加新 commit）。

---

## 6. 本腿边界自证（收尾）

- **零 go 命令**：全程只 `git log`／`git --no-pager grep … HEAD`／`git show HEAD:<path>`／`git ls-tree`／`grep -c`／`wc`／`awk`／`sed -n`／Read／Edit。⛔ 没有一次 `go test|build|vet|run|list|env`。
- **工作树内容面**：`internal/**` 与 `cmd/**` 一律 `HEAD:` 尺（起手锚里那枚未跟踪件 `internal/tools/tasklist_deferred_236r3_teeth_test.go` 本腿**未读一字**，只数过它在 `git status` 里的存在）。
- **禁读前缀**：`.scratch/wisp/probes/236/**`、`/111/r5/**`、`/232/**`、`pool-validity/4d/**` ⇒ 零读、零引用。
- **`frontend/**`／`design/**`／`.gitignore`**：零读零写零转述（票 234 界面半边因此〔量不到〕）。
- **`docs/evidence/s1/**`**：只按**文件名**取（`git ls-tree --name-only` → 见到 `255-tier-registry-r1.md`／`255-tier-registry-v1.md`、`115-*`＝0 枚），⛔ 未读任何内容面。
- **票面零改动**：未翻任何勾、未改 `Status:`、未改名、未追加 Progress log；`docs/**` 与台账一字未碰。
- **唯一写入路径**＝`.scratch/wisp/probes/pool-validity/4e/**`（本腿新建：`batch4e.md`／`msg-skeleton.txt`／`msg-g2.txt`／`rowscan.awk`／`rowscan2.awk`／后续 msg 文件）。
- **只 commit 不 push；显式 pathspec 逐枚点名；无 `--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`／`restore`。**

