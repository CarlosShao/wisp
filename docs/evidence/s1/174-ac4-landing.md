# 174 AC#4 — artifacts 树加入 fs 授权根（落地件，2026-11-12）

写腿：编排者本人（原派 glm 腿 `174-landing` 返回空、无盘上痕迹 ⇒ 按铁律⓪⓪⓪ 自己落地）。
票面：`.scratch/wisp/issues/174-task-output-points-at-a-file-the-model-cannot-read-...md` AC#4。
作业书：`docs/evidence/s1/174-landing-readiness-v1.md`（筛查腿，含最小改点/两向判据/四条风险）。

## 0. 授权（白纸黑字，⛔ 不需要再要人）

- 台账 `docs/reports/pending-and-issues.md:31348` 逐字：`**批复**：A(56/58/60)全批；B(211/215/256/345/350/354/174/177/138/150/148)全批；C 除 55(macOS 暂不做) 全批`
- 同件 `:31352` 逐字：`契约类改动（211 D38d / 215 D36 / 256 C18 / **174 [fs] allowed_dirs** / 168 C17C24 / 350 门禁 / 354 拍形 / 177 选支）按本次批复执行，改契约文字那步照批`
- 批项原文 `docs/evidence/s1/census-199.md:102`：「要不要把 artifacts 加进授权根（改 `[fs] allowed_dirs` 默认值的取舍）」
⇒ 批准覆盖＝「把 artifacts 加进授权根」。且筛查腿已判 `PLAN.md:431`/`:2564` **零字节不改**（落地后"按需再读"自然成真）——**本次落地对 `PLAN.md`/`docs/specs/**` 零字节**。

## 1. 改点（唯一一处）

`cmd/wisp/run.go:762-789`（现读）：

```go
allowed := make([]string, 0, len(cfg.FS.AllowedDirs)+1)
allowed = append(allowed, cfg.FS.AllowedDirs...)
// [注释：授权来路 + 为什么用 mem.ArtifactsDir() + 目录已由 memory.Open 建好]
allowed = append(allowed, mem.ArtifactsDir())
rt.paths = tools.NewPathCanonicalizer(allowed, cfg.FS.ReparsePointExceptions)
```

**三条 gotcha 的处置（筛查腿 §④ 具名，逐条落地）**：

| gotcha | 处置 | 凭据 |
|---|---|---|
| ① 拼写对不齐（根 vs spill 目标） | **用 `mem.ArtifactsDir()`** 而非 `filepath.Join(rt.spec.dataDir,"artifacts")`。理由：`memory.Open` 内部先 `filepath.Abs(dir)`（`internal/memory/open.go:263`）再由该绝对形构造 `artifactsDir`（`:277`），而 `rt.spec.dataDir` 是调用方传入的原始值。问 store 它实际写哪个目录，根与目标必然同拼写 | `internal/memory/open.go:263/:277` 现读 |
| ② 目录可能还不存在（不存在的根会被 `PathCanonicalizer` 丢进 `unusable`＝静默不授权） | **不需要预创建**：`memory.Open` 在 `run.go:691` 就被调用，它在 `open.go:283` 用 `winsec.PrivateDirAll(s.artifactsDir, 0o700)` 已建好并封好该目录，早于 `:789` 构造 canonicalizer ⇒ 目录必然在盘 | `cmd/wisp/run.go:691` ⇄ `internal/memory/open.go:283` |
| ③ workspace 收窄是设计内行为 | **未绕**：`internal/tools/paths.go:42-47` 逐字「It can only ever TIGHTEN InAllowlist」，票 92 AC#3 收窄面在飞时 artifacts 根会被收窄掉，这是对的 | `internal/tools/paths.go:42-47` 现读 |

## 2. 两向判据＋变异控制（`internal/tools/artifacts_authorized_174_test.go`，新件）

| 判据 | 断言 | 读数 |
|---|---|---|
| **正向** `Test174SpilledArtifactIsReadableThroughTheSameCanonicalizer` | ① `Roots()` 含 artifacts 根（防"静默被丢进 unusable"）② 真 `Spiller.Prepare` spill 出 artifact ③ 同一 C26 实例 `fs.read` ⇒ `IsError=false`、`RiskLevel="L0"`、正文含 `artifact-body-174` | **PASS** (0.05s) |
| **反向** `Test174DataRootSiblingsStayOutsideTheArtifactGrant` | artifacts 的三个兄弟项（`secrets\token.txt`／`config.toml`／`wisp.db`）⇒ 逐枚 `IsError=true`、`RiskLevel="L2"` | **PASS** (0.04s) |
| **变异控制** `Test174NoArtifactRootMeansTheArtifactIsUnreadable` | 同一枚 spill 出的 artifact，用**改动前**的根集（不传 artifacts 根）读 ⇒ 必红 `IsError=true`／`L2` | **PASS** (0.02s) |

⇒ **判据有牙**：传根读得回 ⇄ 不传根读不回，两向都钉住。第三枚是自证控制，⛔ 为绿改断言。

## 3. 门禁（我亲跑）

| 门 | 读数 |
|---|---|
| `go build ./cmd/wisp` | rc=0 |
| `go vet ./cmd/wisp` | rc=0 |
| `go vet ./internal/tools` | rc=0 |
| `gofumpt -l cmd/wisp/run.go` | 空（rc=0） |
| `sh scripts/d22scan.sh` | **clean**，rc=0（`ban #8 internal/` examined **1005** Go files、`cmd/` 276、`frontend/` 123、`design/` 39；零违规） |
| `go test ./internal/tools -run Test174 -count=1` | **PASS=3 FAIL=0** rc=0（0.187s） |

## 4. AC#5 两行分账草稿（随件，归编排者结题合入）

> - **本票（授权根那一处断口）**：AC#4 落地后，`task.output` 指向的 artifact 在同一 C26 实例下可读回（正向判据 PASS），越界仍 L2（反向判据 PASS）。**本票管的是"产物文件在不在授权根里"这一处。**
> - **`Q-59`（偏移参数那一处断口）**：**与本票是同一根管子的另一处，未因本票而通**。本票零触碰偏移参数面；`Q-59` 状态按台账 `A1571`（Q-59 已获机主批复、解锁 174/177 等票）继续由它自己的票跟踪。**⛔ 用"AC#4 修好了所以续读问题解决了"顶掉 Q-59。**

## 5. 提交面（显式 pathspec）

- 改：`cmd/wisp/run.go`（仅 `:762-789` 一带）
- 新增：`internal/tools/artifacts_authorized_174_test.go`、本件
- ⛔ 未碰：`PLAN.md`／`docs/specs/**`／`thresholds.go`／golden／`allowlist.txt`／`internal/tools` 产码／`frontend/**`／`design/**`
- ⛔ 未翻 AC 框（勾归非实现者＝编排者，本件只交料）
