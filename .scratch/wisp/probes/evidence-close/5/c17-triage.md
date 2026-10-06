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

### 2.13 票 251 — `.scratch/wisp/issues/251-winsec-pin-is-a-pin-nobody-compares-the-winsec-scope-runs-through-explicit-paths-so-guard-c-never-sees-it-done.md`

① 终裁表（票面 `:36` 自己写明"判语归非实现者"）。
② 现量：
- s1 名下 `0`；`find .scratch/wisp/probes/251 -type f | wc -l` = **19 枚**，verdict 形状只有 **`.scratch/wisp/probes/251/r1/verdict.md`（171 行 / 22,436 字节）**；
- 它 `:1` 逐字 "# 票 251 r1 交件读数（**写码腿 `251-r1`**）"，`:3` 逐字 "> 本节只交**读数**，判语归非实现者。…四枚 AC 的勾选框本腿一枚未碰（本仓硬规矩：`- [ ]`→`- [x]` 归编排者）"
  ⇒ ★这件**自己明说它不是裁决**（与 250 那枚"名叫裁决件"形状相反），所以它是**待裁的料**，不是凭据。
- 反向引用尺：`grep -l -F "251-r1" docs/evidence/s1/*.md` = **零命中**；`grep -l -F "251-v1" docs/evidence/s1/*.md` = **零命中**；`ls -d .scratch/wisp/probes/251/v*` = **无**。
- 票面 `:36` 标题逐字 "## 交件要求（腿只交读数，判语归非实现者）"；`:92` 标题逐字 "## 10-02 10:2x **编排者落地回执（四格我勾**；未结案的残余另立票 254）"；`grep -c "非实现者"` 本票面 = **1**（就是 `:36` 那行"要求"，不是"已有"）。
- 未勾 0；最后提交 `a924fa17 10-03 16:38 票 137/251 结案改名 -done：137 五格全勾＋A1…`。
③ 归口：**票本身不该算结案（凭据矛盾处＝四格由编排者"我勾"，而名下唯一件自陈"判语归非实现者"、那份判语至今 0 枚）**。
按"裁决者≠实现者"这条硬规：编排者自勾＋写码腿自读数 ≠ 凭据。要补的是**一枚 251 名下的非实现者表**（本腿不写）；
251 的落地料齐不齐不是本腿的判断。另登记 §3-2（"另立票 254"与两枚同族票都 `-done` 的关系）。

### 2.14 票 254 — `.scratch/wisp/issues/254-fifth-tier-winsec-no-ci-caller-done.md`

① 终裁表。
② 现量：
- s1 名下 `0`；`find .scratch/wisp/probes/254 -type f | wc -l` = **33 枚**，verdict 形状只有 **`.scratch/wisp/probes/254/r1/verdict.md`（280 行 / 24,115 字节）**；
- 它 `:1` 逐字 "# 票 254 — 交件（腿 `254-r1`）"，`:3` 逐字 "> 本腿只做读数与落地，**不做验收判语**；票面三格 AC 的勾选框一枚未碰" ⇒ 同 2.13 形状＝待裁的料。
- ★票面自己把定性写在标题里：`:36` 逐字 "## 编排者裁定与翻勾（2026-10-02 11:19:06+0800 现量，台账 `A529`）"；
  `:38` 逐字 "⚠ **这一节的定性要说实话：三格是〔编排者现跑自勾〕，不是独立验收腿裁的。**"
