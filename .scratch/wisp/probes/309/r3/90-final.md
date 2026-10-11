# `309-r3` 收尾一页（写腿，⛔ 裁决者）

判语（`AC#2`/`AC#3` 成不成立）⛔ 本腿的活 ⇒ 归 `309-v3`（非实现者）。这里只交"改了什么＋每一发的原始读数＋还欠哪几枚"。

## 1. 两笔代码＋三笔证据（commit 号）

| 笔 | 内容 | 名枚 |
|---|---|---|
| `f4920906` | 起手锚（只 `00-anchor.md`） | 1 枚件 |
| `df8d517d` | (A) 补射程：`internal/tools/toolrequest_nail_309_test.go` 一枚件，`120 加／13 删` | 1 枚文件 |
| `a1f84d39` | (B) 改正文案：`internal/tools/bridge.go` **注释行**，`14 加／4 删` | 1 枚文件 |
| `6fed247e` | 证据件＋`logs/` | 58 枚件，`+10867 −0` |

起手锚＝HEAD `3974ab4d`、`dev`、`internal cmd tools scripts` 工作树 **0 行**、全索引暂存 **0 枚**（只报数⛔ 归因；本腿一枚别家路径⛔ 动过，⛔ `reset`／`stash`／`checkout .`／`clean`）。

## 2. 读数（每把尺带"哪把尺＋射程＋量的还是算的"）

| 承重读数 | 尺 | 射程/面 | 读数 |
|---|---|---|---|
| 三形各有牙 | `logs/rig-shapes.txt`（仓外导出树，`git archive a1f84d39`） | 一枚新 `.go` 加进树 ⇒ 跑指名钉 | 六发＝**红红绿绿红绿**：`new`⇒`rc=1`＋名册带 `new-of-tool-request`／拿掉`rc=0`；`make`⇒`rc=1`＋`make-slice-of-tool-request`（落 `internal/agent`，裸名＋长度 `0` 一形两验）／拿掉`rc=0`；转换⇒`rc=1`＋`conversion-to-tool-request`／拿掉`rc=0`；正控原树`rc=0` |
| 变异⛔ 落进共享工作树 | `git status --porcelain -- internal cmd tools scripts`（脚本末节） | 四根 | **0 行**（empty above＝0 枚）；三发还原各自 `cmp … IDENTICAL` ＋ `sha256[:16] tree=c4aa3d46a6daf3d8/commit=同` 与 `e6c685da6bfbcfdf/commit=同` |
| 整包名册改前⇄改后（同面正控） | `logs/roster-control.txt`，两棵树都 `git archive`（`3974ab4d` ⇄ `a1f84d39`）；名册尺先剥耗时列 | `./internal/tools/` 整包 | 名册行 **279⇄279**、顶层 **211⇄211**、子测试 **68⇄68**、FAIL **1⇄1**、SKIP **0⇄0**、`newly_red_count=0`；added/removed 各 1 枚＝同一枚子测试名里嵌了 `mktemp` 目录名（具名报，⛔ 洗成 0/0）；唯一红＝`TestTicket252R1AlignmentAddsNoAuthorization`（改前那发就红＝导出树缺件族，票 298 `AC#5` 已承接，⛔ 本票⛔ 本腿） |
| (B) 只改注释 | 三把：`git diff -U0` 非 `//` 增删行计数；两版 blob `sed '/^[[:space:]]*\/\//d'` 后 `cmp`；关键语义行 `grep -c` | `internal/tools/bridge.go` | `noncomment_changed=0`、`cmp IDENTICAL`（**854⇄854**，`sha256=496ed6d683e8e070…` ⇄ `309-v2` 同一条）、守卫行与 `WithTaintDetector` 各 **1⇄1**、CR 字节 **0** |
| 射程枚数 | `grep -oE '…' \| sort -u \| wc -l`（七枚 shape 标签，见 `10` §3） | 钉那件 | **7**（⛔ 含 `literalShape309` 那个结构上到不了的 `default` 分支） |
| 自指词面尺（`AC#3` 那半句） | `grep -rn 'ToolRequest{' internal cmd --include=*.go` 双滤口径（`logs/selfref-ruler-final.log`） | `internal`＋`cmd`，工作树面（＝HEAD 面，产码笔已提交） | 总 **72**（v2 那次 71 ⇒ ＋1＝本腿在 `_test.go` 里新增的一枚控制串 `[]agent.ToolRequest{{TaskID: ""}}`）／按**文件名字段**非测试 **2**／简写 `grep -v _test.go` **2**／剥行首 `//` 后 **1**（＝生产构造点 1 枚，结论⛔ 变） |
| 三枚新形今天有没有人真踩 | `grep -rnE 'new\(agent\.ToolRequest' … \| grep -v _test.go \| grep -v '^…//' \| wc -l`（注释行⛔ 数进分子） | `cmd internal tools scripts` | 三枚都 **0**（empty above＝0 枚）⇒ 生产构造点仍＝1 枚，钉绿是"真没人写"而非"尺瞎" |
| 四道门禁 | `GOFLAGS= go build ./...`／`go vet ./internal/tools/`／`sh scripts/d22scan.sh`／`bash scripts/check-path-length-budget.sh --with-self-test` | 各自件末行自落 `rc=` | `rc_build=0`／`rc_vet=0`／`rc_d22scan=0`（live 那句 `clean - no D22 ban violations` 逐字在件里）／`rc_pathlen=0`（`VERDICT GREEN`，`tracked paths=10828`、`over-budget=57 ⇄ roster=57`、`not in roster=0`）；写完自家件后第二轮 `rc_d22scan_with_my_files=0`／`rc_pathlen_with_my_files=0` |
| 格式名册并排 | `gofmt -l internal/tools` ⇄ `gofumpt -l internal/tools`，工作树面 ⇄ HEAD blob 面（逐枚 `git show` 导仓外量） | `internal/tools`（78 枚 tracked `.go`） | **0⇄0**（两面、两把都 0＝empty above＝0 枚）；整仓旁证 `gofmt -l .`＝37 枚，起手就在册、⛔ 归本腿 |
| ⛔ 放宽断言 | `git diff 3974ab4d a1f84d39 --` 钉那件 `grep -E '^[+-][^+-].*(wantTotal\|wantOwner\|examined < 200)' \| grep -c` | 钉那件 | `assertion_lines_changed=0`（empty above＝0 枚）；控制名册 6 → **18** 枚（`grep -c` 数名册），⛔ `t.Run` ⇒ 名册零新增 |

