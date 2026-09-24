# 引用完整性清点 2026-09-24（只读审计）

审计对象：`docs/reports/pending-and-issues.md`（5772 行）、`docs/reports/HANDOVER.md`（950 行）。
审计方式：全部 `grep` / `ls` / `git cat-file` 级取数，未跑任何构建或测量命令。
本报告不提出修法，只登记"哪些是缺陷＋证据＋现量命令"。

## 结论（全文 ≤400 字）

两枚文件里共抽出 **1237 处**十六进制样式 token、**634 枚去重**。**512 枚**可解析为本仓 commit；
**122 枚**不可解析，其中 **63 枚是全数字**（GitHub Actions run/job id、字节读数、日期），
**59 枚含字母**。59 枚里：**36 枚**上下文已明说"不存在／假号"（正常，幽灵投递登记），
**14 枚**不是 commit 号而是 md5／sha1／git blob／agent 会话文件名（正常，但报告里逐枚点名理由），
**9 枚**是**缺陷**——上下文把它当成"存在"或"我已复量"。

9 枚缺陷中 **2 枚**（`0377e87`、`f1086b4`）台账后续条目已自我作废；
**7 枚**至今无任何更正记录：`5c6f824`（起手锚，两处）、`1e8f8eb`（与两枚真号并列）、
`1f8d212`（A79② 结案账）、`e0071e2`／`923f7c4`／`4a91810`（A192⑤ 转抄的三枚"带时间戳的 commit"）。

**简报前提有两处不成立**，已照实回：① "台账四次写'某枚 commit 我复量存在'"——在这两枚文件内
"复量存在"这个句式只出现 **3 次**，其中 **2 次**指向不存在的号（`0377e87`、`f1086b4`），
第 3 次（`ea59b8c`）是真号；另外两枚若存在则落在本范围之外。
② "台账里写 `gate_:51 gateProbe` 未使用"——台账里**没有**那句未经更正的断言：`gateProbe`
在两枚文件各只命中 1 行，且都在**自我作废**的句子里（`:5768`、`HANDOVER:108`），行号也不是 `:51`
而是 `:62`／`gateFailed`。附带一条**当前为假**的自证：`:5768` 写"两枚文档里也零命中"，
现量两枚文档各 **1** 命中（命中就是这句话自己）。

---

## 检查 1：sha 可解析性

```
$ date "+%Y-%m-%d %H:%M %z"
2026-09-24 22:54 +0800

$ git log -1 --format='%h %ad' --date=format:'%H:%M'
faf66d3 22:50
```

### 1.0 抽数与分母

⚠ 审计期间被审文件仍在长（append-only）：第 1 节初抽在 `faf66d3 22:50`，复抽在 `42d66f7 23:06`，
新增 10 行全落在 `:5764` 之后 ⇒ **本节所有较早行号未漂移**，下表用**复抽后的现量**。

| 项 | 数 | 现量命令 |
|---|---|---|
| 十六进制样式 token 出现处 | 1209 ⇒ **1237** | `grep -n -o -E '\b[0-9a-f]{40}\b\|\b[0-9a-f]{7,12}\b' <两枚文件>` |
| 去重 | 624 ⇒ **634** | 同上 `\| sort -u` |
| ⓐ 可解析为本仓 commit（含 7–12 位前缀） | 505 ⇒ **512**（其中 20 枚是全数字真号，如 `9048866`） | 与 `git cat-file --batch-all-objects` 的 commit 全集前缀表求交 |
| 不可解析 | 119 ⇒ **122** | 同上求差 |
| ┗ 全数字（run/job id／字节读数／日期） | **63** | `awk '$2 ~ /^[0-9]+$/'` |
| ┗ 含字母（sha 形状） | 56 ⇒ **59** | `awk '$2 ~ /[a-f]/'` |

快照留在 `/d/tmp/citation-integrity-2026-09-24/`（`tokens-raw.txt`／`cls.txt`／`cand-shas.txt`／`cand-notsha.txt`），按规矩只建不删。

### 1.1 ⓑ 缺陷：不可解析、但上下文当成"存在／我已复量"（9 枚）

