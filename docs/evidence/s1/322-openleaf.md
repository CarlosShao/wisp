# 票 322 — `privateFile` 那次 open 自己跟随被换的叶子链接：落地件与读数（腿 322-r1b）

> **这份文件是什么**：票 322 落地腿 `322-r1b` 的证据件（`AGENTS.md` §1.5 的落点）。
> **它不是什么**：⛔ 裁决表、⛔ 判语。判语与 AC 框归非实现者（`SPEC-12 §4.3` #1/#3）。
> 本腿**没有**翻任何一枚 AC 框（起手现量 6 枚未勾／0 枚已勾，交件时同一把尺复量仍是 6／0），
> 也**没有**动 `Status:`、票名、`-done`、台账。
> 前身＝`322-r1`（死于 150 轮上限、零提交），其字节＝`8baec646`（编排者代提、标〔未验证半成品〕）。
> 本腿读的是那六枚字节，**不是**从零再写一遍；下面每一枚数都是本腿自己现跑的，⛔ 转引。

## 1. 这一格的病（读数，不是感觉）

`privateFile`（`internal/winsec/winsec.go`）的形状是"先 `resolveString(path)` 把拼法问清楚，再 `os.OpenFile` 打开"。
解析与打开之间，叶子被换成链接时，**open 自己就跟随它**，`O_WRONLY|O_TRUNC` 当场把链接指向的那个对象截断。
票 120 那道 fd 闸门拦的是"私密字节落到外来对象"那一跳，拦不住**已经发生的截断**——所以这一格是独立的病，
`A990` 裁①因此另立一票。

判据的形状（票面）：**叶子在检查与 open 之间被换成链接时，树外那个真文件的内容不许被清空**。

## 2. `8baec646` 交的是什么形状（本腿读懂的部分）

形由台账 `A996` 裁死＝**分平台两半**（⛔ 单一形、⛔ 乙两跳、⛔ 丙只报告）。死腿交的正是这一形：

| 面 | 交的东西 | 本腿复认 |
|---|---|---|
| `winsec.go` | `os.OpenFile(resolved, flags, perm)` → `os.OpenFile(resolved, flags\|oNoFollowPlatform, perm)`，并在打开**之前**插一枚 `preOpenLeafGuard(resolved)` 钩子；打开失败那句**原样返回**、不走 `sealError`（保住票 118 那枚"拒绝不许删链接"的钉） | 逐字读过 diff（副本：`.scratch/wisp/probes/322/r1b/8baec646-prod-code.diff.txt`） |
| `winsec_other.go` | `const oNoFollowPlatform = syscall.O_NOFOLLOW`；`preOpenLeafGuard` 在 POSIX 是**有意为空**（理由：这个问题已经由内核在同一发 open 里答了，事前探针只会多加一个会过期的判决）；并把 `:82` 那句"另立票 322、本锚未修"改成说实话的半开陈述 | 即 AC#5 要的那句，**已改**，见 §5 的 AC#5 行 |
| `winsec_windows.go` | `const oNoFollowPlatform = 0` ＋ `preOpenLeafGuard`：`FILE_FLAG_OPEN_REPARSE_POINT` + `OPEN_EXISTING` 开句柄**只问一句"这个名字是不是链接"**，是则走 `preOpenRefusal`（两枚哨兵齐）；名字不存在⇒`nil`（创建路径），其它探测失败⇒**也拒**（收紧方向只许变严） | 注释里明写⛔ 拿它当"不跟随的覆盖开"（`CREATE_ALWAYS` 变体实测会把链接条目换成空普通文件） |

三枚新测试件（`openleaf_322_other_test.go` 262 行／`openleaf_guard_322_windows_test.go` 254 行／
`openleaf_overwrite_322_test.go` 131 行）＝判据本体，本腿**一字未放宽**、⛔ 删用例、⛔ `t.Skip`。

## 3. 本腿补了什么（＝死腿的残余清单，四件全做完）

1. **证据件**＝这一枚（死腿：零枚）。
2. **票尾一条 bullet**（死腿：没写）。
3. **去牙反控**＝§5（死腿：零发；编排者代提时具名留给续程）。
4. **死腿自己那把尺的读数**＝§6 起手那一发（含 `GOTMPDIR` 写法与卷）＋它改到 `3000` 之后**从未跑过**的那一发（§5 表第一列补跑）。
5. 262 行那枚文件的拆分决断＝§7（选"按内容拆"，判据一字未动）。
6. 顺手补完死腿最后一句"注释里有个词打错了正在改"：`winsec_other.go` 那句 `(oNoFollowPlatform below)` 指错方向（常量在 `:29`、在那句话**上面**）⇒ 改 `above`；
   `openleaf_322_other_test.go:9` 指向旧证据件名 `322-open-nofollow.md`（本单权威名 `322-openleaf.md`）⇒ 改正引用。
