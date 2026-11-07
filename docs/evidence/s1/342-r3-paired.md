# 342-r3 落地腿成对读数（票 342 AC#1 两平台＋变异有牙＋门禁；实现者腿，非裁决）

- 腿：`342-r3`（票 342 落地写腿）。接续＝`4a6e70afe`（前一棒 342r2 撞流超时前落盘的 POSIX 仪器件，运行时从未跑过）。
- 起手锚：台树（`git archive` 导出）切自 `766bf71f4`；`internal/winsec` 在 `766bf71f4..HEAD` 零字节（本腿起手现量）。
- 台件：仓外导出树 4 棵（`D:\tmp\wisp\20261107\342r3\`，只建不删）：`treeLand`（落地）⇄ `treeMutA`（回退票 400 那枚 hunk：`pinnedParent(path, resolved)`→`pinnedParent(resolved, resolved)`，即走回 resolver 答案拼写）⇄ `treeMutB`（完整回到 pre-342-r1 按名路由：按名 `os.OpenFile`＋按名 `os.Remove(resolved)`）⇄ 首轮误形树（见"自抓"）。⛔ `-overlay`、⛔ 复用别腿的树、共享工作树零产码字节（唯一 .go 改动见下）。
- POSIX 档：docker `golang:1.27`、`CGO_ENABLED=0`、`GOPROXY=off`、宿主 module 缓存只读挂载、`GOCACHE=/gocache`（容器本地，`TMPDIR`/`GOTMPDIR` 不绑 Windows 卷＝`342-v1` 定式）。Windows 档：本机 go1.27 windows/amd64。
- 本腿唯一 .go 改动＝`internal/winsec/ancestoropen_342_unix_test.go`（见"仪器修正"），winsec.go 一字未动（HEAD blob `cbafe712` ⇄ 台树逐枚同）。

## AC#1 修后读数（两平台各一发，liveness 同发在场）

| 档 | 仪（同枚） | 分子 | liveness | 判 |
|---|---|---|---|---|
| Windows 落地 | `TestAC1WindowsPrivateFileThroughASwappedAncestorNeverTruncatesAForeignFile342r1`（junction＋symlink-dir 两形 ×400 发） | `truncated=0 deleted=0`（两形） | junction `provedReparse=14 formBLegRefusals=5`；symlink-dir `provedReparse=91 formBLegRefusals=40` | 绿 |
| POSIX 落地 | `TestAC1POSIXPrivateFileThroughASwappedAncestorNeverTruncatesAForeignFile342r1`（bare-sibling 祖先链形 ×1500 发，容器首跑＝补 4a6e70afe 欠的运行时） | `truncated=0 deleted=0 untouched=1500`，`bytes=landing` | `provedLinkCycles=270 pinnedParentBeforeOpen=36`（另 `refusedAfterOpen=37` 分栏） | 绿 |

POSIX 落地全包四数（同字节复跑，尺＝`grep -c '^=== RUN' / '^--- PASS' / '^--- FAIL' / '^--- SKIP'` 于 `-v` 日志 `posix-land-v3.log`）：`RUN=76 PASS_TOP=48 FAIL_TOP=0 SKIP=0`，`ok` rc=0。Windows 共享树全包 `-count=1`：`ok` 18.7s rc=0。

## 变异有牙（摘掉祖先那一支的修法 ⇒ 红；还原 ⇒ 绿；两次读数都记）

| 树 | Windows（同仪） | POSIX（同仪） |
|---|---|---|
| `treeMutA`（回退 400 hunk→走 resolved 拼写） | **红**：symlink-dir `truncated=1/400`＝`AC#1 RED`，red 句逐字入库（junction 形该发 0＝竞态摆动量） | 绿（如实：POSIX `resolved==path`，该 hunk 在本平台语义不存在，绿非凭据） |
| `treeMutB`（pre-r1 按名路由） | **红**：symlink-dir `truncated=6/400`＋`RED instrument`（`formBLegRefusals=0`＝落地腿不存在）＋junction 侧同；**漏出 body 逐字＝产码 44 字节载荷 `"raced private payload 342r1 windows ancestor"`**（342-v2 那三条凭据的形状复现） | **红**：`deleted=92/1500`（首跑）⇄ `deleted=78/1500`（修正判语后复跑）＝按名 `os.Remove(resolved)` 穿翻转祖先删掉外来见证（first proof 逐字 `witness ABSENT … no such file or directory`），`pinnedParentBeforeOpen=0`＝归属干净 |
| `treeLand`（还原） | 绿（上表） | 绿（上表） |

