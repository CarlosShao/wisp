# 票 356 普查腿 356-a1 — R3 A 档下 `task.output` 的 notice 对 `fs.read` 的实际行为（只读成对读数）

- 角色：**只读腿**（`AC#0`）。未翻任何格、未改票名/`Status:`/台账、未落任何 `.go` 字节进仓、未 push。
- 起手 `HEAD=a8f5ad0a`；跑门禁时 `HEAD=c0c8aaa9`；落件之后本腿每一笔都记在第 7 节（**不在此处写死终态 HEAD**，那会随每一笔自我过期）。期间编排者插了 `c8da185f`（只动 `docs/reports/**`）。`git diff --stat a8f5ad0a..HEAD -- internal/ cmd/ scripts/ tools/ .github/` 输出 **0 行**（代码面从起手到现在零字节变动），六枚生产件 md5 HEAD 对工作树逐枚 `same=yes`（见第 5 节）。分支 `dev`。
- 尺册：`D:\tmp\wisp\20261021\356a1\sh\*.sh`（10/11/20/21/30/31/40/50 八枚）；读数 `logs/`；探针**只**写在 `D:\tmp\wisp\20261021\356a1\copy\`（`git archive HEAD` 导出的仓外私有副本，**无 `.git`**），两枚 `zzz_356a1_*_test.go`。共享树 `git status --porcelain -- internal/ cmd/ scripts/ tools/ .github/` 起手与终态皆 **0 行**。

## 1. 票面前提新鲜度（逐条；未采信 `174-v1` 任何数）

尺 `sh/11-premise-full.sh`，读数 `logs/11-premise-full.log`。候选路径逐枚 `git cat-file -e HEAD:<path>`：

| 票面点名 | `git cat-file -e HEAD:` |
|---|---|
| `internal/tools/task.go` / `internal/agent/spill.go` / `docs/evidence/s1/174-v1-verdict.md` | AT_HEAD |
| `internal/risk/{mode,assessor,rules_gateway}.go` 与 `internal/tools/{fs,bridge,paths,refusal}.go` | AT_HEAD（六个逐枚） |
| `internal/tools/task_output_pointer_two_roads_174r1_test.go` / `internal/observe/thresholds.go` / `tools/d22scan/main.go` / `tools/d22scan/allowlist.txt` / `docs/PLAN.md` | AT_HEAD |
| **幻影** `internal/tools/task_output_pointer_notice_174rr1_test.go` | **MISSING**（HEAD 0 命中、盘上 0 命中；而它作为字面出现在 HEAD 的 `task.go` 里 **1** 次） |

票面可判定断言，**对 HEAD blob** 复量（不是对工作树）：

| 断言（票面原话） | HEAD 现量 | 判 |
|---|---|---|
| allowlist 支只看 `InAllowlist`、不看 R3 | `task.go:847` 逐字 `case !d.Paths.InAllowlist(canon):`；`pointerNotice` 全体（`831-871` 行）里 `risk.` 出现 **0** 次 | 真 |
| `174-r1` 那句逐字（`task.go:858-860`） | 三行拼接后与票面引文一致；`才读得回来` **0** 次、`要么批下那一张卡` **1** 次 | 真 |
| `task.go:852` 写着幻影文件名 | 该行逐字 `// task_output_pointer_notice_174rr1_test.go measures that on the same` | 真 |
| `a8eb0d51` 是那句的来源 | 对象存在；`--name-only` 恰 2 枚：`internal/tools/task.go`、`internal/tools/task_output_pointer_two_roads_174r1_test.go` | 真 |
| R3 A 档「绝对禁止、任何授权都不可豁免」 | `rules_gateway.go` R3 支：`case TierA:` 出 `level: Deny`、`"R3: 命中 A 档敏感路径，拒绝（不可豁免）: %s"` | 真 |
| 硬约束里的通配树 | `internal/risk/` 49 枚 / `internal/panel/` 46 / `internal/tools/` 94 / `internal/agent/` 82 / `cmd/wisp/` 139 / `frontend/` 116 / `design/` 14 / `tools/d22scan/` 9，均在 HEAD | 真 |
| 「19 枚 gofumpt 脏件」（票面硬约束） | 本机 `gofumpt` **不在 PATH**，该枚**不可判**（未拿它当读数） | 未判 |

**前提成立，无过期，不停手。**（唯一未判项已具名：`gofumpt` 缺席。）

