# 票 308 — 拼法不变量那枚用例在 `filepath.VolumeName` 返回空串时**直接 panic 掉整枚测试二进制**：`vol[:1]` 零守卫 ⇒ 一次"临时目录没有卷标"的运行里，同包其余用例**零求值**（形状＝看着像红、实际是"这一发根本没跑完"）

**立票**：2026-10-11 11:1x 编排者（来路＝只读普查腿 `72-a2` 的自报第 4 问＋**我自己读了崩点与它的原始件**；台账 `A840` §2）
**性质**★：**仪器票**（测试面健壮性），⛔ 产码票。触发面窄（要一枚无卷标的 temp），但坏形状是"整包静默少跑"，与票 300／301／302 那一族"看着像读数、其实是⛔ 跑"同族。
**为什么要紧（用户可见形状）**：这枚用例的名字就是"同一份文件、换任何拼法都该给同一个分类"——它本来的职责是**守住拼法不变量**。它一 panic，那一发里同包**其余几十枚**（含票 307 刚落地的别名夹具）⛔ 出结果，而 CI／本地都只看到 `rc=1` 与两行名册 ⇒ 拿它当"这个包红了两枚"来读会**读错射程**。⚠ ⛔ 生产写被绕过：产码那两处同款切片我自己看过（下表现量）。

## 现量（每条带尺；行号一律当 11:1x 快照）

- ★**崩点＝测试面两枚同形**（尺＝`grep -rn '\[:1\]' internal/risk`）：
  - `internal/risk/pathresolver_anchor_spelling_windows_test.go:77 vol := filepath.VolumeName(long)` → `:78 share := strings.ToLower(vol[:1]) + "$"`；
  - `internal/risk/pathresolver_junction_windows_test.go:151 vol := filepath.VolumeName(long) // e.g. C:` → `:152 share := strings.ToLower(vol[:1]) + "$"`；
  - `tmpHome(t)`＝`pathresolver_junction_windows_test.go:22-25`，体内 `home := t.TempDir()`＋`t.Setenv("USERPROFILE", home)` ⇒ `long` 的卷标完全**由 temp env 决定**。
- ★**panic 逐字**（件＝`probes/72/a2/logs/t03-broken-surface-note.txt`，24 行，腿把被覆盖的那发原样保了一份）＝`panic: runtime error: slice bounds out of range [:1] with length 0 [recovered, repanicked]`；那一发四数＝**`rc=1`／名册只 2 枚／`PASS=0`／`FAIL=2`**（腿自己在件里写明"`PASS=0` ⛔ 代表'全红'，是测试二进制在**第 2 枚**就 panic 死了"）；⚠ 那一发的命令**没设任何 temp env**（工具跨命令⛔ 保留 shell 状态，Go 从 Git Bash 继承无卷标的 POSIX 形 `TMPDIR`）⇒ 触发条件＝**env 形状**，⛔ 是"某台机器坏了"。
- ★**产码⛔ 在射程内**（我这发自己读的，两条都带尺）：`internal/risk/pathresolver.go:262-263` 切 `drive[:1]` 之前有 `if len(drive) == 2 && drive[1] == '$' && …` 守卫；`internal/risk/pathresolver_windows.go:61 reparseComponents()` 取了 `vol` 但**⛔ 做 `[:1]` 切片** ⇒ **本票⛔ 动产码**，射程＝上面那两枚测试文件。
- ⚠ **一枚注释过期（登记、⛔ 顺手改对）**：`pathresolver_anchor_spelling_windows_test.go:70-73` 逐字写着"every other spelling is unconditional, so the test never quietly stops measuring"——在 `vol == ""` 那一形下这半句为假（它⛔ 是 quietly，但它确实 **stops measuring**，因为整枚二进制没了）。
- **枚数**：两枚用例会走那对切片（尺＝上面那条 `grep`，命中＝2 枚）；同文件其余用例（`TestClassifyAnchorSpellingIsNotVerdict` 等）⛔ 碰 `VolumeName`。

## 要建什么

