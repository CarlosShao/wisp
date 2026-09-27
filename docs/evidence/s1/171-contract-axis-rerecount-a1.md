# 171-a1（只读·非实现者）复算表：票 161 AC#5「契约轴零字节」独立重走

- 程＝票 171 **AC#4 一格**（票 171 面 `:29` 现量行号）；被复核格＝票 161 面 `:29` 的 AC#5。
- 时刻 09-27 11:57–12:0x｜锚＝本程 step-0 现量 `d2e84f39ccd1be7a21859d0e408e5febac1ff5c9`（未抄任何号）。
- 台件与原始读数：`.scratch/wisp/probes/171/a1/**`（`commands.sh` 是可复制命令序列）。
- 总裁一句话：**「0 命中」与「main.go +12/-0」两数一致成立；「37 枚名册」这一数在票面支持的 5 种口径下均无法复算，按无法复算报回，其实质结论不受影响。**

## 1. step-0 四件

- `date` -> 2026-09-27 11:57:44 +0800
- `git rev-parse --abbrev-ref HEAD` -> **`dev`**（全程未切分支）
- `git rev-parse HEAD` -> `d2e84f39ccd1be7a21859d0e408e5febac1ff5c9`
- `git status --porcelain -- tools/d22scan/ docs/PLAN.md docs/specs/ internal/ cmd/` -> **空**
- 派单前提核真：票 161 文件名以 `-done` 结尾（成立，改名发生在 `e0405d57` 10:29）；`docs/evidence/s1/` 下 `171-*` 枚数为 0（`ls | grep -E '^171'` 空）＋票 161 全部 8 枚证据件与票面 Progress log 中不存在任何 AC#5 的独立名册表（v1 表裁 AC#1/2/3、v2 表裁 AC#4/6/7）⇒ **「这格至今缺独立表」成立，本件是第一份**——与 171 面 `:29` "没有独立表"一致。

## 2. 本程没测什么

- 没裁票 161 其它任何一格，没碰票 171 的 AC#1/2/3/5/6；没读 `probes/154/**`（射程外，一眼未看）。
- 没跑 `go test ./internal/...`（派单明令：`160-v1` 同树在跑，避免互采半成品）。
- 没验 `ci.yml` 两步的内容合法性（那是 AC#2 已裁格）；没重做摘样本造响（AC#4 格，v2 已裁）。
- 命中口径＝**提交所改文件清单**（`--name-only`），不是 blob 内容 diff——对"零字节"这种"这些提交对这批路径动没动"的命题，文件级即其判据；未对 268 枚展开文件逐枚 `git show <sha>:<path>` 比 blob（无必要：名册里没有任何一枚带这些路径）。
- 编排者"37 枚"的原始命令我没有，也不该猜——见 §4 的无法复算档。

## 3. 名册怎么来的（可复制命令＋枚数）＋逐枚表

号源＝票 161 面 Progress log＋`docs/evidence/s1/161-*.md` 语料自登记的号，**每枚先 `git cat-file -t` 复核**：

- 语料 2842 行抽 7–10 位 hex token -> 58 枚；其中 **41 行解成 commit、去重（同 commit 不同长度）= 34 枚**；**17 枚取数前即拒**（`123a124` 是 diff hunk、`20260926` 是日期、其余 15 枚是日志数值噪声，`cat-file -t` 无类型）。
- 34 枚里 6 枚是语料点名到的**他票号**：`15ff2be`(票77) `5afa666`/`fbe12c7`(前端log系，后者＝r1 锚) `611ae8b`(票169) `506cbae`(票158) `4267bb3c`(票164) ⇒ **161 系自登记 28 枚**。
- `git log dev` 按主题补登记（写了 161 却没在语料自登记号）另得 **13 枚**：`0492cd5d cd999794 8995b366 feed0a17 22ecf9db 71cff1cc 0d176473 70a106e5 339b67a1 bc65d975 a7f2d790 1a95eaa6 e0405d57`（含 r1/r2 余下提交、r4 补格、r5 代落、v2 三格、结案改名本身）。
- ⇒ 本程按**最宽口径 47 枚**（28＋6＋13）逐枚量。最宽成立则任何窄子集成立；他票 6 枚单列归属。
- 名册命令一律带 **`--no-walk`**（`git log --no-walk --reverse --name-only --format='@C %h %s' <47枚>` -> `roster-files-raw.txt`，463 行）。

逐枚表（号｜该枚所改文件枚数｜冻结 11 项命中〔文件〕｜r1 特别名单命中）：

