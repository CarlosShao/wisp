# census.md — 票 269 的只读普查（腿 `269-a1`，三形各四问，【禁】 本腿不裁形）

交付件＝本文件。写面＝仅这一枚。工单票面 AC 框与 `docs/reports/pending-and-issues.md` 一字未动。

## §0 起手锚（我自己现读，不抄票面）

- 起手现量时刻：`2026-10-05 11:52:00 +0800`（`date`）。落笔时刻：`2026-10-05 12:1x +0800`（同一条命令内两次 `date` 之差＝本腿全程）。
- 起手 HEAD：`05680e44`（`git log -1 --format=%h`）；分支 `dev`。**写到一半 HEAD 漂到 `50e29c`→`50e0299c`**（别的腿入库），故 §1 每一把都标了读数时刻，并在 `12:0x` 复认 `ci.yml` 那 7 个行号在新 HEAD 上仍成立（`git diff --numstat -- .github/workflows/ci.yml`＝空、该文件最后一笔＝`519e1ade` 10-04 09:42）。
- 工作树脏度（【禁】 不是我的）：起手 `git status --porcelain` 前 40 行即含 `internal/agent/approval/approval.go`／`queue.go`／`ticket242_binding_test.go` 三枚 M＝在飞写腿 `259-r1` 的包；`git diff --name-only HEAD | wc -l` 现量 **32**。
- 【禁】遵守情况：全程零 `go build`／`go vet`／`go test`／`go run`；用过的仪器＝`grep`／`sed`／`awk`／`git ls-files`／`git ls-tree`／`git log`／`git cat-file`／`git archive`（只读导出到 `/tmp/269a1/`，【禁】 未在仓库目录内建 worktree 或 checkout）＋**`gofumpt` 本体**（`D:\work\base\gopath\bin\gofumpt.exe`，`--version` 现读 `v0.12.0 (go1.27.1)`）。派单允许清单里没有 `go env`，所以我**没有真跑** `attrib.sh`（它第 137 行要调 `go env GOPATH`），而是用等价命令复现它那一条管线（见 §4）。
- 取数口径总则：**每一枚数都带"命令＋分母＋时刻"**；行号一律现读到才写（票面与台账给的行号我只当待验断言）。
- 一把通用的量具（下面多处依赖，具名一次）：CI 那台检出的形状＝**tracked 字节、LF 行尾**。本机工作树被别的腿写过的 `.go` 有相当一部分盘上是 **CRLF**（`core.autocrlf=true` 而 `.gitattributes:4` 写 `*.go text eol=lf`，blob 侧 CR＝0），于是**本机直接跑 `gofumpt -l` 会把纯行尾差异读成"要格式化"**。尺：`tr -dc '\r' < 盘上文件 | wc -c` 对 `git cat-file blob HEAD:<file> | tr -dc '\r' | wc -c`，现量 `cmd/wisp/models.go` 334 对 0、`internal/risk/provenance.go` 1115 对 0、`internal/tools/bridge.go` 1280 对 0、`internal/agent/tools.go` 225 对 0、`internal/agent/approval/pending_read.go` 131 对 0、`probes/257/r2/mut/M1/firstrun.go` 185 对 0。⇒ 本文件凡"CI 会读到什么"的读数一律取自 **HEAD 的整树归档**（`git archive HEAD | tar -x` 到 `/tmp/269a1/full`，`find -type f | wc -l`＝**5839**＝`git ls-files | wc -l`＝**5839**，两数相等＝归档完整），本机工作树读数只用于对照并具名标注。

## §1 复跑票面四把现量（逐把：我自己的命令／时刻／读数／一致与否）

### 现量 #1（样本在库＋预存判定）＝**一致**，但票面那枚"53"复跑不出

- 命令与读数（`11:55 +08`）：`git ls-files -- .scratch/wisp/probes/185/c1/mut/` → **3 枚**：`fs.go`／`fs_broken.go`／`taintmatch.go`。⇒ 与票面 现量#1 逐字一致。
- 首入库：`git log --format='%h %ad %an %s' --date=format:'%Y-%m-%d %H:%M' --diff-filter=A -- .scratch/wisp/probes/185/c1/mut/fs_broken.go` → **`4813567e 2026-09-28 14:42 CarlosShao docs(evidence/185-c1): reread-owner census ...`** ⇒ 票面"首入库笔＝`4813567e`（09-28 14:42）"**一致**。
- "不在本批"这一句：票面用的是枚数尺（`21bec8a1..HEAD` 那 53 枚）。我复跑的是**祖先尺**（更硬，且不随 HEAD 漂移）：`git merge-base --is-ancestor 4813567e 21bec8a1` → **rc=0＝是真祖先** ⇒ "预存"成立。枚数尺我读三发：`git log --format='%h' 21bec8a1..HEAD | wc -l`＝**77**（`11:55` 那次）、`git rev-list --count 21bec8a1..HEAD`＝**78**（`12:10`，HEAD 已漂）、`git rev-list --count 21bec8a1..c6cf66e6`＝**49**。
- ⚠ **不一致（一处，具名）**：`53` 这枚数我复跑不出来——`ci-attr-1` 自己 `§0` 第 19 行写的是 `git rev-list --count 21bec8a1..HEAD`＝53 且它起手 HEAD 就是 `c6cf66e6`（同一区间我这发＝**49**）。差 4 枚的原因本腿查不到（可能在两次读之间 HEAD 又走了 4 枚、或它把 merge/其它计数口径混进来）⇒ **不改判、只登记**：不影响"预存"结论（那靠祖先尺）。

### 现量 #2（`ci.yml` 七个绝对行号＋守卫分层）＝**七个行号全一致；分层清单少一枚**

命令（`11:55:01 +08`，`grep -n` 现读）与读数：

| 票面断言 | 我这发读数 | 判定 |
|---|---|---|
| `:168 - name: gofmt (gofumpt)` | `168:      - name: gofmt (gofumpt)` | 一致 |
| `:171 OUT="$(gofumpt -l . tools/d22scan tools/mockllm)"` | `171:          OUT="$(gofumpt -l . tools/d22scan tools/mockllm)"` | 一致 |
| 步名在 `:184`（形 (A) 那步） | `184:      - name: "gofmt (gofumpt) - the tracked set is the denominator (ticket 161 AC#7 form A)"` | 一致 |
| `:204 run: sh .scratch/wisp/probes/161/r5/attrib.sh --tracked-only` | `204:        run: sh .scratch/wisp/probes/161/r5/attrib.sh --tracked-only` | 一致 |
| `:206 - name: go vet (module)` | `206:      - name: go vet (module)` | 一致 |
| `:209 - name: go vet (tools/d22scan module)` | `209:      - name: go vet (tools/d22scan module)` | 一致 |
| `:213 - name: staticcheck` | `213:      - name: staticcheck` | 一致 |

- "`:184`/`:206`/`:209` 三步无 `if:` 守卫"＝**一致**（§3 全表逐枚扫过，这三步 `if=no`）。
- ⚠ **不一致（一处）**：票面列的"都带 `if: ${{ !cancelled() }}`"的七枚（`:263`/`:300`/`:436`/`:474`/`:484`/`:506`/`:542`）**实测是八枚**，漏的是 **`:558`**（`558:        if: ${{ !cancelled() }}`，属 test-windows 的 `:545 PathResolver junction placeholder` 那步）；另有一枚 `:376:        if: always()`（test-core 的 Stop compose）票面没算进"分层"叙述。⇒ 分层结论**只会更强**（带守卫的步＝9 枚，不是 7 枚），票面那句"守卫在这棵树上是分层的，不是全有或全无"成立。

### 现量 #3（两句自陈）＝**实质一致、行号引用错位**

- `:191` 逐字（现读）：`# ticket 161 AC#7 was decided on after `gofumpt -l . tools/... must be EMPTY`` ⇒ 票面引那句**一致**（票面把反引号写成单引号，词面同）。
- ⚠ **不一致（行号，实质不变）**：票面写"`:195-198` 逐字承认新那步 **adds NO REACH TODAY**、并要求两把尺 must agree"。我现读这段注释的落点是：
  - `:186` `# it adds NO REACH TODAY. On a checkout every file is tracked, so the step`
  - `:187-188` `# above and this one read the same bytes and must agree - if that ever stops` / `# being true, the disagreement is the finding. What this step buys is that`
  - `:189-190` `# the sentence "the denominator of the formatting gate is the TRACKED set,` / `# not the working tree" stops living in human memory. That sentence is what`
  - `:194-195` 讲的是"只在人记忆里＝下一个程序格式化'整棵树'时它会悄悄变"；`:196-203` 讲的是"这不是第二份拷贝、它复用同一把尺、并复用那把尺的空心守卫：递给零枚受跟踪 `.go` 就 exit 2，CI 读成红、永不读成绿"。
  ⇒ 票面真正要引的是 **`:185-188`**（NO REACH／must agree／disagreement is the finding）＋**`:189-193`**（分母＝tracked 集合）。**结论不变、行号要改**；本文件下面一律按现读行号引，并在 §8 具名顶回。

### 现量 #4（词面雷：`grep -rln fs_broken`）＝**一致且远超票面**

- 命令（`12:00:00 +08`）：`grep -rln --exclude-dir=.git fs_broken .` → **60 枚文件**（票面只说"命中至少票 236"，方向一致）。逐枚 tracked 判定：`git ls-files -- <p> | wc -l` → **tracked 37 枚／未入库 23 枚**。
- 三种引用形态的枚数（同一把尺，口径写清）：

| 词面（`grep -rl -F`，排除 `.git`） | 命中文件 | 其中 tracked |
|---|---|---|
| `fs_broken`（任何形） | 60 | 37 |
| `mut/fs_broken.go` | 52 | 30 |
| `185/c1/mut/fs_broken.go`（全路径形） | 51 | 29 |
| `fs_broken.go:4:1`（路径＋行列钉在一起） | 40 | 19 |

