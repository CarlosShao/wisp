# 票 322 — 非实现者终裁验收（腿 322-v1，只读）：六格判语与去牙/反形/还原全果

> **这份文件是什么**：票 322 六枚 AC 格子的**验收判语**（`SPEC-12 §4.3` #1/#3、`AGENTS.md` §0.3 双角色）。
> 实现者＝`322-r1`（死腿，字节代提 `8baec646`）＋`322-r1b`；本腿⛔ 任何一方，⛔ 翻任何 AC 框，
> ⛔ 动 `internal/winsec/**` 一字节。框翻不翻、票走不走，全归编排者。
> **凭据规则**：`322-r1b` 的证据件与读数件里每一句都按未验证处理；下面每一枚数都是本腿自己现跑的，
> 台面＝共享工作树只读＋`git archive` 仓外导出树物理突变（⛔ `-overlay`、⛔ 仓内 worktree）。

## 0. 起手（现量，写进每件读数抬头）

- 起树锚 `HEAD=71002d00`（`git rev-parse HEAD`，`2026-10-16 10:17:10 +0800`）。
- 票框尺现跑：`grep -c '^[[:space:]]*- \[ \]'` ＝ **6 未勾**；`'- \[x\]'` ＝ **0 已勾**。与派单预期相符，无第三腿动票。
- 台面漂移一笔具名：本腿突变台件开跑时 dev 上 `HEAD` 已推进到 `e396301d`（`220-r2` 在飞提交）。
  尺＝`git rev-parse 71002d00:internal/winsec` ＝ `e396301d:internal/winsec` ＝ `8ddcb8cb` **同一棵树对象**
  ⇒ 导出树里的 winsec 字节与本腿起树锚逐字节一致，读数有效。
- 本机整包起手一发（`go test ./internal/winsec/ -count=1 -v`，`GOTMPDIR=C:\Users\swq\AppData\Local\Temp\gotmp322v1`
  原生反斜杠形、**C 卷**，前置 sherpa-onnx 于 PATH）：**`RUN=110 PASS_TOP=61 FAIL_TOP=0 SKIP_TOP=0` rc=0**
  （`logs/00-host-full-v.log`）。派单点名的两枚存量红（`TestGateJudgesThePrivateSetByResolvedSID`、
  `TestSealNarrowsAndNamesThePrincipalItRemovedBySID`）在本腿这一发起手即绿；这一族随卷/残留变色的归属
  **本腿⛔ 重判、⛔ 裁生死、⛔ 计入本票红**（票 310/311 的形状，归编排者）。
- 口径具名（两把尺⛔ 混用）：`^=== RUN` **含子测试行**、`^--- PASS` 只数顶层，所以容器全量那批发 `RUN=71/topPASS=43`
  差 28 枚——本腿现跑差集证明 28 枚**全部**是 `父/子` 形式的子测试 RUN 行，无失踪、无台面故障。

## 1. AC#0 三形代价表 —— 判：**成立**（读验，⛔ 重做）

普查件在 `.scratch/wisp/probes/322/a1/`（`index.md`＋`logs.md`＋三枚探针源副本），本腿逐格读了 §1 主表、§2 Windows
答、§3 顶红名册、§4 调用者承诺表、§5 未测 11 格。对照落地产码（`winsec.go:141-157`、`winsec_other.go:14-39`、
`winsec_windows.go:322-405`，本腿逐字读过）：**落的形＝表里"甲（POSIX 拆常量 `O_NOFOLLOW`）＋表 §2.6 具名的
Windows 唯一可用形状（`FILE_FLAG_OPEN_REPARSE_POINT`+`OPEN_EXISTING` 探针即拒）"**——表说"甲/乙两形 Windows 都
等于没修、reparse 探针给判据不给写路径、`CREATE_ALWAYS` 变体是第二种损伤"，产码注释逐条照说（`winsec_windows.go:341-348`
明写 CREATE_ALWAYS banned），表与落地是同一件事，无暗中换形。`O_EXCL` 稳态格由 `openleaf_overwrite_322_test.go`
的 exclusive 子项钉住（本腿跑绿，见 §4）。

