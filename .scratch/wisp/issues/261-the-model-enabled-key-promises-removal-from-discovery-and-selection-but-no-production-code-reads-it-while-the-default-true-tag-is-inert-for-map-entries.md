# 票 261 — 模型条目的 `enabled` 注释承诺"从发现/选择里移除"，而枚举那一处根本没读它；同时 `default:"true"` 对 map 条目在解码期不生效 ⇒ 手写配置里少写一行＝那枚模型落到 `false`，**而这个 false 今天也没人执行**

**立票时刻**：2026-10-03 10:5x +08，锚点 HEAD 见本票 commit（编排者自取，⛔ 别让下一程引用别人写的号）
**来路**：只读腿 `145-a2`（`.scratch/wisp/probes/145/a2/census.md` 186 行／26,429 字节，骨架 `bfdd910a` → 交件 `675987c7`，零产码零 Go 命令）在量"面板快照那两维的生产者侧"时实锤报出；编排者（＝我）三处现跑复认后才立本票，落账 `A566`。
**性质**：⚠ **算产品行为**，不是纯仪器——它决定"用户在设置里把一枚模型关掉，到底关没关掉"。

## 现量（编排者本机自跑，2026-10-03 10:5x，逐字）

1. **承诺在注释里**：`internal/config/schema.go:348` 起是 `type ModelSpec struct {`，末段逐字
   `// Enabled removes the model from discovery/selection without deleting its catalog entry.`
   紧跟字段 `Enabled bool ` + "`toml:\"enabled\" default:\"true\"`"。
   ⇒ 这句话是**行为承诺**（"removes … from discovery/selection"），不是"仅作记录"。
2. **枚举那一处不读它**：`internal/llm/resolver.go:287-288` 逐字
   `out := make([]string, 0, len(p.Models))` ／ `for id := range p.Models` —— **循环体里没有任何 `Enabled` 判据**；
   `internal/llm/resolver.go:119` 那处 `spec, ok := p.Models[model]` 同样只问"在不在目录里"。
   ⇒ 我把 `grep -rn "\.Enabled" --include=*.go internal/llm internal/config internal/agent cmd/wisp`（排除 `_test.go`）跑了一遍：
   命中的 `Enabled` **全是别的结构**（`internal/config/catalog.go:126` 是 `c.Voice.Realtime.Enabled`；
   `internal/config/manager.go:380/394/512` 是 `Voice`／`WakeWord`；`cmd/wisp/logsink.go:205` 是日志开关）。
   ⚠ **口径写死**：`ModelSpec.Enabled` 的**生产读者＝零枚**；上面那些命中⛔ 不许被下一程当反证。
3. **默认 tag 对 map 条目不生效**：`internal/config/defaults.go:75-80` 的 `applyDefaults` 在
   `case reflect.Map:` 分支逐字只有一句注释 `// leave nil (see NewDefaults)` —— **它不进 map**，
   所以 `[llm.providers.<名>].models.<id>]` 这类**手写条目**若没写 `enabled = true`，解码后落 **`false`**（Go 零值），
   `default:"true"` 那句 tag 对它**没有任何执法者**。
   ⚠ 与票 257 咬合：票 257 选的 ⓒ 形＝**首启引导用户手写 provider**（`cmd/wisp/firstrun.go` 那一支），
   那条路造出来的配置正是"手写条目"形状 ⇒ 本票不是"以后可能撞上"，是**现在式**。
4. **发现那一路会写这个字段吗**：`internal/llm/discover.go:107-118` 里 `p.Models[m.ID] = config.ModelSpec{...}`
   是发现器**建行**的地方（我看到 `:117` 那行字面是构造体，⚠ **是否显式写 `Enabled: true` 我没逐行读**——
   这条留给 AC#0 判死，不许按"应该会写"结案）。
   ★**10-03 10:5x 编排者就地补读（原话不改、追加具名）**：那一行没读成，我现跑 `grep -n "Enabled" internal/llm/discover.go`
   ⇒ 真身＝**`internal/llm/discover.go:119` 逐字 `Enabled: true,`**（构造体在 `:118`，行号比我票面写的 `:117` 漂一行）。
   ⇒ AC#0 的 ⓑ 已有答案＝**发现那一路显式写 `true`**，所以"缺键落 `false`"只发生在**手写条目**那一支（＝票 257 ⓒ 形造出来的形状）。
   ⚠ 但这**不改变**本票的缺口性质：写侧填了 `true`，**读侧零枚读者**才是"关掉＝没关"那件事实；`ⓑ` 那一格因此降级成"复认＋判它够不够"。

## 为什么它单独成票（不并进票 255／180／145 的理由，具名）

- 票 180 AC#4 那把尺数的是**段级字段**（"`internal/config/schema.go` 里还有几枚有 `default` 但生产零读者"），
  本缺陷是 **map 条目里的字段**——两层都不在它那把尺的射程里：① 它的粒度是"一枚字段名"，
  `ModelSpec.Enabled` 只在 `for id := range p.Models` 这种**动态键**里存在；② "default tag 在 map 上不生效"这一维它压根没问。
