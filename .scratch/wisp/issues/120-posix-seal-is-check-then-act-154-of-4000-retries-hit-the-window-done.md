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

- [x] **AC#1** 把那枚重试形状做成**可重跑的用例**：基线允许"抓到窗口的次数 > 0"（诚实钉现状），
      修完之后必须 **0**；多样本**全报**，不许只报最好的那次、不许用平均抹掉坏尾部。
- [x] **AC#2** 先答一个前置问题（决定这张票要不要做机制）：**用 `openat`/`O_NOFOLLOW`/按已验证的文件描述符密封
      这一类"检查与作用于同一对象"的做法，能不能在本仓现有 API 形状下做到？**
      做到 ⇒ 竞态从根上关；做不到 ⇒ 本票降级成"把边界写进注释 + 一条响亮用例"，**不许装作关掉了**。
- [x] **AC#3** 反半边：票 113 已交的那 16 枚容器用例（含链接在祖先各位、含反半边的"合法路径照旧成功"）
      **一枚都不许变红、一枚都不许被 `Skip`**。
- [x] **AC#4** 若 AC#2 的答案要改 `winsec.go` 那段边界话 ⇒ **先登记交回编排者**，不要自行改。
- [x] **AC#5** 门禁：容器内 `-count=2 -v` 四数逐条点名（`-count=2` 不缓存；`=== RUN` == 不同测试名 × 2；
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
- 2026-10-15 21:1x（agent=120-v1，验收腿＝非实现者，零产码）：did=对 `120-r2` 自报三枚具名缺口给三态 + 独立复跑成对读数。
  件＝`docs/evidence/s1/120-v1.md`。起手锚 `400927c4`、`git status --porcelain -- internal/winsec` 空、docker 29.6.2。
  **成对读数独立复跑**（`/d/tmp/probe120v1/{fixed,prebase}`，先数 `.go`=1073/1072 证非空）：改前 `WISP_120_BASELINE=1`
  三形 hit-foreign **9/12/22 合 43/6000 > 0**（旁站分母 2000 与活性 refused/sealed-own 非零＝坏世界真被造），改后默认严格档
  三形 **0/2000**；与腿 26/6000、编排者缩小档**不可相减**，只作方向独立复认。**三枚缺口三态**：① `sealHandle` inline
  `SameFile`（`winsec_other.go:269-274`）＝**缺世界**——本腿物理去这一支的牙（clean `137442cd…`→变异 `57e29c0b…`→还原
  `137442cd…` `RESTORE MATCHES`、`.bak-120v1-sealhandle-control` 留底、仓内工作树全程未碰）跑整包 ⇒ **RUN=55 PASS=55 FAIL=0
  SKIP=0 rc=0**，零枚变色 ⇒ 无任何用例走这一支；须落地腿补具名台件 + 单独去牙，本腿⛔ 代补。② 形 1/2 去牙后仍 0＝**成立**
  （射程限制、非藏）——Root 对绝对目标链接 `openat …: path escapes from parent` 在 open 处即拒，判据**只买到形 3**
  （兄弟相对链接才逼 `sameObjectPinned` 开火）。③ `refused-at-act>0` 只 `t.Logf`＝**成立**（建议⛔ 整体升硬：本腿一发
  shape2 逐样本 `refused-at-act=1,1,0,0` 已证升硬会 flaky/形 1、2 恒红；闸门有牙已由去牙反控钉住）。
  **注释说实话两处**：`winsec_other.go:52-74` r 位收紧写成现状 + `A990` 裁②具名、无新造承诺＝成立；`:82` 票 322 "unfixed"
  现核 `322-…md` 存在、Status "已立，未派落地腿"⇒ 仍属实＝成立。**名册口径**：`-run 'TestAC[1-4]POSIX'` 现数**顶层 20 枚**（票里 16 已过期）、
  容器 RUN=31/PASS=31/0 FAIL/0 SKIP rc=0，腿报名册数与逐枚点名一致。**门禁 rc**：容器 fixed 0、prebase 0、整包去牙 0、名册 0、
  `gofmt -l` 0、`gofumpt -l` 0、`go vet ./internal/winsec/` 0、`GOOS=darwin GOARCH=arm64 go vet ./internal/winsec/` 0、
  `sh scripts/d22scan.sh` 0（首跑 rc=2 系 GOTMPDIR 落 D 卷触 ticket-311 体积守卫，改指 C 卷后 0；分母 internal/=564 ≥ r2 的 563，各 scope **未降**）。
  **存量债具名记色**（⛔ 算进本票）：`cmd/wisp` 的 255 名册尺起手即红（归票 324，本腿**面限制⛔ 未跑** `go test ./cmd/wisp`，〔待面空补验〕）、
  `internal/panel` 5 枚既有红（未跑、非本票）。**没裁完**：`sealHandle` 补集台件（本腿⛔ 产码）、`internal/risk`/`cmd/wisp` 在 Linux CI 上因
  r 位收紧的射程外变色（面外，〔待面空补验〕）、票面 154/4000 与"16→20"措写（属改票面须人工批准）、`privateFile` 开 leg 归票 322。
  ⛔ 翻 AC 框、⛔ 改票名/Status/已有正文。**next=** 落地腿补 `sealHandle` 具名单元台件 + 单独去牙反控（缺口 1 去处）。**⛔ push。**
