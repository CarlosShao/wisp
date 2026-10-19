# 348-v1（裁决腿 342-v2）本票专属读数 · 删除支（`deleted`）逐发点名

★原始逐字 `-v` 日志的**主存**在 `.scratch/wisp/probes/342/v2/logs/`（两票同批＝同一枚实现腿的交件，本腿只建一份台件）；本件＝本票分子的逐发摘录＋本票专属的两发（确定性件 ⇄ Windows 稳态探针）。
台面＝自建仓外导出树 `D:\tmp\wisp\20261019\342v2\{treeB,treeL}`（基线 `fde1df02` blob `67783ae0` ⇄ 落地 `9581a48c` blob `07300bff`）；容器 `golang:1.27`＋`CGO_ENABLED=0`＋`GOPROXY=off`＋module 缓存只读＋`TMPDIR`/`GOTMPDIR` 留容器内 `/tmp`；Windows＝本机 `go1.27.1 windows/amd64`。四数尺＝`RUN/PASS_TOP/FAIL_TOP/SKIP`。竞态读数⛔ 要求逐字相同。

## 一 · POSIX 竞态侧 `TestAC2POSIXPrivateFileRaceThroughASwappedAncestorNeverDeletesAForeignFile342r1`（attempts=1500）

| 侧 | 发 | cycles | provedLinkCycles | pinnedWalkRefusals | pinnedParentBeforeOpen | **deleted** | truncated | narrowed | untouched | 色 |
|---|---|---|---|---|---|---|---|---|---|---|
| 基线 | 1 | 144 | 144 | 32 | 0 | **86** | 0 | 0 | 1414 | FAIL |
| 基线 | 2 | 131 | 131 | 25 | 0 | **71** | 0 | 0 | 1429 | FAIL |
| 基线 | 3 | 135 | 135 | 22 | 0 | **73** | 0 | 0 | 1427 | FAIL |
| 基线 | 4 | 138 | 139 | 20 | 0 | **73** | 0 | 0 | 1427 | FAIL |
| 落地 | 1 | 159 | 159 | 28 | 22 | **0** | 0 | 0 | 1500 | PASS |
| 落地 | 2 | 152 | 152 | 32 | 26 | **0** | 0 | 0 | 1500 | PASS |
| 落地 | 3 | 146 | 147 | 43 | 25 | **0** | 0 | 0 | 1500 | PASS |
| 落地 | 4 | 156 | 156 | 26 | 15 | **0** | 0 | 0 | 1500 | PASS |

⇒ 本票分子＝`deleted`：修前 **65~86/1500**（腿报 67/68/73~85＝量级与形状都复现）⇄ 修后 **0** ×4 **且同发 liveness 非零**（三把钥匙逐发点名＝`provedLinkCycles`／`pinnedWalkRefusals`／`pinnedParentBeforeOpen`）。`narrowed=0` 两侧＝票 120 那一支在此形没被重新打开。

## 二 · POSIX 确定性侧 `TestAC2POSIXRefusalPathWithASwappedAncestorNeverDeletesAForeignFile342r1`（本腿重打 4 发＝把〔仅腿报〕升成本腿收件复跑）

- 基线 **4/4 红**，红句逐字（shot1）＝
  `AC#2 RED: the foreign witness is GONE after the refusal-path unlink (lstat …/001/data/someone-elses-tree/leaf-342r1.bin: no such file or directory), so the removal re-resolved the caller's spelling through the swapped ancestor and deleted a name in somebody else's tree - ticket 348's shape, the second damage of the pair`
- 落地 **4/4 绿**，同发拒绝逐字＝
  `AC#2' refusal verbatim (the leg that answered): winsec: path is not provably resolved, refusing to seal: …/data/sub/leaf-342r1.bin: pinned object at …/data/sub/leaf-342r1.bin is not what the name carries: this descriptor is not the object the name carries now, so the bytes would not land in what was just sealed (winsec: cannot apply a private security descriptor)`
  绿侧 survivor 逐字＝`AC#2' survivor read-back: size 75 (was 75), mode -rw-rw-rw- (was -rw-rw-rw-), mtime moved=false`
- ⇒ 这一发**⛔ 竞态**＝确定性可达（342-a1 的 A3 形状在本腿台面复现）＝本票 `AC#0` 的那一味。

## 三 · Windows 删除支 `TestAC2Windows…NeverDeletes…342r1`（symlink-dir 形，attempts=400）

| 侧 | 发 | provedReparse | formBLegRefusals | truncated | **deleted** | 色 |
|---|---|---|---|---|---|---|
| 基线 | 1/2/3 | 51/–/53 | 0/0/0 | 2/–/10 | **0/0/0** | FAIL（instrument 红） |
| 落地 | 1/2/3 | 69/67/66 | 19/13/16 | 2/1/2 | **0/0/0** | PASS×1・FAIL×2 |

⇒ **"修前分子>0"在 Windows 取⛔ 到＝缺世界**（⛔ 仪器没做对、⛔ 真漏）。决定性一发＝本腿仓外只读探针件（⛔ 入库，`logs/ii02-mut-W-tree{B,L}-probe.raw.txt`；★同一枚打进两棵树、稳态摆好后再打产码入口＝⛔ 靠 racer）：

| 探针 | 基线树 | 落地树 |
|---|---|---|
| W1a 祖先＝相对目标目录符号链接 | `err=winsec: refusing to seal …: risk: path traverses a reparse point (junction/symlink) not covered by reparse_point_exceptions`、见证 body 原样 | **逐字同**、见证原样 |
| W1b 祖先＝junction（绝对替换名） | 上一字同串、见证原样 | 逐字同、见证原样 |
| W2 叶子＝裸兄弟目标链接（复验唯一能答的形） | 同串（resolver 抢先） | 同串 ⇒ 落地那支叶子复验**没轮到答** |

⇒ 本平台稳态里这句话由**已装好的 C26 resolver** 先答（两侧同串＝⛔ 落地形造成）⇒ 路由到⛔ 达拒绝路径那两枚删除 ⇒ `deleted` 恒 0。腿报名的那桶 `some other failure`（86~90/400）逐字＝同一句 risk 拒绝＝腿第 14 笔的归因有效，本腿独立复现。
★续体（⛔ 结案的唯一通道）＝Windows 侧需要一枚"祖先只在 open 之后翻"的确定性窗口件（POSIX 用本包自己的 seam `inIdentityWindow120r3` 做到了；Windows 今天⛔ 对应物）。

## 四 · 去牙发 1＝摘**本票**那支（本腿独立打、落点⛔ 同腿报）

`sh/10-teeth-A-unlink-byname.sh`＝`unix.Unlinkat(int(p.fd()), base, 0)`→`syscall.Unlink(按名)`；pristine `189e55e8613249db`／旧锚 1→0／新锚 1／mut `ec949b0d7f9f5f0b`／还原 `189e55e8613249db`（`PRISTINE_MATCH_AFTER=YES`）。
读数（`logs/iii-A-unlink-byname.raw.txt`）＝`RUN=10 PASS_TOP=1 FAIL_TOP=3 SKIP=0`：本票两枚钉**逐枚回落红**（确定性 `the foreign witness is GONE…`＋竞态 `AC#2 RED: 52 of 1500 attempts found the foreign witness GONE`），票 342 那枚也红（`deleted=47`）＝⛔ 借绿；覆盖写那枚照绿＝尺指向新守卫。
⇒ 本票 `AC#3`"这一支有自己那枚牙"＝成立。
