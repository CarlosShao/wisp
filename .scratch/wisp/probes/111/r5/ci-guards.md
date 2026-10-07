# 111 / r5 — CI 步守卫普查与补守卫（写码腿 `111-r5`，2026-10-07 09:06 +08 起手）

> 本件是这一枚写腿的台件。写面只有 `.github/workflows/ci.yml` 与本目录；
> **零 `.go` 读写**；只 commit 不 push。凡编排者转述与本腿读到的原文冲突，一律以原文为准并具名报回（见 §8）。

## §0 起手锚（现量，不是抄来的）

| 尺 | 读数 |
|---|---|
| `date` | `2026-10-07 09:06 +0800` |
| `git log --oneline -1` | `4dab3fbc probes(A651 落账)：两枚隔夜死腿都零提交⇒"第45次调用先commit"` |
| `git status --porcelain \| wc -l` | **753** 行（共树在飞的未跟踪件，本腿一枚都没 add） |
| 票面 `- [ ]` / `- [x]` 枚数 | **4 / 6**（编排者说的"十枚框、四枚未勾"复现；未勾＝`:27` AC#3、`:31` AC#5、`:75` AC#7、`:77` AC#8） |
| `.github/workflows/ci.yml` 起手状态 | 工作树与 HEAD **逐字节相同**：`git diff --numstat HEAD -- .github/workflows/ci.yml` 空；blob 与 `git hash-object` 同为 `d080b5b51b3fa47cce97b795629551ab27371266` |
| 起手期间 HEAD 是否移动 | **移动了**：本腿第一次读的是 `82a10da4`，第二次读到 `4dab3fbc`；`git log 82a10da4..HEAD -- .github/workflows/ci.yml` = **0 枚** ⇒ 那几枚提交没碰 `ci.yml`，下面 §1 的名册对本腿实际编辑的那份文件成立 |

⚠ 上面那行"起手行数"这一尺本腿换了写法：`git status --porcelain` 共树里有 753 行未跟踪/在飞件，
本腿只按**显式 pathspec** 取自己那几枚（`-- .github/workflows/ci.yml`、`-- scripts`），不引全仓计数当自己的足迹。

## §1 名册：`ci.yml` 每一步 ×（有没有 `if:` / 跑什么尺）

两份产出，都是**从文件逐行抽出来的**，不是人记的：

- `logs/step-roster-before.txt` — **39 行**（＝改前**有名字的**步枚数），抽法 `logs/roster.awk`（awk，纯文本，零 go 命令）。
  列义：`job|步序|- name: 所在行|守卫|类型|步名 + run/uses 摘要`。
- `logs/yaml-census-before.txt` — 同一份文件经**真 YAML 解析器**（`python` + PyYAML，尺＝`logs/yaml-guard-census.py`）的读数：
  **含未命名的 `uses:` 步**，所以是 **51 枚步**（39 有名字 + 12 只有 `uses:`）。逐作业：

| 作业 | 步数 | 带 `if:` | 带 `if:` 的是哪几枚 |
|---|---|---|---|
| `lint` | 13 | 3 | `staticcheck`、`mockllm module vet`、`Portable tests carrier self-test`（票 111 AC#3 正控载体，r5b 那笔 `1309757b` 加的） |
| `test-core` | 7 | 1 | `Stop compose services`（`if: always()`，清容器用的） |
| `test-windows` | 9 | 7 | 该作业**全部**具名步都带 `if: ${{ !cancelled() }}`（票 111 AC#6 那一笔） |
| `slo-smoke` | 6 | 0 | — |
| `slo-full` | 5 | 0 | — |
| `lint-frontend` | 11 | 0 | — |
| 合计 | **51** | **11** | 另：`continue-on-error` 在整份文件里出现 **0** 次（解析器数的，不是 grep 数的） |

改前那枚**本腿的靶子**（`lint` 作业，名册里第 6 枚具名步）：

```
lint|6|184|NOGUARD|sh .scratch/wisp/probes/161/r5/attrib.sh --tracked-only
     name: "gofmt (gofumpt) - the tracked set is the denominator (ticket 161 AC#7 form A)"
     前一步 = lint|5|168|NOGUARD|gofmt (gofumpt)   ← 票面 §305-306 记的正是这一步今天在红
```

