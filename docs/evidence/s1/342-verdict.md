# 票 342 `AC#1` 裁决 · `342` 判语（非实现者裁决腿·只读产码·只写本件）

裁决者＝本腿（`342` 判语腿；⛔ 实现者、⛔ 产码、⛔ 翻框、⛔ 改票、⛔ push）。起手锚＝`a2b3b80c97`（`git rev-parse --short=10 HEAD`，2026-11-05 15:1x）。射程＝**票 342 两枚未勾格里本腿裁的那一格 `AC#1`**（"成对读数要有 liveness：修前分子 > 0、修后 0 且 liveness 非零，两平台各一发"）；`AC#4`（门禁＋消费者五包＋`cmd/wisp` 那一发）本腿只给射程内能取到的那几枚，⛔ 判它闭环。
★被裁对象＝落地腿 `342-r3` 的交件 `docs/evidence/s1/342-r3-impl.md`（锚 `7eb0e3f4b`；该腿自报 "`internal/winsec/**` 零字节改动、只交读数"，本腿现量复核为真）。

## 0. 起手现量（⛔ 抄腿）

| 事实 | 读数 | 我这把尺 |
|---|---|---|
| 票面框 | 未勾 **2**（`AC#1`／`AC#4`）／已勾 **5** | `grep -cE "^- \[ \] \*\*AC#"`／`"^- \[x\] \*\*AC#"` |
| `AC#1` 那两枚仪器在盘上 | `ancestoropen_342_unix_test.go:75`／`ancestoropen_342_windows_test.go:38` | `grep -n "func TestAC1POSIX…\|func TestAC1Windows…"` |
| `internal/winsec` 工作树干净 | `git status --porcelain -- internal/winsec` = **0 行** | 同命令 |
| 该包最后一笔 | `4a6e70afe`（`342r2` 的测试件收尾）＝r3 之后无 winsec 提交 | `git log --oneline -1 -- internal/winsec` |
| 被测量的产码 | `internal/winsec/winsec.go` sha256＝`f0bbf9fff8634f0cc49b2cf692ea7be57ecd2c2f722a8db3372b6f8ae49d55a1` | `sha256sum` |

台件＝**仓外导出树**（`git archive HEAD | tar -x`，`/d/tmp/wisp/20261105/adj342/tree`，⛔ 复用腿的树、⛔ `-overlay`）。Windows 档＝本机 `go1.27.1 windows/amd64`；POSIX 档＝`golang:1.27` 容器、`CGO_ENABLED=0`、`GOPROXY=off`、宿主 module 缓存只读挂、`GOCACHE=/tmp/gocache`、`TMPDIR` 留容器内 `/tmp`（⛔ 绑 Windows 卷）。
**修前树怎么造的（本腿自己造，⛔ 用腿的）**：HEAD 字节**只换 `privateFile` 一个函数的体**，函数体逐字取自 `fde1df02:internal/winsec/winsec.go`（按名 `os.OpenFile(resolved, …)` ＋ 两处 `_ = os.Remove(resolved)` 的拒绝路径）；造完逐枚自证：`os.OpenFile(resolved` 命中 **1**、`pinnedParent(path, resolved)` 命中 **0**、`_ = os.Remove(resolved)` 命中 **2**、`go build ./internal/winsec/` **rc=0**，该件 sha256＝`13f5fced02c28b7a04e15cbb8086da0592f774d2c49e6ee0b19f4ab4265f0b00`；收尾用 `git show HEAD:` 还原，sha256 逐字回 `f0bbf9ff…`＝**match=YES**。

## 1. `AC#1` 成对读数（四数只在 `^=== RUN`／`^--- PASS`／`^--- FAIL`／`^--- SKIP` 取；分子逐发点名）

### 1.1 POSIX（`TestAC1POSIXPrivateFileThroughASwappedAncestorNeverTruncatesAForeignFile342r1`，1500 发/发）

