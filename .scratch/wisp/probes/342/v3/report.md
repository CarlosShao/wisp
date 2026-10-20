# 342-v3 交件摘要（票 342 续体终裁腿＝非实现者·只读产码·两票同批＝实现腿 `342-r2` 的 10 笔交件）

判决＝`docs/evidence/s1/342-v3-verdict.md`⇄票 348 专属摘要＝`.scratch/wisp/probes/348/v2/readings-348.md`。
起手锚现认 `677d4837`；台件＝自建四棵仓外导出树（基线 `fde1df02` ⇄ 落地 `677d4837` ⇄ 落地⛔ 本腿字节 `3abbb639` ⇄ 去牙专用 `treeT`，八件套＋`docs/`，`.go` 766⇄777⇄775 非空尺）＋★三枚本腿自写仓外台架（`sh/06`/`sh/07`/`sh/08`＝纯平台原语、⛔ 任何产码字节、⛔ 模拟）；`go1.27.1 windows/amd64`；`GOTMPDIR`⇄`TMPDIR`＝`%LOCALAPPDATA%\Temp\v3-342v3`（`cygpath -w` 短形）；★⛔ `cygpath -wp` 整条 PATH（本机 gcc 真身 `E:\work\base\msys64\mingw64\bin\gcc.exe`，共享树跑色⛔ 这样转）。

## 三件最硬（盘上凭据）

1. ★**第二道人工批准＝本腿签 ⛔ 该批**。理由⛔ 是成本、是"批了也拿不到"：`sh/06` 把派单那枚钩子可能站的**四枚瞬间**（H1 钉后开前／H2 开完 seal 前／H3 seal 后删除前／H4 提前放钉）全打了一遍＝**祖先 rename 四枚全拒**（`Access is denied`／`in use by another process`）×3 发，同发 H0 空手＝`err=<nil>`＝正控在场 ⇒ 那枚世界**缺⛔ 是钩子、缺是内核事实**。平台锁两态另由 `sh/07` 独立复现（基线态两枚移动 `err=<nil>`×3 ⇄ 落地态 `Access is denied`＋`in use`×3、释放钉后又落得动）＋仪器件三 B/L 两行×两树各 3 发同向。⇒ 裁法按"平台句柄锁＝结构性进⛔ 去"落进判语、⛔ 悬成"待补"、⛔ 圆成绿。
2. ★**`FILE_SHARE_DELETE` 顶名窗口＝本腿签 新立一票**（⛔ 属 342/348 祖先判据族）。凭据＝`sh/08` **跨进程**攻击者×2 发：holder 握父钉时 `rename(data)`＝`Access is denied`，★同发 `remove(caller's name)=<nil>`＋`rename(mine -> caller's name)=<nil>`＋holder 执行产码那一步 `Root.Remove(base)=<nil>`＝**攻击者那枚文件的唯一名字被删**；两枚去牙把机制钉死＝摘 `unlinkPinnedLeaf` 那一步 ⇒ damage `3→0` 而 `entered` 照 3（损伤确由产码那一步、⛔ 仪器自造）／`openPinnedLeaf` 换按名开 ⇒ `entered 3→0`（窗口确由句柄带 DELETE 打开）。族内⛔ 通用＝本票族那两枚分子（`truncated`/`deleted` 走翻转祖先）在这里逐发都 0＝同族判据既⛔ 能满足也⛔ 能证伪。名与判据形状写在判语末节（⛔ 本腿立票＝归编排者）。
3. ★**liveness 用实现腿自己具名的等价钥匙、逐发点名**〔量的〕：`pinnedParentBeforeOpen`/`provedLinkCycles` 在它两枚 Windows 文件里本腿现量＝**各 0 命中**（真身住三枚 `…_unix_test.go`）＝它⛔ 借 POSIX 族、本腿也⛔ 借。本腿用的钥匙＝件四 `proved flips=56/60/56`＋`flipErrs=6/7/9`（⇄基线 `41/35/34`＋`17/18/12`）、件一 `realSealRanFirst=true`×3×**两树**＋`worlds entered=3/3`⇄基线 `0`＋`refused=3`。★AC#5 据实判＝**那段今天⛔ 算"关"**：link 面量到（`re-check=45/56/54`/400 ⇄ 基线 0×3 有牙红）但 **SameFile 面 0 命中**（落地三发逐枚 0）＋★本腿新增一味＝同发 `truncated=45/56/55` 与 re-check 计数贴住＝**复验答话、⛔ 防住**（发 C 只空 link 面 ⇒ link `→0`、SameFile 接住 37、`truncated=38` ⛔ 归零）＝`winsec.go:118-140` 的"两平台都量到"仍⛔ 满足、⛔ 任何腿有权改成"关"。

## 成对四数（★只在仓外树、⛔ 混用⛔ 相减、竞态⛔ 要求逐字相同）

