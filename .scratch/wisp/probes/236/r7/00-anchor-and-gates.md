# 236-r7 — 票 236：钉内 "run.go 行号族" 整族刷新（写腿／实现者）

> 判语归后续非实现者腿；⛔ 本腿未翻任何框、未动 `Status:`／票名、零 push。
> 落点＝**单件** `internal/tools/taskid_source_236r4_teeth_test.go`。裁定出处＝账 `A923`／派单归档
> `98546a1d:.scratch/wisp/probes/orch/20261013/txt/236r7-dispatch.txt`（逐字＝本腿收到的派单）。
> 起手锚 `11cfbc93`（13:48:00）；执行窗内 HEAD 由编排者/他腿移到 `f1676bfe`（零与本件重叠：`scripts/*.sh`＝317-r1 在飞，未碰未读）。

## §0 结论先行

- **整族 11 处全刷**（r6 顶回表列注释 4 处 `:17-19`/`:136`＋字符串 3 处 `:205-206`/`:244`（r6 时编号）＋②票 224 pin `:28`；**现读全族实为 11 处** ⇒ 多出 3 处一并刷并具名＝`:228`/`:235`/`:236`，见 §2 的 ★）。⛔ 断言逻辑/锁/id 集合一字未动。
- **只数字变**：`git diff --numstat`＝`11 11`；判据①（剥掉 `[0-9]+` 后两版逐字同）**ruler rc=0**（空输出）；11 对改动行逐对核＝`STRIP-EQUAL` 11/11。
- **非注释改动行位置＝4 处**（`:208`／`:209`／`:228`／`:247`），全部＝"字符串里仅数字变"的 `t.Errorf/Fatalf` 文字行（判据⑥逐行点名；-/+ 两侧共 8 行）。
- **门禁全绿**：`-run 'Test236'` 前 rc=0（2.042s）／后 rc=0（3.134s）；`-race` rc=0（4.193s）；gofumpt 空（前/后）；
  `go build ./...` rc=0（前/后）；d22scan clean rc=0（正控 35/0/0）；长度门 GREEN（over=59／roster=59／not-in-roster=0）。
- **★一处与裁定列表的具名差异（由判据①本身逼出）**：r6/A923 列的 `:1139-1142`／`:1156-1160`／`:441-443`
  是**块范围**（实体所在区域的起止）。被引**文本行**本身的现读＝`1142`／`1156`／`442`——与 r4 被引行**逐字同串**（logs/10 §same-text）。
  按"只许数字 token 变"＋判据①（剥数字后逐字同）：**单号只能引单号、范围只能引范围**。把 `1124`/`433`
  写成 `1156-1160`/`441-443` 会**新增连字符**（非数字字符）⇒ 判据① 会有 **6 行**不通过（`:18`/`:19`/`:209`/`:228`/`:236`/`:28`）。
  故按派单规矩 3"以现读为准"取精确单号；块范围信息原样保留在 §2 备注列。若编排者要的是范围写法，须明示改判
  （代价＝判据① 口径须同步改，本腿不自行放宽任何判据）。

## §1 现读重跑（起手全部现读，不引用顶回数）

- run.go blob `154279f7`（worktree＝HEAD 同串，r6 窗口后**未再漂**）；ticket224 blob `ebe020dd`；loop.go blob `0d53277d`。
- 被引行现读：run.go `:1137 rt.tasks.MarkRoot(bg.ID, task)`；`:1142 }).Backfill(bg.Root().Ctx, res.TaskID, res.Text)`；
  `:1156 ID: res.TaskID, State: state, QueryText: task,`；`:1167 if err := rt.store.FinishTaskLog(c, res.TaskID, state, summary,`；
  `:1177 res.TaskID, res.Status, res.Rounds, res.ToolCalls, costLine(res))`；块＝`:1139-1142`（Backfill 复合字面量起→调用行）。
  ticket224 `:442 t.Fatalf("boot 1 minted %q, which the shape lock rejects: …")`（块＝`:441-443`，`:441`＝`if !…Valid()`、`:443`＝续串）。
- 同串核对（logs/10 §same-text）：`1105→1137`、`1124→1156`、`1135→1167`、`1145→1177`、`1107-1110→1139-1142`（四行块）、
  ticket224 `433→442` **全部 IDENTICAL**（r4 blob `f88d98ae` 逐字比对）。

## §2 逐处新旧行号表（新号＝现读精确单号；r6/A923 块范围记在备注）

