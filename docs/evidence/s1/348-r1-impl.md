# 票 348 落地腿 `348-r1` 交件 —— `AC#2` 的成对读数（两平台各一发）

- 身份：写腿（实现者）。⛔ 翻框、⛔ 改票名／`Status:`／`-done`、⛔ 改台账、⛔ push。
- 写面：`internal/winsec/**`（只动测试字节）＋本件。★**产码零改动**：`winsec.go`（`23c2f1e4…`）／`winsec_windows.go`（`5cdef881…`）／`pinnedleaf_342_windows.go`（`66a22239…`）逐枚 sha256 与本腿起手锚一致（尺＝`sha256sum`，SHA-256／64 位）。
- 起手锚 `e94a5b177`；本腿四笔 `195ce6d0f`／`e9caddff`／`5476e2417`／`73aa88049`（§7 逐枚）。
- 全部逐字日志与台架住 `D:\tmp\wisp\20261022\348r1\`（`logs\` / `trees\` / `bench\` / `sh\`），仓外、只建不删。
- 台面：Windows 本机 `go1.27.1 windows/amd64`；POSIX 面＝容器 `golang:1.27`（`3680233e3204`），宿主 module 缓存**只读**挂入、容器 `/tmp` 是唯一临时根、`GOFLAGS=`、`GOTMPDIR/TMPDIR/TEMP/TMP` unset。
- 三棵**仓外**导出树（`D:\tmp\wisp\20261022\348r1\trees\`）：`base2`＝`fde1df02`（修前，`winsec.go=41ba9a9b…`；`pinnedleaf_342_*` 尚不存在）⇄ `land2`＝落地（`winsec.go=23c2f1e4…`、`pinnedleaf_342_windows.go=66a22239…`）⇄ `mut`＝`land2` 只把**本腿那一跳**换回修前那一句（`pinnedleaf_342_windows.go=269a84bc…`）。★四枚仪器件在 `base2`／`land2` **逐枚同 sha256**（如 `ancestorunlink_348_windows_test.go=40aec991…`、`ancestorworld_342_348_windows_test.go=9c195e5a…`）⇒ 成对差异**只有产码**。

## 1. POSIX（容器，一发＝一次 `docker run`）

**①修前见证文件被删 ⇒ 分子 > 0**（`base2`）〔量的〕
- 确定性那一枚（`TestAC2POSIXRefusalPathWithASwappedAncestor…`）**红**，红句逐字：`AC#2 RED: the foreign witness is GONE after the refusal-path unlink (lstat /tmp/…/data/someone-elses-tree/leaf-342r1.bin: no such file or directory)`；引它的那一跳逐字：`… pinned object at …/data/sub/leaf-342r1.bin is not what the name carries: this descriptor is not the object the name carries now …`
- 竞态那一枚，两发：`deleted=76/1500`（`provedLinkCycles=228`／`pinnedWalkRefusals=27`／`pinnedParentBeforeOpen=0`／`truncated=1`）⇄ `deleted=62/1500`（`168`／`27`／`0`／`0`）。

**②修后 0**（`land2`）〔量的〕：确定性那枚 **绿**，survivor 逐字 `size 75 (was 75), mode -rw-rw-rw- (was -rw-rw-rw-), mtime moved=false`；竞态两发 `deleted=0/1500`。

**③同一发证明窗口真到达过链接态**〔量的〕：竞态两发 `provedLinkCycles=267`／`256`、`pinnedWalkRefusals=50`／`54`、★`pinnedParentBeforeOpen=29`／`39`（落地自己那半边在 open **之前**答的）——这些数与 `deleted=0` **在同一行日志里**。

## 2. Windows（本机，一发＝一次 `go test -run …`）

**①修前分子 > 0：git 基线（`base2`）上取不到，读数就是 0，而且这个 0 是读数不是钝尺**〔量的，3 发〕
`deleted=0/400` ×3，但**同一发里窗口到达过链接态**：`provedReparse=66`／`59`／`71`、`truncated=13`／`4`／`7`（＝`privateFile` 的 open 真的跟着翻转的祖先落进了别人的树）。⇒ 缺的不是窗口，是**拒绝路径**：这一平台 `sealHandle → applyDescriptor(f.Name())` 按名跑在这个进程自己拥有的对象上会**成功**，于是 `_ = os.Remove(resolved)` 那一跳根本不执行。两枚仓外台架把这个说法量成事实：
- 祖先 rename：什么都不开＝`err=<nil>`（正控）；按名开着叶子＝`rename …\data\sub …\parked: Access is denied`；经 `os.Root` 开着叶子＝`The process cannot access the file because it is being used by another process`。⇒ 本包**唯一**声明的缝（`var applyDescriptor`，那时叶子句柄正开着）**装不下**"open 之后祖先才翻"的世界。
- 给外来对象装 `DENY(WRITE_DAC|WRITE_OWNER)` 再让 seal 去封它：`applyDescriptorWindows err=<nil>`（属主的隐式 `WRITE_DAC` 压过 deny ACE，并顺手把 deny ACE 清掉）⇒ 造不出"合法且封不动"的外来对象。
⇒ ★**Windows 的 git 基线 ① 今天做不到**；⛔ 拿"0 命中"凑数，⛔ 拿它冒充满足。