- [x] **`AC#0` 先把"哪一发真能崩"现量清楚（只读）**：在**无卷标 temp**（如 `TMPDIR=/tmp` 那形）与**有卷标**两把台面上各跑一次那两枚用例所在的名册，逐枚交 `rc`＋`RUN/PASS/FAIL/SKIP`＋panic 行逐字；⛔ 预设"只有 `anchor_spelling` 会崩"（两枚同形，另一枚的实测颜色⛔ 知）。
- [ ] **`AC#1` 加一枚守卫**：判据＝无卷标那一形下，这枚用例**以可解释的颜色结束**（⛔ panic、⛔ 整枚二进制没了）；★**守卫的落点⛔ 是"把整枚用例 `t.Skip` 掉"**——那正是它守的不变量还剩下长/大写/小写/扩展长度几种拼法照跑，只有 `unc-admin` 那一枚拼法在⛔ 卷标的 temp 上无法构造（可具名去掉**那一项**，⛔ 去掉整枚）。硬约束＝⛔ 改任何判级逻辑、⛔ 动 `pathresolver.go`／`pathresolver_windows.go`、⛔ 把 `t.TempDir()` 换成手写常量（票 72 `:201` 那条禁线同形）。
- [ ] **`AC#2` 反形敏感（⛔ 只交"改⇒绿"）**：必须各交一发凭据＝(a) 守卫**拔掉** ⇒ 无卷标台面那一发**回到 panic**（红句逐字）；(b) `cmp`／`sha256` 逐字节还原 ⇒ 复绿；(c) ★**别名台面颜色⛔ 因这枚守卫而变**（改前⇄改后逐名同集合）——否则"守卫"其实是把用例关掉了。
- [ ] **`AC#3` 分母与档位**：两枚件都带 `//go:build windows` ⇒ 具名回答"CI 哪一档真跑它们、`windows` tag 在 `test-windows` 里求不求值"（先例＝票 301 包级 guard D、票 111 那 12 枚的"从未求值"）；⛔ 用"本机跑过"顶替 CI 色；拿不到 CI 色就具名登记欠读数（⛔ 由落地腿自己盖章）。
- [ ] **`AC#4` 门禁与越界**：`GOFLAGS= go build ./...`／`go vet ./internal/risk/`／`sh scripts/d22scan.sh` 各 `rc=0` 且**每把件自落一行 `rc=N`**；格式名册**两把并排且各写射程**（⚠ 本轮刚又实证一次：射程 `internal/risk` 工作树＝1 枚 CRLF 假枚、放宽到 `internal cmd`＝5 枚，见 `A839` §1 末行）；`git show --name-only --format=` **逐笔**量越界；⛔ push；commit 显式 pathspec；写面＝自家 `probes/308/**` ＋裁决表 `docs/evidence/s1/**` ＋本票末尾 bullet（⛔ 再拿"只落 `probes/**`"那把窄尺误报，见 `A838` §4／`A839` §0）。
- [ ] **`AC#5` 对账**：⛔ 翻票 72 与票 307 任何一格；本票与票 72 的关系＝**同一枚文件、不同件事**（票 72 管"期望侧写死了 temp 拼法"，本票管"崩点让整包零求值"）；⛔ 扩票 307（它已 `-done`，残余见 `A839` §3）。

## 边界（⛔ 塞进本票）

- ⛔ 甲族那 7 枚的期望侧（归票 72，`A840` §2 已换好取证尺并派腿）。
- ⛔ `overrideApplies` 那条放宽（票 72 `:95` 逐字作废，且今天它有真生产 caller `internal/tools/mode.go:91/:93/:98`，见 `A840` §1）。
- ⛔ 把票 300／301／302 那三档"windows 档在 CI 恒零求值"的搬运活并进本票（本票只交自己那两枚的档位读数）。

## 规矩（本票全程）

- 子代理**只 commit、⛔ push**；⛔ `git add -A`／`.`；⛔ `--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`／`--no-verify`；临时件**只建不删**且从一开始就落自家 `logs/`（⚠ 仓库根现存 **48** 枚深度 1 未跟踪件，`A840` §1 已登记，⛔ 再往根上丢）；证据件⛔ 叫 `.out`。
- ⛔ 为变绿放宽任何断言；⛔ 造 `--- SKIP` 换绿（本票 `AC#1` 那枚"具名去掉一项拼法"要写清它⛔ 是 skip 整枚用例）；⛔ 动 SLO 阈值／golden／`thresholds.go`／`allowlist.txt`／D43 转移表／C1–C32／D1–D47／三枚冻结件。
- ★票名⛔ 超 100 字符（长度门＝票 262 那把尺 `sh scripts/check-path-length-budget.sh --with-self-test`，本轮 11:0x 读数 **GREEN**、`worst=228`、`over the hat=57 ⇄ roster=57`）⇒ 改名或新建前先跑它。
- 裁决者≠实现者（`SPEC-12 §4.3`）：`AC#1`/`AC#2`/`AC#3` 的判语⛔ 由落地腿自勾。

