# 246-raisercensus-2 / 00 起手锚 + 停手条件核对（⇒ **不命中"已裁定"，继续跑第 1/2 步**）

腿名：`246-raisercensus-2`（只读普查续程腿）。工作语言中文。本件只落起手锚与停手条件核对结论，⛔ 不裁任何东西。

## 1. 现量锚（全部自己跑，未照抄派单）

| 尺 | 命令原文 | 现量 |
|---|---|---|
| 钟 | `date '+%Y-%m-%d %H:%M:%S %z'` | `2026-10-08 18:39:47 +0800` |
| HEAD | `git log -1 --format='%H %ad %s'` | `cb9f1663b17da79064616e518de9450033130981` ＝ `Thu Oct 8 18:39:28 2026 +0800` |
| 分支 | （由 commit 上下文现取）`dev` | 与派单一致 |
| 台账行数 | `wc -l docs/reports/pending-and-issues.md` | **14325**、rc=0 |

HEAD 提交消息逐字（长，原样抄）：

> A729 落账：收 ledger-audit-1（6538b3c1/d568ad8d/d892ecba，十行判语 9 对 1 偏）——唯一偏移＝A720 §3① 引 resident_task_source_windows.go:14，实测含 AskOnTaskRoot 那行是 :13（:14 是相邻的 askConfirmation 句）；我自己 sed 现验两处（同文件 :13/:14；另 This call is the caller 落 :249 而我一族引 :248）⇒ 定式（今日第三次同族·新形）：引文件:行号必须与被引的那句短语同一发现取，⛔ 不许把相邻行行号回填；⛔ 不回改 A720，就地打旧；它另给铁证 a04a095f 起未改、收件时点同为 :13 ⇒ 非漂移是我错位。九行对上里最有价值两条＝A721 的 13 枚污染源闭环到 ci-head.bin、A722 §2② 五禁词名册。§3 续派 246-raisercensus-2（前一枚被我含糊停手条件卡在第 0 步；新派单写死只有命中『已裁定』才停、⛔ 零 Go 命令因 181-v3 正种突变）§4 在飞＝comment-fix-prep-1／181-v3／246-raisercensus-2；⛔ 零 push

## 2. 在飞登记（只登记，⛔ 未读 `design/**`、`frontend/**`）

`git status --porcelain | head -20`（rc=0）前 20 行逐字：

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

另两把登记尺：
- `git status --porcelain | wc -l` ＝ **760**（rc=0；大数主要来自 `design/**` 的在飞改动，属别人名下，本程一行没读）；
- `git status --porcelain -- cmd internal`（rc=0）＝ **1 行**：` M internal/panel/workspace.go` ⇒ `cmd/**` 干净、`internal/` 仅此一枚（与 `181-v3` 正种突变相关的在飞面；⛔ 本程零 Go 命令，避免互洗读数）。

## 3. 停手条件核对：**不命中"已裁定" ⇒ 继续跑**

派单写死：只有命中"已裁定"（有人**明确裁过**"起管线算不算满足 `AC#7`"这一问）才停；"已登记/已收件/待裁"一律继续跑。逐处现取：

- **`A722 §3`（台账 :14168 起）＝登记待人裁**。逐字（`:14191` 一带）："**处置＝我⛔ 不回改票 246、⛔ 不撤那个勾、⛔ 不派腿顺手"修注释"**——这一格要的是**一枚能跑突变的非实现者裁决**（它要回答："起了管线就算满足 `AC#7`，还是必须真能举出一张卡才算？"），而那一发要种产码突变…**排在 253-r1 退出之后**，与 `259-v1` 同批排队。"
- **`A724 §3`（:14221 起）＝具名纠正"那不是裁定"**。逐字："而 `A722 §3` 是**登记**（"处置＝不回改不撤勾，登记成一枚**待人裁**的问题"）**不是裁定** ⇒ 它按字面停手是**对**的，但这一问**今天仍然未裁**。…⇒ `A722 §3` 那一问（起管线 vs 真举卡）**原样挂着**，归下一波裁决腿。"
- **`A725`（:14250 起）＝只收窄、不裁**。逐字节头："★**举卡那一跳量清了：有一条会举卡的产码路，但"今天会不会真举出一张"仍无凭据**——`A722 §3` 的措辞按此收窄（⛔ 不撤、⛔ 不翻框）"；文末："那一格要裁的从此是"起了管线算不算满足 `AC#7`"＋"`Prompt` 这一跳要不要行为凭据""。
- **`A729 §3`（:14306 起）＝本腿任务书**。逐字："把"常驻今天会不会真举出一张卡"这条链逐跳跃量清"。

⇒ 四处一致：**"起管线算不算满足 `AC#7`"这一问从未有人裁定，只被登记并收窄**。本腿**继续**跑第 1 步（逐跳调用形状尺）与第 2 步（`AC#7` 判据逐句核对）。

## 4. 本件边界

⛔ 零 Go 命令（`build`/`vet`/`test`/`list`/`doc`/`env` 全禁；`181-v3` 种突变中）；⛔ 未改产码/票面/台账；⛔ 未种突变；⛔ 未翻框；⛔ 未 `-done`；⛔ 未 push；⛔ 未新建票/判据。写面＝本目录 `.md` 只新建。
