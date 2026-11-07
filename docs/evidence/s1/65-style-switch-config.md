# 票65 · 外观设置项切换球风格（默认玻璃球）——接线取证与裁定

> 机主原话（逐字）：「两个都要，但是设置项里面外观可以配置切换风格，默认就是玻璃球，切换Jarvis才换球的风格」。
> 本件＝65-style-switch-config 腿的取证（步1）＋裁定（步2）＋落地读数（步3-5，后半追加）。
> 时间：2026-11-07T19:0x~19:2x +08，现跑取证。

## 1. 步1 取证（全部现跑，判据可复跑）

### 1.1 ApplyStyle 的签名与它做什么

`internal/ball/orb_windows.go:653`：

```go
func (b *Ball) ApplyStyle(name string) bool {
	if !SetStyle(name) {
		return false
	}
	b.sta.PostTask(func() {
		if o := b.orb.Load(); o != nil {
			o.setStyle(ActiveStyle())
		}
	})
	return true
}
```

- 它做两件事：①进程级 `SetStyle`（`internal/ball/tokens.go:391`，D2D 渲染器重画的家族开关，词表 `styleGlassName|"glass"` / `styleJarvisName|"jarvis"`，`StyleNames()` 在 `tokens.go:382`）；②把 orb 页面导航交给 STA 线程（`orbWindow.setStyle`，`orb_windows.go:457`，内部走 `orbPageForStyle(style)` 纯导航换页）。任何 goroutine 可调；名字不在词表返回 false 且什么都没动。
- 词表枚举（`tokens.go:374-377`）：`StyleGlass`(iota 0,默认) / `StyleJarvis`(1)。`activeStyle` 包级默认＝`StyleGlass`（`tokens.go:380`）——**球默认就是玻璃球，包级已成立**。
- `orbWindow` 创建时种子取自进程级：`orb_windows.go:280` `style: ActiveStyle()` ⇒ **New 之前 SetStyle，orb 第一帧就画对风格（无闪玻璃球）**；New 之后 ApplyStyle 也能换（setStyle 纯导航）。
- **产品调用方＝0（核实）**：`grep -rn ApplyStyle --include=*.go internal/ cmd/` ⇒ 仅测试 `orb_jarvis_live_windows_test.go:83/:86`。缺的一环确证。

### 1.2 config 键现状：键在，setter/reader 不在

- `[ball].style` 键**早已在**（票65实现腿 W65-impl，commit `30d8f2e38`）：
  - 词表：`internal/config/schema.go:97-113`（`BallStyleGlass="glass"` / `BallStyleJarvis="jarvis"`，`ballStyles` 列表，注释明说「internal/config must not import internal/ball」+ 双向漂移钉）；
  - 结构：`schema.go:126` `Ball BallSection`；校验：`validate.go:70` 枚举拒绝（报错含 `ball.style` 与 `glass|jarvis`）；
  - 默认 glass：`ball_style_65_test.go`（`TestBallStyleDefaultIsGlass`）钉死——无 style 行的老文件 load 即 glass。
- **没有的不在**：`grep -rniE 'orb.*style|appearance|外观' internal/config/ cmd/wisp/ --include=*.go` ⇒ **零命中**。`SetBallStyle` / 任何 reader ⇒ 零命中。即：config 有键、有校验、有默认，但**没有写入腿、没有读出面**。

### 1.3 面板今天怎么表达配置（213 的 reader-gated 模式）

- 读侧：`internal/panel/composer.go` `Snapshot` 的指针+`omitempty` 字段族（`Instructions`(200) / `Tasks`(197) / `Commands`(213) / `ModelSwitch`(187，在飞)）；`pump.go` `PumpSources` 挂 `func() XxxSection` reader，**nil reader＝包里不带这个键**（四键 byte 钉 `pump_test.go:123/:291` 保形）；reader 由 `cmd/wisp/run.go` 装配根注入。composer.go 头注释逐字：*「A key without that reader is the constant this file was moved here to stop」*。
- 写侧：面板门 `internal/panel/config_handlers.go:448` `knownWritableField`＝三张 roster 的并集（248 表 + 313 + 187，各 roster 独立文件防冻结钉连锁红：`config_write_187.go` 头注释）；门放行后 `panel.SettingWrite` → `cmd/wisp/panel_config_store.go:199` `ApplySetting` switch → `config.Manager.SetXxx`（一律走 `writeOneKey` 守卫写：预读、validate、单键、收据报键路径）；生效档（`EffectiveTier`）按实测标（187 先例： roster 带每字段档位，门在 leg 答复后纠正 cmd/wisp 的笼统 restart 戳）。
- 构造点：`newConfigStore`（`panel_config_store.go:81`）两处调用：`panel_inbound.go:314` 与 `run.go:553`。

