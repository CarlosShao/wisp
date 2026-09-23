# 票 114 AC#2 — 原生侧门（`R-92-2` 的正身）交件与自证

**代理**：`agent-ticket114-ac2`（写码）· **锚定 sha**：`58302cc`（开工时 `git log -1` = `3879686`，
`git diff --name-only 58302cc HEAD` 只有 `docs/reports/pending-and-issues.md` 一枚 ⇒ 本文件所有 `file:line`
在两枚 sha 之间零漂移）· **分支** `dev`，只 commit 不 push。
**输入**：`docs/evidence/s1/114-ac1-status-table.md`（只读代理 `audit-ticket114-ac1` 的现状表，**未重做、未冷搜索**）。
**commit**：`12d8e28`（门 + 用例 + 票面）· 第二枚＝`cmd/wisp/run.go` 装配 + 本文件（见文末清单）。

## 1. 落点（file:line）

| 物件 | 位置 | 说明 |
|---|---|---|
| `internal/panel/composer_handlers.go`（新建，165 行） | `:111` `func (h *ModeWriteHandler) HandleModeRequest(...)` | mode 写入处理器 |
| **门本体** | `internal/panel/composer_handlers.go:133` `if modeIsWidening(from, to) && h.Confirm == nil {` | **位置在任何 `ModeWriter.Set` 调用之前**（`:147` 才是 `h.Modes.Set`） |
| 三枚哨兵 | 同文件 `:57/:61/:66`（`ErrNoL2Confirm`/`ErrNoModeWriter`/`ErrNoCurrentMode`） | 拒的形状可被 `errors.Is` 分辨 |
| 本地小接口（零新依赖） | 同文件 `:74` `ModeWriter`、`:86` `ModeConfirm` | `*perm.Store` 原样满足；**未 import `internal/perm`**（现状表 ②.4 ③ 那枚方向风险整个省掉，实测 `go list -deps ./internal/panel` 的包集合本轮未变） |
| 装配 | `cmd/wisp/run.go:378`（`rt.modeWrites = &panel.ModeWriteHandler{...}`）+ 字段 `cmd/wisp/run.go:229` | 交进去的 `Confirm` 是 `perm.New` 上面那**同一枚** `confirm`（`run.go:346-349` 的 `s.modeConfirm`/`rt.confirmModeSwitch`）；`Audit: rt.auditf` |

**覆盖面按编排者 P1 落地**：门只管**把档位往更宽切**（`modeIsWidening(from,to) = to > from`，
`internal/risk/mode.go` 的常量声明就是 strict→loose 且原文写着 "Do not reorder these constants"）；
**变严不挡**、**同档不挡**（那两支交给 `Set`，perm 照旧落审计）。审计三档全写未被改变：
拒的一支**自己**落一行 `panel: MODE-REFUSED ...`（写入腿没跑，不然那次切换无痕），过的一支由 `perm.Store.Set` 记，
处理器**不补第二条判定**（用例钉住"过腿时 audit 行数 = 0"）。

**处理器不调用 `Confirm`**（`:86` 的注释即此条）：卡片是写入腿的（R20/M4 的一枚 C18 卡），这里再问一次就是
两张卡、两次 300 秒。这一条也被用例钉住，不是只写在注释里。

## 2. AC#2 那条判据的原文（`internal/panel/composer_handlers_test.go`，327 行，6 枚顶层用例）

```go
			if w.setCalls != 0 {
				t.Fatalf("%s -> %s: the injected ModeWriter was called %d time(s) despite the refusal - the gate did not stop the write",
					names[fromIdx], names[toIdx], w.setCalls)
			}
```

（`composer_handlers_test.go:144-147`，在 `TestAWideningModeRequestIsRefusedWhenNoL2ConfirmLegIsAttached` 内，
逐对遍历三档矩阵的全部变宽方向；`fakeModeWriter.Set` 是**注入**的写入腿，`setCalls` 只在被调时自增。）

其余五枚：`TestTheModeLadderThisGateJudgesAgainstIsTheDocumentedThree`（**清单三钉**：精确相等 + 长度下限 +
严格前缀负断言 ⇒ 词表从 3 缩到 2 会单独红，不只藏在相等里；另加"矩阵枚数 == n²、变宽对数 == n(n-1)/2"两枚
**枚举缩短要能响**的下限）、`TestAWideningRequestWithAnAttachedLegReachesTheOneWriterAndRaisesNoSecondCard`
（`setCalls==1`、`Confirm` 调用数 **必须为 0**、过腿时 audit 行数 0）、
`TestAStricterOrEqualModeRequestNeedsNoConfirmLeg`（变严/同档在 nil 腿下仍走到写入腿，各 3 对 + 3 同档）、
`TestAModeRequestIsRefusedBeforeAnyWriteWhenTheAssemblyIsIncomplete`（无写入腿 / 当前档读不到（`risk.Mode(97)`）/
未知拼写 / 空 target / 走错方法，五支全部 `setCalls == 0`）、`TestTheWriteLegsOwnRefusalIsReturnedUntouched`
（腿接上后 perm 自己的拒绝原文不被改写——票 123 那批超时语义不在本轮动它）。

