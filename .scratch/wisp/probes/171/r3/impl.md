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
并且这一味今天还连着判据：`G2` 今天之所以在聚合表里是一枚 `BAD`（见 §2.4），我第一版把它写成"一半是他票新增的
真调用点"——**那句错了，逐字推翻在 §6.2**：那一枚"新增"里有一行是散文，门在指控一枚不存在的东西。
两者在"命中行数"这一个数里本来分不开，这正是尺洞②的形状。
⇒ **尺洞②在 run 这扇门上的真树读数：3 枚命中里 2 枚是注释。**

### 2.3 合成台件（两族腿各两发对照，未修码读数）

台件本体＝`.scratch/wisp/probes/171/r3/fixtures/s1_*.go.txt`／`s2_*.go.txt`／`s3_*.go.txt`／`s4_*.go.txt`
（落盘写成 `.go.txt`：它们不是任何包的一部分、不进 `internal/**`／`cmd/**` 那两条射程，
由台件脚本 `r3-trailing-comment.sh` 拷进 `$TMPDIR` 的一棵丢弃树里再提交，尺才吃得到——
形状照 `probes/171/r1/r1-pair-shapes.sh:14`–`:20`，仓内不建 worktree、不 checkout、不删任何件）。
跑法与逐枚读数＝`.scratch/wisp/probes/171/r3/logs/harness-pre-mode.txt`（8 条 ok，rc=0）＋
每态每尺的整份读数 `.scratch/wisp/probes/171/r3/logs/gate-{pre,post}-S[1-4].txt`。

**未修码（尺＝§0 那份快照）的读数，两族腿同一次跑出来**：

| 台件 | 形状 | pair 腿 `G6 未成对枚数` | run 腿 `G2 命中行数` |
|---|---|---|---|
| S1 `s1_only_trailing_comment.go` | 词表只出现在**一行真代码的行尾注释**里 | **1 枚＝被计入**（点名 `internal/fake171/s1_only_trailing_comment.go`） | **1 行＝被计入** |
| S2 `s2_real_call.go` | 真调用点 `rt.OpenTask(taskID)`，全文件没合方 | 1 枚＝被计入（该计入） | 1 行＝被计入（该计入） |
| S3 `s3_fake_close_comment.go` | 真开方调用＋**合方只写进行尾注释** | **0 枚＝漏被抹平**（未点名） | 2 行（其中 1 行是注释） |
| S4 `s4_full_line_comment.go` | 词表只出现在**整行注释**里 | 0 枚（改前 pair 已剔这一形） | **1 行＝被计入**（改前 run 不剔） |

⇒ 票面 AC#1 那一格说的"未修码上第一发今天就是响的"**逐字成立**，而且是两扇门同时成立：
pair 腿（S1 被计入 1 枚、S3 的活漏今天静默）与 run 腿（S1/S4 各被计入 1 行）。
S2 与 S4 是这两行的对照：真调用点两把尺都计入，整行注释两把尺**只有一把**剔得掉——
那正是"分子口径要一次覆盖两族腿"的证据形状。

（顺带一枚本程自己犯的、当场纠掉的仪器错：S2 第一版在 `type taskRouter interface { OpenTask(string) error }`
里也写了词根，于是"真调用点"那一发的 `命中行数` 未修码上是 2 不是 1——那是把**声明**当调用点，
与 `discarded_files()` 剔定义行是同一件事（`gate-clauses.sh:175`–`:177` 的自纠）。现版台件只留一行代码词根，
`probes/171/r3/logs/gate-pre-S2.txt` 的 1 是修好台件之后的读数。）

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

### 3.1 改的是分子，不是射程（尺本体 `probes/154/gate-clauses.sh`，逐处行号）