**①′ 把本腿那一跳摘掉（`mut`）⇒ 分子 > 0**〔量的，8 发，反向判据〕：`deleted=3,1,1,0,0,0,1,1`（8 发里 **5 发红**），同发 `provedReparse=84/88/87/74/79/83/73/76`、`formBLegRefusals=33/29/27/28/29/35/23/34`。红句逐字：`AC#2 RED: 3 of 400 attempts found the foreign witness GONE, so the refusal-path unlink in privateFile re-resolved the caller's spelling through a swapped ancestor and removed a name in somebody else's tree`。

**②修后 0 ＋ ③liveness**〔量的，4 发〕：`deleted=0/400` ×4；同发 `provedReparse=75/81/81/79`、`formBLegRefusals=32/36/37/36`、`truncated=0,0,0,1`（第 4 发那枚 1＝open 跟着翻转祖先走了，正是"窗口到达"的旁证）。仪器闸逐字：`if st.provedReparse == 0 || st.formRefusals == 0 → AC#2 RED instrument`。

★两平台各一发：POSIX 的一发＝§1 那四行；Windows 的一发＝§2 的 ②③（`land2` 4 发）＋①′（`mut` 5/8 发）。**Windows 的"修前"用摘跳树而不是 git 基线**，理由在 §2 ①，⛔ 混写成一件事。

## 3. 新用例逐名 ＋ 为什么今天响

- `internal/winsec/ancestorunlink_348_windows_test.go:TestAC2WindowsPinRefusalOnAnAbsentAncestorIsNotFiledAsASealRefusal348r1`（确定性，不靠竞态）。
- 为什么今天响：pin 拒与 seal 拒**都**包 `ErrNotSealable`，旧名册只按 `cannot apply a private security descriptor` 归类 ⇒ **缺席的祖先**被记成 "a seal refusal on this route"。本腿现量：一发 400 次尝试里 **395 次**落在"祖先被搬走、还没搬回来"的缺席态，被记成 395 次 seal 拒，只证到 3 次翻转，**却仍过了活性下限**。新用例把祖先直接删掉（全流程无 reparse point），断言落在 pin 行、且**不计**入 form 乙 的 reparse 腿。

## 4. 反向判据三读数 ＋ 还原核

| 读数 | 树 | 逐字 |
|---|---|---|
| **变异红**（删掉新加的那一行名册） | `norow` | `AC#2 RED: an ABSENT ancestor was filed as "a seal refusal on this route"; want the pin-walk row "…"` ＋ `--- FAIL: TestAC2WindowsPinRefusalOnAnAbsentAncestorIsNotFiledAsASealRefusal348r1` |
| **还原绿**（把名册件逐字节放回） | `norow` | `--- PASS` ＋ `AC#2' attribution: the absent ancestor is answered by the pin-walk row, not by the seal row` |
| **原样绿**（未动过的 `land2`） | `land2` | `--- PASS`（§2 ②③ 那 4 发同源） |

★**还原核**：`sha256sum norow/internal/winsec/ancestorworld_342_348_windows_test.go` ＝ `9c195e5ac7fd5a6a91657bac97ea767c1eec433cd81a5a059036517980998d6b` ＝ `git show HEAD:internal/winsec/ancestorworld_342_348_windows_test.go | sha256sum` ⇒ **MATCH**。★**谁算的、什么算法、几位**：两边都是 `sha256sum`＝**SHA-256／64 位十六进制**；`git hash-object` 是 **SHA-1／40 位**，本件**没有**拿它跟上面那个数比（`A1149②`）。
★另两枚同形反向判据（台损修的）：`toReal` 幂等块摘掉 ⇒ `settle the world back to the real state: remove …\data\sub: The directory is not empty.`（10 发里 2 发命中、两批各 1 发；修前仪器在主检出上 3 发 1 发）；摘跳树那 5 发红见 §2 ①′。还原后 `--- PASS`，字节回 pristine（上面同一把 sha256 尺）。

## 5. 门禁 rc（共享工作树、`GOFLAGS=`、四个临时变量 unset；四数分开、⛔ 相减）

