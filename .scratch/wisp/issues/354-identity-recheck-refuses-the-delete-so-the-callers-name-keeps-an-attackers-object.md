# 354 — 同一性复验拒删之后，调用者那枚名字原地留着【攻击者自己的对象】

- Status: **ready-for-agent**（★但 `AC#1` 拍形**必须先过编排者＋机主**，⛔ 未拍形任何实现都⛔ 开工）
- 立案：2026-10-21 09:2x +08 编排者立。**非实现者签名**＝`353-v1` 判语第 5 件（`docs/evidence/s1/353-v1-verdict.md`）；
  草案＝`.scratch/wisp/probes/353/v1/02-ticket-draft-secC.md`（七格判据原文）。★**本票⛔ 由实现腿推荐立**（规矩同 `342-v3`）。
- 来源：票 353〔乙〕落地（机主 2026-10-20 19:2x 拍形）认下的**残余**。票 353 的判据只读四枚分子
  （`entered`／`damage`＋正控），那两族（`ancestorunlink_348_*`／`ancestorworld_342_348_*`）**⛔ 在本票射程**。
- 关联：票 353（已 `-done`）· 票 348／342（两族仪器宿主）· 票 118／120-r3（"link 只 unlink、⛔ 跟随"＝已批、均 `-done`，本票**⛔ 重述**）

## 这是什么（人话）

〔乙〕把"拒绝路径按名删会删掉调用者⛔ 点名的对象"这件事关死了：删除前先做一次"名字⇄句柄对象"同一性复验，**不匹配就不删**。
代价是——**不删之后，攻击者顶进来的那个对象就原地留在调用者那枚名字下面了**。调用者以为那个名字里还是自己的文件，其实已经是别人的。

这不是"窗口没关"（那个窗口要钩进平台内部的 `NtOpenFile`，本机测不到、票 353 已逐字承认）；
这是**关窗之后的残留**：窗口关上了，但屋子里留下了一件不属于你的东西，而门牌还写着你的名字。

## 结案判据（七格，判据原文在草案 SEC-C＝`.scratch/wisp/probes/353/v1/02-ticket-draft-secC.md`）

- [x] **AC#0 现读＋归因（读腿，⛔ 动产码字节）**：把"残余①留着攻击者对象"那一跳今天读成读数并归因。
      ★两族读数必须**现跑**、⛔ 抄任何腿的数（`353-v1` 量到的 `flipErrs` 基线 `9,130..31,366` ⇄ 落地 `6..16`；
      `provedReparse` 基线 `2..98` ⇄ 落地 `102..153` 只是"待复跑"的线索，⛔ 当凭据）。
      必须先跑**正控**证明尺能命中非零（本仓已有"钝尺读 0"的多次教训）。
      - 〔354 收口腿 2026-11-08 21:3x +08〕翻勾（凭 `docs/evidence/s1/354-a1-residual-readings-r1.md`＋收口腿自己复跑，读腿与收口腿分属两枚 agent＝非实现者裁）：导出树 `git clone --shared`（HEAD `66a03a320`；`ca790c46e..HEAD` 于 `internal/winsec` 零字节）定向 `TestAC2Windows` 六发全 PASS，关键读数逐字同我手——`truncated=0 deleted=0`×2（symlink-dir／leaf-symlink 两形）、`refusals that left the caller's name standing=3`、DAMAGE 行 3 发逐字 `the caller's name=present (size 64 body [an object the caller never named, standing in the same directory])` 且 `the swapped-in object's own name=absent`×3（top-up 后调用者名字＝外来对象唯一名字）、负控 `RemovalStillRuns` `name=absent`（尺两向非钝）、`ancestorunlink_348_windows_test.go:169` 残留句运行期打出（bytes refused、NAME never removed）；`./internal/winsec` 整包 `ok 20.631s` rc0。归因对上＝复验拒的那支在 Remove 前返回（`unlinkPinnedLeaf` 对码 353-v1 第 1 件）⇒ 残留①是 owner 拍〔乙〕可观测方向的代价形、非回归；GREEN SCOPE（recheck 后 NtOpenFile 前 top-up）照抄不冒认。★unix 半边判＝**本格不要求**：AC#0 原文无「两平台」字样（同票 AC#4 才写「两平台成对」）；且残余①在 unix 机制上不存在——`pinnedleaf_342_unix.go:173-174` 的 `unlinkPinnedLeaf` 无同一性复验（直接 `Unlinkat`），「复验拒删留下名字」那一跳 unix 无从发生（`353-v1` 判语同口径：unix 行为照旧、unix 是否也该复验已上交未拍）。判语全文＝`docs/evidence/s1/211-354-closure.md`。
