# 114-v1 探针锚 — 票 114 验收腿（非实现者）

- 身份／时刻：`114-v1`，2026-10-14 16:0x–16:4x +08；起手 HEAD `9b51859b`；窗口内他腿 `20-r2`（internal/tools+config）在飞、`frontend/` 被界面侧持续写（109→113 枚）。
- 被裁：`114-a5` 产码 `d8dd4b38`（`internal/panel/composer_test.go` +146/−22、`bridge_test.go` +73）＋交付 `37eb6992`（`probes/114/a5/00-split-and-landing.md` 78 行＋logs）。
- 方法：仓内零突变；"改前"＝`go test -overlay=<仓外 json>` 换编译输入；种形＝仓外 `git archive HEAD` 树**真种盘**（A713：overlay 换不到 `os.ReadFile`）；全部命令含 `GOTMPDIR='C:\tmp\wisp114v1gotmp'`（反斜杠形）。
- 仓内锚定：突前/突后 `sha256sum` 同值 —— `composer_test.go` `d2b2e848ab12e4b29b31ee6a62964507404252c670fa6048e1ff0c89ea67edcf`；`bridge_test.go` `6ed0d1a41956bcc2e85eee5751f494a75667ec7a504fe6c601e661b59f16384a`。
- 仓外件（只建不删）：`D:\tmp\wisp114v1\mut\`（overlay json 4＋变异件 3＋88c7e752 对照件 2）、`D:\tmp\wisp114v1\ac8tree\`（页侧漂移种形）、`D:\tmp\wisp114v1\ac9tree\`（public 门种形）；logs/01…05 各含原命令行。
- 变异名册：m1＝`rendererSourceFile` 单删 `.json`；m2＝回退旧 7 类集；ac9-oldroot＝单行 root 回退 `frontend/src`；ac7-control＝两测试件回退 `88c7e752`；ac8 页侧 `"config.set"×2 ⇒ "config.sets"`。
- 交付：本目录 logs 5 枚＋`docs/evidence/s1/ticket114-114-v1-accept-v1.md`＋票 114 尾 1 bullet。⛔ 框未翻、未动 `Status:`/票名、零 push、只 commit。
