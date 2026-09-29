# 223 — **`CheckAndReload` 在生产里一次都没被调用过**：配置改了不重读、D33 那条"放宽必须 L2 复确认"接的是一根没通电的线

- Status（**09-29 16:5x 收写腿 `223-r2`，表＝`docs/evidence/s1/223-hot-reload-wiring-r2.md`（骨架 `28475620`→代码＋读数 `5a755c3c`→收工读数若干 commit），台账待编排者落 `A##`）：⛔ 仍七格不翻勾、不翻 `-done`（写腿不自裁，D22 双角色）。r2 只做派单两件事、零扩射程，交回的是读数**：①**AC#5 那枚本票红已改确定性**——`TestTicket223RestartTierSaysItWillNotApply` 的 stdout 单发读改成 `awaitStdout` 有界轮询（monotonic timer／零断言放宽／产品文案一字未动）；**改前本腿复现 25 发 1 红**（红句 `:453` 与历史逐字同形）、**改后同一把尺 20 连发 exit 全 0**（`.scratch/wisp/probes/223/r2/flake-20.txt`）；**正控**：注释掉产品那句"本次运行不会生效"后用例 **42.18s 转红**（`stdout never carried "本次运行不会生效" within 40s`，同发 stderr 两行 restart 审计仍在＝轮询非恒真），`git cat-file` 还原、前后 SHA256 逐字相同（`542f705a…34f4670`）。②**AC#4 两句归位**——`internal/config/loader.go` 新增一支：读得出声明版本（当前版或更新）而正文解析不开 ⇒ 包成 `config.toml parse`＝**语法错那句**；读不出版本仍走 branch 1；**只有声明低于当前版才交给迁移管线**（`migrate_test.go:123` 一字未动、`internal/config` 整包 `ok 0.787s`）。新增常驻用例 `TestTicket223R2FailureSentenceRouting`（8 形双向钉：C1/E1/G1/L1/J1/A1→`cause=syntax` 且不得带「语法没问题」；声明旧版→仍 `cause=migration`；解析得开＋未知键→不抢语法错）首跑 8/8 绿；**牙**：还原成改前 loader 后恰红被接管的 5 形（红句逐字 `the two sentences swapped again`，`teeth-loader-removed.txt`）、对照 3 形纹丝不动，还原哈希逐字相同（`14251d9f…b28b0d30`）。三把门：`go build ./...`＝0／`gofumpt -l`（四枚）＝净／`d22scan` clean；整包单发 `internal/config`＋`cmd/wisp` 双 `ok`（110.376s）。**待裁点具名交回**：J1（声明未来版＋正文坏）归语法错＝r2 按派单条形 1 的现读裁，若裁另保 newer-build 句，改法见 r2 表"没做完"第 2 条；票 227 那枚 226 时序账本腿未观测到红、未动。**下一枚＝非实现者对新两格的复核表＋编排者复跑，然后才谈翻勾。**
- Status（**09-29 16:2x 收非实现者验收腿 `223-v1`，表＝`docs/evidence/s1/223-hot-reload-wiring-v1.md`（56,522 字节／265 行／6 枚提交 `4fc03f21`→`ec5d63fc`），台账 `A439`）：⛔ 七格一格不勾、不翻 `-done`**。**它的判语：AC#1/AC#3/AC#6/AC#7 成立，AC#2/AC#4 成立带注，AC#5 不成立。**⇒ **下一枚＝`223-r2`（写腿，只做两件事）**：①**AC#5 那枚本票自己的红**——`cmd/wisp/config_reload_223_test.go:431 TestTicket223RestartTierSaysItWillNotApply` 是**时序竞态**（等 stderr 审计行、却对 stdout 单发读，窗口在 `config_reload.go:284→:288` 之间），改法＝**把那枚断言改成有界轮询**（用 `context` 限期，⛔ 不许写减法、不许放宽断言）；⚠ **命中率我自己复跑过**：它 5 发 1 红、**我 15 发 1 红（第 13 发，红句在 `:453`）＝合计 2/20**，所以是**间歇红不是必红**，判据要按"连续 20 发全绿"来定，别拿一发绿交差。②**AC#4 那句会说反的话**——17 发 overlay 实测只有 1 发真走迁移管线，`C1/E1/G1/L1` 四发**读得出版本 2＝当前版**、根本不进迁移分支，落到 `cause=invalid` 那句"语法没问题"，而它们的毛病**正是语法错** ⇒ 两句归反了；补一发常驻判据（`migrate_test.go:123` 一字不许动、它今天仍绿）。⚠ **对我自己的两条更正（它列的，我复认）**：**"AC#5 由编排者代跑补齐"这个标题为过**——我那发缺 `-v`、把单发当终态，**降格为"对照读数"**，并立一条通则：**今后"代落"（收死腿现场）可以，"代跑"（替它出终态读数）不作先例**；另**跨票那枚时间炸弹**（票 226 唯一走真宿主的用例靠"跑不满 1 秒"过关）已在**票 227 新加 AC#6**。
- Status（**上一口径，保留不抹**）：**产码已交并由编排者代落（09-29 15:33，commit `248095d1`）；七格里实现腿自陈六格成立、AC#5 由我代跑到终态，⛔ 一格都不翻勾、不翻 `-done`**——裁决表与"没做完/判不了"十格见 `docs/evidence/s1/223-hot-reload-wiring-r1.md`（26,252 字节＋我补的 AC#5 一节，五枚提交 `386f7752`→`ec7a034d`→`029f4763`→`248095d1`），下一枚＝**非实现者对抗验收 `223-v1`**（它要补两件事：**带 `-v` 的终态逐名绿册**，我这发没带；以及 AC#6"裁决在锁外"那一格**缺一发专属突变**——把裁决挪回持锁段、断两枚 `Config()`-readable 用例转红）。⚠ **它交回的现场是半改状态**：五枚文件未提交、含既有产码 `internal/config/loader.go`（＋72/−6，方向＝AC#4"语法错不许被说成`迁移失败`"）⇒ 我**一字未改**、只验门（`go build` 净／`gofumpt -l` 对那五枚零输出／三包 `-count=1` 全 ok：`internal/config 0.794s`、`internal/perm 0.243s`、`cmd/wisp 108.047s`）后**代落**（先例＝`201-r1`、票 35）。**骨架先落盘那条模板升级第二次生效**：这次死在中途只损失一格读数，不是整张表。
- Status（**旧口径，保留不抹**）：**待派，且它是票 219（三枚答复按钮）与票 222 之后第一批的前置之一**——⚠ **不修它，任何"放宽 `[fs]` 要重新确认"的判据都只能在测试里绿**。
- ⚠⚠ **写腿 `223-r1` 顶回我票面两格，我都复认（原句不抹，见下面两条）**：**①AC#2 那句"三档＝立即／重启后／下次会话"里的"下次会话"今天不存在**（`internal/config/schema.go:33-44` 只有 `TierHot`／`TierReload`／`TierRestart` 三枚，`TierReload` 自陈"立即生效＋发重载事件"）；**②我派单给的基线"83 枚具名用例"口径不完整**（那把尺数的是顶层用例；含 `--- PASS` 子测试行是 **127 条**）⇒ **终态比红名册一律用"用例名集合"、不用枚数**。
- 来源：非实现者设计复核 `219-v0`（`docs/evidence/s1/219-approval-reply-design-adversarial.md`，27,833 字节）第 5 条＋编排者自己复跑。

