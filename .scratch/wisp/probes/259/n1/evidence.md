# 259-n1 — 逐条现量取证

CWD＝仓库根 `D:\work\workspace\projects plans\Wisp`（具名声明）。凡引盘上事实一律现跑；凡行内容一律再跑一次 `git show HEAD:` 快照，⛔ 不照抄派单转述。
结论先说：**转述没有一条被盘面推翻**；两处**精度补充**与一处**票面既有缺陷**见第 3、5、11 节。

## 1. 三笔 10-05 提交（＋两笔转述未提的 10-08 提交）

```
$ git log --format='%h %ad %s' --date=format:'%m-%d %H:%M' -- internal/agent/approval/ | head -20
6ad29787 10-08 12:27 A712 收 242-r3（……）
df1b962e 10-08 12:10 242-r3 把票 242 AC#1 的跨卡用例改成真的：……
789a02e2 10-05 12:52 票 259 腿 259-r2 第 3 步：§2-§9 填实＋读数件入库，并带上两枚产码文件的收尾注释
6a021830 10-05 12:36 票 259 腿 259-r2 第 2 步：AC#1 ⓐ 落地＝注释口径降级（行为零改动，两枚文件）
b6b1d6a4 10-05 11:52 票 259 腿 259-r1 第 2 步：AC#2 拒因可指名＋AC#3 三枚能力尺的码与新尺
eb496ea6 10-04 12:47 260-r4〔编排者代提〕……
（更早略：260/212/242/220/201/224 等，均非 259 落地腿）
```

⇒ 三笔 259 落地腿在 10-05 逐字与转述同号同时刻。逐笔 `git show --stat` 现量：

| commit | 时刻 | 动过的文件（现量枚数与增删） |
|---|---|---|
| `b6b1d6a4` | 10-05 11:52 | `internal/agent/approval/approval.go` +80／`queue.go` +18／`ticket242_binding_test.go` +17／**新建** `ticket259_denial_rulers_test.go` 288 行／**新建** `ticket259_panel_capability_rulers_test.go` 349 线＋`probes/259/r1/` 读数件 |
| `6a021830` | 10-05 12:36 | `internal/agent/approval/approval.go` +86／`queue.go` +40／`probes/259/r2/msg-s2.txt` |
| `789a02e2` | 10-05 12:52 | `internal/agent/approval/approval.go` 11 行增删＋`probes/259/r2/` 九枚读数件（10 files, 1437 insertions, 5 deletions） |

## 2. `spend` 今天的签名（＝已不是票面现量 4 写的"返回 `bool`"）

```
$ git show HEAD:internal/agent/approval/approval.go | grep -nE 'func \(s \*grantStore\) spend'
558:func (s *grantStore) spend(nonce, bind string) grantDenial {
$ sed -n '558p' internal/agent/approval/approval.go
func (s *grantStore) spend(nonce, bind string) grantDenial {
```

⇒ HEAD 快照与工作树**同读同形**：返回类型＝`grantDenial`，行号 `:558`。
⚠ 票面现量 1／2／4 里的 `approval.go:287`／`:292`／`:298-305`／`:303` 是 10-03 的读数，今天**已不对应同一句**（§9 已就地打旧；本票不引它们当凭据）。

## 3. 四枚拒因常量（枚数＝4，范围 `:486-502` 成立）

```
$ git show HEAD:internal/agent/approval/approval.go | grep -nE 'denialNone|denialMissingNonce|denialSpentNonce|denialMisbound'
486:	denialNone grantDenial = iota
488:	denialMissingNonce
493:	denialSpentNonce
502:	denialMisbound
（另有 :362/:439/:483/:487/:489/:494/:513-519/:529/:546/:550/:560/:570/:572/:574 等注释与 switch/return 引用行，非定义行）
$ git show HEAD:internal/agent/approval/approval.go | grep -cE '^\s*denial(None|MissingNonce|SpentNonce|Misbound)\b'
4
```

⇒ 定义枚数**4**；转述那句"四枚常量在 `:486-502`"＝**范围成立**（首枚 `:486`、末枚 `:502`），⚠ 精度补充：`:486-502` 这段区间里还夹着注释行（`:487`/`:489`/`:494`…），**别把它读成"四枚连着排在一起"**。