| 树 | truncated | deleted | pinnedParentBeforeOpen | provedLinkCycles | rc | 判 |
|---|---|---|---|---|---|---|
| HEAD（修后）第 1 发 | **0** | 0 | **32** | **239** | 0 | PASS |
| HEAD（修后）第 2 发 | **0** | 0 | **21** | **195** | 0 | PASS |
| HEAD（修后）第 3 发 | **0** | 0 | **30** | **183** | 0 | PASS |
| 修前（`treeOld`，删支在）第 1／2 发 | 0／0 | **97／89** | **0／0** | 204／200 | 1 | RED（见第 3 节①） |
| 修前隔离（删支摘掉）第 1／2 发 | **9／5** | 0／0 | 0／0 | 185／190 | 1 | RED（截断支的修前分子 > 0） |

- 修后三发的 liveness 与零**同发在场**：`pinnedParentBeforeOpen=32/21/30`（父钉在 open 之前答话）＋ `provedLinkCycles=239/195/183`（racer 同发证明祖先真到达过链接态），`truncated=0 deleted=0 narrowed=0 untouched=1500`。
- 修前隔离发那条是本表最要紧的一句：预置路线在 POSIX 上**多数先被拒绝路径的删除把名字拿走**，所以删支留着时 `truncated` 读 0、`deleted` 读 89~97；把两处 `_ = os.Remove(resolved)` 摘掉的隔离发才让那次 `O_TRUNC` 单独露出来（5/9 发）。⇒ **"修前 truncated=0"在这条路上 ⛔ 能读成"没伤"**。

### 1.2 Windows（`TestAC1WindowsPrivateFileThroughASwappedAncestorNeverTruncatesAForeignFile342r1`，两支形各 400 发/发）

| 树 | 形 | truncated | provedReparse | formBLegRefusals | rc | 判 |
|---|---|---|---|---|---|---|
| HEAD（修后） | junction | **0** | **16** | **7** | 0 | PASS |
| HEAD（修后） | symlink-dir | **0** | **85** | **27** | 0 | PASS |
| 修前 | junction | **3** | 13 | 0 | 1 | RED（损伤支＋仪器支各自点名） |
| 修前 | symlink-dir | **8** | 55 | 0 | 1 | RED |

- Windows 修前那两发是**真损伤红**：`AC#1 RED: 3 of 400 attempts found the foreign witness's body replaced, so the open inside privateFile resolved through an ancestor swapped into a link and its truncate landed in somebody else's tree`（symlink-dir 同句数 8）；同发 `provedReparse=13/55` ＝ 翻转真的到过链接态（⛔ 空跑）。修前侧 `formBLegRefusals=0` 是**应该的**（父钉那一族在修前字节上不存在）。
- 修后的 liveness 与零同发：junction `provedReparse=16`＋`formBLegRefusals=7`、symlink-dir `provedReparse=85`＋`formBLegRefusals=27`。

### 1.3 与腿报的对照（竞态读数⛔ 要求逐字相同，只核形状与量级）

POSIX 修后 `32/21/30` ⇄ 腿报 `32/16/31`；POSIX 修前 `deleted=97/89` ⇄ 腿报 `83/100`；隔离发 `truncated=9/5` ⇄ 腿报 `6/9`；Windows 修前 `3 / 8` ⇄ 腿报 `1/1/3 / 8/8/14`；Windows 修后 junction `16/7` ⇄ 腿报 `18/17/15 / 7/6/3`。⇒ **形状、方向、量级逐项复现**；具体枚数是 racer 摆动量，本腿⛔ 拿它判腿报真假，也⛔ 拿腿的数当自己的读数。

## 2. `AC#1` 判＝**通过**（两平台各一发以上，成对且 liveness 非零）

票面三件逐条对上：① 修前"外来见证被截断"分子 **> 0**（POSIX：隔离发 `truncated=9/5`，删支在时同一损伤以 `deleted=89~97` 露出；Windows `truncated=3/400` 与 `8/400`）＝量的；② 修后 **0**（POSIX 3 发全 0、Windows 两形全 0）＝量的；③ 同一发里 liveness 非零（POSIX 32/21/30 ＋ 239/195/183；Windows 16/7 与 85/27）＝量的，⛔ 拿"0 命中"当凭据。⇒ **这一格今天成立**，可翻（翻框归编排者，本腿一枚没翻）。

