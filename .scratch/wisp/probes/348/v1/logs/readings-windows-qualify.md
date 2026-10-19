# 342-v2 读数 2/6 · Windows 三件定性（必答 ii）

台面＝本机出货平台 `go1.27.1 windows/amd64`；树＝与 POSIX 半边**同一对**自建的仓外导出树（基线 `treeB`＝`fde1df02`／落地 `treeL`＝`9581a48c`，两棵都补同串六枚新钉）。尺＝`go test ./internal/winsec/ -v -count=1 -run '<三枚 Windows 钉>'`，attempts=400/形；★仓外树跑 `.exe` 的 `PATH` 用 `cygpath -w`＋`;`（公共块 A1042④ 定式）。逐字＝`logs/ii01-tree{B,L}-win-shot{1..3}.raw.txt`。

## 表一 · 票 342 截断支（`TestAC1Windows…NeverTruncates…342r1`，两形子测试）

| 形 | 侧 | 发 | provedReparse | 落地自己那支答 formBLegRefusals | **truncated** | deleted | returnedNil | flipErrs |
|---|---|---|---|---|---|---|---|---|
| junction | 基线 | 1/2/3 | 12/11/– | **0/0/0** | 0/**1**/– | 0 | 18/20/– | 5/5/– |
| junction | 落地 | 1/2/3 | 12/12/12 | **5/7/4** | **0/0/0** | 0 | 12/12/11 | 5/4/5 |
| symlink-dir | 基线 | 1/2/3 | 53/55/53 | **0/0/0** | **4/7/10** | 0 | 101/104/108 | 23/13/13 |
| symlink-dir | 落地 | 1/2/3 | 66/68/68 | **17/15/15** | **2/2/0** | 0 | 71/77/69 | 5/8/4 |
| symlink-dir(348 那枚件的同形) | 落地 | 1/2/3 | 69/67/66 | 19/13/16 | **2/1/2** | 0 | 72/68/66 | 2/–/– |

⇒ **件 1 ＝ 真漏**（不是台损、不是仪器形状）。三条独立凭据：
- 基线侧每形每发都量到截断（`truncated=1~10`/400）＝这枚钉⛔ 恒绿、有牙；
- 落地侧仍有 0~2/400 的漏出（symlink-dir 形 3 发里 2/2/0 与 2/1/2＝量级稳定、形状复现），**且漏出的字节签名＝产码的载荷**（逐字 `size 44 of 44, content "raced private payload 342r1 windows ancestor"`）＝只有 `PrivateFile` 那只手写这 44 字节；仪器自己的 restock 写的是 74 字节的 witness body、`toReal` 写的是 45 字节的 own body ⇒ **⛔ 可能是采样/复位串扰**；
- 我另打一枚"拔掉产码那只手"的空转对照（`sh/14`，落刀＝`return PrivateFile(...)`→`return nil`，racer/采样/复位全照跑）：junction 形两发都 **`truncated=0` / `returnedNil=400`** ⇒ 残余⛔ 是仪器自造的。（★该对照的 symlink-dir 子测试在拔掉手之后死在仪器自己的复位上＝`settle the world back to the real state: remove …\data\sub: The directory is not empty`＝**台损**、⛔ 计入判据，只取 junction 那一行。）
⇒ 归因⛔ 做到底（落地链每一跳的窗口次序还缺一枚确定性件），但**分类已定**＝残余是真漏的窗口、量级 1~2/400、⛔ 稳态可达。

## 表二 · 票 348 删除支（`TestAC2Windows…NeverDeletes…342r1`，symlink-dir 形）

| 侧 | 发 | provedReparse | formBLegRefusals | truncated | **deleted** |
|---|---|---|---|---|---|
| 基线 | 1/2/3 | 51/–/53 | **0/0/0** | 2/–/10 | **0/0/0** |
| 落地 | 1/2/3 | 69/67/66 | 19/13/16 | 2/1/2 | **0/0/0** |

⇒ **件 2 ＝ 缺世界**（⛔ 仪器没做对、⛔ 真漏）。决定性一发＝我那枚仓外只读探针件（`sh/13`，同一枚打进两棵树、⛔ 靠 racer 撞运气）：

