# 80 — `blacklist_overrides` 配置被解析、被方向审计，但**从来没有接到 `risk.Gate` 上**

**Status:** claimed（写码代理已认领，正在跑 AC#1 只读穷尽清单）
**Type:** 契约与实现脱节（一个面向用户的开关是死线）
**Blocks:** nothing · **Blocked by:** nothing（`internal/risk` 此刻无人写：票 75 已交、票 72 只碰过 `pathresolver.go`/`blacklist.go`）
**Packages:** `internal/config/`、`risk.Gate` 的构造点、以及新建的接线测试。**禁改冻结面**：
`internal/risk/assessor.go`、`internal/risk/pathresolver*.go`、`rules_gateway.go`、`docs/PLAN.md`、`docs/specs/*.md`。
**发现的**：票 75 的接续代理（2026-09-21 收尾报告第 6 节第 1 条）。

## 缺陷

票 75 的实测链条：
- `internal/config/schema.go:457` 有 `BlacklistOverrides`，会被解析；
- `internal/config/manager.go:356` 甚至对它做了**方向审计**（说明写它的人预期它生效）；
- 但 `risk.Gate` 的 `bOverrides` map **在生产路径上没有任何调用者往里灌数据** ——
  那个"单个文件的 B 档豁免"能力**只活在测试里**。

后果有两种方向都不好的读法：
1. 用户/文档以为可以豁免某条 B 档，实际**静默无效**（配置项说谎）；
2. 反过来，若哪天接上去，等于**没人审过的放行通道**被打开。

## 判据前先做的一件事（AC#1，只读）

**"grep 零命中"不等于"没有保护"，也不等于"没有生效"**。所以 AC#1 必须先穷尽列举
`BlacklistOverrides` / `bOverrides` / 任何等价的 override 容器在**全仓**的每一次读与写
（含 `cmd/`、`internal/tools/`、桥接层、以及测试里的构造），并给出**逐条位置 + 该行是读还是写**的清单，
然后才允许下"它是死线"的结论。若清单证明**其实有生效路径**（例如通过另一个字段名接上），
那本票就地关闭并登记为"报告者计数错误"，**不许为了有活干而造修复**。

## AC（1:1，一框一行裁决表，验收在 `docs/evidence/s1/80-*.md`）

- [ ] **AC#1** 上面那张"每次读/写"清单进本票 Progress log，并明确结论：死线 / 已生效（后者 ⇒ 本票转登记，不写码）。
- [ ] **AC#2** **契约判定**：在 `docs/PLAN.md` 与 `docs/specs/SPEC-06*.md`（C19/C26 相关条目）里查
  `blacklist_overrides` 是否被写成**用户可用能力**。
  - 若**是**（契约要求它生效）⇒ 接线是**修 bug**，本票继续 AC#3。
  - 若**契约没有这条**、或它属于"放行侧放宽"⇒ **停下来上报编排者**，由编排者判是否 D22。
    ⚠ 不许用"先接上再说"绕过：放行侧的开关是安全边界，不是配置便利。
- [ ] **AC#3** 接线落地 + 一条**真跑在 `risk.Gate` 上**的用例：同一份配置里写 override，
  证明 (a) 被点名的那一个文件从 B 变成放行，(b) **同目录下其它文件仍然 B**，
  (c) override 的形状必须是 C26 交回的 canonical 形式（票 72 的 R17 裁定：
  **放行侧只认解析后的形式**，`bOverrides` 的 key 不接受拼写变体）。
- [ ] **AC#4** 变异双向：(i) 把接线断开 ⇒ AC#3 必须红；(ii) 把 (b) 的"其它文件仍 B"断言指向 override
  ⇒ 必须红。两次都要在**同一条 `&&` 链里先 grep 证明改动真落地**再跑，还原后 `git diff --quiet` 证干净。
- [ ] **AC#5** 门禁（**只跑自己碰到的包**，共树不跑整仓）：`gofmt -l <pkgs>` 空、
  `go vet <pkgs>` rc=0、`go test -count=2 <pkgs>` rc=0，且**逐跑点名** `--- SKIP`/`--- FAIL` 的行数与测试名；
  `GOOS=linux go vet <pkgs>` rc=0（Linux 编译哑弹是 A44/A49 的老坑）。

## Rules

`git commit -q -F - -- <显式路径> <<'MSGEOF'`（引号 heredoc，A31）；**禁** `git add -A`/`.`；
每次 commit 前 `git diff --cached --name-only` 核对清单；
**禁** `--amend`/`reset`/`rebase`/`stash`/`checkout .`（A34，共树）；**不 push**（编排者的活）；
不在仓目录里建 worktree（A38④），要纯净树用 `git archive HEAD | tar -x -C /tmp/...`；
不碰别人的票面（尤其票 75/72/79 的 `.md`）；15 次工具调用内交第一枚 checkpoint commit，
之后每次 commit 同步 Status + 勾框 + `next=`。

## Progress log（append-only，最新在下）

### L1 — 认领 + checkpoint（2026-09-21）

已把 Status 从 `open` 改为 `claimed`。本轮次序：AC#1 只读穷尽清单 ⇒ AC#2 契约判定 ⇒ 只有两者都过才允许写实现码。
`next=` 跑 `BlacklistOverrides` / `bOverrides` / 等价 override 容器在全仓的每次读与写清单。
