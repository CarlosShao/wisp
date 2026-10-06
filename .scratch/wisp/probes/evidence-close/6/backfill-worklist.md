# evidence-close-6 — 真缺凭据的 5 枚（78／115／250／254／267）逐枚补档工作单

> 本腿＝只读普查腿 `evidence-close-6`。唯一问题＝**每一枚要补哪一档、补它要什么形状的证据、能不能用盘上已有件顶**。
> ⛔ 本腿不重做 `evidence-close-5` 的 17 枚分诊（只读了它的**名册**（`5/c17-triage.md` §1，`:133-163`）与**归口表**（同件 §2.18，`:400-418`）两块；
> 它对票面的判语一枚未读，所以下面每一枚都是本腿自己看票。
> ⛔ 本腿零 `go` 命令、票面零改动、不翻勾、不改名、不动台账与 `docs/**` 一字；唯一写面＝本目录 `.scratch/wisp/probes/evidence-close/6/**`。
> 代号 `6` 起手时未用过（`ls .scratch/wisp/probes/evidence-close/`＝`1 2 3 4 5`），未启用 `6b`。

---

## §0 起手锚（三把尺原文＋取数时刻）

- `date '+%m-%d %H:%M'` ＝ **10-06 16:56**（起手）／落笔时复量 **10-06 17:05**；分支 `dev`。
- `git rev-parse --short HEAD` ＝ 起手 **`de084db7`** → 写件时 **`34ae7b05`** ⇒ **普查期间 HEAD 漂了一枚**（有腿在飞）。
  ⚠ 本件所有枚数一律"按落笔时刻引"，不复用别腿的读数。
- `git status --porcelain -- docs .scratch/wisp/probes/evidence-close` 原文（起手 16:56 与 17:05 两遍**逐字相同**）：

```
?? .scratch/wisp/probes/evidence-close/1/267-boundary-names.txt
?? .scratch/wisp/probes/evidence-close/1/gate-d22scan.log
?? .scratch/wisp/probes/evidence-close/1/gate-pathlen.log
?? .scratch/wisp/probes/evidence-close/1/msg-01-skeleton.txt
?? .scratch/wisp/probes/evidence-close/1/msg-02-sec1.txt
?? .scratch/wisp/probes/evidence-close/1/msg-03-sec2.txt
?? .scratch/wisp/probes/evidence-close/1/msg-04-sec3.txt
?? .scratch/wisp/probes/evidence-close/1/msg-05-sec4.txt
```

⇒ `docs/` 零行未提交（本腿对 `docs/**` 零写），evidence-close 下只有 `1/` 那批未跟踪件（别腿的 msg/log），本腿新建的 `6/` 在 17:05 那一遍之后才出现。

### 0.1 五枚票面文件真名（`ls .scratch/wisp/issues/<号>-*.md` 现量，全部已带 `-done`）

| 票 | 票面全路径 | 行数／字节 |
|---|---|---|
| 78 | `.scratch/wisp/issues/78-linux-vet-buildtags-unblocks-d22-gate-done.md` | 178／17,295 |
| 115 | `.scratch/wisp/issues/115-seal-notices-carry-the-resolvers-answer-while-the-cases-compare-caller-spelling-done.md` | 461／82,966 |
| 250 | `.scratch/wisp/issues/250-portable-tests-guard-c-merges-go-list-stderr-into-its-own-denominator-so-a-cold-module-cache-kills-the-whole-core-scope-reading-done.md` | 65／10,525 |
| 254 | `.scratch/wisp/issues/254-fifth-tier-winsec-no-ci-caller-done.md` | 47／7,779 |
| 267 | `.scratch/wisp/issues/267-confirm-timeout-config-has-no-bounds-so-the-c18-warning-can-vanish-done.md` | 87／24,037 |

改名（`-done`）那笔 commit 逐枚现量（`git log --diff-filter=R --all --name-status`）：
78＝`294f35d7`（09-21 12:18）／115＝`f5bbccd8`（09-22 23:09）／250＝`e39386b7`（10-02 09:27）／
254＝`21a15d27`（10-02 11:33，后又由 `46079fcc` 10-04 09:18 缩短 slug）／267＝`c6cf66e6`（10-05 10:58）。

