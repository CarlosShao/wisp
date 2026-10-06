# evidence-close-5 — C 档 17 枚逐枚分流：凭据是"没有"还是"没归口"

- 腿：`evidence-close-5`（只读普查）。产件路径 `.scratch/wisp/probes/evidence-close/5/c17-triage.md`。
- 射程：只对前一枚索引腿（原件 `docs/evidence/s1/closed-tickets-evidence-index.md`）判为 **C 档**的票逐枚回答"缺凭据"还是"没归口"。
- 本腿零 go、零写已有文件、零翻勾、零改名、不读任何裁决表的内容面（只用 `wc -l -c` 量它存在与多长）。
- 不重判 A/B 档，不补写任何裁决表（写表必须是非实现者；本腿是只读腿，更不许替编排者签）。

---

## §0 起手锚（三把尺原文 + 取数时刻）

取数时刻：`date '+%m-%d %H:%M'` = **10-06 16:10**。

尺一 · 锚与脏项

```
$ git log -1 --format='%h %ad %s' --date=format:'%m-%d %H:%M'
d9aff5fd 10-06 15:25 probes(236-r3 死腿半成品代提)：一枚能编译、自己会红，但证件只有骨架——编排者只代提、不代判，标〔未验证半成品〕
$ date '+%m-%d %H:%M'
10-06 16:10
$ git status --porcelain -- docs cmd internal scripts .github
(空)
$ git rev-parse --abbrev-ref HEAD
dev
```

- 与派单给的现量比对：派单说"只有 `.scratch/wisp/probes/pool-validity/4e/**` 在动"，本腿这把尺**连那一处都没命中**
  （作用域是 `docs cmd internal scripts .github`，那批改动不在此作用域 ⇒ 零命中不矛盾）。
- `docs/evidence/s1/` **零脏项** ⇒ 不构成停手条件，本腿继续。

尺零 · 落点代号

```
$ ls .scratch/wisp/probes/evidence-close
1
2
3
4
$ ls -d .scratch/wisp/probes/evidence-close/5
ls: cannot access '.scratch/wisp/probes/evidence-close/5': No such file or directory
```

⇒ 代号 `5` 未用过、目录不存在 ⇒ **本腿用 `5`，不需要退到 `5b`**。目录由本腿 `mkdir -p` 新建。

尺二 · C 档名册自复量（不抄派单给的 17）

```
$ grep -nE '^\| [0-9]+ \| C \|' docs/evidence/s1/closed-tickets-evidence-index.md
157:| 01 | C | `NOFILE` |
158:| 02 | C | `NOFILE` |
159:| 03 | C | `NOFILE` |
160:| 04 | C | `NOFILE` |
162:| 06 | C | `NOFILE` |
168:| 13 | C | `NOFILE` |
169:| 14 | C | `NOFILE` |
172:| 19 | C | `NOFILE` |
181:| 78 | C | `NOFILE` |
206:| 115 | C | `NOFILE`（凭据在他名下，见 §3.3-3） |
241:| 243 | C | `NOFILE` |
242:| 250 | C | `NOFILE` |
243:| 251 | C | `NOFILE` |
244:| 254 | C | `NOFILE` |
246:| 263 | C | `NOFILE` |
248:| 267 | C | `NOFILE`（凭据在他名下，见 §3.3-3） |
249:| 268 | C | `NOFILE`（凭据在盘但不落 `s1/`，见 §3.3-3） |

$ grep -cE '^\| [0-9]+ \| C \|' docs/evidence/s1/closed-tickets-evidence-index.md
17
```

⇒ **本腿现量枚数 = 17，与派单给的 17 一致，名册零差**（票号 01 02 03 04 06 13 14 19 78 115 243 250 251 254 263 267 268）。

逐枚复量"名下几枚件"（两把尺：宽松／严格）：

