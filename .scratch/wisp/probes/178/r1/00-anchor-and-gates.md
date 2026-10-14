# 178-r1 交件 —— 成对普查尺 G6 家族：判定从「一枚混数」改成四堆可分（写腿＝改仪器）

结论先行：
- **改前 rc=4 → 改后 rc=2**（G1b／G2 两枚＝别的账，具名报回、一字未修）；G6／G6neg 一对按四堆口径
  同批改形＋重登（`quiet`＋`want_n 0`），账实相符：**175-r2 那枚只写 CloseTask 的合法测试件不再顶红**（AC#3 ✓）。
- **AC#2 值钱那发两向都做了**：在生产码真造一枚【只开】之后——**改前尺沉默**（G6 主尺实测 1 恰等于旧
  基线 1 ⇒ ok）、**改后尺红**（G6 与 G6neg 双双 BAD，实测 1 > 基线 0）。四态逐字在 `logs/`。
- 合成树＝锚树 `1675e844`（`605838b0^{tree}`）＋既有生产件 `internal/tools/subagent_197.go` 的变异
  blob `f534915a` ⇒ mutant tree **`8ca0a52e`**；对象只落 `$TMPDIR` 的 scratch object dir，
  **工作树／索引／ref／仓对象库零写**（`git status --porcelain -- internal/ cmd/`＝空、`git diff` 同）。
- ⛔ 未 push；⛔ 未动 `internal/**`／`cmd/**` 产码、`thresholds.go`／SLO／golden、票面任何 Status／AC 框。

## 0. 锚与身份

- 派发锚 `93452d45`；**本腿起手现读＝`605838b0`（漂移逐报：以本腿四态读数钉 `605838b0` 为准）**。
  期间 HEAD 被别的腿推了两次（`aa539691`＝255-r8 三枚提交；`401191b7`＝18:2x 另一枚）——本腿全部读数
  显式传锚 `605838b0`，不受影响（改前那发日志头部逐字：`锚点 605838b07388cb64997fa1022bb1ea8e4394402f`）。
- 仪器：`.scratch/wisp/probes/154/gate-clauses.sh`——改前 sha256 `cda7c3f6…`（669 行）／改后
  `20e45dc8…`（735 行，+66）。GOTMPDIR=`C:\tmp\wisp178r1gotmp`、PATH 带 sherpa-onnx 目录（按派单）；
  本票不跑 go（纯 grep 仪器），`tasklist` 无重活可报。

## 1. 起点复认（改前 × 净树，`logs/gate-pre.txt`）：rc=4，四枚 BAD 逐字

```
# BAD  腿=G1b 声明=quiet 基线=0枚 实测=1枚 因=空/非空那一味与声明不符（票 161 AC#6①）
# BAD  腿=G2  声明=ring  基线=1枚 实测=2枚 因=新增命中（票 171 AC#7：实测 > 基线）
# BAD  腿=G6  声明=ring  基线=1枚 实测=0枚 因=空/非空那一味与声明不符（票 161 AC#6①）
# BAD  腿=G6neg 声明=ring 基线=1枚 实测=2枚 因=新增未成对（票 171 AC#2：实测 > 基线）
# 腿数＝14 声明与实测不符＝4
```
与编排者随单现量逐字同值（起点复认 ✓）。

## 2. 四枚 BAD 的归因表（第 0 步；逐个具名）

| 腿 | 改前读数 | 归因（现量） | 应归的票／账 |
|---|---|---|---|
| G1b | quiet/0，实测 1 | `internal/tools/subagent_197.go:591 return p.inner.Execute(ctx, req)`——票 197 落地的装饰器转发调用（文件由 `b9fa815b` 新增）；该腿 want_n 0 是 171-r3 时的树 | **别的账**：门自己那一票（票 154 族）／票面 U-4 那格已判「归口＝门自己那一票」；台账 `A943` 把它标为 161 AC#6① 那一味。⛔ 未动 |
| G2 | ring/1，实测 2 | 票 315（`bfc74ae9`）在 `admitTask` 补的 `rt.bridge.OpenTask(taskID)`（`cmd/wisp/run.go:986`）⇒ 命中行数 2（986 调用＋992 调用）；基线 1 是票 171-r3 在锚 `4f2a777b` 现量登记的 | **别的账**：票 171 AC#7 的口径（实测>基线）；重登记归口＝门自己那一票。⛔ 未动 |
| G6 | ring/1，实测 0 | 双成因同根：09-27 登记的 `ring 1` 本身就是一发【只关】（`run.go:563` 当时只关不开——正是票 178 判「关的是别处开的、不是缺陷」的那一堆）；票 315 让 run.go 开合成对后连包子数也归 0 | **本票 charter**（G6/G6neg 同一根因，派单允许一并改） |
| G6neg | ring/1，实测 2 | 本票立案形状：反向 2 枚全在【只关】——`internal/tools/ticket175r2_stamp_live_test.go`（5 处 CloseTask，175-r2＝`58db7400` 09-28 落）＋`internal/tools/pointer_183_cli_seam_test.go`（2 处，183-r2＝`afc766ee` 后落） | **本票** |

