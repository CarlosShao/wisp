# 普查件 `266-a1` — 为票 266 量齐「`scripts/slo-check.ps1` 两枚自锁钉有没有外牙」的料

取数时刻 `2026-10-04 14:0x +08`｜锚点 HEAD 见文末「取数锚点」节｜分支 `dev`
本腿＝只读普查：**零产码改动、零 AC 框改动、零裁形**。全部"能不能做到／会不会响"的结论**只到读码级**；需要真跑才知道的一律标〔仅读码，未跑〕。

⛔ 全程未执行任何 `go` 命令（`go test` / `go build` / `go vet` / `scripts/d22scan.sh` 一发都没跑），原因＝`256-v1`／`260-r5` 正在 `cmd/wisp`、`internal/agent/approval`、`internal/ball` 取数。

---

## §1 现量复认（对票面 §现量 1/2/3/4 逐条：对／错＋本腿自己的真读数）

打算答什么：四条逐条给「本腿自己重跑的尺＋逐字读数＋判定（对／错／部分对，具名更正）」。

### 现量 1 —— `tools/d22scan` 的 scope 名册里没有 `scripts/` ⇒ **判定：对（尺本腿自跑）**

- 尺（逐字）：`grep -c '"scripts' tools/d22scan/main.go` ⇒ 输出 `0`，`rc=1`（无命中）。票面写的"＝ **0**"与本腿读数一致；⚠ 补一句票面没写的：**`grep -c` 在无命中时 rc=1**，所以这条尺的"绿"长得像"错"，下一任别把 rc 当枚数。
- scope 名册真身（本腿现读，行号＝14:0x 取数时刻）：
  - `tools/d22scan/main.go:405-437` `declaredScopes()` —— 逐枚：`bans #1-5 internal/`（`:409`）、`bans #1-5 cmd/`（`:413`）、`ban #6 frontend/`（`:417`）、`ban #7 internal/tools/`（`:433`），末尾 `:436` 追加 `ban8Scopes(root)`。
  - `tools/d22scan/main.go:579-586` `emojiScopes()` —— `design/`（`:581`）、`frontend/`（`:582`）、`internal/`（`:583`）、`cmd/`（`:584`）。
  - 走的路径也只有这四棵树：`tools/d22scan/main.go:264-281`（`walkGo internal`、`walkGo cmd`、`walkText frontend`、`walkText internal/tools`、`walkEmoji` 循环）。⇒ **`scripts/` 既不在名册里，也没有任何一条 walk 指向它**，票面那句"任何 PowerShell 文件的形状不在任何一把既有扫描器射程内"成立。
- ⚠ **具名更正一处（更正的是 266 引用的上游说法，不是 266 本身）**：票 263 票面归口节把名册写成"internal/cmd/frontend/design/**tools**"（`.scratch/wisp/issues/263-slo-check-dies-on-unset-lastexitcode.md:53`）。真名册里**没有 `tools/`**：`declaredScopes()`／`emojiScopes()` 都不含它，且 `main.go:55-57` 逐字写着 ban #9 的射程 "over the production Go files of `internal/` and `cmd/` ONLY - `tools/**` and `_test.go` are outside this ban's range"。⇒ 票 266 §现量 1 的清单（internal/＋cmd/＋frontend/＋internal/tools/＋design/）**是对的**；`tools/` 那一枚是上游名册句的误列，本腿按"以盘上为准"记此差异。
- ⚠ 第二处补刀（对候选③ 直接有用）：`tools/d22scan/main.go:1312-1318` `isTextFile()` 的后缀白名单＝`.md .txt .html .css .js .ts .tsx .jsx .json .yaml .yml .toml .go .svg .vue` —— **`.ps1` 不在里面**。`walkText`（`:271`、`:274` 用的那把）过滤走 `isTextFile`，所以"把 `scripts/` 塞进现有 `walkText` 形"并不会真读到 `.ps1` 的字节；只有 `everyFile:true` 那条形（字段在 `main.go:290-299`，`ban #8 frontend/` 与 `ban #6` 用的形状）才逐文件不挑后缀。

### 现量 2 —— 删钉那一发 ⇒ rc=0 全绿、本机零外牙：**判定：读码级复认成立，⛔ 本腿未跑（禁 go／禁跑门本体）**

- 为什么"删掉即绿"在读码层面是**必然**而不是实测：两枚钉的**唯一在场证据**就是它们自己打印的两行 `ok`（`scripts/slo-check.ps1:230`、`:258`），而打印它们的语句与调用它们的语句（`:291-292`）**同在本文件内**；本仓没有任何外部件记这两枚钉的枚数（问① 的名册逐条给）。删除 `:291`／`:292` 中任意一行或整个函数，剩下的脚本仍走采样、仍写 `slo-report.json`、仍 rc=0 —— 〔仅读码，未跑〕。
- 票面 §现量 2 的原始凭据是 `263-v1` 的突变名册（`delete-nails` 一发）。⛔ 本腿不把它当凭据（票面 §现量 2 自己就这么规定），只当来路；本腿给的是上面那条**读码级机制复认**。
- ⚠ 一枚票面没写、但会改变"删钉之后剩什么"的形状：`scripts/slo-check.ps1:594`、`:596-598` 一带的 **leak 自检要求 `exit==1`**——票 263 票面 `Progress log`（`:52`）把它记为"名册外的第三牙"，并具名残余风险"哪天 leak 的形状不再要求 1，这一族就只剩可被就地删除的自锁钉"。⇒ 也就是说"删钉也删 wait"那一形今天还有一枚**语义上相邻、但不是为钉而存在**的牙会响；本腿读码级同意该说法（那几行判的是被采样程序的退出码是否为 1，不是钉的枚数），⛔ 未跑。

