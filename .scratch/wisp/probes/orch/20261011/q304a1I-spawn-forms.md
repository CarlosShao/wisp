# `Q-304-a1-I` 由编排者补的读数（18:0x，只读，⛔ 跑 Go）——非 `os/exec` 的直调起子进程形状

**为什么我现在跑这一发**：`304-a1` 具名欠了 `-I`（它只打了 `os.StartProcess`＝0 命中，没打全）。
这一问对票 304 要紧，因为甲形（Job Object）只在**经过 `exec.Cmd` 起的那一支**才接得上；
若仓里还有绕过 `os/exec` 直接 `CreateProcess` 的测试，甲形就留了第二个泄漏源，`304-v1` 得知道。

尺＝`grep -rn --include='*.go' -E <名>`，射程＝`internal cmd tools scripts` 四枚目录，
分母现量＝该射程内 **`.go` 682 枚**。件＝`20261011/logs/orch-A852-spawn-nonexec2.txt`（`done_rc=0` 自落）。

| 形状 | 命中 | 判读 |
|---|---|---|
| `os.StartProcess` | **0** | 没有绕过 `os/exec` 的直调 |
| `syscall.StartProcess` | **0** | 同上 |
| `ShellExecute` | **0** | 没有 |
| `syscall.Exec` | **0** | 没有 |
| `CreateProcess` | **1** | ★那一枚是**注释行**（`cmd/wisp/slo_exit_os_156_windows_test.go:335` "handed over by CreateProcess"）⇒ 真直调 **0 枚** |

⇒ **结论（给 `304-v1`）**：本仓测试起子进程**全部走 `os/exec`**，甲形两枚 spawn 点同批改完就没有第三支泄漏源；
`-I` 这一问**销掉**，⛔ 再欠。

另有两枚顺手量到的数（⚠ 别和上面混，也⛔ 拿来当"谁漏了"的证据）：

- 同一把尺（`-E 'exec\.Command'`）在**我这四枚目录射程**内命中 **74 行**：其中 `exec.CommandContext` **3 行**（前缀会被 `exec\.Command` 一起吸进来），
  其余 **71 行**；逐行按"这行是不是注释"复量＝**注释 0 行**（尺＝`sed` 取出冒号后那段再 `grep -cE '^[[:space:]]*(//|/\*)'`）。
- 换到**全仓 tracked `.go`**（分母 997 枚）＝**89 行**，其中 **15 行在 `.scratch/wisp/probes/**` 里**＝别家留下的变异拷贝样本，⛔ 是活测试。
  ⇒ 74 ＋ 15 ＝ 89 对上。
- ⚠ 与 `304-a1` 表①那 **72 枚 spawn 语句**差 **2**：这是**两把不同口径的尺**（它数"语句"、我数"命中行"，且我的模式把 `CommandContext` 也算进前缀），
  **成因我此刻没量，⛔ 在这里解释**。要作差请先由 `304-v1` 把两把尺的名册逐枚对上再说话。


⚠ **我这一发的前一发是坏的，留在盘上具名**：`logs/orch-A852-spawn-nonexec.txt` 里 `exec\.Command(` 那行写着 `hits=0`＝**假 0**。
成因＝我把裸 `(` 放进 ERE，它开了一个没闭合的组 ⇒ `grep` 报错，而那句被我 `2>/dev/null` 吞了 ⇒ 管道里只剩 `wc` 数到的 0 行。
⇒ 定式又长一节：**凡一把尺报 0，先看它的 stderr**（这发起我把 `grep_stderr=` 并排落进每一行，⛔ 再拿"空"当"没有"）。