## 3. AC#1 四堆名册（读数件 `ac1-census.txt`，脚本 `ac1-census.sh` 可复算）

- 射程＝`git grep -nEw 'OpenTask|CloseTask' <锚 605838b0> -- internal/ cmd/`（含 `_test.go` 与 bridge.go 定义本体）；
  **工作树 vs 锚点去 sha 前缀后逐行相同**（票面那条字面命令＝同一读数）。
- 名册 39 行 ⇒ ① 开且关＝**5 枚文件**（run.go／bridge.go／bridge_scope…158_test.go／pointer_185…test.go／
  subagent_222_test.go）／② 只开（真缺陷）＝**0 枚**／③ 只关＝**2 枚**（ticket175r2、pointer_183）／
  ④ 注释里提到＝**3 枚文件、14 行**（`cmd/wisp/run.go:945/:947` 那类；`cmd/wisp/run.go:564` 时代的行号已位移）。
- ②＋③＝2＝改前 G6neg 混成的那一枚数——**票面那句「四堆混成一个枚数」现量成立**。

## 4. 修法要旨（instrument diff：只加打印＋一处口径＋同批重登）

- 新味 **`piles`**（`want` 第三味，与 `rev` 互斥；打错一味死在起跑线 rc=3）：pair 腿把名册拆成四堆逐堆打出
  （①开且关 BALANCED 名单／②只开＝原「成对名册基准」/③只关＝原「反向名册」/④注释里提到＝原 COMMENT-ONLY），
  **判据只吃 ② 只开堆的文件枚数**；③ 逐枚点名、不计进退码（票 178 判语：关的是别处开的）。
- **基线枚数语义（AC#3 要求写进仪器头）＝② 只开堆里的文件枚数（不是堆数、不是行数）**——文件头 §178 逐字、want_n 注释、
  腿自己那一格、聚合段注释四处同写。
- G6／G6neg 重登 `quiet`＋`want_n 0`（与形状同批；AC#4 要求的「改完之后 AC#2 那一发仍会红」逐字写在 G6 那一格的注释里）。
- G5／G7 两族腿（rev 味）口径原样：正＋反相加，**读数逐字节不变**（改后净树聚合行 G5 1/1、G5neg 8/8、G5pos 0/0、
  G7 1/1、G7pos 0/0、G7neg 2/2 逐枚同改前）。
- ⚠ 变异发在 G6/G6neg（quiet）上以**① 空/非空那一味**落 BAD（声明 quiet 却响了 1 枚）——book_count 两味并集、① 先响；
  红句里 `基线=0枚 实测=1枚` 两列同现（「多出来的那一枚」没被洗掉）。②「新增未成对」那味在 want=ring 的腿上才是首因。

## 5. 四态读数（逐字在 `logs/gate-pre.txt`、`logs/gate-post.txt`、`logs/gate-mutation-ac2.txt`）

| 态 | 尺 | 锚／tree | rc | G6 主尺 | G6-负一负 |
|---|---|---|---|---|---|
| 改前×净树 | `cda7c3f6` | `605838b0` | **4** | BAD 空/非空（实测 0） | BAD 新增未成对（2>1） |
| 改前×变异树 | `cda7c3f6` | `8ca0a52e` | 3 | **ok（实测 1 = 旧基线 1 ⇒ 沉默）** | BAD（3>1；但它净树本来就恒红⇒不可分辨） |
| 改后×净树 | `20e45dc8` | `605838b0` | **2**（G1b/G2 两枚他账） | **ok quiet 0/0**（四堆①1/②0/③0/④3） | **ok quiet 0/0**（①5/②0/③2/④14） |
| 改后×变异树 | `20e45dc8` | `8ca0a52e` | 4 | **BAD（② 只开 1 枚 > 基线 0）** | **BAD（② 只开 1 枚 > 基线 0；③ 只关 2 枚照旧点名）** |

## 6. AC 逐格自答（每格都答「这一发在未修仪器上响不响」）

- **AC#1 ✓**：四堆名册落 `probes/178/r1/`（票面写 `a1/`，本腿按派单用 `r1/`，⛔ 票面一字未改）。
- **AC#2 ✓**：未修尺上这一发**不可能给出可分辨的响**——G6 主尺在变异树上恰读 1=旧基线 1 ⇒ **沉默**，G6neg
  恒红（净树 2>1）⇒ 决定不了。修后尺：净树双双 quiet、变异树双双 BAD ⇒ 有牙、且只对【只开】有牙。
