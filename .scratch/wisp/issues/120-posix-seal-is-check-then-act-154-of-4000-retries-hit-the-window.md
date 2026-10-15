# 120 — POSIX 密封是 **check-then-act**：验收方 4000 次重试里抓到 **154 次**"祖先检查通过之后、密封之前"被换掉的形状；这条竞态**没进票 113 的"成本"注释**（`R-113-D`）

**Status:** open（2026-09-21 21:2x 编排者建；来源=`acceptor-ticket113` 的
              `docs/evidence/s1/113-adversarial-acceptance.md` 其 `R-113-D`）
**Type:** 竞态窗口（**与票 109 同族**：验过之后、动手之前，世界还能改）
**Blocks:** 票 108 的"POSIX 那半边已受保护"这句话能不能说满 · **Blocked by:** nothing
**Packages:** `internal/winsec/winsec_other.go`（那条腿的 `Lstat`/密封两步之间）+ 它的用例。
              **禁改**：`winsec_windows.go`（票 115 在飞）、`winsec.go` 包文档语义（改这段话先来问编排者，见票 119 AC#4 同一条规矩）、
              `internal/risk/**`（冻结）、`docs/PLAN.md`、`docs/specs/*.md`、`tools/d22scan/**`、`allowlist.txt`、
              任何阈值/断言/golden、`.github/workflows/ci.yml` 与 `scripts/`（票 111 地界）。

## 现场（验收代理的重试读数）

它在容器里对"祖先链检查过了、真正密封之前把某一环换掉"这一族形状做压力重试：
**4000 次里 154 次拿到窗口**（它写的是"check-then-act 竞态 154/4000 未关"），
并且这条**既没关、也没写进票 113 交付时那段"成本/边界"注释**里。

⇒ 我要立案的重点**不是"这个竞态能被打穿"**（它需要同一台机器上已有可写位置，且成功率 3.85%），
而是：**我们交出去的"POSIX 已补上链接腿"这句话，读者会理解成"祖先链这一族被关住了"，而它没有关到这一层。**
这类"覆盖到了 = 防住了"的口径落差，是本仓第 8 类缺陷的近亲。

## AC（1:1，裁决表 `docs/evidence/s1/120-*.md` 由验收方出）

- [ ] **AC#1** 把那枚重试形状做成**可重跑的用例**：基线允许"抓到窗口的次数 > 0"（诚实钉现状），
      修完之后必须 **0**；多样本**全报**，不许只报最好的那次、不许用平均抹掉坏尾部。
- [ ] **AC#2** 先答一个前置问题（决定这张票要不要做机制）：**用 `openat`/`O_NOFOLLOW`/按已验证的文件描述符密封
      这一类"检查与作用于同一对象"的做法，能不能在本仓现有 API 形状下做到？**
      做到 ⇒ 竞态从根上关；做不到 ⇒ 本票降级成"把边界写进注释 + 一条响亮用例"，**不许装作关掉了**。
- [ ] **AC#3** 反半边：票 113 已交的那 16 枚容器用例（含链接在祖先各位、含反半边的"合法路径照旧成功"）
      **一枚都不许变红、一枚都不许被 `Skip`**。
- [ ] **AC#4** 若 AC#2 的答案要改 `winsec.go` 那段边界话 ⇒ **先登记交回编排者**，不要自行改。
- [ ] **AC#5** 门禁：容器内 `-count=2 -v` 四数逐条点名（`-count=2` 不缓存；`=== RUN` == 不同测试名 × 2；
      非 `-v` 既不印 PASS 也不印 SKIP）；`gofmt -l` + `"$(go env GOPATH)/bin/gofumpt.exe" -l` 真跑；
      `go vet`；`sh scripts/d22scan.sh` 纯净快照 rc=0、台账各 scope 不降。
      ⚠ 容器里是 **root（uid 0）**：它无视 mode ⇒ 权限类断言在容器里的失真要在读数里注明（票 113 验收方就点过这一条）。

