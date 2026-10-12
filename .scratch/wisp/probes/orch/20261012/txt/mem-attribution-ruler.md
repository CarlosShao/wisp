
## 10-12 补两枚小尺洞（同一族：归属尺与"分母行"）

在 10-11 那条"`tasklist //FI` 只认精确名，名册要有家族尺"之上，今天现跑抓到两枚：

- **测试进程的"专用根"⛔ 在 `CommandLine` 里，在 `ExecutablePath` 里。** 我用 `Get-CimInstance Win32_Process` 抽 `CommandLine` 去匹配 `wisp-test-gotmp/<run_root>` ⇒ `run_root=n/a`（GOTMPDIR 是**环境变量**，不进命令行）；而那枚根确实出现在**exe 自己的路径**里（`…\wisp-test-gotmp\<run>\…\b001\wisp.test.exe`）。⇒ 归属尺定式改成**两把并排**：`ExecutablePath` 抽根名 ＋ 父 PID 在不在。⚠ 派单里凡是写"用命令行里的具名根路径认归属"的条款，下一轮要顺手补这句，否则腿会照我原话跑出一把哑尺。
- **PowerShell 里单条对象上的 `.Count` 会打出空串**（`$rows = Get-CimInstance …; "$($rows.Count)"` ⇒ 空白），看起来像"0 枚"，其实那把尺**根本没数**。⇒ 一律 `$rows = @(Get-CimInstance …)` 再取 `.Count`。⚠ 这与我踩过的"空 stderr 被读成 0 命中"同形：**空⛔ 等于零**；分母行必须自带"数出来是几"的字面数字，跑完先看那一行是不是空。
