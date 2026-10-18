# 348 — 拒绝路径那句 `_ = os.Remove(resolved)` 自己**沿祖先链接删掉树外对象**＝票 342 的截断之外**第二枚损伤形状**（甲／乙两形都⛔ 防它）

- Status: **已立，未派**（2026-10-18 23:3x 编排者立；来源＝票 342 只读普查腿 `342-a1` 具名上报（`docs`外的 `probes/342/a1/census.md:20`/`:40`/`:76`）＋编排者自己现读产码坐实机制锚；账＝台账 `A1043`）
- 类型：**安全必修邻接**（`D33` 那一族"路径/链接竞态"的**删除**半边）；⛔ 不改契约（`C26` 入口语义不动）
- ★**为什么单独立票、⛔ 并进 342**：342 管的是"那次 `O_TRUNC` 跟着祖先链接把树外文件**清空**"；本票管的是"密封失败时那句 `os.Remove(resolved)` 沿祖先链接把树外文件**删掉**"。⇒ **两枚不同的系统调用、两种不同的损伤、两个不同的落点行**（342＝`winsec.go:149` 的 `open`；本票＝`winsec.go:160` 的 `Remove`）。★更要紧＝`342-a1` 量到**甲形与乙形都⛔ 防住这一支**（落点后残余 `damage` 全部归它：甲 `damage=66/1500`、乙 `damage=59`、基线 `A2 damage=106`）⇒ 若把它写成 342 的一条 AC，会让"选形落地"必须同时结两件事、把一张已够宽的票再撑宽；★但它⛔ 是独立问题——**任何**装进 `PrivateFile` 的祖先守卫都顺手把 unlink 那一支也变成了"跟着链接删"，所以两者要**同批验、分别有钉**。
- ⚠ 行号只当**读数**、⛔ 当判据；每条附可重跑的尺。

## 现量（编排者 2026-10-18 23:3x 现读；引用前重跑）

| 事实 | 读数 | 我这把尺 |
|---|---|---|
| 那一行的现场 | `internal/winsec/winsec.go` 密封失败支逐字：`if err := sealHandle(f); err != nil {` ⇄ `_ = f.Close()` ⇄ ★**`_ = os.Remove(resolved)`** ⇄ `return sealError(resolved, err)` ⇒ 拿的是**名字**、⛔ 已钉住的 fd | `git show HEAD:internal/winsec/winsec.go` 现读（`:158-162` 一带；我 23:2x 亲自读到该四行） |
| 那句 `_ =` 吞掉了什么 | 删除的**错误被丢弃** ⇒ 删成功删失败都一律往上报 `sealError`，日志⛔ 说得清"我把别的东西删了" | 现读同一支 |
| `O_NOFOLLOW` 救⛔ 到这里 | `winsec_other.go:29` 逐字 `const oNoFollowPlatform = syscall.O_NOFOLLOW`（★自陈只说最后一枚成分）＋`winsec_windows.go:335` 逐字 `const oNoFollowPlatform = 0`；`os.Remove` 走 `unlink` 语义、**那枚 open 旗标跟它无关** | 现读两枚平台件＋`git ls-files internal/winsec` 现量＝`winsec` 只有 `other`＋`windows` **两枚平台件**（⛔ 第三枚细分件 ⇒ `!windows` 一路涵盖 plan9/js） |
| 损伤已被谁量到 | `342-a1` 的 A3 确定性件现量：`unlink err=<nil>` 且见证文件 `no such file`（`probes/342/a1/logs/24-container-A3-unlink-through-ancestor.log`）⇒ **不是竞态、是确定性可达** | ★**编排者⛔ 复跑这一发**（那是一枚非确定性容器竞态台件、复跑要另起一棵树＋一轮容器；我复跑的是**机制锚**＝上面那三行代码）。⇒ 本行读数在 `342-a1` 交件复跑坐实前＝**〔仅腿报〕** |

## 判据（逐枚可机读）

