# 440 — 仓里有**两个互相独立的 ACL 解析器**（`internal/models/parseACELines` 与 `internal/winsec/aclSIDs`），漂移无人认领

- Status: **done（判甲：合并进 internal/winsec/aclparse.go，测试全绿）**。
- 来路：票 **238** 的结案腿 `238-seal-r2`（2026-11-13，提交 `673da11df`，件 `docs/evidence/s1/238-data-root-seal.md` §AC#4③）。它逐条给三条欠账归属时，这条**唯一悬空**——issues 树里除了票 238 自己，零票号、零台账号。
- 关联：票 238（owner 否决，已结案）· 票 132（封根，已落地 `logging.go:99`）· 票 121（ACL SID 判据）· 票 95（`winsec` 不碰祖先，避开模型夹）。

## 现量（票 238 腿本程跑，我复核过 `git grep`）

| 解析器 | 位置 | 返回什么 |
|---|---|---|
| `parseACELines` | `internal/models/acl_sid_121_test.go:57` | **整条 ACE 行**（`[]aceLine`，保留 flag 文本） |
| `aclSIDs` | `internal/winsec/acl_windows_test.go:86` | **只取 SID 与名字**（`sids, names []string`） |

**两者都在 `_test.go` 里**（测试辅助），**互相不知道对方存在**，**没有任何票认领它们的关系**。

## 为什么这是个真问题（不是洁癖）

票 238 腿的原话点出的形状：**「沿用现有助手就等于沿用它的盲区」**——`aclSIDs` 只取 SID，**丢掉了 flag 文本**，所以任何"这枚 ACE 是允许还是拒绝、是不是继承来的"的判据，用它写就**判不了**。
反过来 `parseACELines` 保留 flag，但**只在 `internal/models` 的测试里**，`internal/winsec` 的测试拿不到它。

**⇒ 后果：写封条判据的人会随手挑一个，而两个的射程不同，挑错了判据就是空的。** 这与票 180／439 那一族同根（"仪器射程与规格射程不是一回事"）。

## 本票要做的

- [x] **AC#0 只读比对（⛔ 产码）**：逐条列出两个解析器的**射程差**：
  - 各自能回答什么问题、⛔ 回答什么问题
  - 对同一份 `icacls` 输出，两者各返回什么（**给一份真实 dump 的逐行对照**）
  - **特别验**：`parseACELines` 的 `aceLine` 结构体里有没有 flag 字段；`aclSIDs` 丢了哪些信息
- [x] **AC#1 判归（三选一，写理由）**：
  - **(甲) 合并成一个**：抽到某个非测试包（如 `internal/winsec` 或新 `internal/aclparse`），两边共用
  - **(乙) 明确分工**：两个都留，但**各自写清射程**，并在注释里指回对方（防下一个人挑错）
  - **(丙) 归口他票**：已有票认领 ⇒ 具名票号
  - **⚠ 判据＝主流 harness 怎么做**（⛔ "动作小"／"风险低"）。测试辅助抽成共享包是常见做法，但**先量清楚两者是不是真的该是同一个东西**。
- [x] **AC#2 若判"甲"**：合并后**两边原有判据全部仍绿**（⛔ 删任何断言），并给成对判据（合并前红／合并后绿 或 反之）。
- [x] **AC#3 契约轴**：⛔ 改 `docs/PLAN.md`／`docs/specs/**`／`thresholds.go`／golden／`allowlist.txt`。
- [x] **AC#4 门禁**：`go test -count=1 ./internal/models/ ./internal/winsec/` ＋ `gofumpt -l internal cmd tools scripts frontend`（无输出）＋ `scripts/d22scan.sh`（干净）。
- [x] **AC#5 边界声明**：明确写"本票收什么、⛔ 收什么"。

## 写面
`internal/models/**` 与 `internal/winsec/**` 的 `_test.go`（判"甲"时可能新增一个非测试包）。
⚠ **开工前先 `git status --porcelain internal/models internal/winsec` 记基线**。

## 交付四件（缺一算失败）
1. 改动（只对判"甲"的）
2. 证据件 `docs/evidence/s1/440-acl-parser-drift.md`
3. 票面追加 Progress log。**⛔ 翻 AC 框**（编排者做）。
4. commit，**显式 pathspec**，**⛔ push**。

## Progress log

