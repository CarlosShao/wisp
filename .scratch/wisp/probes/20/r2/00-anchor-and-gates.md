# 20-r2 — 票 20 第 7 框（`[fs] allowed_dirs` 首次使用询问流最小闭子集）：锚、三格读数、门禁

- 腿＝`20-r2`（实现者）。起手锚 `125e2255`（15:2x +08，`git status` 里 `internal/tools`/`internal/config` 零他人脏件）；
  交付时 HEAD 已移 `9b51859b`（期间只有 `internal/panel` 两枚测试件被 `114-a5`/编排者改，`git diff --stat 125e2255..9b51859b -- internal/ cmd/` 零命中 `internal/tools`/`internal/config`）。
- 窗口（逐报）：`16:01:56–16:02:17` 整包 `./internal/tools/`（rc=0）；`16:07:53–16:09:08` 复跑（rc=0）；其余窄跑窗口见 logs 内时间戳。开工/收工 `tasklist` 均无 `go.exe`（其间 `114-a5` 只动 `internal/panel` 测试件，与 `internal/tools` 零文件重叠）。

## 0. 落点与边界（三条逐条）

- **(a)** 新产码 `internal/tools/refusal.go`（117 行）：`RefusalKind` 两个机器可辨值 ——
  `RefusalOutOfAllowlist = "out-of-allowlist"`（refusal.go:53，**唯一**允许被首次使用询问取用的类）、
  `RefusalRedTeamDenied = "redteam-denied"`（:58）。判定 `refusalKindOf`（:72）跑在**同一批 raw paths**（C19 判过的输入）上：
  任一 `Canonicalize` 出错（junction/symlink 穿越等解析器拒绝）或等级是 Deny ⇒ 红队类；全部规范化成功且至少一条不在任何根下 ⇒ 越界类；其余 ⇒ 空。
  取值在 `bridge.go:361`（routing 拒绝分支）只对**被拒**调用设置，经 `Decision.Refusal`（`gate.go:37`，OnDecision 交给宿主）
  与审计行 `refusal=%s`（`bridge.go:1162`）两个宿主可达面离开桥；人话 `Reason` 一字未动。
  命名安全：`Decision` 字段扫描（`ticket90_test.go` 禁 `allow/approve/grant` 字样字段）不受影响，既有断言零红。
- **(b)** **零新授权产码**：批准后的追加只走既有通道 —— `config.Manager.AddAllowedDir`（ticket 201 已建、`cmd/wisp/approval_always.go:134` 已在驱动）
  与文件级放宽的**方向引擎 + `ConfirmLocked`**（`manager.go` 的 `fsDirection` fi.allowed_dirs 行＋`planLocked`）。本腿不新增第二条路：
  `PathCanonicalizer` 方法名册钉（只读/判定 + 仅收窄的 workspace 面）＋“`internal/tools` 产码不 import `internal/config`”钉，两条都在 `Test20r2NoSecondAuthorizationPath`。
- **(c)** 三断言＝下节逐格读数。**“下一次”＝下一次装配（重启形状），不是同运行内热换**——这是仓里已落地的裁定，不是本腿选的口径：
  `cmd/wisp/config_reload.go` 头部逐字写着“swapping rt.paths mid-run would need an exported reset on internal/tools, which this leg may not mint (no new exported names)”，
  ticket 201 的卡面也写「下一次启动生效」。故 approve 格**两半都断言**：已拒过的那枚桥复检仍拒（禁「这次先放行」），下一枚装配同路径 L0 直读。

## 1. 逐格读数（改前红 / 突变红，逐字；绿＝`logs/10-baseline-green.txt`／`24`）

**基线**：`go test -run 'Test20r2' ./internal/tools/ -count=1 -v` ⇒ PASS 8 顶层（含 4 子项）0 FAIL 0 SKIP（10-baseline-green.txt）。

- **(a) 改前红**（MUT-B0：`dec.Refusal = …` 一行改成不设值，即改前树行为；日志 `20-mut-*`）：
  - `refusal_20r2_test.go:120: 越界拒绝的机器可辨值 = "", want "out-of-allowlist"（措辞不参与判定）`
  - `refusal_20r2_test.go:169: A 档拒绝的机器可辨值 = "", want "redteam-denied"（Deny 永远不可问）`
  - `refusal_20r2_windows_test.go:50: junction 拒绝的机器可辨值 = "", want "redteam-denied"`
  - `refusal_20r2_windows_test.go:100: 8.3 拒绝的机器可辨值 = "", want "redteam-denied"`
  - （共 10 句红，另含流程三格与 8.3 边界格；全文与红句逐字见 `logs/20-mut-b0.txt`）
