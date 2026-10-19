# 342-v1 终裁（票 342 `AC#0` 三形代价表复核＋两平台竞态重打｜裁决者≠实现者）

腿＝`342-v1`（票 342 第一枚终裁腿；起手 `git log --all --grep='票 342 普查腿 342-a1·'`＝0 笔、`probes/342/v1` 与本判决件起手均不存在＝三键尺全清）。
锚＝起手现量 `251ab436`；台面＝本腿自建 `git archive` 导出树 `D:\tmp\wisp\20261018\342v1\{treeV,treeC,treeH,treeW,mut-*,col-*}`（⛔ 复用 `342-a1` 任何一棵树、⛔ 抄读数）；`internal/winsec` 40 枚 blob 在 `c596df55`⇄现 HEAD `0b5bd8c6` 逐枚同串（现量 checked=40 diff=0）＝腿的行号仍可用。
★选形权⛔ 在本腿：全文无选形倾向句，也⛔ 把被禁词的字面写进正文（派单第 87 行那条禁令，本腿第一版自己也犯过一次＝那句话里带着尺要找的词，已改写；尺与正控见 `scripts/v342v1-selfcheck.sh` 的 (g) 段，正控＝同一把尺在票面 342 命中 1＝尺在跑）。细节逐枚在 `.scratch/wisp/probes/342/v1/logs/`。

## AC#0 三形代价表＝**成立**（承重竞态读数已由〔仅腿报〕升成本腿收件复跑）

1. **改动面对得上量级与形状、⛔ 要求逐字**：`winsec.go:149` 侧＝**1/1** 逐字对得上（本腿内容锚旧 1→0、新 0→1，落刀前后各验，四张落点脸 `winsec.go` sha 全部 `3c786d1d…`）。新件行数本腿自建＝甲 linux **43**／乙 linux **32**／丙 linux **33**／甲 win **54**／乙 win **53**，对腿报 `41/34/53/47`＝同量级、逐枚不同（各腿自写实现，行数不同是常态、⛔ 是缺陷）。
2. **POSIX 竞态（本腿 3 发×1500，逐次分子在 `10-race-posix-keyrows.txt`／raw `10-race-r1..r3`）**：VBASE(by-name) `damage=484/458/479`＝**敞着，复现**；VJIA(甲) `0/0/0`、VYI(乙) `0/0/0`＝**截断那一支被关掉，复现**；VBING(丙) `123/189/206`＝**报告了、没防住，复现**；liveness `flips>1000` 逐发在场（⛔ 拿 0 命中当凭据）。
3. **真入口那一格＝形状复现、量级⛔ 复现，本腿具名归因**：VPROD `damage=2/2/1`（腿报 106）。归因＝竞态节奏差异：本腿 race 线程换手 ~1000–1600/1500，`ENOENT/body` 占 1353/1500（名字多半处在"被 rename 消费掉的空档"，截断根本没机会发生）⇒ **⛔ 腿报造假**。补证＝本腿加 `V342V1_RACER_US` 节流扫（`15-vprod-pace-*.log`）：pace=0 `damage=56`、100us `18`、300us `32`、1000us `28`，且**逐发 `truncated=0 / deleted>0`**＝真入口残余几乎全是票 348 那一支（unlink），不是截断；`flips≈90–270` 与腿报 358 同档。
4. **Windows 竞态（本机 `go1.27.1 windows/amd64`，3 发×1500，`31-win-race-keyrows.txt`）**：WBASE `damage=975/871/878`（腿报 897）＝**复现**；WPROD `7/6/4`（腿报 3）＝同档；WYI(乙 `os.Root` 链) `0/0/0`、`root-refuse≈1341–1390`＝**复现**；★**WJIA(甲) `damage=3/5/5`＝腿报 `0` 没复现**，本腿读数＝**甲在 Windows 没把截断那一支关掉**（残余小、非零），与"判据能关但仍是 check-then-act"自陈**相容**⇒ 判该格应记"残余非零"，⛔ 记 0。
5. **顶红名册＝实测 0 红，两平台各自复现**：容器 B 档（⛔ `GOTMPDIR` bind mount，见第 7 条）基线 `71/43/0/0` **安静 3 发逐字同号**；甲/乙/丙落点脸**逐字同号 `71/43/0/0`、红名册 0 枚**。本机 `treeH` 基线 `110/61/0/0`；甲/乙 Windows 落点脸 `110/61/0/0` **逐字同号、0 红**。⇒ 腿报"两平台 0 顶红"成立。
6. **尺敏感（必答 iii）＝复现**：误形 `yinofb`（钉住父目录＋叶子强 `O_EXCL`）打出 **4 枚真红**、四数 `71/39/4/0`，与腿报名册**逐名相同**：`TestAC1POSIXOpenThroughASwappedLeafLinkNeverTruncatesAForeignFile322r1`／`TestAC2OverwritePromiseHoldsForEveryProductionCaller322r1`／`TestAC1POSIXPrivateFileRefusesWhenNameCarriesAnotherRealObject120r3`／`TestAC1POSIXPrivateFileRefusesWhenNameBecomesLinkAfterPinnedSeal120r3`。★**一处具名更正（腿报红句形状）**：本腿实测 `322` 那枚钉的红**不是** `:199/:200` 的"damaged 红"，而是 `openleaf_322_other_test.go:206` 的 **RED instrument（仪器死）**：逐字 `AC#1 RED instrument: across 3000 raced attempts the leaf leg ("too many levels of symbolic links") answered 0 and the racer landed 562 flips (0 errors, first ""), so the damaged count above is not evidence about a window that was never reached`，且同一发 `damagedWitnesses=0`。⇒ 结论方向不变（尺会照红），但"它红在损伤支"⛔ 成立、红在 liveness 支；覆盖写那枚才真红在损伤：`openleaf_overwrite_322_test.go:70` 逐字 `AC#2 RED: PrivateFile refused an existing name (open …/blob-322r1.bin: v342v1 mis-shape on the pinned parent: file exists), which breaks storing the same id twice replaces the blob in place`（四枚子用例 :70/:73 齐红、`:98/:101` 第五枚），`sealhandle_privateroute_120_other_test.go:83/:85` 与 `:110/:120` 两枚＝`the refusal does not carry ErrNotSealable/ErrUnresolvedPath … want the leg carrying "this descriptor is not the object the name carries now…"`。
7. **本腿自抓一枚台损（⛔ 计入判据、⛔ 记成祖先那一支新造的红）**：本腿首发容器档多挂了 `-v $S/gotmp:/gotmp`＋`GOTMPDIR=/gotmp`（落 D 卷＝drvfs/9p）。该档下 pristine `treeC` 读到 `71/41/2/0`，两枚红＝`TestAC1POSIXSealHandleNameCarriesAnotherRealObject120r3` 与 `TestPOSIX120SealWindowAncestorRelativeLinkIsClosed`。归因＝独立正控（`probe2/samefile.go` 容器内跑）：同一份 rename＋`os.SameFile` 判定在 `/tmp`＝`true`、在 bind mount＝`false` ⇒ `sealHandle` 的身份门在该台面上被内核语义摘掉，红是**台面**造的。改回与腿 `driver-g3.sh` 同 env（只 `GOMODCACHE/GOCACHE`）＝`71/43/0/0`×3。⇒ ★**给编队的定式**：POSIX 容器的 `TMPDIR/GOTMPDIR` ⛔ 绑 Windows 卷，否则 `SameFile`/mode 回读失真、造出长得像判据响的红（与票 311 的"卷形"红同源、新面）。