- 逐枚分类（哪些是注释性引用、哪些是脚本真会去 open）＝**§6**，本条只报数。票面 现量#4 说"别的证据件是否逐字引这条路径，归 AC#0 穷举"——穷举在此：60 枚逐枚读过判定行。
- ⚠ **票面没提到的一枚**：`.scratch/wisp/probes/185/c1/overlay-proof.json` 是**机器消费件**（`go build -overlay` 的输入），它逐字节把 Windows 绝对路径 `D:/work/workspace/projects plans/Wisp/.scratch/wisp/probes/185/c1/mut/fs_broken.go` 写死为 `Replace` 的值 ⇒ 改名那一形唯一"会真去 open"的连带面就是它（尺：`grep -n "fs_broken\|185/c1/mut" .scratch/wisp/probes/185/c1/*.json` 现读三枚 overlay JSON，只有 `overlay-proof.json` 指 `fs_broken.go`，另两枚指 `fs.go`／`taintmatch.go`）。

### 现量 #5（票面没跑的第五把：门的分母里到底有几枚脏样）＝**票面题面被这一把推翻**

尺（CI 形状，`12:07:00`—`12:07:02 +08`）：`git archive HEAD | tar -x -C /tmp/269a1/full` → `git ls-files -z '*.go'`（**921 枚**，`12:0x` 现读 `wc -l`）逐条前缀归档路径 → `xargs -0 gofumpt -l`。

- 读数：**19 枚文件**（stdout 18 行＋stderr 1 行解析诊断，去重后 19 枚不同文件），`xargs` 聚合退码 **123**。逐枚名册见 §5。
- 其中 **1 枚不是样本**：`tools/d22scan/selftest.go`（`gofumpt -d` 现读要改的就是 map 字面量键的对齐两行；入库笔 `git log -1 -- tools/d22scan/selftest.go`＝`5e8748b3` 10-03 21:07）。
- 同一把尺在 `c6cf66e6`（`ci-attr-1` 那枚 run 的 headSha）上复跑：`git ls-tree -r --name-only c6cf66e6 | grep '\.go$'`＝**904 枚**，同法 `xargs -0 gofumpt -l` → **同 19 枚**（`diff` 空）。⇒ **那枚 run 的 lint::gofmt 步今天同时在红于 19 枚文件**，而它的日志只露出 1 枚（原因见 §2 乙④ 与 §8 第 2 条：`$OUT` 被命令替换吞掉＋`bash -e` 在退 2 时当场掐死该步，`echo "$OUT"` 那两行永远走不到）。
- 与台账的差：`docs/reports/pending-and-issues.md:9610`（10-01 那次复跑）读的是"分母 726 枚 tracked `.go`、stdout 7＋stderr 1＝**8 枚脏样**"⇒ 到今天 **8→19**，多出的 11 枚是 10-01 之后陆续入库的（含那枚非样本）。这一条是"甲法确实能让门变绿"那句旧结论过期的直接凭据（§8 第 5 条）。

## §2 三形各答四问（【禁】 不选形、不裁形；裁形归编排者）

### 甲＝给 `ci.yml:168` 那一步换分母（只喂"tracked 集合里非样本"的 `.go`）

**①要动哪几枚文件（逐枚路径）**

1. `.github/workflows/ci.yml` —— 只改 `:169-176` 那一步的 `run:` 体（把 `gofumpt -l . tools/d22scan tools/mockllm` 换成"tracked 全集剔掉具名样本名册"）。⚠ **这一步是"既有步骤"，不是新增步骤**：票 171 `:31` 与票面写面句（`269:32-33`）对"动既有步骤"给的是不同的界，见 ② 末与 §7。
2. 新增 `scripts/gofumpt-tracked-gate.sh`（分母计算＋名册守卫＋自报读数；本仓"门自己数自己的分母"的既有形状＝`scripts/check-path-length-budget.sh:298` 打 `denominator read: %s tracked paths`）。
3. 新增名册件（或把名册内联进上面那枚脚本的数据段，照 `scripts/check-path-length-budget.sh:121`「Exemptions are DATA inside this file, not comments: each roster line assigns」＋`:169-170` 与 `:314`／`:372` 的 `#ROSTER-AWK-BEGIN`／`#ROSTER-AWK-END` 标记定式）。
4. `ci.yml:178-204` 那步（形 (A)）**不动**——它按票 161 的裁定必须继续吃"tracked 全集"，否则 ② 那两句自陈被改动两次。

**"排除怎么写才不会顺手把真违规也排除掉"——具名形状（五条，缺一条就漏）**

- (一)**逐枚全路径、禁目录通配**。名册条目只能是 `.scratch/wisp/probes/<NN>/<...>/x.go` 这种完整相对路径；写成 `probes/**`／`mut/**`／`*.go` 一类模式＝把以后落在同一目录里的真伤一起排掉（本仓同款病的既有定式＝`check-path-length-budget.sh:124-125` 的 A／B 双向守卫）。
- (二)**准入条件＝两条同时成立**：路径必须以 `.scratch/wisp/probes/` 开头，**且** `.scratch/wisp/issues/<NN>-*.md` 真存在。这条直接抄 `attrib.sh:174-193`（`ticket_of`＋`ticket_known`），它的注释逐字理由：「The first half alone would be self-fulfilling (any number passes)」。
- (三)**四形守卫**（照 `check-path-length-budget.sh:120-131` 逐条搬）：A 任何未入册的脏样 ⇒ 红；B 名册里一条今天不再是脏样／不在盘上 ⇒ 红（防过期豁免＝防 phantom exemption）；C `count(名册条数) != count(tracked 脏样条数)` ⇒ 红（连重复键一起挡，原文 `:128-130`）；D 理由为空或全空白 ⇒ 红。**B 与 C 是"只排样本"与"排掉一切"之间唯一的机械差别**，没有 B／C 的名册一定烂成豁免清单。
- (四)**分母读数照常打印**：剔样本之前先把 `git ls-files '*.go' | wc -l`（现量 **921 @`05680e44`／922 @`65142d9a`**，见 §5 末条）打进输出，并保留空心守卫（`attrib.sh:337-338`：递给尺 0 枚 ⇒ `exit 2`，CI 读成红、永不读成绿）。剔的是"被判定的集合"，不是"被数出来的集合"。
- (五)**名册只许覆盖 `.scratch/**`，一枚都不许多覆盖**——⚠ 现量后果：按 §1 现量#5，tracked 分母脏样＝**19 枚**，其中 **18 枚**在 `.scratch/wisp/probes/**`、**1 枚**是 `tools/d22scan/selftest.go`（真违规，非样本）。名册写到 18 枚封顶 ⇒ **`:168` 那一步仍然红**，红的是一枚真违规；把第 19 枚也写进名册＝正是 (一)(三) 要挡的那一形，本腿按字面判它"不许"。而 `tools/d22scan/**` 不在 269 写面（票面 `269:33` 地界＝`ci.yml`＋`scripts/**`＋`.scratch/wisp/probes/**`；票 171 `:35` 还额外要求 `tools/d22scan/**` 零字节）⇒ 修那枚格式要么具名扩面、要么归票 122／新票，**都不是甲自己能顺手做的**。

**②会不会推翻 `:191`／`:195-198`（现读应写作 `:189-193` 与 `:185-188`）那句已裁的"分母＝tracked 集合"＝会，而且是两句（我的待验断言"甲会"成立）**

- `:189-190` 逐字：`# the sentence "the denominator of the formatting gate is the TRACKED set,` ／ `# not the working tree" stops living in human memory. That sentence is what`。甲把 `:168` 的分母变成「**tracked 集合减去具名样本**」——那一句从此对 `:168` 不真。
- `:186-188` 逐字：`# it adds NO REACH TODAY. On a checkout every file is tracked, so the step` ／ `# above and this one read the same bytes and must agree - if that ever stops` ／ `# being true, the disagreement is the finding.` 甲使这两步**按设计必须不同意**（`:168` 剔 18 枚、`:184` 一枚不剔），而那一句把"不同意"定义成"that is the finding"（＝一次要立案的发现）。⇒ 甲不是"用 finding 换绿"，而是**把这个 finding 变成每天必发**——除非同时把 `:184` 的分母也改成剔样本，而那正是票 161「裁二」（`.scratch/wisp/issues/161-*.md:141`：「加，且只加甲形……它买到的是**把"分母＝已跟踪集合"这句话变成机器执行的**，不再靠人的记忆」）明令不许动的东西。
- ⚠ 还要摆一层：票 161 `:124` 逐字写着「⚠ 这**不是放宽阈值**、也不是给 `.scratch` 加豁免——是把判据改成它本来就该量的那个东西」。**甲恰恰就是那句被 161 点名"不是我的修法"的形状（给 `.scratch` 加豁免）**。⇒ 甲不只推翻两行注释，是与 161 AC#7 的裁定文本正面对撞；票面 `269:33` 那句"裁甲则先摆机主"按这一条读是**必需**而非保守。

**③会不会咬到引用那枚样本路径的件＝不咬样本，但新增一枚需要同步维护的路径清单**

- 甲一字不改 `.scratch/wisp/probes/185/**`：`fs_broken.go` 原地、原名、原字节（`git rev-parse HEAD:.scratch/wisp/probes/185/c1/mut/fs_broken.go`＝`ee58385366d8fb6693c12fe0a4367225538f9dd4`，本腿现读）。⇒ §6 那 60 枚引用件全部继续有效，AC#1「那枚坏样本必须还在且还要能被指到」自动满足。
- 代价具名：名册＝这 18 条路径的**第二次落地**。§6 数过"全路径形"引用已有 **29 枚 tracked 件**，名册是第 30 枚，而且是**会被机器照着剔文件**的那一枚——名册过期由 B 守卫判红（这是要买的行为），但"多一处必须同步维护的路径清单"这句话必须写进裁形代价，不能只报收益。

**④修完之后那扇门还剩什么牙（判语，不是偏好）**

- `:168` 的牙：还剩，且**第一次咬在真东西上**——门仍红，红因是 `tools/d22scan/selftest.go` 这枚真违规（§1 现量#5）。但"常红＝没人看"那一格**不因甲而解**：红的形状没变（每天仍红），变的是红的理由。
- `:184` 的牙：**甲碰不到它**（它那步入参写死 tracked 全集，`attrib.sh:341`）。⇒ 甲之后 `:184` 仍红（`exit 2`；机理见 §4）。
- **链头红从 `:168` 挪到 `:184`**，而 `:184` 无 `if:`（§3 表）⇒ `:206`／`:209` **继续采不到读数**。⇒ **甲单独做不解决本票第二半（吞三步）**，只把第一半改好一半，还要推翻两句已裁自陈＋与 161 裁定文本对撞。这是甲的形状与代价，不是本腿的建议。

