# 95 — 票 89 只封了四类数据，**同族的另外七处 `0o600` 落点还开着**：先裁"日志/模型缓存算不算私有数据"，再一行级接线

**Status:** open（2026-09-21 15:3x 编排者建，来源=票 89 交件时**自己点名**的范围外残留；我没塞回票 89，那是扩界）
**Type:** 安全一致性（同一台机器上的其他账户能读到什么）+ 一条**需要裁定的定义问题**
**Blocks:** nothing · **Blocked by:** 票 **89** 验收（`acceptor-ticket89` 正在跑，判它的替代判据够不够）
              ⇒ 本票**等它的结论**再开工：如果 89 的"目录链先封"被证明有洞，本票接的是错的形状。
**Packages:** `internal/models/`、`internal/config/`、`internal/secret/`、`internal/observe/`、`internal/ball/`、
              `cmd/wisp/`（**⚠ 后两处是禁改区，见下**）。**复用** `internal/winsec/`（票 89/94 的地界，
              **不要修改它的语义**，只调用）。**禁改**：`docs/PLAN.md`、`docs/specs/*.md`、`internal/risk/**`（冻结）、
              `tools/d22scan/**` 与 `allowlist.txt`（**5 行非注释，只许变短**）。

## 现场（票 89 报的，逐条 file:line 你要自己核一遍再动手——本仓已抓到过两起引用腐坏）

票 89 封了四类：**artifacts / DPAPI blob / SQLite `db`+`-wal`+`-shm` / staging 临时文件**。
剩下的同族落点（今天仍然只写着 `0o600`/`0o755`，而 Windows 上那串数字**不落地**）：

| 落点 | file:line | 我的初步判断 |
|---|---|---|
| 模型下载 staging | `internal/models/downloader.go:224`、`:287`、`:566` | **要封**（票 89 的 AC#1 实测过同类 staging 带外来 `(M,DC)`） |
| 解包出来的模型文件 | `internal/models/archive.go:41`、`:73`、`:78` | **要封**（同上，且这些文件之后被当作可信模型读回） |
| 配置迁移写的文件 | `internal/config/migrate.go:83`、`internal/config/parse.go:213` | **待定**：配置里可能有 provider 端点与 ref，但不是 key 本体 |
| 凭据迁移 | `internal/secret/migrate.go:154`、`:168` | **要封**（碰的是凭据存储目录） |
| 日志 | `internal/observe/logging.go:72`、`:244`（**0644**） | ⚠ **这就是要裁的那一条**，见下 |
| 球位置 | `internal/ball/position.go:73`、`:77` | **待定**：读它只泄露"球在哪儿" |
| doctor / SLO | `cmd/wisp/doctor.go:248`、`cmd/wisp/slo_windows.go:559` | ⚠ `cmd/wisp/` 与 `internal/ball/` **现在是别人的地界**（票 64/65/68 的球、票 77 的面板宿主）⇒ **先登记、后接线**，别抢 |

## 为什么这张票不是"把 7 处换成一调用"那么机械

**"日志算不算私有数据"是一个会改变产品行为的决定**，不是风格决定：
日志里会出现路径、工具名、错误原文（可能含用户数据的路径），封严了 ⇒ **用户自己用别的工具 tail 不到**、
`wisp doctor` 这类诊断可能读不到、票 66 刚修好的 SLO 观测链可能瞎。
⇒ **AC#1 先出裁定**（每一类：封 / 不封 / 封但保留一条**显式**的诊断通路），**裁定没做完不许开始接线**。

## AC（1:1，裁决表 `docs/evidence/s1/95-*.md` 由验收方出）

- [ ] **AC#1** 上面七类**逐类裁定**"算不算私有数据"，每类给：里面**实际会出现什么**（去读一处真实写入的格式，
      给 file:line）、泄露给**同一台机器上另一个账户**的后果、以及封了之后**谁会读不到**（诊断通路代价）。
      ⚠ 不许整族打包贴同一个结论——那正是票 82 禁过的"整族打包"。
- [ ] **AC#2** 对判"要封"的每一处：**接线 + 一条真实读回的 `icacls` 证据**（照票 89 的形状：主体→SID 白名单），
      并在**同一枚 commit**里给"没接之前它确实宽"的**基线读数**（票 89 的 AC#1 就是这个形状，抄它）。
- [ ] **AC#3** 判"不封"的每一处要有**反向钉子**：一条用例说明"为什么这里宽是故意的"
      （例：日志的读者是用户自己的其他工具 ⇒ 断言它**必须**保持可读）。**"没接"与"故意不接"在代码上要能区分。**
