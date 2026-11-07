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

（步2 commit 后追加）
