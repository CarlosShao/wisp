# 342-r1 Windows 半边成对读数（第 16 笔修正后重打＝本版取代旧版）

★本机＝出货平台 `go1.27.1 windows/amd64`；三棵 `git archive fde1df02` 单差异仓外树（`treePre`＝基线产码＋我这六枚新钉／`treeLand`＝我的落地全量／`treeBase`＝锚本身当基线色）。逐字日志＝仓外 `D:\tmp\wisp\20261019\342r1\logs\60-winpair-{treePre,treeLand}-r{1,2,3}.log`（★本轮＝第 16 笔 `7abffe98` 之后重打的那批）。竞态读数⛔ 要求逐字相同＝复现的是量级与形状；attempts＝400/发。

## 一句话终态

**POSIX 半边＝关掉并量到（成对 3/3 红⇄3/3 绿、liveness 同发、去牙两发照红）；Windows 半边＝落码入库、仪器有牙（基线侧每一形都量到损伤），但落地侧⛔ 归零、且 leaf 那形这轮撞上仪器自己的空转 ⇒ ★两平台"各一发"⛔ 满足 ⇒ 本腿⛔ 更新 `winsec.go:118-140` 的自陈、该段至今一字未动。**

## 表一 · 票 342 截断支钉（`TestAC1Windows…NeverTruncates…342r1`，两形各子测试）

| 形 | 侧 | 发 | provedReparse | formBLegRefusals（本落地那支答） | truncated | deleted | 子测试 |
|---|---|---|---|---|---|---|---|
| junction（绝对替换名） | pre | r1/r2/r3 | 13/12/12 | 0/0/0 | **2/1/3** | 0 | FAIL×3＝有牙 |
| junction | land | r1/r2/r3 | 13/–/13 | –/–/4 | **–/–/2** | 0 | r1 FAIL・r2 FAIL・r3 PASS |
| symlink-dir（相对目标） | pre | r1/r2/r3 | 59/–/– | 0/0/0 | **12/–/–** | 0 | FAIL×3＝有牙 |
| symlink-dir | land | r1/r2/r3 | 66/–/67 | –/–/– | 0/–/2 | 0 | r1 PASS・r2 FAIL・r3 FAIL |

⇒ 基线侧★每形每发都量到截断（`truncated=1～12`/400）＝这枚 Windows 钉⛔ 恒绿、它有牙；⇒ 落地侧★⛔ 稳定归零（junction 与 symlink-dir 都还有 1～2/400 的漏出）。★这枚残余本腿⛔ 归因到底＝续点第 1 条（候选见下）。

## 表二 · 票 348 删除支钉（`TestAC2Windows…NeverDeletes…342r1`，`symlink-dir` 形）

| 侧 | 发 | provedReparse | formBLegRefusals | truncated | **deleted** | 子测试 |
|---|---|---|---|---|---|---|
| pre | r1/r2/r3 | 61/–/– | 0/0/0 | 7/–/– | **0/0/0** | FAIL（instrument 红：本落地那支答 0） |
| land | r1/r2/r3 | –/68/72 | –/16/19 | –/1/0 | **0/0/0** | r1 FAIL・r2 PASS・r3 PASS |

⇒ ★本票那一支（`deleted`）在 Windows 两侧都⛔ 取到"修前 >0"的分子＝票 348 `AC#2` 的"两平台各一发"**⛔ 满足**，原因是窗口形状：本平台稳态里总是 C26 resolver／`preOpenLeafGuard` 先答、路由⛔ 到达拒绝路径的删除那一跳（名册逐行：pre 侧 180–245/400 落 `the installed C26 resolver's own refusal`）。⇒ ★本腿⛔ 把这条读成"Windows 上删除那一支无关"＝它读成"这一维今天这台机器⛔ 可达"，并把构造更宽窗口的做法列为续点。

## 表三 · 叶子复验钉（`TestAC2WindowsLeafIdentityReCheck…342r1`，`leaf-symlink` 形）

