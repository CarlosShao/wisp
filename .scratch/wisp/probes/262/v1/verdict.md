# 票 262 门禁仪器 — 独立验收裁决表（腿号 `262-v1`）

**本件是什么**：非实现者验收腿 `262-v1` 对票 262 落地的「全仓跟踪路径长度门禁」的**逐格独立取证**。
⛔ **本腿不勾任何 AC 框**（翻框权在编排者）；⛔ 零产码／零脚本／零 CI／零票面／零台账改动。

- 验收腿：`262-v1`（写面唯一＝`.scratch/wisp/probes/262/v1/**`）
- 起手 HEAD：`dea66c57`（`git log --oneline -1` 复认，`dev`）
- 起手锚时刻：`2026-10-07 14:09:10 +08`（`date` 现取）
- 票面行数（起手）：`83`（`wc -l`）
- 被验收仪器：`scripts/check-path-length-budget.sh`（**577 行**，只读）
- CI 接线：`.github/workflows/ci.yml`（941 行，只按行区间读）
- 落地腿读数件：`.scratch/wisp/probes/262/r1/census.md`（530 行，按节 `sed -n` 取）
- 仓外工作拷贝（本腿 AC#3/#4 突变面）：`D:/tmp/262v1/`（⛔ 不入库、不删）

---

## 0. 锚读数（本腿自己现跑，取数时刻 2026-10-07 14:09:16→14:09:20 +08）

命令逐字：

```
sh scripts/check-path-length-budget.sh --with-self-test
```

原始输出＝`logs/gate-head.txt`（19 行）。关键行：

- `rc=0`／`VERDICT GREEN - every over-budget tracked path is rostered by name with a reason, and the roster equals the tree`
- `denominator: tracked paths=7405  over-budget=57  covered by roster=57  not in roster=0`
- `longest=180 chars relative (.scratch/wisp/issues/252-the-two-spellings-of-one-path-split-inside-the-allowlist-check-short-name-stays-short-long-name-folds-long-so-in-root-writes-read-as-out-of-bounds-r2-l2.md)`
- `worst full path on the self-hosted runner=224 chars (hat budget 165, wall open interval (206,217])`
- `bands: over the hat=57  of which in the 122..180 middle=57  past the old debt line(180)=0  in the wall interval=0  roster entries=57`
- `unit=characters (probe: length of the 5-byte sample "a U+25FF b" = 3; LANG=C.UTF-8 LC_ALL=C.UTF-8)`
- `hat: rule 9 name cap=100 + issues dir prefix=21 -> relative hat=121 ; worst checkout prefix=44 -> full-path budget=165`
- `unit cross-check: the byte view and the character view select the SAME 57 paths on this tree, recomputed just now in both units`
- 内建正控三发：`control 1/3 ok - bench baseline green`／`control 2/3 ok - the planted over-budget tracked path is rejected and named verbatim:` + 一行 `RED - over budget and NOT in the roster: .scratch/wisp/issues/seeded-…over-the-hat.md  (relative 125 chars, name 104 chars, full path on the self-hosted runner 169 chars, budget 165)`／`control 3/3 ok - green again once the planted path is gone` ⇒ `positive control PASSED`
- bench 落点：`positive control bench=/tmp/tmp.Eeu8TV0TzH (kept on disk; this project never deletes temp artifacts)`

**与编排者 13:5x 那发的对照**：分母四项（7405／57／57／0）与本腿 14:09 现量**逐字一致** ⇒ 本腿不推翻该数；
但本腿把它当**此刻读数**而非常量记账（见 §AC#2 的"分母会动"条款）。

---

## 2. AC#0（前置：具名解冻＋§现量 1/2/5/6 四把尺由本腿现跑复认）——**判：成立（其中"墙"那一半只能到读码级，已归 AC#5）**

- **具名 `A##` 解冻记录**：`docs/reports/pending-and-issues.md A585`（票面 `:35` 具名由编排者落在该条；
  脚本头 `:25-33` 与 `ci.yml:138-141` 的注释都逐字把它绑成"只加一步"的唯一授权面）。本腿未动台账，只核"这条记录在不在"。
- **§现量 1 的本机可测三件**（14:11:3x 现跑，`logs/ac0-local.txt`）：
  `git config --get core.longpaths` ⇒ **rc=1（未设）**；
  `LongPathsEnabled` 注册表 ⇒ **0**（尺＝`D:/tmp/262v1/read-longpaths.ps1` 走 `powershell -File`，⛔ 未内联 `$`）；
  本机 git ⇒ **2.52.0.windows.1**；
  字串尺＝`printf '%s' 'E:\work\base\actions-runner\_work\wisp\wisp' | wc -c` ⇒ **43**（ASCII ⇒ 字符＝字节），
  并且 `E:/work/base/actions-runner/_work/wisp/wisp` **本机存在**（`ls -d` 命中）⇒ 43/44 那两枚数不是抄来的。
  ⛔ **本腿没有跑过一次 self-hosted checkout**："runner 上 checkout 到底在哪一长度断掉"这半只能到**读码／日志级**，
  与 §现量 2 同格 ⇒ **归 AC#5（待取数）**，本腿不假装量过。
- **§现量 2（墙在开区间 (206,217]）**：本腿**没有**实测凭据（要真跑一次 checkout）。
  仪器侧本腿只核到"这三枚常量在脚本里是**打印用、不是判定用**"：`DEBT_CAP=180 / WALL_LOW=206 / WALL_HIGH=217`
  （脚本 `:166-169`），且 `bands:` 那一行确实打它们 ⇒ **判语：这一半待取数，⛔ 不许拿推断当实测**。
