# 440 — 仓里有**两个互相独立的 ACL 解析器**（`internal/models/parseACELines` 与 `internal/winsec/aclSIDs`），漂移无人认领

- Status: **ready-for-agent（先只读比对两者射程，再判合并/分工）**。
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

- [ ] **AC#0 只读比对（⛔ 产码）**：逐条列出两个解析器的**射程差**：
  - 各自能回答什么问题、⛔ 回答什么问题
  - 对同一份 `icacls` 输出，两者各返回什么（**给一份真实 dump 的逐行对照**）
  - **特别验**：`parseACELines` 的 `aceLine` 结构体里有没有 flag 字段；`aclSIDs` 丢了哪些信息
- [ ] **AC#1 判归（三选一，写理由）**：
  - **(甲) 合并成一个**：抽到某个非测试包（如 `internal/winsec` 或新 `internal/aclparse`），两边共用
  - **(乙) 明确分工**：两个都留，但**各自写清射程**，并在注释里指回对方（防下一个人挑错）
  - **(丙) 归口他票**：已有票认领 ⇒ 具名票号
  - **⚠ 判据＝主流 harness 怎么做**（⛔ "动作小"／"风险低"）。测试辅助抽成共享包是常见做法，但**先量清楚两者是不是真的该是同一个东西**。
- [ ] **AC#2 若判"甲"**：合并后**两边原有判据全部仍绿**（⛔ 删任何断言），并给成对判据（合并前红／合并后绿 或 反之）。
- [ ] **AC#3 契约轴**：⛔ 改 `docs/PLAN.md`／`docs/specs/**`／`thresholds.go`／golden／`allowlist.txt`。
- [ ] **AC#4 门禁**：`go test -count=1 ./internal/models/ ./internal/winsec/` ＋ `gofumpt -l internal cmd tools scripts frontend`（无输出）＋ `scripts/d22scan.sh`（干净）。
- [ ] **AC#5 边界声明**：明确写"本票收什么、⛔ 收什么"。

## 写面
`internal/models/**` 与 `internal/winsec/**` 的 `_test.go`（判"甲"时可能新增一个非测试包）。
⚠ **开工前先 `git status --porcelain internal/models internal/winsec` 记基线**。

## 交付四件（缺一算失败）
1. 改动（只对判"甲"的）
2. 证据件 `docs/evidence/s1/440-acl-parser-drift.md`
3. 票面追加 Progress log。**⛔ 翻 AC 框**（编排者做）。
4. commit，**显式 pathspec**，**⛔ push**。
