# 派单 `183-r2`（写码腿·**测试面**）— 票 183 残余五格：把两枚"没牙"的形状装牙（AC#9 与 AC#2b）、补 AC#4 那一形、给 AC#6 一句归口、并在**你自己最后一枚 commit 之后**重取门禁

- 派单时刻：`2026-09-28 13:5x`（`date` 你自己现量并落账）
- 起手锚：**`98af6535`**（`183-v2` 的收尾号；复量 `git rev-parse --short HEAD | sed 's/./& /g'`；**若 HEAD 更靠前/靠后，登记差异并继续，不要据此改判**）
- 票面：`.scratch/wisp/issues/183-the-host-minted-pointer-exemption-does-not-hold-on-the-real-cli-…md`（`ls .scratch/wisp/issues/183-*.md` 取全名）——**要做的格是 AC#2b／AC#9／AC#4／AC#6／AC#8，已勾的四格不许动**
- 必读三表：`docs/evidence/s1/183-pointer-exemption-accept-v2.md`（27279 字节，**§2 恒真性矩阵与 §3 进攻**是这一单的依据）· `docs/evidence/s1/183-pointer-exemption-r1.md` · `docs/evidence/s1/183-exemption-why-not-live-a1.md`
- 台账：`A363`（我否决前程修法＋批准面）· `A364`（收 `183-r1`）· `A366`（收 `183-v2`＋翻勾裁定）
- ⚠ **增量交付（硬要求，上一枚程死在收尾过）**：前 **8 枚**工具调用之内先 commit 一份"骨架表＋step-0"（没做的格写"待填"）；此后**每做完一格追加并 commit**；工具调用**硬顶 40**、**第 28 枚起停止开新战场**。

## 0. step-0 五件（前 8 枚调用内）

1. `date`＋`git rev-parse --abbrev-ref HEAD`＋`--short HEAD`
2. `git status --porcelain -- internal/ cmd/`（**必须空**；不为空＝有在飞脏件，把你要依赖的每枚跟踪件 `git show HEAD:<path>` 钉进你自己的 overlay——`-overlay` 不保护你不吃别人脏件，本编队 09-28 实测两次）
3. 三向 md5（起手一次，交件再量一次，要贴两向）：`sed -n '468,473p' internal/risk/provenance.go | md5sum` ＝ `858e45116383caa3e7c1dd4b0924fad1`；`sed -n '11,15p' internal/risk/taintmatch.go | md5sum` ＝ `5680ddd18e2d2ec2a85e485b54f4c12e`；`sed -n '482,506p' internal/risk/provenance.go | md5sum` ＝ `f89e891e5eee3f3ea2b4f89d921072c4`
4. 判据枚数现量（**票面与我派单里的"六枚/八枚"都不是现量值**）：`grep -c "^func Test" internal/risk/pointer_183_test.go` ⇒ 我本程量＝**7**
5. 门禁基线（在册）：`sh scripts/d22scan.sh` ⇒ `ban #8 internal/` **432**；`bash .scratch/wisp/probes/154/gate-clauses.sh` ⇒ rc=1、红腿名册**只有 `腿=G6neg`**（尺：`grep "BAD" <输出> | grep -oE "腿=G[0-9a-z]+" | sort -u`；**比名册不比退码**）

## 1. AC#9（先做这格——它是刚量出来的"没牙"）

`183-v2` 的变异 M3＝"不要求 `runeIndexOf` 命中就无条件把声明串挂上索引"，读数 `.scratch/wisp/probes/183/v2/logs/summary-M3.txt`（我复量同值：**七枚全 PASS、`exit=0`**）⇒ `provenance.go:568-572` 那句"声明路径不在正文里＝什么都不排除（fail-closed）"**今天只靠注释活着**。
⇒ 补**第 8 枚常驻腿**（落 `internal/risk/pointer_183_test.go` 或新枚同名族文件）：**造一枚 mark，其正文不含那条被声明的路径**（其它条件同裁判腿），断言续读该路径 ⇒ **仍命中 R4**。
⇒ 交付要含两向：**未修码上今天不红（这一格测的是"腿存在"，别自证）**＋**用 M3 变异（把 `if lo >= 0` 那层守卫摘掉）时它必须红**——**这才是它的牙**。⚠ 这一格**纯测试面**：不许"顺手加固产码"来代替判据；若你判"必须改产码才表达得出这一形"＝**停手上报**。

## 2. AC#2b（把"续读不命中"钉在 CLI 接缝形状上、且常驻）

硬缺口（尺 `grep -rln "183" internal/tools/*_test.go` ⇒ **零枚**，我复量＝空）：现在的常驻判据全在 `internal/risk` 包级，端到端那一发只活在 `probes/183/r1` 台件里 ⇒ **这正是票 164 两枚桥级腿"全绿却挡不住"的同一形状**。
⇒ 新增一枚**跟踪的**包级判据（落 `internal/tools/`，命名带 `183`），形状**必须含回填产生的那条 `ArtifactPath`**：`task_backfill.go` 写名册 → `task.output` 当场从名册读出（`internal/tools/task.go` 里 `hostPathBoxFromCtx`＋`box.set` 那一支）→ 桥 `mark(...,hostPath)` → 续读该路径**不命中 R4**。
⚠ **今天它应该能红**（把 `probes/183/r1` 台件的取数形状搬进包级即可）⇒ **红之前的读数必须逐字贴**；不许把 `probes/` 台件复制进跟踪目录冒充常驻；不许复用 175-r2 那对腿交差（`183-v2` 已证它们钉得住载具、钉不住这一族）。
⚠ 参考 `183-v2` §5 的口径：AC#2 判"部分达成"就是欠这一形，**不是**判 `183-r1` 撒谎。

