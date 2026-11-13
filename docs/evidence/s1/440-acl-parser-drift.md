# 440 — 两个独立 ACL 解析器的射程比对与判归（`parseACELines` vs `aclSIDs`）

- 腿：440-v1（写腿，2026-11-13）
- 判归：**甲（合并）**，落点 `internal/winsec/aclparse.go`（新非测试包文件，既有 `internal/winsec` 包内）
- 基线：开工时 `git status --porcelain internal/models internal/winsec` ＝ **空**（读数本腿现跑，HEAD `673da11df`）
- 结案对象：票 238 AC#4③（`docs/evidence/s1/238-data-root-seal.md` §AC#4③：三欠账中唯一零归属者）；本件交付后该格有了具名归属＝票 440（本票）

---

## AC#0 只读比对：射程差逐条

### 0.1 各自能答什么、不能答什么

| 问题 | `parseACELines`（models, `acl_sid_121_test.go:57`） | `aclSIDs`（winsec, `acl_windows_test.go:86`） |
|---|---|---|
| 这对象上有哪些 trustee（主方）？ | ✅ `[]aceLine{Trustee, Rights}` | ✅ `names []string`（逐行解析后即时翻译） |
| 这些 trustee 的身份（SID）是什么？ | ✅（要配 resolver：`aceLinesForSID`/`sidIdentity`） | ✅（内置 PowerShell 翻译，`sidOf`） |
| 某 SID 名下有几枚 ACE、flag 文本是什么？ | ✅（`aceLinesForSID` 返回整行） | ❌ 只回 `sids`，行与 flag 全丢 |
| 这枚 ACE 是**允许还是拒绝**？ | ✅ `Rights` 里有 `(DENY)`（本腿实测，见 0.2） | ❌ flag 在 `:(` 之后的载荷里，它取 `line[:i]` 只留 `:(` 之前 |
| 这枚 ACE 是不是**继承来的**？ | ✅ `Rights` 里有 `(I)` | ❌ 同上 |
| 这枚 ACE 是修改（M）还是只读（RX/F）？ | ✅ `Rights` 原文 | ❌ 同上 |
| 名字解析不了（icacls 打印裸 SID）时怎么办？ | ✅ 报错不静默丢（`TestAC49` 钉死） | ✅ `sidOf` 对 `S-1-` 前缀直通（`acl_windows_test.go:63`-65） |

一句话：`aclSIDs`＝**身份半边**（谁的授予）；`parseACELines`＝**身份＋类型半边**（谁的授予＋什么类型的授予）。票 238 结案腿那句「沿用现有助手就等于沿用它的盲区」的物理位置就是 `line[:i]` 这个切片——冒号之前的 trustee 留下，冒号之后的 rights 载荷（含 `(DENY)`、`(I)`）整段丢弃。

### 0.2 同一份真实 dump 的逐行对照（本腿实测，2026-11-13，临时目录 `%TEMP%\wisp440-6ff42f5d`）

**夹具**（命令原样）：
```
icacls <dir> /inheritance:r /grant:r *S-1-5-32-545:(OI)(CI)(M) *S-1-5-18:(OI)(CI)(F)
icacls <dir>\probe.txt /deny *S-1-5-32-545:(DC)
```

**icacls 对 probe.txt 的原文输出**（删去 trailer 前的 ACE 行）：
```
C:\...\probe.txt BUILTIN\Users:(DENY)(DC)
                 NT AUTHORITY\SYSTEM:(I)(F)
                 BUILTIN\Users:(I)(M)
Successfully processed 1 files; Failed processing 0 files
```

**逐行对照**（两个解析器各自的返回）：

| 行 | icacls 原文 | `parseACELines` 返回 | `aclSIDs` 返回（进 `sids`/`names` 各一元素） |
|---|---|---|---|
| 头行 | `<path> BUILTIN\Users:(DENY)(DC)` | `{Trustee:"BUILTIN\Users", Rights:"(DENY)(DC)"}` | `names[0]="BUILTIN\Users"`、`sids[0]="S-1-5-32-545"` |
| 行 2 | `NT AUTHORITY\SYSTEM:(I)(F)` | `{Trustee:"NT AUTHORITY\SYSTEM", Rights:"(I)(F)"}` | `names[1]="NT AUTHORITY\SYSTEM"`、`sids[1]="S-1-5-18"` |
| 行 3 | `BUILTIN\Users:(I)(M)` | `{Trustee:"BUILTIN\Users", Rights:"(I)(M)"}` | `names[2]="BUILTIN\Users"`、`sids[2]="S-1-5-32-545"` |
| trailer | `Successfully processed...` | 双方都跳过（`i <= 0` 判据同源） | 同左 |

