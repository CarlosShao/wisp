# 171-r3 实现程证据件 —— AC#1（尺洞②＝尾随注释被当成调用点，两族腿）＋ AC#9（flip 新格的两列／删样本不删腿）

本件＝票 171 的 **AC#1** 与 **AC#9** 两格的落地记录。写码腿代号 `171-r3`。
凡"响／不响"一律给当场跑出来的读数；凡推的、量不到的，逐枚具名写在 §7。

---

## 0. 锚点与口径

**取数时刻**：`2026-10-03 10:09 +0800`（`date` 自取）。
**开工锚点**：`git log -1 --format=%H` ＝ `f78d3cb183b041e6211d3955fdfb0f5e96098d6b`（分支 `dev`）。
**读数锚点（本件所有"真树"读数共用的一枚）**：`4f2a777b60770744adabbd92144d261812a025a0`
　　＝ 我第二次 `git rev-parse HEAD` 落下来的号（`2026-10-03 10:13:12 +0800`，`228-a5: fill section 3 and 4 ...`）。
　　**两枚号不同这件事本身就是读数**：同机另外两程在我取数期间各落了提交，HEAD 在我脚下前移过。
　　所以本件里"改前 vs 改后"全部走**显式锚点参数**（`sh <尺> 4f2a777b...`），不拿两次不同 HEAD 比——
　　这条规矩的出处＝票 171 票面 `Progress log` 09-27 18:0x 那节（`171-r2` 同一枚坑的处置）。

**口径三条**
1. **尺读的永远是提交树**（`git grep ... "$A"`），不是我工作副本里的未提交件。同机两程此刻在 `internal/config`、
   `cmd/wisp`、`internal/panel`、`internal/agent/approval` 里有未提交改动 ⇒ 它们**不会**洗掉我的读数，
   但**会**洗掉"下一次在 HEAD 上重跑"的可比性。凡是"今天响不响"的话，都带锚点号说。
2. **未修码的读数先落盘再动尺**。§2 的每一行都出自
   `.scratch/wisp/probes/171/r3/gate-prefix-f78d3cb.sh` ＝ `git show HEAD:.scratch/wisp/probes/154/gate-clauses.sh`
   的字节（577 行，与当时工作副本逐字节同：`cmp` 未报差异），它在修码之前就跑完并 tee 进
   `.scratch/wisp/probes/171/r3/logs/gate-pre-at-anchor.txt`。
3. **分子口径一次覆盖两族腿**（pair 与 run 各至少一发对照）——票面 AC#1 那格 09-27 19:1x 追加的提醒逐字要求这件事，
   出处还有 `docs/evidence/s1/171-run-legs-accept-v2.md:141`（"run 腿分子＝命中行数含注释 ⇒ AC#1 那枚洞今天开在 AC#7 这扇门上"）。

**写面自报（逐枚列路径，越界＝本件作废）**
- `.scratch/wisp/probes/154/gate-clauses.sh` —— AC#1 的尺本体（`pair()` 与 `run()` 的分子）。
- `.scratch/wisp/probes/161/r6/flip-declaration.sh` —— AC#9① 的"新格"本体（`:151`–`:178` 那一格）。
- `.scratch/wisp/probes/171/r3/**` —— 本程台件、快照尺、读数。
- 本件 `.scratch/wisp/probes/171/r3/impl.md`。

> **一处派单文字与票面的射程差，具名在此，不自行扩权也不自行作废**：
> 派单写"写面只有 `.scratch/wisp/probes/171/**`（尺本体与台件）"，而 AC#1 的尺本体在 `probes/154/`、
> AC#9① 那枚新格在 `probes/161/r6/`；票 171 票面 AC#5 的原文是"只许动 `probes/**` ＋ `ci.yml` 的新增步骤 ＋ 自己的证据件"，
> 派单首段也写着"这枚票的尺本体全在 `.scratch/wisp/probes/**` 与仓根脚本层"。
> 我按**票面**（AC#5 的 `probes/**`）做，只碰上面那两枚尺本体的**分子/判据文字**，一枚 pathspec、一条腿、
> 一发声明的语义都没减；若编排者判这两枚尺本体的写面不该动，回退面＝`git cat-file blob HEAD:<path> > <path>` 两枚文件，
> 台件与读数不受影响。**AC#3 需要碰的 `.github/workflows/ci.yml` 我一字节未碰**（见 §7）。