7. 死腿那 66 枚只挂在盘上的读数（tracked 0 枚）⇒ 逐枚 `git add` 入库，正文标〔322-r1 读数遗产·编排者代提〕。

## 4. AC#1 成对读数（POSIX，容器 golang:1.27，同一台件、同一发、一屏两数）

`leafRaceAttempts322r1 = 3000` 这一版（死腿只跑过 800 那一版，且那一发自己被判"仪器红"）。
两棵树：`git archive 8baec646^`（未修码，`.go=1090` 数过才信）＋同三枚测试件与两枚符号的**改前值**；`git archive HEAD`（已修码）。

```
base-322   rc=1  attempts=3000 flips=2164 damagedWitnesses=63  leafLegRefusals=0  foreignsNarrowed=0 foreignsHoldingPayload=0
fixed-322  rc=0  attempts=3000 flips=1772 damagedWitnesses=0   leafLegRefusals=48 foreignsNarrowed=0 foreignsHoldingPayload=0
```

修前分子 63 > 0 ⇒ **坏世界真被造出来**（不是"没到达链接态"）；修后 0 **且** liveness=48 非零 ⇒ 成对读数成立。
同四发（base/fixed/m1/m2）的工具词自检 `walk=0 buildfailed=0 noteststorun=0 setupfailed=0` ⇒ 没有红句是台面故障。
原始正文：`.scratch/wisp/probes/322/r1b/logs/{base,fixed}-322.body.log`＋`mut-posix-summary.txt`。

## 5. 去牙反控（首要攻击点＝判据仪器自己）

⛔ `-overlay`（这一族尺运行时读盘，overlay 恒绿＝假读数）⇒ 全部**物理改字节**，
在 `git archive` 出的**仓外导出树**里做（⛔ 拿共享工作树当改前那一发），每枚改前留 `.bak`、记 sha256、还原后复跑。

### 5.1 牙摘掉必响

| 发 | 摘法（物理） | 结果 | 必须红的用例 |
|---|---|---|---|
| POSIX m1 | `winsec.go` 里 `os.OpenFile(resolved, flags\|oNoFollowPlatform, perm)` → 去掉 `\|oNoFollowPlatform` | rc=1；`attempts=3000 flips=2402 damagedWitnesses=70 leafLegRefusals=0` | `TestAC1POSIXOpenThroughASwappedLeafLinkNeverTruncatesAForeignFile322r1` |
| Windows w1 | `winsec_windows.go` 的 `preOpenLeafGuard` 体首行注入 `return nil` | rc=1 topFAIL=2；`attempts=400 flips=268 guardRefusals=0` | `TestAC3WindowsLeafGuardRefusesLinkNamesAndAcceptsRealOnes322r1` ＋ `TestAC3WindowsPrivateFileRouteHasTheLeafGuardSpeaking322r1` |

摘前的 sha256：`winsec.go=41ba9a9b8fcf7050…`、`winsec_windows.go=5cdef88183deb576…`；
还原后实测逐字相同（`M1-RESTORE sha=41ba9a9b…`、`w1 RESTORE sha=5cdef881…`、`w2 RESTORE sha=5cdef881…`），
复跑绿：`m1-restore-322 rc=0`（damaged=0／liveness=52／flips=1856）、`w1-restore rc=0`（2 PASS）、`w2-restore rc=0`（2 PASS）。

红句逐字（POSIX，来自 `logs/m1-322.body.log`）：

```
openleaf_322_other_test.go:199: AC#1 RED: 70 of 3000 foreign witnesses lost their own bytes or mtime, so the open in privateFile followed a leaf link swapped into the caller's name and truncated a file no caller named (first: foreign-witness-0005-322r1.bin: content is "", want its own bytes). This is ticket 322's shape and the outcome the leaf leg exists to prevent.
openleaf_322_other_test.go:205: AC#1 RED instrument: across 3000 raced attempts the leaf leg ("too many levels of symbolic links") answered 0 and the racer landed 2402 flips (0 errors, first ""), so the damaged count above is not evidence about a window that was never reached
```

红句逐字（Windows，来自 `logs/w1.body.log`）：

