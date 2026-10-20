# 票 348 续体终裁（`342-v3`＝非实现者裁决腿·只读产码·与票 342 同批＝同一枚实现腿 `342-r2` 的 10 笔交件·本票＝删除支）

裁决者＝`342-v3`（⛔ 实现者、⛔ 产码、⛔ 翻框、⛔ push）。起手锚＝本腿第一笔 `git log -1` 现认 `677d48372ea07bda4e77fd70bb4b8cd16dacb443`（10-20 10:08）；实现腿终态＝`bfcac9ed`。前序终裁＝`348-v1-verdict.md`（裁腿 `342-v2`）＝`AC#0`/`AC#1`/`AC#3`/`AC#4`/`AC#5`/`AC#6` 成立、**`AC#2` 部分成立＝Windows 缺世界**。本腿只裁续体那一格＋纪律复核，⛔ 重述前裁（除下面两处具名）。
台件＝自建四棵仓外导出树（基线 `fde1df02`⇄落地 `677d4837`⇄落地⛔ 本腿字节 `3abbb639`⇄去牙专用 `treeT`）＋★三枚本腿自写仓外台架（`probes/342/v3/sh/06`/`07`/`08`＝纯平台原语、⛔ 产码字节、⛔ 模拟）；临时根 `%LOCALAPPDATA%\Temp\v3-342v3` 短形；`go1.27.1 windows/amd64`。本票散件＝`.scratch/wisp/probes/348/v2/**`（分子摘要＝`readings-348.md`）；同批全量逐字日志主存在 `.scratch/wisp/probes/342/v3/logs/`。

## 逐格判定（★每条带"哪把尺＋射程目录＋量的还是算的"）

- **AC#2 ★liveness、⛔ 拿 0 命中当凭据 ＋ 两平台各一发 ＝ 部分成立**（POSIX 成立、Windows＝缺世界，★成因本腿已钉死）〔量的〕。
  POSIX＝成立（前裁在册、本腿无矛盾事实）：`deleted` 基线 `86/71/73/73` ⇄ 落地 `0`×4 且同发 liveness 在场。
  Windows〔量的，本腿两棵树各 3 发＝`probes/342/v3/logs/04..09`〕＝REFUSAL 栏三发×两树**逐字同一条** `risk: path traverses a reparse point (junction/symlink) not covered by reparse_point_exceptions`（`refused=true`/`isUnresolved=false`/`isNotSealable=false`）；DAMAGE 栏 `truncated=0 gone=0`（外来见证 `present (size 72)`、调用者自己的 `present (size 42)`）；★**阳性对照**逐发 `by-name remove err=<nil>`＋`control name is now=absent`×两树＝那枚 0 是读数⛔ 钝尺。
  ⇒ 照实判＝"修前分子>0"在 Windows **仍取⛔ 到**＝`342-r2` 那句本腿复现、⛔ 读成达成；本腿⛔ 圆场、⛔ 硬判⛔ 成立。**本格⛔ 结案。**
  ★续体形状本腿据实改写（⛔ 放宽任何断言）＝**⛔ 再押"补一枚确定性窗口件"**（那枚世界经下面成因 2 已证⛔ 可造＝押它等于押一枚买⛔ 到读数的批准）；押＝① 新票（名字那一跳）落地；② 或机主仍要 Windows 祖先那一格的读数＝须同时接受把删除改成⛔ 按名字走（＝行为变化＝本票 `AC#1` 三形表里的新第四形＝另批另拍）。按 `A1042` 票 339 `AC#2` 同形＝登记续体、⛔ 结案。

- **AC#0 起手现跑 ＝ 成立（维持前裁；本腿只补一枚号对账）**〔算的＋量的〕。本腿现读两棵树逐字确认票面那四行：基线 `winsec.go:158-160`＝`if err := sealHandle(f); err != nil {` ⇄ `_ = f.Close()` ⇄ `_ = os.Remove(resolved)`，★第二枚 `:165`（写失败支）同形＝命中 2；落地命中 **0**（`:175-177`/`:180-182`＝两处 `unlinkPinnedLeaf(parent, filepath.Base(resolved))`）。★A3 那发前裁已升成"收件复跑"＝本腿⛔ 重翻；★号对账见下面成因 1（`sealHandle` 基线⇄落地同住 `winsec_windows.go:413`）。

