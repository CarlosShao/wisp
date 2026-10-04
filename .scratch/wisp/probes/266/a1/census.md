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

- 问① 今天有没有任何一把外部尺**真的读 `scripts/*.ps1` 的字节**？（`scripts/**` ＋ `.github/workflows/ci.yml` 全步名册，逐条答读不读 `slo-check.ps1` 的内容；谁最接近、差在哪一行）
- 问② `slo-freshness.sh` 的 P3 射程到底多大？（阈值常量名／默认天数／判产物还是判步的颜色／两形各响不响、时间维度给天数）
- 问③ 三形候选各的外部件放哪、凭什么发现不了"删检查也删声明"（最小写面逐枚具名文件／是否在被检文件之外／同动作可否一起绕过／走③ 会动到谁的分母）
- 问④ 词面型外牙的误报代价（"注释里写那三个词也打死门"这枚反方向读数意味着什么；仓里同源先例＝票 212 ban #9 的注释豁免 vs 字符串豁免、`emojiRe`/`shorthandRegionRe` 的豁免真值）
- 问⑤ 有没有别的门也是"自锁、外面无牙"（`scripts/**` 里其它自检形状名册＋各自有无外牙，⛔ 只登记不改）

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