读法：三行的 `names`/`sids` 逐元素相同；差别全在 `parseACELines` 每行多出的 `Rights`。第 1 行是 deny、第 2 行是继承、第 3 行是继承的 Modify——这三个判断**只有前者答得了**。写成封条判据的人如果拿到的是 `aclSIDs` 形状，这三列答案就是空的。

**顺手量的一个新事实（票 440 的问题库因此多一条）**：deny ACE 在 icacls 输出里的形状是 `BUILTIN\Users:(DENY)(DC)`——`(DENY)` 标记在 `:(` **之后**的载荷里。本仓现存的五份 `:(` 解析器全部把冒号后的载荷或丢或只做 `(I)` 单项 contains（`observe`/`agent` 两份），⇒ **全仓今天没有任何一份解析器把「允许还是拒绝」递到判据手里**，合并后的 `winsec.ParseACELines` 是第一份。

**附带观察**：对本机 `/grant "*S-1-5-21-3623186960-731165060-4091685855-1717338598:(R)"`（裸 SID 授予）icacls 直接拒绝（`No mapping between account names and security IDs was done`），与票 121 注释里引用的「新查目录上出现未解析 SID 行」互为印证——未解析 trustee 在真实 profile 里出现，但**不是**通过 /grant 裸 SID 这条路制造出来的。

### 0.3 特别验：结构体字段清单

**`aceLine`（合并前 `internal/models/acl_sid_121_test.go:41`-48，合并后 `internal/winsec/aclparse.go` 的 `ACELine`，字段一字未动）**：
1. `Trustee string` — icacls 打印的主方原文（名字或裸 SID，从不单独当身份用）
2. `Rights string` — `"(...)"` 载荷原文（如 `"(I)(OI)(CI)(M)"`、`"(DENY)(DC)"`）

**⇒ 票面要验的答案：有 flag 字段**，就是 `Rights`——不是独立命名的 flag 字段，而是载荷原文整段保留（`DENY`、`I`、`OI`、`CI`、`M`、`DC` 等全部标记都在里面）。

**`aclSIDs` 丢了什么**：整个 `Rights` 载荷。切面就是 `line[:i]`（`acl_windows_test.go` 合并前 ：95）——冒号前的 trustee 进 `names`/`sids`，冒号后的 `:(DENY)(DC)`/`(I)(F)`/`(I)(M)` 不进任何返回值。跟着丢掉的判断力：deny/allow、inherited/explicit、M/RX/F/R 等级——即 0.1 表的后三行。

### 0.4 比对中的额外发现：`:(` 解析器全仓五份，不是两份

票面只点了 models/winsec 两份。全仓 grep（`LastIndex\(line, ":\("\)`）现跑读数＝**5 处**：
1. `internal/models/acl_sid_121_test.go:62`（本票主角之一，→ 已合并）
2. `internal/winsec/acl_windows_test.go:91`（本票主角之一，→ 已改为调共享件）
3. `internal/observe/logprivate_132_windows_test.go:110`（票 132 判据，只取 `(I)` 单项）
4. `internal/config/private_acl_windows_test.go:166`（`principals`，只取名字）
5. `internal/agent/spill_acl_windows_test.go:51`（`assertNoForeignACE`，只取名字）

三份内联副本**不归本票收**（见 AC#5），但票面「两个互相独立」的计数今天修正为「两个命名解析器＋三份内联惯用法」。

---

## AC#1 判归：甲（合并），落点 `internal/winsec/aclparse.go`

**理由（按主流 harness 判据「最有价值」，非「动作小」）**：