## 2. 核心：R3 A 档那一发的成对读数（承重那一跳 = approval 计数）

**判级路径不是造的**：用本仓既有生产分类器 `risk.Classify`（`internal/tools/paths.go` 的 `sensitiveClassifier.Classify` 就是它的 seam，逐句直译 ClassA 到 TierA/ClassB 到 TierB）。对 `...\.ssh\id_ed25519_356a1` 返回 `A`，`risk.Gate` 给 `Allow=false Class=A reason="A-list sensitive path (non-overridable): ~/.ssh/**"`。桥用真 `*PathCanonicalizer`（真 `allowed_dirs` 根，**目标在根外**）、真 C19 assessor、真 `fs.read`、真 `task.output`；`Gate` 只做**计数与答复**（seam），它不替路径判定。

| 读数 | 值（逐字） |
|---|---|
| `task.output` notice（`announceOf`） | `[...输出已落文件：省略 17200 字符，总长 20000 字节 / 约 5000 token；注意：这条路径现在读不到，它不在你被授权的目录范围内：fs.read 会要一张 L2 卡，没人批就是直接拒；回来的路有两条——要么用户把所属目录加进 [fs] allowed_dirs，要么批下那一张卡，全文见 C:\...\.ssh\id_ed25519_356a1...]` |
| notice 本体（至 `，全文见 ` 之前） | **243 字节**（含前导 `；` 为 **246**，与 `174-r1` 的 246=246 同形，见第 4 节） |
| `fs.read` `isError` | **true** |
| `fs.read` `level`（`ToolOutcome.RiskLevel`） | **`L2`** —— **不是 `Deny`**：`levelString`（`bridge.go:1383`）把 `Deny` 折进 `L2` 一列（账表 `risk_level` 只认 L0/L1/L2 三值，`memory/models.go:130-133`） |
| `fs.read` 真判级（`OnDecision`） | `Level=Deny`、`RulesHit=[R2 R3]`、`DecisionColumn=reject` |
| `fs.read` `class` | **`permission_denied`** |
| `fs.read` 正文前 40 字节（hex） | `e8af a5e7 9bae e6a0 87e8 a2ab e7bb 9de5 afb9 e7a6 81e6 ada2 e8ae bfe9 97ae efbc 8852 3320 4120 e6a1` = `该目标被绝对禁止访问（R3 A ` 加「档」字的**前两字节**；byte[39]=`0xA1` 是续字节，40 字节处切在 rune 中间。正文长度随 fixture 临时目录名的位数浮动（多次跑读到 337 / 339 字节），前 44 字节逐次相同 |
| `fs.read` 正文逐字 | `该目标被绝对禁止访问（R3 A 档，任何授权都不可豁免）：R2: 目标路径在授权目录之外: C:\...\.ssh\id_ed25519_356a1; R3: 命中 A 档敏感路径，拒绝（不可豁免）: C:\...\.ssh\id_ed25519_356a1` |
| **`approval` 被问次数** | **0**（`fs.read` 前后 `windows 0->0 / approvals 0->0`，累计 0；`task.output` 那一发也 0） |
| 机器可辨值（对比用） | `refusalKindOf(..., Deny)` = **`redteam-denied`** |

**同一发：notice 说「会要一张 L2 卡、批下那一张卡」，`fs.read` 实测 `permission_denied`、approval 0 次、真判级 `Deny`。卡不存在。**

**真机正控（证明这不是 fixture 角落）**：用本机**真实存在的** `C:\Users\swq\.git-credentials`（120 字节，`Classify=A`，未搬 `USERPROFILE`、未造锚点）重跑：notice 同为 243 字节且逐字相同；`fs.read` `isError=true level="L2" class="permission_denied"` 233 字节、前 40 字节同形；**approval=0**；真判级 `Deny [R2 R3]`。

**第三方旁证（本仓既有、已验收的判据，本腿现跑，不是我的探针）**：`TestSensitiveFileIsDeniedNotEscalated`（`fs_test.go:157`）与 `TestBridgeRefusesTheRealShortNameOfAnAListFile`（`bridge_junction_windows_test.go:475`）各 `--- PASS`，两者断言逐字相同：A 档读必须 `permission_denied` 且 `window=0 approval=0`。也就是说 `approval=0` 这个承重读数不是我这台子造出来的。

## 3. 正控：同一台子上的 R2-only 形状，`approval` = **1**（尺分得清两种形状）