同族对照（编排者给的 `:650` 那一枚；本腿逐行读到的形状＝**同一枚步的 `if:` 落在 `:650`、`run:` 落在 `:651`、
`- name:` 落在 `:564`**，编排者只给了守卫那一行，行号对得上）：
`test-windows` 的 `Package coverage census (ticket 111 AC#1 + GUARD D)` **带** `if: ${{ !cancelled() }}`
⇒ 同一把普查尺，windows 腿那枚有守卫、`lint` 腿这枚没有。

## §2 票面四格原文（逐字抽自 `.scratch/wisp/issues/111-ci-tests-20-of-33-packages.md`，全件在 `logs/ac-verbatim.txt`，8 行 / 1,161 字节）

`:27-28` **AC#3**

> - [ ] **AC#3** `session`/`watchdog` 这种"空分母"要**响亮**：scope 校验加上
>       "**声明在范围内但该平台没有任何测试文件 ⇒ 直接失败**"的守卫（与票 93 的"条目腐坏即红"同族），并人为抽掉一个包证明它会红。

`:31-32` **AC#5**

> - [ ] **AC#5** 门禁：`bash -n` 改动脚本 rc=0；`sh scripts/d22scan.sh` 纯净快照 rc=0 且各 scope 不降；
>       ⚠ 新步若排在"会失败的步骤"之后 ⇒ 必须放前面或 `if: always()`（本仓实测过这道门因此从未执行）。

`:75-76` **AC#7**

> - [ ] **AC#7** `R-110-4`：票 110 承诺过"`-run TestSyncRegistryProbeLive` 经 `runtests.sh` 接进 windows 腿（step7）"，
>       验收复算发现**至今 0 次** ⇒ 要么落地并给 step7 读数，要么在票 110/111 面把它**当众改口径**（不许留在原地当已做）。

`:77-78` **AC#8**

> - [ ] **AC#8** `R-110-3`：包匹配式缺前缀锚定（`"winsec"` 宽松串曾让我把 0 命中说成 18 次）
>       ⇒ 匹配式要能区分"**被测包**"与"日志里出现过这个词"，并用一次阳性自证（种一个只在字符串里出现的包名 ⇒ 不许计入分母）。

⚠ 原文里 AC#5 那一行写的是 **`if: always()`**，而本仓**实际落地的形状**是 `if: ${{ !cancelled() }}`
（票 111 自己在 `ci.yml` 的 test-windows 作业头注释里写明了为什么不取 `always()`——**改前** `:438-465`、
本腿插入 31 行之后是 **`:469-496`**，那一段的起句现在落在 `:488`；AC#6 的凭据读数也按 `!cancelled()` 记）。
本腿**照仓里已落地的形状**补守卫，不照派单字面把 `always()` 引进来，理由与出处见 §3 末。

## §3 名册分三类（逐枚列，⛔ 不是"顺手改"）

判据来自本腿自己数的两份名册（§1）＋下面三把现量的尺。39 枚具名步里，改后**带守卫 12 枚 / 不带 27 枚**；
不带守卫的 27 枚按"为什么不带"分成三类，每一类都给尺的出处。

### 甲类＝产码门禁（**不加守卫**，11 枚）

这类步的形状是"红就得停在这儿"：`go vet` / `gofmt` / `d22scan` / 路径长度预算 / SLO 采样门 / 前端 tsc-oxlint-vite。
本票 AC#5/AC#6 要的不是它们的结论，是**它们后面那些读数**的结论；文件里已经把"谁先红"写成了因果（`:98-104`、`:122-124`、`:397-403` 三段注释），
所以动它们的守卫＝改别人的因果，而不是补自己的尺。

| 步 | 位置 |
|---|---|
| `D22 scanner positive control` / `D22 scanner self-test` / `D22 seven-ban + emoji scan` / `Tracked path-length budget` | `lint` :74 :83 :108 :136 |
| `gofmt (gofumpt)` | `lint` :168 |
| `go vet (module)` / `go vet (tools/d22scan module)` | `lint` :237 :240（本腿改后行号，原 :206/:209） |
| `Environment fork assertion` / `Portable package tests (core scope…)` | `test-core` :421 :437 |
| `SLO smoke gate` / `SLO full gate` | `slo-smoke` :759、`slo-full` :822 |
| `npm ci` / `typecheck` / `lint (oxlint)` / `token drift guard` / `build (vite build)` | `lint-frontend` :887 :890 :893 :896 :902 |