```
openleaf_guard_322_windows_test.go:127: AC#3 RED file-symlink: preOpenLeafGuard accepted a name carrying a reparse point, so the open that follows it can still truncate whatever that link targets
openleaf_guard_322_windows_test.go:129: AC#3 RED directory-symlink: preOpenLeafGuard accepted a name carrying a reparse point, so the open that follows it can still truncate whatever that link targets
openleaf_guard_322_windows_test.go:131: AC#3 RED junction: preOpenLeafGuard accepted a name carrying a reparse point, so the open that follows it can still truncate whatever that link targets
openleaf_guard_322_windows_test.go:241: AC#3 RED instrument: 400 raced attempts produced guard-leg refusals=0 and racer flips=268 (flipErrs=0, firstFlipErr=""), so nothing below proves the guard is on this route
```

### 5.2 反形也必响（⛔ 只做"改前必红"＝恒真三形之一）

把守卫换成"逢名就拒"（POSIX m2：`preOpenLeafGuard` 体改成为任何名字都返回带 `ErrUnresolvedPath` 的拒绝；
Windows w2：体首行注入 `return preOpenRefusal(path, "MUTANT w2: blanket refusal")`）：

- POSIX m2：rc=1。AC#1 那枚**仍然红**——不是因为它看到截断（`damagedWitnesses=0`），而是因为它拒绝把"有任何 error"当成窗口关了：
  `attempts=3000 flips=350 damagedWitnesses=0 leafLegRefusals=0` ⇒ 仪器行自己红。
  同发 `TestAC2OverwritePromiseHoldsForEveryProductionCaller322r1` 六枚子项全红（原句逐字见 `logs/m2-322.body.log`，形如
  `AC#2 RED: PrivateFile refused an existing name (… MUTANT blanket refusal on every name, fixed or not: …), which breaks storing the same id twice replaces the blob in place`）。
- Windows w2：rc=1 topFAIL=2。`guardRefusals=0 otherRefusals=400` ⇒ 路由那枚用例**按守卫自己的原句词判**（`whyGuardRefusal322r1`），
  反形在它眼里等于"另一条腿答了话"，照样红；直调那枚用例另在"真文件不该被拒""名字还不存在不该被拒"两行红。

⇒ 两平台的判据都不是恒真尺：摘牙红、反形也红，只有"叶子那一支被正确回答"那一形绿。

### 5.3 一枚单边尺（具名，不是缺世界）

`TestAC322LeafFlagRefusesALeafLinkWithoutTouchingTheTarget322r1`（现居 `openleaf_flag_322_other_test.go`）
在 m1 与 m2 两发里**都绿**。原因写进了它自己的文件头：它自己按名 `os.OpenFile`，**不经 `privateFile`**，
所以它钉的是 `oNoFollowPlatform` 的内核效果（ELOOP、拒绝不删链接、悬挂叶子不创造树外对象），
接线由 §4 那发 race 钉。⇒ 分工正确，但**验收腿别把它当接线的凭据**。

**缺世界＝零枚**：三枚新测试件里没找到"夹具根本没造出坏世界"的那一支——
POSIX 有（base 63 枚见证被截）、Windows 有（w1 的 268 次翻转、`witnessTruncated=true`）、AC#2 每枚调用者各起一个"名字已经站在那儿"的世界。

### 5.4 本腿量到、但⛔ 由本腿判生死的一枚残留（Windows）

`w-fixed`／`w1-restore`／`w2-restore` 三发里 `AC#3 local paired reading` 那一行都写着 `witnessTruncated=true`
（`witnessSize=0` 一枚、`=35` 两枚＝外来见证被截空、甚至被换上 raced payload）。
这不是去牙，这是**已修码上本机量到的残留**：探针答完话到 `os.OpenFile` 之间那扇窗今天仍然能截断外来文件。
用例把这一步**写成报告而不是闸门**（`openleaf_guard_322_windows_test.go:249-253` 逐字："Reported, because the probe cannot lock its own answer"），
本腿⛔ 改它的判据、也⛔ 替它翻框；只把数抄这里，交验收腿与机主拍"这一格算不算已闭合"。

## 6. 起手读数与门禁（每把尺自己落 `rc=`）

- 起手那一发（`go test ./internal/winsec/ -count=1 -v`，本机 `go1.27.1 windows/amd64`，
  **`GOTMPDIR` 未设**、`TMP` 解析到 `C:\Users\swq\AppData\Local\Temp`＝**C 卷**）：
  `RUN=110 / PASS_TOP=61 / FAIL_TOP=0 / SKIP_TOP=0`，rc=0（口径具名：`^=== RUN` 数顶层枚、`^--- PASS` 只数顶层）。
  ★归属：编排者 09:0x 那发的两枚红（`TestGateJudgesThePrivateSetByResolvedSID`、`TestSealNarrowsAndNamesThePrincipalItRemovedBySID`）
  在本腿这一发起手就**已经是绿**（同 C 卷、同默认 TMP）⇒ ⛔ 算本腿功劳、⛔ 算本腿病、那两枚用例**一字未动**；
  本族颜色随卷/残留摆动这一维按派单具名转给验收腿。