next＝`308-a1`（只读答 `AC#0`：两把台面各自的颜色＋panic 逐字＋CI 档位那问）→ 我裁守卫落点 → `308-r1`（写腿）→ `308-v1`（非实现者终裁）。⚠ 排在票 72 写腿**之后**（同抢 `internal/risk` 的 `*_test.go` 写面）。

- `308-r1`（写腿，2026-10-11 13:0x，锚 `1cd57ec1`）已按编排者裁语 ⓐ–ⓓ 落 `AC#1`：新增共享 helper `driveShare`/`volumeInvariantSpellings`（`internal/risk/pathresolver_volumeless_spelling_308_test.go`，无 build tag）＋直喂无卷标字面量 `/tmp/...` 的单元用例（断言不 panic、`unc-admin`/`unc-extended` 具名去掉、长/大写/小写/扩展等仍求值，无 `t.Skip`）；两枚 `//go:build windows` 崩点件改经 helper（产码 `pathresolver.go`/`pathresolver_windows.go` 逐字节未动）；同带 `A843` 归口的两枚注释缺陷（① `probes/72/r1/20-canonical-form-choice.md`→`10-class-a-expected-side.md` §#2，② "byte-for-byte untouched"→此刻可核，皆注释行）。`AC#2` 三件齐、`AC#4` 三门禁 rc=0、`AC#5` 未翻票 72/307。读数与判语见 `probes/308/r1/{10,20,30,40,90}*.md`；`AC#1`/`AC#2`/`AC#3` 判语归 `308-v1`（⛔ 本腿自勾）；CI 实跑色＝欠账 Q-308-r1-A/B（⛔ push）。


## 编排者裁定（2026-10-11 13:1x，收 **`308-a1`** 只读普查腿，笔 `c023e8e0`→`dd638f98`→`24929768`，件 `.scratch/wisp/probes/308/a1/**`；台账 `A846`；⛔ 上面各节原句一字不改，本节追加）

★**本轮只翻 `AC#0` 一格**＝现尺 `grep -cE '^[[:space:]]*- \[x\]'`＝**1**／`- [ ]`＝**5**；`AC#1`/`AC#2`/`AC#3` 归非实现者验收（落地腿 `308-r1` 在飞）。裁语 ⓐ–ⓓ 之前只在我派单里说过 ⇒ **本节把它盘上化**（否则下一程只剩"派单里批过"那句⛔ 可核的转述）。

### 1. `AC#0` 我自己复跑的两发＋一处**我自己写错的口径**（具名回自己）

尺＝`go test ./internal/risk/ -count=1 -v -run "TestClassifySpellingInvariance|TestPathResolverUNCAListDenied"`，件＝`probes/orch/20261011/logs/orch-308a1-verify.txt`：

| 发 | env 写法 | 读数 |
|---|---|---|
| A | `TMPDIR=/tmp` | 两枚各 `--- PASS`、rc=0、`panic` 命中 **0**、终态行 `ok github.com/CarlosShao/wisp/internal/risk 0.081s` |
| B | `TMP=/tmp TEMP=/tmp` | 两枚各 `--- PASS`、rc=0、`panic` **0**、`ok … 0.070s` |

⚠★**我那句"两把台面"⛔ 成立**：Windows 档的 Go **⛔ 读 `TMPDIR`**（A 发走的是本机默认 temp），而 `TMP=/tmp` 会被补成 **`D:\tmp\…`**（腿的在场尺逐枚坐实，件 `probes/308/a1/logs/probe-T1-fwd-slash.txt`）⇒ 我那两发**都在"带盘符"那一侧**＝**两写法、一把台面**。⇒ 票面 `AC#0` 要的那发"无卷标台面"，**我和腿都⛔ 交出来**；本格闭合的是这一问的**答案**（哪一枚都不崩、为什么），⛔ 是"我复现了崩"。