| 处 | 改前 | 改后 | 干了什么 |
|---|---|---|---|
| `strip_comment()`／`call_lines()`／`comment_only_lines()` | —（不存在） | 现 `:186`–`:222`（新加三枚 helper） | 两族腿共用的"只算代码那一半"：① 剔整行注释（沿用改前 pair 那一味，逐字）② 剔行尾注释（新）③ 剔"切完不再匹配词表"的行。第三味缺位＝只靠注释进场的行仍留在调用行里，等于没修 |
| `pair()` 的 `calls=` | 改前 `:248`：`grep -vE ':[0-9]+:[[:space:]]*(//|/\*|\*)'`（只剔整行） | 现 `:322`：`calls="$(call_lines "$roster" "$open|$close")"`，另 `:323`–`:333` 打 `#   COMMENT-ONLY` 逐枚被剔行 | pair 族分子 |
| `run()` 的 `hits=` | 改前 `:221`：`hits="$(printf '%s' "$roster" | grep -c .)"`（＝整张名册行数，注释全算） | 现 `:252`–`:269`：分子走 `call_lines ... -E`，并把 `# 名册行数＝`（未动）与 `# 命中行数＝`（分子）分成两格打，附 `#   COMMENT-ONLY` 表 | run 族分子；名册本体照旧逐行原样打印，下一程仍按行 diff |
| `--diffsets` 文本喂入那一味 | 改前 `:329`（自己抄了一遍"只剔整行注释"） | 现 `:394`（改吃 `call_lines`） | 第三处同源拷贝收掉（"同源拷贝逐枚有归属"那起事故的形状） |
| `want_n` 复算 | `:359`=2／`:490`=3／`:508`=4 | 同一格改成 1／1／2，各附现量命令与退场理由 | 见 §3.4 |

匹配那一味与每一腿自己捞名册那一味一致：pair 腿 `-Ew`（与 `diffsets()` 同味）、run 腿 `-E`
（与它自己的 `git grep -nE` 同味）——`call_lines` 的第三参就是为这件事留的，不然"命中"与"仍匹配"是两套词法。
按"空白字符＋`//`"切而不按"第一个 `//`"切：gofmt 之后行尾注释前必有一个空格，字符串里的 `https://x` 前面是冒号——
URL 不会被当成注释切掉（这是 S1/S3 那种形状能被稳定切开的理由，不是设想）。

### 3.2 修后台件读数（同一棵树、同一条射程、只有尺不同）

`.scratch/wisp/probes/171/r3/logs/harness-post-mode.txt`（16 条 ok，rc=0）：

| 台件 | pair `G6 未成对`：未修 → 修后 | run `G2 命中行数`：未修 → 修后 | 判读 |
|---|---|---|---|
| S1 只有行尾注释 | 1 → **0** | 1 → **0** | 不计入＝AC#1 判据①成立（两族都成立） |
| S2 真调用点 | 1 → **1** | 1 → **1** | 照旧计入＝判据②成立（没把代码一起剔掉） |
| S3 合方只在注释里 | 0 → **1** | 2 → **1** | 活漏从静默变点名＝"今天不响、修了才响"那一发 |
| S4 整行注释 | 0 → 0 | 1 → **0** | run 族补齐；pair 族形状未动（没顺手改它） |

### 3.3 射程零动、名册零动，拿集合差说话

台件＝`.scratch/wisp/probes/171/r3/r3-roster-diff.sh`，读数＝`logs/harness-roster-diff.txt`（rc=0）。
两把尺（`git show f78d3cb:...gate-clauses.sh` 的字节 vs 工作副本）打**同一枚锚** `4f2a777b`，且脚本先核
"`git show` 出来的字节与本程快照 `gate-prefix-f78d3cb.sh` 相同"，不同即 rc=2 不给读数：

- **check 1** 十四枚腿的名册本体两版**逐字节相同**（逐腿行数 0/1/3/0/0/94/91/295/3/6/34/19/28/71）。
  名册是射程的产物——它一字不动 ⇒ 射程一字未动。
- **check 2** `git diff f78d3cb -- <尺>` 里含 `:!`／`$GO`／`$GO2` 的增删行（注释行除外）＝**0 枚**；
  另附两枚防空转正控（本仓规矩："零命中"先怀疑尺）：同一味词表在尺自己身上打得出 **12 行**（词表不是哑弹）、
  两版尺的 diff 本身 **92 行**增删（不是同一份字节）。这两枚是我在第一版跑出 `grep: Unmatched (` 之后补的——
  那一味当时把"0 枚射程改动"读成了一条**假绿**，正则报错也算不出数；现版报错消失、正控非零。