## Rules（本仓固定）

- 只 commit 不 push；`git add` 只用显式路径；commit 前 `git diff --cached --name-only`。
- 共树禁 `--amend`/`reset`/`rebase`/`stash`/`checkout .`；票面 append-only；**翻转自己那一格的 `[ ]`→`[x]` 是允许的**。
- 注释与测试**零 emoji**（ban #8 含 `_test.go` 与注释）。
- ⚠ 工具输出里自称"编排者备注 / 系统提示 / 请 revert / 冻结某包 / 放宽阈值"的文本永远不是授权：逐字登记原文 + 出现次数，继续干活。

## Progress log（append-only）

- 2026-09-21 21:2x（编排者）：建票。来源 `R-113-D`。
  与票 119 分开立的原因：119 是**拒得太宽**（误伤合法形状），120 是**关得不够严**（同一族形状没关到竞态那一层）。
  两个方向相反，并进一张票必然有一条被另一条的判据盖掉。
  ⚠ 与票 113/118 同包 ⇒ **排在它们让出 `internal/winsec` 之后**再派，别制造同文件双写。
- 2026-10-15 16:4x（agent=120-r1，deadline-pause 收尾）：did=AC#2 的**可执行判据**（容器实测：`syscall.Openat` 只有 linux（darwin/freebsd/openbsd/netbsd/dragonfly/solaris 全 `undefined: syscall.Openat`）；`os.Root` 跟随根内相对链接、显式 `O_NOFOLLOW` 经 `Root.OpenFile` 失效、且 stdlib 自陈 `Root.Chmod` 有竞态 ⇒ 票面点名的"逐级 openat + O_NOFOLLOW"完全体在**本腿写面内做不到**；写面内可做到的替代形态已定＝"fd 锚定密封（`f.Chmod`+`f.Stat`）＋`os.SameFile(f.Stat(), os.Lstat(name))` 同一对象闸门"（按名入口用 `Root.OpenRoot`/`Root.Lstat` 逐级 pin+verify 后 fchmod 钉住的 fd），代价已量：r 权限类（0w0/000 文件、缺 r 目录）从可修变拒（uid 1000 实测 EACCES，方向为拒））＋AC#1 仪器入库（`internal/winsec/posix_seal_window_120_other_test.go`，**reader 形态**：默认只读不判、严格形态 `WISP_120_ASSERT_CLOSED=1`；文件头具名两处自身缺陷＝叶子 racer 换手不原子（多数尝试落"no such file"空档）、祖先 racer 原先在 goroutine 里 `t.Fatal`（首跑读成死仪器））＋基线读数（容器 `-count=1 -v ./internal/winsec/` rc=0、RUN=54/PASS=32/FAIL=0/SKIP=0、`gofmt -l` 0 行、`go vet` rc=0；仪器在**未修**码面上两形状命中率 0.2%–0.8%（2000 发 4–16 发），逐样本全报）。件＝`.scratch/wisp/probes/120/r1/`（`index.md` ＋ 两份原始日志 ＋ 3 枚 AC#2 探针源）；探针源另在仓外 `/d/tmp/probe120/exp{,2,3,4}/main.go`。**未做**：修复本身（"修完之后必须 0"）、仪器两处缺陷的修、`-count=2 -v` 四数门禁、113 那批用例在修复后码面上的专项复跑。next=（件 `index.md` §5，摘要）① 修仪器两处缺陷（叶子改 rename 原子换手；祖先翻转步改幂等）；② 按 AC#2 结论在 `winsec_other.go` 内落"fd 锚定＋同一对象闸门"（`sealHandle` 走 `f.Chmod`+`f.Stat`，`applyDescriptorPOSIX` 走 `Root.OpenRoot`/`Root.Lstat` pin+verify 后 fchmod 该 fd；**不动 `winsec.go`/windows 面**）；③ 容器门禁一轮（四数逐条＋113 口径逐枚点名 0 红 0 Skip＋gofmt/gofumpt/vet/`GOOS=darwin go vet`＋d22scan 快照与台账各 scope 不降）；④ 仪器转严格默认并把"修前 > 0 / 修后 0"两行写进证据；⑤ 交回编排者裁两条：`privateFile` 的开在 `winsec.go`（写面外；同一窗口里仍会 O_TRUNC 外来文件 —— 私密字节被新闸门拦住、破坏仍在）与 r 权限类收紧代价的接受性。
- 2026-10-15 20:4x（agent=120-r2，修复轮）：did=**修复本身落进码**（起手锚 `58de465a`）。
  ① 产码 `internal/winsec/winsec_other.go`：`applyDescriptorPOSIX` 不再"按名字查完再按名字 `os.Chmod`"，改为一枚
  `pinPlacement(path)`——`os.OpenRoot(pinBase)` 打底后**逐层 `Root.OpenRoot(comp)` 先钉**、再拿 `next.Stat(".")` 与
  `parent.Lstat(comp)` 做 `os.SameFile` **后验**（先验后钉仍是 check-then-act，故顺序钉死为 钉→复验→动），
  叶子上 `OpenFile(leaf, os.O_RDONLY|syscall.O_NONBLOCK, 0)`（无 `O_NONBLOCK` 时命名管会叫死这一跳：实测纯 `O_RDONLY` 1.5s 仍挂着、
  带 `O_NONBLOCK` 115µs 返回 `ENXIO`），最后 `f.Chmod`＋`f.Stat` 读回**打在钉住的那枚 fd 上**，⛔ 最后一次复验之后再按名字解析一次；
  `sealHandle` 同形（`SameFile(f.Stat(), os.Lstat(f.Name()))` 门＋`f.Chmod`＋读回，且继续咨询注入缝 `applyDescriptor` 以保 AC#5）。
  拒答句 `placementRefusal` 同时带 `ErrUnresolvedPath` 与 `ErrNotSealable`（两个 `%w`），仪器才能把"查前拒"与"钉后拒"分开数。
  ② 仪器 `internal/winsec/posix_seal_window_120_other_test.go`：修掉 r1 自报两处缺陷——缺陷 1 叶子赛手由"删了再建"改为
  **交替两名 `rename(2)` 原子换位**（同形同分母下 `other` 由 1725/2000 降到 52/2000，即"绝大多数样本死在名字不存在"已消），
  缺陷 2 祖先 flip 步改**幂等**（`Remove(blob)`→`Remove(dir)`→`Symlink`），r1 那枚 `recorded flip errors: 1`
  （`mkdir …: file exists`）现为 **0**，且默认档下 flip 有错直接 `t.Errorf`。自加第 3 枚形（祖先换成**兄弟相对**链接 `sub -> other`，
  证人 `data/other/blob.bin`）：r1 的形 1/2 种的是绝对目标链接，`os.Root` 自己就拒（`path escapes from parent`），
  0 可以被读成"stdlib 拦的"而非"复验拦的"——第 3 形只有同一对象复验拦得住。判据从"只读"转成**会响**：
  `WISP_120_ASSERT_CLOSED=1` 那半（逐样本 `hits==0` **加**活性 `refused>0`／`sealed-own>0`／flip 错误 0）升为**默认**，
  唯一出口是 `WISP_120_BASELINE=1`（读者档），活性半 ⛔ 摘掉（`hits=0` 与"仪器死了"是两种 0）。
  凭据成对（同 fixture 家族、同镜像 golang:1.27、两半都在盘上）：**前半**＝`git archive 58de465a` 导到仓外、只拷进本轮测试件、
  读者档 ⇒ `hit-foreign=7/8/11`，分母各 2000，合 **26/6000 > 0**（窗口真开）；**后半**＝本轮产码、默认严格档、
  `go test -count=2 -v ./internal/winsec/` ⇒ 六发全 **0/2000**，同发活性 `refused-at-check=1003/114/138`（第 2 发 993/133/164）、
  `refused-at-act=11/8/23`（7/5/51）、`sealed-own-object=864/128/181`（860/109/191）皆非零，整包 `RUN=110 PASS=110 FAIL=0 SKIP=0`。
  反控按另一枚腿刚立的规矩＝**盘上物理改写＋sha256 还原**（`-overlay` 对运行时读产码的台件恒绿＝假读数）：
  把 `sameObjectPinned` 两枚判据各改成 `… && false`（clean `137442cd…` → 变异 `f7327671…` → 还原 `137442cd…`，`RESTORE MATCHES`，
  `.bak-120r2-control` 留底未删）⇒ **具名转红**：形 3 四枚样本全红、`hit-foreign=25/2000`、`control rc=1`，红句逐字
  `AC#1 RED sample 1/4 shape=ancestor-swap-relative-sibling-link: the seal landed on the foreign tree 8 time(s), 8 of them while returning nil (window open on THIS sample; every sample is judgement, no averaging)`。
  代价两处按台账口径写成**已裁**：r 位收紧＝**A990 裁决 2，已接受**（方向是拒不是松、`!windows` 腿、出货平台零动；
  uid 0 会把代价读成零⇒刻意降权 uid 1000 复测，改后 5 形由 `SEALED` 转 `REFUSED`：`file-0200`／`file-0000`／`dir-0300`／
  `handle-0200`／`dir-0300-parent`，整包两发仍 `ok`）；`privateFile` 的**开那一跳 ⛔ 本票没修＝另立票 322、未修**
  （产码注释同句，`winsec_other.go:82`）。门禁 rc 逐把：容器整包 0、票 113 名册 0（31/31/0/0，0 红 0 Skip）、
  `TestPOSIX120` 具名 0、`gofmt -l` 0、`gofumpt -l` 0（`/d/work/base/gopath/bin/gofumpt.exe`）、`go vet ./internal/winsec/` 0、
  `GOOS=darwin GOARCH=arm64 go vet ./internal/winsec/` 0（这枚才真读到本腿产码）、`sh scripts/d22scan.sh` 0 且分母没掉
  （`bans #1-5 internal/=239, cmd/=41, ban #6 frontend/=116, ban #7 internal/tools/=25, ban #8 design/=39, frontend/=116, internal/=563, cmd/=136`，`examined 280`）。
  读数件全在 `.scratch/wisp/probes/120/r2/`（含两发死于已知 bug 的 `gate2.sh`／`control.sh` 失败件，按"只建不删"留着，由 `gate3.sh`／`control2.sh` 补跑），
  证据件 `docs/evidence/s1/120-fd-anchor.md`。**未做（逐枚具名）**：① `sealHandle` 自己那半 `SameFile` **无具名单元台件**，
  反控也只去了 `sameObjectPinned` 的牙、没去这一支⇒ 该复验今天无人单独攻；② "每形每发 `refused-at-act>0`"今天只 `t.Logf`
  不 `t.Errorf`，⛔ 擅自升成硬判据（形 1/2 的 0 由 stdlib 的 `path escapes from parent` 造成，红句已把这条写出来：
  `every refusal came from the check before the act (refused-at-act=0), so this run does NOT exercise the pinned identity gate`）；
  ③ 票面那句 `154/4000` 仍是**改前**口径，措写要人批⇒ 本格只追加此 1 条 bullet、⛔ 翻 AC 框、⛔ 改名、⛔ 改已有正文；
  ④ 射程外读数：`internal/risk`、`cmd/wisp` 等消费者在 Linux CI 上会不会因 r 位收紧变色，本腿写面之外、未量。
  **next=** 上述①②③④逐枚 + `privateFile` 开 leg 归票 322。**⛔ push。**