- **(a) 突变证明（题面要求的“把可辨原因值并回红队字符串”）**＝MUT-M1：`RefusalOutOfAllowlist = "redteam-denied"` ⇒ 4 枚钉红（`logs/21-mut-m1.txt`）：
  - `refusal_20r2_test.go:116: 两类拒绝共用同一个字符串 "redteam-denied" —— 被拒原因不再可机器区分（票 20 第 7 框 (a)）`
  - `refusal_20r2_test.go:172: A 档拒绝带上了可问类值：等于让红队门自己开口要社工点击 (c)③`
  - `refusal_20r2_windows_test.go:53: junction 拒绝带上了可问类值：注入文档就此能驱动「要不要授权这个目录」的询问 (c)③`
  - `refusal_20r2_windows_test.go:103: 8.3 拒绝带上了可问类值 (c)③`
- **(c)② / (b)**（`Test20r2FirstUseAskFlowGoesOnlyThroughConfig` 四子项 + `Test20r2DirectionEngineAdjudicatesFileLevelAdditions`）：
  空名单下 fs.read 越界 ⇒ 可问类拒绝；**拒绝**格＝名单仍空、下一次装配仍拒；**批准**格＝`AddAllowedDir` 落盘（config.toml 里逐字有该行、
  新开一枚 Manager 读得回、同一枚旧桥复检仍拒、**下一枚装配同路径 L0 直读**）；**收紧**格＝兄弟目录仍拒（放宽只授权那一个目录）。
  重载格＝手改文件 → `CheckAndReload`：hook 收 `section=fs, keys=[fs.allowed_dirs]`、答否 ⇒ 仍空且同路径仍拒；同一次放宽再写一遍（10ms 间隔过 mtime+size 判据）答可 ⇒ 生效且同路径成功。
- **(c)③ junction / 8.3 绝不产生询问**（真夹具，零 `t.Skip`，`refusal_20r2_windows_test.go`）：
  - 真 `mklink /J`（jun 在授权根内、目标在外，先做阳性对照读出 94 字节级真内容）→ L2/R2 拒绝 ⇒ `redteam-denied`、`askable20r2` 为假；
  - 真 `GetShortPathNameW` 短名打 A 档文件（在授权根内，只 R3 能拦）⇒ Deny/`permission_denied`、门 0/0 ⇒ `redteam-denied`、不可问；
  - 边界格（已裁定、如实钉住）：**非黑名单**越界文件的短名拼法 ⇒ 与长名**同一类**＝越界可问类。依据＝同包已提交的
    `TestShortNameSpellingGetsTheSameVerdictAsTheLongOne` 逐字注释 “8.3 is deliberately NOT a resolver-refused class”。
    即：红队类＝{解析器拒绝 ∪ Deny}（票 18 Case 1–4 的四种拼法在这两支落齐），可问类只跟**目标**走、不跟拼法走。
    ⚠ 若 owner 要更严读法（“任何 8.3 拼法一律红队”），是 refusal.go 一行的事，该边界格就会红（测试注释已写明）。
- **(c)③ 的突变证明**＝MUT-M2：解析器出错支改成返回越界类（junction 变可问）⇒ `refusal_20r2_windows_test.go:50: junction 拒绝的机器可辨值 = "out-of-allowlist", want "redteam-denied"`（唯一红，正是安全内核那枚钉；`logs/22-mut-m2.txt`）。
- **(b) 钉咬证明**＝MUT-M3：给 `PathCanonicalizer` 临时加 `AddRoot`（禁形的第二扇门）⇒
  `refusal_20r2_test.go:363: PathCanonicalizer grew AddRoot. If it WIDENS what InAllowlist authorizes outside NewPathCanonicalizer, that is the second authorization path ticket 20 frame 7 (b) forbids.`（`logs/23-mut-m3.txt`）。
- 三处突变全部还原：`grep -c MUTATION internal/tools/*.go` 零命中、`paths.go` 零 diff、复跑全绿（`24-after-mutations-restored.txt`）。

## 2. 门禁（全部现跑，rc 在 `logs/30-gates-merged.txt`）

