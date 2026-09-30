# 票 246 AC#7 — 非实现者验收腿 `246-v2` 裁决表

- 腿：`246-v2`（对抗验收，只裁 **AC#7** 一格；不写产码、不修任何东西、**不勾任何 AC 框**）
- 起手 HEAD：`aec06d24`（`git rev-parse --short HEAD` 现跑，2026-09-30 23:43；含台账 `A486` 那枚 `45202bb0`，满足"起手＝`A486` 或更新"）
- 被验收的产码：`bb37fac2`／`3edc11d4`／`d8f3b27e`／`2071f59e`／`8af4347d`（写腿 `246-r2`，死于 150 轮上限）
- 判据原文出处：`.scratch/wisp/issues/246-*.md`「AC#7（新补）」那一格
- 上一枚验收腿：`246-v1`（裁 AC#1..AC#6，已被翻勾，不在本表射程）
- ⛔ 本表不引 `246-r2` 的自述当凭据：它表内 §5（突变与正控）与 §6/§7 未写完，那些格由本腿**自己现跑取数**。
- ⛔ 本表不含任何凭据值（密钥一律只写变量名）。

## §1 门禁读数

### §1.1 进程与前置（本腿现跑，2026-09-30 23:43 +08）

| 尺 | 命令（逐字） | 读数 |
|---|---|---|
| 球调试进程占位 | `tasklist //FI "IMAGENAME eq balldebug.exe"` | `INFO: No tasks are running which match the specified criteria.` ＝ **0 枚** |
| Wisp 进程占位 | `tasklist //FI "IMAGENAME eq wisp.exe"` | 同上一行一字不差 ＝ **0 枚** |
| 起手 HEAD | `git rev-parse --short HEAD` | `aec06d24`（含 `45202bb0 ledger(A486)` ⇒ 满足"起手＝`A486` 或更新"） |
| 工作树（本腿地界） | `git status --porcelain -- cmd internal` | **0 行**（起手即干净；脏项全在 `.gitignore`／`design/**`／`probes/**`／别人名下两枚 evidence 文件，一枚未动） |
| 真机用例前置 | `GOFLAGS= go build -o build/wisp.exe ./cmd/wisp` | rc=**0**；`build/wisp.exe` mtime `23:43:42` > `cmd/wisp/run.go` mtime `23:37:24` ⇒ 不是拿旧 exe 读数 |

### §1.2 起手基线（本腿本人这轮取，⛔ 不引编排者的数当凭据）

命令逐字：`PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" GOFLAGS= go test ./cmd/wisp -count=1`

| 项 | 读数 |
|---|---|
| 起跑 | `2026-09-30 23:44:09 +0800` |
| 终态 | `2026-09-30 23:47:06 +0800` |
| 结果 | **`ok github.com/CarlosShao/wisp/cmd/wisp 173.655s`＝全绿，`--- FAIL` 计数 0** |
| rc | 0 |

口径对照（编排者 23:37–23:41 那发，锚 `18270ad7`）＝ `ok ... 209.743s`。**两发都是"各自那一次的读数"**，不是上一版 vs 新版的关系；本腿只把自己这发当凭据。

PATH 坑本腿复认过它的形状：不带那两枚目录时 `go test` 报 `exit status 0xc0000135`、`0.0xxs`、**一行 `--- FAIL` 都没有**＝一枚用例都没跑（仪器坑，不是红）。

### §1.3 起手红名册归因（AC#5 的同一把尺，本腿重跑）

起手基线 **0 枚红** ⇒ `A482` 那枚 `TestAC228ExitRequestDuringBootStillLeavesThroughD38E` 在本腿起手这一发**未复现**。它是 boot 时序 flake，命中率不因这一发改变；**归口＝票 228（AC#8）**，本腿不顺手修、不当回归报。

（本节以下随各攻推进补满；交付时本表每一节都必须是有内容的格子，⛔ 不许留空段与占位句——本票已连续两枚腿死在最后两节，那是我最在意的失败形状）

## §2 逐格判语（AC#7）

## §3 六攻读数（A 生产调用者／B 一进程一枚门／C 四步真机读数／D 定向突变／E 重锚的钉／F 零新依赖边）

## §4 突变台账（改了哪一行 → 指名哪一枚用例红 → 还原凭据 → 复跑回到绿）

## §5 我判不动的地方

## §6 我对实现者表的更正

## §7 收尾三把尺

## §8 next