### 现量 3 —— 唯一外牙＝`slo-freshness.sh` 的 P3，按天钝：**判定：对，且比票面更钝（本腿补三枚射程洞）**

- 阈值常量（现读）：`scripts/slo-freshness.sh:97` `sample_max_age_days=${SLO_FULL_SAMPLE_MAX_AGE_DAYS:-3}` —— 常量名 `sample_max_age_days`，默认 **3 天**；判定式 `scripts/slo-freshness.sh:368` `if [ "$sample_age_days" -gt "$sample_max_age_days" ]`。
- 判的是**产物**不是步的颜色：`:282` `sample_artifact=slo-full-report`（逐字注释 `:275-281` 说明它为什么**不是** env 可覆盖的旋钮）；记录来源是 Actions 产物列表 API（`:310-315`，`select(.name=="slo-full-report")`、`select(.expired==false)`）。⇒ 票面 §现量 3 那句"它看的是产物不是钉的存在"成立。
- 时间维度要给天数（本腿按判定式算，⛔ 未跑）：`sample_age_days` 是整数除法 `age/86400` 向下取整（`:364`），比较是**严格大于** 3 ⇒ 需要 `age_days>=4` ⇒ 最早在产物诞生后**整 4 天（345600 秒）**的那一次读数才红。
- 补刀一：P3 之外还有 `slo-full-sample-never`（`:359-360`），但那一支只在**一枚产物都扫不到**时才走，删钉不影响它。
- 补刀二（票面没写，直接决定"钝"的下界）：这枚外牙**不在本机跑**。它唯一的自动入口是 `.github/workflows/slo-fresh.yml:46-52` 的 `schedule: cron '23 */6 * * *'` ＋ `runs-on: ubuntu-latest`（`:62`），而该文件 `:41-44` 逐字承认 "Like every other cron here, this schedule is only evaluated for the config on the **DEFAULT BRANCH**, so it is **inert until this file reaches main**"。当前分支是 `dev` ⇒ 〔仅读码，未跑〕按这段文字，dev 上那把 cron **根本不响**，只有 `workflow_dispatch`（人工）会跑它。
- 补刀三（同方向）：P3 需要 `gh` + token —— `scripts/slo-freshness.sh:195-213` `can_look()` 在无 `GITHUB_REPOSITORY`／无 `GH_TOKEN`/`GITHUB_TOKEN`／无 `gh` 时 **exit 2**（不是 0）。⇒ "本机有没有这枚外牙"的诚实读数：**本机没有，外牙住在托管 runner 上**；本腿这台机器拿不到它。

### 现量 4 —— 两枚钉今天的真身行号：**判定：对（本腿现读，逐枚行号吻合）**

| 票面主张 | 本腿现读 | 判定 |
|---|---|---|
| 词面钉 `:249`（`PSCommandPath` 空） | `scripts/slo-check.ps1:247` `$path = $PSCommandPath`、`:248` `if (-not $path) {`、`:249` `Fail-InstrumentBroken 'shape nail: this script does not know its own path (PSCommandPath empty) …'` | 对 |
| `:256`（命中即 `Fail-InstrumentBroken`） | `:255` `if ($hits.Count -gt 0) {`、`:256` `Fail-InstrumentBroken ('shape nail: {0} direct read(s) …')` | 对 |
| `:258`（ok 句） | `:258` `Write-Host 'slo-check.ps1: shape nail ok - this file reads no automatic exit-code variable (ticket 263)'` | 对 |
| 能力钉的探测在 `:223` 一形 | `:223` `$probe = Invoke-ExternalProgram -FilePath $comSpec -ArgumentList @('/c', 'exit 7')`；判据 `:224-229`；ok 句 `:230` | 对 |
| 两枚都在任何测量之前跑（`:286-287` 注释具名） | `:286-290` 那段注释逐字 "Both nails run before anything is measured (ticket 263)"；调用点 `:291` `Test-ScriptShapeNail`、`:292` `Test-ExitCodeInstrument` | 对（注释实际占 5 行 `:286-290`，票面只点了头两行，不是错） |

- 本腿自跑的枚数尺（逐字）：`grep -n 'LASTEXITCODE' scripts/slo-check.ps1 | wc -l` ⇒ **2**，命中行＝`:241`（注释，无 sigil）、`:252`（模式串，无 sigil）；带 sigil 的字面＝`grep -c -E '[$][{]?([A-Za-z]+:)?LASTEXITCODE'` ⇒ **0**、`rc=1`。⇒ 与票 263 AC#4 那格的凭据更正（"字面 2 行、带 sigil 0 处"）逐字一致，本腿独立复现。

## §2 五问逐答

打算答什么：每问一节，逐条 `file:line`，名册式，⛔ 不许只答"没有"。