### 1.4 启动装配现状与 balldebug 先例

- 球装配：`cmd/wisp/resident_windows.go:406` `startResidentBall(rt.Registry, …, hotCfg258, hotReload258, hooks…)`。**球启动点 Manager 尚不存在**（258-a1 census：「task pipeline that owns rt.mgr does not exist until startResidentTaskSource below」），所以 `hotCfg258`（resident_windows.go:372）是**按次重读 config.toml**（`config.LoadFile(filepath.Join(rt.Layout.DataDir, configFileName))`）——这就是本腿启动读取该走的同一条缝。
- balldebug 先例：`cmd/balldebug/main.go:97/:130-131`——`-style` flag 在 `ball.New` **之前** `ball.SetStyle(*style)`。运行时切面（新球第一帧即对风格）正是这个次序。
- `cmd/wisp/run.go` 与 `internal/panel/{composer,pump,config_handlers}.go` 本波有 188/187 腿未提交 hunks 在盘（`git status --porcelain` 实见）——这些是共享脏区。

## 2. 步2 裁定：选甲（config 键＋启动读取＋面板设置项＋运行时即时生效），按 187 冲突拆两段落

**选甲**。三条候选与代价：

- **甲（选它）**：①config setter `SetBallStyle`＋reader（键/校验/默认已在，补写入腿与读出面）；②`cmd/wisp` 启动按 `hotCfg258` 同款按次重读 `[ball].style`，`ball.New` 前 `SetStyle`（无闪）＋装配后 `ApplyStyle` 幂等复同步；③面板设置项按 187 同款 roster＋snapshot 键（reader-gated）；④写入路径 `ApplySetting` 成功后调 `ApplyStyle` 即时生效（applier 注入缝挂球装配处）。
- **乙（不选）**：只做 config＋启动读取，切换要重启。代价＝机主原话「切换Jarvis才换球的风格」被降格成「改文件→重启才换」；`ApplyStyle` 整套机器（SetStyle＋STA 换页导航，测试已验）继续零产品调用＝**已建成的能力被闲置**；主流 harness 的外观项无一不是即时生效。省下的只是「一个注入缝＋一个 case」，收益为负。
- **丙（不选）**：只做面板项不落 config。代价＝重启丢失，直接违反 D36「config.toml 是权威」的三档模型（回票 248 的 EffectiveTier 语汇：一个永远不会 `EffectiveNow` 也不落盘的开关）；且快照键若不接 config reader 就是 composer.go 头注释点名的「constant this file was moved here to stop」。

★**默认＝玻璃球三重成立**（机主逐字「默认就是玻璃球」）：①`[ball].style` 缺省 load 即 glass（validate/schema 默认，`ball_style_65_test.go` 已钉）；②ball 包级 `activeStyle = StyleGlass`（tokens.go:380）；③本腿新件**零玻璃球字面量**——默认从 config schema 默认与 ball 包默认自然流出，不新造一份。

### 冲突拆段（协调纪律）

187 腿未提交 hunks 正压 `internal/panel/composer.go` / `pump.go` / `config_handlers.go`（`git status --porcelain` 实见；`run.go` 被 188 腿压）。`git add` 按文件点名会吞同文件里他们的 hunks ⇒ **本腿⛔碰这四个文件**：

- **本腿落**：`internal/config/settings_65.go`（setter+reader，新件）、`internal/panel/config_write_65.go`（roster+字段名+档位语汇，新件，**不接 knownWritableField 并集**）、`cmd/wisp/ball_style_65_windows.go`（种子读取+applier 注入缝，新件）、`cmd/wisp/resident_windows.go`（干净件，两处 hook 调用）、`cmd/wisp/panel_config_store.go`（干净件，`ApplySetting` 加 `ball_style` case＋档位按 applier 结果如实标）、测试件、证据件、票65 log。
- **留给下一腿**（187/188 落盘后各一小行/一小段）：`composer.go` 加 `OrbStyle *OrbStyleSection` 快照键、`pump.go` 加 reader 字段+fill、`config_handlers.go` 并集加 `writableField65(field)`、`run.go` 注入 reader。**如实报**：本腿落盘后「设置项在页面上出现＋写路径被面板门放行」仍差这四小段；config 键、启动读取、`ApplyStyle` 产品调用三段本腿闭环。

