# 16x-c1 · 六枚契约级票（160／162／163／164／165／168）的冻结文本射程普查

- 程名＝`readonly-16x-c1`（只读普查）｜派单＝`.scratch/wisp/dispatches/2026-09-26-225x-readonly-16x-c1-contract-lines-census.md`
- 被派时刻 `2026-09-26 22:5x`｜**本程锚点（step 0 现量）＝`0d1764731973a998c5d6da5839070151fb9b04a9`**（`dev`）
- **本程只写两枚面**：本文件＋`.scratch/wisp/probes/16x/c1/**`。**不改任何票面、不改 `docs/PLAN.md`／`docs/specs/**`／任何 `.go`／任何 `.yml`。**
- 本文件**只提文字、不落文字**：下面所有"建议逐字改成"都是给编排者执行的草稿，**本程一个字都没往冻结件里写**。
- ⚠ 与前端那支会话**零往来**：§4.7 只答"我们自己的票面有没有要求写 `frontend/**`"，判据＝票面原文＋两把尺，**没有读过 `frontend/**` 的任何内容**（只出现过一次按文件名列举，见 §5 尺C）。

---

## 1. step-0 四件原文

```
$ date
Sat Sep 26 22:51:17 CST 2026

$ git rev-parse HEAD
0d1764731973a998c5d6da5839070151fb9b04a9

$ git rev-parse --abbrev-ref HEAD
dev

$ git status --porcelain | head -20        # 全量枚数见下条
 M .scratch/wisp/issues/169-a-single-decorative-glyph-in-frontend-makes-the-repo-wide-gate-red-and-kills-every-ci-step-after-it.md
 M .scratch/wisp/probes/152/my152.py
 M .scratch/wisp/probes/158/accept-r1/mut/guard.no1.go
 M .scratch/wisp/probes/158/accept-r1/mut/guard.no2.go
 M .scratch/wisp/probes/158/accept-r1/mut/guard.no3.go
 D design/assets/base.css
 D design/assets/icons.js
 D design/assets/theme.js
 D design/assets/tokens.css
 M design/doubao/README.md
 M design/doubao/demo/app.js
 M design/doubao/demo/index.html
 M design/doubao/demo/styles.css
 D design/index.html
 D design/screens/approval.html
 D design/screens/ball.html
 D design/screens/chat.html
 D design/screens/config.html
 D design/screens/cost.html
 D design/screens/firstrun.html

$ git status --porcelain | wc -l
49
```

**六枚票面都在（枚数＝6，与派单 §0 的第 4 件一致）**：

```
$ ls .scratch/wisp/issues/ | grep -E '^(160|162|163|164|165|168)-'
160-scope-open-returns-a-handle-carrying-its-own-closer.md
162-add-the-patch-one-region-tool-instead-of-rewriting-whole-files.md
163-persistent-shell-session-that-still-looks-one-shot.md
164-background-jobs-cannot-be-read-fix-the-roster-then-add-the-output-leg.md
165-plan-first-then-act-a-brake-the-voice-entry-needs.md
168-built-in-browser-as-a-disable-able-plugin-two-rules-to-set-first.md

$ ls .scratch/wisp/issues/ | grep -cE '^(160|162|163|164|165|168)-'
6
```

**版本口径（本程所有行号的前置声明）**：本程引的每一行都按上面的锚点取（`git show <锚>:<文件>`），
落盘件在 `.scratch/wisp/probes/16x/c1/`：

```
$ git show 0d17647:docs/PLAN.md > .scratch/wisp/probes/16x/c1/PLAN.md.0d17647
$ wc -l .scratch/wisp/probes/16x/c1/PLAN.md.0d17647
3595
$ cmp .scratch/wisp/probes/16x/c1/PLAN.md.0d17647 docs/PLAN.md && echo IDENTICAL
IDENTICAL                      # ⇒ 本锚点上 docs/PLAN.md 工作树＝blob，行号两把尺同一读数
$ git diff --stat HEAD -- docs/PLAN.md     # 空输出＋rc=0＝未改动
$ cmp .scratch/wisp/probes/16x/c1/provenance.go.0d17647 internal/risk/provenance.go
$ cmp .scratch/wisp/probes/16x/c1/bridge.go.0d17647 internal/tools/bridge.go
```

---

## 2. 本程**没**测什么（每一格＝只给了推断、没给凭据，逐条列）