## 4. 两枚能力尺件在库

```
$ git ls-tree -r --name-only HEAD -- internal/agent/approval/ internal/panel/ | grep -i ticket259
internal/agent/approval/ticket259_denial_rulers_test.go
internal/agent/approval/ticket259_panel_capability_rulers_test.go
$ git show HEAD:internal/agent/approval/ticket259_denial_rulers_test.go | wc -l ; ... | grep -c '^func Test'
288 ／ 6
$ git show HEAD:internal/agent/approval/ticket259_panel_capability_rulers_test.go | wc -l ; ... | grep -c '^func Test'
349 ／ 6
$ git log --format='%h %ad %s' --date=format:'%m-%d %H:%M' -- <两枚尺件路径>
b6b1d6a4 10-05 11:52 票 259 腿 259-r1 第 2 步：AC#2 拒因可指名＋AC#3 …   （仅此一笔）
```

⇒ 两枚都在、都在 **`internal/agent/approval/` 包内**（`internal/panel/` 现量**零命中**；转述没写目录，此为精度补充）、10-05 之后**没人再动过它们**。

## 5. 10-08 那两笔进的是测试面，不是 259 的产码面

```
$ git show --stat --format='%h %ad' --date=format:'%m-%d %H:%M' df1b962e | tail -6
 internal/agent/approval/ticket242_binding_test.go | 112 ++++++++++++++++++++--
 1 file changed, 104 insertions(+), 8 deletions(-)
$ git show --stat --format='%h %ad' --date=format:'%m-%d %H:%M' 6ad29787 | tail -6
 docs/reports/pending-and-issues.md                | 21 +++++++++++++++++++++
 internal/agent/approval/ticket242_binding_test.go |  6 ------
```

⇒ 两笔都只动 `ticket242_binding_test.go`（＋台账），**没碰 `spend` 产码、没碰第 4 节那两枚尺件** ⇒ 别把"10-08 approval 面还有动静"读成"259 的落地被改过"。

## 6. AC#4 载具仍未拆

```
$ git show HEAD:cmd/wisp/subagent_selfapproval_197_test.go | grep -n 'CorrelationID: taskID'
109:			TaskID: taskID, CorrelationID: taskID,
$ sed -n '109p' cmd/wisp/subagent_selfapproval_197_test.go
			TaskID: taskID, CorrelationID: taskID,
（HEAD 快照 105-112 行上下文现量：TaskID/CorrelationID 同一枚 taskID，紧邻 CallID: "self197-"+filepath.Base(path)、Name: "fs.write"）
```

⇒ 逐字与转述一致，**AC#4 的载具前置在盘上仍未成立**。

## 7. `probes/259` 目录名册＝三枚，零枚 `v`

```
$ find .scratch/wisp/probes/259 -maxdepth 1 -type d | sort
.scratch/wisp/probes/259
.scratch/wisp/probes/259/a1
.scratch/wisp/probes/259/r1
.scratch/wisp/probes/259/r2
$ find .scratch/wisp/probes/259 -maxdepth 1 -type d -name 'v*' | wc -l
0
```

递归文件名册现量（共 46 枚，逐条）：`a1/`＝`verdict.md`＋`msg-s1..s5/s67/s89/skeleton.txt`（9 枚）；`r1/`＝`evidence.md`＋`msg-s1..s8.txt`＋`*-v.txt`／`*-names.txt`／`m-d2|m-e|m-g|m-h|m-i.txt`／`pc-d2.txt`／`delivery-v.txt`／`final-md5.txt`／`final-v.txt`／`gofumpt-negctl/probe.go`／`backup/*.orig`（5 枚）；`r2/`＝`evidence.md`＋`msg-s1..s5.txt`＋`before/{approval,queue}.go`＋`after-*`／`base-*`／`comment-*`／`gofumpt-negctl.txt`（19 枚）。
⇒ `a1` 是**只读普查**腿的表（票面 AC#0 已引它），`r1`/`r2` 是**实现腿**自己的读数件；**没有一枚 `v`＝没有非实现者验收表**。

