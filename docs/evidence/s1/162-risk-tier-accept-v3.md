# 票 162 · AC#5（风险档与门控）＋ AC#6（契约轴） — 非实现者对抗验收表 `162-v3`

- 验收程＝`162-v3`（非实现者，只裁不改）｜派单＝`.scratch/wisp/dispatches/2026-09-27-175x-accept-162-r4-v3-tier-and-contract-axis.md`
- 票面＝`.scratch/wisp/issues/162-add-the-patch-one-region-tool-instead-of-rewriting-whole-files.md`
- 被验证据件＝`docs/evidence/s1/162-risk-tier-and-contract-axis-r4.md`（现量 236 行，与锚点 blob 一致）
- **被验锚点＝`fd8201e7`**（父链 `7a41af51 → 27f0cb8a → 9c00bb33 → b542bc1c → fd8201e7`，中间与 `171-r2` 的 `7f78a163`／`290aa63e` 交错在同一共享树上）
- ⚠ 全程被验码一律 `git show <锚>:<路径>`／`git grep <pat> <锚>`，禁读工作树（同树另一程在写 `probes/**`）；未跑 `probes/154/gate-clauses.sh`、`probes/161/r6/flip-declaration.sh`（在飞）。

## 1. step-0 起手五件

| 件 | 命令 | 读数 |
|---|---|---|
| 时间 | `date` | `Sun Sep 27 18:01:35 CST 2026` |
| 分支 | `git rev-parse --abbrev-ref HEAD` | `dev`（未用 switch/checkout/merge/rebase/reset/stash/worktree/clean） |
| 取数时 tip | `git rev-parse HEAD` | `ace65a5147b5c2c1e2f5e4882bed9a9bd059b7be`（因 `171-r2` 在飞会继续前移；本表判语只钉 `fd8201e7`） |
| 写面基线 | `git status --porcelain -- internal/tools/ docs/evidence/s1/162-*` | **空** |
| 锚点串 | `git log --format='%h %ci %s' -6 fd8201e7` | 见下（逐枚 `git cat-file -t` 已自证为 commit） |

```
fd8201e7 2026-09-27 17:53:58 +0800 162-r4 收尾：票面 162 追加 17:4x 进度一节（零勾）＋证据件 §11.1
290aa63e 2026-09-27 17:53:50 +0800 票 171（171-r2 写码位·件③）……   ← 另一程
b542bc1c 2026-09-27 17:50:14 +0800 162-r4 AC#6: 契约轴凭证＋AC#5／§3 交件件落盘
9c00bb33 2026-09-27 17:46:29 +0800 162-r4 AC#5 收尾：gofumpt 甲形归零
7f78a163 2026-09-27 17:45:45 +0800 票 171（171-r2 写码位·件②）……   ← 另一程
27f0cb8a 2026-09-27 17:43:34 +0800 162-r4 派单 §3: 死进程暂存残件升成用例钉
```

父链核对（`git log --format='%h %p' -6 fd8201e7`）：`fd8201e7←290aa63e←b542bc1c←9c00bb33←7f78a163←27f0cb8a←7a41af51`——`162-r4` 的交件与 `171-r2` 的写在同一线性历史上交错，两程互不越界（各自动自己的写面），本表按提交逐枚归位。

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
| `fd8201e7` | 票面 162 ＋ 上面那枚证据件 |

⇒ 162-r4 的写面**只有** `internal/tools/**` 两枚测试＋自己证据件＋票面；契约轴文件**一枚未动**。

独立筛法（本程自造，与被验程式样不同）——对窗口两端直接问契约轴各路径的净改动：

```
git -c core.quotePath=false diff --numstat f1b99a70 fd8201e7 -- internal/risk tools/d22scan/allowlist.txt   → 空
git -c core.quotePath=false diff --numstat f1b99a70 fd8201e7 -- docs/PLAN.md docs/specs                      → 空
```

窗口全量 `--numstat f1b99a70 fd8201e7`（供点名，⚠ 含 `171-r2` 的 `probes/**`，非本程越界）：162-r4 的四枚测试/证据净改动**删除列一律 0**；`9c00bb33` 的 `8 4` 里那 4 枚"删除"＝两处 composite literal 换行重排（逐行验：记号集完全相同，只把 `reading{...}` 拆成多行，无内容行被丢）。