- [x] **AC#1 拍形（编排者＋机主）**：三形各写代价＝〔甲〕复验拒删后**另走一条删除路径**（碰语义）；
      〔乙〕把"留着的那枚名字"**主动报出来**（只加读数、不删）；〔丙〕登记为已知残缺＋常驻仪器把 `leftStandingForeignName` 钉成可见红。
      ★本格**⛔ 由 agent 自决**：代价表出来后摆机主拍；`AC#2` 的形状以拍定那一形为准。
      **机主拍板（2026-10-21 14:2x +08，原话「按推荐」）＝〔乙〕：把"留着的那枚名字"主动报出来（只加读数、不删）。**
      理由（三条，均可在盘上核）：①**不碰删除语义**——甲会二次改已批落地形（票 348 那条），代价最高；
      ②**可观测**——乙把"名字下站着谁"变成一条读数，比丙的"常驻红"温和且不挡 CI 合并；
      ③**与票 353〔乙〕同族**——票 353 已经用"复验拒删"把窗口关上了，本票只是把**关窗之后留下的东西**说出来。
      ⇒ `AC#2` 的形状按〔乙〕定：**只加读数、不动删除路径、不加常驻红**；⛔ 顺手改删除语义、⛔ 借本票重述 342／348。
- [x] **AC#2 确定性仪器（⛔ racer 作主判据）**：要量的那味＝**拒绝路径之后调用者那枚名字下面站着什么**，
      分子建议 `leftStandingForeignName`；⛔ 借 `entered`／`damage`／`deleted`／`truncated` 任何一枚既有分子。
      尺**必须带轮内序号**（`353-v1` 具名：同轮同树两发写同名原件 ⇒ 12 发里 5 行原件不可复原）。
- [x] **AC#3 敏感度＝成对去牙**：摘复验那步 ⇒ 回到 `damage>0`；摘删除那步 ⇒ damage 变体各自红；落点⛔ 同（⛔ 单发当判据）。
- [ ] **AC#4 颜色与门禁**：`./internal/winsec/` 两平台成对＋逐名红册（摆动名册照公共块）；四道门禁 rc；★⛔ 跑 `./cmd/wisp` 整包。
- [ ] **AC#5 与上游两票的接缝**：那两族的**既有红⇄摆动名**（`…NeverTruncatesA…` 那族）今天什么状态，逐名写清；
      并写明本票**⛔ 重述**票 118／120-r3 已批的"link 只 unlink、⛔ 跟随"。
- [ ] **AC#6 纪律格**：⛔ 翻任何框（本票各格由非实现者裁）、⛔ 改票名／`Status:`／`-done`、⛔ push、⛔ 动 342/348 两票任何原句。

## 硬约束（逐条照做）

- ⛔ **重述**票 118／120-r3 已批的"link 只 unlink、⛔ 跟随"；那两票已 `-done`（尺＝`git ls-tree HEAD -- .scratch/wisp/issues/ | grep -c`）。
- ⛔ 并进票 353／342／348 任何一格；⛔ 动那两票任何原句。
- ⛔ 借本票顺手"修绿"：把 `damage` 断言反写成"有顶入即拒绝"＝新语义、= `AC#1` 拍形才许谈。
- ⛔ 碰 `:11-20`／`:158-170` 那两处机制归因注释（票 342 的〔丙〕＝另批）。
- ⛔ `internal/agent/**`、`cmd/wisp/**`、`frontend/**`、`design/**`、`thresholds.go`、任何 golden。
- ⛔ 写面冲突：`internal/winsec/**` 同包**只许一枚写腿**。
- ★**下一位只拿 `alt-*.raw.txt` 当全部凭据**；本次发现"同轮同树两发写同名原件⇒12 发里 5 行原件不可复原"⇒
  新尺**必须带轮内序号**（`353-v1` 具名，草案 SEC-7）。

