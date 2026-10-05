# census.md — 票 269 的只读普查（腿 `269-a1`，三形各四问，⛔ 本腿不裁形）

交付件＝本文件。写面＝仅这一枚。工单票面 AC 框与 `docs/reports/pending-and-issues.md` 一字未动。

## §0 起手锚（我自己现读，不抄票面）

- 起手现量时刻：`2026-10-05 11:52:00 +0800`（`date`）。落笔时刻：`2026-10-05 12:1x +0800`（同一条命令内两次 `date` 之差＝本腿全程）。
- 起手 HEAD：`05680e44`（`git log -1 --format=%h`）；分支 `dev`。**写到一半 HEAD 漂到 `50e29c`→`50e0299c`**（别的腿入库），故 §1 每一把都标了读数时刻，并在 `12:0x` 复认 `ci.yml` 那 7 个行号在新 HEAD 上仍成立（`git diff --numstat -- .github/workflows/ci.yml`＝空、该文件最后一笔＝`519e1ade` 10-04 09:42）。
- 工作树脏度（⛔ 不是我的）：起手 `git status --porcelain` 前 40 行即含 `internal/agent/approval/approval.go`／`queue.go`／`ticket242_binding_test.go` 三枚 M＝在飞写腿 `259-r1` 的包；`git diff --name-only HEAD | wc -l` 现量 **32**。
- 【禁】遵守情况：全程零 `go build`／`go vet`／`go test`／`go run`；用过的仪器＝`grep`／`sed`／`awk`／`git ls-files`／`git ls-tree`／`git log`／`git cat-file`／`git archive`（只读导出到 `/tmp/269a1/`，⛔ 未在仓库目录内建 worktree 或 checkout）＋**`gofumpt` 本体**（`D:\work\base\gopath\bin\gofumpt.exe`，`--version` 现读 `v0.12.0 (go1.27.1)`）。派单允许清单里没有 `go env`，所以我**没有真跑** `attrib.sh`（它第 137 行要调 `go env GOPATH`），而是用等价命令复现它那一条管线（见 §4）。
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

## §2 三形各答四问（⛔ 不选形、不裁形；裁形归编排者）

### 甲＝给 `ci.yml:168` 那一步换分母（只喂"tracked 集合里非样本"的 `.go`）

**①要动哪几枚文件（逐枚路径）**

1. `.github/workflows/ci.yml` —— 只改 `:169-176` 那一步的 `run:` 体（把 `gofumpt -l . tools/d22scan tools/mockllm` 换成"tracked 全集剔掉具名样本名册"）。⚠ **这一步是"既有步骤"，不是新增步骤**：票 171 `:31` 与票面写面句（`269:32-33`）对"动既有步骤"给的是不同的界，见 ② 末与 §7。
2. 新增 `scripts/gofumpt-tracked-gate.sh`（分母计算＋名册守卫＋自报读数；本仓"门自己数自己的分母"的既有形状＝`scripts/check-path-length-budget.sh:298` 打 `denominator read: %s tracked paths`）。
3. 新增名册件（或把名册内联进上面那枚脚本的数据段，照 `check-path-length-budget.sh:170`「Exemptions are DATA inside this file, not comments」＋`:314`／`:372` 的 `#ROSTER-AWK-BEGIN`／`#ROSTER-AWK-END` 标记定式）。
4. `ci.yml:178-204` 那步（形 (A)）**不动**——它按票 161 的裁定必须继续吃"tracked 全集"，否则 ② 那两句自陈被改动两次。

**"排除怎么写才不会顺手把真违规也排除掉"——具名形状（五条，缺一条就漏）**

- (一)**逐枚全路径、禁目录通配**。名册条目只能是 `.scratch/wisp/probes/<NN>/<...>/x.go` 这种完整相对路径；写成 `probes/**`／`mut/**`／`*.go` 一类模式＝把以后落在同一目录里的真伤一起排掉（本仓同款病的既有定式＝`check-path-length-budget.sh:124-125` 的 A／B 双向守卫）。
- (二)**准入条件＝两条同时成立**：路径必须以 `.scratch/wisp/probes/` 开头，**且** `.scratch/wisp/issues/<NN>-*.md` 真存在。这条直接抄 `attrib.sh:174-193`（`ticket_of`＋`ticket_known`），它的注释逐字理由：「The first half alone would be self-fulfilling (any number passes)」。
- (三)**四形守卫**（照 `check-path-length-budget.sh:120-131` 逐条搬）：A 任何未入册的脏样 ⇒ 红；B 名册里一条今天不再是脏样／不在盘上 ⇒ 红（防过期豁免＝防 phantom exemption）；C `count(名册条数) != count(tracked 脏样条数)` ⇒ 红（连重复键一起挡，原文 `:128-130`）；D 理由为空或全空白 ⇒ 红。**B 与 C 是"只排样本"与"排掉一切"之间唯一的机械差别**，没有 B／C 的名册一定烂成豁免清单。
- (四)**分母读数照常打印**：剔样本之前先把 `git ls-files '*.go' | wc -l`（现量 **921**）打进输出，并保留空心守卫（`attrib.sh:337-338`：递给尺 0 枚 ⇒ `exit 2`，CI 读成红、永不读成绿）。剔的是"被判定的集合"，不是"被数出来的集合"。
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

