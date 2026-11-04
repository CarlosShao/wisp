# 403 — `internal/risk` 测试包今天**编译不过**：`provenance_test.go` import `internal/plugin`，而 `plugin/command.go` import `internal/risk` ⇒ import cycle，`go test ./internal/risk/` = `[setup failed]`

- **Status**：待派（**产码/测试重构活**，非只读）。来源＝编排者 2026-11-04 13:0x 现跑发现（在为票 234 取 AC#3/AC#4 的 risk 门禁读数时撞上；台账 A1479/A1480）。这是**回归**：票 160 以 broken 状态收口（漏检）。

## 复现（原样可跑）

```
cd "/d/work/workspace/projects plans/Wisp"
export PATH="$PWD/build:$PWD/third_party/sherpa-onnx:$PATH"
go test ./internal/risk/ -count=1
```

输出（逐字）：

```
# github.com/CarlosShao/wisp/internal/risk
package github.com/CarlosShao/wisp/internal/risk
	imports github.com/CarlosShao/wisp/internal/plugin from provenance_test.go
	imports github.com/CarlosShao/wisp/internal/risk from command.go: import cycle not allowed in test
FAIL	github.com/CarlosShao/wisp/internal/risk [setup failed]
FAIL
```

## 根因（现读，含引入时刻）

- `internal/risk/provenance_test.go:12` `"github.com/CarlosShao/wisp/internal/plugin"` —— 由 **`f576cf08`（票 160 / 160-r1 第②格 AC#2 换形状）** 引入。
- `internal/plugin/command.go:42` `"github.com/CarlosShao/wisp/internal/risk"` —— 自 **`5438d104e`（票 52 D46 命令插件）** 就存在。
- ⇒ `internal/risk` 的测试二进制两条边成环：risk(test) → plugin → risk。`go test` 在测试构建阶段即判 `import cycle not allowed in test`，**根本没进用例**。

## 影响（为什么这不是小事）

- `go test ./internal/risk/` **整包不可编译**（不是红一两个用例，是 setup failed）⇒ 以 risk 测试为分母的一切门禁读数不可得。
- 直接受阻：**票 234** AC#3 的 G-3（四包 `-count=2 -v`）与 G-6（容器三包）分母含 `internal/risk`，AC#4 终态 roster 也含它 ⇒ 234 无法如实收口（本票修好前 234 相应格保持未勾）。
- 其他以 risk 为分母的读数票（如 170/350 一族）同受影响。
- 这是"**越收越坏**"：票 160 交的 import 把整包测试锁死，验收没跑到 `go test ./internal/risk/`。

## 判据（AC 草稿）

- [ ] **AC#1 复现**：上面命令得 `FAIL [setup failed]` 且含逐字 `import cycle not allowed in test`；记录起手锚 HEAD。
- [ ] **AC#2 断环（先裁形状，五选一并写明选了哪支）**：
  - (甲) `provenance_test.go` 改成**外部测试包** `package risk_test`，只走导出 API（若那组测试不需要 risk 未导出符号）；
  - (乙) 把"需要 plugin 的那组测试"**移出 `internal/risk`** 到 `internal/plugin`（那里 import risk 合法，无环）；
  - (丙) 把 plugin 侧对 risk 的依赖**抽成更小的子包**（如 `internal/risk/classify`）让两边都能 import（最重，需先裁）；
  - (丁) 若 plugin→risk 的那条 import 实际只用到极少符号，给 risk 提供一个**不反向依赖**的窄接口/子包；
  - (戊) 其他能断环且 ⛔ 不牺牲覆盖的形状（须写明）。
  - ⛔ **不许**为过编译而删测试、`t.Skip`、或把断言放宽。
- [ ] **AC#3 修后终态**：`go test ./internal/risk/ -count=1` 达终态 `rc=0`（或如实具名红，但 ⛔ 不许再是 `setup failed`）；`go vet ./internal/risk/` `rc=0`。
- [ ] **AC#4 门禁**：改动文件 `gofumpt.exe -l` 为空（`/d/work/base/gopath/bin/gofumpt.exe`）；`sh scripts/d22scan.sh` `rc=0`；不新增 import 环（`go list -deps` 不报 cycle）。
- [ ] **AC#5 回填 234**：本票修好并推后，回填 234 AC#3(G-3/G-6)/AC#4 的 risk 分母读数使其可收口。

## 禁区

零改产品行为（本票只断测试 import 环，**运行时行为一字不动**）；不动 `PLAN.md`／`docs/specs/**`／`thresholds.go`／golden／`allowlist.txt`；不动三枚冻结件（`internal/panel/tokens_fourway_test.go`／`l2_grant_boundary_test.go`／`internal/perm/ticket90_persist_test.go`）；`frontend/**`／`design/**` 零写；不新增契约（D 号）。`internal/risk` 是共享编译面，本票是**唯一写腿**，派单前先现跑核它没别腿在写。

## 与其它票的关系

- **阻塞票 234**：234 AC#3 的 G-3/G-6 与 AC#4 分母含 `internal/risk`；本票不修，234 无法收口。
- **回归源＝票 160**（`f576cf08`）：160 已 `-done`，不回改它；本票是新回归修复票。修好后可在 160 追加一句指向（append-only）。
- **票 350**（internal/risk prefixed-symbol census）与本票不同：350 是"检测器词表"覆盖，本票是"测试包根本编译不过"。

## Progress log

- [2026-11-04 13:0x +08] agent=orchestrator did=立票（零产码）：现跑 `go test ./internal/risk/` 得 `FAIL [setup failed]`＋`import cycle not allowed in test`（provenance_test.go:12 → plugin → command.go:42 → risk）；根因追到 `f576cf08`（票160）引入 plugin import、`5438d104e`（票52）plugin import risk。next=派产码腿断环（先裁甲乙丙）、修后回填 234。
