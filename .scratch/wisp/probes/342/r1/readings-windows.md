# 342-r1 Windows 半边成对读数（本机＝出货平台，`go1.27.1 windows/amd64`，三棵 `git archive fde1df02` 仓外树）

★本文件只抄日志已有读数。逐字全文＝仓外 `D:\tmp\wisp\20261019\342r1\logs\60-winpair-{treePre,treeLand}-r{1,2,3}.log`＋归因探针 `70-probe-{treePre,treeLand}.log`。竞态读数⛔ 要求逐字相同＝复现的是量级与形状。

## 一句话结论

**POSIX 半边＝关掉并量到；Windows 半边＝落码入库、★两平台"各一发"的口径⛔ 满足、且量到一枚残余。** ⇒ 本腿⛔ 更新 `winsec.go:118-140` 那段自陈为"窗口已关"（票 342 `AC#5`／票 348 同条逐字＝两平台都量到才许写），Windows 残差照实进注释与本页。

## 表一 · `junction` 形（票 342 截断支钉，绝对替换名 ⇒ act-time 面）

| 侧 | 发 | attempts | cycles | provedReparse | formBLegRefusals | truncated | deleted | returnedNil | 名册主体 |
|---|---|---|---|---|---|---|---|---|---|
| pre | r1/r2/r3 | 400 | 11/12/12 | 12/13/13 | **0/0/0** | 0/0/0 | 0/0/0 | 21/21/– | 356＝C26 resolver 稳态拒 |
| land | r1/r2/r3 | 400 | 12/12/12 | 13/13/13 | **7/5/3** | 0/0/0 | 0/0/0 | 12/12/12 | 373/373/372＝C26 resolver 稳态拒 |

⇒ 落地侧★子测试 `…/junction` **3/3 PASS**（损伤 0、本落地那支答 7/5/3、`provedReparse` 同发在场）；pre 侧红的是**"落地那支⛔ 存在"**（形变上没有别的解释＝尺有牙，但红句形状＝instrument⛔ damage，本腿⛔ 把它读成"pre 侧零损伤＝没窗口"：pre 侧 400 发里只有 12 个周期真处于翻转态、且其中多数落进"祖先被移走的 gap"，见下表二 pre 侧 `truncated=7/10/6`＝同一个形的损伤在窗口够宽时量得到）。

## 表二 · `symlink-dir` 形（相对目标、留在根内 ⇒ 只能由身份门答）

| 侧 | 发 | attempts | cycles | provedReparse | formBLegRefusals | **truncated** | deleted | returnedNil |
|---|---|---|---|---|---|---|---|---|
| pre（342 钉） | r1/r2/r3 | 400 | 74/58/57 | 74/58/58 | 0/0/0 | **7/10/6** | 0 | 104/112/– |
| land（342 钉） | r1/r2/r3 | 400 | 74/–/71 | 74/–/71 | 27/–/26 | **1/–/1** | 0 | 77/–/74 |
| pre（348 钉） | r1/r2/r3 | 400 | 52/59/56 | 52/60/56 | 0 | 6/9/6 | 0 | 109/107/107 |
| land（348 钉） | r1/r2/r3 | 400 | 72/–/72 | 72/–/72 | 17/–/17 | 2/–/2 | **0** | 78/–/78 |

⇒ ★两个读数必须分开讲：
1. `deleted`＝票 348 那一支，**落地侧 3/3 全 0**（pre 侧也 0＝这一形上 unlink 那一发⛔ 到达过外来名字；本腿⛔ 用这行声称"348 在 Windows 已量到修前分子>0"＝该维⛔ 满足、下一程要先把这一支的窗口打宽再量）。
2. `truncated`＝票 342 那一支，落地侧**⛔ 归零**：`1/–/1`（342 钉）与 `2/–/2`（348 钉同形）。红句逐字（land r1）：`AC#1 RED: 1 of 400 attempts found the foreign witness's body replaced, so the open inside privateFile resolved through an ancestor swapped into a link and its truncate landed in somebody else's tree (first: answered by the landing's own leg: winsec: path is not provably resolved, refusing to seal: C:\Users\swq\AppData\Local\Temp\…`
⇒ ★这条残余的机制本腿尚未归因到底（两种候选，下一程辨）：(a) 采样与复位本身在这形上有串扰（同一发里 `truncated` 与 `firstProof` 取自同一次读、而 `firstProof` 被先写入的"落地那支答过"占住，导致损伤行与拒答行同源）；(b) 真残余＝祖先在两态之间翻的那一刻，`pinnedParent` 已经拿到 `data/sub` 的句柄、随后 `sub` 被换成根内相对链接，那一发 open 仍在**已钉住的真目录对象**里操作＝按设计⛔ 伤外来对象，若伤的是外来，只能是**另一发**在翻转态里由 `preOpenLeafGuard`＋by-name 之外的路径漏进去。本腿⛔ 在归因完成前把这行写成"已关"，也⛔ 写成"我的字节把包打红"＝`treeBase` 全绿对照与逐名作差在同批 logs 里（见下）。