- **§现量 5＋6 的分母（本腿自己的尺，14:11:27，⛔ 未引用票面 §现量）**：
  `git -c core.quotepath=false ls-files | awk '{L=length($0); n=split($0,c,"/"); nl=length(c[n]); ...}'` ⇒
  **`<=121:7350`／`122-150:38`／`151-180:19`／`>180:（空＝0）`／`last-name>100:57`／total 7407**；
  中间带单尺 `git ls-files | awk 'length($0)>121 && length($0)<=180' | wc -l` ⇒ **57**；
  按门禁定义（`相对>121 或 末段名>100`）单算 ⇒ **57**（两把尺同数 ⇒ 票面 §现量 6"两条线都不点名它们" Today 由门禁点名）；
  最长一枚 **180**＝`.scratch/wisp/issues/252-the-two-spellings-…-r2-l2.md`（`sort -rn | head -3` 现取）；
  非 ASCII 跟踪路径 **28** 枚（`grep -c '[^ -~]'`，与脚本注释里那枚 28 对上）。
  ⚠ **本腿自己让分母动过**：票面/编排者那发是 **7405**，本腿 14:11 读到 **7407** ⇒ 差的 2 枚就是本腿第 1 笔 commit
  落进去的 `verdict.md`＋`logs/gate-head.txt`。⇒ **"分母是常量"这句话在本票今天就被证伪**，任何引用都得带时刻。
- 未勾框枚数尺：`grep -c '^[[:space:]]*- \[ \]' .scratch/wisp/issues/262-*.md` ⇒ **8**（与 README:82-84 那句"8 枚全未勾"一致；本腿一枚不翻）。

**判语**：`A585` 在场＋四把尺里**三把（1 的本机半／5／6）本腿现跑复认成立**；
**"检出目录宽度对 checkout 的真实影响"与"墙的精确位置"两半本腿给不了**（要真跑一次 self-hosted checkout ⇒ 归 AC#5）。
⇒ **本格成立，但成立的范围是"帽这一档＋本机配置面"，⛔ 不含墙那一档。**

## 3. AC#1（脚本存在且被真调用）——**判：成立（接线在场、调用点 2 枚具名；"真跑过"不在本格）**

调用点枚数＝**2**（具名 `file:line`，尺＝`git grep -n "check-path-length-budget" HEAD`＋`git show HEAD:.github/workflows/ci.yml | grep -n`）：

1. **CI（形丁）**＝`.github/workflows/ci.yml:136` 步名 `"Tracked path-length budget (ticket 262)"`／
   **`:166` 命令行** `run: sh scripts/check-path-length-budget.sh --with-self-test`；所在 job＝`lint`，`:66 runs-on: ubuntu-latest`。
2. **本机 pre-push（形丙）**＝约定式调用，具名写在两处：`scripts/check-path-length-budget.sh:25-33`（"编排者——唯一持推送权者——每次 push 前跑同一条命令"）
   与 `.scratch/wisp/issues/README.md:79`（10-07 补句给的定式"开完新票、提第一笔之前跑 `sh scripts/check-path-length-budget.sh`"）。
   ⛔ **没有装 git hook**，且理由可核：本仓无 `.githooks/`、`core.hooksPath` 未设、`A585` 只授权"加一步 CI"。

⚠ **本腿给 AC#1 划的界（与编排者一致，并补一条结构性风险）**：
CI 里有这一步 **≠ 它跑过**。`ci.yml:166` 这一步**没带 `if:`**，而同一个 `lint` job 里它前面还有 **3 步**
（`:74` D22 positive control／`:83` D22 self-test／`:108`+`:134` seven-ban scan），
前面任一步红 ⇒ 这一步在真实 run 里就是 `[skipped]`＝**那次 run 它没执行**（这正是本仓 09 月踩过的那枚坑，`ci.yml:234`/`:294`/`:331`/`:387` 别的步骤都带了 `if: ${{ !cancelled() }}`，⛔ 这一步没有）。
⛔ 本腿**不把它算作缺陷**：补 `if:` 要动既有步骤面，超出 `A585` 的"只加一步"授权 ⇒ 属**编排者裁**（要么走一枚新的具名解冻，要么就承认"这一步的覆盖率＝前面三步都绿的 run"）。
"真跑过"这一格＝**AC#5**。

## 4. AC#2（HEAD 上跑＝绿，且打印分母读数）——**判：成立**

本腿两发独立读数（⛔ 都没引用票面 §现量）：

- **本仓工作树 14:09:16→14:09:20**：`sh scripts/check-path-length-budget.sh --with-self-test` ⇒ **rc=0／VERDICT GREEN**，
  分母四项全打：**`tracked paths=7405  over-budget=57  covered by roster=57  not in roster=0`**，
  最长枚**长度＋全名**都打（180 chars ＋ `252-the-two-spellings-…-r2-l2.md` 全名），
  中间带枚数打（`in the 122..180 middle=57`），名册枚数打（`roster entries=57`），
  外加最差全路径 224／帽预算 165／单位判定 `unit=characters`／字节-字符两单位同集合的 cross-check。原始件＝`logs/gate-head.txt`（2,528 字节）。
- **仓外 clone（HEAD `1d59543e`，含本腿第 1 笔）14:14:10**：同一条尺（不带 self-test）⇒ **rc=0／GREEN／`tracked paths=7407 over-budget=57 covered 57 not-in-roster 0`**
  （`logs/ac3-readings.txt` 同段的 `clone-baseline.txt` 摘要）⇒ **"绿"不依赖工作树脏**。
- **AC 要求的四枚分母读数是否"打印"而非"要人自己数"**：是——本腿把 `run_scan` 的 END 段读码对上了打印行
  （脚本 **`:399-405`**：`denominator:`／`longest=`／`worst full path=`／`bands:` 四条 printf，本腿逐行 `grep -n` 点过）。
- **对照编排者 13:5x 那发**：7405／57／57／0 与本腿 14:09 逐字一致 ⇒ 不推翻；
  本腿 14:11 读到 7407 是**本腿自己的提交**造成的 ⇒ 记为"分母会随任何人的提交动，引用必带时刻"。

## 5. AC#3（正控：种一枚违反规矩 9 的跟踪路径 ⇒ 门必须红 ⇒ 拆掉 ⇒ 绿）——**判：成立**

⛔ **本腿全程在仓外 clone 里种**（`D:/tmp/262v1/clone`，`git clone --depth 1 --branch dev` 于 14:14:05，252 MB），
**本仓工作树自始至终没有出现过一枚越帽跟踪路径**（自证见 §10）。命令与读数逐字＝`logs/ac3-readings.txt`＋`logs/ac3-A.txt`／`ac3-B.txt`／`ac3-C.txt`。