### 0.2 本腿用到的尺（⛔ 不写进被扫文件）

- 框尺四形：A `grep -c '^- \[ \]'`（顶格未勾）／B `grep -cE '^[[:space:]]+- \[[ x]\]'`（**必须有**缩进的容缩进锚，派单原样）／
  C `grep -cE '^[[:space:]]*- \[[ x]\]'`（缩进**可选**的容缩进锚，能同时看见顶格）／D `grep -c '^- \[[ x]\]'`（顶格总框）。
- 凭据存在尺（只按**文件名**，⛔ 不读内容）：`find docs/evidence -name '<号>-*'`（＋`-iname '*<号>*'` 复量）、
  `find docs/evidence -mindepth 2 -name '*.md'`（非 s1 层）、`grep -rl '<号>-' .scratch/wisp/probes/`。
- 独立性尺：`grep -n 'agent=orchestrator' <票面>`（派单指定那把）。
- 界面签收尺：`grep -cE '签收|截图|owner 现场|肉眼|界面看' <票面>`。

### 0.3 ★负向结论前的"尺命中真名"控制（本腿逐条跑过，⛔ 不靠它下判语）

- 凭据名册尺：同一把 `find docs/evidence -name '<号>-*' -name '*.md'` 在**已知在盘**的号上给非零——
  `248→2`／`252→1`／`255→2`／`257→1`／`265→1`；本五枚给 `78→0`／`115→0`／`250→0`／`254→0`／`267→0`
  （`-iname '*<号>*'` 复量同为 0）⇒ **"名下零凭据件"是真负向，不是尺瞎**。
  同形控制另跑一组：`230→0`／`251→0`／`253→0`／`263→0`／`268→0`／`243→0`（这几枚"名下无表"是**已由别腿登记**的账，见 §4）。
- 独立性尺：`grep -c 'agent=orchestrator'` 在早期八枚上逐枚命中——`01→3`／`02→1`／`03→2`／`04→2`／`06→2`／`13→2`／`14→2`，
  全票池命中枚数＝**37 枚文件**（`grep -rl 'agent=orchestrator' .scratch/wisp/issues/ | wc -l`）⇒ 尺有牙；
  本五枚该尺**逐枚 0 命中**（读数见各枚⑤栏）。⚠ 但 0 命中**不等于**"非自裁"，理由写在 §3-6。
- 缩进框尺射程：全票池 `grep -rhE '^[[:space:]]*- \[[ x]\]'`＝**1,486** 行 vs `grep -rhE '^- \[[ x]\]'`＝**1,484** 行（差 2），
  未勾形 777 vs 776（差 1）⇒ **缩进里确实藏框**（藏在 `113-…-done.md`／`62-liquid-glass-ball-visuals.md`），
  本腿据把尺**改形**（C 形，缩进可选）逐枚复量，两形读数都记下（见各枚②栏）。

---

## §1 五枚工作单（六栏逐枚，⛔ 不合并）

### 1.1 票 78

**① 票号＋标题原文一行**（`.scratch/wisp/issues/78-linux-vet-buildtags-unblocks-d22-gate-done.md:1` 逐字）
> `# 78 — Fix the Linux-only `go vet` errors that have been silently skipping the D22 gate in CI since `fd8f838``

**② 它欠哪些格**（四把框尺逐枚现量，⛔ 本腿不挑一枚当准）

| 尺 | 写法 | 读数 |
|---|---|---|
| A 顶格未勾 | `grep -c '^- \[ \]'` | **0** |
| B 派单原样（**必须有**缩进＋容勾态） | `grep -cE '^[[:space:]]+- \[[ x]\]'` | **0** |
| C 缩进**可选**的容缩进锚 | `grep -cE '^[[:space:]]*- \[[ x]\]'` | **4** |
| D 顶格总框 | `grep -c '^- \[[ x]\]'` | **4** |