```
506cbae files=50 axis=0 special=0            [票158]
15ff2be files=7  axis=0 special=0            [票77]
5afa666 files=1  axis=0 special=0            [前端log]
611ae8b files=1  axis=1[frontend/.../right-rail.tsx] special=0  [票169]
fbe12c7 files=1  axis=0 special=0            [r1锚/前端log]
16364e6 files=45 axis=0 special=0            r1
e809235 files=1  axis=0 special=0            r1
0492cd5 files=5  axis=0 special=0            r1
7e8cfdc files=3  axis=0 special=0            r1
8995b36 files=28 axis=0 special=0            r2
cd99979 files=1  axis=0 special=0            r2
feed0a1 files=1  axis=0 special=0            r2
27f798e files=11 axis=0 special=0            编排者(票面AC#6/7+派单存档)
0d17647 files=1  axis=0 special=0            编排者(停车点)
e3f5980 files=46 axis=0 special=0            r3(AC#7#1 止血)
3e56e68 files=19 axis=0 special=1[.github/workflows/ci.yml]    r3代提
70a106e files=1  axis=0 special=0            派单存档
339b67a files=1  axis=0 special=0            停车点
77e9ed2 files=30 axis=0 special=0            r4
22ecf9d files=2  axis=0 special=0            r4
ffbf463 files=1  axis=0 special=0            r4
d49a5be files=2  axis=0 special=0            r4
71cff1c files=3  axis=0 special=0            编排者(ledger A319)
83e2320 files=7  axis=0 special=0            r5
e67a577 files=1  axis=0 special=1[.github/workflows/ci.yml]    r5
c664478 files=3  axis=0 special=0            r5
79394ae files=10 axis=0 special=0            r5
ea415e7 files=1  axis=0 special=0            r5
413ea90 files=2  axis=0 special=0            编排者(代落r5顶回)
8efece3 files=3  axis=0 special=0            编排者(AC#6自纠+派单)
6b0a9f0 files=1  axis=0 special=0            ledger(A321)
aa2c0f4 files=1  axis=0 special=0            ledger(A321自纠)
f17b62f files=30 axis=0 special=0            v1
359bb87 files=1  axis=0 special=0            v1
1532794 files=1  axis=0 special=0            v1
4267bb3 files=2  axis=2[docs/PLAN.md;docs/specs/SPEC-07...] special=0  [票164]
bd7475f files=1  axis=0 special=0            v1
addb760 files=1  axis=0 special=0            v1
5d11b8f files=2  axis=0 special=0            r6
9f53057 files=21 axis=0 special=0            r6
b244730 files=2  axis=0 special=0            编排者(三格翻勾)
b067db0 files=1  axis=0 special=0            r6
30121fe files=3  axis=0 special=0            编排者(代落r6)
bc65d97 files=5  axis=0 special=0            v2
a7f2d79 files=3  axis=0 special=0            v2
1a95eaa files=1  axis=0 special=0            v2
e0405d5 files=5  axis=0 special=0            结案改名+立171
```

读数：**41 枚 161 系提交全部 axis=0**；47 枚并集里仅有的 2 枚命中都是**他票提交**（票 169 的 frontend 修复、票 164 的 owner 批准契约文本），与票 161 无涉。特别名单 3 项里 `ci.yml` 被 `3e56e68`/`e67a577` 动过 2 枚——**这不违反 AC#5**：161 面 `:22-24` 白纸黑字写明 ci.yml 不在 11 项名单内、AC#2 只许加步骤；`internal/observe/` 与 `cmd/wisp/` 在全部 47 枚上 0 命中（r1 当时多报的那本特别账也成立）。

## 4. 契约轴命中 vs 编排者三数（并列摆出）

冻结名单：161 面 `:29` 逐字 11 项指法（`docs/PLAN.md`｜`docs/specs/**`｜`internal/risk/**`｜`internal/panel/**`｜`internal/agent/approval/**`｜`thresholds.go`｜golden｜`allowlist.txt`｜`scripts/slo-check.ps1`｜`frontend/**`｜`design/**`）——**项数 11 与编排者口径一致**。`golden`/`thresholds.go` 本程现量展开（命令见 `commands.sh`）：`thresholds.go`＝1 枚（`internal/observe/thresholds.go`）；`golden` 忽略大小写＝**60 枚**（agent 11 sse＋`loop_golden_test.go`＋llm 41 sse＋`golden/` 包 3＋`harness_golden_test.go`＋mockllm 3＋`156-accept` 2）；其余各项＝1/14/37/20/18/1/1/85/30，**11 项合计展开 268 枚具体文件**；`fbe12c7` 版与 HEAD 版的通配展开清单逐行一致（同跑两发对过）。