- **AC#1 三形代价表 ＝ 成立（维持前裁）＋★具名一味：今天多出一枚实测的第四形**〔量的〕。落形＝甲那半（POSIX `unlinkat(父fd, base, 0)`／Windows `Root.Remove(base)`→`Deleteat` 带 `FILE_OPEN_REPARSE_POINT`），父钉 `defer close` 仍摆路由末端＝形变耦合照口径。★新形＝落地的删除按**钉目录里的名字**走＋落地句柄带 `FILE_SHARE_DELETE` ⇒ 名字在 seal 判语后可被腾空＋顶进一枚⛔ 点名的对象＝顶进来的被删（凭据＝成因 3＋去牙 A/B）＝**⛔ 属本票判据族**（本票两枚分子逐发都 0）、★本腿签**单立一票**（名与判据形状见 `342-v3-verdict.md` 末节＋`probes/348/v2/logs/01`）。

- **AC#3 与票 342 同批、分别有钉、⛔ 复用它的绿 ＝ 成立**〔量的，本腿新证一味"⛔ 借"〕。尺＝`grep` 它那两枚 Windows 仪器件＝`pinnedParentBeforeOpen`/`provedLinkCycles` **各 0 命中**（真身住 `ancestoropen_342_unix_test.go`／`ancestorunlink_348_unix_test.go`／`ancestorworld_342_348_unix_test.go`＝`probes/342/v3/logs/22` 末段）；本腿裁 liveness 也只用它自己具名的等价钥匙（件一 `realSealRanFirst=true`×3×两树＋`worlds entered=3/3`⇄基线 `0`；件二阳性对照行×两树×3）。★尺有牙＝同一枚件四打在基线产码上 3/3 红、红句逐字 `AC#1-leaf NOT YET MEASURED even on the answerable form`＝⛔ 恒绿。

- **AC#4 覆盖写／独占语义不许被悄悄改掉 ＝ 成立（维持前裁）**〔算的〕。★本腿新证＝**三枚去牙都⛔ 碰这一维**：摘 `unlinkPinnedLeaf` 那一步／把 `openPinnedLeaf` 换成按名开／只空 re-check 的 link 面 ⇒ 覆盖写那枚常驻钉（`TestAC2OverwritePromise…322r1`）在落地两发全档⛔ 进红名册（`probes/342/v3/logs/10`/`11`）；`:96`/`:106` 旗标行与承诺段⛔ 在本腿射程内被任何一笔改动＝10 笔逐枚 `winsec.go` blob `07300bff`。

- **AC#5 颜色与门禁 ＝ 成立**〔量的〕。成对四数（`./internal/winsec/`、`-v -count=1 -timeout 40m`、★只在仓外树、⛔ 混用⛔ 相减）＝落地 `119/67/1/0` ⇄ `119/66/2/0` ⇄ 基线 `114/64/1/0`×2 逐字同号 ⇄ 落地⛔ 本腿字节 `115/62/2/0`；红名册逐名＝落地两枚全姓 `342r1`（`…NeverTruncatesAForeignFile342r1`／`TestAC2WindowsLeafIdentityReCheck…342r1`）＝`treeLandMinus` 同两枚在场⇒★`342-r2` 那枚"3 枚红与它字节无关"＝**逐名归属复现、枚数为摆动量**（本腿 1⇄2、minus 树到 2；按读数纪律那个 3⛔ 当通则、⛔ 据此判它报假）；★落地侧新增红＝**0 枚**；基线那枚＝本腿叶子件的有牙红。`--- SKIP`=0（五轮逐枚）＝⛔ 靠 SKIP 蒙的绿、`mklink /J` 建不出 0 枚。门禁（共享工作树、`unset GOTMPDIR TMPDIR TEMP TMP`）＝`gofumpt -l`（写面 0 枚 `.go`；★具名指 `.md` 的 rc=2＝工具误用⛔ 读成红）／`d22scan` **rc=0 clean**／长度门 **rc=0 VERDICT GREEN**（over-budget 60＝roster 60）／`go vet ./internal/winsec/` 本机 **rc=0**＋`GOOS=linux`・`GOOS=darwin` 各 **rc=0**（交叉档⛔ 与本机成对判红绿）。⛔ `./cmd/wisp` 整包、⛔ `./...`＝派单令、射程外。唯一合法存量红 `TestC26RewriteAccountIsConsumedAtEverySecurityLeg` 住 `internal/risk`、本包红名册⛔ 出现＝⛔ 被本票动过。