## 表三 · `leaf-symlink` 形（祖先照旧真、只翻叶子名＝`os.Root` ⛔ 钉叶子的正面）

| 侧 | 发 | attempts | cycles | provedReparse | guard（322 那支、开之前） | landing 复验（开之后） | truncated | deleted |
|---|---|---|---|---|---|---|---|---|
| pre | r1/r2/r3 | 400 | 51/62/60 | 51/62/60 | 0/2/1 | 0/0/0 | 0 | 0 |
| land | r1/r2/r3 | 400 | 86/80/88 | 86/81/88 | 2/1/1 | **0/0/0** | 0 | 0 |

⇒ ★本枚件 3/3 PASS 但 PASS 的是"⛔ 删掉外来名字"＋"有两支叶子腿之一起过作用"；`AC#2' leaf-leg split` 那一行逐字报＝`preOpenLeafGuard (check before the open)=1/1/2, openPinnedLeaf re-check (after the open)=0, provedReparse=86/81/88`。⇒ **本落地的叶子复验在 400×3 发里一发没答**＝它今天⛔ 被证明在路上；稳态里总是 `preOpenLeafGuard` 先答（票 322 那一支），只有落在"判语与开之间"那枚两发系统调用宽的 gap 里才轮到复验，而这枚 gap 在这台机器上⛔ 被这枚 racer 打到。⇒ ★这就是派单点名的"⛔ 把新写一条链当成完事"那味：链我写了、复验我也写了，但**复验的有效性尚未被量到**，⛔ 声称已钉住；下一程的钥匙＝把它做成确定性件（用本包已声明的 seam 在判语与开之间插一次翻转，形同 `inIdentityWindow120r3`），⛔ 靠 racer 撞运气。

## 归因（把 300/400 那桶"some other failure"归到底＝第 14 笔）

探针逐字（两侧同一条串，件 1931 字节，⛔ 入库）：
`winsec: refusing to seal C:\Users\…\data\sub\leaf-342r1w.bin: risk: path traverses a reparse point (junction/symlink) not covered by reparse_point_exceptions`
⇒ 那是**已装好的 C26 resolver 自己的稳态拒绝**（`resolveString` → `internal/risk`），⛔ 我的落地形、⛔ 损伤红；我的名册原先漏了这一条措辞，所以两侧各有 180–253 发落进"other"。补进名册后（标⛔ 计入 `formBLegRefusals`＝⛔ 拿别人的拒答冒充自己那支的 liveness），名册主体归位、残余才看得见＝★这条是本腿交件里"没归因⛔ 读数"的正面例子。

## 台面自抓（Windows 侧新增两枚，⛔ 计入判据）

- racer 自身错误里出现 `The process cannot access the file because it is being used by another process`／`staged-leaf-flip` 的 rename 失败＝Windows 的文件占用语义，`flipErrs` 最高 91（pre 侧 60/71）＝已计数、⛔ 藏，liveness 只数"证过的翻转周期"。
- `junction` 形在 Windows 上⛔ 能靠"移动真目录"翻（`mklink /J` 造的链接被 rename 走之后再改名回来时目标已变），所以本形只保留"翻转后 `reparseAt` 亲证"这一把尺，gap 计数照打。

## 尚未做（⛔ 把没做的说成做了）

★去牙反控两发（票 342 `AC#3`／票 348 `AC#3`，成对落点⛔ 相同）；`winsec.go:118-140` 终态自陈（等 Windows 那两行归因完成）；消费者面成对复认（★`internal/risk` 此刻同波 `349-v1` 在读、本腿⛔ 与它同刻跑同一包）；本机全包成对四数（正在跑，日志名＝本文件表头同目录 `80-host-full-*`，若未到手本腿⛔ 报其数）。
