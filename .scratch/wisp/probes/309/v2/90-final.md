# `309-v2` 终局件（判语摘要／记账／欠账／⛔ 我做了哪些）

## 1. 三格判语（全文与逐枚尺＝`docs/evidence/s1/309-v2-verdicts.md`）

- **`AC#1`＝成立。** `309-v1` 那枚"只装了一肢"的条件按编排者 §4.4（⛔ 造第二枚仪器、改注释）已闭合，
  且注释那两句我实测复现：摘 `b.injected` 整包净新增红 **0**、摘 `b.prov == nil` 净新增 **16** 且⛔ 一枚名册讲这一肢、
  两发里 309 仪器自己照 `--- PASS`（`logs/rigA-instrument.log` A5／A6）。
  本格要求的四形（摘肢／探测器改没挂／整枚守卫 off／肢取反）我各跑一发**全红**，pristine 一发**全绿**，每形 `cmp`＋`sha256sum` 回到原样。
- **`AC#2`＝成立。** `309-v1` 三条条件逐条由**我这一发**复现：整包 `internal/agent` 不再被跳（B2 红）、
  声明＋赋值形与别名形在**盘上**红（E1/E2，⛔ 只在夹具里）、roots 里 `tools`（B3）与 `scripts`（B4）都 live。
  牙还在第二向：移除被点名那枚⇒红（B5，`名册=(空名册)`）、`wantTotal` 放宽一档⇒当场红（B6）。
  新抓四处残余⛔ 旧条件（见 §3），其中两处实测**绿**（B7/B8/B14/E3）＝本格⛔ 是"全覆盖"，只到三形。
- **`AC#3`＝附条件成立，一枚条件。** 注释里那把自指尺的每个数（71／69／2／1，简写形 2）我这发逐字重跑对得上；
  五枚内容锚在盘上各 1 命中；点名段内⛔ 绝对 `loop.go:<数字>`；契约声称文字增删 0 行；"行为逐字节不变"
  ＝改动行 70 枚全为注释、非注释 0、两版剥注释 `cmp` IDENTICAL（854⇄854）。
  ⚠ **条件**＝`bridge.go:959-963` 那句全称"counts **every** production construction point … so a second one rings **wherever** it is written"
  被 **B14** 判死（`new(agent.ToolRequest)` 那枚真第二构造点令钉**绿**），而同文件钉头 `:52-63` 自己列了七条漏形＝两句互相矛盾。
  改法便宜到不该占一格：把那半句换成"rings for the three shapes the matcher names; the blind list is in the test's head comment"。

## 2. 记账尺（§5 要求那两个数并排）

- 起手：`- [x]`＝**3**／`- [ ]`＝**3**；收尾（本腿⛔ 翻框）：`- [x]`＝**3**／`- [ ]`＝**3**（尺＝`grep -cE '^[[:space:]]*- \[x\]'`／`- \[ \]`，含缩进项）。
- 票尾本腿只追加 **1 条** Progress log bullet，`git diff --numstat`＝**＋1／−0**（原句零改动）。
- 四道门禁（起手态）：`rc_build=0`／`rc_vet=0`／`rc_d22scan=0`／`rc_pathlen=0`，件各带 `rc=` 行且⛔ 0 字节
  （`logs/gate-{build,vet,d22scan,pathlen}.log`；build／vet 两枚只有 11／9 字节＝那行 `rc=`，输出本就 0 字节＝无诊断）。
- 写完自家件之后**再跑一遍**那两把会被写面影响的尺：`logs/gate2-d22scan.log`＝`rc_d22scan_with_my_files=0`
  （末段含 live 门禁那句 `clean - no D22 ban violations`，`grep -c`＝2）、
  `logs/gate2-pathlen.log`＝`rc_pathlen_with_my_files=0`（`tracked paths=10697  over-budget=57  covered=57  not in roster=0  VERDICT GREEN`）。
- 突变台件⛔ 落进母仓工作树：全部打在两棵 `/tmp` 仓外导出树（`git archive f6067e5d`，单面、⛔ `cp` 别锚）；
  全程 `git status --porcelain -- internal cmd`＝**0 行**（A/B/D 三组末尾各数一遍）。
