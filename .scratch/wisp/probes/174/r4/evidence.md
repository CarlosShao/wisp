# 票 174 · 落地腿 `174-r4` 证据件

> 写面＝`internal/tools/**` ＋ `internal/agent/spill.go` 的回执文案／断言侧。
> 本件是**第四任落地程**（前人已跑 `c1 c2 r1 r2 r3 v1 a3 a4`）：⛔ 不复用前人任何读数，
> 本件里每一个数都是本程现跑，带时刻；凡引前人结论处一律标〔来路〕并明写"本程未复算"或"本程复算"。
> 规矩：AC 框一枚不碰（勾它＝另一枚非实现者程）。

## 0. 起手锚与写面脏度 —— 本节打算答：我从哪一刻的哪一枚树起手的、我落笔前写面是不是干净的

- 现跑钟点 `date` ＝ **15:54:25+0800**（起手）／commit 时刻见 §4 每枚名。
- 起手 HEAD ＝ **`40aae961`**（265-r1 的骨架枚；台账 `A591` 那批之后的在飞腿）。落盘骨架时 HEAD 已前进到 `9f1e89bf`（`257-r1` 骨架枚）⇒ 按第 82 条只记"此刻取到的是别的程的枚"。
- 起手写面脏度：`git status --porcelain -- internal/ cmd/` ＝ **0 行**（台件 `.scratch/wisp/probes/174/r4/start-dirty-internal-cmd.txt`）。全树脏 **628 行**（共享工作树，别人在飞），⚠ 与我无关，我不清不碰。
- 票面未勾框尺（本程现跑，非抄任何人）：`grep -c "^- \[ \]" <票面>` ＝ **5**、`grep -c "^- \[x\]"` ＝ **3**。五枚 ＝ **AC#3 / AC#2b / AC#2c / AC#4 / AC#5**。

## 1. 我选了哪一格，为什么不是别的格 —— 本节打算答：五枚未勾框里我落哪一枚、其余四枚各自卡在哪（具名到包与在飞腿）

**我落的格＝AC#3 的 (ii) 半枚为主、(i) 半枚为次**：把"判这条路径在不在授权根里"**只可能走 C26**、以及**判定腿永不写盘**这两件事，从散文注释变成**会响的判据**。写面＝`internal/tools/**` 新增一枚 `_test.go`，⛔ **零产码改动**（理由见下 AC#2b 那一条：`task.go` 那两句文案今天受模板冻结钉系着，我一个字不动它）。

现量"今天没有判据钉着这件事"的尺（本程跑，逐字）：

| 尺 | 值 |
|---|---|
| `grep -c "GetShortPathName\|shortName" internal/tools/task_output_pointer_notice_test.go` | **0** |
| 同上 `task_output_canonicalize_fail_174_test.go` | **0** |
| 同上 `internal/agent/spill_pointer_honesty_174_test.go` | **0** |
| `grep -c "Workspace" task_output_pointer_notice_test.go` / `…canonicalize_fail_174_test.go` | **0 / 0** |
| `grep -rln "pointerNotice" internal/` | `internal/agent/spill.go`、`internal/tools/task.go`、`internal/tools/task_output_canonicalize_fail_174_test.go`、`internal/tools/task_output_pointer_notice_test.go`（＝只有"说了什么话"那一族，没有"谁说了算"那一族） |

**为什么不是别的格（逐枚具名卡点）**：