- `grep -c "非实现者"` 本票面 = **0**（票面从未要求过非实现者腿、也从未记它有）；`find .scratch/wisp/probes/254 -type d` = 只有 `r1`／`r1b`，**无 v 腿目录**。
- 反向引用尺：`grep -l -F "254-r1" docs/evidence/s1/*.md` → **`docs/evidence/s1/248-settings-write-path-v1.md`**（别票名下一枚名字带 `-v1` 的件提到本票腿号；本腿未读其内容、不据此改判，见 §3-5）；`grep -l -F "254-v1" …` = **零命中**。
- 未勾 0；最后提交 `46079fcc 10-04 09:18 ledger(A584) 收 262-a1：它推翻我 A583 两处数…`。
③ 归口：**缺凭据（要补）**＝票面已具名自陈"编排者现跑自勾、不是独立验收腿裁的"，且名下与别票名下都没有 254 的裁决表。
这一枚是 17 枚里**最不需要本腿推断**的一枚：它自己的票面就把水分写明了。
⚠ 补它的人要注意票面 `:38` 具名的原因（同一道题上三枚腿连续死于环境），⛔ 别只补一张表就完事。

### 2.15 票 263 — `.scratch/wisp/issues/263-slo-check-dies-on-unset-lastexitcode-done.md`

① 终裁表（＋ AC#5 那一格等的是"一次真 run 读数"）。
② 现量：
- s1 名下 `0`；`find .scratch/wisp/probes/263 -type f | wc -l` = **198 枚**（子目录只有 `r1`／`v1`）；
- ★**名下有一枚署名非实现者的验收表**：`.scratch/wisp/probes/263/v1/verdict.md` = **304 行 / 33,526 字节**；
  `:1` 逐字 "# 票 263 验收表（**263-v1，非实现者**）— 攻 `263-r1` 那五枚 commit"；
  `:4` 逐字 "本腿＝`263-v1`（非实现者验收腿）。⛔ 未改任何产码、⛔ 未碰票 263 的任何 AC 框、⛔ 未 push"。
- 票面引用它：`:57` 逐字 "**2026-10-04 13:1x（编排者翻勾，凭据＝非实现者验收腿 `263-v1` 的 `.scratch/wisp/probes/263/v1/verdict.md`，304 行／33,526 字节／占位 0，件在盘 `515ca5c5`）**…"；
  `:59`、`:60` 两处继续引它（"验收腿自己量的""验收腿量到一枚名册外的第三牙"）⇒ ★**行数／字节数本腿独立 `wc` 复量与票面逐字相同（304／33,526）**，票面那句引用不是空引。
- 最后一格的凭据：`:62` 逐字 "**2026-10-05 09:0x（编排者翻勾 AC#5 ⇒ 本票六格全勾，改 `-done`）**：凭据＝**push 事件 run `37249563077`**（headSha 逐字 `21bec8a1…`…"；
  ★该枚 run 日志**不在** `.scratch/ci-logs/`（那里只有 `run-37166458550-failed.{log,err}`）⇒ 那一格今天不可复跑。
- 反向引用尺：`grep -l -F "263-v1" docs/evidence/s1/*.md` = **零命中**；`grep -l -F "263-r1" docs/evidence/s1/*.md` = **零命中**。
- 未勾 0；最后提交 `19b06ac3 10-05 09:10 票 263 结案：AC#5 翻勾（run 37249563077 slo-ful…`。
③ 归口：**凭据在别处（指路径）＝错档**，不是缺凭据。真表在 `.scratch/wisp/probes/263/v1/verdict.md`
（非实现者署名＋票面逐字引它当翻勾凭据＋行数尺复量相符）。要做的只是**把那份件落到 `docs/evidence/s1/263-*.md`**
（＝落点搬运，符合 `AGENTS.md` §1.5 要求的形状），不是重新验收。
⚠ 一格例外：AC#5 的凭据＝一枚远程 run 且日志不在盘上 ⇒ 那一半属"缺一次可复跑的读数"，见 §4。
★**本腿与前一枚索引腿读数不同、具名登记**：该件 §3.3-3 把 263 与 250/251/254 归为同一形状（"实现侧普查台件，本就不是终裁表"）——
本腿现量证其对 250/251/254 成立、**对 263 不成立**（263 名下确有 304 行／33,526 字节的非实现者验收表）。⇒ §3-3。