### 问① 今天有没有任何一把外部尺真的读 `scripts/*.ps1` 的字节？⇒ **零枚。** 名册如下（⛔ 不是空答，最接近的两枚各差一行什么，写在表后）

尺的口径：`run:` 命令行名册＝`grep -nE "^\s+run: " .github/workflows/ci.yml | wc -l` ＝ **33**（命名步 `- name:` ＝ **37** 枚、`- uses:` ＝ **12** 枚、job 名册＝`^  [a-z-]+:$` 抽出 **6** 枚：`lint:65`／`test-core:309`／`test-windows:419`／`slo-smoke:564`／`slo-full:622`／`lint-frontend:696`）；本腿把 33 条逐条看过。`.github/workflows/` 只有 `ci.yml` 与 `slo-fresh.yml` 两枚文件（`ls` 实测）；PowerShell 分析器（PSScriptAnalyzer / Invoke-ScriptAnalyzer）＝**全仓零命中**（`grep -rni` on `.github scripts tools docs cmd internal`，空）。

| # | 步／件（`file:line`） | 它扫什么 | 读 `slo-check.ps1` 的**内容**吗 |
|---|---|---|---|
| 1 | `ci.yml:81` `bash tools/d22scan/runtests.sh -C tools/d22scan ./...` | Go 测试＋"PASS>0 且 SKIP=0"断言（`tools/d22scan/runtests.sh:86-108`） | 否 |
| 2 | `ci.yml:105` `go run . -self-test`（working-directory `tools/d22scan`） | 读 **`tools/d22scan/main.go` 自己**的字节抽名册（`selftest.go:177-241`，正则 `:117`/`:119`/`:125`） | 否 |
| 3 | `ci.yml:134` `sh scripts/d22scan.sh` → `main.go:264-281` 四棵 walk | `internal/` `cmd/` `frontend/` `internal/tools/` `design/`（`declaredScopes:405-437`、`emojiScopes:579-586`）；**`scripts/` 不在名册、也不在 walk**（§现量 1 复认） | 否 |
| 4 | `ci.yml:166` `sh scripts/check-path-length-budget.sh --with-self-test` | `git ls-files` 的**路径字符串**长度（`check-path-length-budget.sh:376-393`、`:463`、`:468`） | 否（见下方"最接近之二"） |
| 5 | `ci.yml:169-176` gofumpt／`:184-204` `attrib.sh --tracked-only` | 只 `.go`（后者"handed ZERO tracked .go files exits 2"） | 否 |
| 6 | `ci.yml:207`/`:210`/`:301` `go vet`、`:264` staticcheck | Go 源码 | 否 |
| 7 | `ci.yml:349`/`:559` runtests.sh 单包（`internal/proc`、`internal/risk`） | Go 测试输出 | 否 |
| 8 | `ci.yml:373`/`:543` `bash scripts/portable-tests.sh --scope=core|windows` | `go list` 解析出的包名册 vs 自己的 `core_pin`（GUARD A/B/C）；判的是 Go 包 | 否 |
| 9 | `ci.yml:471` `bash scripts/winsec-tests.sh` | Go 测试 stdout（`:123`、GUARD 3 `:156-168`） | 否 |
| 10 | `ci.yml:485`/`:582`/`:645` `powershell ... -File scripts/build.ps1 -Env dev`（＝票面说的"build.ps1 自己那步冒烟"） | **执行**，不断内容；build.ps1 自己的冒烟是 `wisp.exe doctor` 的 rc（`build.ps1:167-169`） | 否 |
| 11 | `ci.yml:507` `bash scripts/wisp-cli-tests.sh` | CLI 行为＋`deps.toml` 的 section 名（`:71`-`:76` 只在注释里提 fetch-deps.ps1） | 否 |
| 12 | `ci.yml:585`/`:648` `powershell ... slo-check.ps1 -Subset smoke|full` | **执行门本体**；步的 `run:` 只有这一行，没有任何外层断言看它的 stdout 里该有两行 `ok` | 否（执行≠读内容） |
| 13 | `ci.yml:587`/`:650` `Upload SLO report`（`name: slo-full-report`、`path: build/slo/slo-report.json`） | 产物文件是否存在（`if-no-files-found: warn` 在 `:690`） | 否 |
| 14 | `ci.yml:713-765` lint-frontend 全部步（npm ci/typecheck/oxlint/tokens:check/build/render:*） | `frontend/**` | 否 |
| 15 | `slo-fresh.yml:70` `sh scripts/slo-freshness.sh` | P1 读 `.github/workflows/ci.yml` 的字节（`slo-freshness.sh:134-172`）；P2/P3 读 API | **否**（见下方"最接近之一"） |
| 16 | `slo-fresh.yml:72-83` `sh -n scripts/slo-freshness.sh` ＋ `shellcheck -s sh scripts/slo-freshness.sh` | 只有 `slo-freshness.sh` **这一枚 .sh** 的语法与 lint；shellcheck 缺件即 rc=1（`:79-82`） | 否（且没有 PowerShell 对等件） |
| 17 | 本机件 `scripts/portable-tests-selftest.sh:90-130`、`:258-281` | **真读别的脚本的字节**并数枚数：awk 从 `scripts/portable-tests.sh` 抽 `core_pin`（`:123-124`）、`grep -c .` 计数（`:125`）、0 枚即 exit 2（`:127-130`）；同形抽 `winsec_pin`（`:261-266`） | 读的是 `.sh`，⛔ 不读任何 `.ps1`；且**没有任何自动触发**（ci.yml／slo-fresh.yml／其它脚本零调用，尺＝`grep -rn "portable-tests-selftest" .github scripts tools docs cmd internal`） |
| 18 | 本机件 `cmd/wisp/config_receipt_255_test.go:283` | walk 里**确实列了 `"scripts"`** 这一棵树 | ⛔ 不读：同一函数 `:296` 逐字 `if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") { return nil }` ⇒ 一枚 `.ps1` 都没打开 |
| 19 | 本机件 `internal/panel/composer_dispatch_test.go:605` | 名单里把 `scripts` 当**排除项** | 否 |