```
$ for n in 01 02 03 04 06 13 14 19 78 115 243 250 251 254 263 267 268; do
    c1=$(ls docs/evidence/s1/ | grep -c "^${n}-")
    c2=$(ls docs/evidence/s1/ | grep -c "^0*${n}-")
    echo "$n s1files=$c1 strict=$c2"
  done
01 s1files=0 strict=0
02 s1files=0 strict=0
03 s1files=0 strict=0
04 s1files=0 strict=0
06 s1files=0 strict=0
13 s1files=0 strict=0
14 s1files=0 strict=0
19 s1files=0 strict=0
78 s1files=0 strict=0
115 s1files=0 strict=0
243 s1files=0 strict=0
250 s1files=0 strict=0
251 s1files=0 strict=0
254 s1files=0 strict=0
263 s1files=0 strict=0
267 s1files=0 strict=0
268 s1files=0 strict=0
```

⇒ 17 枚全部 **NOFILE 复量成立**（两把尺同为 0），与索引件 §3.2 逐枚一致。

负向结论的前置自证（本仓踩过"中文词面 grep 一份真名叫英文的册子、零命中就下结论"的坑）：
本腿任何"凭据不存在"的结论，落笔前先用同族命名证这把尺命中得了**真名**——

```
$ ls docs/evidence/s1/ | grep -c "^05-"    -> 1
$ ls docs/evidence/s1/ | grep -c "^265-"   -> 1
$ ls docs/evidence/s1/ | wc -l             -> 358
$ ls docs/evidence/s1/ | grep -iE 'readings'
111-ci-step-readings.md
111-ci-step-readings-attempt2.md
124-ac1-denominator-readings.md
134-ac4-machine-contended-readings.md
134-ac5-gate-readings.md
136-ac11-orchestrator-readings.md
136-ac11-second-witness-readings.md
ci-runner-readings-2026-09-21.md
ci-step-readings-2026-09-22.md
```

尺对同族命名（`<票号>-adversarial-acceptance.md`、`<票号>-*-readings.md`）命中正常 ⇒ 17 枚的 0 命中是"真没有"，不是"尺瞎了"。

尺三 · 索引件口径声明（只借名册形状，不采信结论）

- `:94`（§1.2 末条）逐字：**"C 档（零凭据）三把尺完全一致＝同那 17 枚"**，并说这条是选尺的"承重墙"。
- `:148`（§3.1 表）：C 档 = "名下完全无凭据（NOFILE＝零枚同名件）" = **17**，占比 18.3%。
- `:121`（口径声明）逐字：**"本尺只认文件名前缀，不认'内容里提到了哪枚票'"**，并自陈这会把"有非实现者凭据、但凭据不在自己名下"的票判进 C。
- `:253` 起 §3.3 第 3 条：17 枚里"5 个具名形状（涉及 8 枚票）"的凭据真存在、只是不落在那个名字＋落点形状里。
- 本腿**不采信**这些结论，只按它给的 17 枚票号逐枚自己现量；本腿独立发现的落点形状比它那 5 个更多（见 §2、§3）。

---

## §1 C 档名册（本腿现量）

枚数：**17**（§0 尺二复量，与索引件／派单三方零差）。票面文件全路径与结案态：