（setup 类不计：`Start compose test services` :406、`Probe mock-llm` :409、`Cache third_party` :750、`Build wisp.exe` :756/:819。）

### 乙类＝读数 / 普查 / 名册步（**该有守卫**，共 4 枚，本腿补其中 1 枚）

| 步 | 改前 | 改后 | 这一枚量什么 |
|---|---|---|---|
| `gofmt (gofumpt) - the tracked set is the denominator`（`lint`） | ⛔ **无守卫** | ✅ `if: ${{ !cancelled() }}`（新 `:234`，`run:` 新 `:235`） | tracked 名册当分母的普查尺（票 161 AC#7 form A） |
| `Portable tests carrier self-test`（`lint`，票 111 AC#3 正控载体） | 已带（r5b 那笔 `1309757b`） | 不动 | selftest 的 24 组用例 |
| `Package coverage census`（`test-windows` :681/:682） | 已带 | 不动 | `--scope=census`＝GUARD D 唯一调用点 |
| `Stop compose services`（`test-core` name `:461` / `if:` `:462`，`always()`） | 已带 | 不动 | 清容器，不是尺 |

### 丙类＝**读数步但"跳过本身是另一枚钉子的判据"**（2 枚，本腿**不加**，理由是能引到行）

- `Upload SLO report`：`slo-smoke` :762、`slo-full` :825。看着是"读数载体排在会红的门后面"，正是该补守卫的形状；**但**
  `scripts/slo-freshness.sh` P1（`:160-170`）把 **`slo-full` 作业体内出现 `if:` 或 `continue-on-error:` 判成 FAIL**
  （逐字 `D22 mode-6: a skippable job is a job that will one day be skipped`），而那枚 P1 是真在 CI 上跑的：
  `.github/workflows/slo-fresh.yml:70` 逐字 `run: sh scripts/slo-freshness.sh`（`slo-fresh.yml:51` 的 cron `23 */6 * * *`）。
  本腿用 P1 自己的尺复量了两向（`sed -n '/^  slo-full:/,/^  [a-z-]*:$/p'` ＋ `grep -cE '^[[:space:]]+if:'`）：
  **真件＝0 命中**；把守卫插进 `slo-full` 体内那份 `/tmp` 拷贝 ⇒ **1 命中**（＝P1 会红）。
  另有一层：同脚本 `:38-45` 逐字写着"FAILED 的样本**不该**留下 artifact——`the Upload step is never reached, so there is no artifact
  and the P3 clock does not move. That is on purpose`"，因为 P3 拿"`slo-full-report` 存在"当"取到了有效样本"。
  ⇒ 给 `slo-full` 那枚加守卫＝**把 P3 的判据改掉**，属放宽方向，本票不许。
  `slo-smoke` 那枚不在 P1 射程内（P1 的范围从 `  slo-full:` 起，本腿现量 `slo-smoke` 体内 `if:` ＝0），
  本腿**仍不动它**：两枚同名步一枚有守卫、一枚没有，下一个人会"顺手对齐"成 P1 那一侧，风险比这条读数值钱。
  ⇒ 记在 §8 待裁，编排者下判再由落地腿做。
- `lint-frontend` 的 4 枚 render evidence 步（`:905` `:924` `:930` `:935`）：名字带 "render evidence"，像读数；
  但判它是"打印证据"还是"断言门禁"要看 `npm run render:*` 里面，而 **`frontend/**` 是本腿的禁读面**；
  加上该作业头注释（`:869` 一带）逐字 `D22 mode-6: no \`if:\`, no \`continue-on-error\`, no skip flag, no path filter`。
  ⇒ **判不动**，不改，登记给编排者（要改得起一枚能读 `frontend/**` 的腿）。

## §4 本笔改了什么（`ci.yml` 唯一一处）

- 插入 1 行 `        if: ${{ !cancelled() }}` ＋ 30 行 ASCII 注释，落在 `lint` 作业那枚 tracked 名册步体内；
  `git diff --numstat` ＝ **31 插入 / 0 删除**；文件 910 → **941 行**；行尾仍是纯 LF（`CRLF=0`）；新增行非 ASCII 枚数 **0**。
- 同一次改动里**没有**：删步、挪步、`continue-on-error`、放宽阈值、动 `.go`、动 `scripts/**`。
- 复量（同一把尺两向，PyYAML 解析器，见 `logs/yaml-census-before.txt` / `logs/yaml-census-after.txt`）：
  `YAML_PARSE=OK` 两向都过；`TOTAL_STEPS=51` 不变；`TOTAL_GUARDED` **11 → 12**；`CONTINUE_ON_ERROR=0` 两向皆 0。
- 具名步名册尺（awk）：改前 28 无守卫 / 11 有守卫 → 改后 **27 / 12**（`logs/step-roster-before.txt`、`logs/step-roster-after.txt`）。

### §4b 顺手量到的一枚真发现：本腿补守卫那枚步，**今天在自己的命令行上是红的**

本腿把该步的命令行逐字复跑（`logs/gates-local.txt` `:18-20`）：

```
== (A) tracked tree: git ls-files -z '*.go' | xargs -0 …/gofumpt.exe -l   (tracked .go files handed to it: 935)
attrib.sh: (A) gofumpt exited 123 on the tracked set (>=2 means a file would not parse) - the tracked tree is not readable by this ruler
attrib.sh --tracked-only rc=2
```

机制本腿逐枚钉住了，不是猜：

- `git ls-files '*.go'` ＝ **935** 枚，其中 **286** 枚住在 `.scratch/**`（别的票的变异样本／快照件，`git ls-files '.scratch/**/*.go'` 同数＝286）；
- 逐枚把 gofumpt 单独拉起来数退出码 ⇒ **935 枚里恰好 1 枚退出码非 0**：`.scratch/wisp/probes/185/c1/mut/fs_broken.go`
  （票 185 故意种坏的那枚，`gofumpt -l` 对它 **rc=2**；`git cat-file -e HEAD:…fs_broken.go` ＝ **EXISTS_AT_HEAD=yes** ⇒ 它在 HEAD 上，不在"别人没保存"那一侧）；
