# 票 172 落地腿 `172-r1` — 读数与凭据（2026-10-13 23:1x–23:4x +08）

- 派单＝`.scratch/wisp/probes/orch/20261013/txt/172r1-dispatch.txt`（＝腿消息）。分支 `dev`。
- 起手 `git log -1`＝`aea0dd0b`（`logs/00-head.txt` 逐字）；取证途中树在动（另一程提交到 `044b9066`/`1c1cac46`），本件每处读数各自带 HEAD。
- 写面＝`internal/tools/tier_consistency_172_{ruler,machine,cases}_test.go`（三枚，193/142/180 行）＋
  `.scratch/wisp/probes/172/r1/**` ＋票面尾 1 条 bullet。契约轴、CI、产码（除瞬态 overlay）零改动。

## 一、尺是什么、在哪

- `tier_consistency_172_ruler_test.go`：D34 表解析（给定路径，可指副本）＋机器侧对账（每枚落地工具走
  `pathArgs -> factsFor -> assessorFor(...).Assess`，即"这次调用机器实际给哪一档"）→ 每枚工具
  **文本档位集合必须 == 机器探针给出的档位集合**（不是"包含"）。
- `tier_consistency_172_machine_test.go`：机器名册（与 `cmd/wisp/run.go` 同三支 builder）＋13 条探针调用
  ＋文本侧"可指副本"演示件（`WISP172_TIER_DOC` 只作演示缝，默认＝`docs/PLAN.md`）。
- `tier_consistency_172_cases_test.go`：三发用例＋逐枚读数打印＋红句。

## 二、AC 逐格

### AC#1（文档说假话必须有人喊）＝**成立**
- **未修读数（本腿自复跑）**：把 `docs/PLAN.md` D34 表 `fs.edit` 那行的 RiskLevel 格 `**L2**` 改成 `**L1**`
  （副本在 `$TMPDIR/172-r1/mut/PLAN.172.docflip.md`，改 1 行），配 `-overlay` 跑 `./internal/tools/...`
  ⇒ **rc=0、零枚红**（`logs/10-baseline-docflip-overlay.txt`）。筛法复跑：`ReadFile|os.Open|go:embed|ParseFS|ReadDir`
  在 `internal/tools/*_test.go`＋`internal/risk/**` 里提到 `PLAN|SPEC|docs` 的 **0 命中**
  （`logs/11-unfixed-grep-census.txt`）。
- **修后**：`WISP172_TIER_DOC=<那份副本> go test -run Test172 -v` ⇒ **rc=1**，红句逐字：
  `fs.edit：文档写 L1、机器给 L2（探针：替换授权目录里既有文件的一小段 -> L2）—— 冻结文本与门给的档位不一致，两边必须由同一个 ticket 改齐`
  （`logs/21-ac1-doccopy-lie.txt`）。真 `PLAN.md` 全程零写入。
- **结构性事实（票 236 教训，本腿现量）**：`-overlay` 换 `docs/PLAN.md` ⇒ 本尺**照绿**（rc=0），因为
  `os.ReadFile` 不经过编译器 ⇒ 读盘型尺对 `-overlay` 不可见；红形演示只能走**可指副本**那条缝
  （`logs/22-ac1-docoverlay-invisible.txt`）。

### AC#2（代码说假话也必须有人喊）＝**成立**
- 变异：`internal/tools/fs_edit.go:390` 的 `Declared:     risk.L2,` → `risk.L1`（overlay 副本，源件 sha256
  前后同为 `7efedfffa504337b0f5fae07e31a7707f8111eefc57aac4ed34862bcfdbb2779`，见 `logs/30-ac4-contract-axis.txt`）。
- 全包 `-overlay` 跑 ⇒ **rc=1**，本尺红句逐字：
  `fs.edit：文档写 L2、机器给 L1（探针：替换授权目录里既有文件的一小段 -> L1）……`（`logs/20-ac2-overlay-fsedit-L1.txt:88`）。
  同发里另有 6 枚既有用例红（`TestAuditAllowlistScopeDoesNotInventAVerdict`、`TestFSEditRoutesToApprovalNotTheL1Window`、
  `TestFSEditRefusesAPathC26CannotCanonicalize`、`TestFSEditAppliesABatchOnDisk`、`TestFSRegistrationIsTheD34Roster`、
  `TestToolsDirectoryShape`）⇒ overlay 确证生效，不是静默 no-op。
- **未修码对照**（同一发变异 + 把本票三枚新件 overlay 删除＝模拟改前树）⇒ rc=1、6 枚旧红、**"文档/机器不一致"
  这句话 0 命中**（`logs/23-ac2-unfixed-tree.txt`，`grep -c 文档写`＝0）。
- **不是空心检**：机器侧不重列期望表，走的是生产判定路径；判据是集合相等 ⇒ 除 `Declared` 外，`PathParams`
  抽取与 Facts 钩子（如 `fs.write` 新建 L1／覆盖 L2 两支）也各自成读数。**摘掉哪处哪条不响**：
  摘/改 `Declared` ⇒ `Test172...GateGives` 红；摘掉注册（或改名）⇒ 落地名册钉红（`9 -> 8`，`landed` 与
  `tier172Landed` 比对）；改真文本那行 ⇒ 同一发红（AC#1）。