| # | 票 | 票面文件全路径（`.scratch/wisp/issues/`） | 未勾框 / 总框 / AC 记号枚数 |
|---|---|---|---|
| 1 | 01 | `01-build-chain-done.md` | 0 / 5 / 0 |
| 2 | 02 | `02-s0-spike-done.md` | 0 / 5 / 0 |
| 3 | 03 | `03-skeleton-runtime-rules-done.md` | 0 / 6 / 0 |
| 4 | 04 | `04-sqlite-core-done.md` | 0 / 6 / 0 |
| 5 | 06 | `06-secretstore-envs-done.md` | 0 / 6 / 0 |
| 6 | 13 | `13-audio-capture-done.md` | 0 / 6 / 0 |
| 7 | 14 | `14-model-distribution-done.md` | 0 / 6 / 0 |
| 8 | 19 | `19-provenance-c25-done.md` | 0 / 5 / 0 |
| 9 | 78 | `78-linux-vet-buildtags-unblocks-d22-gate-done.md` | 0 / 4 / 9 |
| 10 | 115 | `115-seal-notices-carry-the-resolvers-answer-while-the-cases-compare-caller-spelling-done.md` | **2** / 5 / 12 |
| 11 | 243 | `243-a-closed-ticket-number-still-owns-work-in-25-comments-which-is-the-shape-that-fooled-three-legs-done.md` | 0 / 3 / 4 |
| 12 | 250 | `250-portable-tests-guard-c-merges-go-list-stderr-into-its-own-denominator-so-a-cold-module-cache-kills-the-whole-core-scope-reading-done.md` | 0 / 4 / 8 |
| 13 | 251 | `251-winsec-pin-is-a-pin-nobody-compares-the-winsec-scope-runs-through-explicit-paths-so-guard-c-never-sees-it-done.md` | 0 / 4 / 5 |
| 14 | 254 | `254-fifth-tier-winsec-no-ci-caller-done.md` | 0 / 3 / 5 |
| 15 | 263 | `263-slo-check-dies-on-unset-lastexitcode-done.md` | 0 / 7 / 8 |
| 16 | 267 | `267-confirm-timeout-config-has-no-bounds-so-the-c18-warning-can-vanish-done.md` | 0 / 5 / 8 |
| 17 | 268 | `268-resident-leg-folds-missing-config-and-refused-config-into-one-provenance-done.md` | 0 / 5 / 6 |

未勾框总账：顶格尺 `grep -c '^- \[ \]'` 与任意缩进尺 `grep -cE '^[[:space:]]*- \[ \]'` 对 17 枚逐枚同数（缩进里不藏框）
⇒ **只有票 115 一枚还剩物理未勾框**，其余 16 枚全勾。票 115 那两格与索引件 §4.1 第 13/14 行
（原件 `docs/evidence/s1/closed-tickets-evidence-index.md` 自 `:304` 起那一节）**独立复现一致**
⇒ 本腿不新增该线索，也不翻框。

---

## §2 十七枚逐枚三分表

形状＝①在等哪种凭据 ②现量尺 + 读数原文 ③归口结论。

### 2.1 票 01 — `.scratch/wisp/issues/01-build-chain-done.md`

① 在等：**终裁表**（判 AC 框那一类；票面 5 框全勾）。
② 现量：
- `ls docs/evidence/s1/ | grep -c "^01-"` = `0`；`grep -c "^0*01-"` = `0`。
- 本腿新加的跨目录同名尺 `find docs/evidence -type f | grep -E '/(0*01)-'` → **命中 `docs/evidence/s0/01-adversarial-acceptance.md`**。
- `wc -l -c` 该件 = **167 行 / 14,227 字节**；`git log -1 -- …` = `98b5a9ed 09-19 14:53 docs(evidence): T01 adversarial accepta…`；`git ls-files … | wc -l` = `1`（已跟踪）。
- 票面引用它的行：`:49`（"T01-adv VERDICT PASS … report docs/evidence/s0/01-adversarial-acceptance.md"）、`:50`（"AC boxes reconciled against docs/evidence/s0/01-adversarial-acceptance.md"）。
- `Status:`（`:3`）逐字 `**Status:** done`；未勾 0；最后一次动这笔票 = `c992da08 09-20 10:49 chore(tickets): reconcile AC checkboxes against docs/e`。
③ 归口：**凭据在别处（指路径）＝错档**。`docs/evidence/s0/01-adversarial-acceptance.md` 以票号开头、命名与 A 档件同族、
且被票面逐字当成翻勾凭据引了两处 ⇒ 真凭据存在，只是落在 **`s0/` 那一层**，被"只扫 `docs/evidence/s1/` 一层"的索引尺看不见。
要做的不是补凭据，是**把尺的落点扩到 `docs/evidence/{s0,s2,s3}/**`**（改不改档＝编排者判断，本腿不改）。
那张表是谁签的本腿不查（查＝读内容面，被禁）。

### 2.2 票 02 — `.scratch/wisp/issues/02-s0-spike-done.md`

