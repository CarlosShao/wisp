# 50-r1e 实现腿读数 — 票 50 AC6（plugin tool 进 list_tools + 执行记 plugin id）

锚点：起手 HEAD `e5dfd52fd`（编排者 A1443 派单 commit）；交付 commit `56eabe7ae`。
写面：`internal/plugin/**`（新增子包）＋ `internal/tools/**`（新增测试＋桥接 glue）。
⛔ 未碰：AC1/AC2/AC4 已交件（manifest.go / validator.go / tamper.go / hostapi.go 一字未改）、
`internal/tools/web*.go`、`document_text.go`、`search_content.go`、`search_files.go`、
`cmd/wisp`、`internal/observe`（45-r1 在写）、thresholds / golden。
起手核：`git status --porcelain -- internal/` 只报 `internal/observe`（45-r1 的活），
`internal/plugin` 与 `internal/tools` 干净 ⇒ 无写面冲突。⛔ push。

## 1. 起手三核里的关键发现（不是猜测，是编译器给的）

**`internal/plugin` 不能 import `internal/tools`：import cycle。**

```
package github.com/CarlosShao/wisp/internal/plugin
	imports github.com/CarlosShao/wisp/internal/tools from toolprovider.go
	imports github.com/CarlosShao/wisp/internal/agent from bridge.go
	imports github.com/CarlosShao/wisp/internal/memory from journal.go
	imports github.com/CarlosShao/wisp/internal/plugin from retention.go: import cycle not allowed
```

（这一段是 `go build` 原文，落盘前先记：`memory/retention.go` 已经依赖 `internal/plugin`。）

⇒ 第一版把 slot 实现写在 `internal/plugin/toolprovider.go` 是**编译不过的**，已删除（两个文件，
未进过任何 commit），改落到叶子包 `internal/plugin/manifestslot/`：它同时 import 两边，
`internal/plugin` 与 `internal/tools` 谁也不知道对方，环消失。写面仍在 `internal/plugin/**` 内。

**接缝形状是现成的，⛔ 新造接口**：`tools.Provider`（Kind/List/Lookup）、`tools.Entry`、
`tools.KindManifest`、`tools.Decl`（其注释逐字写着 host-side 事实放 Decl、
"read from the manifest by the Tier-1/Tier-2 providers (tickets 50/51)"）。
`registry.go` 里 `providerSlots` 本来就写着 `manifest=ticket 50` 占位——落地这张占位是本票的活，
不是改契约文字（C4 四槽位、C1 名字形状、C3 十一枚能力一字未动）。

## 2. 改动清单（8 枚，全在授权写面内）

| 文件 | 做了什么 |
|---|---|
| `internal/plugin/manifestslot/manifestslot.go` | 新增：C4 Manifest slot over 一枚已装 plugin 的声明式工具（`Run` 执行缝＋`New`/`Kind`/`List`/`Lookup`＋`entryFor` 读 C2/C3/C14 声明） |
| `internal/plugin/manifestslot/manifestslot_50r1e_test.go` | 新增 5 枚定向测试（可见性／未注册不可见／无主拒绝／越界声明拒绝／名字形状归 Registry） |
| `internal/tools/registry.go` | `providerSlots` 摘掉 manifest（goja/mcp 仍 `ErrSlotNotLanded`，`UnlandedSlots()` 3→2）；条目缺省 Provider 改取注册方 slot；manifest 条目缺 `PluginID` ⇒ 注册即拒绝 |
| `internal/tools/tool.go` | `Decl.PluginID`（host-side 边车，⛔ 给 Tool 加方法） |
| `internal/tools/gate.go` | `Decision.PluginID`（与 `Mode` 同理由：审计要能从记录回答"哪个插件跑的"） |
| `internal/tools/bridge.go` | `checkCaps` 带上 PluginID；`run()` 执行派发处写 `tools: PLUGIN-SOURCE task= corr= tool= plugin_id= kind=` |
| `internal/tools/plugin_source_50r1e_test.go` | 新增 3 枚：执行记 plugin id（含两枚反向负例）／无主条目注册被拒／goja·mcp 仍占位＋未知 slot 名被拒 |
| `internal/tools/bridge_test.go` | 原有 `TestC4SlotsWithoutAnImplementationAreNotSilent` 的 manifest 占位断言随之翻面（3→2） |

新日志行是**追加**的一条独立记录，`tools: call` 那行的 token 顺序一字未动（173/377 的形状过滤器按连续
token 匹配，插token 会伤及无辜读者）。builtin 条目不打印这一行 ⇒ 这一行本身就是"这发是插件来源"的正证据。

## 3. 读数（全部本腿现跑）