### 乙＝把那枚样本从"会被 Go 仪器看见"的形状里搬走（改名／移出模块，先例＝票 267 `pristine/models.go`→`models.go.pristine`）

**①要动哪几枚文件**

1. `.scratch/wisp/probes/185/c1/mut/fs_broken.go` → 同目录改后缀（`.go.pristine`／`.go.txt` 一类）。【禁】 不许"搬出 `.scratch`"，更不许删（票面 AC#1＋`issues/README` 规则 8「临时件只建不删」）。先例尺：`git ls-files | grep -c '\.pristine'`＝**17 枚在册**（票面点名的 `267/r2/pristine/models.go.pristine` 是其中一枚；同族＝`248/r1` 3 枚、`248/v1b` 5 枚、`252/v1` 3 枚 `-pristine-v1`、`256/v1` 2 枚、`33/r1` 1 枚、`33/r9` 2 枚）⇒ 先例形状＝**只改后缀、不动字节**。
2. `.scratch/wisp/probes/185/c1/overlay-proof.json` —— **唯一机器消费件**：现读全文是 `{"Replace": {"D:/work/.../internal/tools/fs.go": "D:/work/.../185/c1/mut/fs_broken.go"}}`，值里写死绝对路径。不同批改它，下一次复跑 185 那发"种 X 必响"正控就是 `go build -overlay` 打不开替换目标。
3. 引用面（逐枚见 §6）：**`docs/**` 里 5 枚 tracked 件**（`docs/evidence/s1/185-reread-owner-census-c1.md`、`docs/evidence/s1/197-subagent-entity-r1b.md`、`docs/evidence/s1/200-project-instructions-r2.md`、`docs/reports/HANDOVER.md`、`docs/reports/pending-and-issues.md`）**不在 269 写面**（票面 `269:33` 只列 `ci.yml`＋`scripts/**`＋`.scratch/wisp/probes/**`）。⇒ 乙要真满足 AC#1 那句"所有引用同批改净"，要么编排者具名把写面扩到 `docs/**`，要么台账那 8 处按"只追加不删"追加更正（`pending-and-issues.md:9610` 已示范过一次追加更正）。
4. `ci.yml` 零字节。
5. ⚠ 另两枚 overlay 名册（`overlay-ac7-mutant.json`、`overlay-muta.json`）指 `mut/taintmatch.go`／`mut/fs.go`，**不指 `fs_broken.go`**（现读三枚 JSON 全文）⇒ 不在连带面里；且 `mut/fs.go`／`mut/taintmatch.go` 本身在 §1 现量#5 的 19 枚名册之外（gofumpt 干净），没有任何必要动它们。

**②会不会推翻"分母＝tracked 集合"那两句＝不会（待验断言"乙不会"成立）**

- `:189-193` 那句讲的是**分母**（哪一批字节被喂给尺）。乙一字不改分母，改的是**样本的后缀名**——改名后它自动既不在 `git ls-files '*.go'`（`attrib.sh:341` 的取集式）里，也不在 `gofumpt -l .` 的文件发现里。
- `:186-188` 的"两把尺读同一批字节、必须同意"照旧为真：现量两把尺在 HEAD 上读到的是同一批（§7 那格讲的"同不同意"是另一层），改名后仍是同一批 18 枚。⇒ 分母那句不被触碰。
- `attrib.sh` 的归因规则（`:174-193`）也不被触发——那两条只在尺**报出某行**时才用得上；后缀改掉＝尺根本不报这一枚，既不算豁免也不算归因失败。⇒ 161 AC#7 的两形读法不受影响。

**③会不会咬到引用那枚样本路径的件＝会，三形里唯一真咬的一形**

- 尺＝`grep -rln fs_broken .`＝**60 枚**（tracked 37／未入库 23），逐枚读过分三档（逐枚名单与判语在 §6）：
  - **机器真会去 open 它**：`1` 枚＝`.scratch/wisp/probes/185/c1/overlay-proof.json`。
  - **脚本会去 open 它**：`0` 枚。尺＝那 60 枚清单里 `grep -E '\.(sh|py|ps1|js|mjs|go)$'` → **NONE**；`git grep -ln overlay-proof` → 4 枚（`212/a3/work/classified.tsv`、`236/a1/census.md`、`236/a2/census.md`、`docs/evidence/s1/185-reread-owner-census-c1.md`），全是**记录**过那次执行的件，没有一把尺按它跑。⇒ 本仓没有"自动复跑 185 overlay 正控"的入口，乙的硬连带是 1 枚而不是几十枚。
  - **注释性／记录性引用**：其余 59 枚，内部再分两种承重：把「**路径＋`:4:1`＋错误文案**」钉在同一串里的＝**40 枚文件（tracked 19 枚）**（尺：`grep -rl -F "fs_broken.go:4:1"`），改名后这 40 串的路径段与盘上不符＝**phantom citation 候选面**；只提"这枚文件存在／归属票 185"的＝改后缀后仍真（文件还在，只是名变了）。
  - **引"错误文案"而不引路径的**：全仓含 `imports must appear` 字样＝**57 枚文件（tracked 21 枚）**——它们说的是 gofumpt 的输出串，改名**不影响其为真**（同内容仍产同诊断，只是前缀路径变）。
- ⚠ 台账：`grep -n fs_broken docs/reports/pending-and-issues.md`＝**8 处**（`:8983`、`:9566`、`:9610`、`:9611`、`:9614`、`:10721`、`:12086`、`:12111`）逐字写着那枚路径；按 AGENTS.md §1.5「只追加不删」与票 171 `:35`（`docs/reports/**` 零字节）⇒ 乙落地时这 8 处的处置**必须另批**，不是落地腿能顺手改的。

**④门剩什么牙（判语）**

- **乙单独做不让门变绿**：§1 现量#5 现量——剔掉 `fs_broken.go` 之后 tracked 分母仍有 **18 枚**脏样（17 枚 `.scratch` 样本＋`tools/d22scan/selftest.go`）⇒ `:168` 的 `OUT` 仍非空、仍 `exit 1`；`:184` 的 (A) 尺仍报 `TRACKED-DIRTY=18`、仍 `exit 1`。⇒ **链头红仍在 `:168`，`:206`／`:209` 继续采不到**。这一格与票面题面的预设相反（票面把乙写成"把样本从仪器视野搬走"就能解红），是乙最贵的代价。
- 乙买到的两样真东西（都是"牙的形状"变化，不是偏好）：
  - (一)解析错误消失 ⇒ `gofumpt` 退码从 **2 变 0** ⇒ `bash -e {0}`（CI 日志逐字 `shell: /usr/bin/bash -e {0}`；凭据 `.scratch/wisp/probes/ci/1005-a/logs/lint-full.log`）不再在 `:171` 掐死这一步 ⇒ `:172-175` 那三行**第一次真的执行**：打印 `gofumpt would reformat:` ＋ 18 行名册。**今天那一步的日志只有两行**（`##[error]...fs_broken.go:4:1: imports must appear before other declarations` ＋ `##[error]Process completed with exit code 2.`，`echo "$OUT"` 一字未出现）⇒ 门从"只报退码"变"报名册"，这是 AC#2「牙是可见的」的必要条件之一。
  - (二) `attrib.sh` 的 (A) 尺从「`exit 2`＝拒答（这棵树这把尺读不了）」变「`exit 1`＋逐枚列 `TRACKED-DIRTY`」⇒ `:184` 那步第一次给出它自己设计的读数（即使仍红）。
- ⚠ 乙的射程要说准：它只对**按后缀发现文件的格式仪器**有效。`go vet ./...`（`:206`／`:209`）**本来就看不见 `.scratch/**`**——Go 的包模式跳过以 `.`／`_` 开头的目录，且 `probes/**` 下自带 **14 枚 `go.mod`**（现读 `find /tmp/269a1/full -name go.mod | wc -l`＝14）自成一模块。⇒ 票面标题那句"会被 Go 仪器看见"对 `go vet` 是**过头的**；本条未跑 `go` 证实（禁），按 Go 公开语义推断，具名归 §7。

### 丙＝给 `:184`／`:206`／`:209` 三步补 `if: ${{ !cancelled() }}`（不修红，只让红不再吞读数）

**①要动哪几枚文件**

- 只有 **`.github/workflows/ci.yml`** 一枚：在 `:184`／`:206`／`:209` 三步体内各插一行 `        if: ${{ !cancelled() }}`（八空格缩进，插在 `- name:` 之后、`run:` 之前，照 `:263`／`:300` 的既有排版）。
- **纯插入形状**：预期 `git diff --numstat`＝**3 增 0 删、3 个 hunk**。这格要说白：本仓最硬的"没顺带动别的守卫"尺是 **hunk 枚数**，三处各插一行＝三枚 hunk＝**恰好等于三处**；多一枚 hunk 就是动了第四处。落地腿必须把 `--numstat` 的删除列＝0 一并交出来。
- `:206`／`:209` 两步体内现读只有 `- name:`＋`run: go vet ./...`（＋`:211` 的 `working-directory: tools/d22scan`），**没有任何注释**⇒ 插入点无歧义。
- 【禁】 不许"顺手"给 `:168`／`:136` 之类补守卫（那是第四、第五处，越出票面对丙的定义）。

**②会不会推翻"分母＝tracked 集合"那两句＝不会（待验断言"丙不会"成立）**

- 丙不碰任何一步的命令体、不碰任何一步的入参集合；`:186-188` 的"两把尺读同一批字节、必须同意"与 `:189-193` 的"分母＝TRACKED 集合"逐字照旧为真。
- ⚠ 但丙**改变那两句的今天状态**：`:186` 那句 `adds NO REACH TODAY` 的前提是"`:184` 与 `:168` 同生同死"。丙之后 `:184` 在 `:168` 红时**也会作答** ⇒ 那一步第一次有独立射程（它比 `:168` 多一条空心守卫、多一层 `exit 2` 语义，见 §4）。⇒ 严格说：**丙不推翻那两句，但让 `:186` 那句"今天不增加射程"变成过期描述**（那是自陈而非裁定；改不改注释是可选项，不是必需项，本腿不替落地腿决定）。