① 终裁表。
② s1 名下两把尺均 `0`；`find` 命中 **`docs/evidence/s0/02-adversarial-acceptance.md`（136 行 / 10,116 字节）＋ `docs/evidence/s0/02-spike-report.md`**；
提交 `f11f7063 09-19 17:29 docs(evidence): T02 adversarial accepta…`；票面引用行 `:58`、`:59`；`Status: done`；未勾 0；最后提交 `3cbde3a6 09-20 10:56 docs(reports): audit-B MINOR cleanups…`。
③ **凭据在别处＝错档**（同 2.1，落 `s0/`）。

### 2.3 票 03 — `.scratch/wisp/issues/03-skeleton-runtime-rules-done.md`

① 终裁表。
② s1 名下 `0`/`0`；命中 **`docs/evidence/s0/03-adversarial-acceptance.md`（166 行 / 13,033 字节）**，提交 `6b83e93e 09-19 15:58 docs(evidence): T03 adversarial accepta…`；票面引用 `:64`、`:65`；`Status: done`；未勾 0；最后提交 `c992da08 09-20 10:49`。
③ **凭据在别处＝错档**。

### 2.4 票 04 — `.scratch/wisp/issues/04-sqlite-core-done.md`

① 终裁表。
② s1 名下 `0`/`0`；命中 **`docs/evidence/s0/04-adversarial-acceptance.md`（223 行 / 16,572 字节）**，提交 `642eee8a 09-19 17:29 docs(evidence): T04 adversarial accepta…`；票面引用 `:56`；`Status: done`；未勾 0；最后提交 `3cbde3a6 09-20 10:56`。
③ **凭据在别处＝错档**。

### 2.5 票 06 — `.scratch/wisp/issues/06-secretstore-envs-done.md`

① 终裁表。
② 现量：
- s1 名下两把尺均 `0`；跨目录尺命中 **`docs/evidence/s0/06-adversarial-acceptance.md`** = `wc -l -c` → **52 行 / 3,928 字节**（17 枚候选里最短的一件）。
- 该件最后一笔提交 = `8dbd9a2a 09-20 11:25 chore(tickets): addendum AC ruling for ticket 06 (AC#6…`＝**补裁那一笔**（不是首裁那一笔）。
- 票面引用行：`:47`（逐字 "note: AC#6 **PASS** — 见 docs/evidence/s0/06-adversarial-acceptance.md §"Addendum 裁决（2026-09-20，AC#6 补裁）""）、`:62`、`:63`；`:64` 另记 `agent=ac-addendum-auditor did=补裁 AC#6（该框原为空白裁决）`。
- `Status:`（`:3`）= `**Status:** done`；未勾 0；最后一次动这笔票＝`8dbd9a2a 09-20 11:25`（＝与那笔补裁同一枚提交）。
③ 归口：**凭据在别处＝错档**（落 `s0/`）。★附一条尺寸告警（只有尺、没有判语）：件体 52 行／3,928 字节而票面 6 框全勾；
票面 `:47` 与 `:64` 两处都写明 AC#6 那一格**原本"空白裁决"、由补裁笔事后补上**
⇒ 若有人拿这一件当"06 的整票终裁表"，它对六格的覆盖度要由**非实现者另量**。本腿不读内容面、不判覆盖度。

### 2.6 票 13 — `.scratch/wisp/issues/13-audio-capture-done.md`

① 终裁表。
② s1 名下 `0`/`0`；跨目录尺命中 **`docs/evidence/s2/13-adversarial-acceptance.md`（102 行 / 9,820 字节）**，提交 `e45b9672 09-19 23:09 docs(evidence): T13 adversarial accepta…`；
票面引用 `:61`（"T13-adv VERDICT PASS … report docs/evidence/s2/13-adversarial-acceptance.md"）、`:62`（"AC boxes reconciled against …"）；`Status: done`；未勾 0；最后提交 `8bf0d341 09-20 10:49`。
③ **凭据在别处＝错档**（落 `s2/`）。

### 2.7 票 14 — `.scratch/wisp/issues/14-model-distribution-done.md`

