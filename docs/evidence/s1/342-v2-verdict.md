# 票 342 终裁判语（`342-v2`＝非实现者裁决腿·只读产码·与票 348 同批）

裁决者＝`342-v2`（⛔ 实现者、⛔ 产码、⛔ 翻框、⛔ 选形）。实现腿＝`342-r1`（17 笔，终态 `9581a48c`）。起手锚＝本腿第一笔 `git log -1` 现认 `c3c0ac80`。
台面（两把我自己建、⛔ 复用别腿的树）＝POSIX＝容器 `golang:1.27`＋`CGO_ENABLED=0`＋`GOPROXY=off`＋宿主 module 缓存只读＋`TMPDIR`/`GOTMPDIR` 留容器内 `/tmp`（⛔ 绑 Windows 卷）；Windows＝本机 `go1.27.1 windows/amd64`。仓外导出树 `D:\tmp\wisp\20261019\342v2\{treeB,treeL}`＝`git archive fde1df02`（`winsec.go` blob `67783ae0`）⇄`9581a48c`（`07300bff`），单差异＝`internal/winsec` 落地字节（pristine 尺逐枚 sha256 并排差 5 行＝具名）。四数＝`RUN/PASS_TOP/FAIL_TOP/SKIP`（⛔ 混用⛔ 相减）。逐字读数散件＝`.scratch/wisp/probes/342/v2/**`（`readings-posix-nails.md`／`readings-windows-qualify.md`／`logs/`）；本判语每一格都带"哪把尺＋射程目录＋量的还是算的"。

## 逐格判定（`AC#0` 已由编排者凭 `342-v1` 所翻＝本腿⛔ 重翻、⛔ 评）

- **AC#1 成对读数要有 liveness ＋ 两平台各一发 ＝ 部分成立**（POSIX 半边成立、Windows 半边未达成）。
  POSIX〔量的〕：同一枚钉基线 4/4 红（`deleted=65/71/73/79` 与 `72/86/65/66`，`pinnedParentBeforeOpen=0` 逐发＝父钉那支在基线产码上不存在）⇄ 落地 4/4 绿且 liveness 逐发在场（`pinnedParentBeforeOpen=25/22/14/29`、`provedLinkCycles=171/158/160/169` 同发非零）。
  Windows〔量的〕：`truncated` 基线侧每形每发都量到（junction 0/1、symlink-dir 4/7/10＝这枚钉⛔ 恒绿）；**落地侧仍 0~2/400 漏出**（symlink-dir 三发 2/2/0 与 2/1/2＝量级稳定）。残余＝**真漏**：漏出 body 逐字＝产码载荷 `content "raced private payload 342r1 windows ancestor"`（只有 `PrivateFile` 写这 44 字节），且我另打"拔掉产码那只手"的空转对照（`return PrivateFile(...)`→`return nil`，racer/采样/复位全照跑）＝junction 两发 `truncated=0`／`returnedNil=400` ⇒ ⛔ 是仪器自造的。
  ⇒ 缺哪一味＝落地侧"修后 0"在 Windows ⛔ 成立；续体＝残余那枚窗口的确定性件（⛔ 靠 racer）。**本格⛔ 结案。**

- **AC#2 覆盖写语义不许被悄悄改掉 ＝ 成立**〔量的〕。`:96`/`:106` 旗标行两侧逐字同（`git cat-file blob`＋sha16 `cb95787e98d79789`），"名字已存在⛔ 构成拒绝理由"那段承诺（`:87-94`）两侧 sha16 `6218279b19bcff74`＝一字未动；常驻钉 `TestAC2OverwritePromiseHoldsForEveryProductionCaller322r1` 两平台四发**逐名全绿**（含 4 枚调用者子测试）；`PrivateFileExclusive` 撞已存在名两侧都仍 `isErrExist=true`／`errno:file exists`。那枚"补钉"存在且照绿＝票面这一维今天有仪器。

- **AC#3 去牙反控两发落点不同 ＝ 成立**〔量的，本腿独立打、落点⛔ 同腿报〕。两发都只打仓外拷贝树物理字节（⛔ `-overlay`）、串行、每发跑前跑后验 pristine 哈希＋落刀前后锚命中：
  发 1＝摘票 348 那支（`unix.Unlinkat(父fd, base, 0)`→`syscall.Unlink(按名)`；锚 1→0／新锚 1／mut `ec949b0d7f9f5f0b`）⇒ 三枚新钉**全部回落红**（`deleted=47` 与 `52/1500`＋确定性那枚 `the foreign witness is GONE`），覆盖写那枚照绿＝尺指向新守卫、⛔ 全盘红。
  发 2＝反形"逢名就拒"（`pinPlacement(dir)`→`refuseEveryName342v2(dir)`；锚 1→0／新锚 1＋helper 1／mut `3a1ea0eb1d0787f3`／`go vet` rc=0）⇒ `openleaf_overwrite_322_test.go:70`/`:73` 那**四枚调用者子测试逐枚红**（`secret/store.go:79`／`config/migrate.go:93`／`secret/migrate.go:189`／`agent/spill.go:166`）＋"entry stays private"＋"exclusive 仍拒"两枚也红。
  还原哈希两发都逐字回 pristine `189e55e8613249db`。⇒ 这把尺⛔ 一把两种形状都放行。