- **A｜14:14:32 种"票形"越帽路径**（名 130／相对 151，正是规矩 9 的射程）：
  `PAD=$(printf '%110s' '' | tr ' ' 'x'); TI="262v1-ac3-seeded-$PAD.md"; git add -- ".scratch/wisp/issues/$TI"` ⇒ **rc=1**
  ```
  RED - over budget and NOT in the roster: .scratch/wisp/issues/262v1-ac3-seeded-xxxx…xxxx.md  (relative 151 chars, name 130 chars, full path on the self-hosted runner 195 chars, budget 165)
  RED - count guard: the roster holds 57 entries, the tree has 58 over-budget tracked paths
  ```
  ⇒ **点名那枚路径＋它的长度（相对／末段名／最差全路径）三条都在**，AC#3 的"红句要点名"这一硬条件成立。
- **B｜14:14:41 再种一枚"浅目录长名"（名 108、相对 108＝在相对帽之下、在名帽之上）** ⇒ **rc=1**，同一形状点名它
  （`relative 108 chars, name 108 chars, full path 152 chars, budget 165`）
  ⇒ **这一发是 AC#3 的一个额外分支**：只钉"相对 >121"的门会**放掉**这枚（115 字符的名字在仓根只有 115 的相对长度），
  脚本里那枚 repo-wide"末段名 >100"的第二支帽把它拦住了 ⇒ 第二枚帽不是装饰（敏感性的另一半见 §6 的 M2b）。
- **C｜14:14:42 拆掉两枚**（`git rm --cached --quiet -- …`）⇒ **rc=0／GREEN／57＝roster 57** ⇒ "拆掉恢复绿"成立。
  ⚠ 拆＝**从跟踪集里拆掉**（与脚本内建 self-test 同形，它用的也是 `git rm --cached`）；两枚 `.md` 文件按规矩 8"只建不删"**留在 clone 里未跟踪**（`git status` 显示 `??`），不影响门的读数（门只读 `git ls-files`）。
- **仪器自带的第三发**（14:09:16，本仓）：`--with-self-test` 在一次性 git 仓里种一枚相对 125 ⇒ `control 2/3 ok` 并逐字打红句 ⇒ 拆 ⇒ `control 3/3 ok`。

**关于编排者 10-07 那发红（60 vs 57）——本腿明确不同意"它只是旁证"这句**（具名理由）：
那发红的是 `not in roster` 那一支（守卫 A）＋计数守卫 C，**触发物是真实越帽票名、真实工作树、真实调用**，
形状上已经把"门会红并且点名"证了一遍；本腿 A/B/C 三发只是**把它复认成"我自己种的、我知道长度的"** 那一形。
⇒ **本腿把两发都记为凭据**：A/B/C＝**直接凭据**（自种自拆），编排者 13:5x 那发＝**同一判据的第二次独立读数**（不是旁证）。
差别只在来源，不在判据；AC#3 原文（票面 `:42`）也只要求"种一枚 ⇒ 红 ⇒ 拆 ⇒ 绿，两发读数在交件里"，⛔ 没要求来源必须是验收腿自己造的假名。

## 6. AC#4（名册逐枚具名理由＋删理由门必须变行为）＋两发恒真检查——**判：成立**

**名册 57 枚的逐枚静态审**（尺＝`logs/ac4-roster-audit.txt`，14:16 现跑 python，对象＝HEAD 那枚脚本）：

- 解析出 `R["path"] = "reason"` 形如数据的行：**57 枚**；**重复键 0**／**空理由 0**；
  **57 枚键全部仍是 HEAD 的跟踪路径**（守卫 B 的反向核：名册没被改名落下）；
- 相对长度范围 **122–180**、末段名范围 **101–159** ⇒ 名册与"中间带 57 枚"是同一集合（§2 那把独立 awk 尺同数 57）；
  **没有一枚**是"只靠末段名帽豁免"（`relative<=121` 的名册枚数＝0）⇒ 名册没有藏额外豁免；
- 理由的**事实性**本腿也核了，不只是"非空"：
  写 `still open` 的枚里**没有一枚**文件名带 `-done`（0 处矛盾），写 `closed, name doubles as the -done anti-double-claim key` 的枚**全部**带 `-done`（0 处矛盾）；
  分类＝**54 枚工单＋3 枚 dispatch note**（`hist: dispatch note …` ⇒ 规矩 9 只管票名，dispatch 名走同一帽但理由不同句）。
  ⚠ **一句如实定性**：54 枚共用**同一句类目理由**（"filed before README rule 9"），不是 57 句各不相同的故事。
  票面 AC#4 的原文是"每一枚都有具名存在理由"——**本腿判它成立**，因为①逐枚点名了路径（名册是数据不是名单）②理由逐枚落在数据位③类目句本身可核（上面那两条一致性尺就是它的核）；
  但若编排者要的是"逐枚各写一句只属于它的故事"，**那是比票面更严的一档**，本腿不擅自替票面加严——写在这里供裁。

**删掉理由 ⇒ 行为必须变化（守卫 D）**：14:15:50，在 clone 里把 `252-…` 那枚的理由置空
（`sed '/issues\/252-the-two-spellings/ s/= ".*"$/= ""/'`，⛔ 只改 `D:/tmp/262v1/mut/` 下的副本）⇒
```
rc=1  RED - roster entry carries an EMPTY reason; the reason is data, not a comment: .scratch/wisp/issues/252-the-two-spellings-…
```
⇒ **AC#4 主条件成立**（理由不是注释）。
**补一发（比 AC#4 更狠的一形）**：14:15:51 把那一整行名册**删掉** ⇒ `rc=1`／`covered by roster=56 not in roster=1`／
`RED - over budget and NOT in the roster: …252…`／`count guard: the roster holds 56 entries, the tree has 57` ⇒ 名册每一枚都连着牙。

### 6.1 恒真检查①——阈值与第二枚帽这两形敏不敏感（14:15:48／14:15:49／14:15:52／14:15:53）

- **把 `HAT_NAME_LIMIT=100` 改成 `999`**（`mut/m1-hat999.sh`）⇒ **不是照旧绿，而是 rc=1**：
  `over-budget=0 covered=0` 但 **`RED - roster entry that is not an over-budget tracked path any more` × 57 行**＋`count guard: the roster holds 57 entries, the tree has 0` ⇒ **VERDICT RED**。
  ⇒ 抬帽这一形**在真实树上会被守卫 B＋C 当场抓住**（票面禁区里那一支"写成 999"不是安静的）。
