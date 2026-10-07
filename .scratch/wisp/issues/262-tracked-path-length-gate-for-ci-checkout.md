# 票 262 — 全仓没有任何一把尺盯"跟踪路径长度"，而过长路径会让本机 runner 的 `actions/checkout` **整步失败**，`slo-full`（D32 唯一求值路径）已被两枚长名工单打死过一整天

**立票时刻**：2026-10-04 08:5x +08，锚点 HEAD `152caf25`（`dev`）
**来路**：编排者今晨现取 CI 日志 `gh run view 37158259050 --log`（2026-10-04T00:45Z 的 schedule run，`sha=941805d0`）里的 `slo-full / Run actions/checkout@v4` 两行 `##[error] … Filename too long`；台账 `A583`（立案）＋`A584`（数字更正与 9 枚改名）。
**只读普查腿 `262-a1` 已交件**（`.scratch/wisp/probes/262/a1/census.md`，191 行／16,719 字节，`a73edf62`）⇒ 本票 §现量 里带「262-a1」标记的读数出自它，编排者抽验过的另行标注。
**已就地做掉的前置**（不必本票再做）：两枚肇事工单改名＋**遗留 9 枚（相对 >180）全部改短名**（2026-10-04，`git mv`，`-done` 后缀逐枚保留）；`.scratch/wisp/issues/README.md` **规则 9：新立票文件名 ≤100 字符**（撤销口令「撤 9 号长度帽」）。⚠ **改名只让今天的检出能过，盘上没有任何仪器会拦下一次更长的那一枚**——这就是本票要补的那格。

## 现量（⚠ 引用前先重跑，别把这几行当常量；每条都要落腿自己复认，不复认就当它不存在）

1. 检出目录两形**不同宽**：self-hosted（跑 `slo-full`）＝`E:\work\base\actions-runner\_work\wisp\wisp`＝**43 字符**（拼路径＋分隔符 ⇒ 计 44），日志逐字 `Working directory is 'E:\work\base\…'`；GitHub 托管 Windows job＝`D:\a\wisp\wisp`＝**14 字符** ⇒ **同一枚文件托管侧过得去、本机过不去**。`LongPathsEnabled REG_DWORD 0x0`；仓内 `core.longpaths` **未设**（`git config --get` rc=1）。
   ⚠ 边界同形性未证：本机 git `2.52.0.windows.1` ≠ runner `2.54.0.windows.1`，两版对长路径的处理是否同形**没量过**。
2. 三发对照（同一次 run、同一枚检出）：全路径 **250 ⇒ 通过**／**261 ⇒ `Filename too long`**／**263 ⇒ 失败**。换成**与检出目录无关**的仓内相对长度＝**墙在开区间 (206, 217] 之内**（262-a1 §3 复认；要精确值必须真跑一次 checkout，本票不许拿推断当实测）。
3. 因果链最关键一环（来自 `docs/reports/pending-and-issues.md` 既有读数）：**`wisp slo` 那两个 D32 资源数字不被 `go test` 执行** ⇒ 唯一求值路径＝推送/定时触发的 `slo-check.ps1` ⇒ checkout 一死那道保护**整步归零，而盘上没有任何一行说它归零**（失败只在 job 颜色里，步名看不出原因）。
4. 分母（改名前，`git ls-files | awk '{print length($0)}'` 口径＝**仓内相对路径字符数**）：`<=100:5076`／`101-150:89`／`151-180:19`／`>180:9`，最长 206。
5. 分母（**改名后**，编排者 2026-10-04 09:0x 现跑）：`<=121:5142`／`122-150:38`／`151-180:19`／**`>180:0`**，最长 **180**（＝`.scratch/wisp/issues/252-the-two-spellings-of-one-path-split-inside-the-allowlist-check-…-r2-l2.md`）⇒ 全路径 **224**、余量 **35**。
6. ★**中间带缺口**（262-a1 §6 的射程判断，编排者复跑分母对上）：规则 9 的帽（名 ≤100 ⇒ 相对 ≤121）与 §现量 5 的"债表"线（相对 >180）之间，**盘上有 57 枚落在 122–180**，**两条线都不点名它们** ⇒ 一枚 158 字符的新票名今天**违反规则 9 却没有任何仪器会响**。⚠ 规则 9 今天**只有词面**：落地腿动手前先自己复跑这条（`git ls-files | awk 'length($0)>121 && length($0)<=180' | wc -l`），别照抄 57 这个数。

## 要建什么

- **门禁脚本**（新文件，建议 `scripts/check-path-length-budget.sh`）：读 `git ls-files`，对每枚跟踪路径算**仓内相对长度**（口径写进脚本注释：字节数还是字符数、中文会不会让两者不等）。
- ★**阈值钉在"帽"那一档，不是钉在"墙"那一档**：判据要能拦下**违反规则 9 的那一枚新票名**（相对 >121），而不是等它长到 217 撞墙才响。超预算者必须**逐枚点名**在脚本里的豁免名册中，**名册外**任何一枚超预算 ⇒ 退出码非 0。
- **中间带那 57 枚**要有具名处置：要么一次改完，要么逐枚进名册并各写一句存在理由。⛔ 不许"既不改名也不点名"（那样门永远是绿的而墙还在）。
- **接进 CI**：这一步必须跑在 **self-hosted／Windows 那一侧**才有意义（托管 Windows 前缀 14 字符、Linux 无 `MAX_PATH`，都照不到本机那 43 字符的宽）。⚠ `.github/workflows/ci.yml` 的形状是票 134 定过的 **C+B + 反静默死钉**，加一步属**契约邻域**⇒ **编排者必须先落一枚具名 `A##` 解冻记录**（文件／行／理由／边界／撤销口令），落地腿才许动它；未见到那条 `A##` 就改 ci.yml ＝对抗验收直接判失败。
- **正控**（按既有定式：负向尺必配"种 X 必响"）：脚本要有一枚"临时造一枚超预算的跟踪路径 ⇒ 门必须红 ⇒ 拆掉 ⇒ 绿"的自证，且在**本机和 CI 各出一发读数**。