- [ ] **AC#4** 变异：抽**一处**接线退回 `os.WriteFile(0o600)` ⇒ 它的 `icacls` 判据必须红
      （证明 Windows 上那串数字确实不落地，而不是"这次恰好没别人"）。锚点=承载行为那一行，同链 grep 证落地，
      还原后 `git diff --quiet` 证干净；**编译失败不算变异**。
- [ ] **AC#5** 门禁（按包 scope）**加上收尾必跑 `sh scripts/d22scan.sh` 纯净快照 rc=0**（票 94 AC#4 立的规矩，
      本票碰 5 个包，正是最容易把全仓 ban 弄红的那种票）；`gofmt -l`/`gofumpt -l` 空、`go vet` rc=0、
      `go test -count=2` 各包 rc=0 并逐条点名 SKIP/FAIL。
      ⚠ `go test ./cmd/wisp/` 本机是加载期 `0xc0000135` 的**既有红**，不要去追，如实登记。

## Rules（本仓固定）

`git commit -q -F - -- <显式路径> <<'MSGEOF'`（引号 heredoc）；提交前 `git diff --cached --name-only`；
禁 `git add -A`/`.`、`--amend`/`reset`/`rebase`/`stash`/`checkout .`（A34，共树）；**不 push**；
不在仓内建 worktree（A38④，快照 `git archive <sha> | tar -x -C /tmp/<带会话后缀>`）；
票面 append-only（改行前先读；标题前插段落要把标题重抄，`git diff --numstat` 删除列必须 0）；
四种假绿逐条点名；**真机上测 junction/ACL 只在临时目录里玩，绝不往用户真实数据目录写**；
数字不达标写 FAIL 附数字，不许调阈值、不许挑运气那次、多样本全报。
15 次工具调用内交回第一枚 checkpoint；每次提交同步 Status + 勾框 + `next=`；接近轮数上限主动收尾留断点。

## Progress log（append-only）