产码⛔ 在射程（复认）：`internal/risk/pathresolver.go:262-263` 切 `drive[:1]` 之前逐字有 `if len(drive) == 2 && drive[1] == '$' && drive[0] >= 'a' && drive[0] <= 'z' || …`；`pathresolver_windows.go:61` 取了 `vol` 但⛔ 做 `[:1]`（那文件 `[:1]` 命中 0）。崩点取值源＝`tmpHome` 的 `home := t.TempDir()`。

### 2. ★守卫落点裁语 ⓐ–ⓓ（判语者＝我；写腿＝`308-r1`）

- ⓐ **⛔ 动产码**（射程＝那两枚 `//go:build windows` 测试件）。
- ⓑ **落点＝两枚件共用一枚 helper（把 `vol`→UNC 拼法那一步抽出、`vol == ""` 时具名报错）＋一枚"直接喂无卷标字面量"的单元用例**。★理由（这条是裁语的关键，⛔ 只写在派单里）：机器台面造不出无卷标 temp ⇒ 如果守卫只对 env 形状起作用，它就是一枚**永远不被求值的分支**；只有把 `""` 当**输入**直接喂进去，`AC#2` 的"拔掉守卫⇒回到 panic"才有一发可跑的载体。
- ⓒ **⛔ `t.Skip` 整枚用例**、**⛔ 把 `t.TempDir()` 换成手写常量**（票 72 `:201` 那条禁线同形）；只许**具名去掉 `unc-admin` 那一枚拼法**，且**分母要打印**（同片先例逐字＝`task_state_188_test.go` 那句"a scan whose denominator is zero reports clean for a tree it never read"）。
- ⓓ **两处注释缺陷并入本票写面**（腿复跑 `A843` §2 那把尺报回、我复认）：① `anchor_spelling:157` 指向 `probes/72/r1/20-canonical-form-choice.md`＝**盘上⛔ 此件**（真身＝同目录 `10-class-a-expected-side.md` 的 `## #2` 一节；⚠ 票面 `:77-78` 那两个行号锚今天也已漂到 `:103`，一律改**内容锚**）；② 同文件 `:12-16` 段称 `pathresolver_junction_windows_test.go` "byte-for-byte untouched"，而 `git show --numstat --format= f278c57f -- <那枚>`＝**`36 7`**＝**现时假话**。

### 3. ★历史 live panic 与"本机造不出"这枚矛盾具名登记（⛔ 选一边抹掉）

原始件在盘＝`probes/72/a2/logs/t03-broken-surface-note.txt`（24 行；`rc=1`／名册只 2 枚／`PASS=0`／`FAIL=2`；逐字 `panic: runtime error: slice bounds out of range [:1] with length 0 [recovered, repanicked]`；崩位＝`TestClassifySpellingInvariance`）⇒ **崩发生过**；而同一枚工具台面今天⛔ 能重放（腿 T1/T2/T3/default/bare-root 五把在场尺全非空＋我上面两发）。⇒ 判＝**形状可达已被原始件钉住、⛔ 已"本机可复现"**；`AC#2` 凭据⛔ 用 env 台面去要，按 ⓑ 那枚输入面载体去要。⚠ 那发当时的成因归因（腿写"Go 从 Git Bash 继承无卷标 `TMPDIR`"）我⛔ 复认（Windows 档 Go⛔ 读 `TMPDIR`，见 §1）＝**归因⛔ 立、现象⛔ 否**，留给 `308-v1` 单独裁这一句。

### 4. 欠读数（逐枚具名归谁，⛔ 算在腿账上）

- **两枚崩点的 CI 色**＝只有 push 之后取得到（腿只交"哪一步会打到它们"：`ci.yml:797` 整包 `portable-tests.sh --scope=windows` 会；`ci.yml:813` 那步 `-run TestPathResolverJunctionWindows` **只打 Case 1＝⛔ 覆盖这两枚**）。⇒ 归编排者，与 `A813` 那批同批取；⛔ 用"本机跑过"顶替（`AC#3` 原句）。
- **真无卷标台面那一发**＝本机⛔ 可构造（要它只能换工具/换机器）。⇒ 记为本票**永久欠读数**，`AC#3` 那格的档位读数（`//go:build windows`＋`winlive` 命中 0）腿已交形、判语归 `308-v1`。