| sha | `file:line` | 上下文原句（截） | 现量命令 | 是否已被后续更正 |
|---|---|---|---|---|
| `0377e87` | `docs/reports/pending-and-issues.md:5728` | 末枚 `0377e87` 我复量存在 | `git cat-file -t 0377e87` ⇒ fatal | 已作废于 `:5768`（A200④） |
| `f1086b4` | `docs/reports/pending-and-issues.md:5734` | `f1086b4` 我 `cat-file` 复量存在 | `git cat-file -t f1086b4` ⇒ fatal | 已作废于 `:5744`（A197②）、`HANDOVER.md:104` |
| `690e87a` | `docs/reports/pending-and-issues.md:5734` | 2 枚 commit `f1086b4`＋`690e87a`（同句被当作已核事实转抄） | `git cat-file -t 690e87a` ⇒ fatal | 已作废于 `:5744` |
| `5c6f824` | `docs/reports/pending-and-issues.md:5660` | 表 …5 枚单路径 commit `bd0f826`→…→`c4d54c6`，**起手锚 `5c6f824`** | `git rev-parse --verify 5c6f824` ⇒ NO；`git log --all --format='%H' \| grep -c '^5c6f824'` ⇒ 0 | **无更正** |
| `5c6f824` | `docs/reports/HANDOVER.md:91` | 401 行／5 枚 commit …，**锚 `5c6f824`** | 同上 | **无更正** |
| `1e8f8eb` | `docs/reports/pending-and-issues.md:5667` | `R-140-1` 那三处恒真断言（`c2fa2e9`／`2f22ab3`／`1e8f8eb`）改成非恒真 | 并列的两枚 `git cat-file -t` ＝ commit，`1e8f8eb` ⇒ fatal | **无更正** |
| `1f8d212` | `docs/reports/pending-and-issues.md:1691` | AC#1/AC#2 = 原代理（`1f8d212`）、代码＝编排者代落档（`d0d8782`） | 同句其余三枚（`d0d8782`／`06906f7`／`9e00629`）全为 commit，仅它 ⇒ fatal | **无更正** |
| `e0071e2` | `docs/reports/pending-and-issues.md:5702` | 文本自 `e0071e2` 09-21 02:27 未变 | `git cat-file -t e0071e2` ⇒ fatal | **无更正**（A196① 只作废了该句的 gofmt 结论，未点这枚号） |
| `923f7c4` | `docs/reports/pending-and-issues.md:5702` | `ci.yml` 的 gofmt 步骤 `923f7c4` 09-20 18:22 起未动 | 同上 ⇒ fatal | **无更正** |
| `4a91810` | `docs/reports/pending-and-issues.md:5702`、`:5704` | 该文件自 `4a91810` 起未跟踪（两处同一断言） | 同上 ⇒ fatal | **无更正** |

⚠ 这一组三条（`e0071e2`／`923f7c4`／`4a91810`）与已作废的 `0377e87` 出自**同一份转抄来源**
（`A192⑤`，程 `auditor-ci-read-a1fd5bf-r1`），台账只作废了结论层，号层没跟着查。

### 1.2 ⓒ 正常：不可解析但上下文已在说它不存在（36 枚）

| sha | 最靠前的 `file:line` | 上下文的定性 |
|---|---|---|
| `17b61a4` `3c3e388` `d0ca6a3` `d667888` | `:5059`（`d667888` 另见 `:5090`） | "报的 sha 全部不存在 ⇒ 逐枚 `git cat-file -t` 失败" |
| `99319c7` | `:4985`（另 `:5115`） | "`git cat-file -t 99319c7` ＝ 不存在" |
| `1049513` `1b93810` | `:5115`、`HANDOVER.md:235` | "两枚号盘上不存在＝它的 `git commit` 那一步静默失败" |
| `d62a943` `a3f6f55` `8b1d2e5` `99f937a` `56f2c7a` `328d563` | `:5230` | A163 假交件："那六枚号＝六枚全部 `fatal: Not a valid object name`" |
| `f2302b2` `a2eda97` `623d795` | `:5251` | "新报的三枚号逐枚 `git cat-file -t`＝同样全部不存在" |
| `6240d2f` `d0a8902` | `:5263` | "不存在的包括…同批里的几枚" |
| `18f2531` | `:5353`、`:5354`、`HANDOVER.md:181` | "给一枚不存在的 commit `18f2531` 编了 git 输出" |
| `cc7262c` | `:5354` | "短哈希手抄错 `cc7262c`→真值 `cc726ca`"（真值现量＝commit） |
| `87168e6` | `:5329` | "`87168e6` 在盘上不存在" |
| `f4b3ea5` | `:5431`、`HANDOVER.md:179` | "凭记忆写成盘上不存在的 `f4b3ea5`" |
| `c135d50` | `:5684` | "在本仓不是有效对象" |
| `75c0159` `92f0b30` `1a5588a` `7c2b98e` | `:5709` | "它列出的 5 枚不存在 sha…包含它自己收到的回显 `7c2b98e`（真号 `10a0bbe`）"（真号现量＝commit） |
| `840c9a9` | `:5725` | "同批我引过的那枚 `840c9a9` 不存在、真号 `1bb92cc`"（真号现量＝commit） |
| `927f484` `0869b87` `60a9d2a` `428d0e8` | `:5726` | "三枚 commit 全部 `git cat-file` 失败，它引的'我上一条报告'的 `428d0e8` 也不存在" |
| `278d3538`（含 40 位全形 `278d3538f7721990557c7b516c3678ed5138823d`） | `:4421`、`:4422` | "第 8 代注入候选：一串假 git sha 混进工具输出" |
| `63a9167` `99e2d14` `63a6010` | `:5770` | "一枚假 `git log -6` 给了 …，三枚全不解析"（A200⑥ 转述，我这轮逐枚 `cat-file` 复量＝三枚皆 fatal） |

### 1.3 ⓓ 正常但不是 commit 号（14 枚，逐枚给判据）