两行冻结文本现量在位、一字未由本票实现程改动（`git show fd8201e7:docs/PLAN.md`／`SPEC-07…`）：
- `PLAN.md:2536` = `| **fs.edit** | 字面定位替换已存在文件里的一小段 | **L2** | fs.write | S3 | … **R8 覆盖已存在 → L2**，与上面 fs.write 覆盖行同级 |`
- `SPEC-07…:43` = `| **fs.edit** | 字面定位替换已存在文件里的一小段 | **L2** | fs.write | S3 | 票 162 新增，镜像 PLAN.md D34；… |`

**AC#6 终判〔成立〕**：契约轴对 162-r4 **零越界**（逐枚 `--name-only` 只及 `internal/tools/**` 两枚测试＋自己证据件＋票面；`internal/risk`／`tools/d22scan/allowlist.txt`／`docs/PLAN.md`／`docs/specs` 窗口净改动全为空）；两行冻结文本在位、本票实现程一字未动；`9c00bb33` 的 4 枚"删除"逐行验＝纯 composite literal 换行重排、无内容行被丢；`thresholds.go`／golden／`frontend/**`／`design/**` 窗口内未出现（窗口里出现的 `probes/154/gate-clauses.sh` 与 `probes/171/**` 属同树另一程 `171-r2`，非本票越界）。

---

## 3. AC#5 — 风险档与门控（本程主攻＝"这凭据会不会只是档位字符串相等的空心检"）

AC#5 判据（票面 `:42` 原句＋09:4x `>`① 有效部分）：**越界路径必须被拒、不许新增豁免、不动 `allowlist.txt`**；档位按 **L2** 做（"与 fs.write 同族（L1）"已被 `>`① 改判作废）。

**本程独立造的一发关键变异（被验程没做的那发）**：只把 `internal/tools/fs_edit.go:390` 的 `Declared: risk.L2` 经 `-overlay` 改成 `risk.L1`（盘上文件未动，overlay 副本在 `$TMPDIR` 丢弃树；`-overlay` 与 `-cover` 未同用）。**基线四枚（无 overlay）＝`PASS PASS PASS PASS` rc=0**；叠加 L1 档位后原样输出：

```
--- PASS: TestFSEditAndFSWriteShareTheOutOfScopeVerdict (0.01s)     ← 摘掉档位后仍绿
    fs_edit_ac5_gate_r4_test.go:191: 并排读数 fs.edit: level=L2 rules_hit=[R1 R2] ... reason="R1: 工具声明为下界（L1）; R2: 目标路径在授权目录之外: ..."
--- FAIL: TestFSEditRoutesToApprovalNotTheL1Window (0.03s)          ← 档位真被这枚钉
    fs_edit_ac5_gate_r4_test.go:261: fs.edit must reach the approval channel and never the L1 window: window=2 approval=0 (want 1/1 cumulative)
--- PASS: TestFSEditOutOfScopeIsJudgedFromTheDeclaredPathNotTheToolWord (0.02s)
--- FAIL: TestFSEditRefusesAPathC26CannotCanonicalize (0.00s)       ← 档位也被这枚钉（审计行）
    fs_edit_ac5_gate_r4_test.go:388: wanted an audit line for fs.edit carrying [risk=L2], got none in ["...tool=fs.edit risk=L1 ... rules_hit=[R1] in_allowlist_scope=true ..."]
```

