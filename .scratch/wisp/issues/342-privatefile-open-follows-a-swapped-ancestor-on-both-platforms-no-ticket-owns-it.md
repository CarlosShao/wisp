# 342 — `privateFile` 那次 open 的**祖先那一支**今天两支平台都开着（POSIX 只护最后一枚组件；Windows 连 `O_NOFOLLOW` 都没有）＝票 322 按边界留下、本仓⛔ 主人的一支

- Status: **已立，未派**（2026-10-18 13:5x 编排者立；来源＝票 322 的 `AC#1` 判语与产码注释**逐字自陈**＋我自己现跑；账＝台账 `A1037`）
- 类型：**安全必修邻接**（`D33` 那一族的"路径/链接竞态"）；⛔ 不改契约（`C26` 的入口语义不动）
- ★**这不是"322 没做完"**：322 的 `AC#1` 判语**按"只关叶子那一支"收窄并已勾**，票面逐字写着"只有**叶子**那一支关着、**祖先**那一支仍开着（POSIX 与 Windows 都是），产码注释与本件⛔ 出现任何"窗口已关""。⇒ 322 **零未满足判据、已 `-done`**；本票＝**接那一支的新活**（它此前在本仓⛔ 主人）。
- ⚠ 行号只当**读数**、⛔ 当判据；每条附可重跑的尺。

## ⛔ 别和这几张重复（我开票前逐张读过）

| 票 | 它管的那一支 | 为什么不是本票 |
|---|---|---|
| **103**（`rejected-needs-fix`，4 格未勾） | `SetPathResolver` 谁都能装＋`RemoveUnlinked` 沿 junction 删别人真文件 | 尺＝`grep -icE "privateFile\|O_TRUNC\|openat\|per-component\|逐组件"` 在 103 面＝**0** |
| **108**（`rejected-needs-fix`，5 格未勾） | **103 那枚守卫被三枚探针绕过**（含"全 `/` 或混合分隔符让祖先检查失效"） | 同一族但射程＝**密封守卫的绕过**，⛔ `privateFile` 那次 open 的截断窗口；同尺＝**0** |
| **113**／**116**／**119**（已 `-done`） | 放置校验的 POSIX 链接腿／祖先可行动腿／数据根链接 | 都是 `PlatformVerifyPlacement` 那一族，⛔ open 之后 |
| **120**（已 `-done`） | 密封本身是 check-then-act | 它交回的两支之一＝`120-r2` 的 fd 闸门＝票 322；另一支＝**本票** |
| **137**（已 `-done`） | 12 枚 POSIX 拒绝腿在"宿主自己的 `/var/link`"形上常绿 | 仪器射程，⛔ 那次 open |

⇒ **结论＝"祖先那一支"在本仓今天没有主人**（尺同上，逐张跑）。

## 现量（编排者 2026-10-18 13:5x 自己跑，引用前先重跑）