```
go build ./...                     OK（无输出）
go vet ./internal/tools ./internal/plugin/...  OK
go test ./internal/tools ./internal/plugin/... -count=1
  ok  github.com/CarlosShao/wisp/internal/tools              18.148s
  ok  github.com/CarlosShao/wisp/internal/plugin              0.137s
  ok  github.com/CarlosShao/wisp/internal/plugin/manifestslot 0.041s
定向 -v（-run 'TestAC6|TestC4Slots'）
  --- PASS: TestC4SlotsWithoutAnImplementationAreNotSilent
  --- PASS: TestAC6BridgeExecutionRecordsPluginID
  --- PASS: TestAC6ManifestEntryWithoutAnOwnerIsRefused
  --- PASS: TestAC6ManifestIsNoLongerAnUnlandedSlot
  ok  internal/tools 0.040s
gofumpt -l（改动 8 枚）＝空；仓库级 -l 只剩 .scratch 探针件与两枚他腿存量脏件
  （internal/agent/approval/pending_read.go、internal/agent/tools.go）⇒ 本腿未碰、非本腿造
d22scan -root .  EXIT=0  clean - no D22 ban violations
```

**反向判据（AC6 的两条"不该发生"必须真的红）**

1. *未注册的 plugin tool 不可见* ⇒ `TestAC6UnregisteredPluginToolIsInvisible`：只 `New`（造出条目）不
   `RegisterProvider`，`Registry.List()` 必须为 0、`Lookup("lark.search")` 必须 miss；并且"构造条目本身
   不等于可见"也单独钉了一发（防日后有人把注册偷偷塞进构造函数）。
2. *执行不记 plugin id* ⇒ `TestAC6BridgeExecutionRecordsPluginID` 三面对照：插件调用必须有且仅有一行
   `tools: PLUGIN-SOURCE`、且含 `plugin_id=com.example.lark`/`tool=lark.search`/`kind=manifest`/`task=task-1`，
   `Decision.PluginID` 同步带上；**builtin 调用必须一行都不出**（否则这字段就成噪声了）；
   **被 C3 硬拒的插件调用也不许出现这一行**（没执行过就不许留下执行归属）。
   注册侧同向负例：`TestAC6ManifestEntryWithoutAnOwnerIsRefused`（缺 PluginID 的 manifest 条目注册被拒，
   且没漏进目录）。
3. 声明越界不许被洗白：能力不在 C3 十一枚内／risk 拼错（"L9"）／无 parameters schema ⇒ slot 构建即报错，
   且错误里点名那枚 token（⛔ 静默取最像安全的默认）。
4. 占位没被顺手清空：`KindGoja`/`KindMCP` 仍 `errors.Is(err, ErrSlotNotLanded)`、`UnlandedSlots()==2`、
   未知 slot 名（"sidestep"）仍被拒。

**还原核**：本腿只动 `internal/plugin`（新子包）与 `internal/tools` 8 枚；commit 后
`git status --porcelain internal/tools internal/plugin` 空 ⇒ 无残留脏、无未提交件。
`internal/plugin` 五枚已交件（manifest/validator/tamper/hostapi + 各自测试）在 diff 里零出现。

## 4. 仍待登记（⛔ 本腿不吞）

- **AC3 仍开放**：install 流程 e2e（local path / git URL 下载→schema 校验→host_api 核查→列全能力＋来源→
  显式授权 UI→卸载 teardown）。本腿只交 AC6 的注册＋可见＋执行归属，没有装、没有 UI、没有卸载。
- **组合根接线待后续腿（⛔ cmd/wisp 不在本腿写面）**：`manifestslot.New(...)` 目前没有任何生产调用点，
  `wisp run` 起的进程里 list_tools 还看不见插件工具。要落的是"装好的 manifest → slot → `Registry`"这一段
  （cmd/wisp 或后续 AC3 腿）。形状已经给全，接线是体力活，不是未定义。
- **名字语法接缝要裁（本腿发现，⛔ 自行改校验器）**：C14 侧 `toolNamePair` 允许 ≥2 段
  （`a.b.c`），C1 侧 `validToolName` 只允许恰好一枚点。⇒ 一枚把工具写成 `com.example.lark.search`
  的合法 manifest 会在注册时被 Registry 拒掉（拒得对，但 AC3 的 UI 必须把这一步报成人话）。
  本腿只把现状钉成测试（`TestAC6ToolNameShapeIsTheRegistrysRuleNotTheShims`），⛔ 动 validator 一字。
- **Tier-1 执行体本身**：`Run` 是缝不是实现。声明式工具真跑什么属票 52（D46 命令包装），goja 属票 51。
- 缺口审计／失败预演／对抗验收按 §4.3 由**非实现者**出表，本文件只交实现腿读数。