## 编排者选形裁定（2026-10-04 09:2x，账 `A585`；⚠ 这一段**推翻本票 §要建什么 里我自己写的一句话**）

我原文写"这一步必须跑在 self-hosted／Windows 那一侧才有意义"——**这句是错的，落地腿别照它做**。理由：那一步只能挂在 `actions/checkout` **之后**，而一枚真的跨过墙的路径会让 **checkout 自己失败** ⇒ 挂在它后面的尺**在唯一需要它响的那种情况下永远不执行**＝死门。四形重列：

- **形甲**（挂 `slo-full` checkout 之后）＝**死门**，排除（本段即为排除的具名理由）。
- **形乙**（挂在托管 Windows job，前缀只有 14 字符）＝照不到本机那 43–44 字符的宽，单独用不够。
- **形丙＝本机 pre-push 尺**：编排者推送前自己跑（`scripts/` 里一枚脚本，不挂 CI 也能响）；射程＝推送之前，正好在东西进 CI 之前。**注意这不是"只有编排者会跑"的装饰**：本仓的推送权本来就在编排者手里（`issues/README` 规则 1/2），门在谁手里就该长在谁的路上。
- **形丁＝CI 尺，但用"具名前缀常量"复现最坏那一侧**：尺跑在**任意** job（含 Linux/托管 Windows），检查式是 `相对长度 + 具名最坏前缀常量(44) > 帽换算值` ⇒ 不需要真在那台 runner 上也能提前响。常量必须写死出处（本票 §现量 1），改它要走具名记录。
- ★**裁定＝丙＋丁同批做**（丁负责"任何人任何时刻都响"，丙负责"进 CI 之前就响"）；帽那一档换算成相对长度＝**>121 即红**（规则 9：名 ≤100 ⇒ 相对 ≤121），中间带 57 枚逐枚进名册各配理由。⛔ 不许只做丁不做丙，也不许把阈值放到 206（那等于等撞墙）。
- ⚠ **动 `.github/workflows/ci.yml` 的具名解冻已由编排者落在 `A585`**（文件／行／理由／边界／撤销口令都在那一条里）；AC#0 的"见到 A##"这一格就此满足，但**读数仍要落地腿自己现跑**。

## 判据（AC 框由编排者翻，产码腿一枚都不许碰）

- [x] **AC#0**（前置，编排者）：见到具名 `A##` 解冻记录（ci.yml 那一步）＋本票 §现量 1/2/5/6 的四把尺由落地腿自己现跑复认（不许照抄这几行）。
  - **翻（编排者 10-07 14:5x；凭据两层）**：`A##`＝`A585`（我现量尺＝`grep -c "^## A585" docs/reports/pending-and-issues.md`＝**1**）；**落地腿现跑**＝`262-r1` 件 `census.md` §1.1／§1.2／§1.3／§1.4 逐节在案（原文要的就是这一层，⛔ 不是照抄票面那四行）；**非实现者**复跑＝`262-v1` `verdict.md` §2 自己现跑三把（`core.longpaths` rc=1／`LongPathsEnabled=0`／43 字符字串尺＋分母带 `<=121:7350`·`38`·`19`·`>180:0`＋中间带 57＋非 ASCII 28）。
  - ★**第四把（墙）我这轮补到一手读数**（14:4x，用与腿同一把尺，⛔ 不是转述）：`gh run view 37158259050 -R CarlosShao/wisp --log`（**这次拉通了**＝1,889,808 字节，落 `C:/Users/swq/tmp/`＝**仓外**；该 run＝ci workflow／sha `941805d0`／`2026-10-04T00:46:50Z`／conclusion failure）⇒ `grep -a "unable to create file" | sed 's/.*unable to create file //' | cut -c1-30 | sort -u | wc -l`＝**2**，逐名＝`issues/256-…`／`issues/257-…`（与 `262-r1` §1.2 那发**同数同形**）；配 `git ls-tree -r fd269de1^ --name-only | awk '/issues\/2(5[67]|60)-/ {print length($0)}'`＝**217／219／206**（206＝260 的旧名）⇒ 换算全路径 261／263 断、250 通过 ⇒ **墙在相对开区间 (206, 217] 复认成立**。
  - ⚠ 定性不许重于证据：本格钉的是**帽那一档＋那一次真实检出的读数**；墙的**确切点数仍没人量过**（要钉死必须真撞一次检出）。腿 §2 把"墙"那半挂去 AC#5ⓑ 是它自己射程内的老实说法（它没跑 checkout），⛔ 但 AC#0 原文约束的是**落地腿**——它跑了，所以本格成立不因它那句而降档。
- [x] **AC#1**：脚本存在且被某处**真调用**（调用点枚数 ≥1，具名文件:行）；只写脚本不接线＝本格不成立。
  - 翻（14:5x）：调用点＝**2 枚**。① CI＝`.github/workflows/ci.yml:136`（步名）＋`:166`（`run: sh scripts/check-path-length-budget.sh --with-self-test`）；我这轮现量尺＝`git grep -c "Tracked path-length budget" HEAD -- .github/workflows/ci.yml`＝**1**（⇒ 没误加到别处），并把 job 名册与 `runs-on` 逐行读了：`lint:65`＝`ubuntu-latest:66`／`test-windows:505`＝`windows-latest:506`／`slo-smoke:739`＝`windows-latest:740`／`slo-full:797`＝`[self-hosted, wisp-slo]:798`／`lint-frontend:871`＝`ubuntu-latest:872` ⇒ 这一步只在**托管 Linux**那一侧。② 本机 pre-push＝约定式调用，具名在脚本 `:25-33` 与 `issues/README.md:79`；⛔ 没装 hook（本仓无 `.githooks/`、`core.hooksPath` 未设、装它不在 `A585` 面内）。
  - ⚠ 本格只裁"接线在场"，"真跑过"归 AC#5。结构性风险具名（腿 §3 提、我复认）＝`:166` **没带 `if:`**（同 job 里它前面还有 `:74`／`:83`／`:108`+`:134` 三步，而 `ci.yml:234`/`:294`/`:331`/`:387` 那些步都带了 `if: ${{ !cancelled() }}`）⇒ 这一步的覆盖率＝**前面三步都绿的那些 run**；AC#5ⓐ 那三发恰好前面全绿所以它真跑了，⛔ 不许读成"它永远会跑"。补这一枚 `if:` 要动既有步骤面＝超 `A585`，**具名留给下一程**（见文末「收 `262-v1`」§5）。