### 撤销口令

`git revert` 本腿 commits 即整体退回；单段退回＝删新增四件＋还原 resident_windows.go / panel_config_store.go 两处 hunk（`git log --oneline -- cmd/wisp/ball_style_65_windows.go` 找回滚点）。

## 3. 步3 落地（commit `bc48ae6c7`）

| 件 | 形状 |
|---|---|
| `internal/config/settings_65.go`（新） | `SetBallStyle`（`writeOneKey` 守卫写，键路径 `ball.style`，枚举在 config 侧把关＝320 裁定同形；`validate()` 照跑）＋`BallStyleOptions` reader（`BallStyleView{Style, Vocabulary}`，读 `Manager.Config()` 活内存，词表＝schema 的 `ballStyles` 原列表）。**零玻璃球字面量**。 |
| `internal/panel/config_write_65.go`（新） | `FieldBallStyle = "ball_style"`＋本腿 roster（`WritableFields65`/`writableField65`），照 187 roster 形状独立成件防冻结钉连锁红。**未接 `knownWritableField` 并集**（187 冲突区，见 §4）。 |
| `cmd/wisp/ball_style_65_windows.go`（新） | `ballStyleSeed65`（照 `hotCfg258` 按次 `config.LoadFile`，不可读答 `("", false)`＋Warn 说 fallback，无本地默认）＋`ballStyleApplier65` 注入缝＋`installBallStyleApplier65`（挂 `Ball.ApplyStyle`）＋`applyBallStyle65`。 |
| `cmd/wisp/resident_windows.go`（hunk，`runResident` 内） | `startResidentBall` **前**：seed＋`ball.SetStyle(name)`（orbWindow:280 创建种子 `ActiveStyle` ⇒ jarvis 配置第一帧即粒子球，无闪；balldebug `-style` 同序）；`startResidentBall` **后**：`installBallStyleApplier65(rb.b)`＋`rb.b.ApplyStyle(name)` 幂等复同步＝**`ApplyStyle` 产品调用方从 0 变实**（两处：启动复同步＋设置写入路径）。 |
| `cmd/wisp/panel_config_store.go`（hunk，`ApplySetting`） | `case panel.FieldBallStyle: written, err = s.mgr.SetBallStyle(w.Value)`＋档位实测：applier 在＝`EffectiveNow`（真换过了），不在＝`EffectiveRestart`（上面盖的默认档）——187 的静态 `fieldTier187` 标法不适用，因为这一档取决于本进程有没有球，不是读者行为固定的字段。 |

⓪⓪⓪ 复核：选甲的代价已写清（§2）；`ApplyStyle` 机器是已建成能力，接线是买真价值；未把任何技术选择题丢给机主——默认玻璃球是他逐字令，其余形状（守卫写、roster 独立、reader-gated、档位实测）全部照仓里既有主流形状（248/313/187/320/258）接。

## 4. 如实边界（⚠ 未整段完成的部分）

**面板快照键与门准入＝留给下一腿**，四小段，187/188 落盘后各一处小改：

1. `internal/panel/composer.go`：`Snapshot` 加 `OrbStyle *OrbStyleSection`（指针+omitempty，同 ModelSwitch 形状）；
2. `internal/panel/pump.go`：`PumpSources` 加 reader 字段＋`Snapshot()` fill；
3. `internal/panel/config_handlers.go`：`knownWritableField` 并集加 `|| writableField65(field)`（`takesValueFor` 同理加 `fields65TakesValue`）；
4. `cmd/wisp/run.go`：装配根把 `BallStyleOptions`（config 侧 reader 已就绪）注进 PumpSources。

**为什么留着**：这四个文件是 187 腿未提交 hunks 的载体（`git status --porcelain` 实见 composer/pump/config_handlers 脏；`git add` 按文件点名会吞他们的 hunks，派单纪律明令禁）。**因此今天的状态是**：面板设置页上还看不到外观选择项（门会拒 `ball_style`，快照无键）；但 `wisp run` 启动已按 `[ball].style` 换球风格、手改 config.toml 的 `style = "jarvis"` 即时生效于下次启动、写入路径（setter→applier→`ApplyStyle`）已全线打通并有测试钉——**收件三把尺里「ApplyStyle 产品调用者＝0」这条缺陷已消除，验证口径见 §5 的两处调用点**。

