# 票 215 — 装配根接线（run.go 的 Plugins 读口）· 裁决表

- 交付腿：215-assembly-leg（写面 `cmd/wisp/run.go` 一个读口块＋新测试件＋本件＋票 215 Progress log 一行）
- 时刻：2026-11-08 09:5x +08 · 只 commit 不 push · AC 框与 Status 零触碰
- 欠账出处：`docs/evidence/s1/215-skills-plugins-inventory.md` §⑥.2（「生产装配根一行未接…归下一腿」）＋票 215 Progress log 23:1x 行；载体/泵半已落 `5341ead7a`＋`40074c6c1`

## ① 起手三核（现跑读数，非转述）

| 核 | 命令 | 读数 |
|---|---|---|
| 票面行还在不在 | `sed -n '1049p;1059p;1071p' cmd/wisp/run.go`（起手时）＋`git show --stat` | 三枚同款逐字在位：`BallStyle: func() panel.BallStyleSection`（1049）／`Run: rt.runFeed.Section`（1059）／`ModelSwitch: func() panel.ModelSwitchSection`（1071）；上一腿两笔真在 HEAD 链（5341ead7a＝载体 131 行＋三牙 242 行＋pump.go 39 行＋composer 14 行；40074c6c1＝漂移补录） |
| 引文出处 | `grep -n "Plugins" internal/panel/pump.go` / `plugins_inventory_215.go` | 读口签名 `PumpSources.Plugins func() PluginInventorySection`（pump.go:265）＋装配块（pump.go:481-504，拷贝后按 ID 排序、nil→`[]`）；两句 canonical 常量在载体件 70-80 行 |
| 红/绿归属（谁在飞） | `git status --porcelain` ＋ `--porcelain cmd/wisp` | 起手时 cmd/wisp 只有 323 腿两枚 untracked 测试件在盘，run.go/approval_*.go 零改动；本腿中途 219 腿落盘 `approval_always_219ac2_test.go`（??）＋后来 `approval_reply.go`（M +9/-1）——两度 raw vet 红的归属取证见 ⑤，均非本面 |

## ② 关键判断（铁律 ⓪⓪⓪，本腿自拍）

- **读者从哪来（现跑取证，⛔ 猜）**：`grep -n "PluginsSection" internal/config/*.go`＋Manager 全方法名单 → **internal/config 没有 plugins 专用 Options 族 reader**（Options 族只有 `BallStyleOptions` settings_65.go:86 / `ChatSwitchOptions` settings_187.go:166）。pump.go:252-264 字段文档说死了范式：**internal/panel 不 import internal/config，section 由装配根 pre-shape**（与 Commands/ModelSwitch/BallStyle 同款）。⇒ 读口落装配根：`mgr.Config()`（Manager 锁内深拷，manager.go:127-131，settings 腿同一张写脸）把 `[plugins]` 节逐字段映射进载体；两枚句子引载体常量（`panel.PluginEngineNotWired215`／`panel.SkillsNotInConfigModel215`），零 forked prose。
- **并发面**：pump reader 运行期反复调 → `mgr.Config()` 锁内深拷（reload tick 不竞读）；读口每包新建 sec，泵侧再拷贝排序（pump.go:493-499）——无共享可变、泵不写回读者。
- **装载面（测试喂真 Manager 的依据）**：`[plugins]` 键能进真 Manager——unwired 守卫表 5 行不含 plugins（disposition＝`not built: plugin engine is ticket 50/51`，unwired.go:137-145）；`buildPlugins`（parse.go:93）在 decode 回调里手工重装子表；`internal/plugin/tier2_gate_test.go:388-405` 就是「带 [plugins] 的文件喂真 Manager」的活例，本腿测试照此喂。

## ③ 落盘清单

| 件 | 内容 |
|---|---|
| `cmd/wisp/run.go` | PumpSources 字面量加 `Plugins: func() panel.PluginInventorySection`（接线句在 run.go:1104；净 +37 行零删除：注释块＋读口，落点在 ModelSwitch 块之后、Out 之前，照三枚同款） |
| `cmd/wisp/panel_plugins_215_test.go`（新） | 照 187 三牙形状：字面钉＋真 Manager 接线读数（硬编码期望，`plugins215Base`＝187 基线＋[plugins] 节两枚条目）＋突变（无读者 ⇒ wire bytes 零键） |
| 本件＋票 215 Progress log 一行 | 交付登记 |

