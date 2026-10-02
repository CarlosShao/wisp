# 票 248 · 对抗验收腿 `248-v1` 裁决表 — 设置那条写路径（Go 半边）

> **本腿代号**：`248-v1`（裁决者，≠ 实现者 `248-r1`）。⛔ 本腿不改任何产码或测试：所有变异逐枚回滚并证 SAME。
> **被测对象**：产码提交 `0d87a681`（15 枚路径）＋ 编排者代提尾部 `b644d310`。
> **票面**：`.scratch/wisp/issues/248-the-panel-has-no-settings-route-the-inbound-whitelist-has-zero-config-or-credential-methods.md`（本腿只读；AC 框一枚未碰）。
> ★ **本腿进场即处理的大洞**：实现者证据件 `docs/evidence/s1/248-settings-write-path-r1.md` 的 **§4「门禁与读数」正文是占位（`（待填）`）**，
> 而它 §7 表里 AC#6 那一行写着"四数在 §4"＝**指着一节空的凭据**。
> ⇒ 本腿一律按〔**无凭据**〕处理票 248 的门禁结论，并**自己现跑四把尺**（本件 §2），不引它 §4 一字。
> **两层禁令**：⛔ 未读未写 `frontend/**`／`design/**`；本件不出现那两层的任何行号。
> **凭据纪律**：本件不含任何凭据值；哨兵串是仓里既有的测试常量（假值），只写常量名与形状，不复述其值。

---

## 0. 起手锚（同发取，逐字）

| 尺 | 读数 |
|---|---|
| `date -Iseconds`（进场第一发） | `2026-10-02T10:07:50+08:00` |
| `git rev-parse --abbrev-ref HEAD` | `dev` |
| `git log -1 --format='%H %ad %s'` | `40a75e82ac4d8a86be7dfd35e5a28c2003cd0da9 2026-10-02T10:05:23+08:00 evidence(198-v1)：…` |
| `git status --porcelain=v1 \| wc -l` | **336 行**（非空且永远非空＝共享工作树；整份存档 `/tmp/248v1/porcelain-start.txt`，15,579 字节） |
| 起手脏度（本腿被测四路） | `git status --porcelain -- cmd/wisp internal/config internal/panel docs/evidence/s1` ⇒ **1 行**：`M docs/evidence/s1/248-settings-write-path-r1.md`？—— 见 §0a 现量 |
| 门禁四把尺 | 见 §2（本腿现跑，每发带当时 HEAD） |
| 闸门口径 | ⛔ 不是"终态为空"，而是 **终态名册 ＝ 起手名册 ＋ 只本腿那枚文件**；本腿唯一写面＝`docs/evidence/s1/248-settings-write-path-v1.md` |

**§0a 起手名册里与票 248 有关的那几枚（逐名现量）**

> 上面那张表是**前腿 `248-v1` 自己**在它那个时刻（HEAD `40a75e82`）取的起手锚，本接管腿**一个字不删、也不替它重取**；
> 它 §0 表里第 22 行那句「见 §0a 现量」当时指着一节空的正文 ⇒ **本腿把 §0a 填上**，那一行从此有凭据（凭据是本腿的，不是它的）。

| 尺（本腿 248-v1b 现跑，`2026-10-02 10:37:50 +0800`，HEAD `580d6153`） | 读数 |
|---|---|
| `git rev-parse --abbrev-ref HEAD` | `dev` |
| `git status --porcelain -- cmd/wisp internal/config internal/panel docs/evidence/s1` | **0 行** ⇒ 前腿 §0 表第 22 行那句「1 行？」**在当前 HEAD 未复认**（本腿量得 0 枚脏） |
| `git status --porcelain`（整树）行数 | **311 行**（前腿同尺报 336 行；差 25 行＝中间落了 `613606c0`/`580d6153` 等提交与他人台件，**不是同一时刻的名册**，不可相减成结论） |
| 整树名册里含 `248` 字样的路径枚数 | **8 行**（`git status --porcelain \| grep -i 248 \| wc -l`）⇒ 全是 `.scratch/wisp/probes/248/**` 与本件相关路径，**没有一枚是产码/测试路径的脏** |
| 被测四枚核心文件在 `0d87a681` ↔ `HEAD` 的 blob 全等性 | `internal/panel/bridge.go` `bebe8e702a85`＝同、`internal/panel/config_handlers.go` `19445546ea12`＝同、`internal/config/settings.go` `047e7e4c81c2`＝同、`internal/panel/l2_grant_boundary_test.go` `6668f19bc14f`＝同 ⇒ **本腿在 HEAD 上跑的读数就是 `0d87a681` 的形状**（`internal/panel`/`internal/config` 两层零漂移） |
| `0d87a681..HEAD` 里动过被测三包的提交 | **1 枚**：`613606c0`（票 198-r1，只动 `cmd/wisp/**` 与它自己的台件）⇒ `cmd/wisp` 整包读数**含票 198 的增量**，归因时逐名分栏 |

**§0b 前腿留在盘上的台件：本腿逐名读过了什么（引自前腿日志第几行＋当前 HEAD 复认与否）**

| 台件 | 里面是什么 | 本腿处置 |
|---|---|---|
| `.scratch/wisp/probes/248/r1/cmdwisp-full.txt:1` | `ok cmd/wisp 285.774s` | **本腿不复用为凭据**；已在当前 HEAD 自行重跑（见 §2） |
| `.scratch/wisp/probes/248/r1/cmdwisp-final.txt:1` | `ok cmd/wisp 302.323s` | 同上 |
| `.scratch/wisp/probes/248/r1/panel-config-final.txt:1-26` | `internal/panel` 4 枚 `--- FAIL` 名＋`internal/config ok` | 本腿在 HEAD 重跑，**逐名比对＝同名册**（§2 复认成立） |
| `.scratch/wisp/probes/248/r1/panel-head-baseline.txt:1-9` | 6 枚 `--- FAIL` 名（比终态多 `TestComposerRenderFixtureTellsTheTruth`、`TestReadGitOnThisRepositoryIsSelfConsistent` 两枚） | 本腿 HEAD 现跑**只有那 4 枚**⇒ 多出的两枚本腿未复认（取数时刻 `08:50`，早于 `0d87a681` 落库） |
| `.scratch/wisp/probes/248/r1/panel-mine.txt:1-7` | 4 枚 `--- FAIL` 名 | 与本腿现跑同名册 |
| `.scratch/wisp/probes/248/r1/start-anchors.txt` | **只有起手名册（277 行），零把门禁尺** | ⇒ **四把尺（build/vet/d22scan/gofumpt）前腿根本没取过**，本腿 §2 是唯一凭据 |
| `.scratch/wisp/probes/248/r1/{bridge.go,config_handlers.go,settings.go}.pristine` | 实现者的变异回滚副本 | 本腿 `git hash-object` 与 HEAD blob 逐枚比对 ⇒ **三枚全 SAME**（§3 末段） |

**结论（起手层）**：前腿的死法（99 次工具调用、正文全空）**没有留下任何它已取完但未写进表的门禁读数**——它的日志里只有三包整包与名册；四把尺由本腿现跑。

---

## 1. 逐格 AC 判语

## 2. 四把尺现值（带 HEAD）＋三包整包名册

## 3. 变异清单

## 4. 退回与否

## 5. 我可能判错的条目

## 6. 判不动的地方