## AC#2 覆盖写语义（本腿逐枚现读调用点，⛔ 抄表）＝**七枚全对**

`PrivateFile` 5 枚：`secret/store.go:79` 依赖覆盖写（文档 `:59-63` 逐字 "creating or **overwriting** it"）；`config/migrate.go:93` 依赖（注释 `:89-92` "because this call **rewrites an existing backup** … the **repair path**"）；`secret/migrate.go:169` **不**依赖（`:168` 有 `os.Stat` 前置，存在即走 `SealFile` 支 `:174`）；`secret/migrate.go:189` 依赖（pre-rename temp 固定名 `tmpPath`）；`agent/spill.go:166` 依赖（`.retry` temp＋rename，`:158-172` D31 last-writer-wins）。`PrivateFileExclusive` 2 枚要的就是 `fs.ErrExist`：`spill.go:343`（文档 `:331-337` "FAILING with fs.ErrExist"）、`memory/artifacts.go:75`（`:67-70` "not an overwrite"）。⇒ **4 枚依赖覆盖写＝成立**；`openleaf_overwrite_322_test.go` 逐用例点名这四枚（`:43-59` 四行 label），本腿 `yinofb` 一发把这四枚全打红＝**这枚钉今天咬得住**。★两处票面具名更正本腿核**对**：旗标行 `:87/:97`→现读 `winsec.go:96`（`O_WRONLY|O_CREATE|O_TRUNC`）／`:106`（`O_EXCL`）＝漂 −9；"这个承诺⛔ 任何一枚现有钉钉住"→已被 `322-r1` 的 `openleaf_overwrite_322_test.go`（无 build tag＝两平台同件）推翻，`AC#2` 今天有钉。⛔ 由本腿改票面（只追加、⛔ 撤框）。