⇒ overlay 确证生效（两枚绿翻红，非静默 no-op）。**三处读数一起看，AC#5 的"档位凭据"不是一句空心字符串相等，但也不是它最显眼的那枚读数在钉**：
- **越界被拒**＝真路径、非空心：`bridge.Execute` → C19 `Assess` → `rules_gateway` R2（读 C26 `InAllowlist`）→ 审批通道 → 人拒 → 回执 `user_rejected`，并配**盘上字节逐字未变**＋`noStagingFilesLeft`。这条即使把档位改 L1 也**照红**（R2 与档位无关，见 `assessor.go:16` "R2 … L2"、`fuse` 取最大严重度 `assessor.go:312`），所以"越界必须被拒"是硬凭据。
- **档位 L2 的字符串断言（并排读数里 `level=L2`）是 R2 强制出来的、不是声明出来的**：把声明下界改成 L1，并排读数的 `fs.edit: level` **仍是 L2**（reason 变成"下界（L1）"、`rules_hit` 仍含 R2 撑到 L2）⇒ **那枚最像"档位检"的断言根本不检档位**。真正钉住声明档位的是另外两枚：**in-scope 路由枚**（范围内改已存在文件必须落 L2 审批队列，L1 会翻红）＋**空白 path 审计枚**（审计行 `risk=L2` 子串，L1 翻红）。
- **不许新增豁免／不动 allowlist.txt＝成立**：`wc -l tools/d22scan/allowlist.txt`=7、`grep -c fs.edit`=0；窗口 `diff --numstat … tools/d22scan/allowlist.txt` 为空；`sh scripts/d22scan.sh` rc=0。

## 4. 派单 §2 四处复算（本程结论）

- **① `R8` 那一味差异合不合法（三选一）**：本程复算，同一枚授权树外既有文件 `fs.write`⇒`rules_hit=[R1 R2 R8]`、`fs.edit`⇒`rules_hit=[R1 R2]`（基线与 overlay 两遍皆然）。冻结文本给 `fs.edit` 定 L2 的**书面依据**是"R8 覆盖已存在 → L2"，而**机器实际**把 `fs.edit` 抬到 L2 靠的是 **R1 声明下界**（`fs_edit.go:390 Declared: risk.L2`）＋越界时的 R2，**`fs.edit` 全程不命中 R8**（它没挂 `writeFacts`）。⇒ 结论＝**成立但凭据挂错了名字**：L2 判定本身站得住（in-scope 路由与 audit 行都量到 L2），但契约文字写明的*L2 来路*（R8）与机器实际来路（R1 下界）分歧 ⇒ **这是契约文字与机器读数的分歧，应立账（建议新 `A##`：D34 表 fs.edit 行的"R8→L2"表述与 `FSEditDecl` 的"R1 下界"实现不同名），不是本票 162 的缺陷**（票面 10:2x 编排者已裁"L2 不按工具名走、R8 写进声明"）。

- **② `M-1` 判语理解＋"改冻结文本档位 L2→L1 有没有用例红"**：被验程 M-1（摘 `fs_edit.go:391 PathParams`）⇒ 红在 `rules_hit`/审计行、档位与通道不红——本程复算认同，并对"承重在哪"给出更利的一刀：**用上面那发 `Declared L2→L1` 变异**量到"档位真承重只在 in-scope 路由枚＋空白 path 审计枚两处，并排读数那枚 `level=L2` 是 R2 白送的"。至于派单点名要补的那发——**只改冻结文本（`PLAN.md:2536`／`SPEC-07:43`）的档位 L2→L1**：本程以 `git grep 'ReadFile|os.Open|go:embed|ParseFS|ReadDir'` 扫 `internal/tools/**_test.go`·`internal/risk/**` 对 `PLAN|SPEC|docs` ⇒ **0 命中**，即**没有任何 Go 码在运行时读契约文本**（码里的 `PLAN.md:NNNN` 全是注释/错误串，不解析文件）。⇒ **改冻结文本档位＝零用例红，实测坐实**。含义：**"档位"这一维与"契约文本"之间的挂钩今天无人钉**——测试硬编的是*码侧*的 L2 期望，契约侧的 L2 一旦被人脱本票改动（改回 L1），整套工具用例全绿而契约与行为已相反。这一维**是本格的洞、不是实现程的洞**（被验程按 `>`① 只做行为、无权也没动文档）。处方见 §7。