- **AC#2b（接线＋精确文案）**：双堵。①接线那一行在 `cmd/wisp/run.go:365`——**不在我的写面**，且 `265-r1` 正在写 `cmd/wisp/**`（派单写死我不碰不跑不归因）；②剩下的文案半格在 `internal/tools/task.go:328-329`（现锚），那两句被同文件 `:299`／`:349` 两枚**模板冻结钉**系着，"要不要加'或批一张 L2 卡'那半句"＝契约轴（`Q-63` 甲同一条边界，须 owner 裁）〔来路＝`174-c2` §④ 与编排者 22:3x 第 3.3 条；**本程未复算该裁量权，只按"未裁即不动"执行〕。⇒ 我一动那两句就是替 owner 拍契约面。
- **AC#2c（`spill.go` 的桩说实话）**：已由 `174-r2` 落地（三枚 commit `fad39b2c`→`5b14afe5`→`081019c0`，341 行判据；编排者在 `A571` §1 代跑三发证明牙是活的），**待非实现者翻勾**。我再做＝重条目＝按派单算缺陷。⇒ 不选，且我这一格刻意不去碰 `internal/agent/**`（见 §7 一处我量不到的边界）。
- **AC#4（契约轴）**：要求是"未拿到 owner 批准前对 `PLAN.md`／`specs` **零字节**"＋七面零字节自证。它不是产码格，我做完只有"没动"可报 ⇒ 由 §4 的写面自证覆盖，不占我的选格。
- **AC#5（与 `Q-59` 分开结线）**：票面自己写的是"**本票结题时**必须逐条说明"＝翻勾那一枚程的陈述义务。⇒ 我在 §8 供读数，不顶它的格。
- **AC#3 里我刻意不做的两半**：**(iii)**（`[fs] allowed_dirs` 默认值与 `SPEC-03:35` 一字节不改）——默认值住在 `internal/config/**`，**`257-r1` 此刻正在写那包**⇒ 我不碰不跑不归因，且我若去钉它就是把别人的在飞面包进我的判据；**(i) 的名册那一面**（不得新增一枚受门控的 artifacts 工具）——`D34` 权威表名册钉已在 `internal/tools/fs_test.go`／`bridge_junction_windows_test.go`／`task_cancel_221_legs_test.go` 等六枚文件里（本程 `grep -rln "D34"` 现量），重钉＝重条目 ⇒ 我把 (i) 只做在**判定腿只读性**这一面上（`os.Stat` 那一步绝不变成写入/创建），那一面今天**零判据**。

## 2. 判据形状：这一发在未修码上响不响 —— 本节打算答：我新立的判据问的是什么、"未修码不响／修后响"这两侧各由哪一枚读数支撑

**口径先说清（这是一枚"自证"格，不是"修法"格）**：AC#3 的三枚禁区管的是**任何修法都不许越线**，它不新增可修的行为。所以对这一格，"这一发在未修码上响不响"的正确答法是：
**未修码上，现场行为已经是对的（今天就说实话、今天就不写盘），但"对"这件事零枚判据钉着 ⇒ 下一程把它改错时不会有任何一枚红。** 本程交的就是那枚"改错时会响"的红。凭据＝§5 的三发突变（每一发都是**当前产码**上种下的违例，逐名红），不是任何人的自述。

落点＝`internal/tools/task_pointer_authority_ac3_174r4_windows_test.go`（`//go:build windows`，305 行）＋`internal/agent/spill_pointer_authority_ac3_174r4_test.go`（§9 那一枚，见本节末"两腿"）。三＋一枚判据，各问一件事：

| 判据 | 问什么 | 为什么现有十二枚（tools 侧）／八枚（spill 侧）钉不住 |
|---|---|---|
| `TestPointerAuthorityFollowsC26NotTheSpelling174r4` | 宿主登记的拼写与 C26 的规范形**不一致**时（8.3 短名，向 OS 取的，不手打），回执必须跟 **C26 的答复**（说内＝安静），不跟字面前缀（说外＝误报） | 既有判据登记的**全是 `mustCanonical` 之后的规范拼写**，`InAllowlist(canon)` 与 `InAllowlist(raw)` 在那些输入下**逐字同值**⇒ 摘成 raw 一枚都不红（§5 MUT-A1 现量：恰 1 枚红＝本枚） |
| `TestPointerAuthorityNarrowsWithWorkspaceNotLexicalRoot174r4` | 同一枚真实文件，C26 因**工作区收窄**（票 92 AC#3）说"外"，而字面前缀说"内"⇒ 回执必须说"外"，并与同桥 `fs.read` 的 `L2` 对齐；反形＝`ClearWorkspace` 后必须回到安静＋`L0` 全文 | 这一形是**欠报方向**（把走不通的指针说成走得通），正是本票要防的那句谎话的形状；既有判据里没有一枚动过收窄面（`grep -c "Workspace"`＝0，§1 表） |
| `TestPointerCheckWritesNothing174r4` | 存在性那一步**只读**：指向不存在的副本 ⇒ 该路径回执后仍不存在、目录条目集不变；指向健康副本 ⇒ 字节／size／mtime 三项不变 | `task.go:821-824` 那句 "never creates, moves, truncates or removes" 今天是**散文注释**；`os.Stat`→`os.OpenFile(O_CREATE)` 这一步既过 d22scan（无 `filepath.Clean|Abs` token）也不改任何既有文案断言 |
| `TestSpillPointerAuthorityAsksWithTheCanonicalAnswer174r4`（§9） | spill 那半条腿：问 `InAllowlist` 的**必须是 `Canonicalize` 返回的那一枚串** | 174-r2 的 `fake174Judge.InAllowlist` 是 `return f.inRoot`——**完全不看入参**（现读 `spill_pointer_honesty_174_test.go:79-82`），所以那八枚判据对"问错了对象"全盲 |