- **check 3** 点名句单（`UNPAIRED`／`UNPAIRED-REV`）的**新增点名＝0 枚**；退场＝2 枚，逐枚有解释：
  `internal/observe/goroutine.go`、`internal/proc/shutdown.go`（G7 主尺与 G7-负一负各一次）。
  判据是硬的：每一枚位移都必须出现在**同一腿自己打出的** `#   COMMENT-ONLY` 格里，位移而无解释＝射程动了＝本程判失败。
- **check 4** 五枚 run 腿的"名册行数"两版相同（0/1/3/0/0），修后的分子＝0/1/1/0/0。

与 161-v2 对上的地方：它裁 G7 反向 3 枚里"2 枚不是代码，是同行尾随注释"（`161-gates-accept-r2.md:67`/`:68`），
我这把尺改完退场的正是那两枚；第三枚 `internal/memory/retention.go`（`:98` 签名＋`:100` 字符串）**照旧点名、
合方调用点仍＝2**——它是本程的正控：AC#1 修的是"注释算不算调用点"，不是"谁该被点名"。

### 3.4 十四腿基线复算（改口径同一次做，别留 SHR 当默认）

修后尺打同一枚锚，整份读数＝`logs/gate-post-rebaselined.txt`：

```
# BAD  腿=G1b 声明=quiet 基线=0枚 实测=1枚 因=空/非空那一味与声明不符（票 161 AC#6①）
# BAD  腿=G6neg 声明=ring 基线=1枚 实测=3枚 因=新增未成对（票 171 AC#2：实测 > 基线）
# 腿数＝14 声明与实测不符＝2      # 腿数断言＝相符（名册=14 声明=14 记账=14 缺腿=0 空头声明=0）
# 基线过期枚数＝0                  # 聚合退码＝2
```

- **本程核下来三枚**（改的就是它们的分子）：`G2 2→1`、`G7 3→1`、`G7neg 4→2`，每枚旁边写了现量命令与退场理由。
  `ring/quiet` 声明一枚未换位（三枚修后实测都 >0）。
- **本程拒绝动两枚**：`G1b（0→1）`、`G6neg（1→3）`。它们是新落地代码的真读数，改前尺在同一枚锚上同样报这两枚
  BAD（`logs/gate-pre-at-anchor.txt:915`、`:924`）；抬进基线＝替别程的发现把门做绿，正撞那把门自己文件头那句
  "a committed declaration no longer matches today's reading (a finding, **NOT a knob**: do not fix it by relaxing
  the declaration)"。⇒ 归口见 §7。
- 其余九枚（G1／G3／G4／G5／G5pos／G5neg／G6／G6pos／G7pos）实测与基线**逐枚全等**，一枚未动。
  `G5` 那枚在册活形状（`cmd/wisp/panel_assets.go`）照旧点名，`G5neg` 照旧 8 枚——flip 那一格的 `want_n 8`
  今天仍然落在 G5neg 上，是运气不是设计，故 §4.3 把它改成按腿查号。

---

## 4. AC#9：①那两列的核与对照发 ＋ ②删样本不删腿

### 4.1 未修格读数：两列造假、仍声明存在 ⇒ 今天过关（③-limit 复现）

造假尺＝`.scratch/wisp/probes/171/r3/fixtures/gate-lie-columns.sh`（由 `r3-flip-columns.sh` 从尺的 STALE 打印式
那一行生成，只改那一行：`实测=${n}枚 差=$((COUNT-n))枚` 篡成 `实测=${COUNT}枚 差=0枚`；`diff` 行数≠2 即 rc=2 不给读数）。
未修的门禁＝`git show HEAD:.scratch/wisp/probes/161/r6/flip-declaration.sh` 的字节，快照落盘
`.scratch/wisp/probes/171/r3/flip-prefix-4f2a777b.sh`；把它指到造假尺上只 patch 两行（`gate_rel=` 与 `anchor=`，
两行都核"补丁必须落得下来"），整份读数＝`logs/lie-door-unmodified.txt`，汇总＝`logs/harness-flip-columns-pre.txt`（rc=0）：