- **③ 空白 path ⇒ 审计行 `in_allowlist_scope=true` 该不该记账**：本程独立判＝**是，该记一笔**。复现：overlay 那遍 `TestFSEditRefusesAPathC26CannotCanonicalize` 的审计行逐字为 `...tool=fs.edit risk=L1 decision=allow outcome=error rules_hit=[R1] in_allowlist_scope=true...`（基线里则 `risk=L2`，但 `in_allowlist_scope=true` 两边都在）。⇒ 一次**根本没被法官判过 path**（空白串在 `bridge.go` 的 `stringValues`/`pathArgs` 就被丢）的拒绝，审计字段却写 `in_allowlist_scope=true`＝把"没东西可判"渲染成"判成在范围内"，**会误导取证人**。它**不归 162 实现程**（该形状在 `bridge` 审计行写侧与 `internal/risk` 侧，162 射程不含），**也不该在 162 里另造第二把尺**（AC#4b 原话禁止为同源事项造双尺；真尺已在 `internal/risk/assessor_test.go`＋`bridge_junction_windows_test.go`）。⇒ **处方＝立一枚独立小票（或并进 `internal/risk` 审计字段语义那张已有票）：判据＝"当一次调用的 `pathArgs` 未向法官提交任何可判路径时，其审计行不得写 `in_allowlist_scope=true`"**。最小闭合集合（可执行判据）：桥层加一枚用例，喂空白/缺失 path 的 `fs.edit`，断言审计行满足 `in_allowlist_scope!=true`（改为 `false` 或专门的"未判定"标记）；摘掉 `bridge.go` 里把"无路径"落到 `true` 的那一处赋值 ⇒ 该用例须红。162 本票内**不动**，勾不受影响。

- **④ AC#6 凭证重走**：见上 §8(片段) 与门禁节——本程用**自有筛法**（逐枚 `--name-only`＋窗口 `--numstat` 分别点名契约轴各路径）而非复用被验程 grep；`f1b99a70→fd8201e7` 全量 name-only 里 `probes/154/gate-clauses.sh`（135/35）与 `probes/171/**` 均出自另一程 `171-r2`（`7f78a163`／`290aa63e`），非 162-r4、非契约轴文件；`9c00bb33` 的 `8/4` 逐行验＝换行重排。结论同：**契约轴零越界**。

## 6. §3 单点回退表（5 枚测试逐枚"撤哪处／哪条不响"）

| 用例 | 撤掉哪一处 ⇒ 该条不响（变绿/失去牙齿） | 本程凭据 |
|---|---|---|
| `…AndFSWriteShareTheOutOfScopeVerdict` | 撤 `fs_edit.go` 的越界拒绝整链（不再走 bridge/C19）才会失去意义；**但把声明档位 `Declared` 改 L1 它照绿** ⇒ 它的 `level=L2` 断言对"档位"不承重、只对"越界 R2"承重 | overlay 变异：此枚 PASS |
| `…RoutesToApprovalNotTheL1Window` | 撤 `fs_edit.go:390 Declared:risk.L2`→L1 ⇒ 范围内 edit 落 L1 窗、此枚红；撤 R2/审批通道同理 | overlay 变异：此枚 FAIL(261) |
| `…OutOfScopeIsJudgedFromTheDeclaredPathNotTheToolWord` | 撤 `PathParams:["path"]`（被验程 M-1）⇒ 法官看不到 path、`rules_hit` 掉 R2、此枚红；撤 Declared 档位**不响**（越界仍 L2） | M-1 由被验程量、本程 overlay 证档位不承重 |
| `…RefusesAPathC26CannotCanonicalize` | 撤 `fs_edit.go:110` 的 TrimSpace 硬拒 ⇒ 空白 path 不再被工具拒、此枚红；撤 `Declared` L2→L1 ⇒ 审计行 `risk=L2` 子串翻红（此枚对档位承重） | overlay 变异：此枚 FAIL(388) |
| `…WriteReclaimsADeadWritersOrphan`（清扫枚） | 撤 `fs_write.go:288 reclaimStaging` ⇒ 孤儿不被带走、此枚红（被验程 M-2 已量，红在孤儿断言 `:60`）；本程 SKIP=0 证其在 Windows 真跑 | 名册多 5 之第五枚；M-2 在被验程 logs |

⇒ 每枚都答得出"撤哪处不响"，**无装饰枚**。

## 7. 最小可见单位（本程自量，附方法）

