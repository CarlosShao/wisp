# 281-r1 起手锚（账目归位腿）

腿：`281-r1`。任务：照票 `.scratch/wisp/issues/281-class-b-done-tickets-have-open-boxes-and-status-conflicts-reconcile-one-by-one.md` 四个 AC 逐格做——乙类 8 张（07/92/97/104/105/110/113/115：已 `-done` 却有未勾框／Status 与名相抵）逐张逐框归位。
纪律：⛔ 零 Go 命令（腿 `283-r1` 在 `internal/tools` 在飞）；只用 `git`/`grep`/`sed`/`wc`/`date`。写面＝那 8 张票面＋只新建本目录 `*.md`。⛔ `frontend/**`／`design/**` 不读不引；⛔ 凭据（密钥）不进对话。

## 三枚起手读数（原样登记）

### 1) `date`
```
Fri Oct  9 08:39:20 CST 2026
```

### 2) `git log -1 --format='%H %ad %s'`（逐字）
```
c39e2853b4946fa2c8e3422dedf78b3639c34a57 Thu Oct 8 20:20:03 2026 +0800 补票 283＋A749 落账：收 282-v1b（6c90fa7f/369a8f29）越格块非实现者复核——五格：AC#1 成立（种坏 TaskID→空 红 internal/tools/subagent_197_test.go:302；但种坏 →h.corr 三绿＝回退形无仪器）／AC#2 成立（新钉三形皆红：等值/无前缀红 loop_approval_test.go:221、空形红 :204 ⇒ 旧句拦的形仍红未放宽）⇒『没变松』这一问由非实现者答了／AC#3 判不动（身份改读前后同一形状三枚用例双双全绿，缺真 loop per-call corr→bridge 整链尺）／AC#4 成立两条残余均判留（journal 只在非生产组合＋bridge CRLF 是检出态）／AC#5 成立零碰禁列。我复验五枚 hash 全对（cancel 87ac0624/loop 0d53277d/197 8c9266a7/task 3ec8d498/bridge d46124de；⚠它把 197 路径写成 internal/agent，真身 internal/tools/subagent_197.go，hash 真）。⇒ 具名两处新盲区⇒补票 283（3框：查重→两处各补尺或登记不测→不许放宽）。在飞 0 枚；零翻框零 push
```

### 3) `git status --porcelain | head -20`（只登记，不解读）
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

注：上面第 3 条仅作 `git status` 转抄登记（第 13 行起 `design/**` 为共享工作树里他腿的在飞改动；本腿不读不引 `design/**`，仅登记行号不含其内容）。