- 起手 `git -c core.quotePath=false diff --cached --name-only | wc -l`＝**0**；我没清空过任何东西、没看到大面积 staged 别家路径。

## 3. 具名欠账（`Q-309-v2-*`，⛔ 我给任何格盖章翻勾）

1. `Q-309-v2-1-assembly-limb-blind`＝`b.injected`／`b.prov == nil` 两肢⛔ 具名仪器（编排者 §4.4 裁⛔ 造第二枚；注释那半句我实测如实）。
   残余形状＝"装配坏了"那一族；台件两枚写法在裁决表 ②(iii) 里给了（全用现成夹具，零新导出 API、零新字段）。
2. `Q-309-v2-2-roots-and-floor-unpinned`＝**B8**：roots 被一行 commit 收成 `cmd internal`、同时 `tools/` 里真有一枚构造点⇒**绿**（净少扫 24 枚源）；
   **B7**：分母 floor `200→1`⇒**绿**。⇒ 钉的"射程"本身⛔ 被钉着。修法＝把总分母（现量 293）或每根下限当第二枚钉。
3. `Q-309-v2-3-blind-shape-rings-green`＝**B14**（`new(T)`）＋**E3**（省类型嵌套）两条漏形我补上实测读数；
   钉头其余三条（函数体内别名／改名 import 与 dot-import／定义型＋转换）仍**具名⛔ 实测**＝`Q-309-r2-nail-blind-shapes-unrigged` 未结完。
   本格同时判死 `bridge.go:962` 的 "wherever"（＝`AC#3` 那枚条件）。
4. `Q-309-v2-4-frontend-outside-the-roots`＝`frontend/embed.go` 是 `go list ./...`（35 枚包）里唯一落在那四根之外的模块包（今日 0 命中）＝要不要进 roots 归编排者裁。
5. `Q-309-v2-5-callsite-covered-by-other-tickets`＝**D1/D2**：调用点解耦由票 124/175/183/185 那族测的牙接住（净新增 7／3 枚红），
   309 自家仪器两发全绿 ⇒ `bridge.go:982` 那句"Who rings…"该配一行" rings，但⛔ 由本票那枚仪器"。
6. `Q-309-v2-6-legacy-go-snapshots-under-probes`（⛔ 本票三格账，只登记）＝`.scratch/**` 下 tracked `.go` **314** 枚
   （含 `probes/309/a1/logs/bridge.go`、`loop.go` 等），与"⛔ 在 `probes/**` 下建 `.go`"相反；本轮四把门禁零影响（实测）。

## 4. 交件动作（派单 §6 两条硬规矩）

- `git add -- <显式路径>`：⛔ `git add -A`／`.`／裸 add。
- `git commit -F .scratch/wisp/probes/309/v2/logs/commit-msg.txt -- <显式路径>`：commit 也带 pathspec（今天下午那起 1759 枚别家路径被扫进共享索引的事件就是这条规矩的出处）。
- ⛔ push（`A813` 一次未用）／⛔ `--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`／`--no-verify`／⛔ 翻框／⛔ 动 `Status:`／⛔ 改名 `-done`。

## 5. ⛔ 我做了／⛔ 我复跑（分层）

- ⛔ 跑整仓测试名册／格式名册／eol 三把（`AC#4` 已由编排者翻勾，本腿⛔ 重判）；⛔ 重跑 `309-v1` 那两枚 236 载体的三张台面。
- A6 那 17 枚红的**成因**我只逐字读了 2 枚（其余按名册＋"⛔ 一枚名带 prov/assessorfor/taint/309"这把尺推）＝
  "全是管道噪声"这句是**读码推**，⛔ 16/16 逐枚量过。
- R359-3／R359-5（`filterdemo`、`bridge.go` 枚数逐枚对账 7→10、2→1）⛔ 重跑，只引用腿的读数。
- 工作树⇄HEAD blob 两面 interchangeable 只在 `internal cmd tools scripts` 那四档成立（`git status --porcelain`＝0 行实测）；
  别处（`.gitignore`／`go.mod`／`design/**` 有别腿在飞态）我⛔ 用作任何尺的台面、⛔ 碰过一枚。