## 3. 本腿现量抓到的具名缺陷（⛔ 改＝只报；⛔ 影响第 2 节的结论，但下一位收尾前要处置）

1. **`ancestoropen_342_unix_test.go:80/:90/:91` 的归属支说反了，且它把"修前行"挡在损伤门之外**〔量的，两发复现〕。盘上 `landed := st.pinnedParentLeg > 0`；`if !landed` 判红，红句逐字写 "the parent pin answered **0** attempts BEFORE any open, so this row was produced by bytes that **DO** carry the pinned route"——`pinnedParentBeforeOpen==0` 恰恰是"这些字节**没有**父钉"。⇒ 后果：真修前树上那行是**归属红**（`:91`），损伤门（`truncated>0 || deleted>0`）**根本没被走到**；"修前分子 > 0"只能当 `t.Logf` 的**打印读数**读（本腿就是这么读的），并被腿报第 4.1 节具名。改它属判据面＝另批。
2. **`winsec.go:149-154` 的自陈已和同一函数的产码相反（欠说，⛔ 过说）**〔读的〕：`:149` 仍写 "the ANCESTOR leg is still open on POSIX … This window is not claimed closed anywhere"，而 `:194` 传的是调用者拼写、`:170-200` 落的正是父钉路由、`:202` 的 open 在钉住的父 fd 上。⇒ 今天两平台读数都齐（第 1 节），这段该在结案时改成终态；⛔ 改成"窗口已关"以外的任何"已闭"断言仍须两平台读数在册（今天在册）。腿报第 4.2 节同判，本腿独立复核为真。
3. **共享工作树的 `go vet ./internal/winsec/` 今天 rc=1，红因⛔ 在 winsec**〔量的〕：报错落在 `internal/plugin/c24.go:526`（`bool` 当 `error` 返回）与 `goja_runtime.go:241`（`undefined: scanC24Call`）＝**别腿未入库的在飞件**（`git status --porcelain -- internal/plugin`＝`M goja_runtime.go`＋`?? c24.go` 等 4 项，⛔ 跟踪面）。⇒ 本腿读门禁用**干净导出树**：`go build ./internal/winsec/` rc=0、`go vet ./internal/winsec/` rc=0、`GOOS=linux go vet` rc=0、`GOOS=darwin go vet` rc=0、`go build ./internal/...` rc=0、`gofumpt.exe -l internal/winsec/` **0 行 rc=0**、`go test ./internal/winsec -count=1`＝`ok … 22.819s` **rc=0**。⛔ 把共享树上那两行读成 winsec 的红。

## 4. 逐 AC 表（本腿只裁 `AC#1`；其余＝维持既有判语，并标今天能不能被推翻）