1. **两者是不是该是同一个东西——量完的答案是：是。** 逐行对照（0.2）显示两解析器对同一 dump 的输出是**同一事实的两个前缀**：`aclSIDs` 的每行输出 ＝ `parseACELines` 对应行的 `Trustee` 列＋一次翻译。拆开它们存在的唯一理由是「当时的包边界」，不是射程设计——两个解析器连 `:(`/trailer/路径粘连三个陷阱的判据都逐字同源（winsec 原注释 ：82-85 与 models 原注释 ：52-56 互相引用了对方文件里记录的同一个 bug）。这不是「两个该分工的东西」，是**一个东西的两份漂移副本**。
2. **主流 harness 怎么做**：测试辅助里的重复解析逻辑，向上抽成被测包（或紧邻的小包）的导出面，测试文件只留包装——这是 Go 仓的常规做法（`testing/iotest`、`httptest` 之类「给测试用的共享件住生产包」的先例）。判据是「下一个写判据的人只有一个解析器可挑」，这个价值只有合并买得到；注释互指（乙）买不到——互指挡不住下一次复制粘贴（本仓自己就是证据：winsec 的 `:(` 副本与 models 的 `:(` 副本当初就是这么长出来的）。
3. **为什么不判丙**：238 证据件 §AC#4③ 现跑读数（本腿复核）——issues 树内 `parseACELines` 仅票 238 自身两处提及，零票号、零台账号。没有他票可归。
4. **为什么不判乙**：见理由 2。乙的产出（两段注释）不消灭「挑错就空」的根因，只给它立路牌。
5. **落点为什么是 `internal/winsec` 而不是新包 `internal/aclparse`**：不是「动作小」，是**两个硬机制**——
   a. **依赖边**：models → winsec 干净（winsec 对 models 的唯一命中是 `acl_windows_test.go:234` 的一条注释，grep 现跑核实，无循环）；反向（winsec → models）才成环。
   b. **CI 名册的钉子**：`scripts/portable-tests.sh` 的 core tier 用「静态 core_pin 名单 ↔ glob 解析集」双向核对（GUARD C，:780「resolved to a DIFFERENT package set」即红），GUARD D（:633-646）再全仓清点「有测试文件却不被任何 tier 认领」的包。新增子包当天就该进 core_pin＋scope，而 `scripts/**` 不在票 440 的写面（写面＝`internal/models/**` 与 `internal/winsec/**` 的 `_test.go`＋「判甲时可能新增一个非测试包」）——抽进 winsec 本体，两条 CI tier（core :387 与 winsec :408-409）对 `internal/winsec` 的认领**原样继续为真**，名册零移动。这正是「新测试面被现有仪器自动覆盖」与「不越写面」的交集。

---

## AC#2 合并后两边判据读数＋成对判据两发

**改动清单**（`git diff --name-only` 本腿四个文件）：
- 新增 `internal/winsec/aclparse.go`（生产文件；`ParseACELines`/`ACELine`/`IsNumericSID`/`SIDIdentity`/`ACELinesForSID`，逐字承自 models 副本，字段与判据零改动）
- `internal/winsec/acl_windows_test.go`：`aclSIDs` 改为共享件上的薄包装，签名与调用面不变（同包 33 处调用零改动）；注释写明射程（只答「谁的授予」）＋盲区（载荷丢弃 ⇒ 答不了 deny/继承）＋指回 `winsec.ParseACELines`
- `internal/models/acl_sid_121_test.go`：删本地定义（struct/两函数/一方法），`parseACELines`/`aceLinesForSID`/`isNumericSID` 改 `var x = winsec.X` 委托，`aceLine` 改 `type aceLine = winsec.ACELine` 别名——**文件内其余判据逐字未动**
- `internal/models/acl_sid_121_windows_test.go`：一行连锁（注释里方法名 `sidIdentity`→`SIDIdentity`）

**断言零删除的读数**（每文件 `func Test` 计数，HEAD vs 工作树）：
- `acl_sid_121_test.go`：3 → **4**（＋票 440 成对判据，未删任何旧判据）
- `acl_windows_test.go`：5 → **5**（零测试函数变动，判据全部保留）

**门禁读数（合并后，本腿现跑；数值取 AC#4 的终态一轮）**：
```
go test -count=1 ./internal/models/ ./internal/winsec/
  ok  github.com/CarlosShao/wisp/internal/models   10.642s
  ok  github.com/CarlosShao/wisp/internal/winsec  26.788s
```

**成对判据两发**：