- **AC#6 禁区＋纪律格 ＝ 成立**〔算的，本腿另起两把形〕。名册＝锚定行首 **10 笔** ⇄ `--grep='^票 342＋票 348 落地腿 342-r2'` **10 笔** ⇄ 并集（按 hash）**10**＝两把同号；★本派单"续体"形锚定后 **0 笔**（非锚定 3 笔全是编排者自己的笔＝`A1045` 复现）；★**⛔ 拿区间当名册**＝`3abbb639..bfcac9ed` 实含 **37 笔**（`296-r2`/`296-v2`/`351-a1`/编排者都在里头）⇒ 改逐枚 `git show --name-only` 并集＝**43 枚路径**／非 `_test.go` 的 `.go`＝**0 枚**／写面之外＝**0 枚**。禁区＝`preOpenLeafGuard`/`pinPlacement`/`sealHandle` 语义⛔ 动（10 笔逐枚 `winsec.go`＝`07300bff`＝`%w` seam ⛔ 被碰；`:118-140` 逐枚同 sha）；★本票判据文字⛔ 回头改（只追加）；⛔ 新导出名（两枚仪器件除 `TestXxx` 入口外全非导出）。框终量＝本腿**一枚没翻**：票 348 `[x]`=6／`[ ]`=1（未勾＝`AC#2`）⇄ 票 342 `[x]`=5／`[ ]`=2；尺＝`grep '^- \[x\] \*\*AC'`／`'^- \[ \] \*\*AC'`（射程＝两枚票面，`probes/342/v3/logs/22`）；票名／`Status:`／`-done` 一字未动；两票尾各**一条**只追加 bullet（本腿）。

## ★"缺世界"钉到成因＝四枚凭据（本腿自打；`342-r2` 与前裁都⛔ 有；逐字行见散件）

