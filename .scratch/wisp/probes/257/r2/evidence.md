# 票 257 落地腿 `257-r2` 证据件 — 形 ⓒ 的首启回执（写面＝`cmd/wisp/firstrun.go` 文案＋同包测试）

## 0. 锚与写面

起手三读数（同一次发取，逐字）：

```
$ date
Mon Oct  5 11:12:54 CST 2026                      # = 2026-10-05 11:12 +0800

$ git log -1 --format=%h
c6cf66e6

$ git status --porcelain -- cmd internal
 M cmd/wisp/firstrun.go
```

分支 `dev`。写面只有这两枚（枚枚点名，无第三枚）：

- `cmd/wisp/firstrun.go` — 首启回执文案（产码结构一字未动，见 §2）
- `cmd/wisp/firstrun_257_test.go` — 本腿新测试（§3）

### 0.1 进场时的既有脏件＝我自己前一小时那半段（按 README 规则 3 复认过才继续）

`git status` 报的 ` M cmd/wisp/firstrun.go` 不是别人的活，是本票 `257-r2` 前半段留下的未提交件，
同目录里两枚日志为证（时刻与文件内容都对得上）：

- `.scratch/wisp/probes/257/r2/preflight-ticket198.log`（11:08，改文案**前**跑票 198 全家 11 枚，全 PASS）
- `.scratch/wisp/probes/257/r2/after-text-ticket198.log`（11:12，改文案**后**同一把尺，全 PASS）

`git diff --stat` 逐字 = `cmd/wisp/firstrun.go | 61 ++++++++++++++++++++++++++++++++++++++++++++++++++++`，
纯新增（零删除行），落在 `ensureFirstRunConfig` 末段 `return true, nil` 之前。
证据件本身当时还没建 ⇒ 本节是这一腿的第一份书面件。前段已 commit 的 `internal/config` 半边
（`settings_257_test.go`、`settings.go` 的三枚 tag 与三句 guidance）属 `257-r1c`，账 `A606`/`A607`，
**本腿一个字没碰那两个文件**（§4 逐字抄的是它们的现成句，不是本腿新写的）。

### 0.2 中途 HEAD 被别人推走（共享工作树，登记不归我）

11:15:57 后台重跑基线时 `git log -1 --format=%h` 已变成 `69c1bb82`（起手是 `c6cf66e6`）——
同一工作树里其它腿在我之后提交。**本腿的 commit 一律以"显式 pathspec 只含我这两枚文件"自证边界**（§8）。

### 0.3 一次真撞上的仪器坑（不是我的红，具名登记）

第一次全包基线（11:15:57，后台，未注入 dll）逐字读数：

```
BASELINE START 2026-10-05 11:15:57 +0800 anchor=69c1bb82 dirty=firstrun.go
exit status 0xc0000135
FAIL	github.com/CarlosShao/wisp/cmd/wisp	0.033s
```

`=== RUN` 计数 = **0 条**（进程在加载期就死了，一条用例都没跑到）。根因＝票 98 那一格在册老坑：
`cmd/wisp` 的测试 exe 链接 `third_party/sherpa-onnx` 的三枚 dll
（`onnxruntime.dll` / `sherpa-onnx-c-api.dll` / `sherpa-onnx-cxx-api.dll`，现量都在仓里），
不注入 PATH 就 `STATUS_DLL_NOT_FOUND`。出处逐字：
`docs/evidence/s1/101-adversarial-acceptance.md:129` 与 `:168`（`R-101-6`）、
`docs/evidence/s1/102-adversarial-acceptance.md:148`（同形读数 `=== RUN 0 条`）。
⇒ 本腿**所有** `cmd/wisp` 读数都带这条前提：`export PATH="$PWD/third_party/sherpa-onnx:$PATH"`。
第一发死掉的日志留在 `.scratch/wisp/probes/257/r2/baseline-cmd-wisp.log` 的 `11:15:57` 那一节，
注入后的重跑（`11:1x` 起）覆盖进同一文件（只建不删，故这里点名两发先后同物）。

---

## 1. 撞钉预检：动手前逐枚读到的断言原文

⚠ 最硬的一枚雷是 1.1：它禁的是**被生成的文件**里出现 `[llm.providers`，而我的指引句里必然带这四个字——
所以那句只许走 stderr 这条"给人看的回执"通道。1.1～1.4 逐枚读过才动 1.5。

### 1.1 `cmd/wisp/firstrun_198_test.go:215-219`（票 198 AC#1 第十八节的钉，逐字）

```go
	for _, forbidden := range []string{"[llm.providers", "[models.local_override", "[plugins."} {
		if strings.Contains(text, "\n"+forbidden) {
			t.Errorf("the created file carries the dynamic table %q anyway - that would be an invented entry, not a default", forbidden)
		}
	}
```

