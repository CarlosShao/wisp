# 票 283 / 腿 283-r1 起手锚

- 腿名：283-r1（窄写腿，写面＝`internal/tools/**`＋本目录 `*.md`）
- `date`：Fri Oct  9 08:39:49 CST 2026
- `git log -1 --format='%H %ad %s'`：
  `c39e2853b4946fa2c8e3422dedf78b3639c34a57 Thu Oct 8 20:20:03 2026 +0800 补票 283＋A749 落账：收 282-v1b（6c90fa7f/369a8f29）越格块非实现者复核——五格：AC#1 成立（种坏 TaskID→空 红 internal/tools/subagent_197_test.go:302；但种坏 →h.corr 三绿＝回退形无仪器）／AC#2 成立（新钉三形皆红：等值/无前缀红 loop_approval_test.go:221、空形红 :204 ⇒ 旧句拦的形仍红未放宽）⇒『没变松』这一问由非实现者答了／AC#3 判不动（身份改读前后同一形状三枚用例双双全绿，缺真 loop per-call corr→bridge 整链尺）／AC#4 成立两条残余均判留（journal 只在非生产组合＋bridge CRLF 是检出态）／AC#5 成立零碰禁列。我复验五枚 hash 全对（cancel 87ac0624/loop 0d53277d/197 8c9266a7/task 3ec8d498/bridge d46124de；⚠它把 197 路径写成 internal/agent，真身 internal/tools/subagent_197.go，hash 真）。⇒ 具名两处新盲区⇒补票 283（3框：查重→两处各补尺或登记不测→不许放宽）。在飞 0 枚；零翻框零 push`
- `git status --porcelain | head -20`（只登记，勿动）：
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
- 备注：`design/**` 按本腿纪律不读不引；以上 M/D 均非本腿所写，只作登记。
