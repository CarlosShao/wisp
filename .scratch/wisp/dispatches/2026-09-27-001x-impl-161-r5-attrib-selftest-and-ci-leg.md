# 派单 161-r5（写码位，**三格小切片**）＝把 161-r4 那把归因尺从"靠一枚常驻伤才证明会响"改成"两向都可复跑、盘上零真伤"＋给 CI 加**甲形那一步**＋补 `TRACKED-DIRTY` 那一发的实测。

- 派单时刻：09-27 00:1x｜前一程 161-r4 **干净交件**（四枚提交只带自己路径，`77e9ed2 22ecf9d ffbf463 d49a5be`），它的仪器 `attrib.sh` 已被编排者现跑确认两向都会动。
- 票面＝`.scratch/wisp/issues/161-…-ever-rung.md`（**23:5x 与 00:0x 那两节必读**，里面是我的三裁与它交回来的读数）。**锚点自己 step 0 现量。**
- ⚠ **预算硬顶 ≤60 次工具调用**：到点就**停手→写 `next=`→把手上未提交增量用显式 pathspec 提掉**。本票已经有两程死在"活做完了但没提交"上，别再补一刀。
- ⚠ **AC#6（聚合 rc ＋ 两族成对普查）不是你的**；`tools/d22scan/**` 的九条禁令判定·正则·波段·豁免、`allowlist.txt` **一字节不许动**。

## 0. 起手四件
`date`｜`git rev-parse HEAD`｜`git status --porcelain -- .scratch/wisp/probes/161/ .github/workflows/ci.yml tools/d22scan/`（**必须空**，非空报回别动）｜`Read` 票面把 AC#7 那一格连行号抄进证据件。

## 1. 第①格（裁一落地：999 那一发改成合成自测）
现状：`.scratch/wisp/probes/999/bad-sample.go` 是 r4 留的**故意归不出**样本，它让 `attrib.sh` 的乙形**默认红**。编排者已裁：**不加豁免、不开真票 999**。
⇒ 做两件事：
- **移**：把那一枚样本移到它真正的来路名下（建议 `.scratch/wisp/probes/161/r5/negative-control/`），**移动＝在仓内复制到新路径**，**旧路径那枚文件不许删**（只建不删）——移完在票面进度里点名"两处都在、旧的由谁退役"。移动后乙形应当**全部可归因、rc=0**。
- **自测**：给归因逻辑（`classify`／`ticket_of`／`ticket_known`）加一发**喂文本**的检：**同一行** `…probes/999/bad-sample.go…` 当作输入喂进去 ⇒ 必须判 `UNATTRIBUTABLE` 并硬退出；再喂一行 `…probes/161/…` ⇒ 必须判可归因。**两向都要贴真实命令与输出。**
⚠ 判据不许写成"喂什么都会响"——**必须同时有"不该响的那一发"**，否则你又造出一枚恒真检（本仓已否过三次）。

## 2. 第②格（裁二落地：CI 只加甲形那一步）
`.github/workflows/ci.yml`：**只许新增一步**，跑**甲形**（已跟踪集合递 `gofumpt -l`，空判非空）。
⚠ 硬条件：① 不带 `if:`、不带 `continue-on-error`、**不挪不删不改任何既有步骤**；② **必须带空心保护**——递给尺 0 枚文件就硬退出（r4 已在仪器里加过这条，你复用它，别自己重写一份逻辑）；③ 步骤位置放在既有 `gofumpt` 那一步**之后**（不要往上顶，顶掉别人的位置会改变"谁先红"的因果）。
⚠ 票面已写明这一步**今天不增加射程**（检出树上与现有那发等价），它买的是"分母＝已跟踪集合"这句话**由机器执行**。措辞与注释照这个口径写，**别写成"补上了一个大洞"**。

