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
