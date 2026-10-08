# 282-v1b 起手锚（腿名 282-v1b，非实现者复核）

- 起手时刻（`date`）：`Thu Oct  8 20:00:24 CST 2026`
- `git log -1`：`f5c63b7c15052d398830a7924f86cbc8dd2fd5c3 Thu Oct 8 19:57:59 2026 +0800 补票 282＋A748 落账：收 242-corrland-1（三笔 133b1bfa→0abe217c）correlation 落地。…`
- 分支：`dev`
- `git status --porcelain | wc -l` ＝ **758**（共享工作树，多腿在飞；只登记不动它）
- 前一条腿 `282-v1` 只留空目录 `.scratch/wisp/probes/282/v1/`，零提交；本腿从头做，那个空目录不删。
- 本腿写面＝AC#1 突变行（改完还原）＋新建 `.scratch/wisp/probes/282/v1b/*.md`。
- 票面：`.scratch/wisp/issues/282-corr-landing-over-scope-verdict-bridge-carries-taskid-and-the-new-accessor.md`
- 涉及三笔：`133b1bfa`／`bd124b2a`／`0abe217c`。

## 计划（票面五格）

1. AC#1 行为尺：种坏任务身份改读访问器那一跳（`TaskID(ctx)` 返回空或回退 corr），指名用例必须红；全绿则具名写"该支今天没有仪器"。四件套。
2. AC#2 逐枚核动过的断言是否等价或更强（三形皆红：空／旧形等值／无前缀）。
3. AC#3 深度判定前后各一发读数，判"静默失效"真修还是换写法。
4. AC#4 残余核实：journal 非生产组合行仍 taskID＋bridge.go CRLF。
5. AC#5 越界检查：三笔 `--name-only`，核冻结点面。
