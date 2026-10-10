# 298-r1 — 开工闸门锚件（00-anchor）

钟点尺＝`date '+%Y-%m-%d %H:%M:%S %z'` 的 stdout：`2026-10-10 16:14:59 +0800`
工作目录＝`/d/work/workspace/projects plans/Wisp`（win32 / git bash）

## G1 双锚起手

尺＝`git log -1 --format='%h %ad' --date=iso-strict`

    6414a4bb 2026-10-10T16:10:11+08:00
    rc=0

## G2 工作树名册（逐行留档；本腿⛔认领其中任何一枚非本票件）

尺＝`git status --porcelain -- cmd internal scripts .github docs frontend`

    （无输出）
    rc=0 lines=0

⇒ 本票射程的五个目录（`cmd`／`internal`／`scripts`／`.github`／`docs`）＋`frontend` 在**跟踪状态**上全干净。

尺＝`git status --porcelain`（全仓，只作背景留档，⛔ 本腿不动任何一枚）：脏件集中在
`.gitignore`、`.scratch/wisp/probes/**`（前几轮的日志与证据件）、`design/**`（` D` 删除态与 ` M` 改动态若干枚），
以及 `.scratch/**` 下大量 `??` 未跟踪件（commit-msg 草稿、ci-logs 等）。
⇒ **`cmd/wisp` 不在其中任何一行**（下面 G3 是本票射程的专门判据，这条才是权威）。

## G3 停手条件（票面 AC#1 逐字执行的起手量）

尺＝`git status --porcelain -- cmd/wisp/panel_inbound_guards_35r3_test.go cmd/wisp/panel_transport_35r2_test.go`

    （无输出）
    rc=0 lines=0

⇒ **停手条件未触发**：两枚目标件在工作树里既非 ` M` 也非 `MM`，没有人正改着它们 ⇒ 本票放行。
另证：`diff <(git show HEAD:cmd/wisp/<两枚>) cmd/wisp/<两枚>` 两枚都给 `IDENTICAL` ⇒ 工作树内容＝HEAD blob，
本腿 `gofmt -w` 落下去覆盖的是 blob 本身，⛔ 不存在吞别人未提交改动的可能。

## G4 进程现量（只记录，⛔ 杀任何进程）

尺＝`tasklist //FI "IMAGENAME eq wisp.exe"` ⇒ `INFO: No tasks are running which match the specified criteria.`（0 枚）
尺＝`tasklist //FI "IMAGENAME eq balldebug.exe"` ⇒ `INFO: No tasks are running which match the specified criteria.`（0 枚）

## G5 改前基线（⚠ 本节先落锚笔，长跑基线随后落 `logs/`，但全部在 `gofmt -w` 之前）

### G5.1 尺＝`gofmt -l cmd/wisp`（工作树那一把）＝3 枚

    cmd\wisp\models.go
    cmd\wisp\panel_inbound_guards_35r3_test.go
    cmd\wisp\panel_transport_35r2_test.go
    rc=0 count=3

### G5.2 尺＝HEAD blob 那一把＝2 枚（落点＝仓外临时目录，⛔ 不落仓内）

命令链（逐字，可重跑；`$B` ＝ `mktemp -d` 给出来的仓外目录）：

    B=$(mktemp -d)                                   # 实测 /tmp/tmp.XnATepEG0K ⇒ C:\Users\swq\AppData\Local\Temp\tmp.XnATepEG0K
    mkdir -p "$B/blob"
    for p in cmd/wisp/models.go cmd/wisp/panel_inbound_guards_35r3_test.go cmd/wisp/panel_transport_35r2_test.go; do
      git show "HEAD:$p" > "$B/blob/$(basename $p)"
    done
    gofmt -l "$B/blob"

    C:\Users\swq\AppData\Local\Temp\tmp.XnATepEG0K\blob\panel_inbound_guards_35r3_test.go
    C:\Users\swq\AppData\Local\Temp\tmp.XnATepEG0K\blob\panel_transport_35r2_test.go
    rc=0 count=2

⇒ **票面现量复现**：blob 那一把只有 2 枚，`models.go` 不在列 ⇒ 落点＝那 2 枚；第三枚按 `AC#2` 只做标注、⛔ 不清洗。

### G5.3 换行符三把尺（票面现量节的三条，逐条重跑）

尺＝`git config core.autocrlf` ⇒ `true`

尺＝`git ls-files --eol cmd/wisp/models.go cmd/wisp/panel_inbound_guards_35r3_test.go cmd/wisp/panel_transport_35r2_test.go`

    i/lf    w/crlf  attr/text eol=lf      	cmd/wisp/models.go
    i/lf    w/lf    attr/text eol=lf      	cmd/wisp/panel_inbound_guards_35r3_test.go
    i/lf    w/lf    attr/text eol=lf      	cmd/wisp/panel_transport_35r2_test.go

尺＝`tr -cd '\r' | wc -c`（工作树件）

    models.go CR=334                     # ＝其行数，逐行 CRLF
    panel_inbound_guards_35r3_test.go CR=0
    panel_transport_35r2_test.go CR=0

## G6 引入首笔（AC#0 要求的 `git log --oneline -- <path>` 具名）

尺＝`git log --oneline -- cmd/wisp/panel_inbound_guards_35r3_test.go` ⇒ 名册只 1 笔，最老＝唯一＝**首笔引入**
`7f9d6e40`「35-r3 交件（票35 :63 (d) 三支入向守卫红，只写测试）」（票 35 的 `35-r3` 落地腿）

尺＝`git log --oneline -- cmd/wisp/panel_transport_35r2_test.go` ⇒ 3 笔（新→旧 `3a343bc7`／`2fc5f5c9`／`286a7f30`），
最老＝**首笔引入** `286a7f30`「35-r2 甲③ 落地＋判据换行为尺（票35 AC#6:52）…新增 cmd/wisp/panel_transport_35r2_test.go…」（票 35 的 `35-r2` 落地腿）

⇒ **两枚真债都是票 35 那两笔落地腿带进来的**，与票面标题的归因一致；细节见 `10-ac0-dirt-shape.md`。

## 射程声明（本腿自缚）

写面＝`cmd/wisp/panel_inbound_guards_35r3_test.go` ＋ `cmd/wisp/panel_transport_35r2_test.go` ＋
`.scratch/wisp/probes/298/r1/**`（只新建 `.md`／`.txt`）＋ 票 298 票面**追加**。
⛔ `cmd/wisp/models.go`（禁区节逐字：它没格式债）／⛔ `.gitattributes`／⛔ git 配置／⛔ 整包格式洗／⛔ 任何一行语义。