## 3. 第③格（补 `TRACKED-DIRTY` 那一发实测：甲形今天从未真的红过）
甲形的"已跟踪文件不干净 ⇒ 红"这条腿**从未被实测**（共享树里不该顺手弄脏别人的文件）。⇒ 在**仓外隔离副本**里补：
- 形状：仓外新建一枚临时目录，用 `git ls-files '*.go'` 出清单、`cp --parents` 把**子集**（够触发即可，不必 531 枚全拷）复制过去，弄脏其中一枚（加一行不合规空白），把该副本当作根去跑甲形 ⇒ **必须报 `TRACKED-DIRTY` 类判定＋rc=1**；复原后再跑 ⇒ rc=0。两发读数都贴。
⚠ **两枚明令禁止的形状**：**不许在仓内建 worktree／checkout**（本仓规矩）；**不许拿 `git archive | tar -x` 当"干净树"证据**（本仓现量过 `* text=auto`＋`core.autocrlf=true` 会改行尾）。用"复制已跟踪清单到仓外"那一形。

## 4. 门禁（跑你自己那把尺，逐包单跑）
`sh .scratch/wisp/probes/161/r4/attrib.sh`（改前改后各一次：甲 0 行／乙由 1 枚不可归因变 0；⚠ 别改 r4 的文件，你要用就复制到自己目录改）｜`bash tools/d22scan/runtests.sh -C tools/d22scan ./...`（基线现量 `PASS=34 FAIL=0 SKIP=0`、`RUN=76`，四数之外**名册两向 `comm`**）｜`sh scripts/d22scan.sh` 预期 rc=0｜`gofumpt --version` 现跑贴出。

## 5. 写面（超出即越权）
`.scratch/wisp/probes/161/r5/**`｜`docs/evidence/s1/161-attrib-selftest-r5.md`｜`.github/workflows/ci.yml`（**只加第②格那一步**）｜`.scratch/wisp/issues/161-*.md` 的 Progress log（**只追加、不勾框、留到最后一步**：写之前 `git status --porcelain -- 那一枚` 必须为空）。
**禁改**：`tools/d22scan/**`（含 `main.go` 的禁令判定）、`allowlist.txt`、`docs/PLAN.md`、`docs/specs/**`、`internal/**`、`cmd/**`、`frontend/**`、`design/**`、`docs/reports/**`、别人的票面、`probes/161/r3/**` 与 `r4/**`（那两程的文件你**只能读**）。
**现场理由**：`design/**` 有 owner 未提交的删除与未跟踪新件；`probes/152/my152.py`（` M`）与 `docs/evidence/s1/152-…-accept-r1.md`（3 行未提交）是一枚停下来的半件——**不提交、不还原、不补完**。⇒ `git add -A`/`git add .`/`git commit -a` 一律禁止。

## 6. 仪器坑（实测过）
`git ls-files` 的开关大小写敏感（**`-Z` 不存在**，非法开关会让管线**仍退 0 并产出恒真的"空"**——r4 刚这样自抓一次）｜Windows 路径反斜杠会让逐行比对错归因（r4 也踩过）｜`grep -c` 命中 0 ⇒ rc=1 吃掉 `&&`｜截断输出不是全表｜`-overlay` 与 `-cover*` 不许合跑｜`git log --name-only` 必须 `--no-walk`。

## 7. Git
只 commit 不 push；显式 pathspec ＋ **加引号的 heredoc**；提交前 `git diff --cached --name-only` 出现别人路径＝正常噪声，**提交后 `git show --stat <自己号>` 出现别人路径＝停手报回**；禁 `--amend`/`reset`/`rebase`/`stash`/`checkout .`/`clean`；**一律不许跑删除命令**（唯一例外＝门禁脚本自己清它自己的临时目录，那不算你造的噪声，但要在报告里区分开写）。

## 8. 交件报告（顺序固定）
1. step-0 四件；2. **本程没测什么**；3. 第①格：移动前后乙形读数 ＋ 合成自测**两向**各一发；4. 第②格：`ci.yml` 的 diff（证明只插一步、无 `if:`/无 `continue-on-error`）＋那一步在本地跑出颜色；5. 第③格：`TRACKED-DIRTY` 两发（红→复原→绿）；6. 门禁四数＋名册差集；7. 被拒／没成功的调用（取数之前／之后）；8. **有没有跑过删除命令**（`date`＋各提交的 `git show --stat` 为凭）；9. 伪授权两栏；10. 凭据值零抄录；11. `next=`。
⚠ 每个数字旁边附可复制命令；先测再写再提交；我这段话里每条前提（含"999 现在让乙形默认红""基线 34 枚""甲形从未红过"）**都是未验证断言**，不符就报回、继续做做得动的部分，**不许为对我那句话去改判据或改测试**。