## 2. AC#1 截断窗口关死（POSIX 叶子那一支）—— 判：**成立**

成对读数＝本腿自己现跑，同发同屏（同一行打印 damaged 与 liveness 两数），全部仓外导出树：

```
base1b  (git archive 8baec646^＋预票桩＋本票测试件)  rc=1  71/41/2/0
        attempts=3000 flips=2242 damagedWitnesses=67 leafLegRefusals=0
fixed1b (git archive HEAD)                           rc=0  71/43/0/0
        attempts=3000 flips=1828 damagedWitnesses=0  leafLegRefusals=50
fullfixed-v1（同一 fixed 树整包全量）                 rc=0  71/43/0/0
        attempts=3000 flips=1446 damagedWitnesses=0  leafLegRefusals=34
```

修前分子 67>0 ⇒ **坏世界真被造出来**；修后 0 **且** liveness 非零（50 与 34 两发独立），⛔ "恒真第一形"。
红句逐字（base1b，`logs/base1b.body.log`）：`AC#1 RED: 67 of 3000 foreign witnesses lost their own bytes or
mtime, so the open in privateFile followed a leaf link swapped into the caller's name and truncated a file no
caller named (first: foreign-witness-0007-322r1.bin: content is "", want its own bytes).`
那一发在未修码上响＝响（race 红＋flag 件红两形＋悬空叶在树外造对象，`logs/base1b.body.log` 三句全在）。
祖先那一支仍开＝票面与产码注释都具名写着（`winsec.go:128-133`），本格⛔ 裁它——票面判据只管"叶子被换成链接"。

## 3. AC#2 覆盖写语义不破 —— 判：**成立**

调用点名册＝本腿现跑 `grep -rn "privateFile(" internal/ --include=*.go`（生产面仅 `winsec.go:96/:106` 两枚入口）＋
`grep -rn "winsec.PrivateFile(" internal/ cmd/`（非测试）：**`:87` 一支五枚生产调用者**——`internal/agent/spill.go:166`、
`internal/config/migrate.go:93`、`internal/secret/migrate.go:169`、`internal/secret/migrate.go:189`、
`internal/secret/store.go:79`；**`:97` 独占一支两枚**——`spill.go:343`、`internal/memory/artifacts.go:75`。
判定"允许覆盖"仍在、未被具名改成拒绝：本机 `-run TestAC2Overwrite… -v` **rc=0，6 枚子项全绿**（4 枚依赖覆盖的
调用者各一枚＋entry-stays-private＋exclusive-still-ErrExist，`logs/09-ac2-host.log`）；容器 fixed 档同绿。
两形点名两清：四枚依赖覆盖写（store/migrate×2/spill），第五枚 `secret/migrate.go:169` 被 `os.Stat` 挡在覆盖之前、
依赖的是"名字还不存在⇒不拒"——那一支由 Windows 直调用例的 `absent` 行钉着，本腿 w2 反形里它红过（§5）。
**这一发在未修码上响不响＝不响，且这是正确的**：base1b 全量里 `TestAC2Overwrite…` PASS（覆盖写本来没坏，坏的是截断），
这些钉是本票**新加**的防"未来把允许覆盖改成拒绝"的牙——m2b/w2 反形证明它们咬得住（§5）。

## 4. AC#3 Windows 那一半 —— 判：**以票面第二形成立〔降格为仅本机可量＋残留如实报告〕；任何"Windows 窗口已关"读法 = 不成立**

本腿自己现跑（本机 `go1.27 windows/amd64`，同一台件四发，`logs/mut-win-v1-summary.txt`）：

```
w0-baseline  rc=0 3/0/0  attempts=400 flips=247 guardRefusals=10 returnedNil=126 witnessTruncated=true witnessSize=0
w1-restore   rc=0 3/0/0  attempts=400 flips=293 guardRefusals=13 returnedNil=132 witnessTruncated=true witnessSize=35
w2-restore   rc=0 3/0/0  attempts=400 flips=181 guardRefusals=7  returnedNil=199 witnessTruncated=true witnessSize=35
```