**Status（本程更新）：** `AC#0` **已由只读腿交付＋编排者现跑翻勾**（附 §1 那枚我自己口径错具名、§3 那枚矛盾具名）；`AC#1`–`AC#5` ⛔ 动框。⛔ 零 push。

- [2026-10-11 13:1x] agent=编排者 did=收 `308-a1`（三笔 `c023e8e0`/`dd638f98`/`24929768`）⇒ **`AC#0` 翻勾**＋守卫落点裁语 ⓐ–ⓓ **盘上化**（⛔ 动产码／helper＋"直接喂无卷标字面量"那枚用例＝唯一能让守卫有牙的载体／⛔ skip 整枚、⛔ 换手写常量、只许具名去掉 `unc-admin` 且分母必打印／两处注释缺陷并入写面）；我自己现跑两发 `go test` 复认两枚崩点各自 PASS＋rc=0＋panic 0，★并**抓到我自己那句"两把台面"是假的**（Windows 档 Go⛔ 读 `TMPDIR`、`TMP=/tmp` 被补成 `D:\tmp`＝两写法一把台面）；历史 live panic 与本机不可复现那枚矛盾**逐字两边都留**⛔ 抹一边；欠读数两枚具名（CI 色归我同批取／真无卷标台面＝永久欠）；⛔ 动 `AC#1`–`AC#5`、⛔ 改原句、⛔ push

- [2026-10-11 13:3x] agent=编排者（**代提**，正文归 `308-r1`）did=收 `308-r1` 那枚⛔ 由它自己 commit 的票尾 bullet（它按"⛔ 提交别人在飞态"⛔ `git add` 本票＝`A846` §7(d) 那笔停手上报，我已核过框归属：**`AC#0` 那枚 `[x]`⛔ 它翻的、＝我 13:1x 现跑翻的**）⇒ 腿的交付逐枚在册：五笔 `510c95ca`→`1ad48333`→`ee0a93e4`→`79242df7`→`1270e36e`，写面终态 `git status` 干净、逐笔越界 0；★实现＝新件 `internal/risk/pathresolver_volumeless_spelling_308_test.go`（**无 build tag⇒跨 OS 编译**）里 `driveShare(vol)->(share,ok)`（`len(vol)==0`⇒¬ok 而非 panic）＋`volumeInvariantSpellings(home,long)->(keep,dropped)`，两枚崩点件（内容锚 `:103`/`:152`）改经 helper；产码 `pathresolver.go:263` **逐字节⛔ 动**（我已复跑守卫行逐字）。★裁语 ⓑ 那枚"直喂无卷标字面量"的用例 `TestVolumeInvariantSpellingsVolumeLess308` 真写出＝守卫的牙⛔ 挂在 env 台面挂在输入面✓；ⓒ 那两枚禁线守住＝**⛔ 整枚 `t.Skip`**（`unc-*` 具名去掉＋`t.Logf` 印分母）、**⛔ 换 `t.TempDir()`**（输入用 `/tmp/...` 字面量）；ⓓ 两枚注释缺陷（指针抄错／"byte-for-byte untouched"现时假话）已改且 `git diff -U0` 证只动注释行。`AC#2` 三件齐（拔守卫⇒逐字 `panic: runtime error: slice bounds out of range [:1] with length 0 [recovered, repanicked]` rc=1／sha256 `f2c01307…` 还原复绿／别名台面改前⇄改后逐名作差＝只 ＋1 枚新用例）。`AC#4` 三 rc=0（`go build ./...`／`go vet ./internal/risk/`／`d22scan`）。★⛔ 动任何框：`AC#1`..`AC#5` 判语归非实现者 `308-v1`；⚠ 欠读数四枚具名＝**Q-308-r1-A/B**（三枚用例在 CI `test-windows`／`test-core` 的实跑色，⛔ push 代答不了＝与 `A813` 那批同批取）、**Q-308-r1-C**（票 §现量 第三枚注释"never quietly stops measuring"在 `vol==""` 那形下为假，票面写着"登记、⛔ 顺手改对"⇒ **归我裁要不要单独改写**）、**Q-308-r1-D**（＝本笔，已结）。⛔ push