同一枚用例（`TestTicket198AC1CreatedFileHoldsEveryStaticSectionAndNoInventedTables`）还钉着
18 枚静态节各出现**恰好 1 次**（`:203-206`）与 `[plugins]` 恰好 1 次（`:211-213`）
⇒ 我若把指引写进生成的文件，红的不止一枚。

### 1.2 `cmd/wisp/firstrun_198_test.go:243-248`（回执不复用热加载四句归因，逐字）

```go
	for _, borrowed := range []string{
		"cause=missing",
		"config.toml 读不到：文件不存在",
		"本次运行继续用内存里的旧配置",
	} {
		if strings.Contains(stderrText, borrowed) {
```

⇒ 我的三段新句一个字都没用这三枚串（现量：`grep` 三枚 needle 在 `firstrun.go` 新增段里零命中）。

### 1.3 `internal/config/defaults.go:77-78`（票面"现量 1"的那两行，逐字）

```go
		case reflect.Map:
			// leave nil (see NewDefaults)
```

⇒ 首份默认配置里服务商注册表是 nil，这是本票根因；本腿**不改它**（ⓐ 那形已被毙）。

### 1.4 另两枚钉与那条坑（逐字）

- `internal/config/settings_248_test.go:130-132`（AC#7 表里"未声明的服务商"那一支）：

```go
		{"a provider the config does not declare", func(m *Manager) ([]string, error) {
			return m.SetProviderBaseURL("inventco", "https://x.example.invalid")
		}},
```

  ⇒ 写侧"行不存在则拒"是**被钉着的行为**，本腿保持原样（票 §8 边界②）。
- `internal/config/settings_257_test.go:125`（`257-r1c` 已交的干净机前提钉）：
  `if strings.Contains(body, "[llm.providers") { t.Fatalf("the first default config invented a provider row, which contradicts the ticket 198 pin:...") }`
  ⇒ 与本腿 1.1 同一条线的另一侧。
- `internal/config/parse.go:72` 与 `:123` 逐字都是 `dec.DisallowUnknownFields()` / `fdec.DisallowUnknownFields()`
  ⇒ 顶格写 `[providers.x]` 会让整份文件加载不过（票 §8 记的那处坑），指引里**绝不出现**那个拼法。
- `internal/config/catalog.go:86-110`（`validateRoles`）：provider/model 只设一枚 ⇒ `set both or neither` 拒；
  点名了不存在的 provider ⇒ `references unknown provider` 拒。这条决定了指引第三样的形状（两枚一起点名）。

### 1.5 `defaults.go`／首份文件的现量（本腿自己 dump，不是抄来的）

`.scratch/wisp/probes/257/r2/dumpcfg.go`（`//go:build ignore` 探针，不入构建）＋
`dumpcfg.log`：`config.SaveFile(path, config.NewDefaults())` 产出 **2571 字节**、33 枚节头，
其中**有** `[llm.roles.chat]` 且其内是 `provider = ''` / `model = ''`（空值），
**没有任何** `[llm.providers...]` 节。⇒ 票 §8 那句"就地填已有的那一节，别再追加一节同名的"是真读数量出来的。

---

## 2. 改了什么（逐枚 before-after）

写面产码只有一枚文件、只有**文案**，零结构改动：

| 枚 | before（锚 `c6cf66e6` 已提交态） | after（工作树） |
|---|---|---|
| `cmd/wisp/firstrun.go` | `ensureFirstRunConfig` 末段是票 198 的三句回执（新建确认／key 入口／模型入口＋票 261 的 `enabled` 一句），然后 `return true, nil` | 同三句一字未动；在 `return true, nil` **之前**追加三段 `fmt.Fprint(stderr, ...)`（AC#1 三样指引／AC#2 三因各一句／AC#0 边界①＋A560 的三形状生效时机），+38 行文案、+23 行解释注释，`git diff --stat` = `61 ++++`、删除行数 0 |

结构面的"零改动"是可核的三条，不是口头承诺：

1. 函数签名 `func ensureFirstRunConfig(dataDir string, stderr io.Writer) (bool, error)` 未变（`:72`）。
2. 判定/建文件的三段（`os.Stat` 支、`secret.NewStore`、`config.SaveFile`）逐字未变（`:73-87`）⇒
   票 198 的 AST 钉 `TestTicket198FirstRunCallerIsTheRunEntryOnly`（`:256-281`，要求调用集恰为 `[runTextTask]`）照旧绿。
3. 写面**没有任何**新增 import／新增函数／新增类型（`grep` 新段里只有 `fmt.Fprint*`）。

三段文案的逐字内容见 §4（拒因那三段）与 §3（指引那一段），本腿不在此重复。