① 终裁表。
② s1 名下 `0`/`0`；命中 **`docs/evidence/s2/14-adversarial-acceptance.md`（85 行 / 8,976 字节）**，最后一笔提交 `f5362d5d 09-20 11:26 chore(tickets): addendum AC rulings for ticket 14 (AC#…`＝补裁那一笔；
票面引用 `:50`（AC#2 补裁行，逐字"见 docs/evidence/s2/14-adversarial-acceptance.md §"Addendum 裁决（2026-09-20，AC#2 补裁）""）、`:69`、`:70`；`:71` 记 `agent=ac-addendum-auditor did=补裁 AC#2/AC#3/AC#4/AC#5（四框原为空白裁决）`；`Status: done`；未勾 0。
③ **凭据在别处＝错档**（落 `s2/`）。2.5 那条覆盖度告警在这里**更硬**：票面 `:71` 自陈**四格**原为空白裁决、由补裁笔改勾，
而整件只有 85 行／8,976 字节 ⇒ 覆盖度需非实现者另量，本腿不代量。

### 2.8 票 19 — `.scratch/wisp/issues/19-provenance-c25-done.md`

① 终裁表。
② 现量：
- s1 名下 `0`/`0`；跨目录尺命中三件：**`docs/evidence/s3/19-adversarial-acceptance.md`（353 行 / 51,596 字节，提交 `a1ac6213 09-20 23:12 ci(A26,A27,A29): CI 5/5 job 全红实…`）**、
  `docs/evidence/s3/19-adversarial-fix-round.txt`、`docs/evidence/s3/19-provenance-c25-tests.txt`（提交 `c33d10b4 09-20 10:53 docs(risk): C25 P12 conclusion + self-h…`）。
- 票面 `:52` 逐字把第三件当证据引（"evidence: docs/evidence/s3/19-provenance-c25-tests.txt"）；`:63` 记 "accepted FAIL verdict, starting must-fix set B-1/B-2/M-1/M-3/M-6"。
- `Status: done`；未勾 0；最后一次动这笔票 = `b3f6ea74 09-20 14:31 chore(tickets): ticket 19 DONE (orchestrator acceptanc…`。
③ **凭据在别处＝错档**（落 `s3/`）。★但这枚有一个**必须由另腿看**的时序形状（本腿只登记、不判）：
票 19 改名结案那笔提交在 **09-20 14:31**，而 `docs/evidence/s3/19-adversarial-acceptance.md` 的落地提交在 **09-20 23:12**
⇒ **那张表比这笔票的结案晚约 8 小时 41 分**（两笔提交号与时刻本腿已贴原文）。判它要读表的内容面，本腿被禁。

### 2.9 票 78 — `.scratch/wisp/issues/78-linux-vet-buildtags-unblocks-d22-gate-done.md`

① 在等：**终裁表**。票面 AC 段标题逐字 "## AC（1:1 裁决表）" ⇒ 等的就是名下一张 1:1 表。
② 现量（这枚的负向结论做了三重自证）：
- `ls docs/evidence/s1/ | grep -c "^78-"` = `0`；`find docs/evidence -type f -name "78-*"` = **零命中**；`find .scratch/wisp/probes/78 -type f` = **无目录（0 枚）**。
- 票面自报的凭据是 CI 原文而不是表：`:3` 逐字 `**Status:** **done**（编排者验收 2026-09-21 13:2x）—— AC#1..AC#4 四框全绿`；
  `:4` run `35558750456` 里两步 success；`:17` 逐字 `**Evidence:** registry **A44②**；run 35551819606 / job 106188167868；票 71 报告`。
- 那枚 run 的日志**不在盘上**：`ls .scratch/ci-logs/ | grep -E "35558750456|37166458550"` 只命中 `run-37166458550-failed.log` / `.err`（票 263 那枚的），
  78 引的那一枚**没有**（`.scratch/ci-logs` 整层共 6 枚）⇒ 连"当时贴的读数"现在也不可复跑。