| # | 位置（现读 pin 行） | 类型 | 旧 | 新 | 备注（A923/r6 块范围） |
|---|---|---|---|---|---|
| 1 | `:17` | 注释 | `run.go:1105` | `run.go:1137` | 1137=MarkRoot |
| 2 | `:18` | 注释 | `run.go:1110` | `run.go:1142` | 块 :1139-1142 |
| 3 | `:19` | 注释 | `:1124/1135/1145` | `:1156/1167/1177` | 1156∈:1156-1160 |
| 4 | `:28` | 注释（票 224 pin） | `ticket224…:433` | `ticket224…:442` | 块 :441-443 |
| 5 | `:139` | 注释 | `run.go:1107-1110` | `run.go:1139-1142` | 范围→范围（原就带连字符） |
| 6 | `:208` | **字符串** | `run.go:1105` | `run.go:1137` | （-/+ 两侧，判据⑥点名） |
| 7 | `:209` | **字符串** | `:1110/:1124/:1135/:1145` | `:1142/:1156/:1167/:1177` | （-/+ 两侧） |
| 8 | ★`:228` | **字符串** | `run.go:1105/:1110` | `run.go:1137/:1142` | r6 表未列；同族第三处 |
| 9 | ★`:235` | 注释 | `(run.go:1105 then` | `(run.go:1137 then` | r6 表未列 |
| 10 | ★`:236` | 注释 | `// :1110):` | `// :1142):` | r6 表未列 |
| 11 | `:247` | **字符串（Fatalf）** | `run.go:1105` | `run.go:1137` | （-/+ 两侧） |

★＝本腿现读扩出的 3 处（同族引用，机械一致性的必需项；枚数＝**11 处/18 枚号码 token**）。
⛔ 未刷新、具名理由：`:10` 的 `:892`/`:1099` 是**引票面/a4 文字**的过时引用（句内自注"have drifted; re-check before quoting"），
照刷新会改变语义（那些数字引的是**票面/census 的原文**，不是 run.go 现值）；loop.go 族（`:329/330/333/382/676/1153/1155-1161/1162-1163`）
本腿逐枚现读仍准（`0d53277d`），且不属本族。

## §3 判据① "仅数字变" 尺读数（logs/11）

- `diff <(git show HEAD:<件> | sed -E 's/[0-9]+/#/g') <(sed -E 's/[0-9]+/#/g' <件>)` ⇒ **空输出 rc=0**（剥掉数字后两版逐字同）。
- `git diff --numstat` ⇒ `11 11`；`-U0` 每对 `-`/`+` 逐对剥数字比对 ⇒ `STRIP-EQUAL 11/11`。
- 非注释改动 = 4 处字符串行（`:208`/`:209`/`:228`/`:247`），每行仅数字变（逐行在 logs/11 点名）。
- 行数 278→278；sha256 `7e23ec32…`→`00fd832e…`。

## §4 门禁（logs/12）

| # | 判据 | 读数 |
|---|---|---|
| ② | `-run 'Test236' ./internal/tools/ -count=1` | 前 rc=0（2.042s，8 枚全 PASS）／后 rc=0（3.134s） |
| ③ | `-race` 一发 | rc=0（4.193s） |
| ④ | gofumpt 该件 | 前/后 `-l` 皆空 rc=0 |
| ⑤ | `go build ./...` | 前/后均 rc=0 |
| ⑥ | 非注释行点名 | 仅 `:208`/`:209`/`:228`/`:247` 四处"字符串里仅数字变" |
| 补 | d22scan／长度门 | clean rc=0（正控 PASS=35 FAIL=0）／GREEN over=59 roster=59 not-in-roster=0 |

## §5 窗口逐报（他腿在飞）

- 13:48:00 起手锚 `11cfbc93`；pin sha256 `7e23ec32…`（＝r6 交件值，未动过）；run.go `154279f7`／ticket224 `ebe020dd`／loop.go `0d53277d`。
- 执行窗内 HEAD 经 `7a42e168`→`c92f4dd5`→`f1676bfe`（编排者/174-r6 交件/其他笔，**零 Go 文件与本件重叠**）。
- 窗内脏件：`scripts/portable-tests.sh`／`winsec-tests.sh`／`wisp-cli-tests.sh`＝`317-r1` 在飞（收尾时已不在脏单；本腿未读未碰）。
  本腿自己的包面窗口＝干净（两发测试只带本件）。
- GOTMPDIR＝反斜杠形 `C:\tmp\wisp236r7gotmp`（全程）。

## §6 具名没做

- 未刷 `:10` 的 `:892`/`:1099`（理由见 §2 末）；未碰 loop.go 族（现读仍准、非本族）。
- 未改断言/锁/id 列表/措辞；未碰 run.go／ticket224 源文件（只读）；未跑 CI 复读；⛔ 零 push、未碰他腿脏件。
- 裁定列表的块范围写法未采用（★差异已具名于 §0；若改判范围写法须同时改判据①口径——本腿不自行为之）。
- 未复跑 `internal/tools` 整包与 `internal/agent` 基线（面＝单件数字 token；`-run Test236` 已覆盖本件族；整包读数在 v2 窗口）。

## §7 件清单（committed）与 raw（留盘未提交）

committed：`00-anchor-and-gates.md`（本件）／`logs/10-anchor-reprojection-merged.txt`／`logs/11-digit-only-ruler.txt`／`logs/12-gates-merged.txt`（4 枚 ≤150 行）。
raw（只建不删、未提交）：`logs/00-time.txt`、`01-anchor.txt`、`01-reprojection-raw.txt`、`02-same-text-check.txt`、
`03-pre-Test236.txt`、`04-pre-gofumpt-build.txt`、`05-ruler.txt`、`06-post-Test236-race.txt`、`07-post-gofumpt-build.txt`、
`08-post-truth-projection.txt`（均为上列 committed 的来源快照）；备份 `C:\tmp\wisp236r7gotmp\pin.pre-236r7.go`（仓外）。