同一桥、同一 `*PathCanonicalizer`（同一根外根）、同一读法，只把目标换成 `...\tool-output-356a1-r2only.txt`（`Classify=none`，在根外）。
**口径说明**：票面写的是"同一路径改到 R2-only 形状"，但一条路径不可能同时「命中 A 档」又「在 R3 之外」——`Classify` 的 A 档是按锚点/文件名匹配的，改拼法不会把它移出 A。所以这里把「同一路径」落实为「**同一台子、同一根外形状，只让 R3 不命中**」，并把两个目标各自的分类读数都打出来（`A` 对 `none`），让这一步是可核的而不是被默认的。

| 读数 | R3 A 档（第 2 节） | R2-only（正控） |
|---|---|---|
| notice 本体 | 243 字节 | **243 字节，逐字节相同** |
| `fs.read isError` | true | **false** |
| `fs.read level` | `L2` | `L2` |
| `fs.read class` | `permission_denied` | **空串** |
| 正文前 40 字节 | `该目标被绝对禁止访问（R3 A ...` | `abcdefghijklmnopqrstuvwxyzabcdefghijklmn`（全文 20000 字节 = 登记原文全长） |
| **`approval` 次数** | **0** | **1** |
| 真判级 / `refusalKindOf` | `Deny [R2 R3]` / `redteam-denied` | `L2 [R2]` / `out-of-allowlist` |

**正控得 1，台子搭对了**：两种形状只在 `approval`、`class`、`isError`、`refusalKindOf` 上分开，**notice 完全分不开**。

**旁证（不在票面要求内，供 `AC#1` 参考）**：把 home 树整棵加进 `allowed_dirs`（目标变 in-scope）后，`task.output` **完全沉默**（notice 0 字节），`fs.read` 仍是 `Deny [R3]`、approval 0。也就是**第一条路（加目录）对 R3 A 档同样无效**，notice 那两条路在 R3 形状下**两条都假**，且这一发连「读不到」都不说。

## 4. 反控：两形状 notice 原文并排，逐字指出 R3 形状下错在哪几个字

`task.output` 两侧 notice **逐字节相同（243 = 243，每次跑都相同）**；`announceOf` 的差（428 对 435、426 对 433 等，随 fixture 临时目录名的位数浮动）只来自各自 fixture 的 `全文见 <path>`，那是指针句（`PLAN.md:2564` 冻结），不算回执。

```
R3 A 档：注意：这条路径现在读不到，它不在你被授权的目录范围内：fs.read 会要一张 L2 卡，没人批就是直接拒；回来的路有两条——要么用户把所属目录加进 [fs] allowed_dirs，要么批下那一张卡
R2-only：注意：这条路径现在读不到，它不在你被授权的目录范围内：fs.read 会要一张 L2 卡，没人批就是直接拒；回来的路有两条——要么用户把所属目录加进 [fs] allowed_dirs，要么批下那一张卡
```

逐字点错（偏移按 243 字节本体，0 基）：

| 偏移 | 逐字 | R3 形状下的事实 | 判 |
|---|---|---|---|
| 81 | `fs.read 会要一张 L2 卡` | approval 实测 **0** 次，**没有卡** | **错** |
| 111 | `没人批就是直接拒` | 「批」这一步**根本不会发生**（不是「没人批」） | **错**（把不存在的前提写成条件） |
| 138 | `回来的路有两条` | 真判级 `Deny`，**一条也没有** | **错** |
| 165 | `要么用户把所属目录加进 [fs] allowed_dirs` | 已实测：加进去后仍 `Deny [R3]`、仍拒 | **错**（路一也假） |
| 219 | `要么批下那一张卡` | 没有卡可批，批了也不生效 | **错**（点一张永不出示的卡） |
| 0 / 39 | `注意：这条路径现在读不到` / `它不在你被授权的目录范围内` | 读不到为真；但「不在授权目录范围内」只是**原因之一**（R3 压过它） | 半真（`Deny` 时这句不是充分解释） |

## 5. 门禁（本腿现跑；跑前 `unset GOTMPDIR TMPDIR TEMP TMP`、`GOFLAGS=`，`go1.27.1 windows/amd64`）