| AC | 判 | 依据 |
|---|---|---|
| `AC#0` 三形代价表 | **通过**（维持 `342-v1`／编排者） | 本腿未重打；三形表在 `.scratch/wisp/probes/342/a1/census.md`，选形＝乙已落地（`pinnedParent`＋`openPinnedLeaf`／`unlinkPinnedLeaf`） |
| **`AC#1` 成对读数＋liveness＋两平台** | **★通过（本腿＝决定性格）** | 第 1 节两平台成对读数：POSIX 修后 3 发 0＋liveness、修前 89~97／隔离 9/5；Windows 修后两形 0＋liveness、修前 3 与 8。造树/还原逐枚自证（第 0 节） |
| `AC#2` 覆盖写语义 | **通过**（维持） | `:96` 旗标行与"名字已存在⛔ 构成拒绝理由"两段一字未动（本腿读 `sed -n '90,110p'`）；`openleaf_overwrite_322_test.go` 在整包一发里绿（`go test ./internal/winsec` rc=0） |
| `AC#3` 去牙反控两发落点 | **通过**（维持 `342-v2`／`342-v3`） | 本腿未重打去牙（不在派单射程）；落地字节与那两判时的 `winsec.go` 同族，`unlinkPinnedLeaf`／`openPinnedLeaf` 现读在位 |
| `AC#4` 门禁＋消费者不推红 | **不通过（仍未勾）** | 射程内取到：构建／三档 vet／gofumpt／整包 `./internal/winsec` 全绿（第 3 节③，干净树）。**取不到**：票面点名的"消费者五包成对＋`cmd/wisp` 那一发"（派单⛔ 跑、且 `internal/risk`⇄`config` 归别腿车道）；另 HEAD 上 `sh scripts/d22scan.sh` **rc=1**（红因＝`internal/panel/composer_handlers.go:33` 的 phantom-citation ＋ `frontend/` 射程扫 0 件，最后被动在 `8e0702fdb`／`b4ee24cbe`，⛔ 342 的写面）、长度门 **rc=1**（既有 roster 漂移）＝收结前必须先归位这两枚 |
| `AC#5` 注释与票面口径说实话 | **通过**（维持；见第 3 节②） | 今天没有"窗口已关"这类过说；反过来 `:149-154` 是**欠说**（同函数已钉父），属结案时要更正的终态文字，⛔ 影响本格当年那句裁 |
| `AC#6` 纪律格 | **通过**（维持） | 票面框 5 勾／2 未勾＝r3 与各腿一枚没翻；`internal/winsec/**` `git status` 0 行＝r3 真零字节 |

## 5. 一句话终态

票 342 `AC#1` **成立**：同一枚仪器、两平台，修前分子 > 0（POSIX `deleted=97/89`、隔离 `truncated=9/5`；Windows `truncated=3` 与 `8`，各 400 发）、修后同一发里 **0 且 liveness 非零**（POSIX `pinnedParentBeforeOpen=32/21/30` ＋ `provedLinkCycles=183/195/239`；Windows `provedReparse=16/85` ＋ `formBLegRefusals=7/27`），满是票面三件。⇒ **`AC#1` 可翻**；但 **342 尚⛔ 结案**：`AC#4` 仍开着（消费者五包＋`cmd/wisp` 那一发未在射程内，且 HEAD 的 d22scan／长度门两枚红要先归位），另第 3 节① 那枚归属支说反、第 3 节② 那截自陈欠说，须在结案时一并处置（都属判据／产码注释面＝另批，⛔ 本腿改）。

## 6. 日志清单（仓外，逐枚非 0 字节）

`win-after.log 6839`／`win-before.log 6648`／`win-fullpkg.log 56`／`posix-after-1..3.log 2750/2561/2558`／`posix-before-1..2.log 3306/3317`／`posix-oldnounlink-1..2.log 3190/3386`（路径 `/d/tmp/wisp/20261105/adj342/logs/`；`pristine-winsec.go`＋`oldfunc.txt`＋`new-winsec.go` 同目录）。
⚠ 按派单"只写两枚判语件"，本腿⛔ 把 probes 落进 `.scratch/wisp/probes/342/`；要长期留证由编排者定。

---

## 7. 2026-11-07 续判批次：`342-r3` 交件（`356bbfcb9`）收件复裁 —— `AC#1` PASS（维持）、`AC#4` PASS（升格）

裁决者＝票 342 裁决腿 2026-11-07 批次（非实现者；⛔ 翻框、⛔ 产码、⛔ push；票尾 Progress log 记 agent=`342-v4`）。起手锚＝`dfcd937821`（12:0x 现认；`git ls-tree -r 356bbfcb9`⇄同锚 internal/winsec 56 枚 blob 同串＝与交件同字节）；收尾快照 HEAD＝`c56d8d90b3`（`git diff --name-only dfcd937821..c56d8d90b3 -- internal cmd`＝`cmd/wisp/resident_windows.go` 一枚＝246 腿，我读数涉及的包零漂移）。被裁对象＝`356bbfcb9`＋`docs/evidence/s1/342-r3-paired.md`。

### 7.0 收件三把尺（本腿自跑）

