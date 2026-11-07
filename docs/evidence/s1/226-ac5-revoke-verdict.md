# 票 226 AC#5 裁定表（非实现者＝编排者，2026-11-07 22:3x）

## 一句话
**AC#5 已达成，本格可勾。** 票面那句「今天 SetAllowedDirs 零调用者零用例」**已过期**——它今天有 **1 枚生产调用者 ＋ 3 枚用例**，且我已逐条现跑复认。

## 票面字面判据（逐字抄）
> **AC#5 可撤销那一半要有真路径**：撤销一条长期规则必须**既能落库、又能被答复语法触发**（今天 SetAllowedDirs 零调用者零用例）。

## 两半逐条取证

### 第一半「能落库」
| 项 | 盘上事实（现跑） |
|---|---|
| 落库函数 | `internal/config/allowdirs.go:109` `func (m *Manager) SetAllowedDirs(dirs []string) error` |
| 它改什么 | 只重写 `[fs] allowed_dirs` **这一个键**，文件其余行不动（writeguard 守卫写） |
| 用例 | `internal/config/writeguard_226_test.go:308` TestAC5SetAllowedDirsPersistsARevocationWithoutAClobber；`internal/config/writeguard_revoke_226_test.go:22` TestAC5RevokedRuleIsGoneForTheNextStart |
| 我现跑读数 | **两条都 PASS**（`go test ./internal/config -run TestAC5 -v` ⇒ `ok 0.127s`） |

### 第二半「能被答复语法触发」（★票面说缺的就是这一半）
| 项 | 盘上事实（现跑） |
|---|---|
| 落点 | `cmd/wisp/approval_revoke.go`（commit `48abebc41`，标题逐字：`226:AC#5触发半边落地——答复语法新增`） |
| 动词 | `:54` revokeVerb ＝ revoke；`:55` revokeRoute ＝ native/revoke；`:56` revokeTool ＝ config.allow_dir |
| 语法 | `:59` revokeSyntax ＝ `[fs] allowed_dirs += D:/某目录`（接受两种形状：长期流程印出的原文行，或裸目录） |
| **注册点** | **`cmd/wisp/approval_reply.go:641`** `return s.revokeRule(strings.TrimSpace(corr + " " + arg))` ⇒ **答复语法真的分派到它** |
| 落库调用 | `approval_revoke.go:84` `s.rt.mgr.SetAllowedDirs(listWithoutDir(before.FS.AllowedDirs, dir))` |
| 写后复核 | `:96-102` 重读 manager 视图断言「真的不在了」，**不靠 nil error 推断** |
| 拒绝路径 | ①没指名目录 ⇒ 拒；②rt/mgr 未装配 ⇒ 拒且「配置一个字节都没动」；③进程账上没有这条规则 ⇒ 拒且报当前名单 |
| 用例 | `cmd/wisp/approval_revoke_226_test.go`：`:97` TestReplyRevokeRemovesTheStoredRuleForTheNextStart／`:157` TestReplyRevokeRefusesARuleTheProcessDoesNotHold／`:196` TestRevokeVerbRefusalsDoNotTouchTheConfigWriter |
| **我现跑读数** | **三枚全 PASS**：2.59s／2.12s／0.00s ⇒ `ok 4.888s` |

## 票面那句「零调用者零用例」的现量否证（逐字读数）

```
grep -rn SetAllowedDirs --include=*.go internal/ cmd/ | grep -v _test
internal/config/allowdirs.go:109:func (m *Manager) SetAllowedDirs(dirs []string) error   <- 定义
cmd/wisp/approval_revoke.go:84:  ... SetAllowedDirs(...)                                <- 生产调用者 1 枚
（其余为注释/文档行）
```

⇒ **生产调用者 0 → 1**，用例 **0 → 3**（cmd/wisp 三枚）＋ **2**（internal/config 两枚）。

## 一处必须写明的诚实边界
- **撤销的生效档位是「下一次启动」**，不是立即。源码自己写在返回串里（`approval_revoke.go:110-112` 逐字）：
  > `本次运行仍按旧名单判定（撤销与「一直」同轨，生效都在下一次启动）`
- 这不是缺陷，是**与「一直」那条加宽路径同轨的设计**（`approval_revoke.go:16-24` 逐字解释为什么撤销不弹第二张卡：收紧不需要再确认，`PLAN.md:1645-1646` 冻结的再确认只针对**放宽**）。
- **⇒ 若机主要的是「撤销后立刻生效」**，那是另一个裁定（改生效档位＝改契约语义），**不在本格射程**，须另开票。

## 裁定
**AC#5 勾。** 两半都在盘上、都被我现跑复认、都有拒绝路径与写后复核。票面那句「零调用者零用例」作废（原句不抹，就地更正）。

**撤销口令**：本条裁定＝改回票面勾选（无产码）。