```
   the baseline line this run differs by:
   < want_n 8
   > want_n 9
  FAIL  G5neg 基线 8->9（实测 8）still exits 0 (变好了不算言行不一) -> rc=3 (want zero)
  ok    the stale baseline is its own readable table row: # STALE 腿=G5neg 声明=ring 基线=9枚 实测=9枚 差=0枚 ...
row-claim=GREEN  columns-claim=NO-CLAIM  whole-cell verdict=RED  whole-door rc=1
```

读法两条，都不含糊：
1. **那一行造假照样过关**：`基线=9枚` 与 `实测=9枚`、`差=0枚` 三列互不相容（真读数 8 枚），
   未修格只看"表在不在、点不点 G5neg、枚数是不是 1"⇒ 它给 `ok`；`columns-claim=NO-CLAIM`
   ＝那一格对这两列**根本没有任何断言**。⇒ ③-limit 今天成立，与 `171-run-legs-accept-v2.md:142` 同一数。
2. 同一次跑里那条 `FAIL ... still exits 0 -> rc=3` **不是本格的判据**，它是 §2.4 那三枚他票 BAD 把
   "期望退码 0"打到红上的下游。所以这台件按 label 逐条取判据（`row-claim`／`columns-claim`），
   不按整把门的 rc——否则我会拿别程的脏读数冒充自己那格的颜色。

### 4.2 修后：同一枚造假尺必须不 GREEN，诚实尺必须照旧 GREEN

`.scratch/wisp/probes/171/r3/logs/harness-flip-columns-post.txt`（三发同一枚锚 `4f2a777b`，rc=0）。

**第 (2) 发——修后的门打同一枚造假尺（"两列造假但仍声明存在"）**：

```
  FAIL  G5neg 基线 8->9 still exits 0 (变好了不算言行不一) -> rc=2 (want zero)
  ok    the stale baseline is its own readable table row: # STALE 腿=G5neg 声明=ring 基线=9枚 实测=9枚 差=0枚 ...
  FAIL  the row prints its own numbers - row got [# STALE 腿=G5neg 声明=ring 基线=9枚 实测=9枚 差=0枚 ...] want [基线=9枚 实测=8枚 差=1枚]
row-claim=GREEN  columns-claim=RED  whole-door rc=1
```

⇒ 那一格从"表在不在"升级成"表说的是不是真话"：**那一列写错时新格确实红**（票面 AC#9① 那句判据逐字兑现），
而且红句把 got／want 两串都打出来，下一程不必猜是谁在骗谁。
同一次跑里第 (1) 发照旧 `row-claim=GREEN columns-claim=NO-CLAIM`（未修的门禁对同一枚造假尺一无所觉）＝
两向对照齐了，不是"我改完就不响"那种单脚读数。

**第 (3) 发——修后的门打诚实尺（正控，防"又多印一行字"退化成恒红）**：

```
  ok    the stale baseline is its own readable table row: # STALE 腿=G5neg 声明=ring 基线=9枚 实测=8枚 差=1枚 ...
  ok    the row prints its own numbers: # STALE 腿=G5neg 声明=ring 基线=9枚 实测=8枚 差=1枚 ...
row-claim=GREEN  columns-claim=GREEN  whole-door rc=1
```

⇒ 新那一味不咬诚实行；这一发的 `whole-door rc=1` 仍然全部来自 §2.4 那两枚他票 BAD（`FAIL ... -> rc=2 (want zero)`
那一行就是它），不是来自本格——三发的 rc 我都逐字贴出来了，不挑对我有利的那一发讲。

### 4.3 顺带把那一格从"抄数字"改成"按腿查数字"

改前那一味是 `sed 's/^want_n 8$/want_n 9/'`。AC#1 恰恰就是会移动 `want_n` 的那次改动：硬写的数字要么下次落不进去
（那一格以"changed not one byte"死去＝一格没人看见的装饰），要么落到**别腿**的 `want_n 8` 上（翻错腿）。
修后按腿号找它的 `want` 行、取下一行的数、写回 +1，读不到数字就 FAIL（不是跳过、不是"当没这格"）。
仍是同一格、同一个语义（抬一腿的基线高于它的实测），既有九格与那一格"表在不在"那一味逐字未动。

### 4.4 AC#9②：腿还在、样本被删空——三形读数与"响不响"的答案