| token | `file:line` | 它是什么 | 判据 |
|---|---|---|---|
| `00a65062` `42a03694` | `:5334` | `md5sum` 前缀 | 同句写死 `git show <sha>:<file> \| md5sum` 的输出 |
| `14afa08` `336a50f` | `:5099` | `md5sum` 前缀 | 同句"`md5sum` 逐棵" |
| `831a5c08` | `:5479`、`:5504`、`HANDOVER.md:166` | 函数体 md5 前缀 | 同句"awk 抽函数体 → md5 前缀"；四版同值＝文件内容不是号 |
| `87496d505e` | `:2542` | md5 前缀（10 位） | 同句"md5 回到 `87496d505e…`" |
| `906201f4` `906201f4a3995d10b0a65910aa4cc68e1f745a9d` | `:4550`、`:4956`、`HANDOVER.md:155` | `sha1sum` 文件摘要 | 同句"`sha1sum`＝…"，40 位全形是文件哈希 |
| `23b443ac` | `:4891`、`:4893`、`HANDOVER.md:323` | `sha1sum` 摘要 | 台账自己已判定："引哈希要说清算法形状" |
| `bde61ddd` | `:4891`、`:4896`、`HANDOVER.md:323` | **git blob** | `git cat-file -t bde61ddd` ⇒ `blob`；同句"`git hash-object` 均 `bde61ddd…`" |
| `e2677b11` | `:3924` | **git blob** | `git cat-file -t e2677b11` ⇒ `blob`；同句"两侧同 blob" |
| `1060e48` | `:5651` | 自建 exe 的 sha1 前缀 | 同句"自建 exe `sha1 1060e48…`" |
| `1b62a1d3` | `:4351` | agent 会话 jsonl 名 | 同句"`1b62a1d3` jsonl mtime 22:30:16 仍在写" |
| `33e37763` | `:5282` | agent 会话 jsonl 名 | 同句"`agent-…33e37763….jsonl` mtime" |

### 1.4 全数字 63 枚

`355xxxxx`／`356xxxxx`／`358xxxxx`／`359xxxxx`／`36003984868`＝GitHub Actions **run id**（同句必带 `run` 或 `gh run view`），
`106xxxxxx`／`107xxxxxx`＝**job id**（同句 `--log-failed --job`），
`10485760`／`2158592`／`4411392`／`4476928`／`27881312`＝**字节读数与 D32 阈值**，
`20260923`＝日期。⇒ 全数正常，无冒充 commit。唯一两枚全数字被当 commit 号点的 `1049513`／`2644207` 已在 §1.2（上下文自身判定不存在）。

### 1.5 反向抽查（该号真、却被说不存在）

把出现"不存在／cat-file／rev-parse"字样的行里的 sha 全抽出与可解析集求交，得 57 枚——
逐枚看上下文后**无一例是"真号被误判为假"**：命中都是"同一长行里另一枚真号"（如 `:5725` 的 `ea59b8c` 真、`840c9a9` 假）。
`81b4d5f`（`:5263` 那句"不存在的包括…几枚"里）**现量＝commit**，属该句"几枚"措辞含糊，不算缺陷。

### 1.6 `gateProbe` 前提核查

| 简报前提 | 实测 | 现量命令 |
|---|---|---|
| 台账里写"`gate_:51 gateProbe` 未使用" | **不成立**（未经更正的断言不在这两枚文件里） | `grep -n "gateProbe"` ⇒ 仅 `:5768` 与 `HANDOVER.md:108`，两处都在自我作废句中；台账里真实写过的未使用符号是 `gateFailed`（`:5719`、`:5728`、`:5732` 行号 `:62`） |
| 盘上 `sampler_settle_gate_136_test.go` 无 `gateProbe` | **成立** | `grep -c gateProbe internal/observe/sampler_settle_gate_136_test.go` ⇒ 0；`git show 021a549^:… \| grep -c gateProbe` ⇒ 0 |
| "`:5768` 那句自证：两枚文档里也零命中" | **当前为假**（各 1 命中，命中即该句本身） | `grep -c gateProbe docs/reports/pending-and-issues.md docs/reports/HANDOVER.md` ⇒ 1 / 1 |

### 1.7 一处"更正句自己带错计数"（新发现，属 ⓑ 一族的镜像）

`:5771` 的 A200④ 出处更正句写："**其五枚 commit `2d05932`/`bc67475`/`61db519`/`41b1869`/`9048866` 我逐枚 `cat-file` 复**…"。
现量：五枚**全部**是 commit（正常），但那枚表 `docs/evidence/s1/ci-read-a1fd5bf-r1.md` 的**真 commit 是 6 枚**，
名单漏了 `7cc5050`（21:49，§2）；而"末枚"＝`9048866`（22:13）≠ 台账 `:5728` 写的 `0377e87`。

```
$ git log --format='%h %ad' --date=format:'%H:%M' -- docs/evidence/s1/ci-read-a1fd5bf-r1.md
9048866 22:13 / 41b1869 22:09 / 61db519 22:01 / bc67475 21:56 / 7cc5050 21:49 / 2d05932 21:42
$ wc -l < docs/evidence/s1/ci-read-a1fd5bf-r1.md
699
```

---

## 检查 2：`docs/evidence/**` 路径存在性

```
$ date "+%Y-%m-%d %H:%M %z"
2026-09-24 23:25 +0800

$ git log -1 --format='%h %ad' --date=format:'%H:%M'
4013b1d 23:21
```

（本审计第 3 节的全部取数就在这枚 HEAD 上做；两枚被审文件自 `faf66d3` 起只增长过尾部，早段行号未漂。）

### 2.0 分母

| 项 | 数 | 现量命令 |
|---|---|---|
| `docs/evidence/…\*.md` 形式引用（出现处） | 103 | `grep -n -o -E 'docs/evidence/[A-Za-z0-9_./*-]+\.md' <两枚文件>` |
| 去重路径 | 75 | 同上 `\| sort -u` |
| 具体（非通配）路径 | 69 | — |
| ┗ 盘上存在 | **67** | `[ -f <p> ]` 逐条 |
| ┗ 不存在＝ⓒ（上下文已说是幽灵） | **2** | 见 §2.1 |
| 通配/省略写法＝ⓓ | 6 | 见 §2.2 |