---

## 1. 尺本体在哪（目录级路径＋行段，只给文件名＝不可核）

| 件 | 绝对路径（仓内相对写法，仓根＝`D:\work\workspace\projects plans\Wisp`） | 本程关心的行段（改前字节） |
|---|---|---|
| 成对普查＋命中行数那把尺 | `.scratch/wisp/probes/154/gate-clauses.sh`（577 行） | `want()` `:73`–`:97`；`want_n()` `:103`–`:115`；`book_count()` `:130`–`:168`；`discarded_files()` `:179`–`:184`；`diffsets()` `:186`–`:204`；`run()` `:210`–`:228`；`pair()` `:230`–`:320`，**其中"剔整行注释"那一味＝`:248`**；十四腿声明 `:347`–`:510`；聚合段 `:517`–`:577`；【基线过期】表 `:563`–`:570` |
| 声明门禁（flip 那把） | `.scratch/wisp/probes/161/r6/flip-declaration.sh`（198 行） | baseline 格 `:126`–`:131`；九枚翻声明 `:133`–`:149`；**AC#8② 那枚新格 `:151`–`:178`，其中"核表"的三味 grep＝`:169`–`:171`**；restore `:180`–`:189` |
| 改前尺的字节快照（本程台件） | `.scratch/wisp/probes/171/r3/gate-prefix-f78d3cb.sh` | 与 `git show HEAD:...gate-clauses.sh` 同源，577 行 |
| AC#1 真树读数（改前） | `.scratch/wisp/probes/171/r3/logs/gate-pre-at-anchor.txt` | 939 行；聚合表 `:915`–`:930` |
| 台件先例（同一枚合成树手法，本程照其形状） | `.scratch/wisp/probes/171/r1/r1-pair-shapes.sh`（356 行，`:61`–`:109` 是丢弃树＋逐状态提交那段）；`.scratch/wisp/probes/171/r2/r2-run-legs-and-subtraction.sh`（278 行）；`.scratch/wisp/probes/171/r2/r2-roster-diff.sh`（91 行） | 合成树一律 `mktemp -d` 走 `$TMPDIR`、`git init`、**不在仓内建 worktree/checkout**（`r1-pair-shapes.sh:14`–`:20` 写了理由） |
| 验收程 `171-v2` 留下的三枚退化台件 | `.scratch/wisp/probes/171/v2/mutations/post-stale.sh`／`post-stale-lie.sh`／`post-stale-count0.sh`／`post-stale-wrongleg.sh` ＋ 同名 `res-3-*.txt` 读数 | AC#9① 的"那一列写错"那一发的**先例形状**＝`post-stale-lie.sh:157`（把 `实测=${n}枚 差=$((COUNT-n))枚` 篡成 `实测=${COUNT}枚 差=0枚`），其读数 `res-3-limit.txt:851` |

**尺为什么会把注释当调用点（读码定位，行号现量）**
- `pair()` 的 `:248` 那一味 `grep -vE ':[0-9]+:[[:space:]]*(//|/\*|\*)'` 只剔"**行首**是 `//`"的整行注释；
  行尾那半句它原样留给 `diffsets()`（`:186`–`:204` 里 `grep -Ew "$1"`／`grep -Ew "$2"` 吃的是**整行文本**），
  于是"有人在代码里提了一下开方/关范围"与"这里真调用了"在分子上不可分。
- `run()` 更宽：`:219`–`:221` 的分子 `hits="$(printf '%s' "$roster" | grep -c .)"` 吃的是**整张名册**，
  连整行注释都不剔（票 171 r2 自己在 `:345`–`:346` 写下了"run 腿的枚数吃整张名册，注释行也算——
  这与 pair 腿'只吃调用点'是两回事，票面 AC#1 那枚尺洞不归本程"）。⇒ 同一枚洞在两扇门上，两扇门今天都得修。