| 事实 | 读数 | 我这把尺 |
|---|---|---|
| 那次 open 的现场 | `internal/winsec/winsec.go:149` 逐字 `f, err := os.OpenFile(resolved, flags\|oNoFollowPlatform, perm)` | `sed -n '146,150p'` 现读（`:146`＝`preOpenLeafGuard(resolved)`） |
| ★**产码注释自陈两支状态**（同一文件 `:118-138`） | 逐字：`Ticket 322 is the second half of that sentence, and it is only half closed.`／`the LEAF leg is closed on POSIX, by O_NOFOLLOW carried on this open`／**`the ANCESTOR leg is still open on POSIX. O_NOFOLLOW is a property of the last component only, so an ancestor renamed into a link after resolveString returned is still followed by this open. Closing it needs a per-component open-at-time (openat/O_NOFOLLOW down the whole chain, or the handle pinning this package does for the seal), which is not what was adjudicated for this ticket. This window is not claimed closed anywhere.`**／`on Windows O_NOFOLLOW does not exist (oNoFollowPlatform is 0)`／`the ancestor leg is open there too` | 现读（⛔ 抄票面，抄的是**产码注释**＝盘上真件） |
| 敞着的调用者 | `winsec.go:87` 一带＝`os.O_WRONLY\|os.O_CREATE\|os.O_TRUNC`（覆盖写那一支）；`:97`＝`O_EXCL`（票 322 的 `AC#0` 已实测 `O_CREAT\|O_EXCL` 遇链接 EEXIST／不跟随） | 票 322 `322-a1` 实测＋`322-r1b` 交件；★起手复跑 |
| 生产调用者枚数 | `winsec.PrivateFile(` 生产 **5 枚**（`internal/agent/spill.go:166`／`internal/config/migrate.go:93`／`internal/secret/migrate.go:169`／`:189`／`internal/secret/store.go:79`）＋`PrivateFileExclusive(` **2 枚**（`spill.go:343`／`internal/memory/artifacts.go:75`） | 台账 `A1021` 前后编排者复跑过；起手逐枚重跑 |
| ★**Windows 那一侧更难** | 票 322 的 `322-a1` 实测：`FILE_FLAG_OPEN_REPARSE_POINT` 只对**叶子**有效；祖先那一支在 Windows 上连"内核在同一发里回答"这个对应物都⛔ 现成 ⇒ 落点必然要逐组件下探（`CreateFile` 带 `OPEN_REPARSE_POINT` 逐级，或复用本包为密封做的 handle pinning） | `.scratch/wisp/probes/322/a1/{index.md,logs.md}`＋`docs/evidence/s1/322-openleaf.md` |
| 后果形状 | 祖先目录在 `resolveString` 返回之后被换成链接 ⇒ 那次 `O_TRUNC` **跟随它把树外那枚真文件清空**；fd 闸门（`120-r2`）拦得住"私密字节落到外来对象"、⛔ 拦不住**已经发生的截断** | 票 322 票面「为什么立案」第 2 条＋`120-r1` §2 |

## 判据（逐枚可机读）