- 于是 `xargs` 交回 **123** ⇒ `attrib.sh` 的空分母守卫按设计 **exit 2**（`ci.yml:200-203` 那段注释自己写的形状）。
- 另外那 5 枚"gofumpt 想改格式"的真实代码文件（`cmd/wisp/models.go` 等）本腿**没算进这条结论**：它们是共树里别人在飞的未提交件，
  `git ls-files` 给的是名字、gofumpt 读的是工作树内容，CI 的干净检出不会有它们。

⇒ **后果要说满**：守卫补上之后，这一枚步在 CI 的第一次读数**极可能是红（rc=2）**，红因是"tracked 名册里含一枚语法就不可解析的样本"，
不是"格式化回归"。按本票 AC#2 那句"接入第一天就红 ⇒ 那是发现：红名逐条登记进本票面并开票，不许撤步骤"，
本腿**不撤步、不改 `attrib.sh`（那是票 161 的仪器，也不在本腿可写面）**，把这一条按发现交回（§8）。
顺带一句具名对照：票 161 的注释（现 `ci.yml:191-195`）说 form (A) 治好了"`gofumpt -l .` 在 bench 树上结构性不可满足"，
本腿这量的结果是**没治好**——分母换成 tracked 集合之后，`.scratch/**` 里那枚故意坏掉的样本照样在 tracked 集合里。

## §5 AC#3（`:27-28`）——判据原文／今天哪一步在做／缺的那半条尺／本腿只证了哪一半

**判据原文拆成两句**（逐字在 `logs/ac-verbatim.txt`）：

1. "scope 校验加上『**声明在范围内但该平台没有任何测试文件 ⇒ 直接失败**』的守卫（与票 93 的"条目腐坏即红"同族）"；
2. "**并人为抽掉一个包证明它会红**"。

**第 1 句今天在做的位置（尺在盘上，本腿逐行读到）**：

| 守卫 | 代码在哪 | 谁在 CI 里跑它 |
|---|---|---|
| GUARD A＝"声明了但本平台编译出 0 个测试文件 ⇒ 红" | `scripts/portable-tests.sh:565`（注释逐字 `GUARD A: the loud empty denominator`）报错行 `:574` | `test-core` 的 `--scope=core`（`ci.yml:437`）、`test-windows` 的 `--scope=windows`（`:717`） |
| GUARD C＝命名 scope 把自己解析出的 import path 集钉住（**抽掉一个包＝红，不是少跑一个包**） | `scripts/portable-tests.sh:542-560`（报错行 `:553`） | 同上两步，另加 `--scope=census` |
| GUARD D＝全仓普查：某包本平台编译得出测试文件而**没有任何 tier 认领** ⇒ 拒退 0 | `scripts/portable-tests.sh:407-424`（报错行 `:409`），唯一调用点 `--scope=census` | `test-windows` step 8（本腿自己读的 CI 真读数见下） |