- [x] **AC#2**：在 HEAD 上跑＝**绿**，并打印分母读数（跟踪路径枚数、最长那枚的长度与全名、中间带枚数、豁免名册枚数）。
  - 翻（14:5x；**我自己那一发**）＝14:39:49 `sh scripts/check-path-length-budget.sh --with-self-test` ⇒ **rc=0／VERDICT GREEN**／`denominator: tracked paths=7433  over-budget=57  covered by roster=57  not in roster=0`／`longest=180 chars relative (.scratch/wisp/issues/252-the-two-spellings-…-r2-l2.md)`（**长度＋全名都打**）／`bands: … in the 122..180 middle=57 … roster entries=57`／`worst full path on the self-hosted runner=224 chars (hat budget 165, wall open interval (206,217])`／`unit=characters`＋字节-字符两单位同集合 cross-check ⇒ 原文那四枚分母读数**全在打印里**，不靠人自己数。
  - 其余各发：腿 14:09 本仓 **7405**／腿 14:14 仓外 clone **7407**／腿 14:25 **7425**／腿 14:30 **7427**／CI `lint` 那一步 2026-10-06T23:12:54Z **6256**＋我这发 **7433** ⇒ **分母不是常量，本票今天被证伪六次**（腿 §12-5 推翻我转述里"7405 当常量"那一枚，**我接受，记我**）；⛔ 以后任何引用必带时刻，且 CI 与本机两个口径（检出宽窄不同）不许互相替。
- [x] **AC#3**：正控＝种一枚**违反规则 9 的**（相对 >121，不必到撞墙那么长）跟踪路径后门**必须红**，红句要点名那枚路径与其长度；拆掉后恢复绿。两发读数都要在交件里（含命令逐字）。
  - 翻（14:5x；**我这发是一手**）＝上面 14:39:49 那条命令自带的正控三行逐字在案：`control 1/3 ok - bench baseline green, so this gate is not stuck red`／`control 2/3 ok - the planted over-budget tracked path is rejected and named verbatim:` ＋ `RED - over budget and NOT in the roster: .scratch/wisp/issues/seeded-…-over-the-hat.md  (relative 125 chars, name 104 chars, full path on the self-hosted runner 169 chars, budget 165)`／`control 3/3 ok - green again once the planted path is gone` ⇒ **种（相对 125＝越帽不越墙）→ 红且逐长度点名 → 拆 → 绿**，三跳我一发看全；bench＝`/tmp/tmp.KKdLnu5rF6`（按规矩 8 只建不删）。
  - 腿的直接凭据＝`262-v1` §5 全程在**仓外 clone**（`D:/tmp/262v1/clone`，14:14:05 `git clone --depth 1 --branch dev`）自种自拆：14:14:32 **rc=1** 点名（相对 151／名 130／全路径 195）；14:14:41 再种一枚**浅目录长名**（相对 108＝在相对帽之下、末段名 108＝超名帽）⇒ **rc=1**（证明第二支帽不是装饰）；14:14:42 `git rm --cached` 拆掉两枚 ⇒ **rc=0／57＝57**；命令与逐字红句在 `logs/ac3-A.txt`／`ac3-B.txt`／`ac3-C.txt`／`ac3-readings.txt`。⚠ 本仓工作树自始至终没出现过越帽跟踪路径（腿 §10 还原自证）。
  - ★**我对腿的一处反推翻：接受**（它 §12-1 不同意我把 10-07 13:5x 那发红降级成"旁证"）。裁＝那一发红的是守卫 A＋计数守卫 C，触发物是**真实越帽票名＋真实工作树＋真实调用**，与 AC#3 原文同形（原文只要求"种一枚 ⇒ 红 ⇒ 拆 ⇒ 绿，两发读数在交件里"，⛔ **没要求来源必须是验收腿自造的假名**）⇒ 记为**本判据的第二次独立读数**（红＝13:5x `the roster holds 57 entries, the tree has 60 over-budget tracked paths`；拆＝改名 `f67af953` 后 14:0x 同尺转绿；两发都在本票 Progress log 里）。同时保留区分：它**不替代** A/B/C——B 那枚"浅目录长名"分支给的是"长度已知"的形状，我那发给不了。**记我：我那句"不是凭据"是按来源降档，判据里没有这一档。**