⇒ **最接近的两枚，各差在哪一行**：
1. **最接近＝`scripts/slo-freshness.sh:169`**：`if ! printf '%s\n' "$job_body" | grep -qE 'slo-check\.ps1[[:space:]]+-Subset[[:space:]]+full'; then fail ...`。它真的在断言一句**关于 `slo-check.ps1` 的在场**的话——但它 grep 的字节来自 `:160` 那枚 `sed -n '/^  slo-full:/,/^  [a-z-]*:$/p' "$ci_file"`，也就是**ci.yml 里的调用行**。差的那一行＝"被打开的文件是 `ci.yml`（`:99` `ci_file=${SLO_FRESH_CI:-$root/.github/workflows/ci.yml}`），从来不是 `scripts/slo-check.ps1`"。⇒ 它知道门**被叫到**，不知道门**里面有什么**；把 `:291-292` 两枚钉删掉，`ci.yml:648` 那一行原样还在，P1 照绿。
2. **次接近＝`scripts/check-path-length-budget.sh:376-393`**：它的分母**含** `scripts/slo-check.ps1` 这一枚路径（`git ls-files` 全量），但量的是 `length($0)`（路径字符数）与末段名长，⛔ 一次 `open()` 都不做。差的那一行＝没有 `cat`/`Get-Content` 之类的字节读取。⇒ 把门文件整个删掉它也不会红（路径消失只会让分母变小，方向是"更绿"）。
3. ⚠ 一枚**形状上完全正确但没牙床**的近亲＝名册 #17：`portable-tests-selftest.sh` 已经是"外部件读被检文件字节＋数枚数＋0 即 exit 2"的完整实现，**只差没有人自动跑它**。这一点对候选① 的代价估算直接有用（见 §2 问③）。
### 问② `slo-freshness.sh` 的 P3 射程到底多大？

- 阈值常量：`sample_max_age_days`，定义在 `scripts/slo-freshness.sh:97`，默认 **3 天**（`${SLO_FULL_SAMPLE_MAX_AGE_DAYS:-3}`）；判定式 `:368` `[ "$sample_age_days" -gt "$sample_max_age_days" ]`；年龄＝`:362-364` `sample_age=$((now - best_epoch))` 再 `sample_age_days=$((sample_age / 86400))`（整数向下取整）。姊妹枚 P2 同形：常量 `max_age_days`（`:96`，默认 3），判定式 `:260` 也是严格大于。
- 判的是**产物**还是**步的颜色**：产物。记录来源 `:310-315` 的 Actions artifacts 列表，过滤 `select(.name=="slo-full-report")` ＋ `select(.expired==false)`；名字是硬赋值不是旋钮（`:282`，理由逐字写在 `:275-281`）。⛔ 三处都不读 job conclusion（`:37-45` 明确"This is NOT the job conclusion"）。
- 两形各响不响（⛔ 本腿只给"能不能"，不给"会不会"；全部〔仅读码，未跑〕）：
  - **形 A：钉被删、门还能出结论** ⇒ P3 **不响**（"能响"都谈不上：新产物照旧上传，年龄被重置，`:368` 永远不成立）。同一次红也不会有：P1（`:140-172`）只看 ci.yml 的触发块与调用行，P2（`:244-265`）只看 job 记录。⇒ **三枚探针全绿，删钉这件事在这把尺上没有任何投影。**
  - **形 B：钉被删、门出不了结论** ⇒ 只有这一形会累积年龄，最早的红＝**产物年龄 ≥ 4 个整天**那一发（`age_days` 取整后要 `>3`，即 `age ≥ 345600s`）；再叠上读数机会的窗口：`slo-fresh.yml:51` `cron: '23 */6 * * *'` ⇒ 最坏在 4 天之后再多等 6 小时。⚠ 时间维度不许写成"3 天"：3 是阈值常数，**盘面上最早变红的下界是 4 天**。
  - 补一枚更钝的上界（票面没写）：`slo-fresh.yml:41-44` 逐字承认 cron "only evaluated for the config on the DEFAULT BRANCH, so it is inert until this file reaches main"。⇒ 在 `dev` 这一支上，形 B 的"多久响"读数是**不会自动响**，除非有人手动 `workflow_dispatch`。