POSIX 修前分子两发（92、78）落回 342-v1 VPROD/VBASE 的量级带（deleted 65..86/1500）＝〔仅腿报〕升成本腿收件复跑。

## 仪器修正（本腿唯一 .go 字节；断言零放宽）

`4a6e70afe` 给本文件加的归属支**首跑即暴露两处不实**（该支此前从未运行）：① 触发条件与文件头声明相反（头写"old 侧 `pinnedParentBeforeOpen>0` 判不可归属"，代码在 `==0` 时触发）；② 触发句"this bench cannot re-take it"已被本腿实测证伪（old 侧 `deleted=92/1500` 就是重新取到的修前分子）。本腿把头部段与触发句改写为实测口径（seal 侧 walk 与父侧 walk 共用 `winsec_other.go` 措辞、分类器按"是否点名 caller 的叶子"分栏；old 侧实测 `pinnedParentBeforeOpen=0` 且 seal 侧 open 后仍答 44 发＝归属成立）。**门行为不变**：old 侧照旧红（上表两发复跑仍红）、落地侧照旧只认落地字节；`gofumpt`/`go vet`（本机＋darwin＋linux）逐枚复跑 rc=0。

## 自抓（⛔ 计入判据）

首轮变异树把两处 `removeRefusedLeaf` 一并中性化（零值 `leafPin` 在 POSIX 会 panic）⇒ 两平台**全绿**＝牙齿被我自己摘掉（POSIX 修前分子走的就是按名 Remove）。改用忠实 pre-r1 形（按名 `os.Remove(resolved)`）后照红。教训入册：**变异必须忠实到"连损伤通道一起还原"**，只摘守卫不摘损伤通道的变异是假绿。

## 门禁读数表（共享工作树，跑前 `unset GOTMPDIR TMPDIR TEMP TMP`）

| 门 | 读数 |
|---|---|
| `gofumpt -l internal/winsec`（真身 `/d/work/base/gopath/bin/gofumpt.exe`） | 0 命中，rc=0 |
| `go vet ./internal/winsec` | rc=0 |
| `GOOS=darwin go vet ./internal/winsec` | rc=0 |
| `GOOS=linux go vet ./internal/winsec` | rc=0 |
| `sh scripts/d22scan.sh` | rc=0 clean（分母 internal/=374、cmd/=58，ban#8 internal/=932、cmd/=224） |
| `sh scripts/check-path-length-budget.sh` | **rc=1 VERDICT RED＝先于本腿存在**：roster 两条过期行（182/185 票件被 rename-cleanup 加 `-done` 改名，`8b32733d8` 等，roster `scripts/check-path-length-budget.sh:345` 一带未跟）＝"roster entry that is not an over-budget tracked path any more"；与本腿零关联（本腿新增 tracked path 仅短名证据件；fix 归 roster 主人，scripts/** ⛔ 本腿写面） |
| 消费者面五包成对＋`cmd/wisp` 那一发 | ⛔ 本波未跑（派单令 ⛔ 碰 cmd/wisp＝另一腿在写；射程外如实记，⛔ 读红也⛔ 读绿，同 342-v3 口径） |

## 没做到的（如实）

① AC#1／AC#4 两格**一枚没翻**（翻框归非实现者）；Windows 裸兄弟叶子形那次 open 的残差（342-v3 件四那味）属票 322/353 的叶子腿射程，本腿未动、未声称关闭，`winsec.go` 自陈段一字未动；② Windows 变异矩阵里 M-A 的 junction 形该发 0 是竞态摆动量（symlink-dir 形红＝有牙凭据）；③ POSIX M-A 绿是构造性 no-op，⛔ 当凭据；④ 长度门 RED 成因已定位但未修（写面外）。