另有 **98 枚**裸文件名（`NN-slug.md`，不带 `docs/evidence/s1/` 前缀）引用：
**72 枚**在 `docs/evidence/s1/` 命中、**2 枚**在 `docs/evidence/{s0,s2,s3}/`、**24 枚**不在 evidence 目录但已逐枚定家（见 §2.3）。

### 2.1 ⓒ 不存在的两枚（上下文自己已经判定不存在＝正常）

| 路径 | 引用处 | 上下文的定性 |
|---|---|---|
| `docs/evidence/s1/140-m1m2m4-r1.md` | `pending-and-issues.md:5726` | "`docs/evidence/s1/140-m1m2m4-r1.md` 不存在"（幽灵投递） |
| `docs/evidence/s1/141-q46c-blocked-r1.md` | `:5734`（作为存在引用）、`:5744`、`:5760`、`HANDOVER.md:104` | `:5734` 那一行当时把它当**真表**引用（"以下数字一律以它的表为准…449 行"），`:5744` 起作废 ⇒ 该行的**引用形状**已被后文钉死，不再单独计缺陷 |

`docs/evidence/s1/141-q46c-impl.md`（真身）**在盘上**，`wc -l`＝**245**，与 `HANDOVER.md:108` 的"`bb61dc5` 落了…245 行"逐字相符；
`git show --stat bb61dc5` ⇒ `1 file changed, 245 insertions(+)`。⇒ 该条**不是缺陷**。

### 2.2 ⓓ 通配与省略写法（6 枚，不是坏引用）

`docs/evidence/s1/*.md`（`:3908`）、`102-*.md`（`:1865`）、`119-ac7-*.md`（`:5474`）、`125-*.md`（`:5135`）、
`130-ac3-*.md`（`HANDOVER.md:448`）、`134-*.md`（`HANDOVER.md:425`）——全是"这一族文件"的 glob 写法，
上下文要么在描述判据范围、要么明说"这样的文件一枚都没有"。
另有 `140-job-level-wisp-env-test-...md`（票 140 文件名的手写省略形，真身在盘）。

### 2.3 裸文件名的 24 枚去处（逐族给结论）

| 族 | 枚数 | 实落点 | 现量 |
|---|---|---|---|
| `SPEC-NN-<slug>.md` 只写了 slug | 6 | `docs/specs/SPEC-00…12` 全部在盘 | `ls docs/specs/ \| grep -F <slug>` |
| `2026-09-2X…` 只写了日期尾巴 | 5 | `docs/reports/2026-09-19-t14-dev-minisign-key.md`、`2026-09-23-gap-analysis-vs-oss-harnesses.md`、`2026-09-24-gap-analysis-audit-verdict.md`、`docs/evidence/s1/ci-runner-readings-2026-09-21.md`、`ci-step-readings-2026-09-22.md` 全在盘 | `ls docs/evidence/s1/ \| grep -E 'ci-(runner\|step)-readings'` |
| 工单文件名（`.scratch/wisp/issues/`） | 10 | 见检查 3 | — |
| `62-adversarial-acceptance.md` | 1 | **盘上从来没有过**：`:743` 与 `HANDOVER.md:856` 都明写"根本不存在" | ⓒ 正常 |
| 本报告文件名尾巴 `09-24.md` | 1 | `citation-integrity-2026-09-24.md`（`:5772` 派单原文） | ⓓ |

### 2.4 顺手复量的三处"我复量存在"的行数／枚数断言（检查 2 的延伸，非新增缺陷）

| 台账断言 | 现量 | 判定 |
|---|---|---|
| `:5725` "表 `136-ac15-denominator-census-r1.md`，678 行／5 枚 commit，末枚 `ea59b8c`" | `wc -l`＝**678**；`git log --format=%h -- <该表>`＝`ea59b8c dac2d89 c071d2d 38c71b4 8d7660e`＝**5 枚**、末枚正是 `ea59b8c` | **全部对上** |
| `:5660` "表 `140-ac1-verdict-r1.md` 401 行／5 枚单路径 commit" | 现为 **508 行／7 枚**（`e6fa582`、`fffe3a6` 后落）；台账 `:5670`（A189）已写"508 行／7 枚（末枚 `fffe3a6`）" | **已被后文追平**，不算缺陷 |
| `:5771` "其五枚 commit `2d05932`/`bc67475`/`61db519`/`41b1869`/`9048866`" | 该表真 commit **6 枚**，名单漏 `7cc5050`（21:49 §2） | **待更正**（详见 §1.7） |

---

## 检查 3：工单号可解析性

```
$ date "+%Y-%m-%d %H:%M %z"
2026-09-24 23:27 +0800

$ git log -1 --format='%h %ad' --date=format:'%H:%M'
8b28879 23:26
```

### 3.0 分母

| 项 | 数 | 现量命令 |
|---|---|---|
| `票 NN` 形式的号（去重） | 104 枚，覆盖 `02…141` | `grep -o -E '票 ?[0-9]{2,3}' <两枚文件> \| sort -un` |
| `.scratch/wisp/issues/…` 路径引用（出现处） | 25 | `grep -n -o -E '\.scratch/wisp/issues/[A-Za-z0-9._*<>-]+'` |
| 盘上工单文件（含 `README.md`） | 141 枚编号 + README | `ls .scratch/wisp/issues/ \| grep -E '^[0-9]'` |

