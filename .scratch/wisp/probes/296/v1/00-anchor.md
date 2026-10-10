# 296-v1 · 第 1 笔（commit-first）：射程声明＋HEAD＋内容锚＋尺

腿号 `296-v1`（非实现者对抗验收腿）。钟点＝`date` stdout 现插：`2026-10-10 09:10:21 +0800`。

## 0. 射程（我能裁什么、不能裁什么）

- 只裁三格：`AC#0`（射程普查）／`AC#1`（空串世界夹具到底有没有牙）／`AC#2`（甲-全那一发＋同批改写的那句过期注释）。
- ⛔ `AC#3`（真机）／`AC#4`（门禁）归编排者，他 22:2x 已翻勾 ⇒ 本腿不复跑真机、不重跑他跑过的四把门禁。
- ⛔ 零翻框：票面五格（`AC#0`..`AC#4`，尺＝`grep -cE '^- \[[x ]\]'` 于票面）一枚不改。
- ⛔ 零产码入库（突变台件跑完必还原）、⛔ 不动 `docs/**`／`frontend/**`／`design/**`／`PLAN.md`／`docs/specs/**`／`internal/observe/thresholds.go`／golden／`tools/d22scan/allowlist.txt`／三枚冻结件。
- ⛔ 不起任何进程（`build/wisp.exe`、`build/wisp-296.exe`、`cmd/balldebug` 全不起）。
- ⛔ 不碰只读腿 `299-a1` 的射程：`cmd/wisp/panel_host_windows.go`、`internal/panel/assets.go`、`.scratch/wisp/probes/299/a1/`。
- 写面只新建 `.md`，落点 `.scratch/wisp/probes/296/v1/`。⛔ `.sh`／`.ps1`／`.txt`／`.raw`／`.out`（根 `.gitignore` 第 8 行全仓 `*.out`）。

## 1. HEAD 与起点证据

- `git rev-parse HEAD` ⇒ `28b484ec6129557d72981565adc8214464bec9c3`（短号 `28b484ec`）。
  ⚠ **与派单转述不同处（具名报回）**：派单说"落地腿已交回六笔 `ff5c193e`→…→`448f5a57`（HEAD 现量请自己跑）"。
  盘上现量：`448f5a57` 之后还有 `29081a13`（编排者收口＋翻 AC#3/AC#4）与 `28b484ec`（Q-83/票 299）等笔 ⇒ **HEAD 不是 `448f5a57`，本腿所有读数锚在 `28b484ec`**。
- 起点 `git diff --stat -- cmd/wisp` ⇒ **空**（rc=0），即 `cmd/wisp` 工作树＝HEAD，无人在飞写它。

## 2. 三枚内容锚（⛔ 绝对行号不进任何派单，一律内容锚）

- 锚 A（被调方唯一的跳过条件）＝`internal/ball/hotkey_reload.go` 里 `Check()` 的
  `if r.exists && reflect.DeepEqual(r.applied, cfg) {`。
  尺＝`git grep -n 'DeepEqual\|all-empty\|IsZero' HEAD -- internal/ball/hotkey_reload.go`
  ⇒ HEAD 现量**只有 1 命中**（就是这条 `DeepEqual`），`all-empty`／`IsZero` **0 命中** ⇒ 全文件零条"全空就不重绑"守卫（与票面 18:2x 裁定节第二句一致）。
- 锚 B（甲-全落地点，HEAD 态）＝`cmd/wisp/resident_windows.go` 包级函数 `hotkeyReloadSource296(dataDir)` 的**两枚 `return` 都套 `ball.ApplyHotkeyDefaults`**：
  `return ball.ApplyHotkeyDefaults(ball.HotkeyConfig{})`（读不到文件那一支）与
  `return ball.ApplyHotkeyDefaults(ball.HotkeyConfig{Summon: c.Hotkey.Summon, ...})`（裸映射那一支）。
  尺＝`git grep -n 'ApplyHotkeyDefaults' HEAD -- cmd internal | grep -v '_test.go'`（结果见 `10-ac0-roster-check.md`）。
- 锚 C（夹具）＝`cmd/wisp/resident_hotkey_296_windows_test.go` 两枚用例；跑法＝
  `PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test ./cmd/wisp/ -count=1 -run 'Test296|Test258'`。

## 3. 本腿要交的判语

- `AC#0`：名册是否**整族**（⛔ 抽样当整族）；每枚枚数是否写清"哪把尺＋射程目录＋含不含 `_test.go` 与注释行"；是否漏票面点名的落点；`cmd/balldebug/main.go` 反例是否具名确认。
- `AC#1`：定向突变＝摘掉 `hotkeyReloadSource296` 两枚 `return` 的 `ApplyHotkeyDefaults` 套壳 ⇒ 指名跑那两枚用例 ⇒ **必须红**；并核两枚闭包不同＋四行空串逐字＝机主形状；★具名评估"提为包级函数"这一处偏离。
- `AC#2`：新注释有无新造承诺；`internal/ball` 零字节；`schema.go` 四枚默认值零动；`none`/`off` 零 normalize。

## 4. 台件纪律（本腿承诺）

1. 突变前把 `cmd/wisp/resident_windows.go` 按字节备份到**仓外** `$TEMP`，⛔ 仓内不留 `.go` 拷贝。
2. `trap`/`finally` 还原；还原后 `git diff --stat -- cmd/wisp` 为空＋`sha256sum` 与备份逐字一致。
3. 长跑命令输出先落文件再取；退码自己落一行 `rc=N`。