## 必答 (ii) 两平台各自

- **甲在 Windows**：判语＝**"判据能关但仍是 check-then-act"成立**（现读 `winsec_windows.go:361-369` 自陈"the probe cannot hold the answer … Closing THAT needs a single non-following create-and-truncate, which this platform was measured not to offer through the flags Go exposes; it is booked, not claimed closed"＋`:374-395` 那枚 `CreateFile(…, FILE_FLAG_OPEN_REPARSE_POINT)` 只问不锁）；★再补一味：本腿实测 `WJIA damage=3/5/5`（liveness `flips=59/56/63`）＝**残余是真会动字节的、⛔ 只写在注释里**。
- **乙在 Windows＝真缺世界（⛔ 判探针缺陷）**：现读＝`winsec_windows.go:413` `sealHandle(f) = applyDescriptor(f.Name(), false)` 走名字、`pinPlacement/sameObjectPinned/pinBase` 全在 `winsec_other.go` 的 `//go:build !windows` 后 ⇒ 本包 Windows 侧**没有**该族可复用（＝新写平台码，`v342v1-landing-yi-win.go.txt` 就是那枚最小新写件）。本腿自证⛔ 重复腿那枚 off-by-one：`pathComponents` 不含 volume（`winsec.go:394` 从 `len(filepath.VolumeName(path))` 起），本台下探从 `comps[0]` 起＋`OpenRoot(vol+"\")`＝`WYI damage=0/0/0`。★新读出一味本票**必须**记的：`W-ROOTSYM followed=true`（3/3，容器外本机）＝`os.Root` 对**树内相对 FILE 符号链接照旧跟随**、`W-ROOTJUNC`＝对 junction `path escapes from parent` 直接拒 ⇒ **乙在 Windows 钉祖先、⛔ 钉叶子**，落地仍需 `sameObjectPinned` 同款叶子复验（Windows 今天⛔ 没有）＝工作量⛔ 止于"新写一条链"。

## 必答 (iv) 在册钉抽查（逐行现读两枚＋一枚范围裁语）

`placement_symlink_113_other_test.go:189-205`＝硬判据两半逐字对得上：`:197` `errors.Is(err, winsec.ErrUnresolvedPath)` 红句 "…which does not name ErrUnresolvedPath…"、`:201` `refusalCreditsLink137(err, planted)` 红句 "…does not credit the link this case planted…"；`:277/:286` 经 `winsec.PrivateFile` 驱动、反半 `:376-377` "PrivateFile refused a plain path inside the named tree"⇒ 腿报"甲/乙把裸 `*PathError` 直送会打红 :197/:201、逢名就拒会打红 :377"与**实测落点 0 红**并存＝形状自洽。`openleaf_overwrite_322_test.go:70/:79/:107` 逐字对得上（`:79` 是 mode 收窄那半，Windows 侧 `privateModeIsReadable=false` 走 `:82` 的装饰性日志）。★`openleaf_flag` ⛔ 经 `privateFile`＝**成立**：头注 `openleaf_flag_322_other_test.go:10-17` 逐字 "it calls `os.OpenFile` by name and never goes through `privateFile` … the two mutations … leave THIS case green while the race case next door goes red" ⇒ 那枚绿⛔ 读成"落地接线正确"的凭据。

## 必答 (v) 票 348 归谁＝**编排者处置成立**（附一枚落地交互代价）

现读四行逐字＝`winsec.go:158-162`：`if err := sealHandle(f); err != nil {` / `_ = f.Close()` / ★`_ = os.Remove(resolved)` / `return sealError(resolved, err)`；★**同一枚路由今天有两处 `os.Remove`**＝`:160`（密封失败支）与 `:165`（写失败支），票 348 票面只点名 `:160` ⇒ 落地⛔ 只改一枚。复跑 A3（本腿自己建台、3 发确定性、⛔ 竞态）逐字复现：`A3 row1: sealHandle err=winsec: cannot apply a private security descriptor: …: openat sub: path escapes from parent witness-after-seal=truncated`／`A3 row2: os.Remove(…) err=<nil>; foreign witness now=deleted(size -1); ancestor link still standing=true(<nil>)`／`A3 VERDICT: … deleted the FOREIGN object = damage shape 2, deterministic, not a race`＝**盘上事实**。裁"⛔ 并进 342、但落地同批交同一枚写腿"＝**认同**：两枚不同系统调用（`open` 截断⇄`Remove` 删除）、两种损伤、两个落点行，且本腿 `VPROD-LANDED` 逐发 `truncated=0/deleted>0`（甲 `58/49`、乙 `50/54`）＝**甲乙都关不掉 348 那一支**，并进 342 会把"选形门"与"删除形"两件事捆成一件事、且票面⛔ 由本腿改。★**顶回来的只有附加条件**（⛔ 换形）：票 348 的甲形（`unlinkat(fd, base, 0)`）**依赖 342 选形**——本腿甲形模拟件在叶子开完后 `_ = syscall.Close(fd)`（`v342v1-landing-jia-linux.go.txt:36`），要 fd-relative 删就**必须留父目录 fd 到路由末端**；乙形（`pinPlacement` 已交回 `*os.File`）几乎免费。⇒ 同批落地时两票的甲形行数⛔ 按 342 单票估，须按"甲/乙"分别重报。