**③会不会咬到引用那枚样本路径的件＝一枚都不咬**

- 样本原地不动（`fs_broken.go` 字节未变，`git rev-parse HEAD:...`＝`ee583853...` 现读）⇒ §6 那 60 枚引用件、`overlay-proof.json` 那枚硬连带、台账 8 处全部零风险；AC#1 自动满足，而且是**唯一一枚不需要任何引用治理就能满足 AC#1 的形**。
- 唯一新增的文字面风险：新插的三行若配注释去写样本路径，就把那条路径第二次落地。⇒ 形状要求：注释只指票号（`ticket 269`），【禁】 不写行号（本仓定式凭据＝`docs/reports/pending-and-issues.md:11905` 那句「凡引用代码位置，`grep` 被指的那串字符，别抄任何来源的行号（含 CI 日志）」）。

**④门剩什么牙（判语，不是偏好）**

- **丙单独做＝门仍常红＝等价于没人看**：`:168` 照旧红（19 枚脏样）、`:184` 照旧红（`exit 2`），lint 的 job 级颜色一字不变。⇒ 以"格式有没有坏"为命题的那扇门，在丙之后**仍然每天给出同一个不区分任何东西的红**——红已不含信息量：任何人看它都得到同一答案，因此没有任何人会去看它。票面 `269:20` 那句"丙单独做＝门仍常红＝等价于没人看"按现量**复认成立**。
- 丙真正买到的是**另外三枚门的读数**，且这是本票第二半唯一的解：`:184` 从此每次给出 (A) 尺的红因名册；`:206`／`:209` 从此每次真跑 `go vet ./...` 并出名册——今天这两枚自 09-28（`4813567e`）以来**零枚读数**（凭据：`.scratch/wisp/probes/ci/1005-a/verdict.md:36` 现量三步 conclusion＝`skipped`；台账 `:9566` 逐字「**后果不是"多一枚红"，是它把后面的 `go vet (module)` 与 `go vet (tools/d22scan module)` 一起吃掉 ⇒ 今晚两发零 `go vet` 读数**」）。
- 代价必须同批报：**"门开始作答"会把新的红带进来**。`go vet` 一旦响，lint 的红名册从"1 枚格式门"变成"N 枚 vet 名"，其中可能含从未被 CI 看见过的包（本腿禁 `go vet` ⇒ 几枚、会不会红，**量不到**，具名归 §7）。这不是回归，是新增可见性；但在"CI 得变绿"的压力下它正是最容易被读成"丙把 CI 弄红了"的形状 ⇒ 裁丙的人要把这句话提前写进批语。
- 逐字规矩核（细节在 §3）：文件头 `:5-7` 禁的是 `if: false`／`continue-on-error`／skip flags **三形**，`!cancelled()` 不在其列；`:412-416` 是这枚文件对同一形状的**自我定性**（逐字：`# WHY THIS IS NOT A DOOR REMOVED (the two forbidden shapes, both absent): no` ／ `# step was deleted, and NO step here carries `continue-on-error`, whose entire` ／ `# documented meaning is "prevents a job from failing when a step fails". A step` ／ `# run under `!cancelled()` that exits non-zero still fails the job exactly as` ／ `# before - the only change is that the steps after it also get to answer.`）；`:246-249` 逐字（这四行的注释前缀是 `#` 加四个空格，属上层列表的续行，照原样引）：`#    `gofmt (gofumpt)`). `if: ${{ !cancelled() }}` is the one shape` ／ `#    ticket 85's ruling allows for an observation period: it does not` ／ `#    change this step's own red/green, it only makes sure the step` ／ `#    answers. continue-on-error is NOT used (D22 run-away mode 6,`（下一行 `:250` 接 `ci.yml header). The `mockllm module vet` step below already`／`:251` 接 `carries the same guard from ticket 111 and is left untouched;`）。**这三处都是"逐字那行"，不是"别的步也这么写了"。**
- ⚠ 同一枚文件里还有**两句反向的、逐字的批准面**，裁丙的人必须处置，本腿不自作主张：
  - `:178-180`（票 161 给那一票划的界）逐字：`# ADDED STEP ONLY (ticket 161 AC#7 form (A), decided by the orchestrator's` ／ ``# 00:0x ruling "裁二"): no existing step was moved, edited, deleted, given an`` ／ ``# `if:`, or made continue-on-error.``。读法分歧：这句约束的是"161 那次编辑"，还是"这棵树上的既有步骤从此不许给 `if:`"？按字面是前者；`:263`／`:300`（票 85／111 后来确实给既有步加了 `if:`）说明仓里按前者执行过。
  - 票 171 `:31`（更宽的一句）逐字：`⚠ 动触发表／并发组／既有步骤＝**契约级（票 134 地界）⇒ 停手报回**`。若读成**仓级规矩**，丙（改三枚既有步骤）＝契约级 ⇒ 停手报回；若读成 171 自家写面，丙只需 269 票面 `:33` 的写面授权。
  - 台账 `:11905`（A609 §3 第③行，编排者自己的实操先例）逐字：「③ 给那一步加 `if:` 绕过 skipped＝动 `ci.yml:136` 那行注释自陈的 **A585 批准面**（"no existing step was moved, edited, deleted, given an `if:`"），【禁】（原文那枚符号是 stop-sign emoji，逐字引时不改）我不越。」⇒ 按这条先例的**口径**逐枚判：丙的三处里只有 **`:184`** 的注释（`:178-183`）含那句话，**`:206`／`:209` 两步体内没有任何自陈**⇒ 重叠面＝1/3。"三处一起补 vs 只补不重叠的两处"这一格归编排者裁。

## §3 守卫分层逐枚表（49 枚 step 全表，行号一律现读）

尺：`awk` 扫 `.github/workflows/ci.yml` 全 766 行，凡 `^      - (uses|name):` 起一步、到下一步之前为止，检 `^        if:` 与 `^        continue-on-error:` 两种键。读数时刻 `12:0x +08`，锚 `50e0299c`（该文件最后一笔＝`519e1ade` 10-04 09:42，工作树对 HEAD 干净）。**列里的"if 行号"＝那一步体内 `if:` 键的绝对行号，`—`＝无。**