### 2.16 票 267 — `.scratch/wisp/issues/267-confirm-timeout-config-has-no-bounds-so-the-c18-warning-can-vanish-done.md`

① 在等：**终裁表**，但**只有 AC#4 那一格真的缺**（其余四格有非实现者料）。
② 现量：
- s1 名下 `0`；`find .scratch/wisp/probes/267 -type f | wc -l` = **39 枚**，子目录＝`a1 a2 a3 gate r1 r2` ⇒ ★**没有 `v` 腿目录**。
- 名下料逐枚 `wc -l -c` 与署名行：
  - `.scratch/wisp/probes/267/a2/census.md` **362 行 / 77,770 字节**，`:1` 逐字 "# 票 267 · 267-a2 **只读普查**：`confirm_timeout_sec` 值域带内/带外名册 + 逐枚调用点观察对象"，`:3` 自陈"零 `go` 命令、零已跟踪文件改动、未读 `frontend/**` 与 `design/**`"＝**只读普查腿**；
  - `.scratch/wisp/probes/267/r1/evidence.md` **294 行 / 30,943 字节**，`:1` 逐字 "# 票 267 **落地腿 `267-r1`** — 证据件"；
  - `.scratch/wisp/probes/267/r2/evidence.md` **254 行 / 44,771 字节**＝同为落地腿。
- 票面 AC#4 那一格**当初自己说不翻**（`:59` 末段逐字）："**AC#4 不翻**：…而且是**红的**…⇒ **本票这把 band 把下游打红了**。逐名归因在 `A611`，名册台件＝`.scratch/wisp/probes/267/gate/cmdwisp-HEAD.log`…⛔ 这不是"AC#4 判失败要退回整票"，而是**AC#4 的凭据要等种子迁移之后那一发复跑才齐**"。
- 后来那一格被翻掉了：`:72` 起标题逐字 "## 结案（编排者，2026-10-05 10:5x，锚 `c7bb02be`；**凭据全部编排者现跑**，账 `A615`）"；
  `:74` 逐字 "**五格全勾 ⇒ 本票改 `-done`。** AC#4 的凭据（⛔ 不是腿的自述，是我这两小时自己跑的）"；
  `:80` 同一节自陈一处尺失效"我第一次跑 gofumpt 用了 `~/GOPATH/bin/gofumpt`…那次"空输出"**不是干净，是尺没跑到**"。
- 反向引用尺：`grep -l -F "267-v1" docs/evidence/s1/*.md` → `265-267-evidence-index.md`、`closed-tickets-evidence-index.md`；`grep -l -F "267-r2" …` → `265-267-evidence-index.md`。
- 未勾 0；最后提交 `c6cf66e6 10-05 10:58 ledger(A614＋A615) 票 267 五格全勾改名 -done…`。
③ 归口：**缺凭据（要补）**，射程本腿能钉准＝**只缺 AC#4 那一格的一次非实现者复跑**（AC#0/#1/#2/#3 的料在盘：a2 普查 362 行＋两枚落地腿证据件，
它们都**不是**裁决表形状）。理由具名：那一格在 `-done` 之前按票面口径"不翻、等种子迁移后复跑"，之后由**编排者自己跑两小时**把它翻掉 ⇒ 与"裁决者≠实现者"直接撞。
不许读成"17 枚里 267 最严重"：本腿只说"缺 AC#4 那一格"，其余四格该由谁签属另一笔账（票面自陈账在 `A615`）。

### 2.17 票 268 — `.scratch/wisp/issues/268-resident-leg-folds-missing-config-and-refused-config-into-one-provenance-done.md`