- 票 255 打的是**回执说假话**（"已立即生效"）；本票打的是**承诺没人执行**（关掉＝没关）。同族不同格。
- 票 145 AC#2b 要求快照里的模型清单是"目录 ∩ `enabled=true`" ⇒ **交集那一半今天无执法者**（`145-a2` 报、我 §2 复认）。
  那格我另裁（见下方"排程"），⛔ 不许把它算成"已经修好了"。

## 要建什么

- [ ] **AC#0 先把两件事判死，本票的第一格必须是读数不是码**：
  ⓐ **谁该读 `ModelSpec.Enabled`**——把"发现（`internal/llm/discover.go`）／选择（`internal/llm/resolver.go` 的枚举与 `:119` 那一问）／回执（面板清单）"三条路逐枚列出**今天谁读、谁不读**，
  并判死：**有没有任何一条真路径能让一枚 `enabled=false` 的模型被选中并产生花费**（判据＝一发跑起来的读数，⛔ 不许只按读码推断）；
  ⓑ **`discover.go:117` 那个构造体到底显不显式写 `Enabled: true`**（我上面 §4 明写我没读，这一格必须有人读它并给逐字）。
- [ ] **AC#1 修法二选一（⛔ 不许两形都不落就把本票结案；先让我落具名 `A##` 选边再派落地腿）**：
  ⓐ **让承诺兑现**＝在**唯一那一处枚举点**（`resolver.go:287-288` 那一路）按 `Enabled` 过滤，
  并给一发**会响的尺**：种一枚 `enabled=false` 的条目 ⇒ 它⛔ 不许出现在枚举结果里；补回 `true` ⇒ 出现（**两向都要实测**）；
  ⛔ **不许**把过滤写进三张 provider 适配器（`145-a2` 报：那形会把写面扩到 `internal/llm` 三枚文件＝写面扩张，要另裁）；
  ⛔ **不许**顺手把 `settings.go` 的写侧改成"写 false 就等于删行"（那会改掉"条目还在目录里"这半句承诺）。
  ⓑ **具名降级**＝承认"今天 `enabled` 不生效"，把 `schema.go` 那句承诺注释**改成诚实形**并登记 `DEFERRED`（五字段齐全，`SPEC-12 §5`），
  同时票 145 AC#2b 那句"∩ `enabled=true`"必须跟着摘掉——⚠ **这一支会碰承诺文字**，代价要摊清才许选。
- [ ] **AC#2 `default:"true"` 那一维的射程普查**：`applyDefaults`（`internal/config/defaults.go:75-80`）不进 map ⇒
  **全仓 map 型结构里还有几枚字段带 `default` tag**（尺＝逐枚列 map 字段＋带 tag 的条目枚数，给推导式与口径），
  每一枚答一句"缺键时落零值会不会被当成有意义的值"（例：`false` 会不会被读成"用户主动关了"）。⛔ 不许顺手改 `defaults.go` 的行为——那是**第二件事**，要另格。
- [ ] **AC#3 契约轴**：`docs/PLAN.md`（含 D36 的 `[llm.providers]`／`[app]` 那一节）／`docs/specs/**` 一字节不许动；
  `thresholds.go`／golden／`allowlist.txt` 不许动。⚠ **若有人主张"干脆摘掉 `enabled` 这个键"＝碰 D36 段树＝契约变更**，
  那一支**必须停下来摆 owner**（本票两形都不含它，⛔ 不许把它塞进"顺手清理"）。
- [ ] **AC#4 门禁**：逐包 `go test -count=1 ./internal/llm/ ./internal/config/`＋`go build ./...`＋`sh scripts/d22scan.sh`；
  ⛔ 不为变绿放宽任何断言、⛔ 不动 `internal/panel/tokens_fourway_test.go` 一字；
  `internal/panel` 那 2 枚 colour/token 已知常红按已知读（起因＝`design/**` 在工作树里被删未 staged，与本票无关）。

## 禁区

- `frontend/**`／`design/**` **既不读也不引**（两层禁令；界面那一侧归 owner 自己委托的会话）。
- Git 纪律：只 commit 不 push；commit 必带**显式 pathspec**；⛔ `git add -A`／`.`／`--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`；临时件只建不删；⛔ 不在仓内建 worktree/checkout。
- ⛔ **API 密钥的值**绝不出现在任何读数／证据件／日志里，只允许变量名与配置键名。

## 排程与归属（编排者裁定，写给下一程）

1. **本票排在票 145 AC#2b 之前**：那格要往快照里加"当前可选模型清单"，而清单的**过滤语义**属本票。
   快照那一维我此刻的临时口径＝**先按目录内全部模型出、不承诺 `enabled` 交集**（⛔ 不许写成"已按 enabled 过滤"那种假话）；
   撤销口令**「261 交集口径重开」**。落地腿 `145-r3` 起手若发现本票已裁 ⓐ 形，按 ⓐ 收敛。
2. **写面**＝`internal/llm`＋`internal/config` ⇒ 与在飞的 `255-r3`（`cmd/wisp`）、`252-r2`（`internal/tools`）**不互斥**，
   但**起手必跑** `git status --porcelain internal/llm internal/config` 并具名报回，撞面就停手上报。
3. ⛔ **AC 框归编排者**，产码腿一枚都不许碰；判语必须由**非实现者**验收腿给。