| job（起始行） | step 起始行 | if 行号 | continue-on-error | 步名 |
|---|---|---|---|---|
| lint（`:65`） | `:68` | — | 无 | uses: actions/checkout@v4 |
| | `:70` | — | 无 | uses: actions/setup-go@v5 |
| | `:74` | — | 无 | D22 scanner positive control (tools/d22scan tests, seeded red) |
| | `:83` | — | 无 | D22 scanner self-test (ticket 161 AC#2 - every ban, both directions) |
| | `:108` | — | 无 | D22 seven-ban + emoji scan (tools/d22scan) |
| | `:136` | — | 无 | Tracked path-length budget (ticket 262) |
| | `:168` | — | 无 | **gofmt (gofumpt)** ← 今天的链头红 |
| | `:184` | — | 无 | **gofmt (gofumpt) - the tracked set is the denominator (ticket 161 AC#7 form A)** |
| | `:206` | — | 无 | **go vet (module)** |
| | `:209` | — | 无 | **go vet (tools/d22scan module)** |
| | `:213` | `:263` `!cancelled()` | 无 | staticcheck |
| | `:289` | `:300` `!cancelled()` | 无 | mockllm module vet |
| test-core（`:309`） | `:314` | — | 无 | uses: actions/checkout@v4 |
| | `:316` | — | 无 | uses: actions/setup-go@v5 |
| | `:320` | — | 无 | Start compose test services (mock-llm on 18080) |
| | `:323` | — | 无 | Probe mock-llm |
| | `:335` | — | 无 | Environment fork assertion (WISP_ENV=test data dir) |
| | `:351` | — | 无 | Portable package tests (core scope; list 在 scripts/portable-tests.sh) |
| | `:375` | `:376` `always()` | 无 | Stop compose services |
| test-windows（`:419`） | `:424` | — | 无 | uses: actions/checkout@v4 |
| | `:426` | — | 无 | uses: actions/setup-go@v5 |
| | `:430` | `:436` `!cancelled()` | 无 | Windows ACL sealing gate (internal/winsec's own tests, ticket 110) |
| | `:473` | `:474` `!cancelled()` | 无 | Cache third_party (deps.toml key) |
| | `:480` | `:484` `!cancelled()` | 无 | cgo build smoke (build.ps1 fetch-deps + mingw link + doctor) |
| | `:487` | `:506` `!cancelled()` | 无 | cmd/wisp CLI tests (ticket 111 AC#4) |
| | `:509` | `:542` `!cancelled()` | 无 | Portable windows tests (proc/secret/config/risk/ball/perm/plugin/llmrecord) |
| | `:545` | `:558` `!cancelled()` | 无 | PathResolver junction placeholder (real cases tickets 18/20) |
| slo-smoke（`:564`） | `:569` `:571` `:575` `:581` `:584` `:587` | 全 `—` | 无 | checkout／setup-go／Cache third_party／Build wisp.exe／SLO smoke gate／Upload SLO report |
| slo-full（`:622`） | `:637` `:639` `:644` `:647` `:650` | 全 `—` | 无 | checkout／setup-go／Build wisp.exe (deps cached)／SLO full gate (six states + settle + leak)／Upload SLO report |
| lint-frontend（`:696`） | `:702` `:704` `:712` `:715` `:718` `:721` `:727` `:730` `:749` `:755` `:760` | 全 `—` | 无 | checkout／setup-node／npm ci／typecheck／lint (oxlint)／token drift guard／vite build／L2 card render／Composer states／Streaming SSE row／Nav rail icons |

**分层读数（口径与推导）**

- 总 step＝**49**（`grep -c '^      - (uses|name):'` 口径＝上表行数）。带 `if:`＝**9**：`!cancelled()` 8 枚（`:263`、`:300`、`:436`、`:474`、`:484`、`:506`、`:542`、`:558`）＋`always()` 1 枚（`:376`）。带 `continue-on-error`＝**0**（仪器逐枚扫＝0；`grep -n continue-on-error` 的 13 处命中 `:6`、`:34`、`:104`、`:124`、`:138`、`:180`、`:249`、`:413`、`:457`、`:470`、`:607`、`:658`、`:694` 逐条现读**全在 `#` 注释行**）。
- **分层是真的、且按 job 分**：同一 lint job 内部就分层（`:213`／`:289` 有，`:168`／`:184`／`:206`／`:209` 没有）；test-windows 6/6 全有（`:383` 逐字自陈 `# EVERY STEP BELOW THE SETUP ONES CARRIES `if: ${{ !cancelled() }}` (ticket 111`／`:384` `# AC#6).`）；test-core 只有 `:375` 那枚清理步带 `always()`；slo-smoke／slo-full／lint-frontend 三枚 job 零枚。⇒ 票面 现量#2 那句"守卫在这棵树上是分层的，不是全有或全无"**成立**，但枚数是 9 不是 7。

**文件头 `:6` 逐字禁的是哪三形（原文依据）**

- `:5-7` 逐字：`# (full six-state SLO + settle + leak). Per D22 mode-6 NO job is configured` ／ `# skippable: no `if: false`, no continue-on-error, no skip flags. Runner` ／ `# platform limits (no interactive desktop / audio on hosted runners) are` ／ `# handled by SUBSET choice, never by skipping a job (SPEC-10 §3).`（`:4-8` 全引是为了不截句）⇒ **逐字禁的三形＝`if: false`／`continue-on-error`／skip flags**，`!cancelled()` 不在名册里。
- ⚠ 另一句常被误当整文件规矩的：`:34-35` 逐字 `# No `if:`, no continue-on-error, no skip flag anywhere in this block`／`# (D22 mode-6: "A skippable job is a job that will one day be skipped").`——它的射程写死 **"in this block"**，指 `:10-37` 那整块 `on:` 触发表（`:36` 就是 `schedule:`），不是全文件。⇒ **拿 `:34` 反对丙是引错了射程**；能反对丙的只有 `:6` 那三形（不含 `!cancelled()`）与"动既有步骤"那一层（`:178-183`／票 171 `:31`／台账 `:11905`），见 §2 丙④ 与 §7。

**"给现存三步补 `!cancelled()` 算不算破这条自陈规矩"——逐字判定**

- 判：**不破 `:6`**。依据三处，都是逐字行：(1) `:6` 只名列三形（上引）；(2) `:412-416` 文件自我定性（"WHY THIS IS NOT A DOOR REMOVED … A step run under `!cancelled()` that exits non-zero still fails the job exactly as before - the only change is that the steps after it also get to answer."）；(3) `:246-249`（"is the one shape ticket 85's ruling allows for an observation period: it does not change this step's own red/green, it only makes sure the step answers. continue-on-error is NOT used (D22 run-away mode 6, ci.yml header)"）——这段甚至**明写它没破文件头那条**。
- 判：**"算不算破别的规矩"未决**，两把尺对着摆，归编排者：票 171 `:31` 的"动既有步骤＝契约级（票 134 地界）⇒ 停手报回" vs `:263`／`:300` 两次既有先例（票 85／111 都动过既有步骤并留下逐字理由）。【禁】 本腿不裁，也【禁】 不用"别的步也这么写了"当理由——上面 (1)(2)(3) 指的是逐字行。

## §4 `:184` 那步今天真跑起来会怎样（脚本射程逐枚点名）

尺＝逐行读 `.scratch/wisp/probes/161/r5/attrib.sh`（**440 行**，`wc -l` 现读）＋用等价命令复现它 `:341` 那一条管线（**没有真跑 `attrib.sh`**：它 `:137` 要调 `go env GOPATH`，`go env` 不在派单允许清单里，【禁】 我没有绕过去跑）。

1. **它喂的是什么分母**：`ci.yml:204` 调 `--tracked-only` ⇒ `attrib.sh:108` 置 `MODE=tracked` ⇒ `:405-406` 只跑 `ruler_a`。`ruler_a` 的分母＝`:341` 那一条 `git ls-files -z '*.go' | xargs -0 "$GOFUMPT" -l`——**tracked 全集，一枚样本都不剔**（`:37` 注释逐字：`(A) TRACKED TREE   git ls-files '*.go' | xargs gofumpt -l`，`:38-39`：`MUST BE EMPTY. This is the form CI executes, because a CI checkout contains tracked bytes and nothing else.`）。现量分母＝**921 枚 @起手锚 `05680e44`／922 枚 @`65142d9a`（12:23 复跑，多的那一枚是 259-r1 那枚负控）**（尺＝`git ls-files '*.go' | wc -l`；`c6cf66e6` 那枚 run 锚点上＝**904 枚**，尺＝`git ls-tree -r --name-only c6cf66e6 | grep -c '\.go$'`）。
2. **空心守卫在不在**：在，`:337-338`：先 `TRACKED_GO=$(git ls-files -z '*.go' | tr -dc '\0' | wc -c | tr -d ' ')`，再 `[ "$TRACKED_GO" -gt 0 ] || { … refusing to report 'empty' from a ruler that was handed nothing >&2; exit 2; }`——正是票面要的"handed ZERO tracked .go files 就 exit 2"。`:330-335` 那段注释逐字交代了它为什么是硬停而不是脚注（`git ls-files -Z` 那次恒真检）。
3. **它会不会也被 `fs_broken.go` 顶红＝会，两条路径都红（这一条决定甲形的必要性）**：
   - 今天（含 `fs_broken.go`，它是**解析错误**不是"要格式化"）：`gofumpt -l` 对该文件现读＝stderr 一行 `.scratch\wisp\probes\185\c1\mut\fs_broken.go:4:1: imports must appear before other declarations`、**单发 rc=2**（我这发实测）。在 `:341` 那条 xargs 管线上，921 枚分文件递 ⇒ **聚合 rc＝123**（我复现同一条管线实测 `xargs rc=123`）⇒ `:342-345` 命中 `if [ "$A_RC" -gt 1 ]` ⇒ 打「`(A) gofumpt exited 123 on the tracked set (>=2 means a file would not parse) - the tracked tree is not readable by this ruler`」并 **exit 2**。⇒ `:184` **红**（红＝拒答，不是绿）。
   - 乙之后（解析错误没了，剩下 18 枚"要格式化"）：`gofumpt` 对"只要格式化"的文件现读 rc＝**0**（我实测两发：`gofumpt -l cmd/wisp/models.go` 打出文件名、rc=0；`-l tools/d22scan/selftest.go` 同）⇒ `A_RC=0`、`A_OUT` 非空 18 行 ⇒ `:352` 逐枚计 `TRACKED_DIRTY` ⇒ `:411-417` 打 `RED (tracked-only) - tracked-dirty=18 …` 并 **exit 1**。⇒ `:184` 仍红，但**从"拒答"变成"报名册的红"**。
   ⇒ 结论具名：**甲（只改 `:168` 的分母）对 `:184` 完全无效**，因为 `:184` 的分母写在 `attrib.sh:341` 里，不写在 `ci.yml:171` 里。要么甲的定义扩到"改 `attrib.sh` 的取集式"（那是 `.scratch/wisp/probes/**` 的**尺本体**，票面 `:33` 写面允许），要么甲落地后链头红从 `:168` 挪到 `:184`（`:184` 无 `if:`，见 §3 表）⇒ **`:206`／`:209` 照样被吞**。票 236 `:101` 那句"乙（改 `attrib.sh:341` 的分母）**碰不到红的那一步**"讲的是 09-29 那棵树的红位（当时在 `ci.yml:136`），今天红位在 `:168`＋`:184` 两处，那半句已过期（§8 第 5 条）。
4. **一条隐蔽耦合（三形都要知道，具名）**：`attrib.sh:138-140` 从 `$(go env GOPATH)/bin/gofumpt(.exe)` 取二进制，而 CI 里那枚二进制是**上一步 `:170` 的 `go install mvdan.cc/gofumpt@latest` 装的**（现读 `ci.yml:170`）。⇒ `:184` 悄悄依赖 `:168` 至少跑到过 install 那一行。今天 `:168` 在 `:171` 被 `bash -e` 掐死时 install 已完成（`go install` 在 `:170`，早于 `:171`；CI 日志现读 `go: downloading mvdan.cc/gofumpt v0.12.0` 之后才出现那行 `##[error]`）⇒ 现状成立。但**丙之后**多一种失败模式：若某次 `go install` 自己失败（网络／上游改版），`:184` 会以 `:140`「no gofumpt at … - refusing to answer with a ruler that did not run」**exit 2 红**——这是"新增可见性"，不是回归，裁丙的人要预先认下这一格。
5. **版本一致性（读数可迁移的凭据）**：本机 gofumpt `--version`＝`v0.12.0 (go1.27.1)`；CI 日志 `lint-full.log` 现读 `go: downloading mvdan.cc/gofumpt v0.12.0` ⇒ 同版。⇒ §1 现量#5 的 19 枚名册与 CI 那台检出**同尺同版**，差异只在行尾形状（§0 末条）与 `go.mod` 语言档（我第一版只抽 `*.go` 的归档漏读 `probes/263/v1/src/main.go`，因为同目录那枚 30 字节 `go.mod` 决定 gofumpt 的语言档；改整树归档后名册 18→19，见 §9 第 2 条）。

## §5 样本全名册（口径自己定，三把尺并列，逐枚路径）

【口径】"故意写坏的样本"没有一把通吃的尺，所以给三把，各自可复算、各自只答自己那问：

- **尺①（决定门的那把）**：CI 形状下 tracked 字节里被 `gofumpt -l` 报出来的＝**门真的会因此红的枚数**。尺＝HEAD 整树归档（`git archive HEAD | tar -x`，`find -type f`＝`git ls-files | wc -l` 逐锚点对齐）＋`git ls-files -z '*.go'` 逐枚前缀＋`xargs -0 gofumpt -l`（`v0.12.0`）。
- **尺②（盘上在场、CI 永远看不见）**：`.scratch/wisp/probes/**` 下**不在 git 索引里**的 `.go`，分两档（两档都进不了 CI，但本机 `gofumpt -l .` 会走它们——这就是我本机读数 47 枚 vs CI 形状 19/20 枚的另一半原因）：
  - ②a **未入库也未忽略**＝`git ls-files -o --exclude-standard -- .scratch/wisp/probes | grep -c '\.go$'`＝**46 枚**（含 `161/r2/ctl/neg/n.go`、`161/r2/ctl/pos/p.go`、`161/r2/ctl/pos2/r.go`、`161/r5/negative-control/bad-sample.go` 这 4 枚 161 家）。
  - ②b **被目录级 `.gitignore` 屏蔽**＝`git ls-files -o -i --exclude-standard -- .scratch/wisp/probes | grep -c '\.go$'`＝**325 枚**（逐目录：`161/r1/**` 322 枚、`161/r3/**` 3 枚）。凭据尺：`git check-ignore -v .scratch/wisp/probes/161/r3/pre-format/guard.no1.go` → `​.scratch/wisp/probes/161/r3/.gitignore:1:pre-format/*.go`。⇒ 这正是记忆档那条"只建不入库（配 `probes/<票号>/<程>/.gitignore`，`gofumpt` 仍会列它、但不会进别人的提交）"的**现场读数**，也是"故意坏样本绝不该以 `.go` 后缀入库"这条规矩**今天被违反了 276 次**（tracked 那一档）的反面教材。
  - 合起来：`.scratch/wisp/probes/**` 下的 `.go` 盘上总在场＝**276（tracked）＋46（②a）＋325（②b）＝647 枚**，其中门（CI 形状）真正会读到的只有 tracked 里那 19/20 枚。⇒ 任何"把门改成看工作树"的修法会把这 46＋325 枚一起吃进射程——**这条是甲/乙两形都不许碰 `gofumpt -l .` 那一味的原因之一**（记进 §2 甲②）。
- **尺③（自陈标记、但格式门看不见的）**：tracked `.go` 里内容带"这是故意坏的"字样的（`MUTANT|mutant copy|bad-sample|deliberately (broken|unformatted|unparseable)|故意写坏|故意不解析|故意坏`）＝**9 枚**，其中**只有 1 枚**落在尺①里。⇒ 尺①系统性**少数**："行为坏了但排版干净"的样本（多数突变台件）格式门永远不知道，这本身是对的（门的命题只是格式），【禁】 不许有人拿"尺①只有 20 枚"去当"仓里只有 20 枚样本"。

**尺①名册（逐枚；`<A>`＝该票票面文件 `.scratch/wisp/issues/<NN>-*.md` 是否真存在，这是 `attrib.sh:186-193` 的归因第二条）**

| # | 路径（相对仓根） | 归属票 | 票面在 | 脏因（现读） |
|---|---|---|---|---|
| 1 | `.scratch/wisp/probes/33/p1/q1/main.go` | 33 | 有 | 纯格式 |
| 2 | `.scratch/wisp/probes/33/p1/q2/main.go` | 33 | 有 | 纯格式 |
| 3 | `.scratch/wisp/probes/33/p1/q3/main.go` | 33 | 有 | 纯格式 |
| 4 | `.scratch/wisp/probes/163/a1/main.go` | 163 | 有 | 纯格式 |
| 5 | `.scratch/wisp/probes/174/c2/zz174c2_wiring_pair_windows_test.go` | 174 | 有 | 纯格式 |
| 6 | `.scratch/wisp/probes/183/r2/mut/task-boxset-off.go` | 183 | 有 | 纯格式（突变台件副本） |
| 7 | `.scratch/wisp/probes/185/c1/mut/fs_broken.go` | 185 | 有 | **解析错误** `:4:1 imports must appear before other declarations`，单发 rc=2 ← 本票题面那枚 |
| 8 | `.scratch/wisp/probes/185/r1/mut-m1/hostpath_185.go` | 185 | 有 | 纯格式 |
| 9 | `.scratch/wisp/probes/197/r1c/pre/subagent_197.go` | 197 | 有 | 纯格式（产品码改前快照） |
| 10 | `.scratch/wisp/probes/197/r1c/pre/subagent_197_test.go` | 197 | 有 | 纯格式（同上） |
| 11 | `.scratch/wisp/probes/212/v1/mut/main-noq9.go` | 212 | 有 | 纯格式 |
| 12 | `.scratch/wisp/probes/220/r1/prechange/prechange220_probe_test.go` | 220 | 有 | 纯格式（改前快照） |
| 13 | `.scratch/wisp/probes/222/v1/mutations/m11-budget-50ms-gated/subagent_222_test.go` | 222 | 有 | 纯格式 |
| 14 | `.scratch/wisp/probes/224/v2/dialect_probe_test.go` | 224 | 有 | 纯格式 |
| 15 | `.scratch/wisp/probes/235/v1/mut/prereadoff-m3_222_test.go` | 235 | 有 | 纯格式 |
| 16 | `.scratch/wisp/probes/241/v1/posctl/badly_formatted.go` | 241 | 有 | 纯格式（**文件名就自称是坏样本**） |
| 17 | `.scratch/wisp/probes/241/v1/probe_v1_readings_test.go` | 241 | 有 | 纯格式 |
| 18 | `.scratch/wisp/probes/263/v1/src/main.go` | 263 | 有 | 纯格式（同目录另有 30 字节 `go.mod`，见 §4 第 5 点） |
| 19 | `tools/d22scan/selftest.go` | —— | —— | **非样本**：纯格式（map 字面量键对齐），入库笔 `5e8748b3` 10-03 21:07，**不在 269 写面** |
| 20 | `.scratch/wisp/probes/259/r1/gofumpt-negctl/probe.go` | 259 | 有 | 纯格式（`func   f( )  int`，**目录名就叫 gofumpt-negctl**）——⚠ **本腿工作期间入库**：`50e0299c` 10-05 12:09，由在飞写腿 `259-r1` 提交 |

- 枚数随锚点动（具名，【禁】 不许当常量）：起手锚 `05680e44`（11:52）＝**19 枚**；`50e0299c`（12:09，259-r1 入库那枚负控）起＝**20 枚**（`12:23:03` 复跑，两次归档 `diff` 差集正好这一条）。分母同步漂：tracked `.go`＝**921**（`05680e44`）→**922**（`65142d9a`，12:23）。⇒ 裁甲时"名册封顶值"必须现量重读；**这一格是甲的维护代价的现场证明**（票 236 `:115` 那句"乙法是个追不完的移动靶"讲的是 236 自己的另一支（把脏样逐枚整理干净），今天同形复现、方向反过来：动的是新入库的台件，不是人）。
- 尺③那 9 枚（tracked、内容自陈带坏样本字样）逐枚：`probes/160/c1/shape160.go`、`probes/160/r1/mut/provenance.no-identity.go`、`probes/175/r1/mut-ac3/bridge.name-gated.go`、`probes/185/c1/mut/fs.go`、`probes/185/c1/mut/fs_broken.go`、`probes/185/c1/mut/taintmatch.go`、`probes/185/c1/pointer_185_c1_ac7_probe_test.go`、`probes/221/v1/mutations/m12-finalize-wrong-state.go`、`probes/221/v1/mutations/m6-cancel-before-refusal.go`。⇒ 题面"我数到 185 三枚"核对：`git ls-files -- .scratch/wisp/probes/185/c1/mut/`＝**3 枚**（`fs.go`／`fs_broken.go`／`taintmatch.go`）**一致**；但其中**只有 `fs_broken.go` 在尺①里**（另两枚是 overlay 基准副本、排版干净）。185 名下 tracked `.go` 共 **8 枚**（再加 `pointer_185_c1_ac7_probe_test.go` 与 `r1/mut-m1..m4` 4 枚，其中 `r1/mut-m1/hostpath_185.go` 也在尺①）。
- 题面问的"236／161／267 那族还可能有"逐枚回答（尺＝`git ls-files '<dir>/**/*.go'` 与 `git ls-files -o -i --exclude-standard`）：
  - `236/**` 与 `267/**` 名下 tracked `.go`＝**0 枚**（现读两路 pathspec 都为空）。267 名下那枚先例件是 `267/r2/pristine/models.go.pristine`——**它已经不是 `.go` 后缀**，正是乙要的形状的在册证据。236 名下与本案相关的都是 `.md`／`.txt`／`.bin` 记录件（§6 的 R／H 档）。
  - `161/**` 名下 tracked `.go`＝**11 枚**（`r4/removed/selftestsamples.no-*.go` 2 枚＋`v1/samples/rootB/internal/v1/*.go` 9 枚，尺＝`git ls-files '.scratch/wisp/probes/161/**/*.go'` 逐枚列），**尺①里一枚都没有**（那 11 枚排版干净或本就是 d22scan 的取样件，不是格式样本）；另加 ②a 的 4 枚与 ②b 的 325 枚（`161/r1/**` 322＋`161/r3/pre-format/**` 3）——⚠ **`161/r3/pre-format/guard.no1-3.go` 是"被 `probes/161/r3/.gitignore:1` 屏蔽"而不是"未入库未忽略"**，别与票 161 `:62` 那三枚 `158/accept-r1/mut/guard.no1-3.go` 混掉：后者**是 tracked 件**（现读 `git ls-files -- .scratch/wisp/probes/158/accept-r1/mut/`＝9 枚，含 `bridge.M1-M3.go`／`guard.as-is/stub.go`），且它们按 AC#7① 已被格式化干净 ⇒ **不在尺①**（台账 `:9610` 与票 161 `:62` 讲的都是 158 那一族，不是 r3 的 pre-format 快照）。⇒ 本条把"161 那族还可能有"的答案钉成：**161 名下零枚 tracked 格式样本**。
  - `185/**` 名下 tracked `.go`＝**8 枚**（`c1/mut/` 3＋`c1/pointer_185_c1_ac7_probe_test.go` 1＋`r1/mut-m1..m4/` 4），其中**只有 2 枚**在尺①（`c1/mut/fs_broken.go` 解析错误、`r1/mut-m1/hostpath_185.go` 纯格式）；`r1/mut-m2/mut-m3`（`provenance.go`）与 `mut-m4`（`hostpath_185.go`）排版干净 ⇒ 门的读数里没它们。⇒ 题面"185 三枚"这一把**一致**，但"三枚里只有一枚真的顶红门"是新增读数。

## §6 改名那一形（乙）的连带面——60 枚逐枚读过，判"是不是把路径当盘上真身指"

尺：`grep -rln --exclude-dir=.git fs_broken .`＝**60 枚**；tracked 判定尺＝逐枚 `git ls-files -- <p> | wc -l`；三档口径：**M＝机器真会去 open**／**R＝规范性引用（未来的腿会照它定位这枚文件，改名即 phantom citation）**／**H＝历史记录快照（记的是当时的读数，改名不改它的真值，按"只追加不删"不许改写）**。

**M＝1 枚（唯一硬连带）**

- `.scratch/wisp/probes/185/c1/overlay-proof.json`（tracked）：`{"Replace": {"D:/…/internal/tools/fs.go": "D:/…/185/c1/mut/fs_broken.go"}}`，`go build -overlay` 逐字节 open 右值 ⇒ 乙必须同批改它，否则 185 那发"种 X 必响"的正控复跑即失败（AC#1 的"还要能被指到"就断在这里）。
- **脚本层＝0 枚**：60 枚清单里 `grep -E '\.(sh|py|ps1|js|mjs|go)$'` → **NONE**；`git grep -ln overlay-proof` → 4 枚（`212/a3/work/classified.tsv`、`236/a1/census.md`、`236/a2/census.md`、`docs/evidence/s1/185-reread-owner-census-c1.md`）全是**记录**，无一把尺按它跑。⇒ 本仓没有自动化复跑 185 overlay 的入口（这是运气，不是设计；乙落地时应顺手把它记成一条票面残余）。

**R＝11 枚（改名后必须同批改净，其中 5 枚超出 269 写面）**

| 路径 | hits | 引的是什么 | 写面判定 |
|---|---|---|---|
| `.scratch/wisp/issues/236-*.md` | 7 | AC#6 三支摆明＋"`fs_broken.go` 已被入库"现量表第 4 行＋`:40` "不许用删掉 fs_broken.go 当甲的省事做法" | 票面【禁】 落地腿不许碰（归编排者） |
| `.scratch/wisp/issues/269-*.md` | 4 | 题面、现量#1、AC#1、AC#3 | 票面【禁】 同上 |
| `.scratch/wisp/probes/ci/1005-a/verdict.md` | 6 | §2 (A) 那格的逐字红句＋"点名文件：只有 1 枚"判语 | 在写面（`probes/**`），⚠ 但它是**别人的已提交取证件**（本仓规矩：不编辑别人已提交件）⇒ 只能另件追加更正 |
| `.scratch/wisp/probes/ci-red/ci-red-1.md` | 2 | 更早那轮同族取证 | 同上（别人件） |
| `.scratch/wisp/probes/236/a1/census.md`／`a2/census.md` | 9／13 | **AC#6a 的名册本体**（逐枚脏因＋归属＋引用面） | 同上（别人件） |
| `docs/evidence/s1/185-reread-owner-census-c1.md` | 3 | `:48` 引"路径＋错误文案"（防恒绿自证的凭据）、`:191` 引路径 | **【禁】 `docs/**` 不在 269 写面** |
| `docs/evidence/s1/197-subagent-entity-r1b.md` | 2 | 引用面普查表里的一行 | 【禁】 同上 |
| `docs/evidence/s1/200-project-instructions-r2.md` | 1 | 同上 | 【禁】 同上 |
| `docs/reports/HANDOVER.md` | 3 | 停车点叙述里点名该文件 | 【禁】 同上（且停车点只写最新时间戳那一节，归编排者） |
| `docs/reports/pending-and-issues.md` | 8 | `:8983`、`:9566`、`:9610`、`:9611`、`:9614`、`:10721`、`:12086`、`:12111` 逐字路径 | 【禁】 不在写面＋**只追加不删** ⇒ 只能追加更正 |
| `probes/orchestrator/ci-delta-1/delta.md`、`…/excerpt-verbatim.txt`、`…/job_lint.txt` | 3／1／1 | 名册差集与逐字摘录（编排者件） | 【禁】 编排者名下件，不由本腿改 |

**H＝其余（48 枚，含 23 枚未入库）**：各类 gofumpt 名单快照／stderr 快照／`tracked-go-zlist.bin`（NUL 分隔名册，**逐字节**含该路径）／CI 日志／overlay 输出快照 `probes/185/c1/logs/overlay-proof-broken.txt`／`probes/236/a2/diff-8-185-fs-broken.txt`（**文件名本身**带 `fs-broken`）／`ledger-restore.tmp`／commit-msg 草稿。判语：这些是"当时确实读到过这枚路径"的证据，**改名后它们仍然为真**（它们记录的是历史读数，不是当前盘上承诺），按 `issues/README` 规则 8 与台账规矩**一枚都不该改**；⚠ 但要在裁乙的批语里明写这一句，否则下一位会把"48 枚 H 档没改"读成"乙的引用面没改净"。

**phantom citation 的量化口径（三把尺，别混用）**：路径＋`:4:1` 钉在同一串里＝**40 枚文件（tracked 19 枚）**；引"错误文案"不引路径＝**57 枚（tracked 21 枚）**——后者改名后**仍为真**（同内容仍产同诊断）；只有前者会过期，且其中真需动手的＝R 档那 11 枚。

## §7 判不动／量不到（【禁】 不跑 `go` 造成的射程损失，逐枚具名）

1. **`:206`／`:209` 的内容**：禁 `go vet` ⇒ 那两步真跑起来是绿是红、几枚，**量不到**。今天它们零读数（`verdict.md:36` 现量 `skipped`）。⇒ 丙落地的第一发 run 才知道，那一格归编排者。
2. **"两把尺在同一次 run 里到底同不同意"**：我在**同一棵 HEAD 归档**上分别跑了 (A)（显式递 921／922 枚）与 (B)（目录走查），得**同一批**文件（§4 第 5 点、§1 现量#5）⇒ **内容上同意**。但"同一次 CI run 里两步都执行并各自出结论"从没发生过（`:184` 自建库以来从未跑过），所以**票面 `:15` 那句"今天它们不同意（一把红、一把没机会跑）"在"作答与否"这一维是对的，在"读数"这一维是不成立的**——两把尺没有互相打脸，只是一把被掐死。⇒ 定案需一次真 run（丙或手工推送），归编排者。
3. **`attrib.sh` 本体没跑**：它 `:137` 要 `go env GOPATH`（【禁】 不在派单允许清单），我只用等价命令复现了 `:341` 那一条管线。⇒ 未实测的有：根推导 `:127-135`、gofumpt 定位 `:137-140`、`--self-test` 那 8 枚双向 case、`--classify-line` 单发。这些都是"脚本自己会不会红"的路径，裁甲/丙之前应补跑一发（不需要 `go build`，只需 `sh`＋`go env`）。
4. **枚数是移动靶**：本腿工作期间名册 19→20（§5 末）、分母 921→922、HEAD 漂三次（`05680e44`→`50e0299c`→`65142d9a`）。⇒ 本文件任何枚数都**只在标注的锚点上成立**，落地腿必须现量重读；【禁】 不许把它当常量抄。
5. **gofumpt 版本**：CI 写 `@latest`（`ci.yml:170`）⇒ 上游改版会改"要格式化"集合。本腿与那枚 run 恰好同版（`v0.12.0`，凭据＝CI 日志 `go: downloading mvdan.cc/gofumpt v0.12.0`），读数才可迁移；下一发 run 不保证。
6. **行尾形状**：本机工作树被别的腿写过的部分带 CRLF（blob 是 LF，见 §0 末条）⇒ **本机直接跑 `gofumpt -l .` 的读数不能当 CI 读数**（我这发 47 枚 vs CI 形状 19/20 枚，差集 28 枚里 13 枚纯因行尾）。任何后续腿要复现门的读数，请用整树归档法（含 `go.mod`）。
7. **未读的界面**：`frontend/**`／`design/**` 依约未读未引（`lint-frontend` 那 11 步的守卫后果不属本票）；`docs/evidence/s1/**` 我只读了 `fs_broken` 命中的那 3 枚的相关行。
8. **裁丙那一格本腿不自裁**：三处里 `:184` 与它自己的注释（`:178-183`）重叠、`:206`／`:209` 不重叠——"三处一起补 vs 只补两处"是治理判断（A609③ 先例的口径），【禁】 不是事实问题，留给编排者（§2 丙④ 末）。

## §8 我不同意票面的地方（逐条，配凭据）

1. **题面的因果形状我顶回**：票面标题与"这是什么"把红写成"被**一枚**故意写坏的样本顶成常红"。现量：门的 tracked 分母今天脏于 **19→20 枚**（§5 尺①），其中 **1 枚根本不是样本**（`tools/d22scan/selftest.go`，`5e8748b3` 10-03 入库）、**1 枚是本腿期间新入库的负控**（`probes/259/r1/gofumpt-negctl/probe.go`，`50e0299c` 12:09）。⇒ **"把样本处理掉门就绿"这一预设不成立**，甲／乙单独都不行（§2 甲④／乙④）。这是本件最重要的一格：**本票的第二半（吞三步）与第一半（红）需要不同的形**。
2. **"日志点名只有 1 枚"是可见性假象，不是事实**：`ci-attr-1` §2 那句"点名文件：只有 1 枚"我复算了——它的日志确实只有 1 枚，但机理是 `OUT="$(gofumpt -l …)"` 把 stdout 名单**吞进变量**，而 `bash -e` 在 `:171` 以退码 2 当场掐死该步（凭据：`lint-full.log` 里那一步只有两行 `##[error]`，`echo "$OUT"` 一字未出现），所以 18/19 枚名单**从没被打印过**。⇒ 题面把"1 枚"当现量没错，但当"红因全集"就错了；这条同时给乙记一笔收益（§2 乙④(一)）。
3. **`:195-198` 的引用位置错**（应为 `:185-188` ＋ `:189-193`）；实质不变，见 §1 现量#3。
4. **"分层"清单少一枚**：实测带守卫的 step＝**9 枚**（`!cancelled()` 8 枚，含票面漏的 **`:558`**；`always()` 1 枚 `:376`），见 §3 表。方向上更支持票面结论。
5. **本票与票 236 AC#6 是同一枚缺陷的第二次立案，且 236 的旧结论今天过期**：台账 `:9566`（10-01）已把同一件事立进 **票 236 AC#6**（框 `[ ]` 未勾，现读 `236:31`），并在 `:9611`／`236:100-104` 把它重切成 **AC#6a（逐枚列名＋归属＋引用面）＋AC#6b（二选一，编排者当时推荐甲）**；`236:115` 还逐字记着「**甲法（分母排除台件）确实能让门变绿**」——那枚结论的名册是**当时 8 枚、全在 `.scratch`**（台账 `:9610` 现量）。今天 19→20 枚里有一枚在 `tools/` ⇒ **那句"甲能让门变绿"已过期**。⇒ 裁形时两票必须并，且注意**两票的形编号互相打架**：236 的"乙"＝把那批脏样整理干净，269 的"乙"＝样本搬家；236 的"甲"＝分母排除台件＝269 的"甲"。并案时请重命名字母，【禁】 别让下一位拿 236 的"乙"去执行 269 的"乙"。
6. **题面 ④ 的连带面我扩了**：票面只写"命中至少票 236 的票面"。现量 60 枚／tracked 37 枚／**硬连带 1 枚**（`overlay-proof.json`，票面没提）／**脚本 0 枚**／**超出写面的引用 5 枚**（`docs/**` 4 枚＋台账）——见 §6。"改名前先穷举，穷举不了就别改"这句按现量**能穷举**（我穷举完了），所以乙不是"不可判"，是"要同批扩写面到 `docs/**` 才可判"。
7. **"会被 Go 仪器看见"对 `go vet` 是过头的**（§2 乙④ 末）：`.scratch/**` 以 `.` 开头，Go 的包模式本就跳过，且 `probes/**` 下 14 枚 `go.mod` 自成一模块。⇒ 乙的收益只在格式门一侧；若裁乙的人以为它同时解 `go vet` 那一侧，那是错的期待（本条未跑 `go` 证实，见 §7 第 3/6 点）。
8. **票面 AC#2 的判据形状我复认，但它的凭据源要换**：AC#2 说"同时 `:184`/`:206`/`:209` 三步在同一次 run 里必须出日志（判据形状＝步级，不是 job 级）"——这一条正是 §7 第 2 点那格"只能真 run 一次才知道"的形状，**写成步级是对的**（六枚 run 回看：job 级全 `completed` 而步级 3 枚 `skipped`，`verdict.md:24`/`:36`）。本条**不是不同意，是加一句**：AC#2 的"出日志"必须要求**名单可见**，否则甲落地后 `:168` 仍红于 `selftest.go` 而日志可能仍只报退码（§8 第 2 条那个吞 stdout 的机理不会因为换分母而消失——除非同时把 `OUT` 打印出来，那是 `:172-175` 那三行本来就写好的、只是今天走不到）。

## §9 自我对抗（本腿真改了的东西，逐条给"原来怎么读、现在怎么读、尺是什么"）

1. **行尾尺用错过一把**：我第一版用 `grep -c $'\r' <file>` 数 CR，得到 `cmd/wisp/models.go`＝334，于是写下"blob 里就带 CRLF ⇒ 6 枚生产文件是 CI 真红因"。**错**——在该 shell 里那个模式退化成了"匹配每一行"，334 是**行数**。换成 `tr -dc '\r' | wc -c` 后：盘上 334／blob **0** ⇒ 改判"CRLF 是本机工作树污染，CI 形状是 LF，那 5 枚 `internal/**`＋`cmd/**` 不是 CI 红因"。凭据见 §0 末条与 §7 第 6 点。**这条如果不改，我会把甲的必要性写成"门红于 6 枚产码"，整个 §2 就废了。**
2. **归档尺漏抽 `go.mod`**：我第一次做 CI 形状复现用的是 `git archive HEAD -- '*.go'`（只抽 `.go`），名册读出 **18 枚**；差集里 `probes/263/v1/src/main.go` 明明 `md5sum` 与 blob 全等却被漏读。原因是同目录那枚 **30 字节 `go.mod`** 决定 gofumpt 的语言档，只抽 `.go` 就把它丢了。换成整树归档（`git archive HEAD | tar -x`，`find -type f` 对齐 `git ls-files`）后名册 **18→19**。这条改的是**甲的名册封顶值**（少一枚就会漏排一枚）。
3. **枚数没标锚点**：我先写"分母＝921 枚"当常数；HEAD 在我工作期间漂了三次（`05680e44`→`50e0299c`→`65142d9a`），复跑变 **922**、名册变 **20**。⇒ 全部改写成"锚点＋枚数"成对引用（§5 末条、§7 第 4 点），并按本仓定式在每一把尺旁边留命令。
4. **票面行号我先信后核**：`:195-198` 与那串 7 枚守卫行号我按票面抄过一次，`grep -n` 现读后**两处分岔**（`:186` 才是 NO REACH、守卫实测 8 枚含 `:558`）。⇒ §1 现量#2/#3 改成"逐字读数＋判定"两栏，并把"行号一律现读"写回 §0 总则（凭据同 `pending-and-issues.md:11905` 那条定式）。
5. **"1 枚红因"我一度接受**：ci-attr-1 的"点名文件只有 1 枚"我照抄过一版结论；自查机理（`OUT` 被命令替换吞＋`bash -e` 退 2 掐死）后推翻成"19/20 枚，日志只露解析错误那一枚"。⇒ §8 第 2 条。
6. **"185 三枚"我一度当三枚样本**：题面说 185 名下三枚。按尺①只有 `fs_broken.go` 被门读到，`fs.go`／`taintmatch.go` 排版干净（它们是 overlay 基准副本）；按尺③它俩确实是"故意坏的家族件"（内容自陈）。⇒ §5 末条把"三枚"具名拆成两档，【禁】 没把票面那句当错，只是**它答的是另一个问题**。
7. **乙我一度以为"能解 go vet 那一侧"**：核 Go 包模式的点开头目录跳过＋`probes/**` 自带 14 枚 `go.mod` 之后改成"只对格式门有效"，并如实标注这条**未跑 `go` 证实**（§7 第 3/6 点）。
8. **两族 `guard.no1-3.go` 我混成一族**：第一版 §5 把 `161/r3/pre-format/guard.no1-3.go`（**被 `probes/161/r3/.gitignore:1` 屏蔽**，非 tracked）当成票 161 `:62`／台账 `:9610` 说的那三枚，还顺手写了"台账曾按已入库记过、今天索引里不在"。现读推翻：那三枚讲的是 **`158/accept-r1/mut/guard.no1-3.go`，它们是 tracked**（同目录 9 枚在册），且已按 AC#7① 格式化干净 ⇒ 不在尺①。尺＝`git ls-files -- .scratch/wisp/probes/158/accept-r1/mut/`（9 枚）对 `git ls-files -o -i --exclude-standard -- .scratch/wisp/probes/161`（325 枚 ignored）。**这条如果不改，下一位会拿"161 的样本今天不入库"去反驳"台件不该以 `.go` 入库"那条规矩，而两句话其实各自为真。**
9. **"`gofumpt --version` 打不出来"我一度当成尺坏**：第一次运行 `"$GF" --version && echo …` 时版本行没出现在输出里，我以为这把尺不可信；重跑单发 `"$GF" --version 2>&1; echo rc=$?` 读出 `v0.12.0 (go1.27.1)`、rc＝0 ⇒ 是上一条命令的 stdout 被管道吞了，不是尺坏。这条与 §9 第 1 条同族（**工具没骗我，是我的读法骗我**），一起留在件里。

## §10 判语三行（只报形状与代价，不裁形）

1. **形状**：三形里只有丙解本票第二半（吞三步），只有甲碰本票第一半的分母，只有乙把"红得没有名册"改成"红得有名册"；**没有一枚单独能让 lint 变绿**——tracked 分母脏于 **19 枚（起手锚 `05680e44`）／20 枚（本腿期间 `50e0299c` 入库那枚 259 负控之后）**，其中 **1 枚（`tools/d22scan/selftest.go`）既不是样本、也不在 269 的写面里**（尺与名册＝§1 现量#5／§5 尺①；同一把尺在 run 的 `c6cf66e6` 上名册差集为空＝那枚 run 也是这 19 枚）。
2. **代价**：甲推翻 `ci.yml:189-193`＋与票 161 `:124`／`:141` 的裁定文本正面对撞（"给 `.scratch` 加豁免"正是 161 逐字点名"不是我的修法"那一支），并把链头红从 `:168` 挪到 `:184`＝**第二半没解**，还要背"名册每有一票入库就得追一枚"的移动靶（§5 末条现场证明）；乙不动 `ci.yml` 一字、不碰分母那句，但硬连带 1 枚（`overlay-proof.json`）＋phantom 候选 40 枚（tracked 19 枚，含 `docs/**` 5 枚**超出 269 写面**）＋门照旧红于 18/19 枚；丙只动 `ci.yml` 一枚（3 增 0 删／3 hunk）、零引用连带、AC#1 自动满足，代价是**门仍常红＝没人看**这一格原封不动，外加"vet 开始作答后可能新增红名"与"`go install` 失败时 `:184` 会以 exit 2 新增一种红"这两枚不可预见的可见性（§4 第 4 点／§7 第 1 点）。
3. **牙**：以"修完之后，格式真的坏了那一形会不会让门红"为尺——甲：会红，且第一次红在真违规上（但 `:184` 仍把三步吞着）；乙：仍红、但红得**有名册**（吞三步照旧）；丙：红不变、三步不再被吞（门对"格式坏没坏"仍零区分度）。⇒ **甲与乙作用于"红"，丙作用于"吞"；本票的两半不是同一枚门的两档，是两枚不同的病。**