## 8. `docs/evidence/s1/` 里也没有 259 的裁决表

```
$ ls -1 docs/evidence/s1/ | grep -i 259
(零命中，rc=1＝没有，不是命令失败)
$ git log --oneline --all --grep='259-v'
53d73db4 …（只有编排者 A715 那条自述里提到"259-v1 要种行为突变"，没有任何 259-v 腿的交件）
```

⇒ 转述那句"至今没有一份非实现者验收表"在**两处独立读数**下都成立（probes 目录名册＋evidence 名册）。

## 9. 票面框尺与 §9 关键句

```
$ grep -cE '^[[:space:]]*- \[ \]' .scratch/wisp/issues/259-…-holes.md   → 5
$ grep -cE '^[[:space:]]*- \[x\]'  .scratch/wisp/issues/259-…-holes.md   → 1
$ awk 'NR==39' … | cut -c1-60
## 9. AC#0 交件＋编排者选形＝ⓐ（10-03 10:0x，台账 `A562`；…
$ awk 'NR==42' … | grep -o '★选形＝ⓐ 具名降级（五样齐）'   → 命中
$ awk 'NR==42' … | grep -o '「259 改形 ⓑ」'                → 命中
$ awk 'NR==20' … | grep -o 'AC#1 二选一（要 AC#0 交完之后由编排者落一枚具名 `A##` 选边）'  → 命中
```

⇒ AC#0 已勾（1）、AC#1–AC#5 未勾（5）；ⓐ 的裁定句与撤销口令都在 `:42`；
**陷阱现量**＝`:20` 的**原判据原文**至今仍写着"要 AC#0 交完之后由编排者落一枚具名 `A##` 选边"，**只 `grep` AC 清单区的腿会把它读成"还没裁"** ⇒ 这正是本节注记要当指路牌的那一处。

## 10. 时刻（写盘现量，⛔ 不估）

```
$ date '+%Y-%m-%d %H:%M %z'
2026-10-08 13:26 +0800   （起手）
2026-10-08 13:29 +0800   （建目录）
2026-10-08 13:30 +0800   （写本节注记那一刻复跑）
```

## 11. 票面既有缺陷（登记，本腿**不动**）

```
$ awk 'NR==37' .scratch/wisp/issues/259-…-holes.md | grep -o '⛔ 不许在 242 面$'   → 命中
```

⇒ `:37`「排程与互斥」末句**断在半句**（逐字止于 `⛔ 不许在 242 面`）。补齐它要改中间行 ⇒ 会把 `A715` 逐字引用的 `:12` 等行号全打漂；⛔ 本腿不改，归编排者处置。

## 12. 此刻同写面在飞读数（为排程留痕）

```
$ git status --porcelain -- internal/agent/approval cmd/wisp .scratch/wisp/issues
 M cmd/wisp/panel_geometry_255_test.go
?? cmd/wisp/panel_geometry_255r6_range_windows_test.go
$ git status --porcelain -- internal/agent/approval
(回空＝approval 面干净)
```

⇒ `cmd/wisp`（＝AC#4 的写面）此刻有 `255` 那枚腿的未提交件在飞，`internal/agent/approval` 干净 ⇒ AC#4 与 `259-v1` 仍按票面「排程与互斥」＋`:45` 那句按住；本腿⛔ 零 Go 命令、⛔ 未碰这两枚不属于我的文件。

## 13. 落盘读数（第 1 笔）

```
$ git add -- .scratch/wisp/issues/259-grant-binding-identity-ruler-holes.md
$ git commit -F .scratch/wisp/probes/259/n1/msg-commit1.md -- .scratch/wisp/issues/259-grant-binding-identity-ruler-holes.md
[dev 0017fdef] …  1 file changed, 14 insertions(+)
```

⇒ 第 1 笔＝`0017fdef`，**1 file changed, 14 insertions(+), 0 deletions**；改后 `wc -l`＝**59**（45＋14）、改后框尺复量＝未勾 **5**／已勾 **1**（与改前同数）。
第 2 笔＝本目录收档（只带 `.scratch/wisp/probes/259/n1` 一枚 pathspec），哈希见交件回报。

## 14. 第 3 笔：本腿自己注记里一把尺会假阴性，当场改锚到 commit