两腿都要钉的理由：AC#2 的原文并列点名 `task.output` **与** `spill.go` 的桩（票面 §AC#2 与 `174-r1` 的边界③），同一句"全文见"今天在两条腿上；AC#3 写的是"**任何修法**都要自证"，只钉一条腿＝另一条腿的禁区仍零凭据。⚠ 我在 spill 腿**只新增判据、不动 174-r2 的任何一枚**（那八枚待非实现者翻勾，我改它会让验收腿的对数失配）。

## 3. BEFORE 读数（未修码逐字） —— 本节打算答：动判据之前现场说什么话，逐字贴，含 `grep -n` 落地证明

判据文件落下之后、**产码一字未改**的状态下现跑（`date`＝**17:23:35+0800**，`go test -count=1 -v ./internal/tools/ ./internal/agent/`＝`ok 13.748s`／`ok 1.929s`，顶层 **PASS=275 FAIL=0 SKIP=0**）。三枚新判据在未改产码上**全绿**，逐字读数（`-v` 里的 `t.Logf`，台件 `logs/gate-targeted.txt`，同 `logs/mutA1.txt`/`mutA2.txt` 里的 PASS 段）：

- 短名臂（安静侧）：`two spellings of one file: long="C:\\Users\\swq\\AppData\\Local\\Temp\\TestPointerAuthorityFollowsC26NotTheSpelling174r4…\\001\\tool-output-174r4-alias.txt" short="C:\\Users\\swq\\AppData\\Local\\Temp\\TESTPO~1\\001\\TOOL-O~1.TXT"`
  回执逐字：`[…输出已落文件：省略 17200 字符，总长 20000 字节 / 约 5000 token，全文见 C:\\Users\\swq\\AppData\\Local\\Temp\\TESTPO~1\\001\\TOOL-O~1.TXT…]`
  ⇒ **零 `注意：`**，且同一枚桥对该拼写 `fs.read`＝`isError=false level=L0`、全文 20000 字节（该读数由判据内 `readThroughBridge` 断言钉着，非本程另跑一发）。
- 短名臂·反形（撤根）：`[…输出已落文件：省略 17200 字符，总长 20000 字节 / 约 5000 token；注意：这条路径现在读不到，它不在你被授权的目录范围内，fs.read 会被拒；要用户先把所属目录加进 [fs] allowed_dirs 才读得回来，全文见 …TESTPO~1\001\TOOL-O~1.TXT…]`
  ⇒ 同一枚串，只改授权根，答复翻转＝安静不是恒真。
- 工作区收窄臂：`[…省略 17200 字符…；注意：这条路径现在读不到，它不在你被授权的目录范围内，fs.read 会被拒；要用户先把所属目录加进 [fs] allowed_dirs 才读得回来，全文见 …\artifacts174r4\tool-output-174r4-narrowed.txt…]`，且同桥 `fs.read`＝`isError=true level=L2`（判据内断言）。
- 收窄臂·反形（`ClearWorkspace`）：`[…省略 17200 字符…，全文见 …\artifacts174r4\tool-output-174r4-narrowed.txt…]`＝零 `注意：`＋`fs.read` 回 `L0` 全文。
- 只读臂：ghost 支回执逐字含 `宿主登记的那份副本文件并不存在`；健康支逐字与短名臂的安静形同构（`注意：` 零次）。

**"今天的行为对、但没有尺"的落地证明**：`grep -n` 现场锚（本程跑）——判定腿只有这四支会说话，且全部只断"说了什么"：

```
internal/tools/task.go:834:  canon, err := d.Paths.Canonicalize(raw)
internal/tools/task.go:838:  case !d.Paths.InAllowlist(canon):
internal/tools/task.go:842:  if st, statErr := os.Stat(raw); statErr != nil {
```

