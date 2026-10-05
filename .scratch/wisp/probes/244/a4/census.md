# 244-a4 普查 —— 首启回执文案的「stderr 真到达」通道（只读普查，不产码）

> 腿号 `244-a4`。工作目录 `D:\work\workspace\projects plans\Wisp`。
> 起手锚点 `cfd97636`（`git log -1`，2026-10-05 12:58 +08:00）；本文件起手时刻 2026-10-05 13:01 +0800。
> 硬纪律：**一枚 `go` 命令都没跑**（`go test`／`go build`／`go vet`／`go env` 全部为 0 次调用）；
> 写点只有本文件；`grep`／`find` 的根一律显式限定为 `cmd internal tools scripts docs .scratch`。
> 本文件是**骨架 + 逐节读数**：下面是 §0–§7 的骨架，每节的正文由后续 Edit 逐节补齐并逐节 commit。

---

## §0 起手复认（派单给的每一条 `file:line` 断言的现读结果）

未判。

| 派单断言 | 现读 | 判语 |
|---|---|---|
| 未判 | 未判 | 未判 |

---

## §1 问 1 —— 现成先例：`cmd/wisp` 里「编译一次 wisp.exe 再用 `exec.Command` 起真子进程」那族

未判。

---

## §2 问 2 —— 通道绑定的真身：`os.Stdout`／`os.Stderr` 交给了谁

未判。

---

## §3 问 3 —— subsystem 那一格：`-H=windowsgui` 在 `scripts/` 的哪一行、PE subsystem 真身

未判。

---

## §4 问 4 —— 无父控制台那一支：`console_windows.go` 的 return 与 `logsink.go` 的 mirror 谁先谁后

未判。

---

## §5 问 5 —— 撞钉预检：新增一枚 exec 级用例会撞哪些既有计数尺

未判。

---

## §6 问 6 —— CI 那一侧：windows job 跑哪些包、什么 tag／env

未判。

---

## §7 问 7 —— 最小落地形状与代价：一发行文＋一枚用例

未判。

---

## §8 落盘尺与「⛔ 这轮没动的东西」自证

未判。