| 门 | rc | 末行 / 输出 |
|---|---|---|
| `GOFLAGS= go build ./...` | **0** | 无输出（0 行） |
| `GOFLAGS= go vet ./internal/tools/` | **0** | 无输出（0 行） |
| `GOFLAGS= go test -count=1 -run '174r1' -v ./internal/tools/` | **0** | 首跑 `ok  github.com/CarlosShao/wisp/internal/tools  0.424s`；落件前后多次复跑 `0.162s` / `0.126s`，四数逐次相同 |
| `GOFLAGS= go test -count=1 -run 'TestPointer|TestTaskOutputPointer|TestHealthyPointer|TestSensitiveFileIsDeniedNotEscalated|TestEmptyAllowlistAuthorizesNothing' -v ./internal/tools/` | **0** | `ok  ...  0.723s`（该选择器共跑起来 23 名，四数见下） |

**四数分开（不相减）**：`174r1` 选择器 `RUN=4 / PASS=4 / FAIL=0 / SKIP=0`；宽选择器 `RUN=23 / PASS=23 / FAIL=0 / SKIP=0`（缩进层 `^    --- PASS` 两侧皆 **0**，与上面不同分母）。
代码面 md5（HEAD blob 对工作树，六枚全 `same=yes`）：`task.go 432a52be...`、`bridge.go 0f754371...`、`fs.go cc22d191...`、`paths.go adeb51d9...`、`rules_gateway.go acec8584...`、`assessor.go 94271869...`。

## 6. 红名册（逐名）

`^(---|    ---) (FAIL|SKIP)` 在本腿每一份 `-v` 日志（`logs/30-gates*.log`、`logs/31-gates-wide.log`、`logs/2*-paired*.log`）上命中 **0 行**，称**无一名需要归属**（本腿没造红，也没既有红）。
实际跑起来的 4 名（`174r1` 选择器全量）：`TestPointerOutsideAllowlistNamesEveryRoadBack174r1`、`TestPointerSecondRoadBackIsRealOnThisBridge174r1`、`TestPointerSecondRoadBackIsBookedAsACard174r1`、`TestPointerTwoRoadNoticeLeaksNoPathOrRoots174r1`。宽选择器的 23 名逐名见 `logs/31-gates-wide.log`。

## 7. 提交号 / 逐枚 `.go` / 越界枚数 / 正控尺

- 本腿落件：**只** `docs/evidence/s1/356-a1-census.md`（本件）。逐枚 `.go` = **0**；越界枚数 = **0**。
- 提交：**首笔（承重那笔）＝ `40fddae0`**，`1 file changed, 115 insertions(+)`，`A docs/evidence/s1/356-a1-census.md`。之后每一笔都是对同一文件同一路径的校字笔（本件把"读一次就写死"的口径逐处收掉的过程都留在主题里）。
  权威名册不是本文手写的列表，而是这条可复量命令的输出（**笔次不在此写死**，否则每补一笔就自我过期）：
  `git log --oneline 40fddae0^..HEAD -- docs/evidence/s1/356-a1-census.md`
  ——每一笔的主题都以「票 356 普查腿 356-a1」开头，`--name-only` 都只列本文件一枚，逐笔 `.go` 都是 0。
- 每一笔的 `git add` 与 `git commit` 两处都带**显式 pathspec** `-- docs/evidence/s1/356-a1-census.md`；全程无 `add -A`、无 `add .`、无 amend/reset/rebase/stash/checkout/clean。
- 逐笔 `.go` = **0**、越界 = **0**（`logs/50-writeface-final.log` 对本腿每一笔逐笔打一遍，每一笔都是 `files touched = 1 / .go files = 0 / IN RANGE = 1 / OUT OF RANGE = 0`；正控在同一文件每一轮都打 `a8eb0d51` 读回 `.go = 2`、`OUT OF RANGE = 2`）。
- 正控尺（证明「数 `.go` 的尺」不是恒 0）：同一把 `sh/50-writeface.sh` 打 `a8eb0d51`（`174-r1`）读回 `.go = 2`、`IN RANGE = 0`、`OUT OF RANGE = 2`，尺会红。
- 落件之外的写面：`D:\tmp\wisp\20261021\356a1\**`（`sh/` 八枚、`logs/` 十来枚、`copy/` 一份 HEAD 导出加两枚探针测试件、`commit-msg-*.txt`）。临时件**只建不删**；无 `.out` 后缀、无 `-overlay`。

**本腿未做（具名边界）**：`-race` 未跑；`cmd/wisp` 端到端 `wisp run` 未跑；`gofumpt` 不在 PATH（票面「19 枚脏件」未复量）；`fs.list` 与写半边未量（不在 `AC#0` 射程）。