---

## 3. 干净机那一发怎么真跑的（AC#1）

新件 `cmd/wisp/firstrun_257_test.go`，四枚用例，全部**从生产入口起步**，⛔ 没有一处手工塞 config.toml：

起步态（逐条断言，不是注释）：`t.TempDir()` 做数据根 → `os.Stat(config.toml)` 必须
`fs.ErrNotExist` → 用 `run198(t, dir)`（票 198 现成仪器，真调 `runTextTask`，即 `wisp run` 那条 CLI 缝，
README「Hard global constraints」列的合法注入缝之一）跑**一次**，退码必须仍是 `2`（⛔ 不放宽），
文件由生产自己写出。

- `TestTicket257R2AC1ReceiptNamesTheThreeHandAddedThings`
  回执里必须逐字出现三样的**具体小节名**：`[llm.providers.<名>]`／`[llm.providers.<名>.models.<模型 id>]`／
  `[llm.roles.chat]`，加上 `enabled = true`；⛔ 顶格 `[providers.` 那一形一次都不许出现（1.4 那条坑）；
  ⛔ 不许把没接线的引用解析通道列成入口（账 A560）：`resolveRefs`／`Resolved()` 两个标识符零命中；
  第三样必须是"就地填已有那一节"的措辞（含 `duplicate table` 那句反例）。
- `TestTicket257R2AC1CleanMachineRefusesAllSevenFields`
  拿生产建出的那份文件建 `config.NewManager(path, nil)`（`res=nil`＝产线三处的形状，A560），
  逐枚试写名册那七枚（`provider_base_url`／`provider_api_key_ref`／`model_context_window`／`model_price_in`／
  `model_price_out`／`role_chat_model`／`provider_credential` 的配置侧半枚，走 `panel_config_store.go` 那道
  switch 实际路由的 `Manager.Set*` 腿）⇒ **七枚全拒**，每枚的拒句必须命中"三因"里恰好一枚且不能是"文件没建"，
  七枚之后文件字节不变。
- `TestTicket257R2AC1FollowingTheReceiptUnlocksAllSevenFields` ★本腿真正兑现 ⓒ 的那一发：
  **指引串从回执里正则抠出来**（不是测试里另抄一份），把 `<名>`／`<模型 id>` 两枚占位换成实名后
  拼成手加补丁；`[llm.roles.chat]` 那一半是**在已有节内就地替换** `provider = ''` / `model = ''`；
  写回后必须 `NewManager` 加载得过（这条就是"拼法真走得通"的证明，撞 1.4 的 `DisallowUnknownFields` 会当场红），
  再逐枚试写那七枚 ⇒ **7/7 接受**，且接受后的键真的落在盘上。
- `TestTicket257R2AC1SecondRunSaysTheGuidanceOnlyOnce`
  同一目录再跑一次 `run198` ⇒ 文件字节不变、退码仍 2、那三样指引**不再重复出现**（首启专属）。

AC#2／AC#3／两通道三形状各自的用例见 §4 与 §5 的突变名册（同一件文件里，`TestTicket257R2AC2...`、
`TestTicket257R2AC3...`、`TestTicket257R2AC0...`）。

仪器前提：每条 `cmd/wisp` 命令都带 §0.3 那条 dll 注入，否则 0 条 RUN。

---

## 4. 三句拒写原因逐字（AC#2）

### 4.1 盘上现场那三枚 tag（`internal/config/settings.go:78-80`，`257-r1c` 已提交，本腿未碰）

```go
	refusalFileMissing = "第 1 种拒因：文件没建"
	refusalRowMissing  = "第 2 种拒因：行不存在"
	refusalInvalid     = "第 3 种拒因：校验不过"
```

它们在拒句里的长相（同文件逐字，本腿只读）：`:293` 是校验不过
（`refusalInvalid+"：值本身过不了这份 schema 的校验，这一行要改的是值；文件一个字节都没动。"`）、
`:367`／`:379` 是行不存在配 `guidanceModelRow`／`guidanceProviderRow`（`:87-92` 那三句 handAddGuidance 就是 ⓒ 的界面对话），
`internal/config/loader.go:77` 是文件没建。三枚 tag 互斥这一点由 `257-r1c` 在包内测过（账 `A606`）。

### 4.2 本腿写进回执的那三句（`cmd/wisp/firstrun.go:162-171`，逐字）

```
wisp run: 写不进去的时候有三种原因，各是一句不同的话，不会合成一句「配置未生效」：
  第 1 种拒因：文件没建——这一种刚才那一发已经替你办完，<cfgPath 绝对路径> 现在是真的文件；首启之前没有任何旧配置可言。
  第 2 种拒因：行不存在——那一页改不了服务商与模型的存在性，去这份文件里手加上面那三样，加完才写得进；这一条说的不是你的值不对。
  第 3 种拒因：校验不过——行在，但那个值过不了这份 schema 的校验（引用形没写对前缀、点名的模型不在目录里，都算这一种）；要改的是值，文件一个字节都没动。
```

