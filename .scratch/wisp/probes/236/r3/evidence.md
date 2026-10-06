# 236-r3 — AC#2 的交件（六节）

腿＝`236-r3`；射程＝票 236 的 **AC#2 一格**（那把 DEFERRED 尺的 (b) 支改扫能力）。
⛔ 本件不答也不裁 AC#1／AC#1b／AC#3／AC#4／AC#5／AC#6（派单 §0 第 4 把尺），七枚框一枚不碰。

> **本节以下目前是骨架**：每节标题＋本腿打算答什么。读数逐节回填，占位词一个不留（写完才 commit）。

---

## §0 起手锚（四把尺原文）

打算答：`sed -n '3,5p'` 的撤票口令原文与 grep 计数、`git log -1`／`date`、`git status --porcelain -- internal cmd`（★必须证明起手 `internal/tools` 无未提交脏文件）、`grep -n '^- \[ \]'` 的逐条行号与枚数、`git status --porcelain -- scripts .github`。

## §1 现量（AC#2 那把尺的三点量＋分界）

打算答：①尺本体在哪枚文件哪几行、它到底在数什么（逐行读）；②`BuiltinTaskEntries` 注册名册今天有几枚、`task.list` 在不在（现跑）；③与票 225 的分界（逐字引票面＋225 标题，⛔ 不做 225 的活）；④词面被叙述句撑住的机制复认（`grep -n DEFERRED internal/tools/task.go` 行数、同时含 `task.list` 的行数）；⑤"零命中不是证据"——每把负向尺先证它命中得了真名。

## §2 判据表（跑了什么命令｜原始读数｜判语）

打算答：基线绿名册（用例枚数＋rc）、AC#2 判据逐条对应哪一枚用例、每把尺的口径。

## §3 突变名册（teeth 名册）

打算答：
- **M-2a**＝摘掉 `internal/tools/task.go:23` 那行 DEFERRED 标记（`teeth-m13` 那形）：**未修码读数**＝归档 leg `Test221DeferredMarkerForCancelLiftedButListStillMarked` 仍 PASS ⇒ 今天无牙；**装牙后读数**＝本腿新 leg 具名红。
- **M-2c**＝正控：overlay 替换测试文件里的读路径 ⇒ 归档 leg 变红，证明 overlay 打得进测试文件、也证明 M-2a 的绿是"读盘看不见"而不是"overlay 没落地"。
- **M-2b／M-2e**＝能力侧突变：`BuiltinTaskEntries` 名册里加／减 `task.list`，分别证新判据的**两个方向都有牙**、且合法同批核销（摘标记＋接线同时）**不假红**——这一发同时把票面 §四 说的计数形 `len(names) != 2` 的假红风险实测出来。
- 每发记：摘哪一行＋overlay json 路径＋落地证明＋红句原文＋还原后 `md5sum` 与 `git show HEAD:<file> | md5sum` 对拉。

## §4 门禁读数

打算答：`sh scripts/d22scan.sh`（含正控）、`gofmt -l`、`gofumpt -l` 对本腿新文件、`go vet ./internal/tools/`、`git status --porcelain -- internal cmd` 收尾为空。本腿若动 `scripts/**` 则补 `bash -n`；⛔ 本腿不动 `.github/**`。

## §5 判不动的地方

打算答（实质内容，不写"没有"）：AC#2 射程外但今天量到的东西——票 225 的双向对账、`task.list` 到底该不该注册（D34／§7 :1531，归 owner 与产码腿）、读盘尺在产码里的其它同类（若查到）、以及本腿判据"哪一维今天量不到"。