- **2026-11-13 腿 440-v1**：AC#0 完成——射程差逐条（含 deny/继承两问的归属）＋真实 dump 逐行对照（`%TEMP%\wisp440-6ff42f5d`，`/deny` 实造，实测 deny 形状 `BUILTIN\Users:(DENY)(DC)`、继承形状 `NT AUTHORITY\SYSTEM:(I)(F)`，标记全在 `:(` 之后的载荷里）＋结构体验（`aceLine` 有 flag 字段＝`Rights` 载荷原文；`aclSIDs` 丢的正是它）。**比对额外发现**：`:(` 解析惯用法全仓 **5 份**不是 2 份（observe/config/agent 各一份内联，已登记进证据件 §0.4，不归本票收）。AC#1 判归＝**甲**，落点 `internal/winsec/aclparse.go`（winsec 本体内新增非测试文件，**非新包**——理由：models→winsec 无循环（grep 现跑，反向才有）＋ `portable-tests.sh` GUARD C/D 的静态名册机制＋scripts/ 不在本票写面，详见证据件 AC#1）。AC#2 完成：models 本地定义改 `var/type = winsec.X` 委托、`aclSIDs` 改薄包装（签名不变，33 处调用零改动），测试函数计数 models 3→4 / winsec 5→5（零断言删除）；成对判据 `TestAC440TheSharedParserKeepsTheRightsPayload` 钉进 portable 腿＋变异实证（`Rights` 填空 → AC440+AC47 红 → 还原复绿）。AC#3 契约件零触碰。AC#4 三条门禁合并后终态全绿（models 10.642s / winsec 26.788s；gofumpt 空；d22scan clean，本机经 `D:\work\soft\Git\bin\bash.exe` 调用——裸 `bash` 解析到 WSL 无 bash，环境注记在证据件）。AC#5 边界声明在证据件（收：双副本合并＋射程注释＋238 AC#4③ 归属落账；不收：三份内联副本、aclSIDs 调用面迁移、策略层、CI 名册）。票 238 AC#4③ 的悬账自此有归属＝本票。commit 显式 pathspec（四 Go 件＋证据件＋本票面），未 push。
---

## ★结案（2026-11-13 12:3x，编排者翻框）— **判甲：合并成一个共享解析器**

### 三把尺（我现跑）
- **提交真在**：`b7bc01482`（6 files, +414/-92）
- **件真在**：`docs/evidence/s1/440-acl-parser-drift.md`（14541B）＋ `internal/winsec/aclparse.go`（6148B，新）
- **自己复跑**：`go test -count=1 ./internal/models/ ./internal/winsec/` ⇒ **`ok models 11.752s`** ＋ **`ok winsec 28.913s``

### 腿判"甲"的理由（我认）
1. **它们是同一个东西的两份漂移副本**，不是两个该分工的东西：`aclSIDs` 每行输出＝`parseACELines` 对应行的 Trustee 列＋翻译。连 `:(`／trailer／路径粘连三个陷阱的判据都**逐字同源、注释互引同一个 bug**。
2. **主流 harness 判据**：重复的测试解析逻辑**向上抽成被测包导出面、测试留薄包装**，是 Go 仓常规做法（`iotest`/`httptest` 先例）。
3. **价值只有合并买得到**：下一个写判据的人**只有一个解析器可挑**。注释互指（乙）买不到——**本仓 5 份副本就是互指失败的证据**。
4. **落点选 `winsec` 本体不是新包**：①依赖边只有 `models→winsec` 干净（反向才成环）；②`portable-tests.sh` GUARD C/D 用静态名册核对，新包当天就顶红，而 `scripts/**` 不在本票写面。

### ★腿顺手量出的事实（比票面更严重）
**全仓 `icacls` 出现在 35 个文件里**（我复核）。票面只点了 2 份解析器，腿找到 **5 份** `:(` 解析器（observe/config/agent 各还有一份内联），**全都拿不到 deny/继承**。
⇒ **合并前，全仓没有任何一份解析器把「允许还是拒绝」递到判据手里。** 这条比"两份漂移"严重得多。

### 改动与判据
- 测试函数：models **3→4**（＋成对判据）、winsec **5→5**——**零断言删除**；winsec 侧 **33 处 `aclSIDs` 调用签名不变、零改动**。
- **成对判据** `TestAC440TheSharedParserKeepsTheRightsPayload`：绿半边＝合并后 `(DENY)(DC)`/`(I)(F)` 都在判据手里；**变异实证**＝把 `Rights` 临时填空 ⇒ TestAC440＋TestAC47 四子句全红 ⇒ 精确还原复绿。
- 门禁：`gofumpt` 空输出、`d22scan` clean、两包全绿。

**撤销口令：「440 回退」**——`git revert b7bc01482` 即复原两份解析器（但会重新引入"判据拿不到 flag"的洞）。