- **删掉 `|| namelen > namecap` 那一支**（`mut/m2-no-second-hat.sh`）⇒ 干净树上 **rc=0 照旧绿**（两枚帽今天同集合＝57，删一支不改变现量），
  **但**把 §5-B 那枚"仓根 108 字符名"再跟踪进来时：
  `shipped-copy ⇒ rc=1（点名 262v1-ac3-root-…）` vs `m2 ⇒ rc=0／over-budget=57／GREEN`。
  ⇒ **这一支确实是被读的**：它对"浅目录长名"这一形负责，删掉就静默放掉。**这是本票里唯一一枚"改了门还是绿"的形**，
  ⚠ 所以它的凭据**只能是带种子的对照**（本腿已给），⛔ 不许有人拿"删了它门还绿"当"它没用"。
- 差量自证：每个 mutant 与入库脚本 `diff` 只有 **1–2 行**（中和断言那一枚是 **4 行＝两处改动**；`logs/mutants-log.txt` 的 `changed-line-count=` 段逐枚在案），
  对照组＝入库脚本的逐字节副本 `shipped-copy.sh`（`changed-line-count=0`，同路径同样本跑法跑出的 14:15:47 基线 rc=0/GREEN）。

### 6.2 恒真检查②——`--self-test` 本身是不是恒真（14:15:54／14:15:56／14:15:59／14:16:30）

- **M5**＝把阈值改成 999、断言**原封不动**，跑 `--self-test` ⇒ **rc=1**：
  `CONTROL FAILED - planting a tracked path 125 characters long should be RED (rc=1), got rc=0. The gate may be stuck green`
  ＋它把 bench 的现场也吐出来（`gate.sh: denominator: tracked paths=2  over-budget=0`）＋
  `positive control RED - the scan below is not run, because a gate that cannot go red does not get to print green`。
  ⇒ **`--self-test` 不是恒真**：一枚"再也不会红"的门会被它自己拒掉。
- **M6**＝同样抬帽，**再把 control-2 的断言中和**（`if [ "$rc" -ne 1 ]` ⇒ `if false`），跑 `--self-test` ⇒ **rc=0／positive control PASSED**
  ⇒ 定位到"牙"就在那枚断言行上（**机制证明**，⛔ 不是"门没问题"的结论）。
  **但 M6 补跑 `--with-self-test` ⇒ rc=1**：扫描段仍被 `RED - count guard: the roster holds 57 entries, the tree has 0 over-budget tracked paths` 打死。
  ⇒ **纵容和面双层**：中和自检断言救不了那枚门——名册-树计数守卫（C）在 CI 真正跑的那条命令行（`--with-self-test`）上兜住。
- **M7**＝入库脚本原样 `--self-test` ⇒ **rc=0**，三发 `control 1/3·2/3·3/3 ok` 齐全（对照用，说明 M5 的红不是环境噪声）。

⇒ **AC#4 判成立**，且**两发恒真检查都过**：唯一"改了看不出来"的形是 §6.1 第二支帽，本腿已给它配了带种子的对照。

## 7. AC#5（CI 侧真出颜色）——**判：〔待取数〕，本腿给不了**

- 结构性事实（本腿读码级，`:136/:166`＋§3 那条无 `if:` 的跳过风险）已在 §3 具名。
- 本腿**只读地**试过 `gh`（`gh 2.96.0`，`gh auth status` 报 keyring 失败但 `Active account: true`）：
  - 14:12 一发成功：`gh run list -R CarlosShao/wisp --limit 5 --json …` ⇒ 最新 5 发里 `ci` workflow 两发（**`37545246395`＝2026-10-07T00:42:25Z，sha `cc315261`，conclusion=failure**；`37406757402`＝10-06，failure），另有 3 发 `slo-fresh` success。⇒ **"改名之后的树上这一步跑没跑过"这问题，这几发的 sha 都早于今天 13:5x 的改名**（`f67af953` 未推送），编排者那句"从未在改名之后的树上跑过"与本腿读数不冲突。
  - 14:13 起 `…/actions/runs/<id>/jobs`、`gh run view --json jobs`、`--log-failed` 连试 **6 发全部 `EOF`**（网络/代理侧，不是权限侧；原始错误行在 `logs/ac5-gh-attempts.txt`）⇒ **本腿没能读到任何一次 run 里这一步的 step 级结论**。
- **要什么形状的凭据才算（具名）**：某次 `ci` workflow run 的 `lint` job 里，
  step 名 `"Tracked path-length budget (ticket 262)"` 出现且结论＝`success` 或 `failure`（**⛔ `[skipped]`／`cancelled` 不算**），
  日志里可指到本腿 §4 那组打印行之一（`denominator: tracked paths=…`／`VERDICT GREEN`／`RED - over budget…`）；
  取数方式＝`gh run view <id> -R CarlosShao/wisp --json jobs`（或 `--log` **先落到仓外文件**再 `grep`，⛔ 不入本仓、不入上下文）。
- ⚠ **本格还有一枚"文字与形状不同形"的待裁项**：票面 AC#5 原文（`:44`）写的是"该步在**本机 self-hosted 那一侧**跑过"，
  而 `A585`＋形丁裁定把这一步放在 `lint`/`ubuntu-latest`（`ci.yml:66`），并**明写它照不到 self-hosted 的宽**（`ci.yml:150-158`、脚本 `:36-51`）。
  ⇒ 按现形状，AC#5 那一枚"**self-hosted 侧**跑过"**永远不会被满足**；本腿**不自行改写判据**，
  请编排者二选一：(a) 把 AC#5 读成"CI 侧真出颜色（任意 job）＋托管侧不误响＋具名说明它照不到哪一侧"，或 (b) 另立解冻真把一枚尺放到 self-hosted 侧。
  本腿的判语按 (a) 的形状取数，但**在框被翻之前把这一歧义挂在这里**。

## 8. AC#6（卫生四门读数不扩大）——**判：成立（本腿跑的三门全 0 扩大；`./cmd/wisp/` 那一枚归编排者）**

⛔ 本腿**没有**对 `./cmd/wisp/` 跑任何 `go test`／`go vet`（那条腿独占整包窗口）——**`go vet ./cmd/wisp/` 这一格由编排者补跑销账**，本腿不代跑、也不判它"没做所以不成立"。