- 诚实档（这把尺自己合格的证据，读码级）：拿不到就别判——`:195-213` `can_look()` 三种"看不了"全部 **exit 2**（无 repo／无 token／无 `gh`）；`:134-137` 没有 ci.yml ⇒ exit 2；`:177-181` 时间戳解析不了 ⇒ exit 2；`:310-320` API 失败 ⇒ exit 2。⇒ 这把尺**不是**安静返回绿的形状，票面 §现量 3 的"钝"是**射程钝**，不是**诚实度钝**。
- 一句话射程：**P3 射程＝"D32 连续 ≥4 天没有一份新报告产物"，且只在 main 分支的 cron 上自动评估；它对"钉的存在"零敏感度。**

### 问③ 三形候选：外部件放哪、最小写面、能否被同一动作绕过、动谁的分母

**候选① 独立尺读 `slo-check.ps1` 里两枚钉的调用点枚数**
- 最小写面（逐枚具名，两条路线）：
  - 路线 A（新建外部件）：`scripts/<新尺>.sh`（新，1 枚）＋ `.github/workflows/ci.yml`（加一步）。⚠ 加步需要**新的具名解冻**：`ci.yml:136-143` 那格 A585 的授权文字是"add **one** step to the lint job, suggested right after the tools/d22scan step"并且逐字列了"not the trigger table above, not the concurrency group, not any `runs-on`, not a step in slo-full / test-core / test-windows / slo-smoke"——第二枚步不在其内。
  - 路线 B（复用既有外部件）：`scripts/slo-freshness.sh` 加一枚 P1b（改 1 枚文件），它已经在读文件内容、已有 shellcheck 步（`slo-fresh.yml:78`/`:83`）、已有 CI 入口（`:70`）。**这是三形里最小写面的一形：1 枚文件。**
- 是否真在被检文件之外：是（新尺／`slo-freshness.sh` ≠ `slo-check.ps1`）。
- 同一动作（删检查也删声明）能不能一起绕过：**能。** 删 `slo-check.ps1:291-292` ＋ 把外部件里那条断言（连枚数常量）一起删掉，就是两枚文件各删一段，一个 commit。⛔ 关键递归（本腿认为这是本票最该被裁形时读到的一枚读数）：**路线 B 把问题搬了一层而没有关掉**——`slo-freshness.sh` 自己的"该在场"没有任何东西断言（P1 断的是 ci.yml，⛔ 没有谁断 P1b；`sh -n`/shellcheck 只断语法，删掉整段代码它们照样 rc=0）。任何"再加一枚脚本"的形都把同一个自锁复制到那枚新脚本上；本仓**唯一一处闭上的环**是 `tools/d22scan` 那族——它的"该有几枚 ban"从 `main.go` 现抽（`selftest.go:177-241`）、每枚必须有 ring+silent 成对样本（`selftest.go:548-554` auditSelfCases），**而且由另一枚已有的、CI 每跑都在执行的文件替它盯着**。
- 动谁的分母：无（不碰 d22scan，不碰 path-length 阈值）。若走路线 A，动的是 `ci.yml` 的步名册——那是 262-v1 正在读的同一棵树（票 266 §排程 已写互斥）。

**候选② 成功路径必须打印那两行 `ok`，由外层步断其在场**
- 最小写面：
  - 外层＝CI 步：`.github/workflows/ci.yml`——把 `:648` 那一步改成"跑＋`grep` 两行"＝**编辑现有步**，撞的是本仓反复自述的"ADDED STEP ONLY，no existing step was moved, edited, deleted, given an `if:`, or made continue-on-error"规矩（`ci.yml:103-104`、`:137-143` 两处逐字）⇒ 要么新的具名解冻，要么加第二枚"重跑一遍门"的步（⛔ 代价＝self-hosted 的 D32 门跑两遍）。
  - 外层＝另一枚脚本（⛔ 不碰 ci.yml）：1 枚新脚本，参照现成先例 `scripts/winsec-tests.sh:156-168`（`grep -qF` 断子进程 stdout 两句话在场，否则 rc=2）——本仓**已有这一形且就在 `scripts/` 里**。
- 是否在被检文件之外：是（断言住在 ci.yml 或另一枚脚本；被检文件只负责打印）。
- 同一动作能否一起绕过：**能，而且这形还多一枚假绿方向**——⚠ 打印 ≠ 检查：`slo-check.ps1` 可以留着 `:230`/`:258` 两句 `Write-Host`、把 `:224-229`/`:255-257` 的判定本体与 `:291-292` 的调用删空，两行照打、外层照绿。⇒ 若要真闭上环，外层断的不能是"字符串在场"而是"能力钉真起了一个进程／真读了一次文件"那类**产物或时间戳证据**，而那一旦要记枚数就落回候选①/③ 的记数问题。这一条是本腿给候选② 的核心代价读数，⛔ 不是裁形。
- 动谁的分母：无（不碰 d22scan；若改 `:648` 那一步则动 ci.yml 的步名册）。