⚠ 票面前提已翻的那半句本腿不复述成事实：`:299` 记的"`internal/session` 在 scope 里却 0 个测试文件"今天不成立，
`internal/watchdog` 仍 0 枚且**不在任何 scope 里** ⇒ GUARD A 碰不到它，它今天只是被 census 打印一行 `<-NO-TESTS`（与票面 `:320` 末句一致）。

**第 2 句（缺的那半条尺）的常驻载体在哪**：`scripts/portable-tests-selftest.sh`，本腿逐枚点到用例名与行号——

- `:413` `pin-drift-one-package-missing` ＝ GUARD C 的"抽掉一个包会红"；
- `:422` `guard-a-empty-test-package` ＝ GUARD A 的"空分母直接失败"；
- `:431` `guard-b-no-result-line` ＝ GUARD B 的分母尺（AC#8 那一半也压在这台上，见 §7）；
- `:791` `census-unclaimed-package-goes-red` ＋ `:810` `census-unclaimed-zero-count-stays-green` ＝ GUARD D 的正控与负控；
- `:830` `guard-d-slice-is-what-bites` ＝ 证明 GUARD D 那段切片被换坏时这台机器会咬。

**本腿能证的那一半（现量，`logs/selftest-run.txt`）**：
`bash scripts/portable-tests-selftest.sh` 全量跑完 ＝ **`32 case(s) ran, 0 assertion(s) failed`**、`selftest_rc=0`、
逐字末句 `GREEN - every seeded anomaly was refused, and the clean scope passed.`（起 09:2x／止 09:26:47，载体 scratch 留在 `/tmp/tmp.InkKTBvcS0`，按"临时件只建不删"没清）。
⇒ 上面五枚"必须红才算通过"的用例今天**都在红**（＝"抽掉一个包会红"这台机器**会响**，不是装饰）。

**另一半欠在哪（具名，⛔ 不许说成跑过）**：那台机器**在 CI 上至今一次都没被跑过**。
本腿不引二手：用 `gh` 只读 API 把最近 **14 发 `ci` run** 的 `lint` 作业步名册逐枚拉出来（台件 `logs/ci-step-conclusions.txt`），
`Portable tests carrier self-test` 那一列 **14/14 ＝ `STEP-ABSENT`**；同件里 `test-windows` step 8（带守卫的普查步）＝ **success**（run `37545246395`／job `112547577399`）
⇒ 缺的不是"脚本对不对"，是**载体步还没进过 CI**：它在 `1309757b`（r5b 那笔）之后就没再推过，本腿的守卫补丁同样在未推的一列提交里。
票面要的第二句若按"红/绿各一次"读，那两发读数**只能在推送之后**取（本编队此轮不 push）。⇒ 本腿对 AC#3 的判语：**只成立一半，另一半判不动（射程外，要推送）**。

## §6 AC#5（`:31-32`）——四把尺今天都有真读数

| 票面逐字要求 | 本腿用的尺 | 读数 |
|---|---|---|
| `bash -n` 改动脚本 rc=0 | 本腿**没改任何脚本**；把与本步相关的 8 枚载体脚本全量 `bash -n`（`logs/gates-local.txt:5-14`） | 8/8 **rc=0**；`scripts/runtests.sh` **不存在**（rc=127）⇒ 真件是 `tools/d22scan/runtests.sh`（`bash -n` rc=0），见 §8 冲突③ |
| `sh scripts/d22scan.sh` 纯净快照 rc=0 | `git archive HEAD^`／`git archive HEAD` 两份 `/tmp` 快照（⛔ 不在仓内建 worktree），各跑一发 | 两份都 **rc=0**、都 `clean - no D22 ban violations`（`logs/d22scan-two-snapshots.txt:7-8`） |
| ……且各 scope 不降 | 同两发逐字对比 `live scope work` 行 ＋ 55 枚 `ban #` 行 diff | **逐字相同**：`#1-5 internal/=228 cmd/=38`、`#6 frontend/=85`、`#7 internal/tools/=23`、`#8 design/=30 frontend/=85 internal/=514 cmd/=104`；`SCOPES_IDENTICAL=yes` |
| ⚠ 新步若排在会失败的步骤之后 ⇒ 必须放前面或 `if: always()` | 本腿把它**推广到"已有读数步"**并按"这条规矩实际咬过谁"来量（`logs/ci-step-conclusions.txt`，14 发 `ci` run 的 `lint` 步名册） | 那枚 tracked 普查步 **14/14 全 `skipped`**（其中 3 发连 `gofmt` 自己都是 `skipped`，因为再上面的 `D22 scanner positive control` 先红）；同一 job 里带 `!cancelled()` 的两枚（`staticcheck`＝failure、`mockllm module vet`＝success）**14/14 都有结论** ⇒ 守卫在这个 job 上是**被真实数据证过有效**的，不是纸面推断 |

