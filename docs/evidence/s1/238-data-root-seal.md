# 238-seal-r2 —— STOP / 复核（零实现）：派单前提仍然过期（owner A476/A1558 否决未被撤销）；「谁封根」在代码面已有主

- 腿：编排者 2026-11-13 11:37 派 `238-seal-r2`（写腿形派单，写面 `internal/proc`/`observe`/`secret`/`models`，交付证据件 `docs/evidence/s1/238-data-root-seal.md`）
- 起手锚点（现读）：HEAD `b5193c097`（2026-11-13 11:31:18 +0800）
- 结论：**STOP——零实现、零翻框、AC 原文一字未动**。与 wave108 `238-proc`（`31668c2d7`）同判定，另补三笔新读数：① 撤销口令「Q-72 重开」今巡＝未命中；② `internal/proc` 今日**零 `MkdirAll`**（连测试夹具都没有）⇒ 派单"两形"里的甲形在今天的进程拓扑下**没有封的对象**；③ AC#4 ③（`parseACELines` 双副本归属）今巡**零归属**，是三条欠账里唯一还悬着的。

## 1. 派单前提过期（同 238-proc，本腿全部现读复测）

- **A476**（2026-09-30 16:5x，台账 `pending-and-issues.md:9904` 起）：Q-72 由 owner 当场否决**整件事的目的**——原话逐字「锁个屁，人家都能登录我电脑了，还锁个蛋啊，怎么想的你」（台账 `:9906`）。首启封根两支（甲/乙）都不落地。
- **A1558**（2026-11-05 20:1x，台账 `:31282`）：「238 STOP（owner-vetoed）……勿再派」。
- **A1655**（2026-11-09，台账 `:32998`）：编排者把 `238-proc` 放回 wave108 未引「Q-72 重开」⇒ 过期派单，编排者已自记"我的错"；台账 `:33000` 逐字「**238 永不派**，除非 owner 亲口『Q-72 重开』」。
- **今巡（本腿现跑）**：`Q-72 重开` 在台账中出现 6 处（`:9912`/`:32996`/`:32998`/`:33000`/`:34174`/`:35694`），全部是规则行/复盘行，**零一处是 owner 亲口重开** ⇒ 撤销口令未命中，否决仍有效。
- **A1710/A1712**（2026-11-11）：编排者又派过一枚只读 `238-a1`，随后自记"派单前没读 Progress log"是错（A1712）；该件（`probes/238/a1/census.md`，commit `e436b3a37`）结论与 A476 同向：132-r1 已把"创建即封"落地，AC#1..3 的消费者今天不存在。
- 与 `238-proc` 的不同点：本腿是写腿形派单，故把每一格"如果做、代价是什么"在 §3 答满；五条门禁读数是本腿现跑（§6）。

## 2. 现量复测（本腿现跑读数，非转述；锚点 `b5193c097`）

