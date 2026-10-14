# 133-v 面核（三枚提交）＋ 与 133-r1 读数的对照（含分歧登记）＋ 具名没做

==== 面核 ====
- `git show --stat da389ca5`：恰好两枚 cmd/wisp 测试件
  （`cmd/wisp/dispatch_hop_selfproof_133_test.go` +424；`cmd/wisp/leg_dispatch_gate_133_test.go` +66/-12）。
  逐行看 -12：全是该尺自己的四行旧文（minLegs133=11 的注释与常量、no-args 旧 claim 行、旧下限 t.Errorf 文案、
  `ns = append(n.t...` 一行）被替换 ⇒ **没有放宽任何断言**，没有第二枚生产文件的改动。
- `1740fd48`／`46dee9ef` 越面零：文件表＝票面 + `.scratch/wisp/probes/133/r1/**`（6 枚/4 枚）与中途对
  `05-ac6-ci.txt` 的收回重写；`git diff-tree` 三枚合扫 `^\.github/|thresholds|golden` **零命中**。
- 131 的门 `cmd/wisp/leg_sink_gate_131_test.go` 最后一次动于 `6f702aeb`（131 续单），本票三枚提交零触碰。
- `git status --porcelain cmd/wisp/`（本腿开工时）空 ⇒ 基线取自 HEAD；本腿全程仓内零写、变异全在快照。
- gofmt/gofumpt（v0.12.0 go1.27.1）复量：
  - `gofmt -l cmd/wisp/` = `models.go` 一枚；`gofumpt -l cmd/wisp/` = `models.go` + `panel_transport_35r2_test.go`。
    两枚都是既有债：`models.go` 是 CRLF 伪影（blob LF / 盘上 CRLF，`git status` 看不见），
    `panel_transport_35r2_test.go` 由 `18d43f8e` 自认（"HEAD 真格式债"）。本票两枚仪器文件逐枚 clean。
  - **CI 形原命令 `gofumpt -l . tools/d22scan tools/mockllm` 的复量＝非空**：HEAD 快照上点出
    20 枚 `.scratch/wisp/probes/**` 已入库样品 ＋ `cmd/wisp/panel_transport_35r2_test.go`（`models.go` 在
    LF 快照里干净）；活树上另有别的腿在飞件与 CRLF 伪影。这与 r1 票面那句"CI 逐字同形 … 空"**不一致**，
    且 r1 的 12-ac5-fmt-vet.log 只留了 `gofumpt -l cmd/wisp/` 这一条（仓外全日志目录里也 grep 不到
    那条 CI 形命令的原档）⇒ 登记"该句缺原档"，不当它成立也不当它有牙。仓自己的口径已在
    `ci.yml:189-196` 写死（ticket 161 AC#7：`.scratch/**` 是故意坏样品，raw walk 结构性不可满足；
    真门槛是 tracked-set 那一步 `attrib.sh --tracked-only`），**本票面无增债**。
- d22scan：`sh scripts/d22scan.sh` rc=0 纯净（活树 21:22）：
  bans#1-5 internal/=237 cmd/=40 · #6 frontend/=115 · #7 internal/tools/=24 ·
  #8 design/=39 frontend/=115 internal/=557 cmd/=129。
  同树控制（base 快照，21:23，非 git 仓）：…除 `#8 design/=30`（gitignore 姿态差，r1 已登记同形）
  与 `#8 internal/=556`、`#8 cmd/=128` 外逐数相同 ⇒ 相对同日控制**不降**，+1/+1 归因到 21:2x 活树里
  别的腿在飞/新落件（本腿零产码、零新文件进 cmd/internal）。

==== 与 r1 读数的对照（分歧逐条） ====
1. **行号漂移（信息级）**：r1 01-ac1-five-shots.txt 引的红名行号（`leg_dispatch_gate_133_test.go:203/:209/:219/:185`）
   与其**已提交**的同文件行号（`:216/:222/:241/:199/:202`）不同；原因是那批读数取自 a7b88997 上的
   **未提交工作树版本**（`git diff a7b88997 HEAD -- 该文件` = +66/-12），提交前又加了注释行。
   红名文本与判据逐字相同 ⇒ 无实质分歧；本腿读数一律取自提交后的 `b75089cc`。
2. **gofumpt CI 形"空"那句缺原档**：见上。其余 gofumpt/四数/窗口读数与它一致。
3. **AC#6 的一处口径收窄（不是它的错）**：它引的 run `38554642217` 的 sha 是 `72db0f39`，
   早于仪器 commit `da389ca5`（19:25 提交）⇒ 那趟 CI 的 step 分母里只有
   `TestAC1AC2DispatchHopGate133` 与 `TestAC4EveryLegIsNailedOrRuled` 两枚（全 PASS），
   自证腿/逐枚可区分两枚**从未在任何远程 run 里跑过**。它票面写的"本腿四枚用例两次都在该 step 分母里"
   指的是它**本机**两趟 `bash scripts/wisp-cli-tests.sh`（其 18 号原档可复读，四枚全 PASS）——
   本腿独立核过那份原档。结构上仪器已进 CI（`_test.go` 在 cmd/wisp，step 跑整包），但"远程真跑过的 run"这一枚
   读数**今天给不出**，如实登记。
4. **本腿自己的两处环境修正（不是分歧）**：① 快照无 `.git` ⇒ `TestCleanCheckoutBuilds_AC11` 红
   （`git rev-parse --show-toplevel: exit status 128`），第 1 趟字面全包的 21 枚红里只此一枚属此形；
   ② 快照无 `third_party` ⇒ `no native DLLs in ..\..\third_party\sherpa-onnx` 一族红（第 1 趟），
   拷入三枚 DLL 后同发转 PASS（10.69s）再跑第 2 趟。两处都已写进 03 的分母归因。

==== 具名没做 ====
1. **"两门都摘"的全包读数未重跑**（r1 的 no-doors 306/226/7/3）：本腿只做了 N-3 的"131 单跑"字面复核
   （红名只有 131 的门）。全包在本机是死区（见 03），不重复付 10 分钟。
2. **方法值 `v.M(...)` 支未造**、**跨包 install 支不接**（票面原文划出地界，照抄不测）。
3. **linux 侧分母未取**（本腿全程 windows/amd64；`_windows_test.go` 里 4 枚注册表钉的平台豁免
   只在文档第 3 条限制里被披露过，未在此复核）。
4. **远程 CI 里"本票两枚新用例真跑过"的 run 不存在**（见分歧 3）——不催不补，如实报。
5. 未重跑 `38551683523`（另一条同 step 的 run）；未重跑 123 批审批用例的单点（`da389ca5` 的 diff
   证其零触碰：无 -skip 增改、无阈值改动、文件表只有两枚测试件）。
6. `docs/reports/pending-and-issues.md`、`HANDOVER.md` 未动（派单禁令）；票面只 append 1 条 bullet；
   框架未翻、Status 未动。