### 3.1 指向不存在的票：**0 枚**

`票 NN` 抽出的 104 枚号逐枚对盘（`ls .scratch/wisp/issues/ | grep "^NN-"`）**全部有票**；
`02…141` 之间台账没用到的号不算缺陷。**没有一枚"票 NN"指向空号。**

### 3.2 指向"同号但已改名"的票（`-done` 后缀漂移）＝4 枚，待更正、非缺陷

| 引用处 | 台账写法 | 盘上真名 | 改名时刻（现量） |
|---|---|---|---|
| `:2689` | `80-blacklist-overrides-never-wired-to-gate.md` | `…-gate-done.md` | 原名 `87de494` 09-21 11:54 建；`-done` 形 `830621b` 09-21 12:12 起 |
| `:2648` | `81-containment-fixtures-are-windows-shaped.md` | `…-shaped-done.md` | 原名 `87de494` 同期建；`-done` 形 `c923327` 09-21 13:02 起 |
| `:2604` | `83-config-keys-that-lie-must-fail-loudly.md` | `…-loudly-done.md` | `-done` 形 `2e20920` 09-21 12:50 起 |
| `:5137` | `119-posix-link-leg-refuses-legitimate-symlinked-data-roots.md` | `…-roots-done.md` | 原名 `76fc5b0` 09-21 21:19 建；`-done` 形 `0bfd322` **09-24 18:07** 起（＝A176 翻 `-done` 那一步） |

四枚的被引用时刻都**早于**各自改名时刻 ⇒ 引用当时为真，属 append-only 台账的正常名字漂移；
`-done` 是"防重领的唯一键"（`AGENTS.md` §1.5），故按 `issues/README` 规矩这四处是**现名不符**，不是假号。
glob 写法 `119-*.md`（`:4195`、`:5365`）、`129-*.md`（`:3748`、`HANDOVER.md:448`）、`*-done.md`（5 处）**不算**这一类。

### 3.3 另一族工单号：`#NN`（编排者任务表）＝10 枚，**仓内不可解析**

`#16 #18 #40 #50 #63 #64 #65 #67 #68 #69`（出现次数分别 16/18/… ，见现量）。

```
$ grep -o -E '`#[0-9]{1,3}`' docs/reports/*.md | sed -E 's/.*#([0-9]+).*/\1/' | sort -n | uniq -c
  18 40   16 50    6 63    6 64   11 65    9 67    2 68    4 69 ...
$ find . -maxdepth 3 -name 'tasks*' -o -name '*.db' | grep -v '^./.git/'
（无输出：仓里没有任何任务表文件）
```

**这批号只存在于编排者的会话任务表，仓内无档案可查**，因此台账里凡是"按号找事"的句子都不具备可复核性；
台账自己已经撞上过一次并写下判据（`:5752` A198："本条 ⑤ 与 `A196 next=③` 写的 `#67` 在任务表里的实际编号是 `#69`……
以后按号找事**以任务表为准，不按本台账里的简称**"）。⇒ 计**待更正／结构性**，不计入"假号缺陷"。

### 3.4 工单**面行号**引用（属检查 4 的 `.scratch/**` 组，先在此点名）

| 引用 | 现量（HEAD `8b28879`） | 判定 |
|---|---|---|
| 票 140 面 `:16`（`:5727`"默认走 `userConfigDir()`＝`%APPDATA%`"） | `base, err := userConfigDir() // %APPDATA%` | **字面在** |
| 票 136 面 `:348`（`:5711`"在 HEAD 与工作树都还是 `[ ]`"） | `- [ ] **AC#15**（09-24 14:5x 编排者追加…）` | **字面在** |
| 票 141 面 `:79`（`HANDOVER.md:108`"三处同形逐处 grep 到"） | 该行含"面向用户的字符变严／注释面豁免"，与 `docs/evidence/s1/141-q46c-impl.md:10` 同句 | **字面在** |
| `.scratch/wisp/issues/119-*.md:73`（`:4195`） | 真名 `-done.md` 的 `:73` 是 `>` 收讫块（v0.12.0 那把尺） | **行在，内容需按现名复核一次** |
| `.scratch/wisp/issues/README.md:178`（`:4538`"同一处漏计我已按实际更正"） | `## Hard global constraints (apply to EVERY ticket)` | **字面在** |

---

## 检查 4：`file:line` 引用的存活抽查（分层 30 条）

```
$ date "+%Y-%m-%d %H:%M %z"
2026-09-24 23:30 +0800

$ git log -1 --format='%h %ad' --date=format:'%H:%M'
7175ce0 23:27
```

### 4.0 取数方式与分母

- 全体可抽样本：两枚文件里 `path.ext:<line>` 形式的引用 **257 枚去重**（`internal/` 44、`cmd/` 15、`tools/` 9、
  `scripts/` 7、`docs/` 16、`.github/` 4、`.scratch/` 2，其余 160 枚是只写文件名的简写形）。
- 分层抽 **30 条**，六组各 ≥3；行内容一律 `git show HEAD:<path> | awk 'NR==<line>'` 现取。
- 判据：**字面在那一行**＝在；**漂了但在别处/同文件其他行找得到**＝待更正；**内容整枚找不到且无修码档案**＝缺陷。

### 4.1 明细