> ⚠ **09-29 11:5x 普查腿 `223-c1`（`.scratch/wisp/probes/223/c1/census.md`，49,708 字节，提交 `172bda57`）交回、编排者逐条现跑复认（台账 `A426`）：本票原来那句"只有一枚断口"是半句——真形状是**双重断开＋一枚会说谎的注释**，五处更正如下（原表不抹）**
> ① **第二枚断口（我现跑复认）**：`Manager.ConfirmLocked`（声明 `internal/config/manager.go:29`、读取 `:240`）**生产零赋值点**——`grep -rn "ConfirmLocked = " --include=*.go cmd/ internal/ | grep -v _test` **空**（命中全在 `manager_test.go`×4 与 `internal/perm/ticket90_persist_test.go`×2），而 `:240` 逐字 `approved := m.ConfirmLocked != nil && m.ConfirmLocked(section, loosen)` ⇒ 钩子恒 nil ⇒ **就算把轮询接上，任何 `[fs]` 放宽仍会被 fail-closed 静默拒掉**。⇒ 新增 **AC#6**。
> ② ⛔ **注释与实现不一致，而它决定落点的形状（我逐行读码复认，不是推断）**：`manager.go:19-20` 逐字「Callbacks run OUTSIDE the lock so they may call Config()」——这句**对 `OnReload`／`OnRestartPending` 成立**（`:141` 解锁之后才发），**对 `ConfirmLocked` 不成立**：`CheckAndReload` `:124` 加锁 → `:138` **锁内** `apply` → `apply` 内逐段 `applyLocked("risk"/"fs"/"net"…)` → `:240` 调钩子 → `:141` 才解锁；而 `Config()` 自己 `m.mu.Lock()` ⇒ **钩子里回头读 Manager 即自死锁**，且一张 L2 卡最长 300s 会把整个 Manager 冻住。⇒ **本票的落点不许是"同步等人点卡"那形**（要么异步回报裁决、要么解锁之后再裁决）。
> ③ **门禁不拦 ticker、拦"出生方式"**（读的是 `tools/d22scan/main.go` 正则原文）：`:144` `wallclockRe`＝`\.Sub\(time\.Now\(\)\)` ⇒ **`time.Since` 与 `time.Now().Sub(prev)` 不在射程**；真正拦裸协程的是 `:697` 那条 `*ast.GoStmt`；`:665` 让 bans #1-5 **跳过 `_test.go`**、只有 ban #8 含它 ⇒ **"测试里绿"证明不了生产干净** ⇒ AC#1 必须写"生产调用者点数"，不许拿测试用例数交差。
> ④ **AC#2 的第二档连出口都没有**：`OnRestartPending` 同样**零生产赋值点**（我现跑：只有 `manager.go:33/:36/:146/:147` 与 `manager_test.go:82/:94`）⇒ "改了→不重启→**明确告知未生效**"这一档今天不是"接了没通知"，是**根本没接**。
> ⑤ **那圈 tick 的家早就登记成推迟件**：`internal/watchdog/` **只有 `doc.go`**，`:18` 挂着 `DEFERRED(watchdog loop/thresholds): implemented by ticket 42`。⇒ 落点选"给 watchdog 真建一圈"还是"在 `cmd/wisp` 里自己起"代价不同（后者零新增依赖边：`run.go` 已同时持有 config 与 gate，见普查 C/D 节）；⚠ 顺带撞出"`DEFERRED` 标记与 `SPEC-12 §5` 不 1:1、且没有仪器在查"——**那件不归本票，另立票 225**。