另外两把本腿自己加的尺（⛔ 不是"语法应该没问题"）：

- **YAML 语法**：先证明工具在＝`python -c "import yaml"` → **PyYAML 6.0.3**（`yamllint`、`ruby` 本机缺；`node` 在但没装 js-yaml，故不引）。
  尺＝`logs/yaml-guard-census.py`，对**改前/改后两份**都跑：`YAML_PARSE=OK` 两向、`TOTAL_STEPS=51` 不变、
  `TOTAL_GUARDED` **11 → 12**、`CONTINUE_ON_ERROR` 两向 **0**、逐作业 `lint 13/4`、`test-core 7/1`、`test-windows 9/7`、`slo-smoke 6/0`、`slo-full 5/0`、`lint-frontend 11/0`。
- **被补守卫那枚步自己的命令行**：见 §4b，本机 **rc=2**（机制已钉死到一枚具体文件）。

⇒ 本腿对 AC#5 的判语：**本机这半边成立**（四把尺全有读数）；`bash -n`／`d22scan` 这两把"本腿复跑件"以前欠的是"非实现者复跑"（票面 `:321` 原话），
本件就是实现者本人的复跑记录，**不能替非实现者的那一遍**；本格还牵着 AC#5 依赖的"正控进 CI"，那半在 §5 名下仍欠。⛔ 本腿不翻框。

## §7 AC#7（`:75-76`）与 AC#8（`:77-78`）——逐枚：尺在哪／今天响不响／修它要不要动别的包的测试

### AC#7 `R-110-4`（`-run TestSyncRegistryProbeLive` 经 `runtests.sh` 接进 windows 那一档）

- **尺在哪**：不在"某一步跑它"，在"某一步**点名不跑它**"。三处，全部本腿逐行读到——
  1. `scripts/portable-tests.sh:596` 的 ledger 行：`TestSyncRegistryProbeLive|./internal/risk/|windows|fixture|…`（整行 787 字符，
     末句逐字 `Remedy for a real run: a host with a sync record, or fold the shape check into the fixture-drive…`）；
  2. `scripts/portable-tests.sh:680-688` 把整张 ledger 拼成 `skip_pattern`，`:694` 逐字
     `sh "$strict" "${scope[@]}" -count=1 -skip "$skip_pattern" 2>&1 | tee "$capture"` ⇒ 它是被 **`-skip`** 的那一枚；
  3. `tools/d22scan/runtests.sh:88-99`（`skipped=$(count '^--- SKIP')`；非 0 就点名）＝"SKIP 不当成 pass"那把尺。
- **今天响不响**：**不响（就"它跑了并通过"这句话而言）**。本腿现量：`git grep -c -- '-run TestSyncRegistryProbeLive' HEAD -- scripts tools .github` ＝ **0 命中**；
  `TestSyncRegistryProbeLive` 在 tracked 面上只出现 4 处（`portable-tests.sh:7`/`:596`、`ci.yml:440`/`:695`，后两处是注释）。
  ⇒ 票面"至今 0 次"这句话**今天仍成立**，与本腿独立复算同结论。响的是另一半：ledger 每 run 对着编译出的 windows 测试二进制重验并打印"为什么它不跑"（票面 `:202-206` 已把它写成**当众改口径**）。