- **AC#4 门禁＋消费者不推红 ＝ 部分成立**（本腿射程⛔ 覆盖它全部判据）。
  取到的〔量的〕：门禁只在共享工作树、跑前 `unset GOTMPDIR TMPDIR TEMP TMP`＝`gofumpt -l`（我写面 0 枚 `.go`＝无入参件）rc=0／`sh scripts/d22scan.sh` rc=0 clean／`sh scripts/check-path-length-budget.sh` rc=0 VERDICT GREEN／`go vet ./internal/winsec/` 本机 rc=0＋`GOOS=linux`・`GOOS=darwin` 各 rc=0（交叉档与本机档⛔ 成对判红绿）。成对四数两平台：容器基线 `74/43/3/0` ×2 ⇄ 落地 `74/46/0/0` ×2（红名册逐名＝恰本票那三枚新钉，基线树⛔ 有＝⛔ 存量红）；本机基线 `115/61/3/0` ×2 ⇄ 落地 `115/62/2/0` ×2（落地侧两枚红＝本票截断钉＋票 348 叶子复验钉）。
  取不到的那一味＝票面那句"消费者面复认 `internal/risk`⇄`secret`⇄`config`⇄`agent`⇄`memory` 在 Linux 档与本机全 ok 的成对读数＋`cmd/wisp` 那一发"：本波派单明令⛔ 跑 `./cmd/wisp` 整包、⛔ `go test ./...` ⇒ **本裁射程外、⛔ 读成红、也⛔ 读成绿**。登记续体＝该半味由编排者在写腿全交件那一窗做（与票 227 `AC#5`/`AC#6` 同窗）。**本格⛔ 结案。**

- **AC#5 注释与票面口径说实话 ＝ 成立**〔算的＋量的〕。`winsec.go` 自陈段（`:118-140`）基线⇄落地 sha16 逐字同 `e8a2f33dea8cdeab`＝腿一字未动；本腿判 Windows 半区**未量到**（AC#1）⇒ "⛔ 出现窗口已关"这句**照实**、"仍半开"＝本格按票面第二形成立；落地三枚新件的注释⛔ 许诺未实现能力（`pinnedleaf_342_*` 自陈"WHAT THIS DOES NOT BUY"逐条对得上我的读数）。★另具名一枚**注释缺陷**（⛔ 推翻本格、供下一程只追加更正）：`pinnedleaf_342_windows.go:11-20`/`:158-170` 把祖先那一跳描述为 `OBJ_DONT_REPARSE`＋"复验看得见"，而我仓外 `os.Root` 探针实测＝**两形都答在 containment 面**（相对目标符号链接⇄junction 一视同仁，`openat sub: path escapes from parent`），仪器名册同串即旁证；`sameObjectPinnedWindows` 那两半门在该形今天到⛔ 达。

- **AC#6 纪律格 ＝ 成立**〔算的，我自己另起形数，⛔ 采信编排者数的〕。名册两把独立尺：① 单次 `%H%x09%s` 导出＋内存内 `index==1` 并集（锚只到"票 342 落地腿 342-r1"）＝**17 笔**（裸起头锚 0 笔；非锚定子串 26 笔＝A1045 那味在本腿尺上复现、⛔ 采；首 `6c010667`／末 `9581a48c`／父 `510a1519`）；② 整腿 `git diff --name-status 510a1519..9581a48c`＝**17 枚路径**、唯一 `.go` **10 枚**＝9 新＋1 改、`winsec.go` 20/3（逐枚新旧分类见 `probes/342/v2/logs-02-legdiff.txt`）。逐枚 `git show --name-only` ⊆ 声明写面＝越界 **0**。腿自抓台损 7 枚＋去牙发 2 第一版⛔ 编译（四数全 0）＝逐枚具名换落点重打，**⛔ 任何一枚未归因读数被它当凭据**。本腿另具名 3 枚仪器台损（表三落地侧 `flipErrs` 2.4~2.7 万空转、拔手对照 symlink-dir 复位死、基线 shot3 junction 复位死 `The directory is not empty`）＝⛔ 计入判据。票框终量＝本腿**一枚没翻**（`[x]`=1／`[ ]`=6＝与起手逐字同）、票名／`Status:` 未动、票尾只追加一条。

## 调用者可见的错误分类差（本腿签独立结论＝D37 错误模型 ⇄ C26 边界语义＝契约邻接，归编排者摆机主）