判据（本腿写的，不是抄的）：这一格票面本来就给了两条出路——"给出 Windows 侧等价防护并量到"或"降成〔仅本机可量〕
并登记 CI/容器覆盖不到"。读数支持第二条、不支持第一条：① 守卫**有牙且在产线路由上**（摘掉必红，见 §5 的 w1；
路由每发 `guardRefusals>0` 且 `flips>0`，仪器行自检通过）；② 守卫**关不死窗口**——已修码上本腿四发全部
`witnessTruncated=true`，其中 `witnessSize=35` 三发＝外来见证不止被截空、还被换上了 raced payload
（`"raced private payload 322r1 windows"` 恰 35 字节）＝探针答完话到 `os.OpenFile` 之间那扇窗在本机**今天真能穿过**。
`witnessSize=0` 与 `=35` 之别按形状登记：`=35` 意味着字节写到了外来对象之后没被同一对象复验拦下＝`winsec_windows.go`
的 `sealHandle` 无 same-object 闸门＝**票 120 的 Windows 半格**，属具名外射程（`openleaf_guard_322_windows_test.go:23-25`
与 `winsec_windows.go:371-373` 都这么写，本腿读验一致）。
降格所需登记在位：测试件头 `:27-30`〔仅本机可量〕＋`322-a1/index.md` §5 未测格 3（CI runner 符号链接权限档从未量）
＋产码两处 residual 具名（`winsec_windows.go:359-369`）。用例把截断计数写成报告不写成闸门
（`openleaf_guard_322_windows_test.go`，本腿 `grep -n witnessTruncated` 现定位＝**252 行**，⛔ 信派单给的 249-253，实测在 `:249-253` 的 Logf 语句里），
与票面对"丙形不许写成窗口已关"的禁令一致。**本腿⛔ 拿 POSIX 绿冒充这里任何结论**（上面每一枚数都来自 Windows 档）。
那一发在未修码上响＝响：w1（守卫恒答"不是链接"＝未修形的物理复现）三枚具名链接红＋仪器红，见 §5。

## 5. 牙（首要攻击点＝判据仪器自己）—— 判：**三形齐，仪器不是恒真尺；缺世界＝零枚**

全部突变在 `git archive` 仓外导出树里**物理改字节**（⛔ `-overlay`）；每形改前留 `.bak`、记 sha256、还原后复跑。

| 形 | POSIX（容器 golang:1.27） | Windows（本机） |
|---|---|---|
| 改前必红 | base1b：race 红 `damaged=67`＋flag 件红（`logs/base1b.body.log`） | w1＝未修形：三枚链接红＋仪器红（逐字四句在 `logs/mut-win-v1-summary.txt`） |
| 摘牙必响 | m1b 摘 `|oNoFollowPlatform`：rc=1 `topFAIL=1`，红句 `AC#1 RED: 84 of 3000 foreign witnesses lost their own bytes or mtime …` | w1 `return nil`：rc=1 `topFAIL=2`，`AC#3 RED file-symlink/directory-symlink/junction`＋`AC#3 RED instrument: … guard-leg refusals=0 and racer flips=274` |
| 反形也红 | m2b 逢名就拒：rc=1 `topFAIL=9`——race `damaged=0` 但仪器行红（`liveness=0`），AC#2 八条全红（逐字含 `MUTANT blanket refusal on every name`），另七枚既有测试连坐具名 | w2 逢名就拒：rc=1 `topFAIL=3`——`guardRefusals=0 otherRefusals=400` 路由仪器红（按守卫原句词判，反形不被放行）、真文件/不存在名两枚红、三枚链接红改判"词不匹配"、AC#2 八条红 |
| 还原复绿 | winsec.go sha 还原＝`41ba9a9b…`、winsec_other.go＝`19f3c211…` 逐字相同；还原树复跑 m1-restore1b rc=0 三枚全绿 `damaged=0 liveness=56 flips=1870` | winsec_windows.go sha 还原＝`5cdef881…` 两形都逐字相同；w1-restore／w2-restore 都 rc=0 全绿 |