- **修它要不要动别的包的测试**：**要，而且两处都在本腿禁面**。用例本体在 `internal/risk/syncdirs_windows_test.go`（ledger 行逐字引 `--- SKIP at syncdirs_windows_test.go:133`）：
  要么真机 HKCU `Accounts` 键带同步记录（`slo-full` 那台 self-hosted `wisp-slo` 是唯一候选），要么把那枚 `.go` 用例改成 fixture 驱动 ⇒ **`.go` 文件＋别的包**。
  ⇒ 本腿判语：**不成立（今天 0 次没变）／落地判不动（射程外）**；本票可写的只剩"111 面当众改口径"这一支，它已由票面 `:202-206` 落过一次、并由本件 §7 复算钉到行号；
  **票 110 面上那句话不在本腿可写面**（`110-*.md` 是 `-done` 票，本腿不碰），归编排者下判（票面 `:270` 已把这一支点名给编排者）。

### AC#8 `R-110-3`（包匹配式缺前缀锚定）

- **尺在哪**：两把匹配式，本腿逐字读到——
  1. `scripts/portable-tests.sh:720` `escape_re()`（正则转义）、`:730` `result_tail='[[:space:]]+([0-9]+\.[0-9]+s|\[build failed\])$'`（时长尾巴）、
     `:736` `grep -E "^(ok\|FAIL)[[:space:]]+$(escape_re "$pkg")${result_tail}"` ＝ 行首锚＋转义包名＋尾巴，`:749` 缺行即红；
  2. `scripts/winsec-tests.sh:122-123` 同形状（`:109` 注释逐字"the line is anchored to ^ok/FAIL, the two tokens Go uses ONLY for a…"）。
  ⇒ **"被测包 vs 日志里出现过这个词"这一半在盘上**，`"winsec"` 那类宽松串已经不是匹配式了。
- **今天响不响**：**分母尺响、票面点名的那枚阳性自证不响**。
  响的：GUARD B 的 own-line 表在 `test-core`／`test-windows` 两步每 run 都打（票面 `:313` 记 core 25 行／windows 8 行，本腿不重述成自己的读数）；
  载体侧本腿现量＝`guard-b-no-result-line` 在 selftest 全绿的那 32 枚里（§5），它证的是"**缺** own-line 会红"。
  不响的：票面第二句"**并用一次阳性自证（种一个只在字符串里出现的包名 ⇒ 不许计入分母）**"——本腿按"这句话在说什么"扫（不扫符号名），
  selftest 的 24 枚用例名册里没有这一枚诱饵用例：`grep -nE '9\.999s|decoy|githubXcom|fmt\.Println.*ok' scripts/portable-tests-selftest.sh` ＝ **0 命中**，
  `seed-stdout-*` 那几枚种进的是 `go list` 的 stdout（GUARD C 的地盘），不是 `go test` 的结果行。
  旧腿那两发诱饵实验（合成 capture 4 枚／快照里种 `ac8_decoy_test.go`）只活在票面 `:209-214` 的文字里，**今天不可重放**。
- **修它要不要动别的包的测试**：**不要**。要动的只有 `scripts/portable-tests-selftest.sh`（加一枚用例：经 `FAKEGO_*` 往 capture 里吐一条以假乱真的
  `ok  github.com/CarlosShao/wisp/<pkg>  9.999s`，断 GUARD B 的表里**没有**那一行）＋ 它已在 `ci.yml` 里的那枚载体步。
  ⇒ 本腿判语：**尺在／分母那半响／票面点名的阳性自证不成立（无常驻载体，旧证据不可重放）**；
  写这枚用例是 `scripts/**` 面的活、且属"另立一票"（票面 `:264` 那种"三张要开的票"的形状），**本票没授权本腿自造判据**，故只做判定不动手。

## §8 具名报回（与编排者转述不符之处）＋待裁＋next

**① AC#5 的行号**：派单写"当前未勾＝`:27` AC#3、**`:29` AC#5**、`:75` AC#7、`:77` AC#8"。原文实测 AC#5 在 **`:31`**，
`:29` 那一行是**已勾的 AC#4**（`cmd/wisp` 那一格）。其余三枚行号对。未勾枚数本腿复量＝**4**（`- [ ]`）／已勾 **6**（`- [x]`），与派单"十枚框"一致。