## 现量（起手逐条复算，别信这里的行号）

| 事实 | 读数 | 尺 |
|---|---|---|
| 那枚轮询函数存在 | `internal/config/manager.go:113` 一带 `CheckAndReload`，注释自陈是 SPEC-03 §4.3 的"幂等轮询（watchdog 1s tick）" | 现读 |
| ⛔ **非测试调用者只有一枚，且在旁支程序里** | **`cmd/balldebug/main.go:243`**：`bridge.Refresh = func() error { _, err := mgr.CheckAndReload(); return err }` | `grep -rn 'CheckAndReload()' --include='*.go' . \| grep -v _test` ⇒ 只剩这一行＋`internal/ball/hotkey_reload.go:21` 的**注释** |
| 生产装配不接它 | `cmd/wisp/run.go` 只在**装配时**建 canonicalizer（`grep -n 'Paths\|canonicaliz' cmd/wisp/run.go` 现跑），**没有任何 tick 调 `CheckAndReload`** | 现读 |
| 冻结要求 | `docs/PLAN.md:1644-1645` 逐字「**（新增）D33 配置提权**：热加载放宽 `[risk]`/`[fs]`/`[net]`/`[plugins]` → **必须触发 L2 级重新确认，不得静默生效**」 | `sed -n '1643,1646p'` |