## 必答 (vi) 三形对照（每格带凭据；⛔ 选形、⛔ 建议）

| 维 | 甲（逐组件 `openat`） | 乙（复用密封 pinning） | 丙（动手后复验） |
|---|---|---|---|
| 防住截断？POSIX | 是：`VJIA 0/0/0`＋`VPROD-LANDED truncated=0` | 是：`VYI 0/0/0`＋`VPROD-LANDED truncated=0` | **否**：`VBING 123/189/206`、B 档 `523/498`＝损伤照发生 |
| 防住 unlink？ | **否**（`:160/:165` 未动，残余全 `deleted`） | **否**（同左） | **否**（同左，且报告滞后一支） |
| Windows 要不要新写平台码 | **要**（`oNoFollowPlatform=0`，逐前缀复探＝新写；本腿件 54 行）且**判据形仍非原子**：`WJIA damage=3/5/5` | **要**（本包 Windows⛔ 该族，`winsec_windows.go:413`；本腿件 53 行）且**钉祖先⛔ 钉叶子**：`W-ROOTSYM followed=true` | POSIX 侧今天已在跑（`winsec_other.go:287-307` 两道身份门）；Windows 侧⛔ 无人复验祖先 |
| 要不要新导出名 | ⛔ 不新导出，但**收窄可达面**：`syscall.Openat` 只在 linux（`darwin/freebsd` 走不通；`openbsd/netbsd/solaris/plan9` 今天 `go build` 已 rc=1＝`fstypeOf`） | ⛔ 不新导出（`pinPlacement/sameObjectPinned/pinBase` 同包非导出、签名对得上） | ⛔ 不新导出 |
| 改动面 | `1/1`＋43 行（linux）／54 行（win） | `1/1`＋32 行（linux）／53 行（win） | `1/1` 后插复验＋33 行 |
| 会不会摘掉谁的既有牙 | ⛔ 没摘：两平台落点色逐字同基线（`71/43/0/0`、`110/61/0/0`） | ⛔ 没摘（同上）；误形 `O_EXCL` 才摘→4 枚真红 | ⛔ 没摘（`col-bing 71/43/0/0`） |

## 必答 (vii) `AC#1..AC#6` 逐格〔今天可判／落地后可判／缺世界〕