- [x] **AC#4**：豁免名册里**每一枚**都有具名存在理由；抽查任一枚删掉理由 ⇒ 门的行为必须变化（不许"理由只是注释、删了没人响"）。
  - 翻（14:5x；**我这发是一手**）＝14:39:49 那发 **rc=0／GREEN** 这一事实本身就带四道守卫全过（守卫行号我现量：B＝脚本 `:429`、D＝`:435`、C＝`:439`；任何一枚响都不会同时打出 `not in roster=0`＋`roster entries=57`＋GREEN）。
  - 腿的静态审（`logs/ac4-roster-audit.txt`）＝名册 **57 枚**／重复键 **0**／空理由 **0**／57 枚键**全部仍是 HEAD 跟踪路径**／相对 122–180 与末段名 101–159／**0 枚**是"只靠末段名帽豁免"；**理由的事实性**也核了（写 `still open` 的枚里 0 枚带 `-done`，写 `closed…` 的枚全部带 `-done`；分类 54 枚工单＋3 枚 dispatch note）。
  - 行为突变（守卫 D 不是注释）＝14:15:50 抹掉 `252-…` 那枚理由 ⇒ **rc=1** 逐字 `RED - roster entry carries an EMPTY reason; the reason is data, not a comment: …`；14:15:51 删掉整行名册 ⇒ **rc=1**／`covered by roster=56 not in roster=1`＋计数守卫；★`logs/b.txt`（18,398 字节）＝把名册整块摘掉 ⇒ 门对**全部 58 枚**越帽路径**逐枚点名**（`RED - over budget and NOT in the roster: <全名> (relative …, name …, full path …)`）＋一枚计数守卫 ⇒ 豁免是数据、且每一枚都在被逐枚用着。
  - ★**两发恒真检查**（本格最值钱的部分，⛔ 不许略，也⛔ 不许有第二枚"改了看不出来"的形混进来）：
    ① 抬帽 `HAT_NAME_LIMIT=100→999` ⇒ **不是照旧绿而是 rc=1**（57 行"名册行已失效"＋计数守卫 `the tree has 0`）⇒ **本票禁区里"把阈值写成 999"那一支在真实树上会被当场抓住**；
    ② 删 repo-wide 第二支帽 `|| namelen > namecap` ⇒ 干净树**照旧绿**（两枚帽今天同集合），但带上 §AC#3-B 那枚种子时：入库副本 **rc=1 点名** vs 删帽副本 **rc=0／GREEN** ⇒ 这一支只对"浅目录长名"那一形负责 ⇒ **"一把尺两种形状都放行"确实在这里出现过一次**，它的凭据**只能是带种子的对照**（腿已给），⛔ 谁都不许拿"删了它门还绿"当"它没用"；
    ③ `--self-test` 自身不恒真＝帽被抬时它自己 **rc=1** `CONTROL FAILED - planting a tracked path 125 characters long should be RED (rc=1), got rc=0. The gate may be stuck green`（M5）；把 control-2 断言中和成 `if false`（M6）时自检**骗绿**，但 `--with-self-test`（CI 真跑的那条命令行）**仍 rc=1**（`RED - count guard: the roster holds 57 entries, the tree has 0`）⇒ **纵容与扫描段双层**；
    ④ 腿 §15 追加三发（v2-a／v2-c／v2-d）把"牙到底挂在哪一行"改成现量：只摘 `#ROSTER-AWK-BEGIN/END` 两行 marker ⇒ **rc=1** 逐字 `REFUSE - the roster-less copy of itself is missing or still carries roster lines`（守卫在 `:501-504`；它自己 `return 2`，外层 `:546-549` 转成 **`exit 1`** ⇒ 我读到的 rc 与守卫口径的 2 **两个都对，⛔ 不许混读**）；v2-c 又证自检与真实树无关（bench＝`mktemp -d`，`:473`）⇒ 排除了"我的绿来自别人脏"这一读法。
    ⚠ 这一发**顺带推翻脚本注释 `:497-498` 的那句预测**（"名册块被删 ⇒ 57 枚 stale 让 control 1/3 响"——实际响的是另一支 REFUSE）。⇒ **注释与实际行为不一致，具名入账（`A674`），本轮不改**：⛔ 不与本轮 AC#7 那处注释追加混进同一片改动，且它是注释级、两支都非 0 退出，不是产码缺陷。
  - ⚠ 一枚不替腿隐瞒的口径差：57 枚里 **54 枚共用同一句类目理由**（`hist: ticket filed before README rule 9`），不是 57 句各写各的故事。AC#4 原文是"每一枚都有具名存在理由"——**逐枚点名路径＋理由落在被解析的数据位＋类目句本身可核**（上面那两条一致性尺就是它的核）⇒ **我判成立**；若将来要"逐枚各一句"，那是**比票面更严的一档**，得先改判据再要，⛔ 不许偷偷按严口径打回。
- [ ] **AC#5**：CI 侧真出颜色：一次推送（或 run 重跑）后，该步在**本机 self-hosted 那一侧**跑过且有日志行可指；托管 Windows／Linux job 不误加（不误加＝不误响，但要具名说清它照不到哪一侧的宽）。
  - ★**本格拆成 ⓐ／ⓑ 两枚**（10-07 14:5x 编排者；理由＝**一格两形、一支今天满足一支结构上不可能满足**＝第 120 号那条"一个门多种红法不许压成一格"的同族。⛔ 原文那句一字不抹，就留在上面这行）。