（行号＝本程现锚；票面 §AC#2d 写 `:313-314`、编排者翻勾节写 `:836-837`，**两枚都不是现在这枚**，以本节为准，归口见 §7。）

## 4. 改动清单与写面自证 —— 本节打算答：动了哪几枚文件（`git diff --numstat` 三列）、禁改面是不是真的零字节

本程名下提交（逐枚 `git show --numstat`，三列＝增/删/路径）：

| commit | 时刻 | 文件 | +/− |
|---|---|---|---|
| `8e2e1962` 骨架 | 15:5x | `.scratch/wisp/probes/174/r4/evidence.md` | 42 / 0 |
| | | `.scratch/wisp/issues/174-…-fs-allowed-dirs.md` | 1 / 0 |
| `2d6aed92` §0/§1 | 16:0x | 同上两枚 | 23 / 2、1 / 0 |
| `7c5f2316` 判据 | 17:01:55 | `internal/tools/task_pointer_authority_ac3_174r4_windows_test.go` | **305 / 0** |
| （下一枚）判据修订 | 17:2x | 同上文件（`t.Fatalf`→`t.Errorf`＋两句注释，把"说了没有"与"写没写"两问分离） | **4 / 1** |
| §9 那枚（spill 腿） | — | `internal/agent/spill_pointer_authority_ac3_174r4_test.go` | 见 §9 |

写面自证（**按枚算我自己带上了什么**，不用 `git diff <anchor> HEAD` 那种会把别人提交卷进来的尺）：

- 我名下每一枚的 `git show --numstat` 如上表 ⇒ 除 `internal/tools/**` 一枚**新增**测试文件（删除列 **0**）与 `.scratch/wisp/probes/174/r4/**`、票面 Progress log 追加行之外，**零文件**。
- **产码（非 `_test.go`）＝零字节**：三枚突变全部用 `git cat-file blob HEAD:internal/tools/task.go > internal/tools/task.go` 还原，md5 链见 §5，末态 `git status --porcelain -- internal/` ＝ 只剩我自己的判据文件（提交后归零）。
- **票面 AC 框＝零改动**：`grep -c "^- \[ \]"` 起手 **5**／交件复量仍 **5**，`grep -c "^- \[x\]"` 起手 **3**／复量 **3**（尺现跑于 §8）。
- 一字节未动的禁改面（本程未读未写未跑）：`docs/PLAN.md`、`docs/specs/**`（`SPEC-03:35` 只**读**了那一句用作 §7 的具名事实：逐字 `| 🔒 \`[fs]\` | \`allowed_dirs[]=[]\` \`reparse_point_exceptions[]=[]\`（按具体路径）\`delete_enabled(bool)=false\` | 🔒 同上 |`）、`internal/observe/thresholds.go`、任何 golden、`tools/d22scan/allowlist.txt`、`internal/risk/**`、`internal/config/**`、`cmd/**`。

## 5. 突变自证（每发红句逐字＋还原 md5 三点） —— 本节打算答：每一发突变种下必红、还原必须等于 `git cat-file blob HEAD:<path>`，红句逐字贴

起手基线（`git cat-file blob HEAD:internal/tools/task.go`，现算）＝ **`4138177e29ffff427776a73eb0d61a9b`**，与工作树逐字节同值。还原一律用 `git cat-file blob HEAD:internal/tools/task.go > internal/tools/task.go`（⛔ 不是"重读当前文件再写回"），链在 `logs/md5-chain.txt`：**起手＝还原＝HEAD blob，三次全部 `4138177e…`**，每发还原后 `git status --porcelain -- internal/tools/task.go` 空。

### MUT-174R4-A1 —— 判定改用宿主登记的串（`canon`→`raw`）
落地证明：`internal/tools/task.go:835: _ = canon // MUT-174R4-A1: decision moved off C26's canonical answer`、`:839: case !d.Paths.InAllowlist(raw): // MUT-174R4-A1`
现跑（17:04:19+0800，整包 `-v`）：`PASS=193 FAIL=1 SKIP=0`，红名**恰一枚**＝`--- FAIL: TestPointerAuthorityFollowsC26NotTheSpelling174r4`。
红句逐字（`logs/mutA1.txt:985-987`）：