| 门 | 本腿现跑（命令逐字＋时刻） | 读数 | 扩大否 |
|---|---|---|---|
| `tools/d22scan`（含票 212 ban #9） | `sh scripts/d22scan.sh` @ 14:16:58→（脏工作树） | **rc=0**／`clean - no D22 ban violations`；scopes：`bans #1-5 internal/=228`、`cmd/=38`、`ban #6 frontend/=85`、`ban #7 internal/tools/=23`、`ban #8 design/=39 frontend/=85 internal/=514 cmd/=104` | 红名集合＝**空** ⇒ 无扩大 |
| 同上，干净 HEAD | `sh scripts/d22scan.sh`（在 `D:/tmp/262v1/clone`＝HEAD `1d59543e`）@ 14:19:15 | **rc=0**／同一句 `clean`；`ban #8 design/=30`（clone 里 design/ 的枚数与脏树不同＝别人的未提交增删，**不是门禁变化**） | 红名集合＝**空** |
| `gofumpt -l`（⛔ 不含 `cmd/wisp`） | 自证 `gofumpt -version`＝**v0.12.0 (go1.27.1)**；`gofumpt -l internal tools/d22scan tools/mockllm` @ 14:18:08（脏树）／14:18:26（干净 clone） | 脏树＝**4 枚**：`internal/agent/approval/pending_read.go`、`internal/agent/tools.go`、`internal/risk/provenance.go`、`internal/tools/bridge.go`；**干净 HEAD＝0 枚（输出 0 字节）** | **不扩大**：那 4 枚是**别人在飞的未提交改动**（`git status --porcelain` 755 条脏项、含 `D design/assets/*`），HEAD 里那 4 枚本是干净的 ⇒ 归口＝在飞腿，⛔ 不算票 262 的账 |
| `go vet ./internal/...` | 后台跑，止于 **14:17:39** | **rc=0**，输出 56 字节（只有 rc 与时刻两行） | 无扩大 |
| `go vet ./cmd/wisp/` | ⛔ 本腿未跑（指令硬禁＋整包窗口被另一条腿独占） | — | **归编排者补跑** |
| 票 212 ban #9 单独一发 | 尺＝`tools/d22scan/main.go:50/:148-149/:873`（phantom-citation，随 `go run . -root <repo>` 一起跑）⇒ 本腿的 `sh scripts/d22scan.sh` 两发 rc=0 已含它 | 违规 0 枚 | 无扩大 |

⚠ 口径差异具名：CI 那一步是 `gofumpt -l . tools/d22scan tools/mockllm`（`.` **含 `cmd/wisp`**），本腿只跑非 `cmd/wisp` 子集 ⇒ **本腿的 0 枚是子集意义上的 0 枚**，整包那一枚仍欠编排者。

**读数件索引（同目录 `logs/`，⛔ 无一枚叫 `.out`——根 `.gitignore:8` 是全仓 `*.out`）**：
`gate-head.txt`（§0/§4 锚发 19 行原文）· `ac0-local.txt`（§2 本机四把尺＋分母带）·
`ac3-A.txt`／`ac3-B.txt`／`ac3-C.txt`／`ac3-readings.txt`（§5 三发逐字）· `clone-baseline.txt`（§4 第二发）·
`ac4-roster-audit-and-ac6-summary.txt`（§6 名册逐枚审＋§8 四门摘要）· `mutants-log.txt`（§6 七种 mutant 全现场，33,215 字节）·
`m6-with-self-test.txt`（§6.2 那一发）· `ac5-gh-attempts.txt`（§7 只读取数尝试，含 6 发 EOF 原文）·
`d22scan.txt`／`d22scan-clone.txt`（§8 两发全量）· `gofumpt-noncmdwisp.txt`／`ac6-gofumpt-clean-head.txt`（§8 gofumpt 两口径；
其中 `gofumpt-clone-head.txt` 是**真·0 字节空输出**（＝0 枚文件会被重排）；本腿起初断言"git 不存空文件"是**错判**，
`git ls-tree HEAD` 指着它＝`blob e69de29bb2d1d6434b8b29ae775ad8c2e48c5391`（空 blob 已入库），`ac6-gofumpt-clean-head.txt` 只是同一读数的带出处版本）·
`govet-internal.txt`（§8）· `ac7.txt`（§9 十四枚逐枚）· `restore-proof.txt`（§10）。

## 9. AC#7（旧名可追：11＋3＝14 枚）——**判：成立（1 枚的"现名"已过期一跳，具名为缺陷；"记在票面"本腿判它满足 AC#7）**

尺＝`logs/ac7-pairs.txt`（14:19:4x 现跑）：`git show --name-status -M --format= fd269de1｜46079fcc｜f67af953`
⇒ git 自己报出的 `.scratch/wisp/issues/` 改名对＝**恰好 14 枚**（全部 `R100`＝内容一字未动），
逐枚与本腿的核（**旧名不再是跟踪路径／现名在 HEAD 是否可指／记录在哪份文件／旧名今天还被几枚入库件引着**）：

- **11 枚记在脚本 `:178-213`**（"OLD NAMES ARE TRACEABLE" 那块注释，逐枚对在 `:191-213`；fd269de1 两枚＝256/257；46079fcc 九枚＝149/154/176/177/183/253/254/259/260）：
  **旧名全部已不在跟踪集**，**10/11 枚的"现名"今天仍是跟踪路径**，旧名在入库件里的引用**逐枚都还在**
  （枚数：256→2 份、257→3、149→7、154→8、176→0、177→5、183→9、253→3、254→2、259→4、260→3）⇒ **没人把历史读数"修掉"**（AC#7 那条 ⛔ 成立）。
- **★唯一缺陷（具名）**：**257** 那枚——脚本记的现名 `257-clean-machine-provider-registry-nil-blocks-writes.md`
  **在 HEAD 已经不是跟踪路径**；HEAD 上是 `257-clean-machine-provider-registry-nil-blocks-writes-done.md`
  （`6c96a425`，2026-10-05 结案时又加了一跳 `-done`）。
  ⇒ 拿旧名的人**还是能找到对象**（`git log --follow` 一跳即穿、slug 前缀完全一致），但**"逐枚留『曾名…，现名…』一行"这句按字面要求"现名＝今天的名字"**，
  这一枚过期一跳 ⇒ **本腿判：AC#7 成立、带这一枚具名缺陷**；更正只能由**编排者**动那枚正被验收的脚本（本腿⛔ 不动），
  或者按票面 `:46` 的口径把"现名"读成"改名当时的那一跳"（那样连缺陷都不算）——**这一处留给编排者裁，本腿不替他定。**