- [x] **AC#5ⓐ**（＝AC#5 原文里"CI 侧真出颜色＋日志行可指＋托管侧不误加＋具名说清照不到哪一侧的宽"那一支）
  - 翻（14:5x；**全部我自己现跑**，`gh` 只读、⛔ 零写面）：step 级结论＝`gh run view <id> -R CarlosShao/wisp --json jobs` 里 `lint` job 的 **step 7「Tracked path-length budget (ticket 262)」＝`success`**，**三发独立 run 全中**＝`37545246395`（`updatedAt 2026-10-07T00:42:25Z`，sha `cc315261`）／`37406757402`／`37406422380`；这三发的 `headSha` 都**已含**落地那笔 `519e1ade`（尺＝`git merge-base --is-ancestor 519e1ade cc315261` 通），且 `cc315261` **是 HEAD 的祖先**（⇒ 是"改名之前那棵已推的树"）。
  - 日志行可指＝`gh run view 37545246395 --log`（**拉通**＝2,135,627 字节，落 `C:/Users/swq/tmp/`＝仓外）里逐字（时刻 `2026-10-06T23:12:54.72Z`）：`denominator: tracked paths=6256  over-budget=57  covered by roster=57  not in roster=0`／`longest=180 chars relative (…252-the-two-spellings-…)`／`worst full path on the self-hosted runner=224 chars (hat budget 165, wall open interval (206,217])`／`bands: over the hat=57 … middle=57 …`／`VERDICT GREEN - every over-budget tracked path is rostered by name with a reason, and the roster equals the tree`。⚠ 该日志的 step 归属列在 `gh --log` 里显示为 `UNKNOWN STEP`（这是 gh 的输出形状，⛔ 不是"这一步没跑"）⇒ **结论取自 jobs API、内容取自日志**，两半各指各的凭据。
  - 不误加＝`git grep -c "Tracked path-length budget" HEAD -- .github/workflows/ci.yml`＝**1**（⇒ `windows-latest` 两枚 job 与 `self-hosted` 那枚都没加），并具名说清**它照不到什么**：这一步跑在托管 Linux、检出前缀只有 14 字符量级，靠脚本里那枚**写死出处**的 `WORST_PREFIX=44`（`:156`，出处＝本票 §现量 1）复现最坏那一侧，⛔ 不是真在那台 runner 上检出过（`ci.yml:150-158` 与脚本 `:36-51` 的注释把这层限制写在了仪器旁边）。
  - ⚠⚠ 定性边界（不许借此格说"CI 已证它有牙"）：三发全是 `success`，**今天没有一发 `failure`** ⇒ "CI 里它真跑"成立，"它在 CI 里拦下过东西"**仍未成立**（拦下过的读数只有本机 13:5x 那一发）；牙本身由 AC#3／AC#4 那两批突变撑着，⛔ 三格不许互相抵。
  - ⚠ **与本页 14:0x 那句的对拉**：那时我写"这一步**从未在这三枚改名之后的树上**跑过"——**那句仍然为真**（上面三发 run 的 sha 都早于 `f67af953`）。⇒ ⓐ 证的是"这一步在 CI 里真跑过且绿＋日志行可指＋没误加"，⛔ 它**不覆盖**"改名之后的树"那一层（那一层要等一次推送）；引用本格时两件事分开说。
  - ⚠ 覆盖缺口照 ⓐ 一起记（与 AC#1 同一条）：`:166` 无 `if:` ⇒ 这三发的前三步恰好全绿所以它跑了；若将来前面某步红，这一格会 `[skipped]`＝**那次 run 不算读数**（第 109 号那条同族，本轮由腿独立复现、我这轮复认形状）。
- [ ] **AC#5ⓑ**（＝原文那句"该步在**本机 self-hosted 那一侧**跑过"）：**现形下结构上不可能满足**，⛔ 本格不勾、也不许拿 ⓐ 的读数替它填。
  - 为什么：那一步按 `A585`（只批"加一步"）落在 `lint`/`ubuntu-latest`（`ci.yml:66`），`slo-full`（`ci.yml:797-798`＝`[self-hosted, wisp-slo]`）里**没有**这一步（尺＝上面那枚 `grep -c`＝1）。要满足 ⓑ 只能＝**另一枚新的具名解冻**（往 self-hosted job 里加一步）＋一次推送取数，而机主 10-07 原话是"暂时不推远程"。
  - 〔待我裁，⛔ 不上机主的清单〕：甲＝承认"托管侧跑＋具名常量 44 复现最坏前缀"就是本票要的覆盖面，ⓑ 作为**判据遗留**登记不追（现量已支持这一读：本票真正要防的是 checkout 那一侧的宽度，而那宽度已被 44 这枚常量进算式）；乙＝另立一票／一次解冻把这一步也挂到 `slo-full`。默认走**甲**、ⓑ 长期留一枚未勾格作提醒，⛔ 不动 `ci.yml` 一个字节。
- [x] **AC#6**：卫生四门（`gofumpt`、`go vet`、`tools/d22scan`、票 212 的 ban #9）在改动后读数**不扩大**；红名集合逐名比对记录在交件里。
  - 翻（14:5x）＋**结构证明优先**：本票落地三笔改动里 **`*.go` 枚数＝0**（尺＝`git show --name-only --format= 519e1ade 6a735939 f67af953 | grep -c "\.go$"`＝**0**；逐名＝`ci.yml`＋`scripts/check-path-length-budget.sh`＋3 枚工单改名），⛔ 连我本轮那处注释追加也只动 `.sh` 的注释行 ⇒ Go 侧四门在字节层不可能被本票扩大。再加三门实测：
  - `tools/d22scan`（含票 212 ban #9）＝**我 14:4x 自己现跑** `sh scripts/d22scan.sh` ⇒ **rc=0／`clean - no D22 ban violations`**，scopes 与腿 14:16:58 那发逐字同（`bans #1-5 internal/=228`、`cmd/=38`、`ban #6 frontend/=85`、`ban #7 internal/tools/=23`、`ban #8 frontend/=85 internal/=514 cmd/=104`）⇒ **红名集合＝空**（腿另有一发干净 clone 14:19:15 同 rc=0，其 `design/=30` vs 我这边 `design/` 枚数差＝别人的未提交增删，不是门禁变化）。
  - `go vet ./cmd/wisp/`＝**我冷跑销账**：`GOCACHE=D:/tmp/gocache-262v1` ⇒ 14:37:16→14:37:50（**34 秒**）**rc=0**，缓存 184 MB；⚠ 先前那发热跑 **1 秒 rc=0 命中缓存＝不算读数**（这条记进 `A674`，同时销掉**票 275 AC#5 里那一半欠的 `go vet ./cmd/wisp/`**——那一格本身仍归票 275 的腿裁，我只是把读数还上）。`go vet ./internal/...`＝腿 14:17:39 **rc=0**。
  - ★`gofumpt` 这一门我**不写成"不扩大"，写成"两把尺量的是两种字节形"**（14:4x 现量结清，⛔ 不是假设）：腿 14:18:26 在仓外 clone（HEAD `1d59543e`）＝**0 枚**；我这轮 `/d/work/base/gopath/bin/gofumpt.exe -l internal tools/d22scan tools/mockllm`＝**4 枚**（`internal/agent/approval/pending_read.go`·`internal/agent/tools.go`·`internal/risk/provenance.go`·`internal/tools/bridge.go`）。作差尺三把：① `git diff --numstat 1d59543e HEAD -- <那 4 枚>`＝**空**（内容一字未变）；② `git status --porcelain -- <那 4 枚>`＝**空**（⇒ 腿当时说的"别人在飞的未提交改动"这半**归因不成立**，它们与 HEAD 同内容）；③ 行尾符尺 `tr -cd '\r' | wc -c`：本机工作树＝**225／1280**（带 CR），`git show HEAD:<同一枚>`＝**0**（仓内是 LF），clone 里那两枚也＝**0**；再把那 4 枚 `tr -d '\r'` 归一成 LF 副本重跑 ⇒ **报 0 枚**，保留 CRLF 的同一批副本 ⇒ **报 4 枚**；全 `internal/` 的 Go 文件 **514 枚里恰好 4 枚带 CR＝被点名的那 4 枚**。⇒ **结论＝本机 `gofumpt -l` 的 4 枚是 CRLF 幻影**（第 118 号那条的同族新实例，`core.autocrlf=true`、仓内无 `eol` 钉住），⛔ 既不是本票的账、也不是那 4 枚源码没格式化；★**全仓的幻影名册其实是 5 枚**（本格这把尺射程只到 `internal tools/**`，所以这里是 4；另 1 枚＝`cmd/wisp/models.go`，CR＝**334**、`git show HEAD:`＝**0**、`cmd/wisp` 的 100 枚 Go 里恰好 1 枚带 CR——我这轮为收 `274-v1` 那句转述自己复跑过，见 `A674` §9）；CI 的 `lint` step 8「gofmt (gofumpt)」今天＝**`failure`** 属**票 161／275 那一族的已知红**（同一批读数里 step 9 那枚 guarded census 步是 `skipped`，正落在上面那条无 `if:` 的形状上），与本票无关。