- **AC#3 ✓**：`ticket175r2_stamp_live_test.go` 落 ③ 只关堆、逐枚点名但不进判据 ⇒ 不再顶红；基线语义（②堆文件枚数）写进仪器头。
- **AC#4 ✓**：⛔ 没有抬基线；形状（四堆/piles 味）与重登同批同一 commit；「改完之后 AC#2 那一发仍会红」逐字在仪器表里＋本节；
  `thresholds.go`／SLO／golden 零字节（`git status` 逐字为证）。
- **AC#5 ✓**：与票 171 三枚洞的分工——171＝ⓐ合方词根写死（G5 的 CloseScope→CloseScope|Close()）、ⓑ「只多不少没牙」（升级成
  want_n 基线）、ⓒ run() 腿不带枚数；**178＝第四枚**：pair 腿的**四堆不分**（只开/只关/开且关/注释混成一枚数）。
  171 三格一枚未动、本票未并入其任何已勾格。

## 7. 变异程序（可复算；「不改工作树」的凭证）

```
# 变异内容＝锚树版 internal/tools/subagent_197.go ＋文末 6 行（注释＋真调用 b.OpenTask，无 CloseTask）
# 文件 sha256＝d3dfb3e1cc13866326adb02276cf278774261c4848dfc411aef58f8bf1a38586（原件在本目录 mut-subagent_197.go）
export GIT_OBJECT_DIRECTORY="C:/tmp/wisp178r1-obj"
export GIT_ALTERNATE_OBJECT_DIRECTORIES="D:/work/workspace/projects plans/Wisp/.git/objects"
blob=$(git hash-object -w .scratch/wisp/probes/178/r1/mut-subagent_197.go)          # ⇒ f534915af7a0c49a9435e94035732c12334ec93f
export GIT_INDEX_FILE="C:/tmp/wisp178r1-idx4"
git read-tree 1675e844e370b5045f3c085d141ecaa943189b8b                             # 锚 605838b0^{tree}
git update-index --add --cacheinfo 100644,$blob,internal/tools/subagent_197.go
git write-tree                                                                     # ⇒ 8ca0a52ecd94e1a980e3d7f0689051590e4f3128
# 还原证：git status --porcelain -- internal/ cmd/  ⇒ 前后皆空（0 行）；git diff --stat -- internal/ cmd/ ⇒ 空
```
- ⛔ 改动**没有**进任何 commit（合成树只是对象库里一枚不可达 tree；refs 一枚未动；工作树零改）。
  派单那句「仓外备份＋sha256＋还原证＋git status 空」针对就地改产码那一形——本形**没有就地改动**，
  以「internal/ cmd/ 起终皆空＋对象只落 $TMPDIR」为对应凭据（就地改法在共享树里有毒：318-r1／255-r8／114-v2 在飞，
  一枚临时 .go 会给别腿的编译/判据喂假读数）。

## 8. 随单发现＋依赖件（具名，未修＝不在本腿面，归口下一程）

- **⑤ 射程的文件清单来自工作树**：`pair()/run()` 的 `-- $@` 未加引号 ⇒ `internal/**/*.go` 先被 bash 做路径名展开
  （现量 **496 枚路径**）再喂 git grep。合成树里「新加一枚文件」因此不可见（本腿第一版变异＝新文件 `zz178r1…go`，
  mutant tree `c65cc174` 上全尺 `grep -c zz178r1`＝**0**；样本留在 `mut-zz178r1-open-only.go`，只建不删）。
  改它＝改所有腿在所有锚点的射程 ⇒ 须单独立票（已写进仪器头 ⑤，未动代码）。
- `probes/171/r3/r3-sample-deleted.sh` 把 `^want G6 ring`＋`want_n 1` 钉死 ⇒ 本改后它**rc=2**（"the leg moved, no verdict"
  ＝该门自己的诚实硬退）；修它须改该门，非本腿写面。
- `probes/161/r6/flip-declaration.sh` 末句把 G6/G6neg 列在「按设计就响」名单（打印句；该门今天因 G1b/G2 本来就 RED）。
- 工作树里别的腿的脏件（`probes/161/r6/logs/*`、`design/**`、`go.mod` 等）**本腿一枚未动**；本腿 diff＝`gate-clauses.sh` 一枚。

## 9. 具名没做

- G1b／G2 未修、未重登（不在本票判据；归口见 §2）。
- 未改 `r3-sample-deleted.sh`／`flip-declaration.sh`／任何别的探针或门（非本腿面）。
- 未 push；未 amend/reset/rebase/stash/checkout/clean；未用 `git add -A`。

---
[编排者追加 2026-10-14 19:2x，append-only] 按非实现者 `178-v1` §3 的裁定，本件 §4 里「G5/G7 读数逐字节不变」一句**收窄**为：「两族**口径、数值、判语、名册**逐字节不变；共享的 `# 未成对枚数＝` 行加了新尾注（G5 3 处＋G7 3 处）」。原句不抹。