1. **seam 站错了跳**〔算的，`git cat-file blob` 逐字〕＝`var applyDescriptor`（`winsec_windows.go:38`）**基线 blob `fde1df02` 里同一行同字**＝⛔ 落地才有的 plumbing；`func sealHandle(f *os.File) error { return applyDescriptor(f.Name(), false) }` 基线⇄落地**同住 `:413`** ⇒ 钩子只落在**叶子句柄还活着**那一瞬间，而本票那一支的删除在 `_ = f.Close()` **之后**。句柄活着时祖先动⛔ 了（成因 2），祖先能动的瞬间⛔ seam。
2. **平台锁两态**〔量的，`probes/342/v3/logs/01`，3 发全对〕＝基线态（按名开过、句柄已关、⛔ 钉）`rename(data)`／`rename(sub)` 两枚 `err=<nil>`×3 ⇄ 落地态（父钉句柄站着）`Access is denied`／`The process cannot access the file because it is being used by another process`×3；★正控＝空手两枚都落、**释放钉后又落得动**；★**跨进程**（`logs/03`）同向＝⛔ 同进程假象。
3. **祖先进⛔ 去、名字进得去**〔量的，`logs/02`＋`logs/03`〕＝H1（钉后开前）／H2（开完 seal 前）／H3（seal 后删除前）／H4（提前放钉＝本票 `AC#1` ⛔ 许的那形）**四枚瞬间的祖先 rename 全拒**×3 发（H0 空手＝`err=<nil>`）；★同一次测量里独立进程可 `remove(caller's name)=<nil>`＋`rename(mine -> caller's name)=<nil>`，随后 holder 的产码那一步 `Root.Remove(base)=<nil>`＝**攻击者那枚文件的唯一名字被删**（×2 发）。⇒ 本票 booked 的那一支在 Windows 落地＝结构性进⛔ 去；量到的是另一枚形状＝签新票。
4. **两枚去牙把机制钉死**（★落点⛔ 同腿报；`logs/15`/`16`/`18`）＝摘 `unlinkPinnedLeaf` 里 `_ = p.root.Remove(base)` ⇒ 件一 damage `3→0` 而 `entered` 照 3（＝损伤确由产码那一步、⛔ 仪器自造）；`openPinnedLeaf` 换成按名 `os.OpenFile` ⇒ `entered 3→0`（＝窗口确由句柄带 `FILE_SHARE_DELETE` 打开）。落刀前后锚命中逐发验（1→0／新锚 1）；`go vet` 各 rc=0；★还原逐字回 pristine `B9DC39AEA15670E8F24689FA8D5F8157F2A1FFE310996203670C37DCF4D5643B`（`MATCH=True`×2，⛔ `checkout`）。

## 一句话终态

票 348 那一支＝**POSIX 半边关掉并量到**（前裁在册）；**Windows 半边＝`342-r2` 把"缺世界"从形容词升成带阳性对照的分栏读数（本腿两树各 3 发复现＝做到的一味），但"修前分子>0"照旧取⛔ 到**＝**`AC#2` 部分成立、⛔ 结案**；★本腿把"缺"钉成因＝**seam 站错了跳＋祖先那一支经 H1..H4 实测为结构性进⛔ 去**（⇒ 第二道批准签⛔ 该批），同次测量确定性量到的名字那一跳＝**签新立一票**（⛔ 属本票族）。⇒ `AC#0`/`AC#1`/`AC#3`/`AC#4`/`AC#5`/`AC#6` 成立、`AC#2` **部分成立**＝**本票⛔ 结案**。

## ★本腿更正自笔（追加、⛔ 抹原句；尺＝`probes/342/v3/logs/10|11|14` 的 `^--- FAIL` 逐行现读）

`AC#5` 那格里"＝`treeLandMinus` 同两枚在场"这句**按逐名读不准确**，本腿提交后现跑对回抓到：落地 r2 的两枚名＝`{…NeverTruncatesAForeignFile342r1, TestAC2WindowsLeafIdentityReCheckRefusesAForeignNameAndNeverDeletesIt342r1}` ⇄ minus 那发（本腿只跑过**一发**）＝`{…NeverTruncatesAForeignFile342r1, TestAC2WindowsPrivateFileThroughASwappedAncestorNeverDeletesAForeignFile342r1}`＝**共用第一枚、第二枚名不同**。⇒ 成立的只到**族级**＝三发红⛔ 例外全姓 `342r1`、两棵树红名册里⛔ 出现任何 `342r2` 名＝**"这 3 枚红与它字节无关"本腿复现到族、⛔ 复现到逐名同号**＝本票 `AC#5` 的判定⛔ 改（落地侧新增红仍＝0 枚）；另⛔ 把"`342-r2` 报的 3 枚"当本腿复现过的数＝本腿落地两发＝1⇄2 枚红、⛔ 到过 3（竞态⛔ 要求逐字相同这条同样约束本腿，⛔ 拿它轮的数当自己量过）。同把尺顺手坐实基线侧那枚红＝`{TestAC1WindowsLeafRecheckAnswersOnTheBareSiblingForm342r2}` ×2 发逐字同＝有牙红⛔ 存量。
