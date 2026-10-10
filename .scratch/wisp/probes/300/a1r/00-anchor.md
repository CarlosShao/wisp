# 300-a1r — `00` 锚与起手闸门（只读＋台件腿，零产码）

派单编号：`300-a1r`。腿种：**只读＋台件**（read-only + rig），本程**零产码**、⛔ 改任何跟踪文件。
票面：`.scratch/wisp/issues/300-parsewaveformat-subformat-offset-reads-two-bytes-past-the-guid.md`
本程射程：只新建 `.scratch/wisp/probes/300/a1r/**`（正文 `.md`，原始输出 `logs/*.txt`）＋票面 `Progress log` 追加一行。

## 起手闸门（顺序硬门，逐字读数）

命令（第一条，一字未改）：

```
cd "D:\work\workspace\projects plans\Wisp" && date "+%Y-%m-%d %H:%M:%S %z" && git rev-parse --short HEAD && git status --porcelain -- internal cmd | head
```

stdout 逐字：

```
2026-10-10 12:18:27 +0800
b2933dc9
```

- **锚＝`b2933dc9`**：与派单给的派单锚 **一致**，没有漂 ⇒ 本程全部读数钉在 `b2933dc9`。
- 第三段（`git status --porcelain -- internal cmd | head`）**输出 0 行**（尺＝该命令的 stdout 行数；`head` 未截断任何东西，因为压根没有行）。
  ⇒ 我开工时 `internal/**` 与 `cmd/**` 两棵子树**逐路径看是干净的**。
- ⚠ **口径写死**：整棵工作树**本来就脏约 805 枚**（`design/**`、`frontend/**` 有删除项，⛔ 我造成、⛔ 我复位）。
  所以本件及后续任何一件都**不许**写"工作树干净"；越界回报口径只能是
  "**我动过的路径恢复原状**"／"`git status --porcelain -- internal cmd` ＝ 0 行"。

## 环境（逐字读数，尺名随附）

同一条命令内现跑（第二条，锚跑完、落第 1 笔 commit 之前）：

```
git rev-parse --abbrev-ref HEAD && go version && go env GOOS GOARCH GOFLAGS
```

- 分支：`dev`
- `go version`：`go version go1.27.1 windows/amd64`
- `go env GOOS GOARCH GOFLAGS`：`windows` / `amd64` / （`GOFLAGS` **空行**，即未设）
  ⇒ 派单要求的 `GOFLAGS= go build ./...` 里那个显式空值是**把继承来的 GOFLAGS 清掉**，与本地默认等价，仍照写。
- 平台：Windows 10/11 + Git Bash（POSIX sh）。`/tmp` 是 MSYS 树下的可达路径。

## 落点目录与 `.out` 陷阱（开工前现读）

- `ls .scratch/wisp/probes/300` ⇒ **无输出、命令 rc=2**（目录不存在，本程新建 `300/a1r/`）。
  尺＝`ls` 的退出码＋stdout 行数；同级已存在 `111 114 132 139 145 145-accept 147 149 151 151-accept 152 153 153-ac2 154 154-accept 155 155-accept 156 156-accept 157`（`head -20` 截取的**前 20 行**，⛔ 当成全部）。
- `.gitignore` 第 **8** 行逐字：`*.out`（尺＝`sed -n '8p' .gitignore`；文件总 `wc -c`＝**1404 字节**）。
  ⇒ 本程**证据件⛔ 叫 `.out`**，否则会被静默跳过而 `commit` 仍回显成功（那格等于没交）。
  正文一律 `.md`，原始输出放 `logs/*.txt`。

## 本程要交的两格（票面原文口径，⛔ 我新造判据）

- `AC#0`＝**权威偏移**：`WAVEFORMATEXTENSIBLE.SubFormat` 相对结构起始到底是 **24 还是 26**。
  凭据第一选择＝**盘上 Windows SDK 头文件 `mmreg.h`** 逐字整块抄 typedef；盘上没有再取 Microsoft Learn。
  ⛔ 拿本仓另一枚实现当裁判（`wavinjector.go` / `wasapi_windows.go` 都⛔ 当尺），⛔ 凭记忆。
- `AC#1`＝**决定性一发**：`git archive` 把 `b2933dc9` 导出到**仓外**一棵树，在**那份导出树**里新建
  `package audio` + `//go:build windows` 的包内用例，自造 40 字节 `WAVEFORMATEXTENSIBLE` 字节面喂 `parseWaveFormat`，
  断言**解析出来的 `tag` 值 == 1／3**（⛔ 断"被调用过"）。
  ★**恒红用例⛔ 入库**（票面 12:2x「编排者补格」第一节已写死），真正的跟踪测试随 `AC#2` 同批入库。
  同一棵导出树里还要一发**正控**（坏形状⛔ 让用例绿），每形与颜色配对交回。
- ★**硬门**（票面 `AC#0` 末段）：如果权威结论是 **26 才对**，**停手上报**交回编排者裁，⛔ 自己改夹具布局再跑一遍。

## 与派单／票面的一致性（先声明，读到权威再判）

- 派单与票面在本次读到的范围内**没有冲突**：票面 `AC#1` 原文写"在 `internal/audio` 新增一枚 windows-tagged 用例"，
  而票面 12:2x 补格第一节**已由编排者自己把射程改成仓外导出树**，与本派单同形 ⇒ 我照仓外做，⛔ 具名"冲突"。
- 派单里那句 `cbSize=22`／总长 40／`SubFormat@24` **是编排者推的、⛔ 是读数**（票面 `:13` 逐字写着"⚠算术这一半是我推的"）。
  本件把它登记为**待 `AC#0` 裁的假设**，⛔ 当结论用。
