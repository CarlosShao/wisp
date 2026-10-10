# 300-v3 — `00` 起手锚（非实现者验收腿，只裁票 300 `AC#3` 判语归位那一半＋攻 `AC#0`..`AC#2` 恒真性）

派单射程＝**判语**：⛔ 修产码、⛔ 翻框、⛔ 改票面/台账/HANDOVER、⛔ 动任何 `*_test.go` 的断言。
写面＝只在 `.scratch/wisp/probes/300/v3/**` **新建** `.md`/`.txt`（⛔ 新建 `.go`、⛔ `.out`、⛔ 删）。

## 1. 起手现量（本程第一段 `date`，⛔ 手打）

| 尺（逐字） | 读数 |
|---|---|
| `date` | `Sat Oct 10 16:56:52 CST 2026` |
| `git log -1 --format='%H'` | `dce0f133b5123967e527ffbbcf9f3fb3ed4ce14c` |
| `git log -1 --format='%ad'` | `Sat Oct 10 16:51:56 2026 +0800` |
| `git log -1 --format='%s'` | `16:4x 收三枚交件（301-v3／298-r1／302-a1）＝A819＋票301 翻 AC#5 并追加 AC#4b（AC#4 具名永久留空）＋★★票302 前提就地改写（那三枚里两枚基线是绿的）＋★立票303（页面→Go 回执断了＝一枚本机可复现的真回归）` |
| `git status --porcelain \| wc -l` | **819** 行（他腿在飞，⛔ 我射程） |
| `git status --porcelain -- cmd internal scripts tools \| wc -l` | **0 行** ←★我这把的门禁前置数（产码面此刻干净） |
| `tasklist //FI "IMAGENAME eq wisp.exe" \| grep -c "wisp.exe"` | **0** |
| `tasklist //FI "IMAGENAME eq balldebug.exe" \| grep -c "balldebug.exe"` | **0** |

派单给的锚 vs 现量：⛔ 派单没给号（只给了"这一程＝`300-v3`"），上面那枚＝我起手现取 ⇒ 本程所有导出树**按显式号**建，⛔ 一枚按"当时 HEAD"（`300-v2` 那把尺并回）。

## 2. 与本程冲突的在飞事实（具名，⛔ 我躲）

- 另一枚腿 `303-a1` 正在反复跑 `./cmd/wisp/` 真窗测试 ⇒ 我⛔ 跑 `cmd/wisp` 整包、⛔ 起任何窗。
  我要跑 Go 只**定点 `-run`**，优先 `./internal/audio/`。
- ⇒ **整包级读数本程⛔ 取**（会被染色成假红/假绿）；需要就具名标〔欠，等编队空〕。

## 3. 我这把要用的尺清单（内容锚 `grep -n`，⛔ 信行号）

| 编号 | 尺（逐字） | 射程 | 用于 |
|---|---|---|---|
| R1 | `git grep -n "parseWaveFormat" -- internal cmd` | 产码＋测试全名册 | 必答③ 生产侧调用点枚数 |
| R2 | `git grep -n "GetMixFormat" -- internal cmd` | 同上 | 必答③ 那枚边的另一头 |
| R3 | `git grep -n "convertPacket" -- internal cmd` | 同上 | 297 ②/③ 那一支的下游 |
| R4 | `sed -n` 直读 `internal/audio/parse_wave_format_300_windows_test.go` | 入库用例本体 | 必答① 凭据＋必答② 恒真 |
| R5 | 仓外导出树定向突变（`git archive <显式号>` → `cp` 备份 → `perl` 只打一行 → 跑 → `cp` 写回 → `cmp` → `grep -c` 回到 1） | ⛔ 仓内 | 必答② 攻恒真（换形红不红） |
| R6 | `git show HEAD:<path>` 对象层对拉 | 对象层 | 证我读的＝盘上跟踪的那版 |
| R7 | `sh scripts/d22scan.sh` / `gofmt -l` / `go vet ./internal/audio/` | 全仓 | 卫生（跑前已数 scoped porcelain＝0 行） |
| R8 | `head`/`awk` 抽 `portable-tests.sh` 三档清单 | 只读 | 必答④ CI 可见性那一格在票面上的措辞 |

## 4. 我要答的四格（派单原文的射程，⛔ 越界）

1. `AC#3` 判语：盘上有没有凭据支持票 297 那一支（票面写作"③'切错字节'"）成立 ⇒ 成立／部分成立／⛔ 成立，三段各带**凭据件路径＋尺逐字＋红/绿句逐字**。
2. 攻恒真：`AC#1`／`AC#2` 的判据把**判据本体换成反形**之后还会不会红（⛔ 光看"改前必红"）。
3. `parseWaveFormat`／`GetMixFormat` 那条边**生产侧调用点枚数**现跑一把尺。
4. windows-tagged × `core`(ubuntu) 认领 ⇒ 新用例在 CI 从未执行过一次（票 301/`A817` 已量）⇒ 我只裁**票 300 票面有没有哪句话因此要说成 ⛔ 持续防护**，⛔ 另开普查格。

rc=0