**② `if:` 的字面形状**：派单与票面 `:32` 都写 `if: always()`；本仓**已落地并已被真实数据证过**的形状是 `if: ${{ !cancelled() }}`
（`ci.yml` test-windows 头注释逐字论证为什么不取 `always()`，那段现落在 **`:488-496`**＝本腿插入后的行号，改前是 `:457-465`；
票面 `:43` 与 AC#6 的凭据读数是"同 job 只有 `Stop compose services` 用 `always()`，其余全是 `!cancelled()`"）。
本腿按已落地形状补，并在新增注释里点名这处差别（⛔ 没改票面一字）。

**③ "经 `runtests.sh` 接进 windows 那一档"的措辞**：仓里**没有** `scripts/runtests.sh`（本腿 `bash -n` 撞到 rc=127），
真件是 `tools/d22scan/runtests.sh`，由 `scripts/portable-tests.sh:694` 以 `sh "$strict"` 调用。票面 `:50` 自己就登记过同类"按猜出来的路径取件"的坑，本条与它同族。

**④ 派单引的 run `37406757402` 步 8/步 9**：本腿**自己复验了，结论一致**（`logs/ci-step-conclusions.txt`：
step 8 `gofmt (gofumpt)`＝failure、step 9 `…the tracked set is the denominator`＝skipped，lint job `112085927883`）。
但补一条派单没说的：这**不是个例**——最近 14 发 `ci` run **同一形状 14/14**，其中 3 发（`37249563077`/`37240161874`/`37166458550`）连 `gofmt` 自己都是 skipped，
因为它上面的 `D22 scanner positive control` 先红。⇒ 这把尺**从未在 CI 上执行过**，比"某一发被跳过"更强。

**⑤ 派单说"CI 侧的真实读数这台机器拿不到——必须推送才会有"**：这句本腿要**折扣后**报回——
`gh`（2.96.0，已登录 `CarlosShao`，token scopes 含 `repo`）能**只读**取历史 run 的步级结论，本件 §5/§6/§8 的 CI 读数全部出自它。
拿不到的只是**本腿这笔改动之后**的读数（那才要 push）。本腿没做 push、没改任何 workflow 触发、没跑 `gh run rerun`/`workflow run`。

**待裁（本腿不动手，等编排者下判）**
- 丙类那两枚 `Upload SLO report`（`slo-smoke :762` / `slo-full :825`）：`slo-full` 那枚**不能**加守卫（`scripts/slo-freshness.sh` P1 `:160-170` 判 FAIL，
  而那把尺由 `.github/workflows/slo-fresh.yml:70` 在 CI 真跑、最近三发 slo-fresh 全 `success`＝它今天是活的），
  `slo-smoke` 那枚不在 P1 射程内、可以加——本腿**没加**，理由是不让两枚同名步分叉成下一个人"顺手对齐"到 P1 那一侧。要加请点名。
- `lint-frontend` 的 4 枚 render evidence 步：判不动（要读 `frontend/**` 才能分类，而那是本腿禁读面；该作业头注释逐字写着 `no \`if:\``）。
- §4b 那枚发现要不要另立一票（tracked 名册里含一枚语法不可解析的样本 ⇒ 补守卫后的普查步第一发会红）：
  本腿按 AC#2"红名逐条登记进本票面并开票"的精神写进票面追加行，**不撤步、不改别人的仪器**（`attrib.sh` 是票 161 的）。

**next=（都要 push，本腿不接受拿本地绿冒充 CI 绿）**
1. 推送后第一发 `ci` run：`lint` 作业里 `…the tracked set is the denominator` 那一步**必须有 success 或 failure**（⛔ 再是 skipped 就是没修好）；
   同发里 `Portable tests carrier self-test` 应第一次出现（它带守卫，理应给出结论）——这一枚同时收 AC#3 第二句与 AC#8 的载体欠账。
2. 若那一枚普查步红（rc=2）：按 §4b 的机制记账到 `.scratch/wisp/probes/185/c1/mut/fs_broken.go`，
   归"tracked 分母该不该排除故意坏掉的样本"这一枚**口径**问题（改 `attrib.sh` 属票 161 地界，本腿不动）。
3. `internal/watchdog` 那句"响亮"今天只是打印 `<-NO-TESTS`：票面 `:320` 已具名，本腿复量它**不在任何 scope 里** ⇒ GUARD A 碰不到，属实现票 42 落地时的耦合（票面 `:268` ③）。