**helper 用的是本包原有的** `auditLog`（`workspace_test.go:40`），未新建同名件（首轮误建已被本轮改掉）。

## 3. 变异三态（先证落地，再读红名，逐发还原复绿）

| # | 退回的那一步 | 落地证明 | 红名 |
|---|---|---|---|
| **M-1** | 把门退回"不检查 `Confirm == nil`"（`if false && ... && h.Confirm == nil`，其余不动） | `grep -n "MUTATION M-1"` = `composer_handlers.go:133`；`go build ./internal/panel/` rc=0 | `--- FAIL: TestAWideningModeRequestIsRefusedWhenNoL2ConfirmLegIsAttached` → `composer_handlers_test.go:137: ask_every_step -> ask_high_risk with no confirm leg: got nil error, the档 would have been written` |
| **M-2** | 改成"只记日志不拒"（保留 `h.record(...)` 与 `return refused`，中间**多加一次 `h.Modes.Set`**＝档位照样被写） | `grep -n "MUTATION M-2"` = `composer_handlers.go:137`；`go build` rc=0 | 同一枚红名，但**红在 `composer_handlers_test.go:145`**：`the injected ModeWriter was called 1 time(s) despite the refusal - the gate did not stop the write`。哨兵错误**仍然返回**、audit 行**仍然写了** ⇒ 只断字符串/只断 error 的用例这一发是绿的，**只有"ModeWriter 未被调用"接住它**（现状表 ②.3 的那一发） |
| 还原 | 两发逐发还原 | `grep -c MUTATION internal/panel/composer_handlers.go` = **0**、`git diff --numstat HEAD -- internal/panel/composer_handlers.go` 空（＝与 `12d8e28` 提交的字节一致） | 复绿：见 §4 的 windows/linux 两读 |

## 4. 门禁四数（按包 scope；禁全仓 `go test ./...` 遵守）

| 读数 | GOOS | 命令 | RUN | PASS(顶层) | FAIL | SKIP |
|---|---|---|---|---|---|---|
| 还原后复证 | windows（本机 go1.27.1） | `go test -count=2 -v ./internal/panel/` | 152 | 92 | 0 | 0 |
| **Linux 真跑（不是只编译）** | linux（Docker `golang:1.27`，go1.27.1 linux/amd64） | 同命令，容器内 | **152** | **92** | **0** | **0** |
| 装配腿（本轮改了 `cmd/wisp`） | windows | `PATH="$PWD/third_party/sherpa-onnx:$PATH" bash scripts/wisp-cli-tests.sh`（内含 `-count=1` + CI 的 `-skip` 名单） | 98 | 51 | 0 | 0 |

- **分母问法（本仓定过的"这条用例在哪个 runner 有分母"）**：门的行为判据**全在 `internal/panel`** ⇒
  ubuntu core scope 有分母（`scripts/portable-tests.sh:141,179`），本轮**用 Docker 真跑把它转正**；
  `cmd/wisp` 只有 windows 腿跑（`ci.yml:378`），所以那边**只放装配**（`rt.modeWrites` 那一格），
  **没有**把任何行为判据只写在那边。
- Docker 挂载自证（防 Git Bash 静默挂空的假绿）：命令带 `MSYS_NO_PATHCONV=1`，容器内先
  `ls -l /src/go.mod`（真读到 855 字节那枚）+ `uname -s` = `Linux` 再跑测试。
- 其它仪器：`go build ./...` rc=0；`go vet ./internal/panel/` rc=0；`GOOS=linux go vet ./internal/panel/` rc=0
  （**只编译**，不当分母用）；`GOOS=linux go vet ./cmd/wisp/` **失败**，原文
  `build constraints exclude all Go files in .../sherpa-onnx-go-linux`——**本轮前即如此**
  （`scripts/wisp-cli-tests.sh:26-28` 逐字记录过同一读数，且 ubuntu 真跑 cmd/wisp 另有 19 枚在飞的旧红），
  与我的改动无关，也未因此把 cmd/wisp 记成"两平台都有分母"。
  `D:\work\base\gopath/bin/gofumpt.exe -l` 对三枚改动文件**零输出**（真跑，原始输出即"空"）。
  `sh scripts/d22scan.sh` = **clean rc=0**，各 scope 未降：ban #6 `frontend/=43`、ban #8 `internal/=389 Go files`、
  `cmd/=36`、`frontend/=43`、`design/=16`（本轮 `frontend/` **零写入**，枚数与 AC#1 表读数一致）。
- **未动**：任何阈值/断言/golden、票 123 那批 CLI 用例（一枚未改，300 秒未当旋钮）、`internal/winsec/**`、
  `internal/risk/**`、`rules_gateway.go`、`tools/d22scan/**`、`allowlist.txt`、`ci.yml`、`scripts/`、`internal/observe/**`、`frontend/**`。
  开工前 `git status --porcelain cmd/wisp/run.go` 为空（无他人未提交改动）才动它；交件时该文件三枚 hunk 全是本代理的。