- 机制确认〔算的〕＝`internal/winsec/winsec.go:430` `wrapPath` 用 `fmt.Errorf("%w: %s: %v", ErrNotSealable, path, err)`：`%w` 只包 `ErrNotSealable`、errno 走 `%v` 字符串化 ⇒ 不可穿。派单那句"腿说在 `winsec_other.go`＝位置报错"我复核＝**对**（函数住在 `winsec.go`＝本就在 `A1046` 批准写面内）。
- 实测差〔量的，仓外探针件同时打进两棵树，⛔ 入库、⛔ 动跟踪字节；`logs/v02-errclass-tree{B,L}.raw.txt`〕：**POSIX**"父目录不存在"那一跳基线＝`isErrNotExist=true / isNotSealable=false`／`PathError{Op:open errno:no such file or directory}` ⇒ 落地＝`isErrNotExist=false / isNotSealable=true`／verbatim `winsec: cannot apply a private security descriptor: …/no-such-dir: openat no-such-dir: no such file or directory`。**Windows 同一跳**基线＝`isErrNotExist=true`／落地＝`isNotSealable=true ＋ isUnresolved=true`（`ErrUnresolvedPath` 经 `placementRefusalWindows`→`preOpenRefusal` 双 sentinel）。其余七跳两侧逐字同（覆盖写 `err=nil`、独占 `isErrExist=true`、叶子/祖先链接 `isUnresolved=true`、空串/相对、`SealFile` 缺失 `isNotSealable=true` **两侧都 true**＝那一跳本来就⛔ 可穿）。⇒ ★形状更正：差⛔ 止于 `wrapPath` 一处——`pinnedParent` 的 POSIX 支直接 `return err`（`pinPlacement` 的 `wrapPath`）＝与票面说的同一枚机制；Windows 支多带 `ErrUnresolvedPath`＝语义更宽，须一并摆。
- 消费者侧〔算的，全仓逐枚现读：`logs/v01-errnotexist-consumers.txt`＝116 命中／产码 35／`.scratch` 38 归别腿散件⛔ 消费者；`logs/v03-classify-decision-points.txt`＝产码判定点 38 枚逐枚 ±6 行窗口〕：**零枚**真消费者把 `PrivateFile`/`PrivateFileExclusive`/`SealFile` 的返回错误喂进 `fs.ErrNotExist`／`os.IsNotExist`。判定点的来源逐枚＝`os.Stat`/`os.ReadFile`/`os.Lstat`/`os.Remove`/`filepath.WalkDir`（`secret/store.go:106/140/162`・`secret/migrate.go:88/168`・`memory/artifacts.go:115~307`・`config/loader.go:67`＋`fileMissing`(`parse.go:233`) 三处・`cmd/wisp/firstrun.go:76`・`proc/envfork.go:171/221`・`panel/git.go:318/421`・winsec 自己的 `resolve.go:264`/`winsec.go:251`/`winsec_windows.go:683/732/748`）；`memory/artifacts.go:204/214` 那两枚喂的是 **`RemoveUnlinked`**（⛔ 本票路由、落地一字未动）。测试侧同样零枚（逐窗口扫描空）＝⛔ 一枚现有钉的绿/红依赖这一跳的 `ErrNotExist`。
- **我的签＝〔修回 errno 可穿〕**（⛔ 我动手＝产码＝须另批；★本格⛔ 选形、⛔ 摆选单，只有裁决与形状）。理由：① `PrivateFile` 的 `ErrNotExist` 是**分类正确性**问题（D37 (a) 那张表把 `resource`/`config`/`permission_denied` 分开、`errors.Is` 是唯一可编程分类通道；把 ENOENT 报成 `ErrNotSealable` 让调用方看⛔ 见"路径不存在"）；② 它⛔ 是"已批准落地"的必要件（零消费者依赖＝本票⛔ 因它不成立）；③ 改动形状＝把 `winsec.go:430` 的 errno 那半从 `%v` 换成 `%w` 再包一层（`fmt.Errorf("%w: %s: %w", ErrNotSealable, path, err)`），POSIX／Windows 共用同一枚 `wrapPath`＝**一处**；★连带后果须先具名＝`sealError` 的幂等判定（`errors.Is(err, ErrNotSealable)` 那支）⛔ 变、`TestAC5FailedSealRefusesTheWrite`（`private_fail_test.go:37/50/61/72/75`）仍绿、`openleaf_322`/`placement_symlink_113` 那些读 `ErrUnresolvedPath` 的钉⛔ 受影响；另一面照实登记＝**现行为今天无任何消费者依赖**＝这跳的分类既⛔ 被钉住、也⛔ 被读过。
- ⇒ **归"待人拍板清单"这一枚事实**＝它是 D37 错误模型 ⇄ C26 边界语义的契约邻接面、agent 单方面改⛔ 允许；本腿⛔ 改任何 `.go` 字节，等编排者落人工批准那一拍。

## 一句话终态

票 342 的 POSIX 半边＝**关掉并量到**（本腿独立重打 8 发＋去牙两发照红＋liveness 逐发在场）；Windows 半边＝**落码入库、仪器有牙、但截断那一支⛔ 归零（残余 1~2/400＝真漏、未归因到底）＋叶子复验⛔ 量到** ⇒ `AC#1`/`AC#4` 部分成立、`AC#2`/`AC#3`/`AC#5`/`AC#6` 成立 ⇒ **本票⛔ 结案**（续体＝Windows 残余的确定性件＋消费者包色那一窗＋错误分类差拍板）。