- 2026-10-15 22:1x（agent=120-r3，落地腿＝实现者，零产码：`winsec_other.go` 全程 `sha256=137442cd…` 未变）：did=补掉 `120-v1` 判"缺世界"的那枚缺口 1——给 `sealHandle` 自带的两枚 inline 复验（`:269` 链接那一支、`:272` 真对象那一支）**各**造出能走到的世界，并对每一支**单独**去牙做反控。件＝`docs/evidence/s1/120-sealhandle-fixture.md`＋测试三枚 `internal/winsec/sealhandle_window_120_other_test.go`（夹具与判据）／`sealhandle_identity_120_other_test.go`（直接调 `sealHandle` 四枚）／`sealhandle_privateroute_120_other_test.go`（走 `PrivateFile` 产码路由两枚＋一枚盲跑读数），读数与脚本在 `.scratch/wisp/probes/120/r3/`（23 枚＋`tooth.sh`／`gates.sh`）。**新用例 7 枚**（顶层，全 `TestAC1POSIX…120r3`）。**两形都造出来了**：① `:272`＝fd 开在真对象 A、`rename(2)` 把第二枚真对象换到名字上，全程零链接 ⇒ `pinPlacement` 按名钉到的就是它、`sameObjectPinned` 满足，只有 `sealHandle` 自己的 `os.SameFile` 能拒（确定性、⛔ 接缝）；② `:269`＝名字必须在 **pinned seal 判语出来之后**才变成相对兄弟链接 ⇒ 经包自己声明的接缝变量 `applyDescriptor` 先调**真的** `applyDescriptorPOSIX`、保留判语，之后才动手（⛔ 替换；`private_fail_test.go` 为 AC#5 早换过同一枚 var）。★**对派单 §1 那句"⚠ 必须 Root 会跟的形"要更正为读数**：相对兄弟链接若种在调用**之前**，`Root` 会跟它、`pinPlacement` 在 `:188` 当场拒，`sealHandle` 在 `:257` 就返回——**这正是 120-v1 写的成因**，本腿把它具名成一枚用例（`TestAC1POSIXSealHandleNotReachedByAPrePlantedLink120r3`，逐字 `leg that spoke = pinPlacement sameObjectPinned link face (winsec_other.go:188)`）而不是照字面当前提。**反控成对三发**（同一快照 `git archive HEAD=6eb3d2e8`、同镜像 `golang:1.27`、同 temp；物理改写＋sha256 还原、`.bak` 两支留底未删；锚⛔ 用行号——同文件 `sameObjectPinned` 在 `:188`/`:191` 有一对同形判据，改法是先由 `func sealHandle` 体范围现找行＋硬校 `false &&` 全文件恰好 1 处＋`LANDS` 自证）：未去牙 **rc=0／RUN=62／PASS=62／FAIL=0／SKIP=0**（顶层尺 40/40/0/0；交件锚 `e1973191` 再确认一发同数）；去 `:269` 一支 **rc=1／FAIL=2**＝只有本腿两枚链接形红，红句逐字 `sealhandle_identity_120_other_test.go:86: AC#1 RED name-becomes-a-sibling-link-after-the-pinned-seal: refusal came from sealHandle inline SameFile face (winsec_other.go:272), want the leg carrying "the name carries a link, so this descriptor is not the object the caller named"`（★**这一支摘掉后仍被拒、但拒因换腿** ⇒ 判据必须问"哪一腿说的"、⛔ 只问"有没有拒"）；去 `:272` 一支 **rc=1／FAIL=2**＝只有本腿两枚真对象形红，红句逐字 `…:48: AC#1 RED name-carries-another-real-object: the call returned nil, so the leg carrying "this descriptor is not the object the name carries now, so the bytes would not land in what was just sealed" never spoke`＋盘上事实那一行 `…:60: AC#1 RED: the caller's own object was narrowed to -rw------- …`（没牙时 `chmodPinned` 真落在 fd 的对象上）与路由那发 `…:87: … the refused write left the name carrying a file, want it unlinked; err=<nil>`（产码 `PrivateFile` 会**报成功**）。**零变色的尺**⛔ 看四数：`--- FAIL` 连名逐枚比对，三发名册 `diff` 全等、变色只发生在本腿用例（其余 60 枚含 113/118/119/125 那族与我的成因枚、反方向枚、盲跑枚全绿）。**盲跑读数**（⛔ 接缝帮忙、`PrivateFile`×2000 由赛跑者翻名字）同发指纹＝未去牙 `inline SameFile=25 / inline link=11`，去 `:269` 那发 link 行整行消失、去 `:272` 那发 SameFile 行整行消失；命中数按 120-v1 §2 裁③**只上报⛔ 升硬**。**门禁 rc**：`go build ./...` 0（宿主不编 `winsec_other.go`，⛔ 当已验）、`GOOS=linux go vet` 0、`GOOS=darwin GOARCH=arm64 go vet` 0（真读那一半）、`gofmt -l` 0（0 字节）、`/d/work/base/gopath/bin/gofumpt -l` 0（绝对路径）、`sh scripts/d22scan.sh` 0 且分母无一下降（`internal/=568`≥120-v1 的 564、`cmd/=136`、`tools/=25`、`design/=39`、`frontend/=116`、`examined 280`）、`sh scripts/check-path-length-budget.sh` 0 GREEN、容器名册 `-run 'TestAC[1-4]POSIX'` 0（现量**顶层 27 枚**＝旧 20＋本腿 7，RUN=38/PASS=38/0/0；只跑本腿 7 枚 `-run '120r3$'` 0＝7/7/0/0）。**没做完／没证的（逐枚具名）**：① 自然 landing 的分布只有**一发**、⛔ 多次读数与阈值（判者可攻）；② 盲跑里 `other failure=1` 逐字 `openat caller-object-120r3.bin: too many levels of symbolic links`——只上报、⛔ 归因（本腿没查清那一发的交错顺序）；③ **R-322 那一形在这条路由上是活的**：名字带链接时 `privateFile` 的 `O_TRUNC` 在 walk 之前把外站对象截空（读数 `the witness now holds ""`），产码 `:78-82` 自己写着这一句不做、票 322 已立未派 ⇒ 本腿**只计数上报⛔ 判⛔ 顺手修**，因此枚 7 的内容判据是三态（模式硬判／raced 载荷落进外站对象硬判／只被截空＝计数）——⚠ 那条"内容不许变"的硬断言是本腿自己首发起的、首发即红，改的是**本腿用例的射程**，根据＝产码注释＋322 已立事实，判者要攻这一点请攻这里；④ Windows 半格（`winsec_windows.go:328` 缺同一对象闸门）非本单射程；⑤ 台面故障三形已具名实录（`awk` 替换串里 `&`＝匹配文本 ⇒ 造出 `build failed` 且四数全 0、`MSYS_NO_PATHCONV` 泄漏给宿主尺令 d22scan 拿到 `D:\d\work\…missing go.mod` rc=2、`tail -1` 多文件＋`set -e` 令驱动中途死）——前两发的作废读数件**随 `rm -rf` 重建快照一起消失**（只建不删⛔ 删，这一笔是我自己的重建步骤带走的，具名认错）。**⛔ 翻 AC 框、⛔ 改票名/Status/已有正文、⛔ push。**
- 2026-10-16 11:2x（agent=120-v2，终裁腿＝非实现者·只读，零产码、零字节触碰工作树）：did=五格逐格判语＋★`sealHandle` 那半**裁活**。
  件＝`docs/evidence/s1/120-v2-verdict.md`（244 行／27735 字节，现量）＋读数 23 枚与两枚去牙脚本 `.scratch/wisp/probes/120/v2/`（含 `logs/`）。
  起手锚 `232f309e`、时刻 11:14 +08、票框现量 **5 未勾／0 已勾**、`git status --porcelain -- internal/winsec` 全程空。
  台面＝`git archive` 的**仓外导出树**（`head` 742 枚 `.go`／`prebase` 729 枚，先数后跑）；HEAD 那枚 `winsec_other.go` 现量
  `sha256=19f3c211055011cc699a8d8027ef49ff543ede0da5c863f58aaaad21a640d3e2`（456 行）＝**与 120-r3 全件引用的 `137442cd…` 不同**
  （票 322 两枚写腿 `8baec646`/`1fbc9067` 在 120-r3 交件之后改过它）⇒ 本腿所有读数落在**今天交付的字节**上，比腿报那一版更强。
  **★`sealHandle` 判「成立」＝120-v1 的"缺世界"对今天这枚面不再成立**：7 枚台件复跑 `RUN=7 TOPPASS=7 TOPFAIL=0 SKIP=0 rc=0` 且世界行与读数同发
  （`identity:44 AC#1' world: fd and name = two different objects`；夹具 `:41-43` 的 `t.Fatalf` 让"没造出世界"自己叫红）；
  逐支**单独**去牙（物理改写仓外树＋sha256 还原、⛔ `-overlay`）＝摘 `:300`（链接那一支）恰红 **2 枚**（突变 `989af687…`，红句＝仍被拒**但拒因换腿**），
  摘 `:303`（真对象那一支）恰红 **2 枚**（突变 `76e49806…`，含盘上事实那一行 `caller's own object was narrowed to -rw-------`）；
  ★同形两对之问用一发正面区分＝把**另一对** `sameObjectPinned`（`:219`/`:222`）两支一起摘 ⇒ 红的是**完全不同的 4 枚**（票 322 一枚、成因枚 3、盲跑枚 7、仪器形 3），
  而我那 4 枚链接/真对象形全绿；再拆开单摘：只摘 `:219` 只红成因枚、只摘 `:222` 红 322 枚＋形 3 ⇒ 牙下在对的那对上
  （尺＝整文件 `ModeSymlink != 0` 命中 **3 行**＝219/300/455、`!os.SameFile(pinned, now)` 命中 **2 行**＝222/303，现量；改法＝由 `func sealHandle` 体范围 287..307 现找行＋硬校 `false &&` 全文件恰好 1 处）。
  三发全部 `RESTORE MATCHES`（还原＝`19f3c211…`），还原后整包复跑回基线色 `RUN=71 TOPPASS=43 TOPFAIL=0 SKIP=0 rc=0`。
  **五格**：AC#1 **成立**（多样本＝`sealWindowSamples120=4`×`sealWindowAttempts120=500`；成对读数本腿现跑＝改前 `9/8/17` 合 **34/6000>0** 且 `refused-at-act=0/0/0`，
  HEAD **0/6000** 且活性非零 `check=942/82/133 act=19/9/20 sealed-own=943/121/137`；"抓到窗口的次数"是**量的**＝`posix_seal_window_120_other_test.go:171-172`
  每次 attempt 读回外站证人的盘上 mode，单位＝次数；反形三发都会红＝未修码面走默认严格档 `3/3 全红 rc=1`（样本里只命中 1 次也叫红）＋只摘 `:222` 红形 3＋确定性台件红；
  判据在样本循环内、total 行只额外打印 ⇒ 无平均、无"只报最好那次"）；AC#2 **成立**（可执行答复＝"按已验证 fd 密封"做得到且已是今天落地的形状，"逐级 `openat`+`O_NOFOLLOW` 完全体"做不到 ⇒ 只关到密封那一跳，⛔ 读成"从根上关"）：
  `syscall.O_NOFOLLOW` 容器实测 **`0x20000`(131072)**；11 平台×2 枚交叉编译矩阵自跑（`O_NOFOLLOW` 于 windows/js/plan9 取不到；★`syscall.Openat` 除 linux 外 **aix/ppc64 也 COMPILES**＝更正前任"只有 linux"的口径，⛔ 改文字、只具名记现量）；
  `Root.OpenFile(leaf, O_RDONLY|O_NOFOLLOW)` 对相对兄弟链接读数 `OPENED, and the descriptor stands on the TARGET object=true`（＝经 `os.Root` 失效，我量的）；
  `winsec_windows.go:335` 那枚常量今天值＝**0**，OS 手段 `:374 preOpenLeafGuard`（`:379 FILE_FLAG_OPEN_REPARSE_POINT`）只给判据不给不跟随的覆盖开；
  AC#3 **成立**（名册现量＝`-run 'TestAC[1-4]POSIX'` **顶层 28 枚**、四数 `39/28/0/0 rc=0`＝0 红 0 Skip；票面"16"过期、120-v1 20、120-r3 27、第 28 枚＝票 322 落的 leaf-open 枚；⛔ 改票面只具名；
  逐枚点名两栏见件 §5；★一枚具名口径落差：**"同一个对象被换掉"这一形，113 那批是失明的**（只摘 `:222` 时红的正是 120 自己的形 3 与 322 枚，28 枚名册一枚不红）⇒ 该形今天的凭据只在 120 的仪器与 `120-r3` 的台件上）；
  AC#4 **成立**（尺＝`git diff 58de465a..232f309e -- internal/winsec/winsec.go` 现量 **2 个 hunk／+42−1**，被删那一行逐字 `f, err := os.OpenFile(resolved, flags, perm)`；
  `git log` 该枚件在此窗口**只有 1 枚 commit**＝`8baec646`（票 322 遗产、编排者代提）⇒ 票 120 三枚腿一枚没碰那段话，唯一改动者**已登记**＝台账 **A1001** 第 7 条接线尺，
  其四枚落点本腿逐枚 `grep -n` 复认存在：`winsec.go:146`/`:149`、`winsec_other.go:29`、`winsec_windows.go:335`/`:374`）；
  AC#5 **成立**（容器 `-count=2 -v` 整包四数逐条＝`RUN=142 / TOPPASS=86 / TOPFAIL=0 / TOPSKIP=0 rc=0`，口径写清＝不同顶层名 43 枚 ×2＝86 条 `--- PASS`、`=== RUN` 含 28 枚子测；
  `gofmt -l` 列 **0 枚**、`gofumpt -l` 列 **0 字节**（`/d/work/base/gopath/bin/gofumpt.exe` 真跑 rc=0）、`go vet` host/linux/darwin-arm64 三把 rc=0、
  `sh scripts/d22scan.sh` 纯净全仓导出树 rc=0 `clean - no D22 ban violations`，分母**成对同尺**（58de465a ⇄ HEAD）＝`ban#8 internal/=563→573`、`cmd/=136→139`，其余 239/41/116/25/examined 280 持平 ⇒ 无一 scope 下降；
  ⚠ `design/` 我读 **30**、前任报 39＝尺的射程差（他们跑在带未跟踪件的工作树、我跑在 `git archive` 树；`git ls-tree -r design/` 两枚锚都 30 枚），⛔ 当成 scope 掉了；
  容器实测 `getuid=0` ⇒ 权限类断言在本台面失真＝票 113 常注，这一族夹具头自己写明不碰权限面）。
  **能力类必答**：`sealHandle` 产码调用点**只有 1 枚**＝`winsec.go:158`（在 `privateFile` 体内），`PrivateFile`(`:95`) 与 `PrivateFileExclusive`(`:105`) 两条公开入口都经它 ⇒ 生产调用者上**通**；
  测试侧 5 枚件 29 处命中；定义两枚（POSIX `winsec_other.go:287`／Windows `winsec_windows.go:413`）。⛔ 拿 `go build` 的 rc=0 当"接上了"。
  **交付缺陷三枚具名（⛔ 本腿动一字）**：① `spokeLeg120r3` 名册印的行号过期（印 `188/191/269/272`，今天 **219/222/300/303**；判据按文案子串匹配、四串在产码逐字仍在 ⇒ 红绿与归属不受影响，但会误导下一枚腿）；
  ② `sealhandle_privateroute_120_other_test.go:233-238` 与 `sealhandle_window_120_other_test.go:7/:21-23` 仍写"票 322 unfixed at this anchor"，而 322 的**叶子腿今天已落地**
  （本腿两发盲跑那一形都没出现＝`witness holds its own bytes after 2000 raced attempts`），仍开的只有**祖先腿**（`winsec_other.go:103-113` 逐字）；
  ③ 票 120 的 **Windows 半格仍空**＝`winsec_windows.go:413` 缺同一对象闸门（`120-r3` §5④ 指的 `:328` 今天落在 322 的注释段里＝①同一族）。
  **没裁完的格（具名＋缺哪一发读数）**：Windows 侧＝**不可判**〔缺 windows/amd64 上"名字被换成另一枚真对象"这一形的台件与去牙反控；仓里没有这一枚，且票面逐字禁改那枚件〕；
  射程外消费者（`internal/risk`/`cmd/wisp`/`internal/secret` 因 r 位收紧在 Linux CI 的变色）＝**不可判**〔缺 `go test ./cmd/wisp` 整包（票面 314 死亡区、本腿⛔ 跑）与 `internal/tools` 整包（同批只读腿 `270-v1` 在用、本腿⛔ 跑）〕；
  自然 landing 分布只两发（`:272/:269`＝10/9 与 9/1，另 `returned nil=26/50`、`other=41/21`）＝按 120-v1 §2 裁③只上报⛔ 升硬⛔ 给阈值。
  **本腿自抓两把坏尺具名（⛔ 计入判据）**：把读者档与严格档塞进同一枚 `docker run -e WISP_120_BASELINE=1` ⇒ 那发"严格档 rc=0"是假读数（作废件留底
  `logs/v2-prebase-strict-DUPLICATE-of-reader.log`，只建不删）；`grep 'live scope work' | head -1` 抓到 d22scan 仪器自测那枚临时小树的摘要（`internal/=14`）＝坏尺，改 `tail -1` 才是最终行。
  **⛔ 翻任何 AC 框、⛔ 动 `Status:`/票名/`-done`、⛔ push；本腿只 commit。**