---

## 2. AC#1 未修码读数（两族腿，逐字来自盘上文件）

### 2.1 pair 那一族（真树，锚 `4f2a777b`，尺＝§0 的改前快照）

`G7 主尺` 的**反向**名册今天 3 枚（`.scratch/wisp/probes/171/r3/logs/gate-pre-at-anchor.txt`）：

```
#   UNPAIRED-REV internal/memory/retention.go (合方调用点=2)  # 同文件没见过开方：人来判「真漏／委托」
#   UNPAIRED-REV internal/observe/goroutine.go (合方调用点=1)  # 同文件没见过开方：人来判「真漏／委托」
#   UNPAIRED-REV internal/proc/shutdown.go (合方调用点=1)  # 同文件没见过开方：人来判「真漏／委托」
# 反向未成对枚数＝3   （0＝反向没响）
```

其中两枚的**全部**"合方在场"证据就是行尾注释（同文件里那两条名册行，逐字）：

```
4f2a777b...:internal/observe/goroutine.go:34:	CategoryTemporary GoroutineCategory = "temporary" // inside a DisposalScope
4f2a777b...:internal/proc/shutdown.go:85:	ReleaseSpeechSessions func(ctx context.Context) error // step 5 (speech, via DisposalScope)
```

⇒ **这两枚文件在"合方"那一侧是 0 枚真调用点**（`goroutine.go` 另一处 `:54`、`shutdown.go` 的 `:23`/`:161` 都是整行注释、已被 `:248` 剔掉）。
它们之所以被点名，就是因为 `pair()` 只剔整行、不剔行尾。⇒ **尺洞②在 pair 这扇门上的真树读数：3 枚点名里 2 枚是注释。**
这正好复算并顶齐了 `docs/evidence/s1/161-gates-accept-r2.md:67`/`:68` 那两行（"不是代码，是同行尾随注释"）。

第三枚 `internal/memory/retention.go`（合方调用点=2）**本票不动、也不冒充**：它那两行是
`:98 func (s *Store) StartRetentionJob(scope *plugin.DisposalScope, ...)`（真代码，签名）与
`:100 panic("memory: StartRetentionJob requires a DisposalScope")`（真代码，字符串字面量），
161-v2 判它"形状使然，且是判据词表缺一味（`Go` 不并成方）"（`161-gates-accept-r2.md:66`）。
⇒ 修完之后它**该照旧被点名**；我把这一枚当**正控**用（见 §3.3 的名册集合差），它要是退场了就是我在减射程。

### 2.2 run 那一族（真树，同一枚锚、同一把改前尺）

`G2`（`git grep -nE -w 'OpenTask|CloseTask'`，排 `*_test.go`、排 `internal/tools/bridge.go`）今天 3 行：

```
4f2a777b...:cmd/wisp/run.go:942:// WHY THIS HOOK OWNS THE CLOSE. CloseTask's own comment says the composition
4f2a777b...:cmd/wisp/run.go:960:		rt.bridge.CloseTask(taskID)
4f2a777b...:internal/tools/subagent_197.go:332:	// gap is the open question Q-56 recorded in bridge.go's CloseTask comment,
# 命中行数＝3   （票 171 AC#7：run 腿的枚数就是它——同一条腿里的第二枚从此躲不进「非空」）
# 基线枚数＝2
```

⇒ **3 枚"命中"里只有 `run.go:960` 一枚是代码**；两枚是 `//` 开头的整行注释。run 腿的分子把注释行**逐枚当命中**计，
并且这一味今天还连着判据：`G2` 今天之所以在聚合表里是一枚 `BAD`（见 §2.4），一半原因是他票新增的真调用点，
另一半原因是**注释行**——两者在"命中行数"这一个数里分不开。
⇒ **尺洞②在 run 这扇门上的真树读数：3 枚命中里 2 枚是注释。**

### 2.3 合成台件（两族腿各两发对照）

形状与判读口径＝§3.1 那节写的台件表；读数见 §3.2（台件是**未修码**上先跑的，跑完才动尺）。