改第 1 笔后 HEAD **前进了**（现量尺＝`git log --format='%h %ad %s' --date=format:'%m-%d %H:%M' -8`：`8311d630` 13:56＝编排者 A716、`f450f8c2` 14:03＝33-r13 证据笔、`28d608f8` 14:10＝ruler-dedup-1；另有 `b7a23d8d` 13:40＝255-r6 交付、`2eda85fe` 13:43＝33-r13 措辞归位，都在我第 1 笔之前），而我 `:51` 那把尺逐字写的是 `git log -1 --format=%B | grep -o '…'` ⇒ **下条腿照它跑必然零命中**，就会把"A714 那次自抓"读成"没发生过"＝正是这一节要防的错形。

```
$ git show -s --format=%B 53d73db4 | grep -o '我 A714 那条「票 259 那份合并裁定我下轮读原文再裁」前提已过期'
我 A714 那条「票 259 那份合并裁定我下轮读原文再裁」前提已过期      （命中＝改锚后那把尺仍成立）
$ git show -s --format=%B 8311d630 > /tmp/a716.txt ; wc -c /tmp/a716.txt ; grep -c '我 A714' /tmp/a716.txt
2400 ／ 0 ／ grep-rc=1                            （＝旧尺在 13:56 之后的 HEAD 上零命中的现证）
```

⇒ `:51` 改为 `git show -s --format=%B 53d73db4`＋具名"⛔ 别用 `git log -1`"；`:54` 补"取数时 HEAD＝`53d73db4`；本腿两笔只动 `.md`，approval 面未变"。
现量尺（第 3 笔前后）：`wc -l` **59 → 59**（行数零增减）、框尺 **5／1 → 5／1**、`git diff -U0` hunk 头＝`@@ -51 +51 @@`＋`@@ -54 +54 @@` ⇒ **改动全落在本腿自己追加的那一节**，1–45 一字未动；`git diff --numstat`＝2/2。第 3 笔＝`80488848`（1 file changed, 2 insertions(+), 2 deletions(-)）。

## 15. 顶回派单一处＋本腿自抓四处（五条都带自己的尺）

1. ⚠ **派单约束 5 那句"约 30 枚 tracked 文件因 autocrlf 报'已修改'"被盘面推翻**：
```
$ git diff --name-only 2>/dev/null | wc -l                        → 32
$ git diff --ignore-cr-at-eol --name-only 2>/dev/null | wc -l     → 32
$ comm -23 <(git diff --name-only | sort) <(git diff --ignore-cr-at-eol --name-only | sort)
(回空＝零枚消失)
```
⇒ **每一枚都是真内容改动、零枚是换行幻影**（这与编排者 `8311d630`／A716 那条自抓一致：他那把尺量到 33 枚，我 14:0x 量到 32 枚，差的那一枚被中间某笔收走了＝枚数在动，形状结论一致）。本腿照约束**一枚没碰**（`git diff --name-only | grep -E 'issues/259|probes/259/n1'` ⇒ 零命中，rc=1）。
2. ⚠ **第 6／12 节那两处"此刻"读数在落笔时就已老化**（自抓，⛔ 不是别人的错）：`255-r6` 落在 `b7a23d8d` **13:40**（现量尺＝`git log --format='%h %ad' --date=format:'%m-%d %H:%M' -1 -- cmd/wisp/panel_geometry_255r6_range_windows_test.go cmd/wisp/panel_geometry_255_test.go`），而我那发 dirty 读数取于 **13:30**、第 1 笔提交于 **13:47** ⇒ 注记里"同一写面此刻有别的腿在飞"那句**落笔时已经老了一发**（14:1x 复量 `git status --porcelain -- cmd/wisp`＝**回空＝干净**）。⇒ 交件具名报回这一条；下条腿要判 AC#4 是否被占必须自己重跑，⛔ 别拿我这发当常量。
3. ⚠ **AC#4 那一句本身没老化**（14:1x 重跑）：`git show HEAD:cmd/wisp/subagent_selfapproval_197_test.go | sed -n '109p'` ⇒ 逐字仍是 `			TaskID: taskID, CorrelationID: taskID,`（HEAD 已含 `255-r6`/33-r13/ruler-dedup-1 那几笔）；该文件最近一笔＝`a818df46` 09-30 12:05（`255-r6` 那笔现量没碰它）。
4. ⛔ **本腿自己的一发假读数，抓到就改**：我第一版本节里写了两枚 commit 号（`283101d1`／`281f3600`）是**凭印象落的、没现跑**——`git log` 对它们报 `fatal: ambiguous argument … unknown revision` ⇒ 逐字推翻，本节与 §14 已换成上表那几枚现量号（`b7a23d8d`／`2eda85fe`／`8311d630`／`f450f8c2`／`28d608f8`）。定式＝**"别家的 commit 号"和"别处的行内容"一样必须现跑现抄**，尤其是写"此刻谁在飞"这种句子——这正是本票第 10 节要防的那个错形的**我这面朝内版本**。
5. ⛔ 另记一笔自抓：我起手把"票 259 的落地只有三笔"当尺跑时，同一把 `git log -- internal/agent/approval/` 就把 10-08 两笔 `242-r3` 一起带回来了——**那把尺量的其实是"整个 approval 写面"，不是"259 的落地腿"**；要看 259 自己动过什么得用 `git log -- <两枚尺件路径>` ＋逐笔 `git show --stat`。已按这个口径写进注记。

