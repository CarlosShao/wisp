# 票 245 · 静态验收探针目录（`7f5f36b0d`）

本目录是 `docs/evidence/s1/245-static-accept.md` 的盘上附件，只建不删。

- 腿：`245-static-accept`（非实现者验收腿；只读）
- ⛔ 本腿**没有跑** `go test`／没有建窗／没有跑 bash 测试载具。下面所有文件都是 `git show`／`Select-String`／`grep`／`Get-FileHash` 的静态读数。
- 审的对象：`7f5f36b0d`（父 `137bf535c`）；审时 HEAD `307becc69`。

文件：

| 文件 | 内容 |
|---|---|
| `01-name-status.txt` | `git show 7f5f36b0d --name-status` 的 7 枚改动 |
| `02-skips.txt` | 本 commit 涉及文件里全部 `t.Skip*` 点与判词归类 |
| `03-registerhotkey-sites.txt` | 全部 `RegisterHotKey` 调用点与产品注册路径 |
| `04-citations.txt` | 证据件四条声称对应的 `file:line` 与盘上逐字 |

结论（与证据件同）：A 部分／B 通过／C 通过／D 通过。