1. 提交真在＝`356bbfcb9` 在 log（3 件：票尾 +1、paired.md +50、`ancestoropen_342_unix_test.go` 25+/22-）；winsec 漂移尺 `git diff --name-only 766bf71f4..c56d8d90b3 -- internal/winsec`＝恰该枚仪器件一枚 ⇒ **`winsec.go` 两平台产码零字节属实**；HEAD `winsec.go` sha256＝`f0bbf9ff…` 与第 0 节记录逐字同。
2. 件真在＝`342-r3-paired.md` 6523 字节非 0。
3. 决定性＝本腿自建仓外台树三棵（`/d/tmp/wisp/20261107/342v/`，`git archive`，⛔ 复用 r3 的树、⛔ `-overlay`）：treeLand（落地原样）／treeMutA（一行 `pinnedParent(path, resolved)`→`pinnedParent(resolved, resolved)`＝回退票 400 hunk）／treeMutB（`fde1df02` 的按名 `privateFile` 函数体原样换入＝含两处按名 `os.Remove(resolved)` 的忠实 pre-r1 形）；三树 `go build ./internal/winsec/` 各 rc=0、`diff -rq` 各只差 `winsec.go`。本节读数全部自己重打，⛔ 抄腿。★r3 仪器修正（判语第 3 节①那枚"归属支说反"）本腿逐行读 diff＝触发条件与红句按实测改写、**损伤门与 liveness 门一字未动＝断言零放宽**；门行为不变由 7.3 的 M-B 照红坐实。

### 7.1 ★"容器 golang:1.27 首跑"核验＝真跑过

腿的台树与日志在盘（`/d/tmp/wisp/20261107/342r3/`：treeLand/treeMutA/treeMutB＋三枚 tar＋`posix-land-v3.log` 58039 字节 313 行，mtime 11-07 11:42，入库前 4 分钟）。本腿对那份日志自己数：四数＝`RUN=76 PASS_TOP=48 FAIL_TOP=0 SKIP=0`，与腿报**逐字同**；包行 `ok … 3.288s`。其 AC1 行读数＝`provedLinkCycles=155 pinnedParentBeforeOpen=20 refusedAfterOpen=10 truncated=0 deleted=0 untouched=1500`。
★如实记一笔：腿报引用的 `270/36/37` 与盘上存档的 `155/20/10` 是**同字节不同发**（竞态摆动量，形状逐字同＝零＋liveness 在场；日志名 v3 说明前两发未留档）。判：首跑属实、读数诚实；"引用发与存档发不逐字对应"归日志留档纪律，⛔ 计入判据。

### 7.2 `AC#1` 成对读数——本腿独立复跑（两平台各两发＋全包）

POSIX 档（容器 `golang:1.27`、`CGO_ENABLED=0`、`GOPROXY=off`、宿主 module 缓存只读、`TMPDIR`/`GOTMPDIR` 留容器内 `/tmp`）：

| 树 | truncated | deleted | untouched | liveness（同发） | rc |
|---|---|---|---|---|---|
| treeLand 发 1 | 0 | 0 | 1500 | provedLinkCycles=177、pinnedParentBeforeOpen=20 | 0 |
| treeLand 发 2（全包内） | 0 | 0 | 1500 | 172/19 | 0 |
| treeMutB＝M-B（修前） | 1 | **69**/1500 | 1430 | pinnedParentBeforeOpen=**0**（refusedAfterOpen=32） | 1 RED |
| treeMutA＝M-A | 0 | 0 | 1500 | 180/20 | 0（见 7.3） |

全包四数（treeLand）＝`76/48/0/0`、rc=0。

Windows 档（本机 go1.27.1 windows/amd64，跑前 unset GOTMPDIR/TMPDIR/TEMP/TMP）：

| 树·形 | truncated | liveness（同发） | rc |
|---|---|---|---|
| treeLand 单发·junction | 0 | provedReparse=11、formBLegRefusals=3 | 0 |
| treeLand 单发·symlink-dir | 0 | 60/11 | 0 |
| treeLand 全包·junction／symlink-dir | 0／0 | 10/4 与 58/13 | 0（全包四数 127/76/0/0） |
| treeMutA·symlink-dir | **1**/400 | 73/14 | 1 RED（junction 该发 0＝摆动量） |
| treeMutB·symlink-dir | **6**/400 | formBLegRefusals=**0** | 1 RED（损伤红＋仪器红同发） |
| treeMutB·junction | 0 | formBLegRefusals=**0** | 1 RED（仪器红） |

