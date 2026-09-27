# 票 162 · AC#5（风险档与门控）＋ AC#6（契约轴） — 非实现者对抗验收表 `162-v3`

- 验收程＝`162-v3`（非实现者，只裁不改）｜派单＝`.scratch/wisp/dispatches/2026-09-27-175x-accept-162-r4-v3-tier-and-contract-axis.md`
- 票面＝`.scratch/wisp/issues/162-add-the-patch-one-region-tool-instead-of-rewriting-whole-files.md`
- 被验证据件＝`docs/evidence/s1/162-risk-tier-and-contract-axis-r4.md`（现量 236 行，与锚点 blob 一致）
- **被验锚点＝`fd82107`**（父链 `7a41af51 → 27f0cb8a → 9c00bb33 → b542bc1c → fd82107`，中间与 `171-r2` 的 `7f78a163`／`290aa63e` 交错在同一共享树上）
- ⚠ 全程被验码一律 `git show <锚>:<路径>`／`git grep <pat> <锚>`，禁读工作树（同树另一程在写 `probes/**`）；未跑 `probes/154/gate-clauses.sh`、`probes/161/r6/flip-declaration.sh`（在飞）。

## 1. step-0 起手五件

| 件 | 命令 | 读数 |
|---|---|---|
| 时间 | `date` | `Sun Sep 27 18:01:35 CST 2026` |
| 分支 | `git rev-parse --abbrev-ref HEAD` | `dev`（未用 switch/checkout/merge/rebase/reset/stash/worktree/clean） |
| 取数时 tip | `git rev-parse HEAD` | `ace65147b5c2c1e2f5e4882bed9a9bd059b7be`（因 `171-r2` 在飞会继续前移；本表判语只钉 `fd82107`） |
| 写面基线 | `git status --porcelain -- internal/tools/ docs/evidence/s1/162-*` | **空** |
| 锚点串 | `git log --format='%h %ci %s' -6 fd82107` | 见下（逐枚 `git cat-file -t` 已自证为 commit） |

```
fd82107 2026-09-27 17:53:58 +0800 162-r4 收尾：票面 162 追加 17:4x 进度一节（零勾）＋证据件 §11.1
290aa63e 2026-09-27 17:53:50 +0800 票 171（171-r2 写码位·件③）……   ← 另一程
b542bc1c 2026-09-27 17:50:14 +0800 162-r4 AC#6: 契约轴凭证＋AC#5／§3 交件件落盘
9c00bb33 2026-09-27 17:46:29 +0800 162-r4 AC#5 收尾：gofumpt 甲形归零
7f78a163 2026-09-27 17:45:45 +0800 票 171（171-r2 写码位·件②）……   ← 另一程
27f0cb8a 2026-09-27 17:43:34 +0800 162-r4 派单 §3: 死进程暂存残件升成用例钉
```

父链核对（`git log --format='%h %p' -6 fd82107`）：`fd82107←290aa63e←b542bc1c←9c00bb33←7f78a163←27f0cb8a←7a41af51`——`162-r4` 的交件与 `171-r2` 的写在同一线性历史上交错，两程互不越界（各自动自己的写面），本表按提交逐枚归位。

## 2. 本程没测什么（先声明，随后补）

- 只在 `./internal/tools/` scope 取数（禁全仓 `go test ./...`）；`internal/risk` 单元尺未自跑（本程不碰那块地，只 `git show` 读它的源码定档位来源）。
- 真人 L2 卡片未经人手：沿用被验程的 `gateSpy` 形状，本程判语只到"走到哪条通道、卡片带什么规则/审计行"。

## 8(片段). 越界核查＝AC#6 契约轴（本程自己的筛法，非复用被验程 grep）

判据（票面 AC#6 原文＋09:4x `>`②）：本票**只许**动 `docs/PLAN.md` 的 D34 那一行与 `SPEC-07` 镜像那一行（**两行均已由编排者落完 ⇒ 实现程再动＝越权**）＋ `internal/tools/**` ＋自己证据件；`docs/specs/**`、`internal/risk/**`、`thresholds.go`、golden、`frontend/**`、`design/**` 零字节。

被验程**自己**的四枚提交逐枚 `--name-only`（`git -c core.quotePath=false show --name-only <c>`）：

| commit | 动到的文件 |
|---|---|
| `7a41af51` | `internal/tools/fs_edit_ac5_gate_r4_test.go` |
| `27f0cb8a` | `internal/tools/fs_edit_ac5_sweep_r4_windows_test.go` |
| `9c00bb33` | `internal/tools/fs_edit_ac5_gate_r4_test.go`（gofumpt 形状） |
| `b542bc1c` | `docs/evidence/s1/162-risk-tier-and-contract-axis-r4.md` |
| `fd82107` | 票面 162 ＋ 上面那枚证据件 |

⇒ 162-r4 的写面**只有** `internal/tools/**` 两枚测试＋自己证据件＋票面；契约轴文件**一枚未动**。

独立筛法（本程自造，与被验程式样不同）——对窗口两端直接问契约轴各路径的净改动：

```
git -c core.quotePath=false diff --numstat f1b99a70 fd82107 -- internal/risk tools/d22scan/allowlist.txt   → 空
git -c core.quotePath=false diff --numstat f1b99a70 fd82107 -- docs/PLAN.md docs/specs                      → 空
```

窗口全量 `--numstat f1b99a70 fd82107`（供点名，⚠ 含 `171-r2` 的 `probes/**`，非本程越界）：162-r4 的四枚测试/证据净改动**删除列一律 0**；`9c00bb33` 的 `8 4` 里那 4 枚"删除"＝两处 composite literal 换行重排（逐行验：记号集完全相同，只把 `reading{...}` 拆成多行，无内容行被丢）。

两行冻结文本现量在位、一字未由本票实现程改动（`git show fd82107:docs/PLAN.md`／`SPEC-07…`）：
- `PLAN.md:2536` = `| **fs.edit** | 字面定位替换已存在文件里的一小段 | **L2** | fs.write | S3 | … **R8 覆盖已存在 → L2**，与上面 fs.write 覆盖行同级 |`
- `SPEC-07…:43` = `| **fs.edit** | 字面定位替换已存在文件里的一小段 | **L2** | fs.write | S3 | 票 162 新增，镜像 PLAN.md D34；… |`

**AC#6 初判〔成立〕**（终判并入下方，等门禁数复跑）：契约轴对 162-r4 **零越界**，两行冻结文本未动。