- **3 枚记在票面 `:78-80`**（f67af953＝273/274/275）：**旧名全部不在跟踪集、现名全部在 HEAD**，旧名引用分别还在 **7／7／3 枚入库件**里
  （与 `A672` 说的"13 枚已入库件还写着旧名"同形状：7＋7＋3 有交集，本腿没做并集去重，⛔ 不据此说 A672 错）。
- **"记在票面而不记在脚本里"是否满足 AC#7 本意**——**本腿判：满足**，具名理由三条：
  ① AC#7 原文（票面 `:46`）写的是"**脚本名册或票面**里要逐枚留一行"＝**二择一**，编排者选票面是行使票面给的选择权；
  ② 那三枚**不在名册里**（它们是 10-04 之后开的票、走了改名而不是走豁免 ⇒ 名册里根本没有它们的条目可挂；
     本腿核过名册 57 枚键里**没有** 273/274/275 任何一枚）⇒ 记在脚本里反而要**新造条目**，触碰正被验收的仪器；
  ③ **可指性实测**＝上面那把 `git show -M` 尺**不需要脚本参与**就能把 14 枚逐枚复原（票面/脚本只是人读的入口，git 的对象层才是真相源）⇒ 记录面在哪一文不是牙口所在。
  ⚠ 反过来说一句本腿不替编排者隐瞒的话：**票面 `:78-80` 记的是"文件名片段"而不是全路径**（脚本那 11 枚记的是 `.scratch/wisp/issues/…` 全路径）⇒ 两半形状不齐，
  若编排者要"逐枚同一形状"，那是**追加形状**、不是 AC#7 缺件。

## 10. 还原自证（本腿交件时刻的现场尺）

见 `logs/restore-proof.txt`（取数时刻逐行写在件里）。要点：

- `sh scripts/check-path-length-budget.sh --self-test` 复跑 ⇒ **rc=0**（本腿所有突变都在 `/d/tmp/262v1/mut/` 与 clone 里，入库脚本**一字未动**）。
- `git -C <仓根> ls-files` 里**越帽跟踪路径枚数＝57＝名册 57**（与 §2/§4 同尺），**⛔ 0 枚属于本腿**：
  本腿在仓内只落了 `.scratch/wisp/probes/262/v1/**` 自己写面里的**短名**文件（全部 ≤100 字符）。
- `git status --porcelain -- .scratch scripts tools .github docs` 输出＝见 logs；本腿只在自己探针目录下留新增。
- clone／mutants／bench 产物位置（⛔ 不入库、按规矩 8 **不删**）：
  `D:/tmp/262v1/clone`（252 MB，HEAD `1d59543e`）、`D:/tmp/262v1/mut/`（7 份 mutant＋`mutants-log.txt`）、
  `D:/tmp/262v1/*.txt`（各发原始读数）、门自己的 bench＝`/tmp/tmp.Eeu8TV0TzH`（14:09）等（脚本明写"kept on disk; this project never deletes temp artifacts"）。

## 11. 本腿没做完／判不动的格（具名）

1. **AC#5＝〔待取数〕**：step 级结论没读到（gh API 6 发 `EOF`）；要的形状见 §7。另外 AC#5 原文那一枚"self-hosted 侧跑过"与裁定后的形状不同形，**待编排者二选一**。
2. **AC#0 的"墙"那一半**：self-hosted checkout 的真实断点、开区间 (206,217] 的精确化——本腿**没跑过 checkout**，⛔ 不许任何腿拿推断当实测（归 AC#5 同一次取数）。
3. **`go vet ./cmd/wisp/`（AC#6 第四门）**：**归编排者补跑销账**（本腿被硬禁、且整包窗口被另一条腿独占）；CI 的 `gofumpt -l .` 整包口径同理只欠这一枚目录。
4. **CI 的 ubuntu `unit=bytes`（mawk）分支**：本机只有 multibyte-aware awk（`unit=characters`），⛔ 本腿**没量过** byte 分支那条 REFUSE 路径；它只在真 CI run 里可测 ⇒ 挂 AC#5。
5. **推送权**：本腿全程**只 commit 不 push**（0 次 push），⛔ 任何"CI 已绿"的话都不该从这条腿的交件里读出来。

## 12. 本腿推翻／修正编排者的句子（具名）

1. **"我那发红只是 AC#3/AC#4 的天然旁证、但不是凭据"** ⇒ **本腿不同意**（论证见 §5 末段）：那发红命中同一判据（守卫 A＋C 在真实工作树上点名真实越帽路径），本腿把它记为**第二次独立读数**；本腿仍自种自拆了 A/B/C 三发，⛔ 不是因为它不算，而是为了拿到"长度已知"的那一发。
2. **AC#7 那三枚"记在票面"够不够** ⇒ **够**（§9 三条理由），⛔ 本腿不把它判成缺陷；但**同时报出一处编排者大概没料到的过期**：脚本里 **257** 那枚的"现名"今天已不是跟踪路径（`6c96a425` 又加了一跳 `-done`）——这是**真缺陷**，与那三枚无关。
3. **派单给的行号区间有三枚过期（本腿现号已换成真号，尺＝`grep -n`）**：
   阈值**常量定义**在脚本 **`:143-166`**（`HAT_NAME_LIMIT=100 :149`／`WORST_PREFIX=44 :156`／`ISSUES_PREFIX_LEN=21 :161`／`DEBT_CAP=180 :164`／`WALL_LOW/HIGH :165-166`），
   派单说的 `:55-75` 指的是**文件头那段"THE CHECK — 阈值钉在帽"的注释**（真号 `:52-72`），⛔ 不是常量定义处；
   **名册数据**在 **`:315-371`**（`#ROSTER-AWK-BEGIN :314`／`#ROSTER-AWK-END :372`，`grep -c '^    R\["'`＝**57**），
   派单说的 `:168-230` 其实是**名册的说明注释＋AC#7 旧名对**那块；
   **退出码约定**在 **`:114-118`**（注释）＋**`:224`**（usage 正文），派单说的 `:115-130` 半对；
   **`--self-test` 实现**在 **`:472-544`**（`run_selftest()` 起于 `:472`），派单说的 `:115-130` 里根本没有它。
   ⇒ 这一条不是挑刺：本腿 §6 那两发恒真检查**必须知道断言行在哪**才知道要中和哪一行。