- `go build ./...` rc=0；`go vet ./internal/tools/` rc=0。
- `gofumpt -l`（GOPATH/bin/gofumpt.exe，动过/新增 6 件）空、`gofmt -l` 空。
- `go test ./internal/tools/ -count=1` rc=0（`ok 33.514s`，整包 33 顶层 PASS、0 FAIL、0 SKIP——既有用例零红，含审计行/Decision 家族）；
  `-race -run 'Test20r2'` rc=0；`go test -run 'TestManagerLocked' ./internal/config/` rc=0（先例族复跑）。
- 长度门 `scripts/check-path-length-budget.sh` rc=0 GREEN（60=60；本腿新增件名短，不进帽口）。
- `d22scan`：仪器自测 `-self-test` rc=0（40 向检查全过）；**全仓扫描 rc=1＝61 条 finding，逐条全在 `frontend/fixtures/settings-v2/*.png`**
  （mtime 16:01，**未跟踪件，属机主外部前端 agent 的在飞产物**，`frontend/**` 是本编队明令零写面）；`internal/` 与 `cmd/` 命中 0。
  本腿自己撞出的唯一一条（`refusal_20r2_test.go:224` 字符串里的 🔒，ban#8 不豁免字符串）已改（注释里的保留：注释豁免），改后逐跑该文件 0 命中。
  ⇒ 该 rc=1 是他人脏件所致，非本腿产出；编排者可在前端件落地/清理后复跑取绿。

## 3. 具名没做（不越界，逐条）

1. **询问提示本身（native 卡 / GUI）不在本腿写面**：把「要不要把这个目录加进来」弹成卡、以及这张卡的回答落进 `AddAllowedDir`，
   是宿主侧接线（`cmd/wisp`，本腿明令触则停手上报；GUI 完整版归票 39）。本腿交付的是该流程唯一缺的**判定面**：可辨值 + 可问谓词契约 + 三断言。
2. **邻接险情（上报不动）**：ticket 201 的「一直」通道 `ReplyCard.WideningRule`（`internal/agent/approval/replies.go:503`）从卡面 `Paths` 推目录，
   而 junction 卡面的目标是**带标记串**（`bridge.go displayPaths`：`"…(无法规范化: risk: path traverses a reparse point (junction/symlink)…"`），
   `dirOf`（`batch.go:91`）按最后一个 `/\` 切会把标记里的 `/` 当分隔符 ⇒ 可推出一串含标记碎片的“目录”并让第二张卡显示它。
   这是 `internal/agent/approval`＋`cmd/wisp` 面（本腿零写面），且属“红队形状该不该能驱动放宽询问”的判定层问题 ⇒ **只登记，一行未改**。
3. **同运行内热换根集合**：明令未做（见 §0 (b)），仓里已有具名裁定（`config_reload.go` 头部）。
4. 未跑：`-tags winlive`、整包 `cmd/wisp`（314 死区）、跨卷 junction/硬链接（后者两层都看不见，票 20 第 5 框已登记）。

## 4. 日志清单（`logs/`；提交 7 枚＝下表，raw 逐跑件全部留盘未删、未提交）

| 件 | 内容 |
|---|---|
| `10-baseline-green.txt` | `-run Test20r2 -v`：8 顶层 PASS／0 FAIL／0 SKIP（含 4 子项） |
| `20-mut-b0.txt` | MUT-B0（值不设＝改前行为）：rc=1 全文＋10 句红逐字 |
| `21-mut-m1.txt` | MUT-M1（并回红队字符串）：rc=1 全文＋4 句红逐字 |
| `22-mut-m2.txt` | MUT-M2（junction 变可问）：rc=1 全文＋唯一红句逐字 |
| `23-mut-m3.txt` | MUT-M3（加 `AddRoot` 第二扇门）：rc=1 全文＋钉红句逐字 |
| `24-after-mutations-restored.txt` | 三处突变还原后复跑：8 顶层 PASS，`grep MUTATION` 零命中 |
| `30-gates-merged.txt` | 八项门禁合并（整包/递归 build/vet/d22scan 脚本归因＋直跑 61-finding 逐目录计数/仪器自测/先例/config 长度门/race） |

留盘未提交的 raw 逐跑件（含 `20-mut-b0-reds.txt` 等 4 枚红句抽取、`31/31b/32/33/34` 逐跑原文）都在同一目录，随用随查。