- **AC#1**（成对读数＋liveness、两平台各一发）＝**今天不可判**：缺"修后"的常驻仪器（⛔ 本腿⛔ 落码）。缺的料已备齐＝驱动 `PrivateFile`＋祖先换手的竞态钉（本腿 `v342v1-posix-race.go.txt`／`v342v1-windows-forms.go.txt` 即形状雏形，含 liveness 钥匙：POSIX＝`flips>0` 才不许 DEAD INSTRUMENT，Windows＝`reparseAt(ancestor)` 复证）。★两平台各一发的**分叉**今天已量到：同一形在 POSIX＝原子（`openat(O_NOFOLLOW)`）、在 Windows＝非原子（`WJIA` 残余非零）⇒ ⛔ 一发代两发（票 322 `AC#3` 那味）。
- **AC#2**＝**今天可判形、不可判终态**：七枚名册＋那枚钉在场（上节）；仍缺的裁定＝"祖先守卫是否连带把覆盖写改成拒绝"⛔ 由选形决定；本腿 `yinofb` 证明**若**改成拒绝，四枚调用者立刻红。
- **AC#3**（去牙反控两发落点不同）＝**今天不可判**（落地字节不存在）；★但反控起点本腿已备好且**两发落点不同**已示范：一发＝误形旗标（`yinofb`，打红 4 枚＝仪器活），一发＝内容锚换 open 调用（`1/1`）⇒ 落地腿⛔ 只用一种。
- **AC#4**（门禁＋消费者）＝**今天不可判**；本腿先把门禁四数打平（下节），消费者面缺的是**落地后**的 `risk/secret/config/agent/memory` 成对读数＋`GOOS=linux`/`darwin` 交叉档（`A1032` 第③节：Windows 宿主交叉跑 `./cmd/wisp` 必然红＝⛔ 与本机档成对判红绿）。
- **AC#5**（`winsec.go:118-140` 自陈随终态更新）＝**今天不可判**（无终态）；★本腿核：该段自陈**逐字仍然说实话**（`:133` "This window is not claimed closed anywhere."、`:134-140` Windows 两支 residual）＝若选丙，⛔ 出现"窗口已关"这一条今天没有违例；若选甲/乙，更新⛔ 只写 POSIX 一支（Windows 甲＝check-then-act、乙＝⛔ 钉叶子）。
- **AC#6**（纪律格）＝**今天不可判**（要落地腿提交面）；★本腿自证同尺合规：写面⊆三处、`.go` 命中 0、⛔ 翻框、⛔ 改票名／`Status:`、只追加一条票尾 bullet。
- **落地腿需要的人工批准口令形状＝够、但缺一味**：腿给的形「批 342 落〈甲|乙〉：允许改 `winsec.go`＋新 landing 件；`winsec_windows.go` 是否同批单独点名」⛔ 覆盖 ①票 348 的 `:160` **与 `:165`**（两处 `os.Remove`，只点 `:160` 会漏半）、②"留父目录 fd 到路由末端"这一形变（甲形才需要）、③乙形 Windows 的**叶子身份复验**（`W-ROOTSYM` 那一支）。⇒ 那三味⛔ 由本腿自填，具名报编排者拍（本条＝必答 (vii) 要的"核这形够不够"的答复，与选形无关）。

## 两件事的裁语（派单点名）

- **`342-a1` 越界那处（跑了 `./internal/risk/`⇄`./internal/config/`⇄`./internal/agent/` 色彩）＝处置够、⛔ 判谎报成立**。理由：①载体只有 `logs/29-*`／`32-*`／`33-*`，射程全在**消费者维**＝票面 `AC#4`，而 `AC#4` 今天本来就⛔ 可判 ⇒ **退这一维＝退掉的东西今天没有承重**；②它自己具名"⛔ 当判据"＋基线树 3/3 红＝与落点无关，本腿核其"⛔ 谎报"成立（自抓＝本仓要的形状）。★**污染别的格＝没有**：本腿把它三张承重脸全部重打且**逐字同号**（容器 `71/43/0/0`、本机 `110/61/0/0`、`yinofb` `71/39/4/0`＋同名四枚）⇒ 它的 `AC#0` 读数⛔ 受那三枚色彩影响。★**唯一不许延伸的一格**＝"摆动枚安静复量 3 发全绿"：那两枚名（票 346）本腿在 `GOTMPDIR` bind 档读到**真红**（第 7 条）＝它们对台面配置敏感 ⇒ ⛔ 拿它们的绿当"叶子已钉住"的凭据，⛔ 拿本腿 B 档的绿当"祖先那一支已钉住"的凭据（两枚都量不到祖先那一支：`posix_seal_window_120_other_test.go:479` 驱动的是 `SealFile`，本腿尺 `PrivateFile` 命中 0）。
- **票 348＝⛔ 并进 342 成立、同批同腿成立**；★**新报两味**：同一枚路由第二处 `os.Remove` 在 `winsec.go:165`（写失败支，同样沿祖先链接删）＝票 348 票面只点 `:160`；且甲形的 `unlinkat` 与 342 选形有**依赖耦合**（见上）。⛔ 本腿立新票，具名报编排者。

## 门禁四数（只在共享工作树，跑前 `unset GOTMPDIR TMPDIR TEMP TMP`，逐枚 rc）

`gofumpt.exe -l .scratch/wisp/probes/342/v1` **rc=0**（输出空＝本腿写面⛔ `.go` 件；正控 `-l internal/winsec/winsec.go` 同 rc=0＝尺在跑）／`sh scripts/d22scan.sh` **rc=0**（`clean - no D22 ban violations`）／`sh scripts/check-path-length-budget.sh` **rc=0**（`VERDICT GREEN`，worst 229 chars、roster 60＝存量）／`go vet ./internal/winsec/` **rc=0**（无输出）。★票框终量尺：`grep -c '^- [ ] \*\*AC#'`=7、`grep -c '^- [x]'`=0 ＝**七格一枚没翻**（本腿⛔ 翻框）。两平台四数（⛔ 混用⛔ 相减）：容器 B 档 `71/43/0/0`、本机 `110/61/0/0`，落点脸各自逐字同号。