## 排程

`354-a1`（只读＝`AC#0` 现读＋归因）→ 编排者＋机主拍 `AC#1` 三形 → `354-r1`（写腿＝`AC#2`＋`AC#3`）→ `354-v1`（非实现者裁）→ 色与门禁（`AC#4`）。

**Status 追加（2026-10-21 09:2x）**：★号位＝编排者定（两枚草案都建议 `354`＝撞号；按"先交件先占号"＋"安全残余优先"判给本票，
`355` 给 `139-v1` 那枚"痕到不了盘／盘上无读者"）。⛔ 零翻框、零 push。

## Progress log

（追加式；每条 `- [日期 时刻 +08] agent=… did=… next=…`，钟点由 `date` 插值，不手打。）

- [2026-10-21 09:2x +08] agent=`编排者` did=立案（非实现者签名＝`353-v1` 第 5 件；草案 `02-ticket-draft-secC.md` 七格原文照搬＝本票判据）next=派 `354-a1` 只读腿出 `AC#0` 现读；`AC#1` 三形代价表回来后摆机主拍形。
- [2026-11-08 21:1x +08] agent=354-a1（AC#0 只读腿＝现读＋归因；零产码零框翻；写面＝docs/evidence/s1/354-a1-residual-readings-r1.md＋.scratch/wisp/probes/354/a1/{readings.md,01-ac2windows-instruments.log,02-winsec-full-package.log}＋本条）did=两族仪器现跑全绿（导出树 ca790c46e，两棵 clone --shared 同锚）：ancestorunlink_348 Windows 两形 truncated=0/deleted=0（symlink-dir cycles=67/flipErrs=11；leaf-symlink cycles=81/provedReparse=83）且 :169 残留句逐字（bytes refused、NAME never removed）；nametopup_353 三发 refusals that left the caller name standing=3/3＋DAMAGE 行打出调用者名字下 64 字节攻击者对象本体（其自己名字=absent＝唯一名字）；负控 RemovalStillRuns name=absent（尺两向非钝）；正控先行＝shot 级 body 逐字非零，敏感度引 353-v1 去牙三发未重跑（归 354-r1 AC#3）。归因＝复验拒的那支在 Remove 前返回（unlinkPinnedLeaf 对码 353-v1 第 1 件）⇒ 残留①是 owner 拍〔乙〕可观测方向的代价形、非回归；GREEN SCOPE 那支（recheck 后 NtOpenFile 前 top-up）仍测不到照抄。unix 半边本机不可跑＝登记未跑。winsec 整包 ok 20.148s rc0（02 号日志）。⛔ 翻框（AC#0 归编排者）⛔ 动 internal/winsec 字节（单写腿留给 354-r1）⛔ push。next=编排者收件裁 AC#0；354-r1（AC#2 分子＋AC#3 去牙）待派，winsec 现量 0 在飞。
- [2026-11-08 21:3x +08] agent=354-cls（AC#0 收口腿：凭 354-a1 证据件＋自己复跑翻 AC#0 一格；零产码；写面＝本票 AC#0 翻勾＋格下判语＋本条＋docs/evidence/s1/211-354-closure.md）did=三核全过（票面 AC 原文现读在；证据件＋readings.md＋01/02 原始日志现读；色取自 clone --shared 导出树 HEAD `66a03a320`，`ca790c46e..HEAD` 于 internal/winsec 零字节）。复跑：定向 `TestAC2Windows` 六发全 PASS（我手，非抄读数）＋关键读数逐字对上（truncated=0/deleted=0×2、standing=3、DAMAGE 64 字节攻击者对象×3、own name=absent×3、负控 name=absent、:169 残留句）＋winsec 整包 ok 20.631s rc0。unix 半边判＝AC#0 不要求（原文无两平台字样、AC#4 才有）且残余①unix 机制不存在（pinnedleaf_342_unix.go:173-174 无复验直接 Unlinkat，353-v1 同口径）。AC#0 翻勾；票⛔ `-done`（剩 AC#2/AC#3＝354-r1 写腿、AC#4 色与门禁、AC#5 接缝、AC#6 纪律格）。⛔ push。next=354-r1 待派（winsec 单写腿）。
- [2026-11-08 22:2x +08] agent=354-r1（写腿：AC#2 确定性仪器＋AC#3 成对去牙；写面＝internal/winsec/leftstanding_354_windows_test.go 新建＋docs/evidence/s1/354-r1-deterministic-teeth.md＋本条；产码非 test 字节＝0）did=AC#2＝分子 `leftStandingForeignName` 落成仪器（世界 342r2 逐字复用零新接缝；固定 64 字节 body 长度断言；夹具与读数行带轮内序号；幸存名读回逐字节 `size 64 body [an object the caller never named, standing in the same directory]`＋own name=absent 才计数；负控 RemovalStillRuns354r1 双向非钝；`-count=3` 九行逐字节同＝旁证）。AC#3＝两牙只在导出树（/d/tmp/wisp/20261108/r354/tree-r1，HEAD=393430dc9）各摘各红：牙 A（摘 unlinkPinnedLeaf 复验块）⇒ 分子 3/3→0/3 断言红＋353r1 DAMAGE 回到 3/3（damage>0 判据逐字命中）；牙 B（摘 Remove 行，落点不同）⇒ 两枚负控齐红（present: size 0 body []）而主分子照绿＝如实报（该分子量复验留下的形，见牙 B 的是负控）；落 commit 字节零去牙。门禁（导出树）：gofumpt 空/vet rc0/build rc0/winsec 整包 ok 37.254s rc0/d22scan clean。⛔ 翻 AC 勾⛔ push；unix 半边=残余①机制不存在（同 211-354-closure 口径）。next=354-v1 非实现者裁 AC#2/AC#3，再 AC#4/AC#5 收口。
- [2026-11-08 23:2x +08] agent=354-v1（裁决腿：非实现者裁 AC#2＋AC#3；零产码零去牙入库；写面＝docs/evidence/s1/354-v1-verdict.md＋本条）did=两格判语全部出自导出树（git clone --shared，HEAD `b45ff3339`，/d/tmp/wisp/20261108/v354/tree-v1）本腿手跑：AC#2 成立（确定性三判据码内核过——t.Parallel/sleep/Eventually 仪器件 0 命中、逐字节断言在 :87-94/:100-101/:115/:159-163、轮内序号 :104-105、分子不借既有；-count=3 三轮 entered=3/3 leftStandingForeignName=3 逐字节同，与 r1 证据件逐字对上）；AC#3 成立（牙 A 摘 :277-283⇒分子 3/3→0/3 断言红＋353r1 DAMAGE 回 3/3；牙 B 摘 :284⇒两枚 no-top-up 控制齐红 present: size 0 body []＋主分子照绿＝票面允许——判据原文草案 SEC-C:30 逐字「摘删除那步⇒正控转红」，主分子量复验留下的形、摘 Remove 不改它，非偏差；落点两处不同现核；每发仪器内置 3 发非单发；两牙各还原、git status 空、整包复绿）。门禁导出树：gofumpt 空/vet rc0/build rc0/整包 ok 35.199s/d22scan clean。两点如实记：派单正控模式 -run 'Test354|LeftStanding' 字面匹配零测试（现跑留证）改用精确名；证据件「负控」与草案「正控」称呼不一致（同对象、无实质）。⇒ 编排者可翻 AC#2/AC#3 两勾；⛔ 本腿翻框⛔ push。next=AC#4/AC#5/AC#6 收口腿。