```
task_pointer_authority_ac3_174r4_windows_test.go:135: a pointer C26 authorizes must stay silent whatever spelling the host filed; found "注意：" in "[…输出已落文件：省略 17200 字符，总长 20000 字节 / 约 5000 token；注意：这条路径现在读不到，它不在你被授权的目录范围内，fs.read 会被拒；要用户先把所属目录加进 [fs] allowed_dirs 才读得回来，全文见 C:\\Users\\swq\\AppData\\Local\\Temp\\TESTPO~1\\001\\TOOL-O~1.TXT…]" - that is a decision made on the string, not on C26's answer
…:135: … found "读不到" in "[…同上…]" - that is a decision made on the string, not on C26's answer
…:135: … found "不在你被授权的目录范围内" in "[…同上…]" - that is a decision made on the string, not on C26's answer
```

⇒ **既有十二枚判据一枚不红**（`PASS=193` 里含它们）＝"把判定挪回字面"这一步今天没有任何尺量得到，这一发就是那把尺的牙；`--- SKIP`＝0（没有任何一枚靠跳过换绿）。

### MUT-174R4-A2 —— 授权判定整个搬离 C26（手搓 `foldPath`+`HasPrefix` 对 `Roots()` 做包含）
落地证明：`:835: _ = canon // MUT-174R4-A2: containment recomputed by hand, off C26's answer`、`:847: }(): // MUT-174R4-A2`（中间是 8 行 `for _, r := range d.Paths.Roots() { if f == r || strings.HasPrefix(f, r+pathSep) … }`）。
现跑（17:11:54+0800）：`PASS=192 FAIL=2 SKIP=0`，红名＝`TestPointerAuthorityFollowsC26NotTheSpelling174r4`、`TestPointerAuthorityNarrowsWithWorkspaceNotLexicalRoot174r4`。
收窄臂红句逐字（`logs/mutA2.txt:991-992`）＝**欠报方向**，那枚指针明明被 `fs.read` 拒了：

```
…:205: workspace-narrowed arm verbatim: IsError=false Truncated=true Text="[…输出已落文件：省略 17200 字符，总长 20000 字节 / 约 5000 token，全文见 C:\\Users\\swq\\AppData\\Local\\Temp\\TestPointerAuthorityNarrowsWithWorkspaceNotLexicalRoot174r41504622155\\001\\artifacts174r4\\tool-output-174r4-narrowed.txt…]"
…:211: C26 says this file is outside the narrowed scope, so the pointer must say so; got "[…同上，零句说明…]"
```

**d22scan 盲区复量（这一发的凭据，不是我的一句推断）**：突变在盘上时跑 `sh scripts/d22scan.sh`＝`rc=1`，但**全部 1 条 finding 是 `cmd/wisp/resident_approval_windows.go:38 [phantom-citation]`**（别家腿在飞，归因见 §6），`internal/tools/` 那一段＝`scope ban #7 internal/tools/ examined 23 production Go files`、**零 finding**；`grep -n "task.go"` 在该日志里命中 **0**。⇒ 逐字坐实：`filepath.Clean|Abs` 两枚 token 的静态扫，扫不到"用 `strings.HasPrefix` 手搓包含判定"这一形 ⇒ **AC#3(ii) 不能拿 d22scan 当凭据**，必须有行为级判据。

### MUT-174R4-B —— 存在性探针变成写（`os.OpenFile(raw, O_RDWR|O_CREATE)` 后 `os.Stat`）
落地证明：`:843: _ = f.Close() // MUT-174R4-B: the probe materializes what it only meant to look at`。
第一发（17:14:40，判据仍是 `t.Fatalf` 版）：`PASS=192 FAIL=2 SKIP=0`＝`TestPointerCheckWritesNothing174r4`、`TestPointerToMissingCopyFileSpeaks`——但我的那一枚红在**文案层**（`:268 a pointer at a missing copy file must say so`）就停了，**文件系统那一层没跑到**。⇒ 我把该处 `Fatalf` 改成 `Errorf`（＝§4 那枚 `4/1` 修订，两问分离），重跑第二发（17:21:17，`logs/mutB2.txt`）：`rc=1`，红名同上两枚，我的那一枚**逐字打出三条**，其中两条是文件系统事实：

```
…:271: a pointer at a missing copy file must say so, got "[…省略 17200 字符，总长 20000 字节 / 约 5000 token，全文见 C:\\…\\tool-output-174r4-ghost.txt…]"
…:274: answering a pointer at an absent copy file must NOT create it: Lstat=<nil>
…:288: the pointer check wrote into the artifacts tree: before=[tool-output-174r4-intact.txt] after=[tool-output-174r4-ghost.txt tool-output-174r4-intact.txt]
```