## 5. 步4+5 判据读数（全部现跑，2026-11-07T19:3x +08）

四牙（`cmd/wisp/ball_style_65_windows_test.go`，改名为 `Test65*` 对齐派单门禁 `-run 'Test65|TestStyle'`）：

| 牙 | 用例 | 读数 |
|---|---|---|
| ①默认值是 glass | `Test65SeedDefaultGlass` | 无 style 行的合法 config ⇒ seed 返回 `(glass, true)`——glass 来自 schema 默认，本腿新件零字面量 |
| ②设 jarvis 真变 | `Test65SeedJarvis`＋`Test65ApplySettingBallStyleFlipsActiveStyle` | 写路径（真 Manager 守卫写→applier）后 `ball.ActiveStyle()==StyleJarvis` |
| ③设回 glass 真变 | 同上 | 第二次写 glass ⇒ `ActiveStyle()==StyleGlass`，收据/档位逐枚断言 |
| ④突变 | 摘 `ballStyleSeed65` 的 config 读取（临时 return `"", false`） | `Test65SeedDefaultGlass` **FAIL**＋`Test65SeedJarvis` **FAIL**（`Unreadable` 仍 PASS＝失败路径语义独立）；还原后全绿，突变用例逐字复跑 |

config 侧（`internal/config/settings_65_test.go`）：守卫写收据 `[ball.style]`／盘上 canonical 单引号行 `style = 'jarvis'`／词表外四值拒绝且 config.toml 字节不变／信封空格修剪而大小写折叠拒绝（load 枚举是精确小写）／nil Manager 读口词表完整。

门禁读数表：

| 门 | 命令 | 读数 |
|---|---|---|
| gofumpt | `/d/work/base/gopath/bin/gofumpt.exe -l`（七件点名） | 空（rc0） |
| vet | `go vet ./internal/config ./internal/ball ./internal/panel ./cmd/wisp` | rc0 |
| build | 同四包 `go build` | rc0 |
| config | `go test ./internal/config -count=1` | ok 2.197s |
| ball | `go test ./internal/ball -count=1` | ok 0.242s |
| panel | `go test ./internal/panel -count=1` | rc1，**3 红逐枚归属他腿**：两枚预定红（`TestApprovalCardViewJSONKeysMatchFrontendTypes`/`TestComposerContractTypesMatchFrontend`，红行＝`Go Snapshot emits [commands modelSwitch]`——213 的 commands＋187 在飞的 modelSwitch，**无本腿键**）；`TestC21DesignTokensFourWayAgree`＝c21-native-tokens.md 与 `internal/ball/tokens.go` DarkPalette 缺 GlassRing/GlassHi/GlassHiSoft——ball tokens 面，本腿零触碰 |
| cmd/wisp | `go test ./cmd/wisp -run 'Test65\|TestStyle' -count=1 -timeout 30m` | ok 0.226s，4/4 PASS |
| winlive | `go test -tags winlive ./internal/ball/ -run 'TestBallLiveOrb' -count=1 -timeout 25m` | **7/7 PASS 50.797s**（含 `TestBallLiveOrbJarvisStyleRendersParticles`＝ApplyStyle 往返真窗验证、四角 (51,51,50) 精确零差未破） |
| d22scan | `./tools/d22scan/d22scan.exe`（全量） | clean rc0（internal 954 Go＋cmd 233 含注释与测试） |

**LF**：七个触件 `file` 逐一核过，无 CRLF。

**Bug 榜（写腿自报，两处）**：①config 断言初版用了 `style = "jarvis"` 双引号形状，canonical 写盘是单引号——测试红了才看读写盘真形状（守卫写本身无缺陷），已按真形状改钉；②初版测试名 `TestBallStyleSeed65…` 不含 `Test65` 子串，派单门禁 `-run 'Test65|TestStyle'` 选不中——改名对齐门禁后再跑。

**AC 勾与 Status 未动；⛔ push（推送只由编排者做）。**

## 6. 追加更正（2026-11-07T19:4x，append-only，§4 原文保留不改）

本腿最后 commit（`635922c76`）落盘前后，**并发收口腿**已把 §4 留的四小段落掉了大半：