| 核 | 读数 |
|---|---|
| `internal/proc` 对 `internal/winsec` 导入 | **0**。`grep -rn "CarlosShao/wisp/internal/winsec" --include="*.go" internal/proc/` ⇒ **10 行命中全部在注释里**（`envfork.go:112/113/116/138/145/148/156/289` 共 8 行＋`envfork_test.go:116/148` 共 2 行），代码零命中 |
| `internal/proc` 建目录动作 | **零 `MkdirAll`/`Mkdir`（生产＋测试都是零）**。仅 `envfork_test.go:247/271` 两枚测试夹具 `os.WriteFile(..., 0o644)`；生产代码零建文件动作 |
| `internal/buildinfo` | 零 `MkdirAll`、零 winsec 命中（grep rc=1） |
| `internal/proc → observe` 正向边 | 真实存在：`boot_windows.go:19`、`externalsampler_windows.go:10`、`treemetrics_windows.go:12`（另有 4 枚测试件），共 7 处 import |
| `internal/observe` 日志根创建 | `logging.go:16` `import winsec`；`:99` `winsec.PrivateDirAll(cfg.Dir, 0o700)`——**创建即封**，对已存在的根也封（那半是修旧树的修复路径）；`:297` 滚动档 `os.OpenFile(..., 0o600)`；`:316` `sealLogFile`（`:223` 的 var 缝，持 `winsec.SealFile`） |
| 票面所指 logging.go:72 旧行 | 已不存在（今天 `:72` 是 LogConfig.RollDays 字段注释）；「裸 MkdirAll(0o755)」与「0o644」两个旧形状均已成史 |
| 132-r1 落地 | commit `b22225513`；票 132 证据件 4 枚真在 `docs/evidence/s1/`（132-v1-verdict / 132-r1-impl / 132-r1-impl-part2 / 132-r2-impl） |
| 常驻腿封根链 | `cmd/wisp/resident_windows.go:173` `installLogSink(rt.Layout.DataDir)` → `logsink.go:157` 定义 → `observe.InitLog` → `logging.go:99` PrivateDirAll(0o700) ＋ 滚动档 SealFile |
| 装配点 4 枚 | `resident_windows.go:173` / `run.go:246` / `models.go:339` / `secret.go:226`（grep 现跑） |
| 既有 STOP 件 | `docs/evidence/s1/238-proc-r1.md` 真在（38 行）；其 STOP commit `31668c2d7` 真在（`git log -1 -- <path>` 返回该 commit） |
| 238 探针档 | `probes/238/c1/census.md`（56,056B）与 `probes/238/a1/census.md` 均真在 |

## 3. AC 逐格回应（答满，不落地）

### AC#1 先定"谁负责封根"——两形摆开、各自代价、拍板（没改任何代码，"先别改代码"照做）

**两形＋已在码中的第三形状，代价逐条：**

| 形 | 代价明细 |
|---|---|
| **甲：常驻腿（internal/proc）自己调封条** | ① 新增依赖边 `proc→winsec`（今日导入数 0）：这条边**不成环**（graph 本就按 `risk→observe→secret→winsec` 定向、`winsec/resolve.go:35-40` 逐字记录"the reverse edge closes a cycle and go build refuses"），所以撞的不是编译；真正的冲突是台账 `:8788` 在案裁定「装配根 `cmd/wisp` 是唯一的接缝」——让 proc 自己接 winsec 缝＝多一个接缝，属契约级；② **甲形在今天的拓扑下是空转**：`internal/proc` 零 MkdirAll ⇒ `proc.Boot` 跑完的时刻数据根**还不存在**（根由 `observe.InitLog` 的 `PrivateDirAll` 连父带子一次造出）⇒ "boot 时封根"没有对象可封；要它有对象，得先给 proc 添建目录动作＝造第二个创建者，更糟；③ observe 已在的封根照旧跑 ⇒ 两个封根者并存；④ A476 直接否决目的 |
| **乙：装配根（cmd/wisp）在起任何腿之前把根封一次** | ① 写面越界：cmd/wisp 不在本腿写面（另一条腿在跑），本腿不能落；② `PrivateDirAll` 对**已存在**的根也下封条（`winsec.go:249-250` 逐字「An existing path *is* sealed, because that is the repair path for a tree that predates this package」）⇒ "起腿前封一次"在升级装/老机子上**顺手做了 AC#4①（既有根补封）**＝A476 点名不做的"顺手加固"；③ 乙形今天同样有对象问题：装配根最早的可封时点（main 早期）数据根多半尚不存在，真正把它造出来并封掉的就是后面 `installLogSink → observe.InitLog` 那一跳——乙形只是把同一件事换个地方写第二遍；④ cmd/wisp 每波只许一枚写腿的纪律 |
| **第三形状（已在码中，非本票派单所给）：由创建者即封** | `observe.InitLog` 在 `logging.go:99` 用 `PrivateDirAll(cfg.Dir, 0o700)` 把"造根"与"封根"合成一次调用：目录带窄描述符出生，之后滚动档经 `SealFile`（`:316`）逐枚收口。这就是 132-r1（`b22225513`）落地的形状 |

**拍板：不拍甲、也不拍乙——"谁封根"这个问题在代码面已有主：由创建者即封（whoever creates the root seals it at creation）。**

理由（承重层，按铁律 2/3 本腿拍并写明）：