台件＝`fixtures/d_sample_a.go.txt`＋`d_sample_b.go.txt`（两条真开方调用），脚本＝`r3-sample-deleted.sh`，
整份读数＝`logs/harness-sample-deleted-pre.txt`＋`logs/samples-D[012].txt`（rc=0，八条 ok）。
尺的 G6 基线在**副本**里抬到 2 枚（只改一行，`diff` 行数核过；仓内那把尺一字未为此改），于是"删一枚"落在基线以下、
"删光"落到空——腿与声明一枚未动，全程只在那棵丢弃树里换态。

| 态 | 在场样本 | 腿 G6 那一行 | 【基线过期】表 | G6 在不在 BAD 名册 | 聚合退码 |
|---|---|---|---|---|---|
| D2 | 2 枚 | `# ok 腿=G6 声明=ring 基线=2枚 实测=2枚` | 空（表内 0 枚） | 不在 | 5 |
| D1 | 1 枚（删一枚） | `# SHR 腿=G6 声明=ring 基线=2枚 实测=1枚` | `# STALE 腿=G6 声明=ring 基线=2枚 实测=1枚 差=1枚`（表内 2 枚） | **不在** | 4 |
| D0 | 0 枚（删光） | `# BAD 腿=G6 ... 因=空/非空那一味与声明不符（票 161 AC#6①）` | 空 | **在** | 7 |

票面 AC#9② 要的那一句"未修码上响不响"，答案分两形：
**删光＝响**（今天就有牙，牙在 `空/非空` 那一味，与 AC#8① 的腿数断言无关）；
**只删一枚、剩下的低于基线＝不响**（退码里不含 G6 那一味，只在单列那张表里读得到——这正是 AC#8② 那张表的全部价值）。
D2→D1 的 5→4 差额逐枚可归因：共用同一条射程的 `G6neg` 从 BAD 退回 ok（两张 BAD 名册并排打在读数里），
**不是** G6 那一味进了退码。⇒ "腿或它的样本被摘"至此正面台件齐：摘腿（171-v2 已量 rc 0→1）、
摘一枚样本（D1）、摘光样本（D0）三形各有读数。

---

## 5. 门禁与卫生（本程允许的面）

**本程一字节未碰**：`.github/workflows/ci.yml`、`internal/**`、`cmd/**`、`tools/d22scan/**`（含禁令射程与 `allowlist.txt`）、
`docs/PLAN.md`、`docs/specs/**`、`docs/BUILD.md`、`docs/SLO.md`、`internal/observe/thresholds.go`、golden、三枚冻结件。
`frontend/**` 与 `design/**` 未读、未改，本件不引用其内容。AC 框（`- [ ]`）一枚未碰，票面一字未改。
凭据的值不进输出（`logs/doors/flip-*.txt` 里出现的 `secret_*.go` 只是被 `git grep` 打出的文件名与源码行）。

**跑过的（全是不依赖 Go 的尺）**
- `sh .scratch/wisp/probes/154/gate-clauses.sh 4f2a777b...` ⇒ 修后 **rc=2**、`腿数＝14 声明与实测不符＝2`、
  `腿数断言＝相符`、`基线过期枚数＝0`（§3.4）。另在 `f78d3cb` 上跑过同一把改前尺：读数逐字相同（两枚锚同数）。