⇒ A 与 C 同数（4 枚框全在顶格、缩进里不藏框）；**B 那把尺读 0 是尺形问题不是票面问题**——它的 `[[:space:]]+` 要求至少一枚空白，顶格框天生不命中（§0.3 全池射程读数佐证：B 形在票 115 上也给 0，而 C 形给 5）。
框态逐枚（行号现量）：`:45 [x] AC#1`／`:59 [x] AC#2`／`:61 [x] AC#3`／`:64 [x] AC#4` ⇒ **物理零未勾**，
但**票面自己承认 AC#1 的判据只兑现了一半**：`:10-12` 与 `:48-58` 写"AC#1 那条 CI 侧 ubuntu 原生 `go vet (module)` 绿**当时仍无样本**"，
`:53` 那句逐字「**`GOOS=linux go vet ./...` 在本机上永远不可能 rc=0**，与本票修没修对无关」＋`:57`「**真正未见的判据**是 CI 上 ubuntu 原生那一步」。
⇒ 本票"欠"的不是勾，是**四格共同的、非实现者出的一张 1:1 裁决表**（标题 `:44` 逐字「## AC（1:1 裁决表）」）。

**③ 补一档要什么形状的证据**（逐格按票面判据分三类，⛔ 不预判任何一格的结论）

- **AC#1**（`:45-47`：`GOOS=linux go vet ./...` 干净＋"必须给出真实命令与 exit code"＋变异检验）
  ⇒ 要 **`场景或端到端记录`** 两发：(i) 一次 **CI 侧 ubuntu 原生** `go vet (module)` 步级读数（票面自陈本机那把仪器天生不匹配），
  (ii) 一发**退回旧形状 ⇒ 该命令转红**的突变记录（票面明写"先 grep 证变异落盘再跑，还原后证明确实还原"）。
- **AC#2**（`:59-60`：Windows 侧三包 `-count=2` 全绿＋`gofmt -l` 触及包为空）
  ⇒ 要 **`非实现者终裁表`** 一格（四数逐枚点名＋gofmt 读数）；⚠ 这一格**必须真跑 go** ⇒ 派单时要给安静窗口，别与在飞的 `236-r3c` 并发。
- **AC#3**（`:61-63`：装一条"无 tag 文件引用 windows-only 符号 ⇒ CI 必红"的机械检查，判据＝"它必须自己有一次真红，否则等于没装"）
  ⇒ 要 **`非实现者终裁表`** **且**要 **`场景或端到端记录`**：表判"这枚门禁在位且承载体是 `internal/proc/crossvet_test.go`"
  （本腿只按文件名核到该件**在盘**：`ls internal/proc/crossvet_test.go` 命中，⛔ 未读内容、未读它断什么）；
  读数那半＝**种子一发故意的引用 ⇒ 步骤真红**的复跑。
- **AC#4**（`:64-65`：D22 扫描步骤第一次产出真实 CI 结论，交回 run id＋job id＋该步结论；⛔ 排队被取代的 run 不算样本）
  ⇒ 要 **`场景或端到端记录`**（从 `gh api` 重取一发步级结论并区分 failure／cancelled(0 job)／未跑完）。
- **`owner 现场签收` 这一类＝本票零格**（尺：`grep -cE '签收|截图|owner 现场|肉眼|界面看'` 本票＝**0**）⇒ 无界面判据，不要为它派签收腿。

**④ 能不能用盘上已有件顶**（三把存在尺，只按**文件名**，⛔ 不读内容）
- 尺一 `find docs/evidence -name '78-*'`（＋`-iname '*78*'` 复量）＝**0 枚**；`ls docs/evidence/** | grep '^78-'`＝**0**。
  ⇒ 控制组同批发非零（`248→2`／`252→1`／`255→2`／`257→1`／`265→1`）⇒ **78 名下确实无表**（真负向）。