| 编排者的数（171 面 `:29` 转述） | 本程独立读数 | 档位 |
|---|---|---|
| 名册 **37 枚** | 票面支持的口径全数试遍：自登记 161 系 **28**｜自登记全部 **34**｜含 git log 补登记的 161 系全集 **41**｜本程最宽并集 **47**；区间口径另试 4 发（`fbe12c7..1a95eaa6`=57、`..e0405d57`=59、`27f798e..e0405d57`=51、`..30121fe`=51，且区间混入票 164 的 4 枚 PLAN.md 提交 ⇒ 区间口径下"0 命中"根本不成立，可排除） | **无法复算**——37 取不出任何一枚能同时保住"集合口径＋0 命中"的定义；缺的不是命令是可复算性，**报请编排者交出当时的原始命令**；本表不判他错，判"该数不可复核" |
| 11 项名单 **0 命中** | **0 命中**，且对 47 枚最宽并集成立（唯一 2 枚命中确属他票）⇒ 对任何 ⊆47 的 161 系名册定义恒成立 | **一致** |
| `tools/d22scan/main.go` **+12/-0** | 47 枚并集里 main.go 只被 **1 枚**提交动过：`3e56e680`，numstat 唯一一行 `12	0`（`cat -A` 核过 `12^I0^Itools/d22scan/main.go$`，删除列是真 0 不是空行折叠） | **一致**（出处补认：`3e56e68`＝r3 撞顶代提，与 161 面 23:5x"12 行纯新增、0 行删除"自述吻合） |

坑①按派单预告核过：`tools/d22scan/**` 与 `cmd/**` 本就不在 11 项里，`+12/-0` 是"动了扫描器"的证据而非违例——本表两事分开记。

## 5. 门禁两把读数（只跑不写盘的那两把）

- 12:04:23 `sh scripts/d22scan.sh` -> **rc=0**，末行 `clean - no D22 ban violations`（live scope 计数原样在 `runtests-171a1.log` 同分钟现场输出，0 finding）。
- 12:04:51 `bash tools/d22scan/runtests.sh -C tools/d22scan ./...` -> **rc=0**、**PASS=34 FAIL=0 SKIP=0**、`=== RUN`=76、`[no tests to run]`=0（基线＝本程自量；与 161 面 r5/v1 登记的 34/76 同形）。
- 名册两向 comm：对 `probes/161/r4-pre-runtests.log.roster.txt`（34 行）按**测试名**比 -> `roster-comm-nameonly.txt`＝**空**（34 vs 34 同名）。第一发整行比对差出 68 行系本程格式伪影（耗时后缀+尾空格），已留 `roster-comm-diff.txt` 为证并以 name-only 为准——防的正是"一条假差集吞掉真读数"。
- 没跑 `go test ./internal/...`（派单 §3 明令）。

## 6. 被拒／没成功的调用

- 取数**前**：17 枚假 sha token 全部 `cat-file -t` 拒掉，无一进名册。
- 取数**后**两发自身失误、均已修正且不留疑：① `git log --format='x'` fatal invalid --pretty（求和命令自伤；和值改由全表 numstat 直接读出——main.go 在 47 枚上仅一行 numstat，任何含 `3e56e68` 的子集和都是 +12/-0，不依赖那发失败命令）；② 第一发 comm 因行尾耗时格式全差（§5 已述）。
- 无工具权限拒绝、无被拦调用。

## 7. 删除命令

`rm`／`rmdir`／`del`／`git rm`／`clean`／`--amend`／`reset`／`rebase`／`stash`／`checkout .`／`switch` 全程 **0 次**。台件只建不删（`probes/171/a1/**` 8 枚文件）。在己方台件上用过一发 `sed -i`（剥耗时后缀），只动自己本程所建文件。

## 8. 伪授权两栏

- 看到过什么：本程只读了票面/证据/台账/git 对象；`injection-timeline` 所列诸形一概未遇——无文件冒充用户话、无面板侧"允许"、无工具输出里的新指令。
- 做过什么：未把任何转述当授权；派单里编排者的三数全程按"待复核断言"处理（§4 即按其原样三档裁）；未勾任何票面框（161/171 一字未动）。

## 9. 凭据值

本程未接触、未转写任何凭据值；表内出现的字符串只有路径、提交号、测试名与 rc。

## 10. next=（归编排者，本程到此为止）

1. **37 枚名册**：请给出结案当时的原始命令；给出后本表可补第 6 种口径复算（预期：0 命中与 +12/-0 不动，仅名册数收敛）。若给不出，建议结案件第二见证直接引本表 41 枚全集口径。
2. 本件可作票 161 AC#5 的独立第二见证挂票 171 结案链（勾与 Status 归编排者）。
3. 台账若要记：`A##` 只追加，指"161-AC#5 独立复算完成；37 枚数不可复算待原始命令"。
4. `probes/171/a1/**` 全部留件；`160-v1` 在飞，本件对其无写面影响。