- `fe0f740d6`（收口腿步1）：现跑核实本腿全部交付（含「门 384 拒 ball_style」的机制指认），并发现并集实为**三处缝**（`acceptedWritableFields`/`takesValueFor` 也在 `config_write_313.go`，187 先例三处全接）；
- `1b25d931c`（收口腿步2a）：`knownWritableField` 并集接 `writableField65`（成四表并集）、两处附加并集缝同款接 65、`composer.go` 快照新键 `BallStyle`（reader-gated 指针+omitempty）、`pump.go` reader+fill、新件 `ball_style_65.go`，并把本腿 `config_write_65.go` 的过时拆段注释如实改为已落；
- **仍差最后一段**：`cmd/wisp/run.go` 装配根把 `BallStyleOptions` 注进 PumpSources——该件此刻脏（另有腿在飞），收口腿尚未提交。

组合态复跑（本腿现跑，非转述）：`go vet` 三包 rc0；`go test ./cmd/wisp -run 'Test65|TestStyle'` ok 0.212s（4/4）；`go test ./internal/config` ok 2.446s；panel 两契约红现读数＝`Go Snapshot emits [commands modelSwitch ballStyle]`——**ballStyle 已随收口腿进入同一 Q-51 未声明红**，属既定程序（UI 腿声明 frontend 侧后转绿），非缺陷。

§4 写作时点（本腿 `bc48ae6c7`，~19:35）读数属实；盘面随收口腿推进，以本节为最新状态。

## 3. 步3-5 收口（2026-11-07T19:3x +08，接续腿现跑核实＋落地读数）

### 3.1 步1 现跑核实（五问五答，判据可复跑）

- **① `internal/panel/config_write_65.go` 全文已读（71 行，未跟踪）**：wave71 写的第四张名册件。内容＝`FieldBallStyle="ball_style"` 常量、`fields65` 名册（一项）、`fields65TakesValue`（一项 true）、`WritableFields65()`（排序 accessor）、`writableField65()`（本名册查询）。头注释自述：为什么是第四张名册（248/313/187 各被冻结钉＋互斥断言钉死，不能插行）；接线拆分如实声明；档位不在这里盖（applier 装了答 now，没球答 restart）。
- **② `config_handlers.go:448-450` `knownWritableField` 现状**：`listedWritableField(field) || additionalWritableField(field) || writableField187(field)`——**三表并集，无 65**。⇒ `ball_style` 被门在 `config_handlers.go:384` 明确拒绝（`"不在可写名册内"`），这就是机主原话在活机上不可用的直接机制。
- **③ 证据件 §2「接线拆分」原文**：留给下一腿四小段＝`composer.go` 加快照键（原文占位名 `OrbStyle *OrbStyleSection`）、`pump.go` 加 reader 字段+fill、`config_handlers.go` 并集加 `writableField65(field)`、`run.go` 注入 reader。
- **④ config 侧 `[ball].style` schema/validate/默认（现跑复核）**：词表 `schema.go:97-98/:113`（`BallStyleGlass="glass"`/`BallStyleJarvis="jarvis"`/`ballStyles`），枚举拒绝 `validate.go:70-72`，默认 glass 由 `ball_style_65_test.go`（`TestBallStyleDefaultIsGlass`＋与 ball.StyleNames() 双向漂移钉）钉死。30d8f2e38 早已落，属实。
- **⑤ cmd/wisp `ApplySetting`「already wired」真假**：**真，但在未提交盘面**。`git diff cmd/wisp/panel_config_store.go` 实见：`case panel.FieldBallStyle:` 调 `s.mgr.SetBallStyle(w.Value)`，写后按 `applyBallStyle65` 实测结果把档位从 restart 纠正为 now。配套在盘未提交件：`cmd/wisp/ball_style_65_windows.go`（种子 `ballStyleSeed65`＋进程级 applier 缝 `ballStyleApplier65`/`installBallStyleApplier65`/`applyBallStyle65`）、`resident_windows.go` 两处 hook（New 前 `ball.SetStyle` 无闪种子＋装配后 `ApplyStyle` 幂等复同步＋applier 安装）、`internal/config/settings_65.go`（守卫写 `SetBallStyle`＋读口 `BallStyleOptions`）、两枚测试件（`ball_style_65_windows_test.go` 牙①②③＋档位如实、`settings_65_test.go` 写路径四测）。`go build ./internal/panel ./internal/config ./cmd/wisp ./internal/ball` rc0——半成品自洽。
- **★步1 新发现（简报与 §2 都没列的两处并集缝）**：`config_write_313.go` 里 `acceptedWritableFields()`（:150，拒绝文案「可写的只有 …」的清单＝三表并集硬拼接）与 `takesValueFor()`（:170，空值检查＝三表并集）。187 当年三处全接；65 只接 `knownWritableField` 一处的话，拒绝文案撒谎（不报 ball_style）、空值检查漏问（空值漏到 config 侧才被枚举拒绝，报错形状错层）。⇒ 本腿按 187 同款**三处全接**（写面扩到 config_write_313.go 的这两个函数体，文件本身无在飞 hunks，`git status` 实见干净）。

