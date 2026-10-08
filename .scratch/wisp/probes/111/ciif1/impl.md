# 111-ciif1 — `.github/workflows/ci.yml` 步级 `if:` 守卫补挂（落地腿交付）

工作目录 `D:/work/workspace/projects plans/Wisp`，分支 `dev`。
起手锚：`6320a16a`（不是我被告知的 `e05b8a2b`，见 §6 T-1）。我的第 1 笔：`0c9726f9`（只加证据件，`ci.yml` 未动）。
本机日期 2026-10-08。**零 push。**

---

## 0. 尺（先声明尺子，再给数；两把尺分时刻）

**结构尺**＝`<缩进>＋<if: 语句形状>`，只算非注释行。两条独立实现必须一致：
`.scratch/wisp/probes/111/ciif1/roster.py`（PyYAML `safe_load` 走 `jobs.*.steps[].if` 键 ＋ 纯缩进行扫描，两 pass 互验）。

| 尺 | 改前（`6320a16a`） | 改后（本程） | 读数件 |
|---|---|---|---|
| 步级 `if:` | **13** ＝ 12 `!cancelled()` ＋ 1 `always()` | **36** ＝ 35 `!cancelled()` ＋ 1 `always()` | `logs/roster.txt` / `logs/roster-after.txt` / `logs/rulers-after.txt` |
| 作业级 `if:`（`^    if: `） | **0**（rc=1） | **0**（rc=1） | `logs/rulers.txt` / `logs/rulers-after.txt` |
| 词频尺 `grep -c 'cancelled'`（含注释，**不是**结构尺） | **27**（⛔ 不是派单里的 25，见 §6 T-2） | **50** | 同上 |
| `continue-on-error` | 0 | 0 | 未新增（我只加 `if:` 行） |

`13 → 36` 差 **23 枚＝本程补的枚数**。yaml pass 与缩进 pass 两向皆相等（`roster.txt: == AGREE? == ... -> True`）。

---

## 1. 名册（逐 job、逐 step、带 HEAD 行号；`- name:`/`- uses:` 那行的行号）

行号全部取自我自己 dump 的 `logs/step-open-lines-before.txt`（`git show HEAD:.github/workflows/ci.yml`，rc=0）。
"守卫"列＝**改前**该步有没有步级 `if:`；`+n` ＝本程新增，括号里是**改后**文件里那枚 `if:` 的行号。