- `sh .scratch/wisp/probes/171/r3/r3-roster-diff.sh` ⇒ **rc=0**（§3.3）。
- `sh .scratch/wisp/probes/171/r3/r3-trailing-comment.sh pre`／`post` ⇒ **rc=0**（§2.3／§3.2，8＋16 条 ok）。
- `sh .scratch/wisp/probes/171/r3/r3-sample-deleted.sh` ⇒ **rc=0**（§4.4）。
- `sh .scratch/wisp/probes/171/r3/r3-flip-columns.sh pre` ⇒ **rc=0**（§4.1）；`post` ⇒ 见 §4.2。
- 门禁本体在活 HEAD 上的整份跑动＝`logs/door-insitu-live-head.txt`（它自己取 HEAD；本程量时锚已前移到
  `70a935ce...`，比 §0 那枚又晚一枚——同机两程还在落码，所以它的 baseline 那一格与 §2.4 同源，见 §7）。
  现量：**十格里 12 条断言，8 条 ok、4 条 FAIL**，四枚 FAIL 逐枚都是同一下游（`腿数＝14 声明与实测不符＝2`
  那两枚他票 BAD 把"期望 0"与"两枚同翻期望 2"一起顶歪：baseline `rc=2`、两枚同翻 `rc=4`、stale `rc=2`、
  restore `rc=2`）；**本格那两味在新门禁里是 `ok`＋`ok`**（`the stale baseline is its own readable table row` /
  `the row prints its own numbers`，逐字同 §4.2 第 (3) 发那一份）。
  ⚠ 副作用照 `171-r1` 的先例登记：那把门自己往 `.scratch/wisp/probes/161/r6/logs/flip-*.txt`（tracked 件）覆写日志，
  本程**没提交、也没还原**那批件；对照发因此全部跑在 `probes/171/r3/doors/` 的副本上
  （`logdir` 取脚本自身目录＝票面 09-27 19:1x 那条通用规矩，别继承 CWD）。
- 文本级卫生（不依赖 Go）：emoji 频段扫（`U+1F000–1FAFF`／`U+2200–22FF`／`U+2600–27BF`／`U+2B00–2BFF`／`U+FE0F`，
  带正控——现造的 `✓`＋`≤` 样本打得得出 1 行命中，所以"命中少"不是尺哑了）在本程写过的每一件上跑：
  命中的只有 `⚠` 与 `≥`，且逐枚都在**注释行**（尺上四处 `:21`/`:53`/`:62`/`:516` 全是改前既有的注释行，
  与我新加注释同形；仪器"注释豁免、字符串不豁免"那一味＝`Q-46(c)`，出处 `AGENTS.md` §1.2）。
  裸 `go func(`：本程不写 Go，除本件 §7 引用那把门的名字外命中 0。明文密钥：0（`secret` 那几枚命中全是文件名）。
- 提交之后在自己那枚提交上重跑尺本体：`git rev-parse HEAD` ＝ `2a5489d3e49b0b868074c9dcf155db4e883b8a06`，
  `sh .scratch/wisp/probes/154/gate-clauses.sh`（默认锚＝HEAD）⇒ **rc=2**、`腿数＝14 声明与实测不符＝2`、
  `基线过期枚数＝0`，两枚 BAD 仍是 §3.4 那两枚他票的（`G1b`、`G6neg`），逐字读数＝`logs/gate-at-lownhead.txt`。
- **提交后的确认性复跑：我停在了中途，具名在此**。`sh r3-flip-columns.sh post` 在提交后又起了一遍
  （部分日志留在 `logs/harness-flip-columns-post-rerun.txt`，它跑到第 (1) 发的门禁第 3 格里被我叫停）——
  理由不是结果不对，是这台机器上同机两程正在落码，一把门的十一轮 `git grep` 从提交前的约 20 分钟涨到了
  每轮 2 分钟（三发约一小时），超出本程的轮次预算。叫停的方式与时点：`taskkill //F //T //PID 13468`
  （那是本程自己 nohup 起来的 `r3-flip-columns.sh`，不是别的程）；它当时已经开始覆写
  `logs/lie-door-unmodified.txt`（tracked），我按仓里那唯一许可的方式还原了它：
  `git cat-file blob HEAD:.scratch/wisp/probes/171/r3/logs/lie-door-unmodified.txt > 同路径`，
  还原后 `git status --porcelain -- probes/154 probes/161 probes/171` 在我这几个面上**只剩未跟踪的新件**。
- §4.1／§4.2 那三向读数为什么仍然算"提交后的字节"：
  `cmp <(git cat-file blob HEAD:.scratch/wisp/probes/161/r6/flip-declaration.sh) .scratch/wisp/probes/161/r6/flip-declaration.sh`
  与同样一条对 `probes/154/gate-clauses.sh` 的 cmp 都报相同 ⇒ 落进提交的尺与门，和跑出那三向读数的尺与门，
  是同一份字节。**但复跑整轮仍归下一程**（本件不拿"字节相同"冒充"这一小时也跑过一遍"）。