恒真三形逐一堵住的证据点名：`damaged=0` 本腿从没单独当凭据（每次都同行核 liveness）；"有任何 error 就算防住"
在 m2b/w2 里被**词判**（守卫原句 `whyGuardRefusal322r1` 不匹配即算别的腿答话⇒红）。
单边尺一枚具名（不是缺世界）：`TestAC322LeafFlag…`（`openleaf_flag_322_other_test.go`）不经 `privateFile`，
m1b 里它绿＝分工正确，接线凭据归 race 那一发；但它在 base1b 里红（未修码上响），本腿验过它不是恒绿尺。

## 6. AC#4 门禁＋消费者不推红 —— 判：**已跑的六道全绿；⛔ 一枚格子缺读数（cmd/wisp）**

六道本腿逐枚自跑、各落 rc（`logs/00-gate-rcs.txt`＋`01-gate-rcs-b.txt`）：
`go build ./...` **rc=0**；`go vet ./internal/winsec/` **rc=0**；`GOOS=linux go vet` **rc=0**；`GOOS=darwin go vet` **rc=0**；
`sh scripts/d22scan.sh` **rc=0**（clean；ban #8 `internal/` 枚数在件内）；`sh scripts/check-path-length-budget.sh --with-self-test` **rc=0**。
格式化那把：本机 PATH 无 `gofumpt`（rc=1），按派单改**绝对路径** `/d/work/base/gopath/bin/gofumpt.exe` 跑——
`-l internal/winsec/` **rc=0／列出 0 枚**（`logs/07-gofumpt.log`）；⛔ 拿 `gofmt -l` 的 0 冒充。全仓形（winsec+cmd）列出 2 枚
（`cmd/wisp/models.go`、`cmd/wisp/panel_transport_35r2_test.go`）＝**在飞写腿 220 的写面，⛔ 计入本票、⛔ 判它**，本腿没读它们的内容。
消费者维：本机 `go test ./internal/risk ./internal/secret ./internal/config ./internal/agent ./internal/memory` **全 ok rc=0**
（`logs/10-consumers-host.txt`）；Linux 容器整包 `fullfixed-v1` **71/43/0/0 rc=0**；Windows 本机整包起手 **110/61/0/0 rc=0**。
**缺的这一发具名交回**：`go test ./cmd/wisp` 整包本腿零读数（派单：本波唯一写手在动它＋票 314 死亡区同刻只许一枚）——
派单正文写 `220-r1`、总括口径写 `220-r2`，盘上最新提交署名 `220-r2`；**由编排者向该腿取这一发**，本腿⛔ 代跑、⛔ 判"没问题"。

## 7. AC#5 注释说实话 —— 判：**成立**（点名那句已改对；一枚邻位过期句另报，⛔ 本腿动手）

票面点名的原句（`winsec_other.go:82`，未修码上逐字为 "That shape is booked separately as ticket 322 and is
unfixed at this anchor."——本腿 `git show 8baec646^` 现取）已消失：全仓 `grep "booked separately"` 0 命中（先跑正控
确认尺在 `8baec646^` 档能命中该行）。现句（`winsec_other.go:103-113`）写"叶子已由 `oNoFollowPlatform` 关上、祖先仍开、
⛔ 任何人把这两句读成窗口已关"——与本腿 AC#1/AC#2 的裁定**相符**（叶子确关、祖先确开、覆盖确在）。
`winsec.go:118-140` 三段与 `winsec_windows.go:322-373` 各自点名残留（POSIX 祖先、Windows 探针不锁答案、Windows 祖先），
与本腿 §4 实测一致。"窗口已关"式全称声明的尺本腿现跑：CJK 字样在 `internal/winsec/**` **0 命中**；英文侧仅存的
"the window is closed" 两处都在**否定句**里（`openleaf_322_other_test.go:8` "Nothing here claims…"、
`winsec.go:139` "nothing here may be read as…"），肯定式只有"LEAF leg is closed on POSIX"（`winsec.go:125`）——
该狭义声明为真，由本腿 §2 成对读数背书；没有任何一处把**整扇窗**说成已关。
**邻位上报（不推翻本格，归 120 家族的测试面）**：`sealhandle_privateroute_120_other_test.go:232/:237` 仍写
"stated as unfixed at this anchor in winsec_other.go:78-82" 与 "(ticket 322, unfixed at this anchor; winsec_other.go:78-82)"——
前半已被本票终态越过（本腿判词：POSIX 叶子已关），行号引用也腐烂（那句现在 `:103-113`）。
**应改成的那句**（写在这里供实现面使用，本腿⛔ 碰 `internal/winsec/**`）：把 "unfixed at this anchor" 改为
"the leaf leg is closed on this platform since ticket 322's leaf flag; the branch above can now only be reached
through the still-open ancestor leg or a regression, and the cited sentence lives at winsec_other.go:103-113"。