### 2.4 顺带一枚与本票无关、但会影响下一程读数的现量（具名，不改它）

同一把改前尺、同一枚锚 ⇒ **聚合退码＝3**（`gate-pre-at-anchor.txt:928`、`:915`–`:924`）：

```
# BAD  腿=G1b 声明=quiet 基线=0枚 实测=1枚 因=空/非空那一味与声明不符（票 161 AC#6①）
# BAD  腿=G2 声明=ring 基线=2枚 实测=3枚 因=新增命中（票 171 AC#7：实测 > 基线）
# BAD  腿=G6neg 声明=ring 基线=1枚 实测=3枚 因=新增未成对（票 171 AC#2：实测 > 基线）
```

三枚 BAD 的来路＝同机另外两程（票 255／票 220／票 228 一族）在我取数期间落地的**新调用点**，
不是本票的判据、也不是我把谁放宽了——**本程没有为它们改任何 `want`／`want_n` 之外的东西，
也没有拿"门今天本来就红"当自己那两格的凭据**。它同时意味着：`flip-declaration.sh` 的 baseline 那一格
今天在这枚锚上是 FAIL（它期望 rc=0），而 `want_n 8 -> 9` 那一味（`:160`）能不能命中要现量。
⇒ 这条账落在 §4.3 与 §7，处置＝我把新格改成**按腿号定位声明行**、不再抄某一枚硬写的数字（理由与读数见 §4）。

---

## 3. 修法与修后读数

（本节 §3.1 台件表、§3.2 未修码读数、§3.3 修后读数与名册集合差、§3.4 十四腿基线复算——
台件在未修码上跑完之后才动尺，动尺之后的读数与集合差逐枚贴在这里。）

## 4. AC#9①：flip 新格连"实测／差"两列一起核

（含"两列造假但仍声明存在"那一发的**未修格 GREEN／修后格 RED** 两向读数。）

## 5. 门禁与卫生（本程允许的面）

本程**一字节未碰**：`.github/workflows/ci.yml`、`internal/**`、`cmd/**`、`tools/d22scan/**`（含 `allowlist.txt`）、
`docs/PLAN.md`、`docs/specs/**`、`docs/BUILD.md`、`docs/SLO.md`、`internal/observe/thresholds.go`、golden、三枚冻结件。
`frontend/**` 与 `design/**` 本程未读、未改，本件也不引用其内容。
Go 编译类门禁（`go build`／`go test`／走 `go` 的 `scripts/d22scan.sh`）今天按派单禁令**不跑**，逐枚具名在 §7。
文本级卫生（不依赖 Go 的那几味）由我自己跑，读数贴在本节。

## 6. 我推翻前人（票 171 票面与 171-v2 那张表）哪几句

（逐句给"前人的原话／我的现量／差在哪一枚号"。）

## 7. 量不到／判不动的地方

- `scripts/d22scan.sh` 与 `tools/d22scan/runtests.sh`：**量不到**。理由逐字＝派单硬约束"不许跑 `go build`／`go vet`／`go test`（任何包）"，
  而 `scripts/d22scan.sh:1`–`:30` 自述它两步都要 `go`（`runtests.sh -C tools/d22scan ./...` 与 `go run . -root`）。
  ⇒ 本程用不依赖 Go 的文本尺代跑 emoji／裸 `go func(`／明文密钥那几味（读数在 §5），
  但**这不等于 d22scan 门过了**，那一格留给能跑 Go 的程。
- `.github/workflows/ci.yml`：本程未碰（AC#3 归编排者另裁），因此 AC#1/AC#9 的修法**没有**"非动 ci.yml 不可"那一说——
  两格的尺都是 shell/文本尺，跑它不需要 Go、也不需要 CI 步骤。
- "下一次在 HEAD 上重跑本尺会得到什么数"：**判不动**。同机两程正在往 `internal/**`、`cmd/**` 落码，
  §2.4 那三枚 BAD 就是他们造成的；本件的十四腿基线只在与 §0 同一枚锚上成立。