- 同族命名的正向对照见 §0 尺二末段（`readings` 那一族 9 枚命中正常）⇒ 尺认得这类件名，78 名下确实零枚。
- 反向引用尺（只用 `-l`、不读内容）：`grep -l -F "35558750456" docs/evidence/s1/*.md` → 唯一命中 **`docs/evidence/s1/75-independent-verification.md`**。
- 票面还留着一格它自己承认无样本（`:55` 附近逐字）："**真正未见的判据**是 CI 上 ubuntu 原生的那一步 `go vet (module)`：run `35558750456` 里它 **skipped**…所以**AC#1 的"CI 侧绿"仍无样本**"。
③ 归口：**缺凭据（要补）**，且要补的形状本腿能给准——不是重新发明证据，而是**把已存在的那枚 run 读数正式落成 78 名下的一件**
（`docs/evidence/s1/78-*.md`）。这与票 115 那族"证据在别处、就是没有名下的 1:1 表"同形。
★两条必须点名的残余：(a) AC#1 的"CI 侧 ubuntu 原生绿"那一半**票面自己写"仍无样本"**，那不是归口问题、是缺读数；
(b) `75-independent-verification.md` 反向提到同一 run，要不要据此把 78 改记成"凭据在他名下"，属编排者／裁决腿的判断，见 §3-5。

### 2.10 票 115 — `.scratch/wisp/issues/115-seal-notices-carry-the-resolvers-answer-while-the-cases-compare-caller-spelling-done.md`

① 在等：**终裁表**——而且是"两半都要"：既缺 115 名下的 1:1 表，**AC#4（变异两向）与 AC#6（门禁按包四数）这两格连读数都没有**。
② 现量：
- `ls docs/evidence/s1/ | grep -c "^115-"` = `0`（严格尺同 0）；`find .scratch/wisp/probes/115 -type f` = **无目录**。
- 未勾框原文（顶格尺命中 2 枚，缩进尺同数）：
  - `:58:- [ ] **AC#4** 变异：把你选的修法退回原状 ⇒ AC#3 新用例必须红；再试一发"只比大小写不敏感"（`EqualFold`）这种**半修** ⇒ 也要红`
  - `:64:- [ ] **AC#6** 门禁：按包 `-count=2 -v` 四数逐条点名（报 SKIP 要说是不是 `-v` 量的；`-count=2` 不缓存）；`
- `Status:`（`:3`）逐字开头 = `**Status:** BLOCKED 部分（2026-09-21 21:2x agent-ticket115：AC#2 已裁=**方向 B（比对按树不按拼写）**…）`；
  全票 `grep -nE 'WITHDRAWN|作废|撤'` 只命中 `:347`（"票面点名的 run 已作废"——说的是那枚 run，**不是这张票**）。
- 票面自己指的凭据文件名：`grep -oE 'docs/evidence/[A-Za-z0-9._/-]+' … | sort -u` → `docs/evidence/s1/ci-step-readings-2026-09-22.md`（在盘）＋ `docs/evidence/s1/115-`（＝它自己承认名下没有）。
- 最后一次动这笔票 = `6d22dd6c 09-29 19:01 gate-rerun-1 交件：12 枚全部升成〔有读数…〕`。
- 关联票仍开放：`ls .scratch/wisp/issues/ | grep "^230"` → **`230-four-cells-left-unfinished-inside-closed-tickets.md`（无 `-done` 后缀＝还开着）**；票面 `:44`、`:46`、`:71` 三处把账指向票 230。
③ 归口：**缺凭据（要补）**，两件事分属两类腿：
(a) 115 名下的 1:1 终裁表——**必须非实现者写**，本腿不写、更不替编排者签；
(b) AC#4 的变异两向与 AC#6 的门禁四数是**缺读数**（不是缺归口），需要一枚能跑 go 的腿（本腿零 go，做不到）。
★另登记：`-done` 后缀与票面 `Status: BLOCKED 部分` 自相矛盾 ⇒ 见 §3-4。

### 2.11 票 243 — `.scratch/wisp/issues/243-a-closed-ticket-number-still-owns-work-in-25-comments-which-is-the-shape-that-fooled-three-legs-done.md`