4. **"AC#1 只裁接线在场"这一条本腿加了结构性风险**：`:166` 那一步**没带 `if:`**、前面还有 3 步 ⇒ 存在"那一步在红的 run 里 `[skipped]`"的覆盖缺口（§3）；本腿没把它算作 AC#1 的失败，⛔ 但它是 AC#5 取数时必须一起看的东西。
5. **分母不是常量**：编排者转述里 `tracked paths=7405` 与本腿 14:11 的 **7407** 差 2 枚 ⇒ **差的就是本腿自己的第 1 笔提交**；后续任何腿引用这一数都得现跑。

## 13. 裁决一览（⛔ 本腿未翻任何 AC 框；票面 8 枚 `- [ ]` 现仍未勾，尺＝§2）

| 格 | 判语 | 一句话凭据（本腿自己跑的） |
|---|---|---|
| **AC#0** | **成立（半格归 AC#5）** | `A585` 在场；`core.longpaths` rc=1、`LongPathsEnabled=0`、43 字符字串尺、分母 7350/38/19/0＋最长 180＋中间带 57＋非 ASCII 28 **全部本腿现跑**；"墙"那一半本腿没跑 checkout ⇒ 待取数 |
| **AC#1** | **成立** | 调用点 **2 枚**具名＝`ci.yml:136/:166`＋脚本 `:25-33`／`README.md:79`；未装 hook 的理由可核；⚠ 该步无 `if:` 的覆盖缺口具名 |
| **AC#2** | **成立** | 14:09 rc=0/GREEN 打出全部四枚分母读数（7405/57/57/0、longest 180＋全名、middle 57、roster 57）；14:14 干净 clone 同绿（7407） |
| **AC#3** | **成立** | 仓外 clone 自种自拆：14:14:32 **rc=1 点名 151/130/195**、14:14:41 浅目录长名 **rc=1**、14:14:42 拆掉 **rc=0**；另有仪器自带 3/3 正控与编排者 13:5x 那发（本腿记为凭据） |
| **AC#4** | **成立** | 名册 57 枚逐枚审：0 重复／0 空理由／全部可指／open-closed 说法与 `-done` 事实 0 矛盾；14:15:50 抹理由 ⇒ **RED 守卫 D**；14:15:51 删整枚 ⇒ **RED 守卫 A＋C**；恒真①：999 ⇒ RED（B＋C），删第二支帽 ⇒ 干净树绿但带种子树绿/红分形（⇒ 它有牙）；恒真②：`--self-test` 在帽被抬时 **rc=1 CONTROL FAILED** ⇒ 不恒真，断言被中和时靠计数守卫仍 RED |
| **AC#5** | **〔待取数〕** | gh step 级读数 6 发 EOF；已取到的只有 run 级（10-07 00:42Z `37545246395` sha `cc315261` **早于今天的改名**）⇒ 需要"该 step 在某 run 出现 success/failure＋日志行可指"；并挂一枚"AC#5 原文与裁定形状不同形"的待裁项 |
| **AC#6** | **成立（`./cmd/wisp/` 归编排者补跑）** | `d22scan.sh` rc=0 clean ×2（脏树 14:16:58／干净 clone 14:19:15）；`gofumpt -l`（非 cmd/wisp）干净 HEAD **0 枚**、脏树那 4 枚归在飞腿；`go vet ./internal/...` **rc=0**（14:17:39） |
| **AC#7** | **成立，带 1 枚具名缺陷** | `git show -M --name-status` 复原 **14/14 对**（全 R100）；旧名引用逐枚仍在（0–9 份/枚）无人"修掉"；缺陷＝脚本里 **257** 的"现名"今天已过期一跳（HEAD 上是 `-done` 那一枚，`6c96a425`）；"记在票面"本腿判它满足原文 |

## 14. 收尾读数与本腿自纠（原始件＝`logs/gate-final.txt`＋`logs/final-selfproof.txt`，14:25:28→14:25:48 +08）

- **门禁在 HEAD 上的第三发（⛔ 不带 self-test，纯扫）**：`sh scripts/check-path-length-budget.sh` ⇒ **rc=0／VERDICT GREEN／
  `tracked paths=7425  over-budget=57  covered by roster=57  not in roster=0`**，`longest=180`、`bands` 中间带 57、越墙 0 枚。
- **分母三次现量**：`7405`（14:09，本腿第 1 笔之前）→ `7407`（14:11，**差 2 枚＝本腿第 1 笔自己提的 verdict＋log**）→
  `7425`（14:25，本腿第 2 笔的 19 枚件）。⇒ **本票今天第三次证成"引用必带时刻"**，编排者 13:5x 那发的 7405 在它自己的时刻上是对的。
- **本腿写面纪律尺**（逐笔核，不是自述；尺＝`git show --name-only --format= <commit> | grep -vc '^\.scratch/wisp/probes/262/v1/'`）：
  本腿三笔 commit 的**非本腿写面文件枚数＝0／0／0**（文件总数 2／19／3，全在 `probes/262/v1/**`）；
  五枚禁改路径的 `git hash-object` 与 `HEAD` blob **逐枚全等**、`git status --porcelain` 对它们**逐枚 0 条**（14:26 现跑）：
  仪器 `scripts/check-path-length-budget.sh`＝`aa4cd3aa6295`、`ci.yml`＝`fc015a3cea32`、
  票面 262＝`2bceaa5107aa`、`docs/reports/pending-and-issues.md`＝`6ef8062e0dba`、`docs/reports/HANDOVER.md`＝`1add26c03a77`（左右两边同值＝本腿一字未动）；
  票面 `- [ ]` 枚数收尾仍＝ **8**（本腿 0 次翻框）；`git push` 次数＝**0**。
- **`go` 命令面**：本腿只跑过 `go vet ./internal/...`（rc=0，止 14:17:39）与 `go install mvdan.cc/gofumpt@latest`（工具安装，⛔ 不落仓内）；
  **`./cmd/wisp/` 整包本腿 0 次 `go test`／0 次 `go vet`**（那一格在 §8 具名归编排者补跑销账）。