| 探针（稳态摆好后再打一发产码入口） | 基线树 | 落地树 |
|---|---|---|
| W1a 祖先＝相对目标目录符号链接 | `err=winsec: refusing to seal …: risk: path traverses a reparse point (junction/symlink) not covered by reparse_point_exceptions`、见证 body 原样 | **逐字同**、见证原样 |
| W1b 祖先＝junction（绝对替换名） | 同上一字、见证原样 | 逐字同、见证原样 |

⇒ 两平台的删除那一支要跑，必须先让路由**到达**拒绝路径；本平台稳态里这句话由**已装好的 C26 resolver** 先答（⛔ 产码落地与否都一样＝基线侧同串）＝路由⛔ 到达 `unlink` 那一跳 ⇒ 修前分子恒 0。⇒ 票 348 `AC#2` 的"两平台各一发"在 Windows ＝**缺世界**（需要一枚"祖先只在 open 之后翻"的确定性窗口件；POSIX 侧用本包自己的 seam `inIdentityWindow120r3` 做到了，Windows 侧今天⛔ 对应物）。腿报名的那桶 `some other failure`（300/400 上下）逐字＝**同一句 risk 拒绝**＝腿第 14 笔的归因有效，我这发独立复现。

## 表三 · 叶子复验钉（`TestAC2WindowsLeafIdentityReCheck…342r1`，leaf-symlink 形）

| 侧 | 发 | provedReparse | preOpenLeafGuard 答 | ★落地那支 re-check 答 | truncated | deleted | flipErrs |
|---|---|---|---|---|---|---|---|
| 基线 | 1/2/3 | 64/63/71 | 0/1/1 | **0/0/0** | **9/10/10** | 0 | 24/–/29 |
| 落地 | 1/2/3 | **2/13/8** | 0/0/0 | **0/0/0** | 0/1/2 | 0 | **27070/24364/25693** |

⇒ **件 3 ＝ 落地那支今天量⛔ 到（六发逐名 0 命中），且"没答"⛔ 等于"不可能答"＝仍不可判**。
- 基线侧有牙：`truncated=9~10`/400＝旧产码那发 open 跟着相对叶子链接把外来文件截断了（＝342-a1 `W-ROOTSYM followed=true` 的同一形状，本腿独立复现）。
- 落地侧 racer 空转复现（`flipErrs` 2.4~2.7 万 ⇄ `provedReparse` 只 2~13、`returnedNil` 336~344）＝**仪器台损**，与腿表三第 2 条同因；所以"复验答没答过"这发仍答⛔ 了。
- 我另打两发把它从"未知"推到"机制清楚"：
  - 仓外 `os.Root` 谓词探针 v2/v4（卷根逐段下降＝落地链自己的形）：**任何**带 reparse 位的组件（相对目标符号链接⇄junction 一视同仁）在 `OpenRoot` 那一跳就被拒＝`openat sub: path escapes from parent` ⇒ 落地祖先链在这一形上答的是 **containment 面**、`sameObjectPinnedWindows` 那两半门今天这台⛔ 到达（仪器名册第 0 行"act-time refusal of a reparse target"＝文本对、机制名⛔ 对）；而**叶子**若是"目标留在同一目录内的裸兄弟名"这一形，`Root.OpenFile` **照跟随**、写完 `B3 IN-PINNED-DIR sibling now="PROBE-BYTES-B"`＝复验⛔ 是装饰。
  - 确定性探针 W2（同一枚打进落地树，叶子＝裸兄弟目标链接）：产码入口答的仍是 `risk: path traverses a reparse point`＝**resolver 抢先**、复验没轮到。
⇒ 综合＝票 342 `AC#5` 押的"两平台都量到"在 Windows **未量到**；`winsec.go:118-140` 那段自陈"⛔ 出现窗口已关"＝**照实**（本格按票面第二形成立）。腿⛔ 更新那段的决定＝**对**。

## 摆动名册（本腿 Windows 六发里的处置）

`…342r1` 三枚是**本裁在重打的新钉**（⛔ 摆动件）；表三落地侧那三发的 `flipErrs=2.4 万` 已具名为仪器台损、⛔ 当任何一维的凭据；表一"拔手"对照的 symlink-dir 死亡已具名为台损、⛔ 计入判据。票 346 两枚与 `…116` 族在本腿成对四数里逐名 PASS（见 `读数 3/6`）。