- [2026-10-16 11:5x +08] **编排者翻框（五格全翻，凭据＝非实现者终裁腿 `120-v2` ＋我 11:4x–11:5x 逐字重放它的三把尺全部对上）**：AC#1／AC#2／AC#3／AC#4／AC#5 ⇒ 判语一律「成立」（⛔ 附条件）。
  ★**我自己复跑的**（尺＝容器 `golang:1.27`＋挂宿主 mod 缓存＋`GOPROXY=off`，台面＝`git archive HEAD` 的仓外导出树 `/d/tmp/orch-120v2-replay2/head`，⛔ 动工作树；腿的原树 `/d/tmp/probe120v2` 与它的 driver.log 一字未动）：① **牙**：摘 `sealHandle` 链接面（`:300`）⇒ `RUN=71 TOPPASS=41 TOPFAIL=2`，红在 `TestAC1POSIXSealHandleNameBecomesLinkAfterPinnedSeal120r3` ＋ `TestAC1POSIXPrivateFileRefusesWhenNameBecomesLinkAfterPinnedSeal120r3`；摘同一对象面（`:303`）⇒ 也 `71/41/2`，红在那两枚 `…NameCarriesAnotherRealObject120r3`；★**区分发**＝摘**另一对同形**（`sameObjectPinned` 的 `:219`＋`:222`）⇒ `71/39/4` 且红名**完全不同**（322 那枚＋120 的另三枚）⇒ 证明"我打的是哪一对"是读数、不是假设；还原 `sha256=19f3c211…d3e2` 三发 MATCHES、复确认 `71/43/0/0 rc=0`。② **AC#3 那把尺**我原样重放＝`-run 'TestAC[1-4]POSIX' -count=1 -v` ⇒ **`RUN=39 TOPPASS=28 TOPFAIL=0 TOPSKIP=0 rc=0`**（顶层枚数尺＝`grep '^--- PASS' | awk '{print $3}' | sort -u | wc -l`），与腿的 28 逐枚对上。③ **AC#5 那把尺**＝整包 `-count=2 -v` ⇒ **`RUN=142 TOPPASS=86 TOPFAIL=0 TOPSKIP=0 rc=0`**，也与腿同。④ **AC#4** 我自己现跑＝`git diff 58de465a..HEAD -- internal/winsec/winsec.go` **2 个 hunk、+42/−1**，且该区间内 `winsec.go` 只有 **1 枚 commit**＝`8baec646`（票 322 的盘上遗产、编排者代提，已登记台账 `A1001`）⇒ 票 120 的腿一枚都没碰那段边界话。⑤ **能力类**（⛔ 拿 `go build` rc=0 冒充接线）＝`sealHandle` 产码调用点 **1 枚**（`internal/winsec/winsec.go:158`，在 `privateFile` 体内，`PrivateFile` 与 `PrivateFileExclusive` 两条公开入口都经它）＋定义 2 枚（POSIX `winsec_other.go:287`／Windows `winsec_windows.go:413`）＋测试 5 枚件。
  ★**更正我自己的一条旧账**（项目记忆 10-15 那条"sealHandle 那半＝缺世界、两发去牙都全绿"）：**已被今天的字节推翻**——`120-r3` 之后同文件里存在**两对**同形判据（`sealHandle` 的 `:300`/`:303` 与 `sameObjectPinned` 的 `:219`/`:222`），今天摘 `sealHandle` 任一面**恰红 2 枚**（我 11:5x 自己复跑，见上 ①）。旧读数不是假的：那时**那对判据还没被造出来**（缺世界是真的，世界后来被补上了）。教训记我：**"缺世界"是有保质期的结论，引用前必须重放，不许当永久裁决**。
  ⇒ 五格全勾、⛔ 一枚残余留本票 ⇒ 本票**改名 `-done`**（尺＝翻后名册 `unchecked=0 / checked=5`；防重领键＝文件名后缀）。逐枚残余**具名归口**，⛔ 抹平：
  ① **Windows 那半格**（裁决腿判"不可判〔缺 windows 台件与去牙反控〕"）＝**票 322 的射程**——我现读：`322` 的 **AC#3「Windows 那一半单独答」已勾**，只剩 **AC#4**（门禁＋消费者不推红）未勾、它在等一枚 `cmd/wisp` 整包**终态**读数（那格归正派出的 `220-v1` 与我随后复跑）。票 322 自己开头一句就把这跳指回本票："票 120 关掉的是『私密字节落到外来对象』那一跳（`sealHandle`/`applyDescriptorPOSIX` 改到钉住的 fd＋同一对象复验）"。
  ② **枚数过期只具名记现量、⛔ 改票面原文**：票面写"16"、`120-v1` 20、`120-r3` 27、`120-v2` 28、**我也是 28**（尺见上 ②）；改措写要人批。
  ③ 裁决腿具名的一枚**口径落差**：票 113 那 28 枚钉的是"链接站着时的拒绝"，对"换过去又换回另一枚真对象"这一形**失明** ⇒ ⛔ 读成"祖先被换掉这一族都钉住了"（这一形今天的凭据只落在 120 自己的仪器＋`120-r3` 的台件上，而这两枚都已随本票结案）。
  ④ 容器是 root（`getuid=0`）⇒ 权限类断言在本台面失真（票 113 常注，本票⛔ 依赖）。
  ④ **台面故障具名（⛔ 计入判据，记我）**：我第一次把树建在 `mktemp -d`＝`/tmp/tmp.4xaEwDqEu5/head`，Docker 把容器外那层的 `/tmp/…' 读成它自己 VM 的路径、⛔ Windows 那枚 ⇒ 挂载指到空目录、`go test` 连 `go.mod` 都看不见、`rc=1` 且⛔ 落 `$tag.log`。**定式＝凡 rc 非 0 又拿不到日志，先证"日志该出现的那枚目录真被挂上"**；容器可挂的树要建在 `/d/…` 这类 Docker 认得的卷上。

- [2026-10-18 11:4x +08] **编排者追加更正（⛔ 改上文一字；来源＝收尾腿 `322-r2` 具名上报＋我自己现跑复认，台账 `A1027` 第⑥节）**：上面第 2 条「过期"未修"文字」那格里，**`sealhandle_window_120_other_test.go:7/:21-23` 这一处的指认不实**——它**从来没写过** "unfixed at this anchor"、也**从不提票 322**。尺（锚＝`322-r2` 交件侧，我复跑）：`git grep -c 'unfixed at this anchor' -- 该件`＝**0**；`git grep -c '322' -- 该件`＝**0**；`git log -S 'unfixed at this anchor' -- 该件`＝**零枚提交**（＝那句话历史上也没在那儿出现过）。该件 `:7/:21-23` 讲的是 `sealHandle` 那两枚 inline 复验没被走到，与票 322 无关。⇒ **本条更正只收窄指认范围：真过期的只有 `sealhandle_privateroute_120_other_test.go` 那一处**（`:230-238`，含 `winsec_other.go:78-82` 的过期行号指针），而它**已由 `322-r2` 落笔改准**（提交 `ab49e9f9`，7+/7−、净 0 行、CR 0 枚＝LF、只改词、⛔ 动任何断言；`winsec_other.go` 本体一字未动、复量仍 `19f3c211…`＝`322-r1b` 的 pristine 同串）。★另把机制说清免得下一程再猜：票 322 终态＝**叶子那一支已关、祖先那一支仍开**（POSIX 与 Windows 都是），所以那句改口写成带条件事实句，⛔ 任何"窗口已关"的读法。本票（120）的 Windows 半格（上面第 3 条）不受影响、仍空。