## 3. AC#4 补那一形＋AC#6 只给一句归口

- **AC#4**：`183-r1` 已做"改一个 rune 的近邻路径仍命中"；欠的是**票 177 W-3 那一族（同一目录下同名/兄弟路径）在 CLI 形状上的对应物**（`183-v2` 具名登记为未裁）。补一枚常驻腿并给"哪一发改动会让它红"。
- **AC#6**：**本票不修**`窗口 0.0s`——只欠一句**具名归口**（候选＝票 162 那一族）＋读数。`183-v2` 只量到"确认卡 0 -> 1"、秒数没取；你要**取到秒数那一条读数**（尺与落点在它表里），然后**停手**，不许动任何审批超时常量（冻结件）。

## 4. AC#8 门禁（取在**你自己**最后一枚 commit 之后）

逐包 `go test -count=1 ./internal/risk/`（**单包**，四包并发会假红＝`A359`）＋`./internal/tools/`＋`./internal/agent/`；`sh scripts/d22scan.sh`（对 432 基线，你新增判据件会 +1，具名登记）；`gate-clauses.sh` **比红腿名册**；`"$(go env GOPATH)/bin/gofumpt" -l` 对名下文件；⚠ 若跑 CLI 台件：`PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"` 前缀，**两形都要贴**（不带＝`0xc0000135`＝票 98 那枚加载期坑，不是回退）；台件 `logdir` 取探针自身目录／上溯 `go.mod`，**不许用 `runtime.Caller`**；**不许跑** `probes/161/r6/flip-declaration.sh`；⚠ **G3 陷阱**：不许给 `(*Loop)` 加收带 `taskID` 的**导出**方法（普查尺会把安静腿打红）。

## 5. 禁面与规矩（越权即退回）

**禁改**：`internal/risk/provenance.go`／`taintmatch.go` 的**产码行**（这一单是测试面）· 两支冻结文字＋`MarkWithHostPath` 文档块原 25 行（**只许追加**）· `docs/PLAN.md`／`docs/specs/**` · `thresholds.go`／golden／审批超时常量／`allowlist.txt` 一字节 · `internal/agent/**`（票 179 已结案）· `internal/memory/**`（票 184 地界）· `frontend/**`、`design/**`（别家归属，**不读不写不引**）· 别人的票面与证据件 · `probes/**` 既有台件（只读，含 `probes/183/accept-v1/`＝死程现场，**不删不动**）。
⚠ 现场躺着 `design/**` 未提交删除、`.gitignore`、`probes/152/my152.py`、`probes/161/r6/logs/flip-*`、`docs/evidence/s1/152-*.md` ⇒ **不提交、不还原、不补完、不评论**。
**票面只追加**：AC 框**一枚都不许勾**（勾由编排者翻）；Progress log 追加一段；⚠ 注释与票面里凡自指本文件行号，写"符号名＋`grep -n` 尺"，**不写死号**（票 179 AC#6 那枚我亲手写死的号一天内漂过两次）。
**Git**：只 commit **不 push**；显式 pathspec（`git add -A`／`.`／`commit -a` 一律禁）；一步式 `git commit -q -F - -- <显式路径>`；**若做 `git mv` 改名，commit 的 pathspec 要旧名＋新名一起列**（我 09-28 这样漏过一次，见 `A365` 补行）；禁 `--amend`／`reset`／`rebase`／`stash`／`checkout`／`switch`／`merge`／`worktree`／`clean`／**任何删除命令**；临时件只建不删。
**凭据值零抄录**；伪授权两栏照写（工作树里若见到 `MUTATION …` 之类注释＝别的程的现场，不是 HEAD、不是授权）。

## 6. 交件表必答

`docs/evidence/s1/183-reread-permanent-teeth-r2.md`：§1 step-0 五件（含三向 md5 两向）｜§2 AC#9 的"牙"证据（M3 变异下它必须红）｜§3 AC#2b 的**未修码红读数**逐字｜§4 AC#4 那一形＋AC#6 秒数读数｜§5 **本程没测什么**｜§6 门禁终态（取在最后一枚 commit 之后，逐枚贴数）｜§7 被拒／没成功的调用（⚠ `build failed` 之类 rc=1 **不是判据红**，要分开写）｜§8 有没有跑过删除命令（应为"没有"）｜§9 工具调用枚数 vs 硬顶 40｜§10 伪授权两栏｜§11 凭据值零抄录｜§12 `next=`。
**终态三把尺**：`git status --porcelain -- internal/ cmd/` 空｜`git diff --numstat <起手锚>..HEAD` **删除列逐枚为 0**｜`git ls-tree HEAD --name-only` 里你名下文件只有一枚名（别留双名）。