| # | 格子 | 本程做到哪一步 | 没做到的那半 | 档位 |
|---|---|---|---|---|
| 2.1 | 六枚票的**AC 可达成性** | 只读票面＋读被票面点名的那几处码/文字 | **没有跑任何一发变异/用例**（本程是普查，未跑 `go test`、未跑门） | 〔我本轮现跑过的只有 grep/cmp/git show〕 |
| 2.2 | "改后不减少确认次数"这一判断（票 165 的代价栏） | 只核到 `internal/agent/approval/batch.go` 的**建造者注释与阈值常量**（D45-1 的聚合发生在**单次调用内**、且 L2 永不聚合） | **没有实跑一发**证"批完方案后仍逐次弹卡"；这一句是**从注释与常量推的** | 〔我未验，只是推断〕 |
| 2.3 | 票 160 的**跨包误关**那一发（AC#3 反向判据） | 只读了 `OpenScope`/`CloseScope`/`Detector` 的声明与两处生产调用点 | **没有构造那发变异**，也没量"摘掉认身份会怎样" | 〔我未验，只是推断〕 |
| 2.4 | 票 162 的 L1／L2 之争 | 量到 **R8 的 `overwrite` 是生产码里的结构化事实类**（`internal/risk/rules_irreversible.go:10,25`）＋`fs.write` 覆盖行现写 L2（`PLAN.md:2534`） | **没有量到"若真要 L1 需要什么样的可撤销机制"**——我给的"必须先有备份/撤销"是推理，不是读数 | 读数硬、推论那半〔我未验〕 |
| 2.5 | 票 163 的会话形状可行性（PTY／stdin／退出码标记） | 只量了我们侧的约束（ban #4 墙钟、`D33` argv 强制在 `PLAN.md:2507`、`internal/proc/**` 存在） | **外部三仓的原文（dsh 的 `tool-bash-persistent`、`node-pty` 那三条）我一字未读**——它们不在本仓，票面自标的档位我也没复算 | 〔仅指外部读数未取〕；本仓侧有凭据 |
| 2.6 | 票 164 的"名存实无" | **两把尺量过**：生产 `Name() string` 只有 6 枚 `fs.*`；`"task.` 字面量在生产码 0 命中（详见 §4.4 第 5 问） | 没有量"注册表运行时实际条目"（要跑起来才知道）——我引的是 A308 的读数＋我自己的静态两把尺 | 静态读数〔我本轮现跑过〕／运行时那半〔A308 读数，我未复算〕 |
| 2.7 | 票 168 的两张白名单**该定成什么内容** | 只量到"它们尚未定稿"这件事的**六处原文**（§4.6 第 2 问） | **没有提议任何一个具体方法名之外的内容**；我给的两枚方法名（`browser.view.open`/`browser.view.navigate`）是我造的名字，**仓里从未出现过**，属提案不属读数 | 提案（不是读数） |
| 2.8 | `D34` 声明枚数 | 我用两把尺都**没能复现出 A308 那句"36 枚"**（我按行切分得 34 个去重串，且其中 `list`/`cancel` 是被 `reminder.` 前缀折出来的碎片 ⇒ 我的尺子不合格） | **所以本文件不用"36"这个数**；我用的是可复现的**枚举行数 33／33**（命令在 §6） | 〔A308 的 36＝编排者读数，我未复算，且我这把尺证不了〕 |
| 2.9 | `frontend/**` 与 `design/**` 的实际内容 | **零读取**（只在 §5 尺C 里按**文件名**列了 8 枚命中 `shell.exec` 字面的件，未打开任何一枚） | 因此"面板/界面上会显示什么字"这一问，我只从**冻结件已写死的文案**回答（如 `PLAN.md:2025` 那句「点击悬浮球以批准」），**没有看任何现有界面** | 〔按派单 §4 令：故意不测〕 |
| 2.10 | 兄弟程 `161-r3` 的产出 | **零采信**：`tools/d22scan/**`、`ci.yml`、`probes/158/accept-r1/mut/**` 一律未进任何结论 | 无 | 〔按派单 §2 令：故意不测〕 |
| 2.11 | 同源拷贝普查里 `docs/evidence/**` 与 `docs/reports/**` | **主动排除在语料外**（它们是裁决史／历史读数，不是"自称逐字抄自"的契约拷贝） | 所以"某句在证据件里被引过 N 次"本程不答；那里头很多是**带版本的旧行号**，数进来只会污染这张表 | 〔故意不测，理由已核：尺A-1 语料三条〕 |

---

（§3 六行表与 §4 各节在最后几条 commit 里逐枚补上；本文件写到哪一格，commit 记到哪一格。）