- 删除命令：**仓内 0 次**（丢弃树内换态不算）；worktree／checkout／branch 切换：**0 次**；权限拒绝：**0 次**。
- 盘上留在本程目录里、**故意不入提交**的一件：`probes/171/r3/doors/`（三枚门禁副本
  `door-lie-pre.sh`／`door-lie-post.sh`／`door-honest-post.sh` ＋它们自己写出的 `doors/logs/*.txt`，
  合计约 1.2 MB 的整份门读数）。理由：三枚副本都由 `r3-flip-columns.sh` 从" tracked 门禁的字节 ＋ 两行 patch"
  现生成，其**结论性部分已经逐字抄进 §4.1／§4.2**，整份 dumps 可由
  `sh .scratch/wisp/probes/171/r3/r3-flip-columns.sh pre`／`post` 一键重跑复现；
  留在树里只是不删（临时件只建不删），编排者要收就按这一句收，别按"丢了"处理。

## 6. 我推翻前人（票 171 票面与 171-v2／161-v2 那两张表）哪几句

1. **推翻 `171-v2` 表 §14 给 ③-limit 的定档理由**（原文："不构成退回……该数字与门控 SHR 分支的是同一批变量，
   正常不可矛盾"）。现量：造假尺让同一格行说 `基线=9枚 实测=9枚 差=0枚`，门照样 `ok`（§4.1）。
   "同一批变量"讲的是**分支**，那张表打的是**字符串**——字符串可以独立于分支被写坏，
   所以"不可矛盾"是对意图的判断、不是对读数的判断。⚠ 我没把它的"非阻塞"推翻成"该退回 AC#8"：本格只是做它点名的 (A)。
