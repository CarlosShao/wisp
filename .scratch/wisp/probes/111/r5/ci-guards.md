# 111 / r5 — CI 步守卫普查与补守卫（写码腿 `111-r5`，2026-10-07 09:06 +08 起手）

> 本件是这一枚写腿的台件。写面只有 `.github/workflows/ci.yml` 与本目录；
> **零 `.go` 读写**；只 commit 不 push。凡编排者转述与本腿读到的原文冲突，一律以原文为准并具名报回（见 §7）。

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
（票 111 自己在 `ci.yml:438-465` 写明了为什么不取 `always()`；AC#6 的凭据读数也按 `!cancelled()` 记）。
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