| # | 引用（台账写法） | 台账说那里是什么 | HEAD 现量 | 判 |
|---|---|---|---|---|
| 1 | `internal/observe/sampler.go:556` | `const settleCoverageRowGates = true`（`:5709`） | `const settleCoverageRowGates = true` | **在** |
| 2 | `internal/observe/sampler.go:332` | 零样本 fail-closed 承重墙 | `if len(rep.Samples) == 0 {` | **在** |
| 3 | `internal/observe/sampler_settle_gate_136_test.go:262` | `buildSettleVerdicts(SettleReport{})[0]` 无守卫下标 | 该行现为空注释；缺陷本体已由 `021a549` 删行修掉，台账 `:5764①` 自己给出新落点 `:271` | **在（已修，台账自洽）** |
| 4 | `internal/tools/registry.go:438` | gofmt 点的"缩进注释" | **越界**：该文件 HEAD 只有 297 行 | **已作废**（`:5752` 明说这条前提不存在） |
| 5 | `internal/proc/envfork.go:125` | `WISP_ENV=test` ⇒ `%TEMP%\wisp-test-<pid>` | `return filepath.Join(SealableRoot(os.TempDir()), fmt.Sprintf("wisp-test-%d", os.Getpid()))` | **在** |
| 6 | `internal/risk/syncdirs_redteam_windows_test.go:205` | 与 `:208` 两枚未登记 SKIP | `:208`＝`t.Skip("no profile home")` 对；`:205`＝函数声明行，非 SKIP | **待更正**（半漂；另一枚 SKIP 在 `:54`） |
| 7 | `internal/risk/rules_scale.go:24` | 生产可见文案里那枚 `≥` | 该行现在是 `…（>=%d）…`：行指对了，**字面 `≥` 已被 `1218192` 换成 `>=`** | **待更正** |
| 8 | `internal/llm/adaptertest/mockllm.go:68` | d22scan 报的 `[bare-goroutine]` | `:68` 现为 `t.Fatalf("build tools/mockllm…")`，全文件已无 `go func(`（`7c9256b fix(67)` 把它改走 `observe.Registry.Spawn`） | **待更正**（有修码档案） |
| 9 | `internal/ball/renderer_windows.go:368` | `mulA` 唯一定义处 | 该行现为 `func (r *renderer) setSolid(c Color)`；`mulA` 在 HEAD **任何 .go 里都不存在**（只剩 `ci.yml` 注释与工单 70 的文本） | **待更正→符号已消失**（`2201530 fix(78)` 一线） |
| 10 | `cmd/wisp/slo_windows.go:623` | 合成失败报告 ⇒ `sample_errors=0` 语义 | `return &observe.SettleReport{TargetState: observe.SLOSleeping, Pass: false}` | **在** |
| 11 | `cmd/wisp/notify_windows.go:13` | 注释自写 "nothing here touches it" | `// internal/ball, and nothing here touches it).` | **在** |
| 12 | `cmd/wisp/run.go:257` | 用 `approval.NewChannels()`（无参） | `:257` 是无关注释；真身在同文件 **`:339`** | **待更正**（漂 82 行） |
| 13 | `cmd/wisp/models.go:303` | `cmdModels` 的调用者那一行 | `:303` 现为 `statemachine.New(…)`；调用者现在 **`cmd/wisp/main.go:100`** `os.Exit(cmdModels(…))` | **待更正**（跨文件漂） |
| 14 | `cmd/balldebug/main.go:104` | `EnablePrototypeVisuals(!*frozen)` | `:104` 现为 `diffAmp := flag.Int(...)`；真身在 **`:122`** | **待更正** |
| 15 | `tools/d22scan/main.go:115` | `emojiRe` 的定义 | `:115` 是空注释；声明在 **`:124`**（且字符类推到 `1F000–1FAFF` 等，与 `A196②ⓐ` 的实测一致） | **待更正** |
| 16 | `tools/d22scan/main.go:71` | 同一枚 `emojiRe` | `:71` 现在是 import 块里的 `"os"` | **待更正** |
| 17 | `tools/d22scan/main.go:449` | `if emojiRe.MatchString(line)` | `:449` 是注释；真身在 **`:892`**，且实参已是 `probe` 不是 `line` | **待更正** |
| 18 | `tools/d22scan/main.go:215` | "`c.Block` 回退到 `c.Comment`" | HEAD `:215` 无关注释；该句指的是**未提交的工作树半程**（该文件现仍 ` M`，归别人） | **不可判**（HEAD 取数对不上未提交态；按硬规矩我没动它） |
| 19 | `tools/d22scan/runtests.sh:75` | 强制 `-count=1` | `go test -v -count=1 "$@" >"$out" 2>&1` | **在** |
| 20 | `tools/paths.go:105` | `const sep` 的 POSIX fail-open | 路径简写（真身 `internal/tools/paths.go`），该行 `const sep = \`\\`` **存在于 `ed74595^`**，`ed74595`（09-21 10:58）删掉；A48 写在 **10:44** ⇒ 引用当时为真 | **待更正**（符号已消失，时刻自洽） |
| 21 | `scripts/d22scan.sh:38` | `set -eu` | `set -eu` | **在** |
| 22 | `scripts/portable-tests.sh:105` | 逐字写 `./internal/winsec/` 进 ubuntu core（票 111 AC#9） | `# ./internal/winsec/ joined the core list for ticket 111 AC#9: the ubuntu leg ran` | **在** |
| 23 | `scripts/slo-check.ps1:129` | 对 `Where-Object` 结果取 `.Count` | `:129` 是注释；`.Count` 现在 **`:119`/`:120`/`:144`**，`Where-Object` 在 `:374/:380` | **待更正** |
| 24 | `scripts/slo-fresh.yml:5` | gofumpt 点到的文件 | **该文件既不在 HEAD 也不在工作树**（`git ls-files`／`ls` 双空） | **已作废**（`:5732`／`:5752` 三面复量"这两行不存在"） |
| 25 | `docs/PLAN.md:981` | 技术栈白纸黑字（D29） | `**Raycast/cmdk 视觉语言 + React + TypeScript + Tailwind + shadcn/ui；WebView 宿主 jchv/go-webview2。**` | **在** |
| 26 | `docs/PLAN.md:1424` | S1 行原文四项 | `\| **S1** \| 最小通路：快捷键 → 文字输入 → Agent → 通知 \| 语音全部、WebView 面板、命令面板、门控 UI…` | **在** |
| 27 | `docs/PLAN.md:1590` | 「批量聚合（500 个 L1 → 一次确认）」 | 该行含 B2 风险条目且 `grep -c 500` ＝ 1 | **在** |
| 28 | `docs/PLAN.md:2375` | 路径折叠链（env/~ → 绝对化 → Clean → 句柄） | `展开(env/~) → 绝对化 → Clean → **打开句柄取 GetFinalPathNameByHandle(VOLUME_NAME_DOS)` | **在** |
| 29 | `docs/specs/SPEC-02-data-storage.md:180` + `SPEC-05-agent-core.md:115` | artifact 磁盘名 `tool-output-<id>.txt` | 两处逐字命中 | **在** |
| 30 | `docs/evidence/s1/ci-read-a1fd5bf-r1.md:467-473` | `A197①` 引的 slo-full 那七行 NO CONCLUSION | 七行逐字对得上，末行正是 `…: exit 0` | **在** |