2. **推翻票面 AC#7 结案里 G2 那一枚的成色**（原文：`run() 腿的分子＝命中行数、含注释行——真树 G2 那两枚命中里
   internal/agent/run.go:546 是一行注释`，以及 `九枚 pair 腿基线一枚未过期`）。
   现量（锚 `4f2a777b`）：G2 名册 3 行，其中 **2 行是整行注释**（`cmd/wisp/run.go:942`、
   `internal/tools/subagent_197.go:332`），只有 `run.go:960` 是代码；行号也从 `:546` 位移到 `:942`
   （台账 A329 说的就是这种位移）。⇒ 后果比"含注释"更硬：改前尺今天把 G2 报成一枚 `BAD 因=新增命中`，
   而那枚"新增"里有一行是散文——**门在指控一枚不存在的东西**。AC#1 修完 G2 的基线是往下（2→1），不是往上抬。
3. **推翻票面 09-27 18:0x 那一格给 flip 新格的写法**（原文：`核它的那一格进了 flip-declaration.sh——只新增那一格、
   既有九格语义一字节未动`，那一味＝`sed 's/^want_n 8$/want_n 9/'`）。现量与理由见 §4.3：硬写的数字在任何一次
   `want_n` 复算之后要么落不进去、要么翻错腿；本程正是"任何一次"（今天 G5neg 恰好还是 8 才没炸）。
   这一条落在票面 AC#8② 那一格的**写法**上，不落在它的判据上。
4. **推翻我自己 §0 的第一版口径**：我原写"开工锚点＝读数锚点"。不成立——两次 `git rev-parse HEAD` 之间 HEAD
   已前移（`f78d3cb`→`4f2a777b`），§5 那把门自己取 HEAD 时已是 `70a935ce`。处置照 `171-r2` 的先例：
   所有"改前 vs 改后"都打同一枚显式锚，绝不拿两次不同 HEAD 比。
5. **另一处自我推翻（仪器坑，不是我写的尺）**：`r3-roster-diff.sh` 第一版 check 2 的 ERE 里写了未转义的 `(`，
   grep 直接报 `Unmatched (`，而那一味照样输出"增删行＝0"＝**一枚假绿**。第二版把词表修好并补了两枚正控（§3.3）。
   这与票面 AC#8 那条"零命中先怀疑尺"是同一枚坑，只是这次是我自己踩的，落盘在读数里而不是抹掉——
   第一版的报错行还在 `logs/harness-roster-diff.txt` 的历史里？不在（同一路径被第二版覆掉），所以此处按
   **正文记账、盘上无原件**的口径写：〔仅自述＋修法可从 diff 复算〕。
6. **不推翻、只对上的**：161-v2 判 G7 反向 3 枚里 2 枚是尾随注释（`161-gates-accept-r2.md:67`/`:68`）；
   本程改完退场的正是那两枚，第三枚 `retention.go` 照旧点名（§3.3）。

## 7. 量不到／判不动的地方

- `scripts/d22scan.sh` 与 `tools/d22scan/runtests.sh`：**量不到**。理由逐字＝派单硬约束"不许跑 `go build`／`go vet`／
  `go test`（任何包）"，而 `scripts/d22scan.sh:1`–`:30` 自述它两步都要 `go`（`runtests.sh -C tools/d22scan ./...`
  与 `go run . -root`）。⇒ 本程用不依赖 Go 的文本尺代跑 emoji／裸 `go func(`／明文密钥（读数在 §5），
  但**这不等于 d22scan 门过了**，那一格留给能跑 Go 的程。
- `gofumpt`：**量不到**（不在 PATH；`r3-trailing-comment.sh` 每次都打一条
  `gofumpt not on PATH; sample formatting left unchecked`，那是具名、不是默认通过）。
  ⇒ 票面"派单前置条件 ②：台件落盘必须 gofumpt 干净"这一味我只能自证形状（tab 缩进、`package` 后空行、
  无未用声明；且台件是 `.go.txt`，本就不在任何包内）。**门我没跑过，别按"过了"记账。**
- 两枚他票 BAD（`G1b 0→1`、`G6neg 1→3`）：**判不动也不动**。归口＝落地那程在自己的切片里核 `want_n`
  （票面 AC#8 那条"谁改口径就在同一格把 want_n 核下来"反过来读也成立：谁改了生产码，谁核自己的基线）。
  ⇒ 连带：`flip-declaration.sh` 的 baseline 那一格今天在任何新 HEAD 上都是 FAIL（它期望聚合退码 0），
  本程**没有**为它放宽那一味；§4.1/§4.2 的判据因此按 label 取行，不按整把门的颜色。
- `.github/workflows/ci.yml`：未碰（AC#3 归编排者另裁）。AC#1／AC#9 两格的修法**不需要**碰 `ci.yml`、
  也**不需要**碰 `tools/d22scan`——两格全是 shell/文本尺，跑它们不需要 Go（本程全部读数都没用过 Go）。
- 仓根那 8 枚未跟踪 `part*.txt`／`state*.txt`（票面 09-27 19:1x 归口给"下一程碰 `probes/171/**`"的那笔账）：
  本程**未收**。它不在我这两格射程里，搬运＝动他程留下的读数件。现量：收尾时仍在原地
  （`ls -1 part*.txt state*.txt | wc -l` ＝ 8，全部未跟踪）。登记在此，等编排者另派。
- 尺 `probes/161/r6/flip-declaration.sh:19`–`:20`（"6 of 14 legs ringing"）与 `:194`（"7 ringing／7 silent"）
  那两句过期文字（171-r2 的 next=④）：一字未动。它不是判据是打印；本程量到的今天的分法是
  响 7 静 7（`基线过期枚数＝0`、BAD 两枚属他票），但那两格句子的归属仍在 AC#8 的 `>` 账上。
- **一格留给编排者的残余（本程没修、也不擅自修）**：AC#9② 跑完剩下的第三形是
  "删掉一枚样本、而剩下的枚数**恰好等于** `want_n`" ⇒ `ok`、表里也没有它——枚数基线这一形在数学上就看不见它
  （要看得见，得让每一腿再登记"样本名册逐枚是谁"，那是新的一味，不在 AC#9② 的判据文字里：票面那一格只要求
  "造一发台件，并答未修码上响不响"）。本程按判据交台件与答案（§4.4），这一形按**发现**登记在此等定档，
  ⛔ 不记成"已收"。
- "下一次在更新的 HEAD 上重跑本尺会得到什么数"：**判不动**。同机两程仍在落码；本件的十四腿基线只在与 §0
  同一枚锚（`4f2a777b`）那棵树上成立，越靠近它们落地的提交，`G1b`／`G6neg` 那一类差额越可能再动。