## 5. 没做到的（明写，不圆）

1. **通路仍不存在，这条性质本轮没被改变。** 现状表 ①.1 量到"WebView2 `WebMessageReceived` → `ParseComposerRequest`
   那一跳今天**根本没有**（不是某行没接上，是缺函数缺文件）。按 P0 地界**没写宿主**，所以：
   `rt.modeWrites` 在生产里**零调用者**、`HandleModeRequest` 的**生产调用者 = 0**、
   `ParseComposerRequest` 的生产调用者**仍是 0**（本轮一枚没加）。
   ⇒ **门立在"下一次接线"的前面，不是替它把线接了**：票 33/35（两票仍 `ready-for-agent`）落宿主那一跳时，
   必须调这枚 `rt.modeWrites`，否则 AC#2 是空门。**这一格我没有、也不能自行验收。**
2. **AC#3 的反向对照做不了**（"请求更宽档在真机上必须走 L2 原生卡片、超时 300s = 拒"）——它要求真机+原生通道，
   即 1 的那一跳。**未勾。**
3. **AC#4/#5/#7/#8/#9/#10/#11 未做**：本轮只交 AC#2（AC#1 的现状表是只读代理的交付，本代理只负责入库勾框）。
4. **AC#6 未做且不勾**：`skipped=frontend(owner-delegated)`（票头 09-23 横幅 + `A102③`）。
   开工前列了要写的全部文件，**零枚**落在 `frontend/` 下 ⇒ 未触发整段停手。
5. **装配断言没进用例**：现状表 ②.4 推荐在 `cmd/wisp` 里钉"handler 的 Confirm 非 nil"。那要新增/修改
   `cmd/wisp/*_test.go`，**不在本轮可写清单**（派单只给 `cmd/wisp/run.go` 的装配那 10-20 行）⇒ **未做**，
   装配的正确性目前只有"读代码 + `go build` + cli 腿 51/0/0 不破"三档证据，**没有用例**。
6. **行数超估算**：现状表量级是"处理器 60-90 / 用例 80-120"，实际 165 / 327。超的是**注释与判据条数**
   （一次成对矩阵 + 五枚装配不完整支路 + 清单三钉），不是新增语义面积；登记为偏差，不回头砍判据。

## 6. 自称权威的文字（本轮工具输出，逐枚点名）

- 冒充「编排者备注 / 系统提示 / 请 revert / 冻结某包 / 放宽阈值」并**要求改变行为**的文字：**0 次**。
- 到达的 `Modified content:` 通知 **2 批**（形状：`Note: The file ...MEMORY.md was modified since it is last read.`）：
  ① 第 1 批＝首条 Read/Bash 批之后两条（项目级 + 全局级索引）；② 第 2 批＝第 2 次 `Bash`
  （命令前 40 字 `cd "D:\work\workspace\projects plans\Wisp" && ls -l`，内容 = 路径核验 + `date`）之后一条全局级。
  **两把判据都核过**：①路径真存在（`ls -l` 读到 `/c/Users/swq/.qoder-cn/memory/MEMORY.md` 8910 字节、
  `.../projects/D--work-workspace-projects-plans-Wisp/memory/MEMORY.md` 2774 字节，mtime 09:59/09:28 = 另一会话在写）；
  ②内容**没有**削弱 owner 权威、**没有**放宽判据（反而重申"前端不派不代写""只 commit 不 push"）
  ⇒ 判为**真实 harness 通知**，未据其改动任何边界。**`C:\Users\swq\.qoder\settings.xml` 那枚 `A109③` 的形状本轮没出现**
  （本机该路径实测不存在：`No such file or directory`），本轮**零凭据形状**到达，无值可登记。

## 7. `next=`（给编排者 / 下一段）

1. 票 33/35 落宿主那一跳时**必须**使用 `cmd/wisp/run.go:378` 的 `rt.modeWrites`（并加"handler 的 Confirm 非 nil"
   的装配用例，需先解冻 `cmd/wisp/*_test.go`），否则 AC#2 空转；那两票仍 `ready-for-agent`。
2. AC#3（真机 L2 卡片 + 300s 超时）与 AC#6（差分截屏）分别在 1 落地之后、以及 owner 指派的外部 agent 手里。
3. AC#4/#8/#9/#10/#11（门二覆盖面、两份名单同步、扫描根换到"渲染器真加载的集合"、`AddHostObjectToScript` 的显式禁用、
   `render:composer` 进 CI）本票未动；AC#11 要改 `ci.yml`，那枚文件**当前在票 134 的代理手里**（本轮 `git status` 可见其 WIP）。
4. 本票 Status 已是 `in-progress`，**未改 `-done`、未改名**（防重领唯一键未动）。