- [ ] **AC#0 三形代价表（⛔ 未出表任何一形都不许开工）**：**甲**＝POSIX 侧逐组件 `openat(O_NOFOLLOW)` 下探到父目录、再在钉住的 fd 上开叶子（＝产码注释点名的第一形）；**乙**＝复用本包**为密封已经做过的那套 handle pinning**（`120-r2` 的 `sealHandle`/`applyDescriptorPOSIX` 同一族）把父目录钉住 ⇒ 祖先换手后 fd 与 name 不再是同一对象即被复验抓住；**丙**＝**只诚实报告**（`f.Stat()` ⇄ 逐组件 `Lstat` 复验，截断已发生后报错）＝⛔ 防住形，若选它则票面与注释一律⛔ 出现"窗口已关"。★每形现量：改动面（文件∶行）、POSIX 与 **Windows 各自**的实测读数（⛔ 用 POSIX 绿冒充 Windows 绿＝票 322 `AC#3` 那味）、会顶红谁现有的钉（逐名读断言：`placement_symlink_113_other_test.go`／`ancestor_separator_108_other_test.go`／`dataroot_symlink_119_other_test.go`／`seam_probe_root_125_other_test.go`／`posix_seal_window_120_other_test.go`／票 322 的 `openleaf_*_322_*_test.go`）。
- [ ] **AC#1 ★成对读数要有 liveness**：容器／真机台件里成对——修前"外来见证文件被截断"分子 > 0，修后 **0 且 liveness 非零**（⛔ 拿"0 命中"当凭据，除非同一发证明窗口真到达过链接态＝票 322 `AC#1` 逐字口径）。★**两平台各一发**（POSIX 与 Windows 的祖先形**不同机制**，⛔ 一发代两发）。
- [ ] **AC#2 覆盖写语义不许被悄悄改掉**：`:87` 那一支的既有承诺（"写私密内容并允许覆盖"）在选形下仍成立、或被**具名改成拒绝**；两种都要点名现有调用者（上面那 5＋2 枚逐枚给"它依赖覆盖写吗"）。★票 322 `322-a1` 已量到"这个承诺**没有任何一枚现有钉钉住**"⇒ 本票要么补那枚钉、要么具名说明为何不补。
- [ ] **AC#3 去牙反控两发落点不同**：⛔ `-overlay`、只打仓外导出树物理字节、每发前后各验一次 pristine 哈希、两发之间⛔ 并发 shell（`R-228-3b-1` 那味）。★成对第二发＝把新加的祖先守卫换成"逢名就拒"（反形）⇒ 仪器**必须照红**（⛔ 一把两种形状都放行的尺）。
- [ ] **AC#4 门禁＋消费者不推红**：`gofumpt.exe` 真身／`sh scripts/d22scan.sh`／长度门／`go vet`（含 `GOOS=darwin`、`GOOS=linux`）逐枚落 `rc=`；★**跑共享树四道门前先 `unset GOTMPDIR TMPDIR`**（票 311 AC#1 shape B GUARD，带 env 时 rc=2 射程拒绝＝编排者 10-18 坐实，⛔ 读成红）；★`GOOS=linux go vet ./cmd/wisp/` 那一发**在 Windows 宿主上做交叉档＝Go 默认关 cgo ⇒ 必然红**（`sherpa-onnx-go-linux` 全文件被排除）＝⛔ 与本机档成对判红绿、⛔ 当新红（台账 `A1032` 第③节机制已归到底）。消费者面复认：`internal/risk`／`secret`／`config`／`agent`／`memory` 在 Linux 档与本机全 ok 的成对读数＋`cmd/wisp` 那一发（⛔ 抢车道，按波次具名申请）。⛔ `go test ./...`。
- [ ] **AC#5 注释与票面口径说实话**：`winsec.go:118-138` 那段自陈随本票终态更新（★**⛔ 出现"窗口已关"，除非两平台都量到**）；⛔ 许诺未实现的能力。
- [ ] **AC#6 纪律格**：⛔ 翻任何框（本票七格由非实现者裁）、⛔ 改票名／`Status:`／加 `-done`；写面＝`internal/winsec/**`＋`.scratch/wisp/probes/342/r1/**`（短名、⛔ 分节符、逐枚非 0 字节）＋本票末尾**一条**只追加 bullet；`add` 与 `commit` **两处 pathspec 逐枚手写字面路径**（⛔ 从共享 index 反推＝`A1027`／`A1029` 两枚腿踩过）；中文正文⛔ 走 PowerShell（`A1030` 那枚被打坏成逐字 `?` 的主题）；⛔ push。

## 禁区

`docs/PLAN.md`／`docs/specs/**`／`internal/observe/thresholds.go`／golden／`tools/d22scan/allowlist.txt`／`.github/**`／`frontend/**`／`design/**`；三枚冻结件一字不动；⛔ 为变绿放宽票 113／116／119／120／**322** 的既有钉；⛔ 顺手改票 120 的"r 权限类"口径（`A990` 裁②已定）；⛔ 回头改票 322 的判据文字（它已 `-done`，要更正只追加）；⛔ 把叶子那一支当没做（`322-r1b` 已关、`322-v1` 已判）。

**撤销口令**：「撤 342」（⇒ 本票退回未派；票 322 的结案不受影响——它的判据零未满足项）。

## 为什么现在不派

`internal/winsec/**` 本波刚被 `322-r2` 用过（那一枚测试注释面，已交）；★更硬的约束＝**它与票 339 的落地腿都要跑同一批消费者读数**（`internal/risk`⇄`winsec` 交叉），而本包装箱在册的争用假红件（`TestResolvePerCallBudget`）要求安静复量 ⇒ 排在票 339 之后。⚠ 另：本票**甲／乙两形都要动 `internal/winsec` 的产码**＝不是仪器活 ⇒ 落地腿必须是**写腿**、与 `cmd/wisp`⇄`internal/panel` 每波一枚写腿的规矩一起排。