1. **A476 否决的是目的本身**（"首启封根"这件事），不是在甲/乙之间选形。目的被否 ⇒ 甲、乙、以及 A1710 裁定②的注入形全部失去消费者；AC#1 的"由编排者裁后再派落地腿"被 A476 的"不再排落地腿"取代。
2. **"由创建者即封"不碰模型夹**：`winsec.go:246-248` 逐字「Ancestors above path are never touched」⇒ 不做全树走查 ⇒ 避开票 95 那句「decisions like that do not get made by an unrelated hardening sweep」（`no_seal_ruling_windows_test.go`）——正是 A476 的裁定里点名要保的形状。
3. **甲形空转论证（读码硬读数）**：proc 零 MkdirAll ⇒ boot 时刻无对象可封。乙形换处重写同一件事，多一条边、多一个封根者，买不到新性质。
4. **主流 harness 对齐（读码推论层，不承重）**：主流常驻 harness（Docker daemon 对 `/var/lib/docker` 一类）的形状就是"创建者创建时落权限"，而不是"每次启动把既有根重封一遍"。乙形的"起腿前再封一次"在主流里没有对应物。本腿未做 web 取证，这句按惯例不承重；承重的是上面 1-3 与 A476。
5. **撤销口令**：本拍板的撤销口令＝owner 亲口「Q-72 重开」。重开后按 `probes/238/a1/census.md` 的乙形材料＋本件 §2 现量重新派单。

### AC#2 一把能判"这棵树的根是谁造的"的常驻尺——不造，消费者不存在

- 这把尺的存在理由＝守 AC#1 裁定产生的**新**封条。今天没有新封条：根的创建已经走 `PrivateDirAll`（logging.go:99），"根是谁造的"已由 132-r1 的实现＋132 的证据件回答。
- 每个新东西必须有消费者 ⇒ **不造**。
- **若重开后要造**（设计要点存档，防止重犯 132-c3 顶正过的错）：判据**不许**写"至少一枚带非 (I) 的 ACE"——`PrivateDirAll` 的 Windows 侧 `sealDir→propagatePrivate`（`winsec_windows.go:588/592/620`）会在**事后**给早先的子档补显式窄描述符（132-c3 §6.4：8 份日志含 CI 里 `bak-0-1` 显式、`bak-1-2` 仅继承的分裂读数），"至少一枚"型判据今天就已绿＝测空气。能区分时序的形状＝**逐枚**型：封条覆盖范围内的**每一枚**都必须带显式（非继承）窄描述符；X 变异（注释掉 `logging.go:99` 的 PrivateDirAll）下，靠继承出生的档全部退回 `(I)` ⇒ 逐枚判据必红。判据注释里逐字写明"每枚 `*.jsonl` 的收窄属票 132"。

### AC#3 POSIX 半边——按 A476/A1710 裁定：无事可做，理由写进判据注释

- `sealDir` POSIX 侧＝单档 chmod、不递归（`winsec_other.go:165-170`），该段的存在理由逐字写明"the hole being closed does not exist on this platform"。
- 派单点名的 `logging.go:244 0o644` 已不存在：今天的滚动档是 `logging.go:297` `os.OpenFile(..., 0o600)`，POSIX 侧残余洞已被 132-r1 顺手收掉。
- **递归方案＋代价**：递归 chmod 就是 `propagatePrivate` 的 POSIX 对应物——A476 之下属"顺手加固"，不许做；且 POSIX 这台本来不验（winsec_other.go 的原注释即判据注释）。故选"把 POSIX 侧这台不验写进判据注释"，且该注释已由 132-r1 落在 `logging.go:95-98`（0o700 注释）与 `winsec_other.go:165-170`。没有用"容器里跑过了"代替——根本没跑，也不需要跑。

### AC#4 三条"刻意不做"的欠账归属