| 侧 | 发 | provedReparse | guard（开之前，322 那支） | ★landing re-check（开之后） | truncated | deleted | flipErrs |
|---|---|---|---|---|---|---|---|
| pre | r1/r2/r3 | 76/64/68 | 2/1/0 | 0/0/0 | 14/6/7 | 0 | 18/24/28 |
| land | r1/r2/r3 | 9/14/12 | 0/0/0 | 0/0/0 | 0 | 0 | **25389/23848/24691** |

⇒ ★这轮暴露两件事，本腿都⛔ 圆场：
1. **pre 侧有牙**（`truncated=6～14`/400＝祖先真、只翻叶子名时，旧产码那发 open 跟着相对链接把外来文件截断了＝342-a1 `W-ROOTSYM followed=true` 的同一形状，本腿在自己仪器上复现）。★第 16 笔那枚修正（链接目标改指外来**文件**而⛔ 目录）是让这一形真正驱动其 gate 的关键；旧行（`provedReparse=86/81/88`、`guard=1/1/2`、`truncated=0`）⛔ 有效、已具名作废。
2. **land 侧的 racer 空转＝一枚仪器自己的台损**（`flipErrs` 2.4 万⇄`provedReparse` 只 9～14、`returnedNil` 293～329＝周期⛔ 在翻、racer 在错误里打转）。机制候选＝落地后路由会把叶子名删掉（拒绝路径 unlink 生效）⇒ `stageFlip` 在⛔ 存在的 `sub` 里建链接失败⇒`toFlip` 连拒；这使这三行**既⛔ 能当"复验在路上"的证据、也⛔ 能当"复验没答"的证据**＝作废、列为续点第 2 条（做法＝把这一支改成确定性件、借 `inIdentityWindow120r3` 那一族 seam 在"判语与开之间"插一次翻转，⛔ 靠 racer 撞运气）。
⇒ ★因此"本落地的叶子复验到底答没答过"今天＝**未知**，⛔ 已证明在路上、也⛔ 证明它没在路上。`AC#5` 那句"两平台都量到才许写关"因此成立⛔。

## 归因（第 14 笔的结论照旧有效，逐字）

探针两侧同一条串＝`winsec: refusing to seal …: risk: path traverses a reparse point (junction/symlink) not covered by reparse_point_exceptions`；独立正控＝那句话就写在源码里：`internal/risk/pathresolver.go:36 var ErrReparseDenied = errors.New("risk: path traverses a reparse point (junction/symlink) not covered by reparse_point_exceptions")`。⇒ 那⛔ 我的落地形＝已装好的 C26 resolver 的稳态拒绝；补名册时★⛔ 计入 `formBLegRefusals`（⛔ 拿别人的拒答冒充自己那支的 liveness）。

## 本机全包成对四数（★带 `-v`）

- `treeBase`（本轮重跑＝待补，若⛔ 到手本腿⛔ 报其数；上一批＝`RUN=110 PASS_TOP=61 FAIL_TOP=0 SKIP=0` 两发逐字同号）
- `treeLand`＝上一批 v1 `116/63/2/0`・v2 `116/64/1/0`，红名册逐名＝恰两枚 Windows 新钉（基线树⛔ 有＝⛔ 存量红、也⛔ 台面红）。★第 16 笔之后这数会因表一/表三变化，本腿在交件里标为"待重打"。

## 续点（下一枚腿照这三行走）

1. 归因表一落地侧那 1～2/400 的 `truncated`：把"损伤证据"与"拒答证据"分栏各存一份逐字（现在两枚都写进同一枚 `firstProof`、会互相顶），再各连打 3 发判形；候选机制＝(a) 采样/复位串扰；(b) `preOpenLeafGuard` 判语与 `Root` 下降之间那枚窗口（祖先翻在判语之后）。
2. 表三那枚空转：`toFlip` 的失败计数要能区分"名字被路由删掉了"与"链接建不起来"；并把叶子复验改成**确定性**件。
3. 表二那一支要在 Windows 上做出可到达的拒绝路径（⛔ 靠名册里"resolver 先答"充数）。
4. ★以上三条都齐了才动 `winsec.go:118-140`；本腿⛔ 动过它。