① 终裁表。
② 现量：
- s1 名下 `0`；`find .scratch/wisp/probes/268 -type f | wc -l` = **56 枚**，子目录＝`a1 a2 r1 v1` ⇒ ★**有 `v1` 腿目录**（与 267 不同）。
- 名下两件（都不落 `s1/`、都不是 `<票号>-` 名字）：
  - `.scratch/wisp/probes/268/v1/evidence.md` = **459 行 / 62,585 字节**＝验收腿料；★它的提交史＝`git log --format='%h %ad %s' -- <该件>` 现量**十一笔**，署名依次是
    `268-v1 笔1/笔2`（`1c12712a` 20:08、`2875930b` 20:11）、`268-v1b 笔1—3`（`01b38812` 20:32／`554cb945` 20:35／`438f5050` 20:41）、`268-v2 笔1—6`（`689bc241` 21:57 → `ba2482c9` 22:02 → `862556b3` 22:10 → `0d207e14` 22:15 → `f73558c0` 22:28 → `e8df2643` 22:43）
    ⇒ **一枚文件、三枚腿接力署名**；
  - `.scratch/wisp/probes/268/a2/verdict.md` = **205 行 / 26,456 字节**，`:1` 逐字 "# 268-a2 独立复核：对 `268-a1/census.md` 复跑 AC#0 四问那把尺——逐问判"复现／不复现／口径不同""，`:3` 逐字 "编队：`268-a2`（**非实现者腿，裁决腿**）…"
- 票面的凭据陈述：`:52` 标题逐字 "## 5. 编排者翻勾记录（2026-10-05 22:5x +08，锚 HEAD `e8df2643`）"；`:54` 逐字
  "**凭据＝非实现者验收腿 `268-v2`**（`.scratch/wisp/probes/268/v1/evidence.md` **459 行／62,585 字节**，六笔 `689bc241`→…→`e8df2643`；实现腿＝`268-r1`…）"
  ⇒ ★行数／字节本腿独立 `wc` 复量**完全相符**（459／62,585），六笔提交号也与 `git log` 逐枚对上；
  但**代号写 `268-v2`、路径写 `v1`**。票面 `:73` 又自陈"终裁＝`268-v1`／`268-v1b`／`268-v2` 三腿接力（前两枚死于服务故障、第三枚撞 150 轮帽但六笔全落）…⚠ 同目录代号撞车已在 A626 记我"。
- 反向引用尺：`grep -l -F "268-v2" docs/evidence/s1/*.md` = **零命中**；`grep -l -F "268-r1" …` → `265-267-evidence-index.md`；`grep -l -F "268-a2" …` → `closed-tickets-evidence-index.md`。
- 未勾 0；`Status:`（`:3`）逐字 `**done（2026-10-05 22:5x 编排者翻勾＋改名，凭据见下面第 5 节…**；最后提交 `11105265 10-05 22:56 ## 5. 编排者翻勾记录（2026-10-05 22:5x +08，…`。
③ 归口：**凭据在别处（指路径）＝错档**。非实现者料在盘且可指：AC#0 那一格＝`.scratch/wisp/probes/268/a2/verdict.md`（205 行，署名"裁决腿"），
其余四格＝`.scratch/wisp/probes/268/v1/evidence.md`（459 行，v1／v1b／v2 三腿接力）。缺的是**把它们落到 `docs/evidence/s1/268-*.md`**（`AGENTS.md` §1.5 的落点规矩）。
⚠ 一条必须转出去的形式账（不改变"料在不在"的判断）：票面 §5 把凭据代号称作 `268-v2` 而路径是 `v1`，且 `:73` 自陈同目录代号撞车
⇒ "这 459 行里哪一节出自哪一枚腿"要按提交号逐笔核，那是 owner 或非实现者表的活，见 §3-6 与 §4。
★与前一枚索引腿的分歧：该件 §3.3-3 形状③只登记 `268/a2/verdict.md` 一件、并说"票 268 AC#0 的翻勾注记逐字引它"——
本腿现量另见 `268/v1/evidence.md`（459 行／62,585 字节）那一件才是**其余四格**的凭据，AC#0 只是 a2 那一格 ⇒ §3-3 一并登记为"指错路径（不全）"。