- [ ] **AC#0 起手现跑**：上面四行逐枚重跑（一律取具名提交树）；★并把 A3 那一发**自己重打**（容器台件由本腿自己建、⛔ 复用 `342-a1` 那棵），把"腿报"升成"收件复跑"。与表不符⇒当场具名更正。
- [ ] **AC#1 三形代价表（⛔ 未出表不许动码）**：**甲**＝删改成 **fd-relative**（`unlinkat(fd, base, 0)` 在**已钉住的父目录 fd** 上删叶子名 ⇒ 名与对象同一性由内核在同一发里回答）；**乙**＝**根本不删**（保留失败文件＋把"我留下了什么"打进 `sealError`／日志＝行为变化：密封失败的私密残片会留在盘上 ⇒ 与"私密字节不落外来对象"的承诺相冲）；**丙**＝先 `firstLinkAncestor` 复验再删（⛔ 防住形＝还是 check-then-act、换手窗口照在）。★每形现量：改动面（文件:行）／POSIX 与 **Windows 各自**实测读数（Windows `os.Root` 链＝act-time 拒、但那族今天⛔ 在场；POSIX `unlinkat` 今天在场）／会不会顶红谁现有的钉（逐名读断言：`openleaf_*_322_*`／`sealhandle_*_120_*`／`placement_symlink_113_other_test.go`——`322-a1`/`342-a1` 都点过"拒不得删链接"那一支）。
- [ ] **AC#2 ★liveness 要有、⛔ 拿 0 命中当凭据**：成对读数里"修前见证文件被删 ⇒ 分子 > 0、修后 0 **且同一发证明窗口真到达过链接态**"（＝票 322 `AC#1`／票 342 `AC#1` 同一条口径）；★**两平台各一发**（⛔ 一发代两发）。
- [ ] **AC#3 与票 342 同批、分别有钉**：★这一支与截断那一支**共用同一枚祖先换手竞态**，落地时**两枚都要各自驱动 `PrivateFile` 的常驻仪器**（`342-a1` 已证今天**没有任何一枚在场仪器看得见祖先那一支**：`TestPOSIX120SealWindowLeafSwapIsClosed` 只驱动 `SealFile`、`TestAC1POSIX…322r1` 的 liveness 钥匙是**叶子** ELOOP）。⇒ 判据＝本票那枚新钉⛔ 复用截断那枚的绿；两枚红名册逐名分开。
- [ ] **AC#4 覆盖写／独占语义不许被悄悄改掉**：★`342-a1` H2 已量＝7 枚生产调用者里 **4 枚依赖"允许覆盖写"**（`secret/store.go:79`／`config/migrate.go:93`／`secret/migrate.go:189`／`agent/spill.go:166`）＋2 枚 `PrivateFileExclusive` 要的正是 `fs.ErrExist` ⇒ 乙形（不删）对这四枚意味着"失败的私密残留会留着"，要逐枚点名它接受不接受。
- [ ] **AC#5 颜色与门禁**：`./internal/winsec/` 两平台成对（容器 ⇄ 本机；同一 `-run` 名册、四数逐枚点名⛔ 混用⛔ 相减）＋红名册逐名作差；★在册摆动枚（`TestPOSIX120SealWindowLeafSwapIsClosed`／`TestAC1POSIXOpenThroughASwappedLeafLinkNeverTruncatesAForeignFile322r1`，出处＝票 346）读到红⇒**安静复量 3 发再归因**、⛔ 记成本票新造的红；门禁四数只在共享工作树、★跑前 `unset GOTMPDIR TMPDIR`（票 311 shape B GUARD＝rc=2 射程拒绝）。⛔ `./cmd/wisp` 整包、⛔ `go test ./...`。
- [ ] **AC#6 禁区＋纪律格**：⛔ 动 `preOpenLeafGuard`／`pinPlacement`／`sealHandle` 一族的**语义**（本票落点＝那一枚 `Remove` 的调用方式）；⛔ 为变绿放宽票 113／116／119／120／322 的既有钉；⛔ 把 `winsec.go:118-140` 那段自陈写成"窗口已关"除非两平台都量到；★票 342／票 322 判据文字⛔ 回头改（只追加）。⛔ 翻任何框（本票七格由非实现者裁）、⛔ 改票名／`Status:`／加 `-done`；写面＝`internal/winsec/**`＋`.scratch/wisp/probes/348/r1/**`（短名、⛔ U+00A7 那枚字符、逐枚非 0 字节、⛔ `.out` 扩展名）＋本票末尾**一条**只追加 bullet；`add`／`commit` 两处 pathspec 逐枚手写字面路径（⛔ 从共享 index 反推＝`A1027`／`A1029`／`A1040`）；中文正文先 Write 成文件再 `git commit -F`、★路径必须 Windows 形（`A1042`）；⛔ push。

## ★落地形态（先写清楚，免得两条腿抢同一枚文件）

`342-a1` 已量死：甲／乙**都动 `internal/winsec` 产码＝写腿**，且两形都需要口令级批准（形状＝「批 342 落〈甲|乙〉：允许改 `winsec.go`＋新 landing 件；`winsec_windows.go` 是否同批单独点名」）。⇒ 本票与票 342 **同一枚文件、同一族竞态** ⇒ ★编排者裁：**两票同批交给同一枚写腿**（⛔ 两腿先后各动一次＝第二次要重打第一次的读数、纯浪费；⛔ 两腿同刻动＝互相洗脏）。落地腿的判据＝两票的 `AC#1` 选形各自点名。

## 为什么现在不派

①`internal/winsec/**` 此刻由本波刚落地的写腿 `343-r1`⇄在飞裁腿 `341-v1` 之外的**排程冲突**＋★票 342 的选形（甲／乙／丙）今天还没拍；②落地形态已裁＝与票 342 同批单腿，那一批要等 `342-a1` 的表被非实现者复核过一轮；③宽度帽≤8 与"每波一枚写腿同包"。

**撤销口令**：「撤 348」（⇒ 本票退回未派；那一行代码一字未动）。

## Progress log (append-only, newest last)

- [2026-10-18 23:3x +08] 编排者立票：来路＝`342-a1` 具名上报（它⛔ 自填、把立票权交回＝做对了）；★我现读坐实机制锚（`winsec.go:160` 那句 `_ = os.Remove(resolved)`＋两枚平台件的 `oNoFollowPlatform` 值＋`winsec` 只有两枚平台件）；A3 那一发**我没复跑**＝票面逐行标〔仅腿报〕。名 78 字符≤100、★开票前跑两门（长度门／d22scan）＋点名路径逐枚 `cat-file -e HEAD:<path>` 现验＝**8 枚全 rc=0**（含 `probes/342/a1/census.md` 与那枚 A3 日志＝腿报的原始读数确已在库、可复跑）。