## 后果（为什么这不是"少一个便利功能"）

1. **D36 的三档生效级别（立即／重启后／下次会话）今天没有任何一条走"立即"那档的实现路径**——没人轮询，就没有"热加载"。
2. **D33 的安全轨接不上**：想"放宽可读写范围时强制再确认一次"，接的是一根不转的线；于是"长期允许"这一支今天**要么做不到、要么做出来就是静默生效**（后者直接违反上面那句逐字要求）。
3. **owner 第一次用就会得出"这是个 stub"的结论**：改了 `config.toml` 里的东西，产品不重启不生效，而且**没有任何一句告诉他这件事**。

## 判据（每格都要现跑读数）

- [ ] **AC#1 生产里真有人在轮询**：`CheckAndReload`（或等价的显式刷新入口）在 `wisp run` 的装配路径上有**生产调用者**（现跑点数：改前非测试命中＝`cmd/balldebug` 那一枚、改后 `cmd/wisp` 里 ≥1 枚），并具名写清触发形状（watchdog tick／文件通知／显式命令）与**为什么不是"用墙钟时间差实现超时"**（`AGENTS.md` §1.2 硬禁）。
- [ ] **AC#2 三档生效级别各有读数**：立即生效的那一档要能观测到"改了→不重启→行为变了"；重启后生效的那一档要能观测到"改了→不重启→**明确告知未生效**"。⛔ **不许用"静默不生效"充当第二档**。
- [ ] **AC#3 放宽必带 L2 复确认（D33 正控）**：种一条把 `[fs]` 放宽的改动 ⇒ **必须产生一张 L2 卡**；把那次复确认拿掉 ⇒ 判据**必须红**。⚠ 收紧（变严）那一向**不许**也弹卡（否则把安全轨变成骚扰）。
- [ ] **AC#4 不生效与读不到是两句话**（同票 216 的纪律）：配置文件缺失／语法错／权限不够／热加载被禁用，**四种各一句**，不许合成一句"配置未生效"。
- [ ] **AC#5 整包终态读数**：`PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test ./cmd/wisp ./internal/...` 到终态＋逐名比红名集合（历史在册的 `internal/panel` 那几枚红**不算本票新增**，也不许顺手修）。
- [ ] **AC#6 第二枚断口必须一起接上（09-29 由 `223-c1` 顶出，见上面第 ① 格）**：`ConfirmLocked` 要有**生产赋值点**（现跑点数：改前＝0、改后 ≥1），且 AC#3 那枚"放宽必带 L2 复确认"的正控**要走这条生产路径**、不许只在测试里把钩子手动设成 `return true`。⚠ 同时**不许把裁决做成锁内同步等待**（第 ② 格那条自死锁）⇒ 判据要能证明"卡挂着的时候 `Config()` 仍可读"（种一发：另起一路调 `Config()`，若在等卡期间超时/卡住＝本格红）。
- [ ] **AC#7 "重启后生效"那一档要有出口（第 ④ 格）**：`OnRestartPending`（或等价的显式提示面）在生产里有赋值点，且用户能观测到"我改了、这次运行不生效、**为什么**"这句；⛔ 不许用"静默不生效"充当这一档，也不许把这格与 AC#2 第一档合成一枚用例。

## 禁区

不动 `PLAN.md`／`docs/specs/**`／`thresholds.go`／golden／`allowlist.txt`；不新增 C17 方法名；不新增导出名（要新增先落 `A##`）；`internal/panel/l2_grant_boundary_test.go`／`tokens_fourway_test.go`／`internal/perm/ticket90_persist_test.go` 三枚冻结件一字不动；`frontend/**`／`design/**` 零写零转述；⚠ 与票 201／222 同撞 `cmd/wisp` ⇒ **串行**。
