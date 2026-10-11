
## 编排者更正（2026-10-11 13:1x，来路＝只读普查腿 `writeface-1` 的逐腿文件名册（它被授权顶回我），件 `.scratch/wisp/probes/orch/20261011/writeface-1/**`；台账 `A846`；⛔ 上面各节原句一字不改，本节追加）

本票**⛔ 翻任何格**（`AC#2` 那格按 `A795` 口径仍不翻，`AC#5` 承接它）；本节只更正**我派单里的一处前提**＋落一条排面结论。

- **★前提更正：`AC#5` 那三处过头话住在⛔ `cmd/wisp` 的测试件里，而是住在产码文件**（我派 `writeface-1` 时写的是"测试件"）。我自己现跑两把尺（HEAD `510c95ca`）：
  - 正尺＝`sed -n '36,44p' cmd/wisp/resident_windows.go` ⇒ 三处锚逐枚在**产码文件的文档块**里＝**`:38`**（`// an empty answer: Check()'s only skip is the DeepEqual diff against what is`）／**`:40`**（`// here would re-register the set and leave summon, mute and panel dead until`）／**`:42`**（`// owner's machine (an empty [hotkey] slot means "unset", not "turned off").`），它们是 `func hotkeyReloadSource296(...)`（`:44`）那段注释的一部分。
  - 反证尺＝`git show HEAD:cmd/wisp/resident_hotkey_296_windows_test.go | grep -nE "only skip|dead until|turned off|DeepEqual"` ⇒ **0 命中**（那枚测试件⛔ 转述这些句子；它 `:22` 那行只是**转述 `hotkeyReloadSource296` 的来历**，腿另记一条"改注释时顺带扫它一眼，但今天⛔ 同步义务"）。
- **连带排面结论（这是我要的那条）**：⇒ `296-r2` 的写面＝**一枚产码文件的注释面**（`cmd/wisp/resident_windows.go`）⇒ 它与**票 298 那两枚测试件**（`cmd/wisp/panel_inbound_guards_35r3_test.go`／`panel_transport_35r2_test.go`）**⛔ 同名册交集、⛔ 是同一枚文件** ⇒ 我先前"这两件事撞同一枚文件、⛔ 能同批"那句⛔ 成立。⚠ 但**排面仍按住**：`296-r2` 要取 `./cmd/wisp/` 的颜色，而此刻 `internal/tools`（`309-r1` 在写）与 `internal/risk`（`308-r1` 在写）都在 `cmd/wisp` 的导入图里 ⇒ **它排在两枚写腿交完之后**（同波⛔ 给第二枚要 `cmd/wisp` 色的腿）。
- ⚠ 顺带一枚腿报的、与本票相关的口径提醒（⛔ 我改任何件）：`.scratch/**` 下的 `.go` blob 副本会挪**别人**那把"tracked `.go` 分母"的尺（腿现量 `git ls-files '*.go'`＝**992**，它自己那批副本又添了 5 枚 ⇒ 今后引"22 枚 `.scratch` 故意脏拷贝"那类名册要写明**含⛔ 含探针目录的 blob 副本**）。★我本程自己就现场犯了一次（往自家探针目录丢了一枚⛔ 用的 `.go`，当场改名 `.go.txt`）⇒ 定式＝**探针里的 blob 副本从一开始就⛔ 叫 `.go`**（先例＝`1eaa91f6` 那笔 `git mv` 把两枚 `.go` 改成 `.txt`、内容逐字节 `cmp` rc=0）。

- [2026-10-11 13:1x] agent=编排者 did=收 `writeface-1` 对本票那条顶回⇒ **派单前提具名更正**（`AC#5` 三处锚在**产码** `cmd/wisp/resident_windows.go:38/:40/:42`＝`hotkeyReloadSource296` 文档块；反证尺在同名测试件上 4 个词 **0 命中**）＋**排面结论落定**（`296-r2` 与票 298 两枚测试件⛔ 同名册、⛔ 同一枚文件，但它要 `cmd/wisp` 色 ⇒ 仍按在 `309-r1`＋`308-r1` 之后）＋一条探针卫生定式（blob 副本从一开始就⛔ 叫 `.go`，我自己现场犯过一次）；★本票六格⛔ 动、⛔ 改原句、⛔ push