- 拆分后复跑同一把尺：`RUN=110 / PASS_TOP=61 / FAIL_TOP=0 / SKIP_TOP=0`，rc=0 ⇒ 四数一字不变。
- 容器全量 `-count=2 -v`（HEAD＋工作树字节）：`RUN=142 / topPASS=86 / topFAIL=0 / topSKIP=0`，rc=0；
  两样本各自的 AC#1 同屏两数：`damaged=0 liveness=68 flips=1850` 与 `damaged=0 liveness=53 flips=1690`。
- 门禁五件（本机，`GOTMPDIR` 全程未设；具名两枚台面雷都在）：
  `go build ./...` **rc=0**；`gofmt -l internal/winsec/` **rc=0／列出 0 枚**（本机没装 `gofumpt`：`which gofumpt` 无命中、rc=127 ⇒ 用具名更窄的 `gofmt -l` 代跑并写明差别）；
  `go vet ./internal/winsec/` **rc=0**；`sh scripts/d22scan.sh` **rc=0**（clean；⚠ ban #8 的 `internal/` 覆盖面 571→**572** 枚 Go 文件＝本腿新增那枚测试件，阈值与词表一字未动）；
  `sh scripts/check-path-length-budget.sh --with-self-test` **rc=0**（VERDICT GREEN，over=60／roster=60）。
  跨档 vet：`GOOS=darwin`／`GOOS=windows`／`GOOS=linux` 三档各 **rc=0**。
- AC#4 消费者色（容器 Linux 档，同一发）：`internal/risk`／`internal/secret`／`internal/config`／`internal/agent`／`internal/memory` **全 ok**，rc=0。
  ⚠ 台面自纠：本腿第一发把 `GOMODCACHE` 挂成命名卷 `probe322mod:ro`，五包全 `[setup failed]`（`mkdir /gomod/cache: read-only file system`）
  ＝台面故障、⛔ 计入判据；按 `322-r1/cons.sh` 的挂法（宿主 `/d/work/base/gopath/pkg/mod` 只读）重跑才是上面那一发。

## 7. AC 格子逐枚推进（⛔ 翻框；每格答"这一发在未修码上响不响"）

- **AC#0 三形代价表**：`322-a1` 已交（`.scratch/wisp/probes/322/a1/`），票面 progress log 有读数摘要；形由 `A996` 裁死并已落。未修码上响不响＝本腿 §4 的 base 行（63 枚）现量复认。
- **AC#1 截断窗口**：**叶子那一支**在 POSIX 关着（§4：修前 63／修后 0 且 liveness 48；§5.1 摘牙复响 70）。
  **祖先那一支今天仍开着**——票面、产码注释（`winsec.go` 的"ANCESTOR leg is still open on POSIX"、`winsec_other.go` 末段）与本件一律具名写着，⛔ 任何"窗口已关"。
  Windows 那一半见 §5.4：探针只答不锁，本机仍量到截断 ⇒ 这一格**不由 POSIX 绿冒充**。
- **AC#2 覆盖写语义**：全仓现量 `privateFile` 的 `:87` 那一支有**五枚**生产调用者（本腿自己 grep）：
  `internal/agent/spill.go:166`、`internal/config/migrate.go:93`、`internal/secret/migrate.go:169`、`internal/secret/migrate.go:189`、`internal/secret/store.go:79`。
  其中**四枚**依赖"或替换"，各有一枚钉（`openleaf_overwrite_322_test.go`），摘掉必红＝§5.2 m2 那发实证。
  第五枚 `secret/migrate.go:169` 不依赖覆盖（它被 `os.Stat` 的 `ErrNotExist` 挡在前面，只在备份不存在时写），
  它依赖的是"名字还不存在⇒不拒"那一支——Windows 直调用例的 `absent` 行正钉着它（w2 那发里它红过，证有牙）。
  未修码上响不响：base 那发里 `TestAC2Overwrite…` **绿**（覆盖写本来就没坏，坏的是截断）⇒ 这枚钉不为本票而红，这是**正确的**。