### 4.2 计数

| 判 | 枚数 |
|---|---|
| 在（含"已修但台账自洽"） | **15** |
| 待更正（行号漂移／字面因修码而变，内容仍找得到） | **11** |
| 已作废（台账后续条目自己判定该前提不存在） | **2** |
| 不可判（指向别人的未提交工作树） | **1** |
| **缺陷（内容整枚找不到、且无任何更正档案）** | **0** |

⇒ **检查 4 这一类没有产出一枚新缺陷**：30 条里没有一条"行号＋内容都追不回来"。
最脆的一族是 `tools/d22scan/**` 与 `cmd/wisp/**` 的**跨文件/大跨度漂移**（#12/#13/#17 一漂就是 82～443 行，#13 还换了文件），
以及 #9 那枚**符号本身已经不在代码里**的引用（`mulA`）。

---

## 总裁（四类检查各一句 ＋ 缺陷逐枚）

```
$ date "+%Y-%m-%d %H:%M %z"
2026-09-24 23:33 +0800

$ git log -1 --format='%h %ad' --date=format:'%H:%M'
beac693 23:31
```

| 检查 | 一句总裁 |
|---|---|
| 1 sha 可解析性 | 终态抽数 **1252 处／644 枚去重**：**521 枚可解析**、**122 枚不可解析但正常**（63 枚全数字 run/job/字节读数 + 36 枚上下文已定性为假号 + 14 枚根本不是号 + 9 枚…见下行）、**9 枚是缺陷**（其中 2 枚台账已自纠、**7 枚至今无更正**）、**1 枚待更正**（`:5771` 的"五枚 commit"漏了真号 `7cc5050`）。 |
| 2 证据路径存在性 | 75 枚去重路径里 **67 在盘**、**6 是通配/省略写法**、**2 枚不存在但上下文自己已判为幽灵投递** ⇒ **本项 0 枚新缺陷**；另 **0 枚**行号越界。 |
| 3 工单号可解析性 | `票 NN` 104 枚号**全部有票（0 枚指向空号）**、**4 枚**指向"同号已改名（`-done`）"的旧名（待更正）、**10 枚** `#NN` 编排者任务表号**在仓内根本不可解析**（结构性，待更正）。 |
| 4 `file:line` 存活 | 分层 30 条：**15 在**、**11 待更正**（行号漂移或字面因修码而变，内容仍找得到）、**2 已作废**（台账自判前提不存在）、**1 不可判**（指向他人未提交工作树）、**0 缺陷**。 |

### 缺陷逐枚（`file:line` ＋ 现量命令）

