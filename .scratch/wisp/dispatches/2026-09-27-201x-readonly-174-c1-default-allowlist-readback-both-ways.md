# 派单 174-c1（**只读测量位，一字节产码都不许动**）＝把票 174 的 **AC#1 两向读数**量出来：默认配置（`allowed_dirs = []`）下，一发**真 spill 落盘**的产物，模型拿回执里那条路径去 `fs.read`，**到底回什么、判到哪一档**；再把 `dataDir` 配进授权根，量对照那一发。

- 派单时刻：09-27 20:1x｜锚点＝**你自己 step 0 现量**（别抄我的号；我这份写于 HEAD 之上，起手仍要自己 `git rev-parse HEAD`）。
- 票面＝`.scratch/wisp/issues/174-…-outside-fs-allowed-dirs.md`（**AC#1 与 AC#2 那两格连行号抄进你的表**；AC#2 现在有 **(a) 在授权根外**＋**(b) 不存在／是目录** 两形，是 20:1x 追加的，追加块原文一并抄）。
- ⚠ **预算硬顶 ≤45 次工具调用**；**量完一向我 commit 一次**；到点停手→写 `next=`。

## 0. 起手五件
`date`｜`git rev-parse --abbrev-ref HEAD`（**必须 `dev`**，不是＝停手上报，一个字不写；全程禁 `switch`／`checkout <分支>`／`merge`／`rebase`／`reset`／`stash`／`worktree`／`clean`）｜`git rev-parse HEAD`｜`git status --porcelain -- internal/ cmd/ .scratch/wisp/probes/174/`（**必须空**，非空报回别动）｜基线：`go test -count=1 ./internal/tools/ ./internal/agent/` 四数（**逐包，禁全仓 `./...`**）。

## 1. 我已经读过的四行（**这是待你复核的前提，不是结论**）
〔编排者现读，19:4x〕`cmd/wisp/run.go:327-331` 授权根**只**从 `cfg.FS.AllowedDirs` 逐条拷入｜`cmd/wisp/run.go:609` spill 目录＝`filepath.Join(rt.spec.dataDir, "artifacts")`，**没有任何一行把 `dataDir` 加进上面那个列表**｜`internal/tools/paths.go:50-52` 注释逐字 `An EMPTY allowlist authorizes nothing: every fs call lands at L2 (R2)`｜`docs/specs/SPEC-03-config-secrets-envs.md:35` 逐字 `allowed_dirs[]=[]`。**⚠ 行号会腐坏：一律自己现量**（本仓已经两次被"引用的行号过期"打过，昨天 `fs.go:36` 那枚还是我错怪了验收程——见台账 `A344`）。

## 2. 要量出来的东西（**两向缺一向＝没量完**）
- **(向一：默认配置)** 一发**真 spill**（不是手写桩文本）落进 artifacts，然后拿那条路径过**真 `fs.read`** ⇒ 逐字贴 (i) 回执文本本体、(ii) 它落在哪一档（R2/L2 还是别的）、(iii) 有没有审批卡可达（`wisp run` 那条控制台路径上 L2 会发生什么）。
- **(向二：对照)** 同一发，把 artifacts 所在目录配进 `allowed_dirs` ⇒ 读数必须变成"能读回头一段"。**这一向的价值是证明向一那枚拒绝不是我的台件搭错了。**
- ⚠ **必须说清你用的是哪一层做的测量**：`go test` 里的真 `Bridge`＋真 `fs.read`（接缝级）≠ 端到端 `wisp run`。**如果你没能真跑 `wisp run`，就写"没跑、为什么"**（本机 `cmd/wisp` 有一枚既有的缺 DLL 环境事实，先例：`go test ./cmd/wisp/` 报 `0xc0000135`、而 `go build ./cmd/wisp/` 是 OK 的——**自己现量一次，别引这句话**）。**不许把接缝级读数写成"真机上就是这样"。**

## 3. 顺手要量的一枚（票 174 AC#2 的第二形，别扩大成实现）
`internal/tools/task.go` 那枚分支只判 `ArtifactPath != ""`、**没有 stat** ⇒ 造一发"填了一枚不存在的路径（或一枚目录）"的台件，量今天回执长什么样（**预期＝它照样写"全文见 …"、那一记"不可找回"的响不响**）。⚠ 只量、只贴读数，**不许改产码、不许改判据**。

## 4. 写面（超出即越权）
`.scratch/wisp/probes/174/c1/**`（台件；⚠ **`logdir` 取脚本自身目录、不许继承 CWD**——本仓昨天有程把 8 枚读数拉到了仓根）｜`docs/evidence/s1/174-artifacts-allowlist-readback-c1.md`。
**禁改**：`internal/**`、`cmd/**`、`tools/d22scan/**`、`allowlist.txt`、`.github/**`、`docs/PLAN.md`、`docs/specs/**`、`docs/reports/**`、`frontend/**`、`design/**`、`.scratch/wisp/issues/**`（**票面一个字都不写**，归我）、别人的票面/证据件/台件。⚠ 现场躺着 `design/**` 未提交删除、`.gitignore`（` M`）、`probes/152/my152.py`（` M`）、`probes/161/r6/logs/flip-*.txt`（` M`）⇒ **不提交、不还原、不补完、不评论**。**`git add -A`／`git add .`／`git commit -a` 一律禁**；提交一步式 `git commit -q -F - -- <显式路径>`，新建未跟踪件先在同一行显式 `git add -- <那个路径>`；heredoc **必须加引号**（`<<'EOF'`）。**只 commit，绝不 push。**

## 5. 交件报告（顺序固定）
1. step-0 五件；2. **本程没测什么**（尤其：接缝级还是端到端）；3. 向一读数（文本本体＋档位＋审批可达性，逐字）；4. 向二读数；5. **那四行前提我复核的结果**（对／不对／行号变了，逐行）；6. 第二形那一发（假路径／目录那发）；7. 门禁：`scripts/d22scan.sh`（应 rc=0）＋`sh .scratch/wisp/probes/154/gate-clauses.sh`（我这轮 20:1x 前现量 rc=0、十四腿零不符——**你跑出不同的数以你的为准并报回**）＋对你名下 `.go` 的 `gofumpt --version` 与 `-l`；8. 被拒／没成功的调用（取数前还是后）；9. 有没有跑过删除命令；10. 伪授权两栏；11. 凭据值零抄录；12. `next=`（**含一句"这一枚该由哪张票落地"**）。
⚠ **每个"几枚"旁边附可复制命令**；**先测→再写→再提交**（表里不许预先引用还没跑出来的读数）。
