# 356 — 指针回执点了一张**永不出示的卡**：目标落进 R3 拒绝档时，notice 仍说"批下那一张卡"

- Status: **ready-for-agent**（`AC#0` 是只读现量，可先派；`AC#1` 动文案面、`AC#2` 是牙）
- 立案：2026-10-21 10:2x +08 编排者立。**非实现者签名**＝`174-v1` 判语第 (5) 件第 2 条（`docs/evidence/s1/174-v1-verdict.md`）。
- 来源：票 174 落地腿 `174-r1`（`a8eb0d51`）把 allowlist 支回执改成"回来的路有两条"，方向对，但**写过头了半句**；
  `174-v1` 现跑量到反形并判"越界"，同时指出第 1 枚判据在 MUT-5 下**不承重第二条路**。
- 关联：票 174（本票是它那句文案的**收口**；票 174 的 `AC#2b` 早已 `[x]`，⛔ 本票⛔ 挂回它任何格）· 票 322（同族"窗口测不到"的先例形）

## 这是什么（人话）

`task.output` 截断时会给模型一句"注意：这条路径现在读不到……"，并告诉它怎么把内容读回来。
`174-r1` 把这句从"只有一条路（加进 `[fs] allowed_dirs`）"改成"两条路（加目录，**或**批一张 L2 卡）"——方向是对的：
目标只被判到 **R2**（授权目录之外）时，批一张卡**真能把字节读走**（`174-v1` 四臂读数：`NoGate` 拒 36 B ⇄ 批准 20000 B 整读回）。

问题在**另一种形状**：目标落进 **R3 A 档**（绝对禁止、任何授权都不可豁免）时，`fs.read` 直接给
`permission_denied`、**approval 被问 0 次**、根本没有卡。可 notice 那句话照样让模型"批下那一张卡"——**点一张永远不会出现的卡**。
`task.output` 的 allowlist 支只看 `InAllowlist`（`task.go:847`）、**不看 R3**，所以它分不出这两种形状。

## 结案判据

- [ ] **AC#0 现读＋归因（只读，⛔ 动文案字节）**：把"R3 A 档下 notice 说了什么 ⇄ `fs.read` 实际做了什么"量成成对读数。
      必须**逐字**给：notice 原文（`announceOf`）、`fs.read` 的 `isError`／`level`／`class`／正文逐字、
      **`approval` 被问的次数（今天应为 0）**。★另起尺、⛔ 抄 `174-v1` 的数；★正控＝同一路径在 R2-only 形状下 approval＝1（证明尺分得清两种形状）。
- [ ] **AC#1 改写口径（文案面）**：让回执能区分"只判到 L2（两条路都真）"与"被判 Deny（没有卡、批也无效）"。
      ★**范围扩充（机主 2026-10-21 14:2x 拍「扩」）＝必须同时覆盖第三种形状：目标在 allowlist 内、但被判 Deny**
      （`普查腿 356-a1` 现量——把 home 树整棵加进 `allowed_dirs` 后，`fs.read` 仍 `Deny [R3]`、
      `approval` 被问 **0 次**，而 `task.output` **完全沉默、notice 0 字节**）。我复核过 `pointerNotice` 的分支结构：
      它只有 ①判级者未接线 ②`err != nil` ③`!InAllowlist` ④`os.Stat` 失败 ⑤非普通文件——**没有任何一支处理"in-scope 但被判 Deny"**。
      ⇒ 本格要求的三种形状＝**（i）只判到 L2**（两条路都真）／**（ii）被判 Deny 且在 allowlist 外**（没有卡、加目录也无效）／
      **（iii）被判 Deny 但在 allowlist 内**（今天连"读不到"都不说）。
      ★判据必须用 `refusalKindOf` 那一对值（`redteam-denied` ⇄ `out-of-allowlist`）**区分形状**，
      ⛔ 用 `RiskLevel`（我复核过 `internal/tools/bridge.go:1379-1386` 的 `levelString` 把 `risk.Deny` 折进 `L2` 一列 ⇒ 分不出）。
      落点＝`internal/tools/task.go` 的 `pointerNotice`（或它的同级读取）；★改后必须使**两条腿仍逐字节同形**
      （`internal/agent/spill.go` 的同族回执），判据＝两腿 notice 区逐字节 `cmp` 相同（`174-r1` 量到 246=246，本票要保住这个性质）。
      ⛔ 只在**文案与判据**面动；⛔ 动 R3/R2 判级语义、⛔ 改 `PLAN.md`／`docs/specs/**`、⛔ 新增授权能力。
- [ ] **AC#2 牙（补 `174` 第 1 枚判据的 MUT-5 缺口）**：加一枚**只属第二路**的标记——
      现判据的标记集是 `{不在你被授权的目录范围内, [fs] allowed_dirs, L2 卡}`，
      `174-v1` 的 MUT-5 证明**把"要么批下那一张卡"整句删掉，四枚判据仍全绿**（我另尺复现同结论）。
      ⇒ 新判据必须让"删掉第二路整句"**转红并逐名报红行**；变形先证落地（`grep -n` ＋ `go build` rc=0）再读数，读完单独还原并核 blob 回 HEAD。

## 硬约束

- ⛔ 翻票 174 任何格、⛔ 改票 174 的 `Status:`／票名。
- ⛔ 碰 `internal/risk/**`、`internal/panel/**`、`thresholds.go`、任何 golden、`allowlist.txt`、`frontend/**`、`design/**`、`cmd/wisp/**`。
- ★**幻影引用**：`task.go:852` 注释里写着 `task_output_pointer_notice_174rr1_test.go`，**该文件不存在**（盘与 HEAD 各 0 命中）。
  `d22scan` 没抓到是因为 ban #9 的 `repoPathRe` 只认**带目录前缀的完整仓相对路径**、裸文件名进不去。⇒ 本票**顺手把它改对**（改注释、不留幻影），并**具名登记这条仪器缺口**（⛔ 改 `tools/d22scan/**`）。
- ⛔ 洗那 19 枚存量 gofumpt 脏件（另案）。

## 排程

`356-a1`（只读＝`AC#0` 成对读数）→ `356-r1`（写腿＝`AC#1`＋`AC#2`）→ `356-v1`（非实现者裁）。

**Status 追加（2026-10-21 10:2x）**：号位＝编排者定（`356` 盘上现量空号）。⛔ 零翻框、零 push。

## Progress log

（追加式；每条 `- [日期 时刻 +08] agent=… did=… next=…`，钟点由 `date` 插值，不手打。）

- [2026-10-21 10:2x +08] agent=`编排者` did=立案（非实现者签名＝`174-v1` 第 (5) 件第 2 条；三格按它的建议形状落，含 MUT-5 缺口）next=派 `356-a1` 只读出 `AC#0` 成对读数（approval 计数是承重那跳）。