- **两处本腿自纠**（写重了的话以本段为准）：
  ① §8 索引里"git 不存空文件"是错判（空 blob 可入库，尺＝`git ls-tree HEAD`）；
  ② §2 里"§现量 6 的两条线都不点名它们"这句本腿复述时补的半句——原文那 57 枚**今天由门禁点名**，
     规矩 9 那条词面线（名 ≤100）**仍然只有词面**（没有独立仪器去读 `README.md` 那一行），⛔ 别把这两件事读成一件。
- **本件的性质**：这是**裁决表**，不是结案文件。⛔ 本件不翻任何 AC 框、不替编排者裁 §7 那枚"AC#5 原文与裁定形状不同形"的歧义、
  也不替编排者裁 §9 那枚"257 现名过期"要不要改脚本。三处待编排者动作已在 §11／§12 逐枚点名。

## 15. 追加三发（自检的牙到底挂在哪）＋本腿两处自纠（14:28:33／14:28:38／14:29:47 +08，全在 `/d/tmp/262v1/v2/` 的拷贝里）

这三发是**恒真检查②的第二轮**，问的是"§6.2 那套读数有没有一枚其实是我自己造的假象"。原始件＝
`logs/v2-a-rosterless-selftest.txt`／`logs/v2-c-shipped-selftest-overhat-baseline.txt`／`logs/v2-d-marker-removed-refuse.txt`。

1. **★更正 §6.2 的一处机制描述**：入库副本在**真实树上满是 58 枚越帽路径**（本腿的种子仍在跟踪集）时跑 `--self-test` ⇒ **rc=0／`control 1/3 ok`**。
   ⇒ **`--self-test` 的牙不依赖真实树的任何读数**（bench＝`/tmp` 里一枚独立小仓，尺＝`BENCH=$(mktemp -d)` 在 `:473`、`run_selftest()` 起于 `:472`），
   所以名册块如果被删，**不会**像脚本注释 **`:492-498`**（那句预测逐字在 `:497-498`"the bench baseline would then be red on 57 stale entries and control 1/3 would fail out loud"）
   预测的那样"靠 57 枚 stale 让 control 1/3 响"。
   **实际响法是另一支**：只删两行 `#ROSTER-AWK-BEGIN/END` marker（`:314`/`:372`，defs 留在 `:499-500` 不动）⇒
   **整体 rc=1**，逐字 `REFUSE - the roster-less copy of itself is missing or still carries roster lines`——
   该守卫本身在 **`:501-504`**（`sed` 造副本在 `:501`、判据 `:502`、`return 2` 在 `:504`），外层 `run_selftest` 的 `return 2` 被 `:546-549` 那段转成 **`exit 1`**（⇒ 我读到的 rc 是 1，守卫自己的口径是 2，两者都对，⛔ 别混读）。
   ⇒ 判语不变（自检**确实**不恒真、且**对自己赖以工作的突变面也带守卫**），
   但**"响在哪一行"这件事本腿以现量为准、不采信脚本注释里的那句预测**——那是**注释与实际行为的一处不一致**，
   ⛔ 不是产码缺陷（两支都会非 0 退出），记在这里供编排者裁是否要改那几行注释。
2. **本腿自纠一枚突变具的缺陷**：`logs/v2-a-rosterless-selftest.txt` 那发（`rc=1`／`line 440: ROSTER_MARK_BEGIN: unbound variable`）
   **是本腿自己的 sed 写错了射程**（我用的是**未锚定**的 `/#ROSTER-AWK-BEGIN/,/#ROSTER-AWK-END/d`，它连带删掉了 `:499-500` 两行 defs），
   ⛔ **不是仪器的缺陷**——入库那枚用的正是**锚定**形式 `/^${ROSTER_MARK_BEGIN}$/,/^${ROSTER_MARK_END}$/`（`:501`），删得干净。
   本腿一度把它读成"名册块被删时自检会崩"，v2-d 那发把它否掉了。**结论以 v2-d 为准。**
3. **拆干净没**（收尾现量，14:30:1x 同一把尺）：
   **本仓**＝`tracked 7427`／`over-budget 57 枚`／**本腿名下 0 枚**／票面未勾框 **8 枚**；
   **clone**＝`tracked 7407`（与 14:14 那发起的基线**逐枚同数**）、`257-…-done.md` **仍在跟踪集**（本腿一度怀疑它被 M8–M10 摘掉，
   ⛔ **那是一句没验的猜测，现量否掉了**）、本腿两枚种子跟踪集里 **0 枚**、未跟踪残留 **2 枚**（按规矩 8 只建不删、⛔ 不入库）；
   入库五枚禁改路径的 hash 与 HEAD **逐枚全等**（`logs/final-selfproof.txt` [2]/[7]）。

⇒ **AC#4 那格的两发恒真检查到此才算完整**：恒真①＝§6.1（阈值／第二枚帽各自有牙），
恒真②＝§6.2（帽被抬时自检自拒 rc=1）＋**本节的 v2-d（marker 被动时自检 REFUSE rc=1）与 v2-c（自检与真实树无关 ⇒ 排除了"我的绿来自别人脏"这一读法）**。

⚠ **一条关于自家证据件的如实声明（不掩盖）**：本腿**入库的** `logs/mutants-log.txt`（264 行／33,215 字节，`grep -c 'RUN M8'`＝**0**）只含 §6 那批 M0–M7。
同一台机器上另有一枚**同名 bench 日志** `D:/tmp/262v1/gate-mutants/mutants-log.txt`，里面还写着 M8／M9／M10 三段（时刻 14:23:58→14:24:02，
另一份 mutant 树 `gate-mutants/`、另一处 bench `/tmp/tmp.9Y355F…`），**本腿不声称那三段的驱动语句是本腿写的**，
因此**不入库、不列为凭据**；它们与 §15 那三发（v2-a／v2-c／v2-d，**本腿亲手跑、命令逐字在 `logs/` 里**）覆盖同一组问题，
**一切以 §15 为准**（v2-a 那发还顺带暴露本腿自己 sed 的射程错误，见本节第 2 条）。
写这段只为让下一位读者知道那三段在哪、为什么不该引它——⛔ 不是新读数、⛔ 不改任何格的判语。