- **AC#3 Windows 那一半**：本机三形（file symlink／dir symlink／junction）真量到、真拒；`[仅本机可量]` 口径保留（容器是 linux，CI runner 的软链接权限档从未测＝`322-a1` 未测格 3）。⛔ 用 POSIX 绿冒充；残留见 §5.4。
- **AC#4 门禁＋消费者不推红**：§6 全件齐，唯 **`cmd/wisp` 那一维本腿零读数**——派单具名票 314 整包死亡区、只许定向 `-run` 且同刻一枚腿，本腿⛔ 碰；
  这一格缺的那一发**具名交回编排者**（不是"没问题"，是没跑）。
- **AC#5 注释说实话**：`winsec_other.go` 那句"booked separately as ticket 322 and is unfixed at this anchor"已被改成"叶子已关／祖先仍开"的陈述（`8baec646` 交的），
  本腿复量它的方位词错误并改正（`below`→`above`）；`winsec.go` 的新注释三段各自点名残留（POSIX 祖先、Windows 探针不锁答案、Windows 祖先），
  全仓⛔ "窗口已关"字样（本腿读 diff 时逐枚看过）。未修码上响不响＝注释不响，它是口径。

## 8. 本腿自己的偏差（具名自纠，不等谁来抓）

1. 第 0 步第 4 项那笔代提（`1ef8c485`）里，`git add` 做到了路径逐枚（66 枚），但**同一笔的 `git commit` 我用了一枚目录形 pathspec**
   （`.scratch/wisp/probes/322/r1/logs`）——违反派单"`git add`／`git commit` 两处逐枚显式 pathspec"。
   事后尺：`git show --name-only` 该笔列出的**恰好是那 66 枚、非该目录者 0 枚**⇒ 没有卷入别人的字节，但**形状违规成立**，
   从这里往后每一笔（`d9e1498b`／`d425d951`／`1fbc9067`）都是逐枚 pathspec，已复量。
2. `mut-posix.sh` 第一版把 base 树导出自 `HEAD`，会与 `HEAD` 里已有的同名符号重复声明⇒ 已在脚本里改成 `8baec646^` 并把踩坑原因写进注释；
   改前那一发的 `base` 树从未产生过读数（它连编译都没过），⛔ 被上面任何数引用。
3. 本枚证据件第一版 §9 的**枚数自己算错**（"logs/ 21 枚／00~11 十一枚"，现量＝`logs/` 25 枚、顶层 24 枚、共 49 枚／680790 字节）
   ⇒ 由本腿复量后改正，改正只动文档计数、⛔ 动任何判据与读数；这一条按"自抓不等谁来抓"记在这里。
   ⚠ 票 322 那条 bullet（本腿同批写的）里也带着错的那个数"logs/ 21 枚"⇒ **不回改**（Progress log 只追加、且本单只许追加**一条** bullet），
   以本件 §9 的现量为准；编排者翻框要引枚数请引这里。

## 9. 落点索引（本腿交的全部字节）

- 证据件：`docs/evidence/s1/322-openleaf.md`（这一枚）。
- 读数件：`.scratch/wisp/probes/322/r1b/`——共 **49 枚文件、680790 字节**（`find -type f`＋`du -sb` 现量）：
  顶层 24 枚（`00`～`11` 的门/env/色/四数件 13 枚、`8baec646-prod-code.diff.txt`、`mut-posix.sh`／`mut-win.sh` 两台件、
  两份 driver 日志、六枚提交正文 `msg-*.txt`）＋`logs/` 25 枚（容器壳与容器内正文、两族 summary、`final-count2`/`cons2`）。
  ⚠ 本枚证据件第一版这里写的是"logs/ 21 枚／00~11 十一枚"＝**枚数算错**，由本腿自己复量后改正（改动只在文档计数，⛔ 判据、⛔ 读数）。
  树只建不删：`/d/tmp/probe322r1b/{base,fixed,m1-teethout,m2-blanket,split-worktree,final-count2,wfixed,w1-teethout,w2-blanket}.*`。
- 产码面：`internal/winsec/openleaf_flag_322_other_test.go`（新，78 行）、
  `openleaf_322_other_test.go`（214 行，只剩 race 一族＋常量）、`winsec_other.go`（一枚方位词）。
- 票面：`.scratch/wisp/issues/322-…md` 末尾**追加一条** bullet（⛔ 翻框、⛔ 动 Status/票名）。
- 提交：`1ef8c485`（死腿 66 枚读数遗产）／`d9e1498b`（起手三把尺＋两件台件）／`d425d951`（去牙读数 20 枚）／
  `1fbc9067`（拆分＋注释改正＋门禁五件）／`443a2fa7`（这一枚证据件＋AC#4 余下读数）／`a12eb185`（票尾一条 bullet）／
  本腿最后一笔（枚数自纠＋提交正文入库）。⛔ push。