**方法**：数每枚 AC#5 用例里 `b.Execute(` 的次数，并看断言读的是"调用"还是别的。实测（`git show fd8201e7:…gate_r4_test.go` awk 计数）：并排枚 Execute=2、路由枚 Execute=2、判定枚 Execute=1、空白 path 枚 Execute=1。断言全部落在**每次调用产出的一份 Decision / 一条审计行 / 一次盘上效果**上（`g.approvalDecision()` 取该次审批调用的判决、`r4AuditLine` 按 `tool=` 取匹配审计行末条、盘上字节断言读该次调用的后果）。⇒ **AC#5 的最小可见单位＝一枚调用**（一次 `bridge.Execute` 走完 C19 判定），**不是"一个文件"，也不是"一份审计行"**（审计行只是那枚调用判决的渲染，且同一枚调用还同时被盘上字节/通道计数断言覆盖）。与 162-v2 对 r3 量出的口径一致。

## 8. 门禁四数＋名册差集（本程现跑，逐包 scope，未跑在飞的 `probes/154`·`probes/161` 两把尺）

```
go test -count=1 -v ./internal/tools/                       → rc=0
  === RUN(incl 子测) = 155 ; 顶层 PASS=106 FAIL=0 SKIP=0 panic=0
  名册两向 comm：改前 f1b99a70 声明 101 枚 vs 本遍 106 枚 ⇒ 缺 0 / 多 5
    多出的 5 枚逐一点名：…ShareTheOutOfScopeVerdict / …OutOfScopeIsJudgedFromTheDeclaredPath… /
    …RefusesAPathC26CannotCanonicalize / …RoutesToApprovalNotTheL1Window / …WriteReclaimsADeadWritersOrphan
sh scripts/d22scan.sh                                       → rc=0（clean，无 D22 ban 违规）
tools/d22scan/allowlist.txt                                 → 7 行，grep -c 'fs.edit'=0（未新增豁免）
bash tools/d22scan/runtests.sh -C tools/d22scan ./...       → rc=0，PASS=34 FAIL=0 SKIP=0 RUN=76 '[no tests to run]'=0
gofumpt v0.12.0 (go1.27.1) -l(已跟踪 .go)                   → 甲形 0 行；已跟踪 .go=570 枚
```

SKIP=0 ⇒ Windows 清扫枚 `TestFSEditWriteReclaimsADeadWritersOrphan` 本机**真跑且过**，非靠整包绿充数。原始读数落 `$TMPDIR` 丢弃树 `logs/`（`tools-v3.txt`·`d22scan.txt`·`runtests.txt`·`gofumpt.txt`·`logs_m1_tier.txt`；本程零删除命令，件只建不删）。

## 9. 本程自身越界核查（`git -c core.quotePath=false show --name-only` 逐枚）

本程三枚提交，逐枚 `show --name-only HEAD`（原样贴在各次提交的取数里）写面**只两路**：

| commit | 动到 |
|---|---|
| `7a0d7e4f` | `docs/evidence/s1/162-risk-tier-accept-v3.md`（本表） |
| `6aef4164` | `docs/evidence/s1/162-risk-tier-accept-v3.md` |
| `416cd523` | `.scratch/wisp/probes/162/v3/m1-tier-floor-overlay.sh` |

⚠ `git diff --name-only 7a0d7e4f~1 416cd523`（**跨程范围**）会多显 `171-three-ruler-holes…` 票面、`probes/171/r2/logs/*`、`docs/evidence/s1/171-…-r2.md`——**那些是同树另一程 `171-r2` 在我三枚之间落地的交错提交**（tip 因其在飞而前移，派单 step-0 已预告），非本程所写、本程未动未提未评。判据＝**只看本程三枚各自的 `show --name-only`**，不只看范围端点。`git status --porcelain -- docs/evidence/s1/162-risk-tier-accept-v3.md .scratch/wisp/probes/162/v3/` 现量＝**空**（本程两路全提交）。盘上 `internal/tools/fs_edit.go` 与锚点 `fd8201e7` blob **逐字相同**（overlay 变异只在 `$TMPDIR`，未回写）。

## 10. 被拒／没成功的调用（发生在取数前/后）