## ④ 牙与突变读数

| 发 | 刀 | 读数 |
|---|---|---|
| 接线正控 | `Test215WiredPluginsReaderSendsTheLiveValues`：plugins215Base 喂真 `config.NewManager` → 生产泵 marshal，期望硬编码 | 键在；rows 按泵排序恰 `[clip zed]`；行字段逐字 `reflect.DeepEqual`（zed 的缺席授权表以 config 的 nil 原样上线）；`tier2Enabled=true` 逐字；两句 canonical 逐字与 enabled=true 同包（AC#2 正控语义延续） |
| M1（真突变） | `sed` 摘 run.go 接线块（盘上真删）→ 跑字面钉 → 还原 | **真红**：「run.go's PumpSources literal no longer carries the plugins reader…(ticket 215a)」；还原后 `cmp` 字节一致（md5 32bd4af7cb57d5b6 复位、CR=0）→ 复绿 ok 0.085s |
| 作废一发（如实记） | 先试 `-overlay` 把 run.go 换成 pristine 副本 | **没红＝作废**：字面钉读的是盘上 run.go（`os.ReadFile`），不读 overlay 编译视图——假绿不记数，重打为盘上 sed 法（M1） |
| M2（缺席半，常驻钉） | `NewSnapshotPump(PumpSources{})` | wire bytes 零 `"plugins"`（`Test215MutationReaderRemovedSendsNoPluginsKeyAtAll`） |

## ⑤ 门禁读数

| 门 | 读数 |
|---|---|
| gofumpt -l（点名 2 件，绝对路径 `/d/work/base/gopath/bin/gofumpt.exe`） | 空，rc0 |
| `go vet ./cmd/wisp` | **提交时 raw rc0**。★诚实记：本腿中途 raw vet 两度红（`undefined: sync` → 复验换 `undefined: repoRoot219`），归属＝219 腿在飞件 `approval_always_219ac2_test.go` 正在被实时编辑（mtime 00:01:22 晚于本腿 23:41 的基线 vet rc0；porcelain 里 `??` 未提交、非本腿写面、零触碰）——照 62-v1 先例用 `go test -overlay` 把该件隔离成 package stub 再跑（隔离件存 `.scratch/wisp/probes/215/r1/`，工作树零触碰），overlay 下 vet rc0 |
| `go build ./cmd/wisp` | rc0（raw 与 overlay 双读数一致） |
| `go test ./cmd/wisp -run 'Test215|TestPlugins' -count=1 -timeout 30m` | **3/3 PASS 0.087s** |
| `go test ./internal/panel -count=1` | 改前基线＝改后＝**恰 1 红** `TestC21DesignTokensFourWayAgree`（C21/ball 族既存：GlassRing/GlassHi/GlassHiSoft 不在 ball tokens.go，与本腿零字节关系） |
| d22scan | **clean**（963 internal＋241 cmd＋117 frontend＋39 design，七禁零违规） |
| LF | 2 件改动 CR=0（`tr -dc "\r" | wc -c` 读数） |

## ⑥ 诚实边界（具名，不许被读成别的）

1. **界面点不动**照票面前置：本腿交付的是 wire 上的键与它的牙；网页→Go 入向今天仍只有命令行路（账 `A410`）。`plugins` 键的 TS 声明已由 65 界面腿落（`a3fdb50dd`，panel.ts PluginsSection），渲染归界面侧既有程序（Q-51）。
2. **生效位没有**：`EngineState` 句说「今天零读者」；引擎（50/51）落地腿必须同步改 unwired.go disposition 行与这句句子——诚实句子不得活得比它描述的缺席更久（载体注释已具名此义务）。
3. **交付 2（开关）与 215b（技能）未做**：登记于上一腿证据件 ②，本腿零触碰。
4. **他腿写面全程零触碰**：219 的 `approval_reply.go`（M +9/-1）与 `approval_always_219ac2_test.go`（??）、323 的两枚 untracked 测试件、`frontend/**`、`internal/panel/**`、`internal/ball/**`、`internal/config/**`、`design/**` 本腿零字节；共享包互斥照规矩先看 porcelain 再落笔。
5. **在飞件造成的红只隔离、不修复**：219 半成品的编译错由他们自己落盘解决（提交时 raw vet 已自愈 rc0）；本腿只用自己的 overlay 读数证明自己的面。
