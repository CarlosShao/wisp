# 收件腿 r1 —「已交件但没翻框」批次的逐票裁决表（collection-r1）

- 时刻：2026-11-09 09:0x–10:0x +08｜腿名 `收件r1`（编排者派，**非实现者**＝裁决者角色；READ/ADJUDICATION）
- 锚：派单时主树 HEAD `9dbd468e3`；开工后主树漂到 `85800170a`（唯一新增＝354-cls2 收口文书，产码零字节）；**全部复跑读数出自共享克隆导出树 `git clone --shared` → `D:\tmp\wisp\20261108\163recv\repo`，checkout `cda38fd52`**（起手 `git status --porcelain`＝0 行）。导出树与主树 HEAD 之间，本批 9 票的票面/证据件 `git diff --stat`＝空 ⇒ 读数不受漂移影响。
- 环境令照办：`export PATH=$PWD/build:$PWD/third_party/sherpa-onnx:$PATH`；`unset GOTMPDIR TMPDIR TEMP TMP`；未设 `CGO_ENABLED`；gofumpt 用绝对路径 `/d/work/base/gopath/bin/gofumpt.exe`。⚠ 导出树缺未跟踪的 `build/`／`third_party/sherpa-onnx/` DLL（`git ls-files` 两处均 0 枚），首次跑 `./cmd/wisp` 撞 `0xc0000135`（STATUS_DLL_NOT_FOUND）⇒ 手补 3 枚 DLL 入导出树（onnxruntime／sherpa-onnx-c-api／sherpa-onnx-cxx-api，只建不删，不触主树）后复跑。
- ★本腿亲证编队级活缺陷仍在：导出树 `sh .scratch/wisp/probes/154/gate-clauses.sh` ⇒ 聚合退码＝**10**（与派单一致；`.go` 未引号展开超 Windows 命令行上限，读数件 `.scratch/wisp/probes/154/r20261108/00-argv-overflow.md`）⇒ **今天任何票面凡引 gate-clauses 读数的一律取不到真值**，本表逐票如实记。

## 逐票判语（收件三把尺：提交真在／件真在／自己复跑）

| 票 | 交件 | 三把尺 | 判 |
|---|---|---|---|
| **163 AC#1b** | `8ae13b4a0`（11-06，裁定表＋台件 3 枚，+240） | 全过：R1 四行在册（`PLAN.md:2574/2575`＋`SPEC-07:74/75`）；R2 台件复跑双 `Lookup=false`＋对照 `task.output=true`（名册 21→**22**＝`task.list` 新增，只多不少）；R3 `shell.session`=0、`shell.exec` 非测试 6 枚全非实现 | **翻〔成立〕** |
| **177 AC#0** | `aee58549c`（11-05，131 行独立普查件） | 全过：四形状导出树逐发自跑 PASS（A 0.21s／B 0.00s／C 34.78s／D 2.34s）；164 两枚续读判据 PASS；三问全答、独立起手 | **翻〔成立〕** |
| **195** | `e6791e6bb`（11-05 死腿遗产代提，4 文件 +909/−50） | 大体过：config 包整包 ok（十枚 AC 命名用例逐名 PASS）；调用链 7 处现跑在；AC#4 名册 8→9（+1 纯注释具名）；AC#6 零触碰 | **AC#1/2/3/4/6 翻；AC#5 具名不翻**（gate-clauses 仪器坏，那一条取不到真值；其余三把绿）；AC#0 照 10-02 裁定永缺 |
| **71 AC#1/AC#2** | `08c4c51e2`（11-05，5 文件 +715） | 全过：balldebug 包 8/8 PASS；进程级四发真退出码 rc=**1/1/2/0**（真红×2＋拒调低×1＋阴性对照×1） | **双翻〔成立〕** |
| **138** | 在盘：`138-cancel-kills-cascade.md`／`138-cancel-kills-children-r1.md`／`138-cancel-kills-children-r1-accept-r1.md`／`138-cascade-ruling.md` 四件证据 | 本腿未复跑（时间盒，优先 163/177/195/71） | **报状态**：交件在盘、票面有 log；收框待下一腿逐格复跑 |
| **140** | 在盘：11-06 两条 log（READ-MEASURE 四包真跑＋AC#2(b) 实现腿 `140-ac2-fixes.md`） | 本腿未复跑 | **报状态**：交件在盘、log 已记账；AC#1/AC#2 翻框待独立复跑（变异三发读数在表，可照抄复跑） |
| **230** | 在盘：`230-cells-closure.md`（11-08 腿）＋更早 5 件 | 本腿未复跑；11-08 腿自己零框翻、AC#5 如实报「cmd/wisp 整包未到终态（两发两死点＝真窗热键活体挂起）」 | **报状态**：AC#3 复核乙成立／AC#4 机主 A1549 已裁走 B——两格可翻凭据在件，留编排者或下一收件腿；AC#5 照实未达 |
| **346** | 在盘：`346-a1-census.md`（10-19，编排者当日已收、零翻框）＋`346-r1-impl.md` | 本腿未复跑 | **报状态**：两趟交件在盘；票面 7 格的翻框凭据（(i)–(iv) 四问答案）在普查件内，待非实现者逐格裁 |
| **178** | `-done` 改名 | 双名核过：`git ls-tree HEAD` 里旧名（无 `-done` 后缀）0 枚、只 `-done` 版在场；外加 `71-*` 同名族三枚 `-done` 在场不混淆（slug 本身含 "test"） | **核毕，状态正常** |

## 本腿改动面（本 commit 的全部）

- 翻框＋log 行：163（AC#1b）、177（AC#0）、195（AC#1/2/3/4/6＋AC#5 具名不翻）、71（AC#1/AC#2）——每票 Progress log 追加一行 `agent=收件r1 did=复跑读数 verdict=…`；AC 原文零字节未动（只翻勾）。
- 本表（新建）。
- 零 `.go`、零 `scripts/**`、零 `internal/winsec/**`、零 `frontend/**`、零 `forks/**`、零 thresholds/golden、零 push。LF 核：编辑后四票＋本表 CR=0。

## 本腿没测什么（如实）

- 138/140/230/346 四票只做了「交件真在／件真在」两把尺＋票面 log 现读，**没有复跑各自格点名的读数**（时间盒：派单要求至少收 4 票，优先级 163/177/195/71 已全收）。
- 177 证据件自己声明的边界照抄继承：300s 生产默认等待与丙形红集合是〔读码〕非〔现跑〕；31s 是生产同一条路（`risk.confirm_timeout_sec` 可配下限）。
- 195：主树有一枚在飞未跟踪件 `docs/evidence/s1/195-r1-direction-verdict.md`（另一腿的 195-r1 方向判语），本腿未读未引未动——AC#1..6 的判语只基于 `e6791e6bb` 已提交的实现与我自己导出树的复跑。
- gate-clauses.sh 的 argv 溢出缺陷本腿只做了亲证（退码 10），未修（`scripts/**` 另一腿独占）。

## 撤销口令

「撤 collection-r1」＝revert 本笔即可：diff 只有 4 枚票面（勾位＋log 行）＋本表，产码零字节。
