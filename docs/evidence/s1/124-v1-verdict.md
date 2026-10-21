# 票 124 终裁腿 124-v1 — 批 1 出口裁决（非实现者）

裁决者 NOT 实现者。被裁对象＝实现腿 `124-r1` 本轮：码 `d32b459c`、件 `abd09226`。
派单 `8b4989c0`（`.scratch/wisp/probes/orch/20261021/txt/124v1-dispatch.txt`）。
开工首读 HEAD `8b4989c0`，收尾复读 HEAD `43cca3a2`（其间漂过 `ac2847d5`/`6633c45d`，全是编排者存档与台账）；
`git diff --name-only 8b4989c0 43cca3a2 -- internal cmd go.mod go.sum` **零 hunk**、
`git diff --name-only d32b459c HEAD -- internal cmd go.mod go.sum` **零 hunk**、
`git status --porcelain -- internal cmd` **空** => 被裁码与读数同版。
台面全部我自己现跑：纯净件 `git archive <sha> | tar -x`（PRE 锚 `4ac4549c`、POST 锚 `d32b459c`，各 17488 文件），
容器 `golang:1.27` / `go1.27.1 linux/amd64` / `GOPROXY=off`，挂载自证 `go.mod` 1026 字节（md5 `0c266f4d...` 两锚同字）。
软链形 `ln -s /realpriv /varlink` + `TMPDIR=/varlink/w124tmp` + `exit 97/98`；普通形 `[ ! -e /varlink ] || exit 99` + `/plainroot`。
落点 `D:\tmp\wisp\20261021\124v1\`（`logs/` 日志、`rosters/` 名册、`harness/` 脚本）。

名册尺前缀**从盘上现量**：`git log --all --grep='^票 124 落地腿 124-r1'` => **2 笔**（`d32b459c` / `abd09226`）。
松尺 `--grep='124-r1'` => 9 笔（含编排者派单与裁决存档，非本腿）；`--grep='落地腿'` => 357 笔（跨票）。
=> 本件一切枚数走严格前缀那 **2 笔**。

> 本件因 160 行帽拆三份：本件 = sec 0 与 sec 一（七项判语）；`124-v1-verdict-part2.md` = sec 二、三、四；
> `124-v1-verdict-part3.md` = sec 五..八（红名册逐名、批 2/3 读数、建议、提交与正控）。

## 一、七项逐项判语与读数

### 1) 递根点是否接上票 119 纪律 + 零新增走法 — **成立，但"6 处"这个数是错的**

- 尺 A（helper 来路）：`tempdir_resolved_124_test.go:31-34` 逐字 `return proc.SealableRoot(t.TempDir())`，
  `git log --oneline -- internal/config/tempdir_resolved_124_test.go` => 唯一一笔 `fafe2b4a`（2b-1 批已入树）
  => **本腿零新增 helper、零新增走法**，`R-125-2` 副本账不增行。**这是复用，不是新写。**
- 尺 B（本轮实际改了几处）：`git show --numstat d32b459c` => `added=13 deleted=9`；`--unified=0` 的 `@@` 头 **9 枚**。
  删除侧 9 行**逐行**都是 `t.TempDir()` 那一个表达式（`cat -A` 核过，无一行是注释）；
  新增侧 13 行＝9 枚同名替换 ＋ `manager_test.go` 的 4 行注释。
  => **实数是 9 处 / 9 hunk / 4 文件**，NOT 6。
- 尺 C（自报内部是否自洽）：件 sec 3 那张逐文件表自己写的是 manager 2 ＋ settings 5 ＋ enabled 1 ＋ writeguard 1 = **9**，
  而 sec 0、sec 4、commit 主题写 **6** => **同一份交件里 6 与 9 并存，不成立的是那个"6"**。派单里的"6 处"同源，一并作废。
- 尺 D（零新增走法正控）：HEAD blob 上 `sealableTempDir124(t)` 调用点 **14 处**（8 文件），本腿加 9 处；
  `proc.SealableRoot` 全树 9 个包各一枚一行委托、无第二形。8 枚红文件**逐枚**经由那两枚改过的共享 helper
  （`newTestManager` 覆盖 7 文件、`cleanMachineManager` 覆盖 settings_257）=> 覆盖率解释得住。
- 判语：**纪律接上了、零新增走法成立；"6 处"应更正为 9 处。** 这是读数错，不是行为的错。

### 2) 批 1 出口 — 软链形 only-in-link 归零**成立**；"普通形一枚没多红"**单发不成立**（见 sec 二、三）

四发四数我自己现量（每数分列，未相减；PASS＝顶层与子测试两个独立计数）：

| 发 | 顶层 PASS | 顶层 FAIL | 顶层 SKIP | 子 PASS | 子 FAIL | RUN 行 | mech 串 |
|---|---|---|---|---|---|---|---|
| PRE-LINK  | 77  | 35 | 1 | 73 | 0 | 186 | 35 |
| MUT-LINK  | 71  | 41 | 1 | 72 | 1 | 186 | 41 |
| POST-LINK | 110 | 2  | 1 | 81 | 0 | 194 | 0  |
| PRE-PLAIN | 113 | 0  | 0 | 83 | 0 | 196 | 0  |
| POST-PLAIN| 111 | 2  | 0 | 83 | 0 | 196 | 0  |

名册作差（`comm`，NOT 相减）：
- `only-in-link PRE` = PRE-LINK FAIL 减 PRE-PLAIN FAIL => **35**；扣掉那 2 枚摇摆名 => **33**。
- `only-in-link POST` = POST-LINK FAIL 减 POST-PLAIN FAIL => **0**（两份名册均只剩那 2 枚摇摆名）。
- `新增红` = (PRE 两形并集) 减 (POST 两形并集) => **0 名**；反向 `went-green` => **33 名**。
- `mech` 串（`not provably resolved` / `refusing to seal`）**35 => 0**，`/varlink` 39 => 4，`test timed out` 0，`panic` 0。
- 33 枚 `went-green` **逐名**在 PRE-PLAIN 与 POST-PLAIN **都是 PASS** => 它们从来只在软链形红，不是"变绿"换"跳过"。
- SKIP 名册两形恒同一枚（`TestAC2POSIXSeamInstallSurvivesASymlinkSpelledTemp125`，票 125 seam 探针）=> 不冒充绿。

=> **出口前半句成立：软链形 only-in-link 33 => 0，且 0 枚新增红名。**
=> **出口后半句今天测不出来**：普通形红数本身在 0..2 之间摇（sec 三的 8 发重复量），我这一发 PRE-PLAIN 得 **0**、
POST-PLAIN 得 **2**，与实现腿的自报（PRE 2 / POST 0）**正好相反** => "改前普通形就红 2、改后仍 2" 这个叙事是单发巧合。

### 3) AC#3（不许放行侧放宽换绿）— **成立，四条独立尺都过**

- 尺 1（删除侧逐行）：9 行删除**全部**是 `t.TempDir()` 表达式，0 行触及断言/阈值/golden/分支；
  用 `Fatal|Errorf|want|expect|assert|Skip|threshold|golden|if |return` 扫全部 `+/-` 行，唯一命中是注释里的 "assertion" 一词。
- 尺 2（地界）：`--name-only` 4 枚全在 `internal/config/**` 且全是 `*_test.go`；
  `internal/{winsec,proc,risk}` **0 hunk**、`internal/` 非 config **0 hunk**、产码 0 hunk。
- 尺 3（两族拒绝用例在哪）：`assertRefused113` / `leafLinkTo118` / `resolveProbeRoot` 在 `internal/config/` **0 命中**
  => 票 113/119 那两族**根本不在本包**（在 `internal/winsec/`），本腿结构上碰不到它们。
  本包内唯一的拒绝腿用例户 `c26_seam_posix_125_test.go` 的 git log 只有 `4824bb82`/`ff3faf92`，**本腿未碰**。
- 尺 4（决定性那条）：33 枚 went-green 在**改前普通形**就已经逐名 PASS => 它们的断言在改前本就成立，
  这次摘掉的只是 harness 形状造成的那次拒。**放宽断言换绿在这个形状下无从发生。**

### 4) 变异/还原/原样三读数的可信度 — **变异落地与红名都复现；"还原"缺 sha 证词**

- 我自己重打一发变异（`MUTATION-124-V1`：把 helper 体改回 `return t.TempDir()`，保留 `_ = proc.SealableRoot` 守 import）：
  先证落地 `grep -c 'proc.SealableRoot(t.TempDir())'` => **0**、14 调用点仍在、`go build ./internal/config/` **rc=0**、
  `go vet ./internal/config/` **rc=0**，再读数 => **MUT-LINK 顶层 FAIL=41 / 子 FAIL=1**（合计 42），mech 41。
- 差集：PRE-LINK 那 35 名**逐名全数回归**（`comm -23` 0 行）；额外翻红 **6 名**，与实现腿 sec 6 点名的 6 名**逐名相同**。
  => 它那个 `42`（41 顶层 + 1 子）我自己量到同一个 42。**变异读数可信。**
- 它的处置我判**对**：第一发 `go vet` rc=1 是 unused import 的编译错，**不是行为读数**；补 `_ = proc.SealableRoot`
  保 import 后 rc=0 再读，这是正确做法（否则把"编译不过"读成"变异生效"就是假读数）。
- **扣分项**：派单要求"每发读完立刻单独还原并核 `sha256` 回 HEAD blob"，`124-r1-impl-part2.md` **没有任何 sha256/sha1 行**。
  我用 `git status --porcelain -- internal cmd` **空** 替代证明它最终还原了（共享树里这枚文件与 HEAD blob 逐字节同）——
  结论不变，但**它的还原证词本身不合格**，属交件缺陷。
- **它自报的自相矛盾**：sec 5/6 写 `RESTORED-LINK 194/191/2/1` 与 `POST-LINK 194/192/1/1`，随后 sec 5 结论句写
  "RESTORED-LINK 与 POST-LINK 同形同数" —— **191/2 与 192/1 不同数**。这两发差的正是那枚摇摆名，即 flake 从一发挪到另一发，
  它把 flake 读成了"还原"与"原样"的差。

### 5) 票面前提新鲜度（硬闸）— 那把尺**确认恒等于已修态**，AC#2b 的票面判据今天必须换尺

- 对 **HEAD blob**（`git grep ... HEAD`，NOT 工作树）现量「`winsec.` 与 `TempDir` 同行」=> **20 行**。
  逐行分类：**18 行是 `winsec.SealableTempDirForTest124(t)` 调用**、**2 行是注释**
  （`ancestor_separator_108_other_test.go:19`、`placement_symlink_113_other_test.go:25`）。
  单行形 `winsec\.[A-Za-z0-9_]*\([^)]*TempDir` => **0 命中**；`winsec.` 同行 `t.TempDir|os.TempDir` => **0 命中**；
  两行窗口内 `winsec.` 紧邻裸 `TempDir` => **0 命中**。
  => **派单那条腐坏告警成立**：这把尺量的 100% 是"已转换后"的形态，永远读到 20，**再也测不出"还剩多少没转"**。
  * 一并更正一处转述：`124r1-ruling-batch.txt:5` 写"19 处调用、1 行注释"，**实测是 18 处调用、2 行注释**（总数 20 相同）。
- **AC#2b 的票面判据今天该按什么量**：按票面 `:81`/`:169` 逐字那条**终判据**量，而不是按任何"出现次数"尺量 ——
  即**纯净快照 + 两形各一枚新容器 + 逐包 `-count=1 -v`，对两形名册作 `comm` 差集**，取
  `only-in-link 红名数 = 0`（扣票面钉死的豁免名单）。理由：这是唯一今天仍然**会动**的尺（本包 35 => 0 就靠它量出来），
  而"`winsec.` 邻 TempDir 出现次数"这类静态尺已经饱和、分母恒 20。

### 6) 红名册逐名 — 见 part3 sec 五；新增红 **0 名**

### 7) 门禁 — 我自己现跑，见 part2 sec 四；五发全 rc=0