- 尺二 `grep -rl '78-' .scratch/wisp/probes/`＝**18 枚文件**（全是别的腿的 census/verdict 里**提到**这个号，⛔ 没有一枚是 `78-` 名下件）。
- ★ **"有读数但没落 78 名下"这一族确在盘**：本票那两枚 CI run 在**别的票名下的件里**逐枚在册——
  run `35558750456` 命中 `docs/evidence/s1/75-independent-verification.md`、`docs/reports/HANDOVER.md`、`docs/reports/pending-and-issues.md`、
  票 70／75／78／81／82 五枚票面；run `35562680354` 命中 `docs/evidence/s1/81-adversarial-acceptance.md`、`82-adversarial-acceptance.md`、`85-preflight-staticcheck.md`、台账与三枚票面。
  ⇒ 这与票 263 那形（`.scratch/wisp/probes/263/v1/verdict.md` 304 行、`:1` 署名"非实现者"、**曾是有效凭据只是没归口**）**不同形**：
  263 那枚是"给那张票自己出的表"，这三枚是"给票 75／81／82／85 出的表顺手含了同一发 run"。
  ⛔ **本腿不读它们的内容、不判断能不能顶**（那属另一枚愿意读内容的腿）；只把"搬运候选"具名摆出来。
  ⚠ 本腿越界自陈：起手为了定位，`head -1` 取了这四枚的**首行标题**（拿到"独立验收（acceptor，非实现者、非编排者）"／"编排者亲自复跑"／"独立对抗验收裁决表（acceptor-ticket82）"／"票 85 预做"四句署名），
  并因 `grep -n` 顺带看见 `docs/evidence/s1/closed-tickets-evidence-index.md`（**索引件**，不是裁决表）里 78／115／250／254／267 五行都写着 `NOFILE`。除这些名册级信息外一枚裁决表正文未读。

**⑤ 裁决者独立性初判**（尺＝`grep -n 'agent=orchestrator' <票面>`，⛔ 不读表内容）
- 本票读数＝**0 命中**（该尺在全池命中 37 枚文件、在早期八枚逐枚 1—3 命中，见 §0.3 ⇒ 不是尺瞎）。
- ⚠ **但 0 命中在这枚上不给"清白"**：本票的"自裁"形状是**另一种笔迹**——Progress log 三行全署 `agent=ticket78`（`:91`／`:107`／`:152`），
  而四枚勾的落笔人是**编排者**：`:3` Status 逐字「**done**（编排者验收 2026-09-21 13:2x）」、`:48`「编排者复跑后的口径更正…**不撤勾**、改判据读法」、
  `:67`「编排者回填（2026-09-21 13:1x）」（AC#4 那一格正是实现腿 `:66`／`:147` 明写"留着不勾、我推完回填"的那一格）、`:173`「编排者补一条刷新」。
  ⇒ **按 A 档那 7 枚同判读法：勾与判语出自同一方 ⇒ 即使盘上有表，也可能是自裁。** 本票连"表"都没有，只有编排者写在票面里的回填。

**⑥ 建议派单形状**（一句话）
**该立验收腿 `78-v1`**（非实现者，⛔ 本腿不代判），射程＝②③那两格真缺的仪器读数：一次**安静窗口**内复跑 AC#1（CI 侧步级读数＋退回旧形状那发突变）＋AC#2 三包四数＋AC#3 种子必红，产一张 `docs/evidence/s1/78-*.md` 1:1 表；
**同一轮另派一枚搬运笔**（只读、愿意读内容）判断 ④ 那三枚他票名下的件里含同一发 run 的行号**能不能顶 AC#4**，能顶就把行号指进 78 名下件（属归口不属补做）；
⛔ **不建议摘 `-done`**：四框物理全勾（A＝0/C＝4/D＝4）、Status 与 `-done` 一致，无"名实矛盾"那形的尺证据（那形是票 115 的 `:58`／`:64` 两格）。

### 1.2 票 115 —— 待写

### 1.3 票 250 —— 待写

### 1.4 票 254 —— 待写

### 1.5 票 267 —— 待写

---

## §2 治理账三把尺（只量读数，⛔ 不下"该不该算违规"的裁）

## §3 判得心虚的枚数与具名理由

## §4 判不动／要 owner／要搬运的清单