### 3.2 冲突解除核实

187 腿已落盘（`007de34f0`，13 枚），`git status --porcelain internal/panel` 现跑＝仅 `config_write_65.go` 未跟踪，`composer.go`/`pump.go`/`config_handlers.go`/`config_write_313.go` 全干净 ⇒ §2 预留的四小段现在可以动，无吞 hunk 风险。`cmd/wisp` 在盘未提交的 65 件与 323 测试件（`panel_resident_323*_test.go`，他腿）并存——本腿 commit 只点名 65 件，⛔ 吞 323。

### 3.3 步2 落地读数

**并发对账（诚实记录）**：本腿步1 commit（`fe0f740d6`）落盘后，**另一条并发收口腿**把 §2「本腿落」清单整段 commit 了（`bc48ae6c7` 步3+4 ＋ `635922c76` 步5收件 ＋ `d86e591ba` 追加更正，即上文的 §4/§5/§6——故本文件节序为 §1§2§4§5§6§3，append-only 未重排）。四小段由本腿落（`1b25d931c` 步2a：knownWritableField 四表并集＋两附加缝＋composer 快照键＋pump reader＋新件 `ball_style_65.go`；`7858108a0` 步2b：run.go 注入——⚠ §6「仍差 run.go 注入」一句在写作时点属实，`7858108a0` 是 `d86e591ba` 的祖先，**终态 run.go 注入已在库**）。命名一处偏离 §2 占位：快照键落成 `BallStyle *BallStyleSection`/wire 名 `ballStyle`（对齐 `FieldBallStyle` 与 `[ball].style` 线名，§2 的 `OrbStyle` 是未定占位）。

四小段落点（`1b25d931c`，6 件，gofumpt 空/build rc0/全 LF）：

1. `config_handlers.go` `knownWritableField` ＝ `listedWritableField || additionalWritableField || writableField187 || writableField65`（四表并集，248/313/187 钉子原样）；
2. `config_write_313.go` `acceptedWritableFields` 并入 `WritableFields65()`、`takesValueFor` 并入 `fields65TakesValue`（步1 新发现的两缝，187 三处全接同款）；
3. `composer.go` `BallStyle *BallStyleSection`（指针+omitempty，`Q-51` 未声明红同既定程序）＋新件 `ball_style_65.go` 载体（Style/Vocabulary 两字段，零词表字面量）；
4. `pump.go` `PumpSources.BallStyle` reader＋fill（nil 词表渲染 `[]`）；`run.go`（`7858108a0`）装配根注入 `BallStyle: func() panel.BallStyleSection`（读 `rt.settings` 同一个 `mgr` 的 `BallStyleOptions`）。

### 3.4 步3 牙读数（panel 门级＋突变④逐字）

新件 `internal/panel/ball_style_65_test.go`（3 测）＋`cmd/wisp/panel_ball_style_65_test.go`（run.go 源钉，213 同款）：roster 四表并集且三旧名册零污染／门级 `config.set ball_style` 直达 leg（applyCalls=1、值逐字、收据含 `ball.style`）＋空值门拒（正控）＋未列名拒（负控）／ballStyle 键 reader-gated 双向＋nil 词表渲染 `[]`。全部 PASS。

**突变④（牙四条之一，逐字）**：临时把 `knownWritableField` 摘成三表（去掉 `|| writableField65(field)`）：