### AC#3（分母＝整张带档位工具行）＝**成立，读数如下**（`logs/60-census-green.txt`，未修码绿）
- 两处冻结文本逐字对账一致：各 **37 行 markdown**（其中 3 行含多枚工具）⇒ **(工具, 档位) 条 41 条**、
  去重工具名 **40 枚**；两文本逐枚档位集合相等（`Test172RulerReadsBothFrozenTablesAndTheyMirror`）。
- **对得上 9 枚**：`fs.delete` L2、`fs.edit` L2、`fs.list` L0/L2、`fs.move` L1、`fs.read` L0/L2、
  `fs.trash` L1、`fs.write` L1/L2、`task.cancel` L1、`task.output` L0。
- **对不上 0 枚**（今天没有两边档位不同的行）。
- **未落地 31 枚**（文本有行、机器无名册＝还没注册成 C4 工具，如实登记不跳过）：`app.launch` `clipboard.read`
  `clipboard.write` `doc.read` `file.open` `input.type` `list_tools` `media.control` `memory.recall` `memory.save`
  `notify` `plan.present` `reminder.cancel` `reminder.list` `reminder.set` `screen.capture` `search.apps`
  `search.content` `search.files` `shell.exec` `shell.session` `system.get` `system.power` `system.set`
  `task.list` `web.fetch` `web.open` `web.search` `window.focus` `window.list` `window.manage`。
  ⚠ 其中 `notify`/`list_tools` 是 S1 行而 `cmd/wisp/run.go` 有注释说明为何注册不进去（名字形状）——本尺只登记读数，不裁。
- **机器有而两处契约文本都没有行 ＝ 1 枚：`task.spawn`**（票 197 落地，L0 声明无契约行来路）。已按名册钉住
  （`tier172Undocumented`）：将来出现第二枚这样的工具 ⇒ 红并要求"要么带契约行，要么显式登记"。
  旁证：票 236 的只读普查腿 `236-a4` 也量到过"表里没有 `task.spawn` 行"且声明不判它（本尺把它从"没人提"变成"被钉住"）。

### AC#4（契约轴）＝**成立**（`logs/30-ac4-contract-axis.txt`）
- `fs_edit.go` sha256 起手/收手同一串（变异全在 `$TMPDIR`，源件零写入）；
  `git status --porcelain -- docs/PLAN.md docs/specs tools/d22scan` **空**；`allowlist.txt` 仍 7 行、`fs.edit` 0 命中；
  两处冻结文本 `fs.edit` 行现量在位（`PLAN.md:2537`、`SPEC-07:43`，均 L2）。本尺只读它们。

### AC#5（接不接 CI）＝**默认不接，未动任何 CI**
- `.github/**` 零改动（`git status --porcelain -- .github` 空）。若以后要接：只许**新增步骤**，动触发表/并发组/
  既有步骤＝停手上报（票 134 地界）。

## 三、门禁（每件自落 rc）

| 件 | 命令 | 读数 | 日志 |
|---|---|---|---|
| 构建 | `GOFLAGS= go build ./...` | rc=0 | `50-build-gofumpt.txt` |
| 形状 | `gofumpt -l <本票三枚新件>` / `<internal/tools/*.go>` | 空/空 | `50-build-gofumpt.txt` |
| 包测试（改后） | `go test -count=1 ./internal/tools/...` | rc=0（13.6s） | `51-postchange-test-d22scan.txt` |
| 包测试（改前） | 同上，配 `-overlay`（改文档那行） | rc=0（未修读数） | `10-baseline-docflip-overlay.txt` |
| 本票三发 | `go test -count=1 -run Test172 -v ./internal/tools/` | rc=0（3 PASS） | `60-census-green.txt` |
| D22 | `sh scripts/d22scan.sh` | rc=0（正控 PASS=35 FAIL=0；clean） | `51-postchange-test-d22scan.txt` |
| 长度门 | `sh scripts/check-path-length-budget.sh` | **rc=1（⛔ 非本票）** | `40-lengthgate.txt`、`41-lengthgate-red-attribution.txt` |

**长度门那发红不属本票**：它读的是 **git 跟踪路径**（本票三枚新件取证时还是 untracked），红因是名册漂移——
票 173 收口那次 `git mv … -done`（提交 `2a8ab9c4`）把 `.scratch/wisp/issues/173-…-done.md` 变成"超预算且不在名册"，
另一条反向红是"名册里有而树里已无"。修它＝动名册/改别人的票面名，越界 ⇒ 停手登记。

## 四、本腿知道没说到的

- **备注列的升级条款不在读数里**：D34 的"跨盘 -> L2"（`fs.move`）、"关闭 = L2"（`window.manage`）等写在**备注/动作**格
  （`window.manage` 的档位格确实给了 `L1 / L2`，`fs.move` 的没给）⇒ 本尺只对账 RiskLevel 格与两边都拿得到的读数；
  这些腿以及未落地工具的越界/升级腿，本尺今天读不到，未假装覆盖。
- **`fs.write` 覆盖支靠 Facts 钩子（R8）抬 L2**，本尺读到的是 `[L1 L2]` 与文本 `[L1 L2]` 相等 ⇒ 覆盖支被真的探到了；
  但"哪个机制抬的"（R8 vs R1 下界）不在本尺判据内（那是 `162-v3` §4① 的账）。
- **13 条探针调用**里 12 条是"范围内/越界"两种情境的胶水，不含档位期望值；第 13 条（`task.spawn`）只为
  "机器有而契约无"那一格服务。
