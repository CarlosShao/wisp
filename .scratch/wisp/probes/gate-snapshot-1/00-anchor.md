# gate-snapshot-1 · 00 起手锚（只登记）

- 腿：gate-snapshot-1（只读门禁腿）。本轮 ⛔ 不跑 build/test/vet、⛔ 不改任何既有文件；唯一写面＝本目录 `*.md`。
- 锚时刻（本地）：2026-10-08 18:56:27 +0800
- 分支：`dev`
- HEAD（原文）：`0f7d4c53c5644423a2c3edcc07409f2b2dad7e1e 2026-10-08 18:54 票 259 四格翻勾＋A734 落账（凭据＝非实现者腿 259-v1 终裁）：AC#1 成立（ⓐ 已落地、改前句 cannot be replayed onto a different request 已撤、编译树零能挡跨卡、四把钉尺 rc=1）／AC#2 成立（approval.go:558 spend 现回 grantDenial，编排者现取复认；四因各 label+审计行+指名用例 12/12；对外合并句 :250 仍逐字钉）／AC#3 成立（两枚尺件 6+6 枚在库 ticket259_denial_rulers_test.go＋ticket259_panel_capability_rulers_test.go；三发突变 Mu-1 真花令牌红 :73/:82/:101/:111、Mu-2 Permitted bool 红 :174、Mu-3 Grant string+活 nonce 红 :199/:229/:301/:308；hash 逐字回 porcelain 空复绿 12/12）／AC#5 成立（15 笔提交对禁列零命中 rc=1 正控命中 docs/specs）。AC#4 不成立＝cmd/wisp/subagent_selfapproval_197_test.go:109 仍 TaskID: taskID, CorrelationID: taskID,（我复认）⇒ 留未勾并派 259-r4。票面 59→69 行/未勾 5→1/已勾 1→5，撤销口令逐格；⛔ 原句未改只追加 §11；它没试出＝残余风险两发未种（登记可选攻法）＋Mu-1/Mu-3 首种 sed 双命中 build failed 自认；在飞＝card-proof-prep-1；新派 259-r4；零 push`

## `git status --porcelain | head -20`（只登记，不解读；第 21+ 行未登记）

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

（含 `design/` 路径行系 `git status` 原样登记，本腿不读不引其内容、不深入。）

## 环境事实（本腿自检）

- `.scratch/wisp/probes/gate-snapshot-1/` 开跑时不存在 ⇒ 本腿新建，只放 `.md`。
- `git check-ignore -v .scratch/wisp/probes/gate-snapshot-1/00-anchor.md` → 无输出，rc=1 ⇒ 未被忽略，可 add。
- `/tmp` 可写（临时件只建不删、建仓库外）。

## 下一步

五把逐把测量（d22scan／gofmt -l／名册两尺／三面状态／未推计数）→ 写入 `01-snapshot.md`；全程零修改（除本目录 .md）。
