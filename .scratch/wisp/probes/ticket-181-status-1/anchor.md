# ticket-181-status-1 · 起手锚（1/2 件：锚＋四把尺）

腿名：`ticket-181-status-1`（窄射程整档腿，只核票 181 `AC#7` 那一格的历史叙述是否过期）。
本件只登记锚数据与票面四把尺；线索复跑与结论在第二件。

## 0. 起手锚（全部现量）

- `date '+%Y-%m-%d %H:%M %z'` ⇒ `2026-10-08 16:15 +0800`
- `git log -1 --format='%H %ad %s' --date=format:'%Y-%m-%d %H:%M %z'` ⇒
  `7a452d0aca1bf5157a97fb5329ca6eb38b9f0694` ｜ `2026-10-08 16:10 +0800` ｜
  标题起手 `A721 落账（节头 16:0x＝追加前再跑一次 date＝16:06:41 现量…）：收 ci-if-eval-1（8e085ad4…）`
- `git branch --show-current` ⇒ `dev`
- 起手 HEAD（本程唯一比较基准）＝ `7a452d0aca1bf5157a97fb5329ca6eb38b9f0694`

### 别人的在飞（`git status --porcelain | head -20`，rc=0；⛔ 本程一枚都不动）

```
 M .gitignore
 M .scratch/wisp/probes/152/my152.py
 M .scratch/wisp/probes/161/r6/logs/flip-1.txt
 M .scratch/wisp/probes/161/r6/logs/flip-2.txt
 M .scratch/wisp/probes/161/r6/logs/flip-3.txt
 M .scratch/wisp/probes/161/r6/logs/flip-4.txt
 M .scratch/wisp/probes/161/r6/logs/flip-5.txt
 M .scratch/wisp/probes/161/r6/logs/flip-6.txt
 M .scratch/wisp/probes/161/r6/logs/flip-baseline.txt
 M .scratch/wisp/probes/161/r6/logs/flip-restored.txt
 M .scratch/wisp/probes/242/r3/logs/probe-routed.txt
 M .scratch/wisp/probes/268/v1/evidence.md
 D design/assets/base.css
 D design/assets/icons.js
 D design/assets/theme.js
 D design/assets/tokens.css
 M design/doubao/README.md
 M design/doubao/demo/app.js
 M design/doubao/demo/index.html
 M design/doubao/demo/styles.css
```

- 登记要点：在飞面里有 `design/**`（本程既不读也不写；`frontend/**` 同样不读不引）、
  `.gitignore` 与多枚他人证据件；本程写面只有票 181 那一枚 `.md` 的**追加**＋本腿自己的两枚新建 `.md`。
- ⚠ 一条写腿 `253-r1` 正在 `cmd/wisp` 上跑 ⇒ 本程**禁一切突变、禁一切 Go 命令**
  （`internal/panel`／`internal/tools` 任何源码改动都会进 `cmd/wisp` 编译图洗掉它的读数）。

## 1. 票 181 票面四把尺（现量，逐字命令）

票面文件（下称 `<T>`）＝
`.scratch/wisp/issues/181-nothing-in-go-ever-reads-git-so-the-panel-cannot-show-the-branch-and-the-two-mutating-actions-would-need-c17-whitelist-methods.md`

- 票名无 `-done` 后缀（`ls .scratch/wisp/issues/ | grep -i 181` 只回这一枚）⇒ 本票仍是未结案状态，
  不适用"已结案就停手上报"那一条。
- 尺①（行数）：`wc -l <T>` ⇒ `77`（rc=0）
- 尺②（未勾框）：`grep -cE '^[[:space:]]*- \[ \]' <T>` ⇒ `2`（rc=0）
  - 两枚分别是 `:33` **AC#6 门禁** 与 `:35` **AC#7（本程射程）**
- 尺③（已勾框）：`grep -cE '^[[:space:]]*- \[x\]' <T>` ⇒ `5`（rc=0）
  - 五枚＝`:28` AC#1 ／ `:29` AC#2 ／ `:30` AC#3 ／ `:31` AC#4 ／ `:32` AC#5
- 尺④（本程要动的那一行）：**票面真末尾追加，起于第 78 行**（现量末三行＝`:75` 新立 AC#7 那条子项、
  `:76` 「写面漏记（具名不补写…）」、`:77` 「工具坑（具名不美化）」；即末节是 `## 6` 那族历史叙述，
  本程新节挂在它之后，⛔ 不插进任何旧节中间）。
- 追加后的期望读数（自证尺）：行数 `77 → 77+N`（N＝我写的行数）、未勾框仍 `2`、已勾框仍 `5`。

### 尺④更正（同件末追加，16:5x 现量后写；⛔ 不改上面那两行原文）

上面那句「现量末三行＝… `:76` 写面漏记 … `:77` 工具坑」与「末节是 `## 6` 那族」**是我未现量的转述**，
第二把尺跑完即证伪：`grep -n '^## ' <T>` ⇒ 票面 5 枚节头且**无编号**＝`:7`／`:22`／`:26`／`:37`／`:43 Progress log`（rc=0）；
`tail -3` 现量真尾＝`:75` 新立 AC#7 子项、`:76` 纪律面（工具调用 31/30＝超 1 枚＋护栏偏离自报）、`:77` 它的证据件有一处过期。
⇒ 追加节最终写 `## 6.` 仅作可指认序号，并在节内具名声明票面原本无编号体系；本件保留上面那两行错话原样不删（只追加更正）。

## 2. AC#7 那一格的原文位置

`grep -nE '^[[:space:]]*- \[[ x]\]' <T>` ⇒ 框行号＝`28,29,30,31,32,33,35`（rc=0）。
本程射程＝**第 35 行**那一格（`- [ ]` **AC#7**），其原话断言（2026-09-28 写）：
「⚠⚠ **今天这枚"消费改写账户"在结构上是空转的**——`panel.WorkspaceView.Rewritten` 这枚字段
**没有任何生产者会把它填成 `true`**」。本程不改这一行一字（⛔ 不回改已提交的行），只在票尾核现状并具名打旧。
