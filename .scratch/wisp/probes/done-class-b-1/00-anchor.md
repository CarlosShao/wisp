# done-class-b-1 起手锚（只读盘点腿）

腿：done-class-b-1。任务：逐张核乙类 8 张票（07/92/97/104/105/110/113/115）——已加 `-done` 却仍有未勾框：判 漏勾 / 真未做 / 判不动。⛔ 只交判语与缺什么，不动任何票面（不撤名、不改 Status、不补框）。⛔ 零 Go 命令（写腿 242-corrland-1 在飞）。写面＝只新建本目录 `*.md`。

## 三枚起手读数（原样登记）

### 1) `date '+%Y-%m-%d %H:%M:%S %z'`
```
2026-10-08 19:31:34 +0800
```

### 2) `git log -1 --format='%H %ad %s'`（全文逐字）
```
bbbc18fa23c663a79634cf67f18cc44e64b2711a Thu Oct 8 19:31:01 2026 +0800 A745 落账：收 242-corrcensus-1（9ae812d6/0e0ac7b9 四枚件）correlation 语义与读者普查。我核三处：loop.go:363 逐字 newTaskJournal(l.opt.Journal, taskID, taskID) // C18: correlation == task id；loop_approval_test.go:213-214 逐字 r.CorrelationID != res.TaskID … want the task id (C18)（唯一等值断言）；*_test.go 里 CorrelationID: taskID 恰 3 命中（pointer_183:127/pointer_185:120/ticket175r2:95）。采它读数（C18 原文半标仅腿报）：C18＝correlationId 是队列项字段+答复路由键、与任务标识并列未要求同值（PLAN.md:1368/SPEC-06:106）；决定行为读点 29/审计 13/标识 45；其它产码赋值点 22 组。★裁定方向＝应改（每次工具调用铸独立 corr；依据 C18 未要求同值+queue.go:163 corr 是主键⇒同任务并发两卡会塌成一枚）＋四边界（不动契约文本/不许放宽那枚断言须具名重判/loop.go:363 注释属引用与原文不符候选同批重判/三枚测试构造点按上下文重判不批量改），撤销口令「撤 242 改 corr 方向」。下一波 242-corrland-1（写腿铸 id+重判钉子注释+新用例同任务并发两问两 corr 不同）+非实现者验收（普查件不算凭据）。在飞 0 枚；零翻框零 push
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

注：`242-corrland-1` 写腿在飞（与 A745 口径一致），本腿全程零 Go 命令。`frontend/**`、`design/**` 不读不引（上面第 3 条仅作 `git status` 登记转抄，非阅读未越界）。