**候选③ 把 `scripts/*.ps1` 纳入 `tools/d22scan` 的 scope 并补一条 ban**
- 最小写面（逐枚具名，**5 枚文件起步**，全部实测自读码）：
  1. `tools/d22scan/main.go`：`// Bans` 名册加一枚 tab-indented `<N> <tag>` 行（正则 `selftest.go:117` 要求 `^//\t([1-9][0-9]*)[ \t]+([a-z][a-z0-9-]+)`）＋检查函数（内含 `s.add("<tag>"` 字面站点，`selftest.go:119`）＋一条 walk（`main.go:264-281` 区）＋`declaredScopes()` 一枚 `live:true` 条目（`main.go:405-437`）。
  2. `tools/d22scan/selftest.go`：`selfCoverForTag`（`:519-530`）加映射，否则 `:508-514` 直接 `case "": return error`；fixture 种子（`:362-380` 那张 map，今天只有 go.mod／internal／cmd／frontend／allowlist）必须加一枚 `scripts/...` 文件——**不加就 `live:true` × 0 文件 → `main.go:1447-1452` exit 2**（每个 selftest case 都用同一份 fixture 跑）。
  3. `tools/d22scan/selftestsamples.go`：一枚 ring ＋ 一枚 silent 成对样本，缺任一枚 ⇒ `auditSelfCases` 报 HOLE ⇒ `runSelfTest` rc=2（`selftest.go:548-554`）。
  4. `tools/d22scan/scan_test.go`：分母／脚注钉要跟着重读——本腿命中的按名册生成的断言＝`:686`（真仓里不许有空 ban #8 scope）、`:695`（clean 那行必须逐枚点名它真扫过的 scope）、`:1090`（`scopeSummary(...)` 里 `ban #8 frontend/=%d` 的拼法）、`:2183`、`:793`、`:887`。⚠ `main.go:49` 逐字写着"the footer is **generated from the scope list**"，所以脚注文字会变，不是测试文字变不变的问题。
  5. `.github/workflows/ci.yml`：**不需要动**（`:83`/`:105`/`:134` 三步已经会执行到新的 walk 与新 ban）。⇒ 这是三形里唯一"外部件不新增执行入口"的一形。
- ⛔ 禁区确认：`tools/d22scan/allowlist.txt` 本腿一字未读内容以外的改动（只确认它在 `checkRoot`（`main.go:1331-1341`）与 fixture 种子（`selftest.go:373`）里是**必需件**）。走 ③ 时如果新 ban 带出既有命中，**没有豁免通道可写**——见下面这枚实测尺：`grep -c -E '[$][{]?([A-Za-z]+:)?LASTEXITCODE'` 逐枚＝`build.ps1` **4**、`fetch-deps.ps1` **2**、`sign-models.ps1` **3**、`dev/ball-cycle.ps1` **1**、`spike/run.ps1` **4**、`slo-check.ps1` **0** ⇒ 若 ban 的射程是"全 `scripts/**.ps1` 里不许直读自动退出码变量"，则**一发 14 处红，且 allowlist 是禁区**（本腿未跑 d22scan，读码级推论）。
- 是否在被检文件之外：**是，且外部件自己被另一枚已有的门盯着**（`-self-test` 的 roster 从 `main.go` 现抽；`undeclaredKeys`/`emptyLiveScope`/`driftedAbsentScope` 三道 guard 在 `main.go:1432`/`:1447`/`:1456`）。⇒ 三形里只有这一形把"删检查也删声明"变成**必须同时改 4 枚文件且不得漏枚**的动作：只删名册行 → `readRosters` 与 emitted tag 不一致 → `main.go`/`selftest.go` 的 accounting（`:226-236`）与 auditSelfCases 撞车；只删 walk 不删名册 → `undeclaredKeys`／`emptyLiveScope` 撞车。〔仅读码，未跑〕
- 动谁的分母（必须具名报给编排者，⛔ 本腿不加白名单）：
  - `declaredScopes()` 的**枚数**（今天 4 ＋ ban #8 的 4 ＝ 8 条 ledger 行，`main.go:405-437`/`:442-455`）⇒ self-report 的 scope 名册与"clean"那句逐字变。先例族＝`A207`／`5d463bb`（`main.go:30-37` 记的 37/40/43 三读不可互认，AGENTS.md §1.2 也引同一处）。
  - `main.go:55-59` 逐字写着"adding either tree is **a scope change for an owner to approve**, not a test fix"。⇒ ③ 天然是人工批准动作，不是实现腿可以自取的。
  - 第二层分母：`isTextFile()`（`main.go:1312-1318`）**不含 `.ps1`**，所以"加 scope"必须二选一——(a) 用 `everyFile:true` 那条 walk 形（字段与理由在 `main.go:290-299`）只对新 ban 生效；(b) 把 `.ps1` 加进 `isTextFile` 白名单，那会**同时扩大 `walkText` 的所有用者**（ban #6 `frontend/`、ban #7 `internal/tools/`，`main.go:271`/`:274`）的文件级射程 ⇒ 那是分母的第二次变化，票面 AC#4 说的"不许自行翻任何豁免"同样管它。〔仅读码，未跑〕
- 三形绕过成本对比（本腿读数，⛔ 不是裁形）：**③ 需要同时改的枚数最多（4 枚文件＋2 处 guard 语义，且外部件自己已有外牙），② 最少（1 枚文件）但有一个"留句删体"的假绿方向，① 居中（1–2 枚文件）且外部件自身无牙（问题搬一层）。**