## 3. ★还欠的读数（逐枚具名 `Q-309-r3-*`，⛔ 由本腿盖章任何格）

1. **`Q-309-r3-reflection-decode-not-statically-coverable`**（本票真实残余，⛔ 由任何一程宣布解决）：反射造值／`json.Unmarshal` 解码进结构体这两形**没有静态尺抓得住**——`AC#0` 已裁形 Ⓑ（入口⛔ 拒），所以这两形今天仍只有注释登记。本腿只把它们**留在名册里逐枚具名**（钉 `:68-70`、`bridge.go:966-972`），⛔ 声称覆盖。要不要另立一枚运行时闸门归编排者摆。
2. **`Q-309-r3-conversion-branch-is-syntactic`**（本腿自抓、新登记）：转换那一支是纯句法判定（`Fun` 解析得到 `agent.ToolRequest` ∧ 实参一枚）⇒ 一枚恰好同名同包、只吃一枚实参的**函数**调用会被计进分子（假阳方向，⛔ 假阴）。要区分必须引 `go/types`，与件 `:80` 那句"⛔ 类型级完美"同一口径；本腿⛔ 造反证台件去量这枚假阳有多易触发＝**没量的读数**。
3. **`Q-309-r3-three-blind-shapes-still-unrigged-on-disk`**：名册里"函数体内声明的别名／改名 import／dot-import"三形**至今没有一发在盘上的绿读数**（`309-v2` 的 `Q-309-v2-3` 那半条未结完；另两形各有 r2 的一发绿＝RIG-5 省类型嵌套、RIG-6 跨包别名）。⇒ 全称判据的残余枚数：本腿把 9 形削到 **7 形仍漏**，⛔ 归零。
4. **`Q-309-r3-roots-and-denominator-unpinned`**（照收 `Q-309-v2-2`，本腿⛔ 动 roots/地板）：roots 被一行 commit 收成 `cmd internal`、或分母地板 `200` 被改小到仍高于实际枚数，都⛔ 响；本腿只在 `10`/`20` 里报了自己看到的面（分母 **295**，比 `309-r2` 的 293 多 2＝别腿 `internal/llm/adaptertest/**` 已提交所致）。
5. **`Q-309-r3-roster-interval-not-pure`**（口径欠账，⛔ 判语）：第 2 行那把"改前⇄改后"名册尺打在 `3974ab4d..a1f84d39`，区间里除本腿两笔还有 `2a6f11a5`（编排者）与 `304-r1` 两笔（`b9f37892`、`e875f00f`，写面 `internal/llm/**`）⇒ 它证明的是"那一段区间⛔ 让 `internal/tools` 多红"，⛔ 单独证明本腿因果（本腿因果＝那六发＋`go build`/`go vet`）。
6. **`Q-309-r3-assembly-limbs-still-blind`**（照收 `Q-309-v2-1`，编排者 §4.4 已裁⛔ 造第二枚仪器）：`b.prov == nil`／`b.injected` 两肢仍⛔ 具名仪器，注释那半句本腿一字⛔ 改（只改了 `:959-972` 那一段射程句）。

## 4. 指回件名

`00-anchor.md`（起手三把尺）· `10-matcher-and-comment.md`（逐枚 文件:行＋七枚 shape 标签＋被拒形逐枚引原文行号）· `20-teeth.md`（六发红绿＋还原两行＋整包名册）· `30-gates.md`（四道 `rc=`＋格式名册两面＋越界名册）· 台件与原始读数＝`logs/`（`rig-shapes.txt`/`.log`、`rig-*.log` 各带 `rc=`、`roster-control.txt`/`.log`、`gate-*.log`、`selfref-ruler-final.log`、`bridge-comment-scope*.log`）。

⛔ 翻任何框（现量＝`- [x]` **4**／`- [ ]` **2**，尺＝`grep -cE '^[[:space:]]*- \[x\]'` ⇄ `- [ ]` 现跑于收尾那发；⚠ 这⛔ 是 `309-v2` 报的 `3／3`——差的那一枚＝编排者 17:1x 那节自己翻的 `AC#1`，本腿⛔ 动过任何框）· ⛔ 碰 `Status:` · ⛔ 改名 `-done` · ⛔ push。