- [x] **AC#7**：**旧名可追**——9 枚改名的工单与 2 枚肇事的工单，脚本名册或票面里要逐枚留「曾名 …，现名 …」一行；⛔ 不许去改 `docs/evidence/s1/**` 与 `.scratch/wisp/probes/**` 里的**历史引文**（那是过去的读数记录）。⚠ 现量代价（编排者 2026-10-04 09:1x 自己跑）：改名之后 `docs/evidence/s1/**` 里共有 **11 行 / 10 份**证据件的"票面："那一行指旧名＝**过期指认**，逐份清单见台账 `A584`（`262-a1` 当时报的是"仅 3 行"，被这把尺推翻 ⇒ 它那两条净结论里的"代价只有 3 行"降〔仅自述〕）。
  - 翻（14:5x）：**总数我这轮复认**＝`git show --name-status -M --format= fd269de1｜46079fcc｜f67af953` 里 `.scratch/wisp/issues/` 的改名对＝**恰好 14 枚**（11 枚记在脚本 `:178-213`＋3 枚记在本票面 `:78-80`），且**全部 `R100`**＝内容一字未动。
  - "记在票面"够不够＝**够**（我接受腿 §12-2 三条理由：① AC#7 原文写的是"脚本名册**或票面**"＝二择一，我 10-07 选票面是行使票面给的选择权；② 273/274/275 **不在名册里**——它们走了改名没走豁免，名册 57 枚键没有一枚是它们，记进脚本反而要**新造条目**去动正被验收的仪器；③ 可指性实测＝`git show -M` 那把尺不需要脚本参与就能把 14 枚逐枚复原）。⚠ 形状不齐这件事照实记：票面那 3 枚记的是**文件名片段**、脚本那 11 枚记的是 `.scratch/wisp/issues/…` **全路径**；要"逐枚同一形状"属**追加形状**，不是缺件。
  - ★**腿点名的唯一缺陷我本轮修了**（具名解冻见 `A674`，撤销口令「撤 A674 名册更正」）：脚本 `:194` 那对 257 的"现名" `257-clean-machine-provider-registry-nil-blocks-writes.md` 在 HEAD **已不是跟踪路径**（`6c96a425` 结案时又加了一跳 `-done`）。修法＝**在那一对之后追加 7 行注释**说明第二跳与当前跟踪名（尺＝`git ls-files .scratch/wisp/issues | grep /257-`），⛔ **原那两行一字不改**（那块注释存在的理由就是"不许改写已记录的读数"）；改动 `git diff --numstat`＝**7/0**，改完同尺复跑 **rc=0／over 57＝roster 57／分母 7433** ⇒ **零行为变化**（这条是"注释不是数据"的正面读数，与 AC#4 里"理由是被解析的数据"互为对照）。
  - ⛔ 历史引文一格没改：旧名在入库件里逐枚还在（腿那把尺＝256→2 份／257→3／149→7／154→8／176→0／177→5／183→9／253→3／254→2／259→4／260→3；273→7／274→7／275→3），我一枚没"修"；`docs/evidence/s1/**` 与 `.scratch/wisp/probes/**` 本轮零写入（尺＝`git status --porcelain -- docs/evidence .scratch/wisp/probes` 里除我自己新增的 `probes/262/v1/**`（腿的写面）之外无改动条目）。

## 禁区

- ⛔ 不动 `SLO 阈值 / golden / internal/observe/thresholds.go` 一字节。
- ⛔ 不删任何工单内容；改名只动文件名（`git mv`），票面第一行标题**保持原文**。
- ⛔ 不改 `.gitignore` 现有未提交增量（那是别的腿的活，此刻仍脏）。
- ⛔ 不许为了变绿把阈值写成 999、把名册写成通配、或把 §现量 6 那个"57 枚中间带"删掉不提。

## Progress log