### 问④ 词面型外牙的误报代价

- 那枚反方向读数的真身（读码级复认，⛔ 未跑 PowerShell）：`slo-check.ps1:251-252` 拼出的正则是 `(?i)` ＋ `[regex]::Escape([char]36)` ＋ `\{?(?:[A-Za-z]+:)?LASTEXITCODE`。两件事由此定死：
  1. sigil 是**必选**的（`Escape('$')` 给的是 `\$`，后面没有 `?`）⇒ 不带 sigil 的裸词不命中；本腿用等价尺复现：`grep -c -E '[$][{]?([A-Za-z]+:)?LASTEXITCODE' scripts/slo-check.ps1` ＝ **0**（`rc=1`），同文件另有 2 处裸词（`:241` 注释、`:252` 模式串）⇒ 与票 263 AC#4 那格"字面 2 行／带 sigil 0 处"逐字一致。
  2. 匹配面＝`Get-Content -Raw`（`:253`）的**整份文件字节**，⛔ 没有任何"注释豁免"的词法分类器 ⇒ **注释里写一句带 sigil 的字面，`:255` 命中、`:256` `Fail-InstrumentBroken`、rc=1，整道 D32 门当场死。**
- 判：**是，会。** 如果外牙也做成"扫字面、无分类器"，那么"给下一任留一句解释"就是有代价动作，而且代价已经在本文件落地成文字：`:236-243` 那段注释全部刻意写成**无 sigil** 拼法（`:241` 的 `Variable: LASTEXITCODE` 就是被付过费的那一行），票 263 的实现腿还专门把 `$global:` 那一形补进模式（`cb8dcecc` 那次 commit 的标题逐字："词面钉学会 scope 限定拼法（`$global:LASTEXITCODE` 那一枚反形是我自己量出来的）"）。⇒ 后果是**下一任只能用语义含糊的写法描述那个危险形状**，可读性被钉吃掉了。这就是"词面型外牙"的隐性税，裁形时应把它计进候选③（若 ban 做成扫字面）。
- 仓里可参照的同源先例（三枚，带 `file:line`；详见 §4）：
  1. `tools/d22scan/main.go:1065-1079`＋`:1150-1163`＋`:1168-1183`：ban #8 的"**注释豁免、字符串从严**"＝ Q-46(c)，是**人工签字**的口径（`main.go:1069-1076` 逐字："that is a signature, not an improvement ... what happened here is a human authorising it"），实现用的是**真词法分类器**（.go 走 `go/ast`，其它走 `--`/`/* */` 形状），豁免由 `commentRangesFor` ＋ `removeRanges` 两枚函数承担（调用点 `:1154`、`:1157-1159`，命中行 `:1160-1161`）。⇒ 先例的教训：**要豁免注释就写一枚分类器并把它钉进测试，别靠措辞。**
  2. `tools/d22scan/main.go:800-828`＋`:893-904`（`shorthandRegionRe`）＋`:913`／`:934`（豁免键在**区域起点**）：ban #9 的三种自由形状各自带一枚 expect-silent 样本——`selftestsamples.go:384`、`:396`、`:417`、`:429`、`:441`（ring 那枚在 `:374`）。⇒ 教训：**豁免的范围用正反成对样本钉住，一次收窄就会一次红。**
  3. `scripts/slo-check.ps1:251` 自己那枚"从 char code 36 现拼以免自匹配"＝**没有分类器时**的最低成本自免疫写法；它的价格就是上面那句"注释也打死门"。⇒ 三枚放在一起读：**d22scan 已经付过"词法 vs 词面"这笔学费两次**，新的外牙若走词面形，应显式声明它选的是第 3 种（无分类器）并承担"注释有代价"这一条，或者复用第 1/2 种的形状。
- ⚠ 本腿量不到的一枚：票面问的是"那两行豁免真值有几行"。本腿给不出干净的枚数——豁免横跨 `main.go` 三处（`walkEmoji` 内联 `:1154-1161`、`commentRangesFor` 起 `:1168-…`、`removeRanges`/`stripCommentLines` 另在两处），用行数当尺会把分类器本体与调用点混在一起数。已登记进 §6。

### 问⑤ 还有别的门是"自锁、外面无牙"吗？⇒ 有，6 枚；逐枚登记，⛔ 本腿不顺手改