- 2026-09-21 18:0x（`agent-ticket95`，**第一枚 checkpoint：AC#1 前三类裁定 + 那条模型矛盾的取证**）：
  先跑 `date`：本机 18:01 CST。开工前现场核对（`git status`）：`scripts/d22scan.sh` 与
  `internal/winsec/seam_guard_windows_test.go` 有别人的未提交改动 ⇒ 我这两个文件/目录一字未动。

  ### 一、那条自相矛盾：我用证据判了，判给**台账**（票 95 表里那一格错），并登记"该更正"

  判据问题是编排者给的两句：**校验发生在读取时还是写入时？对象是每个模型文件还是一份索引/清单？**
  两处原文我都读了：票 95 表第 20 行（`archive.go:41/:73/:78` 标"要封"，理由是"这些文件之后被当作可信
  模型读回"）与台账 `docs/reports/pending-and-issues.md:1433`（A68⑤："**模型缓存不封**——公开签名的模型文件
  （C29 走签名清单），读到不产生泄露，封了反而挡住多实例复用"）、`:1124`（A76③ 登记这条矛盾本身）。

  取证（全部为本机实读源码，file:line 为我这次核过的当前行号）：
  1. **读取时**校验：`internal/models/downloader.go:170-174`——`Ensure` 的缓存命中分支在返回安装目录之前
     必须过 `m.VerifyDir(entry, installDir)`；`:161` 的 `local_override` 分支同样过它。校验**不在写入侧独活**。
  2. **每个文件**校验：`VerifyDir` 在 `downloader.go:519-531` 遍历 `entry.InstalledFiles()` 逐个
     `verifyFileHash(p, want.SHA256, want.SizeBytes)`；`InstalledFiles()` 的定义在
     `internal/models/manifest.go:218-229`——归档件走 `a.Archive.Files`（清单里逐条列出的成员），
     非归档件走 `a.Path/SHA256/SizeBytes` ⇒ 覆盖面是**每个落盘文件**，不是一份索引。
  3. 校验基准的权威来源：`internal/models/doc.go:6-16`——清单本身用 minisign 离线验签
     （`manifest.go`/`minisign.go`，公钥在 buildinfo），**先验签再联网**，hash 只从清单取（`downloader.go:53-64`、
     `manager` 构造期双门 `:99-107`）⇒ 换一份清单要换构建期公钥，不是同机另一账户能做的事。
  4. 解包产物离开 staging 之前还有一道全量门：`downloader.go:263-276`（staging 内逐个 `verifyFileHash`，
     失败 ⇒ `RemoveAll(staging)` 拒绝）。
  ⇒ **结论**："未过签名的文件也会被当作可信模型读回"这个前提**不成立**；票 95 表第 20 行括号里那句
  "之后被当作可信模型读回"恰恰是该类**已被逐文件校验**的位置，不是密封的理由。
  ⇒ 判：**`internal/models/` 四类落点（staging `:224`/`:287`、解包 `:73`/`:78`/`:41`、安装目录 `:566`）不封**，
  按 AC#3 给**反向钉子**（钉住"读取时逐文件校验"这两条事实本身，任一条退化 ⇒ 用例红 ⇒ 裁定必须重开）。
  ⚠ **改台账不是我的职权**：`docs/reports/**` 一字未动；这里登记的是"票 95 表第 20 行该更正"，交编排者落账。
  ⚠ 顺带一条诚实修正：台账那句"封了反而挡住多实例复用"的**代价**说得不准——`winsec` 封的是
  当前用户 SID（+SYSTEM+Administrators，见 `winsec.go:16-19`），**同一用户的其它实例照样读得到**，
  被封的只有别的账户。⇒ 不封的真正理由是"读到不产生泄露（内容公开且逐文件校验）"，不是"怕挡到复用"。

  ### 二、AC#1 其余两类（本 checkpoint 判完）

  - **④ 凭据迁移 `internal/secret/migrate.go:154`/`:168`——票面引用已腐坏，本票无活可接**。
    实读当前源码：该两处今天**已经**是 `winsec.PrivateFile(backupPath, raw, 0o600)`（`:169`）、
    `winsec.SealFile(backupPath)`（`:174`，既存备份的修复腿）、`winsec.PrivateFile(tmpPath, out, 0o600)`（`:189`），
    是票 89 退回单补的。**泄露后果**：备份是迁移前那份配置，装着用户迁移前的**全部明文 api_key**
    （源码注释 `:157-163` 自己点名），同机另一账户读到＝密钥泄露。**封了谁会读不到**：只有
    `secret.MigratePlaintext` 自己（写后不再读回）与用户手看的备份文件——同一用户仍可读。⇒ 无代价。
    ⚠ 但票面 17:4x ②那条"**`MigratePlaintext` 至今无生产调用方**"我复核为**真**：
    全仓 `grep -rn MigratePlaintext` 只有 `internal/secret/migrate.go:66/:86` 的定义与 `*_test.go`
    ⇒ 票 89 那三处密封今天是"能力就绪、没人叫它"。**接线点在装配根（`cmd/wisp/`）＝别人地界**
    ，且它改的是"启动即改写用户配置文件"这种产品行为 ⇒ **本票不接、登记**（见末尾 R-95 清单），
    `next=` 编排者裁装配根归属。
  - **③ 配置写路径 `internal/config/migrate.go:83` + `internal/config/parse.go:200/:213`——判"要封"，本票接线**。
    真实写入格式（读一处）：`migrate.go:82-85` 把**迁移前的原始 TOML 字节** `raw` 整份写进
    `path + ".bak-" + strconv.Itoa(from)`，模式参数 `0o600` 在 Windows 不落地；
    `parse.go:200-218` 是 `atomicWrite`：`os.CreateTemp(parent, ".wisp-config-*.tmp")` 写完后
    `os.Chmod(tmpName, 0o600)`（同样不落地）再 `os.Rename` 到 `config.toml`——**rename 保留描述符**
    （同一事实在 `secret/migrate.go:185-188` 已被写死），所以那个 pre-rename 临时文件的宽 ACL
    **就是最终 `config.toml` 的 ACL**。
    逐字段判"里面实际会出现什么"（票面 17:4x ①要求）：`internal/config/schema.go` 的
    `providers.<n>.api_key` 走 `RefPrefixDPAPI` 引用（`secret/configrefs.go:16` 一族），
    但**迁移前**的配置里同一个键可以是**明文**（正是 `secret.MigratePlaintext` 存在的理由，
    `secret/migrate.go:157-163`），另有 `[models]` 镜像 URL、`local_override` 绝对路径、
    `[llm]` 端点、日志目录路径。⇒ 泄露给同机另一账户的后果＝**可读的配置在票 89 之后仍然宽**：
    `config.toml` 与其 `.bak-<schema>` 明文备份在别的账户眼里是明文文件。
    **封了谁会读不到**：`config.Load`/`atomicWrite` 与配置编辑器都以**同一用户**身份跑（`winsec` 授的是当前用户
    SID），⇒ 只有"别的本地账户读配置"这条路被关；`wisp doctor` 若以另一账户跑会读不到（本票不跑那种形态）。
    ⇒ 接 `migrate.go:83`（`winsec.PrivateFile`）与 `parse.go` 的 temp（`winsec.PrivateFile` 取代
    `CreateTemp`+`Chmod`，保留独占语义）。

  未判：⑤ 日志（**做证据但不接线、不写反向钉子**，代价两条路写清交回）、⑥ 球位置、⑦ doctor/SLO（**只登记**）。
  next= 见本条末尾。