```
--- FAIL: Test65DoorAdmitsBallStyleAndRoutesItToTheLeg (0.00s)
    ball_style_65_test.go:77: the door refused ball_style: panel: 配置字段不在可写名册内: "ball_style" 不在可写名册内；可写的只有 ball_style/model_context_window/model_price_in/model_price_out/model_row_create/model_row_delete/provider_api_key_ref/provider_base_url/provider_billing/provider_credential/provider_row_create/provider_row_create_custom/provider_row_delete/role_chat_model/role_chat_pair/role_chat_thinking/voice_cloud_asr_chain/voice_cloud_tts_chain/voice_realtime_model/voice_realtime_provider（requestId="r-65"）
FAIL
FAIL	github.com/CarlosShao/wisp/internal/panel	0.076s
```

`go test rc=1`。注意拒绝文案里 ball_style **仍在列**——这恰好证明 `acceptedWritableFields` 缝独立起作用（文案缝与门缝是两条），摘门缝只红门级行为。还原后 sha256 逐字节一致（`2048692e72db376c2cc41407c2191ac525b7e8bd83d38262fc23ac2a22817616`，突变前后同值），门级测试复绿。牙①②③由并发腿测试承担（§5 表，本腿复跑见 §3.5）。

### 3.5 步4 终态门禁（组合态全量复跑，2026-11-07T19:5x +08）

| 门 | 命令 | 读数 |
|---|---|---|
| gofumpt | `/d/work/base/gopath/bin/gofumpt.exe -l`（本腿九件点名） | 空（rc0） |
| vet | `go vet ./internal/panel ./internal/config ./internal/ball ./cmd/wisp` | rc0 |
| build | 同四包 `go build` | rc0 |
| panel | `go test ./internal/panel -count=1` | rc1（17.124s），**3 红逐枚归属**：`TestApprovalCardViewJSONKeysMatchFrontendTypes`/`TestComposerContractTypesMatchFrontend`（红行＝`Go Snapshot emits [commands modelSwitch ballStyle] that interface PanelSnapshot does not declare`——ballStyle 已按 Q-51 既定程序入列，UI 腿声明后转绿）＋`TestC21DesignTokensFourWayAgree`（ball tokens 面，本票零触碰） |
| config | `go test ./internal/config -count=1` | ok 3.136s |
| cmd/wisp | `go test ./cmd/wisp -run 'Test65\|TestStyle\|TestComposer' -count=1 -timeout 30m` | rc0，**5/5 PASS**（Test65SeedDefaultGlass / Test65SeedJarvis / Test65SeedUnreadable / Test65ApplySettingBallStyleFlipsActiveStyle / Test65AssemblyRootCarriesTheBallStyleReader；后两 pattern 在 cmd/wisp 无匹配函数，如实记） |
| winlive | `go test -tags winlive ./internal/ball/ -run 'TestBallLiveOrb' -count=1 -timeout 25m` | **7/7 PASS rc0 56.819s**（含 `TestBallLiveOrbJarvisStyleRendersParticles`＝ApplyStyle 真窗往返，四角 (51,51,50) 判据未破） |
| d22scan | `./tools/d22scan/d22scan.exe` | clean rc0（internal 956 Go＋cmd 234，注释与测试含） |

**LF**：本腿全部触件核过，无 CRLF。

### 3.6 收口终态（机主原话兑现路径）

设置页外观项现在的完整链：**快照键 `ballStyle`**（`composer.go` ← pump reader ← `run.go` 注入 ← `config.BallStyleOptions`）→ 页面拿到当前风格＋词表（glass 在列）；**写**：`config.set ball_style` → 门 `knownWritableField` 放行 → `ApplySetting` case → `config.SetBallStyle`（枚举先拒后写，守卫路径）→ `applyBallStyle65` → `Ball.ApplyStyle`（进程级换家族＋STA 换页导航）→ **即时生效**（有球答 now；无球答 restart，boot 种子下次启动落地）。`ApplyStyle` 产品调用者＝**2**（boot 幂等复同步＋设置写入路径），不再是 0。默认＝玻璃球三重成立（config schema 默认＋ball 包默认＋新件零字面量）。

本票收口链 commits（双腿交错，全只 commit 不 push）：`110e3bdcf`（步1+2 取证裁定）→ `fe0f740d6`（本腿步1）→ `bc48ae6c7`/`635922c76`/`d86e591ba`（并发腿步3+4/步5/更正）→ `1b25d931c`（步2a）→ `7858108a0`（步2b run.go）→ `432d5ab7b`（步3）→ 本笔（步4 收口读数＋log）。