1. `.scratch/wisp/probes/185/c1/mut/fs_broken.go` → 同目录改后缀（`.go.pristine`／`.go.txt` 一类）。⛔ 不许"搬出 `.scratch`"，更不许删（票面 AC#1＋`issues/README` 规则 8「临时件只建不删」）。先例尺：`git ls-files | grep -c '\.pristine'`＝**17 枚在册**（票面点名的 `267/r2/pristine/models.go.pristine` 是其中一枚；同族＝`248/r1` 3 枚、`248/v1b` 5 枚、`252/v1` 3 枚 `-pristine-v1`、`256/v1` 2 枚、`33/r1` 1 枚、`33/r9` 2 枚）⇒ 先例形状＝**只改后缀、不动字节**。
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
- ⛔ 不许"顺手"给 `:168`／`:136` 之类补守卫（那是第四、第五处，越出票面对丙的定义）。

**②会不会推翻"分母＝tracked 集合"那两句＝不会（待验断言"丙不会"成立）**

- 丙不碰任何一步的命令体、不碰任何一步的入参集合；`:186-188` 的"两把尺读同一批字节、必须同意"与 `:189-193` 的"分母＝TRACKED 集合"逐字照旧为真。
- ⚠ 但丙**改变那两句的今天状态**：`:186` 那句 `adds NO REACH TODAY` 的前提是"`:184` 与 `:168` 同生同死"。丙之后 `:184` 在 `:168` 红时**也会作答** ⇒ 那一步第一次有独立射程（它比 `:168` 多一条空心守卫、多一层 `exit 2` 语义，见 §4）。⇒ 严格说：**丙不推翻那两句，但让 `:186` 那句"今天不增加射程"变成过期描述**（那是自陈而非裁定；改不改注释是可选项，不是必需项，本腿不替落地腿决定）。

**③会不会咬到引用那枚样本路径的件＝一枚都不咬**

- 样本原地不动（`fs_broken.go` 字节未变，`git rev-parse HEAD:...`＝`ee583853...` 现读）⇒ §6 那 60 枚引用件、`overlay-proof.json` 那枚硬连带、台账 8 处全部零风险；AC#1 自动满足，而且是**唯一一枚不需要任何引用治理就能满足 AC#1 的形**。
- 唯一新增的文字面风险：新插的三行若配注释去写样本路径，就把那条路径第二次落地。⇒ 形状要求：注释只指票号（`ticket 269`），⛔ 不写行号（本仓定式凭据＝`docs/reports/pending-and-issues.md:11905` 那句「凡引用代码位置，`grep` 被指的那串字符，别抄任何来源的行号（含 CI 日志）」）。

**④门剩什么牙（判语，不是偏好）**