- 2026-09-21 17:4x（编排者，**票 89 复验把两笔账明确转到本票**，登记不扩写）：
  ① `internal/config` 的写路径与它的 `.bak-<schema>` 备份**今天仍然宽**（复验代理实测点名，
     与票 89 已修的 `secret/migrate.go:154` 同族、同形状）⇒ 进本票 AC#1 的裁定清单，
     **注意它比 `secret` 那处更微妙：配置里可能是端点 + 密钥 ref，也可能含明文**，要逐字段读一遍再判。
  ② **`MigratePlaintext` 至今无生产调用方** ⇒ 票 89 给它加的密封今天是"能力就绪、没人叫它"。
     ⚠ 本票接线时**顺手把这条一起接**（否则密封代码会永远绿着却永不生效，
     而台账上会留着"票 89 已修迁移备份"这句**只在被调用时成立**的话）。
  next= 不变（等票 89/94 的下游稳定后开工；AC#1 的裁定优先）。

- 2026-09-21 16:1x（编排者，**AC#1 的裁定先落一半：日志=要，模型缓存=不要**）：
  `acceptor-ticket89` 在给票 89 收尾时顺手答了我预留给本票的那句问（它的原话档位=**代理建议，未背书**）：
  - **日志（`internal/observe/logging.go:72/:244`，今天 `0o644`）⇒ 要封**，
    但它加了一条**施工约束**：必须用 `winsec.SealFile` 那条路，
    **不能对日志目录用 `SealDir` 的传播**——否则"清空/截断日志"会被权限改动带着变成写失败
    （这正是我留给本票 AC#1 要的"封了之后谁会读不到"的具体一半：**读者是日志自己所在的文件句柄**）。
  - **模型缓存/下载产物（`internal/models/`）⇒ 不要封**，理由它写的是"代价零"：
    那些东西是**公开签名的模型文件**（C29 走签名清单），别人读到不产生泄露，封了反而挡住多实例复用。
  ⚠ 两条都**不改变本票 AC#1 仍需独立成立**：派单时要求接单人自己把"里面实际会出现什么"读一遍再签，
  不许照抄我这两行。
  - 另外它把票 89 的一处遗漏算进了本票的邻居：`config.toml.bak`（**迁移前的明文密钥备份**）今天仍走
    装饰性 `0o600` —— 那条**归票 89 补**（同包、它 own 该写路径），**不在本票**，见票 89 的退回单第 2 条。
  next= 仍然等票 89 的验收判据落定（它现在被退回修复）⇒ 本票排在其后；`internal/winsec/` 由票 94 先走完。

- 2026-09-21 15:3x（编排者）：建票。票 89 交件时列了"本票范围外的同族落点"七处并写"**接线是一行级替换，
  但要有人裁决日志/模型缓存算不算私有数据**"。我把它落成这张票而**不**塞回 89：89 的六框已经按它自己的界做完了，
  把 7 个包追加进去会让一张已交件的票重新变胖（而且 `cmd/wisp/`、`internal/ball/` 都有活人）。
  ⚠ **AC#1 排在工作之前**：裁定没做完不许开始接线，因为"封日志"是可观察的行为变化，不是纯加固。
  ⚠ 本票**依赖票 89 的验收结论**（替代判据够不够、`SealDir` 传播那条覆盖面主张是否被独立证实），
  所以现在**不派单**；`internal/winsec/` 的语义也在票 94 手里，别同时改。
  next= 等 `acceptor-ticket89` 与 `agent-ticket94` 交件 ⇒ 那时再派，派之前我先自己把"日志算不算私有数据"
  的**代价**量一遍（谁在读日志：`wisp doctor`、票 66 的 SLO 链、用户自己的 tail）。