- 被权限系统拒绝：**0 次**。
- 本程自伤：**1 次**——首个 commit（`7a0d7e4f`）之前我在表里把锚点误抄成 `fd82107`（正字 `fd8201e7`）并连带抄错一次 tip；`git show fd82107` 报 `invalid object name` 当场暴露，取数主链（各 `git show/grep/diff`）此前一直用正确的 `fd8201e7`，**错误 sha 未进任何判语**；已用后续 Edit 全表纠正（12 处 sha＋1 处 tip），并**追加**修正 commit（不 `--amend`）。这一自伤发生在取数**之前**。

## 11. 有没有跑过删除命令＝**0**

全程未跑 `rm`/`del`/`Remove-Item`/`git clean`/`git checkout .`/`restore`；overlay 变异副本与 logs 只建不删，落 `$TMPDIR` 仓外丢弃树（`mktemp -d`，非 `internal/**`、非仓内）。盘上 `internal/tools/fs_edit.go` 复核仍 `Declared: risk.L2`（grep=1）⇒ overlay 未回写。

## 12. 伪授权两栏（真通知回显／判为注入分开）

- **真通知回显数：0**——本程工具输出里未出现"编排者说／系统通知／请确认真／不必再跑测试直接给结论"类插入语。
- **判为注入数：0**。若出现，按三条判据（路径真不真／内容越权否／盘上可核否）处理＝**只登记不服从**。派单与票面里所有行号（`PLAN.md:2534/2536`、`bridge.go:1004/1020` 等）本程一律现量、未照抄。

## 13. 凭据值零抄录

本程未读、未贴、未提交任何 API 密钥／令牌／DPAPI 密文；出现字符串只有测试自造夹具路径与 `"ALPHA"/"outside"` 类内容；`d22scan.sh` rc=0 含明文密钥扫描项，本程无需绕禁。

## 14. `next=`（AC#5 终判见 §3 上方裁决行；本段给入账处方与最小闭合）

- **AC#5 终判〔成立〕**（凭两条读数）：① 越界经真 C26/C19 走 R2 被拒、盘上字节逐字未变、无暂存件、回执 `user_rejected`——**撤 R2 链即红、与档位无关**，硬凭据；② 声明 L2 由 in-scope 路由枚＋空白 path 审计枚**双点钉住**（`Declared L2→L1` overlay 使这两枚翻红），非纯字符串相等。**不退回**：AC#5 剩余三条（越界必拒／不新增豁免／不动 allowlist）实测全中。
- **随裁三笔入账（都不勾死 AC#5）**：
  - **(a) §2①**：契约文字"R8 覆盖已存在→L2"与机器实现"R1 声明下界→L2"分歧（`fs.edit` 不命中 R8）⇒ **立账**，非本票缺陷。
  - **(b) §2② 本格的洞**：**改冻结文本档位 L2→L1 ⇒ 零用例红**（无 Go 码运行时读文档，实测）⇒ "档位"与"契约文本"的挂钩无人钉。**最小闭合集合**：新增一枚契约↔码一致性尺——从 `docs/PLAN.md` D34 表解析 `fs.edit` 行的档位，与 `risk` 评估器对**范围内已存在文件**的 `fs.edit` 调用所得 `Decision.Level` 比对；**摘掉 `fs_edit.go:390 Declared` 那一味 ⇒ 该尺须红**（今天只 in-scope 路由枚隐式承红，未与文档对账）。此尺落点归**契约轴**（编排者的 `PLAN.md`/`SPEC-07` 那两行）或一枚新票，**不在 162 实现程射程内**。
  - **(c) §2③**：空白 path 的审计行写 `in_allowlist_scope=true` 是误导取证人的假字段 ⇒ **立独立小票／并进 `internal/risk` 审计字段语义已有票**；闭合判据见 §4③。
- 若编排者要把 AC#5 从〔成立〕降到〔退回附条件入账〕，唯一触发条件是认定"档位＝契约文本 L2"这半条属 AC#5 本体而非 AC#6；本程判它属**契约轴↔码一致性**（§2②(b)），故 AC#5 〔成立〕、该洞以 (b) 入账处理。