| # | 门／形状（`file:line`） | 它检查的是谁 | 有没有外牙 |
|---|---|---|---|
| 1 | `scripts/slo-check.ps1:230`＋`:258`（两行 ok）／调用点 `:291-292` | 自己文件的字节＋自己的 instrument | ⛔ 无（本票对象；§现量 2/3 复认） |
| 2 | `scripts/slo-check.ps1:581-591` leak fixture self-test（要求被采样程序 `exit==1`） | 门自己的"能不能变红" | ⚠ 半枚：票 263 票面 `:52` 把它记为"名册外的第三牙"，但它判的是**被采样程序的退出码**，不是钉的在场；⛔ 无任何外部件知道这段该存在。〔仅读码，未跑〕 |
| 3 | `scripts/check-path-length-budget.sh:499-505` 正控读**自己**（`sed "/^${ROSTER_MARK_BEGIN}.../,/^...$/d" "$0" > "$BENCH/gate.sh"`，再 `grep -qE '^    R\["'` 断"删干净"） | 自己（用自抄的副本） | ⛔ 无：`ci.yml:136-166` 那一步若被删（A585 那次解冻的产物整格消失），`slo-freshness.sh` P1 只看 `on:`/`push:`/`cron:`/slo-full 的 `if:`/`continue-on-error:`/`slo-check.ps1 -Subset full`（`:140-172`），**不看这一步**；`tools/d22scan` 也不看 `.github/`。⇒ 本腿没跑那把尺（禁 go＋它是 262-v1 的读面），只登记形状 |
| 4 | `scripts/portable-tests.sh` 的 `core_pin`/GUARD A/B/C | 自己声明的枚数 vs `go list` 实解析 | ✅ 有一枚真外牙：`scripts/portable-tests-selftest.sh:123-130`（抽字节＋数枚数＋0 即 exit 2），`:261-266` 同形抽 `winsec_pin`。⚠ 但那枚载体**没有任何自动触发**（名册 #17 那把尺：ci.yml／slo-fresh.yml／其它脚本零调用）⇒ 记数对了，牙床缺一枚 |
| 5 | `scripts/winsec-tests.sh:156-168` GUARD 3（断子进程 stdout 两句在场） | 别人（`portable-tests.sh` 的输出） | ⚠ 它是**别人的**外牙；它自己那 12 行若被删，`ci.yml:471` 那一步仍绿（谁都不看它）⇒ 同族自锁 |
| 6 | `scripts/build.ps1:71`/`:100`/`:130`/`:169`（4 处直读自动退出码变量）＋`:115` 的 `-H=windowsgui`＋`:168` 的 `&` 不等 | 无人 | ⛔ 无：词面钉的射程由 `$PSCommandPath`（`slo-check.ps1:247`）**限定在它自己那一枚文件**，`build.ps1` 一次都不进它的视野。⇒ 票 263 的病在姊妹文件里原样活着；本腿只登记，不改，⛔ 不在本票射程 |
| 7 | `tools/d22scan/main.go` 的"// Bans 名册＋`s.add` 站点＋walk＋declaredScopes"四件套 | 自己，但由 `selftest.go` 反向抽 | ✅ 有，而且是本仓**唯一闭合**的一族：`ci.yml:81`（runtests.sh 正控）、`ci.yml:105`＋`ci.yml:83`（`-self-test` 两枚独立入口，理由逐字在 `ci.yml:92-96`："two independent ways to be run means deleting one is a reading, not a silence"）、`selftest.go:226-236`（每个 `s.add` 站点必须被认出来）、`:548-554`（缺成对样本即 rc=2） |

⇒ 净读数为编排者裁形用：**表里 #7 那一族是唯一"外部件的在场也被另一枚已有门自动检查"的形状**；#4 是"外部件存在、牙床缺失"；#1/#3/#5/#6 是"自锁且外面无牙"。本票要找的那枚外牙，盘上只有两种归宿：落在 #7 那族里（于是自动被盯），或者落在 #1/#3/#5 那种新文件里（于是把同一个洞复制一份）。

## §3 三形代价表

打算答什么：把 §2 问③ 的三形收成一张表——最小写面枚数／是否在被检文件之外记数／同一动作可否把它一起绕过／需要具名解冻谁。只给读数与枚数，⛔ 不裁形。

## §4 可参照先例（带 `file:line`）

打算答什么：仓里已经存在的"外部件替别人记数"的形状，逐枚给行号与它记的是什么（例：`slo-freshness.sh` P1 读 `ci.yml` 的 `if:`/`continue-on-error:`；d22scan 的 `-self-test` 从 `main.go` 自己抽名册再要求每枚 ban 有成对样本）。

## §5 别的「自锁无外牙」门名册

打算答什么：`scripts/**` 里所有"检查自己文件内容／打印自己的 ok"的形状，逐枚列 `file:line`＋有没有外部件知道它该在场。只登记。

## §6 判不动／量不到

打算答什么：两小节写满——(a) 本腿没跑的尺逐具名（⛔ 禁 go 的连带射程）；(b) 形判不动的地方，具名说"这里需要真跑，本腿只到读码级"。

## §7 记我自己写错的尺

打算答什么：本腿每一步用过的尺逐枚登记，写错的（枚数算错／根目录限定错／`head` 当枚数）**全数留下不删**，供下一任核。

---

## 取数锚点

- 本腿开工时刻：`2026-10-04 14:0x +08`；开工 HEAD＝`86440487`（`256-v1 骨架：验收表九节标题与每节打算答什么（非实现者腿，未跑产码改动）`，`git log -1 --date=iso` ＝ 2026-10-04 14:05:04 +0800）。
- 立票锚点（票面 §标题）＝`515ca5c5`，与本腿锚点不同代；⚠ 本件所有 `file:line` 一律取**本腿时刻**的盘上内容，`scripts/slo-check.ps1`／`scripts/slo-freshness.sh`／`tools/d22scan/main.go`／`.github/workflows/*.yml` 四棵在本腿全程 `git status --porcelain -- scripts tools .github` 为空（工作树＝HEAD，未被人动过），所以行号可对 `86440487` 复核。