如实报重叠：`TestPointerToMissingCopyFileSpeaks`（174-r1 名下）在这一发也红——它红在文案，我红在**写盘事实**（目录条目集＋`Lstat` 存在性），两枚管的不是一层；我没有为此改动它一字。

### §9 那枚（spill 腿）的突变

未判——见 §9（本节末追加，读数先落 §9 再引）。

## 6. 门禁读数（带时刻） —— 本节打算答：`d22scan` / `check-path-length-budget` / `gofumpt -l` / `go vet` / 定向 `go test -count=1 -v` 的现跑值与红名集合逐名比对（哪几枚是既有红、不归我）

| 门 | 时刻（现跑 `date`） | 读数 |
|---|---|---|
| `gofumpt -l`（`$(go env GOPATH)/bin/gofumpt.exe`） | 17:22:18+0800 | 对我的判据文件＝**空**（rc=0，归零） |
| `go vet ./internal/tools/` | 17:22:18+0800 | **rc=0**，零行输出 |
| `sh scripts/check-path-length-budget.sh --with-self-test` | 17:22:18+0800 | **rc=0**，逐字 `denominator: tracked paths=5634  over-budget=57  covered by roster=57  not in roster=0`、`VERDICT GREEN`、`unit cross-check: … SAME 57 paths`（⚠ 最长那枚＝180 字符，是票 252 的工单名，不是我新增的文件；我的新文件在预算内，否则 `not in roster` 会跳） |
| `sh scripts/d22scan.sh` | 17:22:18+0800（干净树复量） | **rc=1**，两枚红＝`TestScannerSelfScanOfRealRepoIsGreen`、`TestRealRepoLedgerIsHonest`，finding 逐字一条：`cmd/wisp/resident_approval_windows.go:38: [phantom-citation] comment cites repo path "internal/panel/panel_pump.go" which does not exist on disk (ticket 212 ban #9)` |
| `go test -count=1 -v ./internal/tools/ ./internal/agent/`（定向，⛔ 未跑 `./cmd/wisp/`） | 17:23:35+0800 | **rc=0**，`ok internal/tools 13.748s`／`ok internal/agent 1.929s`，顶层 **PASS=275 FAIL=0 SKIP=0**；三枚新尺逐名 PASS |
| `go test -count=3 -run 'Pointer|Truncat|TaskOutput|Spill' ./internal/tools/ ./internal/agent/`（安静复量） | 17:24:12+0800 | **rc=0**，两包 `ok 1.593s`／`ok 1.488s`＝三连不抖 |

**d22scan 那一枚红的归因（不归我，我也没去修它）**：尺的措辞写着 "repo HEAD violates"，但它实读**工作树**——现量两侧：工作树里 `grep -c "panel_pump" cmd/wisp/resident_approval_windows.go`＝**1**，`git cat-file blob HEAD:cmd/wisp/resident_approval_windows.go | grep -c panel_pump`＝**0**，且 `git status --porcelain -- cmd/wisp/resident_approval_windows.go`＝` M`（未提交）。⇒ 那一行只在 `265-r1` 此刻的在飞工作树里，**HEAD 上没有**，我这发是并发期跑的，这一格按派单"不碰不跑不归因 `cmd/wisp`"处理，⛔ 未做任何"顺手修红"。同发里我自己的面＝`scope ban #7 internal/tools/ examined 23 production Go files`、零 finding。
⚠ 既有常红名单（`internal/panel` 的 `TestC21DesignTokensFourWayAgree`／`internal/ball` 同族）＝**本程未跑这两包**（定向口径），所以不在我的红名集合里，也不由我判。

## 7. 我攻不动／判不动的地方（具名＋归口） —— 本节打算答：本程量不到、写不得、须他人拍板的每一格，各自归谁

未判。

## 8. 结题账：本格对票面各 AC 的推进位（含与 `Q-59` 的分开结线） —— 本节打算答：我这枚做完之后 AC#3 三半各是什么状态、`Q-59` 与本票为何不是一处断口

未判。

## 9. spill 腿（`internal/agent`）—— 本节打算答：同一条禁区在另一条腿上有没有牙

未判。