第一发（钉进仓，`internal/models/acl_sid_121_test.go` 的 `TestAC440TheSharedParserKeepsTheRightsPayload`，便携腿，两 CI 腿都跑）：
- **合并前（丢弃副本的形状）红**：identity-only 行（SID＋名字，无载荷）对 `(DENY)(DC)`/`(I)(F)` 两个问题**答案不存在**——测试红半边把该形状钉死（「两问的匹配检查永远不火＝盲区本体」）。
- **合并后绿**：`parseACELines` 对同一 dump 交出 `{BUILTIN\Users, (DENY)(DC)}` 与 `{NT AUTHORITY\SYSTEM, (I)(F)}`，两个标记都在判据手里。
- **变异实证**（红半边不只靠形状推演）：把 `aclparse.go` 的 `Rights: line[i+1:]` 临时填 `""` 再跑 → `TestAC440` 红（`the deny marker did not survive the parse`）＋ `TestAC47` 四子句全红（`rights payload lost while parsing: ""`）→ 精确还原后复绿。即「丢弃载荷」这个旧副本的行为差，被现存判据当场抓住。

第二发（winsec 侧，同包现跑）：`aclSIDs` 包装化后全包 33 处调用判据不动，`go test ./internal/winsec/` 绿（含票 89/108/112/118/120/129 全部 ACL 判据）——身份半边判据在共享件上逐一仍真。

---

## AC#3 契约轴

未改 `docs/PLAN.md`、`docs/specs/**`、`thresholds.go`、golden、`allowlist.txt`。本腿 diff 仅 0.2 节所列四件 Go 文件＋本证据件＋票面 Progress log。复核命令：`git diff --name-only`（本腿文件清单里无契约件）。

## AC#4 门禁（三条，合并后终跑读数）

```
go test -count=1 ./internal/models/ ./internal/winsec/
  ok  github.com/CarlosShao/wisp/internal/models   10.642s
  ok  github.com/CarlosShao/wisp/internal/winsec  26.788s
gofumpt -l internal cmd tools scripts frontend
  （空输出，exit 0）
scripts/d22scan.sh
  d22scan: clean - no D22 ban violations  （exit 0；internal 生产面 403→404、含注释面 1004→1005，即 aclparse.go 已入扫）
```

基线（开工时同三条）同绿：models 10.972s / winsec 21.990s；gofumpt 空；d22scan clean。环境注：本机裸 `bash` 解析到 WSL（无 /bin/bash），d22scan 实跑用 `D:\work\soft\Git\bin\bash.exe scripts/d22scan.sh`——脚本未动，只是调用方找到了真 bash。

## AC#5 边界声明：本票收什么、不收什么

**收**：
- `parseACELines`/`aceLine`/`sidIdentity`/`aceLinesForSID`/`isNumericSID`（models 副本）与 `aclSIDs` 的解析内核（winsec 副本）合并为 `internal/winsec/aclparse.go` 一份；
- 两处测试包装的射程注释（含互相指回）；
- 票 238 AC#4③ 的归属落账（具名票 440）；
- 成对判据一枚（TestAC440）。

**不收**：
- **三份内联 `:(` 副本**（`internal/observe/logprivate_132_windows_test.go:110`、`internal/config/private_acl_windows_test.go:166`、`internal/agent/spill_acl_windows_test.go:51`）：它们是各票判据的内联惯用法（0.4 节），各有各的票（132 等）；动它们要过各自判据的验收，且 `internal/config` 有另一条腿在写。已在此登记，够下一位定位。
- **`aclSIDs` 现有调用面的语义**：包装保留 `(sids, names)` 签名与 `sidOf` 翻译行为，33 处调用零改动；把它们迁到 `ParseACELines` 直调不是本票的活。
- **deny/allow、继承语义的策略层**：共享件只交载荷原文；「deny 压过 allow」「integrity 标签不是授予」这类策略仍归各判据（`privateACLError` 等）。
- **CI 名册**：未动 `scripts/portable-tests.sh`（也不在其写面）；抽进 winsec 本体使名册无需移动（AC#1 理由 5b）。
- **契约件**：AC#3 清单零触碰。

## 撤销口令

`git revert <本票 commit>`——单 commit，反向即复原（解析器回到双副本、TestAC440 随之删除，各判据回到原定义点）。