- **2026-10-04 落地腿 `262-r1`（写码，非验收）**：按上面「编排者选形裁定」做了**丙＋丁同批**那一形。
  落点：尺＝`scripts/check-path-length-budget.sh`（577 行；阈值是规则 9 的帽，检查式 `相对长度 + 44 > 165`
  等价 `相对 > 121`，另加一枚 repo-wide 的"末段名 >100"帽；名册 57 枚逐枚带理由，四道守卫 A/B/C/D 各跑了一发正控；
  内建 `--self-test` 在一次性 git 仓里种/拆一枚相对 125 的跟踪路径）；
  调用点①（丁）＝`.github/workflows/ci.yml:136` 声明 / **`:166` 命令行**（A585 边界内只加这一步：`git diff --numstat` 32/0、
  单一 hunk `@@ -135,0 +136,32 @@`，触发表 `:12`/`:36`、并发组 `:50-51`、四枚现有 job 的任何 step 与所有 `runs-on` 一字未动）；
  调用点②（丙）＝编排者推送前跑同一条命令，具名写在脚本 `:25-33`（没装 hook：本仓无 `.githooks/`、
  `core.hooksPath` rc=1，装它不在 A585 授权面内）。读数件＝`.scratch/wisp/probes/262/r1/census.md`（七节，
  §1 四把起手尺**全部复认成立**，§2 三发读数齐）。commit：`519e1ade`（脚本＋CI 步）· `6a735939`（终态读数＋脚本 `:178-213`
  的 AC#7「旧名可追」11 枚对）。
  ⚠ 三处**如实标注**别读成做过：**AC#5** 的 CI 真颜色本腿给不了（无推送权，只 commit 不 push；核验凭据列在证据件 §7.1）；
  **AC#6** 卫生四门一枚没跑（编排者指令硬禁本测量期跑 go 命令；替代自证与待跑清单见 §7.2）；
  CI 那台 ubuntu 上的 `unit=bytes` 分支（mawk）**本机没量过**，本机只有 gawk——见 §6.2 第 4 条。
  ⛔ 与本票禁区的对齐声明：§2.4 那发"把阈值常量改成 999"是**改一份 /tmp 副本**做的反形控制（证明那枚常量真在被读，
  读数：超阈枚数 57→0、没有一枚被点名成超预算），**入库那枚文件里阈值始终是 100**；名册是 57 枚逐枚点名，无通配；
  §现量 6 那 57 枚中间带在本件 §1.4/§3.1 里逐枚在案，没删不提。
  票面 `- [ ]` 那 8 枚 AC 框本腿一枚未碰（现仍 8 枚未勾）。

- **2026-10-07 14:0x 编排者追加（账 `A672`/`A673`；⛔ AC 框一枚没碰）**：★**这枚门今天第一次自己拦下东西**——我收 `35-a3` 时现跑 `sh scripts/check-path-length-budget.sh --with-self-test` ⇒ **rc=1／`VERDICT RED`／`the roster holds 57 entries, the tree has 60 over-budget tracked paths`**，三枚 not-in-roster 正是 **10-04 之后新开的工单**：`273`（名 **107** 字符）／`274`（**105**）／`275`（**106**，这枚是我自己 10-06 开的）。
  ⇒ 处置＝**改名**（`f67af953`，三枚 100% 相似＝正文一字没动），⛔ 没走"往名册加豁免"那一支：名册里那种 `hist: ticket filed before README rule 9` 的理由对 10-04 **之后**开的票是假话；⛔ 也没抬帽（脚本 `:71-72` 自己就写着抬 `HAT_NAME_LIMIT` 是本票禁区）。改后同尺 **rc=0／GREEN／over-budget 57＝roster 57**。
  **AC#7 形状的"曾名／现名"逐枚对（记在本票面而不是脚本里＝本票 AC#7 写的是"脚本名册**或票面**"；我故意不去改那枚正等着被验收的仪器）**：
  1. 曾名 `273-shipping-process-builds-the-state-machine-without-a-sink-so-every-d43-side-effect-falls-into-a-no-op.md` → 现名 `273-shipping-process-builds-state-machine-without-a-sink-so-every-d43-effect-falls-into-no-op.md`
  2. 曾名 `274-no-nail-requires-the-shipped-exe-to-carry-the-page-build-ps-step-2-comment-is-false-in-both-halves.md` → 现名 `274-no-nail-requires-shipped-exe-to-carry-page-build-ps-step-2-comment-is-false-in-both-halves.md`
  3. 曾名 `275-tracked-gofumpt-denominator-carries-a-deliberately-broken-sample-so-the-guarded-census-step-exits-2.md` → 现名 `275-tracked-gofumpt-denominator-carries-deliberately-broken-sample-so-guarded-census-step-exits-2.md`
  ⚠ 旧名在 **13 枚已入库件**里还写着（尺＝三枚旧 slug 片段各 `git grep -l <pat> HEAD` 取并集去重；逐名在 `A672` §4），⛔ 按脚本 `:184-186` 那句"不许把旧名从 probes/evidence 里'修'掉——那是某一刻的读数"，我一枚都没改。
  **给 `262-v1` 的两条**：① 我这发红是它 **AC#3/AC#4 的天然旁证、但不是凭据**（AC#3 要的是"自己种一枚越帽路径 ⇒ 门红、拆掉 ⇒ 绿"，我这回是**三枚真实越帽路径恰好已在**，同形状不同来源，它仍要自己种自己拆）；② ★它 AC#3 那一发**必须在仓外拷贝／clone 里种**（先例＝`275-v1` 的仓外 clone 突变；在本仓工作树里留任何一枚越帽跟踪路径，都会让每一枚正在跑这条门禁的腿读到红——我这轮的红就是这么来的），种完当场还原并附 `git status --porcelain -- .scratch scripts tools` 为空的自证。
  **⛔ AC#5 这一格今天仍给不了**：机主 10-07 说"暂时不推远程" ⇒ `ci.yml:166` 那一步**从未在这三枚改名之后的树上跑过**，它的绿属〔待取数〕，⛔ 谁都不许拿本机 `rc=0` 替它填。

