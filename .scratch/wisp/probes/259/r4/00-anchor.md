# 259-r4 起手锚（票 259 AC#4 载具前置）

起手现取时刻：`2026-10-08 18:54:41 +08`

## HEAD 与工作树

- `git log -1`：`0f7d4c53c5644423a2c3edcc07409f2b2dad7e1e Thu Oct 8 18:54:19 2026 +0800 票 259 四格翻勾＋A734 落账（凭据＝非实现者腿 259-v1 终裁）：…（AC#4 不成立＝cmd/wisp/subagent_selfapproval_197_test.go:109 仍两枚同值）`
- `git status --porcelain | head -20`（只登记，不处置）：
```
 M .gitignore
 M .scratch/wisp/probes/152/my152.py
 M .scratch/wisp/probes/161/r6/logs/flip-1.txt
 M .scratch/wisp/probes/161/r6/logs/flip-2.txt
 M .scratch/wisp/probes/161/r6/logs/flip-3.txt
 M .scratch/wisp/probes/161/r6/logs/flip-4.txt
 M .scratch/wisp/probes/161/r6/logs/flip-5.txt
 M .scratch/wisp/probes/161/r6/logs/flip-6.txt
 M .scratch/wisp/probes/161/r6/logs/flip-baseline.txt
 M .scratch/wisp/probes/161/r6/logs/flip-restored.txt
 M .scratch/wisp/probes/242/r3/logs/probe-routed.txt
 M .scratch/wisp/probes/268/v1/evidence.md
 D design/assets/base.css
 D design/assets/icons.js
 D design/assets/theme.js
 D design/assets/tokens.css
 M design/doubao/README.md
 M design/doubao/demo/app.js
 M design/doubao/demo/index.html
 M design/doubao/demo/styles.css
```

## 写面现量（`cmd/wisp` 此刻另腿状态：`git status --porcelain -- cmd/wisp` 空）

- `wc -l cmd/wisp/subagent_selfapproval_197_test.go` = `710`
- `:109` HEAD 与工作树同读、逐字：
```
			TaskID: taskID, CorrelationID: taskID,
```
- 该用例名（`:194`）：`Test197SubagentSelfApprovalIsRefusedAndTheHostAnswerLands`
- 载具调用点（为本次改动计）：`:94` 定义、`:389` `w1`、`:456` `w2`

## 周边事实（起手核查，只登记）

- `third_party/` 现量有 `sherpa-onnx`（无 `onnxruntime` 目录）；`cmd/wisp` 无 `TestMain`。
- 本腿写面：只 `cmd/wisp/subagent_selfapproval_197_test.go` ＋ 本目录 `*.md`（`00-anchor.md` 本件、后件一份读数）。

## 起手判语

待读现文后给出「既有断言依不依赖两者相等」的判语与凭据行；若依赖相等 ⇒ 停手上报（本件后件里记）。不预设结论。