三句的区别不是修辞而是**动作不同**，且各自主语指向唯一：第 1 种＝已由这一发办完（首启场景里它不可能再发生）；
第 2 种＝缺行 ⇒ 去手加那三样；第 3 种＝值不对 ⇒ 改值、盘没动。
"不会合成一句"这半句是**引用**那枚被禁的折迭句（票 223 AC#4 定式），⛔ 它自己不是那句折迭话——
牙怎么咬这一点在 §5 的 M2 里被证明：把三句真折成一句时，从生产拒句反推 tag 的那枚用例当场红。

### 4.3 测试怎么把"三句不同"钉住而不自造第二名册

⛔ 本腿测试**不硬编码**那三枚 tag 串（那会变成 schema 之外的第二真相源）。做法是：
在干净机上真取一发生产拒句（`SetProviderBaseURL` 打未声明的 provider），用
`第 [123] 种拒因：` 这个**形状**正则从拒句里把 tag 抠出来，再断言回执里同一枚 tag 出现且只出现一次、
三枚 tag 两两不同、每枚 tag 后面跟着的解释子句也两两不同。
⇒ `internal/config` 若哪天改名，这枚测试跟着生产改；⇒ 若有人把三句折成一句，回执里就凑不齐三枚不同 tag，必红。

---

## 5. 突变名册（每枚判据一发会响的牙）

名册先落这里，红句逐枚抄进 §5.1 之后的每一小节；改回后逐枚 `md5`／`git diff` 归零的证明同处登记。

| 号 | 把被钉的那一支改成什么 | 应红的用例 |
|---|---|---|
| M1 | 指引第一样的小节名换成顶格 `[providers.<名>]`（票面那个错拼法） | AC#1 名串枚＋"手加真走得通"枚（后者还额外撞 `DisallowUnknownFields`） |
| M2 | 三因三句折成一句「写不进去就是配置未生效」 | AC#2 两枚（三句各一句／从生产拒句反推 tag 那枚） |
| M3 | 生效时机折成一句「改完重启就好」 | AC#0 两通道三形状枚 |
| M4 | 把指引从 stderr 改写进生成的 `config.toml` | 票 198 的 `:215` 那枚钉（证明最硬的雷真会响） |
| M5 | 第三样改成"再追加一节 `[llm.roles.chat]`"（duplicate table 那个谎） | AC#1 名串枚＋手加走得通枚 |
| M6 | 凭据那句里加一条明文样例（`sk-` 开头那种值） | AC#3 凭据枚 |

---

## 7. 判不动／量不到（具名归口，本腿一律没自己划掉）

| 格 | 判语 | 归口 |
|---|---|---|
| AC#1 的"设置页那七枚"里 `provider_credential` 的**真凭据腿**（`StoreCredential`→DPAPI） | 本腿够不着：AC#3 明令凭据面一字不动，且该腿在 `cmd/wisp/panel_config_store.go`（不在本腿写面）。本腿只对它的**配置侧半枚**（那条引用落盘）负责 | 编排者：需要一次真机 DPAPI 腿的裁决表行，`docs/evidence/s1/` |
| 面板侧那三句拒因**在界面上的长相**（信封键、回执串） | ⛔ 写面禁 `internal/panel/**`（本票 §8 边界＋票 257 票面），且新增 C17 面字段属契约面（ⓑ 那形已被毙） | 编排者：若 owner 要"界面上也这三句"，需单开票（票 §8 边界③已预告 ⓑ 的两件契约面另立） |
| 全包红名册里 `TestAC1ResidentLegInstallsItsLogListenerOnDisk`（票 127 `:497` 不 poll 竞态）／`TestAC14GoSideEvalPushReachesThePage`（票 33 真建窗） | 带载偶发、早已在册，不是本腿造成的红，本腿零动作 | 编排者（原有归属：票 127／票 33） |
| `gofumpt -l cmd/wisp` 预存报 `cmd/wisp/models.go`（CRLF，不在本腿写面） | 明令不许顺手格式化 ⇒ 本腿不动，逐字读数留在 §6 | 编排者（属全仓 gofumpt 那一 sweep，票 70 家族） |
| `frontend/**` 与 `design/**` | 连读都没读（README 规则 7 ＋本票 AC#4）；`git status` 起手就报着别人的一批 `D design/**`，本腿零动作、零引用 | 编排者：那批 `design/` 删除不属本票，请确认是谁的活 |