## 8. 本腿自己的台面故障与偏差（自抓，不等谁来抓）

1. 第一版 POSIX 台件 `mut-posix-v1.sh` 的 `container()` 把 `-run` 过滤串拼进了单引号 `sh -c` 内层，
   过滤器里的 `|` 被容器 sh 当**管道**吃掉 ⇒ base-v1/fixed-v1/m1-v1/m2-v1 四发 rc=127、零判据读数——
   **全部按台面故障处理、⛔ 计入任何格的判据**（同脚本的整包全量发 `fullfixed-v1` 无过滤串、幸而有效，已单列 §2/§6）。
   修正版 `mut-posix-v1b.sh` 重跑了 base/fixed/m1/m2 全四发，§2/§5 的数全部出自 v1b。故障件按"只建不删"留在盘上入库。
2. v1b 因 `$1` 未传参把 `runfilter` 落成空串 ⇒ 四发都是**整包全量**色（信息更宽，无碍判据；口径已在各发写明）。
3. 导出树取自推进中的 dev（`e396301d`），winsec 字节与起树锚同一棵树对象（§0 的 `8ddcb8cb` 尺），有效。
4. `git add`/`git commit` 每一笔逐枚显式 pathspec（c1 17 枚、c2 21 枚、c3 12 枚），⛔ push、⛔ amend/reset/stash。

## 9. 六格总判语（翻框权⛔ 在本腿）

| 格 | 判语 | 本腿凭据 |
|---|---|---|
| AC#0 | **成立**（表与落地同一件事；⛔ 重做表） | §1 |
| AC#1 | **成立**（POSIX 叶子那一支；成对读数自跑、liveness 同行；祖先那一支票面本来就写"仍开"，归下一票） | §2 §5 |
| AC#2 | **成立**（允许覆盖仍在；5+2 枚调用者本腿现 grep 点名；钉有牙＝反形双平台红） | §3 §5 |
| AC#3 | **以票面第二形成立＝〔仅本机可量〕＋残留如实报告＋CI 覆盖不到已登记**；"Windows 防护＝窗口已关"读法**不成立**（本腿四发 `witnessTruncated=true` 亲量） | §4 §5 |
| AC#4 | **部分成立**：六道门禁＋本机消费者＋容器/本机整包全绿；**`cmd/wisp` 档缺这一发读数**，具名交编排者向 220 写腿取 | §6 |
| AC#5 | **成立**（点名那句已改对且与本腿实测相符）；120 测试面一枚过期句按 §7 给定句子另行处置 | §7 |

## 10. 落点索引

读数件＝`.scratch/wisp/probes/322/v1/`（脚本三枚＋`logs/` 全量原始正文；提交 `755fb974`／`8b8fa0be`／`f918d81e` 与本件同链）。
导出树（只建不删）＝`/d/tmp/probe322v1/{base.A9XC3y,fixed.Ti4jeF,m1-teethout.uFeBIY,m2-blanket.mGNGul,win.kH0tLw,base1b.OBdSIE,fixed1b.8l9hKI,m1b-teethout.pu8XeM,m2b-blanket.BZAkIw}`。
本件＝`docs/evidence/s1/322-v1-acceptance.md`。票尾 bullet＝追加一条。⛔ push。