- **丙单独做＝门仍常红＝等价于没人看**：`:168` 照旧红（19 枚脏样）、`:184` 照旧红（`exit 2`），lint 的 job 级颜色一字不变。⇒ 以"格式有没有坏"为命题的那扇门，在丙之后**仍然每天给出同一个不区分任何东西的红**——红已不含信息量：任何人看它都得到同一答案，因此没有任何人会去看它。票面 `269:20` 那句"丙单独做＝门仍常红＝等价于没人看"按现量**复认成立**。
- 丙真正买到的是**另外三枚门的读数**，且这是本票第二半唯一的解：`:184` 从此每次给出 (A) 尺的红因名册；`:206`／`:209` 从此每次真跑 `go vet ./...` 并出名册——今天这两枚自 09-28（`4813567e`）以来**零枚读数**（凭据：`.scratch/wisp/probes/ci/1005-a/verdict.md:36` 现量三步 conclusion＝`skipped`；台账 `:9566` 逐字「**后果不是"多一枚红"，是它把后面的 `go vet (module)` 与 `go vet (tools/d22scan module)` 一起吃掉 ⇒ 今晚两发零 `go vet` 读数**」）。
- 代价必须同批报：**"门开始作答"会把新的红带进来**。`go vet` 一旦响，lint 的红名册从"1 枚格式门"变成"N 枚 vet 名"，其中可能含从未被 CI 看见过的包（本腿禁 `go vet` ⇒ 几枚、会不会红，**量不到**，具名归 §7）。这不是回归，是新增可见性；但在"CI 得变绿"的压力下它正是最容易被读成"丙把 CI 弄红了"的形状 ⇒ 裁丙的人要把这句话提前写进批语。
- 逐字规矩核（细节在 §3）：文件头 `:5-7` 禁的是 `if: false`／`continue-on-error`／skip flags **三形**，`!cancelled()` 不在其列；`:412-416` 是这枚文件对同一形状的**自我定性**（逐字：`# WHY THIS IS NOT A DOOR REMOVED (the two forbidden shapes, both absent): no` ／ `# step was deleted, and NO step here carries `continue-on-error`, whose entire` ／ `# documented meaning is "prevents a job from failing when a step fails". A step` ／ `# run under `!cancelled()` that exits non-zero still fails the job exactly as` ／ `# before - the only change is that the steps after it also get to answer.`）；`:246-249` 逐字：``# `if: ${{ !cancelled() }}` is the one shape`` ／ `# ticket 85's ruling allows for an observation period: it does not` ／ `# change this step's own red/green, it only makes sure the step` ／ `# answers. continue-on-error is NOT used (D22 run-away mode 6,`。**这三处都是"逐字那行"，不是"别的步也这么写了"。**
- ⚠ 同一枚文件里还有**两句反向的、逐字的批准面**，裁丙的人必须处置，本腿不自作主张：
  - `:178-180`（票 161 给那一票划的界）逐字：`# ADDED STEP ONLY (ticket 161 AC#7 form (A), decided by the orchestrator's` ／ ``# 00:0x ruling "裁二"): no existing step was moved, edited, deleted, given an`` ／ ``# `if:`, or made continue-on-error.``。读法分歧：这句约束的是"161 那次编辑"，还是"这棵树上的既有步骤从此不许给 `if:`"？按字面是前者；`:263`／`:300`（票 85／111 后来确实给既有步加了 `if:`）说明仓里按前者执行过。
  - 票 171 `:31`（更宽的一句）逐字：`⚠ 动触发表／并发组／既有步骤＝**契约级（票 134 地界）⇒ 停手报回**`。若读成**仓级规矩**，丙（改三枚既有步骤）＝契约级 ⇒ 停手报回；若读成 171 自家写面，丙只需 269 票面 `:33` 的写面授权。
  - 台账 `:11905`（A609 §3 第③行，编排者自己的实操先例）逐字：「③ 给那一步加 `if:` 绕过 skipped＝动 `ci.yml:136` 那行注释自陈的 **A585 批准面**（"no existing step was moved, edited, deleted, given an `if:`"），⛔ 我不越。」⇒ 按这条先例的**口径**逐枚判：丙的三处里只有 **`:184`** 的注释（`:178-183`）含那句话，**`:206`／`:209` 两步体内没有任何自陈**⇒ 重叠面＝1/3。"三处一起补 vs 只补不重叠的两处"这一格归编排者裁。

## §10 判语三行（只报形状与代价，不裁形）

1. **形状**：三形里只有丙解本票第二半（吞三步），只有甲碰本票第一半的分母，只有乙把"红得没有名册"改成"红得有名册"；**没有一枚单独能让 lint 变绿**——tracked 分母今天脏于 **19 枚**，其中 **1 枚（`tools/d22scan/selftest.go`）既不是样本、也不在 269 的写面里**（尺与名册＝§1 现量#5／§5 尺①；同一把尺在 run 的 `c6cf66e6` 上名册差集为空）。
2. **代价**：甲推翻 `:189-193`＋与票 161 `:124`／`:141` 的裁定文本正面对撞（"给 `.scratch` 加豁免"正是 161 逐字点名"不是我的修法"那一支），并把链头红从 `:168` 挪到 `:184`＝**第二半没解**；乙不动 `ci.yml` 一字、不碰分母那句，但硬连带 1 枚（`overlay-proof.json`）＋phantom 候选 40 枚（tracked 19 枚，含 `docs/**` 5 枚**超出 269 写面**）＋门照旧红；丙只动 `ci.yml` 一枚（3 增 0 删／3 hunk）、零引用连带、AC#1 自动满足，代价是**门仍常红＝没人看**原封不动，外加"vet 开始作答后可能新增红名"这枚不可预见的可见性（§7 量不到）。
3. **牙**：以"修完之后，格式真的坏了那一形会不会让门红"为尺——甲：会红，且第一次红在真违规上（但 `:184` 仍把三步吞着）；乙：仍红于 18 枚、且红得**有名册**（吞三步照旧）；丙：红不变、三步不再被吞（门对"格式坏没坏"仍零区分度）。⇒ **甲与乙作用于"红"，丙作用于"吞"；本票的两半不是同一枚门的两档，是两枚不同的病。**