① 等的凭据形状特殊：票面 `:10` 标题逐字 "## 要交的东西（一张表，不是代码）" ⇒ 本票要的**不是再一张表**，
而是**对那张已交的逐处判定表的非实现者验收／签收**。
② 现量：
- s1 名下 `0`；`find .scratch/wisp/probes/243 -type f` = **1 枚** = `.scratch/wisp/probes/243/c1/census.md`（**177 行 / 35,375 字节**；
  `:1` 逐字 "# 票 243 ／ 腿 `243-c1` — 「ticket 07」25 处注释逐处判定表（只读普查，零产码）"；`:4` 自陈纪律"零 `.go` 改动、零注释改动、零 `go test`、零 `go build`、零新仪器"）。
- 未勾 0；`Status:`（`:3`）逐字 `**已立，未派**（09-30 15:3x，编排者立；由「票 228 AC#7」拆出…）`。
- 最后提交 `08cec5a6 09-30 16:10 ledger(A473 收 243-c1：票 243 三格全成立、翻…`；
  票面 `:17` 标题逐字 "## 编排者收表 — 09-30 16:2x（普查腿 `243-c1` 交件；账 `A473`；三格全成立、翻满、改名结案）"，同节自陈"**我按三把尺自认，没采信通知正文**"。
- 反向引用尺：`grep -l -F "243-c1" docs/evidence/s1/*.md` → **`docs/evidence/s1/228-resident-ball-v1.md`**（别票名下的件引用它）。
③ 归口：**票本身不该算结案（凭据矛盾处＝判定表在盘、被别票名下的件引用，但"三格全成立、翻满、改名结案"这一判语出自编排者本人之手）**。
⇒ 缺的是**名下一张非实现者表**那一半，不是缺普查材料。归口动作＝需非实现者补签（本腿不签），不需要重做那 25 处判定。
★同一枚票面还有第二处矛盾：`:3` 的 `Status:` 停在"已立，未派"、与 `-done` 与"翻满、改名结案"三处互相打脸 ⇒ §3-1。

### 2.12 票 250 — `.scratch/wisp/issues/250-portable-tests-guard-c-merges-go-list-stderr-into-its-own-denominator-so-a-cold-module-cache-kills-the-whole-core-scope-reading-done.md`

① 终裁表（四格都是"尺／门禁"形状，不涉界面）。
② 现量：
- s1 名下 `0`；`find .scratch/wisp/probes/250 -type f | wc -l` = **29 枚**，其中 `.scratch/wisp/probes/250/r1/verdict.md` = **273 行 / 25,344 字节**；
- 那份件的 `:1` 逐字 "# 票 250 · **落地腿 `250-r1`** 裁决件 —— `portable-tests.sh` GUARD C 的分母被 `go list` 的 stderr 污染"，`:3` 起即 "## ① 起手锚与写面"
  ⇒ ★**文件名叫 verdict、署名是落地腿自己**（本腿只读前 4 行取署名，不读判语）。
- 反向引用尺：`grep -l -F "250-r1" docs/evidence/s1/*.md` = **零命中**；`grep -l -F "250-v1" docs/evidence/s1/*.md` = **零命中**。
- `find .scratch/wisp/probes/250 -type f | grep -iE 'verdict|accept|裁决|验收'` 只回上面那一枚
  ⇒ **250 名下连"落错地方"的非实现者件都没有**（与票 263/268 形状不同，那两枚至少各有一枚署名非实现者的表在 `.scratch` 里）。
- 未勾 0；`Status:`（`:3`）逐字 `**已立，待派**（10-02 08:5x，编排者立；来源＝只读取证腿 ci-delta-1…）`；最后提交 `e39386b7 10-02 09:27 ledger(A516)：33-r8 死腿收尾＋250-r1／195-a1 …`。
③ 归口：**缺凭据（要补）**。尺证据＝三把（s1 名下 0／两枚腿号 `grep -l` 零命中／probes 内 verdict 形状只有落地腿自署那一枚）。
不许读成"没做过工作"：29 枚台件确在盘，缺的是**另一个体的读数**。另登记 §3-1（状态行"已立，待派"与 `-done` 矛盾）。