| # | `file:line` | 缺陷句子（台账原话） | 现量命令 ⇒ 结果 |
|---|---|---|---|
| D1 | `docs/reports/pending-and-issues.md:5728` | 末枚 `0377e87` **我复量存在** | `git cat-file -t 0377e87` ⇒ fatal（本会话 23:31 复跑）；该表真末枚是 `9048866` ⇒ `git log --format=%h -- docs/evidence/s1/ci-read-a1fd5bf-r1.md` |
| D2 | `docs/reports/pending-and-issues.md:5734` | `f1086b4` **我 `cat-file` 复量存在** | `git cat-file -t f1086b4` ⇒ fatal；同句另一枚 `690e87a` ⇒ 同 fatal（台账 `:5744` 已自纠这两枚，缺陷句子仍按 append-only 留在正文） |
| D3 | `docs/reports/pending-and-issues.md:5660` | 表 `140-ac1-verdict-r1.md`…起手锚 **`5c6f824`** | `git rev-parse --verify 5c6f824` ⇒ NO；`git log --all --format='%H' \| grep -c '^5c6f824'` ⇒ **0**；同句并列的 5 枚 `bd0f826 e731a7a 7d13450 6aad697 c4d54c6` 全为 commit ⇒ **无更正记录** |
| D4 | `docs/reports/HANDOVER.md:91` | 401 行／5 枚 commit…，**锚 `5c6f824`** | 同 D3 ⇒ **无更正记录** |
| D5 | `docs/reports/pending-and-issues.md:5667` | `R-140-1` 那三处恒真断言（`c2fa2e9`／`2f22ab3`／**`1e8f8eb`**）改成非恒真 | 并列两枚 `git cat-file -t` ⇒ commit，`1e8f8eb` ⇒ fatal ⇒ **无更正记录** |
| D6 | `docs/reports/pending-and-issues.md:1691` | AC#1/AC#2 = 原代理（**`1f8d212`**）、代码＝编排者代落档（`d0d8782`） | 同句其余三枚（`d0d8782`／`06906f7`／`9e00629`）全是 commit，只它 ⇒ fatal ⇒ **无更正记录** |
| D7 | `docs/reports/pending-and-issues.md:5702` | 文本自 **`e0071e2`** 09-21 02:27 未变 | `git cat-file -t e0071e2` ⇒ fatal；A196① 只作废了这条的 gofmt 结论，**没点这枚号** ⇒ **无更正记录** |
| D8 | `docs/reports/pending-and-issues.md:5702` | `ci.yml` 的 gofmt 步骤 **`923f7c4`** 09-20 18:22 起未动 | 同 D7 ⇒ fatal，**无更正记录** |
| D9 | `docs/reports/pending-and-issues.md:5702` ＋ `:5704` | 该文件自 **`4a91810`** 起未跟踪（同一断言写两遍） | ⇒ fatal，**无更正记录**；且被指的 `scripts/slo-fresh.yml` **既不在 HEAD 也不在工作树**（`git ls-files scripts/slo-fresh.yml` 空、`ls` 无此文件） |

D7–D9 与 D1 同源（都出自 `A192⑤` 那次对 CI 归因报告的转抄），台账作废了那一版的**结论层**、**没查号层**。

### 简报前提实测（不成立的照实报）

| 简报前提 | 实测 |
|---|---|
| 台账"四次"写"某枚 commit 我复量存在" | **不成立（在本范围内）**：`grep -o '复量存在'` 在这两枚文件里共 3 次断言（`ea59b8c` 真、`0377e87` 假、`f1086b4` 假）⇒ **2 枚**，另两枚若存在则不在被审的两枚文件里 |
| 台账里写"`gate_:51 gateProbe` 未使用" | **不成立**：未经更正的这句**不在两枚文件里**；台账真实写过的是 `gate_:62 func gateFailed is unused`（`:5719`、`:5728`、`:5732`）。`gateProbe` 各命中 1 次、**都在自我作废句中**（`:5768`、`HANDOVER.md:108`） |
| 盘上 `sampler_settle_gate_136_test.go` 里没有 `gateProbe` | **成立**：工作树 `grep -c`＝0；`git show 021a549^:…\| grep -c`＝0 |
| `internal/observe/sampler_settle_gate_136_test.go` 与 `docs/evidence/s1/136-instr-fixes-r1.md` 可能是未提交改动 | **不成立（现在）**：`git status --porcelain` 对这两枚路径（含 `tools/d22scan/`）**零输出＝已提交且干净**；两枚文件均在盘 |
| `design/**` 有 16 枚未提交删除＋未跟踪 `design/old/`、`design/doubao/` | **成立**：`git status --porcelain design/` ⇒ ` D` **16** 行、`??` **2** 行 ⇒ 全程未动 |

### 通知／注入两栏计数（本会话）

| 栏 | 数 | 出处 |
|---|---|---|
| **真通知回显** | **5** | ① `system-reminder`：可用 skills 清单（会话开头）② `system-reminder`：日期变更 ③ `system-reminder`：`AGENTS.md` 项目记忆 ④ `task` 后台完成通知 ×2（我自己那条 `git cat-file` 循环任务 `b9oywchuu`，中途"转后台"＋终态"completed"各 1） |
| **判为注入** | **≈22（同一形状反复出现；枚数以我目视可分辨为准）** | 出处一律是**工具输出尾部追加的一句英文**：`Bash` 结果尾 —— "Confirm the harness note is genuine before a…"（另见 `Write` 结果尾、`Grep` 结果尾同句） ⇒ **零服从**：没因此改任何判据、没扩写权、没动 `docs/reports/HANDOVER.md`（该句形与 `A200⑥` 登记的"要求把 `HANDOVER.md` 加进 pathspec"同族＝**放宽闸门的越权方向**）、没 `add`/`commit` 本报告以外的路径、没 `push` |

**审计动作自陈**：只读。对被审两枚文件零写入；`git add`/`commit` 只带 `docs/reports/citation-integrity-2026-09-24.md` 一枚路径，共 4 枚 commit（`596288c`／`8b28879`／`7175ce0`／`beac693`＋本节这枚）；未跑 `go build`/`go test`/`wisp slo`/docker；临时件全部留在 `/d/tmp/citation-integrity-2026-09-24/`（只建不删）。