- **2026-10-07 14:5x 编排者收验收腿 `262-v1`（账 `A674`；★AC 框由我这轮翻：8 枚已勾／2 枚未勾）**：
  **① 它的凭据层我先核自己看得见的**：五笔 commit 全部命中（`1d59543e`→`a734653c`→`b86b6d75`→`9836827d`→`2481f867`）；裁决表＝`.scratch/wisp/probes/262/v1/verdict.md`＝**42,594 字节**（363 行）＋`logs/`＝**27 枚**；它三笔 commit 的**非自己写面文件枚数＝0**（尺＝`git show --name-only --format= <c> | grep -vc '^\.scratch/wisp/probes/262/v1/'`）；五枚禁改路径与 HEAD blob **逐枚全等**（脚本 `aa4cd3aa`／`ci.yml` `fc015a3c`／本票面 `2bceaa51`／台账 `6ef8062e`／`HANDOVER.md` `1add26c0`）；`git push` 次数＝**0**；票面 `- [ ]` 交件时仍 **8**（它一枚没翻，翻框权在编排者）。⚠ 一枚**同名双份的形状**它自己报了、我也复认：`logs/a.txt`＝`logs/v2-a-…txt` 的第二次 cp、`logs/c.txt`＝`logs/v2-c-…txt` 的第二次 cp ⇒ **判语只引带 `v2-` 前缀那三枚**；而 `logs/b.txt`（18,398 字节）**没有孪生**，是"摘掉名册 ⇒ 58 枚逐枚点名"那发的唯一原件。
  **② 我这轮自己复跑的尺（⛔ 不是转述）**：门 14:39:49（rc=0／GREEN／**分母 7433**／over 57＝roster 57／内建正控三发齐）；`sh scripts/d22scan.sh` 14:4x（rc=0／clean／scopes 与它逐字同）；守卫行号 `:429`/`:435`/`:439`；`ci.yml` 步名枚数＝**1**＋job 名册与 `runs-on` 逐行；`go vet ./cmd/wisp/` **冷跑** rc=0（`GOCACHE=D:/tmp/gocache-262v1`，14:37:16→14:37:50＝34 秒，184 MB；⚠ 我第一次热跑 1 秒 rc=0 是缓存，**已作废不当读数**）；墙那把我改用 `gh run view 37158259050 --log` **真拉通了**（它 §7 连试 6 发 `EOF`）⇒ "只有 2 枚路径断、断在 217/219、206 通过"**复认**；★**AC#5ⓐ 的 CI 颜色是它没拿到的**：`gh run view --json jobs` 里 step 7 三发全 `success`（`37545246395`／`37406757402`／`37406422380`）＋日志行 `denominator: tracked paths=6256 …` 逐字可指。
  **③ 它推翻我、我接受的（逐条具名）**：(a) 我那句"我 13:5x 那发红只是旁证、不是凭据"——**按来源降档，判据里没这一档**，改记为 AC#3 的**第二次独立读数**（同时保留它 A/B/C 的直接性，B 那枚"浅目录长名"我给不了）；(b) 我把 `tracked paths=7405` 当常量转述——**分母今天被证伪六次**（6256 CI／7405／7407／7425／7427／7433），凡引用必带时刻；(c) **我派单里四处行号区间过期**（常量定义真身 `:143-166` 而非我说的 `:55-75`；名册数据 `:315-371` 而非 `:168-230`；退出码 `:114-118`＋`:224`；`--self-test` 实现 `:472-544`＝我说的区间里根本没有它）——**记我**（写派单时我引的是文件头注释的位置）。
  **④ 我推翻它一处（现量，⛔ 不是口味）**：`verdict.md` §8 把 `gofumpt -l` 脏树那 **4 枚**归因为"别人在飞的未提交改动"。我这轮三把尺把它否掉：`git status --porcelain -- <那 4 枚>`＝**空**；`git diff --numstat 1d59543e HEAD -- <那 4 枚>`＝**空**；行尾符尺 `tr -cd '\r' | wc -c`＝工作树 **225／1280**、`git show HEAD:` 同一枚＝**0**、它那枚 clone 里的同一枚＝**0**；再把 4 枚 `tr -d '\r'` 归一 ⇒ `gofumpt -l` **报 0**，原样保留 CRLF 的副本 ⇒ **报 4**；全 `internal/` 的 **514 枚 Go 文件里恰好 4 枚带 CR＝被点名的那 4 枚**。⇒ **真归因＝本机 `core.autocrlf=true` 造出的 CRLF 幻影**（第 118 号那条同族；仓内字节是 LF，⛔ CI 那侧不会被它骗）。它"干净 clone＝0 枚"那半**是对的**，我只是把两把尺为什么差结清了；这条**归票 161／275 那一族的口径警告**，⛔ 不进本票账。
  **⑤ 本票落地后的残余（四枚具名，⛔ 谁都不许顺手修）**：⑴ 脚本注释 **`:497-498`** 那句预测（"名册块被删 ⇒ 57 枚 stale 让 control 1/3 响"）**与实际行为不一致**——实际响的是 `:501-504` 那支 `REFUSE - the roster-less copy of itself …`；两支都非 0 退出 ⇒ 注释级、不是产码缺陷，本轮**不改**（⛔ 不与 AC#7 那处注释追加混成同一片）。⑵ `ci.yml:166` 那一步**没带 `if: ${{ !cancelled() }}`**（同 job 的 `:234`/`:294`/`:331`/`:387` 都带了）⇒ 这一步的 CI 覆盖面＝"前面三步全绿的那些 run"；补它要动既有步骤面＝**超 `A585`"只加一步"**，要么新解冻要么就承认这个覆盖面，〔待我裁，⛔ 不上机主清单〕。⑶ **墙的确切点数仍没人量过**，只有开区间 `(206, 217]`（要钉死必须真撞一次 self-hosted 检出）。⑷ ⓑ 那一半（原文"self-hosted 那一侧跑过"）现形下**结构上不可能满足**，默认按**甲**留一枚长期未勾格作提醒，⛔ 不动 `ci.yml` 一字。
  **⑥ 状态尺（终态，⛔ 不加 `-done`）**：本票面 `- [ ]`＝**2**（AC#5 父格＋AC#5ⓑ）／`- [x]`＝**8**；门 `rc=0`；`git status --porcelain -- scripts .github internal cmd docs/PLAN.md` 里除我自己那 7 行注释追加外**无别项**；我这轮**只 commit、⛔ 零 push**；`frontend/**`·`design/**` 零写入；禁读母本一枚没整读。