落地 `119/67/1/0` ⇄ `119/66/2/0` ⇄ 基线 `114/64/1/0`×2 逐字同号 ⇄ 落地⛔ 本腿字节 `115/62/2/0`。红名册逐名＝落地两枚全姓 `342r1`（`…NeverTruncatesAForeignFile342r1`／`TestAC2WindowsLeafIdentityReCheck…342r1`）＝`treeLandMinus` 同两枚在场 ⇒ ★`342-r2` 那枚"3 枚红与它字节无关"的**逐名归属复现、枚数是摆动量**（本腿 1⇄2，另在 minus 树到 2）＝按读数纪律裁"方向与形状复现、那个 3 ⛔ 当通则"；★落地侧新增红＝**0 枚**；基线那枚红＝本腿叶子件的有牙红（尺⛔ 恒绿）。唯一合法存量红 `TestC26RewriteAccountIsConsumedAtEverySecurityLeg` 本包红名册⛔ 出现。`--- SKIP`＝0（五轮逐枚）；`mklink /J` 建不出＝0。

## 去牙⇄反形（★三枚落点全⛔ 同腿报；只打仓外导出树物理字节、⛔ `-overlay`、⛔ `checkout`）

发 A 摘 `unlinkPinnedLeaf` 那一步 ⇒ 件一 damage `3→0`＋`entered` 照 3、红句 `AC#5's half-made file survived the refusal`×3；发 B 换按名开 ⇒ `entered 3→0`；发 C 只空 re-check 的 **link 面**（`pinnedleaf_342_windows.go:190`）⇒ link `45~56→0`、SameFile 接住 37、`truncated=38` ⛔ 归零。落刀前后锚命中逐发验（前 1／后旧 0＋新 1）；`go vet` 三发各 rc=0；★三发还原逐字回 pristine `B9DC39AEA15670E8F24689FA8D5F8157F2A1FFE310996203670C37DCF4D5643B`（`MATCH=True`×3）。

## 门禁（只在共享工作树、跑前 `unset GOTMPDIR TMPDIR TEMP TMP`）＝`logs/22`

`gofumpt -l`（写面 0 枚 `.go`；★具名⛔ 把指 `.md` 的 rc=2 读成红＝工具误用）／`sh scripts/d22scan.sh` **rc=0** clean／`sh scripts/check-path-length-budget.sh` **rc=0** VERDICT GREEN（over-budget 60＝roster 60＝⛔ 新增）／`go vet ./internal/winsec/` 本机 **rc=0** ＋ `GOOS=linux`・`GOOS=darwin` 各 **rc=0**（交叉档⛔ 与本机成对判红绿）。★⛔ `go test ./cmd/wisp` 整包、⛔ `./...`（派单令＝射程外、⛔ 读成红⛔ 读成绿）。

## 越界与"声称⇄实际"复核（两把锚并集＋独立整腿形）＝`logs/20`/`logs/21`/`logs/25`

锚定行首名册 **10 笔** ⇄ `--grep='^票 342＋票 348 落地腿 342-r2'` **10 笔** ⇄ 并集（按 hash）**10**；★本派单"续体"形锚定后 **0 笔**（非锚定子串 3 笔全是编排者自己的笔＝`A1045` 那味本腿复现并⛔ 采）；★**⛔ 拿区间当名册**＝`3abbb639..bfcac9ed` 实含 **37 笔**（`296-r2`/`296-v2`/`351-a1`/编排者的笔都在里头）⇒ 改逐枚 name-only 并集＝**43 枚路径**／非 `_test.go` 的 `.go`＝**0 枚**／写面之外＝**0 枚**／`.go` 只有那两枚仪器件；`winsec.go` 10 笔逐枚 `07300bff`（`%w` seam ⛔ 被碰）。它三回自抓本腿逐枚复量＝`probes/342/r2` 现量 **39 枚／0 字节 0／`.out` 0**、`report.md` HEAD blob **119 行** ⇄ 它自己那行 `REPORT_MD_LINES=119` **对得上**（笔⑩真入库）＝★本腿⛔ 再抓到一枚未归因读数当凭据。

## 号对账（一处本腿⛔ 穿派单的号、当场具名）

派单/`342-v2`/`348-v1` 三处引 `:118-140` 自陈段 sha16＝`e8a2f33dea8cdeab`。本腿第一把尺用 `sha1sum` 得 `88f51af530a75de9`＝**相冲**。逐档试（`logs/23`/`logs/24`）＝**同一枚 `sed -n '118,140p'` 区间、`sha256sum | cut -c1-16` 才复现 `e8a2f33dea8cdeab`** ⇒ 号⛔ 错、错的是本腿第一把尺的算法；实质（HEAD ⇄ 基线该区间**逐字同**、1646⇄1646 字节、10 笔逐枚同 sha）两把都对得上＝实现腿一字未动。具名登记＝"同一枚事实两把尺"的形、⛔ 拿任何一枚未归因读数当凭据（这条对我自己也执行）。

## ⛔ 做／⛔ 宣称的（防过度读数）

⛔ 动任何 `.go` 字节；⛔ 翻框（框终量＝票 342 `[x]`5/`[ ]`2、票 348 `[x]`6/`[ ]`1＝与起手逐字同）；⛔ 改票名／`Status:`／`-done`；⛔ push；⛔ 宣布任何一格闭环；⛔ 把件一/件四读成 342 或 348 的 AC 达成；★⛔ 在判语里写"建议机主选 X"（第二道批准那道本腿给独立结论＋理由＝⛔ 该批，摆选单归编排者）。
