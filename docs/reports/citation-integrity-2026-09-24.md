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