- `GOFLAGS= go build ./...` **rc=0**｜`go vet ./internal/winsec/` **rc=0**｜`sh scripts/d22scan.sh` **rc=0**｜`sh scripts/check-path-length-budget.sh` **rc=0**。
- `go test -count=1 ./internal/winsec/`（本机、整包）：`RUN=122`／`PASS=70`／`FAIL=1`／`SKIP=0`（子测试另计：`sub PASS=51`／`sub FAIL=0`）。
- ★格式**甲乙两尺都报**：**甲（判据）**＝`git archive HEAD` 到仓外树跑 `gofumpt -l internal/winsec` ⇒ **0 行**、rc=0（整棵树 20 行全是 `.scratch/` 下存量件，⛔ 本腿写面）；**乙（只作诊断）**＝工作树 `gofumpt -l internal/winsec` ⇒ **0 行**、rc=0；工作树整棵树 rc=2 是 `.scratch/` 里存量的**故意不可解析**探针件（`expected ')'` 一类），⛔ 读成红。归属判定用 `[\\/]`（gofumpt 在 Windows 打反斜杠）。

## 6. 红名册逐名（★取主检出）＋ 新增红 = 0

- 整包那一发（`RUN=122`）**唯一红**：`--- FAIL: TestAC2WindowsLeafIdentityReCheckRefusesAForeignNameAndNeverDeletesIt342r1`。
- ★这枚**不是本腿新造**：本腿改前／改后**同款同率**——改后 5 发 1 红（`AC#2 RED instrument: no leaf leg answered (guard=0, landing's re-check=0) while reparseAt proved 91 flipped leaf cycles across 400`）；改前那棵仪器树 5 发 1 红（逐字同形，`90 flipped leaf cycles`）。⇒ **家族级存量摆动枚**，⛔ 本腿因果。
- 改前整包那发（`e94a5b177`，同一枚仪器字节）有 **2 枚红**：`TestAC1WindowsPrivateFileThroughASwappedAncestorNeverTruncatesAForeignFile342r1`（真残差 `truncated=1..2/400`，本腿改前改后都在）＋`TestAC2WindowsPrivateFileThroughASwappedAncestorNeverDeletesAForeignFile342r1`（台损 `settle the world … The directory is not empty`）。改后前者照旧按竞态出现（`truncated=1,2,1` 三发），后者 3 发 0 红 ⇒ **新增红 0 枚**，且台损那枚被消掉。

## 7. 提交号 ＋ 逐枚 `git show --name-only` ＋ 正控

- `195ce6d0f2737a356114c2acb97550e411e7559a`：2 枚路径，`.go`=2，写面外=0。
  `internal/winsec/ancestorunlink_348_windows_test.go`／`internal/winsec/ancestorworld_342_348_windows_test.go`。
- `e9caddff4dccb3e5a3fb7b4d2f622890de5968b6`：2 枚路径，`.go`=2，写面外=0（同上两枚）。
- `5476e2417fe8847f91ce864af00f7001f1362e97`：1 枚路径，`.go`=0，写面外=0（`docs/evidence/s1/348-r1-impl.md`，本件）。
- `73aa88049`：1 枚路径，`.go`=0，写面外=0（本件两处读数更正）。
- ★从 `5476e2417` 起，本腿其余各笔**只碰本件**：判据＝`git log --format='%H %s' --grep='^票 348 落地腿 348-r1'` 现跑得名册、再逐枚 `git show --name-only` 数 `.go` 与写面外（都是 0）——这条**不依赖枚数**，本件自身再改几笔也不必回来改它。
- 前两笔只 **＋行**（`195ce6d0f` `2 files changed, 104 insertions(+)`；`e9caddff` `+7 −2`，那 −2 是本腿**自己**前一笔那 7 行里挪常量，⛔ 别人字节）。
- ★**正控**（同尺打一笔**含 `.go`** 的真提交）：同一把 `git show --name-only`＋`grep -c '\.go$'` 打 `195ce6d0f` ⇒ `go=2`（尺有牙）；再打一笔**不含 `.go`** 的真提交 `efa950a3`（票 164 续跑腿的判语件）⇒ `go=0`（尺不虚报）。
- 本腿写面现状：`git status --porcelain -- internal/winsec/ docs/evidence/s1/` 为空；本腿新增的唯一路径就是本件（`docs/evidence/s1/348-r1-impl.md`）。

## 8. ★停手项（未定义即停，照实报）

1. **Windows 的 git 基线 ① 今天做不到**＝§2 ① 那两枚台架读数（祖先 rename 三态、属主隐式 `WRITE_DAC`）。⇒ 本腿**不声称** `AC#2` 满足：POSIX 半边是"基线 >0 ⇄ 落地 0 ＋ 同发 liveness"；Windows 半边是"落地 0 ＋ 同发 liveness"，其"修前 >0"用**摘掉本腿那一跳的树**给出（5/8 发），⛔ 用 git 基线的 0 冒充。
2. 要补 Windows 那枚"git 基线 ①"，只有一条路：在**产码里**加一枚**删除那一跳的接缝**（`f.Close()` 之后、`os.Remove` 之前）——那是新增产码名字／改函数形状＝**第二道人工批准**，⛔ 本腿自扩。
3. `TestAC1WindowsPrivateFile…342r1`（票 342 的截断残差）今天仍按竞态红（`truncated=1..2/400`）＝**票 342 的判据族**，⛔ 本腿修、⛔ 记成本腿新增红。