★M-B 红句与漏出 body 逐字入库：Windows `AC#1 RED: 6 of 400 attempts found the foreign witness's body replaced…`，body＝`"raced private payload 342r1 windows ancestor"`（产码 44 字节载荷＝342-v2 三凭据的形状独立复现）；POSIX 红在归属门（修正后 wording）`AC#1 RED (not a landing reading): pinnedParentBeforeOpen answered 0 attempts…leg 342-r3 re-took deleted=69 of 1500 on the by-name route…`——deleted=69 落回 342-v1 VPROD/VBASE 的 65..86/1500 量级带。

**判＝`AC#1` PASS（维持 2026-11-05 第 2 节）**。票面三件全部本腿现量：①修前分子>0（M-B：POSIX deleted=69/1500、Windows truncated=6/400）；②修后 0（两平台各两发）；③同发 liveness 非零（POSIX 177/20、172/19；Windows 10/4、58/13）。零命中⛔ 当凭据一事由变异侧同发反向坐实（下节）。

### 7.3 ★变异有牙——本腿亲放 M-A/M-B，还原全绿

| 变异 | Windows | POSIX |
|---|---|---|
| M-A（回退票 400 hunk） | **红**：symlink-dir truncated=1/400（junction 0＝摆动量） | **绿（⛔ 当凭据）**：`resolved==path` 时该 hunk 构造性 no-op——本腿机制确认＝父钉腿照答 20 发（pinnedParentBeforeOpen=20）、损伤 0 |
| M-B（pre-r1 按名路由） | **红**：truncated=6/400＋RED instrument（formBLegRefusals=0）；junction 侧同发仪器红 | **红**：deleted=69/1500、pinnedParentBeforeOpen=0（归属干净） |
| 还原（treeLand） | 绿（7.2 表） | 绿（7.2 表） |

⇒ 同仪三态齐：修前字节红、落地字节绿、还原复绿。腿报"自抓"（中性化 removeRefusedLeaf 会假绿、变异须连损伤通道一起还原）本腿认可——M-B 用忠实 pre-r1 函数体正是为此；本腿四格读数无一假绿。

### 7.4 `AC#4` 门禁＋消费者——本腿现跑（共享工作树，各门前 `unset GOTMPDIR TMPDIR TEMP TMP`）

| 门 | 读数 |
|---|---|
| `gofumpt -l internal/winsec`（真身 /d/work/base/gopath/bin/gofumpt.exe） | 0 行 rc=0 |
| `gofumpt -l` 五包＋winsec | 0 行 rc=0 |
| `go vet ./internal/winsec/` 本机／GOOS=darwin／GOOS=linux | 各 rc=0 |
| `go vet ./cmd/wisp/` 本机 | rc=0（两发：`dfcd937821` 时与 `c56d8d90b3` 现认） |
| `GOOS=linux`／`GOOS=darwin go vet ./cmd/wisp/` | 各 rc=1＝**票面逐字点名的交叉档机制**（Windows 宿主交叉 cgo 关：linux 报 `sherpa-onnx-go-linux`、darwin 报 `sherpa-onnx-go-macos`，`build constraints exclude all Go files`）＝⛔ 新红，与本机 rc=0 成对 |
| `sh scripts/d22scan.sh` | rc=0 clean（bans#1-5 internal/=376 cmd/=58、ban#8 internal/=934 cmd/=224） |
| ★`sh scripts/check-path-length-budget.sh` | **rc=0 VERDICT GREEN**（tracked=18914、over-budget=65、covered=65、not-in-roster=0，byte/char 双口径同 65）＝腿报那枚 RED 已被 `532f936d6`（roster 6 条过期行跟进 -done 改名）修掉，本腿现跑核绿 |