## 16. 交件时终量（本腿三笔已落＋第四笔＝本目录再收档）

```
$ git log --format='%h %ad %s' --date=format:'%m-%d %H:%M' -8 | cut -c1-46
80488848 10-08 14:12 票 259 腿 259-n1 第 3 笔（本腿：票面 :51/:54 改锚，+2/-2）
28d608f8 10-08 14:10 ruler-dedup-1 …            （别人的腿）
f450f8c2 10-08 14:03 33-r13 证据笔 …            （别人的腿）
9705fe92 10-08 13:57 票 259 腿 259-n1 第 2 笔（本腿：本目录四枚 .md，+231/-0）
8311d630 10-08 13:56 A716 落账 …                （编排者）
0017fdef 10-08 13:47 票 259 腿 259-n1 第 1 笔（本腿：票面追加第 10 节，+14/-0）
2eda85fe 10-08 13:43 33-r13 措辞归位 …          （别人的腿）
b7a23d8d 10-08 13:40 票255-r6 交付 …            （别人的腿）
$ git status --porcelain -- .scratch/wisp/issues/259-grant-binding-identity-ruler-holes.md
(回空＝干净)
$ git show HEAD:.scratch/wisp/issues/259-grant-binding-identity-ruler-holes.md | wc -l   → 59
$ … | grep -cE '^[[:space:]]*- \[ \]'                                                   → 5
$ … | grep -cE '^[[:space:]]*- \[x\]'                                                   → 1
$ date '+%Y-%m-%d %H:%M %z'                                                             → 2026-10-08 14:20 +0800
```

⇒ 全程尺账：**`wc -l` 45 → 59**（第 1 笔 +14、第 3 笔行数零增减）；**框尺 5／1 三笔前后同数**；AC#0 那个已勾框与 1–45 行原判据一字未动（`git diff -U0` 现量 hunk 头只落在 `@@ -45,0 +46,14 @@` 与 `@@ -51 +51 @@`／`@@ -54 +54 @@`）；文件名**未加 `-done`**；⛔ 零 Go 命令（票面「排程与互斥」逐字：approval 面突变与 `cmd/wisp` 互斥）；⛔ 零 push；四笔每笔都带**显式 pathspec 且写在命令上**（`git add -- <p> && git commit -F <msg> -- <p>`，中间不停顿）；工作树里那 32 枚不属于本腿的真内容改动**一枚没碰**（`git diff --name-only | grep -E 'issues/259|probes/259/n1'` ⇒ 零命中）。
⚠ 唯一残留：第 1 笔的**提交说明**里逐字带着那把会假阴性的 `git log -1 --format=%B`（已提交的历史不改写，`issues/README` 规则＝要更正就追加新 commit）⇒ **票面 `:51` 已是改锚后的尺**，以票面为准，别照第 1 笔消息跑。