### lint（job 键 `:65`，13 步，runs-on ubuntu-latest）
| # | step | 行号 | 改前守卫 | 本程 |
|---|---|---|---|---|
| 1 | `uses: actions/checkout@v4` | :68 | 无 | setup，不动（§3 定式①） |
| 2 | `uses: actions/setup-go@v5` | :70 | 无 | setup，不动 |
| 3 | D22 scanner positive control | :74 | 无 | **+ :81** |
| 4 | D22 scanner self-test (161 AC#2) | :83 | 无 | **+ :106** |
| 5 | D22 seven-ban + emoji scan | :108 | 无 | **+ :136** |
| 6 | Tracked path-length budget (262) | :136 | 无 | **+ :169** |
| 7 | gofmt (gofumpt) | :168 | 无 | **+ :173** |
| 8 | gofmt (gofumpt) - tracked set denominator (161 AC#7 A) | :184 | `!cancelled()` :234 | 已有 |
| 9 | go vet (module) | :237 | 无 | **+ :243** |
| 10 | go vet (tools/d22scan module) | :240 | 无 | **+ :247** |
| 11 | staticcheck | :244 | `!cancelled()` :294 | 已有 |
| 12 | mockllm module vet | :320 | `!cancelled()` :331 | 已有 |
| 13 | Portable tests carrier self-test (111 AC#3) | :335 | `!cancelled()` :387 | 已有 |

### test-core（job 键 `:395`，7 步）
| # | step | 行号 | 改前守卫 | 本程 |
|---|---|---|---|---|
| 1 | checkout | :400 | 无 | setup，不动 |
| 2 | setup-go | :402 | 无 | setup，不动 |
| 3 | Start compose test services (mock-llm 18080) | :406 | 无 | **+ :414** |
| 4 | Probe mock-llm | :409 | 无 | **+ :418** |
| 5 | Environment fork assertion (WISP_ENV=test) | :421 | 无 | **+ :444** |
| 6 | Portable package tests (`--scope=core`) | :437 | 无 | **+ :469** |
| 7 | Stop compose services | :461 | `always()` :462 → 改后 :473 | 已有，不动 |

### test-windows（job 键 `:505`，10 步）——**本程 0 枚**
| # | step | 行号 | 改前守卫 |
|---|---|---|---|
| 1 | checkout | :510 | 无（setup，不动） |
| 2 | setup-go | :512 | 无（setup，不动） |
| 3 | Windows ACL sealing gate (110) | :516 | `!cancelled()` :522 |
| 4 | Cache third_party | :559 | :560 |
| 5 | cgo build smoke | :566 | :570 |
| 6 | cmd/wisp CLI tests | :573 | :592 |
| 7 | winlive compile gate (111 AC#11) | :595 | :643 |
| 8 | Package coverage census (111 AC#1) | :646 | :732 |
| 9 | Portable windows tests | :735 | :768 |
| 10 | PathResolver junction placeholder | :771 | :784 |

`:469` 那句 "EVERY STEP BELOW THE SETUP ONES CARRIES `if: ${{ !cancelled() }}`" 在这枚 job 上**已经成立**（8/8），我复量后确认无事可做。

### slo-smoke（job 键 `:790`，6 步）
| # | step | 行号 | 改前守卫 | 本程 |
|---|---|---|---|---|
| 1 | checkout | :795 | 无 | setup，不动 |
| 2 | setup-go | :797 | 无 | setup，不动 |
| 3 | Cache third_party | :801 | 无 | **+ :813** |
| 4 | Build wisp.exe | :807 | 无 | **+ :820** |
| 5 | SLO smoke gate | :810 | 无 | **+ :824** |
| 6 | Upload SLO report | :813 | 无 | **不动**（§3 定式③，`ci.yml:229-233` 已登记为刻意） |

### slo-full（job 键 `:848`，5 步）——**本程 0 枚，硬约束**
| # | step | 行号 | 改前守卫 | 本程 |
|---|---|---|---|---|
| 1 | checkout | :863 | 无 | 不动（setup） |
| 2 | setup-go | :865 | 无 | 不动（setup） |
| 3 | Build wisp.exe (deps cached on the runner) | :870 | 无 | **不动**（定式②） |
| 4 | SLO full gate (six states + settle + leak) | :873 | 无 | **不动**（定式②） |
| 5 | Upload SLO report | :876 | 无 | **不动**（定式②＋③） |

### lint-frontend（job 键 `:922`，11 步）
| # | step | 行号 | 改前守卫 | 本程 |
|---|---|---|---|---|
| 1 | checkout | :928 | 无 | setup，不动 |
| 2 | setup-node@v4 | :930 | 无 | setup，不动 |
| 3 | npm ci | :938 | 无 | **+ :953** |
| 4 | typecheck (tsc -b) | :941 | 无 | **+ :957** |
| 5 | lint (oxlint) | :944 | 无 | **+ :961** |
| 6 | token drift guard (C21 table) | :947 | 无 | **+ :968** |
| 7 | build (vite build -> frontend/dist) | :953 | 无 | **+ :972** |
| 8 | L2 card render evidence (AC#3) | :956 | 无 | **+ :986** |
| 9 | Composer states render evidence | :975 | 无 | **+ :999** |
| 10 | Streaming output SSE render evidence | :981 | 无 | **+ :1005** |
| 11 | Nav rail icons render evidence | :986 | 无 | **+ :1013** |

**合计：52 步 / 6 job；改前 13 枚步级 `if:`，本程 +23，改后 36。**
改后仍**未带守卫的 16 枚**逐枚点名＝12 枚 setup（6 个 job 各自的 checkout/setup-go/setup-node，定式①）＋ slo-full 的 `Build wisp.exe` / `SLO full gate` / `Upload SLO report`（定式②＋③）＋ slo-smoke 的 `Upload SLO report`（定式③）。**没有一枚是漏的。**

---

## 2. 哪些格历史上真的 `[skipped]`（真读数，不是推测）

取数法＝仓里既有的 `gh`：
`gh run list --workflow=ci.yml --limit 14`（`logs/run-list.json`，rc=0）→ 14 发**全部** `conclusion=failure`，最近一发 2026-10-07T23:44Z，headSha 全是已推的 `cc315261`（⛔ 不是我本机 HEAD，见 §6 T-6）。
逐 step 结论：`gh run view <R> --json jobs --jq '.jobs[] | .steps[]? | "\(.name)|\(.conclusion)"'`，R ∈ {`37703959747`, `37545246395`, `37406757402`} → **156 行读数**落 `logs/step-conclusions.txt`（rc=0；每发各自 `view rc=0`）。⛔ 只取步名＋conclusion 两列，一行日志正文都没进上下文。

**三发同形，`[skipped]` 的格只有这些：**

| job | 步 | 三发结论 |
|---|---|---|
| lint | `gofmt (gofumpt)` | **failure**（＝吃掉下面三枚的那枚红） |
| lint | `gofmt (gofumpt) - the tracked set is the denominator` | **skipped ×3**（它带守卫，但守卫在 `cc315261` 上还不存在 → 见 T-6） |
| lint | `go vet (module)` | **skipped ×3** |
| lint | `go vet (tools/d22scan module)` | **skipped ×3** |
| lint | `staticcheck` | failure ×3（有守卫，照跑） |
| lint | `mockllm module vet` | success ×3 |
| test-core | `Portable package tests (--scope=core)` | failure ×3 |
| test-windows | `cmd/wisp CLI tests`／`Portable windows tests` | failure ×3；其余 6 枚 success ×3（AC#6 的形态在真读数里成立） |
| slo-smoke / slo-full / lint-frontend | 每一步 | **success ×3**（含 `SLO smoke gate`、`SLO full gate`、11 枚 frontend 门） |
| 三枚 job 的 `Post Run actions/setup-go@v5`／`Post Cache third_party` | **skipped ×3** | runner 自动生成的 post 步，**不在 yaml 里、我改不了**；不算"死门"，但它是这次普查里唯一还在 skip 的形状 |

⇒ **历史上真的从未被求值过的判据只有 3 枚**：`gofmt denominator`、`go vet (module)`、`go vet (tools/d22scan module)`。这 3 枚与票 111 `:394` 那句（"在 12 发里全是 skipped"）逐字对得上，也与 `:311`/`:321` 登记的"lint 那三枚今天仍是 skipped 且无日志，归票 236 AC#6／票 85 地界"是同一笔账。
⇒ 其余 **20 枚**我补的守卫属于另一类：**这三发没 skip，但结构上随时会被 skip**（它们头顶的步全是 `run:`/`uses:`，都可能红；`gofmt` 与 `staticcheck` 12/12 红已经示范过这件事）。`success` 不等于这道门拦过任何东西，只有 `failure` 才是拦下的读数——这两条定式我按派单原话执行，没扩写。

**没取到的读数（具名，⛔ 未编）**：
- 我想再取一发**带 job 列**的读数与一发更早的 run（`35600043583`，票 111 `:128` 说 `go vet (module)` 在那里＝success）做交叉验证：`gh` 从 2026-10-08 起连续 `net/http: TLS handshake timeout`，`logs/step-conclusions-named.txt` 三条 `view rc=1` 就是那三发的失败现场。job 列我改用**分组边界**恢复（每发以 `view <R> rc=0` 分隔，job 顺序与 `logs/step-open-lines-before.txt` 的名册一一对得上），这条推断**写在证据件里可核**，不是新读数。
- 21 枚"结构上会被 skip"的步，**在过去若干发里是否真的红过**：〔待取数〕——我只取了 3 发。

---

## 3. 我补了哪 23 枚、以及三条不动的定式

**动因（结构判据）**：GitHub 给每一步默认 `success()` 检查 ⇒ 同一 job 内任一前置步变红，后面的步整枚 `[skipped]`＝那道门的判据**从未被求值**。我按"该步头顶是否存在可能变红的步"逐枚判定。

**三条我拒绝扩写射程的定式（出处都是原文，不是我的判断）：**

① **setup 两步不补**（每个 job 的 1、2 枚，共 10 枚）。`ci.yml:469` 把这条例写明写："EVERY STEP BELOW THE SETUP ONES CARRIES `if: ${{ !cancelled() }}`"；`:521` 给的理由是"它自身带守卫，这样一次红的 checkout/setup 不能让它变成 inconclusive"——守卫的语义是**从 setup 那一层往下**开始的。
② **slo-full 整个 job 一步都不补**（`:863 :865 :870 :873 :876` 五枚全不动，其中 3/4/5 结构上确实会被 skip）。硬约束在仪器里：`scripts/slo-freshness.sh:160-167`
```
job_body=$(sed -n '/^  slo-full:/,/^  [a-z-]*:$/p' "$ci_file")
for banned in 'if:' 'continue-on-error:'; do
    if printf '%s\n' "$job_body" | grep -qE "^[[:space:]]+${banned%:}:"; then fail ...
```
⇒ 在 slo-full 体内出现任何一行 `if:`（含其内注释之外的所有真实键）都会把这枚钉子判红：`slo-full-trigger-missing: the slo-full job now carries 'if:' (D22 mode-6)`。票 111 `:229-233` 早已把这一条写进文件正文。
③ **两枚 `Upload SLO report`（slo-smoke `:813`、slo-full `:876`）不补**——`ci.yml:229-233` 逐字登记为刻意："the two `Upload SLO report` steps are deliberately NOT given a guard even though they carry readings … a FAILED sample must not produce an artifact, because P3 ages on 'a valid sample happened'"。失败的样不该产出工件；给了守卫就造得出工件。

**补的 23 枚**＝lint 3/4/5/6/7/9/10（7）＋ test-core 3/4/5/6（4）＋ slo-smoke 3/4/5（3）＋ lint-frontend 3–11（9）。
每枚都是**一行、缩进 8、内容 `if: ${{ !cancelled() }}`**，落点在该步自己的 `run:`/`uses:` 键**紧前面**——这跟文件既有形状一致（`:522` 在 `shell: bash`＋注释块之后、`:784` 在 `shell: bash` 之后）。形状选择由 `logs/insert-plan-dryrun.txt` 先跑 dry-run 落件，`--apply` 才写（`logs/insert-applied.txt`）。

⛔ 未动：任何命令内容、阈值、断言、job 级 `if:`（本来 0 枚，改后仍 0 枚）、`runs-on`、matrix、action 版本、步序、`continue-on-error`（仍 0 枚）。⛔ 未加注释（加注释＝非 `if:` 行，越界）。

---

## 4. 自证：只有 `if:` 行

```
git diff --numstat .github/workflows/ci.yml   ->   23      0      .github/workflows/ci.yml   (rc=0)
```
**23 插入 / 0 删除**——没有"位移"可言，因为没有任何行被删或改写。
逐字核：把 diff 里所有新增行拿去匹配 `^\+ *if: \$\{\{ !cancelled\(\) \}\}$`，**剩下 0 行**
（`git diff -U0 .github/workflows/ci.yml | grep -vE '^(\+\+\+|diff |index |--- |@@)' | grep -vcE '^\+ *if: \$\{\{ !cancelled\(\) \}\}$'` → `0`，`rc=1`＝grep 无命中，是**期望值**，不是失败）。
行数账：`992 → 1015`（`wc -l`，rc=0；差 23）。
行尾账：工作树与 HEAD blob 都是 **CRLF 0 / bare LF 1015|992**（python 数过，`rc=0`）⇒ 我的 23 行没有把文件写成混合行尾；`git` 那句 "LF will be replaced by CRLF" 只是 `core.autocrlf=true` 的既有告警，起手时就存在（不是我的改动带出来的）。
起手洁净：`git status --porcelain -- .github/workflows/ci.yml` 起手无输出（`logs/git-status-start.txt` 同件里 `ci.yml` 不在 31 条在飞改动里），⛔ 所以 diff 里的 23 行只能是我写的。

**YAML 语法自检（只解析，不执行）**：
`python -c "import yaml; d=yaml.safe_load(open('.github/workflows/ci.yml')); print(len(d['jobs']), steps)"` → `jobs 6 steps 52`，**rc=0**（`logs/final-checks.txt`；步数与改前一致＝我没加也没丢任何步）。
另：`roster.py` 的两条独立 pass（yaml 键扫描 ＋ 纯缩进扫描）在改前改后**都相等**（13/13、36/36），见 `logs/roster.txt`、`logs/roster-after.txt` 的 `== AGREE? ==` 行。

**在飞改动**：起手 `git status --porcelain` ＝ **773 行**（`?? 742`／` M 15`／` D 16`；非未跟踪的 31 条单列 `logs/git-status-start-tracked.txt`，全是 `.gitignore`、`.scratch/wisp/probes/{152,161,268}/**`、`design/**`）。⛔ 一条都没提交、一条都没还原。收尾 status 见 `logs/final-checks.txt`：我的两条路径之外**没有新增改动**。

---

## 5. 改完后各该步的本机等价命令颜色读数（会不会一上 CI 就红）

先给总结论：**预期零枚新红**。理由分两类：
- 23 枚里 **21 枚在过去三发 CI 上本来就是 `success`**（它们当时没被 skip，只是"随时可能被 skip"）。给它们守卫不改变它们自己的颜色。
- 只有 **2 枚**是"补完守卫后**从此开始求值**"的：`go vet (module)`、`go vet (tools/d22scan module)`（lint 9/10，历史上 3/3＝skipped）。这两枚我**真的在本机跑了**：

| 步 | 本机等价命令 | rc | 读数 |
|---|---|---|---|
| lint 9 | `CGO_ENABLED=1 go vet ./...`（根 module，go1.27.1 windows/amd64） | **0** | 绿；`logs/local-govet-root.txt`（空文件＝零诊断） |
| lint 10 | `cd tools/d22scan && go vet ./...` | **0** | 绿；`logs/local-govet-d22scan.txt`（空） |

交叉佐证（不是我测的，具名引用）：`ci.yml:127-129` 逐字写着 "`go vet (module)` is GREEN on ubuntu today (run 35600043583 step 7 = success; and in a real container `go vet ./...` is rc=0 with cgo on, rc=1 only for cmd/wisp with cgo off)"。两个方向（我的 windows 本机＋文件自录的 ubuntu 读数）都是绿 ⇒ 这两枚一上 CI **不该红**。
⚠ 诚实边界：本机是 windows/amd64，CI 那枚 job 是 ubuntu-latest；我的 rc=0 **不是** CI 读数，只是"它现在有没有明显红相"的旁证。真正的 CI 读数要等推送（⛔ 本程零 push）。

其余 21 枚的"会不会红"我**用的是历史 CI 真读数**（同一枚步的 success/failure），本机未复跑的理由逐条写明：

| 步 | 我跑了吗 | 颜色读数（来源） |
|---|---|---|
| lint 3 `bash tools/d22scan/runtests.sh -C tools/d22scan ./...` | ⛔ 未跑（`go test` 全量，预算；且它不写文件，可跑未跑） | success ×3（CI） |
| lint 4 `cd tools/d22scan && go run . -self-test` | ⛔ 未跑 | success ×3（CI） |
| lint 5 `sh scripts/d22scan.sh` | **跑了，rc=0** | 绿（本机）＋success ×3（CI）；`logs/d22scan.txt` |
| lint 6 `sh scripts/check-path-length-budget.sh --with-self-test` | **跑了，rc=0** | 绿（本机）＋success ×3（CI）；`logs/pathlen.txt` |
| lint 7 `gofmt (gofumpt)`（`go install mvdan.cc/gofumpt@latest` + `gofumpt -l . tools/d22scan tools/mockllm`） | ⛔ 未本机跑（要联网 go install；⛔ 不取镜像站哈希） | **failure ×3（CI，本来就红）**；票 111 `:394` 记了那枚红：`.scratch/wisp/probes/185/c1/mut/fs_broken.go:4:1` |
| test-core 3 `docker compose -f docker/compose.test.yml --profile test up -d` | ⛔ 未跑（本机未验 docker 可用性） | success ×3（CI） |
| test-core 4 `Probe mock-llm` | ⛔ 未跑（依赖上一枚起的 18080） | success ×3（CI） |
| test-core 5 `bash tools/d22scan/runtests.sh ./internal/proc/ -run Test…` | ⛔ 未跑 | success ×3（CI） |
| test-core 6 `bash scripts/portable-tests.sh --scope=core` | ⛔ 未跑（长跑；且它已在评估中，守卫不改它颜色） | **failure ×3（CI，本来就红）** |
| slo-smoke 3/4/5 cache + `Build wisp.exe` + `SLO smoke gate`（powershell scripts/…） | ⛔ 未跑（self 形状是 windows powershell＋真机采样；本机跑一次会污染 P3 的"有效样"账） | success ×3（CI） |
| lint-frontend 3–11 `npm ci` / `typecheck` / `lint` / `tokens:check` / `build` / 4 枚 render evidence | ⛔ **故意不跑**：`npm ci`/`npm run build`/`render:*` 会往 `frontend/**` 写（`node_modules`、`dist`、render 证据件），而 `frontend/**` 是我的禁区 | success ×3（CI，11 枚全绿） |

**新暴露的既有病：无。** 我没有把任何一步从"skipped"推成"failure"——唯一开始求值的两枚（`go vet` ×2）两向读数都绿。
**另：我没有修任何红。** `gofmt` 的红、`staticcheck` 的红（票 111 `:394`："staticcheck 12 发全红"）、test-core `--scope=core` 的红、test-windows 两枚 failure（`cmd/wisp CLI tests`、`Portable windows tests`）——**原样留着，具名交回给你裁**。

---

## 6. 顶回 / 我没量到的（逐条编号）

**T-1 起手锚号不符。** 你说 HEAD＝`e05b8a2b`，我实测 `6320a16a`（`182-c2 逐框六问…`）。我以现量为准。若 `e05b8a2b` 是你另开的一发，请具名指路；本程一切读数都锚在 `6320a16a`。

**T-2 词频尺的数你写错了：不是 25，是 27。** 你说"本项目实测：词频尺给 25 枚"。`grep -c 'cancelled' .github/workflows/ci.yml` 在 `6320a16a` 上＝**27**（`logs/rulers.txt` rc=0），票 111 `:390` 独立记的也是"**词频尺 27**（含注释，与 12 差 15）"。结构尺那三条我与你**逐枚对上**：步级 13＝12 `!cancelled()`＋1 `always()`、作业级 0。我没有拿你的 25 当基线。

**T-3 你的第 3 步"只补缺守卫且会被 skip 的步"如果按字面全补，会自己判红一枚钉子。** slo-full 的 3/4/5 满足你的判据（缺守卫＋结构上会被 skip），但 `scripts/slo-freshness.sh:160-167` 用 `^[[:space:]]+if:` 扫 slo-full 作业体，任何一行 `if:` 都 fail 那枚 `slo-full-trigger-missing`。我**一枚没补**，并复刻了那三条正则自证：`matched rc=1 / rc=1 / slo-check-present rc=0`（`logs/final-checks.txt`）。同一段还牵出 `ci.yml:229-233` 登记的"两枚 Upload 刻意不守卫"。这两处是**原文 vs 你的转述**冲突，我按原文做。

**T-4 `scripts/slo-freshness.sh:156-158` 那段自述已经过期**（"ci.yml carries 9 indented `if:` lines, 8 of them `${{ !cancelled() }}` and one `always()`, at :178 :215 :291 :351 :389 :399 :421 :457 :473"）——现量 13→36，行号也整体位移。它只是**注释**，不参与 P1 判定（P1 锚的是 slo-full 体内），所以我改不动也不该改（⛔ 非 `if:` 行＝越界；且 `scripts/**` 本来就不在我的两条路径里）。**这条要落账的话是你的活。**

**T-5 我推翻了 `ci.yml:226-228` 的一句正文，需要你确认。** 那句写着 "gofmt / go vet / staticcheck are production gates whose ORDER is itself a claim this file makes twice … **so they keep failing fast**"（111-r5 给它自己不去守卫那几枚的理由）。我照你的射程把 gofmt / go vet ×2 也补了。我的依据：①同一句里 staticcheck 其实**已经带守卫**（`:294`），"failing fast" 这条在本文件内并没有一致执行；②票 111 `:311`/`:321` 把那三枚 skipped 明写成**残余缺陷**（"今天仍是 skipped 且无日志，推前推后一模一样"），归票 236 AC#6／票 85 地界——不是设计目标。**但那段注释正文现在与文件实际状态矛盾**，我不改注释，交回给你。

**T-6 我取到的所有 CI 读数都在已推的 `cc315261` 上，不是在我本机 HEAD 上。** 差别是可点名的：票 111 `:390` 记"结构尺只 `!cancelled()` HEAD 12 / pushed tip 9"——pushed tip 上 `gofmt denominator` 的守卫**还不存在**（所以它三发都 skipped），lint 13 `Portable tests carrier self-test` 更是**在 156 行读数里一次都没出现**。⇒ "哪几步历史上被跳过"这份账是**上一版文件**的账；我推上去之后，那 23 枚的颜色**没有任何一发 run 能证明**。你说"改完先在本机跑那一步等价命令"正是为这一格补的，我按 §5 逐枚标了跑没跑。

**T-7 网络半途中断，两格〔待取数〕。** ①带 job 列的逐 step 读数重取＋`35600043583` 交叉验证：`gh` 连续 TLS 超时（`logs/step-conclusions-named.txt` 三条 `view rc=1` 为证）。②完整跑 `scripts/slo-freshness.sh` 取 P1/P2/P3 三探针合成结论：本机 `GITHUB_REPOSITORY` 设了之后仍 rc=2「no GH_TOKEN/GITHUB_TOKEN — the freshness probes cannot look (this is not a pass)」；gh 的 keyring token 不暴露成环境变量。我**没有**把 rc=2 读成"过"。
**T-8 runner 自动生成的 post 步（`Post Run actions/checkout@v4`、`Post Run actions/setup-go@v5`、`Post Cache third_party`）在 3/3 发里 `skipped`，我在 yaml 里没有载体可以给它守卫**——`if:` 只能挂在声明的步上。这一格是"你的判据覆盖不到"的形状，不是你写错，但按"判某门在 CI 有牙"的定式，它今天确实无牙。要不要记台账你定。
**T-9 票面的 A585 那类"逐票解冻"惯例：`ci.yml:137-143` 写着每次改这个文件都有一条 A### 授权（那次是"add one step to the lint job … and nothing else … not a step in slo-full / test-core / test-windows / slo-smoke"）。本程我动了 4 个 job 的 23 枚步，⛔ 我没有在 `docs/reports/pending-and-issues.md` 里找到对应本程的授权号可引（`docs/**` 与台账是我的禁区，我只读不写）。** 我以你的派单为授权，**但这一笔没有台账锚**，请你补登记。
**T-10 我没量：21 枚"结构上会被 skip"的步在过去若干发里到底有没有真被 skip 过。** 我只取了 3 发（票 111 用的是 12 发）。这直接影响 §2 那句"历史上真的从未被求值过的只有 3 枚"的强度——它是**3 发窗口**的结论，不是 12 发的。要扩到 12 发请再派一发。