消费者五包成对（本腿自找消费者＝`PrivateFile`/`PrivateFileExclusive` 调用面所在五包：internal/risk、internal/secret、internal/config、internal/agent、internal/memory；干净 treeLand，本机档＋容器档各一发 `-count=1`）：

| 包 | Windows rc（RUN/PASS_TOP/FAIL/SKIP） | POSIX rc（RUN/PASS_TOP/FAIL/SKIP） |
|---|---|---|
| risk | 0（289/174/0/1，SKIP=TestSyncRegistryProbeLive） | 0（263/156/0/0） |
| secret | 0（33/25/0/0） | 0（11/11/0/0） |
| config | 0（238/140/0/0） | 0（237/137/0/0） |
| agent | 0（191/158/0/0） | 0（190/157/0/0） |
| memory | 0（193/103/0/2，SKIP=两枚 crash-writer） | 0（192/102/0/2，SKIP=同两枚） |

**十发全 rc=0、顶层红 0**（成对判据按 342-ac4-readings 第 1 节口径＝包级 rc＋顶层红名，⛔ 按四数相等）。★对照 2026-11-05 读数件的三处不成立今天全部消失：`TestC26RewriteAccountIsConsumedAtEverySecurityLeg`（两平台，A380 面）／`TestHandDeletedSchemaVersion…`（容器轮换，票 341 面）／`TestDsn196r3…`＋`TestMigration31V5ToV6…`（POSIX 专属，票 196/31 面）——本腿日志里四枚结果行全 `--- PASS`＝各自 owner 面的修复已落在 HEAD（权属本腿不判）。`go build ./internal/...` rc=0。

`cmd/wisp` 那一发：本波派单令⛔ 碰 cmd/wisp（246 腿在写）＝票面"⛔ 抢车道，按波次具名申请"条款，整包未跑；射程内补的读数＝本机 `go vet ./cmd/wisp/` rc=0（现 HEAD）。整包门的既有读数在册（342-ac4-readings 第 4 节：rc=1 归 A1177 本机环境上限、零 `--- FAIL`、342 邻接 9 枚子集 9/9 绿）。

**判＝`AC#4` PASS（升格）**：票面门禁清单逐枚 rc 在册且全绿（含长度门 GREEN 现核）；消费者五包两平台全 ok 零红＝"不推红"实质成立；`cmd/wisp` 那一发的缺席是纪律性缺席（车道条款）而非读数缺席，交叉档红按票面自己的口径归底。当年搁置该格的三个理由（4 枚既有红／长度门 RED／五包成对缺口）今天逐一消除。

### 7.5 本批次一句话终态与翻框建议（⛔ 本腿自己翻）

- **`AC#1` 建议翻 [x]**：两平台成对读数＋liveness＋变异三态齐，全部本腿独立现量（交件读数形状逐项复现）。
- **`AC#4` 建议翻 [x]**：门禁逐枚绿（长度门已由 `532f936d6` 修绿现核）、五包两平台十发全 ok、cmd/wisp 交叉档红按票面口径归底、整包门既有读数在册。
- 两格翻完＝7/7，本票可结案。结案时具名残留（均⛔ 本票射程）：`winsec.go:149-154` 自陈欠说（第 3 节②，两平台读数今天在册、终态化归编排者/产码主）；Windows 裸兄弟叶子形残差（票 322/353 叶子腿）；`FILE_SHARE_DELETE` 顶名窗（342-v3 具名的新票族）；`cmd/wisp` 整包本机环境上限（A1177）。

本节日志（仓外 `/d/tmp/wisp/20261107/342v/`，逐枚非 0 字节，最小 1018）：`win-land-ac1.log`／`win-land-full.log`／`win-mutA-ac1.log`／`win-mutB-ac1.log`／`win-cons-{risk,secret,config,agent,memory}.log`／`posix-land-ac1.log`／`posix-land-full.log`／`posix-mutA-ac1.log`／`posix-mutB-ac1.log`／`posix-cons-{risk,secret,config,agent,memory}.log`／`make_muts.py`／`pre-r1-winsec.go.txt`；对照读腿的首跑存档 `/d/tmp/wisp/20261107/342r3/posix-land-v3.log`。