1. **① 既有数据根补封（升级装/老机子）→ 归票 132**。`PrivateDirAll` 的"existing path is sealed"语义（`winsec.go:249-250`）就是修复路径，且每次 `InitLog` 都对日志子树跑一遍（logging.go:99）、数据根由 `memory.Open`（132-c3 记 `open.go:176`）覆盖——132-r1 已把这一格按自己的票面收口。A476 点名的"不许顺手加固"正是拒绝在这之外再造一枚独立补封形状。
2. **② `--data-dir`/`WISP_DATA_DIR` 对外口子 → 归台账 `A111④`**（票面原文具名的缺口）。A476 裁定③"失败降级"之后它是逃生口；本票零改动口径下不新增归属，仍挂在 A111④。
3. **③ `parseACELines` 双副本漂移 → 今巡零归属，是三条里唯一悬着的**。现跑读数：`.scratch/wisp/issues/` 树内 `parseACELines` 仅票 238 自身两处提及（`:29` 刀5、`:48` ③）⇒ 无票号、无台账号。按 AC#4 完成判据"三行各自有票号或台账号"，③ 需要**编排者/owner 落一枚归属**（新票号或台账号）——这是"改台账/立票"动作，不在写腿权限内，本腿把读数给满：双副本（`internal/models/acl_sid_121_test.go:50` 与 `internal/winsec/acl_windows_test.go:86`）今天都在树；漂移与否的逐字节比对本腿未跑（零实现腿的诚实边界，留给重开后的落地腿先跑再动）。

## 4. 写面与纪律

- 基线（起手现跑）：`git status --porcelain cmd/wisp internal/ball` ＝ **空**；写面四包＋票面＋证据目录起手零未提交改动（证据目录里已有 8 枚**别人家**的 untracked 文件，本腿零接触）。
- 零实现：零 `.go` 改动、零测试新增、AC 框未翻、AC 原文一字未动。
- 零越界：零 frontend/design 读与写；零 PLAN/specs/thresholds/golden/allowlist 触碰；零 internal/config、scripts、cmd/wisp 写；未动 `build/mut_rvw.go`。
- 禁令：未用 `git add -A`/`git add .`/`--amend`/`reset`/`rebase`/`stash`/`checkout .`/`clean`；commit 用显式 pathspec；不 push。
- 交付形先例：`238-proc` 的 STOP 交付＝票面追加＋证据件同一枚 commit（`31668c2d7`），本腿同形。
- 本腿足迹（仅两枚）：本证据件＋票 238 Progress log 追加一行（append-only）。

## 6. 门禁读数（本腿现跑，2026-11-13 11:44-11:50 +0800；零实现 ⇒ 后读数＝基线读数）

| # | 门禁 | 读数 |
|---|---|---|
| 1 | `go build ./cmd/wisp/` | rc=0 |
| 2 | `go vet ./cmd/wisp/` | rc=0 |
| 3 | `go test -count=1 ./internal/proc/ ./internal/observe/ ./internal/secret/ ./internal/models/` | exit 0：proc 4.110s ok / observe 6.203s ok / secret 0.441s ok / models 13.986s ok |
| 4 | `gofumpt -l internal cmd tools scripts frontend` | rc=0，零输出（gofumpt 不在 PATH，走 `D:\work\base\gopath\bin\gofumpt.exe`） |
| 5 | `scripts/d22scan.sh`（Git bash 跑） | rc=0，干净：positive control PASS=45 FAIL=0；scan `clean - no D22 ban violations`（bans #1-5 internal/=403、cmd/=77、ban #6 frontend/=123、ban #7 tools/=52、ban #10 risk/=21、ban #8 design=39/frontend=123/internal=1004/cmd=276）。**走了 scripts/d22scan.sh，没有用 `go run ./tools/d22scan`** |

## 7. 交付与重开条件

- 交付：本证据件＋票 238 Progress 追加，一枚 commit（显式两枚 pathspec），不 push；推送归编排者。
- 重开条件：仅 owner 亲口「Q-72 重开」⇒ 编排者重派落地腿（届时写面含 cmd/wisp seam 与 winsec 非递归导出名按 Q-72 原案）；取用材料＝`probes/238/a1/census.md`（乙形三文件方案）＋本件 §2 现量＋§3 AC#2 的"逐枚型"判据设计要点＋AC#4③ 漂移检查（未跑，先跑再动）。
- 若 owner 不重开：AC#4③ 的归属格（parseACELines 双副本）仍悬着——它是三条欠账里唯一无票号/台账号的，须由编排者/owner 落账；本腿读数已给满。
